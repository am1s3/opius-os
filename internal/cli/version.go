package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/build"
	"github.com/opius-os/opius/internal/ui/style"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Opius OS version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(style.Brand.Render("OPIUS OS"))
		fmt.Printf("  version : %s\n", style.Value.Render(build.Version))
		fmt.Printf("  channel : %s\n", style.Value.Render(build.Channel))
		fmt.Printf("  built   : %s\n", style.Value.Render(build.BuildDate))
	},
}