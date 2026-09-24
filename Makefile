SHELL := /bin/bash
.DEFAULT_GOAL := build
.NOTPARALLEL:

# Build44 physically qualifies Go 1.26.0 with PixSeal-owned deterministic JPEG
# ingest. Use $(GO) for every normal Go command in this Makefile. Go 1.25.1
# remains a historical Build43 qualification reference only.
QUALIFIED_GO_TOOLCHAIN := go1.26.0
BUILD44_CANDIDATE_GO_TOOLCHAIN ?= go1.26.0
PIXSEAL_GO_TOOLCHAIN ?= $(QUALIFIED_GO_TOOLCHAIN)
EXPECTED_GO_TOOLCHAIN ?= $(PIXSEAL_GO_TOOLCHAIN)
GO := env GOTOOLCHAIN=$(PIXSEAL_GO_TOOLCHAIN) go

PIXSEAL := dist/pixseal
ORIGINAL_PICS_DIR := original pics
TEST_KEY ?= pixseal-test-key
ACTIVE_CORPUS_MANIFEST ?= private-corpus-active.tsv
TEST_PROFILES ?= robust balanced capacity
TEST_MESSAGE_ROBUST ?= PixSeal robust
TEST_MESSAGE_BALANCED ?= PixSeal balanced profile test
TEST_MESSAGE_CAPACITY ?= PixSeal capacity profile test message for regression coverage
ROBUST_RESIZES ?= 95 85 75 65 55
ROBUST_CROPS ?= 90 75 50
RANDOM_CROPS ?= 75 50
RANDOM_CROP_COUNT ?= 1
RANDOM_SEED ?= 20260907
JPEG_QUALITY ?= 82
ROBUST_MAX_MPIX ?= 100
EXTRACT_TIMEOUT ?= 60
TEST_IMAGES_MAX_MPIX ?= 100
GEOMETRY_EXTRACT_TIMEOUT ?= 120
STRICT ?= 0
LIMIT_START ?= 95
LIMIT_MIN ?= 10
LIMIT_STEP ?= 5
GEOMETRY_ANGLES ?= -45 -30 -15 -10 -5 -1 1 5 10 15 30 45 90 180 270
GEOMETRY_COMBINED_ANGLES ?= 12.3
GEOMETRY_COMBINED_MODES ?= rotate-resize75 resize75-rotate rotate-crop80 crop80-rotate rotate-resize75-crop80
GEOMETRY_MAX_MPIX ?= 50
AFFINE_MODES ?= scale110x90 scale90x110 shearx8 sheary8
AFFINE_MAX_MPIX ?= 50
COMPOSITION_PROFILES ?= robust
COMPOSITION_ANGLES ?= 12.3
COMPOSITION_MODES ?= scale110x90-rotate
COMPOSITION_MAX_MPIX ?= 50
LATTICE_ANGLES ?= 12.3
LATTICE_MODES ?= scale110x90-rotate scale90x110-rotate scale105x95-rotate scale95x105-rotate
LATTICE_MAX_MPIX ?= 50
PERSPECTIVE_MODES ?= top-narrow-4 bottom-narrow-4
PERSPECTIVE_MAX_MPIX ?= 50
PRINT_CAMERA_DIR ?= print-camera private
PRINT_CAMERA_KEY ?= Piccotti
PRINT_CAMERA_TIMEOUT ?= 180
PRINT_SCAN_DIR ?= print-scan private
V4_PILOT_CORPUS_DIR ?= $(ORIGINAL_PICS_DIR)
V4_PHYSICAL_SOURCE_DIR ?= $(ORIGINAL_PICS_DIR)
V4_PHYSICAL_OUTPUT_DIR ?= v4-physical private/build35-generated
V4_PHYSICAL_FIXTURE_DIR ?= $(V4_PHYSICAL_OUTPUT_DIR)
V4_PHYSICAL_ACQUISITION_DIR ?= v4-physical private/build35-acquired
V4_PHYSICAL_KEY ?= PixSeal-v4-TestKey-2026
V4_PHYSICAL_ROLE ?= MQ
V4_PHYSICAL_MESSAGE_A ?= v4-b35-phys-a
V4_PHYSICAL_MESSAGE_B ?= v4-b35-phys-b
V4_PHYSICAL_PROFILE ?= robust
V4_PHYSICAL_STRENGTH ?= 24
V4_PHYSICAL_PRINT_PPI ?= 300
V4_SCANNER_CONTROL ?=
V4_SCANNER_MARKED_A ?=
V4_SCANNER_MARKED_B ?=
V4_SCANNER_CANONICAL_WIDTH ?= 1632
V4_SCANNER_CANONICAL_HEIGHT ?= 1632
V4_PHONE_SOURCE_DIR ?= $(ORIGINAL_PICS_DIR)
V4_PHONE_OUTPUT_DIR ?= v4-phone private/build38-generated
V4_PHONE_KEY ?= PixSeal-v4-TestKey-2026
V4_PHONE_ROLE ?= MQ
V4_PHONE_MESSAGE_A ?= v4-b38-phone-a
V4_PHONE_MESSAGE_B ?= v4-b38-phone-b
V4_PHONE_PROFILE ?= robust
V4_PHONE_STRENGTH ?= 48
V4_PHONE_PRINT_PPI ?= 300
V4_PHONE_ACQUISITION_DIR ?= v4-phone private/build38-acquired
V4_PHONE_DIAGNOSTIC_DIR ?= v4-phone private/build40-diagnostics
V4_PHONE_BUILD41_DIAGNOSTIC_DIR ?= v4-phone private/build41-diagnostics
V4_PHONE_BUILD42_DIAGNOSTIC_DIR ?= v4-phone private/build42-diagnostics
V4_PHONE_BUILD43_DIAGNOSTIC_DIR ?= v4-phone private/build43-diagnostics
V4_PHONE_BUILD44_DIAGNOSTIC_DIR ?= v4-phone private/build44-diagnostics
V4_PHONE_BUILD45_DIAGNOSTIC_DIR ?= v4-phone private/build45-diagnostics
V4_PHONE_BUILD46_DIAGNOSTIC_DIR ?= v4-phone private/build46-diagnostics
V4_PHONE_BUILD47_DIAGNOSTIC_DIR ?= v4-phone private/build47-diagnostics
V4_PHONE_BUILD48_DIAGNOSTIC_DIR ?= v4-phone private/build48-diagnostics
V4_PHONE_BUILD49_DIAGNOSTIC_DIR ?= v4-phone private/build49-diagnostics
V4_PHONE_BUILD50_DIAGNOSTIC_DIR ?= v4-phone private/build50-diagnostics
V4_PHONE_BUILD51_DIAGNOSTIC_DIR ?= v4-phone private/build51-diagnostics
V4_PHONE_BUILD51_TIMEOUT ?= 2400
V4_PHONE_BUILD52_DIAGNOSTIC_DIR ?= v4-phone private/build52-diagnostics
V4_PHONE_BUILD52_TIMEOUT ?= 3600
V4_PHONE_BUILD53_DIAGNOSTIC_DIR ?= v4-phone private/build53-diagnostics
V4_PHONE_BUILD53_TIMEOUT ?= 3600
V4_PHONE_BUILD54_DIAGNOSTIC_DIR ?= v4-phone private/build54-diagnostics
V4_PHONE_BUILD54_TIMEOUT ?= 3600
V4_PHONE_BUILD55_DIAGNOSTIC_DIR ?= v4-phone private/build55-diagnostics
V4_PHONE_BUILD55_TIMEOUT ?= 5400
V4_PHONE_BUILD56_DIAGNOSTIC_DIR ?= v4-phone private/build56-diagnostics
V4_PHONE_BUILD56_TIMEOUT ?= 7200
V4_PHONE_BUILD57_DIAGNOSTIC_DIR ?= v4-phone private/build57-diagnostics
V4_PHONE_BUILD57_TIMEOUT ?= 10800
V4_PHONE_BUILD58_DIAGNOSTIC_DIR ?= v4-phone private/build58-diagnostics
V4_PHONE_BUILD58_TIMEOUT ?= 14400
V4_PHONE_BUILD59_DIAGNOSTIC_DIR ?= v4-phone private/build59-diagnostics
V4_PHONE_BUILD59_TIMEOUT ?= 18000
V4_PHONE_CANONICAL_WIDTH ?= 1632
V4_PHONE_CANONICAL_HEIGHT ?= 1632
V4_PHONE_TIMEOUT ?= 600
ALL_TEST_REPORT ?=
ALL_TEST_TARGETS ?=
ALL_TEST_STRICT ?=1
GO_SOURCES := $(shell find cmd internal watermark -type f -name '*.go')

