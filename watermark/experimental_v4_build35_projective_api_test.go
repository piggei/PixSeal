package watermark

import (
	"bytes"
	"testing"
)

func TestExperimentalV4Build35ProjectiveAPI(t *testing.T) {
	const payloadText = "v4-b35-api"
	key := []byte("PixSeal-v4-TestKey-2026")
	source := testImage(888, 768)
	normalized, w, h := experimentalV4Build29NormalizeSource(source)
	carrier, _, err := ExperimentalV4EmbedWithInfo(normalized, []byte(payloadText), key, Options{Profile: ProfileRobust, Strength: 24})
	if err != nil {
		t.Fatal(err)
	}
	truth := experimentalV4BlindGeometryParams{angleDeg: 9.3, scaleX: 1.08, scaleY: .92, topInset: .03, bottomInset: .015}
	geometry, fw, fh := experimentalV4BlindHomography(w, h, truth)
	observedH, ow, oh := experimentalV4PlacementObservedHomographyForTest(geometry, fw, fh, [4]int{101, 69, 83, 47}, [4]int{})
	observed := experimentalV4WarpHomographyForTest(carrier, observedH, ow, oh)

	payload, info, projective, err := ExperimentalV4ExtractProjective(observed, key, w, h)
	if err != nil {
		t.Fatalf("Build35 public projective API failed: %v; info=%+v projective=%+v", err, info, projective)
	}
	if !bytes.Equal(payload, []byte(payloadText)) {
		t.Fatalf("payload=%q want %q", payload, payloadText)
	}
	if !projective.Accepted {
		t.Fatalf("projective telemetry reports rejected geometry: %+v", projective)
	}
	if info.OriginXBlocks != 0 || info.OriginYBlocks != 0 {
		t.Fatalf("canonical origin=(%d,%d), want (0,0)", info.OriginXBlocks, info.OriginYBlocks)
	}
	if projective.PlacementValidation < experimentalV4Build29MinProjectiveValidation {
		t.Fatalf("placement validation %.6f below Build29 floor %.6f", projective.PlacementValidation, experimentalV4Build29MinProjectiveValidation)
	}
	if info.PilotMargin < experimentalV4Build29MinProjectiveMargin {
		t.Fatalf("pilot margin %.6f below Build29 floor %.6f", info.PilotMargin, experimentalV4Build29MinProjectiveMargin)
	}
	t.Logf("payload=%q validation=%.3f score=%.3f margin=%.3f origin=(%d,%d) hypotheses=%d",
		payload, projective.PlacementValidation, info.PilotScore, info.PilotMargin, info.OriginXBlocks, info.OriginYBlocks, projective.HypothesesEvaluated)
}

func TestExperimentalV4Build35ProjectiveAPIRejectsInvalidCanonicalExtent(t *testing.T) {
	_, _, _, err := ExperimentalV4ExtractProjective(testImage(320, 320), []byte("12345678"), 319, 320)
	if err == nil {
		t.Fatal("projective API accepted a non-block-aligned canonical width")
	}
}
