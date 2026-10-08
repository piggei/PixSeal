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
V4_PHONE_BUILD60_DIAGNOSTIC_DIR ?= v4-phone private/build60-diagnostics
V4_PHONE_BUILD60_TIMEOUT ?= 28800
V4_PHONE_BUILD61_DIAGNOSTIC_DIR ?= v4-phone private/build61-diagnostics
V4_PHONE_BUILD61_TIMEOUT ?= 43200
V4_PHONE_BUILD62_DIAGNOSTIC_DIR ?= v4-phone private/build62-diagnostics
V4_PHONE_BUILD62_TIMEOUT ?= 64800
V4_PHONE_BUILD63_DIAGNOSTIC_DIR ?= v4-phone private/build63-diagnostics
V4_PHONE_BUILD63_TIMEOUT ?= 86400
V4_PHONE_BUILD64_DIAGNOSTIC_DIR ?= v4-phone private/build64-diagnostics
V4_PHONE_BUILD64_TIMEOUT ?= 86400
V4_PHONE_BUILD65_DIAGNOSTIC_DIR ?= v4-phone private/build65-diagnostics
V4_PHONE_BUILD65_TIMEOUT ?= 86400
V4_PHONE_BUILD66_DIAGNOSTIC_DIR ?= v4-phone private/build66-diagnostics
V4_PHONE_BUILD66_TIMEOUT ?= 86400
V4_PHONE_BUILD65_BASELINE_TSV ?= v4-phone private/build65-diagnostics/build65-phone-matrix.tsv
V4_PHONE_BUILD67_DIAGNOSTIC_DIR ?= v4-phone private/build67-diagnostics
V4_PHONE_BUILD67_TIMEOUT ?= 86400
V4_PHONE_BUILD66_BASELINE_TSV ?= v4-phone private/build66-diagnostics/build66-phone-matrix.tsv
V4_PHONE_BUILD68_DIAGNOSTIC_DIR ?= v4-phone private/build68-diagnostics
V4_PHONE_BUILD68_TIMEOUT ?= 86400
V4_PHONE_BUILD67_BASELINE_TSV ?= v4-phone private/build67-diagnostics/build67-phone-profile.tsv
V4_PHONE_BUILD69_DIAGNOSTIC_DIR ?= v4-phone private/build69-diagnostics
V4_PHONE_BUILD69_TIMEOUT ?= 86400
V4_PHONE_BUILD70_DIAGNOSTIC_DIR ?= v4-phone private/build70-diagnostics
V4_PHONE_BUILD70_TIMEOUT ?= 86400
V4_PHONE_BUILD71_DIAGNOSTIC_DIR ?= v4-phone private/build71-diagnostics
V4_PHONE_BUILD71_TIMEOUT ?= 86400
V4_PHONE_BUILD72_DIAGNOSTIC_DIR ?= v4-phone private/build72-diagnostics
V4_PHONE_BUILD72_TIMEOUT ?= 86400
V4_PHONE_BUILD73_DIAGNOSTIC_DIR ?= v4-phone private/build73-diagnostics
V4_PHONE_BUILD73_TIMEOUT ?= 86400
V4_PHONE_BUILD74_DIAGNOSTIC_DIR ?= v4-phone private/build74-diagnostics
V4_PHONE_BUILD74_TIMEOUT ?= 86400
V4_PHONE_BUILD75_DIAGNOSTIC_DIR ?= v4-phone private/build75-diagnostics
V4_PHONE_BUILD75_TIMEOUT ?= 86400
V4_PHONE_BUILD76_DIAGNOSTIC_DIR ?= v4-phone private/build76-diagnostics
V4_PHONE_BUILD76_TIMEOUT ?= 86400
V4_PHONE_BUILD77_DIAGNOSTIC_DIR ?= v4-phone private/build77-diagnostics
V4_PHONE_BUILD77_TIMEOUT ?= 86400
V4_PHONE_BUILD78_DIAGNOSTIC_DIR ?= v4-phone private/build78-diagnostics
V4_PHONE_BUILD78_TIMEOUT ?= 86400
V4_PHONE_BUILD79_DIAGNOSTIC_DIR ?= v4-phone private/build79-diagnostics
V4_PHONE_BUILD79_TIMEOUT ?= 86400
V4_PHONE_BUILD80_DIAGNOSTIC_DIR ?= v4-phone private/build80-diagnostics
V4_PHONE_BUILD80_TIMEOUT ?= 86400
V4_PHONE_BUILD81_DIAGNOSTIC_DIR ?= v4-phone private/build81-diagnostics
V4_PHONE_BUILD81_TIMEOUT ?= 86400
V4_PHONE_BUILD82_DIAGNOSTIC_DIR ?= v4-phone private/build82-diagnostics
V4_PHONE_BUILD82_TIMEOUT ?= 86400
V4_BUILD83_PROFILE_DIR ?= v4-phone private/build83-profile
V4_BUILD83_BENCHTIME ?= 2s
V4_BUILD83_COUNT ?= 5
V4_BUILD83_PPROF_TIME ?= 20s
V4_PHONE_BUILD84_DIAGNOSTIC_DIR ?= v4-phone private/build84-diagnostics
V4_PHONE_BUILD84_TIMEOUT ?= 86400
V4_BUILD84_BENCH_DIR ?= v4-phone private/build84-benchmark
V4_BUILD84_BENCHTIME ?= 2s
V4_BUILD84_COUNT ?= 5
V4_BUILD85_PROFILE_DIR ?= v4-phone private/build85-profile
V4_BUILD85_BENCHTIME ?= 2s
V4_BUILD85_COUNT ?= 5
V4_BUILD85_PPROF_TIME ?= 20s
V4_BUILD86_BENCH_DIR ?= v4-phone private/build86-benchmark
V4_BUILD86_BENCHTIME ?= 2s
V4_BUILD86_COUNT ?= 5
V4_PHONE_BUILD86_DIAGNOSTIC_DIR ?= v4-phone private/build86-diagnostics
V4_PHONE_BUILD86_TIMEOUT ?= 86400
V4_BUILD87_PROFILE_DIR ?= v4-phone private/build87-profile
V4_BUILD87_BENCHTIME ?= 2s
V4_BUILD87_COUNT ?= 5
V4_BUILD87_PPROF_TIME ?= 20s
V4_BUILD88_BENCH_DIR ?= v4-phone private/build88-benchmark
V4_BUILD88_BENCHTIME ?= 2s
V4_BUILD88_COUNT ?= 5
V4_PHONE_BUILD88_DIAGNOSTIC_DIR ?= v4-phone private/build88-diagnostics
V4_PHONE_BUILD88_TIMEOUT ?= 86400
V4_PHONE_BUILD89_PROFILE_DIR ?= v4-phone private/build89-profile
V4_PHONE_BUILD89_TIMEOUT ?= 86400
V4_BUILD90_BENCH_DIR ?= v4-phone private/build90-benchmark
V4_BUILD90_BENCHTIME ?= 2s
V4_BUILD90_COUNT ?= 5
V4_PHONE_BUILD90_DIAGNOSTIC_DIR ?= v4-phone private/build90-diagnostics
V4_PHONE_BUILD90_TIMEOUT ?= 86400
V4_PHONE_BUILD75_BASELINE_TSV ?= v4-phone private/build75-diagnostics/build75-phone-basin-parallel.tsv
V4_PHONE_BUILD76_BASELINE_TSV ?= docs/qualified-baselines/build76-phone-performance.tsv
V4_PHONE_BUILD84_BASELINE_TSV ?= docs/qualified-baselines/build84-phone-performance.tsv
V4_PHONE_BUILD73_BASELINE_TSV ?= v4-phone private/build73-diagnostics/build73-phone-gen3-parallel.tsv
V4_PHONE_BUILD74_BASELINE_TSV ?= v4-phone private/build74-diagnostics/build74-phone-freeze-profile.tsv
V4_PHONE_BUILD71_BASELINE_TSV ?= v4-phone private/build71-diagnostics/build71-phone-gen4-parallel.tsv
V4_PHONE_BUILD68_BASELINE_TSV ?= v4-phone private/build68-diagnostics/build68-phone-plane-reuse.tsv
V4_PHONE_CANONICAL_WIDTH ?= 1632
V4_PHONE_CANONICAL_HEIGHT ?= 1632
V4_PHONE_TIMEOUT ?= 600
ALL_TEST_REPORT ?=
ALL_TEST_TARGETS ?=
ALL_TEST_STRICT ?=1
GO_SOURCES := $(shell find cmd internal watermark -type f -name '*.go')

