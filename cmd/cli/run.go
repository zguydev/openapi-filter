package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/zguydev/openapi-filter/internal/config"
	"github.com/zguydev/openapi-filter/internal/filter"
	"github.com/zguydev/openapi-filter/internal/utils"
)

func run(cmd *cobra.Command, args []string) {
	fallbackLogger := utils.NewFallbackLogger()
	defer fallbackLogger.Sync() //nolint:errcheck

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		fallbackLogger.Fatal("failed to get config flag", zap.Error(err))
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fallbackLogger.Fatal("failed to load config",
			zap.Error(fmt.Errorf("config.LoadConfig: %w", err)))
	}

	logger, err := utils.NewLogger(cfg.Tool.Logger)
	if err != nil {
		fallbackLogger.Fatal("failed to init logger",
			zap.Error(fmt.Errorf("utils.NewLogger: %w", err)))
	}

	oaf := filter.NewOpenAPISpecFilter(cfg, logger)
	if err := oaf.Filter(args[0], args[1]); err != nil {
		logger.Error("filter on spec failed",
			zap.Error(fmt.Errorf("oaf.Filter: %w", err)))
		os.Exit(1)
	}
}
