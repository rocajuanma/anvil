# Anvil CLI

[![Go Version](https://img.shields.io/badge/go-1.17+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey.svg)](#platform-support)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)
[![Version](https://img.shields.io/badge/version-1.0.1-blue.svg)](docs/CHANGELOG.md)

> A powerful automation CLI tool designed to streamline development workflows and personal tool configuration

Anvil provides a comprehensive suite of commands for managing development environments, automating installations, and maintaining consistent configurations across different systems. Whether you're setting up a new development machine or maintaining tool consistency across a team, Anvil makes it simple and reliable.

## ✨ Features

- **🚀 Automated Environment Setup** - Bootstrap development environments with a single command
- **📦 Group-based Tool Installation** - Install sets of tools organized by purpose (dev, new-laptop, custom)
- **🔧 Individual Tool Management** - Install and configure specific tools with dedicated flags
- **⚙️ Configuration Management** - Centralized configuration with sensible defaults
- **🌍 Cross-platform Support** - Works on macOS, Linux, and Windows (optimized for macOS)
- **📋 Dry-run Capabilities** - Preview changes before execution
- **🎨 Beautiful Terminal Output** - Colored, structured output with progress indicators
- **📚 Comprehensive Documentation** - Detailed guides and examples for every feature

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/rocajuanma/anvil.git
cd anvil

# Build the binary
go build -o anvil main.go

# Move to your PATH (optional)
sudo mv anvil /usr/local/bin/
```

### First Steps

```bash
# Initialize Anvil (run this first!)
anvil init

# See all available commands
anvil --help

# Install development tools
anvil setup dev

# Install individual tools
anvil setup --git --zsh --vscode
```

## 📋 Commands Overview

| Command                         | Description                                         | Example               |
| ------------------------------- | --------------------------------------------------- | --------------------- |
| [`init`](docs/init-readme.md)   | Bootstrap and initialize Anvil environment          | `anvil init`          |
| [`setup`](docs/setup-readme.md) | Install development tools in groups or individually | `anvil setup dev`     |
| `pull`                          | Download assets and configurations from GitHub      | `anvil pull configs`  |
| `push`                          | Upload assets and configurations to GitHub          | `anvil push dotfiles` |
| `draw`                          | Generate ASCII art text for terminal output         | `anvil draw "Hello"`  |

## 🏗️ Tool Groups

Anvil organizes tools into logical groups for easy batch installation:

### Development Group (`dev`)

Essential tools for software development:

- **Git** - Version control system
- **Zsh** - Advanced shell with oh-my-zsh
- **iTerm2** - Enhanced terminal (macOS)
- **VS Code** - Popular code editor

```bash
anvil setup dev
```

### New Laptop Group (`new-laptop`)

Essential applications for setting up a new machine:

- **Slack** - Team communication
- **Chrome** - Web browser
- **1Password** - Password manager

```bash
anvil setup new-laptop
```

### Custom Groups

Define your own tool groups in `~/.anvil/settings.yaml`:

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
```

## 💻 Individual Tool Installation

Install specific tools with dedicated flags:

```bash
# Install individual tools
anvil setup --git --zsh --vscode

# Preview what would be installed
anvil setup --git --zsh --dry-run

# List all available tools and groups
anvil setup --list
```

## ⚙️ Configuration

Anvil stores its configuration in `~/.anvil/settings.yaml`:

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

## 🎯 Usage Examples

### Complete Development Setup

```bash
# Initialize Anvil
anvil init

# Install all development tools
anvil setup dev

# Add communication tools
anvil setup --slack --chrome

# Generate a project banner
anvil draw "My Project"
```

### Team Onboarding

```bash
# Quick team member setup
anvil init
anvil setup dev
anvil setup new-laptop

# Custom team-specific tools
anvil setup --docker --kubectl --slack
```

### Preview Changes

```bash
# See what would be installed without actually installing
anvil setup dev --dry-run
anvil setup --git --zsh --dry-run
```

## 🖥️ Platform Support

| Platform    | Support Level | Notes                                   |
| ----------- | ------------- | --------------------------------------- |
| **macOS**   | ✅ Full       | Optimized with Homebrew integration     |
| **Linux**   | ⚠️ Partial    | Basic tools supported, limited GUI apps |
| **Windows** | ⚠️ Limited    | Command-line tools only                 |

## 📚 Documentation

- **[Getting Started Guide](docs/GETTING_STARTED.md)** - Comprehensive setup and usage guide
- **[Installation Guide](docs/INSTALLATION.md)** - Detailed installation instructions
- **[Init Command](docs/init-readme.md)** - Complete init command documentation
- **[Setup Command](docs/setup-readme.md)** - Complete setup command documentation
- **[Examples & Tutorials](docs/EXAMPLES.md)** - Real-world usage examples
- **[Contributing](docs/CONTRIBUTING.md)** - How to contribute to Anvil
- **[Changelog](docs/CHANGELOG.md)** - Version history and updates

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details on:

- Setting up the development environment
- Code style and standards
- Submitting pull requests
- Reporting issues

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

**Permission errors:**

```bash
# Fix Homebrew permissions (macOS)
sudo chown -R $(whoami) /usr/local/Homebrew
```

**Tool installation failures:**

```bash
# Update Homebrew and try again
brew update
anvil setup --git --dry-run  # Preview first
```

**Configuration issues:**

```bash
# Reinitialize Anvil
rm -rf ~/.anvil
anvil init
```

See our [troubleshooting guide](docs/GETTING_STARTED.md#troubleshooting) for more solutions.

## 🏗️ Architecture

Anvil is built with a modular architecture:

```
anvil/
├── cmd/           # Command implementations
│   ├── initcmd/   # Init command
│   ├── setup/     # Setup command
│   ├── draw/      # Draw command
│   ├── pull/      # Pull command
│   ├── push/      # Push command
│   └── root.go    # Root command configuration
├── pkg/           # Reusable packages
│   ├── brew/      # Homebrew integration
│   ├── config/    # Configuration management
│   ├── constants/ # Application constants and error types
│   ├── figure/    # ASCII art generation
│   ├── system/    # System command execution
│   ├── terminal/  # Terminal output formatting
│   └── tools/     # Tool validation and installation
├── docs/          # Documentation
└── main.go        # Application entry point
```

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [go-figure](https://github.com/common-nighthawk/go-figure) - ASCII art generation
- [Homebrew](https://brew.sh/) - Package management for macOS

## 🔗 Links

- **Repository**: [github.com/rocajuanma/anvil](https://github.com/rocajuanma/anvil)
- **Issues**: [Report a bug or request a feature](https://github.com/rocajuanma/anvil/issues)
- **Discussions**: [Community discussions](https://github.com/rocajuanma/anvil/discussions)

---

<div align="center">

**[⬆ Back to Top](#anvil-cli)**

Made with ❤️ by [Juanma Roca](https://github.com/rocajuanma)

</div>
