#!/bin/bash
# Pre-download Bazel dependencies using curl (which properly handles proxy auth)
# This populates --distdir so Bazel doesn't need to download through proxy

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_DIR="$(dirname "$SCRIPT_DIR")"
DISTDIR="${WORKSPACE_DIR}/.bazel-distdir"

mkdir -p "$DISTDIR"

echo "Downloading Bazel dependencies to $DISTDIR..."
echo "Using proxy: ${HTTP_PROXY:-none}"

# Extract download URLs from WORKSPACE
# Format: urls = ["https://..."]
URLS=$(grep -oE 'https://[^"]+\.(tar\.gz|zip)' "$WORKSPACE_DIR/WORKSPACE" | sort -u)

# Resolve Starlark format strings like {0} by extracting version variables
buildtools_version=$(grep -oP 'buildtools_version = "\K[^"]+' "$WORKSPACE_DIR/WORKSPACE" || true)
if [ -n "$buildtools_version" ]; then
  URLS="${URLS//\{0\}/$buildtools_version}"
fi

downloaded=0
skipped=0
failed=0

while IFS= read -r url; do
  if [ -z "$url" ]; then
    continue
  fi

  # Skip URLs that still contain unresolved format strings
  if [[ "$url" == *"{"* ]]; then
    echo "⚠ Skipping unresolved URL: $url"
    failed=$((failed + 1))
    continue
  fi

  filename=$(basename "$url")
  output_file="$DISTDIR/$filename"

  if [ -f "$output_file" ]; then
    echo "✓ Already have: $filename"
    skipped=$((skipped + 1))
    continue
  fi

  echo "Downloading: $filename"
  if curl -f -L --retry 3 -o "$output_file.tmp" "$url" 2>/dev/null; then
    mv "$output_file.tmp" "$output_file"
    echo "✓ Downloaded: $filename"
    downloaded=$((downloaded + 1))
  else
    echo "✗ Failed: $filename"
    rm -f "$output_file.tmp"
    failed=$((failed + 1))
  fi
done <<< "$URLS"

echo ""
echo "Summary:"
echo "  Downloaded: $downloaded"
echo "  Skipped (cached): $skipped"
echo "  Failed: $failed"
echo "  Location: $DISTDIR"
echo ""
if [ "$failed" -eq 0 ]; then
  echo "✓ All dependencies downloaded successfully!"
  exit 0
else
  echo "⚠ Some downloads failed. Bazel may still work with cached/distdir files."
  exit 0  # Don't fail the script, downloads may be optional
fi
