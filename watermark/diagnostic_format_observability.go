package watermark

// DiagnosticFormatObservability is a static, key-independent audit of the
// public Format-v3 structure. It does not inspect payload/header/HMAC values and
// cannot authenticate or modify decoding. Build21 uses it to ask whether a
// +/-1 integer block-cycle shift is theoretically distinguishable from the
// format topology itself.
type DiagnosticFormatObservability struct {
	Method                                string                                 `json:"method"`
	Profiles                              []DiagnosticFormatObservabilityProfile `json:"profiles"`
	ProfilesWithRepetitionAnchor          int                                    `json:"profiles_with_repetition_anchor"`
	ProfilesWithSyntheticHammingAlias     int                                    `json:"profiles_with_synthetic_hamming_alias"`
	WeakestPairTopologyGap                float64                                `json:"weakest_pair_topology_gap"`
	WeakestSyntheticHammingSyndrome       float64                                `json:"weakest_synthetic_hamming_syndrome_fraction"`
	AuditedUnitShiftsUniversallySeparated bool                                   `json:"audited_unit_shifts_universally_separated"`
	Note                                  string                                 `json:"note"`
}

type DiagnosticFormatObservabilityProfile struct {
	Profile                               Profile                                  `json:"profile"`
	CodedBits                             int                                      `json:"coded_bits"`
	RepetitionPairs                       int                                      `json:"repetition_pairs"`
	RepetitionAbsoluteAnchor              bool                                     `json:"repetition_absolute_anchor"`
	UnitShifts                            []DiagnosticFormatObservabilityUnitShift `json:"unit_shifts"`
	MinimumPairTopologyGap                float64                                  `json:"minimum_pair_topology_gap"`
	MinimumSyntheticSyndrome              float64                                  `json:"minimum_synthetic_hamming_syndrome_fraction"`
	SyntheticHammingUnitShiftAliasPresent bool                                     `json:"synthetic_hamming_unit_shift_alias_present"`
}

type DiagnosticFormatObservabilityUnitShift struct {
	DX                              int     `json:"dx"`
	DY                              int     `json:"dy"`
	RepetitionPairOverlapFraction   float64 `json:"repetition_pair_overlap_fraction"`
	RepetitionPairTopologyGap       float64 `json:"repetition_pair_topology_gap"`
	MixedTargetCodeIndices          int     `json:"mixed_target_code_indices"`
	PureCodeIndexPermutation        bool    `json:"pure_code_index_permutation"`
	HammingWordPermutationInvariant bool    `json:"hamming_word_permutation_invariant"`
	SyntheticHammingSyndrome        float64 `json:"synthetic_hamming_syndrome_fraction"`
}

const diagnosticFormatObservabilitySyntheticFrames = 32

