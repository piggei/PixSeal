package watermark

import (
	"bytes"
	"image"
	"math"
	"reflect"
	"sync"
	"testing"
)

type experimentalV4Build93Fixture struct {
	image   image.Image
	plane   *pixelPlane
	hyp     experimentalV4PhoneHypothesis
	width   int
	height  int
	key     []byte
	payload []byte
	profile Profile
}

func experimentalV4Build93EmbeddedFixture(tb testing.TB, profile Profile) experimentalV4Build93Fixture {
	tb.Helper()
	key := []byte("Build93-Public-Key")
	var payload []byte
	switch profile {
	case ProfileRobust:
		payload = []byte("build93-robust")
	case ProfileBalanced:
		payload = bytes.Repeat([]byte{'B'}, 29)
	case ProfileCapacity:
		payload = bytes.Repeat([]byte{'C'}, 61)
	default:
		tb.Fatalf("unsupported Build93 profile %q", profile)
	}
	source := testImage(592, 512)
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, Options{Strength: 48, Profile: profile})
	if err != nil {
		tb.Fatal(err)
	}
	w, h := marked.Bounds().Dx(), marked.Bounds().Dy()
	anchor := [4]ImagePoint{{X: 0, Y: 0}, {X: float64(w - 1), Y: 0}, {X: 0, Y: float64(h - 1)}, {X: float64(w - 1), Y: float64(h - 1)}}
	hom, ok := experimentalV4PhoneQuadHomography(w, h, anchor)
	if !ok {
		tb.Fatal("Build93 identity homography")
	}
	detection := experimentalV4DetectPilotProjective(marked, experimentalV4Prototype2Candidate(), w, h, hom)
	return experimentalV4Build93Fixture{
		image: marked,
		plane: newPixelPlane(marked),
		hyp: experimentalV4PhoneHypothesis{
			quad:      anchor,
			h:         hom,
			detection: detection,
		},
		width: w, height: h, key: key, payload: payload, profile: profile,
	}
}

func experimentalV4Build93AssertDecodeEqual(t *testing.T, want experimentalV4PhoneBuild67DecodeResult, got experimentalV4PhoneBuild93DecodeResult) {
	t.Helper()
	g := got.experimentalV4PhoneBuild67DecodeResult
	if (want.err == nil) != (g.err == nil) {
		t.Fatalf("decode outcome mismatch historical=%v build93=%v", want.err, g.err)
	}
	if !bytes.Equal(want.payload, g.payload) || want.info != g.info || want.profiles != g.profiles || want.frames != g.frames || math.Float64bits(want.maxConfidence) != math.Float64bits(g.maxConfidence) {
		t.Fatalf("decode semantic mismatch\nhistorical=%+v\nbuild93=%+v", want.experimentalV4PhoneBuild66DecodeResult, g.experimentalV4PhoneBuild66DecodeResult)
	}
}

func TestExperimentalV4Build93PublicEntrypointsShareQualifiedValidation(t *testing.T) {
	validImage := image.NewRGBA(image.Rect(0, 0, 16, 16))
	cases := []struct {
		name   string
		src    image.Image
		key    []byte
		width  int
		height int
	}{
		{name: "nil-image", src: nil, key: []byte("Build93-Key"), width: 1632, height: 1632},
		{name: "short-key", src: validImage, key: []byte("short"), width: 1632, height: 1632},
		{name: "invalid-canonical", src: validImage, key: []byte("Build93-Key"), width: 288, height: 248},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			productionPayload, productionInfo, productionPhone, productionErr := ExperimentalV4ExtractPhone(tc.src, tc.key, tc.width, tc.height)
			aliasPayload, aliasInfo, aliasPhone, aliasErr := ExperimentalV4ExtractPhoneBuild93Candidate(tc.src, tc.key, tc.width, tc.height)
			if !bytes.Equal(productionPayload, aliasPayload) || !reflect.DeepEqual(productionInfo, aliasInfo) || !reflect.DeepEqual(productionPhone, aliasPhone) {
				t.Fatalf("Build93 public entry points diverged on %s", tc.name)
			}
			if (productionErr == nil) != (aliasErr == nil) {
				t.Fatalf("Build93 public entry point error outcome diverged on %s: production=%v alias=%v", tc.name, productionErr, aliasErr)
			}
			if productionErr != nil && productionErr.Error() != aliasErr.Error() {
				t.Fatalf("Build93 public entry point error text diverged on %s: production=%q alias=%q", tc.name, productionErr, aliasErr)
			}
		})
	}
}

