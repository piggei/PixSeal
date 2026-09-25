# Build66 deterministic ordered-parallel protected-data decode

Build66 is the **current qualified smartphone performance baseline**, derived from the qualified Build64 recovery semantics and the semantically equivalent Build65 seed-parallel experiment. It does not change the Format-v4 encoder, public pilot, strength 48, proposal score, geometry search bounds, held-out qualification, ECC/Hamming logic, whitening/HMAC domains, thresholds, quorums or fallback ordering.

## Build65 result

The qualified Go 1.26.0 host run of Build65 passed the complete retained nine-photo equivalence gate. Every deep-fallback case reproduced the exact Build64 seed/evaluation/bank/qualification/logical-decode/list-frame counts and all payload/HMAC outcomes. The fallback used eight seed workers on that host.

Build65 is not promoted. The Build64 qualification artifact did not record matched `elapsed_ms`, so no rigorous same-host speedup can be computed retroactively. More importantly, the measured Build65 B/mild recovery still required **828036 ms** while logically evaluating **691 qualified decode candidates and 2,120,047 list frames**. Seed-parallel geometry therefore does not address the dominant successful-case cost strongly enough by itself. At the end of Build65, Build64 therefore remained the qualified baseline.

## Build66 change

Build66 keeps the Build65 blind-bank generation exactly:

1. run unchanged Build47 freeze and Build48 top4-per-pair seed selection;
2. run the unchanged Build64 per-seed deep search concurrently;
3. store each seed result in its original slot;
4. concatenate the complete blind bank in original Build64 seed order;
5. qualify the complete frozen bank serially in original Build64 order.

Only after this unchanged freeze and qualification does Build66 change scheduling.

Qualified protected-data candidates are processed in bounded batches of
`min(GOMAXPROCS, qualified_candidates)`. All candidates in one batch may be physically evaluated concurrently, but the batch is then **consumed strictly in original Build64 candidate order**.

If a batch contains the first logical HMAC success:

- candidates before it contribute exactly their Build64 profiles/list frames/max-confidence;
- the successful candidate contributes exactly its Build64 result;
- candidates after it in the same physical batch are discarded semantically;
- later batches are never started.

Therefore the first logical HMAC success, returned payload, logical decode-candidate count, logical list-frame count, profile and max-confidence must remain exactly Build64-identical. At most `workers-1` candidates may perform speculative physical work around the winning batch, but they cannot affect selection or semantic telemetry.

Setting `GOMAXPROCS=1` reduces both Build65 seed scheduling and Build66 data scheduling to the serial form.

## Non-negotiable equivalence

For the retained Build38 nine-photo corpus, Build66 must preserve the qualified Build64 deep-fallback telemetry exactly:

| image | seeds | geometry evals | bank | qualified | logical decode candidates | logical list frames | result |
|---|---:|---:|---:|---:|---:|---:|---|
| control/front | 24 | 18021 | 26 | 0 | 0 | 0 | reject |
| control/mild | 24 | 117609 | 810 | 6 | 6 | 18432 | reject |
| control/angle | 24 | 21203 | 11 | 0 | 0 | 0 | reject |
| B/mild | 24 | 79259 | 937 | 935 | 691 | 2120047 | authenticate `v4-b38-phone-b` |
| B/angle | 24 | 334857 | 6198 | 0 | 0 | 0 | reject |

The four already-fast authenticated cases (`A/front`, `A/mild`, `A/angle`, `B/front`) must continue to return before the deep fallback and must therefore show neither Build64/65/66 recovery telemetry.

## Telemetry

When deep recovery runs, Build66 keeps the semantic Build64 line unchanged and retains Build65 scheduling telemetry:

```text
build64-recovery: attempted=true seeds=... geometry-evals=... bank=... qualified=... decode-candidates=... list-frames=... authenticated=...
build65-parallel: attempted=true workers=N order=seed-stable
```

It adds:

```text
build66-decode-parallel: attempted=true workers=N order=candidate-stable-batches
```

For a fallback with zero qualified candidates, Build66 reports `workers=0` because no protected-data decode is scheduled.

## Qualification and performance gate

Run on the qualified Go 1.26.0 host:

```sh
make v4-build66-phone-parallel-decode-test
make v4-build66-phone-physical-test
```

The physical gate writes:

```text
v4-phone private/build66-diagnostics/build66-phone-matrix.tsv
v4-phone private/build66-diagnostics/build66-phone-matrix.md
```

A correctness PASS requires:

- the same nine-photo outcomes as Build64;
- exact Build64 deep seed/evaluation/bank/qualification/logical-decode/logical-list-frame counts;
- Build65 seed-parallel telemetry on every fallback case and nowhere else;
- Build66 ordered-decode telemetry on every fallback case and nowhere else;
- exact payload/HMAC behavior preserved.

If `v4-phone private/build65-diagnostics/build65-phone-matrix.tsv` exists, the Build66 report also includes the previous Build65 `elapsed_ms` and computes `speedup_vs_build65` per image. Timing is informational and never affects correctness.

Build66 should supersede Build64 only if both conditions are true:

1. semantic equivalence is exact;
2. B/mild and/or the complete retained gate show a clearly useful same-host wall-clock gain without unacceptable resource cost.

The qualified host result below satisfies both conditions; **Build66 supersedes Build64 as the current qualified smartphone baseline**.


## Qualified host result — 2026-09-25

The complete retained nine-photo Go 1.26.0 gate passed with exact Build64/65 semantic equivalence. All five deep-fallback cases reproduced the same seed, geometry-evaluation, frozen-bank, qualified-bank, logical decode-candidate and logical list-frame counts. B/mild remained `937/935/691/2120047`, authenticated exact `v4-b38-phone-b`, and B/angle remained `6198/0` with no HMAC.

On the same host used for Build65, B/mild improved from **828036 ms** to **544038 ms** (**1.522x**, **-34.3%**). The complete nine-photo matrix improved from **2372255 ms** to **2175687 ms** (**1.090x**, **-8.3%**, saving 196568 ms). Control/mild and B/angle were slower in this single run, but neither changes semantic outcome; B/angle has zero qualified decode candidates, so its variation is outside the Build66 decode-parallel work itself.

The qualification artifact was generated by the qualified binary before a packaging-only correction renamed the matrix files from the accidental `build65-phone-matrix.*` names to `build66-phone-matrix.*`. No decoder, search, qualification, HMAC, scheduling or telemetry semantics changed in that correction.

**Decision:** promote Build66 as the current qualified smartphone baseline.
