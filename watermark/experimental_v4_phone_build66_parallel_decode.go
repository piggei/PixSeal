package watermark

import (
	"errors"
	"image"
	"runtime"
	"sync"
)

// experimentalV4PhoneBuild66RecoveryTelemetry extends Build65 scheduling
// telemetry with the bounded ordered-parallel protected-data decode worker
// count. Build66 changes wall-clock scheduling only: blind geometry, bank order,
// qualification order, logical candidate order, list-frame order and HMAC
// semantics remain Build64-identical.
type experimentalV4PhoneBuild66RecoveryTelemetry struct {
	experimentalV4PhoneBuild65RecoveryTelemetry
	DecodeWorkers int
}

type experimentalV4PhoneBuild66DecodeResult struct {
	payload       []byte
	info          ExperimentalV4ExtractInfo
	profiles      int
	frames        int
	maxConfidence float64
	err           error
}

func experimentalV4PhoneBuild66DecodeWorkerCount(candidates int) int {
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

func experimentalV4PhoneBuild66ConsumeBatch(results []experimentalV4PhoneBuild66DecodeResult, telemetry *experimentalV4PhoneBuild66RecoveryTelemetry) ([]byte, ExperimentalV4ExtractInfo, bool) {
	for _, result := range results {
		telemetry.DecodeCandidatesTried++
		telemetry.ProfilesTried += result.profiles
		telemetry.ListFramesTried += result.frames
		if result.maxConfidence > telemetry.MaxDataConfidence {
			telemetry.MaxDataConfidence = result.maxConfidence
		}
		if result.err == nil {
			telemetry.Authenticated = true
			telemetry.Profile = result.info.Profile
			return result.payload, result.info, true
		}
	}
	return nil, ExperimentalV4ExtractInfo{}, false
}

// experimentalV4PhoneBuild66Recover keeps the complete Build64/65 freeze
// barrier unchanged. Build65's seed-parallel blind bank is generated and
// reassembled in deterministic seed order, then qualification remains serial in
// bank order. Only the already-qualified protected-data candidate loop is
// scheduled in bounded batches. Every batch is fully evaluated concurrently,
// but results are consumed strictly in original candidate order. Consequently
// the first logical HMAC success, payload and public telemetry are exactly the
// same as serial Build64; at most workers-1 later candidates may be physically
// evaluated inside the winning batch, and those speculative results are never
// allowed to affect logical telemetry or selection.
func experimentalV4PhoneBuild66Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild66RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild66RecoveryTelemetry{}
	telemetry.Attempted = true

	bank, evals, seeds, seedWorkers := experimentalV4PhoneBuild65BlindBank(work, boundary, cw, ch)
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = seedWorkers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build66 recovery bank empty")
	}

	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	qualified := make([]experimentalV4PhoneHypothesis, 0, len(bank))
	for _, h := range bank {
		q, n, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, h, 0)
		telemetry.QualificationEvaluations += n
		if ok {
			qualified = append(qualified, q)
		}
	}
	telemetry.QualifiedCandidates = len(qualified)
	telemetry.DecodeWorkers = experimentalV4PhoneBuild66DecodeWorkerCount(len(qualified))
	if len(qualified) == 0 {
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build66 recovery authentication failed")
	}

	workers := telemetry.DecodeWorkers
	for start := 0; start < len(qualified); start += workers {
		end := start + workers
		if end > len(qualified) {
			end = len(qualified)
		}
		results := make([]experimentalV4PhoneBuild66DecodeResult, end-start)
		var wg sync.WaitGroup
		wg.Add(end - start)
		for i := start; i < end; i++ {
			i := i
			go func() {
				defer wg.Done()
				payload, info, profiles, frames, confidence, err := experimentalV4PhoneBuild64DecodeSingle(plane, qualified[i], key, cw, ch)
				results[i-start] = experimentalV4PhoneBuild66DecodeResult{
					payload:       payload,
					info:          info,
					profiles:      profiles,
					frames:        frames,
					maxConfidence: confidence,
					err:           err,
				}
			}()
		}
		wg.Wait()

		// Consume only in the exact serial Build64 order. Results after the first
		// logical success in this batch were speculative physical work and are
		// intentionally excluded from all semantic telemetry.
		if payload, info, ok := experimentalV4PhoneBuild66ConsumeBatch(results, &telemetry); ok {
			return payload, info, telemetry, nil
		}
	}

	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build66 recovery authentication failed")
}

func experimentalV4PhoneBuild66ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild66RecoveryTelemetry) {
	experimentalV4PhoneBuild65ApplyTelemetry(public, recovery.experimentalV4PhoneBuild65RecoveryTelemetry)
	public.Build66Attempted = recovery.Attempted
	public.Build66DecodeWorkers = recovery.DecodeWorkers
}
