package watermark

import (
	"errors"
	"image"
	"math"
	"sort"
)

// ExperimentalV4PhoneBuild49Candidate records proposal-only observables for a
// frozen Build47 candidate. Held-out/public-pilot fields are evaluated only
// after the complete observable bank is frozen and never affect any rank here.
type ExperimentalV4PhoneBuild49Candidate struct {
	Index                  int           `json:"index"`
	SourcePair             string        `json:"source_pair"`
	SourcePairRank         int           `json:"source_pair_rank"`
	SourceTier             string        `json:"source_tier"`
	CellRank               int           `json:"cell_rank"`
	PairRobustScore        float64       `json:"pair_robust_score"`
	CellMeanScore          float64       `json:"cell_mean_score"`
	CellRobustScore        float64       `json:"cell_robust_score"`
	ProposalScore          float64       `json:"proposal_score"`
	Fold1Score             float64       `json:"fold1_score"`
	Fold2Score             float64       `json:"fold2_score"`
	FoldMeanScore          float64       `json:"fold_mean_score"`
	FoldMinScore           float64       `json:"fold_min_score"`
	FoldGap                float64       `json:"fold_gap"`
	BalancedFoldScore      float64       `json:"balanced_fold_score"`
	ProposalTileCount      int           `json:"proposal_tile_count"`
	ProposalTileMean       float64       `json:"proposal_tile_mean"`
	ProposalTileStdDev     float64       `json:"proposal_tile_stddev"`
	ProposalTileMin        float64       `json:"proposal_tile_min"`
	ProposalTileMax        float64       `json:"proposal_tile_max"`
	ProposalTilePositive   float64       `json:"proposal_tile_positive_fraction"`
	TileConsistencyScore   float64       `json:"tile_consistency_score"`
	ProposalRankGlobal     int           `json:"proposal_rank_global"`
	ProposalRankWithinPair int           `json:"proposal_rank_within_pair"`
	FoldMinRankWithinPair  int           `json:"fold_min_rank_within_pair"`
	BalancedRankWithinPair int           `json:"balanced_rank_within_pair"`
	TileRankWithinPair     int           `json:"tile_rank_within_pair"`
	WorkingQuad            [4]ImagePoint `json:"working_quad"`
	SourceQuad             [4]ImagePoint `json:"source_quad"`
	ValidationScore        float64       `json:"validation_score"`
	PilotScore             float64       `json:"pilot_score"`
	PilotMargin            float64       `json:"pilot_margin"`
	OriginXBlocks          int           `json:"origin_x_blocks"`
	OriginYBlocks          int           `json:"origin_y_blocks"`
	Qualified              bool          `json:"qualified"`
}

type ExperimentalV4PhoneBuild49Report struct {
	WorkingWidth       int                                   `json:"working_width"`
	WorkingHeight      int                                   `json:"working_height"`
	Downsampled        bool                                  `json:"downsampled"`
	BoundaryDetected   bool                                  `json:"boundary_detected"`
	BoundaryConfidence float64                               `json:"boundary_confidence"`
	PairRanking        []ExperimentalV4PhonePairScore        `json:"pair_ranking,omitempty"`
	FrozenCandidates   int                                   `json:"frozen_candidates"`
	Candidates         []ExperimentalV4PhoneBuild49Candidate `json:"candidates,omitempty"`
}

func experimentalV4PhoneBuild49TileStats(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography) (count int, mean, stddev, minv, maxv, positive float64) {
	if plane == nil {
		return 0, 0, 0, 0, 0, 0
	}
	bw, bh := cw/blockSize, ch/blockSize
	tilesX, tilesY := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	vals := make([]float64, 0, tilesX*tilesY)
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			if experimentalV4PhoneBuild41FoldForTile(tx, ty) == 0 { // held-out excluded
				continue
			}
			num, den := 0.0, 0.0
			baseX, baseY := tx*experimentalV4TileWidthBlocks, ty*experimentalV4TileHeightBlocks
			for i, pos := range candidate.positions {
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				v, ok := readProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
				if !ok {
					continue
				}
				num += float64(candidate.signs[i]) * v
				den += math.Abs(v)
			}
			if den > 0 {
				vals = append(vals, num/den)
			}
		}
	}
	if len(vals) == 0 {
		return 0, 0, 0, 0, 0, 0
	}
	minv, maxv = vals[0], vals[0]
	pos := 0
	for _, v := range vals {
		mean += v
		if v < minv {
			minv = v
		}
		if v > maxv {
			maxv = v
		}
		if v > 0 {
			pos++
		}
	}
	mean /= float64(len(vals))
	for _, v := range vals {
		d := v - mean
		stddev += d * d
	}
	stddev = math.Sqrt(stddev / float64(len(vals)))
	return len(vals), mean, stddev, minv, maxv, float64(pos) / float64(len(vals))
}

