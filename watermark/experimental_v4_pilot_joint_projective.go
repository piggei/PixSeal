package watermark

import (
	"image"
	"math"
	"runtime"
	"sort"
	"sync"
)

// Build29 extends the public Format-v4 geometry research to cases where mild
// projective distortion and placement are simultaneously unknown. Geometry
// proposal is deliberately spatially staged: one middle-row tile finds cyclic
// origin, the remaining middle-row tile repetitions rank a basin, and top/bottom
// repetitions are held out until the final geometry decision. Payload/header,
// key material, ECC and HMAC are never consulted.
type experimentalV4JointProjectiveResult struct {
	available           bool
	params              experimentalV4BlindGeometryParams
	placement           experimentalV4PlacementResult
	detection           ExperimentalV4PilotDetection
	proposalObjective   float64
	geometryValidation  float64
	structuralScore     float64
	hypothesesEvaluated int
}

type experimentalV4JointProjectiveCandidate struct {
	params           experimentalV4BlindGeometryParams
	structural       float64
	proposal         float64
	validation       float64
	h                homography
	originX, originY int
	fullW, fullH     int
}

func experimentalV4JointProjectiveWorkers() int {
	n := runtime.GOMAXPROCS(0)
	if n > 8 {
		n = 8
	}
	if n < 1 {
		n = 1
	}
	return n
}
func experimentalV4DetectPilotSingleTilePlane(plane *pixelPlane, candidate experimentalV4PilotCandidate, startX, startY int, h homography) ExperimentalV4PilotDetection {
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var values [residues]float64
	var valid [residues]bool
	for y := 0; y < experimentalV4TileHeightBlocks; y++ {
		for x := 0; x < experimentalV4TileWidthBlocks; x++ {
			v, ok := readProjectiveBlockValue(plane, h, (startX+x)*blockSize, (startY+y)*blockSize, blockSize)
			if !ok {
				continue
			}
			idx := y*experimentalV4TileWidthBlocks + x
			values[idx], valid[idx] = v, true
		}
	}
	best, runner := math.Inf(-1), math.Inf(-1)
	bestX, bestY, bestVisible := 0, 0, 0
	for oy := 0; oy < experimentalV4TileHeightBlocks; oy++ {
		for ox := 0; ox < experimentalV4TileWidthBlocks; ox++ {
			num, den := 0.0, 0.0
			visible := 0
			for i, pos := range candidate.positions {
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				rx := positiveMod(px-ox, experimentalV4TileWidthBlocks)
				ry := positiveMod(py-oy, experimentalV4TileHeightBlocks)
				idx := ry*experimentalV4TileWidthBlocks + rx
				if !valid[idx] {
					continue
				}
				visible++
				num += float64(candidate.signs[i]) * values[idx]
				den += math.Abs(values[idx])
			}
			if den <= 0 {
				continue
			}
			s := num / den
			if s > best {
				runner, best = best, s
				bestX, bestY, bestVisible = ox, oy, visible
			} else if s > runner {
				runner = s
			}
		}
	}
	if math.IsInf(best, -1) {
		return ExperimentalV4PilotDetection{}
	}
	if math.IsInf(runner, -1) {
		runner = 0
	}
	return ExperimentalV4PilotDetection{Available: true, BlockSizePixels: blockSize, OriginXBlocks: bestX, OriginYBlocks: bestY, Score: best, RunnerUpScore: runner, Margin: best - runner, VisiblePilotPositions: bestVisible, PilotSamples: bestVisible}
}

