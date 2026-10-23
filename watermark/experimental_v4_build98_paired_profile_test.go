package watermark

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	jpeglegacy "github.com/pj/pixseal/internal/jpeglegacy"
)

type experimentalV4Build98PairedRun struct {
	variant                     string
	outcome                     string
	wall                        time.Duration
	total                       time.Duration
	geometry                    time.Duration
	plane                       time.Duration
	freeze                      time.Duration
	prefix1                     time.Duration
	gen2                        time.Duration
	gen3                        time.Duration
	gen4                        time.Duration
	qualification               time.Duration
	decodeWall                  time.Duration
	geometryEvals               int
	bank                        int
	qualEvals                   int
	qualified                   int
	decodeCandidates            int
	listFrames                  int
	physicalDecode              int
	speculativeDecode           int
	physicalProfiles            int
	physicalListFrames          int
	protectedSamplingPasses     int
	protectedProfileDerivations int
	authenticated               bool
	qualificationWorkers        int
	qualificationTasks          int
	foldScores                  int
	blockReads                  int
	blockSuccess                int
	blockFailed                 int
	gen3FetchFoldScores         int
	gen3FetchBlockReads         int
	gen3FetchBlockSuccess       int
	gen3FetchBlockFailed        int
	payload                     []byte
	info                        ExperimentalV4ExtractInfo
	err                         error
}

func experimentalV4Build98RunFromBuild93(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) experimentalV4Build98PairedRun {
	started := time.Now()
	payload, info, telemetry, err := experimentalV4PhoneBuild93Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	gp := telemetry.Build84GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	return experimentalV4Build98PairedRun{
		variant: "build93", outcome: map[bool]string{true: "accept", false: "reject"}[err == nil], wall: wall,
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
		payload: payload, info: info, err: err,
	}
}

func experimentalV4Build98RunFromBuild98(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) experimentalV4Build98PairedRun {
	started := time.Now()
	payload, info, telemetry, err := experimentalV4PhoneBuild98Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	gp := telemetry.Build98GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	return experimentalV4Build98PairedRun{
		variant: "build98", outcome: map[bool]string{true: "accept", false: "reject"}[err == nil], wall: wall,
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
		gen3FetchFoldScores: gp.Gen3FetchFoldScores, gen3FetchBlockReads: gp.Gen3FetchBlockReads,
		gen3FetchBlockSuccess: gp.Gen3FetchBlockSuccess, gen3FetchBlockFailed: gp.Gen3FetchBlockFailed,
		payload: payload, info: info, err: err,
	}
}

func experimentalV4Build98AssertExpected(t *testing.T, imageName string, r experimentalV4Build98PairedRun, expectedPayload string) {
	t.Helper()
	exp, ok := experimentalV4Build91ExpectedForImage(imageName)
	if !ok {
		t.Fatalf("Build98 has no expected telemetry for %s", imageName)
	}
	if r.geometryEvals != exp.geometryEvals || r.bank != exp.bank || r.qualEvals != exp.qualEvals || r.qualified != exp.qualified || r.decodeCandidates != exp.decode || r.listFrames != exp.frames || r.physicalDecode != exp.physical || r.speculativeDecode != exp.speculative || r.foldScores != exp.foldScores || r.blockReads != exp.blockReads || r.blockSuccess != exp.blockSuccess || r.blockFailed != exp.blockFailed {
		t.Fatalf("Build98 %s %s telemetry mismatch: %+v", imageName, r.variant, r)
	}
	if r.variant == "build93" {
		if r.gen3FetchFoldScores != 0 || r.gen3FetchBlockReads != 0 || r.gen3FetchBlockSuccess != 0 || r.gen3FetchBlockFailed != 0 {
			t.Fatalf("Build98 %s baseline unexpectedly reports Gen3 fetch telemetry: %+v", imageName, r)
		}
	} else {
		if r.gen3FetchFoldScores <= 0 || r.gen3FetchBlockReads <= 0 || r.gen3FetchBlockReads != r.gen3FetchBlockSuccess+r.gen3FetchBlockFailed {
			t.Fatalf("Build98 %s candidate Gen3 fetch telemetry invalid: %+v", imageName, r)
		}
	}
	wantWorkers := experimentalV4PhoneBuild90QualificationWorkerCount(exp.bank)
	if r.qualificationWorkers != wantWorkers || r.qualificationTasks != exp.bank {
		t.Fatalf("Build98 %s %s qualification scheduling mismatch workers=%d want=%d tasks=%d want=%d", imageName, r.variant, r.qualificationWorkers, wantWorkers, r.qualificationTasks, exp.bank)
	}
	if exp.accept {
		if r.err != nil || !r.authenticated || string(r.payload) != expectedPayload {
			t.Fatalf("Build98 %s %s accept mismatch: err=%v authenticated=%t payload=%q", imageName, r.variant, r.err, r.authenticated, string(r.payload))
		}
		if r.physicalProfiles != 2086 || r.physicalListFrames != 2135407 || r.protectedSamplingPasses != 696 || r.protectedProfileDerivations != 2086 {
			t.Fatalf("Build98 %s %s protected telemetry mismatch: %+v", imageName, r.variant, r)
		}
	} else {
		if r.err == nil || r.authenticated || len(r.payload) != 0 {
			t.Fatalf("Build98 %s %s reject mismatch: err=%v authenticated=%t payload-bytes=%d", imageName, r.variant, r.err, r.authenticated, len(r.payload))
		}
		if r.physicalProfiles != 0 || r.physicalListFrames != 0 || r.protectedSamplingPasses != 0 || r.protectedProfileDerivations != 0 {
			t.Fatalf("Build98 %s %s unexpectedly reached protected decode: %+v", imageName, r.variant, r)
		}
	}
}

