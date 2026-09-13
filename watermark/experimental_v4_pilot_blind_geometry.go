package watermark

import (
	"image"
	"math"
	"sort"
)

// Build26 keeps the v4 geometry search experimental and completely detached
// from the production v3 decoder. The search is intentionally bounded and
// deterministic: it asks whether the public pilot can rank plausible geometry
// hypotheses without payload, header, ECC or HMAC evidence.

type experimentalV4BlindGeometryParams struct {
	angleDeg    float64
	scaleX      float64
	scaleY      float64
	shearXDeg   float64
	shearYDeg   float64
	topInset    float64
	bottomInset float64
	shiftX      float64
	shiftY      float64
}

type experimentalV4BlindGeometryResult struct {
	available           bool
	params              experimentalV4BlindGeometryParams
	h                   homography
	detection           ExperimentalV4PilotDetection
	objective           float64
	runnerUpObjective   float64
	hypothesesEvaluated int
}

type experimentalV4BlindProposal struct {
	params         experimentalV4BlindGeometryParams
	dimensionError float64
}

var experimentalV4Build26DeformationStates = [...]experimentalV4BlindGeometryParams{
	{scaleX: 1.00, scaleY: 1.00},
	{scaleX: 1.10, scaleY: 0.90},
	{scaleX: 0.90, scaleY: 1.10},
	{scaleX: 0.80, scaleY: 0.80},
	{scaleX: 1.00, scaleY: 1.00, shearXDeg: 8},
	{scaleX: 1.00, scaleY: 1.00, shearXDeg: -8},
	{scaleX: 1.00, scaleY: 1.00, shearYDeg: 8},
	{scaleX: 1.00, scaleY: 1.00, shearYDeg: -8},
}

// experimentalV4BlindGeometrySearch performs a bounded pilot-only geometry
// search. canonicalWidth/canonicalHeight describe the pre-transform carrier
// extent; the observed image itself supplies all pilot evidence. The current
// Build26 search assumes auto-framed transforms without arbitrary crop/offset.
type experimentalV4BlindScoredHypothesis struct {
	params        experimentalV4BlindGeometryParams
	h             homography
	proposalScore float64
	objective     float64
	detection     ExperimentalV4PilotDetection
}

