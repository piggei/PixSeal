package watermark

import (
	"fmt"
	"math"
	"testing"
)

type experimentalV4Build90Fixture struct {
	pilot       experimentalV4PilotCandidate
	plane       *pixelPlane
	width       int
	height      int
	highPass    []experimentalV4PhoneHypothesis
	earlyReject []experimentalV4PhoneHypothesis
}

func experimentalV4Build90QualificationFixture(tb testing.TB, highPassCount, earlyRejectCount int) experimentalV4Build90Fixture {
	tb.Helper()
	pilot := experimentalV4Prototype2Candidate()
	source := testImage(592, 512)
	carrier, err := experimentalV4RenderSyntheticCarrier(source, pilot, 48, experimentalV4SyntheticDataSeed)
	if err != nil {
		tb.Fatal(err)
	}
	width, height := carrier.Bounds().Dx(), carrier.Bounds().Dy()
	plane := newPixelPlane(carrier)
	quad := [4]ImagePoint{
		{X: 0, Y: 0},
		{X: float64(width - 1), Y: 0},
		{X: 0, Y: float64(height - 1)},
		{X: float64(width - 1), Y: float64(height - 1)},
	}
	h, ok := experimentalV4PhoneQuadHomography(width, height, quad)
	if !ok {
		tb.Fatal("Build90 fixture homography rejected")
	}
	proposal, _ := experimentalV4PhoneBuild41FoldScore(plane, pilot, width, height, h, 0, true, false)
	base := experimentalV4PhoneHypothesis{quad: quad, h: h, proposal: proposal}
	qualified, evals, accepted := experimentalV4PhoneBuild68Qualify(plane, pilot, width, height, base, 0)
	if !accepted || evals != 2 {
		tb.Fatalf("Build90 high-pass fixture does not reach full detector: proposal=%.9f validation=%.9f pilot=%.9f margin=%.9f evals=%d ok=%t", proposal, qualified.validation, qualified.detection.Score, qualified.detection.Margin, evals, accepted)
	}

	highPass := make([]experimentalV4PhoneHypothesis, highPassCount)
	for i := range highPass {
		highPass[i] = base
		highPass[i].build43Pair = fmt.Sprintf("build90-public-high-pass-%04d", i)
		highPass[i].build43PairRank = i + 1
	}
	earlyReject := make([]experimentalV4PhoneHypothesis, earlyRejectCount)
	for i := range earlyReject {
		earlyReject[i] = base
		earlyReject[i].proposal = experimentalV4PhoneProposalFloor - 1
		earlyReject[i].build43Pair = fmt.Sprintf("build90-public-early-reject-%04d", i)
		earlyReject[i].build43PairRank = i + 1
	}
	return experimentalV4Build90Fixture{pilot: pilot, plane: plane, width: width, height: height, highPass: highPass, earlyReject: earlyReject}
}

func experimentalV4Build90SerialResults(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, bank []experimentalV4PhoneHypothesis, heldout int) []experimentalV4PhoneBuild90QualificationResult {
	results := make([]experimentalV4PhoneBuild90QualificationResult, len(bank))
	for i := range bank {
		q, n, ok := experimentalV4PhoneBuild68Qualify(plane, pilot, cw, ch, bank[i], heldout)
		results[i] = experimentalV4PhoneBuild90QualificationResult{hypothesis: q, evals: n, ok: ok}
	}
	return results
}

