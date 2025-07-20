# Configuration Management

Anvil's configuration management system allows you to sync dotfiles and configuration files across machines using GitHub repositories.

## Overview

The `config` command provides centralized management of configuration files and dotfiles for your development environment. It serves as a parent command for configuration-related operations.

### Features

- **📁 Centralize configurations** - Store all your dotfiles and configuration files in a GitHub repository
- **🔄 Sync across machines** - Keep consistent configurations across all your development environments
- **📦 Directory-specific operations** - Pull only the configuration directory you need
- **🛡️ Version control** - Track changes to your configurations over time
- **👥 Team sharing** - Share team configurations and best practices

### Current Implementation Status

- ✅ **Pull**: Fully implemented with directory-specific pulling
- ✅ **Show**: View configurations and settings
- ✅ **Sync**: Reconcile configuration state with system reality
- ✅ **Push**: Upload configurations to GitHub repository

## Commands

### anvil config pull [directory]

Pull configuration files from a specific directory in your GitHub repository to your local machine.

```bash
# Pull cursor configuration
anvil config pull cursor

# Pull VS Code settings
anvil config pull vscode
```

**How it works:**

- Automatically fetches the latest changes from your repository
- Copies all files from the specified directory to `~/.anvil/temp/[directory]`
- Guarantees you get the most up-to-date configurations every time

### anvil config show [directory]

Display configuration files and settings for easy viewing and inspection.

```bash
# Show main anvil settings
anvil config show

# Show pulled configuration files
anvil config show cursor
```

**Features:**

- 📄 Single file display: Shows file content directly in terminal
- 📁 Multiple files: Shows tree structure with file listings
- ✅ Smart file detection: Automatically determines best display method

### anvil config sync [directory]

Reconcile configuration state between settings.yaml and system reality.

```bash
# Sync anvil settings (install missing apps)
anvil config sync

# Preview changes without applying them
anvil config sync --dry-run
```

**Features:**

- 📋 Smart difference analysis: Shows what's installed vs what's missing
- ✅ Confirmation prompts: Ask before making changes to system
- 🔍 Dry-run support: Preview changes without applying them

### anvil config push [app-name]

Push configuration files to your GitHub repository with automated branch creation and change tracking.

```bash
# Push anvil settings to repository
anvil config push

# Push application-specific configs (in development)
anvil config push cursor
anvil config push vscode
```

**How it works:**

#### Option 1: Anvil Settings Push (`anvil config push`)

- 🔍 **Smart Detection**: Compares local and remote configurations before proceeding
- 📁 **Organized Storage**: Always commits to `/anvil` directory in repository
- 🌿 **Timestamped Branches**: Creates branches with format `config-push-DDMMYYYY-HHMM`
- 💬 **Standardized Commits**: Uses commit message `anvil[push]: anvil`
- 🔗 **PR Ready**: Provides direct link to create pull request

**Example workflow:**

```bash
$ anvil config push

=== Push Anvil Configuration ===

🔧 Preparing to push anvil configuration...
Repository: username/dotfiles
Branch: main
Settings file: /Users/username/.anvil/settings.yaml
? Do you want to push your anvil settings to the repository? (y/N): y

# If no changes:
✅ Configuration is up-to-date!
Local anvil settings match the remote repository.
No changes to push.

# If changes detected:
Differences detected between local and remote configuration
Created and switched to branch: config-push-18072025-2147
Changes detected, proceeding with commit...
✅ Committed changes: anvil[push]: anvil
✅ Pushed branch 'config-push-18072025-2147' to origin

✅ Configuration push completed successfully!

📋 Push Summary:
  • Branch created: config-push-18072025-2147
  • Commit message: anvil[push]: anvil
  • Files committed: [anvil/settings.yaml]

🔗 Repository: https://github.com/username/dotfiles
🌿 Branch: config-push-18072025-2147

✅ You can now create a Pull Request on GitHub to merge these changes!
Direct link: https://github.com/username/dotfiles/compare/main...config-push-18072025-2147
```

#### Option 2: Application Config Push (`anvil config push <app-name>`)

🚧 **Status**: In Development

```bash
$ anvil config push cursor

=== Push 'cursor' Configuration ===

🚧 Application-specific configuration push is currently in development
📅 Expected: Future release

For now, use 'anvil config push' to push anvil settings only.
```

### Smart Filtering Configuration

Anvil's smart filtering allows you to control what gets synchronized between machines, solving the challenge of machine-specific vs. universal settings.

#### Problem Solved

When syncing configurations across machines, some settings should be universal (tool groups, GitHub repos) while others are machine-specific (git credentials, project paths). Smart filtering enables selective synchronization.

#### Configuration Structure

Add a `_sync_config` section to your `~/.anvil/settings.yaml`:

```yaml
version: 1.0.0
_sync_config:
  exclude_sections:
    - git # Don't sync git credentials
    - environment.machine_specific # Don't sync machine-specific environment vars
  template_sections:
    - git # Convert git section to template placeholders
  include_override: [] # Force include sections (overrides exclude)
  template_values: {} # Local template replacement values

# Universal settings (will be synced)
tools:
  required_tools: [git, curl, brew]
  installed_apps: [terraform, postman]
groups:
  dev: [git, zsh, iterm2, visual-studio-code]
github:
  config_repo: "username/dotfiles"
  branch: "main"

# Machine-specific settings (will be excluded/templated)
git:
  username: "Your Work Name" # Will become {{ REPLACE_USERNAME }}
  email: "you@company.com" # Will become {{ REPLACE_EMAIL }}
environment:
  EDITOR: "code" # Will be synced (universal)
  machine_specific: "/work/path" # Will be excluded
```

#### Smart Filtering Options

