#!/bin/sh
# Opius OS Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/am1s3/opius-os/main/install.sh | sh

set -e

REPO="am1s3/opius-os"
BINARY_NAME="opius"
INSTALL_DIR=""

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()  { echo "${BLUE}→${NC} $1"; }
ok()    { echo "${GREEN}✓${NC} $1"; }
warn()  { echo "${YELLOW}⚠${NC} $1"; }
fail()  { echo "${RED}✗${NC} $1"; exit 1; }

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) fail "Unsupported OS: $OS" ;;
esac

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) fail "Unsupported architecture: $ARCH" ;;
esac

echo ""
echo "╔═══════════════════════════════════════╗"
echo "║        Opius OS Installer             ║"
echo "╚═══════════════════════════════════════╝"
echo ""
info "OS:   $OS"
info "ARCH: $ARCH"
echo ""

# Determine install directory (prefer user-local if no sudo)
if [ -w "/usr/local/bin" ] 2>/dev/null; then
    INSTALL_DIR="/usr/local/bin"
    NEED_SUDO=""
elif [ -w "$HOME/.local/bin" ] 2>/dev/null || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
    INSTALL_DIR="$HOME/.local/bin"
    NEED_SUDO=""
elif [ -w "$HOME/bin" ] 2>/dev/null || mkdir -p "$HOME/bin" 2>/dev/null; then
    INSTALL_DIR="$HOME/bin"
    NEED_SUDO=""
else
    INSTALL_DIR="/usr/local/bin"
    NEED_SUDO="sudo"
fi

# Try to download prebuilt binary
try_download_release() {
    info "Checking for prebuilt binaries..."

    LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
    RELEASE_INFO=$(curl -fsSL "$LATEST_URL" 2>/dev/null) || return 1

    ASSET_URL=$(echo "$RELEASE_INFO" | grep -o "\"browser_download_url\": \"[^\"]*${OS}_${ARCH}[^\"]*\"" | cut -d '"' -f 4)

    if [ -z "$ASSET_URL" ]; then
        return 1
    fi

    info "Downloading prebuilt binary..."

    TEMP_DIR=$(mktemp -d)
    ARCHIVE="$TEMP_DIR/opius.tar.gz"

    curl -fsSL "$ASSET_URL" -o "$ARCHIVE" || return 1

    info "Extracting..."
    tar -xzf "$ARCHIVE" -C "$TEMP_DIR"

    BINARY_PATH=$(find "$TEMP_DIR" -name "$BINARY_NAME" -type f | head -1)

    if [ -z "$BINARY_PATH" ]; then
        rm -rf "$TEMP_DIR"
        return 1
    fi

    chmod +x "$BINARY_PATH"

    if [ -n "$NEED_SUDO" ]; then
        warn "Need sudo access to install to $INSTALL_DIR"
        $NEED_SUDO mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME" || return 1
    else
        mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME" || return 1
    fi

    rm -rf "$TEMP_DIR"
    return 0
}