func experimentalV4BlindGeometrySearch(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int) experimentalV4BlindGeometryResult {
	if img == nil || canonicalWidth < experimentalV4TileWidthBlocks*blockSize || canonicalHeight < experimentalV4TileHeightBlocks*blockSize {
		return experimentalV4BlindGeometryResult{}
	}
	observedWidth, observedHeight := img.Bounds().Dx(), img.Bounds().Dy()
	plane := newPixelPlane(img)
	hypotheses := 0

	// Build26 deliberately separates evidence:
	//   1. repeated data-plane self-consistency proposes geometry without knowing
	//      any data symbols, payload, key or pilot signs;
	//   2. the public pilot ranks a small bounded proposal bank using central rows;
	//   3. spatially held-out pilot repetitions validate the surviving geometry.
	// This mirrors the intended v4 coarse-lattice -> pilot -> origin pipeline.
	type repeatHypothesis struct {
		params experimentalV4BlindGeometryParams
		h      homography
		score  float64
		dimErr float64
	}
	repeatBeam := make([]repeatHypothesis, 0, 24)
	repeatSeen := make(map[[7]int]bool)
	keyForRepeat := func(p experimentalV4BlindGeometryParams) [7]int {
		return [7]int{
			int(math.Round(p.angleDeg * 1000)), int(math.Round(p.scaleX * 1000)), int(math.Round(p.scaleY * 1000)),
			int(math.Round(p.shearXDeg * 1000)), int(math.Round(p.shearYDeg * 1000)),
			int(math.Round(p.topInset * 10000)), int(math.Round(p.bottomInset * 10000)),
		}
	}
	evaluateRepeat := func(p experimentalV4BlindGeometryParams) (repeatHypothesis, bool) {
		p.shiftX, p.shiftY = 0, 0
		h, pw, ph := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
		dimErr := experimentalV4BlindDimensionError(pw, ph, observedWidth, observedHeight)
		score, pairs := experimentalV4DataRepeatConsistencyPlane(plane, candidate, canonicalWidth, canonicalHeight, h)
		hypotheses++
		if pairs < 24 || math.IsInf(score, -1) {
			return repeatHypothesis{}, false
		}
		return repeatHypothesis{params: p, h: h, score: score - 0.20*dimErr, dimErr: dimErr}, true
	}
	insertRepeat := func(p experimentalV4BlindGeometryParams) {
		p.shiftX, p.shiftY = 0, 0
		key := keyForRepeat(p)
		if repeatSeen[key] {
			return
		}
		repeatSeen[key] = true
		hyp, ok := evaluateRepeat(p)
		if !ok {
			return
		}
		repeatBeam = append(repeatBeam, hyp)
		sort.Slice(repeatBeam, func(i, j int) bool { return repeatBeam[i].score > repeatBeam[j].score })
		if len(repeatBeam) > 24 {
			repeatBeam = repeatBeam[:24]
		}
	}

	// Affine scale/rotation branch. The canvas solves scale cheaply for every
	// 0.1-degree angle; repeat consistency decides which basin is real.
	for angle := -20.0; angle <= 20.0001; angle += 0.1 {
		insertRepeat(experimentalV4BestScaleForAngle(canonicalWidth, canonicalHeight, observedWidth, observedHeight, angle).params)
		insertRepeat(experimentalV4BlindGeometryParams{angleDeg: angle, scaleX: 1, scaleY: 1})
	}

	// Shear branches. For every fine angle, dimension-only enumeration chooses
	// one shear magnitude; the repeated data plane then ranks that proposal.
	for axis := 0; axis < 2; axis++ {
		for angle := -20.0; angle <= 20.0001; angle += 0.1 {
			best := experimentalV4BlindProposal{dimensionError: math.Inf(1)}
			for shear := -10.0; shear <= 10.0001; shear += 0.5 {
				p := experimentalV4BlindGeometryParams{angleDeg: angle, scaleX: 1, scaleY: 1}
				if axis == 0 {
					p.shearXDeg = shear
				} else {
					p.shearYDeg = shear
				}
				_, w, h := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
				err := experimentalV4BlindDimensionError(w, h, observedWidth, observedHeight)
				if err < best.dimensionError {
					best = experimentalV4BlindProposal{params: p, dimensionError: err}
				}
			}
			insertRepeat(best.params)
		}
	}

	// Mild projective branch. Insets remain a small public bounded family while
	// angle is fine enough for the 8-pixel DCT carrier. Scale is solved from the
	// canvas for each (angle, top, bottom) tuple before any image-domain scoring.
	perspectiveInsets := []float64{0, 0.015, 0.030, 0.045, 0.060}
	for _, top := range perspectiveInsets {
		for _, bottom := range perspectiveInsets {
			best := repeatHypothesis{score: math.Inf(-1)}
			// Coarse angle scan stays inside this projective family so a valid
			// family cannot be pruned by an unrelated false maximum.
			for angle := -20.0; angle <= 20.0001; angle += 0.25 {
				proposal := experimentalV4BestScaleForAnglePerspective(canonicalWidth, canonicalHeight, observedWidth, observedHeight, angle, top, bottom)
				if hyp, ok := evaluateRepeat(proposal.params); ok && hyp.score > best.score {
					best = hyp
				}
			}
			if math.IsInf(best.score, -1) {
				continue
			}
			baseAngle := best.params.angleDeg
			for da := -0.25; da <= 0.2501; da += 0.05 {
				angle := baseAngle + da
				proposal := experimentalV4BestScaleForAnglePerspective(canonicalWidth, canonicalHeight, observedWidth, observedHeight, angle, top, bottom)
				if hyp, ok := evaluateRepeat(proposal.params); ok && hyp.score > best.score {
					best = hyp
				}
			}
			insertRepeat(best.params)
		}
	}

	if len(repeatBeam) == 0 {
		return experimentalV4BlindGeometryResult{hypothesesEvaluated: hypotheses}
	}

	// Local repeat-only refinement. This stage is still symbol-blind. It sharpens
	// the geometry before the pilot is allowed to see any candidate.
	coarseSeeds := append([]repeatHypothesis(nil), repeatBeam...)
	if len(coarseSeeds) > 10 {
		coarseSeeds = coarseSeeds[:10]
	}
	for _, seed := range coarseSeeds {
		base := seed.params
		for da := -0.10; da <= 0.1001; da += 0.05 {
			p := base
			p.angleDeg = base.angleDeg + da
			insertRepeat(p)
		}
		if math.Abs(base.shearXDeg) > 1e-9 {
			for ds := -0.5; ds <= 0.5001; ds += 0.25 {
				p := base
				p.shearXDeg = base.shearXDeg + ds
				insertRepeat(p)
			}
		} else if math.Abs(base.shearYDeg) > 1e-9 {
			for ds := -0.5; ds <= 0.5001; ds += 0.25 {
				p := base
				p.shearYDeg = base.shearYDeg + ds
				insertRepeat(p)
			}
		}
		if base.topInset > 0 || base.bottomInset > 0 {
			for dt := -0.005; dt <= 0.0051; dt += 0.005 {
				for db := -0.005; db <= 0.0051; db += 0.005 {
					p := base
					p.topInset = math.Max(0, base.topInset+dt)
					p.bottomInset = math.Max(0, base.bottomInset+db)
					insertRepeat(p)
				}
			}
			for dx := -0.005; dx <= 0.0051; dx += 0.005 {
				for dy := -0.005; dy <= 0.0051; dy += 0.005 {
					p := base
					p.scaleX = clampV4Scale(base.scaleX + dx)
					p.scaleY = clampV4Scale(base.scaleY + dy)
					insertRepeat(p)
				}
			}
		}
	}

	// Public-pilot ranking on the bounded repeat-proposed bank. A small global
	// translation bank absorbs auto-bound/subpixel rounding without conflating it
	// with cyclic tile origin.
	type pilotHypothesis struct {
		params        experimentalV4BlindGeometryParams
		h             homography
		proposalScore float64
		repeatScore   float64
	}
	pilotBeam := make([]pilotHypothesis, 0, 12)
	phaseOffsets := []float64{-4, -2, 0, 2, 4}
	insertPilot := func(seed repeatHypothesis) {
		bestScore := math.Inf(-1)
		bestP := seed.params
		bestH := seed.h
		for _, dx := range phaseOffsets {
			for _, dy := range phaseOffsets {
				p := seed.params
				p.shiftX, p.shiftY = dx, dy
				h, _, _ := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
				score, visible := experimentalV4PilotFixedOriginScorePlane(plane, candidate, canonicalWidth, canonicalHeight, h)
				if visible >= 56 && score > bestScore {
					bestScore, bestP, bestH = score, p, h
				}
			}
		}
		if math.IsInf(bestScore, -1) {
			return
		}
		pilotBeam = append(pilotBeam, pilotHypothesis{params: bestP, h: bestH, proposalScore: bestScore, repeatScore: seed.score})
		sort.Slice(pilotBeam, func(i, j int) bool {
			if pilotBeam[i].proposalScore == pilotBeam[j].proposalScore {
				return pilotBeam[i].repeatScore > pilotBeam[j].repeatScore
			}
			return pilotBeam[i].proposalScore > pilotBeam[j].proposalScore
		})
		if len(pilotBeam) > 12 {
			pilotBeam = pilotBeam[:12]
		}
	}
	for _, seed := range repeatBeam {
		insertPilot(seed)
	}
	if len(pilotBeam) == 0 {
		return experimentalV4BlindGeometryResult{hypothesesEvaluated: hypotheses}
	}

	// Pilot-guided fine geometry refinement around only the best few proposals.
	// This is still proposal evidence; the final decision below uses held-out
	// pilot repetitions that were not sampled here.
	refinePilot := func(base pilotHypothesis, p experimentalV4BlindGeometryParams) pilotHypothesis {
		best := base
		for _, dx := range []float64{base.params.shiftX - 1, base.params.shiftX, base.params.shiftX + 1} {
			for _, dy := range []float64{base.params.shiftY - 1, base.params.shiftY, base.params.shiftY + 1} {
				q := p
				q.shiftX, q.shiftY = dx, dy
				h, _, _ := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, q)
				score, visible := experimentalV4PilotFixedOriginScorePlane(plane, candidate, canonicalWidth, canonicalHeight, h)
				if visible >= 56 && score > best.proposalScore {
					best = pilotHypothesis{params: q, h: h, proposalScore: score, repeatScore: base.repeatScore}
				}
			}
		}
		return best
	}
	refined := append([]pilotHypothesis(nil), pilotBeam...)
	seedCount := minInt(5, len(pilotBeam))
	for i := 0; i < seedCount; i++ {
		base := pilotBeam[i]
		winner := base
		for da := -0.10; da <= 0.1001; da += 0.05 {
			p := base.params
			p.angleDeg = base.params.angleDeg + da
			winner = refinePilot(winner, p)
		}
		if math.Abs(base.params.shearXDeg) > 1e-9 {
			for ds := -0.25; ds <= 0.2501; ds += 0.125 {
				p := base.params
				p.shearXDeg = base.params.shearXDeg + ds
				winner = refinePilot(winner, p)
			}
		} else if math.Abs(base.params.shearYDeg) > 1e-9 {
			for ds := -0.25; ds <= 0.2501; ds += 0.125 {
				p := base.params
				p.shearYDeg = base.params.shearYDeg + ds
				winner = refinePilot(winner, p)
			}
		}
		if base.params.topInset > 0 || base.params.bottomInset > 0 {
			for dt := -0.005; dt <= 0.0051; dt += 0.0025 {
				for db := -0.005; db <= 0.0051; db += 0.0025 {
					p := base.params
					p.topInset = math.Max(0, base.params.topInset+dt)
					p.bottomInset = math.Max(0, base.params.bottomInset+db)
					winner = refinePilot(winner, p)
				}
			}
		}
		refined = append(refined, winner)
	}
	sort.Slice(refined, func(i, j int) bool { return refined[i].proposalScore > refined[j].proposalScore })
	if len(refined) > 12 {
		refined = refined[:12]
	}

	// Final held-out validation. The validation objective contains no repeat
	// score and no central-row pilot score: those only proposed candidates.
	validated := make([]experimentalV4BlindScoredHypothesis, 0, len(refined))
	validatedSeen := make(map[[9]int]bool)
	validationKey := func(p experimentalV4BlindGeometryParams) [9]int {
		return [9]int{
			int(math.Round(p.angleDeg * 1000)), int(math.Round(p.scaleX * 1000)), int(math.Round(p.scaleY * 1000)),
			int(math.Round(p.shearXDeg * 1000)), int(math.Round(p.shearYDeg * 1000)),
			int(math.Round(p.topInset * 10000)), int(math.Round(p.bottomInset * 10000)),
			int(math.Round(p.shiftX * 1000)), int(math.Round(p.shiftY * 1000)),
		}
	}
	for _, proposal := range refined {
		key := validationKey(proposal.params)
		if validatedSeen[key] {
			continue
		}
		validatedSeen[key] = true
		detection := experimentalV4DetectPilotProjectiveValidationPlane(plane, candidate, canonicalWidth, canonicalHeight, proposal.h)
		if !detection.Available {
			continue
		}
		_, pw, ph := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, proposal.params)
		dimErr := experimentalV4BlindDimensionError(pw, ph, observedWidth, observedHeight)
		validated = append(validated, experimentalV4BlindScoredHypothesis{
			params: proposal.params, h: proposal.h, proposalScore: proposal.proposalScore,
			detection: detection, objective: experimentalV4BlindObjective(detection) - 0.20*dimErr,
		})
	}
	if len(validated) == 0 {
		return experimentalV4BlindGeometryResult{hypothesesEvaluated: hypotheses}
	}
	sort.Slice(validated, func(i, j int) bool { return validated[i].objective > validated[j].objective })
	winner := validated[0]
	runner := 0.0
	if len(validated) > 1 {
		runner = validated[1].objective
	}
	return experimentalV4BlindGeometryResult{
		available: true, params: winner.params, h: winner.h, detection: winner.detection,
		objective: winner.objective, runnerUpObjective: runner, hypothesesEvaluated: hypotheses,
	}
}

