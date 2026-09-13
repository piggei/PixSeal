package watermark

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"testing"
)

func TestExperimentalV4Prototype2ImageDomainPilotChannel(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		image       image.Image
		blockSize   int
		wantOriginX int
		wantOriginY int
	}{
		{name: "native", image: carrier, blockSize: 8, wantOriginX: 0, wantOriginY: 0},
		{name: "jpeg-82", image: experimentalV4JPEGForTest(t, carrier, 82), blockSize: 8, wantOriginX: 0, wantOriginY: 0},
		{name: "blur-1", image: experimentalV4BoxBlurForTest(carrier), blockSize: 8, wantOriginX: 0, wantOriginY: 0},
		{name: "noise-4", image: experimentalV4NoiseForTest(carrier, 4), blockSize: 8, wantOriginX: 0, wantOriginY: 0},
		{name: "gamma-1.15", image: experimentalV4GammaForTest(carrier, 1.15), blockSize: 8, wantOriginX: 0, wantOriginY: 0},
		{name: "resize-75", image: resizeBilinear(carrier, carrier.Bounds().Dx()*3/4, carrier.Bounds().Dy()*3/4), blockSize: 6, wantOriginX: 0, wantOriginY: 0},
		{name: "resize-50", image: resizeBilinear(carrier, carrier.Bounds().Dx()/2, carrier.Bounds().Dy()/2), blockSize: 4, wantOriginX: 0, wantOriginY: 0},
		{name: "aligned-crop", image: cropCopy(carrier, image.Rect(5*8, 3*8, 105*8, 91*8)), blockSize: 8, wantOriginX: 5, wantOriginY: 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			detection := experimentalV4DetectPilotGrid(tc.image, candidate, tc.blockSize, 0, 0)
			t.Logf("detection: %+v", detection)
			if !detection.Available {
				t.Fatal("pilot unavailable")
			}
			if detection.OriginXBlocks != tc.wantOriginX || detection.OriginYBlocks != tc.wantOriginY {
				t.Fatalf("origin=(%d,%d), want (%d,%d); detection=%+v", detection.OriginXBlocks, detection.OriginYBlocks, tc.wantOriginX, tc.wantOriginY, detection)
			}
			if detection.Margin <= 0 {
				t.Fatalf("non-positive pilot runner-up margin: %+v", detection)
			}
		})
	}

	negative := experimentalV4DetectPilotGrid(source, candidate, 8, 0, 0)
	t.Logf("unmarked negative control: %+v", negative)
	if !negative.Available {
		t.Fatal("negative control telemetry unavailable")
	}
}

func experimentalV4JPEGForTest(t *testing.T, src image.Image, quality int) image.Image {
	t.Helper()
	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, src, &jpeg.Options{Quality: quality}); err != nil {
		t.Fatal(err)
	}
	decoded, err := jpeg.Decode(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	return decoded
}

func experimentalV4BoxBlurForTest(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			var r, g, b, count int
			for dy := -1; dy <= 1; dy++ {
				sy := y + dy
				if sy < 0 || sy >= bounds.Dy() {
					continue
				}
				for dx := -1; dx <= 1; dx++ {
					sx := x + dx
					if sx < 0 || sx >= bounds.Dx() {
						continue
					}
					c := color.NRGBAModel.Convert(src.At(bounds.Min.X+sx, bounds.Min.Y+sy)).(color.NRGBA)
					r += int(c.R)
					g += int(c.G)
					b += int(c.B)
					count++
				}
			}
			out.SetNRGBA(x, y, color.NRGBA{R: uint8(r / count), G: uint8(g / count), B: uint8(b / count), A: 255})
		}
	}
	return out
}

func experimentalV4NoiseForTest(src image.Image, amplitude int) *image.NRGBA {
	bounds := src.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	state := uint64(0x76346e6f69736531)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			state ^= state << 13
			state ^= state >> 7
			state ^= state << 17
			delta := int(state%uint64(2*amplitude+1)) - amplitude
			c := color.NRGBAModel.Convert(src.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)
			out.SetNRGBA(x, y, color.NRGBA{
				R: experimentalV4ClampByte(int(c.R) + delta),
				G: experimentalV4ClampByte(int(c.G) + delta),
				B: experimentalV4ClampByte(int(c.B) + delta),
				A: 255,
			})
		}
	}
	return out
}

func experimentalV4GammaForTest(src image.Image, gamma float64) *image.NRGBA {
	bounds := src.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	inverse := 1 / gamma
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			c := color.NRGBAModel.Convert(src.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)
			convert := func(value uint8) uint8 {
				return experimentalV4ClampByte(int(math.Round(255 * math.Pow(float64(value)/255, inverse))))
			}
			out.SetNRGBA(x, y, color.NRGBA{R: convert(c.R), G: convert(c.G), B: convert(c.B), A: 255})
		}
	}
	return out
}

func experimentalV4ClampByte(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}