func experimentalV4PhoneBuild49AssignRanks(c []ExperimentalV4PhoneBuild49Candidate) {
	global := make([]int, len(c))
	for i := range c {
		global[i] = i
	}
	sort.SliceStable(global, func(i, j int) bool { return c[global[i]].ProposalScore > c[global[j]].ProposalScore })
	for rank, idx := range global {
		c[idx].ProposalRankGlobal = rank + 1
	}

	groups := map[int][]int{}
	for i := range c {
		groups[c[i].SourcePairRank] = append(groups[c[i].SourcePairRank], i)
	}
	rankBy := func(idxs []int, score func(ExperimentalV4PhoneBuild49Candidate) float64, descending bool, set func(*ExperimentalV4PhoneBuild49Candidate, int)) {
		tmp := append([]int(nil), idxs...)
		sort.SliceStable(tmp, func(i, j int) bool {
			a, b := score(c[tmp[i]]), score(c[tmp[j]])
			if descending {
				return a > b
			}
			return a < b
		})
		for r, idx := range tmp {
			set(&c[idx], r+1)
		}
	}
	for _, idxs := range groups {
		rankBy(idxs, func(x ExperimentalV4PhoneBuild49Candidate) float64 { return x.ProposalScore }, true, func(x *ExperimentalV4PhoneBuild49Candidate, r int) { x.ProposalRankWithinPair = r })
		rankBy(idxs, func(x ExperimentalV4PhoneBuild49Candidate) float64 { return x.FoldMinScore }, true, func(x *ExperimentalV4PhoneBuild49Candidate, r int) { x.FoldMinRankWithinPair = r })
		rankBy(idxs, func(x ExperimentalV4PhoneBuild49Candidate) float64 { return x.BalancedFoldScore }, true, func(x *ExperimentalV4PhoneBuild49Candidate, r int) { x.BalancedRankWithinPair = r })
		rankBy(idxs, func(x ExperimentalV4PhoneBuild49Candidate) float64 { return x.TileConsistencyScore }, true, func(x *ExperimentalV4PhoneBuild49Candidate, r int) { x.TileRankWithinPair = r })
	}
}

// ExperimentalV4PhoneBuild49Diagnose measures proposal-only observables across
// the corrected Build47 extended bank. Geometry is not changed or refined.
// Held-out/public-pilot qualification is evaluated only after all proposal
// observables and ranks are fixed.
func ExperimentalV4PhoneBuild49Diagnose(src image.Image, cw, ch int) (ExperimentalV4PhoneBuild49Report, error) {
	var report ExperimentalV4PhoneBuild49Report
	if src == nil {
		return report, errors.New("nil image")
	}
	if cw < experimentalV4TileWidthBlocks*blockSize || ch < experimentalV4TileHeightBlocks*blockSize || cw%blockSize != 0 || ch%blockSize != 0 {
		return report, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}
	work, boundary, _, _, _, down := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
	report.WorkingWidth, report.WorkingHeight = work.Bounds().Dx(), work.Bounds().Dy()
	report.Downsampled = down
	report.BoundaryDetected, report.BoundaryConfidence = boundary.Detected, boundary.Confidence
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return report, nil
	}

	frozen, pairs, _ := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	report.PairRanking = append([]ExperimentalV4PhonePairScore(nil), pairs...)
	report.FrozenCandidates = len(frozen)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	srcBounds := src.Bounds()

	// Phase 1: proposal-only observability. No held-out/full-pilot evidence.
	report.Candidates = make([]ExperimentalV4PhoneBuild49Candidate, 0, len(frozen))
	for i, fh := range frozen {
		f1 := experimentalV4PhoneBuild43ProposalFold(plane, pilot, cw, ch, fh.h.quad, 1)
		f2 := experimentalV4PhoneBuild43ProposalFold(plane, pilot, cw, ch, fh.h.quad, 2)
		mean := (f1 + f2) / 2
		minFold := math.Min(f1, f2)
		gap := math.Abs(f1 - f2)
		n, tm, ts, tmin, tmax, tpos := experimentalV4PhoneBuild49TileStats(plane, pilot, cw, ch, fh.h.h)
		pairScore := 0.0
		if fh.h.build43PairRank > 0 && fh.h.build43PairRank <= len(pairs) {
			pairScore = pairs[fh.h.build43PairRank-1].Score
		}
		c := ExperimentalV4PhoneBuild49Candidate{
			Index: i, SourcePair: fh.h.build43Pair, SourcePairRank: fh.h.build43PairRank,
			SourceTier: experimentalV4PhoneBuild47TierName(fh.tier), CellRank: fh.cellRank,
			PairRobustScore: pairScore, CellMeanScore: fh.cellMean, CellRobustScore: fh.cellRobust,
			ProposalScore: fh.h.proposal, Fold1Score: f1, Fold2Score: f2, FoldMeanScore: mean,
			FoldMinScore: minFold, FoldGap: gap, BalancedFoldScore: minFold - gap,
			ProposalTileCount: n, ProposalTileMean: tm, ProposalTileStdDev: ts, ProposalTileMin: tmin,
			ProposalTileMax: tmax, ProposalTilePositive: tpos, TileConsistencyScore: tm - ts,
			WorkingQuad: fh.h.quad,
		}
		c.SourceQuad = experimentalV4PhoneBuild46SourceQuad(fh.h.quad, srcBounds, work)
		report.Candidates = append(report.Candidates, c)
	}
	experimentalV4PhoneBuild49AssignRanks(report.Candidates)

	// Phase 2: ranks are frozen. Held-out/full-pilot fields are annotation only.
	for i, fh := range frozen {
		q, _, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, fh.h, 0)
		c := &report.Candidates[i]
		c.ValidationScore, c.PilotScore, c.PilotMargin = q.validation, q.detection.Score, q.detection.Margin
		c.OriginXBlocks, c.OriginYBlocks, c.Qualified = q.detection.OriginXBlocks, q.detection.OriginYBlocks, ok
	}
	return report, nil
}
