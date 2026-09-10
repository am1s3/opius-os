package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SHA256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func VerifyFile(path, expectedSHA256 string) error {
	actual, err := SHA256File(path)
	if err != nil {
		return err
	}
	if actual != expectedSHA256 {
		return fmt.Errorf("sha256 mismatch: expected %s, got %s", expectedSHA256, actual)
	}
	return nil
}

func DownloadAndVerify(url, destPath, expectedSHA256 string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
		return err
	}

	actual := hex.EncodeToString(h.Sum(nil))
	if actual != expectedSHA256 {
		return fmt.Errorf("sha256 mismatch: expected %s, got %s", expectedSHA256, actual)
	}

	return nil
}