# Format-v4 Build86 — exact DCT table hoist / bounds-check candidate

Build86 is an **equivalence-preserving performance candidate** over the current qualified Build84 smartphone baseline. It is selected directly from the qualified-host Build85 post-promotion profile. Build84 remains qualified unless Build86 passes all exactness gates and then demonstrates a repeatable material improvement in two independent Go 1.26.0 physical runs.

## Evidence from Build85

On the qualified 13th Gen Intel Core i3-13100T / Go 1.26.0 host, Build85 measured the promoted Build84 reader at mean values of:

- front: **1384.6 ns/block**;
- mild: **1462.4 ns/block**;
- angle: **1441.2 ns/block**.

The isolated angle decomposition measured approximately 211.68 ns/block for 64 `mapPoint` calls, 251.2 ns for address/floor/clamp work, 491.68 ns for prepared RGB fetch + luminance + bilinear interpolation, and 64.53 ns for DCT accumulation. The integrated CPU profile is more authoritative for candidate selection: line-level pprof attributed **3.55 s + 1.11 s = 4.66 s**, about **18.6% of total CPU samples**, to the two DCT accumulation statements in the qualified Build84 reader.

Go 1.26.0 BCE diagnostics also reported four remaining cosine-table checks at those lines:

```text
experimental_v4_phone_build84_rgb_fetch.go:152:21: Found IsInBounds
experimental_v4_phone_build84_rgb_fetch.go:152:33: Found IsInBounds
experimental_v4_phone_build84_rgb_fetch.go:153:21: Found IsInBounds
experimental_v4_phone_build84_rgb_fetch.go:153:33: Found IsInBounds
```

That gives Build86 a narrow, measured target without reopening the previously rejected cache/LUT/manual-inline ideas.

## Exact candidate change

Build84 evaluates:

```go
c23 += l * table3[x] * table2[y]
c32 += l * table2[x] * table3[y]
```

Build86 keeps the exact same multiplication and accumulation order but loads the row constants once per `y` and the column constants once per sample:

```go
t2y := table2[y]
t3y := table3[y]
...
t2x := table2[x]
t3x := table3[x]
c23 += l * t3x * t2y
c32 += l * t2x * t3y
```

This is **not** a precomputed DCT-product optimization. In particular, Build86 does not replace `(l * a) * b` with `l * (a * b)`, because that would reassociate floating-point arithmetic and could change the final bit pattern.

Everything else remains the qualified Build84 path:

- same RGB address/slice shape;
- same `homography.mapPoint` calls;
- same `math.Floor`, clamps and fractions;
- same RGB-to-luminance `float64` arithmetic;
- same bilinear interpolation order;
- same `c23`/`c32` update order and final `abs(c23)-abs(c32)`;
- same single4/pair4 and generation 1-3 behavior;
- same ordered scheduling/commit, freeze, qualification and protected-data decode;
- no `unsafe`, LUT, cache, precomputed luminance plane, DCT product table, pruning or new fallback.

## Exactness gates

`make v4-build86-dct-hoist-test` covers the established regression chain plus Build86-specific checks:

1. Build86 block reader vs qualified Build84 `math.Float64bits`, including border/failure origins;
2. Build86 FoldScore vs Build84 bit-for-bit with identical visible count;
3. Build86 continuation4 vs Build84 state-for-state and eval-for-eval;
4. complete Build86 blind bank vs Build84: identical eval count, seed/worker counts, bank length/order, quad, homography, proposal and validation;
5. telemetry compatibility: Build76 historical scheduler view and Build84 qualified sampler view remain populated, with separate Build86 candidate counters;
6. CLI/buildinfo consistency and `go vet ./...`.

## Benchmark-first decision gate

Run on the qualified Go 1.26.0 host:

```bash
make clean
make v4-build86-dct-hoist-test
make v4-build86-dct-hoist-benchmark
```

The benchmark uses only the public deterministic Build83 fixture. It compares Build84 and Build86 reader front/mild/angle in the same session, compares historical vs hoisted DCT accumulation in isolation, and captures compiler BCE/inlining evidence for both complete readers.

Build86 should **not** proceed to the private nine-photo gate unless the complete reader shows a stable and useful improvement. The working decision threshold is approximately **5% whole-reader speedup** across the relevant geometries; isolated DCT speedup alone is not sufficient.

## Conditional physical gate

Only after a positive benchmark:

```bash
make V4_PHONE_BUILD86_DIAGNOSTIC_DIR='v4-phone private/build86-diagnostics1' \
  v4-build86-phone-physical-test
```

A promotable Build86 must preserve the complete historical 9/9 matrix and deep counters/workload, then repeat the result in a second independent run. Physical timing is compared against the source-controlled qualified Build84 reference in `docs/qualified-baselines/build84-phone-performance.tsv`.

## Status

**Build84 remains the CURRENT QUALIFIED SMARTPHONE BASELINE. Build86 is benchmark-first candidate only.**
