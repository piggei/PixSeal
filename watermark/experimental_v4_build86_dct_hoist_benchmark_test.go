package watermark

import (
	"math"
	"testing"
)

// Build86 benchmark work is public and deterministic. It compares the current
// qualified Build84 reader with the exact Build86 DCT-table-hoist candidate on
// the same fixture and separately measures only the DCT accumulation shape.

var experimentalV4Build86BenchmarkFloat float64
var experimentalV4Build86BenchmarkBool bool

func TestExperimentalV4Build86PublicBenchmarkFixture(t *testing.T) {
	f := experimentalV4Build83Fixture(t)
	for name, h := range map[string]homography{"front": f.front, "mild": f.mild, "angle": f.angle} {
		for _, origin := range f.origins {
			want, wok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			got, gok := experimentalV4PhoneBuild86ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
				t.Fatalf("%s origin=%v exact mismatch ok=%t/%t bits=%x/%x", name, origin, gok, wok, math.Float64bits(got), math.Float64bits(want))
			}
		}
	}
}

func experimentalV4Build86BenchmarkReader(b *testing.B, candidate bool, h homography) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		if candidate {
			value, ok = experimentalV4PhoneBuild86ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
		} else {
			value, ok = experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
		}
	}
	experimentalV4Build86BenchmarkFloat = value
	experimentalV4Build86BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build86Build84BlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build86BenchmarkReader(b, false, f.front)
}
func BenchmarkExperimentalV4Build86Build84BlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build86BenchmarkReader(b, false, f.mild)
}
func BenchmarkExperimentalV4Build86Build84BlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build86BenchmarkReader(b, false, f.angle)
}
func BenchmarkExperimentalV4Build86DCTHoistBlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build86BenchmarkReader(b, true, f.front)
}
func BenchmarkExperimentalV4Build86DCTHoistBlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build86BenchmarkReader(b, true, f.mild)
}
func BenchmarkExperimentalV4Build86DCTHoistBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build86BenchmarkReader(b, true, f.angle)
}

func BenchmarkExperimentalV4Build86Build84DCTAccumulation(b *testing.B) {
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
	experimentalV4Build86BenchmarkFloat = out
}

func BenchmarkExperimentalV4Build86HoistedDCTAccumulation(b *testing.B) {
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
			t2y := table2[y]
			t3y := table3[y]
			for x := 0; x < blockSize; x++ {
				t2x := table2[x]
				t3x := table3[x]
				l := f.dctValues[idx]
				c23 += l * t3x * t2y
				c32 += l * t2x * t3y
				idx++
			}
		}
		out = math.Abs(c23) - math.Abs(c32)
	}
	experimentalV4Build86BenchmarkFloat = out
}
