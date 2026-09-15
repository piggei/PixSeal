package watermark

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestExperimentalV4BlindGeometryCorpus runs a compact Build26 blind-geometry
// matrix on every local original. The source archive intentionally omits these
// images; PIXSEAL_V4_CORPUS_DIR enables the private/local qualification.
func TestExperimentalV4BlindGeometryCorpus(t *testing.T) {
	paths := experimentalV4ActiveCorpusPaths(t)

	// Keep the local-corpus gate deliberately compact. The full family matrix is
	// exercised by the synthetic Build26 target; the private originals cover two
	// representative compound cases here so all-test does not multiply an already
	// expensive blind search unnecessarily.
	cases := []experimentalV4BlindGeometryCase{
		{name: "rotate-scale", p: experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: 0.93}},
		{name: "combined-perspective", p: experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: 0.92, topInset: 0.030, bottomInset: 0.015}},
	}
	candidate := experimentalV4Prototype2Candidate()
	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			source, _, err := image.Decode(file)
			file.Close()
			if err != nil {
				t.Fatal(err)
			}
			width := source.Bounds().Dx() / blockSize * blockSize
			height := source.Bounds().Dy() / blockSize * blockSize
			if width < 2*experimentalV4TileWidthBlocks*blockSize || height < 3*experimentalV4TileHeightBlocks*blockSize {
				t.Skipf("%dx%d does not provide enough disjoint proposal/validation tiles for Build26", width, height)
			}
			normalized := cropCopy(source, image.Rect(source.Bounds().Min.X, source.Bounds().Min.Y, source.Bounds().Min.X+width, source.Bounds().Min.Y+height))
			carrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
			if err != nil {
				t.Fatal(err)
			}
			for _, tc := range cases {
				truth, outW, outH := experimentalV4BlindHomography(width, height, tc.p)
				positiveImage := experimentalV4WarpHomographyForTest(carrier, truth, outW, outH)
				negativeImage := experimentalV4WarpHomographyForTest(normalized, truth, outW, outH)
				positive := experimentalV4BlindGeometrySearch(positiveImage, candidate, width, height)
				negative := experimentalV4BlindGeometrySearch(negativeImage, candidate, width, height)
				cornerError := experimentalV4HomographyCornerErrorRatio(truth, positive.h, width, height)
				t.Logf("%s positive obj=%.6f score=%.6f margin=%.6f origin=(%d,%d) params=%+v corner-error=%.5f; negative obj=%.6f score=%.6f margin=%.6f",
					tc.name, positive.objective, positive.detection.Score, positive.detection.Margin, positive.detection.OriginXBlocks, positive.detection.OriginYBlocks, positive.params, cornerError,
					negative.objective, negative.detection.Score, negative.detection.Margin)
				if !positive.available || !positive.detection.Available {
					t.Fatalf("%s: blind geometry unavailable", tc.name)
				}
				if positive.detection.OriginXBlocks != 0 || positive.detection.OriginYBlocks != 0 {
					t.Fatalf("%s: origin=(%d,%d), want (0,0)", tc.name, positive.detection.OriginXBlocks, positive.detection.OriginYBlocks)
				}
				if positive.detection.Margin < 0.10 {
					t.Fatalf("%s: positive margin %.6f below Build26 corpus floor 0.10", tc.name, positive.detection.Margin)
				}
				if positive.detection.Margin-negative.detection.Margin < 0.05 {
					t.Fatalf("%s: positive/negative margin separation %.6f below Build26 corpus floor 0.05", tc.name, positive.detection.Margin-negative.detection.Margin)
				}
				if cornerError > 0.04 {
					t.Fatalf("%s: corner error %.5f above Build26 corpus floor 0.04", tc.name, cornerError)
				}
			}
		})
	}
}
