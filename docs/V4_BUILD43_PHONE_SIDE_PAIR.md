# PixSeal v0.3.0 Build43 — Smartphone side-pair geometry fallback

## Scope

Build43 addresses the remaining blind-registration failure on the private Build38 strength-48 `A/front` smartphone photograph. It does **not** change Format-v4, the encoder, carrier/data mapping, public pilot, strength, Hamming/ECC, whitening, framing, HMAC or the Build42 data/list decoder.

## Security / methodology boundary

Geometry uses only visible image structure and the public pilot. Secret key, payload, decoded bits, ECC success and HMAC never create, refine, rank or select a geometry. HMAC is final authentication only. The candidate bank is fully frozen before the held-out public-pilot fold is evaluated.

## Fallback order

1. Run the qualified Build41 geometry path unchanged.
2. If Build41 accepts, retain Build41/42 behavior unchanged.
3. If Build41 rejects, run Build43 side-pair geometry search.
4. Freeze at most 32 proposal-only candidates.
5. Apply unchanged Build41 held-out/full-pilot qualification.
6. Feed the qualified bank to the unchanged Build42 data/list decoder.
7. Authenticate only with the unchanged Format-v4 HMAC.

## Geometry model

Build43 represents the artwork by four lines and evaluates all six unordered side pairs: top+bottom, top+left, top+right, bottom+left, bottom+right and left+right. Each selected side can receive a normal offset and small angular correction. Complementary sides can use a conservative image-only paper→artwork refinement when robust long-line contrast improves sufficiently.

The search preserves multiple coarse cells, angular families, fine selected-side offsets and joint complementary-side variants. Proposal maxima are not trusted as the sole basin selector; the bounded freeze deliberately preserves geometric diversity.

## Qualification discipline

The fixed Build41 1/3 held-out fold is excluded from candidate generation/ranking. Only after the bank is frozen are candidates accepted/rejected by unchanged Build41 thresholds. This preserves the separation:

```text
structure/public proposal → bounded bank → FREEZE → held-out qualification → data/ECC → HMAC
```

## Qualified Go/JPEG toolchain

Build43 physical qualification is pinned to **Go 1.25.1**. During final qualification, Go 1.26.0 was shown to decode the canonical `phone-a-mild.jpg` into different Y/Cb/Cr samples because Go 1.26 replaced `image/jpeg`. That upstream raster difference changes the public-structure/pilot geometry scores and can cause A/mild to reject before the unchanged Build42 data path. The same Build43 source compiled with Go 1.25.1 restores the expected HMAC PASS.

This is treated as an input-decoder/toolchain compatibility boundary, not as a reason to retune Build43 geometry. Repository `make` targets select Go 1.25.1 and rebuild `dist/pixseal` before physical tests. See [`GO_TOOLCHAIN_COMPATIBILITY.md`](GO_TOOLCHAIN_COMPATIBILITY.md).

## Public regression

```bash
make v4-build43-phone-side-pair-test
```

The public synthetic regression checks the six-pair model, the <=32 freeze budget, authenticated marked recovery and Build42 decoder reuse.

## Private Build38 target

```bash
make v4-build43-phone-physical-test
```

Expected milestone:

```text
control/front  REJECT
control/mild   REJECT
control/angle  REJECT
A/front        PASS via Build43 geometry fallback
A/mild         PASS (existing Build42 path)
A/angle        PASS (existing direct Build41 path)
B/front        PASS (existing direct Build41 path)
B/mild         INFO/reject
B/angle        INFO/reject
```

Private photographs are never packaged in source/release archives.
