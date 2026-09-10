package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/flash"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	flashPort    string
	flashChip    string
	flashFile    string
	flashBaud    int
	flashErase   bool
	flashDryRun  bool
	flashConfirm bool
)

var deviceFlashCmd = &cobra.Command{
	Use:   "flash",
	Short: "Flash a firmware binary to a device",
	Long: `Flash a firmware binary to a connected microcontroller.

Examples:
  opius device flash --file firmware.bin --chip esp32
  opius device flash --file firmware.bin --chip esp32-s3 --port /dev/cu.usbserial-0001
  opius device flash --file firmware.bin --chip esp32 --dry-run
  opius device flash --file firmware.bin --chip esp32 --erase`,
	RunE: func(cmd *cobra.Command, args []string) error {
		components.PrintSection("Flash Plan")

		// Resolve port if not specified
		port := flashPort
		if port == "" {
			devs, err := device.List()
			if err != nil {
				return fmt.Errorf("device list failed: %w", err)
			}
			if len(devs) == 0 {
				fmt.Println("  " + components.BadgeErr() + " no devices found")
				return fmt.Errorf("no devices connected")
			}
			if len(devs) > 1 {
				fmt.Println("  " + components.BadgeWarn() + " multiple devices found, specify one:")
				for _, d := range devs {
					fmt.Printf("    %s · %s\n", style.Dim.Render(d.Port), style.Value.Render(d.Chip))
				}
				return fmt.Errorf("use --port to select one")
			}
			port = devs[0].Port
		}

		if flashFile == "" {
			return fmt.Errorf("--file is required")
		}
		if _, err := os.Stat(flashFile); err != nil {
			return fmt.Errorf("firmware file not found: %s", flashFile)
		}
		if flashChip == "" {
			return fmt.Errorf("--chip is required")
		}

		target := flash.Target{
			Chip:  strings.ToLower(flashChip),
			Port:  port,
			File:  flashFile,
			Baud:  flashBaud,
			Erase: flashErase,
		}

		registry := flash.NewRegistry()
		if !registry.HasDriver() {
			fmt.Println("  " + components.BadgeErr() + " no flash drivers available")
			fmt.Println()
			fmt.Println(style.Info.Render("  Install esptool:"))
			fmt.Println(style.Dim.Render("    sudo port install py312-esptool"))
			fmt.Println(style.Dim.Render("    # or: pip install esptool"))
			return fmt.Errorf("no drivers")
		}

		driver, err := registry.Select(target)
		if err != nil {
			fmt.Println("  " + components.BadgeErr() + " " + err.Error())
			return err
		}

		plan, err := driver.Plan(context.Background(), target)
		if err != nil {
			fmt.Println("  " + components.BadgeErr() + " plan failed: " + err.Error())
			return err
		}

		printPlan(plan, flashDryRun)

		if flashDryRun {
			fmt.Println()
			fmt.Println(style.Info.Render("  Dry run complete. No changes were made to the device."))
			return nil
		}

		if !flashConfirm {
			fmt.Println()
			fmt.Println(style.Warning.Render("  Destructive action: this will write to the device."))
			if !components.PromptConfirm("  Proceed?", false) {
				fmt.Println("  " + components.BadgeWarn() + " aborted by user")
				return nil
			}
		}

		fmt.Println()
		components.PrintSection("Executing")

		evCh := make(chan flash.ProgressEvent)
		errCh := make(chan error, 1)

		go func() {
			errCh <- driver.Execute(context.Background(), plan, false, evCh)
		}()

		start := time.Now()
		for ev := range evCh {
			phaseColor := phaseStyle(ev.Phase)
			fmt.Printf("  %s %s: %s\n",
				phaseColor.Render("[" + ev.Phase + "]"),
				style.Value.Render(ev.Step),
				style.Dim.Render(ev.Message),
			)
		}

		if err := <-errCh; err != nil {
			fmt.Println()
			fmt.Println("  " + components.BadgeErr() + " " + err.Error())
			return err
		}

		elapsed := time.Since(start).Round(time.Millisecond)
		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " flashed successfully in " + style.Value.Render(elapsed.String()))
		fmt.Println()
		fmt.Println(style.Dim.Render("  Open serial monitor with:"))
		fmt.Println(style.Dim.Render("    opius device monitor --port " + port))
		return nil
	},
}

func printPlan(p *flash.Plan, dryRun bool) {
	fmt.Printf("  %-12s %s\n", "Driver", style.Value.Render(p.Driver))
	fmt.Printf("  %-12s %s\n", "Tool", style.Dim.Render(p.Tool))
	fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(p.Target.Port))
	fmt.Printf("  %-12s %s\n", "Chip", style.Value.Render(p.Target.Chip))
	fmt.Printf("  %-12s %s\n", "File", style.Dim.Render(p.Target.File))
	if p.Target.Erase {
		fmt.Printf("  %-12s %s\n", "Erase", style.Warning.Render("YES (full)"))
	}
	if dryRun {
		fmt.Printf("  %-12s %s\n", "Mode", style.Accent.Render("dry-run"))
	}
	fmt.Println()
	for i, s := range p.Steps {
		phaseColor := phaseStyle(s.Phase)
		fmt.Printf("  %d. %s %s\n", i+1, phaseColor.Render(s.Phase), style.Value.Render(s.Name))
		fmt.Printf("     %s\n", style.Dim.Render(s.Description))
		fmt.Printf("     %s\n", style.Dim.Render("$ " + strings.Join(s.Command, " ")))
	}
}

func phaseStyle(phase string) lipgloss.Style {
	switch phase {
	case "prepare":
		return style.Info
	case "erase":
		return style.Warning
	case "write":
		return style.Accent
	case "verify":
		return style.Success
	case "finalize":
		return style.Success
	}
	return style.Dim
}

func init() {
	deviceFlashCmd.Flags().StringVarP(&flashPort, "port", "p", "", "Serial port (auto-detected if only one device)")
	deviceFlashCmd.Flags().StringVarP(&flashChip, "chip", "c", "", "Target chip: esp32, esp32-s3, esp32-c3, esp32-c2, esp8266")
	deviceFlashCmd.Flags().StringVarP(&flashFile, "file", "f", "", "Path to firmware binary")
	deviceFlashCmd.Flags().IntVarP(&flashBaud, "baud", "b", 0, "Serial speed (default 460800)")
	deviceFlashCmd.Flags().BoolVarP(&flashErase, "erase", "e", false, "Full erase before write")
	deviceFlashCmd.Flags().BoolVar(&flashDryRun, "dry-run", false, "Only show the plan, do not touch the device")
	deviceFlashCmd.Flags().BoolVarP(&flashConfirm, "yes", "y", false, "Skip confirmation prompt")

	_ = deviceFlashCmd.MarkFlagRequired("file")
	_ = deviceFlashCmd.MarkFlagRequired("chip")
}