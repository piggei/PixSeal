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

// Build78 is observability-only over the qualified Build76 baseline. It runs
// the exact Build76 freeze/prefix1/gen2/gen3/gen4 schedule and bank order, but
// decomposes generation-four worker time and evaluation counts into the existing
// single4, pair4 and continuation4 operations. Fine-grained timers intentionally
// perturb runtime, so Build78 is not promotable from timing.

type experimentalV4PhoneBuild78ContinueCallProfile struct {
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
}

type experimentalV4PhoneBuild78Gen4Result struct {
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
	continueProfiles []experimentalV4PhoneBuild78ContinueCallProfile
}

type experimentalV4PhoneBuild78Gen4Profile struct {
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

type experimentalV4PhoneBuild78Continuation4Profile struct {
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

func experimentalV4PhoneBuild78ContinueProfiled(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild55BlindContinuationState, int, experimentalV4PhoneBuild78ContinueCallProfile) {
	profile := experimentalV4PhoneBuild78ContinueCallProfile{}
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
				scoreStarted := time.Now()
				ss, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
				profile.scoreElapsed += time.Since(scoreStarted)
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

func experimentalV4PhoneBuild78Generation4Profiled(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) experimentalV4PhoneBuild78Gen4Result {
	r := experimentalV4PhoneBuild78Gen4Result{bank: make([]experimentalV4PhoneHypothesis, 0, 8)}
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
	r.continueProfiles = make([]experimentalV4PhoneBuild78ContinueCallProfile, 0, len(pairs4))
	for _, pair4 := range pairs4 {
		cont4, ce4, cp := experimentalV4PhoneBuild78ContinueProfiled(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
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

func experimentalV4PhoneBuild78BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild76GeometryTelemetry, experimentalV4PhoneBuild75FreezeTelemetry, experimentalV4PhoneBuild78Gen4Profile, experimentalV4PhoneBuild78Continuation4Profile) {
	var profile experimentalV4PhoneBuild76GeometryTelemetry
	var freezeProfile experimentalV4PhoneBuild75FreezeTelemetry
	var gen4Profile experimentalV4PhoneBuild78Gen4Profile
	var continueProfile experimentalV4PhoneBuild78Continuation4Profile
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0, 0, profile, freezeProfile, gen4Profile, continueProfile
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
		return nil, evals, 0, 0, profile, freezeProfile, gen4Profile, continueProfile
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile
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
		return nil, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile
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
	gen4Results := make([]experimentalV4PhoneBuild78Gen4Result, len(flat4))
	gen4Jobs := make(chan int)
	wg = sync.WaitGroup{}
	wg.Add(gen4Workers)
	started = time.Now()
	for w := 0; w < gen4Workers; w++ {
		go func() {
			defer wg.Done()
			for i := range gen4Jobs {
				s := time.Now()
				r := experimentalV4PhoneBuild78Generation4Profiled(plane, pilot, cw, ch, anchor, flat4[i])
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
	return bank, evals, len(seeds), workers, profile, freezeProfile, gen4Profile, continueProfile
}

type experimentalV4PhoneBuild78RecoveryTelemetry struct {
	experimentalV4PhoneBuild76RecoveryTelemetry
	Build78Gen4Profile          experimentalV4PhoneBuild78Gen4Profile
	Build78Continuation4Profile experimentalV4PhoneBuild78Continuation4Profile
}

func experimentalV4PhoneBuild78Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild78RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild78RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }
	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, gp, fp, g4p, c4p := experimentalV4PhoneBuild78BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.Build76GeometryProfile = gp
	telemetry.Build78Gen4Profile = g4p
	telemetry.Build78Continuation4Profile = c4p
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build78 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build78 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build78 recovery authentication failed")
}

func experimentalV4PhoneBuild78ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild78RecoveryTelemetry) {
	experimentalV4PhoneBuild76ApplyTelemetry(public, recovery.experimentalV4PhoneBuild76RecoveryTelemetry)
	p := recovery.Build78Gen4Profile
	public.Build78Attempted = recovery.Attempted
	public.Build78Gen4Tasks = p.Tasks
	public.Build78Gen4Workers = p.Workers
	public.Build78Single4Calls = p.SingleCalls
	public.Build78Single4Accepted = p.SingleAccepted
	public.Build78Pair4Calls = p.PairCalls
	public.Build78Pair4Outputs = p.PairOutputs
	public.Build78Continue4Calls = p.ContinueCalls
	public.Build78Continue4Outputs = p.ContinueOutputs
	public.Build78Single4Evals = p.SingleEvaluations
	public.Build78Pair4Evals = p.PairEvaluations
	public.Build78Continue4Evals = p.ContinueEvaluations
	public.Build78Single4WorkerMs = p.SingleWorkerElapsed.Milliseconds()
	public.Build78Pair4WorkerMs = p.PairWorkerElapsed.Milliseconds()
	public.Build78Continue4WorkerMs = p.ContinueWorkerElapsed.Milliseconds()
	public.Build78DominantGen4Task = p.DominantTaskIndex
	public.Build78DominantGen4Evals = p.DominantTaskEvaluations
	public.Build78DominantGen4Bank = p.DominantTaskBank
	public.Build78DominantSingle4Ms = p.DominantSingleElapsed.Milliseconds()
	public.Build78DominantPair4Ms = p.DominantPairElapsed.Milliseconds()
	public.Build78DominantContinue4Ms = p.DominantContinueElapsed.Milliseconds()
	public.Build78DominantSingle4Evals = p.DominantSingleEvaluations
	public.Build78DominantPair4Evals = p.DominantPairEvaluations
	public.Build78DominantContinue4Evals = p.DominantContinueEvaluations
	c := recovery.Build78Continuation4Profile
	public.Build78Continue4InputStates = c.InputStates
	public.Build78Continue4ProbeAttempts = c.ProbeAttempts
	public.Build78Continue4LimitRejected = c.LimitRejected
	public.Build78Continue4HomographyRejected = c.HomographyRejected
	public.Build78Continue4ScoreEvals = c.ScoreEvaluations
	public.Build78Continue4InvalidScores = c.InvalidScores
	public.Build78Continue4NonImprovingScores = c.NonImprovingScores
	public.Build78Continue4ImprovingProbes = c.ImprovingScoreProbes
	public.Build78Continue4AcceptedStates = c.AcceptedStates
	public.Build78Continue4DuplicateAccepted = c.DuplicateAccepted
	public.Build78Continue4Passes = c.Passes
	public.Build78Continue4WorkerMs = c.WorkerElapsed.Milliseconds()
	public.Build78Continue4ScoreWorkerMs = c.ScoreWorkerElapsed.Milliseconds()
	public.Build78Continue4PrepareWorkerMs = c.PrepareWorkerElapsed.Milliseconds()
	public.Build78Continue4MinMs = c.MinElapsed.Milliseconds()
	public.Build78Continue4MedianMs = c.MedianElapsed.Milliseconds()
	public.Build78Continue4MaxMs = c.MaxElapsed.Milliseconds()
	public.Build78Continue4MinEvals = c.MinEvaluations
	public.Build78Continue4MaxEvals = c.MaxEvaluations
	public.Build78Continue4MinAccepted = c.MinAcceptedStates
	public.Build78Continue4MaxAccepted = c.MaxAcceptedStates
	public.Build78DominantContinueCall = c.DominantCallIndex
	public.Build78DominantContinueGen4Task = c.DominantGen4Task
	public.Build78DominantContinuePairRank = c.DominantPairRank
	public.Build78DominantContinueEvals = c.DominantEvaluations
	public.Build78DominantContinueAccepted = c.DominantAcceptedStates
	public.Build78DominantContinuePasses = c.DominantPasses
	public.Build78DominantContinueMs = c.DominantElapsed.Milliseconds()
	public.Build78DominantContinueScoreMs = c.DominantScoreElapsed.Milliseconds()
	public.Build78DominantContinueProbes = c.DominantProbeAttempts
	public.Build78DominantContinueLimitRejected = c.DominantLimitRejected
	public.Build78DominantContinueHomographyRejected = c.DominantHomographyRejected
	public.Build78DominantContinueInvalidScores = c.DominantInvalidScores
	public.Build78DominantContinueNonImproving = c.DominantNonImprovingScores
	public.Build78DominantContinueImprovingProbes = c.DominantImprovingScoreProbes
	public.Build78DominantContinueDuplicateAccepted = c.DominantDuplicateAccepted

}
