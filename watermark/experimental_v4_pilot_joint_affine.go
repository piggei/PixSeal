package watermark

import (
	"image"
	"math"
	"runtime"
	"sort"
	"sync"
)

// Build28 joins the two independently qualified Build26/Build27 problems for
// the affine crop case: rotation/anisotropic scale and negative crop placement
// are both unknown. The search remains experimental and public-signal only.
//
// Evidence is deliberately staged:
//   1. absolute DCT-difference phase contrast proposes affine geometry without
//      knowing pilot signs, payload, key, ECC or HMAC;
//   2. a compact coupled DCT refinement resolves angle/scale coupling on the
//      final structural basins, still without consulting pilot symbols;
//   3. only after geometry is fixed does the Build27 placement search expose
//      pilot partition 0 for placement proposal and partition 1 for held-out
//      placement validation;
//   4. the full public pilot supplies absolute-origin diagnostics only.
//
// Build28 qualifies negative crop/translation only. Padded-canvas joint search
// and joint projective search remain separate future work.

type experimentalV4JointAffineResult struct {
	available           bool
	params              experimentalV4BlindGeometryParams
	placement           experimentalV4PlacementResult
	detection           ExperimentalV4PilotDetection
	structuralScore     float64
	proposalScore       float64
	hypothesesEvaluated int
}

type experimentalV4JointAffineScored struct {
	params experimentalV4BlindGeometryParams
	score  float64
}

func experimentalV4JointAffineWorkers() int {
	n := runtime.GOMAXPROCS(0)
	if n > 8 {
		n = 8
	}
	if n < 1 {
		n = 1
	}
	return n
}

// experimentalV4JointAffineBlockEnergy measures only |DCT(2,3)-DCT(3,2)|.
// Therefore it is independent of embedded symbol polarity, pilot signs,
// payload bytes and key material. The middle 60% trimmed mean reduces the
// leverage of natural high-contrast edges.
func experimentalV4JointAffineBlockEnergy(plane *pixelPlane, canonicalWidth, canonicalHeight int, h homography, samples int) (float64, int) {
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	total := blocksWide * blocksHigh
	if total <= 0 || samples <= 0 {
		return math.Inf(-1), 0
	}
	values := make([]float64, 0, samples)
	for k := 0; k < total && len(values) < samples; k++ {
		idx := (k*17 + 11) % total
		x := idx % blocksWide
		y := idx / blocksWide
		v, ok := readProjectiveBlockValue(plane, h, x*blockSize, y*blockSize, blockSize)
		if ok {
			values = append(values, math.Abs(v))
		}
	}
	if len(values) < 16 {
		return math.Inf(-1), len(values)
	}
	sort.Float64s(values)
	lo := len(values) / 5
	hi := len(values) * 4 / 5
	if hi <= lo {
		return math.Inf(-1), len(values)
	}
	sum := 0.0
	for _, v := range values[lo:hi] {
		sum += v
	}
	return sum / float64(hi-lo), len(values)
}

// Phase contrast asks whether one sub-block phase is substantially more
// DCT-carrier-like than the typical phase. Natural directional texture can have
// high absolute DCT energy; it is much less likely to produce a sharp 8-pixel
// phase maximum over the whole image.
func experimentalV4JointAffinePhaseContrast(plane *pixelPlane, canonicalWidth, canonicalHeight int, p experimentalV4BlindGeometryParams, samples int) float64 {
	base, _, _ := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, p)
	offsets := [...]int{-4, -2, 0, 2}
	values := make([]float64, 0, len(offsets)*len(offsets))
	best := math.Inf(-1)
	for _, dy := range offsets {
		for _, dx := range offsets {
			h := experimentalV4PlacementShiftHomography(base, float64(dx), float64(dy))
			score, _ := experimentalV4JointAffineBlockEnergy(plane, canonicalWidth, canonicalHeight, h, samples)
			values = append(values, score)
			if score > best {
				best = score
			}
		}
	}
	if len(values) != 16 || math.IsInf(best, -1) {
		return math.Inf(-1)
	}
	sort.Float64s(values)
	median := 0.5 * (values[7] + values[8])
	return best - median
}

