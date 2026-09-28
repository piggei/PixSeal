#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD70_DIAGNOSTIC_DIR:-v4-phone private/build70-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD70_TIMEOUT:-86400}"
BUILD68_BASELINE_TSV="${V4_PHONE_BUILD68_BASELINE_TSV:-v4-phone private/build68-diagnostics/build68-phone-plane-reuse.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD70_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|required-build70|$MESSAGE_B"
    "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)

for spec in "${specs[@]}"; do
    IFS='|' read -r file _ <<< "$spec"
    [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }
done

mkdir -p "$OUTPUT_DIR/logs"
tsv="$OUTPUT_DIR/build70-phone-seed-genealogy.tsv"
md="$OUTPUT_DIR/build70-phone-seed-genealogy.md"
printf 'image\tclass\trole\tgeometry\tdata_decode\tbuild43_attempted\tbuild43_frozen\tbuild43_qualified\tbuild43_authenticated\tbuild64_attempted\tbuild64_seeds\tbuild64_geometry_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode_candidates\tbuild64_list_frames\tbuild64_authenticated\tbuild65_attempted\tbuild65_workers\tbuild66_attempted\tbuild66_workers\tbuild68_attempted\tbuild68_total_ms\tbuild68_geometry_ms\tbuild68_qualification_ms\tbuild68_decode_wall_ms\tbuild68_physical_decode_candidates\tbuild68_speculative_candidates\tbuild69_attempted\tbuild69_geometry_plane_prep_ms\tbuild69_geometry_freeze_ms\tbuild69_geometry_seed_wall_ms\tbuild69_geometry_seed_worker_ms\tbuild69_geometry_seed_min_ms\tbuild69_geometry_seed_median_ms\tbuild69_geometry_seed_max_ms\tbuild69_geometry_seed_min_evals\tbuild69_geometry_seed_max_evals\tbuild69_geometry_seed_min_bank\tbuild69_geometry_seed_max_bank\tbuild70_attempted\tbuild70_freeze_evals\tbuild70_seed_evals\tbuild70_max_seed_index\tbuild70_max_seed_pair_rank\tbuild70_max_seed_rank_within_pair\tbuild70_max_seed_evals\tbuild70_max_seed_bank\tbuild70_max_seed_ms\tbuild70_max_seed_stage_evals\tbuild70_max_seed_stage_states\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\tbuild68_elapsed_ms\tratio_vs_build68\n' > "$tsv"

extract_re() {
    local pattern="$1" file="$2" default_value="${3:-}"
    local value
    value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n 1)"
    if [[ -z "$value" ]]; then printf '%s' "$default_value"; else printf '%s' "$value"; fi
}

