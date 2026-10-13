# Format-v4 Build90 — ordered-parallel qualification

Build90 is the **current qualified smartphone baseline**. It is the scheduling-only performance successor selected from the full-pipeline Build89 retained-corpus profile and promoted after public benchmark/race evidence, three independent full physical semantic PASS runs, and the contemporaneous Build91 AB/BA paired physical decision.

## Why Build90

Build89 ran the unchanged qualified Build84 runtime over the retained nine-photo strength-48 smartphone matrix and passed all semantic gates. Its deep-path timing split showed that the dominant stage now depends on the difficult case:

- `B/mild`: deep recovery 149,504 ms; geometry 58,976 ms; **qualification 61,193 ms**; ordered decode 29,119 ms. Qualification is 40.93% of deep wall time and is the largest single stage.
- `B/angle`: deep recovery 193,441 ms; geometry 186,352 ms; qualification only 6,867 ms. Generation 4 remains dominant inside geometry on this reject case.

The `B/mild` frozen bank contains 937 candidates and 935 qualify. Build84 evaluates those candidates serially even though the bank is already frozen and each qualification uses only the shared read-only pixel plane, public pilot structure and the candidate's own geometry. This creates a high-value scheduling opportunity without changing any numerical computation or security boundary.

Build89 itself remains observability-only and is not a baseline. Build84 was the qualified baseline when Build90 was selected and is now the previous qualified baseline.

## Exact change

Build90 keeps Build84 geometry and protected-data decode unchanged. Only the qualification loop is scheduled differently:

1. the complete Build84 geometry bank is generated and frozen exactly as before;
2. a bounded worker pool, capped by `runtime.GOMAXPROCS(0)` and bank length, evaluates `experimentalV4PhoneBuild68Qualify` once for each frozen bank entry;
3. every result is stored at its original bank index;
4. after all workers finish, evaluation counts and accepted hypotheses are committed **strictly in original bank order**;
5. the existing ordered-parallel protected-data decode then runs unchanged.

The authoritative qualification function remains `experimentalV4PhoneBuild68Qualify`. Build90 does not change:

- proposal or held-out validation arithmetic;
- full public-pilot detector arithmetic or thresholds;
- bank contents or bank order;
- geometry proposal/ranking/refinement;
- Build84 continuation4 RGB sampler;
- protected-data candidate order;
- list-frame order or first-logical-HMAC semantics;
- ECC/Hamming, whitening/HMAC domains, strength 48, wire format or encoder.

No key, payload, ECC result, HMAC result or known message content is available to qualification scheduling. The complete geometry bank is frozen before Build90 starts qualification.

## Concurrency model

The shared `pixelPlane` is immutable after construction. `experimentalV4PhoneBuild68Qualify` allocates its detector accumulators on the worker stack and reads only the public pilot candidate, the candidate homography and the pixel plane. Workers never append to the authoritative qualified slice.

The result array is pre-sized to bank length. Worker `i` writes only `results[i]`. After `WaitGroup.Wait`, one serial ordered-commit pass reconstructs the exact Build84 qualification evaluation count and qualified candidate sequence.

Low-overhead public telemetry adds only:

- `Build90QualificationWorkers`;
- `Build90QualificationTasks`;
- `Build90QualificationEvaluations`.

No per-candidate timers are added to the physical candidate path.

## Exactness and race gates

Before any private physical test, Build90 must prove equivalence on public deterministic fixtures:

- candidate-by-candidate `evals` identical to serial Build84 qualification;
- candidate-by-candidate `ok` identical;
- hypothesis geometry/proposal/validation/detection fields bit-identical;
- ordered committed qualified slice identical;
- total qualification evaluations identical;
- targeted `go test -race` clean while all workers share one read-only `pixelPlane`.

The normal source target also re-runs the historical Build47/68/71/73/75/76 and Build82-90 regression chain, CLI help, version/buildinfo consistency and `go vet`.

## Public benchmark gate

Build90 provides two public deterministic pilot-only workloads. Neither uses a private acquisition, secret key, payload, ECC outcome or HMAC result.

### High-pass workload

Twenty-four identical-but-index-tagged frozen hypotheses use a correct public synthetic homography and reach the full public-pilot detector. This models the expensive shape of `B/mild`, where almost the entire bank passes the preliminary gates.

Compare:

- `BenchmarkExperimentalV4Build90SerialHighPass`
- `BenchmarkExperimentalV4Build90ParallelHighPass`

The candidate should show a large, stable wall-time reduction before a private physical run is justified.

### Early-reject workload

