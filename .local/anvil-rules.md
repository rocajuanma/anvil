# Anvil Development Rules & Guidelines

## Overview

This document outlines the development methodology, patterns, and guidelines used for implementing Anvil CLI commands. It serves as a reference for maintaining consistency, quality, and best practices across all command implementations.

## Development Philosophy

### Core Principles

1. **Idiomatic Go** - Follow Go best practices and conventions
2. **User-Centric Design** - Prioritize user experience and clear feedback
3. **Modular Architecture** - Organize code into logical, reusable packages
4. **Comprehensive Documentation** - Provide thorough documentation for all features
5. **Error Resilience** - Handle errors gracefully with actionable guidance
6. **Cross-Platform Awareness** - Consider platform differences and limitations

## Project Structure Rules

### Directory Organization

```
anvil/
├── cmd/                    # Command implementations
│   ├── commandname/        # Individual command packages
│   │   └── command.go      # Command definition and logic
│   └── root.go            # Root command configuration
├── pkg/                   # Reusable packages
│   ├── packagename/       # Logical groupings by functionality
│   │   └── package.go     # Package implementation
├── docs/                  # Command documentation
│   └── command-readme.md  # Comprehensive command docs
└── .local/               # Internal development files
    └── rules.md          # Development guidelines (this file)
```

### Package Organization Rules

1. **Single Responsibility** - Each package should have one clear purpose
2. **Clear Naming** - Package names should be descriptive and concise
3. **Logical Grouping** - Group related functionality together
4. **Minimal Dependencies** - Avoid circular dependencies and excessive coupling

### Required Packages for Commands

Every command implementation should utilize these core packages:

- `pkg/constants` - Command descriptions and constants
- `pkg/terminal` - User interface and output formatting
- `pkg/config` - Configuration management
- Platform-specific packages as needed (e.g., `pkg/brew`, `pkg/system`)

## Command Development Process

### Phase 1: Analysis & Requirements

#### 1.1 Requirements Gathering

- [ ] **Analyze user requirements** from specifications, images, or descriptions
- [ ] **Identify core functionality** and user goals
- [ ] **Determine input/output requirements** (arguments, flags, expected behavior)
- [ ] **Assess integration points** with existing commands and packages

#### 1.2 Existing Code Analysis

- [ ] **Read current project structure** to understand patterns
- [ ] **Examine existing commands** for consistency patterns
- [ ] **Identify reusable components** and packages
- [ ] **Note dependency patterns** and import structures

#### 1.3 Technical Design

- [ ] **Define command interface** (arguments, flags, subcommands)
- [ ] **Plan package dependencies** and new package requirements
- [ ] **Design error handling strategy** and user feedback approach
- [ ] **Consider platform compatibility** and limitations

### Phase 2: Package Development

#### 2.1 Core Package Implementation

Follow this order for package development:

1. **System/Utility Packages** (`pkg/system`, `pkg/terminal`)

   - Implement foundational functionality first
   - Ensure proper error handling and cross-platform support
   - Add comprehensive helper functions

2. **Platform-Specific Packages** (`pkg/brew`, `pkg/tools`)

   - Implement platform-specific functionality
   - Include proper platform detection and warnings
   - Handle missing dependencies gracefully

3. **Configuration Packages** (`pkg/config`)
   - Implement configuration management
   - Support both reading and writing configurations
   - Include validation and default value handling

#### 2.2 Package Development Rules

**Structure Requirements:**

```go
// Package declaration with clear purpose
package packagename

// Standard library imports first
import (
    "fmt"
    "os"
)

// Third-party imports
import (
    "github.com/spf13/cobra"
)

// Local imports last
import (
    "github.com/rocajuanma/anvil/pkg/otherpkg"
)

// Public types and constants
type PublicType struct {
    Field string
}

// Public functions with comprehensive documentation
// FunctionName performs a specific action and returns result or error
func FunctionName(param string) error {
    // Implementation with proper error handling
    return nil
}
```

**Error Handling Standards:**

