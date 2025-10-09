# Changelog

All notable changes to Anvil CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed
- **Package Architecture Refactor** 🏗️ - Migrated codebase from `pkg/` to `internal/` package structure
  - **Breaking Change**: Improved encapsulation boundaries using Go's internal package visibility
  - Enhanced maintainability and refactoring safety for internal components

- **Terminal Output Package Refactor** 🔄 - Replaced local terminal package with external Palantir dependency
  - Removed `internal/terminal` package (296 lines deleted) in favor of `github.com/rocajuanma/palantir v1.0.0`
  - Updated all commands and packages to use `palantir.GetGlobalOutputHandler()`
  - Improved consistency and maintainability of terminal output across the codebase
  - **Breaking Change**: Applications directly using `internal/terminal` will need to migrate to the `palantir` package

### Fixed
- **Test Suite Compatibility** 🔧 - Fixed GitHub test case for updated `CopyFileSimple` behavior
  - Updated test to align with new directory creation behavior in filesystem utilities

## [1.5.1] - 2025-10-04

### Added

### Changed

### Fixed
- **Installation Script Directory Creation** 🔧 - Fixed installation failure on fresh macOS systems
  - Installation script now creates `/usr/local/bin` directory if it doesn't exist
  - Resolves "No such file or directory" error when installing on systems without Homebrew or other package managers
  - **Issue**: Installation failed on fresh Mac systems where `/usr/local/bin` directory didn't exist

## [1.5.0] - 2025-10-01

### Added
- **Enhanced Repository Badges** 🏆 - Enhance badge collection to README for better project visibility
  - Go Report Card badge for code quality assessment

### Changed
- **Command Descriptions Simplified** 📝 - Streamlined all command descriptions for better readability
  - Reduced description length while preserving essential usage information
  - Removed verbose feature lists for cleaner, more scannable descriptions
  - **Benefit**: Improved user experience with concise, focused command documentation

- **Documentation Improvements** 📚 - Enhanced getting started guide and examples
  - Updated configuration examples to reference repository files instead of user-specific paths
  - Improved import workflow examples with real GitHub URLs
  - Streamlined authentication setup instructions
  - Enhanced backup and configuration management guidance
  - **Benefit**: Clearer onboarding experience and more actionable documentation

### Fixed

## [1.4.1] - 2025-09-25

### Added

### Changed

### Fixed
- **Code Complexity Reduction** 🔧 - Refactored high-complexity functions & fix goreportcard issues to improve maintainability
  - Reduced cyclomatic complexity on push, clean and output logic while maintaining identical functionality
  - **Benefit**: Improved code organization, testability, and maintainability

## [1.4.0] - 2025-09-18

### Added
- **Installation Performance Optimizations** ⚡ - Major speed improvements across the installation pipeline
  - Eliminated redundant availability checks by checking tool availability only once per installation
  - Optimized cask detection with 3-tier lookup (static table → runtime cache → dynamic search) reducing network calls by 90%+
  - Added batch configuration loading for group installations to eliminate repeated file I/O
  - Reordered application availability checks from fastest to slowest operations for early returns
  - Cached Homebrew installation status to eliminate redundant system calls
  - **Result**: Near-instant detection for installed tools and 3-5x faster group installations

### Changed

### Fixed

## [1.3.3] - 2025-09-16

### Added

### Changed
- **Terminal Output Formatting** 🎨 - Improved error message clarity and fixed redundant output formatting
  - Eliminated cascading "failed to install" error wrapping for cleaner messages
  - Fixed duplicate emoji display in error outputs
  - Enhanced progress display formatting to prevent line concatenation issues

### Fixed
- **Homebrew Installation Reliability** 🔧 - Fixed PATH detection and prerequisites validation
  - Added post-installation verification and enhanced PATH checking
  - Improved error reporting with actual installation script output
  - Added Xcode Command Line Tools validation before installation
  - **Issue**: PATH refresh problems and missing prerequisites caused installation failures

- **Duplicate Homebrew Installation** 🔧 - Fixed `anvil init` attempting to install Homebrew twice causing exit status 1 failures
  - Removed Homebrew from regular tools validation loop, now handled only as prerequisite  
  - **Issue**: Users without Homebrew pre-installed experienced initialization failures due to duplicate installation attempts

