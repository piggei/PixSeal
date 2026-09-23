package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	jpeglegacy "github.com/pj/pixseal/internal/jpeglegacy"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pj/pixseal/internal/buildinfo"
	"github.com/pj/pixseal/watermark"
)

const maxCLISourcePixels int64 = 300_000_000

var usageHeader = fmt.Sprintf(`PixSeal %s - robust image steganography

Usage:
  pixseal <command> [options]

Commands:
  embed      Hide an authenticated v3 message in an image
  extract    Recover and authenticate a v3 message with bounded geometric recovery
  v4-embed   EXPERIMENTAL: hide a Build31 Format-v4 message using the locked public pilot
  v4-extract EXPERIMENTAL: recover a Build31 v4 message on an aligned native 8px lattice
  v4-extract-projective EXPERIMENTAL: blind Build34 projective recovery + authenticated v4 decode
  v4-extract-scanner EXPERIMENTAL: Build37 blind paper/scanner registration + authenticated v4 decode
  v4-extract-phone EXPERIMENTAL: Build44 deterministic JPEG ingest + Build43 smartphone recovery + authenticated v4 decode
  v4-diagnose-phone EXPERIMENTAL: Build45 phone failure decomposition and optional lab-only supplied-geometry oracle
  v4-diagnose-phone-handoff EXPERIMENTAL: Build46 qualified-geometry handoff diagnostic
  v4-diagnose-phone-frozen EXPERIMENTAL: Build47 frozen-candidate bank observability diagnostic
  v4-diagnose-phone-refine EXPERIMENTAL: Build48 proposal-only local projective refinement diagnostic
  v4-diagnose-phone-ranking EXPERIMENTAL: Build49 proposal-ranking observability diagnostic
  v4-diagnose-phone-refine4 EXPERIMENTAL: Build50 top-4-per-pair local refinement diagnostic
  capacity   Show the usable payload capacity of an image
  analyze    Recommend a v3 profile and embedding settings
  diagnose   Experimental bounded local-lattice diagnostics (v0.3 research)

Supported image formats:
  Input       PNG (.png), JPEG (.jpg, .jpeg)
  Output      PNG (.png) only

Embed options:
  -in FILE           Input JPEG or PNG (required)
  -out FILE          Output PNG (required)
  -key TEXT          Secret key, minimum 8 bytes (required)
  -message TEXT      Message to hide (required)
  -profile NAME      auto, robust, balanced or capacity (default auto)
  -strength N        DCT embedding strength, 4 to 120 (default 24)
  -force             Allow replacing an existing output file

Extract options:
  -in FILE           Carrier JPEG or PNG (required)
  -key TEXT          Secret key, minimum 8 bytes (required)
  -raw               Write only authenticated payload bytes to stdout

Experimental v4 options:
  v4-embed accepts the same -in/-out/-key/-message/-profile/-strength/-force options as embed.
  v4-extract accepts -in/-key/-raw but requires an already aligned native 8px lattice.
  v4-extract-projective accepts -in/-key/-raw plus canonical -width/-height from the pre-print carrier.
  v4-extract-scanner accepts the same dimensions but expects a full-page scan with visible paper/artwork boundary.
  v4-extract-phone accepts the same dimensions and expects a smartphone photo with visible paper around the complete artwork.
  v4-diagnose-phone accepts the same input/key/dimensions and reports Build45 stage diagnostics; -oracle-quad-json is lab-only.
  v4-diagnose-phone-handoff inspects already-qualified Build43 candidates; optional oracle geometry is comparison-only.
  v4-diagnose-phone-frozen expands only the diagnostic frozen bank to 128 candidates; production remains capped at 32.
  v4-diagnose-phone-refine selects up to two proposal-only seeds per side-pair and locally refines them before held-out qualification.
  v4-diagnose-phone-ranking measures proposal-only ranking observables across the Build47 extended bank; no key is used.
  v4-diagnose-phone-refine4 selects up to four proposal-ranked seeds per side-pair, then applies the same proposal-only refinement as Build48.

Capacity options:
  -in FILE           Input JPEG or PNG (required)
  -details           Show per-profile capacities and image dimensions

Analyze options:
  -in FILE           Input JPEG or PNG (required)
  -message TEXT      Message whose UTF-8 byte length should be analyzed
  -bytes N           Payload byte count to analyze instead of -message

Diagnose options:
  -in FILE           Input JPEG or PNG (required)
  -key TEXT          Optional key for an independent baseline HMAC attempt
  -json              Emit machine-readable JSON
  -regions N         Diagnostic region grid, N x N (default 3, maximum 4)
  -max-dim N         Maximum diagnostic pyramid dimension (default 2048)
  -levels N          Maximum diagnostic pyramid levels (default 2, maximum 3)

Run "pixseal <command> -help" to show the options for a command.

Examples:
  pixseal embed -in photo.png -out sealed.png -key "a long secret" -message "hello"
  pixseal extract -in sealed.png -key "a long secret"
  pixseal v4-embed -in photo.png -out sealed-v4.png -key "a long secret" -message "hello"
  pixseal v4-extract -in sealed-v4.png -key "a long secret"
  pixseal v4-extract-projective -in acquired.png -key "a long secret" -width 1632 -height 1632
  pixseal v4-extract-scanner -in scan.jpg -key "a long secret" -width 1632 -height 1632
  pixseal v4-diagnose-phone -in phone.jpg -key "a long secret" -width 1632 -height 1632 -json
  pixseal v4-diagnose-phone-handoff -in phone.jpg -key "a long secret" -width 1632 -height 1632 -json
  pixseal v4-diagnose-phone-frozen -in phone.jpg -key "a long secret" -width 1632 -height 1632 -json
  pixseal v4-diagnose-phone-refine -in phone.jpg -key "a long secret" -width 1632 -height 1632 -json
  pixseal v4-diagnose-phone-ranking -in phone.jpg -width 1632 -height 1632 -json
  pixseal v4-diagnose-phone-refine4 -in phone.jpg -key "a long secret" -width 1632 -height 1632 -json
  pixseal capacity -in photo.png -details
  pixseal analyze -in photo.png -message "hidden message"
  pixseal diagnose -in captured.jpg -json
`, buildinfo.String())

func main() {
	if len(os.Args) < 2 {
		rootUsage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "embed":
		err = embed(os.Args[2:])
	case "extract":
		err = extract(os.Args[2:])
	case "v4-embed":
		err = v4Embed(os.Args[2:])
	case "v4-extract":
		err = v4Extract(os.Args[2:])
	case "v4-extract-projective":
		err = v4ExtractProjective(os.Args[2:])
	case "v4-extract-scanner":
		err = v4ExtractScanner(os.Args[2:])
	case "v4-extract-phone":
		err = v4ExtractPhone(os.Args[2:])
	case "v4-diagnose-phone":
		err = v4DiagnosePhone(os.Args[2:])
	case "v4-diagnose-phone-handoff":
		err = v4DiagnosePhoneHandoff(os.Args[2:])
	case "v4-diagnose-phone-frozen":
		err = v4DiagnosePhoneFrozen(os.Args[2:])
	case "v4-diagnose-phone-refine":
		err = v4DiagnosePhoneRefine(os.Args[2:])
	case "v4-diagnose-phone-ranking":
		err = v4DiagnosePhoneRanking(os.Args[2:])
	case "v4-diagnose-phone-refine4":
		err = v4DiagnosePhoneRefine4(os.Args[2:])
	case "capacity":
		err = capacity(os.Args[2:])
	case "analyze":
		err = analyze(os.Args[2:])
	case "diagnose":
		err = diagnose(os.Args[2:])
	case "help", "-help", "--help", "-h":
		rootUsage()
		return
	default:
		if len(os.Args[1]) > 0 && os.Args[1][0] == '-' {
			fmt.Fprintf(os.Stderr, "Missing command before option %q; did you mean \"pixseal embed %s\"?\n\n", os.Args[1], os.Args[1])
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command %q.\n\n", os.Args[1])
		}
		rootUsage()
		os.Exit(2)
	}

	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func rootUsage() {
	fmt.Fprint(os.Stderr, usageHeader)
}

func newFlagSet(command, summary string) *flag.FlagSet {
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: pixseal %s [options]\n\n%s\n\nOptions:\n", command, summary)
		fs.PrintDefaults()
	}
	return fs
}

