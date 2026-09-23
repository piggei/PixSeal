package watermark

import (
	"image"
	"image/color"
	"math"
	"sort"
)

const (
	experimentalV4PhoneBuild43MaxFrozen = 32
	experimentalV4PhoneBuild43PairKeep  = 2
)

type experimentalV4PhoneBuild43Telemetry struct {
	Attempted           bool
	EdgeRefined         bool
	PairsScanned        int
	PairsSelected       int
	GeometryEvaluations int
	ProposalCandidates  int
	FrozenCandidates    int
	QualifiedCandidates int
	Pair0               string
	Pair1               string
}

type experimentalV4PhoneBuild43Line struct {
	p, u ImagePoint
}

type experimentalV4PhoneBuild43Pair struct {
	a, b int
	name string
}

type experimentalV4PhoneBuild43Cell struct {
	oa, ob float64
	mean   float64
	robust float64
}

var experimentalV4PhoneBuild43Pairs = []experimentalV4PhoneBuild43Pair{
	{0, 1, "top+bottom"},
	{0, 2, "top+left"},
	{0, 3, "top+right"},
	{1, 2, "bottom+left"},
	{1, 3, "bottom+right"},
	{2, 3, "left+right"},
}

func experimentalV4PhoneBuild43LineFrom(a, b ImagePoint) experimentalV4PhoneBuild43Line {
	dx, dy := b.X-a.X, b.Y-a.Y
	l := math.Hypot(dx, dy)
	if l <= 0 {
		return experimentalV4PhoneBuild43Line{p: a, u: ImagePoint{X: 1}}
	}
	// Rotate around the side center. Rotating around an endpoint silently adds
	// a large translation on long phone edges and was rejected during Build43.
	return experimentalV4PhoneBuild43Line{
		p: ImagePoint{X: (a.X + b.X) / 2, Y: (a.Y + b.Y) / 2},
		u: ImagePoint{X: dx / l, Y: dy / l},
	}
}

func experimentalV4PhoneBuild43TransformLine(l experimentalV4PhoneBuild43Line, offset, angleDeg float64) experimentalV4PhoneBuild43Line {
	nx, ny := -l.u.Y, l.u.X
	p := ImagePoint{X: l.p.X + offset*nx, Y: l.p.Y + offset*ny}
	th := angleDeg * math.Pi / 180
	c, s := math.Cos(th), math.Sin(th)
	return experimentalV4PhoneBuild43Line{
		p: p,
		u: ImagePoint{X: l.u.X*c - nx*s, Y: l.u.Y*c - ny*s},
	}
}

func experimentalV4PhoneBuild43Intersect(a, b experimentalV4PhoneBuild43Line) (ImagePoint, bool) {
	cross := a.u.X*b.u.Y - a.u.Y*b.u.X
	if math.Abs(cross) < 1e-8 {
		return ImagePoint{}, false
	}
	dx, dy := b.p.X-a.p.X, b.p.Y-a.p.Y
	t := (dx*b.u.Y - dy*b.u.X) / cross
	return ImagePoint{X: a.p.X + t*a.u.X, Y: a.p.Y + t*a.u.Y}, true
}

func experimentalV4PhoneBuild43Quad(lines [4]experimentalV4PhoneBuild43Line) ([4]ImagePoint, bool) {
	tl, ok0 := experimentalV4PhoneBuild43Intersect(lines[0], lines[2])
	tr, ok1 := experimentalV4PhoneBuild43Intersect(lines[0], lines[3])
	bl, ok2 := experimentalV4PhoneBuild43Intersect(lines[1], lines[2])
	br, ok3 := experimentalV4PhoneBuild43Intersect(lines[1], lines[3])
	return [4]ImagePoint{tl, tr, bl, br}, ok0 && ok1 && ok2 && ok3
}

