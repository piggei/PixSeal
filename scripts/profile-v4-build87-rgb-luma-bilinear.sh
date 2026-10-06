#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD87_PROFILE_DIR:-v4-phone private/build87-profile}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD87_BENCHTIME:-2s}"
COUNT="${V4_BUILD87_COUNT:-5}"
PPROF_TIME="${V4_BUILD87_PPROF_TIME:-20s}"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build87.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build87 qualified RGB/luminance/bilinear compiler profiling"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "pprof_time=$PPROF_TIME"
  echo "runtime_path=Build84-qualified"
  echo "current_qualified_baseline=Build84"
  echo "build86_status=exact-benchmark-negative-non-promoted"
  echo "profile_workload=public-deterministic-no-key-no-payload"
  echo "build87_status=observability-only-no-runtime-algorithm-change"
} > "$stage_dir/build87-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build87PublicProfileFixture$' -count=1 \
  | tee "$stage_dir/build87-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build87' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build87-benchmark.txt"

"${GO[@]}" test ./watermark -run '^$' \
  -bench '^BenchmarkExperimentalV4Build87QualifiedBlockReaderAngle$' \
  -benchtime="$PPROF_TIME" -count=1 \
  -outputdir="$stage_dir" -cpuprofile=build87-cpu.pprof \
  | tee "$stage_dir/build87-profile-run.txt"

if [[ ! -s "$stage_dir/build87-cpu.pprof" ]]; then
  echo "error: Build87 CPU profile was not created" >&2
  exit 1
fi

profile_binary="$stage_dir/watermark.test"
if [[ -x "$profile_binary" ]]; then
  "${GO[@]}" tool pprof -top -nodecount=60 "$profile_binary" "$stage_dir/build87-cpu.pprof" \
    > "$stage_dir/build87-pprof-top.txt"
  "${GO[@]}" tool pprof -list='experimentalV4PhoneBuild84ReadProjectiveBlockValue' "$profile_binary" "$stage_dir/build87-cpu.pprof" \
    > "$stage_dir/build87-pprof-reader-list.txt" || true
  "${GO[@]}" tool objdump -s 'experimentalV4PhoneBuild84ReadProjectiveBlockValue' "$profile_binary" \
    > "$stage_dir/build87-reader-asm.txt" || true
else
  "${GO[@]}" tool pprof -top -nodecount=60 "$stage_dir/build87-cpu.pprof" \
    > "$stage_dir/build87-pprof-top.txt"
  "${GO[@]}" tool pprof -list='experimentalV4PhoneBuild84ReadProjectiveBlockValue' "$stage_dir/build87-cpu.pprof" \
    > "$stage_dir/build87-pprof-reader-list.txt" || true
  "${GO[@]}" test -c ./watermark -o "$stage_dir/watermark.test"
  profile_binary="$stage_dir/watermark.test"
  "${GO[@]}" tool objdump -s 'experimentalV4PhoneBuild84ReadProjectiveBlockValue' "$profile_binary" \
    > "$stage_dir/build87-reader-asm.txt" || true
fi

awk '/experimental_v4_phone_build84_rgb_fetch\.go:(140|141|142|143|144|145|146|147|148|149|150)/' \
  "$stage_dir/build87-reader-asm.txt" > "$stage_dir/build87-reader-asm-rgb.txt" || true

"${GO[@]}" test -a ./watermark -run '^$' -gcflags='-m=2 -d=ssa/check_bce/debug=1' \
  > "$stage_dir/build87-compiler.txt" 2>&1 || true

grep -E 'experimental_v4_phone_build84_rgb_fetch\.go:|experimentalV4PhoneBuild84ReadProjectiveBlockValue|homography.mapPoint' \
  "$stage_dir/build87-compiler.txt" > "$stage_dir/build87-compiler-reader.txt" || true

{
  echo '# Build87 RGB/luminance/bilinear BCE evidence'
  echo
  echo '## Qualified Build84 reader lines 131-150'
  awk -F: '$1=="watermark/experimental_v4_phone_build84_rgb_fetch.go" && $2+0>=131 && $2+0<=150 {print}' \
    "$stage_dir/build87-compiler.txt" || true
} > "$stage_dir/build87-bce-rgb.txt"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build87 qualified RGB/luminance/bilinear profile written to:"
echo "  $OUTPUT_DIR/build87-benchmark.txt"
echo "  $OUTPUT_DIR/build87-cpu.pprof"
echo "  $OUTPUT_DIR/build87-pprof-top.txt"
echo "  $OUTPUT_DIR/build87-pprof-reader-list.txt"
echo "  $OUTPUT_DIR/build87-compiler-reader.txt"
echo "  $OUTPUT_DIR/build87-bce-rgb.txt"
echo "  $OUTPUT_DIR/build87-reader-asm-rgb.txt"
echo "Build87 profiling gate: PASS"
