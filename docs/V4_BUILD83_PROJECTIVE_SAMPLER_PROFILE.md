# Format-v4 Build83 — projective sampler CPU profiling

Build83 is an **observability-only** checkpoint after Build82 closed as exact but non-promoted. The active smartphone deep fallback is restored to the qualified Build76 implementation. Build83 changes no recovery algorithm, geometry order, score, threshold, protected-data decode, ECC/Hamming, whitening/HMAC domain, encoder, wire format, pilot or strength.

## Why Build83 exists

Build82 answered the helper-call hypothesis cleanly. Go 1.26.0 does not inline `samplePlaneLuminance`, `homography.mapPoint` or `readProjectiveBlockValue`; Build82 manually incorporated only the historical luminance sampler inside continuation4. Two independent physical runs preserved exact 9/9 semantics and the deep workload, but performance was not repeatable. Full-matrix elapsed was 933,537 ms and 729,213 ms (831,375 ms mean) versus the qualified Build76 mean 758,102 ms. B/mild generation four was 24,305 / 19,991 ms versus 18,473 ms Build76 mean; B/angle generation four was 136,040 / 107,359 ms versus 110,830 ms. Build82 is therefore **semantic PASS x2 / non-promoted**.

The correct next step is measurement, not another speculative optimization.

## Runtime rule

`v4-extract-phone` uses the qualified Build76 deep recovery again. Build82 code remains in the tree only for archived exactness tests and as a controlled comparator in Build83 microbenchmarks. Build83 profiling cannot influence geometry generation, ranking, pruning, qualification or authentication.

## Public deterministic benchmark fixture

Build83 adds `experimental_v4_build83_projective_sampler_benchmark_test.go`. The fixture:

- creates a deterministic 1632x1632 RGB plane from a public arithmetic pattern;
- uses fixed public front/mild/angle-like homographies;
- exercises only interior 8x8 blocks so benchmark validity is deterministic;
- verifies historical `readProjectiveBlockValue` and the closed Build82 exact-inline reader agree bit-for-bit on all fixture blocks;
- contains no private acquisition, key, payload, ECC result or HMAC selector.

The benchmark set measures:

- historical block reader, front/mild/angle-like geometry;
- closed Build82 inline block reader on the same angle-like workload;
- `homography.mapPoint`;
- `samplePlaneLuminance`;
- DCT accumulation over precomputed luminance values.

The mapPoint, luminance-sampling and DCT stage microbenchmarks define one benchmark operation as one complete 8x8 block (64 points/samples/accumulations), so their `ns/op` values are directly comparable at block granularity. The microbenchmarks are intentionally diagnostic. Absolute nanoseconds are host-dependent; compare repeated results on the same qualified host and use CPU-profile flat/cumulative attribution before choosing another implementation candidate.

## Commands

Run the source/regression gate first:

```bash
make clean
make v4-build83-projective-sampler-test
```

Then run the benchmark/profile pack on the qualified Go 1.26.0 host:

```bash
make v4-build83-projective-sampler-profile
```

Optional controls:

```bash
make V4_BUILD83_BENCHTIME=3s V4_BUILD83_COUNT=7 V4_BUILD83_PPROF_TIME=30s \
  v4-build83-projective-sampler-profile
```

Default output:

```text
v4-phone private/build83-profile/
  build83-metadata.txt
  build83-fixture-test.txt
  build83-benchmark.txt
  build83-profile-run.txt
  build83-cpu.pprof
  build83-pprof-top.txt
  build83-inlining.txt
  watermark.test          # retained when Go emits the profile binary
```

The output is staged atomically and published only after all commands succeed.

## Decision rule for Build84

Do not choose Build84 from microbenchmark ranking alone. Use Build83 to answer which component dominates the authoritative historical block reader under the controlled workload:

1. If `samplePlaneLuminance` still dominates flat/cumulative CPU despite Build82's mixed physical result, inspect its generated code/memory access before changing arithmetic.
2. If `homography.mapPoint` is material, a **separate** manual-inline candidate may be justified, but it must preserve operation order and pass a new bit-exact block/FoldScore/continuation/bank gate.
3. If DCT accumulation or surrounding reader structure dominates, optimize only that measured component.
4. If no component has a material isolated share or microbenchmarks are unstable, stop micro-optimizing and preserve Build76.

Any later performance candidate must again pass two independent physical semantic gates before promotion.
