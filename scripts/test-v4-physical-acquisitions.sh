#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
FIXTURE_DIR="${V4_PHYSICAL_FIXTURE_DIR:-v4-physical private/build35-generated}"
ACQUISITION_DIR="${V4_PHYSICAL_ACQUISITION_DIR:-v4-physical private/build35-acquired}"
KEY="${V4_PHYSICAL_KEY:-PixSeal-v4-TestKey-2026}"
PLAN="${V4_PHYSICAL_PLAN:-$FIXTURE_DIR/acquisition-plan.tsv}"

if [[ ! -x "$PIXSEAL" ]]; then echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; fi
if [[ ! -f "$PLAN" ]]; then echo "error: acquisition plan not found: $PLAN" >&2; exit 1; fi
if [[ ! -d "$ACQUISITION_DIR" ]]; then echo "error: acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; fi

failures=0
seen=0
while IFS=$'\t' read -r fixture capture width height expected payload extra; do
    [[ "$fixture" == "fixture" || -z "$fixture" || "$fixture" == \#* ]] && continue
    seen=$((seen+1))
    if [[ -n "${extra:-}" || -z "$capture" || -z "$expected" ]]; then
        echo "FAIL invalid acquisition plan row for $fixture"
        failures=$((failures+1)); continue
    fi
    image="$ACQUISITION_DIR/$capture"
    if [[ ! -f "$image" ]]; then
        echo "FAIL $capture missing"
        failures=$((failures+1)); continue
    fi
    errlog="$(mktemp)"
    if output="$("$PIXSEAL" v4-extract-projective -in "$image" -key "$KEY" -width "$width" -height "$height" -raw 2>"$errlog")"; then
        if [[ "$expected" == "REJECT" ]]; then
            echo "FAIL $capture control unexpectedly authenticated payload=$(printf '%q' "$output")"
            failures=$((failures+1))
        elif [[ "$output" == "$payload" ]]; then
            echo "PASS $capture authenticated payload=$payload"
        else
            echo "FAIL $capture authenticated wrong payload=$(printf '%q' "$output") expected=$(printf '%q' "$payload")"
            failures=$((failures+1))
        fi
    else
        if [[ "$expected" == "REJECT" ]]; then
            echo "PASS $capture rejected as expected"
        else
            echo "FAIL $capture did not authenticate"
            sed 's/^/  /' "$errlog"
            failures=$((failures+1))
        fi
    fi
    rm -f "$errlog"
done < "$PLAN"

if (( seen == 0 )); then echo "error: acquisition plan contains no cases" >&2; exit 1; fi
if (( failures != 0 )); then
    echo "Physical qualification: FAIL ($failures case(s))"
    exit 1
fi
echo "Physical qualification: PASS ($seen cases)"
