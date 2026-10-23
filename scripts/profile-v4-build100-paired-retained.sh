#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD100_PROFILE_DIR:-v4-phone private/build100-paired-profile}"
TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
export GOMAXPROCS="${GOMAXPROCS:-8}"
[[ ! -e "$OUTPUT_DIR" && ! -L "$OUTPUT_DIR" ]] || { echo 'error: existing retained output is preserved; select another directory' >&2; exit 1; }
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"
for image in phone-b-mild.jpg phone-b-angle.jpg; do [[ -f "$ACQUISITION_DIR/$image" ]] || { echo "error: missing retained acquisition $image" >&2; exit 1; }; done
parent="$(dirname "$OUTPUT_DIR")"; base="$(basename "$OUTPUT_DIR")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.${base}.build100.XXXXXX")"
published=false
cleanup(){
  local rc=$?
  trap - EXIT
  if [[ "$published" != true && -d "$stage" ]]; then mv "$stage" "${OUTPUT_DIR}-failed-$(basename "$stage")"; fi
  exit "$rc"
}
trap cleanup EXIT
{
 echo 'PixSeal Build100 contemporaneous Build98/Build100 retained AB/BA'
 echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
 echo "toolchain=$TOOLCHAIN"
 echo "gomaxprocs=$GOMAXPROCS"
 echo 'qualified_baseline=Build98'
 echo 'candidate=Build100 ordered-parallel public pair scoring'
 echo 'sequence=mild 98-100; angle 98-100; angle 100-98; mild 100-98'
 echo 'additional_image_sized_memory_bytes=0'
 echo 'promotion_rule=freeze at least 5% faster in every paired comparison and pair scoring faster in every pair; exact semantics; no systematic unmodified-stage regression'
 echo 'qualification=paired decision evidence only; two independent complete nine-photo gates still required for promotion'
 env GOTOOLCHAIN="$TOOLCHAIN" go version
} > "$stage/build100-metadata.txt"
sha256sum watermark/experimental_v4_phone_build100_*.go watermark/experimental_v4_build100*go > "$stage/build100-source.sha256"
timeout --foreground "${V4_PHONE_BUILD100_TIMEOUT:-7200}s" env GOTOOLCHAIN="$TOOLCHAIN" \
 PIXSEAL_BUILD100_ACQUISITION_DIR="$ACQUISITION_DIR" \
 PIXSEAL_BUILD100_KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}" \
 PIXSEAL_BUILD100_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}" \
 PIXSEAL_BUILD100_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}" \
 PIXSEAL_BUILD100_MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}" \
 PIXSEAL_BUILD100_COOLDOWN_MS="${V4_PHONE_BUILD100_COOLDOWN_MS:-5000}" \
 go test ./watermark -run '^TestExperimentalV4Build100PairedRetainedProfile$' -count=1 -v -timeout=0 2>&1 | tee "$stage/build100-paired-raw.txt"
grep $'^build100-header\t' "$stage/build100-paired-raw.txt" | cut -f2- > "$stage/build100-paired.tsv"
grep $'^build100-row\t' "$stage/build100-paired-raw.txt" | cut -f2- >> "$stage/build100-paired.tsv"
awk -F '\t' 'NR==1{width=NF;next}NF!=width{bad=1}END{if(NR!=9||bad)exit 1}' "$stage/build100-paired.tsv"
mv -T "$stage" "$OUTPUT_DIR"
published=true
echo "Build100 paired semantic gate: PASS; review all four pairs in $OUTPUT_DIR/build100-paired.tsv before any promotion gate."
