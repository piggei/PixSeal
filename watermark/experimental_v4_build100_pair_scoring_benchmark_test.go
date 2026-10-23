package watermark

import "testing"

var experimentalV4Build100BenchSink int

func BenchmarkExperimentalV4Build100Freeze(b *testing.B) {
	carrier, _, _, w, h, _, boundary := experimentalV4Build100Fixture(b)
	for _, variant := range []string{"Build98", "Build100"} {
		b.Run(variant, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if variant == "Build98" {
					bank, _, evals, _ := experimentalV4PhoneBuild75FreezeParallel(carrier, boundary, w, h, 128)
					experimentalV4Build100BenchSink = len(bank) + evals
				} else {
					bank, _, evals, _, _ := experimentalV4PhoneBuild100FreezeParallel(carrier, boundary, w, h, 128)
					experimentalV4Build100BenchSink = len(bank) + evals
				}
			}
		})
	}
}
