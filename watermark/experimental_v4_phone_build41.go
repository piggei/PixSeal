package watermark

import (
	"image"
	"math"
	"sort"
)

// Build41 keeps the Build39/40 phone registration checkpoint intact and adds a
// new bounded basin finder for real smartphone captures. Geometry remains
// public-pilot/structure only: neither payload, ECC, key nor HMAC participates
// in boundary generation, refinement, ranking or qualification.
const (
	experimentalV4PhoneBuild41FoldCount       = 3
	experimentalV4PhoneBuild41CornerLimitPx   = 56.0
	experimentalV4PhoneBuild41StartOffsetPx   = 12.0
	experimentalV4PhoneBuild41SparsePilots    = 16
	experimentalV4PhoneBuild41SparseTiles     = 8
	experimentalV4PhoneBuild41MinPairDistance = 0.75
)

type experimentalV4PhoneBuild41Candidate struct {
	hyp      experimentalV4PhoneHypothesis
	fold     int
	boundary int
}

func experimentalV4PhoneBuild41Quad(boundary PrintBoundaryEstimate) [4]ImagePoint {
	return [4]ImagePoint{boundary.TopLeft, boundary.TopRight, boundary.BottomLeft, boundary.BottomRight}
}

func experimentalV4PhoneBuild41QuadDistance(a, b [4]ImagePoint) float64 {
	sum := 0.0
	for i := range a {
		dx := a[i].X - b[i].X
		dy := a[i].Y - b[i].Y
		sum += dx*dx + dy*dy
	}
	return math.Sqrt(sum / 4)
}

func experimentalV4PhoneBuild41BoundarySaneForImage(src image.Image, q PrintBoundaryEstimate) bool {
	if src == nil || !q.Detected {
		return false
	}
	b := src.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())
	pts := [4]ImagePoint{q.TopLeft, q.TopRight, q.BottomLeft, q.BottomRight}
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, p := range pts {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
			return false
		}
		if p.X < -0.03*w || p.X > 1.03*w || p.Y < -0.03*h || p.Y > 1.03*h {
			return false
		}
		minX = math.Min(minX, p.X)
		maxX = math.Max(maxX, p.X)
		minY = math.Min(minY, p.Y)
		maxY = math.Max(maxY, p.Y)
	}
	return maxX-minX > 0.30*w && maxY-minY > 0.30*h
}