func experimentalV4PhoneBuild43BaseLines(q [4]ImagePoint) [4]experimentalV4PhoneBuild43Line {
	return [4]experimentalV4PhoneBuild43Line{
		experimentalV4PhoneBuild43LineFrom(q[0], q[1]),
		experimentalV4PhoneBuild43LineFrom(q[2], q[3]),
		experimentalV4PhoneBuild43LineFrom(q[0], q[2]),
		experimentalV4PhoneBuild43LineFrom(q[1], q[3]),
	}
}

func experimentalV4PhoneBuild43Pixel(img image.Image, x, y float64) color.NRGBA {
	b := img.Bounds()
	ix, iy := int(math.Round(x)), int(math.Round(y))
	if ix < b.Min.X {
		ix = b.Min.X
	}
	if ix >= b.Max.X {
		ix = b.Max.X - 1
	}
	if iy < b.Min.Y {
		iy = b.Min.Y
	}
	if iy >= b.Max.Y {
		iy = b.Max.Y - 1
	}
	return color.NRGBAModel.Convert(img.At(ix, iy)).(color.NRGBA)
}

func experimentalV4PhoneBuild43StructuralMetric(img image.Image, a, b ImagePoint, offset, angleDeg float64) float64 {
	mx, my := (a.X+b.X)/2, (a.Y+b.Y)/2
	dx, dy := b.X-a.X, b.Y-a.Y
	length := math.Hypot(dx, dy)
	if length < 10 {
		return 0
	}
	ux, uy := dx/length, dy/length
	nx, ny := -uy, ux
	th := angleDeg * math.Pi / 180
	c, s := math.Cos(th), math.Sin(th)
	rux, ruy := ux*c-nx*s, uy*c-ny*s
	rnx, rny := -ruy, rux
	vals := make([]float64, 0, 39)
	for i := 5; i <= 43; i++ {
		t := float64(i)/48 - .5
		x := mx + t*length*rux + offset*rnx
		y := my + t*length*ruy + offset*rny
		c1 := experimentalV4PhoneBuild43Pixel(img, x-10*rnx, y-10*rny)
		c2 := experimentalV4PhoneBuild43Pixel(img, x+10*rnx, y+10*rny)
		dr := float64(c1.R) - float64(c2.R)
		dg := float64(c1.G) - float64(c2.G)
		db := float64(c1.B) - float64(c2.B)
		vals = append(vals, math.Sqrt(dr*dr+dg*dg+db*db))
	}
	sort.Float64s(vals)
	return vals[len(vals)/4]
}

func experimentalV4PhoneBuild43StructuralLines(img image.Image, q [4]ImagePoint) ([4]experimentalV4PhoneBuild43Line, [4]bool) {
	lines := experimentalV4PhoneBuild43BaseLines(q)
	ends := [4][2]ImagePoint{{q[0], q[1]}, {q[2], q[3]}, {q[0], q[2]}, {q[1], q[3]}}
	var accepted [4]bool
	for side := 0; side < 4; side++ {
		base := experimentalV4PhoneBuild43StructuralMetric(img, ends[side][0], ends[side][1], 0, 0)
		best, bestOff, bestAng := base, 0.0, 0.0
		for ang := -.8; ang <= .8001; ang += .2 {
			for off := -40.; off <= 40.001; off += 2 {
				v := experimentalV4PhoneBuild43StructuralMetric(img, ends[side][0], ends[side][1], off, ang)
				if v > best {
					best, bestOff, bestAng = v, off, ang
				}
			}
		}
		// Structure is only a proposal seed. Reject weak improvements so bright
		// artwork texture cannot replace an already-good Build41 side.
		if best >= base*1.5 && best-base >= 4 {
			lines[side] = experimentalV4PhoneBuild43TransformLine(lines[side], bestOff, bestAng)
			accepted[side] = true
		}
	}
	return lines, accepted
}

