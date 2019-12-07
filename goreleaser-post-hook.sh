#!/usr/bin/env bash
set -eo pipefail

echo "moving bazel outputs to goreleaser dist directory for packaging..."

cp bdist/deleterious-darwin dist/deleterious_darwin_amd64/deleterious
cp bdist/deleterious-linux dist/deleterious_linux_amd64/deleterious
cp bdist/deleterious-windows.exe dist/deleterious_windows_amd64/deleterious.exe

