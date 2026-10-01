package watermark

import (
	"image"
	"runtime"
	"testing"
	"time"
)

func TestExperimentalV4Build75ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build75 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build75 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build75BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, geometry, freeze := experimentalV4PhoneBuild75BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if geometry != (experimentalV4PhoneBuild73GeometryTelemetry{}) || freeze != (experimentalV4PhoneBuild75FreezeTelemetry{}) {
		t.Fatalf("invalid boundary produced telemetry geometry=%+v freeze=%+v", geometry, freeze)
	}
}

func TestExperimentalV4Build75ParallelFreezeMatchesBuild47Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	boundary := PrintBoundaryEstimate{
		Detected:    true,
		Confidence:  1,
		TopLeft:     ImagePoint{X: 0, Y: 0},
		TopRight:    ImagePoint{X: float64(w - 1), Y: 0},
		BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)},
		BottomLeft:  ImagePoint{X: 0, Y: float64(h - 1)},
	}
	want, wantScores, wantProposals := experimentalV4PhoneBuild47Freeze(carrier, boundary, w, h, 128)
	got, gotScores, gotProposals, profile := experimentalV4PhoneBuild75FreezeParallel(carrier, boundary, w, h, 128)
	if gotProposals != wantProposals || len(got) != len(want) || len(gotScores) != len(wantScores) {
		t.Fatalf("freeze differs proposals=%d/%d frozen=%d/%d scores=%d/%d", gotProposals, wantProposals, len(got), len(want), len(gotScores), len(wantScores))
	}
	for i := range gotScores {
		if gotScores[i] != wantScores[i] {
			t.Fatalf("pair score differs at %d: %+v / %+v", i, gotScores[i], wantScores[i])
		}
	}
	for i := range got {
		g, w := got[i], want[i]
		if g.tier != w.tier || g.cellRank != w.cellRank || g.cellMean != w.cellMean || g.cellRobust != w.cellRobust || g.h.quad != w.h.quad || g.h.h != w.h.h || g.h.proposal != w.h.proposal || g.h.build43Pair != w.h.build43Pair || g.h.build43PairRank != w.h.build43PairRank {
			t.Fatalf("frozen candidate differs at %d", i)
		}
	}
	if profile.BasinTasks == 0 || profile.BasinWorkers == 0 || profile.BasinWallElapsed <= 0 || profile.BasinWorkerElapsed <= 0 || profile.BasinEvaluations == 0 {
		t.Fatalf("incomplete parallel basin telemetry: %+v", profile)
	}
	wantWorkers := runtime.GOMAXPROCS(0)
	if wantWorkers > profile.BasinTasks {
		wantWorkers = profile.BasinTasks
	}
	if wantWorkers < 1 {
		wantWorkers = 1
	}
	if profile.BasinWorkers != wantWorkers {
		t.Fatalf("basin workers=%d want %d", profile.BasinWorkers, wantWorkers)
	}
	if profile.BasinTasks != profile.ProductionBasinTasks+profile.DepthBasinTasks+profile.AllPairsBasinTasks || profile.BasinEvaluations != profile.ProductionBasinEvaluations+profile.DepthBasinEvaluations+profile.AllPairsBasinEvaluations {
		t.Fatalf("tier accounting mismatch: %+v", profile)
	}
}

func TestExperimentalV4Build75ApplyTelemetryPreservesQualifiedView(t *testing.T) {
	recovery := experimentalV4PhoneBuild75RecoveryTelemetry{}
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
	recovery.TotalElapsed = 240 * time.Second
	recovery.GeometryElapsed = 110 * time.Second
	gp := &recovery.GeometryProfile
	gp.Prefix2Workers = 8
	gp.Gen3Workers = 8
	gp.Gen4Workers = 8
	gp.Gen3Tasks = 275
	gp.Gen4Tasks = 593
	gp.FreezeElapsed = 22 * time.Second
	gp.Prefix2WallElapsed = 24 * time.Second
	gp.Gen3WallElapsed = 12 * time.Second
	gp.Gen4WallElapsed = 20 * time.Second
	fp := &recovery.FreezeProfile
	fp.TotalElapsed = 22 * time.Second
	fp.PairScoreElapsed = 12 * time.Second
	fp.CellsElapsed = 700 * time.Millisecond
	fp.BasinWorkers = 8
	fp.BasinTasks = 16
	fp.BasinWallElapsed = 7 * time.Second
	fp.BasinWorkerElapsed = 44 * time.Second
	fp.ProductionBasinTasks = 4
	fp.DepthBasinTasks = 4
	fp.AllPairsBasinTasks = 8
	fp.BasinEvaluations = 10000
	fp.ProductionBasinEvaluations = 2000
	fp.DepthBasinEvaluations = 2000
	fp.AllPairsBasinEvaluations = 6000
	fp.BasinMaxElapsed = 4 * time.Second

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild75ApplyTelemetry(&public, recovery)
	if !public.Build75Attempted || !public.Build73Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build75 did not preserve inherited attempted telemetry")
	}
	if public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build75FreezeTotalMs != 22000 || public.Build75BasinWorkers != 8 || public.Build75BasinTasks != 16 || public.Build75BasinWallMs != 7000 || public.Build75BasinWorkerMs != 44000 || public.Build75BasinMaxMs != 4000 {
		t.Fatalf("Build75 basin telemetry mismatch: %+v", public)
	}
}

func TestExperimentalV4Build75BlindBankMatchesBuild73Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	boundary := PrintBoundaryEstimate{
		Detected:    true,
		Confidence:  1,
		TopLeft:     ImagePoint{X: 0, Y: 0},
		TopRight:    ImagePoint{X: float64(w - 1), Y: 0},
		BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)},
		BottomLeft:  ImagePoint{X: 0, Y: float64(h - 1)},
	}
	want, wantEvals, wantSeeds, wantWorkers, _ := experimentalV4PhoneBuild73BlindBank(carrier, boundary, w, h)
	got, gotEvals, gotSeeds, gotWorkers, _, _ := experimentalV4PhoneBuild75BlindBank(carrier, boundary, w, h)
	if gotEvals != wantEvals || gotSeeds != wantSeeds || gotWorkers != wantWorkers || len(got) != len(want) {
		t.Fatalf("blind bank differs evals=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", gotEvals, wantEvals, gotSeeds, wantSeeds, gotWorkers, wantWorkers, len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || got[i].proposal != want[i].proposal {
			t.Fatalf("blind bank differs at %d", i)
		}
	}
}
