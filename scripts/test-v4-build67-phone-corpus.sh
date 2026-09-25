#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD67_DIAGNOSTIC_DIR:-v4-phone private/build67-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD67_TIMEOUT:-86400}"
BUILD66_BASELINE_TSV="${V4_PHONE_BUILD66_BASELINE_TSV:-v4-phone private/build66-diagnostics/build66-phone-matrix.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD67_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|required-build67|$MESSAGE_B"
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
tsv="$OUTPUT_DIR/build67-phone-profile.tsv"
md="$OUTPUT_DIR/build67-phone-profile.md"
printf 'image\tclass\trole\tgeometry\tdata_decode\tbuild43_attempted\tbuild43_frozen\tbuild43_qualified\tbuild43_authenticated\tbuild64_attempted\tbuild64_seeds\tbuild64_geometry_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode_candidates\tbuild64_list_frames\tbuild64_authenticated\tbuild65_attempted\tbuild65_workers\tbuild66_attempted\tbuild66_workers\tbuild67_attempted\tbuild67_total_ms\tbuild67_geometry_ms\tbuild67_plane_prep_ms\tbuild67_qualification_ms\tbuild67_decode_wall_ms\tbuild67_physical_decode_candidates\tbuild67_speculative_candidates\tbuild67_physical_profiles\tbuild67_physical_list_frames\tbuild67_sampling_worker_ms\tbuild67_list_worker_ms\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\tbuild66_elapsed_ms\tspeedup_vs_build66\n' > "$tsv"

extract_re() {
    local pattern="$1" file="$2" default_value="${3:-}"
    local value
    value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n 1)"
    if [[ -z "$value" ]]; then printf '%s' "$default_value"; else printf '%s' "$value"; fi
}

resolve_build66_baseline() {
    if [[ -f "$BUILD66_BASELINE_TSV" ]]; then
        printf '%s' "$BUILD66_BASELINE_TSV"
        return
    fi
    # The originally qualified Build66 diagnostics were packaged before the
    # naming-only correction and may still carry the accidental Build65 name.
    local legacy
    legacy="$(dirname "$BUILD66_BASELINE_TSV")/build65-phone-matrix.tsv"
    if [[ -f "$legacy" ]]; then
        printf '%s' "$legacy"
        return
    fi
    printf '%s' ''
}

BASELINE_TSV="$(resolve_build66_baseline)"

