package watermark

import (
	"errors"
	"fmt"
	"image"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// Build72 is observability-only over the qualified Build71 baseline. It keeps
// the exact Build71 two-barrier scheduler and final bank/order, while profiling
// the prefix through sibling3 by public-only stage. No stage timing, protected
// data, key material or HMAC result can alter geometry generation or ranking.

const experimentalV4PhoneBuild72PrefixEvalStageCount = 14
const experimentalV4PhoneBuild72PrefixStateStageCount = 10

type experimentalV4PhoneBuild72PrefixStageProfile struct {
	Evals   [experimentalV4PhoneBuild72PrefixEvalStageCount]int
	States  [experimentalV4PhoneBuild72PrefixStateStageCount]int
	Elapsed [experimentalV4PhoneBuild72PrefixEvalStageCount]time.Duration
}

// Eval/elapsed order:
// baseline,roots,single1,pair1,cont1,sib1,single2,pair2,cont2,sib2,
// single3,pair3,cont3,sib3.
// States order:
// roots,pair1,cont1,sib1,pair2,cont2,sib2,pair3,cont3,sib3.
func experimentalV4PhoneBuild72SeedPrefixProfile(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed experimentalV4PhoneBuild48Seed) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild72PrefixStageProfile) {
	var profile experimentalV4PhoneBuild72PrefixStageProfile
	inputs := make([]experimentalV4PhoneHypothesis, 0, 64)
	evals := 0

	started := time.Now()
	baseline, n := experimentalV4PhoneBuild41Refine(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
	profile.Elapsed[0] += time.Since(started)
	profile.Evals[0] += n
	evals += n

	started = time.Now()
	roots, n := experimentalV4PhoneBuild53TwoPxRoots(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
	profile.Elapsed[1] += time.Since(started)
	profile.Evals[1] += n
	profile.States[0] += len(roots)
	evals += n
	if baseline.h.h[8] == 0 || len(roots) == 0 {
		return inputs, evals, profile
	}

	for _, root := range roots {
		started = time.Now()
		se, si := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
		profile.Elapsed[2] += time.Since(started)
		profile.Evals[2] += se
		evals += se
		if si != 0 {
			continue
		}

		started = time.Now()
		pairs, pe, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
		profile.Elapsed[3] += time.Since(started)
		profile.Evals[3] += pe
		profile.States[1] += len(pairs)
		evals += pe
		for _, pair := range pairs {
			started = time.Now()
			cont, ce := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
			profile.Elapsed[4] += time.Since(started)
			profile.Evals[4] += ce
			profile.States[2] += len(cont)
			evals += ce
			for _, cs := range cont {
				started = time.Now()
				sibs, ne := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, cs.hyp, 0)
				profile.Elapsed[5] += time.Since(started)
				profile.Evals[5] += ne
				profile.States[3] += len(sibs)
				evals += ne
				for _, sib := range sibs {
					started = time.Now()
					se2, si2 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, sib.hyp, 0)
					profile.Elapsed[6] += time.Since(started)
					profile.Evals[6] += se2
					evals += se2
					if si2 != 0 {
						continue
					}

					started = time.Now()
					pairs2, pe2, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, sib.hyp, 0)
					profile.Elapsed[7] += time.Since(started)
					profile.Evals[7] += pe2
					profile.States[4] += len(pairs2)
					evals += pe2
					for _, pair2 := range pairs2 {
						started = time.Now()
						cont2, ce2 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair2.hyp, 0)
						profile.Elapsed[8] += time.Since(started)
						profile.Evals[8] += ce2
						profile.States[5] += len(cont2)
						evals += ce2
						for _, c2 := range cont2 {
							started = time.Now()
							sibs2, ne2 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c2.hyp, 0)
							profile.Elapsed[9] += time.Since(started)
							profile.Evals[9] += ne2
							profile.States[6] += len(sibs2)
							evals += ne2
							for _, s2 := range sibs2 {
								started = time.Now()
								se3, si3 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s2.hyp, 0)
								profile.Elapsed[10] += time.Since(started)
								profile.Evals[10] += se3
								evals += se3
								if si3 != 0 {
									continue
								}

								started = time.Now()
								pairs3, pe3, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s2.hyp, 0)
								profile.Elapsed[11] += time.Since(started)
								profile.Evals[11] += pe3
								profile.States[7] += len(pairs3)
								evals += pe3
								for _, pair3 := range pairs3 {
									started = time.Now()
									cont3, ce3 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair3.hyp, 0)
									profile.Elapsed[12] += time.Since(started)
									profile.Evals[12] += ce3
									profile.States[8] += len(cont3)
									evals += ce3
									for _, c3 := range cont3 {
										started = time.Now()
										sibs3, ne3 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c3.hyp, 0)
										profile.Elapsed[13] += time.Since(started)
										profile.Evals[13] += ne3
										profile.States[9] += len(sibs3)
										evals += ne3
										for _, s3 := range sibs3 {
											inputs = append(inputs, s3.hyp)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return inputs, evals, profile
}

type experimentalV4PhoneBuild72PrefixResult struct {
	inputs  []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
	stage   experimentalV4PhoneBuild72PrefixStageProfile
}

type experimentalV4PhoneBuild72GeometryTelemetry struct {
	experimentalV4PhoneBuild71GeometryTelemetry
	PrefixMinElapsed            time.Duration
	PrefixMedianElapsed         time.Duration
	PrefixMaxElapsed            time.Duration
	PrefixMinEvaluations        int
	PrefixMaxEvaluations        int
	PrefixMinInputs             int
	PrefixMaxInputs             int
	PrefixStageEvaluations      [experimentalV4PhoneBuild72PrefixEvalStageCount]int
	PrefixStageWorkerElapsed    [experimentalV4PhoneBuild72PrefixEvalStageCount]time.Duration
	MaxPrefixSeedIndex          int
	MaxPrefixSeedPairRank       int
	MaxPrefixSeedRankWithinPair int
	MaxPrefixSeedEvaluations    int
	MaxPrefixSeedInputs         int
	MaxPrefixSeedElapsed        time.Duration
	MaxPrefixSeedStage          experimentalV4PhoneBuild72PrefixStageProfile
}

func experimentalV4PhoneBuild72BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild72GeometryTelemetry) {
	var profile experimentalV4PhoneBuild72GeometryTelemetry
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0, 0, profile
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)

	started := time.Now()
	plane := newPixelPlane(work)
	profile.PlanePrepElapsed = time.Since(started)
	pilot := experimentalV4Prototype2Candidate()

	started = time.Now()
	frozen, _, evals := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	profile.FreezeElapsed = time.Since(started)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild63SeedsPerPair)
	if len(seeds) == 0 {
		return nil, evals, 0, 0, profile
	}

	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(seeds) {
		workers = len(seeds)
	}
	profile.PrefixWorkers = workers
	prefixes := make([]experimentalV4PhoneBuild72PrefixResult, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				seedStarted := time.Now()
				inputs, n, stage := experimentalV4PhoneBuild72SeedPrefixProfile(plane, pilot, cw, ch, anchor, seeds[i])
				prefixes[i] = experimentalV4PhoneBuild72PrefixResult{inputs: inputs, evals: n, elapsed: time.Since(seedStarted), stage: stage}
			}
		}()
	}
	for i := range seeds {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	profile.PrefixWallElapsed = time.Since(started)

	totalTasks := 0
	durations := make([]time.Duration, 0, len(prefixes))
	maxIndex := 0
	for i := range prefixes {
		r := prefixes[i]
		evals += r.evals
		profile.PrefixWorkerElapsed += r.elapsed
		totalTasks += len(r.inputs)
		durations = append(durations, r.elapsed)
		if i == 0 || r.elapsed < profile.PrefixMinElapsed {
			profile.PrefixMinElapsed = r.elapsed
		}
		if r.elapsed > profile.PrefixMaxElapsed {
			profile.PrefixMaxElapsed = r.elapsed
		}
		if i == 0 || r.evals < profile.PrefixMinEvaluations {
			profile.PrefixMinEvaluations = r.evals
		}
		if r.evals > profile.PrefixMaxEvaluations {
			profile.PrefixMaxEvaluations = r.evals
		}
		if i == 0 || len(r.inputs) < profile.PrefixMinInputs {
			profile.PrefixMinInputs = len(r.inputs)
		}
		if len(r.inputs) > profile.PrefixMaxInputs {
			profile.PrefixMaxInputs = len(r.inputs)
		}
		for stage := 0; stage < experimentalV4PhoneBuild72PrefixEvalStageCount; stage++ {
			profile.PrefixStageEvaluations[stage] += r.stage.Evals[stage]
			profile.PrefixStageWorkerElapsed[stage] += r.stage.Elapsed[stage]
		}
		if i == 0 || r.evals > prefixes[maxIndex].evals {
			maxIndex = i
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	if len(durations) > 0 {
		profile.PrefixMedianElapsed = durations[len(durations)/2]
	}
	maxSeed := seeds[maxIndex]
	maxResult := prefixes[maxIndex]
	profile.MaxPrefixSeedIndex = maxSeed.index
	profile.MaxPrefixSeedPairRank = maxSeed.frozen.h.build43PairRank
	profile.MaxPrefixSeedRankWithinPair = maxSeed.rankWithinPair
	profile.MaxPrefixSeedEvaluations = maxResult.evals
	profile.MaxPrefixSeedInputs = len(maxResult.inputs)
	profile.MaxPrefixSeedElapsed = maxResult.elapsed
	profile.MaxPrefixSeedStage = maxResult.stage

	profile.Gen4Tasks = totalTasks
	if totalTasks == 0 {
		return nil, evals, len(seeds), workers, profile
	}

	flat := make([]experimentalV4PhoneHypothesis, 0, totalTasks)
	for i := range prefixes {
		flat = append(flat, prefixes[i].inputs...)
	}
	gen4Workers := runtime.GOMAXPROCS(0)
	if gen4Workers < 1 {
		gen4Workers = 1
	}
	if gen4Workers > len(flat) {
		gen4Workers = len(flat)
	}
	profile.Gen4Workers = gen4Workers
	results := make([]experimentalV4PhoneBuild71Gen4Result, len(flat))
	genJobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen4Workers)
	started = time.Now()
	for w := 0; w < gen4Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range genJobs {
				taskStarted := time.Now()
				bank, n := experimentalV4PhoneBuild71Generation4(plane, pilot, cw, ch, anchor, flat[i])
				results[i] = experimentalV4PhoneBuild71Gen4Result{bank: bank, evals: n, elapsed: time.Since(taskStarted)}
			}
		}()
	}
	for i := range flat {
		genJobs <- i
	}
	close(genJobs)
	wg.Wait()
	profile.Gen4WallElapsed = time.Since(started)

	bank := make([]experimentalV4PhoneHypothesis, 0, 1024)
	gen4Durations := make([]time.Duration, 0, len(results))
	for i := range results {
		r := results[i]
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
	return bank, evals, len(seeds), workers, profile
}

