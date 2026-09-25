## v0.3.0-build44 deterministic JPEG / Go 1.26 qualification — 2026-09-23

Build44 physically closes the toolchain reproducibility issue discovered during Build43. The PixSeal-owned `internal/jpeglegacy` decoder reproduces the frozen pre-Go-1.26 raster contract when the project itself is compiled with Go 1.26.0.

On the Surface/WSL2 qualification host:

```text
make v4-build44-go126-jpeg-compat-test            PASS
make v4-build44-go126-phone-physical-test        PASS
controls                                          3/3 REJECT
A/front                                           PASS
A/mild                                            PASS
A/angle                                           PASS
B/front                                           PASS
B/mild                                            INFO/reject
B/angle                                           INFO/reject
```

No encoder, Format-v4 pilot/data mapping, strength, Hamming/ECC, whitening/HMAC domain, Build41/43 geometry or Build42 data/list behavior changed. Go 1.26.0 is promoted as the qualified Build44 toolchain.

# PixSeal v0.2.0 — qualification and measured results

This document separates the **stable release baseline** from experimental
geometry measurements. Results are corpus-specific observations, not universal
recovery guarantees.

## RC2 qualification reference — 2026-09-09

The preserved `v0.2.0-rc2` all-test report was produced on Linux/WSL2 with Go
1.25.1 and ImageMagick, using strict mode for pass/fail baseline and geometry
suites. It reported:

```text
Release baseline: PASS
Research suites: ATTENTION (2 experimental targets failed)
Overall: FAIL (2 experimental targets failed)
```

Release baseline:

```text
version-check          PASS
vet                    PASS
release-unit           PASS
test-images            PASS
deep-test              72/72 PASS
core-target-check      PASS
```

Research measurements:

```text
extreme-test           completed progressive limit mapping
research-unit          PASS
geometry-test          65/120 PASS, 55 FAIL
affine-test            22/24 PASS, 2 FAIL
composition-test       2/2 PASS
lattice-test           8/8 PASS
perspective-test       4/4 PASS
```

The geometry/affine FAIL rows are intentionally visible measurements of the
experimental search bank. They are not failures of the stable authenticated
JPEG/resize/crop release baseline.

## Baseline digital robustness

On the RC2 qualification corpus, `deep-test` recovered all 72 tested cases
across two real images and all three profiles. The suite covered JPEG quality 82,
resize 95/85/75/65/55/50%, center crop 90/75/50% and deterministic random crop
75/50%. `extreme-test` is intentionally non-strict and maps the progressive
failure envelope.

## Adaptive profiles

| Profile | Max payload | Protected bits | Tile redundancy |
|---|---:|---:|---:|
| robust | 16 B | 448 | 2.50x |
| balanced | 32 B | 672 | 1.67x |
| capacity | 64 B | 1120 | 1.00x |

## Resize recovery budget

```text
13 fixed non-direct scale hypotheses
<= 3 phases per scale
<= 6 full-grid shortlisted scales
<= 18 full-carrier virtual aggregations
```

A decisive scale may receive one physical bilinear fallback over the expected
inverse size and its eight ±1-pixel neighbours.

## Experimental lattice composition budget

The frozen bank contains 110x90, 90x110, 105x95 and 95x105 basis shapes over
-45..+45 degrees at 0.25-degree steps. The search is evaluated in two sequential
shape groups (±10% anchors, then ±5% moderates):

```text
1444 sparse probes maximum
192 stronger coherence evaluations maximum
12 full authenticated grids maximum per group
24 full authenticated grids maximum overall
```

The RC2 qualification recovered 8/8 configured robust-profile lattice cases.

## Experimental perspective

The frozen v0.2 bank contains only `top-narrow-4` and `bottom-narrow-4`, one
aligned full projective aggregation each. RC2 recovered 4/4 configured cases.
RC2 already contained an end-to-end Go recovery regression; RC3 strengthened the
synthetic warp implementation while retaining explicit correction-path assertion.
General homography inference and physical print-camera recovery are not claimed.

## Final v0.2.0 qualification — 2026-09-09

The RC4 source candidate was exercised on the same private two-image corpus used
for the RC2 reference. The complete all-test report recorded:

```text
version-check          PASS
vet                    PASS
release-unit           PASS
test-images            PASS
deep-test              72/72 PASS
extreme-test           completed progressive limit mapping
research-unit          PASS
geometry-test          63 PASS, 55 FAIL, 2 TIMEOUT (60-second run)
affine-test            22/24 PASS, 2 FAIL
composition-test       2/2 PASS
lattice-test           8/8 PASS
perspective-test       4/4 PASS
core-target-check      PASS

Release baseline: PASS
Research suites: ATTENTION
```

The two RC4 geometry timeouts were both `PJ_lingua.PNG`, profile `capacity`, at
-10 and -5 degrees. A targeted rerun with:

```text
STRICT=1 EXTRACT_TIMEOUT=120 make geometry-test
```

recovered both cases and produced:

```text
65 passed, 55 failed, 0 timeouts, 0 skipped, 0 errors, 120 total
```

This exactly reproduces the RC2 experimental geometry reference. The timeout
difference is therefore treated as a wall-clock harness sensitivity, not a
decoder capability regression. The final source keeps the same decoder search
banks and uses a 120-second default only for the geometry research harness.

Final release conclusions:

- stable release baseline fully qualified;
- Format v3 golden fingerprints unchanged;
- no encoder or decoder-search functional change after RC4;
- experimental affine/composition/lattice/perspective counts match the RC2
  reference;
- geometry capability matches the RC2 65/120 reference when given the validated
  120-second research timeout;
- reusable core portability checks pass for linux/amd64, windows/amd64,
  android/arm64 and ios/arm64.

## Large images

The v0.2 line has round-tripped a private ~200 MP carrier. v0.2.0 retains the
300,000,000-pixel CLI/core source policy from RC2 while keeping RC3's
overflow-safe arithmetic. Geometry suites independently use lower research
limits (typically 50/100 MP). The 300 MP guard is a safety policy, not a promise
that every machine has sufficient RAM for every image below it.

## Interpretation

- **Deterministic:** format sizes, profile capacities, search-space caps and
  authentication rules.
- **Experimental:** success rates on particular images/transforms.
- **Heuristic:** `analyze` detail score and recommended strength.

No experimental PASS guarantees recovery on an arbitrary future image, and no
failed experimental geometry case invalidates the authenticated digital baseline.

## v0.3.0-build1 research checkpoint — 2026-09-09

The build1 change is isolated from Format v3, embedding and `ExtractWithInfo`; the
v0.2 qualification numbers above therefore remain the stable comparison
baseline. New diagnostic tests add two explicit properties:

- a deterministic marked synthetic carrier must produce bounded local-lattice
  evidence while its corresponding unmarked source remains negative;
- optional correct-key authentication uses the existing HMAC path, while a wrong
  key cannot turn lattice evidence into an authenticated payload.

Local manual qualification on the private digital corpus is recorded only as a
research observation, not a release guarantee. With robust-profile marked copies
of the two qualification images, the 2026-09-09 build1 run measured:

| Input | Lattice evidence | Global consistency | Global `u` | Global `v` |
|---|---:|---:|---:|---:|
| marked `PJ_piccolo.png` | true | 0.939 | (7.95, -0.02) | (-0.02, 8.03) |
| unmarked `PJ_piccolo.png` | false | 0.414 | (5.50, -5.30) | (4.70, 5.23) |
| marked `PJ_lingua.PNG` | true | 0.981 | (8.07, 0.00) | (0.00, 8.00) |
| unmarked `PJ_lingua.PNG` | false | 0.442 | (6.31, -1.91) | (1.76, 6.50) |

The marked canonical carriers therefore resolve very close to the native 8-pixel
lattice while the corresponding unmarked originals remain below the evidence
gate. Exact values are content dependent and are not frozen as interoperability
requirements.

The real `foto stampa.jpg` / `foto stampa storta.jpg` corpus is intentionally not
distributed. `make print-camera-test` therefore SKIPs when that private corpus is
not installed. When it is installed, only an authenticated Format v3 HMAC can
produce PASS; useful lattice evidence without HMAC is reported as research, not
success.

Known build1 limit: automatic arbitrary-angle local-basis estimation is not yet
reliable enough to fit a general homography. Synthetic probing shows measurable
repetition evidence near the true transformed lattice, but candidate generation
and bounded refinement need further work before the homography milestone.
## v0.3.0-build2 projective research checkpoint — 2026-09-09

Build2 remains isolated from Format v3, embedding and production
`ExtractWithInfo`. Its deterministic additions are tested separately through the
research targets. The synthetic projective test verifies that the virtual
homography sampler can reach a valid existing Format-v3 HMAC with the correct
key and that a wrong key fails.

On the two undistributed original smartphone photographs, the current automatic
diagnostic measured:

| Input | Boundary confidence | Global lattice consistency | Lattice evidence | Real HMAC |
|---|---:|---:|---:|---:|
| `foto stampa.jpg` | 0.728 | 0.760 | true | not recovered |
| `foto stampa storta.jpg` | 0.652 | 0.812 | true | not recovered |

For comparison, build1 measured about 0.637/true on the frontal capture and
0.505/false on the inclined capture. Build2 therefore materially improves the
geometric initializer on the harder perspective case. These values are research
observations on two images, not frozen thresholds or interoperability promises.

The bounded projective path probes at most eight canonical-size candidates and
permits at most four full virtual v3 decodes. On the real captures it reaches
that path without materializing a full rectified image, but no candidate yet
authenticates. Key-known header fractions/z-scores are deliberately not reported
as watermark detection because the bounded candidate search still creates a
multiple-testing effect. The hidden message remains unknown.

The qualified v0.2 release baseline remains the reference: build2 changes only
research diagnostics and does not promote these hypotheses into
`ExtractWithInfo`.


## v0.3.0-build3 phase-refinement checkpoint — 2026-09-10

Build3 keeps the build2 boundary/lattice results and adds spatial phase
consensus plus bounded scale/DLT refinement. On the two private ~200 MP
smartphone captures, the automatic diagnostic still reports global lattice
consistency about 0.760 (frontal) and 0.812 (inclined), with lattice evidence on
both.

A notable cross-capture observation is a closely matching coarse scale family:
about 1668x1250 on the frontal capture and 1653x1254 on the inclined capture.
These values are produced independently by the estimator and are not coded as
target dimensions. The frontal member has phase consistency about 0.716; on the
inclined capture the corresponding family is about 0.575 and a bounded height
refinement near 1653x1241 reaches about 0.615. Other local maxima remain and are
retained rather than silently discarded.

For the frontal family, a four-point phase-DLT refinement raises its sparse
known-header probe from roughly z=3.50 to z=4.45. On the inclined capture, the
bounded refinement bank raises the best sparse probe from the build2 value near
z=3.70 to about z=3.97. These are multiple-testing-affected research scores, not
watermark detection.

Both photographs still exhaust at most four virtual full-grid attempts without
a valid Format-v3 HMAC. `print-camera-test` therefore remains a deliberate
research failure and the hidden message remains unknown.


## v0.3.0-build4 fundamental-scale checkpoint — 2026-09-10

Build4 keeps the build3 lattice evidence (global consistency about 0.760 frontal
and 0.812 inclined) and adds an explicit fundamental-scale selector before
key-assisted phase refinement. The selector independently promotes about
1668.02x1250.43 on the frontal capture and 1653.11x1253.88 on the inclined
capture. These values come from multilevel/multiregion support and are not
private-corpus constants.

On the frontal capture, bounded fundamental-scale refinement does not beat the
seed (best fundamental z≈3.50). Modulo-8 refinement also keeps an unshifted best
origin for the strongest refined branch, while the overall known-header maximum
remains balanced z≈4.454, unchanged from build3.

On the inclined capture, the fundamental branch refines to about
1686.17x1278.96 and reaches z≈3.732. The overall known-header maximum remains
balanced z≈3.973, also unchanged from build3. The single lattice-derived residual
fits use nine controls, with RMS residual about 2.28 canonical pixels frontal and
1.24 inclined; neither fit authenticates.

Both real photographs still exhaust at most four complete virtual Format-v3
decodes without a valid HMAC. Therefore `print-camera-test` remains a deliberate
research failure and the hidden message remains unknown. Header z-scores are
reported only as bounded search diagnostics and are not watermark-detection
claims.


## v0.3.0-build5 photometric checkpoint — 2026-09-10

Build5 preserves build4 geometry and adds a three-view photometric experiment.
On `foto stampa.jpg`, the retained fundamental phase-DLT geometry improves from
raw balanced z≈4.454 / fraction≈0.768 to mild-highpass balanced z≈4.695 /
fraction≈0.783. Its normalized signed sync margin also rises from about 0.621 to
0.651. No HMAC is recovered.

On `foto stampa storta.jpg`, the photometric views do not exceed the build4 raw
global best z≈3.973. This is kept as a negative result rather than triggering
additional tuning. The full-decode ceiling remains four and the hidden payload
remains unknown.

