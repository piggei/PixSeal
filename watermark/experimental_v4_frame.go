package watermark

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"image"
	"math"
)

// Build31 introduces the first experimental Format-v4 framing/data path. It is
// intentionally isolated from the frozen v3 Embed/Extract API. The locked
// prototype-2 pilot is public geometry metadata; authenticated success still
// requires a valid v4 HMAC-protected frame.
const (
	experimentalV4VersionNibble byte = 0x40
	experimentalV4DataStride         = 251 // coprime with 448, 672 and 1120
	experimentalV4WhitenLabel        = "pixseal-whiten-v4"
	experimentalV4MACLabel           = "pixseal-frame-v4"
)

// ExperimentalV4EmbedInfo reports deterministic properties of one Build31 v4
// embed. These fields are research telemetry, not a normative wire contract.
type ExperimentalV4EmbedInfo struct {
	Version         int
	Profile         Profile
	PayloadBytes    int
	ProfileCapacity int
	ProtectedBits   int
	PilotPositions  int
	DataPositions   int
	PilotName       string
	PilotHash       string
}

// ExperimentalV4ExtractInfo reports the evidence used by the Build31 aligned
// decoder. Pilot score/margin are proposal telemetry only; HMAC authentication
// is the success criterion.
type ExperimentalV4ExtractInfo struct {
	Version       int
	Profile       Profile
	Confidence    float64
	PilotScore    float64
	PilotMargin   float64
	OriginXBlocks int
	OriginYBlocks int
	PilotName     string
	PilotHash     string
}

// ExperimentalV4ProjectiveInfo reports public, non-authenticating geometry
// telemetry from the Build34 blind projective decoder. Payload/header/CRC/key
// and HMAC evidence are not used to select these parameters.
type ExperimentalV4ProjectiveInfo struct {
	Accepted            bool
	AngleDegrees        float64
	ScaleX              float64
	ScaleY              float64
	TopInset            float64
	BottomInset         float64
	PlacementShiftX     float64
	PlacementShiftY     float64
	ProposalObjective   float64
	GeometryValidation  float64
	PlacementValidation float64
	StructuralScore     float64
	HypothesesEvaluated int
}

// ExperimentalV4Capacity returns the current Build31 payload ceiling when a
// complete 37x32 v4 tile fits in the image. It does not affect v3 Capacity.
func ExperimentalV4Capacity(img image.Image) int {
	if img == nil {
		return 0
	}
	bounds := img.Bounds()
	if bounds.Dx()/blockSize < experimentalV4TileWidthBlocks || bounds.Dy()/blockSize < experimentalV4TileHeightBlocks {
		return 0
	}
	return maxPayload
}

