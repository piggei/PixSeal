package watermark

import "testing"

func TestExperimentalV4Build64RecoveryUsesFrozenBuild63Bounds(t *testing.T) {
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build64 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
	if experimentalV4PhoneBuild55Step != 1 || experimentalV4PhoneBuild55Dimensions != 8 || experimentalV4PhoneBuild55MaxPasses != 8 || experimentalV4PhoneBuild55MaxStatesBranch != 64 {
		t.Fatalf("Build64 continuation bounds step/dims/passes/states=%v/%d/%d/%d want 1/8/8/64", experimentalV4PhoneBuild55Step, experimentalV4PhoneBuild55Dimensions, experimentalV4PhoneBuild55MaxPasses, experimentalV4PhoneBuild55MaxStatesBranch)
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build64 pair bounds combos/keep=%d/%d want 112/8", experimentalV4PhoneBuild53PairCombosPerRoot, experimentalV4PhoneBuild53PairKeepPerRoot)
	}
}

func TestExperimentalV4Build64LeavesQualifiedProductionThresholdsFrozen(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("phone qualification thresholds changed")
	}
}

func TestExperimentalV4Build64BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := testImage(128, 128)
	bank, _, _ := experimentalV4PhoneBuild64BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 {
		t.Fatalf("invalid boundary produced %d recovery candidates", len(bank))
	}
}
