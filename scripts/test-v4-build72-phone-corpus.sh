#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD72_DIAGNOSTIC_DIR:-v4-phone private/build72-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD72_TIMEOUT:-86400}"
BUILD71_BASELINE_TSV="${V4_PHONE_BUILD71_BASELINE_TSV:-v4-phone private/build71-diagnostics/build71-phone-gen4-parallel.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD72_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|required-build72|$MESSAGE_B"
    "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)

for spec in "${specs[@]}"; do
    IFS='|' read -r file _ <<< "$spec"
    [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }
done

mkdir -p "$OUTPUT_DIR/logs"
tsv="$OUTPUT_DIR/build72-phone-prefix-profile.tsv"
md="$OUTPUT_DIR/build72-phone-prefix-profile.md"
printf 'image\tclass\trole\tbuild64_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode\tbuild64_frames\tbuild64_authenticated\tbuild71_attempted\tbuild71_prefix_workers\tbuild71_gen4_workers\tbuild71_gen4_tasks\tbuild71_prefix_wall_ms\tbuild71_gen4_wall_ms\tbuild72_attempted\tbuild72_prefix_min_ms\tbuild72_prefix_median_ms\tbuild72_prefix_max_ms\tbuild72_prefix_min_evals\tbuild72_prefix_max_evals\tbuild72_prefix_min_inputs\tbuild72_prefix_max_inputs\tbuild72_prefix_stage_evals\tbuild72_prefix_stage_worker_ms\tbuild72_max_prefix_seed_index\tbuild72_max_prefix_seed_pair_rank\tbuild72_max_prefix_seed_rank_within_pair\tbuild72_max_prefix_seed_evals\tbuild72_max_prefix_seed_inputs\tbuild72_max_prefix_seed_ms\tbuild72_max_prefix_seed_stage_evals\tbuild72_max_prefix_seed_stage_states\tbuild72_max_prefix_seed_stage_ms\thmac\tpayload_match\ttelemetry_equivalent\texit_code\telapsed_ms\tbuild71_elapsed_ms\tbuild71_prefix_baseline_ms\tqualification\n' > "$tsv"

extract_re() {
    local pattern="$1" file="$2" default_value="${3:-}"
    local value
    value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n 1)"
    if [[ -z "$value" ]]; then printf '%s' "$default_value"; else printf '%s' "$value"; fi
}

