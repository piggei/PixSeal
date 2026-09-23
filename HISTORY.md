## Build44 — deterministic JPEG ingest

Build44 begins by isolating input rasterization from the Go toolchain. Build43 qualification demonstrated that Go 1.26.0 and Go 1.25.1 decode the exact same canonical smartphone JPEG into different Y/Cb/Cr planes, enough to move A/mild outside the qualified geometry path. Build44 does not retune geometry around that difference. Instead it vendors the pre-Go-1.26 pure-Go JPEG decoder as `internal/jpeglegacy` and routes CLI JPEG input through that project-controlled path.

The watermark algorithm remains the Build43 baseline: Format-v3 stays frozen; Format-v4 encoding, locked pilot, strength 48, data mapping, Hamming/ECC, whitening/HMAC domains, Build41/43 geometry and Build42 list decoding are unchanged. A public deterministic JPEG fixture locks the decoder raster contract. On the Surface/WSL2 qualification host, Go 1.26.0 passed both the deterministic-raster regression and the complete private Build43 physical matrix (controls 3/3 reject; A/front, A/mild, A/angle and B/front authenticate). Go 1.26.0 is therefore promoted as the qualified Build44 toolchain; Go 1.25.1 remains the historical Build43 reference.

# PixSeal history

This file records technical evolution, including experiments that were later
superseded or deliberately not promoted.

## 2026-09-18 — licensing transition

Current and future PixSeal source distributions move from MIT to the PolyForm
Noncommercial License 1.0.0 (`PolyForm-Noncommercial-1.0.0`). The project is
therefore described as source-available for noncommercial use; commercial use
requires a separate commercial license. Copies already distributed under MIT
retain their original MIT grant. This transition changes licensing and
documentation only and does not create a new algorithmic build or alter any
Format-v3/Format-v4 behavior. See [`docs/LICENSING.md`](docs/LICENSING.md).

## v0.1.0

First stable public-development baseline: classical DCT embedding for short
messages, bounded extraction, JPEG/resize/crop robustness harnesses and pure-Go
CLI/core architecture.

## v0.2.0 build 1 — adaptive Format v3

Introduced `robust`, `balanced`, `capacity` and `auto`, fixed 35x32 tile,
authenticated v3 headers, adaptive redundancy and the `analyze` command.

## build 2 — v3-only runtime

Removed runtime v1/v2 decoding because those experimental formats had no known
external compatibility population. Their history remained documented.

## build 3 — arbitrary digital rotation

Added bounded quarter-turn and arbitrary-angle recovery without changing the
v3 on-image format.

## build 4 — combined digital geometry

Extended rotation to selected resize/crop combinations and added
`geometry-test`. Aggressive rotation+50% resize remained a documented signal
limit rather than a claimed feature.

## build 5 — axis-aligned affine recovery

Added a fixed 36-hypothesis anisotropic-scale/shear bank with virtual affine
sampling and bounded full-grid promotion. Added `affine-test`.

## build 6 — first composed affine experiment

Demonstrated one bounded anisotropic-scale+rotation composition on a synthetic
case, but the rotation-first heuristic did not generalize to PJ's real corpus.
The failure was retained as a research result.

## build 7 — direct composed-lattice recovery

Reworked the composed case so geometry was scored directly rather than inherited
from a global rotation peak. The two real regression carriers that had produced
FAIL/TIMEOUT in build 6 both authenticated successfully.

## build 8 — explicit lattice basis bank

Generalized the direct bank to both 110x90 and 90x110 anisotropies under
rotation.

## build 9 — broader fixed lattice bank and unified reports

Added 105x95 and 95x105 shapes, giving four fixed direct lattice bases. Added
`make all-test` for sequential comparative reporting. A more continuous local
`u`,`v` refinement prototype improved positive cases but made negative extraction
too expensive and was deliberately not promoted.

## build 10 — resize regression recovery

The build-9 report exposed fractional-resize regressions caused by speculative
geometry suppressing the mature resize path. Build 10 restored ordering with
virtual isotropic-scale recovery and a targeted bilinear fallback. It also fixed
large-PNG dimension probing in shell suites.

## build 11 — first projective step

Added exactly two fixed 4% vertical-keystone homographies using virtual
projective sampling. ImageMagick-generated regression cases authenticated on the
private LQ/MQ corpus. This was recorded as a stepping stone, not general
perspective support.

## build 12 — local-consensus experiment not promoted

A reconstructed local orientation-consensus fallback failed to improve the
private geometry matrix and increased runtime. Build 11 therefore remained the
algorithmic basis for release consolidation.

## v0.2.0-rc1 — first release candidate

Consolidated build 11, froze the Format v3 deterministic encoder and prepared
the v0.2 line for release hardening.

## External audit 1

An independent source/package audit found no reason to change Format v3 but
identified CLI validation, output-path safety/umask, perspective-harness,
large-image-preflight, regression-coverage and documentation issues.

## v0.2.0-rc2 — first audit reconciliation

RC2 implemented those fixes while keeping the on-image format and encoder
fingerprints unchanged. It also separated stable `release-unit` tests from
experimental `research-unit` regressions.

PJ's 2026-09-09 Linux/WSL2 qualification report recorded:

```text
Release baseline: PASS
deep-test:        72/72 PASS
geometry-test:    65/120 PASS, 55 FAIL (experimental)
affine-test:      22/24 PASS, 2 FAIL (experimental)
composition-test: 2/2 PASS
lattice-test:     8/8 PASS
perspective-test: 4/4 PASS
```

## External audit 2

The RC2 audit confirmed the core and most first-round fixes, then identified a
Windows no-clobber assumption, non-strict release-gate semantics, zero-case
qualification, help exit status, capacity full decoding and additional hardening
items.

## v0.2.0-rc3 — second hardening candidate

RC3 introduced a stronger cross-platform no-clobber design using hard links with
an exclusive-create copy fallback, overflow-safe image-size arithmetic, shared
alpha flattening and a stronger projective synthetic regression. During package
assembly, however, several valid RC2 release-engineering/script changes were
reverted and the RC2 qualification history was misattributed.

## External audit 3

The RC3 audit found no Format v3, encoder or decoder regression. It identified
the release-engineering reversions, unresolved strict/zero-case/help/capacity
items, the raw stdout/stderr regression and an incorrect direct-lattice budget.

## v0.2.0-rc4 — reconciliation candidate

RC4 is deliberately non-algorithmic: it keeps RC3 filesystem/overflow/alpha
hardening, restores the valid RC2 release/test-harness changes, fixes the
remaining CLI/release-gate issues, restores the 300 MP v0.2 source policy and
repairs release chronology/documentation. Format v3 and the frozen search banks
remain unchanged.

## v0.2.0 — final release

RC4 passed the stable release baseline on the private qualification corpus and
was promoted to v0.2.0 without further functional changes to Format v3, the
encoder or decoder search banks. `deep-test` recovered 72/72 baseline cases.
The initial RC4 research geometry run contained two 60-second timeouts; a
targeted 120-second rerun recovered both and reproduced the RC2 experimental
reference of 65/120 with zero timeouts. The remaining geometry/affine failures
are retained as measured research limits rather than release-baseline failures.

## v0.3.0-build1 — local lattice diagnostic

Started the v0.3 line from the final v0.2.0 baseline without changing Format v3,
the deterministic encoder or the production extraction search order. Added a
separate bounded diagnostic estimator that measures local lattice basis vectors,
phase, DCT margin, periodic/tile repetition coherence and cross-region consensus.
The diagnostic path uses sampled luminance working planes and exposes explicit
search budgets/timings through CLI JSON.

The first checkpoint intentionally stops before integration into `ExtractWithInfo`.
Canonical marked/unmarked discrimination is covered by new regressions, while
arbitrary-angle automatic estimation remains a measured limitation. A private
print-camera harness was added with HMAC-only PASS semantics and clean SKIP when
the undistributed real corpus is unavailable.
## v0.3.0-build2 — bounded projective diagnostic