.PHONY: toolchain-check test-list corpus-manifest-check private-corpus-manifest v4-physical-fixtures v4-physical-qualification v4-phone-fixtures v4-build40-phone-corpus-diagnostic v4-build41-phone-physical-test v4-build42-phone-physical-test print-scan-test build test test-unit release-unit v3-freeze-check v4-pilot-lock-check research-unit lattice-estimator-test homography-test photometric-test bit-channel-test reliability-test spatial-channel-test phase-surface-test blind-phase-test lattice-phase-test global-unwrap-test crossfit-unwrap-test stability-unwrap-test cycle-anchor-test observability-audit-test physical-topology-test v4-design-study-test v4-foundation-test v4-pilot-search-test v4-pilot-channel-test v4-pilot-corpus-test v4-pilot-geometry-test v4-pilot-geometry-corpus-test v4-pilot-blind-geometry-test v4-pilot-blind-geometry-corpus-test v4-pilot-placement-test v4-pilot-placement-corpus-test v4-pilot-joint-affine-test v4-pilot-joint-affine-corpus-test v4-pilot-joint-projective-test v4-pilot-joint-projective-corpus-test v4-pilot-joint-projective-rank-diagnostic v4-build34-projective-frame-corpus-test v4-build35-projective-api-test v4-build36-soft-channel-test v4-build37-scanner-registration-test v4-build38-phone-channel-test v4-build39-phone-registration-test v4-build40-phone-residual-test v4-build41-phone-basin-test v4-build42-phone-data-test v4-build43-phone-side-pair-test v4-build43-phone-physical-test v4-build44-jpeg-compat-test v4-build44-go126-jpeg-compat-test v4-build44-phone-physical-test v4-build44-go126-phone-physical-test v4-build45-phone-diagnostic-test v4-build45-phone-diagnostic v4-build45-phone-oracle-diagnostic v4-build45-phone-study v4-build46-phone-handoff-test v4-build46-phone-handoff-diagnostic v4-build47-phone-frozen-bank-test v4-build47-phone-frozen-bank-diagnostic v4-build48-phone-local-refine-test v4-build48-phone-local-refine-diagnostic v4-build49-phone-proposal-ranking-test v4-build49-phone-proposal-ranking-diagnostic v4-build50-phone-top4-refine-test v4-build50-phone-top4-refine-diagnostic v4-build51-phone-surface-test v4-build51-phone-surface-diagnostic v4-build52-phone-optimizer-test v4-build52-phone-optimizer-diagnostic v4-build53-phone-pair-escape-test v4-build53-phone-pair-escape-diagnostic v4-build54-phone-pair-continuation-test v4-build54-phone-pair-continuation-diagnostic v4-build55-phone-sibling-stencil-test v4-build55-phone-sibling-stencil-diagnostic v4-build56-phone-sibling-pair-escape-test v4-build56-phone-sibling-pair-escape-diagnostic v4-build57-phone-sibling-pair-continuation-test v4-build57-phone-sibling-pair-continuation-diagnostic v4-build58-phone-second-pair-sibling-stencil-test v4-build58-phone-second-pair-sibling-stencil-diagnostic v4-build59-phone-third-pair-escape-test v4-build59-phone-third-pair-escape-diagnostic v4-build60-phone-third-pair-continuation-test v4-build60-phone-third-pair-continuation-diagnostic v4-build61-phone-third-pair-sibling-stencil-test v4-build61-phone-third-pair-sibling-stencil-diagnostic v4-build62-phone-fourth-pair-escape-test v4-build62-phone-fourth-pair-escape-diagnostic v4-build63-phone-fourth-pair-continuation-test v4-build63-phone-fourth-pair-continuation-diagnostic v4-build64-phone-recovery-test v4-build64-phone-physical-test v4-build65-phone-parallel-test v4-build65-phone-physical-test v4-build66-phone-parallel-decode-test v4-build66-phone-physical-test v4-build67-phone-profile-test v4-build67-phone-physical-test v4-build68-phone-plane-reuse-test v4-build68-phone-physical-test v4-build69-phone-geometry-profile-test v4-build69-phone-physical-test v4-build70-phone-seed-genealogy-test v4-build70-phone-physical-test v4-build71-phone-gen4-parallel-test v4-build71-phone-physical-test v4-build72-phone-prefix-profile-test v4-build72-phone-physical-test v4-build73-phone-gen3-parallel-test v4-build73-phone-physical-test v4-build74-phone-freeze-profile-test v4-build74-phone-physical-test v4-build75-phone-basin-parallel-test v4-build75-phone-physical-test v4-build76-phone-gen2-parallel-test v4-build76-phone-physical-test v4-build77-phone-gen4-profile-test v4-build77-phone-physical-test v4-build78-phone-continuation4-profile-test v4-build78-phone-physical-test v4-build79-phone-foldscore-profile-test v4-build79-phone-physical-test v4-build80-phone-luminance-cache-test v4-build80-phone-physical-test v4-build81-phone-luminance-lut-test v4-build81-phone-physical-test v4-build82-phone-inline-sampler-test v4-build82-phone-physical-test v4-build83-projective-sampler-test v4-build83-projective-sampler-profile v4-build84-phone-rgb-fetch-test v4-build84-phone-rgb-fetch-benchmark v4-build84-phone-physical-test v4-build85-qualified-sampler-test v4-build85-qualified-sampler-profile v4-build86-dct-hoist-test v4-build86-dct-hoist-benchmark v4-build86-phone-physical-test v4-build87-rgb-luma-bilinear-test v4-build87-rgb-luma-bilinear-profile v4-build88-direct-rgb-test v4-build88-direct-rgb-benchmark v4-build88-phone-physical-test v4-build89-full-pipeline-test v4-build89-full-pipeline-profile v4-build90-qualification-parallel-test v4-build90-qualification-benchmark v4-build90-phone-physical-test v4-build37-physical-scanner-test v4-pilot-lock-corpus-test v4-frame-test v4-frame-corpus-test smooth-phase-test print-camera-test test-images deep-test extreme-test geometry-test affine-test composition-test lattice-test perspective-test all-test release-check version-check all build-all core-target-check vet clean

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

