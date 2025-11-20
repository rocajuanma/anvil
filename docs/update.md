# Update Command

The `anvil update` command provides a simple and safe way to update your Anvil installation to the latest version.

## Overview

The update command downloads and executes the official Anvil installation script from GitHub releases, ensuring you always get the latest version with all the newest features and bug fixes.

## Usage

### Basic Syntax

```bash
anvil update [flags]
```

### Available Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview the update process without actually updating |
| `--help` | Show help information for the update command |

## Examples

### Update to Latest Version

```bash
# Update Anvil to the latest version
anvil update
```

This will:
- Download the latest release information from GitHub
- Fetch and execute the official installation script
- Replace your current Anvil binary with the latest version
- Provide clear feedback on the update process

### Preview Update Process

```bash
# See what would be updated without actually updating
anvil update --dry-run
```

This will:
- Show the command that would be executed
- Display information about the update process
- Exit without making any changes

## Features

### Safe Update Process

The update command uses the same trusted installation script that you used for the initial installation, ensuring consistency and reliability.

### Dry-Run Support

Preview exactly what will happen during the update process without making any changes to your system.

### Automatic Detection

The update script automatically detects your system architecture (Intel or Apple Silicon) and downloads the appropriate binary.

### macOS Optimized

Specifically designed for macOS environments, following the same patterns as other Anvil commands.

### Clear Feedback

Provides step-by-step progress updates and helpful instructions throughout the update process.

## How It Works

The update command follows these steps:

1. **Platform Check**: Verifies you're running on macOS (required for Anvil)
2. **Dependency Check**: Ensures `curl` is available on your system
3. **Download Script**: Fetches the latest installation script from GitHub releases
4. **Execute Update**: Runs the script to install the latest Anvil binary
5. **Validation**: Confirms the update was successful

### The Update Command

Under the hood, the update process executes:

```bash
curl -sSL https://github.com/0xjuanma/anvil/releases/latest/download/install.sh | bash
```

This is the same command recommended in the README for fresh installations and updates.

## Safety Considerations

### Safe Practices

- **Trusted Source**: Downloads only from official GitHub releases
- **Same Script**: Uses identical installation script as initial setup
- **Validation**: Verifies the update process completed successfully
- **No Data Loss**: Your `~/.anvil/settings.yaml` and configurations remain unchanged

### Important Notes

- **Terminal Restart**: You may need to restart your terminal session for changes to take effect
- **Admin Permissions**: The script may request admin permissions to install to `/usr/local/bin/`
- **Internet Required**: Requires active internet connection to download the latest version

## Troubleshooting

### Common Issues

#### "curl is required but not available"

**Cause**: The `curl` command is not installed or not in your PATH.

**Solution**: 
```bash
# Install curl using Homebrew
brew install curl

# Or check if it's already available
which curl
```

#### "Update script failed"

**Cause**: Network issues, GitHub API rate limits, or permission problems.

**Solutions**:
1. **Check Internet Connection**: Ensure you have a stable internet connection
2. **Try Again**: GitHub API has rate limits; wait a few minutes and retry
3. **Manual Installation**: Download the binary manually from [GitHub releases](https://github.com/0xjuanma/anvil/releases)

#### "Permission Denied"

**Cause**: Insufficient permissions to write to `/usr/local/bin/`.

**Solution**: The installation script will automatically request admin permissions when needed.

#### Update Doesn't Take Effect

**Cause**: Your terminal is using a cached version of the Anvil binary.

**Solutions**:
1. **Restart Terminal**: Close and reopen your terminal application
2. **Reload PATH**: Run `hash -r` to clear the command cache
3. **Verify Installation**: Check with `which anvil` and `anvil --version`

### Getting Help

If you continue to experience issues:

1. **Check Version**: Run `anvil --version` to verify your current version
2. **Check Installation**: Run `which anvil` to see where Anvil is installed
3. **Manual Download**: Visit [GitHub releases](https://github.com/0xjuanma/anvil/releases) for manual installation
4. **Report Issues**: Open an issue on the [GitHub repository](https://github.com/0xjuanma/anvil/issues)

## Integration with Other Commands

The update command works seamlessly with other Anvil commands:

```bash
# Update Anvil and verify the new version
anvil update
anvil --version

# Run health checks after update
anvil doctor

# Continue with your normal workflow
anvil install dev
anvil config pull
```

---

**Next Steps**: After updating, consider running `anvil doctor` to ensure your environment is properly configured with the latest version.
