/*
Copyright © 2022 Juanma Roca juanmaxroca@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package constants

// Long descriptions for commands
const ANVIL_LONG_DESCRIPTION = `Anvil is a powerful macOS automation CLI tool designed to streamline development workflows 
and personal tool configuration. It provides a comprehensive suite of commands for managing
development environments, automating installations, and maintaining consistent configurations.

What Anvil Solves:
• 🎯 Effortlessly install and manage any macOS application or CLI tool using Homebrew
• 📝 Automatically track all individually installed apps and tools in your settings.yaml for easy reproducibility
• 🗂️ Organize and install tools with flexible group management for different workflows or machine setups
• ⚡ Zero configuration required: works out of the box with smart, sensible defaults
• 🍺 Seamless Homebrew integration for automated installation, upgrades, and uninstalls
• 🔄 Manage, sync, and version your configuration files and dotfiles with simple commands
• 🔍 Dry-run mode to preview all changes before they happen
• 🌈 Beautiful, structured, and colored output for clear progress and results`

const INIT_COMMAND_LONG_DESCRIPTION = `The init command bootstraps your Anvil CLI environment by performing a complete
initialization process. This is the first command you should run after installing Anvil.

What it does:
• ✅ Validates and installs required system tools (Git, cURL, Homebrew)
• 📁 Creates necessary configuration directory (~/.anvil)
• ⚙️ Generates a default settings.yaml configuration file with your system preferences
• 🔍 Checks your local environment for common development configurations
• 💡 Provides actionable recommendations for completing your setup
• 🎨 Displays beautiful ASCII banner for visual confirmation

This command is designed specifically for macOS and requires Homebrew for tool management.`

const INSTALL_COMMAND_LONG_DESCRIPTION = `The install command provides dynamic installation of development tools and applications
for macOS using Homebrew. It supports both group-based and individual installations with intelligent detection.

Installation Modes:
• anvil install [group-name]    - Install all tools in a predefined group
• anvil install [app-name]      - Install any individual application via brew

Available Groups: 
• dev - Essential development tools
• new-laptop - Essential applications for new machines
• Custom groups you define in settings.yaml

Key Features:
• 📝 Automatic App Tracking: Every app you install individually is automatically recorded in settings.yaml under tools.installed_apps for easy environment reproduction
• 🔍 Intelligent App Detection: Uses unified hybrid approach (Homebrew check → cask search → /Applications scan → Spotlight search → PATH detection) to verify app availability regardless of installation method
• 🎯 Manual Install Recognition: Detects apps installed outside Homebrew (manually downloaded, Mac App Store, etc.) preventing unnecessary reinstallation attempts
• 🚦 Consistent Dry-Run: Preview mode performs identical availability checks as real installation for accurate previews
• 🗂️ Group Management: Install tool collections with single commands or define custom groups in settings.yaml
• ⚡ Concurrent Installation: Use --concurrent flag for parallel installation with significant speed improvements
• 🧠 Smart Deduplication: Apps already in groups or required_tools are not redundantly tracked in installed_apps

Flags: Use --list to see available groups, --dry-run to preview, --concurrent for faster parallel installation.`

const CONFIG_COMMAND_LONG_DESCRIPTION = `The config command provides centralized management of configuration files and dotfiles
for your development environment. It serves as a parent command for configuration-related operations.

Subcommands:
• anvil config pull [directory]    - Pull configuration files from remote repository
• anvil config push [directory]    - Push configuration files to remote repository  
• anvil config show [directory]    - Show configuration files from anvil settings or pulled directories
• anvil config sync [directory]    - Sync configuration state with system reality

Key Features:
• 📁 Directory-specific operations for granular configuration management
• 🔄 Version-controlled dotfiles and settings via GitHub repositories
• 🛡️ Automated backup and recovery of development environments
• 👥 Team configuration sharing and collaboration
• 🔍 Smart change detection with pre-push diff analysis
• ⚡ Cross-machine synchronization for consistent development environments

GitHub Repository Configuration:
The 'github.config_repo' field in settings.yaml should be in the format 'username/repository'.

This command structure ensures all configuration operations are properly organized with clear
separation between configuration management and other system operations.`

const PUSH_COMMAND_LONG_DESCRIPTION = `The push command enables you to upload and synchronize your local configuration files
to GitHub for backup and sharing with automated branch creation and change tracking.

Features:
• 🔍 Smart Change Detection: Compares local and remote configurations before proceeding to avoid unnecessary commits
• 🌿 Timestamped Branches: Creates branches with format 'config-push-DDMMYYYY-HHMM' for organized version control
• 📁 Organized Storage: Commits anvil settings to '/anvil' directory in repository for clear structure  
• 💬 Standardized Commits: Uses consistent commit messages for easy tracking and identification
• 🔗 PR-Ready Workflow: Provides direct GitHub links to create pull requests after successful push
• ⚙️ Automated Git Operations: Handles repository cloning, branch creation, committing, and pushing automatically

Implementation Status:
• ✅ Option 1: Anvil settings push (anvil config push) - Fully functional
• ✅ Option 2: Application config push (anvil config push <app-name>) - Fully functional

Perfect for maintaining consistent development environments and sharing configurations across teams.`

const PULL_COMMAND_LONG_DESCRIPTION = `The pull command allows you to download and synchronize configuration files
from a specific directory in your GitHub repository to your local machine.

Usage: anvil config pull [directory]

How it works:
• 📥 Automatically fetches the latest changes from your repository (git fetch/pull)  
• 📁 Copies all files from the specified directory to a temporary location (~/.anvil/temp/[directory])
• ✅ Guarantees you get the most up-to-date configurations every time you pull
• 🔄 Supports multiple repository formats with automatic URL validation and correction
• 🛡️ Secure authentication via SSH keys, GitHub tokens, or public repository access
• 📋 Clear progress feedback with detailed status information

Perfect for:
• Setting up new development environments quickly and consistently
• Synchronizing specific configurations across multiple machines  
• Restoring configurations after system changes or updates
• Sharing configurations with team members and collaborators

GitHub Repository Setup:
Configure your repository in ~/.anvil/settings.yaml with format 'username/repository'.
Supports various URL formats including SSH, HTTPS, and domain-prefixed formats.`

const SHOW_COMMAND_LONG_DESCRIPTION = `The show command displays configuration files and settings for easy viewing and inspection
with intelligent formatting based on content type and structure.

Usage Modes:
• anvil config show              - Display the main anvil settings.yaml file with syntax highlighting
• anvil config show [directory]  - Show configuration files from a pulled directory

Features:
• 📄 Single File Display: Shows file content directly in terminal with proper formatting
• 📁 Multiple Files: Shows tree structure with comprehensive file listings and organization
• ✅ Smart Content Detection: Automatically determines best display method based on file type and count
• 🎨 Syntax Highlighting: Provides clear visual formatting for YAML, JSON, and other configuration formats  
• 💡 Helpful Error Messages: Clear guidance with suggestions for missing directories or invalid paths
• 🔍 Detailed File Information: Shows file sizes, modification dates, and directory structures

Perfect for reviewing pulled configurations before applying them, checking current anvil settings,
and understanding repository structure and organization.`

const SYNC_COMMAND_LONG_DESCRIPTION = `The sync command moves pulled configuration files from the temp directory to their local destinations
with automatic archiving and safety confirmations.

Usage Modes:
• anvil config sync              - Sync anvil settings.yaml file from pulled configs
• anvil config sync [app-name]   - Sync specific app configurations to their local paths

Features:
• 📋 Safe Configuration Override: Archives existing configs before applying new ones
• ✅ Interactive Confirmation: Always asks permission before overriding local files  
• 🔍 Comprehensive Dry-Run: Preview all changes without applying them using --dry-run flag
• 📦 Automatic Archiving: Backs up existing configs to timestamped archive directories
• 🎯 Smart Path Resolution: Uses configs section in settings.yaml for app-specific destinations
• 💡 Clear Error Messages: Helpful guidance when pulled configs or paths are missing

Perfect for applying pulled configurations safely, maintaining config consistency across machines,
and ensuring you always have recovery options through the automatic archiving system.`

// Doctor command descriptions
const DOCTOR_COMMAND_LONG_DESCRIPTION = `Run comprehensive health checks to validate your anvil environment with real-time progress feedback.

The doctor command performs validation across four key areas with live progress indicators,
so you always know what's happening. You can run checks at different levels of granularity:

CATEGORIES (groups of related checks):
• environment    - Verify anvil initialization and directory structure  
• dependencies   - Check required tools and Homebrew installation
• configuration  - Validate git and GitHub settings
• connectivity   - Test GitHub access and repository connections

SPECIFIC CHECKS (individual validators):
Run 'anvil doctor --list' to see all 12 available individual checks.

KEY FEATURES:
✨ Real-time progress indicators with counters (e.g., [1/12] 8% - Running init-run)
🔍 Two output modes: brief default output and detailed verbose mode
🔒 Secure non-interactive authentication (no credential prompts)
🎨 Professional user experience with consistent formatting

Examples:
  anvil doctor                    # Run all health checks with progress feedback
  anvil doctor --list             # Show available categories and checks
  anvil doctor environment        # Run all environment checks (3 checks)
  anvil doctor dependencies       # Run all dependency checks (3 checks)
  anvil doctor git-config         # Run only the git configuration check
  anvil doctor homebrew           # Run only the Homebrew check
  anvil doctor --fix              # Auto-fix detected issues
  anvil doctor dependencies --fix # Auto-fix dependency issues
  anvil doctor --verbose          # Show detailed output with step-by-step results

The doctor provides actionable recommendations for any issues found, shows real-time
progress so you never wait in silence, and can automatically fix many common problems.`

// Clean command descriptions
const CLEAN_COMMAND_LONG_DESCRIPTION = `The clean command removes all content inside .anvil directories while preserving
the essential settings.yaml configuration file.

What it does:
• 🧹 Removes all temporary files, archives, and downloaded configurations
• 📁 Cleans up contents inside temp/ and archive/ directories while preserving their structure
• 🗂️ Completely removes dotfiles/ directory to ensure clean git repository state
• ✅ Preserves your main settings.yaml file with all your configurations
• 🔍 Shows exactly what will be deleted before proceeding
• 🛡️ Requires confirmation unless --force flag is used
• 🧪 Supports --dry-run mode to preview changes without deletion

Use cases:
• Free up disk space by removing old temporary files and archives
• Clean up after testing or development work
• Reset to a clean state while keeping your settings
• Remove downloaded configurations that are no longer needed

Safety features:
• Always preserves settings.yaml to maintain your configuration
• Preserves essential directory structure (temp/, archive/) for tool functionality
• Completely removes dotfiles/ directory to ensure clean git repository state
• Interactive confirmation prevents accidental deletions
• Dry-run mode shows exactly what would be deleted
• Force flag available for automated scripts

This command is safe and will never delete your main configuration file. The dotfiles directory will be completely removed to ensure the next pull/push operation works seamlessly.`
