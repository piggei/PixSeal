package watermark

import (
	"errors"
	"image"
	"math"
)

const (
	experimentalV4PhoneBuild46DirectEnsembleRequired = experimentalV4PhoneEnsembleSize
	experimentalV4PhoneBuild46Build42BankRequired    = 3
)

// ExperimentalV4PhoneBuild46Classification describes the diagnostic-only
// interpretation of the Build43 -> production ensemble -> Build42 handoff.
// It never changes the production decoder.
type ExperimentalV4PhoneBuild46Classification string

const (
	ExperimentalV4PhoneBuild46Build41Accepted            ExperimentalV4PhoneBuild46Classification = "build41-accepted"
	ExperimentalV4PhoneBuild46NoBuild43Candidate         ExperimentalV4PhoneBuild46Classification = "no-build43-candidate"
	ExperimentalV4PhoneBuild46HeldoutQualification       ExperimentalV4PhoneBuild46Classification = "heldout-qualification"
	ExperimentalV4PhoneBuild46QualifiedEnsembleShortfall ExperimentalV4PhoneBuild46Classification = "qualified-ensemble-shortfall"
	ExperimentalV4PhoneBuild46QualifiedCandidateMismatch ExperimentalV4PhoneBuild46Classification = "qualified-candidate-mismatch"
	ExperimentalV4PhoneBuild46Build42BankShortfall       ExperimentalV4PhoneBuild46Classification = "build42-bank-shortfall"
	ExperimentalV4PhoneBuild46QualifiedHandoffAvailable  ExperimentalV4PhoneBuild46Classification = "qualified-handoff-available"
)

// ExperimentalV4PhoneBuild46ClassifyCounts is a pure diagnostic helper used by
// regressions and reporting. Production decoding never consults this label.
func ExperimentalV4PhoneBuild46ClassifyCounts(build41Direct, frozen, qualified int, singleAuthenticated bool) ExperimentalV4PhoneBuild46Classification {
	switch {
	case build41Direct >= experimentalV4PhoneBuild46DirectEnsembleRequired:
		return ExperimentalV4PhoneBuild46Build41Accepted
	case frozen == 0:
		return ExperimentalV4PhoneBuild46NoBuild43Candidate
	case qualified == 0:
		return ExperimentalV4PhoneBuild46HeldoutQualification
	case qualified < experimentalV4PhoneBuild46DirectEnsembleRequired && singleAuthenticated:
		return ExperimentalV4PhoneBuild46QualifiedEnsembleShortfall
	case qualified < experimentalV4PhoneBuild46DirectEnsembleRequired:
		return ExperimentalV4PhoneBuild46QualifiedCandidateMismatch
	case qualified < experimentalV4PhoneBuild46Build42BankRequired:
		return ExperimentalV4PhoneBuild46Build42BankShortfall
	default:
		return ExperimentalV4PhoneBuild46QualifiedHandoffAvailable
	}
}

// ExperimentalV4PhoneBuild46Candidate describes one geometry that already
// passed Build43 held-out public-pilot qualification. Single-candidate data
// decoding is laboratory evidence only: production still requires its existing
// ensemble/bank thresholds and is unchanged by Build46.
type ExperimentalV4PhoneBuild46Candidate struct {
	Index                      int           `json:"index"`
	SourcePair                 string        `json:"source_pair,omitempty"`
	SourcePairRank             int           `json:"source_pair_rank,omitempty"`
	ProposalScore              float64       `json:"proposal_score"`
	ValidationScore            float64       `json:"validation_score"`
	PilotScore                 float64       `json:"pilot_score"`
	PilotMargin                float64       `json:"pilot_margin"`
	OriginXBlocks              int           `json:"origin_x_blocks"`
	OriginYBlocks              int           `json:"origin_y_blocks"`
	WorkingQuad                [4]ImagePoint `json:"working_quad"`
	SourceQuad                 [4]ImagePoint `json:"source_quad"`
	SingleProfilesTried        int           `json:"single_profiles_tried"`
	SingleListFramesTried      int           `json:"single_list_frames_tried"`
	SingleMaxDataConfidence    float64       `json:"single_max_data_confidence"`
	SingleHMACAuthenticated    bool          `json:"single_hmac_authenticated"`
	SingleProfile              Profile       `json:"single_profile,omitempty"`
	OracleCompared             bool          `json:"oracle_compared"`
	OracleMeanCornerErrorPx    float64       `json:"oracle_mean_corner_error_px,omitempty"`
	OracleMaxCornerErrorPx     float64       `json:"oracle_max_corner_error_px,omitempty"`
	OracleMeanCornerErrorRatio float64       `json:"oracle_mean_corner_error_ratio,omitempty"`
}