// experimentalV4DataRepeatConsistencyPlane is a key/payload/pilot-symbol
// independent coarse lattice observable. It samples only data-plane coordinates
// and compares DCT signs at homologous positions of adjacent repeated v4 tiles.
// The actual symbol value is never consulted: only self-consistency matters.
func experimentalV4DataRepeatConsistencyPlane(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) (float64, int) {
	if plane == nil {
		return math.Inf(-1), 0
	}
	blocksWide, blocksHigh := canonicalWidth/blockSize, canonicalHeight/blockSize
	tilesX := blocksWide / experimentalV4TileWidthBlocks
	tilesY := blocksHigh / experimentalV4TileHeightBlocks
	if tilesX < 2 && tilesY < 2 {
		return math.Inf(-1), 0
	}

	pilot := make(map[int]bool, experimentalV4PilotCount)
	for _, position := range candidate.positions {
		pilot[position] = true
	}
	// 64 deterministic data-plane residues, distributed pseudo-uniformly. The
	// permutation step is coprime with 1184, so no residue repeats.
	residues := make([]int, 0, 64)
	for k := 0; k < experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks && len(residues) < 64; k++ {
		position := (k*73 + 17) % (experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks)
		if !pilot[position] {
			residues = append(residues, position)
		}
	}

	numerator, denominator := 0.0, 0.0
	pairs := 0
	compare := func(ax, ay, bx, by int) {
		for _, position := range residues {
			x := position % experimentalV4TileWidthBlocks
			y := position / experimentalV4TileWidthBlocks
			va, oka := readProjectiveBlockValue(plane, h, (ax*experimentalV4TileWidthBlocks+x)*blockSize, (ay*experimentalV4TileHeightBlocks+y)*blockSize, blockSize)
			vb, okb := readProjectiveBlockValue(plane, h, (bx*experimentalV4TileWidthBlocks+x)*blockSize, (by*experimentalV4TileHeightBlocks+y)*blockSize, blockSize)
			if !oka || !okb {
				continue
			}
			weight := math.Min(math.Abs(va), math.Abs(vb))
			if weight <= 0 {
				continue
			}
			if va*vb >= 0 {
				numerator += weight
			} else {
				numerator -= weight
			}
			denominator += weight
			pairs++
		}
	}
	// Adjacent central pairs deliberately reduce sensitivity to coarse angle
	// error; held-out pilot validation later spans distant corner tiles.
	if tilesX >= 2 {
		row := tilesY / 2
		left := (tilesX - 2) / 2
		compare(left, row, left+1, row)
	}
	if tilesY >= 2 {
		col := tilesX / 2
		top := (tilesY - 2) / 2
		compare(col, top, col, top+1)
	}
	if denominator <= 0 {
		return math.Inf(-1), pairs
	}
	return numerator / denominator, pairs
}