.PHONY: toolchain-check test-list corpus-manifest-check private-corpus-manifest v4-physical-fixtures v4-physical-qualification v4-phone-fixtures v4-build40-phone-corpus-diagnostic v4-build41-phone-physical-test v4-build42-phone-physical-test print-scan-test build test test-unit release-unit v3-freeze-check v4-pilot-lock-check research-unit lattice-estimator-test homography-test photometric-test bit-channel-test reliability-test spatial-channel-test phase-surface-test blind-phase-test lattice-phase-test global-unwrap-test crossfit-unwrap-test stability-unwrap-test cycle-anchor-test observability-audit-test physical-topology-test v4-design-study-test v4-foundation-test v4-pilot-search-test v4-pilot-channel-test v4-pilot-corpus-test v4-pilot-geometry-test v4-pilot-geometry-corpus-test v4-pilot-blind-geometry-test v4-pilot-blind-geometry-corpus-test v4-pilot-placement-test v4-pilot-placement-corpus-test v4-pilot-joint-affine-test v4-pilot-joint-affine-corpus-test v4-pilot-joint-projective-test v4-pilot-joint-projective-corpus-test v4-pilot-joint-projective-rank-diagnostic v4-build34-projective-frame-corpus-test v4-build35-projective-api-test v4-build36-soft-channel-test v4-build37-scanner-registration-test v4-build38-phone-channel-test v4-build39-phone-registration-test v4-build40-phone-residual-test v4-build41-phone-basin-test v4-build42-phone-data-test v4-build43-phone-side-pair-test v4-build43-phone-physical-test v4-build44-jpeg-compat-test v4-build44-go126-jpeg-compat-test v4-build44-phone-physical-test v4-build44-go126-phone-physical-test v4-build45-phone-diagnostic-test v4-build45-phone-diagnostic v4-build45-phone-oracle-diagnostic v4-build45-phone-study v4-build46-phone-handoff-test v4-build46-phone-handoff-diagnostic v4-build47-phone-frozen-bank-test v4-build47-phone-frozen-bank-diagnostic v4-build48-phone-local-refine-test v4-build48-phone-local-refine-diagnostic v4-build49-phone-proposal-ranking-test v4-build49-phone-proposal-ranking-diagnostic v4-build50-phone-top4-refine-test v4-build50-phone-top4-refine-diagnostic v4-build51-phone-surface-test v4-build51-phone-surface-diagnostic v4-build52-phone-optimizer-test v4-build52-phone-optimizer-diagnostic v4-build53-phone-pair-escape-test v4-build53-phone-pair-escape-diagnostic v4-build54-phone-pair-continuation-test v4-build54-phone-pair-continuation-diagnostic v4-build55-phone-sibling-stencil-test v4-build55-phone-sibling-stencil-diagnostic v4-build56-phone-sibling-pair-escape-test v4-build56-phone-sibling-pair-escape-diagnostic v4-build57-phone-sibling-pair-continuation-test v4-build57-phone-sibling-pair-continuation-diagnostic v4-build58-phone-second-pair-sibling-stencil-test v4-build58-phone-second-pair-sibling-stencil-diagnostic v4-build59-phone-third-pair-escape-test v4-build59-phone-third-pair-escape-diagnostic v4-build37-physical-scanner-test v4-pilot-lock-corpus-test v4-frame-test v4-frame-corpus-test smooth-phase-test print-camera-test test-images deep-test extreme-test geometry-test affine-test composition-test lattice-test perspective-test all-test release-check version-check all build-all core-target-check vet clean

# Print a categorized index of all test/check targets without running them.
test-list:
	@bash ./scripts/test-list.sh


corpus-manifest-check:
	@echo "Checking Build32 active private corpus manifest..."
	@PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" bash ./scripts/check-active-corpus.sh

# Generate a local, Git-ignored SHA-256 manifest for the private physical corpus.
private-corpus-manifest:
	@PRINT_CAMERA_DIR="$(PRINT_CAMERA_DIR)" PRINT_SCAN_DIR="$(PRINT_SCAN_DIR)" bash ./scripts/private-corpus-manifest.sh

# Generate the private Build35 scanner-first v4 qualification pack. The MQ
# source is block-normalized before embedding, then emitted as one unmarked
# control plus two authenticated robust carriers with different payloads.
# Outputs are git-ignored and never included in source/evidence archives.
v4-physical-fixtures: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHYSICAL_SOURCE_DIR="$(V4_PHYSICAL_SOURCE_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	V4_PHYSICAL_OUTPUT_DIR="$(V4_PHYSICAL_OUTPUT_DIR)" \
	V4_PHYSICAL_KEY="$(V4_PHYSICAL_KEY)" \
	V4_PHYSICAL_ROLE="$(V4_PHYSICAL_ROLE)" \
	V4_PHYSICAL_MESSAGE_A="$(V4_PHYSICAL_MESSAGE_A)" \
	V4_PHYSICAL_MESSAGE_B="$(V4_PHYSICAL_MESSAGE_B)" \
	V4_PHYSICAL_PROFILE="$(V4_PHYSICAL_PROFILE)" \
	V4_PHYSICAL_STRENGTH="$(V4_PHYSICAL_STRENGTH)" \
	V4_PHYSICAL_PRINT_PPI="$(V4_PHYSICAL_PRINT_PPI)" \
	bash ./scripts/prepare-v4-physical-fixtures.sh

# Evaluate the three Build35 physical captures listed by the generated
# acquisition plan. Marked images must authenticate the exact expected payload;
# the unmarked control must reject. This target is intentionally opt-in.
v4-physical-qualification: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHYSICAL_FIXTURE_DIR="$(V4_PHYSICAL_FIXTURE_DIR)" \
	V4_PHYSICAL_ACQUISITION_DIR="$(V4_PHYSICAL_ACQUISITION_DIR)" \
	V4_PHYSICAL_KEY="$(V4_PHYSICAL_KEY)" \
	bash ./scripts/test-v4-physical-acquisitions.sh

# Generate the private Build38 smartphone qualification pack at strength 48.
# The Format-v4 wire format, pilot, Hamming code and HMAC remain unchanged.
v4-phone-fixtures: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_SOURCE_DIR="$(V4_PHONE_SOURCE_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_ROLE="$(V4_PHONE_ROLE)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	V4_PHONE_PROFILE="$(V4_PHONE_PROFILE)" \
	V4_PHONE_STRENGTH="$(V4_PHONE_STRENGTH)" \
	V4_PHONE_PRINT_PPI="$(V4_PHONE_PRINT_PPI)" \
	bash ./scripts/prepare-v4-phone-fixtures.sh

# Opt-in Build40 matrix over the nine private strength-48 smartphone captures.
# This target does not change geometry or decoding behavior; it only records
# stage-by-stage telemetry needed to decide the first Build41 algorithm change.
v4-build40-phone-corpus-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_DIAGNOSTIC_DIR="$(V4_PHONE_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/diagnose-v4-phone-corpus.sh

