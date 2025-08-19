.PHONY: test test-unit test-integration test-coverage test-coverage-html mocks clean build run lint fmt vet check-all

# Run all tests
test:
	go run github.com/vektra/mockery/v2@latest --all
	go test -v ./...

# Run only unit tests (domain and use case layers)
test-unit:
	go test -v ./internal/core/domain/...
	go test -v ./internal/core/usecase/...

# Run only integration tests
test-integration:
	go test -v ./test/integration/...

# Run tests with coverage report
test-coverage:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Run tests with coverage and generate HTML report
test-coverage-html:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	go test -race -v ./...

# Run tests with short flag for CI/CD
test-short:
	go test -short -v ./...

# Generate mocks
mocks:
	go run github.com/vektra/mockery/v2@latest --all

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		@echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Vet code
vet:
	go vet ./...

# Run all code quality checks
check-all: fmt vet test-coverage
	@echo "All quality checks completed"

# Build the application
build:
	go build -o bin/expense-tracker cmd/app/main.go

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/expense-tracker-linux-amd64 cmd/app/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/expense-tracker-windows-amd64.exe cmd/app/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/expense-tracker-darwin-amd64 cmd/app/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/expense-tracker-darwin-arm64 cmd/app/main.go

# Run the application
run:
	go run cmd/app/main.go

# Run the application in development mode with verbose logging
run-dev:
	DEBUG=true go run cmd/app/main.go

# Install development dependencies
dev-deps:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Clean generated files
clean:
	rm -f expense-tracker coverage.out coverage.html
	rm -rf test/mocks/*
	rm -rf bin/

# Show help
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests"
	@echo "  test-unit         - Run only unit tests"
	@echo "  test-integration  - Run only integration tests"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-coverage-html- Generate HTML coverage report"
	@echo "  test-race         - Run tests with race detection"
	@echo "  test-short        - Run tests with short flag"
	@echo "  mocks             - Generate mocks"
	@echo "  fmt               - Format code"
	@echo "  lint              - Lint code"
	@echo "  vet               - Vet code"
	@echo "  check-all         - Run all quality checks"
	@echo "  build             - Build application"
	@echo "  build-all         - Build for multiple platforms"
	@echo "  run               - Run application"
	@echo "  run-dev           - Run in development mode"
	@echo "  dev-deps          - Install development dependencies"
	@echo "  clean             - Clean generated files"
	@echo "  help              - Show this help"