# Build67 deep-recovery stage profiling

Build67 is an **observability-only research build** derived from the qualified Build66 smartphone baseline. It is not a new baseline and it does not introduce a performance optimization. Its purpose is to determine where the remaining deep-recovery wall time is actually spent before Build68 changes implementation details.

The current qualified smartphone baseline remains **Build66**.

## Frozen semantics

Build67 changes none of the following:

- Format-v4 encoder or wire format;
- locked public pilot;
- strength 48;
- Build64 deep-search bounds, proposal score or retained bank;
- Build65 seed-parallel scheduling or seed-stable bank reconstruction;
- held-out/full-pilot qualification thresholds and order;
- Build66 decode worker count or ordered batch boundaries;
- profile order;
- projective protected-data sampling;
- Hamming/list enumeration;
- whitening or HMAC domains;
- first logical HMAC winner;
- logical Build64 telemetry;
- fallback ordering.

The Build67 recovery path therefore remains conceptually:

```text
Build65 seed-parallel blind bank
        |
        v
complete bank freeze in Build64 order
        |
        v
serial held-out/full-pilot qualification
        |
        v
Build66 candidate-stable parallel batches
        |
        v
strict logical consumption in Build64 order
        |
        v
first logical HMAC success
```

Only timers and physical-work counters are added around those existing operations.

## Why profile now

The qualified same-host Build66 gate reduced B/mild from 828036 ms to 544038 ms and the complete nine-photo matrix from 2372255 ms to 2175687 ms while preserving exact Build64/65 semantics. The remaining cost is still large enough that another scheduling or implementation change should not be chosen from intuition alone.

Build67 answers three questions:

1. how much deep-fallback wall time is geometry generation versus rebuilding the shared pixel plane versus qualification versus protected-data decode;
2. inside protected-data decode, how much worker time is spent sampling projective margins versus list/Hamming/HMAC enumeration;
3. how much extra physical work Build66 performs speculatively in the winning ordered batch beyond the logical first-HMAC prefix.

## New telemetry

The existing qualified lines remain unchanged:

```text
build64-recovery: ...
build65-parallel: attempted=true workers=N order=seed-stable
build66-decode-parallel: attempted=true workers=N order=candidate-stable-batches
```

Build67 adds one line only when the deep fallback runs:

```text
build67-profile: attempted=true total-ms=... geometry-ms=... plane-prep-ms=... qualification-ms=... decode-wall-ms=... physical-decode-candidates=... speculative-candidates=... physical-profiles=... physical-list-frames=... sampling-worker-ms=... list-worker-ms=...
```

Definitions:

- `total-ms`: wall time inside the complete Build67 deep-recovery fallback;
- `geometry-ms`: wall time for the unchanged Build65 blind-bank generation, including its internal pixel-plane construction;
- `plane-prep-ms`: wall time for the unchanged second `newPixelPlane(work)` used by qualification and protected-data decode after the blind bank has been frozen;
- `qualification-ms`: wall time for the unchanged serial held-out/full-pilot qualification pass;
- `decode-wall-ms`: wall time for the unchanged Build66 ordered-parallel candidate batches;
- `physical-decode-candidates`: candidate goroutines actually executed, including speculative members of a winning batch;
- `speculative-candidates`: `physical-decode-candidates - logical Build64 decode-candidates`;
- `physical-profiles`: profiles physically sampled/decoded across all candidate work actually executed;
- `physical-list-frames`: list frames physically enumerated across all executed work, including speculative work;
- `sampling-worker-ms`: sum of per-worker time spent in `experimentalV4ReadProtectedProjectiveMargins`;
- `list-worker-ms`: sum of per-worker time spent in the unchanged Build42 list decoder, including whitening, frame parsing and HMAC authentication.

`sampling-worker-ms` and `list-worker-ms` are **summed worker times**, not wall times. Because Build66 runs candidates concurrently, their sum can exceed `decode-wall-ms`.

## Logical versus physical telemetry

The distinction is intentional.

The qualified Build64 counters remain logical and must stay exact. If B/mild succeeds on logical candidate 691 and the winning Build66 batch also physically evaluates later candidates, those later candidates:

- are counted by Build67 physical telemetry;
- are not counted by Build64 logical telemetry;
- cannot affect payload selection;
- cannot affect max-confidence semantic telemetry;
- cannot affect HMAC winner selection.

This makes speculative scheduling cost observable without changing semantics.

## Required retained-corpus equivalence

The nine-photo Build38 matrix must remain exactly the Build66 matrix:

| image | Build64 evals | bank | qualified | logical decode candidates | logical list frames | result |
|---|---:|---:|---:|---:|---:|---|
| control/front | 18021 | 26 | 0 | 0 | 0 | reject |
| control/mild | 117609 | 810 | 6 | 6 | 18432 | reject |
| control/angle | 21203 | 11 | 0 | 0 | 0 | reject |
| A/front | historical path | - | - | - | - | authenticate |
| A/mild | historical path | - | - | - | - | authenticate |
| A/angle | direct | - | - | - | - | authenticate |
| B/front | direct | - | - | - | - | authenticate |
| B/mild | 79259 | 937 | 935 | 691 | 2120047 | authenticate `v4-b38-phone-b` |
| B/angle | 334857 | 6198 | 0 | 0 | 0 | reject |

The four historical fast authenticated cases must still return before Build64/65/66/67 telemetry appears.

## Build67 gate

Run on the qualified Go 1.26.0 host with the retained private phone corpus:

```sh
make v4-build67-phone-profile-test
make v4-build67-phone-physical-test
```

The physical gate writes:

```text
v4-phone private/build67-diagnostics/build67-phone-profile.tsv
v4-phone private/build67-diagnostics/build67-phone-profile.md
```

If a Build66 matrix is present, the report also shows the whole-command Build66 elapsed time and an informational Build66/Build67 ratio. The script accepts both the corrected Build66 filename and the historical packaging-only accidental `build65-phone-matrix.tsv` filename inside `build66-diagnostics`.

Timing never determines correctness.

## Decision rule after the host run

Build67 is **not promoted** even if the profiling run happens to be faster than Build66. The timers themselves perturb execution and Build67 exists only to choose the next implementation experiment.

Use the retained B/mild profile to select Build68:

- if `list-worker-ms` dominates physical worker time, optimize deterministic list/Hamming/HMAC implementation without changing the list or its order;
- if `sampling-worker-ms` dominates, optimize projective margin sampling/buffer reuse without changing sampled positions or arithmetic semantics;
- if `plane-prep-ms` is material, test equivalence-preserving pixel-plane reuse so Build68 avoids rebuilding the same raster after the blind bank freezes;
- if `qualification-ms` dominates negative cases, profile/reduce implementation overhead inside the existing qualification computation without changing thresholds or candidate order;
- if `geometry-ms` remains dominant, optimize allocation/caching within the exact Build64/65 search rather than introducing new pruning;
- if speculative work is material, consider a scheduling-only Build68 experiment that reduces wasted physical work while retaining strict Build66 logical order.

No Build68 optimization should be chosen before the Build67 physical profile is available.
