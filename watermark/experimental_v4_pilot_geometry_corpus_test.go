package watermark

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestExperimentalV4PilotGeometryCorpus applies the Build25 known-geometry
// qualification matrix to every local original image. Source releases omit
// this corpus; the test is opt-in through PIXSEAL_V4_CORPUS_DIR.
func TestExperimentalV4PilotGeometryCorpus(t *testing.T) {
	paths := experimentalV4ActiveCorpusPaths(t)

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
			if width < experimentalV4TileWidthBlocks*blockSize || height < experimentalV4TileHeightBlocks*blockSize {
				t.Skipf("%dx%d is below the experimental v4 minimum", width, height)
			}
			normalized := cropCopy(source, image.Rect(source.Bounds().Min.X, source.Bounds().Min.Y, source.Bounds().Min.X+width, source.Bounds().Min.Y+height))
			carrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
			if err != nil {
				t.Fatal(err)
			}

			for _, tc := range experimentalV4Build25GeometryCases(t, width, height) {
				positiveImage := experimentalV4WarpHomographyForTest(carrier, tc.h, tc.w, tc.hgt)
				negativeImage := experimentalV4WarpHomographyForTest(normalized, tc.h, tc.w, tc.hgt)
				if tc.post != nil {
					positiveImage = toNRGBA(tc.post(t, positiveImage))
					negativeImage = toNRGBA(tc.post(t, negativeImage))
				}
				positive := experimentalV4DetectPilotProjective(positiveImage, candidate, width, height, tc.h)
				negative := experimentalV4DetectPilotProjective(negativeImage, candidate, width, height, tc.h)
				t.Logf("%s positive score=%.6f runner=%.6f margin=%.6f visible=%d samples=%d origin=(%d,%d); negative score=%.6f runner=%.6f margin=%.6f origin=(%d,%d)",
					tc.name, positive.Score, positive.RunnerUpScore, positive.Margin, positive.VisiblePilotPositions, positive.PilotSamples, positive.OriginXBlocks, positive.OriginYBlocks,
					negative.Score, negative.RunnerUpScore, negative.Margin, negative.OriginXBlocks, negative.OriginYBlocks)
				if !positive.Available {
					t.Fatalf("%s: pilot unavailable", tc.name)
				}
				if positive.OriginXBlocks != 0 || positive.OriginYBlocks != 0 {
					t.Fatalf("%s: origin=(%d,%d), want (0,0)", tc.name, positive.OriginXBlocks, positive.OriginYBlocks)
				}
				if positive.Margin <= 0 {
					t.Fatalf("%s: non-positive runner-up margin: %+v", tc.name, positive)
				}
				if positive.Margin-negative.Margin < experimentalV4Build25SeparationFloor {
					t.Fatalf("%s: pilot/negative margin separation %.6f below Build25 development floor %.2f (positive %.6f, negative %.6f)", tc.name, positive.Margin-negative.Margin, experimentalV4Build25SeparationFloor, positive.Margin, negative.Margin)
				}
			}
		})
	}
}
