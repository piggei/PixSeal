# Format-v4 Build84 — qualified exact RGB fetch / bounds-check baseline

Build84 is the **current qualified smartphone baseline**, superseding Build76 after a public Go 1.26.0 benchmark/BCE gate and two independent retained physical 9/9 semantic-equivalence runs. The change is intentionally narrow: only the RGB fetch/address/bounds-check shape inside continuation4 FoldScore projective block reads changes. Geometry search, candidate ordering, qualification, protected-data decode and authentication semantics remain those qualified by Build76.

## Evidence from Build83

Build83 ran on the qualified Go 1.26.0 host (`13th Gen Intel(R) Core(TM) i3-13100T`). The historical angle-like projective reader measured about 1.68–1.77 us/block across repeated runs, with a separate long profile run at 1.710 us/block. The closed Build82 inline comparator measured about 1.66–1.68 us/block, showing that helper-call removal alone was only a small microbenchmark improvement and did not produce a repeatable physical gain.

The stage benchmark measured roughly 243 ns/block for 64 `mapPoint` calls, about 1.09–1.14 us/block for 64 `samplePlaneLuminance` calls, and about 72 ns/block for DCT accumulation. CPU `pprof` attributed approximately 72.33% cumulative CPU to `samplePlaneLuminance`, only 5.65% cumulative to `homography.mapPoint`, and 21.15% flat to the surrounding historical `readProjectiveBlockValue` loop. Build84 therefore targeted RGB fetch/memory/bounds-check structure rather than `mapPoint` or DCT.

## Qualified implementation change

Build84 starts from the exact Build82 continuation4 reader arithmetic but changes only integer address formation and the bounds-check shape:

1. compute `rowStride := width*3` once;
2. compute the two row bases once per bilinear sample;
3. compute `x0*3` and `x1*3` once;
4. form four exact three-byte slices for the `(x0,y0)`, `(x1,y0)`, `(x0,y1)`, `(x1,y1)` pixels;
5. evaluate the historical `.299*R + .587*G + .114*B - 128` formula in the same order from `[0]`, `[1]`, `[2]` of each three-byte slice.

The slice form changes the compiler-visible address/bounds-check structure without `unsafe` and without changing floating-point arithmetic. Go 1.26.0 still reports `IsSliceInBounds` in this function, so Build84 is **not** documented as complete bounds-check elimination; the measured result is a cheaper exact RGB fetch/address shape on the qualified toolchain.

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

The pre-timing source gate is:

```bash
make clean
make v4-build84-phone-rgb-fetch-test
```

It proves:

- projective block `math.Float64bits` equality against the historical reader over multiple homographies, interior blocks, borders and failure cases;
- FoldScore bit equality against Build41;
- Build55 continuation state-for-state and eval-for-eval equality;
- complete Build84-vs-Build76 blind-bank equality in count, order, quad, homography, proposal and validation;
- compatibility of Build76 public telemetry plus Build84-specific block-read counters;
- CLI/buildinfo consistency and `go vet ./...`.

The qualified-host source gate passed under Go 1.26.0 before benchmark and physical testing.

## Qualified Go 1.26.0 benchmark/BCE result

The public deterministic benchmark command is:

```bash
make v4-build84-phone-rgb-fetch-benchmark
```

It compares the historical Build76 reader, the closed Build82 inline reader and Build84 on the same public deterministic front/mild/angle fixture. Five-run means on the qualified i3-13100T host were:

| reader fixture | historical Build76 | Build84 | improvement |
|---|---:|---:|---:|
| front | 1500.4 ns/block | 1413.2 | 5.81% |
| mild | 1514.0 ns/block | 1408.0 | 7.00% |
| angle | 1522.2 ns/block | 1410.0 | 7.37% |

The angle-like Build82 comparator averaged 1533.8 ns/block, so Build84 was also about 8.1% faster than that closed helper-inline experiment in the same benchmark session.

Compiler `check_bce` evidence showed that Build82 retained many direct `IsInBounds` checks and Build84 still retained `IsSliceInBounds` checks for the four exact pixel slices. Therefore the benchmark result is treated as empirical compiler/layout evidence, not as proof that all bounds checks disappeared.

All benchmarked readers reported 0 B/op and 0 allocs/op.

## Physical qualification

Two independent retained physical runs were executed with Go 1.26.0 using separate diagnostic directories. Both passed the complete nine-photo semantic gate:

- controls: 3/3 REJECT;
- A/front: PASS through the historical Build43 path;
- A/mild: PASS through an accepted historical path;
- A/angle: direct PASS;
- B/front: direct PASS;
- B/mild: PASS + HMAC/payload through the deep path;
- B/angle: REJECT.

The deep invariants were identical in both runs:

- B/mild: `79259` geometry evals / bank `937` / qualified `935` / logical decode `691` / logical frames `2120047`, HMAC and payload PASS;
- B/angle: `334857` geometry evals / bank `6198` / qualified `0`, REJECT;
- continuation4 FoldScore/block-read counters matched the exact workload established by Build81/82.

### Timing

| metric | Build76 mean | Build84 run 1 | Build84 run 2 | Build84 mean | improvement |
|---|---:|---:|---:|---:|---:|
| full matrix | 758,102 ms | 748,902 | 746,902 | 747,902 | 1.35% |
| B/mild elapsed | 188,750.5 | 186,961 | 183,780 | 185,370.5 | 1.79% |
| B/mild geometry | 60,351.5 | 59,360 | 58,987 | 59,173.5 | 1.95% |
| B/mild gen4 | 18,473 | 18,181 | 17,869 | 18,025 | 2.43% |
| B/angle elapsed | 226,813 | 225,207 | 222,634 | 223,920.5 | 1.28% |
| B/angle geometry | 183,806.5 | 182,485 | 181,789 | 182,137 | 0.91% |
| B/angle gen4 | 110,829.5 | 109,647 | 109,745 | 109,696 | 1.02% |

The end-to-end gain is intentionally smaller than the isolated reader gain because Build84 changes only continuation4 FoldScore block reads. Both physical runs are faster than the best retained Build76 full-matrix run and preserve every qualified semantic counter.

## Promotion decision

**PROMOTED. Build84 is the current qualified smartphone baseline.**

The promotion is based on correctness first and performance second:

1. exact block/FoldScore/continuation/bank equivalence passed;
2. the public Go 1.26.0 reader benchmark showed a stable 5.81–7.37% improvement in the modified component;
3. two independent physical runs passed 9/9 with exact workload;
4. both physical runs improved the full matrix and all tracked deep timing metrics relative to the Build76 qualified mean;
5. the implementation adds no `unsafe`, cache, LUT, geometry heuristic or protected-data dependency.

Build76 remains the previous qualified baseline and historical scheduler/runtime reference. Build80 is rejected; Build81 and Build82 remain semantic-valid but non-promoted experiments; Build83 remains observability-only.

## Baseline artifacts and future candidates

Source-controlled reference timing files are kept under:

```text
docs/qualified-baselines/build76-phone-performance.tsv
docs/qualified-baselines/build84-phone-performance.tsv
```

They contain only public case names and qualified timing summaries; no private images, keys, payloads or diagnostic logs are shipped. The Build84 physical harness uses the source-controlled Build76 reference for historical comparison rather than relying on the mutable/private Build76 diagnostic directory whose TSV was later truncated by an interrupted invocation.

Any future smartphone performance candidate must derive from Build84 or a path proven Build84-equivalent, preserve all security/semantic invariants, use measurement-first selection, and pass two independent Go 1.26.0 9/9 physical runs with a repeatable material improvement before it can supersede Build84.
