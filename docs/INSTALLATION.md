# Installation Guide

This guide provides detailed installation instructions for Anvil CLI across different platforms and scenarios.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Installation](#quick-installation)
- [Platform-Specific Installation](#platform-specific-installation)
- [Development Installation](#development-installation)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)
- [Uninstallation](#uninstallation)

## Prerequisites

### System Requirements

- **Go**: Version 1.17 or higher
- **Git**: Required for asset synchronization features
- **Internet Connection**: For downloading dependencies and tools
- **Terminal/Command Line**: Basic familiarity recommended

### Platform-Specific Requirements

#### macOS

- **Homebrew**: Installed automatically by Anvil if missing
- **Xcode Command Line Tools**: `xcode-select --install`
- **Admin privileges**: For tool installations

#### Linux

- **Package manager**: apt, yum, or equivalent
- **Build tools**: `build-essential` (Ubuntu) or `Development Tools` (CentOS)
- **curl**: Usually pre-installed

#### Windows

- **Git Bash** or **WSL**: Recommended for best experience
- **PowerShell**: Windows PowerShell 5.1 or PowerShell Core 6+

## Quick Installation

### Method 1: Download and Build (Recommended)

```bash
# Clone the repository
git clone https://github.com/rocajuanma/anvil.git
cd anvil

# Build the binary
go build -o anvil main.go

# Make it executable (Linux/macOS)
chmod +x anvil

# Move to PATH (optional)
sudo mv anvil /usr/local/bin/
```

### Method 2: Go Install (if available)

```bash
go install github.com/rocajuanma/anvil@latest
```

### Method 3: Direct Binary Download

Download the latest release binary from [GitHub Releases](https://github.com/rocajuanma/anvil/releases):

```bash
# macOS (Intel)
curl -L https://github.com/rocajuanma/anvil/releases/latest/download/anvil-darwin-amd64 -o anvil

# macOS (Apple Silicon)
curl -L https://github.com/rocajuanma/anvil/releases/latest/download/anvil-darwin-arm64 -o anvil

# Linux (x86_64)
curl -L https://github.com/rocajuanma/anvil/releases/latest/download/anvil-linux-amd64 -o anvil

# Make executable and move to PATH
chmod +x anvil
sudo mv anvil /usr/local/bin/
```

## Platform-Specific Installation

### macOS Installation

#### Using Homebrew (Future)

```bash
# Once available in Homebrew
brew tap rocajuanma/anvil
brew install anvil
```

#### Manual Installation

```bash
# Install dependencies
xcode-select --install

# Clone and build
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go build -o anvil main.go

# Install to system
sudo mv anvil /usr/local/bin/
sudo chmod +x /usr/local/bin/anvil

# Verify installation
anvil --version
```

#### macOS Security Notice

If you get a security warning when running Anvil:

1. Go to **System Preferences** → **Security & Privacy**
2. Click **Allow Anyway** next to the Anvil security notice
3. Or run: `sudo spctl --add /usr/local/bin/anvil`

### Linux Installation

#### Ubuntu/Debian

```bash
# Install dependencies
sudo apt update
sudo apt install -y git build-essential curl

# Install Go (if not installed)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Clone and build Anvil
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go build -o anvil main.go

# Install system-wide
sudo mv anvil /usr/local/bin/
sudo chmod +x /usr/local/bin/anvil
```

#### CentOS/RHEL/Fedora

```bash
# Install dependencies
sudo yum groupinstall -y "Development Tools"
sudo yum install -y git curl

# Or for newer versions
sudo dnf groupinstall -y "Development Tools"
sudo dnf install -y git curl

# Install Go and Anvil (same as Ubuntu steps above)
```

#### Arch Linux

```bash
# Install dependencies
sudo pacman -S git base-devel go

# Clone and build
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go build -o anvil main.go
sudo mv anvil /usr/local/bin/
```

### Windows Installation

#### Using Git Bash (Recommended)

```bash
# Install Git for Windows first: https://git-scm.com/download/win
# Open Git Bash

# Install Go: https://golang.org/dl/
# Clone and build
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go build -o anvil.exe main.go

# Move to a directory in your PATH
mv anvil.exe /c/Windows/System32/
```

#### Using PowerShell

```powershell
# Clone repository
git clone https://github.com/rocajuanma/anvil.git
cd anvil

# Build
go build -o anvil.exe main.go

# Add to PATH (run as Administrator)
$env:PATH += ";C:\path\to\anvil"
[Environment]::SetEnvironmentVariable("PATH", $env:PATH, [EnvironmentVariableTarget]::Machine)
```

#### Using WSL (Windows Subsystem for Linux)

```bash
# Follow the Linux installation instructions inside WSL
# WSL provides the best Anvil experience on Windows
```

## Development Installation

For contributing to Anvil or customizing functionality:

### Setup Development Environment

```bash
# Clone with development setup
git clone https://github.com/rocajuanma/anvil.git
cd anvil

# Install dependencies
go mod download
go mod tidy

# Build development version
go build -o anvil-dev main.go

# Run tests
go test ./...

# Install development tools (optional)
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
```

### Development Dependencies

```bash
# Code formatting
go install golang.org/x/tools/cmd/goimports@latest

# Linting
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Testing
go install gotest.tools/gotestsum@latest
```

### Building for Multiple Platforms

```bash
# Build for all platforms
./scripts/build-all.sh

# Or manually
GOOS=darwin GOARCH=amd64 go build -o dist/anvil-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o dist/anvil-darwin-arm64 main.go
GOOS=linux GOARCH=amd64 go build -o dist/anvil-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o dist/anvil-windows-amd64.exe main.go
```

## Verification

After installation, verify Anvil is working correctly:

### Basic Verification

```bash
# Check version and help
anvil --help

# Verify main commands are available
anvil init --help
anvil install --help
anvil config --help
anvil config show --help
anvil config sync --help

# Test initialization (this should work without errors)
anvil init
```

### Extended Verification

```bash
# Initialize Anvil
anvil init

# List available tools
anvil install --list

# Test dry run
anvil install dev --dry-run

# Check configuration
cat ~/.anvil/settings.yaml
```

### Health Check

```bash
# Verify all components
anvil init
anvil install --list

# Check configuration directory
ls -la ~/.anvil/
```

## Troubleshooting

### Common Installation Issues

#### Go Not Found

```bash
# Install Go
# macOS with Homebrew
brew install go

# Ubuntu/Debian
sudo apt install golang-go

# Or download from https://golang.org/dl/
```

#### Permission Denied

```bash
# Fix permissions
chmod +x anvil
sudo chown $(whoami) anvil

# For system installation
sudo mv anvil /usr/local/bin/
```

#### PATH Issues

```bash
# Add to PATH temporarily
export PATH=$PATH:/path/to/anvil

# Add to PATH permanently
echo 'export PATH=$PATH:/path/to/anvil' >> ~/.bashrc
source ~/.bashrc
```

#### Build Failures

```bash
# Update Go modules
go mod download
go mod tidy

# Clean build cache
go clean -cache
go clean -modcache

# Rebuild
go build -o anvil main.go
```

### Platform-Specific Issues

#### macOS: "Cannot be opened because it is from an unidentified developer"

```bash
# Allow the binary
sudo spctl --add /path/to/anvil

# Or use system preferences method described above
```

#### Linux: Missing Dependencies

```bash
# Ubuntu/Debian
sudo apt install -y build-essential git curl

# CentOS/RHEL
sudo yum groupinstall -y "Development Tools"
sudo yum install -y git curl
```

#### Windows: Execution Policy

```powershell
# Allow script execution
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Getting Help

If you encounter issues not covered here:

1. **Check existing issues**: [GitHub Issues](https://github.com/rocajuanma/anvil/issues)
2. **Create new issue**: Include your platform, Go version, and error messages
3. **Join discussions**: [GitHub Discussions](https://github.com/rocajuanma/anvil/discussions)

## Uninstallation

### Remove Anvil Binary

```bash
# If installed to /usr/local/bin/
sudo rm /usr/local/bin/anvil

# If installed elsewhere
which anvil  # Find location
rm /path/to/anvil
```

### Remove Configuration

```bash
# Remove Anvil configuration
rm -rf ~/.anvil

# Remove from PATH (edit your shell configuration)
# Remove anvil PATH entry from ~/.bashrc, ~/.zshrc, etc.
```

### Complete Cleanup

```bash
# Remove binary
sudo rm /usr/local/bin/anvil

# Remove configuration
rm -rf ~/.anvil

# Remove source code (if cloned)
rm -rf /path/to/anvil-source

# Clean Go cache (optional)
go clean -cache
```

---

## Next Steps

After successful installation:

1. **[Get Started](GETTING_STARTED.md)** - Learn basic Anvil usage
2. **[Initialize Anvil](init-readme.md)** - Set up your environment
3. **[Install Tools](setup-readme.md)** - Start installing development tools
4. **[View Examples](EXAMPLES.md)** - See real-world usage examples

---

**Need help?** Check our [troubleshooting guide](GETTING_STARTED.md#troubleshooting) or [open an issue](https://github.com/rocajuanma/anvil/issues).
