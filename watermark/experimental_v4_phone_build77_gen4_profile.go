package watermark

import (
	"errors"
	"image"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Build77 is observability-only over the qualified Build76 baseline. It runs
// the exact Build76 freeze/prefix1/gen2/gen3/gen4 schedule and bank order, but
// decomposes generation-four worker time and evaluation counts into the existing
// single4, pair4 and continuation4 operations. Fine-grained timers intentionally
// perturb runtime, so Build77 is not promotable from timing.

type experimentalV4PhoneBuild77Gen4Result struct {
	bank            []experimentalV4PhoneHypothesis
	evals           int
	elapsed         time.Duration
	singleElapsed   time.Duration
	pairElapsed     time.Duration
	continueElapsed time.Duration
	singleEvals     int
	pairEvals       int
	continueEvals   int
	singleAccepted  bool
	pairOutputs     int
	continueCalls   int
	continueOutputs int
}

type experimentalV4PhoneBuild77Gen4Profile struct {
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

func experimentalV4PhoneBuild77Generation4Profiled(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) experimentalV4PhoneBuild77Gen4Result {
	r := experimentalV4PhoneBuild77Gen4Result{bank: make([]experimentalV4PhoneHypothesis, 0, 8)}
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
	for _, pair4 := range pairs4 {
		started = time.Now()
		cont4, ce4 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
		r.continueElapsed += time.Since(started)
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

func experimentalV4PhoneBuild77BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild76GeometryTelemetry, experimentalV4PhoneBuild75FreezeTelemetry, experimentalV4PhoneBuild77Gen4Profile) {
	var profile experimentalV4PhoneBuild76GeometryTelemetry
	var freezeProfile experimentalV4PhoneBuild75FreezeTelemetry
	var gen4Profile experimentalV4PhoneBuild77Gen4Profile
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0, 0, profile, freezeProfile, gen4Profile
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
		return nil, evals, 0, 0, profile, freezeProfile, gen4Profile
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile
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
	gen4Results := make([]experimentalV4PhoneBuild77Gen4Result, len(flat4))
	gen4Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen4Workers)
	started = time.Now()
	for w := 0; w < gen4Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen4Jobs {
				s := time.Now()
				r := experimentalV4PhoneBuild77Generation4Profiled(plane, pilot, cw, ch, anchor, flat4[i])
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
	return bank, evals, len(seeds), workers, profile, freezeProfile, gen4Profile
}

type experimentalV4PhoneBuild77RecoveryTelemetry struct {
	experimentalV4PhoneBuild76RecoveryTelemetry
	Build77Gen4Profile experimentalV4PhoneBuild77Gen4Profile
}

func experimentalV4PhoneBuild77Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild77RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild77RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }
	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, gp, fp, g4p := experimentalV4PhoneBuild77BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.Build76GeometryProfile = gp
	telemetry.Build77Gen4Profile = g4p
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build77 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build77 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build77 recovery authentication failed")
}

func experimentalV4PhoneBuild77ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild77RecoveryTelemetry) {
	experimentalV4PhoneBuild76ApplyTelemetry(public, recovery.experimentalV4PhoneBuild76RecoveryTelemetry)
	p := recovery.Build77Gen4Profile
	public.Build77Attempted = recovery.Attempted
	public.Build77Gen4Tasks = p.Tasks
	public.Build77Gen4Workers = p.Workers
	public.Build77Single4Calls = p.SingleCalls
	public.Build77Single4Accepted = p.SingleAccepted
	public.Build77Pair4Calls = p.PairCalls
	public.Build77Pair4Outputs = p.PairOutputs
	public.Build77Continue4Calls = p.ContinueCalls
	public.Build77Continue4Outputs = p.ContinueOutputs
	public.Build77Single4Evals = p.SingleEvaluations
	public.Build77Pair4Evals = p.PairEvaluations
	public.Build77Continue4Evals = p.ContinueEvaluations
	public.Build77Single4WorkerMs = p.SingleWorkerElapsed.Milliseconds()
	public.Build77Pair4WorkerMs = p.PairWorkerElapsed.Milliseconds()
	public.Build77Continue4WorkerMs = p.ContinueWorkerElapsed.Milliseconds()
	public.Build77DominantGen4Task = p.DominantTaskIndex
	public.Build77DominantGen4Evals = p.DominantTaskEvaluations
	public.Build77DominantGen4Bank = p.DominantTaskBank
	public.Build77DominantSingle4Ms = p.DominantSingleElapsed.Milliseconds()
	public.Build77DominantPair4Ms = p.DominantPairElapsed.Milliseconds()
	public.Build77DominantContinue4Ms = p.DominantContinueElapsed.Milliseconds()
	public.Build77DominantSingle4Evals = p.DominantSingleEvaluations
	public.Build77DominantPair4Evals = p.DominantPairEvaluations
	public.Build77DominantContinue4Evals = p.DominantContinueEvaluations
}
