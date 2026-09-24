package watermark

import "testing"

func TestExperimentalV4Build61ThirdPairSiblingStencilIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild55Step != 1 {
		t.Fatalf("Build61 sibling step=%v want 1", experimentalV4PhoneBuild55Step)
	}
	if experimentalV4PhoneBuild55MaxSiblingEvals != 16 {
		t.Fatalf("Build61 sibling eval budget=%d want 16", experimentalV4PhoneBuild55MaxSiblingEvals)
	}
	if experimentalV4PhoneBuild55Dimensions != 8 {
		t.Fatalf("Build61 dimensions=%d want 8", experimentalV4PhoneBuild55Dimensions)
	}
	if experimentalV4PhoneBuild55MaxPasses != 8 || experimentalV4PhoneBuild55MaxStatesBranch != 64 {
		t.Fatalf("Build61 inherited continuation bounds passes/states=%d/%d want 8/64", experimentalV4PhoneBuild55MaxPasses, experimentalV4PhoneBuild55MaxStatesBranch)
	}
}

func TestExperimentalV4Build61UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild61SeedsPerPair != 4 || experimentalV4PhoneBuild60SeedsPerPair != 4 {
		t.Fatalf("Build60/61 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild60SeedsPerPair, experimentalV4PhoneBuild61SeedsPerPair)
	}
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
}

func TestExperimentalV4Build61RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild61Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image accepted")
	}
}
