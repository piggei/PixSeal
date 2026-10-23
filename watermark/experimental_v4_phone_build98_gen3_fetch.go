package watermark

import (
	"math"
	"sort"
)

// Build98 is a benchmark-first candidate selected after Build97 retained
// evidence failed its predeclared promotion gate. Build94 CPU profiles show
// that generation three still spends almost all of its sampled CPU in the
// historical generic projective FoldScore path: about 76.19 CPU-s on B/mild
// and 283.69 CPU-s on B/angle. Build98 therefore changes only generation-three
// scoring and reuses the already-qualified Build84 exact RGB reader/FoldScore.
// Prefix1, generation two, generation four, freeze, qualification and protected
// decode remain the qualified Build93 path. No image-sized memory is added.

type experimentalV4PhoneBuild98FetchTelemetry struct {
	Single       experimentalV4PhoneBuild84FetchTelemetry
	Pair         experimentalV4PhoneBuild84FetchTelemetry
	Continuation experimentalV4PhoneBuild84FetchTelemetry
	Sibling      experimentalV4PhoneBuild84FetchTelemetry
}

func experimentalV4PhoneBuild98SingleProbeFetch(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, root experimentalV4PhoneHypothesis, heldout int) (int, int, experimentalV4PhoneBuild84FetchTelemetry) {
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

func experimentalV4PhoneBuild98PairStencilFetch(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, root experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild53BlindPairState, int, int, experimentalV4PhoneBuild84FetchTelemetry) {
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

func experimentalV4PhoneBuild98SiblingStencilFetch(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, parent experimentalV4PhoneHypothesis, heldout int) ([]experimentalV4PhoneBuild55BlindSiblingState, int, experimentalV4PhoneBuild84FetchTelemetry) {
	states := make([]experimentalV4PhoneBuild55BlindSiblingState, 0, experimentalV4PhoneBuild55MaxSiblingEvals)
	evals := 0
	var fetch experimentalV4PhoneBuild84FetchTelemetry
	for dim := 0; dim < experimentalV4PhoneBuild55Dimensions; dim++ {
		for _, sign := range []float64{-1, 1} {
			qq := experimentalV4PhoneBuild51ApplyDim(parent.quad, dim, sign*experimentalV4PhoneBuild55Step)
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
			if math.IsInf(ss, 0) || math.IsNaN(ss) || ss <= parent.proposal+1e-7 {
				continue
			}
			states = append(states, experimentalV4PhoneBuild55BlindSiblingState{
				dim: dim, delta: sign * experimentalV4PhoneBuild55Step,
				hyp: experimentalV4PhoneHypothesis{quad: qq, h: hh, proposal: ss},
			})
		}
	}
	return states, evals, fetch
}

// experimentalV4PhoneBuild98Generation3Fetch reproduces the exact Build73
// generation-three subtree while routing every FoldScore call through Build84.
func experimentalV4PhoneBuild98Generation3Fetch(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, input experimentalV4PhoneHypothesis) ([]experimentalV4PhoneHypothesis, int, experimentalV4PhoneBuild98FetchTelemetry) {
	outputs := make([]experimentalV4PhoneHypothesis, 0, 8)
	evals := 0
	var fetch experimentalV4PhoneBuild98FetchTelemetry
	se3, si3, sf := experimentalV4PhoneBuild98SingleProbeFetch(plane, pilot, cw, ch, anchor, input, 0)
	experimentalV4PhoneBuild84MergeFetch(&fetch.Single, sf)
	evals += se3
	if si3 != 0 {
		return outputs, evals, fetch
	}
	pairs3, pe3, _, pf := experimentalV4PhoneBuild98PairStencilFetch(plane, pilot, cw, ch, anchor, input, 0)
	experimentalV4PhoneBuild84MergeFetch(&fetch.Pair, pf)
	evals += pe3
	for _, pair3 := range pairs3 {
		cont3, ce3, cf := experimentalV4PhoneBuild84ContinueFetch(plane, pilot, cw, ch, anchor, pair3.hyp, 0)
		experimentalV4PhoneBuild84MergeFetch(&fetch.Continuation, cf)
		evals += ce3
		for _, c3 := range cont3 {
			sibs3, ne3, nf := experimentalV4PhoneBuild98SiblingStencilFetch(plane, pilot, cw, ch, anchor, c3.hyp, 0)
			experimentalV4PhoneBuild84MergeFetch(&fetch.Sibling, nf)
			evals += ne3
			for _, s3 := range sibs3 {
				outputs = append(outputs, s3.hyp)
			}
		}
	}
	return outputs, evals, fetch
}

func experimentalV4PhoneBuild98TotalFetch(fetch experimentalV4PhoneBuild98FetchTelemetry) experimentalV4PhoneBuild84FetchTelemetry {
	var out experimentalV4PhoneBuild84FetchTelemetry
	experimentalV4PhoneBuild84MergeFetch(&out, fetch.Single)
	experimentalV4PhoneBuild84MergeFetch(&out, fetch.Pair)
	experimentalV4PhoneBuild84MergeFetch(&out, fetch.Continuation)
	experimentalV4PhoneBuild84MergeFetch(&out, fetch.Sibling)
	return out
}
