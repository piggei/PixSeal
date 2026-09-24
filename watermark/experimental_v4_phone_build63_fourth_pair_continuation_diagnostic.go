package watermark

import (
	"errors"
	"image"
)

const (
	experimentalV4PhoneBuild63SeedsPerPair = experimentalV4PhoneBuild62SeedsPerPair
)

// ExperimentalV4PhoneBuild63PairEscapeState is one proposal-improving coupled
// +/-1px two-coordinate escape from a Build55 sibling that is a local maximum
// under the complete independent one-coordinate +/-1px stencil.
type ExperimentalV4PhoneBuild63PairEscapeState struct {
	ExperimentalV4PhoneBuild53PairState
	ContinuationEvaluations int                                                   `json:"continuation_evaluations"`
	ContinuationStates      []ExperimentalV4PhoneBuild63PostPairContinuationState `json:"continuation_states,omitempty"`
}

// ExperimentalV4PhoneBuild63ThirdPairState is one frozen Build59 third-pair
// escape plus every proposal-improving 1px intermediate retained while
// continuing from it. Geometry/proposal are created blind. Qualification and
// HMAC are attached only after the complete Build63 bank has been frozen.
type ExperimentalV4PhoneBuild63ThirdPairState struct {
	ExperimentalV4PhoneBuild53PairState
	ContinuationEvaluations int                                                    `json:"continuation_evaluations"`
	ContinuationStates      []ExperimentalV4PhoneBuild63ThirdPairContinuationState `json:"continuation_states,omitempty"`
}

// ExperimentalV4PhoneBuild63FourthPairState is one proposal-improving coupled
// +/-1px two-coordinate escape from a frozen Build61 third-pair sibling that
// is proposal-local under the complete independent one-coordinate +/-1px stencil.
type ExperimentalV4PhoneBuild63FourthPairState struct {
	ExperimentalV4PhoneBuild53PairState
	ContinuationEvaluations int                                                     `json:"continuation_evaluations"`
	ContinuationStates      []ExperimentalV4PhoneBuild63FourthPairContinuationState `json:"continuation_states,omitempty"`
}

// ExperimentalV4PhoneBuild63FourthPairContinuationState is one accepted 1px
// proposal-only intermediate after a frozen Build62 fourth-pair escape.
type ExperimentalV4PhoneBuild63FourthPairContinuationState struct {
	ExperimentalV4PhoneBuild55ContinuationState
}

// ExperimentalV4PhoneBuild63ThirdPairSiblingState preserves one frozen Build61
// sibling together with its proposal-only local-max probe and any bounded
// fourth-pair escapes retained from that exact parent.
type ExperimentalV4PhoneBuild63ThirdPairSiblingState struct {
	Rank              int                                         `json:"rank"`
	Parent            ExperimentalV4PhoneBuild55SiblingState      `json:"parent"`
	SingleEvaluations int                                         `json:"single_evaluations"`
	SingleImproving   int                                         `json:"single_improving"`
	CoordinateLocal   bool                                        `json:"coordinate_local"`
	PairEvaluations   int                                         `json:"pair_evaluations"`
	PairImproving     int                                         `json:"pair_improving"`
	PairStates        []ExperimentalV4PhoneBuild63FourthPairState `json:"pair_states,omitempty"`
}

// ExperimentalV4PhoneBuild63ThirdPairContinuationState is one frozen Build60
// post-third-pair continuation state plus every proposal-improving independent
// +/-1px sibling evaluated from that identical parent. Geometry/proposal are
// created blind; qualification and HMAC are attached only after the complete
// Build63 bank has been frozen.
type ExperimentalV4PhoneBuild63ThirdPairContinuationState struct {
	ExperimentalV4PhoneBuild55ContinuationState
	SiblingEvaluations int                                               `json:"sibling_evaluations"`
	SiblingStates      []ExperimentalV4PhoneBuild63ThirdPairSiblingState `json:"sibling_states,omitempty"`
}

// ExperimentalV4PhoneBuild63SecondPairSiblingState preserves one frozen
// Build58 sibling, probes its complete independent one-coordinate +/-1px
// neighborhood with proposal only, and retains bounded coupled pair escapes
// only when that sibling is proposal-local.
type ExperimentalV4PhoneBuild63SecondPairSiblingState struct {
	Rank              int                                        `json:"rank"`
	Parent            ExperimentalV4PhoneBuild55SiblingState     `json:"parent"`
	SingleEvaluations int                                        `json:"single_evaluations"`
	SingleImproving   int                                        `json:"single_improving"`
	CoordinateLocal   bool                                       `json:"coordinate_local"`
	PairEvaluations   int                                        `json:"pair_evaluations"`
	PairImproving     int                                        `json:"pair_improving"`
	PairStates        []ExperimentalV4PhoneBuild63ThirdPairState `json:"pair_states,omitempty"`
}

