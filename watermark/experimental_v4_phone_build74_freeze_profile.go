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

// Build74 is an observability-only snapshot over the qualified Build73 phone
// baseline. It keeps the exact Build73 scheduler and recovery semantics but
// replaces the Build47 freeze call with a line-for-line equivalent profiled
// implementation so the residual freeze cost can be decomposed before any
// optimization is attempted.

type experimentalV4PhoneBuild74FreezeTelemetry struct {
	TotalElapsed               time.Duration
	StructuralElapsed          time.Duration
	PlanePrepElapsed           time.Duration
	PairScoreElapsed           time.Duration
	CellsElapsed               time.Duration
	BasinElapsed               time.Duration
	ProductionBasinElapsed     time.Duration
	DepthBasinElapsed          time.Duration
	AllPairsBasinElapsed       time.Duration
	PairScoreEvaluations       int
	CellEvaluations            int
	BasinEvaluations           int
	ProductionBasinEvaluations int
	DepthBasinEvaluations      int
	AllPairsBasinEvaluations   int
	PairScoreTasks             int
	CellTasks                  int
	BasinTasks                 int
	ProductionBasinTasks       int
	DepthBasinTasks            int
	AllPairsBasinTasks         int
	BasinMinElapsed            time.Duration
	BasinMedianElapsed         time.Duration
	BasinMaxElapsed            time.Duration
	BasinMinEvaluations        int
	BasinMaxEvaluations        int
	BasinMinOutputs            int
	BasinMaxOutputs            int
}

type experimentalV4PhoneBuild74BasinSample struct {
	elapsed time.Duration
	evals   int
	outputs int
}

