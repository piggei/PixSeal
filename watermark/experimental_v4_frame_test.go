package watermark

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"testing"
)

func TestExperimentalV4FrameRoundTripProfiles(t *testing.T) {
	key := []byte("Piccotti")
	source := testImage(592, 512)
	cases := []struct {
		profile Profile
		payload []byte
	}{
		{ProfileRobust, []byte("v4-robust-test")},
		{ProfileBalanced, bytes.Repeat([]byte{'B'}, 29)},
		{ProfileCapacity, bytes.Repeat([]byte{'C'}, 61)},
	}

	for _, tc := range cases {
		t.Run(string(tc.profile), func(t *testing.T) {
			marked, embedInfo, err := ExperimentalV4EmbedWithInfo(source, tc.payload, key, Options{Strength: 24, Profile: tc.profile})
			if err != nil {
				t.Fatal(err)
			}
			if embedInfo.Version != 4 || embedInfo.Profile != tc.profile {
				t.Fatalf("unexpected embed info: %+v", embedInfo)
			}
			if embedInfo.PilotHash != experimentalV4PilotCandidateLockHash {
				t.Fatalf("pilot hash=%s, want locked %s", embedInfo.PilotHash, experimentalV4PilotCandidateLockHash)
			}
			got, info, err := ExperimentalV4ExtractAligned(marked, key)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, tc.payload) {
				t.Fatalf("payload=%q, want %q", got, tc.payload)
			}
			if info.Profile != tc.profile || info.OriginXBlocks != 0 || info.OriginYBlocks != 0 {
				t.Fatalf("unexpected extract info: %+v", info)
			}
			if info.PilotMargin <= 0 {
				t.Fatalf("non-positive pilot margin: %+v", info)
			}
		})
	}
}

func TestExperimentalV4AlignedCropRoundTrip(t *testing.T) {
	key := []byte("Piccotti")
	payload := []byte("aligned crop")
	source := testImage(888, 768)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, Options{Strength: 24, Profile: ProfileRobust})
	if err != nil {
		t.Fatal(err)
	}
	cropped := cropCopy(marked, image.Rect(5*8, 3*8, 5*8+592, 3*8+512))
	got, info, err := ExperimentalV4ExtractAligned(cropped, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload=%q, want %q", got, payload)
	}
	if info.OriginXBlocks != 5 || info.OriginYBlocks != 3 {
		t.Fatalf("origin=(%d,%d), want (5,3)", info.OriginXBlocks, info.OriginYBlocks)
	}
}

func TestExperimentalV4WrongKeyAndCrossVersionIsolation(t *testing.T) {
	key := []byte("Piccotti")
	payload := []byte("version isolated")
	spec, _ := profileSpecFor(ProfileRobust)
	v4Frame := makeExperimentalV4Frame(payload, key, spec)
	if _, err := parseV3Frame(v4Frame, key, spec); err == nil {
		t.Fatal("v3 parser accepted an experimental v4 frame")
	}
	v3Frame := makeV3Frame(payload, key, spec)
	if _, err := parseExperimentalV4Frame(v3Frame, key, spec); err == nil {
		t.Fatal("experimental v4 parser accepted a v3 frame")
	}

	source := testImage(592, 512)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, Options{Strength: 24, Profile: ProfileRobust})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := ExperimentalV4ExtractAligned(marked, []byte("WrongKey")); err == nil {
		t.Fatal("experimental v4 aligned decoder authenticated with the wrong key")
	}
	if _, _, err := ExtractWithInfo(marked, key); err == nil {
		t.Fatal("v3 production decoder authenticated an experimental v4 carrier")
	}
}

func TestExperimentalV4HammingBaselineCorrectsSingleBitPerCodeword(t *testing.T) {
	key := []byte("Piccotti")
	payload := []byte("ecc baseline")
	spec, _ := profileSpecFor(ProfileRobust)
	frame := makeExperimentalV4Frame(payload, key, spec)
	protected := hammingEncode(whiten(bytesToBits(frame), key, experimentalV4WhitenLabel))
	for word := 0; word < 12; word++ {
		protected[word*7+(word%7)] ^= 1
	}
	decoded := hammingDecode(protected)
	raw := bitsToBytes(whiten(decoded, key, experimentalV4WhitenLabel))
	got, err := parseExperimentalV4Frame(raw, key, spec)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload=%q, want %q", got, payload)
	}
}

func TestExperimentalV4FrameDeterministicVector(t *testing.T) {
	key := []byte("Piccotti")
	spec, _ := profileSpecFor(ProfileRobust)
	frame := makeExperimentalV4Frame([]byte("Build31-vector"), key, spec)
	protected := hammingEncode(whiten(bytesToBits(frame), key, experimentalV4WhitenLabel))
	sum := sha256.Sum256(append(append([]byte(nil), frame...), protected...))
	got := fmt.Sprintf("%x", sum[:])
	// Build31 freezes this vector only as an experimental compatibility guard.
	// It is not yet a normative Format-v4 interoperability identifier.
	const want = "52f15752f1829e39ecae215932b56ab4108b9790eb9fb093705e357b68a7d0ea"
	if got != want {
		t.Fatalf("experimental v4 Build31 vector hash=%s, want %s", got, want)
	}
}

