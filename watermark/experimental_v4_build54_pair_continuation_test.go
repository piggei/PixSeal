package watermark

import "testing"

func TestExperimentalV4Build54ContinuationIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild54Step != 1 {
		t.Fatalf("Build54 step=%v want 1", experimentalV4PhoneBuild54Step)
	}
	if experimentalV4PhoneBuild54Dimensions != 8 || experimentalV4PhoneBuild54MaxPasses != 8 {
		t.Fatalf("Build54 dimensions/passes=%d/%d want 8/8", experimentalV4PhoneBuild54Dimensions, experimentalV4PhoneBuild54MaxPasses)
	}
	if experimentalV4PhoneBuild54MaxStatesBranch != 64 {
		t.Fatalf("Build54 max states/branch=%d want 64", experimentalV4PhoneBuild54MaxStatesBranch)
	}
	if experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build53 pair keep changed=%d want 8", experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build54UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild54SeedsPerPair != 4 || experimentalV4PhoneBuild53SeedsPerPair != 4 {
		t.Fatalf("Build53/54 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild53SeedsPerPair, experimentalV4PhoneBuild54SeedsPerPair)
	}
}

func TestExperimentalV4Build54RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild54Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image unexpectedly accepted")
	}
}
