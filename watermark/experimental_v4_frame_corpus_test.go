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

// TestExperimentalV4FrameCorpus exercises the first real v4 framing/encoder
// path on the private local development originals. The corpus remains opt-in
// and is never distributed with source/evidence archives.
func TestExperimentalV4FrameCorpus(t *testing.T) {
	paths := experimentalV4ActiveCorpusPaths(t)

	key := []byte("PixSeal-v4-TestKey-2026")
	payloads := map[Profile][]byte{
		ProfileRobust:   []byte("Build31-v4-real"),
		ProfileBalanced: bytes.Repeat([]byte{'B'}, 29),
		ProfileCapacity: bytes.Repeat([]byte{'C'}, 61),
	}
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

			for _, profile := range []Profile{ProfileRobust, ProfileBalanced, ProfileCapacity} {
				payload := payloads[profile]
				marked, embedInfo, err := ExperimentalV4EmbedWithInfo(normalized, payload, key, Options{Strength: 24, Profile: profile})
				if err != nil {
					t.Fatal(err)
				}
				got, extractInfo, err := ExperimentalV4ExtractAligned(marked, key)
				if err != nil {
					t.Fatalf("%s native: %v", profile, err)
				}
				if !bytes.Equal(got, payload) {
					t.Fatalf("%s native payload mismatch", profile)
				}
				t.Logf("%s native pilot=%.6f margin=%.6f data-confidence=%.6f coded=%d", profile, extractInfo.PilotScore, extractInfo.PilotMargin, extractInfo.Confidence, embedInfo.ProtectedBits)

				if profile == ProfileRobust {
					jpegImage := experimentalV4JPEGForTest(t, marked, 82)
					jpegPayload, jpegInfo, err := ExperimentalV4ExtractAligned(jpegImage, key)
					if err != nil {
						t.Fatalf("jpeg-q82: %v", err)
					}
					if !bytes.Equal(jpegPayload, payload) {
						t.Fatal("jpeg-q82 payload mismatch")
					}
					t.Logf("jpeg-q82 pilot=%.6f margin=%.6f data-confidence=%.6f", jpegInfo.PilotScore, jpegInfo.PilotMargin, jpegInfo.Confidence)

					if width >= (experimentalV4TileWidthBlocks+12)*blockSize && height >= (experimentalV4TileHeightBlocks+10)*blockSize {
						cropped := cropCopy(marked, image.Rect(5*blockSize, 3*blockSize, width-7*blockSize, height-7*blockSize))
						cropPayload, cropInfo, err := ExperimentalV4ExtractAligned(cropped, key)
						if err != nil {
							t.Fatalf("aligned-crop: %v", err)
						}
						if !bytes.Equal(cropPayload, payload) {
							t.Fatal("aligned-crop payload mismatch")
						}
						if cropInfo.OriginXBlocks != 5 || cropInfo.OriginYBlocks != 3 {
							t.Fatalf("aligned-crop origin=(%d,%d), want (5,3)", cropInfo.OriginXBlocks, cropInfo.OriginYBlocks)
						}
						t.Logf("aligned-crop pilot=%.6f margin=%.6f data-confidence=%.6f", cropInfo.PilotScore, cropInfo.PilotMargin, cropInfo.Confidence)
					}
				}
			}
		})
	}
}
