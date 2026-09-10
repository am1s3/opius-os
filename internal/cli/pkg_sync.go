package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/config"
	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync packages from remote registry",
	Long: `Download and sync packages from the remote registry (GitHub).

This command fetches the latest package manifests and files from
the configured registry URL and updates the local registry.

The default registry is: ` + config.DefaultRegistryURL + `

Examples:
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

		components.PrintSection("Registry Sync")
		fmt.Printf("  %-12s %s\n", "Source", style.Value.Render(cfg.RegistryURL))
		fmt.Printf("  %-12s %s\n", "Local", style.Dim.Render(regDir))
		fmt.Println()

		progress := make(chan string)
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

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " registry synced successfully")
		fmt.Println()
		fmt.Println(style.Dim.Render("  Search packages with:"))
		fmt.Println(style.Dim.Render("    opius pkg search"))
		return nil
	},
}