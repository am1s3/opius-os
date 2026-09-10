package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/opius-os/opius/internal/config"
	"github.com/opius-os/opius/internal/device"
	"github.com/opius-os/opius/internal/flash"
	"github.com/opius-os/opius/internal/pkg"
	"go.bug.st/serial"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	devs, err := device.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: devs})
}

func (s *Server) handlePackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	regPath, err := pkg.DefaultRegistryPath()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	registry := pkg.NewRegistry(regPath)
	packages, err := registry.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: packages})
}

func (s *Server) handleInstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	regPath, err := pkg.DefaultRegistryPath()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	registry := pkg.NewRegistry(regPath)
	m, err := registry.Find(req.Name)
	if err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: err.Error()})
		return
	}

	storePath, err := pkg.DefaultStorePath()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	store := pkg.NewStore(storePath)
	if err := store.Install(m, regPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": fmt.Sprintf("package %s@%s installed successfully", m.Name, m.Version)},
	})
}

func (s *Server) handleFlash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		PackageName string `json:"package_name"`
		Port        string `json:"port"`
		Chip        string `json:"chip"`
		Erase       bool   `json:"erase"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	regPath, err := pkg.DefaultRegistryPath()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	registry := pkg.NewRegistry(regPath)
	m, err := registry.Find(req.PackageName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: err.Error()})
		return
	}

	var firmwareFile *pkg.File
	for i := range m.Files {
		if m.Files[i].Type == "firmware" {
			firmwareFile = &m.Files[i]
			break
		}
	}

	if firmwareFile == nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "no firmware file found"})
		return
	}

	storePath, err := pkg.DefaultStorePath()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	store := pkg.NewStore(storePath)
	installDir := store.PackageDir(m.Name, m.Version)
	firmwarePath := filepath.Join(installDir, firmwareFile.Path)

	chip := req.Chip
	if chip == "" && len(m.Chips) > 0 {
		chip = m.Chips[0]
	}
	if chip == "" || chip == "serial" || chip == "unknown" {
		chip = "esp32"
	}

	target := flash.Target{
		Chip:  chip,
		Port:  req.Port,
		File:  firmwarePath,
		Erase: req.Erase,
	}

	flashRegistry := flash.NewRegistry()
	driver, err := flashRegistry.Select(target)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	plan, err := driver.Plan(context.Background(), target)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	evCh := make(chan flash.ProgressEvent)
	errCh := make(chan error, 1)

	go func() {
		errCh <- driver.Execute(context.Background(), plan, false, evCh)
	}()

	for range evCh {
	}

	if err := <-errCh; err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": fmt.Sprintf("package %s flashed to %s", m.Name, req.Port)},
	})
}

func (s *Server) handleFlashCustom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		Port     string              `json:"port"`
		Chip     string              `json:"chip"`
		Binaries []flash.BinaryEntry `json:"binaries"`
		Erase    bool                `json:"erase"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	chip := req.Chip
	if chip == "" || chip == "serial" || chip == "unknown" {
		chip = "esp32"
	}

	target := flash.Target{
		Chip:     chip,
		Port:     req.Port,
		Binaries: req.Binaries,
		Erase:    req.Erase,
	}

	registry := flash.NewRegistry()
	driver, err := registry.Select(target)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	plan, err := driver.Plan(context.Background(), target)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	evCh := make(chan flash.ProgressEvent)
	errCh := make(chan error, 1)

	go func() {
		errCh <- driver.Execute(context.Background(), plan, false, evCh)
	}()

	for range evCh {
	}

	if err := <-errCh; err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": fmt.Sprintf("custom flash completed: %d binaries", len(req.Binaries))},
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"status": "ok", "version": "1.0.0"},
	})
}

func (s *Server) handleSystem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	home, _ := os.UserHomeDir()
	opiusDir := filepath.Join(home, ".opius")

	var totalSize int64
	filepath.Walk(opiusDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"go_version": runtime.Version(),
			"cpus":       runtime.NumCPU(),
			"goroutines": runtime.NumGoroutine(),
			"opius_dir":  opiusDir,
			"opius_size": totalSize,
			"time":       time.Now().Format(time.RFC3339),
			"hostname":   hostname(),
		},
	})
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

func (s *Server) handleFSList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".opius")
	}

	home, _ := os.UserHomeDir()
	opiusRoot := filepath.Join(home, ".opius")
	absPath, _ := filepath.Abs(path)
	if !isSubPath(opiusRoot, absPath) {
		writeJSON(w, http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "access denied: only ~/.opius directory is accessible",
		})
		return
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	type Entry struct {
		Name  string `json:"name"`
		IsDir bool   `json:"is_dir"`
		Size  int64  `json:"size"`
		Mode  string `json:"mode"`
	}

	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		info, _ := e.Info()
		size := int64(0)
		mode := ""
		if info != nil {
			size = info.Size()
			mode = info.Mode().String()
		}
		result = append(result, Entry{
			Name:  e.Name(),
			IsDir: e.IsDir(),
			Size:  size,
			Mode:  mode,
		})
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"path":    absPath,
			"entries": result,
		},
	})
}

func isSubPath(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	return rel != ".." && !filepath.IsAbs(rel) && rel[:2] != ".."
}

func (s *Server) handleErase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		Port string `json:"port"`
		Chip string `json:"chip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	if req.Port == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "port is required"})
		return
	}

	chip := req.Chip
	if chip == "" || chip == "serial" || chip == "unknown" {
		chip = "esp32"
	}

	registry := flash.NewRegistry()
	driver, err := registry.Select(flash.Target{Chip: chip})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	target := flash.Target{
		Chip:  chip,
		Port:  req.Port,
		Erase: true,
	}

	plan, err := driver.Plan(context.Background(), target)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	evCh := make(chan flash.ProgressEvent)
	errCh := make(chan error, 1)

	go func() {
		errCh <- driver.Execute(context.Background(), plan, false, evCh)
	}()

	for range evCh {
	}

	if err := <-errCh; err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "device flash erased successfully"},
	})
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		Port string `json:"port"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	if req.Port == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "port is required"})
		return
	}

	serPort, err := serial.Open(req.Port, &serial.Mode{BaudRate: 115200})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "cannot open port: " + err.Error()})
		return
	}
	defer serPort.Close()

	serPort.SetDTR(false)
	serPort.SetRTS(true)
	time.Sleep(100 * time.Millisecond)
	serPort.SetRTS(false)
	time.Sleep(100 * time.Millisecond)

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "reset signal sent to " + req.Port},
	})
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	cfg, err := config.Load()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	regDir, err := cfg.RegistryDir()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	remote := pkg.NewRemoteRegistry(cfg.RegistryURL, regDir)
	progress := make(chan string, 10)

	go func() {
		remote.Sync(progress)
	}()

	for range progress {
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "registry synced successfully"},
	})
}

func (s *Server) handleRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	storePath, err := pkg.DefaultStorePath()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	store := pkg.NewStore(storePath)
	if err := store.Remove(req.Name, req.Version); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "package removed successfully"},
	})
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "failed to parse form"})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "file is required"})
		return
	}
	defer file.Close()

	home, _ := os.UserHomeDir()
	uploadDir := filepath.Join(home, ".opius", "uploads")
	os.MkdirAll(uploadDir, 0755)

	destPath := filepath.Join(uploadDir, handler.Filename)
	dest, err := os.Create(destPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"path": destPath, "filename": handler.Filename},
	})
}