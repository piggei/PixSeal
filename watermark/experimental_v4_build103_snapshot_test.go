package watermark

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"image"
	"io"
	"math"
	"os"
)

// All task floating-point fields use explicit IEEE bits. A normal struct
// encoder may omit negative zero or unexported geometry fields.
type experimentalV4Build103Task struct {
	Quad               [8]uint64
	H                  [9]uint64
	Scores             [2]uint64
	DetectionAvailable bool
	DetectionInts      [7]int
	DetectionFloats    [3]uint64
	Warp               *experimentalV4Build103Warp
	ResidualBools      [2]bool
	ResidualInts       [4]int
	ResidualFloats     [7]uint64
	Pair               string
	Rank               int
}
type experimentalV4Build103Warp struct {
	Bits     [16]uint64
	Controls int
}
type experimentalV4Build103Geometry struct {
	Base experimentalV4PhoneBuild84GeometryTelemetry
	Gen3 [4]int
}
type experimentalV4Build103Snapshot struct {
	Version                                    int
	Image                                      string
	Bounds                                     image.Rectangle
	RGB                                        []byte
	CW, CH                                     int
	Anchor                                     [8]uint64
	Tasks                                      []experimentalV4Build103Task
	PilotName                                  string
	PilotPositions                             [experimentalV4PilotCount]int
	PilotSigns                                 [experimentalV4PilotCount]int8
	Evals, Seeds, Workers                      int
	Prefix                                     experimentalV4Build103Geometry
	Pair                                       experimentalV4PhoneBuild100PairTelemetry
	Freeze                                     experimentalV4PhoneBuild75FreezeTelemetry
	TaskHash, PlaneHash, MetaHash              string
	BankHash                                   string
	Bank, EvalsTotal, SeedsTotal, WorkersTotal int
	Expected                                   experimentalV4Build103Geometry
	ExpectedFreeze                             experimentalV4PhoneBuild75FreezeTelemetry
}

