# Format-v4 Build41 blind smartphone basin recovery

Build41 closes the first staged blind-smartphone milestone on the existing private strength-48 physical corpus. It does **not** change the Format-v4 encoder, wire format, locked public pilot, 1120-position data partition, Hamming(7,4), whitening/HMAC domains, robust strength 48, canonical 1632x1632 carrier, Build40 residual model, or frozen Format-v3 implementation.

The Build40 nine-photo diagnostic matrix showed that the remaining marked failures happened before protected-data recovery: the decoder was losing the correct global projective basin. Build41 therefore changes only the blind smartphone global-registration front end.

## Geometry discipline

The central invariant remains unchanged:

- image structure and the public pilot may determine geometry;
- secret key, payload/header contents, ECC result and HMAC may not generate, refine, rank or select geometry;
- HMAC is evaluated only after the geometry ensemble has been frozen and the protected channel has been decoded.

Build41 uses one deterministic spatial split of the public pilot. Tile fold `0` of `(tx + 2*ty) mod 3` is permanently held out; folds `1` and `2` are proposal evidence. The held-out fold is not sampled during boundary refinement, cyclic-shape fitting, absolute-phase search, fine proposal refinement or shortlist construction.

## Basin search

The phone path now performs:

1. bounded internal downsample, retaining the Build39 long-side limit;
2. structure-only artwork-boundary selection, preferring the dedicated inner-artwork detector and falling back to the established phone/component/scanner structural detectors when required;
3. an eight-coordinate corner refinement bounded to +/-56 working-image pixels around the structural anchor;
4. cyclic public-pilot fitting on proposal folds only, so small absolute phase error does not pull the projective shape into a wrong whole-tile alias;
5. bounded translation search that restores physical origin `(0,0)` using proposal evidence only;
6. proposal-only fine variants for the three strongest phase candidates;
7. freeze of the complete geometry shortlist;
8. held-out qualification followed by complete-pilot `(0,0)` origin, score and margin gates;
9. selection of a two-hypothesis ensemble using only proposal, held-out and complete-pilot evidence.

The ensemble deliberately rejects near-duplicate geometries. Its two members must be separated by at least `max(0.75 px, 0.20 * observed DCT-block scale)`. This small scale-normalized diversity brackets sub-pixel registration uncertainty instead of averaging two copies of the same local optimum.

Only after this ensemble is frozen may Build40's bounded residual field run. The residual remains independently proposal/held-out qualified and capped at six canonical pixels. If it does not qualify, the Build41 global mappings are left unchanged.

## Synthetic qualification

`make v4-build41-phone-basin-test` creates a strength-48 carrier, places it into a strongly projective synthetic phone page, and requires the public API `ExperimentalV4ExtractPhone` to recover the exact HMAC-authenticated payload. The matching unmarked control must reject. The test exercises the same Build41 basin front end and Build40 residual/data path exposed by `v4-extract-phone`.

## Private strength-48 physical matrix

The canonical private corpus remains the nine original ~200 MP JPEGs acquired for Build38. The photographs remain private and are not shipped in source/release archives.

| capture | boundary | basin | ensemble | proposal | held-out | full pilot / margin | data decode | HMAC |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| control/front | yes | no | 0 | - | - | - | no | REJECT |
| control/mild | yes | no | 0 | - | - | - | no | REJECT |
| control/angle | yes | no | 0 | - | - | - | no | REJECT |
| A/front | yes | no | 0 | - | - | - | no | FAIL |
| A/mild | yes | yes | 2 | 0.278123 | 0.191022 | 0.242107 / 0.130763 | yes | FAIL |
| A/angle | yes | yes | 2 | 0.386485 | 0.270798 | 0.321164 / 0.181862 | yes | **PASS: `v4-b38-phone-a`** |
| B/front | yes | yes | 2 | 0.329084 | 0.266604 | 0.295112 / 0.149755 | yes | **PASS: `v4-b38-phone-b`** |
| B/mild | yes | no | 0 | - | - | - | no | FAIL |
| B/angle | yes | no | 0 | - | - | - | no | FAIL |

The Build40 residual fitter was attempted on A/mild, A/angle and B/front. It was **not applied** in any of them because held-out residual evidence did not justify the correction. Therefore both physical HMAC PASS results come from Build41 global registration plus the unchanged Build36 soft data channel, not from a residual overfit or HMAC-guided fallback.

## Staged physical success

Build41 satisfies the smartphone phase's first physical milestone:

- all three controls reject;
- at least one A capture authenticates blind;
- at least one B capture authenticates blind.

The private target `make v4-build41-phone-physical-test` reproduces this rule over all nine files and writes a TSV/Markdown matrix plus per-image stderr telemetry under the Git-ignored Build41 diagnostics directory.

Build41 does **not** claim 6/6 marked robustness. The remaining envelope is explicit: A/front, B/mild and B/angle still need a better basin; A/mild already reaches qualified geometry/data sampling but still fails final authentication. Future work must preserve the three control rejections and the existing A/B physical PASS cases before broadening robustness.
