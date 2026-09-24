package watermark

import (
	"errors"
	"image"
	"math"
)

const experimentalV4PhoneBuild52SeedsPerPair = experimentalV4PhoneBuild50SeedsPerPair

var experimentalV4PhoneBuild52FineSteps = [...]float64{2, 1}

// ExperimentalV4PhoneBuild52State is one geometry from the Build52 diagnostic
// bank. Geometry/proposal fields are created in the blind phase. Validation,
// full-pilot and HMAC fields are filled only after the complete baseline/fine
// bank for all seeds has been frozen.
type ExperimentalV4PhoneBuild52State struct {
	Index                   int           `json:"index"`
	StepSizePx              float64       `json:"step_size_px"`
	Pass                    int           `json:"pass"`
	Dimension               int           `json:"dimension"`
	Corner                  int           `json:"corner"`
	Axis                    string        `json:"axis,omitempty"`
	DeltaPx                 float64       `json:"delta_px"`
	Proposal                float64       `json:"proposal"`
	Validation              float64       `json:"validation"`
	PilotScore              float64       `json:"pilot_score"`
	PilotMargin             float64       `json:"pilot_margin"`
	Qualified               bool          `json:"qualified"`
	OriginXBlocks           int           `json:"origin_x_blocks"`
	OriginYBlocks           int           `json:"origin_y_blocks"`
	WorkingQuad             [4]ImagePoint `json:"working_quad"`
	SourceQuad              [4]ImagePoint `json:"source_quad"`
	MeanMovementFromSeedPx  float64       `json:"mean_movement_from_seed_px"`
	MaxMovementFromSeedPx   float64       `json:"max_movement_from_seed_px"`
	SingleDecodeAttempted   bool          `json:"single_decode_attempted"`
	SingleProfilesTried     int           `json:"single_profiles_tried,omitempty"`
	SingleListFramesTried   int           `json:"single_list_frames_tried,omitempty"`
	SingleMaxDataConfidence float64       `json:"single_max_data_confidence,omitempty"`
	SingleHMACAuthenticated bool          `json:"single_hmac_authenticated"`
	SingleProfile           Profile       `json:"single_profile,omitempty"`
}

// ExperimentalV4PhoneBuild52Candidate compares the unchanged Build41
// coarse-to-fine endpoint with a fine-first restart bank from the exact same
// blind top-4 seed. The fine bank keeps the seed plus every proposal-improving
// accepted state produced by the fixed 2px -> 1px schedule.
type ExperimentalV4PhoneBuild52Candidate struct {
	Index              int                               `json:"index"`
	SeedIndex          int                               `json:"seed_index"`
	SeedRankWithinPair int                               `json:"seed_rank_within_pair"`
	SourcePair         string                            `json:"source_pair,omitempty"`
	SourcePairRank     int                               `json:"source_pair_rank,omitempty"`
	SourceTier         string                            `json:"source_tier,omitempty"`
	CellRank           int                               `json:"cell_rank"`
	SeedProposal       float64                           `json:"seed_proposal"`
	SeedWorkingQuad    [4]ImagePoint                     `json:"seed_working_quad"`
	SeedSourceQuad     [4]ImagePoint                     `json:"seed_source_quad"`
	Baseline           ExperimentalV4PhoneBuild52State   `json:"baseline"`
	FineRestart        []ExperimentalV4PhoneBuild52State `json:"fine_restart,omitempty"`
}

