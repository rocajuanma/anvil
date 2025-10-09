# Import Examples

This directory contains example group configurations for different developer personas. These files can be used with the `anvil config import` command to quickly set up tool groups for various roles and workflows.

## Available Personas

### 🎨 Frontend Developer (`frontend-developer.yaml`)
Perfect for web developers working with modern frontend technologies.
- **frontend-core**: JavaScript/TypeScript ecosystem
- **design-tools**: Design and prototyping applications
- **productivity**: Communication and organization tools

### ⚙️ Backend Developer (`backend-developer.yaml`)
Ideal for server-side developers and API builders.
- **backend-core**: Server technologies and databases
- **devops**: Infrastructure and deployment tools
- **productivity**: Development and collaboration tools

### 📊 Data Scientist (`data-scientist.yaml`)
Tailored for data analysis and machine learning professionals.
- **data-analysis**: Data processing and visualization tools
- **ml-ops**: Machine learning and deployment tools
- **productivity**: Collaboration and documentation tools

### 🚀 DevOps Engineer (`devops-engineer.yaml`)
Designed for infrastructure and deployment specialists.
- **infrastructure**: Container orchestration and automation
- **cloud-tools**: Multi-cloud management tools
- **monitoring**: Observability and monitoring stack

### 🎭 Designer (`designer.yaml`)
Perfect for UI/UX designers and creative professionals.
- **design-core**: Primary design applications
- **prototyping**: Interactive prototyping tools
- **productivity**: Collaboration and file management

### 🏢 Startup Founder (`startup-founder.yaml`)
Comprehensive setup for technical founders and early-stage teams.
- **development**: Full-stack development tools
- **business-tools**: Communication and project management
- **analytics**: User behavior and business metrics

## Usage

Import any of these configurations using the `anvil config import` command:

```bash
# Import frontend developer setup
anvil config import import-examples/frontend-developer.yaml

# Import from remote URL
anvil config import https://raw.githubusercontent.com/rocajuanma/anvil/master/import-examples/backend-developer.yaml

# Import multiple configurations
anvil config import https://raw.githubusercontent.com/rocajuanma/anvil/master/import-examples/data-scientist.yaml
anvil config import https://raw.githubusercontent.com/rocajuanma/anvil/master/import-examples/devops-engineer.yaml
```

## Customization

These examples serve as starting points. You can:

1. **Modify existing groups**: Edit the YAML files to add/remove tools
2. **Create new groups**: Add additional groups for specific workflows
3. **Combine personas**: Import multiple configurations and merge groups
4. **Share with teams**: Use these as templates for team standardization

## Group Structure

Each configuration follows this structure:

```yaml
groups:
  group-name:
    - tool1
    - tool2
    - tool3
```

- **Group names**: Use kebab-case (e.g., `frontend-core`, `ml-ops`)
- **Tool names**: Use lowercase with hyphens (e.g., `visual-studio-code`, `aws-cli`)
- **Size**: Keep groups focused with 3-8 tools each

## Contributing

To add new persona configurations:

1. Create a new YAML file with the persona name
2. Define 3-5 relevant groups with 3-8 tools each
3. Update this README with the new persona description
4. Follow the established naming conventions

## Security Note

These example files contain only tool group definitions. No sensitive configuration data, API keys, or personal information is included. The import command extracts only the `groups` section and ignores all other data.
