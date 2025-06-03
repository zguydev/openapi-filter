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

func (oaf *OpenAPISpecFilter) setOperation(p *openapi3.PathItem, method, path string, operation *openapi3.Operation) (ok bool) {
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
	var (
		dc = oaf.doc.Components
		fc = oaf.filtered.Components
	)
	def, name := parseRef(ref)

	switch def {
	case "schemas":
		initializeDef(&fc.Schemas)
		copyComponent(oaf.logger, dc.Schemas, fc.Schemas, name, ref)
	case "parameters":
		initializeDef(&fc.Parameters)
		copyComponent(oaf.logger, dc.Parameters, fc.Parameters, name, ref)
	case "headers":
		initializeDef(&fc.Headers)
		copyComponent(oaf.logger, dc.Headers, fc.Headers, name, ref)
	case "requestBodies":
		initializeDef(&fc.RequestBodies)
		copyComponent(oaf.logger, dc.RequestBodies, fc.RequestBodies, name, ref)
	case "responses":
		initializeDef(&fc.Responses)
		copyComponent(oaf.logger, dc.Responses, fc.Responses, name, ref)
	case "securitySchemes":
		initializeDef(&fc.SecuritySchemes)
		copyComponent(oaf.logger, dc.SecuritySchemes, fc.SecuritySchemes, name, ref)
	case "examples":
		initializeDef(&fc.Examples)
		copyComponent(oaf.logger, dc.Examples, fc.Examples, name, ref)
	case "links":
		initializeDef(&fc.Links)
		copyComponent(oaf.logger, dc.Links, fc.Links, name, ref)
	case "callbacks":
		initializeDef(&fc.Callbacks)
		copyComponent(oaf.logger, dc.Callbacks, fc.Callbacks, name, ref)
	default:
		oaf.logger.Warn("unhandled $ref", zap.String("ref", ref))
	}
}

func (oaf *OpenAPISpecFilter) filterComponents() {
	componentsCfg := oaf.cfg.Components
	if componentsCfg == nil {
		return
	}
	var (
		dc = oaf.doc.Components
		fc = oaf.filtered.Components
	)

	for _, name := range componentsCfg.Schemas {
		initializeDef(&fc.Schemas)
		copyComponent(oaf.logger, dc.Schemas, fc.Schemas, name, refName("schemas", name))
	}
	for _, name := range componentsCfg.Parameters {
		initializeDef(&fc.Parameters)
		copyComponent(oaf.logger, dc.Parameters, fc.Parameters, name, refName("parameters", name))
	}
	for _, name := range componentsCfg.SecuritySchemes {
		initializeDef(&fc.SecuritySchemes)
		copyComponent(oaf.logger, dc.SecuritySchemes, fc.SecuritySchemes, name, refName("securitySchemes", name))
	}
	for _, name := range componentsCfg.RequestBodies {
		initializeDef(&fc.RequestBodies)
		copyComponent(oaf.logger, dc.RequestBodies, fc.RequestBodies, name, refName("requestBodies", name))
	}
	for _, name := range componentsCfg.Responses {
		initializeDef(&fc.Responses)
		copyComponent(oaf.logger, dc.Responses, fc.Responses, name, refName("responses", name))
	}
	for _, name := range componentsCfg.Headers {
		initializeDef(&fc.Headers)
		copyComponent(oaf.logger, dc.Headers, fc.Headers, name, refName("headers", name))
	}
	for _, name := range componentsCfg.Examples {
		initializeDef(&fc.Examples)
		copyComponent(oaf.logger, dc.Examples, fc.Examples, name, refName("examples", name))
	}
	for _, name := range componentsCfg.Links {
		initializeDef(&fc.Links)
		copyComponent(oaf.logger, dc.Links, fc.Links, name, refName("links", name))
	}
	for _, name := range componentsCfg.Callbacks {
		initializeDef(&fc.Callbacks)
		copyComponent(oaf.logger, dc.Callbacks, fc.Callbacks, name, refName("callbacks", name))
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
