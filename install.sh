#!/bin/sh
# Opius OS Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/am1s3/opius-os/main/install.sh | sh

set -e

REPO="am1s3/opius-os"
BINARY_NAME="opius"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "Opius OS Installer"
echo "=================="
echo "OS:   $OS"
echo "ARCH: $ARCH"
echo ""

# Get latest release URL
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
ASSET_URL=$(curl -fsSL "$LATEST_URL" | grep -o "\"browser_download_url\": \"[^\"]*${OS}_${ARCH}[^\"]*\"" | cut -d '"' -f 4)

if [ -z "$ASSET_URL" ]; then
    echo "Error: Could not find release for $OS/$ARCH"
    echo "Falling back to source build..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo "Error: Go is not installed. Please install Go from https://go.dev/dl/"
        exit 1
    fi
    
    echo "Building from source..."
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    git clone https://github.com/$REPO.git
    cd opius-os
    go build -o "$BINARY_NAME" .
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    else
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    cd /
    rm -rf "$TEMP_DIR"
    
    echo ""
    echo "✓ Opius OS installed to $INSTALL_DIR/$BINARY_NAME"
    echo ""
    echo "Run 'opius init' to get started."
    exit 0
fi

echo "Downloading from: $ASSET_URL"

# Download and extract
TEMP_DIR=$(mktemp -d)
ARCHIVE="$TEMP_DIR/opius.tar.gz"

curl -fsSL "$ASSET_URL" -o "$ARCHIVE"

echo "Extracting..."
tar -xzf "$ARCHIVE" -C "$TEMP_DIR"

# Find binary
BINARY_PATH=$(find "$TEMP_DIR" -name "$BINARY_NAME" -type f | head -1)

if [ -z "$BINARY_PATH" ]; then
    echo "Error: Could not find binary in archive"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Install
echo "Installing to $INSTALL_DIR..."
chmod +x "$BINARY_PATH"

if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
fi

# Cleanup
rm -rf "$TEMP_DIR"

echo ""
echo "✓ Opius OS installed successfully!"
echo ""
echo "Get started:"
echo "  opius init          # Initialize workspace"
echo "  opius doctor        # Check system health"
echo "  opius device list   # List connected devices"
echo "  opius web           # Start web interface"
echo ""