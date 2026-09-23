package watermark

import "testing"

func TestExperimentalV4Build50Top4SelectionIncludesRanksThreeAndFour(t *testing.T) {
	mk := func(pairRank int, proposal, validation float64) experimentalV4PhoneBuild47Frozen {
		return experimentalV4PhoneBuild47Frozen{h: experimentalV4PhoneHypothesis{
			proposal: proposal, validation: validation, build43Pair: "test", build43PairRank: pairRank,
		}}
	}
	frozen := []experimentalV4PhoneBuild47Frozen{
		mk(1, 0.50, 0.01), mk(1, 0.40, 0.99), mk(1, 0.30, 0.50), mk(1, 0.20, 0.90), mk(1, 0.10, 0.95),
		mk(2, 0.45, 0.01), mk(2, 0.35, 0.99), mk(2, 0.25, 0.50), mk(2, 0.15, 0.90), mk(2, 0.05, 0.95),
	}
	got := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild50SeedsPerPair)
	if len(got) != 8 {
		t.Fatalf("selected=%d want 8", len(got))
	}
	wantIdx := []int{0, 1, 2, 3, 5, 6, 7, 8}
	wantRanks := []int{1, 2, 3, 4, 1, 2, 3, 4}
	for i := range wantIdx {
		if got[i].index != wantIdx[i] || got[i].rankWithinPair != wantRanks[i] {
			t.Fatalf("seed[%d]=index %d rank %d; want index %d rank %d", i, got[i].index, got[i].rankWithinPair, wantIdx[i], wantRanks[i])
		}
	}
}

func TestExperimentalV4Build50ProductionAndBuild48RemainFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild48SeedsPerPair != 2 {
		t.Fatalf("Build48 historical diagnostic changed: seeds per pair=%d want 2", experimentalV4PhoneBuild48SeedsPerPair)
	}
	if experimentalV4PhoneBuild50SeedsPerPair != 4 {
		t.Fatalf("Build50 diagnostic seeds per pair=%d want 4", experimentalV4PhoneBuild50SeedsPerPair)
	}
}