// fixed-origin score over selected complete tile rows. Geometry proposal uses
// only the middle row; validation uses top and bottom rows and is therefore
// spatially disjoint whenever at least three tile rows are visible canonically.
func experimentalV4PilotFixedOriginRowsPlane(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography, originX, originY int, middle bool) (float64, int) {
	bw, bh := canonicalWidth/blockSize, canonicalHeight/blockSize
	tilesX := bw / experimentalV4TileWidthBlocks
	tilesY := bh / experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return math.Inf(-1), 0
	}
	midY := tilesY / 2
	num, den := 0.0, 0.0
	visible := 0
	for ty := 0; ty < tilesY; ty++ {
		use := false
		if middle {
			use = ty == midY
		} else {
			use = ty == 0 || ty == tilesY-1
		}
		if !use {
			continue
		}
		for tx := 0; tx < tilesX; tx++ {
			baseX := tx * experimentalV4TileWidthBlocks
			baseY := ty * experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				rx := positiveMod(px-originX, experimentalV4TileWidthBlocks)
				ry := positiveMod(py-originY, experimentalV4TileHeightBlocks)
				v, ok := readProjectiveBlockValue(plane, h, (baseX+rx)*blockSize, (baseY+ry)*blockSize, blockSize)
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

func experimentalV4JointProjectiveCenteredShift(fullW, fullH, obsW, obsH int) (float64, float64) {
	return 0.5 * float64(obsW-fullW), 0.5 * float64(obsH-fullH)
}

// Geometry proposal at a fixed projective family. Whole-block placement error
// is absorbed by the cyclic origin search; only a bounded sub-block phase cross
// is scanned around the center placement implied by frame dimensions.
func experimentalV4ProjectiveLocalPilotScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) (float64, int, int, int) {
	bw, bh := canonicalWidth/blockSize, canonicalHeight/blockSize
	if bw < experimentalV4TileWidthBlocks || bh < experimentalV4TileHeightBlocks {
		return math.Inf(-1), 0, 0, 0
	}
	startX := ((bw - experimentalV4TileWidthBlocks) / 2 / experimentalV4TileWidthBlocks) * experimentalV4TileWidthBlocks
	startY := ((bh - experimentalV4TileHeightBlocks) / 2 / experimentalV4TileHeightBlocks) * experimentalV4TileHeightBlocks
	best := math.Inf(-1)
	bestX, bestY, bestVisible := 0, 0, 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			ox := positiveMod(dx, experimentalV4TileWidthBlocks)
			oy := positiveMod(dy, experimentalV4TileHeightBlocks)
			num, den := 0.0, 0.0
			visible := 0
			for i, pos := range candidate.positions {
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				rx := positiveMod(px-ox, experimentalV4TileWidthBlocks)
				ry := positiveMod(py-oy, experimentalV4TileHeightBlocks)
				v, ok := readProjectiveBlockValue(plane, h, (startX+rx)*blockSize, (startY+ry)*blockSize, blockSize)
				if !ok {
					continue
				}
				visible++
				num += float64(candidate.signs[i]) * v
				den += math.Abs(v)
			}
			if den <= 0 {
				continue
			}
			score := num / den
			if score > best {
				best, bestX, bestY, bestVisible = score, ox, oy, visible
			}
		}
	}
	return best, bestX, bestY, bestVisible
}

// Geometry proposal at a fixed projective family. Whole-block placement error
// is absorbed by a tiny local cyclic-origin search. Only five sub-block phases
// are evaluated, so thousands of coarse geometries remain tractable.
func experimentalV4JointProjectiveProposal(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight, obsW, obsH int, p experimentalV4BlindGeometryParams) experimentalV4JointProjectiveCandidate {
	base, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
	cx, cy := experimentalV4JointProjectiveCenteredShift(fw, fh, obsW, obsH)
	phases := [][2]float64{{0, 0}, {-4, 0}, {4, 0}, {0, -4}, {0, 4}}
	best := experimentalV4JointProjectiveCandidate{params: p, proposal: math.Inf(-1), fullW: fw, fullH: fh}
	for _, ph := range phases {
		h := experimentalV4PlacementShiftHomography(base, cx+ph[0], cy+ph[1])
		score, ox, oy, visible := experimentalV4ProjectiveLocalPilotScore(plane, candidate, canonicalWidth, canonicalHeight, h)
		if visible < 48 || math.IsInf(score, -1) {
			continue
		}
		// Same local origin must explain the complete middle row. This is a
		// spatial cross-check, not a second origin optimization.
		mid, vis := experimentalV4PilotFixedOriginRowsPlane(plane, candidate, canonicalWidth, canonicalHeight, h, ox, oy, true)
		if vis < 48 || math.IsInf(mid, -1) {
			continue
		}
		objective := math.Min(score, mid)
		if objective > best.proposal {
			best.proposal = objective
			best.h = h
			best.originX = ox
			best.originY = oy
		}
	}
	return best
}

func experimentalV4ProjectiveCheapStructural(plane *pixelPlane, canonicalWidth, canonicalHeight, obsW, obsH int, p experimentalV4BlindGeometryParams) float64 {
	base, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
	cx, cy := experimentalV4JointProjectiveCenteredShift(fw, fh, obsW, obsH)
	phases := [][2]float64{{0, 0}, {4, 0}, {0, 4}, {4, 4}}
	values := make([]float64, 0, len(phases))
	bw, bh := canonicalWidth/blockSize, canonicalHeight/blockSize
	for _, ph := range phases {
		h := experimentalV4PlacementShiftHomography(base, cx+ph[0], cy+ph[1])
		vals := make([]float64, 0, 12)
		for k := 0; k < 12; k++ {
			x := bw/4 + ((k*17 + 3) % (maxIntV4Placement(1, bw/2)))
			y := bh/4 + ((k*11 + 5) % (maxIntV4Placement(1, bh/2)))
			v, ok := readProjectiveBlockValue(plane, h, x*blockSize, y*blockSize, blockSize)
			if ok {
				vals = append(vals, math.Abs(v))
			}
		}
		if len(vals) < 8 {
			values = append(values, 0)
			continue
		}
		sort.Float64s(vals)
		sum := 0.0
		for _, v := range vals[2 : len(vals)-2] {
			sum += v
		}
		values = append(values, sum/float64(len(vals)-4))
	}
	if len(values) != 4 {
		return math.Inf(-1)
	}
	sort.Float64s(values)
	return values[3] - 0.5*(values[1]+values[2])
}

