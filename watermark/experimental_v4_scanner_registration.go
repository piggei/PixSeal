package watermark

import (
	"errors"
	"image"
	"math"
	"sort"
)

// Build37 scanner registration deliberately treats the visible artwork/paper
// boundary only as a geometric initializer. The public v4 pilot remains the
// watermark-specific evidence and payload/header/ECC/HMAC never participate in
// geometry proposal or refinement.
const (
	experimentalV4ScannerBoundaryMaxDimension = 2400
	experimentalV4ScannerForegroundThreshold  = 25.0
	experimentalV4ScannerPersistentRun        = 6
	experimentalV4ScannerBoundaryInset        = 2.5 // analysis-plane pixels
)

type experimentalV4ScannerLine struct {
	// Horizontal lines use v = intercept + slope*u. Vertical lines use
	// u = intercept + slope*v.
	intercept float64
	slope     float64
}

type experimentalV4ScannerRegistration struct {
	available  bool
	boundary   PrintBoundaryEstimate
	h          homography
	proposal   float64
	validation float64
	detection  ExperimentalV4PilotDetection
	hypotheses int
}

func experimentalV4ScannerBoundary(src image.Image) PrintBoundaryEstimate {
	plane := experimentalV4BuildScannerColorPlane(src, experimentalV4ScannerBoundaryMaxDimension)
	if plane == nil || plane.width < 300 || plane.height < 300 {
		return PrintBoundaryEstimate{}
	}
	paper, ok := estimateBoundaryPaperColor(plane)
	if !ok {
		return PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	}

	foreground := make([]bool, plane.width*plane.height)
	rowCounts := make([]int, plane.height)
	colCounts := make([]int, plane.width)
	for y := 0; y < plane.height; y++ {
		for x := 0; x < plane.width; x++ {
			if boundaryColorDistance(plane, paper, float64(x), float64(y)) <= experimentalV4ScannerForegroundThreshold {
				continue
			}
			foreground[y*plane.width+x] = true
			rowCounts[y]++
			colCounts[x]++
		}
	}

	// A printed photo occupies a large fraction of each row/column. Small text
	// labels above the artwork deliberately fail these occupancy floors.
	rowFloor := int(0.22 * float64(plane.width))
	colFloor := int(0.18 * float64(plane.height))
	top, bottom := -1, -1
	for y, count := range rowCounts {
		if count >= rowFloor {
			if top < 0 {
				top = y
			}
			bottom = y
		}
	}
	if top < 0 || bottom < 0 || bottom-top < plane.height/5 {
		return PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	}

	// Recompute column occupancy only across the coarse photo-height band.
	// This rejects scanner-bed/page-border streaks that span the full sheet.
	for x := range colCounts {
		colCounts[x] = 0
	}
	for y := top; y <= bottom; y++ {
		for x := plane.width / 20; x < plane.width*19/20; x++ {
			if foreground[y*plane.width+x] {
				colCounts[x]++
			}
		}
	}
	colFloor = int(0.55 * float64(bottom-top+1))
	left, right := -1, -1
	for x := plane.width / 20; x < plane.width*19/20; x++ {
		if colCounts[x] >= colFloor {
			if left < 0 {
				left = x
			}
			right = x
		}
	}
	if top < 0 || left < 0 || bottom-top < plane.height/5 || right-left < plane.width/5 {
		return PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	}

	padX := maxInt(16, (right-left)/18)
	padY := maxInt(16, (bottom-top)/18)
	topPoints := make([]ImagePoint, 0, right-left)
	bottomPoints := make([]ImagePoint, 0, right-left)
	for x := left + padX/3; x <= right-padX/3; x++ {
		if y, ok := scannerFirstPersistentVertical(foreground, plane.width, plane.height, x, maxInt(0, top-padY), minInt(plane.height-1, top+padY), true); ok {
			topPoints = append(topPoints, ImagePoint{X: float64(x), Y: float64(y)})
		}
		if y, ok := scannerFirstPersistentVertical(foreground, plane.width, plane.height, x, maxInt(0, bottom-padY), minInt(plane.height-1, bottom+padY), false); ok {
			bottomPoints = append(bottomPoints, ImagePoint{X: float64(x), Y: float64(y)})
		}
	}
	leftPoints := make([]ImagePoint, 0, bottom-top)
	rightPoints := make([]ImagePoint, 0, bottom-top)
	for y := top + padY/3; y <= bottom-padY/3; y++ {
		if x, ok := scannerFirstPersistentHorizontal(foreground, plane.width, plane.height, y, maxInt(0, left-padX), minInt(plane.width-1, left+padX), true); ok {
			leftPoints = append(leftPoints, ImagePoint{X: float64(x), Y: float64(y)})
		}
		if x, ok := scannerFirstPersistentHorizontal(foreground, plane.width, plane.height, y, maxInt(0, right-padX), minInt(plane.width-1, right+padX), false); ok {
			rightPoints = append(rightPoints, ImagePoint{X: float64(x), Y: float64(y)})
		}
	}

	topLine, topScore, okTop := experimentalV4ScannerRobustLine(topPoints, true)
	bottomLine, bottomScore, okBottom := experimentalV4ScannerRobustLine(bottomPoints, true)
	leftLine, leftScore, okLeft := experimentalV4ScannerRobustLine(leftPoints, false)
	rightLine, rightScore, okRight := experimentalV4ScannerRobustLine(rightPoints, false)
	result := PrintBoundaryEstimate{TopScore: topScore, RightScore: rightScore, BottomScore: bottomScore, LeftScore: leftScore, AnalysisDivisor: plane.divisor}
	if !(okTop && okBottom && okLeft && okRight) {
		return result
	}

	// Thresholding sees a small halo just outside the actual printed artwork.
	// Move all four fitted edges inward by a fixed analysis-plane distance.
	topLine.intercept += experimentalV4ScannerBoundaryInset
	bottomLine.intercept -= experimentalV4ScannerBoundaryInset
	leftLine.intercept += experimentalV4ScannerBoundaryInset
	rightLine.intercept -= experimentalV4ScannerBoundaryInset

	tl, okTL := experimentalV4ScannerIntersect(topLine, leftLine)
	tr, okTR := experimentalV4ScannerIntersect(topLine, rightLine)
	bl, okBL := experimentalV4ScannerIntersect(bottomLine, leftLine)
	br, okBR := experimentalV4ScannerIntersect(bottomLine, rightLine)
	if !(okTL && okTR && okBL && okBR) {
		return result
	}
	if !boundaryQuadSane(plane, tl, tr, br, bl) {
		return result
	}

	divisor := float64(plane.divisor)
	result.TopLeft = ImagePoint{X: tl.X * divisor, Y: tl.Y * divisor}
	result.TopRight = ImagePoint{X: tr.X * divisor, Y: tr.Y * divisor}
	result.BottomLeft = ImagePoint{X: bl.X * divisor, Y: bl.Y * divisor}
	result.BottomRight = ImagePoint{X: br.X * divisor, Y: br.Y * divisor}
	fitScore := math.Min(math.Min(topScore, bottomScore), math.Min(leftScore, rightScore))
	result.Confidence = clampUnit(fitScore)
	result.Detected = fitScore >= 0.45
	return result
}

