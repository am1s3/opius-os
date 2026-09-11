package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available firmware packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg := newPkgRegistry()

		packages, err := reg.List()
		if err != nil {
			if err := reg.Sync(); err != nil {
				return fmt.Errorf("failed to load registry: %w", err)
			}
			packages, err = reg.List()
			if err != nil {
				return err
			}
		}

		components.PrintSection("Available Packages")

		if len(packages) == 0 {
			fmt.Println("  " + components.BadgeWarn() + " No packages available")
			fmt.Println(style.Dim.Render("  Run 'opius update' to update Opius and sync the registry"))
			return nil
		}

		fmt.Printf("  %-30s %-12s %s\n",
			style.Header.Render("NAME"),
			style.Header.Render("CHIPS"),
			style.Header.Render("DESCRIPTION"),
		)

		for _, p := range packages {
			chips := "-"
			if len(p.Chips) > 0 {
				chips = joinStrings(p.Chips)
			}
			desc := p.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}

			fmt.Printf("  %-30s %-12s %s\n",
				style.Value.Render(p.Name),
				style.Dim.Render(chips),
				desc,
			)
		}

		fmt.Println()
		fmt.Printf("  %d package(s) available\n", len(packages))
		fmt.Println(style.Dim.Render("  Info:     opius pkg info <name>"))
		fmt.Println(style.Dim.Render("  Download: opius pkg download <name>"))
		fmt.Println(style.Dim.Render("  Flash:    opius pkg flash <name>"))

		return nil
	},
}

func joinStrings(items []string) string {
	result := ""
	for i, s := range items {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}