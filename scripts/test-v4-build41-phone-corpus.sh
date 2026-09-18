#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD41_DIAGNOSTIC_DIR:-v4-phone private/build41-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-600}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$CANONICAL_WIDTH" =~ ^[0-9]+$ && "$CANONICAL_HEIGHT" =~ ^[0-9]+$ ]] || { echo "error: canonical dimensions must be integers" >&2; exit 1; }
(( CANONICAL_WIDTH >= 296 && CANONICAL_HEIGHT >= 256 && CANONICAL_WIDTH % 8 == 0 && CANONICAL_HEIGHT % 8 == 0 )) || {
    echo "error: canonical dimensions must be block-aligned and at least 296x256" >&2
    exit 1
}
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|REJECT|-"
    "phone-control-mild.jpg|control|REJECT|-"
    "phone-control-angle.jpg|control|REJECT|-"
    "phone-a-front.jpg|A|ACCEPT|$MESSAGE_A"
    "phone-a-mild.jpg|A|ACCEPT|$MESSAGE_A"
    "phone-a-angle.jpg|A|ACCEPT|$MESSAGE_A"
    "phone-b-front.jpg|B|ACCEPT|$MESSAGE_B"
    "phone-b-mild.jpg|B|ACCEPT|$MESSAGE_B"
    "phone-b-angle.jpg|B|ACCEPT|$MESSAGE_B"
)

missing=0
for spec in "${specs[@]}"; do
    IFS='|' read -r file _ <<< "$spec"
    if [[ ! -f "$ACQUISITION_DIR/$file" ]]; then
        echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2
        missing=1
    fi
done
(( missing == 0 )) || exit 1

mkdir -p "$OUTPUT_DIR/logs"
tsv="$OUTPUT_DIR/build41-phone-matrix.tsv"
md="$OUTPUT_DIR/build41-phone-matrix.md"

printf 'image\tclass\texpected\tboundary\tboundary_confidence\tprojective_basin\tgeometry_accepted\tensemble\tproposal\theldout\tpilot_origin\tresidual_attempted\tresidual_fitted\tresidual_applied\tresidual_controls\tresidual_rms_px\theldout_before\theldout_after\theldout_delta\tdata_decode\tsoft_hamming_profiles\tmax_data_confidence\thmac\tfallback_attempted\tfallback_authenticated\tpayload_match\tqualification\texit_code\n' > "$tsv"

extract_re() {
    local pattern="$1" file="$2" default_value="${3:-}"
    local value
    value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n 1)"
    if [[ -z "$value" ]]; then
        printf '%s' "$default_value"
    else
        printf '%s' "$value"
    fi
}

