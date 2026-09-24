package watermark

import "testing"

func TestExperimentalV4Build62FourthPairEscapeIsFixedAndBounded(t *testing.T) {
	if experimentalV4PhoneBuild53ProbeStep != 1 {
		t.Fatalf("Build62 pair step=%v want 1", experimentalV4PhoneBuild53ProbeStep)
	}
	if experimentalV4PhoneBuild53Dimensions != 8 {
		t.Fatalf("Build62 dimensions=%d want 8", experimentalV4PhoneBuild53Dimensions)
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 {
		t.Fatalf("Build62 pair eval budget=%d want 112", experimentalV4PhoneBuild53PairCombosPerRoot)
	}
	if experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build62 retained pair budget=%d want 8", experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build62UsesSameTop4AndLeavesProductionFrozen(t *testing.T) {
	if experimentalV4PhoneBuild62SeedsPerPair != 4 || experimentalV4PhoneBuild61SeedsPerPair != 4 {
		t.Fatalf("Build61/62 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild61SeedsPerPair, experimentalV4PhoneBuild62SeedsPerPair)
	}
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
}

func TestExperimentalV4Build62RejectsInvalidInputs(t *testing.T) {
	if _, err := ExperimentalV4PhoneBuild62Diagnose(nil, []byte("12345678"), 1632, 1632); err == nil {
		t.Fatal("nil image accepted")
	}
}
