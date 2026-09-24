# Build63 — Post-Fourth-Pair Continuation

## Status

Research-only. At this research checkpoint, Build44 was the latest qualified production milestone. Build63 changes no production encoder/decoder path, proposal score, threshold, quorum, pilot, ECC/Hamming, whitening or HMAC domain.

## Evidence entering Build63

The qualified Go 1.26.0 Build62 artifact closes B/mild as `fourth-pair-qualified-gain`: 492 fourth-pair states are frozen, 491 qualify, 0 authenticate, and the best qualified target state reaches **24.056 px** from its exact **25.641 px** parent. B/angle freezes **3746** fourth-pair states with **0 qualified and 0 authenticated**.

## Blind algorithm

For every retained Build62 fourth-pair state, Build63 runs the established 1px proposal-only coordinate descent across 8 dimensions for at most 8 passes, retaining every accepted intermediate (maximum 64 states per branch). The complete B/mild+B/angle continuation bank is frozen before qualification, protected-data decode or HMAC. SIFT/reference geometry is generated only after both blind JSON files exist.

## Commands

```sh
make v4-build63-phone-fourth-pair-continuation-test
make v4-build63-phone-fourth-pair-continuation-diagnostic
```

The diagnostic writes `build63-fourth-pair-parents.tsv`, `build63-continuation-states.tsv`, `build63-fourth-pair-continuation-summary.tsv` and `build63-fourth-pair-continuation.md` under `v4-phone private/build63-diagnostics/`.

## Qualified host result

The Go 1.26.0 host run closes Build63 as **`fourth-pair-continuation-recovery`**.

| image | continuation states | qualified | authenticated | best oracle state | interpretation |
|---|---:|---:|---:|---:|---|
| B/mild | 937 | 935 | **1** | 22.099 px | recovery |
| B/angle | 6198 | 0 | 0 | n/a | not triggered / no qualified bank |

The single authenticated B/mild state is candidate 10 / side-pair rank 3 / seed rank 3 / fourth-pair rank 7 / continuation index 7. It has proposal `0.302815`, held-out validation `0.240486`, is qualified, and authenticates with the existing HMAC. Its post-hoc oracle error is **28.868 px** from an exact 28.586 px parent.

The best geometric state in the same bank is **22.099 px** and does **not** authenticate. This is useful evidence: the independent oracle did not select the successful frame, and minimizing corner error is not equivalent to maximizing protected-data recoverability.

Build63 is therefore the first blind smartphone recovery of retained B/mild. It justifies a separate production-candidate promotion experiment; it does not by itself promote the decoder.
