package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"hash/crc32"
	"image"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pj/pixseal/internal/buildinfo"
	"github.com/pj/pixseal/watermark"
)

func TestSubcommandHelpReturnsFlagErrHelp(t *testing.T) {
	for name, fn := range map[string]func([]string) error{
		"embed": embed, "extract": extract, "v4-embed": v4Embed, "v4-extract": v4Extract, "v4-extract-projective": v4ExtractProjective, "v4-extract-scanner": v4ExtractScanner, "v4-extract-phone": v4ExtractPhone, "v4-diagnose-phone": v4DiagnosePhone, "v4-diagnose-phone-handoff": v4DiagnosePhoneHandoff, "capacity": capacity, "analyze": analyze, "diagnose": diagnose,
	} {
		t.Run(name, func(t *testing.T) {
			if err := fn([]string{"-help"}); !errors.Is(err, flag.ErrHelp) {
				t.Fatalf("%s -help error = %v; want flag.ErrHelp", name, err)
			}
		})
	}
}

func TestV4ExtractProjectiveRequiresCanonicalBlockDimensions(t *testing.T) {
	for _, tc := range []struct {
		w, h string
	}{
		{"1635", "1632"},
		{"1632", "1635"},
		{"288", "1632"},
		{"1632", "248"},
	} {
		err := v4ExtractProjective([]string{"-in", "does-not-matter.png", "-key", "12345678", "-width", tc.w, "-height", tc.h})
		if err == nil || !strings.Contains(err.Error(), "divisible by 8") {
			t.Fatalf("dimensions %sx%s error=%v", tc.w, tc.h, err)
		}
	}
}

func TestV4ExtractScannerRequiresCanonicalBlockDimensions(t *testing.T) {
	for _, tc := range []struct {
		w, h string
	}{
		{"1635", "1632"},
		{"1632", "1635"},
		{"288", "1632"},
		{"1632", "248"},
	} {
		err := v4ExtractScanner([]string{"-in", "does-not-matter.jpg", "-key", "12345678", "-width", tc.w, "-height", tc.h})
		if err == nil || !strings.Contains(err.Error(), "divisible by 8") {
			t.Fatalf("dimensions %sx%s error=%v", tc.w, tc.h, err)
		}
	}
}

func TestV4ExtractPhoneRequiresCanonicalBlockDimensions(t *testing.T) {
	for _, tc := range []struct{ w, h string }{{"1635", "1632"}, {"1632", "1635"}, {"288", "1632"}, {"1632", "248"}} {
		err := v4ExtractPhone([]string{"-in", "does-not-matter.jpg", "-key", "12345678", "-width", tc.w, "-height", tc.h})
		if err == nil || !strings.Contains(err.Error(), "divisible by 8") {
			t.Fatalf("dimensions %sx%s error=%v", tc.w, tc.h, err)
		}
	}
}

func TestV4PhoneDiagnosticsExposeBuild42MatrixFields(t *testing.T) {
	var buf bytes.Buffer
	info := watermark.ExperimentalV4ExtractInfo{Confidence: 12.5, Profile: watermark.ProfileRobust, PilotScore: 0.3, PilotMargin: 0.08, OriginXBlocks: 0, OriginYBlocks: 0, PilotName: "pilot", PilotHash: "hash"}
	phone := watermark.ExperimentalV4PhoneInfo{
		WorkingWidth: 1200, WorkingHeight: 900, Downsampled: true,
		BoundaryDetected: true, BoundaryConfidence: 0.8, ProjectiveBasinFound: true,
		Accepted: true, EnsembleCandidates: 2, ProposalScore: 0.25, ValidationScore: 0.2,
		ResidualAttempted: true, ResidualFitted: true, ResidualApplied: false, ResidualControls: 9, ResidualRMSPixels: 1.2,
		ResidualProposalBefore: 0.2, ResidualProposalAfter: 0.3, ResidualValidationBefore: 0.18, ResidualValidationAfter: 0.16,
		DataDecodeAttempted: true, SoftHammingProfiles: 3, MaxDataConfidence: 11.2,
		HMACAuthenticated: false, FallbackAttempted: true, FallbackAuthenticated: false,
		Build42DataAttempted: true, Build42BankCandidates: 5, Build42EnsemblesTried: 2, Build42ListFramesTried: 1030, Build42DataAuthenticated: true,
	}
	printV4PhoneDiagnostics(&buf, "pixseal-jpeg-pre-go1.26-v1", info, phone)
	out := buf.String()
	for _, want := range []string{
		"input-decoder: pixseal-jpeg-pre-go1.26-v1",
		"projective-basin: found=true",
		"phone-residual: fitted=true applied=false",
		"data-decode: attempted=true soft-hamming-profiles=3 max-confidence=11.20",
		"build42-data: attempted=true bank=5 ensembles=2 list-frames=1030 authenticated=true",
		"hmac: authenticated=false fallback-attempted=true fallback-authenticated=false",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostics missing %q:\n%s", want, out)
		}
	}
}