func experimentalV4PilotFixedOriginScore(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) (float64, int) {
	if img == nil {
		return math.Inf(-1), 0
	}
	return experimentalV4PilotFixedOriginScorePlane(newPixelPlane(img), candidate, canonicalWidth, canonicalHeight, h)
}

func experimentalV4PilotFixedOriginScorePlane(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) (float64, int) {
	if plane == nil {
		return math.Inf(-1), 0
	}
	blocksWide, blocksHigh := canonicalWidth/blockSize, canonicalHeight/blockSize
	if blocksWide < experimentalV4TileWidthBlocks || blocksHigh < experimentalV4TileHeightBlocks {
		return math.Inf(-1), 0
	}
	maxStartX := ((blocksWide - experimentalV4TileWidthBlocks) / experimentalV4TileWidthBlocks) * experimentalV4TileWidthBlocks
	maxStartY := ((blocksHigh - experimentalV4TileHeightBlocks) / experimentalV4TileHeightBlocks) * experimentalV4TileHeightBlocks
	centerX := (maxStartX / (2 * experimentalV4TileWidthBlocks)) * experimentalV4TileWidthBlocks
	centerY := (maxStartY / (2 * experimentalV4TileHeightBlocks)) * experimentalV4TileHeightBlocks

	// Proposal evidence uses the middle tile row only. Held-out validation uses
	// the four corner tiles, so the two stages remain spatially disjoint whenever
	// the carrier has at least three tile rows (the Build26 qualification case).
	starts := [][2]int{{0, centerY}, {centerX, centerY}, {maxStartX, centerY}}
	unique := make(map[[2]int]bool)
	numerator, denominator := 0.0, 0.0
	visibleSamples := 0
	for _, start := range starts {
		if unique[start] {
			continue
		}
		unique[start] = true
		for i, pos := range candidate.positions {
			x := pos % experimentalV4TileWidthBlocks
			y := pos / experimentalV4TileWidthBlocks
			v, ok := readProjectiveBlockValue(plane, h, (start[0]+x)*blockSize, (start[1]+y)*blockSize, blockSize)
			if !ok {
				continue
			}
			visibleSamples++
			numerator += float64(candidate.signs[i]) * v
			denominator += math.Abs(v)
		}
	}
	if denominator <= 0 {
		return math.Inf(-1), visibleSamples
	}
	return numerator / denominator, visibleSamples
}

