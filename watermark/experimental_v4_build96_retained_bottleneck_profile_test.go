package watermark

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"

	jpeglegacy "github.com/pj/pixseal/internal/jpeglegacy"
)

type experimentalV4Build96Mem struct {
	totalAlloc   uint64
	mallocs      uint64
	frees        uint64
	heapAlloc    uint64
	heapInuse    uint64
	heapSys      uint64
	numGC        uint32
	pauseTotalNs uint64
}

func experimentalV4Build96MemSnapshot() experimentalV4Build96Mem {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return experimentalV4Build96Mem{m.TotalAlloc, m.Mallocs, m.Frees, m.HeapAlloc, m.HeapInuse, m.HeapSys, m.NumGC, m.PauseTotalNs}
}

type experimentalV4Build96Run struct {
	base    experimentalV4Build95PairedRun
	gp      experimentalV4PhoneBuild84GeometryTelemetry
	before  experimentalV4Build96Mem
	after   experimentalV4Build96Mem
	profile string
}

func experimentalV4Build96RunVariant(t *testing.T, variant string, work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int, profilePath string) experimentalV4Build96Run {
	t.Helper()
	runtime.GC()
	f, err := os.Create(profilePath)
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		_ = f.Close()
		t.Fatalf("start profile: %v", err)
	}
	before := experimentalV4Build96MemSnapshot()

	var base experimentalV4Build95PairedRun
	var gp experimentalV4PhoneBuild84GeometryTelemetry
	if variant == "build93" {
		started := time.Now()
		payload, info, telemetry, recoverErr := experimentalV4PhoneBuild93Recover(work, boundary, key, cw, ch)
		wall := time.Since(started)
		gp = telemetry.Build84GeometryProfile
		physical := telemetry.PhysicalDecodeCandidates
		base = experimentalV4Build95PairedRun{
			variant: variant, outcome: map[bool]string{true: "accept", false: "reject"}[recoverErr == nil], wall: wall,
			total: telemetry.TotalElapsed, geometry: telemetry.GeometryElapsed, plane: telemetry.PlanePrepElapsed,
			freeze: gp.FreezeElapsed, prefix1: gp.Prefix1WallElapsed, gen2: gp.Gen2WallElapsed, gen3: gp.Gen3WallElapsed, gen4: gp.Gen4WallElapsed,
			qualification: telemetry.QualificationElapsed, decodeWall: telemetry.DecodeWallElapsed,
			geometryEvals: telemetry.GeometryEvaluations, bank: telemetry.BankCandidates, qualEvals: telemetry.QualificationEvaluations,
			qualified: telemetry.QualifiedCandidates, decodeCandidates: telemetry.DecodeCandidatesTried, listFrames: telemetry.ListFramesTried,
			physicalDecode: physical, speculativeDecode: physical - telemetry.DecodeCandidatesTried,
			physicalProfiles: telemetry.PhysicalProfilesTried, physicalListFrames: telemetry.PhysicalListFramesTried,
			protectedSamplingPasses: telemetry.ProtectedSamplingPasses, protectedProfileDerivations: telemetry.ProtectedProfileDerivations,
			authenticated: telemetry.Authenticated, qualificationWorkers: telemetry.QualificationWorkers, qualificationTasks: telemetry.QualificationTasks,
			foldScores: gp.FetchFoldScores, blockReads: gp.FetchBlockReads, blockSuccess: gp.FetchBlockSuccess, blockFailed: gp.FetchBlockFailed,
			payload: payload, info: info, err: recoverErr,
		}
	} else {
		started := time.Now()
		payload, info, telemetry, recoverErr := experimentalV4PhoneBuild95Recover(work, boundary, key, cw, ch)
		wall := time.Since(started)
		gp95 := telemetry.Build95GeometryProfile
		gp = gp95.experimentalV4PhoneBuild84GeometryTelemetry
		physical := telemetry.PhysicalDecodeCandidates
		base = experimentalV4Build95PairedRun{
			variant: variant, outcome: map[bool]string{true: "accept", false: "reject"}[recoverErr == nil], wall: wall,
			total: telemetry.TotalElapsed, geometry: telemetry.GeometryElapsed, plane: telemetry.PlanePrepElapsed,
			freeze: gp.FreezeElapsed, prefix1: gp.Prefix1WallElapsed, gen2: gp.Gen2WallElapsed, gen3: gp.Gen3WallElapsed, gen4: gp.Gen4WallElapsed,
			luminancePrep: gp95.LuminancePlanePrepElapsed, luminanceBytes: gp95.LuminancePlaneBytes, luminanceUsed: gp95.LuminancePlaneUsed, luminanceFallback: gp95.LuminancePlaneFallback,
			qualification: telemetry.QualificationElapsed, decodeWall: telemetry.DecodeWallElapsed,
			geometryEvals: telemetry.GeometryEvaluations, bank: telemetry.BankCandidates, qualEvals: telemetry.QualificationEvaluations,
			qualified: telemetry.QualifiedCandidates, decodeCandidates: telemetry.DecodeCandidatesTried, listFrames: telemetry.ListFramesTried,
			physicalDecode: physical, speculativeDecode: physical - telemetry.DecodeCandidatesTried,
			physicalProfiles: telemetry.PhysicalProfilesTried, physicalListFrames: telemetry.PhysicalListFramesTried,
			protectedSamplingPasses: telemetry.ProtectedSamplingPasses, protectedProfileDerivations: telemetry.ProtectedProfileDerivations,
			authenticated: telemetry.Authenticated, qualificationWorkers: telemetry.QualificationWorkers, qualificationTasks: telemetry.QualificationTasks,
			foldScores: gp.FetchFoldScores, blockReads: gp.FetchBlockReads, blockSuccess: gp.FetchBlockSuccess, blockFailed: gp.FetchBlockFailed,
			payload: payload, info: info, err: recoverErr,
		}
	}
	after := experimentalV4Build96MemSnapshot()
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		t.Fatalf("close profile: %v", err)
	}
	return experimentalV4Build96Run{base: base, gp: gp, before: before, after: after, profile: filepath.Base(profilePath)}
}

