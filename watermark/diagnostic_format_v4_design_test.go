package watermark

import "testing"

func TestDiagnosticFormatV4DesignStudyKeepsV3FrozenAndQuantifiesTradeoffs(t *testing.T) {
	study := diagnosticFormatV4DesignStudy()
	if study.Method != "v4-absolute-pilot-layout-trade-study" || len(study.Candidates) != 3 {
		t.Fatalf("unexpected study: %+v", study)
	}
	var compact, preserve *DiagnosticFormatV4DesignCandidate
	for i := range study.Candidates {
		candidate := &study.Candidates[i]
		switch candidate.Name {
		case "compact-35x32-p64":
			compact = candidate
		case "preserve-capacity-37x32-p64":
			preserve = candidate
		}
	}
	if compact == nil || compact.PreservesV3MaxPayload || compact.ComparableMaxPayloadBytes != 59 {
		t.Fatalf("unexpected compact candidate: %+v", compact)
	}
	if preserve == nil || !preserve.PreservesV3MaxPayload || preserve.DataPositions != eccBits || preserve.ComparableMaxPayloadBytes != maxPayload {
		t.Fatalf("unexpected preserve candidate: %+v", preserve)
	}
	if preserve.MinimumWidthPixels != 296 || preserve.MinimumHeightPixels != 256 {
		t.Fatalf("unexpected preserve minimum geometry: %dx%d", preserve.MinimumWidthPixels, preserve.MinimumHeightPixels)
	}
}
