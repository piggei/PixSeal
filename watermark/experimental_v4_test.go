package watermark

import "testing"

func TestExperimentalV4PrototypeFoundationInvariants(t *testing.T) {
	if len(experimentalV4PrototypePilotPositions) != experimentalV4PilotCount || len(experimentalV4PrototypePilotSigns) != experimentalV4PilotCount {
		t.Fatalf("unexpected pilot lengths: positions=%d signs=%d", len(experimentalV4PrototypePilotPositions), len(experimentalV4PrototypePilotSigns))
	}

	seen := make(map[int]bool, experimentalV4PilotCount)
	for i, position := range experimentalV4PrototypePilotPositions {
		if position < 0 || position >= experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks {
			t.Fatalf("pilot %d has out-of-range position %d", i, position)
		}
		if seen[position] {
			t.Fatalf("duplicate pilot position %d", position)
		}
		seen[position] = true
		sign := experimentalV4PrototypePilotSigns[i]
		if sign != -1 && sign != 1 {
			t.Fatalf("pilot %d has invalid sign %d", i, sign)
		}
	}

	metrics := experimentalV4PrototypePilotMetrics()
	if metrics.Version != 4 || metrics.TileWidthBlocks != 37 || metrics.TileHeightBlocks != 32 {
		t.Fatalf("unexpected v4 geometry: %+v", metrics)
	}
	if metrics.PilotPositions != 64 || metrics.DataPositions != 1120 {
		t.Fatalf("unexpected pilot/data split: %+v", metrics)
	}
	if metrics.PositiveSigns != 32 || metrics.NegativeSigns != 32 {
		t.Fatalf("pilot is not sign-balanced: %+v", metrics)
	}
	if metrics.MinimumQuadrantPilots != 16 || metrics.MaximumQuadrantPilots != 16 {
		t.Fatalf("unexpected quadrant distribution: %+v", metrics)
	}
	if metrics.PerfectCyclicAlias {
		t.Fatalf("prototype must not have a perfect non-zero cyclic alias: %+v", metrics)
	}
	if metrics.MaximumCyclicOverlap > 8 {
		t.Fatalf("prototype cyclic overlap regressed: %+v", metrics)
	}
	if metrics.MaximumWrongSignedCorrelation > 5 {
		t.Fatalf("prototype signed autocorrelation regressed: %+v", metrics)
	}

	data := experimentalV4PrototypeDataPositions()
	if len(data) != experimentalV4DataCount {
		t.Fatalf("unexpected data-position count: got %d want %d", len(data), experimentalV4DataCount)
	}
	for _, position := range data {
		if seen[position] {
			t.Fatalf("data position %d overlaps the pilot", position)
		}
	}
}

func TestExperimentalV4PrototypeUsesOnePilotPerStratum(t *testing.T) {
	counts := [64]int{}
	for _, position := range experimentalV4PrototypePilotPositions {
		x := position % experimentalV4TileWidthBlocks
		y := position / experimentalV4TileWidthBlocks
		gx := 0
		for gx < 7 && x >= ((gx+1)*experimentalV4TileWidthBlocks+4)/8 {
			gx++
		}
		gy := y / 4
		if gy > 7 {
			gy = 7
		}
		counts[gy*8+gx]++
	}
	for stratum, count := range counts {
		if count != 1 {
			t.Fatalf("stratum %d has %d pilots, want exactly 1", stratum, count)
		}
	}
}