// ExperimentalV4PhoneBuild46Report exposes why already-qualified Build43
// geometry does or does not reach the production data paths. Oracle geometry,
// when supplied, is used only after blind candidates have been generated and
// qualified, exclusively for distance measurement.
type ExperimentalV4PhoneBuild46Report struct {
	WorkingWidth                 int                                      `json:"working_width"`
	WorkingHeight                int                                      `json:"working_height"`
	Downsampled                  bool                                     `json:"downsampled"`
	BoundaryDetected             bool                                     `json:"boundary_detected"`
	BoundaryConfidence           float64                                  `json:"boundary_confidence"`
	Build41DirectCandidates      int                                      `json:"build41_direct_candidates"`
	Build41QualifiedBank         int                                      `json:"build41_qualified_bank"`
	Build43Attempted             bool                                     `json:"build43_attempted"`
	Build43FrozenCandidates      int                                      `json:"build43_frozen_candidates"`
	Build43QualifiedCandidates   int                                      `json:"build43_qualified_candidates"`
	Build43PairRanking           []ExperimentalV4PhonePairScore           `json:"build43_pair_ranking,omitempty"`
	DirectEnsembleRequired       int                                      `json:"direct_ensemble_required"`
	DirectEnsembleAvailable      bool                                     `json:"direct_ensemble_available"`
	Build42BankRequired          int                                      `json:"build42_bank_required"`
	Build42BankAvailable         bool                                     `json:"build42_bank_available"`
	SingleCandidateAuthenticated bool                                     `json:"single_candidate_authenticated"`
	Classification               ExperimentalV4PhoneBuild46Classification `json:"classification"`
	Candidates                   []ExperimentalV4PhoneBuild46Candidate    `json:"candidates,omitempty"`
}

func experimentalV4PhoneBuild46SourceQuad(q [4]ImagePoint, srcBounds image.Rectangle, work image.Image) [4]ImagePoint {
	out := q
	sx := float64(srcBounds.Dx()) / float64(work.Bounds().Dx())
	sy := float64(srcBounds.Dy()) / float64(work.Bounds().Dy())
	for i := range out {
		out[i].X = float64(srcBounds.Min.X) + q[i].X*sx
		out[i].Y = float64(srcBounds.Min.Y) + q[i].Y*sy
	}
	return out
}

func experimentalV4PhoneBuild46OracleError(q, oracle [4]ImagePoint) (mean, maxErr, ratio float64) {
	for i := range q {
		d := math.Hypot(q[i].X-oracle[i].X, q[i].Y-oracle[i].Y)
		mean += d
		if d > maxErr {
			maxErr = d
		}
	}
	mean /= 4
	d1 := math.Hypot(oracle[3].X-oracle[0].X, oracle[3].Y-oracle[0].Y)
	d2 := math.Hypot(oracle[2].X-oracle[1].X, oracle[2].Y-oracle[1].Y)
	diag := (d1 + d2) / 2
	if diag > 0 {
		ratio = mean / diag
	}
	return mean, maxErr, ratio
}

func experimentalV4PhoneBuild46SingleCandidateDecode(plane *pixelPlane, h experimentalV4PhoneHypothesis, key []byte, cw, ch int) (profiles, frames int, maxConfidence float64, authenticated bool, profile Profile) {
	for _, spec := range v3Profiles {
		var margins []float64
		var confidence float64
		var ok bool
		if h.warp != nil {
			margins, confidence, ok = experimentalV4ReadProtectedPhoneWarpMargins(plane, cw, ch, h.h, h.warp, spec.codedBits)
		} else {
			margins, confidence, ok = experimentalV4ReadProtectedProjectiveMargins(plane, cw, ch, h.h, spec.codedBits)
		}
		if !ok {
			continue
		}
		profiles++
		if confidence > maxConfidence {
			maxConfidence = confidence
		}
		_, tried, auth := experimentalV4PhoneBuild42ListDecode(margins, key, spec)
		frames += tried
		if auth {
			return profiles, frames, maxConfidence, true, spec.profile
		}
	}
	return profiles, frames, maxConfidence, false, ""
}

