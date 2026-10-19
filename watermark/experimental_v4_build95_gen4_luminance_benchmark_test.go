package watermark

import (
	"image"
	"runtime"
	"testing"
)

var experimentalV4Build95BenchmarkFloat float64
var experimentalV4Build95BenchmarkBool bool
var experimentalV4Build95BenchmarkInt int
var experimentalV4Build95BenchmarkPlane *experimentalV4PhoneBuild95LuminancePlane

func TestExperimentalV4Build95PublicBenchmarkFixture(t *testing.T) {
	TestExperimentalV4Build95ProjectiveBlockMatchesBuild84BitExactly(t)
	TestExperimentalV4Build95ContinueMatchesBuild84Exactly(t)
}

func benchmarkExperimentalV4Build95Reader(b *testing.B, candidate bool) {
	f := experimentalV4Build83Fixture(b)
	luma, _, fallback := experimentalV4PhoneBuild95NewLuminancePlane(f.plane)
	if fallback {
		b.Fatal("unexpected luminance-plane fallback")
	}
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		if candidate {
			value, ok = experimentalV4PhoneBuild95ReadProjectiveBlockValue(luma, f.angle, origin[0], origin[1], blockSize)
		} else {
			value, ok = experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, f.angle, origin[0], origin[1], blockSize)
		}
	}
	experimentalV4Build95BenchmarkFloat = value
	experimentalV4Build95BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build95Build84BlockReaderAngle(b *testing.B) {
	benchmarkExperimentalV4Build95Reader(b, false)
}

func BenchmarkExperimentalV4Build95LuminanceBlockReaderAngle(b *testing.B) {
	benchmarkExperimentalV4Build95Reader(b, true)
}

func BenchmarkExperimentalV4Build95LuminancePlanePreparation(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		luma, _, fallback := experimentalV4PhoneBuild95NewLuminancePlane(f.plane)
		if fallback {
			b.Fatal("unexpected luminance-plane fallback")
		}
		experimentalV4Build95BenchmarkInt = len(luma.luma)
		experimentalV4Build95BenchmarkPlane = luma
	}
}

func BenchmarkExperimentalV4Build95RetainedSizeLuminancePlanePreparation(b *testing.B) {
	const width, height = 3464, 4600
	plane := &pixelPlane{
		bounds: image.Rect(0, 0, width, height),
		rgb:    make([]uint8, width*height*3),
	}
	for i := range plane.rgb {
		plane.rgb[i] = uint8((i*37 + 11) & 0xff)
	}
	wantBytes := int64(width) * int64(height) * 8
	if wantBytes > experimentalV4PhoneBuild95MaxLuminancePlaneBytes {
		b.Fatalf("retained-size plane bytes=%d exceeds cap=%d", wantBytes, experimentalV4PhoneBuild95MaxLuminancePlaneBytes)
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(plane.rgb)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		luma, bytes, fallback := experimentalV4PhoneBuild95NewLuminancePlane(plane)
		if fallback || bytes != wantBytes {
			b.Fatalf("retained-size preparation fallback=%t bytes=%d want=%d", fallback, bytes, wantBytes)
		}
		experimentalV4Build95BenchmarkPlane = luma
	}
	runtime.KeepAlive(plane)
}

func benchmarkExperimentalV4Build95Continue(b *testing.B, candidate bool) {
	pilot := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, pilot, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		b.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	plane := newPixelPlane(carrier)
	anchor := [4]ImagePoint{{X: 0, Y: 0}, {X: float64(w - 1), Y: 0}, {X: 0, Y: float64(h - 1)}, {X: float64(w - 1), Y: float64(h - 1)}}
	hom, ok := experimentalV4PhoneQuadHomography(w, h, anchor)
	if !ok {
		b.Fatal("anchor homography")
	}
	luma, _, fallback := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	if fallback {
		b.Fatal("unexpected luminance-plane fallback")
	}
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, pilot, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if candidate {
			states, evals, fetch := experimentalV4PhoneBuild95ContinueFetch(luma, pilot, w, h, anchor, start, 0)
			experimentalV4Build95BenchmarkInt = len(states) + evals + fetch.BlockSuccess
		} else {
			states, evals, fetch := experimentalV4PhoneBuild84ContinueFetch(plane, pilot, w, h, anchor, start, 0)
			experimentalV4Build95BenchmarkInt = len(states) + evals + fetch.BlockSuccess
		}
	}
}

func BenchmarkExperimentalV4Build95Build84Continuation(b *testing.B) {
	benchmarkExperimentalV4Build95Continue(b, false)
}

func BenchmarkExperimentalV4Build95LuminanceContinuation(b *testing.B) {
	benchmarkExperimentalV4Build95Continue(b, true)
}
