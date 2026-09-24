package watermark

import (
	"errors"
	"image"
)

const (
	experimentalV4PhoneBuild56SeedsPerPair = experimentalV4PhoneBuild55SeedsPerPair
)

// ExperimentalV4PhoneBuild56PairEscapeState is one proposal-improving coupled
// +/-1px two-coordinate escape from a Build55 sibling that is a local maximum
// under the complete independent one-coordinate +/-1px stencil.
type ExperimentalV4PhoneBuild56PairEscapeState struct {
	ExperimentalV4PhoneBuild53PairState
}

// ExperimentalV4PhoneBuild56SiblingState preserves one frozen Build55 sibling
// and records the proposal-only local-max probe plus any retained pair escapes.
type ExperimentalV4PhoneBuild56SiblingState struct {
	Rank              int                                         `json:"rank"`
	Parent            ExperimentalV4PhoneBuild55SiblingState      `json:"parent"`
	SingleEvaluations int                                         `json:"single_evaluations"`
	SingleImproving   int                                         `json:"single_improving"`
	CoordinateLocal   bool                                        `json:"coordinate_local"`
	PairEvaluations   int                                         `json:"pair_evaluations"`
	PairImproving     int                                         `json:"pair_improving"`
	PairStates        []ExperimentalV4PhoneBuild56PairEscapeState `json:"pair_states,omitempty"`
}

type ExperimentalV4PhoneBuild56ContinuationState struct {
	Parent   ExperimentalV4PhoneBuild55ContinuationState `json:"parent"`
	Siblings []ExperimentalV4PhoneBuild56SiblingState    `json:"siblings,omitempty"`
}

type ExperimentalV4PhoneBuild56Branch struct {
	PairRank           int                                           `json:"pair_rank"`
	ParentPair         ExperimentalV4PhoneBuild53PairState           `json:"parent_pair"`
	ContinuationStates []ExperimentalV4PhoneBuild56ContinuationState `json:"continuation_states,omitempty"`
}

type ExperimentalV4PhoneBuild56Root struct {
	RootIndex         int                                `json:"root_index"`
	Root              ExperimentalV4PhoneBuild52State    `json:"root"`
	SingleEvaluations int                                `json:"single_evaluations"`
	SingleImproving   int                                `json:"single_improving"`
	CoordinateLocal   bool                               `json:"coordinate_local"`
	PairEvaluations   int                                `json:"pair_evaluations"`
	PairImproving     int                                `json:"pair_improving"`
	PairRetained      int                                `json:"pair_retained"`
	Branches          []ExperimentalV4PhoneBuild56Branch `json:"branches,omitempty"`
}

type ExperimentalV4PhoneBuild56Candidate struct {
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
	Roots              []ExperimentalV4PhoneBuild56Root `json:"roots,omitempty"`
}

