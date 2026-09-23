#!/usr/bin/env bash

set -uo pipefail

REPORT="${ALL_TEST_REPORT:-}"
if [[ -n "$REPORT" ]]; then
    mkdir -p "$(dirname "$REPORT")"
    exec > >(tee "$REPORT") 2>&1
fi

default_targets="toolchain-check version-check vet release-unit v3-freeze-check v4-pilot-lock-check corpus-manifest-check test-images deep-test extreme-test research-unit lattice-estimator-test homography-test photometric-test bit-channel-test reliability-test spatial-channel-test phase-surface-test blind-phase-test lattice-phase-test global-unwrap-test crossfit-unwrap-test stability-unwrap-test cycle-anchor-test observability-audit-test physical-topology-test v4-design-study-test v4-foundation-test v4-pilot-search-test v4-pilot-channel-test v4-pilot-corpus-test v4-pilot-geometry-test v4-pilot-geometry-corpus-test v4-pilot-blind-geometry-test v4-pilot-blind-geometry-corpus-test v4-pilot-placement-test v4-pilot-placement-corpus-test v4-pilot-joint-affine-test v4-pilot-joint-affine-corpus-test v4-pilot-joint-projective-test v4-pilot-joint-projective-corpus-test v4-pilot-joint-projective-rank-diagnostic v4-build34-projective-frame-corpus-test v4-build35-projective-api-test v4-build36-soft-channel-test v4-build37-scanner-registration-test v4-build38-phone-channel-test v4-build39-phone-registration-test v4-build40-phone-residual-test v4-build41-phone-basin-test v4-build42-phone-data-test v4-build43-phone-side-pair-test v4-build44-jpeg-compat-test v4-build45-phone-diagnostic-test v4-build46-phone-handoff-test v4-build47-phone-frozen-bank-test v4-build48-phone-local-refine-test v4-build49-phone-proposal-ranking-test v4-pilot-lock-corpus-test v4-frame-test v4-frame-corpus-test smooth-phase-test geometry-test affine-test composition-test lattice-test perspective-test core-target-check"
read -r -a targets <<< "${ALL_TEST_TARGETS:-$default_targets}"
if (( ${#targets[@]} == 0 )); then
    echo "error: ALL_TEST_TARGETS resolved to an empty list" >&2
    exit 2
fi

label_for() {
    case "$1" in
        toolchain-check) echo "qualified Go toolchain" ;;
        version-check) echo "version consistency" ;;
        vet) echo "go vet" ;;
        release-unit) echo "release-gate Go tests" ;;
        v3-freeze-check) echo "Format-v3 frozen core manifest" ;;
        v4-pilot-lock-check) echo "Format-v4 Build30 pilot candidate identity lock" ;;
        corpus-manifest-check) echo "Build32 active private corpus manifest" ;;
        test-images) echo "local image round-trip" ;;
        deep-test) echo "baseline transformations" ;;
        extreme-test) echo "progressive limits" ;;
        research-unit) echo "experimental geometry Go regressions" ;;
        lattice-estimator-test) echo "v0.3 local lattice estimator" ;;
        homography-test) echo "v0.3 bounded homography/projective decoder" ;;
        photometric-test) echo "v0.3 bounded print-camera photometric bank" ;;
        bit-channel-test) echo "v0.3 protected-bit/ECC channel diagnostics" ;;
        reliability-test) echo "v0.3 bounded reliability/full-grid decoder" ;;
        spatial-channel-test) echo "v0.3 spatial protected-bit stability" ;;
        phase-surface-test) echo "v0.3 confidence-weighted phase surface" ;;
        blind-phase-test) echo "v0.3 repetition/cross-cell blind phase" ;;
        lattice-phase-test) echo "v0.3 local fractional lattice phase" ;;
        global-unwrap-test) echo "v0.3 global discrete phase unwrap" ;;
        crossfit-unwrap-test) echo "v0.3 disjoint repetition cross-fit" ;;
        stability-unwrap-test) echo "v0.3 multi-partition integer-cycle stability" ;;
        cycle-anchor-test) echo "v0.3 independent cross-cell cycle anchor" ;;
        observability-audit-test) echo "v0.3 Format-v3 observability audit" ;;
        physical-topology-test) echo "v0.3 physical differential topology observability" ;;
        v4-design-study-test) echo "Format-v4 absolute-pilot design study" ;;
        v4-foundation-test) echo "Format-v4 isolated pilot foundation" ;;
        v4-pilot-search-test) echo "Format-v4 Build24 pilot search/partial visibility" ;;
        v4-pilot-channel-test) echo "Format-v4 Build24 synthetic pilot image channel" ;;
        v4-pilot-corpus-test) echo "Format-v4 Build24 local image pilot corpus" ;;
        v4-pilot-geometry-test) echo "Format-v4 Build25 known-geometry pilot channel" ;;
        v4-pilot-geometry-corpus-test) echo "Format-v4 Build25 known-geometry local corpus" ;;
        v4-pilot-blind-geometry-test) echo "Format-v4 Build26 blind pilot-assisted geometry" ;;
        v4-pilot-blind-geometry-corpus-test) echo "Format-v4 Build26 blind geometry local corpus" ;;
        v4-pilot-placement-test) echo "Format-v4 Build27 unknown crop/placement" ;;
        v4-pilot-placement-corpus-test) echo "Format-v4 Build27 unknown-placement local corpus" ;;
        v4-pilot-joint-affine-test) echo "Format-v4 Build28 joint blind affine+crop" ;;
        v4-pilot-joint-affine-corpus-test) echo "Format-v4 Build28 joint blind affine+crop local corpus" ;;
        v4-pilot-joint-projective-test) echo "Format-v4 Build29 joint projective+crop / padded placement" ;;
        v4-pilot-joint-projective-corpus-test) echo "Format-v4 Build29 joint projective/padded local corpus" ;;
        v4-pilot-joint-projective-rank-diagnostic) echo "Format-v4 Build33 MQ projective ranking observability" ;;
        v4-build34-projective-frame-corpus-test) echo "Format-v4 Build34 authenticated MQ projective recovery" ;;
        v4-build35-projective-api-test) echo "Format-v4 Build35 projective API/CLI qualification" ;;
        v4-build36-soft-channel-test) echo "Format-v4 Build36 soft-decision physical channel" ;;
        v4-build37-scanner-registration-test) echo "Format-v4 Build37 blind scanner registration" ;;
        v4-build38-phone-channel-test) echo "Format-v4 Build38 smartphone-channel carrier" ;;
        v4-build39-phone-registration-test) echo "Format-v4 Build39 blind smartphone registration checkpoint" ;;
        v4-build40-phone-residual-test) echo "Format-v4 Build40 pilot-only smartphone residual warp" ;;
        v4-build41-phone-basin-test) echo "Format-v4 Build41 blind smartphone basin recovery" ;;
        v4-build42-phone-data-test) echo "Format-v4 Build42 qualified-bank/list data recovery" ;;
        v4-build43-phone-side-pair-test) echo "Format-v4 Build43 proposal-only side-pair geometry fallback" ;;
        v4-build44-jpeg-compat-test) echo "Build44 deterministic JPEG ingest" ;;
        v4-build45-phone-diagnostic-test) echo "Build45 phone failure-decomposition diagnostics" ;;
        v4-build46-phone-handoff-test) echo "Build46 qualified-geometry handoff diagnostics" ;;
        v4-build47-phone-frozen-bank-test) echo "Build47 frozen-candidate bank observability" ;;
        v4-build48-phone-local-refine-test) echo "Build48 proposal-only local projective refinement" ;;
        v4-build49-phone-proposal-ranking-test) echo "Build49 proposal-ranking observability" ;;
        v4-pilot-lock-corpus-test) echo "Format-v4 Build30 pilot lock known-mapping corpus audit" ;;
        v4-frame-test) echo "Format-v4 Build31 experimental frame/encoder" ;;
        v4-frame-corpus-test) echo "Format-v4 Build31 frame local corpus" ;;
        smooth-phase-test) echo "v0.3 bounded smooth phase field" ;;
        geometry-test) echo "rotation/combined geometry" ;;
        affine-test) echo "axis-aligned affine" ;;
        composition-test) echo "build-7 composition regression" ;;
        lattice-test) echo "direct lattice bank" ;;
        perspective-test) echo "mild projective perspective" ;;
        core-target-check) echo "desktop/mobile core portability" ;;
        *) echo "$1" ;;
    esac
}