// experimentalV4PhoneBuild74FreezeProfile is intentionally the exact Build47
// freeze algorithm with timing/evaluation accounting around existing public-only
// operations. The third return value remains Build47's historical proposal count
// (not the helper evaluation total) so qualified geometry counters stay exact.
func experimentalV4PhoneBuild74FreezeProfile(work image.Image, boundary PrintBoundaryEstimate, cw, ch, maxFrozen int) ([]experimentalV4PhoneBuild47Frozen, []ExperimentalV4PhonePairScore, int, experimentalV4PhoneBuild74FreezeTelemetry) {
	var profile experimentalV4PhoneBuild74FreezeTelemetry
	totalStarted := time.Now()
	if work == nil || !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) || maxFrozen <= 0 {
		profile.TotalElapsed = time.Since(totalStarted)
		return nil, nil, 0, profile
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)
	base := experimentalV4PhoneBuild43BaseLines(anchor)

	started := time.Now()
	structural, structOK := experimentalV4PhoneBuild43StructuralLines(work, anchor)
	profile.StructuralElapsed = time.Since(started)

	started = time.Now()
	plane := newPixelPlane(work)
	profile.PlanePrepElapsed = time.Since(started)
	candidate := experimentalV4Prototype2Candidate()
	type rankedPair struct {
		pair  experimentalV4PhoneBuild43Pair
		score float64
		seed  [4]experimentalV4PhoneBuild43Line
	}
	ranked := make([]rankedPair, 0, len(experimentalV4PhoneBuild43Pairs))
	started = time.Now()
	for _, pair := range experimentalV4PhoneBuild43Pairs {
		seed := experimentalV4PhoneBuild43PairSeed(base, structural, structOK, pair)
		score, n := experimentalV4PhoneBuild43PairRobustScore(plane, candidate, cw, ch, anchor, seed, pair)
		profile.PairScoreEvaluations += n
		profile.PairScoreTasks++
		if !math.IsInf(score, -1) {
			ranked = append(ranked, rankedPair{pair: pair, score: score, seed: seed})
		}
	}
	profile.PairScoreElapsed = time.Since(started)
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	scores := make([]ExperimentalV4PhonePairScore, 0, len(ranked))
	for _, rp := range ranked {
		scores = append(scores, ExperimentalV4PhonePairScore{Pair: rp.pair.name, Score: rp.score})
	}

	frozen := make([]experimentalV4PhoneBuild47Frozen, 0, maxFrozen)
	proposalCount := 0
	cellsCache := make(map[int][]experimentalV4PhoneBuild43Cell, len(ranked))
	cellsFor := func(pairRank int) []experimentalV4PhoneBuild43Cell {
		if cells, ok := cellsCache[pairRank]; ok {
			return cells
		}
		rp := ranked[pairRank]
		s := time.Now()
		cells, n := experimentalV4PhoneBuild43Cells(plane, candidate, cw, ch, anchor, rp.seed, rp.pair)
		profile.CellsElapsed += time.Since(s)
		profile.CellEvaluations += n
		profile.CellTasks++
		cellsCache[pairRank] = cells
		return cells
	}
	basinSamples := make([]experimentalV4PhoneBuild74BasinSample, 0, 16)
	appendTier := func(pairStart, pairEnd int, cellRanks []int, tier int) {
		if pairStart < 0 {
			pairStart = 0
		}
		if pairEnd > len(ranked) {
			pairEnd = len(ranked)
		}
		for pairRank := pairStart; pairRank < pairEnd && len(frozen) < maxFrozen; pairRank++ {
			rp := ranked[pairRank]
			cells := cellsFor(pairRank)
			for _, cr := range cellRanks {
				if cr < 0 || cr >= len(cells) || len(frozen) >= maxFrozen {
					continue
				}
				s := time.Now()
				bank, n := experimentalV4PhoneBuild43BasinCandidates(plane, candidate, cw, ch, anchor, rp.seed, rp.pair, cells[cr])
				elapsed := time.Since(s)
				profile.BasinElapsed += elapsed
				profile.BasinEvaluations += n
				profile.BasinTasks++
				switch tier {
				case experimentalV4PhoneBuild47TierProduction:
					profile.ProductionBasinElapsed += elapsed
					profile.ProductionBasinEvaluations += n
					profile.ProductionBasinTasks++
				case experimentalV4PhoneBuild47TierDepth:
					profile.DepthBasinElapsed += elapsed
					profile.DepthBasinEvaluations += n
					profile.DepthBasinTasks++
				case experimentalV4PhoneBuild47TierAllPairs:
					profile.AllPairsBasinElapsed += elapsed
					profile.AllPairsBasinEvaluations += n
					profile.AllPairsBasinTasks++
				}
				basinSamples = append(basinSamples, experimentalV4PhoneBuild74BasinSample{elapsed: elapsed, evals: n, outputs: len(bank)})
				for _, h := range bank {
					proposalCount++
					if len(frozen) >= maxFrozen {
						break
					}
					h.build43Pair = rp.pair.name
					h.build43PairRank = pairRank + 1
					frozen = append(frozen, experimentalV4PhoneBuild47Frozen{h: h, tier: tier, cellRank: cr, cellMean: cells[cr].mean, cellRobust: cells[cr].robust})
				}
			}
		}
	}

	productionPairs := experimentalV4PhoneBuild43PairKeep
	if productionPairs > len(ranked) {
		productionPairs = len(ranked)
	}
	appendTier(0, productionPairs, []int{0, 3}, experimentalV4PhoneBuild47TierProduction)
	appendTier(0, productionPairs, []int{1, 2}, experimentalV4PhoneBuild47TierDepth)
	appendTier(productionPairs, len(ranked), []int{0, 3}, experimentalV4PhoneBuild47TierAllPairs)

	if len(basinSamples) > 0 {
		sort.Slice(basinSamples, func(i, j int) bool { return basinSamples[i].elapsed < basinSamples[j].elapsed })
		profile.BasinMinElapsed = basinSamples[0].elapsed
		profile.BasinMedianElapsed = basinSamples[len(basinSamples)/2].elapsed
		profile.BasinMaxElapsed = basinSamples[len(basinSamples)-1].elapsed
		profile.BasinMinEvaluations = basinSamples[0].evals
		profile.BasinMaxEvaluations = basinSamples[0].evals
		profile.BasinMinOutputs = basinSamples[0].outputs
		profile.BasinMaxOutputs = basinSamples[0].outputs
		for _, b := range basinSamples[1:] {
			if b.evals < profile.BasinMinEvaluations {
				profile.BasinMinEvaluations = b.evals
			}
			if b.evals > profile.BasinMaxEvaluations {
				profile.BasinMaxEvaluations = b.evals
			}
			if b.outputs < profile.BasinMinOutputs {
				profile.BasinMinOutputs = b.outputs
			}
			if b.outputs > profile.BasinMaxOutputs {
				profile.BasinMaxOutputs = b.outputs
			}
		}
	}
	profile.TotalElapsed = time.Since(totalStarted)
	return frozen, scores, proposalCount, profile
}

