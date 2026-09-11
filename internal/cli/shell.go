package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Enter the Opius interactive shell",
	Long: `Enter the Opius interactive shell environment.

The prompt shows [(username)/Opius] and you can run Opius commands
without typing the 'opius' prefix.

Type 'exit' or press Ctrl+D to leave.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		components.PrintBanner()

		usr, err := user.Current()
		username := "user"
		if err == nil {
			username = usr.Username
		}

		prompt := style.Accent.Render(fmt.Sprintf("[(%s)/Opius]", username)) + " "
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println(style.Dim.Render("  Type 'help' for commands, 'exit' to leave."))
		fmt.Println()

		for {
			fmt.Print(prompt)
			if !scanner.Scan() {
				break
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			if line == "exit" || line == "quit" {
				fmt.Println(style.Dim.Render("  Goodbye!"))
				return nil
			}

			if line == "help" {
				printShellHelp()
				continue
			}

			if line == "clear" {
				fmt.Print("\033[H\033[2J")
				continue
			}

			// Run as opius subcommand
			parts := strings.Fields(line)
			runShellCommand(parts)
		}

		return scanner.Err()
	},
}

func printShellHelp() {
	fmt.Println()
	fmt.Println(style.Header.Render("  Available commands:"))
	fmt.Println()
	fmt.Printf("  %-24s %s\n", style.Value.Render("device list"), "List connected devices")
	fmt.Printf("  %-24s %s\n", style.Value.Render("device flash"), "Flash firmware")
	fmt.Printf("  %-24s %s\n", style.Value.Render("device monitor"), "Open serial monitor")
	fmt.Printf("  %-24s %s\n", style.Value.Render("device erase"), "Erase device flash")
	fmt.Printf("  %-24s %s\n", style.Value.Render("device reset"), "Reset device")
	fmt.Printf("  %-24s %s\n", style.Value.Render("pkg list"), "List available packages")
	fmt.Printf("  %-24s %s\n", style.Value.Render("pkg download <name>"), "Download package")
	fmt.Printf("  %-24s %s\n", style.Value.Render("pkg flash <name>"), "Flash package to device")
	fmt.Printf("  %-24s %s\n", style.Value.Render("doctor"), "Check system health")
	fmt.Printf("  %-24s %s\n", style.Value.Render("update"), "Update Opius OS")
	fmt.Printf("  %-24s %s\n", style.Value.Render("web"), "Start web interface")
	fmt.Printf("  %-24s %s\n", style.Value.Render("clear"), "Clear screen")
	fmt.Printf("  %-24s %s\n", style.Value.Render("exit"), "Exit shell")
	fmt.Println()
}

func runShellCommand(args []string) {
	if len(args) == 0 {
		return
	}

	// Build the opius command
	opiusArgs := append([]string{}, args...)

	cmd := exec.Command("opius", opiusArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Error already printed by the subcommand
	}
}