#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD88_BENCH_DIR:-v4-phone private/build88-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD88_BENCHTIME:-2s}"
COUNT="${V4_BUILD88_COUNT:-5}"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build88.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build88 exact direct RGB loads / dominating BCE benchmark"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "qualified_baseline=Build84"
  echo "candidate=Build88-exact-direct-rgb-dominating-bce"
  echo "workload=public-deterministic-no-key-no-payload"
  echo "physical_gate=deferred-until-reader-speed-review"
} > "$stage_dir/build88-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build88PublicBenchmarkFixture$' -count=1 \
  | tee "$stage_dir/build88-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build88' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build88-benchmark.txt"

"${GO[@]}" test -a ./watermark -run '^$' -gcflags='-m=2 -d=ssa/check_bce/debug=1' \
  > "$stage_dir/build88-compiler.txt" 2>&1 || true

grep -E 'experimental_v4_phone_build84_rgb_fetch\.go:|experimental_v4_phone_build88_direct_rgb\.go:|experimentalV4PhoneBuild84ReadProjectiveBlockValue|experimentalV4PhoneBuild88ReadProjectiveBlockValue' \
  "$stage_dir/build88-compiler.txt" > "$stage_dir/build88-bce-inlining.txt" || true

{
  echo '# Build88 compiler BCE evidence'
  echo
  echo '## Qualified Build84 RGB access lines'
  awk -F: '$1=="watermark/experimental_v4_phone_build84_rgb_fetch.go" && $2+0>=131 && $2+0<=150 {print}' "$stage_dir/build88-compiler.txt" || true
  echo
  echo '## Build88 direct RGB access lines'
  awk -F: '$1=="watermark/experimental_v4_phone_build88_direct_rgb.go" && $2+0>=125 && $2+0<=175 {print}' "$stage_dir/build88-compiler.txt" || true
} > "$stage_dir/build88-bce-reader.txt"

"${GO[@]}" test -c ./watermark -o "$stage_dir/watermark.test"
"${GO[@]}" tool objdump -s 'experimentalV4PhoneBuild84ReadProjectiveBlockValue' "$stage_dir/watermark.test" \
  > "$stage_dir/build88-build84-reader-asm.txt" || true
"${GO[@]}" tool objdump -s 'experimentalV4PhoneBuild88ReadProjectiveBlockValue' "$stage_dir/watermark.test" \
  > "$stage_dir/build88-candidate-reader-asm.txt" || true

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build88 direct-RGB benchmark written to:"
echo "  $OUTPUT_DIR/build88-benchmark.txt"
echo "  $OUTPUT_DIR/build88-bce-reader.txt"
echo "  $OUTPUT_DIR/build88-bce-inlining.txt"
echo "  $OUTPUT_DIR/build88-build84-reader-asm.txt"
echo "  $OUTPUT_DIR/build88-candidate-reader-asm.txt"
echo "Build88 benchmark gate: PASS (correctness only; review whole-reader speed and BCE before physical testing)"
