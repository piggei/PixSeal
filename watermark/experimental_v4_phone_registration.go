package watermark

import (
	"errors"
	"image"
	"image/color"
	"math"
	"sort"
)

// Build39/40 extend the Build37 physical decoder to single-shot smartphone
// captures. The phone path is deliberately limited to public evidence:
// visible paper/artwork geometry plus disjoint halves of the public v4 pilot.
// Payload/header/ECC/HMAC do not participate in geometry proposal/refinement.
const (
	experimentalV4PhoneMaxDimension        = 4600
	experimentalV4PhoneBoundaryDimension   = 2000
	experimentalV4PhoneForegroundThreshold = 22.0
	experimentalV4PhoneEnsembleSize        = 2
	experimentalV4PhoneProposalFloor       = 0.16
	experimentalV4PhoneValidationFloor     = 0.10
	experimentalV4PhonePilotScoreFloor     = 0.12
	experimentalV4PhonePilotMarginFloor    = 0.025
)

type ExperimentalV4PhonePairScore struct {
	Pair  string  `json:"pair"`
	Score float64 `json:"score"`
}

type ExperimentalV4PhoneInfo struct {
	WorkingWidth               int
	WorkingHeight              int
	Downsampled                bool
	BoundaryDetected           bool
	BoundaryConfidence         float64
	ProjectiveBasinFound       bool
	Accepted                   bool
	ProposalScore              float64
	ValidationScore            float64
	PilotScore                 float64
	PilotMargin                float64
	OriginXBlocks              int
	OriginYBlocks              int
	EnsembleCandidates         int
	HypothesesEvaluated        int
	ResidualAttempted          bool
	ResidualFitted             bool
	ResidualApplied            bool
	ResidualControls           int
	ResidualRMSPixels          float64
	ResidualProposalBefore     float64
	ResidualProposalAfter      float64
	ResidualValidationBefore   float64
	ResidualValidationAfter    float64
	DataDecodeAttempted        bool
	SoftHammingProfiles        int
	MaxDataConfidence          float64
	HMACAuthenticated          bool
	FallbackAttempted          bool
	FallbackAuthenticated      bool
	Build42DataAttempted       bool
	Build42BankCandidates      int
	Build42EnsemblesTried      int
	Build42ListFramesTried     int
	Build42DataAuthenticated   bool
	Build41DirectAccepted      bool
	Build41QualifiedCandidates int
	Build43Attempted           bool
	Build43EdgeRefined         bool
	Build43PairsScanned        int
	Build43PairsSelected       int
	Build43GeometryEvaluations int
	Build43ProposalCandidates  int
	Build43FrozenCandidates    int
	Build43QualifiedCandidates int
	Build43Pair0               string
	Build43Pair1               string
	Build43Authenticated       bool
	Build43PairRanking         []ExperimentalV4PhonePairScore
	Build64Attempted           bool
	Build64SeedsSelected       int
	Build64GeometryEvaluations int
	Build64BankCandidates      int
	Build64QualifiedCandidates int
	Build64DecodeCandidates    int
	Build64ListFramesTried     int
	Build64MaxDataConfidence   float64
	Build64Authenticated       bool
	Build65Attempted           bool
	Build65Workers             int
	Build66Attempted           bool
	Build66DecodeWorkers       int
}

type experimentalV4PhoneDecodeTelemetry struct {
	Attempted     bool
	ProfilesTried int
	MaxConfidence float64
	Authenticated bool
	Profile       Profile
}

type experimentalV4PhoneHypothesis struct {
	quad       [4]ImagePoint // TL, TR, BL, BR in working-image coordinates
	h          homography
	proposal   float64
	validation float64
	detection  ExperimentalV4PilotDetection
	warp       *experimentalV4PhoneResidualWarp
	residual   experimentalV4PhoneResidualInfo

	// Diagnostic provenance only. These fields are never consulted by
	// production ranking, qualification, data decode or authentication.
	build43Pair     string
	build43PairRank int
}