func experimentalV4BestScaleForAngle(canonicalWidth, canonicalHeight, observedWidth, observedHeight int, angle float64) experimentalV4BlindProposal {
	return experimentalV4BestScaleForAnglePerspective(canonicalWidth, canonicalHeight, observedWidth, observedHeight, angle, 0, 0)
}

func experimentalV4BestScaleForAnglePerspective(canonicalWidth, canonicalHeight, observedWidth, observedHeight int, angle, topInset, bottomInset float64) experimentalV4BlindProposal {
	best := experimentalV4BlindProposal{dimensionError: math.Inf(1)}
	evaluate := func(sx, sy float64) {
		p := experimentalV4BlindGeometryParams{angleDeg: angle, scaleX: sx, scaleY: sy, topInset: topInset, bottomInset: bottomInset}
		_, w, h := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
		err := experimentalV4BlindDimensionError(w, h, observedWidth, observedHeight)
		if err < best.dimensionError {
			best = experimentalV4BlindProposal{params: p, dimensionError: err}
		}
	}
	for sx := 0.75; sx <= 1.1501; sx += 0.025 {
		for sy := 0.75; sy <= 1.1501; sy += 0.025 {
			evaluate(sx, sy)
		}
	}
	coarseX, coarseY := best.params.scaleX, best.params.scaleY
	for sx := math.Max(0.70, coarseX-0.03); sx <= math.Min(1.20, coarseX+0.03)+1e-9; sx += 0.005 {
		for sy := math.Max(0.70, coarseY-0.03); sy <= math.Min(1.20, coarseY+0.03)+1e-9; sy += 0.005 {
			evaluate(sx, sy)
		}
	}
	return best
}

