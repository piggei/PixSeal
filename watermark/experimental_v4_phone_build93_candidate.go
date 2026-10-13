package watermark

import (
	"errors"
	"image"
)

// ExperimentalV4ExtractPhoneBuild93Candidate is the retained explicit Build93
// qualification/regression entry point. It pins the Build93 deep fallback even
// when the generic smartphone command evolves in later research snapshots. At the
// Build93 qualified snapshot its deep behavior matches ExperimentalV4ExtractPhone.
func ExperimentalV4ExtractPhoneBuild93Candidate(src image.Image, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, ExperimentalV4PhoneInfo, error) {
	if src == nil {
		return nil, ExperimentalV4ExtractInfo{}, ExperimentalV4PhoneInfo{}, errors.New("nil image")
	}
	if len(key) < 8 {
		return nil, ExperimentalV4ExtractInfo{}, ExperimentalV4PhoneInfo{}, errors.New("key must contain at least 8 bytes")
	}
	if cw < experimentalV4TileWidthBlocks*blockSize || ch < experimentalV4TileHeightBlocks*blockSize || cw%blockSize != 0 || ch%blockSize != 0 {
		return nil, ExperimentalV4ExtractInfo{}, ExperimentalV4PhoneInfo{}, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}

	work, boundary, hyp, build42Bank, evals, down := experimentalV4PhoneSearchBuild41Detailed(src, cw, ch)
	public := ExperimentalV4PhoneInfo{
		WorkingWidth:         work.Bounds().Dx(),
		WorkingHeight:        work.Bounds().Dy(),
		Downsampled:          down,
		BoundaryDetected:     boundary.Detected,
		BoundaryConfidence:   boundary.Confidence,
		ProjectiveBasinFound: len(hyp) > 0,
		HypothesesEvaluated:  evals,
		EnsembleCandidates:   len(hyp),
	}
	public.Build41DirectAccepted = len(hyp) >= experimentalV4PhoneEnsembleSize
	public.Build41QualifiedCandidates = len(build42Bank)
	if len(hyp) > 0 {
		q := hyp[0]
		public.ProposalScore = q.proposal
		public.ValidationScore = q.validation
		public.PilotScore = q.detection.Score
		public.PilotMargin = q.detection.Margin
		public.OriginXBlocks = q.detection.OriginXBlocks
		public.OriginYBlocks = q.detection.OriginYBlocks
	}
	public.Accepted = len(hyp) >= experimentalV4PhoneEnsembleSize
	if !public.Accepted {
		bank43, tele43 := experimentalV4PhoneSearchBuild43(work, boundary, cw, ch)
		public.Build43Attempted = tele43.Attempted
		public.Build43EdgeRefined = tele43.EdgeRefined
		public.Build43PairsScanned = tele43.PairsScanned
		public.Build43PairsSelected = tele43.PairsSelected
		public.Build43GeometryEvaluations = tele43.GeometryEvaluations
		public.Build43ProposalCandidates = tele43.ProposalCandidates
		public.Build43FrozenCandidates = tele43.FrozenCandidates
		public.Build43QualifiedCandidates = tele43.QualifiedCandidates
		public.Build43Pair0 = tele43.Pair0
		public.Build43Pair1 = tele43.Pair1
		public.Build43PairRanking = append([]ExperimentalV4PhonePairScore(nil), tele43.PairRanking...)
		public.HypothesesEvaluated += tele43.GeometryEvaluations
		if len(bank43) >= experimentalV4PhoneEnsembleSize {
			build42Bank = bank43
			hyp = append([]experimentalV4PhoneHypothesis(nil), bank43[:experimentalV4PhoneEnsembleSize]...)
			public.ProjectiveBasinFound = true
			public.Accepted = true
			public.EnsembleCandidates = len(hyp)
			q := bank43[0]
			public.ProposalScore = q.proposal
			public.ValidationScore = q.validation
			public.PilotScore = q.detection.Score
			public.PilotMargin = q.detection.Margin
			public.OriginXBlocks = q.detection.OriginXBlocks
			public.OriginYBlocks = q.detection.OriginYBlocks
		}
	}
	if !public.Accepted {
		// Qualified Build93: preserve the exact Build90 frozen geometry and
		// ordered-parallel qualification, then reuse one exact protected-grid residue
		// sample across this candidate's Robust/Balanced/Capacity derivations.
		payload64, info64, recovery64, err64 := experimentalV4PhoneBuild93Recover(work, boundary, key, cw, ch)
		experimentalV4PhoneBuild93ApplyTelemetry(&public, recovery64)
		if err64 == nil {
			public.ProjectiveBasinFound = true
			public.Accepted = true
			public.HMACAuthenticated = true
			return payload64, info64, public, nil
		}
		candidate := experimentalV4Prototype2Candidate()
		info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate), PilotScore: public.PilotScore, PilotMargin: public.PilotMargin, OriginXBlocks: public.OriginXBlocks, OriginYBlocks: public.OriginYBlocks}
		return nil, info, public, errors.New("experimental v4 phone geometry not accepted")
	}

	// Build42 ordering starts with the untouched Build41 global mapping. This is
	// both the cheapest accepted data path and the source of the Build41 physical
	// A/angle and B/front passes; residual fitting is unnecessary when it already
	// authenticates.
	payload, info, decode, baseErr := experimentalV4PhoneDecodeEnsemble(work, key, cw, ch, hyp)
	public.DataDecodeAttempted = decode.Attempted
	public.SoftHammingProfiles = decode.ProfilesTried
	public.MaxDataConfidence = decode.MaxConfidence
	public.HMACAuthenticated = decode.Authenticated
	if baseErr == nil {
		if public.Build43Attempted {
			public.Build43Authenticated = true
		}
		return payload, info, public, nil
	}

	// Build42 data fallback. Reuse the complete bank frozen and qualified during
	// the same Build41 geometry search, then enumerate only deterministic
	// three-geometry data ensembles and ML-ordered Hamming alternatives. HMAC is
	// final authentication only; it never creates or ranks geometry.
	public.Build42DataAttempted = true
	public.Build42BankCandidates = len(build42Bank)
	if len(build42Bank) >= 3 {
		payload42, info42, data42, err42 := experimentalV4PhoneDecodeBuild42Bank(work, key, cw, ch, build42Bank)
		public.Build42DataAttempted = data42.Attempted
		public.Build42BankCandidates = data42.BankCandidates
		public.Build42EnsemblesTried = data42.EnsemblesTried
		public.Build42ListFramesTried = data42.ListFramesTried
		public.Build42DataAuthenticated = data42.Authenticated
		if err42 == nil {
			public.HMACAuthenticated = true
			if public.Build43Attempted {
				public.Build43Authenticated = true
			}
			return payload42, info42, public, nil
		}
	}

	// Build40 residual warp remains available as the final pilot-only geometry
	// fallback. It is now paid only when both the accepted global ensemble and
	// Build42's post-geometry data recovery fail.
	residualHyp := append([]experimentalV4PhoneHypothesis(nil), hyp...)
	plane := newPixelPlane(work)
	candidate := experimentalV4Prototype2Candidate()
	for i := range residualHyp {
		warp, residual := experimentalV4PhoneFitResidualPilotOnly(plane, candidate, cw, ch, residualHyp[i].h)
		residualHyp[i].residual = residual
		public.ResidualAttempted = true
		if residual.Fitted {
			public.ResidualFitted = true
		}
		public.HypothesesEvaluated += residual.HypothesesEvaluated
		if i == 0 || residual.ValidationAfter > public.ResidualValidationAfter {
			public.ResidualControls = residual.Controls
			public.ResidualRMSPixels = residual.RMSPixels
			public.ResidualProposalBefore = residual.ProposalBefore
			public.ResidualProposalAfter = residual.ProposalAfter
			public.ResidualValidationBefore = residual.ValidationBefore
			public.ResidualValidationAfter = residual.ValidationAfter
		}
		if warp != nil && residual.Applied {
			residualHyp[i].warp = warp
			residualHyp[i].proposal = residual.ProposalAfter
			residualHyp[i].validation = residual.ValidationAfter
			residualHyp[i].detection = experimentalV4PhoneResidualFullDetection(plane, candidate, cw, ch, residualHyp[i].h, warp)
			public.ResidualApplied = true
		}
	}
	if public.ResidualApplied {
		public.FallbackAttempted = true
		payloadR, infoR, residualDecode, errR := experimentalV4PhoneDecodeEnsemble(work, key, cw, ch, residualHyp)
		if residualDecode.ProfilesTried > public.SoftHammingProfiles {
			public.SoftHammingProfiles = residualDecode.ProfilesTried
		}
		if residualDecode.MaxConfidence > public.MaxDataConfidence {
			public.MaxDataConfidence = residualDecode.MaxConfidence
		}
		if errR == nil {
			public.HMACAuthenticated = true
			public.FallbackAuthenticated = true
			if public.Build43Attempted {
				public.Build43Authenticated = true
			}
			return payloadR, infoR, public, nil
		}
	}

	// Build76 ordered-generation-two scheduling wraps the qualified Build75 deep-recovery fallback. This remains deliberately last so all
	// Build44-qualified paths remain untouched. The complete deep geometry bank
	// is generated proposal-only and frozen before held-out qualification or
	// HMAC are consulted. HMAC remains final frame authentication only.
	payload64, info64, recovery64, err64 := experimentalV4PhoneBuild75Recover(work, boundary, key, cw, ch)
	experimentalV4PhoneBuild75ApplyTelemetry(&public, recovery64)
	if err64 == nil {
		public.ProjectiveBasinFound = true
		public.Accepted = true
		public.HMACAuthenticated = true
		return payload64, info64, public, nil
	}
	return nil, info, public, baseErr
}
