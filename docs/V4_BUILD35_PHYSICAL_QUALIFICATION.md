# Build35 Format-v4 physical qualification protocol

Build35 is the handoff from digital geometry qualification to the first controlled **physical** Format-v4 test. The decoder itself is the Build34 blind projective+crop path; Build35 exposes that path through a public experimental API/CLI and freezes a reproducible scanner-first acquisition protocol.

## Why scanner first

A smartphone photograph combines printer/paper degradation with arbitrary camera distance, framing, perspective, lens distortion, autofocus/sharpening and unknown resampling. A scanner can hold the digital-to-paper-to-digital scale close to the decoder's qualified range. The first physical gate therefore uses a scanner so a failure can be attributed more narrowly. Phone captures follow only after this gate is measured.

## Fixture set

Run:

```sh
make v4-physical-fixtures
```

The default pack uses the active **MQ** source because its block-normalized carrier has at least 4x4 complete 37x32 v4 tiles and therefore exercises the Build34 large-carrier decoder path. The source is cropped only to an 8-pixel-aligned canonical extent; no resampling is performed before embedding.

The private output directory contains:

- `pixseal-build35-mq-control.png` — unmarked negative control;
- `pixseal-build35-mq-marked-a.png` — robust, strength 24, payload `v4-b35-phys-a`;
- `pixseal-build35-mq-marked-b.png` — robust, strength 24, payload `v4-b35-phys-b`;
- `build35-v4-physical-fixtures.tsv` — exact dimensions, SHA-256, expected result and locked pilot hash;
- `acquisition-plan.tsv` — canonical dimensions and required scan filenames;
- `README.txt` — machine-generated instructions including the physical print size for the selected PPI.

The public development key is `PixSeal-v4-TestKey-2026`. It is intentionally reproducible test material, not a production secret.

## Print and scan contract

Default density is 300 ppi/dpi.

1. Print all three PNGs with the **same** printer, paper and quality settings.
2. Use actual size / 100%. Disable fit-to-page and automatic photo enhancement where possible.
3. Record printer model, paper type, print quality and software used to print.
4. Scan each print at 300 dpi with geometric/photo enhancement disabled where possible.
5. Save the scan losslessly as PNG.
6. Crop to the printed artwork edges. Do **not** manually perspective-correct, sharpen, denoise or resize.
7. Save files exactly as `scan-control.png`, `scan-marked-a.png`, `scan-marked-b.png` in `v4-physical private/build35-acquired/` (or override the Make variable).

The printed dimensions are derived from the canonical pixel size and PPI and are written to the generated `README.txt`; do not estimate or fit the artwork to a page.

## Qualification

Run:

```sh
make v4-physical-qualification
```

Internally this invokes:

```sh
pixseal v4-extract-projective \
  -in scan-marked-a.png \
  -key PixSeal-v4-TestKey-2026 \
  -width <canonical-width> \
  -height <canonical-height>
```

The canonical dimensions are those of the pre-print block-normalized carrier, **not** the scan dimensions.

A physical PASS requires all three outcomes:

1. marked A authenticates exactly as `v4-b35-phys-a`;
2. marked B authenticates exactly as `v4-b35-phys-b`;
3. the unmarked control does not authenticate.

Pilot score, geometry validation, CRC plausibility or partial ECC recovery are diagnostic only. They never count as success without the v4 HMAC.

## What to return after acquisition

Return the three original scan files without JPEG conversion or intermediate resizing, together with the generated `build35-v4-physical-fixtures.tsv` and a short note containing printer model, paper type, print settings, scanner model and scan dpi. The source fixtures themselves remain private and need not enter a release archive.

## Next gate

If the scanner set authenticates, repeat the same three printed sheets with a frontal smartphone photograph and then an intentionally oblique/perspective photograph. Those become separate physical gates; scanner success must not be silently generalized to phone-camera success.

## Actual first campaign result

The available office scanner could not produce color lossless PNG/TIFF. The retained first corpus therefore consists of full-page 600-dpi color JPEG scans of the same three prints. Build36 proved those files contain enough protected data to authenticate when geometry is supplied independently. Build37 then closes the blind scanner gate on those exact captures: control REJECT, marked A HMAC PASS, marked B HMAC PASS. The next physical channel is smartphone photography of the same paper set.
