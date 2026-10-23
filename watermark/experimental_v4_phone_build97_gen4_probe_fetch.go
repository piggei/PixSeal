package watermark

import (
	"math"
	"sort"
)

// Build97 is a benchmark-first candidate selected from Build96 retained
// profiling. Build96 showed that, inside retained B/angle generation four,
// the already-qualified Build84 continuation reader still consumes about one
// third of sampled CPU, while a second large share remains in the historical
// generic samplePlaneLuminance path used by generation-four single/pair probes.
//
// Build97 does not change the reader itself. It reuses the already-qualified,
// bit-exact Build84 RGB reader/FoldScore implementation for generation-four
// single-probe and pair-stencil scoring as well as continuation4. Prefix1,
// generations 2/3, candidate ordering, thresholds, freeze, qualification and
// protected decode remain Build93. No extra image-sized memory is allocated.

// experimentalV4PhoneBuild97FetchTelemetry keeps the qualified historical
// continuation4 counters separate from the new single4/pair4 probe work.
// Build84GeometryTelemetry.Fetch* is a frozen qualified counter family and
// must continue to describe continuation4 only; Build97 probe work is exposed
// through a candidate-local counter family instead of changing that history.
type experimentalV4PhoneBuild97FetchTelemetry struct {
	Continuation experimentalV4PhoneBuild84FetchTelemetry
	Probe        experimentalV4PhoneBuild84FetchTelemetry
}

func experimentalV4PhoneBuild97SingleProbeFetch(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, root experimentalV4PhoneHypothesis, heldout int) (int, int, experimentalV4PhoneBuild84FetchTelemetry) {
	evals, improving := 0, 0
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	for dim := 0; dim < experimentalV4PhoneBuild53Dimensions; dim++ {
		for _, sign := range []float64{-1, 1} {
			qq := experimentalV4PhoneBuild51ApplyDim(root.quad, dim, sign*experimentalV4PhoneBuild53ProbeStep)
			if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
				continue
			}
			hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
			if !ok {
				continue
			}
			ss, _, f := experimentalV4PhoneBuild84FoldScoreFetch(plane, candidate, cw, ch, hh, heldout, true, false)
			experimentalV4PhoneBuild84MergeFetch(&fetch, f)
			evals++
			if !math.IsInf(ss, 0) && !math.IsNaN(ss) && ss > root.proposal+1e-7 {
				improving++
			}
		}
	}
	return evals, improving, fetch
}

func experimentalV4PhoneBuild97PairStencilFetch(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, root experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild53BlindPairState, int, int, experimentalV4PhoneBuild84FetchTelemetry) {
	evals := 0
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	all := make([]experimentalV4PhoneBuild53BlindPairState, 0, experimentalV4PhoneBuild53PairKeepPerRoot)
	for dimA := 0; dimA < experimentalV4PhoneBuild53Dimensions; dimA++ {
		for dimB := dimA + 1; dimB < experimentalV4PhoneBuild53Dimensions; dimB++ {
			for _, signA := range []float64{-1, 1} {
				for _, signB := range []float64{-1, 1} {
					deltaA := signA * experimentalV4PhoneBuild53ProbeStep
					deltaB := signB * experimentalV4PhoneBuild53ProbeStep
					qq := experimentalV4PhoneBuild51ApplyDim(root.quad, dimA, deltaA)
					qq = experimentalV4PhoneBuild51ApplyDim(qq, dimB, deltaB)
					if !experimentalV4PhoneBuild41WithinLimit(qq, anchor) {
						continue
					}
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					ss, _, f := experimentalV4PhoneBuild84FoldScoreFetch(plane, candidate, cw, ch, hh, heldout, true, false)
					experimentalV4PhoneBuild84MergeFetch(&fetch, f)
					evals++
					if !math.IsInf(ss, 0) && !math.IsNaN(ss) && ss > root.proposal+1e-7 {
						all = append(all, experimentalV4PhoneBuild53BlindPairState{
							dimA: dimA, dimB: dimB, deltaA: deltaA, deltaB: deltaB,
							hyp: experimentalV4PhoneHypothesis{quad: qq, h: hh, proposal: ss},
						})
					}
				}
			}
		}
	}
	improving := len(all)
	sort.SliceStable(all, func(i, j int) bool {
		if math.Abs(all[i].hyp.proposal-all[j].hyp.proposal) > 1e-12 {
			return all[i].hyp.proposal > all[j].hyp.proposal
		}
		if all[i].dimA != all[j].dimA {
			return all[i].dimA < all[j].dimA
		}
		if all[i].dimB != all[j].dimB {
			return all[i].dimB < all[j].dimB
		}
		if all[i].deltaA != all[j].deltaA {
			return all[i].deltaA < all[j].deltaA
		}
		return all[i].deltaB < all[j].deltaB
	})
	if len(all) > experimentalV4PhoneBuild53PairKeepPerRoot {
		all = all[:experimentalV4PhoneBuild53PairKeepPerRoot]
	}
	return all, evals, improving, fetch
}

func experimentalV4PhoneBuild97Generation4Fetch(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild97FetchTelemetry) {
	bank := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	var fetch experimentalV4PhoneBuild97FetchTelemetry
	se4, si4, sf := experimentalV4PhoneBuild97SingleProbeFetch(plane, pilot, cw, ch, anchor, input, 0)
	experimentalV4PhoneBuild84MergeFetch(&fetch.Probe, sf)
	evals += se4
	if si4 != 0 {
		return bank, evals, fetch
	}
	pairs4, pe4, _, pf := experimentalV4PhoneBuild97PairStencilFetch(plane, pilot, cw, ch, anchor, input, 0)
	experimentalV4PhoneBuild84MergeFetch(&fetch.Probe, pf)
	evals += pe4
	for _, pair4 := range pairs4 {
		cont4, ce4, cf := experimentalV4PhoneBuild84ContinueFetch(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
		experimentalV4PhoneBuild84MergeFetch(&fetch.Continuation, cf)
		evals += ce4
		for _, c4 := range cont4 {
			bank = append(bank, c4.hyp)
		}
	}
	return bank, evals, fetch
}
