#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD100_DIAGNOSTIC_DIR:-v4-phone private/build100-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD100_TIMEOUT:-86400}"
BUILD98_PERFORMANCE_TSV="${V4_PHONE_BUILD98_PERFORMANCE_TSV:-docs/qualified-baselines/build98-phone-performance.tsv}"
ENTRYPOINT="${V4_PHONE_BUILD100_ENTRYPOINT:-v4-extract-phone-build100}"
GATE_PURPOSE="${V4_PHONE_BUILD100_GATE_PURPOSE:-physical-promotion-candidate}"
CURRENT_BASELINE="${V4_PHONE_BUILD100_CURRENT_BASELINE:-Build98}"
PREVIOUS_BASELINE="${V4_PHONE_BUILD100_PREVIOUS_BASELINE:-Build93}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-unreported}"
export GOMAXPROCS="${GOMAXPROCS:-8}"
BUILD98_BASELINE_TSV="${V4_PHONE_BUILD98_BASELINE_TSV:-docs/qualified-baselines/build98-production-integration.tsv}"
[[ -f "$BUILD98_BASELINE_TSV" ]] || { echo "error: missing Build98 semantic baseline" >&2; exit 1; }
[[ ! -e "$OUTPUT_DIR" && ! -L "$OUTPUT_DIR" ]] || { echo "error: existing physical output is preserved; choose another directory" >&2; exit 1; }

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD100_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
 "phone-control-front.jpg|control|control|-"
 "phone-control-mild.jpg|control|control|-"
 "phone-control-angle.jpg|control|control|-"
 "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
 "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
 "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
 "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
 "phone-b-mild.jpg|B|required-build100|$MESSAGE_B"
 "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do
  IFS='|' read -r file _ <<< "$spec"
  [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }
done

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build98.XXXXXX")"
published=false
cleanup_stage(){ local rc=$?; trap - EXIT; if [[ "$published" != true && -d "$stage_dir" ]]; then mv "$stage_dir" "${OUTPUT_DIR}-failed-$(basename "$stage_dir")"; fi; exit "$rc"; }
trap cleanup_stage EXIT
mkdir -p "$stage_dir/logs"
tsv="$stage_dir/build100-phone-pair-scoring.tsv"
md="$stage_dir/build100-phone-pair-scoring.md"
metadata="$stage_dir/build100-phone-run-metadata.txt"
failed_output_dir="${OUTPUT_DIR}-failed"

pixseal_sha256="-"
if command -v sha256sum >/dev/null 2>&1; then
  pixseal_sha256="$(sha256sum "$PIXSEAL" | awk '{print $1}')"
fi
{
  echo "PixSeal Build100 ordered-parallel public pair scoring semantic matrix"
  echo "created_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "gate_purpose=$GATE_PURPOSE"
  echo "entrypoint=$ENTRYPOINT"
  echo "current_qualified_baseline=$CURRENT_BASELINE"
  echo "previous_qualified_baseline=$PREVIOUS_BASELINE"
  echo "implementation=Build100"
  echo "toolchain=$TOOLCHAIN"
  echo "pixseal_sha256=$pixseal_sha256"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "case_count=${#specs[@]}"
  echo "additional_image_sized_memory_bytes=0"
  echo "paired_evidence=full-freeze and pair-score ABBA comparison must be reviewed before this promotion gate"
} > "$metadata"

header=(build100_attempted build100_pair_workers build100_pair_tasks build100_pair_evaluations build100_pair_wall_ms build100_pair_worker_ms image class role build64_evals build64_bank build64_qualified build64_decode build64_frames build64_authenticated build68_physical_decode build68_speculative build68_physical_profiles build68_physical_list_frames build84_attempted build90_attempted qualification_workers qualification_tasks qualification_evaluations qualification_ms build93_attempted protected_sampling_passes protected_profile_derivations fetch_fold_scores fetch_block_reads fetch_block_success fetch_block_failed build98_attempted gen3_fetch_fold_scores gen3_fetch_block_reads gen3_fetch_block_success gen3_fetch_block_failed hmac payload_match telemetry_equivalent qualification exit_code elapsed_ms deep_total_ms geometry_ms gen3_ms gen4_ms decode_ms build98_reference_deep_ms speedup_vs_build98 gate_reason)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"

