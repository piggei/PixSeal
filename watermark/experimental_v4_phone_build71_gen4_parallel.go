package watermark

import (
	"errors"
	"image"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Build71 is an equivalence-preserving performance experiment selected from
// Build69/70 physical profiling. The qualified Build68 semantics remain the
// reference. Build71 changes only proposal-only scheduling: every seed executes
// the exact Build64 prefix through the third sibling stencil, the resulting
// fourth-generation inputs are frozen in original seed/traversal order, and the
// independent single4 -> pair4 -> cont4 subtrees are evaluated in one bounded
// global worker pool. Results are committed strictly in original task order.

type experimentalV4PhoneBuild71PrefixResult struct {
	inputs  []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild71Gen4Result struct {
	bank    []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild71GeometryTelemetry struct {
	PlanePrepElapsed    time.Duration
	FreezeElapsed       time.Duration
	PrefixWallElapsed   time.Duration
	PrefixWorkerElapsed time.Duration
	Gen4WallElapsed     time.Duration
	Gen4WorkerElapsed   time.Duration
	PrefixWorkers       int
	Gen4Workers         int
	Gen4Tasks           int
	Gen4MinElapsed      time.Duration
	Gen4MedianElapsed   time.Duration
	Gen4MaxElapsed      time.Duration
	Gen4MinEvaluations  int
	Gen4MaxEvaluations  int
	Gen4MinBank         int
	Gen4MaxBank         int
}

// experimentalV4PhoneBuild71SeedPrefix reproduces Build64 exactly through the
// third sibling stencil. Instead of immediately entering generation four, it
// freezes the s3 hypotheses in the same traversal order.
func experimentalV4PhoneBuild71SeedPrefix(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed experimentalV4PhoneBuild48Seed) ([]experimentalV4PhoneHypothesis, int) {
	inputs := make([]experimentalV4PhoneHypothesis, 0, 64)
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
					se2, si2 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, sib.hyp, 0)
					evals += se2
					if si2 != 0 {
						continue
					}
					pairs2, pe2, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, sib.hyp, 0)
					evals += pe2
					for _, pair2 := range pairs2 {
						cont2, ce2 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair2.hyp, 0)
						evals += ce2
						for _, c2 := range cont2 {
							sibs2, ne2 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c2.hyp, 0)
							evals += ne2
							for _, s2 := range sibs2 {
								se3, si3 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s2.hyp, 0)
								evals += se3
								if si3 != 0 {
									continue
								}
								pairs3, pe3, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s2.hyp, 0)
								evals += pe3
								for _, pair3 := range pairs3 {
									cont3, ce3 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair3.hyp, 0)
									evals += ce3
									for _, c3 := range cont3 {
										sibs3, ne3 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c3.hyp, 0)
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
	return inputs, evals
}

func experimentalV4PhoneBuild71Generation4(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int) {
	bank := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	se4, si4 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, input, 0)
	evals += se4
	if si4 != 0 {
		return bank, evals
	}
	pairs4, pe4, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, input, 0)
	evals += pe4
	for _, pair4 := range pairs4 {
		cont4, ce4 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
		evals += ce4
		for _, c4 := range cont4 {
			bank = append(bank, c4.hyp)
		}
	}
	return bank, evals
}

// experimentalV4PhoneBuild71BlindBank preserves the exact Build65/Build68
// logical bank while scheduling all fourth-generation subtrees in a single
// bounded worker pool. Prefix results and generation-four results are both
// stored by original index and concatenated only after all workers complete.
func experimentalV4PhoneBuild71BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild71GeometryTelemetry) {
	var profile experimentalV4PhoneBuild71GeometryTelemetry
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
	prefixes := make([]experimentalV4PhoneBuild71PrefixResult, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				seedStarted := time.Now()
				inputs, n := experimentalV4PhoneBuild71SeedPrefix(plane, pilot, cw, ch, anchor, seeds[i])
				prefixes[i] = experimentalV4PhoneBuild71PrefixResult{inputs: inputs, evals: n, elapsed: time.Since(seedStarted)}
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
	for i := range prefixes {
		evals += prefixes[i].evals
		profile.PrefixWorkerElapsed += prefixes[i].elapsed
		totalTasks += len(prefixes[i].inputs)
	}
	profile.Gen4Tasks = totalTasks
	if totalTasks == 0 {
		return nil, evals, len(seeds), workers, profile
	}

	// Flatten in exact seed/traversal order. The flat index is therefore also
	// the final commit order for the qualified Build64 bank.
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
	durations := make([]time.Duration, 0, len(results))
	for i := range results {
		r := results[i]
		evals += r.evals
		bank = append(bank, r.bank...)
		profile.Gen4WorkerElapsed += r.elapsed
		durations = append(durations, r.elapsed)
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
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	if len(durations) > 0 {
		profile.Gen4MedianElapsed = durations[len(durations)/2]
	}
	return bank, evals, len(seeds), workers, profile
}

type experimentalV4PhoneBuild71RecoveryTelemetry struct {
	experimentalV4PhoneBuild68RecoveryTelemetry
	GeometryProfile experimentalV4PhoneBuild71GeometryTelemetry
}

func experimentalV4PhoneBuild71Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild71RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild71RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, geometryProfile := experimentalV4PhoneBuild71BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.GeometryProfile = geometryProfile
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = prefixWorkers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build71 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build71 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build71 recovery authentication failed")
}

func experimentalV4PhoneBuild71ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild71RecoveryTelemetry) {
	experimentalV4PhoneBuild68ApplyTelemetry(public, recovery.experimentalV4PhoneBuild68RecoveryTelemetry)
	public.Build71Attempted = recovery.Attempted
	public.Build71PrefixWorkers = recovery.GeometryProfile.PrefixWorkers
	public.Build71Gen4Workers = recovery.GeometryProfile.Gen4Workers
	public.Build71Gen4Tasks = recovery.GeometryProfile.Gen4Tasks
	public.Build71GeometryPlanePrepMs = recovery.GeometryProfile.PlanePrepElapsed.Milliseconds()
	public.Build71GeometryFreezeMs = recovery.GeometryProfile.FreezeElapsed.Milliseconds()
	public.Build71PrefixWallMs = recovery.GeometryProfile.PrefixWallElapsed.Milliseconds()
	public.Build71PrefixWorkerMs = recovery.GeometryProfile.PrefixWorkerElapsed.Milliseconds()
	public.Build71Gen4WallMs = recovery.GeometryProfile.Gen4WallElapsed.Milliseconds()
	public.Build71Gen4WorkerMs = recovery.GeometryProfile.Gen4WorkerElapsed.Milliseconds()
	public.Build71Gen4MinMs = recovery.GeometryProfile.Gen4MinElapsed.Milliseconds()
	public.Build71Gen4MedianMs = recovery.GeometryProfile.Gen4MedianElapsed.Milliseconds()
	public.Build71Gen4MaxMs = recovery.GeometryProfile.Gen4MaxElapsed.Milliseconds()
	public.Build71Gen4MinEvals = recovery.GeometryProfile.Gen4MinEvaluations
	public.Build71Gen4MaxEvals = recovery.GeometryProfile.Gen4MaxEvaluations
	public.Build71Gen4MinBank = recovery.GeometryProfile.Gen4MinBank
	public.Build71Gen4MaxBank = recovery.GeometryProfile.Gen4MaxBank
}
