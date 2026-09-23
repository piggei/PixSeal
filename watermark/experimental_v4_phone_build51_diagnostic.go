package watermark

import (
	"errors"
	"image"
	"math"
)

const experimentalV4PhoneBuild51SeedsPerPair = experimentalV4PhoneBuild50SeedsPerPair

// ExperimentalV4PhoneBuild51TraceStep records one proposal-only candidate
// evaluation performed by the Build41 coordinate-descent refiner. The held-out
// validation score is filled only after the complete trace/stencil bank has
// been frozen. No oracle, key, payload, ECC or HMAC evidence can affect the
// step generation or acceptance decision.
type ExperimentalV4PhoneBuild51TraceStep struct {
	EvalIndex           int           `json:"eval_index"`
	StepSizePx          float64       `json:"step_size_px"`
	Pass                int           `json:"pass"`
	Dimension           int           `json:"dimension"`
	Corner              int           `json:"corner"`
	Axis                string        `json:"axis"`
	DeltaPx             float64       `json:"delta_px"`
	ProposalBefore      float64       `json:"proposal_before"`
	ProposalAfter       float64       `json:"proposal_after"`
	ProposalAvailable   bool          `json:"proposal_available"`
	Accepted            bool          `json:"accepted"`
	ValidationScore     float64       `json:"validation_score"`
	ValidationAvailable bool          `json:"validation_available"`
	WorkingQuadBefore   [4]ImagePoint `json:"working_quad_before"`
	WorkingQuadAfter    [4]ImagePoint `json:"working_quad_after"`
	SourceQuadBefore    [4]ImagePoint `json:"source_quad_before"`
	SourceQuadAfter     [4]ImagePoint `json:"source_quad_after"`
}

// ExperimentalV4PhoneBuild51StencilSample is one deterministic local sample
// around the original proposal seed. The basis/delta list is fixed in code;
// proposal score is computed before any held-out evidence is read.
type ExperimentalV4PhoneBuild51StencilSample struct {
	Index               int           `json:"index"`
	Basis               string        `json:"basis"`
	DeltaPx             float64       `json:"delta_px"`
	ProposalScore       float64       `json:"proposal_score"`
	ValidationScore     float64       `json:"validation_score"`
	ValidationAvailable bool          `json:"validation_available"`
	WorkingQuad         [4]ImagePoint `json:"working_quad"`
	SourceQuad          [4]ImagePoint `json:"source_quad"`
}

// ExperimentalV4PhoneBuild51Candidate contains one top-4 proposal seed, the
// exact Build41 refinement trace, a deterministic local score stencil around
// the unrefined seed, and downstream qualification/HMAC annotations.
type ExperimentalV4PhoneBuild51Candidate struct {
	Index                   int                                       `json:"index"`
	SeedIndex               int                                       `json:"seed_index"`
	SeedRankWithinPair      int                                       `json:"seed_rank_within_pair"`
	SourcePair              string                                    `json:"source_pair,omitempty"`
	SourcePairRank          int                                       `json:"source_pair_rank,omitempty"`
	SourceTier              string                                    `json:"source_tier,omitempty"`
	CellRank                int                                       `json:"cell_rank"`
	PreProposal             float64                                   `json:"pre_proposal"`
	PreValidation           float64                                   `json:"pre_validation"`
	PrePilotScore           float64                                   `json:"pre_pilot_score"`
	PrePilotMargin          float64                                   `json:"pre_pilot_margin"`
	PreQualified            bool                                      `json:"pre_qualified"`
	PostProposal            float64                                   `json:"post_proposal"`
	PostValidation          float64                                   `json:"post_validation"`
	PostPilotScore          float64                                   `json:"post_pilot_score"`
	PostPilotMargin         float64                                   `json:"post_pilot_margin"`
	PostQualified           bool                                      `json:"post_qualified"`
	OriginXBlocks           int                                       `json:"origin_x_blocks"`
	OriginYBlocks           int                                       `json:"origin_y_blocks"`
	WorkingQuadBefore       [4]ImagePoint                             `json:"working_quad_before"`
	WorkingQuadAfter        [4]ImagePoint                             `json:"working_quad_after"`
	SourceQuadBefore        [4]ImagePoint                             `json:"source_quad_before"`
	SourceQuadAfter         [4]ImagePoint                             `json:"source_quad_after"`
	RefineMeanMovementPx    float64                                   `json:"refine_mean_movement_px"`
	RefineMaxMovementPx     float64                                   `json:"refine_max_movement_px"`
	Trace                   []ExperimentalV4PhoneBuild51TraceStep     `json:"trace,omitempty"`
	Stencil                 []ExperimentalV4PhoneBuild51StencilSample `json:"stencil,omitempty"`
	SingleDecodeAttempted   bool                                      `json:"single_decode_attempted"`
	SingleProfilesTried     int                                       `json:"single_profiles_tried,omitempty"`
	SingleListFramesTried   int                                       `json:"single_list_frames_tried,omitempty"`
	SingleMaxDataConfidence float64                                   `json:"single_max_data_confidence,omitempty"`
	SingleHMACAuthenticated bool                                      `json:"single_hmac_authenticated"`
	SingleProfile           Profile                                   `json:"single_profile,omitempty"`
}

