package watermark

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func experimentalV4Build100Config(t *testing.T) (string, []byte, int, int, string, time.Duration) {
	t.Helper()
	acquisitionDir := os.Getenv("PIXSEAL_BUILD100_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build100 private retained profile requires PIXSEAL_BUILD100_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD100_KEY"))
	if len(key) < 8 {
		t.Fatal("Build100 PIXSEAL_BUILD100_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD100_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD100_WIDTH: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD100_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD100_HEIGHT: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD100_MESSAGE_B")
	if expectedPayload == "" {
		t.Fatal("PIXSEAL_BUILD100_MESSAGE_B is required")
	}
	cooldownMs := 5000
	if raw := os.Getenv("PIXSEAL_BUILD100_COOLDOWN_MS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 || v > 60000 {
			t.Fatalf("invalid PIXSEAL_BUILD100_COOLDOWN_MS=%q", raw)
		}
		cooldownMs = v
	}
	return acquisitionDir, key, cw, ch, expectedPayload, time.Duration(cooldownMs) * time.Millisecond
}

var experimentalV4Build100PairedColumns = append(append([]string{"cycle", "order", "position", "variant"}, experimentalV4Build99RowColumns...), "freeze_pair_ms", "pair_workers", "pair_tasks", "pair_evaluations", "pair_worker_ms")

func TestExperimentalV4Build100PairedRetainedProfile(t *testing.T) {
	dir, key, cw, ch, expected, cooldown := experimentalV4Build100Config(t)
	fmt.Printf("build100-header\t%s\n", strings.Join(experimentalV4Build100PairedColumns, "\t"))
	scenarios := []struct {
		cycle   int
		image   string
		reverse bool
	}{
		{1, "phone-b-mild.jpg", false}, {1, "phone-b-angle.jpg", false},
		{2, "phone-b-angle.jpg", true}, {2, "phone-b-mild.jpg", true},
	}
	for _, sc := range scenarios {
		work, boundary, prep := experimentalV4Build94Prepare(t, dir, sc.image, cw, ch)
		variants := []string{"build98", "build100"}
		order := "98-100"
		if sc.reverse {
			variants = []string{"build100", "build98"}
			order = "100-98"
		}
		var previous experimentalV4PhoneBuild98RecoveryTelemetry
		var previousPayload []byte
		var previousInfo ExperimentalV4ExtractInfo
		var previousOK bool
		for position, variant := range variants {
			runtime.GC()
			if cooldown > 0 {
				time.Sleep(cooldown)
			}
			var payload []byte
			var info ExperimentalV4ExtractInfo
			var tele experimentalV4PhoneBuild98RecoveryTelemetry
			var pp experimentalV4PhoneBuild100PairTelemetry
			var err error
			started := time.Now()
			if variant == "build98" {
				payload, info, tele, err = experimentalV4PhoneBuild98Recover(work, boundary, key, cw, ch)
				pp = experimentalV4PhoneBuild100PairTelemetry{Workers: 1, Tasks: tele.FreezeProfile.PairScoreTasks, Evaluations: tele.FreezeProfile.PairScoreEvaluations, WallElapsed: tele.FreezeProfile.PairScoreElapsed, WorkerElapsed: tele.FreezeProfile.PairScoreElapsed}
			} else {
				var candidate experimentalV4PhoneBuild100RecoveryTelemetry
				payload, info, candidate, err = experimentalV4PhoneBuild100Recover(work, boundary, key, cw, ch)
				tele = candidate.experimentalV4PhoneBuild98RecoveryTelemetry
				pp = candidate.PairScoring
				wantWorkers := runtime.GOMAXPROCS(0)
				if wantWorkers > 6 {
					wantWorkers = 6
				}
				if pp.Workers != wantWorkers || pp.Tasks != 6 || pp.Evaluations != tele.FreezeProfile.PairScoreEvaluations {
					t.Fatal("candidate pair work accounting changed")
				}
			}
			wall := time.Since(started)
			experimentalV4Build99AssertExpected(t, sc.image, payload, tele, err, expected)
			semantic := tele
			experimentalV4Build100ClearDurations(reflect.ValueOf(&semantic).Elem())
			if position == 0 {
				previous = semantic
				previousPayload = append([]byte(nil), payload...)
				previousInfo = info
				previousOK = err == nil
			} else {
				if !reflect.DeepEqual(previous, semantic) || !bytes.Equal(previousPayload, payload) || previousOK != (err == nil) || (err == nil && !reflect.DeepEqual(previousInfo, info)) {
					t.Fatalf("%s cycle %d pair semantics differ", sc.image, sc.cycle)
				}
			}
			fields := []string{strconv.Itoa(sc.cycle), order, strconv.Itoa(position + 1), variant}
			fields = append(fields, experimentalV4Build99RowFields(sc.image, prep, wall, tele, err)...)
			fields = append(fields, strconv.FormatInt(pp.WallElapsed.Milliseconds(), 10), strconv.Itoa(pp.Workers), strconv.Itoa(pp.Tasks), strconv.Itoa(pp.Evaluations), strconv.FormatInt(pp.WorkerElapsed.Milliseconds(), 10))
			if len(fields) != len(experimentalV4Build100PairedColumns) {
				t.Fatal("paired row shape")
			}
			fmt.Printf("build100-row\t%s\n", strings.Join(fields, "\t"))
		}
	}
}