func experimentalV4SortDimensionProposals(proposals []experimentalV4BlindProposal) {
	sort.Slice(proposals, func(i, j int) bool {
		if proposals[i].dimensionError == proposals[j].dimensionError {
			return proposals[i].params.angleDeg < proposals[j].params.angleDeg
		}
		return proposals[i].dimensionError < proposals[j].dimensionError
	})
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func experimentalV4BlindDimensionError(predictedWidth, predictedHeight, observedWidth, observedHeight int) float64 {
	if predictedWidth <= 0 || predictedHeight <= 0 || observedWidth <= 0 || observedHeight <= 0 {
		return math.Inf(1)
	}
	dw := math.Abs(float64(predictedWidth-observedWidth)) / float64(observedWidth)
	dh := math.Abs(float64(predictedHeight-observedHeight)) / float64(observedHeight)
	return dw + dh
}

func clampV4Scale(value float64) float64 {
	if value < 0.70 {
		return 0.70
	}
	if value > 1.20 {
		return 1.20
	}
	return value
}

func experimentalV4BlindObjective(d ExperimentalV4PilotDetection) float64 {
	visibleWeight := float64(d.VisiblePilotPositions) / float64(experimentalV4PilotCount)
	return d.Score + 1.50*d.Margin + 0.05*visibleWeight
}

// experimentalV4BlindHomography constructs the canonical->observed mapping for
// one bounded Build26 hypothesis. Affine order is scale -> shear -> rotation,
// matching the Build25 synthetic families, followed by a mild projective warp.
func experimentalV4BlindHomography(width, height int, p experimentalV4BlindGeometryParams) (homography, int, int) {
	sx, sy := p.scaleX, p.scaleY
	if sx == 0 {
		sx = 1
	}
	if sy == 0 {
		sy = 1
	}
	shx := math.Tan(p.shearXDeg * math.Pi / 180)
	shy := math.Tan(p.shearYDeg * math.Pi / 180)
	angle := p.angleDeg * math.Pi / 180
	c, s := math.Cos(angle), math.Sin(angle)

	// shear*scale
	m00 := sx
	m01 := shx * sy
	m10 := shy * sx
	m11 := sy
	// rotation*(shear*scale)
	a := c*m00 - s*m10
	b := c*m01 - s*m11
	cc := s*m00 + c*m10
	d := s*m01 + c*m11

	affine, outW, outH := experimentalV4BlindAutoBoundAffine(width, height, a, b, cc, d)
	if p.topInset <= 0 && p.bottomInset <= 0 {
		if p.shiftX != 0 || p.shiftY != 0 {
			shift := homography{h: [9]float64{1, 0, p.shiftX, 0, 1, p.shiftY, 0, 0, 1}}
			affine = experimentalV4BlindMultiplyHomography(shift, affine)
		}
		return affine, outW, outH
	}

	top := math.Max(0, math.Min(0.12, p.topInset))
	bottom := math.Max(0, math.Min(0.12, p.bottomInset))
	dst := [4][2]float64{
		{top * float64(outW-1), 0},
		{(1 - top) * float64(outW-1), 0},
		{bottom * float64(outW-1), float64(outH - 1)},
		{(1 - bottom) * float64(outW-1), float64(outH - 1)},
	}
	src := [4][2]float64{{0, 0}, {float64(outW - 1), 0}, {0, float64(outH - 1)}, {float64(outW - 1), float64(outH - 1)}}
	perspective, ok := diagnosticHomographyFromFourPoints(src, dst)
	if !ok {
		return affine, outW, outH
	}
	combined := experimentalV4BlindMultiplyHomography(perspective, affine)
	if p.shiftX != 0 || p.shiftY != 0 {
		shift := homography{h: [9]float64{1, 0, p.shiftX, 0, 1, p.shiftY, 0, 0, 1}}
		combined = experimentalV4BlindMultiplyHomography(shift, combined)
	}
	return combined, outW, outH
}

func experimentalV4BlindAutoBoundAffine(width, height int, a, b, c, d float64) (homography, int, int) {
	base := homography{h: [9]float64{a, b, 0, c, d, 0, 0, 0, 1}}
	corners := [4][2]float64{{0, 0}, {float64(width - 1), 0}, {0, float64(height - 1)}, {float64(width - 1), float64(height - 1)}}
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, point := range corners {
		x, y, _ := base.mapPoint(point[0], point[1])
		minX = math.Min(minX, x)
		minY = math.Min(minY, y)
		maxX = math.Max(maxX, x)
		maxY = math.Max(maxY, y)
	}
	const margin = 3.0
	translate := homography{h: [9]float64{1, 0, -minX + margin, 0, 1, -minY + margin, 0, 0, 1}}
	h := experimentalV4BlindMultiplyHomography(translate, base)
	outW := int(math.Ceil(maxX - minX + 1 + 2*margin))
	outH := int(math.Ceil(maxY - minY + 1 + 2*margin))
	return h, outW, outH
}

func experimentalV4BlindMultiplyHomography(a, b homography) homography {
	var out [9]float64
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			for k := 0; k < 3; k++ {
				out[row*3+col] += a.h[row*3+k] * b.h[k*3+col]
			}
		}
	}
	return homography{h: out}
}

