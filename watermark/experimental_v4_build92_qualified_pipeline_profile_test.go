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

type experimentalV4Build92DerivedTiming struct {
	geometryAccounted   time.Duration
	geometryUnaccounted time.Duration
	postAccounted       time.Duration
	totalAccounted      time.Duration
	totalUnaccounted    time.Duration
}

func experimentalV4Build92NonNegativeDuration(v time.Duration) time.Duration {
	if v < 0 {
		return 0
	}
	return v
}

func experimentalV4Build92DeriveTiming(telemetry experimentalV4PhoneBuild90RecoveryTelemetry) experimentalV4Build92DerivedTiming {
	gp := telemetry.Build84GeometryProfile
	geometryAccounted := gp.PlanePrepElapsed + gp.FreezeElapsed + gp.Prefix1WallElapsed + gp.Gen2WallElapsed + gp.Gen3WallElapsed + gp.Gen4WallElapsed
	postAccounted := telemetry.PlanePrepElapsed + telemetry.QualificationElapsed + telemetry.DecodeWallElapsed
	totalAccounted := telemetry.GeometryElapsed + postAccounted
	return experimentalV4Build92DerivedTiming{
		geometryAccounted:   geometryAccounted,
		geometryUnaccounted: experimentalV4Build92NonNegativeDuration(telemetry.GeometryElapsed - geometryAccounted),
		postAccounted:       postAccounted,
		totalAccounted:      totalAccounted,
		totalUnaccounted:    experimentalV4Build92NonNegativeDuration(telemetry.TotalElapsed - totalAccounted),
	}
}

func experimentalV4Build92Config(t *testing.T) (string, []byte, int, int, string, time.Duration) {
	t.Helper()
	acquisitionDir := os.Getenv("PIXSEAL_BUILD92_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build92 private retained profile requires PIXSEAL_BUILD92_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD92_KEY"))
	if len(key) < 8 {
		t.Fatal("Build92 PIXSEAL_BUILD92_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD92_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD92_WIDTH: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD92_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD92_HEIGHT: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD92_MESSAGE_B")
	if expectedPayload == "" {
		t.Fatal("PIXSEAL_BUILD92_MESSAGE_B is required")
	}
	cooldownMs := 5000
	if raw := os.Getenv("PIXSEAL_BUILD92_COOLDOWN_MS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 || v > 60000 {
			t.Fatalf("invalid PIXSEAL_BUILD92_COOLDOWN_MS=%q", raw)
		}
		cooldownMs = v
	}
	return acquisitionDir, key, cw, ch, expectedPayload, time.Duration(cooldownMs) * time.Millisecond
}

func experimentalV4Build92Prepare(t *testing.T, acquisitionDir, imageName string, cw, ch int) (image.Image, PrintBoundaryEstimate, time.Duration) {
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
		t.Fatalf("Build92 %s unexpectedly accepted before deep recovery", imageName)
	}
	return work, boundary, elapsed
}

func experimentalV4Build92AsBuild91Run(wall time.Duration, payload []byte, info ExperimentalV4ExtractInfo, telemetry experimentalV4PhoneBuild90RecoveryTelemetry, err error) experimentalV4Build91Run {
	gp := telemetry.Build84GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	outcome := "reject"
	if err == nil {
		outcome = "accept"
	}
	return experimentalV4Build91Run{
		variant: "build90", outcome: outcome, wall: wall,
		total: telemetry.TotalElapsed, geometry: telemetry.GeometryElapsed, plane: telemetry.PlanePrepElapsed,
		freeze: gp.FreezeElapsed, prefix1: gp.Prefix1WallElapsed, gen2: gp.Gen2WallElapsed, gen3: gp.Gen3WallElapsed, gen4: gp.Gen4WallElapsed,
		qualification: telemetry.QualificationElapsed, decodeWall: telemetry.DecodeWallElapsed,
		geometryEvals: telemetry.GeometryEvaluations, bank: telemetry.BankCandidates, qualEvals: telemetry.QualificationEvaluations,
		qualified: telemetry.QualifiedCandidates, decodeCandidates: telemetry.DecodeCandidatesTried, listFrames: telemetry.ListFramesTried,
		physicalDecode: physical, speculativeDecode: physical - telemetry.DecodeCandidatesTried, authenticated: telemetry.Authenticated,
		qualificationWorkers: telemetry.QualificationWorkers, qualificationTasks: telemetry.QualificationTasks,
		foldScores: gp.FetchFoldScores, blockReads: gp.FetchBlockReads, blockSuccess: gp.FetchBlockSuccess, blockFailed: gp.FetchBlockFailed,
		payload: payload, info: info, err: err,
	}
}