baseline_elapsed_ms() {
    local image="$1"
    [[ -n "$BASELINE_TSV" && -f "$BASELINE_TSV" ]] || { printf '%s' '-'; return; }
    awk -F '\t' -v image="$image" '
        NR==1 { for (i=1; i<=NF; i++) if ($i=="elapsed_ms") col=i; next }
        $1==image && col>0 { print $col; found=1; exit }
        END { if (!found) print "-" }
    ' "$BASELINE_TSV"
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

    b66_attempted=false
    b66_workers=0
    if grep -q '^build66-decode-parallel:' "$log"; then
        b66_attempted=true
        b66_workers="$(extract_re 'build66-decode-parallel: .* workers=([0-9]+).*' "$log" 0)"
    fi

    b67_attempted=false
    b67_total_ms=0
    b67_geometry_ms=0
    b67_plane_prep_ms=0
    b67_qualification_ms=0
    b67_decode_wall_ms=0
    b67_physical_decode=0
    b67_speculative=0
    b67_physical_profiles=0
    b67_physical_frames=0
    b67_sampling_worker_ms=0
    b67_list_worker_ms=0
    if grep -q '^build67-profile:' "$log"; then
        b67_attempted=true
        b67_total_ms="$(extract_re 'build67-profile: .* total-ms=([0-9]+).*' "$log" 0)"
        b67_geometry_ms="$(extract_re 'build67-profile: .* geometry-ms=([0-9]+).*' "$log" 0)"
        b67_plane_prep_ms="$(extract_re 'build67-profile: .* plane-prep-ms=([0-9]+).*' "$log" 0)"
        b67_qualification_ms="$(extract_re 'build67-profile: .* qualification-ms=([0-9]+).*' "$log" 0)"
        b67_decode_wall_ms="$(extract_re 'build67-profile: .* decode-wall-ms=([0-9]+).*' "$log" 0)"
        b67_physical_decode="$(extract_re 'build67-profile: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"
        b67_speculative="$(extract_re 'build67-profile: .* speculative-candidates=([0-9]+).*' "$log" 0)"
        b67_physical_profiles="$(extract_re 'build67-profile: .* physical-profiles=([0-9]+).*' "$log" 0)"
        b67_physical_frames="$(extract_re 'build67-profile: .* physical-list-frames=([0-9]+).*' "$log" 0)"
        b67_sampling_worker_ms="$(extract_re 'build67-profile: .* sampling-worker-ms=([0-9]+).*' "$log" 0)"
        b67_list_worker_ms="$(extract_re 'build67-profile: .* list-worker-ms=([0-9]+).*' "$log" 0)"
    fi

    case "$file" in
        phone-control-front.jpg) exp_evals=18021; exp_bank=26; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
        phone-control-mild.jpg)  exp_evals=117609; exp_bank=810; exp_qualified=6; exp_decode=6; exp_frames=18432 ;;
        phone-control-angle.jpg) exp_evals=21203; exp_bank=11; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
        phone-b-mild.jpg)        exp_evals=79259; exp_bank=937; exp_qualified=935; exp_decode=691; exp_frames=2120047 ;;
        phone-b-angle.jpg)       exp_evals=334857; exp_bank=6198; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
        *)                       exp_evals=0; exp_bank=0; exp_qualified=0; exp_decode=0; exp_frames=0 ;;
    esac

    telemetry_equivalent=true
    if [[ "$exp_evals" -gt 0 ]]; then
        if [[ "$b64_attempted" != true || "$b64_seeds" -ne 24 || "$b64_evals" -ne "$exp_evals" || "$b64_bank" -ne "$exp_bank" || "$b64_qualified" -ne "$exp_qualified" || "$b64_decode" -ne "$exp_decode" || "$b64_frames" -ne "$exp_frames" || "$b65_attempted" != true || "$b65_workers" -lt 1 || "$b66_attempted" != true || "$b67_attempted" != true ]]; then
            telemetry_equivalent=false
        fi
        if [[ "$exp_qualified" -gt 0 ]]; then
            if [[ "$b66_workers" -lt 1 || "$b67_physical_decode" -lt "$b64_decode" || "$b67_speculative" -ne $((b67_physical_decode - b64_decode)) || "$b67_physical_frames" -lt "$b64_frames" ]]; then
                telemetry_equivalent=false
            fi
            if [[ "$b64_auth" == true && "$b67_speculative" -ge "$b66_workers" ]]; then
                telemetry_equivalent=false
            fi
            if [[ "$b64_auth" == false && "$b67_speculative" -ne 0 ]]; then
                telemetry_equivalent=false
            fi
        else
            if [[ "$b66_workers" -ne 0 || "$b67_physical_decode" -ne 0 || "$b67_speculative" -ne 0 || "$b67_physical_frames" -ne 0 || "$b67_decode_wall_ms" -ne 0 ]]; then
                telemetry_equivalent=false
            fi
        fi
    else
        if [[ "$b64_attempted" != false || "$b65_attempted" != false || "$b66_attempted" != false || "$b67_attempted" != false ]]; then
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
            required-build67)
                if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == true && "$b64_auth" == true && "$b65_attempted" == true && "$b66_attempted" == true && "$b66_workers" -gt 0 && "$b67_attempted" == true ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
            required-reject)
                if [[ $rc -ne 0 && "$hmac" == false && "$b64_attempted" == true && "$b64_auth" == false && "$b65_attempted" == true && "$b66_attempted" == true && "$b67_attempted" == true ]]; then qualification='PASS'; else qualification='FAIL'; fi
                ;;
        esac
    fi
    if [[ "$telemetry_equivalent" != true ]]; then qualification='FAIL'; fi

    baseline_ms="$(baseline_elapsed_ms "$file")"
    speedup='-'
    if [[ "$baseline_ms" =~ ^[0-9]+$ && "$elapsed_ms" -gt 0 ]]; then
        speedup="$(awk -v b="$baseline_ms" -v c="$elapsed_ms" 'BEGIN { printf "%.3f", b/c }')"
    fi

    printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$file" "$class" "$role" "$geometry" "$data_decode" "$b43_attempted" "$b43_frozen" "$b43_qualified" "$b43_auth" "$b64_attempted" "$b64_seeds" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b65_attempted" "$b65_workers" "$b66_attempted" "$b66_workers" "$b67_attempted" "$b67_total_ms" "$b67_geometry_ms" "$b67_plane_prep_ms" "$b67_qualification_ms" "$b67_decode_wall_ms" "$b67_physical_decode" "$b67_speculative" "$b67_physical_profiles" "$b67_physical_frames" "$b67_sampling_worker_ms" "$b67_list_worker_ms" "$hmac" "$payload_match" "$telemetry_equivalent" "$qualification" "$rc" "$elapsed_ms" "$baseline_ms" "$speedup" >> "$tsv"