Extended the separate v0.3 diagnostic path with a coarse print-boundary prior,
boundary-normalized local lattice consensus, multiscale native-period support, a
bounded projective scale shortlist and explicit homographies. Added a virtual
projective DCT sampler so large smartphone captures can reach the existing v3
frame/HMAC decoder without creating a full rectified image.

A synthetic translated/projective diagnostic regression authenticates through
the new virtual path and rejects a wrong key. On PJ's two original ~200 MP
print-camera captures, both photographs now produce coherent lattice evidence
and bounded projective candidates; neither yet produces a valid HMAC. The
hidden message therefore remains unknown, as required.


## v0.3.0-build3 — spatial phase refinement

Added a second geometric discriminator after the build2 lattice/scale stage:
Format-v3 sync phase is measured independently in spatially separated complete
tiles. Correct geometry should keep profile and lattice phase coherent across
the print; isolated aggregate header peaks are therefore no longer sufficient
to rank a candidate.

The strongest bounded phase seeds receive ±1%/±2% scale refinement and, when
four phase observations are available, a deterministic DLT homography fit from
phase correspondences. Canonical phase corrections are capped and only four
full virtual v3 decodes may reach HMAC. On the two private smartphone captures
this materially sharpens geometric evidence but still does not authenticate the
unknown payload.


## v0.3.0-build4 — fundamental-scale and block-origin refinement

Build4 converted the build3 observation that phase maxima can be harmonics into
an explicit selection stage. Projective scale clusters now retain cross-level
and cross-region support, and a low-frequency candidate is marked fundamental
only when it has sufficient independent support. This independently selects the
closely matching ~1668x1250 and ~1653x1254 families in the two private captures.

The diagnostic authentication path then performs two bounded refinements without
changing Format v3: a ±2%/±1%/±0.5% search confined to the selected fundamental
family, and a modulo-8 block-origin search confined to one DCT cell. Sparse probe
winners are always checked against the unmodified seed with the ordinary probe,
so a refinement cannot be promoted merely because a cheaper sampler produced a
false maximum. A single key-independent smooth residual warp can additionally be
fit from the local 3x3 lattice field.

These changes improve geometric interpretation but do not yet authenticate the
real corpus. The frontal best global header score remains z≈4.45 and the inclined
best remains z≈3.97; neither reaches a valid Format-v3 HMAC. The hidden payload
therefore remains unknown.


## v0.3.0-build5 — bounded photometric channel experiment

Build5 freezes the build4 geometry and tests whether the remaining real-camera
loss is partly photometric. Three deterministic post-geometry views are used:
raw luminance, local mean/variance normalization and a mild high-pass. Only
already-retained geometries are evaluated, so the bank cannot expand geometric
search; at most nine header probes and four complete v3/HMAC decodes are
allowed.

The frontal private capture shows a modest but reproducible improvement under
mild high-pass (z≈4.45 to z≈4.70), while the inclined capture does not improve
over raw. Neither authenticates, so the experiment provides channel evidence
without claiming recovery. The all-test harness is also reorganized so release
invariants, qualification-corpus observations and experimental geometry are
reported separately when different workstations use different image corpora.


## v0.3.0-build6 — lattice-first fallback and adaptive pyramid

Build6 used the first truly held-out print-camera pair to remove an architectural
dependency exposed by build5. A strong physical print boundary is no longer a
mandatory gate. Large images without a strong boundary may receive one bounded
finer lattice level when the cheap pass is already borderline; only after that
independent lattice evidence exists can a sane weak quadrilateral initialize the
existing projective path.

The held-out frontal bicycle capture moves from no projective attempt to a full
four-attempt bounded projective run, while the held-out inclined capture and the
two original private captures remain on their previous strong-boundary paths. No
real capture authenticates yet. A wrong-key control can still obtain an apparently
strong sync score without HMAC, reinforcing the project rule that sync/lattice
evidence is diagnostic only.


## v0.3.0-build7 — protected-bit channel instrumentation

Build7 turns the remaining print-camera failures into measurable coding-channel
failures without using the hidden message as an oracle. The diagnostic reuses
only the four geometry/photometric grids already budgeted for complete v3 decode.
Because Format v3 fixes magic plus version/profile before the payload, six full
Hamming(7,4) words are known exactly after key-derived whitening and can be
measured independently of payload contents.

The original frontal photograph is now quantitatively the closest current case:
its best bounded view leaves only one of six known Hamming words with more than
one error and two of 24 known header bits wrong after ordinary ECC. Its wrong
known bits also have lower DCT margin than correct bits. The other three captures
show two or three known words beyond the one-bit correction radius, with the
inclined bicycle capture showing high-confidence wrong bits. This separates a
likely soft-reliability opportunity from cases that still look geometrically or
optically misregistered.

The private print-camera shell harness is also generalized from two hardcoded
historical filenames to deterministic discovery of every PNG/JPEG in the corpus.
PASS semantics remain per-image HMAC only; the corpus itself is never packaged.


## v0.3.0-build8 — bounded scanner authentication and reliability experiment

Build8 used the first scanner acquisitions as a controlled contrast with the
smartphone channel. The scanner images are only about 34.8 MP, so build7 attempted
the mature production extractor before projective research; that path could exceed
120 seconds even though key-independent lattice geometry finished in roughly seven
seconds. Build8 gives `diagnose` a smaller diagnostic-only baseline budget and bounds
large virtual-grid aggregation by block count/spatially distributed complete tiles.
The same scanner inputs then complete full diagnostics in roughly 10--11 seconds.

Build8 also tests a single-output maximum-likelihood Hamming(7,4) decoder driven by
signed DCT margins. Synthetic tests show the expected advantage on a weak double
error, including successful Format-v3 HMAC where hard Hamming fails. The real frontal
smartphone capture qualifies for one such decode but does not improve its fixed known
header and still fails HMAC. The scanner cases likewise do not authenticate and show
that lattice consistency near 0.87 can coexist with multiple protected-bit errors.
This result argues against treating geometry consistency as a proxy for payload
recoverability and against promoting soft Hamming as a universal fix.

## v0.3.0-build9 — spatial repetition diagnosis

Added zero-extra-read 3x3 spatial accumulation to the diagnostic full-grid path.
The six real scanner/smartphone acquisitions showed only about 0.64--0.68 mean
key-independent tile-position sign agreement despite strong macro-lattice evidence.
Known protected header errors were overwhelmingly spatially mixed rather than
stably wrong. A bounded +/-2-block local phase oracle helped selected cases but
hurt others, so independent local re-phasing was deliberately not promoted. The
result redirected the next experiment toward a smooth low-DOF spatial phase/warp
model instead of more candidate or photometric expansion.

## v0.3.0-build10 — bounded smooth phase-field experiment

Build9 showed that protected-bit errors are overwhelmingly mixed across spatial
regions and that independent local phase choices are too unstable to promote.
Build10 therefore fitted the smallest useful continuous model: two affine fields
with three coefficients each, one for X and one for Y correction in canonical
block coordinates.

The field is derived from the already-bounded +/-2-block cell phase observations.
At least six cells are required. A loose gate controls whether the field is worth
resampling at all, while a stricter fit/leave-one-out gate controls whether the
corrected grid may replace one existing HMAC-decode slot. The four-slot HMAC
ceiling is unchanged and smooth-field resampling is capped at two.

Synthetic testing proves that an exact affine drift can be inverted and the
original Format-v3 payload authenticated. Real data are mixed: four of six current
acquisitions show a better known prefix on at least one measured smooth field, but
two worsen and only one candidate passes the strict decode gate. That candidate
reduces the known post-ECC prefix from 6 to 5 errors and still fails HMAC.

The evidence therefore supports a real smooth component but not a single global
affine correction as a complete print-acquisition model. Future work should first
improve the independence/quality of spatial control observations or test a still
bounded regularized low-order model; it should not simply expand per-cell phase
search or HMAC candidate count.