func experimentalV4Build98AssertPairEqual(t *testing.T, imageName string, a, b experimentalV4Build98PairedRun) {
	t.Helper()
	if a.geometryEvals != b.geometryEvals || a.bank != b.bank || a.qualEvals != b.qualEvals || a.qualified != b.qualified || a.decodeCandidates != b.decodeCandidates || a.listFrames != b.listFrames || a.physicalDecode != b.physicalDecode || a.speculativeDecode != b.speculativeDecode || a.physicalProfiles != b.physicalProfiles || a.physicalListFrames != b.physicalListFrames || a.protectedSamplingPasses != b.protectedSamplingPasses || a.protectedProfileDerivations != b.protectedProfileDerivations || a.authenticated != b.authenticated || a.foldScores != b.foldScores || a.blockReads != b.blockReads || a.blockSuccess != b.blockSuccess || a.blockFailed != b.blockFailed {
		t.Fatalf("Build98 %s semantic telemetry differs between %s and %s", imageName, a.variant, b.variant)
	}
	if (a.err == nil) != (b.err == nil) || !bytes.Equal(a.payload, b.payload) {
		t.Fatalf("Build98 %s outcome differs between %s and %s", imageName, a.variant, b.variant)
	}
	if a.err == nil && !reflect.DeepEqual(a.info, b.info) {
		t.Fatalf("Build98 %s extract info differs between %s and %s", imageName, a.variant, b.variant)
	}
}

func experimentalV4Build98RowFields(cycle int, imageName, order string, position int, r experimentalV4Build98PairedRun) []string {
	return []string{
		strconv.Itoa(cycle), imageName, order, strconv.Itoa(position), r.variant, r.outcome,
		strconv.FormatInt(r.wall.Milliseconds(), 10), strconv.FormatInt(r.total.Milliseconds(), 10),
		strconv.FormatInt(r.geometry.Milliseconds(), 10), strconv.FormatInt(r.plane.Milliseconds(), 10),
		strconv.FormatInt(r.freeze.Milliseconds(), 10), strconv.FormatInt(r.prefix1.Milliseconds(), 10),
		strconv.FormatInt(r.gen2.Milliseconds(), 10), strconv.FormatInt(r.gen3.Milliseconds(), 10), strconv.FormatInt(r.gen4.Milliseconds(), 10),
		strconv.FormatInt(r.qualification.Milliseconds(), 10), strconv.FormatInt(r.decodeWall.Milliseconds(), 10),
		strconv.Itoa(r.geometryEvals), strconv.Itoa(r.bank), strconv.Itoa(r.qualEvals), strconv.Itoa(r.qualified),
		strconv.Itoa(r.decodeCandidates), strconv.Itoa(r.listFrames), strconv.Itoa(r.physicalDecode), strconv.Itoa(r.speculativeDecode),
		strconv.Itoa(r.physicalProfiles), strconv.Itoa(r.physicalListFrames), strconv.Itoa(r.protectedSamplingPasses), strconv.Itoa(r.protectedProfileDerivations),
		strconv.FormatBool(r.authenticated), strconv.Itoa(r.qualificationWorkers), strconv.Itoa(r.qualificationTasks), strconv.Itoa(r.foldScores),
		strconv.Itoa(r.blockReads), strconv.Itoa(r.blockSuccess), strconv.Itoa(r.blockFailed),
		strconv.Itoa(r.gen3FetchFoldScores), strconv.Itoa(r.gen3FetchBlockReads), strconv.Itoa(r.gen3FetchBlockSuccess), strconv.Itoa(r.gen3FetchBlockFailed),
	}
}

