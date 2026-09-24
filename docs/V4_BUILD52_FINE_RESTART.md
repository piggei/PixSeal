# Format-v4 Build52 — fine-first restart / accepted-state retention

Build52 is a research-only optimizer experiment. **At this research checkpoint, Build44 was the latest qualified production milestone.** Production Build43/42 behavior, Format-v4 framing, encoder, locked public pilot, strength 48, ECC/Hamming, whitening/HMAC domains and deterministic Build44 JPEG ingest remain unchanged.

## Evidence entering Build52

Build51 closes the immediate seed-depth and local-direction questions for B/mild:

- the useful blind seed is available proposal-only at `top+left`, side-pair rank 3, seed rank 3;
- its untouched mean oracle error is 34.049 px;
- the standard Build41 coarse-to-fine refinement (`16, 8, 4, 2, 1`) ends at 43.168 px;
- the fixed local stencil contains `corner-2-x +2`, which improves **both** proposal (`+0.006966`) and post-hoc oracle error (`-1.607 px`, to 32.442 px);
- the standard schedule never evaluates that +2 px move from the untouched seed because earlier 4 px accepted moves have already changed the local state;
- B/angle has no better-both stencil sample and remains a distant false-basin control.

The next justified change is therefore optimizer scheduling/state retention, not a new score.

## Build52 question

Can a bounded blind optimizer preserve the fine local region that Build51 proves is proposal-improving before the coarse-to-fine greedy path destroys it?

Build52 tests this without changing the proposal objective. For every unchanged top-4-per-side-pair seed it records two independent proposal-only paths:

1. **baseline** — the unchanged Build41 coarse-to-fine endpoint;
2. **fine restart** — restart from the untouched seed with only `2 px -> 1 px` coordinate steps.

The fine restart retains the untouched seed plus **every accepted intermediate state**, rather than only its final endpoint. This is the smallest bounded experiment that directly tests the Build51 scheduling observation.

## Blind ordering and leakage barrier

For each difficult B acquisition:

1. reproduce the Build47 extended proposal bank from public structure/pilot only;
2. select the same top4 seeds per side pair used by Build50/51, using the unchanged raw proposal score;
3. generate the unchanged Build41 baseline endpoint;
4. independently run the fixed `2, 1` px fine restart from the untouched seed and retain every accepted proposal-improving state;
5. freeze the complete baseline/fine geometry bank for all seeds;
6. only then evaluate held-out/public-pilot qualification;
7. only qualified frozen states may enter diagnostic protected-data/HMAC decoding;
8. only after both blind JSON files exist may the private lab script create SIFT/reference geometry and compute oracle error.

Secret key, payload, ECC, HMAC and reference/SIFT never create, move, rank or accept geometry.

## Why intermediate-state retention matters

Build51 shows that a proposal-improving fine state can be useful even when later proposal ascent leaves the better geometric neighborhood. Returning only the final restart endpoint would therefore repeat the same information loss in a different schedule. Build52 treats accepted fine states as a small frozen diagnostic bank and lets qualification/HMAC observe them only after geometry generation is complete.

This is research observability, not a production quorum change.

## Commands

```bash
make v4-build52-phone-optimizer-test
make v4-build52-phone-optimizer-diagnostic
```

Private outputs belong under `v4-phone private/build52-diagnostics/` and are never packaged:

- `blind/phone-b-mild-restart.json`
- `blind/phone-b-angle-restart.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build52-optimizer-states.tsv`
- `build52-optimizer-summary.tsv`
- `build52-optimizer.md`

## Decision rule

Build52 is successful evidence for the optimizer hypothesis only if the fine restart preserves at least one B/mild state that improves proposal and independent oracle error relative to the untouched seed. Authentication is stronger evidence but is not required to learn whether scheduling/state retention is directionally correct.

Post-hoc classifications are:

- `optimizer-recovery` — at least one frozen fine-restart state authenticates;
- `optimizer-partial-gain` — no fine state authenticates, but at least one improves both proposal and oracle error versus the untouched seed;
- `optimizer-no-gain` — the fine restart fails to preserve a better-both state.

No Build52 result may justify lowering qualification/quorum or changing strength, ECC/Hamming, HMAC, Format-v4 or production geometry without a separate prospective qualification step over the complete Build44 matrix and controls.

## Qualified host result — 2026-09-24

The Go 1.26.0 qualification host ran both Build52 targets successfully. The post-hoc result is:

| image / target | seed oracle | coarse baseline | fine best | best qualified fine | any better-both | HMAC | classification |
|---|---:|---:|---:|---:|---|---|---|
| B/mild c10 | 34.049 px | 43.168 px | 31.747 px | 32.442 px | yes | no | `optimizer-partial-gain` |
| B/angle c18 | 5448.391 px | 5453.906 px | 5448.391 px | n/a | no | no | `optimizer-no-gain` |

For B/mild the retained state at 32.442 px is held-out-qualified and is the same `corner-2-x +2` opportunity identified by Build51. The next accepted 2px move improves proposal and post-hoc oracle further to 31.747 px but loses qualification. Continuing the 1px stage from that later state does not authenticate. This confirms that optimizer scheduling/state retention matters, but the `2 -> 1` greedy restart is not sufficient for recovery.

No Build52 result is promoted to production. At this research checkpoint, Build44 was the latest qualified milestone.
