# Release-related targets
# This file contains all release functionality adapted from Kustomize release pipeline

# Release configuration
BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
RELEASE_DIR = ./dist
CHANGELOG_FILE = $(shell mktemp)

# Architecture lists
LINUX_ARCHS = amd64 arm64 s390x ppc64le
DARWIN_ARCHS = amd64 arm64  
WINDOWS_ARCHS = amd64 arm64
ALL_OSES = linux darwin windows

.PHONY: release build-release check-release-tag create-github-release compile-changelog
.PHONY: build-linux-release build-darwin-release build-windows-release package-releases generate-checksums

check-release-tag:
	@echo "Release version: $(VERSION)"

release: check-release-tag
	@echo "Creating complete release for $(VERSION)"
	@$(MAKE) build-release
	@$(MAKE) compile-changelog
	@$(MAKE) create-github-release

build-release: prepare-release-dir
	@echo "Building release binaries for $(VERSION)"
	@$(MAKE) package-releases
	@$(MAKE) generate-checksums

prepare-release-dir:
	@echo "Preparing release directory: $(RELEASE_DIR)"
	@mkdir -p $(RELEASE_DIR)
	@rm -f $(RELEASE_DIR)/*

package-releases:
	@echo "Packaging release binaries"
	@for os in $(ALL_OSES); do \
		if [ "$$os" = "linux" ]; then \
			archs="$(LINUX_ARCHS)"; \
		elif [ "$$os" = "darwin" ]; then \
			archs="$(DARWIN_ARCHS)"; \
		else \
			archs="$(WINDOWS_ARCHS)"; \
		fi; \
		for arch in $$archs; do \
			$(MAKE) build-single GOOS=$$os GOARCH=$$arch; \
			cd $(RELEASE_DIR); \
			if [ "$$os" = "windows" ]; then \
				binary_name="$(BINARY_NAME).exe"; \
				if [ -f "$$binary_name" ]; then \
					zip -j "$(BINARY_NAME)_$(VERSION)_$${os}_$${arch}.zip" "$$binary_name"; \
					rm "$$binary_name"; \
				fi; \
			else \
				binary_name="$(BINARY_NAME)"; \
				if [ -f "$$binary_name" ]; then \
					tar czf "$(BINARY_NAME)_$(VERSION)_$${os}_$${arch}.tar.gz" "$$binary_name"; \
					rm "$$binary_name"; \
				fi; \
			fi; \
			cd ..; \
		done; \
	done

generate-checksums:
	@echo "Generating checksums"
	@cd $(RELEASE_DIR) && \
	rm -f checksums.txt && \
	for file in *; do \
		if [ "$$file" != "checksums.txt" ] && [ -f "$$file" ]; then \
			echo "Generating checksum for: $$file"; \
			sha256sum "$$file" >> checksums.txt; \
		fi; \
	done

compile-changelog:
	@echo "Generating changelog for $(VERSION)"
	@./hack/makes/changelog.sh "$(PROJECT_NAME)" "$(VERSION)" "$(CHANGELOG_FILE)"	

create-github-release:
	@echo "Creating GitHub release for $(VERSION)"
	@if command -v gh >/dev/null 2>&1; then \
		gh release create "$(VERSION)" \
			--title "$(VERSION)" \
			--draft \
			--notes-file "$(CHANGELOG_FILE)" \
			$(RELEASE_DIR)/* || \
		echo "Failed to create GitHub release. Please ensure 'gh' is installed and authenticated."; \
	else \
		echo "GitHub CLI (gh) not found. Skipping GitHub release creation."; \
		echo "Release artifacts are available in: $(RELEASE_DIR)"; \
		echo "Changelog:"; \
		cat "$(CHANGELOG_FILE)"; \
	fi

# Clean up release artifacts
clean-release:
	@echo "Cleaning release artifacts"
	@rm -rf $(RELEASE_DIR)
	@rm -f "$(CHANGELOG_FILE)" 2>/dev/null || true