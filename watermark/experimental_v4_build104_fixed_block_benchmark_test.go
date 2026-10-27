package watermark

import (
	"os"
	"runtime"
	"sync"
	"testing"
)

// Same bounded original-index pool for both conditions. An operation is exactly
// eight public tasks, with eight workers on the qualified benchmark host.
func experimentalV4Build104Pool(plane *pixelPlane, pilot experimentalV4PilotCandidate, w, h int, anchor [4]ImagePoint, root experimentalV4PhoneHypothesis, candidate bool) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild84FetchTelemetry) {
	const tasks = 8
	workers := runtime.GOMAXPROCS(0)
	if workers > tasks {
		workers = tasks
	}
	jobs := make(chan int)
	results := make([]experimentalV4PhoneBuild84Gen4Result, tasks)
	var group sync.WaitGroup
	group.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer group.Done()
			for index := range jobs {
				var bank []experimentalV4PhoneHypothesis
				var evals int
				var fetch experimentalV4PhoneBuild84FetchTelemetry
				if candidate {
					bank, evals, fetch = experimentalV4Build104Generation4Fetch(plane, pilot, w, h, anchor, root)
				} else {
					bank, evals, fetch = experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
				}
				results[index] = experimentalV4PhoneBuild84Gen4Result{bank: bank, evals: evals, fetch: fetch}
			}
		}()
	}
	for i := 0; i < tasks; i++ {
		jobs <- i
	}
	close(jobs)
	group.Wait()
	bank := make([]experimentalV4PhoneHypothesis, 0, 64)
	evals := 0
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	for _, result := range results {
		bank = append(bank, result.bank...)
		evals += result.evals
		experimentalV4PhoneBuild84MergeFetch(&fetch, result.fetch)
	}
	return bank, evals, fetch
}

var experimentalV4Build104BenchFloat float64
var experimentalV4Build104BenchInt int
var experimentalV4Build104BenchBool bool

// One benchmark selected per fresh child; validation/setup is excluded from its
// measured window. Direct calls preserve each reader's actual compiler shape.
func BenchmarkExperimentalV4Build104Slot(b *testing.B) {
	metric, variant := os.Getenv("PIXSEAL_BUILD104_METRIC"), os.Getenv("PIXSEAL_BUILD104_VARIANT")
	if metric == "" {
		b.Skip("fresh-process benchmark is opt-in")
	}
	if variant != "Q" && variant != "C" {
		b.Fatal("invalid benchmark variant")
	}
	if runtime.GOMAXPROCS(0) != 8 || os.Getenv("GODEBUG") != "asyncpreemptoff=0" || os.Getenv("GOGC") != "100" || os.Getenv("GOMEMLIMIT") != "off" {
		b.Fatal("benchmark runtime environment differs")
	}
	b.ReportAllocs()
	if metric == "ReaderFront" || metric == "ReaderMild" || metric == "ReaderAngle" {
		f := experimentalV4Build83Fixture(b)
		h := f.front
		if metric == "ReaderMild" {
			h = f.mild
		}
		if metric == "ReaderAngle" {
			h = f.angle
		}
		for _, origin := range f.origins {
			q, qok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1], blockSize)
			c, cok := experimentalV4Build104ReadProjectiveBlockValue(f.plane, h, origin[0], origin[1])
			if !qok || !cok || q != c {
				b.Fatal("invalid reader fixture")
			}
		}
		var v float64
		var ok bool
		b.ResetTimer()
		if variant == "Q" {
			for i := 0; i < b.N; i++ {
				o := f.origins[i%len(f.origins)]
				v, ok = experimentalV4PhoneBuild84ReadProjectiveBlockValue(f.plane, h, o[0], o[1], blockSize)
			}
		} else {
			for i := 0; i < b.N; i++ {
				o := f.origins[i%len(f.origins)]
				v, ok = experimentalV4Build104ReadProjectiveBlockValue(f.plane, h, o[0], o[1])
			}
		}
		experimentalV4Build104BenchFloat, experimentalV4Build104BenchBool = v, ok
		return
	}
	if metric != "Generation4" && metric != "Pool8" {
		b.Fatal("unknown metric")
	}
	plane, pilot, w, h, anchor, root := experimentalV4Build104StageFixture(b)
	q, qe, qf := experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
	c, ce, cf := experimentalV4Build104Generation4Fetch(plane, pilot, w, h, anchor, root)
	if experimentalV4Build102Hash(q) != experimentalV4Build102Hash(c) || qe != ce || qf != cf || qf.FoldScores != 128 || qe != 256 || len(q) != 6 {
		b.Fatal("invalid complete-stage fixture")
	}
	b.ResetTimer()
	sum := 0
	for i := 0; i < b.N; i++ {
		var bank []experimentalV4PhoneHypothesis
		var evals int
		var fetch experimentalV4PhoneBuild84FetchTelemetry
		if metric == "Pool8" {
			bank, evals, fetch = experimentalV4Build104Pool(plane, pilot, w, h, anchor, root, variant == "C")
		} else if variant == "C" {
			bank, evals, fetch = experimentalV4Build104Generation4Fetch(plane, pilot, w, h, anchor, root)
		} else {
			bank, evals, fetch = experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
		}
		sum += len(bank) + evals + fetch.BlockReads
	}
	experimentalV4Build104BenchInt = sum
}