release_targets=(toolchain-check version-check vet release-unit v3-freeze-check v4-pilot-lock-check corpus-manifest-check test-images v4-build44-jpeg-compat-test core-target-check)
qualification_targets=(deep-test extreme-test)
research_targets=(research-unit lattice-estimator-test homography-test photometric-test bit-channel-test reliability-test spatial-channel-test phase-surface-test blind-phase-test lattice-phase-test global-unwrap-test crossfit-unwrap-test stability-unwrap-test cycle-anchor-test observability-audit-test physical-topology-test v4-design-study-test v4-foundation-test v4-pilot-search-test v4-pilot-channel-test v4-pilot-corpus-test v4-pilot-geometry-test v4-pilot-geometry-corpus-test v4-pilot-blind-geometry-test v4-pilot-blind-geometry-corpus-test v4-pilot-placement-test v4-pilot-placement-corpus-test v4-pilot-joint-affine-test v4-pilot-joint-affine-corpus-test v4-pilot-joint-projective-test v4-pilot-joint-projective-corpus-test v4-pilot-joint-projective-rank-diagnostic v4-build34-projective-frame-corpus-test v4-build35-projective-api-test v4-build36-soft-channel-test v4-build37-scanner-registration-test v4-build38-phone-channel-test v4-build39-phone-registration-test v4-build40-phone-residual-test v4-build41-phone-basin-test v4-build42-phone-data-test v4-build43-phone-side-pair-test v4-build44-jpeg-compat-test v4-build45-phone-diagnostic-test v4-build46-phone-handoff-test v4-build47-phone-frozen-bank-test v4-build48-phone-local-refine-test v4-build49-phone-proposal-ranking-test v4-pilot-lock-corpus-test v4-frame-test v4-frame-corpus-test smooth-phase-test geometry-test affine-test composition-test lattice-test perspective-test)