func experimentalV4JointAffineRefineMid(plane *pixelPlane, canonicalWidth, canonicalHeight int, seed experimentalV4JointAffineScored) experimentalV4JointAffineScored {
	best := experimentalV4JointAffineScored{params: seed.params, score: math.Inf(-1)}
	eval := func(p experimentalV4BlindGeometryParams) {
		s := experimentalV4JointAffinePhaseContrast(plane, canonicalWidth, canonicalHeight, p, 64)
		if s > best.score {
			best = experimentalV4JointAffineScored{params: p, score: s}
		}
	}
	eval(seed.params)
	base := best.params
	for a := base.angleDeg - 0.5; a <= base.angleDeg+0.5001; a += 0.1 {
		q := best.params
		q.angleDeg = a
		eval(q)
	}
	base = best.params
	for sx := math.Max(0.85, base.scaleX-0.035); sx <= math.Min(1.15, base.scaleX+0.035)+1e-9; sx += 0.005 {
		q := best.params
		q.scaleX = sx
		eval(q)
	}
	base = best.params
	for sy := math.Max(0.85, base.scaleY-0.035); sy <= math.Min(1.15, base.scaleY+0.035)+1e-9; sy += 0.005 {
		q := best.params
		q.scaleY = sy
		eval(q)
	}
	base = best.params
	for a := base.angleDeg - 0.15; a <= base.angleDeg+0.1501; a += 0.05 {
		q := best.params
		q.angleDeg = a
		eval(q)
	}
	return best
}

func experimentalV4JointAffineRefineFull(plane *pixelPlane, canonicalWidth, canonicalHeight int, seed experimentalV4JointAffineScored) experimentalV4JointAffineScored {
	best := experimentalV4JointAffineScored{params: seed.params, score: math.Inf(-1)}
	eval := func(p experimentalV4BlindGeometryParams) {
		s := experimentalV4JointAffinePhaseContrast(plane, canonicalWidth, canonicalHeight, p, 256)
		if s > best.score {
			best = experimentalV4JointAffineScored{params: p, score: s}
		}
	}
	eval(seed.params)
	for pass := 0; pass < 2; pass++ {
		base := best.params
		for a := base.angleDeg - 0.35; a <= base.angleDeg+0.3501; a += 0.05 {
			q := best.params
			q.angleDeg = a
			eval(q)
		}
		base = best.params
		for sx := math.Max(0.85, base.scaleX-0.015); sx <= math.Min(1.15, base.scaleX+0.015)+1e-9; sx += 0.0025 {
			q := best.params
			q.scaleX = sx
			eval(q)
		}
		base = best.params
		for sy := math.Max(0.85, base.scaleY-0.015); sy <= math.Min(1.15, base.scaleY+0.015)+1e-9; sy += 0.0025 {
			q := best.params
			q.scaleY = sy
			eval(q)
		}
	}
	return best
}

// experimentalV4JointAffineRefineCoupled is intentionally reserved for the
// final two DCT basins.  The preceding coordinate refinements are cheap and
// effective at getting close to the correct affine family, but under crop a
// small angle/scale coupling can leave one axis in a nearby local maximum.
// Evaluating a compact 3-D neighborhood with the full 256-block structural
// support removes that last dependency without making the global bank cubic.
func experimentalV4JointAffineRefineCoupled(plane *pixelPlane, canonicalWidth, canonicalHeight int, seed experimentalV4JointAffineScored) experimentalV4JointAffineScored {
	best := experimentalV4JointAffineScored{params: seed.params, score: math.Inf(-1)}
	for a := seed.params.angleDeg - 0.25; a <= seed.params.angleDeg+0.2501; a += 0.05 {
		for sx := math.Max(0.85, seed.params.scaleX-0.0125); sx <= math.Min(1.15, seed.params.scaleX+0.0125)+1e-9; sx += 0.0025 {
			for sy := math.Max(0.85, seed.params.scaleY-0.0200); sy <= math.Min(1.15, seed.params.scaleY+0.0200)+1e-9; sy += 0.0025 {
				p := seed.params
				p.angleDeg = a
				p.scaleX = sx
				p.scaleY = sy
				s := experimentalV4JointAffinePhaseContrast(plane, canonicalWidth, canonicalHeight, p, 256)
				if s > best.score {
					best = experimentalV4JointAffineScored{params: p, score: s}
				}
			}
		}
	}
	return best
}