**Exclude Sections** - Remove sections entirely from sync:

- `git` - Exclude entire git configuration
- `environment` - Exclude all environment variables
- `environment.KEY` - Exclude specific environment variable
- `tool_configs` - Exclude tool-specific configurations

**Template Sections** - Convert values to templates:

- `git` - Replace username/email with `{{ REPLACE_USERNAME }}` / `{{ REPLACE_EMAIL }}`
- `environment` - Replace paths with template placeholders

**Include Override** - Force include normally excluded sections:

- Useful for temporary inclusion of specific sections

#### Enhanced Push Commands

```bash
# Preview what would be synced
anvil config push --show-filtered

# Override filtering with flags
anvil config push --include git --exclude tool_configs

# Normal push with filtering applied
anvil config push
```

#### Example Workflow

**Work Machine Setup:**

```yaml
_sync_config:
  exclude_sections: [git, environment.machine_specific]
  template_sections: [git]
git:
  username: "John Work"
  email: "john@company.com"
environment:
  machine_specific: "/company/projects"
```

**Personal Machine Setup:**

```yaml
_sync_config:
  exclude_sections: [git, environment.machine_specific]
  template_sections: [git]
git:
  username: "John Personal"
  email: "john@personal.com"
environment:
  machine_specific: "/home/john/projects"
```

**Synced Configuration** (what gets pushed to GitHub):

```yaml
# Universal settings synced
tools: [...]
groups: [...]
github: [...]

# Templated sections for setup guidance
git:
  username: "{{ REPLACE_USERNAME }}"
  email: "{{ REPLACE_EMAIL }}"
# Machine-specific sections excluded
# (environment.machine_specific not present)
```

This approach ensures your tool groups, GitHub configurations, and universal settings sync across machines while keeping machine-specific credentials and paths local.

**Key Features:**

- 🎯 **Pre-Push Validation**: Detects differences before creating branches or commits
- 🚫 **No Unnecessary Operations**: Skips Git operations when configurations are up-to-date
- 📦 **Repository Organization**: Maintains clean directory structure
- 🔄 **Workflow Integration**: Seamless integration with GitHub pull request workflow
- 🛡️ **Safe Operations**: Always creates new branches, never pushes directly to main

## Setup

### 1. Initialize Anvil

```bash
anvil init
```

### 2. Configure GitHub Repository

Edit your `~/.anvil/settings.yaml` file:

```yaml
github:
  config_repo: "username/repository" # Your GitHub repository
  branch: "main" # Branch to use (main/master)
  token_env_var: "GITHUB_TOKEN" # Environment variable for authentication
```

### 3. Set Up Authentication

#### Option 1: GitHub Token (Recommended)

1. Create a GitHub Personal Access Token:
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Generate new token with `repo` scope
2. Set environment variable:
   ```bash
   export GITHUB_TOKEN="your_token_here"
   ```

#### Option 2: SSH Keys

Ensure your SSH key is added to your GitHub account:

```bash
ssh -T git@github.com
```

### 4. Repository Structure

Organize your repository with directory-based configurations:

```
your-dotfiles-repo/
├── cursor/
│   ├── settings.json
│   └── keybindings.json
├── vscode/
│   ├── settings.json
│   └── extensions.json
└── zsh/
    ├── .zshrc
    └── .zsh_aliases
```

## Examples

### Personal Configuration Setup

```bash
# Pull your cursor settings
anvil config pull cursor

# Review what was pulled
anvil config show cursor

# Install missing apps from your settings
anvil config sync

# Push any local changes back to repository
anvil config push
```

### Team Configuration Sharing

```bash
# Pull team's development setup
anvil config pull team-dev

# See what tools the team uses
anvil config show team-dev

# Install team's recommended tools
anvil config sync team-dev --dry-run
anvil config sync team-dev
```

### Configuration Backup and Sync Workflow

```bash
# 1. Make changes to your anvil settings locally
# 2. Push changes to repository for backup
anvil config push

# 3. On another machine, pull latest settings
anvil config pull anvil

# 4. Install any missing applications
anvil config sync
```

### Repository Organization Example

After using `anvil config push`, your repository structure will look like:

```
your-dotfiles-repo/
├── anvil/
│   └── settings.yaml          # Anvil configuration (pushed via anvil config push)
├── cursor/
│   ├── settings.json          # Cursor settings (manual or future push)
│   └── keybindings.json
├── vscode/
│   ├── settings.json          # VS Code settings (manual or future push)
│   └── extensions.json
└── zsh/
    ├── .zshrc                 # Zsh configuration (manual)
    └── .zsh_aliases
```

## Repository URL Formats

Anvil automatically validates and corrects repository URLs. Supported formats:

- `username/repository` (preferred)
- `https://github.com/username/repository`
- `https://github.com/username/repository.git`
- `git@github.com:username/repository.git`
- `github.com/username/repository`

All formats are automatically normalized to `username/repository`.

## Troubleshooting

### Common Issues

**Authentication Failed**

```bash
# Check if token is set
echo $GITHUB_TOKEN

# Verify SSH access
ssh -T git@github.com
```

**Repository Not Found**

- Verify repository name in settings.yaml
- Check repository permissions
- Ensure repository exists and is accessible

**Branch Not Found**

- Verify branch name in settings.yaml
- Check if branch exists in repository
- Use `git branch -r` to list remote branches

**Directory Not Found**

- Verify directory exists in repository
- Check case sensitivity
- Use `anvil config show` to see available directories

### Getting Help

For detailed troubleshooting and examples, see:

- [Getting Started Guide](GETTING_STARTED.md)
- [Examples & Tutorials](EXAMPLES.md)
- [Contributing Guidelines](CONTRIBUTING.md)