func experimentalV4Build90HypothesisEqualBits(a, b experimentalV4PhoneHypothesis) bool {
	if a.quad != b.quad || a.h != b.h || math.Float64bits(a.proposal) != math.Float64bits(b.proposal) || math.Float64bits(a.validation) != math.Float64bits(b.validation) || a.warp != b.warp || a.residual != b.residual || a.build43Pair != b.build43Pair || a.build43PairRank != b.build43PairRank {
		return false
	}
	ad, bd := a.detection, b.detection
	return ad.Available == bd.Available && ad.BlockSizePixels == bd.BlockSizePixels && ad.PhaseX == bd.PhaseX && ad.PhaseY == bd.PhaseY && ad.OriginXBlocks == bd.OriginXBlocks && ad.OriginYBlocks == bd.OriginYBlocks && math.Float64bits(ad.Score) == math.Float64bits(bd.Score) && math.Float64bits(ad.RunnerUpScore) == math.Float64bits(bd.RunnerUpScore) && math.Float64bits(ad.Margin) == math.Float64bits(bd.Margin) && ad.VisiblePilotPositions == bd.VisiblePilotPositions && ad.PilotSamples == bd.PilotSamples
}

func experimentalV4Build90AssertResultsEqual(t *testing.T, serial, parallel []experimentalV4PhoneBuild90QualificationResult) {
	t.Helper()
	if len(serial) != len(parallel) {
		t.Fatalf("result length %d/%d", len(parallel), len(serial))
	}
	for i := range serial {
		if parallel[i].evals != serial[i].evals || parallel[i].ok != serial[i].ok || !experimentalV4Build90HypothesisEqualBits(parallel[i].hypothesis, serial[i].hypothesis) {
			t.Fatalf("qualification result %d differs: parallel evals=%d ok=%t serial evals=%d ok=%t", i, parallel[i].evals, parallel[i].ok, serial[i].evals, serial[i].ok)
		}
	}
}

func TestExperimentalV4Build90WorkerCount(t *testing.T) {
	if experimentalV4PhoneBuild90QualificationWorkerCount(0) != 0 {
		t.Fatal("zero candidates must use zero workers")
	}
	if got := experimentalV4PhoneBuild90QualificationWorkerCount(1); got != 1 {
		t.Fatalf("one candidate workers=%d want 1", got)
	}
	if got := experimentalV4PhoneBuild90QualificationWorkerCount(3); got < 1 || got > 3 {
		t.Fatalf("three candidate workers=%d", got)
	}
}

func TestExperimentalV4Build90ParallelQualificationMatchesSerial(t *testing.T) {
	f := experimentalV4Build90QualificationFixture(t, 17, 257)
	for name, bank := range map[string][]experimentalV4PhoneHypothesis{"high-pass": f.highPass, "early-reject": f.earlyReject} {
		t.Run(name, func(t *testing.T) {
			serial := experimentalV4Build90SerialResults(f.plane, f.pilot, f.width, f.height, bank, 0)
			parallel, workers := experimentalV4PhoneBuild90EvaluateQualification(f.plane, f.pilot, f.width, f.height, bank, 0)
			if workers != experimentalV4PhoneBuild90QualificationWorkerCount(len(bank)) {
				t.Fatalf("workers=%d", workers)
			}
			experimentalV4Build90AssertResultsEqual(t, serial, parallel)

			serialQualified, serialEvals := experimentalV4PhoneBuild90CommitQualification(serial)
			parallelQualified, parallelEvals := experimentalV4PhoneBuild90CommitQualification(parallel)
			if parallelEvals != serialEvals || len(parallelQualified) != len(serialQualified) {
				t.Fatalf("commit evals=%d/%d qualified=%d/%d", parallelEvals, serialEvals, len(parallelQualified), len(serialQualified))
			}
			for i := range serialQualified {
				if !experimentalV4Build90HypothesisEqualBits(parallelQualified[i], serialQualified[i]) {
					t.Fatalf("ordered qualified candidate %d differs", i)
				}
				if parallelQualified[i].build43PairRank != serialQualified[i].build43PairRank {
					t.Fatalf("ordered rank %d: %d/%d", i, parallelQualified[i].build43PairRank, serialQualified[i].build43PairRank)
				}
			}
		})
	}
}

