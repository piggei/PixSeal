# Build101 — contemporaneous Build98/Build100 runtime effects

**Closed / observability PASS, 2026-10-27. Observability only. Qualified baseline Build98; previous Build93. Build100 remains closed / not promoted.**

Build100 improved its changed freeze stage and whole-pipeline retained means, but failed its predeclared unchanged-stage guard: B/angle Gen4 mean was slower by 3.82% and 1.28% in the two independent sessions. Build101 investigates this observation without reopening that decision or changing production kernels.

## Protocol fixed before retained execution

1. Eight low-overhead full-recovery executions call the unchanged pinned Build98 and Build100 APIs: B/mild 98→100, B/angle 98→100, B/angle 100→98, B/mild 100→98. Preparation, explicit pre-run GC and five-second cooldown occur outside timing. All full non-timing telemetry, payload and successful extraction info must match exactly; inherited frozen B/mild/B/angle expectations and separate Gen3 counters must also pass.
2. Four separate attribution executions use test-only structural copies with unchanged kernels, task queues and ordered collection. B/mild runs 98→100; B/angle runs 100→98. Stage labels propagate to workers for CPU sampling; memory/GC and runtime scheduler counters are sampled at freeze, prefix1, Gen2, Gen3, Gen4, qualification and (when reached) decode boundaries. Their synchronization and profiling overhead excludes these runs from speed decisions.
3. Record full Gen4 wall/worker accounting and CPU attribution, allocation/heap and GC deltas, and runnable scheduler-latency histogram deltas. Runtime CPU-class values are estimates, not OS process CPU. Histogram thresholds use buckets whose lower bounds are at least 1/10 ms; they do not reconstruct exact latency percentiles. Before/after heap sizes are snapshots, not a peak-memory bound.
4. Close the observability gate only if all 12 runs preserve semantics, structural-copy checks pass, source/corpus hashes stay fixed, and all profiles/rows are retained. Report each timing pair and both case means, without changing Build100's historical guard. GC/CPU correlation is descriptive; no causal conclusion from wall timing alone, and no claim of scheduling causation from one profile pair.

The allocation/GC/scheduler deltas are process-wide and include profiler/test infrastructure. Summed worker elapsed time is an occupancy proxy, not CPU utilization. In the local qualified Go 1.26.0 source, `runtime/metrics.go:cpuStatsAggregate.compute` copies the last STW CPU snapshot rather than accumulating to the read instant. Zero CPU-class deltas between GC updates therefore do not mean zero CPU consumption; sampled stage-labelled CPU profiles supply the independent CPU attribution. A zero completed-GC delta is reported as such, not as proof excluding all background GC work.

This build introduces no runtime candidate, physical qualification or production rewiring. A later runtime change still needs its own predeclared gates. A diagnostic pass cannot promote Build100.

