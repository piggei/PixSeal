# Build68 qualification pixel-plane reuse

Build68 is the **current qualified smartphone baseline**, derived from the qualified Build66 path and selected directly from the retained Build67 physical profile. Its only implementation change is exact reuse of the already-materialized post-freeze `pixelPlane` during full-pilot qualification. Two independent Go 1.26.0 retained-corpus runs passed exact semantic equivalence and reproduced a large qualification-time reduction.

Build66 remains the historical qualified ordered-parallel milestone; Build64 remains the historical deep-recovery semantic checkpoint.

## Evidence from Build67

The Build67 retained-corpus gate passed semantic equivalence on all nine smartphone cases. On `B/mild` it reproduced the qualified logical counters exactly:

- geometry evaluations: `79259`;
- frozen bank: `937`;
- qualified: `935`;
- logical decode candidates: `691`;
- logical list frames: `2120047`;
- HMAC/payload: PASS.

Build67 stage timing on the same run was:

| stage | B/mild |
|---|---:|
| geometry | 288049 ms |
| second plane preparation | 272 ms |
| qualification | 307583 ms |
| ordered decode wall | 35828 ms |
| total deep fallback | 631733 ms |

Only five extra protected-data candidates were physically evaluated beyond the logical first-HMAC prefix (`696` physical versus `691` logical). Summed list/Hamming/HMAC worker time was only `7061 ms`, while qualification alone consumed `307583 ms`.

Inspection of the existing qualification path identified a concrete implementation redundancy: `experimentalV4PhoneBuild41Qualify` already receives a reusable `pixelPlane`, but the full-pilot gate calls `experimentalV4DetectPilotProjective(img, ...)`, which internally executes `newPixelPlane(img)` again for every candidate that reaches the full-pilot stage.

## Build68 change

Build68 adds a separate detector:

```text
experimentalV4DetectPilotProjectivePlane(...)
```

It performs the same full-pilot residue accumulation, cyclic-origin search, scoring and margin computation as the historical detector, but consumes an already materialized `pixelPlane`.

The Build68 deep fallback then uses:

```text
Build65 seed-parallel blind bank
        |
        v
complete Build64-order bank freeze
        |
        v
one post-freeze newPixelPlane(work)
        |
        v
serial Build68 qualification
  - same held-out fold score
  - same proposal/validation thresholds
  - same full-pilot detector arithmetic
  - reused pixelPlane
        |
        v
Build66 candidate-stable parallel decode batches
        |
        v
strict Build64 logical consumption
        |
        v
first logical HMAC success
```

No historical Build41/64/66/67 helper is rewritten to hide the experiment. The Build68 detector and qualifier are separate functions so exact equivalence can be regression-tested directly.

## Frozen semantics

Build68 changes none of the following:

- Format-v4 encoder or wire format;
- locked public pilot;
- strength 48;
- Build64 geometry bounds, proposal objective or bank content;
- Build65 seed scheduling or seed-stable reassembly;
- bank freeze barrier;
- candidate order;
- held-out validation arithmetic;
- proposal/validation/pilot thresholds;
- full-pilot residue/origin/score arithmetic;
- Build66 decode worker count or batch boundaries;
- protected-data sampling;
- profile order;
- Hamming/list enumeration;
- whitening or HMAC domains;
- first logical HMAC winner;
- fallback ordering;
- Build64 logical telemetry.

## Static equivalence regressions

Build68 includes tests that compare the new plane-consuming full-pilot detector against the historical image-consuming detector on the same raster and multiple projective mappings. The complete `ExperimentalV4PilotDetection` value must compare exactly, including floating-point scores and margin.

The Build68 qualifier is also compared directly with `experimentalV4PhoneBuild41Qualify` for the same plane/image/hypothesis and must return identical:

- validation score;
- full-pilot detection;
- evaluation count;
- qualification decision.

## Runtime telemetry

When Build68 deep recovery runs, the CLI adds:

```text
build68-plane-reuse: attempted=true total-ms=... geometry-ms=... plane-prep-ms=... qualification-ms=... decode-wall-ms=... physical-decode-candidates=... speculative-candidates=... physical-profiles=... physical-list-frames=... sampling-worker-ms=... list-worker-ms=...
```

The Build64/65/66 lines remain the authoritative semantic counters. Build68 physical counters remain observational and may include speculative work already executing in the winning Build66 batch.

## Qualification commands

Run on the qualified Go 1.26.0 host:

```sh
make v4-build68-phone-plane-reuse-test
make v4-build68-phone-physical-test
```

The physical gate writes:

```text
v4-phone private/build68-diagnostics/build68-phone-plane-reuse.tsv
v4-phone private/build68-diagnostics/build68-phone-plane-reuse.md
```

The report compares whole-command time against the qualified Build66 matrix and, when the Build67 profile is present, compares Build68 deep/qualification stage time directly against Build67.


## Qualified physical result — 2026-09-28

Build68 completed two independent runs of the retained nine-photo Go 1.26.0 physical gate with exact semantic equivalence to Build66. In both runs:

- all three controls remained rejected;
- `A/front`, `A/mild`, `A/angle` and `B/front` remained authenticated on their historical paths;
- `B/mild` remained exactly `79259` geometry evaluations / `937` frozen / `935` qualified / `691` logical decode candidates / `2120047` logical list frames, with HMAC and payload PASS;
- `B/angle` remained `6198` frozen / `0` qualified and rejected;
- physical ordered-decode work on `B/mild` remained `696` candidates with exactly `5` speculative candidates in both runs.

The qualification optimization reproduced its benefit:

| metric | Build67 / Build66 reference | Build68 run 1 | Build68 run 2 | Build68 mean |
|---|---:|---:|---:|---:|
| B/mild qualification | 307583 ms (Build67) | 58374 ms | 64253 ms | 61313.5 ms |
| B/mild elapsed | 544038 ms (Build66) | 317089 ms | 362286 ms | 339687.5 ms |
| complete matrix elapsed | 2175687 ms (Build66) | 1827596 ms | 1902847 ms | 1865221.5 ms |

The mean B/mild qualification cost is therefore about **5.02x faster / 80.1% lower** than Build67. Mean B/mild command wall time is about **1.60x faster / 37.6% lower** than Build66, and the complete matrix is about **1.166x faster / 14.3% lower**. Geometry timing varies materially between runs and was not changed by Build68, so the promotion decision relies primarily on the repeated qualification reduction plus exact semantic equivalence, not on attributing all wall-clock variation to plane reuse.

**Decision:** Build68 is promoted as the current qualified smartphone baseline. Build66 remains the historical qualified ordered-parallel milestone; Build64 remains the historical deep-recovery semantic checkpoint.

## Qualified result and promotion

The promotion rule was correctness first, then performance. The nine-photo outcomes and all Build64 logical counters remained exactly qualified:

```text
controls       reject
A/front        pass historical path
A/mild         pass historical path
A/angle        pass direct
B/front        pass direct
B/mild         pass, 937 bank / 935 qualified / 691 logical candidates / 2120047 logical frames
B/angle        reject, 6198 bank / 0 qualified
```

Two independent retained-corpus runs met both requirements. Build68 is therefore promoted as the current qualified smartphone baseline.
