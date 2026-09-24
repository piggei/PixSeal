package watermark

import "testing"

func TestExperimentalV4Build63FourthPairContinuationIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild55Step != 1 {
		t.Fatalf("Build63 continuation step=%v want 1", experimentalV4PhoneBuild55Step)
	}
	if experimentalV4PhoneBuild55Dimensions != 8 || experimentalV4PhoneBuild55MaxPasses != 8 || experimentalV4PhoneBuild55MaxStatesBranch != 64 {
		t.Fatalf("Build63 continuation bounds dims/passes/states=%d/%d/%d want 8/8/64", experimentalV4PhoneBuild55Dimensions, experimentalV4PhoneBuild55MaxPasses, experimentalV4PhoneBuild55MaxStatesBranch)
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build63 inherited pair bounds combos/keep=%d/%d want 112/8", experimentalV4PhoneBuild53PairCombosPerRoot, experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build63UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 || experimentalV4PhoneBuild62SeedsPerPair != 4 {
		t.Fatalf("Build62/63 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild62SeedsPerPair, experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build63RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild63Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image unexpectedly accepted")
	}
}
