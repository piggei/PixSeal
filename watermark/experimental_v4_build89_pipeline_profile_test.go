package watermark

import (
	"testing"
	"time"
)

func TestExperimentalV4Build89DerivedTimingAccounting(t *testing.T) {
	base := experimentalV4PhoneBuild84RecoveryTelemetry{}
	base.Attempted = true
	base.TotalElapsed = 1000 * time.Millisecond
	base.GeometryElapsed = 700 * time.Millisecond
	base.PlanePrepElapsed = 20 * time.Millisecond
	base.QualificationElapsed = 80 * time.Millisecond
	base.DecodeWallElapsed = 150 * time.Millisecond
	base.Build84GeometryProfile = experimentalV4PhoneBuild84GeometryTelemetry{
		PlanePrepElapsed:   10 * time.Millisecond,
		FreezeElapsed:      100 * time.Millisecond,
		Prefix1WallElapsed: 50 * time.Millisecond,
		Gen2WallElapsed:    60 * time.Millisecond,
		Gen3WallElapsed:    70 * time.Millisecond,
		Gen4WallElapsed:    300 * time.Millisecond,
		Prefix1Workers:     8,
		Gen2Workers:        8,
		Gen3Workers:        8,
		Gen4Workers:        8,
		Gen2Tasks:          12,
		Gen3Tasks:          34,
		Gen4Tasks:          56,
		FetchFoldScores:    90,
		FetchBlockReads:    1000,
		FetchBlockSuccess:  980,
		FetchBlockFailed:   20,
	}

	r := experimentalV4PhoneBuild89RecoveryTelemetry{experimentalV4PhoneBuild84RecoveryTelemetry: base}
	gp := base.Build84GeometryProfile
	r.GeometryAccountedElapsed = gp.PlanePrepElapsed + gp.FreezeElapsed + gp.Prefix1WallElapsed + gp.Gen2WallElapsed + gp.Gen3WallElapsed + gp.Gen4WallElapsed
	r.GeometryUnaccountedElapsed = experimentalV4PhoneBuild89NonNegativeDuration(base.GeometryElapsed - r.GeometryAccountedElapsed)
	r.PostGeometryAccountedElapsed = base.PlanePrepElapsed + base.QualificationElapsed + base.DecodeWallElapsed
	r.TotalAccountedElapsed = base.GeometryElapsed + r.PostGeometryAccountedElapsed
	r.TotalUnaccountedElapsed = experimentalV4PhoneBuild89NonNegativeDuration(base.TotalElapsed - r.TotalAccountedElapsed)

	if r.GeometryAccountedElapsed != 590*time.Millisecond || r.GeometryUnaccountedElapsed != 110*time.Millisecond {
		t.Fatalf("geometry accounting=%s/%s want 590ms/110ms", r.GeometryAccountedElapsed, r.GeometryUnaccountedElapsed)
	}
	if r.PostGeometryAccountedElapsed != 250*time.Millisecond || r.TotalAccountedElapsed != 950*time.Millisecond || r.TotalUnaccountedElapsed != 50*time.Millisecond {
		t.Fatalf("total accounting post=%s accounted=%s unaccounted=%s", r.PostGeometryAccountedElapsed, r.TotalAccountedElapsed, r.TotalUnaccountedElapsed)
	}

	var p ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild89ApplyTelemetry(&p, r)
	if !p.Build89Attempted || !p.Build84Attempted || !p.Build76Attempted {
		t.Fatalf("attempted chain lost: %+v", p)
	}
	if p.Build89TotalMs != 1000 || p.Build89GeometryMs != 700 || p.Build89GeometryAccountedMs != 590 || p.Build89GeometryUnaccountedMs != 110 {
		t.Fatalf("geometry public accounting: %+v", p)
	}
	if p.Build89PostPlanePrepMs != 20 || p.Build89QualificationMs != 80 || p.Build89DecodeWallMs != 150 || p.Build89PostGeometryAccountedMs != 250 || p.Build89TotalAccountedMs != 950 || p.Build89TotalUnaccountedMs != 50 {
		t.Fatalf("post/total public accounting: %+v", p)
	}
	if p.Build84FetchFoldScores != 90 || p.Build84FetchBlockReads != 1000 || p.Build84FetchBlockSuccess != 980 || p.Build84FetchBlockFailed != 20 {
		t.Fatalf("Build84 view changed: %+v", p)
	}
}

func TestExperimentalV4Build89NonNegativeDuration(t *testing.T) {
	if got := experimentalV4PhoneBuild89NonNegativeDuration(-time.Nanosecond); got != 0 {
		t.Fatalf("negative duration clamp=%s want 0", got)
	}
	if got := experimentalV4PhoneBuild89NonNegativeDuration(7 * time.Millisecond); got != 7*time.Millisecond {
		t.Fatalf("positive duration=%s want 7ms", got)
	}
}
