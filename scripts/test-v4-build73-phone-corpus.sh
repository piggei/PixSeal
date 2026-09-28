#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD73_DIAGNOSTIC_DIR:-v4-phone private/build73-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD73_TIMEOUT:-86400}"
BUILD71_BASELINE_TSV="${V4_PHONE_BUILD71_BASELINE_TSV:-v4-phone private/build71-diagnostics/build71-phone-gen4-parallel.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD73_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|required-build73|$MESSAGE_B"
    "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do IFS='|' read -r file _ <<< "$spec"; [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }; done

mkdir -p "$OUTPUT_DIR/logs"
tsv="$OUTPUT_DIR/build73-phone-gen3-parallel.tsv"
md="$OUTPUT_DIR/build73-phone-gen3-parallel.md"
printf 'image\tclass\trole\tbuild64_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode\tbuild64_frames\tbuild64_authenticated\tbuild68_physical_decode\tbuild68_speculative\tbuild73_attempted\tbuild73_prefix2_workers\tbuild73_gen3_workers\tbuild73_gen4_workers\tbuild73_gen3_tasks\tbuild73_gen4_tasks\tbuild73_plane_prep_ms\tbuild73_freeze_ms\tbuild73_prefix2_wall_ms\tbuild73_prefix2_worker_ms\tbuild73_gen3_wall_ms\tbuild73_gen3_worker_ms\tbuild73_gen4_wall_ms\tbuild73_gen4_worker_ms\tbuild73_gen3_min_ms\tbuild73_gen3_median_ms\tbuild73_gen3_max_ms\tbuild73_gen3_min_evals\tbuild73_gen3_max_evals\tbuild73_gen3_min_outputs\tbuild73_gen3_max_outputs\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\tbuild71_elapsed_ms\tspeedup_vs_build71\tgeometry_ms\tbuild71_geometry_ms\tgeometry_speedup_vs_build71\n' > "$tsv"

