package watermark

import "math"

// DiagnosticPhysicalTopologyObservability is the build22 image-domain probe of
// the weak absolute asymmetry identified by build21. For each unit shift it
// partitions repetition edges into (a) edges shared by both origin hypotheses
// and (b) edges unique to the canonical Format-v3 topology. Shared edges alone
// register a symmetric pair of phases separated by the tested shift. The held-
// out unique edges then choose between those two phases. No key, known header,
// payload, CRC or HMAC is used, and the result is telemetry only.
type DiagnosticPhysicalTopologyObservability struct {
	Method                        string                              `json:"method"`
	Available                     bool                                `json:"available"`
	Profiles                      []DiagnosticPhysicalTopologyProfile `json:"profiles,omitempty"`
	ProfilesAvailable             int                                 `json:"profiles_available"`
	AuditedShifts                 int                                 `json:"audited_shifts"`
	MeanWinnerModalFraction       float64                             `json:"mean_winner_modal_phase_fraction"`
	MinimumWinnerModalFraction    float64                             `json:"minimum_winner_modal_phase_fraction"`
	MeanDirectionalConsistency    float64                             `json:"mean_directional_cell_consistency"`
	MinimumDirectionalConsistency float64                             `json:"minimum_directional_cell_consistency"`
	MeanAbsoluteCellDelta         float64                             `json:"mean_absolute_cell_delta"`
	Note                          string                              `json:"note"`
}

type DiagnosticPhysicalTopologyProfile struct {
	Profile                       Profile                           `json:"profile"`
	Available                     bool                              `json:"available"`
	Shifts                        []DiagnosticPhysicalTopologyShift `json:"unit_shifts,omitempty"`
	ShiftsCompared                int                               `json:"shifts_compared"`
	WinnerPhases                  int                               `json:"winner_phases"`
	UniqueWinnerPhases            int                               `json:"unique_winner_phases"`
	ModalWinnerPhaseX             int                               `json:"modal_winner_phase_x"`
	ModalWinnerPhaseY             int                               `json:"modal_winner_phase_y"`
	ModalWinnerCount              int                               `json:"modal_winner_count"`
	ModalWinnerFraction           float64                           `json:"modal_winner_fraction"`
	MeanDirectionalConsistency    float64                           `json:"mean_directional_cell_consistency"`
	MinimumDirectionalConsistency float64                           `json:"minimum_directional_cell_consistency"`
	MeanAbsoluteCellDelta         float64                           `json:"mean_absolute_cell_delta"`
}

type DiagnosticPhysicalTopologyShift struct {
	DX                     int     `json:"dx"`
	DY                     int     `json:"dy"`
	CommonPairs            int     `json:"registration_common_pairs"`
	ExclusivePairs         int     `json:"heldout_exclusive_pairs"`
	PhaseAX                int     `json:"phase_a_x"`
	PhaseAY                int     `json:"phase_a_y"`
	PhaseBX                int     `json:"phase_b_x"`
	PhaseBY                int     `json:"phase_b_y"`
	RegistrationScoreA     float64 `json:"registration_score_a"`
	RegistrationScoreB     float64 `json:"registration_score_b"`
	RegistrationScore      float64 `json:"registration_min_score"`
	AggregateScoreA        float64 `json:"heldout_score_a"`
	AggregateScoreB        float64 `json:"heldout_score_b"`
	AggregateDeltaAminusB  float64 `json:"heldout_delta_a_minus_b"`
	AggregateAbsoluteDelta float64 `json:"heldout_absolute_delta"`
	Winner                 string  `json:"winner"`
	WinnerPhaseX           int     `json:"winner_phase_x"`
	WinnerPhaseY           int     `json:"winner_phase_y"`
	CellsCompared          int     `json:"cells_compared"`
	MeanCellDeltaAminusB   float64 `json:"mean_cell_delta_a_minus_b"`
	MeanAbsoluteCellDelta  float64 `json:"mean_absolute_cell_delta"`
	CellDeltaStdDev        float64 `json:"cell_delta_stddev"`
	CellEffectSize         float64 `json:"cell_effect_size"`
	CellZScore             float64 `json:"cell_z_score"`
	CellsFavorA            int     `json:"cells_favor_a"`
	CellsFavorB            int     `json:"cells_favor_b"`
	DirectionalConsistency float64 `json:"directional_cell_consistency"`
}

