package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build77ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build77 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build77 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build77BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, gp, fp, g4 := experimentalV4PhoneBuild77BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if gp != (experimentalV4PhoneBuild76GeometryTelemetry{}) || fp != (experimentalV4PhoneBuild75FreezeTelemetry{}) || g4 != (experimentalV4PhoneBuild77Gen4Profile{}) {
		t.Fatalf("invalid boundary produced telemetry")
	}
}

func TestExperimentalV4Build77BlindBankMatchesBuild76Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, wantEvals, wantSeeds, wantWorkers, _, _ := experimentalV4PhoneBuild76BlindBank(carrier, boundary, w, h)
	got, gotEvals, gotSeeds, gotWorkers, gp, fp, g4 := experimentalV4PhoneBuild77BlindBank(carrier, boundary, w, h)
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
		t.Fatalf("Build77 task accounting mismatch g4=%+v gp=%+v", g4, gp)
	}
	if gp.Gen4Tasks > 0 {
		if g4.SingleEvaluations+g4.PairEvaluations+g4.ContinueEvaluations <= 0 {
			t.Fatalf("Build77 stage evaluations not recorded: %+v", g4)
		}
		if g4.SingleWorkerElapsed <= 0 || g4.PairWorkerElapsed < 0 || g4.ContinueWorkerElapsed < 0 {
			t.Fatalf("Build77 stage timing not recorded: %+v", g4)
		}
	}
}

func TestExperimentalV4Build77ApplyTelemetryPreservesQualifiedView(t *testing.T) {
	recovery := experimentalV4PhoneBuild77RecoveryTelemetry{}
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
	g4 := &recovery.Build77Gen4Profile
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
	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild77ApplyTelemetry(&public, recovery)
	if !public.Build77Attempted || !public.Build76Attempted || !public.Build75Attempted || !public.Build73Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build77 did not preserve inherited attempted telemetry")
	}
	if public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build77Gen4Tasks != 593 || public.Build77Single4Calls != 593 || public.Build77Continue4Outputs != 937 || public.Build77DominantGen4Task != 12 || public.Build77DominantContinue4Evals != 720 {
		t.Fatalf("Build77 telemetry mismatch: %+v", public)
	}
}