# Build60 research-only regression. Reproduce the complete Build59 bank and
# continue every retained third-pair state with bounded 1px proposal-only
# coordinate descent, retaining every accepted intermediate.
v4-build60-phone-third-pair-continuation-test:
	@echo "Running Build60 third-pair continuation regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build60' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build60 private study. Both image banks are frozen before oracle generation.
v4-build60-phone-third-pair-continuation-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD60_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD60_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD60_TIMEOUT="$(V4_PHONE_BUILD60_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build60-phone-third-pair-continuation.sh

# Build61 research-only regression. Reproduce the complete Build60 bank and
# evaluate all independent +/-1px siblings from every retained post-third-pair
# continuation state against the identical frozen parent.
v4-build61-phone-third-pair-sibling-stencil-test:
	@echo "Running Build61 third-pair sibling-stencil regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build61' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build61 private study. Both image banks are frozen before oracle generation.
v4-build61-phone-third-pair-sibling-stencil-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD61_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD61_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD61_TIMEOUT="$(V4_PHONE_BUILD61_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build61-phone-third-pair-sibling-stencil.sh

# Build62 research-only regression. Reproduce the complete Build61 bank, probe
# every frozen third-pair sibling for one-coordinate locality, and evaluate
# bounded coupled +/-1px fourth-pair escapes only from proposal-local siblings.
v4-build62-phone-fourth-pair-escape-test:
	@echo "Running Build62 fourth-pair escape regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build62' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build62 private study. Both image banks are frozen before oracle generation.
