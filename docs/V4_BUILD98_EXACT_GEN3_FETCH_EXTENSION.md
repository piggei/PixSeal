# Build98 — exact Generation3 fetch extension

Build98 is the **current qualified smartphone baseline**. Build93 is the previous qualified baseline. Two independent complete Go 1.26.0 nine-photo gates passed 9/9 with exact frozen semantics; normal production is now rewired to the unchanged qualified Build98 implementation.

## Changed scope

Build98 reuses the already-qualified Build84 exact RGB/FoldScore reader for every FoldScore inside Generation3:

- single3;
- pair3;
- continuation3;
- sibling3.

Prefix1, Generation2, the entire Generation4, Build75 freeze, Build90 qualification, Build93 protected-residue reuse, candidate order, ECC/Hamming, whitening/HMAC and first-logical-HMAC remain unchanged. No image-sized memory is added. Geometry still uses public pilot/structure only.

## Public benchmark

The Go 1.26.0/i3-13100T source, exactness and race gates passed. Five benchmark samples showed clean separation:

| metric | historical mean | Build98 mean | reduction | allocations |
|---|---:|---:|---:|---:|
| sibling stencil | 3,088,615.4 ns/op | 2,793,988.4 ns/op | 9.54% | 1 -> 1 |
| complete Generation3 | 23,512,975.6 ns/op | 21,303,687.8 ns/op | 9.40% | 3 -> 3 |

## Retained AB/BA result

The same decoded working image was reused inside each Build93/Build98 pair and order was reversed in cycle 2. Exact outcome, bank, qualification, protected decode, historical continuation4 telemetry and payload/HMAC behavior matched. Candidate-local Generation3 reads are recorded separately as `Gen3Fetch*`.

| case | order | Build93 Gen3 ms | Build98 Gen3 ms | reduction |
|---|---|---:|---:|---:|
| B/mild | 93->98 | 10,391 | 9,929 | 4.45% |
| B/mild | 98->93 | 13,074 | 10,004 | 23.48% |
| B/angle | 93->98 | 39,788 | 37,723 | 5.19% |
| B/angle | 98->93 | 43,075 | 35,981 | 16.47% |

Generation3 is directionally faster in all four pairs. Total/geometry timing is mixed because unmodified Gen4 and downstream stages retain substantial host drift; those rows are preserved without cherry-picking. The exact non-sensitive evidence is frozen in `docs/qualified-baselines/build98-paired-summary.tsv`.

## Qualified entry points

`v4-extract-phone-build98` pins the qualified implementation. Normal `v4-extract-phone` delegates to the same unchanged Build98 implementation; `v4-extract-phone-build93` pins the previous qualified baseline. Public diagnostics retain historical `build84-rgb-fetch`, `build90-qualification-parallel` and `build93-residue-reuse` lines and add:

```text
build98-gen3-fetch: attempted=true fold-scores=... block-reads=... block-success=... block-failed=...
```

The physical harness requires the exact historical nine-photo semantics and exact frozen B/mild/B/angle counters. B/mild additionally requires candidate authentication; B/angle must remain REJECT. Failed staging is preserved atomically.

## Completed physical qualification

The following commands were run twice independently and both matrices passed 9/9:

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build98-gen3-fetch-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build98-phone-physical-test
```

The run1/run2 console reports and diagnostic ZIPs are retained privately. Non-sensitive rows are frozen in `docs/qualified-baselines/build98-phone-performance.tsv`. The separate normal-command integration completed 9/9 PASS on 2026-10-20; archived checksums and semantic equivalence were verified on 2026-10-23. Non-sensitive integration rows are frozen in `docs/qualified-baselines/build98-production-integration.tsv`. Future wiring requalification uses `make v4-build98-phone-production-integration-test`.
