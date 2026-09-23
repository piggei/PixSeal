package watermark

import (
	"errors"
	"image"
	"math"
	"sort"
)

var experimentalV4PhoneBuild47Limits = []int{32, 64, 128}

const (
	experimentalV4PhoneBuild47TierProduction = 1
	experimentalV4PhoneBuild47TierDepth      = 2
	experimentalV4PhoneBuild47TierAllPairs   = 3
)

func experimentalV4PhoneBuild47TierName(tier int) string {
	switch tier {
	case experimentalV4PhoneBuild47TierProduction:
		return "production"
	case experimentalV4PhoneBuild47TierDepth:
		return "selected-pair-depth"
	case experimentalV4PhoneBuild47TierAllPairs:
		return "all-pair-extension"
	default:
		return "unknown"
	}
}

// ExperimentalV4PhoneBuild47Candidate is one proposal-only geometry frozen by
// the diagnostic extended bank. Oracle and HMAC fields are populated only after
// the blind bank has been frozen and qualified; they never affect ordering.
type ExperimentalV4PhoneBuild47Candidate struct {
	Index                      int           `json:"index"`
	DiagnosticTier             int           `json:"diagnostic_tier"`
	DiagnosticTierName         string        `json:"diagnostic_tier_name"`
	CellRank                   int           `json:"cell_rank"`
	SourcePair                 string        `json:"source_pair,omitempty"`
	SourcePairRank             int           `json:"source_pair_rank,omitempty"`
	ProposalScore              float64       `json:"proposal_score"`
	ValidationScore            float64       `json:"validation_score"`
	PilotScore                 float64       `json:"pilot_score"`
	PilotMargin                float64       `json:"pilot_margin"`
	OriginXBlocks              int           `json:"origin_x_blocks"`
	OriginYBlocks              int           `json:"origin_y_blocks"`
	Qualified                  bool          `json:"qualified"`
	WorkingQuad                [4]ImagePoint `json:"working_quad"`
	SourceQuad                 [4]ImagePoint `json:"source_quad"`
	SingleDecodeAttempted      bool          `json:"single_decode_attempted"`
	SingleHMACAuthenticated    bool          `json:"single_hmac_authenticated"`
	SingleListFramesTried      int           `json:"single_list_frames_tried,omitempty"`
	OracleCompared             bool          `json:"oracle_compared"`
	OracleMeanCornerErrorPx    float64       `json:"oracle_mean_corner_error_px,omitempty"`
	OracleMaxCornerErrorPx     float64       `json:"oracle_max_corner_error_px,omitempty"`
	OracleMeanCornerErrorRatio float64       `json:"oracle_mean_corner_error_ratio,omitempty"`
}

type ExperimentalV4PhoneBuild47LimitSummary struct {
	Limit                    int     `json:"limit"`
	Stage                    string  `json:"stage"`
	Available                int     `json:"available"`
	Qualified                int     `json:"qualified"`
	ProductionBestIndex      int     `json:"production_best_index"`
	ProductionBestValidation float64 `json:"production_best_validation"`
	OracleNearestIndex       int     `json:"oracle_nearest_index,omitempty"`
	OracleNearestMeanErrorPx float64 `json:"oracle_nearest_mean_error_px,omitempty"`
	OracleBestScoreIndex     int     `json:"oracle_best_score_index,omitempty"`
	OracleBestScore          float64 `json:"oracle_best_score,omitempty"`
}

type ExperimentalV4PhoneBuild47Report struct {
	WorkingWidth       int                                      `json:"working_width"`
	WorkingHeight      int                                      `json:"working_height"`
	Downsampled        bool                                     `json:"downsampled"`
	BoundaryDetected   bool                                     `json:"boundary_detected"`
	BoundaryConfidence float64                                  `json:"boundary_confidence"`
	PairRanking        []ExperimentalV4PhonePairScore           `json:"pair_ranking,omitempty"`
	ProposalCandidates int                                      `json:"proposal_candidates"`
	SuperBankLimit     int                                      `json:"super_bank_limit"`
	FrozenCandidates   int                                      `json:"frozen_candidates"`
	Candidates         []ExperimentalV4PhoneBuild47Candidate    `json:"candidates,omitempty"`
	Limits             []ExperimentalV4PhoneBuild47LimitSummary `json:"limits,omitempty"`
}

