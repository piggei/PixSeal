# Format-v4 Build78 — continuation4 internal profiling

Build78 is **observability-only** over the qualified Build76 smartphone baseline. Build76 was authoritative at that historical checkpoint; Build84 is the current qualified baseline. Build78 preserves the exact Build76/Build75 recovery bank, task scheduling, ordered commits, qualification and protected-data semantics while instrumenting the existing fourth-generation continuation loop.

## Why Build78 exists

Build77 passed the retained nine-photo semantic-equivalence gate and showed that generation four is already scheduled almost ideally across eight workers. On the retained physical run:

```text
B/mild: gen4 wall 20,868 ms, worker 165,708 ms (~7.94x)
B/angle: gen4 wall 121,551 ms, worker 971,006 ms (~7.99x)
```

The internal split was:

```text
B/mild: single4 39,065 ms; pair4 62,994 ms; continuation4 63,643 ms
B/angle: single4 126,386 ms; pair4 355,904 ms; continuation4 488,696 ms
```

Continuation4 is therefore the largest measured generation-four component on B/angle and tied with pair4 on B/mild. Another scheduling barrier is not justified; the next question is what continuation4 itself spends work on.

## Frozen algorithm

Build78 preserves exactly:

- Build75 ordered-parallel Build47 freeze;
- Build76 prefix through `sibling1` and ordered gen2 pool;
- Build73 ordered gen3 pool;
- Build71 generation-four task creation, worker count and ordered commit;
- the exact Build55 continuation coordinate-descent arithmetic and acceptance rule;
- proposal/validation scores, thresholds, search bounds and seed depth;
- Build68 qualification-plane reuse;
- Build66 ordered protected-data decode and first-logical-HMAC semantics.

Secret key, payload, ECC and HMAC never create, schedule, rank or retain geometry.

## New measurements

Build78 wraps the existing `experimentalV4PhoneBuild55Continue` loop without changing its output. For every continuation4 call it records:

- input/call count;
- total +/- coordinate probe attempts;
- probes rejected by the movement limit;
- homography-construction rejects;
- actual fold-score evaluations;
- invalid scores;
- valid non-improving probes;
- proposal-improving probes;
- accepted continuation states;
- duplicate accepted quads (diagnostic only);
- passes through the eight coordinate dimensions;
- total continuation worker time;
- time spent inside `experimentalV4PhoneBuild41FoldScore`;
- residual preparation/control-flow time outside FoldScore;
- min/median/max call time, evaluations and accepted-state count.

The dominant continuation call is selected **deterministically by evaluation count**, with original traversal order as the tie-break. Timing, HMAC outcome and oracle/reference geometry cannot select it.

## Equivalence contract

Two regressions protect the measurement:

1. the profiled continuation wrapper must return exactly the same states, state order, quads, homographies, proposal values and evaluation count as `experimentalV4PhoneBuild55Continue`;
2. the complete Build78 blind bank must match Build76 exactly in evaluation count, bank size/order and geometric values.

## Physical gate

```bash
make clean
make v4-build78-phone-continuation4-profile-test
make v4-build78-phone-physical-test
```

Diagnostics are published atomically only after a complete 9/9 PASS:

```text
v4-phone private/build78-diagnostics/build78-phone-continuation4-profile.tsv
v4-phone private/build78-diagnostics/build78-phone-continuation4-profile.md
```

Build78 is not promotable from timing. Its measurements exist only to decide whether a later equivalence-preserving implementation optimization is justified.
