package watermark

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
)

// Build24 keeps Format v4 experimental and non-normative. This file provides
// deterministic pilot-candidate search and qualification helpers only; no
// production encoder/decoder path consults these values.
const (
	experimentalV4Build24MaskSearchSeed   uint64 = 0x50585345414c5634 // "PXSEALV4"
	experimentalV4Build24MaskSearchBudget        = 200000
	experimentalV4Build24SignSearchSeed   uint64 = 0x76347369676e7331 // "v4signs1"
	experimentalV4Build24SignSearchBudget        = 200000
	experimentalV4RandomSubsetSeed        uint64 = 0x763463726f703031 // "v4crop01"
	experimentalV4RandomSubsetCases              = 256
)

type experimentalV4PilotCandidate struct {
	name      string
	positions [experimentalV4PilotCount]int
	signs     [experimentalV4PilotCount]int8
}

// ExperimentalV4CyclicMetrics reports the exhaustive toroidal wrong-origin
// audit for one pilot candidate. Correlations are raw +/-1 symbol sums.
type ExperimentalV4CyclicMetrics struct {
	MaximumMaskOverlap             int
	MeanMaskOverlap                float64
	MaximumWrongSignedCorrelation  int
	RunnerUpWrongSignedCorrelation int
	WorstOverlapDX                 int
	WorstOverlapDY                 int
	WorstCorrelationDX             int
	WorstCorrelationDY             int
	PerfectNonZeroCyclicAliases    int
}

// ExperimentalV4PartialVisibilityMetrics measures idealized structural origin
// separation when only part of the public pilot is visible. These metrics use
// pilot symbols only; image-domain photometric/channel tests are separate.
type ExperimentalV4PartialVisibilityMetrics struct {
	VisiblePilots                   int
	ContiguousCases                 int
	WorstContiguousWrongCorrelation int
	WorstContiguousMargin           int
	MedianContiguousMargin          float64
	RandomCases                     int
	WorstRandomWrongCorrelation     int
	WorstRandomMargin               int
	MedianRandomMargin              float64
	FalseOriginCases                int
	FalseOriginRate                 float64
}

// ExperimentalV4PilotQualification is a deterministic structural report for a
// candidate. CandidateHash identifies the complete geometry/mask/sign sequence.
type ExperimentalV4PilotQualification struct {
	Name          string
	CandidateHash string
	Cyclic        ExperimentalV4CyclicMetrics
	Partial       []ExperimentalV4PartialVisibilityMetrics
}

// ExperimentalV4Build24SearchReport records enough information to reproduce
// the bounded Build24 search. The search has a joint coordinate/sign stage and
// a second sign-refinement stage on the winning mask.
type ExperimentalV4Build24SearchReport struct {
	MaskSearchSeed   uint64
	MaskSearchBudget int
	SignSearchSeed   uint64
	SignSearchBudget int
	Baseline         ExperimentalV4PilotQualification
	StageOneWinner   ExperimentalV4PilotQualification
	StageOneRunnerUp ExperimentalV4PilotQualification
	Winner           ExperimentalV4PilotQualification
	RunnerUp         ExperimentalV4PilotQualification
}

// prototype-2 is the reproducible Build24 search result. It is deliberately
// not a normative Format-v4 constant; later image-domain/physical evidence may
// still reject or replace it.
var experimentalV4Prototype2PilotPositions = [...]int{
	3, 44, 13, 88, 94, 27, 28, 32,
	259, 156, 159, 238, 279, 246, 290, 181,
	333, 376, 418, 312, 391, 431, 399, 439,
	521, 523, 492, 573, 574, 505, 474, 551,
	631, 674, 675, 719, 687, 617, 660, 664,
	781, 859, 790, 757, 871, 841, 880, 811,
	963, 969, 1008, 903, 907, 1022, 990, 920,
	1151, 1118, 1119, 1165, 1092, 1061, 1103, 1070,
}

var experimentalV4Prototype2PilotSigns = [...]int8{
	-1, 1, -1, -1, 1, -1, 1, -1,
	-1, 1, -1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, -1, 1, -1, -1,
	1, -1, -1, -1, -1, 1, -1, 1,
	-1, 1, 1, 1, 1, -1, 1, 1,
	1, -1, -1, -1, -1, -1, 1, 1,
	-1, 1, 1, -1, -1, -1, -1, -1,
	-1, 1, -1, 1, 1, -1, -1, 1,
}

