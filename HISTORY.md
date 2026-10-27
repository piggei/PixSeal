## Build102 — frozen Gen4 A/A replay observability PASS, 2026-10-27

Build102 introduces a test-only split at the existing Gen3→Gen4 barrier. Qualified Build98 full-bank references and both Build98/Build100 prefix captures agree exactly; all 20 isolated replays preserve complete ordered-bank, task/raster fingerprints and non-timing counters. Source/exactness/race/CLI/vet, structural/mutation/CPU-parser and frozen-core/pilot gates passed with stable source/corpus hashes. The banks retain 937 / 6,198 candidates from 593 / 1,902 Gen4 tasks.

Sixteen identical-code/data timing slots show pool range/mean 3.414699% / 7.136622% and apparent B/A mean changes -1.157366% / +0.533301% on B/mild / B/angle. All eight paired values and all four separate CPU-profile rows are retained. Live total process CPU replaces the stale CPU-class snapshots for this study. These descriptive ranges are not historical confidence bounds or causal attribution. Pprof assigns 5.39–7.66% samples to asyncPreempt; a separate research-only preemption-attribution control is proposed, without equating samples to removable cost or recommending a production debug setting. Build98 stays qualified, Build93 previous and Build100 closed / not promoted. See `docs/V4_BUILD102_FROZEN_GEN4_REPLAY.md` and the four `docs/qualified-baselines/build102-*.tsv` tables.

## Build101 — runtime-effects observability PASS, 2026-10-27

Build101 preserves the pinned Build98/Build100 runtime sources and introduces only test-side observation and version metadata. Source/exactness/race/CLI/vet, frozen-core/pilot and structural-copy gates passed. All 12 retained executions preserve full exact telemetry, payload and successful info; all source/corpus hashes remained fixed. Eight untouched AB/BA timing runs give Build100 mean deep changes of -11.92% (B/mild) / -5.41% (B/angle), freeze -39.03% / -37.68% and Gen4 -1.35% / -2.21%, favorable in each timing pair.

The separate four-run attribution pass retains the opposite Gen4 sign (+1.99% / +3.36% wall; +1.74% / +2.91% sampled CPU). Both variants spend CPU in unchanged projective readers, show nearly identical allocation counts and high worker occupancy; no unique GC/scheduler cause is established. CPU-class runtime counters have a snapshot cadence limitation and are not treated as live phase CPU. Build101 closes as observability PASS; it neither reopens Build100's historical guard nor promotes it. Build98 stays qualified and Build93 previous. The proposed next study isolates replay of a frozen Gen4 task vector from upstream context. See `docs/V4_BUILD101_RUNTIME_EFFECTS_OBSERVABILITY.md` and the five `docs/qualified-baselines/build101-*.tsv` tables.

## Build100 — ordered-parallel pair scoring, closed without promotion

Build99 identified serial six-pair public proposal scoring within freeze. Build100 changes only scheduling of these fixed tasks, retaining original-index commit and unchanged stable ranking/downstream recovery. Source/exactness/race and the public complete-freeze benchmark passed (-48.43% mean). Two independent retained sessions preserved exact semantics in all 16 executions; deep means improved by 11.45% / 15.20% on B/mild and 1.31% / 4.15% on B/angle. Nonetheless the predeclared unchanged-stage guard failed because B/angle Gen4 remained slower on both session means (+3.82% / +1.28%). Build100 is closed without physical gates or promotion; qualified normal production remains Build98. See `docs/V4_BUILD100_ORDERED_PARALLEL_PAIR_SCORING.md`.

## Build99 — post-promotion observability of qualified Build98

Build98's archived production integration is now documented as closed, with exact nine-row semantic equivalence to both promotion runs. Build99 measures unchanged Build98 using low-overhead B/mild/B/angle timing followed by separate CPU profiles. Historical Build94 preparation/accounting helpers are reused, while Build98 Generation3 telemetry is additionally frozen. The qualified algorithm and normal production dispatch remain Build98; Build93 remains previous. Source and complete retained timing/profile gates passed on 2026-10-23 with exact telemetry. B/mild deep time is 83,618 ms and B/angle 187,601 ms; Gen4 remains 60.56% of B/angle deep wall. Both CPU profiles expose serial six-pair proposal scoring (12.29 / 10.27 CPU-s), selecting ordered-parallel scoring as the next bounded candidate. At the Build99 checkpoint, Build100 was not yet implemented. See `docs/V4_BUILD99_QUALIFIED_BUILD98_PIPELINE_PROFILE.md` for results and reproduction commands.