func experimentalV4BuildScannerColorPlane(src image.Image, maxDimension int) *boundaryColorPlane {
	if src == nil {
		return nil
	}
	bounds := src.Bounds()
	largest := bounds.Dx()
	if bounds.Dy() > largest {
		largest = bounds.Dy()
	}
	divisor := 1
	for (largest+divisor-1)/divisor > maxDimension {
		divisor *= 2
	}
	width := (bounds.Dx() + divisor - 1) / divisor
	height := (bounds.Dy() + divisor - 1) / divisor
	plane := &boundaryColorPlane{width: width, height: height, divisor: divisor, rgb: make([]uint8, width*height*3)}
	for y := 0; y < height; y++ {
		sy0, sy1 := y*divisor, minInt((y+1)*divisor, bounds.Dy())
		for x := 0; x < width; x++ {
			sx0, sx1 := x*divisor, minInt((x+1)*divisor, bounds.Dx())
			var rs, gs, bs, n int
			for sy := sy0; sy < sy1; sy++ {
				for sx := sx0; sx < sx1; sx++ {
					pixel := flattenedNRGBA(src.At(bounds.Min.X+sx, bounds.Min.Y+sy))
					rs += int(pixel.R)
					gs += int(pixel.G)
					bs += int(pixel.B)
					n++
				}
			}
			idx := (y*width + x) * 3
			if n > 0 {
				plane.rgb[idx] = uint8(rs / n)
				plane.rgb[idx+1] = uint8(gs / n)
				plane.rgb[idx+2] = uint8(bs / n)
			}
		}
	}
	return plane
}