- Always return meaningful errors with context
- Use `fmt.Errorf` with verb `%w` for error wrapping
- Provide actionable error messages when possible
- Log intermediate steps for debugging

**Documentation Standards:**

- Document all public functions, types, and constants
- Use standard Go documentation format
- Include usage examples for complex functions
- Explain non-obvious behavior or limitations

### Phase 3: Command Implementation

#### 3.1 Command Structure Template

```go
/*
Copyright © 2022 Juanma Roca juanmaxroca@gmail.com
[Standard license header]
*/

package commandname

import (
    // Standard imports
    // Third-party imports
    // Local imports
)

// CommandCmd represents the command with clear documentation
// This command performs [specific function] and [expected outcome]
var CommandCmd = &cobra.Command{
    Use:   "command [args]",
    Short: "Brief description",
    Long:  constants.COMMAND_LONG_DESCRIPTION,
    Args:  cobra.ExactArgs(1), // or appropriate validation
    Run: func(cmd *cobra.Command, args []string) {
        runCommandFunction(cmd, args)
    },
}

// Command flags
var (
    flagName bool
    flagValue string
)

// runCommandFunction executes the main command logic
// This function orchestrates all command operations and handles errors
func runCommandFunction(cmd *cobra.Command, args []string) {
    // Stage-based implementation with clear progression
    terminal.PrintHeader("Command Execution")

    // Stage 1: Validation
    terminal.PrintStage("Validating inputs...")
    if err := validateInputs(args); err != nil {
        terminal.PrintError("Validation failed: %v", err)
        terminal.PrintError("Use 'anvil command --help' for usage information")
        os.Exit(1)
    }
    terminal.PrintSuccess("Input validation complete")

    // Additional stages as needed...

    // Final stage: Completion
    terminal.PrintHeader("Command Complete!")
    terminal.PrintInfo("Next steps and guidance...")
}

// Helper functions with clear purpose and documentation
func validateInputs(args []string) error {
    // Implementation
    return nil
}

func init() {
    // Flag definitions with clear descriptions
    CommandCmd.Flags().BoolVar(&flagName, "flag", false, "Clear description of flag purpose")
    CommandCmd.Flags().StringVar(&flagValue, "value", "", "Description with expected format")

    // Required flags if applicable
    CommandCmd.MarkFlagRequired("required-flag")
}
```

#### 3.2 Command Implementation Rules

**User Experience Standards:**

- Use staged output with `terminal.PrintStage()` for long operations
- Provide clear progress indicators for multi-step processes
- Include helpful error messages with suggested remediation
- Show completion status and next steps guidance

**Flag Design Rules:**

- Use descriptive long-form flag names
- Provide clear flag descriptions
- Group related flags logically
- Support both boolean and value flags as appropriate

**Error Handling in Commands:**

- Validate inputs early and provide clear error messages
- Continue processing other items if one fails (when appropriate)
- Provide summary information for batch operations
- Exit with appropriate status codes

### Phase 4: Enhancement & Documentation

#### 4.1 Description Enhancement

**Constants Update Process:**

1. **Review existing descriptions** in `pkg/constants/constants.go`
2. **Enhance with comprehensive details** including:
   - Clear purpose statement
   - Feature highlights with bullet points
   - Usage context and benefits
   - Integration information with other commands

**Description Template:**

```go
const COMMAND_LONG_DESCRIPTION = `The command provides [main purpose] by performing [key actions].

What it does:
• [Key feature 1] - [benefit/explanation]
• [Key feature 2] - [benefit/explanation]
• [Key feature 3] - [benefit/explanation]

The command is designed to be [key characteristic] and [additional benefit].
It [integration context] and [usage guidance].

This command is [importance statement] for [target use case].`
```

#### 4.2 Code Documentation Enhancement

**File Header Requirements:**

- Standard copyright and license header
- Clear package purpose documentation
- Import organization (standard, third-party, local)

**Function Documentation Standards:**

