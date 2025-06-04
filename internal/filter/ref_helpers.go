package filter

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func parseRef(ref string) (def, name string, ok bool) {
	elems := strings.Split(ref, "/")
	if len(elems) != 4 {
		return "", "", false
	}
	def, name = elems[2], elems[3]
	return def, name, true
}

func initializeComponentMap[M ~map[string]V, V any](targetMap *M) bool {
	if *targetMap == nil {
		*targetMap = make(M)
		return true
	}
	return false
}

func copyComponent[M ~map[string]V, V any](
	docCompMap, filteredCompMap M,
	name string,
) (ok bool) {
	if component, ok := docCompMap[name]; ok {
		filteredCompMap[name] = component
		return true
	}
	return false
}

func processCopyComponentByType[T any](
	docComps, filteredComps *openapi3.Components,
	typ ComponentType,
	name string,
) (ok bool) {
	docCompMap := ComponentTypeToComponentMap[T](docComps, typ)
	filteredCompMap := ComponentTypeToComponentMap[T](filteredComps, typ)
	if initializeComponentMap(&filteredCompMap) {
		setComponentMapInComponents(filteredComps, typ, filteredCompMap)
	}
	return copyComponent(docCompMap, filteredCompMap, name)
}

func processCopyComponent(
	docComps, filteredComps *openapi3.Components,
	typ ComponentType,
	name string,
) (ok bool) {
	switch typ {
	case ComponentTypeSchema:
		return processCopyComponentByType[*openapi3.SchemaRef](docComps, filteredComps, typ, name)
	case ComponentTypeParameter:
		return processCopyComponentByType[*openapi3.ParameterRef](docComps, filteredComps, typ, name)
	case ComponentTypeHeader:
		return processCopyComponentByType[*openapi3.HeaderRef](docComps, filteredComps, typ, name)
	case ComponentTypeRequestBody:
		return processCopyComponentByType[*openapi3.RequestBodyRef](docComps, filteredComps, typ, name)
	case ComponentTypeResponse:
		return processCopyComponentByType[*openapi3.ResponseRef](docComps, filteredComps, typ, name)
	case ContentTypeSecuritySchema:
		return processCopyComponentByType[*openapi3.SecuritySchemeRef](docComps, filteredComps, typ, name)
	case ContentTypeExample:
		return processCopyComponentByType[*openapi3.ExampleRef](docComps, filteredComps, typ, name)
	case ContentTypeLink:
		return processCopyComponentByType[*openapi3.LinkRef](docComps, filteredComps, typ, name)
	case ContentTypeCallback:
		return processCopyComponentByType[*openapi3.CallbackRef](docComps, filteredComps, typ, name)
	default:
		panic(fmt.Errorf("unsupported component type: %v", typ))
	}
}

func setComponentMapInComponents[T any](
	components *openapi3.Components,
	typ ComponentType,
	compMap map[string]T,
) {
	switch typ {
	case ComponentTypeSchema:
		components.Schemas = any(compMap).(map[string]*openapi3.SchemaRef)
	case ComponentTypeParameter:
		components.Parameters = any(compMap).(map[string]*openapi3.ParameterRef)
	case ComponentTypeHeader:
		components.Headers = any(compMap).(map[string]*openapi3.HeaderRef)
	case ComponentTypeRequestBody:
		components.RequestBodies = any(compMap).(map[string]*openapi3.RequestBodyRef)
	case ComponentTypeResponse:
		components.Responses = any(compMap).(map[string]*openapi3.ResponseRef)
	case ContentTypeSecuritySchema:
		components.SecuritySchemes = any(compMap).(map[string]*openapi3.SecuritySchemeRef)
	case ContentTypeExample:
		components.Examples = any(compMap).(map[string]*openapi3.ExampleRef)
	case ContentTypeLink:
		components.Links = any(compMap).(map[string]*openapi3.LinkRef)
	case ContentTypeCallback:
		components.Callbacks = any(compMap).(map[string]*openapi3.CallbackRef)
	default:
		panic(fmt.Errorf("unsupported component type: %v", typ))
	}
}

func isComponentMapEmpty(
	components *openapi3.Components,
	typ ComponentType,
) bool {
	switch typ {
	case ComponentTypeSchema:
		return len(components.Schemas) != 0
	case ComponentTypeParameter:
		return len(components.Parameters) != 0
	case ComponentTypeHeader:
		return len(components.Headers) != 0
	case ComponentTypeRequestBody:
		return len(components.RequestBodies) != 0
	case ComponentTypeResponse:
		return len(components.Responses) != 0
	case ContentTypeSecuritySchema:
		return len(components.SecuritySchemes) != 0
	case ContentTypeExample:
		return len(components.Examples) != 0
	case ContentTypeLink:
		return len(components.Links) != 0
	case ContentTypeCallback:
		return len(components.Callbacks) != 0
	default:
		panic(fmt.Errorf("unsupported component type: %v", typ))
	}
}

func isEmptyComponents(c *openapi3.Components) bool {
	if c == nil {
		return true
	}
	for _, compTyp := range ComponentTypes() {
		if !isComponentMapEmpty(c, compTyp) {
			return false
		}
	}
	return true
}
