# Install Command

The `anvil install` command provides automated installation of development tools and applications using Homebrew (on macOS). This command streamlines the process of setting up consistent development environments by organizing tools into logical groups.

## Overview

The install command serves several critical functions:

- **🎯 Dynamic installation** - Install any macOS application with `anvil install [app-name]`
- **📝 Smart tracking** - Individual apps automatically tracked in `tools.installed_apps`
- **📦 Group management** - Predefined and custom tool groups for common scenarios
- **🚀 Zero configuration** - Works out of the box with sensible defaults
- **🔍 Dry-run support** - Preview installations before execution
- **🍺 Homebrew integration** - Automatic installation and management

## Usage Modes

### Individual Application Installation

Install any application by name with automatic tracking:

```bash
anvil install firefox
anvil install slack
anvil install figma
anvil install visual-studio-code
```

**Features:**

- Apps are automatically tracked in `tools.installed_apps` in your settings.yaml
- Smart deduplication prevents tracking apps already in groups or required_tools
- Works with any Homebrew package name

### Individual Application Installation with Group Assignment

Install an application and automatically add it to a group:

```bash
# Add to existing group
anvil install firefox --group-name essentials

# Create new group and add app
anvil install final-cut --group-name editing
```

**Features:**

- Installs the application using existing logic
- Adds the app to the specified group if installation is successful
- Creates the group if it doesn't exist
- Prevents duplicate apps within the same group
- Falls back to normal `installed_apps` tracking if group operation fails

### Group Installation

Install all tools in a predefined or custom group:

```bash
anvil install dev         # Development tools
anvil install new-laptop  # Essential applications for new machines
```

### List Available Options

See all available groups and tracked apps:

```bash
anvil install --list
```

### Preview Mode

Preview installations before execution:

```bash
anvil install docker --dry-run
anvil install dev --dry-run
```

## Available Groups

### Default Groups

**dev** - Essential development tools:

```yaml
dev:
  - git
  - zsh
  - iterm2
  - visual-studio-code
```

**new-laptop** - Essential applications for new machines:

```yaml
new-laptop:
  - slack
  - google-chrome
  - 1password
```

### Custom Groups

Define your own groups in `~/.anvil/settings.yaml`:

```yaml
groups:
  frontend:
    - git
    - node
    - visual-studio-code
    - figma
  design:
    - figma
    - sketch
    - adobe-creative-cloud
  devops:
    - docker
    - kubectl
    - terraform
```

## How It Works

### Group Installation Process

1. **Reads configuration** - Loads groups from settings.yaml
2. **Validates tools** - Checks if tools are available in Homebrew
3. **Shows preview** - Lists what will be installed
4. **Requests confirmation** - Asks for user approval
5. **Installs tools** - Uses Homebrew to install each tool
6. **Reports results** - Shows success/failure status

### Individual App Installation Process

1. **Validates app name** - Checks if app exists in Homebrew
2. **Checks tracking** - Determines if app should be tracked
3. **Installs app** - Uses Homebrew to install
4. **Updates settings** - Adds to `tools.installed_apps` if applicable
5. **Reports status** - Shows installation result

### Smart Tracking Logic

Apps are automatically tracked in `tools.installed_apps` when installed individually, UNLESS they are already present in:

- `tools.required_tools`
- `tools.optional_tools`
- Any group in `groups`

This prevents duplication and keeps your settings clean.

## Examples

### Setting Up Development Environment

```bash
# Initialize Anvil first
anvil init

# Install development tools
anvil install dev

# Add individual tools as needed
anvil install docker
anvil install postman

# Add tools to specific groups
anvil install terraform --group-name devops
anvil install kubernetes-cli --group-name devops
```

### New Machine Setup

```bash
# Install essential applications
anvil install new-laptop

# Add personal productivity apps
anvil install notion
anvil install spotify
anvil install discord

# Organize apps into logical groups
anvil install firefox --group-name browsers
anvil install chrome --group-name browsers
anvil install safari --group-name browsers
```

### Custom Workflow

```bash
# Preview what would be installed
anvil install --list

# Install custom group
anvil install frontend

# Add specific tools
anvil install figma --dry-run
anvil install figma

# Create and populate custom groups
anvil install sketch --group-name design
anvil install adobe-creative-cloud --group-name design
anvil install figma --group-name design
```

### Team Setup

```bash
# Install team's standard development setup
anvil install dev

# Add team-specific tools
anvil install slack
anvil install zoom
anvil install jira
```

## Configuration

The install command reads from your `~/.anvil/settings.yaml`:

```yaml
tools:
  required_tools: [git, curl, brew]
  optional_tools: [docker, kubectl]
  installed_apps: [figma, notion, spotify] # Auto-populated

groups:
  dev: [git, zsh, iterm2, visual-studio-code]
  new-laptop: [slack, google-chrome, 1password]

  # Your custom groups
  frontend: [git, node, visual-studio-code, figma]
  design: [figma, sketch, adobe-creative-cloud]
```

## Under-the-Hood: Intelligent App Detection

Anvil uses a unified installation architecture that ensures consistent behavior across all installation modes (individual, group serial, and group concurrent). The system employs a hybrid approach for maximum reliability:

### Detection Process

1. **Homebrew Package Check** - Quick check if app is managed by Homebrew
2. **Installed Cask Detection** - Verify if cask is already installed via `brew list --cask`
3. **Dynamic Cask Search** - Use `brew search --cask` to find actual cask names and install locations
4. **Intelligent /Applications Search** - Smart name transformations (e.g., "visual-studio-code" → "Visual Studio Code.app")
5. **System-wide Spotlight Search** - macOS `mdfind` fallback for comprehensive detection
6. **PATH-based Detection** - Command-line tools via `which` command

### Key Benefits

- **No Hardcoded Mappings** - Dynamically detects any macOS application without manual configuration
- **Consistent Dry-run** - Preview mode performs identical availability checks as real installation
- **Manual Install Detection** - Recognizes apps installed outside of Homebrew
- **Future-proof** - Works with new applications without code updates

This architecture ensures that whether you run `anvil install postman`, `anvil install coding`, or `anvil install coding --dry-run`, you get consistent and accurate detection of already-installed applications.

## Troubleshooting

### Common Issues

**App Not Found**

```bash
# Search for correct app name
brew search firefox
brew search visual-studio

# Use exact Homebrew name
anvil install --dry-run visual-studio-code
```

**Permission Denied**

```bash
# Make sure Homebrew has proper permissions
brew doctor

# Check if you need admin privileges
sudo anvil install [app-name]
```

**Network Issues**

```bash
# Update Homebrew
brew update

# Check internet connection
brew search git
```

**Group Not Found**

```bash
# List available groups
anvil install --list

# Check settings.yaml syntax
anvil config show
```

### Getting Help

For more examples and detailed guides, see:

- [Getting Started Guide](GETTING_STARTED.md)
- [Examples & Tutorials](EXAMPLES.md)
- [Configuration Management](config.md)
