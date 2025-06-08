package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openapi-filter input_spec output_spec [--config filter_config]",
	Short: "Filter an OpenAPI spec to only include specified paths/methods or components",
	Args:  checkArgs,
	Run:   run,
}

func checkArgs(cmd *cobra.Command, args []string) error {
	if ok, _ := cmd.Flags().GetBool("version"); ok {
		return nil
	}
	const exactArgs = 2
	return cobra.ExactArgs(exactArgs)(cmd, args)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("config", ".openapi-filter.yaml", "Path to filter config")
	rootCmd.Flags().Bool("version", false, "Print version and exit")
}
