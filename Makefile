BINARY_NAME=terraform-provider-census
VERSION=0.2.9
BUILD_DIR=bin
LDFLAGS=-ldflags "-X main.version=${VERSION}"

.PHONY: build test clean install fmt vet lint help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the provider binary
	@echo "Building ${BINARY_NAME} v${VERSION}..."
	@mkdir -p ${BUILD_DIR}
	@go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} .

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./... -short

test-unit: ## Run unit tests (no credentials needed)
	@echo "Running unit tests (no credentials required)..."
	@go test -v ./census/tests/provider/unit ./census/tests/client -short

test-integration: ## Run integration tests (creates all resources in staging)
	@echo "Running integration tests against Census staging API..."
	@echo "This will create workspaces, sources, destinations, and syncs in staging"
	@echo "Requires: .env.test with staging credentials (see .env.test.example)"
	@if [ ! -f .env.test ]; then \
		echo "Error: .env.test not found. Copy .env.test.example and fill in your credentials."; \
		exit 1; \
	fi
	@set -a && . ./.env.test && set +a && TF_ACC=1 go test -v ./census/tests/provider/acceptance -timeout 60m

test-acc: test-integration ## Alias for test-integration (Terraform convention)

test-coverage: ## Generate test coverage report
	@echo "Generating test coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total:

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf ${BUILD_DIR}
	@go clean

install: build ## Install the provider locally
	@echo "Installing ${BINARY_NAME} locally..."
	@go install ${LDFLAGS} .

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@golangci-lint run

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

docs: ## Generate documentation
	@echo "Generating documentation..."
	@terraform fmt -recursive examples/

build-all: clean ## Build binaries for all platforms
	@echo "Building ${BINARY_NAME} v${VERSION} for all platforms..."
	@mkdir -p ${BUILD_DIR}
	@echo "Building for macOS (Apple Silicon)..."
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}_darwin_arm64 .
	@echo "Building for macOS (Intel)..."
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}_darwin_amd64 .
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}_linux_amd64 .
	@echo "Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}_linux_arm64 .
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}_windows_amd64.exe .
	@echo "Build complete! Binaries are in ${BUILD_DIR}/"
	@ls -lh ${BUILD_DIR}/

release: clean test build-all ## Build and test release binaries for all platforms
	@echo "Release binaries built successfully!"
	@echo "Upload these files to GitHub release:"
	@ls -1 ${BUILD_DIR}/

dev: build ## Build and set up for local development
	@echo "Setting up for local development..."
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/sutrolabs/census/${VERSION}/$(shell go env GOOS)_$(shell go env GOARCH)/
	@cp ${BUILD_DIR}/${BINARY_NAME} ~/.terraform.d/plugins/registry.terraform.io/sutrolabs/census/${VERSION}/$(shell go env GOOS)_$(shell go env GOARCH)/${BINARY_NAME}_v${VERSION}
	@echo "Provider installed to local Terraform plugins directory"

install-local: build ## Install provider to local plugin directory (filesystem mirror)
	@echo "Installing provider to ~/.terraform.d/plugins/..."
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/sutrolabs/census/${VERSION}/$(shell go env GOOS)_$(shell go env GOARCH)/
	@cp ${BUILD_DIR}/${BINARY_NAME} ~/.terraform.d/plugins/registry.terraform.io/sutrolabs/census/${VERSION}/$(shell go env GOOS)_$(shell go env GOARCH)/${BINARY_NAME}_v${VERSION}
	@echo "Installed to: ~/.terraform.d/plugins/registry.terraform.io/sutrolabs/census/${VERSION}/$(shell go env GOOS)_$(shell go env GOARCH)/${BINARY_NAME}_v${VERSION}"

.DEFAULT_GOAL := help