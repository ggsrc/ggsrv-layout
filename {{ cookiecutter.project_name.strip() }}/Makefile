# =============================================================================
# Configuration
# =============================================================================

resolve-bin = $(shell command -v $(1) 2>/dev/null)
GO := $(call resolve-bin,go)
DOCKER := $(call resolve-bin,docker)
ifeq ($(strip $(GO)),)
GO := go
endif
ifeq ($(strip $(DOCKER)),)
DOCKER := docker
endif

TEST_DIRS := $(shell $(GO) list ./...)
CGO_ENABLED=0

# Determine OS and go build args
UNAME := $(shell uname -m)
ifeq ($(UNAME), arm64)
  # support CGO on apple silicon
  GOOS=darwin
else
  GOOS=linux
endif

# =============================================================================
# Default Target
# =============================================================================

.DEFAULT_GOAL := help

# =============================================================================
# Build
# =============================================================================

.PHONY: build local-run

build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=$(GOOS) $(GO) build -o bin/ ./...

local-run:
	$(GO) run cmd/server/main.go

# =============================================================================
# Proto (gRPC/Protobuf)
# =============================================================================

.PHONY: proto proto-lint proto-format proto-breaking proto-validate proto-clean proto-deps proto-update

proto:
	cd api && make proto

proto-lint:
	cd api && make lint

proto-format:
	cd api && make format

proto-breaking:
	cd api && make breaking

proto-validate:
	cd api && make validate

proto-clean:
	cd api && make clean

proto-deps:
	cd api && buf mod deps

proto-update:
	cd api && buf mod update

# =============================================================================
# Code Generation (SQL, Converters)
# =============================================================================

.PHONY: repos convert

repos:
	cd internal/repos && sqlc generate

convert:
	goverter gen ./internal/convert/...

# =============================================================================
# Formatting
# =============================================================================

.PHONY: fmt fmt-check fmt-sql format format-check

fmt:
	@echo "--> Formatting Go code"
	@$(GO) fmt ./...
	@$(GO) mod tidy

fmt-check:
	@echo "--> Checking Go code formatting"
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix"; \
		exit 1; \
	fi

fmt-sql:
	@echo "--> Formatting SQL files"
	@if ! command -v sql-formatter >/dev/null 2>&1; then \
		echo "sql-formatter not found. Run scripts/install-dev-tools.sh."; \
		exit 1; \
	fi
	@echo "Note: Some SQL files may contain sqlc-specific syntax that cannot be formatted"
	@find internal/repos -type f -name "*.sql" -print0 | xargs -0 -I{} sh -c 'echo "Formatting {}" && sql-formatter --fix -c .sql-formatter.json "{}" 2>/dev/null || echo "  Skipped {} (contains unsupported syntax)"'

format: fmt fmt-sql
	@echo "--> All formatting complete"

format-check: fmt-check lint-sql
	@echo "--> All format checks complete"

# =============================================================================
# Linting
# =============================================================================

.PHONY: lint lint-fix lint-sql lint-all

lint: proto
	@echo "--> Running linter"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found. Run scripts/install-dev-tools.sh."; \
		exit 1; \
	fi
	@golangci-lint run

lint-fix:
	@echo "--> Running linter auto fix"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found. Run scripts/install-dev-tools.sh."; \
		exit 1; \
	fi
	@golangci-lint run --fix

lint-sql:
	@echo "--> Linting SQL files"
	@if ! command -v sql-formatter >/dev/null 2>&1; then \
		echo "sql-formatter not found. Run scripts/install-dev-tools.sh."; \
		exit 1; \
	fi
	@echo "Note: Some SQL files may contain sqlc-specific syntax that cannot be linted"
	@find internal/repos -type f -name "*.sql" -print0 | xargs -0 -I{} sh -c 'echo "Checking {}" && sql-formatter -c .sql-formatter.json "{}" 2>/dev/null || echo "  Skipped {} (contains unsupported syntax)"'

