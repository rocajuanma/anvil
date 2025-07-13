# Anvil Setup Command Documentation

## Overview

The `anvil setup` command provides automated batch installation of development tools and applications using Homebrew (on macOS) or other package managers. This command is designed to streamline the process of setting up consistent development environments across different machines by organizing tools into logical groups.

## Purpose and Importance

The setup command serves several critical functions:

- **Standardized Development Environments** - Ensures consistent tool configurations across team members and machines
- **Automated Tool Installation** - Eliminates the tedious manual process of installing development tools one by one
- **Group-based Organization** - Organizes tools into logical groups for different use cases (development, new laptop setup, etc.)
- **Selective Installation** - Allows installation of individual tools or entire groups based on needs
- **Configuration Management** - Reads from your Anvil configuration to maintain tool preferences

## Command Modes

### Group Installation Mode

Install all tools in a predefined group:

```bash
anvil setup [group-name]
```

### Individual Tool Mode

Install specific tools using flags:

```bash
anvil setup --git --zsh --vscode
```

### Utility Mode

List available groups or perform dry runs:

```bash
anvil setup --list
anvil setup --dry-run
```

## Available Groups

### Default Groups

The setup command comes with two predefined groups:

#### 1. Development Group (`dev`)

**Purpose**: Essential tools for software development
**Tools included**:

- **git** - Version control system
- **zsh** - Advanced shell with oh-my-zsh configuration
- **iterm2** - Enhanced terminal emulator for macOS
- **vscode** - Visual Studio Code editor

**Usage**:

```bash
anvil setup dev
```

#### 2. New Laptop Group (`new-laptop`)

**Purpose**: Essential applications for setting up a new laptop
**Tools included**:

- **slack** - Team communication platform
- **chrome** - Google Chrome web browser
- **1password** - Password manager

**Usage**:

```bash
anvil setup new-laptop
```

### Custom Groups

You can define custom groups in your `settings.yaml` file:

```yaml
groups:
  custom:
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
```

## Individual Tool Installation

### Available Tool Flags

The setup command supports individual installation of the following tools:

| Flag          | Tool               | Description                                      |
| ------------- | ------------------ | ------------------------------------------------ |
| `--git`       | Git                | Version control system with global configuration |
| `--zsh`       | Zsh                | Advanced shell with oh-my-zsh and plugins        |
| `--iterm2`    | iTerm2             | Enhanced terminal emulator for macOS             |
| `--vscode`    | Visual Studio Code | Popular code editor                              |
| `--slack`     | Slack              | Team communication platform                      |
| `--chrome`    | Google Chrome      | Web browser                                      |
| `--1password` | 1Password          | Password manager                                 |

### Examples

**Install only Git:**

```bash
anvil setup --git
```

**Install multiple individual tools:**

```bash
anvil setup --git --zsh --vscode
```

**Install with dry run to preview:**

```bash
anvil setup --git --zsh --dry-run
```

## Command Options and Flags

### Utility Flags

- `--list` - List all available groups and their tools
- `--dry-run` - Show what would be installed without actually installing
- `--help` - Show command help and usage information

### Usage Examples

**List all available groups:**

```bash
anvil setup --list
```

**Preview what would be installed:**

```bash
anvil setup dev --dry-run
```

**Show help information:**

```bash
anvil setup --help
```

## Tool Installation Details

### Git Installation

- Installs Git via Homebrew if not already present
- Preserves existing Git configuration
- Works with global Git settings configured during `anvil init`

### Zsh Installation

- Installs Zsh shell via Homebrew
- Automatically installs oh-my-zsh framework
- Configures sensible defaults for enhanced terminal experience
- Installs in unattended mode to prevent interactive prompts

### iTerm2 Installation

- Installs iTerm2 terminal emulator as a cask
- Provides enhanced terminal features over default Terminal.app
- Supports themes, profiles, and advanced terminal functionality

### Visual Studio Code Installation

- Installs VS Code as a cask application
- Provides a full-featured code editor
- Supports extensions and customization

### Slack Installation

- Installs Slack desktop application
- Enables team communication and collaboration
- Integrates with macOS notifications and features

### Google Chrome Installation

- Installs Chrome web browser
- Provides modern web browsing capabilities
- Supports development tools and extensions

### 1Password Installation

- Installs 1Password password manager
- Provides secure password storage and management
- Integrates with browsers and applications

## Expected Output

### Group Installation Output

```
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

### Individual Tool Installation Output

```
=== Individual Tool Setup ===

Installing individual tools: git, zsh
[1/2] 50% - Installing git
✅ git installed successfully
[2/2] 100% - Installing zsh
Installing oh-my-zsh...
✅ zsh installed successfully

=== Individual Tool Setup Complete! ===
```

### Dry Run Output

```
Dry run mode - no actual installations will be performed

=== Setting up 'dev' group ===

Installing tools for group 'dev': git, zsh, iterm2, vscode
[1/4] 25% - Installing git
Would install: git
[2/4] 50% - Installing zsh
Would install: zsh
[3/4] 75% - Installing iterm2
Would install: iterm2
[4/4] 100% - Installing vscode
Would install: vscode

=== Group Setup Complete! ===

Successfully installed 4 of 4 tools in group 'dev'
```

## Configuration Management

### settings.yaml Structure

The setup command reads group configurations from `~/.anvil/settings.yaml`:

```yaml
groups:
  dev:
    - git
    - zsh
    - iterm2
    - vscode
  new-laptop:
    - slack
    - chrome
    - 1password
  custom:
    frontend:
      - git
      - node
      - yarn
      - chrome
    backend:
      - git
      - docker
      - postgresql
