package watermark

import (
	"errors"
	"image"
	"math"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Build80 is an equivalence-preserving performance candidate over the qualified
// Build76 smartphone baseline. It preserves the complete Build76 freeze, prefix1,
// generation-two, generation-three and generation-four scheduling/commit order.
// Only FoldScore evaluations inside continuation4 use a bounded per-FoldScore
// scratch cache of exact integer-source luminance values. mapPoint coordinates,
// RGB->luma arithmetic, bilinear interpolation and DCT accumulation stay in the
// original order; oversized local source rectangles fall back to the unchanged
// readProjectiveBlockValue path. Build76 remains qualified until physical evidence
// confirms both exact semantics and a material performance benefit.

type experimentalV4PhoneBuild80Prefix1Result struct {
	inputs  []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild80Gen2Result struct {
	outputs []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild80GeometryTelemetry struct {
	PlanePrepElapsed     time.Duration
	FreezeElapsed        time.Duration
	Prefix1WallElapsed   time.Duration
	Prefix1WorkerElapsed time.Duration
	Gen2WallElapsed      time.Duration
	Gen2WorkerElapsed    time.Duration
	Gen3WallElapsed      time.Duration
	Gen3WorkerElapsed    time.Duration
	Gen4WallElapsed      time.Duration
	Gen4WorkerElapsed    time.Duration
	Prefix1Workers       int
	Gen2Workers          int
	Gen3Workers          int
	Gen4Workers          int
	Gen2Tasks            int
	Gen3Tasks            int
	Gen4Tasks            int
	Gen2MinElapsed       time.Duration
	Gen2MedianElapsed    time.Duration
	Gen2MaxElapsed       time.Duration
	Gen2MinEvaluations   int
	Gen2MaxEvaluations   int
	Gen2MinOutputs       int
	Gen2MaxOutputs       int
	Gen3MinElapsed       time.Duration
	Gen3MedianElapsed    time.Duration
	Gen3MaxElapsed       time.Duration
	Gen3MinEvaluations   int
	Gen3MaxEvaluations   int
	Gen3MinOutputs       int
	Gen3MaxOutputs       int
	Gen4MinElapsed       time.Duration
	Gen4MedianElapsed    time.Duration
	Gen4MaxElapsed       time.Duration
	Gen4MinEvaluations   int
	Gen4MaxEvaluations   int
	Gen4MinBank          int
	Gen4MaxBank          int
	CacheBlockReads      int
	CacheHits            int
	CacheFallbacks       int
	CacheFailed          int
	CachePixelsPrepared  int64
	CacheMaxArea         int
}

const experimentalV4PhoneBuild80MaxCachePixels = 196

type experimentalV4PhoneBuild80BlockScratch struct {
	x0   [blockSize * blockSize]int
	y0   [blockSize * blockSize]int
	x1   [blockSize * blockSize]int
	y1   [blockSize * blockSize]int
	fx   [blockSize * blockSize]float64
	fy   [blockSize * blockSize]float64
	luma [experimentalV4PhoneBuild80MaxCachePixels]float64
}

type experimentalV4PhoneBuild80CacheTelemetry struct {
	BlockReads     int
	Hits           int
	Fallbacks      int
	Failed         int
	PixelsPrepared int64
	MaxArea        int
}

func experimentalV4PhoneBuild80MergeCache(dst *experimentalV4PhoneBuild80CacheTelemetry, src experimentalV4PhoneBuild80CacheTelemetry) {
	if dst == nil {
		return
	}
	dst.BlockReads += src.BlockReads
	dst.Hits += src.Hits
	dst.Fallbacks += src.Fallbacks
	dst.Failed += src.Failed
	dst.PixelsPrepared += src.PixelsPrepared
	if src.MaxArea > dst.MaxArea {
		dst.MaxArea = src.MaxArea
	}
}

// experimentalV4PhoneBuild80ReadProjectiveBlockCached preserves the exact
// Build41 block arithmetic while caching integer-source luminance values only
// within one projective 8x8 block. The same mapPoint coordinates, RGB->luma
// expression, bilinear interpolation order and DCT accumulation order are used.
// If the local source rectangle exceeds the bounded scratch cache, the original
// readProjectiveBlockValue path is used unchanged.
func experimentalV4PhoneBuild80ReadProjectiveBlockCached(src *pixelPlane, h homography, originX, originY, size int, scratch *experimentalV4PhoneBuild80BlockScratch) (float64, bool, bool, int) {
	if src == nil || scratch == nil || size != blockSize {
		v, ok := readProjectiveBlockValue(src, h, originX, originY, size)
		return v, ok, false, 0
	}
	width, height := src.bounds.Dx(), src.bounds.Dy()
	minX, minY := width, height
	maxX, maxY := -1, -1
	idx := 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx, sy, ok := h.mapPoint(float64(originX+x), float64(originY+y))
			if !ok || sx < 0 || sy < 0 || sx > float64(width-1) || sy > float64(height-1) {
				return 0, false, false, 0
			}
			x0, y0 := int(math.Floor(sx)), int(math.Floor(sy))
			x1, y1 := x0+1, y0+1
			if x1 >= width {
				x1 = width - 1
			}
			if y1 >= height {
				y1 = height - 1
			}
			scratch.x0[idx], scratch.y0[idx] = x0, y0
			scratch.x1[idx], scratch.y1[idx] = x1, y1
			scratch.fx[idx], scratch.fy[idx] = sx-float64(x0), sy-float64(y0)
			if x0 < minX {
				minX = x0
			}
			if x1 < minX {
				minX = x1
			}
			if y0 < minY {
				minY = y0
			}
			if y1 < minY {
				minY = y1
			}
			if x0 > maxX {
				maxX = x0
			}
			if x1 > maxX {
				maxX = x1
			}
			if y0 > maxY {
				maxY = y0
			}
			if y1 > maxY {
				maxY = y1
			}
			idx++
		}
	}
	cacheW, cacheH := maxX-minX+1, maxY-minY+1
	area := cacheW * cacheH
	if area <= 0 || area > experimentalV4PhoneBuild80MaxCachePixels {
		v, ok := readProjectiveBlockValue(src, h, originX, originY, size)
		return v, ok, false, area
	}
	for py := minY; py <= maxY; py++ {
		for px := minX; px <= maxX; px++ {
			srcIndex := (py*width + px) * 3
			cacheIndex := (py-minY)*cacheW + (px - minX)
			scratch.luma[cacheIndex] = .299*float64(src.rgb[srcIndex]) + .587*float64(src.rgb[srcIndex+1]) + .114*float64(src.rgb[srcIndex+2]) - 128
		}
	}
	table := readCosTables[size]
	c23, c32 := 0.0, 0.0
	idx = 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			x0, y0, x1, y1 := scratch.x0[idx], scratch.y0[idx], scratch.x1[idx], scratch.y1[idx]
			fx, fy := scratch.fx[idx], scratch.fy[idx]
			l00 := scratch.luma[(y0-minY)*cacheW+(x0-minX)]
			l10 := scratch.luma[(y0-minY)*cacheW+(x1-minX)]
			l01 := scratch.luma[(y1-minY)*cacheW+(x0-minX)]
			l11 := scratch.luma[(y1-minY)*cacheW+(x1-minX)]
			top := l00*(1-fx) + l10*fx
			bottom := l01*(1-fx) + l11*fx
			l := top*(1-fy) + bottom*fy
			c23 += l * table[3][x] * table[2][y]
			c32 += l * table[2][x] * table[3][y]
			idx++
		}
	}
	return math.Abs(c23) - math.Abs(c32), true, true, area
}

