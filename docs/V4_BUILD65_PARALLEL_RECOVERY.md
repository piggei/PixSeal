# Build65 deterministic seed-parallel deep recovery

Build65 was a **performance-only candidate** derived from the qualified Build64 smartphone baseline. It does not change the Format-v4 encoder, public pilot, strength 48, proposal score, search bounds, held-out qualification, protected-data decoder, ECC/Hamming logic, whitening/HMAC domains, thresholds, quorums or fallback ordering.

## Motivation

Build64 qualified the full Build63-derived deep recovery and recovered the retained `B/mild` photograph, but the fallback is computationally expensive on difficult rejects and on the recovered B/mild case. The search starts from up to 24 Build48 seeds (top4 independently for each of six side pairs). Once those seeds have been frozen proposal-only, each seed branch is independent until the final geometry banks are concatenated.

Build65 therefore changes **scheduling only**:

1. run the unchanged Build47 freeze and Build48 seed selection serially;
2. execute each already-selected seed branch concurrently;
3. store each result in its original seed slot;
4. wait for all seed workers;
5. concatenate seed results in the original Build64 seed order;
6. run the unchanged Build64 qualification and protected-data/HMAC loop serially in that preserved order.

The worker count is bounded by `min(GOMAXPROCS, number_of_seeds)`. Setting `GOMAXPROCS=1` therefore gives the serial scheduling form of the same Build65 code path.

## Non-negotiable equivalence

Build65 is not allowed to prune, deduplicate, reorder or rescore geometry. For the retained Build38 nine-photo corpus the deep-fallback telemetry must be exactly identical to the qualified Build64 artifact:

| image | seeds | geometry evals | bank | qualified | decode candidates | list frames | result |
|---|---:|---:|---:|---:|---:|---:|---|
| control/front | 24 | 18021 | 26 | 0 | 0 | 0 | reject |
| control/mild | 24 | 117609 | 810 | 6 | 6 | 18432 | reject |
| control/angle | 24 | 21203 | 11 | 0 | 0 | 0 | reject |
| B/mild | 24 | 79259 | 937 | 935 | 691 | 2120047 | authenticate `v4-b38-phone-b` |
| B/angle | 24 | 334857 | 6198 | 0 | 0 | 0 | reject |

The four already-fast authenticated cases (`A/front`, `A/mild`, `A/angle`, `B/front`) must still return before deep recovery and must therefore show neither Build64 nor Build65 telemetry.

## Telemetry

When the deep fallback runs, `v4-extract-phone` keeps the qualified Build64 line unchanged and adds:

```text
build65-parallel: attempted=true workers=N order=seed-stable
```

`build64-recovery:` remains the semantic telemetry because the generated bank, qualification order and decode order are required to be Build64-identical.

## Qualification gate

Run on the qualified Go 1.26.0 host:

```sh
make v4-build65-phone-parallel-test
make v4-build65-phone-physical-test
```

The physical gate records per-image elapsed milliseconds in:

```text
v4-phone private/build65-diagnostics/build65-phone-matrix.tsv
v4-phone private/build65-diagnostics/build65-phone-matrix.md
```

A PASS requires:

- the same nine-photo outcomes as Build64;
- exact Build64 deep-recovery seed/evaluation/bank/qualification/decode/list-frame counts on every fallback case;
- Build65 parallel telemetry on every fallback case and nowhere else;
- exact payload/HMAC behavior preserved.

A correctness PASS alone does **not** automatically supersede Build64. Build65 should be promoted only after the host run also demonstrates a useful wall-clock improvement. At that Build65 checkpoint, Build64 remained the qualified baseline.

## Qualified-host result

The Go 1.26.0 host run completed after this candidate was prepared. The complete nine-photo gate passed and every deep-fallback semantic counter matched Build64 exactly. The host used eight seed workers on fallback cases. B/mild remained `bank=937`, `qualified=935`, `decode-candidates=691`, `list-frames=2120047` and authenticated exact `v4-b38-phone-b`; B/angle remained `bank=6198`, `qualified=0` and rejected.

Build65 is therefore closed as **equivalence-proven but not promoted**. The historical Build64 gate did not record matched per-image timing, so a rigorous same-host Build64 speedup cannot be reconstructed. The Build65 B/mild run itself required 828036 ms, showing that serial protected-data decoding remains a major cost. Build64 stays the qualified baseline and Build66 moves the performance experiment to ordered-parallel decode while preserving Build65's already-proven bank semantics.
