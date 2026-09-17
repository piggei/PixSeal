package watermark

import (
	"math"
	"sort"
)

// Build40 adds a deliberately small, smooth residual field on top of the
// Build39 smartphone homography. The model is fitted from checkerboard-A
// repetitions of the public v4 pilot only; checkerboard-B repetitions remain
// held out until validation. Payload/header/ECC/HMAC never participate in
// fitting or acceptance.
const (
	experimentalV4PhoneResidualMaxPixels       = 6.0
	experimentalV4PhoneResidualMinControls     = 7
	experimentalV4PhoneResidualMinLocalScore   = 0.22
	experimentalV4PhoneResidualMinLocalGap     = 0.025
	experimentalV4PhoneResidualMaxFitRMS       = 2.25
	experimentalV4PhoneResidualValidationFloor = 0.11
	experimentalV4PhoneResidualValidationGain  = 0.015
	experimentalV4PhoneResidualProposalFloor   = 0.14
	experimentalV4PhoneResidualFullScoreFloor  = 0.12
	experimentalV4PhoneResidualFullMarginFloor = 0.012
)

type experimentalV4PhoneResidualWarp struct {
	width, height float64
	dx            [6]float64
	dy            [6]float64
	controls      int
	rmsPixels     float64
	maxPixels     float64
}

type experimentalV4PhoneResidualControl struct {
	x, y       float64
	dx, dy     float64
	weight     float64
	localScore float64
	gap        float64
}

type experimentalV4PhoneResidualInfo struct {
	Applied             bool
	Controls            int
	RMSPixels           float64
	ProposalBefore      float64
	ProposalAfter       float64
	ValidationBefore    float64
	ValidationAfter     float64
	FullPilotScore      float64
	FullPilotMargin     float64
	FullOriginXBlocks   int
	FullOriginYBlocks   int
	HypothesesEvaluated int
}

func (warp *experimentalV4PhoneResidualWarp) correction(x, y float64) (float64, float64) {
	if warp == nil || warp.width <= 0 || warp.height <= 0 {
		return 0, 0
	}
	nx := 2*x/warp.width - 1
	ny := 2*y/warp.height - 1
	basis := [6]float64{1, nx, ny, nx * ny, nx * nx, ny * ny}
	dx, dy := 0.0, 0.0
	for i := range basis {
		dx += warp.dx[i] * basis[i]
		dy += warp.dy[i] * basis[i]
	}
	limit := warp.maxPixels
	if limit <= 0 || limit > experimentalV4PhoneResidualMaxPixels {
		limit = experimentalV4PhoneResidualMaxPixels
	}
	dx = math.Max(-limit, math.Min(limit, dx))
	dy = math.Max(-limit, math.Min(limit, dy))
	return dx, dy
}

func experimentalV4PhoneWarpedBlockValue(src *pixelPlane, h homography, warp *experimentalV4PhoneResidualWarp, canonicalWidth, canonicalHeight, originX, originY, size int) (float64, bool) {
	table := readCosTables[size]
	c23, c32 := 0.0, 0.0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cx := float64(originX + x)
			cy := float64(originY + y)
			if warp != nil {
				dx, dy := warp.correction(cx, cy)
				cx += dx
				cy += dy
			}
			sx, sy, ok := h.mapPoint(cx, cy)
			if !ok {
				return 0, false
			}
			l, ok := samplePlaneLuminance(src, sx, sy)
			if !ok {
				return 0, false
			}
			c23 += l * table[3][x] * table[2][y]
			c32 += l * table[2][x] * table[3][y]
		}
	}
	return math.Abs(c23) - math.Abs(c32), true
}

