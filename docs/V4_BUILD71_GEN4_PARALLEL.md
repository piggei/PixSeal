# Format-v4 Build71 — qualified ordered generation-four parallel baseline

## Status

Build71 is the **current qualified smartphone physical-recovery baseline**. It preserves the complete Build68/Build66/Build64 semantics and changes only proposal-only scheduling inside the deep-recovery geometry search. Promotion is supported by two independent retained nine-photo runs on the qualified Go 1.26.0 host.

Build68 remains an important historical qualified milestone for qualification-plane reuse; Build66 remains the historical ordered-parallel decode milestone; Build64 remains the historical deep-recovery semantic checkpoint.

## Why Build71 exists

Build69 and Build70 showed that the post-Build68 geometry bottleneck was not broadly distributed work. A single seed dominated the seed-parallel critical path, and within that branch the exact fourth generation (`single4 + pair4 + continuation4`) contributed 56.5% of dominant-seed evaluations on `B/mild` and 73.4% on `B/angle`, fed by 593 and 1363 independent `sibling3` inputs respectively.

Build71 therefore changes **scheduling only**:

```text
24 selected seeds
      ↓
bounded parallel prefix through sibling3
      ↓
freeze every sibling3 input in exact seed/traversal order
      ↓
one bounded global generation-four worker pool
      ↓
single4 → pair4 → continuation4
      ↓
store each result by frozen task index
      ↓
strict original-index concatenation
      ↓
complete bank FREEZE
      ↓
unchanged Build68 qualification plane reuse
      ↓
unchanged Build66 ordered protected-data decode
      ↓
HMAC final authentication only
```

The design deliberately avoids nested pools. CPUs released by short prefixes can assist the expensive fourth-generation region without changing the qualified logical traversal or final bank order.

## Frozen invariants

Build71 preserves exactly:

- Build47 freeze and Build48 seed selection/order;
- 24 selected seeds on retained deep-recovery cases;
- every Build64 proposal-only operation, bound and keep count;
- total geometry-evaluation counts;
- final bank contents and ordering;
- Build68 full-pilot qualification, thresholds and pixel-plane reuse;
- qualification counts/order;
- Build66 candidate-stable decode batching;
- logical decode-candidate and list-frame counts;
- payload/HMAC outcomes and first-logical-HMAC semantics;
- control rejection and `B/angle` rejection.

Secret key, payload, ECC and HMAC never create, rank, stop or schedule geometry.

## Exact semantic result

Both qualification runs preserve the retained matrix exactly. The deep cases remain:

| case | geometry evals | bank | qualified | logical decode | logical frames | result |
|---|---:|---:|---:|---:|---:|---|
| `control/front` | 18,021 | 26 | 0 | 0 | 0 | REJECT |
| `control/mild` | 117,609 | 810 | 6 | 6 | 18,432 | REJECT |
| `control/angle` | 21,203 | 11 | 0 | 0 | 0 | REJECT |
| `B/mild` | 79,259 | 937 | 935 | 691 | 2,120,047 | HMAC/PAYLOAD PASS |
| `B/angle` | 334,857 | 6,198 | 0 | 0 | 0 | REJECT |

Historical positive cases `A/front`, `A/mild`, `A/angle` and `B/front` remain authenticated through their earlier qualified paths.

## Two-run performance evidence

The same Build68 timing baseline retained by the gate is used for both Build71 runs.

| metric | Build68 | Build71 run 1 | Build71 run 2 | two-run mean |
|---|---:|---:|---:|---:|
| `B/mild` geometry | 242,405 ms | 150,469 ms | 153,775 ms | 152,122 ms |
| `B/mild` command | 362,286 ms | 277,693 ms | 282,958 ms | 280,325.5 ms |
| `B/angle` geometry | 839,049 ms | 352,838 ms | 367,572 ms | 360,205 ms |
| `B/angle` command | 857,452 ms | 396,068 ms | 412,755 ms | 404,411.5 ms |
| full 9-photo matrix | 1,902,847 ms | 1,233,478 ms | 1,263,623 ms | 1,248,550.5 ms |

Two-run mean improvement versus Build68:

- full matrix: about **1.524x**, **-34.4%**;
- `B/mild` command: about **1.292x**, **-22.6%**;
- `B/angle` command: about **2.120x**, **-52.8%**.

The fourth-generation worker pool scales close to the eight-worker limit in the first run: `B/mild` compresses 147,565 ms cumulative worker time into 18,586 ms wall time (~7.94x), while `B/angle` compresses 918,031 ms into 115,083 ms (~7.98x). The second run reproduces the same scheduling behavior and overall speedup.

## Regression guarantees

The deterministic regression reconstructs a full seed as:

```text
Build71 prefix + ordered Build71 generation-four subtrees
```

and compares it directly with `experimentalV4PhoneBuild64SeedBank` for evaluation count, bank size, hypothesis order, quad, homography and proposal score. Race detection is also required because workers share a read-only `pixelPlane`.

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

## Decision

**Build71 is promoted as the current qualified smartphone baseline.** Future performance work must start from Build71 semantics. The next measured bottleneck is the geometry prefix through `sibling3`, especially on `B/angle`; any Build72+ work must remain separate until exact equivalence is proven.
