#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD101_PROFILE_DIR:-v4-phone private/build101-runtime-effects}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
export GOMAXPROCS="${GOMAXPROCS:-8}"
[[ "$GOMAXPROCS" =~ ^[1-9][0-9]*$ ]] || { echo 'error: invalid GOMAXPROCS' >&2; exit 1; }
[[ ! -e "$OUTPUT_DIR" && ! -L "$OUTPUT_DIR" ]] || { echo 'error: existing evidence is preserved; select another output directory' >&2; exit 1; }
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"
for image in phone-b-mild.jpg phone-b-angle.jpg; do [[ -f "$ACQUISITION_DIR/$image" ]] || { echo "error: missing retained acquisition $image" >&2; exit 1; }; done
python3 scripts/check-v4-build101-observed-copies.py
parent="$(dirname "$OUTPUT_DIR")"; base="$(basename "$OUTPUT_DIR")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build101.XXXXXX")"
stage="$(cd "$stage" && pwd -P)"
published=false
cleanup(){
  local rc=$?
  trap - EXIT
  if [[ "$published" != true && -d "$stage" ]]; then
    local failed="${OUTPUT_DIR}-failed-$(basename "$stage")"
    mv -T "$stage" "$failed"
    echo "Build101 failed evidence preserved: $failed" >&2
  fi
  exit "$rc"
}
trap cleanup EXIT
{
 echo 'PixSeal Build101 runtime-effects observability; no runtime candidate or promotion'
 echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
 echo "toolchain=$TOOLCHAIN"
 echo "gomaxprocs=$GOMAXPROCS"
 echo "gogc=${GOGC:-default-100}"
 echo "gomemlimit=${GOMEMLIMIT:-off}"
 echo 'qualified_baseline=Build98; previous=Build93; comparator=unpromoted Build100'
 echo 'timing=unchanged pinned APIs; full AB/BA, both retained cases, 8 runs'
 echo 'attribution=test-only structural copies; stage labels plus ReadMemStats/runtime metrics; 4 separate runs'
 echo 'attribution_order=mild 98-100; angle 100-98'
 echo 'interpretation=attribution timing cannot qualify speed; estimated runtime CPU is not OS CPU; scheduler buckets are lower-bound thresholds'
 echo 'closure=all 12 semantic runs exact, source/corpus hashes stable, complete timing/stage/profile artifacts'
 echo 'decision=report signs and worker saturation; no causal claim without matching CPU/GC/scheduler evidence; no promotion or relaxed Build100 guard'
 env GOTOOLCHAIN="$TOOLCHAIN" go version
 uname -sm
 awk -F ': ' '/model name/{print "cpu=" $2;exit}' /proc/cpuinfo
} > "$stage/build101-metadata.txt"
{
  rg --files -g '*.go' watermark internal cmd
  if [[ -f go.sum ]]; then printf '%s\n' go.sum; fi
  printf '%s\n' go.mod VERSION scripts/profile-v4-build101-runtime-effects.sh scripts/check-v4-build101-observed-copies.py
} | LC_ALL=C sort | while IFS= read -r path; do sha256sum "$path"; done > "$stage/build101-source.sha256"
sha256sum "$ACQUISITION_DIR/phone-b-mild.jpg" "$ACQUISITION_DIR/phone-b-angle.jpg" > "$stage/build101-corpus.sha256"
env GOTOOLCHAIN="$TOOLCHAIN" go test -c ./watermark -o "$stage/build101-watermark.test"
(cd "$stage" && sha256sum build101-watermark.test) > "$stage/build101-binary.sha256"
timeout --foreground "${V4_PHONE_BUILD101_TIMEOUT:-7200}s" env \
 PIXSEAL_BUILD101_ACQUISITION_DIR="$ACQUISITION_DIR" \
 PIXSEAL_BUILD101_OUTPUT_DIR="$stage" \
 PIXSEAL_BUILD101_KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}" \
 PIXSEAL_BUILD101_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}" \
 PIXSEAL_BUILD101_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}" \
 PIXSEAL_BUILD101_MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}" \
 PIXSEAL_BUILD101_COOLDOWN_MS="${V4_PHONE_BUILD101_COOLDOWN_MS:-5000}" \
 "$stage/build101-watermark.test" -test.run '^TestExperimentalV4Build101RetainedStudy$' -test.count=1 -test.v -test.timeout=0 2>&1 | tee "$stage/build101-raw.txt"
awk -F '\t' '$1=="build101-header" || $1=="build101-row" {sub(/^[^\t]*\t/,"");print}' "$stage/build101-raw.txt" > "$stage/build101-runs.tsv"
awk -F '\t' '$1=="build101-stage-header" || $1=="build101-stage" {sub(/^[^\t]*\t/,"");print}' "$stage/build101-raw.txt" > "$stage/build101-stages.tsv"
awk -F '\t' 'NR==1{w=NF;next}NF!=w{bad=1}END{if(NR!=13||bad)exit 1}' "$stage/build101-runs.tsv"
awk -F '\t' 'NR==1{w=NF;next}NF!=w{bad=1}END{if(NR!=27||bad)exit 1}' "$stage/build101-stages.tsv"
for profile in "$stage"/*.cpu.prof; do
 env GOTOOLCHAIN="$TOOLCHAIN" go tool pprof -top -cum -nodecount=25 -tagfocus='stage=gen4' "$stage/build101-watermark.test" "$profile" > "${profile%.cpu.prof}-gen4-cpu.txt"
 env GOTOOLCHAIN="$TOOLCHAIN" go tool pprof -tags "$stage/build101-watermark.test" "$profile" > "${profile%.cpu.prof}-tags.txt"
done
sha256sum -c "$stage/build101-source.sha256" > "$stage/build101-source-verified.txt"
sha256sum -c "$stage/build101-corpus.sha256" > "$stage/build101-corpus-verified.txt"
mv -T "$stage" "$OUTPUT_DIR"
published=true
echo "Build101 observability semantic gate: PASS; evidence: $OUTPUT_DIR"
