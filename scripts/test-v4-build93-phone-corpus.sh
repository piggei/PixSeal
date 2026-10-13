#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD93_DIAGNOSTIC_DIR:-v4-phone private/build93-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD93_TIMEOUT:-86400}"
BUILD90_BASELINE_TSV="${V4_PHONE_BUILD90_BASELINE_TSV:-docs/qualified-baselines/build90-phone-performance.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD93_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
 "phone-control-front.jpg|control|control|-"
 "phone-control-mild.jpg|control|control|-"
 "phone-control-angle.jpg|control|control|-"
 "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
 "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
 "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
 "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
 "phone-b-mild.jpg|B|required-build93|$MESSAGE_B"
 "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do
  IFS='|' read -r file _ <<< "$spec"
  [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }
done

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build93.XXXXXX")"
published=false
cleanup_stage(){ if [[ "$published" != true && -d "$stage_dir" ]]; then rm -rf "$stage_dir"; fi; }
trap cleanup_stage EXIT
mkdir -p "$stage_dir/logs"
tsv="$stage_dir/build93-phone-residue-reuse.tsv"
md="$stage_dir/build93-phone-residue-reuse.md"
failed_output_dir="${OUTPUT_DIR}-failed"

header=(image class role build64_evals build64_bank build64_qualified build64_decode build64_frames build64_authenticated build68_physical_decode build68_speculative build68_physical_profiles build68_physical_list_frames build84_attempted build90_attempted qualification_workers qualification_tasks qualification_evaluations qualification_ms build93_attempted protected_sampling_passes protected_profile_derivations fetch_fold_scores fetch_block_reads fetch_block_success fetch_block_failed hmac payload_match telemetry_equivalent qualification exit_code elapsed_ms deep_total_ms geometry_ms decode_ms build90_reference_total_ms speedup_vs_build90 gate_reason)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"