// ExperimentalV4EmbedWithInfo writes the locked prototype-2 public pilot and
// an authenticated v4 frame into the 1120 non-pilot positions. The profile
// sizes and Hamming(7,4) code are intentionally inherited for the first v4
// framing experiment so pilot/framing effects can be isolated from a future
// ECC redesign.
func ExperimentalV4EmbedWithInfo(src image.Image, payload, key []byte, options Options) (*image.NRGBA, ExperimentalV4EmbedInfo, error) {
	if err := validateWorkingImageSize(src); err != nil {
		return nil, ExperimentalV4EmbedInfo{}, err
	}
	options, err := normalizeEmbedOptions(options)
	if err != nil {
		return nil, ExperimentalV4EmbedInfo{}, err
	}
	if len(key) < 8 {
		return nil, ExperimentalV4EmbedInfo{}, errors.New("key must contain at least 8 bytes")
	}
	spec, err := experimentalV4SelectProfile(len(payload), options.Profile)
	if err != nil {
		return nil, ExperimentalV4EmbedInfo{}, err
	}
	if ExperimentalV4Capacity(src) == 0 {
		return nil, ExperimentalV4EmbedInfo{}, errors.New("image must be at least 296x256 pixels for an experimental v4 payload")
	}

	frame := makeExperimentalV4Frame(payload, key, spec)
	protected := hammingEncode(whiten(bytesToBits(frame), key, experimentalV4WhitenLabel))
	if len(protected) != spec.codedBits {
		return nil, ExperimentalV4EmbedInfo{}, errors.New("internal experimental v4 protected-frame size mismatch")
	}

	candidate := experimentalV4Prototype2Candidate()
	var pilotSigns [experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks]int8
	var isPilot [experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks]bool
	for i, position := range candidate.positions {
		pilotSigns[position] = candidate.signs[i]
		isPilot[position] = true
	}
	dataOrdinal := experimentalV4DataOrdinalByTilePosition()

	out := toNRGBA(src)
	bounds := out.Bounds()
	blocksWide := bounds.Dx() / blockSize
	blocksHigh := bounds.Dy() / blockSize
	for blockY := 0; blockY < blocksHigh; blockY++ {
		for blockX := 0; blockX < blocksWide; blockX++ {
			tilePosition := (blockY%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + blockX%experimentalV4TileWidthBlocks
			bit := byte(0)
			if isPilot[tilePosition] {
				if pilotSigns[tilePosition] > 0 {
					bit = 1
				}
			} else {
				ordinal := dataOrdinal[tilePosition]
				if ordinal < 0 {
					return nil, ExperimentalV4EmbedInfo{}, errors.New("internal experimental v4 data-position mapping failure")
				}
				bit = protected[experimentalV4CodeIndex(ordinal, spec.codedBits)]
			}
			embedBlock(out, point{blockX * blockSize, blockY * blockSize}, bit, options.Strength)
		}
	}

	return out, ExperimentalV4EmbedInfo{
		Version:         experimentalV4Version,
		Profile:         spec.profile,
		PayloadBytes:    len(payload),
		ProfileCapacity: spec.maxPayload,
		ProtectedBits:   spec.codedBits,
		PilotPositions:  experimentalV4PilotCount,
		DataPositions:   experimentalV4DataCount,
		PilotName:       candidate.name,
		PilotHash:       experimentalV4PilotCandidateHash(candidate),
	}, nil
}

// ExperimentalV4ExtractAligned authenticates a Build31 v4 carrier when its
// native 8-pixel lattice is already aligned with the image axes. The public
// pilot resolves cyclic tile origin before any profile/frame decode attempt.
// Rotation/scale/projective recovery is deliberately outside this Build31 API.
func ExperimentalV4ExtractAligned(src image.Image, key []byte) ([]byte, ExperimentalV4ExtractInfo, error) {
	if err := validateWorkingImageSize(src); err != nil {
		return nil, ExperimentalV4ExtractInfo{}, err
	}
	if len(key) < 8 {
		return nil, ExperimentalV4ExtractInfo{}, errors.New("key must contain at least 8 bytes")
	}
	if ExperimentalV4Capacity(src) == 0 {
		return nil, ExperimentalV4ExtractInfo{}, errors.New("image must contain at least one complete 37x32 experimental v4 tile")
	}

	candidate := experimentalV4Prototype2Candidate()
	pilot := experimentalV4DetectPilotGrid(src, candidate, blockSize, 0, 0)
	if !pilot.Available {
		return nil, ExperimentalV4ExtractInfo{}, errors.New("experimental v4 pilot unavailable on aligned lattice")
	}

	plane := newPixelPlane(src)
	for _, spec := range v3Profiles {
		coded, confidence, ok := experimentalV4ReadProtectedAligned(plane, pilot.OriginXBlocks, pilot.OriginYBlocks, spec.codedBits)
		if !ok {
			continue
		}
		decoded := hammingDecode(coded)
		raw := bitsToBytes(whiten(decoded, key, experimentalV4WhitenLabel))
		payload, err := parseExperimentalV4Frame(raw, key, spec)
		if err != nil {
			continue
		}
		return payload, ExperimentalV4ExtractInfo{
			Version:       experimentalV4Version,
			Profile:       spec.profile,
			Confidence:    confidence,
			PilotScore:    pilot.Score,
			PilotMargin:   pilot.Margin,
			OriginXBlocks: pilot.OriginXBlocks,
			OriginYBlocks: pilot.OriginYBlocks,
			PilotName:     candidate.name,
			PilotHash:     experimentalV4PilotCandidateHash(candidate),
		}, nil
	}
	return nil, ExperimentalV4ExtractInfo{
		Version:       experimentalV4Version,
		PilotScore:    pilot.Score,
		PilotMargin:   pilot.Margin,
		OriginXBlocks: pilot.OriginXBlocks,
		OriginYBlocks: pilot.OriginYBlocks,
		PilotName:     candidate.name,
		PilotHash:     experimentalV4PilotCandidateHash(candidate),
	}, errors.New("experimental v4 aligned payload authentication failed")
}

func experimentalV4SelectProfile(payloadBytes int, requested Profile) (profileSpec, error) {
	if payloadBytes < 0 {
		return profileSpec{}, errorsPayloadSize(payloadBytes)
	}
	minimum, ok := minimumProfile(payloadBytes)
	if !ok {
		return profileSpec{}, fmt.Errorf("payload is %d bytes; experimental PixSeal v4 maximum is %d bytes", payloadBytes, maxPayload)
	}
	if requested == "" {
		requested = ProfileAuto
	}
	if requested == ProfileAuto {
		return minimum, nil
	}
	spec, ok := profileSpecFor(requested)
	if !ok {
		return profileSpec{}, fmt.Errorf("invalid profile %q; expected auto, robust, balanced or capacity", requested)
	}
	if payloadBytes > spec.maxPayload {
		return profileSpec{}, fmt.Errorf(
			"payload is %d bytes; experimental v4 profile %s capacity is %d bytes; minimum compatible profile is %s",
			payloadBytes, requested, spec.maxPayload, minimum.profile,
		)
	}
	return spec, nil
}

func experimentalV4HeaderByte(spec profileSpec) byte {
	return experimentalV4VersionNibble | spec.id
}

func makeExperimentalV4Frame(payload, key []byte, spec profileSpec) []byte {
	frame := make([]byte, spec.frameBytes)
	frame[0], frame[1] = magic[0], magic[1]
	frame[2] = experimentalV4HeaderByte(spec)
	frame[3] = byte(len(payload))
	binary.BigEndian.PutUint32(frame[4:8], crc32.ChecksumIEEE(payload))
	copy(frame[headerSize:], payload)
	tagOffset := headerSize + len(payload)
	mac := experimentalV4FrameMAC(key, frame[:tagOffset])
	copy(frame[tagOffset:], mac[:tagSize])
	return frame
}

func parseExperimentalV4Frame(frame, key []byte, spec profileSpec) ([]byte, error) {
	if len(frame) < spec.frameBytes || frame[0] != magic[0] || frame[1] != magic[1] || frame[2] != experimentalV4HeaderByte(spec) {
		return nil, errors.New("experimental v4 hidden payload header mismatch")
	}
	payloadLength := int(frame[3])
	if payloadLength > spec.maxPayload {
		return nil, errors.New("invalid experimental v4 hidden payload length")
	}
	tagOffset := headerSize + payloadLength
	if tagOffset+tagSize > len(frame) {
		return nil, errors.New("invalid experimental v4 hidden payload frame")
	}
	mac := experimentalV4FrameMAC(key, frame[:tagOffset])
	if !hmac.Equal(frame[tagOffset:tagOffset+tagSize], mac[:tagSize]) {
		return nil, errors.New("experimental v4 hidden payload authentication failed")
	}
	payload := frame[headerSize:tagOffset]
	if crc32.ChecksumIEEE(payload) != binary.BigEndian.Uint32(frame[4:8]) {
		return nil, errors.New("experimental v4 hidden payload damaged: CRC mismatch")
	}
	return append([]byte(nil), payload...), nil
}

func experimentalV4FrameMAC(key, authenticatedPrefix []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(experimentalV4MACLabel))
	mac.Write([]byte{0})
	mac.Write(authenticatedPrefix)
	return mac.Sum(nil)
}

