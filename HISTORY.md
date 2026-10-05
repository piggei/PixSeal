## Build83 — projective sampler CPU profiling

Build82 closed with two independent Go 1.26.0 semantic PASS runs but no repeatable performance gain: 933,537 / 729,213 ms full-matrix elapsed (831,375 ms mean) versus 758,102 ms for qualified Build76. Build83 therefore restores the active deep phone runtime to Build76 and performs no algorithmic optimization. It adds only a public deterministic benchmark/pprof harness for the historical projective reader and its mapPoint, luminance and DCT components, retaining Build82 solely as an exact comparator. Build83 timing is research evidence, never a semantic gate.

## Build82 — specialized exact inline block sampler performance candidate

Build81 closed as a semantic success but not a performance promotion. Its two independent Go 1.26.0 physical runs both passed the complete 9/9 matrix and preserved the deep workload exactly, but full-matrix elapsed was 911,984 and 798,601 ms (855,293 ms mean) versus the Build76 qualified mean 758,102 ms. The faster Build81 run made B/mild gen4 essentially equal to Build76 and B/angle gen4 about 1.8% faster, but the effect did not repeat. Build76 therefore remains qualified.

The Go 1.26.0 inliner report then confirmed the next implementation hypothesis: `samplePlaneLuminance` has cost 348, `homography.mapPoint` 88 and `readProjectiveBlockValue` 263, all above the inline budget 80. Build82 deliberately tests only the largest helper first. Inside continuation4 FoldScore block reads it manually incorporates the exact `samplePlaneLuminance` arithmetic and hoists immutable plane/DCT references, while retaining the historical `homography.mapPoint` call and every float64 operation/order. No LUT, cache, bounding rectangle, fallback, geometry pruning or protected-data change is introduced.

## Build81 — exact luminance LUT experiment (semantic PASS / non-promoted)

Build80 physically preserved the qualified B/mild and B/angle geometry counters but was slower: the nine-photo matrix took 806,849 ms versus the Build76 two-run mean 758,102 ms. The cause is structural: each cached block first maps all 64 destination samples to build a source bounding rectangle and, on fallback, the original reader maps the same 64 samples again. B/angle had zero cache hits, 149,370,664 valid fallbacks and 1,674,456 failed reads, so the candidate was rejected despite semantic equality. Build81 removed all spatial-cache logic and replaced only `.299*r`, `.587*g` and `.114*b` with exact 256-entry float64 product tables inside continuation4 FoldScore sampling. Both retained physical runs preserved exact semantics, including B/mild `79259 / 937 / 935 / 691 / 2120047` and B/angle `334857 / 6198 / 0`, but did not show a repeatable speedup. Build81 is archived as a valid negative performance experiment and Build76 remains qualified.

## Build80 — exact local luminance cache performance candidate

Build79 passed the retained Go 1.26.0 9/9 semantic-equivalence gate and measured projective block reads at 1,041 / 1,071 ms (97.20%) of sampled B/mild FoldScore time and 7,236 / 7,424 ms (97.47%) on B/angle. Detailed replay indicated bilinear luminance sampling is the largest subcomponent. Build80 then changed only continuation4 FoldScore block sampling: a FoldScore-local scratch buffer cached exact float64 RGB→luminance values for mapped source rectangles up to 196 pixels and otherwise fell back to the unchanged reader. Exact semantic regressions and the physical matrix passed, but runtime regressed; Build80 is closed as rejected and Build76 remains qualified.

## Build79 — FoldScore kernel profiling

Build78 passed the retained Go 1.26.0 9/9 semantic-equivalence matrix and showed that continuation4 bookkeeping is negligible: FoldScore consumed 67,777 / 67,832 ms on B/mild and 500,211 / 500,652 ms on B/angle. Build79 therefore keeps the complete qualified Build76 recovery schedule and bank/order frozen and instruments only the FoldScore kernel used by continuation4. Exact counters cover every FoldScore call; a deterministic 1/64 sample times authoritative block reads, and one successful block per sampled FoldScore is replayed after the score has completed to separate `mapPoint`, bilinear luminance sampling and DCT accumulation. Build79 is diagnostic-only; Build76 remains qualified.


## Build78 — continuation4 internal profiling

Build77 localized the largest remaining generation-four worker component to continuation4 on B/angle while showing that the outer gen4 pool is already saturated. Build78 keeps Build76 qualified semantics frozen and profiles the exact Build55 continuation loop internally. It is diagnostic-only and cannot be promoted from timing.

