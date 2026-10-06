package watermark

import (
	"math"
	"testing"
)

// Build87 is observability-only over the qualified Build84 runtime. These
// benchmarks decompose the remaining RGB fetch -> RGB-to-luminance -> bilinear
// hot path on the same public deterministic angle-like fixture used by Builds
// 83-86. They never use a secret key, payload, ECC/HMAC outcome, or private
// acquisition. One benchmark operation always represents one complete 8x8
// block (64 samples), so stage ns/op values are comparable at block granularity.

var experimentalV4Build87BenchmarkFloat float64
var experimentalV4Build87BenchmarkInt int
var experimentalV4Build87BenchmarkBool bool

type experimentalV4Build87PreparedSample struct {
	index00 int
	index10 int
	index01 int
	index11 int
	fx      float64
	fy      float64
	rgb12   [12]uint8
	l00     float64
	l10     float64
	l01     float64
	l11     float64
	top     float64
	bottom  float64
}

func experimentalV4Build87PrepareAngleSamples(tb testing.TB) (experimentalV4Build83SamplerFixture, []experimentalV4Build87PreparedSample) {
	tb.Helper()
	f := experimentalV4Build83Fixture(tb)
	width, height := f.plane.bounds.Dx(), f.plane.bounds.Dy()
	rowStride := width * 3
	out := make([]experimentalV4Build87PreparedSample, 0, len(f.anglePoints))
	for _, p := range f.anglePoints {
		sx, sy := p[0], p[1]
		if sx < 0 || sy < 0 || sx > float64(width-1) || sy > float64(height-1) {
			tb.Fatalf("Build87 prepared point outside plane: %.17g,%.17g", sx, sy)
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
		s := experimentalV4Build87PreparedSample{
			index00: row0 + x03,
			index10: row0 + x13,
			index01: row1 + x03,
			index11: row1 + x13,
			fx:      sx - float64(x0),
			fy:      sy - float64(y0),
		}
		indices := [4]int{s.index00, s.index10, s.index01, s.index11}
		for q, index := range indices {
			copy(s.rgb12[q*3:q*3+3], f.plane.rgb[index:index+3])
		}
		s.l00 = .299*float64(s.rgb12[0]) + .587*float64(s.rgb12[1]) + .114*float64(s.rgb12[2]) - 128
		s.l10 = .299*float64(s.rgb12[3]) + .587*float64(s.rgb12[4]) + .114*float64(s.rgb12[5]) - 128
		s.l01 = .299*float64(s.rgb12[6]) + .587*float64(s.rgb12[7]) + .114*float64(s.rgb12[8]) - 128
		s.l11 = .299*float64(s.rgb12[9]) + .587*float64(s.rgb12[10]) + .114*float64(s.rgb12[11]) - 128
		s.top = s.l00*(1-s.fx) + s.l10*s.fx
		s.bottom = s.l01*(1-s.fx) + s.l11*s.fx
		out = append(out, s)
	}
	return f, out
}

func TestExperimentalV4Build87PublicProfileFixture(t *testing.T) {
	f, prepared := experimentalV4Build87PrepareAngleSamples(t)
	if len(prepared) != len(f.anglePoints) {
		t.Fatalf("Build87 prepared sample count=%d want=%d", len(prepared), len(f.anglePoints))
	}
	for i, p := range f.anglePoints {
		want, ok := samplePlaneLuminance(f.plane, p[0], p[1])
		if !ok {
			t.Fatalf("historical luminance rejected sample %d", i)
		}
		s := prepared[i]
		got := s.top*(1-s.fy) + s.bottom*s.fy
		if math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("Build87 staged luminance mismatch sample=%d bits=%x/%x", i, math.Float64bits(got), math.Float64bits(want))
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

func experimentalV4Build87BenchmarkQualifiedReader(b *testing.B, h homography) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		value, ok = experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
	}
	experimentalV4Build87BenchmarkFloat = value
	experimentalV4Build87BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build87QualifiedBlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build87BenchmarkQualifiedReader(b, f.front)
}
func BenchmarkExperimentalV4Build87QualifiedBlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build87BenchmarkQualifiedReader(b, f.mild)
}
func BenchmarkExperimentalV4Build87QualifiedBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build87BenchmarkQualifiedReader(b, f.angle)
}

