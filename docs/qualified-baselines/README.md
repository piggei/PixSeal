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

Build98 promotion is based on exact source/race equivalence, the public complete-Generation3 benchmark, contemporaneous AB/BA changed-phase evidence, and two independent retained 9/9 semantic physical gates. Build93 is the previous qualified baseline. Historical Build93 promotion was based on exact source/race equivalence, the public residue-reuse benchmark, and two independent retained 9/9 semantic physical gates with repeatable B/mild decode improvement. The later production-entry-point matrix is integration evidence only and is kept separate from promotion-performance evidence. Build90 remains an earlier qualified reference. Cross-run absolute geometry timing is retained rather than cherry-picked.

## Experimental decision evidence

`build95-decision.tsv` and `build97-decision.tsv` freeze the non-sensitive retained AB/BA rows used to reject those candidates. They are evidence about experimental candidates, **not** qualified baselines. `build98-benchmark-summary.tsv` and `build98-paired-summary.tsv` are the pre-physical decision evidence; `build98-phone-performance.tsv` records the two completed physical gates that qualify Build98.
