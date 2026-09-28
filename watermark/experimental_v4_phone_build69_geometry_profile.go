package watermark

import (
	"errors"
	"image"
	"runtime"
	"sort"
	"sync"
	"time"
)

// experimentalV4PhoneBuild69SeedResult records only observational timing around
// the unchanged Build64 per-seed geometry search. Bank content and eval counts
// remain produced by experimentalV4PhoneBuild64SeedBank.
type experimentalV4PhoneBuild69SeedResult struct {
	bank    []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
}

type experimentalV4PhoneBuild69GeometryTelemetry struct {
	PlanePrepElapsed   time.Duration
	FreezeElapsed      time.Duration
	SeedWallElapsed    time.Duration
	SeedWorkerElapsed  time.Duration
	SeedMinElapsed     time.Duration
	SeedMedianElapsed  time.Duration
	SeedMaxElapsed     time.Duration
	SeedMinEvaluations int
	SeedMaxEvaluations int
	SeedMinBank        int
	SeedMaxBank        int
}

// experimentalV4PhoneBuild69BlindBank reproduces Build65BlindBank exactly but
// adds timing around the already-existing stages and per-seed calls. Results are
// still stored by original seed index and concatenated only after all workers
// finish, preserving Build64/65 bank content and order.
func experimentalV4PhoneBuild69BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild69GeometryTelemetry) {
	var profile experimentalV4PhoneBuild69GeometryTelemetry
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
	results := make([]experimentalV4PhoneBuild69SeedResult, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				seedStarted := time.Now()
				bank, n := experimentalV4PhoneBuild64SeedBank(plane, pilot, cw, ch, anchor, seeds[i])
				results[i] = experimentalV4PhoneBuild69SeedResult{bank: bank, evals: n, elapsed: time.Since(seedStarted)}
			}
		}()
	}
	for i := range seeds {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	profile.SeedWallElapsed = time.Since(started)

	bank := make([]experimentalV4PhoneHypothesis, 0, 1024)
	durations := make([]time.Duration, 0, len(results))
	for i := range results {
		r := results[i]
		evals += r.evals
		bank = append(bank, r.bank...)
		profile.SeedWorkerElapsed += r.elapsed
		durations = append(durations, r.elapsed)
		if i == 0 || r.elapsed < profile.SeedMinElapsed {
			profile.SeedMinElapsed = r.elapsed
		}
		if r.elapsed > profile.SeedMaxElapsed {
			profile.SeedMaxElapsed = r.elapsed
		}
		if i == 0 || r.evals < profile.SeedMinEvaluations {
			profile.SeedMinEvaluations = r.evals
		}
		if r.evals > profile.SeedMaxEvaluations {
			profile.SeedMaxEvaluations = r.evals
		}
		if i == 0 || len(r.bank) < profile.SeedMinBank {
			profile.SeedMinBank = len(r.bank)
		}
		if len(r.bank) > profile.SeedMaxBank {
			profile.SeedMaxBank = len(r.bank)
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	if len(durations) > 0 {
		profile.SeedMedianElapsed = durations[len(durations)/2]
	}
	return bank, evals, len(seeds), workers, profile
}

type experimentalV4PhoneBuild69RecoveryTelemetry struct {
	experimentalV4PhoneBuild68RecoveryTelemetry
	GeometryPlanePrepElapsed   time.Duration
	GeometryFreezeElapsed      time.Duration
	GeometrySeedWallElapsed    time.Duration
	GeometrySeedWorkerElapsed  time.Duration
	GeometrySeedMinElapsed     time.Duration
	GeometrySeedMedianElapsed  time.Duration
	GeometrySeedMaxElapsed     time.Duration
	GeometrySeedMinEvaluations int
	GeometrySeedMaxEvaluations int
	GeometrySeedMinBank        int
	GeometrySeedMaxBank        int
}

// experimentalV4PhoneBuild69Recover preserves the qualified Build68 path and
// adds observability inside geometry generation only. Qualification still uses
// Build68 plane reuse and protected-data decode still uses Build66 ordered
// batches. No geometry score, search bound, bank order or selection rule changes.
func experimentalV4PhoneBuild69Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild69RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild69RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, seedWorkers, geometryProfile := experimentalV4PhoneBuild69BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = seedWorkers
	telemetry.BankCandidates = len(bank)
	telemetry.GeometryPlanePrepElapsed = geometryProfile.PlanePrepElapsed
	telemetry.GeometryFreezeElapsed = geometryProfile.FreezeElapsed
	telemetry.GeometrySeedWallElapsed = geometryProfile.SeedWallElapsed
	telemetry.GeometrySeedWorkerElapsed = geometryProfile.SeedWorkerElapsed
	telemetry.GeometrySeedMinElapsed = geometryProfile.SeedMinElapsed
	telemetry.GeometrySeedMedianElapsed = geometryProfile.SeedMedianElapsed
	telemetry.GeometrySeedMaxElapsed = geometryProfile.SeedMaxElapsed
	telemetry.GeometrySeedMinEvaluations = geometryProfile.SeedMinEvaluations
	telemetry.GeometrySeedMaxEvaluations = geometryProfile.SeedMaxEvaluations
	telemetry.GeometrySeedMinBank = geometryProfile.SeedMinBank
	telemetry.GeometrySeedMaxBank = geometryProfile.SeedMaxBank
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build69 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build69 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build69 recovery authentication failed")
}

func experimentalV4PhoneBuild69ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild69RecoveryTelemetry) {
	experimentalV4PhoneBuild68ApplyTelemetry(public, recovery.experimentalV4PhoneBuild68RecoveryTelemetry)
	public.Build69Attempted = recovery.Attempted
	public.Build69GeometryPlanePrepMs = recovery.GeometryPlanePrepElapsed.Milliseconds()
	public.Build69GeometryFreezeMs = recovery.GeometryFreezeElapsed.Milliseconds()
	public.Build69GeometrySeedWallMs = recovery.GeometrySeedWallElapsed.Milliseconds()
	public.Build69GeometrySeedWorkerMs = recovery.GeometrySeedWorkerElapsed.Milliseconds()
	public.Build69GeometrySeedMinMs = recovery.GeometrySeedMinElapsed.Milliseconds()
	public.Build69GeometrySeedMedianMs = recovery.GeometrySeedMedianElapsed.Milliseconds()
	public.Build69GeometrySeedMaxMs = recovery.GeometrySeedMaxElapsed.Milliseconds()
	public.Build69GeometrySeedMinEvals = recovery.GeometrySeedMinEvaluations
	public.Build69GeometrySeedMaxEvals = recovery.GeometrySeedMaxEvaluations
	public.Build69GeometrySeedMinBank = recovery.GeometrySeedMinBank
	public.Build69GeometrySeedMaxBank = recovery.GeometrySeedMaxBank
}