## v0.3.0-build11 — confidence-weighted continuous phase controls

Build10 established that a low-order spatial phase field can sometimes improve a
real candidate, but the control observations were still integer-valued local maxima
chosen independently inside a +/-2-block window. Build11 improves those observations
before granting the model more freedom.

For every populated 3x3 cell the diagnostic keeps the same bounded integer phase
search, then evaluates a clipped signed-sync correlation surface around the selected
integer maximum. One-dimensional parabolic interpolation along X and Y estimates a
fractional-block offset when the maximum is internal to the search window. Each
control receives an explicit confidence derived from correlation strength, peak
prominence and local curvature. Maxima that land on the +/-2 boundary are flagged and
down-weighted because the true peak may continue outside the bounded search.

The field fitter now combines those confidence weights with a Huber robust reweighting
step. It fits the affine six-parameter model first. With at least eight controls it
also evaluates a ridge-regularized quadratic basis, but the higher-order model is
selected only when leave-one-out RMS improves by both an absolute and relative
margin. The maximum correction remains bounded and the existing four HMAC decode
slots are unchanged.

Synthetic regressions verify fractional parabolic recovery, confidence suppression
of an ambiguous control, Huber rejection of a high-confidence geometric outlier,
search-boundary confidence reduction and cross-validated selection of a truly curved
field. The original exact affine phase-drift regression continues to recover an
ordinary authenticated Format-v3 payload.

The six real acquisitions give a useful negative result. After boundary-truncated
controls are down-weighted, the build10 `foto bici dritta` field that previously
changed 6->5 in an HMAC slot now changes 6->8 and is rejected before decode. No real
case consumes a smooth-field HMAC slot in final build11. `foto bici storta` still
shows a promising same-candidate 4->3 diagnostic correction, but its leave-one-out
RMS is about 2.18 blocks and fails the strict decode gate. A scanner `0270_002`
candidate is the only current real case where the regularized quadratic model wins
the model-selection comparison, but the predicted correction exceeds the bounded
maximum and is rejected.

This checkpoint therefore weakens the evidence that a globally fitted field is ready
for decoding and strengthens the conclusion that control-point quality/boundary
ambiguity matters materially. Future work should prefer measurements that constrain
local phase continuously and independently of the known header before adding still
more model degrees of freedom.


## v0.3.0-build12 — key-independent structural phase controls

Build11 showed that the quality of locally selected, key-assisted phase controls was
the limiting factor for smooth-field fitting. Build12 therefore changes the observer
rather than increasing polynomial order. The primary blind observer uses Format-v3
structural repetition in the robust and balanced mappings: repeated positions should
produce coherent observed DCT margins at the correct tile phase without knowing the
key or the repeated bit value. Capacity has no repetition and deliberately contributes
no such evidence.

The aggregate grid is searched over one complete 35x32 phase for robust/balanced,
after which each 3x3 spatial cell receives only the existing +/-2-block local search.
Sub-block interpolation, boundary penalties and confidence weighting remain bounded.
A direct pairwise cross-cell correlation observer is retained only as a fallback; on
the real corpus its correlation peaks are too weak to be the primary source. All blind
searches operate on grids already collected for the four retained decode candidates
and therefore add no source-image DCT sampling.

The key-assisted build11 local phase remains available only to measure blind-oracle
distance after both have been estimated independently. The projective smooth fitter
now takes blind controls exclusively. A controlled synthetic test verifies exact blind
structural phase recovery on a clean repeated tile, and a separate blind-control field
regression verifies that known affine blind controls can drive the existing smooth
correction to an ordinary authenticated Format-v3 HMAC. These component tests do not
claim that the complete blind observer already reconstructs a real print-camera warp.

Across the six private real acquisitions only `foto stampa storta.jpg` passes the
blind diagnostic resample gate. The best same-candidate corrected prefix changes
5->4 post-ECC known-header errors with fit RMS about 1.27 blocks and leave-one-out RMS
about 2.04 blocks. The strict 2.0-block LOO decode gate rejects it, so no blind field
consumes an HMAC slot and no real image authenticates. The best frontal hard case
remains 2/24. Build12 therefore establishes independent structural registration
evidence while showing that blind control quality/cycle ambiguity remains insufficient
for promotion.


## v0.3.0-build13 — guided dual-observer blind phase

Build13 tested a second key-independent observation for the build12 blind controls.
Rather than repeating the weak unrestricted +/-4 cross-cell search, it uses repetition
coherence to predict each pair's relative shift and evaluates only a +/-1 image-domain
correlation neighbourhood. Strong secondary evidence can confirm controls or resolve a
single integer-block cycle slip, but only after an independent absolute-strength gate.

The controlled synthetic regression gives a strong guided pair score (about 0.61) and
exercises the consensus path. Across the six private real acquisitions, guided pair
scores remain only about 0.053--0.065. Although the resulting secondary offsets are
closer to the primary than the old unrestricted observer (roughly 0.51--0.85 blocks
mean distance), none clear the 0.10 strength gate. Build13 therefore applies no real
cycle-slip correction and makes no new HMAC attempt. The negative result is useful: it
rules out direct cross-cell DCT-pattern correlation, even when repetition-guided, as a
sufficiently strong second measurement on the current physical corpus.

The build also adds `make test-list` to make the growing test harness discoverable.


## v0.3.0-build14 — local fractional lattice phase and bounded unwrap

Build13 demonstrated that direct cross-cell DCT-pattern correlation is too weak on the
physical corpus to resolve repetition-observer cycle ambiguity. Build14 therefore
changes domains. It reuses the key-independent local lattice estimator that already
runs before projective authentication and asks a narrower question: after projecting a
measured regional lattice intersection into a candidate's canonical plane, what is its
position **modulo one 8-pixel Format-v3 block**?

That measurement is deliberately not described as an absolute local phase. A periodic
lattice cannot distinguish one integer block cycle from the next. The observer therefore
supplies only the fractional component. Repetition coherence remains the initial cycle
index. Fractional agreement may refine a control, while a low-confidence repetition
cycle may move by at most +/-1 only when a robust affine prediction from the other cells
reduces residual error by at least 0.45 block and leaves at most 0.55 block residual.
Two unwrap passes are the hard ceiling. None of these decisions uses the key or known
Format-v3 prefix.

Several implementation variants were kept as research history rather than erased:

1. The first implementation interpreted the canonical lattice residual as the inverse
   sampler correction. On `foto stampa storta.jpg` it improved geometric fit but left
   the measured same-candidate prefix at 5 -> 5. Reviewing the internal smooth-phase
   convention showed that controls represent **observed phase drift**, while the fitter
   applies the inverse; the sign convention was corrected accordingly.
2. With the corrected observed-drift convention but integer unwrap disabled, the
   modulo-only field was not useful: the same storta candidate became 5 -> 6 and a bike
   candidate 7 -> 11. This confirms the theoretical limitation that fractional phase
   alone cannot resolve integer cycle ambiguity.
3. With the final bounded unwrap enabled, the storta candidate becomes 5 -> 2, the
   straight-bike candidate 7 -> 3, and scanner 001 candidate 9 -> 4. These are
   same-candidate research comparisons, not new global hard baselines. Scanner 002 also
   demonstrates the opposite case (7 -> 10), so the method is not generally reliable.

The lattice observer itself is substantially stronger than build13 cross-cell evidence,
with mean regional confidence around 0.38--0.52 on the current real best-hard spatial
cases. Blind-vs-key-assisted-oracle distance improves modestly in several acquisitions,
but not uniformly. Crucially, all apparently strong same-candidate bit improvements
remain below the strict smooth decode confidence gate and/or outside the strict LOO
gate. No build14 smooth field consumes an HMAC slot and no real acquisition authenticates.

