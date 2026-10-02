# Format-v4 Build77 — generation-four internal-cost profiling

Build77 is **observability-only** over the qualified Build76 smartphone baseline. Build76 remains the authority. Fine-grained timers intentionally perturb runtime, so Build77 cannot be promoted from timing.

## Why profile generation four

After Build76, another barrier shift is no longer justified by the measured scheduling data. On the second retained Build76 run:

```text
B/mild: gen4 wall 18,365 ms, worker 146,841 ms  (~8.00x)
B/angle: gen4 wall 110,801 ms, worker 884,306 ms (~7.98x)
```

Generation four is the dominant B/angle wall region, but its worker pool is already almost perfectly saturated. The next useful question is therefore **which existing operation consumes the work**, not how to schedule more workers.

## Frozen algorithm

Build77 preserves exactly:

- Build75 ordered-parallel Build47 freeze;
- Build76 prefix through `sibling1` and ordered gen2 pool;
- Build73 ordered gen3 pool;
- Build71 gen4 task creation, worker count and original-index commit;
- proposal/validation scores, thresholds, search bounds and seed depth;
- Build68 qualification-plane reuse;
- Build66 ordered protected-data decode and first-logical-HMAC semantics.

Protected data, key, ECC and HMAC cannot create, schedule, rank or retain geometry.

## New measurements

For every existing gen4 task Build77 times the unchanged:

1. `single4` call (`experimentalV4PhoneBuild53SingleProbe`);
2. `pair4` call (`experimentalV4PhoneBuild53PairStencil`) when single4 accepts;
3. each `continuation4` call (`experimentalV4PhoneBuild55Continue`).

The profile records cumulative worker milliseconds, evaluation counts, call/output counts and the deterministic maximum-evaluation gen4 task. The dominant task is selected by evaluation count with original task index as tie-break, never by timing, oracle distance or HMAC outcome.

## Equivalence contract

The direct regression requires Build77 and Build76 to return the same:

- geometry evaluation count;
- selected seed/worker-compatible view;
- blind bank size and order;
- quad and homography;
- proposal and validation values.

## Physical gate

```bash
make clean
make v4-build77-phone-gen4-profile-test
make v4-build77-phone-physical-test
```

Diagnostics:

```text
v4-phone private/build77-diagnostics/build77-phone-gen4-profile.tsv
v4-phone private/build77-diagnostics/build77-phone-gen4-profile.md
```

The gate still requires the exact qualified nine-photo semantics. Timing is diagnostic only and must be used solely to select a later equivalence-preserving implementation experiment.
