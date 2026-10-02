package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build78ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build78 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build78 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build78BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, gp, fp, g4, c4 := experimentalV4PhoneBuild78BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if gp != (experimentalV4PhoneBuild76GeometryTelemetry{}) || fp != (experimentalV4PhoneBuild75FreezeTelemetry{}) || g4 != (experimentalV4PhoneBuild78Gen4Profile{}) || c4 != (experimentalV4PhoneBuild78Continuation4Profile{}) {
		t.Fatalf("invalid boundary produced telemetry")
	}
}

func TestExperimentalV4Build78BlindBankMatchesBuild76Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, wantEvals, wantSeeds, wantWorkers, _, _ := experimentalV4PhoneBuild76BlindBank(carrier, boundary, w, h)
	got, gotEvals, gotSeeds, gotWorkers, gp, fp, g4, c4 := experimentalV4PhoneBuild78BlindBank(carrier, boundary, w, h)
	if gotEvals != wantEvals || gotSeeds != wantSeeds || gotWorkers != wantWorkers || len(got) != len(want) {
		t.Fatalf("blind bank differs evals=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", gotEvals, wantEvals, gotSeeds, wantSeeds, gotWorkers, wantWorkers, len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || got[i].proposal != want[i].proposal || got[i].validation != want[i].validation {
			t.Fatalf("blind bank differs at %d", i)
		}
	}
	if gp.Gen2Tasks == 0 || gp.Gen2Workers == 0 || gp.Gen2WallElapsed <= 0 || gp.Gen3Tasks == 0 {
		t.Fatalf("incomplete Build76-compatible geometry telemetry: %+v", gp)
	}
	if fp.BasinTasks != 16 || fp.BasinWorkers == 0 {
		t.Fatalf("Build75 freeze telemetry not preserved: %+v", fp)
	}
	if g4.Tasks != gp.Gen4Tasks || g4.Workers != gp.Gen4Workers || g4.SingleCalls != gp.Gen4Tasks {
		t.Fatalf("Build78 task accounting mismatch g4=%+v gp=%+v", g4, gp)
	}
	if gp.Gen4Tasks > 0 {
		if g4.SingleEvaluations+g4.PairEvaluations+g4.ContinueEvaluations <= 0 {
			t.Fatalf("Build78 stage evaluations not recorded: %+v", g4)
		}
		if g4.SingleWorkerElapsed <= 0 || g4.PairWorkerElapsed < 0 || g4.ContinueWorkerElapsed < 0 {
			t.Fatalf("Build78 stage timing not recorded: %+v", g4)
		}
		if c4.Calls != g4.ContinueCalls || c4.InputStates != g4.ContinueCalls || c4.ScoreEvaluations != g4.ContinueEvaluations || c4.AcceptedStates != g4.ContinueOutputs {
			t.Fatalf("Build78 continuation accounting mismatch c4=%+v g4=%+v", c4, g4)
		}
		if c4.ProbeAttempts < c4.ScoreEvaluations || c4.WorkerElapsed <= 0 || c4.ScoreWorkerElapsed <= 0 {
			t.Fatalf("Build78 continuation internals not recorded: %+v", c4)
		}
	}
}

func TestExperimentalV4Build78ContinueProfiledMatchesBuild55Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, height := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	plane := newPixelPlane(carrier)
	anchor := [4]ImagePoint{{X: 0, Y: 0}, {X: float64(w - 1), Y: 0}, {X: 0, Y: float64(height - 1)}, {X: float64(w - 1), Y: float64(height - 1)}}
	hom, ok := experimentalV4PhoneQuadHomography(w, height, anchor)
	if !ok {
		t.Fatal("anchor homography")
	}
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, height, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, wantEvals := experimentalV4PhoneBuild55Continue(plane, candidate, w, height, anchor, start, 0)
	got, gotEvals, profile := experimentalV4PhoneBuild78ContinueProfiled(plane, candidate, w, height, anchor, start, 0)
	if gotEvals != wantEvals || len(got) != len(want) {
		t.Fatalf("continuation differs evals=%d/%d states=%d/%d", gotEvals, wantEvals, len(got), len(want))
	}
	for i := range got {
		if got[i].index != want[i].index || got[i].pass != want[i].pass || got[i].dim != want[i].dim || got[i].delta != want[i].delta || got[i].hyp.quad != want[i].hyp.quad || got[i].hyp.h != want[i].hyp.h || got[i].hyp.proposal != want[i].hyp.proposal {
			t.Fatalf("continuation differs at %d", i)
		}
	}
	if profile.evals != wantEvals || profile.acceptedStates != len(want) || profile.probeAttempts < wantEvals {
		t.Fatalf("profile accounting mismatch: %+v", profile)
	}
}

