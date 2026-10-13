# Format-v4 Build91 — contemporaneous Build84 / Build90 paired physical profiling

Build91 is a **completed observability-only decision build**. It adds no recovery algorithm. Its contemporaneous AB/BA evidence resolved the cross-day timing ambiguity and **promoted the unchanged Build90 implementation as the current qualified smartphone baseline**. Build84 is now the previous qualified baseline; Build91 itself is not a runtime baseline.

## Why Build91 exists

Build90 passed its public Go 1.26.0 benchmark and race gate decisively:

- high-pass qualification: 207.413 ms serial mean versus 48.429 ms parallel mean, **4.283x / 76.65% wall-time reduction**;
- early-reject qualification: 111.935 ms serial mean versus 28.120 ms parallel mean, **3.981x / 74.88% reduction**;
- targeted shared-`pixelPlane` race test: PASS.

Three retained physical runs then passed the complete 9/9 semantic gate. The qualification stage itself was strongly repeatable in the two non-degraded runs:

| run | B/mild qualification | B/angle qualification | B/mild elapsed | B/angle elapsed |
|---|---:|---:|---:|---:|
| Build90 run 1 | 14,013 ms | 1,515 ms | 155,665 ms | 245,758 ms |
| Build90 run 2 | 27,985 ms | 3,336 ms | 231,582 ms | 533,422 ms |
| Build90 run 3 | 14,314 ms | 1,569 ms | 155,578 ms | 262,345 ms |

Run 2 suffered a broad host slowdown: unchanged freeze/gen3/gen4/decode work slowed by roughly the same factor as qualification, so it cannot isolate a Build90-specific regression. Runs 1 and 3 reproduce the new qualification speed strongly, but their unchanged geometry timings are also slower than the older Build84 qualification runs. Comparing an 8 October candidate against an older baseline therefore no longer cleanly separates code effect from host drift.

Build91 removes that ambiguity with **contemporaneous paired A/B measurement**.

Harness note: before invoking `go test ./watermark`, the shell wrapper canonicalizes `V4_PHONE_ACQUISITION_DIR` to an absolute path. Go package tests execute with the package directory as their working directory, so forwarding the repository-relative `v4-phone private/build38-acquired` path would otherwise resolve below `./watermark` and fail even when the retained corpus exists. This correction affects only test-path resolution, not recovery code or timing semantics.

## Paired design

Only the two retained deep cases are needed:

- `phone-b-mild.jpg` — authenticated deep positive;
- `phone-b-angle.jpg` — deep negative/reject.

For each pair, the JPEG is decoded once with PixSeal's qualified project-owned `internal/jpeglegacy` decoder and the historical pre-deep Build41 work image/boundary are prepared once. Build84 and Build90 then receive the **same in-memory `work` image and `boundary`**. Each recovery still independently regenerates its complete deep geometry bank, qualification and protected-data decode exactly as its historical implementation specifies.

The order is reversed to reduce systematic drift:

1. cycle 1: B/mild, Build84 -> Build90;
2. cycle 1: B/angle, Build84 -> Build90;
3. cycle 2: B/angle, Build90 -> Build84;
4. cycle 2: B/mild, Build90 -> Build84.

A short configurable idle interval (default 5 s) is inserted before each deep run. Timing is **evidence only**; no timing threshold is part of the semantic gate.

## What is compared

Each run records:

- deep total wall time and internally measured total;
- geometry total;
- Build47 freeze;
- prefix1, gen2, gen3 and gen4 wall time;
- post-freeze plane preparation;
- qualification wall time;
- ordered decode wall time;
- geometry/bank/qualification/decode/list-frame counts;
- physical/speculative decode counts;
- Build84 continuation4 fold/read counters;
- authentication result;
- qualification workers/tasks.

Build84 qualification is the unchanged serial loop (`workers=1`). Build90 uses the unchanged ordered-parallel worker pool selected by `GOMAXPROCS` and commits results in original bank order.

