package watermark

import "testing"

func TestExperimentalV4Build53PairEscapeIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild53RootStep != 2 || experimentalV4PhoneBuild53ProbeStep != 1 {
		t.Fatalf("Build53 root/probe steps=%v/%v want 2/1", experimentalV4PhoneBuild53RootStep, experimentalV4PhoneBuild53ProbeStep)
	}
	if experimentalV4PhoneBuild53MaxPasses != 2 || experimentalV4PhoneBuild53Dimensions != 8 {
		t.Fatalf("Build53 passes/dimensions=%d/%d want 2/8", experimentalV4PhoneBuild53MaxPasses, experimentalV4PhoneBuild53Dimensions)
	}
	if experimentalV4PhoneBuild53MaxRootsPerSeed != 17 {
		t.Fatalf("Build53 max roots=%d want 17", experimentalV4PhoneBuild53MaxRootsPerSeed)
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build53 pair bounds combos/keep=%d/%d want 112/8", experimentalV4PhoneBuild53PairCombosPerRoot, experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build53UsesTop4WithoutChangingProductionOrHistoricalDepths(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild48SeedsPerPair != 2 {
		t.Fatalf("Build48 historical depth=%d want 2", experimentalV4PhoneBuild48SeedsPerPair)
	}
	if experimentalV4PhoneBuild50SeedsPerPair != 4 || experimentalV4PhoneBuild51SeedsPerPair != 4 || experimentalV4PhoneBuild52SeedsPerPair != 4 || experimentalV4PhoneBuild53SeedsPerPair != 4 {
		t.Fatalf("Build50/51/52/53 diagnostic depth=%d/%d/%d/%d want 4/4/4/4", experimentalV4PhoneBuild50SeedsPerPair, experimentalV4PhoneBuild51SeedsPerPair, experimentalV4PhoneBuild52SeedsPerPair, experimentalV4PhoneBuild53SeedsPerPair)
	}
}

func TestExperimentalV4Build53RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild53Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image unexpectedly accepted")
	}
}