func TestExperimentalV4Build78ApplyTelemetryPreservesQualifiedView(t *testing.T) {
	recovery := experimentalV4PhoneBuild78RecoveryTelemetry{}
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
	recovery.FreezeProfile.TotalElapsed = 23 * time.Second
	recovery.FreezeProfile.BasinWorkers = 8
	recovery.FreezeProfile.BasinTasks = 16
	recovery.FreezeProfile.BasinWallElapsed = 10 * time.Second
	gp := &recovery.Build76GeometryProfile
	gp.Prefix1Workers = 8
	gp.Gen2Workers = 8
	gp.Gen3Workers = 8
	gp.Gen4Workers = 8
	gp.Gen2Tasks = 93
	gp.Gen3Tasks = 275
	gp.Gen4Tasks = 593
	gp.FreezeElapsed = 23 * time.Second
	gp.Prefix1WallElapsed = 4 * time.Second
	gp.Prefix1WorkerElapsed = 20 * time.Second
	gp.Gen2WallElapsed = 5 * time.Second
	gp.Gen2WorkerElapsed = 27 * time.Second
	gp.Gen3WallElapsed = 10 * time.Second
	gp.Gen3WorkerElapsed = 80 * time.Second
	gp.Gen4WallElapsed = 20 * time.Second
	gp.Gen4WorkerElapsed = 160 * time.Second
	compat := &recovery.GeometryProfile
	compat.Prefix2Workers = 8
	compat.Prefix2WallElapsed = 9 * time.Second
	compat.Prefix2WorkerElapsed = 47 * time.Second
	compat.Gen3Workers = 8
	compat.Gen4Workers = 8
	compat.Gen3Tasks = 275
	compat.Gen4Tasks = 593
	compat.Gen3WallElapsed = 10 * time.Second
	compat.Gen3WorkerElapsed = 80 * time.Second
	compat.Gen4WallElapsed = 20 * time.Second
	compat.Gen4WorkerElapsed = 160 * time.Second
	g4 := &recovery.Build78Gen4Profile
	g4.Tasks = 593
	g4.Workers = 8
	g4.SingleCalls = 593
	g4.SingleAccepted = 491
	g4.PairCalls = 491
	g4.PairOutputs = 492
	g4.ContinueCalls = 492
	g4.ContinueOutputs = 937
	g4.SingleEvaluations = 9488
	g4.PairEvaluations = 102400
	g4.ContinueEvaluations = 84000
	g4.SingleWorkerElapsed = 30 * time.Second
	g4.PairWorkerElapsed = 70 * time.Second
	g4.ContinueWorkerElapsed = 60 * time.Second
	g4.DominantTaskIndex = 12
	g4.DominantTaskEvaluations = 1248
	g4.DominantTaskBank = 8
	g4.DominantSingleElapsed = 100 * time.Millisecond
	g4.DominantPairElapsed = 300 * time.Millisecond
	g4.DominantContinueElapsed = 500 * time.Millisecond
	g4.DominantSingleEvaluations = 16
	g4.DominantPairEvaluations = 512
	g4.DominantContinueEvaluations = 720
	c4 := &recovery.Build78Continuation4Profile
	c4.Calls = 492
	c4.InputStates = 492
	c4.ProbeAttempts = 90000
	c4.ScoreEvaluations = 84000
	c4.NonImprovingScores = 70000
	c4.ImprovingScoreProbes = 14000
	c4.AcceptedStates = 937
	c4.Passes = 6000
	c4.WorkerElapsed = 60 * time.Second
	c4.ScoreWorkerElapsed = 55 * time.Second
	c4.PrepareWorkerElapsed = 5 * time.Second
	c4.DominantCallIndex = 11
	c4.DominantGen4Task = 12
	c4.DominantPairRank = 1
	c4.DominantEvaluations = 128
	c4.DominantAcceptedStates = 6
	c4.DominantPasses = 8
	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild78ApplyTelemetry(&public, recovery)
	if !public.Build78Attempted || !public.Build76Attempted || !public.Build75Attempted || !public.Build73Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build78 did not preserve inherited attempted telemetry")
	}
	if public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build78Gen4Tasks != 593 || public.Build78Single4Calls != 593 || public.Build78Continue4Outputs != 937 || public.Build78DominantGen4Task != 12 || public.Build78DominantContinue4Evals != 720 {
		t.Fatalf("Build78 telemetry mismatch: %+v", public)
	}
	if public.Build78Continue4InputStates != 492 || public.Build78Continue4ScoreEvals != 84000 || public.Build78Continue4AcceptedStates != 937 || public.Build78DominantContinueCall != 11 || public.Build78DominantContinueEvals != 128 {
		t.Fatalf("Build78 continuation telemetry mismatch: %+v", public)
	}
}
