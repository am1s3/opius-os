package device

import (
	"fmt"
	"runtime"
	"strings"

	"go.bug.st/serial/enumerator"
)

type Device struct {
	ID           int    `json:"id"`
	Port         string `json:"port"`
	IsUSB        bool   `json:"is_usb"`
	VID          string `json:"vid,omitempty"`
	PID          string `json:"pid,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	Product      string `json:"product,omitempty"`
	Chip         string `json:"chip"`
	Board        string `json:"board"`
	Mode         string `json:"mode"`
	Status       string `json:"status"`
	Hint         string `json:"hint,omitempty"`
}

func List() ([]Device, error) {
	return ListWithFilter(true) // by default only USB
}

func ListAll() ([]Device, error) {
	return ListWithFilter(false) // show everything including Bluetooth
}

func ListWithFilter(usbOnly bool) ([]Device, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}

	devices := make([]Device, 0, len(ports))
	idx := 1

	for _, p := range ports {
		// Filter out non-USB ports (Bluetooth, virtual, etc.)
		if usbOnly && !p.IsUSB {
			continue
		}

		d := Device{
			ID:           idx,
			Port:         p.Name,
			IsUSB:        p.IsUSB,
			VID:          normalizeHex(p.VID),
			PID:          normalizeHex(p.PID),
			SerialNumber: p.SerialNumber,
			Product:      cleanText(p.Product),
			Chip:         "unknown",
			Board:        "unknown",
			Mode:         "normal",
			Status:       "unknown",
		}

		guessDevice(&d)
		devices = append(devices, d)
		idx++
	}

	return devices, nil
}

func FindByPort(port string) (*Device, error) {
	devices, err := ListAll()
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.Port == port {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("device on port %s not found", port)
}

func PreferredPortNote() string {
	if runtime.GOOS == "darwin" {
		return "On macOS prefer /dev/cu.* ports for serial communication."
	}
	return ""
}

func normalizeHex(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	v = strings.TrimPrefix(v, "0x")
	if v == "" {
		return ""
	}
	return v
}

func cleanText(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	return v
}

func guessDevice(d *Device) {
	vid := strings.ToLower(d.VID)
	pid := strings.ToLower(d.PID)
	product := strings.ToLower(d.Product)
	port := strings.ToLower(d.Port)

	switch vid {
	case "303a":
		d.Chip = "esp32"
		d.Board = "Espressif USB/JTAG device"
		d.Status = "ready"
		d.Hint = "Native Espressif USB device detected."

	case "10c4":
		d.Chip = "esp/serial"
		d.Board = "CP210x USB-UART board"
		d.Status = "ready"
		d.Hint = "Common on ESP32/ESP8266 development boards."

	case "1a86":
		d.Chip = "serial"
		d.Board = "CH340/CH341 USB-UART board"
		d.Status = "ready"
		d.Hint = "Common on Arduino Nano clones and ESP boards."

	case "0403":
		d.Chip = "serial"
		d.Board = "FTDI USB-UART board"
		d.Status = "ready"
		d.Hint = "Generic FTDI serial adapter."

	case "2341", "2a03":
		d.Chip = "avr/samd"
		d.Board = "Arduino board"
		d.Status = "ready"
		d.Hint = "Official Arduino VID detected."

	case "2e8a":
		d.Chip = "rp2040"
		d.Board = "Raspberry Pi RP2040"
		d.Status = "ready"
		if strings.Contains(product, "boot") || pid == "0003" {
			d.Mode = "bootsel"
			d.Status = "boot"
			d.Hint = "RP2040 BOOTSEL mode detected."
		}
	}

	if strings.Contains(product, "rp2") || strings.Contains(product, "pico") {
		d.Chip = "rp2040"
		d.Board = "Raspberry Pi Pico / RP2040"
		d.Status = "ready"
	}

	if strings.Contains(product, "esp") {
		d.Chip = "esp32"
		d.Board = "Espressif device"
		d.Status = "ready"
	}

	if strings.Contains(product, "arduino") {
		d.Chip = "arduino"
		d.Board = "Arduino board"
		d.Status = "ready"
	}

	if strings.Contains(port, "/dev/cu.") && runtime.GOOS == "darwin" && d.Status == "unknown" {
		d.Status = "ready"
		d.Hint = "macOS serial callout port detected."
	}

	if d.Status == "unknown" && d.IsUSB {
		d.Status = "ready"
		d.Hint = "USB serial device detected, but board is unknown."
	}
}