package watermark

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"

	jpeglegacy "github.com/pj/pixseal/internal/jpeglegacy"
)

// Build94 is observability-only. These tests execute the unchanged qualified
// Build93 recovery implementation and derive stage accounting from telemetry
// that already exists in the qualified path. No Build94 recovery algorithm
// exists and none of these measurements participates in recovery decisions.

type experimentalV4Build94DerivedTiming struct {
	geometryAccounted   time.Duration
	geometryUnaccounted time.Duration
	postAccounted       time.Duration
	totalAccounted      time.Duration
	totalUnaccounted    time.Duration
}

func experimentalV4Build94NonNegativeDuration(v time.Duration) time.Duration {
	if v < 0 {
		return 0
	}
	return v
}

func experimentalV4Build94DeriveTiming(telemetry experimentalV4PhoneBuild93RecoveryTelemetry) experimentalV4Build94DerivedTiming {
	gp := telemetry.Build84GeometryProfile
	geometryAccounted := gp.PlanePrepElapsed + gp.FreezeElapsed + gp.Prefix1WallElapsed + gp.Gen2WallElapsed + gp.Gen3WallElapsed + gp.Gen4WallElapsed
	postAccounted := telemetry.PlanePrepElapsed + telemetry.QualificationElapsed + telemetry.DecodeWallElapsed
	totalAccounted := telemetry.GeometryElapsed + postAccounted
	return experimentalV4Build94DerivedTiming{
		geometryAccounted:   geometryAccounted,
		geometryUnaccounted: experimentalV4Build94NonNegativeDuration(telemetry.GeometryElapsed - geometryAccounted),
		postAccounted:       postAccounted,
		totalAccounted:      totalAccounted,
		totalUnaccounted:    experimentalV4Build94NonNegativeDuration(telemetry.TotalElapsed - totalAccounted),
	}
}

func experimentalV4Build94Config(t *testing.T) (string, []byte, int, int, string, time.Duration) {
	t.Helper()
	acquisitionDir := os.Getenv("PIXSEAL_BUILD94_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build94 private retained profile requires PIXSEAL_BUILD94_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD94_KEY"))
	if len(key) < 8 {
		t.Fatal("Build94 PIXSEAL_BUILD94_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD94_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD94_WIDTH: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD94_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD94_HEIGHT: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD94_MESSAGE_B")
	if expectedPayload == "" {
		t.Fatal("PIXSEAL_BUILD94_MESSAGE_B is required")
	}
	cooldownMs := 5000
	if raw := os.Getenv("PIXSEAL_BUILD94_COOLDOWN_MS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 || v > 60000 {
			t.Fatalf("invalid PIXSEAL_BUILD94_COOLDOWN_MS=%q", raw)
		}
		cooldownMs = v
	}
	return acquisitionDir, key, cw, ch, expectedPayload, time.Duration(cooldownMs) * time.Millisecond
}

func experimentalV4Build94Prepare(t *testing.T, acquisitionDir, imageName string, cw, ch int) (image.Image, PrintBoundaryEstimate, time.Duration) {
	t.Helper()
	path := filepath.Join(acquisitionDir, imageName)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	src, err := jpeglegacy.Decode(f)
	_ = f.Close()
	if err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	started := time.Now()
	work, boundary, hyp, _, _, _ := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
	elapsed := time.Since(started)
	if len(hyp) >= experimentalV4PhoneEnsembleSize {
		t.Fatalf("Build94 %s unexpectedly accepted before deep recovery", imageName)
	}
	return work, boundary, elapsed
}

