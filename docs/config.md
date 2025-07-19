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
- 🚧 **Push**: In development (coming soon)

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

### anvil config push [directory]

_Coming soon_ - Push local configuration files to your GitHub repository.

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
