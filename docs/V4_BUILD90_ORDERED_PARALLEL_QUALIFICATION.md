# Format-v4 Build90 — ordered-parallel qualification

Build90 is a **performance candidate over the current qualified Build84 smartphone baseline**. It is the first optimization selected from the full-pipeline Build89 retained-corpus profile rather than from the projective-reader microbenchmark line.

## Why Build90

Build89 ran the unchanged qualified Build84 runtime over the retained nine-photo strength-48 smartphone matrix and passed all semantic gates. Its deep-path timing split showed that the dominant stage now depends on the difficult case:

- `B/mild`: deep recovery 149,504 ms; geometry 58,976 ms; **qualification 61,193 ms**; ordered decode 29,119 ms. Qualification is 40.93% of deep wall time and is the largest single stage.
- `B/angle`: deep recovery 193,441 ms; geometry 186,352 ms; qualification only 6,867 ms. Generation 4 remains dominant inside geometry on this reject case.

The `B/mild` frozen bank contains 937 candidates and 935 qualify. Build84 evaluates those candidates serially even though the bank is already frozen and each qualification uses only the shared read-only pixel plane, public pilot structure and the candidate's own geometry. This creates a high-value scheduling opportunity without changing any numerical computation or security boundary.

Build89 itself remains observability-only and is not a baseline. Build84 remains the current qualified smartphone baseline.

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

## Deferred physical gate

Only after the Go 1.26.0 benchmark is reviewed should the retained private matrix be run:

```bash
make v4-build90-phone-physical-test
```

The physical gate keeps the complete nine-case semantic matrix, exact historical deep counters, Build84 continuation4 workload counters, HMAC/payload expectations and the Build84 timing reference. It additionally requires the qualification worker/task telemetry to be structurally valid.

Promotion still requires **two independent physical PASS runs** with identical semantics and a material repeatable speedup. Build84 remains qualified until that evidence exists.

## Build89 retained evidence used to select Build90

| case | deep total ms | geometry ms | qualification ms | decode ms | bank | qualified |
|---|---:|---:|---:|---:|---:|---:|
| B/mild | 149,504 | 58,976 | **61,193** | 29,119 | 937 | 935 |
| B/angle | 193,441 | 186,352 | 6,867 | 0 | 6,198 | 0 |

The complete Build89 command matrix was 745,821 ms and preserved the expected 9/9 behavior. These timings are profiling evidence only; they do not change the Build84 qualification status.
