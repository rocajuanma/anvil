# Doctor Command

The `anvil doctor` command provides comprehensive health checks to validate your development environment and troubleshoot common configuration issues with **real-time progress feedback**.

## Overview

The doctor command performs systematic validation across four key areas with **live progress indicators** so you always know what's happening. You can run checks at different levels of granularity:

### 🏷️ **Categories** (groups of related checks)

- **environment** - Verify anvil initialization and directory structure (3 checks)
- **dependencies** - Check required tools and Homebrew installation (2 checks)
- **configuration** - Validate git and GitHub settings (3 checks)
- **connectivity** - Test GitHub access and repository connections (3 checks)

### 🔍 **Specific Checks** (individual validators)

Run `anvil doctor --list` to see all 11 available individual checks like `git-config`, `homebrew`, `github-auth`, etc.

## Key Features

### ✨ **Real-time Progress**

- **Live feedback** - See progress as each validation runs
- **Progress counters** - Know exactly how many checks are remaining
- **Stage indicators** - Understand what phase the doctor is in
- **No more hanging terminals** - Always know what's happening

### 🔍 **Two Output Modes**

- **Default mode** - Brief but informative progress messages
- **Verbose mode (`--verbose`)** - Detailed descriptions, categories, and step-by-step results

### 🔒 **Secure Authentication**

- **Non-interactive operations** - No credential prompts ever
- **Environment-based auth** - Uses `GITHUB_TOKEN` or SSH keys only
- **Private repository enforcement** - Clear warnings about public repos

## Usage

### Basic Commands

```bash
# Run all health checks (11 total) with real-time progress
anvil doctor

# List available categories and checks with explanations
anvil doctor --list

# Run all checks in a category with progress feedback
anvil doctor environment        # 3 environment checks
anvil doctor dependencies       # 2 dependency checks
anvil doctor configuration      # 3 configuration checks
anvil doctor connectivity       # 3 connectivity checks

# Run a specific individual check with detailed feedback
anvil doctor git-config
anvil doctor homebrew
anvil doctor github-auth

# Show detailed output with descriptions and step-by-step results
anvil doctor --verbose
anvil doctor environment --verbose
```

### Auto-Fix Functionality

```bash
# Auto-fix all fixable issues
anvil doctor --fix

# Auto-fix issues in a specific category
anvil doctor dependencies --fix

# Auto-fix a specific check
anvil doctor homebrew --fix
```

## Understanding Categories vs Specific Checks

**Categories** are groups of related checks that test a particular area:

- When you run `anvil doctor environment`, it runs 3 checks: `anvil-init`, `settings-valid`, and `directory-structure`
- When you run `anvil doctor dependencies`, it runs 2 checks: `homebrew` and `required-tools`

**Specific checks** are individual validators that test one particular thing:

- When you run `anvil doctor git-config`, it only checks if your git configuration is properly set
- When you run `anvil doctor homebrew`, it only checks if Homebrew is installed and functional

Use categories when you want to check an entire area, and use specific checks when you want to focus on one particular aspect.

## Health Check Categories

### Environment Checks

| Check                 | Description                                     | Auto-Fix |
| --------------------- | ----------------------------------------------- | -------- |
| `anvil-init`          | Verify anvil initialization has been completed  | ❌       |
| `settings-valid`      | Validate settings.yaml structure and content    | ❌       |
| `directory-structure` | Check ~/.anvil directory structure              | ❌       |

### Dependencies Checks

| Check            | Description                                         | Auto-Fix |
| ---------------- | --------------------------------------------------- | -------- |
| `homebrew`       | Verify Homebrew installation and updates            | ✅       |
| `required-tools` | Check git and curl are installed                    | ❌       |

### Configuration Checks

| Check           | Description                                  | Auto-Fix |
| --------------- | -------------------------------------------- | -------- |
| `git-config`    | Validate git user.name and user.email        | ✅       |
| `github-config` | Verify GitHub repository configuration       | ❌       |
| `sync-config`   | Check config sync settings (not implemented) | ❌       |

