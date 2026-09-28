package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build72ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build72 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build72 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build72BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, profile := experimentalV4PhoneBuild72BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if profile != (experimentalV4PhoneBuild72GeometryTelemetry{}) {
		t.Fatalf("invalid boundary produced profile=%+v", profile)
	}
}

func TestExperimentalV4Build72SeedPrefixMatchesBuild71Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	anchor := [4]ImagePoint{{X: 0, Y: 0}, {X: float64(w - 1), Y: 0}, {X: 0, Y: float64(h - 1)}, {X: float64(w - 1), Y: float64(h - 1)}}
	seed := experimentalV4PhoneBuild48Seed{index: 3, rankWithinPair: 2, frozen: experimentalV4PhoneBuild47Frozen{h: experimentalV4PhoneHypothesis{quad: anchor, build43PairRank: 1}}}
	plane := newPixelPlane(carrier)

	want, wantEvals := experimentalV4PhoneBuild71SeedPrefix(plane, candidate, w, h, anchor, seed)
	got, gotEvals, stage := experimentalV4PhoneBuild72SeedPrefixProfile(plane, candidate, w, h, anchor, seed)
	if gotEvals != wantEvals || len(got) != len(want) {
		t.Fatalf("Build72 prefix differs: evals=%d/%d inputs=%d/%d", gotEvals, wantEvals, len(got), len(want))
	}
	stageTotal := 0
	for _, n := range stage.Evals {
		stageTotal += n
	}
	if stageTotal != gotEvals {
		t.Fatalf("Build72 stage evals sum=%d want %d", stageTotal, gotEvals)
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || got[i].proposal != want[i].proposal {
			t.Fatalf("Build72 prefix input differs at %d", i)
		}
	}
}

func TestExperimentalV4Build72ApplyTelemetryPreservesQualifiedBuild71View(t *testing.T) {
	recovery := experimentalV4PhoneBuild72RecoveryTelemetry{}
	recovery.Attempted = true
	recovery.SeedsSelected = 24
	recovery.GeometryEvaluations = 79259
	recovery.BankCandidates = 937
	recovery.QualifiedCandidates = 935
	recovery.DecodeCandidatesTried = 691
	recovery.ListFramesTried = 2120047
	recovery.Authenticated = true
	recovery.Workers = 8
	recovery.DecodeWorkers = 8
	recovery.PhysicalDecodeCandidates = 696
	recovery.TotalElapsed = 300 * time.Second
	recovery.GeometryElapsed = 200 * time.Second
	recovery.PlanePrepElapsed = 250 * time.Millisecond
	recovery.QualificationElapsed = 60 * time.Second
	recovery.DecodeWallElapsed = 30 * time.Second
	recovery.GeometryProfile.PrefixWorkers = 8
	recovery.GeometryProfile.Gen4Workers = 8
	recovery.GeometryProfile.Gen4Tasks = 593
	recovery.GeometryProfile.FreezeElapsed = 60 * time.Second
	recovery.GeometryProfile.PrefixWallElapsed = 80 * time.Second
	recovery.GeometryProfile.PrefixWorkerElapsed = 100 * time.Second
	recovery.GeometryProfile.Gen4WallElapsed = 20 * time.Second
	recovery.GeometryProfile.Gen4WorkerElapsed = 150 * time.Second
	recovery.PrefixProfile.PrefixMinElapsed = time.Second
	recovery.PrefixProfile.PrefixMedianElapsed = 2 * time.Second
	recovery.PrefixProfile.PrefixMaxElapsed = 75 * time.Second
	recovery.PrefixProfile.PrefixMinEvaluations = 10
	recovery.PrefixProfile.PrefixMaxEvaluations = 37000
	recovery.PrefixProfile.PrefixMinInputs = 0
	recovery.PrefixProfile.PrefixMaxInputs = 593
	recovery.PrefixProfile.MaxPrefixSeedIndex = 67
	recovery.PrefixProfile.MaxPrefixSeedPairRank = 3
	recovery.PrefixProfile.MaxPrefixSeedRankWithinPair = 3
	recovery.PrefixProfile.MaxPrefixSeedEvaluations = 32146
	recovery.PrefixProfile.MaxPrefixSeedInputs = 593
	recovery.PrefixProfile.MaxPrefixSeedElapsed = 74 * time.Second
	recovery.PrefixProfile.PrefixStageEvaluations[13] = 1234
	recovery.PrefixProfile.PrefixStageWorkerElapsed[13] = 12 * time.Second
	recovery.PrefixProfile.MaxPrefixSeedStage.Evals[13] = 1000
	recovery.PrefixProfile.MaxPrefixSeedStage.States[9] = 593
	recovery.PrefixProfile.MaxPrefixSeedStage.Elapsed[13] = 10 * time.Second

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild72ApplyTelemetry(&public, recovery)
	if !public.Build72Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build72 did not preserve inherited attempted telemetry")
	}
	if public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build72MaxPrefixSeedIndex != 67 || public.Build72MaxPrefixSeedPairRank != 3 || public.Build72MaxPrefixSeedRankWithinPair != 3 || public.Build72MaxPrefixSeedEvals != 32146 || public.Build72MaxPrefixSeedInputs != 593 || public.Build72PrefixMaxMs != 75000 {
		t.Fatalf("Build72 prefix telemetry mismatch: %+v", public)
	}
	if public.Build72PrefixStageEvals == "" || public.Build72PrefixStageWorkerMs == "" || public.Build72MaxPrefixSeedStageEvals == "" || public.Build72MaxPrefixSeedStageMs == "" {
		t.Fatalf("Build72 stage telemetry missing: %+v", public)
	}
}
