#!/bin/bash

set -e

REPO="Pitchouneee/koomos"
VERSION=${1:-"latest"}
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture name
if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
  ARCH="arm64"
fi

# If latest, fetch real version tag
if [[ "$VERSION" == "latest" ]]; then
  VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | cut -d'"' -f4)
fi

# Remove "v" prefix for archive naming
VERSION_NO_PREFIX="${VERSION#v}"

TAR_NAME="koomos_${VERSION_NO_PREFIX}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${TAR_NAME}"

echo "Downloading: $URL"
if ! curl -sL "$URL" | tar -xz; then
  echo "Failed to download or extract archive. URL may be invalid."
  exit 1
fi

chmod +x koomos
sudo mv koomos /usr/local/bin/

echo "Koomos installed!"
koomos --version || echo "Run 'koomos --help' to get started."
