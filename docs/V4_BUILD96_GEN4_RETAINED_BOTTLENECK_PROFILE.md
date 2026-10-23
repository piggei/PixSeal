# Build96 — Gen4 retained bottleneck profiling

Build96 is an **observability-only** checkpoint. It changes no production or experimental recovery algorithm. Build93 remains the current qualified smartphone baseline. Build95 remains available only as a reproducibility contrast and is explicitly **rejected / not promoted** after the retained AB/BA decision gate.

## Why Build96 exists

Build95 passed exactness/race and showed strong same-session local gains: the exact luminance reader was roughly twice as fast as Build84 and continuation4 improved by roughly forty percent. The retained paired pipeline did not reproduce that advantage strongly enough to justify ~121.6 MiB of additional memory. B/mild was effectively neutral; B/angle showed a small mean gain but strong AB/BA order sensitivity.

Build96 asks a narrower question: **where does the local reader gain disappear in the retained eight-worker generation-four workload?**

## Method

Only `phone-b-angle.jpg` is used because Build94 showed it is overwhelmingly geometry/gen4-bound and performs no protected decode. Two cycles are run in one process/session:

1. Build93 -> Build95
2. Build95 -> Build93

Each deep run receives its own CPU profile and runtime memory/GC snapshots. The harness records total/geometry/gen4 wall time, cumulative gen4 worker time and task distribution, luminance-plane preparation, allocation deltas, heap state, GC count/pause and the frozen FoldScore/projective-read workload.

Focused pprof reports are generated for generation4, continuation4, FoldScore, projective reader, `homography.mapPoint`, quad homography, single probe and pair stencil.

## Frozen semantic gate

Every row must remain:

- outcome REJECT
- geometry evaluations: 334857
- bank: 6198
- qualification evaluations: 6198
- qualified: 0
- FoldScore: 118004
- block reads: 151045120
- block success: 149370664
- block failed: 1674456

Build96 has no physical promotion gate because it changes no runtime.

## Commands

```bash
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build96-gen4-retained-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build96-gen4-retained-profile
```

Private output is staged atomically under `v4-phone private/build96-profile/`; failed runs are preserved under `build96-profile-failed/`.

## Decision rule for the next build

Do not select Build97 before reviewing both CPU and wall-time evidence. If Build95 reader/continuation CPU falls materially while gen4 wall barely changes, investigate retained parallel memory/cache pressure. If retained CPU itself does not reproduce the microbenchmark gain, close the luminance-plane branch and select a different structural gen4 target from the new pprof.
