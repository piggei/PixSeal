#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD86_BENCH_DIR:-v4-phone private/build86-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD86_BENCHTIME:-2s}"
COUNT="${V4_BUILD86_COUNT:-5}"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build86.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build86 exact DCT table-hoist/BCE benchmark"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "candidate=Build86-exact-dct-table-hoist"
  echo "current_qualified_baseline=Build84"
  echo "benchmark_workload=public-deterministic-no-key-no-payload"
  echo "decision_rule=benchmark-first; no physical test unless whole-reader speedup is stable and useful"
} > "$stage_dir/build86-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build86PublicBenchmarkFixture$' -count=1 \
  | tee "$stage_dir/build86-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build86' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build86-benchmark.txt"

# Compiler inliner/BCE evidence is diagnostic only. It is not a semantic gate.
"${GO[@]}" test -a ./watermark -run '^$' -gcflags='-m=2 -d=ssa/check_bce/debug=1' \
  > "$stage_dir/build86-compiler.txt" 2>&1 || true

grep -E 'experimental_v4_phone_build8(4_rgb_fetch|6_dct_hoist)\.go:|experimentalV4PhoneBuild(84|86)ReadProjectiveBlockValue|homography.mapPoint' \
  "$stage_dir/build86-compiler.txt" > "$stage_dir/build86-bce-inlining.txt" || true

reader_section() {
  local file="$1" func="$2" label="$3" start end
  start="$(grep -n "^func ${func}" "$file" | head -n1 | cut -d: -f1)"
  [[ -n "$start" ]] || return 0
  end="$(awk -v s="$start" 'NR>s && /^}$/ {print NR; exit}' "$file")"
  echo "## $label lines $start-$end"
  awk -F: -v f="$file" -v s="$start" -v e="$end" '$1==f && $2+0>=s && $2+0<=e {print}' "$stage_dir/build86-compiler.txt"
}
{
  echo '# Build86 compiler BCE evidence'
  echo
  reader_section 'watermark/experimental_v4_phone_build84_rgb_fetch.go' 'experimentalV4PhoneBuild84ReadProjectiveBlockValue' 'Build84 qualified reader'
  echo
  reader_section 'watermark/experimental_v4_phone_build86_dct_hoist.go' 'experimentalV4PhoneBuild86ReadProjectiveBlockValue' 'Build86 DCT table-hoist candidate'
} > "$stage_dir/build86-bce-reader.txt"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build86 DCT-hoist benchmark written to:"
echo "  $OUTPUT_DIR/build86-benchmark.txt"
echo "  $OUTPUT_DIR/build86-bce-reader.txt"
echo "  $OUTPUT_DIR/build86-bce-inlining.txt"
echo "Build86 benchmark gate: PASS (correctness only; review whole-reader speed before physical testing)"
