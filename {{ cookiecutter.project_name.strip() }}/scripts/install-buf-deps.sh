#!/usr/bin/env bash
set -euo pipefail

ONLY="${ONLY:-all}"

BUF_VERSION="${BUF_VERSION:-1.56.0}"
PROTOC_VERSION="${PROTOC_VERSION:-31.1}"
PROTOC_GEN_GO_VERSION="${PROTOC_GEN_GO_VERSION:-v1.36.6}"
PROTOC_GEN_GO_GRPC_VERSION="${PROTOC_GEN_GO_GRPC_VERSION:-v1.5.1}"

install_buf() {
  if command -v buf >/dev/null 2>&1; then
    current_version="$(buf --version | cut -d' ' -f2)"
    if [ "$current_version" = "$BUF_VERSION" ]; then
      echo "buf $BUF_VERSION already installed"
      return
    fi
    echo "buf $current_version found, but $BUF_VERSION required. Skipping install."
    echo "To update manually: curl -sSL https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m) -o /tmp/buf && sudo mv /tmp/buf /usr/local/bin/buf && sudo chmod +x /usr/local/bin/buf"
    return
  fi

  os="$(uname -s)"
  arch="$(uname -m)"
  url="https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${os}-${arch}"
  echo "Downloading ${url}"
  curl -sSL "${url}" -o /tmp/buf
  sudo mv /tmp/buf /usr/local/bin/buf
  sudo chmod +x /usr/local/bin/buf
  echo "buf $BUF_VERSION installed"
}

install_protoc() {
  if command -v protoc >/dev/null 2>&1; then
    current_version="$(protoc --version | cut -d' ' -f2)"
    if [ "$current_version" = "$PROTOC_VERSION" ]; then
      echo "protoc $PROTOC_VERSION already installed"
      return
    fi
    echo "protoc $current_version found, but $PROTOC_VERSION required. Skipping install."
    echo "To update manually: curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/arm64/aarch_64/').zip -o /tmp/protoc.zip && cd /tmp && unzip -q protoc.zip && sudo mv bin/protoc /usr/local/bin/ && sudo chmod +x /usr/local/bin/protoc && rm -rf protoc.zip bin include"
    return
  fi

  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  if [ "$arch" = "arm64" ]; then
    arch="aarch_64"
  fi
  url="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-${os}-${arch}.zip"
  echo "Downloading ${url}"

  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT

  curl -sSL "${url}" -o "${tmp_dir}/protoc.zip"
  unzip -q "${tmp_dir}/protoc.zip" -d "${tmp_dir}"
  sudo mv "${tmp_dir}/bin/protoc" /usr/local/bin/
  sudo chmod +x /usr/local/bin/protoc
  echo "protoc $PROTOC_VERSION installed"
}

install_go_plugins() {
  echo "Installing protoc-gen-go ${PROTOC_GEN_GO_VERSION}..."
  go install google.golang.org/protobuf/cmd/protoc-gen-go@"${PROTOC_GEN_GO_VERSION}"
  echo "Installing protoc-gen-go-grpc ${PROTOC_GEN_GO_GRPC_VERSION}..."
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@"${PROTOC_GEN_GO_GRPC_VERSION}"
  echo "Go plugins installed"
}

case "$ONLY" in
  all)
    install_buf
    install_protoc
    install_go_plugins
    ;;
  buf)
    install_buf
    ;;
  protoc)
    install_protoc
    ;;
  go-plugins)
    install_go_plugins
    ;;
  *)
    echo "Unknown ONLY=${ONLY}. Expected one of: all, buf, protoc, go-plugins" >&2
    exit 1
    ;;
esac
