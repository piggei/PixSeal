#!/usr/bin/env bash
set -euo pipefail
PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD76_DIAGNOSTIC_DIR:-v4-phone private/build76-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD76_TIMEOUT:-86400}"
BUILD75_BASELINE_TSV="${V4_PHONE_BUILD75_BASELINE_TSV:-v4-phone private/build75-diagnostics/build75-phone-basin-parallel.tsv}"
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD76_TIMEOUT must be an integer number of seconds" >&2; exit 1; }
specs=(
 "phone-control-front.jpg|control|control|-"
 "phone-control-mild.jpg|control|control|-"
 "phone-control-angle.jpg|control|control|-"
 "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
 "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
 "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
 "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
 "phone-b-mild.jpg|B|required-build76|$MESSAGE_B"
 "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do IFS='|' read -r file _ <<< "$spec"; [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }; done
output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build76.XXXXXX")"
published=false
cleanup_stage(){
 if [[ "$published" != true && -d "$stage_dir" ]]; then rm -rf "$stage_dir"; fi
}
trap cleanup_stage EXIT
mkdir -p "$stage_dir/logs"
tsv="$stage_dir/build76-phone-gen2-parallel.tsv"; md="$stage_dir/build76-phone-gen2-parallel.md"
final_tsv="$OUTPUT_DIR/build76-phone-gen2-parallel.tsv"; final_md="$OUTPUT_DIR/build76-phone-gen2-parallel.md"
printf 'image\tclass\trole\tbuild64_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode\tbuild64_frames\tbuild64_authenticated\tbuild68_physical_decode\tbuild68_speculative\tbuild75_attempted\tbuild76_attempted\tprefix1_workers\tgen2_workers\tgen3_workers\tgen4_workers\tgen2_tasks\tgen3_tasks\tgen4_tasks\tfreeze_ms\tprefix1_wall_ms\tprefix1_worker_ms\tgen2_wall_ms\tgen2_worker_ms\tgen3_wall_ms\tgen3_worker_ms\tgen4_wall_ms\tgen4_worker_ms\tgen2_min_ms\tgen2_median_ms\tgen2_max_ms\tgen2_min_evals\tgen2_max_evals\tgen2_min_outputs\tgen2_max_outputs\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\tbuild75_elapsed_ms\tspeedup_vs_build75\tgeometry_ms\tbuild75_geometry_ms\tgeometry_speedup_vs_build75\n' > "$tsv"
extract_re(){ local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
baseline_metric(){ local path="$1" image="$2" column="$3"; [[ -f "$path" ]] || { printf '%s' '-'; return; }; awk -F '\t' -v image="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==image&&col>0{print $col;found=1;exit}END{if(!found)print "-"}' "$path"; }
speedup(){ local base="$1" now="$2"; if [[ "$base" =~ ^[0-9]+$ && "$now" =~ ^[0-9]+$ && "$now" -gt 0 ]]; then awk -v b="$base" -v n="$now" 'BEGIN{printf "%.3f",b/n}'; else printf '%s' '-'; fi; }
for spec in "${specs[@]}"; do
 IFS='|' read -r file class role expected_payload <<< "$spec"; input="$ACQUISITION_DIR/$file"; stem="${file%.*}"; log="$stage_dir/logs/$stem.stderr.txt"; payload_file="$stage_dir/logs/$stem.payload.bin"
 start_ns="$(date +%s%N)"; set +e; timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"; rc=$?; set -e; end_ns="$(date +%s%N)"; elapsed_ms=$(( (end_ns-start_ns)/1000000 ))
 hmac="$(extract_re 'hmac: authenticated=([^ ]+).*' "$log" false)"; b43_attempted=false; b43_auth=false; if grep -q '^build43-geometry:' "$log"; then b43_attempted=true; b43_auth="$(extract_re 'build43-geometry: .* authenticated=([^ ]+).*' "$log" false)"; fi
 b64_attempted=false;b64_evals=0;b64_bank=0;b64_qualified=0;b64_decode=0;b64_frames=0;b64_auth=false
 if grep -q '^build64-recovery:' "$log"; then b64_attempted=true; b64_evals="$(extract_re 'build64-recovery: .* geometry-evals=([0-9]+).*' "$log" 0)"; b64_bank="$(extract_re 'build64-recovery: .* bank=([0-9]+).*' "$log" 0)"; b64_qualified="$(extract_re 'build64-recovery: .* qualified=([0-9]+).*' "$log" 0)"; b64_decode="$(extract_re 'build64-recovery: .* decode-candidates=([0-9]+).*' "$log" 0)"; b64_frames="$(extract_re 'build64-recovery: .* list-frames=([0-9]+).*' "$log" 0)"; b64_auth="$(extract_re 'build64-recovery: .* authenticated=([^ ]+).*' "$log" false)"; fi
 b65=false; grep -q '^build65-parallel:' "$log" && b65=true; b66=false;b66_workers=0; if grep -q '^build66-decode-parallel:' "$log"; then b66=true;b66_workers="$(extract_re 'build66-decode-parallel: .* workers=([0-9]+).*' "$log" 0)"; fi
 b68=false;b68_physical=0;b68_spec=0;geometry_ms=0;if grep -q '^build68-plane-reuse:' "$log";then b68=true;b68_physical="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)";b68_spec="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)";geometry_ms="$(extract_re 'build68-plane-reuse: .* geometry-ms=([0-9]+).*' "$log" 0)";fi
 b73=false;grep -q '^build73-gen3-parallel:' "$log"&&b73=true;b75=false;grep -q '^build75-basin-parallel:' "$log"&&b75=true
 b76=false;p1w=0;g2w=0;g3w=0;g4w=0;g2t=0;g3t=0;g4t=0;freeze=0;p1wall=0;p1worker=0;g2wall=0;g2worker=0;g3wall=0;g3worker=0;g4wall=0;g4worker=0;g2min=0;g2med=0;g2max=0;g2mine=0;g2maxe=0;g2mino=0;g2maxo=0
 if grep -q '^build76-gen2-parallel:' "$log";then b76=true;p1w="$(extract_re 'build76-gen2-parallel: .* prefix1-workers=([0-9]+).*' "$log" 0)";g2w="$(extract_re 'build76-gen2-parallel: .* gen2-workers=([0-9]+).*' "$log" 0)";g3w="$(extract_re 'build76-gen2-parallel: .* gen3-workers=([0-9]+).*' "$log" 0)";g4w="$(extract_re 'build76-gen2-parallel: .* gen4-workers=([0-9]+).*' "$log" 0)";g2t="$(extract_re 'build76-gen2-parallel: .* gen2-tasks=([0-9]+).*' "$log" 0)";g3t="$(extract_re 'build76-gen2-parallel: .* gen3-tasks=([0-9]+).*' "$log" 0)";g4t="$(extract_re 'build76-gen2-parallel: .* gen4-tasks=([0-9]+).*' "$log" 0)";freeze="$(extract_re 'build76-gen2-parallel: .* freeze-ms=([0-9]+).*' "$log" 0)";p1wall="$(extract_re 'build76-gen2-parallel: .* prefix1-wall-ms=([0-9]+).*' "$log" 0)";p1worker="$(extract_re 'build76-gen2-parallel: .* prefix1-worker-ms=([0-9]+).*' "$log" 0)";g2wall="$(extract_re 'build76-gen2-parallel: .* gen2-wall-ms=([0-9]+).*' "$log" 0)";g2worker="$(extract_re 'build76-gen2-parallel: .* gen2-worker-ms=([0-9]+).*' "$log" 0)";g3wall="$(extract_re 'build76-gen2-parallel: .* gen3-wall-ms=([0-9]+).*' "$log" 0)";g3worker="$(extract_re 'build76-gen2-parallel: .* gen3-worker-ms=([0-9]+).*' "$log" 0)";g4wall="$(extract_re 'build76-gen2-parallel: .* gen4-wall-ms=([0-9]+).*' "$log" 0)";g4worker="$(extract_re 'build76-gen2-parallel: .* gen4-worker-ms=([0-9]+).*' "$log" 0)";g2min="$(extract_re 'build76-gen2-parallel: .* gen2-min-ms=([0-9]+).*' "$log" 0)";g2med="$(extract_re 'build76-gen2-parallel: .* gen2-median-ms=([0-9]+).*' "$log" 0)";g2max="$(extract_re 'build76-gen2-parallel: .* gen2-max-ms=([0-9]+).*' "$log" 0)";g2mine="$(extract_re 'build76-gen2-parallel: .* gen2-min-evals=([0-9]+).*' "$log" 0)";g2maxe="$(extract_re 'build76-gen2-parallel: .* gen2-max-evals=([0-9]+).*' "$log" 0)";g2mino="$(extract_re 'build76-gen2-parallel: .* gen2-min-outputs=([0-9]+).*' "$log" 0)";g2maxo="$(extract_re 'build76-gen2-parallel: .* gen2-max-outputs=([0-9]+).*' "$log" 0)";fi
 case "$file" in phone-control-front.jpg) exp_evals=18021;exp_bank=26;exp_qual=0;exp_decode=0;exp_frames=0;; phone-control-mild.jpg) exp_evals=117609;exp_bank=810;exp_qual=6;exp_decode=6;exp_frames=18432;; phone-control-angle.jpg) exp_evals=21203;exp_bank=11;exp_qual=0;exp_decode=0;exp_frames=0;; phone-b-mild.jpg) exp_evals=79259;exp_bank=937;exp_qual=935;exp_decode=691;exp_frames=2120047;; phone-b-angle.jpg) exp_evals=334857;exp_bank=6198;exp_qual=0;exp_decode=0;exp_frames=0;; *) exp_evals=0;exp_bank=0;exp_qual=0;exp_decode=0;exp_frames=0;; esac
 telemetry=true
 if [[ "$exp_evals" -gt 0 ]];then [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b65" == true && "$b66" == true && "$b68" == true && "$b73" == true && "$b75" == true && "$b76" == true && "$p1w" -gt 0 && "$g2w" -gt 0 && "$g2t" -gt 0 && "$freeze" -gt 0 && "$g2wall" -gt 0 && "$g2worker" -ge "$g2wall" ]]||telemetry=false; if [[ "$exp_qual" -gt 0 ]];then [[ "$b66_workers" -gt 0 && "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]]||telemetry=false; [[ "$b64_auth" == false || "$b68_spec" -lt "$b66_workers" ]]||telemetry=false;else [[ "$b66_workers" -eq 0 && "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]]||telemetry=false;fi; else [[ "$b64_attempted" == false && "$b73" == false && "$b75" == false && "$b76" == false ]]||telemetry=false;fi
 payload_match='-';qualification=INFO
 if [[ "$role" == control ]];then [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false ]]&&qualification=PASS||qualification=FAIL;else if [[ $rc -eq 0 ]];then actual="$(cat "$payload_file")";[[ "$actual" == "$expected_payload" ]]&&payload_match=true||payload_match=false;else payload_match=false;fi;case "$role" in required-direct) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false && "$b64_attempted" == false ]]&&qualification=PASS||qualification=FAIL;; required-build43) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true && "$b64_attempted" == false ]]&&qualification=PASS||qualification=FAIL;; required-any) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false ]]&&qualification=PASS||qualification=FAIL;; required-build76) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b76" == true ]]&&qualification=PASS||qualification=FAIL;; required-reject) [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b76" == true ]]&&qualification=PASS||qualification=FAIL;; esac;fi
 [[ "$telemetry" == true ]]||qualification=FAIL
 base_elapsed="$(baseline_metric "$BUILD75_BASELINE_TSV" "$file" elapsed_ms)";base_geom="$(baseline_metric "$BUILD75_BASELINE_TSV" "$file" geometry_ms)";speed="$(speedup "$base_elapsed" "$elapsed_ms")";gspeed="$(speedup "$base_geom" "$geometry_ms")"
 printf '%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%d\t%d\t%s\t%s\t%d\t%s\t%s\n' "$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b75" "$b76" "$p1w" "$g2w" "$g3w" "$g4w" "$g2t" "$g3t" "$g4t" "$freeze" "$p1wall" "$p1worker" "$g2wall" "$g2worker" "$g3wall" "$g3worker" "$g4wall" "$g4worker" "$g2min" "$g2med" "$g2max" "$g2mine" "$g2maxe" "$g2mino" "$g2maxo" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$base_elapsed" "$speed" "$geometry_ms" "$base_geom" "$gspeed" >> "$tsv"
