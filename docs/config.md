# Configuration Management

Anvil's configuration management system allows you to sync dotfiles and configuration files across machines using **PRIVATE** GitHub repositories for security.

## 🚨 CRITICAL SECURITY REQUIREMENT

**Anvil REQUIRES private repositories for configuration management.**

Configuration files contain sensitive data that must NEVER be exposed publicly:

- **🔑 API keys and tokens** - GitHub tokens, cloud provider credentials
- **📁 Personal paths** - Home directories, SSH key paths, private file locations
- **⚙️ System information** - Usernames, email addresses, development settings
- **🔐 Authentication data** - SSH configurations, git credentials

**🛡️ Security Guarantees:**

- ✅ Anvil **BLOCKS** all pushes to public repositories
- ✅ Repository privacy is **verified before every push**
- ✅ Clear error messages guide users to make repositories private
- ⚠️ **Push operations will FAIL** if repository is public

## Overview

The `config` command provides centralized management of configuration files and dotfiles for your development environment using **private GitHub repositories only**.

### Features

- **📁 Centralize configurations** - Store all your dotfiles and configuration files in a **PRIVATE** GitHub repository
- **🔄 Sync across machines** - Keep consistent configurations across all your development environments
- **📦 Directory-specific operations** - Pull only the configuration directory you need
- **🛡️ Version control** - Track changes to your configurations over time with full privacy protection
- **👥 Team sharing** - Share team configurations and best practices through private repositories
- **🔒 Security-first** - Mandatory private repository validation prevents data exposure

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

🔧 Loading anvil configuration...
✅ Configuration loaded successfully

🚨 SECURITY REMINDER: Configuration files contain sensitive data
   • API keys, tokens, and credentials
   • Personal file paths and system information
   • Private development environment details

🛡️  Anvil REQUIRES private repositories for security
   • Repository 'username/dotfiles' must be PRIVATE
   • Public repositories will be BLOCKED
   • Verify at: https://github.com/username/dotfiles/settings

🔧 Setting up authentication...
✅ GitHub token found in environment

🔧 Preparing to push anvil configuration...
Repository: username/dotfiles
Branch: main
Settings file: /Users/username/.anvil/settings.yaml

🔧 Requesting user confirmation...
? Do you want to push your anvil settings to the repository? (y/N): y

🔧 Pushing configuration to repository...

# If no changes:
✅ Configuration is up-to-date!
Local anvil settings match the remote repository.
No changes to push.

# If changes detected:
🔒 Repository privacy verified - safe to push configuration data
Differences detected between local and remote configuration
Created and switched to branch: config-push-18072025-2147
Changes detected, proceeding with commit...
✅ Committed changes: anvil[push]: anvil
✅ Pushed branch 'config-push-18072025-2147' to origin

=== Push Complete! ===
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

**🚧 Status: In Development**

Push application-specific configurations to your repository:

```bash
$ anvil config push cursor

Expected functionality:
  • Create timestamped branch: config-push-<DDMMYYYY>-<HHMM>
  • Commit message: anvil[push]: cursor
  • Push cursor configs to /cursor directory in repository
  • Create pull request for review

🚧 Status: In Development
📅 Expected: Future release

For now, use 'anvil config push' to push anvil settings only.
```

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

### 2. Create a PRIVATE GitHub Repository

**🚨 CRITICAL: Repository MUST be private for security**

1. **Create a new repository on GitHub:**

   - Go to https://github.com/new
   - ⚠️ **IMPORTANT**: Select **"Private"** repository
   - Name it something like `dotfiles`, `configs`, or `dev-environment`
   - ✅ Verify the repository shows as **Private** before proceeding

2. **Verify repository privacy:**
   - Check repository settings at `https://github.com/username/repository/settings`
   - Ensure "Visibility" section shows **"Private repository"**
   - If public, change it: **Danger Zone** → **Change repository visibility** → **Private**

### 3. Configure GitHub Repository

Edit your `~/.anvil/settings.yaml` file:

```yaml
github:
  config_repo: "username/repository" # Your PRIVATE GitHub repository
  branch: "main" # Branch to use (main/master)
  token_env_var: "GITHUB_TOKEN" # Environment variable for authentication
```

### 4. Set Up Authentication

#### Option 1: GitHub Token (Recommended)

```bash
# Create a personal access token at: https://github.com/settings/tokens
# Select these scopes: repo, read:user

# Set environment variable (add to ~/.zshrc or ~/.bashrc)
export GITHUB_TOKEN="your_token_here"

# Verify token works
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/user
```

#### Option 2: SSH Keys

```bash
# Generate SSH key
ssh-keygen -t ed25519 -C "your.email@example.com"

# Add to SSH agent
ssh-add ~/.ssh/id_ed25519

# Copy public key and add to GitHub
cat ~/.ssh/id_ed25519.pub | pbcopy
# Add at: https://github.com/settings/keys

# Test SSH connection
ssh -T git@github.com
```

## Example Workflows

### Basic Configuration Management

```bash
# 1. Set up repository and authentication
anvil init
# Edit ~/.anvil/settings.yaml with your GitHub repository

# 2. Verify connectivity
anvil doctor connectivity

# 3. Pull application configs from repository
anvil config pull cursor
anvil config show cursor

# 4. Make local changes to anvil settings
# 5. Push anvil settings to repository
anvil config push

# 6. On another machine, pull and sync
anvil config pull team-dev
anvil config show team-dev
```

### Team Development Workflow

```bash
# Team member sets up their environment
anvil init

# Configure team repository
# Edit ~/.anvil/settings.yaml:
#   github.config_repo: "team/dev-configs"

# Pull shared configurations
anvil config pull shared-dev
anvil config sync team-dev

# Push personal anvil settings (if desired)
anvil config push

# Verify setup is complete
anvil doctor
```

### Configuration Backup and Sync Workflow

```bash
# 1. Verify connectivity before config operations
anvil doctor connectivity

# 2. Make changes to your anvil settings locally
# 3. Push changes to repository for backup
anvil config push

# 4. On another machine, pull latest settings
anvil config pull anvil

# 5. Install any missing applications
anvil config sync

# 6. Verify setup is complete
anvil doctor
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