func experimentalV4Build94AssertExpected(t *testing.T, imageName string, payload []byte, telemetry experimentalV4PhoneBuild93RecoveryTelemetry, err error, expectedPayload string) {
	t.Helper()
	exp, ok := experimentalV4Build91ExpectedForImage(imageName)
	if !ok {
		t.Fatalf("Build94 has no frozen expected telemetry for %s", imageName)
	}
	gp := telemetry.Build84GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	if telemetry.GeometryEvaluations != exp.geometryEvals || telemetry.BankCandidates != exp.bank || telemetry.QualificationEvaluations != exp.qualEvals || telemetry.QualifiedCandidates != exp.qualified || telemetry.DecodeCandidatesTried != exp.decode || telemetry.ListFramesTried != exp.frames || physical != exp.physical || physical-telemetry.DecodeCandidatesTried != exp.speculative || gp.FetchFoldScores != exp.foldScores || gp.FetchBlockReads != exp.blockReads || gp.FetchBlockSuccess != exp.blockSuccess || gp.FetchBlockFailed != exp.blockFailed {
		t.Fatalf("Build94 %s qualified Build93 semantic telemetry mismatch: %+v", imageName, telemetry)
	}
	wantQualificationWorkers := experimentalV4PhoneBuild90QualificationWorkerCount(exp.bank)
	if telemetry.QualificationWorkers != wantQualificationWorkers || telemetry.QualificationTasks != exp.bank {
		t.Fatalf("Build94 %s qualification scheduling mismatch: workers=%d want=%d tasks=%d want=%d", imageName, telemetry.QualificationWorkers, wantQualificationWorkers, telemetry.QualificationTasks, exp.bank)
	}
	wantDecodeWorkers := experimentalV4PhoneBuild66DecodeWorkerCount(exp.qualified)
	if telemetry.DecodeWorkers != wantDecodeWorkers {
		t.Fatalf("Build94 %s decode workers=%d want=%d", imageName, telemetry.DecodeWorkers, wantDecodeWorkers)
	}
	if exp.accept {
		if err != nil || !telemetry.Authenticated || string(payload) != expectedPayload {
			t.Fatalf("Build94 %s accept mismatch: err=%v authenticated=%t payload=%q", imageName, err, telemetry.Authenticated, string(payload))
		}
	} else if err == nil || telemetry.Authenticated || len(payload) != 0 {
		t.Fatalf("Build94 %s reject mismatch: err=%v authenticated=%t payload-bytes=%d", imageName, err, telemetry.Authenticated, len(payload))
	}

	if telemetry.ProtectedSamplingPasses != telemetry.PhysicalDecodeCandidates {
		t.Fatalf("Build94 %s protected sampling passes=%d physical decode=%d", imageName, telemetry.ProtectedSamplingPasses, telemetry.PhysicalDecodeCandidates)
	}
	if telemetry.ProtectedProfileDerivations != telemetry.PhysicalProfilesTried {
		t.Fatalf("Build94 %s profile derivations=%d physical profiles=%d", imageName, telemetry.ProtectedProfileDerivations, telemetry.PhysicalProfilesTried)
	}
	switch imageName {
	case "phone-b-mild.jpg":
		if telemetry.PhysicalProfilesTried != 2086 || telemetry.PhysicalListFramesTried != 2135407 || telemetry.ProtectedSamplingPasses != 696 || telemetry.ProtectedProfileDerivations != 2086 {
			t.Fatalf("Build94 B/mild physical decode telemetry mismatch: physical=%d profiles=%d frames=%d sampling=%d derivations=%d", telemetry.PhysicalDecodeCandidates, telemetry.PhysicalProfilesTried, telemetry.PhysicalListFramesTried, telemetry.ProtectedSamplingPasses, telemetry.ProtectedProfileDerivations)
		}
	case "phone-b-angle.jpg":
		if telemetry.PhysicalProfilesTried != 0 || telemetry.PhysicalListFramesTried != 0 || telemetry.ProtectedSamplingPasses != 0 || telemetry.ProtectedProfileDerivations != 0 {
			t.Fatalf("Build94 B/angle unexpectedly reached protected decode: physical=%d profiles=%d frames=%d sampling=%d derivations=%d", telemetry.PhysicalDecodeCandidates, telemetry.PhysicalProfilesTried, telemetry.PhysicalListFramesTried, telemetry.ProtectedSamplingPasses, telemetry.ProtectedProfileDerivations)
		}
	}
}

var experimentalV4Build94RowColumns = []string{
	"image", "outcome", "prepare_ms", "wall_ms", "total_ms", "geometry_ms",
	"geometry_plane_ms", "freeze_ms", "prefix1_ms", "gen2_ms", "gen3_ms", "gen4_ms",
	"prefix1_worker_ms", "gen2_worker_ms", "gen3_worker_ms", "gen4_worker_ms",
	"geometry_accounted_ms", "geometry_unaccounted_ms", "post_plane_ms", "qualification_ms", "decode_ms", "post_accounted_ms", "total_accounted_ms", "total_unaccounted_ms",
	"geometry_evals", "bank", "qualification_evals", "qualified", "decode_candidates", "list_frames", "physical_decode", "speculative_decode",
	"physical_profiles", "physical_list_frames", "protected_sampling_passes", "protected_profile_derivations", "sampling_worker_ms", "list_worker_ms", "authenticated",
	"prefix1_workers", "gen2_workers", "gen3_workers", "gen4_workers", "gen2_tasks", "gen3_tasks", "gen4_tasks", "qualification_workers", "qualification_tasks", "decode_workers",
	"fold_scores", "block_reads", "block_success", "block_failed",
}

