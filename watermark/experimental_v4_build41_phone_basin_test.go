package watermark

import (
	"bytes"
	"image"
	"testing"
)

func experimentalV4Build41PageForTest(src image.Image) image.Image {
	q := [4]ImagePoint{
		{X: 520, Y: 250},
		{X: 2840, Y: 520},
		{X: 310, Y: 2380},
		{X: 2920, Y: 2180},
	}
	h, ok := experimentalV4PhoneQuadHomography(src.Bounds().Dx(), src.Bounds().Dy(), q)
	if !ok {
		panic("invalid Build41 phone test homography")
	}
	return experimentalV4WarpHomographyForTest(src, h, 3300, 2700)
}

func TestExperimentalV4Build41PhoneBasinRecovery(t *testing.T) {
	const width = 1184
	const height = 1024
	key := []byte("PixSeal-v4-TestKey-2026")
	payload := []byte("v4-b41-phone")
	opts := DefaultOptions()
	opts.Profile = ProfileRobust
	opts.Strength = 48

	source := testImage(width, height)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, opts)
	if err != nil {
		t.Fatal(err)
	}

	observed := experimentalV4Build41PageForTest(marked)
	got, _, phone, err := ExperimentalV4ExtractPhone(observed, key, width, height)
	if err != nil {
		t.Fatalf("Build41 marked phone recovery failed: %v info=%+v", err, phone)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload=%q want=%q", got, payload)
	}
	if !phone.BoundaryDetected || !phone.ProjectiveBasinFound || !phone.Accepted || phone.EnsembleCandidates < experimentalV4PhoneEnsembleSize {
		t.Fatalf("Build41 geometry checkpoint not accepted: %+v", phone)
	}
	if phone.OriginXBlocks != 0 || phone.OriginYBlocks != 0 || phone.PilotScore < experimentalV4PhonePilotScoreFloor || phone.PilotMargin < experimentalV4PhonePilotMarginFloor {
		t.Fatalf("Build41 pilot qualification=%+v", phone)
	}
	if !phone.DataDecodeAttempted || !phone.HMACAuthenticated {
		t.Fatalf("Build41 authenticated channel not reached: %+v", phone)
	}

	control := experimentalV4Build41PageForTest(source)
	if _, _, info, err := ExperimentalV4ExtractPhone(control, key, width, height); err == nil || info.HMACAuthenticated {
		t.Fatalf("unmarked Build41 phone control unexpectedly authenticated: %+v", info)
	}
}