## Build98 — Exact Generation3 Fetch Extension (qualified)

Build97 completed two corrected retained Build93/Build97 AB/BA runs with exact semantics, but failed its predeclared performance criterion. B/mild Gen4 was favorable across all four pairs only because of the first run; the independent replication was 0.76% slower on its two-row mean. B/angle Gen4 was 0.98% slower on pair-normalized percentages and 1.27% slower on pooled means, with direction changing by order. End-to-end totals were effectively neutral. Build97 is therefore closed without a physical 9/9 gate and Build93 remains qualified.

Build98 targets a different stage selected from the completed Build94 CPU profiles. Generation3 consumes 76.56 CPU-s on B/mild and 285.63 CPU-s on B/angle, while its historical `readProjectiveBlockValue` path accounts for 76.19 / 283.69 CPU-s. Build98 reuses the qualified Build84 exact RGB/FoldScore reader for Generation3 single, pair, continuation and sibling scoring only, with no image-sized memory and no change to Generation4 or downstream semantics.

The qualified Go 1.26.0/i3-13100T public benchmark passed cleanly: the complete Generation3 mean fell from 23.513 to 21.304 ms/op (-9.40%), sibling stencil improved 9.54%, allocations were unchanged, and every Build98 complete-stage sample beat every historical sample. The subsequent retained Build93/Build98 AB/BA semantic gate also passed. Generation3 was faster in all four paired comparisons: 4.45% and 23.48% on B/mild, 5.19% and 16.47% on B/angle. Unmodified Gen4 and downstream phases still exhibited host drift, so promotion is judged on the isolated changed phase plus complete semantics. Two independent nine-photo physical gates then passed 9/9 with exact frozen semantics. Build98 is promoted as the current qualified smartphone baseline, Build93 becomes previous, and normal `v4-extract-phone` is rewired to the unchanged qualified Build98 implementation. The separate normal-command production-entry-point matrix also completed 9/9 PASS. Its archived report, metadata, checksum and exact semantic equivalence with both promotion matrices were verified on 2026-10-23; `docs/qualified-baselines/build98-production-integration.tsv` closes wiring evidence without being counted as a third performance run.

## Build97 — Exact Gen4 Probe Fetch Extension (benchmark-first candidate)

Build96 retained profiling closed the Build95 luminance-plane branch and showed a second large Gen4 CPU share still flowing through historical generic luminance sampling in single4/pair4. Build97 therefore reuses the already-qualified Build84 exact RGB projective reader for those two Gen4 probe stages only, with zero image-sized memory overhead. Build93 remains qualified and Build97 cannot advance beyond public benchmark without material same-session evidence.

## v0.3.0-build96 — Gen4 retained bottleneck profiling

Build95 is closed as exact/race PASS but not promoted: retained AB/BA evidence did not justify the added 127,475,200-byte luminance plane. Build96 is observability-only and profiles B/angle Build93/Build95 AB/BA with per-run CPU profiles plus memory/GC accounting. Build93 remains the current qualified smartphone baseline.

## v0.3.0-build95 — exact generation-four luminance-plane candidate

Build94 completed the required post-promotion profile of unchanged qualified Build93. Its low-overhead pass measured B/mild at 80,152 ms deep total (57,175 ms geometry, 17,463 ms gen4, 12,511 ms qualification, 10,251 ms decode) and B/angle at 178,627 ms (176,938 ms geometry, 106,635 ms gen4, 1,453 ms qualification, no protected decode). All frozen semantic counters remained exact.

The separate B/angle profile sampled 1,384.63 CPU-s. Generation 4 accounted for 920.46 CPU-s cumulative (66.48%), continuation4 for 459.29 CPU-s and the exact Build84 continuation4 reader for 456.30 CPU-s (32.95% of the complete profile). Major pools already used all eight logical CPUs, so another outer scheduling change was not justified. Qualification and protected decode were also not B/angle targets.

Build95 therefore tests one bounded structural amortization. Immediately before generation 4, it computes the exact historical RGB-to-luminance expression once for every pixel in the immutable work plane. Only continuation4 FoldScore projective reads use this `float64` plane; mapping, bounds, clamping, bilinear interpolation and DCT accumulation retain their exact order. Prefix1, gen2, gen3, generation-four single/pair probes, bank order, Build90 qualification and Build93 protected decode are untouched.

