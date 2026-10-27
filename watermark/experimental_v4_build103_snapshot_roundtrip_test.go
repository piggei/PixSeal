package watermark

import (
	"crypto/sha256"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExperimentalV4Build103SnapshotBitsAndMutationGuard(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	root.quad[0].X = math.Copysign(0, -1)
	root.h.h[1] = math.Float64frombits(0x7ff8000000000042)
	root.proposal = math.Copysign(0, -1)
	root.validation = math.Inf(-1)
	root.detection = ExperimentalV4PilotDetection{Available: true, BlockSizePixels: 8, PhaseX: 1, PhaseY: 2, OriginXBlocks: 3, OriginYBlocks: 4, Score: math.Copysign(0, -1), RunnerUpScore: math.Inf(-1), Margin: math.NaN(), VisiblePilotPositions: 99, PilotSamples: 88}
	root.warp = &experimentalV4PhoneResidualWarp{width: 123, height: 321, controls: 6, rmsPixels: math.Copysign(0, -1), maxPixels: 2}
	root.warp.dx[2] = 1.25
	root.warp.dy[3] = math.Float64frombits(0x7ff8000000000043)
	root.residual = experimentalV4PhoneResidualInfo{Applied: true, Fitted: true, Controls: 6, RMSPixels: math.Copysign(0, -1), ProposalBefore: 1, ProposalAfter: 2, ValidationBefore: 3, ValidationAfter: 4, FullPilotScore: 5, FullPilotMargin: 6, FullOriginXBlocks: 7, FullOriginYBlocks: 8, HypothesesEvaluated: 9}
	root.build43Pair = "fixture"
	root.build43PairRank = 3
	f := experimentalV4Build102Frozen{plane: plane, pilot: pilot, cw: w, ch: h, anchor: anchor, tasks: []experimentalV4PhoneHypothesis{root, root}, evals: 19, seeds: 2, workers: 2}
	f.profile.Gen3FetchFoldScores = 7
	f.profile.Gen3FetchBlockReads = 8
	f.profile.Gen3FetchBlockSuccess = 6
	f.profile.Gen3FetchBlockFailed = 2
	f.profile.Gen2WallElapsed = 17 * time.Millisecond
	f.profile.PairScoring.Workers = 6
	f.freeze.PairScoreTasks = 6
	s := experimentalV4Build103Capture(f)
	path := filepath.Join(t.TempDir(), "snapshot.gob")
	if err := experimentalV4Build103Write(path, s); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(b))
	loaded, err := experimentalV4Build103Read(path, hash)
	if err != nil {
		t.Fatal(err)
	}
	g, err := loaded.frozen()
	if err != nil {
		t.Fatal(err)
	}
	if experimentalV4Build102Hash(f.tasks) != experimentalV4Build102Hash(g.tasks) || experimentalV4Build103Meta(f) != experimentalV4Build103Meta(g) {
		t.Fatal("snapshot lost task/private field/float bits or metadata")
	}
	if _, err := experimentalV4Build103Read(path, "wrong"); err == nil {
		t.Fatal("accepted wrong snapshot checksum")
	}
	loaded.Tasks[0].H[0]++
	if _, err := loaded.frozen(); err == nil {
		t.Fatal("accepted mutated homography bits")
	}
	loaded = s
	loaded.RGB = append([]byte(nil), s.RGB...)
	loaded.RGB[0] ^= 1
	if _, err := loaded.frozen(); err == nil {
		t.Fatal("accepted mutated raster")
	}
	loaded = s
	loaded.CW++
	if _, err := loaded.frozen(); err == nil {
		t.Fatal("accepted mutated metadata")
	}
}
