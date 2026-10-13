package watermark

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	jpeglegacy "github.com/pj/pixseal/internal/jpeglegacy"
)

type experimentalV4Build91Expected struct {
	geometryEvals int
	bank          int
	qualEvals     int
	qualified     int
	decode        int
	frames        int
	physical      int
	speculative   int
	foldScores    int
	blockReads    int
	blockSuccess  int
	blockFailed   int
	accept        bool
}

type experimentalV4Build91Run struct {
	variant              string
	outcome              string
	wall                 time.Duration
	total                time.Duration
	geometry             time.Duration
	plane                time.Duration
	freeze               time.Duration
	prefix1              time.Duration
	gen2                 time.Duration
	gen3                 time.Duration
	gen4                 time.Duration
	qualification        time.Duration
	decodeWall           time.Duration
	geometryEvals        int
	bank                 int
	qualEvals            int
	qualified            int
	decodeCandidates     int
	listFrames           int
	physicalDecode       int
	speculativeDecode    int
	authenticated        bool
	qualificationWorkers int
	qualificationTasks   int
	foldScores           int
	blockReads           int
	blockSuccess         int
	blockFailed          int
	payload              []byte
	info                 ExperimentalV4ExtractInfo
	err                  error
}

func experimentalV4Build91ExpectedForImage(name string) (experimentalV4Build91Expected, bool) {
	switch name {
	case "phone-b-mild.jpg":
		return experimentalV4Build91Expected{geometryEvals: 79259, bank: 937, qualEvals: 1872, qualified: 935, decode: 691, frames: 2120047, physical: 696, speculative: 5, foldScores: 16496, blockReads: 21114880, blockSuccess: 21114880, blockFailed: 0, accept: true}, true
	case "phone-b-angle.jpg":
		return experimentalV4Build91Expected{geometryEvals: 334857, bank: 6198, qualEvals: 6198, qualified: 0, decode: 0, frames: 0, physical: 0, speculative: 0, foldScores: 118004, blockReads: 151045120, blockSuccess: 149370664, blockFailed: 1674456, accept: false}, true
	default:
		return experimentalV4Build91Expected{}, false
	}
}

func experimentalV4Build91RunFromBuild84(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) experimentalV4Build91Run {
	started := time.Now()
	payload, info, telemetry, err := experimentalV4PhoneBuild84Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	gp := telemetry.Build84GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	return experimentalV4Build91Run{
		variant: "build84", outcome: map[bool]string{true: "accept", false: "reject"}[err == nil], wall: wall,
		total: telemetry.TotalElapsed, geometry: telemetry.GeometryElapsed, plane: telemetry.PlanePrepElapsed,
		freeze: gp.FreezeElapsed, prefix1: gp.Prefix1WallElapsed, gen2: gp.Gen2WallElapsed, gen3: gp.Gen3WallElapsed, gen4: gp.Gen4WallElapsed,
		qualification: telemetry.QualificationElapsed, decodeWall: telemetry.DecodeWallElapsed,
		geometryEvals: telemetry.GeometryEvaluations, bank: telemetry.BankCandidates, qualEvals: telemetry.QualificationEvaluations,
		qualified: telemetry.QualifiedCandidates, decodeCandidates: telemetry.DecodeCandidatesTried, listFrames: telemetry.ListFramesTried,
		physicalDecode: physical, speculativeDecode: physical - telemetry.DecodeCandidatesTried, authenticated: telemetry.Authenticated,
		qualificationWorkers: 1, qualificationTasks: telemetry.BankCandidates,
		foldScores: gp.FetchFoldScores, blockReads: gp.FetchBlockReads, blockSuccess: gp.FetchBlockSuccess, blockFailed: gp.FetchBlockFailed,
		payload: payload, info: info, err: err,
	}
}