func experimentalV4PhoneResize(src image.Image, maxDimension int) (image.Image, bool) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	largest := w
	if h > largest {
		largest = h
	}
	if largest <= maxDimension {
		return src, false
	}
	scale := float64(maxDimension) / float64(largest)
	nw := int(math.Round(float64(w) * scale))
	nh := int(math.Round(float64(h) * scale))
	if nw < 2 {
		nw = 2
	}
	if nh < 2 {
		nh = 2
	}
	dst := image.NewNRGBA(image.Rect(0, 0, nw, nh))
	sx := float64(w) / float64(nw)
	sy := float64(h) / float64(nh)
	for y := 0; y < nh; y++ {
		fy := (float64(y)+0.5)*sy - 0.5
		y0 := int(math.Floor(fy))
		dy := fy - float64(y0)
		if y0 < 0 {
			y0 = 0
			dy = 0
		}
		y1 := y0 + 1
		if y1 >= h {
			y1 = h - 1
		}
		for x := 0; x < nw; x++ {
			fx := (float64(x)+0.5)*sx - 0.5
			x0 := int(math.Floor(fx))
			dx := fx - float64(x0)
			if x0 < 0 {
				x0 = 0
				dx = 0
			}
			x1 := x0 + 1
			if x1 >= w {
				x1 = w - 1
			}
			c00 := color.NRGBAModel.Convert(src.At(b.Min.X+x0, b.Min.Y+y0)).(color.NRGBA)
			c10 := color.NRGBAModel.Convert(src.At(b.Min.X+x1, b.Min.Y+y0)).(color.NRGBA)
			c01 := color.NRGBAModel.Convert(src.At(b.Min.X+x0, b.Min.Y+y1)).(color.NRGBA)
			c11 := color.NRGBAModel.Convert(src.At(b.Min.X+x1, b.Min.Y+y1)).(color.NRGBA)
			interp := func(a, b, c, d uint8) uint8 {
				top := float64(a)*(1-dx) + float64(b)*dx
				bot := float64(c)*(1-dx) + float64(d)*dx
				v := top*(1-dy) + bot*dy
				if v < 0 {
					v = 0
				}
				if v > 255 {
					v = 255
				}
				return uint8(v + 0.5)
			}
			dst.SetNRGBA(x, y, color.NRGBA{R: interp(c00.R, c10.R, c01.R, c11.R), G: interp(c00.G, c10.G, c01.G, c11.G), B: interp(c00.B, c10.B, c01.B, c11.B), A: 255})
		}
	}
	return dst, true
}

func experimentalV4PhoneLocalPaperBoundary(plane *boundaryColorPlane) PrintBoundaryEstimate {
	result := PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	w, h := plane.width, plane.height
	bandX := maxInt(8, w/30)
	bandY := maxInt(8, h/30)
	const threshold = 34.0
	const run = 8
	rgbAt := func(x, y int) (float64, float64, float64) {
		i := (y*w + x) * 3
		return float64(plane.rgb[i]), float64(plane.rgb[i+1]), float64(plane.rgb[i+2])
	}
	dist := func(r, g, b, rr, gg, bb float64) float64 {
		dr := r - rr
		dg := g - gg
		db := b - bb
		return math.Sqrt(dr*dr + dg*dg + db*db)
	}
	leftPts, rightPts := make([]ImagePoint, 0, h), make([]ImagePoint, 0, h)
	for y := bandY; y < h-bandY; y++ {
		var lr, lg, lb, rr, rg, rb float64
		for x := 0; x < bandX; x++ {
			r, g, b := rgbAt(x, y)
			lr += r
			lg += g
			lb += b
			r, g, b = rgbAt(w-1-x, y)
			rr += r
			rg += g
			rb += b
		}
		d := float64(bandX)
		lr /= d
		lg /= d
		lb /= d
		rr /= d
		rg /= d
		rb /= d
		first, last := -1, -1
		for x := bandX; x < w-bandX-run; x++ {
			ok := true
			for k := 0; k < run; k++ {
				r, g, b := rgbAt(x+k, y)
				if dist(r, g, b, lr, lg, lb) < threshold {
					ok = false
					break
				}
			}
			if ok {
				first = x
				break
			}
		}
		for x := w - bandX - 1; x >= bandX+run; x-- {
			ok := true
			for k := 0; k < run; k++ {
				r, g, b := rgbAt(x-k, y)
				if dist(r, g, b, rr, rg, rb) < threshold {
					ok = false
					break
				}
			}
			if ok {
				last = x
				break
			}
		}
		if first >= 0 && last > first+w/5 {
			leftPts = append(leftPts, ImagePoint{X: float64(first), Y: float64(y)})
			rightPts = append(rightPts, ImagePoint{X: float64(last), Y: float64(y)})
		}
	}
	topPts, bottomPts := make([]ImagePoint, 0, w), make([]ImagePoint, 0, w)
	for x := bandX; x < w-bandX; x++ {
		var tr, tg, tb, br, bg, bb float64
		for y := 0; y < bandY; y++ {
			r, g, b := rgbAt(x, y)
			tr += r
			tg += g
			tb += b
			r, g, b = rgbAt(x, h-1-y)
			br += r
			bg += g
			bb += b
		}
		d := float64(bandY)
		tr /= d
		tg /= d
		tb /= d
		br /= d
		bg /= d
		bb /= d
		first, last := -1, -1
		for y := bandY; y < h-bandY-run; y++ {
			ok := true
			for k := 0; k < run; k++ {
				r, g, b := rgbAt(x, y+k)
				if dist(r, g, b, tr, tg, tb) < threshold {
					ok = false
					break
				}
			}
			if ok {
				first = y
				break
			}
		}
		for y := h - bandY - 1; y >= bandY+run; y-- {
			ok := true
			for k := 0; k < run; k++ {
				r, g, b := rgbAt(x, y-k)
				if dist(r, g, b, br, bg, bb) < threshold {
					ok = false
					break
				}
			}
			if ok {
				last = y
				break
			}
		}
		if first >= 0 && last > first+h/5 {
			topPts = append(topPts, ImagePoint{X: float64(x), Y: float64(first)})
			bottomPts = append(bottomPts, ImagePoint{X: float64(x), Y: float64(last)})
		}
	}
	topLine, ts, okT := experimentalV4ScannerRobustLine(topPts, true)
	bottomLine, bs, okB := experimentalV4ScannerRobustLine(bottomPts, true)
	leftLine, ls, okL := experimentalV4ScannerRobustLine(leftPts, false)
	rightLine, rs, okR := experimentalV4ScannerRobustLine(rightPts, false)
	result.TopScore, result.BottomScore, result.LeftScore, result.RightScore = ts, bs, ls, rs
	if !(okT && okB && okL && okR) {
		return result
	}
	tl, o1 := experimentalV4ScannerIntersect(topLine, leftLine)
	tr, o2 := experimentalV4ScannerIntersect(topLine, rightLine)
	bl, o3 := experimentalV4ScannerIntersect(bottomLine, leftLine)
	br, o4 := experimentalV4ScannerIntersect(bottomLine, rightLine)
	if !(o1 && o2 && o3 && o4) || !boundaryQuadSane(plane, tl, tr, br, bl) {
		return result
	}
	d := float64(plane.divisor)
	result.TopLeft = ImagePoint{X: tl.X * d, Y: tl.Y * d}
	result.TopRight = ImagePoint{X: tr.X * d, Y: tr.Y * d}
	result.BottomLeft = ImagePoint{X: bl.X * d, Y: bl.Y * d}
	result.BottomRight = ImagePoint{X: br.X * d, Y: br.Y * d}
	fit := math.Min(math.Min(ts, bs), math.Min(ls, rs))
	result.Confidence = clampUnit(fit)
	result.Detected = fit >= 0.30
	// Reject a page/background fit. A phone artwork must leave visible paper on all sides.
	guardX := 0.015 * float64(w*plane.divisor)
	guardY := 0.015 * float64(h*plane.divisor)
	pts := []ImagePoint{result.TopLeft, result.TopRight, result.BottomLeft, result.BottomRight}
	for _, p := range pts {
		if p.X < guardX || p.X > float64(w*plane.divisor)-guardX || p.Y < guardY || p.Y > float64(h*plane.divisor)-guardY {
			result.Detected = false
			break
		}
	}
	return result
}

