package watermark

import (
	"errors"
	"image"
	"math"
	"sort"
)

const experimentalV4PhoneBuild48SeedsPerPair = 2
const experimentalV4PhoneBuild50SeedsPerPair = 4

// ExperimentalV4PhoneBuild48Candidate records one proposal-only seed and its
// proposal-only local projective refinement. Held-out qualification and HMAC
// are evaluated only after every refined geometry has been frozen.
type ExperimentalV4PhoneBuild48Candidate struct {
	Index                   int           `json:"index"`
	SeedIndex               int           `json:"seed_index"`
	SeedRankWithinPair      int           `json:"seed_rank_within_pair"`
	SourcePair              string        `json:"source_pair,omitempty"`
	SourcePairRank          int           `json:"source_pair_rank,omitempty"`
	SourceTier              string        `json:"source_tier,omitempty"`
	CellRank                int           `json:"cell_rank"`
	PreProposal             float64       `json:"pre_proposal"`
	PreValidation           float64       `json:"pre_validation"`
	PrePilotScore           float64       `json:"pre_pilot_score"`
	PrePilotMargin          float64       `json:"pre_pilot_margin"`
	PreQualified            bool          `json:"pre_qualified"`
	PostProposal            float64       `json:"post_proposal"`
	PostValidation          float64       `json:"post_validation"`
	PostPilotScore          float64       `json:"post_pilot_score"`
	PostPilotMargin         float64       `json:"post_pilot_margin"`
	PostQualified           bool          `json:"post_qualified"`
	OriginXBlocks           int           `json:"origin_x_blocks"`
	OriginYBlocks           int           `json:"origin_y_blocks"`
	WorkingQuadBefore       [4]ImagePoint `json:"working_quad_before"`
	WorkingQuadAfter        [4]ImagePoint `json:"working_quad_after"`
	SourceQuadBefore        [4]ImagePoint `json:"source_quad_before"`
	SourceQuadAfter         [4]ImagePoint `json:"source_quad_after"`
	RefineMeanMovementPx    float64       `json:"refine_mean_movement_px"`
	RefineMaxMovementPx     float64       `json:"refine_max_movement_px"`
	SingleDecodeAttempted   bool          `json:"single_decode_attempted"`
	SingleProfilesTried     int           `json:"single_profiles_tried,omitempty"`
	SingleListFramesTried   int           `json:"single_list_frames_tried,omitempty"`
	SingleMaxDataConfidence float64       `json:"single_max_data_confidence,omitempty"`
	SingleHMACAuthenticated bool          `json:"single_hmac_authenticated"`
	SingleProfile           Profile       `json:"single_profile,omitempty"`
}

type ExperimentalV4PhoneBuild48Report struct {
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
	PreQualified       int                                   `json:"pre_qualified"`
	PostQualified      int                                   `json:"post_qualified"`
	PostAuthenticated  int                                   `json:"post_authenticated"`
	Candidates         []ExperimentalV4PhoneBuild48Candidate `json:"candidates,omitempty"`
}

type experimentalV4PhoneBuild48Seed struct {
	index          int
	rankWithinPair int
	frozen         experimentalV4PhoneBuild47Frozen
}

// experimentalV4PhoneBuild48SelectSeeds is intentionally proposal-only. It
// never reads validation, full-pilot, HMAC or oracle evidence. At most N seeds
// are retained independently for each side-pair rank.
func experimentalV4PhoneBuild48SelectSeeds(frozen []experimentalV4PhoneBuild47Frozen, perPair int) []experimentalV4PhoneBuild48Seed {
	if perPair <= 0 {
		return nil
	}
	groups := make(map[int][]experimentalV4PhoneBuild48Seed)
	ranks := make([]int, 0, len(experimentalV4PhoneBuild43Pairs))
	seen := make(map[int]bool)
	for i, fh := range frozen {
		rank := fh.h.build43PairRank
		if rank <= 0 {
			continue
		}
		groups[rank] = append(groups[rank], experimentalV4PhoneBuild48Seed{index: i, frozen: fh})
		if !seen[rank] {
			seen[rank] = true
			ranks = append(ranks, rank)
		}
	}
	sort.Ints(ranks)
	out := make([]experimentalV4PhoneBuild48Seed, 0, len(ranks)*perPair)
	for _, rank := range ranks {
		g := groups[rank]
		sort.SliceStable(g, func(i, j int) bool {
			return g[i].frozen.h.proposal > g[j].frozen.h.proposal
		})
		n := perPair
		if len(g) < n {
			n = len(g)
		}
		for i := 0; i < n; i++ {
			g[i].rankWithinPair = i + 1
			out = append(out, g[i])
		}
	}
	return out
}

func experimentalV4PhoneBuild48QuadMovement(a, b [4]ImagePoint) (mean, maxMove float64) {
	for i := range a {
		d := math.Hypot(a[i].X-b[i].X, a[i].Y-b[i].Y)
		mean += d
		if d > maxMove {
			maxMove = d
		}
	}
	mean /= 4
	return mean, maxMove
}

