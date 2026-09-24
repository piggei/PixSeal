package watermark

import (
	"errors"
	"image"
)

const (
	experimentalV4PhoneBuild59SeedsPerPair = experimentalV4PhoneBuild56SeedsPerPair
)

// ExperimentalV4PhoneBuild59PairEscapeState is one proposal-improving coupled
// +/-1px two-coordinate escape from a Build55 sibling that is a local maximum
// under the complete independent one-coordinate +/-1px stencil.
type ExperimentalV4PhoneBuild59PairEscapeState struct {
	ExperimentalV4PhoneBuild53PairState
	ContinuationEvaluations int                                                   `json:"continuation_evaluations"`
	ContinuationStates      []ExperimentalV4PhoneBuild59PostPairContinuationState `json:"continuation_states,omitempty"`
}

// ExperimentalV4PhoneBuild59PostPairContinuationState is one proposal-improving 1px
// intermediate retained while continuing from a frozen Build56 second-pair
// escape. Geometry/proposal are created blind. Qualification and HMAC are
// attached only after the complete Build59 bank has been frozen.
type ExperimentalV4PhoneBuild59ThirdPairState struct {
	ExperimentalV4PhoneBuild53PairState
}

// ExperimentalV4PhoneBuild59SecondPairSiblingState preserves one frozen
// Build58 sibling, probes its complete independent one-coordinate +/-1px
// neighborhood with proposal only, and retains bounded coupled pair escapes
// only when that sibling is proposal-local.
type ExperimentalV4PhoneBuild59SecondPairSiblingState struct {
	Rank              int                                        `json:"rank"`
	Parent            ExperimentalV4PhoneBuild55SiblingState     `json:"parent"`
	SingleEvaluations int                                        `json:"single_evaluations"`
	SingleImproving   int                                        `json:"single_improving"`
	CoordinateLocal   bool                                       `json:"coordinate_local"`
	PairEvaluations   int                                        `json:"pair_evaluations"`
	PairImproving     int                                        `json:"pair_improving"`
	PairStates        []ExperimentalV4PhoneBuild59ThirdPairState `json:"pair_states,omitempty"`
}

type ExperimentalV4PhoneBuild59PostPairContinuationState struct {
	ExperimentalV4PhoneBuild55ContinuationState
	SiblingEvaluations int                                                `json:"sibling_evaluations"`
	SiblingStates      []ExperimentalV4PhoneBuild59SecondPairSiblingState `json:"sibling_states,omitempty"`
}

// ExperimentalV4PhoneBuild59SiblingState preserves one frozen Build55 sibling
// and records the proposal-only local-max probe plus any retained pair escapes.
type ExperimentalV4PhoneBuild59SiblingState struct {
	Rank              int                                         `json:"rank"`
	Parent            ExperimentalV4PhoneBuild55SiblingState      `json:"parent"`
	SingleEvaluations int                                         `json:"single_evaluations"`
	SingleImproving   int                                         `json:"single_improving"`
	CoordinateLocal   bool                                        `json:"coordinate_local"`
	PairEvaluations   int                                         `json:"pair_evaluations"`
	PairImproving     int                                         `json:"pair_improving"`
	PairStates        []ExperimentalV4PhoneBuild59PairEscapeState `json:"pair_states,omitempty"`
}

type ExperimentalV4PhoneBuild59ContinuationState struct {
	Parent   ExperimentalV4PhoneBuild55ContinuationState `json:"parent"`
	Siblings []ExperimentalV4PhoneBuild59SiblingState    `json:"siblings,omitempty"`
}

type ExperimentalV4PhoneBuild59Branch struct {
	PairRank           int                                           `json:"pair_rank"`
	ParentPair         ExperimentalV4PhoneBuild53PairState           `json:"parent_pair"`
	ContinuationStates []ExperimentalV4PhoneBuild59ContinuationState `json:"continuation_states,omitempty"`
}

type ExperimentalV4PhoneBuild59Root struct {
	RootIndex         int                                `json:"root_index"`
	Root              ExperimentalV4PhoneBuild52State    `json:"root"`
	SingleEvaluations int                                `json:"single_evaluations"`
	SingleImproving   int                                `json:"single_improving"`
	CoordinateLocal   bool                               `json:"coordinate_local"`
	PairEvaluations   int                                `json:"pair_evaluations"`
	PairImproving     int                                `json:"pair_improving"`
	PairRetained      int                                `json:"pair_retained"`
	Branches          []ExperimentalV4PhoneBuild59Branch `json:"branches,omitempty"`
}

type ExperimentalV4PhoneBuild59Candidate struct {
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
	Roots              []ExperimentalV4PhoneBuild59Root `json:"roots,omitempty"`
}