func experimentalV4PhoneGradientBoundary(plane *boundaryColorPlane, rough PrintBoundaryEstimate) PrintBoundaryEstimate {
	result := PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	if !rough.Detected {
		return result
	}
	d := float64(plane.divisor)
	q := [4]ImagePoint{{rough.TopLeft.X / d, rough.TopLeft.Y / d}, {rough.TopRight.X / d, rough.TopRight.Y / d}, {rough.BottomLeft.X / d, rough.BottomLeft.Y / d}, {rough.BottomRight.X / d, rough.BottomRight.Y / d}}
	w, h := plane.width, plane.height
	lum := func(x, y int) float64 {
		i := (y*w + x) * 3
		return 0.299*float64(plane.rgb[i]) + 0.587*float64(plane.rgb[i+1]) + 0.114*float64(plane.rgb[i+2])
	}
	verticalGrad := func(x, y int) float64 {
		if y < 2 || y+2 >= h {
			return 0
		}
		return math.Abs(lum(x, y+2) - lum(x, y-2))
	}
	interpY := func(a, b ImagePoint, x float64) float64 {
		den := b.X - a.X
		if math.Abs(den) < 1 {
			return (a.Y + b.Y) / 2
		}
		t := (x - a.X) / den
		return a.Y + t*(b.Y-a.Y)
	}
	topPts, bottomPts := make([]ImagePoint, 0, w), make([]ImagePoint, 0, w)
	minX := maxInt(4, int(math.Max(math.Min(q[0].X, q[2].X), 0)))
	maxX := minInt(w-5, int(math.Min(math.Max(q[1].X, q[3].X), float64(w-1))))
	radY := maxInt(30, h/12)
	for x := minX; x <= maxX; x += 2 {
		for side := 0; side < 2; side++ {
			var pred float64
			if side == 0 {
				pred = interpY(q[0], q[1], float64(x))
			} else {
				pred = interpY(q[2], q[3], float64(x))
			}
			lo := maxInt(3, int(pred)-radY)
			hi := minInt(h-4, int(pred)+radY)
			bestY := -1
			best := 18.0
			for y := lo; y <= hi; y++ {
				g := verticalGrad(x, y)
				if g > best {
					best = g
					bestY = y
				}
			}
			if bestY >= 0 {
				p := ImagePoint{X: float64(x), Y: float64(bestY)}
				if side == 0 {
					topPts = append(topPts, p)
				} else {
					bottomPts = append(bottomPts, p)
				}
			}
		}
	}
	leftDen := q[2].Y - q[0].Y
	rightDen := q[3].Y - q[1].Y
	if math.Abs(leftDen) < 1 || math.Abs(rightDen) < 1 {
		return result
	}
	leftSlope := (q[2].X - q[0].X) / leftDen
	rightSlope := (q[3].X - q[1].X) / rightDen
	leftLine := experimentalV4ScannerLine{intercept: q[0].X - leftSlope*q[0].Y, slope: leftSlope}
	rightLine := experimentalV4ScannerLine{intercept: q[1].X - rightSlope*q[1].Y, slope: rightSlope}
	topLine, ts, okT := experimentalV4ScannerRobustLine(topPts, true)
	bottomLine, bs, okB := experimentalV4ScannerRobustLine(bottomPts, true)
	ls, rs, okL, okR := rough.LeftScore, rough.RightScore, true, true
	result.TopScore, result.BottomScore, result.LeftScore, result.RightScore = ts, bs, ls, rs
	if !(okT && okB && okL && okR) {
		return result
	}
	tl, o1 := experimentalV4ScannerIntersect(topLine, leftLine)
	tr, o2 := experimentalV4ScannerIntersect(topLine, rightLine)
	bl, o3 := experimentalV4ScannerIntersect(bottomLine, leftLine)
	br, o4 := experimentalV4ScannerIntersect(bottomLine, rightLine)
	if !(o1 && o2 && o3 && o4) || !boundaryQuadSane(plane, tl, tr, br, bl) {
		return result
	}
	result.TopLeft = ImagePoint{X: tl.X * d, Y: tl.Y * d}
	result.TopRight = ImagePoint{X: tr.X * d, Y: tr.Y * d}
	result.BottomLeft = ImagePoint{X: bl.X * d, Y: bl.Y * d}
	result.BottomRight = ImagePoint{X: br.X * d, Y: br.Y * d}
	fit := math.Min(math.Min(ts, bs), math.Min(ls, rs))
	result.Confidence = clampUnit(fit)
	result.Detected = fit >= 0.30
	return result
}

