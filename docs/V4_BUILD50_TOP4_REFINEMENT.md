# Format-v4 Build50 — top-4-per-pair local refinement

Build50 is a research-only diagnostic checkpoint. **At this research checkpoint, Build44 was the latest qualified production milestone.** No production candidate generation, proposal score, held-out threshold, quorum, Build42 data path, Format-v4 field, ECC/Hamming, whitening/HMAC domain or strength setting changes.

## Evidence entering Build50

The qualification-host Build49 run located the Build47 B/mild oracle-nearest candidate at 34.049 px mean corner error (`top+left`, side-pair rank 3, all-pair-extension, cell rank 0). Its raw proposal rank is **3 within its side pair** and 17 globally. The alternative proposal-only observables tested by Build49 rank the same candidate worse within the pair: fold-min 7, balanced-fold 12 and tile-consistency 13.

Therefore the immediate Build48 miss has a simple explanation: Build48 retained only the proposal top 2 candidates per side pair. The best known B/mild basin was excluded by exactly one within-pair rank. Build49 does not justify replacing the raw proposal score; among the measured observables, the existing raw proposal rank is already the best for this candidate.

B/angle remains informational. Its oracle-nearest candidate is still ~5.45 kpx from the reference geometry even though it is also proposal rank 3 within its side pair.

## Build50 question

If the unchanged proposal score is allowed to keep the blind **top 4** seeds per side pair, does local projective refinement materially improve the B/mild 34.049 px basin and can the refined geometry pass held-out qualification and diagnostic HMAC?

Build50 deliberately tests seed depth before inventing a new score or widening production search.

## Ordering and leakage barrier

For each acquisition:

1. generate the corrected Build47 extended bank using public structure/pilot only;
2. rank candidates independently within each side pair using the unchanged proposal score;
3. keep at most four seeds per pair;
4. refine every selected seed using proposal tiles only;
5. freeze the complete refined bank;
6. evaluate held-out/public-pilot qualification;
7. only for qualified candidates, run diagnostic single-candidate data/HMAC;
8. only after both blind JSON files exist, generate SIFT/reference oracle geometry and measure pre/post corner error.

The secret key cannot create, rank, move or qualify geometry. Reference/SIFT information is post-hoc only.

## Exact top2/top4 comparison

Build50 assigns each selected seed `seed_rank_within_pair`. The private report compares the nested subsets `rank <= 2` and `rank <= 4` **inside the same blind top-4 run**. This avoids a second geometry search and makes the Build48-like top2 and Build50 top4 rows directly comparable.

Build48 remains frozen at two seeds per pair. Build43 production remains frozen at two side pairs and at most 32 candidates.

## Commands

```bash
make v4-build50-phone-top4-refine-test
make v4-build50-phone-top4-refine-diagnostic
```

Private output belongs under `v4-phone private/build50-diagnostics/` and is never packaged.

## Decision rule

For B/mild:

- top4 includes the ~34 px basin and refinement reduces it substantially toward a few pixels, with held-out/HMAC success: continue toward a prospective bounded recovery experiment;
- top4 includes the ~34 px basin but refinement leaves it near ~30–40 px: seed pruning is no longer the immediate bottleneck; investigate refinement/proposal-surface reach;
- top4 changes qualification counts but not geometric precision: do not weaken qualification or quorum merely to accept more candidates.

B/angle is telemetry only and cannot independently justify production changes.

## Qualification-host outcome

The Build50 study was executed on the qualified Go 1.26.0 host. B/mild top4 did include the expected 34.049 px rank-3 seed, proving Build49's coverage prediction. The nearest post-refinement geometry was **43.168 px**, not an improvement. Pre/post held-out-qualified counts were 7/11 and no candidate authenticated. The nested top2 view remained 46.202 -> 44.436 px with 4/7 qualified and HMAC 0.

B/angle remained a distant false basin (5448.391 px pre, 5452.835 px nearest post) with no authentication.

Therefore Build50 closes the seed-depth hypothesis negatively: including the correct known basin is necessary but insufficient, and the current proposal-only local refinement moves that basin away from the independent physical geometry. Build51 studies the local proposal surface and exact refinement trajectory without changing production.