func experimentalV4JointProjectiveStructuralBank(plane *pixelPlane, canonicalWidth, canonicalHeight, obsW, obsH int) ([]experimentalV4JointProjectiveCandidate, int) {
	// Coarse public family. Insets include 0.015 explicitly; Build29 diagnostics
	// showed that omitting this half-step can erase the correct projective basin.
	angles := make([]float64, 0, 49)
	for a := -12.0; a <= 12.0001; a += 0.5 {
		angles = append(angles, a)
	}
	scalesX := []float64{0.95, 0.975, 1.0, 1.025, 1.05, 1.075, 1.10, 1.125, 1.15}
	scalesY := []float64{0.85, 0.875, 0.90, 0.925, 0.95, 0.975, 1.0, 1.025, 1.05}
	insets := []float64{0, 0.015, 0.030, 0.045, 0.060}
	type job struct{ a, sx, sy, t, b float64 }
	jobs := make(chan job, 256)
	out := make(chan experimentalV4JointProjectiveCandidate, 256)
	workers := experimentalV4JointProjectiveWorkers()
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				p := experimentalV4BlindGeometryParams{angleDeg: j.a, scaleX: j.sx, scaleY: j.sy, topInset: j.t, bottomInset: j.b}
				_, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
				// Build29 projective branch is crop-oriented. Purely dimensional public
				// bounds remove impossible hypotheses without looking at image content.
				if fw < obsW || fh < obsH || fw-obsW > 384 || fh-obsH > 384 {
					continue
				}
				s := experimentalV4JointAffinePhaseContrast(plane, canonicalWidth, canonicalHeight, p, 16)
				if math.IsInf(s, -1) {
					continue
				}
				out <- experimentalV4JointProjectiveCandidate{params: p, structural: s, fullW: fw, fullH: fh}
			}
		}()
	}
	go func() {
		for _, a := range angles {
			for _, sx := range scalesX {
				for _, sy := range scalesY {
					for _, t := range insets {
						for _, b := range insets {
							jobs <- job{a, sx, sy, t, b}
						}
					}
				}
			}
		}
		close(jobs)
		wg.Wait()
		close(out)
	}()
	all := make([]experimentalV4JointProjectiveCandidate, 0, 50000)
	for c := range out {
		all = append(all, c)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].structural > all[j].structural })
	// Diagnostics established that difficult photographic basins can sit around
	// rank ~2000 structurally. Keep a conservative 5000 for public-pilot ranking.
	if len(all) > 5000 {
		all = all[:5000]
	}
	return all, len(angles) * len(scalesX) * len(scalesY) * len(insets) * len(insets)
}

func experimentalV4ProjectiveBasinNear(a, b experimentalV4BlindGeometryParams) bool {
	return math.Abs(a.angleDeg-b.angleDeg) < 0.76 && math.Abs(a.scaleX-b.scaleX) < 0.031 && math.Abs(a.scaleY-b.scaleY) < 0.031 && math.Abs(a.topInset-b.topInset) < 0.016 && math.Abs(a.bottomInset-b.bottomInset) < 0.016
}

func experimentalV4JointProjectiveRefine(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight, obsW, obsH int, seed experimentalV4JointProjectiveCandidate) experimentalV4JointProjectiveCandidate {
	best := seed
	eval := func(p experimentalV4BlindGeometryParams) {
		q := experimentalV4JointProjectiveProposal(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, p)
		if q.proposal > best.proposal {
			q.structural = experimentalV4JointAffinePhaseContrast(plane, canonicalWidth, canonicalHeight, p, 64)
			best = q
		}
	}
	// Coordinate refinement around a basin; pilot is proposal evidence here.
	for pass := 0; pass < 2; pass++ {
		base := best.params
		for a := base.angleDeg - 0.5; a <= base.angleDeg+0.5001; a += 0.1 {
			p := best.params
			p.angleDeg = a
			eval(p)
		}
		base = best.params
		for sx := math.Max(0.90, base.scaleX-0.04); sx <= math.Min(1.18, base.scaleX+0.04)+1e-9; sx += 0.01 {
			p := best.params
			p.scaleX = sx
			eval(p)
		}
		base = best.params
		for sy := math.Max(0.82, base.scaleY-0.04); sy <= math.Min(1.08, base.scaleY+0.04)+1e-9; sy += 0.01 {
			p := best.params
			p.scaleY = sy
			eval(p)
		}
		base = best.params
		for t := math.Max(0, base.topInset-0.015); t <= math.Min(0.075, base.topInset+0.015)+1e-9; t += 0.005 {
			p := best.params
			p.topInset = t
			eval(p)
		}
		base = best.params
		for b := math.Max(0, base.bottomInset-0.015); b <= math.Min(0.075, base.bottomInset+0.015)+1e-9; b += 0.005 {
			p := best.params
			p.bottomInset = b
			eval(p)
		}
	}
	return best
}

