package watermark

import (
	"errors"
	"image"
	"math"
	"sort"
)

const (
	experimentalV4PhoneBuild42DataMaxBank   = 6
	experimentalV4PhoneBuild42ListWeakWords = 10
)

type experimentalV4PhoneBuild42DataTelemetry struct {
	Attempted          bool
	BankCandidates     int
	EnsemblesTried     int
	ListFramesTried    int
	Authenticated      bool
	Profile            Profile
	EnsembleCandidates int
}

type experimentalV4PhoneBuild42WordAlternative struct {
	word int
	alt  [4]byte
	gap  float64
}

func experimentalV4PhoneBuild42NibbleAlternatives(margins []float64) ([]byte, []experimentalV4PhoneBuild42WordAlternative) {
	base := experimentalV4SoftHammingDecodeMargins(margins)
	words := len(margins) / 7
	alts := make([]experimentalV4PhoneBuild42WordAlternative, 0, words)
	for word := 0; word < words; word++ {
		type candidate struct {
			value int
			score float64
		}
		scores := make([]candidate, 0, 16)
		for value := 0; value < 16; value++ {
			nibble := []byte{byte((value >> 3) & 1), byte((value >> 2) & 1), byte((value >> 1) & 1), byte(value & 1)}
			encoded := hammingEncode(nibble)
			score := 0.0
			for bit := 0; bit < 7; bit++ {
				sign := -1.0
				if encoded[bit] != 0 {
					sign = 1
				}
				score += sign * margins[word*7+bit]
			}
			scores = append(scores, candidate{value: value, score: score})
		}
		sort.Slice(scores, func(i, j int) bool {
			if scores[i].score == scores[j].score {
				return scores[i].value < scores[j].value
			}
			return scores[i].score > scores[j].score
		})
		v := scores[1].value
		alts = append(alts, experimentalV4PhoneBuild42WordAlternative{
			word: word,
			alt:  [4]byte{byte((v >> 3) & 1), byte((v >> 2) & 1), byte((v >> 1) & 1), byte(v & 1)},
			gap:  scores[0].score - scores[1].score,
		})
	}
	sort.Slice(alts, func(i, j int) bool {
		if alts[i].gap == alts[j].gap {
			return alts[i].word < alts[j].word
		}
		return alts[i].gap < alts[j].gap
	})
	return base, alts
}

// experimentalV4PhoneBuild42ListDecode authenticates only frames produced by
// a deterministic soft-decision list. The list is ordered solely by Hamming ML
// ambiguity; key/HMAC do not create or rank geometry or data hypotheses.
func experimentalV4PhoneBuild42ListDecode(margins []float64, key []byte, spec profileSpec) ([]byte, int, bool) {
	base, alts := experimentalV4PhoneBuild42NibbleAlternatives(margins)
	tried := 1
	raw := bitsToBytes(whiten(base, key, experimentalV4WhitenLabel))
	if payload, err := parseExperimentalV4Frame(raw, key, spec); err == nil {
		return payload, tried, true
	}

	limit := experimentalV4PhoneBuild42ListWeakWords
	if limit > len(alts) {
		limit = len(alts)
	}
	for mask := 1; mask < (1 << limit); mask++ {
		bits := append([]byte(nil), base...)
		for rank := 0; rank < limit; rank++ {
			if mask&(1<<rank) == 0 {
				continue
			}
			a := alts[rank]
			copy(bits[a.word*4:a.word*4+4], a.alt[:])
		}
		tried++
		raw = bitsToBytes(whiten(bits, key, experimentalV4WhitenLabel))
		if payload, err := parseExperimentalV4Frame(raw, key, spec); err == nil {
			return payload, tried, true
		}
	}
	return nil, tried, false
}

// experimentalV4PhoneDecodeBuild42Bank is a post-geometry data fallback. It
// enumerates deterministic three-geometry ensembles from a bank that was
// already fully qualified by Build41 public pilot evidence. HMAC is consulted
// only as final frame authentication.
func experimentalV4PhoneDecodeBuild42Bank(src image.Image, key []byte, cw, ch int, bank []experimentalV4PhoneHypothesis) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild42DataTelemetry, error) {
	candidate := experimentalV4Prototype2Candidate()
	info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate)}
	telemetry := experimentalV4PhoneBuild42DataTelemetry{Attempted: true, BankCandidates: len(bank)}
	if len(bank) < 3 {
		return nil, info, telemetry, errors.New("experimental v4 Build42 data bank too small")
	}

	// Bound runtime deterministically using only public qualification evidence.
	bank = append([]experimentalV4PhoneHypothesis(nil), bank...)
	sort.SliceStable(bank, func(i, j int) bool {
		if bank[i].validation == bank[j].validation {
			if bank[i].detection.Score == bank[j].detection.Score {
				return bank[i].proposal > bank[j].proposal
			}
			return bank[i].detection.Score > bank[j].detection.Score
		}
		return bank[i].validation > bank[j].validation
	})
	if len(bank) > experimentalV4PhoneBuild42DataMaxBank {
		bank = bank[:experimentalV4PhoneBuild42DataMaxBank]
	}
	telemetry.EnsembleCandidates = len(bank)
	plane := newPixelPlane(src)

	for _, spec := range v3Profiles {
		margins := make([][]float64, len(bank))
		valid := make([]bool, len(bank))
		for i, q := range bank {
			var m []float64
			var ok bool
			if q.warp != nil {
				m, _, ok = experimentalV4ReadProtectedPhoneWarpMargins(plane, cw, ch, q.h, q.warp, spec.codedBits)
			} else {
				m, _, ok = experimentalV4ReadProtectedProjectiveMargins(plane, cw, ch, q.h, spec.codedBits)
			}
			if ok {
				margins[i], valid[i] = m, true
			}
		}

		for i := 0; i < len(bank)-2; i++ {
			if !valid[i] {
				continue
			}
			for j := i + 1; j < len(bank)-1; j++ {
				if !valid[j] {
					continue
				}
				for k := j + 1; k < len(bank); k++ {
					if !valid[k] {
						continue
					}
					telemetry.EnsemblesTried++
					agg := make([]float64, spec.codedBits)
					for bit := range agg {
						agg[bit] = (margins[i][bit] + margins[j][bit] + margins[k][bit]) / 3.0
					}
					payload, tried, ok := experimentalV4PhoneBuild42ListDecode(agg, key, spec)
					telemetry.ListFramesTried += tried
					if ok {
						info.Profile = spec.profile
						info.PilotScore = math.Max(bank[i].detection.Score, math.Max(bank[j].detection.Score, bank[k].detection.Score))
						telemetry.Authenticated = true
						telemetry.Profile = spec.profile
						return payload, info, telemetry, nil
					}
				}
			}
		}
	}
	return nil, info, telemetry, errors.New("experimental v4 Build42 data list authentication failed")
}
