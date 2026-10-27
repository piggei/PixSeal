# Build102 — frozen Gen4 replay and A/A measurement control

**Closed / observability PASS, 2026-10-27. Observability only. Qualified Build98 / previous Build93 unchanged; Build100 closed / not promoted.**

Build101 found opposite Gen4 signs between contemporaneous untouched timing and separate instrumented passes, with the same kernels/counters and no unique GC/scheduler explanation. Build102 isolates the work at the original Gen3→Gen4 barrier. It does not create a performance candidate or change the historical Build100 guard.

## Protocol fixed before retained results

For each retained B/mild and B/angle acquisition:

1. Run the unchanged qualified Build98 full blind bank once. Check frozen Build101 bank/evaluation/Gen4/fetch counts plus exact Build98 Gen3 counters. Fingerprint the entire ordered bank, including all unexported geometry/proposal fields and exact floating-point bits. This is bank equivalence, not a new payload/HMAC qualification campaign.
2. Capture the complete Gen4 vector from the unchanged Build98 and Build100 prefixes. Verify identical count/order/score fingerprints, full task values, raster bytes and non-timing prefix counters. A structural source check chains the test-only split to the verified Build101 copy and pinned Build100. Release the second capture before replay and retain one shared immutable raster/vector.
3. Invoke the identical Gen4 kernel, bounded pool and original-index collection on that shared capture eight times, in two balanced four-slot blocks: **A B B A | B A A B**. A and B are measurement labels only: same function, inputs, pointers and worker count. Explicit GC and five-second cooldown occur before each slot. There is no upstream freeze/prefix work between replays. Record the complete replay wall, original pool wall, summed worker durations, live Linux process user/system CPU, allocation and completed-GC/pause deltas. Snapshot reads and bitwise fingerprints stay outside the measured kernel window; boundary process counters include small read overhead. CLK_TCK is obtained from getconf and its quantization is recorded.
4. After that case's eight timing slots, run two separate CPU-profile slots, A→B on mild and B→A on angle. All four profiled executions are excluded from timing means. Complete output-bank fingerprints, semantic telemetry, frozen task/raster fingerprints and qualified reference counters must match on every replay. Process-wide allocation/CPU deltas include measurement/profiler infrastructure; summed worker time is an occupancy proxy.
5. Close only after **20 exact complete-bank replays**, four matching captures, two qualified full-bank references, all profiles/rows and stable source/corpus hashes. Retain all eight placebo A/B contrasts, each case's min/max/mean and within-code variability. No significance or causal attribution follows from an apparent A/B difference or overlap with historical 1–4% values. These are two balanced blocks in one process per case, not independent host-session replications; ranges and coefficients of variation are descriptive. Detached replay deliberately removes upstream context; it cannot alone identify which previous environmental or instrumentation effect caused a historical slowdown.