// The strict component detector is deliberately separate from Build40's
// boundary chooser. Build40 may accept a locally strong paper edge early; for
// Build41 we also retain the largest interior non-paper component as an
// independent structural seed, which is particularly useful under keystone.
func experimentalV4PhoneBuild41ComponentBoundary(src image.Image) PrintBoundaryEstimate {
	plane := experimentalV4BuildScannerColorPlane(src, experimentalV4PhoneBoundaryDimension)
	if plane == nil || plane.width < 300 || plane.height < 300 {
		return PrintBoundaryEstimate{}
	}
	paper, ok := estimateBoundaryPaperColor(plane)
	if !ok {
		return PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	}
	n := plane.width * plane.height
	raw := make([]bool, n)
	marginX, marginY := plane.width/50, plane.height/50
	for y := marginY; y < plane.height-marginY; y++ {
		for x := marginX; x < plane.width-marginX; x++ {
			if boundaryColorDistance(plane, paper, float64(x), float64(y)) > experimentalV4PhoneForegroundThreshold {
				raw[y*plane.width+x] = true
			}
		}
	}
	mask := append([]bool(nil), raw...)
	for pass := 0; pass < 2; pass++ {
		next := append([]bool(nil), mask...)
		for y := 1; y < plane.height-1; y++ {
			for x := 1; x < plane.width-1; x++ {
				if mask[y*plane.width+x] {
					continue
				}
				hit := false
				for dy := -1; dy <= 1 && !hit; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if mask[(y+dy)*plane.width+x+dx] {
							hit = true
							break
						}
					}
				}
				if hit {
					next[y*plane.width+x] = true
				}
			}
		}
		mask = next
	}

	seen := make([]bool, n)
	best := make([]int, 0)
	queue := make([]int, 0, 4096)
	for y := marginY; y < plane.height-marginY; y++ {
		for x := marginX; x < plane.width-marginX; x++ {
			idx := y*plane.width + x
			if !mask[idx] || seen[idx] {
				continue
			}
			queue = queue[:0]
			queue = append(queue, idx)
			seen[idx] = true
			comp := make([]int, 0, 1024)
			minx, maxx, miny, maxy := x, x, y, y
			for head := 0; head < len(queue); head++ {
				p := queue[head]
				comp = append(comp, p)
				px, py := p%plane.width, p/plane.width
				minx = minInt(minx, px)
				maxx = maxInt(maxx, px)
				miny = minInt(miny, py)
				maxy = maxInt(maxy, py)
				if px > 0 {
					q := p - 1
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
				if px+1 < plane.width {
					q := p + 1
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
				if py > 0 {
					q := p - plane.width
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
				if py+1 < plane.height {
					q := p + plane.width
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
			}
			if maxx-minx < plane.width/6 || maxy-miny < plane.height/6 {
				continue
			}
			// The table/background surrounding the sheet is also non-paper,
			// but it is connected to the analysis frame. Build41 needs the
			// largest interior non-paper component: the printed artwork.
			if minx <= marginX+2 || miny <= marginY+2 || maxx >= plane.width-marginX-3 || maxy >= plane.height-marginY-3 {
				continue
			}
			if len(comp) > len(best) {
				best = append(best[:0], comp...)
			}
		}
	}

	result := PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	if len(best) < n/80 {
		return result
	}
	inComp := make([]bool, n)
	for _, p := range best {
		inComp[p] = true
	}
	topPts, bottomPts := make([]ImagePoint, 0), make([]ImagePoint, 0)
	leftPts, rightPts := make([]ImagePoint, 0), make([]ImagePoint, 0)
	for x := marginX; x < plane.width-marginX; x++ {
		top, bottom := -1, -1
		for y := marginY; y < plane.height-marginY; y++ {
			if inComp[y*plane.width+x] {
				if top < 0 {
					top = y
				}
				bottom = y
			}
		}
		if top >= 0 {
			topPts = append(topPts, ImagePoint{X: float64(x), Y: float64(top)})
			bottomPts = append(bottomPts, ImagePoint{X: float64(x), Y: float64(bottom)})
		}
	}
	for y := marginY; y < plane.height-marginY; y++ {
		left, right := -1, -1
		for x := marginX; x < plane.width-marginX; x++ {
			if inComp[y*plane.width+x] {
				if left < 0 {
					left = x
				}
				right = x
			}
		}
		if left >= 0 {
			leftPts = append(leftPts, ImagePoint{X: float64(left), Y: float64(y)})
			rightPts = append(rightPts, ImagePoint{X: float64(right), Y: float64(y)})
		}
	}
	topLine, ts, okT := experimentalV4ScannerRobustLine(topPts, true)
	bottomLine, bs, okB := experimentalV4ScannerRobustLine(bottomPts, true)
	leftLine, ls, okL := experimentalV4ScannerRobustLine(leftPts, false)
	rightLine, rs, okR := experimentalV4ScannerRobustLine(rightPts, false)
	result.TopScore, result.BottomScore, result.LeftScore, result.RightScore = ts, bs, ls, rs
	if !(okT && okB && okL && okR) {
		return result
	}
	tl, o1 := experimentalV4ScannerIntersect(topLine, leftLine)
	tr, o2 := experimentalV4ScannerIntersect(topLine, rightLine)
	bl, o3 := experimentalV4ScannerIntersect(bottomLine, leftLine)
	br, o4 := experimentalV4ScannerIntersect(bottomLine, rightLine)
	if !(o1 && o2 && o3 && o4) || !boundaryQuadSane(plane, tl, tr, br, bl) {
		return result
	}
	d := float64(plane.divisor)
	result.TopLeft = ImagePoint{X: tl.X * d, Y: tl.Y * d}
	result.TopRight = ImagePoint{X: tr.X * d, Y: tr.Y * d}
	result.BottomLeft = ImagePoint{X: bl.X * d, Y: bl.Y * d}
	result.BottomRight = ImagePoint{X: br.X * d, Y: br.Y * d}
	fit := math.Min(math.Min(ts, bs), math.Min(ls, rs))
	result.Confidence = clampUnit(fit)
	result.Detected = fit >= 0.30
	return result
}

func experimentalV4PhoneBuild41ArtworkBoundary(src image.Image) PrintBoundaryEstimate {
	plane := experimentalV4BuildScannerColorPlane(src, experimentalV4PhoneBoundaryDimension)
	if plane == nil || plane.width < 300 || plane.height < 300 {
		return PrintBoundaryEstimate{}
	}
	paper, ok := estimateBoundaryPaperColor(plane)
	if !ok {
		return PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	}
	const foreground = 28.0
	const paperLimit = 20.0
	gapX := maxInt(5, plane.width/180)
	gapY := maxInt(5, plane.height/180)
	paperWinX := maxInt(12, plane.width/90)
	paperWinY := maxInt(12, plane.height/90)
	marginX := maxInt(8, plane.width/40)
	marginY := maxInt(8, plane.height/40)

	meanDistRow := func(y, x0, x1 int) float64 {
		if x0 < 0 {
			x0 = 0
		}
		if x1 > plane.width {
			x1 = plane.width
		}
		if x1 <= x0 {
			return 100
		}
		s := 0.0
		for x := x0; x < x1; x++ {
			s += boundaryColorDistance(plane, paper, float64(x), float64(y))
		}
		return s / float64(x1-x0)
	}
	meanDistCol := func(x, y0, y1 int) float64 {
		if y0 < 0 {
			y0 = 0
		}
		if y1 > plane.height {
			y1 = plane.height
		}
		if y1 <= y0 {
			return 100
		}
		s := 0.0
		for y := y0; y < y1; y++ {
			s += boundaryColorDistance(plane, paper, float64(x), float64(y))
		}
		return s / float64(y1-y0)
	}

	type rowRun struct{ y, left, right, width int }
	rows := make([]rowRun, 0, plane.height)
	maxRowWidth := 0
	for y := marginY; y < plane.height-marginY; y++ {
		bestL, bestR, bestW := -1, -1, 0
		for x := marginX; x < plane.width-marginX; {
			for x < plane.width-marginX && boundaryColorDistance(plane, paper, float64(x), float64(y)) <= foreground {
				x++
			}
			if x >= plane.width-marginX {
				break
			}
			start := x
			lastFG := x
			gap := 0
			for x < plane.width-marginX {
				if boundaryColorDistance(plane, paper, float64(x), float64(y)) > foreground {
					lastFG = x
					gap = 0
				} else {
					gap++
					if gap > gapX {
						break
					}
				}
				x++
			}
			end := lastFG
			width := end - start + 1
			if width > bestW && width > plane.width/7 &&
				meanDistRow(y, start-paperWinX, start-2) < paperLimit &&
				meanDistRow(y, end+3, end+paperWinX) < paperLimit {
				bestL, bestR, bestW = start, end, width
			}
		}
		if bestW > 0 {
			rows = append(rows, rowRun{y: y, left: bestL, right: bestR, width: bestW})
			if bestW > maxRowWidth {
				maxRowWidth = bestW
			}
		}
	}
	leftPts, rightPts := make([]ImagePoint, 0, len(rows)), make([]ImagePoint, 0, len(rows))
	for _, r := range rows {
		if maxRowWidth > 0 && float64(r.width) >= 0.75*float64(maxRowWidth) {
			leftPts = append(leftPts, ImagePoint{X: float64(r.left), Y: float64(r.y)})
			rightPts = append(rightPts, ImagePoint{X: float64(r.right), Y: float64(r.y)})
		}
	}

	type colRun struct{ x, top, bottom, height int }
	cols := make([]colRun, 0, plane.width)
	maxColHeight := 0
	for x := marginX; x < plane.width-marginX; x++ {
		bestT, bestB, bestH := -1, -1, 0
		for y := marginY; y < plane.height-marginY; {
			for y < plane.height-marginY && boundaryColorDistance(plane, paper, float64(x), float64(y)) <= foreground {
				y++
			}
			if y >= plane.height-marginY {
				break
			}
			start := y
			lastFG := y
			gap := 0
			for y < plane.height-marginY {
				if boundaryColorDistance(plane, paper, float64(x), float64(y)) > foreground {
					lastFG = y
					gap = 0
				} else {
					gap++
					if gap > gapY {
						break
					}
				}
				y++
			}
			end := lastFG
			height := end - start + 1
			if height > bestH && height > plane.height/7 &&
				meanDistCol(x, start-paperWinY, start-2) < paperLimit &&
				meanDistCol(x, end+3, end+paperWinY) < paperLimit {
				bestT, bestB, bestH = start, end, height
			}
		}
		if bestH > 0 {
			cols = append(cols, colRun{x: x, top: bestT, bottom: bestB, height: bestH})
			if bestH > maxColHeight {
				maxColHeight = bestH
			}
		}
	}
	topPts, bottomPts := make([]ImagePoint, 0, len(cols)), make([]ImagePoint, 0, len(cols))
	for _, c := range cols {
		if maxColHeight > 0 && float64(c.height) >= 0.75*float64(maxColHeight) {
			topPts = append(topPts, ImagePoint{X: float64(c.x), Y: float64(c.top)})
			bottomPts = append(bottomPts, ImagePoint{X: float64(c.x), Y: float64(c.bottom)})
		}
	}

	result := PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	topLine, ts, okT := experimentalV4ScannerRobustLine(topPts, true)
	bottomLine, bs, okB := experimentalV4ScannerRobustLine(bottomPts, true)
	leftLine, ls, okL := experimentalV4ScannerRobustLine(leftPts, false)
	rightLine, rs, okR := experimentalV4ScannerRobustLine(rightPts, false)
	result.TopScore, result.BottomScore, result.LeftScore, result.RightScore = ts, bs, ls, rs
	if !(okT && okB && okL && okR) {
		return result
	}
	tl, o1 := experimentalV4ScannerIntersect(topLine, leftLine)
	tr, o2 := experimentalV4ScannerIntersect(topLine, rightLine)
	bl, o3 := experimentalV4ScannerIntersect(bottomLine, leftLine)
	br, o4 := experimentalV4ScannerIntersect(bottomLine, rightLine)
	if !(o1 && o2 && o3 && o4) || !boundaryQuadSane(plane, tl, tr, br, bl) {
		return result
	}
	d := float64(plane.divisor)
	result.TopLeft = ImagePoint{X: tl.X * d, Y: tl.Y * d}
	result.TopRight = ImagePoint{X: tr.X * d, Y: tr.Y * d}
	result.BottomLeft = ImagePoint{X: bl.X * d, Y: bl.Y * d}
	result.BottomRight = ImagePoint{X: br.X * d, Y: br.Y * d}
	fit := math.Min(math.Min(ts, bs), math.Min(ls, rs))
	result.Confidence = clampUnit(fit)
	result.Detected = fit >= 0.20
	return result
}

func experimentalV4PhoneBuild41BoundaryCandidates(src image.Image) []PrintBoundaryEstimate {
	candidates := make([]PrintBoundaryEstimate, 0, 4)
	add := func(q PrintBoundaryEstimate) {
		if !experimentalV4PhoneBuild41BoundarySaneForImage(src, q) {
			return
		}
		qq := experimentalV4PhoneBuild41Quad(q)
		for _, prior := range candidates {
			if experimentalV4PhoneBuild41QuadDistance(qq, experimentalV4PhoneBuild41Quad(prior)) < 3.0 {
				return
			}
		}
		candidates = append(candidates, q)
	}
	add(experimentalV4PhoneBoundary(src))
	add(experimentalV4ScannerBoundary(src))
	add(experimentalV4PhoneBuild41ComponentBoundary(src))
	add(experimentalV4PhoneBuild41ArtworkBoundary(src))
	return candidates
}

func experimentalV4PhoneBuild41FoldForTile(tx, ty int) int {
	return positiveMod(tx+2*ty, experimentalV4PhoneBuild41FoldCount)
}

// score fixed physical origin (0,0). The visible artwork boundary already
// anchors absolute phase; Build41 never lets a proposal-side cyclic alias move
// the quadrilateral before independent validation.
func experimentalV4PhoneBuild41FoldScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, heldout int, proposal bool, sparse bool) (float64, int) {
	if plane == nil {
		return math.Inf(-1), 0
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return math.Inf(-1), 0
	}
	num, den := 0.0, 0.0
	visible := 0
	usedTiles := 0
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			isHeldout := experimentalV4PhoneBuild41FoldForTile(tx, ty) == heldout
			if proposal == isHeldout {
				continue
			}
			if sparse {
				// deterministic spatial thinning; proposal-only evidence.
				if (tx+3*ty+heldout)%4 != 0 {
					continue
				}
				if usedTiles >= experimentalV4PhoneBuild41SparseTiles {
					continue
				}
			}
			usedTiles++
			baseX := tx * experimentalV4TileWidthBlocks
			baseY := ty * experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				if sparse && i%4 != 0 {
					continue
				}
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				v, ok := readProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
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

func experimentalV4PhoneBuild41Starts(seed [4]ImagePoint) [][4]ImagePoint {
	offsets := [...]float64{-experimentalV4PhoneBuild41StartOffsetPx, 0, experimentalV4PhoneBuild41StartOffsetPx}
	out := make([][4]ImagePoint, 0, 9)
	for _, dy := range offsets {
		for _, dx := range offsets {
			q := seed
			for i := range q {
				q[i].X += dx
				q[i].Y += dy
			}
			out = append(out, q)
		}
	}
	return out
}

func experimentalV4PhoneBuild41WithinLimit(q, anchor [4]ImagePoint) bool {
	for i := range q {
		if math.Abs(q[i].X-anchor[i].X) > experimentalV4PhoneBuild41CornerLimitPx || math.Abs(q[i].Y-anchor[i].Y) > experimentalV4PhoneBuild41CornerLimitPx {
			return false
		}
	}
	return true
}

func experimentalV4PhoneBuild41Refine(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor, start [4]ImagePoint, heldout int) (experimentalV4PhoneHypothesis, int) {
	q := start
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return experimentalV4PhoneHypothesis{}, 0
	}
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, heldout, true, false)
	evals := 1
	if math.IsInf(score, -1) {
		return experimentalV4PhoneHypothesis{}, evals
	}
	steps := [...]float64{16, 8, 4, 2, 1}
	apply := func(in [4]ImagePoint, dim int, delta float64) [4]ImagePoint {
		out := in
		corner := dim / 2
		if dim%2 == 0 {
			out[corner].X += delta
		} else {
			out[corner].Y += delta
		}
		return out
	}
	for _, step := range steps {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 8; dim++ {
				bestQ, bestH, best := q, h, score
				for _, sgn := range []float64{-1, 1} {
					qq := apply(q, dim, sgn*step)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					ss, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
					evals++
					if ss > best+1e-7 {
						bestQ, bestH, best = qq, hh, ss
					}
				}
				if best > score+1e-7 {
					q, h, score = bestQ, bestH, best
					improved = true
				}
			}
			if !improved {
				break
			}
		}
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals
}

func experimentalV4PhoneBuild41BestStart(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, heldout int) ([4]ImagePoint, int, bool) {
	starts := experimentalV4PhoneBuild41Starts(anchor)
	best := math.Inf(-1)
	bestQ := anchor
	evals := 0
	okAny := false
	for _, st := range starts {
		h, ok := experimentalV4PhoneQuadHomography(cw, ch, st)
		if !ok {
			continue
		}
		score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, heldout, true, true)
		evals++
		if !math.IsInf(score, -1) && (!okAny || score > best) {
			best, bestQ, okAny = score, st, true
		}
	}
	return bestQ, evals, okAny
}

