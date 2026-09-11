package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Create a new package manifest template",
	Long: `Create a template manifest.json for a new firmware package.

This generates a starter manifest file that you can edit and submit
to the Opius package registry.

Examples:
  opius pkg init my-firmware
  opius pkg init bruce-custom`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "my-firmware"
		if len(args) > 0 {
			name = args[0]
		}

		// Create manifest template
		manifest := pkg.PackageManifest{
			Name:        name,
			Description: "Description of " + name,
			Category:    "custom",
			Chips:       []string{"esp32"},
			Boards:      []string{},
			Homepage:    "https://github.com/your-username/" + name,
			License:     "MIT",
			Versions: []pkg.PackageVersion{
				{
					Version:     "1.0.0",
					ReleaseDate: "2026-09-11",
					DownloadURL: "https://github.com/your-username/" + name + "/releases/download/v1.0.0/firmware.bin",
					SHA256:      "",
					Size:        0,
				},
			},
		}

		// Create output directory
		outputDir := filepath.Join(".", name)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Write manifest.json
		manifestPath := filepath.Join(outputDir, "manifest.json")
		data, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal manifest: %w", err)
		}

		if err := os.WriteFile(manifestPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write manifest: %w", err)
		}

		components.PrintSection("Package Initialized")
		fmt.Printf("  %-12s %s\n", "Name", style.Value.Render(name))
		fmt.Printf("  %-12s %s\n", "Manifest", style.Dim.Render(manifestPath))
		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " Template created!")
		fmt.Println()
		fmt.Println(style.Info.Render("  Next steps:"))
		fmt.Println(style.Dim.Render("    1. Edit " + manifestPath))
		fmt.Println(style.Dim.Render("    2. Upload your firmware.bin to GitHub releases"))
		fmt.Println(style.Dim.Render("    3. Update download_url and sha256 in manifest"))
		fmt.Println(style.Dim.Render("    4. Submit a PR to https://github.com/am1s3/opius-os-pkg"))
		fmt.Println()

		return nil
	},
}