## [1.3.2] - 2025-09-15

### Added

### Changed
- **Filesystem Utilities Consolidation** 🔧 - Consolidated duplicate file operations across codebase
  - Unified multiple `copyFile` and `copyDirRecursive` implementations into single utilities
  - Added configurable options for copy behavior while maintaining backward compatibility
  - **Benefit**: Follows DRY principle and provides consistent filesystem operations

### Fixed
- **Pull Command Argument Consistency** 🔧 - Made directory argument optional for `anvil config pull` command
  - Now matches `anvil config show` and `anvil config push` commands which accept optional arguments
  - When no directory is specified, defaults to pulling from "anvil" directory
  - `anvil config pull` is now equivalent to `anvil config pull anvil`
  - Updated command help text to reflect the optional argument behavior
  - **Issue**: Previously required a directory argument while other config commands made it optional

## [1.3.1] - 2025-09-14

### Added

### Changed

### Fixed
- **Brew Cask Detection** 🔧 - Fixed incorrect cask detection causing installation failures
  - Fix string matching that treated error messages as valid cask names
  - Add success checks and error message filtering
  - Improve error reporting to show actual brew output
  - Remove redundant formula search logic
  - **Issue**: Resolves installation failures for formula packages incorrectly detected as casks

- **Push Command Independence** 🔧 - Fixed bug where cancelled push commands left staged changes affecting subsequent pushes
  - Ensures each push operation starts from clean repository state
  - Automatically cleans up staged changes when push is cancelled or fails
  - Prevents stale changes from previous operations appearing in new push commits
  - **Issue**: Running `anvil config push app-one`, cancelling, then running `anvil config push app-two` would include changes from both apps

## [1.3.0] - 2025-09-10

### Added
- **Group Import from Files/URLs** 📥 - New `anvil config import` command for importing tool groups
  - **Flexible Sources** - Import from local files or remote URLs
  - **Security-First Design** - Extracts only group definitions, ignores sensitive data
  - **Comprehensive Validation** - Validates group names, structure, and detects conflicts
  - **Interactive Confirmation** - Shows preview and requires user approval before import
  - **Usage**: `anvil config import ./groups.yaml`, `anvil config import https://example.com/groups.yaml`

### Changed

### Fixed

## [1.2.0] - 2025-09-09

### Fixed
- **Homebrew Update Loop** 🔧 - Fixed infinite warning cycle in `anvil doctor --fix`
  - **Conservative Approach** - Only updates formulae database, doesn't auto-upgrade packages
  - **Enhanced User Experience** - Detailed outdated package listing with numbered format
  - **Manual Upgrade Instructions** - Clear guidance to respect user control over package versions
  - **Safety Explanation** - Explains why auto-upgrades are avoided to prevent compatibility issues
  - **Issue Resolution** - Resolves cycle where doctor detected updates but fix didn't resolve them
### Added
- **Update Command** 🔄 - New `anvil update` command for seamless version updates
  - **One-Command Updates** - Update to the latest Anvil version with `anvil update`
  - **Safe Update Process** - Uses the same trusted installation script as initial setup
  - **Dry-Run Preview** - Preview updates with `--dry-run` flag without making changes
  - **macOS Optimized** - Specifically designed for macOS environments
  - **Comprehensive Documentation** - Complete usage guide with troubleshooting
  - **Usage**: `anvil update`, `anvil update --dry-run`

## [1.1.2] - 2025-09-03

### Added
- **Group Assignment for Individual App Installation** 📦 - New `--group-name` flag for the `anvil install` command
  - **Seamless Group Organization** - Install apps and automatically add them to existing or new groups
  - **Dynamic Group Creation** - Creates new groups automatically if they don't exist
  - **Duplicate Prevention** - Prevents adding the same app multiple times to a group
  - **Graceful Fallback** - Falls back to normal `installed_apps` tracking if group operations fail
  - **Success-Only Operations** - Group operations only occur if installation is successful
  - **Usage Examples**: `anvil install firefox --group-name essentials`, `anvil install final-cut --group-name editing`
