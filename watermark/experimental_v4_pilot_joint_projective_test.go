package watermark

import (
	"image"
	"testing"
)

const (
	experimentalV4Build29MaxProjectiveHypotheses = 500000
	experimentalV4Build29MaxPaddedHypotheses     = 350000
	experimentalV4Build29MaxCornerError          = 0.015
)

func experimentalV4Build29NormalizeSource(src image.Image) (image.Image, int, int) {
	w := src.Bounds().Dx() / blockSize * blockSize
	h := src.Bounds().Dy() / blockSize * blockSize
	n := cropCopy(src, image.Rect(src.Bounds().Min.X, src.Bounds().Min.Y, src.Bounds().Min.X+w, src.Bounds().Min.Y+h))
	return n, w, h
}

func experimentalV4Build29ProjectiveFixture(t *testing.T, src image.Image) (image.Image, image.Image, int, int, homography) {
	t.Helper()
	candidate := experimentalV4Prototype2Candidate()
	normalized, w, h := experimentalV4Build29NormalizeSource(src)
	carrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	truth := experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: 0.92, topInset: 0.030, bottomInset: 0.015}
	geometry, fw, fh := experimentalV4BlindHomography(w, h, truth)
	observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, [4]int{101, 69, 83, 47}, [4]int{})
	return experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh), experimentalV4WarpHomographyForTest(normalized, observedH, ow, oh), w, h, observedH
}

func experimentalV4Build29PaddedFixture(t *testing.T, src image.Image) (image.Image, image.Image, int, int, homography) {
	t.Helper()
	candidate := experimentalV4Prototype2Candidate()
	normalized, w, h := experimentalV4Build29NormalizeSource(src)
	carrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	truth := experimentalV4BlindGeometryParams{angleDeg: -7.4, scaleX: 1.03, scaleY: 0.97}
	geometry, fw, fh := experimentalV4BlindHomography(w, h, truth)
	observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, [4]int{}, [4]int{73, 41, 61, 37})
	return experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh), experimentalV4WarpHomographyForTest(normalized, observedH, ow, oh), w, h, observedH
}

func TestExperimentalV4Build29JointQualification(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)

	t.Run("projective-crop", func(t *testing.T) {
		positive, negative, w, h, truthH := experimentalV4Build29ProjectiveFixture(t, source)
		got := experimentalV4JointProjectiveSearch(positive, candidate, w, h)
		if !experimentalV4Build29ProjectiveAccepted(got) {
			t.Fatalf("Build29 projective search was not accepted: %+v", got)
		}
		cornerError := experimentalV4HomographyCornerErrorRatio(truthH, got.placement.h, w, h)
		negativeDetection := experimentalV4DetectPilotProjective(negative, candidate, w, h, got.placement.h)
		t.Logf("params=%+v shift=(%.0f,%.0f) structural=%.6f proposal=%.6f geometry-validation=%.6f placement-validation=%.6f score=%.6f margin=%.6f origin=(%d,%d) corner-error=%.5f hypotheses=%d; matched-negative score=%.6f margin=%.6f",
			got.params, got.placement.shiftX, got.placement.shiftY, got.structuralScore, got.proposalObjective, got.geometryValidation,
			got.placement.validationScore, got.detection.Score, got.detection.Margin, got.detection.OriginXBlocks, got.detection.OriginYBlocks,
			cornerError, got.hypothesesEvaluated, negativeDetection.Score, negativeDetection.Margin)
		if cornerError > experimentalV4Build29MaxCornerError {
			t.Fatalf("Build29 projective corner error %.5f above %.5f", cornerError, experimentalV4Build29MaxCornerError)
		}
		if got.detection.Margin-negativeDetection.Margin < 0.10 {
			t.Fatalf("Build29 projective marked/negative margin separation %.6f below 0.10", got.detection.Margin-negativeDetection.Margin)
		}
		if got.hypothesesEvaluated > experimentalV4Build29MaxProjectiveHypotheses {
			t.Fatalf("Build29 projective search budget exceeded: %d > %d", got.hypothesesEvaluated, experimentalV4Build29MaxProjectiveHypotheses)
		}
	})

	t.Run("affine-padded-placement", func(t *testing.T) {
		positive, negative, w, h, truthH := experimentalV4Build29PaddedFixture(t, source)
		got := experimentalV4JointPaddedAffineSearch(positive, candidate, w, h)
		if !experimentalV4Build29PaddedAccepted(got) {
			t.Fatalf("Build29 padded search was not accepted: %+v", got)
		}
		cornerError := experimentalV4HomographyCornerErrorRatio(truthH, got.placement.h, w, h)
		negativeDetection := experimentalV4DetectPilotProjective(negative, candidate, w, h, got.placement.h)
		t.Logf("params=%+v shift=(%.0f,%.0f) structural=%.6f proposal=%.6f placement-validation=%.6f score=%.6f margin=%.6f origin=(%d,%d) corner-error=%.5f hypotheses=%d; matched-negative score=%.6f margin=%.6f",
			got.params, got.placement.shiftX, got.placement.shiftY, got.structuralScore, got.proposalScore,
			got.placement.validationScore, got.detection.Score, got.detection.Margin, got.detection.OriginXBlocks, got.detection.OriginYBlocks,
			cornerError, got.hypothesesEvaluated, negativeDetection.Score, negativeDetection.Margin)
		if cornerError > experimentalV4Build29MaxCornerError {
			t.Fatalf("Build29 padded corner error %.5f above %.5f", cornerError, experimentalV4Build29MaxCornerError)
		}
		if got.detection.Margin-negativeDetection.Margin < 0.10 {
			t.Fatalf("Build29 padded marked/negative margin separation %.6f below 0.10", got.detection.Margin-negativeDetection.Margin)
		}
		if got.hypothesesEvaluated > experimentalV4Build29MaxPaddedHypotheses {
			t.Fatalf("Build29 padded search budget exceeded: %d > %d", got.hypothesesEvaluated, experimentalV4Build29MaxPaddedHypotheses)
		}
	})
}