baseline68_elapsed_ms() {
    local image="$1"
    [[ -f "$BUILD68_BASELINE_TSV" ]] || { printf '%s' '-'; return; }
    awk -F '\t' -v image="$image" '
        NR==1 { for (i=1; i<=NF; i++) if ($i=="elapsed_ms") col=i; next }
        $1==image && col>0 { print $col; found=1; exit }
        END { if (!found) print "-" }
    ' "$BUILD68_BASELINE_TSV"
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

    b43_attempted=false; b43_frozen=0; b43_qualified=0; b43_auth=false
    if grep -q '^build43-geometry:' "$log"; then
        b43_attempted=true
        b43_frozen="$(extract_re 'build43-geometry: .* frozen=([0-9]+).*' "$log" 0)"
        b43_qualified="$(extract_re 'build43-geometry: .* qualified=([0-9]+).*' "$log" 0)"
        b43_auth="$(extract_re 'build43-geometry: .* authenticated=([^ ]+).*' "$log" false)"
    fi

    b64_attempted=false; b64_seeds=0; b64_evals=0; b64_bank=0; b64_qualified=0; b64_decode=0; b64_frames=0; b64_auth=false
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

    b65_attempted=false; b65_workers=0
    if grep -q '^build65-parallel:' "$log"; then
        b65_attempted=true
        b65_workers="$(extract_re 'build65-parallel: .* workers=([0-9]+).*' "$log" 0)"
    fi

    b66_attempted=false; b66_workers=0
    if grep -q '^build66-decode-parallel:' "$log"; then
        b66_attempted=true
        b66_workers="$(extract_re 'build66-decode-parallel: .* workers=([0-9]+).*' "$log" 0)"
    fi

    b68_attempted=false; b68_total_ms=0; b68_geometry_ms=0; b68_qualification_ms=0; b68_decode_wall_ms=0; b68_physical_decode=0; b68_speculative=0
    if grep -q '^build68-plane-reuse:' "$log"; then
        b68_attempted=true
        b68_total_ms="$(extract_re 'build68-plane-reuse: .* total-ms=([0-9]+).*' "$log" 0)"
        b68_geometry_ms="$(extract_re 'build68-plane-reuse: .* geometry-ms=([0-9]+).*' "$log" 0)"
        b68_qualification_ms="$(extract_re 'build68-plane-reuse: .* qualification-ms=([0-9]+).*' "$log" 0)"
        b68_decode_wall_ms="$(extract_re 'build68-plane-reuse: .* decode-wall-ms=([0-9]+).*' "$log" 0)"
        b68_physical_decode="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"
        b68_speculative="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)"
    fi

    b69_attempted=false; b69_plane=0; b69_freeze=0; b69_seed_wall=0; b69_seed_worker=0; b69_seed_min=0; b69_seed_median=0; b69_seed_max=0; b69_min_evals=0; b69_max_evals=0; b69_min_bank=0; b69_max_bank=0
    if grep -q '^build69-geometry-profile:' "$log"; then
        b69_attempted=true
        b69_plane="$(extract_re 'build69-geometry-profile: .* plane-prep-ms=([0-9]+).*' "$log" 0)"
        b69_freeze="$(extract_re 'build69-geometry-profile: .* freeze-ms=([0-9]+).*' "$log" 0)"
        b69_seed_wall="$(extract_re 'build69-geometry-profile: .* seed-wall-ms=([0-9]+).*' "$log" 0)"
        b69_seed_worker="$(extract_re 'build69-geometry-profile: .* seed-worker-ms=([0-9]+).*' "$log" 0)"
        b69_seed_min="$(extract_re 'build69-geometry-profile: .* seed-min-ms=([0-9]+).*' "$log" 0)"
        b69_seed_median="$(extract_re 'build69-geometry-profile: .* seed-median-ms=([0-9]+).*' "$log" 0)"
        b69_seed_max="$(extract_re 'build69-geometry-profile: .* seed-max-ms=([0-9]+).*' "$log" 0)"
        b69_min_evals="$(extract_re 'build69-geometry-profile: .* seed-min-evals=([0-9]+).*' "$log" 0)"
        b69_max_evals="$(extract_re 'build69-geometry-profile: .* seed-max-evals=([0-9]+).*' "$log" 0)"
        b69_min_bank="$(extract_re 'build69-geometry-profile: .* seed-min-bank=([0-9]+).*' "$log" 0)"
        b69_max_bank="$(extract_re 'build69-geometry-profile: .* seed-max-bank=([0-9]+).*' "$log" 0)"
    fi

    b70_attempted=false; b70_freeze_evals=0; b70_seed_evals=0; b70_max_seed_index=0; b70_max_seed_pair_rank=0; b70_max_seed_rank_pair=0; b70_max_seed_evals=0; b70_max_seed_bank=0; b70_max_seed_ms=0; b70_stage_evals='-'; b70_stage_states='-'
    if grep -q '^build70-seed-genealogy:' "$log"; then
        b70_attempted=true
        b70_freeze_evals="$(extract_re 'build70-seed-genealogy: .* freeze-evals=([0-9]+).*' "$log" 0)"
        b70_seed_evals="$(extract_re 'build70-seed-genealogy: .* seed-evals=([0-9]+).*' "$log" 0)"
        b70_max_seed_index="$(extract_re 'build70-seed-genealogy: .* max-seed-index=([0-9]+).*' "$log" 0)"
        b70_max_seed_pair_rank="$(extract_re 'build70-seed-genealogy: .* max-seed-pair-rank=([0-9]+).*' "$log" 0)"
        b70_max_seed_rank_pair="$(extract_re 'build70-seed-genealogy: .* max-seed-rank-within-pair=([0-9]+).*' "$log" 0)"
        b70_max_seed_evals="$(extract_re 'build70-seed-genealogy: .* max-seed-evals=([0-9]+).*' "$log" 0)"
        b70_max_seed_bank="$(extract_re 'build70-seed-genealogy: .* max-seed-bank=([0-9]+).*' "$log" 0)"
        b70_max_seed_ms="$(extract_re 'build70-seed-genealogy: .* max-seed-ms=([0-9]+).*' "$log" 0)"
        b70_stage_evals="$(extract_re 'build70-seed-genealogy: .* max-seed-stage-evals=([^ ]+).*' "$log" '-')"
        b70_stage_states="$(extract_re 'build70-seed-genealogy: .* max-seed-stage-states=([^ ]+).*' "$log" '-')"
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
    if (( exp_evals > 0 )); then
        if [[ "$b64_attempted" != true || "$b64_seeds" -ne 24 || "$b64_evals" -ne "$exp_evals" || "$b64_bank" -ne "$exp_bank" || "$b64_qualified" -ne "$exp_qualified" || "$b64_decode" -ne "$exp_decode" || "$b64_frames" -ne "$exp_frames" || "$b65_attempted" != true || "$b65_workers" -lt 1 || "$b66_attempted" != true || "$b68_attempted" != true || "$b69_attempted" != true || "$b70_attempted" != true ]]; then
            telemetry_equivalent=false
        fi
        if (( exp_qualified > 0 )); then
            if [[ "$b66_workers" -lt 1 || "$b68_physical_decode" -lt "$b64_decode" || "$b68_speculative" -ne $((b68_physical_decode - b64_decode)) ]]; then telemetry_equivalent=false; fi
        else
            if [[ "$b66_workers" -ne 0 || "$b68_physical_decode" -ne 0 || "$b68_speculative" -ne 0 || "$b68_decode_wall_ms" -ne 0 ]]; then telemetry_equivalent=false; fi
        fi
    else
        if [[ "$b64_attempted" != false || "$b65_attempted" != false || "$b66_attempted" != false || "$b68_attempted" != false || "$b69_attempted" != false || "$b70_attempted" != false ]]; then telemetry_equivalent=false; fi
    fi

    payload_match='-'; qualification='INFO'
    if [[ "$role" == control ]]; then
        if [[ $rc -ne 0 && "$hmac" == false && "$b64_attempted" == true && "$b64_auth" == false ]]; then qualification='PASS'; else qualification='FAIL'; fi
    else
        if [[ $rc -eq 0 ]]; then
            actual_payload="$(cat "$payload_file")"
            [[ "$actual_payload" == "$expected_payload" ]] && payload_match=true || payload_match=false
        else
            payload_match=false
        fi
        case "$role" in
            required-direct)  [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false && "$b64_attempted" == false ]] && qualification='PASS' || qualification='FAIL' ;;
            required-build43) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true && "$b64_attempted" == false ]] && qualification='PASS' || qualification='FAIL' ;;
            required-any)     [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false ]] && qualification='PASS' || qualification='FAIL' ;;
            required-build70) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b69_attempted" == true && "$b70_attempted" == true ]] && qualification='PASS' || qualification='FAIL' ;;
            required-reject)  [[ $rc -ne 0 && "$hmac" == false && "$b64_attempted" == true && "$b64_auth" == false && "$b69_attempted" == true && "$b70_attempted" == true ]] && qualification='PASS' || qualification='FAIL' ;;
        esac
    fi
    [[ "$telemetry_equivalent" == true ]] || qualification='FAIL'

    baseline_ms="$(baseline68_elapsed_ms "$file")"
    ratio='-'
    if [[ "$baseline_ms" =~ ^[0-9]+$ && "$elapsed_ms" -gt 0 ]]; then ratio="$(awk -v b="$baseline_ms" -v c="$elapsed_ms" 'BEGIN { printf "%.3f", b/c }')"; fi

    printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$file" "$class" "$role" "$geometry" "$data_decode" "$b43_attempted" "$b43_frozen" "$b43_qualified" "$b43_auth" "$b64_attempted" "$b64_seeds" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b65_attempted" "$b65_workers" "$b66_attempted" "$b66_workers" "$b68_attempted" "$b68_total_ms" "$b68_geometry_ms" "$b68_qualification_ms" "$b68_decode_wall_ms" "$b68_physical_decode" "$b68_speculative" "$b69_attempted" "$b69_plane" "$b69_freeze" "$b69_seed_wall" "$b69_seed_worker" "$b69_seed_min" "$b69_seed_median" "$b69_seed_max" "$b69_min_evals" "$b69_max_evals" "$b69_min_bank" "$b69_max_bank" "$b70_attempted" "$b70_freeze_evals" "$b70_seed_evals" "$b70_max_seed_index" "$b70_max_seed_pair_rank" "$b70_max_seed_rank_pair" "$b70_max_seed_evals" "$b70_max_seed_bank" "$b70_max_seed_ms" "$b70_stage_evals" "$b70_stage_states" "$hmac" "$payload_match" "$telemetry_equivalent" "$qualification" "$rc" "$elapsed_ms" "$baseline_ms" "$ratio" >> "$tsv"
