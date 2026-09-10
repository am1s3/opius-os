package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/config"
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
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		regDir, err := cfg.RegistryDir()
		if err != nil {
			return err
		}

		remote := pkg.NewRemoteRegistry(cfg.RegistryURL, regDir)

		components.PrintSection("Registry Update")
		fmt.Printf("  %-12s %s\n", "Source", style.Value.Render(cfg.RegistryURL))
		fmt.Printf("  %-12s %s\n", "Local", style.Dim.Render(regDir))
		fmt.Println()

		progress := make(chan string, 10)
		errCh := make(chan error, 1)

		go func() {
			errCh <- remote.Sync(progress)
		}()

		for msg := range progress {
			fmt.Println("  " + style.Info.Render("→") + " " + msg)
		}

		if err := <-errCh; err != nil {
			fmt.Println()
			fmt.Println("  " + components.BadgeErr() + " " + err.Error())
			return err
		}

		// Count packages
		registry := pkg.NewRegistry(regDir)
		packages, _ := registry.List()

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " registry updated successfully")
		fmt.Printf("  %-12s %d\n", "Packages", len(packages))
		fmt.Println()
		fmt.Println(style.Dim.Render("  Search packages with:"))
		fmt.Println(style.Dim.Render("    opius pkg search"))
		return nil
	},
}