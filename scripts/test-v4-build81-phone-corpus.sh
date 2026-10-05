#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD81_DIAGNOSTIC_DIR:-v4-phone private/build81-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD81_TIMEOUT:-86400}"
BUILD76_BASELINE_TSV="${V4_PHONE_BUILD76_BASELINE_TSV:-v4-phone private/build76-diagnostics/build76-phone-gen2-parallel.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD81_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
 "phone-control-front.jpg|control|control|-"
 "phone-control-mild.jpg|control|control|-"
 "phone-control-angle.jpg|control|control|-"
 "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
 "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
 "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
 "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
 "phone-b-mild.jpg|B|required-build81|$MESSAGE_B"
 "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do
  IFS='|' read -r file _ <<< "$spec"
  [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }
done

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build81.XXXXXX")"
published=false
cleanup_stage(){ if [[ "$published" != true && -d "$stage_dir" ]]; then rm -rf "$stage_dir"; fi; }
trap cleanup_stage EXIT
mkdir -p "$stage_dir/logs"
tsv="$stage_dir/build81-phone-luminance-lut.tsv"
md="$stage_dir/build81-phone-luminance-lut.md"
final_tsv="$OUTPUT_DIR/build81-phone-luminance-lut.tsv"
final_md="$OUTPUT_DIR/build81-phone-luminance-lut.md"
failed_output_dir="${OUTPUT_DIR}-failed"

header=(image class role build64_evals build64_bank build64_qualified build64_decode build64_frames build64_authenticated build68_physical_decode build68_speculative build75_attempted build76_attempted build81_attempted prefix1_workers gen2_workers gen3_workers gen4_workers gen2_tasks gen3_tasks gen4_tasks freeze_ms prefix1_ms gen2_ms gen3_ms gen4_ms lut_fold_scores lut_block_reads lut_block_success lut_block_failed hmac payload_match telemetry_equivalent qualification exit_code elapsed_ms build76_elapsed_ms speedup_vs_build76 geometry_ms build76_geometry_ms geometry_speedup_vs_build76 gate_reason)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"

