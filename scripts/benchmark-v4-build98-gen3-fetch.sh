#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
OUTPUT_DIR="${V4_BUILD98_BENCH_DIR:-v4-phone private/build98-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD98_BENCHTIME:-2s}"
COUNT="${V4_BUILD98_COUNT:-5}"
command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ "$COUNT" =~ ^[1-9][0-9]*$ ]] || { echo "error: V4_BUILD98_COUNT must be positive" >&2; exit 1; }
parent="$(dirname "$OUTPUT_DIR")"
base="$(basename "$OUTPUT_DIR")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build98.XXXXXX")"
published=false
trap 'rc=$?; set +e; if [[ "$published" != true && -d "$stage" ]]; then if ((rc!=0)); then rm -rf "${OUTPUT_DIR}-failed"; mv "$stage" "${OUTPUT_DIR}-failed"; else rm -rf "$stage"; fi; fi; exit $rc' EXIT
GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)
{
  echo "PixSeal Build98 exact Generation3 fetch extension public benchmark"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "toolchain=$TOOLCHAIN"
  "${GO[@]}" version
  echo "host=$(uname -a)"
  echo "bench_time=$BENCHTIME"
  echo "count=$COUNT"
  echo "current_qualified_baseline=Build93"
  echo "candidate=Build98 Exact Generation3 Fetch Extension"
  echo "selection_evidence=Build94 retained CPU profile plus Build97 retained rejection"
  echo "changed_scope=Generation3 single/pair/continuation/sibling FoldScore calls reuse qualified Build84 exact RGB reader"
  echo "unchanged_scope=prefix1 generation2 generation4 freeze qualification protected decode"
  echo "extra_image_memory=0"
  echo "build98_source_sha256=$(sha256sum watermark/experimental_v4_phone_build98_gen3_fetch.go | awk '{print $1}')"
} > "$stage/build98-metadata.txt"
"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build98(SiblingStencilMatchesBuild55Exactly|Generation3MatchesBuild73Exactly)$' -count=1 | tee "$stage/build98-fixture-test.txt"
"${GO[@]}" test -race ./watermark -run '^TestExperimentalV4Build98ParallelGeneration3RaceFixture$' -count=1 | tee "$stage/build98-race-test.txt"
"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build98' -benchmem -benchtime="$BENCHTIME" -count="$COUNT" | tee "$stage/build98-benchmark.txt"
{
  printf 'benchmark\tns_per_op\tbytes_per_op\tallocs_per_op\traw\n'
  awk '/^BenchmarkExperimentalV4Build98/ {name=$1;ns="";bytes="";allocs="";for(i=2;i<=NF;i++){if($i=="ns/op")ns=$(i-1);if($i=="B/op")bytes=$(i-1);if($i=="allocs/op")allocs=$(i-1)}raw=$0;gsub(/\t/," ",raw);printf "%s\t%s\t%s\t%s\t%s\n",name,ns,bytes,allocs,raw}' "$stage/build98-benchmark.txt"
} > "$stage/build98-benchmark.tsv"
cat > "$stage/build98-benchmark.md" <<'MD'
# PixSeal Build98 public benchmark

Build98 is benchmark-first. It reuses the already-qualified exact Build84 RGB projective reader for every FoldScore inside generation three: single probe, pair stencil, continuation and sibling stencil. Prefix1, generation two, generation four, freeze, qualification and protected decode remain the qualified Build93 path. Build98 adds no image-sized memory.

The decision metric is the complete Generation3 benchmark. Component results are diagnostic only. No retained or physical gate is authorized unless the full Generation3 result is materially and directionally faster across the repeated qualified-host samples.
MD
rm -rf "$OUTPUT_DIR"
mv "$stage" "$OUTPUT_DIR"
published=true
echo "Build98 benchmark written to: $OUTPUT_DIR"
echo 'Build98 benchmark gate: PASS (semantic/race only; review full Generation3 performance before retained profiling)'
