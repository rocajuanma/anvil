# Import Groups

The `anvil config import` command allows you to import tool group definitions from local files or remote URLs into your anvil configuration. This feature is particularly useful for teams that want to share standardized tool groups across their development environments.

## Overview

The import functionality provides a secure and validated way to add new tool groups to your anvil configuration without manually editing configuration files. It supports both local files and remote URLs, with comprehensive validation and conflict detection.

### Key Features

- **🔗 Flexible Sources**: Import from local files or publicly accessible URLs
- **✅ Comprehensive Validation**: Validates group names, application names, and structure
- **🚫 Conflict Detection**: Prevents overwriting existing groups with clear error messages
- **🌳 Tree Display**: Shows visual preview of groups and applications before import
- **📋 Interactive Confirmation**: Requires user approval before making changes
- **🛡️ Security-First**: Only imports group definitions, ignoring sensitive configuration data

## Usage

### Basic Syntax

```bash
anvil config import [file-or-url]
```

### Examples

#### Import from Local File

```bash
# Import from a local YAML file
anvil config import ./team-groups.yaml

# Import from a file in your home directory
anvil config import ~/Downloads/company-groups.yaml

# Import from a file with absolute path
anvil config import /path/to/shared-groups.yaml
```

#### Import from Remote URL

```bash
# Import from a public GitHub repository
anvil config import https://raw.githubusercontent.com/company/shared-configs/main/groups.yaml

# Import from any publicly accessible URL
anvil config import https://example.com/team-tools.yaml

# Import from a company internal URL
anvil config import https://internal.company.com/configs/development-groups.yaml
```

#### Quick Start with Example Configurations

Anvil includes a comprehensive set of example configurations to help you get started quickly. These examples are designed for different developer personas and can serve as templates for your own custom groups.

```bash
# Import frontend developer setup
anvil config import import-examples/frontend-developer.yaml

# Import backend developer configuration
anvil config import import-examples/backend-developer.yaml

# Import data scientist tools
anvil config import import-examples/data-scientist.yaml

# Import DevOps engineer setup
anvil config import import-examples/devops-engineer.yaml

# Import designer tools
anvil config import import-examples/designer.yaml

# Import startup founder configuration
anvil config import import-examples/startup-founder.yaml

# Import team configuration for small startups
anvil config import import-examples/team-startup.yaml
```

**Available Example Personas:**

- **🎨 Frontend Developer**: Modern web development tools and design applications
- **⚙️ Backend Developer**: Server-side technologies, databases, and DevOps tools
- **📊 Data Scientist**: Data analysis, machine learning, and visualization tools
- **🚀 DevOps Engineer**: Infrastructure, cloud tools, and monitoring stack
- **🎭 Designer**: UI/UX design tools and prototyping applications
- **🏢 Startup Founder**: Comprehensive setup for technical founders
- **👥 Team Startup**: Multi-role configuration for small development teams

> **💡 Pro Tip**: These examples are starting points only! Feel free to modify them, combine multiple configurations, or create your own custom groups that better fit your specific workflow and needs.

## File Format

The import file must be a valid YAML file containing a `groups` section. The structure should follow this format:

```yaml
groups:
  group-name:
    - tool1
    - tool2
    - tool3
  another-group:
    - tool4
    - tool5
```

### Example Import File

```yaml
groups:
  backend-dev:
    - docker
    - postgresql
    - redis
    - nodejs
    - git
  frontend-dev:
    - nodejs
    - npm
    - chrome
    - vscode
    - git
  devops:
    - docker
    - kubernetes
    - terraform
    - aws-cli
    - git
```

## Import Process

The import command follows a structured process to ensure data integrity and user control:

### Stage 1: File Fetching
- **Local Files**: Validates file existence and accessibility
- **Remote URLs**: Downloads file with timeout protection (30 seconds)
- **Temporary Storage**: Creates secure temporary file for processing

### Stage 2: Parsing and Validation
- **YAML Parsing**: Validates YAML syntax and structure
- **Group Extraction**: Extracts only the `groups` section from the file
- **Structure Validation**: Ensures proper group and tool name formatting

### Stage 3: Conflict Detection
- **Existing Groups**: Checks for conflicts with current configuration
- **Clear Errors**: Provides specific error messages for conflicting group names
- **Safe Import**: Prevents accidental overwrites of existing configurations

### Stage 4: Preview and Confirmation
- **Visual Summary**: Displays tree structure of groups and tools to be imported
- **Statistics**: Shows total groups and applications count
- **User Confirmation**: Requires explicit approval before proceeding

### Stage 5: Import Execution
- **Safe Addition**: Adds new groups to existing configuration
- **Configuration Save**: Persists changes to settings.yaml
- **Success Confirmation**: Provides clear feedback on import completion

## Example Workflow

Here's a complete example of importing team groups:

```bash
$ anvil config import https://raw.githubusercontent.com/company/shared-configs/main/team-groups.yaml

=== Import Groups from File ===
🔧 Fetching source file...
✅ Source file fetched successfully
🔧 Parsing import file...
✅ Import file parsed successfully
🔧 Validating group structure...
✅ Group structure validation passed
🔧 Checking for conflicts...
✅ No conflicts detected
🔧 Preparing import summary...

📋 Import Summary:
═══════════════════
├── 📁 backend-dev (5 tools)
│   ├── 🔧 docker
│   ├── 🔧 git
│   ├── 🔧 nodejs
│   ├── 🔧 postgresql
│   └── 🔧 redis
│
├── 📁 frontend-dev (5 tools)
│   ├── 🔧 chrome
│   ├── 🔧 git
│   ├── 🔧 nodejs
│   ├── 🔧 npm
│   └── 🔧 vscode
│
📊 Total: 2 groups, 10 applications

? Proceed with importing these groups? (y/N): y
🔧 Importing groups...
✅ Groups imported successfully

✨ Import completed! 2 groups have been added to your configuration.
```

## Error Handling

The import command provides comprehensive error handling for various scenarios:

### File Access Errors

```bash
$ anvil config import ./nonexistent.yaml
Import failed: file does not exist: ./nonexistent.yaml
```

### Network Errors

```bash
$ anvil config import https://invalid-url.com/groups.yaml
Import failed: failed to download file: Get "https://invalid-url.com/groups.yaml": dial tcp: lookup invalid-url.com: no such host
```

### YAML Parsing Errors

```bash
$ anvil config import ./invalid-yaml.yaml
Import failed: failed to parse YAML: yaml: line 2: found character that cannot start any token
```

### Validation Errors

```bash
$ anvil config import ./invalid-groups.yaml
Import failed: invalid group name 'invalid group name': group name contains invalid characters
```

### Conflict Detection

```bash
$ anvil config import ./conflicting-groups.yaml
Import failed: groups already exist: backend-dev, frontend-dev
```

## Security Considerations

### What Gets Imported

- ✅ **Group Names**: Tool group identifiers
- ✅ **Tool Names**: Application names within groups
- ❌ **Sensitive Data**: No API keys, tokens, or personal information
- ❌ **System Paths**: No file paths or system-specific configurations
- ❌ **Authentication**: No credentials or authentication data

### Network Security

- **HTTPS Only**: Remote imports should use HTTPS URLs
- **Timeout Protection**: 30-second timeout prevents hanging requests
- **User Agent**: Identifies requests as coming from anvil-cli
- **Temporary Files**: Secure cleanup of downloaded files

### Data Validation

- **Group Name Validation**: Ensures valid group identifiers
- **Tool Name Validation**: Validates application names
- **Structure Validation**: Confirms proper YAML structure
- **Conflict Prevention**: Prevents overwriting existing configurations

## Customization and Best Practices

### Creating Your Own Groups

While the example configurations provide excellent starting points, you should customize them to match your specific needs:

#### Start with Examples, Then Customize

```bash
# 1. Import an example configuration
anvil config import import-examples/frontend-developer.yaml

# 2. View your current groups
anvil config show

# 3. Install tools from imported groups
anvil install frontend-core

# 4. Create additional custom groups as needed
# (Edit settings.yaml or use anvil commands)
```

#### Custom Group Creation Strategies

**Option 1: Modify Existing Groups**
- Import an example configuration
- Edit `~/.anvil/settings.yaml` to add/remove tools from groups
- Save and test your changes

**Option 2: Create New Custom Groups**
- Import multiple example configurations
- Create new groups that combine tools from different examples
- Organize by workflow rather than role

**Option 3: Team-Specific Groups**
- Start with team examples
- Add company-specific tools and workflows
- Share custom configurations with your team

#### Example Customization Workflow

```bash
# Start with a base configuration
anvil config import import-examples/backend-developer.yaml

# Add additional tools to existing groups
# (Edit settings.yaml to add tools like 'mongodb', 'elasticsearch')

# Create new custom groups
# (Add groups like 'microservices', 'testing', 'monitoring')

# Test your custom setup
anvil install backend-core
anvil install microservices
```

### Best Practices

#### For Team Leaders

1. **Start with Examples**: Use the provided examples as templates for team standards
2. **Centralized Groups**: Maintain group definitions in a shared repository
3. **Version Control**: Use version-controlled files for group definitions
4. **Documentation**: Document the purpose of each group and when to use them
5. **Regular Updates**: Keep group definitions current with team needs
6. **Customization Guidelines**: Provide guidelines for team members to customize their setups

#### For Team Members

1. **Start Simple**: Begin with an example configuration that matches your role
2. **Verify Sources**: Only import from trusted sources
3. **Review Before Import**: Always review the import summary
4. **Backup Configuration**: Consider backing up settings.yaml before major imports
5. **Test After Import**: Verify imported groups work as expected
6. **Iterate and Improve**: Continuously refine your groups based on your workflow

#### File Organization

