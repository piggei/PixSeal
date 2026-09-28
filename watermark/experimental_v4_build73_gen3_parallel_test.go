package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build73ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build73 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build73 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build73BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, profile := experimentalV4PhoneBuild73BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if profile != (experimentalV4PhoneBuild73GeometryTelemetry{}) {
		t.Fatalf("invalid boundary produced profile=%+v", profile)
	}
}

func TestExperimentalV4Build73SeedReassemblyMatchesBuild64Exactly(t *testing.T) {
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
	s2, gotEvals := experimentalV4PhoneBuild73SeedPrefix2(plane, candidate, w, h, anchor, seed)
	s3 := make([]experimentalV4PhoneHypothesis, 0, 64)
	for _, input := range s2 {
		out, n := experimentalV4PhoneBuild73Generation3(plane, candidate, w, h, anchor, input)
		gotEvals += n
		s3 = append(s3, out...)
	}
	gotBank := make([]experimentalV4PhoneHypothesis, 0, len(wantBank))
	for _, input := range s3 {
		bank, n := experimentalV4PhoneBuild71Generation4(plane, candidate, w, h, anchor, input)
		gotEvals += n
		gotBank = append(gotBank, bank...)
	}
	if gotEvals != wantEvals || len(gotBank) != len(wantBank) {
		t.Fatalf("Build73 seed search differs: evals=%d/%d bank=%d/%d s2=%d s3=%d", gotEvals, wantEvals, len(gotBank), len(wantBank), len(s2), len(s3))
	}
	for i := range gotBank {
		if gotBank[i].quad != wantBank[i].quad || gotBank[i].h != wantBank[i].h || gotBank[i].proposal != wantBank[i].proposal {
			t.Fatalf("Build73 bank differs at %d", i)
		}
	}
}

func TestExperimentalV4Build73ApplyTelemetryPreservesQualifiedView(t *testing.T) {
	recovery := experimentalV4PhoneBuild73RecoveryTelemetry{}
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
	recovery.TotalElapsed = 250 * time.Second
	recovery.GeometryElapsed = 140 * time.Second
	recovery.PlanePrepElapsed = 250 * time.Millisecond
	recovery.QualificationElapsed = 60 * time.Second
	recovery.DecodeWallElapsed = 30 * time.Second
	p := &recovery.GeometryProfile
	p.Prefix2Workers = 8
	p.Gen3Workers = 8
	p.Gen4Workers = 8
	p.Gen3Tasks = 593
	p.Gen4Tasks = 593
	p.PlanePrepElapsed = 200 * time.Millisecond
	p.FreezeElapsed = 50 * time.Second
	p.Prefix2WallElapsed = 40 * time.Second
	p.Prefix2WorkerElapsed = 70 * time.Second
	p.Gen3WallElapsed = 20 * time.Second
	p.Gen3WorkerElapsed = 150 * time.Second
	p.Gen4WallElapsed = 18 * time.Second
	p.Gen4WorkerElapsed = 140 * time.Second
	p.Gen3MaxElapsed = 400 * time.Millisecond
	p.Gen3MaxEvaluations = 512
	p.Gen3MaxOutputs = 8

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild73ApplyTelemetry(&public, recovery)
	if !public.Build73Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build73 did not preserve inherited attempted telemetry")
	}
	if public.Build64SeedsSelected != 24 || public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build73Prefix2Workers != 8 || public.Build73Gen3Workers != 8 || public.Build73Gen4Workers != 8 || public.Build73Gen3Tasks != 593 || public.Build73Gen3WallMs != 20000 || public.Build73Gen3MaxEvals != 512 || public.Build73Gen3MaxOutputs != 8 {
		t.Fatalf("Build73 scheduling telemetry mismatch: %+v", public)
	}
}