v4-build62-phone-fourth-pair-escape-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD62_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD62_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD62_TIMEOUT="$(V4_PHONE_BUILD62_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build62-phone-fourth-pair-escape.sh

# Build63 research-only regression. Reproduce Build62 and continue every retained
# fourth-pair state with bounded proposal-only 1px coordinate descent.
v4-build63-phone-fourth-pair-continuation-test:
	@echo "Running Build63 fourth-pair continuation regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build63' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build63 private study. Both image banks are frozen before oracle generation.
v4-build63-phone-fourth-pair-continuation-diagnostic: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_OUTPUT_DIR="$(V4_PHONE_OUTPUT_DIR)" \
	V4_PHONE_BUILD63_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD63_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD63_TIMEOUT="$(V4_PHONE_BUILD63_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	bash ./scripts/diagnose-v4-build63-phone-fourth-pair-continuation.sh

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

# Build64 qualified-baseline regression. The deep recovery is additive and
# must preserve every Build44-qualified production constant and bound.
v4-build64-phone-recovery-test:
	@echo "Running Build64 qualified deep-recovery regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build64' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1

# Build64 qualified physical regression gate over the complete nine-photo
# Build38 corpus. Existing Build44 passes must return before Build64; controls
# and B/angle must exercise the fallback and reject; B/mild must authenticate.
v4-build64-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD64_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD64_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD64_TIMEOUT="$(V4_PHONE_BUILD64_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build64-phone-corpus.sh

# Build65 performance-only candidate. Geometry, proposal ranking, freeze order,
# qualification, decode order and HMAC semantics must remain exactly Build64.
v4-build65-phone-parallel-test:
	@echo "Running Build65 deterministic seed-parallel recovery regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build6[45]' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo equivalence gate. For every deep-fallback case this
# also requires exact Build64 seed/evaluation/bank/qualification/decode telemetry.
v4-build65-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD65_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD65_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD65_TIMEOUT="$(V4_PHONE_BUILD65_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build65-phone-corpus.sh

# Build66 performance-only candidate. It preserves Build65 seed-parallel blind
# geometry and Build64 logical decode semantics, but evaluates qualified data
# candidates concurrently in deterministic candidate-order batches.
v4-build66-phone-parallel-decode-test:
	@echo "Running Build66 ordered parallel decode regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build6[456]' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo equivalence/performance gate. Semantic Build64
# telemetry must remain exact. When the Build65 baseline TSV is present, the
# report also computes same-host per-image speedup without making timing part of
# correctness.
v4-build66-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD66_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD66_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD66_TIMEOUT="$(V4_PHONE_BUILD66_TIMEOUT)" \
	V4_PHONE_BUILD65_BASELINE_TSV="$(V4_PHONE_BUILD65_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build66-phone-corpus.sh

# Build67 is an observability-only successor to the qualified Build66 baseline.
# It must preserve Build66 scheduling and exact Build64 logical semantics while
# exposing stage timings and physical-vs-logical decode work.
v4-build67-phone-profile-test:
	@echo "Running Build67 deep-recovery profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build6[4567]' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo semantic-equivalence + profiling gate. The Build66
# timing baseline is informational only; Build67 is not promoted by this gate.
v4-build67-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD67_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD67_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD67_TIMEOUT="$(V4_PHONE_BUILD67_TIMEOUT)" \
	V4_PHONE_BUILD66_BASELINE_TSV="$(V4_PHONE_BUILD66_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build67-phone-corpus.sh

# Build68 is the first implementation optimization selected by the returned
# Build67 profile. It reuses the already materialized post-freeze pixel plane
# during full-pilot qualification and must remain detector/qualification exact.
v4-build68-phone-plane-reuse-test:
	@echo "Running Build68 qualification plane-reuse regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build6[45678]' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo equivalence/performance gate. Build68 is the
# qualified baseline; reruns remain useful as reproducibility checks.
v4-build68-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD68_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD68_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD68_TIMEOUT="$(V4_PHONE_BUILD68_TIMEOUT)" \
	V4_PHONE_BUILD66_BASELINE_TSV="$(V4_PHONE_BUILD66_BASELINE_TSV)" \
	V4_PHONE_BUILD67_BASELINE_TSV="$(V4_PHONE_BUILD67_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build68-phone-corpus.sh

