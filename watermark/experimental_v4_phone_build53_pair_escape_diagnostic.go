package watermark

import (
	"errors"
	"image"
	"math"
	"sort"
)

const (
	experimentalV4PhoneBuild53SeedsPerPair      = experimentalV4PhoneBuild52SeedsPerPair
	experimentalV4PhoneBuild53MaxPasses         = 2
	experimentalV4PhoneBuild53Dimensions        = 8
	experimentalV4PhoneBuild53MaxRootsPerSeed   = 1 + experimentalV4PhoneBuild53MaxPasses*experimentalV4PhoneBuild53Dimensions
	experimentalV4PhoneBuild53PairCombosPerRoot = (experimentalV4PhoneBuild53Dimensions * (experimentalV4PhoneBuild53Dimensions - 1) / 2) * 4
	experimentalV4PhoneBuild53PairKeepPerRoot   = 8
)

var (
	experimentalV4PhoneBuild53RootStep  = 2.0
	experimentalV4PhoneBuild53ProbeStep = 1.0
)

// ExperimentalV4PhoneBuild53PairState is one proposal-improving coupled 1px
// move from a 2px root that is a coordinate-local maximum at 1px. Geometry and
// proposal fields are created blind. Validation/full-pilot/HMAC are populated
// only after the complete Build53 bank for all seeds has been frozen.
type ExperimentalV4PhoneBuild53PairState struct {
	Rank                    int           `json:"rank"`
	DimensionA              int           `json:"dimension_a"`
	CornerA                 int           `json:"corner_a"`
	AxisA                   string        `json:"axis_a"`
	DeltaAPx                float64       `json:"delta_a_px"`
	DimensionB              int           `json:"dimension_b"`
	CornerB                 int           `json:"corner_b"`
	AxisB                   string        `json:"axis_b"`
	DeltaBPx                float64       `json:"delta_b_px"`
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

// ExperimentalV4PhoneBuild53Root records one retained proposal-only 2px state.
// Pairwise scanning is enabled only when no +/-1px single-coordinate neighbor
// improves the unchanged Build41 proposal score from this exact root.
type ExperimentalV4PhoneBuild53Root struct {
	RootIndex         int                                   `json:"root_index"`
	Root              ExperimentalV4PhoneBuild52State       `json:"root"`
	SingleEvaluations int                                   `json:"single_evaluations"`
	SingleImproving   int                                   `json:"single_improving"`
	CoordinateLocal   bool                                  `json:"coordinate_local"`
	PairEvaluations   int                                   `json:"pair_evaluations"`
	PairImproving     int                                   `json:"pair_improving"`
	PairRetained      int                                   `json:"pair_retained"`
	PairStates        []ExperimentalV4PhoneBuild53PairState `json:"pair_states,omitempty"`
}

type ExperimentalV4PhoneBuild53Candidate struct {
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
	Roots              []ExperimentalV4PhoneBuild53Root `json:"roots,omitempty"`
}

type ExperimentalV4PhoneBuild53Report struct {
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
	RootEvaluations       int                                   `json:"root_evaluations"`
	SingleEvaluations     int                                   `json:"single_evaluations"`
	PairEvaluations       int                                   `json:"pair_evaluations"`
	RootsFrozen           int                                   `json:"roots_frozen"`
	CoordinateLocalRoots  int                                   `json:"coordinate_local_roots"`
	PairImproving         int                                   `json:"pair_improving"`
	PairStatesFrozen      int                                   `json:"pair_states_frozen"`
	BaselineQualified     int                                   `json:"baseline_qualified"`
	BaselineAuthenticated int                                   `json:"baseline_authenticated"`
	RootsQualified        int                                   `json:"roots_qualified"`
	RootsAuthenticated    int                                   `json:"roots_authenticated"`
	PairQualified         int                                   `json:"pair_qualified"`
	PairAuthenticated     int                                   `json:"pair_authenticated"`
	Candidates            []ExperimentalV4PhoneBuild53Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild53BlindPairState struct {
	dimA, dimB     int
	deltaA, deltaB float64
	hyp            experimentalV4PhoneHypothesis
}

type experimentalV4PhoneBuild53BlindRoot struct {
	root              experimentalV4PhoneBuild52BlindState
	singleEvaluations int
	singleImproving   int
	pairEvaluations   int
	pairImproving     int
	pairs             []experimentalV4PhoneBuild53BlindPairState
}

// experimentalV4PhoneBuild53TwoPxRoots retains the untouched seed plus every
// accepted 2px intermediate produced by the unchanged Build41 coordinate order.
func experimentalV4PhoneBuild53TwoPxRoots(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor, start [4]ImagePoint, heldout int) ([]experimentalV4PhoneBuild52BlindState, int) {
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
	for pass := 0; pass < experimentalV4PhoneBuild53MaxPasses; pass++ {
		improved := false
		for dim := 0; dim < experimentalV4PhoneBuild53Dimensions; dim++ {
			bestQ, bestH, best := q, h, score
			bestDelta := 0.0
			for _, sign := range []float64{-1, 1} {
				delta := sign * experimentalV4PhoneBuild53RootStep
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
					step: experimentalV4PhoneBuild53RootStep, pass: pass + 1, dim: dim, delta: bestDelta,
					hyp: experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score},
				})
			}
		}
		if !improved {
			break
		}
	}
	return states, evals
}

