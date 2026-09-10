package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type Store struct {
	RootPath string
}

func DefaultStorePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".opius", "packages"), nil
}

func NewStore(rootPath string) *Store {
	return &Store{RootPath: rootPath}
}

func (s *Store) Ensure() error {
	return os.MkdirAll(s.RootPath, 0o755)
}

func (s *Store) PackageDir(name, version string) string {
	return filepath.Join(s.RootPath, name, version)
}

func (s *Store) Install(m *Manifest, srcDir string) error {
	if err := s.Ensure(); err != nil {
		return err
	}

	destDir := s.PackageDir(m.Name, m.Version)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	// Copy each declared file from srcDir to destDir and verify hashes
	for _, f := range m.Files {
		src := filepath.Join(srcDir, f.Path)
		dst := filepath.Join(destDir, f.Path)

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("copy %s: %w", f.Path, err)
		}

		if f.SHA256 != "" {
			if err := VerifyFile(dst, f.SHA256); err != nil {
				return fmt.Errorf("verify %s: %w", f.Path, err)
			}
		}
	}

	// Write manifest inside package dir
	manifestPath := filepath.Join(destDir, "package.toml")
	if err := writeManifest(manifestPath, m); err != nil {
		return err
	}

	return nil
}

func (s *Store) List() ([]InstalledPackage, error) {
	if err := s.Ensure(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.RootPath)
	if err != nil {
		return nil, err
	}

	var result []InstalledPackage

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		nameDir := filepath.Join(s.RootPath, name)

		versions, err := os.ReadDir(nameDir)
		if err != nil {
			continue
		}

		for _, v := range versions {
			if !v.IsDir() {
				continue
			}
			version := v.Name()
			installDir := filepath.Join(nameDir, version)

			info, err := os.Stat(installDir)
			installedAt := time.Time{}
			if err == nil {
				installedAt = info.ModTime()
			}

			ip := InstalledPackage{
				Name:        name,
				Version:     version,
				InstallDir:  installDir,
				InstalledAt: installedAt,
			}

			if m, err := readManifest(filepath.Join(installDir, "package.toml")); err == nil {
				ip.Manifest = m
			}

			result = append(result, ip)
		}
	}

	return result, nil
}

func (s *Store) Remove(name, version string) error {
	dir := s.PackageDir(name, version)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("package %s@%s is not installed", name, version)
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	// Clean up empty parent dir
	parent := filepath.Dir(dir)
	entries, _ := os.ReadDir(parent)
	if len(entries) == 0 {
		os.Remove(parent)
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}

func writeManifest(path string, m *Manifest) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	return enc.Encode(m)
}

func readManifest(path string) (*Manifest, error) {
	var m Manifest
	if _, err := toml.DecodeFile(path, &m); err != nil {
		return nil, err
	}
	return &m, nil
}