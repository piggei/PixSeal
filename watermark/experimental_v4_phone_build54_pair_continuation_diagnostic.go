package watermark

import (
	"errors"
	"image"
	"math"
)

const (
	experimentalV4PhoneBuild54SeedsPerPair    = experimentalV4PhoneBuild53SeedsPerPair
	experimentalV4PhoneBuild54Dimensions      = 8
	experimentalV4PhoneBuild54MaxPasses       = 8
	experimentalV4PhoneBuild54MaxStatesBranch = experimentalV4PhoneBuild54Dimensions * experimentalV4PhoneBuild54MaxPasses
)

var experimentalV4PhoneBuild54Step = 1.0

// ExperimentalV4PhoneBuild54ContinuationState is one proposal-improving 1px
// intermediate retained while continuing from a frozen Build53 pair state.
// Geometry/proposal are created blind. Qualification and HMAC are attached only
// after the complete Build54 continuation bank has been frozen.
type ExperimentalV4PhoneBuild54ContinuationState struct {
	Index                   int           `json:"index"`
	Pass                    int           `json:"pass"`
	Dimension               int           `json:"dimension"`
	Corner                  int           `json:"corner"`
	Axis                    string        `json:"axis"`
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

// ExperimentalV4PhoneBuild54Branch continues one retained Build53 coupled
// pair state with bounded 1px proposal-only coordinate descent. Every accepted
// intermediate is retained; the endpoint is not privileged.
type ExperimentalV4PhoneBuild54Branch struct {
	PairRank                int                                           `json:"pair_rank"`
	ParentPair              ExperimentalV4PhoneBuild53PairState           `json:"parent_pair"`
	ContinuationEvaluations int                                           `json:"continuation_evaluations"`
	ContinuationStates      []ExperimentalV4PhoneBuild54ContinuationState `json:"continuation_states,omitempty"`
}

type ExperimentalV4PhoneBuild54Root struct {
	RootIndex         int                                `json:"root_index"`
	Root              ExperimentalV4PhoneBuild52State    `json:"root"`
	SingleEvaluations int                                `json:"single_evaluations"`
	SingleImproving   int                                `json:"single_improving"`
	CoordinateLocal   bool                               `json:"coordinate_local"`
	PairEvaluations   int                                `json:"pair_evaluations"`
	PairImproving     int                                `json:"pair_improving"`
	PairRetained      int                                `json:"pair_retained"`
	Branches          []ExperimentalV4PhoneBuild54Branch `json:"branches,omitempty"`
}

type ExperimentalV4PhoneBuild54Candidate struct {
	Index              int                              `json:"index"`
	SeedIndex          int                              `json:"seed_index"`
	SeedRankWithinPair int                              `json:"seed_rank_within_pair"`
	SourcePair         string                           `json:"source_pair,omitempty"`
	SourcePairRank     int                              `json:"source_pair_rank,omitempty"`
	SourceTier         string                           `json:"source_tier,omitempty"`
	CellRank           int                              `json:"cell_rank"`
	SeedProposal       float64                          `json:"seed_proposal"`
	SeedWorkingQuad    [4]ImagePoint                    `json:"seed_working_quad"`
	SeedSourceQuad     [4]ImagePoint                    `json:"seed_source_quad"`
	Baseline           ExperimentalV4PhoneBuild52State  `json:"baseline"`
	Roots              []ExperimentalV4PhoneBuild54Root `json:"roots,omitempty"`
}

type ExperimentalV4PhoneBuild54Report struct {
	WorkingWidth              int                                   `json:"working_width"`
	WorkingHeight             int                                   `json:"working_height"`
	Downsampled               bool                                  `json:"downsampled"`
	BoundaryDetected          bool                                  `json:"boundary_detected"`
	BoundaryConfidence        float64                               `json:"boundary_confidence"`
	PairRanking               []ExperimentalV4PhonePairScore        `json:"pair_ranking,omitempty"`
	FrozenCandidates          int                                   `json:"frozen_candidates"`
	SeedsPerPair              int                                   `json:"seeds_per_pair"`
	SeedsSelected             int                                   `json:"seeds_selected"`
	BaselineEvaluations       int                                   `json:"baseline_evaluations"`
	RootEvaluations           int                                   `json:"root_evaluations"`
	SingleEvaluations         int                                   `json:"single_evaluations"`
	PairEvaluations           int                                   `json:"pair_evaluations"`
	ContinuationEvaluations   int                                   `json:"continuation_evaluations"`
	RootsFrozen               int                                   `json:"roots_frozen"`
	CoordinateLocalRoots      int                                   `json:"coordinate_local_roots"`
	PairImproving             int                                   `json:"pair_improving"`
	PairStatesFrozen          int                                   `json:"pair_states_frozen"`
	ContinuationStatesFrozen  int                                   `json:"continuation_states_frozen"`
	BaselineQualified         int                                   `json:"baseline_qualified"`
	BaselineAuthenticated     int                                   `json:"baseline_authenticated"`
	RootsQualified            int                                   `json:"roots_qualified"`
	RootsAuthenticated        int                                   `json:"roots_authenticated"`
	PairQualified             int                                   `json:"pair_qualified"`
	PairAuthenticated         int                                   `json:"pair_authenticated"`
	ContinuationQualified     int                                   `json:"continuation_qualified"`
	ContinuationAuthenticated int                                   `json:"continuation_authenticated"`
	Candidates                []ExperimentalV4PhoneBuild54Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild54BlindContinuationState struct {
	index, pass, dim int
	delta            float64
	hyp              experimentalV4PhoneHypothesis
}

type experimentalV4PhoneBuild54BlindBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild54BlindContinuationState
}

type experimentalV4PhoneBuild54BlindRoot struct {
	root              experimentalV4PhoneBuild52BlindState
	singleEvaluations int
	singleImproving   int
	pairEvaluations   int
	pairImproving     int
	branches          []experimentalV4PhoneBuild54BlindBranch
}

func experimentalV4PhoneBuild54Continue(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild54BlindContinuationState, int) {
	q, h, score := start.quad, start.h, start.proposal
	states := make([]experimentalV4PhoneBuild54BlindContinuationState, 0, experimentalV4PhoneBuild54MaxStatesBranch)
	evals := 0
	for pass := 0; pass < experimentalV4PhoneBuild54MaxPasses; pass++ {
		improved := false
		for dim := 0; dim < experimentalV4PhoneBuild54Dimensions; dim++ {
			bestQ, bestH, best := q, h, score
			bestDelta := 0.0
			for _, sign := range []float64{-1, 1} {
				delta := sign * experimentalV4PhoneBuild54Step
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
				states = append(states, experimentalV4PhoneBuild54BlindContinuationState{
					index: len(states) + 1,
					pass:  pass + 1,
					dim:   dim,
					delta: bestDelta,
					hyp:   experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score},
				})
			}
		}
		if !improved {
			break
		}
	}
	return states, evals
}

