package watermark

import (
	"image"
	"math"
	"testing"
)

func build86SyntheticFixture(t *testing.T) (experimentalV4PilotCandidate, image.Image, *pixelPlane, int, int, [4]ImagePoint, homography) {
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

func TestExperimentalV4Build86ProjectiveBlockMatchesBuild84BitExactly(t *testing.T) {
	_, _, plane, w, h, anchor, _ := build86SyntheticFixture(t)
	quads := [][4]ImagePoint{
		anchor,
		{{X: 3, Y: 2}, {X: float64(w - 5), Y: 5}, {X: 6, Y: float64(h - 4)}, {X: float64(w - 7), Y: float64(h - 6)}},
		{{X: 14, Y: 5}, {X: float64(w - 18), Y: 11}, {X: 4, Y: float64(h - 15)}, {X: float64(w - 9), Y: float64(h - 7)}},
		{{X: -6, Y: 9}, {X: float64(w + 4), Y: -3}, {X: 11, Y: float64(h + 7)}, {X: float64(w - 14), Y: float64(h - 12)}},
	}
	origins := [][2]int{{0, 0}, {1, 1}, {w - blockSize, 0}, {0, h - blockSize}, {w - blockSize, h - blockSize}, {-1, 0}, {0, -1}, {w - blockSize + 1, h - blockSize + 1}}
	for qi, q := range quads {
		hom, ok := experimentalV4PhoneQuadHomography(w, h, q)
		if !ok {
			t.Fatalf("hom %d", qi)
		}
		for oy := 0; oy+blockSize <= h; oy += 31 {
			for ox := 0; ox+blockSize <= w; ox += 37 {
				origins = append(origins, [2]int{ox, oy})
			}
		}
		for _, origin := range origins {
			want, wok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(plane, hom, origin[0], origin[1], blockSize)
			got, gok := experimentalV4PhoneBuild86ReadProjectiveBlockValue(plane, hom, origin[0], origin[1], blockSize)
			if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
				t.Fatalf("q=%d block=%d,%d ok=%t/%t bits=%x/%x", qi, origin[0], origin[1], gok, wok, math.Float64bits(got), math.Float64bits(want))
			}
		}
	}
}

func TestExperimentalV4Build86FoldScoreMatchesBuild84BitExactly(t *testing.T) {
	candidate, _, plane, w, h, _, hom := build86SyntheticFixture(t)
	for _, tc := range []struct{ proposal, sparse bool }{{true, false}, {false, false}, {true, true}} {
		want, wv, _ := experimentalV4PhoneBuild84FoldScoreFetch(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		got, gv, fetch := experimentalV4PhoneBuild86FoldScoreFetch(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		if gv != wv || math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("fold proposal=%t sparse=%t visible=%d/%d bits=%x/%x", tc.proposal, tc.sparse, gv, wv, math.Float64bits(got), math.Float64bits(want))
		}
		if fetch.FoldScores != 1 || fetch.BlockSuccess+fetch.BlockFailed != fetch.BlockReads {
			t.Fatalf("telemetry %+v", fetch)
		}
	}
}

func TestExperimentalV4Build86ContinueMatchesBuild84Exactly(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build86SyntheticFixture(t)
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, we, _ := experimentalV4PhoneBuild84ContinueFetch(plane, candidate, w, h, anchor, start, 0)
	got, ge, fetch := experimentalV4PhoneBuild86ContinueFetch(plane, candidate, w, h, anchor, start, 0)
	if ge != we || len(got) != len(want) {
		t.Fatalf("eval/state %d/%d %d/%d", ge, len(got), we, len(want))
	}
	for i := range got {
		if got[i].index != want[i].index || got[i].pass != want[i].pass || got[i].dim != want[i].dim || got[i].delta != want[i].delta || got[i].hyp.quad != want[i].hyp.quad || got[i].hyp.h != want[i].hyp.h || math.Float64bits(got[i].hyp.proposal) != math.Float64bits(want[i].hyp.proposal) {
			t.Fatalf("state %d", i)
		}
	}
	if we > 0 && (fetch.FoldScores != we || fetch.BlockReads == 0 || fetch.BlockSuccess+fetch.BlockFailed != fetch.BlockReads) {
		t.Fatalf("fetch telemetry %+v", fetch)
	}
}

func TestExperimentalV4Build86BlindBankMatchesBuild84Exactly(t *testing.T) {
	_, carrier, _, w, h, _, _ := build86SyntheticFixture(t)
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, we, ws, ww, _, _ := experimentalV4PhoneBuild84BlindBank(carrier, boundary, w, h)
	got, ge, gs, gw, gp, fp := experimentalV4PhoneBuild86BlindBank(carrier, boundary, w, h)
	if ge != we || gs != ws || gw != ww || len(got) != len(want) {
		t.Fatalf("bank summary eval=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", ge, we, gs, ws, gw, ww, len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || math.Float64bits(got[i].proposal) != math.Float64bits(want[i].proposal) || math.Float64bits(got[i].validation) != math.Float64bits(want[i].validation) {
			t.Fatalf("bank %d", i)
		}
	}
	if fp.BasinTasks != 16 || gp.Gen2Tasks == 0 || gp.Gen3Tasks == 0 || gp.FetchBlockSuccess+gp.FetchBlockFailed != gp.FetchBlockReads {
		t.Fatalf("telemetry gp=%+v fp=%+v", gp, fp)
	}
}

func TestExperimentalV4Build86ApplyTelemetryPreservesBuild84View(t *testing.T) {
	r := experimentalV4PhoneBuild86RecoveryTelemetry{}
	r.Attempted = true
	r.Build86GeometryProfile = experimentalV4PhoneBuild86GeometryTelemetry{Prefix1Workers: 8, Gen2Workers: 8, Gen3Workers: 8, Gen4Workers: 8, Gen2Tasks: 12, Gen3Tasks: 34, Gen4Tasks: 56, FetchFoldScores: 90, FetchBlockReads: 1000, FetchBlockSuccess: 980, FetchBlockFailed: 20}
	var p ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild86ApplyTelemetry(&p, r)
	if !p.Build86Attempted || !p.Build76Attempted || p.Build76Gen4Tasks != 56 || !p.Build84Attempted || p.Build84FetchFoldScores != 90 || p.Build84FetchBlockReads != 1000 || p.Build84FetchBlockSuccess != 980 || p.Build84FetchBlockFailed != 20 || p.Build86DCTFoldScores != 90 || p.Build86DCTBlockReads != 1000 || p.Build86DCTBlockSuccess != 980 || p.Build86DCTBlockFailed != 20 {
		t.Fatalf("telemetry %+v", p)
	}
}
