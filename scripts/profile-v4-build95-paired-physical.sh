#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD95_PROFILE_DIR:-v4-phone private/build95-paired-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PAIR_TIMEOUT="${V4_PHONE_BUILD95_TIMEOUT:-7200}"
COOLDOWN_MS="${V4_PHONE_BUILD95_COOLDOWN_MS:-5000}"
GO_TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ -d "$ACQUISITION_DIR" ]] || { echo "error: phone acquisition directory not found: $ACQUISITION_DIR" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-mild.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-mild.jpg" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-angle.jpg" ]] || { echo "error: missing $ACQUISITION_DIR/phone-b-angle.jpg" >&2; exit 1; }
[[ "$PAIR_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD95_TIMEOUT must be an integer number of seconds" >&2; exit 1; }
[[ "$COOLDOWN_MS" =~ ^[0-9]+$ ]] || { echo "error: V4_PHONE_BUILD95_COOLDOWN_MS must be a non-negative integer" >&2; exit 1; }

# Package tests execute below ./watermark, so retain an absolute corpus path.
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build95.XXXXXX")"
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
      echo "Build95 FAILED paired evidence preserved in: $failed_output_dir" >&2
    else
      rm -rf "$stage_dir"
    fi
  fi
  exit "$rc"
}
trap cleanup_stage EXIT

raw="$stage_dir/build95-paired-raw.txt"
tsv="$stage_dir/build95-paired.tsv"
md="$stage_dir/build95-paired.md"
meta="$stage_dir/build95-metadata.txt"

{
  echo "PixSeal Build95 contemporaneous Build93/Build95 paired retained profile"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go_toolchain=$GO_TOOLCHAIN"
  echo "host=$(uname -a)"
  echo "go_version=$(env GOTOOLCHAIN="$GO_TOOLCHAIN" go version 2>/dev/null || true)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "acquisition_dir=$ACQUISITION_DIR"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "cooldown_ms=$COOLDOWN_MS"
  echo "sequence=cycle1 B/mild 93->95; cycle1 B/angle 93->95; cycle2 B/angle 95->93; cycle2 B/mild 95->93"
  echo "current_qualified_baseline=Build93"
  echo "candidate=Build95 Exact Gen4 Luminance Plane"
  echo "selection_evidence=Build94 qualified Build93 post-promotion profile"
  echo "purpose=paired retained decision evidence; not a promotion or physical qualification gate"
  echo "build93_candidate_sha256=$(sha256sum watermark/experimental_v4_phone_build93_candidate.go | awk '{print $1}')"
  echo "build95_candidate_sha256=$(sha256sum watermark/experimental_v4_phone_build95_candidate.go | awk '{print $1}')"
  echo "build95_luminance_sha256=$(sha256sum watermark/experimental_v4_phone_build95_gen4_luminance.go | awk '{print $1}')"
} > "$meta"

set +e
timeout --foreground "${PAIR_TIMEOUT}s" env \
  GOTOOLCHAIN="$GO_TOOLCHAIN" \
  PIXSEAL_BUILD95_ACQUISITION_DIR="$ACQUISITION_DIR" \
  PIXSEAL_BUILD95_KEY="$KEY" \
  PIXSEAL_BUILD95_WIDTH="$CANONICAL_WIDTH" \
  PIXSEAL_BUILD95_HEIGHT="$CANONICAL_HEIGHT" \
  PIXSEAL_BUILD95_MESSAGE_B="$MESSAGE_B" \
  PIXSEAL_BUILD95_COOLDOWN_MS="$COOLDOWN_MS" \
  go test ./watermark -run '^TestExperimentalV4Build95PairedPhysicalProfile$' -count=1 -v -timeout=0 2>&1 | tee "$raw"
rc=${PIPESTATUS[0]}
set -e
if (( rc != 0 )); then
  exit "$rc"
fi

header=(cycle image order position variant outcome wall_ms total_ms geometry_ms plane_ms freeze_ms prefix1_ms gen2_ms gen3_ms gen4_ms luminance_prep_ms luminance_bytes luminance_used luminance_fallback qualification_ms decode_ms geometry_evals bank qualification_evals qualified decode_candidates list_frames physical_decode speculative_decode physical_profiles physical_list_frames protected_sampling_passes protected_profile_derivations authenticated qualification_workers qualification_tasks fold_scores block_reads block_success block_failed)
(IFS=$'\t'; echo "${header[*]}") > "$tsv"
grep $'^build95-row\t' "$raw" | cut -f2- >> "$tsv" || true

rows=$(( $(wc -l < "$tsv") - 1 ))
if (( rows != 8 )); then
  echo "error: Build95 expected 8 paired rows, found $rows" >&2
  exit 1
fi

awk -F '\t' -v expected_fields="${#header[@]}" '
NR==1 { next }
NF != expected_fields { printf "error: Build95 TSV field mismatch on row %d: got %d want %d\n", NR, NF, expected_fields > "/dev/stderr"; bad=1 }
END { exit bad }
' "$tsv"

