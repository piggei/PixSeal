package watermark

import (
	"image"
	"math"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Build100 changes only the execution schedule of the six fixed public pair
// scores in Build75 freeze. All results are committed in original pair order
// before the historical stable sort. No protected input is accepted here.
type experimentalV4PhoneBuild100PairTelemetry struct {
	Workers       int
	Tasks         int
	Evaluations   int
	WallElapsed   time.Duration
	WorkerElapsed time.Duration
}

type experimentalV4PhoneBuild100RankedPair struct {
	pair  experimentalV4PhoneBuild43Pair
	score float64
	seed  [4]experimentalV4PhoneBuild43Line
}

func experimentalV4PhoneBuild100ScorePairs(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, base, structural [4]experimentalV4PhoneBuild43Line, structOK [4]bool, requestedWorkers int) ([]experimentalV4PhoneBuild100RankedPair, experimentalV4PhoneBuild100PairTelemetry) {
	started := time.Now()
	// Snapshot immutable public work before starting any score calculation.
	pairs := append([]experimentalV4PhoneBuild43Pair(nil), experimentalV4PhoneBuild43Pairs...)
	type result struct {
		ranked  experimentalV4PhoneBuild100RankedPair
		evals   int
		elapsed time.Duration
	}
	results := make([]result, len(pairs))
	for i, pair := range pairs {
		results[i].ranked.pair = pair
		results[i].ranked.seed = experimentalV4PhoneBuild43PairSeed(base, structural, structOK, pair)
	}
	workers := requestedWorkers
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	if workers > len(pairs) {
		workers = len(pairs)
	}
	profile := experimentalV4PhoneBuild100PairTelemetry{Workers: workers, Tasks: len(pairs)}
	if len(pairs) == 0 {
		profile.WallElapsed = time.Since(started)
		return nil, profile
	}
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				item := &results[i]
				s := time.Now()
				item.ranked.score, item.evals = experimentalV4PhoneBuild43PairRobustScore(plane, candidate, cw, ch, anchor, item.ranked.seed, item.ranked.pair)
				item.elapsed = time.Since(s)
			}
		}()
	}
	for i := range pairs {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	ranked := make([]experimentalV4PhoneBuild100RankedPair, 0, len(pairs))
	for _, item := range results {
		profile.Evaluations += item.evals
		profile.WorkerElapsed += item.elapsed
		if !math.IsInf(item.ranked.score, -1) {
			ranked = append(ranked, item.ranked)
		}
	}
	profile.WallElapsed = time.Since(started)
	return ranked, profile
}

