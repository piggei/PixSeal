package watermark

import "testing"

func TestExperimentalV4Build52FineRestartScheduleIsFixedAndBounded(t *testing.T) {
	want := []float64{2, 1}
	if len(experimentalV4PhoneBuild52FineSteps) != len(want) {
		t.Fatalf("fine steps=%v want %v", experimentalV4PhoneBuild52FineSteps, want)
	}
	for i := range want {
		if experimentalV4PhoneBuild52FineSteps[i] != want[i] {
			t.Fatalf("fine step[%d]=%v want %v", i, experimentalV4PhoneBuild52FineSteps[i], want[i])
		}
	}
}

func TestExperimentalV4Build52UsesTop4WithoutChangingProductionOrHistoricalDepths(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild48SeedsPerPair != 2 {
		t.Fatalf("Build48 historical depth=%d want 2", experimentalV4PhoneBuild48SeedsPerPair)
	}
	if experimentalV4PhoneBuild50SeedsPerPair != 4 || experimentalV4PhoneBuild51SeedsPerPair != 4 || experimentalV4PhoneBuild52SeedsPerPair != 4 {
		t.Fatalf("Build50/51/52 diagnostic depth=%d/%d/%d want 4/4/4", experimentalV4PhoneBuild50SeedsPerPair, experimentalV4PhoneBuild51SeedsPerPair, experimentalV4PhoneBuild52SeedsPerPair)
	}
}

func TestExperimentalV4Build52RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild52Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image unexpectedly accepted")
	}
}
