# Expense Tracker Makefile
# Build configuration
APP_NAME := expense-tracker
BIN_DIR := bin
SRC_DIR := cmd/app
MAIN_FILE := $(SRC_DIR)/main.go

# Go configuration
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOFMT := gofmt
GOLINT := golangci-lint

# Build flags
LDFLAGS := -ldflags "-s -w"
BUILD_FLAGS := -v $(LDFLAGS)

# OS and architecture detection for cross-compilation
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

.PHONY: all build clean test fmt lint vet deps install-tools run dev help

# Default target
all: clean fmt lint test build

# Build the application
build: deps
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"

# Build for multiple platforms
build-all: clean deps
	@echo "Building for multiple platforms..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "Cross-compilation complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	$(GOCLEAN)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

# Lint code
lint: install-lint
	@echo "Running linter..."
	$(GOLINT) run

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOCMD) mod download
	$(GOCMD) mod tidy

# Install development tools
install-tools: install-lint
	@echo "Development tools installed"

# Install golangci-lint
install-lint:
	@which golangci-lint > /dev/null || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	./$(BIN_DIR)/$(APP_NAME)

# Development mode (run without building binary)
dev:
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_FILE)

# Install the binary to $GOPATH/bin or $GOBIN
install:
	@echo "Installing $(APP_NAME)..."
	$(GOCMD) install $(BUILD_FLAGS) $(MAIN_FILE)

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GOCMD) get -u ./...
	$(GOCMD) mod tidy

# Security audit
audit:
	@echo "Running security audit..."
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# Full quality check
quality: fmt vet lint test
	@echo "Quality checks passed!"

# Help target
help:
	@echo "Available targets:"
	@echo "  build         - Build the application binary"
	@echo "  build-all     - Build for multiple platforms (Linux, macOS, Windows)"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format code and tidy modules"
	@echo "  lint          - Run linter (installs golangci-lint if needed)"
	@echo "  vet           - Run go vet"
	@echo "  deps          - Install and tidy dependencies"
	@echo "  install-tools - Install development tools"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run in development mode (no binary)"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  update-deps   - Update all dependencies"
	@echo "  audit         - Run security vulnerability audit"
	@echo "  quality       - Run all quality checks (fmt, vet, lint, test)"
	@echo "  help          - Show this help message"