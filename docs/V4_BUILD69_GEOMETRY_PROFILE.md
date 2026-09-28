# Build69 geometry-generation profiling

Build69 is a **completed observability-only research snapshot** over the qualified Build68 smartphone baseline. Its Go 1.26.0 retained-corpus gate passed exact semantic equivalence and identified strong single-seed load imbalance. Build68 remains the current qualified smartphone baseline.

## Why this build exists

Build67 showed that deep recovery was dominated by geometry and qualification. Build68 removed the qualification-side repeated `pixelPlane` materialization and reproduced about a 5.02x qualification speedup across two physical runs while preserving exact semantics. After that change, geometry is the dominant remaining stage; on difficult negative `B/angle`, essentially the entire deep fallback is geometry generation.

A static audit found no comparably obvious safe optimization. The geometry path already parallelizes the selected independent seeds and reconstructs their results by original seed index. A second plane materialization exists inside the historical freeze path, but one plane build is only a few hundred milliseconds and is far too small to explain hundreds of seconds of geometry time. The justified next step is therefore profiling, not algorithm redesign.

## Frozen semantics

Build69 must not change:

- Format-v4 wire format, encoder, public pilot or strength 48;
- Build64 geometry score, search bounds, move families, seed count, candidate retention or bank order;
- Build65 seed worker count/reassembly semantics;
- the complete bank freeze barrier;
- Build68 plane-reuse qualification and all qualification thresholds;
- Build66 ordered parallel protected-data decode and first-logical-HMAC semantics;
- ECC/Hamming, whitening or HMAC domains;
- fallback order or historical scanner/Format-v3 paths.

The Build69 wrapper calls the unchanged per-seed Build64 search and stores each result in its original seed slot before concatenation, exactly as Build65 does. Timing is observational only.

## Geometry telemetry

For every deep-fallback case Build69 reports:

```text
build69-geometry-profile: attempted=true \
  plane-prep-ms=... \
  freeze-ms=... \
  seed-wall-ms=... \
  seed-worker-ms=... \
  seed-min-ms=... seed-median-ms=... seed-max-ms=... \
  seed-min-evals=... seed-max-evals=... \
  seed-min-bank=... seed-max-bank=...
```

Interpretation:

- `plane-prep-ms`: first reusable geometry plane construction;
- `freeze-ms`: historical Build47 freeze/seed-source stage;
- `seed-wall-ms`: wall time for the parallel seed-search region;
- `seed-worker-ms`: sum of all individual seed durations;
- `seed-min/median/max-ms`: seed runtime dispersion;
- `seed-min/max-evals`: per-seed search-work dispersion;
- `seed-min/max-bank`: per-seed retained-output dispersion.

The existing Build64/65/66/68 telemetry remains the semantic authority.

## Physical gate

Run on the qualified Go 1.26.0 host:

```sh
make clean
make v4-build69-phone-geometry-profile-test
make v4-build69-phone-physical-test
```

The gate writes:

```text
v4-phone private/build69-diagnostics/build69-phone-geometry-profile.tsv
v4-phone private/build69-diagnostics/build69-phone-geometry-profile.md
```

It requires the same nine-photo outcomes and exact deep logical counters as the qualified lineage. The optional Build68 TSV is used only for informational wall-clock ratios.

## Physical result — 2026-09-28

The retained nine-photo gate passed exact semantics. The maximum seed consumed 96.2–99.7% of the seed-parallel wall time in every deep-recovery case, while median seed durations stayed below one second:

```text
control/front  freeze=66981 ms  seed-wall=40108 ms   max-seed=38572 ms
control/mild   freeze=75176 ms  seed-wall=360285 ms  max-seed=358584 ms
control/angle  freeze=61116 ms  seed-wall=17667 ms   max-seed=17182 ms
B/mild         freeze=56083 ms  seed-wall=169014 ms  max-seed=168398 ms
B/angle        freeze=54826 ms  seed-wall=695422 ms  max-seed=693395 ms
```

The slowest seed also accounts for 73922/79259 total geometry evaluations on `B/mild` and 267775/334857 on `B/angle`; because total geometry evaluations still include the freeze stage, its share of seed-only work is even larger. This closes Build69 as **seed-work dominated with extreme imbalance** and selects Build70 dominant-seed genealogy profiling.

## Decision rule

Build69 itself is not a promotion candidate. The physical profile classifies the remaining cost as case 2 below:

1. **freeze-dominated** — investigate deterministic reuse/caching within the frozen public geometry stage;
2. **seed-work dominated with strong imbalance** — **observed**; Build70 profiles the dominant branch before any scheduling/decomposition change;
3. **seed-work dominated and broadly balanced** — profile inside the unchanged per-seed search before changing implementation;
4. **other/fixed overhead** — measure that component directly before optimization.

Do not use HMAC, payload, ECC, oracle/reference geometry or known correct messages to rank, prune, stop or retune geometry.
