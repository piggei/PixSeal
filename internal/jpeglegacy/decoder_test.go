package jpeglegacy

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"os"
	"testing"
)

func TestDeterministicLegacyJPEGFixture(t *testing.T) {
	encoded, err := os.ReadFile("testdata/legacy-sample.jpg")
	if err != nil {
		t.Fatal(err)
	}
	fixtureSum := sha256.Sum256(encoded)
	if got := hex.EncodeToString(fixtureSum[:]); got != "53bfd738a872db16c91409e174ea7666217b5018a7b01921bfdf474105627005" {
		t.Fatalf("fixture sha256=%s", got)
	}
	cfg, err := DecodeConfig(bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 37 || cfg.Height != 29 {
		t.Fatalf("config=%dx%d", cfg.Width, cfg.Height)
	}
	decoded, err := Decode(bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	y, ok := decoded.(*image.YCbCr)
	if !ok {
		t.Fatalf("decoded type %T, want *image.YCbCr", decoded)
	}
	if y.SubsampleRatio != image.YCbCrSubsampleRatio420 {
		t.Fatalf("subsample=%v", y.SubsampleRatio)
	}
	checks := []struct {
		name string
		data []byte
		want string
	}{
		{"Y", y.Y, "eb6ec297f45a764be2c668febd833ad55dda2c75313589fd7d3909d188543147"},
		{"Cb", y.Cb, "a6b18db6ec85f1427dc260c09ceaad4e1b2c83e090d8d7e28c6dafc7a7c9dcda"},
		{"Cr", y.Cr, "1ae8023e61d94ef8fb4079c58f848fae4676d1d697820d7e8a77b47afcc7cccf"},
	}
	for _, check := range checks {
		sum := sha256.Sum256(check.data)
		if got := hex.EncodeToString(sum[:]); got != check.want {
			t.Fatalf("%s sha256=%s want=%s", check.name, got, check.want)
		}
	}
}
