#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD85_PROFILE_DIR:-v4-phone private/build85-profile}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD85_BENCHTIME:-2s}"
COUNT="${V4_BUILD85_COUNT:-5}"
PPROF_TIME="${V4_BUILD85_PPROF_TIME:-20s}"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build85.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build85 qualified-sampler post-promotion profiling"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "pprof_time=$PPROF_TIME"
  echo "runtime_path=Build84-qualified"
  echo "current_qualified_baseline=Build84"
  echo "profile_workload=public-deterministic-no-key-no-payload"
  echo "build85_status=observability-only-no-runtime-algorithm-change"
} > "$stage_dir/build85-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build85PublicProfileFixture$' -count=1 \
  | tee "$stage_dir/build85-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build85' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build85-benchmark.txt"

"${GO[@]}" test ./watermark -run '^$' \
  -bench '^BenchmarkExperimentalV4Build85QualifiedBlockReaderAngle$' \
  -benchtime="$PPROF_TIME" -count=1 \
  -outputdir="$stage_dir" -cpuprofile=build85-cpu.pprof \
  | tee "$stage_dir/build85-profile-run.txt"

if [[ ! -s "$stage_dir/build85-cpu.pprof" ]]; then
  echo "error: Build85 CPU profile was not created" >&2
  exit 1
fi

profile_binary="$stage_dir/watermark.test"
if [[ -x "$profile_binary" ]]; then
  "${GO[@]}" tool pprof -top -nodecount=50 "$profile_binary" "$stage_dir/build85-cpu.pprof" \
    > "$stage_dir/build85-pprof-top.txt"
  "${GO[@]}" tool pprof -list='experimentalV4PhoneBuild84ReadProjectiveBlockValue' "$profile_binary" "$stage_dir/build85-cpu.pprof" \
    > "$stage_dir/build85-pprof-reader-list.txt" || true
else
  "${GO[@]}" tool pprof -top -nodecount=50 "$stage_dir/build85-cpu.pprof" \
    > "$stage_dir/build85-pprof-top.txt"
  "${GO[@]}" tool pprof -list='experimentalV4PhoneBuild84ReadProjectiveBlockValue' "$stage_dir/build85-cpu.pprof" \
    > "$stage_dir/build85-pprof-reader-list.txt" || true
fi

"${GO[@]}" test -a ./watermark -run '^$' -gcflags='-m=2 -d=ssa/check_bce/debug=1' \
  > "$stage_dir/build85-compiler.txt" 2>&1 || true

grep -E 'experimental_v4_phone_build84_rgb_fetch\.go:|experimentalV4PhoneBuild84ReadProjectiveBlockValue|homography.mapPoint' \
  "$stage_dir/build85-compiler.txt" > "$stage_dir/build85-compiler-reader.txt" || true

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build85 qualified-sampler profile written to:"
echo "  $OUTPUT_DIR/build85-benchmark.txt"
echo "  $OUTPUT_DIR/build85-cpu.pprof"
echo "  $OUTPUT_DIR/build85-pprof-top.txt"
echo "  $OUTPUT_DIR/build85-pprof-reader-list.txt"
echo "  $OUTPUT_DIR/build85-compiler-reader.txt"
echo "Build85 profiling gate: PASS"