func experimentalV4Build98PrintRow(cycle int, imageName, order string, position int, r experimentalV4Build98PairedRun) {
	fmt.Printf("build98-row\t%s\n", strings.Join(experimentalV4Build98RowFields(cycle, imageName, order, position, r), "\t"))
}

func TestExperimentalV4Build98PairedRowSchema(t *testing.T) {
	r := experimentalV4Build98PairedRun{variant: "build98", outcome: "reject", wall: time.Millisecond}
	fields := experimentalV4Build98RowFields(2, "phone-b-angle.jpg", "98-93", 1, r)
	if len(fields) != 40 {
		t.Fatalf("Build98 paired row fields=%d want=40", len(fields))
	}
	if fields[0] != "2" || fields[1] != "phone-b-angle.jpg" || fields[2] != "98-93" || fields[4] != "build98" || fields[39] != "0" {
		t.Fatalf("Build98 paired row schema mismatch: %q", fields)
	}
}

func TestExperimentalV4Build98PairedRetainedProfile(t *testing.T) {
	acquisitionDir := os.Getenv("PIXSEAL_BUILD98_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build98 private paired profile requires PIXSEAL_BUILD98_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD98_KEY"))
	if len(key) < 8 {
		t.Fatal("Build98 PIXSEAL_BUILD98_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD98_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid Build98 width: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD98_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid Build98 height: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD98_MESSAGE_B")
	cooldownMS, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD98_COOLDOWN_MS"))
	if err != nil || cooldownMS < 0 {
		t.Fatalf("invalid Build98 cooldown: %v", err)
	}
	cooldown := time.Duration(cooldownMS) * time.Millisecond

	type scenario struct {
		cycle int
		image string
		order string
	}
	scenarios := []scenario{
		{1, "phone-b-mild.jpg", "93-98"},
		{1, "phone-b-angle.jpg", "93-98"},
		{2, "phone-b-angle.jpg", "98-93"},
		{2, "phone-b-mild.jpg", "98-93"},
	}
	for _, s := range scenarios {
		path := filepath.Join(acquisitionDir, s.image)
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("open %s: %v", path, err)
		}
		src, err := jpeglegacy.Decode(f)
		_ = f.Close()
		if err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
		prepStarted := time.Now()
		work, boundary, hyp, _, _, _ := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
		prepElapsed := time.Since(prepStarted)
		if len(hyp) >= experimentalV4PhoneEnsembleSize {
			t.Fatalf("Build98 %s unexpectedly accepted before deep recovery", s.image)
		}
		fmt.Printf("build98-prepare\t%d\t%s\t%s\t%d\t%d\t%d\n", s.cycle, s.image, s.order, prepElapsed.Milliseconds(), work.Bounds().Dx(), work.Bounds().Dy())

		variants := []string{"build93", "build98"}
		if s.order == "98-93" {
			variants = []string{"build98", "build93"}
		}
		runs := make([]experimentalV4Build98PairedRun, 0, 2)
		for pos, variant := range variants {
			if cooldown > 0 {
				time.Sleep(cooldown)
			}
			runtime.GC()
			var r experimentalV4Build98PairedRun
			if variant == "build93" {
				r = experimentalV4Build98RunFromBuild93(work, boundary, key, cw, ch)
			} else {
				r = experimentalV4Build98RunFromBuild98(work, boundary, key, cw, ch)
			}
			experimentalV4Build98AssertExpected(t, s.image, r, expectedPayload)
			experimentalV4Build98PrintRow(s.cycle, s.image, s.order, pos+1, r)
			runs = append(runs, r)
		}
		experimentalV4Build98AssertPairEqual(t, s.image, runs[0], runs[1])
	}
}
