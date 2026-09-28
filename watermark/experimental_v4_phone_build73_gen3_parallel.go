package watermark

import (
	"errors"
	"image"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Build73 is the qualified smartphone performance baseline selected from Build72
// physical prefix profiling and two independent Go 1.26.0 physical gates. It moves
// the ordered parallel barrier one generation earlier: each seed executes the
// exact Build64 prefix through sibling2, all generation-three inputs are then
// frozen in original seed/traversal order and processed by one bounded global
// worker pool, and the resulting sibling3 hypotheses are committed in original
// task order before the already-qualified Build71 generation-four pool runs.
// No proposal score, threshold, pruning rule, bank contents/order, protected
// decode order or HMAC semantics are changed.

type experimentalV4PhoneBuild73Prefix2Result struct {
	inputs  []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild73Gen3Result struct {
	outputs []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild73GeometryTelemetry struct {
	PlanePrepElapsed     time.Duration
	FreezeElapsed        time.Duration
	Prefix2WallElapsed   time.Duration
	Prefix2WorkerElapsed time.Duration
	Gen3WallElapsed      time.Duration
	Gen3WorkerElapsed    time.Duration
	Gen4WallElapsed      time.Duration
	Gen4WorkerElapsed    time.Duration
	Prefix2Workers       int
	Gen3Workers          int
	Gen4Workers          int
	Gen3Tasks            int
	Gen4Tasks            int
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
}

// experimentalV4PhoneBuild73SeedPrefix2 reproduces Build64 exactly through the
// second sibling stencil and freezes s2 hypotheses in traversal order.
func experimentalV4PhoneBuild73SeedPrefix2(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed experimentalV4PhoneBuild48Seed) ([]experimentalV4PhoneHypothesis, int) {
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
								inputs = append(inputs, s2.hyp)
							}
						}
					}
				}
			}
		}
	}
	return inputs, evals
}

// experimentalV4PhoneBuild73Generation3 executes exactly the Build64 third
// single/pair/continuation/sibling subtree for one frozen s2 input.
func experimentalV4PhoneBuild73Generation3(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int) {
	outputs := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	se3, si3 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, input, 0)
	evals += se3
	if si3 != 0 {
		return outputs, evals
	}
	pairs3, pe3, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, input, 0)
	evals += pe3
	for _, pair3 := range pairs3 {
		cont3, ce3 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair3.hyp, 0)
		evals += ce3
		for _, c3 := range cont3 {
			sibs3, ne3 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c3.hyp, 0)
			evals += ne3
			for _, s3 := range sibs3 {
				outputs = append(outputs, s3.hyp)
			}
		}
	}
	return outputs, evals
}

func experimentalV4PhoneBuild73BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild73GeometryTelemetry) {
	var profile experimentalV4PhoneBuild73GeometryTelemetry
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
		return nil, evals, len(seeds), workers, profile
	}

	// Freeze generation-three inputs in exact seed/traversal order.
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
		return nil, evals, len(seeds), workers, profile
	}

	// Commit generation-three outputs strictly by frozen task index. This is
	// exactly the Build71 s3 traversal order and therefore the gen4 input order.
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
	return bank, evals, len(seeds), workers, profile
}

type experimentalV4PhoneBuild73RecoveryTelemetry struct {
	experimentalV4PhoneBuild68RecoveryTelemetry
	GeometryProfile experimentalV4PhoneBuild73GeometryTelemetry
}

func experimentalV4PhoneBuild73Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild73RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild73RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, geometryProfile := experimentalV4PhoneBuild73BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.GeometryProfile = geometryProfile
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = prefixWorkers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build73 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build73 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build73 recovery authentication failed")
}

func experimentalV4PhoneBuild73ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild73RecoveryTelemetry) {
	experimentalV4PhoneBuild68ApplyTelemetry(public, recovery.experimentalV4PhoneBuild68RecoveryTelemetry)
	// Build73 inherits the qualified Build71 logical semantics. Populate the
	// inherited Build71 view with equivalent aggregate prefix/gen4 measurements
	// while exposing the actual three-barrier scheduler in Build73 fields below.
	p := recovery.GeometryProfile
	public.Build71Attempted = recovery.Attempted
	public.Build71PrefixWorkers = p.Prefix2Workers
	public.Build71Gen4Workers = p.Gen4Workers
	public.Build71Gen4Tasks = p.Gen4Tasks
	public.Build71GeometryPlanePrepMs = p.PlanePrepElapsed.Milliseconds()
	public.Build71GeometryFreezeMs = p.FreezeElapsed.Milliseconds()
	public.Build71PrefixWallMs = (p.Prefix2WallElapsed + p.Gen3WallElapsed).Milliseconds()
	public.Build71PrefixWorkerMs = (p.Prefix2WorkerElapsed + p.Gen3WorkerElapsed).Milliseconds()
	public.Build71Gen4WallMs = p.Gen4WallElapsed.Milliseconds()
	public.Build71Gen4WorkerMs = p.Gen4WorkerElapsed.Milliseconds()
	public.Build71Gen4MinMs = p.Gen4MinElapsed.Milliseconds()
	public.Build71Gen4MedianMs = p.Gen4MedianElapsed.Milliseconds()
	public.Build71Gen4MaxMs = p.Gen4MaxElapsed.Milliseconds()
	public.Build71Gen4MinEvals = p.Gen4MinEvaluations
	public.Build71Gen4MaxEvals = p.Gen4MaxEvaluations
	public.Build71Gen4MinBank = p.Gen4MinBank
	public.Build71Gen4MaxBank = p.Gen4MaxBank
	public.Build73Attempted = recovery.Attempted
	public.Build73Prefix2Workers = p.Prefix2Workers
	public.Build73Gen3Workers = p.Gen3Workers
	public.Build73Gen4Workers = p.Gen4Workers
	public.Build73Gen3Tasks = p.Gen3Tasks
	public.Build73Gen4Tasks = p.Gen4Tasks
	public.Build73GeometryPlanePrepMs = p.PlanePrepElapsed.Milliseconds()
	public.Build73GeometryFreezeMs = p.FreezeElapsed.Milliseconds()
	public.Build73Prefix2WallMs = p.Prefix2WallElapsed.Milliseconds()
	public.Build73Prefix2WorkerMs = p.Prefix2WorkerElapsed.Milliseconds()
	public.Build73Gen3WallMs = p.Gen3WallElapsed.Milliseconds()
	public.Build73Gen3WorkerMs = p.Gen3WorkerElapsed.Milliseconds()
	public.Build73Gen4WallMs = p.Gen4WallElapsed.Milliseconds()
	public.Build73Gen4WorkerMs = p.Gen4WorkerElapsed.Milliseconds()
	public.Build73Gen3MinMs = p.Gen3MinElapsed.Milliseconds()
	public.Build73Gen3MedianMs = p.Gen3MedianElapsed.Milliseconds()
	public.Build73Gen3MaxMs = p.Gen3MaxElapsed.Milliseconds()
	public.Build73Gen3MinEvals = p.Gen3MinEvaluations
	public.Build73Gen3MaxEvals = p.Gen3MaxEvaluations
	public.Build73Gen3MinOutputs = p.Gen3MinOutputs
	public.Build73Gen3MaxOutputs = p.Gen3MaxOutputs
}
