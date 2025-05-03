package utils

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
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
