package config

type Config struct {
	Tool         ToolConfig `mapstructure:"x-openapi-filter"`
	FilterConfig `mapstructure:",squash"`
}

type FilterConfig struct {
	Servers      bool                    `mapstructure:"servers"`
	Paths        map[string][]string     `mapstructure:"paths"`
	Components   *FilterComponentsConfig `mapstructure:"components"`
	Security     []string                `mapstructure:"security"`
	Tags         []string                `mapstructure:"tags"`
	ExternalDocs bool                    `mapstructure:"externalDocs"`
}

type FilterComponentsConfig struct {
	Schemas         []string `mapstructure:"schemas"`
	Parameters      []string `mapstructure:"parameters"`
	SecuritySchemes []string `mapstructure:"securitySchemes"`
	RequestBodies   []string `mapstructure:"requestBodies"`
	Responses       []string `mapstructure:"responses"`
	Headers         []string `mapstructure:"headers"`
	Examples        []string `mapstructure:"examples"`
	Links           []string `mapstructure:"links"`
	Callbacks       []string `mapstructure:"callbacks"`
}

type ToolConfig struct {
	Logger *LoggerConfig `mapstructure:"logger"`
	Loader *LoaderConfig `mapstructure:"loader"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type LoaderConfig struct {
	IsExternalRefsAllowed bool `mapstructure:"external_refs_allowed"`
}
