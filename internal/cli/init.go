package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/config"
	"github.com/opius-os/opius/internal/hostpkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Opius workspace and config",
	RunE: func(cmd *cobra.Command, args []string) error {
		components.PrintBanner()
		
		fmt.Println(style.Info.Render("Detecting system..."))
		fmt.Printf("  %-12s %s\n", "OS", style.Value.Render("macOS"))
		
		// Detect backends
		backends := hostpkg.DetectAll()
		if len(backends) == 0 {
			fmt.Println(components.BadgeWarn(), "No package manager found")
			fmt.Println("  Opius can still work using its own managed tools.")
		} else if len(backends) == 1 {
			fmt.Printf("  %-12s %s\n", "Backend", style.Value.Render(backends[0].Name))
			fmt.Println(components.BadgeOK(), "Using", backends[0].Name, "as default")
		} else {
			// Multiple backends - ask user
			options := make([]string, len(backends))
			for i, b := range backends {
				options[i] = b.Name
			}
			
			choice := components.PromptChoice(
				"Multiple package managers detected. Choose default backend:",
				options,
			)
			
			// Find the selected backend
			var selectedBackend string
			for _, b := range backends {
				if b.Name == choice {
					selectedBackend = b.Name
					break
				}
			}
			
			fmt.Println(components.BadgeOK(), "Selected backend:", selectedBackend)
			
			// Create config
			cfg, err := config.Init()
			if err != nil {
				return err
			}
			cfg.PackageBackend = selectedBackend
			
			// Ask for workspace
			defaultWorkspace := cfg.Workspace
			workspace := components.PromptString("Workspace directory", defaultWorkspace)
			if workspace != "" {
				cfg.Workspace = workspace
			}
			
			// Save config
			if err := cfg.Save(); err != nil {
				return err
			}
			
			fmt.Println()
			fmt.Println(components.BadgeOK(), "Config saved to", cfg.ConfigFile())
			fmt.Println(components.BadgeOK(), "Workspace:", cfg.Workspace)
			
			return nil
		}
		
		// Single backend or no backend
		cfg, err := config.Init()
		if err != nil {
			return err
		}
		
		if len(backends) == 1 {
			cfg.PackageBackend = backends[0].Name
		}
		
		if err := cfg.Save(); err != nil {
			return err
		}
		
		fmt.Println()
		fmt.Println(components.BadgeOK(), "Config written to", cfg.ConfigFile())
		fmt.Println(components.BadgeOK(), "Workspace directories created under ~/.opius")
		
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("backend", "b", "", "Package backend (macports/homebrew)")
}