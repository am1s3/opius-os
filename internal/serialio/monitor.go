package serialio

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/opius-os/opius/internal/ui/style"
	"go.bug.st/serial"
)

type Monitor struct {
	Port     string
	Baud     int
	ShowTime bool
	LogFile  *os.File
	port     serial.Port
}

func NewMonitor(port string, baud int, logPath string) (*Monitor, error) {
	m := &Monitor{
		Port:     port,
		Baud:     baud,
		ShowTime: true,
	}

	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		m.LogFile = f
	}

	return m, nil
}

func (m *Monitor) Connect() error {
	mode := &serial.Mode{BaudRate: m.Baud}
	p, err := serial.Open(m.Port, mode)
	if err != nil {
		return err
	}
	m.port = p
	return nil
}

func (m *Monitor) Close() {
	if m.port != nil {
		m.port.Close()
		m.port = nil
	}
	if m.LogFile != nil {
		m.LogFile.Close()
		m.LogFile = nil
	}
}

func (m *Monitor) RunLoop() error {
	scanner := bufio.NewScanner(m.port)
	for scanner.Scan() {
		line := scanner.Text()
		m.printLine(line)
	}
	return scanner.Err()
}

func (m *Monitor) printLine(line string) {
	ts := time.Now().Format("15:04:05.000")
	colored := colorize(line)

	if m.ShowTime {
		fmt.Printf("%s %s\n", style.Dim.Render(ts), colored)
	} else {
		fmt.Println(colored)
	}

	if m.LogFile != nil {
		fmt.Fprintf(m.LogFile, "[%s] %s\n", ts, line)
	}
}

func colorize(line string) string {
	lower := strings.ToLower(line)

	switch {
	case strings.Contains(lower, "error"),
		strings.Contains(lower, "fail"),
		strings.Contains(lower, "panic"),
		strings.Contains(lower, "fatal"):
		return style.Error.Render(line)

	case strings.Contains(lower, "warn"):
		return style.Warning.Render(line)

	case strings.Contains(lower, "info"),
		strings.Contains(lower, "boot"),
		strings.Contains(lower, "connect"):
		return style.Info.Render(line)

	case strings.Contains(lower, "ok"),
		strings.Contains(lower, "success"),
		strings.Contains(lower, "done"):
		return style.Success.Render(line)

	default:
		return line
	}
}