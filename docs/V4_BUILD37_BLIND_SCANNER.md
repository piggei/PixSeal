# Build37 blind scanner qualification

Build37 is the first PixSeal Format-v4 checkpoint that authenticates the existing physical paper/scanner corpus **blindly**, without registering the scan against the digital original.

## Physical corpus

The corpus is unchanged from Build35:

- one unmarked MQ control;
- one robust/strength-24 carrier with payload `v4-b35-phys-a`;
- one robust/strength-24 carrier with payload `v4-b35-phys-b`;
- public development key `PixSeal-v4-TestKey-2026`.

All three were printed at actual size from the same 1632x1632 canonical MQ carrier on smooth white low-porosity paper. The available office scanner produced full-page 600-dpi color JPEG captures. These private scans are regression evidence and are never shipped in source archives.

## Decoder stages

`ExperimentalV4ExtractScanner` / `v4-extract-scanner` perform:

1. detect the large artwork rectangle against the white page;
2. fit four robust boundary lines and derive a canonical-to-scan seed homography;
3. search a tightly bounded scanner-affine correction;
4. use even-indexed public pilot symbols for proposal;
5. use odd-indexed public pilot symbols as held-out ranking evidence;
6. require the held-out winner to recover complete-pilot cyclic origin `(0,0)` with explicit scanner score/margin floors;
7. freeze five pilot-qualified nearby geometries;
8. aggregate signed protected-bit DCT margins across the frozen ensemble using held-out validation weights;
9. apply Build36 soft Hamming, v4 dewhitening/frame parsing and HMAC.

The key, payload bytes, frame header, CRC, ECC outcome and HMAC result are unavailable to stages 1-7. HMAC therefore verifies the already-frozen result; it does not select geometry.

## Current gates

The scanner-specific experimental floors are:

```text
pilot-A proposal      >= 0.30
pilot-B validation    >= 0.15
complete pilot score  >= 0.15
complete pilot margin >= 0.05
origin                == (0,0)
ensemble size         == 5
```

They are regression gates for the current scanner path, not normative Format-v4 constants.

## Private physical result

Typical Build37 telemetry on the original full-page scans:

| capture | proposal | validation | pilot score | margin | origin | result |
| --- | ---: | ---: | ---: | ---: | --- | --- |
| control | ~0.306 | ~0.093 | ~0.080 | ~0.005 | non-zero | REJECT |
| marked A | ~0.476 | ~0.294 | ~0.297 | ~0.174 | `(0,0)` | HMAC PASS `v4-b35-phys-a` |
| marked B | ~0.483 | ~0.199 | ~0.213 | ~0.088 | `(0,0)` | HMAC PASS `v4-b35-phys-b` |

This is the first fully blind `print -> paper -> scanner -> JPEG -> PixSeal -> HMAC` success in the v4 branch.

## Commands

For one scan:

```sh
pixseal v4-extract-scanner \
  -in full-page-scan.jpg \
  -key PixSeal-v4-TestKey-2026 \
  -width 1632 -height 1632
```

For the private three-file regression:

```sh
make v4-build37-physical-scanner-test \
  V4_SCANNER_CONTROL=/path/to/control.jpg \
  V4_SCANNER_MARKED_A=/path/to/marked-a.jpg \
  V4_SCANNER_MARKED_B=/path/to/marked-b.jpg
```

A PASS requires control rejection and exact HMAC recovery of both expected marked payloads.

## Scope

Build37 qualifies this controlled scanner channel only. It does not prove frontal or oblique smartphone-camera recovery. The same printed sheets should be reused for those next tests so printer/paper variables remain fixed.
