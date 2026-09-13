package watermark

import "testing"

func TestExperimentalV4Prototype2StructuralQualification(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	qualification := experimentalV4PilotQualification(candidate)

	if qualification.Cyclic.PerfectNonZeroCyclicAliases != 0 {
		t.Fatalf("prototype-2 has a non-zero cyclic alias: %+v", qualification.Cyclic)
	}
	if qualification.Cyclic.MaximumMaskOverlap != 7 {
		t.Fatalf("prototype-2 max overlap = %d, want 7", qualification.Cyclic.MaximumMaskOverlap)
	}
	if qualification.Cyclic.MaximumWrongSignedCorrelation != 4 {
		t.Fatalf("prototype-2 max wrong correlation = %d, want 4", qualification.Cyclic.MaximumWrongSignedCorrelation)
	}
	if qualification.Cyclic.RunnerUpWrongSignedCorrelation > qualification.Cyclic.MaximumWrongSignedCorrelation {
		t.Fatalf("runner-up exceeds maximum: %+v", qualification.Cyclic)
	}

	wantMargins := map[int]int{64: 60, 48: 43, 32: 28, 24: 20, 16: 12}
	for _, partial := range qualification.Partial {
		want, ok := wantMargins[partial.VisiblePilots]
		if !ok {
			t.Fatalf("unexpected partial visibility level: %+v", partial)
		}
		if partial.WorstContiguousMargin != want {
			t.Fatalf("visible=%d contiguous margin=%d, want %d: %+v", partial.VisiblePilots, partial.WorstContiguousMargin, want, partial)
		}
		if partial.FalseOriginCases != 0 || partial.FalseOriginRate != 0 {
			t.Fatalf("visible=%d has random-subset false origin: %+v", partial.VisiblePilots, partial)
		}
	}
}

func TestExperimentalV4Prototype2ImprovesBuild23Baseline(t *testing.T) {
	baseline := experimentalV4PilotQualification(experimentalV4Prototype1Candidate())
	candidate := experimentalV4PilotQualification(experimentalV4Prototype2Candidate())

	if candidate.Cyclic.MaximumMaskOverlap >= baseline.Cyclic.MaximumMaskOverlap {
		t.Fatalf("prototype-2 did not improve max overlap: baseline=%d candidate=%d", baseline.Cyclic.MaximumMaskOverlap, candidate.Cyclic.MaximumMaskOverlap)
	}
	if candidate.Cyclic.MaximumWrongSignedCorrelation >= baseline.Cyclic.MaximumWrongSignedCorrelation {
		t.Fatalf("prototype-2 did not improve max wrong correlation: baseline=%d candidate=%d", baseline.Cyclic.MaximumWrongSignedCorrelation, candidate.Cyclic.MaximumWrongSignedCorrelation)
	}

	baselineMargins := make(map[int]int, len(baseline.Partial))
	for _, partial := range baseline.Partial {
		baselineMargins[partial.VisiblePilots] = partial.WorstContiguousMargin
	}
	for _, partial := range candidate.Partial {
		if partial.WorstContiguousMargin < baselineMargins[partial.VisiblePilots] {
			t.Fatalf("prototype-2 regressed visible=%d crop margin: baseline=%d candidate=%d", partial.VisiblePilots, baselineMargins[partial.VisiblePilots], partial.WorstContiguousMargin)
		}
	}
}

func TestExperimentalV4Build24SearchReproducesPrototype2(t *testing.T) {
	winner, report := experimentalV4Build24Search()
	expected := experimentalV4Prototype2Candidate()

	if winner.positions != expected.positions {
		t.Fatalf("Build24 search mask is not reproducible\ngot:  %v\nwant: %v", winner.positions, expected.positions)
	}
	if winner.signs != expected.signs {
		t.Fatalf("Build24 sign refinement is not reproducible\ngot:  %v\nwant: %v", winner.signs, expected.signs)
	}
	if report.MaskSearchSeed != experimentalV4Build24MaskSearchSeed || report.MaskSearchBudget != experimentalV4Build24MaskSearchBudget {
		t.Fatalf("unexpected mask-search parameters: %+v", report)
	}
	if report.SignSearchSeed != experimentalV4Build24SignSearchSeed || report.SignSearchBudget != experimentalV4Build24SignSearchBudget {
		t.Fatalf("unexpected sign-search parameters: %+v", report)
	}
	if report.StageOneWinner.Cyclic.MaximumMaskOverlap != 7 || report.StageOneWinner.Cyclic.MaximumWrongSignedCorrelation != 5 {
		t.Fatalf("unexpected stage-one winner: %+v", report.StageOneWinner)
	}
	if report.Winner.Cyclic.MaximumMaskOverlap != 7 || report.Winner.Cyclic.MaximumWrongSignedCorrelation != 4 {
		t.Fatalf("unexpected final winner: %+v", report.Winner)
	}
	if report.Winner.CandidateHash == "" || report.RunnerUp.CandidateHash == "" || report.StageOneRunnerUp.CandidateHash == "" {
		t.Fatalf("search report did not retain winner/runner-up identifiers: %+v", report)
	}
}
