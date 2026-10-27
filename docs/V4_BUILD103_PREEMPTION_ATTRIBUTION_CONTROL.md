# Build103 — fresh-process preemption-attribution control

**Closed / semantic and observability PASS, 2026-10-27. Diagnostic only. Current qualified Build98 / previous Build93; Build100 closed / not promoted.**

Build102's identical-work A/A experiment found percent-scale timing variability and 5.39–7.66% flat CPU samples assigned to `runtime.asyncPreempt`. Sample share does not establish removable work. Build103 tests one diagnostic runtime condition against explicit default behavior, using identical frozen Gen4 inputs and the unchanged Build102 replay kernel.

## Predeclared protocol and decision

1. In one default-preemption preparation process, obtain a qualified Build98 full-bank reference and capture the prefix at the original Gen3→Gen4 barrier for each retained B/mild and B/angle image. Check frozen Build101 counters and Build102 full task/raster/bank fingerprints. Serialize one private snapshot per case, including every private hypothesis field, original order, exact IEEE floating bits, raster, pilot, prefix metadata and expected non-timing telemetry. Verify its disk round trip. Snapshots contain acquisition pixels and remain private.
2. Compile one Go 1.26.0 test binary. Each measured replay starts in a **fresh process** and loads the same case snapshot. Replace inherited `GODEBUG` with exactly `asyncpreemptoff=0` (D, explicit default) or `asyncpreemptoff=1` (O, diagnostic). Fix GOMAXPROCS=8, GOGC=100, GOMEMLIMIT=off and GOTRACEBACK=single. Explicit GC and a five-second cooldown precede measurement. Startup, loading, input checks and final output checks remain outside the measured replay.
3. For each case run eight timing processes in **D O O D | O D D O** order, four per condition. Then run two additional CPU-profile processes, mild D→O and angle O→D. Keep all four profiled rows separate from timing means. Measure complete replay wall, original pool wall, summed worker duration, live Linux user/system CPU and allocation/GC deltas. CLK_TCK comes from getconf; process counters include small boundary-read overhead. Worker duration is occupancy, not CPU utilization.
4. Every process must preserve the complete ordered bank, exact task/raster/metadata fingerprints, all non-timing geometry/fetch/freeze telemetry, evaluation/seed/worker counts and original snapshot checksum. Require **20/20** exact replays, all four CPU profiles and stable source/corpus hashes for semantic/observability closure. A public synthetic fresh-process smoke test exercises both runtime conditions without private data. External child limits are 1,800 seconds for preparation and 600 seconds for a replay, with an external 7,200-second study limit and process-group cleanup; they do not depend on Go's in-process preemption.
5. **Support a material diagnostic effect only if O improves mean Gen4 pool wall and live total process CPU by at least 5% on both cases, and every time-adjacent paired wall/CPU comparison is favorable.** Otherwise classify attribution as mixed/unsupported. Record all eight pairs and both condition means. This is a prospectively fixed descriptive rule for this one session, not a statistical significance test or production-promotion threshold. Disappearance of asyncPreempt samples alone is neither a speedup nor proof of a pure sampling artifact.

Go 1.26.0's installed `runtime/extern.go` documents that this flag disables signal-based asynchronous preemption, can delay scheduling/GC in some loops, and changes conservative stack scanning. Thus even favorable results concern the entire diagnostic runtime condition; they do not isolate trampoline instructions or explain a historical Build100 slowdown. Fresh process address/layout variation, one host/session, unmeasured frequency/thermal/cache behavior and profiling attribution limit causal claims. No protected decoding, payload/HMAC or nine-image physical qualification is run or claimed. Production configuration, qualified sources and historical Build100 guards remain unchanged.