This checkpoint establishes that independent lattice-phase information is useful, but
also that **integer unwrapping is now the primary unresolved problem**. Future work
should validate cycle choices using a third independent smoothness/gradient/spectral
constraint or a global discrete unwrap objective before any promotion. Lowering the
0.20 confidence or 2.0-block LOO decode gates merely to admit the current cases is
explicitly rejected.



## v0.3.0-build15 — global discrete phase unwrap

Build15 replaces build14's local affine-guided integer-cycle decisions with a bounded
global discrete optimization over the 3x3 blind control field. Each axis is solved
independently with a beam search over only -1/0/+1 cycle shifts. Candidate states are
scored from key-independent geometry: robust affine residual, confidence-weighted cycle
change cost and spatial second-difference curvature. An explicit first-vs-second solution
margin is required before any axis is accepted.

A critical methodological change is transactional fusion. Build14 showed that fractional
lattice phase without a trustworthy integer unwrap can make the channel worse, so build15
never leaves fractional corrections behind after a rejected unwrap. The lattice proposal
is built, globally evaluated, then either committed as a whole or rolled back to the
pre-lattice repetition controls.

Controlled regressions recover known synthetic cycle slips and restore an ordinary
Format-v3 HMAC, while a symmetric ambiguous layout is rejected and leaves the original
controls byte-for-byte equivalent in value. On representative real acquisitions the
global solver frequently finds a much lower geometric objective but with a small
first/second margin; those solutions are intentionally rejected. In particular, the
large build14 `foto stampa storta.jpg` 5 -> 2 diagnostic cannot yet be justified by an
independent unique global unwrap and is not promoted.

## v0.3.0-build19 — multi-partition integer-cycle stability

Build19 tests the remaining build18 hypothesis directly: perhaps the A/B split was merely
unlucky, and a persistent integer-cycle field would emerge if the repetition constraints
were partitioned in several independent deterministic ways. Eight coded-bit-group
partitions are fixed in advance and evaluated in both directions. The experiment remains
best-candidate-only, ambiguous-only and diagnostic-only.

The physical result is strongly negative for repetition-only cycle identification. Every
ambiguous acquisition yields all 16 trials, yet all 16 complete local cycle fields are
different. Filtering to held-out-supported trials does not create a stable field: the
inclined photo has 9 supported fields and 9 unique fields, scanner 001 has 6/6 unique,
and scanner 002 has 12/12 unique. Per-cell modal fractions remain low and pairwise field
agreement is about 7--8%. No vote is promoted.

This checkpoint closes the simple repartition/voting branch on the current four-case
physical corpus. Further decoder research should introduce a genuinely independent
integer-cycle anchor or a stronger physical model, not lower ambiguity thresholds.

## v0.3.0-build18 — lazy cross-fit and cell-level instability map

Build18 does not introduce another cycle solver. It preserves build17's all-pairs exact
unwrap and held-out criterion, but defers the expensive research-only cross-fit until the
bit-channel ranking has identified the single best candidate. Lower-ranked candidates
retain all earlier diagnostics but cannot consume duplicate A/B cross-fit searches.

The checkpoint also records the local cycle proposed by each independent fold in every
3x3 control cell together with primary/lattice confidence. This directly tests whether
the small A/B agreement observed in build17 is concentrated in unusually strong cells.
It is not: on the inclined smartphone acquisition the two agreeing cells have lower mean
joint fold confidence than the seven disagreeing cells. Scanner 001 and 002 have zero
agreement. No partial field is promoted and the production/HMAC path remains unchanged.

## v0.3.0-build17 — held-out repetition cross-fit

Build17 addresses the limitation left explicitly open by build16: its two
`split-repetition-top2` folds were evaluated after the primary controls had already
consumed the complete repetition-pair set. Build17 partitions repetition evidence by
coded-bit group before control estimation. Fold A constructs a proposal that fold B
alone validates, then the roles are reversed. The normal all-pairs exact solver remains
the authoritative conservative gate; cross-fit is research telemetry and cannot create
an HMAC candidate.

The four available real acquisitions produce three distinct negative/diagnostic
patterns. The inclined smartphone case receives positive held-out support in both
directions but the proposed local integer fields disagree strongly (2/9 agreement).
Scanner 001 receives negative held-out support in both directions, and scanner 002
splits directionally. No case meets the deliberately strict condition of symmetric
held-out support plus full proposal agreement. The result narrows the remaining problem
from "find independent evidence" to "obtain a stable cycle proposal across independent
evidence partitions".

## v0.3.0-build16 — exact global top-2 certification

Build16 closes two correctness/reporting issues found in the build15 audit. The global
axis solver no longer uses a width-64 beam: every bounded `-1/0/+1` assignment over the
fixed 3x3 control grid is enumerated, giving a true top-1/top-2 within at most 19683
states per axis. A deterministic regression demonstrates a field that the old beam
would accept but the exact margin correctly rejects as ambiguous. No-eligible paths
are now explicitly not-applicable and remain JSON-finite.

Reporting now separates X/Y eligibility and status and records whether a second
solution exists. Mixed accepted/ambiguous axes are handled atomically: any ambiguous
axis prevents every integer-cycle change from being applied. A diagnostic
split-repetition top-2 comparison partitions the
key-independent repeated-position evidence into two deterministic folds. It is not an
acceptance gate because the primary repetition observer used the complete pair set.

The source package keeps the intentional corpus test key `Piccotti` as the default
`PRINT_CAMERA_KEY` for reproducible physical tests, while protecting both private
acquisition directories in `.gitignore`. It also documents canonical `.jpg` names for
the two historically mis-suffixed smartphone JPEGs and provides a local SHA-256
corpus-manifest target. The key is a public test constant; the acquisition files remain
private research material.

On the four physical acquisitions available during development, exact ranking preserves
the build15 conservative conclusions. Scanner 0270_002's first/second margin tightens
from the beam value around 0.02634 to the exact value around 0.01780; the proposal
remains ambiguous and no smooth HMAC slot is consumed. Full six-image validation still
requires the two bicycle photographs.

## v0.3.0-build20 — independent pairwise cycle-anchor experiment

Build20 tested whether the pre-existing unguided cross-cell image-domain registration graph
could break the integer-cycle ambiguity independently of repetition evidence. The comparison
was made gauge-invariant and diagnostic-only against the exact all-pairs top-1 and runner-up.
All three ambiguous real acquisitions produced a continuous preference for top-1, but none
produced even one rounded-cycle agreement cell after gauge alignment (0/9 for both candidates).
Scanner 001 contradicted held-out repetition, while negative-control scanner 002 also preferred
top-1. The observer is therefore retained only as relative-shape evidence; it is not promoted
as a cycle anchor and no decoder/HMAC budget changed.


## v0.3.0-build21 — Format-v3 observability audit

Build21 stopped proposing another cycle solver and instead audited the information already
present in the frozen public format. For each concrete profile it compared the repetition-pair
graph under every unit block shift and passed deterministic valid Hamming codewords through the
same wrong-origin tile mapping used by decoding.

The audit shows that robust/balanced repetition topology is technically absolute but weak: the
closest one-block shift preserves roughly 94% of expected equal-code-index pairs. Capacity has
no repetition topology. Hamming parity supplies strong horizontal/diagonal discrimination, but
vertical unit shifts expose aliases: robust is a zero-syndrome synthetic alias and capacity is
an exact whole-word Hamming permutation alias; balanced is only weakly separated vertically.
This explains why increasingly sophisticated unwrap solvers cannot manufacture a uniformly
observable absolute cycle from the current signals. Build21 is diagnostic-only and leaves the
format and production extractor untouched.


## v0.3.0-build22 — physical topology observability + quantified v4 branch

Build21 established that robust/balanced repetition topology is technically non-invariant but
leaves only a 5–8% unit-shift gap. Build22 tests whether that residual survives the real image
channel without reusing the evidence that chooses phase. For each unit shift, pair edges shared
by both hypotheses register a symmetric phase pair; topology-exclusive edges are held out and
select between the two phases.

