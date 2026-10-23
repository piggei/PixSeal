package watermark

import (
	"bytes"
	"image"
	"reflect"
	"testing"
)

func TestExperimentalV4Build98SingleProbeMatchesBuild53Exactly(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	we, wi := experimentalV4PhoneBuild53SingleProbe(plane, pilot, w, h, anchor, root, 0)
	ge, gi, fetch := experimentalV4PhoneBuild98SingleProbeFetch(plane, pilot, w, h, anchor, root, 0)
	if ge != we || gi != wi {
		t.Fatalf("single mismatch evals=%d/%d improving=%d/%d", ge, we, gi, wi)
	}
	if fetch.FoldScores != ge || fetch.BlockReads <= 0 || fetch.BlockReads != fetch.BlockSuccess+fetch.BlockFailed {
		t.Fatalf("single fetch telemetry invalid: %+v", fetch)
	}
}

func TestExperimentalV4Build98PairStencilMatchesBuild53Exactly(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	want, we, wi := experimentalV4PhoneBuild53PairStencil(plane, pilot, w, h, anchor, root, 0)
	got, ge, gi, fetch := experimentalV4PhoneBuild98PairStencilFetch(plane, pilot, w, h, anchor, root, 0)
	if ge != we || gi != wi || !reflect.DeepEqual(got, want) {
		t.Fatalf("pair mismatch evals=%d/%d improving=%d/%d got=%d want=%d", ge, we, gi, wi, len(got), len(want))
	}
	if fetch.FoldScores != ge || fetch.BlockReads <= 0 || fetch.BlockReads != fetch.BlockSuccess+fetch.BlockFailed {
		t.Fatalf("pair fetch telemetry invalid: %+v", fetch)
	}
}

func TestExperimentalV4Build98SiblingStencilMatchesBuild55Exactly(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	want, we := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, w, h, anchor, root, 0)
	got, ge, fetch := experimentalV4PhoneBuild98SiblingStencilFetch(plane, pilot, w, h, anchor, root, 0)
	if ge != we || !reflect.DeepEqual(got, want) {
		t.Fatalf("sibling mismatch evals=%d/%d got=%d want=%d", ge, we, len(got), len(want))
	}
	if fetch.FoldScores != ge || fetch.BlockReads <= 0 || fetch.BlockReads != fetch.BlockSuccess+fetch.BlockFailed {
		t.Fatalf("sibling fetch telemetry invalid: %+v", fetch)
	}
}

func TestExperimentalV4Build98Generation3MatchesBuild73Exactly(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	want, we := experimentalV4PhoneBuild73Generation3(plane, pilot, w, h, anchor, root)
	got, ge, fetch := experimentalV4PhoneBuild98Generation3Fetch(plane, pilot, w, h, anchor, root)
	if ge != we || !reflect.DeepEqual(got, want) {
		t.Fatalf("gen3 mismatch evals=%d/%d got=%d want=%d", ge, we, len(got), len(want))
	}
	total := experimentalV4PhoneBuild98TotalFetch(fetch)
	if total.FoldScores != ge || total.BlockReads <= 0 || total.BlockReads != total.BlockSuccess+total.BlockFailed {
		t.Fatalf("gen3 fetch telemetry invalid: %+v total=%+v", fetch, total)
	}
}

func TestExperimentalV4Build98ParallelGeneration3RaceFixture(t *testing.T) {
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	const workers = 8
	errs := make(chan string, workers)
	for i := 0; i < workers; i++ {
		go func() {
			want, we := experimentalV4PhoneBuild73Generation3(plane, pilot, w, h, anchor, root)
			got, ge, _ := experimentalV4PhoneBuild98Generation3Fetch(plane, pilot, w, h, anchor, root)
			if ge != we || !reflect.DeepEqual(got, want) {
				errs <- "mismatch"
				return
			}
			errs <- ""
		}()
	}
	for i := 0; i < workers; i++ {
		if e := <-errs; e != "" {
			t.Fatal(e)
		}
	}
}

