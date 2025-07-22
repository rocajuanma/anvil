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

> Cast aside the burden of manual setup—let Anvil be your Samwise, carrying your configs to the very end.

**Anvil** is the complete macOS development environment automation tool. Stop manually setting up machines, hunting for configs, and dealing with inconsistent environments. With Anvil, you get zero-config tool installation, cross-machine configuration sync, and team-wide environment standardization—all in one powerful CLI.

## ✨ What Anvil Solves

### 🚀 **Installation & Tool Management**

- **Environment Inconsistency** → Smart tool installation with automatic tracking
- **Manual Setup Pain** → Group-based installation for common scenarios
- **Tool Discovery** → Install any macOS app with `anvil install [app-name]`
- **Setup Documentation** → Self-documenting configuration in `settings.yaml`

### ⚙️ **Configuration Management & Sync**

- **Config Drift** → Version-controlled dotfiles and settings via GitHub
- **Team Onboarding** → Shared configuration repositories for instant setup
- **Machine Migrations** → Cross-machine sync with `anvil config pull/push`
- **Configuration Loss** → Automated backup and recovery of development environments

## 🎯 Key Features

### 📦 **Smart Installation**

- **🎯 Dynamic Installation** - Install any macOS application with `anvil install [app-name]`
- **📝 Intelligent Tracking** - Apps automatically tracked in `tools.installed_apps`
- **📦 Group Management** - Predefined and custom tool groups for common scenarios
- **🍺 Homebrew Integration** - Automatic installation with intelligent app verification (detects Homebrew, manual, and system installations)

### 🔄 **Configuration Sync**

- **🌐 Cross-Machine Sync** - Keep configs consistent across all your development environments
- **👥 Team Collaboration** - Share configurations via GitHub repositories
- **🔍 Smart Change Detection** - Pre-push diff analysis prevents unnecessary commits
- **📁 Directory Organization** - App-specific config management (cursor, vscode, zsh)

### 🛠 **Developer Experience**

- **🚀 Zero Configuration** - Works out of the box with sensible defaults
- **🔍 Dry-run Support** - Preview installations and changes before execution
- **🎨 Beautiful Output** - Colored, structured progress indicators
- **⚡ Fast Operations** - Concurrent installation and smart caching
- **🩺 Health Checks** - Comprehensive environment validation with real-time progress feedback via `anvil doctor`

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

### Basic Workflows

#### **🔧 Tool Installation Workflow**

```bash
# Initialize Anvil (run this first!)
anvil init

# Verify setup is working correctly with real-time progress feedback
anvil doctor

# Install applications dynamically
anvil install firefox
anvil install visual-studio-code

# Install predefined tool groups
anvil install dev        # git, zsh, iterm2, visual-studio-code
anvil install new-laptop # slack, google-chrome, 1password

# Preview before installing
anvil install docker --dry-run
```

#### **⚙️ Configuration Management Workflow**

```bash
# Set up GitHub repository (one-time setup)
# Edit ~/.anvil/settings.yaml with your repo details

# Verify connectivity with detailed progress feedback
anvil doctor connectivity

# Pull configurations from your repository
anvil config pull cursor
anvil config pull vscode

# View pulled configurations
anvil config show cursor

# Sync configuration state with system
anvil config sync        # Install missing apps from settings

# Push local changes back to repository
anvil config push       # Push anvil settings to GitHub
```

## 📦 Installation Methods

Install individual applications or predefined groups with automatic Homebrew integration.

### Individual Applications

Install any Homebrew package by name with automatic tracking:

```bash
anvil install terraform
anvil install kubernetes-cli
anvil install postman
anvil install obsidian
```

**🎯 Smart Features:**

