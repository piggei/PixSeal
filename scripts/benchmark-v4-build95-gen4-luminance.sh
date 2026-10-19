#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OUTPUT_DIR="${V4_BUILD95_BENCH_DIR:-v4-phone private/build95-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD95_BENCHTIME:-2s}"
COUNT="${V4_BUILD95_COUNT:-5}"

command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ "$COUNT" =~ ^[1-9][0-9]*$ ]] || { echo "error: V4_BUILD95_COUNT must be a positive integer" >&2; exit 1; }

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build95.XXXXXX")"
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
      echo "Build95 FAILED benchmark evidence preserved in: $failed_output_dir" >&2
    else
      rm -rf "$stage_dir"
    fi
  fi
  exit "$rc"
}
trap cleanup_stage EXIT

GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)

{
  echo "PixSeal Build95 exact Gen4 luminance-plane public benchmark"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "host=$(uname -a)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "candidate=Build95 Exact Gen4 Luminance Plane"
  echo "selection_evidence=Build94 qualified Build93 post-promotion profile"
  echo "current_qualified_baseline=Build93"
  echo "previous_qualified_baseline=Build90"
  echo "changed_scope=continuation4 projective RGB-to-luminance reads only"
  echo "fallback=unchanged Build84 reader when the public-dimension 128 MiB luminance-plane cap is exceeded"
  echo "benchmark_workload=public deterministic synthetic fixtures, including 3464x4600 retained-size plane preparation; no key, payload, HMAC, ECC or private acquisition"
  echo "promotion_status=candidate only; review same-session benchmark before retained paired profile"
  echo "build95_source_sha256=$(sha256sum watermark/experimental_v4_phone_build95_gen4_luminance.go | awk '{print $1}')"
  echo "build84_reader_sha256=$(sha256sum watermark/experimental_v4_phone_build84_rgb_fetch.go | awk '{print $1}')"
} > "$stage_dir/build95-metadata.txt"

"${GO[@]}" test ./watermark \
  -run '^TestExperimentalV4Build95(PublicBenchmarkFixture|LuminancePlaneBitExact|LuminancePlaneBoundedFallback|ProjectiveBlockMatchesBuild84BitExactly|FoldScoreMatchesBuild84BitExactly|ContinueMatchesBuild84Exactly|Generation4MatchesBuild84Exactly)$' \
  -count=1 | tee "$stage_dir/build95-fixture-test.txt"

"${GO[@]}" test -race ./watermark \
  -run '^TestExperimentalV4Build95ParallelContinuationRaceFixture$' \
  -count=1 | tee "$stage_dir/build95-race-test.txt"

"${GO[@]}" test ./watermark -run '^$' \
  -bench '^BenchmarkExperimentalV4Build95' \
  -benchmem -benchtime="$BENCHTIME" -count="$COUNT" \
  | tee "$stage_dir/build95-benchmark.txt"

# Keep a compact, machine-readable list of the benchmark rows. Raw Go output is
# authoritative; this TSV is only a convenient index for later review.
{
  printf 'benchmark\tns_per_op\tbytes_per_op\tallocs_per_op\traw\n'
  awk '
    /^BenchmarkExperimentalV4Build95/ {
      name=$1; ns=""; bytes=""; allocs="";
      for (i=2; i<=NF; i++) {
        if ($i == "ns/op" && i>1) ns=$(i-1);
        if ($i == "B/op" && i>1) bytes=$(i-1);
        if ($i == "allocs/op" && i>1) allocs=$(i-1);
      }
      raw=$0; gsub(/\t/, " ", raw);
      printf "%s\t%s\t%s\t%s\t%s\n", name, ns, bytes, allocs, raw;
    }
  ' "$stage_dir/build95-benchmark.txt"
} > "$stage_dir/build95-benchmark.tsv"

{
  echo '# PixSeal Build95 public benchmark'; echo
  echo 'Build95 is a benchmark-first candidate selected from the qualified Build94 profile. The semantic and race gates are authoritative. Benchmark timings are same-session decision evidence only; they do not promote Build95 and do not replace retained-corpus comparison.'; echo
  echo '- Baseline reader: Build84 exact continuation4 RGB fetch, inherited by qualified Build93.'
  echo '- Candidate reader: one immutable exact float64 luminance plane, prepared immediately before generation 4.'
  echo '- Preparation rows include both the compact public fixture and the public 3464x4600 retained-size dimensions (127,475,200 bytes).'
  echo '- Changed scope: continuation4 FoldScore reads only. Prefix1, generations 2/3, generation4 single/pair probes, bank order, qualification and protected decode are unchanged.'
  echo '- Memory guard: deterministic 128 MiB additional-plane cap based only on public image dimensions; over-cap inputs use the unchanged Build84 reader.'
  echo
  echo 'Review `build95-benchmark.txt` or the convenience `build95-benchmark.tsv` for all repeated samples. A retained B/mild/B/angle paired AB/BA run is permitted only after the candidate is both exact and materially faster in this public benchmark.'
} > "$stage_dir/build95-benchmark.md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build95 exact Gen4 luminance-plane benchmark written to:'
echo "  $OUTPUT_DIR/build95-benchmark.txt"
echo "  $OUTPUT_DIR/build95-benchmark.tsv"
echo "  $OUTPUT_DIR/build95-benchmark.md"
echo "  $OUTPUT_DIR/build95-fixture-test.txt"
echo "  $OUTPUT_DIR/build95-race-test.txt"
echo "  $OUTPUT_DIR/build95-metadata.txt"
echo 'Build95 benchmark gate: PASS (semantic/race correctness only; review performance before paired physical profiling)'