var experimentalV4Build96Columns = []string{
	"cycle", "order", "position", "variant", "outcome", "wall_ms", "total_ms", "geometry_ms", "gen4_wall_ms", "gen4_worker_ms", "gen4_workers", "gen4_tasks", "gen4_min_ms", "gen4_median_ms", "gen4_max_ms", "luminance_prep_ms", "luminance_bytes", "total_alloc_delta", "mallocs_delta", "frees_delta", "heap_alloc_before", "heap_alloc_after", "heap_inuse_before", "heap_inuse_after", "heap_sys_before", "heap_sys_after", "num_gc_delta", "pause_ns_delta", "geometry_evals", "bank", "qualification_evals", "qualified", "fold_scores", "block_reads", "block_success", "block_failed", "profile",
}

func experimentalV4Build96Delta(a, b uint64) uint64 {
	if b >= a {
		return b - a
	}
	return 0
}

func experimentalV4Build96Fields(cycle int, order string, position int, r experimentalV4Build96Run) []string {
	b, g := r.base, r.gp
	return []string{
		strconv.Itoa(cycle), order, strconv.Itoa(position), b.variant, b.outcome,
		strconv.FormatInt(b.wall.Milliseconds(), 10), strconv.FormatInt(b.total.Milliseconds(), 10), strconv.FormatInt(b.geometry.Milliseconds(), 10), strconv.FormatInt(b.gen4.Milliseconds(), 10), strconv.FormatInt(g.Gen4WorkerElapsed.Milliseconds(), 10),
		strconv.Itoa(g.Gen4Workers), strconv.Itoa(g.Gen4Tasks), strconv.FormatInt(g.Gen4MinElapsed.Milliseconds(), 10), strconv.FormatInt(g.Gen4MedianElapsed.Milliseconds(), 10), strconv.FormatInt(g.Gen4MaxElapsed.Milliseconds(), 10),
		strconv.FormatInt(b.luminancePrep.Milliseconds(), 10), strconv.FormatInt(b.luminanceBytes, 10),
		strconv.FormatUint(experimentalV4Build96Delta(r.before.totalAlloc, r.after.totalAlloc), 10), strconv.FormatUint(experimentalV4Build96Delta(r.before.mallocs, r.after.mallocs), 10), strconv.FormatUint(experimentalV4Build96Delta(r.before.frees, r.after.frees), 10),
		strconv.FormatUint(r.before.heapAlloc, 10), strconv.FormatUint(r.after.heapAlloc, 10), strconv.FormatUint(r.before.heapInuse, 10), strconv.FormatUint(r.after.heapInuse, 10), strconv.FormatUint(r.before.heapSys, 10), strconv.FormatUint(r.after.heapSys, 10),
		strconv.FormatUint(uint64(r.after.numGC-r.before.numGC), 10), strconv.FormatUint(experimentalV4Build96Delta(r.before.pauseTotalNs, r.after.pauseTotalNs), 10),
		strconv.Itoa(b.geometryEvals), strconv.Itoa(b.bank), strconv.Itoa(b.qualEvals), strconv.Itoa(b.qualified), strconv.Itoa(b.foldScores), strconv.Itoa(b.blockReads), strconv.Itoa(b.blockSuccess), strconv.Itoa(b.blockFailed), r.profile,
	}
}