func TestOpenImageUsesDeterministicLegacyJPEGDecoder(t *testing.T) {
	path := filepath.Join("..", "..", "internal", "jpeglegacy", "testdata", "legacy-sample.jpg")
	img, format, err := openImageWithFormat(path)
	if err != nil {
		t.Fatal(err)
	}
	if format != "jpeg" {
		t.Fatalf("format=%q", format)
	}
	y, ok := img.(*image.YCbCr)
	if !ok {
		t.Fatalf("decoded type %T, want *image.YCbCr", img)
	}
	sum := sha256.Sum256(y.Y)
	if got := hex.EncodeToString(sum[:]); got != "eb6ec297f45a764be2c668febd833ad55dda2c75313589fd7d3909d188543147" {
		t.Fatalf("Y sha256=%s", got)
	}
	if got := inputDecoderID(format); got != "pixseal-jpeg-pre-go1.26-v1" {
		t.Fatalf("decoder=%q", got)
	}
}

func TestPNGOutputPath(t *testing.T) {
	tests := map[string]string{
		"marked.png":  "marked.png",
		"marked.PNG":  "marked.PNG",
		"marked.jpg":  "marked.png",
		"marked.jpeg": "marked.png",
		"marked":      "marked.png",
		".sealed":     ".sealed.png",
	}
	for input, want := range tests {
		if got := pngOutputPath(input); got != want {
			t.Errorf("pngOutputPath(%q) = %q; want %q", input, got, want)
		}
	}
}

func TestOpenImageRejectsOversizedPNGFromConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "huge.png")
	var ihdr [13]byte
	binary.BigEndian.PutUint32(ihdr[0:4], 20_000)
	binary.BigEndian.PutUint32(ihdr[4:8], 20_000) // 400 MP > 300 MP policy
	ihdr[8], ihdr[9], ihdr[10], ihdr[11], ihdr[12] = 8, 2, 0, 0, 0
	data := append([]byte("\x89PNG\r\n\x1a\n"), 0, 0, 0, 13)
	data = append(data, []byte("IHDR")...)
	data = append(data, ihdr[:]...)
	crc := crc32.ChecksumIEEE(append([]byte("IHDR"), ihdr[:]...))
	var crcBytes [4]byte
	binary.BigEndian.PutUint32(crcBytes[:], crc)
	data = append(data, crcBytes[:]...)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := openImageWithFormat(path); err == nil || !strings.Contains(err.Error(), "CLI safety limit") {
		t.Fatalf("oversized DecodeConfig guard error=%v", err)
	}
}

func TestCapacityUsesDecodeConfigOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config-only.png")
	var ihdr [13]byte
	binary.BigEndian.PutUint32(ihdr[0:4], 560)
	binary.BigEndian.PutUint32(ihdr[4:8], 512)
	ihdr[8], ihdr[9], ihdr[10], ihdr[11], ihdr[12] = 8, 2, 0, 0, 0
	data := append([]byte("\x89PNG\r\n\x1a\n"), 0, 0, 0, 13)
	data = append(data, []byte("IHDR")...)
	data = append(data, ihdr[:]...)
	crc := crc32.ChecksumIEEE(append([]byte("IHDR"), ihdr[:]...))
	var crcBytes [4]byte
	binary.BigEndian.PutUint32(crcBytes[:], crc)
	data = append(data, crcBytes[:]...)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	got := captureStdout(t, func() error { return capacity([]string{"-in", path}) })
	if got != "64 bytes\n" {
		t.Fatalf("capacity config-only output = %q; want %q", got, "64 bytes\n")
	}
}

func TestEmbedRequiresMessage(t *testing.T) {
	err := embed([]string{"-in", "input.png", "-out", "output.png", "-key", "12345678"})
	if err == nil {
		t.Fatal("embed accepted a missing -message")
	}
}

func TestEmbedRejectsPositionalMessage(t *testing.T) {
	err := embed([]string{"-in", "input.png", "-out", "output.png", "-key", "12345678", "message"})
	if err == nil {
		t.Fatal("embed accepted an unlabelled positional message")
	}
}

func TestEmbedRejectsInvalidStrengthAtCLI(t *testing.T) {
	for _, value := range []string{"0", "NaN", "+Inf", "-Inf", "3.9", "120.1"} {
		t.Run(value, func(t *testing.T) {
			err := embed([]string{
				"-in", "does-not-matter.png", "-out", "out.png", "-key", "12345678",
				"-message", "hello", "-strength", value,
			})
			if err == nil || !strings.Contains(err.Error(), "-strength") {
				t.Fatalf("strength %q error = %v", value, err)
			}
		})
	}
}

