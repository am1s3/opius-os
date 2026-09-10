package hostpkg

import (
	"os/exec"
)

type Backend struct {
	Name string
	Path string
}

func Detect() Backend {
	if p, err := exec.LookPath("port"); err == nil {
		return Backend{Name: "macports", Path: p}
	}
	if p, err := exec.LookPath("brew"); err == nil {
		return Backend{Name: "homebrew", Path: p}
	}
	return Backend{Name: "none", Path: ""}
}

func DetectAll() []Backend {
	var backends []Backend
	
	if p, err := exec.LookPath("port"); err == nil {
		backends = append(backends, Backend{Name: "macports", Path: p})
	}
	if p, err := exec.LookPath("brew"); err == nil {
		backends = append(backends, Backend{Name: "homebrew", Path: p})
	}
	
	return backends
}