type ExperimentalV4PhoneBuild52Report struct {
	WorkingWidth          int                                   `json:"working_width"`
	WorkingHeight         int                                   `json:"working_height"`
	Downsampled           bool                                  `json:"downsampled"`
	BoundaryDetected      bool                                  `json:"boundary_detected"`
	BoundaryConfidence    float64                               `json:"boundary_confidence"`
	PairRanking           []ExperimentalV4PhonePairScore        `json:"pair_ranking,omitempty"`
	FrozenCandidates      int                                   `json:"frozen_candidates"`
	SeedsPerPair          int                                   `json:"seeds_per_pair"`
	SeedsSelected         int                                   `json:"seeds_selected"`
	BaselineEvaluations   int                                   `json:"baseline_evaluations"`
	FineEvaluations       int                                   `json:"fine_evaluations"`
	FineStatesFrozen      int                                   `json:"fine_states_frozen"`
	BaselineQualified     int                                   `json:"baseline_qualified"`
	BaselineAuthenticated int                                   `json:"baseline_authenticated"`
	FineQualified         int                                   `json:"fine_qualified"`
	FineAuthenticated     int                                   `json:"fine_authenticated"`
	Candidates            []ExperimentalV4PhoneBuild52Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild52BlindState struct {
	step  float64
	pass  int
	dim   int
	delta float64
	hyp   experimentalV4PhoneHypothesis
}

// experimentalV4PhoneBuild52FineRestart starts again from the untouched seed
// with a deliberately fine-only 2px -> 1px schedule. It uses the exact Build41
// proposal fold and the same per-dimension +/- greedy decision, but retains the
// seed and every accepted intermediate state instead of returning only the last
// state. Held-out/full-pilot/key/oracle evidence is unavailable here.
func experimentalV4PhoneBuild52FineRestart(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor, start [4]ImagePoint, heldout int) ([]experimentalV4PhoneBuild52BlindState, int) {
	q := start
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return nil, 0
	}
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, heldout, true, false)
	evals := 1
	if math.IsInf(score, 0) || math.IsNaN(score) {
		return nil, evals
	}
	states := []experimentalV4PhoneBuild52BlindState{{step: 0, pass: 0, dim: -1, hyp: experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}}}
	for _, step := range experimentalV4PhoneBuild52FineSteps {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 8; dim++ {
				bestQ, bestH, best := q, h, score
				bestDelta := 0.0
				for _, sign := range []float64{-1, 1} {
					delta := sign * step
					qq := experimentalV4PhoneBuild51ApplyDim(q, dim, delta)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					ss, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
					evals++
					if !math.IsInf(ss, 0) && !math.IsNaN(ss) && ss > best+1e-7 {
						bestQ, bestH, best, bestDelta = qq, hh, ss, delta
					}
				}
				if best > score+1e-7 {
					q, h, score = bestQ, bestH, best
					improved = true
					states = append(states, experimentalV4PhoneBuild52BlindState{
						step: step, pass: pass + 1, dim: dim, delta: bestDelta,
						hyp: experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score},
					})
				}
			}
			if !improved {
				break
			}
		}
	}
	return states, evals
}

func experimentalV4PhoneBuild52StateFromQualified(index int, blind experimentalV4PhoneBuild52BlindState, qualified experimentalV4PhoneHypothesis, ok bool, seedSource [4]ImagePoint, srcBounds image.Rectangle, work image.Image) ExperimentalV4PhoneBuild52State {
	axis := ""
	corner := -1
	if blind.dim >= 0 {
		corner = blind.dim / 2
		axis = "y"
		if blind.dim%2 == 0 {
			axis = "x"
		}
	}
	st := ExperimentalV4PhoneBuild52State{
		Index: index, StepSizePx: blind.step, Pass: blind.pass, Dimension: blind.dim, Corner: corner, Axis: axis, DeltaPx: blind.delta,
		Proposal: qualified.proposal, Validation: qualified.validation, PilotScore: qualified.detection.Score, PilotMargin: qualified.detection.Margin,
		Qualified: ok, OriginXBlocks: qualified.detection.OriginXBlocks, OriginYBlocks: qualified.detection.OriginYBlocks,
		WorkingQuad: qualified.quad,
	}
	st.SourceQuad = experimentalV4PhoneBuild46SourceQuad(qualified.quad, srcBounds, work)
	st.MeanMovementFromSeedPx, st.MaxMovementFromSeedPx = experimentalV4PhoneBuild48QuadMovement(seedSource, st.SourceQuad)
	return st
}

