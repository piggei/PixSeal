# Format-v4 Build85 — qualified Build84 sampler post-promotion profiling

Build85 is an **observability-only** checkpoint over the qualified Build84 smartphone baseline. It does not change the active recovery algorithm, the projective reader, geometry search, ordering, thresholds, qualification, protected-data decode, ECC/Hamming, whitening/HMAC domains, or Format-v4 wire format.

Its purpose is to measure the **promoted Build84 reader itself** before selecting another performance candidate. Build83 profiled the historical Build76 reader and led to the successful Build84 RGB-fetch change; after promotion, the next optimization must be chosen from fresh measurements of the new baseline rather than from stale pre-Build84 percentages.

## Why re-profile after Build84

Build84 improved the isolated reader by 5.81–7.37% on the qualified Go 1.26.0 host and produced a repeatable 1.35% full-matrix physical improvement. That means the cost distribution inside the reader has changed. The previous Build83 pprof percentages are no longer authoritative for the next candidate.

Build85 therefore asks one question only:

> where does CPU time go inside the **qualified Build84 projective block reader** now?

No Build85 physical qualification is required because no runtime algorithm changes.

## Public deterministic workload

Build85 reuses the public deterministic Build83 1632x1632 RGB plane and fixed front/mild/angle-like homographies. The benchmark workload is independent of:

- the private smartphone corpus;
- secret keys or payloads;
- ECC/Hamming results;
- HMAC outcomes;
- oracle/reference geometry;
- knowledge of the expected message.

The fixture first proves that Build84 remains bit-exact with the historical reader on all benchmark origins and that the prepared sample representation used by stage benchmarks reproduces `samplePlaneLuminance` bit-for-bit.

## Measurements

Build85 records repeated `testing.B` measurements for:

1. the complete **qualified Build84 block reader** on front/mild/angle-like homographies;
2. 64 `homography.mapPoint` calls per benchmark operation;
3. bounds/floor/clamp plus row/x byte-address formation for 64 already-projected samples;
4. exact RGB slice fetch + four luminance conversions + bilinear interpolation for 64 prepared samples;
5. DCT accumulation over one 8x8 block.

The stage measurements are diagnostic decompositions. They are not assumed to sum exactly to the integrated reader because compiler optimization, instruction scheduling and cache locality differ between isolated loops and the complete function.

## CPU profile and compiler evidence

The profile target also runs a longer CPU pprof on `BenchmarkExperimentalV4Build85QualifiedBlockReaderAngle` and preserves:

- `build85-pprof-top.txt`;
- `build85-pprof-reader-list.txt`, line-level pprof evidence for `experimentalV4PhoneBuild84ReadProjectiveBlockValue`;
- `build85-compiler-reader.txt`, Go 1.26.0 inlining/BCE evidence for the qualified reader.

The next performance candidate must be selected from this evidence. Build85 itself is never promotable from timing because it changes no runtime algorithm.

### Qualified-host result

The Go 1.26.0 i3-13100T profile completed successfully. Whole-reader means were **1384.6 / 1462.4 / 1441.2 ns/block** for front/mild/angle-like fixtures and the long angle profile run was **1427 ns/block**. The integrated 24.99 s CPU profile attributed **3.55 s + 1.11 s = 4.66 s (~18.6%)** to the two DCT accumulation lines; `homography.mapPoint` was **5.80% cumulative**. Compiler BCE diagnostics showed four remaining `IsInBounds` checks across the two cosine-table accumulation lines. This evidence selects Build86: exact DCT table-value hoisting only, with no floating-point reassociation.

## Commands

Source/fixture gate:

```bash
make clean
make v4-build85-qualified-sampler-test
```

Qualified-host profile:

```bash
make v4-build85-qualified-sampler-profile
```

The profile is staged atomically under:

```text
v4-phone private/build85-profile/
```

with metadata, fixture result, repeated benchmark output, CPU profile, pprof top/list and compiler evidence.

## Decision rule

After the qualified Go 1.26.0 profile:

- select **one** hot component only if flat/cumulative and line-level evidence identify a meaningful exact optimization opportunity;
- do not reopen Build80/81 cache/LUT ideas without new evidence;
- do not optimize `mapPoint`, DCT or another stage merely because it was hot before Build84;
- any Build86 performance candidate must derive from Build84 or a path proven Build84-equivalent and must pass block/FoldScore/continuation/bank exactness before any physical timing;
- superseding Build84 still requires two independent Go 1.26.0 9/9 physical PASS runs with a repeatable material improvement.

Build85 evidence selected Build86 as a benchmark-first exact DCT table-hoist experiment. Build86 subsequently passed exactness but regressed the qualified-host whole-reader and isolated DCT benchmarks, so it was closed without a physical gate. **Build84 remains the current qualified smartphone baseline; Build87 now profiles the remaining qualified RGB/luminance/bilinear path.**
