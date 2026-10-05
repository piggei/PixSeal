package watermark

import (
	"image"
	"math"
	"testing"
)

func build82SyntheticFixture(t *testing.T) (experimentalV4PilotCandidate, image.Image, *pixelPlane, int, int, [4]ImagePoint, homography) {
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

func TestExperimentalV4Build82ProjectiveBlockMatchesOriginalBitExactly(t *testing.T) {
	_, _, plane, w, h, anchor, _ := build82SyntheticFixture(t)
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
			want, wok := readProjectiveBlockValue(plane, hom, origin[0], origin[1], blockSize)
			got, gok := experimentalV4PhoneBuild82ReadProjectiveBlockValue(plane, hom, origin[0], origin[1], blockSize)
			if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
				t.Fatalf("q=%d block=%d,%d ok=%t/%t bits=%x/%x", qi, origin[0], origin[1], gok, wok, math.Float64bits(got), math.Float64bits(want))
			}
		}
	}
}

func TestExperimentalV4Build82FoldScoreMatchesBuild41BitExactly(t *testing.T) {
	candidate, _, plane, w, h, _, hom := build82SyntheticFixture(t)
	for _, tc := range []struct{ proposal, sparse bool }{{true, false}, {false, false}, {true, true}} {
		want, wv := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		got, gv, inline := experimentalV4PhoneBuild82FoldScoreInline(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		if gv != wv || math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("fold proposal=%t sparse=%t visible=%d/%d bits=%x/%x", tc.proposal, tc.sparse, gv, wv, math.Float64bits(got), math.Float64bits(want))
		}
		if inline.FoldScores != 1 || inline.BlockSuccess+inline.BlockFailed != inline.BlockReads {
			t.Fatalf("telemetry %+v", inline)
		}
	}
}

func TestExperimentalV4Build82ContinueMatchesBuild55Exactly(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build82SyntheticFixture(t)
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, we := experimentalV4PhoneBuild55Continue(plane, candidate, w, h, anchor, start, 0)
	got, ge, inline := experimentalV4PhoneBuild82ContinueInline(plane, candidate, w, h, anchor, start, 0)
	if ge != we || len(got) != len(want) {
		t.Fatalf("eval/state %d/%d %d/%d", ge, len(got), we, len(want))
	}
	for i := range got {
		if got[i].index != want[i].index || got[i].pass != want[i].pass || got[i].dim != want[i].dim || got[i].delta != want[i].delta || got[i].hyp.quad != want[i].hyp.quad || got[i].hyp.h != want[i].hyp.h || math.Float64bits(got[i].hyp.proposal) != math.Float64bits(want[i].hyp.proposal) {
			t.Fatalf("state %d", i)
		}
	}
	if we > 0 && (inline.FoldScores != we || inline.BlockReads == 0 || inline.BlockSuccess+inline.BlockFailed != inline.BlockReads) {
		t.Fatalf("inline telemetry %+v", inline)
	}
}

func TestExperimentalV4Build82BlindBankMatchesBuild76Exactly(t *testing.T) {
	_, carrier, _, w, h, _, _ := build82SyntheticFixture(t)
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, we, ws, ww, _, _ := experimentalV4PhoneBuild76BlindBank(carrier, boundary, w, h)
	got, ge, gs, gw, gp, fp := experimentalV4PhoneBuild82BlindBank(carrier, boundary, w, h)
	if ge != we || gs != ws || gw != ww || len(got) != len(want) {
		t.Fatalf("bank summary eval=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", ge, we, gs, ws, gw, ww, len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || math.Float64bits(got[i].proposal) != math.Float64bits(want[i].proposal) || math.Float64bits(got[i].validation) != math.Float64bits(want[i].validation) {
			t.Fatalf("bank %d", i)
		}
	}
	if fp.BasinTasks != 16 || gp.Gen2Tasks == 0 || gp.Gen3Tasks == 0 || gp.InlineBlockSuccess+gp.InlineBlockFailed != gp.InlineBlockReads {
		t.Fatalf("telemetry gp=%+v fp=%+v", gp, fp)
	}
}

func TestExperimentalV4Build82ApplyTelemetryPreservesBuild76View(t *testing.T) {
	r := experimentalV4PhoneBuild82RecoveryTelemetry{}
	r.Attempted = true
	r.Build82GeometryProfile = experimentalV4PhoneBuild82GeometryTelemetry{Prefix1Workers: 8, Gen2Workers: 8, Gen3Workers: 8, Gen4Workers: 8, Gen2Tasks: 12, Gen3Tasks: 34, Gen4Tasks: 56, InlineFoldScores: 90, InlineBlockReads: 1000, InlineBlockSuccess: 980, InlineBlockFailed: 20}
	var p ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild82ApplyTelemetry(&p, r)
	if !p.Build82Attempted || !p.Build76Attempted || p.Build76Gen4Tasks != 56 || p.Build82InlineFoldScores != 90 || p.Build82InlineBlockReads != 1000 || p.Build82InlineBlockSuccess != 980 || p.Build82InlineBlockFailed != 20 {
		t.Fatalf("telemetry %+v", p)
	}
}
