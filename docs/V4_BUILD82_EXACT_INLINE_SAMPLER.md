# Format-v4 Build82 — specialized exact inline block sampler

Build82 is a **closed equivalence-preserving performance experiment** over the qualified Build76 smartphone baseline. Build76 remains authoritative. Two independent Go 1.26.0 physical runs passed the complete 9/9 semantic matrix and preserved the deep workload exactly, but the performance benefit was not repeatable; Build82 is therefore **semantic PASS x2 / non-promoted**.

## Why Build82 exists

Build77-79 localized the remaining continuation4 cost progressively:

- Build77: the generation-four worker pool is already close to ideal saturation;
- Build78: continuation4 worker time is ~99.9% FoldScore;
- Build79: projective block reads account for ~97% of sampled FoldScore time, with bilinear luminance sampling the largest replay component.

Build80 proved that a local spatial luminance cache adds too much preparation/fallback overhead. Build81 proved that replacing the three RGB multiplications with exact lookup tables is semantically safe but does not provide a repeatable speedup.

Before writing Build82, the real qualified toolchain was checked directly:

```bash
GOTOOLCHAIN=go1.26.0 go test ./watermark -run '^$' -gcflags='-m=2' 2>&1 | \
  grep -E 'samplePlaneLuminance|readProjectiveBlockValue|homography.mapPoint'
```

Relevant Go 1.26.0 results:

```text
samplePlaneLuminance        cannot inline: cost 348 > budget 80
homography.mapPoint         cannot inline: cost 88  > budget 80
readProjectiveBlockValue    cannot inline: cost 263 > budget 80
```

The small closure inside `samplePlaneLuminance` is inlined, but the outer helper is not. This makes helper-call/structure overhead a justified next hypothesis.

## Exact optimization under test

Build82 changes only the projective block reader used by FoldScore inside generation-four `continuation4`.

The specialized reader:

1. keeps the historical `homography.mapPoint` call unchanged;
2. manually incorporates the exact body of `samplePlaneLuminance` into the block loop;
3. hoists immutable `width`, `height`, `rgb` and DCT row references outside the per-sample luminance path;
4. preserves the historical RGB→luminance formula and operation order:

```go
.299*float64(r) + .587*float64(g) + .114*float64(b) - 128
```

5. preserves `math.Floor`, x1/y1 clamps, interpolation fractions, top/bottom bilinear interpolation order, `c23`/`c32` accumulation order and the final `abs(c23)-abs(c32)` result.

Build82 deliberately does **not** manually inline `homography.mapPoint`. That helper is also above the compiler inline budget, but changing it in the same build would make the performance result harder to attribute and would increase floating-point equivalence risk.

## Explicit non-changes

Build82 introduces:

- no LUT;
- no source bounding rectangle;
- no local cache or scratch buffer;
- no fallback path;
- no projective coordinate recurrence;
- no pruning, new threshold, new score or early exit;
- no encoder, wire-format, pilot, strength, ECC/Hamming, whitening/HMAC or protected-data change.

`single4`, `pair4`, Build76 prefix1/gen2, Build73 gen3, Build75 freeze, Build68 qualification and Build66 ordered protected-data decode stay on the previously qualified behavior.

## Exactness gates

Before any physical timing is considered, Build82 must pass:

1. `math.Float64bits` equality of the specialized projective block reader against historical `readProjectiveBlockValue` across multiple homographies, regular blocks, boundary blocks and failure cases;
2. Build41 FoldScore bit equality;
3. Build55 continuation4 state-for-state and eval-for-eval equality;
4. complete Build82-vs-Build76 blind-bank equality: evaluation count, seed/worker view, bank length/order, quad, homography, proposal and validation;
5. existing Build47/68/71/73/75/76 regressions, CLI help/buildinfo and `go vet`.

The candidate exposes only low-overhead integer counters for FoldScore calls and block reads/success/failure. No fine-grained timer is inserted into the hot sampler. The physical gate checks those counters against the retained exact Build81 workload (which is itself Build76-equivalent): controls `1104/1413120/1413120/0`, `20368/26071040/26071040/0`, `304/389120/389120/0`; B/mild `16496/21114880/21114880/0`; B/angle `118004/151045120/149370664/1674456` (`FoldScore / reads / success / failed`).

## Physical gate

Run on the qualified Go 1.26.0 host:

```bash
make clean
make v4-build82-phone-inline-sampler-test
make v4-build82-phone-physical-test
```

Expected output directory:

```text
v4-phone private/build82-diagnostics/
```

with:

```text
build82-phone-inline-sampler.tsv
build82-phone-inline-sampler.md
logs/
```

The harness stages output atomically. A failing semantic gate is preserved under `build82-diagnostics-failed/`. TSV header and row field counts are checked explicitly before publication.

## Promotion rule

One fast run is not sufficient. Build82 may replace Build76 only if **two independent physical runs**:

- pass the complete 9/9 semantic matrix;
- preserve the qualified deep counters and protected-data outcomes exactly;
- show a material, repeatable performance benefit, especially in B/mild/B/angle geometry and generation four;
- do not introduce regressions elsewhere in the matrix.

If Build82 is exact but not materially faster, do not add another cache/LUT layer. The next step is a dedicated benchmark/pprof measurement; manual `mapPoint` inlining is considered only after that evidence and must receive its own bit-exact gate.

## Final physical result

Both independent physical runs passed 9/9 and reproduced the qualified deep counters exactly. Run 1 / run 2 full-matrix elapsed was 933,537 / 729,213 ms, giving an 831,375 ms mean versus the 758,102 ms Build76 mean. B/mild generation four was 24,305 / 19,991 ms versus 18,473 ms Build76 mean; B/angle generation four was 136,040 / 107,359 ms versus 110,830 ms. The second run alone contains some faster measurements, but the pair is not a material repeatable speedup.

**Decision: do not promote Build82.** Restore the active runtime to Build76 and continue with Build83 observability-only benchmark/pprof work.
