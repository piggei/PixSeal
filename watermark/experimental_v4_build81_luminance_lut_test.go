package watermark

import (
	"image"
	"math"
	"testing"
)

func build81SyntheticFixture(t *testing.T) (experimentalV4PilotCandidate, image.Image, *pixelPlane, int, int, [4]ImagePoint, homography) {
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

func TestExperimentalV4Build81ProductTablesMatchDirectMultiplicationBitExactly(t *testing.T) {
	for i := 0; i < 256; i++ {
		v := float64(i)
		if math.Float64bits(experimentalV4PhoneBuild81LumaR[i]) != math.Float64bits(.299*v) {
			t.Fatalf("R %d", i)
		}
		if math.Float64bits(experimentalV4PhoneBuild81LumaG[i]) != math.Float64bits(.587*v) {
			t.Fatalf("G %d", i)
		}
		if math.Float64bits(experimentalV4PhoneBuild81LumaB[i]) != math.Float64bits(.114*v) {
			t.Fatalf("B %d", i)
		}
	}
}

func TestExperimentalV4Build81RGBLuminanceMatchesOriginalForEveryRGBTriplet(t *testing.T) {
	for r := 0; r < 256; r++ {
		for g := 0; g < 256; g++ {
			for b := 0; b < 256; b++ {
				want := .299*float64(r) + .587*float64(g) + .114*float64(b) - 128
				got := experimentalV4PhoneBuild81RGBLuminance(uint8(r), uint8(g), uint8(b))
				if math.Float64bits(got) != math.Float64bits(want) {
					t.Fatalf("RGB=(%d,%d,%d)", r, g, b)
				}
			}
		}
	}
}

func TestExperimentalV4Build81SamplePlaneLuminanceMatchesOriginalBitExactly(t *testing.T) {
	_, _, plane, w, h, _, _ := build81SyntheticFixture(t)
	coords := [][2]float64{{0, 0}, {float64(w - 1), float64(h - 1)}, {10.25, 11.75}, {101.5, 207.125}, {float64(w) - 1.001, 3.75}, {4.5, float64(h) - 1.125}, {-0.001, 10}, {10, float64(h) + 0.01}}
	for _, c := range coords {
		want, wok := samplePlaneLuminance(plane, c[0], c[1])
		got, gok := experimentalV4PhoneBuild81SamplePlaneLuminance(plane, c[0], c[1])
		if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
			t.Fatalf("coord=%v", c)
		}
	}
}

func TestExperimentalV4Build81ProjectiveBlockMatchesOriginalBitExactly(t *testing.T) {
	_, _, plane, w, h, anchor, _ := build81SyntheticFixture(t)
	quads := [][4]ImagePoint{anchor, {{X: 3, Y: 2}, {X: float64(w - 5), Y: 5}, {X: 6, Y: float64(h - 4)}, {X: float64(w - 7), Y: float64(h - 6)}}, {{X: 14, Y: 5}, {X: float64(w - 18), Y: 11}, {X: 4, Y: float64(h - 15)}, {X: float64(w - 9), Y: float64(h - 7)}}}
	for qi, q := range quads {
		hom, ok := experimentalV4PhoneQuadHomography(w, h, q)
		if !ok {
			t.Fatalf("hom %d", qi)
		}
		for oy := 0; oy+blockSize <= h; oy += 31 {
			for ox := 0; ox+blockSize <= w; ox += 37 {
				want, wok := readProjectiveBlockValue(plane, hom, ox, oy, blockSize)
				got, gok := experimentalV4PhoneBuild81ReadProjectiveBlockValue(plane, hom, ox, oy, blockSize)
				if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
					t.Fatalf("q=%d block=%d,%d", qi, ox, oy)
				}
			}
		}
	}
}

func TestExperimentalV4Build81FoldScoreMatchesBuild41BitExactly(t *testing.T) {
	candidate, _, plane, w, h, _, hom := build81SyntheticFixture(t)
	for _, tc := range []struct{ proposal, sparse bool }{{true, false}, {false, false}, {true, true}} {
		want, wv := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		got, gv, lut := experimentalV4PhoneBuild81FoldScoreLUT(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		if gv != wv || math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("fold")
		}
		if lut.FoldScores != 1 || lut.BlockSuccess+lut.BlockFailed != lut.BlockReads {
			t.Fatalf("telemetry %+v", lut)
		}
	}
}

func TestExperimentalV4Build81ContinueMatchesBuild55Exactly(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build81SyntheticFixture(t)
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, we := experimentalV4PhoneBuild55Continue(plane, candidate, w, h, anchor, start, 0)
	got, ge, lut := experimentalV4PhoneBuild81ContinueLUT(plane, candidate, w, h, anchor, start, 0)
	if ge != we || len(got) != len(want) {
		t.Fatalf("eval/state %d/%d %d/%d", ge, len(got), we, len(want))
	}
	for i := range got {
		if got[i].index != want[i].index || got[i].pass != want[i].pass || got[i].dim != want[i].dim || got[i].delta != want[i].delta || got[i].hyp.quad != want[i].hyp.quad || got[i].hyp.h != want[i].hyp.h || math.Float64bits(got[i].hyp.proposal) != math.Float64bits(want[i].hyp.proposal) {
			t.Fatalf("state %d", i)
		}
	}
	if we > 0 && (lut.FoldScores != we || lut.BlockReads == 0) {
		t.Fatalf("lut %+v", lut)
	}
}

func TestExperimentalV4Build81BlindBankMatchesBuild76Exactly(t *testing.T) {
	_, carrier, _, w, h, _, _ := build81SyntheticFixture(t)
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, we, ws, ww, _, _ := experimentalV4PhoneBuild76BlindBank(carrier, boundary, w, h)
	got, ge, gs, gw, gp, fp := experimentalV4PhoneBuild81BlindBank(carrier, boundary, w, h)
	if ge != we || gs != ws || gw != ww || len(got) != len(want) {
		t.Fatalf("bank summary")
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || math.Float64bits(got[i].proposal) != math.Float64bits(want[i].proposal) || math.Float64bits(got[i].validation) != math.Float64bits(want[i].validation) {
			t.Fatalf("bank %d", i)
		}
	}
	if fp.BasinTasks != 16 || gp.Gen2Tasks == 0 || gp.Gen3Tasks == 0 || gp.LUTBlockSuccess+gp.LUTBlockFailed != gp.LUTBlockReads {
		t.Fatalf("telemetry")
	}
}

func TestExperimentalV4Build81ApplyTelemetryPreservesBuild76View(t *testing.T) {
	r := experimentalV4PhoneBuild81RecoveryTelemetry{}
	r.Attempted = true
	r.Build81GeometryProfile = experimentalV4PhoneBuild81GeometryTelemetry{Prefix1Workers: 8, Gen2Workers: 8, Gen3Workers: 8, Gen4Workers: 8, Gen2Tasks: 12, Gen3Tasks: 34, Gen4Tasks: 56, LUTFoldScores: 90, LUTBlockReads: 1000, LUTBlockSuccess: 980, LUTBlockFailed: 20}
	var p ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild81ApplyTelemetry(&p, r)
	if !p.Build81Attempted || !p.Build76Attempted || p.Build76Gen4Tasks != 56 || p.Build81LUTFoldScores != 90 || p.Build81LUTBlockReads != 1000 || p.Build81LUTBlockSuccess != 980 || p.Build81LUTBlockFailed != 20 {
		t.Fatalf("telemetry %+v", p)
	}
}
