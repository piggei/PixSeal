package watermark

import (
	"errors"
	"image"
	"runtime"
	"sync"
	"time"
)

// Build90 is a scheduling-only performance candidate over the qualified Build84
// smartphone baseline. Build89 full-pipeline profiling showed that qualification
// dominates B/mild wall time (~61.2 s of ~149.5 s deep recovery) while the
// complete frozen bank is already available before qualification begins.
//
// Build90 evaluates each already-frozen candidate independently in a bounded
// worker pool, stores the result at the candidate's original bank index, and
// commits results strictly in bank order after every task finishes. It does not
// change geometry generation, qualification arithmetic/gates, candidate order,
// protected-data decode ordering, ECC/Hamming, whitening/HMAC domains, strength,
// or wire format. The shared pixelPlane and public pilot candidate are read-only.

type experimentalV4PhoneBuild90QualificationResult struct {
	hypothesis experimentalV4PhoneHypothesis
	evals      int
	ok         bool
}

type experimentalV4PhoneBuild90RecoveryTelemetry struct {
	experimentalV4PhoneBuild84RecoveryTelemetry
	QualificationWorkers int
	QualificationTasks   int
}

func experimentalV4PhoneBuild90QualificationWorkerCount(candidates int) int {
	if candidates <= 0 {
		return 0
	}
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > candidates {
		workers = candidates
	}
	return workers
}

// experimentalV4PhoneBuild90EvaluateQualification executes every frozen bank
// entry exactly once. The returned slice is indexed identically to bank; worker
// completion order is deliberately not exposed as semantic order.
func experimentalV4PhoneBuild90EvaluateQualification(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, bank []experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild90QualificationResult, int) {
	workers := experimentalV4PhoneBuild90QualificationWorkerCount(len(bank))
	if workers == 0 {
		return nil, 0
	}

	results := make([]experimentalV4PhoneBuild90QualificationResult, len(bank))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				q, n, ok := experimentalV4PhoneBuild68Qualify(plane, pilot, cw, ch, bank[i], heldout)
				results[i] = experimentalV4PhoneBuild90QualificationResult{hypothesis: q, evals: n, ok: ok}
			}
		}()
	}
	for i := range bank {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results, workers
}

func experimentalV4PhoneBuild90CommitQualification(results []experimentalV4PhoneBuild90QualificationResult) ([]experimentalV4PhoneHypothesis, int) {
	qualified := make([]experimentalV4PhoneHypothesis, 0, len(results))
	evals := 0
	for i := range results {
		evals += results[i].evals
		if results[i].ok {
			qualified = append(qualified, results[i].hypothesis)
		}
	}
	return qualified, evals
}

// experimentalV4PhoneBuild90QualifyBank preserves the exact serial Build84
// semantic order. Physical qualification work may complete in any order, but
// every result is stored by original bank index and the authoritative eval count
// and qualified slice are reconstructed only by an ordered commit pass.
func experimentalV4PhoneBuild90QualifyBank(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, bank []experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneHypothesis, int, int) {
	results, workers := experimentalV4PhoneBuild90EvaluateQualification(plane, pilot, cw, ch, bank, heldout)
	qualified, evals := experimentalV4PhoneBuild90CommitQualification(results)
	return qualified, evals, workers
}

// experimentalV4PhoneBuild90Recover is Build84 with one scheduling change: the
// already-frozen qualification bank is evaluated by the ordered worker pool
// above. Geometry and protected-data decode continue to use the qualified Build84
// implementation and telemetry semantics.
func experimentalV4PhoneBuild90Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild90RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild90RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, gp, fp := experimentalV4PhoneBuild84BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.Build84GeometryProfile = gp
	telemetry.FreezeProfile = fp
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = prefixWorkers
	telemetry.BankCandidates = len(bank)

	// Preserve the exact Build84/Build76-compatible public geometry view.
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build90 recovery bank empty")
	}

	planeStarted := time.Now()
	plane := newPixelPlane(work)
	telemetry.PlanePrepElapsed = time.Since(planeStarted)
	pilot := experimentalV4Prototype2Candidate()

	telemetry.QualificationTasks = len(bank)
	qualificationStarted := time.Now()
	qualified, qualificationEvals, qualificationWorkers := experimentalV4PhoneBuild90QualifyBank(plane, pilot, cw, ch, bank, 0)
	telemetry.QualificationElapsed = time.Since(qualificationStarted)
	telemetry.QualificationEvaluations = qualificationEvals
	telemetry.QualificationWorkers = qualificationWorkers
	telemetry.QualifiedCandidates = len(qualified)

	telemetry.DecodeWorkers = experimentalV4PhoneBuild66DecodeWorkerCount(len(qualified))
	if len(qualified) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build90 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build90 recovery authentication failed")
}

func experimentalV4PhoneBuild90ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild90RecoveryTelemetry) {
	experimentalV4PhoneBuild84ApplyTelemetry(public, recovery.experimentalV4PhoneBuild84RecoveryTelemetry)
	public.Build90Attempted = recovery.Attempted
	public.Build90QualificationWorkers = recovery.QualificationWorkers
	public.Build90QualificationTasks = recovery.QualificationTasks
	public.Build90QualificationEvaluations = recovery.QualificationEvaluations
}
