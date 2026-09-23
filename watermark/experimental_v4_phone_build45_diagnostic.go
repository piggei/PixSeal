package watermark

import (
	"errors"
	"image"
	"math"
)

// ExperimentalV4PhoneBuild45Classification is diagnostic-only. It describes
// where the unchanged Build44 production path stopped; it never changes
// candidate generation, ranking, qualification, ECC or HMAC behavior.
type ExperimentalV4PhoneBuild45Classification string

const (
	ExperimentalV4PhoneBuild45Recovered     ExperimentalV4PhoneBuild45Classification = "recovered"
	ExperimentalV4PhoneBuild45Geometry      ExperimentalV4PhoneBuild45Classification = "geometry"
	ExperimentalV4PhoneBuild45Qualification ExperimentalV4PhoneBuild45Classification = "qualification"
	ExperimentalV4PhoneBuild45DataChannel   ExperimentalV4PhoneBuild45Classification = "data-channel"
)

// ExperimentalV4PhoneBuild45Classify decomposes a production phone result.
// HMAC is only observed after the production decoder has completed; it does not
// feed geometry generation or ranking.
func ExperimentalV4PhoneBuild45Classify(p ExperimentalV4PhoneInfo) ExperimentalV4PhoneBuild45Classification {
	if p.HMACAuthenticated {
		return ExperimentalV4PhoneBuild45Recovered
	}
	if p.Accepted {
		return ExperimentalV4PhoneBuild45DataChannel
	}
	if p.Build43Attempted && p.Build43FrozenCandidates > 0 && p.Build43QualifiedCandidates < experimentalV4PhoneEnsembleSize {
		return ExperimentalV4PhoneBuild45Qualification
	}
	return ExperimentalV4PhoneBuild45Geometry
}

// ExperimentalV4PhoneBuild45OracleInfo reports a laboratory-only decode under
// an independently supplied acquisition-space artwork quadrilateral. The quad
// must be TL, TR, BL, BR in original source-image pixel coordinates. This path
// is intentionally separate from ExperimentalV4ExtractPhone and can never be
// entered by production decoding.
type ExperimentalV4PhoneBuild45OracleInfo struct {
	WorkingWidth      int     `json:"working_width"`
	WorkingHeight     int     `json:"working_height"`
	Downsampled       bool    `json:"downsampled"`
	GeometryAvailable bool    `json:"geometry_available"`
	ProposalScore     float64 `json:"proposal_score"`
	ValidationScore   float64 `json:"validation_score"`
	PilotAvailable    bool    `json:"pilot_available"`
	PilotQualified    bool    `json:"pilot_qualified"`
	PilotScore        float64 `json:"pilot_score"`
	PilotMargin       float64 `json:"pilot_margin"`
	OriginXBlocks     int     `json:"origin_x_blocks"`
	OriginYBlocks     int     `json:"origin_y_blocks"`
	ProfilesTried     int     `json:"profiles_tried"`
	ListFramesTried   int     `json:"list_frames_tried"`
	MaxDataConfidence float64 `json:"max_data_confidence"`
	HMACAuthenticated bool    `json:"hmac_authenticated"`
	Profile           Profile `json:"profile"`
}

// ExperimentalV4PhoneBuild45OracleDecode tests only channel survivability under
// supplied geometry. It does not refine the quad and does not allow payload,
// ECC or HMAC evidence to modify geometry. It is a research oracle, never a
// production fallback.
func ExperimentalV4PhoneBuild45OracleDecode(src image.Image, key []byte, cw, ch int, sourceQuad [4]ImagePoint) ([]byte, ExperimentalV4ExtractInfo, ExperimentalV4PhoneBuild45OracleInfo, error) {
	candidate := experimentalV4Prototype2Candidate()
	info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate)}
	oracle := ExperimentalV4PhoneBuild45OracleInfo{}
	if src == nil {
		return nil, info, oracle, errors.New("nil image")
	}
	if len(key) < 8 {
		return nil, info, oracle, errors.New("key must contain at least 8 bytes")
	}
	if cw < experimentalV4TileWidthBlocks*blockSize || ch < experimentalV4TileHeightBlocks*blockSize || cw%blockSize != 0 || ch%blockSize != 0 {
		return nil, info, oracle, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}

	srcBounds := src.Bounds()
	work, down := experimentalV4PhoneResize(src, experimentalV4PhoneMaxDimension)
	oracle.WorkingWidth = work.Bounds().Dx()
	oracle.WorkingHeight = work.Bounds().Dy()
	oracle.Downsampled = down

	sx := float64(work.Bounds().Dx()) / float64(srcBounds.Dx())
	sy := float64(work.Bounds().Dy()) / float64(srcBounds.Dy())
	q := sourceQuad
	for i := range q {
		q[i].X = (q[i].X - float64(srcBounds.Min.X)) * sx
		q[i].Y = (q[i].Y - float64(srcBounds.Min.Y)) * sy
	}
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return nil, info, oracle, errors.New("oracle quadrilateral does not define a valid homography")
	}
	oracle.GeometryAvailable = true

	plane := newPixelPlane(work)
	oracle.ProposalScore, _ = experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, 0, true, false)
	oracle.ValidationScore, _ = experimentalV4PhoneBuild41FoldScore(plane, candidate, cw, ch, h, 0, false, false)
	det := experimentalV4DetectPilotProjective(work, candidate, cw, ch, h)
	oracle.PilotAvailable = det.Available
	oracle.PilotScore = det.Score
	oracle.PilotMargin = det.Margin
	oracle.OriginXBlocks = det.OriginXBlocks
	oracle.OriginYBlocks = det.OriginYBlocks
	oracle.PilotQualified = det.Available && det.OriginXBlocks == 0 && det.OriginYBlocks == 0 && det.Score >= experimentalV4PhonePilotScoreFloor && det.Margin >= experimentalV4PhonePilotMarginFloor
	info.PilotScore = det.Score
	info.PilotMargin = det.Margin
	info.OriginXBlocks = det.OriginXBlocks
	info.OriginYBlocks = det.OriginYBlocks

	// Decode each profile from this one externally supplied mapping. The
	// Build42 ML-ordered Hamming list is reused only after geometry is fixed.
	for _, spec := range v3Profiles {
		margins, confidence, ok := experimentalV4ReadProtectedProjectiveMargins(plane, cw, ch, h, spec.codedBits)
		if !ok {
			continue
		}
		oracle.ProfilesTried++
		if confidence > oracle.MaxDataConfidence {
			oracle.MaxDataConfidence = confidence
		}
		payload, tried, authenticated := experimentalV4PhoneBuild42ListDecode(margins, key, spec)
		oracle.ListFramesTried += tried
		if authenticated {
			oracle.HMACAuthenticated = true
			oracle.Profile = spec.profile
			info.Profile = spec.profile
			info.Confidence = confidence
			return payload, info, oracle, nil
		}
	}

	if math.IsInf(oracle.ProposalScore, -1) || math.IsInf(oracle.ValidationScore, -1) {
		return nil, info, oracle, errors.New("oracle geometry does not provide enough visible pilot evidence")
	}
	return nil, info, oracle, errors.New("oracle geometry did not authenticate a Format-v4 frame")
}