const diagnosticPhysicalTopologyProfiles = 2
const diagnosticPhysicalTopologyUnitShifts = 8

func diagnosticPhysicalTopologyObservability(cells []diagnosticSpatialGridCell, aggregate []float64) DiagnosticPhysicalTopologyObservability {
	result := DiagnosticPhysicalTopologyObservability{
		Method:                        "v3-image-domain-heldout-repetition-topology",
		MinimumWinnerModalFraction:    1,
		MinimumDirectionalConsistency: 1,
		Note:                          "diagnostic only: shared repetition edges register two shift-separated phase hypotheses; held-out topology-exclusive edges choose between them; no key/header/payload/CRC/HMAC",
	}
	if len(cells) < 3 || len(aggregate) < eccBits {
		result.MinimumWinnerModalFraction = 0
		result.MinimumDirectionalConsistency = 0
		return result
	}
	for _, profile := range []Profile{ProfileRobust, ProfileBalanced} {
		profileResult := diagnosticPhysicalTopologyProfileObservability(profile, cells, aggregate)
		result.Profiles = append(result.Profiles, profileResult)
		if !profileResult.Available {
			continue
		}
		result.ProfilesAvailable++
		result.AuditedShifts += profileResult.ShiftsCompared
		result.MeanWinnerModalFraction += profileResult.ModalWinnerFraction
		result.MeanDirectionalConsistency += profileResult.MeanDirectionalConsistency
		result.MeanAbsoluteCellDelta += profileResult.MeanAbsoluteCellDelta
		if profileResult.ModalWinnerFraction < result.MinimumWinnerModalFraction {
			result.MinimumWinnerModalFraction = profileResult.ModalWinnerFraction
		}
		if profileResult.MinimumDirectionalConsistency < result.MinimumDirectionalConsistency {
			result.MinimumDirectionalConsistency = profileResult.MinimumDirectionalConsistency
		}
	}
	result.Available = result.ProfilesAvailable > 0 && result.AuditedShifts > 0
	if result.ProfilesAvailable > 0 {
		den := float64(result.ProfilesAvailable)
		result.MeanWinnerModalFraction /= den
		result.MeanDirectionalConsistency /= den
		result.MeanAbsoluteCellDelta /= den
	} else {
		result.MinimumWinnerModalFraction = 0
		result.MinimumDirectionalConsistency = 0
	}
	return result
}

func diagnosticPhysicalTopologyProfileObservability(profile Profile, cells []diagnosticSpatialGridCell, aggregate []float64) DiagnosticPhysicalTopologyProfile {
	out := DiagnosticPhysicalTopologyProfile{
		Profile:                       profile,
		MinimumDirectionalConsistency: 1,
	}
	winnerCounts := map[[2]int]int{}
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			shift, ok := diagnosticPhysicalTopologyShiftObservability(profile, cells, aggregate, dx, dy)
			if !ok {
				continue
			}
			out.Shifts = append(out.Shifts, shift)
			out.ShiftsCompared++
			out.MeanDirectionalConsistency += shift.DirectionalConsistency
			out.MeanAbsoluteCellDelta += shift.MeanAbsoluteCellDelta
			if shift.DirectionalConsistency < out.MinimumDirectionalConsistency {
				out.MinimumDirectionalConsistency = shift.DirectionalConsistency
			}
			winnerCounts[[2]int{shift.WinnerPhaseX, shift.WinnerPhaseY}]++
		}
	}
	out.Available = out.ShiftsCompared > 0
	out.WinnerPhases = out.ShiftsCompared
	out.UniqueWinnerPhases = len(winnerCounts)
	for phase, count := range winnerCounts {
		if count > out.ModalWinnerCount || (count == out.ModalWinnerCount && (phase[1] < out.ModalWinnerPhaseY || (phase[1] == out.ModalWinnerPhaseY && phase[0] < out.ModalWinnerPhaseX))) {
			out.ModalWinnerPhaseX = phase[0]
			out.ModalWinnerPhaseY = phase[1]
			out.ModalWinnerCount = count
		}
	}
	if out.ShiftsCompared > 0 {
		den := float64(out.ShiftsCompared)
		out.ModalWinnerFraction = float64(out.ModalWinnerCount) / den
		out.MeanDirectionalConsistency /= den
		out.MeanAbsoluteCellDelta /= den
	} else {
		out.MinimumDirectionalConsistency = 0
	}
	return out
}

