package watermark

import (
	"image"
	"math"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Build95 is an experimental, benchmark-first successor candidate selected
// from the qualified Build94 profile. Build94 showed generation 4 dominating
// B/angle wall time and the Build84 continuation4 RGB-to-luminance reader
// accounting for about one third of sampled CPU. Build95 therefore changes only
// continuation4 projective reads: it materializes one exact full-image float64
// luminance plane immediately before generation 4, then reuses those immutable
// values for every Build84 continuation4 bilinear sample.
//
// Prefix1, generations 2/3, generation4 single/pair probes, bank ordering,
// qualification, protected decode and all security boundaries remain Build93.
// If the deterministic public-dimension memory cap would be exceeded, Build95
// falls back to the unchanged Build84 continuation4 reader.

const experimentalV4PhoneBuild95MaxLuminancePlaneBytes int64 = 128 << 20

type experimentalV4PhoneBuild95LuminancePlane struct {
	bounds image.Rectangle
	luma   []float64
}

type experimentalV4PhoneBuild95GeometryTelemetry struct {
	experimentalV4PhoneBuild84GeometryTelemetry
	LuminancePlanePrepElapsed time.Duration
	LuminancePlaneBytes       int64
	LuminancePlaneUsed        bool
	LuminancePlaneFallback    bool
}

// experimentalV4PhoneBuild95NewLuminancePlane preserves the exact Build84
// per-pixel expression and evaluation order. It is keyed only by public image
// dimensions; no key, payload, ECC, HMAC or correctness signal can affect
// allocation or fallback.
func experimentalV4PhoneBuild95NewLuminancePlane(src *pixelPlane) (*experimentalV4PhoneBuild95LuminancePlane, int64, bool) {
	if src == nil || src.bounds.Dx() <= 0 || src.bounds.Dy() <= 0 {
		return nil, 0, true
	}
	width, height := src.bounds.Dx(), src.bounds.Dy()
	pixels := int64(width) * int64(height)
	bytes := pixels * 8
	if pixels <= 0 || bytes <= 0 || bytes > experimentalV4PhoneBuild95MaxLuminancePlaneBytes || int64(len(src.rgb)) != pixels*3 {
		return nil, bytes, true
	}
	out := &experimentalV4PhoneBuild95LuminancePlane{
		bounds: src.bounds,
		luma:   make([]float64, int(pixels)),
	}
	for i := 0; i < int(pixels); i++ {
		j := i * 3
		out.luma[i] = .299*float64(src.rgb[j]) + .587*float64(src.rgb[j+1]) + .114*float64(src.rgb[j+2]) - 128
	}
	return out, bytes, false
}

func experimentalV4PhoneBuild95ReadProjectiveBlockValue(src *experimentalV4PhoneBuild95LuminancePlane, h homography, originX, originY, size int) (float64, bool) {
	if src == nil {
		return 0, false
	}
	width, height := src.bounds.Dx(), src.bounds.Dy()
	table := readCosTables[size]
	table2, table3 := table[2], table[3]
	c23, c32 := 0.0, 0.0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
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

			row0 := y0 * width
			row1 := y1 * width
			l00 := src.luma[row0+x0]
			l10 := src.luma[row0+x1]
			l01 := src.luma[row1+x0]
			l11 := src.luma[row1+x1]
			top := l00*(1-fx) + l10*fx
			bottom := l01*(1-fx) + l11*fx
			l := top*(1-fy) + bottom*fy

			c23 += l * table3[x] * table2[y]
			c32 += l * table2[x] * table3[y]
		}
	}
	return math.Abs(c23) - math.Abs(c32), true
}

