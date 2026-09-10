package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	pkgRemoveVersion string
	pkgRemoveYes     bool
)

var pkgRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an installed package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := pkgRemoveVersion

		storePath, err := pkg.DefaultStorePath()
		if err != nil {
			return err
		}
		store := pkg.NewStore(storePath)

		// If version not specified, find latest installed
		if version == "" {
			installed, err := store.List()
			if err != nil {
				return err
			}
			for _, ip := range installed {
				if ip.Name == name {
					version = ip.Version
					break
				}
			}
			if version == "" {
				return fmt.Errorf("package %s is not installed", name)
			}
		}

		components.PrintSection("Remove Package")
		fmt.Printf("  %-12s %s\n", "Package", style.Value.Render(name))
		fmt.Printf("  %-12s %s\n", "Version", style.Value.Render(version))
		fmt.Printf("  %-12s %s\n", "Path", style.Dim.Render(store.PackageDir(name, version)))

		if !pkgRemoveYes {
			if !components.PromptConfirm("  Remove package?", false) {
				fmt.Println("  " + components.BadgeWarn() + " aborted by user")
				return nil
			}
		}

		if err := store.Remove(name, version); err != nil {
			return err
		}

		fmt.Println("  " + components.BadgeOK() + " package removed")
		return nil
	},
}

func init() {
	pkgRemoveCmd.Flags().StringVarP(&pkgRemoveVersion, "version", "v", "", "Package version to remove (defaults to latest)")
	pkgRemoveCmd.Flags().BoolVarP(&pkgRemoveYes, "yes", "y", false, "Skip confirmation prompt")
}