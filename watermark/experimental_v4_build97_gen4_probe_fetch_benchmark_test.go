package watermark

import "testing"

var experimentalV4Build97BenchmarkInt int

func BenchmarkExperimentalV4Build97Build53SingleProbe(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e, n := experimentalV4PhoneBuild53SingleProbe(plane, pilot, w, h, anchor, root, 0)
		experimentalV4Build97BenchmarkInt = e + n
	}
}
func BenchmarkExperimentalV4Build97Build84SingleProbe(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e, n, f := experimentalV4PhoneBuild97SingleProbeFetch(plane, pilot, w, h, anchor, root, 0)
		experimentalV4Build97BenchmarkInt = e + n + f.BlockSuccess
	}
}
func BenchmarkExperimentalV4Build97Build53PairStencil(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e, n := experimentalV4PhoneBuild53PairStencil(plane, pilot, w, h, anchor, root, 0)
		experimentalV4Build97BenchmarkInt = len(s) + e + n
	}
}
func BenchmarkExperimentalV4Build97Build84PairStencil(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e, n, f := experimentalV4PhoneBuild97PairStencilFetch(plane, pilot, w, h, anchor, root, 0)
		experimentalV4Build97BenchmarkInt = len(s) + e + n + f.BlockSuccess
	}
}
func BenchmarkExperimentalV4Build97Build84Generation4(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e, f := experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
		experimentalV4Build97BenchmarkInt = len(s) + e + f.BlockSuccess
	}
}
func BenchmarkExperimentalV4Build97ProbeFetchGeneration4(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e, f := experimentalV4PhoneBuild97Generation4Fetch(plane, pilot, w, h, anchor, root)
		experimentalV4Build97BenchmarkInt = len(s) + e + f.Continuation.BlockSuccess + f.Probe.BlockSuccess
	}
}
