package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/zguydev/openapi-filter/internal"
	"github.com/zguydev/openapi-filter/internal/utils"
	"github.com/zguydev/openapi-filter/pkg/config"
	"github.com/zguydev/openapi-filter/pkg/filter"
	"github.com/zguydev/openapi-filter/pkg/loader"
)

func run(cmd *cobra.Command, args []string) {
	fallbackLogger := utils.NewFallbackLogger()
	defer fallbackLogger.Sync() //nolint:errcheck

	if ok, _ := cmd.Flags().GetBool("version"); ok {
		info, _ := internal.GetInfo()
		fmt.Println(info)
		return
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		fallbackLogger.Fatal("failed to get config flag", zap.Error(err))
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fallbackLogger.Fatal("failed to load config", zap.Error(err))
	}

	logger, err := utils.NewLogger(cfg.Tool.Logger)
	if err != nil {
		fallbackLogger.Fatal("failed to init logger", zap.Error(err))
	}

	inputSpecPath, outSpecPath := args[0], args[1]

	inputSpec, err := internal.LoadSpecFromFile(
		loader.NewLoader(cfg.Tool.Loader), inputSpecPath)
	if err != nil {
		logger.Error("failed to load spec from file",
			zap.Error(err), zap.String("path", inputSpecPath))
		os.Exit(1)
	}
	oaf := filter.NewOpenAPISpecFilter(cfg, logger)
	outSpec, err := oaf.Filter(inputSpec)
	if err != nil {
		logger.Error("filter on spec failed", zap.Error(err))
		os.Exit(1)
	}

	if err := internal.WriteSpecToFile(outSpec, outSpecPath); err != nil {
		logger.Error("failed to write filtered spec file",
			zap.Error(err), zap.String("path", outSpecPath))
		os.Exit(1)
	}
	logger.Info("filtered and saved spec", zap.String("path", outSpecPath))
}
