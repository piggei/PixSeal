#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD97_PROFILE_DIR:-v4-phone private/build97-paired-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PAIR_TIMEOUT="${V4_PHONE_BUILD97_TIMEOUT:-7200}"
COOLDOWN_MS="${V4_PHONE_BUILD97_COOLDOWN_MS:-5000}"
GO_TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-mild.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-mild.jpg" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-angle.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-angle.jpg" >&2; exit 1; }
[[ "$PAIR_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD97_TIMEOUT must be an integer number of seconds" >&2; exit 1; }
[[ "$COOLDOWN_MS" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD97_COOLDOWN_MS must be a non-negative integer" >&2; exit 1; }

ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"
output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build97.XXXXXX")"
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
      echo "Build97 FAILED paired evidence preserved in: $failed_output_dir" >&2
    else
      rm -rf "$stage_dir"
    fi
  fi
  exit "$rc"
}
trap cleanup_stage EXIT

raw="$stage_dir/build97-paired-raw.txt"
tsv="$stage_dir/build97-paired.tsv"
md="$stage_dir/build97-paired.md"
meta="$stage_dir/build97-metadata.txt"

{
  echo "PixSeal Build97 contemporaneous Build93/Build97 paired retained profile"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go_toolchain=$GO_TOOLCHAIN"
  echo "host=$(uname -a)"
  echo "go_version=$(env GOTOOLCHAIN="$GO_TOOLCHAIN" go version 2>/dev/null || true)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "acquisition_dir=$ACQUISITION_DIR"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "cooldown_ms=$COOLDOWN_MS"
  echo "sequence=cycle1 B/mild 93->97; cycle1 B/angle 93->97; cycle2 B/angle 97->93; cycle2 B/mild 97->93"
  echo "current_qualified_baseline=Build93"
  echo "candidate=Build97 Exact Gen4 Probe Fetch Extension"
  echo "selection_evidence=Build97 public full-Generation4 benchmark"
  echo "purpose=paired retained decision evidence; not a promotion or physical qualification gate"
  echo "additional_image_sized_memory_bytes=0"
  echo "telemetry_remediation=historical Build84 Fetch counters remain continuation4-only; Build97 single4/pair4 use separate ProbeFetch counters"
  echo "build93_candidate_sha256=$(sha256sum watermark/experimental_v4_phone_build93_candidate.go | awk '{print $1}')"
  echo "build97_bank_sha256=$(sha256sum watermark/experimental_v4_phone_build97_bank.go | awk '{print $1}')"
  echo "build97_fetch_sha256=$(sha256sum watermark/experimental_v4_phone_build97_gen4_probe_fetch.go | awk '{print $1}')"
  echo "build97_recover_sha256=$(sha256sum watermark/experimental_v4_phone_build97_recover.go | awk '{print $1}')"
} > "$meta"

set +e
timeout --foreground "${PAIR_TIMEOUT}s" env \
  GOTOOLCHAIN="$GO_TOOLCHAIN" \
  PIXSEAL_BUILD97_ACQUISITION_DIR="$ACQUISITION_DIR" \
  PIXSEAL_BUILD97_KEY="$KEY" \
  PIXSEAL_BUILD97_WIDTH="$CANONICAL_WIDTH" \
  PIXSEAL_BUILD97_HEIGHT="$CANONICAL_HEIGHT" \
  PIXSEAL_BUILD97_MESSAGE_B="$MESSAGE_B" \
  PIXSEAL_BUILD97_COOLDOWN_MS="$COOLDOWN_MS" \
  go test ./watermark -run '^TestExperimentalV4Build97PairedRetainedProfile$' -count=1 -v -timeout=0 2>&1 | tee "$raw"
rc=${PIPESTATUS[0]}
set -e
if (( rc != 0 )); then exit "$rc"; fi

header=(cycle image order position variant outcome wall_ms total_ms geometry_ms plane_ms freeze_ms prefix1_ms gen2_ms gen3_ms gen4_ms qualification_ms decode_ms geometry_evals bank qualification_evals qualified decode_candidates list_frames physical_decode speculative_decode physical_profiles physical_list_frames protected_sampling_passes protected_profile_derivations authenticated qualification_workers qualification_tasks fold_scores block_reads block_success block_failed probe_fold_scores probe_block_reads probe_block_success probe_block_failed)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"
grep $'^build97-row\t' "$raw" | cut -f2- >> "$tsv" || true
rows=$(( $(wc -l < "$tsv") - 1 ))
if (( rows != 8 )); then echo "error: Build97 expected 8 paired rows, found $rows" >&2; exit 1; fi

awk -F '\t' -v expected_fields="${#header[@]}" '
NR==1 { next }
NF != expected_fields { printf "error: Build97 TSV field mismatch on row %d: got %d want %d\n", NR, NF, expected_fields > "/dev/stderr"; bad=1 }
END { exit bad }
' "$tsv"

for cycle in 1 2; do
  for image in phone-b-mild.jpg phone-b-angle.jpg; do
    count93="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build93"{n++}END{print n+0}' "$tsv")"
    count97="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build97"{n++}END{print n+0}' "$tsv")"
    [[ "$count93" -eq 1 && "$count97" -eq 1 ]] || { echo "error: Build97 incomplete pair cycle=$cycle image=$image build93=$count93 build97=$count97" >&2; exit 1; }
  done
done

ratio(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$b" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.3f", a/b}'; else printf '-'; fi; }
pct(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$a" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.2f", 100*(a-b)/a}'; else printf '-'; fi; }
cell(){ local cycle="$1" image="$2" variant="$3" column="$4"; awk -F '\t' -v c="$cycle" -v im="$image" -v v="$variant" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==c&&$2==im&&$5==v{print $col;exit}' "$tsv"; }

{
  echo '# PixSeal Build97 contemporaneous Build93/Build97 paired retained profile'; echo
  echo 'Build97 remains a candidate. Each pair reuses the same decoded in-memory working image. Cycle 2 reverses both candidate and image order. Exact telemetry/outcome equality is the semantic gate; timing is decision evidence only.'; echo
  echo '| cycle | image | order | Build93 total | Build97 total | total reduction | Build93 geometry | Build97 geometry | geometry reduction | Build93 gen4 | Build97 gen4 | gen4 speedup | Build93 qualification | Build97 qualification | Build93 decode | Build97 decode |'
  echo '|---:|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  for cycle in 1 2; do
    for image in phone-b-mild.jpg phone-b-angle.jpg; do
      order="$(cell "$cycle" "$image" build93 order)"
      b93t="$(cell "$cycle" "$image" build93 total_ms)"; b97t="$(cell "$cycle" "$image" build97 total_ms)"
      b93g="$(cell "$cycle" "$image" build93 geometry_ms)"; b97g="$(cell "$cycle" "$image" build97 geometry_ms)"
      b93g4="$(cell "$cycle" "$image" build93 gen4_ms)"; b97g4="$(cell "$cycle" "$image" build97 gen4_ms)"
      echo "| $cycle | $image | $order | $b93t | $b97t | $(pct "$b93t" "$b97t")% | $b93g | $b97g | $(pct "$b93g" "$b97g")% | $b93g4 | $b97g4 | $(ratio "$b93g4" "$b97g4")x | $(cell "$cycle" "$image" build93 qualification_ms) | $(cell "$cycle" "$image" build97 qualification_ms) | $(cell "$cycle" "$image" build93 decode_ms) | $(cell "$cycle" "$image" build97 decode_ms) |"
    done
  done
  echo
  echo 'Promotion review criterion: Build97 should show a material, directionally stable Gen4/geometry gain in both AB/BA cycles. Isolated microbenchmarks are not sufficient. Build97 adds zero image-sized memory.'; echo
  echo 'Frozen semantics: B/mild 79,259 geometry evaluations / 937 bank / 935 qualified / 691 logical decode / 696 physical decode / authenticated PASS; B/angle 334,857 geometry evaluations / 6,198 bank / 0 qualified / REJECT. Historical fold/block counters remain continuation4-only and must match Build93 row-for-row; Build97 single4/pair4 work is reported separately as probe_* telemetry.'; echo
} > "$md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build97 paired retained profile written to:'
echo "  $OUTPUT_DIR/build97-paired.tsv"
echo "  $OUTPUT_DIR/build97-paired.md"
echo "  $OUTPUT_DIR/build97-paired-raw.txt"
echo "  $OUTPUT_DIR/build97-metadata.txt"
echo 'Build97 paired semantic gate: PASS (timing remains review evidence; candidate not promoted)'