done

{
    echo '# PixSeal Build67 deep-recovery profiling matrix'
    echo
    echo 'Build67 is an observability-only successor to the qualified Build66 smartphone baseline. It preserves Build66 scheduling and Build64 logical semantics while measuring stage wall time and physical protected-data work.'
    echo
    echo '| image | role | B64 bank | B64 qual | logical decode | logical frames | B66 workers | physical decode | speculative | physical frames | geometry ms | plane prep ms | qualification ms | decode wall ms | sampling worker ms | list worker ms | HMAC | telemetry eq | ms | B66 ms | ratio | qualification |'
    echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image class role geometry data_decode b43_attempted b43_frozen b43_qualified b43_auth b64_attempted b64_seeds b64_evals b64_bank b64_qualified b64_decode b64_frames b64_auth b65_attempted b65_workers b66_attempted b66_workers b67_attempted b67_total_ms b67_geometry_ms b67_plane_prep_ms b67_qualification_ms b67_decode_wall_ms b67_physical_decode b67_speculative b67_physical_profiles b67_physical_frames b67_sampling_worker_ms b67_list_worker_ms hmac payload_match telemetry_equivalent qualification rc elapsed_ms baseline_ms speedup; do
        echo "| $image | $role | $b64_bank | $b64_qualified | $b64_decode | $b64_frames | $b66_workers | $b67_physical_decode | $b67_speculative | $b67_physical_frames | $b67_geometry_ms | $b67_plane_prep_ms | $b67_qualification_ms | $b67_decode_wall_ms | $b67_sampling_worker_ms | $b67_list_worker_ms | $hmac | $telemetry_equivalent | $elapsed_ms | $baseline_ms | $speedup | $qualification |"
    done
    echo
    echo 'The Build67 correctness gate is semantic only: every Build64/65/66 outcome and logical counter must remain unchanged. Profiling times and Build66/Build67 wall-clock ratios are informational and must not be used to promote Build67.'
    if [[ -n "$BASELINE_TSV" ]]; then echo "Build66 timing baseline: $BASELINE_TSV"; else echo 'Build66 timing baseline: unavailable'; fi
} > "$md"

rm -f "$OUTPUT_DIR/logs/"*.payload.bin

control_passes="$(awk -F '\t' 'NR>1 && $3=="control" && $37=="PASS" {n++} END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-front.jpg" && $37=="PASS" {n++} END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-mild.jpg" && $37=="PASS" {n++} END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-angle.jpg" && $37=="PASS" {n++} END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-front.jpg" && $37=="PASS" {n++} END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-mild.jpg" && $37=="PASS" {n++} END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-angle.jpg" && $37=="PASS" {n++} END{print n+0}' "$tsv")"

echo "Build67 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
echo "Build67 profiling matrix written to:"
echo "  $tsv"
echo "  $md"

if (( control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 || b_mild_pass != 1 || b_angle_pass != 1 )); then
    echo 'error: Build67 semantic-equivalence physical gate not met' >&2
    exit 1
fi
echo 'Build67 semantic-equivalence physical gate: PASS'
