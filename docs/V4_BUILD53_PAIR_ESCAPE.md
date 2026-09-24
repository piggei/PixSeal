# Format-v4 Build53 — coupled pair escape at coordinate-local roots

Build53 is a research-only optimizer-neighborhood experiment. **Build44 remains the latest qualified production milestone.** Encoder, Format-v4 layout, locked public pilot, strength 48, ECC/Hamming, whitening/HMAC domains, JPEG ingest, production geometry thresholds and quorum are unchanged.

## Evidence entering Build53

The qualified Build52 host run closes the fine-restart question without recovery:

- B/mild candidate 10 starts at 34.049 px mean post-hoc oracle error;
- the unchanged coarse-to-fine baseline ends at 43.168 px;
- the `2 -> 1 px` fine restart retains a best-oracle state at 31.747 px, but it is not qualified;
- retained state 1 (`corner-2-x +2`) is held-out-qualified at 32.442 px;
- no fine state authenticates, so B/mild is `optimizer-partial-gain`;
- B/angle is `optimizer-no-gain` and remains an informational distant false-basin control.

A local development probe, explicitly **not** a qualification result, shows why another optimizer-only experiment is justified: the qualified B/mild 32.442 px root has no proposal-improving individual +/-1px coordinate move, yet a coupled `corner-1-y -1` + `corner-2-x +1` move improves the same proposal score, remains qualified in the preflight, and reaches about 31.842 px. No preflight pair state authenticates.

## Build53 question

Can the unchanged public proposal objective escape a 1px coordinate-wise local maximum when two 1px coordinates are changed together?

Build53 tests only that question. It does not add a new score and does not use held-out, protected-data or oracle information to decide where to search.

## Blind algorithm

For each unchanged top4-per-side-pair seed:

1. compute the unchanged Build41 coarse-to-fine baseline for comparison;
2. independently run only the 2px coordinate stage for at most two passes, retaining the untouched seed plus every accepted 2px intermediate as a **root**;
3. around every root evaluate the complete 1px single-coordinate stencil: 8 dimensions x 2 signs = 16 possible moves (subject to the existing geometry bound);
4. if **any** single-coordinate move improves proposal, do not run pair escape at that root;
5. only when zero single-coordinate moves improve proposal, classify the root as proposal-defined `coordinate_local`;
6. around that root evaluate every unordered pair of distinct dimensions with both signs: `C(8,2) x 4 = 112` coupled +/-1px moves;
7. retain at most the eight proposal-improving pair states with highest proposal score, using deterministic dimension/sign order as the tie break;
8. repeat for every seed and freeze the complete root/pair geometry bank.

The worst-case retained pair bank is bounded by 8 states per coordinate-local root. Pair scanning is also gated by a blind proposal-only locality condition, so ordinary roots do not pay the 112-combination scan.

## Leakage barrier

Before the freeze barrier, Build53 may use only:

- visible image structure and public pilot information already allowed by Build41;
- the unchanged proposal fold;
- fixed step sizes, dimension order and existing movement bound.

It may **not** use held-out validation, full-pilot qualification, payload, secret key, ECC/Hamming, HMAC, SIFT/reference geometry or expected message contents to create, trigger, rank or accept geometry.

After every root and retained pair state for both difficult B images is frozen:

1. evaluate held-out/full-pilot qualification;
2. run diagnostic protected-data/HMAC decoding only for qualified frozen states;
3. after both blind JSON files exist, generate the private SIFT/reference oracle and compute post-hoc geometry error.

## Commands

```bash
make v4-build53-phone-pair-escape-test
make v4-build53-phone-pair-escape-diagnostic
```

Private outputs belong under `v4-phone private/build53-diagnostics/` and are never packaged:

- `blind/phone-b-mild-pair-escape.json`
- `blind/phone-b-angle-pair-escape.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build53-roots.tsv`
- `build53-pair-states.tsv`
- `build53-pair-escape-summary.tsv`
- `build53-pair-escape.md`

## Post-hoc decision labels

- `pair-recovery` — at least one retained pair state authenticates;
- `pair-qualified-gain` — no HMAC, but a retained pair state improves proposal and independent oracle error relative to its own coordinate-local root and passes unchanged qualification;
- `pair-geometric-gain` — proposal + oracle improve, but the pair state does not qualify;
- `pair-no-proposal-escape` — the target has a coordinate-local root but no coupled move improves proposal;
- `pair-no-oracle-gain` — pair moves improve proposal but not oracle;
- `pair-not-triggered` — the target candidate has no 1px coordinate-local retained 2px root.

These labels are laboratory interpretation only. Even `pair-recovery` would not by itself modify production: any production candidate would require a separate implementation and the complete unchanged Build44 physical gate, including all three control rejections and the four qualified marked passes.