done

{
    echo '# PixSeal Build70 dominant-seed genealogy profiling matrix'
    echo
    echo 'Build70 is an observability-only snapshot over the qualified Build68 smartphone baseline. It preserves the Build69 geometry bank/order and reports the identity plus per-stage genealogy of the most expensive seed.'
    echo
    echo '| image | role | evals | bank | qual | seed wall ms | max seed ms | max seed id (pair/rank) | max seed evals | max seed bank | stage evals | stage states | HMAC | telemetry eq | gate |'
    echo '|---|---|---:|---:|---:|---:|---:|---|---:|---:|---|---|---:|---:|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r file class role geometry data_decode b43_attempted b43_frozen b43_qualified b43_auth b64_attempted b64_seeds b64_evals b64_bank b64_qualified b64_decode b64_frames b64_auth b65_attempted b65_workers b66_attempted b66_workers b68_attempted b68_total_ms b68_geometry_ms b68_qualification_ms b68_decode_wall_ms b68_physical_decode b68_speculative b69_attempted b69_plane b69_freeze b69_seed_wall b69_seed_worker b69_seed_min b69_seed_median b69_seed_max b69_min_evals b69_max_evals b69_min_bank b69_max_bank b70_attempted b70_freeze_evals b70_seed_evals b70_max_seed_index b70_max_seed_pair_rank b70_max_seed_rank_pair b70_max_seed_evals b70_max_seed_bank b70_max_seed_ms b70_stage_evals b70_stage_states hmac payload_match telemetry_equivalent qualification rc elapsed_ms baseline_ms ratio; do
        echo "| $file | $role | $b64_evals | $b64_bank | $b64_qualified | $b69_seed_wall | $b70_max_seed_ms | $b70_max_seed_index ($b70_max_seed_pair_rank/$b70_max_seed_rank_pair) | $b70_max_seed_evals | $b70_max_seed_bank | $b70_stage_evals | $b70_stage_states | $hmac | $telemetry_equivalent | $qualification |"
    done
    echo
    echo 'Stage eval order: baseline,roots,single1,pair1,cont1,sib1,single2,pair2,cont2,sib2,single3,pair3,cont3,sib3,single4,pair4,cont4. Stage state order: roots,pair1,cont1,sib1,pair2,cont2,sib2,pair3,cont3,sib3,pair4,cont4,bank. Build70 is diagnostic only; do not change geometry semantics from these data alone.'
    if [[ -f "$BUILD68_BASELINE_TSV" ]]; then echo "Build68 timing baseline: $BUILD68_BASELINE_TSV"; else echo 'Build68 timing baseline: unavailable'; fi
} > "$md"