func TestWritePNGAtomicRejectsDirectoryAndSymlink(t *testing.T) {
	dir := t.TempDir()
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))

	directoryTarget := filepath.Join(dir, "target.png")
	if err := os.Mkdir(directoryTarget, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := writePNGAtomic(directoryTarget, img, true); err == nil {
		t.Fatal("writePNGAtomic replaced a directory with -force")
	}
	if info, err := os.Lstat(directoryTarget); err != nil || !info.IsDir() {
		t.Fatalf("directory target changed after rejected write: info=%v err=%v", info, err)
	}

	symlinkTarget := filepath.Join(dir, "link.png")
	if err := os.Symlink("missing-target", symlinkTarget); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if err := writePNGAtomic(symlinkTarget, img, false); err == nil {
		t.Fatal("writePNGAtomic replaced a dangling symlink without -force")
	}
	if info, err := os.Lstat(symlinkTarget); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("symlink target changed after rejected write: info=%v err=%v", info, err)
	}
}

func TestCommitNoClobberRejectsLateTarget(t *testing.T) {
	dir := t.TempDir()
	temporary := filepath.Join(dir, "temporary")
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(temporary, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := commitNoClobber(temporary, target); err == nil {
		t.Fatal("commitNoClobber overwrote a late-created target")
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "existing" {
		t.Fatalf("late-created target changed: %q", got)
	}
}

func TestExtractRawPreservesMultilinePayload(t *testing.T) {
	dir := t.TempDir()
	carrier := filepath.Join(dir, "carrier.png")
	sealed := filepath.Join(dir, "sealed.png")
	img := image.NewNRGBA(image.Rect(0, 0, 320, 288))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(x), G: uint8(y), B: uint8(x + y), A: 255})
		}
	}
	if err := writePNGAtomic(carrier, img, true); err != nil {
		t.Fatal(err)
	}
	message := "line1\nline2"
	_ = captureStdout(t, func() error {
		return embed([]string{"-in", carrier, "-out", sealed, "-key", "12345678", "-message", message, "-profile", "robust"})
	})
	got := captureStdout(t, func() error {
		return extract([]string{"-in", sealed, "-key", "12345678", "-raw"})
	})
	if got != message {
		t.Fatalf("raw extracted payload = %q; want %q", got, message)
	}
}

func TestVersionFileMatchesBuildInfo(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "VERSION"))
	if err != nil {
		t.Fatal(err)
	}
	want := "PixSeal " + buildinfo.String()
	if got := strings.TrimSpace(string(data)); got != want {
		t.Fatalf("VERSION = %q; want %q", got, want)
	}
}

func TestWritePNGAtomicProtectsExistingOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "carrier.png")
	original := []byte("do not overwrite")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.SetNRGBA(0, 0, color.NRGBA{R: 255, A: 255})

	if err := writePNGAtomic(path, img, false); err == nil {
		t.Fatal("writePNGAtomic overwrote an existing file without -force")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, original) {
		t.Fatal("existing output changed after rejected write")
	}

	if err := writePNGAtomic(path, img, true); err != nil {
		t.Fatalf("forced replacement failed: %v", err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, _, err := image.Decode(f); err != nil {
		t.Fatalf("forced output is not a valid image: %v", err)
	}
	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("forced replacement permissions=%v err=%v; want 0600", info, err)
	}
}

func TestAnalyzeRejectsInvalidPayloadOptions(t *testing.T) {
	if err := analyze([]string{"-in", "carrier.png"}); err == nil {
		t.Fatal("analyze accepted neither -message nor -bytes")
	}
	if err := analyze([]string{"-in", "carrier.png", "-message", "hello", "-bytes", "5"}); err == nil {
		t.Fatal("analyze accepted both -message and -bytes")
	}
	if err := analyze([]string{"-in", "carrier.png", "-bytes", "0"}); err == nil {
		t.Fatal("analyze accepted zero bytes")
	}
}

