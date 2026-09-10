# Opius OS

Terminal-first embedded workspace for microcontrollers.

Opius OS — это кроссплатформенный инструмент для работы с микроконтроллерами: прошивка, мониторинг, управление пакетами прошивок, и полноценный веб-интерфейс в стиле десктопной ОС.

![Opius OS](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Features

- 🔌 **Device Management** — автоматическое определение подключённых устройств (ESP32, ESP8266, Arduino, RP2040, STM32)
- ⚡ **Flash Firmware** — прошивка из пакетов или кастомных `.bin` файлов с указанием адресов памяти
- 📦 **Package Manager** — реестр прошивок с проверкой контрольных сумм
- 🖥️ **Serial Monitor** — цветной монитор с отправкой команд в устройство
- 🌐 **Web UI** — полноценный десктоп в браузере с терминалом, файловым менеджером и системным монитором
- 🎨 **Dark Theme** — профессиональный тёмный интерфейс
- 🔄 **Auto-sync** — синхронизация реестра прошивок из удалённого источника

---

## Installation

### One-liner (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/am1s3/opius-os/main/install.sh | sh
```

### From source

```bash
# Requirements: Go 1.21+
git clone https://github.com/am1s3/opius-os.git
cd opius-os
go build -o opius .
sudo mv opius /usr/local/bin/opius
```

### Initialize

```bash
opius init
```

---

## Quick Start

```bash
# Check system health
opius doctor

# List connected devices
opius device list

# Start web interface
opius web
```

Open http://127.0.0.1:8420 in your browser.

---

## CLI Commands

### Device Management

```bash
# List devices
opius device list
opius device list --all        # Include Bluetooth/virtual ports
opius device list --json       # JSON output

# Device info
opius device info
opius device info --port /dev/cu.usbserial-14110

# Diagnose connection issues
opius device doctor
```

### Flashing

```bash
# Flash from package
opius device flash --chip esp32 --file firmware.bin

# Flash with erase
opius device flash --chip esp32 --file firmware.bin --erase

# Flash custom binaries with memory addresses
opius device flash-custom --chip esp32 \
  --bin bootloader.bin@0x1000 \
  --bin partition-table.bin@0x8000 \
  --bin app.bin@0x10000

# Dry run (show plan without flashing)
opius device flash-custom --chip esp32 --bin app.bin@0x10000 --dry-run
```

### Serial Monitor

```bash
# Open monitor
opius device monitor

# With specific port and baud rate
opius device monitor --port /dev/cu.usbserial-14110 --baud 115200

# Log to file
opius device monitor --log monitor.log
```

### Device Operations

```bash
# Erase flash
opius device erase --chip esp32

# Reset device
opius device reset
```

### Package Manager

```bash
# Update registry from remote
opius pkg update

# Search packages
opius pkg search bruce

# Package info
opius pkg info bruce

# Install package
opius pkg install bruce

# List installed
opius pkg list

# Remove package
opius pkg remove bruce

# Install and flash in one command
opius pkg flash bruce
```

### Web Interface

```bash
# Start web server
opius web

# Custom port
opius web --port 9000

# Allow external access
opius web --host 0.0.0.0
```

---

## Web UI Apps

| App | Description |
|-----|-------------|
| 🖥️ Terminal | Full shell via WebSocket (zsh/bash) |
| 🔌 Devices | List and info about connected boards |
| 📱 Device Studio | Flash, erase, reset, serial monitor |
| 📦 Packages | Browse and sync firmware registry |
| 📁 Files | File manager for ~/.opius |
| 📊 System Monitor | CPU, memory, runtime stats |
| ⚙️ Settings | Configuration reference |
| ℹ️ About | Project info and links |

---

## Package Format

Packages are stored in `~/.opius/registry/` as TOML manifests:

```toml
name = "bruce"
version = "1.5.0"
description = "Multi-tool firmware for ESP32 boards"
license = "MIT"
chips = ["esp32", "esp32-s3"]
boards = ["lilygo-t-embed", "m5stack-cardputer"]

[[files]]
path = "firmware.bin"
sha256 = "9f2c..."
size = 1843200
type = "firmware"
```

---

## Memory Addresses Reference

Common ESP32 flash layout:

| Component | Address |
|-----------|---------|
| Bootloader | `0x1000` |
| Partition Table | `0x8000` |
| Application | `0x10000` |
| NVS Data | `0x9000` |
| PHY Init Data | `0xf000` |

---

## Configuration

Config file: `~/.config/opius/opius.toml`

```toml
package_backend = "macports"
workspace = "~/.opius/workspace"
registry_url = "https://github.com/am1s3/opius-os-pkg"
repository = "https://github.com/am1s3/opius-os"
```

---

## Requirements

- **Go** 1.21 or later (for building)
- **esptool** (for ESP flashing)
- **avrdude** (for AVR flashing)
- **picotool** (for RP2040 flashing)

### Install tools via MacPorts

```bash
sudo port install esptool avrdude picotool
```

### Install tools via Homebrew

```bash
brew install esptool avrdude picotool
```

---

## Project Structure

```
opius-os/
├── main.go                     # Entry point
├── install.sh                  # One-liner installer
├── Makefile                    # Build commands
├── internal/
│   ├── cli/                    # CLI commands
│   │   ├── root.go
│   │   ├── device.go
│   │   ├── device_flash.go
│   │   ├── device_flash_custom.go
│   │   ├── device_monitor.go
│   │   ├── device_erase.go
│   │   ├── device_reset.go
│   │   ├── pkg.go
│   │   ├── pkg_search.go
│   │   ├── pkg_install.go
│   │   ├── pkg_update.go
│   │   └── web.go
│   ├── device/                 # Device detection
│   │   └── detect.go
│   ├── flash/                  # Flash drivers
│   │   ├── driver.go
│   │   ├── esp.go
│   │   └── registry.go
│   ├── pkg/                    # Package manager
│   │   ├── manifest.go
│   │   ├── registry.go
│   │   ├── store.go
│   │   ├── hash.go
│   │   └── remote.go
│   ├── config/                 # Configuration
│   │   └── config.go
│   ├── hostpkg/                # Host package manager detection
│   │   └── detect.go
│   ├── tools/                  # Tool detection
│   │   └── detect.go
│   ├── build/                  # Build info
│   │   └── build.go
│   └── web/                    # Web server
│       ├── server.go
│       ├── api.go
│       ├── terminal.go
│       ├── serial_monitor.go
│       └── assets/
│           ├── index.html
│           ├── style.css
│           ├── window.js
│           ├── apps.js
│           └── desktop.js
└── README.md
```

---

## Troubleshooting

### Port busy error

```bash
# Find what's holding the port
lsof /dev/cu.usbserial-14110

# Kill the process
kill -9 <PID>

# Or release via web UI: Device Studio → Release Port
```

### Device not detected

```bash
opius device doctor
```

Check:
1. USB cable supports data (not charge-only)
2. Drivers installed (CP210x, CH340, FTDI)
3. Board is in normal mode (not bootloader)

### esptool not found

```bash
# MacPorts
sudo port install py312-esptool

# Homebrew
brew install esptool

# pip
pip install esptool
```

---

## Links

- Main Repository: https://github.com/am1s3/opius-os
- Package Registry: https://github.com/am1s3/opius-os-pkg

---

## License

MIT License — see [LICENSE](LICENSE) file for details.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss.

---

## Roadmap

- [x] Terminal CLI with colored output
- [x] Device detection (ESP32, Arduino, RP2040)
- [x] Flash firmware via esptool
- [x] Custom binary flashing with addresses
- [x] Serial monitor with WebSocket
- [x] Package manager with SHA-256 verification
- [x] Web UI desktop environment
- [x] Dark theme
- [x] Device Studio app
- [ ] OTA flashing
- [ ] Multi-device operations
- [ ] Debug probe support (JTAG/SWD)
- [ ] Firmware signing verification
- [ ] Plugin system

---

Built with Go, xterm.js, and ❤️ for embedded hackers.