```

### Adding Custom Groups

You can add custom groups by editing the `settings.yaml` file:

```yaml
groups:
  custom:
    data-science:
      - python
      - jupyter
      - pandas
      - matplotlib
    mobile-dev:
      - xcode
      - android-studio
      - flutter
```

After adding custom groups, they become available for installation:

```bash
anvil setup data-science
anvil setup mobile-dev
```

## Platform Support

### macOS (Primary Platform)

- Full support for all tools via Homebrew
- Optimized for macOS-specific applications
- Integrates with macOS package management

### Other Platforms

- Limited support with warnings displayed
- Basic tools may work on Linux/Windows
- Platform-specific tools may not be available

## Error Handling

### Common Error Scenarios

**Tool Installation Failure:**

```
❌ Failed to install vscode: package not found
```

**Permission Issues:**

```
❌ Failed to install git: permission denied
```

**Network Connectivity:**

```
❌ Failed to install slack: network timeout
```

### Error Recovery

The setup command continues installing other tools even if some fail:

```
=== Group Setup Complete! ===

Successfully installed 3 of 4 tools in group 'dev'
⚠️  Some tools failed to install. Check the output above for details.
```

## Best Practices

### Pre-Setup Checklist

Before running setup commands:

- [ ] Ensure you have administrative privileges
- [ ] Run `anvil init` to initialize your configuration
- [ ] Verify network connectivity
- [ ] Check available disk space for installations

### Recommended Workflow

1. **Start with initialization:**

   ```bash
   anvil init
   ```

2. **Review available groups:**

   ```bash
   anvil setup --list
   ```

3. **Test with dry run:**

   ```bash
   anvil setup dev --dry-run
   ```

4. **Install your desired group:**

   ```bash
   anvil setup dev
   ```

5. **Add individual tools as needed:**
   ```bash
   anvil setup --slack --chrome
   ```

### Group Organization Strategy

**For Teams:**

- Create team-specific groups in settings.yaml
- Share common configurations across team members
- Use descriptive group names (frontend, backend, qa, etc.)

**For Personal Use:**

- Create role-specific groups (work, personal, experiments)
- Organize by project type or technology stack
- Keep groups focused and manageable

## Advanced Usage

### Combining with Other Commands

**Complete development environment setup:**

```bash
# Initialize Anvil
anvil init

# Install development tools
anvil setup dev

# Configure additional settings
anvil setup --git --zsh

# Sync with GitHub (if applicable)
anvil pull configs
```

### Custom Tool Installation

For tools not included in the predefined lists, you can:

1. **Add to custom groups in settings.yaml**
2. **Install directly via Homebrew integration**
3. **Extend the setup command with custom installers**

### Automation and Scripting

The setup command is designed to be scriptable:

```bash
#!/bin/bash
# New machine setup script

echo "Setting up new development machine..."
anvil init
anvil setup new-laptop
anvil setup dev
echo "Setup complete!"
```

## Troubleshooting

### Common Issues

**Homebrew Not Installed:**

```bash
# Run init first to install Homebrew
anvil init
```

**Tool Already Installed:**

- The command safely skips already installed tools
- No action needed - this is expected behavior

**Permission Denied:**

```bash
# Fix Homebrew permissions
sudo chown -R $(whoami) /usr/local/Homebrew
```

**Network Issues:**

```bash
# Check network connectivity
ping -c 3 github.com
brew update
```

### Debugging

**Verbose Output:**

- Use `--dry-run` to see what would be installed
- Check individual tool installation with single flags
- Review Homebrew logs for detailed error information

**Manual Installation:**
If automated installation fails, you can install tools manually:

```bash
# Install via Homebrew directly
brew install git
brew install --cask visual-studio-code
```

## Security Considerations

### Tool Source Verification

- All tools are installed from official Homebrew formulae
- Cask applications are verified through Homebrew's curation process
- No custom or unofficial packages are installed automatically

### Permission Management

- Installations use standard user permissions
- No sudo required for most operations
- Follows macOS security guidelines

### Network Security

- All downloads use HTTPS connections
- Tool signatures are verified by Homebrew
- No execution of arbitrary remote scripts

## Performance Considerations

### Installation Speed

- Tools are installed sequentially for stability
- Progress indicators show installation status
- Large applications may take several minutes

### Resource Usage

- Disk space requirements vary by tool
- Some tools (like Xcode) require significant space
- Check available disk space before installing large groups

### Parallel Installation

- Currently uses sequential installation for reliability
- Future versions may support parallel installation
- Use individual tool flags for faster targeted installation

## Future Enhancements

### Planned Features

- **Parallel Installation** - Install multiple tools simultaneously
- **Version Management** - Specify tool versions in groups
- **Dependency Resolution** - Automatic handling of tool dependencies
- **Update Management** - Update tools and groups to latest versions
- **Template Groups** - Predefined groups for common scenarios

### Extensibility

The setup command is designed to be extensible:

- New tools can be easily added to installation logic
- Custom installers can be integrated
- Platform-specific enhancements can be implemented
- Community-contributed groups can be supported

## Conclusion

The `anvil setup` command provides a powerful and flexible way to install and manage development tools. Its group-based approach, combined with individual tool installation options, makes it suitable for both individual developers and teams. The command's integration with Anvil's configuration system ensures consistency across different environments while providing the flexibility needed for diverse development workflows.

Whether you're setting up a new development machine, onboarding team members, or maintaining consistent tool configurations, the setup command provides the automation and reliability needed for modern development environments.
