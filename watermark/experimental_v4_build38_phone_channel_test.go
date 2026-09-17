package watermark

import (
	"bytes"
	"image/jpeg"
	"testing"
)

func TestExperimentalV4Build38PhoneStrengthRoundTrip(t *testing.T) {
	const width = 1184
	const height = 1024
	key := []byte("PixSeal-v4-TestKey-2026")
	payload := []byte("v4-b38-phone")
	options := DefaultOptions()
	options.Profile = ProfileRobust
	options.Strength = 48

	source := testImage(width, height)
	marked, info, err := ExperimentalV4EmbedWithInfo(source, payload, key, options)
	if err != nil {
		t.Fatal(err)
	}
	if info.Profile != ProfileRobust {
		t.Fatalf("profile=%s want robust", info.Profile)
	}

	got, _, err := ExperimentalV4ExtractAligned(marked, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload=%q want=%q", got, payload)
	}

	var encoded bytes.Buffer
	if err := jpeg.Encode(&encoded, marked, &jpeg.Options{Quality: 82}); err != nil {
		t.Fatal(err)
	}
	jpegCarrier, err := jpeg.Decode(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	got, _, err = ExperimentalV4ExtractAligned(jpegCarrier, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("jpeg payload=%q want=%q", got, payload)
	}
}
