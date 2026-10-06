#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD84_BENCH_DIR:-v4-phone private/build84-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD84_BENCHTIME:-2s}"
COUNT="${V4_BUILD84_COUNT:-5}"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build84.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build84 exact RGB-fetch/BCE benchmark"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "runtime_under_test=Build84-qualified"
  echo "historical_comparison_baseline=Build76"
  echo "current_qualified_baseline=Build84"
  echo "benchmark_workload=public-deterministic-no-key-no-payload"
  echo "qualification_status=Build84 promoted after two independent 9/9 physical PASS runs; benchmark performance is not a correctness gate"
} > "$stage_dir/build84-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build84PublicBenchmarkFixture$' -count=1 \
  | tee "$stage_dir/build84-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build84' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build84-benchmark.txt"

# Preserve compiler inliner and bounds-check evidence. check_bce output is
# diagnostic compiler evidence; no count is treated as semantic correctness.
"${GO[@]}" test -a ./watermark -run '^$' -gcflags='-m=2 -d=ssa/check_bce/debug=1' \
  > "$stage_dir/build84-compiler.txt" 2>&1 || true

grep -E 'experimental_v4_phone_build8(2_exact_inline_sampler|4_rgb_fetch)\\.go:|samplePlaneLuminance|readProjectiveBlockValue|homography.mapPoint' \
  "$stage_dir/build84-compiler.txt" > "$stage_dir/build84-bce-inlining.txt" || true

reader_section() {
  local file="$1" func="$2" label="$3" start end
  start="$(grep -n "^func ${func}" "$file" | head -n1 | cut -d: -f1)"
  [[ -n "$start" ]] || return 0
  end="$(awk -v s="$start" 'NR>s && /^}$/ {print NR; exit}' "$file")"
  echo "## $label lines $start-$end"
  awk -F: -v f="$file" -v s="$start" -v e="$end" '$1==f && $2+0>=s && $2+0<=e {print}' "$stage_dir/build84-compiler.txt"
}
{
  echo '# Build84 compiler BCE evidence'
  echo
  reader_section 'watermark/experimental_v4_phone_build82_exact_inline_sampler.go' 'experimentalV4PhoneBuild82ReadProjectiveBlockValue' 'Build82 direct RGB indexing'
  echo
  reader_section 'watermark/experimental_v4_phone_build84_rgb_fetch.go' 'experimentalV4PhoneBuild84ReadProjectiveBlockValue' 'Build84 three-byte pixel slices'
} > "$stage_dir/build84-bce-reader.txt"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build84 RGB-fetch benchmark written to:"
echo "  $OUTPUT_DIR/build84-benchmark.txt"
echo "  $OUTPUT_DIR/build84-bce-reader.txt"
echo "  $OUTPUT_DIR/build84-bce-inlining.txt"
echo "Build84 benchmark gate: PASS (correctness only; review speed before physical testing)"
