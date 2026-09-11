package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var (
	scriptFlashPort  string
	scriptFlashChip  string
	scriptFlashBoard string
)

var deviceFlashScriptCmd = &cobra.Command{
	Use:   "flash-script <file.ino|file.cpp>",
	Short: "Compile and flash an Arduino-style sketch",
	Long: `Compile and flash an Arduino-style .ino or .cpp sketch.

This command uses arduino-cli to compile the sketch for the
specified board and then flashes it to the device.

Requirements:
  - arduino-cli must be installed
  - Appropriate board core must be installed

Examples:
  opius device flash-script blink.ino
  opius device flash-script sketch.cpp --board esp32:esp32:esp32
  opius device flash-script sketch.ino --port /dev/cu.usbserial-14110`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sketchPath := args[0]

		if _, err := os.Stat(sketchPath); os.IsNotExist(err) {
			return fmt.Errorf("sketch file not found: %s", sketchPath)
		}

		if _, err := exec.LookPath("arduino-cli"); err != nil {
			return fmt.Errorf("arduino-cli not found. Install it with:\n" +
				"  brew install arduino-cli\n" +
				"  or: curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | sh")
		}

		board := scriptFlashBoard
		if board == "" {
			board = detectBoardFromChip(scriptFlashChip)
		}

		port := scriptFlashPort
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

		components.PrintSection("Flash Script")
		fmt.Printf("  %-12s %s\n", "Sketch", style.Value.Render(sketchPath))
		fmt.Printf("  %-12s %s\n", "Board", style.Value.Render(board))
		fmt.Printf("  %-12s %s\n", "Port", style.Value.Render(port))
		fmt.Println()

		components.PrintSection("Compiling")
		fmt.Println(style.Dim.Render("  Running arduino-cli compile..."))

		buildDir := filepath.Join(os.TempDir(), "opius-build")
		os.MkdirAll(buildDir, 0755)

		compileCmd := exec.Command("arduino-cli", "compile",
			"--fqbn", board,
			"--build-path", buildDir,
			sketchPath,
		)
		compileCmd.Stdout = os.Stdout
		compileCmd.Stderr = os.Stderr

		if err := compileCmd.Run(); err != nil {
			return fmt.Errorf("compilation failed: %w", err)
		}

		fmt.Println("  " + components.BadgeOK() + " Compilation successful")

		binPath := findCompiledBinary(buildDir)
		if binPath == "" {
			return fmt.Errorf("could not find compiled binary in %s", buildDir)
		}

		fmt.Println()
		components.PrintSection("Flashing")
		fmt.Printf("  %-12s %s\n", "Binary", style.Dim.Render(binPath))

		uploadCmd := exec.Command("arduino-cli", "upload",
			"--fqbn", board,
			"--port", port,
			"--input-dir", buildDir,
		)
		uploadCmd.Stdout = os.Stdout
		uploadCmd.Stderr = os.Stderr

		if err := uploadCmd.Run(); err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " Sketch flashed successfully!")

		return nil
	},
}

func detectBoardFromChip(chip string) string {
	switch strings.ToLower(chip) {
	case "esp32":
		return "esp32:esp32:esp32"
	case "esp32-s3":
		return "esp32:esp32:esp32s3"
	case "esp32-c3":
		return "esp32:esp32:esp32c3"
	case "esp8266":
		return "esp8266:esp8266:generic"
	case "rp2040":
		return "rp2040:rp2040:rpipico"
	case "avr", "arduino":
		return "arduino:avr:uno"
	default:
		return "arduino:avr:uno"
	}
}

func findCompiledBinary(buildDir string) string {
	var found string
	filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".bin" || ext == ".hex" {
				found = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	return found
}

func init() {
	deviceFlashScriptCmd.Flags().StringVarP(&scriptFlashPort, "port", "p", "", "Serial port")
	deviceFlashScriptCmd.Flags().StringVarP(&scriptFlashChip, "chip", "c", "", "Target chip")
	deviceFlashScriptCmd.Flags().StringVar(&scriptFlashBoard, "board", "", "Arduino FQBN board identifier")
}