The physical signal exists but does not become a reliable absolute reference. The inclined
smartphone case has only 0.3125 mean winner modal-phase fraction, while scanner 001 and scanner
002 both reach 0.4375. Directional cell consistency is similarly overlapping (0.6944, 0.7292,
0.6597). Scanner 002 therefore remains a decisive negative control: the new topology metric can
look at least as stable on a failing acquisition as on the useful one. No decoder gate is added.

Because the remaining v3 structural anchor is now both theoretically weak and empirically
non-discriminative, build22 opens a formal Format-v4 design study. The preferred first prototype
uses a 37×32 tile and reserves 64 blocks for a public absolute pilot. This raises tile area by
5.7% and minimum carrier width from 280 to 296 pixels, while keeping exactly 1120 data positions.
Under v3-like framing/ECC arithmetic this preserves the 16/32/64-byte ceilings and profile
redundancy. Format v4 is not implemented in build22.

## v0.3.0-build25 — geometric pilot qualification

Build25 keeps Format v3 frozen and tests whether `prototype-2-search-p64` survives geometric resampling once an independent mapping is available. The new projective pilot sampler reads canonical 8x8 carrier blocks through a supplied homography, so rotation, anisotropic scale, shear and perspective do not need to be rectified into a temporary image first. This remains geometry evidence only; it has no payload or authentication role.

A fixed 16-case matrix covers two arbitrary rotations, two anisotropic scales, X/Y shear, two perspective shapes, rotation+scale, rotation+shear, rotated 75% resize, perspective+JPEG/blur/noise and perspective+crop with/without JPEG. All synthetic cases and all 32 local-image cases recover cyclic origin `(0,0)`. The weakest local positive runner-up margin is 0.322069; the corresponding negative-control margin is 0.063984.

This closes the "does the pilot survive known geometric distortion?" question for the current development matrix. It does **not** close blind geometry recovery and does not freeze prototype-2. The next branch must use the pilot itself to score bounded geometry hypotheses without payload/header/HMAC oracles.

## v0.3.0-build24 — reproducible v4 pilot qualification begins

Build24 leaves the frozen Format-v3 core unchanged and turns the Build23 pilot idea into a reproducible
qualification experiment. The search is deliberately bounded and replayable: a fixed-seed 200,000-candidate
joint coordinate/sign stage explores one pilot per 8x8 spatial stratum, then a second fixed-seed 200,000-candidate
balanced-sign refinement operates only on the winning mask. Every retained candidate has a SHA-256 identity.

The resulting non-normative `prototype-2-search-p64` improves prototype-1's exhaustive 37x32 toroidal maxima
from 8 to 7 overlapping pilot positions and from 5 to 4 absolute signed correlation, with zero exact non-zero
cyclic aliases. Worst contiguous-crop margins for 64/48/32/24/16 visible pilots are 60/43/28/20/12 symbols;
256 deterministic random subsets at each visibility level produce no false-origin tie or win.

Build24 also adds the first image-domain pilot-only channel. A synthetic qualification carrier writes the
candidate pilot and deterministic pseudo-random data-plane signs into the existing DCT coefficient pair, but
it does not contain a payload, framing, ECC or HMAC. The detector assumes an already-resolved 8/6/4-pixel
lattice and exhaustively scores all 1184 cyclic origins. Deterministic JPEG, blur, noise, gamma, resize and
aligned-crop regressions recover the correct origin. The two local original images supplied for development
also recover the correct origin across the same initial transform set; their unmarked forms remain controls.

This is not a Format-v4 freeze. General rotation/affine/perspective geometry, combined degradation, framing,
ECC selection and real print-camera/scanner v4 evidence are still open. Production v3 APIs and CLI behavior
remain unchanged.

## v0.3.0-build23 — v3 research closure + v4 experimental foundation

Build23 records the project boundary reached after the qualified build22 experiment. Format v3 remains
the only implemented/interoperable format and keeps its encoder, production decoder, profiles, whitening,
Hamming mapping, HMAC semantics and golden fingerprints frozen. The absolute-cycle research line is closed
because the remaining public topology asymmetry is not discriminative against scanner 002; additional
thresholds, repartitioning or voting would reuse the same insufficient evidence.

The build adds a dedicated v3 closure document and reclassifies future v3 work as maintenance/regression
unless a genuinely independent source of information appears. It then starts a separate v4 prototype layer.
No v4 carrier can yet be emitted or decoded by the CLI. The provisional v4 geometry is 37×32 blocks with
64 public pilot positions and 1120 data positions. Prototype-1 is spatially stratified and sign-balanced;
its exhaustive toroidal audit finds no perfect non-zero cyclic alias, maximum shifted-mask overlap 8/64 and
maximum wrong-shift signed correlation 5/64. These numbers are frozen only as build23 prototype regressions,
not as a final Format-v4 interoperability contract.

## v0.3.0-build26 — bounded blind v4 geometry proposal + pilot validation

Build26 removes the independently supplied homography used by Build25 for a bounded auto-framed research case. The new coarse stage uses only self-consistency of repeated v4 **data-plane coordinates**: corresponding DCT observations in adjacent 37x32 tile repetitions should agree when the candidate geometry is correct. The actual data symbols are never known or scored, pilot signs are excluded from this stage, and no payload/key/header/ECC/HMAC evidence is available.

A deterministic rotation/affine/shear/mild-perspective family is searched with canvas dimensions as a weak public prior. The repeat-consistency stage keeps a small geometry bank. Only then does `prototype-2-search-p64` rank sub-block translations on central tile rows; final selection is performed again on held-out corner pilot repetitions. This preserves proposal/validation separation instead of turning pilot correlation into an authentication oracle.

The synthetic matrix recovers the declared geometry families within the Build26 corner-error bound, including 9.3 degree + 1.08/0.92 anisotropic scale + 0.030/0.015 perspective. Both local originals also recover the representative rotate-scale and combined-perspective cases with correct origin. The search remains experimental and assumes known canonical extent plus auto-framed transforms; arbitrary crop/translation and the real print-camera problem remain future work.



## v0.3.0-build27 — unknown crop/translation placement qualification

Build27 deliberately does not combine every remaining unknown at once. Following the same staged method used by Build25/26, it freezes the geometry input for one experiment and asks whether the public v4 pilot can recover where the transformed carrier lies after arbitrary crop or padded-canvas placement.

The search receives the canonical-to-full-frame homography and full transformed dimensions but no crop offset or canvas offset. Translation bounds follow from full-vs-observed extent and remain explicitly capped. Even-index pilot symbols propose placement on a coarse 4-pixel grid and integer refinement; odd-index symbols are held out for final validation. No payload, header, ECC, key or HMAC evidence participates.

Synthetic qualification recovers the declared offsets for affine crop, combined-perspective crop, a deeper crop and two padded-canvas placements. The same representative crop/pad cases pass on `PJ_lingua.PNG` and `PJ_piccolo.png`. A one-degree wrong-geometry control drops validation from roughly 0.994 to 0.236, so placement search is not allowed to hide geometric error.

Because v4 repeats tiles, translations separated by a complete tile can be observationally equivalent. Build27 therefore validates the recovered mapping by full-pilot canonical sampling instead of requiring a non-zero placement top-2 gap. The remaining research problem is **joint** crop-tolerant geometry/extent estimation plus placement; Build27 does not yet freeze the pilot or enable a v4 codec.


## v0.3.0-build28 — first joint blind affine+crop recovery

Build28 combines the previously separated Build26 geometry and Build27 placement problems for one deliberately bounded family: rotation, anisotropic scale and arbitrary **negative crop translation** are all unknown. Early prototypes tried to use repeated-tile consistency or sparse DCT energy as a single global selector. Those variants either overfit photographic texture or pruned the correct basin before the pilot could validate it; lowering thresholds was rejected.