for spec in "${specs[@]}"; do
    IFS='|' read -r file class expected expected_payload <<< "$spec"
    input="$ACQUISITION_DIR/$file"
    stem="${file%.*}"
    log="$OUTPUT_DIR/logs/$stem.stderr.txt"
    payload_file="$OUTPUT_DIR/logs/$stem.payload.bin"

    set +e
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone \
        -in "$input" \
        -key "$KEY" \
        -width "$CANONICAL_WIDTH" \
        -height "$CANONICAL_HEIGHT" \
        -raw >"$payload_file" 2>"$log"
    rc=$?
    set -e

    boundary="$(extract_re 'boundary: detected=([^ ]+).*' "$log" false)"
    boundary_conf="$(extract_re 'boundary: detected=[^ ]+ confidence=([^ ]+).*' "$log" 0)"
    projective="$(extract_re 'projective-basin: found=([^ ]+).*' "$log" false)"
    geometry="$(extract_re 'phone-geometry-accepted: ([^ ]+).*' "$log" false)"
    ensemble="$(extract_re 'phone-geometry-accepted: [^ ]+ ensemble=([0-9]+).*' "$log" 0)"
    proposal="$(extract_re 'proposal: ([^ ]+).*' "$log" 0)"
    heldout="$(extract_re 'validation: ([^ ]+).*' "$log" 0)"
    pilot_origin="$(extract_re 'pilot-origin: \(([-0-9]+,[-0-9]+)\) blocks.*' "$log" '-')"

    residual_attempted=false
    residual_fitted=false
    residual_applied=false
    residual_controls=0
    residual_rms=0
    heldout_before=0
    heldout_after=0
    heldout_delta=0
    if grep -q '^phone-residual:' "$log"; then
        residual_attempted=true
        residual_fitted="$(extract_re 'phone-residual: fitted=([^ ]+).*' "$log" false)"
        residual_applied="$(extract_re 'phone-residual: fitted=[^ ]+ applied=([^ ]+).*' "$log" false)"
        residual_controls="$(extract_re 'phone-residual: .* controls=([0-9]+).*' "$log" 0)"
        residual_rms="$(extract_re 'phone-residual: .* rms=([^ ]+) px.*' "$log" 0)"
        heldout_before="$(extract_re 'phone-residual: .* validation=([^ ]+)->[^ ]+.*' "$log" 0)"
        heldout_after="$(extract_re 'phone-residual: .* validation=[^ ]+->([^ ]+).*' "$log" 0)"
        heldout_delta="$(awk -v a="$heldout_after" -v b="$heldout_before" 'BEGIN {printf "%.6f", a-b}')"
    fi

    data_decode="$(extract_re 'data-decode: attempted=([^ ]+).*' "$log" false)"
    soft_profiles="$(extract_re 'data-decode: .* soft-hamming-profiles=([0-9]+).*' "$log" 0)"
    max_conf="$(extract_re 'data-decode: .* max-confidence=([^ ]+).*' "$log" 0)"
    hmac="$(extract_re 'hmac: authenticated=([^ ]+).*' "$log" false)"
    fallback_attempted="$(extract_re 'hmac: .* fallback-attempted=([^ ]+).*' "$log" false)"
    fallback_authenticated="$(extract_re 'hmac: .* fallback-authenticated=([^ ]+).*' "$log" false)"

    payload_match='-'
    qualification='ERROR'
    if grep -q '^hmac:' "$log"; then
        qualification='FAIL'
        if [[ "$expected" == 'REJECT' ]]; then
            if [[ "$hmac" == 'false' && $rc -ne 0 ]]; then
                qualification='PASS'
            fi
        else
            if [[ $rc -eq 0 ]]; then
                actual_payload="$(cat "$payload_file")"
                if [[ "$actual_payload" == "$expected_payload" ]]; then
                    payload_match='true'
                    qualification='PASS'
                else
                    payload_match='false'
                fi
            else
                payload_match='false'
            fi
        fi
    fi

    printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$file" "$class" "$expected" "$boundary" "$boundary_conf" "$projective" "$geometry" "$ensemble" "$proposal" "$heldout" "$pilot_origin" \
        "$residual_attempted" "$residual_fitted" "$residual_applied" "$residual_controls" "$residual_rms" "$heldout_before" "$heldout_after" "$heldout_delta" \
        "$data_decode" "$soft_profiles" "$max_conf" "$hmac" "$fallback_attempted" "$fallback_authenticated" "$payload_match" "$qualification" "$rc" >> "$tsv"

done

{
    echo '# PixSeal Build41 strength-48 phone qualification matrix'
    echo
    echo "Generated from the Build41 blind phone decoder. Canonical carrier: ${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}."
    echo
    echo '| image | expected | boundary | projective | geom | proposal | held-out | origin | residual fit/applied | held-out delta | soft-Hamming | HMAC | qualification |'
    echo '|---|---|---:|---:|---:|---:|---:|---|---|---:|---:|---:|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image class expected boundary boundary_conf projective geometry ensemble proposal heldout pilot_origin residual_attempted residual_fitted residual_applied residual_controls residual_rms heldout_before heldout_after heldout_delta data_decode soft_profiles max_conf hmac fallback_attempted fallback_authenticated payload_match qualification rc; do
        residual_cell="${residual_fitted}/${residual_applied}"
        echo "| $image | $expected | $boundary | $projective | $geometry ($ensemble) | $proposal | $heldout | $pilot_origin | $residual_cell | $heldout_delta | $soft_profiles | $hmac | $qualification |"
    done
    echo
    echo 'Build41 staged milestone: every control must reject, and at least one A plus one B capture must recover the exact authenticated payload. Pilot/Hamming telemetry never ranks geometry.'
} > "$md"

rm -f "$OUTPUT_DIR/logs/"*.payload.bin

control_passes="$(awk -F '\t' 'NR>1 && $2=="control" && $27=="PASS" {n++} END{print n+0}' "$tsv")"
a_passes="$(awk -F '\t' 'NR>1 && $2=="A" && $27=="PASS" {n++} END{print n+0}' "$tsv")"
b_passes="$(awk -F '\t' 'NR>1 && $2=="B" && $27=="PASS" {n++} END{print n+0}' "$tsv")"

echo "Build41 staged result: controls=${control_passes}/3 A=${a_passes}/3 B=${b_passes}/3"

echo "Build41 phone qualification matrix written to:"
echo "  $tsv"
echo "  $md"
echo "Per-image decoder stderr logs: $OUTPUT_DIR/logs/"

if (( control_passes != 3 || a_passes < 1 || b_passes < 1 )); then
    echo "error: Build41 staged physical milestone not met" >&2
    exit 1
fi
echo "Build41 staged physical milestone: PASS"
