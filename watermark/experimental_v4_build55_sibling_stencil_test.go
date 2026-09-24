package watermark

import "testing"

func TestExperimentalV4Build55SiblingStencilIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild55Step != 1 {
		t.Fatalf("Build55 step=%v want 1", experimentalV4PhoneBuild55Step)
	}
	if experimentalV4PhoneBuild55Dimensions != 8 || experimentalV4PhoneBuild55MaxPasses != 8 {
		t.Fatalf("Build55 dimensions/passes=%d/%d want 8/8", experimentalV4PhoneBuild55Dimensions, experimentalV4PhoneBuild55MaxPasses)
	}
	if experimentalV4PhoneBuild55MaxSiblingEvals != 16 {
		t.Fatalf("Build55 sibling eval bound=%d want 16", experimentalV4PhoneBuild55MaxSiblingEvals)
	}
	if experimentalV4PhoneBuild55MaxStatesBranch != 64 {
		t.Fatalf("Build55 continuation bound=%d want 64", experimentalV4PhoneBuild55MaxStatesBranch)
	}
}

func TestExperimentalV4Build55UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild55SeedsPerPair != 4 || experimentalV4PhoneBuild54SeedsPerPair != 4 {
		t.Fatalf("Build54/55 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild54SeedsPerPair, experimentalV4PhoneBuild55SeedsPerPair)
	}
}

func TestExperimentalV4Build55RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild55Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image unexpectedly accepted")
	}
}
