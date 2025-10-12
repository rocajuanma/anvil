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

const ANVIL_LONG_DESCRIPTION = `Anvil is a macOS automation CLI tool for managing development environments, installing tools via Homebrew, and syncing configuration files. It is designed to help developers set up, maintain, and reproduce their development environments with ease and consistency.

By automating the installation of essential tools and the synchronization of configuration files, Anvil reduces manual setup steps and helps ensure that your environment is always up to date. This not only minimizes the risk of configuration drift but also saves you valuable time, especially when setting up new machines or restoring your environment after changes.

Key capabilities:
• Install and manage macOS applications and CLI tools
• Sync configuration files and dotfiles across machines`

const INIT_COMMAND_LONG_DESCRIPTION = `Initialize your Anvil environment. This is the first command you should run after installing Anvil.

What it does:
• Installs required system tools (Git, cURL, Homebrew)
• Creates configuration directory (~/.anvil) and settings.yaml
• Validates your development environment`

const INSTALL_COMMAND_LONG_DESCRIPTION = `Install development tools individually or in groups using Homebrew.

Define custom groups in settings.yaml`

const CONFIG_COMMAND_LONG_DESCRIPTION = `Manage configuration files and dotfiles for your anvil environment.

Configure 'github.config_repo' in settings.yaml to use this command.`

const PUSH_COMMAND_LONG_DESCRIPTION = `Upload local configuration files to GitHub with automated branch creation.

Configure 'github.config_repo' in settings.yaml to use this command.`

const PULL_COMMAND_LONG_DESCRIPTION = `Download configuration files from your GitHub repository.

Configure 'github.config_repo' in settings.yaml to use this command.`

const SHOW_COMMAND_LONG_DESCRIPTION = `Display configuration files and settings with intelligent formatting.`

const SYNC_COMMAND_LONG_DESCRIPTION = `Apply pulled configuration files to their local destinations with automatic archiving.

Safely applies configs with automatic backup of existing files.`

const DOCTOR_COMMAND_LONG_DESCRIPTION = `Run health checks to validate your anvil environment.

Health Check Categories:

ENVIRONMENT (3 checks)
  • anvil-init       - Verify anvil initialization is complete
  • settings-valid   - Validate settings.yaml structure and content
  • directory-structure - Check ~/.anvil directory structure

DEPENDENCIES (2 checks)
  • homebrew         - Verify Homebrew installation and updates (auto-fixable)
  • required-tools   - Check git and curl are installed

CONFIGURATION (3 checks)
  • git-config       - Validate git user.name and user.email (auto-fixable)
  • github-config    - Verify GitHub repository configuration
  • sync-config      - Check config sync settings (not yet implemented)

CONNECTIVITY (3 checks)
  • github-auth      - Test GitHub authentication and access
  • github-repo      - Verify repository accessibility
  • git-operations   - Test git clone and pull operations

Each check can be run independently by name or grouped by category.
Add --fix flag to auto-fix issues where supported.

Examples:
  anvil doctor                    # Run all 11 checks
  anvil doctor environment        # Run category (3 checks)
  anvil doctor git-config         # Run specific check
  anvil doctor git-config --fix   # Run check and auto-fix
  anvil doctor --fix              # Run all checks and auto-fix issues`

// Clean command descriptions
const CLEAN_COMMAND_LONG_DESCRIPTION = `Remove all content inside .anvil directories while preserving settings.yaml.

What it does:
• Removes temporary files, archives, and downloaded configurations
• Cleans temp/ and archive/ directories
• Removes dotfiles/ directory for clean git state
• Preserves settings.yaml file

Safe operation that never deletes your main configuration file.`

// Update command descriptions
const UPDATE_COMMAND_LONG_DESCRIPTION = `Update Anvil to the latest version from GitHub releases.

What it does:
• Downloads latest release information
• Runs official installation script
• Replaces current installation`