func TestExperimentalV4Build98BlindBankMatchesBuild84Exactly(t *testing.T) {
	_, carrier, _, w, h, _, _ := build84SyntheticFixture(t)
	boundary := PrintBoundaryEstimate{
		Detected:    true,
		Confidence:  1,
		TopLeft:     ImagePoint{X: 0, Y: 0},
		TopRight:    ImagePoint{X: float64(w - 1), Y: 0},
		BottomRight: ImagePoint{X: float64(w - 1), Y: float64(h - 1)},
		BottomLeft:  ImagePoint{X: 0, Y: float64(h - 1)},
	}
	want, wantEvals, wantSeeds, wantWorkers, wantProfile, _ := experimentalV4PhoneBuild84BlindBank(carrier, boundary, w, h)
	got, gotEvals, gotSeeds, gotWorkers, gotProfile, freeze := experimentalV4PhoneBuild98BlindBank(carrier, boundary, w, h)
	if gotEvals != wantEvals || gotSeeds != wantSeeds || gotWorkers != wantWorkers || !reflect.DeepEqual(got, want) {
		t.Fatalf("bank mismatch evals=%d/%d seeds=%d/%d workers=%d/%d bank=%d/%d", gotEvals, wantEvals, gotSeeds, wantSeeds, gotWorkers, wantWorkers, len(got), len(want))
	}
	if gotProfile.FetchFoldScores != wantProfile.FetchFoldScores || gotProfile.FetchBlockReads != wantProfile.FetchBlockReads || gotProfile.FetchBlockSuccess != wantProfile.FetchBlockSuccess || gotProfile.FetchBlockFailed != wantProfile.FetchBlockFailed {
		t.Fatalf("historical continuation4 telemetry changed: got=%+v want=%+v", gotProfile.experimentalV4PhoneBuild84GeometryTelemetry, wantProfile)
	}
	if gotProfile.Gen3Tasks > 0 {
		if gotProfile.Gen3FetchFoldScores <= 0 || gotProfile.Gen3FetchBlockReads <= 0 || gotProfile.Gen3FetchBlockReads != gotProfile.Gen3FetchBlockSuccess+gotProfile.Gen3FetchBlockFailed {
			t.Fatalf("Build98 Gen3 fetch telemetry invalid: %+v", gotProfile)
		}
	}
	if freeze.BasinTasks != 16 {
		t.Fatalf("freeze telemetry changed: %+v", freeze)
	}
}

func TestExperimentalV4Build98RecoveryMatchesBuild93PublicCarrier(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	boundary := PrintBoundaryEstimate{
		Detected:    true,
		Confidence:  1,
		TopLeft:     f.hyp.quad[0],
		TopRight:    f.hyp.quad[1],
		BottomLeft:  f.hyp.quad[2],
		BottomRight: f.hyp.quad[3],
	}
	wantPayload, wantInfo, wantTelemetry, wantErr := experimentalV4PhoneBuild93Recover(f.image, boundary, f.key, f.width, f.height)
	gotPayload, gotInfo, gotTelemetry, gotErr := experimentalV4PhoneBuild98Recover(f.image, boundary, f.key, f.width, f.height)
	if (wantErr == nil) != (gotErr == nil) || !bytes.Equal(gotPayload, wantPayload) || !reflect.DeepEqual(gotInfo, wantInfo) {
		t.Fatalf("recovery outcome mismatch Build93 err=%v info=%+v Build98 err=%v info=%+v", wantErr, wantInfo, gotErr, gotInfo)
	}
	if wantTelemetry.GeometryEvaluations != gotTelemetry.GeometryEvaluations || wantTelemetry.BankCandidates != gotTelemetry.BankCandidates || wantTelemetry.QualificationEvaluations != gotTelemetry.QualificationEvaluations || wantTelemetry.QualifiedCandidates != gotTelemetry.QualifiedCandidates || wantTelemetry.DecodeCandidatesTried != gotTelemetry.DecodeCandidatesTried || wantTelemetry.ListFramesTried != gotTelemetry.ListFramesTried || wantTelemetry.Authenticated != gotTelemetry.Authenticated || wantTelemetry.PhysicalDecodeCandidates != gotTelemetry.PhysicalDecodeCandidates || wantTelemetry.PhysicalProfilesTried != gotTelemetry.PhysicalProfilesTried || wantTelemetry.PhysicalListFramesTried != gotTelemetry.PhysicalListFramesTried || wantTelemetry.ProtectedSamplingPasses != gotTelemetry.ProtectedSamplingPasses || wantTelemetry.ProtectedProfileDerivations != gotTelemetry.ProtectedProfileDerivations || wantTelemetry.QualificationWorkers != gotTelemetry.QualificationWorkers || wantTelemetry.QualificationTasks != gotTelemetry.QualificationTasks {
		t.Fatalf("recovery semantic telemetry mismatch\nBuild93=%+v\nBuild98=%+v", wantTelemetry, gotTelemetry)
	}
	wantGeometry := wantTelemetry.Build84GeometryProfile
	gotGeometry := gotTelemetry.Build98GeometryProfile
	if wantGeometry.FetchFoldScores != gotGeometry.FetchFoldScores || wantGeometry.FetchBlockReads != gotGeometry.FetchBlockReads || wantGeometry.FetchBlockSuccess != gotGeometry.FetchBlockSuccess || wantGeometry.FetchBlockFailed != gotGeometry.FetchBlockFailed {
		t.Fatalf("historical continuation4 telemetry mismatch: want=%+v got=%+v", wantGeometry, gotGeometry)
	}
	if gotGeometry.Gen3Tasks > 0 {
		if gotGeometry.Gen3FetchFoldScores <= 0 || gotGeometry.Gen3FetchBlockReads <= 0 || gotGeometry.Gen3FetchBlockReads != gotGeometry.Gen3FetchBlockSuccess+gotGeometry.Gen3FetchBlockFailed {
			t.Fatalf("Build98 Gen3 fetch telemetry invalid: %+v", gotGeometry)
		}
	} else if gotGeometry.Gen3FetchFoldScores != 0 || gotGeometry.Gen3FetchBlockReads != 0 || gotGeometry.Gen3FetchBlockSuccess != 0 || gotGeometry.Gen3FetchBlockFailed != 0 {
		t.Fatalf("Build98 reports Gen3 fetch telemetry without Gen3 tasks: %+v", gotGeometry)
	}
}