extract_re() { local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
baseline71_metric() {
    local image="$1" column="$2"
    [[ -f "$BUILD71_BASELINE_TSV" ]] || { printf '%s' '-'; return; }
    awk -F '\t' -v image="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==image&&col>0{print $col;found=1;exit}END{if(!found)print "-"}' "$BUILD71_BASELINE_TSV"
}

for spec in "${specs[@]}"; do
    IFS='|' read -r file class role expected_payload <<< "$spec"
    input="$ACQUISITION_DIR/$file"; stem="${file%.*}"; log="$OUTPUT_DIR/logs/$stem.stderr.txt"; payload_file="$OUTPUT_DIR/logs/$stem.payload.bin"
    start_ns="$(date +%s%N)"; set +e
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"
    rc=$?; set -e; end_ns="$(date +%s%N)"; elapsed_ms=$(( (end_ns-start_ns)/1000000 ))
    hmac="$(extract_re 'hmac: authenticated=([^ ]+).*' "$log" false)"
    b43_attempted=false; b43_auth=false
    if grep -q '^build43-geometry:' "$log"; then
        b43_attempted=true
        b43_auth="$(extract_re 'build43-geometry: .* authenticated=([^ ]+).*' "$log" false)"
    fi

    b64_attempted=false; b64_evals=0; b64_bank=0; b64_qualified=0; b64_decode=0; b64_frames=0; b64_auth=false
    if grep -q '^build64-recovery:' "$log"; then
        b64_attempted=true
        b64_evals="$(extract_re 'build64-recovery: .* geometry-evals=([0-9]+).*' "$log" 0)"; b64_bank="$(extract_re 'build64-recovery: .* bank=([0-9]+).*' "$log" 0)"; b64_qualified="$(extract_re 'build64-recovery: .* qualified=([0-9]+).*' "$log" 0)"; b64_decode="$(extract_re 'build64-recovery: .* decode-candidates=([0-9]+).*' "$log" 0)"; b64_frames="$(extract_re 'build64-recovery: .* list-frames=([0-9]+).*' "$log" 0)"; b64_auth="$(extract_re 'build64-recovery: .* authenticated=([^ ]+).*' "$log" false)"
    fi
    b65_attempted=false
    grep -q '^build65-parallel:' "$log" && b65_attempted=true
    b66_attempted=false; b66_workers=0
    if grep -q '^build66-decode-parallel:' "$log"; then
        b66_attempted=true
        b66_workers="$(extract_re 'build66-decode-parallel: .* workers=([0-9]+).*' "$log" 0)"
    fi
    b68_attempted=false; b68_physical=0; b68_spec=0; geometry_ms=0
    if grep -q '^build68-plane-reuse:' "$log"; then
        b68_attempted=true
        b68_physical="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"; b68_spec="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)"; geometry_ms="$(extract_re 'build68-plane-reuse: .* geometry-ms=([0-9]+).*' "$log" 0)"
    fi

    b73_attempted=false; p2w=0; g3w=0; g4w=0; g3tasks=0; g4tasks=0; plane=0; freeze=0; p2wall=0; p2worker=0; g3wall=0; g3worker=0; g4wall=0; g4worker=0; g3min=0; g3med=0; g3max=0; g3mine=0; g3maxe=0; g3mino=0; g3maxo=0
    if grep -q '^build73-gen3-parallel:' "$log"; then
        b73_attempted=true
        p2w="$(extract_re 'build73-gen3-parallel: .* prefix2-workers=([0-9]+).*' "$log" 0)"; g3w="$(extract_re 'build73-gen3-parallel: .* gen3-workers=([0-9]+).*' "$log" 0)"; g4w="$(extract_re 'build73-gen3-parallel: .* gen4-workers=([0-9]+).*' "$log" 0)"; g3tasks="$(extract_re 'build73-gen3-parallel: .* gen3-tasks=([0-9]+).*' "$log" 0)"; g4tasks="$(extract_re 'build73-gen3-parallel: .* gen4-tasks=([0-9]+).*' "$log" 0)"
        plane="$(extract_re 'build73-gen3-parallel: .* plane-prep-ms=([0-9]+).*' "$log" 0)"; freeze="$(extract_re 'build73-gen3-parallel: .* freeze-ms=([0-9]+).*' "$log" 0)"; p2wall="$(extract_re 'build73-gen3-parallel: .* prefix2-wall-ms=([0-9]+).*' "$log" 0)"; p2worker="$(extract_re 'build73-gen3-parallel: .* prefix2-worker-ms=([0-9]+).*' "$log" 0)"; g3wall="$(extract_re 'build73-gen3-parallel: .* gen3-wall-ms=([0-9]+).*' "$log" 0)"; g3worker="$(extract_re 'build73-gen3-parallel: .* gen3-worker-ms=([0-9]+).*' "$log" 0)"; g4wall="$(extract_re 'build73-gen3-parallel: .* gen4-wall-ms=([0-9]+).*' "$log" 0)"; g4worker="$(extract_re 'build73-gen3-parallel: .* gen4-worker-ms=([0-9]+).*' "$log" 0)"
        g3min="$(extract_re 'build73-gen3-parallel: .* gen3-min-ms=([0-9]+).*' "$log" 0)"; g3med="$(extract_re 'build73-gen3-parallel: .* gen3-median-ms=([0-9]+).*' "$log" 0)"; g3max="$(extract_re 'build73-gen3-parallel: .* gen3-max-ms=([0-9]+).*' "$log" 0)"; g3mine="$(extract_re 'build73-gen3-parallel: .* gen3-min-evals=([0-9]+).*' "$log" 0)"; g3maxe="$(extract_re 'build73-gen3-parallel: .* gen3-max-evals=([0-9]+).*' "$log" 0)"; g3mino="$(extract_re 'build73-gen3-parallel: .* gen3-min-outputs=([0-9]+).*' "$log" 0)"; g3maxo="$(extract_re 'build73-gen3-parallel: .* gen3-max-outputs=([0-9]+).*' "$log" 0)"
    fi

    case "$file" in
        phone-control-front.jpg) exp_evals=18021; exp_bank=26; exp_qual=0; exp_decode=0; exp_frames=0 ;;
        phone-control-mild.jpg) exp_evals=117609; exp_bank=810; exp_qual=6; exp_decode=6; exp_frames=18432 ;;
        phone-control-angle.jpg) exp_evals=21203; exp_bank=11; exp_qual=0; exp_decode=0; exp_frames=0 ;;
        phone-b-mild.jpg) exp_evals=79259; exp_bank=937; exp_qual=935; exp_decode=691; exp_frames=2120047 ;;
        phone-b-angle.jpg) exp_evals=334857; exp_bank=6198; exp_qual=0; exp_decode=0; exp_frames=0 ;;
        *) exp_evals=0; exp_bank=0; exp_qual=0; exp_decode=0; exp_frames=0 ;;
    esac
    telemetry=true
    if [[ "$exp_evals" -gt 0 ]]; then
        [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b65_attempted" == true && "$b66_attempted" == true && "$b68_attempted" == true && "$b73_attempted" == true && "$p2w" -gt 0 && "$g3w" -gt 0 && "$g4w" -gt 0 && "$g3tasks" -gt 0 && "$g4tasks" -gt 0 ]] || telemetry=false
        if [[ "$exp_qual" -gt 0 ]]; then
            [[ "$b66_workers" -gt 0 && "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]] || telemetry=false
            if [[ "$b64_auth" == true && "$b68_spec" -ge "$b66_workers" ]]; then telemetry=false; fi
            if [[ "$b64_auth" == false && "$b68_spec" -ne 0 ]]; then telemetry=false; fi
        else
            [[ "$b66_workers" -eq 0 && "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]] || telemetry=false
        fi
    else
        [[ "$b64_attempted" == false && "$b65_attempted" == false && "$b66_attempted" == false && "$b68_attempted" == false && "$b73_attempted" == false ]] || telemetry=false
    fi

    payload_match='-'; qualification='INFO'
    if [[ "$role" == control ]]; then [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false ]] && qualification=PASS || qualification=FAIL
    else
        if [[ $rc -eq 0 ]]; then actual="$(cat "$payload_file")"; [[ "$actual" == "$expected_payload" ]] && payload_match=true || payload_match=false; else payload_match=false; fi
        case "$role" in
            required-direct) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL ;;
            required-build43) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL ;;
            required-any) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL ;;
            required-build73) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b73_attempted" == true ]] && qualification=PASS || qualification=FAIL ;;
            required-reject) [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b73_attempted" == true ]] && qualification=PASS || qualification=FAIL ;;
        esac
    fi
    [[ "$telemetry" == true ]] || qualification=FAIL

    base_elapsed="$(baseline71_metric "$file" elapsed_ms)"; base_geom="$(baseline71_metric "$file" build68_geometry_ms)"; speed='-'; gspeed='-'
    if [[ "$base_elapsed" =~ ^[0-9]+$ && "$elapsed_ms" -gt 0 ]]; then speed="$(awk -v b="$base_elapsed" -v c="$elapsed_ms" 'BEGIN{printf "%.3f",b/c}')"; fi
    if [[ "$base_geom" =~ ^[0-9]+$ && "$geometry_ms" -gt 0 ]]; then gspeed="$(awk -v b="$base_geom" -v c="$geometry_ms" 'BEGIN{printf "%.3f",b/c}')"; fi
    printf '%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%d\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%d\t%d\t%s\t%s\t%d\t%s\t%s\n' \
      "$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b73_attempted" "$p2w" "$g3w" "$g4w" "$g3tasks" "$g4tasks" "$plane" "$freeze" "$p2wall" "$p2worker" "$g3wall" "$g3worker" "$g4wall" "$g4worker" "$g3min" "$g3med" "$g3max" "$g3mine" "$g3maxe" "$g3mino" "$g3maxo" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$base_elapsed" "$speed" "$geometry_ms" "$base_geom" "$gspeed" >> "$tsv"
done

{
 echo '# PixSeal Build73 ordered generation-three parallel matrix'; echo
 echo 'Build73 is an equivalence-preserving performance candidate over the qualified Build71 baseline. It freezes the exact prefix through sibling2, runs independent generation-three subtrees in one bounded global pool, commits strictly in original order, then reuses the qualified Build71 generation-four pool.'; echo
 echo '| image | role | evals | bank | qual | logical decode | frames | gen3 tasks | gen4 tasks | prefix2 wall ms | gen3 wall ms | gen4 wall ms | geometry ms | HMAC | telemetry eq | elapsed ms | B71 ms | speedup | B71 geom ms | geom speedup | gate |'
 echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
 tail -n +2 "$tsv" | awk -F '\t' '{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$1,$3,$4,$5,$6,$7,$8,$16,$17,$20,$22,$24,$41,$33,$35,$38,$39,$40,$42,$43,$36}'
 echo; echo 'Correctness is authoritative. Build64/66/68 logical counters and HMAC/payload outcomes must remain unchanged. Build71 remains qualified until Build73 performance is reproduced.'
 [[ -f "$BUILD71_BASELINE_TSV" ]] && echo "Build71 timing baseline: $BUILD71_BASELINE_TSV" || echo 'Build71 timing baseline: unavailable'
} > "$md"
rm -f "$OUTPUT_DIR/logs/"*.payload.bin

pass_count="$(awk -F '\t' 'NR>1&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
control_passes="$(awk -F '\t' 'NR>1&&$3=="control"&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-front.jpg"&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-mild.jpg"&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-angle.jpg"&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-front.jpg"&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-mild.jpg"&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-angle.jpg"&&$36=="PASS"{n++}END{print n+0}' "$tsv")"
echo "Build73 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
echo "Build73 ordered generation-three matrix written to:"; echo "  $tsv"; echo "  $md"
if (( pass_count != 9 || control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 || b_mild_pass != 1 || b_angle_pass != 1 )); then echo 'error: Build73 semantic-equivalence physical gate not met' >&2; exit 1; fi
echo 'Build73 semantic-equivalence physical gate: PASS'