Qualification-corpus handling is reported separately from release invariants in
`all-test`, allowing different workstations to add low/medium/high-resolution
images without redefining the stable Format-v3 release gate.


## v0.3.0-build6 lattice-first held-out checkpoint — 2026-09-10

Build6 was frozen after build5 and then evaluated on the new private bicycle
print-camera pair. No thresholds or target dimensions were derived from that pair
before the test.

| Input | Strong boundary | Adaptive 1/4 | Lattice consistency | Projective source | Best photometric sync z | HMAC |
|---|---:|---:|---:|---|---:|---:|
| `foto bici dritta.jpg` | false (0.337) | yes | 0.615 | lattice+weak-boundary | 4.454 (raw) | not recovered |
| `foto bici storta.jpg` | true (0.780) | no | 0.781 | print-boundary+lattice | 4.079 (mild-highpass) | not recovered |
| `foto stampa.jpg` | true (0.728) | no | 0.760 | print-boundary+lattice | 4.695 (mild-highpass) | not recovered |
| `foto stampa storta.jpg` | true (0.652) | no | 0.812 | print-boundary+lattice | 3.973 (raw) | not recovered |

The new frontal bicycle capture is the key build6 result: build5 stopped before
projective authentication at consistency≈0.533. The adaptive finer level raises
consistency to ≈0.615 and the lattice-positive weak quadrilateral opens the same
bounded projective/HMAC path used by strong-boundary captures. Four complete
decode attempts are made, but none authenticate.

Negative controls remain essential. The unmarked 16320x12288 qualification image
can itself reach diagnostic lattice evidence after adaptive escalation
(consistency≈0.592), but its boundary confidence is only ≈0.082 and therefore it
does not open the lattice-first projective fallback. A wrong-key run on the new
frontal capture can also produce a high bounded header z-score while still
failing HMAC. These observations are why lattice, phase and sync scores remain
research evidence only and never constitute watermark detection.


## v0.3.0-build7 protected-bit channel checkpoint — 2026-09-10

Build7 keeps the build6 geometry, build5 photometric bank and four-full-decode
budget fixed, then measures the protected-bit/ECC state of those same decode
candidates. Only the fixed three-byte v3 prefix is used as known data; the hidden
payload remains unknown.

| Input | Best known coded errors (42) | Known Hamming words >1 error (6) | Known header errors after ECC (24) | Wrong/correct margin ratio | HMAC |
|---|---:|---:|---:|---:|---:|
| `foto stampa.jpg` | 7 | 1 | 2 | 0.670 | not recovered |
| `foto stampa storta.jpg` | 11 | 3 | 6 | 1.088 | not recovered |
| `foto bici dritta.jpg` | 12 | 3 | 4 | 0.876 | not recovered |
| `foto bici storta.jpg` | 8 | 2 | 4 | 1.567 | not recovered |

The original frontal capture is the most promising reliability-aware case: five
of six known header codewords are already inside the hard Hamming correction
radius and its wrong known bits have materially lower DCT margin than correct
bits. The inclined bicycle capture is different: despite only eight known coded
errors, wrong bits are stronger on average than correct known bits, suggesting
that geometry/phase or optical distortion is still a more plausible limitation
than simple low-confidence channel noise.

These measurements are deliberately not extrapolated to unknown payload bits and
are not a detection criterion. None of the four captures authenticates.

`make print-camera-test` in build7 discovers every PNG/JPEG directly inside the
private corpus directory and reports a corpus summary. The private photographs
remain excluded from source archives.


## v0.3.0-build8 scanner/reliability checkpoint — 2026-09-10

Build8 keeps the four smartphone cases and adds two 600-dpi Canon/JFIF scanner JPEGs
(4960x7015, about 34.8 MP each). The supplied scans are physically rotated 90 degrees
in their pixel data and are intentionally tested as supplied. One scan corresponds to
the bicycle print and the other to the original lake/bicycle print.

| Acquisition | Lattice consistency | Best sync z | Known coded errors /42 | Known words >1 error /6 | Hard post-ECC errors /24 | HMAC |
|---|---:|---:|---:|---:|---:|---|
| original smartphone frontal | 0.760 | 4.695 photometric | 7 | 1 | 2 | not recovered |
| original smartphone inclined | 0.812 | 3.973 | 10 | 3 | 7 | not recovered |
| bicycle smartphone frontal | 0.615 | 4.454 | 10 | 3 | 6 | not recovered |
| bicycle smartphone inclined | 0.781 | 4.079 photometric | 8 | 2 | 4 | not recovered |
| scanner bicycle print (`0270_001`) | 0.869 | 4.320 | 8 | 2 | 4 | not recovered |
| scanner original print (`0270_002`) | 0.875 | 3.732 | 9 | 3 | 5 | not recovered |

The build8 full-grid block cap changes some large-candidate bit-channel aggregates
relative to build7 while leaving the underlying geometry/probe search and Format-v3
semantics unchanged. This is expected: large full grids now use a fixed spatial sample
instead of every canonical block.

The original frontal smartphone capture is the only current case that passes the
conservative reliability-decode gate. One of its four complete decode slots therefore
uses weighted Hamming, but the best known prefix remains two bits wrong after both
hard and soft decoding and HMAC still fails. On the scanner cases, weighted Hamming
does not improve the known prefix. Strong scanner lattice consistency therefore does
not imply a clean protected-bit channel; the remaining problem cannot be attributed
only to smartphone perspective/ISP processing.

Performance is materially bounded: both scanner images previously exceeded a 120 s
full diagnostic run because build7 entered the mature baseline extractor. Build8 skips
that pre-pass above 16 Mi pixels and records full diagnostic wall times around 10--11 s
on the development environment, with projective authentication itself about 3.4--3.6
s. No private acquisition file is part of the source archive.


## v0.3.0-build9 spatial-stability checkpoint — 2026-09-10

Build9 was evaluated without changing any of the build8 geometry, photometric,
full-grid or HMAC budgets. The six private acquisitions are not distributed. The
figures below refer to the spatial evidence attached to the candidate selected by
the existing hard bit-channel ordering, not a new spatially optimized candidate.

| Acquisition | Lattice consistency | Cells | key-independent tile sign agreement | unstable tile positions | hard post-ECC known header | known stable-wrong / mixed | global-phase majority post-ECC | local-phase cells matching global | mean local phase offset (blocks) | bounded local-phase majority post-ECC | HMAC |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `foto stampa.jpg` | 0.760 | 4 | 0.683 | 0.389 | 2/24 | 0 / 38 | 5/24 | 1/4 | 1.37 | 10/24 | no |
| `foto stampa storta.jpg` | 0.812 | 6 | 0.669 | 0.743 | 7/24 | 0 / 39 | 6/24 | 1/6 | 1.82 | 5/24 | no |
| `foto bici dritta.jpg` | 0.615 | 9 | 0.635 | 0.821 | 6/24 | 0 / 42 | 10/24 | 2/9 | 1.35 | 4/24 | no |
| `foto bici storta.jpg` | 0.781 | 4 | 0.682 | 0.397 | 4/24 | 1 / 41 | 12/24 | 1/4 | 1.06 | 10/24 | no |
| scanner `0270_001.jpg` | 0.869 | 6 | 0.659 | 0.773 | 4/24 | 0 / 39 | 10/24 | 3/6 | 1.04 | 8/24 | no |
| scanner `0270_002.jpg` | 0.875 | 6 | 0.655 | 0.783 | 5/24 | 0 / 41 | 9/24 | 1/6 | 1.48 | 5/24 | no |

The key-independent result is the central build9 observation: repeated tile
positions do not preserve a common sign reliably across the print surface. Mean
agreement remains only about 0.64--0.68 even on the two scanner acquisitions whose
macro-lattice consistency is about 0.87. This rules out the simple model in which
a small fixed subset of coded bits is systematically inverted everywhere.

Known-prefix evidence points in the same direction: almost every known protected
header bit is mixed across cells, whereas stable-wrong bits are essentially absent.
A simple cell majority is not a solution and often performs worse than the existing
DCT-margin aggregation, showing that spatial cells are neither equally reliable nor
perfectly registered.

The bounded local-phase oracle finds that most cells prefer a phase within one to
two blocks of the global phase. It can reduce known-prefix majority errors in some
cases (for example `foto bici dritta.jpg` from 10 to 4 post-ECC errors and the
original inclined photograph from 6 to 5), while making other cases worse (the
original frontal photograph from 5 to 10). Because 25 key-known phase positions are
compared per cell, these improvements are not authentication evidence and are not
used for HMAC attempts. The result motivates a future smooth, low-degree-of-freedom
spatial phase/warp model rather than independent per-cell re-phasing.

Build9 also makes the build8 soft-Hamming summary candidate-paired. In particular,
`best_hard_candidate_soft_*` refers to the same candidate that produced the best
hard bit-channel result, while `best_soft_candidate_hard_*` reports the hard error
count on the candidate independently selected by the soft metric. This removes the
ambiguous cross-candidate comparison exposed by the build8 report.

No private acquisition produces a valid Format-v3 HMAC.


## v0.3.0-build10 smooth-field checkpoint — 2026-09-10

Build10 was evaluated on the same six private print-acquisition cases used by
build9. No private image is distributed with the source archive. All figures below
refer to the final build10 code and the unchanged physical-corpus test key `Piccotti`. `Best hard` is the
best hard bit-channel candidate overall. `Smooth before -> after` compares the same
candidate on both sides and therefore must not be confused with `Best hard`.

| Acquisition | Best hard post-ECC /24 | Best measured smooth same-candidate /24 | Smooth fit RMS | LOO RMS | Smooth resamples | Smooth HMAC slots | HMAC |
|---|---:|---:|---:|---:|---:|---:|---|
| scanner `0270_001.jpg` | 4 | 4 -> 9 | 1.362 | 1.969 | 2 | 0 | no |
| scanner `0270_002.jpg` | 5 | 8 -> 6 | 1.536 | 2.435 | 2 | 0 | no |
| `foto bici dritta.jpg` | 6 | 7 -> 3 | 1.500 | 2.461 | 2 | 1 | no |
| `foto bici storta.jpg` | 4 | 3 -> 8 | 1.555 | 2.432 | 2 | 0 | no |
| `foto stampa storta.jpg` | 7 | 8 -> 6 | 1.432 | 2.211 | 2 | 0 | no |
| `foto stampa.jpg` | 2 | 7 -> 5 | 1.525 | 2.389 | 2 | 0 | no |

The diagnostic field improves at least one measured candidate in four of six cases
and worsens the selected measured candidate in two. The strongest raw numerical
improvement is the bicycle frontal photograph, 7 -> 3 known post-ECC header errors,
but that field fails the strict leave-one-out decode gate and is not used for HMAC.

A different bicycle-frontal candidate (`fundamental-sync-subpixel`, mild-highpass)
has fit RMS about 1.215 blocks and leave-one-out RMS about 1.910 blocks. It passes
the strict gate, improves the same candidate from 6 -> 5 known post-ECC errors and
therefore replaces that candidate's ordinary grid in one of the existing four HMAC
slots. Authentication still fails.

The real-corpus conclusion is deliberately narrower than “smooth correction works”.
A continuous component is present strongly enough to improve selected candidates,
but one global affine field is not sufficiently predictive across the corpus. The
best-hard frontal photograph remains at 2/24 known post-ECC errors because its best
candidate contains only four complete spatial cells and therefore cannot support a
six-control field. Scanner cases likewise show that a very strong macro lattice does
not guarantee that one affine phase field describes the protected-bit channel.

Runtime remains bounded. On the development environment the two scanner cases
complete in about 10.2--10.3 s, the smartphone cases in about 12.7--17.4 s, and no
case exceeds the two-smooth-resample or four-HMAC-slot ceilings. No real case
produces a valid Format-v3 HMAC.

## v0.3.0-build11 continuous-control checkpoint — 2026-09-10

Build11 was evaluated on the same six private print-acquisition cases as builds 8--10
with the physical-corpus test key `Piccotti`. The images remain private and are not distributed. The table uses
the final boundary-aware confidence implementation. `Best hard` is the existing hard
bit-channel summary. `Best measured smooth` compares the same candidate before and
after a build11 corrected-grid resample; `--` means no candidate passed the diagnostic
resample gate. A resampled field still receives no HMAC unless it also passes the
strict decode gate and improves the known prefix without increasing known coded-bit
errors.

| Acquisition | Best hard post-ECC /24 | Best measured smooth same-candidate /24 | Mean control confidence (measured field) | Fit RMS | LOO RMS | Smooth resamples | Smooth HMAC slots | HMAC |
|---|---:|---:|---:|---:|---:|---:|---:|---|
| `foto stampa.jpg` | 2 | -- | -- | -- | -- | 0 | 0 | no |
| `foto stampa storta.jpg` | 7 | 7 -> 7 | 0.346 | 1.477 | 2.223 | 2 | 0 | no |
| `foto bici dritta.jpg` | 6 | 6 -> 8 | 0.475 | 1.067 | 1.766 | 1 | 0 | no |
| `foto bici storta.jpg` | 4 | **4 -> 3** | 0.354 | 1.175 | 2.182 | 2 | 0 | no |
| scanner `0270_001.jpg` | 4 | 4 -> 6 | 0.330 | 1.036 | 1.698 | 1 | 0 | no |
| scanner `0270_002.jpg` | 5 | 7 -> 7 | 0.314 | 1.293 | 2.103 | 2 | 0 | no |