type experimentalV4PhoneBuild72RecoveryTelemetry struct {
	experimentalV4PhoneBuild71RecoveryTelemetry
	PrefixProfile experimentalV4PhoneBuild72GeometryTelemetry
}

func experimentalV4PhoneBuild72Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild72RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild72RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, geometryProfile := experimentalV4PhoneBuild72BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.GeometryProfile = geometryProfile.experimentalV4PhoneBuild71GeometryTelemetry
	telemetry.PrefixProfile = geometryProfile
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = prefixWorkers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build72 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build72 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build72 recovery authentication failed")
}

func experimentalV4PhoneBuild72IntsCSV(values []int) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ",")
}

func experimentalV4PhoneBuild72DurationsCSV(values []time.Duration) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("%d", v.Milliseconds())
	}
	return strings.Join(parts, ",")
}

func experimentalV4PhoneBuild72ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild72RecoveryTelemetry) {
	experimentalV4PhoneBuild71ApplyTelemetry(public, recovery.experimentalV4PhoneBuild71RecoveryTelemetry)
	profile := recovery.PrefixProfile
	public.Build72Attempted = recovery.Attempted
	public.Build72PrefixMinMs = profile.PrefixMinElapsed.Milliseconds()
	public.Build72PrefixMedianMs = profile.PrefixMedianElapsed.Milliseconds()
	public.Build72PrefixMaxMs = profile.PrefixMaxElapsed.Milliseconds()
	public.Build72PrefixMinEvals = profile.PrefixMinEvaluations
	public.Build72PrefixMaxEvals = profile.PrefixMaxEvaluations
	public.Build72PrefixMinInputs = profile.PrefixMinInputs
	public.Build72PrefixMaxInputs = profile.PrefixMaxInputs
	public.Build72PrefixStageEvals = experimentalV4PhoneBuild72IntsCSV(profile.PrefixStageEvaluations[:])
	public.Build72PrefixStageWorkerMs = experimentalV4PhoneBuild72DurationsCSV(profile.PrefixStageWorkerElapsed[:])
	public.Build72MaxPrefixSeedIndex = profile.MaxPrefixSeedIndex
	public.Build72MaxPrefixSeedPairRank = profile.MaxPrefixSeedPairRank
	public.Build72MaxPrefixSeedRankWithinPair = profile.MaxPrefixSeedRankWithinPair
	public.Build72MaxPrefixSeedEvals = profile.MaxPrefixSeedEvaluations
	public.Build72MaxPrefixSeedInputs = profile.MaxPrefixSeedInputs
	public.Build72MaxPrefixSeedMs = profile.MaxPrefixSeedElapsed.Milliseconds()
	public.Build72MaxPrefixSeedStageEvals = experimentalV4PhoneBuild72IntsCSV(profile.MaxPrefixSeedStage.Evals[:])
	public.Build72MaxPrefixSeedStageStates = experimentalV4PhoneBuild72IntsCSV(profile.MaxPrefixSeedStage.States[:])
	public.Build72MaxPrefixSeedStageMs = experimentalV4PhoneBuild72DurationsCSV(profile.MaxPrefixSeedStage.Elapsed[:])
}
