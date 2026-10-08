# Format-v4 Build88 — exact direct RGB loads / dominating BCE

Build88 is a **closed exact / benchmark-strongly-negative / non-promoted experiment** over the current qualified Build84 smartphone baseline. It is selected from the completed Build87 compiler/profile evidence and changes only the four RGB pixel accesses inside continuation4 FoldScore projective block reads.

## Evidence selecting the experiment

On the qualified Go 1.26.0 / Intel Core i3-13100T host, Build87 measured the qualified Build84 reader at five-run means of approximately **1394.8 / 1374.8 / 1402.6 ns/block** for front/mild/angle. The prepared angle-like stage means were approximately **236.9 ns/block** for the four three-byte RGB slice fetches, **337.3 ns/block** for four RGB-to-luminance conversions, **144.8 ns/block** for horizontal bilinear interpolation, **131.3 ns/block** for vertical interpolation, and **514.6 ns/block** for the combined prepared fetch+luminance+bilinear path.

The Build87 Go compiler report shows exactly eight `IsSliceInBounds` checks on the four qualified Build84 pixel slices: two checks on each of the four `rgb[index:index+3]` expressions. CPU pprof confirms that the qualified reader is still dominated by the exact RGB/luminance/bilinear kernel while `homography.mapPoint` remains a small cumulative component.

## Single Build88 hypothesis

Build84 reads each RGB pixel through a three-byte slice:

```go
pixel00 := rgb[index00 : index00+3]
l00 := .299*float64(pixel00[0]) + .587*float64(pixel00[1]) + .114*float64(pixel00[2]) - 128
```

Build88 uses one explicit dominating proof for the highest byte followed by direct scalar byte loads:

```go
_ = rgb[index00+2]
r00, g00, b00 := rgb[index00], rgb[index00+1], rgb[index00+2]
l00 := .299*float64(r00) + .587*float64(g00) + .114*float64(b00) - 128
```

The same shape is used independently for the other three pixels. This is intended to let Go 1.26.0 prove the three channel loads from one dominating bounds check per pixel instead of constructing four three-byte slices.

## Invariants

Build88 does **not** change:

- `homography.mapPoint`, `math.Floor`, clamp or fractional-coordinate arithmetic;
- byte offsets or the RGB values read;
- `.299*R + .587*G + .114*B - 128` evaluation order;
- horizontal or vertical bilinear evaluation order;
- DCT table accesses, multiplication order or c23/c32 accumulation order;
- single4/pair4 or generations 1–3;
- Build75 freeze, Build76 scheduling/ordered commit, Build68 qualification reuse or protected-data decode;
- pilot, strength 48, ECC/Hamming, whitening/HMAC domains, thresholds or Format-v4 wire format.

There is no `unsafe`, LUT, cache, precomputed luminance plane, FMA, arithmetic reassociation or coordinate recurrence.

## Qualified-host result

The Go 1.26.0 i3-13100T exactness gate passed, but the same-session performance gate failed decisively. Five-run Build84/Build88 means were **1415.4/1529.2 ns** front, **1434.0/1593.8 ns** mild and **1450.8/1611.8 ns** angle. Build88 is therefore about **8.0% / 11.1% / 11.1% slower**. The isolated prepared fetch also failed to improve: 255.68 ns for the Build84 slice form versus 257.46 ns for direct scalar loads. No private physical Build88 run is justified.

**Decision:** keep Build84 qualified, close Build88, and move to Build89 full-pipeline profiling.

## Gates

Before any private physical run, Build88 must pass:

1. `math.Float64bits` equality against the qualified Build84 reader over broad projective/border/failure cases;
2. FoldScore bit equality;
3. continuation4 state/evaluation equality;
4. complete blind bank equality and workload telemetry preservation;
5. CLI/buildinfo/version checks and `go vet ./...`;
6. same-session Go 1.26.0 Build84-vs-Build88 reader benchmarks on front/mild/angle;
7. Go BCE/inlining and objdump evidence showing what actually changed in generated code.

Because Build87 attributes only a small single-digit percentage of the reader to the slice construction/check region, a 5% whole-reader threshold is unrealistic. The candidate should reach the private corpus only if the complete-reader benchmark shows a **stable useful improvement, approximately 2% or better on all three geometries, with no regression**, and compiler evidence confirms a simpler bounds-check shape.

Build84 remains the current qualified smartphone baseline until two independent physical 9/9 runs prove exact semantics and a repeatable worthwhile end-to-end improvement.

## Commands

```bash
make clean
make v4-build88-direct-rgb-test
make v4-build88-direct-rgb-benchmark
```

Only after the benchmark gate is accepted:

```bash
make V4_PHONE_BUILD88_DIAGNOSTIC_DIR='v4-phone private/build88-diagnostics1' v4-build88-phone-physical-test
```
