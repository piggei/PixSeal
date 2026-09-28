package watermark

import (
	"errors"
	"fmt"
	"image"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// Build70 is observability-only. It reproduces the exact Build68/69 geometry
// search, but records the genealogy of the most expensive seed so the next
// performance experiment can target the stage that actually explodes.

const experimentalV4PhoneBuild70EvalStageCount = 17
const experimentalV4PhoneBuild70StateStageCount = 13

type experimentalV4PhoneBuild70SeedStageProfile struct {
	Evals  [experimentalV4PhoneBuild70EvalStageCount]int
	States [experimentalV4PhoneBuild70StateStageCount]int
}

// Evals order:
// baseline,roots,single1,pair1,cont1,sib1,single2,pair2,cont2,sib2,
// single3,pair3,cont3,sib3,single4,pair4,cont4.
// States order:
// roots,pair1,cont1,sib1,pair2,cont2,sib2,pair3,cont3,sib3,pair4,cont4,bank.
func experimentalV4PhoneBuild70SeedBankProfile(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, seed experimentalV4PhoneBuild48Seed) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild70SeedStageProfile) {
	var profile experimentalV4PhoneBuild70SeedStageProfile
	bank := make([]experimentalV4PhoneHypothesis, 0, 64)
	evals := 0

	baseline, n := experimentalV4PhoneBuild41Refine(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
	profile.Evals[0] += n
	evals += n
	roots, n := experimentalV4PhoneBuild53TwoPxRoots(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
	profile.Evals[1] += n
	profile.States[0] += len(roots)
	evals += n
	if baseline.h.h[8] == 0 || len(roots) == 0 {
		profile.States[12] = len(bank)
		return bank, evals, profile
	}
	for _, root := range roots {
		se, si := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
		profile.Evals[2] += se
		evals += se
		if si != 0 {
			continue
		}
		pairs, pe, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
		profile.Evals[3] += pe
		profile.States[1] += len(pairs)
		evals += pe
		for _, pair := range pairs {
			cont, ce := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
			profile.Evals[4] += ce
			profile.States[2] += len(cont)
			evals += ce
			for _, cs := range cont {
				sibs, ne := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, cs.hyp, 0)
				profile.Evals[5] += ne
				profile.States[3] += len(sibs)
				evals += ne
				for _, sib := range sibs {
					se2, si2 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, sib.hyp, 0)
					profile.Evals[6] += se2
					evals += se2
					if si2 != 0 {
						continue
					}
					pairs2, pe2, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, sib.hyp, 0)
					profile.Evals[7] += pe2
					profile.States[4] += len(pairs2)
					evals += pe2
					for _, pair2 := range pairs2 {
						cont2, ce2 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair2.hyp, 0)
						profile.Evals[8] += ce2
						profile.States[5] += len(cont2)
						evals += ce2
						for _, c2 := range cont2 {
							sibs2, ne2 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c2.hyp, 0)
							profile.Evals[9] += ne2
							profile.States[6] += len(sibs2)
							evals += ne2
							for _, s2 := range sibs2 {
								se3, si3 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s2.hyp, 0)
								profile.Evals[10] += se3
								evals += se3
								if si3 != 0 {
									continue
								}
								pairs3, pe3, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s2.hyp, 0)
								profile.Evals[11] += pe3
								profile.States[7] += len(pairs3)
								evals += pe3
								for _, pair3 := range pairs3 {
									cont3, ce3 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair3.hyp, 0)
									profile.Evals[12] += ce3
									profile.States[8] += len(cont3)
									evals += ce3
									for _, c3 := range cont3 {
										sibs3, ne3 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c3.hyp, 0)
										profile.Evals[13] += ne3
										profile.States[9] += len(sibs3)
										evals += ne3
										for _, s3 := range sibs3 {
											se4, si4 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s3.hyp, 0)
											profile.Evals[14] += se4
											evals += se4
											if si4 != 0 {
												continue
											}
											pairs4, pe4, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s3.hyp, 0)
											profile.Evals[15] += pe4
											profile.States[10] += len(pairs4)
											evals += pe4
											for _, pair4 := range pairs4 {
												cont4, ce4 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
												profile.Evals[16] += ce4
												profile.States[11] += len(cont4)
												evals += ce4
												for _, c4 := range cont4 {
													bank = append(bank, c4.hyp)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	profile.States[12] = len(bank)
	return bank, evals, profile
}

type experimentalV4PhoneBuild70SeedResult struct {
	bank    []experimentalV4PhoneHypothesis
	evals   int
	elapsed time.Duration
	stage   experimentalV4PhoneBuild70SeedStageProfile
}

type experimentalV4PhoneBuild70GeometryTelemetry struct {
	experimentalV4PhoneBuild69GeometryTelemetry
	FreezeEvaluations     int
	SeedEvaluations       int
	MaxSeedIndex          int
	MaxSeedPairRank       int
	MaxSeedRankWithinPair int
	MaxSeedEvaluations    int
	MaxSeedBank           int
	MaxSeedElapsed        time.Duration
	MaxSeedStage          experimentalV4PhoneBuild70SeedStageProfile
}

func experimentalV4PhoneBuild70BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild70GeometryTelemetry) {
	var profile experimentalV4PhoneBuild70GeometryTelemetry
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
	profile.FreezeEvaluations = evals
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
	results := make([]experimentalV4PhoneBuild70SeedResult, len(seeds))
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	started = time.Now()
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				seedStarted := time.Now()
				bank, n, stage := experimentalV4PhoneBuild70SeedBankProfile(plane, pilot, cw, ch, anchor, seeds[i])
				results[i] = experimentalV4PhoneBuild70SeedResult{bank: bank, evals: n, elapsed: time.Since(seedStarted), stage: stage}
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
	maxIndex := 0
	for i := range results {
		r := results[i]
		evals += r.evals
		profile.SeedEvaluations += r.evals
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
		if i == 0 || r.evals > results[maxIndex].evals {
			maxIndex = i
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	if len(durations) > 0 {
		profile.SeedMedianElapsed = durations[len(durations)/2]
	}
	maxSeed := seeds[maxIndex]
	maxResult := results[maxIndex]
	profile.MaxSeedIndex = maxSeed.index
	profile.MaxSeedPairRank = maxSeed.frozen.h.build43PairRank
	profile.MaxSeedRankWithinPair = maxSeed.rankWithinPair
	profile.MaxSeedEvaluations = maxResult.evals
	profile.MaxSeedBank = len(maxResult.bank)
	profile.MaxSeedElapsed = maxResult.elapsed
	profile.MaxSeedStage = maxResult.stage
	return bank, evals, len(seeds), workers, profile
}

type experimentalV4PhoneBuild70RecoveryTelemetry struct {
	experimentalV4PhoneBuild69RecoveryTelemetry
	GeometryFreezeEvaluations     int
	GeometrySeedEvaluations       int
	GeometryMaxSeedIndex          int
	GeometryMaxSeedPairRank       int
	GeometryMaxSeedRankWithinPair int
	GeometryMaxSeedEvaluations    int
	GeometryMaxSeedBank           int
	GeometryMaxSeedElapsed        time.Duration
	GeometryMaxSeedStage          experimentalV4PhoneBuild70SeedStageProfile
}

func experimentalV4PhoneBuild70Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild70RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild70RecoveryTelemetry{}
	telemetry.Attempted = true
	totalStarted := time.Now()
	finish := func() { telemetry.TotalElapsed = time.Since(totalStarted) }

	geometryStarted := time.Now()
	bank, evals, seeds, seedWorkers, geometryProfile := experimentalV4PhoneBuild70BlindBank(work, boundary, cw, ch)
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
	telemetry.GeometryFreezeEvaluations = geometryProfile.FreezeEvaluations
	telemetry.GeometrySeedEvaluations = geometryProfile.SeedEvaluations
	telemetry.GeometryMaxSeedIndex = geometryProfile.MaxSeedIndex
	telemetry.GeometryMaxSeedPairRank = geometryProfile.MaxSeedPairRank
	telemetry.GeometryMaxSeedRankWithinPair = geometryProfile.MaxSeedRankWithinPair
	telemetry.GeometryMaxSeedEvaluations = geometryProfile.MaxSeedEvaluations
	telemetry.GeometryMaxSeedBank = geometryProfile.MaxSeedBank
	telemetry.GeometryMaxSeedElapsed = geometryProfile.MaxSeedElapsed
	telemetry.GeometryMaxSeedStage = geometryProfile.MaxSeedStage
	if len(bank) == 0 {
		finish()
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build70 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build70 recovery authentication failed")
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build70 recovery authentication failed")
}

func experimentalV4PhoneBuild70CSV(values []int) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ",")
}

func experimentalV4PhoneBuild70ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild70RecoveryTelemetry) {
	experimentalV4PhoneBuild69ApplyTelemetry(public, recovery.experimentalV4PhoneBuild69RecoveryTelemetry)
	public.Build70Attempted = recovery.Attempted
	public.Build70GeometryFreezeEvals = recovery.GeometryFreezeEvaluations
	public.Build70GeometrySeedEvals = recovery.GeometrySeedEvaluations
	public.Build70MaxSeedIndex = recovery.GeometryMaxSeedIndex
	public.Build70MaxSeedPairRank = recovery.GeometryMaxSeedPairRank
	public.Build70MaxSeedRankWithinPair = recovery.GeometryMaxSeedRankWithinPair
	public.Build70MaxSeedEvals = recovery.GeometryMaxSeedEvaluations
	public.Build70MaxSeedBank = recovery.GeometryMaxSeedBank
	public.Build70MaxSeedMs = recovery.GeometryMaxSeedElapsed.Milliseconds()
	public.Build70MaxSeedStageEvals = experimentalV4PhoneBuild70CSV(recovery.GeometryMaxSeedStage.Evals[:])
	public.Build70MaxSeedStageStates = experimentalV4PhoneBuild70CSV(recovery.GeometryMaxSeedStage.States[:])
}
