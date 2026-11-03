resolve-bin = $(shell command -v $(1) 2>/dev/null)
GO := $(call resolve-bin,go)
ifeq ($(strip $(GO)),)
GO := go
endif

check-tools:
	@command -v buf >/dev/null 2>&1 || { echo "buf not found. Run ../scripts/install-buf-deps.sh (from repo root: scripts/install-buf-deps.sh)."; exit 1; }
	@command -v protoc >/dev/null 2>&1 || { echo "protoc not found. Run ../scripts/install-buf-deps.sh (from repo root: scripts/install-buf-deps.sh)."; exit 1; }
	@command -v protoc-gen-go >/dev/null 2>&1 || { echo "protoc-gen-go not found. Run ../scripts/install-buf-deps.sh (from repo root: scripts/install-buf-deps.sh)."; exit 1; }
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "protoc-gen-go-grpc not found. Run ../scripts/install-buf-deps.sh (from repo root: scripts/install-buf-deps.sh)."; exit 1; }

generate: check-tools
	cd proto && buf build .
	cd proto && buf generate . --template buf.gen.yaml -o ../go
	@if [ ! -f go/go.mod ]; then \
		cd go && $(GO) mod init {{ cookiecutter.go_module_path.strip('/') }}/api/go; \
		echo "✅ go.mod created in api/go"; \
	else \
		cd go && $(GO) mod tidy; \
		echo "✅ go.mod updated in api/go"; \
	fi

lint: check-tools
	cd proto && buf lint .

format: check-tools
	cd proto && buf format -w .

breaking: check-tools
	cd proto && buf breaking . --against '.git#branch=main'

validate: check-tools
	cd proto && buf build .

clean:
	rm -rf go/
	rm -rf openapi/

proto: generate

.PHONY: check-tools generate proto lint format breaking validate clean
