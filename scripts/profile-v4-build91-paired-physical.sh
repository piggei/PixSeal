#!/usr/bin/env bash
set -euo pipefail

ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD91_PROFILE_DIR:-v4-phone private/build91-paired-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PAIR_TIMEOUT="${V4_PHONE_BUILD91_TIMEOUT:-7200}"
COOLDOWN_MS="${V4_PHONE_BUILD91_COOLDOWN_MS:-5000}"
GO_TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-mild.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-mild.jpg" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-angle.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-angle.jpg" >&2; exit 1; }
# go test runs package tests with the package directory as the working directory.
# Convert the repo-relative retained-corpus path to an absolute path before
# passing it through PIXSEAL_BUILD91_ACQUISITION_DIR, otherwise the Go test
# would resolve it below ./watermark and report a false ENOENT.
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"
[[ "$PAIR_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD91_TIMEOUT must be an integer number of seconds" >&2; exit 1; }
[[ "$COOLDOWN_MS" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD91_COOLDOWN_MS must be a non-negative integer" >&2; exit 1; }

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build91.XXXXXX")"
published=false
cleanup_stage(){ if [[ "$published" != true && -d "$stage_dir" ]]; then rm -rf "$stage_dir"; fi; }
trap cleanup_stage EXIT
failed_output_dir="${OUTPUT_DIR}-failed"
raw="$stage_dir/build91-paired-raw.txt"
tsv="$stage_dir/build91-paired.tsv"
md="$stage_dir/build91-paired.md"
meta="$stage_dir/build91-metadata.txt"

{
  echo "PixSeal Build91 contemporaneous Build84/Build90 paired physical profile"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go_toolchain=$GO_TOOLCHAIN"
  echo "host=$(uname -a)"
  echo "go_version=$(env GOTOOLCHAIN="$GO_TOOLCHAIN" go version 2>/dev/null || true)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "acquisition_dir=$ACQUISITION_DIR"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "cooldown_ms=$COOLDOWN_MS"
  echo "sequence=cycle1 B/mild 84->90; cycle1 B/angle 84->90; cycle2 B/angle 90->84; cycle2 B/mild 90->84"
  echo "purpose=paired physical promotion evidence; Build90 is current qualified baseline and Build84 is previous qualified baseline"
} > "$meta"

set +e
timeout --foreground "${PAIR_TIMEOUT}s" env \
  GOTOOLCHAIN="$GO_TOOLCHAIN" \
  PIXSEAL_BUILD91_ACQUISITION_DIR="$ACQUISITION_DIR" \
  PIXSEAL_BUILD91_KEY="$KEY" \
  PIXSEAL_BUILD91_WIDTH="$CANONICAL_WIDTH" \
  PIXSEAL_BUILD91_HEIGHT="$CANONICAL_HEIGHT" \
  PIXSEAL_BUILD91_MESSAGE_B="$MESSAGE_B" \
  PIXSEAL_BUILD91_COOLDOWN_MS="$COOLDOWN_MS" \
  go test ./watermark -run '^TestExperimentalV4Build91PairedPhysicalProfile$' -count=1 -v -timeout=0 2>&1 | tee "$raw"
rc=${PIPESTATUS[0]}
set -e
if (( rc != 0 )); then
  rm -rf "$failed_output_dir"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  echo "Build91 FAILED paired profile preserved in: $failed_output_dir" >&2
  exit "$rc"
fi

header=(cycle image order position variant outcome wall_ms total_ms geometry_ms plane_ms freeze_ms prefix1_ms gen2_ms gen3_ms gen4_ms qualification_ms decode_ms geometry_evals bank qualification_evals qualified decode_candidates list_frames physical_decode speculative_decode authenticated qualification_workers qualification_tasks fold_scores block_reads block_success block_failed)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"

grep $'^build91-row\t' "$raw" | cut -f2- >> "$tsv" || true
rows=$(( $(wc -l < "$tsv") - 1 ))
if (( rows != 8 )); then
  echo "error: Build91 expected 8 paired rows, found $rows" >&2
  rm -rf "$failed_output_dir"
  mv "$stage_dir" "$failed_output_dir"
  published=true
  exit 1
fi

# Validate field count and pair completeness without using timing as a correctness gate.
awk -F '\t' -v expected_fields="${#header[@]}" '
NR==1 { next }
NF != expected_fields { printf "error: Build91 TSV field mismatch on row %d: got %d want %d\n", NR, NF, expected_fields > "/dev/stderr"; bad=1 }
END { exit bad }
' "$tsv"

for cycle in 1 2; do
  for image in phone-b-mild.jpg phone-b-angle.jpg; do
    count84="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build84"{n++}END{print n+0}' "$tsv")"
    count90="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build90"{n++}END{print n+0}' "$tsv")"
    [[ "$count84" -eq 1 && "$count90" -eq 1 ]] || { echo "error: Build91 incomplete pair cycle=$cycle image=$image build84=$count84 build90=$count90" >&2; exit 1; }
  done
done

ratio(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$b" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.3f", a/b}'; else printf '-'; fi; }
pct(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$a" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.2f", 100*(a-b)/a}'; else printf '-'; fi; }
cell(){ local cycle="$1" image="$2" variant="$3" column="$4"; awk -F '\t' -v c="$cycle" -v im="$image" -v v="$variant" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==c&&$2==im&&$5==v{print $col;exit}' "$tsv"; }

{
  echo '# PixSeal Build91 contemporaneous Build84/Build90 paired physical profile'; echo
  echo 'Build91 is observability-only. It does not change the active Build90 candidate or the qualified Build84 implementation. Each pair decodes the same retained physical image once into the same in-memory working image, then runs Build84 and Build90 deep recovery consecutively. Cycle 2 reverses both variant order and case order to reduce bias from host drift.'; echo
  echo '| cycle | image | order | Build84 total | Build90 total | total reduction | Build84 geometry | Build90 geometry | geometry ratio 84/90 | Build84 qualification | Build90 qualification | qualification speedup | Build84 decode | Build90 decode |'
  echo '|---:|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  for cycle in 1 2; do
    for image in phone-b-mild.jpg phone-b-angle.jpg; do
      order="$(cell "$cycle" "$image" build84 order)"
      b84t="$(cell "$cycle" "$image" build84 total_ms)"; b90t="$(cell "$cycle" "$image" build90 total_ms)"
      b84g="$(cell "$cycle" "$image" build84 geometry_ms)"; b90g="$(cell "$cycle" "$image" build90 geometry_ms)"
      b84q="$(cell "$cycle" "$image" build84 qualification_ms)"; b90q="$(cell "$cycle" "$image" build90 qualification_ms)"
      b84d="$(cell "$cycle" "$image" build84 decode_ms)"; b90d="$(cell "$cycle" "$image" build90 decode_ms)"
      echo "| $cycle | $image | $order | $b84t | $b90t | $(pct "$b84t" "$b90t")% | $b84g | $b90g | $(ratio "$b84g" "$b90g")x | $b84q | $b90q | $(ratio "$b84q" "$b90q")x | $b84d | $b90d |"
    done
  done
  echo
  echo '## Raw authoritative rows'; echo
  echo 'The TSV is authoritative for exact telemetry. Timing is evidence only and is not part of the semantic PASS gate. Build90 may be promoted only after reviewing both paired cycles together with the prior public benchmark and physical semantic runs.'; echo
  echo 'Expected retained-corpus semantics are enforced inside the Go test: B/mild must remain 79,259 geometry evaluations / 937 bank / 935 qualified / 691 logical decode candidates / 2,120,047 logical list frames / 696 physical decode candidates / 5 speculative and authenticated payload PASS; B/angle must remain 334,857 / 6,198 / 0 and reject.'
} > "$md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build91 paired physical profile written to:'
echo "  $OUTPUT_DIR/build91-paired.tsv"
echo "  $OUTPUT_DIR/build91-paired.md"
echo "  $OUTPUT_DIR/build91-paired-raw.txt"
echo "  $OUTPUT_DIR/build91-metadata.txt"
echo 'Build91 paired semantic gate: PASS'
