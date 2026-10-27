# Build104 — exact fixed-block reader candidate

**Closed exact / public benchmark criterion not met / NOT PROMOTED, 2026-10-27. Test-only candidate. Current qualified Build98, previous Build93; Build100 unchanged.**

Build103 did not satisfy its predeclared material diagnostic-effect rule. It left both projective readers dominant: default profiles assign the exact Build84 continuation reader 55.80 / 437.14 sampled CPU-s on mild/angle. These shares identify a kernel to test, not guaranteed removable cost. Build104 keeps normal preemption and tests a compiler-shape hypothesis in the continuation4 reader only.

## Single hypothesis

Continuation4 always calls the reader with the Format-v4 block size, eight. The candidate makes that existing fixed-size contract explicit: remove the variable size argument, use constant eight loop bounds, and view the existing cosine-table rows through safe Go slice-to-array pointers `*[8]float64`. Their length checks happen at conversion, while fixed loop indices can be proven in bounds. No cosine values are copied or precomputed. RGB three-byte slice loads, homography/mapPoint, all floating-point operations and their evaluation order, accumulation order and rejection decisions remain unchanged. No unsafe, recurrence, LUT, image-sized allocation, DCT-product reassociation or new runtime flag.

This differs from Build86's rejected scalar row/column hoist: no cosine scalar load is explicitly hoisted and no new arithmetic ordering is introduced. It differs from Build88's rejected direct RGB loads: all four qualified RGB slice expressions remain identical. Only continuation4 is targeted; single4/pair4 retain their historical reader. Closed Build95/97 branches are not reopened.

The old Build97 public fixture reaches zero continuation FoldScores and is unsuitable for this changed kernel. Before any timing, a public seeded-noise coverage search selected the first root with nonempty downstream work: seed107, fixed quad `(7,8),(584,9),(8,503),(583,504)` on 592×512 pixels. The fixture uses its actual computed proposal score and requires 256 total evaluations, 128 continuation FoldScores, 16,384 continuation block reads and six ordered output candidates. Reader microbenchmarks retain the historical 1632×1632 public textures. This semantic-coverage correction precedes all benchmark data and changes no performance criterion.

## Bounded gates fixed before measurements

1. Verify exact float bits / success decisions across public deterministic front/mild/angle fixtures, random finite homographies, border and near-singular failures, different plane bounds and tiny images. Compare full FoldScore scores/visibility/fetch counters for all heldout/proposal/sparse modes, full continuation states/evaluations/telemetry, and full Generation4 ordered-bank fingerprints/telemetry. Race-test shared immutable inputs. A structural check must show the copied FoldScore/continuation/Generation4 glue differs only in reader dispatch and names. Keep qualified source files, pilot and frozen v3 unchanged.
2. Compile one Go1.26.0 binary. On the qualified i3-13100T host, fix GOMAXPROCS=8, explicit default preemption, GOGC=100 and GOMEMLIMIT=off. Use only public deterministic fixtures, no acquired photographs, keys or payloads. Measure reader front/mild/angle, complete single-task Generation4, and an ordered eight-task Generation4 pool. Fixture setup and correctness checks stay outside timed loops; report ns/op, bytes/op and allocations/op. The pool includes the original bounded-pool/order overhead in both conditions.
3. Run **two complete balanced benchmark sessions**, all metrics in each. For each metric use Q C C Q | Q Q | C Q Q C: four adjacent qualified/candidate pairs plus a separate identical-qualified Q/Q control between the two four-slot blocks. Fresh child per slot; same binary and inputs; one-second minimum measured duration per slot, one-second cooldown before launch. Preserve all raw output, iterations, pairs, control contrasts and session means. External child timeouts and stable source/binary hashes are required. CPU/compiler inspection is separate from timing. Benchmark order/time/metric choices are fixed here; no tuning after seeing results.
4. **Authorize a later retained comparison only if complete Generation4 and the eight-task pool improve at least 3% on the mean in each of both sessions, every paired comparison for those two stage metrics is favorable, all three reader case means are non-regressive in both sessions, and allocations satisfy the fixed gate: exact mean bytes/counts for reader/single-task metrics, identical mean pool allocation count and pool mean bytes within 1% in each session (goroutine bookkeeping can vary).** Every exactness/source/race gate must pass. Q/Q contrasts describe measurement variation and do not relax the criterion. Otherwise close Build104 exact / benchmark-negative / not promoted, without a private physical campaign. This is a descriptive gate, not a significance test or proof of private-image performance.

Build104 is the public benchmark candidate checkpoint. A positive gate authorizes a separately documented retained comparison; it does not promote a runtime or change production dispatch. Any future promotion still needs exact retained comparison, two independent physical semantic matrices and normal-command integration. No protected decoding or HMAC qualification is claimed here.

## Reproduction

```sh
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build104-fixed-block-test
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build104-fixed-block-benchmark
```

Evidence is retained under `v4-phone private/build104-fixed-block-benchmark`; existing outputs are refused and failed/partial results preserved. Only reviewed numeric tables and hashes will be committed publicly. Candidate functions are test-only and have no CLI wiring.

## Results — exactness PASS / complete public benchmark / not promoted

