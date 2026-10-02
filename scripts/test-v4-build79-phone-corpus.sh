#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD79_DIAGNOSTIC_DIR:-v4-phone private/build79-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD79_TIMEOUT:-86400}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD79_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
 "phone-control-front.jpg|control|control|-"
 "phone-control-mild.jpg|control|control|-"
 "phone-control-angle.jpg|control|control|-"
 "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
 "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
 "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
 "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
 "phone-b-mild.jpg|B|required-build79|$MESSAGE_B"
 "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do
  IFS='|' read -r file _ <<< "$spec"
  [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }
done

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build79.XXXXXX")"
published=false
cleanup_stage(){ if [[ "$published" != true && -d "$stage_dir" ]]; then rm -rf "$stage_dir"; fi; }
trap cleanup_stage EXIT
mkdir -p "$stage_dir/logs"
tsv="$stage_dir/build79-phone-foldscore-profile.tsv"
md="$stage_dir/build79-phone-foldscore-profile.md"
final_tsv="$OUTPUT_DIR/build79-phone-foldscore-profile.tsv"
final_md="$OUTPUT_DIR/build79-phone-foldscore-profile.md"

printf 'image\tclass\trole\tbuild64_evals\tbuild64_bank\tbuild64_qualified\tbuild64_decode\tbuild64_frames\tbuild64_authenticated\tbuild68_physical_decode\tbuild68_speculative\tbuild75_attempted\tbuild76_attempted\tbuild79_attempted\tgen4_workers\tgen4_tasks\tcontinue4_calls\tcontinue4_outputs\tcontinue4_evals\tcontinue4_worker_ms\tscore_worker_ms\tprepare_worker_ms\tfold_calls\ttiles\tpilot_positions\tblock_reads\tblock_success\tblock_failed\tvisible\tfold_worker_ms\tsample_mod\tsampled_fold_calls\tsampled_fold_ms\tsampled_block_read_ms\tsampled_overhead_ms\tdetailed_blocks\tdetailed_pixels\treplay_map_ms\treplay_sample_ms\treplay_dct_ms\treplay_failures\thmac\tpayload_match\ttelemetry_equivalent\tqualification\texit_code\telapsed_ms\tgeometry_ms\tfreeze_ms\tprefix1_ms\tgen2_ms\tgen3_ms\tgen4_ms\n' > "$tsv"

