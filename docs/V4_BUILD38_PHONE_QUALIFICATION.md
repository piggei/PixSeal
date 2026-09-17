# Format-v4 Build38 smartphone qualification

Build37 closed the first blind physical scanner path on the Build35 robust/strength-24 paper set. Build38 starts a separate smartphone-camera qualification channel. It does **not** change the Format-v4 wire format, locked pilot, 1120 data positions, Hamming(7,4), whitening/HMAC domains or frozen Format-v3 implementation.

## Strength-24 phone corpus result

The first phone corpus contains nine native camera JPEGs at approximately 200 MP: control, marked A and marked B, each acquired in front, mild-perspective and stronger-perspective conditions. These files are private and are not shipped.

Build37's scanner boundary detector is not a general camera detector: in several phone frames it can confuse the paper edge with the printed artwork. More importantly, geometry is not the only blocker. Reference-assisted diagnostics deliberately allowed the digital original to provide a homography and then a dense non-rigid residual registration. This is **not** a production decoder and is used only to measure channel capacity.

Representative strength-24 results after dense reference-assisted registration:

| capture | coded-bit errors | coded BER | post-soft-Hamming errors |
| --- | ---: | ---: | ---: |
| marked A, front | 80 / 448 | 17.86% | 42 / 256 |
| marked B, angle | 65 / 448 | 14.51% | 20 / 256 |

These values remain far outside exact frame/HMAC recovery. Pilot-only per-tile warp searches can produce high apparent local pilot correlation by fitting natural image texture, while the protected data plane remains weak. Build38 therefore does not lower acceptance thresholds and does not promote that overfit as a decoder.

## Build38 carrier decision

The next phone experiment changes only embedding strength from 24 to **48**. On the canonical MQ qualification carrier, digital comparison against the unmarked control gives approximately 32.49 dB PSNR. The change remains visually subtle in the qualification image while doubling the DCT separation requested by the encoder.

Generate the private pack with:

```bash
make v4-phone-fixtures
```

Default output:

```text
v4-phone private/build38-generated/
```

The pack contains one control plus two different authenticated robust carriers (`v4-b38-phone-a`, `v4-b38-phone-b`), a SHA-256 manifest, an acquisition plan and `README.txt`.

## Acquisition

Print at actual size / 100%, 300 ppi, using the same smooth white paper and printer settings as the successful scanner campaign. For each sheet acquire three native phone files: front, mild perspective (~10-15 degrees) and angle (~25-30 degrees). Use the main 1x camera, maximum native resolution, no digital zoom, no flash, no document mode and no messaging-app recompression. Keep the full artwork and visible white paper around all four sides.

A camera PASS requires exact HMAC-authenticated payload recovery. Pilot correlation alone is never authentication.

## Next decision

First use independently supplied/reference-assisted geometry on the new strength-48 photographs. If the data channel now reaches exact HMAC recovery, the next build can concentrate on blind camera boundary/projective/residual registration. If it still does not, the evidence will justify revisiting physical strength/ECC rather than hiding the problem behind a more permissive geometry search.
