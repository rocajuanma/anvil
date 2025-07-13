# Anvil CLI Code Improvements Report

## Executive Summary

This report analyzes the Anvil CLI codebase for idiomatic Go usage, redundancy, and opportunities for improvement. The analysis covers 12 Go files across commands and packages, identifying improvement opportunities organized into 4 priority levels.

**Key Findings:**

- ✅ **Priority 1 (COMPLETED)**: 3 critical issues have been fixed
  - Error handling improvements in setup and init commands
  - Resource management with context and timeout support
  - Elimination of global variables in favor of struct-based approach
- 9 Medium Priority items for code quality enhancement
- 6 Low Priority items for long-term maintainability
- 5 Enhancement items for future development
- 3 Documentation items for developer experience

## Priority Classification

### 🟡 **Priority 2: High Impact Improvements**

#### 2.1 Reduce Code Duplication in Installation Functions

**File:** `cmd/setup/setup.go`

**Current Issue:**

```go
func installGit() error {
    if system.CommandExists("git") {
        return nil // Already installed
    }
    return brew.InstallPackage("git")
}

func installZsh() error {
    if err := brew.InstallPackage("zsh"); err != nil {
        return fmt.Errorf("failed to install zsh: %w", err)
    }
    // ... additional logic
}
```

**Improvement:**

```go
type InstallConfig struct {
    PackageName    string
    PreCheck       func() bool
    PostInstall    func() error
    SkipIfExists   bool
}

func installWithConfig(config InstallConfig) error {
    if config.SkipIfExists && config.PreCheck() {
        return nil
    }

    if err := brew.InstallPackage(config.PackageName); err != nil {
        return fmt.Errorf("failed to install %s: %w", config.PackageName, err)
    }

    if config.PostInstall != nil {
        return config.PostInstall()
    }
    return nil
}
```

**Why Important:** Reduces maintenance burden and ensures consistent behavior across installations.

#### 2.2 Configuration Loading Optimization

**File:** `pkg/config/config.go`

**Current Issue:**

```go
func GetGroupTools(groupName string) ([]string, error) {
    config, err := LoadConfig() // Loads entire config file
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }
    // ... rest of function
}
```

**Improvement:**

```go
var configCache *AnvilConfig
var configCacheMutex sync.RWMutex

func GetGroupTools(groupName string) ([]string, error) {
    config, err := getCachedConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }
    // ... rest of function
}

func getCachedConfig() (*AnvilConfig, error) {
    configCacheMutex.RLock()
    if configCache != nil {
        configCacheMutex.RUnlock()
        return configCache, nil
    }
    configCacheMutex.RUnlock()

    configCacheMutex.Lock()
    defer configCacheMutex.Unlock()

    if configCache != nil {
        return configCache, nil
    }

    var err error
    configCache, err = LoadConfig()
    return configCache, err
}
```

**Why Important:** Prevents repeated file I/O operations and improves performance.

#### 2.3 Consistent Error Types

**Files:** All command files

**Current Issue:**

```go
// Inconsistent error handling across files
os.Exit(1) // In some files
return err // In others
fmt.Errorf("...") // Various formats
```

**Improvement:**

```go
// Create custom error types
type AnvilError struct {
    Op      string
    Command string
    Err     error
}

func (e *AnvilError) Error() string {
    return fmt.Sprintf("anvil %s %s: %v", e.Op, e.Command, e.Err)
}

func (e *AnvilError) Unwrap() error {
    return e.Err
}
```

**Why Important:** Provides consistent error handling and better error messages for users.

### 🟢 **Priority 3: Medium Impact Improvements**

#### 3.1 Extract Magic Strings to Constants

**Files:** `pkg/brew/brew.go`, `cmd/setup/setup.go`, `pkg/config/config.go`

**Current Issue:**

```go
// Magic strings scattered throughout code
result, err := system.RunCommand("brew", "install", packageName)
sshDir := filepath.Join(homeDir, ".ssh")
return filepath.Join(homeDir, ".anvil")
```

**Improvement:**

```go
// In pkg/constants/constants.go
const (
    BrewCommand        = "brew"
    AnvilConfigDir     = ".anvil"
    SSHDir            = ".ssh"
    ConfigFileName    = "settings.yaml"
    CacheSubDir       = "cache"
    DataSubDir        = "data"
)
```