func experimentalV4PhoneBuild41Qualify(plane *pixelPlane, img image.Image, candidate experimentalV4PilotCandidate, cw, ch int, hyp experimentalV4PhoneHypothesis, heldout int) (experimentalV4PhoneHypothesis, int, bool) {
	validation, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hyp.h, heldout, false, false)
	evals := 1
	hyp.validation = validation
	if hyp.proposal < experimentalV4PhoneProposalFloor || validation < experimentalV4PhoneValidationFloor {
		return hyp, evals, false
	}
	det := experimentalV4DetectPilotProjective(img, candidate, cw, ch, hyp.h)
	evals++
	hyp.detection = det
	ok := det.Available && det.OriginXBlocks == 0 && det.OriginYBlocks == 0 && det.Score >= experimentalV4PhonePilotScoreFloor && det.Margin >= experimentalV4PhonePilotMarginFloor
	return hyp, evals, ok
}

func experimentalV4PhoneBuild41PairScore(a, b experimentalV4PhoneBuild41Candidate) float64 {
	// Held-out is primary; full-pilot is already an independent final gate.
	return math.Min(a.hyp.validation, b.hyp.validation) + 0.5*math.Min(a.hyp.detection.Score, b.hyp.detection.Score) + 0.25*math.Min(a.hyp.proposal, b.hyp.proposal)
}

