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

// Build79 is observability-only over the qualified Build76 baseline. It runs
// the exact Build76 freeze/prefix1/gen2/gen3/gen4 schedule and bank order, while
// instrumenting only the FoldScore kernel invoked by continuation4. Exact work
// counters cover every score evaluation; rare deterministic samples decompose
// block reads into mapPoint, bilinear luminance sampling and DCT work. Timing is
// diagnostic only, so Build79 is not promotable from runtime.

type experimentalV4PhoneBuild79FoldCallProfile struct {
	elapsed                 time.Duration
	tilesVisited            int
	pilotPositions          int
	blockReads              int
	blockSuccess            int
	blockFailed             int
	visible                 int
	sampled                 bool
	sampledBlockReadElapsed time.Duration
	detailedBlocks          int
	detailedPixels          int
	replayMapElapsed        time.Duration
	replaySampleElapsed     time.Duration
	replayDCTElapsed        time.Duration
	replayFailures          int
}

type experimentalV4PhoneBuild79ContinueCallProfile struct {
	elapsed              time.Duration
	scoreElapsed         time.Duration
	passes               int
	probeAttempts        int
	limitRejected        int
	homographyRejected   int
	invalidScores        int
	nonImprovingScores   int
	improvingScoreProbes int
	acceptedStates       int
	duplicateAccepted    int
	evals                int
	foldProfiles         []experimentalV4PhoneBuild79FoldCallProfile
}

type experimentalV4PhoneBuild79Gen4Result struct {
	bank             []experimentalV4PhoneHypothesis
	evals            int
	elapsed          time.Duration
	singleElapsed    time.Duration
	pairElapsed      time.Duration
	continueElapsed  time.Duration
	singleEvals      int
	pairEvals        int
	continueEvals    int
	singleAccepted   bool
	pairOutputs      int
	continueCalls    int
	continueOutputs  int
	continueProfiles []experimentalV4PhoneBuild79ContinueCallProfile
}

type experimentalV4PhoneBuild79Gen4Profile struct {
	Tasks                       int
	Workers                     int
	SingleCalls                 int
	SingleAccepted              int
	PairCalls                   int
	PairOutputs                 int
	ContinueCalls               int
	ContinueOutputs             int
	SingleEvaluations           int
	PairEvaluations             int
	ContinueEvaluations         int
	SingleWorkerElapsed         time.Duration
	PairWorkerElapsed           time.Duration
	ContinueWorkerElapsed       time.Duration
	DominantTaskIndex           int
	DominantTaskEvaluations     int
	DominantTaskBank            int
	DominantSingleElapsed       time.Duration
	DominantPairElapsed         time.Duration
	DominantContinueElapsed     time.Duration
	DominantSingleEvaluations   int
	DominantPairEvaluations     int
	DominantContinueEvaluations int
}

type experimentalV4PhoneBuild79Continuation4Profile struct {
	Calls                        int
	InputStates                  int
	ProbeAttempts                int
	LimitRejected                int
	HomographyRejected           int
	ScoreEvaluations             int
	InvalidScores                int
	NonImprovingScores           int
	ImprovingScoreProbes         int
	AcceptedStates               int
	DuplicateAccepted            int
	Passes                       int
	WorkerElapsed                time.Duration
	ScoreWorkerElapsed           time.Duration
	PrepareWorkerElapsed         time.Duration
	MinElapsed                   time.Duration
	MedianElapsed                time.Duration
	MaxElapsed                   time.Duration
	MinEvaluations               int
	MaxEvaluations               int
	MinAcceptedStates            int
	MaxAcceptedStates            int
	DominantCallIndex            int
	DominantGen4Task             int
	DominantPairRank             int
	DominantEvaluations          int
	DominantAcceptedStates       int
	DominantPasses               int
	DominantElapsed              time.Duration
	DominantScoreElapsed         time.Duration
	DominantProbeAttempts        int
	DominantLimitRejected        int
	DominantHomographyRejected   int
	DominantInvalidScores        int
	DominantNonImprovingScores   int
	DominantImprovingScoreProbes int
	DominantDuplicateAccepted    int
}

// Build79 samples one FoldScore evaluation in every 64 deterministic score
// ordinals (mixed with the public gen4 task and pair indices). The sample rule
// depends only on public traversal order; timing and protected data never
// influence it.
const experimentalV4PhoneBuild79KernelSampleMod = 64

