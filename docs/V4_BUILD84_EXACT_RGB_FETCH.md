# Format-v4 Build84 — exact RGB fetch / bounds-check reduction

Build84 is an **equivalence-preserving performance candidate** over the qualified Build76 smartphone baseline. It is selected directly from the Build83 public CPU profile and changes only the RGB fetch/bounds-check shape inside continuation4 FoldScore projective block reads. Build76 remains the qualified authority until Build84 proves exact semantics and a material repeatable speedup on two independent physical runs.

## Evidence from Build83

The qualified-host Build83 run used Go 1.26.0 on a 13th Gen Intel Core i3-13100T. The historical angle-like projective reader measured about 1.68–1.77 us/block across repeated runs, with a separate long profile run at 1.710 us/block. The closed Build82 inline comparator measured about 1.66–1.68 us/block in four of five repetitions and 1.677 us in the fifth: only a small microbenchmark improvement, consistent with its non-promotable physical result.

The stage benchmark measured roughly 243 ns/block for 64 `mapPoint` calls, about 1.09–1.14 us/block for 64 `samplePlaneLuminance` calls, and about 72 ns/block for DCT accumulation. More importantly, CPU `pprof` attributed approximately 72.33% cumulative CPU to `samplePlaneLuminance`, only 5.65% cumulative to `homography.mapPoint`, and 21.15% flat to the surrounding historical `readProjectiveBlockValue` loop. Therefore Build84 targets the RGB fetch/memory/bounds-check portion rather than `mapPoint` or DCT.

## Candidate change

Build84 starts from the exact Build82 continuation4 reader but changes only integer address formation and the bounds-check shape:

1. compute `rowStride := width*3` once;
2. compute the two row bases once per bilinear sample;
3. compute `x0*3` and `x1*3` once;
4. form four exact three-byte slices for the `(x0,y0)`, `(x1,y0)`, `(x0,y1)`, `(x1,y1)` pixels;
5. evaluate the historical `.299*R + .587*G + .114*B - 128` formula in the same order from `[0]`, `[1]`, `[2]` of each three-byte slice.

The slice construction makes the three fixed channel indexes inside each length-three pixel slice statically bounded, changing the compiler's bounds-check shape without assuming how many slice checks a particular Go version will emit. Build84 uses no `unsafe` and does not change floating-point arithmetic.

## Explicit non-changes

Build84 does **not** change:

- `homography.mapPoint` or its floating-point operation order;
- `math.Floor`, x1/y1 clamps, `fx/fy`, top/bottom bilinear interpolation, or `c23/c32` accumulation order;
- `float64` arithmetic;
- single4, pair4, generation 1–3, freeze, qualification or protected-data decode;
- geometry candidate order, ordered commit or first-logical-HMAC semantics;
- encoder, Format-v4 wire format, locked public pilot, strength 48, ECC/Hamming, whitening/HMAC domains or qualified thresholds;
- any private-corpus, payload, key or HMAC information as a geometry selector.

There is no LUT, cache, precomputed luminance plane, bounding rectangle, new fallback or coordinate recurrence.

## Correctness gate

Before timing is considered:

```bash
make clean
make v4-build84-phone-rgb-fetch-test
```

The gate proves:

- projective block `math.Float64bits` equality against the historical reader over multiple homographies, interior blocks, borders and failure cases;
- FoldScore bit equality against Build41;
- Build55 continuation state-for-state and eval-for-eval equality;
- complete Build84-vs-Build76 blind-bank equality in count, order, quad, homography, proposal and validation;
- compatibility of Build76 public telemetry plus Build84-specific block-read counters;
- CLI/buildinfo consistency and `go vet ./...`.

## Go 1.26 benchmark/BCE gate

Run on the qualified host before any physical corpus run:

```bash
make v4-build84-phone-rgb-fetch-benchmark
```

Default output:

```text
v4-phone private/build84-benchmark/
  build84-metadata.txt
  build84-fixture-test.txt
  build84-benchmark.txt
  build84-compiler.txt
  build84-bce-inlining.txt
  build84-bce-reader.txt
```

The benchmark compares the historical Build76 reader, the closed Build82 inline reader and Build84 on the same public deterministic front/mild/angle fixture. `build84-bce-reader.txt` extracts compiler bounds-check evidence for the Build82 and Build84 reader functions.

**Performance is not a correctness gate.** A slow benchmark still exits successfully if exactness/compiler commands succeed. Review the repeated Go 1.26.0 numbers before spending time on the physical corpus. As an engineering decision rule, a candidate should show a stable, useful reader-level gain (roughly >=5% is a reasonable threshold) before the physical test is justified; this threshold is not part of authentication or semantic correctness.

## Physical gate

Only after benchmark review:

```bash
make V4_PHONE_BUILD84_DIAGNOSTIC_DIR='v4-phone private/build84-diagnostics1'   v4-build84-phone-physical-test
```

If run 1 is exact and materially faster, repeat independently into `build84-diagnostics2`. Promotion still requires two independent physical 9/9 semantic PASS runs and a material repeatable improvement over Build76. One fast run is insufficient.

The deep invariants remain the Build76 reference values, including B/mild `79259 / 937 / 935 / 691 / 2120047` with HMAC/payload PASS and B/angle `334857 / 6198 / 0` with REJECT. Build84 FoldScore/block-read counters must also match the exact continuation4 workload established by Build81/82.

## Stop rule

If the Go 1.26.0 microbenchmark is neutral/regressive, do **not** run the expensive physical matrix and do not stack another speculative cache/LUT on top. Close Build84 as a negative compiler/layout experiment and preserve Build76. If Build84 is clearly faster in microbenchmark but fails to reproduce physically, close it as non-promoted just like Build82.