Qualified Go1.26.0 / i3-13100T / Linux amd64 / eight workers. Reader bit/border/random/near-singular/tiny-plane cases, all FoldScore modes, full continuation states and complete nonempty Gen4 bank/counters match exactly. The shared eight-task pool and race checks pass. Inherited Build98 exactness/race/CLI/buildinfo/vet, frozen v3 and pilot checks pass. The structural checker permits only the declared fixed8 reader contract and named dispatch. The public fixture freezes **256 evaluations / 128 continuation FoldScores / 16,384 continuation block reads / six ordered candidates** per task. An eight-task pool operation is eight copies of that same synthetic task on one shared immutable plane. Its test-only bounded queue restores original task-index order in both conditions; it does not reproduce the heterogeneous retained task distribution or every duration-instrumentation call of the full qualified pool. It is a concurrency/throughput benchmark, not full-bank recovery or a private-image prediction.

Both complete balanced sessions produced all **100/100** slots: 80 comparison samples (40 adjacent Q/C pairs) plus 20 identical-qualified control samples (ten Q/Q contrasts). Fresh children used the same binary and public inputs, with explicit default preemption. Source/binary/protocol hashes pass before and after measurements. No private photographs, protected decode, payload/HMAC or physical qualification were used.

### Reader and complete-stage means

Only comparison slots enter means; controls remain separate. Negative change is faster candidate execution. Reader operations are one 8×8 block; Generation4 is a complete task; Pool8 includes eight tasks, bounded queue and original-index collection.

| session | reader | qualified ns/block | candidate ns/block | change |
|---|---|---:|---:|---:|
| 1 | Front | 1,570.250 | 1,431.000 | -8.87% |
| 1 | Mild | 1,530.750 | 1,443.000 | -5.73% |
| 1 | Angle | 1,564.000 | 1,461.750 | -6.54% |
| 2 | Front | 1,538.250 | 1,443.000 | -6.19% |
| 2 | Mild | 1,548.250 | 1,438.250 | -7.10% |
| 2 | Angle | 1,557.250 | 1,500.000 | -3.68% |

| session | work | qualified ms/op | candidate ms/op | change | favorable pairs |
|---|---|---:|---:|---:|---:|
| 1 | Generation4 | 52.285 | 51.198 | -2.08% | 4/4 |
| 1 | Pool8 | 101.423 | 98.939 | -2.45% | 4/4 |
| 2 | Generation4 | 52.573 | 51.558 | -1.93% | 3/4 |
| 2 | Pool8 | 102.485 | 98.308 | -4.08% | 4/4 |

Reader means are favorable in both sessions. Complete Generation4 remains below the fixed 3% requirement in both. Pool8 meets 3% only in session2. Session2 also contains an adverse complete-Generation4 pair. Thus the speed conjunction fails independently of allocation interpretation; no pooled mean or isolated kernel speed overrides the gate.

### Controls, allocations and limits

All ten Q/Q changes remain in `build104-controls.tsv`. Their apparent changes range from -3.85% to +0.78% for reader slots and -1.35% to +1.00% for complete-stage/pool slots. They are descriptive timing variation, not confidence bounds or a revised acceptance threshold. These are two balanced batches on one host, not independent-machine replications or proof of private-image speed.

All readers remain 0 bytes/op and 0 allocations/op. Generation4 reports nine allocations/op in both conditions, but session1 mean bytes are 121,706.5 / 121,822.0, failing the prospectively exact equality rule; session2 bytes match at 121,704. Pool allocation counts also vary (session1 Q/C 93.75/91.75, session2 91.25/91.50), failing the unchanged-count rule even where counts decrease. These process-wide benchmark observations are retained; they do not alone identify an allocation mechanism. Pool byte means stay within the fixed 1% allowance. Source inspection and the zero-allocation reader show no new image-sized cache. The fixed allocation guard was not relaxed after seeing these rows.

### Compiler inspection after timing

The default-compiled measured binary's assembly contains four inner DCT bounds-panic call sites in the qualified reader and none in the candidate. Table-row entry sites rise from two to four; RGB slice sites remain eight, and the mapPoint call remains one. Total static `runtime.panicBounds` sites change 14→12. These are static sites classified by source line, not dynamic executions or isolated causal CPU measurements. Safe row conversions replace repeated inner index checks; all arithmetic statements remain structurally identical. `build104-compiler.tsv` records this separate inspection.

An external sandbox/host I/O failure occurred **after** the complete benchmark and source/binary checks. Recovery verified every original child output and all100 rows. The optional full SSA-report file was recovered empty; it is excluded. Preserved assembly from the measured binary supplies the compiler evidence above. No benchmark was interrupted, repeated, dropped or reclassified because of this later failure.

## Decision and next bounded proposal

**Close Build104 exact / benchmark criterion not met / NOT PROMOTED.** Do not run a retained or nine-photo campaign under this failed gate. Qualified normal extraction stays Build98, previous Build93; Build100's historical decision stays closed. The test-only candidate remains reproducible without production/CLI wiring.

The reader's positive microbenchmark result may motivate a separate hypothesis on the other profile-dominant historical projective reader: test that reader's own fixed8/table contract, keeping its sampling expressions unchanged. That would be a new public benchmark candidate, not a reopening of Build95's luminance plane or Build97's reader-dispatch extension. No Build105 implementation or performance is claimed here. Any future candidate needs its own predeclared controls and full gates.

## Frozen evidence

Six public tables under `docs/qualified-baselines` preserve all100 samples, ten means, 40 pairs, ten controls, both session decisions and the separate two-row compiler summary. `build104-evidence.sha256` freezes them plus private `build104-fixed-block-benchmark.zip`. Raw child logs, source/binary/protocol manifests, the measured test binary, original assembly, fixture-coverage provenance and source/race logs remain in the archive. The public summarizer regenerates five numeric tables byte-for-byte; compiler rows are classified from preserved assembly. These are candidate decision references, not qualified-baseline evidence.