for cycle in 1 2; do
  for image in phone-b-mild.jpg phone-b-angle.jpg; do
    count93="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build93"{n++}END{print n+0}' "$tsv")"
    count95="$(awk -F '\t' -v c="$cycle" -v im="$image" 'NR>1&&$1==c&&$2==im&&$5=="build95"{n++}END{print n+0}' "$tsv")"
    [[ "$count93" -eq 1 && "$count95" -eq 1 ]] || { echo "error: Build95 incomplete pair cycle=$cycle image=$image build93=$count93 build95=$count95" >&2; exit 1; }
  done
done

# The Go test already checks exact semantic equality. Recheck candidate-specific
# memory telemetry in the published table so malformed output cannot be mistaken
# for a completed paired run.
awk -F '\t' '
NR==1 { for(i=1;i<=NF;i++) c[$i]=i; next }
$c["variant"]=="build93" && ($c["luminance_used"]!="false" || $c["luminance_fallback"]!="false" || $c["luminance_bytes"]!=0) { bad=1 }
$c["variant"]=="build95" && ($c["luminance_used"]!="true" || $c["luminance_fallback"]!="false" || $c["luminance_bytes"]<=0) { bad=1 }
END { if (bad) { print "error: Build95 luminance telemetry validation failed" > "/dev/stderr"; exit 1 } }
' "$tsv"

ratio(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$b" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.3f", a/b}'; else printf '-'; fi; }
pct(){ local a="$1" b="$2"; if [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ && "$a" -gt 0 ]]; then awk -v a="$a" -v b="$b" 'BEGIN{printf "%.2f", 100*(a-b)/a}'; else printf '-'; fi; }
cell(){ local cycle="$1" image="$2" variant="$3" column="$4"; awk -F '\t' -v c="$cycle" -v im="$image" -v v="$variant" -v wanted="$column" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)col=i;next}$1==c&&$2==im&&$5==v{print $col;exit}' "$tsv"; }

{
  echo '# PixSeal Build95 contemporaneous Build93/Build95 paired retained profile'; echo
  echo 'Build95 is still a candidate. Each row pair reuses the same decoded in-memory working image, then executes the qualified Build93 deep runtime and Build95 consecutively. Cycle 2 reverses both candidate order and image order. Timing is decision evidence only; exact telemetry, outcome and authenticated payload/reject equivalence are the semantic gate.'; echo
  echo '| cycle | image | order | Build93 total | Build95 total | total reduction | Build93 geometry | Build95 geometry | geometry reduction | Build93 gen4 | Build95 gen4 | gen4 speedup | luma prep | extra bytes | Build93 qualification | Build95 qualification | Build93 decode | Build95 decode |'
  echo '|---:|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  for cycle in 1 2; do
    for image in phone-b-mild.jpg phone-b-angle.jpg; do
      order="$(cell "$cycle" "$image" build93 order)"
      b93t="$(cell "$cycle" "$image" build93 total_ms)"; b95t="$(cell "$cycle" "$image" build95 total_ms)"
      b93g="$(cell "$cycle" "$image" build93 geometry_ms)"; b95g="$(cell "$cycle" "$image" build95 geometry_ms)"
      b93g4="$(cell "$cycle" "$image" build93 gen4_ms)"; b95g4="$(cell "$cycle" "$image" build95 gen4_ms)"
      echo "| $cycle | $image | $order | $b93t | $b95t | $(pct "$b93t" "$b95t")% | $b93g | $b95g | $(pct "$b93g" "$b95g")% | $b93g4 | $b95g4 | $(ratio "$b93g4" "$b95g4")x | $(cell "$cycle" "$image" build95 luminance_prep_ms) | $(cell "$cycle" "$image" build95 luminance_bytes) | $(cell "$cycle" "$image" build93 qualification_ms) | $(cell "$cycle" "$image" build95 qualification_ms) | $(cell "$cycle" "$image" build93 decode_ms) | $(cell "$cycle" "$image" build95 decode_ms) |"
    done
  done
  echo
  echo 'The candidate is promotion-worthy only if both AB/BA cycles show a material, directionally stable Gen4/geometry gain that repays luminance-plane preparation and added memory. Cross-session Build94/Build93 timings are context only because host drift is already documented.'; echo
  echo 'Exact retained semantics enforced inside the Go test: B/mild remains 79,259 geometry evaluations, 937-bank, 935 qualified, 691 logical decode candidates, 2,120,047 logical list frames, 696 physical decodes, 2,086 physical profiles/derivations and authenticated payload PASS; B/angle remains 334,857 / 6,198 / 0 and REJECT. Build95 must match Build93 row-for-row.'; echo
} > "$md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build95 paired retained profile written to:'
echo "  $OUTPUT_DIR/build95-paired.tsv"
echo "  $OUTPUT_DIR/build95-paired.md"
echo "  $OUTPUT_DIR/build95-paired-raw.txt"
echo "  $OUTPUT_DIR/build95-metadata.txt"
echo 'Build95 paired semantic gate: PASS (timing remains review evidence; candidate not promoted)'