The additional plane is capped at 128 MiB using public dimensions only. The retained 3464x4600 work image needs 127,475,200 bytes (about 121.57 MiB); over-cap or invalid planes fall back to the unchanged Build84 reader. No secret or correctness signal can affect this choice. Build95 adds bit-exact kernel/bank tests, a targeted shared-plane race fixture, a public benchmark and a contemporaneous Build93/Build95 AB/BA retained profiler. It remains an unqualified candidate; normal `v4-extract-phone` remains Build93.

## v0.3.0-build94 — completed qualified Build93 post-promotion pipeline observability

Build93 is qualified and its production-entry-point integration is closed. Because protected-residue reuse materially reduced B/mild decode, the Build92 priority map taken on Build90 was no longer authoritative. Build94 deliberately changed no recovery algorithm and re-measured the current qualified Build93 deep path before selecting a later optimization.

The Go 1.26.0 source and profile gates passed. B/mild measured 80,152 ms deep total: 57,175 ms geometry (71.33%), 21,768 ms freeze, 9,906 ms gen3, 17,463 ms gen4, 12,511 ms qualification and 10,251 ms decode. B/angle measured 178,627 ms: 176,938 ms geometry (99.05%), 20,536 ms freeze, 36,072 ms gen3, 106,635 ms gen4, 1,453 ms qualification and zero decode.

The separate CPU profiles localized the new target without contaminating those wall-time rows. On B/angle, generation 4 was 920.46 of 1,384.63 CPU-s and the Build84 continuation4 reader was 456.30 CPU-s. On B/mild, the same reader was only 51.60 of 513.24 CPU-s, while Build93 protected sampling remained 69.82 CPU-s. The evidence selected Build95 exact generation-four luminance-plane reuse as a benchmark-first candidate. Build94 remains observability-only and Build93 remains the current qualified baseline.

## v0.3.0-build93 — qualified exact protected-residue reuse baseline

Build92 profiled the promoted Build90 pipeline and found B/mild protected decode dominated by repeated projective sampling: 33,460 ms decode wall time, 230.14 CPU-s cumulative in `experimentalV4ReadProtectedProjectiveMargins`, and only 5.81 CPU-s in Build42 list decode. Build93 therefore makes one narrow post-freeze change: each already-qualified candidate samples the protected projective grid once into the exact historical residue sums/counts, then derives Robust, Balanced and Capacity margins from those immutable residues in unchanged logical order. Geometry, Build90 ordered-parallel qualification, candidate ordering, list/HMAC ordering and wire/security semantics remain unchanged.

The qualified Go 1.26.0 public benchmark passed with exactness/race clean: single Capacity candidate 44.439 -> 10.236 ms/op (**4.341x, -76.97%**) and batch8 49.443 -> 22.324 ms/op (**2.215x, -54.85%**), with effectively unchanged allocations.

Two independent retained physical matrices then passed the complete 9/9 semantic gate. B/mild remained exactly `79,259 / 937 / 935 / 691 / 2,120,047`, physical decode `696`, speculative `5`, HMAC/payload PASS, with `696` protected sampling passes and `2,086` profile derivations in both runs. Run1/run2 B/mild deep totals were 97,913 / 75,704 ms and decode 12,878 / 9,879 ms. Their means are 86,808.5 ms deep (**-26.50%, 1.361x** versus the Build90 paired reference 118,104 ms) and 11,378.5 ms decode (**-65.99%, 2.941x** versus the Build92 Build90 measurement 33,460 ms). B/angle remained `334,857 / 6,198 / 0`, REJECT, with zero protected sampling; its large cross-run geometry variation is retained as host-drift evidence and is not attributed to Build93.

**Decision: PROMOTED. Build93 is the current qualified smartphone baseline; Build90 becomes the previous qualified baseline.** Normal `v4-extract-phone` now uses Build93. The explicit `v4-extract-phone-build93` entry point is retained as a pinned regression/requalification alias.

The post-promotion consolidation removes the duplicated normal-command decoder body: `ExperimentalV4ExtractPhone` now delegates directly to the unchanged, physically qualified `ExperimentalV4ExtractPhoneBuild93Candidate` implementation. The retained physical target runs the normal production command, writes a separate integration artifact and records the selected entry point. On 2026-10-16 the qualified Go 1.26.0 source gate passed and the normal `v4-extract-phone` command completed the full 9/9 matrix with exact semantic telemetry. Metadata confirmed `gate_purpose=production-entrypoint-integration`, `entrypoint=v4-extract-phone`, Build93/Build90 qualified/previous baselines, nine cases and the executable SHA-256. This closes only rewiring/integration risk and does not reinterpret the two-run algorithm qualification.

