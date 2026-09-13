#!/usr/bin/env bash
set -euo pipefail

row() {
    printf '  %-24s %-14s %-18s %s\n' "$1" "$2" "$3" "$4"
}

section() {
    printf '\n%s\n' "$1"
    printf '  %-24s %-14s %-18s %s\n' 'TARGET' 'CLASS' 'INPUT' 'DESCRIPTION'
    printf '  %-24s %-14s %-18s %s\n' '------------------------' '--------------' '------------------' '-----------'
}

printf 'PixSeal available test/check targets\n'
printf 'Legend: INPUT=pics means the local qualification corpus; private means files are never shipped.\n'
printf 'Research/qualification targets may intentionally fail at current robustness boundaries; release-check is the production gate.\n'

section 'Release / baseline'
row version-check RELEASE self 'VERSION/buildinfo consistency'
row vet RELEASE source 'go vet on all packages'
row release-unit RELEASE source 'release-scoped Go unit/regression gate'
row v3-freeze-check RELEASE source 'verify frozen v3 core SHA-256 manifest'
row test-images RELEASE pics 'authenticated round-trip on each image/profile'
row deep-test QUALIFICATION pics 'strict baseline JPEG/resize/crop suite; corpus-sensitive'
row core-target-check PORTABILITY source 'compile reusable core for Linux/Windows/Android/iOS'
row release-check AGGREGATE source+pics 'complete production release baseline gate'

section 'Research / diagnostics'
row research-unit RESEARCH source 'experimental geometry Go regressions'
row lattice-estimator-test RESEARCH source 'local lattice estimator regressions'
row homography-test RESEARCH source 'bounded homography/projective regressions'
row photometric-test RESEARCH source 'bounded photometric bank regressions'
row bit-channel-test RESEARCH source 'protected-bit/ECC diagnostics'
row reliability-test RESEARCH source 'soft Hamming/full-grid regressions'
row spatial-channel-test RESEARCH source 'spatial coded-bit stability regressions'
row phase-surface-test RESEARCH source 'confidence-weighted phase surface regressions'
row blind-phase-test RESEARCH source 'repetition + cross-cell blind phase regressions'
row lattice-phase-test RESEARCH source 'local fractional lattice-phase regressions'
row global-unwrap-test RESEARCH source 'exact global +/-1 unwrap + exact top-2/split-repetition/rollback regressions'
row crossfit-unwrap-test RESEARCH source 'disjoint coded-bit-group proposal/validation cross-fit regressions'
row stability-unwrap-test RESEARCH source '8-partition/16-trial integer-cycle field stability regressions'
row cycle-anchor-test RESEARCH source 'unguided cross-cell image-domain cycle-anchor regressions'
row observability-audit-test RESEARCH source 'static key-independent Format-v3 cycle observability audit'
row physical-topology-test RESEARCH source 'held-out image-domain repetition-topology observability probe'
row v4-design-study-test RESEARCH source 'non-normative absolute-pilot Format-v4 sizing/design study'
row v4-foundation-test RESEARCH source 'isolated build23 provisional v4 pilot/data-plane invariants'
row v4-pilot-search-test RESEARCH source 'Build24 reproducible joint pilot search + partial-visibility qualification'
row v4-pilot-channel-test RESEARCH source 'Build24 pilot-only synthetic image-channel qualification'
row v4-pilot-corpus-test RESEARCH pics 'Build24 pilot-only qualification on local original images'
row v4-pilot-geometry-test RESEARCH source 'Build25 known-geometry rotation/affine/perspective pilot qualification'
row v4-pilot-geometry-corpus-test RESEARCH pics 'Build25 known-geometry pilot qualification on local originals'
row v4-pilot-blind-geometry-test RESEARCH source 'Build26 bounded blind geometry via repeat proposal + pilot validation'
row v4-pilot-blind-geometry-corpus-test RESEARCH pics 'Build26 blind pilot-assisted geometry on local originals'
row v4-pilot-placement-test RESEARCH source 'Build27 unknown crop/translation/placement with independently supplied geometry'
row v4-pilot-placement-corpus-test RESEARCH pics 'Build27 unknown-placement qualification on local originals'
row v4-pilot-joint-affine-test RESEARCH source 'Build28 joint blind rotation+anisotropic-scale plus unknown crop/translation'
row v4-pilot-joint-affine-corpus-test RESEARCH pics 'Build28 joint blind affine+crop qualification on local originals'
row smooth-phase-test RESEARCH source 'bounded smooth phase-field regressions'
row geometry-test RESEARCH pics 'strict rotation/combined geometry suite; corpus-sensitive'
row affine-test RESEARCH pics 'strict axis-aligned affine suite; corpus-sensitive'
row composition-test RESEARCH pics 'strict anisotropic-scale + rotation suite; corpus-sensitive'
row lattice-test RESEARCH pics 'strict direct lattice-basis composition suite; corpus-sensitive'
row perspective-test RESEARCH pics 'strict mild projective perspective suite; corpus-sensitive'

section 'Private physical-channel'
row print-camera-test PRIVATE private/print-camera 'print -> paper -> smartphone; PASS requires Format-v3 HMAC'
row print-scan-test PRIVATE private/print-scan 'print -> scanner; PASS requires Format-v3 HMAC'

section 'Exploratory / compatibility'
row extreme-test QUALIFICATION pics 'progressive resize/crop limit map; non-strict by design'
row test-unit COMPATIBILITY source 'complete go test ./... suite'
row test AGGREGATE source+pics 'build + test-unit + image round-trips'
row all AGGREGATE source+pics 'test + baseline transformation suite'
row all-test AGGREGATE source+pics '43-target qualification matrix with final summary'

printf '\nTip: use make <target>. For all-test, optionally set ALL_TEST_REPORT=path/to/report.txt.\n'