type ExperimentalV4PhoneBuild59Report struct {
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
	Candidates                             []ExperimentalV4PhoneBuild59Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild59BlindSecondPairSibling struct {
	base            experimentalV4PhoneBuild55BlindSiblingState
	singleEvals     int
	singleImproving int
	pairEvals       int
	pairImproving   int
	pairs           []experimentalV4PhoneBuild53BlindPairState
}

type experimentalV4PhoneBuild59BlindPostPairContinuation struct {
	base         experimentalV4PhoneBuild55BlindContinuationState
	siblingEvals int
	siblings     []experimentalV4PhoneBuild59BlindSecondPairSibling
}

type experimentalV4PhoneBuild59BlindSecondPairBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild59BlindPostPairContinuation
}

type experimentalV4PhoneBuild59BlindSibling struct {
	base            experimentalV4PhoneBuild55BlindSiblingState
	singleEvals     int
	singleImproving int
	pairEvals       int
	pairImproving   int
	pairs           []experimentalV4PhoneBuild59BlindSecondPairBranch
}

type experimentalV4PhoneBuild59BlindContinuation struct {
	base     experimentalV4PhoneBuild55BlindContinuationState
	siblings []experimentalV4PhoneBuild59BlindSibling
}

type experimentalV4PhoneBuild59BlindBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild59BlindContinuation
}

type experimentalV4PhoneBuild59BlindRoot struct {
	root              experimentalV4PhoneBuild52BlindState
	singleEvaluations int
	singleImproving   int
	pairEvaluations   int
	pairImproving     int
	branches          []experimentalV4PhoneBuild59BlindBranch
}