func experimentalV4JointProjectiveSearchBuild29(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int) experimentalV4JointProjectiveResult {
	if img == nil {
		return experimentalV4JointProjectiveResult{}
	}
	plane := newPixelPlane(img)
	obsW, obsH := img.Bounds().Dx(), img.Bounds().Dy()
	structural, hyp := experimentalV4JointProjectiveStructuralBank(plane, canonicalWidth, canonicalHeight, obsW, obsH)
	if len(structural) == 0 {
		return experimentalV4JointProjectiveResult{hypothesesEvaluated: hyp}
	}
	pre := experimentalV4JointProjectivePrefilterBank(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, structural)
	hyp += len(structural) * 5 * 9
	if len(pre) > 2000 {
		pre = pre[:2000]
	}
	ranked := experimentalV4JointProjectiveFullProposals(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, pre)
	hyp += len(pre) * 5 * 9
	// Basin NMS prevents many variants of one local false maximum from crowding
	// out a distinct projective family before refinement.
	diverse := make([]experimentalV4JointProjectiveCandidate, 0, 12)
	for _, q := range ranked {
		near := false
		for _, d := range diverse {
			if experimentalV4ProjectiveBasinNear(q.params, d.params) {
				near = true
				break
			}
		}
		if near {
			continue
		}
		diverse = append(diverse, q)
		if len(diverse) >= 12 {
			break
		}
	}
	if len(diverse) == 0 {
		return experimentalV4JointProjectiveResult{hypothesesEvaluated: hyp}
	}
	refined := make([]experimentalV4JointProjectiveCandidate, 0, len(diverse))
	for _, seed := range diverse {
		q := experimentalV4JointProjectiveRefine(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, seed)
		hyp += 160
		refined = append(refined, q)
	}
	// Final geometry decision uses spatially held-out top/bottom tile rows and
	// the fixed cyclic origin proposed from the middle row.
	for i := range refined {
		v, vis := experimentalV4PilotFixedOriginRowsPlane(plane, candidate, canonicalWidth, canonicalHeight, refined[i].h, refined[i].originX, refined[i].originY, false)
		if vis < 48 {
			refined[i].validation = math.Inf(-1)
		} else {
			refined[i].validation = v
		}
	}
	sort.Slice(refined, func(i, j int) bool {
		if refined[i].validation == refined[j].validation {
			return refined[i].proposal > refined[j].proposal
		}
		return refined[i].validation > refined[j].validation
	})
	if len(refined) > 5 {
		refined = refined[:5]
	}
	best := experimentalV4JointProjectiveResult{hypothesesEvaluated: hyp}
	bestValidation := math.Inf(-1)
	for _, q := range refined {
		geometry, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, q.params)
		placement := experimentalV4PilotPlacementSearch(img, candidate, canonicalWidth, canonicalHeight, geometry, fw, fh)
		hyp += placement.hypothesesEvaluated
		if !placement.available {
			continue
		}
		d := experimentalV4DetectPilotProjective(img, candidate, canonicalWidth, canonicalHeight, placement.h)
		if placement.validationScore > bestValidation {
			bestValidation = placement.validationScore
			best = experimentalV4JointProjectiveResult{available: true, params: q.params, placement: placement, detection: d, proposalObjective: q.proposal, geometryValidation: q.validation, structuralScore: q.structural, hypothesesEvaluated: hyp}
		}
	}
	best.hypothesesEvaluated = hyp
	return best
}