func experimentalV4PhoneBuild80FoldScoreCached(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, heldout int, proposal bool, sparse bool) (float64, int, experimentalV4PhoneBuild80CacheTelemetry) {
	var cache experimentalV4PhoneBuild80CacheTelemetry
	if plane == nil {
		return math.Inf(-1), 0, cache
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		return math.Inf(-1), 0, cache
	}
	num, den := 0.0, 0.0
	visible, usedTiles := 0, 0
	var scratch experimentalV4PhoneBuild80BlockScratch
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
				cache.BlockReads++
				v, ok, cached, area := experimentalV4PhoneBuild80ReadProjectiveBlockCached(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize, &scratch)
				if area > cache.MaxArea {
					cache.MaxArea = area
				}
				if cached {
					cache.Hits++
					cache.PixelsPrepared += int64(area)
				} else if ok {
					cache.Fallbacks++
				}
				if !ok {
					cache.Failed++
					continue
				}
				visible++
				num += float64(candidate.signs[i]) * v
				den += math.Abs(v)
			}
		}
	}
	if den <= 0 {
		return math.Inf(-1), visible, cache
	}
	return num / den, visible, cache
}

func experimentalV4PhoneBuild80ContinueCached(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild55BlindContinuationState, int, experimentalV4PhoneBuild80CacheTelemetry) {
	q, h, score := start.quad, start.h, start.proposal
	states := make([]experimentalV4PhoneBuild55BlindContinuationState, 0, experimentalV4PhoneBuild55MaxStatesBranch)
	evals := 0
	var cache experimentalV4PhoneBuild80CacheTelemetry
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
				ss, _, cp := experimentalV4PhoneBuild80FoldScoreCached(plane, candidate, cw, ch, hh, heldout, true, false)
				experimentalV4PhoneBuild80MergeCache(&cache, cp)
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
	return states, evals, cache
}

