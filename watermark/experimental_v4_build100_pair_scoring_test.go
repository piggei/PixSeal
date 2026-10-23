package watermark

import (
	"bytes"
	"image"
	"math"
	"reflect"
	"testing"
	"time"
)

func experimentalV4Build100Fixture(tb testing.TB) (image.Image, *pixelPlane, experimentalV4PilotCandidate, int, int, [4]ImagePoint, PrintBoundaryEstimate) {
	tb.Helper()
	pilot := experimentalV4Prototype2Candidate()
	carrier, err := experimentalV4RenderSyntheticCarrier(testImage(592, 512), pilot, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		tb.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	anchor := [4]ImagePoint{{X: 0, Y: 0}, {X: float64(w - 1), Y: 0}, {X: 0, Y: float64(h - 1)}, {X: float64(w - 1), Y: float64(h - 1)}}
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: anchor[0], TopRight: anchor[1], BottomLeft: anchor[2], BottomRight: anchor[3]}
	return carrier, newPixelPlane(carrier), pilot, w, h, anchor, boundary
}

func experimentalV4Build100SerialPairs(plane *pixelPlane, pilot experimentalV4PilotCandidate, w, h int, anchor [4]ImagePoint, base, structural [4]experimentalV4PhoneBuild43Line, ok [4]bool) ([]experimentalV4PhoneBuild100RankedPair, int) {
	ranked := make([]experimentalV4PhoneBuild100RankedPair, 0, len(experimentalV4PhoneBuild43Pairs))
	evals := 0
	for _, pair := range experimentalV4PhoneBuild43Pairs {
		seed := experimentalV4PhoneBuild43PairSeed(base, structural, ok, pair)
		score, n := experimentalV4PhoneBuild43PairRobustScore(plane, pilot, w, h, anchor, seed, pair)
		evals += n
		if !math.IsInf(score, -1) {
			ranked = append(ranked, experimentalV4PhoneBuild100RankedPair{pair: pair, score: score, seed: seed})
		}
	}
	return ranked, evals
}

func TestExperimentalV4Build100PairScoresExactAcrossWorkerCounts(t *testing.T) {
	carrier, plane, pilot, w, h, anchor, _ := experimentalV4Build100Fixture(t)
	base := experimentalV4PhoneBuild43BaseLines(anchor)
	structural, ok := experimentalV4PhoneBuild43StructuralLines(carrier, anchor)
	want, evals := experimentalV4Build100SerialPairs(plane, pilot, w, h, anchor, base, structural, ok)
	if len(want) == 0 || evals == 0 {
		t.Fatal("empty scoring fixture")
	}
	for _, workers := range []int{1, 2, 6, 16} {
		got, tele := experimentalV4PhoneBuild100ScorePairs(plane, pilot, w, h, anchor, base, structural, ok, workers)
		if !reflect.DeepEqual(got, want) || tele.Evaluations != evals || tele.Tasks != 6 {
			t.Fatalf("workers=%d changed scores/order/counts", workers)
		}
		expected := workers
		if expected > 6 {
			expected = 6
		}
		if tele.Workers != expected {
			t.Fatalf("workers=%d want %d", tele.Workers, expected)
		}
		for i := range got {
			if math.Float64bits(got[i].score) != math.Float64bits(want[i].score) {
				t.Fatalf("score bits %d", i)
			}
		}
	}
}

func TestExperimentalV4Build100ParallelPairScoreRaceFixture(t *testing.T) {
	carrier, plane, pilot, w, h, anchor, _ := experimentalV4Build100Fixture(t)
	base := experimentalV4PhoneBuild43BaseLines(anchor)
	structural, ok := experimentalV4PhoneBuild43StructuralLines(carrier, anchor)
	want, evals := experimentalV4Build100SerialPairs(plane, pilot, w, h, anchor, base, structural, ok)
	got, tele := experimentalV4PhoneBuild100ScorePairs(plane, pilot, w, h, anchor, base, structural, ok, 6)
	if !reflect.DeepEqual(got, want) || tele.Evaluations != evals {
		t.Fatal("concurrent pair scoring changed results")
	}
}

