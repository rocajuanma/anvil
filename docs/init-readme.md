# Anvil Init Command Documentation

## Overview

The `anvil init` command is the foundation of the Anvil CLI tool. It performs a comprehensive initialization process that bootstraps your development environment and ensures all necessary components are properly configured for optimal Anvil usage.

## Purpose and Importance

The init command serves as the essential first step after installing Anvil. It:

- **Establishes a solid foundation** for all other Anvil commands
- **Automates tedious setup tasks** that would otherwise require manual configuration
- **Ensures consistency** across different development environments
- **Provides immediate feedback** on system readiness and configuration issues
- **Creates a standardized configuration structure** that other Anvil commands depend on

## What the Init Command Does

### Stage 1: Tool Validation and Installation

The init command validates and installs required system tools:

**Required Tools:**

- **Git** - Essential for version control and asset synchronization
- **cURL** - Required for downloading resources and API interactions

**Optional Tools (detected if available):**

- **Homebrew** (macOS only) - Package manager for additional tool installations
- **Docker** - Container runtime (detected but not installed)
- **kubectl** - Kubernetes command-line tool (detected but not installed)

**Process:**

1. Checks if each tool is available in the system PATH
2. Attempts to install missing required tools using appropriate package managers
3. Reports the status of optional tools without failing if they're unavailable
4. Provides clear feedback on installation success or failure

### Stage 2: Directory Structure Creation

Creates the necessary directory structure for Anvil:

```
~/.anvil/
└── settings.yaml   # Main configuration file
```

**Directory Purpose:**

- **`~/.anvil/`** - Main configuration directory containing settings and configuration files

### Stage 3: Configuration File Generation

Generates a default `settings.yaml` configuration file with:

**System Detection:**

- Automatically detects your Git username and email from global Git configuration
- Sets up appropriate directory paths for your system
- Configures tool lists based on your platform

**Configuration Structure:**

```yaml
version: 1.0.0
directories:
  config: /Users/username/.anvil
tools:
  required_tools:
    - git
    - curl
  optional_tools:
    - brew
    - docker
    - kubectl
git:
  username: "Your Name"
  email: "your.email@example.com"
environment: {}
```

### Stage 4: Environment Configuration Validation

Performs comprehensive checks of your development environment:

**Git Configuration:**

- Verifies `git config --global user.name` is set
- Verifies `git config --global user.email` is set
- Provides specific commands to fix missing configuration

**SSH Key Management:**

- Checks for the existence of `~/.ssh/` directory
- Validates presence of common SSH key files (`id_rsa`, `id_ed25519`, `id_ecdsa`)
- Recommends SSH key generation if none are found

**Environment Variables:**

- Checks for common development environment variables (`EDITOR`, `SHELL`)
- Suggests setting missing variables for improved development experience

### Stage 5: Completion and Next Steps

Provides comprehensive feedback and guidance:

**Success Indicators:**

- Clear confirmation of successful initialization
- Display of configuration file location
- Summary of completed tasks

**Actionable Recommendations:**

- Specific commands to resolve any configuration warnings
- Suggested next steps for optimal Anvil usage
- Links to relevant documentation and resources

## Command Usage

### Basic Usage

```bash
anvil init
```

### Command Help

```bash
anvil init --help
```

### Expected Output

```
=== Anvil Initialization ===

🔧 Validating and installing required tools...
✓ Git is available
✓ cURL is available
✓ Homebrew is available
✓ Docker is available
✓ kubectl is available
✅ All required tools are available

🔧 Creating necessary directories...
✅ Directories created successfully

🔧 Generating default settings.yaml...
✅ Default settings.yaml generated

🔧 Checking local environment configurations...
⚠️  Environment configuration warnings:
⚠️    - Configure git user.name: git config --global user.name 'Your Name'
⚠️    - Set up SSH keys for GitHub: ssh-keygen -t ed25519 -C 'your.email@example.com'

=== Initialization Complete! ===

Anvil has been successfully initialized and is ready to use.
Configuration files have been created in: /Users/username/.anvil

Recommended next steps to complete your setup:
  • Configure git user.name: git config --global user.name 'Your Name'
  • Set up SSH keys for GitHub: ssh-keygen -t ed25519 -C 'your.email@example.com'

These steps are optional but recommended for the best experience.

You can now use:
  • 'anvil --help' to see all available commands
  • 'anvil install' to install development tools
  • 'anvil config pull', 'anvil config show', and 'anvil config push' to manage configuration files with GitHub
  • Edit ~/.anvil/settings.yaml to customize your configuration
```

## Key Features

### Idempotent Operation

The init command is designed to be **idempotent**, meaning:

- Running it multiple times is safe and won't cause issues
- Existing configurations are preserved and not overwritten
- Only missing components are created or installed
- Warnings are re-evaluated on each run to reflect current system state

### Error Handling

Comprehensive error handling ensures:

- Clear error messages with actionable guidance
- Graceful failure with specific remediation steps
- Permission checks with helpful suggestions
- Network connectivity validation for downloads

