# PixSeal v0.3.0 Build44 — deterministic JPEG ingest

## Goal

Build44 removes the accidental dependency between physical recovery and the Go
standard-library JPEG decoder. It does **not** change the watermark algorithm.
Format-v3 stays frozen; Format-v4 encoding, pilot identity, strength 48, data
mapping, Hamming/ECC, whitening/HMAC domains, Build41/43 geometry, Build40
residual handling and the Build42 data/list decoder remain unchanged.

## Why this build exists

During final Build43 qualification the exact same private `phone-a-mild.jpg`
bytes were decoded on the same host with Go 1.25.1 and Go 1.26.0. Go 1.26
replaced `image/jpeg`; the resulting Y/Cb/Cr planes differed enough to move the
capture outside the already-qualified geometry path. Rebuilding the unchanged
Build43 source with Go 1.25.1 restored HMAC authentication.

That is an ingest reproducibility problem, not evidence that geometry should be
retuned. Build44 therefore freezes JPEG rasterization before returning to harder
phone geometry cases.

## Implementation

`internal/jpeglegacy` contains the pure-Go pre-Go-1.26 decoder family used by
the Build43-qualified path. The checked-in source snapshot comes from Go 1.23.2;
on the canonical A/mild photograph its decoded Y/Cb/Cr hashes were observed to
match Go 1.25.1 exactly. The upstream Go copyright headers are retained and the
BSD-style license is copied as `internal/jpeglegacy/LICENSE_GO`.

The CLI now detects the input type from file magic rather than extension:

- JPEG -> `internal/jpeglegacy.DecodeConfig` / `Decode`
- PNG -> `image/png.DecodeConfig` / `Decode`

This also preserves support for historical files whose extension does not match
their actual encoded format.

The stable decoder identifier is:

```text
pixseal-jpeg-pre-go1.26-v1
```

`v4-extract-phone` prints that identifier as `input-decoder` for JPEG input.

## Public deterministic fixture

Build44 includes a small synthetic JPEG under
`internal/jpeglegacy/testdata/legacy-sample.jpg`. The fixture itself and its
legacy decoded Y/Cb/Cr planes are SHA-256 locked by
`TestDeterministicLegacyJPEGFixture`. `TestOpenImageUsesDeterministicLegacyJPEGDecoder`
separately verifies that the normal CLI ingest path uses that same raster.

These tests do not replace the private physical corpus. They make accidental
changes to the deterministic decoder immediately visible in normal source tests.

## Qualification result

Build44 has completed cross-toolchain qualification. Go 1.26.0 is the qualified
Build44 toolchain; Go 1.25.1 remains the historical Build43 reference.

1. Public baseline regression:

   ```sh
   make v4-build44-jpeg-compat-test
   ```

2. Explicit Go 1.26 deterministic-raster regression:

   ```sh
   make v4-build44-go126-jpeg-compat-test
   ```

3. Go 1.26 private physical gate:

   ```sh
   make v4-build44-go126-phone-physical-test
   ```

The third target builds a separate `dist/pixseal-build44-go126` binary and runs
the unchanged Build43 physical requirements:

```text
controls 3/3 reject
A/front   PASS through unchanged Build43 side-pair fallback
A/mild    PASS
A/angle   PASS
B/front   PASS
B/mild    INFO
B/angle   INFO
```

The deterministic-raster gate and private physical gate both passed on the Surface/WSL2 qualification host. The observed physical matrix was:

```text
controls 3/3 reject
A/front   PASS
A/mild    PASS
A/angle   PASS
B/front   PASS
B/mild    INFO/reject
B/angle   INFO/reject
```

Go 1.26.0 is therefore the default qualified Build44 toolchain.

## Security / methodology boundary

Deterministic ingest changes only how JPEG bytes become pixels. It does not add
any new source of evidence to geometry search. Secret key, payload bytes,
protected-data/ECC results and HMAC remain unavailable to geometry proposal,
refinement and ranking. HMAC remains final authentication only.

## Next step after qualification

With Go 1.26.0 now reproducing the physical Build43 matrix, development can resume the
remaining smartphone geometry envelope on `B/mild` and `B/angle` from a stable
raster baseline. No strength increase, threshold relaxation or larger residual
model is justified by this ingest fix alone.