extract_re(){ local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
baseline_metric(){ local path="$1" image="$2" column="$3"; [[ -f "$path" ]] || { printf '%s' '-'; return; }; awk -F '\t' -v image="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==image&&col>0{print $col;found=1;exit}END{if(!found)print "-"}' "$path"; }
speedup(){ local base="$1" now="$2"; if [[ "$base" =~ ^[0-9]+$ && "$now" =~ ^[0-9]+$ && "$now" -gt 0 ]]; then awk -v b="$base" -v n="$now" 'BEGIN{printf "%.3f",b/n}'; else printf '%s' '-'; fi; }

for spec in "${specs[@]}"; do
  IFS='|' read -r file class role expected_payload <<< "$spec"
  input="$ACQUISITION_DIR/$file"; stem="${file%.*}"
  log="$stage_dir/logs/$stem.stderr.txt"; payload_file="$stage_dir/logs/$stem.payload.bin"
  start_ns="$(date +%s%N)"
  set +e
  timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"
  rc=$?
  set -e
  end_ns="$(date +%s%N)"; elapsed_ms=$(( (end_ns-start_ns)/1000000 ))

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
  b65=false; grep -q '^build65-parallel:' "$log" && b65=true
  b66=false; b66_workers=0
  if grep -q '^build66-decode-parallel:' "$log"; then b66=true; b66_workers="$(extract_re 'build66-decode-parallel: .* workers=([0-9]+).*' "$log" 0)"; fi
  b68=false; b68_physical=0; b68_spec=0; geometry_ms=0
  if grep -q '^build68-plane-reuse:' "$log"; then
    b68=true
    b68_physical="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"
    b68_spec="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)"
    geometry_ms="$(extract_re 'build68-plane-reuse: .* geometry-ms=([0-9]+).*' "$log" 0)"
  fi
  b73=false; grep -q '^build73-gen3-parallel:' "$log" && b73=true
  b75=false; grep -q '^build75-basin-parallel:' "$log" && b75=true
  b76=false; p1w=0; g2w=0; g3w=0; g4w=0; g2t=0; g3t=0; g4t=0; freeze=0; p1wall=0; g2wall=0; g3wall=0; g4wall=0
  if grep -q '^build76-gen2-parallel:' "$log"; then
    b76=true
    p1w="$(extract_re 'build76-gen2-parallel: .* prefix1-workers=([0-9]+).*' "$log" 0)"
    g2w="$(extract_re 'build76-gen2-parallel: .* gen2-workers=([0-9]+).*' "$log" 0)"
    g3w="$(extract_re 'build76-gen2-parallel: .* gen3-workers=([0-9]+).*' "$log" 0)"
    g4w="$(extract_re 'build76-gen2-parallel: .* gen4-workers=([0-9]+).*' "$log" 0)"
    g2t="$(extract_re 'build76-gen2-parallel: .* gen2-tasks=([0-9]+).*' "$log" 0)"
    g3t="$(extract_re 'build76-gen2-parallel: .* gen3-tasks=([0-9]+).*' "$log" 0)"
    g4t="$(extract_re 'build76-gen2-parallel: .* gen4-tasks=([0-9]+).*' "$log" 0)"
    freeze="$(extract_re 'build76-gen2-parallel: .* freeze-ms=([0-9]+).*' "$log" 0)"
    p1wall="$(extract_re 'build76-gen2-parallel: .* prefix1-wall-ms=([0-9]+).*' "$log" 0)"
    g2wall="$(extract_re 'build76-gen2-parallel: .* gen2-wall-ms=([0-9]+).*' "$log" 0)"
    g3wall="$(extract_re 'build76-gen2-parallel: .* gen3-wall-ms=([0-9]+).*' "$log" 0)"
    g4wall="$(extract_re 'build76-gen2-parallel: .* gen4-wall-ms=([0-9]+).*' "$log" 0)"
  fi
  b81=false; lut_folds=0; lut_reads=0; lut_success=0; lut_failed=0
  if grep -q '^build81-luminance-lut:' "$log"; then
    b81=true
    lut_folds="$(extract_re 'build81-luminance-lut: .* fold-scores=([0-9]+).*' "$log" 0)"
    lut_reads="$(extract_re 'build81-luminance-lut: .* block-reads=([0-9]+).*' "$log" 0)"
    lut_success="$(extract_re 'build81-luminance-lut: .* block-success=([0-9]+).*' "$log" 0)"
    lut_failed="$(extract_re 'build81-luminance-lut: .* block-failed=([0-9]+).*' "$log" 0)"
  fi

  case "$file" in
    phone-control-front.jpg) exp_evals=18021; exp_bank=26; exp_qual=0; exp_decode=0; exp_frames=0;;
    phone-control-mild.jpg) exp_evals=117609; exp_bank=810; exp_qual=6; exp_decode=6; exp_frames=18432;;
    phone-control-angle.jpg) exp_evals=21203; exp_bank=11; exp_qual=0; exp_decode=0; exp_frames=0;;
    phone-b-mild.jpg) exp_evals=79259; exp_bank=937; exp_qual=935; exp_decode=691; exp_frames=2120047;;
    phone-b-angle.jpg) exp_evals=334857; exp_bank=6198; exp_qual=0; exp_decode=0; exp_frames=0;;
    *) exp_evals=0; exp_bank=0; exp_qual=0; exp_decode=0; exp_frames=0;;
  esac

  telemetry=true
  if [[ "$exp_evals" -gt 0 ]]; then
    [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b65" == true && "$b66" == true && "$b68" == true && "$b73" == true && "$b75" == true && "$b76" == true && "$b81" == true ]] || telemetry=false
    [[ "$p1w" -gt 0 && "$g2w" -gt 0 && "$g3w" -gt 0 && "$g4w" -gt 0 && "$g2t" -gt 0 && "$g3t" -gt 0 && "$g4t" -gt 0 ]] || telemetry=false
    [[ "$lut_folds" -gt 0 && "$lut_reads" -gt 0 && $((lut_success+lut_failed)) -eq "$lut_reads" ]] || telemetry=false
    if [[ "$exp_qual" -gt 0 ]]; then
      [[ "$b66_workers" -gt 0 && "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]] || telemetry=false
      [[ "$b64_auth" == false || "$b68_spec" -lt "$b66_workers" ]] || telemetry=false
    else
      [[ "$b66_workers" -eq 0 && "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]] || telemetry=false
    fi
  else
    [[ "$b64_attempted" == false && "$b73" == false && "$b75" == false && "$b76" == false && "$b81" == false ]] || telemetry=false
  fi

  payload_match='-'; qualification=INFO; gate_reason=ok
  if [[ "$role" == control ]]; then
    if [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false ]]; then qualification=PASS; else qualification=FAIL; gate_reason=control-semantics; fi
  else
    if [[ $rc -eq 0 ]]; then actual="$(cat "$payload_file")"; [[ "$actual" == "$expected_payload" ]] && payload_match=true || payload_match=false; else payload_match=false; fi
    case "$role" in
      required-direct) if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false && "$b64_attempted" == false ]]; then qualification=PASS; else qualification=FAIL; gate_reason=direct-semantics; fi;;
      required-build43) if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true && "$b64_attempted" == false ]]; then qualification=PASS; else qualification=FAIL; gate_reason=build43-semantics; fi;;
      required-any) if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false ]]; then qualification=PASS; else qualification=FAIL; gate_reason=any-semantics; fi;;
      required-build81) if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b81" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=build81-positive-semantics; fi;;
      required-reject) if [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b81" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=reject-semantics; fi;;
    esac
  fi
  if [[ "$telemetry" != true ]]; then
    qualification=FAIL
    if [[ "$gate_reason" == ok ]]; then gate_reason=telemetry; else gate_reason="${gate_reason}+telemetry"; fi
  fi

  base_elapsed="$(baseline_metric "$BUILD76_BASELINE_TSV" "$file" elapsed_ms)"
  base_geom="$(baseline_metric "$BUILD76_BASELINE_TSV" "$file" geometry_ms)"
  speed="$(speedup "$base_elapsed" "$elapsed_ms")"
  gspeed="$(speedup "$base_geom" "$geometry_ms")"
  row=("$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b75" "$b76" "$b81" "$p1w" "$g2w" "$g3w" "$g4w" "$g2t" "$g3t" "$g4t" "$freeze" "$p1wall" "$g2wall" "$g3wall" "$g4wall" "$lut_folds" "$lut_reads" "$lut_success" "$lut_failed" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$base_elapsed" "$speed" "$geometry_ms" "$base_geom" "$gspeed" "$gate_reason")
  (IFS=$'\t'; echo "${row[*]}") >> "$tsv"
done

{
  echo '# PixSeal Build81 exact luminance-LUT matrix'; echo
  echo 'Build81 is an equivalence-preserving performance candidate over qualified Build76. Only continuation4 FoldScore RGB->luminance products use exact 256-entry float64 lookup tables. Geometry, bank/order and protected-data semantics remain authoritative.'; echo
  echo '| image | role | evals | bank | qual | LUT folds | block reads | success | failed | geometry | B76 geom | geom speedup | elapsed | B76 elapsed | speedup | HMAC | telemetry eq | gate |'
  echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  tail -n +2 "$tsv" | awk -F '\t' '{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$1,$3,$4,$5,$6,$27,$28,$29,$30,$39,$40,$41,$36,$37,$38,$31,$33,$34}'
  echo
  echo 'Build76 remains the qualified baseline until Build81 is reproduced physically. Timing comparison is informative only after exact semantic equivalence passes.'
  [[ -f "$BUILD76_BASELINE_TSV" ]] && echo "Build76 timing baseline: $BUILD76_BASELINE_TSV" || echo 'Build76 timing baseline: unavailable'
} > "$md"
rm -f "$stage_dir/logs/"*.payload.bin

pass_count="$(awk -F '\t' 'NR>1&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
control_passes="$(awk -F '\t' 'NR>1&&$3=="control"&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-front.jpg"&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-mild.jpg"&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-a-angle.jpg"&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-front.jpg"&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-mild.jpg"&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR>1&&$1=="phone-b-angle.jpg"&&$34=="PASS"{n++}END{print n+0}' "$tsv")"
echo "Build81 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
if ((pass_count!=9 || control_passes!=3 || a_front_pass!=1 || a_mild_pass!=1 || a_angle_pass!=1 || b_front_pass!=1 || b_mild_pass!=1 || b_angle_pass!=1)); then
  echo 'Build81 failing cases:' >&2
  awk -F '\t' 'NR>1 && $34!="PASS" {printf "  %s role=%s reason=%s rc=%s hmac=%s telemetry=%s evals=%s bank=%s qual=%s decode=%s frames=%s lut=%s/%s/%s/%s\n",$1,$3,$42,$35,$31,$33,$4,$5,$6,$7,$8,$27,$28,$29,$30}' "$tsv" >&2
  rm -rf "$failed_output_dir"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  echo "Build81 FAILED diagnostics preserved in: $failed_output_dir" >&2
  echo 'error: Build81 semantic-equivalence physical gate not met' >&2
  exit 1
fi
mkdir -p "$OUTPUT_DIR"
rm -rf "$OUTPUT_DIR/logs"
mv "$stage_dir/logs" "$OUTPUT_DIR/logs"
mv "$tsv" "$final_tsv"
mv "$md" "$final_md"
rmdir "$stage_dir"
published=true
echo 'Build81 exact local luminance-lut matrix written to:'
echo "  $final_tsv"
echo "  $final_md"
echo 'Build81 semantic-equivalence physical gate: PASS'
