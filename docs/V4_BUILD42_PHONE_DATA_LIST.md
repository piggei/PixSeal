# Build42 smartphone qualified-bank / list-decoder checkpoint

Build42 is a **post-geometry** recovery build. It does not change the Format-v4 carrier, strength 48, locked public pilot, data-position mapping, Hamming(7,4), whitening/HMAC domains, scanner decoder, or frozen Format-v3 core.

## Evidence entering Build42

Build41 closed the first blind smartphone milestone with direct authenticated recovery of `A/angle` and `B/front`, while all three controls rejected. `A/mild` was different from the remaining failures: its geometry was already accepted and protected-data decoding was attempted, but HMAC failed. `A/front`, `B/mild` and `B/angle` remained geometry rejects.

Geometry experiments during Build42 (tracked cyclic origin, alternate folds, cross-fit corner polish, one-parameter radial distortion and proposal-only quadratic residuals) did not improve held-out evidence consistently on the geometry-reject captures. Those variants were rejected rather than weakening qualification floors.

## A/mild bank diagnosis

The exact Build41 proposal/freeze/held-out process produces five qualified hypotheses for `A/mild`. No pair of those hypotheses authenticates. However, several deterministic three-geometry margin ensembles reduce the soft-Hamming result to one residual information-bit error. In a representative uniform `[Q0,Q2,Q4]` ensemble the erroneous Hamming word has the authentic nibble as its **second** maximum-likelihood candidate.

This is data-plane ambiguity after accepted geometry, not evidence for reopening geometry or increasing strength.

## Build42 decoder

1. Run the unchanged Build41 geometry search and retain both its normal two-member ensemble and the complete qualified bank.
2. Attempt the normal Build41 soft-Hamming/HMAC decode first.
3. If it authenticates, stop. Build42 is not entered.
4. Otherwise retain at most six already-qualified geometries, ordered only by public qualification evidence.
5. Enumerate deterministic three-geometry subsets and average their signed protected margins uniformly.
6. Run ordinary soft Hamming ML first.
7. If that frame does not authenticate, compute the best and second-best nibble score for every Hamming word.
8. Take only the ten smallest best-vs-second gaps and enumerate the bounded second-best substitution list.
9. Dewhiten, parse and HMAC-check each complete frame candidate.
10. Only if both the direct Build41 decoder and Build42 data fallback fail may the older Build40 residual field be attempted.

The key and HMAC never create, refine or rank geometry. HMAC is the final validator of complete frame candidates after the geometry bank and soft-decision list are already fixed.

## Physical qualification

On the unchanged private nine-photo strength-48 corpus:

| capture | result |
|---|---|
| control/front | REJECT before data decode |
| control/mild | REJECT before data decode |
| control/angle | REJECT before data decode |
| A/front | geometry reject |
| A/mild | **HMAC PASS `v4-b38-phone-a` via Build42** |
| A/angle | **HMAC PASS `v4-b38-phone-a` direct Build41** |
| B/front | **HMAC PASS `v4-b38-phone-b` direct Build41** |
| B/mild | geometry reject |
| B/angle | geometry reject |

The direct A/angle and B/front paths do not enter Build42 and no longer pay the residual-fitting cost. A/mild authenticates on the original ~200 MP JPEG after the normal internal resize.

## Test targets

```sh
make v4-build42-phone-data-test
```

runs the public Build41 geometry regression plus the synthetic Build42 list-decoder test.

The opt-in private physical gate is:

```sh
make v4-build42-phone-physical-test
```

It requires three control rejections, direct A/angle and B/front authentication, and Build42-authenticated A/mild. Other marked captures are recorded but are not promoted to passes.

## Next step

Build43 should focus only on the remaining geometry-reject envelope (`A/front`, `B/mild`, `B/angle`). Build42 provides no evidence for increasing carrier strength, weakening pilot thresholds or expanding the residual model.
