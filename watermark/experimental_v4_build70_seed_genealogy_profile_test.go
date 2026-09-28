package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build70ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build70 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build70 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build70BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, profile := experimentalV4PhoneBuild70BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if profile != (experimentalV4PhoneBuild70GeometryTelemetry{}) {
		t.Fatalf("invalid boundary produced profile=%+v", profile)
	}
}

func TestExperimentalV4Build70SeedBankProfileMatchesBuild64Exactly(t *testing.T) {
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
	gotBank, gotEvals, profile := experimentalV4PhoneBuild70SeedBankProfile(plane, candidate, w, h, anchor, seed)
	if gotEvals != wantEvals || len(gotBank) != len(wantBank) {
		t.Fatalf("Build70 seed search differs: evals=%d/%d bank=%d/%d", gotEvals, wantEvals, len(gotBank), len(wantBank))
	}
	if profile.States[12] != len(gotBank) {
		t.Fatalf("Build70 final stage bank=%d want %d", profile.States[12], len(gotBank))
	}
	stageEvals := 0
	for _, n := range profile.Evals {
		stageEvals += n
	}
	if stageEvals != gotEvals {
		t.Fatalf("Build70 stage eval sum=%d want %d", stageEvals, gotEvals)
	}
	for i := range gotBank {
		if gotBank[i].quad != wantBank[i].quad || gotBank[i].h != wantBank[i].h || gotBank[i].proposal != wantBank[i].proposal {
			t.Fatalf("Build70 bank differs at %d", i)
		}
	}
}

func TestExperimentalV4Build70ApplyTelemetryPreservesBuild69AndBuild68Views(t *testing.T) {
	recovery := experimentalV4PhoneBuild70RecoveryTelemetry{}
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
	recovery.TotalElapsed = 100 * time.Second
	recovery.GeometryElapsed = 80 * time.Second
	recovery.PlanePrepElapsed = 250 * time.Millisecond
	recovery.QualificationElapsed = 10 * time.Second
	recovery.DecodeWallElapsed = 9 * time.Second
	recovery.GeometryPlanePrepElapsed = 200 * time.Millisecond
	recovery.GeometryFreezeElapsed = 3 * time.Second
	recovery.GeometrySeedWallElapsed = 76 * time.Second
	recovery.GeometrySeedWorkerElapsed = 590 * time.Second
	recovery.GeometrySeedMinElapsed = 10 * time.Second
	recovery.GeometrySeedMedianElapsed = 20 * time.Second
	recovery.GeometrySeedMaxElapsed = 40 * time.Second
	recovery.GeometrySeedMinEvaluations = 100
	recovery.GeometrySeedMaxEvaluations = 10000
	recovery.GeometrySeedMinBank = 0
	recovery.GeometrySeedMaxBank = 200
	recovery.GeometryFreezeEvaluations = 5000
	recovery.GeometrySeedEvaluations = 74259
	recovery.GeometryMaxSeedIndex = 17
	recovery.GeometryMaxSeedPairRank = 4
	recovery.GeometryMaxSeedRankWithinPair = 2
	recovery.GeometryMaxSeedEvaluations = 70000
	recovery.GeometryMaxSeedBank = 900
	recovery.GeometryMaxSeedElapsed = 39 * time.Second
	for i := range recovery.GeometryMaxSeedStage.Evals {
		recovery.GeometryMaxSeedStage.Evals[i] = i + 1
	}
	for i := range recovery.GeometryMaxSeedStage.States {
		recovery.GeometryMaxSeedStage.States[i] = (i + 1) * 10
	}

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild70ApplyTelemetry(&public, recovery)
	if !public.Build70Attempted || !public.Build69Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build65Attempted || !public.Build64Attempted {
		t.Fatalf("Build70 did not preserve inherited attempted telemetry")
	}
	if public.Build64SeedsSelected != 24 || public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build70GeometryFreezeEvals != 5000 || public.Build70GeometrySeedEvals != 74259 || public.Build70MaxSeedIndex != 17 || public.Build70MaxSeedPairRank != 4 || public.Build70MaxSeedRankWithinPair != 2 || public.Build70MaxSeedEvals != 70000 || public.Build70MaxSeedBank != 900 || public.Build70MaxSeedMs != 39000 {
		t.Fatalf("Build70 genealogy mismatch: %+v", public)
	}
	if public.Build70MaxSeedStageEvals != "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17" {
		t.Fatalf("Build70 stage eval CSV=%q", public.Build70MaxSeedStageEvals)
	}
	if public.Build70MaxSeedStageStates != "10,20,30,40,50,60,70,80,90,100,110,120,130" {
		t.Fatalf("Build70 stage state CSV=%q", public.Build70MaxSeedStageStates)
	}
}
