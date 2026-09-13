package watermark

import "testing"

func TestDiagnosticFormatObservabilityAuditIsDeterministicAndKeyIndependent(t *testing.T) {
	a := diagnosticFormatObservabilityAudit()
	b := diagnosticFormatObservabilityAudit()
	if a.Method != "v3-key-independent-structural-observability" || len(a.Profiles) != 3 || len(b.Profiles) != 3 || a.ProfilesWithRepetitionAnchor != 2 || a.ProfilesWithSyntheticHammingAlias != 2 || a.AuditedUnitShiftsUniversallySeparated {
		t.Fatalf("unexpected audit shape: %+v", a)
	}
	for i := range a.Profiles {
		if len(a.Profiles[i].UnitShifts) != 8 || len(b.Profiles[i].UnitShifts) != 8 {
			t.Fatalf("profile %s unit shifts=%d/%d", a.Profiles[i].Profile, len(a.Profiles[i].UnitShifts), len(b.Profiles[i].UnitShifts))
		}
		for j := range a.Profiles[i].UnitShifts {
			x, y := a.Profiles[i].UnitShifts[j], b.Profiles[i].UnitShifts[j]
			if x != y {
				t.Fatalf("non-deterministic observability result profile=%s shift=%d: %+v != %+v", a.Profiles[i].Profile, j, x, y)
			}
		}
	}
}

func TestDiagnosticFormatObservabilityRepetitionTopologyIsWeakButNonInvariant(t *testing.T) {
	audit := diagnosticFormatObservabilityAudit()
	for _, profile := range audit.Profiles {
		switch profile.Profile {
		case ProfileRobust, ProfileBalanced:
			if profile.RepetitionPairs == 0 || !profile.RepetitionAbsoluteAnchor {
				t.Fatalf("%s should have a non-invariant repetition topology: %+v", profile.Profile, profile)
			}
			if profile.MinimumPairTopologyGap <= 0 || profile.MinimumPairTopologyGap >= 0.10 {
				t.Fatalf("%s nearest unit-cycle topology gap should be real but weak, got %f", profile.Profile, profile.MinimumPairTopologyGap)
			}
		case ProfileCapacity:
			if profile.RepetitionPairs != 0 || profile.RepetitionAbsoluteAnchor || profile.MinimumPairTopologyGap != 0 {
				t.Fatalf("capacity must expose no repetition-pair anchor: %+v", profile)
			}
		}
	}
}

func TestDiagnosticFormatObservabilityCapacityVerticalHammingAlias(t *testing.T) {
	audit := diagnosticFormatObservabilityAudit()
	for _, profile := range audit.Profiles {
		if profile.Profile != ProfileCapacity {
			continue
		}
		aliases := 0
		for _, shift := range profile.UnitShifts {
			if shift.DX == 0 && (shift.DY == -1 || shift.DY == 1) {
				if !shift.PureCodeIndexPermutation || !shift.HammingWordPermutationInvariant || shift.SyntheticHammingSyndrome != 0 {
					t.Fatalf("capacity vertical shift %+v must be an exact Hamming-word permutation alias", shift)
				}
				aliases++
			}
		}
		if aliases != 2 || !profile.SyntheticHammingUnitShiftAliasPresent {
			t.Fatalf("capacity vertical Hamming aliases=%d profile=%+v", aliases, profile)
		}
		return
	}
	t.Fatal("capacity observability profile not found")
}