### Connectivity Checks

| Check             | Description                              | Auto-Fix |
| ----------------- | ---------------------------------------- | -------- |
| `github-auth`     | Test GitHub authentication and access    | ❌       |
| `github-repo`     | Verify repository accessibility          | ❌       |
| `git-operations`  | Test git clone and pull operations       | ❌       |

## Check Results

Each check returns one of four statuses:

- **✅ PASS** - Check completed successfully
- **⚠️ WARN** - Check passed but has warnings or recommendations
- **❌ FAIL** - Check failed and requires attention
- **⏭️ SKIP** - Check was skipped (usually due to missing configuration)

## Output Examples

### Full Health Check with Real-time Progress

```bash
$ anvil doctor

=== Running Anvil Health Check ===

🔍 Validating environment, dependencies, configuration, and connectivity...

🔧 Executing 12 health checks...
[1/12] 8% - Running init-run
   ✅ Anvil initialization complete
[2/12] 17% - Running settings-file
   ✅ Settings file is valid
[3/12] 25% - Running directory-structure
   ✅ Directory structure is correct
[4/12] 33% - Running homebrew
   ✅ Homebrew is installed and functional
[5/12] 42% - Running required-tools
   ✅ All required tools installed (2/2)
[6/12] 50% - Running optional-tools
   ⚠️ Optional tools missing: docker
[7/12] 58% - Running git-config
   ✅ Git configuration is complete
[8/12] 67% - Running github-config
   ✅ GitHub configuration is complete
[9/12] 75% - Running sync-config
   ⏭️ No sync configuration present (optional)
[10/12] 83% - Running github-access
   ✅ GitHub SSH access confirmed
[11/12] 92% - Running repository-access
   ✅ ✅ Private repository accessible with proper authentication
[12/12] 100% - Running git-connectivity
   ✅ Git operations are functional

✅ All validation checks completed

✅ Environment
  ✅ Anvil initialization complete
  ✅ Settings file is valid
  ✅ Directory structure is correct

⚠️ Dependencies
  ✅ Homebrew is installed and functional
  ✅ All required tools installed (2/2)
  ⚠️ Optional tools missing: docker
      💡 Run 'anvil install docker' to install missing tools

✅ Configuration
  ✅ Git configuration is complete
  ✅ GitHub configuration is complete
  ⏭️ No sync configuration present (optional)

✅ Connectivity
  ✅ GitHub SSH access confirmed
  ✅ ✅ Private repository accessible with proper authentication
  ✅ Git operations are functional

=== Health Check Summary ===

Total checks: 12
✅ Passed: 10
⚠️ Warnings: 1
⏭️ Skipped: 1

✅ Overall status: Healthy
```

### Category-Specific Check with Verbose Output

```bash
$ anvil doctor connectivity --verbose

=== Running Connectivity Health Checks ===

🔧 Executing 3 checks in connectivity category...
[1/3] 33% - Running github-access
   Description: Verify GitHub API access and authentication
   Category: connectivity
   Result: ✅ GitHub SSH access confirmed
      Repository: rocajuanma/configs
      Token environment variable: GITHUB_TOKEN
      ✗ No GitHub token found in environment
      Attempting SSH authentication...
      ✓ SSH authentication successful
[2/3] 67% - Running repository-access
   Description: Verify configured repository exists and is accessible
   Category: connectivity
   Result: ✅ ✅ Private repository accessible with proper authentication
      Repository: rocajuanma/configs
      🔒 Repository is private (secure)
      🔑 Git authentication working
      🛡️  Configuration data is protected
[3/3] 100% - Running git-connectivity
   Description: Verify git operations are functional
   Category: connectivity
   Result: ✅ Git operations are functional
      git version 2.39.3 (Apple Git-145)
      Global username: Juanma Roca
      Global email: juanma.roca@zapier.com

✅ All validation checks completed

✅ Connectivity
  ✅ GitHub SSH access confirmed
      Repository: rocajuanma/configs
      Token environment variable: GITHUB_TOKEN
      ✗ No GitHub token found in environment
      Attempting SSH authentication...
      ✓ SSH authentication successful
  ✅ ✅ Private repository accessible with proper authentication
      Repository: rocajuanma/configs
      🔒 Repository is private (secure)
      🔑 Git authentication working
      🛡️  Configuration data is protected
  ✅ Git operations are functional
      git version 2.39.3 (Apple Git-145)
      Global username: Juanma Roca
      Global email: juanma.roca@zapier.com

=== Health Check Summary ===

Total checks: 3
✅ Passed: 3
✅ Overall status: Healthy
```

