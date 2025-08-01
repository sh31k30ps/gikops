# Test-related targets
# This file contains all testing and quality assurance functionality

.PHONY: test coverage lint fmt vet test-verbose test-short test-race

test:
	@echo "Running tests for $(FILE_PATH)"
	@go test -v $(FILE_PATH)

test-short:
	@echo "Running short tests for $(FILE_PATH)"
	@go test -short $(FILE_PATH)

test-verbose:
	@echo "Running verbose tests for $(FILE_PATH)"
	@go test -v -race $(FILE_PATH)

test-race:
	@echo "Running race tests for $(FILE_PATH)"
	@go test -race $(FILE_PATH)

coverage:
	@echo "Generating test coverage for $(FILE_PATH)"
	@go test -coverprofile=coverage.out $(FILE_PATH)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out
	@rm coverage.out

coverage-ci:
	@echo "Generating test coverage for CI"
	@go test -coverprofile=coverage.out $(FILE_PATH)
	@go tool cover -func=coverage.out

lint:
	@echo "Running linting checks"
	@$(MAKE) vet
	@$(MAKE) fmt-check

vet:
	@echo "Running go vet for $(FILE_PATH)"
	@go vet $(FILE_PATH)

fmt:
	@echo "Formatting Go code"
	@go fmt $(FILE_PATH)

fmt-check:
	@echo "Checking Go code formatting"
	@test -z "$$(gofmt -l . | grep -v vendor/)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

# Additional linting tools (if available)
lint-advanced:
	@echo "Running advanced linting (requires golangci-lint)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run $(FILE_PATH); \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		$(MAKE) lint; \
	fi

# Benchmark tests
bench:
	@echo "Running benchmarks for $(FILE_PATH)"
	@go test -bench=. -benchmem $(FILE_PATH)

# Clean test artifacts
clean-test:
	@echo "Cleaning test artifacts"
	@rm -f coverage.out coverage.html
	@rm -f *.test
	@rm -f test.log