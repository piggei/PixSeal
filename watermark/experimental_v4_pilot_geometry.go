package watermark

import (
	"image"
	"math"
)

// experimentalV4DetectPilotProjective evaluates every 37x32 cyclic pilot
// origin after a canonical-to-observed projective mapping has already been
// resolved. Build25 deliberately keeps geometry estimation and absolute-origin
// evidence separate: this function asks whether the public pilot survives the
// geometric channel once the correct mapping is known. It does not authenticate
// a payload and it is not reachable from the production v3 decoder.
func experimentalV4DetectPilotProjective(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) ExperimentalV4PilotDetection {
	if img == nil || canonicalWidth < experimentalV4TileWidthBlocks*blockSize || canonicalHeight < experimentalV4TileHeightBlocks*blockSize {
		return ExperimentalV4PilotDetection{}
	}
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	if blocksWide <= 0 || blocksHigh <= 0 {
		return ExperimentalV4PilotDetection{}
	}

	plane := newPixelPlane(img)
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var sums [residues]float64
	var absSums [residues]float64
	var counts [residues]int
	for blockY := 0; blockY < blocksHigh; blockY++ {
		for blockX := 0; blockX < blocksWide; blockX++ {
			value, ok := readProjectiveBlockValue(plane, h, blockX*blockSize, blockY*blockSize, blockSize)
			if !ok {
				continue
			}
			residue := (blockY%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + blockX%experimentalV4TileWidthBlocks
			sums[residue] += value
			absSums[residue] += math.Abs(value)
			counts[residue]++
		}
	}

	bestScore := math.Inf(-1)
	runnerUpScore := math.Inf(-1)
	bestOriginX, bestOriginY := 0, 0
	bestVisible, bestSamples := 0, 0
	for originY := 0; originY < experimentalV4TileHeightBlocks; originY++ {
		for originX := 0; originX < experimentalV4TileWidthBlocks; originX++ {
			numerator := 0.0
			denominator := 0.0
			visible := 0
			samples := 0
			for index, pilotPosition := range candidate.positions {
				pilotX := pilotPosition % experimentalV4TileWidthBlocks
				pilotY := pilotPosition / experimentalV4TileWidthBlocks
				observedX := positiveMod(pilotX-originX, experimentalV4TileWidthBlocks)
				observedY := positiveMod(pilotY-originY, experimentalV4TileHeightBlocks)
				residue := observedY*experimentalV4TileWidthBlocks + observedX
				if counts[residue] == 0 {
					continue
				}
				visible++
				samples += counts[residue]
				numerator += float64(candidate.signs[index]) * sums[residue]
				denominator += absSums[residue]
			}
			if denominator <= 0 {
				continue
			}
			score := numerator / denominator
			if score > bestScore {
				runnerUpScore = bestScore
				bestScore = score
				bestOriginX, bestOriginY = originX, originY
				bestVisible, bestSamples = visible, samples
			} else if score > runnerUpScore {
				runnerUpScore = score
			}
		}
	}
	if math.IsInf(bestScore, -1) {
		return ExperimentalV4PilotDetection{}
	}
	if math.IsInf(runnerUpScore, -1) {
		runnerUpScore = 0
	}
	return ExperimentalV4PilotDetection{
		Available:             true,
		BlockSizePixels:       blockSize,
		OriginXBlocks:         bestOriginX,
		OriginYBlocks:         bestOriginY,
		Score:                 bestScore,
		RunnerUpScore:         runnerUpScore,
		Margin:                bestScore - runnerUpScore,
		VisiblePilotPositions: bestVisible,
		PilotSamples:          bestSamples,
	}
}
