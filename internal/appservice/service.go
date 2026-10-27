// Package appservice adapts the unchanged v4 engine to native app workflows.
// It has no GUI, network, file-system persistence or third-party dependencies.
package appservice

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"time"

	"github.com/pj/pixseal/internal/buildinfo"
	jpeglegacy "github.com/pj/pixseal/internal/jpeglegacy"
	"github.com/pj/pixseal/watermark"
)

const MaxFileBytes = 64 << 20
const MaxPhotoPixels int64 = 24_000_000
const PrintMaxEdge = 1632

var ErrNotAuthenticated = errors.New("messaggio non autenticato: verifica la chiave e la qualità dell’immagine")

type Document struct {
	Image                image.Image
	Name, Format, SHA256 string
}
type Report struct {
	Version         string `json:"version"`
	Operation       string `json:"operation"`
	Mode            string `json:"mode"`
	SourceSHA256    string `json:"source_sha256"`
	SourceWidth     int    `json:"source_width"`
	SourceHeight    int    `json:"source_height"`
	CanonicalWidth  int    `json:"canonical_width"`
	CanonicalHeight int    `json:"canonical_height"`
	Authenticated   bool   `json:"authenticated"`
	ElapsedMS       int64  `json:"elapsed_ms"`
	Engine          string `json:"engine"`
	Profile         string `json:"profile,omitempty"`
}
type Result struct {
	Report  Report
	Message string      `json:"-"`
	PNG     []byte      `json:"-"`
	Image   image.Image `json:"-"`
}

func Load(r io.Reader, name string) (Document, error) {
	b, err := io.ReadAll(io.LimitReader(r, MaxFileBytes+1))
	if err != nil {
		return Document{}, fmt.Errorf("lettura immagine: %w", err)
	}
	if len(b) > MaxFileBytes {
		return Document{}, errors.New("file troppo grande: massimo 64 MB")
	}
	format := ""
	if len(b) >= 2 && b[0] == 0xff && b[1] == 0xd8 {
		format = "jpeg"
	} else if bytes.HasPrefix(b, []byte("\x89PNG\r\n\x1a\n")) {
		format = "png"
	} else {
		return Document{}, errors.New("scegli un’immagine JPEG o PNG")
	}
	var cfg image.Config
	if format == "jpeg" {
		cfg, err = jpeglegacy.DecodeConfig(bytes.NewReader(b))
	} else {
		cfg, err = png.DecodeConfig(bytes.NewReader(b))
	}
	if err != nil {
		return Document{}, errors.New("immagine danneggiata o non leggibile")
	}
	if cfg.Width < 1 || cfg.Height < 1 || int64(cfg.Width) > MaxPhotoPixels/int64(cfg.Height) {
		return Document{}, errors.New("immagine troppo grande: usa una foto fino a 24 megapixel")
	}
	var img image.Image
	if format == "jpeg" {
		img, err = jpeglegacy.Decode(bytes.NewReader(b))
	} else {
		img, err = png.Decode(bytes.NewReader(b))
	}
	if err != nil {
		return Document{}, errors.New("impossibile decodificare l’immagine")
	}
	if img.Bounds().Dx() != cfg.Width || img.Bounds().Dy() != cfg.Height {
		return Document{}, errors.New("dimensioni immagine non coerenti")
	}
	hash := sha256.Sum256(b)
	return Document{Image: img, Name: name, Format: format, SHA256: hex.EncodeToString(hash[:])}, nil
}
func baseReport(d Document, operation, mode string) Report {
	return Report{Version: buildinfo.String(), Operation: operation, Mode: mode, SourceSHA256: d.SHA256, SourceWidth: d.Image.Bounds().Dx(), SourceHeight: d.Image.Bounds().Dy(), Engine: "Build98; Android experimental"}
}
func validate(d Document, key string) error {
	if d.Image == nil {
		return errors.New("scegli prima un’immagine")
	}
	if len([]byte(key)) < 8 {
		return errors.New("la chiave deve contenere almeno 8 byte")
	}
	return nil
}
func Embed(d Document, message, key string) (Result, error) {
	if err := validate(d, key); err != nil {
		return Result{}, err
	}
	if len([]byte(message)) == 0 || len([]byte(message)) > 64 {
		return Result{}, errors.New("il messaggio deve contenere da 1 a 64 byte UTF-8")
	}
	start := time.Now()
	report := baseReport(d, "embed", "print-v4")
	result := Result{Report: report}
	w, h := d.Image.Bounds().Dx(), d.Image.Bounds().Dy()
	scale := math.Min(1, float64(PrintMaxEdge)/float64(max(w, h)))
	cw, ch := int(float64(w)*scale)/8*8, int(float64(h)*scale)/8*8
	if cw < 296 || ch < 256 {
		return Result{}, errors.New("immagine troppo piccola o troppo stretta: servono almeno 296 × 256 pixel")
	}
	prepared := resize(d.Image, cw, ch)
	marked, info, err := watermark.ExperimentalV4EmbedWithInfo(prepared, []byte(message), []byte(key), watermark.Options{Strength: 48, Profile: watermark.ProfileAuto})
	if err != nil {
		return result, err
	}
	// A saved carrier is offered only after a real authenticated round trip.
	payload, _, err := watermark.ExperimentalV4ExtractAligned(marked, []byte(key))
	if err != nil || !bytes.Equal(payload, []byte(message)) {
		return result, errors.New("verifica del messaggio creato non riuscita")
	}
	var output bytes.Buffer
	if err := png.Encode(&output, marked); err != nil {
		return result, err
	}
	result.PNG = output.Bytes()
	result.Image = marked
	result.Report.CanonicalWidth, result.Report.CanonicalHeight = cw, ch
	result.Report.Profile = string(info.Profile)
	result.Report.Authenticated = true
	result.Report.ElapsedMS = time.Since(start).Milliseconds()
	return result, nil
}
func Extract(d Document, key, mode string, cw, ch int) (Result, error) {
	if err := validate(d, key); err != nil {
		return Result{}, err
	}
	if mode != "photo" && mode != "original" {
		return Result{}, errors.New("modalità di lettura non valida")
	}
	if mode == "photo" && (cw < 296 || ch < 256 || cw%8 != 0 || ch%8 != 0 || cw > 4096 || ch > 4096 || int64(cw)*int64(ch) > 4_000_000) {
		return Result{}, errors.New("inserisci le dimensioni del PNG prima della stampa: multipli di 8, almeno 296 × 256 e fino a 4 megapixel")
	}
	start := time.Now()
	result := Result{Report: baseReport(d, "extract", mode)}
	if mode == "original" {
		cw, ch = d.Image.Bounds().Dx(), d.Image.Bounds().Dy()
	}
	result.Report.CanonicalWidth, result.Report.CanonicalHeight = cw, ch
	var payload []byte
	var info watermark.ExperimentalV4ExtractInfo
	var err error
	if mode == "photo" {
		var phone watermark.ExperimentalV4PhoneInfo
		payload, info, phone, err = watermark.ExperimentalV4ExtractPhone(d.Image, []byte(key), cw, ch)
		if err == nil && !phone.HMACAuthenticated {
			err = ErrNotAuthenticated
		}
	} else {
		payload, info, err = watermark.ExperimentalV4ExtractAligned(d.Image, []byte(key))
	}
	result.Report.ElapsedMS = time.Since(start).Milliseconds()
	if err != nil {
		return result, ErrNotAuthenticated
	}
	result.Report.Authenticated = true
	result.Report.Profile = string(info.Profile)
	result.Message = string(payload)
	return result, nil
}
func Preview(img image.Image) image.Image {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	scale := math.Min(1, 320/float64(max(w, h)))
	return resize(img, max(1, int(float64(w)*scale)), max(1, int(float64(h)*scale)))
}