func experimentalV4Prototype1Candidate() experimentalV4PilotCandidate {
	candidate := experimentalV4PilotCandidate{name: experimentalV4PrototypeName}
	copy(candidate.positions[:], experimentalV4PrototypePilotPositions[:])
	copy(candidate.signs[:], experimentalV4PrototypePilotSigns[:])
	return candidate
}

func experimentalV4Prototype2Candidate() experimentalV4PilotCandidate {
	candidate := experimentalV4PilotCandidate{name: "prototype-2-search-p64"}
	copy(candidate.positions[:], experimentalV4Prototype2PilotPositions[:])
	copy(candidate.signs[:], experimentalV4Prototype2PilotSigns[:])
	return candidate
}

func experimentalV4PilotCandidateHash(candidate experimentalV4PilotCandidate) string {
	h := sha256.New()
	var word [2]byte
	binary.BigEndian.PutUint16(word[:], uint16(experimentalV4TileWidthBlocks))
	h.Write(word[:])
	binary.BigEndian.PutUint16(word[:], uint16(experimentalV4TileHeightBlocks))
	h.Write(word[:])
	for i, position := range candidate.positions {
		binary.BigEndian.PutUint16(word[:], uint16(position))
		h.Write(word[:])
		h.Write([]byte{byte(candidate.signs[i] + 1)})
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func experimentalV4CyclicMetrics(candidate experimentalV4PilotCandidate) ExperimentalV4CyclicMetrics {
	const tilePositions = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var overlap [tilePositions]int
	var correlation [tilePositions]int
	for i, position := range candidate.positions {
		xi := position % experimentalV4TileWidthBlocks
		yi := position / experimentalV4TileWidthBlocks
		for j, other := range candidate.positions {
			if i == j {
				continue
			}
			xj := other % experimentalV4TileWidthBlocks
			yj := other / experimentalV4TileWidthBlocks
			dx := xi - xj
			if dx < 0 {
				dx += experimentalV4TileWidthBlocks
			}
			dy := yi - yj
			if dy < 0 {
				dy += experimentalV4TileHeightBlocks
			}
			index := dy*experimentalV4TileWidthBlocks + dx
			overlap[index]++
			correlation[index] += int(candidate.signs[i] * candidate.signs[j])
		}
	}

	metrics := ExperimentalV4CyclicMetrics{}
	correlations := make([]int, 0, tilePositions-1)
	totalOverlap := 0
	for index := 1; index < tilePositions; index++ {
		dx := index % experimentalV4TileWidthBlocks
		dy := index / experimentalV4TileWidthBlocks
		if overlap[index] > metrics.MaximumMaskOverlap {
			metrics.MaximumMaskOverlap = overlap[index]
			metrics.WorstOverlapDX = dx
			metrics.WorstOverlapDY = dy
		}
		absCorrelation := correlation[index]
		if absCorrelation < 0 {
			absCorrelation = -absCorrelation
		}
		if absCorrelation > metrics.MaximumWrongSignedCorrelation {
			metrics.MaximumWrongSignedCorrelation = absCorrelation
			metrics.WorstCorrelationDX = dx
			metrics.WorstCorrelationDY = dy
		}
		correlations = append(correlations, absCorrelation)
		totalOverlap += overlap[index]
		if overlap[index] == experimentalV4PilotCount && absCorrelation == experimentalV4PilotCount {
			metrics.PerfectNonZeroCyclicAliases++
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(correlations)))
	if len(correlations) > 1 {
		metrics.RunnerUpWrongSignedCorrelation = correlations[1]
	}
	metrics.MeanMaskOverlap = float64(totalOverlap) / float64(tilePositions-1)
	return metrics
}

func experimentalV4PilotQualification(candidate experimentalV4PilotCandidate) ExperimentalV4PilotQualification {
	return ExperimentalV4PilotQualification{
		Name:          candidate.name,
		CandidateHash: experimentalV4PilotCandidateHash(candidate),
		Cyclic:        experimentalV4CyclicMetrics(candidate),
		Partial:       experimentalV4PartialVisibility(candidate),
	}
}

func experimentalV4PartialVisibility(candidate experimentalV4PilotCandidate) []ExperimentalV4PartialVisibilityMetrics {
	levels := []int{64, 48, 32, 24, 16}
	out := make([]ExperimentalV4PartialVisibilityMetrics, 0, len(levels))
	for _, visible := range levels {
		contiguousMargins := make([]int, 0, 64)
		worstContiguousWrong := 0
		for _, shape := range experimentalV4ContiguousShapes(visible) {
			width, height := shape[0], shape[1]
			for y0 := 0; y0 <= 8-height; y0++ {
				for x0 := 0; x0 <= 8-width; x0++ {
					indices := make([]int, 0, visible)
					for gy := y0; gy < y0+height; gy++ {
						for gx := x0; gx < x0+width; gx++ {
							indices = append(indices, gy*8+gx)
						}
					}
					wrong := experimentalV4SubsetWorstWrongCorrelation(candidate, indices)
					if wrong > worstContiguousWrong {
						worstContiguousWrong = wrong
					}
					contiguousMargins = append(contiguousMargins, visible-wrong)
				}
			}
		}

		rng := experimentalV4RNG{state: experimentalV4RandomSubsetSeed ^ uint64(visible)*0x9e3779b97f4a7c15}
		randomMargins := make([]int, 0, experimentalV4RandomSubsetCases)
		worstRandomWrong := 0
		falseOriginCases := 0
		for caseIndex := 0; caseIndex < experimentalV4RandomSubsetCases; caseIndex++ {
			indices := experimentalV4RandomSubsetIndices(&rng, visible)
			wrong := experimentalV4SubsetWorstWrongCorrelation(candidate, indices)
			if wrong > worstRandomWrong {
				worstRandomWrong = wrong
			}
			margin := visible - wrong
			if margin <= 0 {
				falseOriginCases++
			}
			randomMargins = append(randomMargins, margin)
		}

		metric := ExperimentalV4PartialVisibilityMetrics{
			VisiblePilots:                   visible,
			ContiguousCases:                 len(contiguousMargins),
			WorstContiguousWrongCorrelation: worstContiguousWrong,
			WorstContiguousMargin:           visible - worstContiguousWrong,
			MedianContiguousMargin:          experimentalV4MedianInt(contiguousMargins),
			RandomCases:                     len(randomMargins),
			WorstRandomWrongCorrelation:     worstRandomWrong,
			WorstRandomMargin:               visible - worstRandomWrong,
			MedianRandomMargin:              experimentalV4MedianInt(randomMargins),
			FalseOriginCases:                falseOriginCases,
			FalseOriginRate:                 float64(falseOriginCases) / float64(len(randomMargins)),
		}
		out = append(out, metric)
	}
	return out
}

func experimentalV4ContiguousShapes(visible int) [][2]int {
	switch visible {
	case 64:
		return [][2]int{{8, 8}}
	case 48:
		return [][2]int{{6, 8}, {8, 6}}
	case 32:
		return [][2]int{{4, 8}, {8, 4}}
	case 24:
		return [][2]int{{3, 8}, {8, 3}, {4, 6}, {6, 4}}
	case 16:
		return [][2]int{{2, 8}, {8, 2}, {4, 4}}
	default:
		return nil
	}
}

func experimentalV4SubsetWorstWrongCorrelation(candidate experimentalV4PilotCandidate, indices []int) int {
	const tilePositions = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var correlation [tilePositions]int
	for _, i := range indices {
		position := candidate.positions[i]
		xi := position % experimentalV4TileWidthBlocks
		yi := position / experimentalV4TileWidthBlocks
		for j, other := range candidate.positions {
			if i == j {
				continue
			}
			xj := other % experimentalV4TileWidthBlocks
			yj := other / experimentalV4TileWidthBlocks
			dx := xi - xj
			if dx < 0 {
				dx += experimentalV4TileWidthBlocks
			}
			dy := yi - yj
			if dy < 0 {
				dy += experimentalV4TileHeightBlocks
			}
			index := dy*experimentalV4TileWidthBlocks + dx
			correlation[index] += int(candidate.signs[i] * candidate.signs[j])
		}
	}
	worst := 0
	for index := 1; index < tilePositions; index++ {
		value := correlation[index]
		if value < 0 {
			value = -value
		}
		if value > worst {
			worst = value
		}
	}
	return worst
}

func experimentalV4MedianInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	copyValues := append([]int(nil), values...)
	sort.Ints(copyValues)
	middle := len(copyValues) / 2
	if len(copyValues)%2 == 1 {
		return float64(copyValues[middle])
	}
	return float64(copyValues[middle-1]+copyValues[middle]) / 2
}

type experimentalV4RNG struct {
	state uint64
}

func (rng *experimentalV4RNG) next() uint64 {
	x := rng.state
	if x == 0 {
		x = 0x9e3779b97f4a7c15
	}
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	rng.state = x
	return x
}

func (rng *experimentalV4RNG) n(limit int) int {
	return int(rng.next() % uint64(limit))
}

func experimentalV4RandomBalancedSigns(rng *experimentalV4RNG) [experimentalV4PilotCount]int8 {
	var signs [experimentalV4PilotCount]int8
	for index := range signs {
		if index < experimentalV4PilotCount/2 {
			signs[index] = 1
		} else {
			signs[index] = -1
		}
	}
	for index := len(signs) - 1; index > 0; index-- {
		other := rng.n(index + 1)
		signs[index], signs[other] = signs[other], signs[index]
	}
	return signs
}

func experimentalV4RandomSubsetIndices(rng *experimentalV4RNG, visible int) []int {
	var indices [experimentalV4PilotCount]int
	for index := range indices {
		indices[index] = index
	}
	for index := len(indices) - 1; index > 0; index-- {
		other := rng.n(index + 1)
		indices[index], indices[other] = indices[other], indices[index]
	}
	out := append([]int(nil), indices[:visible]...)
	sort.Ints(out)
	return out
}

func experimentalV4StratumXBounds(gx int) (int, int) {
	start := 0
	if gx > 0 {
		start = (gx*experimentalV4TileWidthBlocks + 4) / 8
	}
	end := experimentalV4TileWidthBlocks
	if gx < 7 {
		end = ((gx+1)*experimentalV4TileWidthBlocks + 4) / 8
	}
	return start, end
}

func experimentalV4RandomCandidate(rng *experimentalV4RNG) experimentalV4PilotCandidate {
	candidate := experimentalV4PilotCandidate{name: "build24-search-stage1"}
	index := 0
	for gy := 0; gy < 8; gy++ {
		for gx := 0; gx < 8; gx++ {
			x0, x1 := experimentalV4StratumXBounds(gx)
			y0, y1 := gy*4, (gy+1)*4
			x := x0 + rng.n(x1-x0)
			y := y0 + rng.n(y1-y0)
			candidate.positions[index] = y*experimentalV4TileWidthBlocks + x
			index++
		}
	}
	candidate.signs = experimentalV4RandomBalancedSigns(rng)
	return candidate
}

type experimentalV4SearchScore struct {
	maxOverlap int
	maxCorr    int
	margins    [5]int
}

func experimentalV4SearchScoreFor(candidate experimentalV4PilotCandidate) experimentalV4SearchScore {
	cyclic := experimentalV4CyclicMetrics(candidate)
	partial := experimentalV4PartialVisibilityContiguousOnly(candidate)
	score := experimentalV4SearchScore{maxOverlap: cyclic.MaximumMaskOverlap, maxCorr: cyclic.MaximumWrongSignedCorrelation}
	for index := range partial {
		score.margins[index] = partial[index].WorstContiguousMargin
	}
	return score
}

func experimentalV4PartialVisibilityContiguousOnly(candidate experimentalV4PilotCandidate) [5]ExperimentalV4PartialVisibilityMetrics {
	levels := [5]int{64, 48, 32, 24, 16}
	var out [5]ExperimentalV4PartialVisibilityMetrics
	for levelIndex, visible := range levels {
		worstWrong := 0
		cases := 0
		margins := make([]int, 0, 64)
		for _, shape := range experimentalV4ContiguousShapes(visible) {
			width, height := shape[0], shape[1]
			for y0 := 0; y0 <= 8-height; y0++ {
				for x0 := 0; x0 <= 8-width; x0++ {
					indices := make([]int, 0, visible)
					for gy := y0; gy < y0+height; gy++ {
						for gx := x0; gx < x0+width; gx++ {
							indices = append(indices, gy*8+gx)
						}
					}
					wrong := experimentalV4SubsetWorstWrongCorrelation(candidate, indices)
					if wrong > worstWrong {
						worstWrong = wrong
					}
					margins = append(margins, visible-wrong)
					cases++
				}
			}
		}
		out[levelIndex] = ExperimentalV4PartialVisibilityMetrics{
			VisiblePilots:                   visible,
			ContiguousCases:                 cases,
			WorstContiguousWrongCorrelation: worstWrong,
			WorstContiguousMargin:           visible - worstWrong,
			MedianContiguousMargin:          experimentalV4MedianInt(margins),
		}
	}
	return out
}

func experimentalV4SearchScoreBetter(left, right experimentalV4SearchScore) bool {
	if left.maxOverlap != right.maxOverlap {
		return left.maxOverlap < right.maxOverlap
	}
	if left.maxCorr != right.maxCorr {
		return left.maxCorr < right.maxCorr
	}
	// Give the smallest visible-pilot case highest priority after global
	// structural separation, then work upward toward fuller visibility.
	for index := len(left.margins) - 1; index >= 1; index-- {
		if left.margins[index] != right.margins[index] {
			return left.margins[index] > right.margins[index]
		}
	}
	return false
}

func experimentalV4Build24Search() (experimentalV4PilotCandidate, ExperimentalV4Build24SearchReport) {
	baseline := experimentalV4Prototype1Candidate()
	baselineScore := experimentalV4SearchScoreFor(baseline)

	rng := experimentalV4RNG{state: experimentalV4Build24MaskSearchSeed}
	var best, runnerUp experimentalV4PilotCandidate
	bestScore := experimentalV4SearchScore{maxOverlap: 1 << 30, maxCorr: 1 << 30}
	runnerScore := bestScore
	for attempt := 0; attempt < experimentalV4Build24MaskSearchBudget; attempt++ {
		candidate := experimentalV4RandomCandidate(&rng)
		cyclic := experimentalV4CyclicMetrics(candidate)
		// The Build23 baseline is the admission gate: Build24 does not spend
		// crop scoring on candidates already worse in either global metric.
		if cyclic.MaximumMaskOverlap > baselineScore.maxOverlap || cyclic.MaximumWrongSignedCorrelation > baselineScore.maxCorr {
			continue
		}
		score := experimentalV4SearchScoreFor(candidate)
		if experimentalV4SearchScoreBetter(score, bestScore) {
			runnerUp, runnerScore = best, bestScore
			best, bestScore = candidate, score
		} else if experimentalV4SearchScoreBetter(score, runnerScore) {
			runnerUp, runnerScore = candidate, score
		}
	}

	// Sign-only refinement is permitted only after the joint mask/sign stage.
	// It keeps the winning mask fixed while searching a larger balanced sign
	// family for lower wrong-origin correlation and better crop margins.
	signRNG := experimentalV4RNG{state: experimentalV4Build24SignSearchSeed}
	final := best
	final.name = "prototype-2-search-p64"
	finalScore := experimentalV4SearchScore{maxOverlap: bestScore.maxOverlap, maxCorr: 1 << 30}
	var finalRunner experimentalV4PilotCandidate
	finalRunnerScore := finalScore
	for attempt := 0; attempt < experimentalV4Build24SignSearchBudget; attempt++ {
		candidate := best
		candidate.name = "prototype-2-search-p64"
		candidate.signs = experimentalV4RandomBalancedSigns(&signRNG)
		cyclic := experimentalV4CyclicMetrics(candidate)
		if cyclic.MaximumWrongSignedCorrelation > finalScore.maxCorr {
			continue
		}
		score := experimentalV4SearchScoreFor(candidate)
		if experimentalV4SearchScoreBetter(score, finalScore) {
			finalRunner, finalRunnerScore = final, finalScore
			final, finalScore = candidate, score
		} else if experimentalV4SearchScoreBetter(score, finalRunnerScore) {
			finalRunner, finalRunnerScore = candidate, score
		}
	}

	report := ExperimentalV4Build24SearchReport{
		MaskSearchSeed:   experimentalV4Build24MaskSearchSeed,
		MaskSearchBudget: experimentalV4Build24MaskSearchBudget,
		SignSearchSeed:   experimentalV4Build24SignSearchSeed,
		SignSearchBudget: experimentalV4Build24SignSearchBudget,
		Baseline:         experimentalV4PilotQualification(baseline),
		StageOneWinner:   experimentalV4PilotQualification(best),
		Winner:           experimentalV4PilotQualification(final),
	}
	if runnerUp.name != "" {
		report.StageOneRunnerUp = experimentalV4PilotQualification(runnerUp)
	}
	if finalRunner.name != "" {
		report.RunnerUp = experimentalV4PilotQualification(finalRunner)
	}
	return final, report
}
