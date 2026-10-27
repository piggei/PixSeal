package watermark

import (
	"image"
	"testing"
)

// First deterministic seed found by semantic coverage search, before collecting
// any timing. Pure public noise, no carrier payload or private photograph.
// A true coordinate-local root with four improving coupled pairs exercises
// 128 continuation FoldScores and returns six ordered Gen4 candidates.
func experimentalV4Build104StageFixture(tb testing.TB) (*pixelPlane, experimentalV4PilotCandidate, int, int, [4]ImagePoint, experimentalV4PhoneHypothesis) {
	tb.Helper()
	const w, h = 592, 512
	p := &pixelPlane{bounds: image.Rect(0, 0, w, h), rgb: make([]byte, w*h*3)}
	state := uint64(107)
	for i := range p.rgb {
		state ^= state << 13
		state ^= state >> 7
		state ^= state << 17
		p.rgb[i] = byte(state >> 24)
	}
	pilot := experimentalV4Prototype2Candidate()
	anchor := [4]ImagePoint{{8, 8}, {w - 9, 8}, {8, h - 9}, {w - 9, h - 9}}
	quad := [4]ImagePoint{{7, 8}, {584, 9}, {8, 503}, {583, 504}}
	hom, ok := experimentalV4PhoneQuadHomography(w, h, quad)
	if !ok {
		tb.Fatal("public fixture homography")
	}
	score, _ := experimentalV4PhoneBuild41FoldScore(p, pilot, w, h, hom, 0, true, false)
	return p, pilot, w, h, anchor, experimentalV4PhoneHypothesis{quad: quad, h: hom, proposal: score}
}
