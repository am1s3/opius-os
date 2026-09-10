package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Registry struct {
	LocalPath string
}

func NewRegistry(localPath string) *Registry {
	return &Registry{LocalPath: localPath}
}

func DefaultRegistryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".opius", "registry"), nil
}

func (r *Registry) Ensure() error {
	return os.MkdirAll(r.LocalPath, 0o755)
}

func (r *Registry) List() ([]Manifest, error) {
	if err := r.Ensure(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(r.LocalPath)
	if err != nil {
		return nil, err
	}

	manifests := make([]Manifest, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(e.Name()), ".toml") {
			continue
		}

		full := filepath.Join(r.LocalPath, e.Name())
		var m Manifest
		if _, err := toml.DecodeFile(full, &m); err != nil {
			continue // skip broken manifests
		}
		if m.Name == "" {
			continue
		}
		manifests = append(manifests, m)
	}

	return manifests, nil
}

func (r *Registry) Find(name string) (*Manifest, error) {
	all, err := r.List()
	if err != nil {
		return nil, err
	}

	for i := range all {
		if strings.EqualFold(all[i].Name, name) {
			return &all[i], nil
		}
	}

	return nil, fmt.Errorf("package %q not found in registry %s", name, r.LocalPath)
}

func (r *Registry) Search(query string) ([]Manifest, error) {
	all, err := r.List()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	if query == "" {
		return all, nil
	}

	result := make([]Manifest, 0, len(all))
	for _, m := range all {
		name := strings.ToLower(m.Name)
		desc := strings.ToLower(m.Description)
		if strings.Contains(name, query) || strings.Contains(desc, query) {
			result = append(result, m)
		}
	}

	return result, nil
}

func (r *Registry) Add(m *Manifest) error {
	if err := r.Ensure(); err != nil {
		return err
	}

	filename := filepath.Join(r.LocalPath, m.Name+"-"+m.Version+".toml")

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	return enc.Encode(m)
}