type experimentalV4PhoneBuild47Frozen struct {
	h          experimentalV4PhoneHypothesis
	tier       int
	cellRank   int
	cellMean   float64
	cellRobust float64
}

// experimentalV4PhoneBuild47Freeze keeps the exact Build43 production stream
// as tier 1, then appends two diagnostic-only proposal tiers:
//
//	tier 1 / 32: top two side pairs, cell ranks 0 and 3 (Build43 production)
//	tier 2 / 64: same side pairs, deeper cell ranks 1 and 2
//	tier 3 /128: remaining side-pair ranks 3..6, cell ranks 0 and 3
//
// The extended tiers are diagnostic only. They never change Build43 production
// ordering, qualification, quorum, data decode or authentication.
func experimentalV4PhoneBuild47Freeze(work image.Image, boundary PrintBoundaryEstimate, cw, ch, maxFrozen int) ([]experimentalV4PhoneBuild47Frozen, []ExperimentalV4PhonePairScore, int) {
	if work == nil || !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) || maxFrozen <= 0 {
		return nil, nil, 0
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)
	base := experimentalV4PhoneBuild43BaseLines(anchor)
	structural, structOK := experimentalV4PhoneBuild43StructuralLines(work, anchor)
	plane := newPixelPlane(work)
	candidate := experimentalV4Prototype2Candidate()
	type rankedPair struct {
		pair  experimentalV4PhoneBuild43Pair
		score float64
		seed  [4]experimentalV4PhoneBuild43Line
	}
	ranked := make([]rankedPair, 0, len(experimentalV4PhoneBuild43Pairs))
	for _, pair := range experimentalV4PhoneBuild43Pairs {
		seed := experimentalV4PhoneBuild43PairSeed(base, structural, structOK, pair)
		score, _ := experimentalV4PhoneBuild43PairRobustScore(plane, candidate, cw, ch, anchor, seed, pair)
		if !math.IsInf(score, -1) {
			ranked = append(ranked, rankedPair{pair: pair, score: score, seed: seed})
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	scores := make([]ExperimentalV4PhonePairScore, 0, len(ranked))
	for _, rp := range ranked {
		scores = append(scores, ExperimentalV4PhonePairScore{Pair: rp.pair.name, Score: rp.score})
	}

	frozen := make([]experimentalV4PhoneBuild47Frozen, 0, maxFrozen)
	proposalCount := 0
	cellsCache := make(map[int][]experimentalV4PhoneBuild43Cell, len(ranked))
	cellsFor := func(pairRank int) []experimentalV4PhoneBuild43Cell {
		if cells, ok := cellsCache[pairRank]; ok {
			return cells
		}
		rp := ranked[pairRank]
		cells, _ := experimentalV4PhoneBuild43Cells(plane, candidate, cw, ch, anchor, rp.seed, rp.pair)
		cellsCache[pairRank] = cells
		return cells
	}
	appendTier := func(pairStart, pairEnd int, cellRanks []int, tier int) {
		if pairStart < 0 {
			pairStart = 0
		}
		if pairEnd > len(ranked) {
			pairEnd = len(ranked)
		}
		for pairRank := pairStart; pairRank < pairEnd && len(frozen) < maxFrozen; pairRank++ {
			rp := ranked[pairRank]
			cells := cellsFor(pairRank)
			for _, cr := range cellRanks {
				if cr < 0 || cr >= len(cells) || len(frozen) >= maxFrozen {
					continue
				}
				bank, _ := experimentalV4PhoneBuild43BasinCandidates(plane, candidate, cw, ch, anchor, rp.seed, rp.pair, cells[cr])
				for _, h := range bank {
					proposalCount++
					if len(frozen) >= maxFrozen {
						break
					}
					h.build43Pair = rp.pair.name
					h.build43PairRank = pairRank + 1
					frozen = append(frozen, experimentalV4PhoneBuild47Frozen{h: h, tier: tier, cellRank: cr, cellMean: cells[cr].mean, cellRobust: cells[cr].robust})
				}
			}
		}
	}

	productionPairs := experimentalV4PhoneBuild43PairKeep
	if productionPairs > len(ranked) {
		productionPairs = len(ranked)
	}
	// Tier 1 is byte-for-byte the Build43 proposal structure: same selected
	// pairs, same cell ranks and same basin generation order.
	appendTier(0, productionPairs, []int{0, 3}, experimentalV4PhoneBuild47TierProduction)
	// Tier 2 asks whether the selected side pairs contain a better basin in
	// nearby high-ranked cells that production does not currently visit.
	appendTier(0, productionPairs, []int{1, 2}, experimentalV4PhoneBuild47TierDepth)
	// Tier 3 asks whether side-pair ranking removed the useful basin entirely.
	appendTier(productionPairs, len(ranked), []int{0, 3}, experimentalV4PhoneBuild47TierAllPairs)
	return frozen, scores, proposalCount
}

// ExperimentalV4PhoneBuild47Diagnose expands only the diagnostic frozen bank.
// The oracle is post-hoc comparison evidence and never affects generation,
// ordering, qualification, or production decode.
func ExperimentalV4PhoneBuild47Diagnose(src image.Image, key []byte, cw, ch int, oracle *[4]ImagePoint) (ExperimentalV4PhoneBuild47Report, error) {
	var report ExperimentalV4PhoneBuild47Report
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
	maxLimit := experimentalV4PhoneBuild47Limits[len(experimentalV4PhoneBuild47Limits)-1]
	frozen, pairRanking, proposals := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, maxLimit)
	report.PairRanking = pairRanking
	report.ProposalCandidates = proposals
	report.SuperBankLimit = maxLimit
	report.FrozenCandidates = len(frozen)
	plane := newPixelPlane(work)
	pilotCandidate := experimentalV4Prototype2Candidate()
	srcBounds := src.Bounds()
	for i, fh := range frozen {
		h := fh.h
		q, _, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilotCandidate, cw, ch, h, 0)
		c := ExperimentalV4PhoneBuild47Candidate{
			Index: i, DiagnosticTier: fh.tier, DiagnosticTierName: experimentalV4PhoneBuild47TierName(fh.tier), CellRank: fh.cellRank,
			SourcePair: h.build43Pair, SourcePairRank: h.build43PairRank,
			ProposalScore: q.proposal, ValidationScore: q.validation,
			PilotScore: q.detection.Score, PilotMargin: q.detection.Margin,
			OriginXBlocks: q.detection.OriginXBlocks, OriginYBlocks: q.detection.OriginYBlocks,
			Qualified: ok, WorkingQuad: q.quad,
		}
		c.SourceQuad = experimentalV4PhoneBuild46SourceQuad(q.quad, srcBounds, work)
		if ok {
			c.SingleDecodeAttempted = true
			_, c.SingleListFramesTried, _, c.SingleHMACAuthenticated, _ = experimentalV4PhoneBuild46SingleCandidateDecode(plane, q, key, cw, ch)
		}
		if oracle != nil {
			c.OracleCompared = true
			c.OracleMeanCornerErrorPx, c.OracleMaxCornerErrorPx, c.OracleMeanCornerErrorRatio = experimentalV4PhoneBuild46OracleError(c.SourceQuad, *oracle)
		}
		report.Candidates = append(report.Candidates, c)
	}

	for stageIndex, limit := range experimentalV4PhoneBuild47Limits {
		maxTier := stageIndex + 1
		s := ExperimentalV4PhoneBuild47LimitSummary{Limit: limit, Stage: experimentalV4PhoneBuild47TierName(maxTier), ProductionBestIndex: -1, OracleNearestIndex: -1, OracleBestScoreIndex: -1}
		bestVal := math.Inf(-1)
		bestOracle := math.Inf(1)
		bestScore := math.Inf(-1)
		for _, c := range report.Candidates {
			if c.DiagnosticTier > maxTier {
				continue
			}
			s.Available++
			if c.Qualified {
				s.Qualified++
			}
			if c.ValidationScore > bestVal {
				bestVal = c.ValidationScore
				s.ProductionBestIndex = c.Index
				s.ProductionBestValidation = c.ValidationScore
			}
			if c.OracleCompared && c.OracleMeanCornerErrorPx < bestOracle {
				bestOracle = c.OracleMeanCornerErrorPx
				s.OracleNearestIndex = c.Index
				s.OracleNearestMeanErrorPx = c.OracleMeanCornerErrorPx
			}
			if c.OracleCompared {
				score := c.ValidationScore + c.PilotScore + c.PilotMargin
				if score > bestScore {
					bestScore = score
					s.OracleBestScoreIndex = c.Index
					s.OracleBestScore = score
				}
			}
		}
		report.Limits = append(report.Limits, s)
	}
	return report, nil
}