func diagnosticPhysicalTopologyShiftObservability(profile Profile, cells []diagnosticSpatialGridCell, aggregate []float64, dx, dy int) (DiagnosticPhysicalTopologyShift, bool) {
	out := DiagnosticPhysicalTopologyShift{DX: dx, DY: dy}
	common, exclusive := diagnosticPhysicalTopologyPairPartition(profile, dx, dy)
	out.CommonPairs = len(common)
	out.ExclusivePairs = len(exclusive)
	if len(common) == 0 || len(exclusive) == 0 {
		return out, false
	}
	phaseAX, phaseAY, scoreA, scoreB, ok := diagnosticPhysicalTopologyRegisterPair(aggregate, common, dx, dy, -1, -1)
	if !ok {
		return out, false
	}
	phaseBX := positiveMod(phaseAX+dx, tileWidth)
	phaseBY := positiveMod(phaseAY+dy, tileHeight)
	out.PhaseAX, out.PhaseAY = phaseAX, phaseAY
	out.PhaseBX, out.PhaseBY = phaseBX, phaseBY
	out.RegistrationScoreA, out.RegistrationScoreB = scoreA, scoreB
	out.RegistrationScore = math.Min(scoreA, scoreB)
	out.AggregateScoreA = diagnosticBlindRepetitionScorePairs(aggregate, phaseAX, phaseAY, exclusive)
	out.AggregateScoreB = diagnosticBlindRepetitionScorePairs(aggregate, phaseBX, phaseBY, exclusive)
	out.AggregateDeltaAminusB = out.AggregateScoreA - out.AggregateScoreB
	out.AggregateAbsoluteDelta = math.Abs(out.AggregateDeltaAminusB)
	if out.AggregateDeltaAminusB >= 0 {
		out.Winner = "a"
		out.WinnerPhaseX, out.WinnerPhaseY = phaseAX, phaseAY
	} else {
		out.Winner = "b"
		out.WinnerPhaseX, out.WinnerPhaseY = phaseBX, phaseBY
	}

	deltas := make([]float64, 0, len(cells))
	for _, cell := range cells {
		localAX, localAY, _, _, localOK := diagnosticPhysicalTopologyRegisterPair(cell.grid, common, dx, dy, phaseAX, phaseAY)
		if !localOK {
			continue
		}
		localBX := positiveMod(localAX+dx, tileWidth)
		localBY := positiveMod(localAY+dy, tileHeight)
		scoreLocalA := diagnosticBlindRepetitionScorePairs(cell.grid, localAX, localAY, exclusive)
		scoreLocalB := diagnosticBlindRepetitionScorePairs(cell.grid, localBX, localBY, exclusive)
		delta := scoreLocalA - scoreLocalB
		deltas = append(deltas, delta)
		out.MeanCellDeltaAminusB += delta
		out.MeanAbsoluteCellDelta += math.Abs(delta)
		if delta >= 0 {
			out.CellsFavorA++
		} else {
			out.CellsFavorB++
		}
	}
	out.CellsCompared = len(deltas)
	if out.CellsCompared < 3 {
		return out, false
	}
	out.MeanCellDeltaAminusB /= float64(out.CellsCompared)
	out.MeanAbsoluteCellDelta /= float64(out.CellsCompared)
	variance := 0.0
	for _, delta := range deltas {
		d := delta - out.MeanCellDeltaAminusB
		variance += d * d
	}
	if out.CellsCompared > 1 {
		variance /= float64(out.CellsCompared - 1)
	}
	out.CellDeltaStdDev = math.Sqrt(math.Max(0, variance))
	if out.CellDeltaStdDev > 1e-12 {
		out.CellEffectSize = out.MeanCellDeltaAminusB / out.CellDeltaStdDev
		out.CellZScore = out.MeanCellDeltaAminusB / (out.CellDeltaStdDev / math.Sqrt(float64(out.CellsCompared)))
	}
	majority := out.CellsFavorA
	if out.CellsFavorB > majority {
		majority = out.CellsFavorB
	}
	out.DirectionalConsistency = float64(majority) / float64(out.CellsCompared)
	return out, true
}