## Reproduction

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build101-runtime-effects-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build101-runtime-effects-profile
```

The default private output is `v4-phone private/build101-runtime-effects`. Existing outputs are refused; failed staged evidence is preserved under a unique name. CPU profiles, symbolizable test binary, acquisition hashes and raw logs stay private. The reviewed derived non-sensitive tables are frozen under `docs/qualified-baselines/`.

## Results

Source/exactness/race/CLI/vet, inherited frozen-core/pilot and structural-copy gates passed on the qualified Go 1.26.0 / i3-13100T host. All **12 retained executions** passed full exact semantics: eight low-overhead timing runs and four separate attribution runs, with B/mild authenticated and B/angle rejected. Source/corpus hashes remained stable, all 26 stage rows and four symbolizable CPU profiles were retained, and the predeclared observability closure gate passed. No causal attribution or runtime promotion is claimed.

| case | Build98 mean deep ms | Build100 mean deep ms | deep change | freeze change | Gen4 change |
|---|---:|---:|---:|---:|---:|
| B/mild | 84,566.5 | 74,486.0 | -11.92% | -39.03% | -1.35% |
| B/angle | 187,536.5 | 177,391.5 | -5.41% | -37.68% | -2.21% |

Negative change means faster. Gen4 is faster in each of the four timing pairs, including both B/angle orders. The historical adverse B/angle mean is not reproduced in this session; this does not establish its cause or reopen Build100's closed decision.


## Separate attribution results and limits

These four instrumented runs are deliberately excluded from the timing means above. Unlike the untouched timing comparisons, they show a slower candidate Gen4 in both cases: +1.99% wall / +1.74% sampled CPU on B/mild, and +3.36% wall / +2.91% sampled CPU on B/angle. This differing sign is retained as evidence of measurement/session sensitivity; neither pass is discarded.

| case / variant | Gen4 wall ms | sampled Gen4 CPU s | worker occupancy proxy | Gen4 allocated bytes | completed GC cycles | GC pause ms |
|---|---:|---:|---:|---:|---:|---:|
| B/mild Build98 | 18,636 | 146.58 | 0.990932 | 17,343,544 | 0 | 0 |
| B/mild Build100 | 19,007 | 149.13 | 0.991345 | 17,343,656 | 0 | 0 |
| B/angle Build100 | 116,977 | 920.71 | 0.997568 | 120,057,848 | 1 | 0.349053 |
| B/angle Build98 | 113,171 | 894.71 | 0.997819 | 120,087,208 | 1 | 2.507698 |

Stage-labelled CPU samples occupy about 7.85–7.91 CPU-s per wall second with eight workers. The pool shows no large idle gap in these samples. Allocated bytes are nearly identical within each pair; the candidate does not show a new large allocation cost. B/angle completes one GC cycle in both variants, and the candidate's recorded pause is smaller despite its slower Gen4. These pauses alone do not explain the seconds of extra wall time; this does not exclude every GC or memory-locality effect.

The CPU cost remains concentrated in the same unchanged projective readers. B/angle historical `readProjectiveBlockValue` costs 449.34 / 458.69 sampled CPU-s (Build98 / Build100); the exact Build84 continuation reader costs 439.11 / 455.82 CPU-s. B/mild reader values are 90.82 / 91.29 and 54.79 / 56.85 CPU-s respectively. These are cumulative function values normalized separately from the whole-profile percentages shown by pprof. No new expensive call path or changed workload is identified.

Scheduler counters remain descriptive: B/angle Build100 has 5,913 runnable-latency events, 192 in buckets beginning at ≥1 ms and 102 at ≥10 ms; Build98 has 5,766 / 219 / 96. This single observed pair does not establish scheduler causation. CPU frequency, thermal state, hardware counters and cache locality were not controlled or measured, so none is asserted as the cause. The CPU-class snapshot cadence limitation above also prevents treating those deltas as live phase CPU measurements.

## Decision and next bounded study

**Close Build101 as observability PASS. Keep Build98 qualified, Build93 previous and Build100 closed / not promoted.** The historical slower mean is not reproduced in the untouched contemporary AB/BA timing pass, yet adverse Gen4 values reappear in the separate instrumented pass. Wall timing alone therefore remains insufficient for a causal claim, and this build does not reverse Build100's predeclared decision or authorize a physical qualification campaign.

The next bounded investigation should replay exactly the same frozen Gen4 task vector with the unchanged kernel in repeated balanced orders, separating prefix/freeze context from Gen4 execution. It should freeze the vector's count/order/score hashes and compare CPU/worker/GC evidence without moving Build100's historical guard. This is a proposed next study, not an implemented Build102 or a runtime optimization.

Frozen public evidence: `build101-runs.tsv` (all 12 runs, explicitly labelled timing/attribution), `build101-stages.tsv` (all 26 observed stage rows), `build101-timing-pairs.tsv` (all four timing pairs), `build101-timing-means.tsv` and `build101-gen4-attribution.tsv`. Private raw evidence is archived as `build101-runtime-effects.zip`, with its checksum and the public-table checksums in `build101-evidence.sha256`. No private corpus, key, payload, profile binary or raw log enters the public commit.