func experimentalV4Build94RowFields(imageName string, prep, wall time.Duration, telemetry experimentalV4PhoneBuild93RecoveryTelemetry, err error) []string {
	gp := telemetry.Build84GeometryProfile
	derived := experimentalV4Build94DeriveTiming(telemetry)
	physical := telemetry.PhysicalDecodeCandidates
	outcome := "reject"
	if err == nil {
		outcome = "accept"
	}
	ms := func(v time.Duration) string { return strconv.FormatInt(v.Milliseconds(), 10) }
	i := func(v int) string { return strconv.Itoa(v) }
	return []string{
		imageName, outcome, ms(prep), ms(wall), ms(telemetry.TotalElapsed), ms(telemetry.GeometryElapsed),
		ms(gp.PlanePrepElapsed), ms(gp.FreezeElapsed), ms(gp.Prefix1WallElapsed), ms(gp.Gen2WallElapsed), ms(gp.Gen3WallElapsed), ms(gp.Gen4WallElapsed),
		ms(gp.Prefix1WorkerElapsed), ms(gp.Gen2WorkerElapsed), ms(gp.Gen3WorkerElapsed), ms(gp.Gen4WorkerElapsed),
		ms(derived.geometryAccounted), ms(derived.geometryUnaccounted), ms(telemetry.PlanePrepElapsed), ms(telemetry.QualificationElapsed), ms(telemetry.DecodeWallElapsed), ms(derived.postAccounted), ms(derived.totalAccounted), ms(derived.totalUnaccounted),
		i(telemetry.GeometryEvaluations), i(telemetry.BankCandidates), i(telemetry.QualificationEvaluations), i(telemetry.QualifiedCandidates), i(telemetry.DecodeCandidatesTried), i(telemetry.ListFramesTried), i(physical), i(physical - telemetry.DecodeCandidatesTried),
		i(telemetry.PhysicalProfilesTried), i(telemetry.PhysicalListFramesTried), i(telemetry.ProtectedSamplingPasses), i(telemetry.ProtectedProfileDerivations), ms(telemetry.SamplingWorkerElapsed), ms(telemetry.ListWorkerElapsed), strconv.FormatBool(telemetry.Authenticated),
		i(gp.Prefix1Workers), i(gp.Gen2Workers), i(gp.Gen3Workers), i(gp.Gen4Workers), i(gp.Gen2Tasks), i(gp.Gen3Tasks), i(gp.Gen4Tasks), i(telemetry.QualificationWorkers), i(telemetry.QualificationTasks), i(telemetry.DecodeWorkers),
		i(gp.FetchFoldScores), i(gp.FetchBlockReads), i(gp.FetchBlockSuccess), i(gp.FetchBlockFailed),
	}
}

func experimentalV4Build94PrintRow(imageName string, prep, wall time.Duration, telemetry experimentalV4PhoneBuild93RecoveryTelemetry, err error) {
	fmt.Printf("build94-row\t%s\n", strings.Join(experimentalV4Build94RowFields(imageName, prep, wall, telemetry, err), "\t"))
}

func TestExperimentalV4Build94DerivedTimingAccounting(t *testing.T) {
	telemetry := experimentalV4PhoneBuild93RecoveryTelemetry{}
	telemetry.TotalElapsed = 1000 * time.Millisecond
	telemetry.GeometryElapsed = 700 * time.Millisecond
	telemetry.PlanePrepElapsed = 20 * time.Millisecond
	telemetry.QualificationElapsed = 80 * time.Millisecond
	telemetry.DecodeWallElapsed = 150 * time.Millisecond
	telemetry.Build84GeometryProfile = experimentalV4PhoneBuild84GeometryTelemetry{
		PlanePrepElapsed: 10 * time.Millisecond, FreezeElapsed: 100 * time.Millisecond,
		Prefix1WallElapsed: 50 * time.Millisecond, Gen2WallElapsed: 60 * time.Millisecond,
		Gen3WallElapsed: 70 * time.Millisecond, Gen4WallElapsed: 300 * time.Millisecond,
	}
	got := experimentalV4Build94DeriveTiming(telemetry)
	if got.geometryAccounted != 590*time.Millisecond || got.geometryUnaccounted != 110*time.Millisecond {
		t.Fatalf("geometry accounting=%s/%s want 590ms/110ms", got.geometryAccounted, got.geometryUnaccounted)
	}
	if got.postAccounted != 250*time.Millisecond || got.totalAccounted != 950*time.Millisecond || got.totalUnaccounted != 50*time.Millisecond {
		t.Fatalf("total accounting post=%s accounted=%s unaccounted=%s", got.postAccounted, got.totalAccounted, got.totalUnaccounted)
	}
}