func experimentalV4PhoneSearchBuild41Detailed(src image.Image, cw, ch int) (image.Image, PrintBoundaryEstimate, []experimentalV4PhoneHypothesis, []experimentalV4PhoneHypothesis, int, bool) {
	work, down := experimentalV4PhoneResize(src, experimentalV4PhoneMaxDimension)
	candidate := experimentalV4Prototype2Candidate()
	plane := newPixelPlane(work)

	// Structural boundary selection is independent of the watermark. Prefer the
	// dedicated artwork edge when available; otherwise retain the established
	// Build39 phone boundary. Component/scanner boundaries are structural
	// fallbacks only when the preferred detector is unavailable.
	boundary := experimentalV4PhoneBuild41ArtworkBoundary(work)
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		boundary = experimentalV4PhoneBoundary(work)
	}
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		boundary = experimentalV4PhoneBuild41ComponentBoundary(work)
	}
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		boundary = experimentalV4ScannerBoundary(work)
	}
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return work, boundary, nil, nil, 0, down
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	const heldout = 0 // fixed, deterministic 1/3 spatial holdout
	evals := 0

	// Every geometry operation below uses only the 2/3 proposal tiles. The
	// held-out fold is not sampled until the proposal-only bank is frozen.
	coarse, n := experimentalV4PhoneBuild41RefineShapeFold(plane, candidate, cw, ch, anchor, heldout)
	evals += n
	if coarse.h.h[8] == 0 {
		return work, boundary, nil, nil, evals, down
	}
	shape, n := experimentalV4PhoneBuild41PolishShapeFold(plane, candidate, cw, ch, anchor, coarse, heldout)
	evals += n
	if shape.h.h[8] == 0 {
		return work, boundary, nil, nil, evals, down
	}
	bank, n := experimentalV4PhoneBuild41PhaseBank(plane, candidate, cw, ch, anchor, shape, heldout)
	evals += n
	if len(bank) == 0 {
		return work, boundary, nil, nil, evals, down
	}

	// Expand the proposal-only shortlist before held-out evidence is read. The
	// top three phase candidates receive a small fixed-origin corner polish,
	// still scored exclusively on proposal tiles. Both the original and polished
	// hypotheses are frozen into the bank; validation may only accept/reject them.
	frozen := make([]experimentalV4PhoneHypothesis, 0, len(bank)+3)
	for i, hyp := range bank {
		frozen = append(frozen, hyp)
		if i < 3 {
			fine, m := experimentalV4PhoneBuild41FineFixedOrigin(plane, candidate, cw, ch, anchor, hyp, heldout)
			evals += m
			frozen = append(frozen, fine)
		}
	}

	// The shortlist is now frozen. Held-out evidence may qualify/reject and rank
	// candidates, but it never creates or refines geometry.
	qualified := make([]experimentalV4PhoneBuild41Candidate, 0, len(frozen))
	for i, hyp := range frozen {
		q, m, ok := experimentalV4PhoneBuild41Qualify(plane, work, candidate, cw, ch, hyp, heldout)
		evals += m
		if ok {
			qualified = append(qualified, experimentalV4PhoneBuild41Candidate{hyp: q, fold: heldout, boundary: i})
		}
	}
	qualifiedBank := make([]experimentalV4PhoneHypothesis, len(qualified))
	for i := range qualified {
		qualifiedBank[i] = qualified[i].hyp
	}
	if len(qualified) < experimentalV4PhoneEnsembleSize {
		return work, boundary, nil, qualifiedBank, evals, down
	}

	bestPairScore := math.Inf(-1)
	bestPair := make([]experimentalV4PhoneHypothesis, 0, experimentalV4PhoneEnsembleSize)
	// Two nearly identical geometries do not form a useful ensemble. Require a
	// small, scale-normalized separation so aggregation brackets sub-pixel
	// registration uncertainty instead of duplicating one local optimum.
	minPairDistance := math.Max(experimentalV4PhoneBuild41MinPairDistance, 0.20*experimentalV4PhoneAverageBlockScale(anchor, cw, ch))
	for i := 0; i < len(qualified); i++ {
		for j := i + 1; j < len(qualified); j++ {
			if experimentalV4PhoneBuild41QuadDistance(qualified[i].hyp.quad, qualified[j].hyp.quad) < minPairDistance {
				continue
			}
			score := experimentalV4PhoneBuild41PairScore(qualified[i], qualified[j])
			if score > bestPairScore {
				bestPairScore = score
				bestPair = []experimentalV4PhoneHypothesis{qualified[i].hyp, qualified[j].hyp}
			}
		}
	}
	if len(bestPair) < experimentalV4PhoneEnsembleSize {
		return work, boundary, nil, nil, evals, down
	}
	sort.SliceStable(bestPair, func(i, j int) bool {
		if bestPair[i].validation == bestPair[j].validation {
			return bestPair[i].detection.Score > bestPair[j].detection.Score
		}
		return bestPair[i].validation > bestPair[j].validation
	})
	return work, boundary, bestPair, qualifiedBank, evals, down
}

