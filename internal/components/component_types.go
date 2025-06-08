package components

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/zguydev/openapi-filter/pkg/config"
)

type ComponentType int

const (
	_ ComponentType = iota
	ComponentTypeSchema
	ComponentTypeParameter
	ComponentTypeHeader
	ComponentTypeRequestBody
	ComponentTypeResponse
	ContentTypeSecuritySchema
	ContentTypeExample
	ContentTypeLink
	ContentTypeCallback
)

func ComponentTypes() []ComponentType {
	return []ComponentType{
		ComponentTypeSchema,
		ComponentTypeParameter,
		ComponentTypeHeader,
		ComponentTypeRequestBody,
		ComponentTypeResponse,
		ContentTypeSecuritySchema,
		ContentTypeExample,
		ContentTypeLink,
		ContentTypeCallback,
	}
}

var componentDefToTypeMap = map[string]ComponentType{
	"schemas":         ComponentTypeSchema,
	"parameters":      ComponentTypeParameter,
	"headers":         ComponentTypeHeader,
	"requestBodies":   ComponentTypeRequestBody,
	"responses":       ComponentTypeResponse,
	"securitySchemes": ContentTypeSecuritySchema,
	"examples":        ContentTypeExample,
	"links":           ContentTypeLink,
	"callbacks":       ContentTypeCallback,
}

func ComponentDefToType(def string) (typ ComponentType, ok bool) {
	typ, ok = componentDefToTypeMap[def]
	return typ, ok
}

var componentTypeToDefMap = map[ComponentType]string{
	ComponentTypeSchema:       "schemas",
	ComponentTypeParameter:    "parameters",
	ComponentTypeHeader:       "headers",
	ComponentTypeRequestBody:  "requestBodies",
	ComponentTypeResponse:     "responses",
	ContentTypeSecuritySchema: "securitySchemes",
	ContentTypeExample:        "examples",
	ContentTypeLink:           "links",
	ContentTypeCallback:       "callbacks",
}

func ComponentTypeToDef(typ ComponentType) string {
	return componentTypeToDefMap[typ]
}

func ComponentTypeToComponentMap[T any](
	components *openapi3.Components,
	typ ComponentType,
) (compMap map[string]T) {
	switch typ {
	case ComponentTypeSchema:
		return any(map[string]*openapi3.SchemaRef(components.Schemas)).(map[string]T)
	case ComponentTypeParameter:
		return any(map[string]*openapi3.ParameterRef(components.Parameters)).(map[string]T)
	case ComponentTypeHeader:
		return any(map[string]*openapi3.HeaderRef(components.Headers)).(map[string]T)
	case ComponentTypeRequestBody:
		return any(map[string]*openapi3.RequestBodyRef(components.RequestBodies)).(map[string]T)
	case ComponentTypeResponse:
		return any(map[string]*openapi3.ResponseRef(components.Responses)).(map[string]T)
	case ContentTypeSecuritySchema:
		return any(map[string]*openapi3.SecuritySchemeRef(components.SecuritySchemes)).(map[string]T)
	case ContentTypeExample:
		return any(map[string]*openapi3.ExampleRef(components.Examples)).(map[string]T)
	case ContentTypeLink:
		return any(map[string]*openapi3.LinkRef(components.Links)).(map[string]T)
	case ContentTypeCallback:
		return any(map[string]*openapi3.CallbackRef(components.Callbacks)).(map[string]T)
	default:
		panic(fmt.Errorf("unsupported component type: %T", typ))
	}
}

func ComponentTypeToCfgNames(
	cfg *config.FilterComponentsConfig,
	typ ComponentType,
) []string {
	switch typ {
	case ComponentTypeSchema:
		return cfg.Schemas
	case ComponentTypeParameter:
		return cfg.Parameters
	case ComponentTypeHeader:
		return cfg.Headers
	case ComponentTypeRequestBody:
		return cfg.RequestBodies
	case ComponentTypeResponse:
		return cfg.Responses
	case ContentTypeSecuritySchema:
		return cfg.SecuritySchemes
	case ContentTypeExample:
		return cfg.Examples
	case ContentTypeLink:
		return cfg.Links
	case ContentTypeCallback:
		return cfg.Callbacks
	default:
		panic(fmt.Errorf("unsupported component type: %T", typ))
	}
}