The build11 control estimator materially changes the interpretation of build10. The
single build10 real field that entered an HMAC slot (`foto bici dritta`, 6 -> 5) no
longer improves when phase maxima at the edge of the bounded +/-2 search are explicitly
down-weighted. In final build11 the corresponding confidence-aware field is still
geometrically decode-eligible, but changes the known prefix from 6 -> 8; the existing
improvement requirement therefore prevents it from replacing the hard grid. No real
case consumes a smooth-field HMAC slot.

`foto bici storta.jpg` contains the strongest remaining same-candidate diagnostic
change: 4 -> 3 post-ECC known-header errors. Its robust affine fit RMS is about 1.175
blocks, but leave-one-out RMS is about 2.182 blocks, above the 2.0 decode limit. It is
therefore measured but not promoted. This is evidence that useful local correction
exists without evidence that the fitted global field predicts unseen cells reliably.

The regularized quadratic model is computed only with at least eight controls and is
selected only after an explicit leave-one-out gain. Across the six real acquisitions it
wins model selection on one `0270_002` candidate. That model has fit RMS about 0.775
blocks and leave-one-out RMS about 1.915 blocks, but predicts a maximum correction of
about 2.73 blocks, above the 2.5-block bound. It is rejected before resampling. The
corpus therefore does not currently justify promoting the quadratic model.

The front smartphone capture remains the closest hard channel case at 2/24 known
post-ECC errors. Its best hard candidate has only four complete spatial cells; its
leave-one-out behavior is unstable and no build11 smooth field is resampled. This
preserves the good hard candidate rather than forcing a higher-order correction.

Approximate diagnostic wall times in the development environment remain bounded:
about 9.2--11.3 s for the scanner cases and 11.5--17.4 s for the four smartphone cases.
No image authenticates. The negative result is intentional: build11 improves control
quality and rejects the only build10 real HMAC-field substitution instead of claiming
a fragile gain.


## v0.3.0-build12 blind-registration checkpoint — 2026-09-11

