package watermark

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
)

// Build33 is a ranking-observability checkpoint for the MQ projective+crop
// blocker.  The transform is intentionally known only to this diagnostic so we
// can ask where its basin survives or disappears inside the blind Build29
// search.  No payload/header/HMAC evidence is used to rank geometry.
func TestExperimentalV4Build33ProjectiveRankingDiagnostic(t *testing.T) {
	paths := experimentalV4ActiveCorpusPaths(t)
	var mqPath string
	for _, p := range paths {
		if filepath.Base(p) == "Immagine MQ.png" {
			mqPath = p
			break
		}
	}
	if mqPath == "" {
		t.Skip("active corpus does not contain Immagine MQ.png")
	}

	f, err := os.Open(mqPath)
	if err != nil {
		t.Fatal(err)
	}
	src, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	normalized, w, h := experimentalV4Build29NormalizeSource(src)
	candidate := experimentalV4Prototype2Candidate()

	randomCarrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	frameCarrier, _, err := ExperimentalV4EmbedWithInfo(
		normalized,
		[]byte("v4-b33-auth"),
		[]byte("PixSeal-v4-TestKey-2026"),
		Options{Profile: ProfileRobust, Strength: 24},
	)
	if err != nil {
		t.Fatal(err)
	}

	truth := experimentalV4BlindGeometryParams{
		angleDeg: 9.3, scaleX: 1.08, scaleY: 0.92,
		topInset: 0.030, bottomInset: 0.015,
	}
	geometry, fw, fh := experimentalV4BlindHomography(w, h, truth)
	truthH, ow, oh := experimentalV4PlacementObservedHomographyForTest(
		geometry, fw, fh, [4]int{101, 69, 83, 47}, [4]int{},
	)

	cases := []struct {
		name    string
		carrier image.Image
	}{
		{name: "random-data-plane", carrier: randomCarrier},
		{name: "authenticated-v4-frame", carrier: frameCarrier},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			observed := experimentalV4WarpHomographyForTest(tc.carrier, truthH, ow, oh)
			build33LogProjectiveRanks(t, observed, candidate, w, h, truth)
		})
	}
}

func build33ProjectiveTruthBasin(p, truth experimentalV4BlindGeometryParams) bool {
	return math.Abs(p.angleDeg-truth.angleDeg) <= 0.8 &&
		math.Abs(p.scaleX-truth.scaleX) <= 0.031 &&
		math.Abs(p.scaleY-truth.scaleY) <= 0.031 &&
		math.Abs(p.topInset-truth.topInset) <= 0.016 &&
		math.Abs(p.bottomInset-truth.bottomInset) <= 0.016
}

func build33LogProjectiveRanks(t *testing.T, img image.Image, candidate experimentalV4PilotCandidate, w, h int, truth experimentalV4BlindGeometryParams) {
	t.Helper()
	plane := newPixelPlane(img)
	ow, oh := img.Bounds().Dx(), img.Bounds().Dy()

	structural, _ := experimentalV4JointProjectiveStructuralBank(plane, w, h, ow, oh)
	sRank, sCandidate := build33FirstTruthRankStructural(structural, truth)
	if sRank == 0 {
		t.Fatal("truth basin disappeared before structural top-5000 retention")
	}
	t.Logf("structural truth-basin rank=%d/%d params=%+v score=%.6f", sRank, len(structural), sCandidate.params, sCandidate.structural)

	pre := experimentalV4JointProjectivePrefilterBank(plane, candidate, w, h, ow, oh, structural)
	pRank, pCandidate := build33FirstTruthRankPrefilter(pre, truth)
	if pRank == 0 {
		t.Fatal("truth basin disappeared in half-pilot prefilter")
	}
	t.Logf("half-pilot truth-basin rank=%d/%d params=%+v score=%.6f full=%.6f origin=(%d,%d)",
		pRank, len(pre), pCandidate.c.params, pCandidate.score, pCandidate.local.fullScore,
		pCandidate.local.originX, pCandidate.local.originY)

	if len(pre) > 2000 {
		pre = pre[:2000]
	}
	full := experimentalV4JointProjectiveFullProposals(plane, candidate, w, h, ow, oh, pre)
	fRank, fCandidate := build33FirstTruthRankFull(full, truth)
	if fRank == 0 {
		t.Fatal("truth basin disappeared before full-proposal ranking")
	}
	t.Logf("full-proposal truth-basin rank=%d/%d params=%+v proposal=%.6f origin=(%d,%d)",
		fRank, len(full), fCandidate.params, fCandidate.proposal, fCandidate.originX, fCandidate.originY)

	// Raw top-N lists are misleading when one texture maximum contributes many
	// near-duplicates. Log distinct basins, which is the object Build33 needs to
	// preserve before refinement.
	distinct := make([]experimentalV4JointProjectiveCandidate, 0, 20)
	for _, q := range full {
		near := false
		for _, d := range distinct {
			if experimentalV4ProjectiveBasinNear(q.params, d.params) {
				near = true
				break
			}
		}
		if near {
			continue
		}
		distinct = append(distinct, q)
		if len(distinct) == 20 {
			break
		}
	}
	for i, q := range distinct {
		t.Logf("distinct-basin %02d proposal=%.6f structural=%.6f params=%+v", i+1, q.proposal, q.structural, q.params)
	}
}

func build33FirstTruthRankStructural(items []experimentalV4JointProjectiveCandidate, truth experimentalV4BlindGeometryParams) (int, experimentalV4JointProjectiveCandidate) {
	for i, q := range items {
		if build33ProjectiveTruthBasin(q.params, truth) {
			return i + 1, q
		}
	}
	return 0, experimentalV4JointProjectiveCandidate{}
}

func build33FirstTruthRankPrefilter(items []experimentalV4ProjectivePreScore, truth experimentalV4BlindGeometryParams) (int, experimentalV4ProjectivePreScore) {
	for i, q := range items {
		if build33ProjectiveTruthBasin(q.c.params, truth) {
			return i + 1, q
		}
	}
	return 0, experimentalV4ProjectivePreScore{}
}

func build33FirstTruthRankFull(items []experimentalV4JointProjectiveCandidate, truth experimentalV4BlindGeometryParams) (int, experimentalV4JointProjectiveCandidate) {
	for i, q := range items {
		if build33ProjectiveTruthBasin(q.params, truth) {
			return i + 1, q
		}
	}
	return 0, experimentalV4JointProjectiveCandidate{}
}
