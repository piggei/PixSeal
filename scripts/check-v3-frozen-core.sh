#!/usr/bin/env bash
set -euo pipefail
manifest="${1:-docs/V3_FROZEN_CORE_SHA256.txt}"
if [[ ! -f "$manifest" ]]; then
  echo "error: v3 frozen-core manifest not found: $manifest" >&2
  exit 2
fi
# Ignore comments/blank lines but otherwise use the standard sha256sum format.
tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT
grep -Ev '^[[:space:]]*(#|$)' "$manifest" > "$tmp"
sha256sum -c "$tmp"
echo "Format-v3 frozen core matches the build23 reference manifest."
