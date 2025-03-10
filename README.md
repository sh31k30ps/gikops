# Git Ops Helper

A Go-based tool for managing Kubernetes components and environments, replacing traditional Makefile and bash script approaches with a more robust solution.

## Prerequisites

The following tools must be installed and available in your PATH:

- Docker (>= 20.10.0)
- kubectl (>= 1.20.0)
- Helm (>= 3.0.0)
- Kind (>= 0.20.0)
- Kustomize (>= 5.0.0)

The tool will verify these dependencies before running any command.

## Features

- Kubernetes environment management with Kind
- Component initialization and application
- Helm chart management
- CRD handling
- Post-initialization processing (uploads, resolves, renames)
- Automatic tool verification
- Version information with build details
- Shell completion for Bash, Zsh, Fish, and PowerShell

## Project Structure

```
.
├── cmd/
│   └── gikopsctl/      # Main application entry point
├── pkg/
│   ├── cmd/                 # Command implementations
│   ├── component/           # Component management logic
│   ├── installer/           # Environment installation logic
│   ├── tools/              # Tool verification and utilities
│   ├── version/            # Version information
│   └── types/              # Common type definitions
├── components/             # Component definitions
└── overrides/             # Kubernetes overrides
```

## Installation

```bash
go install github.com/sh31k30ps/gikopsctl/cmd/gikopsctl@latest
```

## Usage

### Shell Completion

Generate shell completion scripts for:

```bash
# Bash
gikopsctl completion bash > /etc/bash_completion.d/gikopsctl

# Zsh
gikopsctl completion zsh > "${fpath[1]}/_gikopsctl"

# Fish
gikopsctl completion fish > ~/.config/fish/completions/gikopsctl.fish

# PowerShell
gikopsctl completion powershell > gikopsctl.ps1
```

For detailed installation instructions for each shell, run:
```bash
gikopsctl completion --help
```

### Version Information

```bash
# Display version and build information
gikopsctl version
```

### Environment Management

```bash
# Install environment
gikopsctl install

# Uninstall environment
gikopsctl uninstall
```

### Component Management

```bash
# Initialize all components
gikopsctl component init

# Initialize specific component
gikopsctl component init [component-name]

# Apply all components
gikopsctl component apply --env local --mode all

# Apply specific component
gikopsctl component apply [component-name] --env local --mode all
```

### Modes

- `all`: Apply both CRDs and manifests
- `crd`: Apply only CRDs
- `manifests`: Apply only manifests

## Component Configuration

Components are defined in `component.yaml` files within the `components/` directory. Example:

```yaml
kind: Component
apiVersion: bbr.k8s.io/v1alpha1
metadata:
  name: example
  namespace: example-ns
helm:
  repo: example-repo
  repo-url: https://example.com/charts
  version: v1.0.0
  chart: example/chart
  post-init:
    uploads:
      - name: crds.yaml
        url: https://example.com/crds.yaml
```

## Development

To build from source:

```bash
# Build with version information
make build

# Install to GOPATH/bin
make install

# Run tests
make test

# Run linting
make lint
