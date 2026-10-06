package watermark

import (
	"math"
	"testing"
)

// Build85 is observability-only over the qualified Build84 smartphone baseline.
// All work uses the public deterministic Build83 sampler fixture. No private
// acquisition, key, payload, ECC result, HMAC result, or expected message selects
// the benchmark workload.

var experimentalV4Build85BenchmarkFloat float64
var experimentalV4Build85BenchmarkBool bool
var experimentalV4Build85BenchmarkInt int

type experimentalV4Build85PreparedSample struct {
	index00 int
	index10 int
	index01 int
	index11 int
	fx      float64
	fy      float64
}

func experimentalV4Build85PreparedAngleSamples(tb testing.TB) (experimentalV4Build83SamplerFixture, []experimentalV4Build85PreparedSample) {
	tb.Helper()
	f := experimentalV4Build83Fixture(tb)
	width, height := f.plane.bounds.Dx(), f.plane.bounds.Dy()
	rowStride := width * 3
	out := make([]experimentalV4Build85PreparedSample, 0, len(f.anglePoints))
	for _, p := range f.anglePoints {
		sx, sy := p[0], p[1]
		if sx < 0 || sy < 0 || sx > float64(width-1) || sy > float64(height-1) {
			tb.Fatalf("Build85 prepared point outside plane: %.17g,%.17g", sx, sy)
		}
		x0, y0 := int(math.Floor(sx)), int(math.Floor(sy))
		x1, y1 := x0+1, y0+1
		if x1 >= width {
			x1 = width - 1
		}
		if y1 >= height {
			y1 = height - 1
		}
		row0 := y0 * rowStride
		row1 := y1 * rowStride
		x03 := x0 * 3
		x13 := x1 * 3
		out = append(out, experimentalV4Build85PreparedSample{
			index00: row0 + x03,
			index10: row0 + x13,
			index01: row1 + x03,
			index11: row1 + x13,
			fx:      sx - float64(x0),
			fy:      sy - float64(y0),
		})
	}
	return f, out
}

func experimentalV4Build85PreparedLuminance(rgb []uint8, s experimentalV4Build85PreparedSample) float64 {
	pixel00 := rgb[s.index00 : s.index00+3]
	pixel10 := rgb[s.index10 : s.index10+3]
	pixel01 := rgb[s.index01 : s.index01+3]
	pixel11 := rgb[s.index11 : s.index11+3]
	l00 := .299*float64(pixel00[0]) + .587*float64(pixel00[1]) + .114*float64(pixel00[2]) - 128
	l10 := .299*float64(pixel10[0]) + .587*float64(pixel10[1]) + .114*float64(pixel10[2]) - 128
	l01 := .299*float64(pixel01[0]) + .587*float64(pixel01[1]) + .114*float64(pixel01[2]) - 128
	l11 := .299*float64(pixel11[0]) + .587*float64(pixel11[1]) + .114*float64(pixel11[2]) - 128
	top := l00*(1-s.fx) + l10*s.fx
	bottom := l01*(1-s.fx) + l11*s.fx
	return top*(1-s.fy) + bottom*s.fy
}

func TestExperimentalV4Build85PublicProfileFixture(t *testing.T) {
	f, prepared := experimentalV4Build85PreparedAngleSamples(t)
	if len(prepared) != len(f.anglePoints) {
		t.Fatalf("prepared sample count=%d want=%d", len(prepared), len(f.anglePoints))
	}
	for i, p := range f.anglePoints {
		want, ok := samplePlaneLuminance(f.plane, p[0], p[1])
		if !ok {
			t.Fatalf("historical luminance rejected sample %d", i)
		}
		got := experimentalV4Build85PreparedLuminance(f.plane.rgb, prepared[i])
		if math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("prepared luminance mismatch sample=%d bits=%x/%x", i, math.Float64bits(got), math.Float64bits(want))
		}
	}
	for name, h := range map[string]homography{"front": f.front, "mild": f.mild, "angle": f.angle} {
		for _, origin := range f.origins {
			historical, hok := readProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			qualified, qok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			if hok != qok || (hok && math.Float64bits(historical) != math.Float64bits(qualified)) {
				t.Fatalf("%s origin=%v qualified mismatch ok=%t/%t bits=%x/%x", name, origin, hok, qok, math.Float64bits(historical), math.Float64bits(qualified))
			}
		}
	}
}

func experimentalV4Build85BenchmarkQualifiedReader(b *testing.B, h homography) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		value, ok = experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
	}
	experimentalV4Build85BenchmarkFloat = value
	experimentalV4Build85BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build85QualifiedBlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build85BenchmarkQualifiedReader(b, f.front)
}
func BenchmarkExperimentalV4Build85QualifiedBlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build85BenchmarkQualifiedReader(b, f.mild)
}
func BenchmarkExperimentalV4Build85QualifiedBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build85BenchmarkQualifiedReader(b, f.angle)
}

func BenchmarkExperimentalV4Build85MapPointAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var sx, sy float64
	var ok bool
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			p := f.canonicalPoints[base+j]
			sx, sy, ok = f.angle.mapPoint(p[0], p[1])
		}
	}
	experimentalV4Build85BenchmarkFloat = sx + sy
	experimentalV4Build85BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build85AddressFloorClampAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	width, height := f.plane.bounds.Dx(), f.plane.bounds.Dy()
	rowStride := width * 3
	b.ReportAllocs()
	b.ResetTimer()
	var sink int
	var frac float64
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			p := f.anglePoints[base+j]
			sx, sy := p[0], p[1]
			if sx < 0 || sy < 0 || sx > float64(width-1) || sy > float64(height-1) {
				continue
			}
			x0, y0 := int(math.Floor(sx)), int(math.Floor(sy))
			x1, y1 := x0+1, y0+1
			if x1 >= width {
				x1 = width - 1
			}
			if y1 >= height {
				y1 = height - 1
			}
			row0 := y0 * rowStride
			row1 := y1 * rowStride
			x03 := x0 * 3
			x13 := x1 * 3
			sink = row0 + row1 + x03 + x13
			frac = (sx - float64(x0)) + (sy - float64(y0))
		}
	}
	experimentalV4Build85BenchmarkInt = sink
	experimentalV4Build85BenchmarkFloat = frac
}

func BenchmarkExperimentalV4Build85PreparedFetchLumaBilinearAngle(b *testing.B) {
	f, prepared := experimentalV4Build85PreparedAngleSamples(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			value = experimentalV4Build85PreparedLuminance(f.plane.rgb, prepared[base+j])
		}
	}
	experimentalV4Build85BenchmarkFloat = value
}

func BenchmarkExperimentalV4Build85DCTAccumulationAngle(b *testing.B) {
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
	experimentalV4Build85BenchmarkFloat = out
}