func TestExperimentalV4Build93MarginsMatchHistoricalBitExactly(t *testing.T) {
	_, _, plane, w, h, anchor, _ := build84SyntheticFixture(t)
	quads := map[string][4]ImagePoint{
		"front": anchor,
		"mild":  {{X: 3, Y: 2}, {X: float64(w - 5), Y: 5}, {X: 6, Y: float64(h - 4)}, {X: float64(w - 7), Y: float64(h - 6)}},
		"angle": {{X: 14, Y: 5}, {X: float64(w - 18), Y: 11}, {X: 4, Y: float64(h - 15)}, {X: float64(w - 9), Y: float64(h - 7)}},
	}
	for name, quad := range quads {
		hom, ok := experimentalV4PhoneQuadHomography(w, h, quad)
		if !ok {
			t.Fatalf("%s homography", name)
		}
		residues, ok := experimentalV4PhoneBuild93SampleProtectedResidues(plane, w, h, hom)
		if !ok {
			t.Fatalf("%s residue sampling", name)
		}
		for _, spec := range v3Profiles {
			want, wantConfidence, wantOK := experimentalV4ReadProtectedProjectiveMargins(plane, w, h, hom, spec.codedBits)
			got, gotConfidence, gotOK := experimentalV4PhoneBuild93MarginsFromResidues(&residues, spec.codedBits)
			if gotOK != wantOK {
				t.Fatalf("%s/%s ok=%t/%t", name, spec.profile, gotOK, wantOK)
			}
			if math.Float64bits(gotConfidence) != math.Float64bits(wantConfidence) {
				t.Fatalf("%s/%s confidence bits=%x/%x", name, spec.profile, math.Float64bits(gotConfidence), math.Float64bits(wantConfidence))
			}
			if len(got) != len(want) {
				t.Fatalf("%s/%s margins len=%d/%d", name, spec.profile, len(got), len(want))
			}
			for i := range got {
				if math.Float64bits(got[i]) != math.Float64bits(want[i]) {
					t.Fatalf("%s/%s margin[%d] bits=%x/%x", name, spec.profile, i, math.Float64bits(got[i]), math.Float64bits(want[i]))
				}
			}
		}
	}
}

func TestExperimentalV4Build93DecodeMatchesHistoricalProfiles(t *testing.T) {
	for profileIndex, profile := range []Profile{ProfileRobust, ProfileBalanced, ProfileCapacity} {
		t.Run(string(profile), func(t *testing.T) {
			f := experimentalV4Build93EmbeddedFixture(t, profile)
			want := experimentalV4PhoneBuild67DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
			got := experimentalV4PhoneBuild93DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
			experimentalV4Build93AssertDecodeEqual(t, want, got)
			if got.err != nil || got.info.Profile != profile || !bytes.Equal(got.payload, f.payload) {
				t.Fatalf("Build93 %s failed authenticated fixture: result=%+v", profile, got.experimentalV4PhoneBuild66DecodeResult)
			}
			wantDerivations := profileIndex + 1
			if got.protectedSamplingPasses != 1 || got.protectedProfileDerivations != wantDerivations || got.profiles != wantDerivations {
				t.Fatalf("Build93 %s sampling=%d derivations=%d profiles=%d want 1/%d/%d", profile, got.protectedSamplingPasses, got.protectedProfileDerivations, got.profiles, wantDerivations, wantDerivations)
			}
		})
	}
}

func TestExperimentalV4Build93DecodeMatchesHistoricalWrongKey(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	wrongKey := []byte("Build93-Wrong-Key")
	want := experimentalV4PhoneBuild67DecodeSingle(f.plane, f.hyp, wrongKey, f.width, f.height)
	got := experimentalV4PhoneBuild93DecodeSingle(f.plane, f.hyp, wrongKey, f.width, f.height)
	experimentalV4Build93AssertDecodeEqual(t, want, got)
	if got.err == nil || got.protectedSamplingPasses != 1 || got.protectedProfileDerivations != len(v3Profiles) || got.profiles != len(v3Profiles) {
		t.Fatalf("wrong-key Build93 result=%+v sampling=%d derivations=%d", got.experimentalV4PhoneBuild66DecodeResult, got.protectedSamplingPasses, got.protectedProfileDerivations)
	}
}

