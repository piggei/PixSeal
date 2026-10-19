# Format-v4 Build95 — exact generation-four luminance-plane candidate

Build95 is a **benchmark-first runtime candidate** selected from the completed Build94 profile of the current qualified Build93 smartphone decoder. It is not qualified, does not replace the normal `v4-extract-phone` path, and cannot be promoted from source tests or microbenchmarks alone.

The candidate changes one narrow geometry kernel: projective block reads inside generation-four `continuation4` FoldScore. It prepares one immutable exact `float64` luminance plane immediately before generation 4 and uses it for those repeated bilinear samples. All other geometry, qualification and protected-data behavior remains Build93.

## Selection evidence from Build94

Build94 completed its Go 1.26.0 source and profiling gates with exact retained semantics. The low-overhead timing pass measured:

| case | deep total | geometry | geometry share | freeze | gen3 | gen4 | qualification | decode |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| B/mild | 80,152 ms | 57,175 ms | 71.33% | 21,768 ms | 9,906 ms | 17,463 ms | 12,511 ms | 10,251 ms |
| B/angle | 178,627 ms | 176,938 ms | 99.05% | 20,536 ms | 36,072 ms | 106,635 ms | 1,453 ms | 0 ms |

The separate B/angle CPU profile sampled 1,384.63 CPU-s over 191.99 s wall time. Generation 4 accounted for 920.46 CPU-s cumulative (66.48%), and the exact Build84 continuation4 reader alone accounted for 456.30 CPU-s cumulative (32.95% of the complete profile). `experimentalV4PhoneBuild84ContinueFetch` accounted for 459.29 CPU-s. Qualification was only 10.06 CPU-s and protected decode was not entered.

B/mild is less concentrated in this target: generation 4 accounted for 135.96 of 513.24 sampled CPU-s, and the Build84 continuation4 reader for 51.60 CPU-s (10.05%). Build93 protected residue reuse remains effective: B/mild decode was 10,251 ms and the exact `696` sampling passes / `2,086` profile derivations were preserved.

This evidence rules out another scheduling change as the immediate next step: the major generation pools already use eight workers on the eight-logical-CPU qualified host, and worker-time/wall-time ratios show those pools are already saturated. It also rules out qualification or protected decode as the B/angle target. Build95 therefore tests one structural amortization hypothesis around the specific reader newly exposed by Build94.

The source-controlled numeric snapshot is `docs/qualified-baselines/build94-profile-summary.tsv`. Raw `.pprof`, retained images, keys, payloads and diagnostic logs remain private.

## Why this is not Build80 or Build81 repeated

Build80 built a local cache per FoldScore. Every block first paid a 64-point mapping pre-pass, and fallback blocks repeated work in the historical reader; B/angle recorded no useful cache hits and performance regressed. Build81 replaced scalar products with LUT lookups but did not remove repeated sampling work and produced no repeatable gain.

Build95 is structurally different:

- it computes each source pixel's exact luminance once for the complete deep geometry run, not once per FoldScore or block;
- it adds no per-block cache lookup, source bounding pass or geometry-dependent cache footprint;
- projective mapping, bounds checks, clamping, bilinear interpolation and DCT accumulation still execute in their original order;
- every generation-four continuation worker reads the same immutable plane concurrently;
- a deterministic public-dimension memory cap chooses between the candidate reader and the unchanged Build84 reader before protected data exists.

The memory cost is explicit rather than hidden. On the retained `3464 x 4600` work image, the additional plane contains `15,934,400` `float64` values, or `127,475,200` bytes (about 121.57 MiB). The cap is `134,217,728` bytes (128 MiB), leaving about 6.43 MiB of headroom for this exact retained size. The existing RGB plane and work image are not included in that additional-plane number.

## Exact candidate behavior

`experimentalV4PhoneBuild95NewLuminancePlane` traverses the immutable Build84 RGB plane once and evaluates, for every pixel, the exact historical expression in the same order:

```text
.299*float64(R) + .587*float64(G) + .114*float64(B) - 128
```

The candidate reader then retains the historical sequence:

1. `homography.mapPoint` for every one of the 64 block samples;
2. the same out-of-bounds rejection;
3. `math.Floor`, `x1/y1` edge clamps and `fx/fy` fractions;
4. the same four luminance loads;
5. horizontal and vertical bilinear interpolation in the same expression order;
6. the same `c23` and `c32` DCT multiplication/accumulation order;
7. `abs(c23) - abs(c32)`.

Only RGB-to-luminance multiplication is moved from every projective sample to the one-time plane construction. No approximate conversion, smaller scalar type, SIMD/unsafe code, coordinate recurrence or floating-point reassociation is used.

## Scope boundary

Build95 uses the luminance plane only in `experimentalV4PhoneBuild95ContinueFetch`, the continuation portion of generation 4. The following remain on the exact qualified implementation:

- Build47/75 geometry freeze and its public-pilot proposal work;
- prefix1;
- generation 2;
- generation 3;
- generation-four single probes and pair stencils;
- candidate ordering and ordered commits;
- Build90 ordered-parallel qualification;
- Build93 protected residue sampling and profile derivation;
- Build66 batch ordering and first-logical-HMAC semantics;
- encoder, wire format, locked pilot, strength 48, ECC/Hamming, whitening/HMAC domains, thresholds and quorum.

The complete geometry bank remains frozen before qualification or protected decode. No key, payload, ECC, HMAC or message correctness signal can affect luminance-plane creation, fallback, geometry proposal, ranking, pruning or scheduling.

## Bounded fallback

The plane is attempted only after generation 3 has produced generation-four tasks. Its requested byte count is derived solely from the public work-image width and height.

Build95 uses the unchanged Build84 continuation4 reader when:

- the RGB plane is absent or invalid;
- dimensions are empty;
- RGB storage length is inconsistent;
- the exact `width * height * 8` request exceeds 128 MiB.

Fallback does not change task order, evaluation count, candidate order or telemetry. The cap is a portability and memory-safety policy for this candidate, not a quality or secret-dependent heuristic.

## Telemetry

The existing Build84/Build93 telemetry is preserved. Build95 adds only:

- `Build95Attempted`;
- `Build95LuminancePlaneUsed`;
- `Build95LuminancePlaneFallback`;
- `Build95LuminancePlaneBytes`;
- `Build95LuminancePlanePrepMs`.

Private paired evidence additionally records generation-four wall time beside plane preparation and byte count. These values are diagnostic only and cannot affect recovery.

## Source and race gate

```bash
PIXSEAL_GO_TOOLCHAIN=go1.26.0 \
  make v4-build95-gen4-luminance-test
```

The gate proves:

- every precomputed luminance value is bit-identical to the Build84 expression;
- projective block outputs and failures match Build84 bit-for-bit across exact, skewed, boundary and invalid cases;
- FoldScore output, visible count and fetch telemetry are exact;
- continuation states/order/evaluations are exact;
- generation-four output bank/order/evaluations are exact, including explicit fallback;
- complete synthetic blind-bank output is exact;
- public entry validation remains Build93-equivalent;
- Build93 telemetry remains present when Build95 telemetry is added;
- concurrent continuation reads of one immutable plane pass `go test -race`;
- CLI help, `VERSION`/buildinfo and `go vet` pass.

## Public benchmark gate

```bash
PIXSEAL_GO_TOOLCHAIN=go1.26.0 \
  make v4-build95-gen4-luminance-benchmark
```

The public deterministic benchmark measures, in the same process/session:

- Build84 block reader versus Build95 luminance-plane block reader;
- Build84 continuation versus Build95 continuation;
- one-time luminance-plane preparation cost on both the compact public fixture and a public `3464 x 4600` retained-size plane.

It uses no private acquisition, key, payload, ECC or HMAC. Correctness and race are hard gates; timing is decision evidence. The retained paired test should not be run merely because the candidate is exact. It is justified only if repeated Go 1.26.0 samples show a material reader/continuation improvement large enough to plausibly repay preparation and additional memory.

## Paired retained decision gate

If the public benchmark is positive:

```bash
PIXSEAL_GO_TOOLCHAIN=go1.26.0 \
  make v4-build95-paired-profile
```

The harness executes two contemporaneous AB/BA cycles:

```text
cycle 1: B/mild Build93 -> Build95
cycle 1: B/angle Build93 -> Build95
cycle 2: B/angle Build95 -> Build93
cycle 2: B/mild Build95 -> Build93
```

Each pair reuses the same decoded in-memory work image. Before timing is considered, both variants must have exactly equal outcome, payload, extract info and frozen semantic telemetry. Build95 must also report the expected plane byte count, `used=true` and `fallback=false` on the retained cases.

The candidate is worth physical qualification only if both orderings show a directionally stable, material generation-four/geometry gain after including plane preparation. A one-sided or host-drift-sensitive result closes Build95 as exact but non-promoted.

## Promotion discipline

Build93 remains the current qualified smartphone baseline throughout Build95 evaluation. The normal command remains:

```text
v4-extract-phone -> Build93
```

The explicit research command is:

```text
v4-extract-phone-build95 -> Build95 candidate
```

A favorable public benchmark and paired profile still do not promote Build95. Promotion requires:

1. exactness/source/race PASS;
2. public same-session benchmark PASS;
3. contemporaneous retained AB/BA evidence showing a stable material benefit;
4. two independent complete Go 1.26.0 nine-photo physical semantic PASS runs through the explicit candidate path;
5. documentation and production-entry-point consolidation only after those gates.

Until then, Build95 is reproducible research evidence and Build93 remains authoritative.
