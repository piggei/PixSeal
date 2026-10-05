package watermark

import (
	"image"
	"math"
	"testing"
)

// Build83 is observability-only. These benchmarks exercise public deterministic
// fixtures and never inspect payload, ECC, HMAC, a secret key, or a private
// acquisition. The historical Build76 projective reader is authoritative; the
// closed Build82 exact-inline reader is retained only as a benchmark comparator.

var experimentalV4Build83BenchmarkFloat float64
var experimentalV4Build83BenchmarkBool bool

type experimentalV4Build83SamplerFixture struct {
	plane           *pixelPlane
	front           homography
	mild            homography
	angle           homography
	origins         [][2]int
	canonicalPoints [][2]float64
	anglePoints     [][2]float64
	dctValues       [blockSize * blockSize]float64
}

func experimentalV4Build83Fixture(tb testing.TB) experimentalV4Build83SamplerFixture {
	tb.Helper()
	const width, height = 1632, 1632
	plane := &pixelPlane{bounds: image.Rect(0, 0, width, height), rgb: make([]uint8, width*height*3)}
	// Public deterministic texture. The benchmark must not depend on a private
	// photo, embedded payload, pilot outcome, or secret-derived pixels.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 3
			plane.rgb[i] = uint8((x*13 + y*7 + 17) & 0xff)
			plane.rgb[i+1] = uint8((x*5 + y*19 + 73) & 0xff)
			plane.rgb[i+2] = uint8((x*23 + y*3 + 151) & 0xff)
		}
	}
	makeHom := func(name string, quad [4][2]float64) homography {
		h, ok := homographyForQuad(width, height, quad)
		if !ok {
			tb.Fatalf("Build83 %s homography rejected", name)
		}
		return h
	}
	front := makeHom("front", [4][2]float64{{0, 0}, {1, 0}, {0, 1}, {1, 1}})
	mild := makeHom("mild", [4][2]float64{{0.020, 0.012}, {0.982, 0.030}, {0.015, 0.985}, {0.986, 0.972}})
	angle := makeHom("angle", [4][2]float64{{0.095, 0.035}, {0.925, 0.105}, {0.030, 0.905}, {0.975, 0.965}})

	origins := make([][2]int, 0, 64)
	for y := 384; y <= 1152; y += 128 {
		for x := 384; x <= 1152; x += 128 {
			origins = append(origins, [2]int{x, y})
		}
	}
	canonicalPoints := make([][2]float64, 0, len(origins)*blockSize*blockSize)
	anglePoints := make([][2]float64, 0, len(origins)*blockSize*blockSize)
	for _, origin := range origins {
		want, wok := readProjectiveBlockValue(plane, angle, origin[0], origin[1], blockSize)
		got, gok := experimentalV4PhoneBuild82ReadProjectiveBlockValue(plane, angle, origin[0], origin[1], blockSize)
		if !wok || !gok || math.Float64bits(want) != math.Float64bits(got) {
			tb.Fatalf("Build83 fixture exact-reader mismatch at %v ok=%t/%t bits=%x/%x", origin, wok, gok, math.Float64bits(want), math.Float64bits(got))
		}
		for by := 0; by < blockSize; by++ {
			for bx := 0; bx < blockSize; bx++ {
				cx, cy := float64(origin[0]+bx), float64(origin[1]+by)
				sx, sy, ok := angle.mapPoint(cx, cy)
				if !ok {
					tb.Fatalf("Build83 fixture mapPoint rejected at %v %d,%d", origin, bx, by)
				}
				if _, ok := samplePlaneLuminance(plane, sx, sy); !ok {
					tb.Fatalf("Build83 fixture luminance rejected at %.17g,%.17g", sx, sy)
				}
				canonicalPoints = append(canonicalPoints, [2]float64{cx, cy})
				anglePoints = append(anglePoints, [2]float64{sx, sy})
			}
		}
	}
	first := origins[0]
	var dctValues [blockSize * blockSize]float64
	idx := 0
	for by := 0; by < blockSize; by++ {
		for bx := 0; bx < blockSize; bx++ {
			sx, sy, _ := angle.mapPoint(float64(first[0]+bx), float64(first[1]+by))
			dctValues[idx], _ = samplePlaneLuminance(plane, sx, sy)
			idx++
		}
	}
	return experimentalV4Build83SamplerFixture{plane: plane, front: front, mild: mild, angle: angle, origins: origins, canonicalPoints: canonicalPoints, anglePoints: anglePoints, dctValues: dctValues}
}

func TestExperimentalV4Build83PublicBenchmarkFixture(t *testing.T) {
	f := experimentalV4Build83Fixture(t)
	for name, h := range map[string]homography{"front": f.front, "mild": f.mild, "angle": f.angle} {
		for _, origin := range f.origins {
			want, wok := readProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			got, gok := experimentalV4PhoneBuild82ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
				t.Fatalf("%s origin=%v exact mismatch ok=%t/%t bits=%x/%x", name, origin, gok, wok, math.Float64bits(got), math.Float64bits(want))
			}
		}
	}
	if len(f.anglePoints) != len(f.origins)*blockSize*blockSize {
		t.Fatalf("Build83 point count=%d", len(f.anglePoints))
	}
}

func BenchmarkExperimentalV4Build83HistoricalBlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		value, ok = readProjectiveBlockValue(f.plane, f.front, origin[0], origin[1], blockSize)
	}
	experimentalV4Build83BenchmarkFloat = value
	experimentalV4Build83BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build83HistoricalBlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		value, ok = readProjectiveBlockValue(f.plane, f.mild, origin[0], origin[1], blockSize)
	}
	experimentalV4Build83BenchmarkFloat = value
	experimentalV4Build83BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build83HistoricalBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		value, ok = readProjectiveBlockValue(f.plane, f.angle, origin[0], origin[1], blockSize)
	}
	experimentalV4Build83BenchmarkFloat = value
	experimentalV4Build83BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build83InlineBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		value, ok = experimentalV4PhoneBuild82ReadProjectiveBlockValue(f.plane, f.angle, origin[0], origin[1], blockSize)
	}
	experimentalV4Build83BenchmarkFloat = value
	experimentalV4Build83BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build83MapPointAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var sx, sy float64
	var ok bool
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			p := f.canonicalPoints[base+j]
			// One benchmark operation is one complete 8x8 block worth of
			// mapPoint calls, making stage ns/op directly comparable by block.
			sx, sy, ok = f.angle.mapPoint(p[0], p[1])
		}
	}
	experimentalV4Build83BenchmarkFloat = sx + sy
	experimentalV4Build83BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build83SamplePlaneLuminanceAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			p := f.anglePoints[base+j]
			// One benchmark operation is one complete 8x8 block worth of
			// bilinear luminance samples, matching the reader's block granularity.
			value, ok = samplePlaneLuminance(f.plane, p[0], p[1])
		}
	}
	experimentalV4Build83BenchmarkFloat = value
	experimentalV4Build83BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build83DCTAccumulationAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	table := readCosTables[blockSize]
	table2, table3 := table[2], table[3]
	b.ReportAllocs()
	b.ResetTimer()
	var out float64
	for i := 0; i < b.N; i++ {
		c23, c32 := 0.0, 0.0
		idx := 0
		for y := 0; y < blockSize; y++ {
			for x := 0; x < blockSize; x++ {
				l := f.dctValues[idx]
				c23 += l * table3[x] * table2[y]
				c32 += l * table2[x] * table3[y]
				idx++
			}
		}
		out = math.Abs(c23) - math.Abs(c32)
	}
	experimentalV4Build83BenchmarkFloat = out
}
