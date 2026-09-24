# Format-v4 Build51 — local proposal-surface / refinement-trajectory observability

Build51 is a research-only diagnostic checkpoint. **At this research checkpoint, Build44 was the latest qualified production milestone.** Production Build43 candidate generation, pair count, 32-candidate cap, proposal score, qualification thresholds, Build42 data path, Format-v4 framing, strength 48, ECC/Hamming, whitening and HMAC domains remain unchanged.

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

## Qualification-host result — Build51 closed

The Build51 private diagnostic was completed on the qualified host. The post-hoc oracle summary is:

| metric | B/mild | B/angle |
|---|---:|---:|
| target candidate | 10 | 18 |
| side pair / rank | `top+left` / 3 | `top+bottom` / 5 |
| seed rank within pair | 3 | 3 |
| seed oracle error | **34.049 px** | **5448.391 px** |
| final coarse-to-fine error | 43.168 px | 5453.906 px |
| accepted moves | 9 | 5 |
| accepted oracle-improving | 2 | 0 |
| accepted oracle-worsening | 7 | 5 |
| best geometry merely visited by trace | **25.172 px** | 5434.772 px |
| fixed-stencil better-both sample | **yes** | no |
| classification | `optimizer-opportunity` | `proposal-surface-misaligned` |

For B/mild, the 25.172 px trace state is `corner-2-x +16` from the untouched seed. It is **rejected** because its proposal falls from `0.201986` to `0.076256`; it therefore proves geometric reach but not a proposal-compatible path. More importantly, the fixed stencil contains exactly one sample that improves both observables from the untouched seed: `corner-2-x +2`, proposal `0.201986 -> 0.208953` (`+0.006966`) and oracle error `34.049 -> 32.442 px` (`-1.607 px`).

The unchanged coarse-to-fine schedule does not evaluate that `+2 px` move from the untouched seed. At the 4 px level it first accepts `corner-1-x +4`, then `corner-2-y -4`; when the 2 px level is reached, the optimizer is already in a different local state. This closes the Build51 question: **B/mild contains a real optimizer/scheduling opportunity under the existing proposal score.** B/angle does not show the same evidence and remains informational.

Build51 does not establish that the proposal score is globally correct, nor does it justify any production promotion. The next experiment must remain optimizer-only, preserve the untouched seed long enough to exercise the fine local opportunity, freeze all resulting geometry before held-out qualification, and keep B/angle plus the Build44 matrix as controls.