func experimentalV4Build34DistinctStructuralAnchors(structural []experimentalV4JointProjectiveCandidate, limit int) []experimentalV4JointProjectiveCandidate {
	out := make([]experimentalV4JointProjectiveCandidate, 0, limit)
	for _, q := range structural {
		near := false
		for _, a := range out {
			if experimentalV4ProjectiveBasinNear(q.params, a.params) {
				near = true
				break
			}
		}
		if near {
			continue
		}
		out = append(out, q)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func experimentalV4JointProjectiveSearchBuild34(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int) experimentalV4JointProjectiveResult {
	if img == nil {
		return experimentalV4JointProjectiveResult{}
	}
	plane := newPixelPlane(img)
	obsW, obsH := img.Bounds().Dx(), img.Bounds().Dy()
	structural, hyp := experimentalV4JointProjectiveStructuralBank(plane, canonicalWidth, canonicalHeight, obsW, obsH)
	if len(structural) == 0 {
		return experimentalV4JointProjectiveResult{hypothesesEvaluated: hyp}
	}
	anchors := experimentalV4Build34DistinctStructuralAnchors(structural, 2)
	best := experimentalV4JointProjectiveCandidate{proposal: math.Inf(-1)}
	for _, anchor := range anchors {
		local := experimentalV4Build34LocalNeighborhood(anchor)
		for i := range local {
			local[i].structural = experimentalV4JointAffinePhaseContrast(plane, canonicalWidth, canonicalHeight, local[i].params, 16)
			hyp++
		}
		sort.Slice(local, func(i, j int) bool { return local[i].structural > local[j].structural })
		if len(local) > 64 {
			local = local[:64]
		}
		for _, seed := range local {
			q := experimentalV4Build34CenteredProposal(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, seed.params)
			hyp += 9
			q.structural = seed.structural
			if q.proposal > best.proposal {
				best = q
			}
		}
	}
	if math.IsInf(best.proposal, -1) {
		return experimentalV4JointProjectiveResult{hypothesesEvaluated: hyp}
	}
	canonH := experimentalV4Build34CanonicalizeOrigin(best.h, best.originX, best.originY)
	detection := experimentalV4DetectPilotProjective(img, candidate, canonicalWidth, canonicalHeight, canonH)
	validation, vis := experimentalV4PilotFixedOriginRowsPlane(plane, candidate, canonicalWidth, canonicalHeight, canonH, 0, 0, false)
	placement := experimentalV4PlacementResult{
		available:           !math.IsInf(validation, -1) && vis >= 48,
		h:                   canonH,
		proposalScore:       best.proposal,
		validationScore:     validation,
		visibleValidation:   vis,
		hypothesesEvaluated: hyp,
	}
	return experimentalV4JointProjectiveResult{
		available:           placement.available,
		params:              best.params,
		placement:           placement,
		detection:           detection,
		proposalObjective:   best.proposal,
		geometryValidation:  validation,
		structuralScore:     best.structural,
		hypothesesEvaluated: hyp,
	}
}

func experimentalV4JointProjectiveSearch(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int) experimentalV4JointProjectiveResult {
	// Build34 applies only when at least four complete v4 tiles are available in
	// both axes. Smaller carriers keep the frozen Build29 research behavior and
	// its conservative SAFE-REJECT boundary.
	if canonicalWidth/blockSize/experimentalV4TileWidthBlocks >= 4 && canonicalHeight/blockSize/experimentalV4TileHeightBlocks >= 4 {
		return experimentalV4JointProjectiveSearchBuild34(img, candidate, canonicalWidth, canonicalHeight)
	}
	return experimentalV4JointProjectiveSearchBuild29(img, candidate, canonicalWidth, canonicalHeight)
}

type experimentalV4ProjectiveLocalResult struct {
	score            float64
	fullScore        float64
	h                homography
	originX, originY int
}

func experimentalV4ProjectiveHalfCrossHypothesis(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight, obsW, obsH int, p experimentalV4BlindGeometryParams) experimentalV4ProjectiveLocalResult {
	base, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
	cx, cy := experimentalV4JointProjectiveCenteredShift(fw, fh, obsW, obsH)
	phases := [][2]float64{{0, 0}, {-4, 0}, {4, 0}, {0, -4}, {0, 4}}
	bw, bh := canonicalWidth/blockSize, canonicalHeight/blockSize
	startX := ((bw - experimentalV4TileWidthBlocks) / 2 / experimentalV4TileWidthBlocks) * experimentalV4TileWidthBlocks
	startY := ((bh - experimentalV4TileHeightBlocks) / 2 / experimentalV4TileHeightBlocks) * experimentalV4TileHeightBlocks
	best := experimentalV4ProjectiveLocalResult{score: math.Inf(-1)}
	for _, ph := range phases {
		h := experimentalV4PlacementShiftHomography(base, cx+ph[0], cy+ph[1])
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				ox := positiveMod(dx, experimentalV4TileWidthBlocks)
				oy := positiveMod(dy, experimentalV4TileHeightBlocks)
				var num [2]float64
				var den [2]float64
				var vis [2]int
				for i, pos := range candidate.positions {
					px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
					rx := positiveMod(px-ox, experimentalV4TileWidthBlocks)
					ry := positiveMod(py-oy, experimentalV4TileHeightBlocks)
					v, ok := readProjectiveBlockValue(plane, h, (startX+rx)*blockSize, (startY+ry)*blockSize, blockSize)
					if !ok {
						continue
					}
					part := i & 1
					vis[part]++
					num[part] += float64(candidate.signs[i]) * v
					den[part] += math.Abs(v)
				}
				if vis[0] < 24 || vis[1] < 24 || den[0] <= 0 || den[1] <= 0 {
					continue
				}
				a, b := num[0]/den[0], num[1]/den[1]
				score := math.Min(a, b)
				if score > best.score {
					best = experimentalV4ProjectiveLocalResult{score: score, fullScore: (num[0] + num[1]) / (den[0] + den[1]), h: h, originX: ox, originY: oy}
				}
			}
		}
	}
	return best
}

