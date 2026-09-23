#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD45_DIAGNOSTIC_DIR:-v4-phone private/build45-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
python3 - <<'PY' >/dev/null 2>&1 || { echo "error: Python OpenCV with SIFT support is required for this lab-only target" >&2; exit 1; }
import cv2
assert hasattr(cv2, 'SIFT_create')
PY
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
[[ -f "$REFERENCE" ]] || { echo "error: Build38 marked-B reference not found: $REFERENCE" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/oracle"

tsv="$OUTPUT_DIR/build45-reference-oracle.tsv"
md="$OUTPUT_DIR/build45-reference-oracle.md"
printf 'image\tmatches\tinliers\tinlier_fraction\tpilot_qualified\tproposal\tvalidation\tpilot_score\tpilot_margin\torigin_x\torigin_y\tprofiles\tlist_frames\thmac\tpayload\n' > "$tsv"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    quad="$OUTPUT_DIR/oracle/$stem-quad.json"
    out="$OUTPUT_DIR/oracle/$stem-oracle.json"
    python3 ./scripts/build45-reference-register.py --reference "$REFERENCE" --acquired "$input" --out "$quad"
    "$PIXSEAL" v4-diagnose-phone -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" \
        -oracle-quad-json "$quad" -oracle-only -json > "$out"
    python3 - "$quad" "$out" "$file" >> "$tsv" <<'PY'
import json, sys
q=json.load(open(sys.argv[1], encoding='utf-8'))
j=json.load(open(sys.argv[2], encoding='utf-8'))
o=j.get('oracle') or {}
vals=[
    sys.argv[3], q.get('matches',0), q.get('inliers',0), f"{q.get('inlier_fraction',0):.6f}",
    str(o.get('pilot_qualified',False)).lower(), f"{o.get('proposal_score',0):.6f}", f"{o.get('validation_score',0):.6f}",
    f"{o.get('pilot_score',0):.6f}", f"{o.get('pilot_margin',0):.6f}", o.get('origin_x_blocks',0), o.get('origin_y_blocks',0),
    o.get('profiles_tried',0), o.get('list_frames_tried',0), str(o.get('hmac_authenticated',False)).lower(), j.get('oracle_payload','')
]
print('\t'.join(map(str,vals)))
PY
done

{
    echo '# PixSeal Build45 reference-assisted oracle'
    echo
    echo '**LAB ONLY — never a production decoder path.** SIFT/reference geometry is supplied independently, then PixSeal tests the unchanged protected data channel.'
    echo
    echo '| image | matches | inliers | frac | pilot qualified | proposal | validation | pilot | margin | origin | profiles | list frames | HMAC | payload |'
    echo '|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|---:|---|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image matches inliers frac pq prop val ps pm ox oy prof lf hmac payload; do
        echo "| $image | $matches | $inliers | $frac | $pq | $prop | $val | $ps | $pm | ($ox,$oy) | $prof | $lf | $hmac | $payload |"
    done
    echo
    echo 'Interpretation: oracle HMAC PASS proves the data channel is viable when geometry is supplied independently; it does not qualify blind recovery.'
} > "$md"

echo "Build45 lab oracle written to:"
echo "  $tsv"
echo "  $md"
cat "$tsv"