extract_re(){
  local pattern="$1" file="$2" default_value="${3:-}" value
  value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"
  [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"
}

pass_count=0; control_passes=0; a_front_pass=0; a_mild_pass=0; a_angle_pass=0; b_front_pass=0; b_mild_pass=0; b_angle_pass=0
for spec in "${specs[@]}"; do
  IFS='|' read -r file class role expected_payload <<< "$spec"
  input="$ACQUISITION_DIR/$file"
  stem="${file%.*}"
  log="$stage_dir/logs/$stem.stderr.txt"
  payload_file="$stage_dir/logs/$stem.payload.bin"
  start_ns="$(date +%s%N)"
  set +e
  timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"
  rc=$?
  set -e
  end_ns="$(date +%s%N)"
  elapsed_ms=$(( (end_ns-start_ns)/1000000 ))

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
  b76=false; grep -q '^build76-gen2-parallel:' "$log" && b76=true
  b79=false; grep -q '^build79-foldscore-profile:' "$log" && b79=true

  gen4_workers=0; gen4_tasks=0; freeze=0; prefix1=0; gen2=0; gen3=0; gen4=0
  if [[ "$b76" == true ]]; then
    gen4_workers="$(extract_re 'build76-gen2-parallel: .* gen4-workers=([0-9]+).*' "$log" 0)"
    gen4_tasks="$(extract_re 'build76-gen2-parallel: .* gen4-tasks=([0-9]+).*' "$log" 0)"
    freeze="$(extract_re 'build76-gen2-parallel: .* freeze-ms=([0-9]+).*' "$log" 0)"
    prefix1="$(extract_re 'build76-gen2-parallel: .* prefix1-wall-ms=([0-9]+).*' "$log" 0)"
    gen2="$(extract_re 'build76-gen2-parallel: .* gen2-wall-ms=([0-9]+).*' "$log" 0)"
    gen3="$(extract_re 'build76-gen2-parallel: .* gen3-wall-ms=([0-9]+).*' "$log" 0)"
    gen4="$(extract_re 'build76-gen2-parallel: .* gen4-wall-ms=([0-9]+).*' "$log" 0)"
  fi

  cont_calls=0; cont_out=0; cont_e=0; cont_ms=0; score_ms=0; prep_ms=0
  fold_calls=0; tiles=0; pilot_positions=0; block_reads=0; block_success=0; block_failed=0; visible=0; fold_ms=0
  sample_mod=0; sampled_calls=0; sampled_ms=0; sampled_block_ms=0; sampled_overhead=0; detailed_blocks=0; detailed_pixels=0; replay_map=0; replay_sample=0; replay_dct=0; replay_failures=0
  if [[ "$b79" == true ]]; then
    cont_calls="$(extract_re 'build79-foldscore-profile: .* continue4-calls=([0-9]+).*' "$log" 0)"
    cont_out="$(extract_re 'build79-foldscore-profile: .* continue4-outputs=([0-9]+).*' "$log" 0)"
    cont_e="$(extract_re 'build79-foldscore-profile: .* continue4-evals=([0-9]+).*' "$log" 0)"
    cont_ms="$(extract_re 'build79-foldscore-profile: .* continue4-worker-ms=([0-9]+).*' "$log" 0)"
    score_ms="$(extract_re 'build79-foldscore-profile: .* score-worker-ms=([0-9]+).*' "$log" 0)"
    prep_ms="$(extract_re 'build79-foldscore-profile: .* prepare-worker-ms=([0-9]+).*' "$log" 0)"
    fold_calls="$(extract_re 'build79-foldscore-profile: .* fold-calls=([0-9]+).*' "$log" 0)"
    tiles="$(extract_re 'build79-foldscore-profile: .* tiles=([0-9]+).*' "$log" 0)"
    pilot_positions="$(extract_re 'build79-foldscore-profile: .* pilot-positions=([0-9]+).*' "$log" 0)"
    block_reads="$(extract_re 'build79-foldscore-profile: .* block-reads=([0-9]+).*' "$log" 0)"
    block_success="$(extract_re 'build79-foldscore-profile: .* block-success=([0-9]+).*' "$log" 0)"
    block_failed="$(extract_re 'build79-foldscore-profile: .* block-failed=([0-9]+).*' "$log" 0)"
    visible="$(extract_re 'build79-foldscore-profile: .* visible=([0-9]+).*' "$log" 0)"
    fold_ms="$(extract_re 'build79-foldscore-profile: .* fold-worker-ms=([0-9]+).*' "$log" 0)"
    sample_mod="$(extract_re 'build79-foldscore-profile: .* sample-mod=([0-9]+).*' "$log" 0)"
    sampled_calls="$(extract_re 'build79-foldscore-profile: .* sampled-fold-calls=([0-9]+).*' "$log" 0)"
    sampled_ms="$(extract_re 'build79-foldscore-profile: .* sampled-fold-ms=([0-9]+).*' "$log" 0)"
    sampled_block_ms="$(extract_re 'build79-foldscore-profile: .* sampled-block-read-ms=([0-9]+).*' "$log" 0)"
    sampled_overhead="$(extract_re 'build79-foldscore-profile: .* sampled-overhead-ms=([0-9]+).*' "$log" 0)"
    detailed_blocks="$(extract_re 'build79-foldscore-profile: .* detailed-blocks=([0-9]+).*' "$log" 0)"
    detailed_pixels="$(extract_re 'build79-foldscore-profile: .* detailed-pixels=([0-9]+).*' "$log" 0)"
    replay_map="$(extract_re 'build79-foldscore-profile: .* replay-map-ms=([0-9]+).*' "$log" 0)"
    replay_sample="$(extract_re 'build79-foldscore-profile: .* replay-sample-ms=([0-9]+).*' "$log" 0)"
    replay_dct="$(extract_re 'build79-foldscore-profile: .* replay-dct-ms=([0-9]+).*' "$log" 0)"
    replay_failures="$(extract_re 'build79-foldscore-profile: .* replay-failures=([0-9]+).*' "$log" 0)"
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
    [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b65" == true && "$b66" == true && "$b68" == true && "$b73" == true && "$b75" == true && "$b76" == true && "$b79" == true ]] || telemetry=false
    [[ "$gen4_tasks" -gt 0 && "$gen4_workers" -gt 0 && "$fold_calls" -eq "$cont_e" && "$block_reads" -eq "$pilot_positions" && "$block_success" -eq "$visible" && $((block_success+block_failed)) -eq "$block_reads" && "$sample_mod" -eq 64 && "$sampled_calls" -gt 0 && "$detailed_blocks" -gt 0 && "$detailed_pixels" -gt 0 && "$replay_failures" -eq 0 ]] || telemetry=false
    if [[ "$exp_qual" -gt 0 ]]; then
      [[ "$b66_workers" -gt 0 && "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]] || telemetry=false
      [[ "$b64_auth" == false || "$b68_spec" -lt "$b66_workers" ]] || telemetry=false
    else
      [[ "$b66_workers" -eq 0 && "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]] || telemetry=false
    fi
  else
    [[ "$b64_attempted" == false && "$b73" == false && "$b75" == false && "$b76" == false && "$b79" == false ]] || telemetry=false
  fi

  payload_match='-'; qualification=INFO
  if [[ "$role" == control ]]; then
    [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false ]] && qualification=PASS || qualification=FAIL
  else
    if [[ $rc -eq 0 ]]; then actual="$(cat "$payload_file")"; [[ "$actual" == "$expected_payload" ]] && payload_match=true || payload_match=false; else payload_match=false; fi
    case "$role" in
      required-direct) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == false && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL;;
      required-build43) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b43_attempted" == true && "$b43_auth" == true && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL;;
      required-any) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_attempted" == false ]] && qualification=PASS || qualification=FAIL;;
      required-build79) [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b79" == true ]] && qualification=PASS || qualification=FAIL;;
      required-reject) [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b79" == true ]] && qualification=PASS || qualification=FAIL;;
    esac
  fi
  [[ "$telemetry" == true ]] || qualification=FAIL

  printf '%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%d\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n' \
    "$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b75" "$b76" "$b79" "$gen4_workers" "$gen4_tasks" "$cont_calls" "$cont_out" "$cont_e" "$cont_ms" "$score_ms" "$prep_ms" "$fold_calls" "$tiles" "$pilot_positions" "$block_reads" "$block_success" "$block_failed" "$visible" "$fold_ms" "$sample_mod" "$sampled_calls" "$sampled_ms" "$sampled_block_ms" "$sampled_overhead" "$detailed_blocks" "$detailed_pixels" "$replay_map" "$replay_sample" "$replay_dct" "$replay_failures" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$geometry_ms" "$freeze" "$prefix1" "$gen2" "$gen3" "$gen4" >> "$tsv"

  if [[ "$qualification" == PASS ]]; then
    pass_count=$((pass_count+1))
    case "$file" in
      phone-control-*.jpg) control_passes=$((control_passes+1));;
      phone-a-front.jpg) a_front_pass=1;; phone-a-mild.jpg) a_mild_pass=1;; phone-a-angle.jpg) a_angle_pass=1;;
      phone-b-front.jpg) b_front_pass=1;; phone-b-mild.jpg) b_mild_pass=1;; phone-b-angle.jpg) b_angle_pass=1;;
    esac
  fi
