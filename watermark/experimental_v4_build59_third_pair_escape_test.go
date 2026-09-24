package watermark

import "testing"

func TestExperimentalV4Build59ThirdPairEscapeIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild55Step != 1 || experimentalV4PhoneBuild53ProbeStep != 1 {
		t.Fatalf("Build59 inherited sibling/pair steps=%v/%v want 1/1", experimentalV4PhoneBuild55Step, experimentalV4PhoneBuild53ProbeStep)
	}
	if experimentalV4PhoneBuild55MaxSiblingEvals != 16 {
		t.Fatalf("Build59 local sibling probe bound=%d want 16", experimentalV4PhoneBuild55MaxSiblingEvals)
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build59 pair bounds combos/keep=%d/%d want 112/8", experimentalV4PhoneBuild53PairCombosPerRoot, experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build59UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild59SeedsPerPair != 4 || experimentalV4PhoneBuild58SeedsPerPair != 4 {
		t.Fatalf("Build58/59 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild58SeedsPerPair, experimentalV4PhoneBuild59SeedsPerPair)
	}
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
}

func TestExperimentalV4Build59RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild59Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image accepted")
	}
}
