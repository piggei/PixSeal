# Format-v4 Build89 — full qualified Build84 deep-pipeline profiling

Build89 is an **observability-only** checkpoint over the current qualified Build84 smartphone baseline. It closes the exact-portable reader micro-optimization branch after Build88 and restores the active deep runtime to Build84.

## Why Build89

Build88 passed exactness but failed its benchmark-first performance criterion on the qualified Go 1.26.0 / Intel Core i3-13100T host. Same-session five-run means were:

| reader | front | mild | angle |
|---|---:|---:|---:|
| Build84 qualified | 1415.4 ns | 1434.0 ns | 1450.8 ns |
| Build88 direct RGB | 1529.2 ns | 1593.8 ns | 1611.8 ns |
| regression | +8.0% | +11.1% | +11.1% |

The isolated prepared RGB fetch also failed to improve: the qualified slice form averaged 255.68 ns/block and the direct-load form 257.46 ns/block. Build88 therefore did not reach the private physical gate and is non-promoted.

Build80/81/82/86/88 together show that further source-level rewrites of the exact projective-reader kernel are unlikely to be a productive next step without fundamentally new evidence. Build89 moves the measurement scope back to the complete qualified deep-recovery pipeline.

## Runtime semantics

Build89 delegates deep recovery directly to `experimentalV4PhoneBuild84Recover`. It does not use the Build86 or Build88 candidate readers. No geometry, ranking, ordering, threshold, qualification, protected-data decode, ECC/Hamming, whitening/HMAC domain, strength or wire-format behavior changes.

Build89 adds only **derived accounting after Build84 returns**. It does not add new `time.Now()` probes inside the hot geometry, qualification, sampling, list-decoder or HMAC loops. The authoritative timers are the ones already present in the qualified Build84 telemetry chain.

## Measured phases

For every retained case that reaches the Build84 deep fallback, Build89 reports:

- deep-recovery total wall time;
- geometry wall time;
- geometry pixel-plane preparation;
- Build47 freeze;
- prefix1;
- generation 2;
- generation 3;
- generation 4;
- accounted and residual geometry wall time;
- post-geometry pixel-plane preparation;
- held-out qualification;
- ordered protected-data decode wall time;
- accounted and residual total deep-recovery wall time;
- inherited physical sampling/list worker sums, candidate counts and exact semantic counters.

The retained nine-photo matrix is still checked semantically so the profiling run cannot silently drift away from qualified Build84 behavior. Those semantic checks are a profiling safety gate, **not a new qualification or promotion event**.

## Result and decision

Build89 completed the retained Go 1.26.0 nine-photo matrix with a 9/9 semantic PASS while running the unchanged qualified Build84 runtime. The decisive phase shares were:

1. Build47 freeze;
2. prefix1 / gen2 / gen3 / gen4 wall time;
3. qualification;
4. protected-data decode;
5. unaccounted residuals.

The result selected Build90: on B/mild, qualification consumed 61,193 ms of 149,504 ms deep recovery (40.93%), exceeding geometry at 58,976 ms and decode at 29,119 ms. B/angle remained geometry-dominated, with qualification only 6,867 ms. Build90 therefore tests ordered-parallel qualification over the already-frozen bank; Build84 remains the qualified baseline until benchmark and physical evidence justify promotion.

## Commands

Source/consistency gate:

```bash
make clean
make v4-build89-full-pipeline-test
```

Qualified-host retained-corpus profile:

```bash
make v4-build89-full-pipeline-profile
```

Artifacts are published atomically under:

```text
v4-phone private/build89-profile/
├── build89-metadata.txt
├── build89-full-pipeline.tsv
├── build89-full-pipeline.md
└── logs/
```

No separate Build89 physical qualification target exists: the profile itself executes the retained matrix only to measure the unchanged qualified Build84 runtime and verify semantic equivalence.
