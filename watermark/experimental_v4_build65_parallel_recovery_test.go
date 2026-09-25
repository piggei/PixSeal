package watermark

import (
	"image"
	"testing"
)

func TestExperimentalV4Build65PreservesBuild64Bounds(t *testing.T) {
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build65 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
	if experimentalV4PhoneBuild55Step != 1 || experimentalV4PhoneBuild55Dimensions != 8 || experimentalV4PhoneBuild55MaxPasses != 8 || experimentalV4PhoneBuild55MaxStatesBranch != 64 {
		t.Fatalf("Build65 continuation bounds changed")
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build65 pair bounds changed")
	}
}

func TestExperimentalV4Build65BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers := experimentalV4PhoneBuild65BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
}

func TestExperimentalV4Build65ProductionThresholdsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build65 changed qualified production thresholds")
	}
}