type ExperimentalV4PhoneBuild63PostPairContinuationState struct {
	ExperimentalV4PhoneBuild55ContinuationState
	SiblingEvaluations int                                                `json:"sibling_evaluations"`
	SiblingStates      []ExperimentalV4PhoneBuild63SecondPairSiblingState `json:"sibling_states,omitempty"`
}

// ExperimentalV4PhoneBuild63SiblingState preserves one frozen Build55 sibling
// and records the proposal-only local-max probe plus any retained pair escapes.
type ExperimentalV4PhoneBuild63SiblingState struct {
	Rank              int                                         `json:"rank"`
	Parent            ExperimentalV4PhoneBuild55SiblingState      `json:"parent"`
	SingleEvaluations int                                         `json:"single_evaluations"`
	SingleImproving   int                                         `json:"single_improving"`
	CoordinateLocal   bool                                        `json:"coordinate_local"`
	PairEvaluations   int                                         `json:"pair_evaluations"`
	PairImproving     int                                         `json:"pair_improving"`
	PairStates        []ExperimentalV4PhoneBuild63PairEscapeState `json:"pair_states,omitempty"`
}

type ExperimentalV4PhoneBuild63ContinuationState struct {
	Parent   ExperimentalV4PhoneBuild55ContinuationState `json:"parent"`
	Siblings []ExperimentalV4PhoneBuild63SiblingState    `json:"siblings,omitempty"`
}

type ExperimentalV4PhoneBuild63Branch struct {
	PairRank           int                                           `json:"pair_rank"`
	ParentPair         ExperimentalV4PhoneBuild53PairState           `json:"parent_pair"`
	ContinuationStates []ExperimentalV4PhoneBuild63ContinuationState `json:"continuation_states,omitempty"`
}

type ExperimentalV4PhoneBuild63Root struct {
	RootIndex         int                                `json:"root_index"`
	Root              ExperimentalV4PhoneBuild52State    `json:"root"`
	SingleEvaluations int                                `json:"single_evaluations"`
	SingleImproving   int                                `json:"single_improving"`
	CoordinateLocal   bool                               `json:"coordinate_local"`
	PairEvaluations   int                                `json:"pair_evaluations"`
	PairImproving     int                                `json:"pair_improving"`
	PairRetained      int                                `json:"pair_retained"`
	Branches          []ExperimentalV4PhoneBuild63Branch `json:"branches,omitempty"`
}

type ExperimentalV4PhoneBuild63Candidate struct {
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
	Roots              []ExperimentalV4PhoneBuild63Root `json:"roots,omitempty"`
}

