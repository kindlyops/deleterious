#!/usr/bin/env bash
set -eo pipefail

export RUNFILES_DIR="$PWD"/..
export PATH="$PWD/external/go_sdk/bin:$PATH"
gazelle="$PWD/$1"

echo "Using these commands"
command -v go
echo "$gazelle"

cd "$BUILD_WORKSPACE_DIRECTORY"

go mod tidy -compat=1.17
go mod vendor
$gazelle

# apply patches to generated BUILD files
# patch -p1 <cobra.BUILD.patch

# Fix vendored rules_go BUILD file to use external repository labels
VENDOR_BUILD_FILE="vendor/github.com/bazelbuild/rules_go/go/tools/bazel/BUILD.bazel"
if [ -f "$VENDOR_BUILD_FILE" ]; then
    echo "Patching $VENDOR_BUILD_FILE..."
    sed -i.bak 's|load("//go:def.bzl"|load("@io_bazel_rules_go//go:def.bzl"|g' "$VENDOR_BUILD_FILE"
    rm -f "$VENDOR_BUILD_FILE.bak"
    echo "Successfully patched vendored rules_go BUILD file"
fi