func TestExperimentalV4Build93SamplingFailureMatchesHistorical(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	want := experimentalV4PhoneBuild67DecodeSingle(nil, f.hyp, f.key, f.width, f.height)
	got := experimentalV4PhoneBuild93DecodeSingle(nil, f.hyp, f.key, f.width, f.height)
	experimentalV4Build93AssertDecodeEqual(t, want, got)
	if got.protectedSamplingPasses != 0 || got.protectedProfileDerivations != 0 {
		t.Fatalf("failed sampling counters=%d/%d", got.protectedSamplingPasses, got.protectedProfileDerivations)
	}
}

func TestExperimentalV4Build93ParallelDecodeRaceFixture(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	const workers = 8
	results := make([]experimentalV4PhoneBuild93DecodeResult, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		i := i
		go func() {
			defer wg.Done()
			results[i] = experimentalV4PhoneBuild93DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
		}()
	}
	wg.Wait()
	for i := range results {
		if results[i].err != nil || !bytes.Equal(results[i].payload, f.payload) || results[i].info.Profile != ProfileCapacity || results[i].protectedSamplingPasses != 1 || results[i].protectedProfileDerivations != 3 {
			t.Fatalf("worker %d result=%+v sampling=%d derivations=%d", i, results[i].experimentalV4PhoneBuild66DecodeResult, results[i].protectedSamplingPasses, results[i].protectedProfileDerivations)
		}
	}
}

func TestExperimentalV4Build93PublicBenchmarkFixture(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	want := experimentalV4PhoneBuild67DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
	got := experimentalV4PhoneBuild93DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
	experimentalV4Build93AssertDecodeEqual(t, want, got)
	if got.err != nil || got.protectedSamplingPasses != 1 || got.protectedProfileDerivations != 3 {
		t.Fatalf("benchmark fixture failed: %+v sampling=%d derivations=%d", got.experimentalV4PhoneBuild66DecodeResult, got.protectedSamplingPasses, got.protectedProfileDerivations)
	}
}

var experimentalV4Build93BenchmarkFrames int
var experimentalV4Build93BenchmarkConfidence float64
var experimentalV4Build93BenchmarkPayloadByte byte

func benchmarkExperimentalV4Build93Single(b *testing.B, reuse bool) {
	f := experimentalV4Build93EmbeddedFixture(b, ProfileCapacity)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if reuse {
			got := experimentalV4PhoneBuild93DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
			experimentalV4Build93BenchmarkFrames = got.frames
			experimentalV4Build93BenchmarkConfidence = got.maxConfidence
			if len(got.payload) > 0 {
				experimentalV4Build93BenchmarkPayloadByte = got.payload[0]
			}
		} else {
			got := experimentalV4PhoneBuild67DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
			experimentalV4Build93BenchmarkFrames = got.frames
			experimentalV4Build93BenchmarkConfidence = got.maxConfidence
			if len(got.payload) > 0 {
				experimentalV4Build93BenchmarkPayloadByte = got.payload[0]
			}
		}
	}
}

func BenchmarkExperimentalV4Build93HistoricalSingleCapacity(b *testing.B) {
	benchmarkExperimentalV4Build93Single(b, false)
}

func BenchmarkExperimentalV4Build93ReuseSingleCapacity(b *testing.B) {
	benchmarkExperimentalV4Build93Single(b, true)
}

func benchmarkExperimentalV4Build93Batch8(b *testing.B, reuse bool) {
	f := experimentalV4Build93EmbeddedFixture(b, ProfileCapacity)
	const workers = 8
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(workers)
		frames := make([]int, workers)
		confidence := make([]float64, workers)
		first := make([]byte, workers)
		for i := 0; i < workers; i++ {
			i := i
			go func() {
				defer wg.Done()
				if reuse {
					got := experimentalV4PhoneBuild93DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
					frames[i] = got.frames
					confidence[i] = got.maxConfidence
					if len(got.payload) > 0 {
						first[i] = got.payload[0]
					}
				} else {
					got := experimentalV4PhoneBuild67DecodeSingle(f.plane, f.hyp, f.key, f.width, f.height)
					frames[i] = got.frames
					confidence[i] = got.maxConfidence
					if len(got.payload) > 0 {
						first[i] = got.payload[0]
					}
				}
			}()
		}
		wg.Wait()
		experimentalV4Build93BenchmarkFrames = frames[workers-1]
		experimentalV4Build93BenchmarkConfidence = confidence[workers-1]
		experimentalV4Build93BenchmarkPayloadByte = first[workers-1]
	}
}

