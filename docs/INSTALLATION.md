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

## Quick Installation

### Method 1: Direct Download (Recommended)

Download the latest release binary from [GitHub Releases](https://github.com/rocajuanma/anvil/releases):

```bash
curl -sSL https://github.com/rocajuanma/anvil/releases/latest/download/install.sh | bash
```

### Method 2: Download and Build

```bash
# Clone the repository
git clone https://github.com/rocajuanma/anvil.git
cd anvil

# Build the binary
go build -o anvil main.go

# Move to PATH (optional)
sudo mv anvil /usr/local/bin/
```

## Development Installation

For contributing to Anvil or customizing functionality:

### Setup Development Environment

```bash
# Build using Method #2

# Install dependencies
go mod download
go mod tidy

# Run tests
go test ./...

# Install development tools (optional)
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
```


## Verification

After installation, verify Anvil is working correctly:

### Basic Verification

```bash
# Check version and help
anvil -v
anvil --help

# Test initialization (this should work without errors)
anvil init

anvil doctor
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
anvil config show
```

### Health Check

```bash
# Verify all components and runs checks
anvil doctor
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
2. **[Initialize Anvil](init.md)** - Set up your environment
3. **[Install Tools](install.md)** - Start installing development tools
4. **[View Examples](EXAMPLES.md)** - See real-world usage examples

## Keeping Anvil Updated

Once you have Anvil v1.2.0 or later installed, you can easily update to newer versions:

```bash
# Update to the latest version
anvil update

# Preview what would be updated
anvil update --dry-run
```

For versions prior to v1.2.0, use the curl installation(Method 1) in the first section to upgrade to the latest version.

