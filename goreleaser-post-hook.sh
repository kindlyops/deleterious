#!/usr/bin/env bash
set -eo pipefail

echo "moving bazel outputs to goreleaser dist directory for packaging..."

cp bazel-bin/darwin_amd64_stripped/deleterious dist/deleterious_darwin_amd64/
cp bazel-bin/linux_amd64_pure_stripped/deleterious dist/deleterious_linux_amd64/

