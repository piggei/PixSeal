package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build79ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build79 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build79 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
	if experimentalV4PhoneBuild79KernelSampleMod != 64 {
		t.Fatalf("Build79 kernel sample modulus=%d want 64", experimentalV4PhoneBuild79KernelSampleMod)
	}
}

func TestExperimentalV4Build79BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, gp, fp, g4, c4, kp := experimentalV4PhoneBuild79BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if gp != (experimentalV4PhoneBuild76GeometryTelemetry{}) || fp != (experimentalV4PhoneBuild75FreezeTelemetry{}) || g4 != (experimentalV4PhoneBuild79Gen4Profile{}) || c4 != (experimentalV4PhoneBuild79Continuation4Profile{}) || kp != (experimentalV4PhoneBuild79FoldKernelProfile{}) {
		t.Fatalf("invalid boundary produced telemetry")
	}
}

func build79SyntheticFixture(t *testing.T) (experimentalV4PilotCandidate, image.Image, *pixelPlane, int, int, [4]ImagePoint, homography) {
	t.Helper()
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	plane := newPixelPlane(carrier)
	anchor := [4]ImagePoint{{X: 0, Y: 0}, {X: float64(w - 1), Y: 0}, {X: 0, Y: float64(h - 1)}, {X: float64(w - 1), Y: float64(h - 1)}}
	hom, ok := experimentalV4PhoneQuadHomography(w, h, anchor)
	if !ok {
		t.Fatal("anchor homography")
	}
	return candidate, carrier, plane, w, h, anchor, hom
}

func TestExperimentalV4Build79FoldScoreProfiledMatchesBuild41Exactly(t *testing.T) {
	candidate, _, plane, w, h, _, hom := build79SyntheticFixture(t)
	wantScore, wantVisible := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	for _, sampled := range []bool{false, true} {
		gotScore, gotVisible, profile := experimentalV4PhoneBuild79FoldScoreProfiled(plane, candidate, w, h, hom, 0, true, false, sampled)
		if gotScore != wantScore || gotVisible != wantVisible {
			t.Fatalf("sampled=%t score/visible differ got=(%.17g,%d) want=(%.17g,%d)", sampled, gotScore, gotVisible, wantScore, wantVisible)
		}
		if profile.blockReads == 0 || profile.blockSuccess == 0 || profile.visible != wantVisible || profile.blockSuccess+profile.blockFailed != profile.blockReads {
			t.Fatalf("sampled=%t incomplete fold accounting: %+v", sampled, profile)
		}
		if sampled {
			if profile.sampledBlockReadElapsed <= 0 || profile.detailedBlocks != 1 || profile.detailedPixels != blockSize*blockSize || profile.replayFailures != 0 {
				t.Fatalf("sampled kernel decomposition missing: %+v", profile)
			}
		}
	}
}

func TestExperimentalV4Build79ContinueProfiledMatchesBuild55Exactly(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build79SyntheticFixture(t)
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, wantEvals := experimentalV4PhoneBuild55Continue(plane, candidate, w, h, anchor, start, 0)
	got, gotEvals, profile := experimentalV4PhoneBuild79ContinueProfiled(plane, candidate, w, h, anchor, start, 0, 0, 0)
	if gotEvals != wantEvals || len(got) != len(want) {
		t.Fatalf("continuation differs evals=%d/%d states=%d/%d", gotEvals, wantEvals, len(got), len(want))
	}
	for i := range got {
		if got[i].index != want[i].index || got[i].pass != want[i].pass || got[i].dim != want[i].dim || got[i].delta != want[i].delta || got[i].hyp.quad != want[i].hyp.quad || got[i].hyp.h != want[i].hyp.h || got[i].hyp.proposal != want[i].hyp.proposal {
			t.Fatalf("continuation differs at %d", i)
		}
	}
	if profile.evals != wantEvals || len(profile.foldProfiles) != wantEvals || profile.acceptedStates != len(want) || profile.probeAttempts < wantEvals {
		t.Fatalf("profile accounting mismatch: %+v", profile)
	}
}

