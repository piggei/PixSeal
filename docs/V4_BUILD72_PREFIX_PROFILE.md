# Format-v4 Build72 — prefix-stage profiling over qualified Build71

## Status

Build72 is **observability-only**. Build71 remains the current qualified smartphone baseline. Build72 must not be promoted: its fine-grained timers intentionally perturb runtime and exist only to select a later equivalence-preserving optimization target.

## Why this profile is needed

Build71 removed the fourth-generation runaway critical path. In the two retained Build71 runs, the remaining geometry time is dominated by the prefix through `sibling3`, especially on `B/angle`:

- `B/mild` prefix wall: 74,480 / 76,105 ms;
- `B/angle` prefix wall: 184,439 / 194,499 ms.

The next change must therefore identify whether the residual prefix cost is concentrated in one generation/operator or distributed broadly. Parallelizing blindly could add memory pressure/oversubscription without attacking the actual critical path.

## Exact stage order

Build72 profiles these existing Build71 prefix operations:

```text
0  baseline
1  roots
2  single1
3  pair1
4  continuation1
5  sibling1
6  single2
7  pair2
8  continuation2
9  sibling2
10 single3
11 pair3
12 continuation3
13 sibling3
```

It records:

- min/median/max prefix seed duration;
- min/max prefix evaluations and `sibling3` outputs;
- total evaluation count by stage across all selected seeds;
- cumulative worker time by stage across all selected seeds;
- deterministic maximum-evaluation prefix seed identity (`seed index`, pair rank, rank within pair);
- that seed's total evaluations, output count and elapsed time;
- that seed's per-stage evaluations, states and elapsed milliseconds.

The dominant seed is selected by **evaluation count**, not elapsed time, with original seed order as the stable tie break. This keeps the identity deterministic and independent of CPU scheduling.

## Semantics and methodology

Build72 reproduces the exact Build71 two-barrier scheduler:

```text
Build47 freeze
  ↓
Build48 seed selection
  ↓
Build71 bounded parallel prefix through sibling3
  ↓
freeze sibling3 inputs in original order
  ↓
Build71 global generation-four worker pool
  ↓
strict original-index commit
  ↓
Build68 qualification plane reuse
  ↓
Build66 ordered protected-data decode
  ↓
HMAC final authentication only
```

Only `time.Now()/Since` and public counter accumulation are added around the existing prefix calls. Timing, key material, payload, ECC, HMAC, expected message and reference/oracle geometry never alter search, ranking, scheduling, stopping or commit order.

## Regression guarantees

Build72 includes a direct regression comparing the profiled prefix with `experimentalV4PhoneBuild71SeedPrefix` for:

- exact evaluation count;
- exact number of frozen generation-four inputs;
- exact input order;
- quad, homography and proposal equality;
- per-stage evaluation sum equal to total prefix evaluations.

Inherited Build71 seed reassembly still verifies the full prefix+generation-four bank against Build64. Race detection is required because the profiling wrappers run inside the same parallel prefix pool.

## Qualified-host commands

```bash
make clean
make v4-build72-phone-prefix-profile-test
make v4-build72-phone-physical-test
```

Outputs:

```text
v4-phone private/build72-diagnostics/build72-phone-prefix-profile.tsv
v4-phone private/build72-diagnostics/build72-phone-prefix-profile.md
```

The gate compares semantic outcomes/counters with the qualified Build71 line. Build71 elapsed/prefix timings may be included as reference only; Build72 timing is **not** a promotion metric.

## Decision rule

After the physical profile, choose one later implementation experiment only if the profile identifies a clear public-only critical region. Do not add pruning, thresholds, new scores or protected-data-guided early exit. Build71 remains qualified until a separate candidate later proves exact equivalence and reproducible performance.

## Retained Go 1.26.0 physical result — 2026-09-28

The complete nine-photo gate passed exact semantic equivalence to Build71. Because Build72 places fine-grained timers around every prefix operator, its command/wall time remains diagnostic only. The relevant cumulative public-only worker-time distribution is:

| case | prefix worker ms | pair3+cont3+sib3 ms | share | dominant prefix seed | dominant total ms | dominant late-gen3 ms | share |
|---|---:|---:|---:|---:|---:|---:|---:|
| control/mild | 156234 | 80980 | 51.8% | 98 (5/4) | 138912 | 80980 | 58.3% |
| B/mild | 89863 | 41319 | 46.0% | 67 (3/3) | 72021 | 41319 | 57.4% |
| B/angle | 293120 | 194940 | 66.5% | 98 (6/4) | 203169 | 149514 | 73.6% |

`B/mild` remains exactly 79,259 geometry evaluations / 937 bank / 935 qualified / 691 logical decode / 2,120,047 logical frames with HMAC/payload PASS. `B/angle` remains 334,857 / 6,198 / 0 and REJECT. Controls and historical positive paths are unchanged.

This closes Build72 as **semantic-equivalence PASS / target-selection complete** and selects ordered generation-three scheduling for Build73. No geometry pruning, re-ranking or threshold change is justified by this profile.
