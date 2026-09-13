package watermark

import (
	"math"
	"testing"
)

func TestDiagnosticPhysicalTopologyPairPartitionIsHeldOut(t *testing.T) {
	for _, profile := range []Profile{ProfileRobust, ProfileBalanced} {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				common, exclusive := diagnosticPhysicalTopologyPairPartition(profile, dx, dy)
				if len(common) == 0 || len(exclusive) == 0 {
					t.Fatalf("profile=%s shift=(%d,%d) empty partition common=%d exclusive=%d", profile, dx, dy, len(common), len(exclusive))
				}
				seen := map[uint32]bool{}
				for _, pair := range common {
					seen[diagnosticObservabilityPairKey(pair[0], pair[1])] = true
				}
				for _, pair := range exclusive {
					key := diagnosticObservabilityPairKey(pair[0], pair[1])
					if seen[key] {
						t.Fatalf("profile=%s shift=(%d,%d) held-out pair reused for registration", profile, dx, dy)
					}
				}
			}
		}
	}
}

func TestDiagnosticPhysicalTopologyObservabilitySeparatesSyntheticCarrier(t *testing.T) {
	aggregate := syntheticPhysicalTopologyGrid(ProfileBalanced, 0)
	cells := make([]diagnosticSpatialGridCell, 0, 9)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			grid := syntheticPhysicalTopologyGrid(ProfileBalanced, y*3+x+1)
			cells = append(cells, diagnosticSpatialGridCell{RegionX: x, RegionY: y, grid: grid, counts: make([]int, eccBits)})
		}
	}
	result := diagnosticPhysicalTopologyObservability(cells, aggregate)
	if !result.Available || result.ProfilesAvailable != 2 || result.AuditedShifts != 16 {
		t.Fatalf("unexpected observability summary: %+v", result)
	}
	var balanced *DiagnosticPhysicalTopologyProfile
	for i := range result.Profiles {
		if result.Profiles[i].Profile == ProfileBalanced {
			balanced = &result.Profiles[i]
			break
		}
	}
	if balanced == nil || !balanced.Available || balanced.ShiftsCompared != 8 {
		t.Fatalf("balanced profile unavailable: %+v", balanced)
	}
	if balanced.ModalWinnerFraction < 0.5 || balanced.MeanDirectionalConsistency < 0.65 || balanced.MeanAbsoluteCellDelta <= 0 {
		t.Fatalf("synthetic balanced topology not separated strongly enough: %+v", balanced)
	}
}

func syntheticPhysicalTopologyGrid(profile Profile, variant int) []float64 {
	spec, _ := profileSpecFor(profile)
	state := uint64(0x4d595df4d0f33173)
	coded := make([]float64, spec.codedBits)
	for i := range coded {
		state ^= state << 13
		state ^= state >> 7
		state ^= state << 17
		if state&1 != 0 {
			coded[i] = 1
		} else {
			coded[i] = -1
		}
	}
	grid := make([]float64, eccBits)
	for position := range grid {
		base := coded[v3CodeIndex(position, spec.codedBits)]
		noise := 0.04 * math.Sin(float64((position+1)*(variant+3)))
		grid[position] = base + noise
	}
	return grid
}
