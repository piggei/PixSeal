#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
OUTPUT_DIR="${V4_BUILD104_BENCH_DIR:-v4-phone private/build104-fixed-block-benchmark}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
[[ ! -e "$OUTPUT_DIR" && ! -L "$OUTPUT_DIR" ]] || { echo 'error: existing Build104 evidence preserved; select a new output directory' >&2; exit 1; }
parent="$(dirname "$OUTPUT_DIR")";base="$(basename "$OUTPUT_DIR")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build104.XXXXXX")";stage="$(cd "$stage" && pwd -P)"
published=false
cleanup(){ local rc=$?;trap - EXIT;if [[ "$published" != true && -d "$stage" ]];then mv -T "$stage" "${OUTPUT_DIR}-failed-$(basename "$stage")";fi;exit "$rc"; }
trap cleanup EXIT
python3 scripts/check-v4-build104-fixed-block-copy.py
cp docs/V4_BUILD104_FIXED_BLOCK_READER.md "$stage/build104-predeclared-protocol.md"
[[ ! -f 'v4-phone private/build104-source-gate.txt' ]] || cp 'v4-phone private/build104-source-gate.txt' "$stage/build104-source-gate.txt"
{
 echo 'Build104 test-only fixed8 continuation reader; public deterministic data'
 date -u +%Y-%m-%dT%H:%M:%SZ
 env GOTOOLCHAIN="$TOOLCHAIN" go version
 uname -sm
 awk -F ': ' '/model name/{print "cpu=" $2;exit}' /proc/cpuinfo
 echo 'runtime=GOMAXPROCS8;GODEBUGasyncpreemptoff0;GOGC100;GOMEMLIMIToff'
 echo 'protocol=two sessions; five metrics; Q C C Q | Q Q | C Q Q C per metric; 100 fresh children'
 echo 'timing=benchtime1s minimum, cooldown1s, external child180s, overall1800s'
 echo 'decision=each session Gen4 and Pool8 mean>=3% faster and every stage pair favorable; reader means non-regressive; reader/Gen4 mean B/op+allocs exact; Pool8 count exact and bytes within1%'
 echo 'control=Q/Q descriptive; never relax speed threshold'
 echo 'fixture=public-seed107; actual proposal; 256evals/6bank/128continuationFoldScores/16384blocks per Gen4'
 echo 'limits=one host; synthetic workloads; no private/physical/HMAC qualification or production promotion'
} > "$stage/build104-metadata.txt"
{
 rg --files -g '*.go' watermark internal cmd
 printf '%s\n' go.mod VERSION Makefile scripts/test-list.sh scripts/check-v4-build104-fixed-block-copy.py scripts/run-v4-build104-fixed-block.py scripts/summarize-v4-build104-fixed-block.py scripts/benchmark-v4-build104-fixed-block.sh
} | LC_ALL=C sort | while IFS= read -r path;do sha256sum "$path";done > "$stage/build104-source.sha256"
env GOTOOLCHAIN="$TOOLCHAIN" go test -c ./watermark -o "$stage/build104-watermark.test"
(cd "$stage" && sha256sum build104-watermark.test build104-predeclared-protocol.md) > "$stage/build104-binary-protocol.sha256"
timeout --foreground --kill-after=15s 1800s env PIXSEAL_BUILD104_OUTPUT_DIR="$stage" python3 scripts/run-v4-build104-fixed-block.py 2>&1 | tee "$stage/build104-run-console.txt"
python3 scripts/summarize-v4-build104-fixed-block.py "$stage" "$stage/reviewed-public" | tee "$stage/build104-decision.txt"
sha256sum -c "$stage/build104-source.sha256" > "$stage/build104-source-verified.txt"
(cd "$stage" && sha256sum -c build104-binary-protocol.sha256) > "$stage/build104-binary-protocol-verified.txt"
mv -T "$stage" "$OUTPUT_DIR";published=true
echo "Build104 public benchmark observability gate: PASS; outcome recorded separately; evidence: $OUTPUT_DIR"