type ExperimentalV4PhoneBuild63Report struct {
	WorkingWidth                           int                                   `json:"working_width"`
	WorkingHeight                          int                                   `json:"working_height"`
	Downsampled                            bool                                  `json:"downsampled"`
	BoundaryDetected                       bool                                  `json:"boundary_detected"`
	BoundaryConfidence                     float64                               `json:"boundary_confidence"`
	PairRanking                            []ExperimentalV4PhonePairScore        `json:"pair_ranking,omitempty"`
	FrozenCandidates                       int                                   `json:"frozen_candidates"`
	SeedsPerPair                           int                                   `json:"seeds_per_pair"`
	SeedsSelected                          int                                   `json:"seeds_selected"`
	BaselineEvaluations                    int                                   `json:"baseline_evaluations"`
	RootEvaluations                        int                                   `json:"root_evaluations"`
	SingleEvaluations                      int                                   `json:"single_evaluations"`
	PairEvaluations                        int                                   `json:"pair_evaluations"`
	ContinuationEvaluations                int                                   `json:"continuation_evaluations"`
	SiblingEvaluations                     int                                   `json:"sibling_evaluations"`
	SiblingLocalProbeEvaluations           int                                   `json:"sibling_local_probe_evaluations"`
	SiblingPairEvaluations                 int                                   `json:"sibling_pair_evaluations"`
	SecondPairContinuationEvaluations      int                                   `json:"second_pair_continuation_evaluations"`
	RootsFrozen                            int                                   `json:"roots_frozen"`
	CoordinateLocalRoots                   int                                   `json:"coordinate_local_roots"`
	PairStatesFrozen                       int                                   `json:"pair_states_frozen"`
	ContinuationStatesFrozen               int                                   `json:"continuation_states_frozen"`
	SiblingStatesFrozen                    int                                   `json:"sibling_states_frozen"`
	CoordinateLocalSiblings                int                                   `json:"coordinate_local_siblings"`
	SiblingPairImproving                   int                                   `json:"sibling_pair_improving"`
	SiblingPairStatesFrozen                int                                   `json:"sibling_pair_states_frozen"`
	SecondPairContinuationStatesFrozen     int                                   `json:"second_pair_continuation_states_frozen"`
	SecondPairSiblingEvaluations           int                                   `json:"second_pair_sibling_evaluations"`
	SecondPairSiblingStatesFrozen          int                                   `json:"second_pair_sibling_states_frozen"`
	SecondPairSiblingLocalProbeEvaluations int                                   `json:"second_pair_sibling_local_probe_evaluations"`
	CoordinateLocalSecondPairSiblings      int                                   `json:"coordinate_local_second_pair_siblings"`
	ThirdPairEvaluations                   int                                   `json:"third_pair_evaluations"`
	ThirdPairImproving                     int                                   `json:"third_pair_improving"`
	ThirdPairStatesFrozen                  int                                   `json:"third_pair_states_frozen"`
	ThirdPairContinuationEvaluations       int                                   `json:"third_pair_continuation_evaluations"`
	ThirdPairContinuationStatesFrozen      int                                   `json:"third_pair_continuation_states_frozen"`
	ThirdPairSiblingEvaluations            int                                   `json:"third_pair_sibling_evaluations"`
	ThirdPairSiblingStatesFrozen           int                                   `json:"third_pair_sibling_states_frozen"`
	ThirdPairSiblingLocalProbeEvaluations  int                                   `json:"third_pair_sibling_local_probe_evaluations"`
	CoordinateLocalThirdPairSiblings       int                                   `json:"coordinate_local_third_pair_siblings"`
	FourthPairEvaluations                  int                                   `json:"fourth_pair_evaluations"`
	FourthPairImproving                    int                                   `json:"fourth_pair_improving"`
	FourthPairStatesFrozen                 int                                   `json:"fourth_pair_states_frozen"`
	FourthPairContinuationEvaluations      int                                   `json:"fourth_pair_continuation_evaluations"`
	FourthPairContinuationStatesFrozen     int                                   `json:"fourth_pair_continuation_states_frozen"`
	BaselineQualified                      int                                   `json:"baseline_qualified"`
	BaselineAuthenticated                  int                                   `json:"baseline_authenticated"`
	RootsQualified                         int                                   `json:"roots_qualified"`
	RootsAuthenticated                     int                                   `json:"roots_authenticated"`
	PairQualified                          int                                   `json:"pair_qualified"`
	PairAuthenticated                      int                                   `json:"pair_authenticated"`
	ContinuationQualified                  int                                   `json:"continuation_qualified"`
	ContinuationAuthenticated              int                                   `json:"continuation_authenticated"`
	SiblingQualified                       int                                   `json:"sibling_qualified"`
	SiblingAuthenticated                   int                                   `json:"sibling_authenticated"`
	SiblingPairQualified                   int                                   `json:"sibling_pair_qualified"`
	SiblingPairAuthenticated               int                                   `json:"sibling_pair_authenticated"`
	SecondPairContinuationQualified        int                                   `json:"second_pair_continuation_qualified"`
	SecondPairContinuationAuthenticated    int                                   `json:"second_pair_continuation_authenticated"`
	SecondPairSiblingQualified             int                                   `json:"second_pair_sibling_qualified"`
	SecondPairSiblingAuthenticated         int                                   `json:"second_pair_sibling_authenticated"`
	ThirdPairQualified                     int                                   `json:"third_pair_qualified"`
	ThirdPairAuthenticated                 int                                   `json:"third_pair_authenticated"`
	ThirdPairContinuationQualified         int                                   `json:"third_pair_continuation_qualified"`
	ThirdPairContinuationAuthenticated     int                                   `json:"third_pair_continuation_authenticated"`
	ThirdPairSiblingQualified              int                                   `json:"third_pair_sibling_qualified"`
	ThirdPairSiblingAuthenticated          int                                   `json:"third_pair_sibling_authenticated"`
	FourthPairQualified                    int                                   `json:"fourth_pair_qualified"`
	FourthPairAuthenticated                int                                   `json:"fourth_pair_authenticated"`
	FourthPairContinuationQualified        int                                   `json:"fourth_pair_continuation_qualified"`
	FourthPairContinuationAuthenticated    int                                   `json:"fourth_pair_continuation_authenticated"`
	Candidates                             []ExperimentalV4PhoneBuild63Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild63BlindFourthPairBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild55BlindContinuationState
}