```go
// FunctionName performs a specific action with given parameters
// This function [detailed explanation of behavior and purpose]
//
// Parameters:
//   param1: description of first parameter
//   param2: description of second parameter
//
// Returns:
//   error: description of possible error conditions
//
// Example usage:
//   if err := FunctionName("value1", "value2"); err != nil {
//       // handle error
//   }
func FunctionName(param1, param2 string) error {
    // Implementation with inline comments for complex logic
    return nil
}
```

#### 4.3 Comprehensive Documentation Creation

**Documentation Structure Template:**

```markdown
# Command Documentation

## Overview

[Brief description and purpose]

## Purpose and Importance

[Why this command exists and its value]

## Command Usage

### Basic Usage

### Advanced Usage

### Examples

## Features and Functionality

### Feature 1

### Feature 2

## Configuration

[How configuration works]

## Expected Output

[What users should expect to see]

## Error Handling

[Common errors and solutions]

## Best Practices

[Recommended usage patterns]

## Troubleshooting

[Common issues and resolution]

## Security Considerations

[Security aspects and guidelines]

## Future Enhancements

[Planned improvements]

## Conclusion

[Summary and final guidance]
```

**Documentation Content Rules:**

- **Be Comprehensive** - Cover all aspects of functionality
- **Include Examples** - Provide real-world usage examples
- **Add Context** - Explain why features exist and when to use them
- **Provide Troubleshooting** - Address common issues and solutions
- **Consider All Users** - From beginners to advanced users
- **Keep Updated** - Ensure documentation reflects current functionality

## Quality Assurance Process

### Phase 5: Testing & Validation

#### 5.1 Build Testing

```bash
# Verify clean build
go build -o anvil main.go

# Test basic functionality
./anvil --help
./anvil command --help
```

#### 5.2 Functional Testing

- [ ] **Test all command modes** (flags, arguments, subcommands)
- [ ] **Test error conditions** (invalid input, missing dependencies)
- [ ] **Test edge cases** (empty input, very long input, special characters)
- [ ] **Test platform compatibility** (verify warnings on unsupported platforms)

#### 5.3 User Experience Testing

- [ ] **Verify help text clarity** and completeness
- [ ] **Test command flow** from user perspective
- [ ] **Validate error messages** are helpful and actionable
- [ ] **Check output formatting** and visual hierarchy

#### 5.4 Integration Testing

- [ ] **Test with existing commands** for consistency
- [ ] **Verify configuration integration** works properly
- [ ] **Test command chaining** scenarios if applicable
- [ ] **Validate package dependencies** resolve correctly

## File Organization Rules

### Command Files

```
cmd/commandname/
└── command.go          # Single file per command
```

### Package Files

```
pkg/packagename/
└── package.go          # Primary implementation
└── types.go           # Type definitions (if extensive)
└── utils.go           # Utility functions (if needed)
```

### Documentation Files

```
docs/
├── command-readme.md   # One comprehensive doc per command
└── overview.md        # General project documentation
```

### Internal Files

```
.local/
├── anvil-rules.md     # This development guide
├── notes.md           # Development notes (optional)
└── templates/         # Code templates (optional)
```

## Code Style Guidelines

### Naming Conventions

- **Packages**: lowercase, single word when possible
- **Files**: lowercase with hyphens for multiple words
- **Functions**: PascalCase for public, camelCase for private
- **Variables**: camelCase, descriptive names
- **Constants**: UPPER_SNAKE_CASE for exported constants

### Code Organization

```go
// Package-level constants
const (
    DefaultValue = "default"
    MaxRetries   = 3
)

// Package-level variables
var (
    globalConfig *Config
    initialized  bool
)

// Type definitions
type ConfigStruct struct {
    Field1 string
    Field2 int
}

// Public functions
func PublicFunction() error {
    return nil
}

// Private functions
func privateHelper() string {
    return "helper"
}
```

### Error Handling Patterns

```go
// Standard error handling
if err := someOperation(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Error with context
if err := validateInput(input); err != nil {
    return fmt.Errorf("invalid input '%s': %w", input, err)
}

// Multiple return values
result, err := processData(data)
if err != nil {
    return nil, fmt.Errorf("data processing failed: %w", err)
}
```

## Terminal Output Standards

### Output Hierarchy

