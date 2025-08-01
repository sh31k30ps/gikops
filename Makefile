# Main Makefile for gikopsctl

# Project configuration
BINARY_NAME = gikopsctl
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.0.0")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION = $(shell go version | cut -d ' ' -f 3)
FILE_PATH ?= ./...

# Build flags
LDFLAGS = -ldflags "-X github.com/sh31k30ps/gikopsctl/pkg/version.Version=${VERSION} \
                    -X github.com/sh31k30ps/gikopsctl/pkg/version.GitCommit=${GIT_COMMIT} \
                    -X github.com/sh31k30ps/gikopsctl/pkg/version.BuildTime=${BUILD_TIME} \
                    -X github.com/sh31k30ps/gikopsctl/pkg/version.GoVersion=${GO_VERSION}"

# Include sub-makefiles
include hack/makes/build.mk
include hack/makes/release.mk
include hack/makes/test.mk
include hack/makes/tools.mk

.PHONY: help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build         - Build the binary"
	@echo "  install       - Install the binary"
	@echo "  clean         - Clean build artifacts"
	@echo ""
	@echo "Release targets:"
	@echo "  release       - Create a complete release (requires TAG=<tag>)"
	@echo "  build-release - Build release binaries only"
	@echo ""
	@echo "Development targets:"
	@echo "  generate      - Generate code (deepcopy)"
	@echo "  test          - Run tests"
	@echo "  coverage      - Generate test coverage"
	@echo "  lint          - Run linting"
	@echo ""
	@echo "Tool targets:"
	@echo "  install-tools - Install development tools"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make release TAG=v1.0.0"
	@echo "  make test FILE_PATH=./pkg/..."

.DEFAULT_GOAL := help