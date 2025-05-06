package config

type FilterConfig struct {
	Tool       *ToolConfig             `mapstructure:"x-openapi-filter"`
	Paths      map[string][]string     `mapstructure:"paths"`
	Components *FilterComponentsConfig `mapstructure:"components"`
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
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}
