package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	pkgInstallSource string
	pkgInstallYes    bool
)

var pkgInstallCmd = &cobra.Command{
	Use:   "install <name>",
	Short: "Install a package from the registry",
	Long: `Install a package from the local registry into ~/.opius/packages.

The package source directory must contain all files declared in the manifest.
Each file is verified against its SHA-256 hash after installation.

Examples:
  opius pkg install bruce
  opius pkg install bruce --source ./my-packages/bruce`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		regPath, err := pkg.DefaultRegistryPath()
		if err != nil {
			return err
		}
		registry := pkg.NewRegistry(regPath)

		m, err := registry.Find(name)
		if err != nil {
			return err
		}

		storePath, err := pkg.DefaultStorePath()
		if err != nil {
			return err
		}
		store := pkg.NewStore(storePath)

		// Determine source directory
		srcDir := pkgInstallSource
		if srcDir == "" {
			srcDir = regPath
		}

		components.PrintSection("Install Plan")
		fmt.Printf("  %-12s %s\n", "Package", style.Value.Render(m.Name))
		fmt.Printf("  %-12s %s\n", "Version", style.Value.Render(m.Version))
		fmt.Printf("  %-12s %s\n", "Source", style.Dim.Render(srcDir))
		fmt.Printf("  %-12s %s\n", "Target", style.Dim.Render(store.PackageDir(m.Name, m.Version)))
		fmt.Printf("  %-12s %d\n", "Files", len(m.Files))

		if !pkgInstallYes {
			if !components.PromptConfirm("  Install package?", true) {
				fmt.Println("  " + components.BadgeWarn() + " aborted by user")
				return nil
			}
		}

		fmt.Println()
		components.PrintSection("Installing")

		if err := store.Install(m, srcDir); err != nil {
			fmt.Println("  " + components.BadgeErr() + " " + err.Error())
			return err
		}

		fmt.Println("  " + components.BadgeOK() + " installed to " + store.PackageDir(m.Name, m.Version))
		fmt.Println()
		fmt.Println(style.Dim.Render("  List installed packages with:"))
		fmt.Println(style.Dim.Render("    opius pkg list"))
		return nil
	},
}

func init() {
	pkgInstallCmd.Flags().StringVarP(&pkgInstallSource, "source", "s", "", "Source directory for package files (defaults to registry path)")
	pkgInstallCmd.Flags().BoolVarP(&pkgInstallYes, "yes", "y", false, "Skip confirmation prompt")

	// Helper to get absolute path if needed
	_ = filepath.Abs
}