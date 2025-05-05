package filter

import (
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/zguydev/openapi-filter/internal/config"
)

type OpenAPISpecFilter struct {
	cfg    *config.FilterConfig
	logger *zap.Logger
	loader *openapi3.Loader
}

func NewOpenAPISpecFilter(cfg *config.FilterConfig, logger *zap.Logger) *OpenAPISpecFilter {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	return &OpenAPISpecFilter{
		cfg:    cfg,
		logger: logger,
		loader: loader,
	}
}

func (o *OpenAPISpecFilter) Filter(inputSpecPath, outSpecPath string) error {
	inputFileData, err := os.ReadFile(inputSpecPath)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	doc, err := o.loader.LoadFromData(inputFileData)
	if err != nil {
		o.logger.Error("failed to load OpenAPI spec",
			zap.Error(fmt.Errorf("loader.LoadFromData: %w", err)))
		return fmt.Errorf("loader.LoadFromData: %w", err)
	}

	filtered := &openapi3.T{
		OpenAPI:      doc.OpenAPI,
		Components:   &openapi3.Components{},
		Info:         doc.Info,
		Paths:        &openapi3.Paths{},
		Security:     doc.Security,
		Servers:      doc.Servers,
		Tags:         doc.Tags,
		ExternalDocs: doc.ExternalDocs,
	}

	usedRefs := map[string]struct{}{}
	for path, methods := range o.cfg.Paths {
		originalPathItem, exists := doc.Paths.Map()[path]
		if !exists {
			o.logger.Warn("path not found in spec", zap.String("path", path))
			continue
		}

		newPathItem := &openapi3.PathItem{}
		for _, method := range methods {
			var op *openapi3.Operation
			switch strings.ToLower(method) {
			case "connect":
				op = originalPathItem.Get
				newPathItem.Get = op
			case "delete":
				op = originalPathItem.Delete
				newPathItem.Delete = op
			case "get":
				op = originalPathItem.Get
				newPathItem.Get = op
			case "head":
				op = originalPathItem.Head
				newPathItem.Head = op
			case "options":
				op = originalPathItem.Options
				newPathItem.Options = op
			case "patch":
				op = originalPathItem.Patch
				newPathItem.Patch = op
			case "post":
				op = originalPathItem.Post
				newPathItem.Post = op
			case "put":
				op = originalPathItem.Put
				newPathItem.Put = op
			case "trace":
				op = originalPathItem.Trace
				newPathItem.Trace = op
			default:
				o.logger.Warn("unknown HTTP method", zap.String("method", method))
				continue
			}

			if op == nil {
				o.logger.Warn("method not exists for specified path",
					zap.String("path", path), zap.String("method", method))
				continue
			}
			o.collectOperationRefs(op, usedRefs)
		}
		filtered.Paths.Set(path, newPathItem)
	}

	for ref := range usedRefs {
		switch {
		case strings.HasPrefix(ref, "#/components/schemas/"):
			name := strings.TrimPrefix(ref, "#/components/schemas/")
			if filtered.Components.Schemas == nil {
				filtered.Components.Schemas = map[string]*openapi3.SchemaRef{}
			}
			if s, ok := doc.Components.Schemas[name]; ok {
				filtered.Components.Schemas[name] = s
			} else {
				o.logger.Warn("schema component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/parameters/"):
			name := strings.TrimPrefix(ref, "#/components/parameters/")
			if filtered.Components.Parameters == nil {
				filtered.Components.Parameters = map[string]*openapi3.ParameterRef{}
			}
			if p, ok := doc.Components.Parameters[name]; ok {
				filtered.Components.Parameters[name] = p
			} else {
				o.logger.Warn("parameters component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/securitySchemes/"):
			name := strings.TrimPrefix(ref, "#/components/securitySchemes/")
			if filtered.Components.SecuritySchemes == nil {
				filtered.Components.SecuritySchemes = map[string]*openapi3.SecuritySchemeRef{}
			}
			if ss, ok := doc.Components.SecuritySchemes[name]; ok {
				filtered.Components.SecuritySchemes[name] = ss
			} else {
				o.logger.Warn("security schema component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/requestBodies/"):
			name := strings.TrimPrefix(ref, "#/components/requestBodies/")
			if filtered.Components.RequestBodies == nil {
				filtered.Components.RequestBodies = map[string]*openapi3.RequestBodyRef{}
			}
			if rb, ok := doc.Components.RequestBodies[name]; ok {
				filtered.Components.RequestBodies[name] = rb
			} else {
				o.logger.Warn("request body component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/responses/"):
			name := strings.TrimPrefix(ref, "#/components/responses/")
			if filtered.Components.Responses == nil {
				filtered.Components.Responses = map[string]*openapi3.ResponseRef{}
			}
			if r, ok := doc.Components.Responses[name]; ok {
				filtered.Components.Responses[name] = r
			} else {
				o.logger.Warn("response component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/headers/"):
			name := strings.TrimPrefix(ref, "#/components/headers/")
			if filtered.Components.Headers == nil {
				filtered.Components.Headers = map[string]*openapi3.HeaderRef{}
			}
			if h, ok := doc.Components.Headers[name]; ok {
				filtered.Components.Headers[name] = h
			} else {
				o.logger.Warn("headers component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/examples/"):
			name := strings.TrimPrefix(ref, "#/components/examples/")
			if filtered.Components.Examples == nil {
				filtered.Components.Examples = map[string]*openapi3.ExampleRef{}
			}
			if e, ok := doc.Components.Examples[name]; ok {
				filtered.Components.Examples[name] = e
			} else {
				o.logger.Warn("examples component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/links/"):
			name := strings.TrimPrefix(ref, "#/components/links/")
			if filtered.Components.Links == nil {
				filtered.Components.Links = map[string]*openapi3.LinkRef{}
			}
			if l, ok := doc.Components.Links[name]; ok {
				filtered.Components.Links[name] = l
			} else {
				o.logger.Warn("links component not found",
					zap.String("ref", ref))
			}

		case strings.HasPrefix(ref, "#/components/callbacks/"):
			name := strings.TrimPrefix(ref, "#/components/callbacks/")
			if filtered.Components.Callbacks == nil {
				filtered.Components.Callbacks = map[string]*openapi3.CallbackRef{}
			}
			if c, ok := doc.Components.Callbacks[name]; ok {
				filtered.Components.Callbacks[name] = c
			} else {
				o.logger.Warn("callback component not found",
					zap.String("ref", ref))
			}

		default:
			o.logger.Warn("unhandled $ref component", zap.String("ref", ref))
		}
	}

	f, err := os.Create(outSpecPath)
	if err != nil {
		o.logger.Error("failed to create output spec file",
			zap.Error(fmt.Errorf("os.Create: %w", err)))
		return fmt.Errorf("os.Create: %w", err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	defer encoder.Close()

	if err := encoder.Encode(filtered); err != nil {
		o.logger.Error("failed to encode output spec file",
			zap.Error(fmt.Errorf("encoder.Encode: %w", err)))
		return fmt.Errorf("encoder.Encode: %w", err)
	}
	return nil
}
