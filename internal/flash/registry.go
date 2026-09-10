package flash

import (
	"fmt"
	"os/exec"
)

// Registry holds all available drivers and resolves tool paths.
type Registry struct {
	drivers []Driver
}

// NewRegistry builds a registry with drivers supported on the current host.
func NewRegistry() *Registry {
	r := &Registry{}

	if p, err := exec.LookPath("esptool"); err == nil {
		r.drivers = append(r.drivers, NewESPDriver(p))
	} else if p, err := exec.LookPath("esptool.py"); err == nil {
		r.drivers = append(r.drivers, NewESPDriver(p))
	}

	// Future: r.drivers = append(r.drivers, NewRP2040Driver())
	// Future: r.drivers = append(r.drivers, NewAVRDriver())

	return r
}

// Select returns the first driver that supports the given target.
func (r *Registry) Select(t Target) (Driver, error) {
	for _, d := range r.drivers {
		if d.Supports(t) {
			return d, nil
		}
	}
	return nil, fmt.Errorf("no driver supports chip=%s; install the required toolchain (esptool for ESP, avrdude for AVR, picotool for RP2040)", t.Chip)
}

// HasDriver reports whether at least one driver is available.
func (r *Registry) HasDriver() bool {
	return len(r.drivers) > 0
}

// Names returns the names of all registered drivers.
func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.drivers))
	for _, d := range r.drivers {
		out = append(out, d.Name())
	}
	return out
}