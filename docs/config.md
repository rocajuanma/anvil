# Configuration Management

Anvil's configuration management system allows you to sync dotfiles and configuration files across machines using **PRIVATE** GitHub repositories for security.

## CRITICAL SECURITY REQUIREMENT

**Anvil REQUIRES private repositories for configuration management.**

Configuration files contain sensitive data that must NEVER be exposed publicly:

- **API keys and tokens** - GitHub tokens, cloud provider credentials
- **Personal paths** - Home directories, SSH key paths, private file locations
- **System information** - Usernames, email addresses, development settings
- **Authentication data** - SSH configurations, git credentials

**Security Guarantees:**

- Anvil **BLOCKS** all pushes to public repositories
- Repository privacy is **verified before every push**
- Clear error messages guide users to make repositories private
- **Push operations will FAIL** if repository is public

## Overview

The `config` command provides centralized management of configuration files and dotfiles for your development environment using **private GitHub repositories only**.

### Features

- **Centralize Configurations** - Store all your dotfiles and configuration files in a **PRIVATE** GitHub repository
- **Sync Across Machines** - Keep consistent configurations across all your development environments
- **Directory-Specific Operations** - Pull only the configuration directory you need
- **Version Control** - Track changes to your configurations over time with full privacy protection
- **Team Sharing** - Share team configurations and best practices through private repositories
- **Security-First** - Mandatory private repository validation prevents data exposure

### Current Implementation Status

- **Pull**: Fully implemented with directory-specific pulling
- **Show**: View configurations and settings
- **Sync**: Reconcile configuration state with system reality
- **Push**: Upload configurations to GitHub repository
- **Import**: Import group definitions from local files or URLs

## Commands

### anvil config pull [directory]

Pull configuration files from a specific directory in your GitHub repository to your local machine.

```bash
anvil config pull cursor
anvil config pull vscode
```

**How it works:**

- Automatically fetches the latest changes from your repository
- Copies all files from the specified directory to `~/.anvil/temp/[directory]`
- Guarantees you get the most up-to-date configurations every time

### anvil config show [directory]

Display configuration files and settings for easy viewing and inspection.

```bash
anvil config show
anvil config show cursor
```

**Features:**

- **Single File Display** - Shows file content directly in terminal
- **Multiple Files** - Shows tree structure with file listings
- **Smart File Detection** - Automatically determines best display method

#### Section-Specific Display Flags

View specific sections of your anvil settings with targeted flags:

```bash
anvil config show --groups          # Show only groups
anvil config show -g                # Short form for groups
anvil config show --configs         # Show only config source directories
anvil config show -c                # Short form for configs
anvil config show --git             # Show only git configuration
anvil config show --github          # Show only GitHub configuration
```

**Available Flags:**

- **`--groups/-g`** - Display only groups (built-in and custom) with tool counts
- **`--configs/-c`** - Display only configured source directories for apps
- **`--git`** - Display only git configuration (username, email, SSH key path)
- **`--github`** - Display only GitHub configuration (repository, branch, local path)

**Use Cases:**

- **Quick Reference** - `anvil config show --groups` for easy group overview
- **Path Management** - `anvil config show --configs` to see configured app paths
- **Git Setup** - `anvil config show --git` to verify git configuration
- **Repository Check** - `anvil config show --github` to confirm GitHub settings

### anvil config sync [app-name]

Move pulled configuration files from the temp directory to their local destinations with automatic archiving.

```bash
anvil config sync
anvil config sync obsidian
anvil config sync cursor
anvil config sync --dry-run
```

**How it works:**

- **Smart Path Resolution** - Uses your settings.yaml configs section for destinations
- **Automatic Archiving** - Backs up existing configurations before overwriting
- **Dry-Run Support** - Preview changes before applying them
- **Clear Error Messages** - Helpful guidance when configs or paths are missing

### anvil config push [app-name]

Push configuration files to your GitHub repository with automated branch creation and change tracking.

```bash
anvil config push
anvil config push cursor
```

**Key Features:**

- **Pre-Push Validation** - Detects differences before creating branches or commits
- **No Unnecessary Operations** - Skips Git operations when configurations are up-to-date
- **Repository Organization** - Maintains clean directory structure
- **Workflow Integration** - Seamless integration with GitHub pull request workflow

### anvil config import [file-or-url]

Import group definitions from local files or URLs with comprehensive validation and conflict detection.

```bash
anvil config import ./team-groups.yaml
anvil config import https://raw.githubusercontent.com/company/shared-configs/main/groups.yaml
```

**Key Features:**

- **Flexible Sources** - Import from local files or publicly accessible URLs
- **Comprehensive Validation** - Validates group names, application names, and structure
- **Conflict Detection** - Prevents overwriting existing groups with clear error messages
- **Tree Display** - Shows visual preview of groups and applications before import
- **Interactive Confirmation** - Requires user approval before making changes
- **Security-First** - Only imports group definitions, ignoring sensitive configuration data

## Setup

### 1. Initialize Anvil

```bash
anvil init
```

### 2. Create a PRIVATE GitHub Repository

Create a new private repository on GitHub to store your configurations.

### 3. Configure Repository in Settings

Edit `~/.anvil/settings.yaml` to add your repository information:

```yaml
github:
  config_repo: "username/<repo_name>"
  branch: "<main_branch_name>"
  token_env_var: "GITHUB_TOKEN" # optional
```

### 4. Set Up Authentication

#### Option 1: GitHub Token (Recommended)

```bash
export GITHUB_TOKEN="your_token_here"
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/user
```

#### Option 2: SSH Keys

```bash
ssh-keygen -t ed25519 -C "your.email@example.com"
ssh-add ~/.ssh/id_ed25519
```

## Example Workflows

### Basic Configuration Management

```bash
anvil init
anvil config pull cursor
anvil config show cursor
anvil config sync cursor
```

### Configuration Inspection and Verification

```bash
# Check your groups setup
anvil config show --groups

# Verify configured app paths
anvil config show --configs

# Check git configuration
anvil config show --git

# Verify GitHub repository settings
anvil config show --github
```

### Team Development Workflow

```bash
anvil init
anvil config import https://raw.githubusercontent.com/team/shared-configs/main/groups.yaml
anvil install dev
anvil config push
```

### Configuration Backup and Sync Workflow

```bash
anvil doctor connectivity
anvil config push
anvil config pull cursor
anvil config sync cursor
```

## Troubleshooting

### Common Issues

**Authentication Failed**

```bash
echo $GITHUB_TOKEN
ssh -T git@github.com
```

**Repository Not Found**

```bash
anvil config show
```

**Permission Denied**

```bash
anvil doctor connectivity
```

**Section Flags Not Working**

```bash
# Ensure you're using the flags with anvil settings (no directory argument)
anvil config show --groups
anvil config show --git

# Check if settings file exists
anvil config show
```

## Best Practices

1. **Always use private repositories** for configuration management
2. **Test connectivity** before major configuration operations
3. **Use dry-run mode** to preview changes before applying them
4. **Backup important configurations** before syncing
5. **Keep repository organized** with clear directory structure