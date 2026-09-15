package watermark

import "testing"

func TestExperimentalV4Build36SoftHammingRecoversWeakDoubleError(t *testing.T) {
	input := []byte{1, 0, 1, 1}
	encoded := hammingEncode(input)
	margins := make([]float64, len(encoded))
	for i, bit := range encoded {
		if bit == 0 {
			margins[i] = -9
		} else {
			margins[i] = 9
		}
	}
	// Two weak sign errors exceed Hamming(7,4)'s guaranteed hard-decision
	// correction radius, while the remaining five observations are strong.
	margins[0] *= -0.08
	margins[1] *= -0.08

	hard := make([]byte, len(margins))
	for i, margin := range margins {
		if margin >= 0 {
			hard[i] = 1
		}
	}
	if got := hammingDecode(hard); equalBitsV4Build36(got, input) {
		t.Fatal("crafted double-error word unexpectedly succeeded with hard Hamming")
	}
	if got := experimentalV4SoftHammingDecodeMargins(margins); !equalBitsV4Build36(got, input) {
		t.Fatalf("soft Hamming=%v want %v", got, input)
	}
}

func TestExperimentalV4Build36SoftHammingMatchesHardOnCleanWords(t *testing.T) {
	input := []byte{1, 0, 0, 1, 0, 1, 1, 0}
	encoded := hammingEncode(input)
	margins := make([]float64, len(encoded))
	for i, bit := range encoded {
		if bit == 0 {
			margins[i] = -4 - float64(i%3)
		} else {
			margins[i] = 4 + float64(i%3)
		}
	}
	soft := experimentalV4SoftHammingDecodeMargins(margins)
	hard := hammingDecode(encoded)
	if !equalBitsV4Build36(soft, hard) || !equalBitsV4Build36(soft, input) {
		t.Fatalf("soft=%v hard=%v input=%v", soft, hard, input)
	}
}

func equalBitsV4Build36(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
