# Changelog

All notable changes to Anvil CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Parallel tool installation support
- Tool version management in groups
- Dependency resolution for tools
- Configuration validation command
- Backup and restore functionality
- Windows package manager support
- Linux distribution-specific package managers

### Changed

- Improved error messages with more context
- Enhanced progress indicators
- Better platform detection

### Fixed

- Memory usage optimization for large tool sets

## [1.0.0] - 2024-01-XX

### Added

- **Init Command** - Complete environment bootstrapping

  - Tool validation and installation (Git, cURL, Homebrew)
  - Directory structure creation (`~/.anvil/`, cache, data)
  - Default configuration generation (`settings.yaml`)
  - Environment configuration checking
  - SSH key and Git configuration validation
  - Comprehensive error handling and user guidance

- **Setup Command** - Group-based and individual tool installation

  - Predefined groups: `dev` (git, zsh, iterm2, vscode), `new-laptop` (slack, chrome, 1password)
  - Individual tool flags: `--git`, `--zsh`, `--iterm2`, `--vscode`, `--slack`, `--chrome`, `--1password`
  - Dry-run capability with `--dry-run` flag
  - List available groups and tools with `--list` flag
  - Custom group support in configuration
  - Progress indicators for installation status
  - Graceful error handling with continuation

- **Configuration Management**

  - YAML-based configuration in `~/.anvil/settings.yaml`
  - Support for custom tool groups
  - Git configuration integration
  - Environment variable management
  - Directory customization support

- **Terminal Output System**

  - Colored, structured output with progress indicators
  - Consistent formatting across all commands
  - Stage-based progress reporting
  - Success, warning, and error message types
  - User guidance and next steps

- **Cross-Platform Support**

  - Full macOS support with Homebrew integration
  - Basic Linux support for command-line tools
  - Limited Windows support with appropriate warnings
  - Platform-specific tool installation logic

- **Package Management Integration**

  - Homebrew integration for macOS
  - Automatic Homebrew installation
  - Package installation verification
  - Error handling for missing package managers

- **Tool Management**

  - Zsh installation with oh-my-zsh configuration
  - Git installation and configuration
  - VS Code, iTerm2, Slack, Chrome, 1Password support
  - Tool availability validation
  - Graceful handling of already-installed tools

- **Documentation**
  - Comprehensive README with quick start guide
  - Detailed command documentation (init, setup)
  - Installation guide for all platforms
  - Getting started tutorial
  - Examples and use cases
  - Contributing guidelines
  - Development rules and standards

### Technical Features

- **Modular Architecture** - Clean separation of concerns
- **Error Resilience** - Comprehensive error handling
- **Idempotent Operations** - Safe to run multiple times
- **Extensible Design** - Easy to add new tools and platforms

### Developer Experience

- **Clear Command Structure** - Intuitive command hierarchy
- **Helpful Error Messages** - Actionable guidance for users
- **Progress Feedback** - Real-time installation progress
- **Comprehensive Help** - Detailed help for all commands

## [0.1.0] - Initial Release

### Added

- Basic project structure
- Simple CLI framework
- Initial command scaffolding

---

## Release Notes

### Version 1.0.0 - Major Release

This is the first major release of Anvil CLI, representing a complete development environment automation solution. The release includes:

**🚀 Core Features**

- Complete environment initialization with the `init` command
- Powerful group-based tool installation with the `setup` command
- Comprehensive configuration management
- Beautiful terminal output with progress indicators

**🛠️ Tool Support**

- Essential development tools (Git, Zsh, VS Code, iTerm2)
- Team collaboration tools (Slack, Chrome, 1Password)
- Custom tool group definitions
- Individual tool installation flags

**📚 Documentation**

- Complete user guides and tutorials
- API documentation for developers
- Platform-specific installation instructions
- Contributing guidelines and development standards

**🌍 Platform Coverage**

- Full macOS support with Homebrew integration
- Linux support for command-line tools
- Windows compatibility with WSL recommendations

### Breaking Changes

None - this is the initial major release.

### Migration Guide

This is the first major release, so no migration is needed. Follow the [Getting Started Guide](GETTING_STARTED.md) for initial setup.

### Deprecations

None in this release.

---

## Contributing to Changelog

When contributing changes, please:

1. **Add entries to [Unreleased]** section
2. **Use appropriate categories**: Added, Changed, Deprecated, Removed, Fixed, Security
3. **Write clear descriptions** of changes
4. **Include relevant links** to issues or PRs
5. **Follow semantic versioning** principles

### Categories

- **Added** - New features
- **Changed** - Changes in existing functionality
- **Deprecated** - Soon-to-be removed features
- **Removed** - Now removed features
- **Fixed** - Any bug fixes
- **Security** - Vulnerability fixes

### Example Entry Format

```markdown
### Added

- New feature description (#123)
- Another feature with more details (#456)

### Changed

- Modified behavior description (#789)

### Fixed

- Bug fix description (#101)
```

---

## Support

- **Documentation**: [docs/](.)
- **Issues**: [GitHub Issues](https://github.com/rocajuanma/anvil/issues)
- **Discussions**: [GitHub Discussions](https://github.com/rocajuanma/anvil/discussions)
