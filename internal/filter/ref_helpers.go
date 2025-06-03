package filter

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

func refName(def, name string) (ref string) {
	return fmt.Sprintf("#/components/%s/%s", def, name)
}

func parseRef(ref string) (def, name string) {
	elems := strings.Split(ref, "/")
	def, name = elems[2], elems[3]
	return def, name
}

func initializeDef[M ~map[string]V, V any](targetMap *M) {
	if *targetMap == nil {
		*targetMap = make(M)
	}
}

func copyComponent[V any](
	logger *zap.Logger,
	docDef, filteredDef map[string]V,
	name, ref string,
) {
	if component, ok := docDef[name]; ok {
		filteredDef[name] = component
	} else {
		logger.Warn("component not found", zap.String("ref", ref))
	}
}

func isEmptyComponents(c *openapi3.Components) bool {
	if c == nil {
		return true
	}
	return len(c.Schemas) == 0 &&
		len(c.Parameters) == 0 &&
		len(c.Headers) == 0 &&
		len(c.RequestBodies) == 0 &&
		len(c.Responses) == 0 &&
		len(c.SecuritySchemes) == 0 &&
		len(c.Examples) == 0 &&
		len(c.Links) == 0 &&
		len(c.Callbacks) == 0
}
