package watermark

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"
)

func experimentalV4Build102Semantic(g experimentalV4PhoneBuild98GeometryTelemetry) experimentalV4PhoneBuild98GeometryTelemetry {
	experimentalV4Build100ClearDurations(reflect.ValueOf(&g).Elem())
	return g
}

func experimentalV4Build102CheckReference(t *testing.T, path, image string, bank, evals int, gp experimentalV4PhoneBuild98GeometryTelemetry) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = '\t'
	rows, err := r.ReadAll()
	if err != nil || len(rows) < 2 {
		t.Fatal("invalid frozen Build101 reference")
	}
	indices := map[string]int{}
	for i, k := range rows[0] {
		indices[k] = i
	}
	for _, row := range rows[1:] {
		if row[indices["image"]] != image || row[indices["variant"]] != "build98" || row[indices["pass"]] != "timing" {
			continue
		}
		for k, want := range map[string]int{"bank": bank, "geometry_evals": evals, "gen4_tasks": gp.Gen4Tasks, "fold_scores": gp.FetchFoldScores, "block_reads": gp.FetchBlockReads, "block_success": gp.FetchBlockSuccess, "block_failed": gp.FetchBlockFailed} {
			i, ok := indices[k]
			if !ok || row[i] != strconv.Itoa(want) {
				t.Fatalf("frozen %s %s reference differs", image, k)
			}
		}
		if err := experimentalV4Build99ValidateGen3(image, gp); err != nil {
			t.Fatal(err)
		}
		return
	}
	t.Fatal("missing frozen case reference")
}