type ExperimentalV4PhoneBuild56Report struct {
	WorkingWidth                 int                                   `json:"working_width"`
	WorkingHeight                int                                   `json:"working_height"`
	Downsampled                  bool                                  `json:"downsampled"`
	BoundaryDetected             bool                                  `json:"boundary_detected"`
	BoundaryConfidence           float64                               `json:"boundary_confidence"`
	PairRanking                  []ExperimentalV4PhonePairScore        `json:"pair_ranking,omitempty"`
	FrozenCandidates             int                                   `json:"frozen_candidates"`
	SeedsPerPair                 int                                   `json:"seeds_per_pair"`
	SeedsSelected                int                                   `json:"seeds_selected"`
	BaselineEvaluations          int                                   `json:"baseline_evaluations"`
	RootEvaluations              int                                   `json:"root_evaluations"`
	SingleEvaluations            int                                   `json:"single_evaluations"`
	PairEvaluations              int                                   `json:"pair_evaluations"`
	ContinuationEvaluations      int                                   `json:"continuation_evaluations"`
	SiblingEvaluations           int                                   `json:"sibling_evaluations"`
	SiblingLocalProbeEvaluations int                                   `json:"sibling_local_probe_evaluations"`
	SiblingPairEvaluations       int                                   `json:"sibling_pair_evaluations"`
	RootsFrozen                  int                                   `json:"roots_frozen"`
	CoordinateLocalRoots         int                                   `json:"coordinate_local_roots"`
	PairStatesFrozen             int                                   `json:"pair_states_frozen"`
	ContinuationStatesFrozen     int                                   `json:"continuation_states_frozen"`
	SiblingStatesFrozen          int                                   `json:"sibling_states_frozen"`
	CoordinateLocalSiblings      int                                   `json:"coordinate_local_siblings"`
	SiblingPairImproving         int                                   `json:"sibling_pair_improving"`
	SiblingPairStatesFrozen      int                                   `json:"sibling_pair_states_frozen"`
	BaselineQualified            int                                   `json:"baseline_qualified"`
	BaselineAuthenticated        int                                   `json:"baseline_authenticated"`
	RootsQualified               int                                   `json:"roots_qualified"`
	RootsAuthenticated           int                                   `json:"roots_authenticated"`
	PairQualified                int                                   `json:"pair_qualified"`
	PairAuthenticated            int                                   `json:"pair_authenticated"`
	ContinuationQualified        int                                   `json:"continuation_qualified"`
	ContinuationAuthenticated    int                                   `json:"continuation_authenticated"`
	SiblingQualified             int                                   `json:"sibling_qualified"`
	SiblingAuthenticated         int                                   `json:"sibling_authenticated"`
	SiblingPairQualified         int                                   `json:"sibling_pair_qualified"`
	SiblingPairAuthenticated     int                                   `json:"sibling_pair_authenticated"`
	Candidates                   []ExperimentalV4PhoneBuild56Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild56BlindSibling struct {
	base            experimentalV4PhoneBuild55BlindSiblingState
	singleEvals     int
	singleImproving int
	pairEvals       int
	pairImproving   int
	pairs           []experimentalV4PhoneBuild53BlindPairState
}

type experimentalV4PhoneBuild56BlindContinuation struct {
	base     experimentalV4PhoneBuild55BlindContinuationState
	siblings []experimentalV4PhoneBuild56BlindSibling
}

type experimentalV4PhoneBuild56BlindBranch struct {
	pair   experimentalV4PhoneBuild53BlindPairState
	evals  int
	states []experimentalV4PhoneBuild56BlindContinuation
}

type experimentalV4PhoneBuild56BlindRoot struct {
	root              experimentalV4PhoneBuild52BlindState
	singleEvaluations int
	singleImproving   int
	pairEvaluations   int
	pairImproving     int
	branches          []experimentalV4PhoneBuild56BlindBranch
}

// ExperimentalV4PhoneBuild56Diagnose is research-only. It reproduces the full
// Build55 blind continuation+sibling bank. Every frozen Build55 sibling is then
// tested against the complete independent one-coordinate +/-1px stencil using
// proposal only. Only proposal-local sibling states receive the bounded Build53
// two-coordinate +/-1px pair stencil, retaining at most the same eight states
// ordered by proposal. The entire expanded geometry bank is frozen before any
// held-out/full-pilot qualification, protected-data decode, HMAC or oracle use.
func ExperimentalV4PhoneBuild56Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild56Report, error) {
	var report ExperimentalV4PhoneBuild56Report
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
	report.SeedsPerPair = experimentalV4PhoneBuild56SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild56SeedsPerPair)
	report.SeedsSelected = len(seeds)

	type blindCandidate struct {
		seed     experimentalV4PhoneBuild48Seed
		baseline experimentalV4PhoneHypothesis
		roots    []experimentalV4PhoneBuild56BlindRoot
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
		bc := blindCandidate{seed: seed, baseline: baseline, roots: make([]experimentalV4PhoneBuild56BlindRoot, 0, len(roots))}
		for _, root := range roots {
			singleEvals, singleImproving := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
			report.SingleEvaluations += singleEvals
			br := experimentalV4PhoneBuild56BlindRoot{root: root, singleEvaluations: singleEvals, singleImproving: singleImproving}
			if singleImproving == 0 {
				report.CoordinateLocalRoots++
				pairs, pairEvals, pairImproving := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
				br.pairEvaluations, br.pairImproving = pairEvals, pairImproving
				report.PairEvaluations += pairEvals
				report.PairStatesFrozen += len(pairs)
				br.branches = make([]experimentalV4PhoneBuild56BlindBranch, 0, len(pairs))
				for _, pair := range pairs {
					cont, contEvals := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
					report.ContinuationEvaluations += contEvals
					report.ContinuationStatesFrozen += len(cont)
					bb := experimentalV4PhoneBuild56BlindBranch{pair: pair, evals: contEvals, states: make([]experimentalV4PhoneBuild56BlindContinuation, 0, len(cont))}
					for _, cs := range cont {
						sibs, sibEvals := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, cs.hyp, 0)
						cs.siblingEvals = sibEvals
						cs.siblings = sibs
						report.SiblingEvaluations += sibEvals
						report.SiblingStatesFrozen += len(sibs)
						outCS := experimentalV4PhoneBuild56BlindContinuation{base: cs, siblings: make([]experimentalV4PhoneBuild56BlindSibling, 0, len(sibs))}
						for _, sib := range sibs {
							se, si := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, sib.hyp, 0)
							report.SiblingLocalProbeEvaluations += se
							outS := experimentalV4PhoneBuild56BlindSibling{base: sib, singleEvals: se, singleImproving: si}
							if si == 0 {
								report.CoordinateLocalSiblings++
								ps, pe, pi := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, sib.hyp, 0)
								outS.pairs, outS.pairEvals, outS.pairImproving = ps, pe, pi
								report.SiblingPairEvaluations += pe
								report.SiblingPairImproving += pi
								report.SiblingPairStatesFrozen += len(ps)
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

	// Freeze barrier: all roots, first pair escapes, continuations, siblings and
	// second pair escapes for every seed are fixed before qualification/decode.
	srcBounds := src.Bounds()
	for i, bc := range blind {
		seedSource := experimentalV4PhoneBuild46SourceQuad(bc.seed.frozen.h.quad, srcBounds, work)
		c := ExperimentalV4PhoneBuild56Candidate{
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
			outRoot := ExperimentalV4PhoneBuild56Root{RootIndex: ri, Root: root, SingleEvaluations: br.singleEvaluations, SingleImproving: br.singleImproving, CoordinateLocal: br.singleImproving == 0, PairEvaluations: br.pairEvaluations, PairImproving: br.pairImproving, PairRetained: len(br.branches)}
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
				outBranch := ExperimentalV4PhoneBuild56Branch{PairRank: pi + 1, ParentPair: parentPair}
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
					outCS := ExperimentalV4PhoneBuild56ContinuationState{Parent: parentState}
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
						outS := ExperimentalV4PhoneBuild56SiblingState{Rank: sri + 1, Parent: parentSibling, SingleEvaluations: bs.singleEvals, SingleImproving: bs.singleImproving, CoordinateLocal: bs.singleImproving == 0, PairEvaluations: bs.pairEvals, PairImproving: bs.pairImproving}
						for pri, bp := range bs.pairs {
							pq, _, pok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bp.hyp, 0)
							ps := experimentalV4PhoneBuild53PairStateFromQualified(pri+1, bp, pq, pok, seedSource, srcBounds, work)
							if pok {
								report.SiblingPairQualified++
								ps.SingleDecodeAttempted = true
								ps.SingleProfilesTried, ps.SingleListFramesTried, ps.SingleMaxDataConfidence, ps.SingleHMACAuthenticated, ps.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, pq, key, cw, ch)
								if ps.SingleHMACAuthenticated {
									report.SiblingPairAuthenticated++
								}
							}
							outS.PairStates = append(outS.PairStates, ExperimentalV4PhoneBuild56PairEscapeState{ExperimentalV4PhoneBuild53PairState: ps})
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
