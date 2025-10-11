# Getting Started with Anvil CLI

Welcome to Anvil! This guide will help you get up and running quickly with Anvil CLI, from installation to your first successful tool setup.

## Table of Contents

- [What is Anvil?](#what-is-anvil)
- [Installation](#installation)
- [First Steps](#first-steps)
- [Basic Usage](#basic-usage)
- [Common Workflows](#common-workflows)
- [Understanding Configuration](#understanding-configuration)
- [Tips and Best Practices](#tips-and-best-practices)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

## What is Anvil?

Anvil is a CLI automation tool that helps developers:

- **🚀 Bootstrap development environments** quickly and reliably
- **📦 Install tools in logical groups** (development, new laptop, custom)
- **⚙️ Manage configurations** across different machines
- **🔧 Automate repetitive setup tasks** for individuals and teams
- **🩺 Validate environment health** with comprehensive diagnostic checks

### Key Concepts

- **Commands**: Actions you can perform (`init`, `install`, `config pull`, `config show`, `config sync`, `config push`, `clean`, `doctor`)
- **Groups**: Collections of related tools (`dev`, `essentials`, custom groups)
- **Configuration**: Settings stored in `~/.anvil/settings.yaml`
- **Health Checks**: Validation and troubleshooting via `anvil doctor`
- **Tools**: Individual applications or utilities that can be installed

## Installation

### Quick Install

#### Latest version
```bash
curl -sSL https://github.com/rocajuanma/anvil/releases/latest/download/install.sh | bash
```

#### From source

```bash
# Clone and build
git clone https://github.com/rocajuanma/anvil.git
cd anvil
go build -o anvil main.go

# Move to PATH (optional)
sudo mv anvil /usr/local/bin/
```

For detailed installation instructions, see our [Installation Guide](INSTALLATION.md).

## First Steps

### 1. Initialize Anvil

Before using any other commands, initialize Anvil:

```bash
anvil init
```

This command will:

- **✅ Validate and install required tools** - Git, cURL, Homebrew on macOS
- **📁 Create necessary configuration directory** - `~/.anvil/`
- **⚙️ Generate default settings.yaml** - Configuration file
- **🔍 Check your environment** - Provide recommendations


### 2. Verify Your Setup

After initialization, it's recommended to run a health check to ensure everything is working correctly:

```bash
anvil doctor
```

This command will:

- **✅ Verify anvil initialization** - Complete with real-time progress feedback
- **🔧 Check dependencies** - All required dependencies are installed and functional
- **⚙️ Validate configuration** - Your configuration settings
- **🔗 Test connectivity** - External services (if configured)
- **📊 Show live progress** - Indicators so you know exactly what's happening


If you see any issues, the doctor command will provide specific fix recommendations. Use `anvil doctor --verbose` for detailed troubleshooting information.

### 3. Explore Available Commands

Get familiar with Anvil's capabilities:

```bash
# See all available commands
anvil --help

# Get help for specific commands
anvil init --help
anvil install --help
anvil config --help
anvil config pull --help
anvil config show --help
anvil config sync --help
anvil config push --help
anvil clean --help
anvil doctor --help
```

### 3. Check Your Configuration

View the generated anvil configuration:

```bash
cat ~/.anvil/settings.yaml
```

or simply run:

```bash
anvil config show
```

You'll see something like:

```yaml
version: 1.0.0
directories:
  config: /Users/username/.anvil
tools:
  required_tools: [git, curl]
groups:
  dev: [git, zsh, iterm2, vscode]
  essentials: [slack, chrome, 1password]
git:
  username: "Your Name"
  email: "your.email@example.com"
environment: {}
```

## Basic Usage

### Installing Tool Groups

Anvil organizes tools into logical groups for easy batch installation:

#### Development Tools

Install essential development tools:

```bash
anvil install dev
```

This installs:

- **📝 Git** - Version control
- **🐚 Zsh** - Advanced shell with oh-my-zsh
- **💻 iTerm2** - Enhanced terminal (macOS)
- **🎨 VS Code** - Code editor

#### New Laptop Essentials

Set up a new machine with essential applications:

```bash
anvil install essentials
```

This installs:

- **💬 Slack** - Team communication
- **🌐 Chrome** - Web browser
- **🔐 1Password** - Password manager

#### Preview Before Installing

Use dry-run to see what would be installed:

```bash
anvil install dev --dry-run
anvil install essentials --dry-run
```
### Custom Group Creation

You can create your own custom groups with any app you would like to organize. Group creation can be donee manually directly in the settings.yaml file or by importing existing gropups. See some examples in [import-examples](../import-examples/README.md)

### Installing Individual Applications

Install any application by name with automatic tracking:

```bash
# Install any application available through Homebrew
anvil install git
anvil install firefox
anvil install slack
anvil install visual-studio-code
anvil install figma

# Preview installation
anvil install firefox --dry-run
```

**🎯 Automatic Tracking**: Individual apps are automatically added to `tools.installed_apps` in your settings.yaml.

### Installing Applications with Group Assignment

Install applications and automatically organize them into groups:

```bash
# Add to existing group
anvil install firefox --group-name essentials
anvil install chrome --group-name browsers

# Create new group and add app
anvil install final-cut --group-name editing
anvil install premiere --group-name editing

# Preview with group assignment
anvil install sketch --group-name design --dry-run
```

**🎯 Group Organization**: Apps are added to specified groups instead of `installed_apps`, creating logical collections for easy management.

### How Individual App Installation Works

1. **Dynamic Detection**: Works with any Homebrew package
2. **Smart Tracking**: Automatically added to settings.yaml
3. **Duplicate Prevention**: Won't track apps already in groups or required tools
4. **Existing App Registration**: Works on already-installed apps too

### Listing Available Options

See all available groups and tools:

```bash
anvil install --list
```

## Common Workflows

### Workflow 1: New Developer Machine Setup

Complete setup for a new development machine:

```bash
# Step 1: Initialize Anvil
anvil init

# Step 2: Verify setup is working
anvil doctor

# Step 3: Import coding essentials group
anvil import https://raw.githubusercontent.com/rocajuanma/anvil/master/import-examples/code-essentials.yaml

# Step 4: Install essentials group
anvil install essentials

# Step 5: Verify all installations
anvil doctor dependencies

# Step 6: Add any additional tools as needed
# (Additional tools can be installed through custom groups)
```

### Workflow 2: Team Onboarding

Quickly onboard a new team member:

```bash
# Initialize
anvil init

# Verify environment before proceeding
anvil doctor

# Import backend team groups
anvil import https://raw.githubusercontent.com/rocajuanma/anvil/master/import-examples/backend-developer.yaml

# Install imported group: slack, postman, 1password
anvil install productivity

# Install backend tools
anvil install backend-core

# Final verification
anvil doctor
```

### Workflow 3: Selective Tool Installation

Install only specific tools you need:

```bash
# Initialize first
anvil init

# Preview what you want
anvil install dev --dry-run

# Install selected tools
anvil install dev
```

### Workflow 4: Custom Group Creation

Create your own tool groups by editing `~/.anvil/settings.yaml`:

```yaml
groups:
  frontend:
    - git
    - node
    - yarn
    - chrome
  backend:
    - git
    - docker
    - postgresql
    - redis
  data-science:
    - python
    - jupyter
    - pandas
```

Then use your custom groups:

```bash
anvil install frontend
anvil install backend
```

## Understanding Configuration

### Configuration File Location

Anvil stores its configuration in `~/.anvil/settings.yaml`. This file contains:

- **Directories**: Path for configuration storage
- **Tools**: Lists of required and optional tools
- **Groups**: Tool collections for batch installation
- **Git**: Your Git configuration
- **Environment**: Custom environment variables

### Customizing Groups

Edit `~/.anvil/settings.yaml` to add custom groups:

```yaml
groups:
  my-workflow:
    - git
    - docker
    - vscode
    - slack
```

### Environment Configuration

Add custom environment variables:

```yaml
environment:
  EDITOR: "code"
  DEVELOPER_MODE: "true"
  MY_CUSTOM_PATH: "/usr/local/custom/bin"
```

### Directory Structure

Anvil creates and uses these directories:

```
~/.anvil/
├── settings.yaml    # Main configuration
├── temp/            # Temporary storage for pulled configurations
└── dotfiles/        # Local repository clone (when using config commands)
```

## Configuration Management

Anvil provides powerful configuration management to sync dotfiles and application settings across machines.

### Setting Up Configuration Management

1. **Create a GitHub repository** for your configurations.

2. **Configure Anvil** by editing `~/.anvil/settings.yaml`:

   ```yaml
   github:
     config_repo: "username/dotfiles" # Your GitHub repository
     branch: "main" # Branch to use
     local_path: "~/.anvil/dotfiles" # Local storage path
     token_env_var: "GITHUB_TOKEN" # Environment variable for token

   git:
     username: "Your Name"
     email: "your.email@example.com"
     ssh_key_path: "~/.ssh/id_ed25519"
   ```

3. **Set up authentication** (choose one):
Anvil uses your existing `.ssh` configs for authentication. The `init` command should configure everything in your `settings.yaml` automatically.

If you run into problems or need to update those values(because of an ssh key update, etc), just run `anvil doctor git-config --fix`

### Pulling Configurations

Pull specific configuration directories from your repository:

```bash
# Pull Anvil settings(if present in remote repo)
anvil config pull

# Pull Cursor editor configurations
anvil config pull cursor

# Pull VS Code configurations
anvil config pull vs-code

# Pull shell configurations
anvil config pull zsh

# View pulled configurations
anvil config show cursor
anvil config show vs-code

# Sync pulled configs to their destinations
anvil config sync              # Apply pulled anvil settings
anvil config sync cursor      # Apply pulled app configs
anvil config sync --dry-run   # Preview sync changes
```

**Current Behavior**: Always fetches the latest changes from your repository and pulls files to `~/.anvil/temp/[directory]` for review before manual application.

### Configuration Push

The `anvil config push` command allows you to upload your anvil configuration changes back to your repository:

```bash
# Push anvil settings to repository
anvil config push

# The command will:
# 1. Compare local and remote configurations
# 2. Create timestamped branch if changes exist
# 3. Commit changes with standardized message
# 4. Push branch and provide PR link
```

**Note**: Application-specific configuration push (e.g., `anvil config push <app-name>`) is also fully supported.

### For More Details

📖 **[Complete Configuration Guide](config.md)** - Detailed setup instructions, troubleshooting, and advanced usage examples.

## Tips and Best Practices

### 🎯 Initialization Best Practices

1. **Always run `anvil init` first** on any new machine
2. **Review the output** and follow any recommendations
3. **Complete environment setup** before installing tools

### 🔧 Tool Installation Best Practices

1. **Use dry-run first** to preview installations:

   ```bash
   anvil install dev --dry-run
   ```

2. **Start with groups**, then add individual tools:

   ```bash
   anvil install dev
   anvil install docker && anvil install kubectl
   ```

3. **Check available options** before installing:
   ```bash
   anvil install --list
   ```

### 📋 Configuration Best Practices

1. **Backup your configuration** before making changes:

   ```bash
   cp ~/.anvil/settings.yaml ~/.anvil/settings.yaml.backup
   ```
  
  or push to remote to have a backup in repo
  ```bash
  anvil config push
  ```

2. **Use descriptive names** for custom groups
3. **Keep groups focused** - don't make them too large
4. **Document custom groups** for team members

### 🚀 Team Usage Best Practices

1. **Share configurations** across team members
2. **Create team-specific groups** in settings.yaml
3. **Document your team's setup process**
4. **Use consistent tool versions** when possible

### 🧹 Maintenance Best Practices

1. **Regular cleanup** with `anvil clean` to free disk space
2. **Clean before major operations** to ensure fresh state
3. **Use dry-run mode** to preview cleanups with `anvil clean --dry-run`
4. **Maintain configurations** by cleaning old archives and temporary files

## Troubleshooting

### First Step: Run Health Check

When experiencing any issues with Anvil, your first step should always be to run the diagnostic command:

```bash
anvil doctor
```

This will automatically detect and report:

- ✅ Environment setup issues
- ✅ Missing or broken dependencies
- ✅ Configuration problems
- ✅ Connectivity issues

**For specific areas:**

```bash
# Check only environment setup
anvil doctor --category environment

# Check only dependencies
anvil doctor --category dependencies

# Check only configuration
anvil doctor --category configuration

# Auto-fix detected issues
anvil doctor --fix
```

### Common Issues and Solutions

#### Command Not Found

```bash
# Check if anvil is in PATH
which anvil

# If not found, add to PATH or move to /usr/local/bin/
export PATH=$PATH:/path/to/anvil
# or
sudo mv anvil /usr/local/bin/
```

#### Permission Errors

```bash
# Fix common permission issues
sudo chown -R $(whoami) ~/.anvil
chmod 755 ~/.anvil
```

#### Tool Installation Failures

```bash
# Run health check first
anvil doctor --category dependencies

# Update Homebrew (macOS)
brew update

# Check internet connectivity
ping -c 3 github.com

# Try individual installation to isolate issues
anvil install git --dry-run
anvil install git
```

#### Configuration Issues

```bash
# Check configuration health
anvil doctor --category configuration

# Reinitialize if configuration is corrupted
rm -rf ~/.anvil
anvil init
```

#### Homebrew Issues (macOS)

```bash
# Fix Homebrew permissions
sudo chown -R $(whoami) $(brew --prefix)/*

# Reinstall Homebrew if needed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Platform-Specific Issues

#### macOS

- **Security warnings**: Allow in System Preferences → Security & Privacy
- **Xcode tools missing**: Run `xcode-select --install`
- **Homebrew PATH issues**: Add `/opt/homebrew/bin` to PATH

### Getting Help

If you encounter issues not covered here:

1. **Check the configuration** in `~/.anvil/settings.yaml`
2. **Search existing issues** on [GitHub](https://github.com/rocajuanma/anvil/issues)
3. **Create a new issue** with:
   - Your platform and version
   - Command you ran
   - Complete error message
   - Output of `anvil --version`

## Next Steps

Now that you're familiar with the basics:

### Explore Advanced Features

- **[Install Command Documentation](install.md)** - Deep dive into tool installation
- **[Init Command Documentation](init.md)** - Understand initialization process
- **[Clean Command Documentation](clean.md)** - Directory cleanup and maintenance
- **[Examples and Tutorials](EXAMPLES.md)** - Real-world usage scenarios

### Customize Your Setup

- **Edit `~/.anvil/settings.yaml`** to create custom tool groups
- **Add environment variables** for your workflow
- **Share configurations** with your team

### Contribute

- **[Contributing Guide](CONTRIBUTING.md)** - Help improve Anvil
- **Report bugs** or **request features** on GitHub

### Stay Updated

- **Use `anvil update`** for automatic updates (v1.2.0+)
- **Watch the repository** for updates
- **Read the [changelog](CHANGELOG.md)** for new features
- **Join discussions** for community support

---

## Quick Reference

### Essential Commands

```bash
anvil init                      # Initialize Anvil
anvil install --list            # List available groups and tools
anvil install dev                # Install development tools
anvil install git && anvil install zsh  # Install specific tools individually
anvil install dev --dry-run     # Preview installations
anvil config pull cursor      # Pull Cursor configurations from remote
anvil config pull vs-code     # Pull VS Code configurations from remote
anvil config show cursor      # View pulled Cursor configurations
anvil config show            # View anvil settings.yaml
anvil config sync            # Apply pulled anvil settings
anvil config sync cursor     # Apply pulled app configs
anvil config sync --dry-run  # Preview sync changes
anvil config push            # Push anvil settings to remote
anvil config push cursor      # Push app configurations (in development)
anvil clean                  # Clean anvil directories and temp files
anvil clean --dry-run        # Preview cleanup without deletion
anvil --help                   # Get help
```

### Key Files

- `~/.anvil/settings.yaml` - Main configuration

### Important Locations

- **Configuration**: `~/.anvil/`
- **Documentation**: `docs/` directory

