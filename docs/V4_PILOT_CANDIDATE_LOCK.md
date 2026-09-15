# Format v4 public pilot candidate lock

Status: **development candidate lock — not yet normative**  
Introduced: **v0.3.0-build30**  
Locked candidate: **`prototype-2-search-p64`**  
Candidate identity SHA-256: **`858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`**

## What is locked

Build30 locks the exact public pilot identity selected and reproduced in Build24:

- tile geometry: **37 x 32 DCT blocks**;
- pilot count: **64** positions;
- data-plane count: **1120** remaining positions;
- exact ordered 64-position mask;
- exact ordered balanced sign sequence (**32 +1 / 32 -1**);
- candidate identity hash defined as SHA-256 over tile width/height followed by the ordered `(position, sign)` sequence.

The source-of-truth position/sign arrays remain in `watermark/experimental_v4_pilot_qualification.go`. `make v4-pilot-lock-check` recomputes the identity and structural invariants; accidental mutation is therefore a test failure.

Locked ordered positions (linear index `y*37+x`):

```text
3, 44, 13, 88, 94, 27, 28, 32,
259, 156, 159, 238, 279, 246, 290, 181,
333, 376, 418, 312, 391, 431, 399, 439,
521, 523, 492, 573, 574, 505, 474, 551,
631, 674, 675, 719, 687, 617, 660, 664,
781, 859, 790, 757, 871, 841, 880, 811,
963, 969, 1008, 903, 907, 1022, 990, 920,
1151, 1118, 1119, 1165, 1092, 1061, 1103, 1070
```

Locked ordered signs:

```text
-1, +1, -1, -1, +1, -1, +1, -1,
-1, +1, -1, +1, +1, +1, +1, +1,
+1, +1, +1, +1, -1, +1, -1, -1,
+1, -1, -1, -1, -1, +1, -1, +1,
-1, +1, +1, +1, +1, -1, +1, +1,
+1, -1, -1, -1, -1, -1, +1, +1,
-1, +1, +1, -1, -1, -1, -1, -1,
-1, +1, -1, +1, +1, -1, -1, +1
```

The lock is intentionally **semantic**, not a source-file SHA lock. Comments, refactoring and qualification code may change without changing the pilot identity. A change to tile geometry, any pilot position or any pilot sign changes the identity hash and must not silently reuse this lock.

## What is not locked

This candidate lock does **not** freeze the whole Format v4 design. In particular it does not freeze:

- decoder geometry-search algorithms;
- confidence/acceptance thresholds;
- framing or version marker;
- data-plane bit mapping;
- ECC choice;
- whitening/interleaving;
- payload layout;
- HMAC/authentication framing;
- encoder strength/profile policy;
- GUI/API behavior.

Format v3 remains the only implemented production interoperability baseline.

## Why the candidate is locked now

The same pilot identity has survived every v4 qualification checkpoint without modification:

1. **Build24 — structural search and partial visibility.** Maximum cyclic mask overlap 7, maximum wrong signed correlation 4, zero perfect non-zero cyclic aliases; deterministic search reproduces the same candidate hash.
2. **Build24 — image-domain channel.** Pilot-only synthetic carrier/channel tests pass with negative controls.
3. **Build25 — independently known geometry.** Rotation, anisotropic scale, shear, perspective, crop and photometric variants retain absolute-origin separation on synthetic and photographic development images.
4. **Build26 — blind geometry.** Data-plane repetition proposals plus pilot validation recover bounded rotation/scale/shear/perspective families.
5. **Build27 — unknown placement.** With geometry supplied independently, arbitrary crop and padded placement are recovered using separated proposal/validation evidence.
6. **Build28 — joint affine + crop.** Geometry and crop are simultaneously unknown on synthetic and both local development originals.
7. **Build29 — projective/padded composition.** The decoder adds explicit ACCEPT/SAFE-REJECT behavior rather than lowering thresholds when geometry ranking is weak.
8. **Build30 — known-mapping lock audit.** On both private photographic development originals, direct Build29 projective-crop and affine-padded mappings retain strong pilot score/margin at origin `(0,0)`, including cases that the joint search intentionally rejects. This isolates remaining failures to search/placement rather than pilot identity.

Build30 known-mapping corpus results:

| carrier | transform | marked score | marked margin | negative margin | origin |
|---|---|---:|---:|---:|---:|
| large development original | projective + crop | 0.986195 | 0.559445 | 0.007237 | (0,0) |
| large development original | affine + padded | 0.988352 | 0.553797 | 0.000598 | (0,0) |
| smaller development original | projective + crop | 0.928412 | 0.490432 | 0.010553 | (0,0) |
| smaller development original | affine + padded | 0.932220 | 0.484086 | 0.010305 | (0,0) |

The private image names are intentionally not part of the format contract or source archive.

## Why this is not yet a normative freeze

PixSeal's target channel includes real print/camera and scanner acquisition. There is still **no physical v4 pilot corpus produced by a real v4 encoder**. Promoting the candidate to a normative on-image contract before that evidence would turn an untested physical-channel assumption into an interoperability promise.

Therefore Build30 uses the following status:

> **candidate identity locked for continued development; normative promotion pending physical v4 qualification.**

A future normative promotion may adopt this exact hash without changing the pilot. If physical qualification instead demonstrates a pilot-level weakness, the replacement must receive a new candidate/version identity and repeat the structural/image/physical qualification chain; the Build30 lock must not be edited in place.


## Build31 framing consequence

Build31 removes one reason the pilot could not yet be exercised physically: PixSeal can now emit a real authenticated v4 carrier through the explicit experimental `v4-embed` path. The pilot identity itself is unchanged and remains development-locked at the hash above.

This does **not** automatically promote the pilot to normative status. Build31 digital framing tests establish that the locked pilot coexists correctly with 1120 real data positions and that the v4 frame authenticates on native/JPEG/aligned-crop digital channels. The missing promotion gate is now concrete: create a separately printed Build31 v4 corpus and recover a valid v4 HMAC through the physical camera/scanner channel.

The Build31 frame is regression-guarded but not locked at the same normative level as the pilot candidate. ECC/framing may still change before v4 promotion; any such change must update the explicit Build31 vector/format documentation rather than silently reusing it.
