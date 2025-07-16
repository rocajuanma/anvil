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

> A powerful macOS automation CLI tool for dynamic development environment setup

Anvil provides a streamlined approach to managing development environments on macOS through intelligent Homebrew integration. Install any application or manage tool groups with a single command - Anvil dynamically determines what you want to install and handles it intelligently.

## ✨ Features

- **🎯 Dynamic Installation** - Install any macOS application with `anvil setup [app-name]`
- **📦 Smart Group Management** - Predefined tool groups for common development scenarios
- **🚀 Zero Configuration** - Works out of the box with sensible defaults
- **🍺 Homebrew Integration** - Automatic Homebrew installation and management
- **⚙️ Intelligent Error Handling** - Helpful error messages with actionable suggestions
- **🔍 Dry-run Support** - Preview installations before execution
- **🎨 Beautiful Terminal Output** - Colored, structured progress indicators

## 🔧 Configuration Management

Anvil provides powerful configuration management through the `config` command, allowing you to sync configuration files and dotfiles across machines using GitHub repositories:

### Pull Configurations (✅ Available)

```bash
# Pull specific configuration directories from your GitHub repository
anvil config pull cursor      # Pull Cursor editor configurations
anvil config pull vs-code     # Pull VS Code configurations
anvil config pull zsh         # Pull shell configurations
anvil config pull git         # Pull Git configurations
```

**Current Implementation**: Always fetches the latest changes from your repository and pulls configurations to `~/.anvil/temp/[directory]` for review before manual application.

**Future Enhancement**: Will automatically apply configurations to their destination directories.

### Push Configurations (🚧 In Development)

```bash
# Upload local configurations to your GitHub repository (coming soon)
anvil config push cursor      # Push Cursor configurations
anvil config push --all       # Push all configuration changes
```

### Key Benefits

- **🔄 Cross-machine sync** - Keep configurations consistent across all your devices
- **📁 Directory-specific pulls** - Pull only the configurations you need
- **🛡️ Version control** - Your configurations are safely stored in GitHub
- **👥 Team sharing** - Share team configurations and best practices
- **🔍 Branch validation** - Automatic branch checking with detailed error messages
- **🔐 Multiple auth methods** - Support for SSH keys and GitHub tokens

### Getting Started with Config Management

1. **Set up your GitHub repository** with organized configuration directories
2. **Configure Anvil** by editing `~/.anvil/settings.yaml`
3. **Start pulling configurations** with `anvil config pull [directory]`

📖 **[Complete Configuration Guide](docs/config-readme.md)** - Detailed setup instructions, examples, and troubleshooting

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

### First Steps

```bash
# Initialize Anvil (run this first!)
anvil init

# Install any application dynamically
anvil setup firefox
anvil setup slack
anvil setup visual-studio-code

# Install predefined tool groups
anvil setup dev
anvil setup new-laptop

# Preview before installing
anvil setup docker --dry-run

# Manage configurations (if you have a config repository)
anvil config pull cursor    # Pull Cursor editor configs
anvil config pull vs-code   # Pull VS Code configs
anvil config push cursor    # Push local changes (coming soon)
```

## 📋 Dynamic Setup Command

### Individual Application Installation

The setup command intelligently installs any macOS application available through Homebrew:

```bash
# Install any application by name
anvil setup firefox
anvil setup slack
anvil setup docker
anvil setup visual-studio-code
anvil setup figma
anvil setup notion
anvil setup chrome
anvil setup spotify
anvil setup zoom

# Preview installation
anvil setup [app-name] --dry-run

# Get helpful error messages for typos
anvil setup firefx  # Suggests using 'brew search firefx'
```

### Predefined Tool Groups

Install curated sets of tools for common scenarios:

#### Development Group (`dev`)

Essential tools for software development:

- **git** - Version control system
- **zsh** - Advanced shell with oh-my-zsh
- **iterm2** - Enhanced terminal
- **visual-studio-code** - Code editor

```bash
anvil setup dev
```

#### New Laptop Group (`new-laptop`)

Essential applications for a new machine:

- **slack** - Team communication
- **google-chrome** - Web browser
- **1password** - Password manager

```bash
anvil setup new-laptop
```

#### Custom Groups

Define your own groups in `~/.anvil/settings.yaml`:

```yaml
groups:
  custom:
    frontend:
      - git
      - node
      - visual-studio-code
      - figma
    content:
      - notion
      - obsidian
      - figma
      - canva
```

Then use them:

```bash
anvil setup frontend
anvil setup content
```

## 💻 Usage Examples

### Complete Development Setup

```bash
# Initialize Anvil
anvil init

# Install all development tools
anvil setup dev

# Add additional tools dynamically
anvil setup docker
anvil setup figma
anvil setup spotify
```

