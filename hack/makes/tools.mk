# Tools and code generation targets
# This file contains all development tools and code generation functionality

TOOLS_DIR = ./bin
HACK_TOOLS_DIR = ./hack/tools

.PHONY: generate install-tools clean-tools deepcopy-gen

generate: deepcopy-gen
	@echo "Running code generation"

deepcopy-gen: install-deepcopy-gen
	@echo "Generating deepcopy code"
	@$(TOOLS_DIR)/deepcopy-gen --output-file zz_generated.deepcopy.go ./api/config/v1alpha1

install-deepcopy-gen:
	@echo "Installing deepcopy-gen tool"
	@if [ ! -f $(TOOLS_DIR)/deepcopy-gen ]; then \
		mkdir -p $(TOOLS_DIR); \
		if [ -d "$(HACK_TOOLS_DIR)" ]; then \
			cd $(HACK_TOOLS_DIR) && \
			go mod tidy && \
			go mod vendor && \
			go build -o ../../$(TOOLS_DIR)/deepcopy-gen k8s.io/code-generator/cmd/deepcopy-gen; \
		else \
			echo "Creating hack/tools directory and go.mod"; \
			mkdir -p $(HACK_TOOLS_DIR); \
			cd $(HACK_TOOLS_DIR) && \
			go mod init tools && \
			go get k8s.io/code-generator/cmd/deepcopy-gen && \
			go mod tidy && \
			go build -o ../../$(TOOLS_DIR)/deepcopy-gen k8s.io/code-generator/cmd/deepcopy-gen; \
		fi; \
	fi

install-tools: install-development-tools install-build-tools

install-development-tools:
	@echo "Installing development tools"
	@$(MAKE) install-deepcopy-gen
	@$(MAKE) install-linter

install-build-tools:
	@echo "Installing build tools"
	@if ! command -v gh >/dev/null 2>&1; then \
		echo "GitHub CLI (gh) not found. Please install it:"; \
		echo "  macOS: brew install gh"; \
		echo "  Linux: See https://cli.github.com/manual/installation"; \
	fi
	@if ! command -v jq >/dev/null 2>&1; then \
		echo "jq not found. Please install it:"; \
		echo "  macOS: brew install jq"; \
		echo "  Linux: apt-get install jq or yum install jq"; \
	fi

install-linter:
	@echo "Installing golangci-lint"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	else \
		echo "golangci-lint already installed"; \
	fi

# Tool version checks
check-tools:
	@echo "Checking tool versions:"
	@echo -n "Go: "; go version 2>/dev/null || echo "not found"
	@echo -n "Git: "; git --version 2>/dev/null || echo "not found"
	@echo -n "GitHub CLI: "; gh --version 2>/dev/null || echo "not found"
	@echo -n "jq: "; jq --version 2>/dev/null || echo "not found"
	@echo -n "golangci-lint: "; golangci-lint --version 2>/dev/null || echo "not found"
	@if [ -f "$(TOOLS_DIR)/deepcopy-gen" ]; then \
		echo "deepcopy-gen: installed"; \
	else \
		echo "deepcopy-gen: not found"; \
	fi

# Clean tools
clean-tools:
	@echo "Cleaning tools"
	@rm -rf $(TOOLS_DIR)
	@rm -rf $(HACK_TOOLS_DIR)/vendor

# Update tools
update-tools:
	@echo "Updating tools"
	@$(MAKE) clean-tools
	@$(MAKE) install-tools