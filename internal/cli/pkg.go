package cli

import (
	"github.com/spf13/cobra"
)

var pkgCmd = &cobra.Command{
	Use:     "pkg",
	Aliases: []string{"opiuspkg"},
	Short:   "Opius package manager",
	Long: `Opius package manager.

Manage firmware packages, tools and configurations through a local
or remote registry. Packages are stored in ~/.opius/packages and
verified with SHA-256 checksums.

Default registry: https://github.com/am1s3/opius-os-pkg`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	pkgCmd.AddCommand(pkgSearchCmd)
	pkgCmd.AddCommand(pkgInfoCmd)
	pkgCmd.AddCommand(pkgInstallCmd)
	pkgCmd.AddCommand(pkgListCmd)
	pkgCmd.AddCommand(pkgRemoveCmd)
	pkgCmd.AddCommand(pkgInitCmd)
	pkgCmd.AddCommand(pkgFlashCmd)
	pkgCmd.AddCommand(pkgSyncCmd)
	pkgCmd.AddCommand(pkgUpdateCmd)
}