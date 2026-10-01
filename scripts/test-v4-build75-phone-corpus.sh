#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD75_DIAGNOSTIC_DIR:-v4-phone private/build75-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD75_TIMEOUT:-86400}"
BUILD73_BASELINE_TSV="${V4_PHONE_BUILD73_BASELINE_TSV:-v4-phone private/build73-diagnostics/build73-phone-gen3-parallel.tsv}"
BUILD74_BASELINE_TSV="${V4_PHONE_BUILD74_BASELINE_TSV:-v4-phone private/build74-diagnostics/build74-phone-freeze-profile.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD75_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|required-build75|$MESSAGE_B"
    "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do IFS='|' read -r file _ <<< "$spec"; [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }; done

mkdir -p "$OUTPUT_DIR/logs"
tsv="$OUTPUT_DIR/build75-phone-basin-parallel.tsv"
md="$OUTPUT_DIR/build75-phone-basin-parallel.md"
printf 'image\tclass\trole\tbuild64_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode\tbuild64_frames\tbuild64_authenticated\tbuild68_physical_decode\tbuild68_speculative\tbuild73_attempted\tbuild75_attempted\tbasin_workers\tbasin_tasks\tfreeze_total_ms\tstructural_ms\tfreeze_plane_ms\tpair_score_ms\tcells_ms\tbasin_wall_ms\tbasin_worker_ms\tproduction_worker_ms\tdepth_worker_ms\tallpairs_worker_ms\tbasin_evals\tproduction_evals\tdepth_evals\tallpairs_evals\tproduction_tasks\tdepth_tasks\tallpairs_tasks\tbasin_min_ms\tbasin_median_ms\tbasin_max_ms\tbasin_min_evals\tbasin_max_evals\tbasin_min_outputs\tbasin_max_outputs\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\tbuild73_elapsed_ms\tspeedup_vs_build73\tgeometry_ms\tbuild73_geometry_ms\tgeometry_speedup_vs_build73\tbuild74_freeze_ms\tfreeze_speedup_vs_build74\tbuild74_basin_ms\tbasin_speedup_vs_build74\n' > "$tsv"

