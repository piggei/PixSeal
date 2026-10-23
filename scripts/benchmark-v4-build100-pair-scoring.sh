#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
OUTPUT_DIR="${V4_BUILD100_BENCH_DIR:-v4-phone private/build100-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD100_BENCHTIME:-2s}"
COUNT="${V4_BUILD100_COUNT:-5}"
export GOMAXPROCS="${GOMAXPROCS:-8}"
[[ "$COUNT" =~ ^[1-9][0-9]*$ ]] || { echo 'error: benchmark count must be positive' >&2; exit 1; }
[[ ! -e "$OUTPUT_DIR" && ! -L "$OUTPUT_DIR" ]] || { echo 'error: existing benchmark output is preserved; choose another directory' >&2; exit 1; }
parent="$(dirname "$OUTPUT_DIR")"
base="$(basename "$OUTPUT_DIR")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build100.XXXXXX")"
published=false
cleanup(){
  local rc=$?
  trap - EXIT
  if [[ "$published" != true && -d "$stage" ]]; then
    mv "$stage" "${OUTPUT_DIR}-failed-$(basename "$stage")"
  fi
  exit "$rc"
}
trap cleanup EXIT
{
  echo 'PixSeal Build100 public complete-freeze benchmark'
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "toolchain=$TOOLCHAIN"
  echo "gomaxprocs=$GOMAXPROCS"
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo 'qualified_baseline=Build98'
  echo 'candidate=Build100 ordered-parallel six-pair public proposal scores'
  echo 'additional_image_sized_memory_bytes=0'
  echo 'decision_rule=complete-freeze mean at least 5% faster, all corresponding repeated samples directionally faster; allocations bounded to small task/worker structures'
  env GOTOOLCHAIN="$TOOLCHAIN" go version
} > "$stage/build100-metadata.txt"
sha256sum watermark/experimental_v4_phone_build100_*.go watermark/experimental_v4_build100_pair_scoring*go > "$stage/build100-source.sha256"
env GOTOOLCHAIN="$TOOLCHAIN" go test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build100Freeze$' -benchmem -benchtime="$BENCHTIME" -count="$COUNT" | tee "$stage/build100-benchmark.txt"
awk 'BEGIN{print "benchmark\tns_per_op\tbytes_per_op\tallocs_per_op"}/^BenchmarkExperimentalV4Build100Freeze/{name=$1;for(i=2;i<=NF;i++){if($i=="ns/op")ns=$(i-1);if($i=="B/op")bytes=$(i-1);if($i=="allocs/op")allocs=$(i-1)}print name"\t"ns"\t"bytes"\t"allocs}' "$stage/build100-benchmark.txt" > "$stage/build100-benchmark.tsv"
mv -T "$stage" "$OUTPUT_DIR"
published=true
echo "Build100 benchmark complete: $OUTPUT_DIR; review the predeclared decision rule before retained comparison."