# Verify the selected Build44 qualified toolchain (Go 1.26.0).
toolchain-check:
	@out="$$( $(GO) version 2>&1 )" || { \
		echo "error: unable to run qualified Go $(QUALIFIED_GO_TOOLCHAIN) toolchain" >&2; \
		echo "$$out" >&2; \
		echo "Install Go $(QUALIFIED_GO_TOOLCHAIN) or allow Go's GOTOOLCHAIN mechanism to download it." >&2; \
		exit 1; \
	}; \
	actual="$$(printf '%s\n' "$$out" | awk '{print $$3}')"; \
	if [[ "$$actual" != "$(EXPECTED_GO_TOOLCHAIN)" ]]; then \
		echo "error: PixSeal expected Go $(EXPECTED_GO_TOOLCHAIN); selected $$actual" >&2; \
		exit 1; \
	fi; \
	echo "Selected Go toolchain: $$actual"

# Default target: always rebuild the native executable with the qualified
# Build44 toolchain. Rebuilding avoids silently reusing a binary from another
# Go release.
build: toolchain-check
	@echo "Building PixSeal with $(PIXSEAL_GO_TOOLCHAIN)..."
	@mkdir -p dist
	@$(GO) build -trimpath -ldflags="-s -w" -o $(PIXSEAL) ./cmd/pixseal
	@echo "Created $(PIXSEAL)"

# Local release test suite: unit tests plus round trips on original pics.
test: build test-unit test-images
	@echo "Local test suite completed."

test-unit:
	@echo "Running complete Go unit/regression suite..."
	@$(GO) test ./...

# Release-gate Go tests deliberately exclude the expensive experimental geometry
# regression group. make test still runs every Go test for compatibility.
release-unit:
	@echo "Running release-gate Go tests..."
	@$(GO) test ./cmd/pixseal ./internal/buildinfo ./internal/jpeglegacy -count=1
	@$(GO) test ./watermark -run 'Test(V3EncoderGoldenFingerprint|StrengthRejectsNonFiniteValues|AnalyzerUsesSameWhiteAlphaFlatteningAsEncoder|WorkingImageLimitRejectsBeforePixelPlaneAllocation|WorkingImageLimitRejectsIntegerOverflow|IsotropicScaleSearchIsFixed|HammingCorrectsSingleBit|ProfileSelectionThresholds|ExplicitProfileCapacityErrors|V3ProfileRoundTrips|V3TransformsByProfile|V3AutoProfileExtraction|WrongKeyAndUnmarkedImageAreBounded|AnalyzeImageMatchesProfileMath|V3FrameIgnoresTrailingPaddingButAuthenticatesHeader|V3SyncPatternObservationCounts|V3TileMappingObservationCounts)$$' -count=1

v3-freeze-check:
	@echo "Checking frozen Format-v3 core hashes..."
	@bash ./scripts/check-v3-frozen-core.sh

# Build30 semantic lock for the exact prototype-2 public pilot identity. This
# protects the candidate hash/mask/sign sequence from accidental mutation while
# normative Format-v4 promotion still waits for physical print-camera evidence.
v4-pilot-lock-check:
	@echo "Checking Build30 Format-v4 pilot candidate lock..."
	@$(GO) test ./watermark -run '^TestExperimentalV4PilotCandidateLock(Identity|DetectsMutation)$$' -count=1 -v


# v0.3 bounded local-lattice diagnostic regressions. These are research tests and
# deliberately remain outside the v0.2-derived release-unit gate.
lattice-estimator-test:
	@echo "Running v0.3 local-lattice estimator regressions..."
	@$(GO) test ./watermark -run '^TestDiagnostic' -count=1

# v0.3 deterministic print-boundary/homography/projective virtual-decoder
# regressions. These stay research-only and do not alter ExtractWithInfo.
homography-test:
	@echo "Running v0.3 bounded homography/projective regressions..."
	@$(GO) test ./watermark -run 'Test(PrintBoundaryEstimatorFindsSyntheticPrint|HomographyForPrintBoundaryMapsCorners|DiagnosticScaleClusteringRewardsCrossRegionSupport|DiagnosticProjectiveVirtualAuthentication|DiagnosticPhaseConsensusPrefersExactScale|DiagnosticPhaseDifferenceIsBoundedModuloTile|DiagnosticFundamentalSelectionRejectsHigherFrequencyAliases|DiagnosticSubpixelOffsetCorrectsBoundaryPhase|DiagnosticResidualWarpFitsSmoothSubBlockField|DiagnosticAdaptiveEscalationAddsOneFinerLevel|DiagnosticWeakBoundaryRequiresPlausibleQuad)$$' -count=1

photometric-test:
	@echo "Running v0.3 bounded print-camera photometric regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticPhotometric' -count=1

bit-channel-test:
	@echo "Running v0.3 protected-bit/ECC diagnostics..."
	@$(GO) test ./watermark -run '^TestDiagnosticBitChannel' -count=1

reliability-test:
	@echo "Running v0.3 bounded reliability/full-grid regressions..."
	@$(GO) test ./watermark -run '^TestDiagnostic(SoftHamming|Reliability|FullGridSampling)' -count=1

spatial-channel-test:
	@echo "Running v0.3 spatial protected-bit stability regressions..."
	@$(GO) test ./watermark -run '^TestDiagnostic(SpatialBitEvidence|FullGridSamplingCollectsSpatialCells)' -count=1

phase-surface-test:
	@echo "Running v0.3 confidence-weighted phase-surface regressions..."
	@$(GO) test ./watermark -run '^TestDiagnostic(PhaseSurface|SmoothPhaseConfidence|SmoothPhaseQuadratic|SmoothPhaseHuber)' -count=1

blind-phase-test:
	@echo "Running v0.3 key-independent blind phase regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticBlind' -count=1

lattice-phase-test:
	@echo "Running v0.3 local fractional lattice-phase regressions..."
	@$(GO) test ./watermark -run '^TestDiagnostic(LocalLatticeFractionalPhase|BlindLatticeFusion)' -count=1

global-unwrap-test:
	@echo "Running v0.3 global discrete phase-unwrapping regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticGlobalDiscreteUnwrap' -count=1

crossfit-unwrap-test:
	@echo "Running v0.3 held-out repetition cross-fit regressions..."
	@$(GO) test ./watermark -run '^TestDiagnostic(Crossfit|GlobalUnwrapCrossfit|ApplyCrossfit)' -count=1

stability-unwrap-test:
	@echo "Running v0.3 multi-partition integer-cycle stability regressions..."
	@$(GO) test ./watermark -run '^TestDiagnostic(StabilityPartitions|GlobalUnwrapPartitionStability|ApplyStability)' -count=1

cycle-anchor-test:
	@echo "Running v0.3 independent cross-cell cycle-anchor regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticCycleAnchor' -count=1

observability-audit-test:
	@echo "Running v0.3 Format-v3 key-independent observability audit regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticFormatObservability' -count=1

physical-topology-test:
	@echo "Running v0.3 held-out physical repetition-topology observability regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticPhysicalTopology' -count=1

v4-design-study-test:
	@echo "Running non-normative Format-v4 absolute-pilot design-study regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticFormatV4DesignStudy' -count=1

v4-foundation-test:
	@echo "Running isolated experimental Format-v4 pilot-foundation regressions..."
	@$(GO) test ./watermark -run 'TestExperimentalV4Prototype(FoundationInvariants|UsesOnePilotPerStratum)$$' -count=1

# Build24 deterministic joint coordinate/sign search plus partial-visibility qualification.
v4-pilot-search-test:
	@echo "Running Build24 deterministic Format-v4 pilot search/qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4(Prototype2StructuralQualification|Prototype2ImprovesBuild23Baseline|Build24SearchReproducesPrototype2)$$' -count=1

