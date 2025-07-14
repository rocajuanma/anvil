# Getting Started with Anvil CLI

Welcome to Anvil! This guide will help you get up and running quickly with Anvil CLI, from installation to your first successful tool setup.

## Table of Contents

- [What is Anvil?](#what-is-anvil)
- [Installation](#installation)
- [First Steps](#first-steps)
- [Basic Usage](#basic-usage)
- [Common Workflows](#common-workflows)
- [Understanding Configuration](#understanding-configuration)
- [Tips and Best Practices](#tips-and-best-practices)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

## What is Anvil?

Anvil is a CLI automation tool that helps developers:

- **🚀 Bootstrap development environments** quickly and reliably
- **📦 Install tools in logical groups** (development, new laptop, custom)
- **⚙️ Manage configurations** across different machines
- **🔧 Automate repetitive setup tasks** for individuals and teams

### Key Concepts

- **Commands**: Actions you can perform (`init`, `setup`, `pull`, `push`, `draw`)
- **Groups**: Collections of related tools (`dev`, `new-laptop`, custom groups)
- **Configuration**: Settings stored in `~/.anvil/settings.yaml`
- **Tools**: Individual applications or utilities that can be installed

## Installation

### Quick Install

```bash
# Clone and build
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go build -o anvil main.go

# Move to PATH (optional)
sudo mv anvil /usr/local/bin/
```

For detailed installation instructions, see our [Installation Guide](INSTALLATION.md).

## First Steps

### 1. Initialize Anvil

Before using any other commands, initialize Anvil:

```bash
anvil init
```

This command will:

- ✅ Validate and install required tools (Git, cURL, Homebrew on macOS)
- ✅ Create necessary directories (`~/.anvil/`, `~/.anvil/cache/`, `~/.anvil/data/`)
- ✅ Generate default `settings.yaml` configuration
- ✅ Check your environment and provide recommendations

**Expected output:**

```
=== Anvil Initialization ===

🔧 Validating and installing required tools...
✅ All required tools are available
🔧 Creating necessary directories...
✅ Directories created successfully
🔧 Generating default settings.yaml...
✅ Default settings.yaml generated
🔧 Checking local environment configurations...
✅ Environment configurations are properly set

=== Initialization Complete! ===
```

### 2. Explore Available Commands

Get familiar with Anvil's capabilities:

```bash
# See all available commands
anvil --help

# Get help for specific commands
anvil init --help
anvil setup --help
```

### 3. Check Your Configuration

View the generated configuration:

```bash
cat ~/.anvil/settings.yaml
```

You'll see something like:

```yaml
version: 1.0.0
directories:
  config: /Users/username/.anvil
  cache: /Users/username/.anvil/cache
  data: /Users/username/.anvil/data
tools:
  required_tools: [git, curl]
  optional_tools: [brew, docker, kubectl]
groups:
  dev: [git, zsh, iterm2, vscode]
  new-laptop: [slack, chrome, 1password]
  custom: {}
git:
  username: "Your Name"
  email: "your.email@example.com"
environment: {}
```

## Basic Usage

### Installing Tool Groups

Anvil organizes tools into logical groups for easy batch installation:

#### Development Tools

Install essential development tools:

```bash
anvil setup dev
```

This installs:

- **Git** - Version control
- **Zsh** - Advanced shell with oh-my-zsh
- **iTerm2** - Enhanced terminal (macOS)
- **VS Code** - Code editor

#### New Laptop Essentials

Set up a new machine with essential applications:

```bash
anvil setup new-laptop
```

This installs:

- **Slack** - Team communication
- **Chrome** - Web browser
- **1Password** - Password manager

#### Preview Before Installing

Use dry-run to see what would be installed:

```bash
anvil setup dev --dry-run
anvil setup new-laptop --dry-run
```

### Installing Individual Tools

Install specific tools using flags:

```bash
# Install just Git
anvil setup --git

# Install multiple specific tools
anvil setup --git --zsh --vscode

# Preview individual tool installation
anvil setup --git --zsh --dry-run
```

### Available Individual Tools

| Flag          | Tool      | Description                   |
| ------------- | --------- | ----------------------------- |
| `--git`       | Git       | Version control system        |
| `--zsh`       | Zsh       | Advanced shell with oh-my-zsh |
| `--iterm2`    | iTerm2    | Enhanced terminal (macOS)     |
| `--vscode`    | VS Code   | Code editor                   |
| `--slack`     | Slack     | Team communication            |
| `--chrome`    | Chrome    | Web browser                   |
| `--1password` | 1Password | Password manager              |

### Listing Available Options

See all available groups and tools:

```bash
anvil setup --list
```

## Common Workflows

### Workflow 1: New Developer Machine Setup

Complete setup for a new development machine:

```bash
# Step 1: Initialize Anvil
anvil init

# Step 2: Install development tools
anvil setup dev

# Step 3: Add essential applications
anvil setup new-laptop

# Step 4: Add any additional tools as needed
# (Additional tools can be installed through custom groups)
```

### Workflow 2: Team Onboarding

Quickly onboard a new team member:

```bash
# Initialize
anvil init

# Install team-standard tools
anvil setup dev

# Add team communication tools
anvil setup --slack

# Additional tools can be defined in custom groups
# See configuration section for custom group setup
```

### Workflow 3: Selective Tool Installation

Install only specific tools you need:

```bash
# Initialize first
anvil init

# Preview what you want
anvil setup --git --zsh --vscode --dry-run

# Install selected tools
anvil setup --git --zsh --vscode
```

### Workflow 4: Custom Group Creation

Create your own tool groups by editing `~/.anvil/settings.yaml`:

```yaml
groups:
  custom:
    frontend:
      - git
      - node
      - yarn
      - chrome
    backend:
      - git
      - docker
      - postgresql
      - redis
    data-science:
      - python
      - jupyter
      - pandas
```

Then use your custom groups:

```bash
anvil setup frontend
anvil setup backend
```

## Understanding Configuration

### Configuration File Location

Anvil stores its configuration in `~/.anvil/settings.yaml`. This file contains:

- **Directories**: Paths for config, cache, and data
- **Tools**: Lists of required and optional tools
- **Groups**: Tool collections for batch installation
- **Git**: Your Git configuration
- **Environment**: Custom environment variables

### Customizing Groups

Edit `~/.anvil/settings.yaml` to add custom groups:

```yaml
groups:
  custom:
    my-workflow:
      - git
      - docker
      - vscode
      - slack
```

### Environment Configuration

Add custom environment variables:

```yaml
environment:
  EDITOR: "code"
  DEVELOPER_MODE: "true"
  MY_CUSTOM_PATH: "/usr/local/custom/bin"
```

### Directory Structure

Anvil creates and uses these directories:

```
~/.anvil/
├── settings.yaml    # Main configuration
├── cache/          # Temporary files and downloads
└── data/           # Persistent data and logs
```

## Tips and Best Practices

### 🎯 Initialization Best Practices

1. **Always run `anvil init` first** on any new machine
2. **Review the output** and follow any recommendations
3. **Complete environment setup** before installing tools

### 🔧 Tool Installation Best Practices

1. **Use dry-run first** to preview installations:

   ```bash
   anvil setup dev --dry-run
   ```

2. **Start with groups**, then add individual tools:

   ```bash
   anvil setup dev
   anvil setup --docker --kubectl
   ```

3. **Check available options** before installing:
   ```bash
   anvil setup --list
   ```

### 📋 Configuration Best Practices

1. **Backup your configuration** before making changes:

   ```bash
   cp ~/.anvil/settings.yaml ~/.anvil/settings.yaml.backup
   ```

2. **Use descriptive names** for custom groups
3. **Keep groups focused** - don't make them too large
4. **Document custom groups** for team members

### 🚀 Team Usage Best Practices

1. **Share configurations** across team members
2. **Create team-specific groups** in settings.yaml
3. **Document your team's setup process**
4. **Use consistent tool versions** when possible

## Troubleshooting

### Common Issues and Solutions

#### Command Not Found

```bash
# Check if anvil is in PATH
which anvil

# If not found, add to PATH or move to /usr/local/bin/
export PATH=$PATH:/path/to/anvil
# or
sudo mv anvil /usr/local/bin/
```

#### Permission Errors

```bash
# Fix common permission issues
sudo chown -R $(whoami) ~/.anvil
chmod 755 ~/.anvil
```

#### Tool Installation Failures

```bash
# Update Homebrew (macOS)
brew update

# Check internet connectivity
ping -c 3 github.com

# Try individual installation to isolate issues
anvil setup --git --dry-run
anvil setup --git
```

#### Configuration Issues

```bash
# Reinitialize if configuration is corrupted
rm -rf ~/.anvil
anvil init
```

#### Homebrew Issues (macOS)

```bash
# Fix Homebrew permissions
sudo chown -R $(whoami) $(brew --prefix)/*

# Reinstall Homebrew if needed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Platform-Specific Issues

#### macOS

- **Security warnings**: Allow in System Preferences → Security & Privacy
- **Xcode tools missing**: Run `xcode-select --install`
- **Homebrew PATH issues**: Add `/opt/homebrew/bin` to PATH

#### Linux

- **Missing build tools**: Install `build-essential` (Ubuntu) or `Development Tools` (CentOS)
- **Permission issues**: Use `sudo` for system-wide installations
- **Package manager**: Some tools may not be available via package managers

#### Windows

- **Use WSL or Git Bash** for best experience
- **PowerShell execution policy**: Run `Set-ExecutionPolicy RemoteSigned`
- **Limited tool support**: Not all tools available on Windows

### Getting Help

If you encounter issues not covered here:

1. **Check the logs** in `~/.anvil/data/`
2. **Search existing issues** on [GitHub](https://github.com/rocajuanma/anvil/issues)
3. **Create a new issue** with:
   - Your platform and version
   - Command you ran
   - Complete error message
   - Output of `anvil --version`

## Next Steps

Now that you're familiar with the basics:

### Explore Advanced Features

- **[Setup Command Documentation](setup-readme.md)** - Deep dive into tool installation
- **[Init Command Documentation](init-readme.md)** - Understand initialization process
- **[Examples and Tutorials](EXAMPLES.md)** - Real-world usage scenarios

### Customize Your Setup

- **Edit `~/.anvil/settings.yaml`** to create custom tool groups
- **Add environment variables** for your workflow
- **Share configurations** with your team

### Contribute

- **[Contributing Guide](CONTRIBUTING.md)** - Help improve Anvil
- **[Development Setup](.local/anvil-rules.md)** - Development guidelines
- **Report bugs** or **request features** on GitHub

### Stay Updated

- **Watch the repository** for updates
- **Read the [changelog](CHANGELOG.md)** for new features
- **Join discussions** for community support

---

## Quick Reference

### Essential Commands

```bash
anvil init                    # Initialize Anvil
anvil setup --list          # List available groups and tools
anvil setup dev              # Install development tools
anvil setup --git --zsh     # Install specific tools
anvil setup dev --dry-run   # Preview installations
anvil --help                 # Get help
```

### Key Files

- `~/.anvil/settings.yaml` - Main configuration
- `~/.anvil/cache/` - Temporary files
- `~/.anvil/data/` - Persistent data

### Important Locations

- **Configuration**: `~/.anvil/`
- **Documentation**: `docs/` directory
- **Development**: `.local/anvil-rules.md`

---

**Ready to start?** Run `anvil init` and begin automating your development workflow!
