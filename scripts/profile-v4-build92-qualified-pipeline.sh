#!/usr/bin/env bash
set -euo pipefail

ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD92_PROFILE_DIR:-v4-phone private/build92-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PROFILE_TIMEOUT="${V4_PHONE_BUILD92_TIMEOUT:-7200}"
COOLDOWN_MS="${V4_PHONE_BUILD92_COOLDOWN_MS:-5000}"
GO_TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-mild.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-mild.jpg" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-angle.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-angle.jpg" >&2; exit 1; }
[[ "$PROFILE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD92_TIMEOUT must be an integer number of seconds" >&2; exit 1; }
[[ "$COOLDOWN_MS" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD92_COOLDOWN_MS must be a non-negative integer" >&2; exit 1; }

# Package tests may execute from a package-local working directory. Always pass
# the retained corpus and profile outputs as absolute paths.
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build92.XXXXXX")"
stage_dir="$(cd "$stage_dir" && pwd -P)"
published=false
cleanup_stage(){ if [[ "$published" != true && -d "$stage_dir" ]]; then rm -rf "$stage_dir"; fi; }
trap cleanup_stage EXIT
failed_output_dir="${OUTPUT_DIR}-failed"

testbin="$stage_dir/build92-watermark.test"
timing_raw="$stage_dir/build92-timing-raw.txt"
tsv="$stage_dir/build92-pipeline.tsv"
md="$stage_dir/build92-pipeline.md"
meta="$stage_dir/build92-metadata.txt"

{
  echo "PixSeal Build92 qualified Build90 post-promotion pipeline profile"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go_toolchain=$GO_TOOLCHAIN"
  echo "host=$(uname -a)"
  echo "go_version=$(env GOTOOLCHAIN="$GO_TOOLCHAIN" go version 2>/dev/null || true)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "acquisition_dir=$ACQUISITION_DIR"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "cooldown_ms=$COOLDOWN_MS"
  echo "timing_sequence=B/mild then B/angle; CPU profiles are separate later passes"
  echo "runtime=unchanged qualified Build90"
  echo "purpose=post-promotion measurement only; select any Build93 from Build90-qualified evidence"
} > "$meta"

# Build one stable test binary and retain it beside the profiles so pprof remains
# symbolizable even after the working tree moves on.
env GOTOOLCHAIN="$GO_TOOLCHAIN" go test -c ./watermark -o "$testbin" 2>&1 | tee "$stage_dir/build92-test-binary-build.txt"

run_env=(
  "PIXSEAL_BUILD92_ACQUISITION_DIR=$ACQUISITION_DIR"
  "PIXSEAL_BUILD92_KEY=$KEY"
  "PIXSEAL_BUILD92_WIDTH=$CANONICAL_WIDTH"
  "PIXSEAL_BUILD92_HEIGHT=$CANONICAL_HEIGHT"
  "PIXSEAL_BUILD92_MESSAGE_B=$MESSAGE_B"
  "PIXSEAL_BUILD92_COOLDOWN_MS=$COOLDOWN_MS"
)

set +e
timeout --foreground "${PROFILE_TIMEOUT}s" env "${run_env[@]}" \
  "$testbin" -test.run '^TestExperimentalV4Build92RetainedTimingProfile$' -test.v -test.timeout=0 2>&1 | tee "$timing_raw"
rc=${PIPESTATUS[0]}
set -e
if (( rc != 0 )); then
  rm -rf "$failed_output_dir"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  echo "Build92 FAILED timing profile preserved in: $failed_output_dir" >&2
  exit "$rc"
fi

header=(image outcome prepare_ms wall_ms total_ms geometry_ms geometry_plane_ms freeze_ms prefix1_ms gen2_ms gen3_ms gen4_ms prefix1_worker_ms gen2_worker_ms gen3_worker_ms gen4_worker_ms geometry_accounted_ms geometry_unaccounted_ms post_plane_ms qualification_ms decode_ms post_accounted_ms total_accounted_ms total_unaccounted_ms geometry_evals bank qualification_evals qualified decode_candidates list_frames physical_decode speculative_decode authenticated prefix1_workers gen2_workers gen3_workers gen4_workers gen2_tasks gen3_tasks gen4_tasks qualification_workers qualification_tasks decode_workers fold_scores block_reads block_success block_failed)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"
grep $'^build92-row\t' "$timing_raw" | cut -f2- >> "$tsv" || true
rows=$(( $(wc -l < "$tsv") - 1 ))
if (( rows != 2 )); then
  echo "error: Build92 expected 2 timing rows, found $rows" >&2
  rm -rf "$failed_output_dir"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  exit 1
fi
awk -F '\t' -v expected_fields="${#header[@]}" '
NR==1 { next }
NF != expected_fields { printf "error: Build92 TSV field mismatch on row %d: got %d want %d\n", NR, NF, expected_fields > "/dev/stderr"; bad=1 }
END { exit bad }
' "$tsv"

for image_name in phone-b-mild.jpg phone-b-angle.jpg; do
  count="$(awk -F '\t' -v im="$image_name" 'NR>1&&$1==im{n++}END{print n+0}' "$tsv")"
  [[ "$count" -eq 1 ]] || { echo "error: Build92 missing unique timing row for $image_name" >&2; exit 1; }
done

run_cpu_profile(){
  local image_name="$1"
  local stem="$2"
  local profile="$stage_dir/${stem}.pprof"
  local raw="$stage_dir/${stem}-raw.txt"
  set +e
  timeout --foreground "${PROFILE_TIMEOUT}s" env "${run_env[@]}" \
    "PIXSEAL_BUILD92_CASE=$image_name" \
    "PIXSEAL_BUILD92_CPU_PROFILE=$profile" \
    "$testbin" -test.run '^TestExperimentalV4Build92CPUProfile$' -test.v -test.timeout=0 2>&1 | tee "$raw"
  local rc=${PIPESTATUS[0]}
  set -e
  if (( rc != 0 )); then
    rm -rf "$failed_output_dir"
    mv "$stage_dir" "$failed_output_dir"
    published=true
    echo "Build92 FAILED CPU profile preserved in: $failed_output_dir" >&2
    exit "$rc"
  fi
  [[ -s "$profile" ]] || { echo "error: Build92 missing CPU profile $profile" >&2; exit 1; }

  env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -nodecount=80 "$testbin" "$profile" > "$stage_dir/${stem}-pprof-top.txt" 2>&1
  env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -cum -nodecount=80 "$testbin" "$profile" > "$stage_dir/${stem}-pprof-cum.txt" 2>&1

  for phase in geometry qualification decode; do
    case "$phase" in
      geometry) focus='experimentalV4PhoneBuild84BlindBank' ;;
      qualification) focus='experimentalV4PhoneBuild68Qualify' ;;
      decode) focus='experimentalV4PhoneBuild67DecodeSingle' ;;
    esac
    env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -nodecount=80 -focus="$focus" "$testbin" "$profile" > "$stage_dir/${stem}-pprof-${phase}.txt" 2>&1 || true
  done

  {
    for fn in experimentalV4PhoneBuild90Recover experimentalV4PhoneBuild84BlindBank experimentalV4PhoneBuild68Qualify experimentalV4PhoneBuild67DecodeSingle experimentalV4PhoneBuild66ConsumeBatch experimentalV4PhoneBuild84ReadProjectiveBlockValue; do
      echo "===== $fn ====="
      env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -list="$fn" "$testbin" "$profile" 2>&1 || true
      echo
    done
  } > "$stage_dir/${stem}-pprof-list.txt"
}