# Build24 image-domain pilot-only synthetic channel. This does not encode or
# authenticate a payload and remains disconnected from production v3 APIs.
v4-pilot-channel-test:
	@echo "Running Build24 experimental Format-v4 pilot image-channel regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Prototype2ImageDomainPilotChannel$$' -count=1

# Optional local image-corpus pilot qualification. Source archives intentionally
# omit the corpus; the target uses ORIGINAL_PICS_DIR by default when present.
v4-pilot-corpus-test:
	@echo "Running Build24 experimental Format-v4 pilot corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4PilotCorpus$$' -count=1 -v

# Build25 known-geometry pilot qualification. Geometry is supplied independently
# so this isolates whether prototype-2 survives the transformed image channel.
v4-pilot-geometry-test:
	@echo "Running Build25 Format-v4 known-geometry pilot qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Prototype2KnownGeometryQualification$$' -count=1 -v

# Optional Build25 known-geometry qualification on the local original-image corpus.
v4-pilot-geometry-corpus-test:
	@echo "Running Build25 Format-v4 known-geometry local corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4PilotGeometryCorpus$$' -count=1 -v

# Build26 bounded blind geometry recovery. Repeated data-plane self-consistency
# proposes geometry without data symbols; the public pilot ranks and validates.
v4-pilot-blind-geometry-test:
	@echo "Running Build26 Format-v4 blind pilot-assisted geometry qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4(Prototype2BlindGeometryQualification|DataRepeatProposalIsPilotSymbolIndependent)$$' -count=1 -v

# Optional Build26 blind-geometry qualification on the local original-image corpus.
v4-pilot-blind-geometry-corpus-test:
	@echo "Running Build26 Format-v4 blind geometry local corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4BlindGeometryCorpus$$' -count=1 -v

# Build27 unknown crop/translation/placement qualification with geometry supplied
# independently. Pilot half A proposes placement; disjoint pilot half B validates.
v4-pilot-placement-test:
	@echo "Running Build27 Format-v4 unknown-placement qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4(Prototype2UnknownPlacementQualification|PlacementWrongGeometryDoesNotMasqueradeAsPlacement)$$' -count=1 -v

# Optional Build27 placement qualification on the local original-image corpus.
v4-pilot-placement-corpus-test:
	@echo "Running Build27 Format-v4 unknown-placement local corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4PlacementCorpus$$' -count=1 -v

# Build28 joint affine + unknown negative-crop qualification. Geometry is
# selected only from public, sign-independent DCT phase contrast; the pilot is
# exposed only after geometry is fixed, through the qualified Build27 placement
# proposal/held-out validation split.
v4-pilot-joint-affine-test:
	@echo "Running Build28 Format-v4 joint blind affine+crop qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4(Prototype2JointAffineCropQualification|JointAffineStructuralProposalIsPilotSymbolIndependent)$$' -count=1 -v

# Optional Build28 joint affine+crop qualification on the local originals.
v4-pilot-joint-affine-corpus-test:
	@echo "Running Build28 Format-v4 joint blind affine+crop local corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4JointAffineCropCorpus$$' -count=1 -v

# Build29 composes unknown projective geometry with crop and unknown affine
# geometry with positive padded placement. Public-pilot evidence is used only
# for geometry/placement proposal and absolute-origin validation; payload, key,
# ECC and HMAC are never consulted.
v4-pilot-joint-projective-test:
	@echo "Running Build29 Format-v4 joint projective+crop / padded qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build29JointQualification$$' -count=1 -v

# Optional Build29 local-corpus gate. PASS includes explicitly documented SAFE
# REJECT outcomes for cases that the bounded joint search does not qualify.
v4-pilot-joint-projective-corpus-test:
	@echo "Running Build29 Format-v4 joint projective/padded local corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4Build29JointCorpus$$' -count=1 -v

# Build33 ranking-observability checkpoint for the MQ blocker. This diagnostic
# compares the same projective+crop transform with a synthetic random data plane
# and a real authenticated v4 frame, then reports the truth-basin rank after
# structural, half-pilot and full-proposal stages. It does not claim decode success.
v4-pilot-joint-projective-rank-diagnostic:
	@echo "Running Build33 Format-v4 MQ projective ranking diagnostic..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4Build33ProjectiveRankingDiagnostic$$' -count=1 -v

# Build34 end-to-end MQ projective+crop gate. Geometry is selected only from
# structure/public-pilot evidence; authenticated frame recovery happens after
# the unchanged projective acceptance gate.
v4-build34-projective-frame-corpus-test:
	@echo "Running Build34 Format-v4 authenticated MQ projective corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4Build34AuthenticatedProjectiveCorpus$$' -count=1 -v

# Build35 exposes the qualified Build34 blind projective decoder as a public
# experimental API/CLI for controlled physical-channel qualification.
v4-build35-projective-api-test:
	@echo "Running Build35 Format-v4 projective API/CLI qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build35ProjectiveAPI' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^(TestSubcommandHelpReturnsFlagErrHelp|TestV4ExtractProjectiveRequiresCanonicalBlockDimensions)$$' -count=1 -v

# Build36 promotes reliability-aware Hamming decoding into the accepted
# projective v4 data path. Geometry/pilot acceptance remains unchanged.
v4-build36-soft-channel-test:
	@echo "Running Build36 Format-v4 soft-decision physical-channel qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build36SoftHamming' -count=1 -v

# Build37 closes blind registration for full-page scanner captures. Boundary
# geometry proposes a narrow affine basin; disjoint public-pilot halves qualify
# a five-geometry ensemble before any data/HMAC work occurs.
v4-build37-scanner-registration-test:
	@echo "Running Build37 Format-v4 blind scanner-registration qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build37' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^(TestSubcommandHelpReturnsFlagErrHelp|TestV4ExtractScannerRequiresCanonicalBlockDimensions)$$' -count=1 -v

# Build38 freezes the smartphone qualification carrier at strength 48 while
# preserving the Format-v4 pilot, data mapping, Hamming code and HMAC.
v4-build38-phone-channel-test:
	@echo "Running Build38 Format-v4 smartphone-channel carrier qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build38PhoneStrengthRoundTrip$$' -count=1 -v

# Build39 introduces the first dedicated blind-smartphone registration checkpoint.
# The gate qualifies bounded downsampling, phone artwork-boundary extraction and
# spatially disjoint public-pilot projective evidence; it deliberately does not
# claim complete blind HMAC closure on the private real-phone corpus yet.
v4-build39-phone-registration-test:
	@echo "Running Build39 Format-v4 blind smartphone-registration checkpoint..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build39' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^(TestSubcommandHelpReturnsFlagErrHelp|TestV4ExtractPhoneRequiresCanonicalBlockDimensions)$$' -count=1 -v

# Build40 tests the hypothesis that the remaining phone error is a small smooth
# local warp after Build39. Proposal controls come only from checkerboard-A
# public-pilot tiles; checkerboard-B tiles validate the fitted field.
v4-build40-phone-residual-test:
	@echo "Running Build40 Format-v4 pilot-only smartphone residual-warp qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build40PilotOnlyResidualWarp$$' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^(TestSubcommandHelpReturnsFlagErrHelp|TestV4ExtractPhoneRequiresCanonicalBlockDimensions)$$' -count=1 -v

# Build41 improves the blind global phone basin while keeping the Build40
# residual layer, carrier and authentication channel unchanged. A fixed 1/3
# spatial fold is held out until the proposal-only shortlist is frozen.
v4-build41-phone-basin-test:
	@echo "Running Build41 Format-v4 blind smartphone basin qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build41PhoneBasinRecovery$$' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^(TestSubcommandHelpReturnsFlagErrHelp|TestV4ExtractPhoneRequiresCanonicalBlockDimensions)$$' -count=1 -v

