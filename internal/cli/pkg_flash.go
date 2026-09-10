package cli

import (
	"context"
	"fmt"
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
	pkgFlashErase   bool
	pkgFlashDryRun  bool
	pkgFlashYes     bool
)

var pkgFlashCmd = &cobra.Command{
	Use:   "flash <name>",
	Short: "Install a package and flash it to a device",
	Long: `Install a package from the registry and flash it to a connected device.

This command combines 'pkg install' and 'device flash' into one step.
It automatically finds the firmware file in the package manifest and
flashes it to the detected or specified device.

Examples:
  opius pkg flash bruce
  opius pkg flash bruce --port /dev/cu.usbserial-14110
  opius pkg flash bruce --erase
  opius pkg flash bruce --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Step 1: Find package in registry
		regPath, err := pkg.DefaultRegistryPath()
		if err != nil {
			return err
		}
		registry := pkg.NewRegistry(regPath)

		m, err := registry.Find(name)
		if err != nil {
			return err
		}

		// Step 2: Install package if not already installed
		storePath, err := pkg.DefaultStorePath()
		if err != nil {
			return err
		}
		store := pkg.NewStore(storePath)

		installDir := store.PackageDir(m.Name, m.Version)
		installed, err := store.List()
		if err != nil {
			return err
		}

		alreadyInstalled := false
		for _, ip := range installed {
			if ip.Name == m.Name && ip.Version == m.Version {
				alreadyInstalled = true
				break
			}
		}

		if !alreadyInstalled {
			components.PrintSection("Step 1: Install Package")
			fmt.Printf("  %-12s %s\n", "Package", style.Value.Render(m.Name))
			fmt.Printf("  %-12s %s\n", "Version", style.Value.Render(m.Version))
			fmt.Printf("  %-12s %s\n", "Target", style.Dim.Render(installDir))

			if err := store.Install(m, regPath); err != nil {
				return fmt.Errorf("install failed: %w", err)
			}

			fmt.Println("  " + components.BadgeOK() + " package installed")
			fmt.Println()
		} else {
			fmt.Println(style.Dim.Render("  Package already installed, skipping install step"))
			fmt.Println()
		}

		// Step 3: Find firmware file in manifest
		var firmwareFile *pkg.File
		for i := range m.Files {
			if m.Files[i].Type == "firmware" {
				firmwareFile = &m.Files[i]
				break
			}
		}

		if firmwareFile == nil {
			return fmt.Errorf("no firmware file found in package manifest (looking for type='firmware')")
		}

		// Step 4: Determine chip from manifest
		chip := ""
		if len(m.Chips) > 0 {
			chip = m.Chips[0]
		} else {
			chip = "esp32" // default fallback
		}

		// Step 5: Detect or use specified port
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
				fmt.Println("  " + components.BadgeWarn() + " multiple devices found:")
				for _, d := range devs {
					fmt.Printf("    %s · %s\n", style.Dim.Render(d.Port), style.Value.Render(d.Chip))
				}
				return fmt.Errorf("use --port to select one")
			}
			port = devs[0].Port
		}

		// Step 6: Build flash plan
		firmwarePath := filepath.Join(installDir, firmwareFile.Path)

		target := flash.Target{
			Chip:  chip,
			Port:  port,
			File:  firmwarePath,
			Erase: pkgFlashErase,
		}

		registry2 := flash.NewRegistry()
		driver, err := registry2.Select(target)
		if err != nil {
			return err
		}

		plan, err := driver.Plan(context.Background(), target)
		if err != nil {
			return err
		}

		components.PrintSection("Step 2: Flash Plan")
		fmt.Printf("  %-12s %s\n", "Package", style.Value.Render(m.Name+"@"+m.Version))
		fmt.Printf("  %-12s %s\n", "Firmware", style.Dim.Render(firmwareFile.Path))
		fmt.Printf("  %-12s %s\n", "Chip", style.Value.Render(chip))
		fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(port))
		if pkgFlashErase {
			fmt.Printf("  %-12s %s\n", "Erase", style.Warning.Render("YES"))
		}
		fmt.Printf("  %-12s %d\n", "Steps", len(plan.Steps))
		fmt.Println()

		if pkgFlashDryRun {
			fmt.Println(style.Info.Render("  Dry run mode — no changes will be made"))
			fmt.Println()
			for i, s := range plan.Steps {
				phaseColor := phaseStyle(s.Phase)
				fmt.Printf("  %d. %s %s\n", i+1, phaseColor.Render(s.Phase), style.Value.Render(s.Name))
				fmt.Printf("     %s\n", style.Dim.Render(s.Description))
			}
			return nil
		}

		if !pkgFlashYes {
			fmt.Println(style.Warning.Render("  This will write firmware to the device."))
			if !components.PromptConfirm("  Proceed?", true) {
				fmt.Println("  " + components.BadgeWarn() + " aborted by user")
				return nil
			}
		}

		// Step 7: Execute flash
		fmt.Println()
		components.PrintSection("Step 3: Flashing")

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
		fmt.Println("  " + components.BadgeOK() + " package flashed successfully")
		fmt.Println()
		fmt.Println(style.Dim.Render("  Open serial monitor with:"))
		fmt.Println(style.Dim.Render("    opius device monitor --port " + port))

		return nil
	},
}

func init() {
	pkgFlashCmd.Flags().StringVarP(&pkgFlashPort, "port", "p", "", "Serial port (auto-detected if only one)")
	pkgFlashCmd.Flags().BoolVarP(&pkgFlashErase, "erase", "e", false, "Full erase before write")
	pkgFlashCmd.Flags().BoolVar(&pkgFlashDryRun, "dry-run", false, "Only show the plan, do not flash")
	pkgFlashCmd.Flags().BoolVarP(&pkgFlashYes, "yes", "y", false, "Skip confirmation prompt")
}