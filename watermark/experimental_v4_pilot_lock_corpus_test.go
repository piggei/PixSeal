package watermark

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

const (
	experimentalV4Build30MinKnownMappingScore      = 0.85
	experimentalV4Build30MinKnownMappingMargin     = 0.40
	experimentalV4Build30MinKnownMappingSeparation = 0.40
	experimentalV4Build30MaxNegativeMargin         = 0.08
)

// TestExperimentalV4PilotCandidateLockCorpus isolates the pilot channel from
// geometry/placement search. The exact mapping used to synthesize each Build29
// fixture is supplied directly, so a failure here would be evidence against the
// pilot itself rather than against a search heuristic. The private images are
// development corpus only and are never shipped in source/evidence archives.
func TestExperimentalV4PilotCandidateLockCorpus(t *testing.T) {
	paths := experimentalV4ActiveCorpusPaths(t)

	candidate := experimentalV4Prototype2Candidate()
	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			source, _, err := image.Decode(f)
			f.Close()
			if err != nil {
				t.Fatal(err)
			}

			cases := []struct {
				name string
				make func(*testing.T, image.Image) (image.Image, image.Image, int, int, homography)
			}{
				{name: "projective-crop", make: experimentalV4Build29ProjectiveFixture},
				{name: "affine-padded", make: experimentalV4Build29PaddedFixture},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					positive, negative, w, h, knownMapping := tc.make(t, source)
					marked := experimentalV4DetectPilotProjective(positive, candidate, w, h, knownMapping)
					control := experimentalV4DetectPilotProjective(negative, candidate, w, h, knownMapping)
					t.Logf("known-mapping marked score=%.6f margin=%.6f origin=(%d,%d); negative score=%.6f margin=%.6f origin=(%d,%d)",
						marked.Score, marked.Margin, marked.OriginXBlocks, marked.OriginYBlocks,
						control.Score, control.Margin, control.OriginXBlocks, control.OriginYBlocks)

					if !marked.Available || marked.Score < experimentalV4Build30MinKnownMappingScore || marked.Margin < experimentalV4Build30MinKnownMappingMargin {
						t.Fatalf("locked pilot known-mapping evidence too weak: %+v", marked)
					}
					if marked.OriginXBlocks != 0 || marked.OriginYBlocks != 0 {
						t.Fatalf("locked pilot known mapping recovered wrong origin: %+v", marked)
					}
					if control.Margin > experimentalV4Build30MaxNegativeMargin {
						t.Fatalf("unmarked known-mapping margin %.6f exceeds %.6f", control.Margin, experimentalV4Build30MaxNegativeMargin)
					}
					if marked.Margin-control.Margin < experimentalV4Build30MinKnownMappingSeparation {
						t.Fatalf("marked/negative known-mapping margin separation %.6f below %.6f", marked.Margin-control.Margin, experimentalV4Build30MinKnownMappingSeparation)
					}
				})
			}
		})
	}
}