lint-all: lint lint-sql
	@echo "--> All linting complete"

# =============================================================================
# Testing
# =============================================================================

.PHONY: test test-coverage test-verbose test-update-golden

test:
	@echo "--> Running tests (Test-Container mode)"
	export ENV=test && \
	export TZ=UTC && \
	CGO_ENABLED=${CGO_ENABLED} GOOS=$(GOOS) $(GO) run gotest.tools/gotestsum@latest --format github-actions -- -tags=integration ./...

test-coverage:
	@echo "--> Running tests with coverage (Test-Container mode)"
	export ENV=test && \
	export TZ=UTC && \
	CGO_ENABLED=${CGO_ENABLED} GOOS=$(GOOS) $(GO) run gotest.tools/gotestsum@latest --format github-actions -- -tags=integration ./... -coverprofile=coverage.txt -covermode=atomic
	@echo "--> Generating HTML coverage report"
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "--> Coverage report: coverage.html"

test-verbose:
	@echo "--> Running tests verbose (Test-Container mode)"
	export ENV=test && \
	export TZ=UTC && \
	CGO_ENABLED=${CGO_ENABLED} GOOS=$(GOOS) $(GO) test ${TEST_DIRS} -test.v -tags=integration

test-update-golden:
	export ENV=test && \
	export TZ=UTC && \
	CGO_ENABLED=${CGO_ENABLED} GOOS=$(GOOS) $(GO) test ${TEST_DIRS} -test.v -p 1 -update

# =============================================================================
# Docker
# =============================================================================

.PHONY: docker-up docker-down

docker-up:
	$(DOCKER) compose up -d

docker-down:
	$(DOCKER) compose down

# =============================================================================
# Maintenance
# =============================================================================

.PHONY: update-template submodule

update-template:
	cruft update --skip-apply-ask --refresh-private-variables
	git restore --staged .

submodule:
	git submodule update --init --recursive

# =============================================================================
# Help
# =============================================================================

.PHONY: help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build:"
	@echo "  build              Build all binaries"
	@echo "  local-run          Run the server locally"
	@echo ""
	@echo "Proto:"
	@echo "  proto              Generate protobuf code"
	@echo "  proto-lint         Lint proto files"
	@echo "  proto-format       Format proto files"
	@echo "  proto-breaking     Check for breaking changes"
	@echo "  proto-validate     Validate proto files"
	@echo "  proto-clean        Clean generated proto files"
	@echo "  proto-deps         Install proto dependencies"
	@echo "  proto-update       Update proto dependencies"
	@echo "  submodule          Update git submodules"
	@echo ""
	@echo "Code Generation:"
	@echo "  repos              Generate sqlc code"
	@echo "  convert            Generate goverter converters"
	@echo ""
	@echo "Formatting:"
	@echo "  fmt                Format Go code"
	@echo "  fmt-check          Check Go code formatting"
	@echo "  fmt-sql            Format SQL files"
	@echo "  format             Format all (Go + SQL)"
	@echo "  format-check       Check all formatting"
	@echo ""
	@echo "Linting:"
	@echo "  lint               Run Go linter (requires proto)"
	@echo "  lint-fix           Run Go linter with auto-fix"
	@echo "  lint-sql           Lint SQL files"
	@echo "  lint-all           Run all linters"
	@echo ""
	@echo "Testing:"
	@echo "  test               Run tests"
	@echo "  test-coverage      Run tests with coverage"
	@echo "  test-verbose       Run tests with verbose output"
	@echo "  test-update-golden Update golden test files"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up          Start Docker containers"
	@echo "  docker-down        Stop Docker containers"
	@echo ""
	@echo "Maintenance:"
	@echo "  update-template    Update cruft template"
	@echo "  submodule          Update git submodules"
	@echo ""
	@echo "Help:"
	@echo "  help               Show this help message"