// ExperimentalV4PhoneBuild46Diagnose reruns only public-evidence Build41/43
// geometry search, then inspects the already-qualified Build43 candidates. The
// optional oracle quad is never used to create, refine, rank or qualify blind
// geometry. Key/HMAC evidence is consulted only in the explicitly diagnostic
// single-candidate data test after geometry has been frozen and qualified.
func ExperimentalV4PhoneBuild46Diagnose(src image.Image, key []byte, cw, ch int, oracle *[4]ImagePoint) (ExperimentalV4PhoneBuild46Report, error) {
	var report ExperimentalV4PhoneBuild46Report
	if src == nil {
		return report, errors.New("nil image")
	}
	if len(key) < 8 {
		return report, errors.New("key must contain at least 8 bytes")
	}
	if cw < experimentalV4TileWidthBlocks*blockSize || ch < experimentalV4TileHeightBlocks*blockSize || cw%blockSize != 0 || ch%blockSize != 0 {
		return report, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}

	work, boundary, direct41, bank41, _, down := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
	report.WorkingWidth = work.Bounds().Dx()
	report.WorkingHeight = work.Bounds().Dy()
	report.Downsampled = down
	report.BoundaryDetected = boundary.Detected
	report.BoundaryConfidence = boundary.Confidence
	report.Build41DirectCandidates = len(direct41)
	report.Build41QualifiedBank = len(bank41)
	report.DirectEnsembleRequired = experimentalV4PhoneBuild46DirectEnsembleRequired
	report.Build42BankRequired = experimentalV4PhoneBuild46Build42BankRequired

	if len(direct41) >= experimentalV4PhoneEnsembleSize {
		report.DirectEnsembleAvailable = true
		report.Build42BankAvailable = len(bank41) >= experimentalV4PhoneBuild46Build42BankRequired
		report.Classification = ExperimentalV4PhoneBuild46Build41Accepted
		return report, nil
	}

	qualified, tele43 := experimentalV4PhoneSearchBuild43(work, boundary, cw, ch)
	report.Build43Attempted = tele43.Attempted
	report.Build43FrozenCandidates = tele43.FrozenCandidates
	report.Build43QualifiedCandidates = tele43.QualifiedCandidates
	report.Build43PairRanking = append([]ExperimentalV4PhonePairScore(nil), tele43.PairRanking...)
	report.DirectEnsembleAvailable = len(qualified) >= experimentalV4PhoneBuild46DirectEnsembleRequired
	report.Build42BankAvailable = len(qualified) >= experimentalV4PhoneBuild46Build42BankRequired

	plane := newPixelPlane(work)
	srcBounds := src.Bounds()
	for i, h := range qualified {
		c := ExperimentalV4PhoneBuild46Candidate{
			Index:           i,
			SourcePair:      h.build43Pair,
			SourcePairRank:  h.build43PairRank,
			ProposalScore:   h.proposal,
			ValidationScore: h.validation,
			PilotScore:      h.detection.Score,
			PilotMargin:     h.detection.Margin,
			OriginXBlocks:   h.detection.OriginXBlocks,
			OriginYBlocks:   h.detection.OriginYBlocks,
			WorkingQuad:     h.quad,
		}
		c.SourceQuad = experimentalV4PhoneBuild46SourceQuad(h.quad, srcBounds, work)
		c.SingleProfilesTried, c.SingleListFramesTried, c.SingleMaxDataConfidence, c.SingleHMACAuthenticated, c.SingleProfile = experimentalV4PhoneBuild46SingleCandidateDecode(plane, h, key, cw, ch)
		if c.SingleHMACAuthenticated {
			report.SingleCandidateAuthenticated = true
		}
		if oracle != nil {
			c.OracleCompared = true
			c.OracleMeanCornerErrorPx, c.OracleMaxCornerErrorPx, c.OracleMeanCornerErrorRatio = experimentalV4PhoneBuild46OracleError(c.SourceQuad, *oracle)
		}
		report.Candidates = append(report.Candidates, c)
	}

	report.Classification = ExperimentalV4PhoneBuild46ClassifyCounts(len(direct41), tele43.FrozenCandidates, len(qualified), report.SingleCandidateAuthenticated)
	return report, nil
}