func scannerFirstPersistentVertical(mask []bool, width, height, x, start, end int, forward bool) (int, bool) {
	if x < 0 || x >= width || start < 0 || end >= height || start > end {
		return 0, false
	}
	run := experimentalV4ScannerPersistentRun
	if forward {
		for y := start; y+run-1 <= end; y++ {
			good := true
			for k := 0; k < run; k++ {
				if !mask[(y+k)*width+x] {
					good = false
					break
				}
			}
			if good {
				return y, true
			}
		}
	} else {
		for y := end; y-run+1 >= start; y-- {
			good := true
			for k := 0; k < run; k++ {
				if !mask[(y-k)*width+x] {
					good = false
					break
				}
			}
			if good {
				return y, true
			}
		}
	}
	return 0, false
}

func scannerFirstPersistentHorizontal(mask []bool, width, height, y, start, end int, forward bool) (int, bool) {
	if y < 0 || y >= height || start < 0 || end >= width || start > end {
		return 0, false
	}
	run := experimentalV4ScannerPersistentRun
	if forward {
		for x := start; x+run-1 <= end; x++ {
			good := true
			for k := 0; k < run; k++ {
				if !mask[y*width+x+k] {
					good = false
					break
				}
			}
			if good {
				return x, true
			}
		}
	} else {
		for x := end; x-run+1 >= start; x-- {
			good := true
			for k := 0; k < run; k++ {
				if !mask[y*width+x-k] {
					good = false
					break
				}
			}
			if good {
				return x, true
			}
		}
	}
	return 0, false
}

func experimentalV4ScannerRobustLine(points []ImagePoint, horizontal bool) (experimentalV4ScannerLine, float64, bool) {
	if len(points) < 40 {
		return experimentalV4ScannerLine{}, 0, false
	}
	work := append([]ImagePoint(nil), points...)
	var line experimentalV4ScannerLine
	for iter := 0; iter < 4; iter++ {
		var ok bool
		line, ok = experimentalV4ScannerLeastSquares(work, horizontal)
		if !ok {
			return experimentalV4ScannerLine{}, 0, false
		}
		residuals := make([]float64, len(work))
		absResiduals := make([]float64, len(work))
		for i, p := range work {
			var r float64
			if horizontal {
				r = p.Y - (line.intercept + line.slope*p.X)
			} else {
				r = p.X - (line.intercept + line.slope*p.Y)
			}
			residuals[i] = r
			absResiduals[i] = math.Abs(r)
		}
		sort.Float64s(absResiduals)
		median := absResiduals[len(absResiduals)/2]
		limit := math.Max(1.25, 3.5*median)
		next := make([]ImagePoint, 0, len(work))
		for i, p := range work {
			if math.Abs(residuals[i]) <= limit {
				next = append(next, p)
			}
		}
		if len(next) < 40 || len(next) == len(work) {
			break
		}
		work = next
	}
	inlierFraction := float64(len(work)) / float64(len(points))
	score := clampUnit((inlierFraction - 0.45) / 0.5)
	return line, score, true
}

func experimentalV4ScannerLeastSquares(points []ImagePoint, horizontal bool) (experimentalV4ScannerLine, bool) {
	if len(points) < 2 {
		return experimentalV4ScannerLine{}, false
	}
	var sx, sy, sxx, sxy float64
	for _, p := range points {
		x, y := p.X, p.Y
		if !horizontal {
			x, y = p.Y, p.X
		}
		sx += x
		sy += y
		sxx += x * x
		sxy += x * y
	}
	n := float64(len(points))
	denom := n*sxx - sx*sx
	if math.Abs(denom) < 1e-9 {
		return experimentalV4ScannerLine{}, false
	}
	slope := (n*sxy - sx*sy) / denom
	intercept := (sy - slope*sx) / n
	return experimentalV4ScannerLine{intercept: intercept, slope: slope}, true
}

