package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var uninstallYes bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Completely remove Opius OS from this system",
	Long: `Remove all Opius OS files, configuration, packages, and cache.

This will delete:
  - ~/.config/opius/       (configuration)
  - ~/.opius/              (workspace, packages, cache, downloads)
  - The opius binary from /usr/local/bin or ~/.local/bin

WARNING: This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		components.PrintSection("Uninstall Opius OS")
		fmt.Println(style.Warning.Render("  This will permanently remove all Opius OS data."))
		fmt.Println()

		home, _ := os.UserHomeDir()
		paths := []string{
			filepath.Join(home, ".config", "opius"),
			filepath.Join(home, ".opius"),
		}

		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				fmt.Printf("  Will remove: %s\n", style.Dim.Render(p))
			}
		}

		fmt.Println()

		if !uninstallYes {
			if !components.PromptConfirm("  Are you sure you want to uninstall?", false) {
				fmt.Println("  " + components.BadgeWarn() + " Uninstall cancelled")
				return nil
			}
		}

		// Remove directories
		for _, p := range paths {
			if err := os.RemoveAll(p); err != nil {
				fmt.Printf("  %s Failed to remove %s: %v\n", components.BadgeErr(), p, err)
			} else {
				fmt.Printf("  %s Removed %s\n", components.BadgeOK(), style.Dim.Render(p))
			}
		}

		// Remove binary
		binaryPaths := []string{
			"/usr/local/bin/opius",
			filepath.Join(home, ".local", "bin", "opius"),
			filepath.Join(home, "bin", "opius"),
		}

		if runtime.GOOS == "windows" {
			binaryPaths = append(binaryPaths, `C:\Windows\System32\opius.exe`)
		}

		for _, bp := range binaryPaths {
			if _, err := os.Stat(bp); err == nil {
				if err := os.Remove(bp); err != nil {
					fmt.Printf("  %s Failed to remove binary %s: %v\n", components.BadgeErr(), bp, err)
				} else {
					fmt.Printf("  %s Removed binary %s\n", components.BadgeOK(), style.Dim.Render(bp))
				}
			}
		}

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " Opius OS has been uninstalled")
		fmt.Println(style.Dim.Render("  Thank you for using Opius OS!"))

		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolVarP(&uninstallYes, "yes", "y", false, "Skip confirmation")
}