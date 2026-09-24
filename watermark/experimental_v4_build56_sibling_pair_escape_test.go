package watermark

import "testing"

func TestExperimentalV4Build56SiblingPairEscapeIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild55Step != 1 || experimentalV4PhoneBuild53ProbeStep != 1 {
		t.Fatalf("Build56 inherited sibling/pair steps=%v/%v want 1/1", experimentalV4PhoneBuild55Step, experimentalV4PhoneBuild53ProbeStep)
	}
	if experimentalV4PhoneBuild55MaxSiblingEvals != 16 {
		t.Fatalf("Build56 local sibling probe bound=%d want 16", experimentalV4PhoneBuild55MaxSiblingEvals)
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build56 pair bounds combos/keep=%d/%d want 112/8", experimentalV4PhoneBuild53PairCombosPerRoot, experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build56UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild56SeedsPerPair != 4 || experimentalV4PhoneBuild55SeedsPerPair != 4 {
		t.Fatalf("Build55/56 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild55SeedsPerPair, experimentalV4PhoneBuild56SeedsPerPair)
	}
}

func TestExperimentalV4Build56RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild56Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image unexpectedly accepted")
	}
}
