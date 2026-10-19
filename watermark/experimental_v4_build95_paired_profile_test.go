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

type experimentalV4Build95PairedRun struct {
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
	luminancePrep               time.Duration
	luminanceBytes              int64
	luminanceUsed               bool
	luminanceFallback           bool
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
	payload                     []byte
	info                        ExperimentalV4ExtractInfo
	err                         error
}

func experimentalV4Build95RunFromBuild93(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) experimentalV4Build95PairedRun {
	started := time.Now()
	payload, info, telemetry, err := experimentalV4PhoneBuild93Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	gp := telemetry.Build84GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	return experimentalV4Build95PairedRun{
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

func experimentalV4Build95RunFromBuild95(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) experimentalV4Build95PairedRun {
	started := time.Now()
	payload, info, telemetry, err := experimentalV4PhoneBuild95Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	gp := telemetry.Build95GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	return experimentalV4Build95PairedRun{
		variant: "build95", outcome: map[bool]string{true: "accept", false: "reject"}[err == nil], wall: wall,
		total: telemetry.TotalElapsed, geometry: telemetry.GeometryElapsed, plane: telemetry.PlanePrepElapsed,
		freeze: gp.FreezeElapsed, prefix1: gp.Prefix1WallElapsed, gen2: gp.Gen2WallElapsed, gen3: gp.Gen3WallElapsed, gen4: gp.Gen4WallElapsed,
		luminancePrep: gp.LuminancePlanePrepElapsed, luminanceBytes: gp.LuminancePlaneBytes, luminanceUsed: gp.LuminancePlaneUsed, luminanceFallback: gp.LuminancePlaneFallback,
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

func experimentalV4Build95AssertExpected(t *testing.T, imageName string, r experimentalV4Build95PairedRun, expectedPayload string, work image.Image) {
	t.Helper()
	exp, ok := experimentalV4Build91ExpectedForImage(imageName)
	if !ok {
		t.Fatalf("Build95 has no expected telemetry for %s", imageName)
	}
	if r.geometryEvals != exp.geometryEvals || r.bank != exp.bank || r.qualEvals != exp.qualEvals || r.qualified != exp.qualified || r.decodeCandidates != exp.decode || r.listFrames != exp.frames || r.physicalDecode != exp.physical || r.speculativeDecode != exp.speculative || r.foldScores != exp.foldScores || r.blockReads != exp.blockReads || r.blockSuccess != exp.blockSuccess || r.blockFailed != exp.blockFailed {
		t.Fatalf("Build95 %s %s telemetry mismatch: %+v", imageName, r.variant, r)
	}
	wantWorkers := experimentalV4PhoneBuild90QualificationWorkerCount(exp.bank)
	if r.qualificationWorkers != wantWorkers || r.qualificationTasks != exp.bank {
		t.Fatalf("Build95 %s %s qualification scheduling mismatch workers=%d want=%d tasks=%d want=%d", imageName, r.variant, r.qualificationWorkers, wantWorkers, r.qualificationTasks, exp.bank)
	}
	if exp.accept {
		if r.err != nil || !r.authenticated || string(r.payload) != expectedPayload {
			t.Fatalf("Build95 %s %s accept mismatch: err=%v authenticated=%t payload=%q", imageName, r.variant, r.err, r.authenticated, string(r.payload))
		}
		if r.physicalProfiles != 2086 || r.physicalListFrames != 2135407 || r.protectedSamplingPasses != 696 || r.protectedProfileDerivations != 2086 {
			t.Fatalf("Build95 %s %s protected telemetry mismatch: %+v", imageName, r.variant, r)
		}
	} else {
		if r.err == nil || r.authenticated || len(r.payload) != 0 {
			t.Fatalf("Build95 %s %s reject mismatch: err=%v authenticated=%t payload-bytes=%d", imageName, r.variant, r.err, r.authenticated, len(r.payload))
		}
		if r.physicalProfiles != 0 || r.physicalListFrames != 0 || r.protectedSamplingPasses != 0 || r.protectedProfileDerivations != 0 {
			t.Fatalf("Build95 %s %s unexpectedly reached protected decode: %+v", imageName, r.variant, r)
		}
	}
	if r.variant == "build95" {
		wantBytes := int64(work.Bounds().Dx()) * int64(work.Bounds().Dy()) * 8
		if !r.luminanceUsed || r.luminanceFallback || r.luminanceBytes != wantBytes || wantBytes > experimentalV4PhoneBuild95MaxLuminancePlaneBytes || r.luminancePrep <= 0 {
			t.Fatalf("Build95 %s luminance plane mismatch: bytes=%d want=%d used=%t fallback=%t prep=%s cap=%d", imageName, r.luminanceBytes, wantBytes, r.luminanceUsed, r.luminanceFallback, r.luminancePrep, experimentalV4PhoneBuild95MaxLuminancePlaneBytes)
		}
	} else if r.luminanceUsed || r.luminanceFallback || r.luminanceBytes != 0 || r.luminancePrep != 0 {
		t.Fatalf("Build93 row unexpectedly carries Build95 luminance telemetry: %+v", r)
	}
}

func experimentalV4Build95AssertPairEqual(t *testing.T, imageName string, a, b experimentalV4Build95PairedRun) {
	t.Helper()
	if a.geometryEvals != b.geometryEvals || a.bank != b.bank || a.qualEvals != b.qualEvals || a.qualified != b.qualified || a.decodeCandidates != b.decodeCandidates || a.listFrames != b.listFrames || a.physicalDecode != b.physicalDecode || a.speculativeDecode != b.speculativeDecode || a.physicalProfiles != b.physicalProfiles || a.physicalListFrames != b.physicalListFrames || a.protectedSamplingPasses != b.protectedSamplingPasses || a.protectedProfileDerivations != b.protectedProfileDerivations || a.authenticated != b.authenticated || a.foldScores != b.foldScores || a.blockReads != b.blockReads || a.blockSuccess != b.blockSuccess || a.blockFailed != b.blockFailed {
		t.Fatalf("Build95 %s semantic telemetry differs between %s and %s", imageName, a.variant, b.variant)
	}
	if (a.err == nil) != (b.err == nil) || !bytes.Equal(a.payload, b.payload) {
		t.Fatalf("Build95 %s outcome differs between %s and %s", imageName, a.variant, b.variant)
	}
	if a.err == nil && !reflect.DeepEqual(a.info, b.info) {
		t.Fatalf("Build95 %s extract info differs between %s and %s", imageName, a.variant, b.variant)
	}
}

func experimentalV4Build95RowFields(cycle int, imageName, order string, position int, r experimentalV4Build95PairedRun) []string {
	return []string{
		strconv.Itoa(cycle), imageName, order, strconv.Itoa(position), r.variant, r.outcome,
		strconv.FormatInt(r.wall.Milliseconds(), 10), strconv.FormatInt(r.total.Milliseconds(), 10),
		strconv.FormatInt(r.geometry.Milliseconds(), 10), strconv.FormatInt(r.plane.Milliseconds(), 10),
		strconv.FormatInt(r.freeze.Milliseconds(), 10), strconv.FormatInt(r.prefix1.Milliseconds(), 10),
		strconv.FormatInt(r.gen2.Milliseconds(), 10), strconv.FormatInt(r.gen3.Milliseconds(), 10),
		strconv.FormatInt(r.gen4.Milliseconds(), 10), strconv.FormatInt(r.luminancePrep.Milliseconds(), 10),
		strconv.FormatInt(r.luminanceBytes, 10), strconv.FormatBool(r.luminanceUsed), strconv.FormatBool(r.luminanceFallback),
		strconv.FormatInt(r.qualification.Milliseconds(), 10), strconv.FormatInt(r.decodeWall.Milliseconds(), 10),
		strconv.Itoa(r.geometryEvals), strconv.Itoa(r.bank), strconv.Itoa(r.qualEvals), strconv.Itoa(r.qualified),
		strconv.Itoa(r.decodeCandidates), strconv.Itoa(r.listFrames), strconv.Itoa(r.physicalDecode),
		strconv.Itoa(r.speculativeDecode), strconv.Itoa(r.physicalProfiles), strconv.Itoa(r.physicalListFrames),
		strconv.Itoa(r.protectedSamplingPasses), strconv.Itoa(r.protectedProfileDerivations), strconv.FormatBool(r.authenticated),
		strconv.Itoa(r.qualificationWorkers), strconv.Itoa(r.qualificationTasks), strconv.Itoa(r.foldScores),
		strconv.Itoa(r.blockReads), strconv.Itoa(r.blockSuccess), strconv.Itoa(r.blockFailed),
	}
}

func experimentalV4Build95PrintRow(cycle int, imageName, order string, position int, r experimentalV4Build95PairedRun) {
	fmt.Printf("build95-row\t%s\n", strings.Join(experimentalV4Build95RowFields(cycle, imageName, order, position, r), "\t"))
}

func TestExperimentalV4Build95PairedRowSchema(t *testing.T) {
	r := experimentalV4Build95PairedRun{
		variant: "build95", outcome: "accept", wall: 1 * time.Millisecond,
		luminanceBytes: 127475200, luminanceUsed: true,
	}
	fields := experimentalV4Build95RowFields(2, "phone-b-mild.jpg", "95-93", 1, r)
	if len(fields) != 40 {
		t.Fatalf("Build95 paired row fields=%d want=40", len(fields))
	}
	if fields[0] != "2" || fields[1] != "phone-b-mild.jpg" || fields[2] != "95-93" || fields[4] != "build95" || fields[16] != "127475200" || fields[17] != "true" || fields[39] != "0" {
		t.Fatalf("Build95 paired row schema mismatch: %q", fields)
	}
	for i, field := range fields {
		if strings.ContainsAny(field, "\t\n\r") {
			t.Fatalf("Build95 paired row field %d contains a delimiter: %q", i, field)
		}
	}
}

func TestExperimentalV4Build95PairedPhysicalProfile(t *testing.T) {
	acquisitionDir := os.Getenv("PIXSEAL_BUILD95_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build95 private paired profile requires PIXSEAL_BUILD95_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD95_KEY"))
	if len(key) < 8 {
		t.Fatal("Build95 PIXSEAL_BUILD95_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD95_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid Build95 width: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD95_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid Build95 height: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD95_MESSAGE_B")
	cooldownMS, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD95_COOLDOWN_MS"))
	if err != nil || cooldownMS < 0 {
		t.Fatalf("invalid Build95 cooldown: %v", err)
	}
	cooldown := time.Duration(cooldownMS) * time.Millisecond

	type scenario struct {
		cycle int
		image string
		order string
	}
	scenarios := []scenario{
		{cycle: 1, image: "phone-b-mild.jpg", order: "93-95"},
		{cycle: 1, image: "phone-b-angle.jpg", order: "93-95"},
		{cycle: 2, image: "phone-b-angle.jpg", order: "95-93"},
		{cycle: 2, image: "phone-b-mild.jpg", order: "95-93"},
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
			t.Fatalf("Build95 %s unexpectedly accepted before deep recovery", s.image)
		}
		fmt.Printf("build95-prepare\t%d\t%s\t%s\t%d\t%d\t%d\n", s.cycle, s.image, s.order, prepElapsed.Milliseconds(), work.Bounds().Dx(), work.Bounds().Dy())

		variants := []string{"build93", "build95"}
		if s.order == "95-93" {
			variants = []string{"build95", "build93"}
		}
		runs := make([]experimentalV4Build95PairedRun, 0, 2)
		for pos, variant := range variants {
			if cooldown > 0 {
				time.Sleep(cooldown)
			}
			runtime.GC()
			var r experimentalV4Build95PairedRun
			if variant == "build93" {
				r = experimentalV4Build95RunFromBuild93(work, boundary, key, cw, ch)
			} else {
				r = experimentalV4Build95RunFromBuild95(work, boundary, key, cw, ch)
			}
			experimentalV4Build95AssertExpected(t, s.image, r, expectedPayload, work)
			experimentalV4Build95PrintRow(s.cycle, s.image, s.order, pos+1, r)
			runs = append(runs, r)
		}
		experimentalV4Build95AssertPairEqual(t, s.image, runs[0], runs[1])
	}
}
