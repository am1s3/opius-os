package flash

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ESPDriver struct {
	ToolPath string
}

func NewESPDriver(toolPath string) *ESPDriver {
	return &ESPDriver{ToolPath: toolPath}
}

func (d *ESPDriver) Name() string { return "esp" }

func (d *ESPDriver) Supports(t Target) bool {
	switch strings.ToLower(t.Chip) {
	case "esp32", "esp32-s3", "esp32-c3", "esp32-c2", "esp32-s2", "esp8266":
		return true
	}
	return false
}

func (d *ESPDriver) Plan(ctx context.Context, t Target) (*Plan, error) {
	if d.ToolPath == "" {
		return nil, errors.New("esptool not found in PATH")
	}

	if t.Port == "" {
		return nil, errors.New("serial port is required")
	}

	baud := t.Baud
	if baud <= 0 {
		baud = 460800
	}

	steps := []Step{}

	// Erase-only plan
	if (t.File == "" && len(t.Binaries) == 0) && t.Erase {
		steps = append(steps, Step{
			Name:        "erase-flash",
			Description: "Full flash erase",
			Phase:       "erase",
			Command: []string{
				d.ToolPath,
				"--port", t.Port,
				"--baud", fmt.Sprintf("%d", baud),
				"erase_flash",
			},
		})
		return &Plan{
			Target:  t,
			Driver:  d.Name(),
			Tool:    d.ToolPath,
			Steps:   steps,
			Summary: "erase-only plan",
		}, nil
	}

	// Single file flash (backward compatibility)
	if t.File != "" && len(t.Binaries) == 0 {
		if _, err := os.Stat(t.File); err != nil {
			return nil, fmt.Errorf("firmware file not found: %s", t.File)
		}

		writeAddr := "0x10000"
		if strings.ToLower(t.Chip) == "esp8266" {
			writeAddr = "0x0"
		}

		steps = append(steps,
			Step{
				Name:        "chip-info",
				Description: "Query chip identification",
				Phase:       "prepare",
				Command:     []string{d.ToolPath, "--port", t.Port, "chip_id"},
			},
		)

		if t.Erase {
			steps = append(steps, Step{
				Name:        "erase-flash",
				Description: "Full flash erase",
				Phase:       "erase",
				Command: []string{
					d.ToolPath,
					"--port", t.Port,
					"--baud", fmt.Sprintf("%d", baud),
					"erase_flash",
				},
			})
		}

		steps = append(steps, Step{
			Name:        "write-firmware",
			Description: fmt.Sprintf("Write %s to flash at %s", filepath.Base(t.File), writeAddr),
			Phase:       "write",
			Command: []string{
				d.ToolPath,
				"--port", t.Port,
				"--baud", fmt.Sprintf("%d", baud),
				"--chip", strings.ToLower(t.Chip),
				"write_flash",
				writeAddr,
				t.File,
			},
		})

		steps = append(steps, Step{
			Name:        "verify-hash",
			Description: "Verify written image hash",
			Phase:       "verify",
			Command: []string{
				d.ToolPath,
				"--port", t.Port,
				"--chip", strings.ToLower(t.Chip),
				"flash_md5sum",
				writeAddr,
				t.File,
			},
		})

		steps = append(steps, Step{
			Name:        "reboot",
			Description: "Reset chip into application mode",
			Phase:       "finalize",
			Command: []string{
				d.ToolPath,
				"--port", t.Port,
				"--chip", strings.ToLower(t.Chip),
				"flash_finish",
			},
		})

		return &Plan{
			Target:  t,
			Driver:  d.Name(),
			Tool:    d.ToolPath,
			Steps:   steps,
			Summary: fmt.Sprintf("%d steps", len(steps)),
		}, nil
	}

	// Multi-binary flash
	if len(t.Binaries) == 0 {
		return nil, errors.New("no firmware file or binaries specified")
	}

	// Verify all binary files exist
	for _, bin := range t.Binaries {
		if _, err := os.Stat(bin.Path); err != nil {
			return nil, fmt.Errorf("binary file not found: %s", bin.Path)
		}
	}

	steps = append(steps, Step{
		Name:        "chip-info",
		Description: "Query chip identification",
		Phase:       "prepare",
		Command:     []string{d.ToolPath, "--port", t.Port, "chip_id"},
	})

	if t.Erase {
		steps = append(steps, Step{
			Name:        "erase-flash",
			Description: "Full flash erase",
			Phase:       "erase",
			Command: []string{
				d.ToolPath,
				"--port", t.Port,
				"--baud", fmt.Sprintf("%d", baud),
				"erase_flash",
			},
		})
	}

	// Build write_flash command with all binaries
	writeCmd := []string{
		d.ToolPath,
		"--port", t.Port,
		"--baud", fmt.Sprintf("%d", baud),
		"--chip", strings.ToLower(t.Chip),
		"write_flash",
	}

	desc := "Write multiple binaries:"
	for _, bin := range t.Binaries {
		writeCmd = append(writeCmd, bin.Address, bin.Path)
		desc += fmt.Sprintf("\n  %s @ %s", filepath.Base(bin.Path), bin.Address)
	}

	steps = append(steps, Step{
		Name:        "write-binaries",
		Description: desc,
		Phase:       "write",
		Command:     writeCmd,
	})

	steps = append(steps, Step{
		Name:        "reboot",
		Description: "Reset chip into application mode",
		Phase:       "finalize",
		Command: []string{
			d.ToolPath,
			"--port", t.Port,
			"--chip", strings.ToLower(t.Chip),
			"flash_finish",
		},
	})

	return &Plan{
		Target:  t,
		Driver:  d.Name(),
		Tool:    d.ToolPath,
		Steps:   steps,
		Summary: fmt.Sprintf("%d binaries, %d steps", len(t.Binaries), len(steps)),
	}, nil
}

func (d *ESPDriver) Execute(ctx context.Context, p *Plan, dryRun bool, ev chan<- ProgressEvent) error {
	defer close(ev)

	if dryRun {
		for i, s := range p.Steps {
			ev <- ProgressEvent{
				Step:    s.Name,
				Phase:   s.Phase,
				Message: "[dry-run] would execute: " + strings.Join(s.Command, " "),
				Percent: int((i + 1) * 100 / len(p.Steps)),
			}
		}
		return nil
	}

	for i, s := range p.Steps {
		ev <- ProgressEvent{
			Step:    s.Name,
			Phase:   s.Phase,
			Message: "running: " + s.Description,
			Percent: int(i * 100 / len(p.Steps)),
		}

		if len(s.Command) == 0 {
			continue
		}

		cmd := exec.CommandContext(ctx, s.Command[0], s.Command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("step %s failed: %w", s.Name, err)
		}

		ev <- ProgressEvent{
			Step:    s.Name,
			Phase:   s.Phase,
			Message: s.Description + " done",
			Percent: int((i + 1) * 100 / len(p.Steps)),
		}
	}

	return nil
}