**Why Important:** Improves maintainability and reduces typos.

#### 3.2 Improve Terminal Output Consistency

**File:** `pkg/terminal/terminal.go`

**Current Issue:**

```go
// Inconsistent emoji usage and formatting
func PrintStage(message string) {
    fmt.Printf("%s%s🔧 %s%s\n", ColorBold, ColorBlue, message, ColorReset)
}
```

**Improvement:**

```go
type OutputLevel int

const (
    LevelInfo OutputLevel = iota
    LevelWarning
    LevelError
    LevelSuccess
    LevelStage
)

func PrintWithLevel(level OutputLevel, format string, args ...interface{}) {
    if !IsTerminalSupported() {
        fmt.Printf(format+"\n", args...)
        return
    }

    var prefix, color string
    switch level {
    case LevelStage:
        prefix, color = "🔧 ", ColorBlue
    case LevelSuccess:
        prefix, color = "✅ ", ColorGreen
    // ... other cases
    }

    fmt.Printf("%s%s%s%s%s\n", ColorBold, color, prefix, fmt.Sprintf(format, args...), ColorReset)
}
```

**Why Important:** Provides consistent user experience and better terminal compatibility.

#### 3.3 Add Validation Functions

**Files:** `pkg/config/config.go`, `cmd/setup/setup.go`

**Current Issue:**

```go
// No validation of user input or configuration
func GetGroupTools(groupName string) ([]string, error) {
    // Direct access without validation
    switch groupName {
    case "dev":
        return config.Groups.Dev, nil
    }
}
```

**Improvement:**

```go
func ValidateGroupName(groupName string) error {
    if groupName == "" {
        return fmt.Errorf("group name cannot be empty")
    }

    validGroups := []string{"dev", "new-laptop"}
    for _, valid := range validGroups {
        if groupName == valid {
            return nil
        }
    }

    return fmt.Errorf("invalid group name: %s", groupName)
}
```

**Why Important:** Prevents runtime errors and provides better user feedback.

### 🔵 **Priority 4: Low Impact Improvements**

#### 4.1 Improve Package Documentation

**Files:** All package files

**Current Issue:**

```go
// Minimal package documentation
package terminal

import (
    "fmt"
    "os"
)
```

**Improvement:**

```go
// Package terminal provides colored terminal output utilities for the Anvil CLI.
//
// This package offers functions for consistent terminal output formatting
// including headers, status messages, progress indicators, and user prompts.
// It automatically detects terminal capabilities and gracefully degrades
// for unsupported terminals.
package terminal
```

**Why Important:** Improves code maintainability and developer experience.

#### 4.2 Add Input Validation

**Files:** `cmd/draw/draw.go`, `cmd/pull/pull.go`, `cmd/push/push.go`

**Current Issue:**

```go
// No validation of command arguments
Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("draw called")
    figure.Draw("anvil", args[0]) // Potential panic if args is empty
},
```

**Improvement:**

```go
Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
        terminal.PrintError("Font argument is required")
        return
    }

    font := args[0]
    if !isValidFont(font) {
        terminal.PrintError("Invalid font: %s", font)
        return
    }

    figure.Draw("anvil", font)
},
```

**Why Important:** Prevents runtime panics and provides better user experience.

#### 4.3 Implement Proper Logging

**Files:** All command files

**Current Issue:**

```go
// No structured logging
terminal.PrintError("Failed to install %s: %v", tool, err)
```

**Improvement:**

```go
// Add structured logging
import "log/slog"

func setupLogger() *slog.Logger {
    opts := &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }

    handler := slog.NewTextHandler(os.Stderr, opts)
    return slog.New(handler)
}
```

**Why Important:** Enables better debugging and monitoring in production.

### 🟣 **Priority 5: Enhancement Opportunities**

#### 5.1 Add Configuration Validation

**File:** `pkg/config/config.go`

**Improvement:**

```go
func (c *AnvilConfig) Validate() error {
    if c.Version == "" {
        return fmt.Errorf("version cannot be empty")
    }

    if len(c.Groups.Dev) == 0 {
        return fmt.Errorf("dev group cannot be empty")
    }

    return nil
}
```

**Why Important:** Prevents invalid configurations from causing runtime errors.

#### 5.2 Add Progress Reporting

