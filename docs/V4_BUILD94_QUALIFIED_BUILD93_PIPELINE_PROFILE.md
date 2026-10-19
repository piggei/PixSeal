# Format-v4 Build94 — qualified Build93 post-promotion pipeline profiling

Build94 is an **observability-only checkpoint** over the integration-closed, current qualified Build93 smartphone runtime. It introduces no recovery algorithm, changes no production path and cannot become a runtime baseline by itself.

## Why Build94

Build93 removed repeated protected-grid sampling after the geometry bank and ordered-parallel qualification were already frozen. Its two independent physical qualification runs reduced B/mild protected decode from the Build92/Build90 measurement of 33,460 ms to 11,378.5 ms mean while preserving exact semantics. The subsequent production-entry-point 9/9 gate closed the Build93 integration risk.

That change invalidates Build92 as the current priority map. Build94 therefore re-measures the **qualified Build93 runtime as it exists after promotion and consolidation** before any Build95 optimization hypothesis is opened.

## Runtime semantics

The smartphone runtime remains Build93. Build94 adds only test/profile harness code, version metadata and documentation. The measured path retains:

- the Build84 exact RGB-fetch geometry reader;
- the complete geometry-bank freeze barrier;
- Build90 ordered-parallel public-pilot qualification with bank-indexed results and serial ordered commit;
- Build93 one-pass protected-residue sampling with unchanged Robust -> Balanced -> Capacity derivation order;
- Build66 ordered batch consumption and first-logical-HMAC semantics;
- unchanged Format-v4 wire format, locked public pilot, strength 48, ECC/Hamming, whitening/HMAC domains, thresholds and quorum rules.

No Build94 timing, profile result or secret-dependent outcome participates in recovery decisions.

## Measurement design

Build94 separates **authoritative low-overhead stage timing** from **diagnostic CPU profiling**.

### Timing pass

The retained B/mild and B/angle cases are decoded/prepared with the deterministic JPEG path and historical pre-deep Build41 preparation. The unchanged `experimentalV4PhoneBuild93Recover` runtime is then executed without CPU profiling. Stage accounting is derived from timers already present in Build93/Build90/Build84 telemetry:

- deep total and geometry total;
- geometry pixel-plane preparation;
- Build47 freeze;
- prefix1, generation 2, generation 3 and generation 4 wall time;
- worker-time sums for prefix1/gen2/gen3/gen4;
- accounted and residual geometry time;
- post-geometry pixel-plane preparation;
- ordered-parallel qualification;
- Build93 protected-data decode wall time;
- accounted and residual post-geometry/total time;
- frozen geometry, bank, qualification, logical decode and list-frame counters;
- physical decode candidates, profiles and list frames;
- protected sampling passes and profile derivations;
- cumulative sampling/list worker time;
- worker/task counts and Build84 FoldScore/block-read workload.

The timing pass is the authoritative source for current stage shares. Sampling/list worker times are cumulative across parallel workers and must not be interpreted as wall-time shares.

### CPU profiling passes

B/mild and B/angle are executed again in separate later passes with Go CPU profiling enabled only around `experimentalV4PhoneBuild93Recover`. Preparation remains outside the profile. Their wall times are diagnostic and are not compared with the low-overhead timing pass.

The harness generates:

- flat and cumulative top reports;
- geometry-, generation3-, generation4-, continuation4- and FoldScore-focused reports;
- qualification-focused reports;
- Build93 decode-, protected-sampling-, profile-derivation- and list-decode-focused reports;
- source listings for Build93 recovery/decode/reuse, Build90 qualification, Build84 geometry/reader, FoldScore, projective sampling and ordered batch consumption.

One stable `watermark` test binary is retained beside the `.pprof` files so profiles remain symbolizable after the working tree advances.

## Semantic safety gate

Build94 reuses the frozen Build91/Build93 retained-corpus expectations and must preserve them exactly before accepting any profile.

B/mild must remain:

- geometry evaluations `79,259`;
- bank `937`;
- qualification evaluations `1,872`;
- qualified `935`;
- logical decode candidates `691`;
- logical list frames `2,120,047`;
- physical decode `696`, including `5` speculative candidates;
- physical profiles `2,086`;
- physical list frames `2,135,407`;
- protected sampling passes `696`;
- profile derivations `2,086`;
- FoldScores `16,496`;
- block reads `21,114,880`, all successful;
- authenticated expected payload PASS.

B/angle must remain:

- geometry evaluations `334,857`;
- bank and qualification evaluations `6,198`;
- qualified/decode/protected-sampling counters all zero;
- FoldScores `118,004`;
- block reads `151,045,120` (`149,370,664` successful / `1,674,456` failed);
- expected authentication REJECT.

The Build93 parallel-decode race regression, CLI help, version consistency and `go vet` are rerun by the source gate.

## Commands

Source/race/consistency gate:

```bash
make clean
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build94-qualified-pipeline-test
```

Qualified-host retained deep profile:

```bash
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build94-qualified-pipeline-profile
```

The profile publishes atomically under:

```text
v4-phone private/build94-profile/
├── build94-metadata.txt
├── build94-test-binary-build.txt
├── build94-watermark.test
├── build94-timing-raw.txt
├── build94-pipeline.tsv
├── build94-pipeline.md
├── build94-b-mild.pprof
├── build94-b-mild-raw.txt
├── build94-b-mild-pprof-top.txt
├── build94-b-mild-pprof-cum.txt
├── build94-b-mild-pprof-geometry.txt
├── build94-b-mild-pprof-gen3.txt
├── build94-b-mild-pprof-gen4.txt
├── build94-b-mild-pprof-continuation4.txt
├── build94-b-mild-pprof-foldscore.txt
├── build94-b-mild-pprof-qualification.txt
├── build94-b-mild-pprof-decode.txt
├── build94-b-mild-pprof-protected-sampling.txt
├── build94-b-mild-pprof-profile-derivation.txt
├── build94-b-mild-pprof-list-decode.txt
├── build94-b-mild-pprof-list.txt
├── build94-b-angle.pprof
├── build94-b-angle-raw.txt
├── build94-b-angle-pprof-top.txt
├── build94-b-angle-pprof-cum.txt
├── build94-b-angle-pprof-geometry.txt
├── build94-b-angle-pprof-gen3.txt
├── build94-b-angle-pprof-gen4.txt
├── build94-b-angle-pprof-continuation4.txt
├── build94-b-angle-pprof-foldscore.txt
├── build94-b-angle-pprof-qualification.txt
├── build94-b-angle-pprof-decode.txt
├── build94-b-angle-pprof-protected-sampling.txt
├── build94-b-angle-pprof-profile-derivation.txt
├── build94-b-angle-pprof-list-decode.txt
└── build94-b-angle-pprof-list.txt
```

If a profile fails, the staged evidence is preserved under `build94-profile-failed`.

## Decision rule

Build94 does **not** preselect a Build95 target. Review the qualified-host output in this order:

1. use the low-overhead timing rows to determine the current largest wall-time stages separately for B/mild and B/angle;
2. use flat/cumulative and focused CPU reports to localize cost inside those stages;
3. distinguish a B/mild improvement from a B/angle-only improvement and state which qualified case is expected to benefit;
4. prefer a structural, bounded change with explicit invariants over another reader micro-rewrite;
5. do not reopen Build80-88 cache/LUT/inlining/BCE experiments without genuinely new profile evidence;
6. if host drift makes a later candidate comparison ambiguous, use contemporaneous paired AB/BA evidence;
7. any Build95 runtime candidate must prove exact semantic equivalence, pass a targeted race gate and public same-session benchmark, then complete two independent Go 1.26.0 retained 9/9 physical PASS runs before it can supersede Build93.

## Qualified-host result — completed 2026-10-16

The Go 1.26.0 source gate and the complete profile both passed. The retained low-overhead timing rows preserved all frozen Build93 semantics and measured:

| case | deep total | geometry | geometry share | freeze | gen3 | gen4 | qualification | decode |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| B/mild | 80,152 ms | 57,175 ms | 71.33% | 21,768 ms | 9,906 ms | 17,463 ms | 12,511 ms | 10,251 ms |
| B/angle | 178,627 ms | 176,938 ms | 99.05% | 20,536 ms | 36,072 ms | 106,635 ms | 1,453 ms | 0 ms |

B/mild kept `79,259 / 937 / 1,872 / 935 / 691 / 2,120,047`, physical/speculative `696/5`, `2,086` physical profiles, `2,135,407` physical list frames, `696` protected sampling passes, `2,086` derivations and authenticated payload PASS. B/angle kept `334,857 / 6,198 / 6,198 / 0`, `118,004` FoldScores, `151,045,120` block reads and REJECT.

The separate CPU profiles localized the next target:

- B/mild: 513.24 CPU-s sampled; generation 4 135.96 CPU-s cumulative (26.49%); Build84 continuation4 reader 51.60 CPU-s (10.05%); Build93 protected sample 69.82 CPU-s (13.60%).
- B/angle: 1,384.63 CPU-s sampled; generation 4 920.46 CPU-s cumulative (66.48%); continuation4 459.29 CPU-s (33.17%); exact Build84 continuation4 reader 456.30 CPU-s (32.95%); qualification only 10.06 CPU-s (0.73%).

All major geometry pools already used eight workers on the eight-logical-CPU host. Build94 therefore did not select another outer scheduling change. It selected **Build95 Exact Gen4 Luminance Plane** as a bounded benchmark-first hypothesis: precompute exact per-pixel luminance once immediately before generation 4 and reuse it only in continuation4 projective reads, with a public-dimension 128 MiB cap and exact Build84 fallback. Build93 remains qualified. See `V4_BUILD95_EXACT_GEN4_LUMINANCE_PLANE.md`.
