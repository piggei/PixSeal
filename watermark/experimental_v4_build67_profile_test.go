package watermark

import (
	"errors"
	"testing"
	"time"
)

func TestExperimentalV4Build67ProductionThresholdsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build67 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build67 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build67PhysicalProfilingKeepsLogicalTelemetrySeparate(t *testing.T) {
	telemetry := experimentalV4PhoneBuild67RecoveryTelemetry{}
	results := []experimentalV4PhoneBuild67DecodeResult{
		{
			experimentalV4PhoneBuild66DecodeResult: experimentalV4PhoneBuild66DecodeResult{profiles: 2, frames: 100, maxConfidence: 1.5, err: errors.New("no auth")},
			samplingElapsed:                        11 * time.Millisecond,
			listElapsed:                            21 * time.Millisecond,
		},
		{
			experimentalV4PhoneBuild66DecodeResult: experimentalV4PhoneBuild66DecodeResult{payload: []byte("winner"), info: ExperimentalV4ExtractInfo{Profile: ProfileRobust}, profiles: 3, frames: 200, maxConfidence: 2.5, err: nil},
			samplingElapsed:                        12 * time.Millisecond,
			listElapsed:                            22 * time.Millisecond,
		},
		{
			experimentalV4PhoneBuild66DecodeResult: experimentalV4PhoneBuild66DecodeResult{payload: []byte("speculative"), info: ExperimentalV4ExtractInfo{Profile: ProfileCapacity}, profiles: 9, frames: 9999, maxConfidence: 99, err: nil},
			samplingElapsed:                        13 * time.Millisecond,
			listElapsed:                            23 * time.Millisecond,
		},
	}

	experimentalV4PhoneBuild67AccumulatePhysical(&telemetry, results)
	if telemetry.PhysicalDecodeCandidates != 3 || telemetry.PhysicalProfilesTried != 14 || telemetry.PhysicalListFramesTried != 10299 {
		t.Fatalf("physical telemetry=%d/%d/%d want 3/14/10299", telemetry.PhysicalDecodeCandidates, telemetry.PhysicalProfilesTried, telemetry.PhysicalListFramesTried)
	}
	if telemetry.SamplingWorkerElapsed != 36*time.Millisecond || telemetry.ListWorkerElapsed != 66*time.Millisecond {
		t.Fatalf("worker timing=%s/%s want 36ms/66ms", telemetry.SamplingWorkerElapsed, telemetry.ListWorkerElapsed)
	}

	semantic := make([]experimentalV4PhoneBuild66DecodeResult, len(results))
	for i := range results {
		semantic[i] = results[i].experimentalV4PhoneBuild66DecodeResult
	}
	payload, info, ok := experimentalV4PhoneBuild66ConsumeBatch(semantic, &telemetry.experimentalV4PhoneBuild66RecoveryTelemetry)
	if !ok || string(payload) != "winner" || info.Profile != ProfileRobust {
		t.Fatalf("unexpected logical winner ok=%t payload=%q profile=%q", ok, payload, info.Profile)
	}
	if telemetry.DecodeCandidatesTried != 2 || telemetry.ProfilesTried != 5 || telemetry.ListFramesTried != 300 {
		t.Fatalf("logical telemetry=%d/%d/%d want 2/5/300", telemetry.DecodeCandidatesTried, telemetry.ProfilesTried, telemetry.ListFramesTried)
	}
	if telemetry.MaxDataConfidence != 2.5 {
		t.Fatalf("logical max confidence=%f want 2.5", telemetry.MaxDataConfidence)
	}
}

func TestExperimentalV4Build67ApplyTelemetryReportsSpeculativePhysicalWork(t *testing.T) {
	recovery := experimentalV4PhoneBuild67RecoveryTelemetry{}
	recovery.Attempted = true
	recovery.DecodeWorkers = 8
	recovery.DecodeCandidatesTried = 691
	recovery.PhysicalDecodeCandidates = 696
	recovery.PhysicalProfilesTried = 1300
	recovery.PhysicalListFramesTried = 2140000
	recovery.TotalElapsed = 10 * time.Second
	recovery.GeometryElapsed = 2 * time.Second
	recovery.PlanePrepElapsed = 750 * time.Millisecond
	recovery.QualificationElapsed = 3 * time.Second
	recovery.DecodeWallElapsed = 5 * time.Second
	recovery.SamplingWorkerElapsed = 7 * time.Second
	recovery.ListWorkerElapsed = 19 * time.Second

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild67ApplyTelemetry(&public, recovery)
	if !public.Build67Attempted || !public.Build66Attempted {
		t.Fatalf("profiling did not preserve inherited Build66 attempted telemetry")
	}
	if public.Build67SpeculativeDecodeCandidates != 5 {
		t.Fatalf("speculative=%d want 5", public.Build67SpeculativeDecodeCandidates)
	}
	if public.Build67TotalMs != 10000 || public.Build67GeometryMs != 2000 || public.Build67PlanePrepMs != 750 || public.Build67QualificationMs != 3000 || public.Build67DecodeWallMs != 5000 {
		t.Fatalf("wall timing=%d/%d/%d/%d/%d", public.Build67TotalMs, public.Build67GeometryMs, public.Build67PlanePrepMs, public.Build67QualificationMs, public.Build67DecodeWallMs)
	}
	if public.Build67SamplingWorkerMs != 7000 || public.Build67ListWorkerMs != 19000 {
		t.Fatalf("worker timing=%d/%d", public.Build67SamplingWorkerMs, public.Build67ListWorkerMs)
	}
}
