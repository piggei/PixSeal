package watermark

import (
	"errors"
	"image"
	"sync"
	"time"
)

type experimentalV4PhoneBuild95RecoveryTelemetry struct {
	experimentalV4PhoneBuild93RecoveryTelemetry
	Build95GeometryProfile experimentalV4PhoneBuild95GeometryTelemetry
}

func experimentalV4PhoneBuild95Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild95RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild95RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, prefixWorkers, gp, fp := experimentalV4PhoneBuild95BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.Build84GeometryProfile = gp.experimentalV4PhoneBuild84GeometryTelemetry
	telemetry.Build95GeometryProfile = gp
	telemetry.FreezeProfile = fp
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = prefixWorkers
	telemetry.BankCandidates = len(bank)

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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build95 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build95 recovery authentication failed")
	}

	workers := telemetry.DecodeWorkers
	decodeStarted := time.Now()
	for start := 0; start < len(qualified); start += workers {
		end := start + workers
		if end > len(qualified) {
			end = len(qualified)
		}
		results := make([]experimentalV4PhoneBuild93DecodeResult, end-start)
		var wg sync.WaitGroup
		wg.Add(end - start)
		for i := start; i < end; i++ {
			i := i
			go func() {
				defer wg.Done()
				results[i-start] = experimentalV4PhoneBuild93DecodeSingle(plane, qualified[i], key, cw, ch)
			}()
		}
		wg.Wait()

		experimentalV4PhoneBuild93AccumulatePhysical(&telemetry.experimentalV4PhoneBuild93RecoveryTelemetry, results)
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build95 recovery authentication failed")
}

func experimentalV4PhoneBuild95ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild95RecoveryTelemetry) {
	experimentalV4PhoneBuild93ApplyTelemetry(public, recovery.experimentalV4PhoneBuild93RecoveryTelemetry)
	p := recovery.Build95GeometryProfile
	public.Build95Attempted = recovery.Attempted
	public.Build95LuminancePlaneUsed = p.LuminancePlaneUsed
	public.Build95LuminancePlaneFallback = p.LuminancePlaneFallback
	public.Build95LuminancePlaneBytes = p.LuminancePlaneBytes
	public.Build95LuminancePlanePrepMs = p.LuminancePlanePrepElapsed.Milliseconds()
}