- Apps automatically tracked in `tools.installed_apps`
- Duplicate prevention (won't track apps already in groups)
- Manual installation detection (works with pre-installed apps)

### Predefined Groups

- **`dev`** - Essential development tools (git, zsh, iterm2, visual-studio-code)
- **`new-laptop`** - Essential applications for new machines (slack, google-chrome, 1password)

### Custom Groups

Define your own in `~/.anvil/settings.yaml`:

```yaml
groups:
  backend:
    - git
    - docker
    - postgresql
    - redis
  frontend:
    - git
    - node
    - visual-studio-code
    - figma
  devops:
    - docker
    - kubernetes-cli
    - terraform
    - vault
```

📖 **[Complete Installation Guide](docs/install.md)** - Detailed installation options, troubleshooting, and examples

## 🔧 Configuration Management

Sync dotfiles and configurations across machines using GitHub repositories with full version control.

### Cross-Machine Synchronization

```bash
# Pull configurations from your repository
anvil config pull neovim
anvil config pull tmux
anvil config pull zsh

# View configurations before applying
anvil config show neovim

# Sync missing apps from your settings
anvil config sync --dry-run # Preview changes
anvil config sync           # Apply changes

# Push local changes to repository
anvil config push           # Creates timestamped branch with PR link
```

### Team Configuration Sharing

```bash
# Pull team's development setup by specialty
anvil config pull team-backend
anvil config pull team-frontend
anvil config pull team-devops

# Install team's recommended tools
anvil config sync team-backend

# View team configurations
anvil config show team-backend
```

### Key Configuration Features

- **🔍 Smart Change Detection** - Only pushes when configurations actually differ
- **🌿 Timestamped Branches** - Creates branches like `config-push-18072025-1234`
- **🔗 PR-Ready Workflow** - Provides direct GitHub PR links
- **📁 Organized Storage** - Directory-based config organization in repositories
- **🔐 Multiple Auth Methods** - SSH keys, GitHub tokens, or public access
- **⚡ Efficient Operations** - Local caching and smart diff algorithms

📖 **[Complete Configuration Guide](docs/config.md)** - Setup, authentication, repository structure, and team workflows

## ⚙️ Settings

Your development environment configuration in `~/.anvil/settings.yaml`:

```yaml
tools:
  required_tools: [git, curl, brew]
  optional_tools: [docker, kubectl]
  installed_apps: [terraform, postman, kubernetes-cli, obsidian] # Auto-tracked individual installs
groups:
  dev: [git, zsh, iterm2, visual-studio-code]
  backend: [git, docker, postgresql, redis] # Custom groups
git:
  username: "Your Name"
  email: "your.email@example.com"
github:
  config_repo: "username/dotfiles" # For config sync
  branch: "main"
  token_env_var: "GITHUB_TOKEN"
```

## 🎯 Command Reference

| Command             | Description            | Example                     |
| ------------------- | ---------------------- | --------------------------- |
| `init`              | Initialize environment | `anvil init`                |
| `install [app]`     | Install application    | `anvil install terraform`   |
| `install [group]`   | Install tool group     | `anvil install dev`         |
| `install --list`    | List available groups  | `anvil install --list`      |
| `config pull [app]` | Pull configurations    | `anvil config pull neovim`  |
| `config show [app]` | Show configurations    | `anvil config show neovim`  |
| `config sync [app]` | Sync configurations    | `anvil config sync`         |
| `config push [app]` | Push configurations    | `anvil config push`         |
| `doctor`            | Run health checks      | `anvil doctor`              |
| `doctor [category]` | Check specific area    | `anvil doctor dependencies` |
| `doctor [check]`    | Run individual check   | `anvil doctor git-config`   |
| `doctor --fix`      | Auto-fix issues        | `anvil doctor --fix`        |

### Useful Flags

- `--dry-run` - Preview installations and changes
- `--list` - Show available groups and tracked apps
- `--concurrent` - Enable parallel installation (faster)
- `--verbose` - Show detailed output (doctor command)
- `--fix` - Automatically fix detected issues (doctor command)

## 📚 Documentation

### **Documentation & Guides**

| Guide                                          | Description                              |
| ---------------------------------------------- | ---------------------------------------- |
| **[Getting Started](docs/GETTING_STARTED.md)** | Complete setup and first workflows       |
| **[Installation Guide](docs/INSTALLATION.md)** | Platform-specific installation           |
| **[Install Command](docs/install.md)**         | Deep-dive on tool installation           |
| **[Configuration Management](docs/config.md)** | Complete config sync setup and workflows |
| **[Doctor Command](docs/doctor.md)**           | Health checks and environment validation |
| **[Examples & Tutorials](docs/EXAMPLES.md)**   | Real-world usage scenarios               |
| **[Contributing](docs/CONTRIBUTING.md)**       | Development guidelines                   |
| **[Changelog](docs/CHANGELOG.md)**             | Version history and updates              |

## 🍺 macOS Focus

Optimized specifically for macOS with:

- **Homebrew Integration** - Automatic installation and cask support
- **Native Terminal Colors** - Beautiful output in Terminal.app and iTerm2
- **GUI Application Support** - Seamless installation of Mac applications
- **Application Detection** - Smart detection of manually installed apps

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Homebrew](https://brew.sh/) - macOS package management

---

<div align="center">

**[⬆ Back to Top](#anvil)**

Made with ❤️ for macOS engineers who value consistency and automation

</div>
