#!/usr/bin/env bash
set -euo pipefail
PICS_DIR="${PICS_DIR:-original pics}"
ACTIVE_CORPUS_MANIFEST="${ACTIVE_CORPUS_MANIFEST:-private-corpus-active.tsv}"
if command -v magick >/dev/null 2>&1; then identify_tool=(magick identify); elif command -v identify >/dev/null 2>&1; then identify_tool=(identify); else identify_tool=(false); fi
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/test-common.sh"
verify_active_corpus_manifest "$PICS_DIR"
