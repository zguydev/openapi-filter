package utils

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/zguydev/openapi-filter/internal/config"
)

func NewFallbackLogger() *zap.Logger {
	var zapCfg zap.Config
	switch strings.ToLower(os.Getenv("APP_ENV")) {
	case "production", "prod":
		zapCfg = zap.NewProductionConfig()
	default:
		zapCfg = zap.NewDevelopmentConfig()
	}
	zapCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)

	logger, err := zapCfg.Build()
	if err != nil {
		panic(fmt.Errorf("failed to create fallback logger: %w", err))
	}
	return logger
}

func NewLogger(cfg *config.LoggerConfig) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("wrong logger level in config: %w", err)
	}

	var zapCfg zap.Config
	switch strings.ToLower(os.Getenv("APP_ENV")) {
	case "production", "prod":
		zapCfg = zap.NewProductionConfig()
	default:
		zapCfg = zap.NewDevelopmentConfig()
	}
	zapCfg.Level = level
	zapCfg.OutputPaths = []string{"stdout"}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("zapCfg.Build: %w", err)
	}
	return logger, nil
}
