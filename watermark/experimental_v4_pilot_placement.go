package watermark

import (
	"image"
	"math"
	"sort"
)

// Build27 isolates the placement problem from geometry, exactly as Build25
// isolated pilot observability from blind geometry. The caller supplies the
// canonical->full-frame geometry (rotation/scale/shear/projective shape) but no
// crop or canvas translation. The public pilot then recovers the unknown
// placement modulo the repeated v4 tile lattice. Payload/header/ECC/HMAC are
// never consulted.

type experimentalV4PlacementResult struct {
	available           bool
	shiftX              float64
	shiftY              float64
	h                   homography
	proposalScore       float64
	validationScore     float64
	validationRunnerUp  float64
	visibleProposal     int
	visibleValidation   int
	hypothesesEvaluated int
}

type experimentalV4PlacementHypothesis struct {
	x, y              int
	proposal          float64
	validation        float64
	visibleProposal   int
	visibleValidation int
	h                 homography
}

// experimentalV4PilotPartitionScore uses disjoint halves of the public pilot.
// partition 0 proposes placement; partition 1 validates it. All repeated tiles
// are allowed to contribute, but an individual pilot symbol belongs to exactly
// one stage, preserving proposal/validation separation without payload or HMAC.
func experimentalV4PilotPartitionScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography, partition int) (float64, int) {
	if plane == nil || canonicalWidth < experimentalV4TileWidthBlocks*blockSize || canonicalHeight < experimentalV4TileHeightBlocks*blockSize {
		return math.Inf(-1), 0
	}
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	tilesX := blocksWide / experimentalV4TileWidthBlocks
	tilesY := blocksHigh / experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return math.Inf(-1), 0
	}

	numerator, denominator := 0.0, 0.0
	visible := 0
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			baseX := tx * experimentalV4TileWidthBlocks
			baseY := ty * experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				if i&1 != partition {
					continue
				}
				px := pos % experimentalV4TileWidthBlocks
				py := pos / experimentalV4TileWidthBlocks
				v, ok := readProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
				if !ok {
					continue
				}
				visible++
				numerator += float64(candidate.signs[i]) * v
				denominator += math.Abs(v)
			}
		}
	}
	if denominator <= 0 {
		return math.Inf(-1), visible
	}
	return numerator / denominator, visible
}

func experimentalV4PlacementShiftHomography(base homography, x, y float64) homography {
	shift := homography{h: [9]float64{1, 0, x, 0, 1, y, 0, 0, 1}}
	return experimentalV4BlindMultiplyHomography(shift, base)
}

// experimentalV4PlacementBounds derives the only translations compatible with
// a crop or a padded canvas when the full transformed extent is known. If the
// observed frame is smaller, the carrier may have been cropped from either
// side; if it is larger, the carrier may be placed anywhere inside the canvas.
// A small guard absorbs integer rounding in synthetic/real resampling.
func experimentalV4PlacementBounds(full, observed int) (int, int) {
	const guard = 6
	delta := observed - full
	if delta >= 0 {
		return -guard, delta + guard
	}
	return delta - guard, guard
}

