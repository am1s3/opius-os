package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgInstallVersion string

var pkgInstallCmd = &cobra.Command{
	Use:   "install <name>",
	Short: "Download and install a package",
	Long: `Download a firmware package from the registry.

This is an alias for 'opius pkg download'.

Examples:
  opius pkg install bruce-cyd-2432s028
  opius pkg install bruce-m5stack-cardputer --version 1.16`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		reg := newPkgRegistry()
		if err := reg.Sync(); err != nil {
			return fmt.Errorf("failed to sync registry: %w", err)
		}

		m, err := reg.Find(name)
		if err != nil {
			return err
		}

		if len(m.Versions) == 0 {
			return fmt.Errorf("no versions available for %s", name)
		}

		var version *pkg.PackageVersion
		if pkgInstallVersion != "" {
			for i := range m.Versions {
				if m.Versions[i].Version == pkgInstallVersion {
					version = &m.Versions[i]
					break
				}
			}
			if version == nil {
				return fmt.Errorf("version %s not found", pkgInstallVersion)
			}
		} else {
			version = &m.Versions[0]
		}

		components.PrintSection("Installing Package")
		fmt.Printf("  %-12s %s\n", "Package", style.Value.Render(m.Name))
		fmt.Printf("  %-12s %s\n", "Version", style.Value.Render(version.Version))
		fmt.Println()

		home, _ := os.UserHomeDir()
		downloadDir := filepath.Join(home, ".opius", "downloads")
		os.MkdirAll(downloadDir, 0755)

		filename := filepath.Base(version.DownloadURL)
		if filename == "" || filename == "/" {
			filename = m.Name + "-" + version.Version + ".bin"
		}
		destPath := filepath.Join(downloadDir, filename)

		if err := downloadWithProgress(version.DownloadURL, destPath, version.Size); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		fmt.Println()

		if version.SHA256 != "" {
			fmt.Print("  Verifying SHA256... ")
			if err := verifyFileSHA256(destPath, version.SHA256); err != nil {
				fmt.Println(components.BadgeErr())
				return fmt.Errorf("checksum mismatch: %w", err)
			}
			fmt.Println(components.BadgeOK())
		}

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " Installed to " + style.Value.Render(destPath))

		return nil
	},
}

func init() {
	pkgInstallCmd.Flags().StringVarP(&pkgInstallVersion, "version", "v", "", "Specific version")
}