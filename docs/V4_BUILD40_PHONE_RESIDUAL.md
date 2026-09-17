# Format-v4 Build40 smartphone residual-warp checkpoint

Build40 reuses the existing Build38 strength-48 printed/photographed corpus. It does not change the Format-v4 pilot, data partition, framing, Hamming(7,4), whitening/HMAC domains, robust strength, print dimensions or frozen Format-v3 implementation.

## Question

Build39 supplies a perspective-capable global canonical-to-phone homography. Build40 asks a narrower question: **if that homography is already in the right basin, can a small smooth residual deformation be recovered from public pilot evidence without using the protected frame as an oracle?**

The residual model is intentionally not another free camera transform. X and Y are represented by six normalized quadratic terms (`1`, `x`, `y`, `xy`, `x²`, `y²`) and the resulting correction is capped at six canonical pixels.

## Proposal / validation separation

The repeated v4 pilot tiles are split spatially as a checkerboard. Only checkerboard-A tiles can create local displacement controls. A local control is retained only when its bounded +/-6 px search has both a minimum pilot score and a separated runner-up. A robust quadratic fit then discards controls inconsistent with the smooth field.

Checkerboard-B tiles are never used to fit the field. They are evaluated only after the field is frozen. The field must improve held-out validation and the complete pilot must subsequently recover cyclic origin `(0,0)` with explicit score and margin floors. Payload bytes, frame header, ECC result and HMAC do not participate in this process.

## Data sampling

If the residual field qualifies, protected DCT margins are sampled through the composed mapping:

`canonical point -> bounded residual correction -> Build39 homography -> observed phone image`

The existing Build36 soft-Hamming decoder and unchanged v4 HMAC then operate normally. `ExperimentalV4ExtractPhone` keeps a copy of the untouched Build39 ensemble; if a residual-qualified bank does not authenticate, the original Build39 bank is retried once. HMAC therefore verifies fixed geometries but never chooses or ranks them.

## Qualification

`make v4-build40-phone-residual-test` synthesizes a strength-48 carrier with a deterministic smooth non-projective deformation. Build40 must recover a residual field from proposal tiles, improve held-out validation, return complete-pilot origin `(0,0)`, recover the exact HMAC-authenticated payload, and reject the matching unmarked control.

## Real phone corpus result

Representative current private strength-48 captures do **not** qualify the residual field from the blind Build39 basin. Local pilot maxima can be strong, but their corrections are not jointly supported by held-out tiles: representative fits raise proposal evidence while held-out validation falls and the complete pilot retains a non-zero origin. Build40 therefore safe-rejects the field rather than increasing the bound or feeding HMAC results back into registration.

The next research step is global projective-basin preservation/ranking on the same existing phone corpus. No reprint or additional acquisition is justified by Build40.
