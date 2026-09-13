package watermark

import (
	"image"
	"math"
	"testing"
)

const (
	experimentalV4Build26MinPositiveMargin   = 0.15
	experimentalV4Build26MinMarginSeparation = 0.08
	experimentalV4Build26MaxCornerErrorRatio = 0.035
)

type experimentalV4BlindGeometryCase struct {
	name string
	p    experimentalV4BlindGeometryParams
	post func(*testing.T, image.Image) image.Image
}

func TestExperimentalV4Prototype2BlindGeometryQualification(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}

	cases := []experimentalV4BlindGeometryCase{
		{name: "rotate-7.4", p: experimentalV4BlindGeometryParams{angleDeg: 7.4, scaleX: 1, scaleY: 1}},
		{name: "rotate--13.2", p: experimentalV4BlindGeometryParams{angleDeg: -13.2, scaleX: 1, scaleY: 1}},
		{name: "rotate-scale-holdout", p: experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: 0.93}},
		{name: "shearx-6.5", p: experimentalV4BlindGeometryParams{scaleX: 1, scaleY: 1, shearXDeg: 6.5}},
		{name: "sheary--5.5", p: experimentalV4BlindGeometryParams{scaleX: 1, scaleY: 1, shearYDeg: -5.5}},
		{name: "perspective-holdout", p: experimentalV4BlindGeometryParams{scaleX: 1, scaleY: 1, topInset: 0.035, bottomInset: 0.010}},
		{name: "combined-perspective", p: experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: 0.92, topInset: 0.030, bottomInset: 0.015}},
		{name: "combined-shear", p: experimentalV4BlindGeometryParams{angleDeg: -8.7, scaleX: 1, scaleY: 1, shearXDeg: 6.0}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			truth, width, height := experimentalV4BlindHomography(source.Bounds().Dx(), source.Bounds().Dy(), tc.p)
			positiveImage := experimentalV4WarpHomographyForTest(carrier, truth, width, height)
			negativeImage := experimentalV4WarpHomographyForTest(source, truth, width, height)
			if tc.post != nil {
				positiveImage = toNRGBA(tc.post(t, positiveImage))
				negativeImage = toNRGBA(tc.post(t, negativeImage))
			}

			positive := experimentalV4BlindGeometrySearch(positiveImage, candidate, source.Bounds().Dx(), source.Bounds().Dy())
			negative := experimentalV4BlindGeometrySearch(negativeImage, candidate, source.Bounds().Dx(), source.Bounds().Dy())
			cornerError := experimentalV4HomographyCornerErrorRatio(truth, positive.h, source.Bounds().Dx(), source.Bounds().Dy())
			t.Logf("positive obj=%.6f runner=%.6f gap=%.6f score=%.6f margin=%.6f origin=(%d,%d) hypotheses=%d params=%+v corner-error=%.5f; negative obj=%.6f score=%.6f margin=%.6f",
				positive.objective, positive.runnerUpObjective, positive.objective-positive.runnerUpObjective,
				positive.detection.Score, positive.detection.Margin, positive.detection.OriginXBlocks, positive.detection.OriginYBlocks,
				positive.hypothesesEvaluated, positive.params, cornerError,
				negative.objective, negative.detection.Score, negative.detection.Margin)

			if !positive.available || !positive.detection.Available {
				t.Fatal("blind pilot geometry unavailable")
			}
			if positive.objective-positive.runnerUpObjective <= 1e-6 {
				t.Fatalf("blind geometry top-2 remains ambiguous: best=%.6f runner=%.6f", positive.objective, positive.runnerUpObjective)
			}
			if positive.detection.Margin < experimentalV4Build26MinPositiveMargin {
				t.Fatalf("pilot origin margin %.6f below Build26 floor %.2f", positive.detection.Margin, experimentalV4Build26MinPositiveMargin)
			}
			if positive.detection.Margin-negative.detection.Margin < experimentalV4Build26MinMarginSeparation {
				t.Fatalf("positive/negative blind margin separation %.6f below Build26 floor %.2f", positive.detection.Margin-negative.detection.Margin, experimentalV4Build26MinMarginSeparation)
			}
			if cornerError > experimentalV4Build26MaxCornerErrorRatio {
				t.Fatalf("blind geometry corner error %.5f above Build26 floor %.3f", cornerError, experimentalV4Build26MaxCornerErrorRatio)
			}
			if positive.hypothesesEvaluated > 6400 || negative.hypothesesEvaluated > 6400 {
				t.Fatalf("Build26 search budget exceeded: positive=%d negative=%d max=6400", positive.hypothesesEvaluated, negative.hypothesesEvaluated)
			}
		})
	}
}

func experimentalV4HomographyCornerErrorRatio(want, got homography, width, height int) float64 {
	corners := [4][2]float64{{0, 0}, {float64(width - 1), 0}, {0, float64(height - 1)}, {float64(width - 1), float64(height - 1)}}
	sum := 0.0
	count := 0
	for _, point := range corners {
		wx, wy, wok := want.mapPoint(point[0], point[1])
		gx, gy, gok := got.mapPoint(point[0], point[1])
		if !wok || !gok {
			return math.Inf(1)
		}
		dx, dy := wx-gx, wy-gy
		sum += dx*dx + dy*dy
		count++
	}
	rms := math.Sqrt(sum / float64(count))
	diagonal := math.Hypot(float64(width), float64(height))
	return rms / diagonal
}

func TestExperimentalV4DataRepeatProposalIsPilotSymbolIndependent(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	params := experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: 0.92, topInset: 0.030, bottomInset: 0.015}
	h, width, height := experimentalV4BlindHomography(source.Bounds().Dx(), source.Bounds().Dy(), params)
	observed := experimentalV4WarpHomographyForTest(carrier, h, width, height)
	plane := newPixelPlane(observed)

	score, pairs := experimentalV4DataRepeatConsistencyPlane(plane, candidate, source.Bounds().Dx(), source.Bounds().Dy(), h)
	mutated := candidate
	for i := range mutated.signs {
		mutated.signs[i] = -mutated.signs[i]
	}
	mutatedScore, mutatedPairs := experimentalV4DataRepeatConsistencyPlane(plane, mutated, source.Bounds().Dx(), source.Bounds().Dy(), h)
	if pairs < 64 || mutatedPairs != pairs {
		t.Fatalf("repeat support=%d mutated=%d", pairs, mutatedPairs)
	}
	if math.Abs(score-mutatedScore) > 1e-12 {
		t.Fatalf("repeat proposal changed with public pilot signs: %.12f vs %.12f", score, mutatedScore)
	}
	if score < 0.70 {
		t.Fatalf("true-geometry repeat score %.6f below Build26 structural floor 0.70", score)
	}

	wrong := params
	wrong.angleDeg += 0.7
	wrongH, _, _ := experimentalV4BlindHomography(source.Bounds().Dx(), source.Bounds().Dy(), wrong)
	wrongScore, _ := experimentalV4DataRepeatConsistencyPlane(plane, candidate, source.Bounds().Dx(), source.Bounds().Dy(), wrongH)
	if score-wrongScore < 0.25 {
		t.Fatalf("repeat geometry separation %.6f too small: correct=%.6f wrong=%.6f", score-wrongScore, score, wrongScore)
	}
}