func experimentalV4Build92AssertExpected(t *testing.T, imageName string, wall time.Duration, payload []byte, info ExperimentalV4ExtractInfo, telemetry experimentalV4PhoneBuild90RecoveryTelemetry, err error, expectedPayload string) {
	t.Helper()
	r := experimentalV4Build92AsBuild91Run(wall, payload, info, telemetry, err)
	experimentalV4Build91AssertExpected(t, imageName, r, expectedPayload)
}

func experimentalV4Build92PrintRow(imageName string, prep, wall time.Duration, telemetry experimentalV4PhoneBuild90RecoveryTelemetry, err error) {
	gp := telemetry.Build84GeometryProfile
	derived := experimentalV4Build92DeriveTiming(telemetry)
	physical := telemetry.PhysicalDecodeCandidates
	outcome := "reject"
	if err == nil {
		outcome = "accept"
	}
	ms := func(v time.Duration) string { return strconv.FormatInt(v.Milliseconds(), 10) }
	i := func(v int) string { return strconv.Itoa(v) }
	fields := []string{
		imageName, outcome, ms(prep), ms(wall), ms(telemetry.TotalElapsed), ms(telemetry.GeometryElapsed),
		ms(gp.PlanePrepElapsed), ms(gp.FreezeElapsed), ms(gp.Prefix1WallElapsed), ms(gp.Gen2WallElapsed), ms(gp.Gen3WallElapsed), ms(gp.Gen4WallElapsed),
		ms(gp.Prefix1WorkerElapsed), ms(gp.Gen2WorkerElapsed), ms(gp.Gen3WorkerElapsed), ms(gp.Gen4WorkerElapsed),
		ms(derived.geometryAccounted), ms(derived.geometryUnaccounted), ms(telemetry.PlanePrepElapsed), ms(telemetry.QualificationElapsed), ms(telemetry.DecodeWallElapsed), ms(derived.postAccounted), ms(derived.totalAccounted), ms(derived.totalUnaccounted),
		i(telemetry.GeometryEvaluations), i(telemetry.BankCandidates), i(telemetry.QualificationEvaluations), i(telemetry.QualifiedCandidates), i(telemetry.DecodeCandidatesTried), i(telemetry.ListFramesTried), i(physical), i(physical - telemetry.DecodeCandidatesTried), strconv.FormatBool(telemetry.Authenticated),
		i(gp.Prefix1Workers), i(gp.Gen2Workers), i(gp.Gen3Workers), i(gp.Gen4Workers), i(gp.Gen2Tasks), i(gp.Gen3Tasks), i(gp.Gen4Tasks), i(telemetry.QualificationWorkers), i(telemetry.QualificationTasks), i(telemetry.DecodeWorkers),
		i(gp.FetchFoldScores), i(gp.FetchBlockReads), i(gp.FetchBlockSuccess), i(gp.FetchBlockFailed),
	}
	fmt.Printf("build92-row\t%s\n", strings.Join(fields, "\t"))
}

