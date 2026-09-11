package cli

import (
	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Opius package manager",
	Long: `Opius package manager — simple firmware management.

Commands:
  list      Show all available firmware packages
  info      Show detailed info about a package
  download  Download a firmware .bin file
  flash     Download and flash a package to device

Examples:
  opius pkg list
  opius pkg info bruce
  opius pkg download bruce
  opius pkg flash bruce`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	pkgCmd.AddCommand(pkgSearchCmd)
	pkgCmd.AddCommand(pkgListCmd)
	pkgCmd.AddCommand(pkgInfoCmd)
	pkgCmd.AddCommand(pkgDownloadCmd)
	pkgCmd.AddCommand(pkgInstallCmd)
	pkgCmd.AddCommand(pkgFlashCmd)
	pkgCmd.AddCommand(pkgInitCmd)
}

func newPkgRegistry() *pkg.Registry {
	cacheDir, err := pkg.DefaultCacheDir()
	if err != nil {
		cacheDir = ""
	}
	return pkg.NewRegistry(pkg.DefaultRegistryURL(), cacheDir)
}