// diagnosticPhysicalTopologyRegisterPair selects a pair of phases separated by
// (dx,dy) using only edges common to both hypotheses. When referenceX/Y are
// non-negative the search is limited to the ordinary local phase radius around
// that reference; otherwise the entire tile is searched. The objective is the
// minimum of the two common-edge scores, so neither hypothesis can win during
// registration.
func diagnosticPhysicalTopologyRegisterPair(grid []float64, common [][2]int, dx, dy, referenceX, referenceY int) (bestX, bestY int, bestA, bestB float64, ok bool) {
	if len(grid) < eccBits || len(common) == 0 {
		return 0, 0, 0, 0, false
	}
	bestObjective := math.Inf(-1)
	evaluate := func(x, y int) {
		x = positiveMod(x, tileWidth)
		y = positiveMod(y, tileHeight)
		bx := positiveMod(x+dx, tileWidth)
		by := positiveMod(y+dy, tileHeight)
		a := diagnosticBlindRepetitionScorePairs(grid, x, y, common)
		b := diagnosticBlindRepetitionScorePairs(grid, bx, by, common)
		objective := math.Min(a, b)
		if objective > bestObjective || (objective == bestObjective && (y < bestY || (y == bestY && x < bestX))) {
			bestObjective = objective
			bestX, bestY, bestA, bestB = x, y, a, b
			ok = true
		}
	}
	if referenceX >= 0 && referenceY >= 0 {
		for oy := -diagnosticSpatialPhaseRadius; oy <= diagnosticSpatialPhaseRadius; oy++ {
			for ox := -diagnosticSpatialPhaseRadius; ox <= diagnosticSpatialPhaseRadius; ox++ {
				evaluate(referenceX+ox, referenceY+oy)
			}
		}
		return bestX, bestY, bestA, bestB, ok
	}
	for y := 0; y < tileHeight; y++ {
		for x := 0; x < tileWidth; x++ {
			evaluate(x, y)
		}
	}
	return bestX, bestY, bestA, bestB, ok
}

// diagnosticPhysicalTopologyPairPartition returns the large common edge set and
// the canonical edges that disappear under the tested unit shift. The latter
// are held out from registration and are the only edges used to choose between
// the two registered phase hypotheses.
func diagnosticPhysicalTopologyPairPartition(profile Profile, dx, dy int) (common, exclusive [][2]int) {
	base := diagnosticBlindRepetitionPairs(profile)
	if len(base) == 0 || (dx == 0 && dy == 0) {
		return nil, nil
	}
	shifted := make(map[uint32]struct{}, len(base))
	for _, pair := range base {
		a := diagnosticObservabilityShiftPosition(pair[0], dx, dy)
		b := diagnosticObservabilityShiftPosition(pair[1], dx, dy)
		shifted[diagnosticObservabilityPairKey(a, b)] = struct{}{}
	}
	for _, pair := range base {
		key := diagnosticObservabilityPairKey(pair[0], pair[1])
		if _, ok := shifted[key]; ok {
			common = append(common, pair)
		} else {
			exclusive = append(exclusive, pair)
		}
	}
	return common, exclusive
}
