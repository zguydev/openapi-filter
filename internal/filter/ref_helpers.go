package filter

import (
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

func initializeComponentMap[M ~map[string]V, V any](targetMap *M) {
	if *targetMap == nil {
		*targetMap = make(M)
	}
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

func processCopyComponent(
	docComps, filteredComps *openapi3.Components,
	typ ComponentType,
	name string,
) (ok bool) {
	docCompMap := ComponentTypeToComponentMap(docComps, typ)
	filteredCompMap := ComponentTypeToComponentMap(filteredComps, typ)
	initializeComponentMap(&filteredCompMap)
	return copyComponent(docCompMap, filteredCompMap, name)
}

func isEmptyComponents(c *openapi3.Components) bool {
	if c == nil {
		return true
	}
	for _, compTyp := range ComponentTypes() {
		compMap := ComponentTypeToComponentMap(c, compTyp)
		if len(compMap) != 0 {
			return false
		}
	}
	return true
}