func experimentalV4PhoneSearchBuild41(src image.Image, cw, ch int) (image.Image, PrintBoundaryEstimate, []experimentalV4PhoneHypothesis, int, bool) {
	work, boundary, pair, _, evals, down := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
	return work, boundary, pair, evals, down
}

// Build41 shape refinement deliberately ignores absolute cyclic phase while it
// adjusts the four artwork corners. This makes the objective stable under a
// few-pixel boundary error. Absolute phase is restored later by a small,
// physically bounded translation scored at fixed origin (0,0).
func experimentalV4PhoneBuild41RefineShape(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, partition int) (experimentalV4PhoneHypothesis, int) {
	q := anchor
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return experimentalV4PhoneHypothesis{}, 0
	}
	det := experimentalV4PhoneBuild41CyclicSparseDetection(plane, candidate, cw, ch, partition, h)
	evals := 1
	if !det.Available {
		return experimentalV4PhoneHypothesis{}, evals
	}
	score := det.Score
	apply := func(in [4]ImagePoint, dim int, delta float64) [4]ImagePoint {
		out := in
		corner := dim / 2
		if dim%2 == 0 {
			out[corner].X += delta
		} else {
			out[corner].Y += delta
		}
		return out
	}
	steps := [...]float64{24, 12, 6, 3, 1.5}
	for _, step := range steps {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 8; dim++ {
				bestQ, bestH, bestDet, best := q, h, det, score
				for _, sign := range []float64{-1, 1} {
					qq := apply(q, dim, sign*step)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					dd := experimentalV4PhoneBuild41CyclicSparseDetection(plane, candidate, cw, ch, partition, hh)
					evals++
					if dd.Available && dd.Score > best+1e-7 {
						bestQ, bestH, bestDet, best = qq, hh, dd, dd.Score
					}
				}
				if best > score+1e-7 {
					q, h, det, score = bestQ, bestH, bestDet, best
					improved = true
				}
			}
			if !improved {
				break
			}
		}
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals
}