func TestEmbedRejectsInvalidProfile(t *testing.T) {
	err := embed([]string{
		"-in", "input.png", "-out", "output.png", "-key", "12345678",
		"-message", "hello", "-profile", "unknown",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid profile") {
		t.Fatalf("embed invalid profile error = %v", err)
	}
}

func TestRemovedCompatibilityFlagsAreRejected(t *testing.T) {
	if err := extract([]string{"-in", "carrier.png", "-key", "12345678", "-repetition", "5"}); err == nil {
		t.Fatal("extract still accepts removed -repetition flag")
	}
	if err := extract([]string{"-in", "carrier.png", "-key", "12345678", "-strength", "24"}); err == nil {
		t.Fatal("extract still accepts removed -strength flag")
	}
	if err := capacity([]string{"-in", "carrier.png", "-repetition", "5"}); err == nil {
		t.Fatal("capacity still accepts removed -repetition flag")
	}
}

func TestAnalyzeCapacityAndEmbedAgreeOnProfiles(t *testing.T) {
	dir := t.TempDir()
	carrier := filepath.Join(dir, "carrier.png")
	output := filepath.Join(dir, "sealed.png")
	img := image.NewNRGBA(image.Rect(0, 0, 320, 288))
	for y := 0; y < 288; y++ {
		for x := 0; x < 320; x++ {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8((x + y) % 256),
				G: uint8((2*x + y) % 256),
				B: uint8((x + 2*y) % 256),
				A: 255,
			})
		}
	}
	if err := writePNGAtomic(carrier, img, true); err != nil {
		t.Fatal(err)
	}

	capacityOutput := captureStdout(t, func() error {
		return capacity([]string{"-in", carrier, "-details"})
	})
	for _, want := range []string{"robust:      16 bytes", "balanced:    32 bytes", "capacity:    64 bytes", "maximum:     64 bytes"} {
		if !strings.Contains(capacityOutput, want) {
			t.Fatalf("capacity output missing %q:\n%s", want, capacityOutput)
		}
	}

	analyzeOutput := captureStdout(t, func() error {
		return analyze([]string{"-in", carrier, "-bytes", "18"})
	})
	if !strings.Contains(analyzeOutput, "Recommended profile:       balanced") ||
		!strings.Contains(analyzeOutput, "Profile capacity:          32 bytes") {
		t.Fatalf("analyze did not recommend balanced/32:\n%s", analyzeOutput)
	}

	embedOutput := captureStdout(t, func() error {
		return embed([]string{
			"-in", carrier, "-out", output, "-key", "12345678",
			"-message", strings.Repeat("x", 18),
		})
	})
	if !strings.Contains(embedOutput, "embedded 18 bytes using profile balanced") {
		t.Fatalf("embed did not auto-select balanced:\n%s", embedOutput)
	}
}

func TestExplicitProfileErrorNamesMinimumCompatibleProfile(t *testing.T) {
	dir := t.TempDir()
	carrier := filepath.Join(dir, "carrier.png")
	img := image.NewNRGBA(image.Rect(0, 0, 320, 288))
	if err := writePNGAtomic(carrier, img, true); err != nil {
		t.Fatal(err)
	}
	err := embed([]string{
		"-in", carrier, "-out", filepath.Join(dir, "out.png"), "-key", "12345678",
		"-message", strings.Repeat("x", 17), "-profile", "robust",
	})
	if err == nil || !strings.Contains(err.Error(), "payload is 17 bytes") ||
		!strings.Contains(err.Error(), "profile robust capacity is 16 bytes") ||
		!strings.Contains(err.Error(), "minimum compatible profile is balanced") {
		t.Fatalf("unexpected explicit-profile error: %v", err)
	}
}

func captureStdout(t *testing.T, fn func() error) string {
	t.Helper()
	old := os.Stdout
	readEnd, writeEnd, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writeEnd
	callErr := fn()
	_ = writeEnd.Close()
	os.Stdout = old
	data, readErr := io.ReadAll(readEnd)
	_ = readEnd.Close()
	if readErr != nil {
		t.Fatal(readErr)
	}
	if callErr != nil {
		t.Fatalf("captured command failed: %v", callErr)
	}
	return string(data)
}

func TestExperimentalV4CLIRoundTrip(t *testing.T) {
	dir := t.TempDir()
	carrier := filepath.Join(dir, "carrier.png")
	sealed := filepath.Join(dir, "sealed-v4.png")
	img := image.NewNRGBA(image.Rect(0, 0, 592, 512))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(x * 3), G: uint8(y * 5), B: uint8(x + y*2), A: 255})
		}
	}
	if err := writePNGAtomic(carrier, img, true); err != nil {
		t.Fatal(err)
	}
	message := "Build31 CLI v4"
	out := captureStdout(t, func() error {
		return v4Embed([]string{"-in", carrier, "-out", sealed, "-key", "Piccotti", "-message", message, "-profile", "robust"})
	})
	if !strings.Contains(out, "EXPERIMENTAL v4") || !strings.Contains(out, "prototype-2-search-p64") {
		t.Fatalf("unexpected v4-embed output: %q", out)
	}
	got := captureStdout(t, func() error {
		return v4Extract([]string{"-in", sealed, "-key", "Piccotti", "-raw"})
	})
	if got != message {
		t.Fatalf("v4 raw extracted payload=%q, want %q", got, message)
	}
	if _, err := os.Stat(sealed); err != nil {
		t.Fatalf("v4 carrier missing: %v", err)
	}
}
