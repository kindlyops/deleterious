#!/usr/bin/env bash
set -eo pipefail

echo "moving bazel outputs to goreleaser dist directory for packaging..."
mkdir -p dist/deleterious_darwin_amd64_v1
mkdir -p dist/deleterious_darwin_arm64
mkdir -p dist/deleterious_linux_amd64_v1
mkdir -p dist/deleterious_linux_arm64
mkdir -p dist/deleterious_windows_amd64_v1

sudo cp -f bdist/deleterious-darwin dist/deleterious_darwin_amd64_v1/deleterious
sudo cp -f bdist/deleterious-darwin-m1 dist/deleterious_darwin_arm64/deleterious
sudo cp -f bdist/deleterious-linux dist/deleterious_linux_amd64_v1/deleterious
sudo cp -f bdist/deleterious-linux-arm dist/deleterious_linux_arm64/deleterious
sudo cp -f bdist/deleterious-windows.exe dist/deleterious_windows_amd64_v1/deleterious.exe