func experimentalV4Build103PointsBits(q [4]ImagePoint) (b [8]uint64) {
	for i, p := range q {
		b[2*i] = math.Float64bits(p.X)
		b[2*i+1] = math.Float64bits(p.Y)
	}
	return
}
func experimentalV4Build103Points(b [8]uint64) (q [4]ImagePoint) {
	for i := range q {
		q[i] = ImagePoint{X: math.Float64frombits(b[2*i]), Y: math.Float64frombits(b[2*i+1])}
	}
	return
}
func experimentalV4Build103PackTask(h experimentalV4PhoneHypothesis) experimentalV4Build103Task {
	d := h.detection
	r := h.residual
	t := experimentalV4Build103Task{Quad: experimentalV4Build103PointsBits(h.quad), Scores: [2]uint64{math.Float64bits(h.proposal), math.Float64bits(h.validation)}, DetectionAvailable: d.Available, DetectionInts: [7]int{d.BlockSizePixels, d.PhaseX, d.PhaseY, d.OriginXBlocks, d.OriginYBlocks, d.VisiblePilotPositions, d.PilotSamples}, DetectionFloats: [3]uint64{math.Float64bits(d.Score), math.Float64bits(d.RunnerUpScore), math.Float64bits(d.Margin)}, ResidualBools: [2]bool{r.Applied, r.Fitted}, ResidualInts: [4]int{r.Controls, r.FullOriginXBlocks, r.FullOriginYBlocks, r.HypothesesEvaluated}, Pair: h.build43Pair, Rank: h.build43PairRank}
	for i, v := range h.h.h {
		t.H[i] = math.Float64bits(v)
	}
	for i, v := range []float64{r.RMSPixels, r.ProposalBefore, r.ProposalAfter, r.ValidationBefore, r.ValidationAfter, r.FullPilotScore, r.FullPilotMargin} {
		t.ResidualFloats[i] = math.Float64bits(v)
	}
	if h.warp != nil {
		w := h.warp
		t.Warp = &experimentalV4Build103Warp{Controls: w.controls}
		v := append([]float64{w.width, w.height}, w.dx[:]...)
		v = append(v, w.dy[:]...)
		v = append(v, w.rmsPixels, w.maxPixels)
		for i, x := range v {
			t.Warp.Bits[i] = math.Float64bits(x)
		}
	}
	return t
}
func experimentalV4Build103UnpackTask(t experimentalV4Build103Task) experimentalV4PhoneHypothesis {
	v := math.Float64frombits
	h := experimentalV4PhoneHypothesis{quad: experimentalV4Build103Points(t.Quad), proposal: v(t.Scores[0]), validation: v(t.Scores[1]), build43Pair: t.Pair, build43PairRank: t.Rank}
	for i, x := range t.H {
		h.h.h[i] = v(x)
	}
	d := t.DetectionInts
	f := t.DetectionFloats
	h.detection = ExperimentalV4PilotDetection{Available: t.DetectionAvailable, BlockSizePixels: d[0], PhaseX: d[1], PhaseY: d[2], OriginXBlocks: d[3], OriginYBlocks: d[4], VisiblePilotPositions: d[5], PilotSamples: d[6], Score: v(f[0]), RunnerUpScore: v(f[1]), Margin: v(f[2])}
	r := t.ResidualInts
	x := t.ResidualFloats
	h.residual = experimentalV4PhoneResidualInfo{Applied: t.ResidualBools[0], Fitted: t.ResidualBools[1], Controls: r[0], FullOriginXBlocks: r[1], FullOriginYBlocks: r[2], HypothesesEvaluated: r[3], RMSPixels: v(x[0]), ProposalBefore: v(x[1]), ProposalAfter: v(x[2]), ValidationBefore: v(x[3]), ValidationAfter: v(x[4]), FullPilotScore: v(x[5]), FullPilotMargin: v(x[6])}
	if t.Warp != nil {
		b := t.Warp.Bits
		h.warp = &experimentalV4PhoneResidualWarp{width: v(b[0]), height: v(b[1]), controls: t.Warp.Controls, rmsPixels: v(b[14]), maxPixels: v(b[15])}
		for i := 0; i < 6; i++ {
			h.warp.dx[i] = v(b[i+2])
			h.warp.dy[i] = v(b[i+8])
		}
	}
	return h
}
func experimentalV4Build103Geo(g experimentalV4PhoneBuild98GeometryTelemetry) experimentalV4Build103Geometry {
	return experimentalV4Build103Geometry{g.experimentalV4PhoneBuild84GeometryTelemetry, [4]int{g.Gen3FetchFoldScores, g.Gen3FetchBlockReads, g.Gen3FetchBlockSuccess, g.Gen3FetchBlockFailed}}
}
func (g experimentalV4Build103Geometry) unpack() experimentalV4PhoneBuild98GeometryTelemetry {
	return experimentalV4PhoneBuild98GeometryTelemetry{g.Base, g.Gen3[0], g.Gen3[1], g.Gen3[2], g.Gen3[3]}
}
func experimentalV4Build103Meta(f experimentalV4Build102Frozen) string {
	return experimentalV4Build102Hash(struct {
		Bounds                image.Rectangle
		Pilot                 experimentalV4PilotCandidate
		CW, CH                int
		Anchor                [4]ImagePoint
		Evals, Seeds, Workers int
		Profile               experimentalV4PhoneBuild100GeometryTelemetry
		Freeze                experimentalV4PhoneBuild75FreezeTelemetry
	}{f.plane.bounds, f.pilot, f.cw, f.ch, f.anchor, f.evals, f.seeds, f.workers, f.profile, f.freeze})
}
func experimentalV4Build103Capture(f experimentalV4Build102Frozen) experimentalV4Build103Snapshot {
	s := experimentalV4Build103Snapshot{Version: 1, Bounds: f.plane.bounds, RGB: f.plane.rgb, CW: f.cw, CH: f.ch, Anchor: experimentalV4Build103PointsBits(f.anchor), PilotName: f.pilot.name, PilotPositions: f.pilot.positions, PilotSigns: f.pilot.signs, Evals: f.evals, Seeds: f.seeds, Workers: f.workers, Prefix: experimentalV4Build103Geo(f.profile.experimentalV4PhoneBuild98GeometryTelemetry), Pair: f.profile.PairScoring, Freeze: f.freeze, TaskHash: experimentalV4Build102Hash(f.tasks), PlaneHash: fmt.Sprintf("%x", sha256.Sum256(f.plane.rgb)), MetaHash: experimentalV4Build103Meta(f)}
	s.Tasks = make([]experimentalV4Build103Task, len(f.tasks))
	for i, h := range f.tasks {
		s.Tasks[i] = experimentalV4Build103PackTask(h)
	}
	return s
}
func (s experimentalV4Build103Snapshot) frozen() (experimentalV4Build102Frozen, error) {
	var f experimentalV4Build102Frozen
	w, h := s.Bounds.Dx(), s.Bounds.Dy()
	if s.Version != 1 || w <= 0 || h <= 0 || w > 20000 || h > 20000 || int64(w)*int64(h)*3 != int64(len(s.RGB)) || s.CW <= 0 || s.CH <= 0 || len(s.Tasks) == 0 {
		return f, fmt.Errorf("invalid frozen snapshot schema/shape")
	}
	f = experimentalV4Build102Frozen{plane: &pixelPlane{bounds: s.Bounds, rgb: s.RGB}, pilot: experimentalV4PilotCandidate{name: s.PilotName, positions: s.PilotPositions, signs: s.PilotSigns}, cw: s.CW, ch: s.CH, anchor: experimentalV4Build103Points(s.Anchor), evals: s.Evals, seeds: s.Seeds, workers: s.Workers, profile: experimentalV4PhoneBuild100GeometryTelemetry{experimentalV4PhoneBuild98GeometryTelemetry: s.Prefix.unpack(), PairScoring: s.Pair}, freeze: s.Freeze}
	f.tasks = make([]experimentalV4PhoneHypothesis, len(s.Tasks))
	for i, t := range s.Tasks {
		f.tasks[i] = experimentalV4Build103UnpackTask(t)
	}
	if experimentalV4Build102Hash(f.tasks) != s.TaskHash || fmt.Sprintf("%x", sha256.Sum256(f.plane.rgb)) != s.PlaneHash || experimentalV4Build103Meta(f) != s.MetaHash {
		return f, fmt.Errorf("frozen snapshot fingerprint differs")
	}
	if experimentalV4Build102Hash(f.pilot) != experimentalV4Build102Hash(experimentalV4Prototype2Candidate()) {
		return f, fmt.Errorf("snapshot pilot lock differs")
	}
	return f, nil
}
func experimentalV4Build103Write(path string, s experimentalV4Build103Snapshot) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	err = gob.NewEncoder(f).Encode(s)
	closed := f.Close()
	if err != nil {
		return err
	}
	return closed
}
func experimentalV4Build103Read(path, want string) (experimentalV4Build103Snapshot, error) {
	var s experimentalV4Build103Snapshot
	f, err := os.Open(path)
	if err != nil {
		return s, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return s, err
	}
	if info.Size() > 512<<20 {
		return s, fmt.Errorf("snapshot exceeds study bound")
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return s, err
	}
	if fmt.Sprintf("%x", h.Sum(nil)) != want {
		return s, fmt.Errorf("snapshot file checksum differs")
	}
	if _, err := f.Seek(0, 0); err != nil {
		return s, err
	}
	err = gob.NewDecoder(f).Decode(&s)
	return s, err
}
