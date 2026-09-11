package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
)

var (
	noColor bool
)

var rootCmd = &cobra.Command{
	Use:   "opius",
	Short: "Opius OS — terminal embedded workspace",
	Long: `Opius OS is a terminal-first workspace for microcontroller flashing,
firmware package management and embedded development workflows.

Version 1.2.0-alpha

Repository: https://github.com/am1s3/opius-os`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if noColor {
			os.Setenv("NO_COLOR", "1")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		components.PrintBanner()
		_ = cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(deviceCmd)
	rootCmd.AddCommand(pkgCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(updateCmd)
}