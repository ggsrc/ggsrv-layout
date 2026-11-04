#!/usr/bin/env bash
set -euo pipefail

# GOLANGCI_LINT_VERSION matches lint-version in .github/workflows/lint.yml.
GOLANGCI_LINT_VERSION="${GOLANGCI_LINT_VERSION:-v2.15.0}"
MOCKERY_VERSION="${MOCKERY_VERSION:-v2.53.8}"
SQL_FORMATTER_VERSION="${SQL_FORMATTER_VERSION:-15.9.0}"

echo "Installing golangci-lint ${GOLANGCI_LINT_VERSION}..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@"${GOLANGCI_LINT_VERSION}"

echo "Installing mockery ${MOCKERY_VERSION}..."
go install github.com/vektra/mockery/v2@"${MOCKERY_VERSION}"

if command -v npm >/dev/null 2>&1; then
  echo "Installing sql-formatter ${SQL_FORMATTER_VERSION}..."
  if [ "$SQL_FORMATTER_VERSION" = "latest" ]; then
    npm install -g sql-formatter
  else
    npm install -g "sql-formatter@${SQL_FORMATTER_VERSION}"
  fi
else
  echo "npm not found. Please install Node.js to use sql-formatter."
fi
