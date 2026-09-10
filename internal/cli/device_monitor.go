package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/serialio"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	monitorPort      string
	monitorBaud      int
	monitorLog       string
	monitorNoTime    bool
	monitorReconnect bool
)

var deviceMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Open serial monitor with colored output",
	Long: `Open a serial monitor with automatic colorization of log lines.

Exit with Ctrl+C. Errors, warnings and boot messages are highlighted automatically.

Examples:
  opius device monitor
  opius device monitor --port /dev/cu.usbserial-14110 --baud 115200
  opius device monitor --log monitor.log`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := monitorPort
		if port == "" {
			devs, err := device.List()
			if err != nil {
				return err
			}
			if len(devs) == 0 {
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

		components.PrintSection("Serial Monitor")
		fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(port))
		fmt.Printf("  %-12s %d\n", "Baud", monitorBaud)
		if monitorLog != "" {
			fmt.Printf("  %-12s %s\n", "Log file", style.Value.Render(monitorLog))
		}
		fmt.Printf("  %-12s %s\n", "Exit", style.Dim.Render("Ctrl+C"))
		fmt.Println()

		// Handle Ctrl+C gracefully
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigCh
			fmt.Println()
			fmt.Println(style.Warning.Render("  Monitor stopped by user"))
			os.Exit(0)
		}()

		// Reconnect loop
		for {
			m, err := serialio.NewMonitor(port, monitorBaud, monitorLog)
			if err != nil {
				return err
			}
			m.ShowTime = !monitorNoTime

			if err := m.Connect(); err != nil {
				fmt.Println(components.BadgeErr(), "cannot open port:", err)
				if !monitorReconnect {
					return err
				}
				fmt.Println(style.Dim.Render("  retrying in 3s..."))
				time.Sleep(3 * time.Second)
				continue
			}

			if err := m.RunLoop(); err != nil {
				fmt.Println(components.BadgeWarn(), "connection lost:", err)
				m.Close()
				if !monitorReconnect {
					return nil
				}
				fmt.Println(style.Dim.Render("  reconnecting in 3s..."))
				time.Sleep(3 * time.Second)
				continue
			}

			m.Close()
			break
		}

		return nil
	},
}

func init() {
	deviceMonitorCmd.Flags().StringVarP(&monitorPort, "port", "p", "", "Serial port (auto-detected if only one)")
	deviceMonitorCmd.Flags().IntVarP(&monitorBaud, "baud", "b", 115200, "Serial speed")
	deviceMonitorCmd.Flags().StringVar(&monitorLog, "log", "", "Write raw output to log file")
	deviceMonitorCmd.Flags().BoolVar(&monitorNoTime, "no-time", false, "Disable timestamps")
	deviceMonitorCmd.Flags().BoolVar(&monitorReconnect, "reconnect", true, "Auto-reconnect on connection loss")
}