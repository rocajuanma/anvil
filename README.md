<div align="center">
  <img src="assets/anvil-logo.png" alt="Anvil Logo" width="200" style="border-radius: 50%;">
  <h1>Anvil</h1>
</div>

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.17+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20only-blue.svg)](#installation)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)
[![Version](https://img.shields.io/badge/version-1.2.0+-blue.svg)](docs/CHANGELOG.md)

</div>

**Anvil** is the complete macOS development automation tool. Setting up and maintaining a consistent macOS dev env can be painful, error-prone and derail attention. Stop manually setting up machines, hunting for configs, and dealing with inconsistent environments. With Anvil, you get zero-config batch tool installation, cross-machine configuration sync, and team-wide environment standardization—all in one powerful CLI.

## What Anvil Does

- **🚀 Batch App Installation** - Install development tools in groups or individually via Homebrew
- **🔄 Configuration Sync** - Sync dotfiles across machines using private GitHub repositories  
- **🩺 Health Checks** - Auto-diagnose and fix common setup issues

## Why Choose Anvil?
- **⏱️ Fast, Automated Setup—Focus on Coding, Not Configuration** – Anvil handles all tool installations and configuration sync automatically, letting you get started in minutes instead of hours.
- **🧑‍💻 Effortless Onboarding & Consistency** – Onboard new machines or teammates with a single command, ensuring everyone has the same reliable, ready-to-code environment—every time, on every Mac.
- **🛡️ Built-in Safety** – Dry-run mode, automatic backups, and smart deduplication protect your system and your work.
- **👥 Seamless Team Collaboration** – Instantly sync dotfiles and configs from private GitHub repositories, making team onboarding and environment sharing simple and secure.

## Quick Start

### Installation

**New installations:**
```bash
curl -sSL https://github.com/rocajuanma/anvil/releases/latest/download/install.sh | bash
```

**Update existing installation:**
```bash
anvil update
```

> **Note**: The `anvil update` command was introduced in v1.2.0. If you have an older version, use the curl command above.

### Try It Out

```bash
# Initialize Anvil
anvil init

# Install development tools
anvil install dev        # git, zsh, iterm2, visual-studio-code
anvil install terraform  # Individual apps

# Check environment health
anvil doctor

# Sync configurations (after setting up GitHub repo)
anvil config pull neovim
anvil config sync neovim
```

## Key Features

- **🎯 Smart Installation** - Install individual apps or predefined groups (`dev`, `new-laptop`)
- **📝 Auto-tracking** - Automatically tracks installed apps and prevents duplicates
- **🔒 Secure Config Sync** - Uses private GitHub repositories with automatic backups
- **🩺 Health Diagnostics** - `anvil doctor` detects and auto-fixes common issues
- **🧹 Environment Cleanup** - Smart cleanup tools that preserve essential configs
- **🚀 Zero Configuration** - Works out of the box with sensible defaults

## Documentation

| Guide | Description |
|-------|-------------|
| **[Getting Started](docs/GETTING_STARTED.md)** | Complete setup and workflows |
| **[Configuration Management](docs/config.md)** | Config sync setup and workflows |
| **[Install Command](docs/install.md)** | Tool installation guide |
| **[Doctor Command](docs/doctor.md)** | Health checks and validation |
| **[Examples & Tutorials](docs/EXAMPLES.md)** | Real-world usage scenarios |

**[📖 View All Documentation →](docs/)**

---

<div align="center">

**[⬆ Back to Top](#anvil)**

Made with ❤️ for macOS developers who value automation and consistency

**[⭐ Star this project](https://github.com/rocajuanma/anvil)** • **[📖 Documentation](docs/)** • **[🐛 Report Issues](https://github.com/rocajuanma/anvil/issues)**

</div>
