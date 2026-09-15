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

func TestExperimentalV4PlacementCorpus(t *testing.T) {
	paths := experimentalV4ActiveCorpusPaths(t)
	candidate := experimentalV4Prototype2Candidate()
	cases := []experimentalV4PlacementCase{
		{name: "rotate-scale-crop", p: experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: .93}, crop: [4]int{117, 83, 95, 61}},
		{name: "combined-perspective-pad", p: experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: .92, topInset: .03, bottomInset: .015}, pad: [4]int{67, 43, 79, 51}},
	}
	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			src, _, err := image.Decode(f)
			f.Close()
			if err != nil {
				t.Fatal(err)
			}
			width := src.Bounds().Dx() / blockSize * blockSize
			height := src.Bounds().Dy() / blockSize * blockSize
			if width < 2*experimentalV4TileWidthBlocks*blockSize || height < 3*experimentalV4TileHeightBlocks*blockSize {
				t.Skipf("%dx%d too small for Build27 2x3 placement corpus", width, height)
			}
			normalized := cropCopy(src, image.Rect(src.Bounds().Min.X, src.Bounds().Min.Y, src.Bounds().Min.X+width, src.Bounds().Min.Y+height))
			carrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
			if err != nil {
				t.Fatal(err)
			}
			for _, tc := range cases {
				geometry, fw, fh := experimentalV4BlindHomography(width, height, tc.p)
				observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, tc.crop, tc.pad)
				posImg := experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh)
				negImg := experimentalV4WarpHomographyForTest(normalized, observedH, ow, oh)
				got := experimentalV4PilotPlacementSearch(posImg, candidate, width, height, geometry, fw, fh)
				neg := experimentalV4PilotPlacementSearch(negImg, candidate, width, height, geometry, fw, fh)
				t.Logf("%s marked shift=(%.0f,%.0f) proposal=%.6f validation=%.6f visible=%d/%d hypotheses=%d; negative=%.6f", tc.name, got.shiftX, got.shiftY, got.proposalScore, got.validationScore, got.visibleProposal, got.visibleValidation, got.hypothesesEvaluated, neg.validationScore)
				if !got.available {
					t.Fatalf("%s placement unavailable", tc.name)
				}
				if got.validationScore < 0.40 {
					t.Fatalf("%s marked validation %.6f below Build27 corpus floor 0.40", tc.name, got.validationScore)
				}
				if got.validationScore-neg.validationScore < 0.15 {
					t.Fatalf("%s marked/negative separation %.6f below Build27 corpus floor 0.15", tc.name, got.validationScore-neg.validationScore)
				}
				d := experimentalV4DetectPilotProjective(posImg, candidate, width, height, got.h)
				if !d.Available || d.OriginXBlocks != 0 || d.OriginYBlocks != 0 || d.Margin < 0.15 {
					t.Fatalf("%s recovered mapping failed full pilot validation: %+v", tc.name, d)
				}
				if math.IsNaN(got.validationScore) || math.IsInf(got.validationScore, 0) {
					t.Fatalf("%s non-finite placement score", tc.name)
				}
			}
		})
	}
}