type experimentalV4PhoneBuild80Gen4Result struct {
	bank    []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
	cache   experimentalV4PhoneBuild80CacheTelemetry
}

func experimentalV4PhoneBuild80Generation4Cached(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild80CacheTelemetry) {
	bank := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	var cache experimentalV4PhoneBuild80CacheTelemetry
	se4, si4 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, input, 0)
	evals += se4
	if si4 != 0 {
		return bank, evals, cache
	}
	pairs4, pe4, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, input, 0)
	evals += pe4
	for _, pair4 := range pairs4 {
		cont4, ce4, cp := experimentalV4PhoneBuild80ContinueCached(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
		experimentalV4PhoneBuild80MergeCache(&cache, cp)
		evals += ce4
		for _, c4 := range cont4 {
			bank = append(bank, c4.hyp)
		}
	}
	return bank, evals, cache
}

// experimentalV4PhoneBuild80SeedPrefix1 reproduces Build64 exactly through the
// first sibling stencil and freezes sibling1 hypotheses in traversal order.
func experimentalV4PhoneBuild80SeedPrefix1(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed experimentalV4PhoneBuild48Seed) ([]experimentalV4PhoneHypothesis, int) {
	inputs := make([]experimentalV4PhoneHypothesis, 0, 32)
	evals := 0
	baseline, n := experimentalV4PhoneBuild41Refine(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
	evals += n
	roots, n := experimentalV4PhoneBuild53TwoPxRoots(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
	evals += n
	if baseline.h.h[8] == 0 || len(roots) == 0 {
		return inputs, evals
	}
	for _, root := range roots {
		se, si := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
		evals += se
		if si != 0 {
			continue
		}
		pairs, pe, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
		evals += pe
		for _, pair := range pairs {
			cont, ce := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
			evals += ce
			for _, cs := range cont {
				sibs, ne := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, cs.hyp, 0)
				evals += ne
				for _, sib := range sibs {
					inputs = append(inputs, sib.hyp)
				}
			}
		}
	}
	return inputs, evals
}

// experimentalV4PhoneBuild80Generation2 executes exactly the Build64 second
// single/pair/continuation/sibling subtree for one frozen sibling1 input.
func experimentalV4PhoneBuild80Generation2(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int) {
	outputs := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	se2, si2 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, input, 0)
	evals += se2
	if si2 != 0 {
		return outputs, evals
	}
	pairs2, pe2, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, input, 0)
	evals += pe2
	for _, pair2 := range pairs2 {
		cont2, ce2 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair2.hyp, 0)
		evals += ce2
		for _, c2 := range cont2 {
			sibs2, ne2 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c2.hyp, 0)
			evals += ne2
			for _, s2 := range sibs2 {
				outputs = append(outputs, s2.hyp)
			}
		}
	}
	return outputs, evals
}