func experimentalV4ScannerIntersect(horizontal, vertical experimentalV4ScannerLine) (ImagePoint, bool) {
	// y = ah + bh*x ; x = av + bv*y
	denom := 1 - vertical.slope*horizontal.slope
	if math.Abs(denom) < 1e-9 {
		return ImagePoint{}, false
	}
	x := (vertical.intercept + vertical.slope*horizontal.intercept) / denom
	y := horizontal.intercept + horizontal.slope*x
	return ImagePoint{X: x, Y: y}, true
}

func experimentalV4HomographyForObservedQuad(canonicalWidth, canonicalHeight int, boundary PrintBoundaryEstimate) (homography, bool) {
	if canonicalWidth < 2 || canonicalHeight < 2 {
		return homography{}, false
	}
	src := [4][2]float64{{0, 0}, {float64(canonicalWidth - 1), 0}, {0, float64(canonicalHeight - 1)}, {float64(canonicalWidth - 1), float64(canonicalHeight - 1)}}
	dst := [4]ImagePoint{boundary.TopLeft, boundary.TopRight, boundary.BottomLeft, boundary.BottomRight}
	var a [8][9]float64
	for i := 0; i < 4; i++ {
		x, y := src[i][0], src[i][1]
		u, v := dst[i].X, dst[i].Y
		a[2*i] = [9]float64{x, y, 1, 0, 0, 0, -u * x, -u * y, u}
		a[2*i+1] = [9]float64{0, 0, 0, x, y, 1, -v * x, -v * y, v}
	}
	s, ok := solveLinear8(a)
	if !ok {
		return homography{}, false
	}
	return homography{h: [9]float64{s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], 1}}, true
}

func experimentalV4ScannerDomainAffine(width, height int, scaleX, scaleY, shearX, shearY, tx, ty float64) homography {
	cx, cy := float64(width-1)/2, float64(height-1)/2
	toOrigin := homography{h: [9]float64{1, 0, -cx, 0, 1, -cy, 0, 0, 1}}
	affine := homography{h: [9]float64{scaleX, shearX, tx, shearY, scaleY, ty, 0, 0, 1}}
	fromOrigin := homography{h: [9]float64{1, 0, cx, 0, 1, cy, 0, 0, 1}}
	return experimentalV4BlindMultiplyHomography(fromOrigin, experimentalV4BlindMultiplyHomography(affine, toOrigin))
}

func experimentalV4ScannerDecodeWithHomography(src image.Image, key []byte, canonicalWidth, canonicalHeight int, h homography) ([]byte, ExperimentalV4ExtractInfo, error) {
	candidate := experimentalV4Prototype2Candidate()
	detection := experimentalV4DetectPilotProjective(src, candidate, canonicalWidth, canonicalHeight, h)
	info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotScore: detection.Score, PilotMargin: detection.Margin, OriginXBlocks: detection.OriginXBlocks, OriginYBlocks: detection.OriginYBlocks, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate)}
	plane := newPixelPlane(src)
	for _, spec := range v3Profiles {
		margins, confidence, ok := experimentalV4ReadProtectedProjectiveMargins(plane, canonicalWidth, canonicalHeight, h, spec.codedBits)
		if !ok {
			continue
		}
		decoded := experimentalV4SoftHammingDecodeMargins(margins)
		raw := bitsToBytes(whiten(decoded, key, experimentalV4WhitenLabel))
		if payload, err := parseExperimentalV4Frame(raw, key, spec); err == nil {
			info.Profile = spec.profile
			info.Confidence = confidence
			return payload, info, nil
		}
	}
	return nil, info, errors.New("experimental v4 scanner payload authentication failed")
}