# Install Go if not present
install_go() {
    info "Go not found. Installing Go automatically..."

    GO_VERSION="1.21.13"

    # Determine download URL
    if [ "$OS" = "darwin" ] && [ "$ARCH" = "arm64" ]; then
        GO_URL="https://go.dev/dl/go${GO_VERSION}.darwin-arm64.tar.gz"
    elif [ "$OS" = "darwin" ] && [ "$ARCH" = "amd64" ]; then
        GO_URL="https://go.dev/dl/go${GO_VERSION}.darwin-amd64.tar.gz"
    elif [ "$OS" = "linux" ] && [ "$ARCH" = "arm64" ]; then
        GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-arm64.tar.gz"
    else
        GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
    fi

    GO_INSTALL_DIR="$HOME/.local/go"
    mkdir -p "$GO_INSTALL_DIR"
    GO_INSTALL_DIR="$(cd "$GO_INSTALL_DIR" && pwd)"
    # Use standard location
    GO_INSTALL_DIR="$HOME/go-install"

    info "Downloading Go $GO_VERSION..."
    GO_ARCHIVE=$(mktemp)
    curl -fsSL "$GO_URL" -o "$GO_ARCHIVE" || fail "Failed to download Go"

    info "Extracting Go..."
    mkdir -p "$GO_INSTALL_DIR"
    tar -xzf "$GO_ARCHIVE" -C "$GO_INSTALL_DIR"
    rm "$GO_ARCHIVE"

    # Set up Go paths
    export GOROOT="$GO_INSTALL_DIR/go"
    export GOPATH="$HOME/go"
    export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"

    # Add to shell profile
    SHELL_RC=""
    if [ -n "$ZSH_VERSION" ] || [ "$(basename "$SHELL")" = "zsh" ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ] || [ "$(basename "$SHELL")" = "bash" ]; then
        SHELL_RC="$HOME/.bashrc"
    else
        SHELL_RC="$HOME/.profile"
    fi

    if [ -f "$SHELL_RC" ] && grep -q "GOROOT=$GOROOT" "$SHELL_RC" 2>/dev/null; then
        : # Already configured
    else
        echo "" >> "$SHELL_RC"
        echo "# Go (installed by Opius)" >> "$SHELL_RC"
        echo "export GOROOT=\"$GOROOT\"" >> "$SHELL_RC"
        echo "export GOPATH=\"\$HOME/go\"" >> "$SHELL_RC"
        echo "export PATH=\"\$GOROOT/bin:\$GOPATH/bin:\$PATH\"" >> "$SHELL_RC"
        info "Go paths added to $SHELL_RC"
    fi

    ok "Go $GO_VERSION installed to $GOROOT"

    # Verify
    if ! go version > /dev/null 2>&1; then
        fail "Go installation failed verification"
    fi
}

# Build from source
build_from_source() {
    info "Building Opius from source..."

    # Check/install Go
    if ! command -v go > /dev/null 2>&1; then
        install_go
    fi

    info "Go version: $(go version)"

    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    info "Cloning repository..."
    git clone --depth 1 "https://github.com/$REPO.git" opius-src || fail "Failed to clone repository"

    cd opius-src

    info "Building binary..."
    go build -ldflags "-X github.com/opius-os/opius/internal/build.Version=1.0.0" -o "$BINARY_NAME" . || fail "Build failed"

    info "Installing to $INSTALL_DIR..."
    chmod +x "$BINARY_NAME"

    if [ -n "$NEED_SUDO" ]; then
        warn "Need sudo access to install to $INSTALL_DIR"
        $NEED_SUDO mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME" || fail "Install failed"
    else
        mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME" || fail "Install failed"
    fi

    cd /
    rm -rf "$TEMP_DIR"
}

# Main logic
if try_download_release; then
    ok "Opius OS installed from prebuilt binary"
else
    warn "No prebuilt binary available, building from source..."
    build_from_source
    ok "Opius OS built and installed from source"
fi

# Add install dir to PATH if needed
case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
        warn "$INSTALL_DIR is not in your PATH"

        SHELL_RC=""
        if [ -n "$ZSH_VERSION" ] || [ "$(basename "$SHELL" 2>/dev/null)" = "zsh" ]; then
            SHELL_RC="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ] || [ "$(basename "$SHELL" 2>/dev/null)" = "bash" ]; then
            SHELL_RC="$HOME/.bashrc"
        else
            SHELL_RC="$HOME/.profile"
        fi

        if [ -f "$SHELL_RC" ] && grep -q "export PATH=.*$INSTALL_DIR" "$SHELL_RC" 2>/dev/null; then
            : # Already configured
        else
            echo "" >> "$SHELL_RC"
            echo "# Opius OS" >> "$SHELL_RC"
            echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_RC"
            info "Added $INSTALL_DIR to PATH in $SHELL_RC"
            info "Run 'source $SHELL_RC' or open a new terminal to use 'opius'"
        fi
        ;;
esac

echo ""
echo "╔═══════════════════════════════════════╗"
echo "║      Installation Complete!           ║"
echo "╚═══════════════════════════════════════╝"
echo ""
echo "  Binary: $INSTALL_DIR/opius"
echo ""
echo "  Get started:"
echo "    opius init          # Initialize workspace"
echo "    opius doctor        # Check system health"
echo "    opius device list   # List connected devices"
echo "    opius web           # Start web interface"
echo ""
