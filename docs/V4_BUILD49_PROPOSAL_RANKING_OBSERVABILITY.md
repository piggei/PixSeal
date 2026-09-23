# Format-v4 Build49 — proposal-ranking observability

Build49 is a research-only diagnostic checkpoint. **Build44 remains the latest qualified production milestone.** No production candidate generation, geometry refinement, qualification threshold, quorum, Format-v4 field, ECC/Hamming, whitening/HMAC or strength setting changes.

## Evidence entering Build49

The qualification-host Build48 run selected 12 proposal-ranked seeds per image (two per side pair). For `B/mild`, local refinement increased held-out-qualified seeds from 4 to 7 but authenticated none. The oracle-nearest selected seed improved only from 46.202 px to 44.436 px mean corner error. Build47 had already exposed a materially better 34.049 px candidate in the all-pair extension (`top+left`, side-pair rank 3, cell rank 0), so Build48's top-2-per-pair proposal selection had discarded the geometrically best known basin before refinement. `B/angle` remained far from the oracle (~5.45 kpx) and stays informational.

## Question

Why does proposal ranking place the geometrically better B/mild candidate too low?

Build49 does **not** invent a replacement score. It freezes the corrected Build47 extended bank and records independent proposal-only observables for every candidate:

- Build43 proposal score;
- proposal fold 1 and fold 2 separately;
- fold mean, minimum and disagreement;
- a balanced-fold diagnostic (`min(fold1,fold2) - |fold1-fold2|`);
- proposal-tile mean, standard deviation, min/max and positive fraction;
- tile-consistency diagnostic (`tile mean - tile stddev`);
- side-pair robust score;
- Build43 cell mean and robust score;
- global and within-pair ranks.

All these quantities use only proposal folds. No key is accepted by the Build49 CLI command.

After all blind JSON files are fixed, the lab script creates the SIFT/reference oracle and measures candidate corner error. It then reports the oracle-nearest candidate's ranks and the coverage of proposal `top 2 / 4 / 6 / 8` candidates per side pair. Spearman correlations with negative oracle error are post-hoc descriptive evidence only.

## Commands

```bash
make v4-build49-phone-proposal-ranking-test
make v4-build49-phone-proposal-ranking-diagnostic
```

Private output belongs under `v4-phone private/build49-diagnostics/` and is never packaged.

## Decision rule

- Oracle-nearest enters at proposal top 4/6/8: seed-depth pruning is plausible; test a bounded increase diagnostically before changing production.
- Another proposal-only observable ranks the oracle-nearest candidate materially better than raw proposal: test that observable prospectively on controls and all qualified PASS cases.
- No measured proposal-only observable improves ranking: the useful geometry is not identifiable by the present pilot surface; the next build should change the public geometric observable rather than simply widen search/refinement.

`B/mild` is the primary case. `B/angle` remains informational and cannot alone justify production changes.