// fixed-origin checkerboard score. Build39 has already chosen/canonicalized the
// absolute pilot origin; Build40 is allowed to move only within one residual
// sub-block basin, so it must not reopen the 37x32 origin search while fitting.
func experimentalV4PhoneResidualSpatialScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, warp *experimentalV4PhoneResidualWarp, partition int) (float64, int) {
	if plane == nil {
		return math.Inf(-1), 0
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX := bw / experimentalV4TileWidthBlocks
	tilesY := bh / experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return math.Inf(-1), 0
	}
	num, den := 0.0, 0.0
	visible := 0
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			if (tx+ty)&1 != partition {
				continue
			}
			baseX := tx * experimentalV4TileWidthBlocks
			baseY := ty * experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				px := pos % experimentalV4TileWidthBlocks
				py := pos / experimentalV4TileWidthBlocks
				v, ok := experimentalV4PhoneWarpedBlockValue(plane, h, warp, cw, ch, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
				if !ok {
					continue
				}
				visible++
				num += float64(candidate.signs[i]) * v
				den += math.Abs(v)
			}
		}
	}
	if den <= 0 {
		return math.Inf(-1), visible
	}
	return num / den, visible
}

func experimentalV4PhoneResidualTileScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, tx, ty int, dx, dy float64) (float64, int) {
	if plane == nil {
		return math.Inf(-1), 0
	}
	baseX := tx * experimentalV4TileWidthBlocks
	baseY := ty * experimentalV4TileHeightBlocks
	// A local constant correction is a right-side/domain translation. It is
	// used only to estimate a control point; the final sampler uses the fitted
	// smooth field directly.
	domain := homography{h: [9]float64{1, 0, dx, 0, 1, dy, 0, 0, 1}}
	localH := experimentalV4BlindMultiplyHomography(h, domain)
	num, den := 0.0, 0.0
	visible := 0
	for i, pos := range candidate.positions {
		px := pos % experimentalV4TileWidthBlocks
		py := pos / experimentalV4TileWidthBlocks
		v, ok := readProjectiveBlockValue(plane, localH, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
		if !ok {
			continue
		}
		visible++
		num += float64(candidate.signs[i]) * v
		den += math.Abs(v)
	}
	if den <= 0 {
		return math.Inf(-1), visible
	}
	return num / den, visible
}

func experimentalV4PhoneResidualControls(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography) ([]experimentalV4PhoneResidualControl, int) {
	bw, bh := cw/blockSize, ch/blockSize
	tilesX := bw / experimentalV4TileWidthBlocks
	tilesY := bh / experimentalV4TileHeightBlocks
	controls := make([]experimentalV4PhoneResidualControl, 0, (tilesX*tilesY+1)/2)
	evals := 0
	coarse := [...]float64{-4, -2, 0, 2, 4}
	fine := [...]float64{-1, 0, 1}
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			if (tx+ty)&1 != 0 { // checkerboard B is held out
				continue
			}
			type local struct{ dx, dy, score float64 }
			bank := make([]local, 0, len(coarse)*len(coarse))
			for _, dy := range coarse {
				for _, dx := range coarse {
					score, visible := experimentalV4PhoneResidualTileScore(plane, candidate, cw, ch, h, tx, ty, dx, dy)
					evals++
					if visible == len(candidate.positions) && !math.IsInf(score, -1) {
						bank = append(bank, local{dx: dx, dy: dy, score: score})
					}
				}
			}
			if len(bank) == 0 {
				continue
			}
			sort.Slice(bank, func(i, j int) bool { return bank[i].score > bank[j].score })
			best := bank[0]
			// Refine only around the coarse winner; this is still bounded to the
			// +/-5 px residual basin and cannot become a hidden geometry search.
			for _, fy := range fine {
				for _, fx := range fine {
					if fx == 0 && fy == 0 {
						continue
					}
					dx, dy := best.dx+fx, best.dy+fy
					if math.Abs(dx) > experimentalV4PhoneResidualMaxPixels || math.Abs(dy) > experimentalV4PhoneResidualMaxPixels {
						continue
					}
					score, visible := experimentalV4PhoneResidualTileScore(plane, candidate, cw, ch, h, tx, ty, dx, dy)
					evals++
					if visible == len(candidate.positions) {
						bank = append(bank, local{dx: dx, dy: dy, score: score})
					}
				}
			}
			sort.Slice(bank, func(i, j int) bool { return bank[i].score > bank[j].score })
			best = bank[0]
			runner := math.Inf(-1)
			for i := 1; i < len(bank); i++ {
				// Nearby samples are the same peak. Runner-up must represent a
				// distinct local basin before the control is considered reliable.
				if math.Hypot(bank[i].dx-best.dx, bank[i].dy-best.dy) >= 2.0 {
					runner = bank[i].score
					break
				}
			}
			gap := best.score - runner
			if math.IsInf(runner, -1) {
				gap = best.score
			}
			if best.score < experimentalV4PhoneResidualMinLocalScore || gap < experimentalV4PhoneResidualMinLocalGap {
				continue
			}
			cx := float64((tx*experimentalV4TileWidthBlocks + experimentalV4TileWidthBlocks/2) * blockSize)
			cy := float64((ty*experimentalV4TileHeightBlocks + experimentalV4TileHeightBlocks/2) * blockSize)
			controls = append(controls, experimentalV4PhoneResidualControl{x: cx, y: cy, dx: best.dx, dy: best.dy, weight: math.Max(0.25, best.score*best.score), localScore: best.score, gap: gap})
		}
	}
	return controls, evals
}