func experimentalV4PhoneBuild41Translate(q [4]ImagePoint, dx, dy float64) [4]ImagePoint {
	out := q
	for i := range out {
		out[i].X += dx
		out[i].Y += dy
	}
	return out
}

func experimentalV4PhoneBuild41PhaseAlign(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, shape experimentalV4PhoneHypothesis, heldout int) (experimentalV4PhoneHypothesis, int) {
	best := math.Inf(-1)
	bestQ := shape.quad
	bestH := shape.h
	evals := 0
	evalAt := func(dx, dy float64) {
		qq := experimentalV4PhoneBuild41Translate(shape.quad, dx, dy)
		if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
			return
		}
		hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
		if !ok {
			return
		}
		score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
		evals++
		if score > best+1e-7 {
			best, bestQ, bestH = score, qq, hh
		}
	}
	bestDX, bestDY := 0.0, 0.0
	for dy := -16.0; dy <= 16; dy += 4 {
		for dx := -16.0; dx <= 16; dx += 4 {
			before := best
			evalAt(dx, dy)
			if best > before+1e-7 {
				bestDX, bestDY = dx, dy
			}
		}
	}
	coarseDX, coarseDY := bestDX, bestDY
	for dy := coarseDY - 3; dy <= coarseDY+3; dy += 1 {
		for dx := coarseDX - 3; dx <= coarseDX+3; dx += 1 {
			before := best
			evalAt(dx, dy)
			if best > before+1e-7 {
				bestDX, bestDY = dx, dy
			}
		}
	}
	if math.IsInf(best, -1) {
		return experimentalV4PhoneHypothesis{}, evals
	}
	return experimentalV4PhoneHypothesis{quad: bestQ, h: bestH, proposal: best}, evals
}