1. **Headers** - `terminal.PrintHeader()` for major sections
2. **Stages** - `terminal.PrintStage()` for process steps
3. **Success** - `terminal.PrintSuccess()` for completion messages
4. **Info** - `terminal.PrintInfo()` for general information
5. **Warnings** - `terminal.PrintWarning()` for non-critical issues
6. **Errors** - `terminal.PrintError()` for problems requiring attention

### Progress Indicators

```go
// For multi-step processes
for i, item := range items {
    terminal.PrintProgress(i+1, len(items), fmt.Sprintf("Processing %s", item))
    // Process item
}
```

### User Guidance

```go
// Always provide next steps
terminal.PrintInfo("You can now use:")
terminal.PrintInfo("  • 'anvil command --help' for more information")
terminal.PrintInfo("  • 'anvil other-command' for related functionality")
```

## Configuration Management Rules

### Settings Structure

```yaml
# Standard configuration sections
version: "1.0.0"
directories:
  config: "/path/to/config"
  cache: "/path/to/cache"
  data: "/path/to/data"
tools:
  required_tools: []
  optional_tools: []
groups:
  group_name: []
  custom: {}
git:
  username: ""
  email: ""
environment: {}
```

### Configuration Functions

```go
// Standard configuration function patterns
func GetConfig() (*Config, error)           // Load configuration
func SaveConfig(*Config) error              // Save configuration
func GetDefaultConfig() *Config             // Generate defaults
func ValidateConfig(*Config) error          // Validate configuration
```

## Dependency Management

### Import Order

1. Standard library imports
2. Third-party imports
3. Local project imports

### Dependency Selection Criteria

- **Stability** - Prefer well-maintained packages
- **Size** - Avoid heavy dependencies for simple tasks
- **Compatibility** - Ensure cross-platform support
- **License** - Verify license compatibility

### Version Management

- Use go.mod for dependency versioning
- Pin to specific versions for stability
- Regular dependency updates with testing

## Security Considerations

### Input Validation

- Validate all user inputs
- Sanitize file paths and command arguments
- Check for injection attempts

### File Operations

- Use proper file permissions (0644 for files, 0755 for directories)
- Validate file paths are within expected boundaries
- Handle symbolic links appropriately

### External Commands

- Validate commands exist before execution
- Use absolute paths when possible
- Sanitize command arguments

## Documentation Standards

### README Structure

1. **Overview** - What the command does
2. **Purpose** - Why it exists and value proposition
3. **Usage** - How to use it with examples
4. **Features** - Detailed feature descriptions
5. **Configuration** - How configuration works
6. **Output** - What to expect
7. **Troubleshooting** - Common issues and solutions
8. **Best Practices** - Recommended usage patterns
9. **Advanced Usage** - Complex scenarios
10. **Future Enhancements** - Planned improvements

### Documentation Maintenance

- Update documentation with code changes
- Include version information
- Provide migration guides for breaking changes
- Keep examples current and tested

## Release Process

### Pre-Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Breaking changes are documented
- [ ] Version numbers are updated
- [ ] Dependencies are reviewed

### Release Validation

- [ ] Build succeeds on target platforms
- [ ] Basic functionality tests pass
- [ ] Help text is accurate
- [ ] Examples in documentation work

## Continuous Improvement

### Code Review Guidelines

- Check adherence to these rules
- Verify user experience quality
- Validate error handling completeness
- Ensure documentation accuracy

### Feedback Integration

- Monitor user feedback and issues
- Update rules based on lessons learned
- Refine processes for better efficiency
- Share improvements across the team

---

## Summary

These rules provide a comprehensive framework for developing high-quality, consistent Anvil CLI commands. By following these guidelines, developers can ensure:

- **Consistent User Experience** across all commands
- **High Code Quality** with proper organization and documentation
- **Maintainable Codebase** with clear patterns and structures
- **Comprehensive Documentation** that serves users effectively
- **Reliable Functionality** with proper error handling and testing

The key to successful Anvil development is balancing technical excellence with user-centric design, always prioritizing the developer experience while maintaining clean, maintainable code.