func validateSourceDimensions(width, height int) error {
	if width <= 0 || height <= 0 {
		return fmt.Errorf("invalid image dimensions %dx%d", width, height)
	}
	if int64(width) > (1<<63-1)/int64(height) {
		return errors.New("image dimensions overflow the PixSeal pixel-count calculation")
	}
	pixels := int64(width) * int64(height)
	if pixels > maxCLISourcePixels {
		return fmt.Errorf("image has %d pixels; PixSeal CLI safety limit is %d pixels", pixels, maxCLISourcePixels)
	}
	return nil
}

func detectInputFormat(f *os.File) (string, error) {
	var signature [8]byte
	n, err := io.ReadFull(f, signature[:])
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	b := signature[:n]
	if len(b) >= 2 && b[0] == 0xff && b[1] == 0xd8 {
		return "jpeg", nil
	}
	pngSignature := [...]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if len(b) >= len(pngSignature) {
		match := true
		for i := range pngSignature {
			if b[i] != pngSignature[i] {
				match = false
				break
			}
		}
		if match {
			return "png", nil
		}
	}
	return "", errors.New("unsupported image format: PixSeal accepts PNG or JPEG")
}

func openImageConfig(path string) (image.Config, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return image.Config{}, "", err
	}
	defer f.Close()
	format, err := detectInputFormat(f)
	if err != nil {
		return image.Config{}, "", err
	}
	var config image.Config
	switch format {
	case "jpeg":
		config, err = jpeglegacy.DecodeConfig(f)
	case "png":
		config, err = png.DecodeConfig(f)
	default:
		err = errors.New("unsupported image format")
	}
	if err != nil {
		return image.Config{}, "", err
	}
	if err := validateSourceDimensions(config.Width, config.Height); err != nil {
		return image.Config{}, "", err
	}
	return config, format, nil
}

func openImageWithFormat(path string) (image.Image, string, error) {
	config, format, err := openImageConfig(path)
	if err != nil {
		return nil, "", err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	var img image.Image
	switch format {
	case "jpeg":
		img, err = jpeglegacy.Decode(f)
	case "png":
		img, err = png.Decode(f)
	default:
		err = errors.New("unsupported image format")
	}
	if err != nil {
		return nil, "", err
	}
	if got := img.Bounds(); got.Dx() != config.Width || got.Dy() != config.Height {
		return nil, "", fmt.Errorf("decoded dimensions %dx%d differ from image config %dx%d", got.Dx(), got.Dy(), config.Width, config.Height)
	}
	return img, format, nil
}

type boundsOnlyImage struct{ rectangle image.Rectangle }

func (img boundsOnlyImage) ColorModel() color.Model { return color.NRGBAModel }
func (img boundsOnlyImage) Bounds() image.Rectangle { return img.rectangle }
func (img boundsOnlyImage) At(x, y int) color.Color { return color.NRGBA{A: 255} }

func openImage(path string) (image.Image, error) {
	img, _, err := openImageWithFormat(path)
	return img, err
}

func pngOutputPath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	// filepath.Ext(".sealed") is ".sealed". A Unix dotfile with no second dot
	// is a basename, not an extension-only filename.
	if strings.HasPrefix(base, ".") && strings.Count(base, ".") == 1 {
		ext = ""
	}
	if strings.EqualFold(ext, ".png") {
		return path
	}
	if ext == "" {
		return path + ".png"
	}
	return strings.TrimSuffix(path, ext) + ".png"
}

func outputTargetInfo(path string, force bool) (bool, os.FileMode, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return false, 0, fmt.Errorf("output path %s exists but is not a regular file", path)
	}
	if !force {
		return true, info.Mode().Perm(), fmt.Errorf("output file %s already exists; use -force to replace it", path)
	}
	return true, info.Mode().Perm(), nil
}

func createTempForOutput(dir, base string) (*os.File, string, error) {
	for attempt := 0; attempt < 32; attempt++ {
		var random [8]byte
		if _, err := rand.Read(random[:]); err != nil {
			return nil, "", err
		}
		path := filepath.Join(dir, "."+base+".pixseal-"+hex.EncodeToString(random[:]))
		// 0666 is intentionally subject to the process umask. This avoids making a
		// newly-created carrier more permissive than the caller requested.
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
		if err == nil {
			return f, path, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, "", err
		}
	}
	return nil, "", errors.New("could not create a unique temporary output file")
}

func commitNoClobber(temporaryPath, path string) error {
	// A hard link publishes the already-synced temporary inode atomically and
	// fails if any directory entry already exists at the destination.
	if err := os.Link(temporaryPath, path); err == nil {
		_ = os.Remove(temporaryPath)
		return nil
	}
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("output file %s appeared while writing; refusing to overwrite it", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Some filesystems/platforms do not support hard links. Fall back to an
	// exclusive destination create. This preserves no-clobber semantics, though
	// it is not crash-atomic while bytes are copied.
	src, err := os.Open(temporaryPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("output file %s appeared while writing; refusing to overwrite it", path)
		}
		return err
	}
	ok := false
	defer func() {
		_ = dst.Close()
		if !ok {
			_ = os.Remove(path)
		}
	}()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	if err := dst.Sync(); err != nil {
		return err
	}
	if err := dst.Close(); err != nil {
		return err
	}
	ok = true
	_ = os.Remove(temporaryPath)
	return nil
}

func commitForcedRegular(temporaryPath, path string, mode os.FileMode) error {
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// The original target disappeared during encoding. Do not turn -force into
			// permission to clobber a newly-created path in a race; publish no-clobber.
			return commitNoClobber(temporaryPath, path)
		}
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return fmt.Errorf("output path %s changed and is no longer a regular file", path)
	}
	if err := os.Chmod(temporaryPath, mode); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err == nil {
		return nil
	} else if runtime.GOOS != "windows" {
		return err
	}

	// Windows cannot always replace an existing destination with os.Rename.
	// Restrict the backup-and-restore fallback to Windows and only to a regular
	// file that has just been revalidated above.
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	backup, backupPath, err := createTempForOutput(dir, base+".backup")
	if err != nil {
		return err
	}
	if err := backup.Close(); err != nil {
		_ = os.Remove(backupPath)
		return err
	}
	if err := os.Remove(backupPath); err != nil {
		return err
	}
	if err := os.Rename(path, backupPath); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		_ = os.Rename(backupPath, path)
		return err
	}
	_ = os.Remove(backupPath)
	return nil
}

func writePNGAtomic(path string, img image.Image, force bool) error {
	existed, mode, err := outputTargetInfo(path, force)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	temporary, temporaryPath, err := createTempForOutput(dir, base)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = temporary.Close()
			_ = os.Remove(temporaryPath)
		}
	}()

	if err := png.Encode(temporary, img); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}

	if existed {
		err = commitForcedRegular(temporaryPath, path, mode)
	} else {
		err = commitNoClobber(temporaryPath, path)
	}
	if err != nil {
		return err
	}
	committed = true
	return nil
}

func embed(args []string) error {
	fs := newFlagSet("embed", "Hide an authenticated message; the output image is always PNG.")
	in := fs.String("in", "", "input JPEG or PNG file (required)")
	out := fs.String("out", "", "output PNG file (required)")
	message := fs.String("message", "", "message to hide, up to 64 bytes (required)")
	profileName := fs.String("profile", "auto", "v3 profile: auto, robust, balanced or capacity")
	force := fs.Bool("force", false, "replace an existing output file")
	key := fs.String("key", "", "secret key (required, minimum 8 bytes)")
	strength := fs.Float64("strength", 24, "DCT embedding strength from 4 to 120")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q; use -message \"text\"", fs.Arg(0))
	}
	if *in == "" || *out == "" || *key == "" || *message == "" {
		fs.Usage()
		return fmt.Errorf("-in, -out, -key and -message are required")
	}

	if math.IsNaN(*strength) || math.IsInf(*strength, 0) || *strength < 4 || *strength > 120 {
		return fmt.Errorf("-strength must be a finite number from 4 to 120")
	}
	profile, err := watermark.ParseProfile(*profileName)
	if err != nil {
		return err
	}
	outputPath := pngOutputPath(*out)
	// Fast-fail before decoding or embedding a potentially huge input. The final
	// writer repeats the check and commits race-safely.
	if _, _, err := outputTargetInfo(outputPath, *force); err != nil {
		return err
	}
	img, err := openImage(*in)
	if err != nil {
		return err
	}
	marked, embedInfo, err := watermark.EmbedWithInfo(img, []byte(*message), []byte(*key), watermark.Options{
		Strength: *strength,
		Profile:  profile,
	})
	if err != nil {
		return err
	}

	if err := writePNGAtomic(outputPath, marked, *force); err != nil {
		return err
	}
	if outputPath != *out {
		fmt.Printf("output renamed to %s (PixSeal output is PNG)\n", outputPath)
	}
	fmt.Printf("embedded %d bytes using profile %s in %s\n", len([]byte(*message)), embedInfo.Profile, outputPath)
	return nil
}