# Build69 is an observability-only successor to the qualified Build68 baseline.
# It wraps the exact Build68 geometry-generation path with stage/per-seed timing
# and must not change any geometry content, ordering, qualification or decode.
v4-build69-phone-geometry-profile-test:
	@echo "Running Build69 geometry-generation profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build6[456789]' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo semantic-equivalence + geometry profiling gate. The
# Build68 timing baseline is informational only; Build69 is not promoted by this
# gate and exists solely to select an evidence-based Build70 optimization.
v4-build69-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD69_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD69_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD69_TIMEOUT="$(V4_PHONE_BUILD69_TIMEOUT)" \
	V4_PHONE_BUILD68_BASELINE_TSV="$(V4_PHONE_BUILD68_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build69-phone-corpus.sh

# Build70 is observability-only. It preserves the qualified Build68 semantics
# and extends Build69 with dominant-seed identity and per-stage genealogy.
v4-build70-phone-seed-genealogy-test:
	@echo "Running Build70 dominant-seed genealogy profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(68|69|70)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo semantic-equivalence + dominant-seed genealogy gate.
# Build68 remains the qualified baseline; Build70 is diagnostic only.
v4-build70-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD70_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD70_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD70_TIMEOUT="$(V4_PHONE_BUILD70_TIMEOUT)" \
	V4_PHONE_BUILD68_BASELINE_TSV="$(V4_PHONE_BUILD68_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build70-phone-corpus.sh

# Build71 is the current qualified smartphone baseline. It was selected from
# Build69/70 physical profiles and confirmed by two retained Go 1.26.0 runs.
# It freezes the exact Build64 prefix through sib3 and schedules independent
# generation-four subtrees in one bounded pool.
v4-build71-phone-gen4-parallel-test:
	@echo "Running Build71 ordered generation-four parallel regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(68|70|71)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo reproducibility gate for the qualified Build71 baseline.
v4-build71-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD71_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD71_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD71_TIMEOUT="$(V4_PHONE_BUILD71_TIMEOUT)" \
	V4_PHONE_BUILD68_BASELINE_TSV="$(V4_PHONE_BUILD68_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build71-phone-corpus.sh

# Build72 is observability-only over the qualified Build71 baseline. It keeps
# the exact Build71 two-barrier scheduler and profiles the prefix through sib3
# by public-only stage without changing bank/order or decode semantics.
v4-build72-phone-prefix-profile-test:
	@echo "Running Build72 prefix-stage profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(71|72)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo semantic-equivalence + prefix-stage observability gate.
# Build71 remains the qualified baseline; Build72 is diagnostic only.
v4-build72-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD72_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD72_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD72_TIMEOUT="$(V4_PHONE_BUILD72_TIMEOUT)" \
	V4_PHONE_BUILD71_BASELINE_TSV="$(V4_PHONE_BUILD71_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build72-phone-corpus.sh

# Build73 is the qualified equivalence-preserving smartphone baseline selected from the
# physical Build72 prefix profile. It freezes the exact Build64 prefix through
# sibling2, evaluates independent generation-three subtrees in one bounded pool,
# commits in original order, then uses the qualified Build71 generation-four pool.
v4-build73-phone-gen3-parallel-test:
	@echo "Running Build73 ordered generation-three parallel regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(68|70|71|73)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Full retained nine-photo semantic-equivalence + performance gate for Build73.
# Build73 is qualified by two independent retained Go 1.26.0 physical runs; this target remains the regression/qualification gate.
v4-build73-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD73_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD73_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD73_TIMEOUT="$(V4_PHONE_BUILD73_TIMEOUT)" \
	V4_PHONE_BUILD71_BASELINE_TSV="$(V4_PHONE_BUILD71_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build73-phone-corpus.sh

# Build74 is observability-only over the qualified Build73 baseline. It profiles
# the exact Build47 freeze internals and must not change frozen bank/order.
v4-build74-phone-freeze-profile-test:
	@echo "Running Build74 Build47-freeze profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(68|71|73|74)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained nine-photo semantic-equivalence + freeze observability gate.
# Build73 remains the qualified baseline; Build74 is not promotable from timing.
v4-build74-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD74_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD74_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD74_TIMEOUT="$(V4_PHONE_BUILD74_TIMEOUT)" \
	V4_PHONE_BUILD73_BASELINE_TSV="$(V4_PHONE_BUILD73_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build74-phone-corpus.sh

# Build75 is the qualified ordered-parallel Build47 basin smartphone baseline.
# It parallelizes only the frozen Build47 basin task stream with ordered commit.
v4-build75-phone-basin-parallel-test:
	@echo "Running Build75 ordered-parallel Build47 basin regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained nine-photo semantic-equivalence + performance gate for Build75.
# Build75 was qualified by two independent retained Go 1.26.0 physical runs; this target remains its regression gate.
v4-build75-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD75_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD75_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD75_TIMEOUT="$(V4_PHONE_BUILD75_TIMEOUT)" \
	V4_PHONE_BUILD73_BASELINE_TSV="$(V4_PHONE_BUILD73_BASELINE_TSV)" \
	V4_PHONE_BUILD74_BASELINE_TSV="$(V4_PHONE_BUILD74_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build75-phone-corpus.sh

