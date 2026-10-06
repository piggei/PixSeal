# Format-v4 Build86 — exact DCT table hoist / bounds-check experiment

Build86 is a **closed exact / benchmark-negative / non-promoted** experiment over the current qualified Build84 smartphone baseline. It was selected from Build85 CPU/BCE evidence and changed only DCT cosine-table load placement inside continuation4 projective block reads.

## Exact candidate change

Qualified Build84:

```go
c23 += l * table3[x] * table2[y]
c32 += l * table2[x] * table3[y]
```

Build86:

```go
t2y := table2[y]
t3y := table3[y]
...
t2x := table2[x]
t3x := table3[x]
c23 += l * t3x * t2y
c32 += l * t2x * t3y
```

The floating-point order remained `(l*a)*b`; c23/c32 update order was unchanged. Build86 did not precompute DCT products and did not change the qualified Build84 RGB sampler, `mapPoint`, bilinear interpolation, geometry, bank/order, qualification or protected-data semantics.

## Correctness result

`make v4-build86-dct-hoist-test` passed on the qualified Go 1.26.0 host. Build84-vs-Build86 block/FoldScore/continuation/bank exactness and CLI/buildinfo/vet checks all passed. The public benchmark fixture also passed, with 0 B/op and 0 allocs/op for both readers.

## Qualified-host benchmark result

Five-run means on the 13th Gen Intel Core i3-13100T / Go 1.26.0 host:

| Benchmark | Build84 | Build86 | Build86 delta |
| --- | ---: | ---: | ---: |
| front reader | 1348.4 ns/block | 1382.4 ns/block | +2.5% slower |
| mild reader | 1380.8 ns/block | 1513.0 ns/block | +9.6% slower |
| angle reader | 1425.4 ns/block | 1435.6 ns/block | +0.7% slower |
| isolated DCT accumulation | 65.588 ns/block | 71.166 ns/block | +8.5% slower |

The benchmark-first criterion therefore failed. Build86 was **not** run on the private nine-photo corpus.

## Compiler lesson

Go 1.26.0 BCE evidence for qualified Build84 showed four DCT-table `IsInBounds` checks on the two accumulation lines. Build86 removed those checks from the inner accumulation statements, but introduced bounds checks on the hoisted row/column table loads. The Build84 reader inlining cost was 562; Build86 rose to 582.

This is a useful negative result: reducing the visible number/location of bounds checks did not imply lower runtime. The transformed source likely changed code generation, scheduling/register pressure or surrounding optimization enough to regress the complete reader.

## Decision

**Do not promote Build86. Do not run a physical gate.** Build84 remains the current qualified smartphone baseline. Build87 restores Build84 as the active deep runtime and profiles the remaining RGB/luminance/bilinear path before any further candidate is selected.
