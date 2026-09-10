package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Enter the Opius interactive shell",
	Run: func(cmd *cobra.Command, args []string) {
		components.PrintBanner()
		fmt.Println("  Opius interactive shell arrives in the next stage.")
		fmt.Println("  This is the foundation build — CLI, config and diagnostics first.")
	},
}