func TestExperimentalV4Build98PublicEntrypointsShareQualifiedValidation(t *testing.T) {
	validImage := image.NewRGBA(image.Rect(0, 0, 16, 16))
	cases := []struct {
		name   string
		src    image.Image
		key    []byte
		width  int
		height int
	}{
		{name: "nil-image", src: nil, key: []byte("Build98-Key"), width: 1632, height: 1632},
		{name: "short-key", src: validImage, key: []byte("short"), width: 1632, height: 1632},
		{name: "invalid-canonical", src: validImage, key: []byte("Build98-Key"), width: 288, height: 248},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantPayload, wantInfo, wantPhone, wantErr := ExperimentalV4ExtractPhoneBuild98Candidate(tc.src, tc.key, tc.width, tc.height)
			gotPayload, gotInfo, gotPhone, gotErr := ExperimentalV4ExtractPhone(tc.src, tc.key, tc.width, tc.height)
			if !bytes.Equal(gotPayload, wantPayload) || !reflect.DeepEqual(gotInfo, wantInfo) || !reflect.DeepEqual(gotPhone, wantPhone) {
				t.Fatalf("entry points diverged on %s", tc.name)
			}
			if (gotErr == nil) != (wantErr == nil) || (gotErr != nil && gotErr.Error() != wantErr.Error()) {
				t.Fatalf("error mismatch on %s: got=%v want=%v", tc.name, gotErr, wantErr)
			}
		})
	}
}

func TestExperimentalV4Build98ApplyTelemetryPreservesBuild93View(t *testing.T) {
	r := experimentalV4PhoneBuild98RecoveryTelemetry{}
	r.Attempted = true
	r.ProtectedSamplingPasses = 696
	r.ProtectedProfileDerivations = 2086
	r.Build98GeometryProfile.Gen3FetchFoldScores = 23024
	r.Build98GeometryProfile.Gen3FetchBlockReads = 29470720
	r.Build98GeometryProfile.Gen3FetchBlockSuccess = 29470720
	var p ExperimentalV4PhoneInfo
	experimentalV4PhoneBuild98ApplyTelemetry(&p, r)
	if !p.Build93Attempted || p.Build93ProtectedSamplingPasses != 696 || p.Build93ProtectedProfileDerivations != 2086 || !p.Build98Attempted || p.Build98Gen3FetchFoldScores != 23024 || p.Build98Gen3FetchBlockReads != 29470720 || p.Build98Gen3FetchBlockSuccess != 29470720 || p.Build98Gen3FetchBlockFailed != 0 {
		t.Fatalf("telemetry %+v", p)
	}
}
