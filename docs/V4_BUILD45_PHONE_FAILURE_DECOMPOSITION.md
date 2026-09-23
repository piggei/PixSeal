# Format-v4 Build45 — B/mild failure decomposition

Build45 is a **diagnostic research build** on top of the qualified Build44 baseline. It does not change Format-v4, the encoder, strength 48, the locked public pilot, protected data mapping, Hamming/ECC, whitening/HMAC domains, Build37 scanner recovery, Build42 data recovery, Build43 side-pair geometry, or Build44 deterministic JPEG ingest.

## Goal

The primary target is `phone-b-mild.jpg`. `phone-b-angle.jpg` remains informational. Before changing production geometry, Build45 determines where each failure occurs:

1. **geometry** — no usable blind geometry bank reaches acceptance;
2. **qualification** — proposal-side geometry exists, but fewer than two frozen candidates survive held-out/public-pilot qualification;
3. **data-channel** — geometry is fully qualified but the unchanged Build42 protected-data path does not authenticate;
4. **recovered** — the production path authenticates normally.

The classification is descriptive telemetry only. It never changes production decisions.

## Blind diagnostic

`v4-diagnose-phone` calls the unchanged production `ExperimentalV4ExtractPhone` path and reports:

- deterministic input decoder;
- Build41 direct acceptance and qualified-bank size;
- Build43 full six-pair proposal-only ranking;
- selected side pairs, frozen count and held-out qualified count;
- Build42 qualified-bank/list-decoder attempts;
- final HMAC result;
- Build45 failure class.

Run the private two-case study with:

```bash
make v4-build45-phone-diagnostic
```

Outputs are written under `v4-phone private/build45-diagnostics/` and remain private.

## Reference-assisted oracle

Build45 also contains a deliberately isolated laboratory oracle. `scripts/build45-reference-register.py` uses OpenCV SIFT + RANSAC against the known digital marked-B carrier to produce a single artwork quadrilateral. That geometry may then be supplied only to:

```text
pixseal v4-diagnose-phone -oracle-quad-json ... -oracle-only
```

The oracle does **not** refine or rank geometry. PixSeal scales the externally supplied quad to its deterministic working raster, measures proposal/held-out/full-pilot evidence, then reads the unchanged protected margins and applies the existing Build42 soft-Hamming/list decoder. HMAC remains final verification.

SIFT/OpenCV is therefore a lab dependency only. It is not linked into PixSeal, not required by normal builds, and cannot be reached from `v4-extract-phone`.

Run both B cases with:

```bash
make v4-build45-phone-oracle-diagnostic
```

or run blind + oracle sequentially:

```bash
make v4-build45-phone-study
```

## Preliminary retained-corpus evidence

A developer lab run on the retained private Build38 corpus produced:

| acquisition | SIFT matches | RANSAC inliers | pilot | origin | proposal | held-out | HMAC |
|---|---:|---:|---:|---:|---:|---:|---:|
| B/mild | 2584 | 2564 | 0.342214 | (0,0) | 0.384618 | 0.333015 | PASS, first list frame |
| B/angle | 3089 | 3069 | 0.261081 | (0,0) | 0.320779 | 0.230254 | PASS, first list frame |

Both recover the exact `v4-b38-phone-b` payload under independently supplied geometry. This is **not a blind qualification**. It is channel-isolation evidence: strength 48, the locked pilot/data mapping, existing Hamming path and HMAC are sufficient for both difficult B photographs when registration is correct.

Therefore Build45 must not tune strength, ECC or HMAC. The next production change, if justified by the blind matrix, belongs in bounded public-evidence geometry acquisition/ranking.

## Research discipline

The production order remains:

```text
structure -> proposal-side geometry -> freeze -> held-out pilot qualification -> protected data -> HMAC
```

The reference oracle is outside that path and exists only to answer a laboratory question after the blind result is recorded.

## Qualification-host blind result

The Build45 study was run on the qualification host after the Build44 commit. Both difficult marked-B captures reach Build43 held-out qualification with exactly one surviving candidate:

| image | frozen | qualified | Build42 bank | Build45 class |
|---|---:|---:|---:|---|
| B/mild | 32 | 1 | 0 | `qualification` |
| B/angle | 28 | 1 | 0 | `qualification` |

The original Build45 class is intentionally coarse. A Build43 singleton is below the production two-geometry ensemble quorum, so the Build43 bank is not promoted downstream; Build42 separately requires at least three geometries. Build46 refines this observation without changing either threshold. See `V4_BUILD46_QUALIFIED_HANDOFF.md`.