func TestExperimentalV4Build90ApplyTelemetryPreservesBuild84View(t *testing.T) {
	r := experimentalV4PhoneBuild90RecoveryTelemetry{}
	r.Attempted = true
	r.QualificationWorkers = 8
	r.QualificationTasks = 937
	r.QualificationEvaluations = 1872
	r.Build84GeometryProfile = experimentalV4PhoneBuild84GeometryTelemetry{
		Prefix1Workers: 8, Gen2Workers: 8, Gen3Workers: 8, Gen4Workers: 8,
		Gen2Tasks: 93, Gen3Tasks: 275, Gen4Tasks: 593,
		FetchFoldScores: 16496, FetchBlockReads: 21114880, FetchBlockSuccess: 21114880,
	}
	var p ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild90ApplyTelemetry(&p, r)
	if !p.Build90Attempted || !p.Build84Attempted || !p.Build76Attempted {
		t.Fatalf("attempted chain lost: %+v", p)
	}
	if p.Build90QualificationWorkers != 8 || p.Build90QualificationTasks != 937 || p.Build90QualificationEvaluations != 1872 {
		t.Fatalf("Build90 telemetry %+v", p)
	}
	if p.Build84FetchFoldScores != 16496 || p.Build84FetchBlockReads != 21114880 || p.Build84FetchBlockSuccess != 21114880 || p.Build84FetchBlockFailed != 0 {
		t.Fatalf("Build84 view changed: %+v", p)
	}
}

func TestExperimentalV4Build90PublicBenchmarkFixture(t *testing.T) {
	f := experimentalV4Build90QualificationFixture(t, 8, 64)
	serial := experimentalV4Build90SerialResults(f.plane, f.pilot, f.width, f.height, f.highPass, 0)
	parallel, workers := experimentalV4PhoneBuild90EvaluateQualification(f.plane, f.pilot, f.width, f.height, f.highPass, 0)
	experimentalV4Build90AssertResultsEqual(t, serial, parallel)
	if workers < 1 || len(f.highPass) != 8 || len(f.earlyReject) != 64 {
		t.Fatalf("fixture workers=%d high=%d reject=%d", workers, len(f.highPass), len(f.earlyReject))
	}
}

var experimentalV4Build90BenchmarkEvals int
var experimentalV4Build90BenchmarkQualified int

func benchmarkExperimentalV4Build90Serial(b *testing.B, bankKind string) {
	f := experimentalV4Build90QualificationFixture(b, 24, 512)
	bank := f.highPass
	if bankKind == "early-reject" {
		bank = f.earlyReject
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := experimentalV4Build90SerialResults(f.plane, f.pilot, f.width, f.height, bank, 0)
		qualified, evals := experimentalV4PhoneBuild90CommitQualification(results)
		experimentalV4Build90BenchmarkEvals = evals
		experimentalV4Build90BenchmarkQualified = len(qualified)
	}
}

func benchmarkExperimentalV4Build90Parallel(b *testing.B, bankKind string) {
	f := experimentalV4Build90QualificationFixture(b, 24, 512)
	bank := f.highPass
	if bankKind == "early-reject" {
		bank = f.earlyReject
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qualified, evals, _ := experimentalV4PhoneBuild90QualifyBank(f.plane, f.pilot, f.width, f.height, bank, 0)
		experimentalV4Build90BenchmarkEvals = evals
		experimentalV4Build90BenchmarkQualified = len(qualified)
	}
}

func BenchmarkExperimentalV4Build90SerialHighPass(b *testing.B) {
	benchmarkExperimentalV4Build90Serial(b, "high-pass")
}

func BenchmarkExperimentalV4Build90ParallelHighPass(b *testing.B) {
	benchmarkExperimentalV4Build90Parallel(b, "high-pass")
}

func BenchmarkExperimentalV4Build90SerialEarlyReject(b *testing.B) {
	benchmarkExperimentalV4Build90Serial(b, "early-reject")
}

func BenchmarkExperimentalV4Build90ParallelEarlyReject(b *testing.B) {
	benchmarkExperimentalV4Build90Parallel(b, "early-reject")
}
