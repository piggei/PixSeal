package watermark

import "math"

// Test-only candidate: fixed8 contract in continuation4 only. Pinned qualified
// sources and production dispatch are unchanged. No arithmetic is reassociated.
func experimentalV4Build104ReadProjectiveBlockValue(src *pixelPlane, h homography, originX, originY int) (float64, bool) {
	// Fixed8 contract; preserve all qualified sampling/arithmetic statements.
	width, height := src.bounds.Dx(), src.bounds.Dy()
	rgb := src.rgb
	rowStride := width * 3
	table := readCosTables[blockSize]
	table2, table3 := (*[blockSize]float64)(table[2]), (*[blockSize]float64)(table[3])
	c23, c32 := 0.0, 0.0
	for y := 0; y < blockSize; y++ {
		for x := 0; x < blockSize; x++ {
			sx, sy, ok := h.mapPoint(float64(originX+x), float64(originY+y))
			if !ok {
				return 0, false
			}
			if sx < 0 || sy < 0 || sx > float64(width-1) || sy > float64(height-1) {
				return 0, false
			}
			x0, y0 := int(math.Floor(sx)), int(math.Floor(sy))
			x1, y1 := x0+1, y0+1
			if x1 >= width {
				x1 = width - 1
			}
			if y1 >= height {
				y1 = height - 1
			}
			fx, fy := sx-float64(x0), sy-float64(y0)

			row0 := y0 * rowStride
			row1 := y1 * rowStride
			x03 := x0 * 3
			x13 := x1 * 3
			index00 := row0 + x03
			index10 := row0 + x13
			index01 := row1 + x03
			index11 := row1 + x13

			pixel00 := rgb[index00 : index00+3]
			pixel10 := rgb[index10 : index10+3]
			pixel01 := rgb[index01 : index01+3]
			pixel11 := rgb[index11 : index11+3]
			l00 := .299*float64(pixel00[0]) + .587*float64(pixel00[1]) + .114*float64(pixel00[2]) - 128
			l10 := .299*float64(pixel10[0]) + .587*float64(pixel10[1]) + .114*float64(pixel10[2]) - 128
			l01 := .299*float64(pixel01[0]) + .587*float64(pixel01[1]) + .114*float64(pixel01[2]) - 128
			l11 := .299*float64(pixel11[0]) + .587*float64(pixel11[1]) + .114*float64(pixel11[2]) - 128
			top := l00*(1-fx) + l10*fx
			bottom := l01*(1-fx) + l11*fx
			l := top*(1-fy) + bottom*fy

			c23 += l * table3[x] * table2[y]
			c32 += l * table2[x] * table3[y]
		}
	}
	return math.Abs(c23) - math.Abs(c32), true
}

func experimentalV4Build104FoldScoreFetch(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, heldout int, proposal bool, sparse bool) (float64, int, experimentalV4PhoneBuild84FetchTelemetry) {
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	fetch.FoldScores = 1
	if plane == nil {
		return math.Inf(-1), 0, fetch
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return math.Inf(-1), 0, fetch
	}
	num, den := 0.0, 0.0
	visible, usedTiles := 0, 0
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			isHeldout := experimentalV4PhoneBuild41FoldForTile(tx, ty) == heldout
			if proposal == isHeldout {
				continue
			}
			if sparse {
				if (tx+3*ty+heldout)%4 != 0 {
					continue
				}
				if usedTiles >= experimentalV4PhoneBuild41SparseTiles {
					continue
				}
			}
			usedTiles++
			baseX, baseY := tx*experimentalV4TileWidthBlocks, ty*experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				if sparse && i%4 != 0 {
					continue
				}
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				fetch.BlockReads++
				v, ok := experimentalV4Build104ReadProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize)
				if !ok {
					fetch.BlockFailed++
					continue
				}
				fetch.BlockSuccess++
				visible++
				num += float64(candidate.signs[i]) * v
				den += math.Abs(v)
			}
		}
	}
	if den <= 0 {
		return math.Inf(-1), visible, fetch
	}
	return num / den, visible, fetch
}

func experimentalV4Build104ContinueFetch(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild55BlindContinuationState, int, experimentalV4PhoneBuild84FetchTelemetry) {
	q, h, score := start.quad, start.h, start.proposal
	states := make([]experimentalV4PhoneBuild55BlindContinuationState, 0, experimentalV4PhoneBuild55MaxStatesBranch)
	evals := 0
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	for pass := 0; pass < experimentalV4PhoneBuild55MaxPasses; pass++ {
		improved := false
		for dim := 0; dim < experimentalV4PhoneBuild55Dimensions; dim++ {
			bestQ, bestH, best := q, h, score
			bestDelta := 0.0
			for _, sign := range []float64{-1, 1} {
				delta := sign * experimentalV4PhoneBuild55Step
				qq := experimentalV4PhoneBuild51ApplyDim(q, dim, delta)
				if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
					continue
				}
				hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
				if !ok {
					continue
				}
				ss, _, ip := experimentalV4Build104FoldScoreFetch(plane, candidate, cw, ch, hh, heldout, true, false)
				experimentalV4PhoneBuild84MergeFetch(&fetch, ip)
				evals++
				if !math.IsInf(ss, 0) && !math.IsNaN(ss) && ss > best+1e-7 {
					bestQ, bestH, best, bestDelta = qq, hh, ss, delta
				}
			}
			if best > score+1e-7 {
				q, h, score = bestQ, bestH, best
				improved = true
				states = append(states, experimentalV4PhoneBuild55BlindContinuationState{index: len(states) + 1, pass: pass + 1, dim: dim, delta: bestDelta, hyp: experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}})
			}
		}
		if !improved {
			break
		}
	}
	return states, evals, fetch
}

func experimentalV4Build104Generation4Fetch(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild84FetchTelemetry) {
	bank := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	se4, si4 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, input, 0)
	evals += se4
	if si4 != 0 {
		return bank, evals, fetch
	}
	pairs4, pe4, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, input, 0)
	evals += pe4
	for _, pair4 := range pairs4 {
		cont4, ce4, ip := experimentalV4Build104ContinueFetch(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
		experimentalV4PhoneBuild84MergeFetch(&fetch, ip)
		evals += ce4
		for _, c4 := range cont4 {
			bank = append(bank, c4.hyp)
		}
	}
	return bank, evals, fetch
}

// experimentalV4PhoneBuild84SeedPrefix1 reproduces Build64 exactly through the
// first sibling stencil and freezes sibling1 hypotheses in traversal order.