func experimentalV4PhoneFitResidual(cw, ch int, controls []experimentalV4PhoneResidualControl) (*experimentalV4PhoneResidualWarp, bool) {
	if cw <= 1 || ch <= 1 || len(controls) < experimentalV4PhoneResidualMinControls {
		return nil, false
	}
	fit := func(in []experimentalV4PhoneResidualControl) (*experimentalV4PhoneResidualWarp, bool) {
		var normal [6][6]float64
		var rhsX, rhsY [6]float64
		for _, c := range in {
			nx := 2*c.x/float64(cw) - 1
			ny := 2*c.y/float64(ch) - 1
			basis := [6]float64{1, nx, ny, nx * ny, nx * nx, ny * ny}
			w := math.Max(0.20, c.weight)
			for r := 0; r < 6; r++ {
				rhsX[r] += w * basis[r] * c.dx
				rhsY[r] += w * basis[r] * c.dy
				for col := 0; col < 6; col++ {
					normal[r][col] += w * basis[r] * basis[col]
				}
			}
		}
		dx, okX := solveDiagnostic6(normal, rhsX)
		dy, okY := solveDiagnostic6(normal, rhsY)
		if !okX || !okY {
			return nil, false
		}
		warp := &experimentalV4PhoneResidualWarp{width: float64(cw), height: float64(ch), dx: dx, dy: dy, controls: len(in), maxPixels: experimentalV4PhoneResidualMaxPixels}
		err2, ws := 0.0, 0.0
		for _, c := range in {
			px, py := warp.correction(c.x, c.y)
			w := math.Max(0.20, c.weight)
			e := math.Hypot(px-c.dx, py-c.dy)
			err2 += w * e * e
			ws += w
		}
		if ws <= 0 {
			return nil, false
		}
		warp.rmsPixels = math.Sqrt(err2 / ws)
		return warp, true
	}

	working := append([]experimentalV4PhoneResidualControl(nil), controls...)
	for pass := 0; pass < 3; pass++ {
		warp, ok := fit(working)
		if !ok {
			return nil, false
		}
		if warp.rmsPixels <= experimentalV4PhoneResidualMaxFitRMS {
			return warp, true
		}
		filtered := working[:0]
		for _, c := range working {
			px, py := warp.correction(c.x, c.y)
			if math.Hypot(px-c.dx, py-c.dy) <= 3.0 {
				filtered = append(filtered, c)
			}
		}
		working = filtered
		if len(working) < experimentalV4PhoneResidualMinControls {
			return nil, false
		}
	}
	warp, ok := fit(working)
	return warp, ok && warp.rmsPixels <= experimentalV4PhoneResidualMaxFitRMS
}