## v0.3.0-build92 — post-promotion Build90 pipeline observability

Build91 closed the promotion decision and made Build90 the current qualified smartphone baseline. Because ordered-parallel qualification removed the former B/mild bottleneck, Build92 deliberately makes no recovery change and re-measures the promoted pipeline before any new optimization is selected.

The profile separates timing from instrumentation: one retained B/mild/B/angle pass runs unchanged Build90 with only its existing telemetry, while separate later passes collect CPU `pprof` around deep recovery and generate flat/cumulative, phase-focused and source-list reports. Build92 reuses the frozen Build91 semantic expectations and Build90 race regression. Its result can select a later Build93 hypothesis, but Build92 itself is observability-only and cannot supersede Build90.

## v0.3.0-build91 — paired physical decision checkpoint closes Build90 promotion

Build90 produced the first large post-Build84 pipeline-level speed signal: public qualification benchmarks on the qualified i3-13100T/Go 1.26.0 host improved by 4.283x on a high-pass workload and 3.981x on an early-reject workload, with the targeted race gate clean. Three retained nine-photo runs were semantically exact. Cross-day absolute timing still drifted materially, including one broad host-degraded run, so Build91 was created to measure Build84 and Build90 contemporaneously rather than cherry-pick historical baselines.

Build91 changes no recovery algorithm. Its AB/BA paired profile decodes each retained deep image once and runs Build84/Build90 on the same in-memory work image, reversing variant and case order in cycle 2. The paired semantic gate passed. B/mild paired deep total improved from 168,444.5 ms to 118,104 ms (-29.89%) and qualification from 64,185.5 to 14,975 ms (4.286x). B/angle improved from 220,870 to 212,025.5 ms (-4.00%) and qualification from 7,602.5 to 1,562.5 ms (4.866x). Geometry/decode stayed comparable within pairs. **Build90 is therefore promoted as the current qualified smartphone baseline; Build84 becomes the previous qualified baseline.** Build91 remains observability-only evidence. The harness-only absolute-path correction is retained for reproducible package-test execution.

## v0.3.0-build90 — qualified ordered-parallel qualification baseline selected from full-pipeline evidence

Build89 completed the full qualified Build84 pipeline profile with 9/9 retained semantic equivalence. It changed the optimization priority: on `B/mild`, serial public-pilot qualification consumed 61,193 ms (40.93% of deep recovery), exceeding both geometry and ordered protected-data decode. The frozen bank has 937 candidates and 935 qualify, making the stage highly parallel without altering the bank or protected-data semantics.

Build90 therefore parallelizes only qualification execution. It preserves the Build84 freeze barrier, writes one result per original bank index, waits for all tasks, then commits evaluation counts and qualified hypotheses serially in bank order. Public benchmark/race, three full physical semantic gates, and the Build91 contemporaneous paired profile subsequently qualified this unchanged implementation. Build90 became the current qualified baseline at that checkpoint and was later superseded by Build93; Build84 remains an earlier qualified baseline.

## v0.3.0-build89 — qualified Build84 full-pipeline profiling

Build88 passed exactness but regressed the qualified Build84 reader by roughly 8–11% and therefore never reached the private physical gate. Build89 restores Build84 as the active deep runtime and moves observability back to the complete retained workload. It derives stage accounting from the qualified timers already present in Build84 and profiles freeze, prefix1, gen2, gen3, gen4, qualification and ordered decode before any Build90 optimization is chosen.

## v0.3.0-build88 — direct RGB bounds-check experiment, closed negative

Build87 completed the qualified RGB/luminance/bilinear compiler profile without changing runtime semantics. The key actionable compiler result is eight `IsSliceInBounds` checks on Build84's four `rgb[index:index+3]` expressions. Build88 therefore tests one final narrow exact-portable reader rewrite: one dominating `rgb[index+2]` proof per pixel followed by direct scalar R/G/B loads. All floating-point expressions and DCT order remain Build84-exact. The candidate is benchmark-first and Build84 remained the qualified smartphone baseline at that checkpoint; Build90 later superseded it, and Build93 is now current.

## v0.3.0-build87 — qualified RGB/luminance/bilinear profiling

