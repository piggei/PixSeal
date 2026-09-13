package watermark

import (
	"image"
	"image/color"
	"math"
	"testing"
)

const experimentalV4Build25SeparationFloor = 0.10

type experimentalV4GeometryCase struct {
	name string
	h    homography
	w    int
	hgt  int
	post func(*testing.T, image.Image) image.Image
}

func TestExperimentalV4Prototype2KnownGeometryQualification(t *testing.T) {
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(888, 768)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 24, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range experimentalV4Build25GeometryCases(t, source.Bounds().Dx(), source.Bounds().Dy()) {
		t.Run(tc.name, func(t *testing.T) {
			positiveImage := experimentalV4WarpHomographyForTest(carrier, tc.h, tc.w, tc.hgt)
			negativeImage := experimentalV4WarpHomographyForTest(source, tc.h, tc.w, tc.hgt)
			if tc.post != nil {
				positiveImage = toNRGBA(tc.post(t, positiveImage))
				negativeImage = toNRGBA(tc.post(t, negativeImage))
			}

			positive := experimentalV4DetectPilotProjective(positiveImage, candidate, source.Bounds().Dx(), source.Bounds().Dy(), tc.h)
			negative := experimentalV4DetectPilotProjective(negativeImage, candidate, source.Bounds().Dx(), source.Bounds().Dy(), tc.h)
			t.Logf("positive score=%.6f runner=%.6f margin=%.6f visible=%d samples=%d origin=(%d,%d); negative score=%.6f runner=%.6f margin=%.6f origin=(%d,%d)",
				positive.Score, positive.RunnerUpScore, positive.Margin, positive.VisiblePilotPositions, positive.PilotSamples, positive.OriginXBlocks, positive.OriginYBlocks,
				negative.Score, negative.RunnerUpScore, negative.Margin, negative.OriginXBlocks, negative.OriginYBlocks)

			if !positive.Available {
				t.Fatal("pilot unavailable")
			}
			if positive.OriginXBlocks != 0 || positive.OriginYBlocks != 0 {
				t.Fatalf("origin=(%d,%d), want (0,0); detection=%+v", positive.OriginXBlocks, positive.OriginYBlocks, positive)
			}
			if positive.Margin <= 0 {
				t.Fatalf("non-positive pilot runner-up margin: %+v", positive)
			}
			if positive.Margin-negative.Margin < experimentalV4Build25SeparationFloor {
				t.Fatalf("pilot/negative margin separation %.6f below Build25 development floor %.2f (positive %.6f, negative %.6f)", positive.Margin-negative.Margin, experimentalV4Build25SeparationFloor, positive.Margin, negative.Margin)
			}
		})
	}
}

func experimentalV4Build25GeometryCases(t *testing.T, width, height int) []experimentalV4GeometryCase {
	t.Helper()
	rotate12, rw12, rh12 := experimentalV4AffineHomographyForTest(width, height, math.Cos(12.3*math.Pi/180), -math.Sin(12.3*math.Pi/180), math.Sin(12.3*math.Pi/180), math.Cos(12.3*math.Pi/180))
	rotateNeg17, rw17, rh17 := experimentalV4AffineHomographyForTest(width, height, math.Cos(-17*math.Pi/180), -math.Sin(-17*math.Pi/180), math.Sin(-17*math.Pi/180), math.Cos(-17*math.Pi/180))
	scale110x90, sw1, sh1 := experimentalV4AffineHomographyForTest(width, height, 1.10, 0, 0, 0.90)
	scale90x110, sw2, sh2 := experimentalV4AffineHomographyForTest(width, height, 0.90, 0, 0, 1.10)
	shearX8, sxw, sxh := experimentalV4AffineHomographyForTest(width, height, 1, math.Tan(8*math.Pi/180), 0, 1)
	shearY8, syw, syh := experimentalV4AffineHomographyForTest(width, height, 1, 0, math.Tan(-8*math.Pi/180), 1)
	combined, cw, ch := experimentalV4AffineHomographyForTest(width, height,
		1.10*math.Cos(12.3*math.Pi/180), -0.90*math.Sin(12.3*math.Pi/180),
		1.10*math.Sin(12.3*math.Pi/180), 0.90*math.Cos(12.3*math.Pi/180))
	combinedShear, csw, csh := experimentalV4AffineHomographyForTest(width, height,
		math.Cos(-9*math.Pi/180), math.Cos(-9*math.Pi/180)*math.Tan(5*math.Pi/180)-math.Sin(-9*math.Pi/180),
		math.Sin(-9*math.Pi/180), math.Sin(-9*math.Pi/180)*math.Tan(5*math.Pi/180)+math.Cos(-9*math.Pi/180))
	combined75, c75w, c75h := experimentalV4AffineHomographyForTest(width, height,
		0.75*math.Cos(12.3*math.Pi/180), -0.75*math.Sin(12.3*math.Pi/180),
		0.75*math.Sin(12.3*math.Pi/180), 0.75*math.Cos(12.3*math.Pi/180))

	perspectiveTop := experimentalV4PerspectiveHomographyForTest(t, width, height, [4][2]float64{
		{0.04 * float64(width-1), 0.01 * float64(height-1)},
		{0.96 * float64(width-1), 0.00 * float64(height-1)},
		{0.00 * float64(width-1), 1.00 * float64(height-1)},
		{1.00 * float64(width-1), 0.99 * float64(height-1)},
	})
	perspectiveBottom := experimentalV4PerspectiveHomographyForTest(t, width, height, [4][2]float64{
		{0.00 * float64(width-1), 0.00 * float64(height-1)},
		{1.00 * float64(width-1), 0.01 * float64(height-1)},
		{0.05 * float64(width-1), 0.99 * float64(height-1)},
		{0.95 * float64(width-1), 1.00 * float64(height-1)},
	})
	postPerspective := experimentalV4PerspectiveHomographyForTest(t, cw, ch, [4][2]float64{
		{0.03 * float64(cw-1), 0.01 * float64(ch-1)},
		{0.97 * float64(cw-1), 0.00 * float64(ch-1)},
		{0.00 * float64(cw-1), 0.99 * float64(ch-1)},
		{1.00 * float64(cw-1), 1.00 * float64(ch-1)},
	})
	combinedPerspective := experimentalV4MultiplyHomography(postPerspective, combined)
	cropLeft, cropTop := cw/12, ch/14
	cropRight, cropBottom := cw/14, ch/12
	combinedCrop, cropW, cropH := experimentalV4CropHomographyForTest(combinedPerspective, cw, ch, cropLeft, cropTop, cropRight, cropBottom)

	return []experimentalV4GeometryCase{
		{name: "rotate-12.3", h: rotate12, w: rw12, hgt: rh12},
		{name: "rotate--17", h: rotateNeg17, w: rw17, hgt: rh17},
		{name: "scale110x90", h: scale110x90, w: sw1, hgt: sh1},
		{name: "scale90x110", h: scale90x110, w: sw2, hgt: sh2},
		{name: "shearx8", h: shearX8, w: sxw, hgt: sxh},
		{name: "sheary8", h: shearY8, w: syw, hgt: syh},
		{name: "perspective-top", h: perspectiveTop, w: width, hgt: height},
		{name: "perspective-bottom", h: perspectiveBottom, w: width, hgt: height},
		{name: "rotate12.3-scale110x90", h: combined, w: cw, hgt: ch},
		{name: "rotate--9-shearx5", h: combinedShear, w: csw, hgt: csh},
		{name: "rotate12.3-resize75", h: combined75, w: c75w, hgt: c75h},
		{name: "combined-perspective-jpeg82", h: combinedPerspective, w: cw, hgt: ch, post: func(t *testing.T, src image.Image) image.Image { return experimentalV4JPEGForTest(t, src, 82) }},
		{name: "combined-perspective-blur", h: combinedPerspective, w: cw, hgt: ch, post: func(_ *testing.T, src image.Image) image.Image { return experimentalV4BoxBlurForTest(src) }},
		{name: "combined-perspective-noise4", h: combinedPerspective, w: cw, hgt: ch, post: func(_ *testing.T, src image.Image) image.Image { return experimentalV4NoiseForTest(src, 4) }},
		{name: "combined-perspective-crop", h: combinedCrop, w: cropW, hgt: cropH},
		{name: "combined-perspective-crop-jpeg82", h: combinedCrop, w: cropW, hgt: cropH, post: func(t *testing.T, src image.Image) image.Image { return experimentalV4JPEGForTest(t, src, 82) }},
	}
}