func BenchmarkExperimentalV4Build87PreparedPixelSliceFetchAngle(b *testing.B) {
	f, prepared := experimentalV4Build87PrepareAngleSamples(b)
	rgb := f.plane.rgb
	b.ReportAllocs()
	b.ResetTimer()
	var sink int
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		local := 0
		for j := 0; j < blockSize*blockSize; j++ {
			s := prepared[base+j]
			p00 := rgb[s.index00 : s.index00+3]
			p10 := rgb[s.index10 : s.index10+3]
			p01 := rgb[s.index01 : s.index01+3]
			p11 := rgb[s.index11 : s.index11+3]
			local += int(p00[0]) + int(p00[1]) + int(p00[2])
			local += int(p10[0]) + int(p10[1]) + int(p10[2])
			local += int(p01[0]) + int(p01[1]) + int(p01[2])
			local += int(p11[0]) + int(p11[1]) + int(p11[2])
		}
		sink = local
	}
	experimentalV4Build87BenchmarkInt = sink
}

func BenchmarkExperimentalV4Build87PreparedRGBToLuminanceAngle(b *testing.B) {
	f, prepared := experimentalV4Build87PrepareAngleSamples(b)
	b.ReportAllocs()
	b.ResetTimer()
	var l00, l10, l01, l11 float64
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			c := prepared[base+j].rgb12
			l00 = .299*float64(c[0]) + .587*float64(c[1]) + .114*float64(c[2]) - 128
			l10 = .299*float64(c[3]) + .587*float64(c[4]) + .114*float64(c[5]) - 128
			l01 = .299*float64(c[6]) + .587*float64(c[7]) + .114*float64(c[8]) - 128
			l11 = .299*float64(c[9]) + .587*float64(c[10]) + .114*float64(c[11]) - 128
		}
	}
	experimentalV4Build87BenchmarkFloat = l00 + l10 + l01 + l11
}

func BenchmarkExperimentalV4Build87PreparedHorizontalBilinearAngle(b *testing.B) {
	f, prepared := experimentalV4Build87PrepareAngleSamples(b)
	b.ReportAllocs()
	b.ResetTimer()
	var top, bottom float64
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			s := prepared[base+j]
			top = s.l00*(1-s.fx) + s.l10*s.fx
			bottom = s.l01*(1-s.fx) + s.l11*s.fx
		}
	}
	experimentalV4Build87BenchmarkFloat = top + bottom
}

func BenchmarkExperimentalV4Build87PreparedVerticalBilinearAngle(b *testing.B) {
	f, prepared := experimentalV4Build87PrepareAngleSamples(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			s := prepared[base+j]
			value = s.top*(1-s.fy) + s.bottom*s.fy
		}
	}
	experimentalV4Build87BenchmarkFloat = value
}

func BenchmarkExperimentalV4Build87PreparedFetchLumaBilinearAngle(b *testing.B) {
	f, prepared := experimentalV4Build87PrepareAngleSamples(b)
	rgb := f.plane.rgb
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	for i := 0; i < b.N; i++ {
		base := (i % len(f.origins)) * blockSize * blockSize
		for j := 0; j < blockSize*blockSize; j++ {
			s := prepared[base+j]
			p00 := rgb[s.index00 : s.index00+3]
			p10 := rgb[s.index10 : s.index10+3]
			p01 := rgb[s.index01 : s.index01+3]
			p11 := rgb[s.index11 : s.index11+3]
			l00 := .299*float64(p00[0]) + .587*float64(p00[1]) + .114*float64(p00[2]) - 128
			l10 := .299*float64(p10[0]) + .587*float64(p10[1]) + .114*float64(p10[2]) - 128
			l01 := .299*float64(p01[0]) + .587*float64(p01[1]) + .114*float64(p01[2]) - 128
			l11 := .299*float64(p11[0]) + .587*float64(p11[1]) + .114*float64(p11[2]) - 128
			top := l00*(1-s.fx) + l10*s.fx
			bottom := l01*(1-s.fx) + l11*s.fx
			value = top*(1-s.fy) + bottom*s.fy
		}
	}
	experimentalV4Build87BenchmarkFloat = value
}
