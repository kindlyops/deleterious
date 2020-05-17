#!/usr/bin/env bash
set -eo pipefail

echo "moving bazel outputs to goreleaser dist directory for packaging..."
mkdir -p dist/deleterious_darwin_amd64
mkdir -p dist/deleterious_linux_amd64
mkdir -p dist/deleterious_windows_amd64

sudo cp bdist/deleterious-darwin dist/deleterious_darwin_amd64/deleterious
sudo cp bdist/deleterious-linux dist/deleterious_linux_amd64/deleterious
sudo cp bdist/deleterious-windows.exe dist/deleterious_windows_amd64/deleterious.exe

