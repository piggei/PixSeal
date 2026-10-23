package watermark

import "testing"

var experimentalV4Build98BenchmarkInt int

func BenchmarkExperimentalV4Build98Build55SiblingStencil(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, w, h, anchor, root, 0)
		experimentalV4Build98BenchmarkInt = len(s) + e
	}
}

func BenchmarkExperimentalV4Build98Build84SiblingStencil(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e, f := experimentalV4PhoneBuild98SiblingStencilFetch(plane, pilot, w, h, anchor, root, 0)
		experimentalV4Build98BenchmarkInt = len(s) + e + f.BlockSuccess
	}
}

func BenchmarkExperimentalV4Build98Build73Generation3(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e := experimentalV4PhoneBuild73Generation3(plane, pilot, w, h, anchor, root)
		experimentalV4Build98BenchmarkInt = len(s) + e
	}
}

func BenchmarkExperimentalV4Build98FetchGeneration3(b *testing.B) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e, f := experimentalV4PhoneBuild98Generation3Fetch(plane, pilot, w, h, anchor, root)
		t := experimentalV4PhoneBuild98TotalFetch(f)
		experimentalV4Build98BenchmarkInt = len(s) + e + t.BlockSuccess
	}
}