func experimentalV4PhoneBuild74BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild73GeometryTelemetry, experimentalV4PhoneBuild74FreezeTelemetry) {
	var profile experimentalV4PhoneBuild73GeometryTelemetry
	var freezeProfile experimentalV4PhoneBuild74FreezeTelemetry
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0, 0, profile, freezeProfile
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)

	started := time.Now()
	plane := newPixelPlane(work)
	profile.PlanePrepElapsed = time.Since(started)
	pilot := experimentalV4Prototype2Candidate()

	frozen, _, evals, freezeProfile := experimentalV4PhoneBuild74FreezeProfile(work, boundary, cw, ch, 128)
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
	profile.Prefix2Workers = workers
	prefixes := make([]experimentalV4PhoneBuild73Prefix2Result, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				s := time.Now()
				inputs, n := experimentalV4PhoneBuild73SeedPrefix2(plane, pilot, cw, ch, anchor, seeds[i])
				prefixes[i] = experimentalV4PhoneBuild73Prefix2Result{inputs: inputs, evals: n, elapsed: time.Since(s)}
			}
		}()
	}
	for i := range seeds {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	profile.Prefix2WallElapsed = time.Since(started)

	totalGen3 := 0
	for i := range prefixes {
		evals += prefixes[i].evals
		profile.Prefix2WorkerElapsed += prefixes[i].elapsed
		totalGen3 += len(prefixes[i].inputs)
	}
	profile.Gen3Tasks = totalGen3
	if totalGen3 == 0 {
		return nil, evals, len(seeds), workers, profile, freezeProfile
	}

	flat3 := make([]experimentalV4PhoneHypothesis, 0, totalGen3)
	for i := range prefixes {
		flat3 = append(flat3, prefixes[i].inputs...)
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

	totalGen4 := 0
	gen3Durations := make([]time.Duration, 0, len(gen3Results))
	for i := range gen3Results {
		r := gen3Results[i]
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
	gen4Results := make([]experimentalV4PhoneBuild71Gen4Result, len(flat4))
	gen4Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen4Workers)
	started = time.Now()
	for w := 0; w < gen4Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen4Jobs {
				s := time.Now()
				bank, n := experimentalV4PhoneBuild71Generation4(plane, pilot, cw, ch, anchor, flat4[i])
				gen4Results[i] = experimentalV4PhoneBuild71Gen4Result{bank: bank, evals: n, elapsed: time.Since(s)}
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
	for i := range gen4Results {
		r := gen4Results[i]
		evals += r.evals
		bank = append(bank, r.bank...)
		profile.Gen4WorkerElapsed += r.elapsed
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

type experimentalV4PhoneBuild74RecoveryTelemetry struct {
	experimentalV4PhoneBuild73RecoveryTelemetry
	FreezeProfile experimentalV4PhoneBuild74FreezeTelemetry
}

func experimentalV4PhoneBuild74Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild74RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild74RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, geometryProfile, freezeProfile := experimentalV4PhoneBuild74BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.GeometryProfile = geometryProfile
	telemetry.FreezeProfile = freezeProfile
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = prefixWorkers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build74 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build74 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build74 recovery authentication failed")
}

func experimentalV4PhoneBuild74ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild74RecoveryTelemetry) {
	experimentalV4PhoneBuild73ApplyTelemetry(public, recovery.experimentalV4PhoneBuild73RecoveryTelemetry)
	p := recovery.FreezeProfile
	public.Build74Attempted = recovery.Attempted
	public.Build74FreezeTotalMs = p.TotalElapsed.Milliseconds()
	public.Build74FreezeStructuralMs = p.StructuralElapsed.Milliseconds()
	public.Build74FreezePlanePrepMs = p.PlanePrepElapsed.Milliseconds()
	public.Build74FreezePairScoreMs = p.PairScoreElapsed.Milliseconds()
	public.Build74FreezeCellsMs = p.CellsElapsed.Milliseconds()
	public.Build74FreezeBasinMs = p.BasinElapsed.Milliseconds()
	public.Build74FreezeProductionBasinMs = p.ProductionBasinElapsed.Milliseconds()
	public.Build74FreezeDepthBasinMs = p.DepthBasinElapsed.Milliseconds()
	public.Build74FreezeAllPairsBasinMs = p.AllPairsBasinElapsed.Milliseconds()
	public.Build74FreezePairScoreEvals = p.PairScoreEvaluations
	public.Build74FreezeCellEvals = p.CellEvaluations
	public.Build74FreezeBasinEvals = p.BasinEvaluations
	public.Build74FreezeProductionBasinEvals = p.ProductionBasinEvaluations
	public.Build74FreezeDepthBasinEvals = p.DepthBasinEvaluations
	public.Build74FreezeAllPairsBasinEvals = p.AllPairsBasinEvaluations
	public.Build74FreezePairScoreTasks = p.PairScoreTasks
	public.Build74FreezeCellTasks = p.CellTasks
	public.Build74FreezeBasinTasks = p.BasinTasks
	public.Build74FreezeProductionBasinTasks = p.ProductionBasinTasks
	public.Build74FreezeDepthBasinTasks = p.DepthBasinTasks
	public.Build74FreezeAllPairsBasinTasks = p.AllPairsBasinTasks
	public.Build74FreezeBasinMinMs = p.BasinMinElapsed.Milliseconds()
	public.Build74FreezeBasinMedianMs = p.BasinMedianElapsed.Milliseconds()
	public.Build74FreezeBasinMaxMs = p.BasinMaxElapsed.Milliseconds()
	public.Build74FreezeBasinMinEvals = p.BasinMinEvaluations
	public.Build74FreezeBasinMaxEvals = p.BasinMaxEvaluations
	public.Build74FreezeBasinMinOutputs = p.BasinMinOutputs
	public.Build74FreezeBasinMaxOutputs = p.BasinMaxOutputs
}
