package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	devices "github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	deviceListJSON bool
	deviceListAll  bool
	deviceInfoJSON bool
	deviceInfoPort string
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Manage microcontroller devices",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected serial/USB devices",
	Run: func(cmd *cobra.Command, args []string) {
		var list []devices.Device
		var err error

		if deviceListAll {
			list, err = devices.ListAll()
		} else {
			list, err = devices.List()
		}

		if err != nil {
			fmt.Println(components.BadgeErr(), err)
			return
		}

		if deviceListJSON {
			out, _ := json.MarshalIndent(list, "", "  ")
			fmt.Println(string(out))
			return
		}

		components.PrintSection("Devices")

		if len(list) == 0 {
			if deviceListAll {
				fmt.Println("  " + components.BadgeWarn() + " no serial/USB devices found (including Bluetooth)")
			} else {
				fmt.Println("  " + components.BadgeWarn() + " no USB devices found")
				fmt.Println()
				fmt.Println(style.Dim.Render("  hints:"))
				fmt.Println(style.Dim.Render("    - check USB cable"))
				fmt.Println(style.Dim.Render("    - some cables are charge-only"))
				fmt.Println(style.Dim.Render("    - put board into bootloader mode"))
				fmt.Println(style.Dim.Render("    - use --all to show Bluetooth and virtual ports"))
			}
			return
		}

		fmt.Printf("  %-3s %-24s %-10s %-28s %-10s %-10s\n",
			style.Header.Render("ID"),
			style.Header.Render("PORT"),
			style.Header.Render("CHIP"),
			style.Header.Render("BOARD"),
			style.Header.Render("MODE"),
			style.Header.Render("STATUS"),
		)

		for _, d := range list {
			fmt.Printf("  %-3d %-24s %-10s %-28s %-10s %-10s\n",
				d.ID,
				style.Dim.Render(d.Port),
				style.Value.Render(d.Chip),
				truncate(d.Board, 28),
				d.Mode,
				statusBadge(d.Status),
			)
		}

		fmt.Println()
		fmt.Printf("  %d device(s) found\n", len(list))
		note := devices.PreferredPortNote()
		if note != "" {
			fmt.Println(style.Dim.Render("  " + note))
		}
		if !deviceListAll {
			fmt.Println(style.Dim.Render("  Use --all to show Bluetooth and virtual ports"))
		}
	},
}

var deviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show detailed information about a device",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := devices.ListAll()
		if err != nil {
			fmt.Println(components.BadgeErr(), err)
			return
		}

		if len(list) == 0 {
			fmt.Println(components.BadgeWarn(), "no devices found")
			return
		}

		var selected *devices.Device

		if deviceInfoPort != "" {
			d, err := devices.FindByPort(deviceInfoPort)
			if err != nil {
				fmt.Println(components.BadgeErr(), err)
				return
			}
			selected = d
		} else {
			if len(list) > 1 {
				fmt.Println(components.BadgeWarn(), "multiple devices found")
				fmt.Println("  specify one with:")
				fmt.Println(style.Dim.Render("    opius device info --port <port>"))
				fmt.Println()
				for _, d := range list {
					fmt.Printf("  %-3d %s\n", d.ID, d.Port)
				}
				return
			}
			selected = &list[0]
		}

		if deviceInfoJSON {
			out, _ := json.MarshalIndent(selected, "", "  ")
			fmt.Println(string(out))
			return
		}

		components.PrintSection("Device Info")
		fmt.Printf("  %-14s %s\n", "Port", style.Value.Render(selected.Port))
		fmt.Printf("  %-14s %v\n", "USB", selected.IsUSB)
		fmt.Printf("  %-14s %s\n", "VID", emptyDash(selected.VID))
		fmt.Printf("  %-14s %s\n", "PID", emptyDash(selected.PID))
		fmt.Printf("  %-14s %s\n", "Serial", emptyDash(selected.SerialNumber))
		fmt.Printf("  %-14s %s\n", "Product", emptyDash(selected.Product))
		fmt.Printf("  %-14s %s\n", "Chip", style.Value.Render(selected.Chip))
		fmt.Printf("  %-14s %s\n", "Board", style.Value.Render(selected.Board))
		fmt.Printf("  %-14s %s\n", "Mode", selected.Mode)
		fmt.Printf("  %-14s %s\n", "Status", statusBadge(selected.Status))
		if selected.Hint != "" {
			fmt.Printf("  %-14s %s\n", "Hint", style.Dim.Render(selected.Hint))
		}
	},
}

var deviceDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose device connection issues",
	Run: func(cmd *cobra.Command, args []string) {
		components.PrintBanner()
		components.PrintSection("Device Doctor")

		list, err := devices.List()
		if err != nil {
			fmt.Println("  "+components.BadgeErr(), err)
			return
		}

		if len(list) == 0 {
			fmt.Println("  " + components.BadgeWarn() + " no devices detected")
			fmt.Println()
			fmt.Println(style.Info.Render("  Checklist:"))
			fmt.Println(style.Dim.Render("    1. Try another USB cable"))
			fmt.Println(style.Dim.Render("    2. Use a data cable, not charge-only"))
			fmt.Println(style.Dim.Render("    3. Try another USB port"))
			fmt.Println(style.Dim.Render("    4. Hold BOOT and press RESET for ESP boards"))
			fmt.Println(style.Dim.Render("    5. Hold BOOTSEL while plugging RP2040 boards"))
			fmt.Println(style.Dim.Render("    6. Close Arduino IDE / PlatformIO serial monitor"))
			return
		}

		fmt.Println("  "+components.BadgeOK(), "serial/USB detection works")
		fmt.Printf("  %-12s %d\n", "devices", len(list))

		fmt.Println()
		for _, d := range list {
			fmt.Printf("  %s %s · %s · %s\n",
				statusBadge(d.Status),
				style.Dim.Render(d.Port),
				style.Value.Render(d.Chip),
				truncate(d.Board, 32),
			)
			if d.Hint != "" {
				fmt.Println(style.Dim.Render("       " + d.Hint))
			}
		}

		note := devices.PreferredPortNote()
		if note != "" {
			fmt.Println()
			fmt.Println(style.Info.Render("  macOS note:"))
			fmt.Println(style.Dim.Render("    " + note))
		}
	},
}

func init() {
	deviceListCmd.Flags().BoolVar(&deviceListJSON, "json", false, "Output devices as JSON")
	deviceListCmd.Flags().BoolVar(&deviceListAll, "all", false, "Show all ports including Bluetooth and virtual")

	deviceInfoCmd.Flags().BoolVar(&deviceInfoJSON, "json", false, "Output device info as JSON")
	deviceInfoCmd.Flags().StringVarP(&deviceInfoPort, "port", "p", "", "Device port")

	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceInfoCmd)
	deviceCmd.AddCommand(deviceDoctorCmd)
	deviceCmd.AddCommand(deviceFlashCmd)
	deviceCmd.AddCommand(deviceFlashCustomCmd)
	deviceCmd.AddCommand(deviceMonitorCmd)
	deviceCmd.AddCommand(deviceEraseCmd)
	deviceCmd.AddCommand(deviceResetCmd)
}

func statusBadge(status string) string {
	switch status {
	case "ready":
		return style.Success.Render("[READY]")
	case "boot":
		return style.Accent.Render("[BOOT]")
	case "locked":
		return style.Error.Render("[LOCK]")
	case "unknown":
		return style.Warning.Render("[UNKWN]")
	default:
		return style.Dim.Render("[----]")
	}
}

func emptyDash(v string) string {
	if v == "" {
		return style.Dim.Render("-")
	}
	return style.Value.Render(v)
}

func truncate(v string, max int) string {
	if len(v) <= max {
		return v
	}
	if max <= 1 {
		return v[:max]
	}
	return v[:max-1] + "…"
}