type ExperimentalV4PhoneBuild51Report struct {
	WorkingWidth       int                                   `json:"working_width"`
	WorkingHeight      int                                   `json:"working_height"`
	Downsampled        bool                                  `json:"downsampled"`
	BoundaryDetected   bool                                  `json:"boundary_detected"`
	BoundaryConfidence float64                               `json:"boundary_confidence"`
	PairRanking        []ExperimentalV4PhonePairScore        `json:"pair_ranking,omitempty"`
	FrozenCandidates   int                                   `json:"frozen_candidates"`
	SeedsPerPair       int                                   `json:"seeds_per_pair"`
	SeedsSelected      int                                   `json:"seeds_selected"`
	RefineEvaluations  int                                   `json:"refine_evaluations"`
	StencilEvaluations int                                   `json:"stencil_evaluations"`
	PreQualified       int                                   `json:"pre_qualified"`
	PostQualified      int                                   `json:"post_qualified"`
	PostAuthenticated  int                                   `json:"post_authenticated"`
	Candidates         []ExperimentalV4PhoneBuild51Candidate `json:"candidates,omitempty"`
}

func experimentalV4PhoneBuild51ApplyDim(in [4]ImagePoint, dim int, delta float64) [4]ImagePoint {
	out := in
	corner := dim / 2
	if dim%2 == 0 {
		out[corner].X += delta
	} else {
		out[corner].Y += delta
	}
	return out
}

// experimentalV4PhoneBuild51TraceRefine exactly mirrors the Build41
// coordinate-descent proposal objective and update ordering, while recording
// every actually-evaluated +/- move. A move is marked Accepted only when it is
// the best direction for that dimension and becomes the new current state.
func experimentalV4PhoneBuild51TraceRefine(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor, start [4]ImagePoint, heldout int) (experimentalV4PhoneHypothesis, int, []ExperimentalV4PhoneBuild51TraceStep) {
	q := start
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return experimentalV4PhoneHypothesis{}, 0, nil
	}
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, heldout, true, false)
	evals := 1
	if math.IsInf(score, 0) || math.IsNaN(score) {
		return experimentalV4PhoneHypothesis{}, evals, nil
	}
	steps := [...]float64{16, 8, 4, 2, 1}
	trace := make([]ExperimentalV4PhoneBuild51TraceStep, 0, 128)
	for _, step := range steps {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 8; dim++ {
				baseQ, baseScore := q, score
				bestQ, bestH, best := q, h, score
				bestTraceIndex := -1
				for _, sgn := range []float64{-1, 1} {
					delta := sgn * step
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
					axis := "y"
					if dim%2 == 0 {
						axis = "x"
					}
					available := !math.IsInf(ss, 0) && !math.IsNaN(ss)
					rec := ExperimentalV4PhoneBuild51TraceStep{
						EvalIndex: len(trace), StepSizePx: step, Pass: pass + 1, Dimension: dim,
						Corner: dim / 2, Axis: axis, DeltaPx: delta,
						ProposalBefore: baseScore, ProposalAvailable: available,
						WorkingQuadBefore: baseQ, WorkingQuadAfter: qq,
					}
					if available {
						rec.ProposalAfter = ss
					}
					trace = append(trace, rec)
					if available && ss > best+1e-7 {
						bestQ, bestH, best = qq, hh, ss
						bestTraceIndex = len(trace) - 1
					}
				}
				if best > score+1e-7 {
					if bestTraceIndex >= 0 {
						trace[bestTraceIndex].Accepted = true
					}
					q, h, score = bestQ, bestH, best
					improved = true
				}
			}
			if !improved {
				break
			}
		}
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals, trace
}

type experimentalV4PhoneBuild51StencilDef struct {
	name  string
	delta float64
	apply func([4]ImagePoint, float64) [4]ImagePoint
}

