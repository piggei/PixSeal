# Build98 — post-promotion consolidation

Build98 exact Generation3 fetch is the **current qualified smartphone baseline**. Build93 is the previous qualified baseline.

## Promotion evidence

Build98 passed, in order:

1. source/exactness/race gates on Go 1.26.0;
2. repeated public complete-Generation3 benchmark: 23.513 -> 21.304 ms/op (`-9.40%`) with unchanged allocations and clean five-sample separation;
3. contemporaneous Build93/Build98 AB/BA retained semantic gate, with Generation3 faster in all four B/mild/B/angle pairs;
4. two independent complete nine-photo physical gates through `v4-extract-phone-build98`, both 9/9 PASS with exact frozen semantic telemetry.

The two physical runs are frozen in `docs/qualified-baselines/build98-phone-performance.tsv`. Absolute cross-run timing retains host drift and is not used to manufacture a performance claim. Promotion is supported by the same-session changed-phase evidence plus two independent physical semantic passes.

## Production rewiring

Normal `v4-extract-phone` now delegates directly to the retained, physically qualified `ExperimentalV4ExtractPhoneBuild98Candidate` implementation. The implementation file is not rewritten. `v4-extract-phone-build98` remains an explicit qualified alias; `v4-extract-phone-build93` remains the previous-qualified alias.

No encoder, Format-v4 wire format, pilot, strength, ECC/Hamming, whitening/HMAC domain, geometry proposal/ranking, freeze barrier, Build90 qualification, Build93 protected-residue reuse, candidate order or first-logical-HMAC behavior changes during consolidation.

## Production integration — CLOSED / PASS

The archived Go 1.26.0 source gate passed and normal `v4-extract-phone` completed the separate nine-photo production-entry-point matrix on 2026-10-20 with **9/9 PASS**. Every semantic field matches both promotion runs; B/mild authenticates, B/angle and all three controls reject, and the four historical successful captures retain their qualified paths.

The non-sensitive matrix and original run metadata are frozen in `docs/qualified-baselines/build98-production-integration.tsv`. Metadata confirms `entrypoint=v4-extract-phone`, `gate_purpose=production-entrypoint-integration`, current Build98 / previous Build93, nine cases and zero additional image-sized memory. This is wiring evidence, not a third promotion-performance sample. Raw logs, images and payloads remain private.

The original console reports and diagnostic ZIP checksums were verified during the 2026-10-23 handoff closure. The retained implementation was not changed. The source/exactness/race/CLI/version/vet gates were also independently rerun successfully on Go 1.26.0 during takeover.

For future requalification:

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build98-gen3-fetch-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build98-phone-production-integration-test
```

The integration harness writes to `v4-phone private/build98-production-integration/`, separately from historical promotion artifacts. Documentation closure is complete; the next research step is observability-only profiling of unchanged Build98 before selecting another runtime optimization.
