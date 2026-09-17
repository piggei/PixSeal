#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
CONTROL="${V4_SCANNER_CONTROL:?set V4_SCANNER_CONTROL to the unmarked full-page scan}"
MARKED_A="${V4_SCANNER_MARKED_A:?set V4_SCANNER_MARKED_A to marked-A full-page scan}"
MARKED_B="${V4_SCANNER_MARKED_B:?set V4_SCANNER_MARKED_B to marked-B full-page scan}"
KEY="${V4_PHYSICAL_KEY:-PixSeal-v4-TestKey-2026}"
WIDTH="${V4_SCANNER_CANONICAL_WIDTH:-1632}"
HEIGHT="${V4_SCANNER_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHYSICAL_MESSAGE_A:-v4-b35-phys-a}"
MESSAGE_B="${V4_PHYSICAL_MESSAGE_B:-v4-b35-phys-b}"

for f in "$CONTROL" "$MARKED_A" "$MARKED_B"; do
    [[ -f "$f" ]] || { echo "error: scanner capture not found: $f" >&2; exit 1; }
done

failures=0
run_marked() {
    local label="$1" file="$2" expected="$3" errlog output
    errlog="$(mktemp)"
    if output="$($PIXSEAL v4-extract-scanner -in "$file" -key "$KEY" -width "$WIDTH" -height "$HEIGHT" -raw 2>"$errlog")"; then
        if [[ "$output" == "$expected" ]]; then
            echo "PASS $label authenticated payload=$expected"
        else
            echo "FAIL $label wrong payload=$(printf '%q' "$output") expected=$(printf '%q' "$expected")"
            failures=$((failures+1))
        fi
    else
        echo "FAIL $label did not authenticate"
        sed 's/^/  /' "$errlog"
        failures=$((failures+1))
    fi
    rm -f "$errlog"
}

errlog="$(mktemp)"
if output="$($PIXSEAL v4-extract-scanner -in "$CONTROL" -key "$KEY" -width "$WIDTH" -height "$HEIGHT" -raw 2>"$errlog")"; then
    echo "FAIL control unexpectedly authenticated payload=$(printf '%q' "$output")"
    failures=$((failures+1))
else
    echo "PASS control rejected as expected"
fi
rm -f "$errlog"

run_marked "marked-a" "$MARKED_A" "$MESSAGE_A"
run_marked "marked-b" "$MARKED_B" "$MESSAGE_B"

if (( failures != 0 )); then
    echo "Build37 blind scanner qualification: FAIL ($failures case(s))"
    exit 1
fi

echo "Build37 blind scanner qualification: PASS (3 cases)"