in_list() {
    local needle="$1"; shift
    local item
    for item in "$@"; do
        [[ "$item" == "$needle" ]] && return 0
    done
    return 1
}

statuses=()
durations=()
started="$(date '+%Y-%m-%d %H:%M:%S %z')"
start_epoch="$(date +%s)"

printf 'PixSeal all-test\n'
if [[ -f VERSION ]]; then printf 'Version: %s\n' "$(head -n 1 VERSION)"; fi
printf 'Started: %s\n' "$started"
printf 'Host: %s\n' "$(uname -srm 2>/dev/null || echo unknown)"
printf 'Host Go: %s\n' "$(go version 2>/dev/null || echo unavailable)"
qualified_toolchain="${PIXSEAL_GO_TOOLCHAIN:-go1.26.0}"
printf 'Qualified Go: %s\n' "$(GOTOOLCHAIN="$qualified_toolchain" go version 2>/dev/null || echo unavailable)"
if command -v magick >/dev/null 2>&1; then
    printf 'ImageMagick: %s\n' "$(magick -version 2>/dev/null | head -n 1)"
elif command -v convert >/dev/null 2>&1; then
    printf 'ImageMagick: %s\n' "$(convert -version 2>/dev/null | head -n 1)"
fi
printf 'STRICT=%s is used for deep/geometry/affine/composition/lattice/perspective; extreme-test remains non-strict.\n' "${ALL_TEST_STRICT:-1}"
if [[ -n "$REPORT" ]]; then printf 'Report: %s\n' "$REPORT"; fi
printf '\n'

for i in "${!targets[@]}"; do
    target="${targets[$i]}"
    label="$(label_for "$target")"
    printf '%s\n' '========================================================================'
    printf '[%d/%d] %s (%s)\n' "$((i+1))" "${#targets[@]}" "$label" "$target"
    printf '%s\n' '------------------------------------------------------------------------'
    section_start="$(date +%s)"

    case "$target" in
        deep-test|geometry-test|affine-test|composition-test|lattice-test|perspective-test)
            make --no-print-directory STRICT="${ALL_TEST_STRICT:-1}" "$target"
            status=$?
            ;;
        *)
            make --no-print-directory "$target"
            status=$?
            ;;
    esac

    elapsed=$(( $(date +%s) - section_start ))
    durations[$i]="$elapsed"
    if (( status == 0 )); then statuses[$i]="PASS"; else statuses[$i]="FAIL($status)"; fi
    printf '\n[%s] %s — %ss\n\n' "${statuses[$i]}" "$target" "$elapsed"