func experimentalV4PhoneBuild54StateFromQualified(blind experimentalV4PhoneBuild54BlindContinuationState, qualified experimentalV4PhoneHypothesis, ok bool, seedSource [4]ImagePoint, srcBounds image.Rectangle, work image.Image) ExperimentalV4PhoneBuild54ContinuationState {
	corner, axis := experimentalV4PhoneBuild53Axis(blind.dim)
	st := ExperimentalV4PhoneBuild54ContinuationState{
		Index: blind.index, Pass: blind.pass, Dimension: blind.dim, Corner: corner, Axis: axis, DeltaPx: blind.delta,
		Proposal: qualified.proposal, Validation: qualified.validation, PilotScore: qualified.detection.Score, PilotMargin: qualified.detection.Margin,
		Qualified: ok, OriginXBlocks: qualified.detection.OriginXBlocks, OriginYBlocks: qualified.detection.OriginYBlocks,
		WorkingQuad: qualified.quad,
	}
	st.SourceQuad = experimentalV4PhoneBuild46SourceQuad(qualified.quad, srcBounds, work)
	st.MeanMovementFromSeedPx, st.MaxMovementFromSeedPx = experimentalV4PhoneBuild48QuadMovement(seedSource, st.SourceQuad)
	return st
}

// ExperimentalV4PhoneBuild54Diagnose is research-only. It reproduces the
// unchanged Build53 blind pair bank, then continues every retained pair state
// with bounded 1px proposal-only coordinate descent while retaining every
// accepted intermediate. The complete geometry bank is frozen before any
// held-out/full-pilot qualification, protected-data decode, HMAC or oracle use.
func ExperimentalV4PhoneBuild54Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild54Report, error) {
	var report ExperimentalV4PhoneBuild54Report
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
	report.SeedsPerPair = experimentalV4PhoneBuild54SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild54SeedsPerPair)
	report.SeedsSelected = len(seeds)

	type blindCandidate struct {
		seed     experimentalV4PhoneBuild48Seed
		baseline experimentalV4PhoneHypothesis
		roots    []experimentalV4PhoneBuild54BlindRoot
	}
	blind := make([]blindCandidate, 0, len(seeds))
	for _, seed := range seeds {
		baseline, n := experimentalV4PhoneBuild41Refine(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		report.BaselineEvaluations += n
		roots, n := experimentalV4PhoneBuild53TwoPxRoots(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		report.RootEvaluations += n
		if baseline.h.h[8] == 0 || len(roots) == 0 {
			continue
		}
		bc := blindCandidate{seed: seed, baseline: baseline, roots: make([]experimentalV4PhoneBuild54BlindRoot, 0, len(roots))}
		for _, root := range roots {
			singleEvals, singleImproving := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
			report.SingleEvaluations += singleEvals
			br := experimentalV4PhoneBuild54BlindRoot{root: root, singleEvaluations: singleEvals, singleImproving: singleImproving}
			if singleImproving == 0 {
				report.CoordinateLocalRoots++
				pairs, pairEvals, pairImproving := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
				br.pairEvaluations, br.pairImproving = pairEvals, pairImproving
				report.PairEvaluations += pairEvals
				report.PairImproving += pairImproving
				report.PairStatesFrozen += len(pairs)
				br.branches = make([]experimentalV4PhoneBuild54BlindBranch, 0, len(pairs))
				for _, pair := range pairs {
					states, continuationEvals := experimentalV4PhoneBuild54Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
					report.ContinuationEvaluations += continuationEvals
					report.ContinuationStatesFrozen += len(states)
					br.branches = append(br.branches, experimentalV4PhoneBuild54BlindBranch{pair: pair, evals: continuationEvals, states: states})
				}
			}
			report.RootsFrozen++
			bc.roots = append(bc.roots, br)
		}
		blind = append(blind, bc)
	}

	// Freeze barrier: all roots, pairs and continuation intermediates for every
	// seed are fixed before qualification, protected-data decode or HMAC.
	srcBounds := src.Bounds()
	for i, bc := range blind {
		seedSource := experimentalV4PhoneBuild46SourceQuad(bc.seed.frozen.h.quad, srcBounds, work)
		c := ExperimentalV4PhoneBuild54Candidate{
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

		for ri, br := range bc.roots {
			rootQ, _, rootOK := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, br.root.hyp, 0)
			root := experimentalV4PhoneBuild52StateFromQualified(ri, br.root, rootQ, rootOK, seedSource, srcBounds, work)
			if rootOK {
				report.RootsQualified++
				root.SingleDecodeAttempted = true
				root.SingleProfilesTried, root.SingleListFramesTried, root.SingleMaxDataConfidence, root.SingleHMACAuthenticated, root.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, rootQ, key, cw, ch)
				if root.SingleHMACAuthenticated {
					report.RootsAuthenticated++
				}
			}
			outRoot := ExperimentalV4PhoneBuild54Root{
				RootIndex: ri, Root: root, SingleEvaluations: br.singleEvaluations, SingleImproving: br.singleImproving,
				CoordinateLocal: br.singleImproving == 0, PairEvaluations: br.pairEvaluations, PairImproving: br.pairImproving, PairRetained: len(br.branches),
			}
			for pi, bb := range br.branches {
				pairQ, _, pairOK := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bb.pair.hyp, 0)
				parent := experimentalV4PhoneBuild53PairStateFromQualified(pi+1, bb.pair, pairQ, pairOK, seedSource, srcBounds, work)
				if pairOK {
					report.PairQualified++
					parent.SingleDecodeAttempted = true
					parent.SingleProfilesTried, parent.SingleListFramesTried, parent.SingleMaxDataConfidence, parent.SingleHMACAuthenticated, parent.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, pairQ, key, cw, ch)
					if parent.SingleHMACAuthenticated {
						report.PairAuthenticated++
					}
				}
				branch := ExperimentalV4PhoneBuild54Branch{PairRank: pi + 1, ParentPair: parent, ContinuationEvaluations: bb.evals}
				for _, bs := range bb.states {
					q, _, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bs.hyp, 0)
					st := experimentalV4PhoneBuild54StateFromQualified(bs, q, ok, seedSource, srcBounds, work)
					if ok {
						report.ContinuationQualified++
						st.SingleDecodeAttempted = true
						st.SingleProfilesTried, st.SingleListFramesTried, st.SingleMaxDataConfidence, st.SingleHMACAuthenticated, st.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, q, key, cw, ch)
						if st.SingleHMACAuthenticated {
							report.ContinuationAuthenticated++
						}
					}
					branch.ContinuationStates = append(branch.ContinuationStates, st)
				}
				outRoot.Branches = append(outRoot.Branches, branch)
			}
			c.Roots = append(c.Roots, outRoot)
		}
		report.Candidates = append(report.Candidates, c)
	}
	return report, nil
}
