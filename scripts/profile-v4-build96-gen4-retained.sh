#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
OUTPUT_DIR="${V4_PHONE_BUILD96_PROFILE_DIR:-v4-phone private/build96-profile}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PROFILE_TIMEOUT="${V4_PHONE_BUILD96_TIMEOUT:-7200}"
COOLDOWN_MS="${V4_PHONE_BUILD96_COOLDOWN_MS:-5000}"
GO_TOOLCHAIN="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v sha256sum >/dev/null 2>&1 || { echo "error: sha256sum is required" >&2; exit 1; }
[[ -f "$ACQUISITION_DIR/phone-b-angle.jpg" ]] || { echo "error: missing B/angle corpus image" >&2; exit 1; }
[[ "$PROFILE_TIMEOUT" =~ ^[0-9]+$ ]] || { echo "error: timeout must be integer seconds" >&2; exit 1; }
[[ "$COOLDOWN_MS" =~ ^[0-9]+$ ]] || { echo "error: cooldown must be non-negative integer ms" >&2; exit 1; }
ACQUISITION_DIR="$(cd "$ACQUISITION_DIR" && pwd -P)"

output_parent="$(dirname "$OUTPUT_DIR")"
output_base="$(basename "$OUTPUT_DIR")"
mkdir -p "$output_parent"
stage_dir="$(mktemp -d "$output_parent/.${output_base}.build96.XXXXXX")"
stage_dir="$(cd "$stage_dir" && pwd -P)"
published=false
failed_output_dir="${OUTPUT_DIR}-failed"
cleanup(){
  local rc=$?
  trap - EXIT
  set +e
  if [[ "$published" != true && -d "$stage_dir" ]]; then
    if (( rc != 0 )); then
      rm -rf "$failed_output_dir"
      mv "$stage_dir" "$failed_output_dir"
      echo "Build96 FAILED staged evidence preserved in: $failed_output_dir" >&2
    else
      rm -rf "$stage_dir"
    fi
  fi
  exit "$rc"
}
trap cleanup EXIT

testbin="$stage_dir/build96-watermark.test"
raw="$stage_dir/build96-profile-raw.txt"
tsv="$stage_dir/build96-profile.tsv"
md="$stage_dir/build96-profile.md"
meta="$stage_dir/build96-metadata.txt"

{
  echo "PixSeal Build96 Gen4 retained bottleneck profile"
  echo "date_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go_toolchain=$GO_TOOLCHAIN"
  echo "go_version=$(env GOTOOLCHAIN="$GO_TOOLCHAIN" go version 2>/dev/null || true)"
  echo "host=$(uname -a)"
  echo "host_logical_cpus=$(nproc 2>/dev/null || true)"
  echo "acquisition_dir=$ACQUISITION_DIR"
  echo "canonical=${CANONICAL_WIDTH}x${CANONICAL_HEIGHT}"
  echo "cooldown_ms=$COOLDOWN_MS"
  echo "current_qualified_baseline=Build93"
  echo "build95_decision=rejected_not_promoted"
  echo "sequence=B/angle cycle1 Build93->Build95; cycle2 Build95->Build93"
  echo "purpose=observability-only retained bottleneck localization; no runtime change"
  echo "build93_candidate_sha256=$(sha256sum watermark/experimental_v4_phone_build93_candidate.go | awk '{print $1}')"
  echo "build95_candidate_sha256=$(sha256sum watermark/experimental_v4_phone_build95_candidate.go | awk '{print $1}')"
  echo "build95_luminance_sha256=$(sha256sum watermark/experimental_v4_phone_build95_gen4_luminance.go | awk '{print $1}')"
} > "$meta"

env GOTOOLCHAIN="$GO_TOOLCHAIN" go test -c ./watermark -o "$testbin" 2>&1 | tee "$stage_dir/build96-test-binary-build.txt"
echo "test_binary_sha256=$(sha256sum "$testbin" | awk '{print $1}')" >> "$meta"

set +e
timeout --foreground "${PROFILE_TIMEOUT}s" env \
  PIXSEAL_BUILD96_ACQUISITION_DIR="$ACQUISITION_DIR" \
  PIXSEAL_BUILD96_PROFILE_DIR="$stage_dir" \
  PIXSEAL_BUILD96_KEY="$KEY" \
  PIXSEAL_BUILD96_WIDTH="$CANONICAL_WIDTH" \
  PIXSEAL_BUILD96_HEIGHT="$CANONICAL_HEIGHT" \
  PIXSEAL_BUILD96_MESSAGE_B="$MESSAGE_B" \
  PIXSEAL_BUILD96_COOLDOWN_MS="$COOLDOWN_MS" \
  "$testbin" -test.run '^TestExperimentalV4Build96Gen4RetainedBottleneckProfile$' -test.v -test.timeout=0 2>&1 | tee "$raw"
rc=${PIPESTATUS[0]}
set -e
(( rc == 0 )) || exit "$rc"

grep $'^build96-header\t' "$raw" | cut -f2- > "$tsv"
grep $'^build96-row\t' "$raw" | cut -f2- >> "$tsv" || true
rows=$(( $(wc -l < "$tsv") - 1 ))
(( rows == 4 )) || { echo "error: Build96 expected 4 rows, found $rows" >&2; exit 1; }
fields="$(awk -F '\t' 'NR==1{print NF}' "$tsv")"
awk -F '\t' -v n="$fields" 'NR>1&&NF!=n{bad=1}END{exit bad}' "$tsv" || { echo "error: Build96 malformed TSV" >&2; exit 1; }

