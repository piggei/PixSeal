package watermark

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Build99 is observability-only over the unchanged qualified Build98 runtime.
// The historical Build94 preparation, stage-accounting and inherited semantic
// assertions are shared deliberately. Build98 Gen3 counters are checked
// separately; no profile or timing influences recovery decisions.

func experimentalV4Build99Config(t *testing.T) (string, []byte, int, int, string, time.Duration) {
	t.Helper()
	acquisitionDir := os.Getenv("PIXSEAL_BUILD99_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build99 private retained profile requires PIXSEAL_BUILD99_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD99_KEY"))
	if len(key) < 8 {
		t.Fatal("Build99 PIXSEAL_BUILD99_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD99_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD99_WIDTH: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD99_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD99_HEIGHT: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD99_MESSAGE_B")
	if expectedPayload == "" {
		t.Fatal("PIXSEAL_BUILD99_MESSAGE_B is required")
	}
	cooldownMs := 5000
	if raw := os.Getenv("PIXSEAL_BUILD99_COOLDOWN_MS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 || v > 60000 {
			t.Fatalf("invalid PIXSEAL_BUILD99_COOLDOWN_MS=%q", raw)
		}
		cooldownMs = v
	}
	return acquisitionDir, key, cw, ch, expectedPayload, time.Duration(cooldownMs) * time.Millisecond
}

func experimentalV4Build99ValidateGen3(imageName string, gp experimentalV4PhoneBuild98GeometryTelemetry) error {
	expected := map[string][4]int{
		"phone-b-mild.jpg":  {23024, 29470720, 29470720, 0},
		"phone-b-angle.jpg": {76817, 98325760, 97394512, 931248},
	}
	want, ok := expected[imageName]
	if !ok {
		return fmt.Errorf("Build99 has no frozen Gen3 expectations for %s", imageName)
	}
	got := [4]int{gp.Gen3FetchFoldScores, gp.Gen3FetchBlockReads, gp.Gen3FetchBlockSuccess, gp.Gen3FetchBlockFailed}
	if got != want {
		return fmt.Errorf("Build99 %s Gen3 counters=%v want=%v", imageName, got, want)
	}
	return nil
}

func experimentalV4Build99AssertExpected(t *testing.T, imageName string, payload []byte, telemetry experimentalV4PhoneBuild98RecoveryTelemetry, err error, expectedPayload string) {
	t.Helper()
	experimentalV4Build94AssertExpected(t, imageName, payload, telemetry.experimentalV4PhoneBuild93RecoveryTelemetry, err, expectedPayload)
	if err := experimentalV4Build99ValidateGen3(imageName, telemetry.Build98GeometryProfile); err != nil {
		t.Fatal(err)
	}
}

var experimentalV4Build99RowColumns = append(append([]string{}, experimentalV4Build94RowColumns...),
	"gen3_fetch_fold_scores", "gen3_fetch_block_reads", "gen3_fetch_block_success", "gen3_fetch_block_failed")

func experimentalV4Build99RowFields(imageName string, prep, wall time.Duration, telemetry experimentalV4PhoneBuild98RecoveryTelemetry, err error) []string {
	fields := experimentalV4Build94RowFields(imageName, prep, wall, telemetry.experimentalV4PhoneBuild93RecoveryTelemetry, err)
	gp := telemetry.Build98GeometryProfile
	return append(fields, strconv.Itoa(gp.Gen3FetchFoldScores), strconv.Itoa(gp.Gen3FetchBlockReads), strconv.Itoa(gp.Gen3FetchBlockSuccess), strconv.Itoa(gp.Gen3FetchBlockFailed))
}

func experimentalV4Build99PrintRow(imageName string, prep, wall time.Duration, telemetry experimentalV4PhoneBuild98RecoveryTelemetry, err error) {
	fmt.Printf("build99-row\t%s\n", strings.Join(experimentalV4Build99RowFields(imageName, prep, wall, telemetry, err), "\t"))
}

