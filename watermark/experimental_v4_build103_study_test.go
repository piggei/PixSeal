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

func experimentalV4Build103ReadTable(t *testing.T, path string) []map[string]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Comma = '\t'
	rows, err := reader.ReadAll()
	if err != nil || len(rows) < 2 {
		t.Fatal("invalid frozen table")
	}
	out := []map[string]string{}
	for _, row := range rows[1:] {
		m := map[string]string{}
		for i, k := range rows[0] {
			m[k] = row[i]
		}
		out = append(out, m)
	}
	return out
}
func TestExperimentalV4Build103PrepareSnapshots(t *testing.T) {
	dir := os.Getenv("PIXSEAL_BUILD103_ACQUISITION_DIR")
	if dir == "" {
		t.Skip("private snapshot preparation is opt-in")
	}
	out := os.Getenv("PIXSEAL_BUILD103_OUTPUT_DIR")
	if !filepath.IsAbs(out) {
		t.Fatal("absolute private output required")
	}
	if os.Getenv("GODEBUG") != "asyncpreemptoff=0" {
		t.Fatal("snapshot preparation must use explicit default preemption")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD103_WIDTH"))
	if err != nil || cw < 1 {
		t.Fatal("invalid width")
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD103_HEIGHT"))
	if err != nil || ch < 1 {
		t.Fatal("invalid height")
	}
	captures := experimentalV4Build103ReadTable(t, os.Getenv("PIXSEAL_BUILD103_CAPTURE_REFERENCE"))
	fmt.Println("build103-capture-header\timage\tsnapshot_sha256\ttask_sha256\tplane_sha256\tmeta_sha256\tbank_sha256\ttasks\tbank\tgeometry_evals\tsemantic_gate")
	for _, image := range []string{"phone-b-mild.jpg", "phone-b-angle.jpg"} {
		fmt.Printf("build103-start\tprepare\t%s\n", image)
		work, boundary, _ := experimentalV4Build94Prepare(t, dir, image, cw, ch)
		bank, evals, seeds, workers, gp, fp := experimentalV4PhoneBuild98BlindBank(work, boundary, cw, ch)
		experimentalV4Build102CheckReference(t, os.Getenv("PIXSEAL_BUILD103_BANK_REFERENCE"), image, len(bank), evals, gp)
		bankHash := experimentalV4Build102Hash(bank)
		n := len(bank)
		bank = nil
		runtime.GC()
		f := experimentalV4Build102Prepare(work, boundary, cw, ch, false)
		s := experimentalV4Build103Capture(f)
		s.Image = image
		s.BankHash = bankHash
		s.Bank = n
		s.EvalsTotal = evals
		s.SeedsTotal = seeds
		s.WorkersTotal = workers
		s.Expected = experimentalV4Build103Geo(gp)
		s.ExpectedFreeze = fp
		matched := false
		for _, r := range captures {
			if r["image"] == image && r["origin"] == "build98" {
				matched = r["task_sha256"] == s.TaskHash && r["plane_sha256"] == s.PlaneHash && r["qualified_bank_sha256"] == bankHash && r["tasks"] == strconv.Itoa(len(f.tasks)) && r["bank"] == strconv.Itoa(n) && r["geometry_evals"] == strconv.Itoa(evals)
			}
		}
		if !matched {
			t.Fatal("snapshot differs from frozen Build102 case")
		}
		path := filepath.Join(out, image+".gob")
		if err := experimentalV4Build103Write(path, s); err != nil {
			t.Fatal(err)
		}
		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(bytes))
		loaded, err := experimentalV4Build103Read(path, hash)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := loaded.frozen(); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("build103-capture\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\tPASS\n", image, hash, s.TaskHash, s.PlaneHash, s.MetaHash, bankHash, len(f.tasks), n, evals)
	}
}

var experimentalV4Build103Columns = []string{"image", "pass", "run", "block", "pair", "position", "variant", "godebug", "pid", "snapshot_sha256", "task_sha256", "plane_sha256", "bank_sha256", "task_count", "bank", "geometry_evals", "replay_ms", "gen4_pool_ms", "gen4_worker_ms", "workers", "process_user_cpu_ms", "process_system_cpu_ms", "clock_ticks_per_s", "allocated_bytes", "mallocs", "completed_gc_cycles", "gc_pause_ns", "gen4_fold_scores", "gen4_block_reads", "gen4_block_success", "gen4_block_failed", "semantic_gate"}

// Each invocation is a fresh process with a fixed runtime condition, one loaded
// snapshot and one replay. No parent process changes GODEBUG mid-execution.
func TestExperimentalV4Build103ReplayProcess(t *testing.T) {
	path := os.Getenv("PIXSEAL_BUILD103_SNAPSHOT")
	if path == "" {
		t.Skip("private replay process is opt-in")
	}
	variant := os.Getenv("PIXSEAL_BUILD103_VARIANT")
	want := "asyncpreemptoff=0"
	if variant == "async_off" {
		want = "asyncpreemptoff=1"
	} else if variant != "default" {
		t.Fatal("invalid runtime variant")
	}
	if runtime.GOOS != "linux" || os.Getenv("GODEBUG") != want || os.Getenv("GOGC") != "100" || os.Getenv("GOMEMLIMIT") != "off" || runtime.GOMAXPROCS(0) != 8 {
		t.Fatal("runtime study environment differs")
	}
	parse := func(name string, low, high int) int {
		v, err := strconv.Atoi(os.Getenv(name))
		if err != nil || v < low || v > high {
			t.Fatalf("invalid %s", name)
		}
		return v
	}
	hz := parse("PIXSEAL_BUILD103_CLOCK_TICKS", 1, 1000000)
	cooldown := parse("PIXSEAL_BUILD103_COOLDOWN_MS", 0, 60000)
	hash := os.Getenv("PIXSEAL_BUILD103_SNAPSHOT_SHA256")
	s, err := experimentalV4Build103Read(path, hash)
	if err != nil {
		t.Fatal(err)
	}
	f, err := s.frozen()
	if err != nil {
		t.Fatal(err)
	}
	if len(f.tasks) == 0 {
		t.Fatal("empty frozen work")
	}
	runtime.GC()
	time.Sleep(time.Duration(cooldown) * time.Millisecond)
	pass := os.Getenv("PIXSEAL_BUILD103_PASS")
	var profile *os.File
	if pass == "cpu" {
		profile, err = os.OpenFile(os.Getenv("PIXSEAL_BUILD103_CPU_PROFILE"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			t.Fatal(err)
		}
		if err = pprof.StartCPUProfile(profile); err != nil {
			_ = profile.Close()
			t.Fatal(err)
		}
	} else if pass != "timing" {
		t.Fatal("invalid study pass")
	}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	u0, s0, err := experimentalV4Build102CPU()
	if err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	bank, evals, seeds, workers, gp, fp := experimentalV4Build102Replay(f)
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
	expectedFreeze := s.ExpectedFreeze
	experimentalV4Build100ClearDurations(reflect.ValueOf(&fp).Elem())
	experimentalV4Build100ClearDurations(reflect.ValueOf(&expectedFreeze).Elem())
	if u1 < u0 || s1 < s0 || len(bank) != s.Bank || evals != s.EvalsTotal || seeds != s.SeedsTotal || workers != s.WorkersTotal || experimentalV4Build102Hash(bank) != s.BankHash || !reflect.DeepEqual(experimentalV4Build102Semantic(gp.experimentalV4PhoneBuild98GeometryTelemetry), experimentalV4Build102Semantic(s.Expected.unpack())) || !reflect.DeepEqual(fp, expectedFreeze) || experimentalV4Build102Hash(f.tasks) != s.TaskHash || fmt.Sprintf("%x", sha256.Sum256(f.plane.rgb)) != s.PlaneHash || experimentalV4Build103Meta(f) != s.MetaHash {
		t.Fatal("replay changed bank, counters or frozen inputs")
	}
	fields := []string{s.Image, pass, os.Getenv("PIXSEAL_BUILD103_RUN"), os.Getenv("PIXSEAL_BUILD103_BLOCK"), os.Getenv("PIXSEAL_BUILD103_PAIR"), os.Getenv("PIXSEAL_BUILD103_POSITION"), variant, want, strconv.Itoa(os.Getpid()), hash, s.TaskHash, s.PlaneHash, s.BankHash, strconv.Itoa(len(f.tasks)), strconv.Itoa(len(bank)), strconv.Itoa(evals), strconv.FormatInt(wall.Milliseconds(), 10), strconv.FormatInt(gp.Gen4WallElapsed.Milliseconds(), 10), strconv.FormatInt(gp.Gen4WorkerElapsed.Milliseconds(), 10), strconv.Itoa(gp.Gen4Workers), fmt.Sprintf("%.3f", float64(u1-u0)*1000/float64(hz)), fmt.Sprintf("%.3f", float64(s1-s0)*1000/float64(hz)), strconv.Itoa(hz), strconv.FormatUint(after.TotalAlloc-before.TotalAlloc, 10), strconv.FormatUint(after.Mallocs-before.Mallocs, 10), strconv.FormatUint(uint64(after.NumGC-before.NumGC), 10), strconv.FormatUint(after.PauseTotalNs-before.PauseTotalNs, 10), strconv.Itoa(gp.FetchFoldScores), strconv.Itoa(gp.FetchBlockReads), strconv.Itoa(gp.FetchBlockSuccess), strconv.Itoa(gp.FetchBlockFailed), "PASS"}
	if len(fields) != len(experimentalV4Build103Columns) {
		t.Fatal("study row schema")
	}
	fmt.Printf("build103-header\t%s\n", strings.Join(experimentalV4Build103Columns, "\t"))
	fmt.Printf("build103-row\t%s\n", strings.Join(fields, "\t"))
}
