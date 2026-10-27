#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD102_PROFILE_DIR:-v4-phone private/build102-frozen-gen4-replay}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
export GOMAXPROCS="${GOMAXPROCS:-8}"
[[ "$GOMAXPROCS" =~ ^[1-9][0-9]*$ ]] || { echo 'error: invalid GOMAXPROCS' >&2; exit 1; }
[[ ! -e "$OUTPUT_DIR" && ! -L "$OUTPUT_DIR" ]] || { echo 'error: existing replay evidence is preserved; select another output directory' >&2; exit 1; }
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"
for image in phone-b-mild.jpg phone-b-angle.jpg; do [[ -f "$ACQUISITION_DIR/$image" ]] || { echo "error: missing retained acquisition $image" >&2; exit 1; }; done
python3 scripts/check-v4-build101-observed-copies.py
python3 scripts/check-v4-build102-replay-copy.py
parent="$(dirname "$OUTPUT_DIR")"; base="$(basename "$OUTPUT_DIR")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build102.XXXXXX")"
stage="$(cd "$stage" && pwd -P)"
published=false
cleanup(){
  local rc=$?
  trap - EXIT
  if [[ "$published" != true && -d "$stage" ]]; then
    local failed="${OUTPUT_DIR}-failed-$(basename "$stage")"
    mv -T "$stage" "$failed"
    echo "Build102 failed evidence preserved: $failed" >&2
  fi
  exit "$rc"
}
trap cleanup EXIT
ticks="$(getconf CLK_TCK)"
{
 echo 'PixSeal Build102 frozen Gen4 A/A replay observability; no candidate or promotion'
 echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
 echo "toolchain=$TOOLCHAIN"
 echo "gomaxprocs=$GOMAXPROCS"
 echo "clock_ticks_per_s=$ticks"
 echo "gogc=${GOGC:-default-100}"
 echo "gomemlimit=${GOMEMLIMIT:-off}"
 echo 'qualified_baseline=Build98; previous=Build93; Build100=closed/not-promoted'
 echo 'protocol=each case: qualified full-bank reference, captures from 98/100 prefixes, shared raster/vector, 8 A/A timing slots then 2 separate CPU-profile slots'
 echo 'timing_order=A B B A | B A A B; labels are measurements of identical code/data, not runtime variants'
 echo 'cpu_order=mild A-B; angle B-A; excluded from timing means'
 echo 'closure=20 complete-bank exact replays; identical task/order/score hashes and raster; frozen reference counters; source/corpus stability; complete rows and 4 profiles'
 echo 'decision=report each placebo pair, range and CPU/worker/GC data; no significance or causal claim from overlap alone; no guard change or physical qualification'
 env GOTOOLCHAIN="$TOOLCHAIN" go version
 uname -sm
 awk -F ': ' '/model name/{print "cpu=" $2;exit}' /proc/cpuinfo
} > "$stage/build102-metadata.txt"
{
 rg --files -g '*.go' watermark internal cmd
 if [[ -f go.sum ]]; then printf '%s\n' go.sum; fi
 printf '%s\n' go.mod VERSION scripts/profile-v4-build102-frozen-gen4-replay.sh scripts/check-v4-build101-observed-copies.py scripts/check-v4-build102-replay-copy.py docs/qualified-baselines/build101-runs.tsv
} | LC_ALL=C sort | while IFS= read -r path; do sha256sum "$path"; done > "$stage/build102-source.sha256"
sha256sum "$ACQUISITION_DIR/phone-b-mild.jpg" "$ACQUISITION_DIR/phone-b-angle.jpg" > "$stage/build102-corpus.sha256"
env GOTOOLCHAIN="$TOOLCHAIN" go test -c ./watermark -o "$stage/build102-watermark.test"
(cd "$stage" && sha256sum build102-watermark.test) > "$stage/build102-binary.sha256"
timeout --foreground "${V4_PHONE_BUILD102_TIMEOUT:-7200}s" env \
 PIXSEAL_BUILD102_ACQUISITION_DIR="$ACQUISITION_DIR" \
 PIXSEAL_BUILD102_OUTPUT_DIR="$stage" \
 PIXSEAL_BUILD102_REFERENCE_TSV="$ROOT/docs/qualified-baselines/build101-runs.tsv" \
 PIXSEAL_BUILD102_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}" \
 PIXSEAL_BUILD102_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}" \
 PIXSEAL_BUILD102_CLOCK_TICKS="$ticks" \
 PIXSEAL_BUILD102_COOLDOWN_MS="${V4_PHONE_BUILD102_COOLDOWN_MS:-5000}" \
 "$stage/build102-watermark.test" -test.run '^TestExperimentalV4Build102RetainedReplay$' -test.count=1 -test.v -test.timeout=0 2>&1 | tee "$stage/build102-raw.txt"
awk -F '\t' '$1=="build102-header" || $1=="build102-row" {sub(/^[^\t]*\t/,"");print}' "$stage/build102-raw.txt" > "$stage/build102-replays.tsv"
awk -F '\t' '$1=="build102-capture-header" || $1=="build102-capture" {sub(/^[^\t]*\t/,"");print}' "$stage/build102-raw.txt" > "$stage/build102-captures.tsv"
awk -F '\t' 'NR==1{w=NF;next}NF!=w{bad=1}END{if(NR!=21||bad)exit 1}' "$stage/build102-replays.tsv"
awk -F '\t' 'NR==1{w=NF;next}NF!=w{bad=1}END{if(NR!=5||bad)exit 1}' "$stage/build102-captures.tsv"
for profile in "$stage"/*.cpu.prof; do
 env GOTOOLCHAIN="$TOOLCHAIN" go tool pprof -top -cum -nodecount=30 "$stage/build102-watermark.test" "$profile" > "${profile%.cpu.prof}-cpu.txt"
done
sha256sum -c "$stage/build102-source.sha256" > "$stage/build102-source-verified.txt"
sha256sum -c "$stage/build102-corpus.sha256" > "$stage/build102-corpus-verified.txt"
mv -T "$stage" "$OUTPUT_DIR"
published=true
echo "Build102 frozen-bank observability gate: PASS; evidence: $OUTPUT_DIR"
