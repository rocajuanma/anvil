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

// Package interfaces provides common abstractions for the Anvil CLI.
//
// This package defines interfaces that promote code reusability and
// testability across different components of the application.
package interfaces

import (
	"context"

	"github.com/rocajuanma/anvil/pkg/system"
)

// CommandExecutor defines the interface for executing system commands
type CommandExecutor interface {
	Execute(ctx context.Context, command string, args ...string) (*system.CommandResult, error)
	ExecuteWithOutput(ctx context.Context, command string, args ...string) error
	ExecuteInDirectory(ctx context.Context, dir, command string, args ...string) (*system.CommandResult, error)
	CommandExists(command string) bool
}

// PackageManager defines the interface for package management operations
type PackageManager interface {
	IsInstalled() bool
	Install() error
	Update() error
	InstallPackage(packageName string) error
	IsPackageInstalled(packageName string) bool
	GetInstalledPackages() ([]Package, error)
	GetPackageInfo(packageName string) (*Package, error)
}

// Package represents a software package
type Package interface {
	GetName() string
	GetVersion() string
	GetDescription() string
	IsInstalled() bool
}

// OutputHandler defines the interface for terminal output operations
type OutputHandler interface {
	PrintHeader(message string)
	PrintStage(message string)
	PrintSuccess(message string)
	PrintError(format string, args ...interface{})
	PrintWarning(format string, args ...interface{})
	PrintInfo(format string, args ...interface{})
	PrintProgress(current, total int, message string)
	Confirm(message string) bool
	IsSupported() bool
}

// ConfigManager defines the interface for configuration management
type ConfigManager interface {
	Load() error
	Save() error
	GetGroupTools(groupName string) ([]string, error)
	GetAvailableGroups() (map[string][]string, error)
	AddInstalledApp(appName string) error
	GetInstalledApps() ([]string, error)
	IsAppTracked(appName string) (bool, error)
	RemoveInstalledApp(appName string) error
}

// ToolInstaller defines the interface for tool installation operations
type ToolInstaller interface {
	Install(toolName string) error
	IsInstalled(toolName string) bool
	GetInfo(toolName string) (*ToolInfo, error)
	ValidateAndInstall() error
}

// ToolInfo represents information about a tool
type ToolInfo struct {
	Name        string
	Command     string
	Required    bool
	InstallWith string
	Description string
}

// CommandRunner defines the interface for command execution with validation
type CommandRunner interface {
	ValidateArgs(args []string) error
	Execute(args []string) error
	GetHelp() string
}

// ProgressReporter defines the interface for progress reporting
type ProgressReporter interface {
	Start(total int, message string)
	Update(current int, message string)
	Finish(message string)
	Error(err error)
}

// Validator defines the interface for input validation
type Validator interface {
	ValidateGroupName(groupName string) error
	ValidateAppName(appName string) error
	ValidateFont(font string) error
	ValidateConfig(config interface{}) error
}

// ErrorHandler defines the interface for error handling
type ErrorHandler interface {
	Handle(err error) error
	Wrap(err error, op, command string) error
	IsRecoverable(err error) bool
	GetUserMessage(err error) string
}

// Logger defines the interface for structured logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// CacheManager defines the interface for cache operations
type CacheManager interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Clear()
	Size() int
}

// FileSystemManager defines the interface for file system operations
type FileSystemManager interface {
	CreateDirectory(path string) error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	FileExists(path string) bool
	DirectoryExists(path string) bool
	DeleteFile(path string) error
	DeleteDirectory(path string) error
}

// PlatformDetector defines the interface for platform detection
type PlatformDetector interface {
	GetOS() string
	GetArch() string
	IsSupported() bool
	GetPackageManager() PackageManager
}

// ServiceManager defines the interface for service management
type ServiceManager interface {
	Start(service string) error
	Stop(service string) error
	Restart(service string) error
	Status(service string) (string, error)
	IsRunning(service string) bool
}