# Reassert the frozen B/angle semantics in the emitted table.
awk -F '\t' '
NR==1{for(i=1;i<=NF;i++)c[$i]=i;next}
$c["outcome"]!="reject" || $c["geometry_evals"]!=334857 || $c["bank"]!=6198 || $c["qualification_evals"]!=6198 || $c["qualified"]!=0 || $c["fold_scores"]!=118004 || $c["block_reads"]!=151045120 || $c["block_success"]!=149370664 || $c["block_failed"]!=1674456 {bad=1}
END{exit bad}
' "$tsv" || { echo "error: Build96 semantic invariant validation failed" >&2; exit 1; }

for profile in "$stage_dir"/build96-cycle*.pprof; do
  [[ -s "$profile" ]] || { echo "error: missing profile $profile" >&2; exit 1; }
  stem="${profile%.pprof}"
  env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -nodecount=120 "$testbin" "$profile" > "${stem}-top.txt" 2>&1
  env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -cum -nodecount=120 "$testbin" "$profile" > "${stem}-cum.txt" 2>&1
  while IFS='|' read -r phase focus; do
    env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -top -nodecount=120 -focus="$focus" "$testbin" "$profile" > "${stem}-${phase}.txt" 2>&1 || true
  done <<'PHASES'
gen4|experimentalV4PhoneBuild(84|95)Generation4Fetch
continuation4|experimentalV4PhoneBuild(84|95)ContinueFetch
foldscore|experimentalV4PhoneBuild(84|95)FoldScoreFetch
reader|experimentalV4PhoneBuild(84|95)ReadProjectiveBlockValue
map-point|homography.mapPoint
quad-homography|experimentalV4PhoneQuadHomography
single-probe|experimentalV4PhoneBuild53SingleProbe
pair-stencil|experimentalV4PhoneBuild53PairStencil
PHASES
  {
    for fn in experimentalV4PhoneBuild84Generation4Fetch experimentalV4PhoneBuild95Generation4Fetch experimentalV4PhoneBuild84ContinueFetch experimentalV4PhoneBuild95ContinueFetch experimentalV4PhoneBuild84FoldScoreFetch experimentalV4PhoneBuild95FoldScoreFetch experimentalV4PhoneBuild84ReadProjectiveBlockValue experimentalV4PhoneBuild95ReadProjectiveBlockValue homography.mapPoint experimentalV4PhoneQuadHomography; do
      echo "===== $fn ====="
      env GOTOOLCHAIN="$GO_TOOLCHAIN" go tool pprof -list="$fn" "$testbin" "$profile" 2>&1 || true
      echo
    done
  } > "${stem}-list.txt"
done

cell(){ local cycle="$1" variant="$2" col="$3"; awk -F '\t' -v cy="$cycle" -v v="$variant" -v wanted="$col" 'NR==1{for(i=1;i<=NF;i++)if($i==wanted)c=i;next}$1==cy&&$4==v{print $c;exit}' "$tsv"; }
pct(){ awk -v a="$1" -v b="$2" 'BEGIN{if(a>0)printf "%.2f",100*(a-b)/a;else printf "-"}'; }
ratio(){ awk -v a="$1" -v b="$2" 'BEGIN{if(b>0)printf "%.3f",a/b;else printf "-"}'; }

{
  echo '# PixSeal Build96 Gen4 retained bottleneck profile'; echo
  echo 'Build96 is observability-only. It repeats B/angle in contemporaneous AB/BA order, but now captures one CPU profile and runtime memory/GC snapshot per deep run. Build93 remains the qualified production baseline; Build95 is retained only as an experimental contrast after its promotion was rejected.'; echo
  echo '| cycle | order | B93 total | B95 total | total reduction | B93 gen4 wall | B95 gen4 wall | gen4 reduction | B93 gen4 worker | B95 gen4 worker | worker reduction | B95 luma prep | B95 alloc delta bytes | B93 GC | B95 GC |'
  echo '|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|'
  for cycle in 1 2; do
    order="$(cell "$cycle" build93 order)"
    a="$(cell "$cycle" build93 total_ms)"; b="$(cell "$cycle" build95 total_ms)"
    g93="$(cell "$cycle" build93 gen4_wall_ms)"; g95="$(cell "$cycle" build95 gen4_wall_ms)"
    w93="$(cell "$cycle" build93 gen4_worker_ms)"; w95="$(cell "$cycle" build95 gen4_worker_ms)"
    echo "| $cycle | $order | $a | $b | $(pct "$a" "$b")% | $g93 | $g95 | $(pct "$g93" "$g95")% | $w93 | $w95 | $(pct "$w93" "$w95")% | $(cell "$cycle" build95 luminance_prep_ms) | $(cell "$cycle" build95 total_alloc_delta) | $(cell "$cycle" build93 num_gc_delta) | $(cell "$cycle" build95 num_gc_delta) |"
  done
  echo
  echo 'Interpretation rule: if Build95 cumulative CPU in reader/continuation falls materially while Gen4 wall remains nearly unchanged, the next investigation should target parallel memory/cache pressure rather than another reader micro-optimization. If cumulative CPU also fails to improve in the retained workload, the public microbenchmark is not representative of the real working set and the luminance-plane branch should stay closed.'
  echo
  echo 'Frozen B/angle semantic invariants for every row: geometry=334857, bank=6198, qualification=6198, qualified=0, FoldScore=118004, block reads=151045120, success=149370664, failed=1674456, outcome=REJECT.'
} > "$md"

rm -rf "$OUTPUT_DIR"
mv "$stage_dir" "$OUTPUT_DIR"
published=true

echo 'Build96 retained bottleneck profile written to:'
echo "  $OUTPUT_DIR/build96-profile.tsv"
echo "  $OUTPUT_DIR/build96-profile.md"
echo "  $OUTPUT_DIR/build96-profile-raw.txt"
echo "  $OUTPUT_DIR/build96-metadata.txt"
echo 'Build96 observability gate: PASS'
