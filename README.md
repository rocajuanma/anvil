<div align="center">
  <img src="assets/anvil-logo.png" alt="Anvil Logo" width="200" style="border-radius: 50%;">
  <h1>Anvil</h1>
</div>

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.17+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20only-blue.svg)](#macos-focus)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)
[![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)](docs/CHANGELOG.md)

</div>

> One CLI to rule them all, forged in the fires of productivity: wield Anvil to command your macOS realm and master your development destiny.

Supercharge your macOS development setup with Anvil — the all-in-one CLI for effortless app installation, smart tracking, group management, and seamless Homebrew integration. Instantly set up tools, preview changes, and enjoy beautiful, actionable feedback with zero configuration required.

## ✨ Features

- **🎯 Dynamic Installation** - Install any macOS application with `anvil install [app-name]`
- **📝 Smart Tracking** - Individual apps automatically tracked in `tools.installed_apps`
- **📦 Group Management** - Predefined and custom tool groups for common scenarios
- **🚀 Zero Configuration** - Works out of the box with sensible defaults
- **🍺 Homebrew Integration** - Automatic installation and management
- **🔍 Dry-run Support** - Preview installations before execution
- **🎨 Beautiful Output** - Colored, structured progress indicators

## 🚀 Quick Start

### Installation

```bash
# Clone and build
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go build -o anvil main.go

# Move to your PATH (optional)
sudo mv anvil /usr/local/bin/
```

### Basic Usage

```bash
# Initialize Anvil (run this first!)
anvil init

# Install applications dynamically
anvil install firefox
anvil install visual-studio-code

# Install predefined tool groups
anvil install dev        # git, zsh, iterm2, visual-studio-code
anvil install new-laptop # slack, google-chrome, 1password

# Preview before installing
anvil install docker --dry-run
```

## 📋 Installation Methods

### Individual Applications

Install any Homebrew package by name with automatic tracking:

```bash
anvil install firefox
anvil install slack
anvil install figma
```

### Predefined Groups

- **`dev`** - Essential development tools
- **`new-laptop`** - Essential applications for new machines

### Custom Groups

Define your own in `~/.anvil/settings.yaml`:

```yaml
groups:
  frontend:
    - git
    - node
    - visual-studio-code
    - figma
```

## 🔧 Configuration Management

Sync dotfiles and configurations across machines using GitHub repositories:

```bash
# Pull configurations from your repository
anvil config pull cursor
anvil config pull vs-code

# View configurations
anvil config show        # View anvil settings
anvil config show cursor # View pulled configs

# Sync configuration state with system
anvil config sync        # Install missing apps from settings
anvil config sync --dry-run # Preview what would be synced

# Push configurations to repository
anvil config push        # Push anvil settings
anvil config push cursor # Push app configs (in development)
```

📖 **[Complete Setup Guide](docs/config.md)** - Authentication, repository structure, and examples

## ⚙️ Settings

Basic settings structure in `~/.anvil/settings.yaml`:

```yaml
tools:
  required_tools: [git, curl, brew]
  optional_tools: [docker, kubectl]
  installed_apps: [figma, notion, spotify] # Auto-tracked
groups:
  dev: [git, zsh, iterm2, visual-studio-code]
  frontend: [git, node, visual-studio-code, figma] # Custom groups
git:
  username: "Your Name"
  email: "your.email@example.com"
github:
  config_repo: "username/dotfiles" # For config sync
```

## 🎯 Command Reference

| Command             | Description            | Example                    |
| ------------------- | ---------------------- | -------------------------- |
| `init`              | Initialize environment | `anvil init`               |
| `install [app]`     | Install application    | `anvil install firefox`    |
| `install [group]`   | Install tool group     | `anvil install dev`        |
| `install --list`    | List available groups  | `anvil install --list`     |
| `config pull [app]` | Pull configurations    | `anvil config pull cursor` |
| `config show [app]` | Show configurations    | `anvil config show cursor` |
| `config sync [app]` | Sync configurations    | `anvil config sync`        |
| `config push [app]` | Push configurations    | `anvil config push`        |

### Useful Flags

- `--dry-run` - Preview installations
- `--list` - Show available groups and tracked apps

## 📚 Documentation

| Guide                                          | Description                    |
| ---------------------------------------------- | ------------------------------ |
| **[Getting Started](docs/GETTING_STARTED.md)** | Complete setup and usage guide |
| **[Configuration Management](docs/config.md)** | Dotfiles sync with GitHub      |
| **[Examples & Tutorials](docs/EXAMPLES.md)**   | Real-world usage scenarios     |
| **[Install Command](docs/install.md)**         | Detailed installation options  |
| **[Contributing](docs/CONTRIBUTING.md)**       | Development guidelines         |
| **[Changelog](docs/CHANGELOG.md)**             | Version history and updates    |

## 🍺 macOS Focus

Optimized specifically for macOS with Homebrew integration, native terminal colors, and automatic GUI application support.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Homebrew](https://brew.sh/) - macOS package management

---

<div align="center">

**[⬆ Back to Top](#anvil)**

Made with ❤️ for macOS engineers

</div>