func extract(args []string) error {
	fs := newFlagSet("extract", "Recover and authenticate a hidden PixSeal message; bounded rotation and supported combined geometry correction are automatic.")
	in := fs.String("in", "", "carrier JPEG or PNG file (required)")
	key := fs.String("key", "", "secret key (required, minimum 8 bytes)")
	raw := fs.Bool("raw", false, "write only the authenticated payload bytes to stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" {
		fs.Usage()
		return fmt.Errorf("-in and -key are required")
	}

	img, err := openImage(*in)
	if err != nil {
		return err
	}
	payload, info, err := watermark.ExtractWithInfo(img, []byte(*key))
	if err != nil {
		return err
	}
	if *raw {
		if _, err := os.Stdout.Write(payload); err != nil {
			return err
		}
		printExtractDiagnostics(os.Stderr, info)
		return nil
	}
	fmt.Printf("%s\n", payload)
	printExtractDiagnostics(os.Stderr, info)
	return nil
}

func printExtractDiagnostics(w io.Writer, info watermark.ExtractInfo) {
	fmt.Fprintf(w, "confidence-margin: %.2f\nprofile: %s\n", info.Confidence, info.Profile)
	if info.RotationCorrectionDegrees != 0 {
		fmt.Fprintf(w, "rotation-correction: %.2f degrees\n", info.RotationCorrectionDegrees)
	}
	if info.ScaleXCorrection != 0 || info.ScaleYCorrection != 0 {
		fmt.Fprintf(w, "scale-correction: x=%.4f y=%.4f\n", info.ScaleXCorrection, info.ScaleYCorrection)
	}
	if info.ShearXCorrection != 0 {
		fmt.Fprintf(w, "shear-x-correction: %.2f degrees\n", math.Atan(info.ShearXCorrection)*180/math.Pi)
	}
	if info.ShearYCorrection != 0 {
		fmt.Fprintf(w, "shear-y-correction: %.2f degrees\n", math.Atan(info.ShearYCorrection)*180/math.Pi)
	}
	if info.PerspectiveCorrection != "" {
		fmt.Fprintf(w, "perspective-correction: %s\n", info.PerspectiveCorrection)
	}
}

func v4Embed(args []string) error {
	fs := newFlagSet("v4-embed", "EXPERIMENTAL Build31 encoder: hide an authenticated Format-v4 message using the development-locked prototype-2 pilot. Output is always PNG.")
	in := fs.String("in", "", "input JPEG or PNG file (required)")
	out := fs.String("out", "", "output PNG file (required)")
	message := fs.String("message", "", "message to hide, up to 64 bytes (required)")
	profileName := fs.String("profile", "auto", "experimental v4 profile: auto, robust, balanced or capacity")
	force := fs.Bool("force", false, "replace an existing output file")
	key := fs.String("key", "", "secret key (required, minimum 8 bytes)")
	strength := fs.Float64("strength", 24, "DCT embedding strength from 4 to 120")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q; use -message \"text\"", fs.Arg(0))
	}
	if *in == "" || *out == "" || *key == "" || *message == "" {
		fs.Usage()
		return fmt.Errorf("-in, -out, -key and -message are required")
	}
	if math.IsNaN(*strength) || math.IsInf(*strength, 0) || *strength < 4 || *strength > 120 {
		return fmt.Errorf("-strength must be a finite number from 4 to 120")
	}
	profile, err := watermark.ParseProfile(*profileName)
	if err != nil {
		return err
	}
	outputPath := pngOutputPath(*out)
	if _, _, err := outputTargetInfo(outputPath, *force); err != nil {
		return err
	}
	img, err := openImage(*in)
	if err != nil {
		return err
	}
	marked, info, err := watermark.ExperimentalV4EmbedWithInfo(img, []byte(*message), []byte(*key), watermark.Options{
		Strength: *strength,
		Profile:  profile,
	})
	if err != nil {
		return err
	}
	if err := writePNGAtomic(outputPath, marked, *force); err != nil {
		return err
	}
	if outputPath != *out {
		fmt.Printf("output renamed to %s (PixSeal output is PNG)\n", outputPath)
	}
	fmt.Printf("EXPERIMENTAL v4: embedded %d bytes using profile %s in %s\n", len([]byte(*message)), info.Profile, outputPath)
	fmt.Printf("pilot: %s sha256=%s\n", info.PilotName, info.PilotHash)
	return nil
}

