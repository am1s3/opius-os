package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show detailed information about a package",
	Args:  cobra.ExactArgs(1),
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

		components.PrintSection("Package Info")
		fmt.Printf("  %-14s %s\n", "Name", style.Value.Render(m.Name))
		fmt.Printf("  %-14s %s\n", "Version", style.Value.Render(m.Version))
		fmt.Printf("  %-14s %s\n", "Description", m.Description)
		fmt.Printf("  %-14s %s\n", "License", emptyDashPkg(m.License))
		fmt.Printf("  %-14s %s\n", "Homepage", emptyDashPkg(m.Homepage))
		fmt.Printf("  %-14s %s\n", "Repository", emptyDashPkg(m.Repository))

		if len(m.Authors) > 0 {
			fmt.Printf("  %-14s %s\n", "Authors", style.Dim.Render(joinStrings(m.Authors)))
		}

		if len(m.Chips) > 0 {
			fmt.Printf("  %-14s %s\n", "Chips", style.Value.Render(joinStrings(m.Chips)))
		}

		if len(m.Boards) > 0 {
			fmt.Printf("  %-14s %s\n", "Boards", style.Value.Render(joinStrings(m.Boards)))
		}

		if len(m.Dependencies) > 0 {
			fmt.Printf("  %-14s\n", "Dependencies")
			for _, d := range m.Dependencies {
				fmt.Printf("    - %s %s\n", style.Value.Render(d.Name), style.Dim.Render(d.Version))
			}
		}

		if len(m.Files) > 0 {
			fmt.Printf("\n  %-14s\n", "Files")
			for _, f := range m.Files {
				size := formatSize(f.Size)
				fmt.Printf("    %-30s %-10s %-8s\n",
					style.Dim.Render(f.Path),
					size,
					style.Dim.Render(f.Type),
				)
			}
		}

		return nil
	},
}

func emptyDashPkg(v string) string {
	if v == "" {
		return style.Dim.Render("-")
	}
	return style.Value.Render(v)
}

func joinStrings(items []string) string {
	out := ""
	for i, s := range items {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}

func formatSize(bytes int64) string {
	switch {
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}