package watermark

import (
	"errors"
	"image"
	"math"
)

const experimentalV4SyntheticDataSeed uint64 = 0x7634646174613031 // "v4data01"

// ExperimentalV4PilotDetection is image-domain pilot-only telemetry. It is an
// experimental geometry observable and never authenticates a payload.
type ExperimentalV4PilotDetection struct {
	Available             bool    `json:"available"`
	BlockSizePixels       int     `json:"block_size_px"`
	PhaseX                int     `json:"phase_x_px"`
	PhaseY                int     `json:"phase_y_px"`
	OriginXBlocks         int     `json:"origin_x_blocks"`
	OriginYBlocks         int     `json:"origin_y_blocks"`
	Score                 float64 `json:"score"`
	RunnerUpScore         float64 `json:"runner_up_score"`
	Margin                float64 `json:"margin"`
	VisiblePilotPositions int     `json:"visible_pilot_positions"`
	PilotSamples          int     `json:"pilot_samples"`
}

// experimentalV4RenderSyntheticCarrier embeds a public pilot plus deterministic
// pseudo-random data-plane signs. It exists only to exercise the pilot against
// realistic DCT interference before a real v4 frame/ECC encoder exists.
func experimentalV4RenderSyntheticCarrier(src image.Image, candidate experimentalV4PilotCandidate, strength float64, dataSeed uint64) (*image.NRGBA, error) {
	if err := validateWorkingImageSize(src); err != nil {
		return nil, err
	}
	if math.IsNaN(strength) || math.IsInf(strength, 0) || strength < 4 || strength > 120 {
		return nil, errors.New("experimental v4 strength must be finite and between 4 and 120")
	}
	bounds := src.Bounds()
	if bounds.Dx() < experimentalV4TileWidthBlocks*blockSize || bounds.Dy() < experimentalV4TileHeightBlocks*blockSize {
		return nil, errors.New("image must be at least 296x256 pixels for the experimental v4 pilot")
	}

	pilotSigns := make(map[int]int8, experimentalV4PilotCount)
	for index, position := range candidate.positions {
		pilotSigns[position] = candidate.signs[index]
	}

	out := toNRGBA(src)
	blocksWide := bounds.Dx() / blockSize
	blocksHigh := bounds.Dy() / blockSize
	for blockY := 0; blockY < blocksHigh; blockY++ {
		for blockX := 0; blockX < blocksWide; blockX++ {
			tilePosition := (blockY%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + blockX%experimentalV4TileWidthBlocks
			sign, pilot := pilotSigns[tilePosition]
			if !pilot {
				sign = experimentalV4SyntheticDataSign(tilePosition, dataSeed)
			}
			bit := byte(0)
			if sign > 0 {
				bit = 1
			}
			embedBlock(out, point{blockX * blockSize, blockY * blockSize}, bit, strength)
		}
	}
	return out, nil
}

func experimentalV4SyntheticDataSign(position int, seed uint64) int8 {
	x := seed + uint64(position+1)*0x9e3779b97f4a7c15
	x ^= x >> 30
	x *= 0xbf58476d1ce4e5b9
	x ^= x >> 27
	x *= 0x94d049bb133111eb
	x ^= x >> 31
	if x&1 == 0 {
		return -1
	}
	return 1
}

// experimentalV4DetectPilotGrid evaluates every 37x32 cyclic origin for one
// already-resolved pixel lattice. Geometry search is intentionally separate:
// Build24 first asks whether the pilot itself identifies absolute tile origin.
func experimentalV4DetectPilotGrid(img image.Image, candidate experimentalV4PilotCandidate, blockSizePixels, phaseX, phaseY int) ExperimentalV4PilotDetection {
	if img == nil || (blockSizePixels != 8 && blockSizePixels != 6 && blockSizePixels != 4) {
		return ExperimentalV4PilotDetection{}
	}
	if phaseX < 0 || phaseY < 0 || phaseX >= blockSizePixels || phaseY >= blockSizePixels {
		return ExperimentalV4PilotDetection{}
	}
	bounds := img.Bounds()
	blocksWide := (bounds.Dx() - phaseX) / blockSizePixels
	blocksHigh := (bounds.Dy() - phaseY) / blockSizePixels
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
			position := point{
				x: bounds.Min.X + phaseX + blockX*blockSizePixels,
				y: bounds.Min.Y + phaseY + blockY*blockSizePixels,
			}
			value := readBlockSized(plane, position, blockSizePixels)
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
		BlockSizePixels:       blockSizePixels,
		PhaseX:                phaseX,
		PhaseY:                phaseY,
		OriginXBlocks:         bestOriginX,
		OriginYBlocks:         bestOriginY,
		Score:                 bestScore,
		RunnerUpScore:         runnerUpScore,
		Margin:                bestScore - runnerUpScore,
		VisiblePilotPositions: bestVisible,
		PilotSamples:          bestSamples,
	}
}