func v4Extract(args []string) error {
	fs := newFlagSet("v4-extract", "EXPERIMENTAL Build31 aligned decoder: authenticate a Format-v4 message on a native 8px lattice. For blind projective/crop recovery use v4-extract-projective.")
	in := fs.String("in", "", "carrier JPEG or PNG file (required)")
	key := fs.String("key", "", "secret key (required, minimum 8 bytes)")
	raw := fs.Bool("raw", false, "write only the authenticated payload bytes to stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" {
		fs.Usage()
		return fmt.Errorf("-in and -key are required")
	}
	img, err := openImage(*in)
	if err != nil {
		return err
	}
	payload, info, err := watermark.ExperimentalV4ExtractAligned(img, []byte(*key))
	if err != nil {
		return err
	}
	if *raw {
		if _, err := os.Stdout.Write(payload); err != nil {
			return err
		}
		printV4ExtractDiagnostics(os.Stderr, info)
		return nil
	}
	fmt.Printf("%s\n", payload)
	printV4ExtractDiagnostics(os.Stderr, info)
	return nil
}

func printV4ExtractDiagnostics(w io.Writer, info watermark.ExperimentalV4ExtractInfo) {
	fmt.Fprintf(w, "EXPERIMENTAL Format-v4 aligned decode\n")
	fmt.Fprintf(w, "data-confidence: %.2f\nprofile: %s\n", info.Confidence, info.Profile)
	fmt.Fprintf(w, "pilot-score: %.6f\npilot-margin: %.6f\npilot-origin: (%d,%d) blocks\n", info.PilotScore, info.PilotMargin, info.OriginXBlocks, info.OriginYBlocks)
	fmt.Fprintf(w, "pilot: %s sha256=%s\n", info.PilotName, info.PilotHash)
}

func v4ExtractProjective(args []string) error {
	fs := newFlagSet("v4-extract-projective", "EXPERIMENTAL Build36 physical-channel decoder: run Build34 blind projective/crop recovery, then use reliability-aware Hamming decoding and authenticate the Format-v4 frame. Canonical dimensions must match the block-normalized carrier before printing.")
	in := fs.String("in", "", "acquired/scanned JPEG or PNG file (required)")
	key := fs.String("key", "", "secret key (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	raw := fs.Bool("raw", false, "write only the authenticated payload bytes to stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	if *width < 296 || *height < 256 || *width%8 != 0 || *height%8 != 0 {
		return fmt.Errorf("canonical -width/-height must be divisible by 8 and at least 296x256")
	}
	img, err := openImage(*in)
	if err != nil {
		return err
	}
	payload, info, projective, err := watermark.ExperimentalV4ExtractProjective(img, []byte(*key), *width, *height)
	if err != nil {
		printV4ProjectiveDiagnostics(os.Stderr, info, projective)
		return err
	}
	if *raw {
		if _, err := os.Stdout.Write(payload); err != nil {
			return err
		}
		printV4ProjectiveDiagnostics(os.Stderr, info, projective)
		return nil
	}
	fmt.Printf("%s\n", payload)
	printV4ProjectiveDiagnostics(os.Stderr, info, projective)
	return nil
}

func v4ExtractScanner(args []string) error {
	fs := newFlagSet("v4-extract-scanner", "EXPERIMENTAL Build37 scanner decoder: detect the printed artwork boundary, refine a tightly bounded affine mapping with disjoint public-pilot partitions, aggregate a five-geometry pilot-qualified ensemble, then authenticate the unchanged Format-v4 frame.")
	in := fs.String("in", "", "full-page/scanner JPEG or PNG file with visible white paper around the artwork (required)")
	key := fs.String("key", "", "secret key (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	raw := fs.Bool("raw", false, "write only the authenticated payload bytes to stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	if *width < 296 || *height < 256 || *width%8 != 0 || *height%8 != 0 {
		return fmt.Errorf("canonical -width/-height must be divisible by 8 and at least 296x256")
	}
	img, err := openImage(*in)
	if err != nil {
		return err
	}
	payload, info, scanner, err := watermark.ExperimentalV4ExtractScanner(img, []byte(*key), *width, *height)
	if err != nil {
		printV4ScannerDiagnostics(os.Stderr, info, scanner)
		return err
	}
	if *raw {
		if _, err := os.Stdout.Write(payload); err != nil {
			return err
		}
		printV4ScannerDiagnostics(os.Stderr, info, scanner)
		return nil
	}
	fmt.Printf("%s\n", payload)
	printV4ScannerDiagnostics(os.Stderr, info, scanner)
	return nil
}

func printV4ScannerDiagnostics(w io.Writer, info watermark.ExperimentalV4ExtractInfo, s watermark.ExperimentalV4ScannerInfo) {
	fmt.Fprintf(w, "EXPERIMENTAL Format-v4 scanner decode\n")
	fmt.Fprintf(w, "boundary: detected=%t confidence=%.6f\n", s.BoundaryDetected, s.BoundaryConfidence)
	fmt.Fprintf(w, "scanner-geometry-accepted: %t ensemble=%d\n", s.Accepted, s.EnsembleCandidates)
	fmt.Fprintf(w, "scanner-affine: scale=(%.6f,%.6f) shear=(%.6f,%.6f) shift=(%.3f,%.3f)\n", s.ScaleX, s.ScaleY, s.ShearX, s.ShearY, s.ShiftX, s.ShiftY)
	fmt.Fprintf(w, "proposal: %.6f\nvalidation: %.6f\nhypotheses: %d\n", s.ProposalScore, s.ValidationScore, s.HypothesesEvaluated)
	fmt.Fprintf(w, "data-confidence: %.2f\nprofile: %s\n", info.Confidence, info.Profile)
	fmt.Fprintf(w, "pilot-score: %.6f\npilot-margin: %.6f\npilot-origin: (%d,%d) blocks\n", info.PilotScore, info.PilotMargin, info.OriginXBlocks, info.OriginYBlocks)
	fmt.Fprintf(w, "pilot: %s sha256=%s\n", info.PilotName, info.PilotHash)
}

func v4ExtractPhone(args []string) error {
	fs := newFlagSet("v4-extract-phone", "EXPERIMENTAL Build44 smartphone decoder: use deterministic PixSeal JPEG rasterization, preserve Build41/42 behavior, then use the Build43 bounded proposal-only side-pair geometry fallback before held-out qualification and unchanged HMAC authentication.")
	in := fs.String("in", "", "smartphone JPEG or PNG with the complete artwork and visible paper around it (required)")
	key := fs.String("key", "", "secret key (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	raw := fs.Bool("raw", false, "write only the authenticated payload bytes to stdout")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	if *width < 296 || *height < 256 || *width%8 != 0 || *height%8 != 0 {
		return fmt.Errorf("canonical -width/-height must be divisible by 8 and at least 296x256")
	}
	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	decoderID := inputDecoderID(format)
	payload, info, phone, err := watermark.ExperimentalV4ExtractPhone(img, []byte(*key), *width, *height)
	if err != nil {
		printV4PhoneDiagnostics(os.Stderr, decoderID, info, phone)
		return err
	}
	if *raw {
		if _, err := os.Stdout.Write(payload); err != nil {
			return err
		}
		printV4PhoneDiagnostics(os.Stderr, decoderID, info, phone)
		return nil
	}
	fmt.Printf("%s\n", payload)
	printV4PhoneDiagnostics(os.Stderr, decoderID, info, phone)
	return nil
}

func inputDecoderID(format string) string {
	switch format {
	case "jpeg":
		return jpeglegacy.DecoderID
	case "png":
		return "go-image-png"
	default:
		return "unknown"
	}
}

func printV4PhoneDiagnostics(w io.Writer, decoderID string, info watermark.ExperimentalV4ExtractInfo, p watermark.ExperimentalV4PhoneInfo) {
	fmt.Fprintf(w, "EXPERIMENTAL Format-v4 phone decode\n")
	fmt.Fprintf(w, "input-decoder: %s\n", decoderID)
	fmt.Fprintf(w, "working-image: %dx%d downsampled=%t\n", p.WorkingWidth, p.WorkingHeight, p.Downsampled)
	fmt.Fprintf(w, "boundary: detected=%t confidence=%.6f\n", p.BoundaryDetected, p.BoundaryConfidence)
	fmt.Fprintf(w, "projective-basin: found=%t\n", p.ProjectiveBasinFound)
	fmt.Fprintf(w, "phone-geometry-accepted: %t ensemble=%d\n", p.Accepted, p.EnsembleCandidates)
	fmt.Fprintf(w, "proposal: %.6f\nvalidation: %.6f\nhypotheses: %d\n", p.ProposalScore, p.ValidationScore, p.HypothesesEvaluated)
	if p.ResidualAttempted {
		fmt.Fprintf(w, "phone-residual: fitted=%t applied=%t controls=%d rms=%.3f px proposal=%.6f->%.6f validation=%.6f->%.6f\n", p.ResidualFitted, p.ResidualApplied, p.ResidualControls, p.ResidualRMSPixels, p.ResidualProposalBefore, p.ResidualProposalAfter, p.ResidualValidationBefore, p.ResidualValidationAfter)
	}
	fmt.Fprintf(w, "data-decode: attempted=%t soft-hamming-profiles=%d max-confidence=%.2f\n", p.DataDecodeAttempted, p.SoftHammingProfiles, p.MaxDataConfidence)
	if p.Build42DataAttempted {
		fmt.Fprintf(w, "build42-data: attempted=true bank=%d ensembles=%d list-frames=%d authenticated=%t\n", p.Build42BankCandidates, p.Build42EnsemblesTried, p.Build42ListFramesTried, p.Build42DataAuthenticated)
	}
	if p.Build43Attempted {
		fmt.Fprintf(w, "build43-geometry: attempted=true edge-refined=%t pairs=%d/%d selected=%s,%s proposal=%d frozen=%d qualified=%d authenticated=%t\n", p.Build43EdgeRefined, p.Build43PairsSelected, p.Build43PairsScanned, p.Build43Pair0, p.Build43Pair1, p.Build43ProposalCandidates, p.Build43FrozenCandidates, p.Build43QualifiedCandidates, p.Build43Authenticated)
	}
	fmt.Fprintf(w, "hmac: authenticated=%t fallback-attempted=%t fallback-authenticated=%t\n", p.HMACAuthenticated, p.FallbackAttempted, p.FallbackAuthenticated)
	fmt.Fprintf(w, "data-confidence: %.2f\nprofile: %s\n", info.Confidence, info.Profile)
	fmt.Fprintf(w, "pilot-score: %.6f\npilot-margin: %.6f\npilot-origin: (%d,%d) blocks\n", info.PilotScore, info.PilotMargin, info.OriginXBlocks, info.OriginYBlocks)
	fmt.Fprintf(w, "pilot: %s sha256=%s\n", info.PilotName, info.PilotHash)
}

func printV4ProjectiveDiagnostics(w io.Writer, info watermark.ExperimentalV4ExtractInfo, p watermark.ExperimentalV4ProjectiveInfo) {
	fmt.Fprintf(w, "EXPERIMENTAL Format-v4 projective decode\n")
	fmt.Fprintf(w, "geometry-accepted: %t\n", p.Accepted)
	fmt.Fprintf(w, "geometry: angle=%.3f scale=(%.5f,%.5f) inset=(%.5f,%.5f)\n", p.AngleDegrees, p.ScaleX, p.ScaleY, p.TopInset, p.BottomInset)
	fmt.Fprintf(w, "placement-shift: (%.1f,%.1f) px\n", p.PlacementShiftX, p.PlacementShiftY)
	fmt.Fprintf(w, "proposal-objective: %.6f\ngeometry-validation: %.6f\nplacement-validation: %.6f\nstructural-score: %.6f\nhypotheses: %d\n", p.ProposalObjective, p.GeometryValidation, p.PlacementValidation, p.StructuralScore, p.HypothesesEvaluated)
	fmt.Fprintf(w, "data-confidence: %.2f\nprofile: %s\n", info.Confidence, info.Profile)
	fmt.Fprintf(w, "pilot-score: %.6f\npilot-margin: %.6f\npilot-origin: (%d,%d) blocks\n", info.PilotScore, info.PilotMargin, info.OriginXBlocks, info.OriginYBlocks)
	fmt.Fprintf(w, "pilot: %s sha256=%s\n", info.PilotName, info.PilotHash)
}

func capacity(args []string) error {
	fs := newFlagSet("capacity", "Show usable v3 payload capacity in bytes.")
	in := fs.String("in", "", "input JPEG or PNG file (required)")
	details := fs.Bool("details", false, "show per-profile capacities and image dimensions")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" {
		fs.Usage()
		return fmt.Errorf("-in is required")
	}

	config, format, err := openImageConfig(*in)
	if err != nil {
		return err
	}
	geometry := boundsOnlyImage{rectangle: image.Rect(0, 0, config.Width, config.Height)}
	maximum := watermark.Capacity(geometry)
	if !*details {
		fmt.Printf("%d bytes\n", maximum)
		return nil
	}

	fmt.Printf("Image:       %s\n", *in)
	fmt.Printf("Format:      %s\n", strings.ToUpper(format))
	fmt.Printf("Dimensions:  %d x %d\n", config.Width, config.Height)
	for _, profile := range watermark.Profiles() {
		value := profile.MaximumPayload
		if maximum == 0 {
			value = 0
		}
		fmt.Printf("%-12s %d bytes\n", string(profile.Name)+":", value)
	}
	fmt.Printf("maximum:     %d bytes\n", maximum)
	return nil
}

func analyze(args []string) error {
	fs := newFlagSet("analyze", "Analyze carrier geometry and recommend a v3 adaptive profile. Recommendations are advisory.")
	in := fs.String("in", "", "input JPEG or PNG file (required)")
	message := fs.String("message", "", "message whose UTF-8 byte length should be analyzed")
	payloadBytes := fs.Int("bytes", -1, "payload byte count to analyze instead of -message")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" {
		fs.Usage()
		return fmt.Errorf("-in is required")
	}
	hasMessage := *message != ""
	hasBytes := *payloadBytes >= 0
	if hasMessage == hasBytes {
		fs.Usage()
		return fmt.Errorf("exactly one of -message or -bytes is required")
	}
	requested := *payloadBytes
	if hasMessage {
		requested = len([]byte(*message))
	}
	if requested <= 0 {
		return fmt.Errorf("requested payload must be at least 1 byte")
	}

	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	result := watermark.AnalyzeImage(img, requested)

	profile := "none"
	if result.RecommendedProfile != "" {
		profile = string(result.RecommendedProfile)
	}
	fmt.Printf("Image:                     %s\n", *in)
	fmt.Printf("Format:                    %s\n", strings.ToUpper(format))
	fmt.Printf("Dimensions:                %d x %d\n", result.Width, result.Height)
	fmt.Printf("Requested payload:         %d bytes\n", result.RequestedBytes)
	fmt.Printf("Recommended profile:       %s  [deterministic]\n", profile)
	fmt.Printf("Profile capacity:          %d bytes  [deterministic]\n", result.ProfileCapacity)
	fmt.Printf("Tile redundancy:           %.2fx  [deterministic]\n", result.ProfileTileRedundancy)
	fmt.Printf("Average observations/bit:  %.2fx  [deterministic, untransformed carrier]\n", result.AverageObservations)
	fmt.Printf("Image detail:              %s (score %.2f)  [heuristic]\n", result.Detail, result.DetailScore)
	fmt.Printf("Recommended strength:      %.0f  [heuristic]\n", result.RecommendedStrength)
	fmt.Printf("Status:                    %s\n", result.Status)
	for _, warning := range result.Warnings {
		fmt.Printf("Warning:                   %s\n", warning)
	}
	fmt.Println("Note:                      robustness results are experimental; this analysis is not a recovery guarantee")
	return nil
}

func diagnose(args []string) error {
	fs := newFlagSet("diagnose", "Experimental v0.3 local-lattice diagnostics. Lattice evidence is not watermark authentication.")
	in := fs.String("in", "", "input JPEG or PNG file (required)")
	key := fs.String("key", "", "optional key for an independent baseline v3 authentication attempt")
	jsonOutput := fs.Bool("json", false, "emit machine-readable JSON")
	regions := fs.Int("regions", 3, "diagnostic region grid, N x N (1 to 4)")
	maxDimension := fs.Int("max-dim", 2048, "maximum diagnostic pyramid dimension (512 to 4096)")
	levels := fs.Int("levels", 2, "maximum diagnostic pyramid levels (1 to 3)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" {
		fs.Usage()
		return fmt.Errorf("-in is required")
	}
	if *regions < 1 || *regions > 4 {
		return fmt.Errorf("-regions must be from 1 to 4")
	}
	if *maxDimension < 512 || *maxDimension > 4096 {
		return fmt.Errorf("-max-dim must be from 512 to 4096")
	}
	if *levels < 1 || *levels > 3 {
		return fmt.Errorf("-levels must be from 1 to 3")
	}

	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	options := watermark.DefaultDiagnosticOptions()
	options.RegionsX = *regions
	options.RegionsY = *regions
	options.MaxAnalysisDimension = *maxDimension
	options.MaxLevels = *levels
	options.AttemptAuthentication = *key != ""
	report, err := watermark.DiagnoseGeometry(img, []byte(*key), options)
	if err != nil {
		return err
	}
	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}

	fmt.Printf("Image:                    %s\n", *in)
	fmt.Printf("Format:                   %s\n", strings.ToUpper(format))
	fmt.Printf("Dimensions:               %d x %d\n", report.Width, report.Height)
	fmt.Printf("Pyramid levels:           ")
	for index, level := range report.Levels {
		if index > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("1/%d=%dx%d", level.Divisor, level.Width, level.Height)
	}
	fmt.Println()
	fmt.Printf("Print boundary:           detected=%t confidence=%.3f prior=%t\n", report.PrintBoundary.Detected, report.PrintBoundary.Confidence, report.BoundaryPriorUsed)
	if report.AdaptiveEscalated {
		fmt.Printf("Adaptive lattice:         escalated=true divisor=1/%d reason=%s\n", report.AdaptiveDivisor, report.AdaptiveReason)
	}
	if report.ProjectiveEstimate.Available {
		fmt.Printf("Projective initializer:   source=%s confidence=%.3f canonical≈%.1fx%.1f candidates=%d fallback=%t\n",
			report.ProjectiveEstimate.Source, report.ProjectiveEstimate.Confidence, report.ProjectiveEstimate.CanonicalWidthPixels,
			report.ProjectiveEstimate.CanonicalHeightPixels, len(report.ProjectiveEstimate.ScaleCandidates), report.LatticeFirstFallbackUsed)
	}
	fmt.Printf("Local regions:            %d\n", len(report.Regions))
	for _, region := range report.Regions {
		fmt.Printf("  [%d,%d] div=%d u=(%.2f,%.2f) v=(%.2f,%.2f) period=(%.2f,%.2f) angle=%.2f axis=%.2f coherence=%.3f repeat=%.3f confidence=%.3f\n",
			region.RegionX, region.RegionY, region.AnalysisDivisor,
			region.U.X, region.U.Y, region.V.X, region.V.Y,
			region.PeriodU, region.PeriodV, region.OrientationDegrees,
			region.InterAxisDegrees, region.PeriodicCoherence, region.TileRepetitionCoherence, region.Confidence)
	}
	fmt.Printf("Global basis:             u=(%.2f,%.2f) v=(%.2f,%.2f)\n", report.GlobalU.X, report.GlobalU.Y, report.GlobalV.X, report.GlobalV.Y)
	fmt.Printf("Global consistency:       %.3f\n", report.GlobalConsistency)
	fmt.Printf("Consensus regions:        %d/%d (%.3f)\n", report.ConsensusRegions, len(report.Regions), report.ConsensusFraction)
	fmt.Printf("Lattice evidence:         %t  [diagnostic only]\n", report.LatticeEvidence)
	fmt.Printf("Authentication status:    %s\n", report.AuthenticationStatus)
	fmt.Printf("Authenticated payload:    %t\n", report.AuthenticatedPayload)
	if report.ProjectiveAuthentication.CandidatesProbed > 0 {
		fmt.Printf("Projective auth probe:    candidates=%d full-decodes=%d best-sync=%s %.3f (z=%.2f)\n",
			report.ProjectiveAuthentication.CandidatesProbed, report.ProjectiveAuthentication.FullDecodeAttempts,
			report.ProjectiveAuthentication.BestSyncProfile, report.ProjectiveAuthentication.BestSyncFraction,
			report.ProjectiveAuthentication.BestSyncZScore)
		if report.ProjectiveAuthentication.BitDiagnosticAttempts > 0 {
			fmt.Printf("Bit-channel diagnostic:   attempts=%d best=%s/%s/%s known-coded-errors=%d/42 multi-error-words=%d/6 post-ECC-header-errors=%d/24 syndrome=%.3f wrong/correct-margin=%.3f\n",
				report.ProjectiveAuthentication.BitDiagnosticAttempts,
				report.ProjectiveAuthentication.BestBitGeometrySource,
				report.ProjectiveAuthentication.BestBitMode,
				report.ProjectiveAuthentication.BestBitProfile,
				report.ProjectiveAuthentication.BestKnownHeaderErrors,
				report.ProjectiveAuthentication.BestKnownHeaderMulti,
				report.ProjectiveAuthentication.BestPostECCHeaderErrors,
				report.ProjectiveAuthentication.BestSyndromeFraction,
				report.ProjectiveAuthentication.BestWrongMarginRatio)
			if report.ProjectiveAuthentication.SmoothPhaseFitAttempts > 0 {
				fmt.Printf("Smooth phase diagnostic:  attempts=%d eligible=%d resamples=%d decodes=%d best=%s/%s/%s model=%s source=%s controls=%d fit=%.2f loo=%.2f base=%d/24 corrected=%d/24 gain=%d\n",
					report.ProjectiveAuthentication.SmoothPhaseFitAttempts,
					report.ProjectiveAuthentication.SmoothPhaseEligibleFits,
					report.ProjectiveAuthentication.SmoothPhaseResampleAttempts,
					report.ProjectiveAuthentication.SmoothPhaseDecodeAttempts,
					report.ProjectiveAuthentication.BestSmoothPhaseGeometry,
					report.ProjectiveAuthentication.BestSmoothPhaseMode,
					report.ProjectiveAuthentication.BestSmoothPhaseProfile,
					report.ProjectiveAuthentication.BestSmoothPhaseModel,
					report.ProjectiveAuthentication.BestSmoothPhaseControlSource,
					report.ProjectiveAuthentication.BestSmoothPhaseControls,
					report.ProjectiveAuthentication.BestSmoothPhaseFitRMS,
					report.ProjectiveAuthentication.BestSmoothPhaseLOORMS,
					report.ProjectiveAuthentication.BestSmoothPhaseBasePostECC,
					report.ProjectiveAuthentication.BestSmoothPhasePostECC,
					report.ProjectiveAuthentication.BestSmoothPhaseImprovement)
			}
			if report.ProjectiveAuthentication.BestSpatialBlindCells > 0 {
				fmt.Printf("Blind phase observer:    method=%s profile=%s global=(%d,%d) score=%.3f cells=%d mean-confidence=%.3f oracle-distance=%.2f blocks\n",
					report.ProjectiveAuthentication.BestSpatialBlindMethod,
					report.ProjectiveAuthentication.BestSpatialBlindProfile,
					report.ProjectiveAuthentication.BestSpatialBlindGlobalX,
					report.ProjectiveAuthentication.BestSpatialBlindGlobalY,
					report.ProjectiveAuthentication.BestSpatialBlindGlobalScore,
					report.ProjectiveAuthentication.BestSpatialBlindCells,
					report.ProjectiveAuthentication.BestSpatialBlindMeanConfidence,
					report.ProjectiveAuthentication.BestSpatialBlindMeanOracleDistance)
			}
			if report.ProjectiveAuthentication.SpatialDiagnosticAttempts > 0 {
				fmt.Printf("Spatial bit diagnostic: attempts=%d best=%s/%s/%s cells=%d tile-sign-agreement=%.3f tile-unstable=%.3f stable-wrong=%d/42 mixed=%d/42 majority-coded-errors=%d/42 majority-post-ECC=%d/24 gain=%d all-coded-agreement=%.3f unstable=%.3f local-phase-same=%d/%d mean-offset=%.2f max-offset=%.2f local-agreement=%.3f local-unstable=%.3f local-majority-post-ECC=%d/24\n",
					report.ProjectiveAuthentication.SpatialDiagnosticAttempts,
					report.ProjectiveAuthentication.BestSpatialGeometrySource,
					report.ProjectiveAuthentication.BestSpatialMode,
					report.ProjectiveAuthentication.BestSpatialProfile,
					report.ProjectiveAuthentication.BestSpatialCells,
					report.ProjectiveAuthentication.BestSpatialTileAgreement,
					report.ProjectiveAuthentication.BestSpatialTileUnstable,
					report.ProjectiveAuthentication.BestSpatialStableWrongBits,
					report.ProjectiveAuthentication.BestSpatialMixedBits,
					report.ProjectiveAuthentication.BestSpatialMajorityErrors,
					report.ProjectiveAuthentication.BestSpatialMajorityPostECC,
					report.ProjectiveAuthentication.BestSpatialMajorityGain,
					report.ProjectiveAuthentication.BestSpatialAllAgreement,
					report.ProjectiveAuthentication.BestSpatialUnstableFraction,
					report.ProjectiveAuthentication.BestSpatialLocalPhaseSame,
					report.ProjectiveAuthentication.BestSpatialCells,
					report.ProjectiveAuthentication.BestSpatialMeanPhaseOffset,
					report.ProjectiveAuthentication.BestSpatialMaxPhaseOffset,
					report.ProjectiveAuthentication.BestSpatialLocalAgreement,
					report.ProjectiveAuthentication.BestSpatialLocalUnstable,
					report.ProjectiveAuthentication.BestSpatialLocalMajorityECC)
			}
		}
	}
	if report.AuthenticatedPayload {
		fmt.Printf("Authenticated profile:    %s\n", report.AuthenticatedProfile)
		fmt.Printf("Authentication confidence: %.2f\n", report.AuthenticationConfidence)
		fmt.Printf("Authenticated message:    %s\n", report.AuthenticatedMessage)
	}
	fmt.Printf("Timing:                   pyramid=%dms boundary=%dms lattice=%dms projective=%dms auth=%dms total=%dms\n",
		report.Timings.PyramidMilliseconds, report.Timings.PrintBoundaryMilliseconds, report.Timings.LocalLatticeMilliseconds,
		report.Timings.ProjectiveFitMilliseconds, report.Timings.AuthenticationMilliseconds, report.Timings.TotalMilliseconds)
	fmt.Printf("Note:                     %s\n", report.Note)
	return nil
}

type v4Build45QuadFile struct {
	Method  string `json:"method"`
	Matches int    `json:"matches"`
	Inliers int    `json:"inliers"`
	Quad    []struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"quad"`
}

type v4Build45DiagnosticOutput struct {
	Version              string                                             `json:"version"`
	InputDecoder         string                                             `json:"input_decoder"`
	Classification       watermark.ExperimentalV4PhoneBuild45Classification `json:"classification,omitempty"`
	ProductionError      string                                             `json:"production_error,omitempty"`
	AuthenticatedPayload string                                             `json:"authenticated_payload,omitempty"`
	Phone                *watermark.ExperimentalV4PhoneInfo                 `json:"phone,omitempty"`
	Extract              *watermark.ExperimentalV4ExtractInfo               `json:"extract,omitempty"`
	OracleMethod         string                                             `json:"oracle_method,omitempty"`
	OracleMatches        int                                                `json:"oracle_matches,omitempty"`
	OracleInliers        int                                                `json:"oracle_inliers,omitempty"`
	OracleError          string                                             `json:"oracle_error,omitempty"`
	OraclePayload        string                                             `json:"oracle_payload,omitempty"`
	Oracle               *watermark.ExperimentalV4PhoneBuild45OracleInfo    `json:"oracle,omitempty"`
}

func readBuild45OracleQuad(path string) ([4]watermark.ImagePoint, v4Build45QuadFile, error) {
	var q [4]watermark.ImagePoint
	var f v4Build45QuadFile
	data, err := os.ReadFile(path)
	if err != nil {
		return q, f, err
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return q, f, fmt.Errorf("parse oracle quad JSON: %w", err)
	}
	if len(f.Quad) != 4 {
		return q, f, fmt.Errorf("oracle quad JSON must contain exactly 4 points (TL, TR, BL, BR)")
	}
	for i := range q {
		q[i] = watermark.ImagePoint{X: f.Quad[i].X, Y: f.Quad[i].Y}
	}
	return q, f, nil
}

func v4DiagnosePhone(args []string) error {
	fs := newFlagSet("v4-diagnose-phone", "EXPERIMENTAL Build45 diagnostic: decompose the unchanged Build44 smartphone path into geometry/qualification/data-channel failure classes. An optional externally supplied quad is a lab-only oracle and is never a production fallback.")
	in := fs.String("in", "", "smartphone JPEG or PNG (required)")
	key := fs.String("key", "", "secret key used only after geometry qualification / for final oracle authentication (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	oracleQuad := fs.String("oracle-quad-json", "", "lab-only JSON quadrilateral (TL,TR,BL,BR) generated independently from the production decoder")
	oracleOnly := fs.Bool("oracle-only", false, "skip the blind production path and run only the supplied-geometry oracle")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	if *oracleOnly && *oracleQuad == "" {
		return fmt.Errorf("-oracle-only requires -oracle-quad-json")
	}
	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	out := v4Build45DiagnosticOutput{Version: buildinfo.String(), InputDecoder: inputDecoderID(format)}

	if !*oracleOnly {
		payload, info, phone, prodErr := watermark.ExperimentalV4ExtractPhone(img, []byte(*key), *width, *height)
		out.Phone = &phone
		out.Extract = &info
		out.Classification = watermark.ExperimentalV4PhoneBuild45Classify(phone)
		if prodErr != nil {
			out.ProductionError = prodErr.Error()
		} else {
			out.AuthenticatedPayload = string(payload)
		}
	}

	if *oracleQuad != "" {
		q, qf, err := readBuild45OracleQuad(*oracleQuad)
		if err != nil {
			return err
		}
		out.OracleMethod = qf.Method
		out.OracleMatches = qf.Matches
		out.OracleInliers = qf.Inliers
		payload, _, oracle, oracleErr := watermark.ExperimentalV4PhoneBuild45OracleDecode(img, []byte(*key), *width, *height, q)
		out.Oracle = &oracle
		if oracleErr != nil {
			out.OracleError = oracleErr.Error()
		} else {
			out.OraclePayload = string(payload)
		}
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	fmt.Printf("Build45 phone diagnostic\n")
	fmt.Printf("input-decoder: %s\n", out.InputDecoder)
	if out.Phone != nil {
		fmt.Printf("classification: %s\n", out.Classification)
		fmt.Printf("build41: direct-accepted=%t qualified-bank=%d\n", out.Phone.Build41DirectAccepted, out.Phone.Build41QualifiedCandidates)
		fmt.Printf("build43: attempted=%t frozen=%d qualified=%d selected=%s,%s\n", out.Phone.Build43Attempted, out.Phone.Build43FrozenCandidates, out.Phone.Build43QualifiedCandidates, out.Phone.Build43Pair0, out.Phone.Build43Pair1)
		for i, p := range out.Phone.Build43PairRanking {
			fmt.Printf("build43-pair-rank-%d: %s score=%.6f\n", i+1, p.Pair, p.Score)
		}
		fmt.Printf("build42-data: attempted=%t bank=%d ensembles=%d list-frames=%d authenticated=%t\n", out.Phone.Build42DataAttempted, out.Phone.Build42BankCandidates, out.Phone.Build42EnsemblesTried, out.Phone.Build42ListFramesTried, out.Phone.Build42DataAuthenticated)
		fmt.Printf("hmac: %t\n", out.Phone.HMACAuthenticated)
		if out.ProductionError != "" {
			fmt.Printf("production-error: %s\n", out.ProductionError)
		}
	}
	if out.Oracle != nil {
		fmt.Printf("oracle: method=%s matches=%d inliers=%d pilot-qualified=%t proposal=%.6f validation=%.6f pilot=%.6f margin=%.6f origin=(%d,%d) list-frames=%d hmac=%t\n", out.OracleMethod, out.OracleMatches, out.OracleInliers, out.Oracle.PilotQualified, out.Oracle.ProposalScore, out.Oracle.ValidationScore, out.Oracle.PilotScore, out.Oracle.PilotMargin, out.Oracle.OriginXBlocks, out.Oracle.OriginYBlocks, out.Oracle.ListFramesTried, out.Oracle.HMACAuthenticated)
		if out.OracleError != "" {
			fmt.Printf("oracle-error: %s\n", out.OracleError)
		}
	}
	return nil
}

type v4Build46HandoffOutput struct {
	Version       string                                     `json:"version"`
	InputDecoder  string                                     `json:"input_decoder"`
	OracleMethod  string                                     `json:"oracle_method,omitempty"`
	OracleMatches int                                        `json:"oracle_matches,omitempty"`
	OracleInliers int                                        `json:"oracle_inliers,omitempty"`
	Handoff       watermark.ExperimentalV4PhoneBuild46Report `json:"handoff"`
}

func v4DiagnosePhoneHandoff(args []string) error {
	fs := newFlagSet("v4-diagnose-phone-handoff", "EXPERIMENTAL Build46 diagnostic: inspect already-qualified Build43 candidates and the unchanged production ensemble/Build42 thresholds. Optional oracle geometry is comparison-only and never guides blind search or qualification.")
	in := fs.String("in", "", "smartphone JPEG or PNG (required)")
	key := fs.String("key", "", "secret key used only for diagnostic post-qualification single-candidate authentication (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	oracleQuad := fs.String("oracle-quad-json", "", "lab-only JSON quadrilateral used only to measure distance after blind candidates are frozen and qualified")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	out := v4Build46HandoffOutput{Version: buildinfo.String(), InputDecoder: inputDecoderID(format)}
	var oracle *[4]watermark.ImagePoint
	if *oracleQuad != "" {
		q, qf, err := readBuild45OracleQuad(*oracleQuad)
		if err != nil {
			return err
		}
		oracle = &q
		out.OracleMethod = qf.Method
		out.OracleMatches = qf.Matches
		out.OracleInliers = qf.Inliers
	}
	report, err := watermark.ExperimentalV4PhoneBuild46Diagnose(img, []byte(*key), *width, *height, oracle)
	if err != nil {
		return err
	}
	out.Handoff = report
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	fmt.Printf("Build46 phone qualified-handoff diagnostic\n")
	fmt.Printf("input-decoder: %s\n", out.InputDecoder)
	fmt.Printf("classification: %s\n", report.Classification)
	fmt.Printf("build41: direct=%d qualified-bank=%d\n", report.Build41DirectCandidates, report.Build41QualifiedBank)
	fmt.Printf("build43: attempted=%t frozen=%d qualified=%d\n", report.Build43Attempted, report.Build43FrozenCandidates, report.Build43QualifiedCandidates)
	fmt.Printf("thresholds: direct-ensemble=%d available=%t build42-bank=%d available=%t\n", report.DirectEnsembleRequired, report.DirectEnsembleAvailable, report.Build42BankRequired, report.Build42BankAvailable)
	for i, p := range report.Build43PairRanking {
		fmt.Printf("build43-pair-rank-%d: %s score=%.6f\n", i+1, p.Pair, p.Score)
	}
	for _, c := range report.Candidates {
		fmt.Printf("qualified-%d: pair=%s pair-rank=%d proposal=%.6f validation=%.6f pilot=%.6f margin=%.6f origin=(%d,%d) single-profiles=%d single-list-frames=%d single-confidence=%.3f single-hmac=%t", c.Index, c.SourcePair, c.SourcePairRank, c.ProposalScore, c.ValidationScore, c.PilotScore, c.PilotMargin, c.OriginXBlocks, c.OriginYBlocks, c.SingleProfilesTried, c.SingleListFramesTried, c.SingleMaxDataConfidence, c.SingleHMACAuthenticated)
		if c.OracleCompared {
			fmt.Printf(" oracle-mean-error=%.2fpx oracle-max-error=%.2fpx oracle-ratio=%.6f", c.OracleMeanCornerErrorPx, c.OracleMaxCornerErrorPx, c.OracleMeanCornerErrorRatio)
		}
		fmt.Println()
	}
	return nil
}

type v4Build47FrozenOutput struct {
	Version       string                                     `json:"version"`
	InputDecoder  string                                     `json:"input_decoder"`
	OracleMethod  string                                     `json:"oracle_method,omitempty"`
	OracleMatches int                                        `json:"oracle_matches,omitempty"`
	OracleInliers int                                        `json:"oracle_inliers,omitempty"`
	Frozen        watermark.ExperimentalV4PhoneBuild47Report `json:"frozen"`
}

func v4DiagnosePhoneFrozen(args []string) error {
	fs := newFlagSet("v4-diagnose-phone-frozen", "EXPERIMENTAL Build47 diagnostic: preserve the exact Build43 production proposal tier, then add selected-pair depth and lower-ranked side-pair tiers up to 128 diagnostic candidates. Production remains unchanged; optional oracle geometry is post-hoc comparison only.")
	in := fs.String("in", "", "smartphone JPEG or PNG (required)")
	key := fs.String("key", "", "secret key used only for post-qualification diagnostic authentication (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	oracleQuad := fs.String("oracle-quad-json", "", "lab-only JSON quadrilateral used only after the blind frozen bank has been generated")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	out := v4Build47FrozenOutput{Version: buildinfo.String(), InputDecoder: inputDecoderID(format)}
	var oracle *[4]watermark.ImagePoint
	if *oracleQuad != "" {
		q, qf, err := readBuild45OracleQuad(*oracleQuad)
		if err != nil {
			return err
		}
		oracle = &q
		out.OracleMethod, out.OracleMatches, out.OracleInliers = qf.Method, qf.Matches, qf.Inliers
	}
	report, err := watermark.ExperimentalV4PhoneBuild47Diagnose(img, []byte(*key), *width, *height, oracle)
	if err != nil {
		return err
	}
	out.Frozen = report
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	fmt.Printf("Build47 phone frozen-bank diagnostic\n")
	fmt.Printf("input-decoder: %s\n", out.InputDecoder)
	fmt.Printf("frozen: %d/%d proposal-candidates=%d\n", report.FrozenCandidates, report.SuperBankLimit, report.ProposalCandidates)
	for i, p := range report.PairRanking {
		fmt.Printf("pair-rank-%d: %s score=%.6f\n", i+1, p.Pair, p.Score)
	}
	for _, s := range report.Limits {
		fmt.Printf("limit-%d: available=%d qualified=%d production-best=%d validation=%.6f", s.Limit, s.Available, s.Qualified, s.ProductionBestIndex, s.ProductionBestValidation)
		if oracle != nil {
			fmt.Printf(" oracle-nearest=%d mean=%.2fpx", s.OracleNearestIndex, s.OracleNearestMeanErrorPx)
		}
		fmt.Println()
	}
	return nil
}

type v4Build48RefineOutput struct {
	Version      string                                     `json:"version"`
	InputDecoder string                                     `json:"input_decoder"`
	Refine       watermark.ExperimentalV4PhoneBuild48Report `json:"refine"`
}

func v4DiagnosePhoneRefine(args []string) error {
	fs := newFlagSet("v4-diagnose-phone-refine", "EXPERIMENTAL Build48 diagnostic: select up to two seeds per side-pair using proposal-only evidence, locally refine every selected geometry using proposal tiles only, freeze the complete refined bank, then evaluate held-out qualification and diagnostic HMAC. Production remains unchanged.")
	in := fs.String("in", "", "smartphone JPEG or PNG (required)")
	key := fs.String("key", "", "secret key used only after geometry freeze/qualification for diagnostic authentication (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	report, err := watermark.ExperimentalV4PhoneBuild48Diagnose(img, []byte(*key), *width, *height)
	if err != nil {
		return err
	}
	out := v4Build48RefineOutput{Version: buildinfo.String(), InputDecoder: inputDecoderID(format), Refine: report}
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	fmt.Printf("Build48 phone local-refinement diagnostic\n")
	fmt.Printf("input-decoder: %s\n", out.InputDecoder)
	fmt.Printf("frozen=%d seeds=%d seeds-per-pair=%d refine-evals=%d pre-qualified=%d post-qualified=%d post-authenticated=%d\n", report.FrozenCandidates, report.SeedsSelected, report.SeedsPerPair, report.RefineEvaluations, report.PreQualified, report.PostQualified, report.PostAuthenticated)
	for i, p := range report.PairRanking {
		fmt.Printf("pair-rank-%d: %s score=%.6f\n", i+1, p.Pair, p.Score)
	}
	for _, c := range report.Candidates {
		fmt.Printf("candidate-%d: seed=%d pair=%s pair-rank=%d tier=%s cell=%d pre-proposal=%.6f post-proposal=%.6f post-validation=%.6f post-pilot=%.6f post-margin=%.6f post-qualified=%t movement=%.2fpx single-hmac=%t\n", c.Index, c.SeedIndex, c.SourcePair, c.SourcePairRank, c.SourceTier, c.CellRank, c.PreProposal, c.PostProposal, c.PostValidation, c.PostPilotScore, c.PostPilotMargin, c.PostQualified, c.RefineMeanMovementPx, c.SingleHMACAuthenticated)
	}
	return nil
}

type v4Build50RefineOutput struct {
	Version      string                                     `json:"version"`
	InputDecoder string                                     `json:"input_decoder"`
	Refine       watermark.ExperimentalV4PhoneBuild48Report `json:"refine"`
}

func v4DiagnosePhoneRefine4(args []string) error {
	fs := newFlagSet("v4-diagnose-phone-refine4", "EXPERIMENTAL Build50 diagnostic: select up to four seeds per side-pair using the unchanged proposal score, locally refine every selected geometry using proposal tiles only, freeze the complete refined bank, then evaluate held-out qualification and diagnostic HMAC. Production and Build48 remain unchanged.")
	in := fs.String("in", "", "smartphone JPEG or PNG (required)")
	key := fs.String("key", "", "secret key used only after geometry freeze/qualification for diagnostic authentication (required, minimum 8 bytes)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *key == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -key, -width and -height are required")
	}
	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	report, err := watermark.ExperimentalV4PhoneBuild50Diagnose(img, []byte(*key), *width, *height)
	if err != nil {
		return err
	}
	out := v4Build50RefineOutput{Version: buildinfo.String(), InputDecoder: inputDecoderID(format), Refine: report}
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	fmt.Printf("Build50 phone top-4 local-refinement diagnostic\n")
	fmt.Printf("input-decoder: %s\n", out.InputDecoder)
	fmt.Printf("frozen=%d seeds=%d seeds-per-pair=%d refine-evals=%d pre-qualified=%d post-qualified=%d post-authenticated=%d\n", report.FrozenCandidates, report.SeedsSelected, report.SeedsPerPair, report.RefineEvaluations, report.PreQualified, report.PostQualified, report.PostAuthenticated)
	for i, p := range report.PairRanking {
		fmt.Printf("pair-rank-%d: %s score=%.6f\n", i+1, p.Pair, p.Score)
	}
	for _, c := range report.Candidates {
		fmt.Printf("candidate-%d: seed=%d seed-rank=%d pair=%s pair-rank=%d tier=%s cell=%d pre-proposal=%.6f post-proposal=%.6f post-validation=%.6f post-pilot=%.6f post-margin=%.6f post-qualified=%t movement=%.2fpx single-hmac=%t\n", c.Index, c.SeedIndex, c.SeedRankWithinPair, c.SourcePair, c.SourcePairRank, c.SourceTier, c.CellRank, c.PreProposal, c.PostProposal, c.PostValidation, c.PostPilotScore, c.PostPilotMargin, c.PostQualified, c.RefineMeanMovementPx, c.SingleHMACAuthenticated)
	}
	return nil
}

type v4Build49RankingOutput struct {
	Version      string                                     `json:"version"`
	InputDecoder string                                     `json:"input_decoder"`
	Ranking      watermark.ExperimentalV4PhoneBuild49Report `json:"ranking"`
}

func v4DiagnosePhoneRanking(args []string) error {
	fs := newFlagSet("v4-diagnose-phone-ranking", "EXPERIMENTAL Build49 diagnostic: measure proposal-only ranking observables across the corrected Build47 extended bank, freeze all proposal ranks, then annotate held-out/public-pilot qualification. No secret key, HMAC or oracle evidence participates in this command.")
	in := fs.String("in", "", "smartphone JPEG or PNG (required)")
	width := fs.Int("width", 0, "canonical pre-print carrier width in pixels, divisible by 8 (required)")
	height := fs.Int("height", 0, "canonical pre-print carrier height in pixels, divisible by 8 (required)")
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return fmt.Errorf("unexpected positional argument %q", fs.Arg(0))
	}
	if *in == "" || *width == 0 || *height == 0 {
		fs.Usage()
		return fmt.Errorf("-in, -width and -height are required")
	}
	img, format, err := openImageWithFormat(*in)
	if err != nil {
		return err
	}
	report, err := watermark.ExperimentalV4PhoneBuild49Diagnose(img, *width, *height)
	if err != nil {
		return err
	}
	out := v4Build49RankingOutput{Version: buildinfo.String(), InputDecoder: inputDecoderID(format), Ranking: report}
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	fmt.Printf("Build49 phone proposal-ranking diagnostic\ninput-decoder: %s\nfrozen=%d\n", out.InputDecoder, report.FrozenCandidates)
	for _, c := range report.Candidates {
		fmt.Printf("candidate-%d pair=%s rank=%d tier=%s cell=%d proposal=%.6f within-pair=%d fold-min=%.6f fold-gap=%.6f tile-mean=%.6f tile-std=%.6f tile-rank=%d qualified=%t\n", c.Index, c.SourcePair, c.SourcePairRank, c.SourceTier, c.CellRank, c.ProposalScore, c.ProposalRankWithinPair, c.FoldMinScore, c.FoldGap, c.ProposalTileMean, c.ProposalTileStdDev, c.TileRankWithinPair, c.Qualified)
	}
	return nil
}