// ExperimentalV4PhoneBuild48Diagnose performs a diagnostic-only local
// projective refinement over proposal-selected seeds from all six side-pair
// ranks. All geometry creation/refinement is proposal-only. The complete set of
// refined candidates is frozen before held-out qualification is evaluated.
func ExperimentalV4PhoneBuild48Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild48Report, error) {
	return experimentalV4PhoneBuild48DiagnoseWithSeedDepth(src, key, cw, ch, experimentalV4PhoneBuild48SeedsPerPair)
}

// ExperimentalV4PhoneBuild50Diagnose repeats the Build48 proposal-only local
// refinement study with a blind top-4-per-side-pair seed depth. It is
// diagnostic-only; production Build43/42 generation, ranking, qualification
// and quorum remain unchanged.
func ExperimentalV4PhoneBuild50Diagnose(src image.Image, key []byte, cw, ch int) (ExperimentalV4PhoneBuild48Report, error) {
	return experimentalV4PhoneBuild48DiagnoseWithSeedDepth(src, key, cw, ch, experimentalV4PhoneBuild50SeedsPerPair)
}

func experimentalV4PhoneBuild48DiagnoseWithSeedDepth(src image.Image, key []byte, cw, ch, perPair int) (ExperimentalV4PhoneBuild48Report, error) {
	var report ExperimentalV4PhoneBuild48Report
	if src == nil {
		return report, errors.New("nil image")
	}
	if len(key) < 8 {
		return report, errors.New("key must contain at least 8 bytes")
	}
	if cw < experimentalV4TileWidthBlocks*blockSize || ch < experimentalV4TileHeightBlocks*blockSize || cw%blockSize != 0 || ch%blockSize != 0 {
		return report, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}

	work, boundary, _, _, _, down := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
	report.WorkingWidth, report.WorkingHeight = work.Bounds().Dx(), work.Bounds().Dy()
	report.Downsampled = down
	report.BoundaryDetected, report.BoundaryConfidence = boundary.Detected, boundary.Confidence
	report.SeedsPerPair = perPair
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilotCandidate := experimentalV4Prototype2Candidate()
	frozen, pairRanking, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairRanking...)
	report.FrozenCandidates = len(frozen)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, perPair)
	report.SeedsSelected = len(seeds)

	// Phase 1: proposal-only refinement for every selected seed. No held-out,
	// full-pilot, key/HMAC or oracle evidence is read in this phase.
	type frozenRefinement struct {
		seed experimentalV4PhoneBuild48Seed
		pre  experimentalV4PhoneHypothesis
		post experimentalV4PhoneHypothesis
	}
	refined := make([]frozenRefinement, 0, len(seeds))
	for _, seed := range seeds {
		post, n := experimentalV4PhoneBuild41Refine(plane, pilotCandidate, cw, ch, anchor, seed.frozen.h.quad, 0)
		report.RefineEvaluations += n
		if post.h.h[8] == 0 {
			continue
		}
		post.build43Pair = seed.frozen.h.build43Pair
		post.build43PairRank = seed.frozen.h.build43PairRank
		refined = append(refined, frozenRefinement{seed: seed, pre: seed.frozen.h, post: post})
	}

	// Phase 2: the refined bank is frozen. Held-out/public-pilot qualification
	// can now observe both before/after geometry, but cannot create or refine it.
	srcBounds := src.Bounds()
	for i, fr := range refined {
		preQ, _, preOK := experimentalV4PhoneBuild41Qualify(plane, work, pilotCandidate, cw, ch, fr.pre, 0)
		postQ, _, postOK := experimentalV4PhoneBuild41Qualify(plane, work, pilotCandidate, cw, ch, fr.post, 0)
		if preOK {
			report.PreQualified++
		}
		if postOK {
			report.PostQualified++
		}
		c := ExperimentalV4PhoneBuild48Candidate{
			Index: i, SeedIndex: fr.seed.index, SeedRankWithinPair: fr.seed.rankWithinPair,
			SourcePair:     fr.seed.frozen.h.build43Pair,
			SourcePairRank: fr.seed.frozen.h.build43PairRank,
			SourceTier:     experimentalV4PhoneBuild47TierName(fr.seed.frozen.tier),
			CellRank:       fr.seed.frozen.cellRank,
			PreProposal:    preQ.proposal, PreValidation: preQ.validation,
			PrePilotScore: preQ.detection.Score, PrePilotMargin: preQ.detection.Margin, PreQualified: preOK,
			PostProposal: postQ.proposal, PostValidation: postQ.validation,
			PostPilotScore: postQ.detection.Score, PostPilotMargin: postQ.detection.Margin, PostQualified: postOK,
			OriginXBlocks: postQ.detection.OriginXBlocks, OriginYBlocks: postQ.detection.OriginYBlocks,
			WorkingQuadBefore: preQ.quad, WorkingQuadAfter: postQ.quad,
		}
		c.SourceQuadBefore = experimentalV4PhoneBuild46SourceQuad(preQ.quad, srcBounds, work)
		c.SourceQuadAfter = experimentalV4PhoneBuild46SourceQuad(postQ.quad, srcBounds, work)
		c.RefineMeanMovementPx, c.RefineMaxMovementPx = experimentalV4PhoneBuild48QuadMovement(c.SourceQuadBefore, c.SourceQuadAfter)
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
