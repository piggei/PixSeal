package watermark

import (
	"testing"
	"time"
)

func TestExperimentalV4Build68ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build68 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build68 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build68PlaneDetectorMatchesHistoricalDetectorExactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	plane := newPixelPlane(carrier)
	cases := []homography{
		{h: [9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}},
		{h: [9]float64{1, 0, 3.25, 0, 1, -2.5, 0, 0, 1}},
		{h: [9]float64{0.997, 0.004, 2.0, -0.003, 1.002, 1.5, 0.000002, -0.000003, 1}},
	}
	for i, h := range cases {
		got := experimentalV4DetectPilotProjectivePlane(plane, candidate, carrier.Bounds().Dx(), carrier.Bounds().Dy(), h)
		want := experimentalV4DetectPilotProjective(carrier, candidate, carrier.Bounds().Dx(), carrier.Bounds().Dy(), h)
		if got != want {
			t.Fatalf("case %d plane detector differs from historical detector\n got: %+v\nwant: %+v", i, got, want)
		}
	}
}

func TestExperimentalV4Build68QualificationMatchesBuild41Exactly(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	plane := newPixelPlane(carrier)
	h := homography{h: [9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}}
	hyp := experimentalV4PhoneHypothesis{h: h, proposal: 1}

	got, gotEvals, gotOK := experimentalV4PhoneBuild68Qualify(plane, candidate, carrier.Bounds().Dx(), carrier.Bounds().Dy(), hyp, 0)
	want, wantEvals, wantOK := experimentalV4PhoneBuild41Qualify(plane, carrier, candidate, carrier.Bounds().Dx(), carrier.Bounds().Dy(), hyp, 0)
	if gotEvals != wantEvals || gotOK != wantOK || got.validation != want.validation || got.detection != want.detection {
		t.Fatalf("Build68 qualification differs: got evals=%d ok=%t validation=%g detection=%+v; want evals=%d ok=%t validation=%g detection=%+v", gotEvals, gotOK, got.validation, got.detection, wantEvals, wantOK, want.validation, want.detection)
	}
}

func TestExperimentalV4Build68ApplyTelemetryPreservesBuild66LogicalView(t *testing.T) {
	recovery := experimentalV4PhoneBuild68RecoveryTelemetry{}
	recovery.Attempted = true
	recovery.DecodeWorkers = 8
	recovery.DecodeCandidatesTried = 691
	recovery.PhysicalDecodeCandidates = 696
	recovery.PhysicalProfilesTried = 2086
	recovery.PhysicalListFramesTried = 2135407
	recovery.TotalElapsed = 410 * time.Second
	recovery.GeometryElapsed = 288 * time.Second
	recovery.PlanePrepElapsed = 272 * time.Millisecond
	recovery.QualificationElapsed = 85 * time.Second
	recovery.DecodeWallElapsed = 36 * time.Second
	recovery.SamplingWorkerElapsed = 263 * time.Second
	recovery.ListWorkerElapsed = 7 * time.Second

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild68ApplyTelemetry(&public, recovery)
	if !public.Build68Attempted || !public.Build66Attempted {
		t.Fatalf("Build68 telemetry did not preserve inherited Build66 attempted state")
	}
	if public.Build67Attempted {
		t.Fatalf("Build68 must not masquerade as a Build67 profiling run")
	}
	if public.Build68SpeculativeDecodeCandidates != 5 {
		t.Fatalf("speculative=%d want 5", public.Build68SpeculativeDecodeCandidates)
	}
	if public.Build64DecodeCandidates != 691 {
		t.Fatalf("logical decode candidates=%d want 691", public.Build64DecodeCandidates)
	}
	if public.Build68QualificationMs != 85000 || public.Build68PlanePrepMs != 272 {
		t.Fatalf("timing qualification=%d plane=%d", public.Build68QualificationMs, public.Build68PlanePrepMs)
	}
}