extract_re(){ local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
qualified_baseline_total(){ local path="$1" image="$2"; [[ -f "$path" ]] || { printf '%s' '-'; return; }; awk -F '\t' -v image="$image" 'NR==1{for(i=1;i<=NF;i++){if($i=="scope")s=i;if($i=="image")im=i;if($i=="total_ms")tm=i}next}$s=="paired_mean"&&$im==image&&tm>0{print $tm;found=1;exit}END{if(!found)print "-"}' "$path"; }
speedup(){ local base="$1" now="$2"; if [[ "$base" =~ ^[0-9]+([.][0-9]+)?$ && "$now" =~ ^[0-9]+([.][0-9]+)?$ ]] && awk -v n="$now" 'BEGIN{exit !(n>0)}'; then awk -v b="$base" -v n="$now" 'BEGIN{printf "%.3f",b/n}'; else printf '%s' '-'; fi; }

for spec in "${specs[@]}"; do
  IFS='|' read -r file class role expected_payload <<< "$spec"
  input="$ACQUISITION_DIR/$file"; stem="${file%.*}"
  log="$stage_dir/logs/$stem.stderr.txt"; payload_file="$stage_dir/logs/$stem.payload.bin"
  start_ns="$(date +%s%N)"
  set +e
  timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone-build93 -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"
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

  case "$file" in
    phone-control-front.jpg) exp_evals=18021; exp_bank=26; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=1104; exp_fetch_reads=1413120; exp_fetch_success=1413120; exp_fetch_failed=0;;
    phone-control-mild.jpg) exp_evals=117609; exp_bank=810; exp_qual=6; exp_decode=6; exp_frames=18432; exp_fetch_folds=20368; exp_fetch_reads=26071040; exp_fetch_success=26071040; exp_fetch_failed=0;;
    phone-control-angle.jpg) exp_evals=21203; exp_bank=11; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=304; exp_fetch_reads=389120; exp_fetch_success=389120; exp_fetch_failed=0;;
    phone-b-mild.jpg) exp_evals=79259; exp_bank=937; exp_qual=935; exp_decode=691; exp_frames=2120047; exp_fetch_folds=16496; exp_fetch_reads=21114880; exp_fetch_success=21114880; exp_fetch_failed=0;;
    phone-b-angle.jpg) exp_evals=334857; exp_bank=6198; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=118004; exp_fetch_reads=151045120; exp_fetch_success=149370664; exp_fetch_failed=1674456;;
    *) exp_evals=0; exp_bank=0; exp_qual=0; exp_decode=0; exp_frames=0; exp_fetch_folds=0; exp_fetch_reads=0; exp_fetch_success=0; exp_fetch_failed=0;;
  esac

  telemetry=true
  if [[ "$exp_evals" -gt 0 ]]; then
    [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b84" == true && "$b90" == true && "$b93" == true ]] || telemetry=false
    [[ "$qtasks" -eq "$exp_bank" && "$qworkers" -gt 0 && "$qworkers" -le "$qtasks" && "$qevals" -ge "$qtasks" && "$qevals" -le $((2*qtasks)) ]] || telemetry=false
    [[ "$fetch_folds" -eq "$exp_fetch_folds" && "$fetch_reads" -eq "$exp_fetch_reads" && "$fetch_success" -eq "$exp_fetch_success" && "$fetch_failed" -eq "$exp_fetch_failed" && $((fetch_success+fetch_failed)) -eq "$fetch_reads" ]] || telemetry=false
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
    [[ "$b64_attempted" == false && "$b84" == false && "$b90" == false && "$b93" == false && "$sampling_passes" -eq 0 && "$profile_derivations" -eq 0 ]] || telemetry=false
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
      required-build93) if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b84" == true && "$b90" == true && "$b93" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=build93-positive-semantics; fi;;
      required-reject) if [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b84" == true && "$b90" == true && "$b93" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=reject-semantics; fi;;
    esac
  fi
  if [[ "$telemetry" != true ]]; then qualification=FAIL; [[ "$gate_reason" == ok ]] && gate_reason=telemetry || gate_reason="${gate_reason}+telemetry"; fi

  base_total="$(qualified_baseline_total "$BUILD90_BASELINE_TSV" "$file")"
  deep_speed="$(speedup "$base_total" "$deep_total_ms")"
  row=("$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b68_profiles" "$b68_physical_frames" "$b84" "$b90" "$qworkers" "$qtasks" "$qevals" "$qualification_ms" "$b93" "$sampling_passes" "$profile_derivations" "$fetch_folds" "$fetch_reads" "$fetch_success" "$fetch_failed" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$elapsed_ms" "$deep_total_ms" "$geometry_ms" "$decode_ms" "$base_total" "$deep_speed" "$gate_reason")
  if (( ${#row[@]} != ${#header[@]} )); then echo "error: Build93 TSV field mismatch for $file: header=${#header[@]} row=${#row[@]}" >&2; exit 1; fi
  (IFS=$'\t'; echo "${row[*]}") >> "$tsv"
done

{
  echo '# PixSeal Build93 exact protected-residue reuse physical matrix'; echo
  echo 'Build93 changes only protected-data sampling after the Build90 frozen bank and ordered-parallel public-pilot qualification. Build93 is now qualified; this retained gate uses the dedicated Build93-pinned entry point for reproducible requalification.'; echo
  echo '| image | role | bank | qualified | physical decode | profiles | sampling passes | derivations | qualification ms | decode ms | deep total ms | Build90 paired ref | speedup | HMAC | telemetry eq | gate |'
  echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$c["image"],$c["role"],$c["build64_bank"],$c["build64_qualified"],$c["build68_physical_decode"],$c["build68_physical_profiles"],$c["protected_sampling_passes"],$c["protected_profile_derivations"],$c["qualification_ms"],$c["decode_ms"],$c["deep_total_ms"],$c["build90_reference_total_ms"],$c["speedup_vs_build90"],$c["hmac"],$c["telemetry_equivalent"],$c["qualification"]}' "$tsv"
  echo
  echo 'Build93 is the current qualified smartphone baseline. This retained matrix is a reproducible requalification sample; promotion evidence is frozen separately under docs/qualified-baselines/.'
} > "$md"
rm -f "$stage_dir/logs/"*.payload.bin

pass_count="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)if($i=="qualification")q=i;next}$q=="PASS"{n++}END{print n+0}' "$tsv")"
control_passes="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="role")r=i}next}$r=="control"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
a_front_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-a-front.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
a_mild_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-a-mild.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
a_angle_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-a-angle.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
b_front_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-b-front.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
b_mild_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-b-mild.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"
b_angle_pass="$(awk -F '\t' 'NR==1{for(i=1;i<=NF;i++){if($i=="qualification")q=i;if($i=="image")im=i}next}$im=="phone-b-angle.jpg"&&$q=="PASS"{n++}END{print n+0}' "$tsv")"

echo "Build93 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
if ((pass_count!=9 || control_passes!=3 || a_front_pass!=1 || a_mild_pass!=1 || a_angle_pass!=1 || b_front_pass!=1 || b_mild_pass!=1 || b_angle_pass!=1)); then
  echo 'Build93 failing cases:' >&2
  awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}$c["qualification"]!="PASS"{printf "  %s role=%s reason=%s rc=%s hmac=%s telemetry=%s evals=%s bank=%s qual=%s samples=%s derivations=%s\n",$c["image"],$c["role"],$c["gate_reason"],$c["exit_code"],$c["hmac"],$c["telemetry_equivalent"],$c["build64_evals"],$c["build64_bank"],$c["build64_qualified"],$c["protected_sampling_passes"],$c["protected_profile_derivations"]}' "$tsv" >&2
  rm -rf "$failed_output_dir"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  echo "Build93 FAILED diagnostics preserved in: $failed_output_dir" >&2
  echo 'error: Build93 semantic-equivalence physical gate not met' >&2
  exit 1
fi

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build93 exact protected-residue reuse matrix written to:'
echo "  $OUTPUT_DIR/build93-phone-residue-reuse.tsv"
echo "  $OUTPUT_DIR/build93-phone-residue-reuse.md"
echo 'Build93 semantic-equivalence physical gate: PASS'
