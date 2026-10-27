package appservice

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"

	jpeglegacy "github.com/pj/pixseal/internal/jpeglegacy"
)

func fixture(w, h int) *image.NRGBA {
	p := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p.SetNRGBA(x, y, color.NRGBA{R: byte(x*3 + y), G: byte(x + y*5), B: byte(x*7 + y*2), A: 255})
		}
	}
	return p
}
func TestCreateReadAuthenticateAndReportPrivacy(t *testing.T) {
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, fixture(641, 513)); err != nil {
		t.Fatal(err)
	}
	doc, err := Load(bytes.NewReader(encoded.Bytes()), "foto.png")
	if err != nil {
		t.Fatal(err)
	}
	const message = "Prova Android 105"
	const key = "chiave-demo-105"
	result, err := Embed(doc, message, key)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Report.Authenticated || result.Report.CanonicalWidth != 640 || result.Report.CanonicalHeight != 512 {
		t.Fatalf("wrong prepared dimensions/result: %+v", result.Report)
	}
	marked, err := Load(bytes.NewReader(result.PNG), "pixseal.png")
	if err != nil {
		t.Fatal(err)
	}
	read, err := Extract(marked, key, "original", 0, 0)
	if err != nil || read.Message != message || !read.Report.Authenticated {
		t.Fatalf("authenticated read: %v %+v", err, read.Report)
	}
	wrong, err := Extract(marked, "altra-chiave-105", "original", 0, 0)
	if err == nil || wrong.Message != "" || wrong.Report.Authenticated {
		t.Fatal("unauthenticated data exposed")
	}
	report, _ := json.Marshal(result.Report)
	if bytes.Contains(report, []byte(key)) || bytes.Contains(report, []byte(message)) {
		t.Fatal("report contains secret key or payload")
	}
}
func TestLoadPinsJPEGDecoderAndChecksHeaderSize(t *testing.T) {
	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, fixture(32, 24), &jpeg.Options{Quality: 85}); err != nil {
		t.Fatal(err)
	}
	doc, err := Load(bytes.NewReader(buffer.Bytes()), "camera.jpg")
	if err != nil {
		t.Fatal(err)
	}
	reference, err := jpeglegacy.Decode(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < 24; y++ {
		for x := 0; x < 32; x++ {
			ar, ag, ab, aa := doc.Image.At(x, y).RGBA()
			br, bg, bb, ba := reference.At(x, y).RGBA()
			if ar != br || ag != bg || ab != bb || aa != ba {
				t.Fatal("JPEG raster differs from CLI codec")
			}
		}
	}
	buffer.Reset()
	if err := png.Encode(&buffer, fixture(1, 1)); err != nil {
		t.Fatal(err)
	}
	bomb := append([]byte(nil), buffer.Bytes()...)
	binary.BigEndian.PutUint32(bomb[16:20], uint32(MaxPhotoPixels+1))
	binary.BigEndian.PutUint32(bomb[29:33], crc32.ChecksumIEEE(bomb[12:29]))
	if _, err := Load(bytes.NewReader(bomb), "oversized.png"); err == nil || !strings.Contains(err.Error(), "24 megapixel") {
		t.Fatalf("oversized image not rejected at config: %v", err)
	}
	if _, err := Load(strings.NewReader("not an image"), "image.jpg"); err == nil {
		t.Fatal("accepted arbitrary bytes")
	}
}
func TestValidationAndPreview(t *testing.T) {
	doc := Document{Image: fixture(592, 512)}
	for _, tc := range []struct {
		mode string
		w, h int
	}{{"photo", 0, 0}, {"photo", 1631, 1632}, {"photo", 4096, 4096}, {"invalid", 1632, 1632}} {
		if _, err := Extract(doc, "12345678", tc.mode, tc.w, tc.h); err == nil {
			t.Fatal("invalid photo request accepted")
		}
	}
	if _, err := Embed(doc, strings.Repeat("è", 33), "12345678"); err == nil {
		t.Fatal("message limit used runes instead of bytes")
	}
	if _, err := Embed(Document{}, "hello", "12345678"); err == nil {
		t.Fatal("nil image accepted")
	}
	if _, err := Embed(doc, "hello", "short"); err == nil {
		t.Fatal("short key accepted")
	}
	p := Preview(doc.Image)
	if p.Bounds().Dx() > 320 || p.Bounds().Dy() > 320 {
		t.Fatal("oversized preview")
	}
}
func TestPrepareDownscaleAndAlpha(t *testing.T) {
	img := image.NewNRGBA(image.Rect(17, 23, 2017, 1023))
	p := resize(img, 1632, 816)
	if p.Bounds() != image.Rect(0, 0, 1632, 816) || p.NRGBAAt(0, 0) != (color.NRGBA{255, 255, 255, 255}) {
		t.Fatal("non-zero bounds/alpha-on-white preparation changed")
	}
}
