package watermark

import (
	"reflect"
	"testing"
)

func experimentalV4Build97Fixture(t testing.TB) (*pixelPlane, experimentalV4PilotCandidate, int, int, [4]ImagePoint, experimentalV4PhoneHypothesis) {
	t.Helper()
	pilot := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, pilot, 48, experimentalV4SyntheticDataSeed)
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
	score, _ := experimentalV4PhoneBuild41FoldScore(plane, pilot, w, h, hom, 0, true, false)
	root := experimentalV4PhoneHypothesis{quad: anchor, h: hom, proposal: score}
	return plane, pilot, w, h, anchor, root
}

func TestExperimentalV4Build97SingleProbeMatchesBuild53Exactly(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	we, wi := experimentalV4PhoneBuild53SingleProbe(plane, pilot, w, h, anchor, root, 0)
	ge, gi, gf := experimentalV4PhoneBuild97SingleProbeFetch(plane, pilot, w, h, anchor, root, 0)
	if ge != we || gi != wi {
		t.Fatalf("single mismatch evals=%d/%d improving=%d/%d", ge, we, gi, wi)
	}
	if gf.FoldScores != ge || gf.BlockReads <= 0 || gf.BlockSuccess <= 0 {
		t.Fatalf("single fetch telemetry invalid: %+v", gf)
	}
}

func TestExperimentalV4Build97PairStencilMatchesBuild53Exactly(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	want, we, wi := experimentalV4PhoneBuild53PairStencil(plane, pilot, w, h, anchor, root, 0)
	got, ge, gi, gf := experimentalV4PhoneBuild97PairStencilFetch(plane, pilot, w, h, anchor, root, 0)
	if ge != we || gi != wi || !reflect.DeepEqual(got, want) {
		t.Fatalf("pair mismatch evals=%d/%d improving=%d/%d got=%d want=%d", ge, we, gi, wi, len(got), len(want))
	}
	if gf.FoldScores != ge || gf.BlockReads <= 0 || gf.BlockSuccess <= 0 {
		t.Fatalf("pair fetch telemetry invalid: %+v", gf)
	}
}

func TestExperimentalV4Build97Generation4MatchesBuild84Exactly(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	want, we, wf := experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
	got, ge, gf := experimentalV4PhoneBuild97Generation4Fetch(plane, pilot, w, h, anchor, root)
	if ge != we || !reflect.DeepEqual(got, want) {
		t.Fatalf("gen4 mismatch evals=%d/%d got=%d want=%d", ge, we, len(got), len(want))
	}
	if gf.Continuation != wf {
		t.Fatalf("Build97 historical continuation telemetry changed: got=%+v Build84=%+v", gf.Continuation, wf)
	}
	if gf.Probe.FoldScores <= 0 || gf.Probe.BlockReads <= 0 || gf.Probe.BlockReads != gf.Probe.BlockSuccess+gf.Probe.BlockFailed {
		t.Fatalf("Build97 probe telemetry invalid: %+v", gf.Probe)
	}
}

func TestExperimentalV4Build97ParallelGeneration4RaceFixture(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	const workers = 8
	errs := make(chan string, workers)
	for i := 0; i < workers; i++ {
		go func() {
			want, we, _ := experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
			got, ge, _ := experimentalV4PhoneBuild97Generation4Fetch(plane, pilot, w, h, anchor, root)
			if ge != we || !reflect.DeepEqual(got, want) {
				errs <- "mismatch"
				return
			}
			errs <- ""
		}()
	}
	for i := 0; i < workers; i++ {
		if e := <-errs; e != "" {
			t.Fatal(e)
		}
	}
}