// ExperimentalV4PhoneBuild59Diagnose is research-only. It reproduces the full
// Build58 blind bank. Every frozen Build58 second-pair sibling is then tested
// against the complete independent one-coordinate +/-1px stencil using proposal
// only. Only proposal-local siblings receive the bounded two-coordinate +/-1px
// pair stencil, retaining the unchanged Build53 top-eight policy. The entire
// expanded B/mild+B/angle geometry bank is frozen before qualification, decode,
// HMAC or oracle use.
func ExperimentalV4PhoneBuild59Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild59Report, error) {
	var report ExperimentalV4PhoneBuild59Report
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
	report.SeedsPerPair = experimentalV4PhoneBuild59SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild59SeedsPerPair)
	report.SeedsSelected = len(seeds)

	type blindCandidate struct {
		seed     experimentalV4PhoneBuild48Seed
		baseline experimentalV4PhoneHypothesis
		roots    []experimentalV4PhoneBuild59BlindRoot
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
		bc := blindCandidate{seed: seed, baseline: baseline, roots: make([]experimentalV4PhoneBuild59BlindRoot, 0, len(roots))}
		for _, root := range roots {
			singleEvals, singleImproving := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
			report.SingleEvaluations += singleEvals
			br := experimentalV4PhoneBuild59BlindRoot{root: root, singleEvaluations: singleEvals, singleImproving: singleImproving}
			if singleImproving == 0 {
				report.CoordinateLocalRoots++
				pairs, pairEvals, pairImproving := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
				br.pairEvaluations, br.pairImproving = pairEvals, pairImproving
				report.PairEvaluations += pairEvals
				report.PairStatesFrozen += len(pairs)
				br.branches = make([]experimentalV4PhoneBuild59BlindBranch, 0, len(pairs))
				for _, pair := range pairs {
					cont, contEvals := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
					report.ContinuationEvaluations += contEvals
					report.ContinuationStatesFrozen += len(cont)
					bb := experimentalV4PhoneBuild59BlindBranch{pair: pair, evals: contEvals, states: make([]experimentalV4PhoneBuild59BlindContinuation, 0, len(cont))}
					for _, cs := range cont {
						sibs, sibEvals := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, cs.hyp, 0)
						cs.siblingEvals = sibEvals
						cs.siblings = sibs
						report.SiblingEvaluations += sibEvals
						report.SiblingStatesFrozen += len(sibs)
						outCS := experimentalV4PhoneBuild59BlindContinuation{base: cs, siblings: make([]experimentalV4PhoneBuild59BlindSibling, 0, len(sibs))}
						for _, sib := range sibs {
							se, si := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, sib.hyp, 0)
							report.SiblingLocalProbeEvaluations += se
							outS := experimentalV4PhoneBuild59BlindSibling{base: sib, singleEvals: se, singleImproving: si}
							if si == 0 {
								report.CoordinateLocalSiblings++
								ps, pe, pi := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, sib.hyp, 0)
								outS.pairEvals, outS.pairImproving = pe, pi
								report.SiblingPairEvaluations += pe
								report.SiblingPairImproving += pi
								report.SiblingPairStatesFrozen += len(ps)
								outS.pairs = make([]experimentalV4PhoneBuild59BlindSecondPairBranch, 0, len(ps))
								for _, p := range ps {
									cont2, ce := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, p.hyp, 0)
									report.SecondPairContinuationEvaluations += ce
									report.SecondPairContinuationStatesFrozen += len(cont2)
									outPair := experimentalV4PhoneBuild59BlindSecondPairBranch{pair: p, evals: ce, states: make([]experimentalV4PhoneBuild59BlindPostPairContinuation, 0, len(cont2))}
									for _, c2 := range cont2 {
										sibs2, se2 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c2.hyp, 0)
										report.SecondPairSiblingEvaluations += se2
										report.SecondPairSiblingStatesFrozen += len(sibs2)
										outC2 := experimentalV4PhoneBuild59BlindPostPairContinuation{base: c2, siblingEvals: se2, siblings: make([]experimentalV4PhoneBuild59BlindSecondPairSibling, 0, len(sibs2))}
										for _, s2 := range sibs2 {
											se3, si3 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s2.hyp, 0)
											report.SecondPairSiblingLocalProbeEvaluations += se3
											outS2 := experimentalV4PhoneBuild59BlindSecondPairSibling{base: s2, singleEvals: se3, singleImproving: si3}
											if si3 == 0 {
												report.CoordinateLocalSecondPairSiblings++
												ps3, pe3, pi3 := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s2.hyp, 0)
												outS2.pairs, outS2.pairEvals, outS2.pairImproving = ps3, pe3, pi3
												report.ThirdPairEvaluations += pe3
												report.ThirdPairImproving += pi3
												report.ThirdPairStatesFrozen += len(ps3)
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
	// second-pair sibling are fixed for every seed before qualification/decode.
	srcBounds := src.Bounds()
	for i, bc := range blind {
		seedSource := experimentalV4PhoneBuild46SourceQuad(bc.seed.frozen.h.quad, srcBounds, work)
		c := ExperimentalV4PhoneBuild59Candidate{
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
			outRoot := ExperimentalV4PhoneBuild59Root{RootIndex: ri, Root: root, SingleEvaluations: br.singleEvaluations, SingleImproving: br.singleImproving, CoordinateLocal: br.singleImproving == 0, PairEvaluations: br.pairEvaluations, PairImproving: br.pairImproving, PairRetained: len(br.branches)}
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
				outBranch := ExperimentalV4PhoneBuild59Branch{PairRank: pi + 1, ParentPair: parentPair}
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
					outCS := ExperimentalV4PhoneBuild59ContinuationState{Parent: parentState}
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
						outS := ExperimentalV4PhoneBuild59SiblingState{Rank: sri + 1, Parent: parentSibling, SingleEvaluations: bs.singleEvals, SingleImproving: bs.singleImproving, CoordinateLocal: bs.singleImproving == 0, PairEvaluations: bs.pairEvals, PairImproving: bs.pairImproving}
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
							outPair := ExperimentalV4PhoneBuild59PairEscapeState{ExperimentalV4PhoneBuild53PairState: ps, ContinuationEvaluations: bp.evals}
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
								outCont := ExperimentalV4PhoneBuild59PostPairContinuationState{ExperimentalV4PhoneBuild55ContinuationState: cs, SiblingEvaluations: bcont.siblingEvals}
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
									outS2 := ExperimentalV4PhoneBuild59SecondPairSiblingState{Rank: si + 1, Parent: ss2, SingleEvaluations: bs2.singleEvals, SingleImproving: bs2.singleImproving, CoordinateLocal: bs2.singleImproving == 0, PairEvaluations: bs2.pairEvals, PairImproving: bs2.pairImproving}
									for pri3, bp3 := range bs2.pairs {
										pq3, _, pok3 := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bp3.hyp, 0)
										ps3 := experimentalV4PhoneBuild53PairStateFromQualified(pri3+1, bp3, pq3, pok3, seedSource, srcBounds, work)
										if pok3 {
											report.ThirdPairQualified++
											ps3.SingleDecodeAttempted = true
											ps3.SingleProfilesTried, ps3.SingleListFramesTried, ps3.SingleMaxDataConfidence, ps3.SingleHMACAuthenticated, ps3.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, pq3, key, cw, ch)
											if ps3.SingleHMACAuthenticated {
												report.ThirdPairAuthenticated++
											}
										}
										outS2.PairStates = append(outS2.PairStates, ExperimentalV4PhoneBuild59ThirdPairState{ExperimentalV4PhoneBuild53PairState: ps3})
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