## Reproduction

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build103-preemption-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build103-preemption-profile
```

Default private output: `v4-phone private/build103-preemption-control`. Existing outputs are refused; failed evidence is retained uniquely. Raw snapshots/profiles/logs/binary remain private. Only reviewed numeric rows and fingerprints will be versioned after completion.

## Results — complete semantic/observability PASS

Qualified Go 1.26.0 / i3-13100T / Linux amd64 / eight workers. All inherited source/exactness/race/CLI/vet, frozen-core/pilot and structural checks passed, as did lossless serialization, mutation rejection and actual default/off fresh-process smoke tests. Both qualified references, both snapshot round trips and **all 20 complete-bank replays** pass exact task/raster/ordered-bank and non-timing telemetry checks. B/mild retains 593 tasks, 937 bank candidates and 79,259 geometry evaluations; B/angle retains 1,902 tasks, 6,198 candidates and 334,857 evaluations. Source/corpus hashes and both snapshot checksums stayed fixed throughout retained execution; all four profiles are complete. No protected decode, payload/HMAC or physical campaign is claimed.

An earlier preparation attempt was interrupted by sandbox/approval infrastructure failure (read-only registry / ENOSPC), before any measured replay. Its recovered snapshot was zero bytes. Logs are retained separately and that attempt is excluded entirely. The complete session restarts the same predeclared protocol and criteria after environment restoration; no interrupted data enter a mean. A second external execution-session interruption occurred during CPU-only run19. All sixteen timing rows and both mild profiles were already complete and are preserved. The partial profile/log remain archived as `.interrupted`; the two angle profiles resumed with the same verified binary and snapshots in the original O-D order. They are a resumed profiling period and remain excluded from timing statistics.

### Unprofiled fresh-process timing

Negative change means diagnostic-off is faster. Only the 16 timing processes contribute; every row is retained.

| case | default pool ms | diagnostic pool ms | change | default live CPU ms | diagnostic live CPU ms | change |
|---|---:|---:|---:|---:|---:|---:|
| B/mild | 20,647.250 | 19,727.250 | -4.46% | 158,260.000 | 151,330.000 | -4.38% |
| B/angle | 127,833.250 | 124,495.750 | -2.61% | 990,280.000 | 964,655.000 | -2.59% |

| case | pair / order | pool change | live CPU change | both favorable |
|---|---|---:|---:|---|
| B/mild | 1 / D-O | -12.45% | -11.65% | true |
| B/mild | 2 / O-D | -4.01% | -3.42% | true |
| B/mild | 3 / O-D | +1.13% | +0.05% | false |
| B/mild | 4 / D-O | -0.61% | -1.02% | true |
| B/angle | 1 / D-O | +2.45% | +2.11% | false |
| B/angle | 2 / O-D | -0.56% | +0.12% | false |
| B/angle | 3 / O-D | -8.07% | -8.27% | true |
| B/angle | 4 / D-O | -3.63% | -3.66% | true |

The public means table also retains complete-replay/worker means, every condition's min/max and population CV. These are descriptive one-session values, not independent host replications or confidence bounds. Early-slot variation is retained rather than removed after inspection. CPU totals are Linux live user + system counters, quantized at 100 ticks/s (10 ms); small boundary-read overhead remains included. Summed worker durations measure occupancy.

### Separate CPU attribution

The following four additional executions are excluded from timing means. Whole-profile samples and live CPU are distinct measures. The final two angle profiles (runs19/20) were collected after the external interruption and are explicitly labelled `resumed_cpu`; they are not contemporaneous replacements for any timing row.

| case | condition | profile pool ms | live CPU ms | sampled CPU s | asyncPreempt flat samples |
|---|---|---:|---:|---:|---:|
| B/mild | default | 19,293 | 148,840 | 148.46 | 2.05% |
| B/mild | async_off | 19,514 | 150,900 | 150.61 | 0.00% |
| B/angle | async_off | 114,620 | 892,570 | 890.40 | 0.00% |
| B/angle | default | 114,856 | 896,220 | 894.21 | 2.40% |

Default profiles assign 2.05% (mild) / 2.40% (angle) flat samples to `runtime.asyncPreempt`; neither diagnostic profile contains samples for that symbol. This disappearance accompanies mixed timing/CPU comparisons and does not satisfy the predeclared material-effect rule.

The full unfiltered profile summaries, both projective readers' cumulative samples, allocation/GC counters and original profiles are retained. A zero asyncPreempt entry means no samples in that profile, not a general guarantee about all runtime work. Removing its sampled attribution does not itself quantify removable elapsed/CPU cost. The diagnostic flag also affects GC scheduling and stack handling, while process address/layout, frequency and cache state are not controlled or causally isolated.

## Decision

**Close semantic/observability PASS; predeclared material diagnostic effect: MIXED/UNSUPPORTED.** The exact prospective conjunction is recorded for each case in `build103-decision.tsv`; it is not relaxed. Favorable individual values remain evidence about this diagnostic condition, not a production optimization, unique historical cause or a revision of Build100's guard. Current qualified Build98 / previous Build93 remain unchanged.

Future work should select a bounded kernel hypothesis from the retained reader profiles and start with exact public tests and benchmarks, keeping contemporaneous identical-work controls in any later retained comparison. The sampled asyncPreempt share alone does not select a runtime setting or authorize promotion.

## Post-study source-test hardening

After retained execution and its source checks, the synthetic fresh-process smoke test was made independent of a caller starting with fewer CPUs: it explicitly matches the child GOMAXPROCS=8 and restores the parent setting afterwards. The Build103 source tests pass with GOMAXPROCS=1. Only this public test changed; the measured kernel, replay, snapshot codec, harness and summarizer remain byte-identical. The private archive preserves the original test/source manifest/binary and a separate final-source manifest, distinguishing measured sources from this test-only fix. No retained row or interpretation was changed.

## Frozen evidence

Six reviewed tables are versioned under `docs/qualified-baselines`: `build103-captures.tsv`, `build103-runs.tsv` (all 20 rows, plus explicitly derived user+system `process_cpu_ms`), `build103-pairs.tsv` (all eight pairs), `build103-means.tsv` (four metrics per case), `build103-decision.tsv` (both case conjunctions), and `build103-cpu-attribution.tsv` (four separate profiles). The reproducible summarizer reads raw tables without overwriting them. Raw 32-column timing rows and reviewed public copies remain distinct. `build103-evidence.sha256` freezes these tables and the private `build103-preemption-control.zip`; snapshots contain acquisition pixels and remain private with binary, logs and profiles. The interrupted preparation is archived separately and its checksum is retained.
