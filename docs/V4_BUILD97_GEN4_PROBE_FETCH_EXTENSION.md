# Build97 — exact Gen4 probe-fetch extension

Build97 is a **benchmark-first research candidate** selected from the completed Build96 retained profile. Build93 remains the current qualified smartphone baseline. Build95 remains rejected/not promoted.

## Build96 decision

Build96 preserved the frozen B/angle semantics in all four AB/BA rows. It showed that Build95's luminance-plane reader does **not** reproduce its public microbenchmark gain in the retained eight-worker workload: continuation-reader cumulative CPU improved only from 508.29 -> 486.96 CPU-s in cycle 1 and 497.92 -> 445.24 CPU-s in cycle 2, while Gen4 wall time was actually slower in cycle 1 (132283 -> 140949 ms) and only 1.16% faster in cycle 2 (130469 -> 128962 ms). Build95 also retained its 127,475,200-byte image-sized allocation. The luminance-plane branch therefore stays closed.

The same Build96 profiles provide new evidence for a different bounded target: within qualified Build93 Gen4, the historical generic `samplePlaneLuminance` path used by **single4/pair4** remains roughly 399 CPU-s cumulative, while continuation4 already uses the qualified Build84 exact RGB reader. Build97 therefore does not invent a new reader and does not allocate a luminance plane. It extends the already-qualified Build84 reader/FoldScore implementation to the two remaining Gen4 probe stages only.

## Exact changed scope

Only generation-four scoring calls change:

- `single4`: historical `Build53SingleProbe` -> exact Build97 wrapper using `Build84FoldScoreFetch`;
- `pair4`: historical `Build53PairStencil` -> exact Build97 wrapper using `Build84FoldScoreFetch`;
- `continuation4`: unchanged Build84 implementation.

Prefix1, generations 2/3, quad construction, proposal thresholds, stable pair ordering, geometry freeze, Build90 qualification, Build93 protected-residue reuse, candidate order, ordered commit and first-logical-HMAC remain unchanged. There is **zero image-sized additional memory**.

## Qualification discipline

Build97 must pass exact single-probe, pair-stencil and full-generation-four equality, race testing and a same-session public benchmark before any retained AB/BA profile is allowed. A local sanity benchmark on the packaging host is intentionally non-authoritative; only the qualified Go 1.26.0 i3-13100T result decides whether Build97 proceeds.

If the full Generation4 benchmark is not materially and directionally faster, Build97 is closed without retained or physical testing. If it is positive, the next gate is a contemporaneous Build93/Build97 retained B/angle AB/BA profile. No physical 9/9 gate is justified before that.

## Qualified-host public benchmark review — retained paired profiling authorized

The Go 1.26.0 / linux-amd64 benchmark on the qualified Intel Core i3-13100T host passed source/exactness/race and produced a directionally positive full-Generation4 result. Across five samples, the historical Build84 Generation4 path averaged 28.314 ms/op and Build97 averaged 26.427 ms/op, a 6.66% mean reduction. Single-probe averaged 3.492 -> 3.325 ms/op (-4.76%) and pair-stencil averaged 24.357 -> 22.970 ms/op (-5.69%). Allocation counts remained unchanged in the compared full-Generation4 benchmark (6552 B/op, 3 allocs/op).

This evidence is strong enough to authorize retained same-session comparison, but it is not promotion evidence. The next gate is a contemporaneous Build93/Build97 AB/BA profile on B/mild and B/angle using the same decoded retained work image per pair. Build97 remains research-only and the normal `v4-extract-phone` remains Build93.

Run:

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build97-gen4-probe-fetch-paired-profile
```

Promotion review requires exact semantic equality plus a material and directionally stable Gen4/geometry gain in both orderings. Do not run the 9/9 physical promotion gate before this retained evidence is reviewed.

## Retained gate telemetry remediation

The first qualified-host retained invocation stopped correctly on the first B/mild Build97 row. Recovery itself authenticated the expected payload and preserved every semantic counter, but the new single4/pair4 calls were merged into the historical Build84 `Fetch*` family. That changed the published counters from `16,496 / 21,114,880` to `41,776 / 53,473,280` FoldScores/block reads. The exact delta was `25,280` FoldScores and `32,358,400` reads (`25,280 * 1,280`), proving that the failure was accounting contamination from the newly instrumented probe stages rather than a geometry/output mismatch.

The harness was right to fail: qualified telemetry is append-only. Build97 remediation therefore does **not** relax the expected values and does not restore counters after the fact. It keeps Build84 `Fetch*` continuation4-only and adds a separate candidate-local `ProbeFetch*` family for single4/pair4. The paired gate still requires the historical counters to match Build93 row-for-row, while also checking that Build97 probe counters are positive and internally consistent. Build93 reports zero probe counters.

The failed first attempt remains non-performance evidence. Its partial timing looked promising, but no decision is allowed until the corrected complete AB/BA gate passes.

## Final retained decision — exact but not promoted

The corrected telemetry gate completed two independent Build93/Build97 AB/BA runs with exact semantic equality. Historical Build84 continuation counters remained unchanged, and candidate-local `ProbeFetch*` counters accounted separately for single4/pair4 work.

The predeclared promotion criterion was not met. The first run showed a strong B/mild Gen4 gain, but the independent second run averaged `17,649 -> 17,783.5 ms` for B/mild Gen4, a 0.76% regression. Across all four B/angle pairs, pair-normalized Gen4 was 0.98% slower and pooled means were 1.27% slower, with the sign changing by order/run. End-to-end totals were effectively neutral. Build97 is therefore closed without a physical 9/9 gate. Build93 remains qualified.

Non-sensitive decision rows are frozen in `docs/qualified-baselines/build97-decision.tsv`. The next benchmark-first candidate is Build98, which targets the separate Generation3 generic-reader cost identified by the completed Build94 pprof.
