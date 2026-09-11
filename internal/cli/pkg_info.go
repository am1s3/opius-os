package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show detailed info about a package",
	Args:  cobra.ExactArgs(1),
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

		components.PrintSection("Package Info")
		fmt.Printf("  %-14s %s\n", "Name", style.Value.Render(m.Name))
		fmt.Printf("  %-14s %s\n", "Description", m.Description)
		fmt.Printf("  %-14s %s\n", "Category", style.Dim.Render(m.Category))
		fmt.Printf("  %-14s %s\n", "License", dashIfEmpty(m.License))
		fmt.Printf("  %-14s %s\n", "Homepage", dashIfEmpty(m.Homepage))

		if len(m.Chips) > 0 {
			fmt.Printf("  %-14s %s\n", "Chips", style.Value.Render(joinStrings(m.Chips)))
		}
		if len(m.Boards) > 0 {
			fmt.Printf("  %-14s %s\n", "Boards", style.Value.Render(joinStrings(m.Boards)))
		}

		if len(m.Versions) > 0 {
			fmt.Println()
			fmt.Println(style.Header.Render("  Available Versions:"))
			for _, v := range m.Versions {
				sizeStr := formatBytes(v.Size)
				fmt.Printf("    %-10s  %-12s  %s\n",
					style.Value.Render(v.Version),
					style.Dim.Render(v.ReleaseDate),
					style.Dim.Render(sizeStr),
				)
			}
		}

		fmt.Println()
		fmt.Println(style.Dim.Render("  Download: opius pkg download " + m.Name))
		fmt.Println(style.Dim.Render("  Flash:    opius pkg flash " + m.Name))

		return nil
	},
}

func dashIfEmpty(v string) string {
	if v == "" {
		return style.Dim.Render("-")
	}
	return v
}

func truncateHash(h string) string {
	if len(h) > 16 {
		return h[:16] + "..."
	}
	return h
}

func formatBytes(bytes int64) string {
	switch {
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}