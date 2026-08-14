#!/usr/bin/env bash
# Build macOS portable binary + .app install package (run on a Mac).
# Requires: Go, Xcode Command Line Tools, fyne CLI.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export PATH="$(go env GOPATH)/bin:$PATH"
if ! command -v fyne >/dev/null 2>&1; then
  go install fyne.io/tools/cmd/fyne@latest
fi

mkdir -p dist/portable dist/install
ARCH="$(uname -m)"
case "$ARCH" in
  arm64) PORTABLE="dist/portable/CopyTool-darwin-arm64" ;;
  x86_64) PORTABLE="dist/portable/CopyTool-darwin-amd64" ;;
  *) PORTABLE="dist/portable/CopyTool-darwin-$ARCH" ;;
esac

echo "==> portable $PORTABLE"
go build -ldflags="-s -w" -o "$PORTABLE" ./cmd/copytool
chmod +x "$PORTABLE"

echo "==> install package (.app)"
fyne package -os darwin -name CopyTool -appID com.vn.copytool -icon Icon.png -src ./cmd/copytool
if [[ -d CopyTool.app ]]; then
  zip -r "dist/install/CopyTool-darwin-${ARCH}.app.zip" CopyTool.app
  rm -rf CopyTool.app
fi

echo "Done."
ls -la dist/portable dist/install
