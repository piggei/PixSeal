package watermark

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestExperimentalV4PilotCorpus is intentionally opt-in because release/source
// archives do not ship the local image corpus. It embeds only a synthetic v4
// pilot/data plane; no payload encoder or authentication path is involved.
func TestExperimentalV4PilotCorpus(t *testing.T) {
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

			negative := experimentalV4DetectPilotGrid(normalized, candidate, 8, 0, 0)
			t.Logf("negative-control score=%.6f runner=%.6f margin=%.6f origin=(%d,%d)", negative.Score, negative.RunnerUpScore, negative.Margin, negative.OriginXBlocks, negative.OriginYBlocks)

			carrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
			if err != nil {
				t.Fatal(err)
			}
			cases := []struct {
				name        string
				img         image.Image
				blockSizePx int
				originX     int
				originY     int
			}{
				{name: "native", img: carrier, blockSizePx: 8},
				{name: "jpeg-82", img: experimentalV4JPEGForTest(t, carrier, 82), blockSizePx: 8},
				{name: "blur-1", img: experimentalV4BoxBlurForTest(carrier), blockSizePx: 8},
				{name: "noise-4", img: experimentalV4NoiseForTest(carrier, 4), blockSizePx: 8},
				{name: "gamma-1.15", img: experimentalV4GammaForTest(carrier, 1.15), blockSizePx: 8},
				{name: "resize-75", img: resizeBilinear(carrier, width*3/4, height*3/4), blockSizePx: 6},
				{name: "resize-50", img: resizeBilinear(carrier, width/2, height/2), blockSizePx: 4},
			}
			if width >= (experimentalV4TileWidthBlocks+12)*blockSize && height >= (experimentalV4TileHeightBlocks+10)*blockSize {
				startX, startY := 5*blockSize, 3*blockSize
				endX := width - 7*blockSize
				endY := height - 7*blockSize
				cases = append(cases, struct {
					name        string
					img         image.Image
					blockSizePx int
					originX     int
					originY     int
				}{name: "aligned-crop", img: cropCopy(carrier, image.Rect(startX, startY, endX, endY)), blockSizePx: 8, originX: 5, originY: 3})
			}

			for _, tc := range cases {
				detection := experimentalV4DetectPilotGrid(tc.img, candidate, tc.blockSizePx, 0, 0)
				t.Logf("%s score=%.6f runner=%.6f margin=%.6f visible=%d samples=%d origin=(%d,%d)", tc.name, detection.Score, detection.RunnerUpScore, detection.Margin, detection.VisiblePilotPositions, detection.PilotSamples, detection.OriginXBlocks, detection.OriginYBlocks)
				if !detection.Available {
					t.Fatalf("%s: pilot unavailable", tc.name)
				}
				if detection.OriginXBlocks != tc.originX || detection.OriginYBlocks != tc.originY {
					t.Fatalf("%s: origin=(%d,%d), want (%d,%d)", tc.name, detection.OriginXBlocks, detection.OriginYBlocks, tc.originX, tc.originY)
				}
				if detection.Margin <= 0 {
					t.Fatalf("%s: non-positive runner-up margin: %+v", tc.name, detection)
				}
			}
		})
	}
}
