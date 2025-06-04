package filter

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	"github.com/zguydev/openapi-filter/internal/config"
)

type OpenAPISpecFilter struct {
	cfg    *config.FilterConfig
	logger *zap.Logger
	loader *openapi3.Loader

	filtered, doc *openapi3.T
}

func NewOpenAPISpecFilter(cfg *config.Config, logger *zap.Logger) *OpenAPISpecFilter {
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
					zap.String("path", path), zap.String("method", method))
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

func (oaf *OpenAPISpecFilter) getOperation(p *openapi3.PathItem, method, path string) (op *openapi3.Operation) {
	defer func() {
		if r := recover(); r != nil {
			oaf.logger.Warn("unknown HTTP method in filter config",
				zap.String("path", path), zap.String("method", method))
			op = nil
		}
	}()
	return p.GetOperation(strings.ToUpper(method))
}

func (oaf *OpenAPISpecFilter) setOperation(
	p *openapi3.PathItem,
	method, path string,
	operation *openapi3.Operation,
) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			oaf.logger.Warn("unknown HTTP method in spec",
				zap.String("path", path), zap.String("method", method))
			ok = false
		}
	}()
	p.SetOperation(strings.ToUpper(method), operation)
	return true
}

func (oaf *OpenAPISpecFilter) filterRefs(refs map[string]struct{}) {
	for ref := range refs {
		oaf.filterRef(ref)
	}
}

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
			zap.String("def", def), zap.String("ref", ref))
		return
	}
	if !processCopyComponent(
		oaf.doc.Components,
		oaf.filtered.Components,
		compType, name) {
		oaf.logger.Warn("component not found",
			zap.String("name", name),
			zap.String("def", def),
			zap.String("ref", ref))
	}
}

func (oaf *OpenAPISpecFilter) filterComponents() {
	if oaf.cfg.Components == nil || oaf.doc.Components == nil {
		return
	}

	for _, compTyp := range ComponentTypes() {
		for _, name := range ComponentTypeToCfgNames(oaf.cfg.Components, compTyp) {
			if !processCopyComponent(
				oaf.doc.Components,
				oaf.filtered.Components,
				compTyp, name) {
				oaf.logger.Warn("component not found",
					zap.String("name", name),
					zap.String("def", ComponentTypeToDef(compTyp)))
			}
		}
	}
}

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