## Semantic gate

The Go test enforces the retained physical invariants independently for both implementations.

### B/mild

- geometry evaluations: `79,259`;
- frozen bank: `937`;
- qualification evaluations: `1,872`;
- qualified candidates: `935`;
- logical decode candidates: `691`;
- logical list frames: `2,120,047`;
- physical decode candidates: `696`;
- speculative candidates: `5`;
- continuation4 FoldScores: `16,496`;
- projective block reads: `21,114,880`, all successful;
- exact expected payload and authentication: PASS.

### B/angle

- geometry evaluations: `334,857`;
- frozen bank: `6,198`;
- qualification evaluations: `6,198`;
- qualified candidates: `0`;
- continuation4 FoldScores: `118,004`;
- projective block reads: `151,045,120` (`149,370,664` successful / `1,674,456` failed);
- authentication: REJECT.

Within every pair Build84 and Build90 must also agree on all semantic telemetry, payload/outcome and authenticated extract information. Only scheduling/timing may differ.

## Commands

Source/race gate:

```bash
make clean
make v4-build91-paired-test
```

Contemporaneous retained-corpus profile:

```bash
make v4-build91-paired-physical-profile
```

The profile writes exactly:

```text
v4-phone private/build91-paired-profile/
├── build91-metadata.txt
├── build91-paired-raw.txt
├── build91-paired.tsv
└── build91-paired.md
```

To archive it, use the actual output directory:

```bash
zip -r build91-paired-profile.zip \
  'v4-phone private/build91-paired-profile'
```

## Completed paired result and promotion decision

The qualified Go 1.26.0 host completed both cycles and the paired semantic gate passed.

| case | cycle/order | Build84 total | Build90 total | reduction | Build84 geometry | Build90 geometry | Build84 qualification | Build90 qualification | q speedup | Build84 decode | Build90 decode |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| B/mild | 1 / 84->90 | 165,717 | 116,044 | 29.97% | 67,159 | 67,655 | 64,274 | 14,816 | 4.338x | 34,055 | 33,314 |
| B/mild | 2 / 90->84 | 171,172 | 120,164 | 29.80% | 71,932 | 70,357 | 64,097 | 15,134 | 4.235x | 34,876 | 34,429 |
| B/angle | 1 / 84->90 | 218,567 | 211,819 | 3.09% | 210,926 | 210,022 | 7,404 | 1,562 | 4.740x | 0 | 0 |
| B/angle | 2 / 90->84 | 223,173 | 212,232 | 4.90% | 215,016 | 210,429 | 7,801 | 1,563 | 4.991x | 0 | 0 |

Paired means are B/mild 168,444.5 -> 118,104 ms (**-29.89%**) with qualification 64,185.5 -> 14,975 ms (**4.286x**), and B/angle 220,870 -> 212,025.5 ms (**-4.00%**) with qualification 7,602.5 -> 1,562.5 ms (**4.866x**). Across all four paired deep executions Build84 totals 778,629 ms and Build90 660,259 ms (**-15.20%**).

The AB/BA reversal reproduces the speedup while unchanged geometry and decode remain comparable within contemporaneous pairs. The prior cross-day geometry slowdown is therefore attributed to host drift rather than to ordered-parallel qualification. Together with the public benchmark/race PASS and three independent full nine-photo semantic PASS runs, this evidence satisfies the promotion discipline.

**Final decision: Build90 is promoted as the current qualified smartphone baseline. Build84 becomes the previous qualified baseline. Build91 remains observability-only and is closed as paired profiling PASS.**

The exact authoritative paired rows are frozen source-side as `docs/qualified-baselines/build91-build84-build90-paired.tsv`; summarized Build90 qualification/performance evidence is frozen in `docs/qualified-baselines/build90-phone-performance.tsv`.

No key, payload, ECC, HMAC result or known message is used to select or schedule geometry or qualification work. HMAC remains final authentication only.
