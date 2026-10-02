package watermark

import (
	"errors"
	"image"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Build76 is the qualified equivalence-preserving smartphone baseline promoted
// after two independent Go 1.26.0 physical runs over the Build75 baseline. Build72 already measured the residual Build73 prefix2
// stage-by-stage; after Build75 made the Build47 freeze substantially cheaper,
// the remaining prefix2 wall time is dominated by the second generation on the
// long-running seed. Build76 therefore moves the ordered global barrier one
// generation earlier: each seed executes the exact Build64 prefix through
// sibling1, generation-two inputs are frozen in original seed/traversal order,
// one bounded pool evaluates single2->pair2->continuation2->sibling2, results
// are committed strictly by task index, and the qualified Build73 gen3/gen4
// pools then run unchanged. Build75 ordered-parallel freeze is reused unchanged.

type experimentalV4PhoneBuild76Prefix1Result struct {
	inputs  []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild76Gen2Result struct {
	outputs []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild76GeometryTelemetry struct {
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
}

// experimentalV4PhoneBuild76SeedPrefix1 reproduces Build64 exactly through the
// first sibling stencil and freezes sibling1 hypotheses in traversal order.
func experimentalV4PhoneBuild76SeedPrefix1(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed experimentalV4PhoneBuild48Seed) ([]experimentalV4PhoneHypothesis, int) {
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

// experimentalV4PhoneBuild76Generation2 executes exactly the Build64 second
// single/pair/continuation/sibling subtree for one frozen sibling1 input.
func experimentalV4PhoneBuild76Generation2(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int) {
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

func experimentalV4PhoneBuild76BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild76GeometryTelemetry, experimentalV4PhoneBuild75FreezeTelemetry) {
	var profile experimentalV4PhoneBuild76GeometryTelemetry
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
	for i, r := range gen4Results {
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

type experimentalV4PhoneBuild76RecoveryTelemetry struct {
	experimentalV4PhoneBuild75RecoveryTelemetry
	Build76GeometryProfile experimentalV4PhoneBuild76GeometryTelemetry
}

func experimentalV4PhoneBuild76Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild76RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild76RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }
	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, gp, fp := experimentalV4PhoneBuild76BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.Build76GeometryProfile = gp
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build76 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build76 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build76 recovery authentication failed")
}

func experimentalV4PhoneBuild76ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild76RecoveryTelemetry) {
	experimentalV4PhoneBuild75ApplyTelemetry(public, recovery.experimentalV4PhoneBuild75RecoveryTelemetry)
	p := recovery.Build76GeometryProfile
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
}