Build77 is an observability-only snapshot over the qualified Build76 baseline. Build76 left generation four as the dominant B/angle wall region (~110.8 s) while already reaching ~7.98x effective concurrency on eight workers, so another scheduling-barrier move is not evidence-based. Build77 preserves the exact Build76 task stream and commit order and decomposes existing generation-four work into single4, pair4 and continuation4 worker time/evaluations, plus call/output counts and a deterministic maximum-evaluation task. Timing is diagnostic only; Build76 remains qualified.

Build76 closes the ordered-generation-two experiment and becomes the current qualified smartphone baseline. Two independent retained Go 1.26.0 nine-photo runs passed exact 9/9 semantic equivalence. Full-matrix elapsed was 764,104 and 752,100 ms (758,102 ms mean) versus the Build75 two-run mean 836,722.5 ms (~1.104x / -9.4%). B/mild remained exactly 79,259 / 937 / 935 / 691 / 2,120,047 with physical decode 696, five speculative candidates and HMAC/payload PASS; B/angle remained 334,857 / 6,198 / 0 and REJECT. B/mild geometry averaged 60,351.5 ms versus Build75 76,149 ms; B/angle averaged 183,806.5 ms versus 209,554.5 ms. Build76 changes only public-only generation-two scheduling and preserves Build75 freeze, Build73 gen3, Build71 gen4, Build68 qualification and Build66 ordered decode semantics. A later interrupted invocation was found to have truncated the already-completed second-run TSV after the successful Markdown had been written; the finalized Build76 gate now stages diagnostics and publishes them only after a complete PASS.

## v0.3.0-build76 — ordered generation-two parallel candidate

After Build75 qualified ordered-parallel basin generation, the residual long-running geometry path is no longer dominated by freeze. On the second Build75 run, B/mild geometry is 76.556 s with freeze 23.251 s, prefix2 21.876 s, gen3 10.243 s and gen4 20.958 s; B/angle geometry is 209.184 s with freeze 20.812 s, prefix2 38.015 s, gen3 37.586 s and gen4 112.540 s. Gen3/gen4 already scale near the eight-worker limit, while prefix2 reaches only ~1.8x effective concurrency. Earlier Build72 stage profiling shows the dominant prefix seed is overwhelmingly second-generation work. Build76 therefore moves the ordered barrier to sibling1 and parallelizes only generation two, committing by frozen task index before unchanged generation-three/four pools. Build75 remains qualified until physical reproduction.

## v0.3.0-build75 — qualified ordered-parallel Build47 basin baseline

Build74 localized ~77–79% of the residual Build47 freeze to 16 independent basin tasks. Build75 freezes the exact historical production/depth/all-pairs stream, computes those tasks with one bounded pool, and commits results strictly by original task index. Two independent Go 1.26.0 physical runs passed the complete 9/9 semantic gate with exact Build73 logical telemetry. The two-run full-matrix mean is 836,722.5 ms versus the retained Build73 reference 1,047,729 ms (-20.1%); B/mild averages 202,688.5 ms (-18.2%) and B/angle 251,578 ms (-17.8%). Basin speedup versus Build74 averages ~4.66x on B/mild and ~4.58x on B/angle, while freeze speedup averages ~2.65x and ~2.72x respectively. Build75 is therefore the current qualified smartphone baseline.

## v0.3.0-build74 — Build47 freeze-stage observability snapshot

Build73 is now the qualified smartphone baseline. Its two retained physical runs leave a stable ~49–64 s Build47 freeze cost across the deep cases, while the larger residual generation-four cost on `B/angle` is already close to ideal eight-worker scaling. Build74 therefore makes no optimization. It runs the exact Build73 scheduler and replaces only the Build47 freeze call with a line-for-line equivalent profiled implementation that times structural seeding, plane construction, pair ranking, cells and basin generation, including the three Build47 diagnostic tiers. A direct regression requires the profiled freeze to reproduce the historical frozen bank/order and pair ranking exactly.

## v0.3.0-build73 — qualified ordered generation-three parallel baseline

The completed Build72 physical profile preserved Build71 semantics and localized the new geometry critical path to the third prefix generation. `B/angle` spends 194,940 / 293,120 ms cumulative prefix worker time in `pair3 + continuation3 + sibling3`, and its deterministic dominant prefix seed spends 149,514 / 203,169 ms there. `B/mild` shows the same pattern at 41,319 / 89,863 ms and 41,319 / 72,021 ms respectively. Build73 therefore leaves the search unchanged but moves the ordered parallel barrier to `sibling2`: all independent generation-three subtrees run in one bounded pool, are committed in original order, and feed the already-qualified Build71 generation-four pool. Two independent retained Go 1.26.0 physical runs subsequently passed exact 9/9 semantic equivalence and reproduced the performance gain. Build73 is promoted as the current qualified smartphone baseline; Build71 becomes the historical qualified ordered-generation-four milestone.

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


