#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD83_PROFILE_DIR:-v4-phone private/build83-profile}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD83_BENCHTIME:-2s}"
COUNT="${V4_BUILD83_COUNT:-5}"
PPROF_TIME="${V4_BUILD83_PPROF_TIME:-20s}"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build83.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build83 projective sampler profiling"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "pprof_time=$PPROF_TIME"
  echo "runtime_path=Build76-qualified"
  echo "profile_workload=public-deterministic-no-key-no-payload"
} > "$stage_dir/build83-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build83PublicBenchmarkFixture$' -count=1 \
  | tee "$stage_dir/build83-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build83' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build83-benchmark.txt"

"${GO[@]}" test ./watermark -run '^$' \
  -bench '^BenchmarkExperimentalV4Build83HistoricalBlockReaderAngle$' \
  -benchtime="$PPROF_TIME" -count=1 \
  -outputdir="$stage_dir" -cpuprofile=build83-cpu.pprof \
  | tee "$stage_dir/build83-profile-run.txt"

if [[ ! -s "$stage_dir/build83-cpu.pprof" ]]; then
  echo "error: Build83 CPU profile was not created" >&2
  exit 1
fi

profile_binary="$stage_dir/watermark.test"
if [[ -x "$profile_binary" ]]; then
  "${GO[@]}" tool pprof -top -nodecount=40 "$profile_binary" "$stage_dir/build83-cpu.pprof" \
    > "$stage_dir/build83-pprof-top.txt"
else
  "${GO[@]}" tool pprof -top -nodecount=40 "$stage_dir/build83-cpu.pprof" \
    > "$stage_dir/build83-pprof-top.txt"
fi

"${GO[@]}" test ./watermark -run '^$' -gcflags='-m=2' 2>&1 \
  | grep -E 'samplePlaneLuminance|readProjectiveBlockValue|homography.mapPoint' \
  > "$stage_dir/build83-inlining.txt" || true

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build83 projective sampler profile written to:"
echo "  $OUTPUT_DIR/build83-benchmark.txt"
echo "  $OUTPUT_DIR/build83-cpu.pprof"
echo "  $OUTPUT_DIR/build83-pprof-top.txt"
echo "  $OUTPUT_DIR/build83-inlining.txt"
echo "Build83 profiling gate: PASS"