// This is an A/A control: every slot uses the same plane, task slice, kernel,
// worker count and collection code. A/B are balanced measurement labels only.
func TestExperimentalV4Build102RetainedReplay(t *testing.T) {
	dir := os.Getenv("PIXSEAL_BUILD102_ACQUISITION_DIR")
	if dir == "" {
		t.Skip("private retained replay is opt-in")
	}
	if runtime.GOOS != "linux" {
		t.Fatal("Build102 process CPU study requires Linux")
	}
	out := os.Getenv("PIXSEAL_BUILD102_OUTPUT_DIR")
	if !filepath.IsAbs(out) {
		t.Fatal("absolute output directory required")
	}
	parse := func(name string, lo, hi int) int {
		v, err := strconv.Atoi(os.Getenv(name))
		if err != nil || v < lo || v > hi {
			t.Fatalf("invalid %s", name)
		}
		return v
	}
	cw := parse("PIXSEAL_BUILD102_WIDTH", 1, 20000)
	ch := parse("PIXSEAL_BUILD102_HEIGHT", 1, 20000)
	hz := parse("PIXSEAL_BUILD102_CLOCK_TICKS", 1, 1000000)
	cooldown := time.Duration(parse("PIXSEAL_BUILD102_COOLDOWN_MS", 0, 60000)) * time.Millisecond
	ref := os.Getenv("PIXSEAL_BUILD102_REFERENCE_TSV")
	columns := []string{"image", "pass", "run", "block", "pair", "position", "slot", "task_count", "task_sha256", "plane_sha256", "bank_sha256", "bank", "geometry_evals", "replay_ms", "gen4_pool_ms", "gen4_worker_ms", "workers", "process_user_cpu_ms", "process_system_cpu_ms", "clock_ticks_per_s", "allocated_bytes", "mallocs", "completed_gc_cycles", "gc_pause_ns", "gen4_fold_scores", "gen4_block_reads", "gen4_block_success", "gen4_block_failed", "semantic_gate"}
	fmt.Printf("build102-header\t%s\n", strings.Join(columns, "\t"))
	fmt.Println("build102-capture-header\timage\torigin\tprepare_ms\tprefix_ms\ttasks\ttask_sha256\tplane_sha256\tqualified_bank_sha256\tbank\tgeometry_evals\tsemantic_gate")
	run := 0
	for _, image := range []string{"phone-b-mild.jpg", "phone-b-angle.jpg"} {
		work, boundary, prep := experimentalV4Build94Prepare(t, dir, image, cw, ch)
		runtime.GC()
		time.Sleep(cooldown)
		fmt.Printf("build102-start\treference\t%s\n", image)
		gold, ge, gs, gw, gg, gf := experimentalV4PhoneBuild98BlindBank(work, boundary, cw, ch)
		experimentalV4Build102CheckReference(t, ref, image, len(gold), ge, gg)
		bankHash := experimentalV4Build102Hash(gold)
		bankCount := len(gold)
		gold = nil
		goldenGP := experimentalV4Build102Semantic(gg)
		experimentalV4Build100ClearDurations(reflect.ValueOf(&gf).Elem())
		var frozen experimentalV4Build102Frozen
		var taskHash, planeHash string
		for _, candidate := range []bool{false, true} {
			runtime.GC()
			time.Sleep(cooldown)
			origin := "build98"
			if candidate {
				origin = "build100"
			}
			fmt.Printf("build102-start\tcapture\t%s\t%s\n", image, origin)
			started := time.Now()
			f := experimentalV4Build102Prepare(work, boundary, cw, ch, candidate)
			wall := time.Since(started)
			if len(f.tasks) == 0 || f.seeds != gs || f.workers != gw {
				t.Fatal("invalid or empty frozen capture")
			}
			hash := experimentalV4Build102Hash(f.tasks)
			plane := fmt.Sprintf("%x", sha256.Sum256(f.plane.rgb))
			if !candidate {
				frozen = f
				taskHash = hash
				planeHash = plane
			} else {
				if hash != taskHash || plane != planeHash || f.plane.bounds != frozen.plane.bounds || !reflect.DeepEqual(f.tasks, frozen.tasks) || f.evals != frozen.evals || !reflect.DeepEqual(experimentalV4Build102Semantic(f.profile.experimentalV4PhoneBuild98GeometryTelemetry), experimentalV4Build102Semantic(frozen.profile.experimentalV4PhoneBuild98GeometryTelemetry)) {
					t.Fatal("prefix origins changed frozen work")
				}
			}
			fmt.Printf("build102-capture\t%s\t%s\t%d\t%d\t%d\t%s\t%s\t%s\t%d\t%d\tPASS\n", image, origin, prep.Milliseconds(), wall.Milliseconds(), len(f.tasks), hash, plane, bankHash, bankCount, ge)
		}
		for _, pass := range []string{"timing", "cpu"} {
			slots := []string{"A", "B", "B", "A", "B", "A", "A", "B"}
			if pass == "cpu" {
				slots = []string{"A", "B"}
				if image == "phone-b-angle.jpg" {
					slots = []string{"B", "A"}
				}
			}
			for i, slot := range slots {
				run++
				runtime.GC()
				time.Sleep(cooldown)
				fmt.Printf("build102-start\t%s\t%d\t%s\t%s\n", pass, run, image, slot)
				var profile *os.File
				if pass == "cpu" {
					var err error
					profile, err = os.OpenFile(filepath.Join(out, fmt.Sprintf("run%02d-%s-%s.cpu.prof", run, image, slot)), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
					if err != nil {
						t.Fatal(err)
					}
					if err := pprof.StartCPUProfile(profile); err != nil {
						_ = profile.Close()
						t.Fatal(err)
					}
				}
				var before, after runtime.MemStats
				runtime.ReadMemStats(&before)
				u0, s0, err := experimentalV4Build102CPU()
				if err != nil {
					t.Fatal(err)
				}
				started := time.Now()
				bank, evals, seeds, workers, gp, fp := experimentalV4Build102Replay(frozen)
				wall := time.Since(started)
				u1, s1, err := experimentalV4Build102CPU()
				if err != nil {
					t.Fatal(err)
				}
				runtime.ReadMemStats(&after)
				if profile != nil {
					pprof.StopCPUProfile()
					if err := profile.Close(); err != nil {
						t.Fatal(err)
					}
				}
				if u1 < u0 || s1 < s0 {
					t.Fatal("nonmonotonic process CPU")
				}
				experimentalV4Build100ClearDurations(reflect.ValueOf(&fp).Elem())
				if experimentalV4Build102Hash(bank) != bankHash || len(bank) != bankCount || evals != ge || seeds != gs || workers != gw || !reflect.DeepEqual(experimentalV4Build102Semantic(gp.experimentalV4PhoneBuild98GeometryTelemetry), goldenGP) || !reflect.DeepEqual(fp, gf) || experimentalV4Build102Hash(frozen.tasks) != taskHash || fmt.Sprintf("%x", sha256.Sum256(frozen.plane.rgb)) != planeHash {
					t.Fatalf("run %d changed frozen work, bank, counters or raster", run)
				}
				vals := []string{image, pass, strconv.Itoa(run), strconv.Itoa(i/4 + 1), strconv.Itoa(i/2 + 1), strconv.Itoa(i%2 + 1), slot, strconv.Itoa(len(frozen.tasks)), taskHash, planeHash, bankHash, strconv.Itoa(len(bank)), strconv.Itoa(evals), strconv.FormatInt(wall.Milliseconds(), 10), strconv.FormatInt(gp.Gen4WallElapsed.Milliseconds(), 10), strconv.FormatInt(gp.Gen4WorkerElapsed.Milliseconds(), 10), strconv.Itoa(gp.Gen4Workers), fmt.Sprintf("%.3f", float64(u1-u0)*1000/float64(hz)), fmt.Sprintf("%.3f", float64(s1-s0)*1000/float64(hz)), strconv.Itoa(hz), strconv.FormatUint(after.TotalAlloc-before.TotalAlloc, 10), strconv.FormatUint(after.Mallocs-before.Mallocs, 10), strconv.FormatUint(uint64(after.NumGC-before.NumGC), 10), strconv.FormatUint(after.PauseTotalNs-before.PauseTotalNs, 10), strconv.Itoa(gp.FetchFoldScores), strconv.Itoa(gp.FetchBlockReads), strconv.Itoa(gp.FetchBlockSuccess), strconv.Itoa(gp.FetchBlockFailed), "PASS"}
				if len(vals) != len(columns) {
					t.Fatal("replay row schema")
				}
				fmt.Printf("build102-row\t%s\n", strings.Join(vals, "\t"))
			}
		}
	}
}