// Fold score for proposal partitions 1 and 2. Fold 3 (held-out) is never
// sampled anywhere in Build43 candidate generation or ranking.
func experimentalV4PhoneBuild43ProposalFold(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, q [4]ImagePoint, fold int) float64 {
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return math.Inf(-1)
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	num, den := 0.0, 0.0
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			if experimentalV4PhoneBuild41FoldForTile(tx, ty) != fold {
				continue
			}
			baseX, baseY := tx*experimentalV4TileWidthBlocks, ty*experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				v, ok := readProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
				if !ok {
					continue
				}
				num += float64(candidate.signs[i]) * v
				den += math.Abs(v)
			}
		}
	}
	if den <= 0 {
		return math.Inf(-1)
	}
	return num / den
}

func experimentalV4PhoneBuild43PairSeed(base, structural [4]experimentalV4PhoneBuild43Line, structOK [4]bool, pair experimentalV4PhoneBuild43Pair) [4]experimentalV4PhoneBuild43Line {
	seed := base
	for side := 0; side < 4; side++ {
		if side != pair.a && side != pair.b && structOK[side] {
			seed[side] = structural[side]
		}
	}
	return seed
}

func experimentalV4PhoneBuild43SparsePhaseScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, q [4]ImagePoint) float64 {
	best := math.Inf(-1)
	for _, dy := range []float64{-8, 0, 8} {
		for _, dx := range []float64{-8, 0, 8} {
			qq := experimentalV4PhoneBuild41Translate(q, dx, dy)
			h, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
			if !ok {
				continue
			}
			score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, 0, true, true)
			if score > best {
				best = score
			}
		}
	}
	return best
}

func experimentalV4PhoneBuild43PairRobustScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed [4]experimentalV4PhoneBuild43Line, pair experimentalV4PhoneBuild43Pair) (float64, int) {
	robust := make([]float64, 0, 1225)
	evals := 0
	for oa := -48.; oa <= 48; oa += 16 {
		for ob := -48.; ob <= 48; ob += 16 {
			for aa := -.8; aa <= .8001; aa += .4 {
				for ab := -.8; ab <= .8001; ab += .4 {
					ls := seed
					ls[pair.a] = experimentalV4PhoneBuild43TransformLine(seed[pair.a], oa, aa)
					ls[pair.b] = experimentalV4PhoneBuild43TransformLine(seed[pair.b], ob, ab)
					q, ok := experimentalV4PhoneBuild43Quad(ls)
					if !ok || !experimentalV4PhoneBuild41WithinLimit(q, anchor) {
						continue
					}
					f1 := experimentalV4PhoneBuild43ProposalFold(plane, candidate, cw, ch, q, 1)
					f2 := experimentalV4PhoneBuild43ProposalFold(plane, candidate, cw, ch, q, 2)
					evals += 2
					if math.IsInf(f1, -1) || math.IsInf(f2, -1) {
						continue
					}
					robust = append(robust, math.Min(f1, f2))
				}
			}
		}
	}
	if len(robust) == 0 {
		return math.Inf(-1), evals
	}
	sort.Slice(robust, func(i, j int) bool { return robust[i] > robust[j] })
	n := 5
	if len(robust) < n {
		n = len(robust)
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += robust[i]
	}
	return sum / float64(n), evals
}

func experimentalV4PhoneBuild43Cells(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed [4]experimentalV4PhoneBuild43Line, pair experimentalV4PhoneBuild43Pair) ([]experimentalV4PhoneBuild43Cell, int) {
	cells := make([]experimentalV4PhoneBuild43Cell, 0, 49)
	evals := 0
	for oa := -48.; oa <= 48; oa += 16 {
		for ob := -48.; ob <= 48; ob += 16 {
			ls := seed
			ls[pair.a] = experimentalV4PhoneBuild43TransformLine(seed[pair.a], oa, 0)
			ls[pair.b] = experimentalV4PhoneBuild43TransformLine(seed[pair.b], ob, 0)
			q, ok := experimentalV4PhoneBuild43Quad(ls)
			if !ok || !experimentalV4PhoneBuild41WithinLimit(q, anchor) {
				continue
			}
			f1 := experimentalV4PhoneBuild43ProposalFold(plane, candidate, cw, ch, q, 1)
			f2 := experimentalV4PhoneBuild43ProposalFold(plane, candidate, cw, ch, q, 2)
			evals += 2
			if math.IsInf(f1, -1) || math.IsInf(f2, -1) {
				continue
			}
			cells = append(cells, experimentalV4PhoneBuild43Cell{oa: oa, ob: ob, mean: (f1 + f2) / 2, robust: math.Min(f1, f2)})
		}
	}
	sort.SliceStable(cells, func(i, j int) bool {
		if cells[i].mean == cells[j].mean {
			return cells[i].robust > cells[j].robust
		}
		return cells[i].mean > cells[j].mean
	})
	return cells, evals
}

