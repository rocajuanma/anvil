# Anvil Examples & Tutorials

This document provides real-world examples and tutorials for using Anvil CLI effectively. From basic scenarios to advanced team setups, these examples will help you master Anvil.

## Table of Contents

- [Basic Examples](#basic-examples)
- [Development Workflows](#development-workflows)
- [Team Scenarios](#team-scenarios)
- [Advanced Configurations](#advanced-configurations)
- [Platform-Specific Examples](#platform-specific-examples)
- [Troubleshooting Examples](#troubleshooting-examples)
- [Integration Examples](#integration-examples)

## Basic Examples

### Example 1: First-Time Setup

**Scenario**: You just installed Anvil and want to set up your development environment.

**Key Points:**
- **🚀 Complete initialization** - Sets up your Anvil environment from scratch
- **🔍 Health verification** - Ensures everything is working correctly
- **📦 Tool installation** - Installs essential development tools

```bash
# Step 1: Initialize Anvil
$ anvil init

=== Anvil Initialization ===
🔧 Validating and installing required tools...
✅ All required tools are available
🔧 Creating necessary directories...
✅ Directories created successfully
🔧 Generating default settings.yaml...
✅ Default settings.yaml generated

# Step 2: Verify setup is working correctly
$ anvil doctor

=== Running Anvil Health Check ===

✅ Environment
  ✅ Anvil initialization complete
  ✅ Settings file is valid
  ✅ Directory structure is correct

✅ Dependencies
  ✅ Homebrew is installed and functional
  ✅ All required tools installed (2/2)

✅ Overall status: Healthy

# Step 3: See what tools are available
$ anvil install --list

=== Available Setup Groups ===
Group: dev
  • git
  • zsh
  • iterm2
  • vscode

Group: new-laptop
  • slack
  • chrome
  • 1password

# Step 4: Install development tools
$ anvil install dev

=== Setting up 'dev' group ===
Installing tools for group 'dev': git, zsh, iterm2, vscode
[1/4] 25% - Installing git
✅ git installed successfully
[2/4] 50% - Installing zsh
Installing oh-my-zsh...
✅ zsh installed successfully
[3/4] 75% - Installing iterm2
✅ iterm2 installed successfully
[4/4] 100% - Installing vscode
✅ vscode installed successfully

=== Group Setup Complete! ===
Successfully installed 4 of 4 tools in group 'dev'
```

### Example 2: Selective Tool Installation

**Scenario**: You only need specific tools, not entire groups.

**Key Points:**
- **🔍 Preview capability** - See what would be installed before proceeding
- **🎯 Selective installation** - Choose only the tools you need
- **📝 Smart tracking** - Individual tools are automatically tracked

```bash
# Preview what would be installed
$ anvil install dev --dry-run

Dry run mode - no actual installations will be performed

=== Individual Tool Setup ===
Installing individual tools: git, zsh
[1/2] 50% - Installing git
Would install: git
[2/2] 100% - Installing zsh
Would install: zsh

# Install the tools
$ anvil install dev

=== Individual Tool Setup ===
Installing individual tools: git, zsh
[1/2] 50% - Installing git
✅ git installed successfully
[2/2] 100% - Installing zsh
Installing oh-my-zsh...
✅ zsh installed successfully

=== Individual Tool Setup Complete! ===
```

### Example 3: New Laptop Setup

**Scenario**: Setting up essential applications on a new machine.

**Key Points:**
- **💻 Essential apps** - Communication, browser, and security tools
- **🚀 Quick setup** - Get productive immediately on new hardware
- **🔧 Team consistency** - Standard applications across team members

```bash
# Initialize first
$ anvil init

# Install essential applications
$ anvil install new-laptop

=== Setting up 'new-laptop' group ===
Installing tools for group 'new-laptop': slack, chrome, 1password
[1/3] 33% - Installing slack
✅ slack installed successfully
[2/3] 67% - Installing chrome
✅ chrome installed successfully
[3/3] 100% - Installing 1password
✅ 1password installed successfully

=== Group Setup Complete! ===
Successfully installed 3 of 3 tools in group 'new-laptop'
```

### Example 3.5: Organizing Applications into Groups

**Scenario**: Installing applications and organizing them into logical groups for better management.

**Key Points:**
- **🗂️ Logical organization** - Group related applications together
- **📋 Better management** - Easier to track and maintain tool collections
- **🎯 Targeted installation** - Install entire groups with one command

```bash
# Initialize first
$ anvil init

# Install and organize browsers
$ anvil install firefox --group-name browsers
✅ firefox installed successfully
✅ Added firefox to group 'browsers'

$ anvil install chrome --group-name browsers
✅ chrome installed successfully
✅ Added chrome to group 'browsers'

# Install and organize design tools
$ anvil install figma --group-name design
✅ figma installed successfully
✅ Added figma to group 'design'

$ anvil install sketch --group-name design
✅ sketch installed successfully
✅ Added sketch to group 'design'

# Install and organize development tools
$ anvil install docker --group-name devops
✅ docker installed successfully
✅ Added docker to group 'devops'

$ anvil install kubernetes-cli --group-name devops
✅ kubernetes-cli installed successfully
✅ Added kubernetes-cli to group 'devops'

# Verify your organized groups
$ anvil install --list

=== Available Groups ===
Built-in Groups:
  • dev: git, zsh, iterm2, visual-studio-code
  • new-laptop: slack, google-chrome, 1password

Custom Groups:
  • browsers: firefox, chrome
  • design: figma, sketch
  • devops: docker, kubernetes-cli

Individually Tracked Apps:
  (none - all apps are organized in groups)
```

### Example 3.6: Environment Maintenance and Cleanup

**Scenario**: Regular maintenance of your Anvil environment to free disk space and organize files.

**Key Points:**
- **🧹 Disk space management** - Remove temporary files and old archives
- **🔍 Preview changes** - See what will be cleaned before deletion
- **🛡️ Safe operations** - Settings and directory structure are preserved

```bash
# Preview what would be cleaned
$ anvil clean --dry-run

=== Cleaning Anvil Directories ===

🔧 Scanning .anvil directory for content to clean

Found 3 root directories to clean:
Directory structure to be cleaned:
  📁 temp (2)
    ├── cursor
    ├── vscode
  📁 archive (1)
    ├── 2025-01-15-1430
  📁 dotfiles (5)
    ├── .git
    ├── cursor
    ├── vscode
    ├── README.md
    ├── .gitignore

DRY RUN: Would clean contents of 3 root directories

# Perform the actual cleanup
$ anvil clean

=== Cleaning Anvil Directories ===

🔧 Scanning .anvil directory for content to clean
🔧 Cleaning directories and files

✅ Cleaned contents of directory temp
✅ Cleaned contents of directory archive  
✅ Removed dotfiles directory completely

Successfully cleaned contents of 3/3 root directories
```

## Development Workflows

### Example 4: Frontend Developer Setup

**Scenario**: Setting up a machine for frontend development.

```bash
# Initialize Anvil
$ anvil init

# Install core development tools
$ anvil install dev

# Add frontend-specific tools
$ anvil install chrome

# Verify installation
$ git --version
$ code --version
$ zsh --version
```

**Custom Configuration**: Edit `~/.anvil/settings.yaml` to add a frontend group:

```yaml
groups:
  frontend:
    - git
    - node
    - yarn
    - vscode
    - chrome
    - figma
```

Then use it:

```bash
$ anvil install frontend
```

### Example 5: Backend Developer Setup

**Scenario**: Setting up for backend development with containers.

```bash
# Initialize
$ anvil init

# Install development tools
$ anvil install dev

# Add Docker through custom group configuration
# (See settings.yaml custom group setup)

# Create custom backend group in settings.yaml
```

**Custom settings.yaml addition**:

```yaml
groups:
  backend:
    - git
    - docker
    - vscode
    - postgresql
    - redis
    - kubectl
```

**Usage**:

```bash
$ anvil install backend
```

### Example 6: Full-Stack Developer Setup

**Scenario**: Complete setup for full-stack development.

```bash
# Initialize
$ anvil init

# Install all development tools
$ anvil install dev

# Add communication tools
$ anvil install slack

# Add browser for testing
$ anvil install chrome

# Optional: Add custom tools through group configuration
# (See settings.yaml for custom group setup)
```

## Team Scenarios

### Example 7: Team Onboarding Script

**Scenario**: Create a script for onboarding new team members.

**File**: `team-setup.sh`

```bash
#!/bin/bash
echo "🚀 Welcome to the team! Setting up your development environment..."

# Initialize Anvil
echo "📋 Initializing Anvil..."
anvil init

# Install core development tools
echo "🔧 Installing development tools..."
anvil install dev

# Install team communication tools
echo "💬 Installing communication tools..."
anvil install slack

# Install additional tools through custom groups
echo "🔧 Installing additional tools..."
# (Define custom groups in settings.yaml for additional tools)

# Custom team tools
echo "🛠️  Installing team-specific tools..."
anvil install --figma --postman

echo "✅ Setup complete! Welcome to the team!"
echo "📚 Next steps:"
echo "  1. Configure your Git credentials: git config --global user.name 'Your Name'"
echo "  2. Set up SSH keys for GitHub"
echo "  3. Join our Slack workspace"
echo "  4. Clone team repositories"
```

**Usage**:

```bash
chmod +x team-setup.sh
./team-setup.sh
```

### Example 8: Team Configuration Sharing

**Scenario**: Share a standard configuration across team members.

**File**: `team-settings.yaml`

```yaml
version: 1.0.0
directories:
  config: ~/.anvil
  cache: ~/.anvil/cache
  data: ~/.anvil/data
tools:
  required_tools: [git, curl]
  optional_tools: [brew, docker, kubectl]
groups:
  dev: [git, zsh, iterm2, vscode]
  new-laptop: [slack, chrome, 1password]
  frontend:
    - git
    - node
    - yarn
    - vscode
    - chrome
    - figma
  backend:
    - git
    - docker
    - vscode
    - postgresql
    - redis
  qa:
    - git
    - chrome
    - postman
    - cypress
git:
  username: ""
  email: ""
environment:
  EDITOR: "code"
  TEAM_NAME: "awesome-team"
```

**Distribution**:

```bash
# Share via Git repository
git clone https://github.com/company/anvil-config.git
cp anvil-config/team-settings.yaml ~/.anvil/settings.yaml

# Or via direct download
curl -o ~/.anvil/settings.yaml https://company.com/anvil/team-settings.yaml

# Then install team tools
anvil install frontend  # or backend, qa
```

### Example 9: Multi-Platform Team Setup

**Scenario**: Team with members on different platforms.

**macOS team member**:

```bash
anvil init
anvil install dev
anvil install --slack
```

**Linux team member**:

```bash
anvil init
# May see platform warnings
anvil install --git --vscode
# Install additional tools manually or via platform package manager
```

**Windows team member**:

```bash
# Use WSL or Git Bash
anvil init
anvil install --git --vscode
# Install Windows-specific alternatives manually
```

## Advanced Configurations

### Configuration Management

**Scenario**: You want to sync your dotfiles and application configurations across multiple machines.

#### Setting Up Configuration Management

**Step 1**: Create a GitHub repository with organized directories

```bash
# Create a new repository for your configs
mkdir ~/my-dotfiles
cd ~/my-dotfiles
git init
git remote add origin https://github.com/yourusername/dotfiles.git

# Create organized directory structure
mkdir -p cursor vs-code zsh git ssh
echo "# My Configuration Repository" > README.md

# Add your configuration files by application
cp ~/Library/Application\ Support/Cursor/User/settings.json cursor/
cp ~/Library/Application\ Support/Cursor/User/keybindings.json cursor/

cp ~/.zshrc zsh/
cp ~/.zsh_aliases zsh/
cp ~/.zsh_functions zsh/

cp ~/.gitconfig git/
cp ~/.gitignore_global git/

cp ~/.ssh/config ssh/

git add .
git commit -m "Initial configuration backup"
git push -u origin main
```

**Step 2**: Configure Anvil for config management

```bash
# Initialize Anvil on your main machine
anvil init

# Edit ~/.anvil/settings.yaml to add GitHub configuration:
```

```yaml
github:
  config_repo: "yourusername/dotfiles" # Your GitHub repository
  branch: "main" # Branch to use
  local_path: "~/.anvil/dotfiles" # Local storage path
  token_env_var: "GITHUB_TOKEN" # Environment variable for token

git:
  username: "Your Name"
  email: "your.email@example.com"
  ssh_key_path: "~/.ssh/id_ed25519"
```

**Step 3**: Set up authentication

```bash
# Option A: SSH Key (Recommended)
# Add your SSH key to GitHub if you haven't already
cat ~/.ssh/id_ed25519.pub | pbcopy
# Go to GitHub > Settings > SSH Keys and add the key

# Option B: GitHub Token
# Create token at github.com/settings/tokens
export GITHUB_TOKEN="your_token_here"
echo 'export GITHUB_TOKEN="your_token_here"' >> ~/.zshrc
```

**Step 4**: Test configuration pulling

```bash
# Pull specific configuration directories
anvil config pull cursor
anvil config pull zsh
anvil config pull git
```

#### Applying Configurations on New Machines

**Complete setup example on a new machine:**

```bash
# 1. Install and initialize Anvil
go install github.com/yourusername/anvil@latest
anvil init

# 2. Configure GitHub repository
# Edit ~/.anvil/settings.yaml with your repository details

# 3. Pull specific configurations as needed
anvil config pull cursor

# Example output:
# 🔧 Using branch: main
#
# === Pulling Configuration Directory: cursor ===
#
# Repository: yourusername/dotfiles
# Branch: main
# Target directory: cursor
# ✅ GitHub token found in environment variable: GITHUB_TOKEN
# 🔧 Validating repository access and branch configuration...
# ✅ Repository and branch configuration validated
# 🔧 Setting up local repository...
# ✅ Local repository ready
# 🔧 Pulling latest changes...
# ✅ Repository updated
# 🔧 Copying configuration directory...
# ✅ Configuration directory copied to temp location
#
# === Pull Complete! ===
#
# Configuration directory 'cursor' has been pulled from: yourusername/dotfiles
# Files are available at: /Users/username/.anvil/temp/cursor
#
# Copied files:
#   • settings.json
#   • keybindings.json

# 4. Review and apply configurations manually
ls ~/.anvil/temp/cursor/
cp ~/.anvil/temp/cursor/settings.json ~/Library/Application\ Support/Cursor/User/
cp ~/.anvil/temp/cursor/keybindings.json ~/Library/Application\ Support/Cursor/User/

# 5. Pull other configurations as needed
anvil config pull vs-code
anvil config pull zsh

# 6. Review configurations before applying
anvil config show vs-code
anvil config show zsh

# 7. Apply pulled configurations
anvil config sync --dry-run     # Preview changes
anvil config sync               # Apply anvil settings
anvil config sync vs-code       # Apply VS Code configs
anvil config sync zsh           # Apply zsh configs
```

#### Current vs. Future Implementation

**Current Behavior**: Always fetches the latest changes from your repository and pulls configuration files to `~/.anvil/temp/[directory]` for manual review and application.

**Future Enhancement**: Anvil will automatically detect application configuration paths and apply configurations directly with backup and rollback capabilities.

### Team Configuration Sharing

**Scenario**: Share team configurations and development standards across team members.

#### Team Setup Example

**Team Lead Setup**:

```bash
# 1. Create team configuration repository
mkdir ~/team-configs
cd ~/team-configs

# 2. Add team-specific configurations
mkdir -p vs-code cursor git zsh homebrew

# VS Code team settings
cat > vs-code/settings.json << 'EOF'
{
  "editor.tabSize": 2,
  "editor.insertSpaces": true,
  "eslint.enable": true,
  "prettier.requireConfig": true,
  "workbench.colorTheme": "Dark+ (default dark)"
}
EOF

# Git team settings
cat > git/.gitconfig << 'EOF'
[core]
    editor = code --wait
    autocrlf = input
[pull]
    rebase = true
[init]
    defaultBranch = main
EOF

git init
git add .
git commit -m "Initial team configuration"
git remote add origin https://github.com/company/team-configs.git
git push -u origin main
```

**Team Member Setup**:

```bash
# 1. Configure Anvil with team repository
anvil init

# Edit ~/.anvil/settings.yaml:
```

```yaml
github:
  config_repo: "company/team-configs"
  branch: "main"
  local_path: "~/.anvil/team-configs"
  token_env_var: "GITHUB_TOKEN"
```

```bash
# 2. Pull team configurations
anvil config pull vs-code
anvil config pull git
anvil config pull zsh

# 3. Review and apply team standards
# Files are now in ~/.anvil/temp/ directories for review
ls ~/.anvil/temp/vs-code/
ls ~/.anvil/temp/git/
```

#### Advanced Team Workflow

```bash
# Team lead updates configurations
# (edit files in team repository)
git add .
git commit -m "Update team ESLint settings"
git push

# Team members pull latest updates
anvil config pull vs-code

# Output shows what changed:
# === Pull Complete! ===
#
# Configuration directory 'vs-code' has been pulled from: company/team-configs
# Files are available at: /Users/username/.anvil/temp/vs-code
#
# Copied files:
#   • settings.json (updated)
#   • extensions.json
#   • keybindings.json
```

### Configuration Repository Organization Examples

#### Personal Developer Setup

```
dotfiles/
├── editors/
│   ├── cursor/
│   │   ├── settings.json
│   │   ├── keybindings.json
│   │   └── snippets/
│   │       ├── javascript.json
│   │       └── typescript.json
│   └── vs-code/
│       ├── settings.json
│       ├── extensions.json
│       └── keybindings.json
├── shell/
│   ├── zsh/
│   │   ├── .zshrc
│   │   ├── .zsh_aliases
│   │   └── .zsh_functions
│   └── bash/
│       └── .bashrc
├── git/
│   ├── .gitconfig
│   └── .gitignore_global
├── ssh/
│   └── config
└── README.md
```

Pull examples:

```bash
anvil config pull editors/cursor
anvil config pull shell/zsh
anvil config pull git
```

#### Team/Company Setup

```
team-configs/
├── frontend/
│   ├── vs-code/
│   │   ├── settings.json        # Frontend team VS Code settings
│   │   └── extensions.json
│   └── cursor/
│       └── settings.json        # Frontend team Cursor settings
├── backend/
│   ├── vs-code/
│   │   └── settings.json        # Backend team VS Code settings
│   └── git/
│       └── .gitconfig           # Backend-specific Git settings
├── shared/
│   ├── git/
│   │   ├── .gitconfig           # Company-wide Git settings
│   │   └── .gitignore_global
│   └── zsh/
│       └── .zshrc               # Company shell configuration
└── README.md
```

Pull examples:

```bash
# Frontend developer
anvil config pull frontend/vs-code
anvil config pull shared/git

# Backend developer
anvil config pull backend/vs-code
anvil config pull shared/git
anvil config pull shared/zsh
```

### Example 10: Custom Tool Groups

**Scenario**: Create specialized tool groups for different projects.

**Configuration**: `~/.anvil/settings.yaml`

```yaml
groups:
  # Mobile development
  mobile:
    - git
    - vscode
    - android-studio
    - xcode
    - flutter

  # Data science
  data-science:
    - git
    - python
    - jupyter
    - vscode
    - r
    - rstudio

  # DevOps
  devops:
    - git
    - docker
    - kubectl
    - terraform
    - aws-cli
    - vscode

  # Design
  design:
    - git
    - figma
    - sketch
    - adobe-creative-cloud
    - chrome
```

**Usage**:

```bash
# Install specific workflow tools
anvil install mobile
anvil install data-science
anvil install devops
anvil install design
```

### Example 11: Environment-Specific Setup

**Scenario**: Different setups for different environments.

**Development Environment**:

```bash
anvil init
anvil install dev
# Additional tools through custom groups
```

**Production Environment**:

```bash
anvil init
anvil install --git  # Minimal tools only
# Production-specific tools via other means
```

**Testing Environment**:

```bash
anvil init
anvil install --git --chrome
anvil install --cypress --postman
```

### Example 12: Project-Specific Tool Installation

**Scenario**: Install tools based on project requirements.

**React Project**:

```bash
# Create project-specific group
cat >> ~/.anvil/settings.yaml << EOF
    react-project:
      - git
      - node
      - yarn
      - vscode
      - chrome
      - react-devtools
EOF

anvil install react-project
```

**Python Project**:

```bash
# Add Python project group
cat >> ~/.anvil/settings.yaml << EOF
    python-project:
      - git
      - python
      - pip
      - vscode
      - jupyter
EOF

anvil install python-project
```

## Platform-Specific Examples

### Example 13: macOS Optimized Setup

**Scenario**: Take advantage of macOS-specific tools.

```bash
# Full macOS setup
anvil init  # Installs Homebrew automatically
anvil install dev  # Includes iTerm2
anvil install new-laptop

# Add macOS-specific productivity tools
anvil install --alfred --spectacle --raycast
```

### Example 14: Linux Development Setup

**Scenario**: Setting up on Ubuntu/Debian.

```bash
# Initialize (may need to install dependencies first)
sudo apt update
sudo apt install -y git curl build-essential

anvil init
anvil install --git --vscode

# Linux-specific additions
sudo apt install -y zsh vim tmux
```

### Example 15: Windows with WSL

**Scenario**: Windows developer using WSL.

```bash
# In WSL
anvil init
anvil install --git --vscode

# Install Windows VS Code extension for WSL integration
# Some tools need Windows-specific installation
```

## Troubleshooting Examples

### Example 16: Using Health Checks for Diagnostics

**Scenario**: Something isn't working and you need to diagnose the issue.

```bash
# Always start with a comprehensive health check
$ anvil doctor

=== Running Anvil Health Check ===

✅ Environment
  ✅ Anvil initialization complete
  ✅ Settings file is valid
  ✅ Directory structure is correct

⚠️ Dependencies
  ⚠️ Homebrew has available updates
      💡 Homebrew will be updated
  ✅ All required tools installed (2/2)

❌ Configuration
  ❌ Git configuration incomplete: username not set
      💡 Git configuration must be set manually in settings.yaml

✅ Connectivity
  ✅ GitHub SSH access confirmed
  ✅ Git operations are functional

=== Health Check Summary ===
Total checks: 12
✅ Passed: 9  ⚠️ Warnings: 1  ❌ Failed: 1

🔧 1 issues can be auto-fixed
Run 'anvil doctor --fix' to automatically fix them

# Fix auto-fixable issues
$ anvil doctor --fix
Successfully fixed: homebrew

# Check specific categories when targeting issues
$ anvil doctor configuration
$ anvil doctor dependencies

# Check individual validators
$ anvil doctor git-config
```

### Example 17: Fixing Installation Failures

**Scenario**: Some tools fail to install.

```bash
# First, run health check to see if it's a dependency issue
anvil doctor dependencies

# Check what failed
anvil install dev

# Output shows:
# ❌ Failed to install vscode: package not found

# Debug individual tool
anvil install --vscode --dry-run

# Check Homebrew
brew update
brew search visual-studio-code

# Manual installation if needed
brew install --cask visual-studio-code

# Verify
code --version
```

### Example 18: Recovering from Bad Configuration

**Scenario**: Configuration file got corrupted.

```bash
# Run health check to see what's wrong
anvil doctor

# If configuration is severely broken:
# Backup current config
cp ~/.anvil/settings.yaml ~/.anvil/settings.yaml.backup

# Reset to defaults
rm -rf ~/.anvil
anvil init

# Verify new setup
anvil doctor

# Restore custom groups from backup if needed
# Edit ~/.anvil/settings.yaml to add custom configurations
```

### Example 18: Permission Issues

**Scenario**: Permission errors during installation.

```bash
# Fix Homebrew permissions (macOS)
sudo chown -R $(whoami) $(brew --prefix)/*

# Fix Anvil directory permissions
sudo chown -R $(whoami) ~/.anvil
chmod 755 ~/.anvil

# Retry installation
anvil install dev
```

## Integration Examples

### Example 19: CI/CD Pipeline Integration

**Scenario**: Use Anvil in automated environments.

**File**: `.github/workflows/setup.yml`

```yaml
name: Development Environment Setup
on: [push, pull_request]

jobs:
  setup:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v2

      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17

      - name: Build Anvil
        run: go build -o anvil main.go

      - name: Initialize Anvil
        run: ./anvil init

      - name: Install Development Tools
        run: ./anvil install dev --dry-run # Dry run in CI
```

### Example 20: Docker Container with Anvil

**Scenario**: Create a development container with Anvil pre-installed.

**File**: `Dockerfile`

```dockerfile
FROM golang:1.17-alpine

# Install dependencies
RUN apk add --no-cache git curl bash

# Install Anvil
WORKDIR /tmp
COPY . .
RUN go build -o /usr/local/bin/anvil main.go

# Setup development user
RUN adduser -D -s /bin/bash developer
USER developer
WORKDIR /home/developer

# Initialize Anvil
RUN anvil init

# Default command
CMD ["/bin/bash"]
```

**Usage**:

```bash
# Build container
docker build -t anvil-dev .

# Run with Anvil ready
docker run -it anvil-dev

# Inside container
anvil install --git --vscode
```

### Example 21: Makefile Integration

**Scenario**: Integrate Anvil into project Makefile.

**File**: `Makefile`

```makefile
.PHONY: setup-dev setup-team clean-anvil

setup-dev:
	@echo "Setting up development environment..."
	anvil init
	anvil install dev
	@echo "Development setup complete!"

setup-team:
	@echo "Setting up team environment..."
	anvil init
	anvil install dev
	anvil install --slack
	@echo "Team setup complete!"

clean-anvil:
	@echo "Cleaning Anvil configuration..."
	rm -rf ~/.anvil
	@echo "Anvil configuration removed."

install-anvil:
	@echo "Installing Anvil..."
	go build -o anvil main.go
	sudo mv anvil /usr/local/bin/
	@echo "Anvil installed to /usr/local/bin/"
```

**Usage**:

```bash
make install-anvil
make setup-dev
# or
make setup-team
```

## Advanced Scripting Examples

### Example 22: Interactive Setup Script

**Scenario**: Create an interactive setup script.

**File**: `interactive-setup.sh`

```bash
#!/bin/bash

echo "🚀 Anvil Interactive Setup"
echo "=========================="

# Check if Anvil is installed
if ! command -v anvil &> /dev/null; then
    echo "❌ Anvil not found. Please install first."
    exit 1
fi

# Initialize Anvil
echo "📋 Initializing Anvil..."
anvil init

# Ask user what they want to install
echo ""
echo "What type of development do you do?"
echo "1) Frontend Development"
echo "2) Backend Development"
echo "3) Full-Stack Development"
echo "4) Mobile Development"
echo "5) Data Science"
echo "6) Custom Selection"

read -p "Enter your choice (1-6): " choice

case $choice in
    1)
        echo "🎨 Setting up Frontend Development..."
        anvil install dev
        anvil install --chrome
        ;;
    2)
        echo "⚙️ Setting up Backend Development..."
        anvil install dev
        # Additional tools through custom groups
        ;;
    3)
        echo "🔄 Setting up Full-Stack Development..."
        anvil install dev
        anvil install --chrome
        ;;
    4)
        echo "📱 Setting up Mobile Development..."
        anvil install --git --vscode
        echo "Note: Install Xcode and Android Studio manually"
        ;;
    5)
        echo "📊 Setting up Data Science..."
        anvil install --git --vscode --python
        ;;
    6)
        echo "🛠️ Custom tool selection..."
        anvil install --list
        echo "Use 'anvil install --tool1 --tool2' to install specific tools"
        ;;
    *)
        echo "Invalid choice. Running basic setup..."
        anvil install dev
        ;;
esac

echo ""
echo "✅ Setup complete!"
echo "📚 Next steps:"
echo "  - Configure Git: git config --global user.name 'Your Name'"
echo "  - Configure Git: git config --global user.email 'you@example.com'"
echo "  - Check installed tools: anvil install --list"
```

### Example 23: Bulk Configuration Sync

**Scenario**: Apply multiple pulled configurations to their local destinations.

**File**: `bulk-config-sync.sh`

```bash
#!/bin/bash

echo "📝 Setting up configuration paths in settings.yaml..."

# Configure local paths for applications
cat >> ~/.anvil/settings.yaml << 'EOF'
configs:
  cursor: "~/Library/Application Support/Cursor"
  obsidian: "~/.config/obsidian"
  zsh: "~/.config/zsh"
  neovim: "~/.config/nvim"
EOF

echo "⬇️ Pulling configurations from repository..."
anvil config pull cursor
anvil config pull obsidian
anvil config pull zsh
anvil config pull neovim

echo "🔍 Preview what will be applied..."
anvil config sync --dry-run

echo "🚀 Apply configurations..."
anvil config sync cursor
anvil config sync obsidian
anvil config sync zsh
anvil config sync neovim

echo "✅ Bulk configuration sync complete!"
anvil install --list | grep "Individually Tracked Apps"
```

**Use Case**: Perfect for setting up development environments or onboarding new team members.

### Example 24: Backup and Restore Configuration

**Scenario**: Backup and restore Anvil configurations.

**File**: `backup-config.sh`

```bash
#!/bin/bash

BACKUP_DIR="$HOME/anvil-backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/anvil-config-$TIMESTAMP.tar.gz"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup configuration
echo "📦 Backing up Anvil configuration..."
tar -czf "$BACKUP_FILE" -C "$HOME" .anvil/

echo "✅ Backup created: $BACKUP_FILE"
echo "📝 Backup includes:"
echo "  - Configuration files"
echo "  - Configuration files"
```

**File**: `restore-config.sh`

```bash
#!/bin/bash

BACKUP_DIR="$HOME/anvil-backups"

# List available backups
echo "📋 Available backups:"
ls -1 "$BACKUP_DIR"/anvil-config-*.tar.gz 2>/dev/null | sort -r

# Ask user to select backup
read -p "Enter backup filename to restore: " backup_file

if [[ -f "$BACKUP_DIR/$backup_file" ]]; then
    echo "🔄 Restoring from $backup_file..."

    # Remove current config
    rm -rf ~/.anvil

    # Restore backup
    tar -xzf "$BACKUP_DIR/$backup_file" -C "$HOME"

    echo "✅ Configuration restored!"
else
    echo "❌ Backup file not found!"
    exit 1
fi
```

---

## Quick Reference

### Common Command Patterns

```bash
# Basic setup
anvil init && anvil install dev

# Preview before install
anvil install dev --dry-run

# Individual tools
anvil install --git --zsh --vscode

# List options
anvil install --list

# Help
anvil --help
anvil install --help
```

### Configuration Locations

- **Main config**: `~/.anvil/settings.yaml`
- **Configuration**: `~/.anvil/settings.yaml`

### Useful File Locations

- **Team configs**: Share `settings.yaml` files
- **Scripts**: Store in project root or `~/bin/`
- **Backups**: Use `~/anvil-backups/` or similar

---

**Need more examples?** Check our [GitHub Discussions](https://github.com/rocajuanma/anvil/discussions) or contribute your own examples!
