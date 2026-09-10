package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
	"github.com/opius-os/opius/internal/web"
)

var (
	webHost string
	webPort int
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the Opius OS web interface",
	Long: `Start a local web server with a graphical interface for managing
devices and packages through your browser.

The web UI provides:
  - Device list and status
  - Package browser
  - One-click flashing
  - Real-time monitoring (coming soon)

Default: http://127.0.0.1:8420

Examples:
  opius web
  opius web --port 9000
  opius web --host 0.0.0.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		components.PrintBanner()

		fmt.Printf("  %-12s %s\n", "Host", style.Value.Render(webHost))
		fmt.Printf("  %-12s %d\n", "Port", webPort)
		fmt.Printf("  %-12s %s\n", "URL", style.Info.Render(fmt.Sprintf("http://%s:%d", webHost, webPort)))
		fmt.Println()

		server := web.NewServer(webHost, webPort)

		// Handle Ctrl+C gracefully
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-sigCh
			fmt.Println()
			fmt.Println(style.Warning.Render("  Stopping web server..."))
			server.Stop()
			os.Exit(0)
		}()

		return server.Start()
	},
}

func init() {
	webCmd.Flags().StringVar(&webHost, "host", "127.0.0.1", "Host to bind to")
	webCmd.Flags().IntVarP(&webPort, "port", "p", 8420, "Port to listen on")
}