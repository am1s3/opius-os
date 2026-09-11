package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/flash"
	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	pkgFlashPort    string
	pkgFlashChip    string
	pkgFlashVersion string
	pkgFlashErase   bool
	pkgFlashYes     bool
)

var pkgFlashCmd = &cobra.Command{
	Use:   "flash <name>",
	Short: "Download and flash a package to device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		reg := newPkgRegistry()
		if err := reg.Sync(); err != nil {
			return fmt.Errorf("failed to sync registry: %w", err)
		}

		m, err := reg.Find(name)
		if err != nil {
			return err
		}

		if len(m.Versions) == 0 {
			return fmt.Errorf("no versions available for %s", name)
		}

		// Select version
		var version *pkg.PackageVersion
		if pkgFlashVersion != "" {
			for i := range m.Versions {
				if m.Versions[i].Version == pkgFlashVersion {
					version = &m.Versions[i]
					break
				}
			}
			if version == nil {
				return fmt.Errorf("version %s not found. Available: %s", pkgFlashVersion, listVersionNames(m.Versions))
			}
		} else {
			version = &m.Versions[0] // Latest
		}

		// Determine chip
		chip := pkgFlashChip
		if chip == "" && len(m.Chips) > 0 {
			chip = m.Chips[0]
		}
		if chip == "" {
			chip = "esp32"
		}

		// Determine port
		port := pkgFlashPort
		if port == "" {
			devs, err := device.List()
			if err != nil {
				return err
			}
			if len(devs) == 0 {
				return fmt.Errorf("no devices connected")
			}
			if len(devs) > 1 {
				fmt.Println("  " + components.BadgeWarn() + " Multiple devices found:")
				for _, d := range devs {
					fmt.Printf("    %s · %s\n", style.Dim.Render(d.Port), style.Value.Render(d.Chip))
				}
				return fmt.Errorf("use --port to select one")
			}
			port = devs[0].Port
		}

		// Download first
		home, _ := os.UserHomeDir()
		downloadDir := filepath.Join(home, ".opius", "downloads")
		os.MkdirAll(downloadDir, 0755)

		filename := filepath.Base(version.DownloadURL)
		if filename == "" || filename == "/" {
			filename = m.Name + "-" + version.Version + ".bin"
		}
		firmwarePath := filepath.Join(downloadDir, filename)

		components.PrintSection("Flash Plan")
		fmt.Printf("  %-12s %s\n", "Package", style.Value.Render(m.Name+"@"+version.Version))
		fmt.Printf("  %-12s %s\n", "Chip", style.Value.Render(chip))
		fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(port))
		fmt.Printf("  %-12s %s\n", "Firmware", style.Dim.Render(firmwarePath))
		if pkgFlashErase {
			fmt.Printf("  %-12s %s\n", "Erase", style.Warning.Render("YES"))
		}
		fmt.Println()

		if !pkgFlashYes {
			fmt.Println(style.Warning.Render("  This will write firmware to the device."))
			if !components.PromptConfirm("  Proceed?", true) {
				fmt.Println("  " + components.BadgeWarn() + " aborted")
				return nil
			}
		}

		// Download if not exists
		if _, err := os.Stat(firmwarePath); os.IsNotExist(err) {
			fmt.Println()
			components.PrintSection("Downloading firmware")
			if err := downloadWithProgress(version.DownloadURL, firmwarePath, version.Size); err != nil {
				return fmt.Errorf("download failed: %w", err)
			}

			if version.SHA256 != "" {
				fmt.Print("  Verifying SHA256... ")
				if err := verifyFileSHA256(firmwarePath, version.SHA256); err != nil {
					fmt.Println(components.BadgeErr())
					return fmt.Errorf("checksum mismatch: %w", err)
				}
				fmt.Println(components.BadgeOK())
			}
		}

		// Flash
		fmt.Println()
		components.PrintSection("Flashing")

		target := flash.Target{
			Chip:  chip,
			Port:  port,
			File:  firmwarePath,
			Erase: pkgFlashErase,
		}

		flashReg := flash.NewRegistry()
		driver, err := flashReg.Select(target)
		if err != nil {
			return err
		}

		plan, err := driver.Plan(context.Background(), target)
		if err != nil {
			return err
		}

		evCh := make(chan flash.ProgressEvent)
		errCh := make(chan error, 1)

		go func() {
			errCh <- driver.Execute(context.Background(), plan, false, evCh)
		}()

		for ev := range evCh {
			phaseColor := phaseStyle(ev.Phase)
			fmt.Printf("  %s %s: %s\n",
				phaseColor.Render("["+ev.Phase+"]"),
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
		fmt.Println("  " + components.BadgeOK() + " Flash completed successfully!")
		fmt.Println()
		fmt.Println(style.Dim.Render("  Monitor with: opius device monitor --port " + port))

		return nil
	},
}

func init() {
	pkgFlashCmd.Flags().StringVarP(&pkgFlashPort, "port", "p", "", "Serial port")
	pkgFlashCmd.Flags().StringVarP(&pkgFlashChip, "chip", "c", "", "Target chip")
	pkgFlashCmd.Flags().StringVarP(&pkgFlashVersion, "version", "v", "", "Specific version (default: latest)")
	pkgFlashCmd.Flags().BoolVarP(&pkgFlashErase, "erase", "e", false, "Erase before flash")
	pkgFlashCmd.Flags().BoolVarP(&pkgFlashYes, "yes", "y", false, "Skip confirmation")
}