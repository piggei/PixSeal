package watermark

import (
	"image"
	"time"
)

// experimentalV4PhoneBuild89RecoveryTelemetry is observability-only over the
// qualified Build84 runtime. It embeds the complete Build84 telemetry and derives
// coarse wall-clock accounting after the qualified recovery returns. No new
// timing probes are inserted inside the hot geometry, qualification or protected
// decode loops.
type experimentalV4PhoneBuild89RecoveryTelemetry struct {
	experimentalV4PhoneBuild84RecoveryTelemetry

	GeometryAccountedElapsed     time.Duration
	GeometryUnaccountedElapsed   time.Duration
	PostGeometryAccountedElapsed time.Duration
	TotalAccountedElapsed        time.Duration
	TotalUnaccountedElapsed      time.Duration
}

func experimentalV4PhoneBuild89NonNegativeDuration(v time.Duration) time.Duration {
	if v < 0 {
		return 0
	}
	return v
}

// experimentalV4PhoneBuild89Recover delegates all recovery work to the current
// qualified Build84 implementation and computes only derived timing accounting
// from telemetry Build84 already measured. The returned payload, error and
// semantic telemetry are therefore authoritative Build84 results.
func experimentalV4PhoneBuild89Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild89RecoveryTelemetry, error) {
	payload, info, base, err := experimentalV4PhoneBuild84Recover(work, boundary, key, cw, ch)
	telemetry := experimentalV4PhoneBuild89RecoveryTelemetry{
		experimentalV4PhoneBuild84RecoveryTelemetry: base,
	}

	gp := base.Build84GeometryProfile
	telemetry.GeometryAccountedElapsed = gp.PlanePrepElapsed +
		gp.FreezeElapsed +
		gp.Prefix1WallElapsed +
		gp.Gen2WallElapsed +
		gp.Gen3WallElapsed +
		gp.Gen4WallElapsed
	telemetry.GeometryUnaccountedElapsed = experimentalV4PhoneBuild89NonNegativeDuration(base.GeometryElapsed - telemetry.GeometryAccountedElapsed)

	telemetry.PostGeometryAccountedElapsed = base.PlanePrepElapsed + base.QualificationElapsed + base.DecodeWallElapsed
	telemetry.TotalAccountedElapsed = base.GeometryElapsed + telemetry.PostGeometryAccountedElapsed
	telemetry.TotalUnaccountedElapsed = experimentalV4PhoneBuild89NonNegativeDuration(base.TotalElapsed - telemetry.TotalAccountedElapsed)

	return payload, info, telemetry, err
}

func experimentalV4PhoneBuild89ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild89RecoveryTelemetry) {
	experimentalV4PhoneBuild84ApplyTelemetry(public, recovery.experimentalV4PhoneBuild84RecoveryTelemetry)
	base := recovery.experimentalV4PhoneBuild84RecoveryTelemetry
	gp := base.Build84GeometryProfile

	public.Build89Attempted = base.Attempted
	public.Build89TotalMs = base.TotalElapsed.Milliseconds()
	public.Build89GeometryMs = base.GeometryElapsed.Milliseconds()
	public.Build89GeometryPlanePrepMs = gp.PlanePrepElapsed.Milliseconds()
	public.Build89FreezeMs = gp.FreezeElapsed.Milliseconds()
	public.Build89Prefix1Ms = gp.Prefix1WallElapsed.Milliseconds()
	public.Build89Gen2Ms = gp.Gen2WallElapsed.Milliseconds()
	public.Build89Gen3Ms = gp.Gen3WallElapsed.Milliseconds()
	public.Build89Gen4Ms = gp.Gen4WallElapsed.Milliseconds()
	public.Build89GeometryAccountedMs = recovery.GeometryAccountedElapsed.Milliseconds()
	public.Build89GeometryUnaccountedMs = recovery.GeometryUnaccountedElapsed.Milliseconds()
	public.Build89PostPlanePrepMs = base.PlanePrepElapsed.Milliseconds()
	public.Build89QualificationMs = base.QualificationElapsed.Milliseconds()
	public.Build89DecodeWallMs = base.DecodeWallElapsed.Milliseconds()
	public.Build89PostGeometryAccountedMs = recovery.PostGeometryAccountedElapsed.Milliseconds()
	public.Build89TotalAccountedMs = recovery.TotalAccountedElapsed.Milliseconds()
	public.Build89TotalUnaccountedMs = recovery.TotalUnaccountedElapsed.Milliseconds()
}
