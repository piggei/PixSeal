#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD103_PROFILE_DIR:-v4-phone private/build103-preemption-control}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
[[ ! -e "$OUTPUT_DIR" && ! -L "$OUTPUT_DIR" ]] || { echo 'error: existing preemption evidence is preserved; select another output directory' >&2; exit 1; }
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"
for image in phone-b-mild.jpg phone-b-angle.jpg; do [[ -f "$ACQUISITION_DIR/$image" ]] || { echo "error: missing retained acquisition $image" >&2; exit 1; }; done
python3 scripts/check-v4-build101-observed-copies.py
python3 scripts/check-v4-build102-replay-copy.py
parent="$(dirname "$OUTPUT_DIR")"; base="$(basename "$OUTPUT_DIR")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build103.XXXXXX")"
stage="$(cd "$stage" && pwd -P)"
published=false
cleanup(){
 local rc=$?
 trap - EXIT
 if [[ "$published" != true && -d "$stage" ]]; then
  local failed="${OUTPUT_DIR}-failed-$(basename "$stage")"
  mv -T "$stage" "$failed"
  echo "Build103 failed evidence preserved: $failed" >&2
 fi
 exit "$rc"
}
trap cleanup EXIT
cp docs/V4_BUILD103_PREEMPTION_ATTRIBUTION_CONTROL.md "$stage/build103-predeclared-protocol.md"
(cd "$stage" && sha256sum build103-predeclared-protocol.md) > "$stage/build103-protocol.sha256"
export GOMAXPROCS=8 GOGC=100 GOMEMLIMIT=off GOTRACEBACK=single
{
 echo 'PixSeal Build103 fresh-process preemption-attribution control; no runtime promotion'
 echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
 echo "toolchain=$TOOLCHAIN"
 echo 'runtime=GOMAXPROCS=8;GOGC=100;GOMEMLIMIT=off;GOTRACEBACK=single'
 echo 'GODEBUG=asyncpreemptoff=0 or asyncpreemptoff=1; inherited GODEBUG replaced per child'
 echo "clock_ticks_per_s=$(getconf CLK_TCK)"
 echo 'timing_order=D O O D | O D D O per case; 16 fresh-process timing replays'
 echo 'profile_order=mild D-O; angle O-D; 4 additional fresh processes excluded from timing means'
 echo 'data=one bit-preserving private snapshot per case; qualified full-bank reference and frozen Build102 hashes; same binary in all conditions'
 echo 'closure=all20 banks and non-timing telemetry exact; snapshot/raster/task/source/corpus stability; both conditions and four profiles complete'
 echo 'attribution_rule=support a material diagnostic effect only if off improves mean pool wall and live total CPU by >=5% on both cases and every paired wall/CPU comparison is favorable'
 echo 'interpretation=otherwise mixed/unsupported; disappearance of asyncPreempt samples alone is not a speedup or proof of pure sampling artifact'
 echo 'limits=fresh-process address/layout variation; one host/session; debug flag also affects GC stack handling; hard external child timeout600s; no physical/production qualification'
 env GOTOOLCHAIN="$TOOLCHAIN" go version
 uname -sm
 awk -F ': ' '/model name/{print "cpu=" $2;exit}' /proc/cpuinfo
} > "$stage/build103-metadata.txt"
{
 rg --files -g '*.go' watermark internal cmd
 if [[ -f go.sum ]]; then printf '%s\n' go.sum; fi
 printf '%s\n' go.mod VERSION Makefile scripts/test-list.sh scripts/profile-v4-build103-preemption-control.sh scripts/run-v4-build103-preemption-control.py scripts/summarize-v4-build103-preemption-control.py scripts/check-v4-build101-observed-copies.py scripts/check-v4-build102-replay-copy.py docs/qualified-baselines/build101-runs.tsv docs/qualified-baselines/build102-captures.tsv
} | LC_ALL=C sort | while IFS= read -r path; do sha256sum "$path"; done > "$stage/build103-source.sha256"
sha256sum "$ACQUISITION_DIR/phone-b-mild.jpg" "$ACQUISITION_DIR/phone-b-angle.jpg" > "$stage/build103-corpus.sha256"
env GOTOOLCHAIN="$TOOLCHAIN" go test -c ./watermark -o "$stage/build103-watermark.test"
(cd "$stage" && sha256sum build103-watermark.test) > "$stage/build103-binary.sha256"
timeout --foreground --kill-after=15s "${V4_PHONE_BUILD103_TIMEOUT:-7200}s" env \
 PIXSEAL_BUILD103_ACQUISITION_DIR="$ACQUISITION_DIR" \
 PIXSEAL_BUILD103_OUTPUT_DIR="$stage" \
 PIXSEAL_BUILD103_BANK_REFERENCE="$ROOT/docs/qualified-baselines/build101-runs.tsv" \
 PIXSEAL_BUILD103_CAPTURE_REFERENCE="$ROOT/docs/qualified-baselines/build102-captures.tsv" \
 PIXSEAL_BUILD103_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}" \
 PIXSEAL_BUILD103_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}" \
 PIXSEAL_BUILD103_CLOCK_TICKS="$(getconf CLK_TCK)" \
 PIXSEAL_BUILD103_COOLDOWN_MS="${V4_PHONE_BUILD103_COOLDOWN_MS:-5000}" \
 python3 scripts/run-v4-build103-preemption-control.py 2>&1 | tee "$stage/build103-raw.txt"
awk -F '\t' 'NR==1{w=NF;next}NF!=w{bad=1}END{if(NR!=21||bad)exit 1}' "$stage/build103-runs.tsv"
awk -F '\t' 'NR==1{w=NF;next}NF!=w{bad=1}END{if(NR!=3||bad)exit 1}' "$stage/build103-captures.tsv"
for profile in "$stage"/*.cpu.prof; do
 env GOTOOLCHAIN="$TOOLCHAIN" go tool pprof -top -cum -nodecount=35 "$stage/build103-watermark.test" "$profile" > "${profile%.cpu.prof}-cpu.txt"
done
sha256sum -c "$stage/build103-source.sha256" > "$stage/build103-source-verified.txt"
sha256sum -c "$stage/build103-corpus.sha256" > "$stage/build103-corpus-verified.txt"
(cd "$stage" && sha256sum phone-b-mild.jpg.gob phone-b-angle.jpg.gob) > "$stage/build103-snapshots.sha256"
mv -T "$stage" "$OUTPUT_DIR"
published=true
echo "Build103 preemption observability gate: PASS; evidence: $OUTPUT_DIR"