baseline71_metric() {
    local image="$1" column="$2"
    [[ -f "$BUILD71_BASELINE_TSV" ]] || { printf '%s' '-'; return; }
    awk -F '\t' -v image="$image" -v wanted="$column" '
        NR==1 { for (i=1; i<=NF; i++) if ($i==wanted) col=i; next }
        $1==image && col>0 { print $col; found=1; exit }
        END { if (!found) print "-" }
    ' "$BUILD71_BASELINE_TSV"
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

    hmac="$(extract_re 'hmac: authenticated=([^ ]+).*' "$log" false)"

    b64_attempted=false; b64_evals=0; b64_bank=0; b64_qualified=0; b64_decode=0; b64_frames=0; b64_auth=false
    if grep -q '^build64-recovery:' "$log"; then
        b64_attempted=true
        b64_evals="$(extract_re 'build64-recovery: .* geometry-evals=([0-9]+).*' "$log" 0)"
        b64_bank="$(extract_re 'build64-recovery: .* bank=([0-9]+).*' "$log" 0)"
        b64_qualified="$(extract_re 'build64-recovery: .* qualified=([0-9]+).*' "$log" 0)"
        b64_decode="$(extract_re 'build64-recovery: .* decode-candidates=([0-9]+).*' "$log" 0)"
        b64_frames="$(extract_re 'build64-recovery: .* list-frames=([0-9]+).*' "$log" 0)"
        b64_auth="$(extract_re 'build64-recovery: .* authenticated=([^ ]+).*' "$log" false)"
    fi

    b71_attempted=false; b71_prefix_workers=0; b71_gen4_workers=0; b71_gen4_tasks=0; b71_prefix_wall=0; b71_gen4_wall=0
    if grep -q '^build71-gen4-parallel:' "$log"; then
        b71_attempted=true
        b71_prefix_workers="$(extract_re 'build71-gen4-parallel: .* prefix-workers=([0-9]+).*' "$log" 0)"
        b71_gen4_workers="$(extract_re 'build71-gen4-parallel: .* gen4-workers=([0-9]+).*' "$log" 0)"
        b71_gen4_tasks="$(extract_re 'build71-gen4-parallel: .* gen4-tasks=([0-9]+).*' "$log" 0)"
        b71_prefix_wall="$(extract_re 'build71-gen4-parallel: .* prefix-wall-ms=([0-9]+).*' "$log" 0)"
        b71_gen4_wall="$(extract_re 'build71-gen4-parallel: .* gen4-wall-ms=([0-9]+).*' "$log" 0)"
    fi

    b72_attempted=false; b72_min_ms=0; b72_median_ms=0; b72_max_ms=0; b72_min_evals=0; b72_max_evals=0; b72_min_inputs=0; b72_max_inputs=0; b72_stage_evals='-'; b72_stage_worker_ms='-'; b72_seed_index=0; b72_pair_rank=0; b72_rank_within=0; b72_seed_evals=0; b72_seed_inputs=0; b72_seed_ms=0; b72_seed_stage_evals='-'; b72_seed_stage_states='-'; b72_seed_stage_ms='-'
    if grep -q '^build72-prefix-profile:' "$log"; then
        b72_attempted=true
        b72_min_ms="$(extract_re 'build72-prefix-profile: .* prefix-min-ms=([0-9]+).*' "$log" 0)"
        b72_median_ms="$(extract_re 'build72-prefix-profile: .* prefix-median-ms=([0-9]+).*' "$log" 0)"
        b72_max_ms="$(extract_re 'build72-prefix-profile: .* prefix-max-ms=([0-9]+).*' "$log" 0)"
        b72_min_evals="$(extract_re 'build72-prefix-profile: .* prefix-min-evals=([0-9]+).*' "$log" 0)"
        b72_max_evals="$(extract_re 'build72-prefix-profile: .* prefix-max-evals=([0-9]+).*' "$log" 0)"
        b72_min_inputs="$(extract_re 'build72-prefix-profile: .* prefix-min-inputs=([0-9]+).*' "$log" 0)"
        b72_max_inputs="$(extract_re 'build72-prefix-profile: .* prefix-max-inputs=([0-9]+).*' "$log" 0)"
        b72_stage_evals="$(extract_re 'build72-prefix-profile: .* prefix-stage-evals=([^ ]+).*' "$log" '-')"
        b72_stage_worker_ms="$(extract_re 'build72-prefix-profile: .* prefix-stage-worker-ms=([^ ]+).*' "$log" '-')"
        b72_seed_index="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-index=([0-9]+).*' "$log" 0)"
        b72_pair_rank="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-pair-rank=([0-9]+).*' "$log" 0)"
        b72_rank_within="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-rank-within-pair=([0-9]+).*' "$log" 0)"
        b72_seed_evals="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-evals=([0-9]+).*' "$log" 0)"
        b72_seed_inputs="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-inputs=([0-9]+).*' "$log" 0)"
        b72_seed_ms="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-ms=([0-9]+).*' "$log" 0)"
        b72_seed_stage_evals="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-stage-evals=([^ ]+).*' "$log" '-')"
        b72_seed_stage_states="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-stage-states=([^ ]+).*' "$log" '-')"
        b72_seed_stage_ms="$(extract_re 'build72-prefix-profile: .* max-prefix-seed-stage-ms=([^ ]+).*' "$log" '-')"
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
        if [[ "$b64_attempted" != true || "$b64_evals" -ne "$exp_evals" || "$b64_bank" -ne "$exp_bank" || "$b64_qualified" -ne "$exp_qualified" || "$b64_decode" -ne "$exp_decode" || "$b64_frames" -ne "$exp_frames" || "$b71_attempted" != true || "$b72_attempted" != true || "$b71_prefix_workers" -lt 1 || "$b71_gen4_workers" -lt 1 || "$b71_gen4_tasks" -lt 1 || "$b72_stage_evals" == "-" || "$b72_stage_worker_ms" == "-" || "$b72_seed_stage_evals" == "-" || "$b72_seed_stage_ms" == "-" ]]; then
            telemetry_equivalent=false
        fi
    else
        if [[ "$b64_attempted" != false || "$b71_attempted" != false || "$b72_attempted" != false ]]; then telemetry_equivalent=false; fi
    fi

    payload_match='-'; qualification='INFO'
    if [[ "$role" == control ]]; then
        [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b72_attempted" == true ]] && qualification='PASS' || qualification='FAIL'
    else
        if [[ $rc -eq 0 ]]; then actual_payload="$(cat "$payload_file")"; [[ "$actual_payload" == "$expected_payload" ]] && payload_match=true || payload_match=false; else payload_match=false; fi
        case "$role" in
            required-direct) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false && "$b72_attempted" == false ]] && qualification='PASS' || qualification='FAIL' ;;
            required-build43) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false && "$b72_attempted" == false ]] && qualification='PASS' || qualification='FAIL' ;;
            required-any) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false && "$b72_attempted" == false ]] && qualification='PASS' || qualification='FAIL' ;;
            required-build72) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b72_attempted" == true ]] && qualification='PASS' || qualification='FAIL' ;;
            required-reject) [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b72_attempted" == true ]] && qualification='PASS' || qualification='FAIL' ;;
        esac
    fi
    [[ "$telemetry_equivalent" == true ]] || qualification='FAIL'

    base_elapsed="$(baseline71_metric "$file" elapsed_ms)"
    base_prefix="$(baseline71_metric "$file" build71_prefix_wall_ms)"

    printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b71_attempted" "$b71_prefix_workers" "$b71_gen4_workers" "$b71_gen4_tasks" "$b71_prefix_wall" "$b71_gen4_wall" "$b72_attempted" "$b72_min_ms" "$b72_median_ms" "$b72_max_ms" "$b72_min_evals" "$b72_max_evals" "$b72_min_inputs" "$b72_max_inputs" "$b72_stage_evals" "$b72_stage_worker_ms" "$b72_seed_index" "$b72_pair_rank" "$b72_rank_within" "$b72_seed_evals" "$b72_seed_inputs" "$b72_seed_ms" "$b72_seed_stage_evals" "$b72_seed_stage_states" "$b72_seed_stage_ms" "$hmac" "$payload_match" "$telemetry_equivalent" "$rc" "$elapsed_ms" "$base_elapsed" "$base_prefix" "$qualification" >> "$tsv"
done

{
    echo '# PixSeal Build72 prefix-stage profiling matrix'
    echo
    echo 'Build72 is observability-only over the qualified Build71 smartphone baseline. It keeps the exact Build71 two-barrier scheduler and profiles only the public proposal-only prefix through sibling3.'
    echo
    echo 'Prefix stage order for the CSV fields: baseline, roots, single1, pair1, cont1, sib1, single2, pair2, cont2, sib2, single3, pair3, cont3, sib3.'
    echo
    echo '| image | role | evals | bank | qual | logical decode | logical frames | prefix wall ms | gen4 wall ms | prefix min/med/max ms | dominant prefix seed | pair/rank | seed evals | seed inputs | seed ms | HMAC | telemetry eq | elapsed ms | B71 ms | B71 prefix ms | gate |'
    echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image class role b64_evals b64_bank b64_qualified b64_decode b64_frames b64_auth b71_attempted b71_prefix_workers b71_gen4_workers b71_gen4_tasks b71_prefix_wall b71_gen4_wall b72_attempted b72_min_ms b72_median_ms b72_max_ms b72_min_evals b72_max_evals b72_min_inputs b72_max_inputs b72_stage_evals b72_stage_worker_ms b72_seed_index b72_pair_rank b72_rank_within b72_seed_evals b72_seed_inputs b72_seed_ms b72_seed_stage_evals b72_seed_stage_states b72_seed_stage_ms hmac payload_match telemetry_equivalent rc elapsed_ms base_elapsed base_prefix qualification; do
        echo "| $image | $role | $b64_evals | $b64_bank | $b64_qualified | $b64_decode | $b64_frames | $b71_prefix_wall | $b71_gen4_wall | $b72_min_ms/$b72_median_ms/$b72_max_ms | $b72_seed_index | $b72_pair_rank/$b72_rank_within | $b72_seed_evals | $b72_seed_inputs | $b72_seed_ms | $hmac | $telemetry_equivalent | $elapsed_ms | $base_elapsed | $base_prefix | $qualification |"
    done
    echo
    echo 'Correctness is authoritative. Build72 timing is diagnostic and must not be used as a promotion benchmark because fine-grained timers add overhead.'
    [[ -f "$BUILD71_BASELINE_TSV" ]] && echo "Build71 timing reference: $BUILD71_BASELINE_TSV" || echo 'Build71 timing reference: unavailable'
} > "$md"

rm -f "$OUTPUT_DIR/logs/"*.payload.bin

pass_count="$(awk -F '\t' 'NR>1 && $42=="PASS" {n++} END{print n+0}' "$tsv")"
control_passes="$(awk -F '\t' 'NR>1 && $3=="control" && $42=="PASS" {n++} END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-front.jpg" && $42=="PASS" {n++} END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-mild.jpg" && $42=="PASS" {n++} END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1 && $1=="phone-a-angle.jpg" && $42=="PASS" {n++} END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-front.jpg" && $42=="PASS" {n++} END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-mild.jpg" && $42=="PASS" {n++} END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR>1 && $1=="phone-b-angle.jpg" && $42=="PASS" {n++} END{print n+0}' "$tsv")"

echo "Build72 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
echo "Build72 prefix-stage profile written to:"
echo "  $tsv"
echo "  $md"
if (( pass_count != 9 || control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 || b_mild_pass != 1 || b_angle_pass != 1 )); then
    echo 'error: Build72 semantic-equivalence physical gate not met' >&2
    exit 1
fi
echo 'Build72 semantic-equivalence physical gate: PASS'
