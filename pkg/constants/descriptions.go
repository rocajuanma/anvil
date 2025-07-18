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

Key Features:
• 🎯 Effortlessly install any macOS application using Homebrew
• 📝 Automatically track individually installed apps in your settings.yaml
• Organize tools with flexible group management for different workflows
• ⚡ Zero configuration required: ready to use with smart defaults
• 🍺 Seamless Homebrew integration for automated installation and updates
• 🔍 Dry-run mode to preview changes before making them
• Beautiful output with colored, structured progress indicators`

const INIT_COMMAND_LONG_DESCRIPTION = `The init command bootstraps your Anvil CLI environment by performing a complete
initialization process. This is the first command you should run after installing Anvil.

What it does:
• ✅ Validates and installs required system tools (Git, cURL, Homebrew)
• Creates necessary configuration directory (~/.anvil)
• Generates a default settings.yaml configuration file with your system preferences
• Checks your local environment for common development configurations
• 💡 Provides actionable recommendations for completing your setup

This command is designed specifically for macOS and requires Homebrew for tool management.`

const INSTALL_COMMAND_LONG_DESCRIPTION = `The install command provides dynamic installation of development tools and applications
for macOS using Homebrew. It supports both group-based and individual installations.

Installation Modes:
• anvil install [group-name]    - Install all tools in a predefined group
• anvil install [app-name]      - Install any individual application via brew

Available Groups: 
• dev - Essential development tools
• new-laptop - Essential applications for new machines
• Custom groups you define in settings.yaml

Special Features:
• 📝 Automatic App Tracking: Individual apps installed via 'anvil install [app-name]' are
  automatically tracked in your settings.yaml under tools.installed_apps, making it easy to
  recreate your environment on new machines
• Smart Deduplication: Apps in groups/required_tools are not tracked separately
• 💡 Smart Discovery: Can't find an app? Get helpful suggestions for alternatives

Use 'anvil install --list' to see all available groups and tracked apps.`

const CONFIG_COMMAND_LONG_DESCRIPTION = `The config command provides centralized management of configuration files and dotfiles
for your development environment. It serves as a parent command for pull and push operations
to ensure all configuration-related actions are properly organized.

Subcommands:
• anvil config pull [directory]    - Pull configuration files from remote repository
• anvil config push [directory]    - Push configuration files to remote repository

GitHub Repository Configuration:
The 'github.config_repo' field in settings.yaml should be in the format 'username/repository'.

Supported input formats (automatically corrected):
• username/repository (preferred format)
• https://github.com/username/repository
• https://github.com/username/repository.git
• git@github.com:username/repository.git  
• github.com/username/repository

Example: 'github.config_repo: octocat/Hello-World'

This command structure ensures that all pull and push operations are scoped to configuration
files only, providing clear separation between configuration management and other operations.
Use this command to maintain consistent configuration across different development environments.`

const PUSH_COMMAND_LONG_DESCRIPTION = `The push command enables you to upload and synchronize your local configuration files
and dotfiles to GitHub for backup and sharing purposes.

Features:
• Selective config pushing based on configuration
• 🔄 Automatic Git repository management
• Conflict resolution and merge handling
• Support for various config types (dotfiles, configs, scripts)

The command takes an argument to specify which type of configurations should be pushed,
allowing for granular control over what gets synchronized to your remote repository.
This command is scoped to configuration files only for security and organization.`

const PULL_COMMAND_LONG_DESCRIPTION = `The pull command allows you to download and synchronize configuration files
from a specific directory in your GitHub repository to your local machine.

Usage: anvil config pull [directory]

How it works:
• 📥 Automatically fetches the latest changes from your repository (git fetch/pull)
• Copies all files from the specified directory to a temporary location (~/.anvil/temp/[directory])
• ✅ Guarantees you get the most up-to-date configurations every time

Perfect for:
• Setting up new development environments
• Synchronizing specific configurations across multiple machines
• Restoring configurations after system changes
• Sharing configurations with team members

GitHub Repository Setup:
Before using this command, configure your repository in ~/.anvil/settings.yaml:

  github:
    config_repo: "username/repository"  # Format: username/repository
    branch: "main"                      # Branch to pull from
    token_env_var: "GITHUB_TOKEN"       # Environment variable for authentication

The repository URL format is automatically validated and corrected if needed.
Supported formats include full URLs, SSH URLs, and domain-prefixed formats.

Example: anvil config pull cursor    # Pulls all files from the 'cursor' directory`

const DRAW_COMMAND_LONG_DESCRIPTION = `The draw command generates beautiful ASCII art representations of text using the
go-figure library for enhanced terminal output and visual appeal.

Features:
• Multiple font options for different styles
• Customizable output formatting
• Integration with Anvil's terminal output system
• Support for various ASCII art styles

Perfect for creating distinctive headers, banners, or decorative elements in scripts and terminal applications.`
