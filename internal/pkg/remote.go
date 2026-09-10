package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RemoteRegistry handles downloading and syncing packages from a remote
// GitHub repository (or any URL serving a zip archive).
type RemoteRegistry struct {
	URL       string // e.g. https://github.com/am1s3/opius-os-pkg
	LocalPath string // local registry dir
}

func NewRemoteRegistry(url, localPath string) *RemoteRegistry {
	return &RemoteRegistry{URL: url, LocalPath: localPath}
}

// Sync downloads the latest archive from the remote registry and extracts
// it into the local registry directory. For GitHub repos, it uses the
// archive/refs/heads/main.zip endpoint.
func (r *RemoteRegistry) Sync(progress chan string) error {
	defer close(progress)

	if progress == nil {
		progress = make(chan string)
		go func() {
			for range progress {
			}
		}()
	}

	// Convert github URL to zip archive URL
	zipURL := toArchiveURL(r.URL)
	progress <- fmt.Sprintf("downloading from %s", zipURL)

	tmpFile, err := os.CreateTemp("", "opius-registry-*.zip")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	resp, err := http.Get(zipURL)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("write download: %w", err)
	}
	tmpFile.Close()

	progress <- "extracting archive"

	// Create temp extraction directory
	tmpDir, err := os.MkdirTemp("", "opius-registry-extract-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := extractZip(tmpFile.Name(), tmpDir); err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	// Find the actual registry content (GitHub nests it in a subdirectory)
	contentDir := findContentDir(tmpDir)
	if contentDir == "" {
		return fmt.Errorf("could not find registry content in archive")
	}

	// Ensure local registry exists
	if err := os.MkdirAll(r.LocalPath, 0o755); err != nil {
		return err
	}

	progress <- "updating local registry"

	// Sync files from extracted content to local registry
	if err := syncDir(contentDir, r.LocalPath); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	progress <- "sync complete"
	return nil
}

func toArchiveURL(url string) string {
	// Handle GitHub URLs
	if strings.Contains(url, "github.com") {
		url = strings.TrimSuffix(url, "/")
		// If already an archive URL, return as-is
		if strings.Contains(url, "archive/refs") {
			return url
		}
		// Convert to archive URL
		return url + "/archive/refs/heads/main.zip"
	}

	// Assume direct zip URL
	return url
}

func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(destDir, f.Name)

		// Protect against zip slip
		if !strings.HasPrefix(path, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func findContentDir(baseDir string) string {
	// GitHub archives have a single top-level directory
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return ""
	}

	// If there's exactly one directory, it's the GitHub root
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(baseDir, entries[0].Name())
	}

	// Otherwise, the content is directly in baseDir
	return baseDir
}

func syncDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}

		// Copy file
		return copyFileSync(path, dstPath)
	})
}

func copyFileSync(src, dst string) error {
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

	_, err = io.Copy(out, in)
	return err
}