done

{
  echo '# PixSeal Build79 FoldScore kernel profile'; echo
  echo 'Build79 is observability-only over the qualified Build76 baseline. It preserves the exact Build76 bank/order and scheduler. Exact counters cover every continuation4 FoldScore; rare deterministic samples decompose projective block reads and replay one successful block into mapPoint, bilinear luminance sampling and DCT phases. Timing overhead makes Build79 non-promotable.'; echo
  echo '| image | role | evals | bank | qual | fold calls | block reads | failed | sampled folds | block-read ms | overhead ms | map ms | sample ms | DCT ms | replay fail | HMAC | telemetry eq | elapsed | geometry | gate |'
  echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  tail -n +2 "$tsv" | awk -F '\t' '{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$1,$3,$4,$5,$6,$23,$26,$28,$32,$34,$35,$38,$39,$40,$41,$42,$44,$47,$48,$45}'
  echo
  echo 'Correctness is authoritative. Build76 remains the qualified baseline; Build79 timing exists only to characterize the continuation4 FoldScore kernel.'
} > "$md"
rm -f "$stage_dir/logs/"*.payload.bin

echo "Build79 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
if ((pass_count!=9 || control_passes!=3 || a_front_pass!=1 || a_mild_pass!=1 || a_angle_pass!=1 || b_front_pass!=1 || b_mild_pass!=1 || b_angle_pass!=1)); then
  echo 'error: Build79 semantic-equivalence physical gate not met' >&2
  exit 1
fi
mkdir -p "$OUTPUT_DIR"
rm -rf "$OUTPUT_DIR/logs"
mv "$stage_dir/logs" "$OUTPUT_DIR/logs"
mv "$tsv" "$final_tsv"
mv "$md" "$final_md"
rmdir "$stage_dir"
published=true
echo 'Build79 FoldScore kernel profile written to:'
echo "  $final_tsv"
echo "  $final_md"
echo 'Build79 semantic-equivalence physical gate: PASS'
