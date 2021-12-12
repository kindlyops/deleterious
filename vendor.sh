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
