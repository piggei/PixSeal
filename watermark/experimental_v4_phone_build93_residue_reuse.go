package watermark

import (
	"errors"
	"image"
	"math"
	"sync"
	"time"
)

// Build93 is the current qualified protected-data decode successor selected
// from the Build92 post-promotion profile. Build92 showed that B/mild decode
// is now the second-largest wall-time stage and that almost all decode CPU is
// spent re-reading the same protected projective grid once per profile.
//
// Build93 changes only that implementation detail. For each already-qualified
// frozen homography, the protected grid is sampled once into the exact residue
// sums/counts historically reconstructed by experimentalV4ReadProtectedProjectiveMargins.
// Robust, Balanced and Capacity margins are then derived from those immutable
// residues in the historical profile order. Geometry, qualification, candidate
// order, list-frame order and first-logical-HMAC semantics remain unchanged.

const experimentalV4PhoneBuild93Residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks

type experimentalV4PhoneBuild93ProtectedResidues struct {
	sums   [experimentalV4PhoneBuild93Residues]float64
	counts [experimentalV4PhoneBuild93Residues]int
}

// experimentalV4PhoneBuild93SampleProtectedResidues is the sampling half of
// experimentalV4ReadProtectedProjectiveMargins, copied without arithmetic or
// traversal changes. It is deliberately profile-independent: codedBits is not
// consulted while reading the projective grid.
func experimentalV4PhoneBuild93SampleProtectedResidues(plane *pixelPlane, canonicalWidth, canonicalHeight int, h homography) (experimentalV4PhoneBuild93ProtectedResidues, bool) {
	var residues experimentalV4PhoneBuild93ProtectedResidues
	if plane == nil {
		return residues, false
	}
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	if blocksWide < experimentalV4TileWidthBlocks || blocksHigh < experimentalV4TileHeightBlocks {
		return residues, false
	}
	for by := 0; by < blocksHigh; by++ {
		for bx := 0; bx < blocksWide; bx++ {
			v, ok := readProjectiveBlockValue(plane, h, bx*blockSize, by*blockSize, blockSize)
			if !ok {
				continue
			}
			residue := (by%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + bx%experimentalV4TileWidthBlocks
			residues.sums[residue] += v
			residues.counts[residue]++
		}
	}
	return residues, true
}

// experimentalV4PhoneBuild93MarginsFromResidues is the profile-dependent half
// of experimentalV4ReadProtectedProjectiveMargins. Loop order, floating-point
// operations and success condition intentionally match the historical reader.
func experimentalV4PhoneBuild93MarginsFromResidues(residues *experimentalV4PhoneBuild93ProtectedResidues, codedBits int) ([]float64, float64, bool) {
	if codedBits <= 0 {
		return nil, 0, false
	}
	bitSums := make([]float64, codedBits)
	bitObservations := make([]int, codedBits)
	for ordinal, pos := range experimentalV4LockedDataPositions() {
		if residues.counts[pos] == 0 {
			continue
		}
		idx := experimentalV4CodeIndex(ordinal, codedBits)
		bitSums[idx] += residues.sums[pos]
		bitObservations[idx] += residues.counts[pos]
	}
	margins := make([]float64, codedBits)
	confidence := 0.0
	observedBits := 0
	for i, sum := range bitSums {
		if bitObservations[i] == 0 {
			continue
		}
		observedBits++
		margins[i] = sum / float64(bitObservations[i])
		confidence += math.Abs(margins[i])
	}
	if observedBits != codedBits {
		return nil, 0, false
	}
	return margins, confidence / float64(codedBits), true
}

type experimentalV4PhoneBuild93DecodeResult struct {
	experimentalV4PhoneBuild67DecodeResult
	protectedSamplingPasses     int
	protectedProfileDerivations int
}

// experimentalV4PhoneBuild93DecodeSingle preserves Build67 semantic ordering
// and list/HMAC behavior. The only change is one protected sampling pass per
// physical candidate instead of one pass per attempted profile.
func experimentalV4PhoneBuild93DecodeSingle(plane *pixelPlane, h experimentalV4PhoneHypothesis, key []byte, cw, ch int) experimentalV4PhoneBuild93DecodeResult {
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

	result := experimentalV4PhoneBuild93DecodeResult{}
	result.info = info

	started := time.Now()
	residues, ok := experimentalV4PhoneBuild93SampleProtectedResidues(plane, cw, ch, h.h)
	result.samplingElapsed += time.Since(started)
	if !ok {
		result.err = errors.New("experimental v4 Build93 single-candidate authentication failed")
		return result
	}
	result.protectedSamplingPasses = 1

	for _, spec := range v3Profiles {
		started = time.Now()
		margins, confidence, ok := experimentalV4PhoneBuild93MarginsFromResidues(&residues, spec.codedBits)
		result.samplingElapsed += time.Since(started)
		result.protectedProfileDerivations++
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

	result.err = errors.New("experimental v4 Build93 single-candidate authentication failed")
	return result
}

type experimentalV4PhoneBuild93RecoveryTelemetry struct {
	experimentalV4PhoneBuild90RecoveryTelemetry
	ProtectedSamplingPasses     int
	ProtectedProfileDerivations int
}

func experimentalV4PhoneBuild93AccumulatePhysical(telemetry *experimentalV4PhoneBuild93RecoveryTelemetry, results []experimentalV4PhoneBuild93DecodeResult) {
	for _, result := range results {
		telemetry.PhysicalDecodeCandidates++
		telemetry.PhysicalProfilesTried += result.profiles
		telemetry.PhysicalListFramesTried += result.frames
		telemetry.SamplingWorkerElapsed += result.samplingElapsed
		telemetry.ListWorkerElapsed += result.listElapsed
		telemetry.ProtectedSamplingPasses += result.protectedSamplingPasses
		telemetry.ProtectedProfileDerivations += result.protectedProfileDerivations
	}
}

// experimentalV4PhoneBuild93Recover preserves the qualified Build90 barriers
// and scheduling. Only the per-candidate protected-data decoder is replaced by
// the exact residue-reuse implementation above. Build93 is now the qualified
// smartphone runtime; Build90 remains its previous-qualified predecessor.
func experimentalV4PhoneBuild93Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild93RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild93RecoveryTelemetry{}
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build93 recovery bank empty")
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
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build93 recovery authentication failed")
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

		experimentalV4PhoneBuild93AccumulatePhysical(&telemetry, results)
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
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build93 recovery authentication failed")
}

// experimentalV4PhoneBuild93ApplyTelemetry preserves the complete qualified
// Build90/Build84 public telemetry view and adds only candidate-specific reuse
// counters used by the deferred physical gate.
func experimentalV4PhoneBuild93ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild93RecoveryTelemetry) {
	experimentalV4PhoneBuild90ApplyTelemetry(public, recovery.experimentalV4PhoneBuild90RecoveryTelemetry)
	public.Build93Attempted = recovery.Attempted
	public.Build93ProtectedSamplingPasses = recovery.ProtectedSamplingPasses
	public.Build93ProtectedProfileDerivations = recovery.ProtectedProfileDerivations
}
