# Format-v4 Build81 — exact luminance LUT

Build81 is a **closed equivalence-preserving performance experiment** over the qualified Build76 smartphone baseline. It passed semantic qualification twice but did not demonstrate a material, repeatable speedup and therefore was **not promoted**. Build76 was authoritative at that historical checkpoint; Build84 is the current qualified baseline.

## Optimization tested

Only continuation4 FoldScore luminance conversion changed. The historical expression:

```go
.299*float64(r) + .587*float64(g) + .114*float64(b) - 128
```

was replaced by three process-wide 256-entry `float64` product tables, preserving the original addition grouping, source coordinates, floor/clamp behavior, bilinear interpolation, DCT accumulation, geometry ordering and protected-data semantics. No cache, bounding rectangle or fallback path was used.

## Exactness result

Local regressions proved:

1. every table entry equals direct multiplication bit-for-bit;
2. all 16,777,216 RGB triples reproduce the original luminance expression bit-for-bit;
3. sample, projective block and Build41 FoldScore equality;
4. Build55 continuation4 state/evaluation equality;
5. complete Build81-vs-Build76 blind-bank equality.

Both retained Go 1.26.0 physical runs then passed the complete 9/9 semantic matrix. Deep workload remained exactly:

```text
B/mild: 79259 evals / bank 937 / qualified 935 / decode 691 / 2120047 frames / HMAC PASS
B/angle: 334857 evals / bank 6198 / qualified 0 / REJECT
```

Continuation4 workload also remained identical between runs:

```text
B/mild:  16496 FoldScore /  21114880 block reads /  21114880 success /       0 failed
B/angle: 118004 FoldScore / 151045120 block reads / 149370664 success / 1674456 failed
```

## Performance result

| metric | Build76 mean | Build81 run 1 | Build81 run 2 | Build81 mean | delta |
|---|---:|---:|---:|---:|---:|
| full matrix | 758,102 ms | 911,984 | 798,601 | 855,293 | +12.8% |
| B/mild elapsed | 188,750.5 | 218,597 | 195,081 | 206,839 | +9.6% |
| B/mild geometry | 60,351.5 | 73,443 | 64,148 | 68,796 | +14.0% |
| B/mild gen4 | 18,473 | 22,342 | 18,385 | 20,364 | +10.2% |
| B/angle elapsed | 226,813 | 285,826 | 228,901 | 257,364 | +13.5% |
| B/angle geometry | 183,806.5 | 232,949 | 186,322 | 209,636 | +14.1% |
| B/angle gen4 | 110,829.5 | 138,009 | 108,872 | 123,441 | +11.4% |

Run 2 alone looked slightly favorable in B/angle gen4 and essentially neutral in B/mild gen4, but the result did not repeat in run 1 and the complete matrix was slower in both runs. The host/runtime variance is too large to attribute a benefit to the LUT.

## Decision

**NON PROMUOVERE / DO NOT PROMOTE.** Retain Build81 as evidence that exact scalar RGB product lookup is semantically safe but not a useful performance optimization for this workload. The next experiment must target helper-call/reader structure rather than another cache or lookup layer.
