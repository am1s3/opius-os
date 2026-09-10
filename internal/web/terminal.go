package web

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// TerminalSession represents a single pty session attached to a WebSocket.
type TerminalSession struct {
	id   string
	cmd  *exec.Cmd
	ptmx *os.File
	mu   sync.Mutex
	ws   *websocket.Conn
}

var sessions = struct {
	sync.RWMutex
	m map[string]*TerminalSession
}{m: make(map[string]*TerminalSession)}

func (s *Server) handleTerminalWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Pick shell
	shell := "/bin/zsh"
	if runtime.GOOS != "darwin" {
		if sh, err := exec.LookPath("bash"); err == nil {
			shell = sh
		}
	}

	cmd := exec.Command(shell, "-l")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color", "OPIUS_WEB=1")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("error: cannot start pty: "+err.Error()))
		return
	}
	defer ptmx.Close()

	id := r.URL.Query().Get("id")
	if id == "" {
		id = "default"
	}

	sess := &TerminalSession{id: id, cmd: cmd, ptmx: ptmx, ws: conn}

	sessions.Lock()
	sessions.m[id] = sess
	sessions.Unlock()

	defer func() {
		sessions.Lock()
		delete(sessions.m, id)
		sessions.Unlock()
	}()

	// pty -> ws (output from shell)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("pty read err: %v", err)
				}
				conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// ws -> pty (input from browser)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Handle resize messages: first byte 0x01 + 2 bytes cols + 2 bytes rows
		if len(msg) >= 5 && msg[0] == 0x01 {
			cols := int(msg[1])<<8 | int(msg[2])
			rows := int(msg[3])<<8 | int(msg[4])
			pty.Setsize(ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
			continue
		}

		sess.mu.Lock()
		_, err = ptmx.Write(msg)
		sess.mu.Unlock()
		if err != nil {
			return
		}
	}
}