Build86 closed as an exact but benchmark-negative DCT-table-hoist experiment and never reached the private physical corpus. Its qualified-host reader means were slower than Build84 on front/mild/angle, and the isolated hoisted DCT kernel itself regressed by about 8.5%. Build87 therefore restores Build84 as the active qualified runtime and performs no recovery change. It decomposes the remaining RGB fetch -> luminance -> bilinear hot path with public deterministic benchmarks, integrated pprof, BCE/inlining output and generated assembly so any later Build88 candidate is selected from fresh evidence rather than another speculative source rewrite.

## Build86 — exact DCT table hoist / bounds-check experiment — closed negative

Build85 selected one exact DCT cosine-table load-placement experiment. Build86 passed block/FoldScore/continuation/bank exactness on the qualified Go 1.26.0 host, but its same-session five-run reader means regressed versus Build84: front 1382.4 vs 1348.4 ns/block, mild 1513.0 vs 1380.8, angle 1435.6 vs 1425.4. The isolated DCT kernel also regressed from 65.588 to 71.166 ns/block (+8.5%). Build86 therefore stopped before the private corpus and is not promoted. Compiler evidence showed that moving/removing inner DCT bounds checks was not sufficient to improve generated-code performance.

## Build85 — qualified Build84 sampler post-promotion profiling

After Build84 became the current qualified smartphone baseline, Build85 deliberately makes no decoder optimization. It reuses the public deterministic sampler fixture and profiles the **qualified Build84 reader** with repeated whole-reader/stage benchmarks, CPU pprof, line-level reader attribution and compiler BCE/inlining evidence. The goal is to select the next single optimization from the new baseline rather than continuing from Build83 measurements taken before the Build84 RGB-fetch improvement. Build84 remains qualified; Build85 requires no physical corpus gate unless a later candidate changes runtime behavior. The qualified-host profile completed successfully: front/mild/angle whole-reader means were 1384.6 / 1462.4 / 1441.2 ns/block; DCT accumulation accounted for ~18.6% of integrated CPU samples and retained four cosine-table bounds checks. This evidence selected Build86; Build85 remains observability-only and requires no physical gate.

## Build84 — qualified exact RGB fetch / bounds-check baseline

Build83 completed successfully as an observability-only checkpoint on the qualified Go 1.26.0 i3-13100T host. CPU pprof attributed 72.33% cumulative time to `samplePlaneLuminance`, 5.65% to `homography.mapPoint` and 21.15% flat to surrounding `readProjectiveBlockValue` work. This evidence rejected mapPoint/DCT as the immediate target and selected RGB fetch/addressing for Build84.

Build84 keeps the exact Build82 continuation4 sampler arithmetic but changes only RGB address formation and bounds-check shape: row bases and x byte offsets are computed once, then each of the four bilinear source pixels is represented by a length-three slice before evaluating the unchanged `.299*R + .587*G + .114*B - 128` expression. No `unsafe`, cache, LUT, precomputed luminance plane, geometry change or floating-point reordering is introduced.

The qualified Go 1.26.0 public benchmark measured Build84 at 1413.2 / 1408.0 / 1410.0 ns/block mean for front/mild/angle-like fixtures versus 1500.4 / 1514.0 / 1522.2 ns/block for the historical reader: 5.81%, 7.00% and 7.37% reader-level improvements. Compiler BCE evidence still reports `IsSliceInBounds`, so the result is documented as a favorable address/bounds-check shape, not as complete bounds-check elimination.

Two independent retained physical runs then passed the complete 9/9 semantic gate with exact workload. Full-matrix elapsed was 748,902 / 746,902 ms (747,902 ms mean) versus the 758,102 ms Build76 mean. B/mild elapsed/geometry/gen4 means were 185,370.5 / 59,173.5 / 18,025 ms; B/angle were 223,920.5 / 182,137 / 109,696 ms. Both runs preserved B/mild `79259 / 937 / 935 / 691 / 2120047` with HMAC/payload PASS and B/angle `334857 / 6198 / 0` with REJECT.

**Decision: Build84 is promoted as the current qualified smartphone baseline.** Build76 remains the previous qualified scheduler/runtime reference. Source-controlled timing reference files under `docs/qualified-baselines/` freeze the Build76 historical mean and the Build84 two-run qualified baseline without shipping private images or diagnostic logs.

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