**Files:** `cmd/setup/setup.go`, `cmd/initcmd/init.go`

**Improvement:**

```go
type ProgressReporter struct {
    total   int
    current int
    ch      chan string
}

func (pr *ProgressReporter) Report(message string) {
    pr.current++
    pr.ch <- fmt.Sprintf("[%d/%d] %s", pr.current, pr.total, message)
}
```

**Why Important:** Provides better user feedback for long-running operations.

#### 5.3 Add Concurrent Installation

**File:** `cmd/setup/setup.go`

**Improvement:**

```go
func installToolsConcurrently(tools []string) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(tools))

    for _, tool := range tools {
        wg.Add(1)
        go func(t string) {
            defer wg.Done()
            if err := installTool(t); err != nil {
                errChan <- fmt.Errorf("failed to install %s: %w", t, err)
            }
        }(tool)
    }

    wg.Wait()
    close(errChan)

    var errors []string
    for err := range errChan {
        errors = append(errors, err.Error())
    }

    if len(errors) > 0 {
        return fmt.Errorf("installation errors: %s", strings.Join(errors, "; "))
    }

    return nil
}
```

**Why Important:** Significantly improves installation performance.

## Implementation Recommendations

### ✅ **Phase 1: Critical Fixes (COMPLETED)**

1. ✅ Fix error handling in all command functions
2. ✅ Add context and timeout support to system commands
3. ✅ Eliminate global variables in setup command
4. ✅ Implement proper resource management

### Phase 2: Code Quality (2-3 weeks)

1. Reduce code duplication in installation functions
2. Add configuration caching
3. Implement consistent error types
4. Extract magic strings to constants

### Phase 3: Enhancements (3-4 weeks)

1. Add input validation across all commands
2. Implement structured logging
3. Add configuration validation
4. Improve terminal output consistency

### Phase 4: Performance & Features (4-6 weeks)

1. Add progress reporting
2. Implement concurrent installation
3. Add comprehensive package documentation
4. Create unit tests for all packages

## Testing Strategy

### Unit Testing Requirements

```go
// Example test structure
func TestInstallTool(t *testing.T) {
    tests := []struct {
        name     string
        toolName string
        wantErr  bool
    }{
        {"valid tool", "git", false},
        {"invalid tool", "nonexistent", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := installTool(tt.toolName)
            if (err != nil) != tt.wantErr {
                t.Errorf("installTool() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Testing

- Test complete init workflow
- Test group installation scenarios
- Test error recovery mechanisms
- Test cross-platform compatibility

## Metrics for Success

### Code Quality Metrics

- **Cyclomatic Complexity**: Target < 10 per function
- **Test Coverage**: Target > 80%
- **Code Duplication**: Target < 5%
- **Go Report Card**: Target A+ grade

### Performance Metrics

- **Installation Time**: 50% reduction with concurrent installation
- **Memory Usage**: < 50MB peak usage
- **Startup Time**: < 100ms for help commands

### User Experience Metrics

- **Error Recovery**: 95% of errors should provide actionable feedback
- **Progress Feedback**: All operations > 2 seconds should show progress
- **Cross-platform Support**: 100% feature parity on supported platforms

## Long-term Architecture Considerations

### Plugin System

Consider implementing a plugin architecture for:

- Custom tool installers
- Additional package managers
- Platform-specific optimizations

### Configuration Management

- Support for multiple configuration formats (YAML, JSON, TOML)
- Environment-specific configurations
- Configuration validation and migration

### CLI Framework Migration

Consider migrating from Cobra to a more modern CLI framework if:

- Advanced features are needed
- Performance becomes a concern
- Maintenance burden increases

## Conclusion

The Anvil CLI codebase has had its critical issues resolved and shows good foundational structure. The remaining prioritized improvements outlined in this report will further enhance code maintainability, user experience, and overall software quality.

**Recently Completed (Priority 1):**

1. ✅ Fixed error handling patterns in setup and init commands
2. ✅ Eliminated global variables in favor of struct-based approach
3. ✅ Added proper resource management with context and timeout support

**Next Steps:**

1. Create issues for Priority 2 items
2. Assign development resources
3. Establish testing framework
4. Set up continuous integration

The estimated effort for implementing all remaining improvements is 9-12 weeks for a single developer, or 5-7 weeks for a team of 2-3 developers working in parallel.
