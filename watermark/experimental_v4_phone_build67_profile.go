package watermark

import (
	"errors"
	"image"
	"sync"
	"time"
)

// experimentalV4PhoneBuild67RecoveryTelemetry extends the qualified Build66
// scheduling telemetry with observability only. Build67 deliberately leaves
// geometry generation, qualification order, ordered decode batches, logical
// telemetry and first-HMAC semantics unchanged. The additional counters measure
// physical work so a later optimization can target the actual runtime cost.
type experimentalV4PhoneBuild67RecoveryTelemetry struct {
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

type experimentalV4PhoneBuild67DecodeResult struct {
	experimentalV4PhoneBuild66DecodeResult
	samplingElapsed time.Duration
	listElapsed     time.Duration
}

// experimentalV4PhoneBuild67DecodeSingle is the Build64/66 single-candidate
// protected-data decoder with timing observability around the two expensive
// implementation stages. The loops, profile order, margin reader, list decoder
// and HMAC authentication are intentionally identical to Build64.
func experimentalV4PhoneBuild67DecodeSingle(plane *pixelPlane, h experimentalV4PhoneHypothesis, key []byte, cw, ch int) experimentalV4PhoneBuild67DecodeResult {
	candidate := experimentalV4Prototype2Candidate()
	info := ExperimentalV4ExtractInfo{
		Version:       experimentalV4Version,
		PilotName:     candidate.name,
		PilotHash:     experimentalV4PilotCandidateHash(candidate),
		PilotScore:    h.detection.Score,
		PilotMargin:   h.detection.Margin,
		OriginXBlocks: h.detection.OriginXBlocks,
		OriginYBlocks: h.detection.OriginYBlocks,
	}

	result := experimentalV4PhoneBuild67DecodeResult{}
	result.info = info
	for _, spec := range v3Profiles {
		started := time.Now()
		margins, confidence, ok := experimentalV4ReadProtectedProjectiveMargins(plane, cw, ch, h.h, spec.codedBits)
		result.samplingElapsed += time.Since(started)
		if !ok {
			continue
		}

		result.profiles++
		if confidence > result.maxConfidence {
			result.maxConfidence = confidence
		}

		started = time.Now()
		payload, tried, authenticated := experimentalV4PhoneBuild42ListDecode(margins, key, spec)
		result.listElapsed += time.Since(started)
		result.frames += tried
		if authenticated {
			result.payload = payload
			result.info.Profile = spec.profile
			result.info.Confidence = confidence
			result.err = nil
			return result
		}
	}

	result.err = errors.New("experimental v4 Build67 single-candidate authentication failed")
	return result
}

func experimentalV4PhoneBuild67AccumulatePhysical(telemetry *experimentalV4PhoneBuild67RecoveryTelemetry, results []experimentalV4PhoneBuild67DecodeResult) {
	for _, result := range results {
		telemetry.PhysicalDecodeCandidates++
		telemetry.PhysicalProfilesTried += result.profiles
		telemetry.PhysicalListFramesTried += result.frames
		telemetry.SamplingWorkerElapsed += result.samplingElapsed
		telemetry.ListWorkerElapsed += result.listElapsed
	}
}

// experimentalV4PhoneBuild67Recover reproduces Build66 scheduling exactly and
// adds wall-clock/physical-work profiling. Physical counters intentionally
// include speculative candidates already running in the winning ordered batch;
// the embedded Build64/65/66 telemetry remains logical and therefore excludes
// those later speculative results exactly as in the qualified Build66 path.
func experimentalV4PhoneBuild67Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild67RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild67RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() {
		telemetry.TotalElapsed = time.Since(totalStarted)
	}

	geometryStarted := time.Now()
	bank, evals, seeds, seedWorkers := experimentalV4PhoneBuild65BlindBank(work, boundary, cw, ch)
	telemetry.GeometryElapsed = time.Since(geometryStarted)
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = seedWorkers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build67 recovery bank empty")
	}

	planeStarted := time.Now()
	plane := newPixelPlane(work)
	telemetry.PlanePrepElapsed = time.Since(planeStarted)
	pilot := experimentalV4Prototype2Candidate()
	qualified := make([]experimentalV4PhoneHypothesis, 0, len(bank))
	qualificationStarted := time.Now()
	for _, h := range bank {
		q, n, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, h, 0)
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build67 recovery authentication failed")
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

		// Profiling counts physical work, including speculative later members of
		// a winning batch. Semantic consumption remains delegated to the exact
		// Build66 helper and therefore keeps the qualified logical prefix.
		experimentalV4PhoneBuild67AccumulatePhysical(&telemetry, results)
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build67 recovery authentication failed")
}

func experimentalV4PhoneBuild67ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild67RecoveryTelemetry) {
	experimentalV4PhoneBuild66ApplyTelemetry(public, recovery.experimentalV4PhoneBuild66RecoveryTelemetry)
	public.Build67Attempted = recovery.Attempted
	public.Build67TotalMs = recovery.TotalElapsed.Milliseconds()
	public.Build67GeometryMs = recovery.GeometryElapsed.Milliseconds()
	public.Build67PlanePrepMs = recovery.PlanePrepElapsed.Milliseconds()
	public.Build67QualificationMs = recovery.QualificationElapsed.Milliseconds()
	public.Build67DecodeWallMs = recovery.DecodeWallElapsed.Milliseconds()
	public.Build67PhysicalDecodeCandidates = recovery.PhysicalDecodeCandidates
	public.Build67SpeculativeDecodeCandidates = recovery.PhysicalDecodeCandidates - recovery.DecodeCandidatesTried
	if public.Build67SpeculativeDecodeCandidates < 0 {
		public.Build67SpeculativeDecodeCandidates = 0
	}
	public.Build67PhysicalProfiles = recovery.PhysicalProfilesTried
	public.Build67PhysicalListFrames = recovery.PhysicalListFramesTried
	public.Build67SamplingWorkerMs = recovery.SamplingWorkerElapsed.Milliseconds()
	public.Build67ListWorkerMs = recovery.ListWorkerElapsed.Milliseconds()
}
