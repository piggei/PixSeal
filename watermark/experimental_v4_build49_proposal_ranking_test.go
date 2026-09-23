package watermark

import "testing"

func TestExperimentalV4Build49AssignRanks(t *testing.T) {
	c := []ExperimentalV4PhoneBuild49Candidate{
		{SourcePairRank: 1, ProposalScore: .2, FoldMinScore: .4, BalancedFoldScore: .1, TileConsistencyScore: .3},
		{SourcePairRank: 1, ProposalScore: .5, FoldMinScore: .1, BalancedFoldScore: .4, TileConsistencyScore: .2},
		{SourcePairRank: 1, ProposalScore: .3, FoldMinScore: .3, BalancedFoldScore: .2, TileConsistencyScore: .5},
		{SourcePairRank: 2, ProposalScore: .4, FoldMinScore: .2, BalancedFoldScore: .2, TileConsistencyScore: .1},
	}
	experimentalV4PhoneBuild49AssignRanks(c)
	if c[1].ProposalRankGlobal != 1 || c[3].ProposalRankGlobal != 2 {
		t.Fatalf("unexpected global proposal ranks: %+v", c)
	}
	if c[0].FoldMinRankWithinPair != 1 || c[1].BalancedRankWithinPair != 1 || c[2].TileRankWithinPair != 1 {
		t.Fatalf("unexpected within-pair ranks: %+v", c[:3])
	}
}

func TestExperimentalV4Build49ProductionRemainsFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
}
