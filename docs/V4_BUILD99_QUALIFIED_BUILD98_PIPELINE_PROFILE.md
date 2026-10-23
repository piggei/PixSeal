# Build99 — qualified Build98 post-promotion pipeline profiling

**Observability-only. Qualified runtime: Build98. Previous qualified baseline: Build93.**

Build98 production integration is closed. Build99 measures the unchanged Build98 deep-recovery implementation on B/mild and B/angle before any new optimization is chosen. The version/buildinfo bump identifies the measurement snapshot; it does not promote a new recovery algorithm. No encoder, wire format, pilot, ECC/Hamming, whitening/HMAC domain, threshold, geometry generation/ranking, bank order, freeze barrier, qualification or protected-decode code changes.

## Method and semantic gates

One timing pass runs B/mild then B/angle without CPU profiling. Existing qualified telemetry supplies geometry-plane preparation, freeze, prefix1, Gen2/3/4, post-geometry plane, qualification and decode durations. Historical Build94 accounting/preparation helpers are reused without alteration. CPU profiles are collected only in two later, separate passes around the unchanged `experimentalV4PhoneBuild98Recover` call. Worker CPU sums are cumulative, not wall-time shares. Preparation and JPEG decoding remain outside the CPU profile.

Every timing and CPU pass enforces the frozen inherited geometry, qualification, ordered-decode, protected sampling and historical continuation4 read counters. In addition, Build98 Generation3 counters must match exactly:

| case | FoldScores | reads | successful | failed |
|---|---:|---:|---:|---:|
| B/mild | 23,024 | 29,470,720 | 29,470,720 | 0 |
| B/angle | 76,817 | 98,325,760 | 97,394,512 | 931,248 |

B/mild must authenticate the expected payload; B/angle must reject with no protected decode. Measurements cannot guide geometry or scheduling. Public source tests deliberately skip the private retained runs unless configured; a full profile requires both timing rows and two actual CPU-test PASS results.

## Reproduction

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build99-qualified-pipeline-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build99-qualified-pipeline-profile
```

Default corpus: `v4-phone private/build38-acquired/`. Default private output: `v4-phone private/build99-profile/`. Override `V4_PHONE_BUILD99_PROFILE_DIR` to retain another run separately. Existing outputs are refused; failed staging is preserved under a unique failed directory. The harness records host, toolchain, GOMAXPROCS (default 8), canonical dimensions, acquisition/source hashes and a retained symbolizable test binary. It never publishes images, keys, payloads or raw logs in the source documentation.

## Result — completed PASS on 2026-10-23

The qualified Go 1.26.0 source gate passed, including Build98 exactness/race/CLI/version/vet, frozen Format-v3 hashes, pilot identity and public Build94/99 harness regressions. The complete retained profile then passed both timing cases and both separate CPU-profile cases on the i3-13100T with GOMAXPROCS=8. All inherited semantic counters and additive Build98 Gen3 counters matched exactly in every pass. B/mild authenticated; B/angle rejected with zero protected decode. This is a two-case observability PASS, not a new nine-photo qualification or a runtime promotion.

The low-overhead timing pass measured:

| case | deep total | geometry | freeze | Gen3 | Gen4 | qualification | decode |
|---|---:|---:|---:|---:|---:|---:|---:|
| B/mild | 83,618 ms | 60,628 ms | 22,567 ms | 10,609 ms | 18,907 ms | 12,678 ms | 10,098 ms |
| B/angle | 187,601 ms | 185,870 ms | 21,276 ms | 36,337 ms | 113,604 ms | 1,479 ms | 0 ms |

Geometry is 72.51% / 99.08% of deep wall time. Gen4 remains the largest B/angle phase at 60.56%; B/mild freeze is 26.99% of deep time. These are current stage shares; no cross-session speedup is inferred.

Separate CPU profiles measured:

| case | total sampled CPU | Gen3 cumulative | Gen4 cumulative | serial pair scoring | Gen4 exact reader | Gen4 generic reader |
|---|---:|---:|---:|---:|---:|---:|
| B/mild | 518.84 s | 74.86 s | 137.39 s | 12.29 s | 52.01 s | 84.52 s |
| B/angle | 1,343.65 s | 285.71 s | 879.78 s | 10.27 s | 433.87 s | 439.71 s |

CPU cumulative values overlap through call chains and must not be added or interpreted as wall time. The exact reader is shared by Gen3 and continuation4 after Build98, so the two Gen4 reader columns are explicitly filtered to stacks containing `experimentalV4PhoneBuild84Generation4Fetch`. A helper-only focus on `experimentalV4PhoneBuild84ContinueFetch` (the historical `continuation4` report filename) includes both generations; it is not an isolated Gen4 measure. The generic geometry focus also matches named worker stacks rather than providing a complete logical-phase CPU sum. Wall-time telemetry remains the phase-share authority.

Non-sensitive timing/CPU rows are frozen in `docs/qualified-baselines/build99-profile-summary.tsv`. Private artifacts retain the test binary, raw logs, profiles, metadata, exact source/corpus hashes and additional phase-filtered reports. The original profile and measurement-source checksum manifest were verified after publication. The complete private evidence is archived as `build99-profile.zip` and referenced by `build99-profile-evidence.sha256`; no private image or raw payload is published in documentation.

## Decision — next candidate: ordered-parallel pair scoring

Do not automatically reopen the rejected Build95 luminance plane or Build97 Gen4 probe-reader extension. Gen4 is still dominant on B/angle, but another reader rewrite does not follow from that fact alone. Source inspection and both new profiles expose a different bounded serial section: `experimentalV4PhoneBuild75FreezeParallel` calls `experimentalV4PhoneBuild43PairRobustScore` for the six fixed public side pairs in sequence. Pair scoring alone consumes 12.29 / 10.27 sampled CPU-s, within the 22.567 / 21.276-second timing-pass freeze stages. This identifies a structural scheduling hypothesis without image-sized allocation. It does not establish its eventual speedup.

The proposed Build100 candidate is limited to computing those six independent public-proposal scores in a bounded worker pool, storing results by original pair index and committing in historical order before the unchanged stable ranking/cell/basin pipeline. Per-pair arithmetic, robust-score reduction, tie order, evaluation counts, complete-bank freeze, downstream Gen3/4, Build90 qualification and Build93 protected decode must remain identical. Key, payload, ECC and HMAC cannot influence the work or scheduling. The existing qualified runtime remains untouched until a separate candidate is implemented and evaluated.

Candidate acceptance must begin with exact pair scores/order/counts and full-bank/recovery equivalence, targeted race tests, a repeated same-session complete-freeze benchmark and contemporaneous retained comparison. Only a favorable result can justify two independent nine-photo promotion gates. Build98 remains current and Build93 previous; Build99 is closed observability-only PASS.
