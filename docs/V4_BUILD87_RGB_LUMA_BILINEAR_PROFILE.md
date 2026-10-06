# Format-v4 Build87 — qualified RGB/luminance/bilinear compiler profiling

Build87 is an **observability-only** checkpoint over the current qualified Build84 smartphone baseline. It is selected after Build86 closed as exact but benchmark-negative. Build87 restores the active deep runtime to Build84 and changes no recovery algorithm, geometry search, ordering, thresholds, qualification, protected-data decode, ECC/Hamming, whitening/HMAC domain, strength, pilot, or Format-v4 wire format.

## Why Build87 exists

Build85 re-profiled the promoted Build84 reader and showed that RGB fetch + luminance + bilinear interpolation remained the largest exact subcomponent. Build86 then tested a narrow DCT-table-hoist hypothesis. On the qualified Go 1.26.0 i3-13100T host, Build86 remained exact but the benchmark regressed instead of improving:

- whole reader front: Build84 **1348.4 ns/block**, Build86 **1382.4 ns/block** (+2.5% slower);
- whole reader mild: Build84 **1380.8 ns/block**, Build86 **1513.0 ns/block** (+9.6% slower);
- whole reader angle: Build84 **1425.4 ns/block**, Build86 **1435.6 ns/block** (+0.7% slower);
- isolated DCT accumulation: Build84 **65.588 ns/block**, Build86 **71.166 ns/block** (+8.5% slower).

Compiler evidence explains an important negative result: Build86 removed the four inner-loop DCT `IsInBounds` checks but introduced checks on the hoisted table loads and increased the reader inlining cost from 562 to 582. Fewer visible bounds checks did **not** produce faster code. Build86 therefore does not proceed to the private physical corpus and is not promoted.

The next justified question is narrower:

> Within the qualified Build84 RGB/luminance/bilinear path, which exact operation shape still dominates on Go 1.26.0, and what bounds-check / generated-assembly evidence is attached to it?

## Public deterministic workload

Build87 reuses the Build83 1632x1632 synthetic RGB plane, fixed front/mild/angle-like homographies and fixed interior block origins. No private image, secret key, payload, ECC result, HMAC result, expected message, oracle geometry or protected-data outcome selects benchmark work.

One benchmark operation is always one complete 8x8 block (64 samples), so stage `ns/op` values are directly comparable at block granularity.

## Measurements

The repeated benchmark records:

1. the complete qualified Build84 block reader on front/mild/angle-like homographies;
2. prepared four-pixel three-byte slice fetch + byte loads;
3. four exact RGB-to-luminance conversions per sample;
4. the two horizontal bilinear interpolations (`top` / `bottom`);
5. the final vertical bilinear interpolation;
6. the combined prepared fetch + luminance + bilinear path.

Prepared inputs are generated before the timed section and are checked bit-for-bit against the historical luminance sampler. They are a decomposition aid only; they are never used by the production decoder.

## Integrated CPU and compiler evidence

`make v4-build87-rgb-luma-bilinear-profile` also preserves:

- a long CPU pprof of `experimentalV4PhoneBuild84ReadProjectiveBlockValue` on the angle-like fixture;
- `pprof -top` and line-level `pprof -list` output;
- Go 1.26.0 inlining and BCE diagnostics for the qualified reader;
- full `go tool objdump` output for the qualified reader;
- an assembly subset covering the source lines that form the four RGB slices, calculate four luminances and perform bilinear interpolation.

This is evidence gathering only. Build87 itself cannot be promoted from timing because it introduces no runtime optimization.

## Commands

```bash
make clean
make v4-build87-rgb-luma-bilinear-test
make v4-build87-rgb-luma-bilinear-profile
```

Expected profile directory:

```text
v4-phone private/build87-profile/
├── build87-metadata.txt
├── build87-fixture-test.txt
├── build87-benchmark.txt
├── build87-profile-run.txt
├── build87-cpu.pprof
├── build87-pprof-top.txt
├── build87-pprof-reader-list.txt
├── build87-compiler.txt
├── build87-compiler-reader.txt
├── build87-bce-rgb.txt
├── build87-reader-asm.txt
├── build87-reader-asm-rgb.txt
└── watermark.test
```

## Decision rule for Build88

Build87 does not assume that another exact optimization exists. A Build88 performance candidate should be opened only if the qualified-host profile shows one narrow source/assembly transformation with a plausible reader-level benefit and with exact arithmetic/order preservation.

Out of scope unless new evidence explicitly justifies a separate research branch:

- `unsafe` RGB access;
- LUT/cache layers already contradicted by Build80/81 evidence;
- FMA or floating-point reassociation;
- precomputed luminance planes;
- projective-coordinate recurrences that change floating-point order;
- geometry pruning/ranking changes;
- any secret/payload/ECC/HMAC-guided workload selection.

If Build87 exposes no credible exact-portable target, Build84 should remain the qualified performance baseline rather than forcing another micro-optimization.
