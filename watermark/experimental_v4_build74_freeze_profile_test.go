package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build74ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build74 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build74 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build74BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, geometry, freeze := experimentalV4PhoneBuild74BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if geometry != (experimentalV4PhoneBuild73GeometryTelemetry{}) || freeze != (experimentalV4PhoneBuild74FreezeTelemetry{}) {
		t.Fatalf("invalid boundary produced telemetry geometry=%+v freeze=%+v", geometry, freeze)
	}
}

func TestExperimentalV4Build74FreezeMatchesBuild47Exactly(t *testing.T) {
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
	got, gotScores, gotProposals, profile := experimentalV4PhoneBuild74FreezeProfile(carrier, boundary, w, h, 128)
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
	if profile.PairScoreTasks != len(experimentalV4PhoneBuild43Pairs) || profile.CellTasks == 0 || profile.BasinTasks == 0 || profile.PairScoreEvaluations == 0 || profile.CellEvaluations == 0 || profile.BasinEvaluations == 0 || profile.TotalElapsed <= 0 {
		t.Fatalf("incomplete freeze profile: %+v", profile)
	}
	if profile.BasinTasks != profile.ProductionBasinTasks+profile.DepthBasinTasks+profile.AllPairsBasinTasks || profile.BasinEvaluations != profile.ProductionBasinEvaluations+profile.DepthBasinEvaluations+profile.AllPairsBasinEvaluations {
		t.Fatalf("tier accounting mismatch: %+v", profile)
	}
}

func TestExperimentalV4Build74ApplyTelemetryPreservesQualifiedView(t *testing.T) {
	recovery := experimentalV4PhoneBuild74RecoveryTelemetry{}
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
	recovery.GeometryElapsed = 120 * time.Second
	gp := &recovery.GeometryProfile
	gp.Prefix2Workers = 8
	gp.Gen3Workers = 8
	gp.Gen4Workers = 8
	gp.Gen3Tasks = 275
	gp.Gen4Tasks = 593
	gp.FreezeElapsed = 60 * time.Second
	gp.Prefix2WallElapsed = 22 * time.Second
	gp.Gen3WallElapsed = 11 * time.Second
	gp.Gen4WallElapsed = 21 * time.Second
	fp := &recovery.FreezeProfile
	fp.TotalElapsed = 60 * time.Second
	fp.StructuralElapsed = 2 * time.Second
	fp.PlanePrepElapsed = 200 * time.Millisecond
	fp.PairScoreElapsed = 12 * time.Second
	fp.CellsElapsed = 3 * time.Second
	fp.BasinElapsed = 42 * time.Second
	fp.PairScoreEvaluations = 12000
	fp.CellEvaluations = 500
	fp.BasinEvaluations = 10000
	fp.BasinTasks = 16
	fp.ProductionBasinTasks = 4
	fp.DepthBasinTasks = 4
	fp.AllPairsBasinTasks = 8
	fp.BasinMaxElapsed = 4 * time.Second

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild74ApplyTelemetry(&public, recovery)
	if !public.Build74Attempted || !public.Build73Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build74 did not preserve inherited attempted telemetry")
	}
	if public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build74FreezeTotalMs != 60000 || public.Build74FreezePairScoreMs != 12000 || public.Build74FreezeBasinMs != 42000 || public.Build74FreezeBasinTasks != 16 || public.Build74FreezeBasinMaxMs != 4000 {
		t.Fatalf("Build74 freeze telemetry mismatch: %+v", public)
	}
}
