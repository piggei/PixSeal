#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD74_DIAGNOSTIC_DIR:-v4-phone private/build74-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD74_TIMEOUT:-86400}"
BUILD73_BASELINE_TSV="${V4_PHONE_BUILD73_BASELINE_TSV:-v4-phone private/build73-diagnostics/build73-phone-gen3-parallel.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD74_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
    "phone-control-front.jpg|control|control|-"
    "phone-control-mild.jpg|control|control|-"
    "phone-control-angle.jpg|control|control|-"
    "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
    "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
    "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
    "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
    "phone-b-mild.jpg|B|required-build74|$MESSAGE_B"
    "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do IFS='|' read -r file _ <<< "$spec"; [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }; done

mkdir -p "$OUTPUT_DIR/logs"
tsv="$OUTPUT_DIR/build74-phone-freeze-profile.tsv"
md="$OUTPUT_DIR/build74-phone-freeze-profile.md"
printf 'image\tclass\trole\tbuild64_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode\tbuild64_frames\tbuild64_authenticated\tbuild68_physical_decode\tbuild68_speculative\tbuild73_attempted\tbuild74_attempted\tfreeze_total_ms\tstructural_ms\tfreeze_plane_ms\tpair_score_ms\tcells_ms\tbasin_ms\tproduction_basin_ms\tdepth_basin_ms\tallpairs_basin_ms\tpair_score_evals\tcell_evals\tbasin_evals\tproduction_basin_evals\tdepth_basin_evals\tallpairs_basin_evals\tpair_score_tasks\tcell_tasks\tbasin_tasks\tproduction_basin_tasks\tdepth_basin_tasks\tallpairs_basin_tasks\tbasin_min_ms\tbasin_median_ms\tbasin_max_ms\tbasin_min_evals\tbasin_max_evals\tbasin_min_outputs\tbasin_max_outputs\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\tbuild73_elapsed_ms\tgeometry_ms\tbuild73_geometry_ms\n' > "$tsv"

