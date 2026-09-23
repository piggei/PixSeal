#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD45_DIAGNOSTIC_DIR:-v4-phone private/build45-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-600}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required to summarize Build45 JSON" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind"

tsv="$OUTPUT_DIR/build45-blind-failure-decomposition.tsv"
md="$OUTPUT_DIR/build45-blind-failure-decomposition.md"
printf 'image\tclassification\tbuild41_direct\tbuild41_qualified\tbuild43_attempted\tbuild43_frozen\tbuild43_qualified\tbuild43_pair0\tbuild43_pair1\tbuild42_bank\tbuild42_ensembles\tbuild42_list_frames\thmac\n' > "$tsv"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem.json"
    set +e
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
    rc=$?
    set -e
    if [[ $rc -eq 124 ]]; then
        echo "error: Build45 blind diagnostic timed out for $file" >&2
        exit 1
    elif [[ $rc -ne 0 ]]; then
        echo "error: Build45 diagnostic command failed for $file (rc=$rc)" >&2
        exit "$rc"
    fi
    python3 - "$json" "$file" >> "$tsv" <<'PY'
import json, sys
j=json.load(open(sys.argv[1], encoding='utf-8'))
p=j.get('phone') or {}
vals=[
    sys.argv[2], j.get('classification',''),
    str(p.get('Build41DirectAccepted',False)).lower(), p.get('Build41QualifiedCandidates',0),
    str(p.get('Build43Attempted',False)).lower(), p.get('Build43FrozenCandidates',0), p.get('Build43QualifiedCandidates',0),
    p.get('Build43Pair0',''), p.get('Build43Pair1',''),
    p.get('Build42BankCandidates',0), p.get('Build42EnsemblesTried',0), p.get('Build42ListFramesTried',0),
    str(p.get('HMACAuthenticated',False)).lower(),
]
print('\t'.join(map(str,vals)))
PY
done

{
    echo '# PixSeal Build45 blind phone failure decomposition'
    echo
    echo 'Diagnostic only. Production Build44/Build43 search and ranking are unchanged.'
    echo
    echo '| image | class | B41 direct | B41 qualified | B43 | frozen | qualified | pair0 | pair1 | B42 bank | ensembles | list frames | HMAC |'
    echo '|---|---|---:|---:|---:|---:|---:|---|---|---:|---:|---:|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image class b41d b41q b43a b43f b43q p0 p1 b42b b42e b42l hmac; do
        echo "| $image | $class | $b41d | $b41q | $b43a | $b43f | $b43q | $p0 | $p1 | $b42b | $b42e | $b42l | $hmac |"
    done
    echo
    echo 'Build45 target: explain B/mild first; B/angle remains informational. No production parameter is changed by this target.'
} > "$md"

echo "Build45 blind failure decomposition written to:"
echo "  $tsv"
echo "  $md"
cat "$tsv"