func experimentalV4PhoneResidualFullDetection(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, warp *experimentalV4PhoneResidualWarp) ExperimentalV4PilotDetection {
	if plane == nil {
		return ExperimentalV4PilotDetection{}
	}
	bw, bh := cw/blockSize, ch/blockSize
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var sums [residues]float64
	var abss [residues]float64
	var counts [residues]int
	for by := 0; by < bh; by++ {
		for bx := 0; bx < bw; bx++ {
			v, ok := experimentalV4PhoneWarpedBlockValue(plane, h, warp, cw, ch, bx*blockSize, by*blockSize, blockSize)
			if !ok {
				continue
			}
			r := (by%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + bx%experimentalV4TileWidthBlocks
			sums[r] += v
			abss[r] += math.Abs(v)
			counts[r]++
		}
	}
	return experimentalV4DetectPilotFromResidues(candidate, sums[:], abss[:], counts[:])
}

func experimentalV4PhoneFitResidualPilotOnly(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography) (*experimentalV4PhoneResidualWarp, experimentalV4PhoneResidualInfo) {
	info := experimentalV4PhoneResidualInfo{}
	info.ProposalBefore, _ = experimentalV4PhoneResidualSpatialScore(plane, candidate, cw, ch, h, nil, 0)
	info.ValidationBefore, _ = experimentalV4PhoneResidualSpatialScore(plane, candidate, cw, ch, h, nil, 1)
	controls, evals := experimentalV4PhoneResidualControls(plane, candidate, cw, ch, h)
	info.Controls = len(controls)
	info.HypothesesEvaluated = evals
	warp, ok := experimentalV4PhoneFitResidual(cw, ch, controls)
	if !ok {
		return nil, info
	}
	info.RMSPixels = warp.rmsPixels
	info.ProposalAfter, _ = experimentalV4PhoneResidualSpatialScore(plane, candidate, cw, ch, h, warp, 0)
	info.ValidationAfter, _ = experimentalV4PhoneResidualSpatialScore(plane, candidate, cw, ch, h, warp, 1)
	full := experimentalV4PhoneResidualFullDetection(plane, candidate, cw, ch, h, warp)
	if full.Available {
		info.FullPilotScore = full.Score
		info.FullPilotMargin = full.Margin
		info.FullOriginXBlocks = full.OriginXBlocks
		info.FullOriginYBlocks = full.OriginYBlocks
	}
	info.Applied = info.ProposalAfter >= experimentalV4PhoneResidualProposalFloor &&
		info.ValidationAfter >= experimentalV4PhoneResidualValidationFloor &&
		info.ValidationAfter >= info.ValidationBefore+experimentalV4PhoneResidualValidationGain &&
		full.Available && full.OriginXBlocks == 0 && full.OriginYBlocks == 0 &&
		full.Score >= experimentalV4PhoneResidualFullScoreFloor && full.Margin >= experimentalV4PhoneResidualFullMarginFloor
	if !info.Applied {
		return nil, info
	}
	return warp, info
}

func experimentalV4ReadProtectedPhoneWarpMargins(plane *pixelPlane, cw, ch int, h homography, warp *experimentalV4PhoneResidualWarp, codedBits int) ([]float64, float64, bool) {
	if plane == nil || codedBits <= 0 {
		return nil, 0, false
	}
	bw, bh := cw/blockSize, ch/blockSize
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var sums [residues]float64
	var counts [residues]int
	for by := 0; by < bh; by++ {
		for bx := 0; bx < bw; bx++ {
			v, ok := experimentalV4PhoneWarpedBlockValue(plane, h, warp, cw, ch, bx*blockSize, by*blockSize, blockSize)
			if !ok {
				continue
			}
			r := (by%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + bx%experimentalV4TileWidthBlocks
			sums[r] += v
			counts[r]++
		}
	}
	bitSums := make([]float64, codedBits)
	bitCounts := make([]int, codedBits)
	for ordinal, pos := range experimentalV4LockedDataPositions() {
		if counts[pos] == 0 {
			continue
		}
		i := experimentalV4CodeIndex(ordinal, codedBits)
		bitSums[i] += sums[pos]
		bitCounts[i] += counts[pos]
	}
	margins := make([]float64, codedBits)
	confidence := 0.0
	seen := 0
	for i := range margins {
		if bitCounts[i] == 0 {
			continue
		}
		seen++
		margins[i] = bitSums[i] / float64(bitCounts[i])
		confidence += math.Abs(margins[i])
	}
	if seen != codedBits {
		return nil, 0, false
	}
	return margins, confidence / float64(codedBits), true
}
