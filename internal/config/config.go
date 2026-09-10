package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	DefaultRegistryURL = "https://github.com/am1s3/opius-os-pkg"
	DefaultRepository  = "https://github.com/am1s3/opius-os"
)

type Config struct {
	HomeDir        string
	ConfigPath     string
	PackageBackend string
	Workspace      string
	RegistryURL    string
	Repository     string
}

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "opius"), nil
}

func (c *Config) ConfigFile() string {
	return filepath.Join(c.HomeDir, "opius.toml")
}

func (c *Config) WorkspaceDir() string {
	return c.Workspace
}

func (c *Config) RegistryDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".opius", "registry"), nil
}

func (c *Config) StoreDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".opius", "packages"), nil
}

func Load() (*Config, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(dir, "opius.toml")

	cfg := &Config{
		HomeDir:        dir,
		ConfigPath:     configPath,
		PackageBackend: "auto",
		RegistryURL:    DefaultRegistryURL,
		Repository:     DefaultRepository,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfg.Workspace = filepath.Join(home, ".opius", "workspace")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	var raw struct {
		PackageBackend string `toml:"package_backend"`
		Workspace      string `toml:"workspace"`
		RegistryURL    string `toml:"registry_url"`
		Repository     string `toml:"repository"`
	}

	if _, err := toml.DecodeFile(configPath, &raw); err != nil {
		return nil, err
	}

	if raw.PackageBackend != "" {
		cfg.PackageBackend = raw.PackageBackend
	}
	if raw.Workspace != "" {
		cfg.Workspace = raw.Workspace
	}
	if raw.RegistryURL != "" {
		cfg.RegistryURL = raw.RegistryURL
	}
	if raw.Repository != "" {
		cfg.Repository = raw.Repository
	}

	return cfg, nil
}

func (c *Config) Save() error {
	if err := os.MkdirAll(c.HomeDir, 0o755); err != nil {
		return err
	}

	content := "# Opius OS configuration\n"
	content += "package_backend = \"" + c.PackageBackend + "\"\n"
	content += "workspace = \"" + c.Workspace + "\"\n"
	content += "registry_url = \"" + c.RegistryURL + "\"\n"
	content += "repository = \"" + c.Repository + "\"\n"

	return os.WriteFile(c.ConfigFile(), []byte(content), 0o644)
}

func Init() (*Config, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	for _, sub := range []string{"tools", "packages", "cache", "backups", "logs", "workspace", "registry"} {
		if err := os.MkdirAll(filepath.Join(home, ".opius", sub), 0o755); err != nil {
			return nil, err
		}
	}

	cfg := &Config{
		HomeDir:        dir,
		PackageBackend: "auto",
		Workspace:      filepath.Join(home, ".opius", "workspace"),
		RegistryURL:    DefaultRegistryURL,
		Repository:     DefaultRepository,
	}

	if err := cfg.Save(); err != nil {
		return nil, err
	}

	return cfg, nil
}