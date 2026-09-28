package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build71ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build71 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build71 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build71BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, profile := experimentalV4PhoneBuild71BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if profile != (experimentalV4PhoneBuild71GeometryTelemetry{}) {
		t.Fatalf("invalid boundary produced profile=%+v", profile)
	}
}

func TestExperimentalV4Build71SeedReassemblyMatchesBuild64Exactly(t *testing.T) {
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

	wantBank, wantEvals := experimentalV4PhoneBuild64SeedBank(plane, candidate, w, h, anchor, seed)
	inputs, gotEvals := experimentalV4PhoneBuild71SeedPrefix(plane, candidate, w, h, anchor, seed)
	gotBank := make([]experimentalV4PhoneHypothesis, 0, len(wantBank))
	for _, input := range inputs {
		bank, n := experimentalV4PhoneBuild71Generation4(plane, candidate, w, h, anchor, input)
		gotEvals += n
		gotBank = append(gotBank, bank...)
	}
	if gotEvals != wantEvals || len(gotBank) != len(wantBank) {
		t.Fatalf("Build71 seed search differs: evals=%d/%d bank=%d/%d inputs=%d", gotEvals, wantEvals, len(gotBank), len(wantBank), len(inputs))
	}
	for i := range gotBank {
		if gotBank[i].quad != wantBank[i].quad || gotBank[i].h != wantBank[i].h || gotBank[i].proposal != wantBank[i].proposal {
			t.Fatalf("Build71 bank differs at %d", i)
		}
	}
}

func TestExperimentalV4Build71ApplyTelemetryPreservesQualifiedBuild68View(t *testing.T) {
	recovery := experimentalV4PhoneBuild71RecoveryTelemetry{}
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
	recovery.GeometryProfile.PlanePrepElapsed = 200 * time.Millisecond
	recovery.GeometryProfile.FreezeElapsed = 60 * time.Second
	recovery.GeometryProfile.PrefixWallElapsed = 80 * time.Second
	recovery.GeometryProfile.PrefixWorkerElapsed = 100 * time.Second
	recovery.GeometryProfile.Gen4WallElapsed = 40 * time.Second
	recovery.GeometryProfile.Gen4WorkerElapsed = 300 * time.Second
	recovery.GeometryProfile.Gen4MinElapsed = time.Millisecond
	recovery.GeometryProfile.Gen4MedianElapsed = 20 * time.Millisecond
	recovery.GeometryProfile.Gen4MaxElapsed = 400 * time.Millisecond
	recovery.GeometryProfile.Gen4MinEvaluations = 16
	recovery.GeometryProfile.Gen4MaxEvaluations = 512
	recovery.GeometryProfile.Gen4MinBank = 0
	recovery.GeometryProfile.Gen4MaxBank = 8

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild71ApplyTelemetry(&public, recovery)
	if !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build65Attempted || !public.Build64Attempted {
		t.Fatalf("Build71 did not preserve inherited attempted telemetry")
	}
	if public.Build64SeedsSelected != 24 || public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build71PrefixWorkers != 8 || public.Build71Gen4Workers != 8 || public.Build71Gen4Tasks != 593 || public.Build71GeometryFreezeMs != 60000 || public.Build71PrefixWallMs != 80000 || public.Build71Gen4WallMs != 40000 || public.Build71Gen4MaxEvals != 512 || public.Build71Gen4MaxBank != 8 {
		t.Fatalf("Build71 scheduling telemetry mismatch: %+v", public)
	}
}
