#!/usr/bin/env bash
set -euo pipefail

ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD94_PROFILE_DIR:-v4-phone private/build94-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PROFILE_TIMEOUT="${V4_PHONE_BUILD94_TIMEOUT:-7200}"
COOLDOWN_MS="${V4_PHONE_BUILD94_COOLDOWN_MS:-5000}"
GO_TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-mild.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-mild.jpg" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-angle.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-angle.jpg" >&2; exit 1; }
[[ "$PROFILE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD94_TIMEOUT must be an integer number of seconds" >&2; exit 1; }
[[ "$COOLDOWN_MS" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD94_COOLDOWN_MS must be a non-negative integer" >&2; exit 1; }

# Package tests may execute from a package-local working directory. Always pass
# the retained corpus and profile outputs as absolute paths.
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build94.XXXXXX")"
stage_dir="$(cd "$stage_dir" && pwd -P)"
published=false
failed_output_dir="${OUTPUT_DIR}-failed"
cleanup_stage(){
  local rc=$?
  trap - EXIT
  set +e
  if [[ "$published" != true && -d "$stage_dir" ]]; then
    if (( rc != 0 )); then
      rm -rf "$failed_output_dir"
      mv "$stage_dir" "$failed_output_dir"
      echo "Build94 FAILED staged evidence preserved in: $failed_output_dir" >&2
    else
      rm -rf "$stage_dir"
    fi
  fi
  exit "$rc"
}
trap cleanup_stage EXIT

testbin="$stage_dir/build94-watermark.test"
timing_raw="$stage_dir/build94-timing-raw.txt"
tsv="$stage_dir/build94-pipeline.tsv"
md="$stage_dir/build94-pipeline.md"
meta="$stage_dir/build94-metadata.txt"

build93_candidate_sha256="$(sha256sum watermark/experimental_v4_phone_build93_candidate.go | awk '{print $1}')"
build93_reuse_sha256="$(sha256sum watermark/experimental_v4_phone_build93_residue_reuse.go | awk '{print $1}')"
build90_qualification_sha256="$(sha256sum watermark/experimental_v4_phone_build90_parallel_qualification.go | awk '{print $1}')"
build84_reader_sha256="$(sha256sum watermark/experimental_v4_phone_build84_rgb_fetch.go | awk '{print $1}')"

{
  echo "PixSeal Build94 qualified Build93 post-promotion pipeline profile"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go_toolchain=$GO_TOOLCHAIN"
  echo "host=$(uname -a)"
  echo "go_version=$(env GOTOOLCHAIN="$GO_TOOLCHAIN" go version 2>/dev/null || true)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "acquisition_dir=$ACQUISITION_DIR"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "cooldown_ms=$COOLDOWN_MS"
  echo "timing_sequence=B/mild then B/angle; CPU profiles are separate later passes"
  echo "profile_build=Build94"
  echo "runtime=unchanged current qualified Build93"
  echo "qualified_baseline=Build93"
  echo "previous_qualified_baseline=Build90"
  echo "purpose=post-promotion observability only; select any Build95 from Build93-qualified evidence"
  echo "build93_candidate_sha256=$build93_candidate_sha256"
  echo "build93_residue_reuse_sha256=$build93_reuse_sha256"
  echo "build90_qualification_sha256=$build90_qualification_sha256"
  echo "build84_reader_sha256=$build84_reader_sha256"
} > "$meta"

# Build one stable test binary and retain it beside the profiles so pprof remains
# symbolizable even after the working tree moves on.
env GOTOOLCHAIN="$GO_TOOLCHAIN" go test -c ./watermark -o "$testbin" 2>&1 | tee "$stage_dir/build94-test-binary-build.txt"
echo "test_binary_sha256=$(sha256sum "$testbin" | awk '{print $1}')" >> "$meta"

run_env=(
  "PIXSEAL_BUILD94_ACQUISITION_DIR=$ACQUISITION_DIR"
  "PIXSEAL_BUILD94_KEY=$KEY"
  "PIXSEAL_BUILD94_WIDTH=$CANONICAL_WIDTH"
  "PIXSEAL_BUILD94_HEIGHT=$CANONICAL_HEIGHT"
  "PIXSEAL_BUILD94_MESSAGE_B=$MESSAGE_B"
  "PIXSEAL_BUILD94_COOLDOWN_MS=$COOLDOWN_MS"
)

set +e
timeout --foreground "${PROFILE_TIMEOUT}s" env "${run_env[@]}" \
  "$testbin" -test.run '^TestExperimentalV4Build94RetainedTimingProfile$' -test.v -test.timeout=0 2>&1 | tee "$timing_raw"
rc=${PIPESTATUS[0]}
set -e
if (( rc != 0 )); then
  exit "$rc"
fi

header_rows="$(awk -F '\t' '$1 == "build94-header" { n++ } END { print n+0 }' "$timing_raw")"
if (( header_rows != 1 )); then
  echo "error: Build94 expected one timing header, found $header_rows" >&2
  exit 1
fi
grep $'^build94-header\t' "$timing_raw" | cut -f2- > "$tsv"
grep $'^build94-row\t' "$timing_raw" | cut -f2- >> "$tsv" || true
rows=$(( $(wc -l < "$tsv") - 1 ))
if (( rows != 2 )); then
  echo "error: Build94 expected 2 timing rows, found $rows" >&2
  exit 1
fi
expected_fields="$(awk -F '\t' 'NR == 1 { print NF; exit }' "$tsv")"
awk -F '\t' -v expected_fields="$expected_fields" '
NR==1 { next }
NF != expected_fields { printf "error: Build94 TSV field mismatch on row %d: got %d want %d\n", NR, NF, expected_fields > "/dev/stderr"; bad=1 }
END { exit bad }
' "$tsv"

for image_name in phone-b-mild.jpg phone-b-angle.jpg; do
  count="$(awk -F '\t' -v im="$image_name" 'NR>1&&$1==im{n++}END{print n+0}' "$tsv")"
  [[ "$count" -eq 1 ]] || { echo "error: Build94 missing unique timing row for $image_name" >&2; exit 1; }
done

run_cpu_profile(){
  local image_name="$1"
  local stem="$2"
  local profile="$stage_dir/${stem}.pprof"
  local raw="$stage_dir/${stem}-raw.txt"
  set +e
  timeout --foreground "${PROFILE_TIMEOUT}s" env "${run_env[@]}" \
    "PIXSEAL_BUILD94_CASE=$image_name" \
    "PIXSEAL_BUILD94_CPU_PROFILE=$profile" \
    "$testbin" -test.run '^TestExperimentalV4Build94CPUProfile$' -test.v -test.timeout=0 2>&1 | tee "$raw"
  local rc=${PIPESTATUS[0]}
  set -e
  if (( rc != 0 )); then
    exit "$rc"
  fi
  [[ -s "$profile" ]] || { echo "error: Build94 missing CPU profile $profile" >&2; exit 1; }

  env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -nodecount=100 "$testbin" "$profile" > "$stage_dir/${stem}-pprof-top.txt" 2>&1
  env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -cum -nodecount=100 "$testbin" "$profile" > "$stage_dir/${stem}-pprof-cum.txt" 2>&1

  while IFS='|' read -r phase focus; do
    env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -nodecount=100 -focus="$focus" "$testbin" "$profile" > "$stage_dir/${stem}-pprof-${phase}.txt" 2>&1 || true
  done <<'PHASES'
geometry|experimentalV4PhoneBuild84BlindBank
gen3|experimentalV4PhoneBuild73Generation3
gen4|experimentalV4PhoneBuild84Generation4Fetch
continuation4|experimentalV4PhoneBuild84ContinueFetch
foldscore|experimentalV4PhoneBuild41FoldScore
qualification|experimentalV4PhoneBuild68Qualify
decode|experimentalV4PhoneBuild93DecodeSingle
protected-sampling|experimentalV4PhoneBuild93SampleProtectedResidues
profile-derivation|experimentalV4PhoneBuild93MarginsFromResidues
list-decode|experimentalV4PhoneBuild42ListDecode
PHASES

  {
    for fn in \
      experimentalV4PhoneBuild93Recover \
      experimentalV4PhoneBuild84BlindBank \
      experimentalV4PhoneBuild73Generation3 \
      experimentalV4PhoneBuild84Generation4Fetch \
      experimentalV4PhoneBuild84ContinueFetch \
      experimentalV4PhoneBuild41FoldScore \
      experimentalV4PhoneBuild84ReadProjectiveBlockValue \
      readProjectiveBlockValue \
      samplePlaneLuminance \
      experimentalV4PhoneBuild90EvaluateQualification \
      experimentalV4PhoneBuild68Qualify \
      experimentalV4PhoneBuild93DecodeSingle \
      experimentalV4PhoneBuild93SampleProtectedResidues \
      experimentalV4PhoneBuild93MarginsFromResidues \
      experimentalV4PhoneBuild42ListDecode \
      experimentalV4PhoneBuild66ConsumeBatch; do
      echo "===== $fn ====="
      env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -list="$fn" "$testbin" "$profile" 2>&1 || true
      echo
    done
  } > "$stage_dir/${stem}-pprof-list.txt"
}

run_cpu_profile phone-b-mild.jpg build94-b-mild
run_cpu_profile phone-b-angle.jpg build94-b-angle

cell(){ local image="$1" column="$2"; awk -F '\t' -v im="$image" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==im{print $col;exit}' "$tsv"; }
pct(){ local part="$1" whole="$2"; if [[ "$part" =~ ^[0-9]+$ && "$whole" =~ ^[0-9]+$ && "$whole" -gt 0 ]]; then awk -v p="$part" -v w="$whole" 'BEGIN{printf "%.2f",100*p/w}'; else printf '-'; fi; }
ratio(){ local num="$1" den="$2"; if [[ "$num" =~ ^[0-9]+$ && "$den" =~ ^[0-9]+$ && "$den" -gt 0 ]]; then awk -v n="$num" -v d="$den" 'BEGIN{printf "%.3f",n/d}'; else printf '-'; fi; }

{
  echo '# PixSeal Build94 qualified Build93 post-promotion pipeline profile'; echo
  echo 'Build94 is observability-only. The timing pass executes the unchanged, integration-closed Build93 deep runtime without CPU profiling. Separate later passes collect CPU pprof for B/mild and B/angle so profiler overhead cannot contaminate the authoritative stage-share rows.'; echo
  echo '| image | total | geometry | geometry share | freeze | prefix1 | gen2 | gen3 | gen4 | qualification | qualification share | decode | decode share | q workers/tasks | decode workers |'
  echo '|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  for image in phone-b-mild.jpg phone-b-angle.jpg; do
    total="$(cell "$image" total_ms)"; geometry="$(cell "$image" geometry_ms)"; qualification="$(cell "$image" qualification_ms)"; decode="$(cell "$image" decode_ms)"
    echo "| $image | $total | $geometry | $(pct "$geometry" "$total")% | $(cell "$image" freeze_ms) | $(cell "$image" prefix1_ms) | $(cell "$image" gen2_ms) | $(cell "$image" gen3_ms) | $(cell "$image" gen4_ms) | $qualification | $(pct "$qualification" "$total")% | $decode | $(pct "$decode" "$total")% | $(cell "$image" qualification_workers)/$(cell "$image" qualification_tasks) | $(cell "$image" decode_workers) |"
  done
  echo
  echo '| image | physical decode | physical profiles | physical list frames | protected sampling passes | profile derivations | derivations/pass | sampling worker ms | list worker ms |'
  echo '|---|---:|---:|---:|---:|---:|---:|---:|---:|'
  for image in phone-b-mild.jpg phone-b-angle.jpg; do
    passes="$(cell "$image" protected_sampling_passes)"; derivations="$(cell "$image" protected_profile_derivations)"
    echo "| $image | $(cell "$image" physical_decode) | $(cell "$image" physical_profiles) | $(cell "$image" physical_list_frames) | $passes | $derivations | $(ratio "$derivations" "$passes") | $(cell "$image" sampling_worker_ms) | $(cell "$image" list_worker_ms) |"
  done
  echo
  echo 'Worker-time sums are cumulative across parallel decode workers and must not be interpreted as wall-time shares. CPU profiles are separate diagnostic passes; their wall times are not promotion evidence.'; echo
  echo 'Semantic invariants are enforced using the frozen retained-corpus expectations plus Build93-specific physical sampling/profile-derivation counters. No key, payload, ECC, HMAC result, timing value or pprof focus selects geometry, qualification work or recovery order.'; echo
} > "$md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build94 qualified Build93 post-promotion profile written to:'
echo "  $OUTPUT_DIR/build94-pipeline.tsv"
echo "  $OUTPUT_DIR/build94-pipeline.md"
echo "  $OUTPUT_DIR/build94-b-mild.pprof"
echo "  $OUTPUT_DIR/build94-b-angle.pprof"
echo "  $OUTPUT_DIR/build94-b-mild-pprof-top.txt / -cum.txt / phase-focused reports / -list.txt"
echo "  $OUTPUT_DIR/build94-b-angle-pprof-top.txt / -cum.txt / phase-focused reports / -list.txt"
echo "  $OUTPUT_DIR/build94-metadata.txt"
echo 'Build94 profiling gate: PASS'
