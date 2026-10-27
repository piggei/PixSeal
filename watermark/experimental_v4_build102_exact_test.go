package watermark

import (
	"math"
	"reflect"
	"runtime"
	"testing"
)

func TestExperimentalV4Build102FingerprintGuard(t *testing.T) {
	a := experimentalV4PhoneHypothesis{proposal: 1, build43Pair: "fixture", quad: [4]ImagePoint{{X: 3}}}
	b := a
	b.proposal = math.Nextafter(a.proposal, 2)
	if experimentalV4Build102Hash(a) == experimentalV4Build102Hash(b) {
		t.Fatal("unexported proposal mutation accepted")
	}
	b = a
	b.quad[0].X++
	if experimentalV4Build102Hash(a) == experimentalV4Build102Hash(b) {
		t.Fatal("geometry mutation accepted")
	}
	if experimentalV4Build102Hash([]experimentalV4PhoneHypothesis{a, b}) == experimentalV4Build102Hash([]experimentalV4PhoneHypothesis{b, a}) {
		t.Fatal("order mutation accepted")
	}
	if experimentalV4Build102Hash(0.0) == experimentalV4Build102Hash(math.Copysign(0, -1)) {
		t.Fatal("float bit mutation accepted")
	}
	x, y := a, a
	if experimentalV4Build102Hash(&x) != experimentalV4Build102Hash(&y) {
		t.Fatal("fingerprint depends on pointer address")
	}
	if experimentalV4Build102Hash([]experimentalV4PhoneHypothesis(nil)) == experimentalV4Build102Hash([]experimentalV4PhoneHypothesis{}) {
		t.Fatal("nil marker lost")
	}
}

func TestExperimentalV4Build102ProcessCPUParser(t *testing.T) {
	u, s, err := experimentalV4Build102ParseCPU("123 (name with ) parens) R 1 2 3 4 5 6 7 8 9 10 12345 6789 99")
	if err != nil || u != 12345 || s != 6789 {
		t.Fatalf("CPU fields: %d %d %v", u, s, err)
	}
	for _, raw := range []string{"bad", "12 (x) R", "12 (x) R 1 2 3 4 5 6 7 8 9 10 invalid 4"} {
		if _, _, err := experimentalV4Build102ParseCPU(raw); err == nil {
			t.Fatal("accepted invalid process stat")
		}
	}
	if runtime.GOOS == "linux" {
		if _, _, err := experimentalV4Build102CPU(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestExperimentalV4Build102SplitBankMatchesQualified(t *testing.T) {
	carrier, _, _, w, h, _, boundary := experimentalV4Build100Fixture(t)
	want, we, ws, ww, wg, wf := experimentalV4PhoneBuild98BlindBank(carrier, boundary, w, h)
	var taskHash string
	for _, candidate := range []bool{false, true} {
		f := experimentalV4Build102Prepare(carrier, boundary, w, h, candidate)
		before := experimentalV4Build102Hash(f.tasks)
		if taskHash == "" {
			taskHash = before
		} else if before != taskHash {
			t.Fatal("prefix contexts changed frozen vector")
		}
		got, ge, gs, gw, gg, gf := experimentalV4Build102Replay(f)
		gp := gg.experimentalV4PhoneBuild98GeometryTelemetry
		wgCopy, wfCopy := wg, wf
		for _, v := range []any{&gp, &wgCopy, &gf, &wfCopy} {
			experimentalV4Build100ClearDurations(reflect.ValueOf(v).Elem())
		}
		if experimentalV4Build102Hash(got) != experimentalV4Build102Hash(want) || ge != we || gs != ws || gw != ww || !reflect.DeepEqual(gp, wgCopy) || !reflect.DeepEqual(gf, wfCopy) {
			t.Fatalf("split candidate=%t changed qualified bank or counters", candidate)
		}
		if before != experimentalV4Build102Hash(f.tasks) {
			t.Fatal("replay mutated frozen vector")
		}
	}
}

func TestExperimentalV4Build102ReplayRaceFixture(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	want, we, wf := experimentalV4PhoneBuild84Generation4Fetch(plane, pilot, w, h, anchor, root)
	f := experimentalV4Build102Frozen{plane: plane, pilot: pilot, cw: w, ch: h, anchor: anchor, tasks: []experimentalV4PhoneHypothesis{root, root}}
	before := experimentalV4Build102Hash(f.tasks)
	got, ge, _, _, gg, _ := experimentalV4Build102Replay(f)
	if experimentalV4Build102Hash(got) != experimentalV4Build102Hash(append(append([]experimentalV4PhoneHypothesis{}, want...), want...)) || ge != we*2 || gg.FetchBlockReads != wf.BlockReads*2 || before != experimentalV4Build102Hash(f.tasks) {
		t.Fatal("concurrent replay differs or mutates inputs")
	}
}