func diagnosticFormatObservabilityAudit() DiagnosticFormatObservability {
	result := DiagnosticFormatObservability{
		Method:                                "v3-key-independent-structural-observability",
		Note:                                  "diagnostic only: public tile/repetition/Hamming structure; no key, known header, payload, CRC or HMAC is used",
		WeakestPairTopologyGap:                1,
		WeakestSyntheticHammingSyndrome:       1,
		AuditedUnitShiftsUniversallySeparated: true,
	}
	for _, spec := range v3Profiles {
		profile := DiagnosticFormatObservabilityProfile{
			Profile:                  spec.profile,
			CodedBits:                spec.codedBits,
			RepetitionPairs:          diagnosticBlindRepetitionPairCount(spec.profile),
			MinimumPairTopologyGap:   1,
			MinimumSyntheticSyndrome: 1,
		}
		pairs := diagnosticBlindRepetitionPairs(spec.profile)
		baseSet := diagnosticObservabilityPairSet(pairs)
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				shift := DiagnosticFormatObservabilityUnitShift{DX: dx, DY: dy}
				if len(pairs) > 0 {
					shifted := make(map[uint32]struct{}, len(pairs))
					for _, pair := range pairs {
						a := diagnosticObservabilityShiftPosition(pair[0], dx, dy)
						b := diagnosticObservabilityShiftPosition(pair[1], dx, dy)
						shifted[diagnosticObservabilityPairKey(a, b)] = struct{}{}
					}
					common := 0
					for key := range baseSet {
						if _, ok := shifted[key]; ok {
							common++
						}
					}
					shift.RepetitionPairOverlapFraction = float64(common) / float64(len(baseSet))
					shift.RepetitionPairTopologyGap = 1 - shift.RepetitionPairOverlapFraction
					if shift.RepetitionPairTopologyGap < profile.MinimumPairTopologyGap {
						profile.MinimumPairTopologyGap = shift.RepetitionPairTopologyGap
					}
				} else {
					profile.MinimumPairTopologyGap = 0
				}
				shift.MixedTargetCodeIndices, shift.PureCodeIndexPermutation, shift.HammingWordPermutationInvariant = diagnosticObservabilityCodeIndexMapping(spec.codedBits, dx, dy)
				shift.SyntheticHammingSyndrome = diagnosticObservabilitySyntheticHammingSyndrome(spec.codedBits, dx, dy)
				if shift.SyntheticHammingSyndrome < profile.MinimumSyntheticSyndrome {
					profile.MinimumSyntheticSyndrome = shift.SyntheticHammingSyndrome
				}
				if shift.HammingWordPermutationInvariant || shift.SyntheticHammingSyndrome == 0 {
					profile.SyntheticHammingUnitShiftAliasPresent = true
				}
				profile.UnitShifts = append(profile.UnitShifts, shift)
			}
		}
		profile.RepetitionAbsoluteAnchor = profile.RepetitionPairs > 0 && profile.MinimumPairTopologyGap > 0
		if profile.RepetitionAbsoluteAnchor {
			result.ProfilesWithRepetitionAnchor++
		}
		if profile.SyntheticHammingUnitShiftAliasPresent {
			result.ProfilesWithSyntheticHammingAlias++
		}
		if profile.RepetitionPairs > 0 && profile.MinimumPairTopologyGap < result.WeakestPairTopologyGap {
			result.WeakestPairTopologyGap = profile.MinimumPairTopologyGap
		}
		if profile.MinimumSyntheticSyndrome < result.WeakestSyntheticHammingSyndrome {
			result.WeakestSyntheticHammingSyndrome = profile.MinimumSyntheticSyndrome
		}
		for _, shift := range profile.UnitShifts {
			if shift.RepetitionPairTopologyGap <= 0 && shift.SyntheticHammingSyndrome <= 0 {
				result.AuditedUnitShiftsUniversallySeparated = false
			}
		}
		result.Profiles = append(result.Profiles, profile)
	}
	return result
}

func diagnosticObservabilityShiftPosition(position, dx, dy int) int {
	x, y := position%tileWidth, position/tileWidth
	x = positiveMod(x+dx, tileWidth)
	y = positiveMod(y+dy, tileHeight)
	return y*tileWidth + x
}

func diagnosticObservabilityPairKey(a, b int) uint32 {
	if a > b {
		a, b = b, a
	}
	return uint32(a)<<16 | uint32(b)
}

func diagnosticObservabilityPairSet(pairs [][2]int) map[uint32]struct{} {
	out := make(map[uint32]struct{}, len(pairs))
	for _, pair := range pairs {
		out[diagnosticObservabilityPairKey(pair[0], pair[1])] = struct{}{}
	}
	return out
}