func experimentalV4PhoneBuild80BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild80GeometryTelemetry, experimentalV4PhoneBuild75FreezeTelemetry) {
	var profile experimentalV4PhoneBuild80GeometryTelemetry
	var freezeProfile experimentalV4PhoneBuild75FreezeTelemetry
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0, 0, profile, freezeProfile
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)
	started := time.Now()
	plane := newPixelPlane(work)
	profile.PlanePrepElapsed = time.Since(started)
	pilot := experimentalV4Prototype2Candidate()

	frozen, _, evals, freezeProfile := experimentalV4PhoneBuild75FreezeParallel(work, boundary, cw, ch, 128)
	profile.FreezeElapsed = freezeProfile.TotalElapsed
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild63SeedsPerPair)
	if len(seeds) == 0 {
		return nil, evals, 0, 0, profile, freezeProfile
	}

	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(seeds) {
		workers = len(seeds)
	}
	profile.Prefix1Workers = workers
	prefixes := make([]experimentalV4PhoneBuild80Prefix1Result, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				s := time.Now()
				inputs, n := experimentalV4PhoneBuild80SeedPrefix1(plane, pilot, cw, ch, anchor, seeds[i])
				prefixes[i] = experimentalV4PhoneBuild80Prefix1Result{inputs: inputs, evals: n, elapsed: time.Since(s)}
			}
		}()
	}
	for i := range seeds {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	profile.Prefix1WallElapsed = time.Since(started)

	totalGen2 := 0
	for i := range prefixes {
		evals += prefixes[i].evals
		profile.Prefix1WorkerElapsed += prefixes[i].elapsed
		totalGen2 += len(prefixes[i].inputs)
	}
	profile.Gen2Tasks = totalGen2
	if totalGen2 == 0 {
		return nil, evals, len(seeds), workers, profile, freezeProfile
	}
	flat2 := make([]experimentalV4PhoneHypothesis, 0, totalGen2)
	for i := range prefixes {
		flat2 = append(flat2, prefixes[i].inputs...)
	}
	gen2Workers := runtime.GOMAXPROCS(0)
	if gen2Workers < 1 {
		gen2Workers = 1
	}
	if gen2Workers > len(flat2) {
		gen2Workers = len(flat2)
	}
	profile.Gen2Workers = gen2Workers
	gen2Results := make([]experimentalV4PhoneBuild80Gen2Result, len(flat2))
	gen2Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen2Workers)
	started = time.Now()
	for w := 0; w < gen2Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen2Jobs {
				s := time.Now()
				outputs, n := experimentalV4PhoneBuild80Generation2(plane, pilot, cw, ch, anchor, flat2[i])
				gen2Results[i] = experimentalV4PhoneBuild80Gen2Result{outputs: outputs, evals: n, elapsed: time.Since(s)}
			}
		}()
	}
	for i := range flat2 {
		gen2Jobs <- i
	}
	close(gen2Jobs)
	wg.Wait()
	profile.Gen2WallElapsed = time.Since(started)
	gen2Durations := make([]time.Duration, 0, len(gen2Results))
	totalGen3 := 0
	for i, r := range gen2Results {
		evals += r.evals
		profile.Gen2WorkerElapsed += r.elapsed
		totalGen3 += len(r.outputs)
		gen2Durations = append(gen2Durations, r.elapsed)
		if i == 0 || r.elapsed < profile.Gen2MinElapsed {
			profile.Gen2MinElapsed = r.elapsed
		}
		if r.elapsed > profile.Gen2MaxElapsed {
			profile.Gen2MaxElapsed = r.elapsed
		}
		if i == 0 || r.evals < profile.Gen2MinEvaluations {
			profile.Gen2MinEvaluations = r.evals
		}
		if r.evals > profile.Gen2MaxEvaluations {
			profile.Gen2MaxEvaluations = r.evals
		}
		if i == 0 || len(r.outputs) < profile.Gen2MinOutputs {
			profile.Gen2MinOutputs = len(r.outputs)
		}
		if len(r.outputs) > profile.Gen2MaxOutputs {
			profile.Gen2MaxOutputs = len(r.outputs)
		}
	}
	sort.Slice(gen2Durations, func(i, j int) bool { return gen2Durations[i] < gen2Durations[j] })
	if len(gen2Durations) > 0 {
		profile.Gen2MedianElapsed = gen2Durations[len(gen2Durations)/2]
	}
	profile.Gen3Tasks = totalGen3
	if totalGen3 == 0 {
		return nil, evals, len(seeds), workers, profile, freezeProfile
	}

	flat3 := make([]experimentalV4PhoneHypothesis, 0, totalGen3)
	for i := range gen2Results {
		flat3 = append(flat3, gen2Results[i].outputs...)
	}
	gen3Workers := runtime.GOMAXPROCS(0)
	if gen3Workers < 1 {
		gen3Workers = 1
	}
	if gen3Workers > len(flat3) {
		gen3Workers = len(flat3)
	}
	profile.Gen3Workers = gen3Workers
	gen3Results := make([]experimentalV4PhoneBuild73Gen3Result, len(flat3))
	gen3Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen3Workers)
	started = time.Now()
	for w := 0; w < gen3Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen3Jobs {
				s := time.Now()
				outputs, n := experimentalV4PhoneBuild73Generation3(plane, pilot, cw, ch, anchor, flat3[i])
				gen3Results[i] = experimentalV4PhoneBuild73Gen3Result{outputs: outputs, evals: n, elapsed: time.Since(s)}
			}
		}()
	}
	for i := range flat3 {
		gen3Jobs <- i
	}
	close(gen3Jobs)
	wg.Wait()
	profile.Gen3WallElapsed = time.Since(started)
	gen3Durations := make([]time.Duration, 0, len(gen3Results))
	totalGen4 := 0
	for i, r := range gen3Results {
		evals += r.evals
		profile.Gen3WorkerElapsed += r.elapsed
		totalGen4 += len(r.outputs)
		gen3Durations = append(gen3Durations, r.elapsed)
		if i == 0 || r.elapsed < profile.Gen3MinElapsed {
			profile.Gen3MinElapsed = r.elapsed
		}
		if r.elapsed > profile.Gen3MaxElapsed {
			profile.Gen3MaxElapsed = r.elapsed
		}
		if i == 0 || r.evals < profile.Gen3MinEvaluations {
			profile.Gen3MinEvaluations = r.evals
		}
		if r.evals > profile.Gen3MaxEvaluations {
			profile.Gen3MaxEvaluations = r.evals
		}
		if i == 0 || len(r.outputs) < profile.Gen3MinOutputs {
			profile.Gen3MinOutputs = len(r.outputs)
		}
		if len(r.outputs) > profile.Gen3MaxOutputs {
			profile.Gen3MaxOutputs = len(r.outputs)
		}
	}
	sort.Slice(gen3Durations, func(i, j int) bool { return gen3Durations[i] < gen3Durations[j] })
	if len(gen3Durations) > 0 {
		profile.Gen3MedianElapsed = gen3Durations[len(gen3Durations)/2]
	}
	profile.Gen4Tasks = totalGen4
	if totalGen4 == 0 {
		return nil, evals, len(seeds), workers, profile, freezeProfile
	}

	flat4 := make([]experimentalV4PhoneHypothesis, 0, totalGen4)
	for i := range gen3Results {
		flat4 = append(flat4, gen3Results[i].outputs...)
	}
	gen4Workers := runtime.GOMAXPROCS(0)
	if gen4Workers < 1 {
		gen4Workers = 1
	}
	if gen4Workers > len(flat4) {
		gen4Workers = len(flat4)
	}
	profile.Gen4Workers = gen4Workers
	gen4Results := make([]experimentalV4PhoneBuild80Gen4Result, len(flat4))
	gen4Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen4Workers)
	started = time.Now()
	for w := 0; w < gen4Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen4Jobs {
				s := time.Now()
				bank, n, cp := experimentalV4PhoneBuild80Generation4Cached(plane, pilot, cw, ch, anchor, flat4[i])
				gen4Results[i] = experimentalV4PhoneBuild80Gen4Result{bank: bank, evals: n, elapsed: time.Since(s), cache: cp}
			}
		}()
	}
	for i := range flat4 {
		gen4Jobs <- i
	}
	close(gen4Jobs)
	wg.Wait()
	profile.Gen4WallElapsed = time.Since(started)
	bank := make([]experimentalV4PhoneHypothesis, 0, 1024)
	gen4Durations := make([]time.Duration, 0, len(gen4Results))
	for i, r := range gen4Results {
		evals += r.evals
		bank = append(bank, r.bank...)
		profile.Gen4WorkerElapsed += r.elapsed
		profile.CacheBlockReads += r.cache.BlockReads
		profile.CacheHits += r.cache.Hits
		profile.CacheFallbacks += r.cache.Fallbacks
		profile.CacheFailed += r.cache.Failed
		profile.CachePixelsPrepared += r.cache.PixelsPrepared
		if r.cache.MaxArea > profile.CacheMaxArea {
			profile.CacheMaxArea = r.cache.MaxArea
		}
		gen4Durations = append(gen4Durations, r.elapsed)
		if i == 0 || r.elapsed < profile.Gen4MinElapsed {
			profile.Gen4MinElapsed = r.elapsed
		}
		if r.elapsed > profile.Gen4MaxElapsed {
			profile.Gen4MaxElapsed = r.elapsed
		}
		if i == 0 || r.evals < profile.Gen4MinEvaluations {
			profile.Gen4MinEvaluations = r.evals
		}
		if r.evals > profile.Gen4MaxEvaluations {
			profile.Gen4MaxEvaluations = r.evals
		}
		if i == 0 || len(r.bank) < profile.Gen4MinBank {
			profile.Gen4MinBank = len(r.bank)
		}
		if len(r.bank) > profile.Gen4MaxBank {
			profile.Gen4MaxBank = len(r.bank)
		}
	}
	sort.Slice(gen4Durations, func(i, j int) bool { return gen4Durations[i] < gen4Durations[j] })
	if len(gen4Durations) > 0 {
		profile.Gen4MedianElapsed = gen4Durations[len(gen4Durations)/2]
	}
	return bank, evals, len(seeds), workers, profile, freezeProfile
}

