# Format-v4 Build92 — qualified Build90 post-promotion pipeline profiling

Build92 is an **observability-only post-promotion checkpoint** over the current qualified Build90 smartphone runtime. It introduces no recovery algorithm and cannot become a runtime baseline by itself.

## Why Build92

Build90 was promoted after its public benchmark/race gate, three complete retained 9/9 semantic PASS runs, and the contemporaneous Build91 AB/BA physical profile. The paired decision showed that ordered-parallel qualification reduces B/mild deep recovery by 29.89% and B/angle by 4.00%, with qualification speedups of 4.286x and 4.866x respectively while geometry/decode remain comparable.

That promotion changes the cost distribution enough that pre-Build90 profiling is no longer an adequate basis for another optimization. Build92 therefore measures the **qualified Build90 runtime as it exists after promotion** before any Build93 hypothesis is opened.

## Runtime semantics

The production smartphone path remains Build90. Build92 adds only test/profile harness code and documentation. `ExperimentalV4ExtractPhone` still executes the unchanged Build90 deep fallback, including:

- the Build84 exact RGB-fetch geometry path;
- the complete freeze barrier before qualification;
- Build90 ordered-parallel public-pilot qualification with bank-indexed results and serial ordered commit;
- Build66/67 ordered protected-data decode and first-logical-HMAC semantics;
- unchanged Format-v4 wire format, pilot, strength 48, ECC/Hamming, whitening/HMAC domains and thresholds.

No Build92 timing or profiling result participates in recovery decisions.

## Measurement design

Build92 separates **wall-time accounting** from **CPU profiling**.

### Timing pass

The retained B/mild and B/angle cases are decoded/prepared with the qualified deterministic JPEG path and historical pre-deep Build41 preparation. The unchanged Build90 recovery is then executed without CPU profiling. Stage accounting is derived from timers already present in Build90/Build84 telemetry:

- deep total;
- geometry total;
- geometry pixel-plane preparation;
- Build47 freeze;
- prefix1;
- generation 2;
- generation 3;
- generation 4;
- worker-time sums for prefix1/gen2/gen3/gen4;
- accounted/residual geometry time;
- post-geometry pixel-plane preparation;
- qualification;
- protected-data decode;
- accounted/residual post-geometry and total time;
- geometry/bank/qualification/decode/list counters;
- qualification/decode worker counts and task counts;
- Build84 continuation4 FoldScore/block-read workload.

The timing pass is the authoritative source for stage shares.

### CPU profiling passes

B/mild and B/angle are then executed again in separate later passes with Go CPU profiling enabled **only around `experimentalV4PhoneBuild90Recover`**. Preparation is outside the CPU profile. These runs are diagnostic only; their wall time is not compared with the timing pass.

The harness generates:

- flat top reports;
- cumulative top reports;
- geometry-focused reports (`experimentalV4PhoneBuild84BlindBank` stacks);
- qualification-focused reports (`experimentalV4PhoneBuild68Qualify` stacks);
- decode-focused reports (`experimentalV4PhoneBuild67DecodeSingle` stacks);
- source-list reports for Build90 recovery, blind-bank geometry, qualification, decode, batch consumption and the qualified Build84 projective block reader.

One stable `watermark` test binary is retained beside the `.pprof` files so profiles remain symbolizable after the working tree advances.

## Semantic safety gate

Build92 reuses the frozen Build91 retained-corpus expectations and therefore requires exact Build90 semantics before accepting a profile.

B/mild must remain:

- geometry evaluations `79,259`;
- bank `937`;
- qualification evaluations `1,872`;
- qualified `935`;
- logical decode candidates `691`;
- list frames `2,120,047`;
- physical decode `696`;
- speculative decode `5`;
- FoldScores `16,496`;
- block reads `21,114,880`, all successful;
- authenticated expected payload PASS.

B/angle must remain:

- geometry evaluations `334,857`;
- bank `6,198`;
- qualification evaluations `6,198`;
- qualified `0`;
- FoldScores `118,004`;
- block reads `151,045,120` (`149,370,664` successful / `1,674,456` failed);
- expected authentication reject.

The Build90 qualification race regression is also rerun by the source gate.

## Commands

Source/race/consistency gate:

```bash
make clean
make v4-build92-qualified-pipeline-test
```

Qualified-host retained deep profile:

```bash
make v4-build92-qualified-pipeline-profile
```

The profile publishes atomically under:

```text
v4-phone private/build92-profile/
├── build92-metadata.txt
├── build92-test-binary-build.txt
├── build92-watermark.test
├── build92-timing-raw.txt
├── build92-pipeline.tsv
├── build92-pipeline.md
├── build92-b-mild.pprof
├── build92-b-mild-raw.txt
├── build92-b-mild-pprof-top.txt
├── build92-b-mild-pprof-cum.txt
├── build92-b-mild-pprof-geometry.txt
├── build92-b-mild-pprof-qualification.txt
├── build92-b-mild-pprof-decode.txt
├── build92-b-mild-pprof-list.txt
├── build92-b-angle.pprof
├── build92-b-angle-raw.txt
├── build92-b-angle-pprof-top.txt
├── build92-b-angle-pprof-cum.txt
├── build92-b-angle-pprof-geometry.txt
├── build92-b-angle-pprof-qualification.txt
├── build92-b-angle-pprof-decode.txt
└── build92-b-angle-pprof-list.txt
```

If a profile fails, the staged evidence is preserved under `build92-profile-failed`.

## Decision rule

Build92 does **not** select a Build93 in advance. After the qualified-host profile:

1. compare timing-stage shares first;
2. use flat/cumulative and phase-focused CPU reports to identify the dominant code inside the largest stage;
3. prefer one narrow, measurable, semantics-preserving hypothesis;
4. do not reopen negative reader micro-optimization branches without new evidence;
5. any later runtime candidate must prove exact equivalence first and then pass at least two independent Go 1.26.0 physical 9/9 gates with a repeatable material improvement before it can supersede Build90.