// Sparse scoring samples one central 37x32 tile instead of every repeated tile.
// This keeps geometry proposal bounded while preserving every pilot residue.
func experimentalV4DetectPilotProjectiveSparse(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) ExperimentalV4PilotDetection {
	if img == nil {
		return ExperimentalV4PilotDetection{}
	}
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	if blocksWide < experimentalV4TileWidthBlocks || blocksHigh < experimentalV4TileHeightBlocks {
		return ExperimentalV4PilotDetection{}
	}

	startBlockX := ((blocksWide - experimentalV4TileWidthBlocks) / 2 / experimentalV4TileWidthBlocks) * experimentalV4TileWidthBlocks
	startBlockY := ((blocksHigh - experimentalV4TileHeightBlocks) / 2 / experimentalV4TileHeightBlocks) * experimentalV4TileHeightBlocks
	plane := newPixelPlane(img)
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var values [residues]float64
	var valid [residues]bool
	for y := 0; y < experimentalV4TileHeightBlocks; y++ {
		for x := 0; x < experimentalV4TileWidthBlocks; x++ {
			value, ok := readProjectiveBlockValue(plane, h, (startBlockX+x)*blockSize, (startBlockY+y)*blockSize, blockSize)
			if !ok {
				continue
			}
			index := y*experimentalV4TileWidthBlocks + x
			values[index] = value
			valid[index] = true
		}
	}

	bestScore := math.Inf(-1)
	runnerUp := math.Inf(-1)
	bestX, bestY := 0, 0
	bestVisible := 0
	for originY := 0; originY < experimentalV4TileHeightBlocks; originY++ {
		for originX := 0; originX < experimentalV4TileWidthBlocks; originX++ {
			numerator, denominator := 0.0, 0.0
			visible := 0
			for i, position := range candidate.positions {
				px := position % experimentalV4TileWidthBlocks
				py := position / experimentalV4TileWidthBlocks
				ox := positiveMod(px-originX, experimentalV4TileWidthBlocks)
				oy := positiveMod(py-originY, experimentalV4TileHeightBlocks)
				index := oy*experimentalV4TileWidthBlocks + ox
				if !valid[index] {
					continue
				}
				visible++
				numerator += float64(candidate.signs[i]) * values[index]
				denominator += math.Abs(values[index])
			}
			if denominator <= 0 {
				continue
			}
			score := numerator / denominator
			if score > bestScore {
				runnerUp = bestScore
				bestScore = score
				bestX, bestY = originX, originY
				bestVisible = visible
			} else if score > runnerUp {
				runnerUp = score
			}
		}
	}
	if math.IsInf(bestScore, -1) {
		return ExperimentalV4PilotDetection{}
	}
	if math.IsInf(runnerUp, -1) {
		runnerUp = 0
	}
	return ExperimentalV4PilotDetection{
		Available:             true,
		BlockSizePixels:       blockSize,
		OriginXBlocks:         bestX,
		OriginYBlocks:         bestY,
		Score:                 bestScore,
		RunnerUpScore:         runnerUp,
		Margin:                bestScore - runnerUp,
		VisiblePilotPositions: bestVisible,
		PilotSamples:          bestVisible,
	}
}

