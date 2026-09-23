package watermark

import "testing"

func TestExperimentalV4Build43PhoneSidePairRecovery(t *testing.T) {
	if len(experimentalV4PhoneBuild43Pairs) != 6 {
		t.Fatalf("side-pair count=%d want 6", len(experimentalV4PhoneBuild43Pairs))
	}
	if experimentalV4PhoneBuild43MaxFrozen > 32 {
		t.Fatalf("freeze budget=%d", experimentalV4PhoneBuild43MaxFrozen)
	}
	const width, height = 1184, 1024
	key := []byte("PixSeal-v4-TestKey-2026")
	opts := DefaultOptions()
	opts.Profile = ProfileRobust
	opts.Strength = 48
	source := testImage(width, height)
	payload := []byte("v4-b43-synth")
	marked, _, err := ExperimentalV4EmbedWithInfo(source, payload, key, opts)
	if err != nil {
		t.Fatal(err)
	}
	observed := experimentalV4Build41PageForTest(marked)
	boundary := experimentalV4PhoneBoundary(observed)
	bank, tele := experimentalV4PhoneSearchBuild43(observed, boundary, width, height)
	if tele.PairsScanned != 6 {
		t.Fatalf("pairs scanned=%d", tele.PairsScanned)
	}
	if tele.FrozenCandidates > 32 {
		t.Fatalf("frozen=%d", tele.FrozenCandidates)
	}
	if len(bank) < 3 {
		t.Fatalf("qualified=%d tele=%+v", len(bank), tele)
	}
	got, _, data, err := experimentalV4PhoneDecodeBuild42Bank(observed, key, width, height, bank)
	if err != nil {
		t.Fatalf("Build43 synthetic decode: %v tele=%+v", err, data)
	}
	if string(got) != string(payload) {
		t.Fatalf("payload=%q want=%q", got, payload)
	}

}
