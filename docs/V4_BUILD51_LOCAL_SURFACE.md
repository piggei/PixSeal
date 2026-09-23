# Format-v4 Build51 — local proposal-surface / refinement-trajectory observability

Build51 is a research-only diagnostic checkpoint. **Build44 remains the latest qualified production milestone.** Production Build43 candidate generation, pair count, 32-candidate cap, proposal score, qualification thresholds, Build42 data path, Format-v4 framing, strength 48, ECC/Hamming, whitening and HMAC domains remain unchanged.

## Evidence entering Build51

The qualification-host Build50 run resolved the seed-depth question for B/mild:

| selection | pre nearest | post nearest | pre qualified | post qualified | HMAC |
|---|---:|---:|---:|---:|---:|
| top2 | 46.202 px | 44.436 px | 4 | 7 | 0 |
| top4 | **34.049 px** | **43.168 px** | 7 | 11 | 0 |

The top4 run therefore does include the Build47/49 rank-3 `top+left` basin without oracle guidance, but the existing proposal-only local refinement moves away from the reference geometry. Qualification becomes easier, not more geometrically accurate. B/angle remains a distant false-basin control (~5.45 kpx) and cannot drive tuning.

This closes the immediate hypotheses that the failure is caused by top2 seed pruning or by insufficient seed breadth.

## Build51 question

Why does the unchanged Build41 local refiner move the B/mild 34.049 px seed away from the independently measured reference geometry?

Build51 distinguishes two possibilities without changing the decoder:

1. **proposal-surface misalignment** — accepted proposal-ascent steps systematically move away from the physical geometry, and the local neighborhood contains no point that simultaneously improves proposal score and oracle error;
2. **optimizer/search-direction opportunity** — a deterministic nearby perturbation has both a higher proposal score and lower post-hoc oracle error than the seed, but the current coordinate descent does not reach it.

The labels are post-hoc laboratory interpretations only. They never affect geometry generation, ranking or acceptance.

## Blind trace

For the same blind top4-per-side-pair seeds used by Build50, Build51 exactly mirrors the existing Build41 coordinate-descent schedule (`16, 8, 4, 2, 1` px, up to two passes per step). Every actually evaluated +/- corner-axis move records:

- step size, pass, dimension, corner, axis and delta;
- proposal score before and for the attempted geometry;
- accepted/rejected state;
- working/source quadrilateral before and after;
- held-out validation score, annotated only after the complete trace/stencil geometry bank is frozen.

The production/refinement implementation itself is not replaced.

## Deterministic local stencil

Around each unrefined top4 seed, Build51 samples a fixed proposal-only local stencil:

- each of the eight individual corner x/y dimensions at `-4, -2, +2, +4` px;
- global x/y translation at `-2, +2` px;
- x/y scale at `-2, +2` px;
- x/y shear at `-2, +2` px;
- top/bottom width and left/right height perspective-like modes at `-2, +2` px;
- the unchanged center seed.

The stencil definition is fixed in code and does not depend on validation, key, payload, HMAC or oracle evidence.

## Leakage barrier

For each image:

1. Build47 extended candidate bank is generated from public structure/pilot only;
2. top4 seeds are selected by the unchanged raw proposal score within each side pair;
3. complete refinement traces and local stencils are generated using proposal folds only;
4. all generated geometry is frozen;
5. held-out validation is annotated on trace/stencil samples;
6. full-pilot qualification and diagnostic HMAC are evaluated only on the original/final seed states;
7. only after blind JSON files for both B/mild and B/angle exist does the lab script create SIFT/reference oracle geometry and compute corner errors/correlations.

Reference geometry cannot guide any blind step.

## Commands

```bash
make v4-build51-phone-surface-test
make v4-build51-phone-surface-diagnostic
```

Private output belongs under `v4-phone private/build51-diagnostics/` and is never packaged.

The report contains:

- `build51-surface-seeds.tsv` — before/after state per top4 seed;
- `build51-refinement-trajectory.tsv` — every evaluated coordinate-descent move;
- `build51-local-stencil.tsv` — every deterministic local stencil sample;
- `build51-surface-summary.tsv` — oracle-nearest seed, accepted-step directionality, best local sample, local proposal/error Spearman correlation and post-hoc interpretation;
- `build51-local-surface.md` — human-readable summary.

## Decision rule

For B/mild:

- if a fixed stencil point has **both** higher proposal score and lower oracle error than the seed, investigate search directions/optimizer reach before changing the score;
- if accepted proposal-ascent steps predominantly increase oracle error and no better-both stencil point exists, the local proposal objective is displaced from the true geometry and the next experiment should seek an additional public geometric observable rather than a stronger optimizer;
- if the evidence is mixed, expand observability before any production change.

No Build51 result may justify lowering qualification/quorum or changing strength/ECC/HMAC.
