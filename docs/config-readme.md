# Configuration Management with Anvil

This guide provides comprehensive documentation for Anvil's configuration management system, which allows you to sync configuration files and dotfiles across machines using GitHub repositories.

## 📋 Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Configuration Setup](#configuration-setup)
- [Pull Command](#pull-command)
- [Push Command](#push-command)
- [Settings Configuration](#settings-configuration)
- [Troubleshooting](#troubleshooting)
- [Advanced Usage](#advanced-usage)
- [Roadmap](#roadmap)

## 🎯 Overview

Anvil's configuration management system enables you to:

- **📁 Centralize configurations** - Store all your dotfiles and configuration files in a GitHub repository
- **🔄 Sync across machines** - Keep consistent configurations across all your development environments
- **📦 Directory-specific pulls** - Pull only the configuration directory you need
- **🛡️ Version control** - Track changes to your configurations over time
- **👥 Team sharing** - Share team configurations and best practices

### Current Implementation Status

- ✅ **Config Pull**: Fully implemented with directory-specific pulling to temp locations
- 🚧 **Config Push**: In development (coming soon)

## 🚀 Quick Start

### 1. Initialize Anvil

```bash
anvil init
```

### 2. Configure GitHub Repository

Edit your `~/.anvil/settings.yaml` file:

```yaml
github:
  config_repo: "username/dotfiles" # Your GitHub repository
  branch: "main" # Branch to use (main/master)
  local_path: "~/.anvil/dotfiles" # Local storage path
  token_env_var: "GITHUB_TOKEN" # Environment variable for token
```

### 3. Set Up Authentication

**Option A: SSH Key (Recommended)**

```bash
# Generate SSH key if you don't have one
ssh-keygen -t ed25519 -C "your.email@example.com"

# Add to your SSH agent
ssh-add ~/.ssh/id_ed25519

# Add SSH key to GitHub account
cat ~/.ssh/id_ed25519.pub | pbcopy  # Copy to clipboard
```

**Option B: GitHub Token**

```bash
# Create a personal access token at github.com/settings/tokens
# Set it as environment variable
export GITHUB_TOKEN="your_token_here"

# Add to your shell profile for persistence
echo 'export GITHUB_TOKEN="your_token_here"' >> ~/.zshrc
```

### 4. Pull Configuration Files

```bash
# Pull a specific configuration directory
anvil config pull cursor
anvil config pull vs-code
anvil config pull zsh

# View pulled configurations
anvil config show cursor
anvil config show vs-code
```

## 🔧 Configuration Setup

### GitHub Repository Structure

Your configuration repository should be organized by application or tool:

```
your-config-repo/
├── cursor/
│   ├── settings.json
│   ├── keybindings.json
│   └── snippets/
├── vs-code/
│   ├── settings.json
│   ├── keybindings.json
│   └── extensions.json
├── zsh/
│   ├── .zshrc
│   ├── .zsh_aliases
│   └── .zsh_functions
├── git/
│   ├── .gitconfig
│   └── .gitignore_global
└── README.md
```

### Settings.yaml Configuration

Complete configuration example:

```yaml
version: 1.0.0
directories:
  config: /Users/username/.anvil

github:
  config_repo: "username/dotfiles" # Required: GitHub repository
  branch: "main" # Required: Git branch to use
  local_path: "~/.anvil/dotfiles" # Required: Local clone location
  token_env_var: "GITHUB_TOKEN" # Optional: Environment variable for token

git:
  username: "Your Name" # Git user configuration
  email: "your.email@example.com" # Git email configuration
  ssh_key_path: "~/.ssh/id_ed25519" # SSH key path
  ssh_dir: "~/.ssh" # SSH directory
```

### GitHub URL Format Support

Anvil automatically normalizes various GitHub URL formats:

```bash
# All these formats are supported and auto-corrected:
username/repository                          # ← Preferred format
https://github.com/username/repository
https://github.com/username/repository.git
git@github.com:username/repository
git@github.com:username/repository.git
github.com/username/repository
```

## 📥 Pull Command

### Basic Usage

### Configuration Commands

The config command provides three main operations:

- **`anvil config pull [directory]`** - Downloads configuration files from your GitHub repository
- **`anvil config show [directory]`** - Views configuration files (anvil settings or pulled configs)
- **`anvil config push [directory]`** - Uploads configuration files to your GitHub repository (coming soon)

#### Pull Command

The `anvil config pull` command downloads configuration files from your GitHub repository:

```bash
anvil config pull [directory]
```

**⚡ Always Fresh**: Every pull command automatically fetches the latest changes from your repository, ensuring you always get the most up-to-date configurations.

### Examples

```bash
# Pull Cursor configuration
anvil config pull cursor

# Pull VS Code configuration
anvil config pull vs-code

# Pull Zsh configuration
anvil config pull zsh

# Pull Git configuration
anvil config pull git
```

### Detailed Process

When you run `anvil config pull cursor`, Anvil:

1. **Validates configuration** - Checks GitHub repository and branch settings
2. **Authenticates** - Uses SSH keys or GitHub token for repository access
3. **Clones/updates repository** - Downloads or updates the local repository copy
4. **⚡ Always fetches latest changes** - Runs `git fetch` and `git pull` to ensure you get the most up-to-date files
5. **Copies directory** - Copies only the specified directory to `~/.anvil/temp/[directory]`
6. **Lists files** - Shows what files were pulled
7. **Provides guidance** - Suggests next steps for applying configurations

### Output Example

```bash
$ anvil config pull cursor

🔧 Using branch: main

=== Pulling Configuration Directory: cursor ===

Repository: username/dotfiles
Branch: main
Target directory: cursor
✅ GitHub token found in environment variable: GITHUB_TOKEN
🔧 Validating repository access and branch configuration...
✅ Repository and branch configuration validated
🔧 Setting up local repository...
✅ Local repository ready
🔧 Pulling latest changes...
✅ Repository updated
🔧 Copying configuration directory...
✅ Configuration directory copied to temp location

=== Pull Complete! ===

Configuration directory 'cursor' has been pulled from: username/dotfiles
Files are available at: /Users/username/.anvil/temp/cursor

Copied files:
  • settings.json
  • keybindings.json
  • snippets/javascript.json

Next steps:
  • Review the pulled configuration files in: /Users/username/.anvil/temp/cursor
  • Use 'anvil config show [directory]' to view configuration content
  • Apply/copy configurations to their destination as needed
  • Use 'anvil config push' to upload any local changes
```

### Current Implementation Details

**⚡ Always Up-to-Date**: Every pull command automatically fetches the latest changes from your GitHub repository using `git fetch` and `git pull`. You're guaranteed to get the most recent version of your configurations.

**⚠️ Important**: The current pull implementation copies files to a temporary location (`~/.anvil/temp/[directory]`) for manual review and application.

**Future Enhancement**: Anvil will eventually automatically apply configurations to their destination directories (e.g., `~/Library/Application Support/Cursor/User/settings.json`).

### Repository Update Behavior

**🔄 Always Synchronized**: Every time you run `anvil config pull`, the command automatically:

1. **Fetches** the latest changes from your repository (`git fetch`)
2. **Pulls** those changes into the local copy (`git pull`)
3. **Copies** the updated files to your temp directory

This means if you:

- Add new files to your repository
- Update existing configurations
- Make any changes remotely

The next `anvil config pull` command will **always** get those latest changes. No manual repository updates needed!

**Example Scenario**:

```bash
# 1. Repository has "app1" directory with 1 file
anvil config pull app1  # Gets 1 file

# 2. You update repository remotely - now "app1" has 3 files
anvil config pull app1  # Gets all 3 files (automatically fetched latest)
```

### Authentication Methods

Anvil supports multiple authentication methods:

1. **SSH Key Authentication (Recommended)**

   - Automatic detection of SSH keys
   - Supports custom SSH key paths
   - No token storage required

2. **GitHub Token Authentication**

   - Uses environment variable for security
   - Supports private repositories
   - Automatic HTTPS URL conversion

3. **Public Repository Access**
   - No authentication required
   - Automatically falls back to HTTPS

## 📄 Show Command

### View Configuration Files

The `anvil config show` command displays configuration files for easy viewing and inspection:

```bash
# View main anvil settings
anvil config show

# View pulled configuration directory
anvil config show [directory]
```

### Features

- **Single File Display**: Shows file content directly in terminal
- **Multiple Files**: Shows tree structure with file listings
- **Smart Detection**: Automatically determines best display method
- **Helpful Errors**: Clear messages for missing directories with suggestions

### Examples

```bash
# View your anvil settings.yaml
anvil config show

# View pulled Cursor configurations
anvil config show cursor

# View pulled VS Code configurations
anvil config show vs-code
```

Perfect for reviewing pulled configurations before applying them or checking your current anvil settings.

## 📤 Push Command

### Status: In Development

The `anvil config push` command is currently under development. When complete, it will:

- Upload local configuration changes to your GitHub repository
- Commit changes with descriptive messages
- Support selective directory pushing
- Maintain proper version history

### Planned Usage

```bash
# Push specific directory (planned feature)
anvil config push cursor

# Push all changes (planned feature)
anvil config push --all

# Push with custom commit message (planned feature)
anvil config push cursor -m "Update cursor settings"
```

**Coming Soon**: Full implementation with automatic change detection, commit generation, and push capabilities.

## ⚙️ Settings Configuration

### Complete Settings Reference

```yaml
version: 1.0.0

directories:
  config: /Users/username/.anvil

# GitHub Configuration (Required for config commands)
github:
  config_repo: "username/dotfiles" # GitHub repository (username/repo format)
  branch: "main" # Git branch (main/master/develop/etc.)
  local_path: "~/.anvil/dotfiles" # Local repository storage path
  token_env_var: "GITHUB_TOKEN" # Environment variable for GitHub token

# Git Configuration (Recommended)
git:
  username: "Your Name" # Git commit author name
  email: "your.email@example.com" # Git commit author email
  ssh_key_path: "~/.ssh/id_ed25519" # Path to SSH private key
  ssh_dir: "~/.ssh" # SSH configuration directory

# Tool Configuration
tools:
  required_tools: [git, curl]
  optional_tools: [brew, docker, kubectl]

# Environment Variables
environment:
  EDITOR: "code"
  GITHUB_TOKEN: "your_token_here" # Alternative to token_env_var
```

### Configuration Validation

Anvil automatically validates your configuration and provides helpful error messages:

- **Repository format validation** - Ensures proper GitHub URL format
- **Branch existence checks** - Verifies the specified branch exists
- **Authentication verification** - Tests repository access
- **Path validation** - Checks directory permissions and accessibility

## 🐛 Troubleshooting

### Common Issues and Solutions

#### Branch Configuration Errors

**Problem**: Branch doesn't exist in repository

```bash
❌ Branch Configuration Error

The branch 'nonexistent-branch' does not exist in repository 'username/repo'.

✅ Available branches in repository:
    - main
    - master
    - development
```

**Solution**: Update your `settings.yaml` with a valid branch:

```yaml
github:
  branch: "main" # Use an existing branch
```

#### Authentication Issues

**Problem**: Repository access denied

```bash
❌ Pull failed: repository validation failed: cannot access repository
```

**Solutions**:

1. **For SSH**: Ensure your SSH key is added to GitHub
2. **For Token**: Check the `GITHUB_TOKEN` environment variable
3. **For Private repos**: Ensure you have repository access

#### Directory Not Found

**Problem**: Specified directory doesn't exist in repository

```bash
❌ Pull failed: directory 'nonexistent-dir' does not exist in repository username/repo
```

**Solution**: Check your repository structure and use existing directories:

```bash
# List repository contents to see available directories
git clone https://github.com/username/repo.git temp
ls temp/
rm -rf temp
```

#### Git Configuration Missing

**Problem**: Git user not configured

```bash
⚠️ Git user configuration is incomplete
```

**Solution**: Add git configuration to your `settings.yaml`:

```yaml
git:
  username: "Your Name"
  email: "your.email@example.com"
```

### Getting Detailed Error Information

When pull fails due to branch issues, Anvil provides:

- **Detailed error explanation** with the specific issue
- **List of available branches** in your repository
- **Step-by-step fix instructions** with examples
- **Direct path to settings file** for easy editing

### Repository Synchronization

**❓ Common Question**: "Do I need to manually update my local repository?"

**✅ Answer**: **No!** Every `anvil config pull` command automatically fetches and pulls the latest changes from your remote repository. You never need to manually run `git pull` or update the repository yourself.

## 🔧 Advanced Usage

### Custom Repository Structures

You can organize your repository however you prefer:

```
dotfiles/
├── applications/
│   ├── cursor/
│   └── vscode/
├── shell/
│   ├── zsh/
│   └── bash/
└── tools/
    ├── git/
    └── tmux/
```

Pull from nested directories:

```bash
anvil config pull applications/cursor
anvil config pull shell/zsh
anvil config pull tools/git
```

### Multiple Environment Support

Use different branches for different environments:

```yaml
# Development environment
github:
  branch: "development"

# Production environment
github:
  branch: "main"

# Personal environment
github:
  branch: "personal"
```

### Team Configuration Sharing

Share configurations across team members:

1. **Create team repository** with shared configurations
2. **Use consistent branch naming** (e.g., `team-frontend`, `team-backend`)
3. **Document configuration setup** in repository README
4. **Maintain environment-specific branches** when needed

## 🗺️ Roadmap

### Current Status: Pull Implementation ✅

- ✅ Directory-specific pulling
- ✅ GitHub repository integration
- ✅ SSH and token authentication
- ✅ Branch validation and error handling
- ✅ Temporary file location management

### Next: Push Implementation 🚧

**Phase 1: Basic Push (In Development)**

- 📋 Upload local changes to repository
- 📋 Automatic commit message generation
- 📋 Selective directory pushing
- 📋 Change detection and validation

**Phase 2: Advanced Push Features (Planned)**

- 📋 Interactive change selection
- 📋 Custom commit messages
- 📋 Merge conflict resolution
- 📋 Push hooks and validation

### Future: Automatic Application 🔮

**Phase 3: Smart Configuration Application (Planned)**

- 📋 Automatic detection of application config locations
- 📋 Safe backup before applying configurations
- 📋 Application restart handling
- 📋 Rollback capabilities
- 📋 Application-specific configuration validators

**Target Application Support**:

- 📋 VS Code / Cursor (`~/Library/Application Support/`)
- 📋 Shell configurations (`~/.zshrc`, `~/.bashrc`)
- 📋 Git (`~/.gitconfig`)
- 📋 SSH (`~/.ssh/config`)
- 📋 Homebrew (`/opt/homebrew/`)
- 📋 Custom application paths

### Integration Features (Future)

- 📋 Configuration diff visualization
- 📋 Automatic backup before applying changes
- 📋 Configuration validation and testing
- 📋 Team configuration templates
- 📋 Configuration dependency management

---

## 💡 Tips and Best Practices

### Repository Organization

1. **Use descriptive directory names** that match application names
2. **Include README files** explaining configuration purposes
3. **Use consistent naming conventions** across your team
4. **Keep application-specific configurations separate**

### Security Best Practices

1. **Never commit secrets** or API keys to your configuration repository
2. **Use environment variables** for sensitive configuration
3. **Use SSH keys** instead of tokens when possible
4. **Review configurations** before committing

### Workflow Recommendations

1. **Test configurations** in the temp directory before applying
2. **Keep backups** of working configurations
3. **Use descriptive commit messages** for configuration changes
4. **Document breaking changes** in your repository

---

**Ready to start managing your configurations?**

1. Run `anvil init` to set up Anvil
2. Configure your GitHub repository in `~/.anvil/settings.yaml`
3. Start pulling configurations with `anvil config pull [directory]`

For more help, see the [main documentation](README.md) or check out [examples](EXAMPLES.md).