# Build42 preserves Build41 geometry and adds only post-geometry recovery: the
# complete already-qualified bank feeds deterministic 3-way data ensembles and
# a bounded soft-Hamming list decoder. HMAC remains final authentication only.
v4-build42-phone-data-test:
	@echo "Running Build42 Format-v4 qualified-bank/list-decoder regression..."
	@$(GO) test ./watermark -run '^(TestExperimentalV4Build41PhoneBasinRecovery|TestExperimentalV4Build42ListDecodeRecoversSecondBestWord)$$' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^(TestSubcommandHelpReturnsFlagErrHelp|TestV4ExtractPhoneRequiresCanonicalBlockDimensions|TestV4PhoneDiagnosticsExposeBuild42MatrixFields)$$' -count=1 -v


# Build43 adds a bounded side-pair geometry fallback only after Build41 geometry
# rejects. Candidate generation/ranking is proposal-only and frozen before held-out.
v4-build43-phone-side-pair-test:
	@echo "Running Build43 Format-v4 proposal-only side-pair regression..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build43PhoneSidePairRecovery$$' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^(TestSubcommandHelpReturnsFlagErrHelp|TestV4ExtractPhoneRequiresCanonicalBlockDimensions)$$' -count=1 -v

# Build44 freezes JPEG rasterization inside PixSeal instead of inheriting the
# host toolchain's image/jpeg implementation. This baseline regression runs on
# the currently qualified Go toolchain and protects both the vendored decoder
# vectors and the CLI ingest path.
v4-build44-jpeg-compat-test:
	@echo "Running Build44 deterministic JPEG-ingest regression..."
	@$(GO) test ./internal/jpeglegacy -run '^TestDeterministicLegacyJPEGFixture$$' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^TestOpenImageUsesDeterministicLegacyJPEGDecoder$$' -count=1 -v

# Explicit Go 1.26 deterministic-ingest check. Kept as a named regression for
# the qualification evidence that promoted Go 1.26.0 in Build44.
v4-build44-go126-jpeg-compat-test:
	@set -euo pipefail; \
	out="$$(env GOTOOLCHAIN=$(BUILD44_CANDIDATE_GO_TOOLCHAIN) go version 2>&1)" || { echo "$$out" >&2; exit 1; }; \
	echo "Build44 qualified toolchain check: $$out"; \
	env GOTOOLCHAIN=$(BUILD44_CANDIDATE_GO_TOOLCHAIN) go test ./internal/jpeglegacy -run '^TestDeterministicLegacyJPEGFixture$$' -count=1 -v; \
	env GOTOOLCHAIN=$(BUILD44_CANDIDATE_GO_TOOLCHAIN) go test ./cmd/pixseal -run '^TestOpenImageUsesDeterministicLegacyJPEGDecoder$$' -count=1 -v

# Opt-in private physical regression over the original full-page scanner files.
# Paths are explicit so the private captures never enter the source archive.
v4-build37-physical-scanner-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_SCANNER_CONTROL="$(V4_SCANNER_CONTROL)" \
	V4_SCANNER_MARKED_A="$(V4_SCANNER_MARKED_A)" \
	V4_SCANNER_MARKED_B="$(V4_SCANNER_MARKED_B)" \
	V4_SCANNER_CANONICAL_WIDTH="$(V4_SCANNER_CANONICAL_WIDTH)" \
	V4_SCANNER_CANONICAL_HEIGHT="$(V4_SCANNER_CANONICAL_HEIGHT)" \
	V4_PHYSICAL_KEY="$(V4_PHYSICAL_KEY)" \
	V4_PHYSICAL_MESSAGE_A="$(V4_PHYSICAL_MESSAGE_A)" \
	V4_PHYSICAL_MESSAGE_B="$(V4_PHYSICAL_MESSAGE_B)" \
	bash ./scripts/test-v4-build37-scanner.sh

# Build30 freeze-readiness evidence with independently supplied mappings. The
# same locked pilot is measured directly on projective-crop and padded fixtures
# over the private photographic corpus, including Build29 SAFE-REJECT cases.
v4-pilot-lock-corpus-test:
	@echo "Running Build30 Format-v4 pilot candidate-lock corpus audit..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4PilotCandidateLockCorpus$$' -count=1 -v

# Build31 first real experimental v4 framing/data path. This target covers
# frame/version separation, Hamming baseline, locked pilot/data partition,
# aligned image round-trip, JPEG/crop channel checks and the explicit CLI.
v4-frame-test:
	@echo "Running Build31 experimental Format-v4 framing/encoder qualification..."
	@$(GO) test ./watermark -run '^TestExperimentalV4(FrameRoundTripProfiles|AlignedCropRoundTrip|WrongKeyAndCrossVersionIsolation|HammingBaselineCorrectsSingleBitPerCodeword|FrameDeterministicVector|CapacityGeometry|LockedDataPartition|AlignedJPEGChannel|MinimumTileRoundTrip|DataMappingProfileCoverage|HeaderBytes)$$' -count=1 -v
	@$(GO) test ./cmd/pixseal -run '^TestExperimentalV4CLIRoundTrip$$' -count=1 -v

# Opt-in Build31 real-image framing/channel qualification. The private originals
# are never distributed; the same experimental frame is tested in all three
# profiles plus robust JPEG-q82 and block-aligned crop.
v4-frame-corpus-test:
	@echo "Running Build31 experimental Format-v4 frame corpus qualification..."
	@PIXSEAL_V4_CORPUS_DIR="$(CURDIR)/$(V4_PILOT_CORPUS_DIR)" \
	PIXSEAL_V4_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
		$(GO) test ./watermark -run '^TestExperimentalV4FrameCorpus$$' -count=1 -v

smooth-phase-test:
	@echo "Running v0.3 bounded smooth phase-field regressions..."
	@$(GO) test ./watermark -run '^TestDiagnosticSmoothPhase' -count=1

# Private real print -> paper -> smartphone regression gate. Every PNG/JPEG in
# the private corpus directory is tested automatically. The photographs are
# never distributed with the source tree. Missing/empty corpus is a clean SKIP;
# when present, each image PASS requires an authenticated Format v3 HMAC.
print-camera-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PRINT_CAMERA_DIR="$(PRINT_CAMERA_DIR)" \
	PRINT_CAMERA_KEY="$(PRINT_CAMERA_KEY)" \
	PRINT_CAMERA_TIMEOUT="$(PRINT_CAMERA_TIMEOUT)" \
	bash ./scripts/test-print-camera.sh


# Build45 source-level diagnostic regression. This does not touch the private
# corpus and does not change the Build44 production decoder.
v4-build45-phone-diagnostic-test:
	@echo "Running Build45 phone failure-classification regressions..."
	@$(GO) test ./watermark -run 'TestExperimentalV4Build45Classification$$' -count=1
	@$(GO) test ./cmd/pixseal -run 'TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build45 blind decomposition over B/mild (primary) and B/angle (informational).
# The command calls the unchanged production decoder and records where it stops.
v4-build45-phone-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD45_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD45_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/diagnose-v4-build45-phone.sh

# Optional laboratory oracle. OpenCV/SIFT and the known digital marked-B carrier
# are used only to supply geometry; they are never linked into PixSeal or used by
# production search/ranking.
v4-build45-phone-oracle-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD45_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD45_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build45-phone-oracle.sh

# Convenience target: blind decomposition first, then the isolated lab oracle.
v4-build45-phone-study: v4-build45-phone-diagnostic v4-build45-phone-oracle-diagnostic

# Build46 source-level diagnostic regression. Production thresholds are constants
# under observation; this target does not change the phone decoder.
v4-build46-phone-handoff-test:
	@echo "Running Build46 qualified-geometry handoff regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build46HandoffClassification$$' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build46 private diagnostic. Blind Build41/43 search completes first. Only after
