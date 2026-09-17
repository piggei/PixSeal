package watermark

import (
	"image"
	"testing"
)

func experimentalV4Build39PhonePageForTest(src image.Image) image.Image {
	q := [4]ImagePoint{
		{X: 420, Y: 330},
		{X: 2740, Y: 270},
		{X: 320, Y: 2310},
		{X: 2810, Y: 2370},
	}
	h, ok := experimentalV4PhoneQuadHomography(src.Bounds().Dx(), src.Bounds().Dy(), q)
	if !ok {
		panic("invalid Build39 phone test homography")
	}
	return experimentalV4WarpHomographyForTest(src, h, 3200, 2700)
}

func TestExperimentalV4Build39PhoneRegistrationCheckpoint(t *testing.T) {
	const width = 1184
	const height = 1024
	key := []byte("PixSeal-v4-TestKey-2026")
	payload := []byte("v4-b39-phone")
	opts := DefaultOptions()
	opts.Profile = ProfileRobust
	opts.Strength = 48
	source := testImage(width, height)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, opts)
	if err != nil {
		t.Fatal(err)
	}

	observed := experimentalV4Build39PhonePageForTest(marked)
	work, boundary, bank, evaluated, downsampled := experimentalV4PhoneSearch(observed, width, height)
	if work == nil || downsampled {
		t.Fatalf("unexpected working image/downsample state")
	}
	if !boundary.Detected {
		t.Fatalf("marked phone boundary not detected: %+v", boundary)
	}
	if len(bank) == 0 {
		t.Fatal("marked phone search produced no pilot-qualified candidate")
	}
	if bank[0].proposal < experimentalV4PhoneProposalFloor || bank[0].validation < experimentalV4PhoneValidationFloor {
		t.Fatalf("marked checkpoint proposal=%.4f validation=%.4f eval=%d", bank[0].proposal, bank[0].validation, evaluated)
	}

	control := experimentalV4Build39PhonePageForTest(source)
	_, controlBoundary, controlBank, _, _ := experimentalV4PhoneSearch(control, width, height)
	if !controlBoundary.Detected {
		t.Fatalf("control boundary not detected: %+v", controlBoundary)
	}
	if len(controlBank) > 0 && controlBank[0].proposal >= bank[0].proposal && controlBank[0].validation >= bank[0].validation {
		t.Fatalf("unmarked control dominates marked public-pilot evidence: marked=(%.4f,%.4f) control=(%.4f,%.4f)", bank[0].proposal, bank[0].validation, controlBank[0].proposal, controlBank[0].validation)
	}
}

func TestExperimentalV4Build39PhoneDownsampleBound(t *testing.T) {
	src := testImage(5000, 360)
	work, down := experimentalV4PhoneResize(src, experimentalV4PhoneMaxDimension)
	if !down {
		t.Fatal("expected oversized phone image to be downsampled")
	}
	if work.Bounds().Dx() > experimentalV4PhoneMaxDimension || work.Bounds().Dy() > experimentalV4PhoneMaxDimension {
		t.Fatalf("working image=%dx%d", work.Bounds().Dx(), work.Bounds().Dy())
	}
}

func TestExperimentalV4Build39PhoneRejectsInvalidCanonicalExtent(t *testing.T) {
	if _, _, _, err := ExperimentalV4ExtractPhone(testImage(320, 320), []byte("12345678"), 319, 320); err == nil {
		t.Fatal("accepted non-block-aligned canonical extent")
	}
}