type experimentalV4ProjectivePreScore struct {
	c     experimentalV4JointProjectiveCandidate
	score float64
	local experimentalV4ProjectiveLocalResult
}

func experimentalV4JointProjectivePrefilterBank(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight, obsW, obsH int, structural []experimentalV4JointProjectiveCandidate) []experimentalV4ProjectivePreScore {
	out := make([]experimentalV4ProjectivePreScore, len(structural))
	workers := experimentalV4JointProjectiveWorkers()
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for i, s := range structural {
		i, s := i, s
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			local := experimentalV4ProjectiveHalfCrossHypothesis(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, s.params)
			out[i] = experimentalV4ProjectivePreScore{c: s, score: local.score, local: local}
		}()
	}
	wg.Wait()
	filtered := out[:0]
	for _, x := range out {
		if !math.IsInf(x.score, -1) {
			filtered = append(filtered, x)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].score > filtered[j].score })
	return filtered
}

func experimentalV4JointProjectiveFullProposals(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight, obsW, obsH int, pre []experimentalV4ProjectivePreScore) []experimentalV4JointProjectiveCandidate {
	out := make([]experimentalV4JointProjectiveCandidate, len(pre))
	workers := experimentalV4JointProjectiveWorkers()
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for i, x := range pre {
		i, x := i, x
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			mid, vis := experimentalV4PilotFixedOriginRowsPlane(plane, candidate, canonicalWidth, canonicalHeight, x.local.h, x.local.originX, x.local.originY, true)
			q := experimentalV4JointProjectiveCandidate{params: x.c.params, structural: x.c.structural, proposal: math.Inf(-1), h: x.local.h, originX: x.local.originX, originY: x.local.originY, fullW: x.c.fullW, fullH: x.c.fullH}
			if vis >= 48 && !math.IsInf(mid, -1) {
				q.proposal = math.Min(x.local.fullScore, mid)
			}
			out[i] = q
		}()
	}
	wg.Wait()
	filtered := out[:0]
	for _, q := range out {
		if !math.IsInf(q.proposal, -1) {
			filtered = append(filtered, q)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].proposal > filtered[j].proposal })
	return filtered
}

func experimentalV4JointPaddedAffineStructuralBank(plane *pixelPlane, canonicalWidth, canonicalHeight, obsW, obsH int) ([]experimentalV4JointProjectiveCandidate, int) {
	angles := make([]float64, 0, 97)
	for a := -12.0; a <= 12.0001; a += 0.25 {
		angles = append(angles, a)
	}
	scales := make([]float64, 0, 25)
	for s := 0.88; s <= 1.1201; s += 0.01 {
		scales = append(scales, s)
	}
	type job struct{ a, sx, sy float64 }
	jobs := make(chan job, 256)
	out := make(chan experimentalV4JointProjectiveCandidate, 256)
	workers := experimentalV4JointProjectiveWorkers()
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				p := experimentalV4BlindGeometryParams{angleDeg: j.a, scaleX: j.sx, scaleY: j.sy}
				_, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
				if fw > obsW || fh > obsH || obsW-fw > 384 || obsH-fh > 384 {
					continue
				}
				score := experimentalV4ProjectiveCheapStructural(plane, canonicalWidth, canonicalHeight, obsW, obsH, p)
				if math.IsInf(score, -1) {
					continue
				}
				out <- experimentalV4JointProjectiveCandidate{params: p, structural: score, fullW: fw, fullH: fh}
			}
		}()
	}
	go func() {
		for _, a := range angles {
			for _, sx := range scales {
				for _, sy := range scales {
					jobs <- job{a, sx, sy}
				}
			}
		}
		close(jobs)
		wg.Wait()
		close(out)
	}()
	all := make([]experimentalV4JointProjectiveCandidate, 0, 60000)
	for q := range out {
		all = append(all, q)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].structural > all[j].structural })
	if len(all) > 5000 {
		all = all[:5000]
	}
	return all, len(angles) * len(scales) * len(scales)
}

type experimentalV4JointPaddedResult struct {
	available           bool
	params              experimentalV4BlindGeometryParams
	placement           experimentalV4PlacementResult
	detection           ExperimentalV4PilotDetection
	proposalScore       float64
	structuralScore     float64
	hypothesesEvaluated int
}

