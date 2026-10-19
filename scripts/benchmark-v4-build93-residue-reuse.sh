#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD93_BENCH_DIR:-v4-phone private/build93-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD93_BENCHTIME:-2s}"
COUNT="${V4_BUILD93_COUNT:-5}"

[[ "$COUNT" =~ ^[1-9][0-9]*$ ]] || { echo "error: V4_BUILD93_COUNT must be a positive integer" >&2; exit 1; }

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build93.XXXXXX")"
cleanup() { rm -rf "$stage_dir"; }
trap cleanup EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build93 exact protected-residue reuse benchmark"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "qualified_baseline=Build93"
  echo "previous_qualified_baseline=Build90"
  echo "selection_evidence=Build92-qualified-Build90-post-promotion-profile"
  echo "subject=Build93-qualified-exact-protected-residue-reuse"
  echo "workload=public-deterministic-capacity-profile-592x512-strength48"
  echo "single_candidate=historical-Build67-vs-Build93"
  echo "batch_candidate_count=8"
  echo "profile_order=robust-balanced-capacity"
  echo "physical_qualification=PASS-9of9-x2"
  echo "production_entrypoint_integration=separate-gate"
} > "$stage_dir/build93-metadata.txt"

"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build93PublicBenchmarkFixture$' -count=1 \
  | tee "$stage_dir/build93-fixture-test.txt"

"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build93' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build93-benchmark.txt"

"${GO[@]}" test -race ./watermark -run '^TestExperimentalV4Build93ParallelDecodeRaceFixture$' -count=1 \
  | tee "$stage_dir/build93-race-test.txt"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
trap - EXIT

echo "Build93 exact protected-residue reuse benchmark written to:"
echo "  $OUTPUT_DIR/build93-benchmark.txt"
echo "  $OUTPUT_DIR/build93-fixture-test.txt"
echo "  $OUTPUT_DIR/build93-race-test.txt"
echo "  $OUTPUT_DIR/build93-metadata.txt"
echo "Build93 benchmark gate: PASS (correctness/race recheck; qualified performance and physical evidence are frozen under docs/qualified-baselines/)"
