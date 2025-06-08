package loader

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/zguydev/openapi-filter/pkg/config"
)

func NewLoader(cfg *config.LoaderConfig) *openapi3.Loader {
	loader := openapi3.NewLoader()
	if cfg == nil {
		return loader
	}

	loader.IsExternalRefsAllowed = cfg.IsExternalRefsAllowed
	return loader
}
