# Format-v4 Build70 — dominant-seed genealogy profiling

## Status

Build70 is **observability-only**. The current qualified smartphone baseline remains **Build68**.

The qualified Go 1.26.0 Build69 gate preserved the complete Build68 semantic matrix and showed that the seed-parallel geometry region is dominated by one branch: `B/mild` measured 168398 ms max-seed time inside 169014 ms seed wall time, while `B/angle` measured 693395/695422 ms. Similar concentration appears in all deep-recovery controls. The median seed is below one second, so simply increasing the number of top-level seed workers cannot remove the critical path.

## What Build70 measures

Build70 reproduces the Build69/Build68 geometry bank and ordering. It adds only:

- Build47 freeze evaluation count;
- total evaluations inside the 24 selected seed branches;
- dominant seed frozen index (the dominant seed is the maximum-evaluation seed; ties keep original seed order);
- dominant seed side-pair rank and rank within that pair;
- dominant seed evaluation count, final-bank contribution and elapsed time;
- dominant-seed per-stage evaluation counts;
- dominant-seed per-stage state counts.

Evaluation-stage order is fixed as:

```text
baseline, roots,
single1, pair1, continuation1, sibling1,
single2, pair2, continuation2, sibling2,
single3, pair3, continuation3, sibling3,
single4, pair4, continuation4
```

State-stage order is:

```text
roots,
pair1, continuation1, sibling1,
pair2, continuation2, sibling2,
pair3, continuation3, sibling3,
pair4, continuation4,
final bank
```

No protected-data information, key, ECC, HMAC or expected message contributes to these counters or to geometry generation.

## Invariants

Build70 must preserve exactly:

- 24 selected seeds in every deep-recovery case;
- Build64 geometry-evaluation totals;
- final bank size and order;
- full-pilot qualification results;
- logical decode-candidate/list-frame counts;
- payload and HMAC outcomes;
- Build68 qualification plane reuse;
- Build66 candidate-stable ordered decode.

Build70 must not change proposal score, search bounds, keep counts, thresholds, pruning or ordering.

## Qualified-host commands

```bash
make clean
make v4-build70-phone-seed-genealogy-test
make v4-build70-phone-physical-test
```

The private gate writes:

```text
v4-phone private/build70-diagnostics/build70-phone-seed-genealogy.tsv
v4-phone private/build70-diagnostics/build70-phone-seed-genealogy.md
```

## Decision rule for Build71

Do not optimize merely because one stage has many evaluations. First confirm on the physical profile whether the runaway seed expands primarily at one generation or accumulates cost broadly. Build71 should make one equivalence-preserving implementation/scheduling change and must reconstruct the exact Build68 bank in original order before qualification.

## Physical result — 2026-09-28

The qualified Go 1.26.0 retained nine-photo gate passed exact semantic equivalence.

The dominant public-only branches were:

- `control/front`: seed 93, pair/rank 5/2, 12782/17901 seed evaluations;
- `control/mild`: seed 98, pair/rank 5/4, 112978/117483;
- `control/angle`: seed 64, pair/rank 3/1, 6962/21075;
- `B/mild`: seed 67, pair/rank 3/3, 73922/79131;
- `B/angle`: seed 98, pair/rank 6/4, 267775/334749.

The largest common independent region is generation four. `single4 + pair4 + continuation4` contributes 41776 evaluations (56.5%) of the dominant `B/mild` seed and 196641 (73.4%) of dominant `B/angle`, fed by 593 and 1363 `sibling3` states. This selects Build71 ordered generation-four parallel scheduling as the next experiment; no pruning or geometry retuning is justified by Build70.
