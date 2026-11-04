#!/usr/bin/env bash
set -euo pipefail

# ONLY=golangci-lint|mockery|sql-formatter installs one tool (default: all).
ONLY="${ONLY:-all}"

# GOLANGCI_LINT_VERSION matches lint-version in .github/workflows/lint.yml.
GOLANGCI_LINT_VERSION="${GOLANGCI_LINT_VERSION:-v2.15.0}"
MOCKERY_VERSION="${MOCKERY_VERSION:-v3.9.0}"
SQL_FORMATTER_VERSION="${SQL_FORMATTER_VERSION:-15.9.0}"

install_golangci_lint() {
  echo "Installing golangci-lint ${GOLANGCI_LINT_VERSION}..."
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@"${GOLANGCI_LINT_VERSION}"
}

install_mockery() {
  echo "Installing mockery ${MOCKERY_VERSION}..."
  go install github.com/vektra/mockery/v3@"${MOCKERY_VERSION}"
}

install_sql_formatter() {
  if command -v npm >/dev/null 2>&1; then
    echo "Installing sql-formatter ${SQL_FORMATTER_VERSION}..."
    npm install -g "sql-formatter@${SQL_FORMATTER_VERSION}"
  else
    echo "npm not found. Please install Node.js to use sql-formatter."
  fi
}

case "$ONLY" in
  all)
    install_golangci_lint
    install_mockery
    install_sql_formatter
    ;;
  golangci-lint) install_golangci_lint ;;
  mockery) install_mockery ;;
  sql-formatter) install_sql_formatter ;;
  *)
    echo "Unknown ONLY=${ONLY}. Expected one of: all, golangci-lint, mockery, sql-formatter" >&2
    exit 1
    ;;
esac