// experimentalV4DetectPilotProjectiveValidation aggregates up to five complete
// tile repetitions that are spatially separated from the single central proposal
// tile. It is intentionally bounded and remains pilot-only evidence.
func experimentalV4DetectPilotProjectiveValidation(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) ExperimentalV4PilotDetection {
	if img == nil {
		return ExperimentalV4PilotDetection{}
	}
	return experimentalV4DetectPilotProjectiveValidationPlane(newPixelPlane(img), candidate, canonicalWidth, canonicalHeight, h)
}

func experimentalV4DetectPilotProjectiveValidationPlane(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) ExperimentalV4PilotDetection {
	if plane == nil {
		return ExperimentalV4PilotDetection{}
	}
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	if blocksWide < experimentalV4TileWidthBlocks || blocksHigh < experimentalV4TileHeightBlocks {
		return ExperimentalV4PilotDetection{}
	}
	maxStartX := ((blocksWide - experimentalV4TileWidthBlocks) / experimentalV4TileWidthBlocks) * experimentalV4TileWidthBlocks
	maxStartY := ((blocksHigh - experimentalV4TileHeightBlocks) / experimentalV4TileHeightBlocks) * experimentalV4TileHeightBlocks
	starts := [][2]int{{0, 0}, {maxStartX, 0}, {0, maxStartY}, {maxStartX, maxStartY}}
	unique := make(map[[2]int]bool)
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var sums [residues]float64
	var absSums [residues]float64
	var counts [residues]int
	for _, start := range starts {
		if unique[start] {
			continue
		}
		unique[start] = true
		for y := 0; y < experimentalV4TileHeightBlocks; y++ {
			for x := 0; x < experimentalV4TileWidthBlocks; x++ {
				value, ok := readProjectiveBlockValue(plane, h, (start[0]+x)*blockSize, (start[1]+y)*blockSize, blockSize)
				if !ok {
					continue
				}
				index := y*experimentalV4TileWidthBlocks + x
				sums[index] += value
				absSums[index] += math.Abs(value)
				counts[index]++
			}
		}
	}
	return experimentalV4DetectPilotFromResidues(candidate, sums[:], absSums[:], counts[:])
}

func experimentalV4DetectPilotFromResidues(candidate experimentalV4PilotCandidate, sums, absSums []float64, counts []int) ExperimentalV4PilotDetection {
	bestScore := math.Inf(-1)
	runnerUp := math.Inf(-1)
	bestX, bestY := 0, 0
	bestVisible, bestSamples := 0, 0
	for originY := 0; originY < experimentalV4TileHeightBlocks; originY++ {
		for originX := 0; originX < experimentalV4TileWidthBlocks; originX++ {
			numerator, denominator := 0.0, 0.0
			visible, samples := 0, 0
			for i, position := range candidate.positions {
				px := position % experimentalV4TileWidthBlocks
				py := position / experimentalV4TileWidthBlocks
				ox := positiveMod(px-originX, experimentalV4TileWidthBlocks)
				oy := positiveMod(py-originY, experimentalV4TileHeightBlocks)
				index := oy*experimentalV4TileWidthBlocks + ox
				if index >= len(counts) || counts[index] == 0 {
					continue
				}
				visible++
				samples += counts[index]
				numerator += float64(candidate.signs[i]) * sums[index]
				denominator += absSums[index]
			}
			if denominator <= 0 {
				continue
			}
			score := numerator / denominator
			if score > bestScore {
				runnerUp = bestScore
				bestScore = score
				bestX, bestY = originX, originY
				bestVisible, bestSamples = visible, samples
			} else if score > runnerUp {
				runnerUp = score
			}
		}
	}
	if math.IsInf(bestScore, -1) {
		return ExperimentalV4PilotDetection{}
	}
	if math.IsInf(runnerUp, -1) {
		runnerUp = 0
	}
	return ExperimentalV4PilotDetection{Available: true, BlockSizePixels: blockSize, OriginXBlocks: bestX, OriginYBlocks: bestY, Score: bestScore, RunnerUpScore: runnerUp, Margin: bestScore - runnerUp, VisiblePilotPositions: bestVisible, PilotSamples: bestSamples}
}
