package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build76ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build76 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build76 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build76BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, gp, fp := experimentalV4PhoneBuild76BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if gp != (experimentalV4PhoneBuild76GeometryTelemetry{}) || fp != (experimentalV4PhoneBuild75FreezeTelemetry{}) {
		t.Fatalf("invalid boundary produced telemetry")
	}
}

func TestExperimentalV4Build76BlindBankMatchesBuild75Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, wantEvals, wantSeeds, wantWorkers, _, _ := experimentalV4PhoneBuild75BlindBank(carrier, boundary, w, h)
	got, gotEvals, gotSeeds, gotWorkers, gp, fp := experimentalV4PhoneBuild76BlindBank(carrier, boundary, w, h)
	if gotEvals != wantEvals || gotSeeds != wantSeeds || gotWorkers != wantWorkers || len(got) != len(want) {
		t.Fatalf("blind bank differs evals=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", gotEvals, wantEvals, gotSeeds, wantSeeds, gotWorkers, wantWorkers, len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || got[i].proposal != want[i].proposal || got[i].validation != want[i].validation {
			t.Fatalf("blind bank differs at %d", i)
		}
	}
	if gp.Gen2Tasks == 0 || gp.Gen2Workers == 0 || gp.Gen2WallElapsed <= 0 || gp.Gen2WorkerElapsed <= 0 || gp.Gen3Tasks == 0 {
		t.Fatalf("incomplete Build76 geometry telemetry: %+v", gp)
	}
	if fp.BasinTasks != 16 || fp.BasinWorkers == 0 {
		t.Fatalf("Build75 freeze telemetry not preserved: %+v", fp)
	}
}

func TestExperimentalV4Build76ApplyTelemetryPreservesQualifiedView(t *testing.T) {
	recovery := experimentalV4PhoneBuild76RecoveryTelemetry{}
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
	gp.Gen2Tasks = 144
	gp.Gen3Tasks = 275
	gp.Gen4Tasks = 593
	gp.FreezeElapsed = 23 * time.Second
	gp.Prefix1WallElapsed = 4 * time.Second
	gp.Prefix1WorkerElapsed = 8 * time.Second
	gp.Gen2WallElapsed = 5 * time.Second
	gp.Gen2WorkerElapsed = 34 * time.Second
	gp.Gen3WallElapsed = 10 * time.Second
	gp.Gen3WorkerElapsed = 80 * time.Second
	gp.Gen4WallElapsed = 20 * time.Second
	gp.Gen4WorkerElapsed = 160 * time.Second
	gp.Gen2MaxElapsed = 900 * time.Millisecond
	gp.Gen2MaxEvaluations = 640
	gp.Gen2MaxOutputs = 16
	compat := &recovery.GeometryProfile
	compat.Prefix2Workers = 8
	compat.Prefix2WallElapsed = 9 * time.Second
	compat.Prefix2WorkerElapsed = 42 * time.Second
	compat.Gen3Workers = 8
	compat.Gen4Workers = 8
	compat.Gen3Tasks = 275
	compat.Gen4Tasks = 593
	compat.Gen3WallElapsed = 10 * time.Second
	compat.Gen3WorkerElapsed = 80 * time.Second
	compat.Gen4WallElapsed = 20 * time.Second
	compat.Gen4WorkerElapsed = 160 * time.Second
	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild76ApplyTelemetry(&public, recovery)
	if !public.Build76Attempted || !public.Build75Attempted || !public.Build73Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build76 did not preserve inherited attempted telemetry")
	}
	if public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build76Gen2Tasks != 144 || public.Build76Gen2Workers != 8 || public.Build76Gen2WallMs != 5000 || public.Build76Gen2WorkerMs != 34000 || public.Build76Gen2MaxMs != 900 {
		t.Fatalf("Build76 telemetry mismatch: %+v", public)
	}
}
