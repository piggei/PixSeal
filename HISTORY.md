# PixSeal history

This file records technical evolution, including experiments that were later
superseded or deliberately not promoted.

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