type experimentalV4PhoneBuild63BlindThirdPairSibling struct {
	base            experimentalV4PhoneBuild55BlindSiblingState
	singleEvals     int
	singleImproving int
	pairEvals       int
	pairImproving   int
	pairs           []experimentalV4PhoneBuild63BlindFourthPairBranch
}

type experimentalV4PhoneBuild63BlindThirdPairContinuation struct {
	base         experimentalV4PhoneBuild55BlindContinuationState
	siblingEvals int
	siblings     []experimentalV4PhoneBuild63BlindThirdPairSibling
}

type experimentalV4PhoneBuild63BlindThirdPairBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild63BlindThirdPairContinuation
}

type experimentalV4PhoneBuild63BlindSecondPairSibling struct {
	base            experimentalV4PhoneBuild55BlindSiblingState
	singleEvals     int
	singleImproving int
	pairEvals       int
	pairImproving   int
	pairs           []experimentalV4PhoneBuild63BlindThirdPairBranch
}

type experimentalV4PhoneBuild63BlindPostPairContinuation struct {
	base         experimentalV4PhoneBuild55BlindContinuationState
	siblingEvals int
	siblings     []experimentalV4PhoneBuild63BlindSecondPairSibling
}

type experimentalV4PhoneBuild63BlindSecondPairBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild63BlindPostPairContinuation
}

type experimentalV4PhoneBuild63BlindSibling struct {
	base            experimentalV4PhoneBuild55BlindSiblingState
	singleEvals     int
	singleImproving int
	pairEvals       int
	pairImproving   int
	pairs           []experimentalV4PhoneBuild63BlindSecondPairBranch
}

type experimentalV4PhoneBuild63BlindContinuation struct {
	base     experimentalV4PhoneBuild55BlindContinuationState
	siblings []experimentalV4PhoneBuild63BlindSibling
}

type experimentalV4PhoneBuild63BlindBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild63BlindContinuation
}

type experimentalV4PhoneBuild63BlindRoot struct {
	root              experimentalV4PhoneBuild52BlindState
	singleEvaluations int
	singleImproving   int
	pairEvaluations   int
	pairImproving     int
	branches          []experimentalV4PhoneBuild63BlindBranch
}