# that does the script optionally create SIFT/reference geometry for post-hoc
# corner-error comparison; the reference never guides blind search or ranking.
v4-build46-phone-handoff-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD46_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD46_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/diagnose-v4-build46-phone-handoff.sh

# Build47 diagnostic-only regression. Production Build43 remains capped at 32
# frozen candidates; the research surface observes deterministic prefixes 32/64/128.
v4-build47-phone-frozen-bank-test:
	@echo "Running Build47 frozen-candidate bank regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build47' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build47 private observability study. Both blind super-banks are written before
# SIFT/reference geometry is generated for post-hoc distance measurement.
v4-build47-phone-frozen-bank-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD47_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD47_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/diagnose-v4-build47-phone-frozen-bank.sh

# Build48 diagnostic-only regression. Seed selection must depend on proposal
# evidence only; production Build43/42 paths and quorum remain unchanged.
v4-build48-phone-local-refine-test:
	@echo "Running Build48 proposal-only local-refinement regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build48' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build48 private study. Both blind refinements are completed before the script
# creates SIFT/reference geometry for post-hoc before/after error measurement.
v4-build48-phone-local-refine-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD48_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD48_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/diagnose-v4-build48-phone-local-refine.sh

# Build49 diagnostic-only regression. All candidate ranks use proposal folds
# only; production Build43 ranking and decoder paths remain unchanged.
v4-build49-phone-proposal-ranking-test:
	@echo "Running Build49 proposal-ranking observability regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build49' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build49 private study. Blind proposal observables for both images are frozen
# before the SIFT/reference oracle is generated for post-hoc rank analysis.
v4-build49-phone-proposal-ranking-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD49_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD49_DIAGNOSTIC_DIR)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/diagnose-v4-build49-phone-proposal-ranking.sh

# Build50 diagnostic-only regression. The historical Build48 top-2 selector
# remains frozen; Build50 extends only the research seed depth to top 4/pair.
v4-build50-phone-top4-refine-test:
	@echo "Running Build50 top-4-per-pair local-refinement regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build50' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build50 private study. Both blind top-4 refined banks are fixed before the
# SIFT/reference oracle is generated. top2/top4 are nested views of one run.
v4-build50-phone-top4-refine-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD50_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD50_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/diagnose-v4-build50-phone-top4-refine.sh

# Build51 research-only regression. The exact Build41 proposal coordinate-descent
# objective is traced; the local stencil is deterministic and oracle-free.
v4-build51-phone-surface-test:
	@echo "Running Build51 local proposal-surface/trajectory regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build51' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build51 private study. Both blind traces/stencils are fully frozen before the
# SIFT/reference oracle is generated for post-hoc surface analysis.
v4-build51-phone-surface-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD51_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD51_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD51_TIMEOUT="$(V4_PHONE_BUILD51_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build51-phone-surface.sh

# Build52 research-only regression. The top-4 seed selection remains unchanged;
# the new fine restart uses only the existing proposal fold and retains every
# accepted 2px -> 1px state before held-out qualification.
v4-build52-phone-optimizer-test:
	@echo "Running Build52 fine-restart optimizer regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build52' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build52 private study. Both blind optimizer banks are completely fixed before
# the SIFT/reference oracle is generated for post-hoc geometry analysis.
v4-build52-phone-optimizer-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD52_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD52_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD52_TIMEOUT="$(V4_PHONE_BUILD52_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build52-phone-optimizer.sh

# Build53 research-only regression. The unchanged top-4 seed bank and proposal
# score are preserved. Coupled +/-1px pair scans are enabled only at proposal-
# defined 1px coordinate-local roots retained from the 2px path.
v4-build53-phone-pair-escape-test:
	@echo "Running Build53 coupled pair-escape regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build53' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build53 private study. Both blind pair-escape banks are completely frozen
# before SIFT/reference oracle geometry is generated post-hoc.
v4-build53-phone-pair-escape-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD53_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD53_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD53_TIMEOUT="$(V4_PHONE_BUILD53_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build53-phone-pair-escape.sh

# Build54 research-only regression. Build53 pair generation remains unchanged;
# every retained pair state is continued with bounded proposal-only 1px
# coordinate descent and every accepted intermediate is retained.
v4-build54-phone-pair-continuation-test:
	@echo "Running Build54 post-pair continuation regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build54' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build54 private study. Both blind continuation banks are completely frozen
# before SIFT/reference oracle geometry is generated post-hoc.
v4-build54-phone-pair-continuation-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD54_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD54_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD54_TIMEOUT="$(V4_PHONE_BUILD54_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build54-phone-pair-continuation.sh

# Build55 research-only regression. Build54 continuation remains unchanged;
# every retained continuation state receives an independent full +/-1px sibling
# stencil evaluated against the same frozen parent.
v4-build55-phone-sibling-stencil-test:
	@echo "Running Build55 continuation sibling-stencil regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build55' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build55 private study. Both blind sibling banks are fully frozen before
# SIFT/reference oracle geometry is generated post-hoc.
v4-build55-phone-sibling-stencil-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD55_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD55_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD55_TIMEOUT="$(V4_PHONE_BUILD55_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build55-phone-sibling-stencil.sh

# Build56 research-only regression. The complete Build55 sibling bank is
# unchanged; only proposal-local siblings receive a bounded coupled +/-1px
# two-coordinate escape scan.
v4-build56-phone-sibling-pair-escape-test:
	@echo "Running Build56 sibling pair-escape regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build56' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build56 private study. Both blind expanded banks are fully frozen before
# SIFT/reference oracle geometry is generated post-hoc.
v4-build56-phone-sibling-pair-escape-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD56_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD56_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD56_TIMEOUT="$(V4_PHONE_BUILD56_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build56-phone-sibling-pair-escape.sh

# Build57 research-only regression. The complete Build56 second-pair bank is
# unchanged; every retained second-pair state receives bounded proposal-only
# 1px continuation and every accepted intermediate is retained.
v4-build57-phone-sibling-pair-continuation-test:
	@echo "Running Build57 second-pair continuation regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build57' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build57 private study. Both blind expanded banks are fully frozen before
# SIFT/reference oracle geometry is generated post-hoc.
v4-build57-phone-sibling-pair-continuation-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD57_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD57_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD57_TIMEOUT="$(V4_PHONE_BUILD57_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build57-phone-sibling-pair-continuation.sh

# Build58 research-only regression. Reproduce the complete Build57 bank and
# evaluate all independent +/-1px siblings from every retained post-second-pair
# continuation state against the same frozen parent.
v4-build58-phone-second-pair-sibling-stencil-test:
	@echo "Running Build58 second-pair continuation sibling-stencil regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build58' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build58 private study. Both blind expanded banks are fully frozen before
# SIFT/reference oracle geometry is generated post-hoc.
v4-build58-phone-second-pair-sibling-stencil-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD58_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD58_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD58_TIMEOUT="$(V4_PHONE_BUILD58_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build58-phone-second-pair-sibling-stencil.sh

# Build59 research-only regression. Reproduce the complete Build58 bank, probe
# every frozen second-pair sibling for one-coordinate locality, and apply the
# bounded pair stencil only at proposal-local sibling states.
v4-build59-phone-third-pair-escape-test:
	@echo "Running Build59 third-pair escape regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build59' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build59 private study. Both image banks are frozen before oracle generation.
v4-build59-phone-third-pair-escape-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD59_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD59_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD59_TIMEOUT="$(V4_PHONE_BUILD59_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build59-phone-third-pair-escape.sh

# Optional private print -> scanner corpus. Uses the same HMAC-only semantics as
# print-camera-test but keeps acquisition methods separate in the filesystem.

# Opt-in Build41 physical milestone over the private nine-photo strength-48
# corpus. All controls must reject and at least one A plus one B capture must
# authenticate the exact expected payload.
v4-build41-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD41_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD41_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/test-v4-build41-phone-corpus.sh