func TestExperimentalV4Build96RowShape(t *testing.T) {
	f := experimentalV4Build96Fields(1, "93-95", 1, experimentalV4Build96Run{base: experimentalV4Build95PairedRun{variant: "build93"}})
	if len(f) != len(experimentalV4Build96Columns) {
		t.Fatalf("Build96 fields=%d columns=%d", len(f), len(experimentalV4Build96Columns))
	}
}

func TestExperimentalV4Build96Gen4RetainedBottleneckProfile(t *testing.T) {
	acquisitionDir := os.Getenv("PIXSEAL_BUILD96_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build96 private retained profile requires corpus")
	}
	outputDir := os.Getenv("PIXSEAL_BUILD96_PROFILE_DIR")
	if outputDir == "" {
		t.Fatal("PIXSEAL_BUILD96_PROFILE_DIR is required")
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD96_KEY"))
	if len(key) < 8 {
		t.Fatal("Build96 key must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD96_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatal("invalid Build96 width")
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD96_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatal("invalid Build96 height")
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD96_MESSAGE_B")
	cooldownMS, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD96_COOLDOWN_MS"))
	if err != nil || cooldownMS < 0 {
		t.Fatal("invalid Build96 cooldown")
	}
	cooldown := time.Duration(cooldownMS) * time.Millisecond

	fmt.Printf("build96-header\t%s\n", strings.Join(experimentalV4Build96Columns, "\t"))
	for cycle := 1; cycle <= 2; cycle++ {
		path := filepath.Join(acquisitionDir, "phone-b-angle.jpg")
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		src, err := jpeglegacy.Decode(f)
		_ = f.Close()
		if err != nil {
			t.Fatal(err)
		}
		work, boundary, hyp, _, _, _ := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
		if len(hyp) >= experimentalV4PhoneEnsembleSize {
			t.Fatal("Build96 B/angle unexpectedly accepted before deep recovery")
		}
		order := "93-95"
		variants := []string{"build93", "build95"}
		if cycle == 2 {
			order = "95-93"
			variants = []string{"build95", "build93"}
		}
		var pair []experimentalV4Build96Run
		for pos, variant := range variants {
			if cooldown > 0 {
				time.Sleep(cooldown)
			}
			profile := filepath.Join(outputDir, fmt.Sprintf("build96-cycle%d-pos%d-%s.pprof", cycle, pos+1, variant))
			r := experimentalV4Build96RunVariant(t, variant, work, boundary, key, cw, ch, profile)
			experimentalV4Build95AssertExpected(t, "phone-b-angle.jpg", r.base, expectedPayload, work)
			fmt.Printf("build96-row\t%s\n", strings.Join(experimentalV4Build96Fields(cycle, order, pos+1, r), "\t"))
			pair = append(pair, r)
		}
		experimentalV4Build95AssertPairEqual(t, "phone-b-angle.jpg", pair[0].base, pair[1].base)
	}
}
