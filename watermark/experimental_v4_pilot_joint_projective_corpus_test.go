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

func TestExperimentalV4Build29JointCorpus(t *testing.T) {
	directory := os.Getenv("PIXSEAL_V4_CORPUS_DIR")
	if directory == "" {
		t.Skip("set PIXSEAL_V4_CORPUS_DIR to run the local Build29 joint corpus")
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
		t.Fatal("no PNG/JPEG files found in local Build29 corpus")
	}
	candidate := experimentalV4Prototype2Candidate()

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
			_, w, h := experimentalV4Build29NormalizeSource(src)
			largeCarrier := w >= 1000 && h >= 1000

			t.Run("projective-crop", func(t *testing.T) {
				positive, negative, w, h, truthH := experimentalV4Build29ProjectiveFixture(t, src)
				got := experimentalV4JointProjectiveSearch(positive, candidate, w, h)
				accepted := experimentalV4Build29ProjectiveAccepted(got)
				if largeCarrier {
					if !accepted {
						t.Fatalf("Build29 large-carrier projective case was not accepted: %+v", got)
					}
					cornerError := experimentalV4HomographyCornerErrorRatio(truthH, got.placement.h, w, h)
					negativeDetection := experimentalV4DetectPilotProjective(negative, candidate, w, h, got.placement.h)
					t.Logf("ACCEPT params=%+v shift=(%.0f,%.0f) validation=%.6f score=%.6f margin=%.6f origin=(%d,%d) corner-error=%.5f hypotheses=%d; matched-negative score=%.6f margin=%.6f",
						got.params, got.placement.shiftX, got.placement.shiftY, got.placement.validationScore, got.detection.Score, got.detection.Margin,
						got.detection.OriginXBlocks, got.detection.OriginYBlocks, cornerError, got.hypothesesEvaluated, negativeDetection.Score, negativeDetection.Margin)
					if cornerError > 0.020 {
						t.Fatalf("Build29 projective corpus corner error %.5f above 0.020", cornerError)
					}
					if got.detection.Margin-negativeDetection.Margin < 0.10 {
						t.Fatalf("Build29 projective corpus marked/negative margin separation %.6f below 0.10", got.detection.Margin-negativeDetection.Margin)
					}
				} else {
					t.Logf("SAFE-REJECT accepted=%v available=%v params=%+v validation=%.6f score=%.6f margin=%.6f origin=(%d,%d) hypotheses=%d",
						accepted, got.available, got.params, got.placement.validationScore, got.detection.Score, got.detection.Margin,
						got.detection.OriginXBlocks, got.detection.OriginYBlocks, got.hypothesesEvaluated)
					if accepted {
						t.Fatalf("Build29 small-carrier projective case must safe-reject, got %+v", got)
					}
				}
			})

			t.Run("affine-padded-placement", func(t *testing.T) {
				positive, _, w, h, _ := experimentalV4Build29PaddedFixture(t, src)
				got := experimentalV4JointPaddedAffineSearch(positive, candidate, w, h)
				accepted := experimentalV4Build29PaddedAccepted(got)
				t.Logf("SAFE-REJECT accepted=%v available=%v params=%+v validation=%.6f score=%.6f margin=%.6f origin=(%d,%d) hypotheses=%d",
					accepted, got.available, got.params, got.placement.validationScore, got.detection.Score, got.detection.Margin,
					got.detection.OriginXBlocks, got.detection.OriginYBlocks, got.hypothesesEvaluated)
				if accepted {
					t.Fatalf("Build29 photographic padded joint search remains intentionally unqualified and must reject, got %+v", got)
				}
				if largeCarrier {
					truth := experimentalV4BlindGeometryParams{angleDeg: -7.4, scaleX: 1.03, scaleY: 0.97}
					geometry, fw, fh := experimentalV4BlindHomography(w, h, truth)
					placement := experimentalV4PilotPlacementSearch(positive, candidate, w, h, geometry, fw, fh)
					if !placement.available {
						t.Fatal("known-geometry padded placement unexpectedly unavailable")
					}
					detection := experimentalV4DetectPilotProjective(positive, candidate, w, h, placement.h)
					t.Logf("known-geometry control validation=%.6f score=%.6f margin=%.6f origin=(%d,%d)", placement.validationScore, detection.Score, detection.Margin, detection.OriginXBlocks, detection.OriginYBlocks)
					if placement.validationScore < 0.80 || detection.Margin < 0.40 || detection.OriginXBlocks != 0 || detection.OriginYBlocks != 0 {
						t.Fatalf("known-geometry padded control should retain strong signal: placement=%+v detection=%+v", placement, detection)
					}
				}
			})
		})
	}
}
