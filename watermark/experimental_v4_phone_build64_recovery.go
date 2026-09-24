package watermark

import (
	"errors"
	"image"
)

// experimentalV4PhoneBuild64RecoveryTelemetry records the additive recovery
// fallback introduced by Build64. Geometry generation remains proposal-only;
// qualification and HMAC are evaluated only after the complete recovery bank
// has been frozen.
type experimentalV4PhoneBuild64RecoveryTelemetry struct {
	Attempted                bool
	SeedsSelected            int
	GeometryEvaluations      int
	BankCandidates           int
	QualificationEvaluations int
	QualifiedCandidates      int
	DecodeCandidatesTried    int
	ProfilesTried            int
	ListFramesTried          int
	MaxDataConfidence        float64
	Authenticated            bool
	Profile                  Profile
}

// experimentalV4PhoneBuild64BlindBank reproduces only the blind geometry phase
// of Build63 and returns the final fourth-pair continuation bank. It never reads
// held-out validation, key material, payload/ECC or HMAC evidence.
func experimentalV4PhoneBuild64BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int) ([]experimentalV4PhoneHypothesis, int, int) {
	if !experimentalV4PhoneBuild41BoundarySaneForImage(work, boundary) {
		return nil, 0, 0
	}
	anchor := experimentalV4PhoneBuild41Quad(boundary)
	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	frozen, _, evals := experimentalV4PhoneBuild47Freeze(work, boundary, cw, ch, 128)
	seeds := experimentalV4PhoneBuild48SelectSeeds(frozen, experimentalV4PhoneBuild63SeedsPerPair)
	bank := make([]experimentalV4PhoneHypothesis, 0, 1024)

	for _, seed := range seeds {
		baseline, n := experimentalV4PhoneBuild41Refine(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		evals += n
		roots, n := experimentalV4PhoneBuild53TwoPxRoots(plane, pilot, cw, ch, anchor, seed.frozen.h.quad, 0)
		evals += n
		if baseline.h.h[8] == 0 || len(roots) == 0 {
			continue
		}
		for _, root := range roots {
			se, si := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, root.hyp, 0)
			evals += se
			if si != 0 {
				continue
			}
			pairs, pe, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, root.hyp, 0)
			evals += pe
			for _, pair := range pairs {
				cont, ce := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair.hyp, 0)
				evals += ce
				for _, cs := range cont {
					sibs, ne := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, cs.hyp, 0)
					evals += ne
					for _, sib := range sibs {
						se2, si2 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, sib.hyp, 0)
						evals += se2
						if si2 != 0 {
							continue
						}
						pairs2, pe2, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, sib.hyp, 0)
						evals += pe2
						for _, pair2 := range pairs2 {
							cont2, ce2 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair2.hyp, 0)
							evals += ce2
							for _, c2 := range cont2 {
								sibs2, ne2 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c2.hyp, 0)
								evals += ne2
								for _, s2 := range sibs2 {
									se3, si3 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s2.hyp, 0)
									evals += se3
									if si3 != 0 {
										continue
									}
									pairs3, pe3, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s2.hyp, 0)
									evals += pe3
									for _, pair3 := range pairs3 {
										cont3, ce3 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair3.hyp, 0)
										evals += ce3
										for _, c3 := range cont3 {
											sibs3, ne3 := experimentalV4PhoneBuild55SiblingStencil(plane, pilot, cw, ch, anchor, c3.hyp, 0)
											evals += ne3
											for _, s3 := range sibs3 {
												se4, si4 := experimentalV4PhoneBuild53SingleProbe(plane, pilot, cw, ch, anchor, s3.hyp, 0)
												evals += se4
												if si4 != 0 {
													continue
												}
												pairs4, pe4, _ := experimentalV4PhoneBuild53PairStencil(plane, pilot, cw, ch, anchor, s3.hyp, 0)
												evals += pe4
												for _, pair4 := range pairs4 {
													cont4, ce4 := experimentalV4PhoneBuild55Continue(plane, pilot, cw, ch, anchor, pair4.hyp, 0)
													evals += ce4
													for _, c4 := range cont4 {
														bank = append(bank, c4.hyp)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return bank, evals, len(seeds)
}

func experimentalV4PhoneBuild64DecodeSingle(plane *pixelPlane, h experimentalV4PhoneHypothesis, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, int, int, float64, error) {
	candidate := experimentalV4Prototype2Candidate()
	info := ExperimentalV4ExtractInfo{
		Version:       experimentalV4Version,
		PilotName:     candidate.name,
		PilotHash:     experimentalV4PilotCandidateHash(candidate),
		PilotScore:    h.detection.Score,
		PilotMargin:   h.detection.Margin,
		OriginXBlocks: h.detection.OriginXBlocks,
		OriginYBlocks: h.detection.OriginYBlocks,
	}
	profiles, frames := 0, 0
	maxConfidence := 0.0
	for _, spec := range v3Profiles {
		margins, confidence, ok := experimentalV4ReadProtectedProjectiveMargins(plane, cw, ch, h.h, spec.codedBits)
		if !ok {
			continue
		}
		profiles++
		if confidence > maxConfidence {
			maxConfidence = confidence
		}
		payload, tried, authenticated := experimentalV4PhoneBuild42ListDecode(margins, key, spec)
		frames += tried
		if authenticated {
			info.Profile = spec.profile
			info.Confidence = confidence
			return payload, info, profiles, frames, maxConfidence, nil
		}
	}
	return nil, info, profiles, frames, maxConfidence, errors.New("experimental v4 Build64 single-candidate authentication failed")
}

// experimentalV4PhoneBuild64Recover is an additive production-candidate
// fallback. It first freezes the complete Build63-derived blind geometry bank,
// then qualifies every geometry with public held-out pilot evidence, and only
// after that attempts data/HMAC authentication in deterministic bank order.
func experimentalV4PhoneBuild64Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneBuild64RecoveryTelemetry, error) {
	telemetry := experimentalV4PhoneBuild64RecoveryTelemetry{Attempted: true}
	bank, evals, seeds := experimentalV4PhoneBuild64BlindBank(work, boundary, cw, ch)
	telemetry.GeometryEvaluations = evals
	telemetry.SeedsSelected = seeds
	telemetry.BankCandidates = len(bank)
	if len(bank) == 0 {
		return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build64 recovery bank empty")
	}

	plane := newPixelPlane(work)
	pilot := experimentalV4Prototype2Candidate()
	qualified := make([]experimentalV4PhoneHypothesis, 0, len(bank))
	for _, h := range bank {
		q, n, ok := experimentalV4PhoneBuild41Qualify(plane, work, pilot, cw, ch, h, 0)
		telemetry.QualificationEvaluations += n
		if ok {
			qualified = append(qualified, q)
		}
	}
	telemetry.QualifiedCandidates = len(qualified)
	for _, q := range qualified {
		telemetry.DecodeCandidatesTried++
		payload, info, profiles, frames, confidence, err := experimentalV4PhoneBuild64DecodeSingle(plane, q, key, cw, ch)
		telemetry.ProfilesTried += profiles
		telemetry.ListFramesTried += frames
		if confidence > telemetry.MaxDataConfidence {
			telemetry.MaxDataConfidence = confidence
		}
		if err == nil {
			telemetry.Authenticated = true
			telemetry.Profile = info.Profile
			return payload, info, telemetry, nil
		}
	}
	return nil, ExperimentalV4ExtractInfo{}, telemetry, errors.New("experimental v4 Build64 recovery authentication failed")
}

func experimentalV4PhoneBuild64ApplyTelemetry(public *ExperimentalV4PhoneInfo, recovery experimentalV4PhoneBuild64RecoveryTelemetry) {
	public.Build64Attempted = recovery.Attempted
	public.Build64SeedsSelected = recovery.SeedsSelected
	public.Build64GeometryEvaluations = recovery.GeometryEvaluations
	public.Build64BankCandidates = recovery.BankCandidates
	public.Build64QualifiedCandidates = recovery.QualifiedCandidates
	public.Build64DecodeCandidates = recovery.DecodeCandidatesTried
	public.Build64ListFramesTried = recovery.ListFramesTried
	if recovery.DecodeCandidatesTried > 0 {
		public.DataDecodeAttempted = true
	}
	if recovery.ProfilesTried > public.SoftHammingProfiles {
		public.SoftHammingProfiles = recovery.ProfilesTried
	}
	public.Build64MaxDataConfidence = recovery.MaxDataConfidence
	public.Build64Authenticated = recovery.Authenticated
	public.HypothesesEvaluated += recovery.GeometryEvaluations + recovery.QualificationEvaluations
	if recovery.MaxDataConfidence > public.MaxDataConfidence {
		public.MaxDataConfidence = recovery.MaxDataConfidence
	}
}
