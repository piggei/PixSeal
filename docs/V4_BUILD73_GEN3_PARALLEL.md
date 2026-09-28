# Format-v4 Build73 — ordered generation-three parallel candidate

## Status

Build73 is an **equivalence-preserving performance candidate** over the qualified Build71 smartphone baseline. Build71 remains qualified until Build73 passes exact physical equivalence and its performance benefit is reproduced.

## Evidence from Build72

Build72 preserved the complete Build71 physical matrix under Go 1.26.0 and measured the remaining prefix through `sibling3` without changing scheduling. Its decisive result is that the late third generation is the dominant remaining prefix region on the hard cases:

| case | prefix worker ms | pair3+cont3+sib3 ms | share | dominant-seed total ms | dominant late-gen3 ms | share |
|---|---:|---:|---:|---:|---:|---:|
| control/mild | 156234 | 80980 | 51.8% | 138912 | 80980 | 58.3% |
| B/mild | 89863 | 41319 | 46.0% | 72021 | 41319 | 57.4% |
| B/angle | 293120 | 194940 | 66.5% | 203169 | 149514 | 73.6% |

The signal is structural rather than timing-only: the same stages also carry most of the late-prefix evaluations, particularly on `B/angle`.

## Scheduling change

Build73 leaves all geometry functions unchanged and changes only where parallel work is scheduled:

```text
Build47 freeze
    ↓
24 seeds: exact Build64 prefix through sibling2
    ↓ barrier 1
freeze all sibling2 hypotheses in exact seed/traversal order
    ↓
one bounded global pool
single3 → pair3 → continuation3 → sibling3
    ↓ barrier 2
commit generation-three result slices by frozen task index
    ↓
freeze sibling3 hypotheses in exact resulting order
    ↓
qualified Build71 bounded generation-four pool
single4 → pair4 → continuation4
    ↓ barrier 3
commit final bank strictly in original order
```

There are no nested worker pools. `GOMAXPROCS` bounds each phase, and only one pool is active at a time.

## Equivalence requirements

Build73 must preserve exactly:

- seed selection and seed order;
- geometry evaluations;
- final bank size and final bank order;
- each retained quad, homography and proposal score;
- Build68 qualification plane reuse and qualification decisions;
- Build66 candidate-stable protected-data batches;
- logical decode-candidate/list-frame telemetry;
- first logical HMAC success and payload.

A direct synthetic regression reconstructs one complete Build64 seed search as:

```text
Build73 prefix2 + ordered gen3 + Build71 gen4
```

and compares evaluation count, bank length, order, quad, homography and proposal score element-by-element with `experimentalV4PhoneBuild64SeedBank`.

## Physical gate

On the qualified Go 1.26.0 host:

```bash
make clean
make v4-build73-phone-gen3-parallel-test
make v4-build73-phone-physical-test
```

The retained diagnostics are written to:

```text
v4-phone private/build73-diagnostics/build73-phone-gen3-parallel.tsv
v4-phone private/build73-diagnostics/build73-phone-gen3-parallel.md
```

Correctness is authoritative. Performance is compared only after the complete 9/9 semantic gate passes. A first strong result is not sufficient for promotion; the run must be repeated before Build73 can supersede Build71.
