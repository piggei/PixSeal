# Build31 experimental Format-v4 framing and encoder

Status: **experimental, non-normative**. Build31 is the first PixSeal build that can render an authenticated v4 payload into a real image. It does not replace Format v3 and it does not yet connect the Build29 blind geometry search to payload decoding.

The public pilot remains the Build30 development-locked candidate:

- name: `prototype-2-search-p64`
- tile: 37x32 DCT blocks (1184 positions)
- pilot positions: 64
- data positions: 1120
- candidate SHA-256: `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`

A change to that pilot identity is not a Build31 framing change. It requires a new pilot identity and requalification under the Build30 lock policy.

## Build31 goals

Build31 deliberately changes one variable at a time. It introduces a real frame, authenticated data symbols and an encoder while retaining the v3 profile sizes and Hamming(7,4) only as an apples-to-apples baseline. A stronger ECC remains an open v4 design question.

The build therefore provides:

- a public structural v4 marker through the locked pilot;
- an authenticated v4 version/profile byte distinct from v3;
- domain-separated whitening and frame authentication;
- an explicit 1120-position data mapping that never overlaps the pilot;
- experimental `v4-embed` and aligned `v4-extract` CLI commands;
- deterministic frame/mapping vectors and v3/v4 cross-version rejection tests;
- local photographic round-trip, JPEG-q82 and block-aligned-crop qualification.

It does **not** yet provide a general v4 extractor for rotated, scaled, projective or arbitrary print-camera images. Build29 geometry/placement research remains separate from the Build31 payload decoder.

## Tile and data mapping

One v4 tile contains 1184 block positions in row-major order:

```text
37 columns x 32 rows = 1184 positions
```

The exact 64 prototype-2 pilot positions are reserved first. The remaining 1120 positions are enumerated in ascending row-major tile-position order and assigned data ordinals `0..1119`.

For a profile with `codedBits` protected bits, data ordinal `d` carries:

```text
codeIndex = (d * 251) mod codedBits
```

251 is coprime with 448, 672 and 1120. This preserves the intended profile repetition without ever writing data into a pilot block.

Current Build31 mapping coverage:

| profile | frame bytes | coded bits | data positions per coded bit |
|---|---:|---:|---:|
| robust | 32 | 448 | 2 or 3 |
| balanced | 48 | 672 | 1 or 2 |
| capacity | 80 | 1120 | exactly 1 |

These profile sizes are retained for controlled comparison with v3; they are not a commitment that normative v4 will keep Hamming(7,4).

## Frame layout

The fixed profile frame is:

```text
offset  size  field
0       2     ASCII magic "PS"
2       1     v4 version/profile byte
3       1     payload length in bytes
4       4     CRC32/IEEE of the payload, big-endian
8       N     payload bytes
8+N     8     truncated HMAC-SHA256 tag
...           zero padding to the fixed profile frame size
```

The authenticated version/profile byte is:

```text
robust   0x41
balanced 0x42
capacity 0x43
```

The upper nibble `0x4` identifies this experimental v4 framing family after key-dependent frame recovery. The public pilot is the pre-key structural marker used before data decoding.

## Authentication

Build31 uses the existing minimum 8-byte key policy. The v4 HMAC is deliberately domain-separated from v3:

```text
HMAC-SHA256(
    key,
    "pixseal-frame-v4" || 0x00 || frame[0 : 8+payloadLength]
)
```

The first 8 bytes of the HMAC are stored in the frame.

The decoder checks HMAC before accepting the payload. CRC32 is checked after HMAC as a damage/integrity consistency check. Pilot score, pilot margin, CRC, header plausibility and ECC success are never authentication by themselves.

## Whitening and ECC baseline

The complete fixed-size frame is converted to MSB-first bits and whitened with the existing SHA-256 counter stream construction, but with a v4-specific label:

```text
pixseal-whiten-v4
```