# Build76 is the current qualified equivalence-preserving smartphone baseline.
# It moves only the ordered public-geometry barrier from sibling2 to sibling1.
v4-build76-phone-gen2-parallel-test:
	@echo "Running Build76 ordered generation-two parallel regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained nine-photo semantic-equivalence + performance gate for Build76.
# Build75 is retained as the historical timing baseline for Build76 reproducibility checks.
v4-build76-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD76_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD76_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD76_TIMEOUT="$(V4_PHONE_BUILD76_TIMEOUT)" \
	V4_PHONE_BUILD75_BASELINE_TSV="$(V4_PHONE_BUILD75_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build76-phone-corpus.sh

# Build77 is observability-only over the qualified Build76 baseline. It decomposes
# existing generation-four work without changing bank/order, scheduling or decode.
v4-build77-phone-gen4-profile-test:
	@echo "Running Build77 generation-four stage profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|77)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained nine-photo semantic-equivalence profiling gate for Build77. Timing is
# diagnostic only; Build76 remains the qualified baseline.
v4-build77-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD77_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD77_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD77_TIMEOUT="$(V4_PHONE_BUILD77_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build77-phone-corpus.sh

# Build78 is observability-only over the qualified Build76 baseline. It profiles
# the existing continuation4 coordinate-descent loop without changing its outputs.
v4-build78-phone-continuation4-profile-test:
	@echo "Running Build78 continuation4 internal profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|78)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained nine-photo semantic-equivalence profiling gate for Build78. Timing is
# diagnostic only; Build76 remains the qualified baseline.
v4-build78-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD78_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD78_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD78_TIMEOUT="$(V4_PHONE_BUILD78_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build78-phone-corpus.sh

# Build79 is observability-only over the qualified Build76 baseline. It profiles
# the FoldScore kernel used by continuation4 without changing arithmetic or output.
v4-build79-phone-foldscore-profile-test:
	@echo "Running Build79 FoldScore kernel profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|79)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained nine-photo semantic-equivalence profiling gate for Build79. Timing is
# diagnostic only; Build76 remains the qualified baseline.
v4-build79-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD79_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD79_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD79_TIMEOUT="$(V4_PHONE_BUILD79_TIMEOUT)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build79-phone-corpus.sh

# Build80 is an equivalence-preserving performance candidate over qualified Build76.
# It changes only continuation4 FoldScore block sampling by reusing exact local
# integer-source luminance values inside a bounded block-local cache.
v4-build80-phone-luminance-cache-test:
	@echo "Running Build80 exact local luminance-cache regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|80)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained historical nine-photo semantic-equivalence/performance gate for Build80.
# Build80 preserved semantics but regressed performance and is closed as rejected; Build76 remains qualified.
v4-build80-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD80_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD80_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD80_TIMEOUT="$(V4_PHONE_BUILD80_TIMEOUT)" \
	V4_PHONE_BUILD76_BASELINE_TSV="$(V4_PHONE_BUILD76_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build80-phone-corpus.sh

# Build81 is a closed equivalence-preserving performance experiment over qualified Build76.
# It changed only continuation4 FoldScore RGB->luminance products using exact
# 256-entry float64 lookup tables; two 9/9 semantic PASS runs did not show repeatable speedup.
v4-build81-phone-luminance-lut-test:
	@echo "Running Build81 exact luminance-LUT regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|81)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check

# Retained nine-photo semantic-equivalence/performance gate for Build81.
# Build76 remains the qualified baseline; Build81 is retained as semantic PASS / non-promoted.
v4-build81-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD81_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD81_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD81_TIMEOUT="$(V4_PHONE_BUILD81_TIMEOUT)" \
	V4_PHONE_BUILD76_BASELINE_TSV="$(V4_PHONE_BUILD76_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build81-phone-corpus.sh

# Build82 is a closed equivalence-preserving experiment over qualified Build76.
# Two physical 9/9 PASS runs preserved exact semantics but did not show a
# repeatable speedup. Go 1.26.0 leaves samplePlaneLuminance/readProjectiveBlockValue uninlined; Build82
# specializes only continuation4 block reads by manually incorporating the exact
# samplePlaneLuminance body while retaining homography.mapPoint unchanged.
v4-build82-phone-inline-sampler-test:
	@echo "Running Build82 exact inline-sampler regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Retained nine-photo semantic-equivalence/performance gate for Build82.
# Build76 remains the qualified baseline. Promotion requires two independent
# physical PASS runs with a material, repeatable performance improvement.
v4-build82-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD82_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD82_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD82_TIMEOUT="$(V4_PHONE_BUILD82_TIMEOUT)" \
	V4_PHONE_BUILD76_BASELINE_TSV="$(V4_PHONE_BUILD76_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build82-phone-corpus.sh


# Build83 is observability-only after Build82 closed as exact but non-promoted.
# Runtime recovery is restored to qualified Build76. The benchmark fixture is
# public/deterministic and contains no private acquisition, key, payload or HMAC
# selector. Build82 remains only as an exact comparator inside the benchmark.
v4-build83-projective-sampler-test:
	@echo "Running Build83 public projective-sampler fixture/regression checks..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Stable microbenchmark + CPU profile for the historical Build76 projective