rm -f "$OUTPUT_DIR/logs/"*.payload.bin

count_pass() {
    local image="$1"
    awk -F '\t' -v image="$image" 'NR==1 {for(i=1;i<=NF;i++){if($i=="qualification")q=i}} NR>1 && $1==image && q>0 && $q=="PASS" {n++} END{print n+0}' "$tsv"
}
control_passes="$(awk -F '\t' 'NR==1 {for(i=1;i<=NF;i++){if($i=="role")r=i;if($i=="qualification")q=i}} NR>1 && $r=="control" && $q=="PASS" {n++} END{print n+0}' "$tsv")"
a_front_pass="$(count_pass phone-a-front.jpg)"; a_mild_pass="$(count_pass phone-a-mild.jpg)"; a_angle_pass="$(count_pass phone-a-angle.jpg)"
b_front_pass="$(count_pass phone-b-front.jpg)"; b_mild_pass="$(count_pass phone-b-mild.jpg)"; b_angle_pass="$(count_pass phone-b-angle.jpg)"

echo "Build70 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
echo "Build70 seed genealogy profile written to:"
echo "  $tsv"
echo "  $md"

if (( control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 || b_mild_pass != 1 || b_angle_pass != 1 )); then
    echo 'error: Build70 semantic-equivalence physical gate not met' >&2
    exit 1
fi
echo 'Build70 semantic-equivalence physical gate: PASS'
