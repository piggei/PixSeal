#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD65_DIAGNOSTIC_DIR:-v4-phone private/build65-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD65_TIMEOUT:-86400}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD65_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|required-build65|$MESSAGE_B"
    "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
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
tsv="$OUTPUT_DIR/build65-phone-matrix.tsv"
md="$OUTPUT_DIR/build65-phone-matrix.md"
printf 'image\tclass\trole\tgeometry\tdata_decode\tbuild43_attempted\tbuild43_frozen\tbuild43_qualified\tbuild43_authenticated\tbuild64_attempted\tbuild64_seeds\tbuild64_geometry_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode_candidates\tbuild64_list_frames\tbuild64_authenticated\tbuild65_attempted\tbuild65_workers\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\n' > "$tsv"

extract_re() {
    local pattern="$1" file="$2" default_value="${3:-}"
    local value
    value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n 1)"
    if [[ -z "$value" ]]; then printf '%s' "$default_value"; else printf '%s' "$value"; fi
}

for spec in "${specs[@]}"; do
    IFS='|' read -r file class role expected_payload <<< "$spec"
    input="$ACQUISITION_DIR/$file"
    stem="${file%.*}"
    log="$OUTPUT_DIR/logs/$stem.stderr.txt"
    payload_file="$OUTPUT_DIR/logs/$stem.payload.bin"

    start_ns="$(date +%s%N)"
    set +e
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw \
        >"$payload_file" 2>"$log"
    rc=$?
    set -e
    end_ns="$(date +%s%N)"
    elapsed_ms=$(( (end_ns - start_ns) / 1000000 ))

    geometry="$(extract_re 'phone-geometry-accepted: ([^ ]+).*' "$log" false)"
    data_decode="$(extract_re 'data-decode: attempted=([^ ]+).*' "$log" false)"
    hmac="$(extract_re 'hmac: authenticated=([^ ]+).*' "$log" false)"

    b43_attempted=false
    b43_frozen=0
    b43_qualified=0
    b43_auth=false
    if grep -q '^build43-geometry:' "$log"; then
        b43_attempted=true
        b43_frozen="$(extract_re 'build43-geometry: .* frozen=([0-9]+).*' "$log" 0)"
        b43_qualified="$(extract_re 'build43-geometry: .* qualified=([0-9]+).*' "$log" 0)"
        b43_auth="$(extract_re 'build43-geometry: .* authenticated=([^ ]+).*' "$log" false)"
    fi

    b64_attempted=false
    b64_seeds=0
    b64_evals=0
    b64_bank=0
    b64_qualified=0
    b64_decode=0
    b64_frames=0
    b64_auth=false
    if grep -q '^build64-recovery:' "$log"; then
        b64_attempted=true
        b64_seeds="$(extract_re 'build64-recovery: .* seeds=([0-9]+).*' "$log" 0)"
        b64_evals="$(extract_re 'build64-recovery: .* geometry-evals=([0-9]+).*' "$log" 0)"
        b64_bank="$(extract_re 'build64-recovery: .* bank=([0-9]+).*' "$log" 0)"
        b64_qualified="$(extract_re 'build64-recovery: .* qualified=([0-9]+).*' "$log" 0)"
        b64_decode="$(extract_re 'build64-recovery: .* decode-candidates=([0-9]+).*' "$log" 0)"
        b64_frames="$(extract_re 'build64-recovery: .* list-frames=([0-9]+).*' "$log" 0)"
        b64_auth="$(extract_re 'build64-recovery: .* authenticated=([^ ]+).*' "$log" false)"
    fi

    b65_attempted=false
    b65_workers=0
    if grep -q '^build65-parallel:' "$log"; then
        b65_attempted=true
        b65_workers="$(extract_re 'build65-parallel: .* workers=([0-9]+).*' "$log" 0)"
    fi

    case "$file" in
        phone-control-front.jpg) exp_evals=18021; exp_bank=26; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
        phone-control-mild.jpg)  exp_evals=117609; exp_bank=810; exp_qualified=6; exp_decode=6; exp_frames=18432 ;;
        phone-control-angle.jpg) exp_evals=21203; exp_bank=11; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
        phone-b-mild.jpg)        exp_evals=79259; exp_bank=937; exp_qualified=935; exp_decode=691; exp_frames=2120047 ;;
        phone-b-angle.jpg)       exp_evals=334857; exp_bank=6198; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
        *)                        exp_evals=0; exp_bank=0; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
    esac

    telemetry_equivalent=true
    if [[ "$exp_evals" -gt 0 ]]; then
        if [[ "$b64_attempted" != true || "$b64_seeds" -ne 24 || "$b64_evals" -ne "$exp_evals" || "$b64_bank" -ne "$exp_bank" || "$b64_qualified" -ne "$exp_qualified" || "$b64_decode" -ne "$exp_decode" || "$b64_frames" -ne "$exp_frames" || "$b65_attempted" != true || "$b65_workers" -lt 1 ]]; then
            telemetry_equivalent=false
        fi
    else
        if [[ "$b64_attempted" != false || "$b65_attempted" != false ]]; then
            telemetry_equivalent=false
        fi
    fi

    payload_match='-'
    qualification='INFO'
    if [[ "$role" == control ]]; then
        if [[ $rc -ne 0 && "$hmac" == false && "$b64_attempted" == true && "$b64_auth" == false ]]; then qualification='PASS'; else qualification='FAIL'; fi
    else
        if [[ $rc -eq 0 ]]; then
            actual_payload="$(cat "$payload_file")"
            if [[ "$actual_payload" == "$expected_payload" ]]; then payload_match=true; else payload_match=false; fi
        else
            payload_match=false
        fi
        case "$role" in
            required-direct)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false && "$b64_attempted" == false ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            required-build43)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true && "$b64_attempted" == false ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            required-any)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            required-build65)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == true && "$b64_auth" == true && "$b65_attempted" == true ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            required-reject)
                if [[ $rc -ne 0 && "$hmac" == false && "$b64_attempted" == true && "$b64_auth" == false && "$b65_attempted" == true ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
        esac
    fi
    if [[ "$telemetry_equivalent" != true ]]; then qualification='FAIL'; fi

    printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$file" "$class" "$role" "$geometry" "$data_decode" "$b43_attempted" "$b43_frozen" "$b43_qualified" "$b43_auth" "$b64_attempted" "$b64_seeds" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b65_attempted" "$b65_workers" "$hmac" "$payload_match" "$telemetry_equivalent" "$qualification" "$rc" "$elapsed_ms" >> "$tsv"
done

{
    echo '# PixSeal Build65 deterministic seed-parallel equivalence matrix'
    echo
    echo 'Build65 preserves the qualified Build64 search/decode semantics exactly, but schedules the 24 independent deep-recovery seeds concurrently and reassembles them in original seed order.'
    echo
    echo '| image | role | geometry | B43 auth | B64 evals | B64 bank | B64 qual | B64 decode | B64 frames | B65 workers | HMAC | payload | telemetry eq | ms | qualification |'
    echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image class role geometry data_decode b43_attempted b43_frozen b43_qualified b43_auth b64_attempted b64_seeds b64_evals b64_bank b64_qualified b64_decode b64_frames b64_auth b65_attempted b65_workers hmac payload_match telemetry_equivalent qualification rc elapsed_ms; do
        echo "| $image | $role | $geometry | $b43_auth | $b64_evals | $b64_bank | $b64_qualified | $b64_decode | $b64_frames | $b65_workers | $hmac | $payload_match | $telemetry_equivalent | $elapsed_ms | $qualification |"
    done
    echo
    echo 'Required gate: complete Build64 outcome equivalence, exact deep-bank/evaluation/decode telemetry on every fallback case, Build65 parallel telemetry present only when deep recovery runs, and the nine-photo qualification matrix remains PASS.'
} > "$md"

rm -f "$OUTPUT_DIR/logs/"*.payload.bin

control_passes="$(awk -F '\t' 'NR>1 && $3=="control" && $23=="PASS" {n++} END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-front.jpg" && $23=="PASS" {n++} END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-mild.jpg" && $23=="PASS" {n++} END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-angle.jpg" && $23=="PASS" {n++} END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-front.jpg" && $23=="PASS" {n++} END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-mild.jpg" && $23=="PASS" {n++} END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-angle.jpg" && $23=="PASS" {n++} END{print n+0}' "$tsv")"

echo "Build65 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
echo "Build65 equivalence phone matrix written to:"
echo "  $tsv"
echo "  $md"

if (( control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 || b_mild_pass != 1 || b_angle_pass != 1 )); then
    echo 'error: Build65 equivalence physical gate not met' >&2
    exit 1
fi
echo 'Build65 equivalence physical gate: PASS'