func TestExperimentalV4Build79BlindBankMatchesBuild76Exactly(t *testing.T) {
	candidate, carrier, _, w, h, _, _ := build79SyntheticFixture(t)
	_ = candidate
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, wantEvals, wantSeeds, wantWorkers, _, _ := experimentalV4PhoneBuild76BlindBank(carrier, boundary, w, h)
	got, gotEvals, gotSeeds, gotWorkers, gp, fp, g4, c4, kp := experimentalV4PhoneBuild79BlindBank(carrier, boundary, w, h)
	if gotEvals != wantEvals || gotSeeds != wantSeeds || gotWorkers != wantWorkers || len(got) != len(want) {
		t.Fatalf("blind bank differs evals=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", gotEvals, wantEvals, gotSeeds, wantSeeds, gotWorkers, wantWorkers, len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || got[i].proposal != want[i].proposal || got[i].validation != want[i].validation {
			t.Fatalf("blind bank differs at %d", i)
		}
	}
	if gp.Gen2Tasks == 0 || gp.Gen2Workers == 0 || gp.Gen3Tasks == 0 || fp.BasinTasks != 16 || fp.BasinWorkers == 0 {
		t.Fatalf("inherited telemetry incomplete gp=%+v fp=%+v", gp, fp)
	}
	if g4.Tasks != gp.Gen4Tasks || c4.ScoreEvaluations != g4.ContinueEvaluations || kp.FoldCalls != c4.ScoreEvaluations {
		t.Fatalf("Build79 kernel accounting mismatch g4=%+v c4=%+v kp=%+v", g4, c4, kp)
	}
	if gp.Gen4Tasks > 0 && (kp.BlockReads == 0 || kp.BlockSuccess == 0 || kp.SampledFoldCalls == 0 || kp.DetailedBlocks == 0) {
		t.Fatalf("Build79 kernel profiler missing deep-work telemetry: %+v", kp)
	}
}

func TestExperimentalV4Build79ApplyTelemetryPreservesQualifiedView(t *testing.T) {
	recovery := experimentalV4PhoneBuild79RecoveryTelemetry{}
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
	gp := &recovery.Build76GeometryProfile
	gp.Prefix1Workers, gp.Gen2Workers, gp.Gen3Workers, gp.Gen4Workers = 8, 8, 8, 8
	gp.Gen2Tasks, gp.Gen3Tasks, gp.Gen4Tasks = 93, 275, 593
	gp.FreezeElapsed = 23 * time.Second
	gp.Prefix1WallElapsed, gp.Gen2WallElapsed, gp.Gen3WallElapsed, gp.Gen4WallElapsed = 4*time.Second, 5*time.Second, 10*time.Second, 20*time.Second
	gp.Prefix1WorkerElapsed, gp.Gen2WorkerElapsed, gp.Gen3WorkerElapsed, gp.Gen4WorkerElapsed = 20*time.Second, 27*time.Second, 80*time.Second, 160*time.Second
	compat := &recovery.GeometryProfile
	compat.Prefix2Workers, compat.Gen3Workers, compat.Gen4Workers = 8, 8, 8
	compat.Prefix2WallElapsed, compat.Prefix2WorkerElapsed = 9*time.Second, 47*time.Second
	compat.Gen3Tasks, compat.Gen4Tasks = 275, 593
	compat.Gen3WallElapsed, compat.Gen3WorkerElapsed = 10*time.Second, 80*time.Second
	compat.Gen4WallElapsed, compat.Gen4WorkerElapsed = 20*time.Second, 160*time.Second
	g4 := &recovery.Build79Gen4Profile
	g4.Tasks, g4.Workers, g4.SingleCalls, g4.ContinueCalls, g4.ContinueOutputs = 593, 8, 593, 492, 937
	g4.ContinueEvaluations = 84000
	c4 := &recovery.Build79Continuation4Profile
	c4.Calls, c4.InputStates, c4.ScoreEvaluations, c4.AcceptedStates = 492, 492, 84000, 937
	kp := &recovery.Build79FoldKernelProfile
	kp.FoldCalls = 84000
	kp.TilesVisited = 168000
	kp.PilotPositions = 10752000
	kp.BlockReads = 10752000
	kp.BlockSuccess = 10752000
	kp.Visible = 10752000
	kp.FoldWorkerElapsed = 55 * time.Second
	kp.SampledFoldCalls = 1300
	kp.SampledFoldElapsed = time.Second
	kp.SampledBlockReadElapsed = 950 * time.Millisecond
	kp.SampledFoldOverhead = 50 * time.Millisecond
	kp.DetailedBlocks = 1300
	kp.DetailedPixels = 1300 * blockSize * blockSize
	kp.ReplayMapElapsed = 10 * time.Millisecond
	kp.ReplaySampleElapsed = 40 * time.Millisecond
	kp.ReplayDCTElapsed = 5 * time.Millisecond
	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild79ApplyTelemetry(&public, recovery)
	if !public.Build79Attempted || !public.Build76Attempted || !public.Build75Attempted || !public.Build73Attempted || !public.Build71Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build64Attempted {
		t.Fatalf("Build79 did not preserve inherited attempted telemetry")
	}
	if public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build79FoldCalls != 84000 || public.Build79FoldBlockReads != 10752000 || public.Build79SampledFoldCalls != 1300 || public.Build79KernelSampleMod != 64 || public.Build79ReplaySampleMs != 40 {
		t.Fatalf("Build79 kernel telemetry mismatch: %+v", public)
	}
}
