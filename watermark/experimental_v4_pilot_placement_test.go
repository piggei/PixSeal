package watermark

import (
	"image"
	"testing"
)

const experimentalV4Build27MinMarkedValidation = 0.45
const experimentalV4Build27MinMarkedNegativeSeparation = 0.20

type experimentalV4PlacementCase struct {
	name string
	p    experimentalV4BlindGeometryParams
	crop [4]int // left, top, right, bottom; zero means no crop
	pad  [4]int // left, top, right, bottom; zero means no padded canvas
}

func TestExperimentalV4Prototype2UnknownPlacementQualification(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	cases := []experimentalV4PlacementCase{
		{name: "rotate-scale-arbitrary-crop", p: experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: .93}, crop: [4]int{113, 77, 91, 55}},
		{name: "combined-perspective-arbitrary-crop", p: experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: .92, topInset: .03, bottomInset: .015}, crop: [4]int{101, 69, 83, 47}},
		{name: "rotate-scale-deep-crop", p: experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: .93}, crop: [4]int{171, 119, 149, 101}},
		{name: "rotate-scale-padded-placement", p: experimentalV4BlindGeometryParams{angleDeg: -7.4, scaleX: 1.03, scaleY: .97}, pad: [4]int{73, 41, 61, 37}},
		{name: "perspective-padded-placement", p: experimentalV4BlindGeometryParams{angleDeg: 6.1, scaleX: 1.02, scaleY: .98, topInset: .035, bottomInset: .010}, pad: [4]int{57, 35, 79, 49}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			geometry, fullW, fullH := experimentalV4BlindHomography(source.Bounds().Dx(), source.Bounds().Dy(), tc.p)
			positiveH, obsW, obsH := experimentalV4PlacementObservedHomographyForTest(geometry, fullW, fullH, tc.crop, tc.pad)
			positive := experimentalV4WarpHomographyForTest(carrier, positiveH, obsW, obsH)
			negative := experimentalV4WarpHomographyForTest(source, positiveH, obsW, obsH)
			got := experimentalV4PilotPlacementSearch(positive, candidate, source.Bounds().Dx(), source.Bounds().Dy(), geometry, fullW, fullH)
			neg := experimentalV4PilotPlacementSearch(negative, candidate, source.Bounds().Dx(), source.Bounds().Dy(), geometry, fullW, fullH)
			t.Logf("marked shift=(%.0f,%.0f) proposal=%.6f validation=%.6f runner=%.6f visible=%d/%d hypotheses=%d; negative validation=%.6f",
				got.shiftX, got.shiftY, got.proposalScore, got.validationScore, got.validationRunnerUp, got.visibleProposal, got.visibleValidation, got.hypothesesEvaluated, neg.validationScore)
			if !got.available {
				t.Fatal("placement unavailable")
			}
			if got.validationScore < experimentalV4Build27MinMarkedValidation {
				t.Fatalf("validation %.6f below Build27 floor %.2f", got.validationScore, experimentalV4Build27MinMarkedValidation)
			}
			if got.validationScore-neg.validationScore < experimentalV4Build27MinMarkedNegativeSeparation {
				t.Fatalf("marked/negative validation separation %.6f below Build27 floor %.2f", got.validationScore-neg.validationScore, experimentalV4Build27MinMarkedNegativeSeparation)
			}
			if got.hypothesesEvaluated > 12000 {
				t.Fatalf("placement budget exceeded: %d", got.hypothesesEvaluated)
			}
			// The recovered mapping must make the complete pilot detectable at its
			// canonical origin. This is stronger than matching one absolute crop
			// translation, which is not unique under whole-tile repetition.
			d := experimentalV4DetectPilotProjective(positive, candidate, source.Bounds().Dx(), source.Bounds().Dy(), got.h)
			if !d.Available || d.OriginXBlocks != 0 || d.OriginYBlocks != 0 || d.Margin < 0.20 {
				t.Fatalf("recovered placement does not restore canonical pilot: %+v", d)
			}
		})
	}
}

func experimentalV4PlacementObservedHomographyForTest(base homography, fullW, fullH int, crop, pad [4]int) (homography, int, int) {
	if crop != [4]int{} {
		left, top, right, bottom := crop[0], crop[1], crop[2], crop[3]
		if left+right >= fullW || top+bottom >= fullH {
			panic("invalid Build27 crop")
		}
		shift := homography{h: [9]float64{1, 0, -float64(left), 0, 1, -float64(top), 0, 0, 1}}
		return experimentalV4BlindMultiplyHomography(shift, base), fullW - left - right, fullH - top - bottom
	}
	if pad != [4]int{} {
		left, top, right, bottom := pad[0], pad[1], pad[2], pad[3]
		shift := homography{h: [9]float64{1, 0, float64(left), 0, 1, float64(top), 0, 0, 1}}
		return experimentalV4BlindMultiplyHomography(shift, base), fullW + left + right, fullH + top + bottom
	}
	return base, fullW, fullH
}

func TestExperimentalV4PlacementWrongGeometryDoesNotMasqueradeAsPlacement(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, _ := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	truth := experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: .93}
	geometry, fw, fh := experimentalV4BlindHomography(888, 768, truth)
	observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, [4]int{113, 77, 91, 55}, [4]int{})
	observed := experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh)
	wrong := truth
	wrong.angleDeg += 1.0
	wrongGeometry, wfw, wfh := experimentalV4BlindHomography(888, 768, wrong)
	good := experimentalV4PilotPlacementSearch(observed, candidate, 888, 768, geometry, fw, fh)
	bad := experimentalV4PilotPlacementSearch(observed, candidate, 888, 768, wrongGeometry, wfw, wfh)
	t.Logf("good validation=%.6f bad=%.6f", good.validationScore, bad.validationScore)
	if !good.available {
		t.Fatal("good placement unavailable")
	}
	if bad.available && good.validationScore-bad.validationScore < 0.20 {
		t.Fatalf("wrong geometry too competitive: good %.6f bad %.6f", good.validationScore, bad.validationScore)
	}
}

// keep image import used while making the test file easy to extend with a
// future post-transform photometric case.
var _ image.Image