type experimentalV4PhoneBuild79FoldKernelProfile struct {
	FoldCalls               int
	TilesVisited            int
	PilotPositions          int
	BlockReads              int
	BlockSuccess            int
	BlockFailed             int
	Visible                 int
	FoldWorkerElapsed       time.Duration
	SampledFoldCalls        int
	SampledFoldElapsed      time.Duration
	SampledBlockReadElapsed time.Duration
	SampledFoldOverhead     time.Duration
	DetailedBlocks          int
	DetailedPixels          int
	ReplayMapElapsed        time.Duration
	ReplaySampleElapsed     time.Duration
	ReplayDCTElapsed        time.Duration
	ReplayFailures          int
}

func experimentalV4PhoneBuild79ShouldSampleFold(gen4Task, pairRank, scoreOrdinal int) bool {
	return positiveMod(gen4Task*131+pairRank*17+scoreOrdinal, experimentalV4PhoneBuild79KernelSampleMod) == 0
}

func experimentalV4PhoneBuild79ProfileBlockReplay(src *pixelPlane, h homography, originX, originY, size int) (mapElapsed, sampleElapsed, dctElapsed time.Duration, pixels int, ok bool) {
	if src == nil || size <= 0 {
		return 0, 0, 0, 0, false
	}
	xs := make([]float64, size*size)
	ys := make([]float64, size*size)
	values := make([]float64, size*size)
	started := time.Now()
	idx := 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx, sy, valid := h.mapPoint(float64(originX+x), float64(originY+y))
			if !valid {
				return time.Since(started), 0, 0, idx, false
			}
			xs[idx], ys[idx] = sx, sy
			idx++
		}
	}
	mapElapsed = time.Since(started)
	started = time.Now()
	for i := range xs {
		v, valid := samplePlaneLuminance(src, xs[i], ys[i])
		if !valid {
			return mapElapsed, time.Since(started), 0, i, false
		}
		values[i] = v
	}
	sampleElapsed = time.Since(started)
	table := readCosTables[size]
	started = time.Now()
	c23, c32 := 0.0, 0.0
	idx = 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			l := values[idx]
			c23 += l * table[3][x] * table[2][y]
			c32 += l * table[2][x] * table[3][y]
			idx++
		}
	}
	_ = math.Abs(c23) - math.Abs(c32)
	dctElapsed = time.Since(started)
	return mapElapsed, sampleElapsed, dctElapsed, len(values), true
}

// experimentalV4PhoneBuild79FoldScoreProfiled is a line-for-line arithmetic
// copy of Build41 FoldScore. Instrumentation counts all work. Rare deterministic
// samples time block reads; a single successful block per sampled FoldScore is
// replayed after the authoritative score is complete to decompose mapPoint,
// bilinear luminance sampling and DCT accumulation without affecting output.
func experimentalV4PhoneBuild79FoldScoreProfiled(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, heldout int, proposal bool, sparse bool, sample bool) (float64, int, experimentalV4PhoneBuild79FoldCallProfile) {
	profile := experimentalV4PhoneBuild79FoldCallProfile{sampled: sample}
	started := time.Now()
	if plane == nil {
		profile.elapsed = time.Since(started)
		return math.Inf(-1), 0, profile
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	if tilesX < 1 || tilesY < 1 {
		profile.elapsed = time.Since(started)
		return math.Inf(-1), 0, profile
	}
	num, den := 0.0, 0.0
	visible := 0
	usedTiles := 0
	detailOriginX, detailOriginY := 0, 0
	haveDetail := false
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
			profile.tilesVisited++
			baseX := tx * experimentalV4TileWidthBlocks
			baseY := ty * experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				if sparse && i%4 != 0 {
					continue
				}
				profile.pilotPositions++
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				originX := (baseX + px) * blockSize
				originY := (baseY + py) * blockSize
				profile.blockReads++
				var v float64
				var ok bool
				if sample {
					blockStarted := time.Now()
					v, ok = readProjectiveBlockValue(plane, h, originX, originY, blockSize)
					profile.sampledBlockReadElapsed += time.Since(blockStarted)
				} else {
					v, ok = readProjectiveBlockValue(plane, h, originX, originY, blockSize)
				}
				if !ok {
					profile.blockFailed++
					continue
				}
				profile.blockSuccess++
				if sample && !haveDetail {
					detailOriginX, detailOriginY = originX, originY
					haveDetail = true
				}
				visible++
				num += float64(candidate.signs[i]) * v
				den += math.Abs(v)
			}
		}
	}
	profile.visible = visible
	profile.elapsed = time.Since(started)
	if sample {
		profileOverhead := profile.elapsed - profile.sampledBlockReadElapsed
		_ = profileOverhead
		if haveDetail {
			m, sp, d, pixels, ok := experimentalV4PhoneBuild79ProfileBlockReplay(plane, h, detailOriginX, detailOriginY, blockSize)
			profile.detailedBlocks = 1
			profile.detailedPixels = pixels
			profile.replayMapElapsed = m
			profile.replaySampleElapsed = sp
			profile.replayDCTElapsed = d
			if !ok {
				profile.replayFailures = 1
			}
		}
	}
	if den <= 0 {
		return math.Inf(-1), visible, profile
	}
	return num / den, visible, profile
}

