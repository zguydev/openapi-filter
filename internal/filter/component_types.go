package filter

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/zguydev/openapi-filter/internal/config"
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

func ComponentTypeToComponentMap(
	components *openapi3.Components,
	typ ComponentType,
) (compMap map[string]any) {
	var raw any

	switch typ {
	case ComponentTypeSchema:
		raw = components.Schemas
	case ComponentTypeParameter:
		raw = components.Parameters
	case ComponentTypeHeader:
		raw = components.Headers
	case ComponentTypeRequestBody:
		raw = components.RequestBodies
	case ComponentTypeResponse:
		raw = components.Responses
	case ContentTypeSecuritySchema:
		raw = components.SecuritySchemes
	case ContentTypeExample:
		raw = components.Examples
	case ContentTypeLink:
		raw = components.Links
	case ContentTypeCallback:
		raw = components.Callbacks
	default:
		panic(fmt.Errorf("not implemented for %T", typ))
	}

	if raw == nil {
		return nil
	}
	return compMap
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
		panic(fmt.Errorf("not implemented for %T", typ))
	}
}
