# Build-related targets
# This file contains all build-related functionality

.PHONY: build build-single install clean

build:
	@echo "Building $(BINARY_NAME) $(VERSION)"
	@go build $(LDFLAGS) -o ./bin/$(BINARY_NAME) ./cmd/gikopsctl

install:
	@echo "Installing $(BINARY_NAME) $(VERSION)"
	@go install $(LDFLAGS) ./cmd/gikopsctl

clean:
	@echo "Cleaning build artifacts"
	@rm -f $(BINARY_NAME)
	@rm -rf ./bin/
	@rm -rf ./dist/
	@rm -f coverage.out coverage.html

# Build for specific OS/ARCH (used by release)
build-single:
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "Error: GOOS and GOARCH must be specified"; \
		exit 1; \
	fi
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)"
	@binary_name="$(BINARY_NAME)"; \
	if [ "$(GOOS)" = "windows" ]; then \
		binary_name="$(BINARY_NAME).exe"; \
	fi; \
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		$(LDFLAGS) \
		-o ./dist/$$binary_name \
		./cmd/gikopsctl