// experimentalV4PilotPlacementSearch is bounded and deterministic. A 4-pixel
// grid proposes candidate translations using pilot half A; the best candidates
// are refined at integer-pixel resolution and ranked only by held-out pilot half
// B. Whole-tile-equivalent translations may legitimately tie because v4 repeats
// the same tile; that is an equivalence, not a decoder ambiguity.
func experimentalV4PilotPlacementSearch(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, geometry homography, fullWidth, fullHeight int) experimentalV4PlacementResult {
	if img == nil || fullWidth <= 0 || fullHeight <= 0 {
		return experimentalV4PlacementResult{}
	}
	plane := newPixelPlane(img)
	observedWidth, observedHeight := img.Bounds().Dx(), img.Bounds().Dy()
	minX, maxX := experimentalV4PlacementBounds(fullWidth, observedWidth)
	minY, maxY := experimentalV4PlacementBounds(fullHeight, observedHeight)
	// Keep the research search explicitly bounded. Current qualification crops
	// and placements are deliberately moderate; larger placement recovery is a
	// future coarse-localization problem, not a reason for unbounded scanning.
	if maxX-minX > 384 || maxY-minY > 384 {
		return experimentalV4PlacementResult{}
	}

	hypotheses := 0
	coarse := make([]experimentalV4PlacementHypothesis, 0, 24)
	insertCoarse := func(x, y int) {
		h := experimentalV4PlacementShiftHomography(geometry, float64(x), float64(y))
		score, visible := experimentalV4PilotPartitionScore(plane, candidate, canonicalWidth, canonicalHeight, h, 0)
		hypotheses++
		if visible < 24 || math.IsInf(score, -1) {
			return
		}
		coarse = append(coarse, experimentalV4PlacementHypothesis{x: x, y: y, proposal: score, visibleProposal: visible, h: h})
		sort.Slice(coarse, func(i, j int) bool {
			if coarse[i].proposal == coarse[j].proposal {
				return coarse[i].visibleProposal > coarse[j].visibleProposal
			}
			return coarse[i].proposal > coarse[j].proposal
		})
		if len(coarse) > 24 {
			coarse = coarse[:24]
		}
	}
	const coarseStep = 4
	for y := minY; y <= maxY; y += coarseStep {
		for x := minX; x <= maxX; x += coarseStep {
			insertCoarse(x, y)
		}
	}
	// Always include exact range endpoints because a 4-pixel grid may not land there.
	for _, y := range []int{minY, maxY} {
		for _, x := range []int{minX, maxX} {
			insertCoarse(x, y)
		}
	}
	if len(coarse) == 0 {
		return experimentalV4PlacementResult{hypothesesEvaluated: hypotheses}
	}

	fineSeen := make(map[[2]int]bool)
	refined := make([]experimentalV4PlacementHypothesis, 0, 96)
	seeds := coarse
	if len(seeds) > 12 {
		seeds = seeds[:12]
	}
	for _, seed := range seeds {
		for y := maxIntV4Placement(minY, seed.y-4); y <= minInt(maxY, seed.y+4); y++ {
			for x := maxIntV4Placement(minX, seed.x-4); x <= minInt(maxX, seed.x+4); x++ {
				key := [2]int{x, y}
				if fineSeen[key] {
					continue
				}
				fineSeen[key] = true
				h := experimentalV4PlacementShiftHomography(geometry, float64(x), float64(y))
				proposal, vp := experimentalV4PilotPartitionScore(plane, candidate, canonicalWidth, canonicalHeight, h, 0)
				hypotheses++
				if vp < 24 || math.IsInf(proposal, -1) {
					continue
				}
				refined = append(refined, experimentalV4PlacementHypothesis{x: x, y: y, proposal: proposal, visibleProposal: vp, h: h})
			}
		}
	}
	if len(refined) == 0 {
		return experimentalV4PlacementResult{hypothesesEvaluated: hypotheses}
	}
	sort.Slice(refined, func(i, j int) bool { return refined[i].proposal > refined[j].proposal })
	if len(refined) > 32 {
		refined = refined[:32]
	}

	validated := make([]experimentalV4PlacementHypothesis, 0, len(refined))
	for _, q := range refined {
		validation, vv := experimentalV4PilotPartitionScore(plane, candidate, canonicalWidth, canonicalHeight, q.h, 1)
		if vv < 24 || math.IsInf(validation, -1) {
			continue
		}
		q.validation, q.visibleValidation = validation, vv
		validated = append(validated, q)
	}
	if len(validated) == 0 {
		return experimentalV4PlacementResult{hypothesesEvaluated: hypotheses}
	}
	// Final ranking uses held-out evidence only. Proposal score breaks exact ties
	// but cannot promote a placement that validation does not support.
	sort.Slice(validated, func(i, j int) bool {
		if validated[i].validation == validated[j].validation {
			return validated[i].proposal > validated[j].proposal
		}
		return validated[i].validation > validated[j].validation
	})
	win := validated[0]
	runner := 0.0
	if len(validated) > 1 {
		runner = validated[1].validation
	}
	return experimentalV4PlacementResult{available: true, shiftX: float64(win.x), shiftY: float64(win.y), h: win.h, proposalScore: win.proposal, validationScore: win.validation, validationRunnerUp: runner, visibleProposal: win.visibleProposal, visibleValidation: win.visibleValidation, hypothesesEvaluated: hypotheses}
}

func maxIntV4Placement(a, b int) int {
	if a > b {
		return a
	}
	return b
}
