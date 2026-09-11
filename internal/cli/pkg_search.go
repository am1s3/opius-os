package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for firmware packages",
	Long: `Search the package registry for firmware packages.

Examples:
  opius pkg search bruce
  opius pkg search m5stack
  opius pkg search cyd`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		reg := newPkgRegistry()
		if err := reg.Sync(); err != nil {
			return fmt.Errorf("failed to sync registry: %w", err)
		}

		results, err := reg.Search(query)
		if err != nil {
			return err
		}

		if query != "" {
			components.PrintSection("Search Results for " + style.Value.Render(query))
		} else {
			components.PrintSection("All Packages")
		}

		if len(results) == 0 {
			fmt.Println("  " + components.BadgeWarn() + " No packages found")
			if query != "" {
				fmt.Println(style.Dim.Render("  Try a different search term"))
			}
			return nil
		}

		fmt.Printf("  %-30s %-12s %s\n",
			style.Header.Render("NAME"),
			style.Header.Render("CHIPS"),
			style.Header.Render("DESCRIPTION"),
		)

		for _, p := range results {
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
		fmt.Printf("  %d package(s) found\n", len(results))
		fmt.Println(style.Dim.Render("  Get details: opius pkg info <name>"))

		return nil
	},
}