func experimentalV4JointPaddedAffineRefine(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight, obsW, obsH int, seed experimentalV4JointProjectiveCandidate) experimentalV4JointProjectiveCandidate {
	best := seed
	eval := func(p experimentalV4BlindGeometryParams) {
		p.topInset, p.bottomInset = 0, 0
		q := experimentalV4JointProjectiveProposal(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, p)
		if q.proposal > best.proposal {
			q.structural = experimentalV4ProjectiveCheapStructural(plane, canonicalWidth, canonicalHeight, obsW, obsH, p)
			best = q
		}
	}
	for pass := 0; pass < 2; pass++ {
		base := best.params
		for a := base.angleDeg - .5; a <= base.angleDeg+.5001; a += .1 {
			p := best.params
			p.angleDeg = a
			eval(p)
		}
		base = best.params
		for sx := math.Max(.85, base.scaleX-.04); sx <= math.Min(1.15, base.scaleX+.04)+1e-9; sx += .01 {
			p := best.params
			p.scaleX = sx
			eval(p)
		}
		base = best.params
		for sy := math.Max(.85, base.scaleY-.04); sy <= math.Min(1.15, base.scaleY+.04)+1e-9; sy += .01 {
			p := best.params
			p.scaleY = sy
			eval(p)
		}
	}
	return best
}

func experimentalV4JointPaddedAffineSearch(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int) experimentalV4JointPaddedResult {
	if img == nil {
		return experimentalV4JointPaddedResult{}
	}
	plane := newPixelPlane(img)
	obsW, obsH := img.Bounds().Dx(), img.Bounds().Dy()
	structural, hyp := experimentalV4JointPaddedAffineStructuralBank(plane, canonicalWidth, canonicalHeight, obsW, obsH)
	if len(structural) == 0 {
		return experimentalV4JointPaddedResult{hypothesesEvaluated: hyp}
	}
	pre := experimentalV4JointProjectivePrefilterBank(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, structural)
	hyp += len(structural) * 45
	if len(pre) > 3000 {
		pre = pre[:3000]
	}
	ranked := experimentalV4JointProjectiveFullProposals(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, pre)
	hyp += len(pre)
	div := make([]experimentalV4JointProjectiveCandidate, 0, 12)
	for _, q := range ranked {
		near := false
		for _, d := range div {
			if math.Abs(q.params.angleDeg-d.params.angleDeg) < .76 && math.Abs(q.params.scaleX-d.params.scaleX) < .031 && math.Abs(q.params.scaleY-d.params.scaleY) < .031 {
				near = true
				break
			}
		}
		if near {
			continue
		}
		div = append(div, q)
		if len(div) >= 12 {
			break
		}
	}
	if len(div) == 0 {
		return experimentalV4JointPaddedResult{hypothesesEvaluated: hyp}
	}
	ref := make([]experimentalV4JointProjectiveCandidate, 0, len(div))
	for _, s := range div {
		q := experimentalV4JointPaddedAffineRefine(plane, candidate, canonicalWidth, canonicalHeight, obsW, obsH, s)
		v, _ := experimentalV4PilotFixedOriginRowsPlane(plane, candidate, canonicalWidth, canonicalHeight, q.h, q.originX, q.originY, false)
		q.validation = v
		ref = append(ref, q)
		hyp += 120
	}
	sort.Slice(ref, func(i, j int) bool { return ref[i].validation > ref[j].validation })
	if len(ref) > 5 {
		ref = ref[:5]
	}
	best := experimentalV4JointPaddedResult{hypothesesEvaluated: hyp}
	bestVal := math.Inf(-1)
	for _, q := range ref {
		geo, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, q.params)
		pl := experimentalV4PilotPlacementSearch(img, candidate, canonicalWidth, canonicalHeight, geo, fw, fh)
		hyp += pl.hypothesesEvaluated
		if !pl.available {
			continue
		}
		d := experimentalV4DetectPilotProjective(img, candidate, canonicalWidth, canonicalHeight, pl.h)
		if pl.validationScore > bestVal {
			bestVal = pl.validationScore
			best = experimentalV4JointPaddedResult{available: true, params: q.params, placement: pl, detection: d, proposalScore: q.proposal, structuralScore: q.structural, hypothesesEvaluated: hyp}
		}
	}
	best.hypothesesEvaluated = hyp
	return best
}

const (
	experimentalV4Build29MinProjectiveValidation = 0.35
	experimentalV4Build29MinProjectiveMargin     = 0.15
	experimentalV4Build29MinPaddedValidation     = 0.55
	experimentalV4Build29MinPaddedMargin         = 0.20
)

func experimentalV4Build29ProjectiveAccepted(r experimentalV4JointProjectiveResult) bool {
	return r.available && r.detection.Available && r.placement.validationScore >= experimentalV4Build29MinProjectiveValidation && r.detection.Margin >= experimentalV4Build29MinProjectiveMargin && r.detection.OriginXBlocks == 0 && r.detection.OriginYBlocks == 0
}

func experimentalV4Build29PaddedAccepted(r experimentalV4JointPaddedResult) bool {
	return r.available && r.detection.Available && r.placement.validationScore >= experimentalV4Build29MinPaddedValidation && r.detection.Margin >= experimentalV4Build29MinPaddedMargin && r.detection.OriginXBlocks == 0 && r.detection.OriginYBlocks == 0
}