func experimentalV4PhoneBuild79ContinueProfiled(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout, gen4Task, pairRank int) ([]experimentalV4PhoneBuild55BlindContinuationState, int, experimentalV4PhoneBuild79ContinueCallProfile) {
	profile := experimentalV4PhoneBuild79ContinueCallProfile{}
	started := time.Now()
	q, h, score := start.quad, start.h, start.proposal
	states := make([]experimentalV4PhoneBuild55BlindContinuationState, 0, experimentalV4PhoneBuild55MaxStatesBranch)
	evals := 0
	seen := make(map[[4]ImagePoint]struct{}, experimentalV4PhoneBuild55MaxStatesBranch)
	for pass := 0; pass < experimentalV4PhoneBuild55MaxPasses; pass++ {
		profile.passes++
		improved := false
		for dim := 0; dim < experimentalV4PhoneBuild55Dimensions; dim++ {
			bestQ, bestH, best := q, h, score
			bestDelta := 0.0
			for _, sign := range []float64{-1, 1} {
				profile.probeAttempts++
				delta := sign * experimentalV4PhoneBuild55Step
				qq := experimentalV4PhoneBuild51ApplyDim(q, dim, delta)
				if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
					profile.limitRejected++
					continue
				}
				hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
				if !ok {
					profile.homographyRejected++
					continue
				}
				sample := experimentalV4PhoneBuild79ShouldSampleFold(gen4Task, pairRank, evals)
				ss, _, foldProfile := experimentalV4PhoneBuild79FoldScoreProfiled(plane, candidate, cw, ch, hh, heldout, true, false, sample)
				profile.foldProfiles = append(profile.foldProfiles, foldProfile)
				profile.scoreElapsed += foldProfile.elapsed
				evals++
				if math.IsInf(ss, 0) || math.IsNaN(ss) {
					profile.invalidScores++
					continue
				}
				if ss > best+1e-7 {
					profile.improvingScoreProbes++
					bestQ, bestH, best, bestDelta = qq, hh, ss, delta
				} else {
					profile.nonImprovingScores++
				}
			}
			if best > score+1e-7 {
				q, h, score = bestQ, bestH, best
				improved = true
				if _, exists := seen[q]; exists {
					profile.duplicateAccepted++
				} else {
					seen[q] = struct{}{}
				}
				states = append(states, experimentalV4PhoneBuild55BlindContinuationState{
					index: len(states) + 1,
					pass:  pass + 1,
					dim:   dim,
					delta: bestDelta,
					hyp:   experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score},
				})
				profile.acceptedStates++
			}
		}
		if !improved {
			break
		}
	}
	profile.evals = evals
	profile.elapsed = time.Since(started)
	return states, evals, profile
}