done

total_elapsed=$(( $(date +%s) - start_epoch ))
printf '%s\n' '========================================================================'
printf 'ALL-TEST SUMMARY\n'
printf '%s\n' '========================================================================'
printf '%-22s %-12s %10s\n' 'Target' 'Result' 'Seconds'
printf '%-22s %-12s %10s\n' '----------------------' '------------' '----------'

failed=0
release_failed=0
research_failed=0
qualification_failed=0
release_selected=0
qualification_selected=0
research_selected=0
for i in "${!targets[@]}"; do
    target="${targets[$i]}"
    status="${statuses[$i]}"
    printf '%-22s %-12s %10s\n' "$target" "$status" "${durations[$i]}"
    if in_list "$target" "${release_targets[@]}"; then
        ((release_selected += 1))
        [[ "$status" == PASS ]] || ((release_failed += 1))
    elif in_list "$target" "${qualification_targets[@]}"; then
        ((qualification_selected += 1))
        [[ "$status" == PASS ]] || ((qualification_failed += 1))
    elif in_list "$target" "${research_targets[@]}"; then
        ((research_selected += 1))
        [[ "$status" == PASS ]] || ((research_failed += 1))
    fi
    [[ "$status" == PASS ]] || ((failed += 1))
done

printf '\nTotal elapsed: %ss\n' "$total_elapsed"

if (( release_selected == 0 )); then
    echo 'Release baseline: NOT RUN'
elif (( release_selected < ${#release_targets[@]} )); then
    if (( release_failed > 0 )); then
        printf 'Release baseline: PARTIAL (%d/%d targets run; %d failed)\n' "$release_selected" "${#release_targets[@]}" "$release_failed"
    else
        printf 'Release baseline: PARTIAL (%d/%d targets run)\n' "$release_selected" "${#release_targets[@]}"
    fi
elif (( release_failed == 0 )); then
    echo 'Release baseline: PASS'
else
    printf 'Release baseline: FAIL (%d release-gate target%s failed)\n' "$release_failed" "$([[ $release_failed -eq 1 ]] && echo '' || echo 's')"
fi

if (( qualification_selected == 0 )); then
    echo 'Qualification corpus: NOT RUN'
elif (( qualification_selected < ${#qualification_targets[@]} )); then
    if (( qualification_failed > 0 )); then
        printf 'Qualification corpus: PARTIAL/ATTENTION (%d/%d targets run; %d failed)\n' "$qualification_selected" "${#qualification_targets[@]}" "$qualification_failed"
    else
        printf 'Qualification corpus: PARTIAL (%d/%d targets run)\n' "$qualification_selected" "${#qualification_targets[@]}"
    fi
elif (( qualification_failed == 0 )); then
    echo 'Qualification corpus: PASS'
else
    printf 'Qualification corpus: ATTENTION (%d qualification target%s failed)\n' "$qualification_failed" "$([[ $qualification_failed -eq 1 ]] && echo '' || echo 's')"
fi

if (( research_selected == 0 )); then
    echo 'Research suites: NOT RUN'
elif (( research_selected < ${#research_targets[@]} )); then
    if (( research_failed > 0 )); then
        printf 'Research suites: PARTIAL/ATTENTION (%d/%d targets run; %d failed)\n' "$research_selected" "${#research_targets[@]}" "$research_failed"
    else
        printf 'Research suites: PARTIAL (%d/%d targets run)\n' "$research_selected" "${#research_targets[@]}"
    fi
elif (( research_failed == 0 )); then
    echo 'Research suites: PASS'
else
    printf 'Research suites: ATTENTION (%d experimental target%s failed)\n' "$research_failed" "$([[ $research_failed -eq 1 ]] && echo '' || echo 's')"
fi

if (( failed > 0 )); then
    printf 'Overall: FAIL (%d target%s failed)\n' "$failed" "$([[ $failed -eq 1 ]] && echo '' || echo 's')"
    exit 1
fi
if (( release_selected < ${#release_targets[@]} || qualification_selected < ${#qualification_targets[@]} || research_selected < ${#research_targets[@]} )); then
    echo 'Overall: PARTIAL (requested target set did not run the complete qualification matrix)'
    exit 0
fi
echo 'Overall: PASS'
