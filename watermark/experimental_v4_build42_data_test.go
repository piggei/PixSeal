package watermark

import "testing"

func TestExperimentalV4Build42ListDecodeRecoversSecondBestWord(t *testing.T) {
	key := []byte("PixSeal-v4-TestKey-2026")
	spec, ok := profileSpecFor(ProfileRobust)
	if !ok {
		t.Fatal("robust profile unavailable")
	}
	payload := []byte("build42-list")
	frame := makeExperimentalV4Frame(payload, key, spec)
	infoBits := whiten(bytesToBits(frame), key, experimentalV4WhitenLabel)
	protected := hammingEncode(infoBits)
	margins := make([]float64, len(protected))
	for i, bit := range protected {
		if bit != 0 {
			margins[i] = 5
		} else {
			margins[i] = -5
		}
	}

	// Make one Hamming word deliberately ambiguous, with a wrong nibble just
	// above the correct one. The list decoder must recover the authenticated
	// frame by trying the second ML nibble for this word.
	word := 7
	base := word * 4
	correctValue := int(infoBits[base])<<3 | int(infoBits[base+1])<<2 | int(infoBits[base+2])<<1 | int(infoBits[base+3])
	wrongValue := correctValue ^ 0x1
	correctNibble := []byte{byte((correctValue >> 3) & 1), byte((correctValue >> 2) & 1), byte((correctValue >> 1) & 1), byte(correctValue & 1)}
	wrongNibble := []byte{byte((wrongValue >> 3) & 1), byte((wrongValue >> 2) & 1), byte((wrongValue >> 1) & 1), byte(wrongValue & 1)}
	correctCode := hammingEncode(correctNibble)
	wrongCode := hammingEncode(wrongNibble)
	for bit := 0; bit < 7; bit++ {
		cs, ws := -1.0, -1.0
		if correctCode[bit] != 0 {
			cs = 1
		}
		if wrongCode[bit] != 0 {
			ws = 1
		}
		margins[word*7+bit] = cs + 1.05*ws
	}

	soft := experimentalV4SoftHammingDecodeMargins(margins)
	if soft[base] == infoBits[base] && soft[base+1] == infoBits[base+1] && soft[base+2] == infoBits[base+2] && soft[base+3] == infoBits[base+3] {
		t.Fatal("test setup did not make the base ML nibble wrong")
	}
	got, tried, authenticated := experimentalV4PhoneBuild42ListDecode(margins, key, spec)
	if !authenticated {
		t.Fatalf("list decode did not authenticate after %d candidates", tried)
	}
	if string(got) != string(payload) {
		t.Fatalf("payload=%q want %q", got, payload)
	}
	if tried <= 1 {
		t.Fatalf("expected list fallback, tried=%d", tried)
	}
}
