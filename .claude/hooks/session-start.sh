#!/bin/bash
set -euo pipefail

# Only run in remote (Claude Code on the web) environments
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

INSTALL_DIR="/usr/local/bin"

# Install Bazelisk (provides the `bazel` command)
if ! command -v bazel &>/dev/null; then
  echo "Installing Bazelisk..."
  BAZELISK_VERSION="v1.25.0"
  curl -fsSL "https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64" \
    -o "${INSTALL_DIR}/bazel"
  chmod +x "${INSTALL_DIR}/bazel"
  echo "Bazelisk installed successfully"
fi

# Install ibazel (incremental Bazel)
if ! command -v ibazel &>/dev/null; then
  echo "Installing ibazel..."
  IBAZEL_VERSION="v0.25.3"
  curl -fsSL "https://github.com/bazelbuild/bazel-watcher/releases/download/${IBAZEL_VERSION}/ibazel_linux_amd64" \
    -o "${INSTALL_DIR}/ibazel"
  chmod +x "${INSTALL_DIR}/ibazel"
  echo "ibazel installed successfully"
fi

echo "Session startup complete: bazel and ibazel are available"