func experimentalV4PhoneBuild41FineFixedOrigin(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) (experimentalV4PhoneHypothesis, int) {
	q, h, score := start.quad, start.h, start.proposal
	evals := 0
	apply := func(in [4]ImagePoint, dim int, delta float64) [4]ImagePoint {
		out := in
		corner := dim / 2
		if dim%2 == 0 {
			out[corner].X += delta
		} else {
			out[corner].Y += delta
		}
		return out
	}
	for _, step := range []float64{2, 1, 0.5} {
		for dim := 0; dim < 8; dim++ {
			bestQ, bestH, best := q, h, score
			for _, sign := range []float64{-1, 1} {
				qq := apply(q, dim, sign*step)
				if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
					continue
				}
				hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
				if !ok {
					continue
				}
				ss, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
				evals++
				if ss > best+1e-7 {
					bestQ, bestH, best = qq, hh, ss
				}
			}
			if best > score+1e-7 {
				q, h, score = bestQ, bestH, best
			}
		}
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals
}

func experimentalV4PhoneBuild41CyclicSparseDetection(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, partition int, h homography) ExperimentalV4PilotDetection {
	if plane == nil {
		return ExperimentalV4PilotDetection{}
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return ExperimentalV4PilotDetection{}
	}
	type tile struct{ x, y int }
	all := make([]tile, 0, tilesX*tilesY/2+1)
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			if (tx+ty)&1 == partition {
				all = append(all, tile{tx, ty})
			}
		}
	}
	if len(all) == 0 {
		return ExperimentalV4PilotDetection{}
	}
	indices := []int{0, len(all) - 1, len(all) / 3, 2 * len(all) / 3}
	selected := make([]tile, 0, 4)
	seen := map[int]bool{}
	for _, idx := range indices {
		if idx < 0 || idx >= len(all) || seen[idx] {
			continue
		}
		seen[idx] = true
		selected = append(selected, all[idx])
	}
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var sums [residues]float64
	var absSums [residues]float64
	var counts [residues]int
	for _, t := range selected {
		baseX := t.x * experimentalV4TileWidthBlocks
		baseY := t.y * experimentalV4TileHeightBlocks
		for y := 0; y < experimentalV4TileHeightBlocks; y++ {
			for x := 0; x < experimentalV4TileWidthBlocks; x++ {
				v, ok := readProjectiveBlockValue(plane, h, (baseX+x)*blockSize, (baseY+y)*blockSize, blockSize)
				if !ok {
					continue
				}
				i := y*experimentalV4TileWidthBlocks + x
				sums[i] += v
				absSums[i] += math.Abs(v)
				counts[i]++
			}
		}
	}
	return experimentalV4DetectPilotFromResidues(candidate, sums[:], absSums[:], counts[:])
}

func experimentalV4PhoneBuild41PolishShapeFull(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, partition int) (experimentalV4PhoneHypothesis, int) {
	q, h := start.quad, start.h
	det := experimentalV4PhoneSpatialPilotDetection(plane, candidate, cw, ch, partition, h)
	evals := 1
	if !det.Available {
		return start, evals
	}
	score := det.Score
	apply := func(in [4]ImagePoint, dim int, d float64) [4]ImagePoint {
		out := in
		c := dim / 2
		if dim%2 == 0 {
			out[c].X += d
		} else {
			out[c].Y += d
		}
		return out
	}
	for _, step := range []float64{6, 3, 1.5} {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 8; dim++ {
				bq, bh, bd, bs := q, h, det, score
				for _, sg := range []float64{-1, 1} {
					qq := apply(q, dim, sg*step)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					dd := experimentalV4PhoneSpatialPilotDetection(plane, candidate, cw, ch, partition, hh)
					evals++
					if dd.Available && dd.Score > bs+1e-7 {
						bq, bh, bd, bs = qq, hh, dd, dd.Score
					}
				}
				if bs > score+1e-7 {
					q, h, det, score = bq, bh, bd, bs
					improved = true
				}
			}
			if !improved {
				break
			}
		}
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals
}

func experimentalV4PhoneBuild41PhaseBank(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, shape experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneHypothesis, int) {
	type item struct {
		dx, dy, score float64
		q             [4]ImagePoint
		h             homography
	}
	evals := 0
	evaluate := func(dx, dy float64) (item, bool) {
		qq := experimentalV4PhoneBuild41Translate(shape.quad, dx, dy)
		if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
			return item{}, false
		}
		hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
		if !ok {
			return item{}, false
		}
		s, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
		evals++
		if math.IsInf(s, -1) {
			return item{}, false
		}
		return item{dx: dx, dy: dy, score: s, q: qq, h: hh}, true
	}
	coarse := make([]item, 0, 81)
	for dy := -16.0; dy <= 16; dy += 4 {
		for dx := -16.0; dx <= 16; dx += 4 {
			if it, ok := evaluate(dx, dy); ok {
				coarse = append(coarse, it)
			}
		}
	}
	sort.SliceStable(coarse, func(i, j int) bool { return coarse[i].score > coarse[j].score })
	centers := make([]item, 0, 3)
	for _, it := range coarse {
		distinct := true
		for _, c := range centers {
			if math.Hypot(it.dx-c.dx, it.dy-c.dy) < 4.0 {
				distinct = false
				break
			}
		}
		if distinct {
			centers = append(centers, it)
			if len(centers) == 3 {
				break
			}
		}
	}
	bank := append([]item(nil), coarse...)
	for _, c := range centers {
		for dy := c.dy - 3; dy <= c.dy+3; dy += 1 {
			for dx := c.dx - 3; dx <= c.dx+3; dx += 1 {
				if it, ok := evaluate(dx, dy); ok {
					bank = append(bank, it)
				}
			}
		}
	}
	sort.SliceStable(bank, func(i, j int) bool { return bank[i].score > bank[j].score })
	out := make([]experimentalV4PhoneHypothesis, 0, 6)
	for _, it := range bank {
		near := false
		for _, q := range out {
			if experimentalV4PhoneBuild41QuadDistance(it.q, q.quad) < 1.5 {
				near = true
				break
			}
		}
		if near {
			continue
		}
		out = append(out, experimentalV4PhoneHypothesis{quad: it.q, h: it.h, proposal: it.score})
		if len(out) >= 6 {
			break
		}
	}
	return out, evals
}

