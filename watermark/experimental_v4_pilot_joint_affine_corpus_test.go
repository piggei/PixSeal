package watermark

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestExperimentalV4JointAffineCropCorpus(t *testing.T) {
	directory := os.Getenv("PIXSEAL_V4_CORPUS_DIR")
	if directory == "" {
		t.Skip("set PIXSEAL_V4_CORPUS_DIR to run the local experimental v4 joint-affine corpus")
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := strings.ToLower(e.Name())
		if strings.HasSuffix(n, ".png") || strings.HasSuffix(n, ".jpg") || strings.HasSuffix(n, ".jpeg") {
			paths = append(paths, filepath.Join(directory, e.Name()))
		}
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		t.Fatal("no PNG/JPEG files found in experimental v4 joint-affine corpus")
	}
	candidate := experimentalV4Prototype2Candidate()
	truth := experimentalV4BlindGeometryParams{angleDeg: 11.2, scaleX: 1.07, scaleY: 0.93}
	crop := [4]int{117, 83, 95, 61}
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
				t.Skipf("%dx%d too small for Build28 joint-affine corpus", width, height)
			}
			normalized := cropCopy(src, image.Rect(src.Bounds().Min.X, src.Bounds().Min.Y, src.Bounds().Min.X+width, src.Bounds().Min.Y+height))
			carrier, err := experimentalV4RenderSyntheticCarrier(normalized, candidate, 24, experimentalV4SyntheticDataSeed)
			if err != nil {
				t.Fatal(err)
			}
			geometry, fw, fh := experimentalV4BlindHomography(width, height, truth)
			observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, crop, [4]int{})
			positiveImage := experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh)
			negativeImage := experimentalV4WarpHomographyForTest(normalized, observedH, ow, oh)
			got := experimentalV4JointAffineSearch(positiveImage, candidate, width, height)
			if !got.available || !got.detection.Available {
				t.Fatalf("joint affine search unavailable: %+v", got)
			}
			cornerError := experimentalV4HomographyCornerErrorRatio(observedH, got.placement.h, width, height)
			negativeValidation, _ := experimentalV4PilotPartitionScore(newPixelPlane(negativeImage), candidate, width, height, got.placement.h, 1)
			negativeDetection := experimentalV4DetectPilotProjective(negativeImage, candidate, width, height, got.placement.h)
			t.Logf("params=%+v shift=(%.0f,%.0f) structural=%.6f proposal=%.6f validation=%.6f full-score=%.6f margin=%.6f origin=(%d,%d) corner-error=%.5f hypotheses=%d; matched-negative validation=%.6f margin=%.6f",
				got.params, got.placement.shiftX, got.placement.shiftY, got.structuralScore, got.proposalScore, got.placement.validationScore,
				got.detection.Score, got.detection.Margin, got.detection.OriginXBlocks, got.detection.OriginYBlocks,
				cornerError, got.hypothesesEvaluated, negativeValidation, negativeDetection.Margin)
			if got.placement.validationScore < 0.50 {
				t.Fatalf("held-out corpus validation %.6f below Build28 floor 0.50", got.placement.validationScore)
			}
			if got.detection.OriginXBlocks != 0 || got.detection.OriginYBlocks != 0 || got.detection.Margin < 0.15 {
				t.Fatalf("corpus full public pilot failed absolute-origin gate: %+v", got.detection)
			}
			if got.placement.validationScore-negativeValidation < 0.15 {
				t.Fatalf("corpus marked/matched-negative separation %.6f below Build28 floor 0.15", got.placement.validationScore-negativeValidation)
			}
			if cornerError > 0.015 {
				t.Fatalf("corpus joint mapping corner error %.5f above Build28 floor 0.015", cornerError)
			}
			if got.hypothesesEvaluated > experimentalV4Build28MaxHypotheses {
				t.Fatalf("Build28 corpus search budget exceeded: %d > %d", got.hypothesesEvaluated, experimentalV4Build28MaxHypotheses)
			}
		})
	}
}