// experimentalV4ScannerPilotProposalScore is a cheap fixed-origin proposal
// score. It samples only four interior/reasonably separated tile repetitions
// and only pilot partition A. Full repeated-field scoring is reserved for the
// shortlist and partition B is held out from proposal.
func experimentalV4ScannerPilotProposalScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, canonicalWidth, canonicalHeight int, h homography) (float64, int) {
	if plane == nil {
		return math.Inf(-1), 0
	}
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	tilesX := blocksWide / experimentalV4TileWidthBlocks
	tilesY := blocksHigh / experimentalV4TileHeightBlocks
	if tilesX < 2 || tilesY < 2 {
		return math.Inf(-1), 0
	}
	coords := [][2]int{{1, 1}, {tilesX - 2, 1}, {1, tilesY - 2}, {tilesX - 2, tilesY - 2}}
	seen := map[[2]int]bool{}
	numerator, denominator := 0.0, 0.0
	visible := 0
	for _, c := range coords {
		if c[0] < 0 || c[1] < 0 || c[0] >= tilesX || c[1] >= tilesY || seen[c] {
			continue
		}
		seen[c] = true
		baseX := c[0] * experimentalV4TileWidthBlocks
		baseY := c[1] * experimentalV4TileHeightBlocks
		for i, pos := range candidate.positions {
			if i&1 != 0 {
				continue
			}
			px := pos % experimentalV4TileWidthBlocks
			py := pos / experimentalV4TileWidthBlocks
			v, ok := readProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)
			if !ok {
				continue
			}
			visible++
			numerator += float64(candidate.signs[i]) * v
			denominator += math.Abs(v)
		}
	}
	if denominator <= 0 {
		return math.Inf(-1), visible
	}
	return numerator / denominator, visible
}

const (
	experimentalV4ScannerCoarseKeep   = 64
	experimentalV4ScannerFineKeep     = 96
	experimentalV4ScannerEnsembleSize = 5

	experimentalV4ScannerProposalFloor    = 0.30
	experimentalV4ScannerValidationFloor  = 0.15
	experimentalV4ScannerPilotScoreFloor  = 0.15
	experimentalV4ScannerPilotMarginFloor = 0.05
)

type experimentalV4ScannerParams struct {
	scaleX float64
	scaleY float64
	shearX float64
	shearY float64
	shiftX float64
	shiftY float64
}

type experimentalV4ScannerHypothesis struct {
	params     experimentalV4ScannerParams
	h          homography
	proposal   float64
	validation float64
	detection  ExperimentalV4PilotDetection
}

// ExperimentalV4ScannerInfo reports only public boundary/pilot evidence from
// the Build37 scanner path. Payload bytes, frame structure and HMAC validity do
// not participate in these fields or in scanner geometry selection.
type ExperimentalV4ScannerInfo struct {
	BoundaryDetected    bool
	BoundaryConfidence  float64
	Accepted            bool
	ProposalScore       float64
	ValidationScore     float64
	PilotScore          float64
	PilotMargin         float64
	OriginXBlocks       int
	OriginYBlocks       int
	ScaleX              float64
	ScaleY              float64
	ShearX              float64
	ShearY              float64
	ShiftX              float64
	ShiftY              float64
	EnsembleCandidates  int
	HypothesesEvaluated int
}

func experimentalV4ScannerCompose(seed homography, canonicalWidth, canonicalHeight int, p experimentalV4ScannerParams) homography {
	return experimentalV4BlindMultiplyHomography(seed, experimentalV4ScannerDomainAffine(canonicalWidth, canonicalHeight, p.scaleX, p.scaleY, p.shearX, p.shearY, p.shiftX, p.shiftY))
}

func experimentalV4ScannerInsertProposal(bank []experimentalV4ScannerHypothesis, q experimentalV4ScannerHypothesis, limit int) []experimentalV4ScannerHypothesis {
	bank = append(bank, q)
	sort.Slice(bank, func(i, j int) bool {
		if bank[i].proposal == bank[j].proposal {
			return bank[i].validation > bank[j].validation
		}
		return bank[i].proposal > bank[j].proposal
	})
	if len(bank) > limit {
		bank = bank[:limit]
	}
	return bank
}

