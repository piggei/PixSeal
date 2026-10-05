package watermark

import (
	"math"
	"testing"
)

// Build84 benchmark work is public and deterministic. It compares the
// authoritative historical Build76 reader, the closed Build82 manual-inline
// comparator, and the Build84 RGB-fetch/BCE candidate on the Build83 fixture.
// No private acquisition, key, payload, ECC result, or HMAC result selects work.

var experimentalV4Build84BenchmarkFloat float64
var experimentalV4Build84BenchmarkBool bool

func TestExperimentalV4Build84PublicBenchmarkFixture(t *testing.T) {
	f := experimentalV4Build83Fixture(t)
	for name, h := range map[string]homography{"front": f.front, "mild": f.mild, "angle": f.angle} {
		for _, origin := range f.origins {
			want, wok := readProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			got, gok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
				t.Fatalf("%s origin=%v exact mismatch ok=%t/%t bits=%x/%x", name, origin, gok, wok, math.Float64bits(got), math.Float64bits(want))
			}
		}
	}
}

func experimentalV4Build84BenchmarkReader(b *testing.B, which string, h homography) {
	f := experimentalV4Build83Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	var value float64
	var ok bool
	for i := 0; i < b.N; i++ {
		origin := f.origins[i%len(f.origins)]
		switch which {
		case "historical":
			value, ok = readProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
		case "build82":
			value, ok = experimentalV4PhoneBuild82ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
		case "build84":
			value, ok = experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
		default:
			b.Fatalf("unknown reader %q", which)
		}
	}
	experimentalV4Build84BenchmarkFloat = value
	experimentalV4Build84BenchmarkBool = ok
}

func BenchmarkExperimentalV4Build84HistoricalBlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build84BenchmarkReader(b, "historical", f.front)
}
func BenchmarkExperimentalV4Build84HistoricalBlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build84BenchmarkReader(b, "historical", f.mild)
}
func BenchmarkExperimentalV4Build84HistoricalBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build84BenchmarkReader(b, "historical", f.angle)
}
func BenchmarkExperimentalV4Build84Build82BlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build84BenchmarkReader(b, "build82", f.angle)
}
func BenchmarkExperimentalV4Build84RGBFetchBlockReaderFront(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build84BenchmarkReader(b, "build84", f.front)
}
func BenchmarkExperimentalV4Build84RGBFetchBlockReaderMild(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build84BenchmarkReader(b, "build84", f.mild)
}
func BenchmarkExperimentalV4Build84RGBFetchBlockReaderAngle(b *testing.B) {
	f := experimentalV4Build83Fixture(b)
	experimentalV4Build84BenchmarkReader(b, "build84", f.angle)
}