### Single Check with Detailed Output

```bash
$ anvil doctor github-access --verbose

=== Running Check: github-access ===

🔧 Executing github-access check...
🔍 Running github-access check...
   Description: Verify GitHub API access and authentication
   Category: connectivity
✅ GitHub SSH access confirmed

✅ GitHub SSH access confirmed
      Repository: rocajuanma/configs
      Token environment variable: GITHUB_TOKEN
      ✗ No GitHub token found in environment
      Attempting SSH authentication...
      ✓ SSH authentication successful
```

## Common Issues and Solutions

### Environment Issues

**Settings file not found**

```bash
# Solution: Run anvil init
anvil init
```

**Directory structure incomplete**

```bash
# Solution: Auto-fix creates missing directories
anvil doctor directory-structure --fix
```

### Dependency Issues

**Homebrew not installed**

```bash
# Solution: Auto-fix installs Homebrew
anvil doctor homebrew --fix
```

**Required tools missing**

```bash
# Solution: Auto-fix installs missing tools
anvil doctor required-tools --fix
```

### Configuration Issues

**Git configuration incomplete**

```bash
# Solution: Manually configure in settings.yaml
# Edit ~/.anvil/settings.yaml:
git:
  username: "Your Name"
  email: "your.email@example.com"
```

**GitHub configuration incomplete**

```bash
# Solution: Configure GitHub repository in settings.yaml
# Edit ~/.anvil/settings.yaml:
github:
  config_repo: "username/repository"
  branch: "main"
  token_env_var: "GITHUB_TOKEN"
```

### Connectivity Issues

**GitHub authentication failed**

```bash
# Solution: Set up GitHub token or SSH keys
export GITHUB_TOKEN="your_token_here"

# Or configure SSH keys:
ssh-keygen -t ed25519 -C "your.email@example.com"
# Add to GitHub: https://github.com/settings/keys
```

**Repository not accessible**

```bash
# Solution: Check repository name and permissions
# Ensure repository exists and you have access
# Update config_repo in settings.yaml if needed
```

## Integration with Other Commands

The doctor command integrates seamlessly with other anvil commands:

```bash
# After init, verify setup
anvil init
anvil doctor

# Before config operations, check connectivity
anvil doctor connectivity
anvil config pull

# After installation, verify dependencies
anvil install dev
anvil doctor dependencies
```

## Best Practices

1. **Run after initialization**: Always run `anvil doctor` after `anvil init` to verify setup
2. **Use real-time feedback**: Watch the progress indicators to understand timing
3. **Pre-flight checks**: Use category-specific checks before major operations
4. **Regular maintenance**: Run periodic health checks to catch configuration drift
5. **Auto-fix safely**: Review auto-fix suggestions before applying them
6. **Verbose for debugging**: Use `--verbose` for detailed troubleshooting information
7. **Monitor authentication**: Watch for authentication method details in verbose mode

## Troubleshooting

If the doctor command itself fails:

```bash
# Check basic anvil functionality
anvil --help

# Verify Go installation
go version

# Rebuild anvil
go build -o anvil main.go
```

For issues with specific validators, check the relevant package documentation and ensure all dependencies are properly installed. The verbose mode (`--verbose`) provides detailed execution information for debugging.

## Security Notes

- **No credential prompts**: All git operations are non-interactive
- **Environment-based auth**: Uses `GITHUB_TOKEN` environment variable or SSH keys
- **Private repository validation**: Warns about public repositories for security
- **Secure by default**: All authentication methods prioritize security