// Source preparation applies only when creating a new carrier (or its preview).
// Acquired photographs reach the qualified decoder without resizing/re-encoding.
func resize(src image.Image, w, h int) *image.NRGBA {
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	bounds := src.Bounds()
	sw, sh := bounds.Dx(), bounds.Dy()
	sample := func(x, y int) (float64, float64, float64) {
		r, g, b, a := src.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
		white := uint32(65535) - a
		return float64(r+white) / 257, float64(g+white) / 257, float64(b+white) / 257
	}
	for y := 0; y < h; y++ {
		sy := math.Max(0, (float64(y)+.5)*float64(sh)/float64(h)-.5)
		y0 := int(sy)
		y1 := min(y0+1, sh-1)
		fy := sy - float64(y0)
		for x := 0; x < w; x++ {
			sx := math.Max(0, (float64(x)+.5)*float64(sw)/float64(w)-.5)
			x0 := int(sx)
			x1 := min(x0+1, sw-1)
			fx := sx - float64(x0)
			r0, g0, b0 := sample(x0, y0)
			r1, g1, b1 := sample(x1, y0)
			r2, g2, b2 := sample(x0, y1)
			r3, g3, b3 := sample(x1, y1)
			blend := func(a, b, c, d float64) uint8 { return uint8(math.Round((a*(1-fx)+b*fx)*(1-fy) + (c*(1-fx)+d*fx)*fy)) }
			out.SetNRGBA(x, y, color.NRGBA{R: blend(r0, r1, r2, r3), G: blend(g0, g1, g2, g3), B: blend(b0, b1, b2, b3), A: 255})
		}
	}
	return out
}
