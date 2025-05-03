package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openapi-filter input output --config filter_config",
	Short: "Filter an OpenAPI spec to only include specified paths/methods",
	Args:  cobra.ExactArgs(2),
	Run:   run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("config", "openapi-filter.yaml", "Path to filter config")
	_ = rootCmd.MarkFlagRequired("config")
}