- **Real-time Progress for Doctor Command** ✨ - Enhanced `anvil doctor` command with live feedback
  - **Live Progress Indicators** - See validation progress as each check runs with progress counters (e.g., `[1/12] 8% - Running init-run`)
  - **Stage-by-stage Feedback** - Clear indication of which validation is currently executing
  - **Enhanced Verbose Mode** - `--verbose` flag now provides detailed descriptions, categories, and step-by-step results
  - **No More Hanging Terminals** - Always know what's happening during validation
  - **Two Output Modes** - Brief default output for quick checks, detailed verbose output for debugging
- **Secure Non-interactive Authentication** 🔒 - Eliminated all credential prompts for enhanced security
  - **No Credential Prompts** - All git operations are now completely non-interactive
  - **Environment-based Authentication** - Uses `GITHUB_TOKEN` environment variable or SSH keys exclusively
  - **Non-interactive Git Operations** - Set `GIT_TERMINAL_PROMPT=0`, `GIT_ASKPASS=/bin/false`, and `SSH_ASKPASS=/bin/false`
  - **SSH BatchMode** - All SSH operations use `BatchMode=yes` for security
  - **Enhanced GitHub Validation** - Detailed verbose output showing authentication method attempts
- **Standardized Output Formatting** 🎨 - Consistent progress formatting across all commands
  - **Config Commands Progress** - `anvil config pull`, `push`, `show` now have stage-by-stage progress indicators
  - **Unified Terminal Output** - All commands use consistent `PrintStage()`, `PrintSuccess()`, `PrintWarning()` formatting
  - **Clear Stage Indicators** - Users understand what phase each command is in
  - **Professional User Experience** - Same visual style and feedback patterns across the entire CLI
- **Enhanced GitHub Security Validation** 🛡️ - Improved repository security checks
  - **Private Repository Enforcement** - Clear warnings and blocks for public repositories
  - **Authentication Method Details** - Verbose mode shows which authentication method is being used
  - **Repository Privacy Verification** - Validates that configuration repositories are private
  - **Security-first Design** - All authentication methods prioritize security over convenience
- **Comprehensive Health Check System** - New `anvil doctor` command for environment validation and troubleshooting
  - **Multi-Category Validation** - Systematic checks across environment, dependencies, configuration, and connectivity
  - **12 Built-in Validators** - Complete coverage from initialization to GitHub repository access
  - **Auto-Fix Capabilities** - Automatic resolution of common issues (Homebrew updates, directory permissions, missing dependencies)
  - **Granular Execution** - Run all checks, specific categories, or individual validators
  - **Interactive Fix Mode** - User-confirmed automatic fixes with verification
  - **Detailed Reporting** - Structured output with actionable fix recommendations
  - **Extensible Architecture** - Scalable validator framework following DRY principles
- **Simplified Groups Structure** - Removed "custom:" header requirement, all groups now at the same level in settings.yaml
- **Automatic App Tracking** - Individual apps installed via `anvil install [app-name]` are automatically tracked in `tools.installed_apps`
- **Smart Duplicate Prevention** - Prevents tracking apps already in groups, required_tools, or optional_tools
- **Dynamic Individual Installation** - Install any Homebrew package by name with automatic tracking
- **Clean Command** 🧹 - New `anvil clean` command for environment maintenance and cleanup
  - **Smart Cleanup** - Removes temporary files, archives, and downloaded configurations while preserving settings.yaml
  - **Directory Structure Preservation** - Maintains essential directories (temp/, archive/) for tool functionality
  - **Complete Git Cleanup** - Fully removes dotfiles/ directory to ensure clean repository state
  - **Safety Features** - Interactive confirmation, dry-run support, and settings preservation
  - **Intelligent Handling** - Different cleanup strategies for different directory types
  - **Usage Examples**: `anvil clean`, `anvil clean --dry-run`, `anvil clean --force`
- **Configuration Pull System** - Complete implementation of `anvil config pull [directory]` command
- **Configuration Show System** - New `anvil config show [directory]` command to view pulled configs and settings
- **Configuration Sync System** - New `anvil config sync [directory]` command to reconcile settings with system state
  - Directory-specific configuration pulling from GitHub repositories
  - SSH key and GitHub token authentication support