run_cpu_profile phone-b-mild.jpg build92-b-mild
run_cpu_profile phone-b-angle.jpg build92-b-angle

cell(){ local image="$1" column="$2"; awk -F '\t' -v im="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==im{print $col;exit}' "$tsv"; }
pct(){ local part="$1" whole="$2"; if [[ "$part" =~ ^[0-9]+$ && "$whole" =~ ^[0-9]+$ && "$whole" -gt 0 ]]; then awk -v p="$part" -v w="$whole" 'BEGIN{printf "%.2f",100*p/w}'; else printf '-'; fi; }

{
  echo '# PixSeal Build92 qualified Build90 post-promotion pipeline profile'; echo
  echo 'Build92 is observability-only. The timing pass executes the unchanged qualified Build90 runtime without CPU profiling. Separate later passes collect CPU pprof for B/mild and B/angle so profiler overhead cannot contaminate the authoritative stage-share timing row.'; echo
  echo '| image | total | geometry | geometry share | freeze | prefix1 | gen2 | gen3 | gen4 | qualification | qualification share | decode | decode share | q workers/tasks | decode workers |'
  echo '|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  for image in phone-b-mild.jpg phone-b-angle.jpg; do
    total="$(cell "$image" total_ms)"; geometry="$(cell "$image" geometry_ms)"; qualification="$(cell "$image" qualification_ms)"; decode="$(cell "$image" decode_ms)"
    echo "| $image | $total | $geometry | $(pct "$geometry" "$total")% | $(cell "$image" freeze_ms) | $(cell "$image" prefix1_ms) | $(cell "$image" gen2_ms) | $(cell "$image" gen3_ms) | $(cell "$image" gen4_ms) | $qualification | $(pct "$qualification" "$total")% | $decode | $(pct "$decode" "$total")% | $(cell "$image" qualification_workers)/$(cell "$image" qualification_tasks) | $(cell "$image" decode_workers) |"
  done
  echo
  echo 'CPU profiles are separate diagnostic passes. Review flat/cumulative top reports together with phase-focused reports before selecting a Build93 hypothesis. Timing from the CPU-profile passes is not promotion evidence.'; echo
  echo 'Semantic invariants are enforced by the Go tests using the frozen Build91 retained-corpus expectations. No key/payload/ECC/HMAC result selects geometry, qualification work, profiling scope or pprof focus.'; echo
} > "$md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build92 qualified Build90 post-promotion profile written to:'
echo "  $OUTPUT_DIR/build92-pipeline.tsv"
echo "  $OUTPUT_DIR/build92-pipeline.md"
echo "  $OUTPUT_DIR/build92-b-mild.pprof"
echo "  $OUTPUT_DIR/build92-b-angle.pprof"
echo "  $OUTPUT_DIR/build92-b-mild-pprof-top.txt / -cum.txt / phase-focused reports / -list.txt"
echo "  $OUTPUT_DIR/build92-b-angle-pprof-top.txt / -cum.txt / phase-focused reports / -list.txt"
echo "  $OUTPUT_DIR/build92-metadata.txt"
echo 'Build92 profiling gate: PASS'