Build12 was evaluated on the same six private acquisitions with the physical-corpus test key `Piccotti`; the
key is used only by the ordinary bit/HMAC diagnostic after the blind observer has
constructed its controls. `Local oracle` is the build11 key-assisted local-phase
majority result and is shown only for comparison. `Blind score/confidence/oracle
distance` come from the best hard candidate's independent blind observation. `Blind
smooth` compares the same candidate before/after a corrected-grid resample; `--` means
no blind field passed the diagnostic resample gate.

| Acquisition | Best hard /24 | Local oracle /24 | Blind structural profile | Blind global score | Mean blind conf. | Mean blind-oracle distance (blocks) | Best blind smooth same-candidate /24 | Smooth HMAC slots | HMAC |
|---|---:|---:|---|---:|---:|---:|---|---:|---|
| `foto stampa.jpg` | **2** | 10 | balanced | 0.193 | 0.080 | 1.625 | -- | 0 | no |
| `foto stampa storta.jpg` | 7 | **2** | balanced | **0.410** | **0.138** | 2.393 | **5 -> 4** | 0 | no |
| `foto bici dritta.jpg` | 6 | 4 | balanced | 0.209 | 0.092 | 2.528 | -- | 0 | no |
| `foto bici storta.jpg` | 4 | 3 | balanced | 0.121 | 0.073 | 1.992 | -- | 0 | no |
| scanner `0270_001.jpg` | 4 | 6 | balanced | 0.189 | 0.067 | 2.043 | -- | 0 | no |
| scanner `0270_002.jpg` | 5 | **3** | robust | 0.107 | 0.035 | 2.587 | -- | 0 | no |

Only the inclined historical smartphone acquisition produces blind fields that pass
the diagnostic resample gate. Its best measured blind field uses nine controls and
`control_source=blind-self-registration`; on a mild-highpass capacity candidate it
changes the same-candidate known post-ECC prefix from 5 to 4 errors. The robust affine
fit RMS is about `1.271` blocks and leave-one-out RMS about `2.041` blocks. The latter
remains above the strict `2.0` decode threshold, so the correction is not promoted to
an HMAC slot. A second field on the same image is also diagnostic-eligible but not
decode-eligible.

The blind-vs-oracle distances remain large (roughly 1.6--2.6 blocks on the reported
best-hard candidates), showing that the two observers do not yet identify the same
local phase consistently. This is particularly important because the key-assisted
oracle can reach 2/24 on `foto stampa storta.jpg` and 3/24 on scanner `0270_002`, while
the independent blind controls cannot yet reproduce those gains. The direct pairwise
cross-cell observer was also measured during development; real pair peaks were around
0.09 and often several blocks from the oracle, so intra-tile repetition became the
primary build12 observer.

A wrong-key run on `foto stampa storta.jpg` remains `projective-failed`, performs the
normal bounded four full-decode attempts and authenticates nothing. Blind structural
scores can still be nonzero under a wrong key because the observer itself is
key-independent; that score is therefore explicitly not a watermark-detection result.
No real acquisition in build12 produces a valid Format-v3 HMAC.


## v0.3.0-build13 guided dual-observer checkpoint — 2026-09-11

Build13 was evaluated on the same six private physical acquisitions. The guided
cross-cell observer uses no key or expected header values. `Observer distance` compares
its zero-mean controls with the primary repetition controls; `blind-oracle distance`
remains the ex-post comparison against the build11 key-assisted local-phase oracle.

| Acquisition | Best hard /24 | Repetition score | Guided pair score | Mean observer distance (blocks) | Mean blind-oracle distance (blocks) | Consensus / cycle fixes | Best blind smooth | HMAC |
|---|---:|---:|---:|---:|---:|---|---|---|
| `foto stampa.jpg` | **2** | 0.193 | 0.065 | 0.853 | 1.625 | 0 / 0 | -- | no |
| `foto stampa storta.jpg` | 7 | **0.410** | 0.062 | 0.587 | 2.393 | 0 / 0 | **5 -> 4** | no |
| `foto bici dritta.jpg` | 6 | 0.209 | 0.058 | 0.744 | 2.528 | 0 / 0 | -- | no |
| `foto bici storta.jpg` | 4 | 0.121 | 0.053 | **0.515** | 1.992 | 0 / 0 | -- | no |
| scanner `0270_001.jpg` | 4 | 0.189 | 0.060 | 0.652 | 2.043 | 0 / 0 | -- | no |
| scanner `0270_002.jpg` | 5 | 0.107 | 0.055 | 0.664 | 2.587 | 0 / 0 | -- | no |

The guided search is substantially better behaved than the unrestricted build12
pairwise control: its offsets stay within roughly one block of the primary rather than
wandering several blocks. However, absolute pair correlation remains weak. The
controlled synthetic consensus regression gives about 0.61 mean pair score, whereas
all six real best-hard observations remain below 0.07. The independent 0.10 strength
gate therefore keeps every real secondary observation in check-only mode. No real
control is cycle-slip adjusted.

Because weak secondary evidence is forbidden from vetoing the primary observer, the
build12 `foto stampa storta.jpg` result is preserved: two blind fields reach the
diagnostic resample stage and the best same-candidate result is 5 -> 4 post-ECC known
header errors. Its leave-one-out error remains above the strict decode gate, so it
receives no HMAC slot. The frontal image remains hard-best at 2/24 and is untouched.
No acquisition authenticates in build13.

## v0.3.0-build14 local lattice-phase checkpoint — 2026-09-11

Build14 was evaluated on the same four private smartphone captures and two private
scanner acquisitions with the physical-corpus test key `Piccotti`. Originals remain outside source archives.
`Best smooth` is always a same-candidate before/after comparison; it must not be
compared directly with `Best hard` when the base value differs.

| Acquisition | Best hard /24 | Oracle local /24 | Lattice mean conf. | Mean fractional disagreement (blocks) | Lattice consensus / cycle fixes | Best smooth same-candidate /24 | Smooth mean conf. | Fit / LOO RMS | Smooth HMAC | HMAC |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `foto stampa.jpg` | **2** | 10 | 0.419 | 0.314 | 3 / 0 | -- | -- | -- | 0 | no |
| `foto stampa storta.jpg` | 7 | **2** | 0.384 | 0.367 | 6 / 1 | **5 -> 2** | 0.174 | 1.240 / 1.952 | 0 | no |
| `foto bici dritta.jpg` | 6 | 4 | 0.407 | 0.398 | 5 / 2 | **7 -> 3** | 0.145 | 1.387 / 2.032 | 0 | no |
| `foto bici storta.jpg` | 4 | 3 | 0.438 | 0.447 | 3 / 3 | 9 -> 8 | 0.196 | 0.727 / 1.279 | 0 | no |
| scanner `0270_001.jpg` | **4** | 6 | 0.404 | 0.375 | 7 / 3 | **9 -> 4** | 0.125 | 1.168 / 2.262 | 0 | no |
| scanner `0270_002.jpg` | 5 | 3 | **0.522** | 0.372 | 6 / 2 | 7 -> 10 | 0.201 | 1.557 / 2.121 | 0 | no |

The lattice phase is substantially stronger than build13 guided cross-cell correlation,
but high geometric confidence is still not equivalent to a correct protected-bit grid.
The storta 5 -> 2 result is especially informative: its fit and LOO both pass the old
numeric geometry limits, yet mean fitted-control confidence is only about 0.174, below
the unchanged 0.20 decode gate. It therefore receives no HMAC slot. The straight-bike
7 -> 3 field misses both confidence and (slightly) the 2.0 LOO limit. Scanner 001 reaches
9 -> 4 but starts from a candidate worse than the global hard 4/24 and also fails the
confidence/LOO promotion gates.

Two explicitly recorded ablations prevent over-interpreting these gains. With the
opposite residual/correction sign interpretation the storta candidate remains 5 -> 5.
With the final observed-drift sign but cycle unwrap disabled, storta becomes 5 -> 6 and
straight-bike 7 -> 11. The strong final improvements therefore depend on the bounded
integer unwrap, which is also the least validated part of the method. Scanner 002's
7 -> 10 degradation confirms that the current unwrap is not ready for promotion.

Wrong-key control on `foto stampa storta.jpg` remains `projective-failed`, executes four
ordinary full-decode attempts, zero smooth decode attempts, and authenticates nothing.
Lattice phase remains measurable under a wrong key, as expected for key-independent
geometry evidence.



## v0.3.0-build15 global discrete unwrap checkpoint — 2026-09-11

Build15 changes the question from "can a local +/-1 unwrap improve the known prefix?"
to "is the integer cycle field uniquely selectable from key-independent geometry?".
Representative final diagnostics answer the second question conservatively: multiple
real cases admit much better geometric objectives but insufficient first/second margins.
Those proposals are rolled back transactionally and receive no HMAC slot.

For `foto stampa storta.jpg`, the global solver evaluates 1584 states, proposes changes
in nine cells and reduces its objective from about 5.480 to 0.804, but the alternative
solution is about 0.826 (margin about 0.021), so the unwrap is ambiguous. Zero changes
are committed. The safe pre-lattice blind field still measures 5 -> 4 on the same
candidate and remains above the strict smooth LOO gate.

This checkpoint intentionally does not reinterpret the build14 5 -> 2 / 7 -> 3 / 9 -> 4
results as successes. They remain valuable evidence that useful cycle assignments exist,
but build15 shows they are not yet independently identifiable with sufficient uniqueness.
Only a valid Format-v3 HMAC is success.

## v0.3.0-build16 exact top-2 checkpoint — 2026-09-11

Build16 first corrects the uniqueness certificate itself. Build15's width-64 beam could
prune a partial assignment that later became the true runner-up, so its first/second
margin was not guaranteed to be global. The new solver enumerates the complete bounded
`{-1,0,+1}` state space independently for X and Y (at most `3^9 = 19683` assignments
per axis), preserves the existing key-independent objective and acceptance thresholds,
and obtains the exact top-1/top-2. A deterministic regression captures an actual
false-accept pattern from the old beam: the optimum remains `0.027408951002`, while the
runner-up changes from the beam's `0.062494434795` to the exact `0.035129429298`, which
correctly makes the field ambiguous under the unchanged gate.

The no-eligible, insufficient-coverage and partial-lattice paths are now explicit
`not-applicable` outcomes rather than malformed/ambiguous solutions, and public numeric
diagnostics remain finite/JSON-serializable. Per-axis status and eligibility are
reported separately. These changes do not consume additional HMAC attempts or smooth
resamples.

Build16 also adds `split-repetition-top2`, a key-independent consistency diagnostic that
compares the exact geometric top-1/top-2 with two deterministic disjoint partitions of
the existing repeated-position evidence. It cannot override ambiguity: the primary
repetition controls were estimated from the complete pair set, so this is useful
structural telemetry but not a fully held-out certificate.

Only four of the six canonical physical acquisitions were available in this build16
session. The two bicycle photographs were absent. The two historical smartphone files
are now canonically named `.jpg`, matching their JPEG encoding; only their filenames
changed.

| Acquisition | Best hard /24 | Exact states | Eligible X/Y | Baseline -> best | Exact second | Margin | Split fold deltas | Split supports best | Same-candidate smooth | Smooth HMAC | HMAC |
|---|---:|---:|---:|---:|---:|---:|---|---|---|---:|---|
| `foto stampa.jpg` | **2** | 0 | 0/0 | n/a | n/a | n/a | n/a | n/a | -- | 0 | no |
| `foto stampa storta.jpg` | 7 | 8748 | 7/8 | 5.4804 -> 0.8045 | 0.8256 | 0.02113 | +0.1002 / -0.0361 | no | **5 -> 4** | 0 | no |
| scanner `0270_001.jpg` | **4** | 26244 | 9/8 | 5.7797 -> 0.8083 | 0.8297 | 0.02144 | -0.0828 / -0.0343 | no | -- | 0 | no |
| scanner `0270_002.jpg` | 5 | **39366** | 9/9 | 7.3178 -> 1.7229 | 1.7407 | **0.01780** | +0.0410 / -0.0929 | no | -- | 0 | no |

`foto stampa storta.jpg` therefore preserves the conservative build15 result: a large
geometric objective improvement exists, but the exact competitor is too close and the
two split-repetition folds disagree. Scanner `0270_001.jpg` is even stronger negative
evidence for the selected geometric top-1 because both split folds prefer its runner-up.
Scanner `0270_002.jpg`, the historical counterexample, tightens from the build15 beam
margin of roughly `0.02634` to an exact margin of `0.01780`; its split folds also
disagree. No proposal is promoted and no smooth HMAC slot is consumed.

No available real acquisition authenticates in build16. The next useful experiment must
construct proposal and validation evidence that are disjoint by design and then be
frozen before the complete six-image corpus is evaluated. A lower known-header error
count remains research evidence only; success still requires an ordinary valid Format-v3
HMAC.

## v0.3.0-build17 held-out cross-fit checkpoint — 2026-09-11

Build17 preserves the exact all-pairs build16 results and adds two genuinely disjoint
repetition proposal/validation directions. The folds are assigned by coded-bit group,
so the proposal fold and held-out fold do not reuse a repetition group or logical
repetition position.

Only the four physical acquisitions available in the current private corpus were
measured; the two bicycle photographs remain absent.

| Acquisition | Best hard /24 | All-pairs unwrap | A->B delta | B->A delta | A/B support | Local-cycle agreement | Cross-fit support | HMAC |
|---|---:|---|---:|---:|---|---:|---|---|
| `foto stampa.jpg` | **2** | not-applicable | n/a | n/a | n/a | n/a | not run | no |
| `foto stampa storta.jpg` | 7 | ambiguous (0.02113) | +0.00465 | +0.20522 | yes / yes | 2/9 (0.222) | no | no |
| scanner `0270_001.jpg` | **4** | ambiguous (0.02144) | -0.18283 | -0.02287 | no / no | 0/9 | no | no |
| scanner `0270_002.jpg` | 5 | ambiguous (0.01780) | +0.05639 | -0.24050 | yes / no | 0/9 | no | no |

The inclined smartphone result shows why both parts of the build17 criterion matter:
held-out validation is positive in both directions, yet the independently estimated
cycle fields mostly disagree. Scanner 001 provides the complementary result: both
held-out folds prefer the runner-up rather than the proposal's top-1. Scanner 002 remains
the historical negative control and is directionally inconsistent.

No cross-fit result overrides the all-pairs ambiguity gate, no smooth HMAC attempt is
added, and no physical acquisition authenticates. The result is therefore scientific
progress in failure characterization rather than decoder success.



## v0.3.0-build18 lazy cross-fit / instability checkpoint — 2026-09-12

Build18 preserves build17's all-pairs and held-out decisions while running cross-fit only
on the final best bit candidate. The four available physical cases were re-evaluated with
the same public test key and no HMAC success:

| Acquisition | Best hard /24 | all-pairs | A->B delta | B->A delta | A/B cycles | cross-fit attempts/skipped | selected cross-fit ms |
|---|---:|---|---:|---:|---:|---:|---:|
| `foto stampa.jpg` | **2** | not-applicable | n/a | n/a | n/a | 0 / 4 | 0 |
| `foto stampa storta.jpg` | 7 | ambiguous (0.02113) | +0.00465 | +0.20522 | 2/9 | 1 / 3 | 26 |
| scanner `0270_001.jpg` | **4** | ambiguous (0.02144) | -0.18283 | -0.02287 | 0/9 | 1 / 3 | 48 |
| scanner `0270_002.jpg` | 5 | ambiguous (0.01780) | +0.05639 | -0.24050 | 0/9 | 1 / 3 | 73 |

For `foto stampa storta.jpg`, the mean of `min(conf_A, conf_B)` is `0.1337` on the two
agreeing cells and `0.3576` on the seven disagreements. Mean lattice confidence is
`0.4059` versus `0.3772`, respectively. Thus the agreement subset is not distinguished
by stronger independent repetition evidence; a partial-consensus decoder change is not
supported. Scanner 001/002 provide no agreeing cells at all.

These timings are measurements from the build18 development environment, not a direct
speed comparison with another host. Their methodological value is that the selected
cross-fit itself is small and explicitly counted, while three lower-ranked candidates
are skipped. Production decode/HMAC budgets remain unchanged.


## v0.3.0-build19 multi-partition stability checkpoint — 2026-09-12

Build19 re-evaluates only the final best ambiguous bit candidate using 8 fixed
coded-bit-group partitions in both directions (16 held-out trials maximum). Production
decode and HMAC budgets are unchanged.

| Acquisition | Best hard /24 | trials available | supported | unique fields | supported unique | mean cell modal | min cell modal | pairwise agreement | stability ms |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `foto stampa.jpg` | **2** | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `foto stampa storta.jpg` | 7 | 16 | 9 | **16** | **9** | 0.2153 | 0.1250 | 0.0731 | 439 |
| scanner `0270_001.jpg` | **4** | 16 | 6 | **16** | **6** | 0.2153 | 0.1250 | 0.0759 | 501 |
| scanner `0270_002.jpg` | 5 | 16 | 12 | **16** | **12** | 0.2361 | 0.1875 | 0.0824 | 1048 |

The full-field modal fraction is therefore exactly `1/16 = 0.0625` in every ambiguous
case. Restricting to positively held-out-supported trials does not improve reproducibility:
all supported complete fields remain unique. No cell is unanimous. The frontal case stays
`not-applicable` and pays no stability cost.

This is negative but decisive evidence: on the current physical corpus repetition
constraints provide useful local scoring, but not a persistent integer-cycle field.
Further repartition/voting is not justified without a new independent signal.

## v0.3.0-build20 independent cycle-anchor checkpoint — 2026-09-12

| acquisition | anchor conf | top-1 obj | second obj | delta second-top1 | top-1 agreement | second agreement |
|---|---:|---:|---:|---:|---:|---:|
| foto stampa storta | 0.0575 | 5.9922 | 6.5005 | +0.5083 | 0/9 | 0/9 |
| foto stampa | n/a | n/a | n/a | n/a | n/a | n/a |
| 0270_001 | 0.0168 | 8.7807 | 9.0084 | +0.2277 | 0/9 | 0/9 |
| 0270_002 | 0.0798 | 14.6072 | 14.7028 | +0.0956 | 0/9 | 0/9 |

All three ambiguous cases produce a lower continuous anchor objective for top-1, but none
shares an integer-cycle cell with either top-1 or runner-up after gauge alignment. Scanner001
contradicts the held-out repetition preference and scanner002 remains a negative control. The
anchor is therefore not promoted. No physical payload authenticates.


## v0.3.0-build21 Format-v3 observability checkpoint — 2026-09-12

This checkpoint is format-static; it does not require a physical acquisition. Every one-block
neighbor of the correct tile origin was audited under public repetition topology and Hamming
validity.

| profile | coded bits | repetition pairs | nearest pair gap | vertical pair gap | min synthetic syndrome | notable alias |
|---|---:|---:|---:|---:|---:|---|
| robust | 448 | 896 | 0.0580 | 0.0781 | **0.0000** | vertical ±1 zero-syndrome synthetic alias |
| balanced | 672 | 448 | 0.0558 | 0.0781 | 0.0882 | vertical ±1 only weakly separated |
| capacity | 1120 | 0 | 0 | n/a | **0.0000** | vertical ±1 exact Hamming-word permutation |

For robust/balanced, a horizontal one-block shift preserves 94.20%/94.42% of the public
repetition-pair graph; a vertical shift preserves 92.19%. Thus the repetition mapping is not
mathematically invariant, but the absolute-origin signal is only a small residual difference
that must survive print-camera noise and local geometry error.

Hamming parity supplies a second key-independent observable because whitening occurs before ECC.
Horizontal and diagonal unit shifts produce roughly 0.75--0.79 non-zero-syndrome fraction in the
deterministic valid-frame probe. The vertical direction is structurally problematic: capacity
maps complete Hamming words to complete Hamming words with the same bit roles, so syndrome cannot
distinguish the shifted origin at all; robust also yields zero syndrome in the deterministic
probe, while balanced rises only to about 0.088--0.093.

The combined audited mechanisms therefore fail the `audited_unit_shifts_universally_separated`
condition. This does not prove that no conceivable image-domain statistic can recover an origin,
but it proves that the two strongest public structural mechanisms already present in Format v3
are not a uniform absolute-cycle reference. No production behavior changes.


## v0.3.0-build22 physical-topology checkpoint — 2026-09-12

The build22 probe uses held-out topology-exclusive repetition edges after shared-edge-only phase
registration. The frontal smartphone image remains exact-unwrap `not-applicable`, so no build22
physical-topology attempt is made there.

| case | topology ms* | profiles | shifts | mean modal winner phase | min modal | mean cell direction consistency | min direction consistency | mean abs cell delta | HMAC |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| foto stampa storta | 239 | 2 | 16 | **0.3125** | 0.2500 | 0.6944 | 0.5556 | 0.2731 | fail |
| foto stampa | 0 | 0 | 0 | n/a | n/a | n/a | n/a | n/a | fail |
| scanner 001 | 237 | 2 | 16 | **0.4375** | 0.3750 | 0.7292 | 0.5556 | 0.2514 | fail |
| scanner 002 | 254 | 2 | 16 | **0.4375** | 0.3750 | 0.6597 | 0.5556 | 0.2125 | fail |

\* development-container timings; qualification-host absolute timing must be measured separately.

The central falsification is not that the held-out edge deltas are zero—they are not. It is that
phase winner stability and cell consistency overlap strongly between the useful acquisition and
the negative control. The residual topology asymmetry therefore does not provide a safe blind gate.

### Format-v4 sizing result

| candidate | tile | pilot | data | comparable max payload | area ratio vs v3 | ideal pilot random-wrong-phase z |
|---|---:|---:|---:|---:|---:|---:|
| compact | 35×32 | 64 | 1056 | 59 B | 1.0000 | 8.00 |
| preserve-capacity | **37×32** | **64** | **1120** | **64 B** | **1.0571** | **8.00** |
| strong-pilot | 38×32 | 96 | 1120 | 64 B | 1.0857 | 9.80 |

The recommended prototype is 37×32/64 pilot because it adds a dedicated absolute reference without
reducing current data positions. These figures are sizing evidence, not a Format-v4 specification
or physical guarantee.

## v0.3.0-build23 transition checkpoint — 2026-09-12

Build23 is based on the user-host build22 qualification. That run completed all 30 targets with release
baseline and qualification corpus PASS; only the historical experimental `geometry-test` and `affine-test`
reported two failures each. `research-unit`, `physical-topology-test`, `v4-design-study-test`, composition,
lattice, perspective and core portability all passed.

The final build22 physical topology measurements were:

| acquisition | modal phase mean | cell consistency mean | mean absolute cell delta |
|---|---:|---:|---:|
| inclined smartphone | 0.3125 | 0.6944 | 0.2731 |
| scanner 001 | 0.4375 | 0.7292 | 0.2514 |
| scanner 002 negative | 0.4375 | 0.6597 | 0.2125 |

Because the negative control is not separated, build23 records v3 absolute-cycle research as closed rather
than tuning another gate. No v3 physical HMAC has been recovered.

Build23 also records and verifies the frozen v3 core SHA-256 baseline. Its new deterministic v4 result is the isolated prototype foundation. `prototype-1-stratified-p64`
uses 64/1184 pilot positions, leaves 1120 data positions, is exactly sign-balanced and quadrant-balanced,
and has no perfect non-zero cyclic alias. Exhaustive wrong-shift maxima are 8 overlapping pilot positions
and absolute signed correlation 5. These metrics are regression anchors for prototype-1 only, not a final
v4 specification or physical guarantee.

## v0.3.0-build24 initial v4 pilot qualification — 2026-09-13

### Structural comparison

| candidate | max cyclic overlap | max wrong correlation | exact aliases | SHA-256 |
|---|---:|---:|---:|---|
| prototype-1-stratified-p64 | 8/64 | 5/64 | 0 | `99042fd9827e0606a92afaf088bd96009bdf1fc6040c2ae95905534a3d57ff4f` |
| **prototype-2-search-p64** | **7/64** | **4/64** | **0** | `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174` |

Build24 uses fixed seeds and 200,000 attempts in both the joint mask/sign stage and sign-refinement stage. The
result is reproducible by `make v4-pilot-search-test`.

### Partial visibility

| visible pilots | worst contiguous margin | median contiguous margin | worst random margin | median random margin | random false origins |
|---:|---:|---:|---:|---:|---:|
| 64 | 60 | 60 | 60 | 60 | 0/256 |
| 48 | 43 | 43.5 | 43 | 44 | 0/256 |
| 32 | 28 | 28 | 27 | 28 | 0/256 |
| 24 | 20 | 20 | 19 | 20 | 0/256 |
| 16 | 12 | 13 | 11 | 13 | 0/256 |

These are structural pilot-only margins, not authentication confidence.

### Initial local image-channel qualification

The synthetic v4 qualification carrier uses prototype-2 plus deterministic pseudo-random data-plane signs; it
does not contain a payload/HMAC. Both supplied local originals selected the correct pilot origin in native,
JPEG-82, blur-1, noise-4, gamma-1.15, exact resize-75, resize-50 and aligned-crop cases.

| local source | unmarked best score | unmarked margin | lowest marked score in matrix | lowest marked margin in matrix | origin errors |
|---|---:|---:|---:|---:|---:|
| PJ_lingua.PNG | 0.306699 | 0.022743 | 0.985886 | 0.552488 | 0 |
| PJ_piccolo.png | 0.232802 | 0.003697 | 0.790946 | 0.381056 | 0 |

No threshold is promoted from this two-image development corpus. Rotation/affine/perspective/combined tests and
real v4 print/scan evidence remain open, so prototype-2 is not yet a normative Format-v4 pilot.



## v0.3.0-build25 v4 known-geometry qualification — 2026-09-13

`prototype-2-search-p64` recovered cyclic origin `(0,0)` in all 16 deterministic synthetic geometry cases and all 32 cases across the two local original images. The synthetic matrix minimum positive runner-up margin was 0.531955; the maximum matching unmarked-control margin was 0.020642. On the local originals, the weakest positive margin was 0.322069 (`PJ_piccolo`, combined perspective+blur), the maximum negative-control margin was 0.067786, and the weakest positive-minus-negative separation was 0.258085.

The Build25 regression floor is 0.10 positive-minus-negative margin separation. It is a development gate only and must not be reused as a production acceptance threshold without independent calibration. No Format-v3 frozen file changed and no v4 payload encoder/decoder exists yet.

## v0.3.0-build26 bounded blind v4 geometry qualification — 2026-09-13

Build26 introduces a pilot-assisted geometry search with three separated evidence stages: repeated data-plane self-consistency for proposal, central-row pilot ranking, and held-out corner-pilot validation. The coarse repeat score uses only non-pilot coordinates and is invariant to the pilot sign sequence.

Representative synthetic recoveries include:

| case | recovered geometry | origin | corner error | pilot margin |
|---|---|---:|---:|---:|
| rotate 7.4 deg | ~7.45 deg | (0,0) | 0.00058 | 0.577991 |
| rotate+scale | 11.2 deg, 1.07/0.93 | (0,0) | 0.00000 | 0.595496 |
| perspective | ~0 deg, 0.035/0.010 | (0,0) | 0.00075 | 0.592342 |
| combined perspective | 9.35 deg, 1.08/0.92, 0.030/0.015 | (0,0) | 0.00057 | 0.559796 |
| combined shear | -8.75 deg, shearX 6.0 | (0,0) | 0.00055 | 0.567819 |

The local corpus gate uses two representative compound cases per source. `PJ_piccolo.png` recovers 11.2 deg + 1.07/0.93 exactly with margin 0.520953, and the combined 9.3 deg + 1.08/0.92 + 0.030/0.015 perspective case exactly with margin 0.425713; corresponding unmarked blind-search margins are 0.094881 and 0.083961. `PJ_lingua.PNG` also recovers both declared geometries and origin `(0,0)`. These are development observations, not production thresholds.

The current search assumes known canonical extent and auto-framed transforms and is bounded to <=6400 repeat hypotheses in the synthetic gate. No v3 frozen-core file changes and no v4 payload/HMAC path exists.



## v0.3.0-build27 unknown-placement qualification — 2026-09-13

Build27 keeps geometry supplied independently and removes only the placement assumption. The deterministic synthetic placement suite currently reports:

| case | recovered shift | marked validation | matching negative | notes |
|---|---:|---:|---:|---|
| rotate+scale arbitrary crop | (-113,-77) | 0.994367 | 0.100612 | exact non-block-aligned crop |
| combined perspective arbitrary crop | (-101,-69) | 0.980152 | 0.192472 | exact crop under projective warp |
| rotate+scale deep crop | (-171,-119) | 1.000000 | 0.094531 | substantially reduced visible carrier |
| rotate+scale padded placement | (73,40) | 0.986240 | 0.221759 | one-pixel-equivalent padded solution |
| perspective padded placement | (57,35) | 0.990739 | 0.176874 | exact padded placement |

The padded rotate+scale result can differ by one pixel from the synthetic insertion while still producing an essentially identical validated carrier mapping; complete pilot validation remains the deciding diagnostic. Placement top-2 gaps are intentionally not acceptance criteria because whole-tile-equivalent placements are expected for a periodic carrier.

The local original-image corpus reports:

| image | case | recovered shift | marked validation | negative validation |
|---|---|---:|---:|---:|
| PJ_lingua.PNG | rotate+scale crop | (-117,-83) | 0.981501 | 0.301624 |
| PJ_lingua.PNG | combined perspective padded | (67,43) | 0.982102 | 0.324639 |
| PJ_piccolo.png | rotate+scale crop | (-117,-83) | 0.932792 | 0.231024 |
| PJ_piccolo.png | combined perspective padded | (67,43) | 0.953116 | 0.266898 |

A wrong-geometry control (true affine geometry plus +1 degree error) remains available to the placement optimizer but reaches only 0.236030 validation versus 0.994367 for the correct geometry. This supports the staged conclusion: the pilot can recover unknown placement once geometry is sufficiently correct, but Build27 does not yet solve geometry and placement jointly.


## v0.3.0-build28 joint blind affine+crop qualification — 2026-09-13

Build28 combines unknown rotation/anisotropic scale with unknown negative crop translation. Geometry is selected only from sign-independent DCT phase contrast; the pilot enters after geometry selection through Build27 placement proposal/validation.

Synthetic results:

| case | recovered affine | recovered shift | validation | full pilot margin | matched-negative validation | corner error | hypotheses |
|---|---|---:|---:|---:|---:|---:|---:|
| rotate+scale+crop | ~11.25 deg, 1.070/0.9325 | (-113,-78) | 0.975547 | 0.585616 | -0.120030 | 0.00102 | 69574 |
| negative rotate+scale+crop | ~-7.30 deg, 1.030/0.9675 | (-90,-63) | 0.881391 | 0.507172 | 0.005611 | 0.00135 | 68989 |

Local development originals:

| image | recovered affine | recovered shift | validation | full pilot margin | matched-negative validation | corner error | hypotheses |
|---|---|---:|---:|---:|---:|---:|---:|
| PJ_lingua.PNG | ~11.20 deg, 1.070/0.935 | (-118,-85) | 0.812461 | 0.432900 | 0.032305 | 0.00189 | 69818 |
| PJ_piccolo.png | ~11.25 deg, 1.0675/0.930 | (-116,-83) | 0.781974 | 0.410485 | -0.136843 | 0.00102 | 69881 |

Negative partition scores may be below zero because the signed public-pilot correlation is not clamped; they are not probabilities. All positive cases recover cyclic origin `(0,0)`. The regression floors remain 0.50 held-out validation for the local corpus, 0.15 positive-minus-matched-negative validation separation, 0.15 local full-pilot margin, and 0.015 local corner-error ratio.

Build28 does not qualify joint projective/perspective or positive padded-canvas geometry, and it does not freeze the pilot or add a v4 payload/HMAC path.

## v0.3.0-build29 joint projective/padded qualification — 2026-09-14

| case | expected outcome | validation | full pilot margin | origin | corner error | hypotheses |
|---|---|---:|---:|---:|---:|---:|
| synthetic projective+crop | ACCEPT | 0.554332 | 0.269007 | (0,0) | 0.00472 | 427982 |
| synthetic affine+padded | ACCEPT | 0.986240 | 0.542529 | (0,0) | 0.00085 | 299613 |
| PJ_lingua projective+crop | ACCEPT | 0.383955 | 0.179718 | (0,0) | 0.00878 | 429470 |
| PJ_piccolo projective+crop | SAFE REJECT | 0.421398 | 0.118781 | (0,0) | not accepted | 426641 |
| PJ_lingua affine+padded joint search | SAFE REJECT | 0.190152 | 0.029623 | (0,0) | not accepted | 302683 |
| PJ_piccolo affine+padded joint search | SAFE REJECT | 0.373953 | 0.029185 | (20,22) | not accepted | 299134 |

The padded `PJ_lingua` known-geometry control reaches validation 0.991263, score 0.988352, margin 0.553797 and origin `(0,0)`. This distinguishes the current joint-search limitation from loss of pilot signal. Projective and padded Build29 thresholds are development regression gates only; SAFE REJECT counts as the correct result for cases outside the currently qualified joint-search envelope.



## v0.3.0-build30 pilot candidate-lock audit — 2026-09-14

Build30 does not change the pilot and does not widen Build29 ACCEPT thresholds. It adds a semantic identity lock plus a known-mapping corpus audit designed to answer whether current SAFE REJECT cases indicate a weak pilot or a weak search.

Locked identity: `prototype-2-search-p64` / `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`.

| carrier class | known mapping | marked score | marked margin | negative margin | origin |
|---|---|---:|---:|---:|---:|
| large local original | projective + crop | 0.986195 | 0.559445 | 0.007237 | (0,0) |
| large local original | affine + padded | 0.988352 | 0.553797 | 0.000598 | (0,0) |
| smaller local original | projective + crop | 0.928412 | 0.490432 | 0.010553 | (0,0) |
| smaller local original | affine + padded | 0.932220 | 0.484086 | 0.010305 | (0,0) |

Development gates for this audit are score >= 0.85, margin >= 0.40, marked-minus-negative margin separation >= 0.40, negative margin <= 0.08 and origin `(0,0)`. These are freeze-readiness regression floors, not decoder acceptance thresholds.

The smaller padded case is decisive: Build29's bounded placement/joint search SAFE REJECTS it, but the exact mapping exposes a strong pilot. The unresolved padded/projective cases therefore remain decoder search problems. A fully translation-invariant 64-phase DCT ranking experiment was also tested and rejected because natural photographic texture still dominated the structural score.

Build30 locks the candidate identity for continued development but does **not** promote it to normative Format v4. Physical v4 print-camera/scanner qualification is still absent.

## Build31 — first real v4 frame/encoder digital qualification

Build31 keeps the locked pilot unchanged and adds the first authenticated data plane. The unit qualification passes all three profiles, wrong-key rejection, v3/v4 parser separation, Hamming single-error correction, exact pilot/data partition, deterministic frame vector, minimum 296x256 capacity round-trip, JPEG-q82 robust recovery and block-aligned crop origin recovery.

Private development-original results with strength 24:

| image class | profile/channel | pilot score | pilot margin | authenticated |
|---|---|---:|---:|---|
| larger original | robust native | 1.000000 | 0.597052 | yes |
| larger original | robust JPEG q82 | 1.000000 | 0.597329 | yes |
| larger original | robust aligned crop | 1.000000 | 0.622984 | yes |
| larger original | balanced native | 1.000000 | 0.595112 | yes |
| larger original | capacity native | 1.000000 | 0.523919 | yes |
| smaller original | robust native | 0.993618 | 0.590285 | yes |
| smaller original | robust JPEG q82 | 0.990079 | 0.594082 | yes |
| smaller original | robust aligned crop | 0.991861 | 0.595633 | yes |
| smaller original | balanced native | 0.993618 | 0.604494 | yes |
| smaller original | capacity native | 0.993618 | 0.550018 | yes |

Every listed result includes a valid v4 HMAC; pilot score is telemetry, not the success criterion. The existing original-image corpus is still a **digital** development corpus. Physical qualification requires new paper/scanner/smartphone captures of carriers produced by Build31 `v4-embed`.

Build31 does not claim that Hamming(7,4) is the final v4 ECC and does not yet connect the blind Build29 geometry search to payload decoding.


## Build32 corpus reset

Active corpus: LQ/MQ/HQ only, manifest-verified. LQ deep baseline is 33/33 PASS at the revised 55% common resize floor. MQ robust/balanced completed without failure and MQ capacity independently passes all 11 baseline transformations. Build31 v4 frame tests on LQ/MQ pass native, JPEG q82 and aligned crop in all profiles; HQ remains a manifest-verified 200.5 MP physical-fixture source and is skipped by bounded Go corpus tests before decode.

## Build33 MQ joint-projective ranking observability

Build33 adds a controlled truth-basin diagnostic for the Build29 projective+crop transform on the active MQ carrier. The same source image and transform are tested with two different marked data planes: the historical deterministic random-data carrier and a real authenticated Format-v4 robust frame carrying `v4-b33-auth`.

On the Build33 development host, the first near-truth basin ranks as follows:

| Carrier | structural top-5000 | half-pilot prefilter | full proposal (after top-2000 prefilter retention) |
| --- | ---: | ---: | ---: |
| random data plane | 1 / 5000 | 98 / 5000 | 22 / 2000 |
| authenticated v4 frame | 1 / 5000 | 35 / 5000 | 2 / 2000 |

The structural stage therefore already contains a strong near-truth basin in both cases. The ordering changes materially when public-pilot proposal evidence is introduced, showing that the current joint ranking is sensitive to which repeated 1120-position data plane happens to surround the 64 fixed pilot positions. This does **not** mean the pilot itself is weak: Build32 known-mapping controls remain strong. It means candidate preservation/refinement cannot be validated on only one data-plane realization.

Exploratory translation-aware and beam-ranking variants were deliberately rejected when they turned one of the two MQ carriers green while selecting a false basin on the other. Those variants are not part of Build33. The next accepted decoder change must preserve/identify the correct basin on both carriers, keep matched negatives rejected, and only then proceed to v4 data sampling/HMAC authentication.

## Build34 MQ projective + authenticated frame result

Build34 turns the Build33 MQ observability result into an end-to-end authenticated decode. On the active MQ carrier under the fixed Build29 transform (angle 9.3 deg, scale 1.08x0.92, top/bottom inset 0.030/0.015 and asymmetric crop), the blind search recovers the exact parameter family, canonical origin `(0,0)` and approximately 0.001 normalized corner error. The corpus regression authenticates both `v4-b34-auth` and `v4-b34-alt` with the public test key after geometry has already been accepted. Typical held-out validation is 0.86-0.88 and pilot margin about 0.43.

The LQ carrier remains a SAFE REJECT and the legacy v3 research-limit reds are not reclassified. The v3 frozen core, pilot lock, tile geometry, data partition, frame/HMAC domains, Hamming baseline and Build29 acceptance thresholds are unchanged.

## Build35 physical-qualification readiness

Build35 does not claim new physical-channel performance yet. Its purpose is to make the next measurement reproducible after Build34 closed the active MQ blind projective+crop/HMAC blocker.

The preceding projective checkpoint is consistently identified as **Build34**. The Build34 corpus result remains unchanged: the MQ projective case is accepted with canonical origin `(0,0)`, validation about 0.873, score about 0.741, margin about 0.459 and ~0.001 normalized corner error; two independent robust v4 payloads authenticate after geometry acceptance.

Build35 exports that path through `ExperimentalV4ExtractProjective` / `v4-extract-projective` and adds a source-only API/CLI qualification target. It also defines a private scanner-first physical pack from the active MQ source: one unmarked control plus two robust/strength-24 carriers (`v4-b35-phys-a`, `v4-b35-phys-b`) using the public development key. The generated manifest freezes carrier dimensions and SHA-256; an acquisition plan freezes expected scan names and payload outcomes.

No physical PASS is recorded in this document until both marked paper scans recover their exact HMAC-authenticated payloads and the control rejects. The initial scanner gate is kept separate from later smartphone capture so channel degradation can be attributed rather than conflated.

## Build36 physical-channel result

- Real Build35 paper corpus acquired: control + two marked MQ prints.
- Primary usable scanner mode: 600-dpi color JPEG.
- Blind Build35 geometry: control REJECT; marked A/B SAFE-REJECT — physical blind qualification still open.
- Reference-assisted geometry diagnostic: marked A and marked B both recover exact v4 payloads and pass HMAC with Build36 soft Hamming.
- Interpretation: physical signal/ECC capacity is sufficient; blind scanner registration is the active blocker. No reprint or strength/ECC change is justified yet.

## Build37 blind scanner physical qualification

Build37 reuses the exact Build35 printed sheets and the first 600-dpi color-JPEG office-scanner captures. Unlike the Build36 channel-capacity diagnostic, no digital original is used as a registration reference.

The full-page paper/artwork boundary supplies only a geometric seed. A bounded scanner-affine search uses one half of the public v4 pilot for proposal and the disjoint half for held-out ranking. A complete-pilot origin competition then requires canonical origin `(0,0)`. The five best pilot-qualified geometries are frozen before protected data are read; their signed DCT margins are validation-weighted, soft-Hamming decoded and authenticated by the unchanged v4 HMAC.

Observed private-corpus result:

| capture | boundary | proposal | held-out | pilot score | margin | origin | result |
| --- | ---: | ---: | ---: | ---: | ---: | --- | --- |
| control | ~0.801 | ~0.306 | ~0.093 | ~0.080 | ~0.005 | non-zero | REJECT |
| marked A | ~0.808 | ~0.476 | ~0.294 | ~0.297 | ~0.174 | `(0,0)` | HMAC PASS `v4-b35-phys-a` |
| marked B | ~0.826 | ~0.483 | ~0.199 | ~0.213 | ~0.088 | `(0,0)` | HMAC PASS `v4-b35-phys-b` |

This is the first fully blind Format-v4 print/scanner/HMAC success. It is a scanner-channel result, not evidence for arbitrary smartphone photographs. Format v3 remains frozen and the v4 pilot/frame remain experimental/development-locked.


## Build40 pilot-only smartphone residual result

Build40 leaves the strength-48 carrier unchanged and adds a post-Build39 smooth residual field. On the deterministic synthetic phone-residual gate, proposal-only pilot controls fit a <=6 px quadratic field, held-out pilot validation improves, complete-pilot origin returns `(0,0)`, and the protected channel recovers the exact HMAC-authenticated v4 payload. The corresponding unmarked control does not qualify.

The private nine-photo strength-48 corpus does not yet satisfy the same residual assumptions under the current blind Build39 homographies. Example candidate fields can have enough proposal controls and sub-2 px internal fit RMS while **held-out validation decreases** and the complete-pilot origin remains non-zero. Those cases are rejected. This result narrows the remaining blind-phone blocker: the decoder must first retain/select a better global projective basin; a larger local warp is not supported by the evidence.


## Build41 blind smartphone basin result

Build41 re-runs the unchanged strength-48 nine-photo physical corpus after replacing only the global smartphone basin front end. The private success rule is staged rather than 6/6: all three unmarked controls must reject, and at least one A plus one B photograph must recover the exact expected HMAC-authenticated payload.

| capture | result |
|---|---|
| control/front | REJECT before protected-data acceptance |
| control/mild | REJECT before protected-data acceptance |
| control/angle | REJECT before protected-data acceptance |
| A/front | geometry reject |
| A/mild | geometry/data path reached; HMAC FAIL |
| A/angle | **HMAC PASS — `v4-b38-phone-a`** |
| B/front | **HMAC PASS — `v4-b38-phone-b`** |
| B/mild | geometry reject |
| B/angle | geometry reject |

Both authenticated cases are blind: the projective shortlist is generated only from structural evidence and the proposal folds of the public pilot, then frozen before the held-out fold is read. The Build40 residual fitter is evaluated downstream but is not applied in either PASS case because held-out residual evidence does not justify it. Thus the new physical success is attributable to Build41 global basin recovery plus the unchanged Build36 soft-Hamming/HMAC channel.

This closes the first smartphone physical milestone but not the whole channel. A/front, B/mild and B/angle still require stronger blind basin acquisition/retention, while A/mild is already close enough geometrically to reach protected-data decoding. Strength 48, ECC and HMAC remain unchanged.


## Build42 qualified-bank/list-decoder result

Build42 reuses the unchanged nine-photo strength-48 corpus and the unchanged Build41 geometry thresholds. The new gate is stricter about what changed: the two existing physical passes must remain direct Build41 decodes, `A/mild` must authenticate through the new post-geometry fallback, and all controls must still reject before protected-data decoding.

| capture | Build42 result |
|---|---|
| control/front | REJECT before data decode |
| control/mild | REJECT before data decode |
| control/angle | REJECT before data decode |
| A/front | geometry reject |
| A/mild | **HMAC PASS — `v4-b38-phone-a` via Build42 data fallback** |
| A/angle | **HMAC PASS — `v4-b38-phone-a`, direct Build41 path** |
| B/front | **HMAC PASS — `v4-b38-phone-b`, direct Build41 path** |
| B/mild | geometry reject |
| B/angle | geometry reject |

For A/mild the frozen Build41 bank contains five qualified geometries. Pair-only diagnostics do not authenticate. Deterministic three-member ensembles reduce the residual soft-Hamming ambiguity enough that the bounded list decoder reaches the authentic frame; the production fallback authenticates after two ensemble attempts. No pilot threshold, strength, ECC or HMAC domain was changed.

The remaining problem is now cleanly separated: A/front, B/mild and B/angle still fail in the geometry stage and therefore cannot benefit from Build42 data recovery. They are deferred to a later geometry-focused build.


## Build43 result

On the private Build38 strength-48 corpus the minimum Build43 milestone is controls 3/3 reject; A/front, A/mild, A/angle and B/front authenticate; B/mild/B-angle remain informational rejects. A/front is recovered through the new geometry fallback and unchanged Build42/HMAC data path.

Final Build43 qualification is defined with Go 1.25.1. A Go 1.26.0 rebuild of the same source changes the decoded raster of the canonical A/mild JPEG (Go 1.26 replaced `image/jpeg`) and can geometry-reject that case. Rebuilding with Go 1.25.1 restores A/mild HMAC PASS. The Build43 source package was therefore pinned to Go 1.25.1; Go 1.26+ JPEG physical behavior was deferred rather than masked by geometry retuning.

## Build44 deterministic JPEG checkpoint

Build44 adds a project-controlled pre-Go-1.26 JPEG decoder and public locked
raster vectors. The deterministic JPEG fixture and CLI-ingest regressions pass
under Go 1.26.0, and the unchanged private strength-48 physical gate also passes:
all three controls reject while A/front, A/mild, A/angle and B/front authenticate.
Build44 therefore promotes Go 1.26.0 as the qualified toolchain and is the latest
qualified smartphone milestone. Go 1.25.1 remains the historical Build43 reference.

## Build45 phone failure decomposition

Build45 introduces no production geometry or channel change. The private B/mild and B/angle captures are first treated as blind failures under the qualified Build44 path, then separately measured under a reference-assisted supplied-geometry oracle.

The retained-corpus oracle yields exact HMAC recovery for both difficult B captures on the first ML-Hamming frame: B/mild has proposal 0.384618, held-out validation 0.333015, full-pilot score 0.342214 / margin 0.207221 at origin `(0,0)`; B/angle has proposal 0.320779, held-out 0.230254, pilot 0.261081 / margin 0.140723 at `(0,0)`. SIFT/RANSAC registration is exceptionally coherent (2564/2584 and 3069/3089 inliers respectively).

This is not a blind decoder qualification. It is decisive channel-isolation evidence: with geometry supplied independently, the unchanged strength-48 carrier, Hamming/ECC, Build42 soft/list decoder and HMAC recover both B payloads. Build45 therefore directs subsequent work to blind geometry acquisition/ranking only.

## Build45 qualification-host blind decomposition — 2026-09-23

The retained private B captures were run through the unchanged Build44 production path on the qualification host after the Build45 oracle study.

| image | classification | Build41 direct | Build41 qualified | Build43 frozen | Build43 qualified | Build42 bank | HMAC |
|---|---|---:|---:|---:|---:|---:|---:|
| B/mild | qualification | false | 0 | 32 | 1 | 0 | false |
| B/angle | qualification | false | 0 | 28 | 1 | 0 | false |

The selected pair leaders differ (`bottom+right` / `left+right` for B/mild; `bottom+left` / `top+left` for B/angle), so the remaining failure is not explained by one universally missing artwork side. The important common feature is instead the **single held-out-qualified Build43 candidate**.

Build45's original `qualification` label is intentionally coarse. Production requires two qualified Build43 geometries before the direct ensemble is accepted; only then can that bank replace the Build41 bank. Build42 separately requires at least three geometries for its three-member list ensembles. Therefore `qualified=1, Build42 bank=0` is consistent with the designed production quorum and is not by itself an implementation handoff bug.

The Build45 reference oracle remains decisive channel evidence: both B/mild and B/angle authenticate `v4-b38-phone-b` on the first list frame under independently supplied geometry. Build46 is therefore restricted to inspecting the precision/data viability of the already-qualified singleton and the production quorum boundary.

## Build46 planned evidence gate

Build46 will classify each difficult B capture from the already-qualified Build43 candidates without changing production behavior. If the singleton authenticates under diagnostic single-candidate decoding, the result is `qualified-ensemble-shortfall`: the next geometry work should seek a second independently qualified basin rather than weaken data/ECC/HMAC. If the singleton does not authenticate, the result is `qualified-candidate-mismatch`: the next geometry work should improve the precision of the surviving blind basin. Post-hoc oracle corner error is recorded only after blind results exist.
## Build47 corrected frozen-bank host result

The corrected three-tier diagnostic was run on the qualification host. B/mild produces 32 candidates in the exact production tier with one qualified candidate and an oracle-nearest mean corner error of 59.734 px. Adding deeper cells on the same two selected side pairs increases the bank to 64 and qualified count to two but **does not improve** oracle-nearest error. Extending to all six side pairs reaches 128 candidates / 12 qualified and exposes a closer 34.049 px candidate from `top+left`, side-pair rank 3, cell rank 0. That oracle-nearest candidate is not held-out qualified.

B/angle produces 28 / 52 / 108 candidates across the three tiers. Its oracle-nearest error changes only from 5499.824 px to 5448.391 px and remains a distant false basin. This does not justify using B/angle to broaden production search.

The Build47 result therefore rejects a simple cell-depth explanation for B/mild and identifies a better basin in a lower-ranked side pair, while simultaneously showing that naive all-pair promotion would multiply held-out-qualified false basins.

## Build48 local-refinement checkpoint

Build48 changes no production decoder decision. From the corrected Build47 extended bank it selects at most two seeds per side-pair rank **only by proposal score**, locally refines all eight corner coordinates using proposal-fold evidence, freezes the complete refined bank, and only then evaluates held-out/full-pilot qualification. Diagnostic HMAC is attempted only after qualification. The reference-assisted oracle is generated after both blind refinement JSON outputs and is used solely to measure before/after corner error.

The primary private gate is B/mild: determine whether the Build47 ~34.049 px lower-ranked-pair basin can be reduced toward a few-pixel geometry without weakening any qualification threshold. B/angle remains informational.



## Build48 host result -> Build49 ranking question

On the qualification host, Build48 selected 12 seeds for B/mild, with 4 pre-qualified and 7 post-qualified candidates but 0 authenticated. The selected oracle-nearest seed moved from 46.202 px to 44.436 px mean corner error. This is worse than Build47's 34.049 px all-pair oracle-nearest candidate, showing that the top-2-per-pair proposal selector removed the best known geometry before refinement. B/angle remained ~5.45 kpx from the oracle. Build49 therefore freezes geometry and measures proposal-only ranking observability.


## Build49 qualification-host result -> Build50

Build49 measured 128 B/mild candidates / 12 held-out-qualified candidates. The oracle-nearest candidate (index 67) is 34.049 px from the reference, `top+left`, side-pair rank 3, all-pair-extension, cell rank 0. It is raw-proposal rank 17 globally but rank **3 within its pair**. The alternative Build49 observables rank it worse within the pair: fold-min 7, balanced 12, tile-consistency 13. Therefore Build48's top2-per-pair selector excluded the best known basin by exactly one rank, while the existing raw proposal score is already the strongest measured selector for that basin.

B/angle remains informational: the oracle-nearest candidate is still 5448.391 px away despite raw-proposal rank 3 within its side pair.

Build50 will test top4-per-pair refinement under the unchanged proposal score. No production ranking, quorum or threshold is changed.

## Build50 qualification-host result -> Build51 local surface study

The qualification host ran both Build50 targets under Go 1.26.0. For B/mild, the nested top2 view reproduces Build48-like behavior (12 seeds, 4 pre-qualified, 7 post-qualified, 0 authenticated, 46.202 -> 44.436 px nearest error). The top4 view includes 24 seeds and, as predicted by Build49, contains the 34.049 px `top+left` / pair-rank-3 / seed-rank-3 basin before refinement. The unchanged proposal-only coordinate descent moves the nearest refined result away to 43.168 px (max corner error 68.871 px); held-out-qualified candidates increase 7 -> 11 while HMAC remains 0.

B/angle remains a distant false-basin control: top4 starts at 5448.391 px and the nearest refined result is 5452.835 px, with 1 qualified candidate and 0 authenticated.

This rules out top2 seed pruning as the immediate B/mild bottleneck and shows that stronger proposal ascent cannot be assumed to improve physical registration. Build51 therefore instruments the exact refinement trajectory and a deterministic local neighborhood. It changes no production decision.


## Build51 qualification-host local-surface result -> Build52 optimizer study

Build51 was completed on the qualification host under the frozen Build44 raster/toolchain baseline. The primary B/mild target is candidate 10 (`top+left`, side-pair rank 3, seed rank 3), the same 34.049 px basin exposed by Builds47/49 and retained by Build50 top4.

| image / target | pre oracle | standard-refine oracle | accepted moves | accepted improve / worsen | best trace oracle | better-both stencil | classification |
|---|---:|---:|---:|---:|---:|---|---|
| B/mild c10 | 34.049 px | 43.168 px | 9 | 2 / 7 | 25.172 px (rejected) | yes: `corner-2-x +2` | `optimizer-opportunity` |
| B/angle c18 | 5448.391 px | 5453.906 px | 5 | 0 / 5 | distant false basin | no | `proposal-surface-misaligned` |

The B/mild `corner-2-x +2` stencil point is decisive because it improves both quantities without oracle guidance during generation: proposal rises `0.201986 -> 0.208953` (`+0.006966`) and post-hoc mean corner error falls `34.049 -> 32.442 px` (`-1.607 px`). The ordinary coarse-to-fine schedule does not test that `+2 px` move from the untouched seed: earlier accepted coarse moves have already changed the geometry by the time the 2 px stage is reached.

Build52 therefore changes only optimizer scheduling/state retention. It compares the unchanged coarse-to-fine endpoint with an independent `2 -> 1 px` proposal-only restart from every untouched top4 seed and retains every accepted intermediate restart state. Geometry is frozen before held-out/full-pilot qualification, and HMAC/oracle remain downstream diagnostics. No Build52 result is a production qualification until the private host study is run and any candidate production change separately re-passes the complete Build44 physical gate.


## Build52 qualification-host result -> Build53 coupled pair escape

The Go 1.26.0 qualification host completed Build52. B/mild candidate 10 remains the primary target. The untouched seed is 34.049 px from the private post-hoc oracle; the unchanged coarse-to-fine baseline ends at 43.168 px. The fine restart retains 8 states. Its best oracle state is 31.747 px (`-2.302 px` versus the seed) but does not qualify; retained state 1 is held-out-qualified at 32.442 px and no fine state authenticates. Build52 therefore classifies B/mild as `optimizer-partial-gain`. B/angle remains a distant false basin and is `optimizer-no-gain`.

The state sequence narrows the optimizer question further. On B/mild the first accepted 2px state (`corner-2-x +2`) is qualified at 32.442 px. The next 2px move raises proposal again and improves oracle to 31.747 px but destroys qualification. A local development probe (not qualification evidence) finds that the qualified 32.442 px root has no proposal-improving individual +/-1px coordinate move, yet it does have proposal-improving **coupled** +/-1px two-coordinate moves. One such pair, `corner-1-y -1` plus `corner-2-x +1`, remains qualified and reduces post-hoc oracle error to ~31.842 px.

Build53 therefore changes only optimizer neighborhood topology. It retains proposal-only 2px roots, uses the complete +/-1px single-coordinate stencil to identify coordinate-local roots, and only there scans all 112 coupled +/-1px two-coordinate combinations. At most eight proposal-improving pair states per root are retained by proposal rank. Geometry is frozen before qualification/HMAC; reference/SIFT remains post-hoc. The qualified host Build53 run is required before this preflight evidence can be interpreted.


## Build53 qualification-host result -> Build54 post-pair continuation

The qualified Go 1.26.0 Build53 run confirms the coupled escape on the primary B/mild target. Candidate 10 (`top+left`, pair rank 3, seed rank 3) contains one 1px coordinate-local 2px root at 32.442 px. The complete coupled stencil retains four proposal-improving pair states and **all four remain held-out-qualified**. Pair rank 2 (`corner-1-y -1` + `corner-2-x +1`) improves proposal from 0.208953 to 0.210454 and reduces post-hoc mean corner error from 32.442 to 31.842 px. No pair state authenticates, so B/mild is `pair-qualified-gain`. The B/angle target has no coordinate-local root and is `pair-not-triggered`.

A local deterministic design probe then continues every retained B/mild pair state using only the unchanged proposal objective. Pair rank 1 re-enters an ordinary 1px coordinate-ascent path; retaining intermediates exposes a held-out-qualified state near 30.539 px before later proposal ascent moves away from the oracle again. No preflight continuation state authenticates. This motivates Build54 but is not qualification evidence.

Build54 therefore reproduces the full Build53 pair bank, continues **every** retained pair state with bounded 1px proposal-only coordinate descent, retains every accepted intermediate, freezes the entire bank, and only then annotates qualification/HMAC and post-hoc oracle error.


## Build54 host result -> Build55 sibling-stencil hypothesis — 2026-09-24

**Build54 result.** B/mild candidate 10 retains 34 continuation states. Pair-rank-1 continuation state 7 reaches **30.539 px** mean post-hoc oracle error from a 32.447 px parent while remaining qualified (`proposal=0.257772`, `validation=0.110310`); HMAC remains false. Classification: `continuation-qualified-gain`. The B/angle target is `continuation-not-triggered`.

**Order-bias probe, not qualification evidence.** Re-evaluating every +/-1px neighbor from the exact frozen 30.539 px state shows that `corner 2 / x +1 px` improves proposal to `0.266832` and oracle error to about `29.778 px`, with `validation=0.163706` and unchanged qualification. A single-candidate diagnostic decode still fails HMAC. The ordinary Build54 continuation misses that sibling because another coordinate is accepted earlier and changes the parent before corner-2-x is evaluated.

**Build55 hypothesis.** The proposal objective still contains useful local directions, but sequential Gauss-Seidel ordering can hide them. Build55 therefore evaluates the complete independent +/-1px sibling stencil from every retained continuation state and freezes every proposal-improving sibling before qualification/HMAC.

## Build55 host result -> Build56 sibling pair-escape hypothesis — 2026-09-24

**Build55 host result.** The qualified Go 1.26.0 artifact confirms `sibling-qualified-gain` on B/mild candidate 10 (`top+left`, side-pair rank 3, seed rank 3). The target retains 91 proposal-improving siblings. The best qualified sibling is rank 3 from continuation state 7 of pair-rank 1: `dimension 4 +1 px`, proposal `0.257772 -> 0.266832`, held-out validation `0.163706`, and independent oracle error `30.539 -> 29.778 px`. HMAC remains false. The B/angle target remains `sibling-not-triggered`. Across the complete blind outputs B/mild freezes 93 siblings (76 qualified, 0 authenticated) and B/angle freezes 106 (0 qualified, 0 authenticated).

**Second local-maximum probe, not qualification evidence.** Re-evaluating the exact 29.778 px sibling against all 16 independent +/-1px single-coordinate moves finds zero proposal improvement. A bounded 112-combination coupled two-coordinate stencil nevertheless finds five proposal-improving pair states; all five are qualified in the local probe and none authenticates. Pair rank 2 (`dimension 2 -1 px` + `dimension 4 +1 px`) raises proposal to `0.273793` and reduces independent oracle error to about `28.514 px`.

**Build56 hypothesis.** The unchanged proposal objective still contains useful directions beyond the Build55 sibling local maximum, but they require another coupled move. Build56 must not choose the 29.778 px state using oracle or held-out evidence. Instead it reproduces every Build55 sibling, tests complete one-coordinate locality proposal-only, and runs the pair stencil on every proposal-local sibling. Retain at most eight pair states by proposal alone, freeze both image banks, then annotate qualification/HMAC and finally oracle.

## Build56 host result -> Build57 second-pair continuation hypothesis — 2026-09-24

The qualified-host Build56 artifact confirms the second pair escape on the primary B/mild target. Candidate 10 (`top+left`, pair rank 3, seed rank 3) contains 91 Build55 siblings, 16 proposal-local target siblings and 65 retained second-pair states. The best qualified second-pair state is root 1 / first-pair rank 1 / continuation state 7 / sibling rank 3 / second-pair rank 2 (`dimension 2 -1 px` + `dimension 4 +1 px`): proposal `0.266832 -> 0.273793`, post-hoc oracle `29.778 -> 28.514 px`, qualification true, HMAC false. This is `sibling-pair-qualified-gain`.

Across the complete Build56 banks, B/mild freezes 68 second-pair states (51 qualified, 0 authenticated) and B/angle freezes 176 (0 qualified, 0 authenticated). Target B/angle candidate 18 remains `sibling-pair-not-triggered`.

**Build57 hypothesis.** The unchanged proposal objective still contains useful geometry after the second pair escape. Build57 must not oracle-select the 28.514 px state. Instead it reproduces every retained Build56 second-pair state and continues each with the established bounded 1px proposal-only coordinate descent, retaining every accepted intermediate. Freeze both image banks before qualification/HMAC; generate/reference SIFT only after both blind outputs exist.


## Build57 host result -> Build58 second-pair sibling-stencil hypothesis — 2026-09-24

**Build57 result.** The qualified Go 1.26.0 host closes B/mild as `second-pair-continuation-qualified-gain`. Target candidate 10 retains 65 Build56 second-pair parents and 154 Build57 continuation states; 147 target continuation states pass unchanged qualification and none authenticate. The best qualified continuation reaches 26.968 px post-hoc oracle error, `-2.581 px` relative to its exact second-pair parent and `-7.080 px` relative to the original 34.049 px seed. The same branch then accepts further proposal-improving coordinates while moving away from that geometric minimum. B/angle target candidate 18 remains `second-pair-continuation-not-triggered`; the complete B/angle continuation bank contains 347 states, 0 qualified and 0 authenticated.

**Build58 hypothesis.** The repeated divergence after the geometric minimum may again be caused by Gauss-Seidel coordinate order rather than by the proposal observable itself. Build58 therefore must not select the 26.968 px state. It reproduces every Build57 continuation state and evaluates all 16 independent +/-1px coordinate siblings from the identical frozen parent. Retain every proposal-improving sibling, freeze both image banks, then apply unchanged qualification/HMAC and only afterwards generate SIFT/reference oracle geometry.

**Leakage barrier.** Sibling generation and retention are proposal-only. Held-out/full-pilot qualification, secret key, payload, ECC/HMAC and oracle geometry cannot affect the sibling bank. No production score, threshold, quorum, encoder, pilot, ECC/Hamming or HMAC domain changes in Build58.

## Build58 host result -> Build59 third-pair escape hypothesis — 2026-09-24

**Build58 result.** The qualified Go 1.26.0 host closes B/mild as `second-pair-sibling-qualified-gain`. Across the complete B/mild bank, Build58 freezes 275 sibling states; 271 pass unchanged qualification and 0 authenticate. Target candidate 10 improves from its exact Build57 parent at 26.968 px to a qualified sibling at 26.149 px while proposal rises from 0.276536 to 0.285435. B/angle freezes 542 siblings, 0 qualified and 0 authenticated; target candidate 18 remains not triggered.

**Build59 hypothesis.** The unchanged proposal objective may still contain useful geometry beyond the Build58 sibling bank, but the next transition should only occur after proving one-coordinate locality. Build59 must not oracle-select the 26.149 px state. It reproduces every Build58 sibling, probes the complete independent +/-1px one-coordinate stencil from that exact sibling, and only for proposal-local siblings evaluates the bounded 112 coupled two-coordinate +/-1px moves, retaining at most eight by proposal. Freeze both image banks before qualification/HMAC and generate SIFT/reference only afterwards.

**Leakage barrier.** All local-max tests, pair generation and pair ranking are proposal-only. Qualification, secret key, payload, ECC/HMAC and oracle geometry cannot affect the Build59 bank. No production score, threshold, quorum, encoder, pilot, ECC/Hamming or HMAC domain changes.

## Build59 host result and Build60 continuation hypothesis — 2026-09-24

**Build59 result.** The qualified Go 1.26.0 host confirms `third-pair-qualified-gain` on B/mild. The full blind B/mild bank contains 180 third-pair states, 175 qualified and 0 authenticated. Target candidate 10 reaches a best qualified third-pair state at 25.407 px post-hoc oracle error, 1.664 px better than its exact 27.071 px Build58 sibling parent while the unchanged proposal also improves. B/angle target candidate 18 remains `third-pair-not-triggered`; its full Build59 bank contains 914 third-pair states, 0 qualified and 0 authenticated.

**Build60 hypothesis.** The third pair escape may cross another proposal-local barrier without reaching the terminal useful basin. Build60 therefore continues every retained Build59 third-pair state, not an oracle-selected subset, with bounded 1px proposal-only coordinate descent and retains every accepted intermediate. The complete B/mild+B/angle continuation bank is frozen before qualification/HMAC; SIFT/reference remains post-hoc only.

**Leakage barrier.** Third-pair continuation generation, acceptance, stopping and retention use only the unchanged proposal objective. Qualification, secret key, payload, ECC/HMAC and oracle geometry cannot affect the Build60 geometry bank. No production score, threshold, quorum, encoder, pilot, ECC/Hamming or HMAC domain changes.



## Build60 host result and Build61 sibling hypothesis — 2026-09-24

**Build60 host result.** The qualified Go 1.26.0 artifact closes B/mild as `third-pair-continuation-qualified-gain`: 406 continuation states are retained, 397 qualify, none authenticate, and the best qualified state reaches 25.150 px. That state improves its exact 29.938 px third-pair parent by 4.788 px while proposal also improves. B/angle target candidate 18 remains not triggered; across the full non-target bank 1467 continuation states are retained, 2 qualify at about 5665 px oracle error and 0 authenticate.

**Build61 hypothesis.** The continuation bank can still hide order-dependent one-coordinate alternatives because Gauss-Seidel updates alter later coordinate probes. Build61 therefore evaluates all 16 independent +/-1px siblings from every retained Build60 continuation state against the identical frozen parent, retains proposal-improving siblings, freezes both image banks, and only then applies qualification/HMAC.

**Control escalation.** Build61 reports whole-bank qualification counts and rates explicitly. The two B/angle qualified states in Build60 are far from oracle and unauthenticated, but any material increase in such leakage is evidence against unconstrained further optimizer expansion.


## Build61 host result and Build62 fourth-pair hypothesis — 2026-09-24

**Build61 host result.** The qualified Go 1.26.0 artifact closes B/mild as `third-pair-sibling-qualified-gain`: 593 proposal-improving siblings are frozen, 591 qualify, none authenticate, and the best qualified sibling reaches 24.614 px. The exact parent of that sibling is 25.359 px, so the local oracle improvement is 0.745 px while proposal rises to 0.288064. The separate best Build60 continuation parent remains 25.150 px and is not the genealogy parent of the winning sibling. B/angle produces 1902 sibling states, 0 qualified and 0 authenticated, so the rare Build60 qualification leakage does not propagate through the sibling stage.

**Build62 hypothesis.** A Build61 sibling may itself be proposal-local in the one-coordinate +/-1px neighborhood while a coupled two-coordinate move still improves the unchanged proposal objective. Build62 therefore probes every frozen Build61 third-pair sibling with the full one-coordinate stencil and applies the bounded 112-combination pair stencil only when zero one-coordinate improvements exist, retaining at most eight proposal-ranked fourth-pair states.

**Leakage barrier.** Locality classification, pair generation, acceptance and ranking use proposal only. The complete B/mild+B/angle fourth-pair bank is frozen before qualification/HMAC; SIFT/reference remains post-hoc. Whole-bank B/angle qualification/HMAC rates remain explicit because Build60 demonstrated rare qualification leakage under larger banks.


## Build62 host result -> Build63 fourth-pair continuation — 2026-09-24

Build62 closes B/mild as `fourth-pair-qualified-gain`: 492 fourth-pair states are frozen on B/mild, 491 qualify and 0 authenticate. The best qualified target state reaches 24.056 px from its exact 25.641 px Build61 sibling parent, a -1.585 px local oracle gain. B/angle freezes 3746 fourth-pair states with 0 qualified and 0 authenticated. Build63 therefore continues every retained Build62 fourth-pair state with bounded proposal-only 1px coordinate descent, retains every accepted intermediate, freezes both image banks, and only then annotates qualification/HMAC and post-hoc oracle geometry.

## Build63 first blind B/mild recovery -> Build64 qualified production baseline — 2026-09-24

The qualified Go 1.26.0 Build63 host artifact is the first blind recovery of retained `phone-b-mild.jpg`. The final fourth-pair continuation bank contains **937 states, 935 held-out/full-pilot-qualified states and exactly one authenticated state**. The authenticated state belongs to candidate 10 (`top+left`, side-pair rank 3, seed rank 3), fourth-pair rank 7, continuation index 7. Its proposal is `0.302815`, validation `0.240486`, post-hoc oracle error `28.868 px`, and HMAC is true.

The same bank contains a geometrically closer state at `22.099 px` that does not authenticate. Therefore oracle closeness neither generated nor selected the successful state; it remains an independent post-hoc measurement. `B/angle` expands to **6198** fourth-pair continuation states with **0 qualified and 0 authenticated**.

This is the first evidence strong enough to stop the optimizer-only diagnostic chain and open a production-candidate experiment. Build64 adds the exact Build63-derived search family only as a final fallback after the Build44-qualified paths. The complete deep geometry bank must be frozen proposal-only before qualification or HMAC, and the full nine-photo physical matrix must pass before any baseline promotion.


## Build64 physical qualification -> qualified baseline — 2026-09-24

The complete retained nine-photo Build38 smartphone matrix passes under Go 1.26.0 with the Build64 deep fallback enabled. All three controls reject. A/front authenticates through the historical Build43 path; A/mild, A/angle and B/front authenticate through their existing paths without invoking Build64. B/mild reaches Build64 after the historical path fails and authenticates exact `v4-b38-phone-b` from a bank of 937 frozen states / 935 qualified states. B/angle reaches Build64, freezes 6198 states, qualifies none and remains rejected.

This closes the Build45–63 optimizer research line as a successful production promotion. Build64 supersedes Build44 as the qualified smartphone physical-recovery baseline while preserving Build44's deterministic JPEG ingest, Go 1.26.0 toolchain qualification and all previously qualified recovery behavior.


## Build64 qualified baseline -> Build65 scheduling-only candidate — 2026-09-24

Build64 closes the difficult B/mild recovery and becomes the qualified smartphone baseline, but its deep fallback is intentionally exhaustive and expensive. Build65 changes no search or acceptance rule. It runs the same already-selected Build48 seed branches concurrently, stores each branch by seed index, concatenates them in Build64 order, then runs unchanged qualification and protected-data/HMAC decode serially. The private Build65 gate requires exact Build64 deep-recovery telemetry counts and the same nine-photo outcomes before performance is considered. At candidate creation time no Build65 host result had been recorded; the qualified-host result is documented immediately below.

## Build65 equivalence host result -> Build66 ordered-decode candidate — 2026-09-24

The complete Build65 nine-photo gate passes and reproduces the qualified Build64 semantics exactly. Every fallback case matches the Build64 seed count, geometry-evaluation count, frozen bank, qualified bank, logical decode-candidate count and logical list-frame count; the four historical fast authenticated cases still return before deep recovery. B/mild remains `937/935/691/2120047` for bank/qualified/decode-candidates/list-frames and authenticates exact `v4-b38-phone-b`; B/angle remains a reject at `6198/0`. Build65 uses 8 seed workers on this host, but it is not promoted: Build64 did not record matched `elapsed_ms`, and Build65 B/mild still costs 828036 ms, so the performance case is insufficiently strong.

The next measurable bottleneck is protected-data decoding rather than blind geometry generation. Build66 therefore preserves Build65's exact seed-parallel bank and runs qualified single-candidate decodes in bounded parallel batches. Results are committed only in original Build64 order, preserving the first logical HMAC success and exact logical telemetry. Build64 remains the qualified baseline while Build66 is evaluated.


## Build66 same-host qualification -> qualified performance baseline — 2026-09-25

The complete Build66 nine-photo gate passed on the same Go 1.26.0 host used for Build65. Semantic equivalence is exact: every deep fallback reproduces Build64/65 seed, geometry-evaluation, bank, qualification, logical decode-candidate and logical list-frame counts, and every payload/HMAC outcome is unchanged. B/mild remains `937/935/691/2120047` and authenticates exact `v4-b38-phone-b`; B/angle remains `6198/0` and rejects.

Performance is materially better on the decode-bound success path: B/mild drops from 828036 ms to 544038 ms, a 1.522x speedup / 34.3% wall-clock reduction. Across the complete nine-photo matrix, total time drops from 2372255 ms to 2175687 ms, an 8.3% reduction. Build66 therefore supersedes Build64 as the current qualified smartphone baseline.

## Build67 profiling candidate — 2026-09-25

Build67 is an observability-only successor to the qualified Build66 smartphone baseline. It preserves the exact Build64/65/66 logical search/decode behavior and adds stage wall timing plus physical-work counters around the deep fallback. In particular, Build67 distinguishes the qualified logical candidate/list-frame prefix from physically executed speculative candidates in the winning Build66 batch and separates summed projective-margin sampling worker time from summed list/Hamming/HMAC worker time.

No Build67 physical result is recorded in this source snapshot yet. Build66 remains the qualified baseline until the Build67 retained-corpus profiling run is returned and a separate Build68 optimization is justified from that evidence.