type experimentalV4PhoneBuild80RecoveryTelemetry struct {
	experimentalV4PhoneBuild75RecoveryTelemetry
	Build80GeometryProfile experimentalV4PhoneBuild80GeometryTelemetry
}

func experimentalV4PhoneBuild80Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild80RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild80RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }
	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, gp, fp := experimentalV4PhoneBuild80BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.Build80GeometryProfile = gp
	telemetry.FreezeProfile = fp
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = prefixWorkers
	telemetry.BankCandidates = len(bank)
	// Populate the qualified Build73-compatible geometry view so all historical
	// telemetry remains meaningful under the earlier barrier.
	compat := &telemetry.GeometryProfile
	compat.PlanePrepElapsed = gp.PlanePrepElapsed
	compat.FreezeElapsed = gp.FreezeElapsed
	compat.Prefix2Workers = gp.Prefix1Workers
	compat.Prefix2WallElapsed = gp.Prefix1WallElapsed + gp.Gen2WallElapsed
	compat.Prefix2WorkerElapsed = gp.Prefix1WorkerElapsed + gp.Gen2WorkerElapsed
	compat.Gen3Workers = gp.Gen3Workers
	compat.Gen4Workers = gp.Gen4Workers
	compat.Gen3Tasks = gp.Gen3Tasks
	compat.Gen4Tasks = gp.Gen4Tasks
	compat.Gen3WallElapsed = gp.Gen3WallElapsed
	compat.Gen3WorkerElapsed = gp.Gen3WorkerElapsed
	compat.Gen4WallElapsed = gp.Gen4WallElapsed
	compat.Gen4WorkerElapsed = gp.Gen4WorkerElapsed
	compat.Gen3MinElapsed = gp.Gen3MinElapsed
	compat.Gen3MedianElapsed = gp.Gen3MedianElapsed
	compat.Gen3MaxElapsed = gp.Gen3MaxElapsed
	compat.Gen3MinEvaluations = gp.Gen3MinEvaluations
	compat.Gen3MaxEvaluations = gp.Gen3MaxEvaluations
	compat.Gen3MinOutputs = gp.Gen3MinOutputs
	compat.Gen3MaxOutputs = gp.Gen3MaxOutputs
	compat.Gen4MinElapsed = gp.Gen4MinElapsed
	compat.Gen4MedianElapsed = gp.Gen4MedianElapsed
	compat.Gen4MaxElapsed = gp.Gen4MaxElapsed
	compat.Gen4MinEvaluations = gp.Gen4MinEvaluations
	compat.Gen4MaxEvaluations = gp.Gen4MaxEvaluations
	compat.Gen4MinBank = gp.Gen4MinBank
	compat.Gen4MaxBank = gp.Gen4MaxBank
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build80 recovery bank empty")
	}
	planeStarted := time.Now()
	plane := newPixelPlane(work)
	telemetry.PlanePrepElapsed = time.Since(planeStarted)
	pilot := experimentalV4Prototype2Candidate()
	qualified := make([]experimentalV4PhoneHypothesis, 0, len(bank))
	qualificationStarted := time.Now()
	for _, h := range bank {
		q, n, ok := experimentalV4PhoneBuild68Qualify(plane, pilot, cw, ch, h, 0)
		telemetry.QualificationEvaluations += n
		if ok {
			qualified = append(qualified, q)
		}
	}
	telemetry.QualificationElapsed = time.Since(qualificationStarted)
	telemetry.QualifiedCandidates = len(qualified)
	telemetry.DecodeWorkers = experimentalV4PhoneBuild66DecodeWorkerCount(len(qualified))
	if len(qualified) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build80 recovery authentication failed")
	}
	workers := telemetry.DecodeWorkers
	decodeStarted := time.Now()
	for start := 0; start < len(qualified); start += workers {
		end := start + workers
		if end > len(qualified) {
			end = len(qualified)
		}
		results := make([]experimentalV4PhoneBuild67DecodeResult, end-start)
		var wg sync.WaitGroup
		wg.Add(end - start)
		for i := start; i < end; i++ {
			i := i
			go func() {
				defer wg.Done()
				results[i-start] = experimentalV4PhoneBuild67DecodeSingle(plane, qualified[i], key, cw, ch)
			}()
		}
		wg.Wait()
		experimentalV4PhoneBuild68AccumulatePhysical(&telemetry.experimentalV4PhoneBuild68RecoveryTelemetry, results)
		semanticResults := make([]experimentalV4PhoneBuild66DecodeResult, len(results))
		for i := range results {
			semanticResults[i] = results[i].experimentalV4PhoneBuild66DecodeResult
		}
		if payload, info, ok := experimentalV4PhoneBuild66ConsumeBatch(semanticResults, &telemetry.experimentalV4PhoneBuild66RecoveryTelemetry); ok {
			telemetry.DecodeWallElapsed = time.Since(decodeStarted)
			finish()
			return payload, info, telemetry, nil
		}
	}
	telemetry.DecodeWallElapsed = time.Since(decodeStarted)
	finish()
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build80 recovery authentication failed")
}

