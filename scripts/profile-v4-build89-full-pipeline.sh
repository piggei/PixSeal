#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD89_PROFILE_DIR:-v4-phone private/build89-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PHONE_TIMEOUT="${V4_PHONE_BUILD89_TIMEOUT:-86400}"
BUILD84_BASELINE_TSV="${V4_PHONE_BUILD84_BASELINE_TSV:-docs/qualified-baselines/build84-phone-performance.tsv}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ "$PHONE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD89_TIMEOUT must be an integer number of seconds" >&2; exit 1; }

specs=(
 "phone-control-front.jpg|control|control|-"
 "phone-control-mild.jpg|control|control|-"
 "phone-control-angle.jpg|control|control|-"
 "phone-a-front.jpg|A|required-build43|$MESSAGE_A"
 "phone-a-mild.jpg|A|required-any|$MESSAGE_A"
 "phone-a-angle.jpg|A|required-direct|$MESSAGE_A"
 "phone-b-front.jpg|B|required-direct|$MESSAGE_B"
 "phone-b-mild.jpg|B|required-build84|$MESSAGE_B"
 "phone-b-angle.jpg|B|required-reject|$MESSAGE_B"
)
for spec in "${specs[@]}"; do
  IFS='|' read -r file _ <<< "$spec"
  [[ -f "$ACQUISITION_DIR/$file" ]] || { echo "error: missing required acquisition: $ACQUISITION_DIR/$file" >&2; exit 1; }
done

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build89.XXXXXX")"
published=false
cleanup_stage(){ if [[ "$published" != true && -d "$stage_dir" ]]; then rm -rf "$stage_dir"; fi; }
trap cleanup_stage EXIT
mkdir -p "$stage_dir/logs"
tsv="$stage_dir/build89-full-pipeline.tsv"
md="$stage_dir/build89-full-pipeline.md"
metadata="$stage_dir/build89-metadata.txt"
final_tsv="$OUTPUT_DIR/build89-full-pipeline.tsv"
final_md="$OUTPUT_DIR/build89-full-pipeline.md"
failed_output_dir="${OUTPUT_DIR}-failed"

{
  echo "PixSeal Build89 full qualified Build84 deep-pipeline profiling"
  echo "toolchain=$TOOLCHAIN"
  echo "current_qualified_baseline=Build84"
  echo "active_deep_runtime=Build84-qualified"
  echo "build86_status=exact-benchmark-negative-non-promoted"
  echo "build88_status=exact-benchmark-strongly-negative-non-promoted"
  echo "build89_status=observability-only"
  echo "profile_scope=retained-nine-photo-matrix"
  echo "timing_probes=derived-from-existing-Build84-telemetry"
  echo "private_corpus_in_source=false"
} > "$metadata"

header=(image class role build64_evals build64_bank build64_qualified build64_decode build64_frames build64_authenticated build68_physical_decode build68_speculative build84_attempted build89_attempted total_ms geometry_ms geometry_plane_prep_ms freeze_ms prefix1_ms gen2_ms gen3_ms gen4_ms geometry_accounted_ms geometry_unaccounted_ms post_plane_prep_ms qualification_ms decode_wall_ms post_accounted_ms total_accounted_ms total_unaccounted_ms sampling_worker_ms list_worker_ms hmac payload_match telemetry_equivalent qualification exit_code command_elapsed_ms build84_elapsed_ms command_ratio_vs_build84 build84_geometry_ms geometry_ratio_vs_build84 gate_reason)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"