func experimentalV4PhoneBuild51Joint(in [4]ImagePoint, basis string, d float64) [4]ImagePoint {
	q := in
	switch basis {
	case "translate-x":
		for i := range q {
			q[i].X += d
		}
	case "translate-y":
		for i := range q {
			q[i].Y += d
		}
	case "scale-x":
		q[0].X -= d
		q[2].X -= d
		q[1].X += d
		q[3].X += d
	case "scale-y":
		q[0].Y -= d
		q[1].Y -= d
		q[2].Y += d
		q[3].Y += d
	case "shear-x":
		q[0].X += d
		q[1].X += d
		q[2].X -= d
		q[3].X -= d
	case "shear-y":
		q[0].Y += d
		q[2].Y += d
		q[1].Y -= d
		q[3].Y -= d
	case "top-width":
		q[0].X -= d
		q[1].X += d
	case "bottom-width":
		q[2].X -= d
		q[3].X += d
	case "left-height":
		q[0].Y -= d
		q[2].Y += d
	case "right-height":
		q[1].Y -= d
		q[3].Y += d
	}
	return q
}

func experimentalV4PhoneBuild51StencilQuads(start [4]ImagePoint) []experimentalV4PhoneBuild51StencilDef {
	defs := make([]experimentalV4PhoneBuild51StencilDef, 0, 53)
	defs = append(defs, experimentalV4PhoneBuild51StencilDef{name: "center", delta: 0, apply: func(q [4]ImagePoint, _ float64) [4]ImagePoint { return q }})
	for dim := 0; dim < 8; dim++ {
		dimCopy := dim
		axis := "y"
		if dim%2 == 0 {
			axis = "x"
		}
		name := "corner-" + string(rune('0'+dim/2)) + "-" + axis
		for _, d := range []float64{-4, -2, 2, 4} {
			defs = append(defs, experimentalV4PhoneBuild51StencilDef{name: name, delta: d, apply: func(q [4]ImagePoint, delta float64) [4]ImagePoint {
				return experimentalV4PhoneBuild51ApplyDim(q, dimCopy, delta)
			}})
		}
	}
	for _, basis := range []string{"translate-x", "translate-y", "scale-x", "scale-y", "shear-x", "shear-y", "top-width", "bottom-width", "left-height", "right-height"} {
		basisCopy := basis
		for _, d := range []float64{-2, 2} {
			defs = append(defs, experimentalV4PhoneBuild51StencilDef{name: basisCopy, delta: d, apply: func(q [4]ImagePoint, delta float64) [4]ImagePoint {
				return experimentalV4PhoneBuild51Joint(q, basisCopy, delta)
			}})
		}
	}
	return defs
}

func experimentalV4PhoneBuild51Stencil(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor, start [4]ImagePoint, heldout int) ([]ExperimentalV4PhoneBuild51StencilSample, int) {
	defs := experimentalV4PhoneBuild51StencilQuads(start)
	out := make([]ExperimentalV4PhoneBuild51StencilSample, 0, len(defs))
	evals := 0
	for _, def := range defs {
		q := def.apply(start, def.delta)
		if def.name != "center" && !experimentalV4PhoneBuild41WithinLimit(q, anchor) {
			continue
		}
		h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
		if !ok {
			continue
		}
		score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, heldout, true, false)
		evals++
		if math.IsInf(score, 0) || math.IsNaN(score) {
			continue
		}
		out = append(out, ExperimentalV4PhoneBuild51StencilSample{Index: len(out), Basis: def.name, DeltaPx: def.delta, ProposalScore: score, WorkingQuad: q})
	}
	return out, evals
}

// ExperimentalV4PhoneBuild51Diagnose freezes a complete proposal-only trace
// and local stencil for the same blind top-4 seeds used by Build50. Only after
// all geometry is fixed are held-out validation, full-pilot qualification and
// diagnostic HMAC annotations evaluated. Reference/SIFT oracle geometry is not
// accepted by this API and remains external/post-hoc in the lab script.

// experimentalV4PhoneBuild51Prepare reproduces only the structural image
// preparation/boundary selection prefix of Build41. It intentionally stops
// before any proposal/held-out pilot search, so Build51 can truthfully freeze
// all trace/stencil geometry before held-out evidence is sampled.
func experimentalV4PhoneBuild51Prepare(src image.Image) (image.Image, PrintBoundaryEstimate, bool) {
	work, down := experimentalV4PhoneResize(src, experimentalV4PhoneMaxDimension)
	boundary := experimentalV4PhoneBuild41ArtworkBoundary(work)
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		boundary = experimentalV4PhoneBoundary(work)
	}
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		boundary = experimentalV4PhoneBuild41ComponentBoundary(work)
	}
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		boundary = experimentalV4ScannerBoundary(work)
	}
	return work, boundary, down
}

