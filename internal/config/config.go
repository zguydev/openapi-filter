package config

type FilterConfig struct {
	Tool  *ToolConfig         `mapstructure:"x-openapi-filter"`
	Paths map[string][]string `mapstructure:"paths"`
}

type ToolConfig struct {
	Logger *LoggerConfig `mapstructure:"logger"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}
