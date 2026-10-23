# Build100 — ordered-parallel public pair scoring

**Closed / not promoted. Current qualified smartphone baseline: Build98. Previous: Build93.**

Build99 measured the six-pair proposal-score loop as a remaining serial section in freeze (12.29 / 10.27 sampled CPU-s on B/mild / B/angle). Build100 changes only the schedule of these six fixed public tasks. Seeds and task slots are prepared before the pool starts. Workers share immutable image/pilot data, write only their own result slot, and all scores/evaluation counts are committed in original pair order before the unchanged stable ranking, cell selection and basin stage.

No key, payload, ECC, CRC or HMAC is accepted by the scorer. Per-pair arithmetic, top-five robust reduction, tie order, search graph, bank freeze, Gen3/Gen4 kernels, Build90 qualification and Build93 protected-residue reuse remain unchanged. No image-sized allocation is introduced. The qualified Build98 implementation stays pinned; normal `v4-extract-phone` continues to use it. The candidate is explicitly available through `v4-extract-phone-build100`.

## Source and public benchmark — PASS

Bit-exact score/seed/order/evaluation comparisons pass with 1, 2, 6 and capped 16 requested workers. Complete freeze comparisons pass with both early and full bank caps. Complete blind-bank and recovery semantic telemetry matches Build98 after removing durations, and the concurrent six-worker fixture passes the race detector. CLI help/dimension tests pass, with inherited Build98, frozen v3 and locked pilot checks preserved.

Five qualified Go 1.26.0 / i3-13100T / GOMAXPROCS=8 complete-freeze benchmark samples measured 2,393.223 → 1,234.146 ms/op mean (**48.43% reduction**). Every candidate sample is faster than every baseline sample. Mean allocations remain about 3.99 MB per freeze with no image-sized addition. This passes the predeclared ≥5% mean improvement and all-corresponding-samples-faster rule for opening retained comparison; it does not promote the candidate.

## Final retained decision — closed, not promoted (2026-10-23)

Both independent complete AB/BA sessions passed exact semantics in all **16 deep-recovery executions**. Every paired freeze exceeded the predeclared 5% improvement threshold and every six-pair score phase improved. Full non-timing recovery telemetry, payload/info and the separate Gen3 counters remained identical. The non-sensitive complete rows and per-pair decision metrics are frozen in `docs/qualified-baselines/build100-paired-summary.tsv` and `build100-decision.tsv`.

| session | case | Build98 mean deep | Build100 mean deep | deep reduction | freeze reduction | unchanged Gen4 change |
|---|---|---:|---:|---:|---:|---:|
| r1 | B/mild | 84,618.0 ms | 74,932.5 ms | 11.45% | 42.52% | +1.37% |
| r1 | B/angle | 183,724.0 ms | 181,320.5 ms | 1.31% | 36.94% | +3.82% |
| r2 | B/mild | 87,380.5 ms | 74,103.0 ms | 15.20% | 41.63% | -10.26% |
| r2 | B/angle | 189,119.5 ms | 181,276.0 ms | 4.15% | 39.72% | +1.28% |

Positive Gen4 change means slower. Favorable overall means are real and are not discarded: Build100 does improve the complete retained mean in both cases/sessions. However, the predeclared guard against a repeatable slowdown in an unchanged stage was not satisfied. The independent replication again measured a slower B/angle Gen4 mean (+1.28% after +3.82% in session 1). The replication rule was fixed before its B/angle results and is not relaxed after observing them.

**Decision: close Build100 without promotion. No Build100 nine-photo physical gate was run, no production rewiring occurred and no Build100 integration PASS is claimed.** Build98 remains current; Build93 remains previous. The physical harness is prepared for possible future requalification only; its presence is not evidence of a completed physical run. The implemented candidate, strong public benchmark and positive changed-phase results are retained for scientific reproducibility, rather than described as an algorithmic failure or a lack of end-to-end benefit.

The cause of the additional Gen4 cost is not established by wall timing. Source inspection confirms the shared Gen4 kernel is unchanged, and structural comparison confirms the downstream bank scheduler and post-score freeze body match the qualified code. A later observability-only investigation should compare baseline/candidate CPU, worker timing, heap/GC and scheduling effects in one session before revising this candidate or selecting another runtime optimization. Do not infer that a slower mean proves an arithmetic or functional regression.

Private evidence archives are `build100-benchmark.zip`, `build100-paired-profile-r1.zip` and `build100-paired-profile-r2.zip`, referenced by `build100-evidence.sha256`. Images, keys, payloads and raw logs remain private. Source/race, CLI, vet, frozen-v3 and pilot gates passed. The candidate's runtime source hashes match both retained sessions exactly.

## Reproduction

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build100-pair-scoring-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build100-pair-scoring-benchmark
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build100-pair-scoring-paired-profile
# Deferred: Build100 did not satisfy its full retained review criterion:
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build100-phone-physical-test
```

Output directories are private, staged and never overwrite an existing run. Use `V4_BUILD100_BENCH_DIR`, `V4_PHONE_BUILD100_PROFILE_DIR` and `V4_PHONE_BUILD100_DIAGNOSTIC_DIR` to retain additional independent runs. Raw logs, acquisition images, keys and payloads remain excluded from public commits.
