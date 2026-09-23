package watermark

import "testing"

func TestExperimentalV4Build48SeedSelectionIsProposalOnly(t *testing.T) {
	mk := func(pairRank int, proposal, validation float64) experimentalV4PhoneBuild47Frozen {
		return experimentalV4PhoneBuild47Frozen{h: experimentalV4PhoneHypothesis{
			proposal: proposal, validation: validation, build43Pair: "test", build43PairRank: pairRank,
		}}
	}
	frozen := []experimentalV4PhoneBuild47Frozen{
		mk(1, 0.20, 0.99),
		mk(1, 0.40, 0.01),
		mk(1, 0.30, 0.80),
		mk(2, 0.10, 0.90),
		mk(2, 0.25, 0.02),
		mk(2, 0.15, 0.95),
	}
	got := experimentalV4PhoneBuild48SelectSeeds(frozen, 2)
	if len(got) != 4 {
		t.Fatalf("selected=%d want 4", len(got))
	}
	want := []int{1, 2, 4, 5}
	for i := range want {
		if got[i].index != want[i] {
			t.Fatalf("selected indexes=%v want %v", []int{got[0].index, got[1].index, got[2].index, got[3].index}, want)
		}
	}
}

func TestExperimentalV4Build48ProductionRemainsFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild48SeedsPerPair != 2 {
		t.Fatalf("Build48 diagnostic seeds per pair=%d want 2", experimentalV4PhoneBuild48SeedsPerPair)
	}
}