// Cyclic proposal score using only the two non-held-out spatial folds. The
// held-out fold is never sampled here.
func experimentalV4PhoneBuild41FoldCyclicDetection(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, heldout int) ExperimentalV4PilotDetection {
	if plane == nil {
		return ExperimentalV4PilotDetection{}
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return ExperimentalV4PilotDetection{}
	}
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var sums [residues]float64
	var abss [residues]float64
	var counts [residues]int
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			if experimentalV4PhoneBuild41FoldForTile(tx, ty) == heldout {
				continue
			}
			baseX := tx * experimentalV4TileWidthBlocks
			baseY := ty * experimentalV4TileHeightBlocks
			for y := 0; y < experimentalV4TileHeightBlocks; y++ {
				for x := 0; x < experimentalV4TileWidthBlocks; x++ {
					v, ok := readProjectiveBlockValue(plane, h, (baseX+x)*blockSize, (baseY+y)*blockSize, blockSize)
					if !ok {
						continue
					}
					i := y*experimentalV4TileWidthBlocks + x
					sums[i] += v
					abss[i] += math.Abs(v)
					counts[i]++
				}
			}
		}
	}
	return experimentalV4DetectPilotFromResidues(candidate, sums[:], abss[:], counts[:])
}

func experimentalV4PhoneBuild41RefineShapeFold(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, heldout int) (experimentalV4PhoneHypothesis, int) {
	q := anchor
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return experimentalV4PhoneHypothesis{}, 0
	}
	det := experimentalV4PhoneBuild41FoldCyclicDetection(plane, candidate, cw, ch, h, heldout)
	evals := 1
	if !det.Available {
		return experimentalV4PhoneHypothesis{}, evals
	}
	score := det.Score
	apply := func(in [4]ImagePoint, dim int, delta float64) [4]ImagePoint {
		out := in
		corner := dim / 2
		if dim%2 == 0 {
			out[corner].X += delta
		} else {
			out[corner].Y += delta
		}
		return out
	}
	for _, step := range []float64{24, 12, 6, 3, 1.5} {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 8; dim++ {
				bestQ, bestH, bestDet, best := q, h, det, score
				for _, sign := range []float64{-1, 1} {
					qq := apply(q, dim, sign*step)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					dd := experimentalV4PhoneBuild41FoldCyclicDetection(plane, candidate, cw, ch, hh, heldout)
					evals++
					if dd.Available && dd.Score > best+1e-7 {
						bestQ, bestH, bestDet, best = qq, hh, dd, dd.Score
					}
				}
				if best > score+1e-7 {
					q, h, det, score = bestQ, bestH, bestDet, best
					improved = true
				}
			}
			if !improved {
				break
			}
		}
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals
}

func experimentalV4PhoneBuild41PolishShapeFold(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) (experimentalV4PhoneHypothesis, int) {
	q, h := start.quad, start.h
	det := experimentalV4PhoneBuild41FoldCyclicDetection(plane, candidate, cw, ch, h, heldout)
	evals := 1
	if !det.Available {
		return start, evals
	}
	score := det.Score
	apply := func(in [4]ImagePoint, dim int, delta float64) [4]ImagePoint {
		out := in
		corner := dim / 2
		if dim%2 == 0 {
			out[corner].X += delta
		} else {
			out[corner].Y += delta
		}
		return out
	}
	for _, step := range []float64{6, 3, 1.5} {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 8; dim++ {
				bestQ, bestH, bestDet, best := q, h, det, score
				for _, sign := range []float64{-1, 1} {
					qq := apply(q, dim, sign*step)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					dd := experimentalV4PhoneBuild41FoldCyclicDetection(plane, candidate, cw, ch, hh, heldout)
					evals++
					if dd.Available && dd.Score > best+1e-7 {
						bestQ, bestH, bestDet, best = qq, hh, dd, dd.Score
					}
				}
				if best > score+1e-7 {
					q, h, det, score = bestQ, bestH, bestDet, best
					improved = true
				}
			}
			if !improved {
				break
			}
		}
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals
}
