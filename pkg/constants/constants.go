package constants

const ANVIL_LONG_DESCRIPTION = `Anvil is a powerful automation CLI tool designed to streamline development workflows 
and personal tool configuration. It provides a comprehensive suite of commands for managing
development environments, automating installations, and maintaining consistent configurations
across different systems.

Key features:
- Automated tool installation and validation
- Environment configuration management
- Asset synchronization with GitHub
- ASCII art generation for enhanced terminal output
- Cross-platform compatibility with intelligent defaults`

const INIT_COMMAND_LONG_DESCRIPTION = `The init command bootstraps your Anvil CLI environment by performing a complete
initialization process. This is the first command you should run after installing Anvil.

What it does:
• Validates and installs required system tools (Git, cURL, Homebrew on macOS)
• Creates necessary configuration directories (~/.anvil, ~/.anvil/cache, ~/.anvil/data)
• Generates a default settings.yaml configuration file with your system preferences
• Checks your local environment for common development configurations
• Provides actionable recommendations for completing your setup

The init command is designed to be idempotent - it can be run multiple times safely
without overwriting existing configurations. It will only create missing components
and warn you about any environment issues that need attention.

This command is essential for ensuring Anvil has the proper foundation to operate
effectively in your development environment.`

const SETUP_COMMAND_LONG_DESCRIPTION = `The setup command provides automated batch installation of development tools and
applications using Homebrew (on macOS) or other package managers.

This command reads from your Anvil configuration to install predefined sets of
tools, ensuring consistent development environments across different machines.
It's particularly useful for setting up new development machines or maintaining
tool consistency across team members.`

const PUSH_COMMAND_LONG_DESCRIPTION = `The push command enables you to upload and synchronize your local assets, configurations,
and dotfiles to GitHub for backup and sharing purposes.

Features:
• Selective asset pushing based on configuration
• Automatic Git repository management
• Conflict resolution and merge handling
• Support for various asset types (dotfiles, configs, scripts)

The command takes an argument to specify which type of assets should be pushed,
allowing for granular control over what gets synchronized to your remote repository.`

const PULL_COMMAND_LONG_DESCRIPTION = `The pull command allows you to download and synchronize assets, configurations,
and dotfiles from your GitHub repository to your local machine.

This is particularly useful for:
• Setting up new development environments
• Synchronizing configurations across multiple machines
• Restoring configurations after system changes
• Sharing configurations with team members

The command takes an argument to specify which type of assets should be retrieved,
providing flexibility in what gets synchronized to your local environment.`

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