func experimentalV4AffineHomographyForTest(width, height int, a, b, c, d float64) (homography, int, int) {
	base := homography{h: [9]float64{a, b, 0, c, d, 0, 0, 0, 1}}
	corners := [4][2]float64{{0, 0}, {float64(width - 1), 0}, {0, float64(height - 1)}, {float64(width - 1), float64(height - 1)}}
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, p := range corners {
		x, y, _ := base.mapPoint(p[0], p[1])
		minX = math.Min(minX, x)
		minY = math.Min(minY, y)
		maxX = math.Max(maxX, x)
		maxY = math.Max(maxY, y)
	}
	const margin = 3.0
	translate := homography{h: [9]float64{1, 0, -minX + margin, 0, 1, -minY + margin, 0, 0, 1}}
	h := experimentalV4MultiplyHomography(translate, base)
	outWidth := int(math.Ceil(maxX - minX + 1 + 2*margin))
	outHeight := int(math.Ceil(maxY - minY + 1 + 2*margin))
	return h, outWidth, outHeight
}

func experimentalV4CropHomographyForTest(h homography, width, height, left, top, right, bottom int) (homography, int, int) {
	if left < 0 || top < 0 || right < 0 || bottom < 0 || left+right >= width || top+bottom >= height {
		panic("invalid Build25 crop")
	}
	translate := homography{h: [9]float64{1, 0, -float64(left), 0, 1, -float64(top), 0, 0, 1}}
	return experimentalV4MultiplyHomography(translate, h), width - left - right, height - top - bottom
}

func experimentalV4PerspectiveHomographyForTest(t *testing.T, width, height int, dst [4][2]float64) homography {
	t.Helper()
	src := [4][2]float64{{0, 0}, {float64(width - 1), 0}, {0, float64(height - 1)}, {float64(width - 1), float64(height - 1)}}
	h, ok := diagnosticHomographyFromFourPoints(src, dst)
	if !ok {
		t.Fatal("cannot construct Build25 perspective homography")
	}
	return h
}

// experimentalV4MultiplyHomography returns a*b, so b is applied first.
func experimentalV4MultiplyHomography(a, b homography) homography {
	var out [9]float64
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			for k := 0; k < 3; k++ {
				out[row*3+col] += a.h[row*3+k] * b.h[k*3+col]
			}
		}
	}
	return homography{h: out}
}

func experimentalV4WarpHomographyForTest(src image.Image, h homography, width, height int) *image.NRGBA {
	inverse, ok := invertDiagnosticHomography(h)
	if !ok {
		panic("non-invertible Build25 test homography")
	}
	out := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sx, sy, ok := inverse.mapPoint(float64(x), float64(y))
			if !ok {
				out.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
				continue
			}
			out.SetNRGBA(x, y, sampleImageNRGBAForTest(src, sx, sy))
		}
	}
	return out
}
