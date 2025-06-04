package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var ErrConfigPathEmpty = errors.New("config path is empty")

func initConfig[C any](configPath string) (*C, error) {
	k := koanf.New(".")

	configExt := strings.TrimLeft(filepath.Ext(configPath), ".")

	var parser koanf.Parser
	switch configExt {
	case "yaml", "yml":
		parser = yaml.Parser()
	case "toml":
		parser = toml.Parser()
	case "json":
		parser = json.Parser()
	default:
		return nil, fmt.Errorf("unsupported config format: %s", configExt)
	}

	if err := k.Load(file.Provider(configPath), parser); err != nil {
		return nil, fmt.Errorf("k.Load: %w", err)
	}

	var cfg C
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("k.Unmarshal: %w", err)
	}
	return &cfg, nil
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, ErrConfigPathEmpty
	}
	cfg, err := initConfig[Config](configPath)
	if err != nil {
		return nil, fmt.Errorf("initConfig[Config]: %w", err)
	}
	return cfg, nil
}