func experimentalV4PhoneBuild43BasinCandidates(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed [4]experimentalV4PhoneBuild43Line, pair experimentalV4PhoneBuild43Pair, cell experimentalV4PhoneBuild43Cell) ([]experimentalV4PhoneHypothesis, int) {
	evals := 0
	type angular struct {
		score float64
		lines [4]experimentalV4PhoneBuild43Line
	}
	angulars := make([]angular, 0, 81)
	for aa := -.8; aa <= .8001; aa += .2 {
		for ab := -.8; ab <= .8001; ab += .2 {
			ls := seed
			ls[pair.a] = experimentalV4PhoneBuild43TransformLine(seed[pair.a], cell.oa, aa)
			ls[pair.b] = experimentalV4PhoneBuild43TransformLine(seed[pair.b], cell.ob, ab)
			q, ok := experimentalV4PhoneBuild43Quad(ls)
			if !ok || !experimentalV4PhoneBuild41WithinLimit(q, anchor) {
				continue
			}
			s := experimentalV4PhoneBuild43SparsePhaseScore(plane, candidate, cw, ch, q)
			evals += 9
			angulars = append(angulars, angular{score: s, lines: ls})
		}
	}
	sort.SliceStable(angulars, func(i, j int) bool { return angulars[i].score > angulars[j].score })
	out := make([]experimentalV4PhoneHypothesis, 0, 8)
	for _, ar := range []int{0, 2, 9, 11} {
		if ar >= len(angulars) {
			continue
		}
		a := angulars[ar]
		type fine struct {
			score float64
			lines [4]experimentalV4PhoneBuild43Line
		}
		fines := make([]fine, 0, 81)
		for da := -6.; da <= 6.001; da += 1.5 {
			for db := -6.; db <= 6.001; db += 1.5 {
				ls := a.lines
				ls[pair.a] = experimentalV4PhoneBuild43TransformLine(a.lines[pair.a], da, 0)
				ls[pair.b] = experimentalV4PhoneBuild43TransformLine(a.lines[pair.b], db, 0)
				q, ok := experimentalV4PhoneBuild43Quad(ls)
				if !ok || !experimentalV4PhoneBuild41WithinLimit(q, anchor) {
					continue
				}
				h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
				if !ok {
					continue
				}
				s, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, 0, true, false)
				evals++
				fines = append(fines, fine{score: s, lines: ls})
			}
		}
		sort.SliceStable(fines, func(i, j int) bool { return fines[i].score > fines[j].score })
		for _, fr := range []int{3, 10} {
			if fr >= len(fines) {
				continue
			}
			fz := fines[fr]
			comps := make([]int, 0, 2)
			for side := 0; side < 4; side++ {
				if side != pair.a && side != pair.b {
					comps = append(comps, side)
				}
			}
			type child struct {
				score float64
				q     [4]ImagePoint
				h     homography
			}
			children := make([]child, 0, 9)
			for _, ca := range []float64{-8, 0, 8} {
				for _, cb := range []float64{-8, 0, 8} {
					ls := fz.lines
					ls[comps[0]] = experimentalV4PhoneBuild43TransformLine(ls[comps[0]], ca, 0)
					ls[comps[1]] = experimentalV4PhoneBuild43TransformLine(ls[comps[1]], cb, 0)
					q, ok := experimentalV4PhoneBuild43Quad(ls)
					if !ok || !experimentalV4PhoneBuild41WithinLimit(q, anchor) {
						continue
					}
					h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
					if !ok {
						continue
					}
					s, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, 0, true, false)
					evals++
					children = append(children, child{score: s, q: q, h: h})
				}
			}
			sort.SliceStable(children, func(i, j int) bool { return children[i].score > children[j].score })
			if len(children) == 0 {
				continue
			}
			shape := experimentalV4PhoneHypothesis{quad: children[0].q, h: children[0].h, proposal: children[0].score}
			aligned, n := experimentalV4PhoneBuild41PhaseAlign(plane, candidate, cw, ch, anchor, shape, 0)
			evals += n
			out = append(out, aligned)
		}
	}
	return out, evals
}

