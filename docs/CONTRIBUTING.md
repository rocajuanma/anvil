# Contributing to Anvil CLI

Thank you for your interest in contributing to Anvil! This guide will help you get started with contributing to the project, whether you're fixing bugs, adding features, or improving documentation.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Development Workflow](#development-workflow)
- [Code Style and Standards](#code-style-and-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Submitting Changes](#submitting-changes)
- [Community](#community)

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. Please read and follow our Code of Conduct:

- **Be respectful** and considerate in all interactions
- **Be constructive** when providing feedback
- **Be patient** with newcomers and help them learn
- **Focus on the issue** rather than personal attacks
- **Respect different viewpoints** and experiences

## Getting Started

### Ways to Contribute

- 🐛 **Bug Reports** - Help us identify and fix issues
- ✨ **Feature Requests** - Suggest new functionality
- 📝 **Documentation** - Improve guides, examples, and API docs
- 🛠️ **Code Contributions** - Fix bugs, implement features
- 🎨 **UI/UX Improvements** - Enhance terminal output and user experience
- 🧪 **Testing** - Add test cases and improve coverage
- 📦 **Package Management** - Support for new platforms or package managers

### Before You Start

1. **Check existing issues** to avoid duplicate work
2. **Join discussions** for clarification on requirements
3. **Start small** with documentation or minor bug fixes
4. **Ask questions** if anything is unclear

## Development Setup

### Prerequisites

- **Go 1.17+** - [Download](https://golang.org/dl/)
- **Git** - Version control
- **Make** - Build automation (optional)
- **Docker** - For testing (optional)

### Clone and Setup

```bash
# Fork the repository on GitHub first
git clone https://github.com/yourusername/anvil.git
cd anvil

# Add upstream remote
git remote add upstream https://github.com/rocajuanma/anvil.git

# Install dependencies
go mod download
go mod tidy

# Build the project
go build -o anvil main.go

# Initialize for development
./anvil init
```

### Development Tools

Install recommended development tools:

```bash
# Code formatting
go install golang.org/x/tools/cmd/goimports@latest

# Linting
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Testing
go install gotest.tools/gotestsum@latest
```

### Project Structure

```
anvil/
├── cmd/                    # Command implementations
│   ├── initcmd/           # Init command
│   ├── setup/             # Setup command
│   ├── config/            # Config command (parent for pull/push)
│   │   ├── pull/          # Pull subcommand (config pull)
│   │   └── push/          # Push subcommand (config push)
│   └── root.go            # Root command configuration
├── internal/                   # Reusable packages
│   ├── brew/              # Homebrew integration
│   ├── config/            # Configuration management
│   ├── constants/         # Application constants
│   ├── errors/            # Error types and structured error handling
│   ├── figure/            # ASCII art generation
│   ├── system/            # System command execution
│   ├── terminal/          # Terminal output formatting
│   └── tools/             # Tool validation and installation
├── docs/                  # Documentation
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── main.go                # Application entry point
└── README.md              # Project overview
```

## Contributing Guidelines

### Issue Guidelines

#### Bug Reports

When reporting bugs, please include:

```
**Bug Description**
A clear description of what went wrong.

**Steps to Reproduce**
1. Step one
2. Step two
3. Step three

**Expected Behavior**
What you expected to happen.

**Actual Behavior**
What actually happened.

**Environment**
- OS: [e.g., macOS 12.0, Ubuntu 20.04]
- Go version: [e.g., 1.17.3]
- Anvil version: [e.g., 1.0.0]

**Additional Context**
- Error messages (full output)
- Configuration files (if relevant)
- Screenshots (if applicable)
```

#### Feature Requests

For feature requests, please include:

```
**Feature Description**
Clear description of the proposed feature.

**Use Case**
Why is this feature needed? What problem does it solve?

**Proposed Solution**
How should this feature work?

**Alternatives Considered**
Other approaches you've considered.

**Additional Context**
Mockups, examples, or related issues.
```

### Pull Request Guidelines

#### Before Submitting

- [ ] **Search existing PRs** to avoid duplicates
- [ ] **Create an issue** for discussion if it's a significant change
- [ ] **Fork the repository** and create a feature branch
- [ ] **Write tests** for new functionality
- [ ] **Update documentation** as needed
- [ ] **Follow code style** guidelines
- [ ] **Test thoroughly** on your platform

#### PR Description Template

```
**Description**
Brief description of changes.

**Type of Change**
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Refactoring

**Related Issues**
Fixes #123, Related to #456

**Testing**
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] Tested on multiple platforms

**Documentation**
- [ ] Code comments updated
- [ ] README updated
- [ ] API documentation updated
- [ ] Examples added/updated

**Checklist**
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Tests pass locally
- [ ] No breaking changes (or documented)
```

## Development Workflow

### 1. Create Feature Branch

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number-short-description
```

### 2. Make Changes

Follow our development standards for:

- **Package organization** - Use appropriate package structure with clear separation of concerns
- **Error handling patterns** - Use `errors.AnvilError` for structured error handling with operation context and error types
- **Terminal output formatting** - Use `terminal` package for consistent output formatting
- **Configuration management** - Use cached configuration access for optimal performance
- **Constants usage** - Use constants from `internal/constants/` instead of magic strings
- **Documentation standards** - Update relevant documentation for any changes

### 3. Test Changes

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Test specific package
go test ./internal/config/

# Manual testing
go build -o anvil-dev main.go
./anvil-dev init
./anvil-dev setup --list
```

### 4. Commit Changes

Use conventional commit messages:

```bash
# Examples
git commit -m "feat: add custom group support to setup command"
git commit -m "fix: resolve homebrew path issues on M1 Macs"
git commit -m "docs: update installation guide for Linux"
git commit -m "test: add unit tests for config package"
git commit -m "refactor: improve error handling in brew package"
```

### 5. Push and Create PR

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create pull request on GitHub
# Fill out the PR template completely
```

## Code Style and Standards

### Go Code Style

Follow standard Go conventions and our specific guidelines:

#### Import Organization

```go
import (
    // Standard library
    "fmt"
    "os"
    "path/filepath"

    // Third-party packages
    "github.com/spf13/cobra"
    "gopkg.in/yaml.v2"

    // Local packages
    "github.com/rocajuanma/anvil/internal/config"
    "github.com/rocajuanma/anvil/internal/terminal"
)
```

#### Error Handling

```go
// Always provide context
if err := someOperation(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Include relevant information
if err := installTool(toolName); err != nil {
    return fmt.Errorf("failed to install %s: %w", toolName, err)
}
```

#### Function Documentation

```go
// InstallTool installs a specific tool using the appropriate package manager
// This function determines the best installation method based on the platform
// and tool availability, providing detailed progress feedback to the user.
//
// Parameters:
//   toolName: the name of the tool to install
//
// Returns:
//   error: nil on success, or an error describing what went wrong
func InstallTool(toolName string) error {
    // Implementation
}
```

### Terminal Output Standards

Use consistent structure output patterns from https://github.com/rocajuanma/palantir

```go
// Use terminal package functions
output = palantir.GetGlobalOutputHandler()
output.PrintHeader("Major Section")
output.PrintStage("Processing step...")
output.PrintSuccess("Operation completed")
output.PrintInfo("General information")
output.PrintWarning("Non-critical issue")
output.PrintError("Error occurred: %v", err)

// Progress indicators for multi-step operations
for i, item := range items {
    output.PrintProgress(i+1, len(items), fmt.Sprintf("Processing %s", item))
    // Process item
}
```

### Configuration Standards

Follow YAML configuration patterns:

```yaml
# Use consistent structure
version: "1.0.0"
directories:
  config: "/path/to/config"
  cache: "/path/to/cache"
  data: "/path/to/data"
tools:
  required_tools: []
  optional_tools: []
groups:
  predefined_group: []
environment: {}
```

## Testing

### Unit Tests

Write unit tests for new functionality:

```go
package config

import (
    "testing"
    "os"
    "path/filepath"
)

func TestGetDefaultConfig(t *testing.T) {
    config := GetDefaultConfig()

    if config.Version == "" {
        t.Error("Version should not be empty")
    }

    if len(config.Tools.RequiredTools) == 0 {
        t.Error("Should have required tools")
    }
}

func TestCreateDirectories(t *testing.T) {
    // Setup test directory
    tmpDir := t.TempDir()
    os.Setenv("HOME", tmpDir)
    defer os.Unsetenv("HOME")

    // Test directory creation
    err := CreateDirectories()
    if err != nil {
        t.Fatalf("CreateDirectories failed: %v", err)
    }

    // Verify directories exist
    configDir := filepath.Join(tmpDir, ".anvil")
    if _, err := os.Stat(configDir); os.IsNotExist(err) {
        t.Error("Config directory was not created")
    }
}
```

### Integration Tests

Test command functionality:

```go
func TestInitCommand(t *testing.T) {
    // Setup isolated test environment
    tmpDir := t.TempDir()
    os.Setenv("HOME", tmpDir)
    defer os.Unsetenv("HOME")

    // Run init command
    cmd := exec.Command("./anvil", "init")
    output, err := cmd.CombinedOutput()

    if err != nil {
        t.Fatalf("Init command failed: %v\nOutput: %s", err, output)
    }

    // Verify expected results
    configFile := filepath.Join(tmpDir, ".anvil", "settings.yaml")
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        t.Error("Settings file was not created")
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestGetDefaultConfig ./internal/config/

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

### Code Documentation

- **Document all public functions, types, and constants**
- **Use standard Go documentation format**
- **Include usage examples for complex functions**
- **Explain non-obvious behavior**

### User Documentation

When adding features that affect users:

- **Update command help text** in constants
- **Add examples** to relevant documentation files
- **Update README** if it affects basic usage
- **Create or update** specific command documentation

### Documentation Files to Update

- `README.md` - For major features or changes
- `docs/GETTING_STARTED.md` - For user-facing changes
- `docs/command-readme.md` - For command-specific features
- `docs/EXAMPLES.md` - Add usage examples
- Help text in `internal/constants/constants.go`

## Submitting Changes

### Pre-submission Checklist

- [ ] **Code compiles** without errors or warnings
- [ ] **Tests pass** on your local machine
- [ ] **Code follows** style guidelines
- [ ] **Documentation updated** as needed
- [ ] **Commit messages** follow conventional format
- [ ] **PR description** is complete and clear
- [ ] **No sensitive information** in commits

### Review Process

1. **Automated checks** run on PR creation
2. **Code review** by maintainers
3. **Discussion and feedback** as needed
4. **Approval and merge** by maintainers

### After Submission

- **Respond to feedback** promptly
- **Make requested changes** in additional commits
- **Ask questions** if feedback is unclear
- **Be patient** - review process may take time

## Community

### Communication Channels

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - General questions and community chat
- **Pull Requests** - Code review and discussion

### Getting Help

- **Check documentation** first
- **Search existing issues** for similar problems
- **Ask specific questions** with context
- **Provide relevant information** (OS, Go version, etc.)

### Mentoring

We welcome new contributors! If you're new to:

- **Go programming** - We can help with Go-specific questions
- **CLI development** - We can guide you through CLI patterns
- **Open source** - We can help with the contribution process

### Recognition

Contributors are recognized in:

- **GitHub contributors** page
- **Release notes** for significant contributions
- **Special thanks** in documentation

---

## Development Resources

### Useful Links

- **[Go Documentation](https://golang.org/doc/)**
- **[Cobra CLI Framework](https://github.com/spf13/cobra)**
- **[Effective Go](https://golang.org/doc/effective_go.html)**
- **[Anvil Development Rules](.local/anvil-rules.md)**

### Tools and Extensions

- **[VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.go)**
- **[GoLand IDE](https://www.jetbrains.com/go/)**
- **[Go Tools](https://golang.org/doc/install#extra-tools)**

---

## Questions?

Don't hesitate to ask questions! We're here to help:

- **Open an issue** for bugs or feature requests
- **Start a discussion** for general questions
- **Comment on PRs** for code-related questions

Thank you for contributing to Anvil! 🎉