# Opt-in Build42 physical gate. Controls must still reject before data decode;
# A/angle and B/front must retain their Build41 direct passes; A/mild must now
# authenticate through the Build42 qualified-bank/list-decoder fallback.
v4-build42-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD42_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD42_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/test-v4-build42-phone-corpus.sh


# Opt-in Build43 physical gate over the original Build38 nine-photo corpus.
v4-build43-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD43_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD43_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/test-v4-build43-phone-corpus.sh

# Build44 qualified Go 1.26 physical gate. The watermark/geometry/data path is
# unchanged from Build43; only JPEG ingest is frozen inside PixSeal.
v4-build44-go126-phone-physical-test:
	@set -euo pipefail; \
	out="$$(env GOTOOLCHAIN=$(BUILD44_CANDIDATE_GO_TOOLCHAIN) go version 2>&1)" || { echo "$$out" >&2; exit 1; }; \
	actual="$$(printf '%s\n' "$$out" | awk '{print $$3}')"; \
	if [[ "$$actual" != "$(BUILD44_CANDIDATE_GO_TOOLCHAIN)" ]]; then echo "error: Build44 qualification requires $(BUILD44_CANDIDATE_GO_TOOLCHAIN); selected $$actual" >&2; exit 1; fi; \
	echo "Build44 qualified toolchain: $$actual"; \
	mkdir -p dist; \
	env GOTOOLCHAIN=$(BUILD44_CANDIDATE_GO_TOOLCHAIN) go build -trimpath -ldflags="-s -w" -o dist/pixseal-build44-go126 ./cmd/pixseal; \
	PIXSEAL="$(abspath dist/pixseal-build44-go126)" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD44_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD44_DIAGNOSTIC_DIR)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	V4_PHONE_TIMEOUT="$(V4_PHONE_TIMEOUT)" \
	bash ./scripts/test-v4-build44-phone-corpus.sh

# Preferred Build44 physical qualification command. Go 1.26.0 is the qualified
# toolchain; the historical go126-named target remains as an explicit alias.
v4-build44-phone-physical-test: v4-build44-go126-phone-physical-test

print-scan-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PRINT_CAMERA_DIR="$(PRINT_SCAN_DIR)" \
	PRINT_CAMERA_KEY="$(PRINT_CAMERA_KEY)" \
	PRINT_CAMERA_TIMEOUT="$(PRINT_CAMERA_TIMEOUT)" \
	PRINT_CAMERA_LABEL="Print-scan" \
	bash ./scripts/test-print-camera.sh

# Deterministic Go regressions for experimental geometry. all-test runs this
# separately so research failures cannot make the release baseline red.
research-unit:
	@echo "Running experimental geometry Go regressions..."
	@$(GO) test ./watermark -run 'Test(V3QuarterTurnRecovery|V3ArbitraryRotationRecovery|V3CombinedGeometryRecovery|V3DirectLatticeBasisRecovery|DirectLatticeBasisSearchIsFixed|RotationProbeRejectsUnmarkedSyntheticImage|WrongKeyOnRotatedCarrierIsBounded|V3AxisAlignedAffineScaleRecovery|AxisAlignedAffineSearchIsFixed|V3MildPerspectiveRecoveryEndToEnd|Build11PerspectiveHypothesisBound)$$' -count=1