func experimentalV4PhoneBoundary(src image.Image) PrintBoundaryEstimate {
	// Scanner boundary is excellent for front-facing phone shots too. Keep it
	// as the first choice, then fall back to a perspective-tolerant component
	// boundary for stronger keystone.
	if q := experimentalV4ScannerBoundary(src); q.Detected {
		b := src.Bounds()
		guardX := 0.02 * float64(b.Dx())
		guardY := 0.02 * float64(b.Dy())
		pts := []ImagePoint{q.TopLeft, q.TopRight, q.BottomLeft, q.BottomRight}
		inside := true
		for _, p := range pts {
			if p.X < guardX || p.X > float64(b.Dx())-guardX || p.Y < guardY || p.Y > float64(b.Dy())-guardY {
				inside = false
				break
			}
		}
		if inside {
			return q
		}
	}
	plane := experimentalV4BuildScannerColorPlane(src, experimentalV4PhoneBoundaryDimension)
	if plane == nil || plane.width < 300 || plane.height < 300 {
		return PrintBoundaryEstimate{}
	}
	if q := experimentalV4PhoneLocalPaperBoundary(plane); q.Detected {
		if refined := experimentalV4PhoneGradientBoundary(plane, q); refined.Detected {
			return refined
		}
		return q
	}
	paper, ok := estimateBoundaryPaperColor(plane)
	if !ok {
		return PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	}
	n := plane.width * plane.height
	raw := make([]bool, n)
	marginX, marginY := plane.width/50, plane.height/50
	for y := marginY; y < plane.height-marginY; y++ {
		for x := marginX; x < plane.width-marginX; x++ {
			if boundaryColorDistance(plane, paper, float64(x), float64(y)) > experimentalV4PhoneForegroundThreshold {
				raw[y*plane.width+x] = true
			}
		}
	}
	// Small morphological dilation joins texture islands inside the printed
	// photograph but does not connect the separate filename label to it.
	mask := append([]bool(nil), raw...)
	for pass := 0; pass < 2; pass++ {
		next := append([]bool(nil), mask...)
		for y := 1; y < plane.height-1; y++ {
			for x := 1; x < plane.width-1; x++ {
				if mask[y*plane.width+x] {
					continue
				}
				hit := false
				for dy := -1; dy <= 1 && !hit; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if mask[(y+dy)*plane.width+x+dx] {
							hit = true
							break
						}
					}
				}
				if hit {
					next[y*plane.width+x] = true
				}
			}
		}
		mask = next
	}
	seen := make([]bool, n)
	best := make([]int, 0)
	queue := make([]int, 0, 4096)
	for y := marginY; y < plane.height-marginY; y++ {
		for x := marginX; x < plane.width-marginX; x++ {
			idx := y*plane.width + x
			if !mask[idx] || seen[idx] {
				continue
			}
			queue = queue[:0]
			queue = append(queue, idx)
			seen[idx] = true
			comp := make([]int, 0, 1024)
			minx, maxx, miny, maxy := x, x, y, y
			for head := 0; head < len(queue); head++ {
				p := queue[head]
				comp = append(comp, p)
				px := p % plane.width
				py := p / plane.width
				if px < minx {
					minx = px
				}
				if px > maxx {
					maxx = px
				}
				if py < miny {
					miny = py
				}
				if py > maxy {
					maxy = py
				}
				if px > 0 {
					q := p - 1
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
				if px+1 < plane.width {
					q := p + 1
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
				if py > 0 {
					q := p - plane.width
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
				if py+1 < plane.height {
					q := p + plane.width
					if mask[q] && !seen[q] {
						seen[q] = true
						queue = append(queue, q)
					}
				}
			}
			if maxx-minx < plane.width/6 || maxy-miny < plane.height/6 {
				continue
			}
			if len(comp) > len(best) {
				best = append(best[:0], comp...)
			}
		}
	}
	result := PrintBoundaryEstimate{AnalysisDivisor: plane.divisor}
	if len(best) < n/80 {
		return result
	}
	inComp := make([]bool, n)
	for _, p := range best {
		inComp[p] = true
	}
	topPts, bottomPts, leftPts, rightPts := make([]ImagePoint, 0), make([]ImagePoint, 0), make([]ImagePoint, 0), make([]ImagePoint, 0)
	for x := marginX; x < plane.width-marginX; x++ {
		top, bottom := -1, -1
		for y := marginY; y < plane.height-marginY; y++ {
			if inComp[y*plane.width+x] {
				if top < 0 {
					top = y
				}
				bottom = y
			}
		}
		if top >= 0 {
			topPts = append(topPts, ImagePoint{X: float64(x), Y: float64(top)})
			bottomPts = append(bottomPts, ImagePoint{X: float64(x), Y: float64(bottom)})
		}
	}
	for y := marginY; y < plane.height-marginY; y++ {
		left, right := -1, -1
		for x := marginX; x < plane.width-marginX; x++ {
			if inComp[y*plane.width+x] {
				if left < 0 {
					left = x
				}
				right = x
			}
		}
		if left >= 0 {
			leftPts = append(leftPts, ImagePoint{X: float64(left), Y: float64(y)})
			rightPts = append(rightPts, ImagePoint{X: float64(right), Y: float64(y)})
		}
	}
	topLine, ts, okT := experimentalV4ScannerRobustLine(topPts, true)
	bottomLine, bs, okB := experimentalV4ScannerRobustLine(bottomPts, true)
	leftLine, ls, okL := experimentalV4ScannerRobustLine(leftPts, false)
	rightLine, rs, okR := experimentalV4ScannerRobustLine(rightPts, false)
	result.TopScore, result.BottomScore, result.LeftScore, result.RightScore = ts, bs, ls, rs
	if !(okT && okB && okL && okR) {
		return result
	}
	tl, o1 := experimentalV4ScannerIntersect(topLine, leftLine)
	tr, o2 := experimentalV4ScannerIntersect(topLine, rightLine)
	bl, o3 := experimentalV4ScannerIntersect(bottomLine, leftLine)
	br, o4 := experimentalV4ScannerIntersect(bottomLine, rightLine)
	if !(o1 && o2 && o3 && o4) || !boundaryQuadSane(plane, tl, tr, br, bl) {
		return result
	}
	d := float64(plane.divisor)
	result.TopLeft = ImagePoint{X: tl.X * d, Y: tl.Y * d}
	result.TopRight = ImagePoint{X: tr.X * d, Y: tr.Y * d}
	result.BottomLeft = ImagePoint{X: bl.X * d, Y: bl.Y * d}
	result.BottomRight = ImagePoint{X: br.X * d, Y: br.Y * d}
	fit := math.Min(math.Min(ts, bs), math.Min(ls, rs))
	result.Confidence = clampUnit(fit)
	result.Detected = fit >= 0.30
	return result
}

func experimentalV4PhoneQuadHomography(cw, ch int, q [4]ImagePoint) (homography, bool) {
	b := PrintBoundaryEstimate{Detected: true, TopLeft: q[0], TopRight: q[1], BottomLeft: q[2], BottomRight: q[3]}
	return experimentalV4HomographyForObservedQuad(cw, ch, b)
}

func experimentalV4PhoneQuadFromHomography(h homography, cw, ch int) ([4]ImagePoint, bool) {
	coords := [4][2]float64{{0, 0}, {float64(cw - 1), 0}, {0, float64(ch - 1)}, {float64(cw - 1), float64(ch - 1)}}
	var q [4]ImagePoint
	for i, c := range coords {
		x, y, ok := h.mapPoint(c[0], c[1])
		if !ok {
			return [4]ImagePoint{}, false
		}
		q[i] = ImagePoint{X: x, Y: y}
	}
	return q, true
}

func experimentalV4PhoneAverageBlockScale(q [4]ImagePoint, cw, ch int) float64 {
	dist := func(a, b ImagePoint) float64 { return math.Hypot(a.X-b.X, a.Y-b.Y) }
	sx := (dist(q[0], q[1]) + dist(q[2], q[3])) / (2 * float64(cw)) * blockSize
	sy := (dist(q[0], q[2]) + dist(q[1], q[3])) / (2 * float64(ch)) * blockSize
	return math.Max(2, (sx+sy)/2)
}

func experimentalV4PhoneSpatialPilotDetection(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, partition int, h homography) ExperimentalV4PilotDetection {
	if plane == nil {
		return ExperimentalV4PilotDetection{}
	}
	bw, bh := cw/blockSize, ch/blockSize
	if bw < experimentalV4TileWidthBlocks || bh < experimentalV4TileHeightBlocks {
		return ExperimentalV4PilotDetection{}
	}
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var sums [residues]float64
	var abss [residues]float64
	var counts [residues]int
	for by := 0; by < bh; by++ {
		ty := by / experimentalV4TileHeightBlocks
		for bx := 0; bx < bw; bx++ {
			tx := bx / experimentalV4TileWidthBlocks
			if (tx+ty)&1 != partition {
				continue
			}
			v, ok := readProjectiveBlockValue(plane, h, bx*blockSize, by*blockSize, blockSize)
			if !ok {
				continue
			}
			r := (by%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + (bx % experimentalV4TileWidthBlocks)
			sums[r] += v
			abss[r] += math.Abs(v)
			counts[r]++
		}
	}
	return experimentalV4DetectPilotFromResidues(candidate, sums[:], abss[:], counts[:])
}

func experimentalV4PhonePartitionOriginScore(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, h homography, partition int) (float64, int, int, int) {
	if plane == nil {
		return math.Inf(-1), 0, 0, 0
	}
	bw, bh := cw/blockSize, ch/blockSize
	txs, tys := bw/experimentalV4TileWidthBlocks, bh/experimentalV4TileHeightBlocks
	if txs < 2 || tys < 2 {
		return math.Inf(-1), 0, 0, 0
	}
	starts := [][2]int{{0, 0}, {(txs - 1) * experimentalV4TileWidthBlocks, 0}, {0, (tys - 1) * experimentalV4TileHeightBlocks}, {(txs - 1) * experimentalV4TileWidthBlocks, (tys - 1) * experimentalV4TileHeightBlocks}}
	residues := experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	sums := make([]float64, residues)
	abss := make([]float64, residues)
	counts := make([]int, residues)
	visibleSamples := 0
	seen := map[[2]int]bool{}
	for _, st := range starts {
		if seen[st] {
			continue
		}
		seen[st] = true
		for y := 0; y < experimentalV4TileHeightBlocks; y++ {
			for x := 0; x < experimentalV4TileWidthBlocks; x++ {
				v, ok := readProjectiveBlockValue(plane, h, (st[0]+x)*blockSize, (st[1]+y)*blockSize, blockSize)
				if !ok {
					continue
				}
				i := y*experimentalV4TileWidthBlocks + x
				sums[i] += v
				abss[i] += math.Abs(v)
				counts[i]++
				visibleSamples++
			}
		}
	}
	best := math.Inf(-1)
	bestX, bestY := 0, 0
	for oy := 0; oy < experimentalV4TileHeightBlocks; oy++ {
		for ox := 0; ox < experimentalV4TileWidthBlocks; ox++ {
			num, den := 0.0, 0.0
			for i, pos := range candidate.positions {
				if i&1 != partition {
					continue
				}
				px, py := pos%experimentalV4TileWidthBlocks, pos/experimentalV4TileWidthBlocks
				rx := positiveMod(px-ox, experimentalV4TileWidthBlocks)
				ry := positiveMod(py-oy, experimentalV4TileHeightBlocks)
				r := ry*experimentalV4TileWidthBlocks + rx
				if counts[r] == 0 {
					continue
				}
				num += float64(candidate.signs[i]) * sums[r]
				den += abss[r]
			}
			if den > 0 {
				v := num / den
				if v > best {
					best, bestX, bestY = v, ox, oy
				}
			}
		}
	}
	return best, bestX, bestY, visibleSamples
}

func experimentalV4PhoneRefineA(plane *pixelPlane, candidate experimentalV4PilotCandidate, cw, ch int, start [4]ImagePoint) (experimentalV4PhoneHypothesis, int) {
	q := start
	h, ok := experimentalV4PhoneQuadHomography(cw, ch, q)
	if !ok {
		return experimentalV4PhoneHypothesis{}, 0
	}
	det := experimentalV4PhoneSpatialPilotDetection(plane, candidate, cw, ch, 0, h)
	evals := 1
	if !det.Available {
		return experimentalV4PhoneHypothesis{}, evals
	}
	score := det.Score
	block := experimentalV4PhoneAverageBlockScale(q, cw, ch)
	steps := []float64{4 * block, 2 * block, block, 0.5 * block, 0.25 * block}
	apply := func(in [4]ImagePoint, dim int, delta float64) [4]ImagePoint {
		out := in
		switch dim {
		case 0:
			out[0].Y += delta
			out[1].Y += delta
		case 1:
			out[2].Y += delta
			out[3].Y += delta
		case 2:
			out[0].X += delta
		case 3:
			out[1].X += delta
		case 4:
			out[2].X += delta
		case 5:
			out[3].X += delta
		}
		return out
	}
	for _, step := range steps {
		for pass := 0; pass < 2; pass++ {
			improved := false
			for dim := 0; dim < 6; dim++ {
				bestQ, bestH, bestDet, best := q, h, det, score
				for _, sgn := range []float64{-1, 1} {
					qq := apply(q, dim, sgn*step)
					hh, ok := experimentalV4PhoneQuadHomography(cw, ch, qq)
					if !ok {
						continue
					}
					dd := experimentalV4PhoneSpatialPilotDetection(plane, candidate, cw, ch, 0, hh)
					evals++
					if dd.Available && dd.Score > best {
						best, bestQ, bestH, bestDet = dd.Score, qq, hh, dd
					}
				}
				if best > score+1e-7 {
					q, h, det, score = bestQ, bestH, bestDet, best
					improved = true
				}
			}
			if !improved {
				break
			}
		}
	}
	h = experimentalV4Build34CanonicalizeOrigin(h, det.OriginXBlocks, det.OriginYBlocks)
	cq, ok := experimentalV4PhoneQuadFromHomography(h, cw, ch)
	if ok {
		q = cq
	}
	return experimentalV4PhoneHypothesis{quad: q, h: h, proposal: score}, evals
}

func experimentalV4PhoneStarts(seed [4]ImagePoint, cw, ch int) [][4]ImagePoint {
	out := make([][4]ImagePoint, 0, 3)
	out = append(out, seed)
	cx, cy := 0.0, 0.0
	for _, p := range seed {
		cx += p.X
		cy += p.Y
	}
	cx /= 4
	cy /= 4
	block := experimentalV4PhoneAverageBlockScale(seed, cw, ch)
	for _, amount := range []float64{-1.0, 1.0} {
		q := seed
		for i, p := range q {
			dx, dy := p.X-cx, p.Y-cy
			n := math.Hypot(dx, dy)
			if n > 0 {
				q[i].X += amount * block * dx / n
				q[i].Y += amount * block * dy / n
			}
		}
		out = append(out, q)
	}
	return out
}

func experimentalV4PhoneSearch(src image.Image, cw, ch int) (image.Image, PrintBoundaryEstimate, []experimentalV4PhoneHypothesis, int, bool) {
	work, down := experimentalV4PhoneResize(src, experimentalV4PhoneMaxDimension)
	boundary := experimentalV4PhoneBoundary(work)
	if !boundary.Detected {
		return work, boundary, nil, 0, down
	}
	seed := [4]ImagePoint{boundary.TopLeft, boundary.TopRight, boundary.BottomLeft, boundary.BottomRight}
	plane := newPixelPlane(work)
	candidate := experimentalV4Prototype2Candidate()
	bank := make([]experimentalV4PhoneHypothesis, 0, 9)
	evals := 0
	for _, st := range experimentalV4PhoneStarts(seed, cw, ch) {
		q, n := experimentalV4PhoneRefineA(plane, candidate, cw, ch, st)
		evals += n
		if q.h.h[8] == 0 {
			continue
		}
		vd := experimentalV4PhoneSpatialPilotDetection(plane, candidate, cw, ch, 1, q.h)
		evals++
		if vd.Available && vd.OriginXBlocks == 0 && vd.OriginYBlocks == 0 {
			q.validation = vd.Score
		} else {
			q.validation = math.Inf(-1)
		}
		bank = append(bank, q)
	}
	sort.Slice(bank, func(i, j int) bool {
		if bank[i].validation == bank[j].validation {
			return bank[i].proposal > bank[j].proposal
		}
		return bank[i].validation > bank[j].validation
	})
	if len(bank) == 0 {
		return work, boundary, nil, evals, down
	}
	bank[0].detection = experimentalV4DetectPilotProjective(work, candidate, cw, ch, bank[0].h)
	win := bank[0]
	good := win.proposal >= experimentalV4PhoneProposalFloor && win.validation >= experimentalV4PhoneValidationFloor
	if !good {
		return work, boundary, bank[:1], evals, down
	}
	accepted := make([]experimentalV4PhoneHypothesis, 0, experimentalV4PhoneEnsembleSize)
	for i := 0; i < len(bank) && len(accepted) < experimentalV4PhoneEnsembleSize; i++ {
		q := bank[i]
		if q.proposal < experimentalV4PhoneProposalFloor || q.validation < experimentalV4PhoneValidationFloor {
			continue
		}
		if i == 0 {
			q.detection = win.detection
		}
		accepted = append(accepted, q)
	}
	return work, boundary, accepted, evals, down
}

func experimentalV4PhoneDecodeEnsemble(src image.Image, key []byte, cw, ch int, hyp []experimentalV4PhoneHypothesis) ([]byte, ExperimentalV4ExtractInfo, experimentalV4PhoneDecodeTelemetry, error) {
	candidate := experimentalV4Prototype2Candidate()
	info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate)}
	telemetry := experimentalV4PhoneDecodeTelemetry{}
	if len(hyp) < experimentalV4PhoneEnsembleSize {
		return nil, info, telemetry, errors.New("experimental v4 phone geometry ensemble not accepted")
	}
	telemetry.Attempted = true
	info.PilotScore = hyp[0].detection.Score
	info.PilotMargin = hyp[0].detection.Margin
	info.OriginXBlocks = hyp[0].detection.OriginXBlocks
	info.OriginYBlocks = hyp[0].detection.OriginYBlocks
	plane := newPixelPlane(src)
	for _, spec := range v3Profiles {
		var agg []float64
		ws, cs := 0.0, 0.0
		for _, q := range hyp[:experimentalV4PhoneEnsembleSize] {
			var m []float64
			var c float64
			var ok bool
			if q.warp != nil {
				m, c, ok = experimentalV4ReadProtectedPhoneWarpMargins(plane, cw, ch, q.h, q.warp, spec.codedBits)
			} else {
				m, c, ok = experimentalV4ReadProtectedProjectiveMargins(plane, cw, ch, q.h, spec.codedBits)
			}
			if !ok {
				continue
			}
			if agg == nil {
				agg = make([]float64, len(m))
			}
			w := math.Max(0.001, q.validation)
			for i := range m {
				agg[i] += w * m[i]
			}
			ws += w
			cs += w * c
		}
		if ws <= 0 || len(agg) == 0 {
			continue
		}
		telemetry.ProfilesTried++
		for i := range agg {
			agg[i] /= ws
		}
		confidence := cs / ws
		if confidence > telemetry.MaxConfidence {
			telemetry.MaxConfidence = confidence
		}
		soft := experimentalV4SoftHammingDecodeMargins(agg)
		raw := bitsToBytes(whiten(soft, key, experimentalV4WhitenLabel))
		if payload, err := parseExperimentalV4Frame(raw, key, spec); err == nil {
			info.Profile = spec.profile
			info.Confidence = confidence
			telemetry.Authenticated = true
			telemetry.Profile = spec.profile
			return payload, info, telemetry, nil
		}
	}
	return nil, info, telemetry, errors.New("experimental v4 phone payload authentication failed")
}

// ExperimentalV4ExtractPhone preserves the Build44-qualified smartphone path:
// Build41/43 geometry, Build42 post-geometry data recovery and Build40 residual
// warp remain first and unchanged. Build64 adds only a final bounded deep
// recovery candidate. Its geometry bank is generated from public proposal
// evidence and frozen before held-out qualification or HMAC are consulted.
func ExperimentalV4ExtractPhone(src image.Image, key []byte, cw, ch int) ([]byte, ExperimentalV4ExtractInfo, ExperimentalV4PhoneInfo, error) {
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
		// Build66 candidate uses the exact Build64 recovery bank and is allowed to run even when Build41/43 cannot form the
		// historical two-geometry production ensemble. Its complete deep bank is
		// frozen proposal-only before held-out qualification or HMAC are read.
		payload64, info64, recovery64, err64 := experimentalV4PhoneBuild66Recover(work, boundary, key, cw, ch)
		experimentalV4PhoneBuild66ApplyTelemetry(&public, recovery64)
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

	// Build66 ordered-parallel decode candidate for the qualified Build64 deep-recovery fallback. This remains deliberately last so all
	// Build44-qualified paths remain untouched. The complete deep geometry bank
	// is generated proposal-only and frozen before held-out qualification or
	// HMAC are consulted. HMAC remains final frame authentication only.
	payload64, info64, recovery64, err64 := experimentalV4PhoneBuild66Recover(work, boundary, key, cw, ch)
	experimentalV4PhoneBuild66ApplyTelemetry(&public, recovery64)
	if err64 == nil {
		public.ProjectiveBasinFound = true
		public.Accepted = true
		public.HMACAuthenticated = true
		return payload64, info64, public, nil
	}
	return nil, info, public, baseErr
}