### Cross-Platform Support

The init command intelligently adapts to different operating systems:

- **macOS**: Full Homebrew integration for tool installation
- **Linux**: Appropriate package manager detection and usage
- **Windows**: Adapted tool detection and installation methods

### Progressive Enhancement

The initialization process follows a progressive enhancement approach:

- **Required tools** must be available for core functionality
- **Optional tools** enhance the experience but don't block initialization
- **Environment warnings** provide guidance without preventing usage

## Integration with Other Commands

The init command prepares your system for all other Anvil commands:

### Setup Command

- Relies on the tool configuration established by init
- Uses the directory structure created during initialization
- References the settings.yaml file for installation preferences

### Config Commands (Pull/Push)

- Depend on Git configuration validated during init
- Use SSH keys whose presence is checked during initialization
- Store configuration files in the main config directory

## Configuration Customization

After initialization, you can customize your Anvil configuration:

### Editing settings.yaml

```yaml
# Add custom environment variables
environment:
  CUSTOM_VAR: "value"
  DEVELOPMENT_MODE: "true"

# Modify tool lists
tools:
  required_tools:
    - git
    - curl
    - custom-tool
  optional_tools:
    - brew
    - docker
    - kubectl
    - your-preferred-tool
```

### Directory Customization

While not recommended, you can modify directory paths:

```yaml
directories:
  config: /custom/path/.anvil
```

## Troubleshooting

### Common Issues and Solutions

**Permission Errors:**

```bash
# Fix home directory permissions
chmod 755 ~/
mkdir -p ~/.anvil
chmod 755 ~/.anvil
```

**Tool Installation Failures:**

```bash
# On macOS, install Homebrew manually
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Update PATH if needed
export PATH="/opt/homebrew/bin:$PATH"
```

**Git Configuration Issues:**

```bash
# Set up basic Git configuration
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

**SSH Key Setup:**

```bash
# Generate new SSH key
ssh-keygen -t ed25519 -C "your.email@example.com"

# Add to SSH agent
ssh-add ~/.ssh/id_ed25519

# Copy public key to clipboard (macOS)
pbcopy < ~/.ssh/id_ed25519.pub
```

### Advanced Troubleshooting

**Debugging Mode:**

```bash
# Run with verbose output (if implemented)
anvil init --verbose

# Check system requirements
anvil init --check-only
```

**Manual Configuration:**
If automatic initialization fails, you can manually create the configuration:

```bash
# Create directory
mkdir -p ~/.anvil

# Create basic settings file
cat > ~/.anvil/settings.yaml << EOF
version: 1.0.0
directories:
  config: $HOME/.anvil
tools:
  required_tools:
    - git
    - curl
  optional_tools: []
git:
  username: ""
  email: ""
environment: {}
EOF
```

## Best Practices

### When to Run Init

- **After first installation** of Anvil
- **On new development machines** to establish consistency
- **After major system updates** that might affect tool availability
- **When experiencing configuration issues** with other Anvil commands

### Pre-initialization Checklist

Before running `anvil init`, ensure you have:

- [ ] Administrative privileges for tool installation
- [ ] Network connectivity for downloading tools
- [ ] Basic Git configuration (or willingness to set it up)
- [ ] Understanding of your development environment needs

### Post-initialization Recommendations

After successful initialization:

1. **Review the generated settings.yaml** and customize as needed
2. **Complete any recommended configuration steps** from the output
3. **Test other Anvil commands** to ensure proper setup
4. **Backup your configuration** for future reference

## Security Considerations

### Tool Installation Security

- Only installs tools from trusted package managers
- Verifies tool signatures when possible
- Uses official installation methods recommended by tool maintainers

### Configuration File Security

- Settings.yaml contains no sensitive information by default
- Uses standard Unix file permissions (644)
- Stores configuration in user's home directory for proper isolation

### Network Security

- All downloads use HTTPS when possible
- Tool installations follow official security guidelines
- No automatic execution of downloaded scripts without user consent

## Future Enhancements

### Planned Features

- **Configuration validation** - Verify settings.yaml format and content
- **Backup and restore** - Save and restore Anvil configurations
- **Template configurations** - Pre-defined configurations for different use cases
- **Plugin system** - Support for custom initialization steps
- **Team configurations** - Share initialization settings across team members

### Extensibility

The init command is designed to be extensible:

- New tools can be easily added to the validation process
- Additional configuration checks can be implemented
- Custom initialization steps can be integrated
- Platform-specific enhancements can be added

## Conclusion

The `anvil init` command is a critical component of the Anvil CLI ecosystem. It provides a robust, user-friendly way to establish a proper development environment foundation. By automating complex setup tasks and providing clear guidance, it ensures that users can quickly and confidently begin using Anvil for their development workflows.

The command's idempotent nature, comprehensive error handling, and progressive enhancement approach make it safe and reliable for users at all skill levels. Whether you're setting up a new development machine or troubleshooting configuration issues, the init command provides the foundation for a successful Anvil experience.
