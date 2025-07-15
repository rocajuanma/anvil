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
	OpInit   = "init"
	OpSetup  = "setup"
	OpConfig = "config"
	OpPull   = "pull"
	OpPush   = "push"
	OpDraw   = "draw"
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
