package watermark

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func experimentalV4Build37ScannerPageForTest(src image.Image) image.Image {
	bounds := src.Bounds()
	const scale = 2
	const offsetX = 360
	const offsetY = 320
	page := image.NewNRGBA(image.Rect(0, 0, 3200, 2700))
	for i := 0; i < len(page.Pix); i += 4 {
		page.Pix[i] = 255
		page.Pix[i+1] = 255
		page.Pix[i+2] = 255
		page.Pix[i+3] = 255
	}
	for y := 0; y < bounds.Dy()*scale; y++ {
		for x := 0; x < bounds.Dx()*scale; x++ {
			c := color.NRGBAModel.Convert(src.At(bounds.Min.X+x/scale, bounds.Min.Y+y/scale)).(color.NRGBA)
			page.SetNRGBA(offsetX+x, offsetY+y, c)
		}
	}
	var compressed bytes.Buffer
	if err := jpeg.Encode(&compressed, page, &jpeg.Options{Quality: 82}); err != nil {
		panic(err)
	}
	decoded, err := jpeg.Decode(bytes.NewReader(compressed.Bytes()))
	if err != nil {
		panic(err)
	}
	return decoded
}

func TestExperimentalV4Build37BlindScannerRegistration(t *testing.T) {
	const width = 1184  // four complete 37-block tiles
	const height = 1024 // four complete 32-block tiles
	key := []byte("PixSeal-v4-TestKey-2026")
	payload := []byte("v4-b37-synth")
	options := DefaultOptions()
	options.Profile = ProfileRobust
	options.Strength = 24

	source := testImage(width, height)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, options)
	if err != nil {
		t.Fatal(err)
	}
	observed := experimentalV4Build37ScannerPageForTest(marked)
	got, info, scanner, err := ExperimentalV4ExtractScanner(observed, key, width, height)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload=%q want=%q", got, payload)
	}
	if !scanner.BoundaryDetected || !scanner.Accepted || scanner.EnsembleCandidates != experimentalV4ScannerEnsembleSize {
		t.Fatalf("scanner info=%+v", scanner)
	}
	if info.OriginXBlocks != 0 || info.OriginYBlocks != 0 {
		t.Fatalf("origin=(%d,%d) want (0,0)", info.OriginXBlocks, info.OriginYBlocks)
	}

	control := experimentalV4Build37ScannerPageForTest(source)
	_, controlInfo, controlScanner, err := ExperimentalV4ExtractScanner(control, key, width, height)
	if err == nil {
		t.Fatal("unmarked scanner control unexpectedly authenticated")
	}
	if controlScanner.Accepted {
		t.Fatalf("unmarked control passed scanner geometry gate: %+v extract=%+v", controlScanner, controlInfo)
	}
}

func TestExperimentalV4Build37ScannerRejectsInvalidCanonicalExtent(t *testing.T) {
	_, _, _, err := ExperimentalV4ExtractScanner(testImage(320, 320), []byte("12345678"), 319, 320)
	if err == nil {
		t.Fatal("accepted non-block-aligned canonical extent")
	}
}
