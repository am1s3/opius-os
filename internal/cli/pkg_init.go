package cli

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a sample package manifest in the registry",
	Long: `Create a sample package manifest in the local registry.

This is useful for learning the package format. After running this
command, edit the generated .toml file and add your firmware files
to the registry directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		regPath, err := pkg.DefaultRegistryPath()
		if err != nil {
			return err
		}
		registry := pkg.NewRegistry(regPath)

		sample := &pkg.Manifest{
			Name:        "example-firmware",
			Version:     "1.0.0",
			Description: "Sample firmware package for learning Opius package format",
			License:     "MIT",
			Homepage:    "https://example.com",
			Repository:  "https://github.com/example/example-firmware",
			Authors:     []string{"Opius Team"},
			Chips:       []string{"esp32", "esp32-s3"},
			Boards:      []string{"esp32-devkit-v1", "esp32-s3-devkit"},
			CreatedAt:   time.Now(),
			Files: []pkg.File{
				{
					Path:   "firmware.bin",
					SHA256: "",
					Size:   0,
					Type:   "firmware",
				},
				{
					Path:   "README.md",
					SHA256: "",
					Size:   0,
					Type:   "doc",
				},
			},
		}

		if err := registry.Add(sample); err != nil {
			return err
		}

		manifestPath := filepath.Join(regPath, sample.Name+"-"+sample.Version+".toml")

		components.PrintSection("Sample Package Created")
		fmt.Printf("  %-12s %s\n", "Name", style.Value.Render(sample.Name))
		fmt.Printf("  %-12s %s\n", "Version", style.Value.Render(sample.Version))
		fmt.Printf("  %-12s %s\n", "Manifest", style.Dim.Render(manifestPath))
		fmt.Println()
		fmt.Println(style.Info.Render("  Next steps:"))
		fmt.Println(style.Dim.Render("    1. Edit the manifest to match your firmware"))
		fmt.Println(style.Dim.Render("    2. Put your files in the registry directory:"))
		fmt.Println(style.Dim.Render("       " + regPath))
		fmt.Println(style.Dim.Render("    3. Update SHA-256 hashes with:"))
		fmt.Println(style.Dim.Render("       shasum -a 256 <file>"))
		fmt.Println(style.Dim.Render("    4. Search for it with:"))
		fmt.Println(style.Dim.Render("       opius pkg search"))
		return nil
	},
}