func TestExperimentalV4Build99Gen3Gate(t *testing.T) {
	cases := []struct {
		image    string
		counters [4]int
	}{
		{"phone-b-mild.jpg", [4]int{23024, 29470720, 29470720, 0}},
		{"phone-b-angle.jpg", [4]int{76817, 98325760, 97394512, 931248}},
	}
	for _, tc := range cases {
		gp := experimentalV4PhoneBuild98GeometryTelemetry{Gen3FetchFoldScores: tc.counters[0], Gen3FetchBlockReads: tc.counters[1], Gen3FetchBlockSuccess: tc.counters[2], Gen3FetchBlockFailed: tc.counters[3]}
		if err := experimentalV4Build99ValidateGen3(tc.image, gp); err != nil {
			t.Fatal(err)
		}
		for i := range tc.counters {
			mutated := gp
			switch i {
			case 0:
				mutated.Gen3FetchFoldScores++
			case 1:
				mutated.Gen3FetchBlockReads++
			case 2:
				mutated.Gen3FetchBlockSuccess++
			case 3:
				mutated.Gen3FetchBlockFailed++
			}
			if experimentalV4Build99ValidateGen3(tc.image, mutated) == nil {
				t.Fatalf("%s accepted changed Gen3 counter %d", tc.image, i)
			}
		}
	}
	if experimentalV4Build99ValidateGen3("unknown.jpg", experimentalV4PhoneBuild98GeometryTelemetry{}) == nil {
		t.Fatal("accepted unknown case")
	}
}

func TestExperimentalV4Build99RowShape(t *testing.T) {
	telemetry := experimentalV4PhoneBuild98RecoveryTelemetry{}
	telemetry.Build98GeometryProfile.Gen3FetchFoldScores = 23024
	fields := experimentalV4Build99RowFields("phone-b-mild.jpg", 0, 0, telemetry, nil)
	if len(fields) != len(experimentalV4Build99RowColumns) {
		t.Fatalf("row fields=%d columns=%d", len(fields), len(experimentalV4Build99RowColumns))
	}
	seen := map[string]bool{}
	for _, name := range experimentalV4Build99RowColumns {
		if name == "" || seen[name] || strings.ContainsAny(name, "\t\r\n") {
			t.Fatalf("invalid or duplicate column %q", name)
		}
		seen[name] = true
	}
	if fields[len(fields)-4] != "23024" {
		t.Fatal("Gen3 fields not appended to inherited row")
	}
}

func TestExperimentalV4Build99RetainedTimingProfile(t *testing.T) {
	acquisitionDir, key, cw, ch, expectedPayload, cooldown := experimentalV4Build99Config(t)
	fmt.Printf("build99-header\t%s\n", strings.Join(experimentalV4Build99RowColumns, "\t"))
	for _, imageName := range []string{"phone-b-mild.jpg", "phone-b-angle.jpg"} {
		work, boundary, prep := experimentalV4Build94Prepare(t, acquisitionDir, imageName, cw, ch)
		fmt.Printf("build99-prepare\t%s\t%d\n", imageName, prep.Milliseconds())
		if cooldown > 0 {
			time.Sleep(cooldown)
		}
		started := time.Now()
		payload, _, telemetry, err := experimentalV4PhoneBuild98Recover(work, boundary, key, cw, ch)
		wall := time.Since(started)
		experimentalV4Build99AssertExpected(t, imageName, payload, telemetry, err, expectedPayload)
		experimentalV4Build99PrintRow(imageName, prep, wall, telemetry, err)
	}
}

func TestExperimentalV4Build99CPUProfile(t *testing.T) {
	acquisitionDir, key, cw, ch, expectedPayload, cooldown := experimentalV4Build99Config(t)
	imageName := os.Getenv("PIXSEAL_BUILD99_CASE")
	if imageName != "phone-b-mild.jpg" && imageName != "phone-b-angle.jpg" {
		t.Fatalf("PIXSEAL_BUILD99_CASE must be phone-b-mild.jpg or phone-b-angle.jpg, got %q", imageName)
	}
	profilePath := os.Getenv("PIXSEAL_BUILD99_CPU_PROFILE")
	if profilePath == "" {
		t.Fatal("PIXSEAL_BUILD99_CPU_PROFILE is required")
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
	payload, _, telemetry, recoverErr := experimentalV4PhoneBuild98Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		t.Fatalf("close CPU profile: %v", err)
	}
	experimentalV4Build99AssertExpected(t, imageName, payload, telemetry, recoverErr, expectedPayload)
	fmt.Printf("build99-cpu-row\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\n",
		imageName, prep.Milliseconds(), wall.Milliseconds(), telemetry.TotalElapsed.Milliseconds(), telemetry.GeometryElapsed.Milliseconds(), telemetry.QualificationElapsed.Milliseconds(), telemetry.DecodeWallElapsed.Milliseconds(), telemetry.SamplingWorkerElapsed.Milliseconds(), telemetry.ListWorkerElapsed.Milliseconds(), telemetry.ProtectedSamplingPasses, telemetry.ProtectedProfileDerivations, profilePath)
}
