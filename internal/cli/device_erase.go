package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/flash"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	erasePort string
	eraseChip string
	eraseYes  bool
)

var deviceEraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "Fully erase device flash memory",
	Long: `Completely erase the flash memory of the connected device.

This is a DESTRUCTIVE operation. All firmware, data and settings
will be permanently removed from the device.

Examples:
  opius device erase --chip esp32
  opius device erase --chip esp32 --port /dev/cu.usbserial-14110`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := erasePort
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

		chip := eraseChip

		components.PrintSection("Erase Plan")
		fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(port))
		fmt.Printf("  %-12s %s\n", "Chip", style.Value.Render(chip))
		fmt.Println()
		fmt.Println(style.Warning.Render("  DESTRUCTIVE: this will completely erase flash memory."))
		fmt.Println(style.Dim.Render("  All data, firmware and settings will be removed."))

		if !eraseYes {
			if !components.PromptConfirm("  Erase device?", false) {
				fmt.Println("  " + components.BadgeWarn() + " aborted by user")
				return nil
			}
		}

		registry := flash.NewRegistry()
		driver, err := registry.Select(flash.Target{Chip: chip})
		if err != nil {
			return err
		}

		target := flash.Target{
			Chip:  chip,
			Port:  port,
			Erase: true,
		}

		plan, err := driver.Plan(context.Background(), target)
		if err != nil {
			return err
		}

		fmt.Println()
		components.PrintSection("Executing erase")

		evCh := make(chan flash.ProgressEvent)
		errCh := make(chan error, 1)
		go func() {
			errCh <- driver.Execute(context.Background(), plan, false, evCh)
		}()

		for ev := range evCh {
			fmt.Printf("  %s %s\n", style.Value.Render(ev.Step), style.Dim.Render(ev.Message))
		}

		if err := <-errCh; err != nil {
			fmt.Println()
			fmt.Println("  " + components.BadgeErr() + " " + err.Error())
			return err
		}

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " flash erased successfully")
		return nil
	},
}

func init() {
	deviceEraseCmd.Flags().StringVarP(&erasePort, "port", "p", "", "Serial port")
	deviceEraseCmd.Flags().StringVarP(&eraseChip, "chip", "c", "esp32", "Chip type")
	deviceEraseCmd.Flags().BoolVarP(&eraseYes, "yes", "y", false, "Skip confirmation")
}