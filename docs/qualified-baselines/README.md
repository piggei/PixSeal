# Qualified smartphone performance references

These files freeze source-controlled numeric references without shipping the private physical images, secret key, payload files, or raw diagnostic logs.

- `build76-phone-performance.tsv` — historical Build76 two-run reference.
- `build84-phone-performance.tsv` — earlier qualified Build84 two-run command-matrix reference.
- `build90-phone-performance.tsv` — previous qualified Build90 evidence: all three full semantic physical runs are retained, including the host-degraded run, plus the contemporaneous Build91 paired means used for its promotion.
- `build91-build84-build90-paired.tsv` — exact authoritative rows from the completed Build91 AB/BA paired physical profile.
- `build93-phone-performance.tsv` — previous qualified Build93 evidence: two independent retained 9/9 physical PASS runs, protected sampling/derivation counts, and frozen B/mild promotion comparisons.
- `build93-production-integration.tsv` — completed Go 1.26.0 post-promotion wiring confirmation through normal `v4-extract-phone`; nine semantic rows plus non-sensitive entry-point/toolchain/executable-hash metadata.
- `build94-profile-summary.tsv` — non-sensitive low-overhead timing and key pprof attribution from the completed Build94 profile of unchanged qualified Build93; this selects the bounded Build95 benchmark candidate but does not change the qualified baseline.
- `build95-decision.tsv` — retained AB/BA evidence that rejected Build95; this is candidate decision evidence, not a qualified baseline.
- `build96-retained-profile.tsv` — observability-only Build93/Build95 retained CPU/memory rows that selected the bounded Build97 hypothesis.
- `build97-benchmark-summary.tsv` — qualified-host public benchmark evidence that authorized retained Build97 comparison only.
- `build97-decision.tsv` — two independent corrected Build93/Build97 AB/BA runs and aggregate decision; Build97 is exact but not promoted.
- `build98-benchmark-summary.tsv` — qualified-host public benchmark evidence: complete Generation3 -9.40% mean with clean five-sample separation.
- `build98-paired-summary.tsv` — retained Build93/Build98 AB/BA semantic PASS; changed Generation3 faster in all four pairs.
- `build98-phone-performance.tsv` — current qualified Build98 evidence: two independent complete 9/9 physical PASS runs with frozen historical and Generation3 telemetry.

- `build98-production-integration.tsv` — completed normal-command Go 1.26.0 9/9 integration evidence; archived semantic rows and non-sensitive metadata verified on 2026-10-23.

- `build99-profile-summary.tsv` — completed observability-only timing and separate CPU attribution for unchanged qualified Build98, 2026-10-23. No new baseline or cross-session speedup is claimed.

- `build100-benchmark-summary.tsv` — public complete-freeze candidate samples; -48.43% mean.
- `build100-paired-summary.tsv` — all 16 exact-semantic retained rows from two independent Build98/Build100 sessions.
- `build100-decision.tsv` — every paired timing comparison, including favorable deep/freeze results and adverse unchanged-Gen4 values. Build100 is closed / not promoted; these are candidate evidence, not a qualified baseline.

Build98 promotion is based on exact source/race equivalence, the public complete-Generation3 benchmark, contemporaneous AB/BA changed-phase evidence, and two independent retained 9/9 semantic physical gates. Build93 is the previous qualified baseline. Historical Build93 promotion was based on exact source/race equivalence, the public residue-reuse benchmark, and two independent retained 9/9 semantic physical gates with repeatable B/mild decode improvement. The later production-entry-point matrix is integration evidence only and is kept separate from promotion-performance evidence. Build90 remains an earlier qualified reference. Cross-run absolute geometry timing is retained rather than cherry-picked.

## Experimental decision evidence

`build95-decision.tsv` and `build97-decision.tsv` freeze the non-sensitive retained AB/BA rows used to reject those candidates. They are evidence about experimental candidates, **not** qualified baselines. `build98-benchmark-summary.tsv` and `build98-paired-summary.tsv` are the pre-physical decision evidence; `build98-phone-performance.tsv` records the two completed physical gates that qualify Build98.

## Build101 observability closure

`build101-runs.tsv` retains all 12 exact-semantic runs with separate timing/attribution labels; `build101-stages.tsv` contains all 26 observed stage rows. `build101-timing-pairs.tsv` and `build101-timing-means.tsv` use only the eight untouched-API timing runs. `build101-gen4-attribution.tsv` combines observed Gen4 CPU, worker, allocation/GC and scheduler evidence for the four separate instrumented runs. CPU values come from stage-labelled sampling; runtime CPU-class estimates can remain stale between STW updates. Worker occupancy is not itself CPU utilization. All raw profiles, logs, binary and corpus remain private; `build101-evidence.sha256` freezes the private archive and these five public tables. See `../V4_BUILD101_RUNTIME_EFFECTS_OBSERVABILITY.md` for limits and the unchanged Build100 non-promotion decision.

## Build102 frozen Gen4 A/A closure

`build102-captures.tsv` retains four exact prefix captures from two qualified bank references. `build102-replays.tsv` retains all 20 frozen bank/raster/task exact replays with separate timing/CPU labels and the derived total `process_cpu_ms = user + system`. Only the 16 timing rows feed `build102-placebo-pairs.tsv` (all eight pairs) and `build102-variation-summary.tsv` (six metrics per case). A/B are identical-code/data measurement labels, not runtime variants. Session ranges/CVs and paired differences are distinct descriptive quantities, not historical confidence bounds. `build102-evidence.sha256` freezes these four public tables and the private raw archive. See `../V4_BUILD102_FROZEN_GEN4_REPLAY.md` for source/semantic gates, profile-attribution limits and unchanged qualification decisions.