The retained path uses sign-independent absolute DCT carrier energy and 8-pixel phase contrast as the geometry observable. A hierarchical bank searches angle and X/Y scale, then 64- and 256-block refinements localize the strongest basins. The key correction was a compact **coupled** angle/scale search on only the last two basins: coordinate descent could leave `PJ_piccolo.png` at approximately `scaleY=0.915` although the correct `~0.93` basin was nearby. The coupled refinement raises the correct photographic basin to first place without using pilot symbols.

After geometry is fixed, Build27 placement recovery is reused unchanged: even pilot indices propose translation, odd indices validate it, and the full pilot reports cyclic origin. Synthetic and both local-original gates recover `(0,0)` with bounded hypothesis counts and strong matched-negative separation. Geometry selection therefore remains structurally independent from the public pilot, while placement and origin use the pilot only after the affine mapping has been chosen.

Build28 is not the final blind v4 decoder. Joint shear, projective/perspective geometry and positive padded-canvas placement are still separate work, as are unknown physical carrier extent, v4 framing/ECC, the payload encoder and fresh print-camera/scanner evidence.

The project also records an architectural decision that had previously been implicit: Go is used in part to keep one reusable implementation compilable across Linux, Windows, Android and iOS targets. A future graphical frontend is planned, especially for mobile use, while the qualified Go core remains independent of UI technology.

## v0.3.0-build29 — bounded joint projective/padded search with safe rejection

Build29 attacks the two composition problems left open by Build28: projective distortion plus unknown crop, and affine geometry plus positive padded placement. Several pilot-first and sparse-repeat prototypes were rejected because they either multiplied geometry by a full placement scan or overfit a single tile. A concrete indexing bug in an early single-tile probe (`y` derived with width 32 instead of 37) was found and removed during the investigation; none of those diagnostic probes enter the final source tree.

The retained projective search keeps a bounded public structural bank, requires two pilot halves to agree on the same local phase/origin hypothesis, preserves geometric basin diversity, refines only a small set and then lets the qualified Build27 placement search evaluate the remaining candidates. Complete public-pilot margin and absolute origin are final development gates. This gives a synthetic projective+crop PASS and a photographic `PJ_lingua` PASS, while `PJ_piccolo` is deliberately rejected because its margin is only about 0.119.

The padded affine branch passes the synthetic case exactly. Both photographic padded joint searches are currently rejected. A known-geometry `PJ_lingua` control reaches roughly 0.991 placement validation and 0.554 pilot margin, showing that the pilot channel remains strong and the open issue is geometry ranking. Build29 therefore treats safe rejection as a qualified outcome and explicitly refuses to lower thresholds to manufacture coverage. Prototype-2 remains non-normative and Format v3 remains frozen.



## v0.3.0-build30 — pilot candidate lock and freeze-readiness audit

Build30 deliberately does not widen the Build29 acceptance envelope. The first experiment tested whether padded photographic ranking could be repaired simply by making the DCT phase bank fully translation-invariant over all 64 sub-block phases. The true odd phase becomes observable, but natural photographic texture still produces much stronger false structural maxima at practical low-support budgets. That variant is rejected rather than added as an expensive heuristic. Build29 ACCEPT/SAFE-REJECT gates remain unchanged.

The build instead isolates the question required before a pilot freeze decision: **does the exact prototype-2 pattern itself remain strong when geometry and placement are supplied correctly, including the cases that the joint search rejects?** A new known-mapping corpus audit applies the exact Build29 projective-crop and affine-padded mappings directly to both development originals. All four marked cases recover origin `(0,0)` with score 0.928412–0.988352 and margin 0.484086–0.559445. Matching unmarked controls have margin at most 0.010553. The smaller padded case is especially informative: the Build29 placement/joint path can select the wrong placement, but sampling at the true mapping still gives score 0.932220 and margin 0.484086.

This evidence supports locking `prototype-2-search-p64` against accidental mutation. Build30 therefore adds a semantic identity check for SHA-256 `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174` and documents the exact lock contract. The decision is intentionally narrower than a normative Format-v4 freeze. There is still no physical v4 print-camera/scanner corpus generated by a real v4 encoder, and framing/ECC/payload/HMAC choices are still open. The locked identity may be promoted unchanged after physical qualification; if physical evidence exposes a pilot-level weakness, a replacement must receive a new identity rather than editing the Build30 lock in place.

## v0.3.0-build31 — first real authenticated v4 frame and encoder

Build31 changes the v4 branch from pilot-only geometry research into a real, still experimental data format. The frozen v3 core remains untouched and the stable `embed` / `extract` commands still mean Format v3. New `v4-embed` and `v4-extract` commands are explicit opt-ins so no existing carrier or script silently changes format.

The Build30 pilot identity remains unchanged. The 64 prototype-2 positions are reserved first; the other 1120 row-major tile positions become data ordinals. Data ordinal `d` maps through `(d*251) mod codedBits`, preserving the existing robust/balanced/capacity repetition arithmetic without colliding with a pilot block. During implementation, an important historical-helper trap was found: the old `experimentalV4PrototypeDataPositions()` still describes prototype-1. Build31 deliberately leaves that historical helper intact and adds a prototype-2-locked partition instead, with a regression that covers all 1184 positions exactly once as pilot or data.

The first v4 frame retains 32/48/80-byte profile sizes and Hamming(7,4) only as a controlled baseline. Its version/profile bytes are `0x41`, `0x42` and `0x43`; whitening and HMAC are domain-separated from v3 using `pixseal-whiten-v4` and `pixseal-frame-v4`. HMAC authenticates the header through actual payload bytes and is truncated to eight bytes as in the current comparison baseline. CRC is retained as a post-authentication consistency check, never as the success criterion.

A native-lattice aligned decoder first resolves cyclic tile origin from the public pilot, then reads only non-pilot data residues, tries the three protected frame lengths, dewhitens and accepts only a valid v4 HMAC. It deliberately does not prepend the Build29 geometry search yet. This preserves the evidence ordering while framing/authentication is qualified separately.

The digital channel passes all three profiles on both private development originals. Robust also authenticates after JPEG quality 82 and an aligned crop, and a single 296x256 tile authenticates a capacity payload in the synthetic gate. Build31 therefore supplies the missing encoder needed for a separately printed v4 physical corpus. The pilot remains development-locked, not normative, until that physical corpus yields HMAC-authenticated recovery.


## v0.3.0-build32 — explicit anonymous corpus reset

Build32 deliberately changes test governance rather than the watermark algorithm. The active private corpus becomes the three anonymous LQ/MQ/HQ originals supplied for the new v4 work session. Corpus membership is now explicit in `private-corpus-active.tsv`; names, dimensions and SHA-256 are verified before qualification, so dropping another image into `original pics/` can no longer silently change results.

The public v4 development key for new tests/physical fixtures is `PixSeal-v4-TestKey-2026`. Historical `Piccotti` vectors are retained only where needed to reproduce older results. The locked `prototype-2-search-p64`, Build31 frame/HMAC/whitening/data mapping and all frozen v3 files remain unchanged. LQ established a new common baseline boundary: capacity authenticates at 55% resize but not universally at 50%, so `deep-test` now requires 55% while `extreme-test` retains the full limit map.

## v0.3.0-build33 — MQ projective ranking observability

Build33 isolates the active MQ joint-projective SAFE REJECT instead of widening the search. Under the exact Build29 projective+crop transform it records the known truth-basin rank after the structural bank, half-pilot prefilter and full-proposal stages for both a deterministic random data plane and a real authenticated v4 frame. The correct basin is already structural rank 1/5000 in both cases but degrades to 98/5000 -> 22/2000 on the random plane and 35/5000 -> 2/2000 on the authenticated frame. This proves the blocker is preservation/refinement of an existing basin rather than insufficient global search. Exploratory ranking variants that solved only one data plane are explicitly rejected.

## v0.3.0-build34 — authenticated MQ blind projective recovery

