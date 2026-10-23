#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD98_PROFILE_DIR:-v4-phone private/build98-paired-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PAIR_TIMEOUT="${V4_PHONE_BUILD98_TIMEOUT:-7200}"
COOLDOWN_MS="${V4_PHONE_BUILD98_COOLDOWN_MS:-5000}"
GO_TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-mild.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-mild.jpg" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-angle.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-angle.jpg" >&2; exit 1; }
[[ "$PAIR_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD98_TIMEOUT must be an integer number of seconds" >&2; exit 1; }
[[ "$COOLDOWN_MS" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD98_COOLDOWN_MS must be a non-negative integer" >&2; exit 1; }

ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"
output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build98.XXXXXX")"
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
      echo "Build98 FAILED paired evidence preserved in: $failed_output_dir" >&2
    else
      rm -rf "$stage_dir"
    fi
  fi
  exit "$rc"
}
trap cleanup_stage EXIT

raw="$stage_dir/build98-paired-raw.txt"
tsv="$stage_dir/build98-paired.tsv"
md="$stage_dir/build98-paired.md"
meta="$stage_dir/build98-metadata.txt"

{
  echo "PixSeal Build98 contemporaneous Build93/Build98 paired retained profile"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go_toolchain=$GO_TOOLCHAIN"
  echo "host=$(uname -a)"
  echo "go_version=$(env GOTOOLCHAIN="$GO_TOOLCHAIN" go version 2>/dev/null || true)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "acquisition_dir=$ACQUISITION_DIR"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "cooldown_ms=$COOLDOWN_MS"
  echo "sequence=cycle1 B/mild 93->98; cycle1 B/angle 93->98; cycle2 B/angle 98->93; cycle2 B/mild 98->93"
  echo "current_qualified_baseline=Build93"
  echo "candidate=Build98 Exact Generation3 Fetch Extension"
  echo "selection_evidence=Build98 public full-Generation3 benchmark"
  echo "purpose=paired retained decision evidence; not a promotion or physical qualification gate"
  echo "additional_image_sized_memory_bytes=0"
  echo "telemetry_remediation=historical Build84 Fetch counters remain continuation4-only; Build98 generation-three uses separate Gen3Fetch counters"
  echo "build93_candidate_sha256=$(sha256sum watermark/experimental_v4_phone_build93_candidate.go | awk '{print $1}')"
  echo "build98_bank_sha256=$(sha256sum watermark/experimental_v4_phone_build98_bank.go | awk '{print $1}')"
  echo "build98_fetch_sha256=$(sha256sum watermark/experimental_v4_phone_build98_gen3_fetch.go | awk '{print $1}')"
  echo "build98_recover_sha256=$(sha256sum watermark/experimental_v4_phone_build98_recover.go | awk '{print $1}')"
} > "$meta"

set +e
timeout --foreground "${PAIR_TIMEOUT}s" env \
  GOTOOLCHAIN="$GO_TOOLCHAIN" \
  PIXSEAL_BUILD98_ACQUISITION_DIR="$ACQUISITION_DIR" \
  PIXSEAL_BUILD98_KEY="$KEY" \
  PIXSEAL_BUILD98_WIDTH="$CANONICAL_WIDTH" \
  PIXSEAL_BUILD98_HEIGHT="$CANONICAL_HEIGHT" \
  PIXSEAL_BUILD98_MESSAGE_B="$MESSAGE_B" \
  PIXSEAL_BUILD98_COOLDOWN_MS="$COOLDOWN_MS" \
  go test ./watermark -run '^TestExperimentalV4Build98PairedRetainedProfile$' -count=1 -v -timeout=0 2>&1 | tee "$raw"
rc=${PIPESTATUS[0]}
set -e
if (( rc != 0 )); then exit "$rc"; fi

header=(cycle image order position variant outcome wall_ms total_ms geometry_ms plane_ms freeze_ms prefix1_ms gen2_ms gen3_ms gen4_ms qualification_ms decode_ms geometry_evals bank qualification_evals qualified decode_candidates list_frames physical_decode speculative_decode physical_profiles physical_list_frames protected_sampling_passes protected_profile_derivations authenticated qualification_workers qualification_tasks fold_scores block_reads block_success block_failed gen3_fetch_fold_scores gen3_fetch_block_reads gen3_fetch_block_success gen3_fetch_block_failed)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"
grep $'^build98-row\t' "$raw" | cut -f2- >> "$tsv" || true
rows=$(( $(wc -l < "$tsv") - 1 ))
if (( rows != 8 )); then echo "error: Build98 expected 8 paired rows, found $rows" >&2; exit 1; fi

awk -F '\t' -v expected_fields="${#header[@]}" '
NR==1 { next }
NF != expected_fields { printf "error: Build98 TSV field mismatch on row %d: got %d want %d\n", NR, NF, expected_fields > "/dev/stderr"; bad=1 }
END { exit bad }
' "$tsv"

for cycle in 1 2; do
  for image in phone-b-mild.jpg phone-b-angle.jpg; do
    count93="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build93"{n++}END{print n+0}' "$tsv")"
    count98="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build98"{n++}END{print n+0}' "$tsv")"
    [[ "$count93" -eq 1 && "$count98" -eq 1 ]] || { echo "error: Build98 incomplete pair cycle=$cycle image=$image build93=$count93 build98=$count98" >&2; exit 1; }
  done
done

ratio(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$b" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.3f", a/b}'; else printf '-'; fi; }
pct(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$a" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.2f", 100*(a-b)/a}'; else printf '-'; fi; }
cell(){ local cycle="$1" image="$2" variant="$3" column="$4"; awk -F '\t' -v c="$cycle" -v im="$image" -v v="$variant" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==c&&$2==im&&$5==v{print $col;exit}' "$tsv"; }

{
  echo '# PixSeal Build98 contemporaneous Build93/Build98 paired retained profile'; echo
  echo 'Build98 remains a candidate. Each pair reuses the same decoded in-memory working image. Cycle 2 reverses both candidate and image order. Exact telemetry/outcome equality is the semantic gate; timing is decision evidence only.'; echo
  echo '| cycle | image | order | Build93 total | Build98 total | total reduction | Build93 geometry | Build98 geometry | geometry reduction | Build93 gen3 | Build98 gen3 | gen3 speedup | Build93 qualification | Build98 qualification | Build93 decode | Build98 decode |'
  echo '|---:|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  for cycle in 1 2; do
    for image in phone-b-mild.jpg phone-b-angle.jpg; do
      order="$(cell "$cycle" "$image" build93 order)"
      b93t="$(cell "$cycle" "$image" build93 total_ms)"; b98t="$(cell "$cycle" "$image" build98 total_ms)"
      b93g="$(cell "$cycle" "$image" build93 geometry_ms)"; b98g="$(cell "$cycle" "$image" build98 geometry_ms)"
      b93g3="$(cell "$cycle" "$image" build93 gen3_ms)"; b98g3="$(cell "$cycle" "$image" build98 gen3_ms)"
      echo "| $cycle | $image | $order | $b93t | $b98t | $(pct "$b93t" "$b98t")% | $b93g | $b98g | $(pct "$b93g" "$b98g")% | $b93g3 | $b98g3 | $(ratio "$b93g3" "$b98g3")x | $(cell "$cycle" "$image" build93 qualification_ms) | $(cell "$cycle" "$image" build98 qualification_ms) | $(cell "$cycle" "$image" build93 decode_ms) | $(cell "$cycle" "$image" build98 decode_ms) |"
    done
  done
  echo
  echo 'Promotion review criterion: Build98 should show a material, directionally stable Gen3/geometry gain in both AB/BA cycles. Isolated microbenchmarks are not sufficient. Build98 adds zero image-sized memory.'; echo
  echo 'Frozen semantics: B/mild 79,259 geometry evaluations / 937 bank / 935 qualified / 691 logical decode / 696 physical decode / authenticated PASS; B/angle 334,857 geometry evaluations / 6,198 bank / 0 qualified / REJECT. Historical fold/block counters remain continuation4-only and must match Build93 row-for-row; Build98 generation-three work is reported separately as gen3_fetch_* telemetry.'; echo
} > "$md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build98 paired retained profile written to:'
echo "  $OUTPUT_DIR/build98-paired.tsv"
echo "  $OUTPUT_DIR/build98-paired.md"
echo "  $OUTPUT_DIR/build98-paired-raw.txt"
echo "  $OUTPUT_DIR/build98-metadata.txt"
echo 'Build98 paired semantic gate: PASS (timing remains review evidence; candidate not promoted)'
