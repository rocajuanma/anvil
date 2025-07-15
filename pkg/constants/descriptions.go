package constants

// Long descriptions for commands
const ANVIL_LONG_DESCRIPTION = `Anvil is a powerful macOS automation CLI tool designed to streamline development workflows 
and personal tool configuration. It provides a comprehensive suite of commands for managing
development environments, automating installations, and maintaining consistent configurations.

Key features:
- Automated tool installation via Homebrew
- Dynamic group and individual app installation
- Environment configuration management
- ASCII art generation for enhanced terminal output
- Optimized specifically for macOS`

const INIT_COMMAND_LONG_DESCRIPTION = `The init command bootstraps your Anvil CLI environment by performing a complete
initialization process. This is the first command you should run after installing Anvil.

What it does:
• Validates and installs required system tools (Git, cURL, Homebrew)
• Creates necessary configuration directories (~/.anvil, ~/.anvil/cache, ~/.anvil/data)
• Generates a default settings.yaml configuration file with your system preferences
• Checks your local environment for common development configurations
• Provides actionable recommendations for completing your setup

This command is designed specifically for macOS and requires Homebrew for tool management.`

const SETUP_COMMAND_LONG_DESCRIPTION = `The setup command provides dynamic installation of development tools and applications 
using Homebrew on macOS.

Usage:
• anvil setup [group-name]    - Install all tools in a predefined group
• anvil setup [app-name]      - Install any individual application via brew

This command intelligently determines if the argument is a group name (defined in settings.yaml)
or an application name. If it's not a group, it attempts to install the application directly
using Homebrew. All installations are validated and errors are handled gracefully.

Supported groups: dev, new-laptop, and custom groups defined in your configuration.

Use 'anvil setup --list' to see all available groups.`

const CONFIG_COMMAND_LONG_DESCRIPTION = `The config command provides centralized management of configuration files and assets
for your development environment. It serves as a parent command for pull and push operations
to ensure all configuration-related actions are properly organized and guarded.

Subcommands:
• anvil config pull    - Pull configuration files from remote repository
• anvil config push    - Push configuration files to remote repository

This command structure ensures that all pull and push operations are scoped to configuration
files only, providing clear separation between configuration management and other operations.
Use this command to maintain consistent configuration across different development environments.`

const PUSH_COMMAND_LONG_DESCRIPTION = `The push command enables you to upload and synchronize your local configuration files,
and dotfiles to GitHub for backup and sharing purposes.

Features:
• Selective config pushing based on configuration
• Automatic Git repository management
• Conflict resolution and merge handling
• Support for various config types (dotfiles, configs, scripts)

The command takes an argument to specify which type of configurations should be pushed,
allowing for granular control over what gets synchronized to your remote repository.
This command is now scoped to configuration files only.`

const PULL_COMMAND_LONG_DESCRIPTION = `The pull command allows you to download and synchronize configuration files,
and dotfiles from your GitHub repository to your local machine.

This is particularly useful for:
• Setting up new development environments
• Synchronizing configurations across multiple machines
• Restoring configurations after system changes
• Sharing configurations with team members

The command takes an argument to specify which type of configurations should be retrieved,
providing flexibility in what gets synchronized to your local environment.
This command is now scoped to configuration files only.`

const DRAW_COMMAND_LONG_DESCRIPTION = `The draw command generates beautiful ASCII art representations of text using the
go-figure library. This command enhances terminal output with visually appealing
text formatting.

Features:
• Multiple font options for different styles
• Customizable output formatting
• Integration with Anvil's terminal output system
• Support for various ASCII art styles

This command is useful for creating distinctive headers, banners, or decorative
elements in scripts and terminal applications.`
