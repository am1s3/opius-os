package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/flash"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	customFlashPort    string
	customFlashChip    string
	customFlashBaud    int
	customFlashErase   bool
	customFlashDryRun  bool
	customFlashYes     bool
	customFlashBins    []string // Format: "file@address"
)

var deviceFlashCustomCmd = &cobra.Command{
	Use:   "flash-custom",
	Short: "Flash custom binary files with specific memory addresses",
	Long: `Flash one or more custom binary files to specific memory addresses.

This is useful for flashing partitioned firmware with bootloader, partition
table, and application binary at different addresses.

Format: file.bin@address (e.g., bootloader.bin@0x1000)

Examples:
  # Single binary
  opius device flash-custom --chip esp32 --bin firmware.bin@0x10000

  # Multiple binaries (partitioned firmware)
  opius device flash-custom --chip esp32 \
    --bin bootloader.bin@0x1000 \
    --bin partition-table.bin@0x8000 \
    --bin app.bin@0x10000

  # With erase
  opius device flash-custom --chip esp32 --erase \
    --bin bootloader.bin@0x1000 \
    --bin app.bin@0x10000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(customFlashBins) == 0 {
			return fmt.Errorf("at least one --bin argument is required")
		}

		// Parse binary entries
		var binaries []flash.BinaryEntry
		for _, entry := range customFlashBins {
			parts := strings.Split(entry, "@")
			if len(parts) != 2 {
				return fmt.Errorf("invalid binary format: %s (expected file@address)", entry)
			}

			path := strings.TrimSpace(parts[0])
			address := strings.TrimSpace(parts[1])

			if !strings.HasPrefix(address, "0x") {
				address = "0x" + address
			}

			binaries = append(binaries, flash.BinaryEntry{
				Path:    path,
				Address: address,
			})
		}

		// Detect or use specified port
		port := customFlashPort
		if port == "" {
			devs, err := device.List()
			if err != nil {
				return err
			}
			if len(devs) == 0 {
				return fmt.Errorf("no devices connected")
			}
			if len(devs) > 1 {
				return fmt.Errorf("multiple devices found, use --port")
			}
			port = devs[0].Port
		}

		target := flash.Target{
			Chip:     customFlashChip,
			Port:     port,
			Binaries: binaries,
			Baud:     customFlashBaud,
			Erase:    customFlashErase,
		}

		registry := flash.NewRegistry()
		driver, err := registry.Select(target)
		if err != nil {
			return err
		}

		plan, err := driver.Plan(context.Background(), target)
		if err != nil {
			return err
		}

		components.PrintSection("Custom Flash Plan")
		fmt.Printf("  %-12s %s\n", "Chip", style.Value.Render(customFlashChip))
		fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(port))
		fmt.Printf("  %-12s %d\n", "Binaries", len(binaries))

		for i, bin := range binaries {
			fmt.Printf("  %-12s %s @ %s\n", fmt.Sprintf("Binary %d", i+1), style.Dim.Render(bin.Path), style.Value.Render(bin.Address))
		}

		if customFlashErase {
			fmt.Printf("  %-12s %s\n", "Erase", style.Warning.Render("YES"))
		}

		fmt.Println()

		if customFlashDryRun {
			fmt.Println(style.Info.Render("  Dry run mode — no changes will be made"))
			fmt.Println()
			for i, s := range plan.Steps {
				phaseColor := phaseStyle(s.Phase)
				fmt.Printf("  %d. %s %s\n", i+1, phaseColor.Render(s.Phase), style.Value.Render(s.Name))
				fmt.Printf("     %s\n", style.Dim.Render(s.Description))
			}
			return nil
		}

		if !customFlashYes {
			fmt.Println(style.Warning.Render("  This will write to the device."))
			if !components.PromptConfirm("  Proceed?", true) {
				fmt.Println("  " + components.BadgeWarn() + " aborted by user")
				return nil
			}
		}

		fmt.Println()
		components.PrintSection("Flashing")

		evCh := make(chan flash.ProgressEvent)
		errCh := make(chan error, 1)

		go func() {
			errCh <- driver.Execute(context.Background(), plan, false, evCh)
		}()

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

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " custom flash completed successfully")
		return nil
	},
}

func init() {
	deviceFlashCustomCmd.Flags().StringVarP(&customFlashPort, "port", "p", "", "Serial port")
	deviceFlashCustomCmd.Flags().StringVarP(&customFlashChip, "chip", "c", "esp32", "Target chip")
	deviceFlashCustomCmd.Flags().IntVarP(&customFlashBaud, "baud", "b", 460800, "Serial speed")
	deviceFlashCustomCmd.Flags().BoolVarP(&customFlashErase, "erase", "e", false, "Erase flash before writing")
	deviceFlashCustomCmd.Flags().BoolVar(&customFlashDryRun, "dry-run", false, "Show plan without flashing")
	deviceFlashCustomCmd.Flags().BoolVarP(&customFlashYes, "yes", "y", false, "Skip confirmation")
	deviceFlashCustomCmd.Flags().StringArrayVar(&customFlashBins, "bin", []string{}, "Binary file with address (format: file@address)")

	_ = deviceFlashCustomCmd.MarkFlagRequired("bin")
	_ = deviceFlashCustomCmd.MarkFlagRequired("chip")
}