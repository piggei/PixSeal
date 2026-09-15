package watermark

import "testing"

func TestExperimentalV4PilotCandidateLockIdentity(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	if candidate.name != experimentalV4PilotCandidateLockName {
		t.Fatalf("locked v4 pilot name=%q, want %q", candidate.name, experimentalV4PilotCandidateLockName)
	}
	if got := experimentalV4PilotCandidateHash(candidate); got != experimentalV4PilotCandidateLockHash {
		t.Fatalf("locked v4 pilot hash=%s, want %s", got, experimentalV4PilotCandidateLockHash)
	}

	seen := make(map[int]struct{}, experimentalV4PilotCount)
	positive, negative := 0, 0
	for i, position := range candidate.positions {
		if position < 0 || position >= experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks {
			t.Fatalf("pilot position %d out of range at index %d", position, i)
		}
		if _, duplicate := seen[position]; duplicate {
			t.Fatalf("duplicate locked pilot position %d", position)
		}
		seen[position] = struct{}{}
		switch candidate.signs[i] {
		case 1:
			positive++
		case -1:
			negative++
		default:
			t.Fatalf("invalid locked pilot sign %d at index %d", candidate.signs[i], i)
		}
	}
	if len(seen) != 64 || positive != 32 || negative != 32 {
		t.Fatalf("unexpected locked pilot balance: unique=%d positive=%d negative=%d", len(seen), positive, negative)
	}

	qualification := experimentalV4PilotQualification(candidate)
	if qualification.Cyclic.MaximumMaskOverlap != 7 || qualification.Cyclic.MaximumWrongSignedCorrelation != 4 || qualification.Cyclic.PerfectNonZeroCyclicAliases != 0 {
		t.Fatalf("locked pilot structural invariants changed: %+v", qualification.Cyclic)
	}
}

func TestExperimentalV4PilotCandidateLockDetectsMutation(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	candidate.signs[0] = -candidate.signs[0]
	if got := experimentalV4PilotCandidateHash(candidate); got == experimentalV4PilotCandidateLockHash {
		t.Fatalf("pilot identity hash did not change after sign mutation")
	}
}
