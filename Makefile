.PHONY: help build lint test-unit test-unit-coverage clean image snapshot pre-commit-install pre-commit-test pre-commit-uninstall

BINARY_NAME := wnc
BUILD_DIR := ./tmp
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)
COVERAGE_DIR := ./coverage
IMAGE_DIR := $(BUILD_DIR)/image
GOARCH := $(shell go env GOARCH)

LDFLAGS := -X github.com/umatare5/wnc/internal/cli.version=$(shell cat VERSION)
BUILD_FLAGS := -trimpath -ldflags "$(LDFLAGS)"

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  build                - Build the binary into $(BINARY_PATH)"
	@echo "  lint                 - Run golangci-lint and go mod tidy"
	@echo "  test-unit            - Run unit tests with coverage"
	@echo "  test-unit-coverage   - Generate the HTML coverage report"
	@echo "  snapshot             - Build a goreleaser snapshot"
	@echo "  image                - Build the Docker image"
	@echo "  clean                - Remove build artifacts"
	@echo "  pre-commit-install   - Install the pre-commit hooks"
	@echo "  pre-commit-test      - Run every hook across the whole tree"
	@echo "  pre-commit-uninstall - Remove the pre-commit hooks"
	@echo ""
	@echo "Requirements:"
	@echo "  - gotestsum:     go install gotest.tools/gotestsum@latest"
	@echo "  - golangci-lint: https://golangci-lint.run/docs/welcome/install/"
	@echo "  - pre-commit:    https://pre-commit.com/#install"
	@echo "  - gitleaks:      https://github.com/gitleaks/gitleaks#installing"

build: $(BINARY_PATH)

$(BINARY_PATH):
	mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd

# config verify comes first because `run` accepts an unknown nested setting key
# silently and reverts that setting to its default, so a typo leaves a rule the
# author believes is on quietly off.
lint:
	golangci-lint config verify
	golangci-lint run
	go mod tidy

# The WNC_* variables are cleared because they are read at flag-parse time, so a developer's own
# shell must not decide what the CLI tests see. A new variable the CLI reads has to be added here
# as well as to the neutralization inside cli_test.go.
test-unit:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Error: gotestsum is not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	mkdir -p $(COVERAGE_DIR)
	env -u WNC_CONTROLLER -u WNC_ACCESS_TOKEN -u WNC_CONFIG -u WNC_USERNAME -u WNC_PASSWORD \
		gotestsum --format testname -- -race -coverprofile=$(COVERAGE_DIR)/report.out ./...

test-unit-coverage: test-unit
	go tool cover -html=$(COVERAGE_DIR)/report.out -o $(COVERAGE_DIR)/report.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/report.html"

snapshot:
	goreleaser release --snapshot --clean

# BUILD_DIR is ./tmp, which is also the worktree root: this deletes any worktree there.
clean:
	rm -rf $(BUILD_DIR) $(COVERAGE_DIR)
	find . -name "*.bak*" -type f -delete 2>/dev/null || true

# Assembles the same context goreleaser hands docker — see the Dockerfile header. The
# linux/$(GOARCH) subdirectory is what makes the shared COPY line resolve in both.
image:
	mkdir -p $(IMAGE_DIR)/linux/$(GOARCH)
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_FLAGS) -o $(IMAGE_DIR)/linux/$(GOARCH)/$(BINARY_NAME) ./cmd
	cp LICENSE $(IMAGE_DIR)/
	docker build --platform linux/$(GOARCH) -f Dockerfile -t $(USER)/$(BINARY_NAME) $(IMAGE_DIR)

# --allow-missing-config is load-bearing: the hook path is the shared git common
# dir, so the hook installed from one worktree also runs on every other one and
# on main, where .pre-commit-config.yaml may not exist yet.
pre-commit-install:
	@command -v pre-commit >/dev/null 2>&1 || { echo "Error: pre-commit is not installed. See: https://pre-commit.com/#install"; exit 1; }
	@pre-commit install --allow-missing-config

pre-commit-test:
	@pre-commit run --all-files

pre-commit-uninstall:
	@pre-commit uninstall