func TestExperimentalV4Build100FreezeMatchesBuild75Exactly(t *testing.T) {
	carrier, _, _, w, h, _, boundary := experimentalV4Build100Fixture(t)
	for _, cap := range []int{1, 128} {
		want, ws, we, wt := experimentalV4PhoneBuild75FreezeParallel(carrier, boundary, w, h, cap)
		got, gs, ge, gt, pt := experimentalV4PhoneBuild100FreezeParallel(carrier, boundary, w, h, cap)
		if !reflect.DeepEqual(got, want) || !reflect.DeepEqual(gs, ws) || ge != we {
			t.Fatalf("freeze cap=%d changed bank/order/scores/evals", cap)
		}
		if gt.PairScoreTasks != wt.PairScoreTasks || gt.PairScoreEvaluations != wt.PairScoreEvaluations || gt.BasinTasks != wt.BasinTasks || gt.BasinEvaluations != wt.BasinEvaluations || pt.Tasks != 6 {
			t.Fatalf("freeze semantic accounting changed cap=%d", cap)
		}
	}
	bank, scores, evals, fp, pp := experimentalV4PhoneBuild100FreezeParallel(nil, PrintBoundaryEstimate{}, w, h, 128)
	if len(bank) != 0 || len(scores) != 0 || evals != 0 || fp.PairScoreTasks != 0 || pp.Tasks != 0 {
		t.Fatal("invalid boundary did work")
	}
}

func experimentalV4Build100ClearDurations(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Struct {
			experimentalV4Build100ClearDurations(f)
		} else if f.Type() == reflect.TypeOf(time.Duration(0)) {
			f.SetInt(0)
		}
	}
}

func TestExperimentalV4Build100BlindBankMatchesBuild98Exactly(t *testing.T) {
	carrier, _, _, w, h, _, boundary := experimentalV4Build100Fixture(t)
	want, we, ws, ww, wg, wf := experimentalV4PhoneBuild98BlindBank(carrier, boundary, w, h)
	got, ge, gs, gw, gg, gf := experimentalV4PhoneBuild100BlindBank(carrier, boundary, w, h)
	if !reflect.DeepEqual(got, want) || ge != we || gs != ws || gw != ww {
		t.Fatal("complete bank or traversal/evals changed")
	}
	pg := gg.experimentalV4PhoneBuild98GeometryTelemetry
	experimentalV4Build100ClearDurations(reflect.ValueOf(&pg).Elem())
	experimentalV4Build100ClearDurations(reflect.ValueOf(&wg).Elem())
	experimentalV4Build100ClearDurations(reflect.ValueOf(&gf).Elem())
	experimentalV4Build100ClearDurations(reflect.ValueOf(&wf).Elem())
	if !reflect.DeepEqual(pg, wg) || !reflect.DeepEqual(gf, wf) {
		t.Fatal("historical semantic telemetry changed")
	}
}

func TestExperimentalV4Build100RecoveryMatchesBuild98Exactly(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: f.hyp.quad[0], TopRight: f.hyp.quad[1], BottomLeft: f.hyp.quad[2], BottomRight: f.hyp.quad[3]}
	wp, wi, wt, we := experimentalV4PhoneBuild98Recover(f.image, boundary, f.key, f.width, f.height)
	gp, gi, gt, ge := experimentalV4PhoneBuild100Recover(f.image, boundary, f.key, f.width, f.height)
	if (we == nil) != (ge == nil) || !bytes.Equal(wp, gp) || !reflect.DeepEqual(wi, gi) {
		t.Fatal("recovery outcome changed")
	}
	baseline := gt.experimentalV4PhoneBuild98RecoveryTelemetry
	experimentalV4Build100ClearDurations(reflect.ValueOf(&wt).Elem())
	experimentalV4Build100ClearDurations(reflect.ValueOf(&baseline).Elem())
	if !reflect.DeepEqual(baseline, wt) {
		t.Fatal("qualified recovery semantic counters changed")
	}
}