func TestExperimentalV4Build92DerivedTimingAccounting(t *testing.T) {
	telemetry := experimentalV4PhoneBuild90RecoveryTelemetry{}
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
	got := experimentalV4Build92DeriveTiming(telemetry)
	if got.geometryAccounted != 590*time.Millisecond || got.geometryUnaccounted != 110*time.Millisecond {
		t.Fatalf("geometry accounting=%s/%s want 590ms/110ms", got.geometryAccounted, got.geometryUnaccounted)
	}
	if got.postAccounted != 250*time.Millisecond || got.totalAccounted != 950*time.Millisecond || got.totalUnaccounted != 50*time.Millisecond {
		t.Fatalf("total accounting post=%s accounted=%s unaccounted=%s", got.postAccounted, got.totalAccounted, got.totalUnaccounted)
	}
}

func TestExperimentalV4Build92ExpectedTelemetryTable(t *testing.T) {
	mild, ok := experimentalV4Build91ExpectedForImage("phone-b-mild.jpg")
	if !ok || mild.bank != 937 || mild.qualified != 935 || mild.decode != 691 || mild.frames != 2120047 || !mild.accept {
		t.Fatalf("Build92 B/mild expected table invalid: %+v", mild)
	}
	angle, ok := experimentalV4Build91ExpectedForImage("phone-b-angle.jpg")
	if !ok || angle.bank != 6198 || angle.qualified != 0 || angle.accept {
		t.Fatalf("Build92 B/angle expected table invalid: %+v", angle)
	}
}

func TestExperimentalV4Build92RetainedTimingProfile(t *testing.T) {
	acquisitionDir, key, cw, ch, expectedPayload, cooldown := experimentalV4Build92Config(t)
	for _, imageName := range []string{"phone-b-mild.jpg", "phone-b-angle.jpg"} {
		work, boundary, prep := experimentalV4Build92Prepare(t, acquisitionDir, imageName, cw, ch)
		fmt.Printf("build92-prepare\t%s\t%d\n", imageName, prep.Milliseconds())
		if cooldown > 0 {
			time.Sleep(cooldown)
		}
		started := time.Now()
		payload, info, telemetry, err := experimentalV4PhoneBuild90Recover(work, boundary, key, cw, ch)
		wall := time.Since(started)
		experimentalV4Build92AssertExpected(t, imageName, wall, payload, info, telemetry, err, expectedPayload)
		experimentalV4Build92PrintRow(imageName, prep, wall, telemetry, err)
	}
}

func TestExperimentalV4Build92CPUProfile(t *testing.T) {
	acquisitionDir, key, cw, ch, expectedPayload, cooldown := experimentalV4Build92Config(t)
	imageName := os.Getenv("PIXSEAL_BUILD92_CASE")
	if imageName != "phone-b-mild.jpg" && imageName != "phone-b-angle.jpg" {
		t.Fatalf("PIXSEAL_BUILD92_CASE must be phone-b-mild.jpg or phone-b-angle.jpg, got %q", imageName)
	}
	profilePath := os.Getenv("PIXSEAL_BUILD92_CPU_PROFILE")
	if profilePath == "" {
		t.Fatal("PIXSEAL_BUILD92_CPU_PROFILE is required")
	}
	work, boundary, prep := experimentalV4Build92Prepare(t, acquisitionDir, imageName, cw, ch)
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
	payload, info, telemetry, recoverErr := experimentalV4PhoneBuild90Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		t.Fatalf("close CPU profile: %v", err)
	}
	experimentalV4Build92AssertExpected(t, imageName, wall, payload, info, telemetry, recoverErr, expectedPayload)
	fmt.Printf("build92-cpu-row\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%s\n", imageName, prep.Milliseconds(), wall.Milliseconds(), telemetry.TotalElapsed.Milliseconds(), telemetry.GeometryElapsed.Milliseconds(), telemetry.QualificationElapsed.Milliseconds(), telemetry.DecodeWallElapsed.Milliseconds(), profilePath)
}