// experimentalV4ScannerSearch is deliberately scanner-specific. The visible
// paper/artwork boundary provides an independent coarse mapping, after which a
// tightly bounded affine correction is proposed with pilot partition A and
// ranked with disjoint pilot partition B. No key, payload/header bits, ECC or
// HMAC are available to this function.
func experimentalV4ScannerSearch(src image.Image, canonicalWidth, canonicalHeight int) (experimentalV4ScannerRegistration, []experimentalV4ScannerHypothesis) {
	if src == nil || canonicalWidth < experimentalV4TileWidthBlocks*blockSize || canonicalHeight < experimentalV4TileHeightBlocks*blockSize || canonicalWidth%blockSize != 0 || canonicalHeight%blockSize != 0 {
		return experimentalV4ScannerRegistration{}, nil
	}
	boundary := experimentalV4ScannerBoundary(src)
	result := experimentalV4ScannerRegistration{boundary: boundary}
	if !boundary.Detected {
		return result, nil
	}
	seed, ok := experimentalV4HomographyForObservedQuad(canonicalWidth, canonicalHeight, boundary)
	if !ok {
		return result, nil
	}
	plane := newPixelPlane(src)
	candidate := experimentalV4Prototype2Candidate()
	hypotheses := 0

	// Coarse scanner correction. These ranges compensate only small edge-fit,
	// platen and resampling errors; they are intentionally far narrower than the
	// general Build34 projective search envelope.
	scales := []float64{1.000, 1.002, 1.004, 1.006, 1.008, 1.010}
	shearX := []float64{-0.004, -0.003, -0.002, -0.001, 0}
	shearY := []float64{-0.002, -0.001, 0, 0.001, 0.002}
	shiftX := []float64{-1, 0, 1, 2, 3}
	shiftY := []float64{-3, -1.5, 0, 1.5, 3}
	coarse := make([]experimentalV4ScannerHypothesis, 0, experimentalV4ScannerCoarseKeep)
	for _, sx := range scales {
		for _, sy := range scales {
			for _, shx := range shearX {
				for _, shy := range shearY {
					for _, tx := range shiftX {
						for _, ty := range shiftY {
							p := experimentalV4ScannerParams{scaleX: sx, scaleY: sy, shearX: shx, shearY: shy, shiftX: tx, shiftY: ty}
							h := experimentalV4ScannerCompose(seed, canonicalWidth, canonicalHeight, p)
							proposal, visible := experimentalV4ScannerPilotProposalScore(plane, candidate, canonicalWidth, canonicalHeight, h)
							hypotheses++
							if visible < 96 || math.IsInf(proposal, -1) {
								continue
							}
							coarse = experimentalV4ScannerInsertProposal(coarse, experimentalV4ScannerHypothesis{params: p, h: h, proposal: proposal}, experimentalV4ScannerCoarseKeep)
						}
					}
				}
			}
		}
	}
	if len(coarse) == 0 {
		result.hypotheses = hypotheses
		return result, nil
	}
	for i := range coarse {
		coarse[i].validation, _ = experimentalV4PilotPartitionScore(plane, candidate, canonicalWidth, canonicalHeight, coarse[i].h, 1)
	}
	sort.Slice(coarse, func(i, j int) bool {
		if coarse[i].validation == coarse[j].validation {
			return coarse[i].proposal > coarse[j].proposal
		}
		return coarse[i].validation > coarse[j].validation
	})

	// Fine 6-D neighborhood around the best held-out coarse basin. Proposal is
	// still partition A only; partition B is used only after the top-A shortlist
	// has been frozen.
	base := coarse[0].params
	offsets := []float64{-1, -0.5, 0, 0.5, 1}
	fine := make([]experimentalV4ScannerHypothesis, 0, experimentalV4ScannerFineKeep)
	for _, ox := range offsets {
		for _, oy := range offsets {
			for _, oshx := range offsets {
				for _, oshy := range offsets {
					for _, otx := range offsets {
						for _, oty := range offsets {
							p := experimentalV4ScannerParams{
								scaleX: base.scaleX + ox*0.001,
								scaleY: base.scaleY + oy*0.001,
								shearX: base.shearX + oshx*0.0005,
								shearY: base.shearY + oshy*0.0005,
								shiftX: base.shiftX + otx*0.75,
								shiftY: base.shiftY + oty*0.75,
							}
							h := experimentalV4ScannerCompose(seed, canonicalWidth, canonicalHeight, p)
							proposal, visible := experimentalV4ScannerPilotProposalScore(plane, candidate, canonicalWidth, canonicalHeight, h)
							hypotheses++
							if visible < 96 || math.IsInf(proposal, -1) {
								continue
							}
							fine = experimentalV4ScannerInsertProposal(fine, experimentalV4ScannerHypothesis{params: p, h: h, proposal: proposal}, experimentalV4ScannerFineKeep)
						}
					}
				}
			}
		}
	}
	if len(fine) == 0 {
		result.hypotheses = hypotheses
		return result, nil
	}
	for i := range fine {
		fine[i].validation, _ = experimentalV4PilotPartitionScore(plane, candidate, canonicalWidth, canonicalHeight, fine[i].h, 1)
	}
	sort.Slice(fine, func(i, j int) bool {
		if fine[i].validation == fine[j].validation {
			return fine[i].proposal > fine[j].proposal
		}
		return fine[i].validation > fine[j].validation
	})

	// One complete-field origin competition is enough to establish that the
	// held-out winner really corresponds to absolute origin (0,0), rather than
	// a texture-induced cyclic alias. The ensemble members themselves are the
	// next held-out-ranked hypotheses in the same already-qualified basin.
	fine[0].detection = experimentalV4DetectPilotProjective(src, candidate, canonicalWidth, canonicalHeight, fine[0].h)
	win := fine[0]
	globalAccepted := win.proposal >= experimentalV4ScannerProposalFloor &&
		win.validation >= experimentalV4ScannerValidationFloor &&
		win.detection.Available && win.detection.OriginXBlocks == 0 && win.detection.OriginYBlocks == 0 &&
		win.detection.Score >= experimentalV4ScannerPilotScoreFloor && win.detection.Margin >= experimentalV4ScannerPilotMarginFloor

	accepted := make([]experimentalV4ScannerHypothesis, 0, experimentalV4ScannerEnsembleSize)
	if globalAccepted {
		for i := 0; i < len(fine) && len(accepted) < experimentalV4ScannerEnsembleSize; i++ {
			q := fine[i]
			if q.proposal < experimentalV4ScannerProposalFloor || q.validation < experimentalV4ScannerValidationFloor {
				continue
			}
			if i == 0 {
				q.detection = win.detection
			}
			accepted = append(accepted, q)
		}
	}

	result.available = true
	result.hypotheses = hypotheses
	if len(fine) > 0 {
		win := fine[0]
		result.h = win.h
		result.proposal = win.proposal
		result.validation = win.validation
		result.detection = win.detection
	}
	if len(accepted) >= experimentalV4ScannerEnsembleSize {
		win := accepted[0]
		result.h = win.h
		result.proposal = win.proposal
		result.validation = win.validation
		result.detection = win.detection
	}
	return result, accepted
}

