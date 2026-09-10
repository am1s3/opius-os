package web

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opius-os/opius/internal/device"
	"go.bug.st/serial"
)

var serialUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	activeSerialSessions = make(map[string]*SerialMonitorSession)
	serialSessionsMu     sync.Mutex
)

type SerialMonitorSession struct {
	mu     sync.Mutex
	port   serial.Port
	ws     *websocket.Conn
	closed bool
	name   string
}

type SerialMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

func closeSession(name string) {
	serialSessionsMu.Lock()
	defer serialSessionsMu.Unlock()

	if sess, ok := activeSerialSessions[name]; ok {
		sess.mu.Lock()
		if !sess.closed {
			sess.closed = true
			if sess.port != nil {
				sess.port.Close()
			}
			if sess.ws != nil {
				sess.ws.Close()
			}
		}
		sess.mu.Unlock()
		delete(activeSerialSessions, name)
	}
}

// findPortHolder tries to find what process is holding the port
func findPortHolder(portName string) string {
	cmd := exec.Command("lsof", portName)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		// Return first few lines as info
		info := strings.Join(lines[1:3], " | ")
		return info
	}
	return ""
}

func (s *Server) handleSerialMonitorWS(w http.ResponseWriter, r *http.Request) {
	portName := r.URL.Query().Get("port")
	baudStr := r.URL.Query().Get("baud")

	if portName == "" {
		http.Error(w, "port parameter required", http.StatusBadRequest)
		return
	}

	baud := 115200
	if baudStr != "" {
		if parsed, err := strconv.Atoi(baudStr); err == nil && parsed > 0 {
			baud = parsed
		}
	}

	// Close any existing session for this port
	closeSession(portName)

	// Longer delay to ensure port is released
	time.Sleep(300 * time.Millisecond)

	// Verify device exists
	devices, err := device.ListAll()
	if err != nil {
		http.Error(w, "failed to list devices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	found := false
	for _, d := range devices {
		if d.Port == portName {
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "device not found: "+portName, http.StatusNotFound)
		return
	}

	conn, err := serialUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("serial monitor ws upgrade failed: %v", err)
		return
	}

	// Try to open serial port with multiple retries
	var serPort serial.Port
	maxRetries := 5
	openErr := error(nil)

	for i := 0; i < maxRetries; i++ {
		mode := &serial.Mode{BaudRate: baud}
		serPort, openErr = serial.Open(portName, mode)
		if openErr == nil {
			break
		}

		// Check if port is held by another process
		holder := findPortHolder(portName)
		if holder != "" {
			log.Printf("Port %s held by: %s", portName, holder)
		}

		if i < maxRetries-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if openErr != nil {
		holder := findPortHolder(portName)
		errMsg := "cannot open port: " + openErr.Error()
		if holder != "" {
			errMsg += "\nPort held by: " + holder
		}
		errMsg += "\nTry closing Arduino IDE or other serial monitors."

		conn.WriteJSON(SerialMessage{
			Type:    "error",
			Content: errMsg,
			Time:    time.Now().Format("15:04:05.000"),
		})
		conn.Close()
		return
	}

	session := &SerialMonitorSession{
		port: serPort,
		ws:   conn,
		name: portName,
	}

	serialSessionsMu.Lock()
	activeSerialSessions[portName] = session
	serialSessionsMu.Unlock()

	defer func() {
		closeSession(portName)
	}()

	conn.WriteJSON(SerialMessage{
		Type:    "status",
		Content: fmt.Sprintf("connected to %s at %d baud", portName, baud),
		Time:    time.Now().Format("15:04:05.000"),
	})

	// Read from serial -> send to ws
	go func() {
		scanner := bufio.NewScanner(serPort)
		for scanner.Scan() {
			session.mu.Lock()
			if session.closed {
				session.mu.Unlock()
				return
			}
			session.mu.Unlock()

			line := scanner.Text()
			msg := SerialMessage{
				Type:    "data",
				Content: line,
				Time:    time.Now().Format("15:04:05.000"),
			}

			session.mu.Lock()
			if !session.closed {
				if err := conn.WriteJSON(msg); err != nil {
					session.mu.Unlock()
					return
				}
			}
			session.mu.Unlock()
		}
	}()

	// Read from ws -> write to serial
	for {
		var cmd struct {
			Action  string `json:"action"`
			Content string `json:"content"`
		}

		if err := conn.ReadJSON(&cmd); err != nil {
			break
		}

		if cmd.Action == "close" {
			break
		}

		if cmd.Action == "send" {
			session.mu.Lock()
			if !session.closed && serPort != nil {
				serPort.Write([]byte(cmd.Content + "\n"))
			}
			session.mu.Unlock()
		}
	}
}