func experimentalV4PhoneBuild100FreezeParallel(work image.Image, boundary PrintBoundaryEstimate, cw, ch, maxFrozen int) ([]experimentalV4PhoneBuild47Frozen, []ExperimentalV4PhonePairScore, int, experimentalV4PhoneBuild75FreezeTelemetry, experimentalV4PhoneBuild100PairTelemetry) {
	var profile experimentalV4PhoneBuild75FreezeTelemetry
	var pairProfile experimentalV4PhoneBuild100PairTelemetry
	totalStarted := time.Now()
	if work == nil || !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) || maxFrozen <= 0 {
		profile.TotalElapsed = time.Since(totalStarted)
		return nil, nil, 0, profile, pairProfile
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
	ranked, pairProfile := experimentalV4PhoneBuild100ScorePairs(plane, candidate, cw, ch, anchor, base, structural, structOK, 0)
	profile.PairScoreEvaluations = pairProfile.Evaluations
	profile.PairScoreTasks = pairProfile.Tasks
	profile.PairScoreElapsed = pairProfile.WallElapsed
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	scores := make([]ExperimentalV4PhonePairScore, 0, len(ranked))
	for _, rp := range ranked {
		scores = append(scores, ExperimentalV4PhonePairScore{Pair: rp.pair.name, Score: rp.score})
	}

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

	type basinTask struct {
		pairRank int
		cellRank int
		tier     int
		pair     experimentalV4PhoneBuild43Pair
		seed     [4]experimentalV4PhoneBuild43Line
		cell     experimentalV4PhoneBuild43Cell
	}
	tasks := make([]basinTask, 0, 16)
	appendTierTasks := func(pairStart, pairEnd int, cellRanks []int, tier int) {
		if pairStart < 0 {
			pairStart = 0
		}
		if pairEnd > len(ranked) {
			pairEnd = len(ranked)
		}
		for pairRank := pairStart; pairRank < pairEnd; pairRank++ {
			rp := ranked[pairRank]
			cells := cellsFor(pairRank)
			for _, cr := range cellRanks {
				if cr < 0 || cr >= len(cells) {
					continue
				}
				tasks = append(tasks, basinTask{pairRank: pairRank, cellRank: cr, tier: tier, pair: rp.pair, seed: rp.seed, cell: cells[cr]})
			}
		}
	}
	productionPairs := experimentalV4PhoneBuild43PairKeep
	if productionPairs > len(ranked) {
		productionPairs = len(ranked)
	}
	appendTierTasks(0, productionPairs, []int{0, 3}, experimentalV4PhoneBuild47TierProduction)
	appendTierTasks(0, productionPairs, []int{1, 2}, experimentalV4PhoneBuild47TierDepth)
	appendTierTasks(productionPairs, len(ranked), []int{0, 3}, experimentalV4PhoneBuild47TierAllPairs)

	profile.BasinTasks = len(tasks)
	if len(tasks) == 0 {
		profile.TotalElapsed = time.Since(totalStarted)
		return nil, scores, 0, profile, pairProfile
	}
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(tasks) {
		workers = len(tasks)
	}
	profile.BasinWorkers = workers
	results := make([]experimentalV4PhoneBuild75BasinResult, len(tasks))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				task := tasks[i]
				s := time.Now()
				bank, n := experimentalV4PhoneBuild43BasinCandidates(plane, candidate, cw, ch, anchor, task.seed, task.pair, task.cell)
				results[i] = experimentalV4PhoneBuild75BasinResult{bank: bank, evals: n, elapsed: time.Since(s)}
			}
		}()
	}
	for i := range tasks {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	profile.BasinWallElapsed = time.Since(started)

	durations := make([]time.Duration, 0, len(results))
	for i, r := range results {
		task := tasks[i]
		profile.BasinWorkerElapsed += r.elapsed
		profile.BasinEvaluations += r.evals
		durations = append(durations, r.elapsed)
		if i == 0 || r.elapsed < profile.BasinMinElapsed {
			profile.BasinMinElapsed = r.elapsed
		}
		if r.elapsed > profile.BasinMaxElapsed {
			profile.BasinMaxElapsed = r.elapsed
		}
		if i == 0 || r.evals < profile.BasinMinEvaluations {
			profile.BasinMinEvaluations = r.evals
		}
		if r.evals > profile.BasinMaxEvaluations {
			profile.BasinMaxEvaluations = r.evals
		}
		if i == 0 || len(r.bank) < profile.BasinMinOutputs {
			profile.BasinMinOutputs = len(r.bank)
		}
		if len(r.bank) > profile.BasinMaxOutputs {
			profile.BasinMaxOutputs = len(r.bank)
		}
		switch task.tier {
		case experimentalV4PhoneBuild47TierProduction:
			profile.ProductionBasinTasks++
			profile.ProductionBasinEvaluations += r.evals
			profile.ProductionBasinWorkerElapsed += r.elapsed
		case experimentalV4PhoneBuild47TierDepth:
			profile.DepthBasinTasks++
			profile.DepthBasinEvaluations += r.evals
			profile.DepthBasinWorkerElapsed += r.elapsed
		case experimentalV4PhoneBuild47TierAllPairs:
			profile.AllPairsBasinTasks++
			profile.AllPairsBasinEvaluations += r.evals
			profile.AllPairsBasinWorkerElapsed += r.elapsed
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	if len(durations) > 0 {
		profile.BasinMedianElapsed = durations[len(durations)/2]
	}

	// Semantic commit is strictly historical Build47 order. Physical results
	// from tasks beyond the point where maxFrozen is reached are ignored.
	frozen := make([]experimentalV4PhoneBuild47Frozen, 0, maxFrozen)
	proposalCount := 0
	for i, r := range results {
		task := tasks[i]
		for _, h0 := range r.bank {
			proposalCount++
			if len(frozen) >= maxFrozen {
				break
			}
			h := h0
			h.build43Pair = task.pair.name
			h.build43PairRank = task.pairRank + 1
			frozen = append(frozen, experimentalV4PhoneBuild47Frozen{h: h, tier: task.tier, cellRank: task.cellRank, cellMean: task.cell.mean, cellRobust: task.cell.robust})
		}
		if len(frozen) >= maxFrozen {
			break
		}
	}
	profile.TotalElapsed = time.Since(totalStarted)
	return frozen, scores, proposalCount, profile, pairProfile
}