func experimentalV4PhoneBuild80ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild80RecoveryTelemetry) {
	experimentalV4PhoneBuild75ApplyTelemetry(public, recovery.experimentalV4PhoneBuild75RecoveryTelemetry)
	p := recovery.Build80GeometryProfile
	// Preserve the qualified Build76 public view exactly so existing physical
	// gates and downstream tooling continue to see the same scheduler semantics.
	public.Build76Attempted = recovery.Attempted
	public.Build76Prefix1Workers = p.Prefix1Workers
	public.Build76Gen2Workers = p.Gen2Workers
	public.Build76Gen3Workers = p.Gen3Workers
	public.Build76Gen4Workers = p.Gen4Workers
	public.Build76Gen2Tasks = p.Gen2Tasks
	public.Build76Gen3Tasks = p.Gen3Tasks
	public.Build76Gen4Tasks = p.Gen4Tasks
	public.Build76PlanePrepMs = p.PlanePrepElapsed.Milliseconds()
	public.Build76FreezeMs = p.FreezeElapsed.Milliseconds()
	public.Build76Prefix1WallMs = p.Prefix1WallElapsed.Milliseconds()
	public.Build76Prefix1WorkerMs = p.Prefix1WorkerElapsed.Milliseconds()
	public.Build76Gen2WallMs = p.Gen2WallElapsed.Milliseconds()
	public.Build76Gen2WorkerMs = p.Gen2WorkerElapsed.Milliseconds()
	public.Build76Gen3WallMs = p.Gen3WallElapsed.Milliseconds()
	public.Build76Gen3WorkerMs = p.Gen3WorkerElapsed.Milliseconds()
	public.Build76Gen4WallMs = p.Gen4WallElapsed.Milliseconds()
	public.Build76Gen4WorkerMs = p.Gen4WorkerElapsed.Milliseconds()
	public.Build76Gen2MinMs = p.Gen2MinElapsed.Milliseconds()
	public.Build76Gen2MedianMs = p.Gen2MedianElapsed.Milliseconds()
	public.Build76Gen2MaxMs = p.Gen2MaxElapsed.Milliseconds()
	public.Build76Gen2MinEvals = p.Gen2MinEvaluations
	public.Build76Gen2MaxEvals = p.Gen2MaxEvaluations
	public.Build76Gen2MinOutputs = p.Gen2MinOutputs
	public.Build76Gen2MaxOutputs = p.Gen2MaxOutputs
	public.Build80Attempted = recovery.Attempted
	public.Build80CacheBlockReads = p.CacheBlockReads
	public.Build80CacheHits = p.CacheHits
	public.Build80CacheFallbacks = p.CacheFallbacks
	public.Build80CacheFailed = p.CacheFailed
	public.Build80CachePixelsPrepared = p.CachePixelsPrepared
	public.Build80CacheMaxArea = p.CacheMaxArea
	public.Build80CacheLimit = experimentalV4PhoneBuild80MaxCachePixels
}