extract_re(){ local pattern="$1" file="$2" default_value="${3:-}" value; value="$(sed -nE "s/$pattern/\\1/p" "$file" | head -n1)"; [[ -n "$value" ]] && printf '%s' "$value" || printf '%s' "$default_value"; }
baseline_metric(){ local path="$1" image="$2" column="$3"; [[ -f "$path" ]] || { printf '%s' '-'; return; }; awk -F '\t' -v image="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==image&&col>0{print $col;found=1;exit}END{if(!found)print "-"}' "$path"; }
ratio(){ local base="$1" now="$2"; if [[ "$base" =~ ^[0-9]+$ && "$now" =~ ^[0-9]+$ && "$base" -gt 0 ]]; then awk -v b="$base" -v n="$now" 'BEGIN{printf "%.4f",n/b}'; else printf '%s' '-'; fi; }

for spec in "${specs[@]}"; do
  IFS='|' read -r file class role expected_payload <<< "$spec"
  input="$ACQUISITION_DIR/$file"; stem="${file%.*}"
  log="$stage_dir/logs/$stem.stderr.txt"; payload_file="$stage_dir/logs/$stem.payload.bin"
  start_ns="$(date +%s%N)"
  set +e
  timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-extract-phone -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -raw >"$payload_file" 2>"$log"
  rc=$?
  set -e
  end_ns="$(date +%s%N)"; command_elapsed_ms=$(( (end_ns-start_ns)/1000000 ))

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

  b68_physical=0; b68_spec=0; sampling_worker=0; list_worker=0
  if grep -q '^build68-plane-reuse:' "$log"; then
    b68_physical="$(extract_re 'build68-plane-reuse: .* physical-decode-candidates=([0-9]+).*' "$log" 0)"
    b68_spec="$(extract_re 'build68-plane-reuse: .* speculative-candidates=([0-9]+).*' "$log" 0)"
    sampling_worker="$(extract_re 'build68-plane-reuse: .* sampling-worker-ms=([0-9]+).*' "$log" 0)"
    list_worker="$(extract_re 'build68-plane-reuse: .* list-worker-ms=([0-9]+).*' "$log" 0)"
  fi

  b84=false
  grep -q '^build84-rgb-fetch:' "$log" && b84=true

  b89=false; total=0; geometry=0; geom_plane=0; freeze=0; prefix1=0; gen2=0; gen3=0; gen4=0; geom_accounted=0; geom_unaccounted=0; post_plane=0; qual_ms=0; decode_wall=0; post_accounted=0; total_accounted=0; total_unaccounted=0
  if grep -q '^build89-pipeline-profile:' "$log"; then
    b89=true
    total="$(extract_re 'build89-pipeline-profile: .* total-ms=([0-9]+).*' "$log" 0)"
    geometry="$(extract_re 'build89-pipeline-profile: .* geometry-ms=([0-9]+).*' "$log" 0)"
    geom_plane="$(extract_re 'build89-pipeline-profile: .* geometry-plane-prep-ms=([0-9]+).*' "$log" 0)"
    freeze="$(extract_re 'build89-pipeline-profile: .* freeze-ms=([0-9]+).*' "$log" 0)"
    prefix1="$(extract_re 'build89-pipeline-profile: .* prefix1-ms=([0-9]+).*' "$log" 0)"
    gen2="$(extract_re 'build89-pipeline-profile: .* gen2-ms=([0-9]+).*' "$log" 0)"
    gen3="$(extract_re 'build89-pipeline-profile: .* gen3-ms=([0-9]+).*' "$log" 0)"
    gen4="$(extract_re 'build89-pipeline-profile: .* gen4-ms=([0-9]+).*' "$log" 0)"
    geom_accounted="$(extract_re 'build89-pipeline-profile: .* geometry-accounted-ms=([0-9]+).*' "$log" 0)"
    geom_unaccounted="$(extract_re 'build89-pipeline-profile: .* geometry-unaccounted-ms=([0-9]+).*' "$log" 0)"
    post_plane="$(extract_re 'build89-pipeline-profile: .* post-plane-prep-ms=([0-9]+).*' "$log" 0)"
    qual_ms="$(extract_re 'build89-pipeline-profile: .* qualification-ms=([0-9]+).*' "$log" 0)"
    decode_wall="$(extract_re 'build89-pipeline-profile: .* decode-wall-ms=([0-9]+).*' "$log" 0)"
    post_accounted="$(extract_re 'build89-pipeline-profile: .* post-accounted-ms=([0-9]+).*' "$log" 0)"
    total_accounted="$(extract_re 'build89-pipeline-profile: .* total-accounted-ms=([0-9]+).*' "$log" 0)"
    total_unaccounted="$(extract_re 'build89-pipeline-profile: .* total-unaccounted-ms=([0-9]+).*' "$log" 0)"
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
    [[ "$b64_attempted" == true && "$b64_evals" -eq "$exp_evals" && "$b64_bank" -eq "$exp_bank" && "$b64_qualified" -eq "$exp_qual" && "$b64_decode" -eq "$exp_decode" && "$b64_frames" -eq "$exp_frames" && "$b84" == true && "$b89" == true ]] || telemetry=false
    [[ "$total" -ge "$geometry" && "$geometry" -ge "$geom_accounted" && "$geom_accounted" -ge "$gen4" && "$total" -ge "$total_accounted" ]] || telemetry=false
    [[ $((geom_accounted+geom_unaccounted)) -le $((geometry+5)) ]] || telemetry=false
    [[ $((total_accounted+total_unaccounted)) -le $((total+5)) ]] || telemetry=false
    if [[ "$exp_qual" -gt 0 ]]; then
      [[ "$b68_physical" -ge "$b64_decode" && "$b68_spec" -eq $((b68_physical-b64_decode)) ]] || telemetry=false
    else
      [[ "$b68_physical" -eq 0 && "$b68_spec" -eq 0 ]] || telemetry=false
    fi
  else
    [[ "$b64_attempted" == false && "$b84" == false && "$b89" == false ]] || telemetry=false
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
      required-build84) if [[ $rc -eq 0 && "$payload_match" == true && "$hmac" == true && "$b64_auth" == true && "$b84" == true && "$b89" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=build84-positive-semantics; fi;;
      required-reject) if [[ $rc -ne 0 && "$hmac" == false && "$b64_auth" == false && "$b84" == true && "$b89" == true ]]; then qualification=PASS; else qualification=FAIL; gate_reason=reject-semantics; fi;;
    esac
  fi
  if [[ "$telemetry" != true ]]; then
    qualification=FAIL
    if [[ "$gate_reason" == ok ]]; then gate_reason=telemetry; else gate_reason="${gate_reason}+telemetry"; fi
  fi

  base_elapsed="$(baseline_metric "$BUILD84_BASELINE_TSV" "$file" elapsed_ms)"
  base_geometry="$(baseline_metric "$BUILD84_BASELINE_TSV" "$file" geometry_ms)"
  elapsed_ratio="$(ratio "$base_elapsed" "$command_elapsed_ms")"
  geometry_ratio="$(ratio "$base_geometry" "$geometry")"
  row=("$file" "$class" "$role" "$b64_evals" "$b64_bank" "$b64_qualified" "$b64_decode" "$b64_frames" "$b64_auth" "$b68_physical" "$b68_spec" "$b84" "$b89" "$total" "$geometry" "$geom_plane" "$freeze" "$prefix1" "$gen2" "$gen3" "$gen4" "$geom_accounted" "$geom_unaccounted" "$post_plane" "$qual_ms" "$decode_wall" "$post_accounted" "$total_accounted" "$total_unaccounted" "$sampling_worker" "$list_worker" "$hmac" "$payload_match" "$telemetry" "$qualification" "$rc" "$command_elapsed_ms" "$base_elapsed" "$elapsed_ratio" "$base_geometry" "$geometry_ratio" "$gate_reason")
  if (( ${#row[@]} != ${#header[@]} )); then
    echo "error: Build89 TSV field mismatch for $file: header=${#header[@]} row=${#row[@]}" >&2
    exit 1
  fi
  (IFS=$'\t'; echo "${row[*]}") >> "$tsv"
done

{
  echo '# PixSeal Build89 full qualified Build84 deep-pipeline profile'; echo
  echo 'Build89 is observability-only. The active deep recovery path is the qualified Build84 runtime; Build89 derives phase accounting from timing telemetry Build84 already records and does not add probes inside geometry, qualification, sampling, list decode or HMAC loops.'; echo
  echo '| image | role | deep total ms | geometry | freeze | prefix1 | gen2 | gen3 | gen4 | post plane | qualification | decode wall | geom unacct | total unacct | command ms | semantic gate |'
  echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}{printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",$c["image"],$c["role"],$c["total_ms"],$c["geometry_ms"],$c["freeze_ms"],$c["prefix1_ms"],$c["gen2_ms"],$c["gen3_ms"],$c["gen4_ms"],$c["post_plane_prep_ms"],$c["qualification_ms"],$c["decode_wall_ms"],$c["geometry_unaccounted_ms"],$c["total_unaccounted_ms"],$c["command_elapsed_ms"],$c["qualification"]}' "$tsv"
  echo
  echo '## Deep-path phase shares'
  echo
  echo '| image | geometry/total | freeze/total | prefix1/total | gen2/total | gen3/total | gen4/total | qualification/total | decode/total |'
  echo '|---|---:|---:|---:|---:|---:|---:|---:|---:|'
  awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}$c["build89_attempted"]=="true"&&$c["total_ms"]+0>0{t=$c["total_ms"]+0;printf "| %s | %.2f%% | %.2f%% | %.2f%% | %.2f%% | %.2f%% | %.2f%% | %.2f%% | %.2f%% |\n",$c["image"],100*$c["geometry_ms"]/t,100*$c["freeze_ms"]/t,100*$c["prefix1_ms"]/t,100*$c["gen2_ms"]/t,100*$c["gen3_ms"]/t,100*$c["gen4_ms"]/t,100*$c["qualification_ms"]/t,100*$c["decode_wall_ms"]/t}' "$tsv"
  echo
  echo 'Worker-time counters for decode sampling/list work are preserved in the TSV but are intentionally not treated as wall-clock shares because parallel worker sums can exceed decode wall time.'
  echo 'Build84 remains the current qualified smartphone baseline. Build89 timing is profiling evidence only and cannot itself promote a new runtime.'
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

echo "Build89 staged result: controls=${control_passes}/3 A/front=${a_front_pass}/1 A/mild=${a_mild_pass}/1 A/angle=${a_angle_pass}/1 B/front=${b_front_pass}/1 B/mild=${b_mild_pass}/1 B/angle-reject=${b_angle_pass}/1"
if ((pass_count!=9 || control_passes!=3 || a_front_pass!=1 || a_mild_pass!=1 || a_angle_pass!=1 || b_front_pass!=1 || b_mild_pass!=1 || b_angle_pass!=1)); then
  echo 'Build89 failing cases:' >&2
  awk -F '\t' 'NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}$c["qualification"]!="PASS"{printf "  %s role=%s reason=%s rc=%s hmac=%s telemetry=%s evals=%s bank=%s qual=%s decode=%s frames=%s\n",$c["image"],$c["role"],$c["gate_reason"],$c["exit_code"],$c["hmac"],$c["telemetry_equivalent"],$c["build64_evals"],$c["build64_bank"],$c["build64_qualified"],$c["build64_decode"],$c["build64_frames"]}' "$tsv" >&2
  rm -rf "$failed_output_dir"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  echo "Build89 FAILED profile preserved in: $failed_output_dir" >&2
  echo 'error: Build89 semantic-equivalence profiling gate not met' >&2
  exit 1
fi

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build89 full qualified Build84 pipeline profile written to:'
echo "  $final_tsv"
echo "  $final_md"
echo "  $OUTPUT_DIR/build89-metadata.txt"
echo 'Build89 profiling gate: PASS'