Build34 closes the Build33 MQ blocker without changing the locked pilot, Build31 frame, Hamming baseline or Build29 acceptance floors. The key defect was phase granularity: the former centered proposal sampled 0/+/-4-pixel phases and could skip a strong +/-2-pixel alignment. For carriers with at least 4x4 complete v4 tiles, Build34 keeps two distinct structural anchors, explores bounded coupled local neighborhoods, retains the strongest structural candidates, and ranks them with the public pilot over interior tile rows using 0/+/-2/+/-4 pixel phases. First/last tile rows stay held out until a geometry is committed.

Non-zero cyclic origin is canonicalized by composing the block-cycle translation into the homography and re-running the detector; the unchanged `(0,0)` gate still decides acceptance. Only then are the locked 1120 data positions sampled and decoded through Hamming, whitening, v4 frame parsing, CRC and HMAC. On the active MQ carrier the exact projective family is recovered with approximately 0.001 normalized corner error, and two different authenticated robust payloads recover exactly. LQ remains a SAFE REJECT on the older bounded path. Build35 preserves this checkpoint under its correct historical identity: Build34.

## v0.3.0-build35 — physical-qualification handoff

Build35 deliberately does not retune the Build34 decoder. It first verifies Build34 numbering consistently across source symbols, filenames, Make targets and documentation, then exposes the qualified blind projective+crop path through the public experimental `ExperimentalV4ExtractProjective` API and `v4-extract-projective` CLI. The caller supplies the block-aligned dimensions of the digital carrier before printing; public structure/pilot evidence must pass the existing geometry gate before frame/HMAC processing begins.

The build also replaces the old generic physical-fixture helper with a scanner-first MQ campaign. `make v4-physical-fixtures` creates an unmarked control and two robust/strength-24 carriers carrying `v4-b35-phys-a` and `v4-b35-phys-b`, plus SHA-256 and acquisition manifests. The first physical contract is actual-size 300-ppi printing, 300-dpi lossless scanning, tight crop to artwork edges and no manual geometric/photo correction. `make v4-physical-qualification` requires exact HMAC-authenticated recovery of both marked payloads and rejection of the control. Smartphone captures are a separate later gate.
## v0.3.0-build36 — physical-channel soft-decision closure

The first scanner-first Build35 acquisition campaign produced three real paper scans. The available office scanner stores color scans as JPEG; the highest-quality useful mode was 600 dpi. Build35's blind projective search correctly continues to reject the control, but it also SAFE-REJECTs the two marked acquisitions because weak physical pilot evidence allows false geometric basins to outrank the true scanner transform.

A private reference-assisted diagnostic was used to separate geometry from channel capacity. The original digital carrier supplied only the geometric registration for this diagnostic; it was not used by the production decoder and the result is therefore explicitly non-blind. Once the correct scan homography is fixed, the current Format-v4 signal is strong enough to recover both physical robust/strength-24 payloads. Hard Hamming remains just beyond its correction radius on some words, while deterministic soft maximum-likelihood Hamming recovery yields valid HMAC for both `v4-b35-phys-a` and `v4-b35-phys-b`.

Build36 therefore changes only post-geometry reliability handling. Signed DCT margins are retained through the projective data sampler and decoded with soft Hamming before the historical hard-decision fallback. Geometry selection, Build29 pilot acceptance, Build30 pilot identity, Build31 wire format/HMAC, encoder strength and all frozen Format-v3 files remain unchanged. The remaining physical blocker is now specifically blind scanner registration; no new printing is required for the next research step.


## v0.3.0-build37 — blind scanner registration and first physical HMAC PASS

Build37 closes the blocker isolated by Build36. The original Build35 paper set is reused unchanged: one unmarked MQ control and two robust/strength-24 v4 carriers printed at actual size. The available office scanner produced 600-dpi color JPEG full-page captures rather than the planned 300-dpi lossless files.

The scanner path does not use the digital original as a registration template. It estimates the visible artwork rectangle against the white page, converts that quadrilateral into a canonical-to-scan seed homography and explores only a tightly bounded scanner-affine correction. Even-indexed public pilot symbols propose candidates; odd-indexed pilot symbols rank the frozen proposal shortlist. The held-out winner must also establish complete-pilot cyclic origin `(0,0)` with explicit score/margin floors. Five independently pilot-qualified geometries are retained as an ensemble. Only then are protected data margins sampled and validation-weighted across the ensemble, followed by Build36 soft Hamming, v4 frame parsing and HMAC.

On the private physical corpus the unmarked control fails the held-out/full-pilot gate, while both marked captures authenticate their exact messages: `v4-b35-phys-a` and `v4-b35-phys-b`. This is the first fully blind print -> paper -> scanner -> JPEG -> HMAC PASS for Format v4. It does not promote the experimental v4 format to normative status and does not claim smartphone-camera robustness. The next physical campaign is frontal then oblique smartphone capture of the same printed sheets.
## v0.3.0-build38 — smartphone channel characterization and strength-48 qualification pack

After Build37 closed the scanner path, nine original 200 MP smartphone photographs were acquired from the same three Build35 paper sheets: control, marked A and marked B, each under frontal, mild-perspective and stronger-perspective conditions. The first result is intentionally recorded as channel evidence rather than hidden by threshold retuning. Build37's scanner-specific boundary detector often locks onto the page edge instead of the artwork in camera frames, but a stronger diagnostic goes further: even when the digital original is allowed to supply a reference-assisted homography, and even after dense non-rigid residual registration, representative strength-24 photographs remain far outside the current Hamming/HMAC correction envelope. Measured examples leave roughly 14.5–17.9% protected coded-bit errors and 20–42 wrong bits out of 256 after soft Hamming.

Pilot-only per-tile translation searches can manufacture high local pilot correlations by fitting natural photo texture, while the protected data bits do not improve correspondingly. Build38 therefore does not promote such a warp to production and does not weaken pilot/HMAC gates. The physical conclusion is that the scanner-qualified strength-24 carrier is too weak for reliable single-shot smartphone recovery under this print/camera pipeline.

Build38 keeps the entire v4 wire format unchanged and changes only the next physical experiment: robust phone fixtures use strength 48. On the canonical MQ carrier this gives about 32.5 dB PSNR relative to the unmarked control and remains visually subtle in the qualification image. `make v4-phone-fixtures` generates a private control plus `v4-b38-phone-a` and `v4-b38-phone-b`, all at the same 1632x1632 / 300-ppi print geometry used by the scanner campaign. The previous strength-24 phone photographs remain a historical private corpus.

The next decoder step is conditional on the new strength-48 photographs: first establish that the protected data channel now reaches the Hamming/HMAC envelope under reference-assisted analysis, then build a blind camera-specific boundary/projective/residual model without allowing payload/HMAC evidence to select geometry.



## v0.3.0-build40 — pilot-only smartphone residual-warp checkpoint

Build40 tests whether the remaining strength-48 phone failures can be explained by a small smooth deformation after Build39's global projective registration. It does not alter the printed carrier. A quadratic canonical displacement field is capped at six pixels, controls are proposed only from checkerboard-A public-pilot tiles, and checkerboard-B tiles remain held out until acceptance. Ambiguous local pilot peaks and incoherent fits are rejected.

The synthetic gate confirms the mechanism is useful when its assumptions are true: a smooth non-projective deformation is recovered, held-out pilot evidence improves, protected margins sampled through the fitted field soft-Hamming decode and authenticate the exact v4 frame, while an unmarked control does not qualify.

Representative private real-phone captures provide the equally important negative half of the experiment. Starting from the current Build39 basins, locally attractive pilot corrections are not spatially coherent: they can raise proposal evidence while reducing held-out validation and leaving a non-zero complete-pilot origin. Build40 therefore refuses to enlarge the residual bound or use HMAC as a search oracle. The next phone build must improve global projective-basin preservation/ranking first; only then is another small residual pass justified.

