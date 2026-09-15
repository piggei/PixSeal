package watermark

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestExperimentalV4Build34AuthenticatedProjectiveCorpus(t *testing.T) {
	paths := experimentalV4ActiveCorpusPaths(t)
	key := []byte("PixSeal-v4-TestKey-2026")
	payloads := [][]byte{[]byte("v4-b34-auth"), []byte("v4-b34-alt")}
	for _, path := range paths {
		if filepath.Base(path) != "Immagine MQ.png" {
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		src, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		norm, w, h := experimentalV4Build29NormalizeSource(src)
		for _, payload := range payloads {
			payload := payload
			t.Run(string(payload), func(t *testing.T) {
				carrier, _, err := ExperimentalV4EmbedWithInfo(norm, payload, key, Options{Profile: ProfileRobust, Strength: 24})
				if err != nil {
					t.Fatal(err)
				}
				truth := experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: .92, topInset: .03, bottomInset: .015}
				geometry, fw, fh := experimentalV4BlindHomography(w, h, truth)
				truthH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, [4]int{101, 69, 83, 47}, [4]int{})
				observed := experimentalV4WarpHomographyForTest(carrier, truthH, ow, oh)
				got, info, search, err := experimentalV4ExtractJointProjective(observed, key, w, h)
				if err != nil {
					t.Fatalf("authenticated Build34 decode failed: %v; search=%+v info=%+v", err, search, info)
				}
				if !bytes.Equal(got, payload) {
					t.Fatalf("payload=%q want %q", got, payload)
				}
				ce := experimentalV4HomographyCornerErrorRatio(truthH, search.placement.h, w, h)
				t.Logf("payload=%q profile=%s confidence=%.3f validation=%.3f score=%.3f margin=%.3f origin=(%d,%d) corner-error=%.5f", got, info.Profile, info.Confidence, search.placement.validationScore, search.detection.Score, search.detection.Margin, search.detection.OriginXBlocks, search.detection.OriginYBlocks, ce)
				if ce > .020 {
					t.Fatalf("corner error %.5f above .020", ce)
				}
			})
		}
	}
}
