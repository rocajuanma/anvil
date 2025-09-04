#!/bin/bash

# Anvil Installation Script
# This script downloads and installs the latest version of Anvil

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="rocajuanma/anvil"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="anvil"

# Functions
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect architecture
detect_arch() {
    local arch=$(uname -m)
    case $arch in
        x86_64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Detect OS
detect_os() {
    local os=$(uname -s)
    case $os in
        Darwin)
            echo "darwin"
            ;;
        Linux)
            echo "linux"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | cut -d'"' -f4
}

# Download and install
install_anvil() {
    local os=$(detect_os)
    local arch=$(detect_arch)
    local version=$(get_latest_version)
    
    if [ -z "$version" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi
    
    print_status "Installing Anvil $version for $os-$arch"
    
    # Use architecture-specific binaries
    local binary_name="anvil-$os-$arch"
    local download_url="https://github.com/$REPO/releases/download/$version/$binary_name"
    
    print_status "Downloading from: $download_url"
    
    # Create temporary directory
    local tmp_dir=$(mktemp -d)
    local tmp_file="$tmp_dir/$BINARY_NAME"
    
    # Download binary
    if ! curl -L "$download_url" -o "$tmp_file"; then
        print_error "Failed to download Anvil"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    # Make executable
    chmod +x "$tmp_file"
    
    # Test the binary
    if ! "$tmp_file" --help > /dev/null 2>&1; then
        print_error "Downloaded binary is not working correctly"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    # Install binary
    print_status "Installing to $INSTALL_DIR/$BINARY_NAME"
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_file" "$INSTALL_DIR/$BINARY_NAME"
    else
        print_status "Need sudo permissions to install to $INSTALL_DIR"
        sudo mv "$tmp_file" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup
    rm -rf "$tmp_dir"
    
    print_success "Anvil $version installed successfully!"
    print_status "Run 'anvil' to verify installation"
    print_status "Run 'anvil init' to get started"
}

# Main
main() {
    echo ""
    echo "🔨 Anvil Installation Script"
    echo "================================"
    echo ""
    
    # Check dependencies
    if ! command -v curl > /dev/null 2>&1; then
        print_error "curl is required but not installed"
        exit 1
    fi
    
    # Check if already installed
    if command -v anvil > /dev/null 2>&1; then
        local current_version=$(anvil --version 2>/dev/null | head -n1 || echo "unknown")
        print_status "Anvil is already installed: $current_version"
        print_status "Updating to latest version..."
    fi
    
    install_anvil
    
    echo ""
    echo "🎉 Installation Complete!"
    echo ""
    echo "Next steps:"
    echo "  1. Run 'anvil init' to initialize your environment"
    echo "  2. Run 'anvil doctor' to verify everything is working"
    echo "  3. Run 'anvil install --list' to see available tool groups. Create your own groups in settings.yaml."
    echo ""
}

main "$@"
