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

TAR_NAME="koomos_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${TAR_NAME}"

echo "Downloading: $URL"
curl -sL "$URL" | tar -xz

chmod +x koomos
sudo mv koomos /usr/local/bin/

echo "Koomos installed!"
koomos --version || echo "Run 'koomos --help' to get started."
