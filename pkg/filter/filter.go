// Package filter provides functionality to filter OpenAPI specs based
// on config. It allows filtering of paths, methods, components,
// and other OpenAPI elements while maintaining the integrity of the spec.
package filter

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	"github.com/zguydev/openapi-filter/internal/components"
	"github.com/zguydev/openapi-filter/internal/refs"
	"github.com/zguydev/openapi-filter/pkg/config"
)

// OpenAPISpecFilter is the main type that handles filtering of OpenAPI specs.
type OpenAPISpecFilter struct {
	cfg       *config.FilterConfig
	logger    *zap.Logger
	collector *refs.RefsCollector

	doc, filtered *openapi3.T
}

// NewOpenAPISpecFilter creates a new OpenAPISpecFilter instance with the
// provided configuration and logger.
func NewOpenAPISpecFilter(
	cfg *config.Config,
	logger *zap.Logger,
) *OpenAPISpecFilter {
	return &OpenAPISpecFilter{
		cfg:       &cfg.FilterConfig,
		logger:    logger,
		collector: refs.NewRefsCollector(),
	}
}

// Filter processes an OpenAPI spec according to the configured
// filters and returns a filtered spec.
// Returns an error if any step of the filtering process fails.
func (oaf *OpenAPISpecFilter) Filter(doc *openapi3.T) (filtered *openapi3.T, err error) {
	oaf.doc = doc

	oaf.filtered = &openapi3.T{
		OpenAPI:    oaf.doc.OpenAPI,
		Components: &openapi3.Components{},
		Info:       oaf.doc.Info,
		Paths:      &openapi3.Paths{},
	}

	oaf.filterPaths()
	oaf.filterComponents()
	oaf.filterOther()
	oaf.filterRefs()
	if components.IsEmptyComponents(oaf.filtered.Components) {
		oaf.filtered.Components = nil
	}
	return oaf.filtered, nil
}

// filterPaths processes the paths specified in the configuration and filters them
// according to the allowed methods. It also collects all references used in the
// filtered paths.
func (oaf *OpenAPISpecFilter) filterPaths() {
	for path, methods := range oaf.cfg.Paths {
		pathItem := oaf.doc.Paths.Find(path)
		if pathItem == nil {
			oaf.logger.Warn("path not found in spec", zap.String("path", path))
			continue
		}

		newPathItem := &openapi3.PathItem{}
		for _, method := range methods {
			op := oaf.getOperation(pathItem, method, path)
			if op == nil {
				oaf.logger.Warn("method not exists for specified path",
					zap.String("method", method),
					zap.String("path", path))
				continue
			}
			if !oaf.setOperation(newPathItem, method, path, op) {
				continue
			}
			oaf.collector.CollectOperation(op)
		}

		oaf.filtered.Paths.Set(path, newPathItem)
	}
}

// getOperation safely retrieves an operation from [openapi3.PathItem] for the specified
// method. It handles unknown HTTP methods gracefully and returns nil if the
// method is invalid.
func (oaf *OpenAPISpecFilter) getOperation(
	p *openapi3.PathItem,
	method, path string,
) (op *openapi3.Operation) {
	defer func() {
		if r := recover(); r != nil {
			oaf.logger.Warn("unknown HTTP method in filter config",
				zap.String("method", method),
				zap.String("path", path))
			op = nil
		}
	}()
	return p.GetOperation(strings.ToUpper(method))
}

// setOperation safely sets an operation in [openapi3.PathItem] for the specified method.
// It handles unknown HTTP methods gracefully and returns false if the method
// is invalid.
func (oaf *OpenAPISpecFilter) setOperation(
	p *openapi3.PathItem,
	method, path string,
	operation *openapi3.Operation,
) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			oaf.logger.Warn("unknown HTTP method in spec",
				zap.String("method", method),
				zap.String("path", path))
			ok = false
		}
	}()
	p.SetOperation(strings.ToUpper(method), operation)
	return true
}

// filterRefs processes all collected references and ensures they are properly
// included in the filtered spec.
func (oaf *OpenAPISpecFilter) filterRefs() {
	for ref := range oaf.collector.Refs() {
		oaf.filterRef(ref)
	}
}

// filterRef processes a single reference and copies the referenced component
// to the filtered spec.
func (oaf *OpenAPISpecFilter) filterRef(ref string) {
	if oaf.doc.Components == nil {
		return
	}

	def, name, ok := refs.ParseRef(ref)
	if !ok {
		oaf.logger.Warn("incorrect ref", zap.String("ref", ref))
		return
	}

	compType, ok := components.ComponentDefToType(def)
	if !ok {
		oaf.logger.Warn("unknown component definition",
			zap.String("def", def),
			zap.String("name", name),
			zap.String("ref", ref))
		return
	}
	if !components.ProcessCopyComponent(
		oaf.doc.Components,
		oaf.filtered.Components,
		compType,
		name,
	) {
		oaf.logger.Warn("component not found",
			zap.String("def", def),
			zap.String("name", name),
			zap.String("ref", ref))
	}
}

// filterComponents processes all components specified in the configuration and
// copies them to the filtered spec.
func (oaf *OpenAPISpecFilter) filterComponents() {
	if oaf.cfg.Components == nil || oaf.doc.Components == nil {
		return
	}

	for _, compTyp := range components.ComponentTypes() {
		for _, name := range components.ComponentTypeToCfgNames(oaf.cfg.Components, compTyp) {
			if !components.ProcessCopyComponent(
				oaf.doc.Components,
				oaf.filtered.Components,
				compTyp,
				name,
			) {
				oaf.logger.Warn("component not found",
					zap.String("def", components.ComponentTypeToDef(compTyp)),
					zap.String("name", name))
				continue
			}
			oaf.collector.CollectComponent(oaf.doc.Components, compTyp, name)
		}
	}
}

// filterOther processes additional OpenAPI elements specified in the configuration,
// including servers, security requirements, tags, and external documentation.
func (oaf *OpenAPISpecFilter) filterOther() {
	if oaf.cfg.Servers {
		oaf.filtered.Servers = oaf.doc.Servers
	}
	if oaf.cfg.Security {
		oaf.filtered.Security = oaf.doc.Security
	}
	if oaf.cfg.Tags {
		oaf.filtered.Tags = oaf.doc.Tags
	}
	if oaf.cfg.ExternalDocs {
		oaf.filtered.ExternalDocs = oaf.doc.ExternalDocs
	}
}
