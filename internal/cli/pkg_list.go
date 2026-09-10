package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		storePath, err := pkg.DefaultStorePath()
		if err != nil {
			return err
		}
		store := pkg.NewStore(storePath)

		installed, err := store.List()
		if err != nil {
			return err
		}

		components.PrintSection("Installed Packages")

		if len(installed) == 0 {
			fmt.Println("  " + components.BadgeWarn() + " no packages installed")
			fmt.Println(style.Dim.Render("  Install one with:"))
			fmt.Println(style.Dim.Render("    opius pkg install <name>"))
			return nil
		}

		fmt.Printf("  %-20s %-12s %-20s %-30s\n",
			style.Header.Render("NAME"),
			style.Header.Render("VERSION"),
			style.Header.Render("INSTALLED"),
			style.Header.Render("PATH"),
		)

		for _, ip := range installed {
			installedAt := ip.InstalledAt.Format("2006-01-02 15:04")
			if ip.InstalledAt.IsZero() {
				installedAt = "-"
			}
			fmt.Printf("  %-20s %-12s %-20s %-30s\n",
				style.Value.Render(ip.Name),
				style.Dim.Render(ip.Version),
				style.Dim.Render(installedAt),
				style.Dim.Render(ip.InstallDir),
			)
		}

		fmt.Println()
		fmt.Printf("  %d package(s) installed\n", len(installed))
		return nil
	},
}