test-images: build
	@set -euo pipefail; \
	if [[ ! -d "$(ORIGINAL_PICS_DIR)" ]]; then \
		echo "error: image directory not found: $(ORIGINAL_PICS_DIR)" >&2; \
		exit 1; \
	fi; \
	source ./scripts/test-common.sh; \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))"; export ACTIVE_CORPUS_MANIFEST; \
	load_active_corpus_images "$(ORIGINAL_PICS_DIR)"; \
	tmp_dir="$$(mktemp -d)"; \
	trap 'rm -rf -- "$$tmp_dir"' EXIT; \
	read -r -a profiles <<< "$(TEST_PROFILES)"; \
	if (( $${#profiles[@]} == 0 )); then \
		echo "error: TEST_PROFILES is empty" >&2; \
		exit 1; \
	fi; \
	if command -v magick >/dev/null 2>&1; then identify_tool=(magick identify); else identify_tool=(identify); fi; \
	tested_images=0; skipped_images=0; \
	for image in "$${images[@]}"; do \
		name="$$(basename "$$image")"; \
		read -r width height < <(read_image_dimensions "$$image") || { echo "error: cannot read dimensions for $$image" >&2; exit 1; }; \
		if (( $(TEST_IMAGES_MAX_MPIX) > 0 && width * height > $(TEST_IMAGES_MAX_MPIX) * 1000000 )); then \
			echo "Testing $$image"; echo "  SKIP all profiles $$(( (width*height+999999)/1000000 )) MP exceeds limit of $(TEST_IMAGES_MAX_MPIX) MP"; skipped_images=$$((skipped_images+1)); continue; \
		fi; \
		tested_images=$$((tested_images+1)); \
		echo "Testing $$image"; \
		for profile in "$${profiles[@]}"; do \
			case "$$profile" in \
				robust) message='$(TEST_MESSAGE_ROBUST)' ;; \
				balanced) message='$(TEST_MESSAGE_BALANCED)' ;; \
				capacity) message='$(TEST_MESSAGE_CAPACITY)' ;; \
				*) echo "error: unsupported TEST_PROFILES entry: $$profile" >&2; exit 1 ;; \
			esac; \
			marked="$$tmp_dir/$$name-$$profile.png"; \
			"$(PIXSEAL)" embed -in "$$image" -out "$$marked" \
				-key "$(TEST_KEY)" -message "$$message" -profile "$$profile" >/dev/null; \
			extracted="$$($(PIXSEAL) extract -raw -in "$$marked" -key "$(TEST_KEY)" 2>/dev/null)"; \
			if [[ "$$extracted" != "$$message" ]]; then \
				echo "error: extracted message differs for $$image ($$profile)" >&2; \
				exit 1; \
			fi; \
			echo "  OK $$profile"; \
		done; \
	done; \
	if (( tested_images == 0 )); then echo "error: no active corpus image was within test-images budget" >&2; exit 1; fi; \
	echo "Image round-trip tests completed: $$tested_images images tested x $${#profiles[@]} profiles; $$skipped_images image(s) skipped"

deep-test: build
	@echo "Running image transformation tests..."
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	TEST_KEY="$(TEST_KEY)" \
	TEST_PROFILES="$(TEST_PROFILES)" \
	TEST_MESSAGE_ROBUST="$(TEST_MESSAGE_ROBUST)" \
	TEST_MESSAGE_BALANCED="$(TEST_MESSAGE_BALANCED)" \
	TEST_MESSAGE_CAPACITY="$(TEST_MESSAGE_CAPACITY)" \
	ROBUST_RESIZES="$(ROBUST_RESIZES)" \
	ROBUST_CROPS="$(ROBUST_CROPS)" \
	RANDOM_CROPS="$(RANDOM_CROPS)" \
	RANDOM_CROP_COUNT="$(RANDOM_CROP_COUNT)" \
	RANDOM_SEED="$(RANDOM_SEED)" \
	JPEG_QUALITY="$(JPEG_QUALITY)" \
	ROBUST_MAX_MPIX="$(ROBUST_MAX_MPIX)" \
	EXTRACT_TIMEOUT="$(EXTRACT_TIMEOUT)" \
	STRICT="$(STRICT)" \
	bash ./scripts/test-robustness.sh

# Explore resize and crop limits; intentionally excluded from make all.
extreme-test: build
	@echo "Running progressive limit tests..."
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	TEST_KEY="$(TEST_KEY)" \
	TEST_PROFILES="$(TEST_PROFILES)" \
	TEST_MESSAGE_ROBUST="$(TEST_MESSAGE_ROBUST)" \
	TEST_MESSAGE_BALANCED="$(TEST_MESSAGE_BALANCED)" \
	TEST_MESSAGE_CAPACITY="$(TEST_MESSAGE_CAPACITY)" \
	ROBUST_MAX_MPIX="$(ROBUST_MAX_MPIX)" \
	EXTRACT_TIMEOUT="$(EXTRACT_TIMEOUT)" \
	LIMIT_START="$(LIMIT_START)" \
	LIMIT_MIN="$(LIMIT_MIN)" \
	LIMIT_STEP="$(LIMIT_STEP)" \
	RANDOM_SEED="$(RANDOM_SEED)" \
	bash ./scripts/test-limits.sh

# Experimental geometric robustness suite; intentionally excluded from make all.
geometry-test: build
	@echo "Running digital rotation and combined geometry tests..."
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	TEST_KEY="$(TEST_KEY)" \
	TEST_PROFILES="$(TEST_PROFILES)" \
	TEST_MESSAGE_ROBUST="$(TEST_MESSAGE_ROBUST)" \
	TEST_MESSAGE_BALANCED="$(TEST_MESSAGE_BALANCED)" \
	TEST_MESSAGE_CAPACITY="$(TEST_MESSAGE_CAPACITY)" \
	GEOMETRY_ANGLES="$(GEOMETRY_ANGLES)" \
	GEOMETRY_COMBINED_ANGLES="$(GEOMETRY_COMBINED_ANGLES)" \
	GEOMETRY_COMBINED_MODES="$(GEOMETRY_COMBINED_MODES)" \
	GEOMETRY_MAX_MPIX="$(GEOMETRY_MAX_MPIX)" \
	EXTRACT_TIMEOUT="$(GEOMETRY_EXTRACT_TIMEOUT)" \
	STRICT="$(STRICT)" \
	bash ./scripts/test-geometry.sh

# Experimental axis-aligned affine suite; intentionally excluded from make all.
affine-test: build
	@echo "Running axis-aligned affine tests..."
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	TEST_KEY="$(TEST_KEY)" \
	TEST_PROFILES="$(TEST_PROFILES)" \
	TEST_MESSAGE_ROBUST="$(TEST_MESSAGE_ROBUST)" \
	TEST_MESSAGE_BALANCED="$(TEST_MESSAGE_BALANCED)" \
	TEST_MESSAGE_CAPACITY="$(TEST_MESSAGE_CAPACITY)" \
	AFFINE_MODES="$(AFFINE_MODES)" \
	AFFINE_MAX_MPIX="$(AFFINE_MAX_MPIX)" \
	EXTRACT_TIMEOUT="$(EXTRACT_TIMEOUT)" \
	STRICT="$(STRICT)" \
	bash ./scripts/test-affine.sh

# Experimental composed geometry suite; intentionally excluded from make all.
composition-test: build
	@echo "Running composed anisotropic-scale + rotation tests..."
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	TEST_KEY="$(TEST_KEY)" \
	TEST_MESSAGE_ROBUST="$(TEST_MESSAGE_ROBUST)" \
	COMPOSITION_PROFILES="$(COMPOSITION_PROFILES)" \
	COMPOSITION_ANGLES="$(COMPOSITION_ANGLES)" \
	COMPOSITION_MODES="$(COMPOSITION_MODES)" \
	COMPOSITION_MAX_MPIX="$(COMPOSITION_MAX_MPIX)" \
	EXTRACT_TIMEOUT="$(EXTRACT_TIMEOUT)" \
	STRICT="$(STRICT)" \
	bash ./scripts/test-composition.sh

# Experimental direct lattice-basis suite; intentionally excluded from make all.
lattice-test: build
	@echo "Running direct lattice-basis composition tests..."
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	TEST_KEY="$(TEST_KEY)" \
	TEST_MESSAGE_ROBUST="$(TEST_MESSAGE_ROBUST)" \
	LATTICE_ANGLES="$(LATTICE_ANGLES)" \
	LATTICE_MODES="$(LATTICE_MODES)" \
	LATTICE_MAX_MPIX="$(LATTICE_MAX_MPIX)" \
	EXTRACT_TIMEOUT="$(EXTRACT_TIMEOUT)" \
	STRICT="$(STRICT)" \
	bash ./scripts/test-lattice.sh

# Experimental mild projective/keystone suite; intentionally excluded from make all.
perspective-test: build
	@echo "Running mild perspective tests..."
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PICS_DIR="$(CURDIR)/$(ORIGINAL_PICS_DIR)" \
	ACTIVE_CORPUS_MANIFEST="$(abspath $(ACTIVE_CORPUS_MANIFEST))" \
	TEST_KEY="$(TEST_KEY)" TEST_MESSAGE_ROBUST="$(TEST_MESSAGE_ROBUST)" \
	PERSPECTIVE_MODES="$(PERSPECTIVE_MODES)" PERSPECTIVE_MAX_MPIX="$(PERSPECTIVE_MAX_MPIX)" \
	EXTRACT_TIMEOUT="$(EXTRACT_TIMEOUT)" STRICT="$(STRICT)" bash ./scripts/test-perspective.sh

# Run every test/check target sequentially, continue after individual failures,
# and print one comparable summary at the end. STRICT=1 is applied to the
# baseline-transform and geometric research suites that support it; extreme-test
# intentionally remains a non-strict progressive limit map.
# Set ALL_TEST_REPORT=path/to/report.txt to tee the complete run to a file.
all-test:
	@ALL_TEST_REPORT="$(ALL_TEST_REPORT)" \
	ALL_TEST_TARGETS="$(ALL_TEST_TARGETS)" \
	ALL_TEST_STRICT="$(ALL_TEST_STRICT)" \
	PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	bash ./scripts/test-all.sh

version-check:
	@echo "Checking VERSION/buildinfo consistency..."
	@$(GO) test ./internal/buildinfo -run TestVersionFileMatchesBuildInfo -count=1

vet:
	@echo "Running go vet..."
	@$(GO) vet ./...

# Release-candidate gate: static analysis, release-scoped unit tests, local
# round trips, baseline transforms and reusable-core portability. Research Go
# regressions remain visible through research-unit / all-test. Requires original pics/.
# Strict mode is target-specific so a plain `make release-check` is self-contained.
release-check: STRICT := 1
release-check: toolchain-check version-check vet release-unit v3-freeze-check v4-pilot-lock-check corpus-manifest-check test-images v4-build44-jpeg-compat-test deep-test core-target-check
	@echo "Release baseline checks passed."

# Run build + local round-trip tests + baseline transformation tests.
all: test deep-test
	@echo "Complete local test suite finished."

build-all:
	@echo "Building PixSeal for Linux, Windows and macOS..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o dist/pixseal-linux-amd64 ./cmd/pixseal
	@GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o dist/pixseal-windows-amd64.exe ./cmd/pixseal
	@GOOS=darwin GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o dist/pixseal-macos-amd64 ./cmd/pixseal
	@GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="-s -w" -o dist/pixseal-macos-arm64 ./cmd/pixseal
	@echo "Created cross-platform binaries in dist/"

# Compile the reusable steganography core for representative future frontend targets.
# This creates no distributable binaries; it is an architectural portability check.
core-target-check:
	@echo "Checking reusable core on desktop/mobile targets..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build ./watermark
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build ./watermark
	@CGO_ENABLED=0 GOOS=android GOARCH=arm64 $(GO) build ./watermark
	@CGO_ENABLED=0 GOOS=ios GOARCH=arm64 $(GO) build ./watermark
	@echo "Core target checks passed: linux/amd64, windows/amd64, android/arm64, ios/arm64"

clean:
	@rm -rf -- dist
	@echo "Removed dist/"