func TestExperimentalV4CapacityGeometry(t *testing.T) {
	if got := ExperimentalV4Capacity(testImage(295, 256)); got != 0 {
		t.Fatalf("capacity below v4 minimum=%d, want 0", got)
	}
	if got := ExperimentalV4Capacity(testImage(296, 256)); got != 64 {
		t.Fatalf("capacity at v4 minimum=%d, want 64", got)
	}
}

func TestExperimentalV4LockedDataPartition(t *testing.T) {
	data := experimentalV4LockedDataPositions()
	if len(data) != experimentalV4DataCount {
		t.Fatalf("locked data positions=%d, want %d", len(data), experimentalV4DataCount)
	}
	seen := make(map[int]string, experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks)
	for _, p := range experimentalV4Prototype2PilotPositions {
		if prior := seen[p]; prior != "" {
			t.Fatalf("duplicate tile position %d (%s then pilot)", p, prior)
		}
		seen[p] = "pilot"
	}
	for _, p := range data {
		if prior := seen[p]; prior != "" {
			t.Fatalf("tile position %d overlaps %s and data", p, prior)
		}
		seen[p] = "data"
	}
	if len(seen) != experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks {
		t.Fatalf("partition covers %d positions, want %d", len(seen), experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks)
	}
}

func TestExperimentalV4AlignedJPEGChannel(t *testing.T) {
	key := []byte("Piccotti")
	payload := []byte("jpeg-v4")
	source := testImage(592, 512)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, Options{Strength: 24, Profile: ProfileRobust})
	if err != nil {
		t.Fatal(err)
	}
	jpegImage := experimentalV4JPEGForTest(t, marked, 82)
	got, info, err := ExperimentalV4ExtractAligned(jpegImage, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload=%q, want %q", got, payload)
	}
	t.Logf("jpeg-q82 pilot score=%.6f margin=%.6f data-confidence=%.6f", info.PilotScore, info.PilotMargin, info.Confidence)
}

func TestExperimentalV4MinimumTileRoundTrip(t *testing.T) {
	key := []byte("Piccotti")
	payload := bytes.Repeat([]byte{'M'}, 61)
	source := testImage(296, 256)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, Options{Strength: 24, Profile: ProfileCapacity})
	if err != nil {
		t.Fatal(err)
	}
	got, info, err := ExperimentalV4ExtractAligned(marked, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload mismatch at minimum v4 tile")
	}
	t.Logf("minimum tile pilot score=%.6f margin=%.6f data-confidence=%.6f", info.PilotScore, info.PilotMargin, info.Confidence)
}

func TestExperimentalV4DataMappingProfileCoverage(t *testing.T) {
	for _, spec := range v3Profiles {
		counts := make([]int, spec.codedBits)
		for ordinal := 0; ordinal < experimentalV4DataCount; ordinal++ {
			counts[experimentalV4CodeIndex(ordinal, spec.codedBits)]++
		}
		minCount, maxCount := experimentalV4DataCount, 0
		total := 0
		for i, count := range counts {
			if count == 0 {
				t.Fatalf("profile %s coded bit %d has no data position", spec.profile, i)
			}
			if count < minCount {
				minCount = count
			}
			if count > maxCount {
				maxCount = count
			}
			total += count
		}
		if total != experimentalV4DataCount {
			t.Fatalf("profile %s maps %d positions, want %d", spec.profile, total, experimentalV4DataCount)
		}
		switch spec.profile {
		case ProfileRobust:
			if minCount != 2 || maxCount != 3 {
				t.Fatalf("robust repetition range=%d..%d, want 2..3", minCount, maxCount)
			}
		case ProfileBalanced:
			if minCount != 1 || maxCount != 2 {
				t.Fatalf("balanced repetition range=%d..%d, want 1..2", minCount, maxCount)
			}
		case ProfileCapacity:
			if minCount != 1 || maxCount != 1 {
				t.Fatalf("capacity repetition range=%d..%d, want 1..1", minCount, maxCount)
			}
		}
	}
}

func TestExperimentalV4HeaderBytes(t *testing.T) {
	want := map[Profile]byte{
		ProfileRobust:   0x41,
		ProfileBalanced: 0x42,
		ProfileCapacity: 0x43,
	}
	for profile, expected := range want {
		spec, _ := profileSpecFor(profile)
		if got := experimentalV4HeaderByte(spec); got != expected {
			t.Fatalf("profile %s header=0x%02x, want 0x%02x", profile, got, expected)
		}
	}
}