// diagnosticObservabilityCodeIndexMapping asks whether a wrong tile origin is
// merely a permutation of coded-bit indices. If so, it also checks whether that
// permutation maps every Hamming(7,4) word to another whole word while
// preserving each bit role 0..6, in which case Hamming syndrome is an exact
// structural alias for all valid codewords.
func diagnosticObservabilityCodeIndexMapping(codedBits, dx, dy int) (mixed int, pure, hammingInvariant bool) {
	sets := make([]map[int]struct{}, codedBits)
	for logical := 0; logical < eccBits; logical++ {
		x, y := logical%tileWidth, logical/tileWidth
		observedX := positiveMod(x-dx, tileWidth)
		observedY := positiveMod(y-dy, tileHeight)
		actualLogical := observedY*tileWidth + observedX
		target := v3CodeIndex(logical, codedBits)
		actual := v3CodeIndex(actualLogical, codedBits)
		if sets[target] == nil {
			sets[target] = make(map[int]struct{}, 2)
		}
		sets[target][actual] = struct{}{}
	}
	perm := make([]int, codedBits)
	pure = true
	for i, set := range sets {
		if len(set) != 1 {
			pure = false
			if len(set) > 1 {
				mixed++
			}
			continue
		}
		for value := range set {
			perm[i] = value
		}
	}
	if !pure || codedBits%7 != 0 {
		return mixed, pure, false
	}
	for word := 0; word < codedBits/7; word++ {
		sourceWord := -1
		for role := 0; role < 7; role++ {
			target := word*7 + role
			source := perm[target]
			if source%7 != role {
				return mixed, pure, false
			}
			if sourceWord < 0 {
				sourceWord = source / 7
			} else if source/7 != sourceWord {
				return mixed, pure, false
			}
		}
	}
	return mixed, pure, true
}

// The synthetic probe uses valid Hamming codewords generated by a fixed local
// PRNG and then applies exactly the same v3 tile mapping and hard aggregation as
// readV3ProtectedFrame. It estimates practical structural separation without
// any secret or known payload bits. A zero result means the tested shift is an
// alias for all sampled valid frames; it is diagnostic evidence, not a proof,
// unless HammingWordPermutationInvariant is also true.
func diagnosticObservabilitySyntheticHammingSyndrome(codedBits, dx, dy int) float64 {
	if codedBits <= 0 || codedBits%7 != 0 {
		return 0
	}
	state := uint64(0x9e3779b97f4a7c15) ^ uint64(codedBits*131+dx*17+dy*29)
	total := 0.0
	for frame := 0; frame < diagnosticFormatObservabilitySyntheticFrames; frame++ {
		input := make([]byte, codedBits/7*4)
		for i := range input {
			state ^= state << 13
			state ^= state >> 7
			state ^= state << 17
			input[i] = byte(state & 1)
		}
		coded := hammingEncode(input)
		grid := make([]float64, eccBits)
		for position := 0; position < eccBits; position++ {
			if coded[v3CodeIndex(position, codedBits)] != 0 {
				grid[position] = 1
			} else {
				grid[position] = -1
			}
		}
		shifted := readV3ProtectedFrameForObservability(grid, dx, dy, codedBits)
		bad := 0
		for start := 0; start+7 <= len(shifted); start += 7 {
			if diagnosticHammingSyndrome(shifted[start:start+7]) != 0 {
				bad++
			}
		}
		total += float64(bad) / float64(len(shifted)/7)
	}
	return total / diagnosticFormatObservabilitySyntheticFrames
}

func readV3ProtectedFrameForObservability(grid []float64, phaseX, phaseY, codedBits int) []byte {
	sums := make([]float64, codedBits)
	for logicalPosition := 0; logicalPosition < eccBits; logicalPosition++ {
		logicalX := logicalPosition % tileWidth
		logicalY := logicalPosition / tileWidth
		observedX := positiveMod(logicalX-phaseX, tileWidth)
		observedY := positiveMod(logicalY-phaseY, tileHeight)
		codeIndex := v3CodeIndex(logicalPosition, codedBits)
		sums[codeIndex] += grid[observedY*tileWidth+observedX]
	}
	coded := make([]byte, codedBits)
	for i, value := range sums {
		if value >= 0 {
			coded[i] = 1
		}
	}
	return coded
}
