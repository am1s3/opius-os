package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/pkg"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/progress"
	"github.com/opius-os/opius/internal/ui/style"
)

var pkgDownloadVersion string

var pkgDownloadCmd = &cobra.Command{
	Use:   "download <name>",
	Short: "Download a firmware .bin file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		reg := newPkgRegistry()
		if err := reg.Sync(); err != nil {
			return fmt.Errorf("failed to sync registry: %w", err)
		}

		m, err := reg.Find(name)
		if err != nil {
			return err
		}

		if len(m.Versions) == 0 {
			return fmt.Errorf("no versions available for %s", name)
		}

		// Select version
		var version *pkg.PackageVersion
		if pkgDownloadVersion != "" {
			for i := range m.Versions {
				if m.Versions[i].Version == pkgDownloadVersion {
					version = &m.Versions[i]
					break
				}
			}
			if version == nil {
				return fmt.Errorf("version %s not found. Available: %s", pkgDownloadVersion, listVersionNames(m.Versions))
			}
		} else {
			version = &m.Versions[0] // Latest
		}

		components.PrintSection("Downloading")
		fmt.Printf("  %-12s %s\n", "Package", style.Value.Render(m.Name))
		fmt.Printf("  %-12s %s\n", "Version", style.Value.Render(version.Version))
		fmt.Printf("  %-12s %s\n", "URL", style.Dim.Render(version.DownloadURL))
		fmt.Println()

		// Create downloads directory
		home, _ := os.UserHomeDir()
		downloadDir := filepath.Join(home, ".opius", "downloads")
		os.MkdirAll(downloadDir, 0755)

		filename := filepath.Base(version.DownloadURL)
		if filename == "" || filename == "/" {
			filename = m.Name + "-" + version.Version + ".bin"
		}
		destPath := filepath.Join(downloadDir, filename)

		// Download with progress
		if err := downloadWithProgress(version.DownloadURL, destPath, version.Size); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		fmt.Println()

		// Verify SHA256
		if version.SHA256 != "" {
			fmt.Print("  Verifying SHA256... ")
			if err := verifyFileSHA256(destPath, version.SHA256); err != nil {
				fmt.Println(components.BadgeErr())
				return fmt.Errorf("checksum mismatch: %w", err)
			}
			fmt.Println(components.BadgeOK())
		}

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " Downloaded to " + style.Value.Render(destPath))
		fmt.Println()
		fmt.Println(style.Dim.Render("  Flash with: opius pkg flash " + m.Name))

		return nil
	},
}

func listVersionNames(versions []pkg.PackageVersion) string {
	names := make([]string, len(versions))
	for i, v := range versions {
		names[i] = v.Version
	}
	return strings.Join(names, ", ")
}

func downloadWithProgress(url, destPath string, expectedSize int64) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	totalSize := expectedSize
	if totalSize == 0 && resp.ContentLength > 0 {
		totalSize = resp.ContentLength
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	bar := progress.NewBar(totalSize, "Downloading")

	buf := make([]byte, 32*1024)
	var downloaded int64

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
			downloaded += int64(n)
			bar.Update(downloaded)
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return readErr
		}
	}

	bar.Finish()
	return nil
}

func verifyFileSHA256(path, expected string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	actual := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actual, expected) {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}

	return nil
}

func init() {
	pkgDownloadCmd.Flags().StringVarP(&pkgDownloadVersion, "version", "v", "", "Specific version to download (default: latest)")
}