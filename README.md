<div align="center">
  <img src="assets/anvil-2.0.png" alt="Anvil Logo" width="200" style="border-radius: 50%;">
  <h1>Anvil CLI</h1>
</div>

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.17+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rocajuanma/anvil)](https://goreportcard.com/report/github.com/rocajuanma/anvil)
[![GitHub Release](https://img.shields.io/github/v/release/rocajuanma/anvil?style=flat&label=Release)](https://github.com/rocajuanma/anvil/releases/latest)
[![Platform](https://img.shields.io/badge/platform-macOS%20only-blue.svg)](#installation)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

</div>

 Tired of wasting time setting up your Mac for development? Anvil automates tool installs, syncs configs, and keeps your environment consistent—no hassle, no manual steps, just one powerful CLI.

<div align="center">
  <img src="assets/anvil.gif" alt="Anvil Demo" width="600">
</div>

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

# Import tool groups from shared configs
anvil config import https://example.com/team-groups.yaml

# Or start with example configurations
anvil config import https://raw.githubusercontent.com/rocajuanma/anvil/master/import-examples/juanma-essentials.yaml

# Check environment health
anvil doctor

# Sync configurations (after setting up GitHub repo)
anvil config pull neovim
anvil config sync neovim
```

## Key Features

- **🎯 Smart Installation** - Install individual apps or predefined groups (`dev`, `new-laptop`)
- **📝 Auto-tracking** - Automatically tracks installed apps and prevents duplicates
- **📥 Group Import** - Import groups from local files or URLs with validation and conflict detection
- **🔒 Secure Config Sync** - Uses private GitHub repositories with automatic backups
- **🩺 Health Diagnostics** - `anvil doctor` detects and auto-fixes common issues
- **🧹 Environment Cleanup** - Smart cleanup tools that preserve essential configs
- **🚀 Zero Configuration** - Works out of the box with sensible defaults

## Documentation

| Guide | Description |
|-------|-------------|
| **[Getting Started](docs/GETTING_STARTED.md)** | Complete setup and workflows |
| **[Configuration Management](docs/config.md)** | Config sync setup and workflows |
| **[Import Groups](docs/import.md)** | Import tool groups from files/URLs |
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
