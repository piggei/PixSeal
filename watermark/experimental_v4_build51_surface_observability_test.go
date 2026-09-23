package watermark

import "testing"

func TestExperimentalV4Build51StencilIsDeterministicAndBounded(t *testing.T) {
	start := [4]ImagePoint{{X: 100, Y: 100}, {X: 500, Y: 100}, {X: 100, Y: 500}, {X: 500, Y: 500}}
	defs := experimentalV4PhoneBuild51StencilQuads(start)
	if len(defs) != 53 {
		t.Fatalf("stencil defs=%d want 53", len(defs))
	}
	if defs[0].name != "center" || defs[0].delta != 0 {
		t.Fatalf("first stencil=%q delta=%v want center/0", defs[0].name, defs[0].delta)
	}
	seen := map[string]int{}
	for _, d := range defs {
		seen[d.name]++
	}
	if seen["corner-0-x"] != 4 || seen["corner-3-y"] != 4 {
		t.Fatalf("corner stencil counts=%v", seen)
	}
	for _, name := range []string{"translate-x", "translate-y", "scale-x", "scale-y", "shear-x", "shear-y", "top-width", "bottom-width", "left-height", "right-height"} {
		if seen[name] != 2 {
			t.Fatalf("basis %s count=%d want 2", name, seen[name])
		}
	}
}

func TestExperimentalV4Build51UsesTop4WithoutChangingHistoricalDepths(t *testing.T) {
	if experimentalV4PhoneBuild43PairKeep != 2 || experimentalV4PhoneBuild43MaxFrozen != 32 {
		t.Fatalf("Build43 production changed: pairs=%d frozen=%d", experimentalV4PhoneBuild43PairKeep, experimentalV4PhoneBuild43MaxFrozen)
	}
	if experimentalV4PhoneBuild48SeedsPerPair != 2 {
		t.Fatalf("Build48 historical depth=%d want 2", experimentalV4PhoneBuild48SeedsPerPair)
	}
	if experimentalV4PhoneBuild50SeedsPerPair != 4 || experimentalV4PhoneBuild51SeedsPerPair != 4 {
		t.Fatalf("Build50/51 diagnostic depth=%d/%d want 4/4", experimentalV4PhoneBuild50SeedsPerPair, experimentalV4PhoneBuild51SeedsPerPair)
	}
}

func TestExperimentalV4Build51JointPerturbationsAreLocal(t *testing.T) {
	start := [4]ImagePoint{{X: 100, Y: 100}, {X: 500, Y: 100}, {X: 100, Y: 500}, {X: 500, Y: 500}}
	q := experimentalV4PhoneBuild51Joint(start, "translate-x", 2)
	for i := range q {
		if q[i].X != start[i].X+2 || q[i].Y != start[i].Y {
			t.Fatalf("translate-x corner %d=%+v", i, q[i])
		}
	}
	q = experimentalV4PhoneBuild51Joint(start, "top-width", 2)
	if q[0].X != 98 || q[1].X != 502 || q[2] != start[2] || q[3] != start[3] {
		t.Fatalf("top-width=%+v", q)
	}
}
