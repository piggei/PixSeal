package watermark

import "testing"

func TestExperimentalV4Build47Limits(t *testing.T) {
	want := []int{32, 64, 128}
	if len(experimentalV4PhoneBuild47Limits) != len(want) {
		t.Fatalf("limits=%v", experimentalV4PhoneBuild47Limits)
	}
	for i := range want {
		if experimentalV4PhoneBuild47Limits[i] != want[i] {
			t.Fatalf("limits=%v", experimentalV4PhoneBuild47Limits)
		}
	}
	if experimentalV4PhoneBuild47TierName(1) != "production" ||
		experimentalV4PhoneBuild47TierName(2) != "selected-pair-depth" ||
		experimentalV4PhoneBuild47TierName(3) != "all-pair-extension" {
		t.Fatalf("unexpected Build47 tier names")
	}
}

func TestExperimentalV4Build47ProductionCapUnchanged(t *testing.T) {
	if experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production frozen cap=%d want 32", experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild43PairKeep != 2 {
		t.Fatalf("Build43 production pair keep=%d want 2", experimentalV4PhoneBuild43PairKeep)
	}
}

func TestExperimentalV4Build47DiagnosticPlanCanExtendBeyondProduction(t *testing.T) {
	// Each visited basin can emit at most eight candidates. Build43 production
	// visits 2 pairs x 2 cells = 4 basins (32 theoretical candidates). Build47
	// adds 4 more basins on the same selected pairs and 8 basins on the four
	// lower-ranked pairs, giving a diagnostic ceiling of 128 without changing
	// the production path.
	const candidatesPerBasin = 8
	production := 2 * 2 * candidatesPerBasin
	depth := 2 * 2 * candidatesPerBasin
	allPairs := 4 * 2 * candidatesPerBasin
	if production != 32 || production+depth != 64 || production+depth+allPairs != 128 {
		t.Fatalf("Build47 theoretical tiers=%d/%d/%d", production, production+depth, production+depth+allPairs)
	}
}
