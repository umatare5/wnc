# Makefile for wnc CLI application

.PHONY: help clean deps lint test-unit test-integration test-mock test-coverage test-coverage-html generate-mocks build build-snapshot run install

# Default target
help:
	@echo "Available targets:"
	@echo "  help             - Show this help message"
	@echo "  clean            - Clean build artifacts and temporary files"
	@echo "  deps             - Install development dependencies"
	@echo "  lint             - Run linting tools"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests (requires environment variables)"
	@echo "  test-mock        - Run GoMock-based tests"
	@echo "  test-coverage    - Run tests with coverage analysis"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo "  generate-mocks   - Generate mock implementations using GoMock"
	@echo "  build            - Build the CLI application"
	@echo "  build-snapshot   - Build snapshot release with goreleaser"
	@echo "  run              - Run the CLI application"
	@echo "  install          - Install the CLI application to GOPATH/bin"

# Clean build artifacts and temporary files
clean:
	@echo "Cleaning build artifacts..."
	rm -f coverage.out
	rm -rf ./tmp
	go clean -cache -testcache
	@echo "Cleanup completed!"

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo "Installing goreleaser..."; \
		go install github.com/goreleaser/goreleaser@latest; \
	fi
	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	fi
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing air for hot reload..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@if ! command -v mockgen >/dev/null 2>&1; then \
		echo "Installing GoMock..."; \
		go install github.com/golang/mock/mockgen@latest; \
	fi
	@echo "Development dependencies installed!"

# Run linting (if tools are available)
lint:
	@echo "Running linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, running go vet instead..."; \
		go vet ./...; \
	fi
	@echo "Linting completed!"

# Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@mkdir -p ./tmp
	@if command -v gotestsum >/dev/null 2>&1; then \
		WNC_CONTROLLERS="" gotestsum --format testname -- -race ./...; \
	else \
		echo "gotestsum not found, running go test with verbose output..."; \
		WNC_CONTROLLERS="" go test -v -race ./...; \
	fi

# Run integration tests
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@if [ -z "$$WNC_CONTROLLERS" ]; then \
		echo "Warning: WNC_CONTROLLERS not set - integration tests will be skipped"; \
		echo "Set WNC_CONTROLLERS to run integration tests"; \
	fi
	@mkdir -p ./tmp
	@if command -v gotestsum >/dev/null 2>&1; then \
		gotestsum --format testname -- -race -run "TestIntegration" ./internal/cli/...; \
	else \
		echo "gotestsum not found, running go test with verbose output..."; \
		go test -v -race -run "TestIntegration" ./internal/cli/...; \
	fi

# Run GoMock-based tests
.PHONY: test-mock
test-mock:
	@echo "Running GoMock-based tests..."
	@mkdir -p ./tmp
	@if command -v gotestsum >/dev/null 2>&1; then \
		gotestsum --format testname -- -race ./test/mock/...; \
	else \
		echo "gotestsum not found, running go test with verbose output..."; \
		go test -v -race ./test/mock/...; \
	fi

# Generate mock implementations using GoMock
.PHONY: generate-mocks
generate-mocks:
	@echo "Generating mock implementations..."
	@if ! command -v mockgen >/dev/null 2>&1; then \
		echo "mockgen not found. Install it with 'make deps' first"; \
		exit 1; \
	fi
	@cd pkg/cisco && go generate
	@echo "Mock generation completed!"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p ./tmp
	@if command -v gotestsum >/dev/null 2>&1; then \
		WNC_CONTROLLERS="" gotestsum --format testname -- -race -coverprofile=./tmp/coverage.out ./...; \
	else \
		echo "gotestsum not found, running go test with verbose output..."; \
		WNC_CONTROLLERS="" go test -v -race -coverprofile=./tmp/coverage.out ./...; \
	fi
	@if [ -f ./tmp/coverage.out ]; then \
		echo "Coverage report generated at ./tmp/coverage.out"; \
		go tool cover -func=./tmp/coverage.out | tail -1; \
	fi

# Generate HTML coverage report
.PHONY: test-coverage-html
test-coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	@mkdir -p ./tmp
	@if [ -f ./tmp/coverage.out ]; then \
		go tool cover -html=./tmp/coverage.out -o ./tmp/coverage.html; \
		echo "HTML coverage report generated at ./tmp/coverage.html"; \
		if command -v open >/dev/null 2>&1; then \
			echo "Opening coverage report in browser..."; \
			open ./tmp/coverage.html; \
		fi; \
	else \
		echo "No coverage file found. Run 'make test-coverage' first."; \
	fi

# Build the CLI application
build:
	@echo "Building wnc CLI application..."
	@mkdir -p ./tmp
	go build -o ./tmp/wnc cmd/main.go
	@echo "Build completed! Binary: ./tmp/wnc"

# Build snapshot release with goreleaser
build-snapshot:
	@echo "Building snapshot release..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not found. Install it with 'make deps' or run 'make build' instead"; \
		exit 1; \
	fi

# Run the CLI application
run: build
	@echo "Running wnc CLI application..."
	./tmp/wnc

# Install the CLI application to GOPATH/bin
install:
	@echo "Installing wnc CLI application..."
	go install cmd/main.go
	@echo "Installation completed! wnc is now available in your PATH"

# Development target with hot reload (requires air)
.PHONY: dev
dev:
	@if command -v air >/dev/null 2>&1; then \
		echo "Starting development server with hot reload..."; \
		air; \
	else \
		echo "air not found. Install it with 'make deps' first"; \
		exit 1; \
	fi
