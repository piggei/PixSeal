package watermark

import (
	"errors"
	"image"
	"sync"
	"time"
)

// experimentalV4PhoneBuild68RecoveryTelemetry keeps Build66 logical telemetry
// authoritative while measuring the Build68 plane-reuse experiment by stage.
// Build68 changes only how full-pilot qualification receives pixel data: the
// already-built post-freeze plane is reused instead of reconstructed once per
// candidate.
type experimentalV4PhoneBuild68RecoveryTelemetry struct {
	experimentalV4PhoneBuild66RecoveryTelemetry

	TotalElapsed         time.Duration
	GeometryElapsed      time.Duration
	PlanePrepElapsed     time.Duration
	QualificationElapsed time.Duration
	DecodeWallElapsed    time.Duration

	PhysicalDecodeCandidates int
	PhysicalProfilesTried    int
	PhysicalListFramesTried  int
	SamplingWorkerElapsed    time.Duration
	ListWorkerElapsed        time.Duration
}

// experimentalV4PhoneBuild68Qualify is intentionally the same qualification
// computation and gate order as Build41/64/66. The only implementation change
// is that full-pilot detection consumes the existing pixelPlane.
func experimentalV4PhoneBuild68Qualify(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, hyp experimentalV4PhoneHypothesis, heldout int) (experimentalV4PhoneHypothesis, int, bool) {
	validation, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hyp.h, heldout, false, false)
	evals := 1
	hyp.validation = validation
	if hyp.proposal < experimentalV4PhoneProposalFloor || validation < experimentalV4PhoneValidationFloor {
		return hyp, evals, false
	}
	det := experimentalV4DetectPilotProjectivePlane(plane, candidate, cw, ch, hyp.h)
	evals++
	hyp.detection = det
	ok := det.Available && det.OriginXBlocks == 0 && det.OriginYBlocks == 0 && det.Score >= experimentalV4PhonePilotScoreFloor && det.Margin >= experimentalV4PhonePilotMarginFloor
	return hyp, evals, ok
}

func experimentalV4PhoneBuild68AccumulatePhysical(telemetry *experimentalV4PhoneBuild68RecoveryTelemetry, results []experimentalV4PhoneBuild67DecodeResult) {
	for _, result := range results {
		telemetry.PhysicalDecodeCandidates++
		telemetry.PhysicalProfilesTried += result.profiles
		telemetry.PhysicalListFramesTried += result.frames
		telemetry.SamplingWorkerElapsed += result.samplingElapsed
		telemetry.ListWorkerElapsed += result.listElapsed
	}
}

// experimentalV4PhoneBuild68Recover preserves the complete Build66 bank freeze,
// qualification order and ordered-parallel protected-data decode. It reuses one
// pixelPlane for every full-pilot qualification candidate and for subsequent
// protected-data sampling. No candidate is pruned, reordered or early-stopped.
func experimentalV4PhoneBuild68Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild68RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild68RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, seedWorkers := experimentalV4PhoneBuild65BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = seedWorkers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build68 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build68 recovery authentication failed")
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

		experimentalV4PhoneBuild68AccumulatePhysical(&telemetry, results)
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build68 recovery authentication failed")
}

func experimentalV4PhoneBuild68ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild68RecoveryTelemetry) {
	experimentalV4PhoneBuild66ApplyTelemetry(public, recovery.experimentalV4PhoneBuild66RecoveryTelemetry)
	public.Build68Attempted = recovery.Attempted
	public.Build68TotalMs = recovery.TotalElapsed.Milliseconds()
	public.Build68GeometryMs = recovery.GeometryElapsed.Milliseconds()
	public.Build68PlanePrepMs = recovery.PlanePrepElapsed.Milliseconds()
	public.Build68QualificationMs = recovery.QualificationElapsed.Milliseconds()
	public.Build68DecodeWallMs = recovery.DecodeWallElapsed.Milliseconds()
	public.Build68PhysicalDecodeCandidates = recovery.PhysicalDecodeCandidates
	public.Build68SpeculativeDecodeCandidates = recovery.PhysicalDecodeCandidates - recovery.DecodeCandidatesTried
	if public.Build68SpeculativeDecodeCandidates < 0 {
		public.Build68SpeculativeDecodeCandidates = 0
	}
	public.Build68PhysicalProfiles = recovery.PhysicalProfilesTried
	public.Build68PhysicalListFrames = recovery.PhysicalListFramesTried
	public.Build68SamplingWorkerMs = recovery.SamplingWorkerElapsed.Milliseconds()
	public.Build68ListWorkerMs = recovery.ListWorkerElapsed.Milliseconds()
}
