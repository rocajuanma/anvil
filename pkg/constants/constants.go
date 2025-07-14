package constants

import "fmt"

// AnvilError represents a structured error with operation and command context
type AnvilError struct {
	Op      string
	Command string
	Err     error
}

// Error implements the error interface
func (e *AnvilError) Error() string {
	if e.Command != "" {
		return fmt.Sprintf("anvil %s %s: %v", e.Op, e.Command, e.Err)
	}
	return fmt.Sprintf("anvil %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error
func (e *AnvilError) Unwrap() error {
	return e.Err
}

// NewAnvilError creates a new AnvilError
func NewAnvilError(op, command string, err error) *AnvilError {
	return &AnvilError{
		Op:      op,
		Command: command,
		Err:     err,
	}
}

// Command operation constants
const (
	OpInit  = "init"
	OpSetup = "setup"
	OpPull  = "pull"
	OpPush  = "push"
	OpDraw  = "draw"
)

// System command constants
const (
	BrewCommand = "brew"
	GitCommand  = "git"
	CurlCommand = "curl"
	ShCommand   = "sh"
)

// Brew subcommand constants
const (
	BrewInstall = "install"
	BrewList    = "list"
	BrewInfo    = "info"
	BrewUpdate  = "update"
	BrewUpgrade = "upgrade"
	BrewSearch  = "search"
)

// Git subcommand constants
const (
	GitConfig    = "config"
	GitGlobal    = "--global"
	GitUserName  = "user.name"
	GitUserEmail = "user.email"
)

// Directory and file constants
const (
	AnvilConfigDir = ".anvil"
	SSHDir         = ".ssh"
	OhMyZshDir     = ".oh-my-zsh"
	ConfigFileName = "settings.yaml"
	CacheSubDir    = "cache"
	DataSubDir     = "data"
)

// macOS package names (Homebrew formulae and casks)
const (
	PkgGit       = "git"
	PkgZsh       = "zsh"
	PkgIterm2    = "iterm2"
	PkgVSCode    = "visual-studio-code"
	PkgSlack     = "slack"
	PkgChrome    = "google-chrome"
	Pkg1Password = "1password"
	PkgDocker    = "docker"
	PkgKubectl   = "kubectl"
	PkgNode      = "node"
	PkgPython    = "python3"
	PkgGo        = "go"
	PkgMysql     = "mysql"
	PkgPostgres  = "postgresql"
	PkgRedis     = "redis"
	PkgVim       = "vim"
	PkgTmux      = "tmux"
	PkgFigma     = "figma"
	PkgNotion    = "notion"
	PkgObsidian  = "obsidian"
)

// Environment variables
const (
	EnvEditor = "EDITOR"
	EnvShell  = "SHELL"
	EnvTerm   = "TERM"
	EnvHome   = "HOME"
	EnvPath   = "PATH"
)

// Oh-my-zsh installation
const (
	OhMyZshInstallURL = "https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh"
	OhMyZshInstallCmd = `sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended`
)

// Common directory permissions
const (
	DirPerm  = 0755
	FilePerm = 0644
)

// Long descriptions for commands
const ANVIL_LONG_DESCRIPTION = `Anvil is a powerful macOS automation CLI tool designed to streamline development workflows 
and personal tool configuration. It provides a comprehensive suite of commands for managing
development environments, automating installations, and maintaining consistent configurations.

Key features:
- Automated tool installation via Homebrew
- Dynamic group and individual app installation
- Environment configuration management
- ASCII art generation for enhanced terminal output
- Optimized specifically for macOS`

const INIT_COMMAND_LONG_DESCRIPTION = `The init command bootstraps your Anvil CLI environment by performing a complete
initialization process. This is the first command you should run after installing Anvil.

What it does:
• Validates and installs required system tools (Git, cURL, Homebrew)
• Creates necessary configuration directories (~/.anvil, ~/.anvil/cache, ~/.anvil/data)
• Generates a default settings.yaml configuration file with your system preferences
• Checks your local environment for common development configurations
• Provides actionable recommendations for completing your setup

This command is designed specifically for macOS and requires Homebrew for tool management.`

const SETUP_COMMAND_LONG_DESCRIPTION = `The setup command provides dynamic installation of development tools and applications 
using Homebrew on macOS.

Usage:
• anvil setup [group-name]    - Install all tools in a predefined group
• anvil setup [app-name]      - Install any individual application via brew

This command intelligently determines if the argument is a group name (defined in settings.yaml)
or an application name. If it's not a group, it attempts to install the application directly
using Homebrew. All installations are validated and errors are handled gracefully.

Supported groups: dev, new-laptop, and custom groups defined in your configuration.

Use 'anvil setup --list' to see all available groups.`

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