# reader. Timing from this target is diagnostic only and is never a semantic or
# promotion gate by itself.
v4-build83-projective-sampler-profile:
	@PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_BUILD83_PROFILE_DIR="$(V4_BUILD83_PROFILE_DIR)" \
	V4_BUILD83_BENCHTIME="$(V4_BUILD83_BENCHTIME)" \
	V4_BUILD83_COUNT="$(V4_BUILD83_COUNT)" \
	V4_BUILD83_PPROF_TIME="$(V4_BUILD83_PPROF_TIME)" \
	bash ./scripts/profile-v4-build83-projective-sampler.sh


# Build84 is the current qualified smartphone baseline. It changes only
# continuation4 RGB address/bounds-check shape; all arithmetic and semantics are exact.
v4-build84-phone-rgb-fetch-test:
	@echo "Running Build84 exact RGB-fetch regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83|84)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Public deterministic microbenchmark plus Go compiler bounds-check evidence.
# Timing is diagnostic: a slow result does not make the correctness target fail.
v4-build84-phone-rgb-fetch-benchmark:
	@PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_BUILD84_BENCH_DIR="$(V4_BUILD84_BENCH_DIR)" \
	V4_BUILD84_BENCHTIME="$(V4_BUILD84_BENCHTIME)" \
	V4_BUILD84_COUNT="$(V4_BUILD84_COUNT)" \
	bash ./scripts/benchmark-v4-build84-rgb-fetch.sh

# Nine-photo semantic-equivalence/performance requalification gate for Build84.
# The original promotion used two independent Go 1.26.0 runs; reruns remain diagnostic.
v4-build84-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD84_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD84_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD84_TIMEOUT="$(V4_PHONE_BUILD84_TIMEOUT)" \
	V4_PHONE_BUILD76_BASELINE_TSV="$(V4_PHONE_BUILD76_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build84-phone-corpus.sh

# Build85 is observability-only over the qualified Build84 runtime. It does not
# change recovery behavior; it re-profiles the promoted sampler before any new
# performance candidate is selected.
v4-build85-qualified-sampler-test:
	@echo "Running Build85 qualified-sampler public profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83|84|85)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Public deterministic repeated benchmarks + CPU pprof of the current Build84
# qualified block reader. Timing is observability evidence only.
v4-build85-qualified-sampler-profile:
	@PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_BUILD85_PROFILE_DIR="$(V4_BUILD85_PROFILE_DIR)" \
	V4_BUILD85_BENCHTIME="$(V4_BUILD85_BENCHTIME)" \
	V4_BUILD85_COUNT="$(V4_BUILD85_COUNT)" \
	V4_BUILD85_PPROF_TIME="$(V4_BUILD85_PPROF_TIME)" \
	bash ./scripts/profile-v4-build85-qualified-sampler.sh


# Build86 is an exact performance candidate over the current qualified Build84
# runtime. Only DCT cosine-table load placement changes inside continuation4.
v4-build86-dct-hoist-test:
	@echo "Running Build86 exact DCT table-hoist regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83|84|85|86)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Benchmark-first decision gate. Timing itself is diagnostic; correctness is
# established by the exact fixture/regressions before any physical run.
v4-build86-dct-hoist-benchmark:
	@PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_BUILD86_BENCH_DIR="$(V4_BUILD86_BENCH_DIR)" \
	V4_BUILD86_BENCHTIME="$(V4_BUILD86_BENCHTIME)" \
	V4_BUILD86_COUNT="$(V4_BUILD86_COUNT)" \
	bash ./scripts/benchmark-v4-build86-dct-hoist.sh

# Nine-photo semantic/performance gate. Run only after the qualified-host
# benchmark shows a stable, useful whole-reader improvement over Build84.
v4-build86-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD86_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD86_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD86_TIMEOUT="$(V4_PHONE_BUILD86_TIMEOUT)" \
	V4_PHONE_BUILD84_BASELINE_TSV="$(V4_PHONE_BUILD84_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build86-phone-corpus.sh

# Build87 is observability-only over the restored qualified Build84 runtime.
# Build86 is retained as exact/benchmark-negative evidence and is not active.
v4-build87-rgb-luma-bilinear-test:
	@echo "Running Build87 qualified RGB/luminance/bilinear profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83|84|85|86|87)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Repeated public benchmarks, integrated CPU pprof, BCE/inliner evidence and
# objdump around the qualified Build84 RGB/luminance/bilinear hot lines.
v4-build87-rgb-luma-bilinear-profile:
	@PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_BUILD87_PROFILE_DIR="$(V4_BUILD87_PROFILE_DIR)" \
	V4_BUILD87_BENCHTIME="$(V4_BUILD87_BENCHTIME)" \
	V4_BUILD87_COUNT="$(V4_BUILD87_COUNT)" \
	V4_BUILD87_PPROF_TIME="$(V4_BUILD87_PPROF_TIME)" \
	bash ./scripts/profile-v4-build87-rgb-luma-bilinear.sh