extract_re() { local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
baseline73_metric() {
    local image="$1" column="$2"
    [[ -f "$BUILD73_BASELINE_TSV" ]] || { printf '%s' '-'; return; }
    awk -F '\t' -v image="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==image&&col>0{print $col;found=1;exit}END{if(!found)print "-"}' "$BUILD73_BASELINE_TSV"
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
        b64_evals="$(extract_re 'build64-recovery: .* geometry-evals=([0-9]+).*' "$log" 0)"
        b64_bank="$(extract_re 'build64-recovery: .* bank=([0-9]+).*' "$log" 0)"
        b64_qualified="$(extract_re 'build64-recovery: .* qualified=([0-9]+).*' "$log" 0)"
        b64_decode="$(extract_re 'build64-recovery: .* decode-candidates=([0-9]+).*' "$log" 0)"
        b64_frames="$(extract_re 'build64-recovery: .* list-frames=([0-9]+).*' "$log" 0)"
        b64_auth="$(extract_re 'build64-recovery: .* authenticated=([^ ]+).*' "$log" false)"
    fi
    b65_attempted=false; grep -q '^build65-parallel:' "$log" && b65_attempted=true
    b66_attempted=false; b66_workers=0
    if grep -q '^build66-decode-parallel:' "$log"; then b66_attempted=true; b66_workers="$(extract_re 'build66-decode-parallel: .* workers=([0-9]+).*' "$log" 0)"; fi
    b68_attempted=false; b68_physical=0; b68_spec=0; geometry_ms=0
    if grep -q '^build68-plane-reuse:' "$log"; then
        b68_attempted=true
        b68_physical="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"
        b68_spec="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)"
        geometry_ms="$(extract_re 'build68-plane-reuse: .* geometry-ms=([0-9]+).*' "$log" 0)"
    fi
    b73_attempted=false; grep -q '^build73-gen3-parallel:' "$log" && b73_attempted=true

    b74_attempted=false
    total=0; structural=0; fplane=0; pairms=0; cellsms=0; basinms=0; prodms=0; depthms=0; allms=0
    paire=0; celle=0; basine=0; prode=0; depthe=0; alle=0
    pairt=0; cellt=0; basint=0; prodt=0; deptht=0; allt=0
    bmin=0; bmed=0; bmax=0; bmine=0; bmaxe=0; bmino=0; bmaxo=0
    if grep -q '^build74-freeze-profile:' "$log"; then
        b74_attempted=true
        total="$(extract_re 'build74-freeze-profile: .* total-ms=([0-9]+).*' "$log" 0)"
        structural="$(extract_re 'build74-freeze-profile: .* structural-ms=([0-9]+).*' "$log" 0)"
        fplane="$(extract_re 'build74-freeze-profile: .* plane-ms=([0-9]+).*' "$log" 0)"
        pairms="$(extract_re 'build74-freeze-profile: .* pair-ms=([0-9]+).*' "$log" 0)"
        cellsms="$(extract_re 'build74-freeze-profile: .* cells-ms=([0-9]+).*' "$log" 0)"
        basinms="$(extract_re 'build74-freeze-profile: .* basin-ms=([0-9]+).*' "$log" 0)"
        prodms="$(extract_re 'build74-freeze-profile: .* production-basin-ms=([0-9]+).*' "$log" 0)"
        depthms="$(extract_re 'build74-freeze-profile: .* depth-basin-ms=([0-9]+).*' "$log" 0)"
        allms="$(extract_re 'build74-freeze-profile: .* allpairs-basin-ms=([0-9]+).*' "$log" 0)"
        paire="$(extract_re 'build74-freeze-profile: .* pair-evals=([0-9]+).*' "$log" 0)"
        celle="$(extract_re 'build74-freeze-profile: .* cell-evals=([0-9]+).*' "$log" 0)"
        basine="$(extract_re 'build74-freeze-profile: .* basin-evals=([0-9]+).*' "$log" 0)"
        prode="$(extract_re 'build74-freeze-profile: .* production-evals=([0-9]+).*' "$log" 0)"
        depthe="$(extract_re 'build74-freeze-profile: .* depth-evals=([0-9]+).*' "$log" 0)"
        alle="$(extract_re 'build74-freeze-profile: .* allpairs-evals=([0-9]+).*' "$log" 0)"
        pairt="$(extract_re 'build74-freeze-profile: .* pair-tasks=([0-9]+).*' "$log" 0)"
        cellt="$(extract_re 'build74-freeze-profile: .* cell-tasks=([0-9]+).*' "$log" 0)"
        basint="$(extract_re 'build74-freeze-profile: .* basin-tasks=([0-9]+).*' "$log" 0)"
        prodt="$(extract_re 'build74-freeze-profile: .* production-tasks=([0-9]+).*' "$log" 0)"
        deptht="$(extract_re 'build74-freeze-profile: .* depth-tasks=([0-9]+).*' "$log" 0)"
        allt="$(extract_re 'build74-freeze-profile: .* allpairs-tasks=([0-9]+).*' "$log" 0)"
        bmin="$(extract_re 'build74-freeze-profile: .* basin-min-ms=([0-9]+).*' "$log" 0)"
        bmed="$(extract_re 'build74-freeze-profile: .* basin-median-ms=([0-9]+).*' "$log" 0)"
        bmax="$(extract_re 'build74-freeze-profile: .* basin-max-ms=([0-9]+).*' "$log" 0)"
        bmine="$(extract_re 'build74-freeze-profile: .* basin-min-evals=([0-9]+).*' "$log" 0)"
        bmaxe="$(extract_re 'build74-freeze-profile: .* basin-max-evals=([0-9]+).*' "$log" 0)"
        bmino="$(extract_re 'build74-freeze-profile: .* basin-min-outputs=([0-9]+).*' "$log" 0)"
        bmaxo="$(extract_re 'build74-freeze-profile: .* basin-max-outputs=([0-9]+).*' "$log" 0)"
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
        [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b65_attempted" == true && "$b66_attempted" == true && "$b68_attempted" == true && "$b73_attempted" == true && "$b74_attempted" == true && "$pairt" -gt 0 && "$basint" -gt 0 && "$total" -gt 0 ]] || telemetry=false
        [[ "$basint" -eq $((prodt+deptht+allt)) && "$basine" -eq $((prode+depthe+alle)) ]] || telemetry=false
        if [[ "$exp_qual" -gt 0 ]]; then
            [[ "$b66_workers" -gt 0 && "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]] || telemetry=false
            if [[ "$b64_auth" == true && "$b68_spec" -ge "$b66_workers" ]]; then telemetry=false; fi
            if [[ "$b64_auth" == false && "$b68_spec" -ne 0 ]]; then telemetry=false; fi
        else
            [[ "$b66_workers" -eq 0 && "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]] || telemetry=false
        fi
    else
        [[ "$b64_attempted" == false && "$b65_attempted" == false && "$b66_attempted" == false && "$b68_attempted" == false && "$b73_attempted" == false && "$b74_attempted" == false ]] || telemetry=false
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
            required-build74) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b74_attempted" == true ]] && qualification=PASS || qualification=FAIL ;;
            required-reject) [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b74_attempted" == true ]] && qualification=PASS || qualification=FAIL ;;
        esac
    fi
    [[ "$telemetry" == true ]] || qualification=FAIL

    base_elapsed="$(baseline73_metric "$file" elapsed_ms)"; base_geom="$(baseline73_metric "$file" geometry_ms)"
    printf '%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%d\t%d\t%s\t%d\t%s\n' \
      "$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b73_attempted" "$b74_attempted" "$total" "$structural" "$fplane" "$pairms" "$cellsms" "$basinms" "$prodms" "$depthms" "$allms" "$paire" "$celle" "$basine" "$prode" "$depthe" "$alle" "$pairt" "$cellt" "$basint" "$prodt" "$deptht" "$allt" "$bmin" "$bmed" "$bmax" "$bmine" "$bmaxe" "$bmino" "$bmaxo" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$base_elapsed" "$geometry_ms" "$base_geom" >> "$tsv"
done

{
 echo '# PixSeal Build74 Build47-freeze profile matrix'; echo
 echo 'Build74 is observability-only over the qualified Build73 baseline. It reproduces the exact Build47 frozen bank and Build73 scheduler while decomposing freeze time into structural seed, pair ranking, cell generation and basin generation. Fine-grained timings are diagnostic and Build74 is not promotable.'; echo
 echo '| image | role | evals | bank | qual | decode | frames | freeze ms | structural | pair | cells | basin | prod basin | depth basin | all-pairs basin | basin tasks | basin evals | HMAC | telemetry eq | elapsed ms | B73 ms | gate |'
 echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
 tail -n +2 "$tsv" | awk -F '\t' '{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$1,$3,$4,$5,$6,$7,$8,$14,$15,$17,$18,$19,$20,$21,$22,$31,$25,$42,$44,$47,$48,$45}'
 echo; echo 'Correctness is authoritative. Build73 remains the qualified baseline; Build74 timing exists only to select a later optimization.'
 [[ -f "$BUILD73_BASELINE_TSV" ]] && echo "Build73 timing reference: $BUILD73_BASELINE_TSV" || echo 'Build73 timing reference: unavailable'
} > "$md"
rm -f "$OUTPUT_DIR/logs/"*.payload.bin

pass_count="$(awk -F '\t' 'NR>1&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
control_passes="$(awk -F '\t' 'NR>1&&$3=="control"&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-front.jpg"&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-mild.jpg"&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-angle.jpg"&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-front.jpg"&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-mild.jpg"&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-angle.jpg"&&$45=="PASS"{n++}END{print n+0}' "$tsv")"
echo "Build74 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
echo "Build74 freeze profile written to:"; echo "  $tsv"; echo "  $md"
if (( pass_count != 9 || control_passes != 3 || a_front_pass != 1 || a_mild_pass != 1 || a_angle_pass != 1 || b_front_pass != 1 || b_mild_pass != 1 || b_angle_pass != 1 )); then echo 'error: Build74 semantic-equivalence physical gate not met' >&2; exit 1; fi
echo 'Build74 semantic-equivalence physical gate: PASS'