func experimentalV4PhoneBuild79Generation4Profiled(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis, gen4Task int) experimentalV4PhoneBuild79Gen4Result {
	r := experimentalV4PhoneBuild79Gen4Result{bank: make([]experimentalV4PhoneHypothesis, 0, 8)}
	started := time.Now()
	se4, si4 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, input, 0)
	r.singleElapsed = time.Since(started)
	r.singleEvals = se4
	r.evals += se4
	if si4 != 0 {
		return r
	}
	r.singleAccepted = true
	started = time.Now()
	pairs4, pe4, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, input, 0)
	r.pairElapsed = time.Since(started)
	r.pairEvals = pe4
	r.evals += pe4
	r.pairOutputs = len(pairs4)
	r.continueProfiles = make([]experimentalV4PhoneBuild79ContinueCallProfile, 0, len(pairs4))
	for pairRank, pair4 := range pairs4 {
		cont4, ce4, cp := experimentalV4PhoneBuild79ContinueProfiled(plane, pilot, cw, ch, anchor, pair4.hyp, 0, gen4Task, pairRank)
		r.continueProfiles = append(r.continueProfiles, cp)
		r.continueElapsed += cp.elapsed
		r.continueEvals += ce4
		r.evals += ce4
		r.continueCalls++
		r.continueOutputs += len(cont4)
		for _, c4 := range cont4 {
			r.bank = append(r.bank, c4.hyp)
		}
	}
	return r
}

