# Clean Command

The `anvil clean` command provides safe cleanup of temporary files, archives, and cached configurations while preserving your essential settings and directory structure.

## Overview

The clean command serves as a maintenance utility that helps you:

- **🧹 Free up disk space** - Removes temporary files, old archives, and cached configurations
- **🔄 Reset to clean state** - Clean slate for configuration operations while preserving settings
- **📁 Smart preservation** - Keeps essential `settings.yaml` and maintains directory structure
- **🛡️ Safe operations** - Interactive confirmations and dry-run support prevent accidental deletions
- **🗂️ Complete git cleanup** - Fully removes dotfiles directory for clean repository state

## Usage

### Basic Cleanup

```bash
# Interactive cleanup with confirmation
anvil clean

# Force cleanup without confirmation
anvil clean --force

# Preview what would be cleaned without deletion
anvil clean --dry-run
```

## What Gets Cleaned

The clean command targets specific content while preserving essential files:

### **Cleaned Content**

- **temp/ directory contents** - All pulled configurations waiting to be synced
- **archive/ directory contents** - Old archived configurations and backups
- **dotfiles/ directory** - Completely removed to ensure clean git repository state
- **Other root files/directories** - Any additional files created in ~/.anvil (except settings.yaml)

### **Preserved Content**

- **settings.yaml** - Your main configuration file with all settings
- **Directory structure** - Essential directories like temp/ and archive/ are preserved for tool functionality

## How It Works

### Stage 1: Discovery and Analysis

The clean command scans your ~/.anvil directory and identifies:

1. **Content to clean** - All items except settings.yaml
2. **Directory structures** - Shows tree view of what will be removed
3. **Item counts** - Displays how many files/directories will be affected

### Stage 2: Safety Confirmation

Unless using `--force` flag:

1. **Interactive confirmation** - Shows exactly what will be cleaned
2. **Clear summary** - Displays count of directories/files to be processed
3. **Cancellation option** - Easy way to abort if unsure

### Stage 3: Smart Cleanup

Different handling based on content type:

- **dotfiles/ directory** - Completely removed to ensure clean git repository state for next pull/push
- **temp/ and archive/ directories** - Contents cleaned but directory structure preserved
- **Other files** - Individual files removed as needed

## Directory Structure Impact

### Before Cleaning

```
~/.anvil/
├── settings.yaml           # Preserved
├── temp/                   # Structure preserved
│   ├── cursor/
│   │   ├── settings.json
│   │   └── keybindings.json
│   └── vscode/
│       └── settings.json
├── archive/                # Structure preserved
│   ├── 2025-01-15-1430/
│   │   └── old-configs/
│   └── 2025-01-10-0900/
└── dotfiles/               # Completely removed
    ├── .git/
    ├── cursor/
    └── vscode/
```

### After Cleaning

```
~/.anvil/
├── settings.yaml           # Preserved
├── temp/                   # Empty but preserved
└── archive/                # Empty but preserved
```

## Use Cases

### Disk Space Management

Free up space consumed by temporary files and archives:

```bash
# Check what's taking up space
anvil clean --dry-run

# Clean up to free disk space
anvil clean
```

### Development and Testing

Reset to clean state during development or testing:

```bash
# Clean up after testing configuration operations
anvil clean --force

# Preview cleanup before development work
anvil clean --dry-run
```

### Configuration Reset

Start fresh with configuration operations:

```bash
# Clean up before major configuration changes
anvil clean

# Ensure clean state before team configuration sync
anvil clean --dry-run
anvil clean
```

### Git Repository Cleanup

Remove dotfiles directory to ensure clean git operations:

```bash
# Clean before pull to avoid git conflicts
anvil clean
anvil config pull cursor

# Clean before push to ensure clean commit state
anvil clean --force
anvil config push
```

## Examples

### Interactive Cleanup with Detailed Preview

```bash
$ anvil clean

=== Cleaning Anvil Directories ===

🔧 Scanning .anvil directory for content to clean

Found 3 root directories to clean:
Directory structure to be cleaned:
  📁 temp (2)
    ├── cursor
    ├── vscode
  📁 archive (1)
    ├── 2025-01-15-1430
  📁 dotfiles (5)
    ├── .git
    ├── cursor
    ├── vscode
    ├── README.md
    ├── .gitignore

Are you sure you want to clean the contents of these 3 root directories? This action cannot be undone [y/N]: y

🔧 Cleaning directories and files

✅ Cleaned contents of directory temp
✅ Cleaned contents of directory archive
✅ Removed dotfiles directory completely

Successfully cleaned contents of 3/3 root directories
```

### Dry-Run Preview

