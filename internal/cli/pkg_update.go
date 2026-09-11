package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"sync"},
	Short:   "Update package registry from remote source",
	Long: `Download and sync the latest package manifests from the remote registry.

This fetches all available firmware packages and updates the local registry.
After running this command, use 'opius pkg search' to see available packages.

Default registry: https://github.com/am1s3/opius-os-pkg

Examples:
  opius pkg update
  opius pkg sync`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheDir, err := pkg.DefaultCacheDir()
		if err != nil {
			return err
		}

		reg := pkg.NewRegistry(pkg.DefaultRegistryURL(), cacheDir)

		components.PrintSection("Registry Update")
		fmt.Printf("  %-12s %s\n", "Source", style.Value.Render(pkg.DefaultRegistryURL()))
		fmt.Printf("  %-12s %s\n", "Local", style.Dim.Render(cacheDir))
		fmt.Println()

		fmt.Print("  Syncing... ")

		if err := reg.Sync(); err != nil {
			fmt.Println(components.BadgeErr())
			fmt.Println("  " + components.BadgeErr() + " " + err.Error())
			return err
		}

		fmt.Println(components.BadgeOK())

		// Count packages
		packages, _ := reg.List()

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " registry updated successfully")
		fmt.Printf("  %-12s %d\n", "Packages", len(packages))
		fmt.Println()
		fmt.Println(style.Dim.Render("  Search packages with:"))
		fmt.Println(style.Dim.Render("    opius pkg search"))
		fmt.Println(style.Dim.Render("    opius pkg list"))
		return nil
	},
}