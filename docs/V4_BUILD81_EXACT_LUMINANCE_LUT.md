# Format-v4 Build81 — exact luminance LUT

Build81 is an **equivalence-preserving performance candidate** over the qualified Build76 smartphone baseline. Build76 remains authoritative until two physical runs reproduce the complete semantic matrix and a material speedup.

## Evidence

Build79 showed that projective block reads consume 97.20% of sampled continuation4 FoldScore time on B/mild and 97.47% on B/angle. Build80 then tested a spatial luminance cache. Its physical run preserved the qualified counters exactly but regressed performance:

```text
full matrix:      806,849 ms vs Build76 mean 758,102 ms   (+6.4%)
B/mild elapsed:   191,873 ms vs 188,750.5 ms              (+1.7%)
B/mild geometry:   62,344 ms vs 60,351.5 ms               (+3.3%)
B/angle elapsed:  257,674 ms vs 226,813 ms                (+13.6%)
B/angle geometry: 211,148 ms vs 183,806.5 ms              (+14.9%)
```

B/angle had `151,045,120` continuation4 block reads, zero cache hits, `149,370,664` valid fallbacks and `1,674,456` failed reads. The Build80 gate originally mislabeled that valid no-hit case as a telemetry failure; the retained diagnostics prove geometry itself remained `334857 / 6198 / 0` and REJECT. The implementation is rejected because of performance, not correctness.

## Build81 optimization

The original integer-source luminance expression is:

```go
.299*float64(r) + .587*float64(g) + .114*float64(b) - 128
```

Each channel is 8-bit. Build81 therefore precomputes three process-wide tables:

```text
R[n] = .299 * float64(n)
G[n] = .587 * float64(n)
B[n] = .114 * float64(n)
```

for `n = 0..255`. The hot path evaluates:

```go
value := R[r] + G[g]
value += B[b]
return value - 128
```

The grouping deliberately matches Go evaluation order of the original expression. No source coordinates, floor/clamp behavior, interpolation fractions, bilinear operations or DCT accumulation are changed. There is no bounding rectangle, no per-block scratch buffer and no fallback path. The LUT occupies 3 × 256 × 8 = 6,144 bytes.

## Scope

The LUT is used only by FoldScore calls inside generation-four `continuation4`. `single4`, `pair4`, Build76 prefix1/gen2, Build73 gen3, Build75 freeze, Build68 qualification and Build66 ordered protected-data decode remain unchanged. No key, payload, ECC, HMAC, oracle or timing data can influence geometry.

## Exactness regressions

Build81 requires:

1. every table entry to equal the direct multiplication bit-for-bit;
2. all 16,777,216 RGB triples to equal the original luminance expression bit-for-bit;
3. bilinear `samplePlaneLuminance` equality at representative in-bounds/edge/out-of-bounds coordinates;
4. projective block equality across multiple homographies and block positions;
5. FoldScore equality against Build41;
6. continuation4 state/evaluation equality against Build55;
7. complete blind-bank equality against Build76.

## Physical qualification

```bash
make clean
make v4-build81-phone-luminance-lut-test
make v4-build81-phone-physical-test
```

Outputs:

```text
v4-phone private/build81-diagnostics/build81-phone-luminance-lut.tsv
v4-phone private/build81-diagnostics/build81-phone-luminance-lut.md
```

The gate publishes diagnostics atomically on PASS and preserves them in `build81-diagnostics-failed/` on failure. Timing is evaluated only after semantic equivalence passes.
