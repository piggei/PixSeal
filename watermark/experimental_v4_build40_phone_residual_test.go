package watermark

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func experimentalV4Build40WarpForTest(src image.Image, warp *experimentalV4PhoneResidualWarp) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	plane := newPixelPlane(src)
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Invert z = c + residual(c) with a few fixed-point iterations.
			zx, zy := float64(x), float64(y)
			cx, cy := zx, zy
			for i := 0; i < 6; i++ {
				dx, dy := warp.correction(cx, cy)
				cx, cy = zx-dx, zy-dy
			}
			l, ok := samplePlaneLuminance(plane, cx, cy)
			if !ok {
				dst.SetNRGBA(x, y, color.NRGBA{128, 128, 128, 255})
				continue
			}
			v := clamp(l + 128)
			dst.SetNRGBA(x, y, color.NRGBA{v, v, v, 255})
		}
	}
	return dst
}

func experimentalV4Build40IdentityHomography() homography {
	return homography{h: [9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}}
}

func TestExperimentalV4Build40PilotOnlyResidualWarp(t *testing.T) {
	const width = 1480
	const height = 1280
	key := []byte("PixSeal-v4-TestKey-2026")
	payload := []byte("v4-b40-residual")
	opts := DefaultOptions()
	opts.Profile = ProfileRobust
	opts.Strength = 48
	source := testImage(width, height)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, opts)
	if err != nil {
		t.Fatal(err)
	}
	truth := &experimentalV4PhoneResidualWarp{
		width: float64(width), height: float64(height), maxPixels: experimentalV4PhoneResidualMaxPixels,
		dx: [6]float64{1.4, 2.0, -1.5, 1.1, 0.8, -0.6},
		dy: [6]float64{-1.2, 0.9, 1.8, -1.0, 0.4, 0.8},
	}
	observed := experimentalV4Build40WarpForTest(marked, truth)
	plane := newPixelPlane(observed)
	candidate := experimentalV4Prototype2Candidate()
	base := experimentalV4Build40IdentityHomography()

	warp, residual := experimentalV4PhoneFitResidualPilotOnly(plane, candidate, width, height, base)
	if warp == nil || !residual.Applied {
		t.Fatalf("residual warp not accepted: %+v", residual)
	}
	if residual.Controls < experimentalV4PhoneResidualMinControls {
		t.Fatalf("too few controls: %+v", residual)
	}
	if residual.ValidationAfter <= residual.ValidationBefore {
		t.Fatalf("held-out validation did not improve: before=%.4f after=%.4f", residual.ValidationBefore, residual.ValidationAfter)
	}

	var got []byte
	for _, spec := range v3Profiles {
		margins, _, ok := experimentalV4ReadProtectedPhoneWarpMargins(plane, width, height, base, warp, spec.codedBits)
		if !ok {
			continue
		}
		soft := experimentalV4SoftHammingDecodeMargins(margins)
		raw := bitsToBytes(whiten(soft, key, experimentalV4WhitenLabel))
		if decoded, err := parseExperimentalV4Frame(raw, key, spec); err == nil {
			got = decoded
			break
		}
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("authenticated payload mismatch: got=%q residual=%+v", got, residual)
	}

	control := experimentalV4Build40WarpForTest(source, truth)
	controlWarp, controlInfo := experimentalV4PhoneFitResidualPilotOnly(newPixelPlane(control), candidate, width, height, base)
	if controlWarp != nil || controlInfo.Applied {
		t.Fatalf("unmarked control unexpectedly qualified residual warp: %+v", controlInfo)
	}
}
