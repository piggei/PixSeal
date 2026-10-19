package watermark

import (
	"bytes"
	"image"
	"math"
	"reflect"
	"sync"
	"testing"
	"time"
)

func experimentalV4Build95AssertStatesEqual(t *testing.T, want, got []experimentalV4PhoneBuild55BlindContinuationState) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("state count=%d want=%d", len(got), len(want))
	}
	for i := range got {
		if got[i].index != want[i].index || got[i].pass != want[i].pass || got[i].dim != want[i].dim || got[i].delta != want[i].delta || got[i].hyp.quad != want[i].hyp.quad || got[i].hyp.h != want[i].hyp.h || math.Float64bits(got[i].hyp.proposal) != math.Float64bits(want[i].hyp.proposal) {
			t.Fatalf("state %d mismatch\nwant=%+v\ngot=%+v", i, want[i], got[i])
		}
	}
}

func experimentalV4Build95AssertHypothesesEqual(t *testing.T, want, got []experimentalV4PhoneHypothesis) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("hypothesis count=%d want=%d", len(got), len(want))
	}
	for i := range got {
		if got[i].quad != want[i].quad || got[i].h != want[i].h || math.Float64bits(got[i].proposal) != math.Float64bits(want[i].proposal) || math.Float64bits(got[i].validation) != math.Float64bits(want[i].validation) || got[i].detection != want[i].detection {
			t.Fatalf("hypothesis %d mismatch\nwant=%+v\ngot=%+v", i, want[i], got[i])
		}
	}
}

func TestExperimentalV4Build95LuminancePlaneBitExact(t *testing.T) {
	_, _, plane, _, _, _, _ := build84SyntheticFixture(t)
	luma, bytes, fallback := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	if fallback || luma == nil {
		t.Fatal("unexpected luminance-plane fallback")
	}
	wantBytes := int64(plane.bounds.Dx()) * int64(plane.bounds.Dy()) * 8
	if bytes != wantBytes || int64(len(luma.luma))*8 != wantBytes {
		t.Fatalf("bytes=%d len-bytes=%d want=%d", bytes, int64(len(luma.luma))*8, wantBytes)
	}
	for i := range luma.luma {
		j := i * 3
		want := .299*float64(plane.rgb[j]) + .587*float64(plane.rgb[j+1]) + .114*float64(plane.rgb[j+2]) - 128
		if math.Float64bits(luma.luma[i]) != math.Float64bits(want) {
			t.Fatalf("pixel %d bits=%x want=%x", i, math.Float64bits(luma.luma[i]), math.Float64bits(want))
		}
	}
}

func TestExperimentalV4Build95LuminancePlaneBoundedFallback(t *testing.T) {
	if luma, bytes, fallback := experimentalV4PhoneBuild95NewLuminancePlane(nil); luma != nil || bytes != 0 || !fallback {
		t.Fatalf("nil result luma=%v bytes=%d fallback=%t", luma != nil, bytes, fallback)
	}
	invalid := &pixelPlane{bounds: image.Rect(0, 0, 2, 2), rgb: make([]uint8, 11)}
	if luma, bytes, fallback := experimentalV4PhoneBuild95NewLuminancePlane(invalid); luma != nil || bytes != 32 || !fallback {
		t.Fatalf("invalid-storage result luma=%v bytes=%d fallback=%t", luma != nil, bytes, fallback)
	}
	oversized := &pixelPlane{bounds: image.Rect(0, 0, experimentalV4PhoneMaxDimension, experimentalV4PhoneMaxDimension)}
	luma, bytes, fallback := experimentalV4PhoneBuild95NewLuminancePlane(oversized)
	if luma != nil || !fallback || bytes <= experimentalV4PhoneBuild95MaxLuminancePlaneBytes {
		t.Fatalf("oversized result luma=%v bytes=%d fallback=%t", luma != nil, bytes, fallback)
	}
	retainedBytes := int64(3464) * int64(4600) * 8
	if retainedBytes > experimentalV4PhoneBuild95MaxLuminancePlaneBytes {
		t.Fatalf("retained 3464x4600 plane=%d exceeds cap=%d", retainedBytes, experimentalV4PhoneBuild95MaxLuminancePlaneBytes)
	}
}

