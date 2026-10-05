# Format-v4 Build80 — exact local luminance cache

Build80 is an **equivalence-preserving performance candidate** over the qualified Build76 smartphone baseline. Build76 remains authoritative until Build80 reproduces the retained nine-photo semantic matrix and a material speedup on the qualified Go 1.26.0 host.

## Why Build80 exists

Build79 passed the retained 9/9 physical semantic-equivalence gate and showed that continuation4 FoldScore time is dominated by projective block reads:

```text
B/mild sampled FoldScore: 1,071 ms total; 1,041 ms block reads (97.20%)
B/angle sampled FoldScore: 7,424 ms total; 7,236 ms block reads (97.47%)
```

The same run counted 21,114,880 authoritative block reads on B/mild and 151,045,120 on B/angle. A detailed post-result replay pointed to bilinear luminance sampling as the largest component inside `readProjectiveBlockValue`; DCT accumulation and homography mapping were smaller. Build79 also showed that only ~1.11% of B/angle block reads fail visibility/bounds, so failed-read pruning is not a useful primary target.

## Optimization

The original `samplePlaneLuminance` computes RGB-to-luminance for four integer source pixels for each projected sample. Neighboring samples inside one 8x8 destination block often reuse the same source pixels.

Build80 preserves the original projective mapping and arithmetic but, for each continuation4 FoldScore block:

1. maps the same 64 destination pixels with the same homography;
2. records the same `floor`, clamp and bilinear fractions;
3. computes the small integer-source bounding rectangle required by those samples;
4. if that rectangle contains at most **196 pixels**, computes the exact float64 RGB-to-luminance value once per integer source pixel in a FoldScore-local scratch buffer;
5. performs the same bilinear interpolation expression in the same order;
6. accumulates the same two DCT coefficients in the same pixel order;
7. if the rectangle exceeds the 196-pixel bound, falls back to the unchanged `readProjectiveBlockValue` implementation.

The scratch storage is allocated once per FoldScore invocation and reused for all block reads in that score. It is not allocated per block.

## Scope

The cache is used **only by FoldScore calls inside generation-four continuation4**. Build80 intentionally leaves unchanged:

- Build75 ordered-parallel Build47 freeze;
- Build76 prefix1 and generation-two scheduling;
- Build73 generation-three pool;
- Build71 single4/pair4 scheduling and generation-four task order/commit;
- all proposal thresholds and acceptance rules;
- frozen bank order and provenance;
- Build68 qualification-plane reuse;
- Build66 ordered protected-data decode;
- ECC, Hamming, whitening and HMAC semantics.

No timing, secret, payload, ECC, HMAC or oracle/reference information affects the cache decision. The cache decision depends only on the public mapped block footprint.

## Exactness contract

Build80 requires bit-level numerical equivalence, not approximate equality.

Local regressions require:

1. cached block value and original `readProjectiveBlockValue` to have identical `ok` results and identical `math.Float64bits` output across multiple projective transforms and block positions;
2. the bounded-cache fallback path to return the exact original result;
3. cached FoldScore to match Build41 FoldScore bit-for-bit and preserve visible count;
4. cached continuation4 to reproduce Build55 state sequence, quads, homographies, proposal bits and evaluation count exactly;
5. the complete Build80 blind bank to equal Build76 in evaluation count, seed/worker counts, bank size/order, quad, homography, proposal and validation values.

## Low-overhead telemetry

Build80 records only integer counters so performance timing remains meaningful:

- block reads attempted by cached continuation4 FoldScore;
- cache hits;
- original-reader fallbacks;
- failed block reads;
- total integer-source luminance pixels prepared;
- maximum local cache rectangle area;
- cache limit (196 pixels).

For every deep path:

```text
cache_hits + cache_fallbacks + cache_failed == cache_reads
```

## Physical gate

```bash
make clean
make v4-build80-phone-luminance-cache-test
make v4-build80-phone-physical-test
```

Diagnostics are staged and published atomically only after a complete 9/9 semantic-equivalence PASS:

```text
v4-phone private/build80-diagnostics/build80-phone-luminance-cache.tsv
v4-phone private/build80-diagnostics/build80-phone-luminance-cache.md
```

The gate compares elapsed and geometry time with the retained Build76 TSV when available. Correctness remains authoritative; timing cannot compensate for any semantic mismatch.

A second independent physical run is required before Build80 can replace Build76 as the qualified smartphone baseline.

## Physical outcome — rejected performance candidate

The retained physical run preserved the qualified geometry counters, including B/mild `79259 / 937 / 935 / 691 / 2120047` with HMAC/payload PASS and B/angle `334857 / 6198 / 0` REJECT. The first gate revision incorrectly required `cache_hits > 0`, so B/angle was labeled `telemetry=false` even though zero hits were a valid consequence of its large mapped footprints. The diagnostic revision proved this was a harness error only.

Performance nevertheless regressed materially: full matrix 806,849 ms (+6.4% versus Build76 mean), B/mild geometry 62,344 ms (+3.3%) and B/angle geometry 211,148 ms (+14.9%). B/angle performed zero cache hits and 149,370,664 valid fallbacks, causing the original reader to repeat mapPoint work after the cache pre-pass. Build80 is therefore **rejected for performance** and is not eligible for promotion.