func ExperimentalV4PhoneBuild51Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild51Report, error) {
	var report ExperimentalV4PhoneBuild51Report
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
	report.SeedsPerPair = experimentalV4PhoneBuild51SeedsPerPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild51SeedsPerPair)
	report.SeedsSelected = len(seeds)

	// Phase 1: create every trace/stencil geometry using proposal tiles only.
	// No validation, full-pilot, key/HMAC or oracle evidence is read here.
	type blindCandidate struct {
		seed    experimentalV4PhoneBuild48Seed
		pre     experimentalV4PhoneHypothesis
		post    experimentalV4PhoneHypothesis
		trace   []ExperimentalV4PhoneBuild51TraceStep
		stencil []ExperimentalV4PhoneBuild51StencilSample
	}
	blind := make([]blindCandidate, 0, len(seeds))
	for _, seed := range seeds {
		post, n, trace := experimentalV4PhoneBuild51TraceRefine(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		report.RefineEvaluations += n
		stencil, sn := experimentalV4PhoneBuild51Stencil(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		report.StencilEvaluations += sn
		if post.h.h[8] == 0 {
			continue
		}
		post.build43Pair = seed.frozen.h.build43Pair
		post.build43PairRank = seed.frozen.h.build43PairRank
		blind = append(blind, blindCandidate{seed: seed, pre: seed.frozen.h, post: post, trace: trace, stencil: stencil})
	}

	// Phase 2: every geometry is frozen. Held-out validation is now annotated
	// on all trace/stencil samples. Full-pilot qualification/HMAC is kept to
	// the before/after seed states to keep this observability study bounded.
	srcBounds := src.Bounds()
	for i, bc := range blind {
		preQ, _, preOK := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bc.pre, 0)
		postQ, _, postOK := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, bc.post, 0)
		if preOK {
			report.PreQualified++
		}
		if postOK {
			report.PostQualified++
		}
		c := ExperimentalV4PhoneBuild51Candidate{
			Index: i, SeedIndex: bc.seed.index, SeedRankWithinPair: bc.seed.rankWithinPair,
			SourcePair: bc.seed.frozen.h.build43Pair, SourcePairRank: bc.seed.frozen.h.build43PairRank,
			SourceTier: experimentalV4PhoneBuild47TierName(bc.seed.frozen.tier), CellRank: bc.seed.frozen.cellRank,
			PreProposal: preQ.proposal, PreValidation: preQ.validation, PrePilotScore: preQ.detection.Score, PrePilotMargin: preQ.detection.Margin, PreQualified: preOK,
			PostProposal: postQ.proposal, PostValidation: postQ.validation, PostPilotScore: postQ.detection.Score, PostPilotMargin: postQ.detection.Margin, PostQualified: postOK,
			OriginXBlocks: postQ.detection.OriginXBlocks, OriginYBlocks: postQ.detection.OriginYBlocks,
			WorkingQuadBefore: preQ.quad, WorkingQuadAfter: postQ.quad,
			Trace: bc.trace, Stencil: bc.stencil,
		}
		c.SourceQuadBefore = experimentalV4PhoneBuild46SourceQuad(preQ.quad, srcBounds, work)
		c.SourceQuadAfter = experimentalV4PhoneBuild46SourceQuad(postQ.quad, srcBounds, work)
		c.RefineMeanMovementPx, c.RefineMaxMovementPx = experimentalV4PhoneBuild48QuadMovement(c.SourceQuadBefore, c.SourceQuadAfter)
		for ti := range c.Trace {
			t := &c.Trace[ti]
			t.SourceQuadBefore = experimentalV4PhoneBuild46SourceQuad(t.WorkingQuadBefore, srcBounds, work)
			t.SourceQuadAfter = experimentalV4PhoneBuild46SourceQuad(t.WorkingQuadAfter, srcBounds, work)
			hh, ok := experimentalV4PhoneQuadHomography(cw, ch, t.WorkingQuadAfter)
			if ok {
				v, _ := experimentalV4PhoneBuild41FoldScore(plane, pilot, cw, ch, hh, 0, false, false)
				if !math.IsInf(v, 0) && !math.IsNaN(v) {
					t.ValidationScore = v
					t.ValidationAvailable = true
				}
			}
		}
		for si := range c.Stencil {
			s := &c.Stencil[si]
			s.SourceQuad = experimentalV4PhoneBuild46SourceQuad(s.WorkingQuad, srcBounds, work)
			hh, ok := experimentalV4PhoneQuadHomography(cw, ch, s.WorkingQuad)
			if ok {
				v, _ := experimentalV4PhoneBuild41FoldScore(plane, pilot, cw, ch, hh, 0, false, false)
				if !math.IsInf(v, 0) && !math.IsNaN(v) {
					s.ValidationScore = v
					s.ValidationAvailable = true
				}
			}
		}
		if postOK {
			c.SingleDecodeAttempted = true
			c.SingleProfilesTried, c.SingleListFramesTried, c.SingleMaxDataConfidence, c.SingleHMACAuthenticated, c.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, postQ, key, cw, ch)
			if c.SingleHMACAuthenticated {
				report.PostAuthenticated++
			}
		}
		report.Candidates = append(report.Candidates, c)
	}
	return report, nil
}
