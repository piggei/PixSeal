package watermark

import (
	"image"
	"testing"
	"time"
)

func TestExperimentalV4Build69ProductionSemanticsRemainFrozen(t *testing.T) {
	if experimentalV4PhoneProposalFloor != 0.16 || experimentalV4PhoneValidationFloor != 0.10 || experimentalV4PhonePilotScoreFloor != 0.12 || experimentalV4PhonePilotMarginFloor != 0.025 {
		t.Fatalf("Build69 changed qualified production thresholds")
	}
	if experimentalV4PhoneBuild63SeedsPerPair != 4 {
		t.Fatalf("Build69 inherited seeds-per-pair=%d want 4", experimentalV4PhoneBuild63SeedsPerPair)
	}
}

func TestExperimentalV4Build69BlindBankRejectsInvalidBoundary(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	bank, evals, seeds, workers, profile := experimentalV4PhoneBuild69BlindBank(img, PrintBoundaryEstimate{}, 1632, 1632)
	if len(bank) != 0 || evals != 0 || seeds != 0 || workers != 0 {
		t.Fatalf("invalid boundary returned bank=%d evals=%d seeds=%d workers=%d", len(bank), evals, seeds, workers)
	}
	if profile != (experimentalV4PhoneBuild69GeometryTelemetry{}) {
		t.Fatalf("invalid boundary produced profile=%+v", profile)
	}
}

func TestExperimentalV4Build69ApplyTelemetryPreservesBuild68View(t *testing.T) {
	recovery := experimentalV4PhoneBuild69RecoveryTelemetry{}
	recovery.Attempted = true
	recovery.SeedsSelected = 24
	recovery.GeometryEvaluations = 79259
	recovery.BankCandidates = 937
	recovery.QualifiedCandidates = 935
	recovery.DecodeCandidatesTried = 691
	recovery.ListFramesTried = 2120047
	recovery.Authenticated = true
	recovery.Workers = 8
	recovery.DecodeWorkers = 8
	recovery.PhysicalDecodeCandidates = 696
	recovery.TotalElapsed = 100 * time.Second
	recovery.GeometryElapsed = 80 * time.Second
	recovery.PlanePrepElapsed = 250 * time.Millisecond
	recovery.QualificationElapsed = 10 * time.Second
	recovery.DecodeWallElapsed = 9 * time.Second
	recovery.GeometryPlanePrepElapsed = 200 * time.Millisecond
	recovery.GeometryFreezeElapsed = 3 * time.Second
	recovery.GeometrySeedWallElapsed = 76 * time.Second
	recovery.GeometrySeedWorkerElapsed = 590 * time.Second
	recovery.GeometrySeedMinElapsed = 10 * time.Second
	recovery.GeometrySeedMedianElapsed = 20 * time.Second
	recovery.GeometrySeedMaxElapsed = 40 * time.Second
	recovery.GeometrySeedMinEvaluations = 100
	recovery.GeometrySeedMaxEvaluations = 10000
	recovery.GeometrySeedMinBank = 0
	recovery.GeometrySeedMaxBank = 200

	var public ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild69ApplyTelemetry(&public, recovery)
	if !public.Build69Attempted || !public.Build68Attempted || !public.Build66Attempted || !public.Build65Attempted || !public.Build64Attempted {
		t.Fatalf("Build69 did not preserve inherited attempted telemetry")
	}
	if public.Build64SeedsSelected != 24 || public.Build64GeometryEvaluations != 79259 || public.Build64BankCandidates != 937 || public.Build64QualifiedCandidates != 935 || public.Build64DecodeCandidates != 691 || public.Build64ListFramesTried != 2120047 || !public.Build64Authenticated {
		t.Fatalf("qualified logical telemetry changed: %+v", public)
	}
	if public.Build68SpeculativeDecodeCandidates != 5 {
		t.Fatalf("Build68 physical/logical split changed: speculative=%d want 5", public.Build68SpeculativeDecodeCandidates)
	}
	if public.Build69GeometryPlanePrepMs != 200 || public.Build69GeometryFreezeMs != 3000 || public.Build69GeometrySeedWallMs != 76000 || public.Build69GeometrySeedWorkerMs != 590000 {
		t.Fatalf("Build69 timing=%d/%d/%d/%d", public.Build69GeometryPlanePrepMs, public.Build69GeometryFreezeMs, public.Build69GeometrySeedWallMs, public.Build69GeometrySeedWorkerMs)
	}
	if public.Build69GeometrySeedMinMs != 10000 || public.Build69GeometrySeedMedianMs != 20000 || public.Build69GeometrySeedMaxMs != 40000 || public.Build69GeometrySeedMinEvals != 100 || public.Build69GeometrySeedMaxEvals != 10000 || public.Build69GeometrySeedMinBank != 0 || public.Build69GeometrySeedMaxBank != 200 {
		t.Fatalf("Build69 seed profile mismatch")
	}
}