func experimentalV4PhoneBuild53SingleProbe(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, root experimentalV4PhoneHypothesis, heldout int) (int, int) {
	evals, improving := 0, 0
	for dim := 0; dim < experimentalV4PhoneBuild53Dimensions; dim++ {
		for _, sign := range []float64{-1, 1} {
			qq := experimentalV4PhoneBuild51ApplyDim(root.quad, dim, sign*experimentalV4PhoneBuild53ProbeStep)
			if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
				continue
			}
			hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
			if !ok {
				continue
			}
			ss, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
			evals++
			if !math.IsInf(ss, 0) && !math.IsNaN(ss) && ss > root.proposal+1e-7 {
				improving++
			}
		}
	}
	return evals, improving
}

func experimentalV4PhoneBuild53PairStencil(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, root experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild53BlindPairState, int, int) {
	evals := 0
	all := make([]experimentalV4PhoneBuild53BlindPairState, 0, experimentalV4PhoneBuild53PairKeepPerRoot)
	for dimA := 0; dimA < experimentalV4PhoneBuild53Dimensions; dimA++ {
		for dimB := dimA + 1; dimB < experimentalV4PhoneBuild53Dimensions; dimB++ {
			for _, signA := range []float64{-1, 1} {
				for _, signB := range []float64{-1, 1} {
					deltaA := signA * experimentalV4PhoneBuild53ProbeStep
					deltaB := signB * experimentalV4PhoneBuild53ProbeStep
					qq := experimentalV4PhoneBuild51ApplyDim(root.quad, dimA, deltaA)
					qq = experimentalV4PhoneBuild51ApplyDim(qq, dimB, deltaB)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					ss, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, hh, heldout, true, false)
					evals++
					if !math.IsInf(ss, 0) && !math.IsNaN(ss) && ss > root.proposal+1e-7 {
						all = append(all, experimentalV4PhoneBuild53BlindPairState{
							dimA: dimA, dimB: dimB, deltaA: deltaA, deltaB: deltaB,
							hyp: experimentalV4PhoneHypothesis{quad: qq, h: hh, proposal: ss},
						})
					}
				}
			}
		}
	}
	improving := len(all)
	sort.SliceStable(all, func(i, j int) bool {
		if math.Abs(all[i].hyp.proposal-all[j].hyp.proposal) > 1e-12 {
			return all[i].hyp.proposal > all[j].hyp.proposal
		}
		if all[i].dimA != all[j].dimA {
			return all[i].dimA < all[j].dimA
		}
		if all[i].dimB != all[j].dimB {
			return all[i].dimB < all[j].dimB
		}
		if all[i].deltaA != all[j].deltaA {
			return all[i].deltaA < all[j].deltaA
		}
		return all[i].deltaB < all[j].deltaB
	})
	if len(all) > experimentalV4PhoneBuild53PairKeepPerRoot {
		all = all[:experimentalV4PhoneBuild53PairKeepPerRoot]
	}
	return all, evals, improving
}

func experimentalV4PhoneBuild53Axis(dim int) (int, string) {
	corner := dim / 2
	axis := "y"
	if dim%2 == 0 {
		axis = "x"
	}
	return corner, axis
}

