// Package config provides configuration structures for the OpenAPI filter tool.
// It defines the configuration format for filtering OpenAPI specs and
// tool-specific settings.
package config

// Config represents the root configuration structure for the OpenAPI filter tool.
// It combines tool-specific settings with filter configuration.
type Config struct {
	Tool         ToolConfig `koanf:"x-openapi-filter"`
	FilterConfig `koanf:",squash"`
}

// FilterConfig defines the configuration for filtering an OpenAPI spec.
// It specifies which parts of the spec should be included in the output.
type FilterConfig struct {
	Servers      bool                    `koanf:"servers"`      // Include servers section
	Paths        map[string][]string     `koanf:"paths"`        // Map of paths to allowed HTTP methods
	Components   *FilterComponentsConfig `koanf:"components"`   // Component filtering configuration
	Security     bool                    `koanf:"security"`     // Include security requirements
	Tags         bool                    `koanf:"tags"`         // Include tags
	ExternalDocs bool                    `koanf:"externalDocs"` // Include external documentation
}

// FilterComponentsConfig specifies which components should be included in the
// filtered OpenAPI spec. Each field is a list of component names to include.
type FilterComponentsConfig struct {
	Schemas         []string `koanf:"schemas"`         // List of schema names to include
	Parameters      []string `koanf:"parameters"`      // List of parameter names to include
	SecuritySchemes []string `koanf:"securitySchemes"` // List of security scheme names to include
	RequestBodies   []string `koanf:"requestBodies"`   // List of request body names to include
	Responses       []string `koanf:"responses"`       // List of response names to include
	Headers         []string `koanf:"headers"`         // List of header names to include
	Examples        []string `koanf:"examples"`        // List of example names to include
	Links           []string `koanf:"links"`           // List of link names to include
	Callbacks       []string `koanf:"callbacks"`       // List of callback names to include
}

// ToolConfig contains tool-specific configuration settings.
type ToolConfig struct {
	Logger *LoggerConfig `koanf:"logger"` // Logger configuration
	Loader *LoaderConfig `koanf:"loader"` // OpenAPI loader configuration
}

// LoggerConfig defines the logging configuration for the tool.
type LoggerConfig struct {
	Level string `koanf:"level"` // Log level (e.g., "debug", "info", "warn", "error")
}

// LoaderConfig defines configuration for the OpenAPI spec loader.
type LoaderConfig struct {
	IsExternalRefsAllowed bool `koanf:"external_refs_allowed"` // Whether to allow external references
}
