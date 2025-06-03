package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var ErrConfigPathEmpty = errors.New("config path is empty")

func initConfig[C any](configPath string) (*C, error) {
	v := viper.New()
	configExt := strings.TrimLeft(filepath.Ext(configPath), ".")
	v.SetConfigFile(configPath)
	v.SetConfigType(configExt)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("v.ReadInConfig: %w", err)
	}

	for _, k := range v.AllKeys() {
		v.Set(k, v.Get(k))
	}

	cfg := new(C)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("v.Unmarshal: %w", err)
	}
	return cfg, nil
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