func experimentalV4PhoneBuild95FoldScoreFetch(luma *experimentalV4PhoneBuild95LuminancePlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, heldout int, proposal bool, sparse bool) (float64, int, experimentalV4PhoneBuild84FetchTelemetry) {
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	fetch.FoldScores = 1
	if luma == nil {
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
				v, ok := experimentalV4PhoneBuild95ReadProjectiveBlockValue(luma, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
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

func experimentalV4PhoneBuild95ContinueFetch(luma *experimentalV4PhoneBuild95LuminancePlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild55BlindContinuationState, int, experimentalV4PhoneBuild84FetchTelemetry) {
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
				ss, _, ip := experimentalV4PhoneBuild95FoldScoreFetch(luma, candidate, cw, ch, hh, heldout, true, false)
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

func experimentalV4PhoneBuild95Generation4Fetch(rgb *pixelPlane, luma *experimentalV4PhoneBuild95LuminancePlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild84FetchTelemetry) {
	bank := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	se4, si4 := experimentalV4PhoneBuild53SingleProbe(rgb, pilot, cw, ch, anchor, input, 0)
	evals += se4
	if si4 != 0 {
		return bank, evals, fetch
	}
	pairs4, pe4, _ := experimentalV4PhoneBuild53PairStencil(rgb, pilot, cw, ch, anchor, input, 0)
	evals += pe4
	for _, pair4 := range pairs4 {
		var cont4 []experimentalV4PhoneBuild55BlindContinuationState
		var ce4 int
		var ip experimentalV4PhoneBuild84FetchTelemetry
		if luma != nil {
			cont4, ce4, ip = experimentalV4PhoneBuild95ContinueFetch(luma, pilot, cw, ch, anchor, pair4.hyp, 0)
		} else {
			cont4, ce4, ip = experimentalV4PhoneBuild84ContinueFetch(rgb, pilot, cw, ch, anchor, pair4.hyp, 0)
		}
		experimentalV4PhoneBuild84MergeFetch(&fetch, ip)
		evals += ce4
		for _, c4 := range cont4 {
			bank = append(bank, c4.hyp)
		}
	}
	return bank, evals, fetch
}

func experimentalV4PhoneBuild95BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild95GeometryTelemetry, experimentalV4PhoneBuild75FreezeTelemetry) {
	var profile experimentalV4PhoneBuild95GeometryTelemetry
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
	prefixes := make([]experimentalV4PhoneBuild84Prefix1Result, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				s := time.Now()
				inputs, n := experimentalV4PhoneBuild84SeedPrefix1(plane, pilot, cw, ch, anchor, seeds[i])
				prefixes[i] = experimentalV4PhoneBuild84Prefix1Result{inputs: inputs, evals: n, elapsed: time.Since(s)}
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
	gen2Results := make([]experimentalV4PhoneBuild84Gen2Result, len(flat2))
	gen2Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen2Workers)
	started = time.Now()
	for w := 0; w < gen2Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen2Jobs {
				s := time.Now()
				outputs, n := experimentalV4PhoneBuild84Generation2(plane, pilot, cw, ch, anchor, flat2[i])
				gen2Results[i] = experimentalV4PhoneBuild84Gen2Result{outputs: outputs, evals: n, elapsed: time.Since(s)}
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

	lumaStarted := time.Now()
	lumaPlane, lumaBytes, lumaFallback := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	profile.LuminancePlanePrepElapsed = time.Since(lumaStarted)
	profile.LuminancePlaneBytes = lumaBytes
	profile.LuminancePlaneUsed = lumaPlane != nil
	profile.LuminancePlaneFallback = lumaFallback

	gen4Workers := runtime.GOMAXPROCS(0)
	if gen4Workers < 1 {
		gen4Workers = 1
	}
	if gen4Workers > len(flat4) {
		gen4Workers = len(flat4)
	}
	profile.Gen4Workers = gen4Workers
	gen4Results := make([]experimentalV4PhoneBuild84Gen4Result, len(flat4))
	gen4Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen4Workers)
	started = time.Now()
	for w := 0; w < gen4Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen4Jobs {
				s := time.Now()
				bank, n, lp := experimentalV4PhoneBuild95Generation4Fetch(plane, lumaPlane, pilot, cw, ch, anchor, flat4[i])
				gen4Results[i] = experimentalV4PhoneBuild84Gen4Result{bank: bank, evals: n, elapsed: time.Since(s), fetch: lp}
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
		profile.FetchFoldScores += r.fetch.FoldScores
		profile.FetchBlockReads += r.fetch.BlockReads
		profile.FetchBlockSuccess += r.fetch.BlockSuccess
		profile.FetchBlockFailed += r.fetch.BlockFailed
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