## v0.3.0-build39 — blind smartphone projective-registration checkpoint

The completed strength-48 phone corpus changes the diagnosis from Build38. Reference-assisted registration can now recover authenticated frames from multiple real marked photographs (including strong-perspective cases), and another marked capture is within one post-soft-Hamming bit. Strength 48 is therefore retained; no Format-v4 framing, ECC or HMAC change is justified.

Build39 introduces the first dedicated blind phone path. Native ~200 MP captures are reduced internally to a bounded working image, the printed artwork is localized independently of the watermark, and the resulting quadrilateral seeds a projective search. Public pilot evidence is split spatially between proposal and held-out validation, preserving the rule that payload bytes, ECC outcomes and HMAC cannot act as a geometry oracle. The Build36 soft-decision channel is reused only after geometry qualification.

The checkpoint also exposes the next physical limitation: a global homography is not sufficient for every phone capture. Camera lens distortion, page/print flatness and resampling leave local residual phase error even when the artwork corners are close. Build39 therefore does not claim complete blind phone recovery; subsequent work should add a tightly bounded pilot-only residual field after projective registration, with held-out spatial validation and unchanged HMAC semantics.


## v0.3.0-build41 — blind smartphone basin milestone

Build40 established that a bounded local residual field is useful when a correct global phone homography already exists, but the complete nine-photo strength-48 matrix showed that the real blind failures were usually earlier: the decoder was losing or mis-ranking the global projective basin. Build41 therefore leaves the encoder, strength 48, locked public pilot, protected data mapping, Hamming(7,4), whitening/HMAC domains, residual model and frozen v3 core untouched and changes only the global smartphone registration front end.

A structure-only inner-artwork boundary anchors physical phase. The public pilot is split deterministically into three spatial folds. Folds 1 and 2 alone generate and refine a bounded eight-coordinate projective shortlist; fold 0 is never sampled until that shortlist has been frozen. Proposal-only cyclic fitting finds the local shape basin without committing to a whole-tile alias, then bounded proposal-only phase restoration returns candidates to physical origin `(0,0)`. The strongest proposal candidates may receive a small proposal-only corner polish before freeze. Held-out and complete-pilot evidence can accept/reject and rank only frozen candidates. Payload/header contents, ECC result, secret key and HMAC never generate or rank geometry.

The final data sampler uses two independently pilot-qualified geometries that must also be separated by a scale-normalized minimum distance. This prevents two near-identical local optima from masquerading as an ensemble and lets the unchanged signed-margin/soft-Hamming path bracket residual sub-pixel uncertainty without HMAC-guided selection. Build40's <=6 px residual remains downstream and is attempted only after the Build41 global ensemble has been frozen.

On the private nine-photo strength-48 corpus Build41 reaches the staged physical milestone defined at the Build40 handoff: all three unmarked controls reject; `phone-a-angle.jpg` authenticates `v4-b38-phone-a`; `phone-b-front.jpg` authenticates `v4-b38-phone-b`. `A/mild` reaches accepted geometry and protected-data decoding but still fails HMAC. `A/front`, `B/mild` and `B/angle` remain outside the accepted basin. The residual layer is not applied in either PASS case because held-out residual evidence does not justify it. Build41 therefore closes the first blind smartphone A/B milestone while explicitly leaving 6/6 marked robustness for later builds.


## v0.3.0-build42 — qualified-bank data-list recovery

Build42 starts from a deliberately narrow observation: `A/mild` already passed the complete Build41 public-pilot geometry gate and reached protected-data decoding, yet failed HMAC. Reopening geometry or increasing strength was therefore not justified. Offline diagnostics over the frozen five-member qualified bank showed that no Build41 pair authenticated, but several deterministic three-member ensembles reduced the residual soft-Hamming error to a single information bit. In the best diagnostic ensemble the correct nibble was exactly the second ML Hamming candidate for one word.

The production change is post-geometry only. Build41 still generates and qualifies geometry from structure/public pilot, and its normal two-member decoder remains first. On failure, Build42 keeps up to six already-qualified geometries, enumerates three-member data ensembles deterministically, averages signed protected margins and applies a bounded list decoder over the second-best nibble of the ten weakest Hamming words. Geometry is already frozen before any of this occurs. HMAC is only the final complete-frame validator.

To avoid the expensive Build40 residual pass when it is unnecessary, the phone order becomes global Build41 decode, then Build42 data-list fallback, then Build40 residual only if both fail. On the unchanged physical corpus A/angle and B/front remain direct passes, A/mild becomes an authenticated pass, and all three controls still reject before data decode. A/front, B/mild and B/angle remain geometry research cases for a later build.


## v0.3.0-build43 — proposal-only side-pair smartphone recovery

Build43 closes the `A/front` blind-registration gap without changing the watermark/data plane. The new fallback executes only after Build41 geometry rejection, proposes side-pair basins from image structure/public-pilot evidence, preserves angular/fine/complement diversity, freezes a bounded bank before held-out evaluation, then reuses Build41 qualification and the Build42 list decoder unchanged. The private Build38 matrix reaches the minimum target: controls 3/3 reject; A/front, A/mild, A/angle and B/front authenticate; B/mild/B/angle remain informational rejects.

Build43 final qualification also identified a toolchain boundary unrelated to the watermark algorithm. Go 1.26 replaced `image/jpeg`; the canonical A/mild JPEG decodes to different Y/Cb/Cr samples than under Go 1.25.1 and can fall outside the qualified geometry path. Rebuilding the same Build43 source with Go 1.25.1 restores the expected A/mild HMAC PASS. Build43 therefore pins Go 1.25.1 for build/test commands and defers deterministic Go 1.26+ JPEG ingest to a later build.

### Build45 — isolate B/mild before changing geometry

After Build44 removed Go-version-dependent JPEG rasterization, the remaining difficult smartphone cases can be studied without toolchain ambiguity. Build45 therefore makes no production algorithm change. It exposes the already existing Build41/43 geometry and Build42 data stages as structured diagnostic telemetry and classifies a failed phone decode according to the stage where the unchanged pipeline stops.

A strictly separate reference-assisted experiment uses SIFT/RANSAC only to supply an external quadrilateral. On the retained private corpus, both B/mild and B/angle authenticate `v4-b38-phone-b` immediately under that supplied mapping. This sharply narrows the research problem: the strength-48 protected channel is viable in both images and the next justified work is blind geometry acquisition/ranking, not stronger embedding, new ECC or HMAC-guided search.

## v0.3.0-build46 — qualified Build43 singleton handoff study

The Build45 qualification-host run refined the failure location for both remaining marked-B cases. B/mild freezes 32 Build43 candidates and B/angle freezes 28; exactly one candidate in each image passes held-out Build43 qualification. Production still reports an empty Build42 bank because the Build43 fallback is not accepted until at least two qualified geometries exist, while the Build42 list-bank path independently requires at least three. The prior Build45 label `qualification` was therefore too coarse to distinguish a true held-out failure from a qualified singleton blocked by the intentional production quorum.

Build46 leaves those quorums unchanged and adds a diagnostic-only inspection of the already-qualified candidates. A candidate may be decoded alone only after blind geometry has been frozen and qualified; the resulting HMAC can describe the candidate but cannot influence geometry. A reference-assisted quadrilateral, when available, is used only after blind search to measure corner distance. This build is intended to decide whether the next bounded geometry experiment should preserve more independent qualified basins or improve the precision of the currently surviving basin.

## v0.3.0-build47 — frozen candidate observability

Build46 ruled out lowering the production quorum: the held-out-qualified singleton for B/mild is still roughly 62.7 px from the reference-assisted oracle and does not authenticate, while the B/angle singleton is a distant false basin. Build47 therefore preserves the exact production proposal tier and extends proposal breadth in two diagnostic stages: deeper cells on the same selected side pairs, then lower-ranked side pairs. No production threshold or decoder path changes.