// Build34 helpers. These keep geometry proposal key-independent and data-plane
// agnostic: only the public pilot and image structure participate before the
// final authenticated frame decode.
func experimentalV4Build34CanonicalizeOrigin(h homography, originX, originY int) homography {
	if originX == 0 && originY == 0 {
		return h
	}
	// Detection origin (ox,oy) means that pilot residue (px,py) is currently
	// observed when sampling canonical residue (px-ox,py-oy). Move the canonical
	// coordinate system by that many blocks so the same physical samples become
	// origin (0,0). This is a right-side/domain translation, unlike placement.
	domain := homography{h: [9]float64{1, 0, -float64(originX * blockSize), 0, 1, -float64(originY * blockSize), 0, 0, 1}}
	return experimentalV4BlindMultiplyHomography(h, domain)
}

func experimentalV4Build34InteriorRowsScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography, originX, originY int) (float64, int) {
	bw, bh := canonicalWidth/blockSize, canonicalHeight/blockSize
	tilesX := bw / experimentalV4TileWidthBlocks
	tilesY := bh / experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 3 {
		return math.Inf(-1), 0
	}
	num, den := 0.0, 0.0
	visible := 0
	for ty := 1; ty < tilesY-1; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			baseX := tx * experimentalV4TileWidthBlocks
			baseY := ty * experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				rx := positiveMod(px-originX, experimentalV4TileWidthBlocks)
				ry := positiveMod(py-originY, experimentalV4TileHeightBlocks)
				v, ok := readProjectiveBlockValue(plane, h, (baseX+rx)*blockSize, (baseY+ry)*blockSize, blockSize)
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

func experimentalV4Build34CenteredProposal(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight, obsW, obsH int, p experimentalV4BlindGeometryParams) experimentalV4JointProjectiveCandidate {
	base, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
	cx, cy := experimentalV4JointProjectiveCenteredShift(fw, fh, obsW, obsH)
	phases := [][2]float64{{0, 0}, {-2, 0}, {2, 0}, {0, -2}, {0, 2}, {-4, 0}, {4, 0}, {0, -4}, {0, 4}}
	best := experimentalV4JointProjectiveCandidate{params: p, proposal: math.Inf(-1), fullW: fw, fullH: fh}
	bw, bh := canonicalWidth/blockSize, canonicalHeight/blockSize
	startX := ((bw - experimentalV4TileWidthBlocks) / 2 / experimentalV4TileWidthBlocks) * experimentalV4TileWidthBlocks
	startY := ((bh - experimentalV4TileHeightBlocks) / 2 / experimentalV4TileHeightBlocks) * experimentalV4TileHeightBlocks
	for _, ph := range phases {
		h := experimentalV4PlacementShiftHomography(base, cx+ph[0], cy+ph[1])
		d := experimentalV4DetectPilotSingleTilePlane(plane, candidate, startX, startY, h)
		if !d.Available || d.VisiblePilotPositions < 48 {
			continue
		}
		interior, vis := experimentalV4Build34InteriorRowsScore(plane, candidate, canonicalWidth, canonicalHeight, h, d.OriginXBlocks, d.OriginYBlocks)
		if vis < 48 || math.IsInf(interior, -1) {
			continue
		}
		objective := math.Min(d.Score, interior)
		if objective > best.proposal {
			best.proposal = objective
			best.h = h
			best.originX, best.originY = d.OriginXBlocks, d.OriginYBlocks
		}
	}
	return best
}

func experimentalV4Build34LocalNeighborhood(anchor experimentalV4JointProjectiveCandidate) []experimentalV4JointProjectiveCandidate {
	dA := []float64{-0.3, 0, 0.3}
	dSX := []float64{-0.03, 0, 0.03}
	dSY := []float64{-0.01, -0.005, 0, 0.005, 0.01}
	dI := []float64{-0.015, 0, 0.015}
	out := make([]experimentalV4JointProjectiveCandidate, 0, 243)
	seen := map[[5]int]bool{}
	for _, da := range dA {
		for _, dsx := range dSX {
			for _, dsy := range dSY {
				for _, dt := range dI {
					for _, db := range dI {
						p := anchor.params
						p.angleDeg += da
						p.scaleX += dsx
						p.scaleY += dsy
						p.topInset += dt
						p.bottomInset += db
						if p.scaleX < .90 || p.scaleX > 1.18 || p.scaleY < .82 || p.scaleY > 1.08 || p.topInset < 0 || p.topInset > .075 || p.bottomInset < 0 || p.bottomInset > .075 {
							continue
						}
						key := [5]int{int(math.Round(p.angleDeg * 10)), int(math.Round(p.scaleX * 1000)), int(math.Round(p.scaleY * 1000)), int(math.Round(p.topInset * 1000)), int(math.Round(p.bottomInset * 1000))}
						if seen[key] {
							continue
						}
						seen[key] = true
						out = append(out, experimentalV4JointProjectiveCandidate{params: p})
					}
				}
			}
		}
	}
	return out
}
