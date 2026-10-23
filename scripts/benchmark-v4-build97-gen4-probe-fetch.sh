#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"; cd "$ROOT"
OUTPUT_DIR="${V4_BUILD97_BENCH_DIR:-v4-phone private/build97-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
BENCHTIME="${V4_BUILD97_BENCHTIME:-2s}"
COUNT="${V4_BUILD97_COUNT:-5}"
command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ "$COUNT" =~ ^[1-9][0-9]*$ ]] || { echo "error: V4_BUILD97_COUNT must be positive" >&2; exit 1; }
parent="$(dirname "$OUTPUT_DIR")"; base="$(basename "$OUTPUT_DIR")"; mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build97.XXXXXX")"; published=false
trap 'rc=$?; set +e; if [[ "$published" != true && -d "$stage" ]]; then if ((rc!=0)); then rm -rf "${OUTPUT_DIR}-failed"; mv "$stage" "${OUTPUT_DIR}-failed"; else rm -rf "$stage"; fi; fi; exit $rc' EXIT
GO=(env "GOTOOLCHAIN=$TOOLCHAIN" go)
{
 echo "PixSeal Build97 exact Gen4 probe-fetch extension public benchmark"
 echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"; echo "toolchain=$TOOLCHAIN"; "${GO[@]}" version
 echo "host=$(uname -a)"; echo "bench_time=$BENCHTIME"; echo "count=$COUNT"
 echo "current_qualified_baseline=Build93"; echo "candidate=Build97 Exact Gen4 Probe Fetch Extension"
 echo "selection_evidence=Build96 retained pprof"
 echo "changed_scope=Gen4 single-probe and pair-stencil FoldScore calls reuse qualified Build84 exact RGB reader; continuation remains Build84"
 echo "extra_image_memory=0"
 echo "build97_source_sha256=$(sha256sum watermark/experimental_v4_phone_build97_gen4_probe_fetch.go | awk '{print $1}')"
} > "$stage/build97-metadata.txt"
"${GO[@]}" test ./watermark -run '^TestExperimentalV4Build97(SingleProbeMatchesBuild53Exactly|PairStencilMatchesBuild53Exactly|Generation4MatchesBuild84Exactly)$' -count=1 | tee "$stage/build97-fixture-test.txt"
"${GO[@]}" test -race ./watermark -run '^TestExperimentalV4Build97ParallelGeneration4RaceFixture$' -count=1 | tee "$stage/build97-race-test.txt"
"${GO[@]}" test ./watermark -run '^$' -bench '^BenchmarkExperimentalV4Build97' -benchmem -benchtime="$BENCHTIME" -count="$COUNT" | tee "$stage/build97-benchmark.txt"
{
 printf 'benchmark\tns_per_op\tbytes_per_op\tallocs_per_op\traw\n'
 awk '/^BenchmarkExperimentalV4Build97/ {name=$1;ns="";bytes="";allocs="";for(i=2;i<=NF;i++){if($i=="ns/op")ns=$(i-1);if($i=="B/op")bytes=$(i-1);if($i=="allocs/op")allocs=$(i-1)}raw=$0;gsub(/\t/," ",raw);printf "%s\t%s\t%s\t%s\t%s\n",name,ns,bytes,allocs,raw}' "$stage/build97-benchmark.txt"
} > "$stage/build97-benchmark.tsv"
cat > "$stage/build97-benchmark.md" <<'MD'
# PixSeal Build97 public benchmark

Build97 is benchmark-first. It extends the already-qualified exact Build84 RGB projective reader to generation-four single-probe and pair-stencil FoldScore calls. Continuation4 remains Build84. There is no luminance plane and no image-sized memory allocation.

Promotion is not implied by this benchmark. Continue to retained paired profiling only if the full Generation4 benchmark is materially and directionally faster, not merely an isolated single/pair microbenchmark.
MD
rm -rf "$OUTPUT_DIR"; mv "$stage" "$OUTPUT_DIR"; published=true
echo "Build97 benchmark written to: $OUTPUT_DIR"; echo 'Build97 benchmark gate: PASS (semantic/race only; review performance before retained paired profiling)'