func experimentalV4ScannerDecodeEnsemble(src image.Image, key []byte, canonicalWidth, canonicalHeight int, hypotheses []experimentalV4ScannerHypothesis) ([]byte, ExperimentalV4ExtractInfo, error) {
	candidate := experimentalV4Prototype2Candidate()
	info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate)}
	if len(hypotheses) < experimentalV4ScannerEnsembleSize {
		return nil, info, errors.New("experimental v4 scanner geometry ensemble not accepted")
	}
	win := hypotheses[0]
	info.PilotScore = win.detection.Score
	info.PilotMargin = win.detection.Margin
	info.OriginXBlocks = win.detection.OriginXBlocks
	info.OriginYBlocks = win.detection.OriginYBlocks
	plane := newPixelPlane(src)
	for _, spec := range v3Profiles {
		var aggregate []float64
		weightSum := 0.0
		confidenceSum := 0.0
		for _, q := range hypotheses[:experimentalV4ScannerEnsembleSize] {
			margins, confidence, ok := experimentalV4ReadProtectedProjectiveMargins(plane, canonicalWidth, canonicalHeight, q.h, spec.codedBits)
			if !ok {
				continue
			}
			if aggregate == nil {
				aggregate = make([]float64, len(margins))
			}
			weight := math.Max(0.001, q.validation)
			for i := range margins {
				aggregate[i] += weight * margins[i]
			}
			weightSum += weight
			confidenceSum += weight * confidence
		}
		if weightSum <= 0 || len(aggregate) == 0 {
			continue
		}
		for i := range aggregate {
			aggregate[i] /= weightSum
		}
		softDecoded := experimentalV4SoftHammingDecodeMargins(aggregate)
		softRaw := bitsToBytes(whiten(softDecoded, key, experimentalV4WhitenLabel))
		if payload, err := parseExperimentalV4Frame(softRaw, key, spec); err == nil {
			info.Profile = spec.profile
			info.Confidence = confidenceSum / weightSum
			return payload, info, nil
		}
		coded := make([]byte, len(aggregate))
		for i, margin := range aggregate {
			if margin >= 0 {
				coded[i] = 1
			}
		}
		hardRaw := bitsToBytes(whiten(hammingDecode(coded), key, experimentalV4WhitenLabel))
		if payload, err := parseExperimentalV4Frame(hardRaw, key, spec); err == nil {
			info.Profile = spec.profile
			info.Confidence = confidenceSum / weightSum
			return payload, info, nil
		}
	}
	return nil, info, errors.New("experimental v4 scanner payload authentication failed")
}