func TestExperimentalV4Build95ProjectiveBlockMatchesBuild84BitExactly(t *testing.T) {
	_, _, plane, w, h, anchor, _ := build84SyntheticFixture(t)
	luma, _, fallback := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	if fallback {
		t.Fatal("unexpected fallback")
	}
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
			t.Fatalf("homography %d", qi)
		}
		for oy := 0; oy+blockSize <= h; oy += 31 {
			for ox := 0; ox+blockSize <= w; ox += 37 {
				origins = append(origins, [2]int{ox, oy})
			}
		}
		for _, origin := range origins {
			want, wok := experimentalV4PhoneBuild84ReadProjectiveBlockValue(plane, hom, origin[0], origin[1], blockSize)
			got, gok := experimentalV4PhoneBuild95ReadProjectiveBlockValue(luma, hom, origin[0], origin[1], blockSize)
			if gok != wok || (gok && math.Float64bits(got) != math.Float64bits(want)) {
				t.Fatalf("q=%d origin=%v ok=%t/%t bits=%x/%x", qi, origin, gok, wok, math.Float64bits(got), math.Float64bits(want))
			}
		}
	}
}

func TestExperimentalV4Build95FoldScoreMatchesBuild84BitExactly(t *testing.T) {
	candidate, _, plane, w, h, _, hom := build84SyntheticFixture(t)
	luma, _, _ := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	for _, tc := range []struct{ proposal, sparse bool }{{true, false}, {false, false}, {true, true}} {
		want, wv, wf := experimentalV4PhoneBuild84FoldScoreFetch(plane, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		got, gv, gf := experimentalV4PhoneBuild95FoldScoreFetch(luma, candidate, w, h, hom, 0, tc.proposal, tc.sparse)
		if gv != wv || math.Float64bits(got) != math.Float64bits(want) || gf != wf {
			t.Fatalf("fold proposal=%t sparse=%t visible=%d/%d bits=%x/%x fetch=%+v/%+v", tc.proposal, tc.sparse, gv, wv, math.Float64bits(got), math.Float64bits(want), gf, wf)
		}
	}
}

func TestExperimentalV4Build95ContinueMatchesBuild84Exactly(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build84SyntheticFixture(t)
	luma, _, _ := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, we, wf := experimentalV4PhoneBuild84ContinueFetch(plane, candidate, w, h, anchor, start, 0)
	got, ge, gf := experimentalV4PhoneBuild95ContinueFetch(luma, candidate, w, h, anchor, start, 0)
	if ge != we || gf != wf {
		t.Fatalf("eval/fetch mismatch eval=%d/%d fetch=%+v/%+v", ge, we, gf, wf)
	}
	experimentalV4Build95AssertStatesEqual(t, want, got)
}

func TestExperimentalV4Build95Generation4MatchesBuild84Exactly(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build84SyntheticFixture(t)
	luma, _, _ := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	input := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	want, we, wf := experimentalV4PhoneBuild84Generation4Fetch(plane, candidate, w, h, anchor, input)
	got, ge, gf := experimentalV4PhoneBuild95Generation4Fetch(plane, luma, candidate, w, h, anchor, input)
	if ge != we || gf != wf {
		t.Fatalf("eval/fetch mismatch eval=%d/%d fetch=%+v/%+v", ge, we, gf, wf)
	}
	experimentalV4Build95AssertHypothesesEqual(t, want, got)

	fallback, fe, ff := experimentalV4PhoneBuild95Generation4Fetch(plane, nil, candidate, w, h, anchor, input)
	if fe != we || ff != wf {
		t.Fatalf("fallback eval/fetch mismatch eval=%d/%d fetch=%+v/%+v", fe, we, ff, wf)
	}
	experimentalV4Build95AssertHypothesesEqual(t, want, fallback)
}

func TestExperimentalV4Build95BlindBankMatchesBuild84Exactly(t *testing.T) {
	_, carrier, _, w, h, _, _ := build84SyntheticFixture(t)
	boundary := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: ImagePoint{X: 0, Y: 0}, TopRight: ImagePoint{X: float64(w - 1), Y: 0}, BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)}, BottomLeft: ImagePoint{X: 0, Y: float64(h - 1)}}
	want, we, ws, ww, _, _ := experimentalV4PhoneBuild84BlindBank(carrier, boundary, w, h)
	got, ge, gs, gw, gp, fp := experimentalV4PhoneBuild95BlindBank(carrier, boundary, w, h)
	if ge != we || gs != ws || gw != ww || len(got) != len(want) {
		t.Fatalf("bank summary eval=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", ge, we, gs, ws, gw, ww, len(got), len(want))
	}
	experimentalV4Build95AssertHypothesesEqual(t, want, got)
	if fp.BasinTasks != 16 || gp.FetchBlockSuccess+gp.FetchBlockFailed != gp.FetchBlockReads {
		t.Fatalf("telemetry gp=%+v fp=%+v", gp, fp)
	}
	if gp.Gen4Tasks > 0 {
		if !gp.LuminancePlaneUsed || gp.LuminancePlaneFallback || gp.LuminancePlaneBytes <= 0 || gp.LuminancePlanePrepElapsed <= 0 {
			t.Fatalf("Gen4 luminance telemetry gp=%+v", gp)
		}
	} else if gp.LuminancePlaneUsed || gp.LuminancePlaneBytes != 0 || gp.LuminancePlanePrepElapsed != 0 {
		t.Fatalf("unexpected luminance preparation without Gen4: %+v", gp)
	}
}