func BenchmarkExperimentalV4Build93HistoricalBatch8Capacity(b *testing.B) {
	benchmarkExperimentalV4Build93Batch8(b, false)
}

func BenchmarkExperimentalV4Build93ReuseBatch8Capacity(b *testing.B) {
	benchmarkExperimentalV4Build93Batch8(b, true)
}

func TestExperimentalV4Build93MarginFailureOnSmallCanonicalMatchesHistorical(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileRobust)
	const smallW = (experimentalV4TileWidthBlocks - 1) * blockSize
	const smallH = (experimentalV4TileHeightBlocks - 1) * blockSize
	residues, ok := experimentalV4PhoneBuild93SampleProtectedResidues(f.plane, smallW, smallH, f.hyp.h)
	if ok {
		t.Fatal("Build93 unexpectedly sampled undersized canonical grid")
	}
	for _, spec := range v3Profiles {
		want, wc, wok := experimentalV4ReadProtectedProjectiveMargins(f.plane, smallW, smallH, f.hyp.h, spec.codedBits)
		got, gc, gok := experimentalV4PhoneBuild93MarginsFromResidues(&residues, spec.codedBits)
		if wok || want != nil || wc != 0 {
			t.Fatalf("historical %s undersized result=%v/%v/%v", spec.profile, want, wc, wok)
		}
		// Derivation from an empty residue set must also reject; the decode path
		// never calls it after sample failure, but this keeps helper behavior safe.
		if gok || got != nil || gc != 0 {
			t.Fatalf("Build93 %s empty-residue result=%v/%v/%v", spec.profile, got, gc, gok)
		}
	}
}

func TestExperimentalV4Build93RecoveryMatchesBuild90PublicCarrier(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	boundary := PrintBoundaryEstimate{
		Detected: true, Confidence: 1,
		TopLeft: f.hyp.quad[0], TopRight: f.hyp.quad[1],
		BottomLeft: f.hyp.quad[2], BottomRight: f.hyp.quad[3],
	}
	work := f.image
	wantPayload, wantInfo, want, wantErr := experimentalV4PhoneBuild90Recover(work, boundary, f.key, f.width, f.height)
	gotPayload, gotInfo, got, gotErr := experimentalV4PhoneBuild93Recover(work, boundary, f.key, f.width, f.height)
	if (wantErr == nil) != (gotErr == nil) || !bytes.Equal(wantPayload, gotPayload) || wantInfo != gotInfo {
		t.Fatalf("recovery outcome mismatch Build90 err=%v info=%+v Build93 err=%v info=%+v", wantErr, wantInfo, gotErr, gotInfo)
	}
	if want.GeometryEvaluations != got.GeometryEvaluations || want.BankCandidates != got.BankCandidates || want.QualificationEvaluations != got.QualificationEvaluations || want.QualifiedCandidates != got.QualifiedCandidates || want.DecodeCandidatesTried != got.DecodeCandidatesTried || want.ProfilesTried != got.ProfilesTried || want.ListFramesTried != got.ListFramesTried || want.Authenticated != got.Authenticated || want.Profile != got.Profile || want.PhysicalDecodeCandidates != got.PhysicalDecodeCandidates || want.PhysicalProfilesTried != got.PhysicalProfilesTried || want.PhysicalListFramesTried != got.PhysicalListFramesTried || want.QualificationWorkers != got.QualificationWorkers || want.QualificationTasks != got.QualificationTasks || want.DecodeWorkers != got.DecodeWorkers {
		t.Fatalf("recovery semantic telemetry mismatch\nBuild90=%+v\nBuild93=%+v", want, got.experimentalV4PhoneBuild90RecoveryTelemetry)
	}
	if got.PhysicalDecodeCandidates > 0 && got.ProtectedSamplingPasses != got.PhysicalDecodeCandidates {
		t.Fatalf("Build93 sampling passes=%d physical candidates=%d", got.ProtectedSamplingPasses, got.PhysicalDecodeCandidates)
	}
	if got.ProtectedProfileDerivations != got.PhysicalProfilesTried {
		t.Fatalf("Build93 derivations=%d physical profiles=%d", got.ProtectedProfileDerivations, got.PhysicalProfilesTried)
	}
}
