package watermark

import "math"

// DiagnosticFormatV4DesignStudy is the non-normative layout comparison first added in build22 of
// dedicated absolute-pilot layouts. It does not define Format v4 and is never
// consulted by Embed/Extract. The arithmetic intentionally keeps v3's 8-byte
// header, 8-byte HMAC tag and Hamming(7,4) expansion so the capacity/geometry
// tradeoff is directly comparable to the frozen v3 format.
type DiagnosticFormatV4DesignStudy struct {
	Method                  string                              `json:"method"`
	AssumedFrameOverhead    int                                 `json:"assumed_frame_overhead_bytes"`
	AssumedECC              string                              `json:"assumed_ecc"`
	V3TilePositions         int                                 `json:"v3_tile_positions"`
	V3MinimumWidthPixels    int                                 `json:"v3_minimum_width_px"`
	V3MinimumHeightPixels   int                                 `json:"v3_minimum_height_px"`
	Candidates              []DiagnosticFormatV4DesignCandidate `json:"candidates"`
	RecommendedPrototype    string                              `json:"recommended_prototype"`
	RecommendationRationale string                              `json:"recommendation_rationale"`
	Note                    string                              `json:"note"`
}

type DiagnosticFormatV4DesignCandidate struct {
	Name                          string  `json:"name"`
	TileWidthBlocks               int     `json:"tile_width_blocks"`
	TileHeightBlocks              int     `json:"tile_height_blocks"`
	TilePositions                 int     `json:"tile_positions"`
	PilotPositions                int     `json:"pilot_positions"`
	PilotFraction                 float64 `json:"pilot_fraction"`
	DataPositions                 int     `json:"data_positions"`
	GeometryAreaRatioVsV3         float64 `json:"geometry_area_ratio_vs_v3"`
	MinimumWidthPixels            int     `json:"minimum_width_px"`
	MinimumHeightPixels           int     `json:"minimum_height_px"`
	ComparableMaxFrameBytes       int     `json:"comparable_max_frame_bytes"`
	ComparableMaxPayloadBytes     int     `json:"comparable_max_payload_bytes"`
	PreservesV3MaxPayload         bool    `json:"preserves_v3_max_payload"`
	RobustRedundancyIfUnchanged   float64 `json:"robust_redundancy_if_unchanged"`
	BalancedRedundancyIfUnchanged float64 `json:"balanced_redundancy_if_unchanged"`
	CapacityRedundancyIfUnchanged float64 `json:"capacity_redundancy_if_unchanged"`
	IdealPilotRandomWrongPhaseZ   float64 `json:"ideal_pilot_random_wrong_phase_z"`
	PilotCorrelationSamples       int     `json:"pilot_phase_correlation_samples"`
}

func diagnosticFormatV4DesignStudy() DiagnosticFormatV4DesignStudy {
	study := DiagnosticFormatV4DesignStudy{
		Method:                  "v4-absolute-pilot-layout-trade-study",
		AssumedFrameOverhead:    headerSize + tagSize,
		AssumedECC:              "Hamming(7,4), unchanged only for apples-to-apples sizing",
		V3TilePositions:         eccBits,
		V3MinimumWidthPixels:    tileWidth * blockSize,
		V3MinimumHeightPixels:   tileHeight * blockSize,
		RecommendedPrototype:    "preserve-capacity-37x32-p64",
		RecommendationRationale: "adds a dedicated 64-block public absolute pilot while retaining 1120 data positions, all v3 payload ceilings and nominal profile redundancy; geometry area grows only 5.7%",
		Note:                    "research sizing model only: build23 adds a separate provisional pilot foundation, but v4 framing, ECC and encoder/decoder remain experimental and are not production-enabled",
	}
	for _, candidate := range []struct {
		name          string
		width, height int
		pilot         int
	}{
		{"compact-35x32-p64", 35, 32, 64},
		{"preserve-capacity-37x32-p64", 37, 32, 64},
		{"strong-pilot-38x32-p96", 38, 32, 96},
	} {
		study.Candidates = append(study.Candidates, diagnosticFormatV4Candidate(candidate.name, candidate.width, candidate.height, candidate.pilot))
	}
	return study
}

func diagnosticFormatV4Candidate(name string, width, height, pilot int) DiagnosticFormatV4DesignCandidate {
	positions := width * height
	data := positions - pilot
	frameBytes := 0
	if data > 0 {
		// Hamming(7,4) consumes 14 carrier positions per source byte.
		frameBytes = data / 14
	}
	payload := frameBytes - (headerSize + tagSize)
	if payload < 0 {
		payload = 0
	}
	return DiagnosticFormatV4DesignCandidate{
		Name:                          name,
		TileWidthBlocks:               width,
		TileHeightBlocks:              height,
		TilePositions:                 positions,
		PilotPositions:                pilot,
		PilotFraction:                 float64(pilot) / float64(positions),
		DataPositions:                 data,
		GeometryAreaRatioVsV3:         float64(positions) / float64(eccBits),
		MinimumWidthPixels:            width * blockSize,
		MinimumHeightPixels:           height * blockSize,
		ComparableMaxFrameBytes:       frameBytes,
		ComparableMaxPayloadBytes:     payload,
		PreservesV3MaxPayload:         payload >= maxPayload,
		RobustRedundancyIfUnchanged:   float64(data) / 448.0,
		BalancedRedundancyIfUnchanged: float64(data) / 672.0,
		CapacityRedundancyIfUnchanged: float64(data) / 1120.0,
		IdealPilotRandomWrongPhaseZ:   math.Sqrt(float64(pilot)),
		PilotCorrelationSamples:       positions * pilot,
	}
}