The whitened frame is then encoded with the existing Hamming(7,4) primitive. This is a Build31 experimental baseline, chosen because it preserves the existing 16/32/64-byte profile ceilings and lets the project compare v3 and v4 with the pilot as the main changed structural variable.

A future ECC change must use a new deterministic vector and repeat channel qualification. It must not be slipped into the existing Build31 vector silently.

## DCT symbols

Pilot and data symbols use the existing PixSeal 8x8 luminance-DCT differential carrier. For the locked pilot:

```text
sign +1 -> embedded bit 1
sign -1 -> embedded bit 0
```

All non-pilot blocks carry the Hamming-protected data bit selected by the mapping above. The whole 37x32 tile repeats periodically over the image.

## Build31 aligned decoder

`ExperimentalV4ExtractAligned` and the CLI command `v4-extract` intentionally support only an already-resolved native 8-pixel lattice in Build31.

The order is:

```text
image on native 8px lattice
    -> public prototype-2 pilot correlation
    -> cyclic tile origin
    -> data-only residue aggregation
    -> try robust/balanced/capacity protected lengths
    -> Hamming decode
    -> v4 dewhitening
    -> v4 header/profile parse
    -> HMAC authentication
    -> CRC consistency check
    -> authenticated payload
```

No profile is supplied to extraction. A profile is accepted only if its v4 frame authenticates.

The Build29 geometry search is not yet used as an oracle or prepended to this function. Connecting those two branches is a later qualification step.

## CLI

The stable commands remain Format v3:

```sh
pixseal embed ...
pixseal extract ...
```

Build31 adds explicit experimental commands:

```sh
pixseal v4-embed \
  -in photo.png \
  -out sealed-v4.png \
  -key Piccotti \
  -message "physical v4 test" \
  -profile robust

pixseal v4-extract \
  -in sealed-v4.png \
  -key Piccotti
```

`v4-extract` is not a print-camera decoder yet. It expects the native 8-pixel lattice to remain aligned. Its purpose is to qualify framing/authentication and to create real v4 carriers for the next private print-camera/scanner experiment.

## Deterministic Build31 test vector

The current experimental vector uses:

```text
key      = "Piccotti"
profile  = robust
payload  = "Build31-vector"
header   = 0x41
```

Fixed 32-byte frame:

```text
5053410ebb936d304275696c6433312d766563746f720084e064803bb6840000
```

SHA-256 of the raw frame:

```text
c28f5a0b85bf14cb1b70dca44041e58dfc41ae5c3020a73226a6555214ea9cb7
```

SHA-256 of the 448 Hamming-protected bits represented as one byte per bit:

```text
b6cc9d25b65ec8f2ffcd40320362d9abd8967113168709a324ffe1cfa8c88445
```

SHA-256 of `frame || protectedBits`:

```text
52f15752f1829e39ecae215932b56ab4108b9790eb9fb093705e357b68a7d0ea
```

This vector is an **experimental compatibility guard**, not a normative Format-v4 standard. A deliberate framing/ECC revision may replace it before v4 promotion, but must do so explicitly and document the compatibility break.

## Build31 qualification result

On the two private development originals, all three profiles authenticate natively. The robust profile also authenticates after JPEG quality 82 and after a block-aligned crop whose cyclic origin is recovered from the public pilot.

The minimum 296x256 one-tile carrier also authenticates a capacity-profile payload in the synthetic unit test.

These results establish that the real v4 frame/data path works on the digital channel. They do not provide physical print-camera evidence.

## Normative boundary

After Build31, two different locks must not be confused:

1. **pilot identity:** development-locked since Build30;
2. **Build31 frame:** deterministic and regression-guarded, but still experimental and revisable.

Normative v4 promotion still requires a separately embedded and printed v4 physical corpus. A valid physical result requires an **HMAC-authenticated v4 payload**, not merely pilot detection. The existing v3 photographs cannot satisfy that gate because they do not contain the v4 pilot or frame.