Build71 qualification leaves the prefix through `sibling3` as the measured geometry critical path, especially on `B/angle`. Build72 makes no scheduling or algorithm change. It times the exact fourteen public-only prefix stages, aggregates stage evaluations/worker time across all seeds, and records the deterministic maximum-evaluation prefix seed plus its per-stage genealogy. Build71 remained the qualified smartphone baseline at that checkpoint; later baselines superseded it and Build93 is current.

## v0.3.0-build71 — qualified smartphone baseline

Two independent retained Go 1.26.0 physical runs reproduce the complete Build68 semantic matrix while confirming the Build71 generation-four scheduling gain. `B/mild` remains exactly 79,259 geometry evaluations / 937 bank / 935 qualified / 691 logical decode candidates / 2,120,047 logical frames and authenticates the expected payload; `B/angle` remains 334,857 / 6,198 / 0 and rejects. Full-matrix elapsed is 1,233,478 and 1,263,623 ms versus the retained Build68 reference of 1,902,847 ms. Build71 is therefore promoted as the current qualified smartphone baseline. Build68 remains the historical qualification-plane-reuse milestone, Build66 the historical ordered-parallel decode milestone, and Build64 the deep-recovery semantic reference.

## v0.3.0-build71 — ordered generation-four parallel performance candidate

Build70 physical profiling passed the complete retained nine-photo semantic-equivalence gate and identified a deterministic public-only scheduling opportunity. The dominant seed contributes 73922/79131 seed evaluations on `B/mild` and 267775/334749 on `B/angle`. Its fourth generation (`single4 + pair4 + continuation4`) contributes 41776/73922 evaluations on `B/mild` and 196641/267775 on `B/angle`, with 593 and 1363 independent `sibling3` inputs respectively. Build71 preserves the exact Build64/68 search and final bank order but splits geometry scheduling into two deterministic barriers: parallel seed prefixes through `sibling3`, then one bounded global pool for the frozen fourth-generation subtrees, followed by original-order commit. Build68 remains the qualified smartphone baseline until Build71 passes the physical equivalence/performance gate.

## v0.3.0-build70 — dominant-seed genealogy profiling snapshot

The qualified Go 1.26.0 Build69 physical gate reproduced Build68 semantics exactly and established that the remaining geometry runtime is dominated by one seed branch, not by broadly distributed work. `B/mild` measured 168398 ms for the slowest seed inside a 169014 ms seed wall region; `B/angle` measured 693395/695422 ms; `control/mild` measured 358584/360285 ms. Build70 deliberately makes no geometry optimization. It reproduces the same search while recording freeze-vs-seed evaluation counts, the dominant seed identity, and per-stage evaluation/state genealogy. Build68 remained the qualified smartphone baseline at that checkpoint; later baselines superseded it and Build93 is current.

## v0.3.0-build69 — geometry-generation profiling result

The retained nine-photo Go 1.26.0 gate passed semantic equivalence. Deep cases showed extreme load imbalance: max-seed/seed-wall time was 96.2% (`control/front`), 99.5% (`control/mild`), 97.3% (`control/angle`), 99.6% (`B/mild`) and 99.7% (`B/angle`). The slowest seed also produced the complete final bank for `control/front`, `control/mild` and `B/mild`, and 5486/6198 states for `B/angle`. This evidence redirects Build70 from an optimization attempt to finer-grained genealogy profiling.



### Build97 benchmark-review extension

Qualified-host Go 1.26.0 benchmarking showed a 6.66% mean reduction for the complete Generation4 benchmark (28.314 -> 26.427 ms/op), with unchanged full-Generation4 allocation counts. This is sufficient to proceed to retained Build93/Build97 AB/BA evidence, not to promotion. Added `v4-build97-gen4-probe-fetch-paired-profile` and its exact semantic gate.

### Build97 retained telemetry remediation

The first retained B/mild pair authenticated correctly and preserved the complete geometry/qualification/decode workload, but the candidate reported `41,776` FoldScores and `53,473,280` reads instead of the frozen continuation4 values `16,496` and `21,114,880`. The exact delta (`25,280` and `32,358,400 = 25,280 * 1,280`) came from single4/pair4 calls newly routed through the Build84 instrumented reader. This is not a search/output mismatch, but changing a historical counter family is still disallowed. Build97 is remediated in place: Build84 `Fetch*` remains continuation4-only; new `ProbeFetch*` counters account for single4/pair4. The paired harness continues to require exact Build93 historical telemetry and must be rerun from the beginning.
