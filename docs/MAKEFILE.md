# Makefile Structure Documentation

This document describes the modular Makefile structure for the gikopsctl project.

## Directory Structure

```
project-root/
├── Makefile                # Main Makefile (entry point)
├── hack/makes/             # Sub-Makefiles directory
│   ├── build.mk            # Build-related targets
│   ├── changelog.sh        # Changelog generation script
│   ├── release.mk          # Release and distribution targets
│   ├── test.mk             # Testing and quality assurance targets
│   └── tools.mk            # Development tools and code generation
├── hack/tools/             # Go modules for development tools
├── bin/                    # Built binaries (created during build)
└── dist/                   # Release artifacts (created during release)
```

## File Responsibilities

### Main Makefile
- Project configuration and global variables
- Includes all sub-makefiles
- Provides help target with all available commands
- Sets default goal to help

### makefiles/build.mk
- **Primary responsibility**: Building binaries
- **Key targets**:
  - `build`: Build the main binary
  - `install`: Install the binary to GOPATH
  - `clean`: Clean build artifacts
  - `build-single`: Build for specific OS/ARCH (used by release)

### makefiles/release.mk
- **Primary responsibility**: Creating releases
- **Key targets**:
  - `release`: Complete release process (requires TAG=version)
  - `build-release`: Build binaries for all platforms
  - `create-github-release`: Create GitHub release with artifacts
  - `generate-checksums`: Generate SHA256 checksums
- **Platforms supported**: Linux, Darwin, Windows
- **Architectures**: amd64, arm64 (plus s390x, ppc64le for Linux in original)

### makefiles/test.mk
- **Primary responsibility**: Testing and code quality
- **Key targets**:
  - `test`: Run standard tests
  - `coverage`: Generate test coverage reports
  - `lint`: Run linting checks (vet + formatting)
  - `fmt`: Format Go code
  - `bench`: Run benchmark tests

### makefiles/tools.mk
- **Primary responsibility**: Development tools and code generation
- **Key targets**:
  - `generate`: Run all code generation
  - `install-tools`: Install all development tools
  - `deepcopy-gen`: Generate Kubernetes deepcopy code
  - `check-tools`: Verify tool installation

## Usage Examples

### Basic Development Workflow
```bash
# Build the project
make build

# Run tests
make test

# Generate code
make generate

# Clean artifacts
make clean
```

### Release Workflow
```bash
# Create a complete release
make release 

# Just build release binaries (no GitHub release)
make build-release

# Clean release artifacts
make clean-release
```

### Testing Workflow
```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run linting
make lint

# Run specific test path
make test FILE_PATH=./pkg/...
```

### Tool Management
```bash
# Install all development tools
make install-tools

# Check tool versions
make check-tools

# Update tools
make update-tools
```

## Configuration Variables

### Global Variables (Main Makefile)
- `BINARY_NAME`: Name of the binary (gikopsctl)
- `VERSION`: Git-based version string
- `GIT_COMMIT`: Short git commit hash
- `BUILD_TIME`: Build timestamp
- `GO_VERSION`: Go compiler version
- `FILE_PATH`: Test file path pattern (default: ./...)

### Release Variables (release.mk)
- `RELEASE_DIR`: Directory for release artifacts (./dist)

## Integration with Original Script

The release functionality has been adapted from the original Kustomize release script with these changes:

1. **Modular structure**: Split into logical sub-makefiles
2. **Project adaptation**: Adapted for gikopsctl instead of kustomize
3. **Simplified changelog**: Uses git log if no custom script exists
4. **Error handling**: Graceful handling of missing tools (gh, jq)
5. **Cross-platform**: Supports the same OS/arch combinations

## Dependencies

### Required for basic functionality:
- Go compiler
- Git

### Required for releases:
- GitHub CLI (`gh`) - for creating GitHub releases
- `jq` - for JSON processing (if using custom changelog script)

### Optional for enhanced development:
- `golangci-lint` - for advanced linting
- Custom changelog script at `./scripts/compile-changelog.sh`

## Extending the Makefile

To add new functionality:

1. **Choose appropriate sub-makefile** based on responsibility
2. **Add new targets** following the existing patterns
3. **Use sub-make calls** with `$(MAKE)` for modularity
4. **Update help target** in main Makefile if needed
5. **Document new targets** in this file