func experimentalV4JointAffineGeometryFinalists(plane *pixelPlane, canonicalWidth, canonicalHeight int) ([]experimentalV4JointAffineScored, int) {
	angles := make([]float64, 0, 97)
	for a := -12.0; a <= 12.0001; a += 0.25 {
		angles = append(angles, a)
	}
	seeds := make([]experimentalV4JointAffineScored, len(angles))
	workers := experimentalV4JointAffineWorkers()
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for i, angle := range angles {
		i, angle := i, angle
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			best := experimentalV4JointAffineScored{score: math.Inf(-1)}
			for sx := 0.88; sx <= 1.1201; sx += 0.01 {
				for sy := 0.88; sy <= 1.1201; sy += 0.01 {
					p := experimentalV4BlindGeometryParams{angleDeg: angle, scaleX: sx, scaleY: sy}
					score := experimentalV4JointAffinePhaseContrast(plane, canonicalWidth, canonicalHeight, p, 32)
					if score > best.score {
						best = experimentalV4JointAffineScored{params: p, score: score}
					}
				}
			}
			seeds[i] = best
		}()
	}
	wg.Wait()
	coarseHypotheses := len(angles) * 25 * 25
	sort.Slice(seeds, func(i, j int) bool {
		if seeds[i].score == seeds[j].score {
			return seeds[i].params.angleDeg < seeds[j].params.angleDeg
		}
		return seeds[i].score > seeds[j].score
	})
	if len(seeds) > 15 {
		seeds = seeds[:15]
	}

	mid := make([]experimentalV4JointAffineScored, len(seeds))
	for i, seed := range seeds {
		i, seed := i, seed
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			mid[i] = experimentalV4JointAffineRefineMid(plane, canonicalWidth, canonicalHeight, seed)
		}()
	}
	wg.Wait()

	full := make([]experimentalV4JointAffineScored, len(mid))
	for i, seed := range mid {
		i, seed := i, seed
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			full[i] = experimentalV4JointAffineRefineFull(plane, canonicalWidth, canonicalHeight, seed)
		}()
	}
	wg.Wait()
	sort.Slice(full, func(i, j int) bool {
		if full[i].score == full[j].score {
			if full[i].params.angleDeg == full[j].params.angleDeg {
				if full[i].params.scaleX == full[j].params.scaleX {
					return full[i].params.scaleY < full[j].params.scaleY
				}
				return full[i].params.scaleX < full[j].params.scaleX
			}
			return full[i].params.angleDeg < full[j].params.angleDeg
		}
		return full[i].score > full[j].score
	})
	if len(full) > 2 {
		full = full[:2]
	}

	// Only the final two basins pay for the coupled refinement.  Running this
	// globally would turn the bounded hierarchical search back into a cubic
	// brute force.  The two refinements are independent and deterministic.
	coupled := make([]experimentalV4JointAffineScored, len(full))
	for i, seed := range full {
		i, seed := i, seed
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			coupled[i] = experimentalV4JointAffineRefineCoupled(plane, canonicalWidth, canonicalHeight, seed)
		}()
	}
	wg.Wait()
	sort.Slice(coupled, func(i, j int) bool {
		if coupled[i].score == coupled[j].score {
			if coupled[i].params.angleDeg == coupled[j].params.angleDeg {
				if coupled[i].params.scaleX == coupled[j].params.scaleX {
					return coupled[i].params.scaleY < coupled[j].params.scaleY
				}
				return coupled[i].params.scaleX < coupled[j].params.scaleX
			}
			return coupled[i].params.angleDeg < coupled[j].params.angleDeg
		}
		return coupled[i].score > coupled[j].score
	})
	full = coupled
	// Count geometry parameter hypotheses rather than individual block/phase reads.
	// The fixed refinement budgets are deterministic and small compared with the
	// coarse 30,625-point bank.
	const midBudgetPerSeed = 11 + 15 + 15 + 7 + 1
	const fullBudgetPerSeed = 2*(15+13+13) + 1
	const coupledBudgetPerFinalist = 11 * 11 * 17
	return full, coarseHypotheses + len(seeds)*midBudgetPerSeed + len(mid)*fullBudgetPerSeed + len(full)*coupledBudgetPerFinalist
}

func experimentalV4JointAffineSearch(img image.Image, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int) experimentalV4JointAffineResult {
	if img == nil || canonicalWidth < experimentalV4TileWidthBlocks*blockSize || canonicalHeight < experimentalV4TileHeightBlocks*blockSize {
		return experimentalV4JointAffineResult{}
	}
	plane := newPixelPlane(img)
	finalists, hypotheses := experimentalV4JointAffineGeometryFinalists(plane, canonicalWidth, canonicalHeight)
	if len(finalists) == 0 {
		return experimentalV4JointAffineResult{hypothesesEvaluated: hypotheses}
	}
	// Geometry selection is deliberately structural-only.  Earlier prototypes
	// used pilot partition 0 to choose between the last two DCT basins; the
	// coupled 256-block refinement makes that unnecessary and gives a cleaner
	// evidence boundary: no pilot symbol participates in geometry selection.
	win := finalists[0]
	geometry, fw, fh := experimentalV4BlindHomography(canonicalWidth, canonicalHeight, win.params)
	placement := experimentalV4PilotPlacementSearch(img, candidate, canonicalWidth, canonicalHeight, geometry, fw, fh)
	hypotheses += placement.hypothesesEvaluated
	if !placement.available {
		return experimentalV4JointAffineResult{hypothesesEvaluated: hypotheses}
	}
	detection := experimentalV4DetectPilotProjective(img, candidate, canonicalWidth, canonicalHeight, placement.h)
	return experimentalV4JointAffineResult{
		available:           true,
		params:              win.params,
		placement:           placement,
		detection:           detection,
		structuralScore:     win.score,
		proposalScore:       placement.proposalScore,
		hypothesesEvaluated: hypotheses,
	}
}