extract_re() { local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
baseline_metric() {
    local path="$1" image="$2" column="$3"
    [[ -f "$path" ]] || { printf '%s' '-'; return; }
    awk -F '\t' -v image="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==image&&col>0{print $col;found=1;exit}END{if(!found)print "-"}' "$path"
}
speedup() {
    local base="$1" now="$2"
    if [[ "$base" =~ ^[0-9]+$ && "$now" =~ ^[0-9]+$ && "$now" -gt 0 ]]; then awk -v b="$base" -v n="$now" 'BEGIN{printf "%.3f",b/n}'; else printf '%s' '-'; fi
}

for spec in "${specs[@]}"; do
    IFS='|' read -r file class role expected_payload <<< "$spec"
    input="$ACQUISITION_DIR/$file"; stem="${file%.*}"; log="$OUTPUT_DIR/logs/$stem.stderr.txt"; payload_file="$OUTPUT_DIR/logs/$stem.payload.bin"
    start_ns="$(date +%s%N)"; set +e
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"
    rc=$?; set -e; end_ns="$(date +%s%N)"; elapsed_ms=$(( (end_ns-start_ns)/1000000 ))
    hmac="$(extract_re 'hmac: authenticated=([^ ]+).*' "$log" false)"
    b43_attempted=false; b43_auth=false
    if grep -q '^build43-geometry:' "$log"; then b43_attempted=true; b43_auth="$(extract_re 'build43-geometry: .* authenticated=([^ ]+).*' "$log" false)"; fi

    b64_attempted=false; b64_evals=0; b64_bank=0; b64_qualified=0; b64_decode=0; b64_frames=0; b64_auth=false
    if grep -q '^build64-recovery:' "$log"; then
        b64_attempted=true
        b64_evals="$(extract_re 'build64-recovery: .* geometry-evals=([0-9]+).*' "$log" 0)"; b64_bank="$(extract_re 'build64-recovery: .* bank=([0-9]+).*' "$log" 0)"; b64_qualified="$(extract_re 'build64-recovery: .* qualified=([0-9]+).*' "$log" 0)"; b64_decode="$(extract_re 'build64-recovery: .* decode-candidates=([0-9]+).*' "$log" 0)"; b64_frames="$(extract_re 'build64-recovery: .* list-frames=([0-9]+).*' "$log" 0)"; b64_auth="$(extract_re 'build64-recovery: .* authenticated=([^ ]+).*' "$log" false)"
    fi
    b65_attempted=false; grep -q '^build65-parallel:' "$log" && b65_attempted=true
    b66_attempted=false; b66_workers=0
    if grep -q '^build66-decode-parallel:' "$log"; then b66_attempted=true; b66_workers="$(extract_re 'build66-decode-parallel: .* workers=([0-9]+).*' "$log" 0)"; fi
    b68_attempted=false; b68_physical=0; b68_spec=0; geometry_ms=0
    if grep -q '^build68-plane-reuse:' "$log"; then
        b68_attempted=true
        b68_physical="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"; b68_spec="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)"; geometry_ms="$(extract_re 'build68-plane-reuse: .* geometry-ms=([0-9]+).*' "$log" 0)"
    fi
    b73_attempted=false; grep -q '^build73-gen3-parallel:' "$log" && b73_attempted=true

    b75_attempted=false; bw=0; bt=0; total=0; structural=0; fplane=0; pairms=0; cellsms=0; bwall=0; bworker=0; prodw=0; depthw=0; allw=0; be=0; prode=0; depthe=0; alle=0; prodt=0; deptht=0; allt=0; bmin=0; bmed=0; bmax=0; bmine=0; bmaxe=0; bmino=0; bmaxo=0
    if grep -q '^build75-basin-parallel:' "$log"; then
        b75_attempted=true
        total="$(extract_re 'build75-basin-parallel: .* freeze-total-ms=([0-9]+).*' "$log" 0)"; structural="$(extract_re 'build75-basin-parallel: .* structural-ms=([0-9]+).*' "$log" 0)"; fplane="$(extract_re 'build75-basin-parallel: .* plane-ms=([0-9]+).*' "$log" 0)"; pairms="$(extract_re 'build75-basin-parallel: .* pair-ms=([0-9]+).*' "$log" 0)"; cellsms="$(extract_re 'build75-basin-parallel: .* cells-ms=([0-9]+).*' "$log" 0)"
        bw="$(extract_re 'build75-basin-parallel: .* basin-workers=([0-9]+).*' "$log" 0)"; bt="$(extract_re 'build75-basin-parallel: .* basin-tasks=([0-9]+).*' "$log" 0)"; bwall="$(extract_re 'build75-basin-parallel: .* basin-wall-ms=([0-9]+).*' "$log" 0)"; bworker="$(extract_re 'build75-basin-parallel: .* basin-worker-ms=([0-9]+).*' "$log" 0)"
        prodw="$(extract_re 'build75-basin-parallel: .* production-worker-ms=([0-9]+).*' "$log" 0)"; depthw="$(extract_re 'build75-basin-parallel: .* depth-worker-ms=([0-9]+).*' "$log" 0)"; allw="$(extract_re 'build75-basin-parallel: .* allpairs-worker-ms=([0-9]+).*' "$log" 0)"
        be="$(extract_re 'build75-basin-parallel: .* basin-evals=([0-9]+).*' "$log" 0)"; prode="$(extract_re 'build75-basin-parallel: .* production-evals=([0-9]+).*' "$log" 0)"; depthe="$(extract_re 'build75-basin-parallel: .* depth-evals=([0-9]+).*' "$log" 0)"; alle="$(extract_re 'build75-basin-parallel: .* allpairs-evals=([0-9]+).*' "$log" 0)"
        prodt="$(extract_re 'build75-basin-parallel: .* production-tasks=([0-9]+).*' "$log" 0)"; deptht="$(extract_re 'build75-basin-parallel: .* depth-tasks=([0-9]+).*' "$log" 0)"; allt="$(extract_re 'build75-basin-parallel: .* allpairs-tasks=([0-9]+).*' "$log" 0)"
        bmin="$(extract_re 'build75-basin-parallel: .* basin-min-ms=([0-9]+).*' "$log" 0)"; bmed="$(extract_re 'build75-basin-parallel: .* basin-median-ms=([0-9]+).*' "$log" 0)"; bmax="$(extract_re 'build75-basin-parallel: .* basin-max-ms=([0-9]+).*' "$log" 0)"; bmine="$(extract_re 'build75-basin-parallel: .* basin-min-evals=([0-9]+).*' "$log" 0)"; bmaxe="$(extract_re 'build75-basin-parallel: .* basin-max-evals=([0-9]+).*' "$log" 0)"; bmino="$(extract_re 'build75-basin-parallel: .* basin-min-outputs=([0-9]+).*' "$log" 0)"; bmaxo="$(extract_re 'build75-basin-parallel: .* basin-max-outputs=([0-9]+).*' "$log" 0)"
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
        [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b65_attempted" == true && "$b66_attempted" == true && "$b68_attempted" == true && "$b73_attempted" == true && "$b75_attempted" == true && "$bw" -gt 0 && "$bt" -gt 0 && "$total" -gt 0 && "$bwall" -gt 0 && "$bworker" -gt 0 ]] || telemetry=false
        [[ "$bt" -eq $((prodt+deptht+allt)) && "$be" -eq $((prode+depthe+alle)) && "$bworker" -ge "$bwall" ]] || telemetry=false
        if [[ "$exp_qual" -gt 0 ]]; then
            [[ "$b66_workers" -gt 0 && "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]] || telemetry=false
            if [[ "$b64_auth" == true && "$b68_spec" -ge "$b66_workers" ]]; then telemetry=false; fi
            if [[ "$b64_auth" == false && "$b68_spec" -ne 0 ]]; then telemetry=false; fi
        else
            [[ "$b66_workers" -eq 0 && "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]] || telemetry=false
        fi
    else
        [[ "$b64_attempted" == false && "$b65_attempted" == false && "$b66_attempted" == false && "$b68_attempted" == false && "$b73_attempted" == false && "$b75_attempted" == false ]] || telemetry=false
    fi

    payload_match='-'; qualification='INFO'
    if [[ "$role" == control ]]; then
        [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false ]] && qualification=PASS || qualification=FAIL
    else
        if [[ $rc -eq 0 ]]; then actual="$(cat "$payload_file")"; [[ "$actual" == "$expected_payload" ]] && payload_match=true || payload_match=false; else payload_match=false; fi
        case "$role" in
            required-direct) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL ;;
            required-build43) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL ;;
            required-any) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL ;;
            required-build75) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b75_attempted" == true ]] && qualification=PASS || qualification=FAIL ;;
            required-reject) [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b75_attempted" == true ]] && qualification=PASS || qualification=FAIL ;;
        esac
    fi
    [[ "$telemetry" == true ]] || qualification=FAIL

    base_elapsed="$(baseline_metric "$BUILD73_BASELINE_TSV" "$file" elapsed_ms)"; base_geom="$(baseline_metric "$BUILD73_BASELINE_TSV" "$file" geometry_ms)"
    base_freeze="$(baseline_metric "$BUILD74_BASELINE_TSV" "$file" freeze_total_ms)"; base_basin="$(baseline_metric "$BUILD74_BASELINE_TSV" "$file" basin_ms)"
    speed="$(speedup "$base_elapsed" "$elapsed_ms")"; gspeed="$(speedup "$base_geom" "$geometry_ms")"; fspeed="$(speedup "$base_freeze" "$total")"; bspeed="$(speedup "$base_basin" "$bwall")"
    printf '%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%d\t%d\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n' \
      "$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b73_attempted" "$b75_attempted" "$bw" "$bt" "$total" "$structural" "$fplane" "$pairms" "$cellsms" "$bwall" "$bworker" "$prodw" "$depthw" "$allw" "$be" "$prode" "$depthe" "$alle" "$prodt" "$deptht" "$allt" "$bmin" "$bmed" "$bmax" "$bmine" "$bmaxe" "$bmino" "$bmaxo" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$base_elapsed" "$speed" "$geometry_ms" "$base_geom" "$gspeed" "$base_freeze" "$fspeed" "$base_basin" "$bspeed" >> "$tsv"
done

{
 echo '# PixSeal Build75 ordered-parallel Build47 basin matrix'; echo
 echo 'Build75 is an equivalence-preserving performance candidate over the qualified Build73 baseline. It freezes the historical Build47 basin task stream, computes independent basin tasks with one bounded pool, and commits results strictly in original tier/pair/cell order. No pruning, score, threshold, bank-order or protected-data semantics change.'; echo
 echo '| image | role | evals | bank | qual | decode | frames | basin workers | tasks | freeze ms | basin wall ms | basin worker ms | geometry ms | B73 geom | geom speedup | HMAC | telemetry eq | elapsed ms | B73 ms | speedup | B74 freeze | freeze speedup | B74 basin | basin speedup | gate |'
 echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
 tail -n +2 "$tsv" | awk -F '\t' '{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$1,$3,$4,$5,$6,$7,$8,$14,$15,$16,$21,$22,$48,$49,$50,$40,$42,$45,$46,$47,$51,$52,$53,$54,$43}'
 echo; echo 'Correctness is authoritative. Build73 remains the qualified baseline until Build75 performance is reproduced independently.'
 [[ -f "$BUILD73_BASELINE_TSV" ]] && echo "Build73 timing baseline: $BUILD73_BASELINE_TSV" || echo 'Build73 timing baseline: unavailable'
 [[ -f "$BUILD74_BASELINE_TSV" ]] && echo "Build74 freeze reference: $BUILD74_BASELINE_TSV" || echo 'Build74 freeze reference: unavailable'
} > "$md"
rm -f "$OUTPUT_DIR/logs/"*.payload.bin

pass_count="$(awk -F '\t' 'NR>1&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
control_passes="$(awk -F '\t' 'NR>1&&$3=="control"&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-front.jpg"&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-mild.jpg"&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-angle.jpg"&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-front.jpg"&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-mild.jpg"&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-angle.jpg"&&$43=="PASS"{n++}END{print n+0}' "$tsv")"
echo "Build75 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
echo "Build75 ordered-parallel basin matrix written to:"; echo "  $tsv"; echo "  $md"
if (( pass_count != 9 || control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 || b_mild_pass != 1 || b_angle_pass != 1 )); then echo 'error: Build75 semantic-equivalence physical gate not met' >&2; exit 1; fi
echo 'Build75 semantic-equivalence physical gate: PASS'
