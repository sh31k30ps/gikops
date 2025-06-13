BINARY_NAME=gikopsctl
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.1")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(shell go version | cut -d ' ' -f 3)

LDFLAGS=-ldflags "-X github.com/sh31k30ps/gikopsctl/pkg/version.Version=${VERSION} \
                  -X github.com/sh31k30ps/gikopsctl/pkg/version.GitCommit=${GIT_COMMIT} \
                  -X github.com/sh31k30ps/gikopsctl/pkg/version.BuildTime=${BUILD_TIME} \
                  -X github.com/sh31k30ps/gikopsctl/pkg/version.GoVersion=${GO_VERSION}"

.PHONY: generate
generate:
	@if [ ! -f bin/deepcopy-gen ]; then \
		cd hack/tools && \
		go mod tidy && \
		go mod vendor && \
		go build -o ../../bin/deepcopy-gen k8s.io/code-generator/cmd/deepcopy-gen; \
	fi
	@bin/deepcopy-gen --output-file zz_generated.deepcopy.go ./api/config/v1alpha1

.PHONY: build
build:
	go build ${LDFLAGS} -o ./bin/${BINARY_NAME} ./cmd/gikopsctl

.PHONY: install
install:
	@go install ${LDFLAGS} ./cmd/gikopsctl

.PHONY: clean
clean:
	@rm -f ${BINARY_NAME}

.PHONY: test
test:
	@go test -v ./...

.PHONY: coverage
coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out


.PHONY: lint
lint:
	@go vet ./...
	@test -z $(gofmt -l .)
