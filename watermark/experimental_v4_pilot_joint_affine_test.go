package watermark

import (
	"image"
	"math"
	"testing"
)

const (
	experimentalV4Build28MinValidation         = 0.55
	experimentalV4Build28MinFullMargin         = 0.20
	experimentalV4Build28MinNegativeSeparation = 0.15
	experimentalV4Build28MaxCornerErrorRatio   = 0.012
	experimentalV4Build28MaxHypotheses         = 260000
)

type experimentalV4JointAffineCase struct {
	name string
	p    experimentalV4BlindGeometryParams
	crop [4]int
}

func TestExperimentalV4Prototype2JointAffineCropQualification(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	cases := []experimentalV4JointAffineCase{
		{name: "rotate-scale-crop", p: experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: 0.93}, crop: [4]int{113, 77, 91, 55}},
		{name: "negative-rotate-scale-crop", p: experimentalV4BlindGeometryParams{angleDeg: -7.4, scaleX: 1.03, scaleY: 0.97}, crop: [4]int{91, 64, 73, 49}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			geometry, fw, fh := experimentalV4BlindHomography(source.Bounds().Dx(), source.Bounds().Dy(), tc.p)
			observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, tc.crop, [4]int{})
			positiveImage := experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh)
			negativeImage := experimentalV4WarpHomographyForTest(source, observedH, ow, oh)

			got := experimentalV4JointAffineSearch(positiveImage, candidate, source.Bounds().Dx(), source.Bounds().Dy())
			if !got.available || !got.detection.Available {
				t.Fatalf("joint affine search unavailable: %+v", got)
			}
			cornerError := experimentalV4HomographyCornerErrorRatio(observedH, got.placement.h, source.Bounds().Dx(), source.Bounds().Dy())
			negPlane := newPixelPlane(negativeImage)
			negativeValidation, _ := experimentalV4PilotPartitionScore(negPlane, candidate, source.Bounds().Dx(), source.Bounds().Dy(), got.placement.h, 1)
			negativeDetection := experimentalV4DetectPilotProjective(negativeImage, candidate, source.Bounds().Dx(), source.Bounds().Dy(), got.placement.h)
			t.Logf("params=%+v shift=(%.0f,%.0f) structural=%.6f proposal=%.6f validation=%.6f full-score=%.6f margin=%.6f origin=(%d,%d) corner-error=%.5f hypotheses=%d; matched-negative validation=%.6f margin=%.6f",
				got.params, got.placement.shiftX, got.placement.shiftY, got.structuralScore, got.proposalScore, got.placement.validationScore,
				got.detection.Score, got.detection.Margin, got.detection.OriginXBlocks, got.detection.OriginYBlocks,
				cornerError, got.hypothesesEvaluated, negativeValidation, negativeDetection.Margin)
			if got.placement.validationScore < experimentalV4Build28MinValidation {
				t.Fatalf("held-out validation %.6f below Build28 floor %.2f", got.placement.validationScore, experimentalV4Build28MinValidation)
			}
			if got.detection.OriginXBlocks != 0 || got.detection.OriginYBlocks != 0 || got.detection.Margin < experimentalV4Build28MinFullMargin {
				t.Fatalf("full public pilot failed absolute-origin gate: %+v", got.detection)
			}
			if got.placement.validationScore-negativeValidation < experimentalV4Build28MinNegativeSeparation {
				t.Fatalf("marked/matched-negative held-out separation %.6f below Build28 floor %.2f", got.placement.validationScore-negativeValidation, experimentalV4Build28MinNegativeSeparation)
			}
			if cornerError > experimentalV4Build28MaxCornerErrorRatio {
				t.Fatalf("joint mapping corner error %.5f above Build28 floor %.3f", cornerError, experimentalV4Build28MaxCornerErrorRatio)
			}
			if got.hypothesesEvaluated > experimentalV4Build28MaxHypotheses {
				t.Fatalf("Build28 search budget exceeded: %d > %d", got.hypothesesEvaluated, experimentalV4Build28MaxHypotheses)
			}
			if math.IsNaN(got.placement.validationScore) || math.IsInf(got.placement.validationScore, 0) {
				t.Fatal("non-finite joint validation")
			}
		})
	}
}

func TestExperimentalV4JointAffineStructuralProposalIsPilotSymbolIndependent(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	truth := experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: 0.93}
	geometry, fw, fh := experimentalV4BlindHomography(source.Bounds().Dx(), source.Bounds().Dy(), truth)
	observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, [4]int{113, 77, 91, 55}, [4]int{})
	observed := experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh)
	plane := newPixelPlane(observed)
	score := experimentalV4JointAffinePhaseContrast(plane, source.Bounds().Dx(), source.Bounds().Dy(), truth, 64)
	mutated := candidate
	for i := range mutated.signs {
		mutated.signs[i] = -mutated.signs[i]
	}
	// The structural primitive has no candidate argument at all; mutating pilot
	// signs therefore cannot change it. Keep the mutation explicit as a guard on
	// the intended evidence separation.
	mutatedScore := experimentalV4JointAffinePhaseContrast(plane, source.Bounds().Dx(), source.Bounds().Dy(), truth, 64)
	if mutated.signs[0] == candidate.signs[0] {
		t.Fatal("pilot-sign mutation did not occur")
	}
	if math.Abs(score-mutatedScore) > 1e-12 {
		t.Fatalf("structural phase contrast depends on pilot signs: %.12f vs %.12f", score, mutatedScore)
	}
	if score <= 0 {
		t.Fatalf("true geometry structural phase contrast not positive: %.6f", score)
	}
}

var _ image.Image
