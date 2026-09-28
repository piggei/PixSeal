# Format-v4 Build71 — ordered generation-four parallel scheduling

## Status

Build71 is an **equivalence-preserving performance candidate** over the current qualified smartphone baseline, **Build68**. Build71 is not promoted by construction; promotion requires the retained nine-photo Go 1.26.0 physical gate and reproducible performance evidence.

## Evidence from Build69/70

Build69 showed that top-level seed parallelism is strongly load-imbalanced. Build70 then identified the dominant branch deterministically by public geometry-evaluation count and decomposed it by generation.

Physical Build70 results:

| case | seed evals | dominant seed | pair/rank | dominant evals | dominant share | `sibling3` inputs | gen4 evals | gen4 share of dominant |
|---|---:|---:|---|---:|---:|---:|---:|---:|
| control/front | 17901 | 93 | 5/2 | 12782 | 71.4% | 33 | 2864 | 22.4% |
| control/mild | 117483 | 98 | 5/4 | 112978 | 96.2% | 666 | 57120 | 50.6% |
| control/angle | 21075 | 64 | 3/1 | 6962 | 33.0% | 10 | 800 | 11.5% |
| B/mild | 79131 | 67 | 3/3 | 73922 | 93.4% | 593 | 41776 | 56.5% |
| B/angle | 334749 | 98 | 6/4 | 267775 | 80.0% | 1363 | 196641 | 73.4% |

`gen4 evals` means the exact Build64 `single4 + pair4 + continuation4` work. These stages are attractive because each `sibling3` hypothesis defines an independent proposal-only subtree and the current algorithm merely appends each subtree output in traversal order.

## Scheduling change

Build71 changes **only scheduling**.

Qualified Build64/68 logical traversal:

```text
seed 0 prefix ... sibling3[0] -> generation4 -> append
                     sibling3[1] -> generation4 -> append
...
seed 1 prefix ...
```

Build71:

```text
all selected seeds
      |
      +-- bounded parallel prefix through sibling3
      |
      +-- freeze every sibling3 input in exact seed/traversal order
      |
      +-- one bounded global generation-four worker pool
      |      single4 -> pair4 -> continuation4
      |
      +-- store each result by frozen task index
      |
      +-- concatenate strictly by original task index
      |
      +-- complete bank FREEZE
             |
             +-- unchanged Build68 qualification plane reuse
             +-- unchanged Build66 ordered protected-data decode
             +-- HMAC final authentication only
```

This two-barrier design avoids nested worker pools. It lets the CPUs that became idle after short seed prefixes assist the expensive generation-four region without changing the qualified bank order.

## Invariants

Build71 must preserve exactly:

- Build47 freeze and Build48 seed selection;
- 24 selected seeds in retained deep-recovery cases;
- every Build64 proposal-only operation and bound;
- total geometry-evaluation count;
- final bank size, hypothesis values and ordering;
- Build68 full-pilot qualification and plane reuse;
- qualification counts and order;
- Build66 candidate-stable decode batches;
- logical decode-candidate and list-frame counts;
- payload/HMAC outcomes;
- control rejection and `B/angle` rejection.

Build71 does **not** use protected payload, ECC, HMAC, key material, reference/oracle geometry or expected messages to choose tasks or schedule geometry.

## Telemetry

Build71 reports:

- prefix worker count;
- generation-four worker count;
- generation-four task count;
- geometry plane-prep and freeze time;
- prefix wall and cumulative worker time;
- generation-four wall and cumulative worker time;
- generation-four task min/median/max duration;
- generation-four task min/max evaluation count;
- generation-four task min/max bank contribution.

The inherited Build64/65/66/68 logical telemetry remains authoritative for semantic equivalence.

## Regressions

The Build71 unit regression reconstructs a complete seed serially from:

```text
Build71 prefix + ordered Build71 generation-four subtrees
```

and compares it directly to `experimentalV4PhoneBuild64SeedBank` for:

- evaluation count;
- bank size;
- hypothesis order;
- quad;
- homography;
- proposal score.

Race detection is also required because the generation-four workers share one read-only `pixelPlane`.

## Qualified-host commands

```bash
make clean
make v4-build71-phone-gen4-parallel-test
make v4-build71-phone-physical-test
```

The physical gate writes:

```text
v4-phone private/build71-diagnostics/build71-phone-gen4-parallel.tsv
v4-phone private/build71-diagnostics/build71-phone-gen4-parallel.md
```

## Promotion rule

Build68 remains the baseline unless Build71:

1. reproduces the complete 9/9 retained physical semantic matrix;
2. preserves every Build64 logical counter exactly;
3. preserves Build68 qualification/decode behavior and HMAC outcomes;
4. provides a useful performance gain on the same qualified host;
5. reproduces that gain before promotion.

If Build71 is semantically exact but the generation-four wall region remains imbalanced, the next experiment should profile the generation-four task distribution before introducing another scheduling layer.