// ExperimentalV4PhoneBuild52Diagnose is research-only. It compares the
// unchanged Build41 coarse-to-fine endpoint with a proposal-only fine-first
// restart path for the exact same Build50/51 top-4 seeds. Every geometry is
// generated and frozen before held-out qualification or diagnostic HMAC. The
// function accepts no reference/oracle geometry.
func ExperimentalV4PhoneBuild52Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild52Report, error) {
	var report ExperimentalV4PhoneBuild52Report
	if src == nil {
		return report, errors.New("nil image")
	}
	if len(key) < 8 {
		return report, errors.New("key must contain at least 8 bytes")
	}
	if cw < experimentalV4TileWidthBlocks*blockSize || ch < experimentalV4TileHeightBlocks*blockSize || cw%blockSize != 0 || ch%blockSize != 0 {
		return report, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}

	work, boundary, down := experimentalV4PhoneBuild51Prepare(src)
	report.WorkingWidth, report.WorkingHeight = work.Bounds().Dx(), work.Bounds().Dy()
	report.Downsampled = down
	report.BoundaryDetected, report.BoundaryConfidence = boundary.Detected, boundary.Confidence
	report.SeedsPerPair = experimentalV4PhoneBuild52SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild52SeedsPerPair)
	report.SeedsSelected = len(seeds)

	type blindCandidate struct {
		seed     experimentalV4PhoneBuild48Seed
		baseline experimentalV4PhoneHypothesis
		fine     []experimentalV4PhoneBuild52BlindState
	}
	blind := make([]blindCandidate, 0, len(seeds))
	for _, seed := range seeds {
		baseline, n := experimentalV4PhoneBuild41Refine(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		report.BaselineEvaluations += n
		fine, n := experimentalV4PhoneBuild52FineRestart(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		report.FineEvaluations += n
		report.FineStatesFrozen += len(fine)
		if baseline.h.h[8] == 0 || len(fine) == 0 {
			continue
		}
		blind = append(blind, blindCandidate{seed: seed, baseline: baseline, fine: fine})
	}

	// Freeze barrier: only after every baseline/fine geometry is fixed may
	// held-out qualification, protected data and HMAC be evaluated.
	srcBounds := src.Bounds()
	for i, bc := range blind {
		seedSource := experimentalV4PhoneBuild46SourceQuad(bc.seed.frozen.h.quad, srcBounds, work)
		c := ExperimentalV4PhoneBuild52Candidate{
			Index: i, SeedIndex: bc.seed.index, SeedRankWithinPair: bc.seed.rankWithinPair,
			SourcePair: bc.seed.frozen.h.build43Pair, SourcePairRank: bc.seed.frozen.h.build43PairRank,
			SourceTier: experimentalV4PhoneBuild47TierName(bc.seed.frozen.tier), CellRank: bc.seed.frozen.cellRank,
			SeedProposal: bc.seed.frozen.h.proposal, SeedWorkingQuad: bc.seed.frozen.h.quad, SeedSourceQuad: seedSource,
		}

		baseQ, _, baseOK := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bc.baseline, 0)
		baseBlind := experimentalV4PhoneBuild52BlindState{step: 16, pass: 0, dim: -1, hyp: bc.baseline}
		c.Baseline = experimentalV4PhoneBuild52StateFromQualified(0, baseBlind, baseQ, baseOK, seedSource, srcBounds, work)
		if baseOK {
			report.BaselineQualified++
			c.Baseline.SingleDecodeAttempted = true
			c.Baseline.SingleProfilesTried, c.Baseline.SingleListFramesTried, c.Baseline.SingleMaxDataConfidence, c.Baseline.SingleHMACAuthenticated, c.Baseline.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, baseQ, key, cw, ch)
			if c.Baseline.SingleHMACAuthenticated {
				report.BaselineAuthenticated++
			}
		}

		for si, fs := range bc.fine {
			q, _, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, fs.hyp, 0)
			st := experimentalV4PhoneBuild52StateFromQualified(si, fs, q, ok, seedSource, srcBounds, work)
			if ok {
				report.FineQualified++
				st.SingleDecodeAttempted = true
				st.SingleProfilesTried, st.SingleListFramesTried, st.SingleMaxDataConfidence, st.SingleHMACAuthenticated, st.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, q, key, cw, ch)
				if st.SingleHMACAuthenticated {
					report.FineAuthenticated++
				}
			}
			c.FineRestart = append(c.FineRestart, st)
		}
		report.Candidates = append(report.Candidates, c)
	}
	return report, nil
}