func experimentalV4PhoneBuild79BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild76GeometryTelemetry, experimentalV4PhoneBuild75FreezeTelemetry, experimentalV4PhoneBuild79Gen4Profile, experimentalV4PhoneBuild79Continuation4Profile, experimentalV4PhoneBuild79FoldKernelProfile) {
	var profile experimentalV4PhoneBuild76GeometryTelemetry
	var freezeProfile experimentalV4PhoneBuild75FreezeTelemetry
	var gen4Profile experimentalV4PhoneBuild79Gen4Profile
	var continueProfile experimentalV4PhoneBuild79Continuation4Profile
	var kernelProfile experimentalV4PhoneBuild79FoldKernelProfile
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0, 0, profile, freezeProfile, gen4Profile, continueProfile, kernelProfile
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
		return nil, evals, 0, 0, profile, freezeProfile, gen4Profile, continueProfile, kernelProfile
	}

	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(seeds) {
		workers = len(seeds)
	}
	profile.Prefix1Workers = workers
	prefixes := make([]experimentalV4PhoneBuild76Prefix1Result, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				s := time.Now()
				inputs, n := experimentalV4PhoneBuild76SeedPrefix1(plane, pilot, cw, ch, anchor, seeds[i])
				prefixes[i] = experimentalV4PhoneBuild76Prefix1Result{inputs: inputs, evals: n, elapsed: time.Since(s)}
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile, kernelProfile
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
	gen2Results := make([]experimentalV4PhoneBuild76Gen2Result, len(flat2))
	gen2Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen2Workers)
	started = time.Now()
	for w := 0; w < gen2Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen2Jobs {
				s := time.Now()
				outputs, n := experimentalV4PhoneBuild76Generation2(plane, pilot, cw, ch, anchor, flat2[i])
				gen2Results[i] = experimentalV4PhoneBuild76Gen2Result{outputs: outputs, evals: n, elapsed: time.Since(s)}
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile, kernelProfile
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile, kernelProfile
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
	gen4Results := make([]experimentalV4PhoneBuild79Gen4Result, len(flat4))
	gen4Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen4Workers)
	started = time.Now()
	for w := 0; w < gen4Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen4Jobs {
				s := time.Now()
				r := experimentalV4PhoneBuild79Generation4Profiled(plane, pilot, cw, ch, anchor, flat4[i], i)
				r.elapsed = time.Since(s)
				gen4Results[i] = r
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
	continueDurations := make([]time.Duration, 0, gen4Profile.ContinueCalls)
	gen4Profile.Tasks = len(gen4Results)
	gen4Profile.Workers = gen4Workers
	for i, r := range gen4Results {
		evals += r.evals
		bank = append(bank, r.bank...)
		profile.Gen4WorkerElapsed += r.elapsed
		gen4Durations = append(gen4Durations, r.elapsed)
		gen4Profile.SingleCalls++
		gen4Profile.SingleEvaluations += r.singleEvals
		gen4Profile.PairEvaluations += r.pairEvals
		gen4Profile.ContinueEvaluations += r.continueEvals
		gen4Profile.SingleWorkerElapsed += r.singleElapsed
		gen4Profile.PairWorkerElapsed += r.pairElapsed
		gen4Profile.ContinueWorkerElapsed += r.continueElapsed
		if r.singleAccepted {
			gen4Profile.SingleAccepted++
			gen4Profile.PairCalls++
		}
		gen4Profile.PairOutputs += r.pairOutputs
		gen4Profile.ContinueCalls += r.continueCalls
		gen4Profile.ContinueOutputs += r.continueOutputs
		for pairRank, cp := range r.continueProfiles {
			callIndex := continueProfile.Calls
			continueDurations = append(continueDurations, cp.elapsed)
			continueProfile.Calls++
			continueProfile.InputStates++
			continueProfile.ProbeAttempts += cp.probeAttempts
			continueProfile.LimitRejected += cp.limitRejected
			continueProfile.HomographyRejected += cp.homographyRejected
			continueProfile.ScoreEvaluations += cp.evals
			continueProfile.InvalidScores += cp.invalidScores
			continueProfile.NonImprovingScores += cp.nonImprovingScores
			continueProfile.ImprovingScoreProbes += cp.improvingScoreProbes
			continueProfile.AcceptedStates += cp.acceptedStates
			continueProfile.DuplicateAccepted += cp.duplicateAccepted
			continueProfile.Passes += cp.passes
			continueProfile.WorkerElapsed += cp.elapsed
			continueProfile.ScoreWorkerElapsed += cp.scoreElapsed
			for _, fold := range cp.foldProfiles {
				kernelProfile.FoldCalls++
				kernelProfile.TilesVisited += fold.tilesVisited
				kernelProfile.PilotPositions += fold.pilotPositions
				kernelProfile.BlockReads += fold.blockReads
				kernelProfile.BlockSuccess += fold.blockSuccess
				kernelProfile.BlockFailed += fold.blockFailed
				kernelProfile.Visible += fold.visible
				kernelProfile.FoldWorkerElapsed += fold.elapsed
				if fold.sampled {
					kernelProfile.SampledFoldCalls++
					kernelProfile.SampledFoldElapsed += fold.elapsed
					kernelProfile.SampledBlockReadElapsed += fold.sampledBlockReadElapsed
					overhead := fold.elapsed - fold.sampledBlockReadElapsed
					if overhead > 0 {
						kernelProfile.SampledFoldOverhead += overhead
					}
					kernelProfile.DetailedBlocks += fold.detailedBlocks
					kernelProfile.DetailedPixels += fold.detailedPixels
					kernelProfile.ReplayMapElapsed += fold.replayMapElapsed
					kernelProfile.ReplaySampleElapsed += fold.replaySampleElapsed
					kernelProfile.ReplayDCTElapsed += fold.replayDCTElapsed
					kernelProfile.ReplayFailures += fold.replayFailures
				}
			}
			prep := cp.elapsed - cp.scoreElapsed
			if prep > 0 {
				continueProfile.PrepareWorkerElapsed += prep
			}
			if callIndex == 0 || cp.elapsed < continueProfile.MinElapsed {
				continueProfile.MinElapsed = cp.elapsed
			}
			if cp.elapsed > continueProfile.MaxElapsed {
				continueProfile.MaxElapsed = cp.elapsed
			}
			if callIndex == 0 || cp.evals < continueProfile.MinEvaluations {
				continueProfile.MinEvaluations = cp.evals
			}
			if cp.evals > continueProfile.MaxEvaluations {
				continueProfile.MaxEvaluations = cp.evals
			}
			if callIndex == 0 || cp.acceptedStates < continueProfile.MinAcceptedStates {
				continueProfile.MinAcceptedStates = cp.acceptedStates
			}
			if cp.acceptedStates > continueProfile.MaxAcceptedStates {
				continueProfile.MaxAcceptedStates = cp.acceptedStates
			}
			if callIndex == 0 || cp.evals > continueProfile.DominantEvaluations {
				continueProfile.DominantCallIndex = callIndex
				continueProfile.DominantGen4Task = i
				continueProfile.DominantPairRank = pairRank
				continueProfile.DominantEvaluations = cp.evals
				continueProfile.DominantAcceptedStates = cp.acceptedStates
				continueProfile.DominantPasses = cp.passes
				continueProfile.DominantElapsed = cp.elapsed
				continueProfile.DominantScoreElapsed = cp.scoreElapsed
				continueProfile.DominantProbeAttempts = cp.probeAttempts
				continueProfile.DominantLimitRejected = cp.limitRejected
				continueProfile.DominantHomographyRejected = cp.homographyRejected
				continueProfile.DominantInvalidScores = cp.invalidScores
				continueProfile.DominantNonImprovingScores = cp.nonImprovingScores
				continueProfile.DominantImprovingScoreProbes = cp.improvingScoreProbes
				continueProfile.DominantDuplicateAccepted = cp.duplicateAccepted
			}
		}
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
		if i == 0 || r.evals > gen4Profile.DominantTaskEvaluations {
			gen4Profile.DominantTaskIndex = i
			gen4Profile.DominantTaskEvaluations = r.evals
			gen4Profile.DominantTaskBank = len(r.bank)
			gen4Profile.DominantSingleElapsed = r.singleElapsed
			gen4Profile.DominantPairElapsed = r.pairElapsed
			gen4Profile.DominantContinueElapsed = r.continueElapsed
			gen4Profile.DominantSingleEvaluations = r.singleEvals
			gen4Profile.DominantPairEvaluations = r.pairEvals
			gen4Profile.DominantContinueEvaluations = r.continueEvals
		}
	}
	sort.Slice(gen4Durations, func(i, j int) bool { return gen4Durations[i] < gen4Durations[j] })
	if len(gen4Durations) > 0 {
		profile.Gen4MedianElapsed = gen4Durations[len(gen4Durations)/2]
	}
	sort.Slice(continueDurations, func(i, j int) bool { return continueDurations[i] < continueDurations[j] })
	if len(continueDurations) > 0 {
		continueProfile.MedianElapsed = continueDurations[len(continueDurations)/2]
	}
	return bank, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile, kernelProfile
}