- **Configuration Push System** - New `anvil config push` command to upload anvil settings to GitHub repository
  - Smart difference detection between local and remote configurations
  - Timestamped branch creation with format `config-push-DDMMYYYY-HHMM`
  - Automated commit with standardized message `anvil[push]: anvil`
  - Direct pull request link generation for workflow integration
  - Pre-push validation to avoid unnecessary Git operations when configurations are up-to-date
  - Automatic GitHub URL format validation and normalization
  - Branch validation with detailed error messages listing available branches
  - Temporary file storage at `~/.anvil/temp/[directory]` for review before manual application
  - Enhanced error handling with step-by-step troubleshooting guidance
- **GitHub Integration Package** - New `internal/github` package for repository operations
  - Repository cloning and pulling with branch support
  - Multiple authentication methods (SSH keys, GitHub tokens, public access)
  - Branch existence validation and available branch listing
  - Git user configuration management
- **Comprehensive Documentation** - Complete configuration management documentation
  - New `docs/config.md` with detailed setup instructions and examples
  - Updated `README.md`, `GETTING_STARTED.md`, and `EXAMPLES.md` with current implementation
  - Real-world examples for personal and team configuration sharing
- Parallel tool installation support
- Tool version management in groups
- Dependency resolution for tools
- Configuration validation command
- Backup and restore functionality
- Windows package manager support
- Linux distribution-specific package managers
### Changed
- **BREAKING CHANGE**: Command renamed from `setup` to `install` - All `anvil setup` commands are now `anvil install`
- **BREAKING CHANGE**: Individual app installation approach - Replaced flag-based installation (`--git`, `--zsh`) with dynamic name-based installation (`anvil install git`, `anvil install zsh`)
- **BREAKING CHANGE**: Groups configuration structure - Removed "custom:" header requirement, all groups now at the same level in settings.yaml
- **BREAKING CHANGE**: Reorganized `pull` and `push` commands as subcommands of `config`
  - `anvil pull` → `anvil config pull [directory]`
  - `anvil push` → `anvil config push [directory]` (now implemented for anvil settings)
  - This change improves command hierarchy and follows Cobra best practices
  - All documentation updated to reflect new command structure with directory-specific operations
- **Enhanced Configuration Structure** - Updated settings.yaml with GitHub and Git configuration sections
  - Added `github.config_repo`, `github.branch`, `github.local_path`, `github.token_env_var`
  - Added `git.ssh_key_path`, `git.ssh_dir` for SSH configuration
  - Automatic GitHub URL format validation and correction
- **Improved Error Messages** - Branch configuration errors now provide detailed guidance
  - Lists available branches when specified branch doesn't exist
  - Step-by-step troubleshooting instructions with settings file path
  - Enhanced validation messages during repository access
- Enhanced progress indicators
- Better platform detection
### Fixed
- **Codebase Cleanup** - Removed unused code and simplified directory structure
  - Removed unused `.anvil/cache` and `.anvil/data` directories from entire codebase
  - Simplified `AnvilDirectories` struct to only include `Config` field
  - Removed concurrent installation features, unused error constructors, and extended package constants
  - Cleaned up placeholder implementations and unused flags
- Memory usage optimization for large tool sets
- Updated import paths for reorganized command structure

## [1.0.1] - 2024-01-XX

### Added
- **Configuration Caching** - Thread-safe configuration caching with double-checked locking pattern
- **Context-Aware Commands** - Command execution with configurable timeouts to prevent hanging
- **Structured Error Handling** - `AnvilError` struct with operation, command context, and error types
- **Constants Extraction** - Centralized constants for system commands, package names, and environment variables
### Changed
- **Improved Error Handling** - Commands now return proper errors instead of calling `os.Exit(1)` directly
- **Unified Installation Patterns** - Consolidated installation logic using `InstallConfig` struct
- **Eliminated Global Variables** - Replaced global flags with struct-based approach in setup command
- **Enhanced Performance** - Configuration loading optimized with caching, eliminating repeated file I/O
### Fixed
- **Resource Management** - Commands now use timeouts to prevent indefinite hanging
- **Code Duplication** - Eliminated ~50 lines of duplicate installation code
- **Magic Strings** - Extracted ~30 magic strings to constants for better maintainability
- **Thread Safety** - Configuration access is now thread-safe with proper locking

## [1.0.0] - 2024-01-XX

### Added
- **Init Command** - Complete environment bootstrapping
  - Tool validation and installation (Git, cURL, Homebrew)
  - Directory structure creation (`~/.anvil/`)
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