func experimentalV4CodeIndex(dataOrdinal, codedBits int) int {
	return (dataOrdinal * experimentalV4DataStride) % codedBits
}

func experimentalV4LockedDataPositions() []int {
	pilot := make(map[int]struct{}, experimentalV4PilotCount)
	for _, position := range experimentalV4Prototype2PilotPositions {
		pilot[position] = struct{}{}
	}
	positions := make([]int, 0, experimentalV4DataCount)
	for position := 0; position < experimentalV4TileWidthBlocks*experimentalV4TileHeightBlocks; position++ {
		if _, reserved := pilot[position]; !reserved {
			positions = append(positions, position)
		}
	}
	return positions
}

func experimentalV4DataOrdinalByTilePosition() [experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks]int {
	var mapping [experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks]int
	for i := range mapping {
		mapping[i] = -1
	}
	for ordinal, position := range experimentalV4LockedDataPositions() {
		mapping[position] = ordinal
	}
	return mapping
}

func experimentalV4ReadProtectedAligned(plane *pixelPlane, originXBlocks, originYBlocks, codedBits int) ([]byte, float64, bool) {
	if plane == nil || codedBits <= 0 {
		return nil, 0, false
	}
	bounds := plane.bounds
	blocksWide := bounds.Dx() / blockSize
	blocksHigh := bounds.Dy() / blockSize
	if blocksWide < experimentalV4TileWidthBlocks || blocksHigh < experimentalV4TileHeightBlocks {
		return nil, 0, false
	}

	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var residueSums [residues]float64
	var residueCounts [residues]int
	for blockY := 0; blockY < blocksHigh; blockY++ {
		for blockX := 0; blockX < blocksWide; blockX++ {
			residue := (blockY%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + blockX%experimentalV4TileWidthBlocks
			residueSums[residue] += readBlockSized(plane, point{x: bounds.Min.X + blockX*blockSize, y: bounds.Min.Y + blockY*blockSize}, blockSize)
			residueCounts[residue]++
		}
	}

	bitSums := make([]float64, codedBits)
	bitObservations := make([]int, codedBits)
	for ordinal, logicalPosition := range experimentalV4LockedDataPositions() {
		logicalX := logicalPosition % experimentalV4TileWidthBlocks
		logicalY := logicalPosition / experimentalV4TileWidthBlocks
		observedX := positiveMod(logicalX-originXBlocks, experimentalV4TileWidthBlocks)
		observedY := positiveMod(logicalY-originYBlocks, experimentalV4TileHeightBlocks)
		residue := observedY*experimentalV4TileWidthBlocks + observedX
		if residueCounts[residue] == 0 {
			continue
		}
		codeIndex := experimentalV4CodeIndex(ordinal, codedBits)
		bitSums[codeIndex] += residueSums[residue]
		bitObservations[codeIndex] += residueCounts[residue]
	}

	coded := make([]byte, codedBits)
	confidence := 0.0
	observedBits := 0
	for i, sum := range bitSums {
		if bitObservations[i] == 0 {
			continue
		}
		observedBits++
		if sum >= 0 {
			coded[i] = 1
		}
		confidence += math.Abs(sum) / float64(bitObservations[i])
	}
	if observedBits != codedBits {
		return nil, 0, false
	}
	return coded, confidence / float64(codedBits), true
}

// experimentalV4ReadProtectedProjective samples the locked v4 data partition
// through an already accepted canonical->observed homography. Geometry has
// already been selected exclusively from public structure/pilot evidence; no
// frame/header/key/HMAC data participates in that selection.
func experimentalV4ReadProtectedProjectiveMargins(plane *pixelPlane, canonicalWidth, canonicalHeight int, h homography, codedBits int) ([]float64, float64, bool) {
	if plane == nil || codedBits <= 0 {
		return nil, 0, false
	}
	blocksWide := canonicalWidth / blockSize
	blocksHigh := canonicalHeight / blockSize
	if blocksWide < experimentalV4TileWidthBlocks || blocksHigh < experimentalV4TileHeightBlocks {
		return nil, 0, false
	}
	const residues = experimentalV4TileWidthBlocks * experimentalV4TileHeightBlocks
	var residueSums [residues]float64
	var residueCounts [residues]int
	for by := 0; by < blocksHigh; by++ {
		for bx := 0; bx < blocksWide; bx++ {
			v, ok := readProjectiveBlockValue(plane, h, bx*blockSize, by*blockSize, blockSize)
			if !ok {
				continue
			}
			residue := (by%experimentalV4TileHeightBlocks)*experimentalV4TileWidthBlocks + bx%experimentalV4TileWidthBlocks
			residueSums[residue] += v
			residueCounts[residue]++
		}
	}
	bitSums := make([]float64, codedBits)
	bitObservations := make([]int, codedBits)
	for ordinal, pos := range experimentalV4LockedDataPositions() {
		if residueCounts[pos] == 0 {
			continue
		}
		idx := experimentalV4CodeIndex(ordinal, codedBits)
		bitSums[idx] += residueSums[pos]
		bitObservations[idx] += residueCounts[pos]
	}
	margins := make([]float64, codedBits)
	confidence := 0.0
	observedBits := 0
	for i, sum := range bitSums {
		if bitObservations[i] == 0 {
			continue
		}
		observedBits++
		margins[i] = sum / float64(bitObservations[i])
		confidence += math.Abs(margins[i])
	}
	if observedBits != codedBits {
		return nil, 0, false
	}
	return margins, confidence / float64(codedBits), true
}

func experimentalV4ReadProtectedProjective(plane *pixelPlane, canonicalWidth, canonicalHeight int, h homography, codedBits int) ([]byte, float64, bool) {
	margins, confidence, ok := experimentalV4ReadProtectedProjectiveMargins(plane, canonicalWidth, canonicalHeight, h, codedBits)
	if !ok {
		return nil, 0, false
	}
	coded := make([]byte, codedBits)
	for i, margin := range margins {
		if margin >= 0 {
			coded[i] = 1
		}
	}
	return coded, confidence, true
}

// experimentalV4SoftHammingDecodeMargins performs deterministic maximum-
// likelihood decoding of each Hamming(7,4) word from signed DCT evidence.
// It is used only after blind geometry has already passed the public pilot
// gate; it never participates in geometry proposal or acceptance.
func experimentalV4SoftHammingDecodeMargins(margins []float64) []byte {
	words := len(margins) / 7
	out := make([]byte, words*4)
	for word := 0; word < words; word++ {
		base := word * 7
		bestScore := math.Inf(-1)
		var best [4]byte
		for value := 0; value < 16; value++ {
			nibble := []byte{byte((value >> 3) & 1), byte((value >> 2) & 1), byte((value >> 1) & 1), byte(value & 1)}
			encoded := hammingEncode(nibble)
			score := 0.0
			for bit := 0; bit < 7; bit++ {
				sign := -1.0
				if encoded[bit] != 0 {
					sign = 1
				}
				score += sign * margins[base+bit]
			}
			if score > bestScore {
				bestScore = score
				copy(best[:], nibble)
			}
		}
		copy(out[word*4:word*4+4], best[:])
	}
	return out
}

// experimentalV4ExtractJointProjective performs blind Build34 geometry
// recovery first and only then attempts authenticated v4 frame decoding.
func experimentalV4ExtractJointProjective(src image.Image, key []byte, canonicalWidth, canonicalHeight int) ([]byte, ExperimentalV4ExtractInfo, experimentalV4JointProjectiveResult, error) {
	if src == nil {
		return nil, ExperimentalV4ExtractInfo{}, experimentalV4JointProjectiveResult{}, errors.New("nil image")
	}
	if len(key) < 8 {
		return nil, ExperimentalV4ExtractInfo{}, experimentalV4JointProjectiveResult{}, errors.New("key must contain at least 8 bytes")
	}
	candidate := experimentalV4Prototype2Candidate()
	geometry := experimentalV4JointProjectiveSearch(src, candidate, canonicalWidth, canonicalHeight)
	info := ExperimentalV4ExtractInfo{Version: experimentalV4Version, PilotScore: geometry.detection.Score, PilotMargin: geometry.detection.Margin, OriginXBlocks: geometry.detection.OriginXBlocks, OriginYBlocks: geometry.detection.OriginYBlocks, PilotName: candidate.name, PilotHash: experimentalV4PilotCandidateHash(candidate)}
	if !experimentalV4Build29ProjectiveAccepted(geometry) {
		return nil, info, geometry, errors.New("experimental v4 projective geometry not accepted")
	}
	plane := newPixelPlane(src)
	for _, spec := range v3Profiles {
		margins, confidence, ok := experimentalV4ReadProtectedProjectiveMargins(plane, canonicalWidth, canonicalHeight, geometry.placement.h, spec.codedBits)
		if !ok {
			continue
		}

		// Soft ML decoding is especially useful after print/scan/JPEG, where a
		// Hamming word can contain more than one weak sign error. HMAC remains
		// the only success criterion and geometry is already frozen above.
		softDecoded := experimentalV4SoftHammingDecodeMargins(margins)
		softRaw := bitsToBytes(whiten(softDecoded, key, experimentalV4WhitenLabel))
		if payload, err := parseExperimentalV4Frame(softRaw, key, spec); err == nil {
			info.Profile = spec.profile
			info.Confidence = confidence
			return payload, info, geometry, nil
		}

		// Preserve the Build31/34 hard-decision path as a deterministic fallback.
		coded := make([]byte, len(margins))
		for i, margin := range margins {
			if margin >= 0 {
				coded[i] = 1
			}
		}
		hardDecoded := hammingDecode(coded)
		hardRaw := bitsToBytes(whiten(hardDecoded, key, experimentalV4WhitenLabel))
		if payload, err := parseExperimentalV4Frame(hardRaw, key, spec); err == nil {
			info.Profile = spec.profile
			info.Confidence = confidence
			return payload, info, geometry, nil
		}
	}
	return nil, info, geometry, errors.New("experimental v4 projective payload authentication failed")
}

// ExperimentalV4ExtractProjective exposes the Build34 blind projective+crop
// path for controlled physical-channel qualification. canonicalWidth and
// canonicalHeight are the pixel dimensions of the digital carrier before
// print/acquisition. Geometry is selected only from public structure/pilot
// evidence; frame authentication happens strictly after geometry acceptance.
func ExperimentalV4ExtractProjective(src image.Image, key []byte, canonicalWidth, canonicalHeight int) ([]byte, ExperimentalV4ExtractInfo, ExperimentalV4ProjectiveInfo, error) {
	if canonicalWidth < experimentalV4TileWidthBlocks*blockSize || canonicalHeight < experimentalV4TileHeightBlocks*blockSize || canonicalWidth%blockSize != 0 || canonicalHeight%blockSize != 0 {
		return nil, ExperimentalV4ExtractInfo{}, ExperimentalV4ProjectiveInfo{}, errors.New("canonical dimensions must be block-aligned and contain at least one complete v4 tile")
	}
	payload, info, search, err := experimentalV4ExtractJointProjective(src, key, canonicalWidth, canonicalHeight)
	projective := ExperimentalV4ProjectiveInfo{
		Accepted:            experimentalV4Build29ProjectiveAccepted(search),
		AngleDegrees:        search.params.angleDeg,
		ScaleX:              search.params.scaleX,
		ScaleY:              search.params.scaleY,
		TopInset:            search.params.topInset,
		BottomInset:         search.params.bottomInset,
		PlacementShiftX:     search.placement.shiftX,
		PlacementShiftY:     search.placement.shiftY,
		ProposalObjective:   search.proposalObjective,
		GeometryValidation:  search.geometryValidation,
		PlacementValidation: search.placement.validationScore,
		StructuralScore:     search.structuralScore,
		HypothesesEvaluated: search.hypothesesEvaluated,
	}
	return payload, info, projective, err
}