type experimentalV4PhoneBuild79RecoveryTelemetry struct {
	experimentalV4PhoneBuild76RecoveryTelemetry
	Build79Gen4Profile          experimentalV4PhoneBuild79Gen4Profile
	Build79Continuation4Profile experimentalV4PhoneBuild79Continuation4Profile
	Build79FoldKernelProfile    experimentalV4PhoneBuild79FoldKernelProfile
}

func experimentalV4PhoneBuild79Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild79RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild79RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }
	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, gp, fp, g4p, c4p, kp := experimentalV4PhoneBuild79BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.Build76GeometryProfile = gp
	telemetry.Build79Gen4Profile = g4p
	telemetry.Build79Continuation4Profile = c4p
	telemetry.Build79FoldKernelProfile = kp
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build79 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build79 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build79 recovery authentication failed")
}

func experimentalV4PhoneBuild79ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild79RecoveryTelemetry) {
	experimentalV4PhoneBuild76ApplyTelemetry(public, recovery.experimentalV4PhoneBuild76RecoveryTelemetry)
	p := recovery.Build79Gen4Profile
	public.Build79Attempted = recovery.Attempted
	public.Build79Gen4Tasks = p.Tasks
	public.Build79Gen4Workers = p.Workers
	public.Build79Single4Calls = p.SingleCalls
	public.Build79Single4Accepted = p.SingleAccepted
	public.Build79Pair4Calls = p.PairCalls
	public.Build79Pair4Outputs = p.PairOutputs
	public.Build79Continue4Calls = p.ContinueCalls
	public.Build79Continue4Outputs = p.ContinueOutputs
	public.Build79Single4Evals = p.SingleEvaluations
	public.Build79Pair4Evals = p.PairEvaluations
	public.Build79Continue4Evals = p.ContinueEvaluations
	public.Build79Single4WorkerMs = p.SingleWorkerElapsed.Milliseconds()
	public.Build79Pair4WorkerMs = p.PairWorkerElapsed.Milliseconds()
	public.Build79Continue4WorkerMs = p.ContinueWorkerElapsed.Milliseconds()
	public.Build79DominantGen4Task = p.DominantTaskIndex
	public.Build79DominantGen4Evals = p.DominantTaskEvaluations
	public.Build79DominantGen4Bank = p.DominantTaskBank
	public.Build79DominantSingle4Ms = p.DominantSingleElapsed.Milliseconds()
	public.Build79DominantPair4Ms = p.DominantPairElapsed.Milliseconds()
	public.Build79DominantContinue4Ms = p.DominantContinueElapsed.Milliseconds()
	public.Build79DominantSingle4Evals = p.DominantSingleEvaluations
	public.Build79DominantPair4Evals = p.DominantPairEvaluations
	public.Build79DominantContinue4Evals = p.DominantContinueEvaluations
	c := recovery.Build79Continuation4Profile
	public.Build79Continue4InputStates = c.InputStates
	public.Build79Continue4ProbeAttempts = c.ProbeAttempts
	public.Build79Continue4LimitRejected = c.LimitRejected
	public.Build79Continue4HomographyRejected = c.HomographyRejected
	public.Build79Continue4ScoreEvals = c.ScoreEvaluations
	public.Build79Continue4InvalidScores = c.InvalidScores
	public.Build79Continue4NonImprovingScores = c.NonImprovingScores
	public.Build79Continue4ImprovingProbes = c.ImprovingScoreProbes
	public.Build79Continue4AcceptedStates = c.AcceptedStates
	public.Build79Continue4DuplicateAccepted = c.DuplicateAccepted
	public.Build79Continue4Passes = c.Passes
	public.Build79Continue4WorkerMs = c.WorkerElapsed.Milliseconds()
	public.Build79Continue4ScoreWorkerMs = c.ScoreWorkerElapsed.Milliseconds()
	public.Build79Continue4PrepareWorkerMs = c.PrepareWorkerElapsed.Milliseconds()
	public.Build79Continue4MinMs = c.MinElapsed.Milliseconds()
	public.Build79Continue4MedianMs = c.MedianElapsed.Milliseconds()
	public.Build79Continue4MaxMs = c.MaxElapsed.Milliseconds()
	public.Build79Continue4MinEvals = c.MinEvaluations
	public.Build79Continue4MaxEvals = c.MaxEvaluations
	public.Build79Continue4MinAccepted = c.MinAcceptedStates
	public.Build79Continue4MaxAccepted = c.MaxAcceptedStates
	public.Build79DominantContinueCall = c.DominantCallIndex
	public.Build79DominantContinueGen4Task = c.DominantGen4Task
	public.Build79DominantContinuePairRank = c.DominantPairRank
	public.Build79DominantContinueEvals = c.DominantEvaluations
	public.Build79DominantContinueAccepted = c.DominantAcceptedStates
	public.Build79DominantContinuePasses = c.DominantPasses
	public.Build79DominantContinueMs = c.DominantElapsed.Milliseconds()
	public.Build79DominantContinueScoreMs = c.DominantScoreElapsed.Milliseconds()
	public.Build79DominantContinueProbes = c.DominantProbeAttempts
	public.Build79DominantContinueLimitRejected = c.DominantLimitRejected
	public.Build79DominantContinueHomographyRejected = c.DominantHomographyRejected
	public.Build79DominantContinueInvalidScores = c.DominantInvalidScores
	public.Build79DominantContinueNonImproving = c.DominantNonImprovingScores
	public.Build79DominantContinueImprovingProbes = c.DominantImprovingScoreProbes
	public.Build79DominantContinueDuplicateAccepted = c.DominantDuplicateAccepted
	k := recovery.Build79FoldKernelProfile
	public.Build79FoldCalls = k.FoldCalls
	public.Build79FoldTiles = k.TilesVisited
	public.Build79FoldPilotPositions = k.PilotPositions
	public.Build79FoldBlockReads = k.BlockReads
	public.Build79FoldBlockSuccess = k.BlockSuccess
	public.Build79FoldBlockFailed = k.BlockFailed
	public.Build79FoldVisible = k.Visible
	public.Build79FoldWorkerMs = k.FoldWorkerElapsed.Milliseconds()
	public.Build79KernelSampleMod = experimentalV4PhoneBuild79KernelSampleMod
	public.Build79SampledFoldCalls = k.SampledFoldCalls
	public.Build79SampledFoldMs = k.SampledFoldElapsed.Milliseconds()
	public.Build79SampledBlockReadMs = k.SampledBlockReadElapsed.Milliseconds()
	public.Build79SampledFoldOverheadMs = k.SampledFoldOverhead.Milliseconds()
	public.Build79DetailedBlocks = k.DetailedBlocks
	public.Build79DetailedPixels = k.DetailedPixels
	public.Build79ReplayMapMs = k.ReplayMapElapsed.Milliseconds()
	public.Build79ReplaySampleMs = k.ReplaySampleElapsed.Milliseconds()
	public.Build79ReplayDCTMs = k.ReplayDCTElapsed.Milliseconds()
	public.Build79ReplayFailures = k.ReplayFailures

}
