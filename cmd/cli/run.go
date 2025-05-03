package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/zguydev/openapi-filter/internal/utils"
)

func run(cmd *cobra.Command, args []string) {
	fallbackLogger := utils.NewFallbackLogger()
	defer fallbackLogger.Sync() //nolint:errcheck

	_, err := cmd.Flags().GetString("config")
	if err != nil {
		fallbackLogger.Fatal("failed to get config flag", zap.Error(err))
	}

	
}
