package config

type Config struct {
	Tool         ToolConfig `koanf:"x-openapi-filter"`
	FilterConfig `koanf:",squash"`
}

type FilterConfig struct {
	Servers      bool                    `koanf:"servers"`
	Paths        map[string][]string     `koanf:"paths"`
	Components   *FilterComponentsConfig `koanf:"components"`
	Security     bool                    `koanf:"security"`
	Tags         bool                    `koanf:"tags"`
	ExternalDocs bool                    `koanf:"externalDocs"`
}

type FilterComponentsConfig struct {
	Schemas         []string `koanf:"schemas"`
	Parameters      []string `koanf:"parameters"`
	SecuritySchemes []string `koanf:"securitySchemes"`
	RequestBodies   []string `koanf:"requestBodies"`
	Responses       []string `koanf:"responses"`
	Headers         []string `koanf:"headers"`
	Examples        []string `koanf:"examples"`
	Links           []string `koanf:"links"`
	Callbacks       []string `koanf:"callbacks"`
}

type ToolConfig struct {
	Logger *LoggerConfig `koanf:"logger"`
	Loader *LoaderConfig `koanf:"loader"`
}

type LoggerConfig struct {
	Level string `koanf:"level"`
}

type LoaderConfig struct {
	IsExternalRefsAllowed bool `koanf:"external_refs_allowed"`
}
