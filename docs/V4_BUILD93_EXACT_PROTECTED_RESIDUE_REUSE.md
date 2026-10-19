# Format-v4 Build93 — exact protected-residue reuse

Build93 is the **current qualified smartphone baseline**, selected from Build92 profiling and promoted after public exactness/race/benchmark gates plus two independent retained 9/9 physical semantic PASS runs. Build90 is now the previous qualified baseline.

## Build92 evidence

The retained Build92 timing pass completed with exact qualified semantics:

| case | deep total | geometry | qualification | decode |
|---|---:|---:|---:|---:|
| B/mild | 109,620 ms | 59,290 ms | 16,644 ms | 33,460 ms |
| B/angle | 193,565 ms | 191,762 ms | 1,431 ms | 0 ms |

B/mild therefore moved from a qualification-dominated pre-Build90 pipeline to one where geometry is the largest stage and protected-data decode is the second-largest stage.

The separate B/mild CPU profile attributes the protected decode almost entirely to repeated projective sampling:

- `experimentalV4PhoneBuild67DecodeSingle`: 235.96 CPU-s cumulative;
- `experimentalV4ReadProtectedProjectiveMargins`: 230.14 CPU-s cumulative;
- `readProjectiveBlockValue`: 228.66 CPU-s cumulative;
- `experimentalV4PhoneBuild42ListDecode`: 5.81 CPU-s cumulative.

The historical decoder calls `experimentalV4ReadProtectedProjectiveMargins` once for each attempted profile. For a candidate that reaches Capacity, the same homography/grid is therefore sampled three times even though the projective residue sums/counts are profile-independent.

## Build93 single change

For each already-qualified physical decode candidate, Build93:

1. executes the exact historical protected-grid traversal once;
2. accumulates the same `residueSums` and `residueCounts` in the same block order;
3. derives Robust margins from those immutable residues;
4. if authentication fails, derives Balanced margins from the same residues;
5. if authentication fails again, derives Capacity margins from the same residues.

The profile-dependent derivation preserves the historical loop order and floating-point expressions. No projective sample is approximated, cached across candidates or shared between homographies.

The old logical order remains:

```text
Robust -> Balanced -> Capacity
```

and the first logical authenticated profile remains authoritative.

## Security and semantic boundaries

Build93 starts only after Build90 has already:

- frozen the complete public-geometry bank;
- completed ordered-parallel public-pilot qualification;
- reconstructed the qualified candidate list in original bank order.

Protected residue reuse therefore cannot influence geometry proposal, ranking, pruning or qualification.

The reuse decision is unconditional and does not depend on key, ECC result, HMAC result, payload contents or message correctness. The key is still used only by the existing protected-frame decode/authentication path after margins exist.

Unchanged:

- encoder and Format-v4 wire format;
- locked public pilot and strength 48;
- geometry search and complete freeze barrier;
- Build90 qualification worker pool and ordered commit;
- candidate order and Build66 ordered batch consumption;
- Hamming/list decoding;
- whitening and HMAC domains;
- profile order and first-logical-HMAC behavior;
- all qualified thresholds/quorums.

## Candidate telemetry

Build93 adds candidate-local diagnostic counters:

- `ProtectedSamplingPasses` — successful protected-grid sampling passes physically executed;
- `ProtectedProfileDerivations` — Robust/Balanced/Capacity margin derivations executed from sampled residues.

Existing physical/logical decode telemetry remains authoritative. For candidates whose protected sampling succeeds, Build93 performs one protected-grid sampling pass regardless of how many profiles are subsequently derived. On all-success workloads such as the retained B/mild path and the public benchmark fixture, profile derivations equal the historical successful physical profile attempts. If a derived profile lacks complete observations, `ProtectedProfileDerivations` can exceed `PhysicalProfilesTried` because the former counts the attempted derivation while the historical counter increments only after a complete margin vector exists.

## Exactness gates

The public source gate requires:

- bit-exact margins and confidence versus the historical reader for Robust/Balanced/Capacity on front/mild/angle-like homographies;
- identical authenticated result, winning profile, payload, profile count, list frames and max confidence for public Robust/Balanced/Capacity carriers;
- identical wrong-key failure behavior;
- identical invalid/undersized sampling behavior;
- complete public deep-recovery semantic equality between Build90 and Build93;
- one protected sampling pass per successfully sampled physical candidate; on the all-success public recovery fixture, profile derivations must equal the historical physical profile attempts;
- failure fixtures must preserve the historical reject path without treating diagnostic counter relationships as semantic output;
- an eight-way shared-plane race test.

## Public benchmark

The benchmark fixture is a deterministic public 592x512 Capacity carrier at strength 48. Capacity intentionally forces the historical decoder through Robust, Balanced and Capacity while Build93 samples the protected grid once.

Two comparisons are measured in the same Go process/toolchain:

- historical Build67 single candidate vs Build93 single candidate;
- historical eight-way decode batch vs Build93 eight-way batch.

The benchmark is decision evidence only. A useful speedup with exact/race PASS can justify a later retained physical gate; it cannot promote Build93 by itself.

### Qualified-host result

Go 1.26.0 on the i3-13100T produced five retained samples per variant:

| workload | historical mean | Build93 mean | speedup | time reduction |
|---|---:|---:|---:|---:|
| single Capacity candidate | 44.439 ms/op | 10.236 ms/op | **4.341x** | **76.97%** |
| eight-way Capacity batch | 49.443 ms/op | 22.324 ms/op | **2.215x** | **54.85%** |

Allocation counts remained effectively unchanged, so the result reflects removal of repeated projective sampling rather than a memory-allocation trade. Source exactness and the targeted race gate both passed. The benchmark decision gate is therefore closed PASS.

## Physical qualification result

The retained gate uses:

```bash
make v4-build93-phone-physical-test
```

and writes `v4-phone private/build93-diagnostics/`. Two independent Go 1.26.0 runs completed the entire 9-photo matrix with semantic-equivalence PASS.

| metric | run1 | run2 | two-run mean |
|---|---:|---:|---:|
| full matrix elapsed | 839,352 ms | 649,960 ms | 744,656 ms |
| B/mild deep total | 97,913 ms | 75,704 ms | **86,808.5 ms** |
| B/mild geometry | 70,435 ms | 53,836 ms | 62,135.5 ms |
| B/mild qualification | 14,344 ms | 11,797 ms | 13,070.5 ms |
| B/mild decode | 12,878 ms | 9,879 ms | **11,378.5 ms** |
| B/angle deep total | 248,022 ms | 164,681 ms | 206,351.5 ms |

B/mild semantics are identical in both runs: `79,259` geometry evaluations, bank `937`, `1,872` qualification evaluations, `935` qualified candidates, `691` logical decode candidates, `2,120,047` logical list frames, `696` physical decode candidates, `5` speculative candidates, HMAC/payload PASS, `16,496` FoldScores and `21,114,880` successful block reads. Build93 adds exactly `696` protected sampling passes and `2,086` profile derivations in each run.

B/angle remains exactly `334,857 / 6,198 / 6,198 / 0`, expected REJECT, with `118,004` FoldScores and `151,045,120` block reads (`149,370,664` successful / `1,674,456` failed). Because no candidate qualifies, Build93 performs zero protected sampling or profile derivation there.

The changed B/mild decode phase improves from the Build92 Build90 measurement 33,460 ms to an 11,378.5 ms two-run mean: **2.941x faster / -65.99%**. B/mild deep recovery averages 86,808.5 ms versus the qualified Build90 contemporaneous paired reference 118,104 ms: **1.361x / -26.50%**. B/angle absolute geometry varies strongly across runs despite executing no Build93 sampling, reinforcing the established host-drift caution; that timing is not used to attribute a Build93 gain.

## Qualified runtime routing

Normal `v4-extract-phone` uses Build93 for the final deep fallback. The explicit `v4-extract-phone-build93` command remains a pinned Build93 regression/requalification alias. All earlier direct Build41/42/43 and residual fallback ordering remains unchanged.

Post-promotion consolidation keeps the physically qualified `ExperimentalV4ExtractPhoneBuild93Candidate` body unchanged and makes `ExperimentalV4ExtractPhone` a narrow wrapper around it. Before consolidation the two bodies were source-equivalent; the cleanup deletes only the duplicate production body and does not rewrite the algorithm.

The physical harness can identify either accepted entry point, records it in `build93-phone-run-metadata.txt`, and never records key or payload values. The Makefile physical target deliberately selects normal `v4-extract-phone` and publishes into `v4-phone private/build93-production-integration/` so the two historical promotion runs are not overwritten.

The historical two-run qualification remains authoritative. The additional qualified Go 1.26.0 production-command pass completed on 2026-10-16 with 9/9 PASS, exact semantic telemetry and metadata confirming normal `v4-extract-phone`. It closes post-promotion wiring/integration risk only. See [`V4_BUILD93_POST_PROMOTION_CONSOLIDATION.md`](V4_BUILD93_POST_PROMOTION_CONSOLIDATION.md).

Qualified-host commands:

```bash
make clean
GOTOOLCHAIN=go1.26.0 make v4-build93-residue-reuse-test
GOTOOLCHAIN=go1.26.0 make v4-build93-phone-production-integration-test
```

The backward-compatible `make v4-build93-phone-physical-test` target runs the same production-entry-point integration gate.

## Promotion decision

Build93 satisfied every promotion requirement:

1. exact public margins/confidence and full recovery equivalence;
2. targeted race PASS;
3. clear qualified-host public benchmark win;
4. two independent retained 9/9 semantic physical PASS runs;
5. repeatable material improvement in the protected decode phase actually changed;
6. no change to geometry, qualification semantics, candidate ordering, ECC/Hamming, whitening/HMAC domains, strength or wire format.

**Decision: PROMOTED. Build93 is the current qualified smartphone baseline; Build90 becomes the previous qualified baseline.**

Frozen promotion evidence is stored in `docs/qualified-baselines/build93-phone-performance.tsv`; the non-sensitive production-entry-point confirmation is stored separately in `docs/qualified-baselines/build93-production-integration.tsv`. Private images, keys, payload files and raw diagnostic logs are not shipped in source archives. Build93 is promoted and its post-promotion integration is closed; Build94 must begin as observability-only profiling before any new optimization is selected.
