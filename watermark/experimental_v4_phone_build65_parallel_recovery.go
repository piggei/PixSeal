package watermark

import (
	"errors"
	"image"
	"runtime"
	"sync"
)

// experimentalV4PhoneBuild65RecoveryTelemetry extends the already-qualified
// Build64 recovery telemetry with scheduling information only. Build65 changes
// wall-clock scheduling, not geometry, ranking, qualification or decode order.
type experimentalV4PhoneBuild65RecoveryTelemetry struct {
	experimentalV4PhoneBuild64RecoveryTelemetry
	Workers int
}

type experimentalV4PhoneBuild65SeedResult struct {
	bank  []experimentalV4PhoneHypothesis
	evals int
}

// experimentalV4PhoneBuild65BlindBank is exactly the Build64 blind search, but
// schedules independent Build48 seeds concurrently. Results are stored by seed
// index and concatenated only after all workers complete, so the final bank and
// deterministic ordering are identical to serial Build64.
func experimentalV4PhoneBuild65BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int) {
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0, 0
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, _, evals := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild63SeedsPerPair)
	if len(seeds) == 0 {
		return nil, evals, 0, 0
	}

	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(seeds) {
		workers = len(seeds)
	}
	results := make([]experimentalV4PhoneBuild65SeedResult, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				bank, n := experimentalV4PhoneBuild64SeedBank(plane, pilot, cw, ch, anchor, seeds[i])
				results[i] = experimentalV4PhoneBuild65SeedResult{bank: bank, evals: n}
			}
		}()
	}
	for i := range seeds {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	bank := make([]experimentalV4PhoneHypothesis, 0, 1024)
	for i := range results {
		evals += results[i].evals
		bank = append(bank, results[i].bank...)
	}
	return bank, evals, len(seeds), workers
}

// experimentalV4PhoneBuild65Recover preserves Build64's freeze barrier and
// decode order. The only difference is concurrent execution of independent
// proposal-only seed branches before the bank is reassembled in original order.
func experimentalV4PhoneBuild65Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild65RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild65RecoveryTelemetry{}
	telemetry.Attempted = true
	bank, evals, seeds, workers := experimentalV4PhoneBuild65BlindBank(work, boundary, cw, ch)
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.Workers = workers
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build65 recovery bank empty")
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
	for _, q := range qualified {
		telemetry.DecodeCandidatesTried++
		payload, info, profiles, frames, confidence, err := experimentalV4PhoneBuild64DecodeSingle(plane, q, key, cw, ch)
		telemetry.ProfilesTried += profiles
		telemetry.ListFramesTried += frames
		if confidence > telemetry.MaxDataConfidence {
			telemetry.MaxDataConfidence = confidence
		}
		if err == nil {
			telemetry.Authenticated = true
			telemetry.Profile = info.Profile
			return payload, info, telemetry, nil
		}
	}
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build65 recovery authentication failed")
}

func experimentalV4PhoneBuild65ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild65RecoveryTelemetry) {
	experimentalV4PhoneBuild64ApplyTelemetry(public, recovery.experimentalV4PhoneBuild64RecoveryTelemetry)
	public.Build65Attempted = recovery.Attempted
	public.Build65Workers = recovery.Workers
}