func experimentalV4PhoneBuild53PairStateFromQualified(rank int, blind experimentalV4PhoneBuild53BlindPairState, qualified experimentalV4PhoneHypothesis, ok bool, seedSource [4]ImagePoint, srcBounds image.Rectangle, work image.Image) ExperimentalV4PhoneBuild53PairState {
	cornerA, axisA := experimentalV4PhoneBuild53Axis(blind.dimA)
	cornerB, axisB := experimentalV4PhoneBuild53Axis(blind.dimB)
	st := ExperimentalV4PhoneBuild53PairState{
		Rank:       rank,
		DimensionA: blind.dimA, CornerA: cornerA, AxisA: axisA, DeltaAPx: blind.deltaA,
		DimensionB: blind.dimB, CornerB: cornerB, AxisB: axisB, DeltaBPx: blind.deltaB,
		Proposal: qualified.proposal, Validation: qualified.validation, PilotScore: qualified.detection.Score, PilotMargin: qualified.detection.Margin,
		Qualified: ok, OriginXBlocks: qualified.detection.OriginXBlocks, OriginYBlocks: qualified.detection.OriginYBlocks,
		WorkingQuad: qualified.quad,
	}
	st.SourceQuad = experimentalV4PhoneBuild46SourceQuad(qualified.quad, srcBounds, work)
	st.MeanMovementFromSeedPx, st.MaxMovementFromSeedPx = experimentalV4PhoneBuild48QuadMovement(seedSource, st.SourceQuad)
	return st
}

// ExperimentalV4PhoneBuild53Diagnose is research-only. It keeps the unchanged
// top-4 seed bank, derives proposal-only 2px roots, identifies 1px coordinate
// local maxima using only the proposal fold, and scans a bounded coupled +/-1px
// two-coordinate neighborhood only around those roots. At most eight pairwise
// proposal-improving states per root are retained, ordered by proposal alone.
// Every geometry is frozen before qualification, protected-data decode or HMAC.
func ExperimentalV4PhoneBuild53Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild53Report, error) {
	var report ExperimentalV4PhoneBuild53Report
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
	report.SeedsPerPair = experimentalV4PhoneBuild53SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild53SeedsPerPair)
	report.SeedsSelected = len(seeds)

	type blindCandidate struct {
		seed     experimentalV4PhoneBuild48Seed
		baseline experimentalV4PhoneHypothesis
		roots    []experimentalV4PhoneBuild53BlindRoot
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
		bc := blindCandidate{seed: seed, baseline: baseline, roots: make([]experimentalV4PhoneBuild53BlindRoot, 0, len(roots))}
		for _, root := range roots {
			singleEvals, singleImproving := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
			report.SingleEvaluations += singleEvals
			br := experimentalV4PhoneBuild53BlindRoot{root: root, singleEvaluations: singleEvals, singleImproving: singleImproving}
			if singleImproving == 0 {
				report.CoordinateLocalRoots++
				pairs, pairEvals, pairImproving := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
				br.pairs, br.pairEvaluations, br.pairImproving = pairs, pairEvals, pairImproving
				report.PairEvaluations += pairEvals
				report.PairImproving += pairImproving
				report.PairStatesFrozen += len(pairs)
			}
			report.RootsFrozen++
			bc.roots = append(bc.roots, br)
		}
		blind = append(blind, bc)
	}

	// Freeze barrier: held-out/full-pilot/key evidence is unavailable until all
	// roots and retained pairwise geometries for every seed are fixed.
	srcBounds := src.Bounds()
	for i, bc := range blind {
		seedSource := experimentalV4PhoneBuild46SourceQuad(bc.seed.frozen.h.quad, srcBounds, work)
		c := ExperimentalV4PhoneBuild53Candidate{
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
			outRoot := ExperimentalV4PhoneBuild53Root{
				RootIndex: ri, Root: root, SingleEvaluations: br.singleEvaluations, SingleImproving: br.singleImproving,
				CoordinateLocal: br.singleImproving == 0, PairEvaluations: br.pairEvaluations, PairImproving: br.pairImproving, PairRetained: len(br.pairs),
			}
			for pi, bp := range br.pairs {
				q, _, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bp.hyp, 0)
				st := experimentalV4PhoneBuild53PairStateFromQualified(pi+1, bp, q, ok, seedSource, srcBounds, work)
				if ok {
					report.PairQualified++
					st.SingleDecodeAttempted = true
					st.SingleProfilesTried, st.SingleListFramesTried, st.SingleMaxDataConfidence, st.SingleHMACAuthenticated, st.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, q, key, cw, ch)
					if st.SingleHMACAuthenticated {
						report.PairAuthenticated++
					}
				}
				outRoot.PairStates = append(outRoot.PairStates, st)
			}
			c.Roots = append(c.Roots, outRoot)
		}
		report.Candidates = append(report.Candidates, c)
	}
	return report, nil
}
