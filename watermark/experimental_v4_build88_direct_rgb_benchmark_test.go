package watermark

import (
	"math"
	"testing"
)

// Build88 benchmark work is public and deterministic. It compares the current
// qualified Build84 reader with the exact Build88 direct-RGB/BCE candidate on
// the same front/mild/angle fixture. The prepared microbenchmarks isolate only
// the four-pixel access shape: Build84 three-byte slices versus Build88 one
// dominating index+2 proof followed by direct scalar byte loads.

var experimentalV4Build88BenchmarkFloat float64
var experimentalV4Build88BenchmarkInt int
var experimentalV4Build88BenchmarkBool bool

func TestExperimentalV4Build88PublicBenchmarkFixture(t *testing.T) {
	f := experimentalV4Build83Fixture(t)
	for name, h := range map[string]homography{"front": f.front, "mild": f.mild, "angle": f.angle} {
		for _, origin := range f.origins {
			want, wok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			got, gok := experimentalV4PhoneBuild88ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
				t.Fatalf("%s origin=%v exact mismatch ok=%t/%t bits=%x/%x", name, origin, gok, wok, math.Float64bits(got), math.Float64bits(want))
			}
		}
	}
}

func experimentalV4Build88BenchmarkReader(b *testing.B, candidate bool, h homography) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		if candidate {
			value, ok = experimentalV4PhoneBuild88ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
		} else {
			value, ok = experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
		}
	}
	experimentalV4Build88BenchmarkFloat = value
	experimentalV4Build88BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build88Build84BlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build88BenchmarkReader(b, false, f.front)
}
func BenchmarkExperimentalV4Build88Build84BlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build88BenchmarkReader(b, false, f.mild)
}
func BenchmarkExperimentalV4Build88Build84BlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build88BenchmarkReader(b, false, f.angle)
}
func BenchmarkExperimentalV4Build88DirectRGBBlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build88BenchmarkReader(b, true, f.front)
}
func BenchmarkExperimentalV4Build88DirectRGBBlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build88BenchmarkReader(b, true, f.mild)
}
func BenchmarkExperimentalV4Build88DirectRGBBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build88BenchmarkReader(b, true, f.angle)
}

func BenchmarkExperimentalV4Build88PreparedSliceFetchAngle(b *testing.B) {
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
	experimentalV4Build88BenchmarkInt = sink
}

func BenchmarkExperimentalV4Build88PreparedDirectFetchAngle(b *testing.B) {
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
			_ = rgb[s.index00+2]
			local += int(rgb[s.index00]) + int(rgb[s.index00+1]) + int(rgb[s.index00+2])
			_ = rgb[s.index10+2]
			local += int(rgb[s.index10]) + int(rgb[s.index10+1]) + int(rgb[s.index10+2])
			_ = rgb[s.index01+2]
			local += int(rgb[s.index01]) + int(rgb[s.index01+1]) + int(rgb[s.index01+2])
			_ = rgb[s.index11+2]
			local += int(rgb[s.index11]) + int(rgb[s.index11+1]) + int(rgb[s.index11+2])
		}
		sink = local
	}
	experimentalV4Build88BenchmarkInt = sink
}