func experimentalV4Build91RunFromBuild90(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) experimentalV4Build91Run {
	started := time.Now()
	payload, info, telemetry, err := experimentalV4PhoneBuild90Recover(work, boundary, key, cw, ch)
	wall := time.Since(started)
	gp := telemetry.Build84GeometryProfile
	physical := telemetry.PhysicalDecodeCandidates
	return experimentalV4Build91Run{
		variant: "build90", outcome: map[bool]string{true: "accept", false: "reject"}[err == nil], wall: wall,
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

func experimentalV4Build91AssertExpected(t *testing.T, imageName string, r experimentalV4Build91Run, expectedPayload string) {
	t.Helper()
	exp, ok := experimentalV4Build91ExpectedForImage(imageName)
	if !ok {
		t.Fatalf("Build91 has no expected telemetry for %s", imageName)
	}
	if r.geometryEvals != exp.geometryEvals || r.bank != exp.bank || r.qualEvals != exp.qualEvals || r.qualified != exp.qualified || r.decodeCandidates != exp.decode || r.listFrames != exp.frames || r.physicalDecode != exp.physical || r.speculativeDecode != exp.speculative || r.foldScores != exp.foldScores || r.blockReads != exp.blockReads || r.blockSuccess != exp.blockSuccess || r.blockFailed != exp.blockFailed {
		t.Fatalf("Build91 %s %s telemetry mismatch: %+v", imageName, r.variant, r)
	}
	if exp.accept {
		if r.err != nil || !r.authenticated || string(r.payload) != expectedPayload {
			t.Fatalf("Build91 %s %s accept mismatch: err=%v authenticated=%t payload=%q", imageName, r.variant, r.err, r.authenticated, string(r.payload))
		}
	} else if r.err == nil || r.authenticated || len(r.payload) != 0 {
		t.Fatalf("Build91 %s %s reject mismatch: err=%v authenticated=%t payload-bytes=%d", imageName, r.variant, r.err, r.authenticated, len(r.payload))
	}
	if r.variant == "build84" {
		if r.qualificationWorkers != 1 || r.qualificationTasks != exp.bank {
			t.Fatalf("Build91 Build84 qualification scheduling mismatch: workers=%d tasks=%d", r.qualificationWorkers, r.qualificationTasks)
		}
	} else if r.variant == "build90" {
		wantWorkers := experimentalV4PhoneBuild90QualificationWorkerCount(exp.bank)
		if r.qualificationWorkers != wantWorkers || r.qualificationTasks != exp.bank {
			t.Fatalf("Build91 Build90 qualification scheduling mismatch: workers=%d want=%d tasks=%d", r.qualificationWorkers, wantWorkers, r.qualificationTasks)
		}
	}
}

func experimentalV4Build91AssertPairEqual(t *testing.T, imageName string, a, b experimentalV4Build91Run) {
	t.Helper()
	if a.geometryEvals != b.geometryEvals || a.bank != b.bank || a.qualEvals != b.qualEvals || a.qualified != b.qualified || a.decodeCandidates != b.decodeCandidates || a.listFrames != b.listFrames || a.physicalDecode != b.physicalDecode || a.speculativeDecode != b.speculativeDecode || a.authenticated != b.authenticated || a.foldScores != b.foldScores || a.blockReads != b.blockReads || a.blockSuccess != b.blockSuccess || a.blockFailed != b.blockFailed {
		t.Fatalf("Build91 %s semantic telemetry differs between %s and %s", imageName, a.variant, b.variant)
	}
	if (a.err == nil) != (b.err == nil) || !bytes.Equal(a.payload, b.payload) {
		t.Fatalf("Build91 %s outcome differs between %s and %s", imageName, a.variant, b.variant)
	}
	if a.err == nil && !reflect.DeepEqual(a.info, b.info) {
		t.Fatalf("Build91 %s extract info differs between %s and %s", imageName, a.variant, b.variant)
	}
}

func experimentalV4Build91PrintRow(cycle int, imageName, order string, position int, r experimentalV4Build91Run) {
	fmt.Printf("build91-row\t%d\t%s\t%s\t%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%t\t%d\t%d\t%d\t%d\t%d\t%d\n",
		cycle, imageName, order, position, r.variant, r.outcome,
		r.wall.Milliseconds(), r.total.Milliseconds(), r.geometry.Milliseconds(), r.plane.Milliseconds(), r.freeze.Milliseconds(), r.prefix1.Milliseconds(), r.gen2.Milliseconds(), r.gen3.Milliseconds(), r.gen4.Milliseconds(), r.qualification.Milliseconds(), r.decodeWall.Milliseconds(),
		r.geometryEvals, r.bank, r.qualEvals, r.qualified, r.decodeCandidates, r.listFrames, r.physicalDecode, r.speculativeDecode, r.authenticated,
		r.qualificationWorkers, r.qualificationTasks, r.foldScores, r.blockReads, r.blockSuccess, r.blockFailed)
}

func TestExperimentalV4Build91ExpectedTelemetryTable(t *testing.T) {
	mild, ok := experimentalV4Build91ExpectedForImage("phone-b-mild.jpg")
	if !ok || !mild.accept || mild.bank != 937 || mild.qualEvals != 1872 || mild.decode != 691 || mild.frames != 2120047 {
		t.Fatalf("Build91 B/mild expected table invalid: %+v", mild)
	}
	angle, ok := experimentalV4Build91ExpectedForImage("phone-b-angle.jpg")
	if !ok || angle.accept || angle.bank != 6198 || angle.qualEvals != 6198 || angle.qualified != 0 {
		t.Fatalf("Build91 B/angle expected table invalid: %+v", angle)
	}
}

func TestExperimentalV4Build91PairedPhysicalProfile(t *testing.T) {
	acquisitionDir := os.Getenv("PIXSEAL_BUILD91_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build91 private paired profile requires PIXSEAL_BUILD91_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD91_KEY"))
	if len(key) < 8 {
		t.Fatal("Build91 PIXSEAL_BUILD91_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD91_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD91_WIDTH: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD91_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD91_HEIGHT: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD91_MESSAGE_B")
	if expectedPayload == "" {
		t.Fatal("PIXSEAL_BUILD91_MESSAGE_B is required")
	}
	cooldownMs := 5000
	if raw := os.Getenv("PIXSEAL_BUILD91_COOLDOWN_MS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 || v > 60000 {
			t.Fatalf("invalid PIXSEAL_BUILD91_COOLDOWN_MS=%q", raw)
		}
		cooldownMs = v
	}
	cooldown := time.Duration(cooldownMs) * time.Millisecond

	type scenario struct {
		cycle int
		image string
		order string
	}
	// Reverse both variant order and case order in cycle 2 to reduce systematic
	// drift from host temperature, scheduler state and long-run frequency changes.
	scenarios := []scenario{
		{cycle: 1, image: "phone-b-mild.jpg", order: "84-90"},
		{cycle: 1, image: "phone-b-angle.jpg", order: "84-90"},
		{cycle: 2, image: "phone-b-angle.jpg", order: "90-84"},
		{cycle: 2, image: "phone-b-mild.jpg", order: "90-84"},
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
			t.Fatalf("Build91 %s unexpectedly accepted before deep recovery", s.image)
		}
		fmt.Printf("build91-prepare\t%d\t%s\t%s\t%d\n", s.cycle, s.image, s.order, prepElapsed.Milliseconds())

		variants := []string{"build84", "build90"}
		if s.order == "90-84" {
			variants = []string{"build90", "build84"}
		}
		runs := make([]experimentalV4Build91Run, 0, 2)
		for pos, variant := range variants {
			if cooldown > 0 {
				time.Sleep(cooldown)
			}
			var r experimentalV4Build91Run
			if variant == "build84" {
				r = experimentalV4Build91RunFromBuild84(work, boundary, key, cw, ch)
			} else {
				r = experimentalV4Build91RunFromBuild90(work, boundary, key, cw, ch)
			}
			experimentalV4Build91AssertExpected(t, s.image, r, expectedPayload)
			experimentalV4Build91PrintRow(s.cycle, s.image, s.order, pos+1, r)
			runs = append(runs, r)
		}
		experimentalV4Build91AssertPairEqual(t, s.image, runs[0], runs[1])
	}
}
