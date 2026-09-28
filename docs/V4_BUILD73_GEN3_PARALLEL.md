# Format-v4 Build73 — qualified ordered generation-three parallel baseline

## Status

Build73 is the **current qualified smartphone baseline**. Promotion is supported by two independent retained Go 1.26.0 physical runs with exact 9/9 semantic equivalence and reproducible performance.

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

Correctness remained authoritative. Both independent retained Go 1.26.0 runs passed the complete 9/9 semantic gate.

## Qualification result

| metric | run 1 | run 2 | retained Build71 |
|---|---:|---:|---:|
| full matrix elapsed ms | 1,014,434 | 1,047,729 | 1,263,623 |
| B/mild elapsed ms | 250,790 | 247,690 | 282,958 |
| B/angle elapsed ms | 276,966 | 306,153 | 412,755 |
| B/mild geometry ms | 120,532 | 113,919 | 153,775 |
| B/angle geometry ms | 236,300 | 258,973 | 367,572 |

Two-run mean full-matrix elapsed is **1,031,081.5 ms** (~1.226x / -18.4% versus Build71). `B/mild` averages **249,240 ms** (~1.135x / -11.9%); `B/angle` averages **291,559.5 ms** (~1.416x / -29.4%).

Semantic counters remain exact: `B/mild` 79,259 / 937 / 935 / 691 / 2,120,047 with HMAC/payload PASS; `B/angle` 334,857 / 6,198 / 0 with REJECT. Build73 is promoted as the current qualified smartphone baseline.