```bash
$ anvil clean --dry-run

=== Cleaning Anvil Directories ===

🔧 Scanning .anvil directory for content to clean

Found 2 root directories to clean:
Directory structure to be cleaned:
  📁 temp (3)
    ├── cursor
    ├── vscode
    ├── neovim
  📁 archive (2)
    ├── 2025-01-15-1430
    ├── 2025-01-10-0900

DRY RUN: Would clean contents of 2 root directories
```

### Force Cleanup for Scripts

```bash
$ anvil clean --force

=== Cleaning Anvil Directories ===

🔧 Scanning .anvil directory for content to clean

Found 1 root directories to clean:
Directory structure to be cleaned:
  📁 dotfiles (3)
    ├── .git
    ├── configurations
    ├── README.md

🔧 Cleaning directories and files

✅ Removed dotfiles directory completely

Successfully cleaned contents of 1/1 root directories
```

### No Content to Clean

```bash
$ anvil clean

=== Cleaning Anvil Directories ===

🔧 Scanning .anvil directory for content to clean

✅ No root directories found to clean. Only settings.yaml exists.
```

## Safety Features

### Interactive Confirmation

- **Clear prompts** - Shows exactly what will be deleted
- **Easy cancellation** - Default behavior is to cancel unless explicitly confirmed
- **Detailed preview** - Tree structure showing all content that will be removed

### Dry-Run Mode

- **Safe preview** - See exactly what would be cleaned without deletion
- **Identical logic** - Dry-run uses same detection logic as real cleanup
- **Risk-free exploration** - Perfect for understanding impact before cleanup

### Settings Preservation

- **Always preserved** - settings.yaml is never touched by clean operation
- **Configuration safety** - All your tool groups, GitHub settings, and preferences remain intact
- **Quick recovery** - No need to reconfigure after cleanup

### Directory Structure Maintenance

- **Essential directories preserved** - temp/ and archive/ remain for tool functionality
- **Clean slate ready** - Directories are empty but ready for immediate use
- **No reinitialization needed** - Tool continues working without `anvil init`

## Integration with Other Commands

The clean command integrates seamlessly with configuration management:

### Before Configuration Operations

```bash
# Clean before pulling to ensure fresh state
anvil clean
anvil config pull cursor

# Clean before pushing to avoid git conflicts
anvil clean --force
anvil config push
```

### After Development Work

```bash
# Clean up after testing configurations
anvil config pull test-configs
# ... testing ...
anvil clean

# Clean up after configuration development
anvil clean --dry-run
anvil clean
```

### Maintenance Workflows

```bash
# Regular maintenance cleanup
anvil clean --dry-run    # Preview
anvil clean             # Clean with confirmation

# Automated cleanup in scripts
anvil clean --force     # No interaction required
```

## Best Practices

### 🎯 Regular Maintenance

1. **Run periodic cleanups** to prevent disk space accumulation
2. **Use dry-run first** to understand what will be removed
3. **Clean before major operations** to ensure clean starting state

### 🔧 Development Workflow

1. **Clean after testing** to remove temporary configurations
2. **Clean before git operations** to avoid conflicts
3. **Use force flag in scripts** for automation

### 📁 Configuration Management

1. **Clean before pulling** to ensure fresh configuration state
2. **Clean before pushing** to avoid git repository conflicts
3. **Regular cleanup** of archived configurations to manage disk space

### 🛡️ Safety Practices

1. **Always preview with dry-run** when unsure
2. **Backup important configurations** outside of ~/.anvil before cleaning
3. **Understand that dotfiles are completely removed** for clean git state

## Troubleshooting

### Common Issues

**Permission Denied**

```bash
# Check directory permissions
ls -la ~/.anvil

# Fix permissions if needed
chmod 755 ~/.anvil
chmod -R 644 ~/.anvil/*
```

**Directory Not Found**

```bash
# Check if .anvil directory exists
ls -la ~/.anvil

# Initialize if needed
anvil init
```

**Partial Cleanup**

```bash
# Some files couldn't be cleaned - check warnings
# Try cleaning individual directories manually
rm -rf ~/.anvil/temp/*
rm -rf ~/.anvil/archive/*
```

**Settings File Accidentally Deleted**

```bash
# Reinitialize to restore default settings
anvil init

# Restore from backup if available
cp ~/.anvil/settings.yaml.backup ~/.anvil/settings.yaml
```

### Getting Help

For detailed setup guidance and examples, see:

- [Getting Started Guide](GETTING_STARTED.md)
- [Configuration Management](config.md)
- [Doctor Command](doctor.md) - For diagnosing cleanup issues

## Security Notes

- **No credential handling** - Clean command never touches authentication files
- **Settings preservation** - Your GitHub tokens and SSH key configurations remain intact
- **Repository cleanup** - Complete dotfiles removal ensures clean git repository state
- **Safe by default** - Interactive confirmation prevents accidental deletions