# Build88 is the final narrow exact-portable reader candidate selected by
# Build87: replace four three-byte RGB slices with dominating index+2 proofs
# plus direct scalar byte loads. Build84 remains the qualified baseline.
v4-build88-direct-rgb-test:
	@echo "Running Build88 exact direct-RGB / dominating-BCE regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83|84|85|86|87|88)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Same-session Build84-vs-Build88 public reader benchmark plus BCE/inliner and
# objdump evidence. Review complete-reader speed before any private corpus run.
v4-build88-direct-rgb-benchmark:
	@PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_BUILD88_BENCH_DIR="$(V4_BUILD88_BENCH_DIR)" \
	V4_BUILD88_BENCHTIME="$(V4_BUILD88_BENCHTIME)" \
	V4_BUILD88_COUNT="$(V4_BUILD88_COUNT)" \
	bash ./scripts/benchmark-v4-build88-direct-rgb.sh

# Nine-photo semantic/performance gate. Run only if the Go 1.26.0 benchmark
# shows a stable useful reader-level advantage over qualified Build84.
v4-build88-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD88_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD88_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD88_TIMEOUT="$(V4_PHONE_BUILD88_TIMEOUT)" \
	V4_PHONE_BUILD84_BASELINE_TSV="$(V4_PHONE_BUILD84_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build88-phone-corpus.sh

# Build89 closes the reader micro-optimization branch after Build88 and restores
# the current qualified Build84 runtime. It adds only derived phase accounting
# from timing telemetry Build84 already records.
v4-build89-full-pipeline-test:
	@echo "Running Build89 qualified Build84 full-pipeline profiling regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83|84|85|86|87|88|89)' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Retained nine-photo workload profile. This is observability-only, not a
# promotion/qualification target: semantic gates are retained only to prove that
# the profiling snapshot still executes the qualified Build84 behavior.
v4-build89-full-pipeline-profile: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD89_PROFILE_DIR="$(V4_PHONE_BUILD89_PROFILE_DIR)" \
	V4_PHONE_BUILD89_TIMEOUT="$(V4_PHONE_BUILD89_TIMEOUT)" \
	V4_PHONE_BUILD84_BASELINE_TSV="$(V4_PHONE_BUILD84_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/profile-v4-build89-full-pipeline.sh

# Build90 is the first pipeline-level performance candidate after Build89.
# It parallelizes only qualification over the already-frozen bank and commits
# results strictly in original bank order. Build84 remains qualified.
v4-build90-qualification-parallel-test:
	@echo "Running Build90 ordered-parallel qualification regressions..."
	@$(GO) test ./watermark -run '^TestExperimentalV4Build(47|68|71|73|75|76|82|83|84|85|86|87|88|89|90)' -count=1
	@$(GO) test -race ./watermark -run '^TestExperimentalV4Build90ParallelQualificationMatchesSerial$$' -count=1
	@$(GO) test ./cmd/pixseal -run '^TestSubcommandHelpReturnsFlagErrHelp$$' -count=1
	@$(MAKE) --no-print-directory version-check
	@$(GO) vet ./...

# Public deterministic high-pass and early-reject qualification workloads.
# Timing is a decision aid only; exactness/race gates remain authoritative.
v4-build90-qualification-benchmark:
	@PIXSEAL_GO_TOOLCHAIN="$(PIXSEAL_GO_TOOLCHAIN)" \
	V4_BUILD90_BENCH_DIR="$(V4_BUILD90_BENCH_DIR)" \
	V4_BUILD90_BENCHTIME="$(V4_BUILD90_BENCHTIME)" \
	V4_BUILD90_COUNT="$(V4_BUILD90_COUNT)" \
	bash ./scripts/benchmark-v4-build90-parallel-qualification.sh

# Deferred nine-photo semantic/performance gate. Run only after the qualified
# host benchmark shows a large high-pass gain without material early-reject cost.
v4-build90-phone-physical-test: build
	@PIXSEAL="$(abspath $(PIXSEAL))" \
	V4_PHONE_ACQUISITION_DIR="$(V4_PHONE_ACQUISITION_DIR)" \
	V4_PHONE_BUILD90_DIAGNOSTIC_DIR="$(V4_PHONE_BUILD90_DIAGNOSTIC_DIR)" \
	V4_PHONE_BUILD90_TIMEOUT="$(V4_PHONE_BUILD90_TIMEOUT)" \
	V4_PHONE_BUILD84_BASELINE_TSV="$(V4_PHONE_BUILD84_BASELINE_TSV)" \
	V4_PHONE_KEY="$(V4_PHONE_KEY)" \
	V4_PHONE_CANONICAL_WIDTH="$(V4_PHONE_CANONICAL_WIDTH)" \
	V4_PHONE_CANONICAL_HEIGHT="$(V4_PHONE_CANONICAL_HEIGHT)" \
	V4_PHONE_MESSAGE_A="$(V4_PHONE_MESSAGE_A)" \
	V4_PHONE_MESSAGE_B="$(V4_PHONE_MESSAGE_B)" \
	bash ./scripts/test-v4-build90-phone-corpus.sh

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
