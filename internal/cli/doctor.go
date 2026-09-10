package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/config"
	"github.com/opius-os/opius/internal/hostpkg"
	"github.com/opius-os/opius/internal/tools"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system, backend and tools health",
	Run: func(cmd *cobra.Command, args []string) {
		components.PrintBanner()
		
		// System info
		components.PrintSection("System")
		fmt.Printf("  %-12s %s %s\n", "OS", runtime.GOOS, components.BadgeOK())
		fmt.Printf("  %-12s %s %s\n", "ARCH", runtime.GOARCH, components.BadgeOK())
		fmt.Printf("  %-12s %s %s\n", "Go", runtime.Version(), components.BadgeOK())
		
		// Config
		components.PrintSection("Config")
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("  %-12s %s\n", "status", components.BadgeErr()+" "+err.Error())
		} else {
			fmt.Printf("  %-12s %s %s\n", "status", "loaded", components.BadgeOK())
			fmt.Printf("  %-12s %s\n", "path", style.Dim.Render(cfg.ConfigFile()))
			fmt.Printf("  %-12s %s\n", "backend", style.Value.Render(cfg.PackageBackend))
			fmt.Printf("  %-12s %s\n", "workspace", style.Dim.Render(cfg.Workspace))
		}
		
		// Backend
		components.PrintSection("Backend")
		b := hostpkg.Detect()
		if b.Path != "" {
			fmt.Printf("  %-12s %s %s\n", "detected", b.Name, components.BadgeOK())
			fmt.Printf("  %-12s %s\n", "path", style.Dim.Render(b.Path))
		} else {
			fmt.Printf("  %-12s %s\n", "detected", components.BadgeWarn()+" none found")
		}
		
		// Tools
		components.PrintSection("Tools")
		allTools := tools.DetectAll()
		
		var installed, missing []tools.Tool
		for _, t := range allTools {
			if t.Installed {
				installed = append(installed, t)
				version := t.Version
				if len(version) > 50 {
					version = version[:47] + "..."
				}
				fmt.Printf("  %-12s %s %s\n", t.Name, style.Dim.Render(version), components.BadgeOK())
			} else {
				missing = append(missing, t)
				fmt.Printf("  %-12s %s\n", t.Name, components.BadgeWarn()+" not found")
			}
		}
		
		// Summary
		fmt.Println()
		if len(missing) > 0 {
			fmt.Println(style.Warning.Render(fmt.Sprintf("  %d tools missing", len(missing))))
			fmt.Println()
			fmt.Println(style.Info.Render("  Install missing tools:"))
			for _, t := range missing {
				cmd := tools.GetInstallCommand(t.Name, b.Name)
				if cmd != "" {
					fmt.Printf("    %s\n", style.Dim.Render(cmd))
				}
			}
		} else {
			fmt.Println(style.Success.Render("  All tools installed"))
		}
	},
}