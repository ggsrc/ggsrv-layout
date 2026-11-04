#!/usr/bin/env bash
# Install sqlc from a Stumble/sqlc release (the fork that generates wpgx code),
# after checking the archive against the SHA-256 pinned below. `make repos`
# requires exactly this version, so every machine generates the same code.
#
#   scripts/install-sqlc.sh                  # into $(go env GOPATH)/bin
#   INSTALL_DIR=/usr/local/bin scripts/install-sqlc.sh
#
# To move to another release, take the values from its checksums.txt.
set -euo pipefail

SQLC_VERSION="v2.4.1"

sha256_for() {
  case "$1" in
    linux_amd64) echo "f9378d4db9ed7aef1aec974abecb4b784ed7d88a51a26031695ac4fbe2fe8e3f" ;;
    linux_arm64) echo "bb447a86aed9389775be86d50fda95153b9d8193272d67f7208fd273df3264bb" ;;
    darwin_amd64) echo "7dd58a8e2568add9406afcffc6b1f495e0be62739a7d69d5e1d50df3d160bfe0" ;;
    darwin_arm64) echo "5e2a8ce9f24f9e6fceff8d0a872810a341e52ccf91ca933f879b2f86406a0f64" ;;
    *) return 1 ;;
  esac
}

BASE_URL="${SQLC_BASE_URL:-https://github.com/Stumble/sqlc/releases/download}"
INSTALL_DIR="${INSTALL_DIR:-$(go env GOPATH)/bin}"

case "$(uname -s)" in
  Linux) os=linux ;;
  Darwin) os=darwin ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  *) echo "unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac
want="$(sha256_for "${os}_${arch}")"

if [ -x "$INSTALL_DIR/sqlc" ] && [ "$("$INSTALL_DIR/sqlc" version)" = "${SQLC_VERSION}-wicked-fork" ]; then
  echo "sqlc ${SQLC_VERSION}-wicked-fork already installed in $INSTALL_DIR"
  exit 0
fi

archive="sqlc_${SQLC_VERSION#v}_${os}_${arch}.tar.gz"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "Downloading ${archive}..."
curl -fsSL -o "$tmp/$archive" "$BASE_URL/$SQLC_VERSION/$archive"
if command -v sha256sum >/dev/null 2>&1; then
  got="$(sha256sum "$tmp/$archive" | awk '{ print $1 }')"
else
  got="$(shasum -a 256 "$tmp/$archive" | awk '{ print $1 }')"
fi
if [ "$got" != "$want" ]; then
  echo "SHA-256 mismatch for $archive: got $got, want $want" >&2
  exit 1
fi

tar -xzf "$tmp/$archive" -C "$tmp" sqlc
mkdir -p "$INSTALL_DIR"
mv "$tmp/sqlc" "$INSTALL_DIR/sqlc"
chmod 0755 "$INSTALL_DIR/sqlc"
echo "Installed $INSTALL_DIR/sqlc: $("$INSTALL_DIR/sqlc" version)"
