package pkg

import (
	"time"
)

type Manifest struct {
	Name         string       `toml:"name"`
	Version      string       `toml:"version"`
	Description  string       `toml:"description"`
	License      string       `toml:"license"`
	Homepage     string       `toml:"homepage"`
	Repository   string       `toml:"repository"`
	Authors      []string     `toml:"authors"`
	Chips        []string     `toml:"chips"`
	Boards       []string     `toml:"boards"`
	Files        []File       `toml:"files"`
	Dependencies []Dependency `toml:"dependencies"`
	CreatedAt    time.Time    `toml:"created_at"`
}

type File struct {
	Path   string `toml:"path"`
	SHA256 string `toml:"sha256"`
	Size   int64  `toml:"size"`
	Type   string `toml:"type"` // firmware, config, tool, doc, asset
}

type Dependency struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

// InstalledPackage describes a package already present in the local store.
type InstalledPackage struct {
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	InstallDir string    `json:"install_dir"`
	InstalledAt time.Time `json:"installed_at"`
	Manifest   *Manifest `json:"manifest,omitempty"`
}