func TestExperimentalV4Build95PublicEntryValidationMatchesBuild93(t *testing.T) {
	validImage := image.NewRGBA(image.Rect(0, 0, 16, 16))
	cases := []struct {
		name   string
		src    image.Image
		key    []byte
		width  int
		height int
	}{
		{name: "nil-image", src: nil, key: []byte("Build95-Key"), width: 1632, height: 1632},
		{name: "short-key", src: validImage, key: []byte("short"), width: 1632, height: 1632},
		{name: "invalid-canonical", src: validImage, key: []byte("Build95-Key"), width: 288, height: 248},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantPayload, wantInfo, wantPhone, wantErr := ExperimentalV4ExtractPhoneBuild93Candidate(tc.src, tc.key, tc.width, tc.height)
			gotPayload, gotInfo, gotPhone, gotErr := ExperimentalV4ExtractPhoneBuild95Candidate(tc.src, tc.key, tc.width, tc.height)
			if !bytes.Equal(gotPayload, wantPayload) || !reflect.DeepEqual(gotInfo, wantInfo) || !reflect.DeepEqual(gotPhone, wantPhone) {
				t.Fatalf("entry points diverged on %s", tc.name)
			}
			if (gotErr == nil) != (wantErr == nil) || (gotErr != nil && gotErr.Error() != wantErr.Error()) {
				t.Fatalf("error mismatch on %s: got=%v want=%v", tc.name, gotErr, wantErr)
			}
		})
	}
}

func TestExperimentalV4Build95ApplyTelemetryPreservesBuild93View(t *testing.T) {
	r := experimentalV4PhoneBuild95RecoveryTelemetry{}
	r.Attempted = true
	r.ProtectedSamplingPasses = 696
	r.ProtectedProfileDerivations = 2086
	r.Build95GeometryProfile.LuminancePlaneUsed = true
	r.Build95GeometryProfile.LuminancePlaneBytes = 123456
	r.Build95GeometryProfile.LuminancePlanePrepElapsed = 7 * time.Millisecond
	var p ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild95ApplyTelemetry(&p, r)
	if !p.Build93Attempted || p.Build93ProtectedSamplingPasses != 696 || p.Build93ProtectedProfileDerivations != 2086 || !p.Build95Attempted || !p.Build95LuminancePlaneUsed || p.Build95LuminancePlaneFallback || p.Build95LuminancePlaneBytes != 123456 || p.Build95LuminancePlanePrepMs != 7 {
		t.Fatalf("telemetry %+v", p)
	}
}

func TestExperimentalV4Build95ParallelContinuationRaceFixture(t *testing.T) {
	candidate, _, plane, w, h, anchor, hom := build84SyntheticFixture(t)
	luma, _, fallback := experimentalV4PhoneBuild95NewLuminancePlane(plane)
	if fallback || luma == nil {
		t.Fatal("unexpected luminance-plane fallback")
	}
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, candidate, w, h, hom, 0, true, false)
	start := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	wantStates, wantEvals, wantFetch := experimentalV4PhoneBuild95ContinueFetch(luma, candidate, w, h, anchor, start, 0)

	const workers = 8
	var wg sync.WaitGroup
	errs := make(chan string, workers)
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func() {
			defer wg.Done()
			for iteration := 0; iteration < 4; iteration++ {
				gotStates, gotEvals, gotFetch := experimentalV4PhoneBuild95ContinueFetch(luma, candidate, w, h, anchor, start, 0)
				if gotEvals != wantEvals || gotFetch != wantFetch || !reflect.DeepEqual(gotStates, wantStates) {
					errs <- "parallel continuation result diverged"
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
}
