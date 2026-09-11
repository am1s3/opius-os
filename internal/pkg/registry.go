package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PackageIndex struct {
	Version     string         `json:"version"`
	LastUpdated string         `json:"last_updated"`
	Packages    []PackageEntry `json:"packages"`
}

type PackageEntry struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Chips       []string `json:"chips"`
	ManifestURL string   `json:"manifest_url"`
}

type PackageManifest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Chips       []string        `json:"chips"`
	Boards      []string        `json:"boards"`
	Homepage    string          `json:"homepage"`
	License     string          `json:"license"`
	Versions    []PackageVersion `json:"versions"`
}

type PackageVersion struct {
	Version     string `json:"version"`
	ReleaseDate string `json:"release_date"`
	DownloadURL string `json:"download_url"`
	SHA256      string `json:"sha256"`
	Size        int64  `json:"size"`
}

type Registry struct {
	IndexURL string
	CacheDir string
	index    *PackageIndex
}

func NewRegistry(indexURL, cacheDir string) *Registry {
	return &Registry{
		IndexURL: indexURL,
		CacheDir: cacheDir,
	}
}

func DefaultRegistryURL() string {
	return "https://raw.githubusercontent.com/am1s3/opius-os-pkg/main/index.json"
}

func DefaultCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".opius", "cache"), nil
}

func (r *Registry) Sync() error {
	resp, err := http.Get(r.IndexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registry returned HTTP %d", resp.StatusCode)
	}

	var index PackageIndex
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return fmt.Errorf("failed to parse registry: %w", err)
	}

	r.index = &index

	if err := r.saveCache(&index); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not cache registry: %v\n", err)
	}

	return nil
}

func (r *Registry) Load() (*PackageIndex, error) {
	if r.index != nil {
		return r.index, nil
	}

	cached, err := r.loadCache()
	if err == nil && cached != nil {
		r.index = cached
		return cached, nil
	}

	if err := r.Sync(); err != nil {
		return nil, err
	}

	return r.index, nil
}

func (r *Registry) List() ([]PackageEntry, error) {
	index, err := r.Load()
	if err != nil {
		return nil, err
	}
	return index.Packages, nil
}

func (r *Registry) Find(name string) (*PackageManifest, error) {
	index, err := r.Load()
	if err != nil {
		return nil, err
	}

	var manifestURL string
	for _, pkg := range index.Packages {
		if strings.EqualFold(pkg.Name, name) {
			manifestURL = pkg.ManifestURL
			break
		}
	}

	if manifestURL == "" {
		return nil, fmt.Errorf("package %q not found", name)
	}

	resp, err := http.Get(manifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest returned HTTP %d", resp.StatusCode)
	}

	var manifest PackageManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

func (r *Registry) Search(query string) ([]PackageEntry, error) {
	packages, err := r.List()
	if err != nil {
		return nil, err
	}

	if query == "" {
		return packages, nil
	}

	query = strings.ToLower(query)
	var result []PackageEntry
	for _, p := range packages {
		name := strings.ToLower(p.Name)
		desc := strings.ToLower(p.Description)
		if strings.Contains(name, query) || strings.Contains(desc, query) {
			result = append(result, p)
		}
	}

	return result, nil
}

func (r *Registry) saveCache(index *PackageIndex) error {
	if err := os.MkdirAll(r.CacheDir, 0755); err != nil {
		return err
	}

	cachePath := filepath.Join(r.CacheDir, "registry.json")
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func (r *Registry) loadCache() (*PackageIndex, error) {
	cachePath := filepath.Join(r.CacheDir, "registry.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var index PackageIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	return &index, nil
}