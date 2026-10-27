package watermark

import (
	"crypto/sha256"
	"image"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func experimentalV4Build104AssertReader(t *testing.T, p *pixelPlane, h homography, x, y int) {
	t.Helper()
	want, wok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(p, h, x, y, blockSize)
	got, gok := experimentalV4Build104ReadProjectiveBlockValue(p, h, x, y)
	if gok != wok || math.Float64bits(got) != math.Float64bits(want) {
		t.Fatalf("reader mismatch origin=%d,%d ok=%t/%t bits=%x/%x", x, y, gok, wok, math.Float64bits(got), math.Float64bits(want))
	}
}
func TestExperimentalV4Build104ReaderBitsAndBoundaries(t *testing.T) {
	f := experimentalV4Build83Fixture(t)
	for _, h := range []homography{f.front, f.mild, f.angle} {
		for _, o := range append(f.origins, [][2]int{{0, 0}, {1624, 1624}, {-1, 0}, {0, -1}, {1625, 1625}}...) {
			experimentalV4Build104AssertReader(t, f.plane, h, o[0], o[1])
		}
	}
	random := rand.New(rand.NewSource(104))
	for i := 0; i < 1500; i++ {
		h := homography{h: [9]float64{.4 + random.Float64(), random.Float64() * .2, random.Float64()*40 - 20, random.Float64() * .2, .4 + random.Float64(), random.Float64()*40 - 20, random.Float64()*.0004 - .0002, random.Float64()*.0004 - .0002, 1}}
		experimentalV4Build104AssertReader(t, f.plane, h, random.Intn(1700)-30, random.Intn(1700)-30)
	}
	for _, size := range []int{1, 2, 7, 8, 9, 17} {
		p := &pixelPlane{bounds: image.Rect(17, 23, 17+size, 23+size), rgb: make([]byte, size*size*3)}
		for i := range p.rgb {
			p.rgb[i] = byte(i*37 + 11)
		}
		for _, d := range []float64{0, math.Copysign(0, -1), math.Nextafter(1e-12, 0), 1e-12, -1e-12, 1, math.Inf(1)} {
			h := homography{h: [9]float64{0, 0, 0, 0, 0, 0, 0, 0, d}}
			// Infinite denominator maps to finite zero; numerator stays finite.
			experimentalV4Build104AssertReader(t, p, h, 0, 0)
		}
		for _, o := range [][2]int{{0, 0}, {size - blockSize, size - blockSize}, {size - 1, size - 1}, {-1, -1}} {
			experimentalV4Build104AssertReader(t, p, homography{h: [9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}}, o[0], o[1])
		}
	}
}
func TestExperimentalV4Build104FoldScoreExact(t *testing.T) {
	plane, pilot, w, h, _, root := experimentalV4Build97Fixture(t)
	for fold := 0; fold < 4; fold++ {
		for _, proposal := range []bool{false, true} {
			for _, sparse := range []bool{false, true} {
				want, wv, wf := experimentalV4PhoneBuild84FoldScoreFetch(plane, pilot, w, h, root.h, fold, proposal, sparse)
				got, gv, gf := experimentalV4Build104FoldScoreFetch(plane, pilot, w, h, root.h, fold, proposal, sparse)
				if math.Float64bits(got) != math.Float64bits(want) || gv != wv || gf != wf {
					t.Fatalf("fold=%d proposal=%t sparse=%t mismatch", fold, proposal, sparse)
				}
			}
		}
	}
	for _, p := range []*pixelPlane{nil, plane} {
		want, wv, wf := experimentalV4PhoneBuild84FoldScoreFetch(p, pilot, 1, 1, root.h, 0, true, false)
		got, gv, gf := experimentalV4Build104FoldScoreFetch(p, pilot, 1, 1, root.h, 0, true, false)
		if math.Float64bits(got) != math.Float64bits(want) || gv != wv || gf != wf {
			t.Fatal("empty/nil fold changed")
		}
	}
}
func TestExperimentalV4Build104ContinuationAndGenerationExact(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build104StageFixture(t)
	for _, proposal := range []float64{root.proposal, root.proposal - .1} {
		input := root
		input.proposal = proposal
		want, we, wf := experimentalV4PhoneBuild84ContinueFetch(plane, pilot, w, h, anchor, input, 0)
		got, ge, gf := experimentalV4Build104ContinueFetch(plane, pilot, w, h, anchor, input, 0)
		if experimentalV4Build102Hash(got) != experimentalV4Build102Hash(want) || ge != we || gf != wf {
			t.Fatal("continuation state/order/evaluation/fetch changed")
		}
	}
	want, we, wf := experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
	got, ge, gf := experimentalV4Build104Generation4Fetch(plane, pilot, w, h, anchor, root)
	if experimentalV4Build102Hash(got) != experimentalV4Build102Hash(want) || ge != we || gf != wf {
		t.Fatal("Gen4 bank/evaluations/telemetry changed")
	}
	if we != 256 || len(want) != 6 || wf.FoldScores != 128 || wf.BlockReads != 16384 {
		t.Fatal("fixture never exercised changed continuation reader")
	}
	t.Logf("public Gen4 fixture: evals=%d bank=%d continuation_folds=%d block_reads=%d", we, len(want), wf.FoldScores, wf.BlockReads)
}
func TestExperimentalV4Build104OrderedPoolRace(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build104StageFixture(t)
	before := experimentalV4Build102Hash(struct {
		RGB  [32]byte
		Root experimentalV4PhoneHypothesis
	}{sha256.Sum256(plane.rgb), root})
	want, we, wf := experimentalV4Build104Pool(plane, pilot, w, h, anchor, root, false)
	got, ge, gf := experimentalV4Build104Pool(plane, pilot, w, h, anchor, root, true)
	if experimentalV4Build102Hash(got) != experimentalV4Build102Hash(want) || ge != we || !reflect.DeepEqual(gf, wf) {
		t.Fatal("ordered pool bank/counters changed")
	}
	after := experimentalV4Build102Hash(struct {
		RGB  [32]byte
		Root experimentalV4PhoneHypothesis
	}{sha256.Sum256(plane.rgb), root})
	if before != after {
		t.Fatal("immutable pool input mutated")
	}
}
