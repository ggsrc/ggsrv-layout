#!/usr/bin/env bash
set -euo pipefail

V8RUNNER_VERSION="${V8RUNNER_VERSION:-v0.0.5}"

CGO_ENABLED=1 go install github.com/stumble/v8runner/cmd/v8runner@"${V8RUNNER_VERSION}"