```bash
# Recommended file structure for shared groups
shared-configs/
├── examples/                    # Team-specific examples
│   ├── senior-developer.yaml
│   ├── junior-developer.yaml
│   └── qa-engineer.yaml
├── custom/                      # Custom team groups
│   ├── microservices.yaml
│   ├── data-pipeline.yaml
│   └── security-tools.yaml
├── README.md
└── CHANGELOG.md
```

### Advanced Usage Patterns

#### Combining Multiple Configurations

```bash
# Import base configuration
anvil config import import-examples/backend-developer.yaml

# Import additional specialized tools
anvil config import import-examples/data-scientist.yaml

# Now you have both backend and data science tools available
anvil install backend-core
anvil install data-analysis
```

#### Role-Based Team Setup

```bash
# Each team member imports their role-specific configuration
# Frontend developers
anvil config import import-examples/frontend-developer.yaml

# Backend developers  
anvil config import import-examples/backend-developer.yaml

# DevOps engineers
anvil config import import-examples/devops-engineer.yaml

# Plus shared team configuration
anvil config import team-shared-tools.yaml
```

## Troubleshooting

### Common Issues

#### Import File Not Found
```bash
# Check file path
ls -la ./team-groups.yaml

# Use absolute path if needed
anvil config import /full/path/to/team-groups.yaml
```

#### Network Connectivity Issues
```bash
# Test URL accessibility
curl -I https://raw.githubusercontent.com/company/shared-configs/main/groups.yaml

# Check network connectivity
ping github.com
```

#### YAML Format Issues
```bash
# Validate YAML syntax
python -c "import yaml; yaml.safe_load(open('groups.yaml'))"

# Check for proper indentation
cat -A groups.yaml
```

#### Permission Issues
```bash
# Check file permissions
ls -la ./team-groups.yaml

# Fix permissions if needed
chmod 644 ./team-groups.yaml
```

### Getting Help

If you encounter issues with the import command:

1. **Check File Format**: Ensure your YAML file follows the correct structure
2. **Verify Network Access**: For remote URLs, ensure internet connectivity
3. **Review Error Messages**: Error messages provide specific guidance
4. **Check Documentation**: Refer to this guide for troubleshooting steps

## Integration with Other Commands

After importing groups, you can use them with other anvil commands:

```bash
# Import groups from examples
anvil config import import-examples/frontend-developer.yaml

# Install tools from imported groups
anvil install frontend-core
anvil install design-tools

# List available groups
anvil install --list

# Show group contents
anvil config show frontend-core

# Install multiple groups at once
anvil install essentials frontend-core productivity
```

### Complete Workflow Example

Here's a typical workflow for a new team member:

```bash
# 1. Initialize anvil
anvil init

# 2. Import role-specific configuration
anvil config import import-examples/backend-developer.yaml

# 3. Install essential tools
anvil install essentials

# 4. Install role-specific tools
anvil install backend-core

# 5. Install productivity tools
anvil install productivity

# 6. Verify installation
anvil doctor

# 7. Set up configuration sync (optional)
anvil config push anvil
```

### Team Onboarding Workflow

For teams setting up standardized environments:

```bash
# Team leader creates custom configuration
# (Based on import-examples/team-startup.yaml)

# Each team member runs:
anvil config import https://raw.githubusercontent.com/company/shared-configs/main/team-groups.yaml
anvil install team-essentials
anvil install [role-specific-group]

# Optional: Set up configuration sync
anvil config push anvil
```

This integration makes the import command a powerful tool for team collaboration and standardized development environment setup. The example configurations provide an excellent foundation that teams can build upon and customize for their specific needs.

## Quick Reference

### Example Configurations

| Persona | File | Description |
|---------|------|-------------|
| 🎨 Frontend Developer | `import-examples/frontend-developer.yaml` | Modern web development tools |
| ⚙️ Backend Developer | `import-examples/backend-developer.yaml` | Server-side technologies |
| 📊 Data Scientist | `import-examples/data-scientist.yaml` | Data analysis and ML tools |
| 🚀 DevOps Engineer | `import-examples/devops-engineer.yaml` | Infrastructure and deployment |
| 🎭 Designer | `import-examples/designer.yaml` | UI/UX design tools |
| 🏢 Startup Founder | `import-examples/startup-founder.yaml` | Technical founder setup |
| 👥 Team Startup | `import-examples/team-startup.yaml` | Multi-role team configuration |

### Quick Start Commands

```bash
# Get started quickly with your role
anvil config import import-examples/[your-role].yaml
anvil install essentials
anvil install [role-specific-group]

# For teams
anvil config import import-examples/team-startup.yaml
anvil install team-essentials
```

### Remember

- ✅ **Examples are starting points** - customize them for your needs
- ✅ **Combine configurations** - import multiple examples and create custom groups
- ✅ **Share with teams** - use examples as templates for team standards
- ✅ **Iterate and improve** - continuously refine your groups based on your workflow