done
{
 echo '# PixSeal Build76 ordered generation-two parallel matrix';echo
 echo 'Build76 is the qualified equivalence-preserving smartphone baseline, originally measured against qualified Build75. It moves only the public geometry scheduling barrier: exact per-seed search through sibling1, ordered global generation-two pool, then unchanged qualified generation-three/generation-four pools. Build75 ordered-parallel freeze and all downstream qualification/data/HMAC semantics are unchanged.';echo
 echo '| image | role | evals | bank | qual | decode | frames | gen2 tasks | prefix1 ms | gen2 ms | gen2 worker ms | gen3 ms | gen4 ms | geometry ms | B75 geom | geom speedup | HMAC | telemetry eq | elapsed ms | B75 ms | speedup | gate |'
 echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
 tail -n +2 "$tsv" | awk -F '\t' '{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$1,$3,$4,$5,$6,$7,$8,$18,$22,$24,$25,$26,$28,$45,$46,$47,$37,$39,$42,$43,$44,$40}'
 echo;echo 'Correctness is authoritative. This target is now a Build76 qualification/reproducibility gate; Build75 timings are retained only as the historical comparison baseline.'; [[ -f "$BUILD75_BASELINE_TSV" ]]&&echo "Build75 timing baseline: $BUILD75_BASELINE_TSV"||echo 'Build75 timing baseline: unavailable'
} > "$md"
rm -f "$stage_dir/logs/"*.payload.bin
pass_count="$(awk -F '\t' 'NR>1&&$40=="PASS"{n++}END{print n+0}' "$tsv")";control_passes="$(awk -F '\t' 'NR>1&&$3=="control"&&$40=="PASS"{n++}END{print n+0}' "$tsv")";a_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-front.jpg"&&$40=="PASS"{n++}END{print n+0}' "$tsv")";a_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-mild.jpg"&&$40=="PASS"{n++}END{print n+0}' "$tsv")";a_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-angle.jpg"&&$40=="PASS"{n++}END{print n+0}' "$tsv")";b_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-front.jpg"&&$40=="PASS"{n++}END{print n+0}' "$tsv")";b_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-mild.jpg"&&$40=="PASS"{n++}END{print n+0}' "$tsv")";b_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-angle.jpg"&&$40=="PASS"{n++}END{print n+0}' "$tsv")"
echo "Build76 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
if ((pass_count!=9||control_passes!=3||a_front_pass!=1||a_mild_pass!=1||a_angle_pass!=1||b_front_pass!=1||b_mild_pass!=1||b_angle_pass!=1));then echo 'error: Build76 semantic-equivalence physical gate not met' >&2;exit 1;fi
mkdir -p "$OUTPUT_DIR"
rm -rf "$OUTPUT_DIR/logs"
mv "$stage_dir/logs" "$OUTPUT_DIR/logs"
mv "$tsv" "$final_tsv"
mv "$md" "$final_md"
rmdir "$stage_dir"
published=true
echo 'Build76 ordered generation-two matrix written to:';echo "  $final_tsv";echo "  $final_md"
echo 'Build76 semantic-equivalence physical gate: PASS'
