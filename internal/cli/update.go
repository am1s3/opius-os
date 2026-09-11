package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opius-os/opius/internal/build"
	"github.com/opius-os/opius/internal/ui/components"
	"github.com/opius-os/opius/internal/ui/style"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Opius OS to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		components.PrintSection("Checking for updates")

		currentVersion := build.Version
		fmt.Printf("  %-12s %s\n", "Current", style.Value.Render(currentVersion))
		fmt.Print("  Checking latest release... ")

		latest, err := getLatestRelease()
		if err != nil {
			fmt.Println(components.BadgeErr())
			return fmt.Errorf("failed to check updates: %w", err)
		}

		fmt.Println(components.BadgeOK())
		fmt.Printf("  %-12s %s\n", "Latest", style.Value.Render(latest.TagName))

		if strings.TrimPrefix(latest.TagName, "v") == strings.TrimPrefix(currentVersion, "v") {
			fmt.Println()
			fmt.Println("  " + components.BadgeOK() + " " + style.Success.Render("You are already on the latest version!"))
			return nil
		}

		fmt.Println()
		fmt.Println(style.Info.Render("  New version available: " + latest.TagName))
		fmt.Println(style.Dim.Render("  " + latest.Body))
		fmt.Println()

		// Find download URL for current platform
		assetName := fmt.Sprintf("opius_%s_%s", runtime.GOOS, runtime.GOARCH)
		var downloadURL string

		for _, asset := range latest.Assets {
			if strings.Contains(asset.Name, assetName) {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}

		if downloadURL == "" {
			fmt.Println("  " + components.BadgeWarn() + " No prebuilt binary for your platform")
			fmt.Println(style.Dim.Render("  Please build from source:"))
			fmt.Println(style.Dim.Render("    git clone https://github.com/am1s3/opius-os.git"))
			fmt.Println(style.Dim.Render("    cd opius-os && go build -o opius ."))
			return nil
		}

		if !components.PromptConfirm("  Download and install "+latest.TagName+"?", true) {
			return nil
		}

		fmt.Println()
		components.PrintSection("Updating")

		if err := downloadAndInstall(downloadURL); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		fmt.Println()
		fmt.Println("  " + components.BadgeOK() + " Updated to " + style.Success.Render(latest.TagName))
		fmt.Println(style.Dim.Render("  Run 'opius version' to verify"))

		return nil
	},
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func getLatestRelease() (*GitHubRelease, error) {
	url := "https://api.github.com/repos/am1s3/opius-os/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned HTTP %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadAndInstall(url string) error {
	fmt.Print("  Downloading... ")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// Save to temp file
	tmpFile := "/tmp/opius-update"
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	if _, err := f.ReadFrom(resp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// Make executable
	os.Chmod(tmpFile, 0755)

	fmt.Println(components.BadgeOK())
	fmt.Print("  Installing... ")

	// Try to replace current binary
	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine current binary path: %w", err)
	}

	// Replace binary
	if err := os.Rename(tmpFile, currentBinary); err != nil {
		// Try copying instead
		if err := copyFile(tmpFile, currentBinary); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
		os.Remove(tmpFile)
	}

	fmt.Println(components.BadgeOK())
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}