func experimentalV4PhoneSearchBuild43(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, experimentalV4PhoneBuild43Telemetry) {
	tele := experimentalV4PhoneBuild43Telemetry{Attempted: true}
	if work == nil || !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, tele
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)
	base := experimentalV4PhoneBuild43BaseLines(anchor)
	structural, structOK := experimentalV4PhoneBuild43StructuralLines(work, anchor)
	for _, ok := range structOK {
		if ok {
			tele.EdgeRefined = true
			break
		}
	}
	plane := newPixelPlane(work)
	candidate := experimentalV4Prototype2Candidate()
	type rankedPair struct {
		pair  experimentalV4PhoneBuild43Pair
		score float64
		seed  [4]experimentalV4PhoneBuild43Line
	}
	ranked := make([]rankedPair, 0, 6)
	for _, pair := range experimentalV4PhoneBuild43Pairs {
		seed := experimentalV4PhoneBuild43PairSeed(base, structural, structOK, pair)
		score, n := experimentalV4PhoneBuild43PairRobustScore(plane, candidate, cw, ch, anchor, seed, pair)
		tele.GeometryEvaluations += n
		tele.PairsScanned++
		if !math.IsInf(score, -1) {
			ranked = append(ranked, rankedPair{pair: pair, score: score, seed: seed})
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	if len(ranked) > experimentalV4PhoneBuild43PairKeep {
		ranked = ranked[:experimentalV4PhoneBuild43PairKeep]
	}
	tele.PairsSelected = len(ranked)
	if len(ranked) > 0 {
		tele.Pair0 = ranked[0].pair.name
	}
	if len(ranked) > 1 {
		tele.Pair1 = ranked[1].pair.name
	}
	frozen := make([]experimentalV4PhoneHypothesis, 0, experimentalV4PhoneBuild43MaxFrozen)
	for _, rp := range ranked {
		cells, n := experimentalV4PhoneBuild43Cells(plane, candidate, cw, ch, anchor, rp.seed, rp.pair)
		tele.GeometryEvaluations += n
		for _, cr := range []int{0, 3} {
			if cr >= len(cells) {
				continue
			}
			bank, n := experimentalV4PhoneBuild43BasinCandidates(plane, candidate, cw, ch, anchor, rp.seed, rp.pair, cells[cr])
			tele.GeometryEvaluations += n
			for _, h := range bank {
				tele.ProposalCandidates++
				if len(frozen) >= experimentalV4PhoneBuild43MaxFrozen {
					break
				}
				frozen = append(frozen, h)
			}
		}
	}
	tele.FrozenCandidates = len(frozen)
	qualified := make([]experimentalV4PhoneHypothesis, 0, len(frozen))
	for _, h := range frozen {
		q, n, ok := experimentalV4PhoneBuild41Qualify(plane, work, candidate, cw, ch, h, 0)
		tele.GeometryEvaluations += n
		if ok {
			qualified = append(qualified, q)
		}
	}
	tele.QualifiedCandidates = len(qualified)
	return qualified, tele
}
