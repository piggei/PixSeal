#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD90_BENCH_DIR:-v4-phone private/build90-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD90_BENCHTIME:-2s}"
COUNT="${V4_BUILD90_COUNT:-5}"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build90.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build90 ordered-parallel qualification benchmark"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "qualified_baseline=Build84"
  echo "selection_evidence=Build89-full-pipeline-profile"
  echo "candidate=Build90-ordered-parallel-qualification"
  echo "workload=public-deterministic-pilot-only-no-key-no-payload"
  echo "high_pass_tasks=24"
  echo "early_reject_tasks=512"
  echo "physical_gate=deferred-until-benchmark-review"
} > "$stage_dir/build90-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build90PublicBenchmarkFixture$' -count=1 \
  | tee "$stage_dir/build90-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build90' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build90-benchmark.txt"

"${GO[@]}" test -race ./watermark -run '^TestExperimentalV4Build90ParallelQualificationMatchesSerial$' -count=1 \
  | tee "$stage_dir/build90-race-test.txt"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build90 ordered-parallel qualification benchmark written to:"
echo "  $OUTPUT_DIR/build90-benchmark.txt"
echo "  $OUTPUT_DIR/build90-fixture-test.txt"
echo "  $OUTPUT_DIR/build90-race-test.txt"
echo "  $OUTPUT_DIR/build90-metadata.txt"
echo "Build90 benchmark gate: PASS (correctness/race only; review high-pass speedup and early-reject overhead before physical testing)"