No protected decode is required or claimed by this observability gate. The pinned runtime source files, encoder/format/pilot/qualification and normal extraction dispatch are untouched. Any future optimization still requires its own gates; no Build102 physical campaign or promotion is authorized by this study.

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build102-frozen-gen4-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build102-frozen-gen4-profile
```

The default private output is `v4-phone private/build102-frozen-gen4-replay`. Existing outputs are refused and failed evidence is preserved uniquely. Raw profiles, binary, raster and candidate details remain private; only derived numeric rows and fingerprints will be frozen publicly after review.

## Results — complete PASS

Qualified Go 1.26.0 / i3-13100T / linux-amd64 / GOMAXPROCS=8 source/exactness/race/CLI/vet, frozen-core/pilot, structural-copy, fingerprint-mutation and live-CPU parser gates passed. Both qualified full-bank references match the frozen Build101 counters; all four prefix captures agree; **all 20 complete-bank replays are exact**. Four CPU profiles, all rows and unchanged source/corpus hashes are retained. No payload/HMAC or nine-photo physical gate was run or claimed.

The two prefixes produce 593 identical Gen4 tasks on B/mild and 1,902 on B/angle. Every replay preserves the qualified ordered banks of 937 / 6,198 candidates and total geometry evaluations of 79,259 / 334,857 respectively. The public capture/replay rows freeze full task/raster/bank fingerprints, not just outcome counts.

### Untouched-kernel timing controls

Only the 16 unprofiled slots contribute here. Every A/B slot uses identical code/data/pointers; negative apparent B/A change means B happened to be faster, not an optimization.

| case | pool mean ms | pool min–max ms | range / mean | population CV | apparent B/A mean change |
|---|---:|---:|---:|---:|---:|
| B/mild | 20,616.750 | 20,304–21,008 | 3.414699% | 0.993298% | -1.157366% |
| B/angle | 129,304.875 | 126,041–135,269 | 7.136622% | 2.687507% | +0.533301% |

| case | pair 1 A→B | pair 2 B→A | pair 3 B→A | pair 4 A→B |
|---|---:|---:|---:|---:|
| B/mild apparent B/A | -2.064441% | -1.238853% | +1.065390% | -2.346725% |
| B/angle apparent B/A | +1.026650% | -0.123458% | -0.343457% | +1.628050% |

The whole-session range must not be confused with a time-adjacent pair difference or a confidence bound for historical experiments. Both tables retain all samples, including B/angle's two approximately 135 s values between approximately 126–129 s values. A/A apparent signs vary despite exact work. These observations demonstrate within-code variation on this host; they do not assign the cause of an earlier Build100 slowdown or revise its historical guard.

Live total process CPU averages are 161,777.500 ms / 1,019,098.750 ms on B/mild / B/angle, with ranges relative to mean of 3.362643% / 6.754988%. Total CPU is the derived sum of Linux user and system CPU, measured at 100 ticks/s (10 ms resolution). Total CPU/wall spans 7.794–7.887 / 7.831–7.895 CPU-s per wall second on eight workers. Summed worker time and live CPU vary alongside wall time; worker occupancy alone is not CPU utilization. Process counters include small boundary-read costs.

B/mild allocates 17,476,096–17,480,720 bytes per timing window, with no completed GC cycles or observed GC pause. B/angle allocates 130,815,344–130,821,088 bytes and completes two GC cycles in every timing replay; summed pause ranges from 0.629 to 5.703 ms. Equal completed-cycle counts do not prove identical GC scheduling or exclude concurrent GC work. These pauses alone do not explain the seconds of wall-time spread. Peak memory, cache locality, CPU frequency and thermal state were not measured; no causal attribution to those factors is made.

### Separate CPU profiles

The four profiled rows remain separate: mild A/B pool wall 20,711 / 20,562 ms; angle B/A 131,124 / 133,687 ms. Their whole-profile sampled CPU totals are 160.33 / 160.60 / 1,019.91 / 1,034.52 s, with live process CPU retained alongside the raw samples. They are not added to timing means or used as candidate qualification.

The same projective readers dominate. In the final B/angle A profile, historical `readProjectiveBlockValue` is 512.99 cumulative sampled CPU-s and the exact Build84 reader is 514.78 CPU-s. Pprof also assigns 5.39–7.66% flat samples to `runtime.asyncPreempt` across the four profiles. This is a profiling observation, not proof that this fraction is removable runtime overhead or that preemption caused the historical adverse stage timing. Sampling/stack attribution and actual runtime work need a separate controlled experiment.

## Decision and next bounded study

**Close Build102 as observability PASS. Keep Build98 qualified, Build93 previous and Build100 closed / not promoted.** There is no new runtime implementation, default runtime option, physical campaign or promotion. Detached replay shows that identical work can produce percent-scale apparent differences, and that whole-session variation differs from paired variation. Future candidate protocols should include contemporaneous identical-work controls and fix their interpretation prospectively; these data do not retrospectively adjust historical thresholds.

The next bounded proposal is a research-only preemption-attribution control on the same frozen replay: compare default behavior with a separately launched diagnostic `GODEBUG=asyncpreemptoff=1` process, keeping exact bank/raster/task gates, live CPU, balanced ordering, hard timeouts and separate CPU sampling. Qualified Go's local `runtime/extern.go` documents that this debug option disables signal-based asynchronous preemption and can delay GC/scheduling in some loops. It is therefore a diagnostic condition, not a recommended production option. The control must test whether wall/live CPU observations support the profile attribution; it must not equate profile sample shares with a guaranteed speedup. No Build103 or runtime configuration change is implemented here.

## Frozen evidence

`docs/qualified-baselines/build102-captures.tsv` retains all four captures; `build102-replays.tsv` retains all 20 timing/profile rows and adds the explicitly derived `process_cpu_ms = user + system`; `build102-placebo-pairs.tsv` retains all eight A/A pair contrasts; `build102-variation-summary.tsv` retains six metrics for each case. Profiled rows are excluded from pair/variation statistics. The four tables and private `build102-frozen-gen4-replay.zip` are frozen by `build102-evidence.sha256`. Raw profiles, binary, logs, source/corpus manifests and acquisition images remain excluded from public commits.