Five hundred twelve hypotheses deliberately fail the proposal gate after the unchanged held-out validation read. This models the short-task side of the scheduling trade-off and checks that channel/goroutine overhead does not create a material regression.

Compare:

- `BenchmarkExperimentalV4Build90SerialEarlyReject`
- `BenchmarkExperimentalV4Build90ParallelEarlyReject`

A strong high-pass win with only small/acceptable early-reject overhead is the required direction. Benchmark timing is diagnostic and never a semantic gate.

Run:

```bash
make clean
make v4-build90-qualification-parallel-test
make v4-build90-qualification-benchmark
```

The benchmark writes:

```text
v4-phone private/build90-benchmark/
├── build90-metadata.txt
├── build90-fixture-test.txt
├── build90-benchmark.txt
└── build90-race-test.txt
```

## Physical qualification gate — completed

After the Go 1.26.0 benchmark/race gate passed, the retained private matrix was run:

```bash
make v4-build90-phone-physical-test
```

The physical gate keeps the complete nine-case semantic matrix, exact historical deep counters, Build84 continuation4 workload counters, HMAC/payload expectations and the Build84 timing reference. It additionally requires the qualification worker/task telemetry to be structurally valid.

Build90 passed the complete retained 9/9 semantic matrix three times. Because cross-day host drift contaminated absolute timing, the promotion decision was intentionally deferred to Build91 contemporaneous AB/BA evidence rather than selecting favorable runs.

## Build89 retained evidence used to select Build90

| case | deep total ms | geometry ms | qualification ms | decode ms | bank | qualified |
|---|---:|---:|---:|---:|---:|---:|
| B/mild | 149,504 | 58,976 | **61,193** | 29,119 | 937 | 935 |
| B/angle | 193,441 | 186,352 | 6,867 | 0 | 6,198 | 0 |

The complete Build89 command matrix was 745,821 ms and preserved the expected 9/9 behavior. These timings were profiling evidence only and selected Build90; Build84 remained qualified at that stage.

## Qualified-host benchmark and physical evidence — 2026-10-08

The public Go 1.26.0 benchmark/race gate completed successfully on the i3-13100T qualification host.

| workload | serial mean | parallel mean | speedup | wall-time reduction |
|---|---:|---:|---:|---:|
| high-pass | 207.413 ms | 48.429 ms | **4.283x** | **76.65%** |
| early-reject | 111.935 ms | 28.120 ms | **3.981x** | **74.88%** |

The targeted race test passed. Build90 then passed the complete retained 9/9 physical semantic gate three times. Deep semantic counters were exact in every run.

| physical run | full matrix | B/mild elapsed | B/mild geometry | B/mild qualification | B/angle elapsed | B/angle geometry | B/angle qualification |
|---|---:|---:|---:|---:|---:|---:|---:|
| run 1 | 741,008 ms | 155,665 | 68,973 | **14,013** | 245,758 | 205,178 | **1,515** |
| run 2 | 1,103,068 ms | 231,582 | 93,577 | 27,985 | 533,422 | 440,043 | 3,336 |
| run 3 | 779,553 ms | 155,578 | 66,989 | **14,314** | 262,345 | 217,206 | **1,569** |

Run 2 is a broad host-degraded execution: unchanged geometry/decode work slows by approximately the same factor as qualification. Runs 1 and 3 strongly reproduce the Build90 qualification gain and nearly identical B/mild end-to-end elapsed, but unchanged geometry timing remains slower than the older Build84 qualification reference. Because promotion requires a repeatable code-specific speedup rather than selecting favorable cross-day runs, the decision moved to Build91 contemporaneous AB/BA profiling. That profile passed semantically and produced the following paired means:

| case | Build84 deep total | Build90 deep total | total reduction | Build84 qualification | Build90 qualification | qualification speedup |
|---|---:|---:|---:|---:|---:|---:|
| B/mild | 168,444.5 ms | **118,104 ms** | **29.89%** | 64,185.5 | **14,975** | **4.286x** |
| B/angle | 220,870 ms | **212,025.5 ms** | **4.00%** | 7,602.5 | **1,562.5** | **4.866x** |

Geometry/decode remained comparable within the same pairs and the order reversal reproduced the result. **Decision: PROMOTED. Build90 is the current qualified smartphone baseline; Build84 is the previous qualified baseline.** Build91 itself remains observability-only evidence. See [`V4_BUILD91_PAIRED_PHYSICAL_AB.md`](V4_BUILD91_PAIRED_PHYSICAL_AB.md) and the frozen TSVs under `docs/qualified-baselines/`.