// ExperimentalV4ExtractScanner is the Build37 blind scanner path. It requires
// a full-page/white-background scan containing the printed artwork. The visible
// artwork boundary seeds a tightly bounded affine refinement. Geometry is
// selected with disjoint halves of the public v4 pilot; only after a five-way
// pilot-qualified ensemble is frozen are protected data margins aggregated and
// authenticated with the unchanged v4 frame/HMAC.
func ExperimentalV4ExtractScanner(src image.Image, key []byte, canonicalWidth, canonicalHeight int) ([]byte, ExperimentalV4ExtractInfo, ExperimentalV4ScannerInfo, error) {
	if src == nil {
		return nil, ExperimentalV4ExtractInfo{}, ExperimentalV4ScannerInfo{}, errors.New("nil image")
	}
	if len(key) < 8 {
		return nil, ExperimentalV4ExtractInfo{}, ExperimentalV4ScannerInfo{}, errors.New("key must contain at least 8 bytes")
	}
	if canonicalWidth < experimentalV4TileWidthBlocks*blockSize || canonicalHeight < experimentalV4TileHeightBlocks*blockSize || canonicalWidth%blockSize != 0 || canonicalHeight%blockSize != 0 {
		return nil, ExperimentalV4ExtractInfo{}, ExperimentalV4ScannerInfo{}, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}
	search, ensemble := experimentalV4ScannerSearch(src, canonicalWidth, canonicalHeight)
	public := ExperimentalV4ScannerInfo{
		BoundaryDetected:    search.boundary.Detected,
		BoundaryConfidence:  search.boundary.Confidence,
		Accepted:            len(ensemble) >= experimentalV4ScannerEnsembleSize,
		ProposalScore:       search.proposal,
		ValidationScore:     search.validation,
		PilotScore:          search.detection.Score,
		PilotMargin:         search.detection.Margin,
		OriginXBlocks:       search.detection.OriginXBlocks,
		OriginYBlocks:       search.detection.OriginYBlocks,
		EnsembleCandidates:  len(ensemble),
		HypothesesEvaluated: search.hypotheses,
	}
	if len(ensemble) > 0 {
		p := ensemble[0].params
		public.ScaleX, public.ScaleY = p.scaleX, p.scaleY
		public.ShearX, public.ShearY = p.shearX, p.shearY
		public.ShiftX, public.ShiftY = p.shiftX, p.shiftY
	}
	if len(ensemble) < experimentalV4ScannerEnsembleSize {
		candidate := experimentalV4Prototype2Candidate()
		info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotScore: search.detection.Score, PilotMargin: search.detection.Margin, OriginXBlocks: search.detection.OriginXBlocks, OriginYBlocks: search.detection.OriginYBlocks, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate)}
		return nil, info, public, errors.New("experimental v4 scanner geometry not accepted")
	}
	payload, info, err := experimentalV4ScannerDecodeEnsemble(src, key, canonicalWidth, canonicalHeight, ensemble)
	return payload, info, public, err
}
