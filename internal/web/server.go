package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"
)

//go:embed assets/*
var assetsFS embed.FS

type Server struct {
	Port    int
	Host    string
	server  *http.Server
	handler http.Handler
}

func NewServer(host string, port int) *Server {
	s := &Server{
		Port: port,
		Host: host,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/devices", s.handleDevices)
	mux.HandleFunc("/api/packages", s.handlePackages)
	mux.HandleFunc("/api/flash", s.handleFlash)
	mux.HandleFunc("/api/flash-custom", s.handleFlashCustom)
	mux.HandleFunc("/api/flash/progress", s.handleFlashProgress)
	mux.HandleFunc("/api/install", s.handleInstall)
	mux.HandleFunc("/api/erase", s.handleErase)
	mux.HandleFunc("/api/reset", s.handleReset)
	mux.HandleFunc("/api/sync", s.handleSync)
	mux.HandleFunc("/api/remove", s.handleRemove)
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/system", s.handleSystem)
	mux.HandleFunc("/api/fs/list", s.handleFSList)
	mux.HandleFunc("/api/script/compile", s.handleScriptCompile)
	mux.HandleFunc("/api/script/flash", s.handleScriptFlash)

	mux.HandleFunc("/ws/terminal", s.handleTerminalWS)
	mux.HandleFunc("/ws/serial", s.handleSerialMonitorWS)

	assetsSubFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		panic(fmt.Sprintf("failed to create sub filesystem: %v", err))
	}
	mux.Handle("/", http.FileServer(http.FS(assetsSubFS)))

	s.handler = mux
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.handler,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	fmt.Printf("  Web UI running at http://%s\n", addr)
	fmt.Printf("  Press Ctrl+C to stop\n\n")

	return s.server.Serve(ln)
}

func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}