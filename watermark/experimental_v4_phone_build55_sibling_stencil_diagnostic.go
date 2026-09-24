package watermark

import (
	"errors"
	"image"
	"math"
)

const (
	experimentalV4PhoneBuild55SeedsPerPair    = experimentalV4PhoneBuild53SeedsPerPair
	experimentalV4PhoneBuild55Dimensions      = 8
	experimentalV4PhoneBuild55MaxPasses       = 8
	experimentalV4PhoneBuild55MaxStatesBranch = experimentalV4PhoneBuild55Dimensions * experimentalV4PhoneBuild55MaxPasses
	experimentalV4PhoneBuild55MaxSiblingEvals = experimentalV4PhoneBuild55Dimensions * 2
)

var experimentalV4PhoneBuild55Step = 1.0

// ExperimentalV4PhoneBuild55SiblingState is one proposal-improving one-step
// sibling of a frozen Build54 continuation state. The sibling stencil evaluates
// all eight coordinates in both +/-1px directions from the same parent state;
// it does not inherit Gauss-Seidel updates from an earlier coordinate.
type ExperimentalV4PhoneBuild55SiblingState struct {
	Rank                    int           `json:"rank"`
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

// ExperimentalV4PhoneBuild55ContinuationState is one proposal-improving 1px
// intermediate retained while continuing from a frozen Build53 pair state.
// Geometry/proposal are created blind. Qualification and HMAC are attached only
// after the complete Build55 continuation bank has been frozen.
type ExperimentalV4PhoneBuild55ContinuationState struct {
	Index                   int                                      `json:"index"`
	Pass                    int                                      `json:"pass"`
	Dimension               int                                      `json:"dimension"`
	Corner                  int                                      `json:"corner"`
	Axis                    string                                   `json:"axis"`
	DeltaPx                 float64                                  `json:"delta_px"`
	Proposal                float64                                  `json:"proposal"`
	Validation              float64                                  `json:"validation"`
	PilotScore              float64                                  `json:"pilot_score"`
	PilotMargin             float64                                  `json:"pilot_margin"`
	Qualified               bool                                     `json:"qualified"`
	OriginXBlocks           int                                      `json:"origin_x_blocks"`
	OriginYBlocks           int                                      `json:"origin_y_blocks"`
	WorkingQuad             [4]ImagePoint                            `json:"working_quad"`
	SourceQuad              [4]ImagePoint                            `json:"source_quad"`
	MeanMovementFromSeedPx  float64                                  `json:"mean_movement_from_seed_px"`
	MaxMovementFromSeedPx   float64                                  `json:"max_movement_from_seed_px"`
	SingleDecodeAttempted   bool                                     `json:"single_decode_attempted"`
	SingleProfilesTried     int                                      `json:"single_profiles_tried,omitempty"`
	SingleListFramesTried   int                                      `json:"single_list_frames_tried,omitempty"`
	SingleMaxDataConfidence float64                                  `json:"single_max_data_confidence,omitempty"`
	SingleHMACAuthenticated bool                                     `json:"single_hmac_authenticated"`
	SingleProfile           Profile                                  `json:"single_profile,omitempty"`
	SiblingEvaluations      int                                      `json:"sibling_evaluations"`
	SiblingImproving        int                                      `json:"sibling_improving"`
	SiblingStates           []ExperimentalV4PhoneBuild55SiblingState `json:"sibling_states,omitempty"`
}

// ExperimentalV4PhoneBuild55Branch continues one retained Build53 coupled
// pair state with bounded 1px proposal-only coordinate descent. Every accepted
// intermediate is retained; the endpoint is not privileged.
type ExperimentalV4PhoneBuild55Branch struct {
	PairRank                int                                           `json:"pair_rank"`
	ParentPair              ExperimentalV4PhoneBuild53PairState           `json:"parent_pair"`
	ContinuationEvaluations int                                           `json:"continuation_evaluations"`
	ContinuationStates      []ExperimentalV4PhoneBuild55ContinuationState `json:"continuation_states,omitempty"`
}

type ExperimentalV4PhoneBuild55Root struct {
	RootIndex         int                                `json:"root_index"`
	Root              ExperimentalV4PhoneBuild52State    `json:"root"`
	SingleEvaluations int                                `json:"single_evaluations"`
	SingleImproving   int                                `json:"single_improving"`
	CoordinateLocal   bool                               `json:"coordinate_local"`
	PairEvaluations   int                                `json:"pair_evaluations"`
	PairImproving     int                                `json:"pair_improving"`
	PairRetained      int                                `json:"pair_retained"`
	Branches          []ExperimentalV4PhoneBuild55Branch `json:"branches,omitempty"`
}

type ExperimentalV4PhoneBuild55Candidate struct {
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
	Roots              []ExperimentalV4PhoneBuild55Root `json:"roots,omitempty"`
}

type ExperimentalV4PhoneBuild55Report struct {
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
	SiblingEvaluations        int                                   `json:"sibling_evaluations"`
	SiblingStatesFrozen       int                                   `json:"sibling_states_frozen"`
	BaselineQualified         int                                   `json:"baseline_qualified"`
	BaselineAuthenticated     int                                   `json:"baseline_authenticated"`
	RootsQualified            int                                   `json:"roots_qualified"`
	RootsAuthenticated        int                                   `json:"roots_authenticated"`
	PairQualified             int                                   `json:"pair_qualified"`
	PairAuthenticated         int                                   `json:"pair_authenticated"`
	ContinuationQualified     int                                   `json:"continuation_qualified"`
	ContinuationAuthenticated int                                   `json:"continuation_authenticated"`
	SiblingQualified          int                                   `json:"sibling_qualified"`
	SiblingAuthenticated      int                                   `json:"sibling_authenticated"`
	Candidates                []ExperimentalV4PhoneBuild55Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild55BlindSiblingState struct {
	dim   int
	delta float64
	hyp   experimentalV4PhoneHypothesis
}

type experimentalV4PhoneBuild55BlindContinuationState struct {
	index, pass, dim int
	delta            float64
	hyp              experimentalV4PhoneHypothesis
	siblingEvals     int
	siblings         []experimentalV4PhoneBuild55BlindSiblingState
}

type experimentalV4PhoneBuild55BlindBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild55BlindContinuationState
}

type experimentalV4PhoneBuild55BlindRoot struct {
	root              experimentalV4PhoneBuild52BlindState
	singleEvaluations int
	singleImproving   int
	pairEvaluations   int
	pairImproving     int
	branches          []experimentalV4PhoneBuild55BlindBranch
}

func experimentalV4PhoneBuild55Continue(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, start experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild55BlindContinuationState, int) {
	q, h, score := start.quad, start.h, start.proposal
	states := make([]experimentalV4PhoneBuild55BlindContinuationState, 0, experimentalV4PhoneBuild55MaxStatesBranch)
	evals := 0
	for pass := 0; pass < experimentalV4PhoneBuild55MaxPasses; pass++ {
		improved := false
		for dim := 0; dim < experimentalV4PhoneBuild55Dimensions; dim++ {
			bestQ, bestH, best := q, h, score
			bestDelta := 0.0
			for _, sign := range []float64{-1, 1} {
				delta := sign * experimentalV4PhoneBuild55Step
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
				states = append(states, experimentalV4PhoneBuild55BlindContinuationState{
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

func experimentalV4PhoneBuild55SiblingStencil(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, parent experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild55BlindSiblingState, int) {
	states := make([]experimentalV4PhoneBuild55BlindSiblingState, 0, experimentalV4PhoneBuild55MaxSiblingEvals)
	evals := 0
	for dim := 0; dim < experimentalV4PhoneBuild55Dimensions; dim++ {
		for _, sign := range []float64{-1, 1} {
			qq := experimentalV4PhoneBuild51ApplyDim(parent.quad, dim, sign*experimentalV4PhoneBuild55Step)
			if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
				continue
			}
			hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
			if !ok {
				continue
			}
			ss, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
			evals++
			if math.IsInf(ss, 0) || math.IsNaN(ss) || ss <= parent.proposal+1e-7 {
				continue
			}
			states = append(states, experimentalV4PhoneBuild55BlindSiblingState{
				dim: dim, delta: sign * experimentalV4PhoneBuild55Step,
				hyp: experimentalV4PhoneHypothesis{quad: qq, h: hh, proposal: ss},
			})
		}
	}
	return states, evals
}

func experimentalV4PhoneBuild55SiblingFromQualified(rank int, blind experimentalV4PhoneBuild55BlindSiblingState, qualified experimentalV4PhoneHypothesis, ok bool, seedSource [4]ImagePoint, srcBounds image.Rectangle, work image.Image) ExperimentalV4PhoneBuild55SiblingState {
	corner, axis := experimentalV4PhoneBuild53Axis(blind.dim)
	st := ExperimentalV4PhoneBuild55SiblingState{
		Rank: rank, Dimension: blind.dim, Corner: corner, Axis: axis, DeltaPx: blind.delta,
		Proposal: qualified.proposal, Validation: qualified.validation, PilotScore: qualified.detection.Score, PilotMargin: qualified.detection.Margin,
		Qualified: ok, OriginXBlocks: qualified.detection.OriginXBlocks, OriginYBlocks: qualified.detection.OriginYBlocks,
		WorkingQuad: qualified.quad,
	}
	st.SourceQuad = experimentalV4PhoneBuild46SourceQuad(qualified.quad, srcBounds, work)
	st.MeanMovementFromSeedPx, st.MaxMovementFromSeedPx = experimentalV4PhoneBuild48QuadMovement(seedSource, st.SourceQuad)
	return st
}

func experimentalV4PhoneBuild55StateFromQualified(blind experimentalV4PhoneBuild55BlindContinuationState, qualified experimentalV4PhoneHypothesis, ok bool, seedSource [4]ImagePoint, srcBounds image.Rectangle, work image.Image) ExperimentalV4PhoneBuild55ContinuationState {
	corner, axis := experimentalV4PhoneBuild53Axis(blind.dim)
	st := ExperimentalV4PhoneBuild55ContinuationState{
		Index: blind.index, Pass: blind.pass, Dimension: blind.dim, Corner: corner, Axis: axis, DeltaPx: blind.delta,
		Proposal: qualified.proposal, Validation: qualified.validation, PilotScore: qualified.detection.Score, PilotMargin: qualified.detection.Margin,
		Qualified: ok, OriginXBlocks: qualified.detection.OriginXBlocks, OriginYBlocks: qualified.detection.OriginYBlocks,
		WorkingQuad: qualified.quad, SiblingEvaluations: blind.siblingEvals, SiblingImproving: len(blind.siblings),
	}
	st.SourceQuad = experimentalV4PhoneBuild46SourceQuad(qualified.quad, srcBounds, work)
	st.MeanMovementFromSeedPx, st.MaxMovementFromSeedPx = experimentalV4PhoneBuild48QuadMovement(seedSource, st.SourceQuad)
	return st
}

// ExperimentalV4PhoneBuild55Diagnose is research-only. It reproduces the
// unchanged Build54 continuation bank and evaluates every independent +/-1px
// coordinate sibling from each retained continuation state against the same
// frozen parent. Every proposal-improving sibling is retained. The complete
// geometry bank is frozen before held-out/full-pilot qualification, protected-
// data decode, HMAC or oracle use.
func ExperimentalV4PhoneBuild55Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild55Report, error) {
	var report ExperimentalV4PhoneBuild55Report
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
	report.SeedsPerPair = experimentalV4PhoneBuild55SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild55SeedsPerPair)
	report.SeedsSelected = len(seeds)

	type blindCandidate struct {
		seed     experimentalV4PhoneBuild48Seed
		baseline experimentalV4PhoneHypothesis
		roots    []experimentalV4PhoneBuild55BlindRoot
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
		bc := blindCandidate{seed: seed, baseline: baseline, roots: make([]experimentalV4PhoneBuild55BlindRoot, 0, len(roots))}
		for _, root := range roots {
			singleEvals, singleImproving := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
			report.SingleEvaluations += singleEvals
			br := experimentalV4PhoneBuild55BlindRoot{root: root, singleEvaluations: singleEvals, singleImproving: singleImproving}
			if singleImproving == 0 {
				report.CoordinateLocalRoots++
				pairs, pairEvals, pairImproving := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
				br.pairEvaluations, br.pairImproving = pairEvals, pairImproving
				report.PairEvaluations += pairEvals
				report.PairImproving += pairImproving
				report.PairStatesFrozen += len(pairs)
				br.branches = make([]experimentalV4PhoneBuild55BlindBranch, 0, len(pairs))
				for _, pair := range pairs {
					states, continuationEvals := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
					report.ContinuationEvaluations += continuationEvals
					report.ContinuationStatesFrozen += len(states)
					for si := range states {
						siblings, siblingEvals := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, states[si].hyp, 0)
						states[si].siblingEvals = siblingEvals
						states[si].siblings = siblings
						report.SiblingEvaluations += siblingEvals
						report.SiblingStatesFrozen += len(siblings)
					}
					br.branches = append(br.branches, experimentalV4PhoneBuild55BlindBranch{pair: pair, evals: continuationEvals, states: states})
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
		c := ExperimentalV4PhoneBuild55Candidate{
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
			outRoot := ExperimentalV4PhoneBuild55Root{
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
				branch := ExperimentalV4PhoneBuild55Branch{PairRank: pi + 1, ParentPair: parent, ContinuationEvaluations: bb.evals}
				for _, bs := range bb.states {
					q, _, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bs.hyp, 0)
					st := experimentalV4PhoneBuild55StateFromQualified(bs, q, ok, seedSource, srcBounds, work)
					if ok {
						report.ContinuationQualified++
						st.SingleDecodeAttempted = true
						st.SingleProfilesTried, st.SingleListFramesTried, st.SingleMaxDataConfidence, st.SingleHMACAuthenticated, st.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, q, key, cw, ch)
						if st.SingleHMACAuthenticated {
							report.ContinuationAuthenticated++
						}
					}
					for sri, blindSibling := range bs.siblings {
						sq, _, sok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, blindSibling.hyp, 0)
						sibling := experimentalV4PhoneBuild55SiblingFromQualified(sri+1, blindSibling, sq, sok, seedSource, srcBounds, work)
						if sok {
							report.SiblingQualified++
							sibling.SingleDecodeAttempted = true
							sibling.SingleProfilesTried, sibling.SingleListFramesTried, sibling.SingleMaxDataConfidence, sibling.SingleHMACAuthenticated, sibling.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, sq, key, cw, ch)
							if sibling.SingleHMACAuthenticated {
								report.SiblingAuthenticated++
							}
						}
						st.SiblingStates = append(st.SiblingStates, sibling)
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