extract_re(){ local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
build98_reference_deep(){ local path="$1" image="$2"; [[ -f "$path" ]] || { printf '%s' '-'; return; }; awk -F '\t' -v image="$image" 'NR==1{for(i=1;i<=NF;i++){if($i=="scope")s=i;if($i=="image")im=i;if($i=="deep_total_ms")tm=i}next}$s=="mean"&&$im==image&&tm>0{print $tm;found=1;exit}END{if(!found)print "-"}' "$path"; }
speedup(){ local base="$1" now="$2"; if [[ "$base" =~ ^[0-9]+([.][0-9]+)?$ && "$now" =~ ^[0-9]+([.][0-9]+)?$ ]] && awk -v n="$now" 'BEGIN{exit !(n>0)}'; then awk -v b="$base" -v n="$now" 'BEGIN{printf "%.3f",b/n}'; else printf '%s' '-'; fi; }

for spec in "${specs[@]}"; do
  IFS='|' read -r file class role expected_payload <<< "$spec"
  input="$ACQUISITION_DIR/$file"; stem="${file%.*}"
  log="$stage_dir/logs/$stem.stderr.txt"; payload_file="$stage_dir/logs/$stem.payload.bin"
  start_ns="$(date +%s%N)"
  set +e
  timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" "$ENTRYPOINT" -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"
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

  b68_physical=0; b68_spec=0; b68_profiles=0; b68_physical_frames=0; deep_total_ms=0; geometry_ms=0; qualification_ms=0; decode_ms=0
  if grep -q '^build68-plane-reuse:' "$log"; then
    b68_physical="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"
    b68_spec="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)"
    b68_profiles="$(extract_re 'build68-plane-reuse: .* physical-profiles=([0-9]+).*' "$log" 0)"
    b68_physical_frames="$(extract_re 'build68-plane-reuse: .* physical-list-frames=([0-9]+).*' "$log" 0)"
    deep_total_ms="$(extract_re 'build68-plane-reuse: .* total-ms=([0-9]+).*' "$log" 0)"
    geometry_ms="$(extract_re 'build68-plane-reuse: .* geometry-ms=([0-9]+).*' "$log" 0)"
    qualification_ms="$(extract_re 'build68-plane-reuse: .* qualification-ms=([0-9]+).*' "$log" 0)"
    decode_ms="$(extract_re 'build68-plane-reuse: .* decode-wall-ms=([0-9]+).*' "$log" 0)"
  fi

  gen3_ms="$(extract_re 'build73-gen3-parallel: .* gen3-wall-ms=([0-9]+).*' "$log" 0)"
  gen4_ms="$(extract_re 'build73-gen3-parallel: .* gen4-wall-ms=([0-9]+).*' "$log" 0)"

  b84=false; fetch_folds=0; fetch_reads=0; fetch_success=0; fetch_failed=0
  if grep -q '^build84-rgb-fetch:' "$log"; then
    b84=true
    fetch_folds="$(extract_re 'build84-rgb-fetch: .* fold-scores=([0-9]+).*' "$log" 0)"
    fetch_reads="$(extract_re 'build84-rgb-fetch: .* block-reads=([0-9]+).*' "$log" 0)"
    fetch_success="$(extract_re 'build84-rgb-fetch: .* block-success=([0-9]+).*' "$log" 0)"
    fetch_failed="$(extract_re 'build84-rgb-fetch: .* block-failed=([0-9]+).*' "$log" 0)"
  fi

  b90=false; qworkers=0; qtasks=0; qevals=0
  if grep -q '^build90-qualification-parallel:' "$log"; then
    b90=true
    qworkers="$(extract_re 'build90-qualification-parallel: .* workers=([0-9]+).*' "$log" 0)"
    qtasks="$(extract_re 'build90-qualification-parallel: .* tasks=([0-9]+).*' "$log" 0)"
    qevals="$(extract_re 'build90-qualification-parallel: .* evaluations=([0-9]+).*' "$log" 0)"
  fi

  b93=false; sampling_passes=0; profile_derivations=0
  if grep -q '^build93-residue-reuse:' "$log"; then
    b93=true
    sampling_passes="$(extract_re 'build93-residue-reuse: .* sampling-passes=([0-9]+).*' "$log" 0)"
    profile_derivations="$(extract_re 'build93-residue-reuse: .* profile-derivations=([0-9]+).*' "$log" 0)"
  fi

  b98=false; gen3_folds=0; gen3_reads=0; gen3_success=0; gen3_failed=0
  if grep -q '^build98-gen3-fetch:' "$log"; then
    b98=true
    gen3_folds="$(extract_re 'build98-gen3-fetch: .* fold-scores=([0-9]+).*' "$log" 0)"
    gen3_reads="$(extract_re 'build98-gen3-fetch: .* block-reads=([0-9]+).*' "$log" 0)"
    gen3_success="$(extract_re 'build98-gen3-fetch: .* block-success=([0-9]+).*' "$log" 0)"
    gen3_failed="$(extract_re 'build98-gen3-fetch: .* block-failed=([0-9]+).*' "$log" 0)"
  fi

  case "$file" in
    phone-control-front.jpg) exp_evals=18021; exp_bank=26; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=1104; exp_fetch_reads=1413120; exp_fetch_success=1413120; exp_fetch_failed=0; exp_gen3_folds=0; exp_gen3_reads=0; exp_gen3_success=0; exp_gen3_failed=0;;
    phone-control-mild.jpg) exp_evals=117609; exp_bank=810; exp_qual=6; exp_decode=6; exp_frames=18432; exp_fetch_folds=20368; exp_fetch_reads=26071040; exp_fetch_success=26071040; exp_fetch_failed=0; exp_gen3_folds=0; exp_gen3_reads=0; exp_gen3_success=0; exp_gen3_failed=0;;
    phone-control-angle.jpg) exp_evals=21203; exp_bank=11; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=304; exp_fetch_reads=389120; exp_fetch_success=389120; exp_fetch_failed=0; exp_gen3_folds=0; exp_gen3_reads=0; exp_gen3_success=0; exp_gen3_failed=0;;
    phone-b-mild.jpg) exp_evals=79259; exp_bank=937; exp_qual=935; exp_decode=691; exp_frames=2120047; exp_fetch_folds=16496; exp_fetch_reads=21114880; exp_fetch_success=21114880; exp_fetch_failed=0; exp_gen3_folds=23024; exp_gen3_reads=29470720; exp_gen3_success=29470720; exp_gen3_failed=0;;
    phone-b-angle.jpg) exp_evals=334857; exp_bank=6198; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=118004; exp_fetch_reads=151045120; exp_fetch_success=149370664; exp_fetch_failed=1674456; exp_gen3_folds=76817; exp_gen3_reads=98325760; exp_gen3_success=97394512; exp_gen3_failed=931248;;
    *) exp_evals=0; exp_bank=0; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=0; exp_fetch_reads=0; exp_fetch_success=0; exp_fetch_failed=0; exp_gen3_folds=0; exp_gen3_reads=0; exp_gen3_success=0; exp_gen3_failed=0;;
  esac

  b100=false; pair_workers=0; pair_tasks=0; pair_evals=0; pair_wall=0; pair_worker=0
  if grep -q '^build100-pair-scoring:' "$log"; then
    b100=true
    pair_workers="$(extract_re 'build100-pair-scoring: .* workers=([0-9]+).*' "$log" 0)"
    pair_tasks="$(extract_re 'build100-pair-scoring: .* tasks=([0-9]+).*' "$log" 0)"
    pair_evals="$(extract_re 'build100-pair-scoring: .* evaluations=([0-9]+).*' "$log" 0)"
    pair_wall="$(extract_re 'build100-pair-scoring: .* wall-ms=([0-9]+).*' "$log" 0)"
    pair_worker="$(extract_re 'build100-pair-scoring: .* worker-ms=([0-9]+).*' "$log" 0)"
  fi
  telemetry=true
  if [[ "$exp_evals" -gt 0 ]]; then
    [[ "$b100" == true && "$pair_workers" -eq 6 && "$pair_tasks" -eq 6 && "$pair_evals" -gt 0 ]] || telemetry=false
  else
    [[ "$b100" == false && "$pair_tasks" -eq 0 ]] || telemetry=false
  fi
  if [[ "$exp_evals" -gt 0 ]]; then
    [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b84" == true && "$b90" == true && "$b93" == true && "$b98" == true ]] || telemetry=false
    [[ "$qtasks" -eq "$exp_bank" && "$qworkers" -gt 0 && "$qworkers" -le "$qtasks" && "$qevals" -ge "$qtasks" && "$qevals" -le $((2*qtasks)) ]] || telemetry=false
    [[ "$fetch_folds" -eq "$exp_fetch_folds" && "$fetch_reads" -eq "$exp_fetch_reads" && "$fetch_success" -eq "$exp_fetch_success" && "$fetch_failed" -eq "$exp_fetch_failed" && $((fetch_success+fetch_failed)) -eq "$fetch_reads" ]] || telemetry=false
    [[ $((gen3_success+gen3_failed)) -eq "$gen3_reads" ]] || telemetry=false
    if [[ "$gen3_reads" -eq 0 ]]; then [[ "$gen3_folds" -eq 0 ]] || telemetry=false; else [[ "$gen3_folds" -gt 0 ]] || telemetry=false; fi
    if [[ "$exp_gen3_folds" -gt 0 ]]; then
      [[ "$gen3_folds" -eq "$exp_gen3_folds" && "$gen3_reads" -eq "$exp_gen3_reads" && "$gen3_success" -eq "$exp_gen3_success" && "$gen3_failed" -eq "$exp_gen3_failed" ]] || telemetry=false
    fi
    if [[ "$exp_qual" -gt 0 ]]; then
      [[ "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]] || telemetry=false
    else
      [[ "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]] || telemetry=false
    fi
    [[ "$sampling_passes" -eq "$b68_physical" ]] || telemetry=false
    if [[ "$sampling_passes" -eq 0 ]]; then
      [[ "$profile_derivations" -eq 0 ]] || telemetry=false
    else
      [[ "$profile_derivations" -ge "$b68_profiles" && "$profile_derivations" -le $((3*sampling_passes)) ]] || telemetry=false
    fi
  else
    [[ "$b64_attempted" == false && "$b84" == false && "$b90" == false && "$b93" == false && "$b98" == false && "$sampling_passes" -eq 0 && "$profile_derivations" -eq 0 && "$gen3_reads" -eq 0 ]] || telemetry=false
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
      required-build100) if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b84" == true && "$b90" == true && "$b93" == true && "$b98" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=build98-positive-semantics; fi;;
      required-reject) if [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b84" == true && "$b90" == true && "$b93" == true && "$b98" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=reject-semantics; fi;;
    esac
  fi
  if [[ "$telemetry" != true ]]; then qualification=FAIL; [[ "$gate_reason" == ok ]] && gate_reason=telemetry || gate_reason="${gate_reason}+telemetry"; fi

  base_total="$(build98_reference_deep "$BUILD98_PERFORMANCE_TSV" "$file")"
  deep_speed="$(speedup "$base_total" "$deep_total_ms")"
  row=("$b100" "$pair_workers" "$pair_tasks" "$pair_evals" "$pair_wall" "$pair_worker" "$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b68_profiles" "$b68_physical_frames" "$b84" "$b90" "$qworkers" "$qtasks" "$qevals" "$qualification_ms" "$b93" "$sampling_passes" "$profile_derivations" "$fetch_folds" "$fetch_reads" "$fetch_success" "$fetch_failed" "$b98" "$gen3_folds" "$gen3_reads" "$gen3_success" "$gen3_failed" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$deep_total_ms" "$geometry_ms" "$gen3_ms" "$gen4_ms" "$decode_ms" "$base_total" "$deep_speed" "$gate_reason")
  if (( ${#row[@]} != ${#header[@]} )); then echo "error: Build100 TSV field mismatch for $file: header=${#header[@]} row=${#row[@]}" >&2; exit 1; fi
  (IFS=$'\t'; echo "${row[*]}") >> "$tsv"
done

{
  echo '# PixSeal Build100 ordered-parallel public pair scoring physical matrix'; echo
  echo "Entry point: \`$ENTRYPOINT\`; gate purpose: \`$GATE_PURPOSE\`; current qualified baseline: \`$CURRENT_BASELINE\`."; echo
  echo 'Build100 changes only Generation3 FoldScore dispatch and adds no image-sized memory. Historical Build84 continuation4, Build90 qualification and Build93 protected-decode counters remain frozen.'; echo
  echo '| image | role | bank | qualified | Gen3 fetch folds | Gen3 fetch reads | gen3 ms | gen4 ms | qualification ms | decode ms | deep total ms | Build98 mean | speedup | HMAC | telemetry eq | gate |'
  echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$c["image"],$c["role"],$c["build64_bank"],$c["build64_qualified"],$c["gen3_fetch_fold_scores"],$c["gen3_fetch_block_reads"],$c["gen3_ms"],$c["gen4_ms"],$c["qualification_ms"],$c["decode_ms"],$c["deep_total_ms"],$c["build98_reference_deep_ms"],$c["speedup_vs_build98"],$c["hmac"],$c["telemetry_equivalent"],$c["qualification"]}' "$tsv"
  echo
  echo 'Build100 promotion is already supported by two independent Go 1.26.0 9/9 runs. This matrix is integration/requalification evidence; timing remains fully retained.'
} > "$md"
rm -f "$stage_dir/logs/"*.payload.bin

# Compare the complete inherited semantic view against the qualified Build98
# integration table. Timing, entrypoint, role names and candidate-only fields are
# deliberately excluded; no tolerance is applied to semantic counters.
awk -F '\t' '
BEGIN {split("build64_evals build64_bank build64_qualified build64_decode build64_frames build64_authenticated build68_physical_decode build68_speculative build68_physical_profiles build68_physical_list_frames build84_attempted build90_attempted qualification_workers qualification_tasks qualification_evaluations build93_attempted protected_sampling_passes protected_profile_derivations fetch_fold_scores fetch_block_reads fetch_block_success fetch_block_failed build98_attempted gen3_fetch_fold_scores gen3_fetch_block_reads gen3_fetch_block_success gen3_fetch_block_failed hmac payload_match", names, " ")}
FNR==1 {delete columns;for(i=1;i<=NF;i++)columns[$i]=i;for(j in names)if(!(names[j] in columns)){print "error: missing semantic column " names[j] > "/dev/stderr";fatal=1;exit 1};next}
NR==FNR {im=$columns["image"];seen[im]=1;for(j in names)values[im SUBSEP names[j]]=$columns[names[j]];next}
{im=$columns["image"];if(!(im in seen)){print "error: unknown semantic baseline image" > "/dev/stderr";bad=1};for(j in names)if($columns[names[j]]!=values[im SUBSEP names[j]]){print "error: " im " changed semantic field " names[j] > "/dev/stderr";bad=1}}
END {exit (bad||fatal)}
' "$BUILD98_BASELINE_TSV" "$tsv"

pass_count="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)if($i=="qualification")q=i;next}$q=="PASS"{n++}END{print n+0}' "$tsv")"
control_passes="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="role")r=i}next}$r=="control"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-a-front.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-a-mild.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-a-angle.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-b-front.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-b-mild.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-b-angle.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"

echo "Build100 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
if ((pass_count!=9 || control_passes!=3 || a_front_pass!=1 || a_mild_pass!=1 || a_angle_pass!=1 || b_front_pass!=1 || b_mild_pass!=1 || b_angle_pass!=1)); then
  echo 'Build100 failing cases:' >&2
  awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}$c["qualification"]!="PASS"{printf "  %s role=%s reason=%s rc=%s hmac=%s telemetry=%s evals=%s bank=%s qual=%s gen3-folds=%s gen3-reads=%s\n",$c["image"],$c["role"],$c["gate_reason"],$c["exit_code"],$c["hmac"],$c["telemetry_equivalent"],$c["build64_evals"],$c["build64_bank"],$c["build64_qualified"],$c["gen3_fetch_fold_scores"],$c["gen3_fetch_block_reads"]}' "$tsv" >&2
  failed_output_dir="${failed_output_dir}-$(basename "$stage_dir")"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  echo "Build100 FAILED diagnostics preserved in: $failed_output_dir" >&2
  echo 'error: Build100 semantic-equivalence physical gate not met' >&2
  exit 1
fi

mv -T "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build100 ordered-parallel public pair scoring matrix written to:'
echo "  $OUTPUT_DIR/build100-phone-pair-scoring.tsv"
echo "  $OUTPUT_DIR/build100-phone-pair-scoring.md"
echo "  $OUTPUT_DIR/build100-phone-run-metadata.txt"
echo "Build100 semantic-equivalence gate: PASS entrypoint=$ENTRYPOINT purpose=$GATE_PURPOSE"
