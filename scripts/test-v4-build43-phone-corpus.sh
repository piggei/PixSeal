#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD43_DIAGNOSTIC_DIR:-v4-phone private/build43-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-600}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|exploratory|$MESSAGE_B"
    "phone-b-angle.jpg|B|exploratory|$MESSAGE_B"
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
tsv="$OUTPUT_DIR/build43-phone-matrix.tsv"
md="$OUTPUT_DIR/build43-phone-matrix.md"
printf 'image\tclass\trole\tgeometry\tdata_decode\tbuild43_attempted\tbuild43_bank\tbuild43_ensembles\tbuild43_list_frames\tbuild43_authenticated\thmac\tpayload_match\tqualification\texit_code\n' > "$tsv"

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

    set +e
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw \
        >"$payload_file" 2>"$log"
    rc=$?
    set -e

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

    payload_match='-'
    qualification='INFO'
    if [[ "$role" == control ]]; then
        if [[ $rc -ne 0 && "$hmac" == false ]]; then
            qualification='PASS'
        else
            qualification='FAIL'
        fi
    else
        if [[ $rc -eq 0 ]]; then
            actual_payload="$(cat "$payload_file")"
            if [[ "$actual_payload" == "$expected_payload" ]]; then payload_match=true; else payload_match=false; fi
        else
            payload_match=false
        fi
        case "$role" in
            required-direct)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            required-build43)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            required-any)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            exploratory)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true ]]; then qualification='PASS'; else qualification='INFO'; fi
                ;;
        esac
    fi

    printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$file" "$class" "$role" "$geometry" "$data_decode" "$b43_attempted" "$b43_frozen" "$b43_qualified" "0" "$b43_auth" "$hmac" "$payload_match" "$qualification" "$rc" >> "$tsv"
done

{
    echo '# PixSeal Build43 smartphone qualification matrix'
    echo
    echo 'Build43 keeps Build41 geometry unchanged and adds a post-geometry qualified-bank/list-decoder fallback.'
    echo
    echo '| image | role | geometry | data | Build43 | bank | ensembles | list frames | HMAC | payload | qualification |'
    echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image class role geometry data_decode b43_attempted b43_frozen b43_qualified _ b43_auth hmac payload_match qualification rc; do
        echo "| $image | $role | $geometry | $data_decode | $b43_attempted/$b43_auth | $b43_frozen | $b43_qualified | - | $hmac | $payload_match | $qualification |"
    done
    echo
    echo 'Required gate: controls reject; A/front authenticates through Build43 geometry fallback; A/mild, A/angle and B/front remain authenticated; B/mild/B-angle are informational.'
} > "$md"

rm -f "$OUTPUT_DIR/logs/"*.payload.bin

control_passes="$(awk -F '\t' 'NR>1 && $3=="control" && $13=="PASS" {n++} END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-front.jpg" && $13=="PASS" {n++} END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-mild.jpg" && $13=="PASS" {n++} END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-angle.jpg" && $13=="PASS" {n++} END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-front.jpg" && $13=="PASS" {n++} END{print n+0}' "$tsv")"

echo "Build43 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1"
echo "Build43 phone qualification matrix written to:"
echo "  $tsv"
echo "  $md"

if (( control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 )); then
    echo 'error: Build43 physical gate not met' >&2
    exit 1
fi
echo 'Build43 physical gate: PASS'
