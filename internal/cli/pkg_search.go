package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgSearchQuery string

var pkgSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search packages in the registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = strings.Join(args, " ")
		}

		regPath, err := pkg.DefaultRegistryPath()
		if err != nil {
			return err
		}
		registry := pkg.NewRegistry(regPath)

		results, err := registry.Search(query)
		if err != nil {
			return err
		}

		components.PrintSection("Package Search")
		if query != "" {
			fmt.Printf("  query: %s\n\n", style.Value.Render(query))
		}

		if len(results) == 0 {
			fmt.Println("  " + components.BadgeWarn() + " no packages found")
			fmt.Println(style.Dim.Render("  Add packages to the registry with:"))
			fmt.Println(style.Dim.Render("    opius pkg init"))
			return nil
		}

		fmt.Printf("  %-20s %-10s %-40s\n",
			style.Header.Render("NAME"),
			style.Header.Render("VERSION"),
			style.Header.Render("DESCRIPTION"),
		)

		for _, m := range results {
			desc := m.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			fmt.Printf("  %-20s %-10s %-40s\n",
				style.Value.Render(m.Name),
				style.Dim.Render(m.Version),
				desc,
			)
		}

		fmt.Println()
		fmt.Printf("  %d package(s) found\n", len(results))
		return nil
	},
}

func init() {
	pkgSearchCmd.Flags().StringVarP(&pkgSearchQuery, "query", "q", "", "Search query")
}