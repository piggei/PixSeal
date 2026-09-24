package watermark

import "testing"

func TestExperimentalV4Build58SecondPairSiblingStencilIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild55Step != 1 {
		t.Fatalf("Build58 sibling step=%v want 1", experimentalV4PhoneBuild55Step)
	}
	if experimentalV4PhoneBuild55MaxSiblingEvals != 16 {
		t.Fatalf("Build58 sibling eval budget=%d want 16", experimentalV4PhoneBuild55MaxSiblingEvals)
	}
	if experimentalV4PhoneBuild55Dimensions != 8 {
		t.Fatalf("Build58 dimensions=%d want 8", experimentalV4PhoneBuild55Dimensions)
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build58 inherited pair bounds combos/keep=%d/%d want 112/8", experimentalV4PhoneBuild53PairCombosPerRoot, experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build58UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild58SeedsPerPair != 4 || experimentalV4PhoneBuild57SeedsPerPair != 4 {
		t.Fatalf("Build57/58 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild57SeedsPerPair, experimentalV4PhoneBuild58SeedsPerPair)
	}
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
}

func TestExperimentalV4Build58RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild58Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image accepted")
	}
}
