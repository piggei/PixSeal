package watermark

import (
	"image"
	"math"
	"testing"
)

func build80SyntheticFixture(t *testing.T) (experimentalV4PilotCandidate, image.Image, *pixelPlane, int, int, [4]ImagePoint, homography) {
	t.Helper()
	candidate := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, candidate, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		t.Fatal(err)
	}
	w, h := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	plane := newPixelPlane(carrier)
	anchor := [4]ImagePoint{{X: 0, Y: 0}, {X: float64(w - 1), Y: 0}, {X: 0, Y: float64(h - 1)}, {X: float64(w - 1), Y: float64(h - 1)}}
	hom, ok := experimentalV4PhoneQuadHomography(w, h, anchor)
	if !ok {
		t.Fatal("anchor homography")
	}
	return candidate, carrier, plane, w, h, anchor, hom
}

func TestExperimentalV4Build80ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneBuild80MaxCachePixels != 196 {
		t.Fatalf("cache limit=%d want 196", experimentalV4PhoneBuild80MaxCachePixels)
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build80CachedBlockMatchesOriginalBitExactly(t *testing.T) {
	_, _, plane, w, h, anchor, _ := build80SyntheticFixture(t)
	quads := [][4]ImagePoint{
		anchor,
		{{X: 3, Y: 2}, {X: float64(w - 5), Y: 5}, {X: 6, Y: float64(h - 4)}, {X: float64(w - 7), Y: float64(h - 6)}},
		{{X: 14, Y: 5}, {X: float64(w - 18), Y: 11}, {X: 4, Y: float64(h - 15)}, {X: float64(w - 9), Y: float64(h - 7)}},
	}
	for qi, q := range quads {
		hom, ok := experimentalV4PhoneQuadHomography(w, h, q)
		if !ok {
			t.Fatalf("homography %d", qi)
		}
		var scratch experimentalV4PhoneBuild80BlockScratch
		for oy := 0; oy+blockSize <= h; oy += 37 {
			for ox := 0; ox+blockSize <= w; ox += 41 {
				want, wantOK := readProjectiveBlockValue(plane, hom, ox, oy, blockSize)
				got, gotOK, _, _ := experimentalV4PhoneBuild80ReadProjectiveBlockCached(plane, hom, ox, oy, blockSize, &scratch)
				if gotOK != wantOK || (gotOK && math.Float64bits(got) != math.Float64bits(want)) {
					t.Fatalf("q=%d block=(%d,%d) got=(%.17g,%t) want=(%.17g,%t)", qi, ox, oy, got, gotOK, want, wantOK)
				}
			}
		}
	}
}

func TestExperimentalV4Build80CachedBlockFallsBackExactly(t *testing.T) {
	_, _, plane, _, _, _, _ := build80SyntheticFixture(t)
	h := homography{h: [9]float64{4, 0, 0, 0, 4, 0, 0, 0, 1}}
	want, wantOK := readProjectiveBlockValue(plane, h, 0, 0, blockSize)
	var scratch experimentalV4PhoneBuild80BlockScratch
	got, gotOK, cached, area := experimentalV4PhoneBuild80ReadProjectiveBlockCached(plane, h, 0, 0, blockSize, &scratch)
	if cached || area <= experimentalV4PhoneBuild80MaxCachePixels {
		t.Fatalf("fallback not exercised cached=%t area=%d", cached, area)
	}
	if gotOK != wantOK || (gotOK && math.Float64bits(got) != math.Float64bits(want)) {
		t.Fatalf("fallback differs got=(%.17g,%t) want=(%.17g,%t)", got, gotOK, want, wantOK)
	}
}

func TestExperimentalV4Build80FoldScoreMatchesBuild41BitExactly(t *testing.T) {
	candidate, _, plane, w, h, _, hom := build80SyntheticFixture(t)
	for _, tc := range []struct{ proposal, sparse bool }{{true, false}, {false, false}, {true, true}} {
		want, wantVisible := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		got, gotVisible, cache := experimentalV4PhoneBuild80FoldScoreCached(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		if gotVisible != wantVisible || math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("proposal=%t sparse=%t score/visible got=(%.17g,%d) want=(%.17g,%d)", tc.proposal, tc.sparse, got, gotVisible, want, wantVisible)
		}
		if cache.Hits+cache.Fallbacks+cache.Failed != cache.BlockReads {
			t.Fatalf("cache accounting incomplete: %+v", cache)
		}
		if wantVisible > 0 && (cache.BlockReads == 0 || cache.Hits == 0) {
			t.Fatalf("cache was not exercised for visible FoldScore: %+v", cache)
		}
	}
}

func TestExperimentalV4Build80ContinueMatchesBuild55Exactly(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build80SyntheticFixture(t)
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, wantEvals := experimentalV4PhoneBuild55Continue(plane, candidate, w, h, anchor, start, 0)
	got, gotEvals, cache := experimentalV4PhoneBuild80ContinueCached(plane, candidate, w, h, anchor, start, 0)
	if gotEvals != wantEvals || len(got) != len(want) {
		t.Fatalf("evals/states got=%d/%d want=%d/%d", gotEvals, len(got), wantEvals, len(want))
	}
	for i := range got {
		if got[i].index != want[i].index || got[i].pass != want[i].pass || got[i].dim != want[i].dim || got[i].delta != want[i].delta || got[i].hyp.quad != want[i].hyp.quad || got[i].hyp.h != want[i].hyp.h || math.Float64bits(got[i].hyp.proposal) != math.Float64bits(want[i].hyp.proposal) {
			t.Fatalf("continuation differs at %d", i)
		}
	}
	if wantEvals > 0 && cache.BlockReads == 0 {
		t.Fatalf("missing cache work: %+v", cache)
	}
}

func TestExperimentalV4Build80BlindBankMatchesBuild76Exactly(t *testing.T) {
	_, carrier, _, w, h, _, _ := build80SyntheticFixture(t)
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, wantEvals, wantSeeds, wantWorkers, _, _ := experimentalV4PhoneBuild76BlindBank(carrier, boundary, w, h)
	got, gotEvals, gotSeeds, gotWorkers, gp, fp := experimentalV4PhoneBuild80BlindBank(carrier, boundary, w, h)
	if gotEvals != wantEvals || gotSeeds != wantSeeds || gotWorkers != wantWorkers || len(got) != len(want) {
		t.Fatalf("blind bank differs evals=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", gotEvals, wantEvals, gotSeeds, wantSeeds, gotWorkers, wantWorkers, len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || math.Float64bits(got[i].proposal) != math.Float64bits(want[i].proposal) || math.Float64bits(got[i].validation) != math.Float64bits(want[i].validation) {
			t.Fatalf("blind bank differs at %d", i)
		}
	}
	if fp.BasinTasks != 16 || gp.Gen2Tasks == 0 || gp.Gen3Tasks == 0 {
		t.Fatalf("telemetry incomplete gp=%+v fp=%+v", gp, fp)
	}
	if gp.CacheHits+gp.CacheFallbacks+gp.CacheFailed != gp.CacheBlockReads {
		t.Fatalf("cache accounting mismatch: %+v", gp)
	}
}

func TestExperimentalV4Build80ApplyTelemetryPreservesBuild76View(t *testing.T) {
	recovery := experimentalV4PhoneBuild80RecoveryTelemetry{}
	recovery.Attempted = true
	recovery.Build80GeometryProfile = experimentalV4PhoneBuild80GeometryTelemetry{Prefix1Workers: 8, Gen2Workers: 8, Gen3Workers: 8, Gen4Workers: 8, Gen2Tasks: 12, Gen3Tasks: 34, Gen4Tasks: 56, CacheBlockReads: 1000, CacheHits: 900, CacheFallbacks: 80, CacheFailed: 20, CachePixelsPrepared: 76543, CacheMaxArea: 121}
	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild80ApplyTelemetry(&public, recovery)
	if !public.Build80Attempted || !public.Build76Attempted || public.Build76Gen4Tasks != 56 || public.Build80CacheBlockReads != 1000 || public.Build80CacheHits != 900 || public.Build80CacheLimit != 196 {
		t.Fatalf("telemetry mismatch: %+v", public)
	}
}
