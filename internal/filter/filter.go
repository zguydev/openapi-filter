// Package filter provides functionality to filter OpenAPI specs based
// on config. It allows filtering of paths, methods, components,
// and other OpenAPI elements while maintaining the integrity of the spec.
package filter

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	"github.com/zguydev/openapi-filter/internal/config"
)

// OpenAPISpecFilter is the main type that handles filtering of OpenAPI specs.
type OpenAPISpecFilter struct {
	cfg    *config.FilterConfig
	logger *zap.Logger
	loader *openapi3.Loader

	doc, filtered *openapi3.T
}

// NewOpenAPISpecFilter creates a new OpenAPISpecFilter instance with the
// provided configuration and logger. It initializes the OpenAPI loader.
func NewOpenAPISpecFilter(
	cfg *config.Config,
	logger *zap.Logger,
) *OpenAPISpecFilter {
	loader := openapi3.NewLoader()
	if cfg.Tool.Loader != nil {
		loader.IsExternalRefsAllowed = cfg.Tool.Loader.IsExternalRefsAllowed
	}

	return &OpenAPISpecFilter{
		cfg:    &cfg.FilterConfig,
		logger: logger,
		loader: loader,
	}
}

// Filter processes an OpenAPI specification file according to the configured 
// filters and writes the filtered result to the specified output path.
// It handles loading, filtering, and writing of the specification while
// maintaining all necessary references and components.
//
// Parameters:
//   - inputSpecPath: Path to the input OpenAPI specification file
//   - outSpecPath: Path where the filtered specification will be written
//
// Returns an error if any step of the filtering process fails.
func (oaf *OpenAPISpecFilter) Filter(inputSpecPath, outSpecPath string) error {
	var err error
	oaf.doc, err = loadSpecFromFile(oaf.loader, inputSpecPath)
	if err != nil {
		oaf.logger.Error("failed to load spec from file", zap.Error(err))
		return fmt.Errorf("loadSpecFromFile: %w", err)
	}

	oaf.filtered = &openapi3.T{
		OpenAPI:    oaf.doc.OpenAPI,
		Components: &openapi3.Components{},
		Info:       oaf.doc.Info,
		Paths:      &openapi3.Paths{},
	}

	collector := NewRefsCollector()

	oaf.filterPaths(collector)
	oaf.filterComponents()
	oaf.filterOther()
	if isEmptyComponents(oaf.filtered.Components) {
		oaf.filtered.Components = nil
	}

	if err := writeSpecToFile(oaf.filtered, outSpecPath); err != nil {
		oaf.logger.Error("failed to write output spec file", zap.Error(err))
		return fmt.Errorf("writeSpecToFile: %w", err)
	}
	oaf.logger.Info("filtered and saved spec",
		zap.String("output", outSpecPath))
	return nil
}

// filterPaths processes the paths specified in the configuration and filters them
// according to the allowed methods. It also collects all references used in the
// filtered paths and processes them.
func (oaf *OpenAPISpecFilter) filterPaths(collector *RefsCollector) {
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
			collector.CollectOperation(op)
		}

		oaf.filtered.Paths.Set(path, newPathItem)
	}

	oaf.filterRefs(collector.Refs())
}

// getOperation safely retrieves an operation from a PathItem for the specified
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

// setOperation safely sets an operation in a PathItem for the specified method.
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
// included in the filtered specification.
func (oaf *OpenAPISpecFilter) filterRefs(refs map[string]struct{}) {
	for ref := range refs {
		oaf.filterRef(ref)
	}
}

// filterRef processes a single reference and copies the referenced component
// to the filtered specification if it exists in the original.
func (oaf *OpenAPISpecFilter) filterRef(ref string) {
	if oaf.doc.Components == nil {
		return
	}

	def, name, ok := parseRef(ref)
	if !ok {
		oaf.logger.Warn("incorrect ref", zap.String("ref", ref))
		return
	}

	compType, ok := ComponentDefToType(def)
	if !ok {
		oaf.logger.Warn("unknown component definition",
			zap.String("def", def),
			zap.String("name", name),
			zap.String("ref", ref))
		return
	}
	if !processCopyComponent(oaf.doc.Components,
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
// copies them to the filtered specification if they exist in the original.
func (oaf *OpenAPISpecFilter) filterComponents() {
	if oaf.cfg.Components == nil || oaf.doc.Components == nil {
		return
	}

	for _, compTyp := range ComponentTypes() {
		for _, name := range ComponentTypeToCfgNames(oaf.cfg.Components, compTyp) {
			if !processCopyComponent(oaf.doc.Components,
				oaf.filtered.Components,
				compTyp,
				name,
			) {
				oaf.logger.Warn("component not found",
					zap.String("def", ComponentTypeToDef(compTyp)),
					zap.String("name", name))
			}
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
