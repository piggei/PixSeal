package watermark

import (
	"errors"
	"runtime"
	"testing"
)

func TestExperimentalV4Build66PreservesBuild64Bounds(t *testing.T) {
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build66 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
	if experimentalV4PhoneBuild55Step != 1 || experimentalV4PhoneBuild55Dimensions != 8 || experimentalV4PhoneBuild55MaxPasses != 8 || experimentalV4PhoneBuild55MaxStatesBranch != 64 {
		t.Fatalf("Build66 continuation bounds changed")
	}
	if experimentalV4PhoneBuild53PairCombosPerRoot != 112 || experimentalV4PhoneBuild53PairKeepPerRoot != 8 {
		t.Fatalf("Build66 pair bounds changed")
	}
}

func TestExperimentalV4Build66DecodeWorkerCount(t *testing.T) {
	old := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(old)
	runtime.GOMAXPROCS(4)
	for _, tc := range []struct {
		candidates int
		want       int
	}{
		{0, 0},
		{1, 1},
		{3, 3},
		{4, 4},
		{9, 4},
	} {
		if got := experimentalV4PhoneBuild66DecodeWorkerCount(tc.candidates); got != tc.want {
			t.Fatalf("candidates=%d got workers=%d want %d", tc.candidates, got, tc.want)
		}
	}
}

func TestExperimentalV4Build66ProductionThresholdsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build66 changed qualified production thresholds")
	}
}

func TestExperimentalV4Build66BatchConsumptionStopsAtFirstLogicalSuccess(t *testing.T) {
	telemetry := experimentalV4PhoneBuild66RecoveryTelemetry{}
	results := []experimentalV4PhoneBuild66DecodeResult{
		{profiles: 2, frames: 100, maxConfidence: 1.5, err: errors.New("no auth")},
		{payload: []byte("winner"), info: ExperimentalV4ExtractInfo{Profile: ProfileRobust}, profiles: 3, frames: 200, maxConfidence: 2.5, err: nil},
		{payload: []byte("speculative"), info: ExperimentalV4ExtractInfo{Profile: ProfileCapacity}, profiles: 9, frames: 9999, maxConfidence: 99, err: nil},
	}
	payload, info, ok := experimentalV4PhoneBuild66ConsumeBatch(results, &telemetry)
	if !ok || string(payload) != "winner" || info.Profile != ProfileRobust {
		t.Fatalf("unexpected winner ok=%t payload=%q profile=%q", ok, payload, info.Profile)
	}
	if telemetry.DecodeCandidatesTried != 2 || telemetry.ProfilesTried != 5 || telemetry.ListFramesTried != 300 {
		t.Fatalf("logical telemetry candidates/profiles/frames=%d/%d/%d want 2/5/300", telemetry.DecodeCandidatesTried, telemetry.ProfilesTried, telemetry.ListFramesTried)
	}
	if telemetry.MaxDataConfidence != 2.5 {
		t.Fatalf("max confidence=%v want 2.5; speculative result leaked", telemetry.MaxDataConfidence)
	}
	if !telemetry.Authenticated || telemetry.Profile != ProfileRobust {
		t.Fatalf("authentication telemetry mismatch authenticated=%t profile=%q", telemetry.Authenticated, telemetry.Profile)
	}
}
