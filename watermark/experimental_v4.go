package watermark

// Build23 starts the experimental Format-v4 branch without changing any v3
// encoder/decoder behavior. Everything in this file is deliberately isolated
// from EmbedWithInfo, ExtractWithInfo and the frozen v3 format.
const (
	experimentalV4Version          = 4
	experimentalV4TileWidthBlocks  = 37
	experimentalV4TileHeightBlocks = 32
	experimentalV4PilotCount       = 64
	experimentalV4DataCount        = 1120
	experimentalV4PrototypeName    = "prototype-1-stratified-p64"
)

// ExperimentalV4PilotMetrics reports structural properties of the provisional
// build23 pilot. These are design invariants, not physical-channel guarantees.
type ExperimentalV4PilotMetrics struct {
	Name                          string
	Version                       int
	TileWidthBlocks               int
	TileHeightBlocks              int
	PilotPositions                int
	DataPositions                 int
	PositiveSigns                 int
	NegativeSigns                 int
	MinimumQuadrantPilots         int
	MaximumQuadrantPilots         int
	MaximumCyclicOverlap          int
	MaximumWrongSignedCorrelation int
	PerfectCyclicAlias            bool
}

// The build23 candidate was selected by a deterministic offline search over a
// stratified 8x8 placement family. There is exactly one pilot in each spatial
// stratum and an exactly balanced +/- sign sequence. This candidate is NOT a
// normative v4 constant yet; future builds may replace it after simulated and
// physical qualification.
var experimentalV4PrototypePilotPositions = [...]int{
	74, 81, 86, 125, 19, 138, 139, 144,
	223, 153, 234, 202, 167, 248, 290, 182,
	370, 413, 381, 386, 318, 359, 399, 369,
	518, 562, 456, 534, 537, 506, 509, 587,
	629, 674, 716, 609, 649, 726, 623, 663,
	815, 858, 752, 792, 870, 876, 805, 884,
	966, 896, 1011, 1016, 982, 911, 1028, 924,
	1073, 1079, 1049, 1051, 1057, 1096, 1176, 1108,
}

var experimentalV4PrototypePilotSigns = [...]int8{
	-1, -1, -1, -1, 1, -1, 1, 1,
	-1, -1, 1, 1, -1, -1, -1, -1,
	1, -1, -1, -1, -1, -1, 1, 1,
	1, -1, -1, -1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, -1, -1, -1, 1, -1,
	1, 1, 1, -1, -1, 1, -1, 1,
	1, -1, -1, -1, 1, -1, -1, -1,
}

func experimentalV4PrototypePilotMetrics() ExperimentalV4PilotMetrics {
	const positions = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	occupied := make(map[int]int8, experimentalV4PilotCount)
	positive := 0
	negative := 0
	quadrants := [4]int{}
	for i, position := range experimentalV4PrototypePilotPositions {
		sign := experimentalV4PrototypePilotSigns[i]
		occupied[position] = sign
		if sign > 0 {
			positive++
		} else if sign < 0 {
			negative++
		}
		x := position % experimentalV4TileWidthBlocks
		y := position / experimentalV4TileWidthBlocks
		qx := 0
		if 2*x >= experimentalV4TileWidthBlocks {
			qx = 1
		}
		qy := 0
		if 2*y >= experimentalV4TileHeightBlocks {
			qy = 1
		}
		quadrants[qy*2+qx]++
	}

	maxOverlap := 0
	maxAbsCorrelation := 0
	perfectAlias := false
	for dy := 0; dy < experimentalV4TileHeightBlocks; dy++ {
		for dx := 0; dx < experimentalV4TileWidthBlocks; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			overlap := 0
			correlation := 0
			for position, sign := range occupied {
				x := position % experimentalV4TileWidthBlocks
				y := position / experimentalV4TileWidthBlocks
				shiftedX := (x + dx) % experimentalV4TileWidthBlocks
				shiftedY := (y + dy) % experimentalV4TileHeightBlocks
				shifted := shiftedY*experimentalV4TileWidthBlocks + shiftedX
				if shiftedSign, ok := occupied[shifted]; ok {
					overlap++
					correlation += int(sign * shiftedSign)
				}
			}
			if overlap > maxOverlap {
				maxOverlap = overlap
			}
			if correlation < 0 {
				correlation = -correlation
			}
			if correlation > maxAbsCorrelation {
				maxAbsCorrelation = correlation
			}
			if overlap == experimentalV4PilotCount && correlation == experimentalV4PilotCount {
				perfectAlias = true
			}
		}
	}

	minQuadrant := experimentalV4PilotCount
	maxQuadrant := 0
	for _, count := range quadrants {
		if count < minQuadrant {
			minQuadrant = count
		}
		if count > maxQuadrant {
			maxQuadrant = count
		}
	}

	return ExperimentalV4PilotMetrics{
		Name:                          experimentalV4PrototypeName,
		Version:                       experimentalV4Version,
		TileWidthBlocks:               experimentalV4TileWidthBlocks,
		TileHeightBlocks:              experimentalV4TileHeightBlocks,
		PilotPositions:                len(occupied),
		DataPositions:                 positions - len(occupied),
		PositiveSigns:                 positive,
		NegativeSigns:                 negative,
		MinimumQuadrantPilots:         minQuadrant,
		MaximumQuadrantPilots:         maxQuadrant,
		MaximumCyclicOverlap:          maxOverlap,
		MaximumWrongSignedCorrelation: maxAbsCorrelation,
		PerfectCyclicAlias:            perfectAlias,
	}
}

func experimentalV4PrototypeDataPositions() []int {
	pilot := make(map[int]struct{}, experimentalV4PilotCount)
	for _, position := range experimentalV4PrototypePilotPositions {
		pilot[position] = struct{}{}
	}
	positions := make([]int, 0, experimentalV4DataCount)
	for position := 0; position < experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks; position++ {
		if _, reserved := pilot[position]; !reserved {
			positions = append(positions, position)
		}
	}
	return positions
}