func TestExperimentalV4Build94RowShape(t *testing.T) {
	fields := experimentalV4Build94RowFields("phone-b-mild.jpg", 2*time.Millisecond, 3*time.Millisecond, experimentalV4PhoneBuild93RecoveryTelemetry{}, nil)
	if len(fields) != len(experimentalV4Build94RowColumns) {
		t.Fatalf("Build94 row fields=%d columns=%d", len(fields), len(experimentalV4Build94RowColumns))
	}
	seen := make(map[string]struct{}, len(experimentalV4Build94RowColumns))
	for _, column := range experimentalV4Build94RowColumns {
		if column == "" || strings.ContainsAny(column, "\t\r\n") {
			t.Fatalf("Build94 invalid row column %q", column)
		}
		if _, duplicate := seen[column]; duplicate {
			t.Fatalf("Build94 duplicate row column %q", column)
		}
		seen[column] = struct{}{}
	}
	for _, field := range fields {
		if strings.ContainsAny(field, "\t\r\n") {
			t.Fatalf("Build94 invalid row field %q", field)
		}
	}
}

func TestExperimentalV4Build94ExpectedTelemetryTable(t *testing.T) {
	mild, ok := experimentalV4Build91ExpectedForImage("phone-b-mild.jpg")
	if !ok || mild.bank != 937 || mild.qualified != 935 || mild.decode != 691 || mild.frames != 2120047 || mild.physical != 696 || !mild.accept {
		t.Fatalf("Build94 B/mild expected table invalid: %+v", mild)
	}
	angle, ok := experimentalV4Build91ExpectedForImage("phone-b-angle.jpg")
	if !ok || angle.bank != 6198 || angle.qualified != 0 || angle.physical != 0 || angle.accept {
		t.Fatalf("Build94 B/angle expected table invalid: %+v", angle)
	}
}

func TestExperimentalV4Build94RetainedTimingProfile(t *testing.T) {
	acquisitionDir, key, cw, ch, expectedPayload, cooldown := experimentalV4Build94Config(t)
	fmt.Printf("build94-header\t%s\n", strings.Join(experimentalV4Build94RowColumns, "\t"))
	for _, imageName := range []string{"phone-b-mild.jpg", "phone-b-angle.jpg"} {
		work, boundary, prep := experimentalV4Build94Prepare(t, acquisitionDir, imageName, cw, ch)
		fmt.Printf("build94-prepare\t%s\t%d\n", imageName, prep.Milliseconds())
		if cooldown > 0 {
			time.Sleep(cooldown)
		}
		started := time.Now()
		payload, _, telemetry, err := experimentalV4PhoneBuild93Recover(work, boundary, key, cw, ch)
		wall := time.Since(started)
		experimentalV4Build94AssertExpected(t, imageName, payload, telemetry, err, expectedPayload)
		experimentalV4Build94PrintRow(imageName, prep, wall, telemetry, err)
	}
}

func TestExperimentalV4Build94CPUProfile(t *testing.T) {
	acquisitionDir, key, cw, ch, expectedPayload, cooldown := experimentalV4Build94Config(t)
	imageName := os.Getenv("PIXSEAL_BUILD94_CASE")
	if imageName != "phone-b-mild.jpg" && imageName != "phone-b-angle.jpg" {
		t.Fatalf("PIXSEAL_BUILD94_CASE must be phone-b-mild.jpg or phone-b-angle.jpg, got %q", imageName)
	}
	profilePath := os.Getenv("PIXSEAL_BUILD94_CPU_PROFILE")
	if profilePath == "" {
		t.Fatal("PIXSEAL_BUILD94_CPU_PROFILE is required")
	}
	work, boundary, prep := experimentalV4Build94Prepare(t, acquisitionDir, imageName, cw, ch)
	if cooldown > 0 {
		time.Sleep(cooldown)
	}
	f, err := os.Create(profilePath)
	if err != nil {
		t.Fatalf("create CPU profile %s: %v", profilePath, err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		_ = f.Close()
		t.Fatalf("start CPU profile: %v", err)
	}
	started := time.Now()
	payload, _, telemetry, recoverErr := experimentalV4PhoneBuild93Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		t.Fatalf("close CPU profile: %v", err)
	}
	experimentalV4Build94AssertExpected(t, imageName, payload, telemetry, recoverErr, expectedPayload)
	fmt.Printf("build94-cpu-row\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\n",
		imageName, prep.Milliseconds(), wall.Milliseconds(), telemetry.TotalElapsed.Milliseconds(), telemetry.GeometryElapsed.Milliseconds(), telemetry.QualificationElapsed.Milliseconds(), telemetry.DecodeWallElapsed.Milliseconds(), telemetry.SamplingWorkerElapsed.Milliseconds(), telemetry.ListWorkerElapsed.Milliseconds(), telemetry.ProtectedSamplingPasses, telemetry.ProtectedProfileDerivations, profilePath)
}