// ExperimentalV4PhoneBuild63Diagnose is research-only. It reproduces the full
// Build62 blind bank and continues every retained fourth-pair state with bounded proposal-only 1px coordinate descent, retaining every accepted intermediate. The inherited Build61 third-pair siblings were tested against
// the complete independent one-coordinate +/-1px stencil using proposal only.
// Only proposal-local siblings receive the bounded Build53 two-coordinate +/-1px
// pair stencil, retaining at most the same eight states ordered by proposal. Every retained fourth-pair escape is then continued blind. The
// entire expanded B/mild+B/angle geometry bank is frozen before qualification,
// decode, HMAC or oracle use.
func ExperimentalV4PhoneBuild63Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild63Report, error) {
	var report ExperimentalV4PhoneBuild63Report
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
	report.SeedsPerPair = experimentalV4PhoneBuild63SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild63SeedsPerPair)
	report.SeedsSelected = len(seeds)

	type blindCandidate struct {
		seed     experimentalV4PhoneBuild48Seed
		baseline experimentalV4PhoneHypothesis
		roots    []experimentalV4PhoneBuild63BlindRoot
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
		bc := blindCandidate{seed: seed, baseline: baseline, roots: make([]experimentalV4PhoneBuild63BlindRoot, 0, len(roots))}
		for _, root := range roots {
			singleEvals, singleImproving := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
			report.SingleEvaluations += singleEvals
			br := experimentalV4PhoneBuild63BlindRoot{root: root, singleEvaluations: singleEvals, singleImproving: singleImproving}
			if singleImproving == 0 {
				report.CoordinateLocalRoots++
				pairs, pairEvals, pairImproving := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
				br.pairEvaluations, br.pairImproving = pairEvals, pairImproving
				report.PairEvaluations += pairEvals
				report.PairStatesFrozen += len(pairs)
				br.branches = make([]experimentalV4PhoneBuild63BlindBranch, 0, len(pairs))
				for _, pair := range pairs {
					cont, contEvals := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
					report.ContinuationEvaluations += contEvals
					report.ContinuationStatesFrozen += len(cont)
					bb := experimentalV4PhoneBuild63BlindBranch{pair: pair, evals: contEvals, states: make([]experimentalV4PhoneBuild63BlindContinuation, 0, len(cont))}
					for _, cs := range cont {
						sibs, sibEvals := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, cs.hyp, 0)
						cs.siblingEvals = sibEvals
						cs.siblings = sibs
						report.SiblingEvaluations += sibEvals
						report.SiblingStatesFrozen += len(sibs)
						outCS := experimentalV4PhoneBuild63BlindContinuation{base: cs, siblings: make([]experimentalV4PhoneBuild63BlindSibling, 0, len(sibs))}
						for _, sib := range sibs {
							se, si := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, sib.hyp, 0)
							report.SiblingLocalProbeEvaluations += se
							outS := experimentalV4PhoneBuild63BlindSibling{base: sib, singleEvals: se, singleImproving: si}
							if si == 0 {
								report.CoordinateLocalSiblings++
								ps, pe, pi := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, sib.hyp, 0)
								outS.pairEvals, outS.pairImproving = pe, pi
								report.SiblingPairEvaluations += pe
								report.SiblingPairImproving += pi
								report.SiblingPairStatesFrozen += len(ps)
								outS.pairs = make([]experimentalV4PhoneBuild63BlindSecondPairBranch, 0, len(ps))
								for _, p := range ps {
									cont2, ce := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, p.hyp, 0)
									report.SecondPairContinuationEvaluations += ce
									report.SecondPairContinuationStatesFrozen += len(cont2)
									outPair := experimentalV4PhoneBuild63BlindSecondPairBranch{pair: p, evals: ce, states: make([]experimentalV4PhoneBuild63BlindPostPairContinuation, 0, len(cont2))}
									for _, c2 := range cont2 {
										sibs2, se2 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c2.hyp, 0)
										report.SecondPairSiblingEvaluations += se2
										report.SecondPairSiblingStatesFrozen += len(sibs2)
										outC2 := experimentalV4PhoneBuild63BlindPostPairContinuation{base: c2, siblingEvals: se2, siblings: make([]experimentalV4PhoneBuild63BlindSecondPairSibling, 0, len(sibs2))}
										for _, s2 := range sibs2 {
											se3, si3 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s2.hyp, 0)
											report.SecondPairSiblingLocalProbeEvaluations += se3
											outS2 := experimentalV4PhoneBuild63BlindSecondPairSibling{base: s2, singleEvals: se3, singleImproving: si3}
											if si3 == 0 {
												report.CoordinateLocalSecondPairSiblings++
												ps3, pe3, pi3 := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s2.hyp, 0)
												outS2.pairEvals, outS2.pairImproving = pe3, pi3
												report.ThirdPairEvaluations += pe3
												report.ThirdPairImproving += pi3
												report.ThirdPairStatesFrozen += len(ps3)
												outS2.pairs = make([]experimentalV4PhoneBuild63BlindThirdPairBranch, 0, len(ps3))
												for _, p3 := range ps3 {
													cont3, ce3 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, p3.hyp, 0)
													report.ThirdPairContinuationEvaluations += ce3
													report.ThirdPairContinuationStatesFrozen += len(cont3)
													outP3 := experimentalV4PhoneBuild63BlindThirdPairBranch{pair: p3, evals: ce3, states: make([]experimentalV4PhoneBuild63BlindThirdPairContinuation, 0, len(cont3))}
													for _, c3 := range cont3 {
														sibs3, se4 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c3.hyp, 0)
														report.ThirdPairSiblingEvaluations += se4
														report.ThirdPairSiblingStatesFrozen += len(sibs3)
														outC3 := experimentalV4PhoneBuild63BlindThirdPairContinuation{base: c3, siblingEvals: se4, siblings: make([]experimentalV4PhoneBuild63BlindThirdPairSibling, 0, len(sibs3))}
														for _, s3 := range sibs3 {
															se5, si5 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s3.hyp, 0)
															report.ThirdPairSiblingLocalProbeEvaluations += se5
															outS3 := experimentalV4PhoneBuild63BlindThirdPairSibling{base: s3, singleEvals: se5, singleImproving: si5}
															if si5 == 0 {
																report.CoordinateLocalThirdPairSiblings++
																ps4, pe4, pi4 := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s3.hyp, 0)
																outS3.pairEvals, outS3.pairImproving = pe4, pi4
																report.FourthPairEvaluations += pe4
																report.FourthPairImproving += pi4
																report.FourthPairStatesFrozen += len(ps4)
																outS3.pairs = make([]experimentalV4PhoneBuild63BlindFourthPairBranch, 0, len(ps4))
																for _, p4 := range ps4 {
																	cont4, ce4 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, p4.hyp, 0)
																	report.FourthPairContinuationEvaluations += ce4
																	report.FourthPairContinuationStatesFrozen += len(cont4)
																	outS3.pairs = append(outS3.pairs, experimentalV4PhoneBuild63BlindFourthPairBranch{pair: p4, evals: ce4, states: cont4})
																}

															}
															outC3.siblings = append(outC3.siblings, outS3)
														}
														outP3.states = append(outP3.states, outC3)
													}
													outS2.pairs = append(outS2.pairs, outP3)
												}
											}
											outC2.siblings = append(outC2.siblings, outS2)
										}
										outPair.states = append(outPair.states, outC2)
									}
									outS.pairs = append(outS.pairs, outPair)
								}
							}
							outCS.siblings = append(outCS.siblings, outS)
						}
						bb.states = append(bb.states, outCS)
					}
					br.branches = append(br.branches, bb)
				}
			}
			report.RootsFrozen++
			bc.roots = append(bc.roots, br)
		}
		blind = append(blind, bc)
	}

	// Freeze barrier: all roots, first pair escapes, continuations, first siblings,
	// second pair escapes, post-second-pair continuations and every independent
	// second-pair sibling, third-pair escape, every retained post-third-pair
	// continuation state, every independent third-pair sibling and every bounded
	// fourth-pair escape are fixed for every seed before qualification/decode.
	srcBounds := src.Bounds()
	for i, bc := range blind {
		seedSource := experimentalV4PhoneBuild46SourceQuad(bc.seed.frozen.h.quad, srcBounds, work)
		c := ExperimentalV4PhoneBuild63Candidate{
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
			outRoot := ExperimentalV4PhoneBuild63Root{RootIndex: ri, Root: root, SingleEvaluations: br.singleEvaluations, SingleImproving: br.singleImproving, CoordinateLocal: br.singleImproving == 0, PairEvaluations: br.pairEvaluations, PairImproving: br.pairImproving, PairRetained: len(br.branches)}
			for pi, bb := range br.branches {
				pairQ, _, pairOK := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bb.pair.hyp, 0)
				parentPair := experimentalV4PhoneBuild53PairStateFromQualified(pi+1, bb.pair, pairQ, pairOK, seedSource, srcBounds, work)
				if pairOK {
					report.PairQualified++
					parentPair.SingleDecodeAttempted = true
					parentPair.SingleProfilesTried, parentPair.SingleListFramesTried, parentPair.SingleMaxDataConfidence, parentPair.SingleHMACAuthenticated, parentPair.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, pairQ, key, cw, ch)
					if parentPair.SingleHMACAuthenticated {
						report.PairAuthenticated++
					}
				}
				outBranch := ExperimentalV4PhoneBuild63Branch{PairRank: pi + 1, ParentPair: parentPair}
				for _, bcs := range bb.states {
					q, _, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bcs.base.hyp, 0)
					parentState := experimentalV4PhoneBuild55StateFromQualified(bcs.base, q, ok, seedSource, srcBounds, work)
					if ok {
						report.ContinuationQualified++
						parentState.SingleDecodeAttempted = true
						parentState.SingleProfilesTried, parentState.SingleListFramesTried, parentState.SingleMaxDataConfidence, parentState.SingleHMACAuthenticated, parentState.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, q, key, cw, ch)
						if parentState.SingleHMACAuthenticated {
							report.ContinuationAuthenticated++
						}
					}
					outCS := ExperimentalV4PhoneBuild63ContinuationState{Parent: parentState}
					for sri, bs := range bcs.siblings {
						sq, _, sok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bs.base.hyp, 0)
						parentSibling := experimentalV4PhoneBuild55SiblingFromQualified(sri+1, bs.base, sq, sok, seedSource, srcBounds, work)
						if sok {
							report.SiblingQualified++
							parentSibling.SingleDecodeAttempted = true
							parentSibling.SingleProfilesTried, parentSibling.SingleListFramesTried, parentSibling.SingleMaxDataConfidence, parentSibling.SingleHMACAuthenticated, parentSibling.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, sq, key, cw, ch)
							if parentSibling.SingleHMACAuthenticated {
								report.SiblingAuthenticated++
							}
						}
						outS := ExperimentalV4PhoneBuild63SiblingState{Rank: sri + 1, Parent: parentSibling, SingleEvaluations: bs.singleEvals, SingleImproving: bs.singleImproving, CoordinateLocal: bs.singleImproving == 0, PairEvaluations: bs.pairEvals, PairImproving: bs.pairImproving}
						for pri, bp := range bs.pairs {
							pq, _, pok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bp.pair.hyp, 0)
							ps := experimentalV4PhoneBuild53PairStateFromQualified(pri+1, bp.pair, pq, pok, seedSource, srcBounds, work)
							if pok {
								report.SiblingPairQualified++
								ps.SingleDecodeAttempted = true
								ps.SingleProfilesTried, ps.SingleListFramesTried, ps.SingleMaxDataConfidence, ps.SingleHMACAuthenticated, ps.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, pq, key, cw, ch)
								if ps.SingleHMACAuthenticated {
									report.SiblingPairAuthenticated++
								}
							}
							outPair := ExperimentalV4PhoneBuild63PairEscapeState{ExperimentalV4PhoneBuild53PairState: ps, ContinuationEvaluations: bp.evals}
							for _, bcont := range bp.states {
								cq, _, cok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bcont.base.hyp, 0)
								cs := experimentalV4PhoneBuild55StateFromQualified(bcont.base, cq, cok, seedSource, srcBounds, work)
								if cok {
									report.SecondPairContinuationQualified++
									cs.SingleDecodeAttempted = true
									cs.SingleProfilesTried, cs.SingleListFramesTried, cs.SingleMaxDataConfidence, cs.SingleHMACAuthenticated, cs.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, cq, key, cw, ch)
									if cs.SingleHMACAuthenticated {
										report.SecondPairContinuationAuthenticated++
									}
								}
								outCont := ExperimentalV4PhoneBuild63PostPairContinuationState{ExperimentalV4PhoneBuild55ContinuationState: cs, SiblingEvaluations: bcont.siblingEvals}
								for si, bs2 := range bcont.siblings {
									sq2, _, sok2 := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bs2.base.hyp, 0)
									ss2 := experimentalV4PhoneBuild55SiblingFromQualified(si+1, bs2.base, sq2, sok2, seedSource, srcBounds, work)
									if sok2 {
										report.SecondPairSiblingQualified++
										ss2.SingleDecodeAttempted = true
										ss2.SingleProfilesTried, ss2.SingleListFramesTried, ss2.SingleMaxDataConfidence, ss2.SingleHMACAuthenticated, ss2.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, sq2, key, cw, ch)
										if ss2.SingleHMACAuthenticated {
											report.SecondPairSiblingAuthenticated++
										}
									}
									outS2 := ExperimentalV4PhoneBuild63SecondPairSiblingState{Rank: si + 1, Parent: ss2, SingleEvaluations: bs2.singleEvals, SingleImproving: bs2.singleImproving, CoordinateLocal: bs2.singleImproving == 0, PairEvaluations: bs2.pairEvals, PairImproving: bs2.pairImproving}
									for pri3, bp3 := range bs2.pairs {
										pq3, _, pok3 := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bp3.pair.hyp, 0)
										ps3 := experimentalV4PhoneBuild53PairStateFromQualified(pri3+1, bp3.pair, pq3, pok3, seedSource, srcBounds, work)
										if pok3 {
											report.ThirdPairQualified++
											ps3.SingleDecodeAttempted = true
											ps3.SingleProfilesTried, ps3.SingleListFramesTried, ps3.SingleMaxDataConfidence, ps3.SingleHMACAuthenticated, ps3.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, pq3, key, cw, ch)
											if ps3.SingleHMACAuthenticated {
												report.ThirdPairAuthenticated++
											}
										}
										outTP := ExperimentalV4PhoneBuild63ThirdPairState{ExperimentalV4PhoneBuild53PairState: ps3, ContinuationEvaluations: bp3.evals}
										for _, bc3 := range bp3.states {
											cq3, _, cok3 := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bc3.base.hyp, 0)
											cs3 := experimentalV4PhoneBuild55StateFromQualified(bc3.base, cq3, cok3, seedSource, srcBounds, work)
											if cok3 {
												report.ThirdPairContinuationQualified++
												cs3.SingleDecodeAttempted = true
												cs3.SingleProfilesTried, cs3.SingleListFramesTried, cs3.SingleMaxDataConfidence, cs3.SingleHMACAuthenticated, cs3.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, cq3, key, cw, ch)
												if cs3.SingleHMACAuthenticated {
													report.ThirdPairContinuationAuthenticated++
												}
											}
											outC3 := ExperimentalV4PhoneBuild63ThirdPairContinuationState{ExperimentalV4PhoneBuild55ContinuationState: cs3, SiblingEvaluations: bc3.siblingEvals}
											for sri4, bs3 := range bc3.siblings {
												sq3, _, sok3 := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bs3.base.hyp, 0)
												ss3 := experimentalV4PhoneBuild55SiblingFromQualified(sri4+1, bs3.base, sq3, sok3, seedSource, srcBounds, work)
												if sok3 {
													report.ThirdPairSiblingQualified++
													ss3.SingleDecodeAttempted = true
													ss3.SingleProfilesTried, ss3.SingleListFramesTried, ss3.SingleMaxDataConfidence, ss3.SingleHMACAuthenticated, ss3.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, sq3, key, cw, ch)
													if ss3.SingleHMACAuthenticated {
														report.ThirdPairSiblingAuthenticated++
													}
												}
												outS3 := ExperimentalV4PhoneBuild63ThirdPairSiblingState{Rank: sri4 + 1, Parent: ss3, SingleEvaluations: bs3.singleEvals, SingleImproving: bs3.singleImproving, CoordinateLocal: bs3.singleImproving == 0, PairEvaluations: bs3.pairEvals, PairImproving: bs3.pairImproving}
												for pri4, bp4 := range bs3.pairs {
													pq4, _, pok4 := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bp4.pair.hyp, 0)
													ps4 := experimentalV4PhoneBuild53PairStateFromQualified(pri4+1, bp4.pair, pq4, pok4, seedSource, srcBounds, work)
													if pok4 {
														report.FourthPairQualified++
														ps4.SingleDecodeAttempted = true
														ps4.SingleProfilesTried, ps4.SingleListFramesTried, ps4.SingleMaxDataConfidence, ps4.SingleHMACAuthenticated, ps4.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, pq4, key, cw, ch)
														if ps4.SingleHMACAuthenticated {
															report.FourthPairAuthenticated++
														}
													}
													outP4 := ExperimentalV4PhoneBuild63FourthPairState{ExperimentalV4PhoneBuild53PairState: ps4, ContinuationEvaluations: bp4.evals}
													for _, bc4 := range bp4.states {
														cq4, _, cok4 := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bc4.hyp, 0)
														cs4 := experimentalV4PhoneBuild55StateFromQualified(bc4, cq4, cok4, seedSource, srcBounds, work)
														if cok4 {
															report.FourthPairContinuationQualified++
															cs4.SingleDecodeAttempted = true
															cs4.SingleProfilesTried, cs4.SingleListFramesTried, cs4.SingleMaxDataConfidence, cs4.SingleHMACAuthenticated, cs4.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, cq4, key, cw, ch)
															if cs4.SingleHMACAuthenticated {
																report.FourthPairContinuationAuthenticated++
															}
														}
														outP4.ContinuationStates = append(outP4.ContinuationStates, ExperimentalV4PhoneBuild63FourthPairContinuationState{ExperimentalV4PhoneBuild55ContinuationState: cs4})
													}
													outS3.PairStates = append(outS3.PairStates, outP4)
												}
												outC3.SiblingStates = append(outC3.SiblingStates, outS3)
											}
											outTP.ContinuationStates = append(outTP.ContinuationStates, outC3)
										}
										outS2.PairStates = append(outS2.PairStates, outTP)
									}
									outCont.SiblingStates = append(outCont.SiblingStates, outS2)
								}
								outPair.ContinuationStates = append(outPair.ContinuationStates, outCont)
							}
							outS.PairStates = append(outS.PairStates, outPair)
						}
						outCS.Siblings = append(outCS.Siblings, outS)
					}
					outBranch.ContinuationStates = append(outBranch.ContinuationStates, outCS)
				}
				outRoot.Branches = append(outRoot.Branches, outBranch)
			}
			c.Roots = append(c.Roots, outRoot)
		}
		report.Candidates = append(report.Candidates, c)
	}
	return report, nil
}
