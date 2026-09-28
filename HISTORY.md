## v0.3.0-build73 — ordered generation-three parallel candidate

The completed Build72 physical profile preserved Build71 semantics and localized the new geometry critical path to the third prefix generation. `B/angle` spends 194,940 / 293,120 ms cumulative prefix worker time in `pair3 + continuation3 + sibling3`, and its deterministic dominant prefix seed spends 149,514 / 203,169 ms there. `B/mild` shows the same pattern at 41,319 / 89,863 ms and 41,319 / 72,021 ms respectively. Build73 therefore leaves the search unchanged but moves the ordered parallel barrier to `sibling2`: all independent generation-three subtrees run in one bounded pool, are committed in original order, and feed the already-qualified Build71 generation-four pool. Build71 remains the current qualified smartphone baseline until Build73 passes reproducible physical qualification.

## v0.3.0-build72 — prefix-stage observability snapshot

Build71 qualification leaves the prefix through `sibling3` as the measured geometry critical path, especially on `B/angle`. Build72 makes no scheduling or algorithm change. It times the exact fourteen public-only prefix stages, aggregates stage evaluations/worker time across all seeds, and records the deterministic maximum-evaluation prefix seed plus its per-stage genealogy. Build71 remains the current qualified smartphone baseline.

## v0.3.0-build71 — qualified smartphone baseline

Two independent retained Go 1.26.0 physical runs reproduce the complete Build68 semantic matrix while confirming the Build71 generation-four scheduling gain. `B/mild` remains exactly 79,259 geometry evaluations / 937 bank / 935 qualified / 691 logical decode candidates / 2,120,047 logical frames and authenticates the expected payload; `B/angle` remains 334,857 / 6,198 / 0 and rejects. Full-matrix elapsed is 1,233,478 and 1,263,623 ms versus the retained Build68 reference of 1,902,847 ms. Build71 is therefore promoted as the current qualified smartphone baseline. Build68 remains the historical qualification-plane-reuse milestone, Build66 the historical ordered-parallel decode milestone, and Build64 the deep-recovery semantic reference.

## v0.3.0-build71 — ordered generation-four parallel performance candidate

Build70 physical profiling passed the complete retained nine-photo semantic-equivalence gate and identified a deterministic public-only scheduling opportunity. The dominant seed contributes 73922/79131 seed evaluations on `B/mild` and 267775/334749 on `B/angle`. Its fourth generation (`single4 + pair4 + continuation4`) contributes 41776/73922 evaluations on `B/mild` and 196641/267775 on `B/angle`, with 593 and 1363 independent `sibling3` inputs respectively. Build71 preserves the exact Build64/68 search and final bank order but splits geometry scheduling into two deterministic barriers: parallel seed prefixes through `sibling3`, then one bounded global pool for the frozen fourth-generation subtrees, followed by original-order commit. Build68 remains the qualified smartphone baseline until Build71 passes the physical equivalence/performance gate.

## v0.3.0-build70 — dominant-seed genealogy profiling snapshot

The qualified Go 1.26.0 Build69 physical gate reproduced Build68 semantics exactly and established that the remaining geometry runtime is dominated by one seed branch, not by broadly distributed work. `B/mild` measured 168398 ms for the slowest seed inside a 169014 ms seed wall region; `B/angle` measured 693395/695422 ms; `control/mild` measured 358584/360285 ms. Build70 deliberately makes no geometry optimization. It reproduces the same search while recording freeze-vs-seed evaluation counts, the dominant seed identity, and per-stage evaluation/state genealogy. Build68 remains the current qualified smartphone baseline.

## v0.3.0-build69 — geometry-generation profiling result

The retained nine-photo Go 1.26.0 gate passed semantic equivalence. Deep cases showed extreme load imbalance: max-seed/seed-wall time was 96.2% (`control/front`), 99.5% (`control/mild`), 97.3% (`control/angle`), 99.6% (`B/mild`) and 99.7% (`B/angle`). The slowest seed also produced the complete final bank for `control/front`, `control/mild` and `B/mild`, and 5486/6198 states for `B/angle`. This evidence redirects Build70 from an optimization attempt to finer-grained genealogy profiling.


