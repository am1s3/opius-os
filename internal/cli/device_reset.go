package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
	"go.bug.st/serial"
)

var (
	resetPort string
)

var deviceResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset device via DTR/RTS signals",
	Long: `Send a reset signal to the connected device by toggling
the DTR and RTS control lines of the serial port.

This works for most Arduino and ESP development boards.

Examples:
  opius device reset
  opius device reset --port /dev/cu.usbserial-14110`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := resetPort
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

		components.PrintSection("Device Reset")
		fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(port))
		fmt.Println(style.Info.Render("  Sending reset signal via DTR/RTS..."))

		p, err := serial.Open(port, &serial.Mode{BaudRate: 115200})
		if err != nil {
			return fmt.Errorf("cannot open port: %w", err)
		}
		defer p.Close()

		// Standard Arduino/ESP reset sequence
		p.SetDTR(false)
		p.SetRTS(true)
		time.Sleep(100 * time.Millisecond)
		p.SetRTS(false)
		time.Sleep(100 * time.Millisecond)

		fmt.Println("  " + components.BadgeOK() + " device reset signal sent")
		return nil
	},
}

func init() {
	deviceResetCmd.Flags().StringVarP(&resetPort, "port", "p", "", "Serial port")
}