### Team Onboarding

```bash
# Quick team member setup
anvil init
anvil setup dev
anvil setup slack
anvil setup zoom
```

### Selective Installation

```bash
# Install only what you need
anvil setup git --dry-run  # Preview first
anvil setup git            # Install
anvil setup firefox
anvil setup notion
```

## ⚙️ Configuration

Anvil stores configuration in `~/.anvil/settings.yaml`:

```yaml
version: 1.0.0
directories:
  config: /Users/username/.anvil
tools:
  required_tools: [git, curl, brew]
  optional_tools: [docker, kubectl]
groups:
  dev: [git, zsh, iterm2, visual-studio-code]
  new-laptop: [slack, google-chrome, 1password]
  custom:
    frontend: [git, node, visual-studio-code, figma]
git:
  username: "Your Name"
  email: "your.email@example.com"
environment: {}
```

## 🍺 macOS Focus

Anvil is optimized specifically for macOS and leverages:

- **Homebrew** - Primary package manager for all installations
- **macOS System Integration** - Native terminal colors and progress indicators
- **Cask Support** - Automatic GUI application installation
- **Shell Integration** - oh-my-zsh setup and configuration

### Why macOS Only?

- **Consistent Package Management** - Homebrew provides reliable, consistent installations
- **Quality Assurance** - Focus on one platform allows for better testing and reliability
- **Native Integration** - Optimal terminal experience and system integration
- **Developer-Focused** - macOS is the primary platform for many development workflows

## 🎯 Command Reference

### Core Commands

| Command             | Description                                | Example                    |
| ------------------- | ------------------------------------------ | -------------------------- |
| `init`              | Initialize Anvil environment               | `anvil init`               |
| `setup [app]`       | Install any application                    | `anvil setup firefox`      |
| `setup [group]`     | Install tool group                         | `anvil setup dev`          |
| `config pull [dir]` | Pull specific config directory from remote | `anvil config pull cursor` |
| `config push [dir]` | Push config directory to remote (dev)      | `anvil config push cursor` |
| `draw [font]`       | Generate ASCII art text                    | `anvil draw banner`        |

### Setup Flags

| Flag        | Description           | Example                        |
| ----------- | --------------------- | ------------------------------ |
| `--dry-run` | Preview installation  | `anvil setup docker --dry-run` |
| `--list`    | List available groups | `anvil setup --list`           |
| `--update`  | Update Homebrew first | `anvil setup --update`         |

### Available Groups

- `dev` - Development tools (git, zsh, iterm2, visual-studio-code)
- `new-laptop` - Essential apps (slack, google-chrome, 1password)
- Custom groups defined in your configuration

## 🏗️ Architecture

```
anvil/
├── cmd/                    # Command implementations
│   ├── initcmd/           # Environment initialization
│   ├── setup/             # Dynamic app installation
│   ├── draw/              # ASCII art generation
│   ├── config/            # Configuration management
│   │   ├── pull/          # Pull configuration files (config pull)
│   │   └── push/          # Push configuration files (config push)
│   └── root.go            # CLI framework setup
├── pkg/                   # Core packages
│   ├── brew/              # Homebrew integration
│   ├── config/            # Configuration management
│   ├── constants/         # Application constants
│   ├── system/            # System command execution
│   ├── terminal/          # Terminal output formatting
│   ├── tools/             # Tool validation
│   └── figure/            # ASCII art generation
├── docs/                  # Documentation
└── main.go                # Application entry point
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

### Quick Development Setup

```bash
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go mod download
go build -o anvil main.go
./anvil init
```

## 🐛 Troubleshooting

### Common Issues

**Application not found:**

```bash
# Use brew search to find the correct name
brew search firefox
anvil setup firefox
```

**Homebrew issues:**

```bash
# Reinstall Homebrew
anvil init  # Will reinstall if needed
```

**Permission errors:**

```bash
# Fix Homebrew permissions
sudo chown -R $(whoami) $(brew --prefix)/*
```

## 📚 Documentation

- **[Getting Started Guide](docs/GETTING_STARTED.md)** - Comprehensive setup guide
- **[Configuration Management](docs/config-readme.md)** - Complete guide to config pull/push functionality
- **[Examples & Tutorials](docs/EXAMPLES.md)** - Real-world usage scenarios
- **[Contributing Guide](docs/CONTRIBUTING.md)** - Development guidelines
- **[Changelog](docs/CHANGELOG.md)** - Version history

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [go-figure](https://github.com/common-nighthawk/go-figure) - ASCII art generation
- [Homebrew](https://brew.sh/) - macOS package management

---

<div align="center">

**[⬆ Back to Top](#anvil-cli)**

Made with ❤️ for macOS developers

</div>
