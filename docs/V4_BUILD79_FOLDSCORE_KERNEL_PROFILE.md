# Format-v4 Build79 — FoldScore kernel profiling

Build79 is **observability-only** over the qualified Build76 smartphone baseline. Build76 was authoritative at that historical checkpoint; Build84 is the current qualified baseline. Build79 preserves the exact Build75 freeze, Build76 prefix1/gen2 scheduling, Build73 generation-three pool, Build71 generation-four task order/commit, Build68 qualification-plane reuse and Build66 ordered protected-data decode. It instruments only the `experimentalV4PhoneBuild41FoldScore` kernel when invoked by generation-four continuation work.

## Why Build79 exists

The retained Build78 Go 1.26.0 physical matrix passed exact 9/9 semantic equivalence and localized essentially all continuation4 cost inside FoldScore:

```text
B/mild: continuation4 worker 67,832 ms; FoldScore 67,777 ms (99.92%)
B/angle: continuation4 worker 500,652 ms; FoldScore 500,211 ms (99.91%)
```

Build78 also showed that movement-limit checks, homography rejects, duplicate accepted states and continuation bookkeeping are negligible relative to scoring. On B/mild, all 16,496 continuation score evaluations were valid and 15,557 (94.31%) were non-improving. On B/angle, 118,004 FoldScore evaluations were performed after 476 movement-limit rejects; 111,801 valid scores (94.74%) were non-improving. Those probes cannot be skipped safely merely because they fail to improve: their outcome is known only after FoldScore and participates in deterministic coordinate-descent choice.

The next question is therefore computational, not scheduling: which part of FoldScore and its projective block reader dominates cost?

## Frozen algorithm

Build79 changes no recovery decision. It preserves exactly:

- the Build41 FoldScore tile/fold selection and arithmetic order;
- the same 64 locked pilot positions/signs;
- the same `readProjectiveBlockValue` result for every authoritative block read;
- Build55 continuation probe order and acceptance rule;
- Build75 ordered-parallel Build47 freeze;
- Build76 ordered generation-two pool;
- Build73 ordered generation-three pool;
- Build71 generation-four task creation, bounded worker count and ordered commit;
- Build68 qualification-plane reuse;
- Build66 candidate-stable protected-data decode and first-logical-HMAC semantics.

Secret key, payload, ECC, HMAC, reference/oracle geometry and measured timing never create, select, rank or retain geometry.

## Exact counters

Every continuation4 FoldScore call records exact, low-overhead counters:

- FoldScore calls;
- selected tile visits;
- pilot-position visits;
- projective block-read attempts;
- successful and failed block reads;
- visible pilot samples used in the score;
- cumulative FoldScore worker time.

These counters are observational only and do not affect the score.

## Deterministic timing sample

Timing every projective pixel sample would distort the workload materially. Build79 therefore uses a deterministic 1/64 FoldScore sample rule:

```text
(gen4Task*131 + pairRank*17 + scoreOrdinal) mod 64 == 0
```

All inputs to the rule are public traversal indices. The rule is independent of timing, payload, key, ECC, HMAC and final authentication.

For sampled FoldScore calls, Build79 measures cumulative time spent in authoritative `readProjectiveBlockValue` calls and derives the residual FoldScore loop/accumulation overhead.

## Detailed block replay

For each sampled FoldScore, the first successful authoritative block is remembered. **After the authoritative FoldScore result is complete**, that block is replayed diagnostically in three separate phases:

1. homography `mapPoint` for all 64 block pixels;
2. bilinear `samplePlaneLuminance` for the mapped coordinates;
3. the two DCT coefficient accumulations used by `readProjectiveBlockValue`.

The replay result is discarded. It cannot modify FoldScore, continuation state, bank order or qualification. Replay timing is deliberately separate from the authoritative FoldScore elapsed counter.

This gives a low-perturbation estimate of whether the next optimization should target homography mapping, bilinear luminance/RGB conversion or DCT accumulation.

## Equivalence contract

Build79 is protected by three direct regressions:

1. `experimentalV4PhoneBuild79FoldScoreProfiled` must return exactly the same score and visible count as Build41 FoldScore with sampling both disabled and enabled;
2. Build79 profiled continuation must return exactly the same state sequence, quads, homographies, proposal values and evaluation count as Build55 continuation;
3. the complete Build79 blind bank must match Build76 exactly in evaluation count, bank size/order, quad, homography, proposal and validation values.

## Physical gate

```bash
make clean
make v4-build79-phone-foldscore-profile-test
make v4-build79-phone-physical-test
```

Diagnostics are staged and published only after a complete 9/9 semantic-equivalence PASS:

```text
v4-phone private/build79-diagnostics/build79-phone-foldscore-profile.tsv
v4-phone private/build79-diagnostics/build79-phone-foldscore-profile.md
```

Build79 is not promotable from timing. At the Build79 checkpoint, Build76 remained the qualified smartphone baseline; Build84 later superseded it. Build79 exists only to select a later implementation experiment from measured FoldScore kernel cost.
