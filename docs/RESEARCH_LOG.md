# PixSeal research log

This file is an **append-only engineering/research notebook** for the v0.3 physical
print-acquisition work. `HISTORY.md` records milestones; `docs/RESULTS.md` records
checkpoint measurements; this file additionally preserves hypotheses, rejected
variants, threshold decisions and negative results so later work can reconstruct why a
path was accepted or abandoned.

Rules for future entries:

- never call a known-header improvement authentication; only Format-v3 HMAC is success;
- distinguish global-best hard values from same-candidate before/after comparisons;
- state whether a measurement is key-independent, key-assisted oracle, or HMAC-gated;
- record failed/rejected variants instead of deleting them from the historical record;
- record bounded search budgets and threshold changes, especially when a threshold is
  deliberately *not* relaxed after a promising result;
- private print-camera/scanner images are referenced by regression-case name only and
  are never copied into source archives.

## Compact chronology

| Build | Main question | Outcome |
|---|---|---|
| build1 | Can local lattice evidence be measured separately from production decode? | Yes; bounded diagnostic lattice path established. |
| build2 | Can print boundary + homography + virtual DCT sampling recover synthetic perspective? | Yes synthetically; real HMAC still absent. |
| build3 | Does spatial phase consensus improve projective geometry? | Useful evidence, insufficient real authentication. |
| build4 | Are scale aliases/harmonics hiding the fundamental carrier period? | Fundamental-scale and modulo-8 refinement added. |
| build5 | Is residual photometric variation the main blocker? | Small bounded bank helps some scores, no real HMAC. |
| build6 | Are large no-boundary photos blocked too early? | Adaptive finer lattice and weak-boundary lattice-first fallback added. |
| build7 | Is the failure before or after ECC? | Six known Hamming words expose protected-bit failure without revealing payload. |
| build8 | Can reliability-aware Hamming rescue weak words? | Synthetic yes; real cases no. |
| build9 | Are errors systematic or spatially local? | Mostly mixed/local; global inversion hypothesis rejected. |
| build10 | Can one low-DOF affine phase field remove spatial drift? | Some candidate improvements, no general solution. |
| build11 | Were integer local maxima/control quality causing false fields? | Yes; confidence/boundary-aware controls reject prior weak promotion. |
| build12 | Can controls be estimated without key/header? | Repetition-coherence blind observer established; first real blind 5 -> 4. |
| build13 | Can guided cross-cell correlation resolve cycle slips? | No; real pair score 0.053--0.065 vs ~0.61 synthetic. Kept diagnostic-only. |
| build14 | Can independent local lattice phase supply fractional registration and assist unwrap? | Stronger signal and large same-candidate gains, but unwrap remains mixed and unpromoted. |

## Build14 experiment notebook — local lattice fractional phase

### Hypothesis

The repetition observer can often identify an approximate integer tile/block phase but
retains cycle ambiguity. The pre-existing local lattice estimator measures physical
8-pixel periodicity independently of the key. Projecting a regional lattice
intersection back into the candidate canonical plane should therefore provide the
**fractional** phase inside one block. Combining integer repetition phase with this
modulo-block measurement may reduce spatial registration error without looking at the
known Format-v3 prefix.

### Fixed constraints

- Format v3 and encoder unchanged.
- Production `ExtractWithInfo` unchanged.
- Four complete/HMAC decode slots maximum.
- Two smooth corrected-grid resamples maximum.
- Nine lattice-phase controls maximum (one per 3x3 region).
- Lattice phase sees no key, magic, header bits, payload, CRC or HMAC result.
- Key-assisted local phase is retained only for ex-post oracle distance.

### Attempt A — build13 guided cross-cell retained as a check

The build13 guided cross-cell observer was not removed, because its weakness is useful
negative evidence. It remains an independent check but does not substitute for the
new lattice phase. Real pair scores previously measured around 0.053--0.065 remain far
below the 0.10 build13 consensus-strength gate.

Decision: **retain for diagnostics, do not tune further**.

### Attempt B — initial residual-as-inverse-correction sign

The first lattice implementation treated the canonical modulo-grid residual as if the
observer should directly encode the inverse sampler correction. On
`foto stampa storta.jpg`, the lattice controls were measurable and the smooth fit
improved geometrically, but the best same-candidate bit result was only **5 -> 5**.

Reviewing the internal smooth-phase convention showed that local controls represent
*observed phase drift*; `diagnosticFitBlindSmoothPhaseField` performs the sign inversion
when constructing the correction. The lattice observer was therefore changed to report
the observed canonical drift directly.

Decision: **reject inverse-correction sign; preserve this trial as convention history**.
The final sign choice is defined by internal phase semantics, not by selecting whichever
known-header result is numerically better.

### Attempt C — observed-drift sign, fractional phase only

With the final sign convention but integer unwrap disabled:

- `foto stampa storta.jpg`: best measured same candidate **5 -> 6**;
- `foto bici dritta.jpg`: best measured same candidate **7 -> 11**.

This is a useful negative result. A lattice phase reduced modulo one block cannot know
which integer block cycle is correct. Replacing/refining only the fractional component
without solving the cycle index is insufficient and can move the sampler in the wrong
periodic basin.

Decision: **fractional-only fusion is insufficient**.

### Attempt D — bounded affine-guided +/-1 unwrap (final build14 experiment)

The final research path keeps repetition as the initial integer-cycle estimate and
lattice phase as independent fractional evidence. A low-confidence primary control may
move by at most +/-1 block per axis only when a robust affine fit to all controls shows:

- lattice confidence >= 0.22;
- primary confidence <= 0.24;
- residual improvement >= 0.45 block;
- final residual <= 0.55 block;
- at most two unwrap passes.

These conditions are evaluated before known-header diagnostics.

Observed final same-candidate results include:

| Acquisition | Best hard /24 | Same candidate before -> after | Lattice consensus / cycle fixes | Smooth mean confidence | Fit / LOO RMS | Promoted to HMAC? |
|---|---:|---:|---:|---:|---:|---|
| `foto stampa.jpg` | 2 | -- | 3 / 0 | -- | -- | no |
| `foto stampa storta.jpg` | 7 | **5 -> 2** | 6 / 1 | 0.174 | 1.240 / 1.952 | no |
| `foto bici dritta.jpg` | 6 | **7 -> 3** | 5 / 2 | 0.145 | 1.387 / 2.032 | no |
| `foto bici storta.jpg` | 4 | 9 -> 8 | 3 / 3 | 0.196 | 0.727 / 1.279 | no |
| scanner `0270_001.jpg` | 4 | **9 -> 4** | 7 / 3 | 0.125 | 1.168 / 2.262 | no |
| scanner `0270_002.jpg` | 5 | **7 -> 10** | 6 / 2 | 0.201 | 1.557 / 2.121 | no |

The storta 5 -> 2 result is especially interesting because it reaches the same error
count as the key-assisted local oracle on that acquisition. It is **not** an HMAC
success and it is **not** the global hard candidate moving from 7 -> 2. It is one
separately ranked candidate moving 5 -> 2 after a blind field.

The method is not ready for promotion. The three large improvements all fail at least
one unchanged decode gate (especially the 0.20 mean-control-confidence threshold and,
in some cases, the 2.0 LOO threshold). Scanner 002 demonstrates a substantial
same-candidate degradation. Therefore no threshold was relaxed, zero build14 smooth
fields receive an HMAC slot, and all six physical acquisitions remain unauthenticated.

### Wrong-key control

On `foto stampa storta.jpg` with an unrelated key, final build14 remains
`projective-failed`, performs four ordinary full-decode attempts, zero smooth decode
attempts and authenticates no payload. Lattice-phase confidence remains measurable,
which is correct because it is key-independent geometry evidence rather than watermark
detection.

### Research conclusion

Build14 supports three statements and rejects a fourth:

1. **Supported:** local lattice phase is a materially stronger independent measurement
   than direct cross-cell DCT-pattern correlation on this corpus.
2. **Supported:** combining fractional lattice phase with a bounded integer unwrap can
   move some real same-candidate protected-bit channels dramatically closer to the
   known-header oracle.
3. **Supported:** a low geometric residual/LOO alone still does not guarantee a better
   bit channel, so confidence and HMAC gates must remain independent.
4. **Rejected:** the current affine-guided +/-1 unwrap is not yet reliable enough for
   decode promotion across acquisitions.

Next experiments should target the **discrete unwrap problem**, preferably with a third
key-independent signal (gradient/spectral phase, neighbourhood consistency, or a global
small-state smoothness objective). Do not increase HMAC slots, lower the 0.20 confidence
gate, lower the LOO gate, or introduce Format v4 solely to admit the current results.


## build15 — global discrete integer unwrap and transactional rollback

### Hypothesis
The build14 fractional lattice observer appears informative, while the integer cycle
assignment is the unstable component. Test whether a global smoothness objective can
choose the integer cycles independently of known Format-v3 bits.

### Attempt A — axis-separable bounded beam search
Each populated cell may retain its repetition cycle or move by exactly one block in
either direction. X and Y are solved separately with beam width 64. The score is a
combination of robust affine residual, confidence-weighted cycle-change penalty and a
second-difference curvature term. The search is deterministic and bounded; the declared
worst-case evaluation budget is 3456 states across both axes.

### Attempt B — accept geometric optimum without ambiguity margin (rejected)
Early real runs showed that a dramatically lower global objective can still correspond
to many nearly equivalent integer assignments. Treating the best state alone as
authoritative would recreate the build14 over-selection problem in a more complicated
form. The solver was therefore changed to retain the second-best state and require both
absolute and relative first/second margins.

### Attempt C — reporting baseline as best after rejection (bug, fixed)
An intermediate diagnostic path rolled the applied controls back correctly but mixed the
rolled-back baseline objective with the second-best proposed objective, allowing a
negative-looking margin in JSON. Reporting now separates baseline, proposed, applied and
second-best objectives. No decoding behavior depended on the reporting bug.

### Attempt D — fractional lattice survives ambiguous unwrap (rejected)
This was methodologically unsafe. Build14's explicit ablation already showed that
fractional-only correction can degrade the bit channel. Build15 therefore makes lattice
fusion transactional: if the integer unwrap is ambiguous/rejected, fractional lattice
changes are rolled back as well. A dedicated regression prevents future leakage.

### Representative real evidence
`foto stampa storta.jpg`: baseline objective about 5.480, proposed about 0.804,
second-best about 0.826, margin about 0.021. The proposal changes nine cells but is
ambiguous, so zero global changes are applied and the safe same-candidate diagnostic
returns to 5 -> 4 rather than build14's 5 -> 2.

`foto bici dritta.jpg`: baseline about 6.159, proposed about 1.267, margin about 0.011;
proposal rejected as ambiguous. Scanner `0270_002.jpg`: baseline about 7.318, proposed
about 1.723, margin about 0.026; proposal likewise rejected. These examples demonstrate
that geometric improvement alone does not establish a uniquely recoverable cycle field.

### Status
The global unwrap machinery is retained as RESEARCH. It improves the falsifiability of
the experiment by refusing non-unique cycle assignments, but it does not yet justify a
new real HMAC attempt. Format v3 and the production decoder remain frozen.

## build16 — exact top-2 certification and conservative structural validation

### Hypothesis
The build15 global formulation is useful, but a width-64 beam cannot certify the
first/second solution margin. First remove that approximation. Then ask whether a
key-independent structural comparison of the exact top-1/top-2 assignments provides
additional evidence without turning HMAC or the known header into a search oracle.

### Audit finding A — beam top-2 is not globally certified
A deterministic synthetic field was found for which the build15 beam accepted an X
unwrap with objective 0.027408951002 and beam runner-up 0.062494434795. Exact bounded
enumeration finds the same optimum but a better runner-up at 0.035129429298. The exact
margin falls below the fixed gate, so the correct result is ambiguous. This is a
false-accept hole in the beam uniqueness certificate, not merely a numerical change.

### Attempt A — exact axis enumeration (retained)
Build16 enumerates all `3^N` assignments for the `N <= 9` eligible controls of each
axis. The objective and thresholds are unchanged; only search completeness changes.
Maximum work is 19683 assignments per axis / 39366 across X/Y. A permanent regression
contains the deterministic former beam false-accept field. Exact top-1/top-2 ranking,
no-slip behavior, high-confidence locking, ambiguous rollback and HMAC restoration on
a controlled synthetic field remain covered.

### Audit finding B — mixed-axis status could be misleading (fixed)
The outer lattice fusion already refused an ambiguous global result, but the inner
solver could temporarily apply an accepted X or Y axis while the other axis was
ambiguous and report the global state as accepted. Build16 makes the transaction
atomic at the solver boundary too: one ambiguous axis means global `ambiguous`, zero
applied axes and baseline applied objective. A dedicated regression preserves the
per-axis evidence while asserting zero committed changes.

### Audit finding B — non-finite no-eligible reporting (fixed)
Build15 could leave `secondObjective=+Inf` when no control was eligible. Go's JSON
encoder rejects non-finite floats. Build16 reports the axis as `not-applicable`, keeps
numeric placeholders finite, and exports whether a second solution is actually
available. A regression serializes the public diagnostic structure for this path.

### Attempt B — split repetition top-2 validator (diagnostic-only)
The exact top-1/top-2 combined assignments are additionally compared using two
deterministic disjoint subsets of the existing Format-v3 repetition pairs. Each fold
measures same-coded-bit agreement only; no expected value, key, CRC or HMAC is used.
The signal is exported as `split-repetition-top2` with per-fold deltas and agreement.
It is not allowed to override an exact ambiguous result because the primary repetition
controls were originally estimated from the complete pair set. Treating this split as
fully held-out would overstate its independence.

### Physical evidence available in this build16 session
Four of the six canonical acquisitions were available. The two bicycle photographs
were absent, so no six-case promotion claim is made. On the available cases exact
ranking preserves the conservative build15 behavior. In particular, `foto stampa
storta.jpg` remains ambiguous with the same exact top-2 margin (~0.02113), while
scanner `0270_002.jpg` exposes the beam approximation clearly: its reported build15
margin ~0.02634 becomes exact ~0.01780. The proposal is still rejected, so the safe
behavior is unchanged.

### Status
Exact bounded enumeration is retained and replaces beam ranking. Split-repetition
validation is retained as research telemetry only. The next integer-cycle experiment
must use proposal and validation evidence that are disjoint by construction and must be
evaluated on the complete six-image corpus before it can relax an ambiguity decision.

## build17 — coded-bit-group-disjoint repetition cross-fit

### Hypothesis

Build16 showed that repetition scores contain information about the geometric top-1 vs
runner-up, but its fold scores were not truly held out because the primary repetition
controls had already been estimated from all repetition pairs. The build17 hypothesis
was that a proposal constructed from one independent subset of Format-v3 repetition
structure could be validated by the other subset strongly enough to discriminate the
integer cycle without expected bits, key material or HMAC.

### Predeclared construction

- Keep Format v3, encoder, production decoder, exact build16 unwrap objective and
  HMAC/smooth-resample budgets frozen.
- Partition by coded-bit group, not by individual pair. All repeated positions belonging
  to one coded bit stay in one fold. The two folds therefore share no logical repetition
  position.
- Direction A->B: estimate global/local repetition controls using fold A only, combine
  with the existing key-independent fractional lattice phase, enumerate exact top-1/
  top-2, then score only those two assignments with fold B.
- Direction B->A repeats the experiment with roles reversed.
- Compare the two independently proposed recentered local integer-cycle fields. Do not
  compare fold-specific absolute global tile anchors.
- `crossfit_supports_best` requires both held-out directions to favor their own top-1,
  at least three comparable cells and complete local-cycle agreement. This flag is
  diagnostic only and cannot override all-pairs ambiguity.
- Run cross-fit only when the normal all-pairs exact solver is already `ambiguous` and
  has an exact runner-up. This bounds duplicate search and prevents the experiment from
  becoming a general candidate multiplier.

### Synthetic regressions

`make crossfit-unwrap-test` verifies:

1. robust and balanced repetition groups partition into two non-empty folds;
2. no logical repetition position appears in both folds;
3. fold-specific repetition estimation uses the bounded proposal subset;
4. a held-out fold can prefer a known correct phase over a one-cycle alternative;
5. both symmetric directions execute through the exact solver on a bounded synthetic
   field.

### Physical evidence (four available acquisitions)

| acquisition | all-pairs result | A->B delta | B->A delta | cycle agreement | overall cross-fit |
|---|---|---:|---:|---:|---|
| `foto stampa.jpg` | not-applicable, best hard 2/24 | n/a | n/a | n/a | not run |
| `foto stampa storta.jpg` | ambiguous, margin 0.02113 | +0.00465 | +0.20522 | 2/9 | reject |
| `0270_001.jpg` | ambiguous, margin 0.02144 | -0.18283 | -0.02287 | 0/9 | reject |
| `0270_002.jpg` | ambiguous, margin 0.01780 | +0.05639 | -0.24050 | 0/9 | reject |

The best-candidate fold profile is balanced in all three ambiguous cases; the group
partition yields 220 proposal pairs in fold A and 228 in fold B. Exact cross-fit state
counts remain below the declared worst-case ceiling. Measured end-to-end diagnostic
runtime remains around 10--15 seconds for the reported best cases once cross-fit is
restricted to already-ambiguous candidates.

### Interpretation

The hypothesis is only partially supported. Truly held-out repetition evidence is
clearly informative: scanner 001 independently rejects the geometric top-1 twice, while
scanner 002 exposes directional instability. However, even the inclined smartphone case
where both held-out directions are positive does not reproduce the same integer-cycle
field across folds. Positive top-1 validation alone is therefore insufficient.

Build17 does **not** promote any unwrap, does not consume an extra HMAC slot and does not
reinterpret lower known-header errors as success. The next problem is candidate
stability/reconciliation across independent evidence partitions, or an additional
key-independent signal that can anchor the cycle field itself.



## build18 — lazy cross-fit and cell-level instability localization

### Hypothesis

Build17 showed useful held-out information but unstable local cycle fields. Before
adding another observer, determine whether agreement is concentrated in high-quality
controls, and remove duplicated cross-fit work from candidates that cannot become the
reported best bit-channel result.

### Implementation

- Freeze Format v3, production extraction, all-pairs exact objective/thresholds,
  build17 fold partition and held-out scoring.
- Run all-pairs exact unwrap on every existing bit diagnostic as before.
- Defer A->B/B->A cross-fit until bit diagnostics are ranked; run it at most once on the
  final best candidate and only for `ambiguous` + exact-runner-up cases.
- Export cross-fit milliseconds, attempts and skipped candidates.
- Export per-cell all-pairs proposed shifts, A/B local cycles, A/B confidence and cycle
  agreement. Summarize agreeing/disagreeing joint fold confidence and lattice confidence.

### Four-case physical result

The scientific decisions are identical to build17. The frontal photograph remains
`2/24` and not-applicable. The inclined photograph remains ambiguous with both held-out
deltas positive and only 2/9 cycle agreement. Scanner 001 remains 4/24 with both
held-out deltas negative and 0/9 agreement. Scanner 002 remains 5/24, directionally
split, with 0/9 agreement.

The new localization falsifies a simple partial-consensus hypothesis. On the inclined
photo, the two agreeing cells are `(0,1)` and `(1,2)`, but their mean joint fold
confidence is only `0.1337`; the seven disagreeing cells average `0.3576`. Lattice
confidence is `0.4059` for agreeing versus `0.3772` for disagreeing cells, not enough to
rescue the weak repetition consensus. The agreement subset is therefore not a
high-confidence core that can be promoted safely.

On the development host the selected held-out experiment costs only tens of milliseconds
(26/48/73 ms on the three ambiguous best candidates), with three lower-ranked bit
candidates skipped. Absolute end-to-end timings are host-dependent; the invariant to
preserve is one cross-fit attempt maximum per image candidate ranking.

### Status

Retain build18 telemetry and lazy scheduling. Do not build a decoder correction from the
2/9 inclined-photo agreement. The next experiment should add or derive a genuinely
independent cycle anchor, or demonstrate stability across additional deterministic
partitions/acquisitions, before relaxing any ambiguity gate.


## build19 — multi-partition integer-cycle stability

### Hypothesis

Build17/18 may have observed an unlucky two-way split. If repetition evidence contains a
real cycle anchor, the same recentered local integer field should recur when coded-bit
groups are repartitioned several deterministic ways, especially among proposals that are
positively validated by their held-out complement.

### Frozen experimental design

- 8 deterministic binary partitions, fixed before physical evaluation.
- Each coded-bit group is indivisible; proposal and validation evidence never share a
  repetition group.
- Both directions are evaluated for each partition: 16 maximum trials.
- Partition 0 is exactly the historical build17 A/B split; the other seven are distinct
  modulo complement.
- The experiment runs only on the final best bit candidate, only after all-pairs exact
  unwrap is ambiguous with an exact runner-up.
- No modal vote, supported-trial filter or stability statistic can alter sampling, smooth
  resampling, full-decode count or HMAC attempts.

### Four-case result

| case | trials | held-out supported | unique complete fields | supported unique fields | mean cell modal | supported mean modal | mean pairwise agreement |
|---|---:|---:|---:|---:|---:|---:|---:|
| frontal photo | not run | n/a | n/a | n/a | n/a | n/a | n/a |
| inclined photo | 16/16 | 9 | 16/16 | 9/9 | 0.2153 | 0.2716 | 0.0731 |
| scanner 001 | 16/16 | 6 | 16/16 | 6/6 | 0.2153 | 0.2963 | 0.0759 |
| scanner 002 | 16/16 | 12 | 16/16 | 12/12 | 0.2361 | 0.2593 | 0.0824 |

No cell is unanimous across all 16 trials in any ambiguous case. Even the most frequent
per-cell cycle receives only a small minority of votes. Scanner 002 is particularly
instructive: 75% of trials positively validate their own top-1, yet all 12 supported
complete fields are different. Positive held-out validation therefore does not imply a
reproducible field.

### Interpretation / decision

The hypothesis is rejected on the current corpus. The build17 A/B instability is not an
artifact of one unlucky split: repetition-derived cycle fields are highly partition
sensitive. Additional repetition voting or threshold relaxation would convert instability
into confidence rather than add independent information. Keep exact rollback and all HMAC
budgets unchanged. The next research branch must add a genuinely independent cycle anchor
or physical model.

## build20 — independent cross-cell cycle-anchor experiment

Hypothesis: unguided cross-cell image-domain registration, which uses neither coded-bit
repetition nor any key/header oracle, may supply an independent absolute-relative anchor for
the exact integer-cycle top-2 ambiguity.

Predeclared construction: run only on the final best candidate when the exact all-pairs solver
is ambiguous and has a runner-up. Build the existing unguided pairwise correlation graph,
solve its zero-mean relative offsets, and score exact top-1 and runner-up after removing the
best common integer x/y gauge. No threshold can promote a field in build20; sign and agreement
are telemetry only.

Physical result: inclined smartphone delta(second-top1)=+0.5083, mean confidence 0.0575;
scanner001 +0.2277, confidence 0.0168; scanner002 +0.0956, confidence 0.0798. All three prefer
top-1 continuously, but both top-1 and runner-up have 0/9 rounded-cycle agreement with the
anchor. Frontal smartphone is not-applicable and the anchor is not run. Scanner001 contradicts
the repetition held-out preference; scanner002 negative control also prefers top-1.

Conclusion: pairwise registration contains continuous relative-shape information but is not an
absolute integer-cycle anchor on the current corpus. No unwrap, HMAC, encoder or format rule is
changed.


## build21 — Format-v3 key-independent observability audit

### Question

After build19 rejected repetition repartition/voting and build20 rejected unguided pairwise
registration as an absolute cycle anchor, is the remaining failure primarily a solver problem,
or does frozen Format v3 expose only weak/non-uniform absolute-origin information?

### Frozen audit

No physical threshold is tuned. For all three concrete profiles and all eight unit block shifts,
measure two public key-independent mechanisms:

1. repetition-pair topology overlap induced only by `v3CodeIndex`;
2. Hamming syndrome survival on deterministic valid whitened-bit codewords after the real wrong
   tile-origin mapping and hard aggregation.

No key, magic/header expectation, payload, CRC or HMAC is used. No audit result can alter decode.

### Result

Robust and balanced do have a unique repetition topology, but the nearest horizontal unit alias
retains ~94% of the pair graph and the vertical alias retains ~92%. Capacity has no repetition
pairs. Hamming parity is strong against most horizontal/diagonal shifts but has a vertical blind
spot: capacity ±1 vertical is an exact whole-word permutation symmetry; robust is zero-syndrome
on the deterministic probe; balanced is only weakly separated vertically.

### Decision

Do not add another threshold around existing repetition/pairwise evidence. The frozen format is
not uniformly observable under the audited absolute-cycle mechanisms. Keep Format v3 frozen for
v0.3 research, and make the next design decision explicit: either derive a stronger physical
model from independent image evidence, or design a future format revision with an intentional
absolute-origin/asymmetry pilot. Any format revision must remain separate from the v3 corpus and
production decoder baseline.


## build22 — held-out physical topology observability and v4 decision branch

### Hypothesis

Build21 found a small but real robust/balanced repetition-topology asymmetry. If that 5–8% residual
survives print-camera noise, it may provide an absolute origin signal once the 92–94% shared pair
evidence is removed from the decision score.

### Predeclared experiment

For each unit shift, split repetition edges into shared registration edges and canonical-exclusive
held-out edges. Shared edges register two phase hypotheses separated by exactly the tested shift;
registration maximizes the weaker of their two shared-edge scores. The held-out edges then score
A versus B globally and in every complete 3×3 spatial cell. No key, header, payload, CRC or HMAC
participates, and the result cannot alter decode.

### Result

The signal is measurable, but winner phases are not stable enough to identify a unique origin.
Inclined smartphone mean modal phase fraction is 0.3125; scanner 001 and scanner 002 are both
0.4375. Mean cell directional consistency is 0.6944, 0.7292 and 0.6597 respectively. Scanner 002
therefore remains fatal to promotion: a negative control can look more phase-stable than the target
smartphone acquisition. No HMAC authenticates.

### Decision

Reject another v3 threshold/weighting round. The remaining v3 asymmetry is too weak and too
non-specific on the current physical corpus. Preserve the probe as evidence, but do not use it to
steer exact unwrap or sampling.

In parallel, begin a serious Format-v4 branch. The first quantified candidate is a 37×32 tile with
64 public absolute-pilot positions and 1120 data positions. This preserves v3-like payload ceilings
while increasing minimum tile area only 5.7%. Before implementation, optimize/freeze a pilot mask
and sign sequence and define promotion tests against physical negative controls.

## build23 — close v3 absolute-cycle research, start isolated v4 foundation

### Decision input

The user-host build22 run confirmed both deterministic and physical development results. The 30-target
qualification retained release baseline PASS and qualification corpus PASS. The held-out physical topology
probe did not separate scanner 002 negative control from scanner 001 or the inclined smartphone image.
Therefore another v3 topology threshold, vote or partition would not add independent information.

### v3 decision

Close the absolute-cycle branch. Keep v3 implemented, interoperable and regression-tested, but treat future
changes as maintenance unless a genuinely independent observable is introduced. A five-file SHA-256 freeze
manifest makes core drift explicit and reviewable. Preserve all negative results
so future work does not rediscover the same repetition/pairwise aliases. HMAC remains the only physical PASS.

### v4 foundation hypothesis

An intentional sparse public pilot can remove the exact integer-cycle symmetry while preserving all 1120 v3-like
data positions if the tile expands from 35×32 to 37×32. Before implementing a codec, test whether a concrete
pilot can satisfy structural invariants with low wrong-shift overlap/correlation.

### Prototype-1 result

A deterministic stratified search family produced `prototype-1-stratified-p64`: one pilot in each 8×8 spatial
stratum, 32 positive and 32 negative signs. Exhaustive 37×32 cyclic audit gives max overlap 8/64 and max absolute
signed wrong-shift correlation 5/64, with no perfect alias. This is promising enough to support the architecture,
but it is deliberately not frozen as Format v4. Partial-crop and simulated print-camera behavior remain untested.

### Next experiment

Search broader pilot mask/sign families with a predeclared multi-objective score including full-tile cyclic
separation, crop survival and simulated print-camera degradation. Only after that should the project define a
version marker and build an experimental v4 encoder/pilot-only detector.

## build24 — reproducible pilot search, partial visibility and first image channel

### Question

Can Build23's 64-symbol public-pilot architecture be improved by a reproducible search, and does the resulting
candidate still identify absolute origin after partial visibility and simple image-channel degradation without
consulting payload, key, header, CRC or HMAC evidence?

### Predeclared search

Use one pilot per 8x8 spatial stratum and exact 32/32 sign balance. Stage 1 evaluates 200,000 fixed-seed joint
coordinate/sign candidates. Candidates worse than prototype-1's 8-position maximum overlap or 5-symbol maximum
wrong correlation are rejected before crop scoring. The retained family is ranked lexicographically by maximum
overlap, maximum absolute wrong correlation, then worst contiguous margins from 16 toward 48 visible pilots.
Stage 2 keeps the winning mask and evaluates 200,000 fixed-seed balanced sign sequences with the same correlation/
crop objective. Candidate hashes, budgets and runner-up reports are retained.

### Structural result

`prototype-2-search-p64` (`858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`) improves
prototype-1 from max overlap 8 to 7 and max absolute wrong correlation 5 to 4. No exact cyclic alias exists.
Worst contiguous margins for 64/48/32/24/16 visible pilots are 60/43/28/20/12. Across 256 deterministic random
subsets per visibility level, no wrong origin ties or beats the correct structural origin.

### Initial image-domain experiment

Render a synthetic carrier with prototype-2 pilot signs and deterministic pseudo-random data-plane signs using
the real PixSeal DCT block embedding primitive. A pilot-only detector assumes the 8/6/4-pixel lattice is already
resolved and scores every 37x32 cyclic origin. The deterministic test image recovers the correct origin after
JPEG-82, blur, +/-4 noise, gamma 1.15, exact 75% and 50% resize and aligned crop. The unmarked negative control
shows a small runner-up margin and is never interpreted as authenticated evidence.

The two supplied local originals produce the correct origin in the same initial matrix. Their unmarked controls
produce best scores 0.306699 / 0.232802 with margins 0.022743 / 0.003697, while marked transformed cases retain
positive correct-origin margins (lowest observed initial local margin 0.381056). These figures are evidence only;
no production threshold is selected from two images.

### Decision

Keep prototype-2 as the current non-normative v4 candidate and continue Build24 qualification. Do not implement
a production v4 codec or freeze interoperability yet. The next independent tests are blind rotation/affine/
perspective geometry, combined degradation and then a separate real v4 print-camera/scanner corpus.



## build25 — prototype-2 geometric survivability with known mapping

**Question.** Does the Build24 pilot remain an absolute-origin signal after the geometric channel, before spending effort on blind geometry search?

**Method.** Add a projective pilot sampler that receives a canonical-to-observed homography independently of pilot scoring. Evaluate 16 fixed transforms (rotation, anisotropic scale, shear, perspective, combined rotation/scale/shear, 75% scale, JPEG/blur/noise after perspective, and crop combinations) on a synthetic background and on both local originals. For each positive, evaluate the correspondingly transformed unmarked image as a negative control. Require correct origin plus a predeclared 0.10 positive-minus-negative margin separation.

**Result.** 16/16 synthetic and 32/32 local-image cases select origin `(0,0)`. Synthetic minimum positive margin 0.531955. Local minimum positive margin 0.322069; maximum negative margin 0.067786; minimum observed separation 0.258085.

**Decision.** Keep prototype-2. Do not freeze it yet. The evidence now supports moving to bounded blind pilot-assisted geometry estimation; it does not justify a production threshold, v4 encoder promotion or authentication claim.

## build26 — blind geometry proposal from repeated v4 data-plane structure

### Question

Can the Build25 pilot survive the removal of the supplied homography without turning the pilot into an unbounded brute-force geometry oracle?

### Initial negative result

A first bounded search ranked thousands of angle/scale/shear/perspective hypotheses directly with pilot correlation. It worked on synthetic images but did not generalize reliably to the local photographic backgrounds: small scale quantization errors could destroy the DCT phase, while unrelated geometries could win locally on image content. Lowering thresholds was rejected.

### Revised hypothesis

The v4 tile already repeats spatially. Geometry can therefore be proposed from key-independent self-consistency of corresponding **data-plane positions** across repeated tiles without knowing their symbols. The pilot can then do the job it was designed for: absolute origin and candidate validation.

### Method

Sample 64 deterministic non-pilot residues. For each candidate homography, compare weighted DCT signs at homologous positions of adjacent repeated tiles. Use the repeat score plus only a weak canvas-dimension prior to retain a small bank. Evaluate pilot phase on central tile rows, then validate survivors on held-out corner repetitions. Search rotation/scale/shear/perspective deterministically and cap the synthetic regression at 6400 evaluated repeat hypotheses per image.

### Result

The revised search recovers the complete synthetic Build26 family and the representative local-corpus cases. In particular, the previously failing `PJ_piccolo` combined-perspective case now recovers the declared 9.3 degree, 1.08/0.92, 0.030/0.015 transform exactly with origin `(0,0)` and pilot margin 0.425713. The data-repeat proposal score is unchanged if every pilot sign is inverted, confirming that proposal does not depend on pilot symbols.

### Decision

Accept Build26 as the first bounded blind pilot-assisted geometry checkpoint, but do not freeze the pilot or call the problem solved. The next independent problem is unknown placement/crop/translation and eventually physical print-camera geometry. Keep payload/header/ECC/HMAC out of geometry proposal and validation.



## build27 — isolate crop/translation placement before joint blind recovery

### Hypothesis

Once Build26 has shown that the public pilot can rank bounded geometry hypotheses, the next unknown should be isolated rather than folded immediately into a larger brute force. If the geometric mapping is correct but crop/translation is unknown, the public pilot should recover placement without payload/header/HMAC evidence.

### Method

Supply the canonical-to-full-frame geometry but withhold the crop or padded-canvas translation. Derive a finite translation interval from the difference between full transformed extent and observed extent. Pilot indices `i mod 2 = 0` propose placement on a 4-pixel grid and integer refinement; indices `i mod 2 = 1` validate the shortlist. Require finite scores, bounded hypothesis count, full-pilot recovery under the selected mapping, matching unmarked controls, and a deliberately wrong-geometry control.

Do **not** require a placement runner-up gap: because tiles repeat, translations differing by a full tile can be equivalent and forcing one to win would manufacture confidence from an unobservable quantity.

### Result

All synthetic crop/pad cases pass, including mild perspective and a deeper crop. Both local originals pass the representative crop and padded-perspective cases. The wrong-geometry control is strongly separated from the correct transform.

### Decision

Accept placement recovery as independently qualified for Build27, but keep `prototype-2-search-p64` non-normative. The next experiment must combine a crop-tolerant coarse lattice/extent proposal with placement and pilot validation. Do not call Build27 a fully blind decoder and do not introduce payload/ECC/HMAC evidence to bridge the remaining gap.


## build28 — joint blind affine geometry plus unknown crop

### Question

Can the geometry and placement problems qualified separately in Build26/27 be solved in the same bounded search on both synthetic and photographic backgrounds without using pilot symbols to choose geometry?

### Rejected paths

A repeat-only global proposal under crop was too narrow in angle/scale and became texture-sensitive on `PJ_piccolo.png`. A sparse absolute-DCT bank produced false photographic maxima. Coordinate-wise final refinement also stalled at a nearby `scaleY~0.915` basin even after the correct positive-angle family reached the final shortlist. Lowering pilot/validation floors was explicitly rejected.

### Retained method

Use absolute DCT carrier differential and **phase contrast** as a public structural observable. Search angle at 0.25-degree coarse spacing with bounded anisotropic scales; refine the best angle basins with 64 and then 256 distributed blocks. Apply one compact coupled angle/scale grid only to the final two structural candidates. The coupled step moves `PJ_piccolo.png` to approximately 11.25 degrees / 1.0675 / 0.93 and makes it the highest structural candidate.

Do not use pilot partition 0 to select between geometry basins. Once the DCT winner is fixed, run the unchanged Build27 placement search: partition 0 proposes crop translation, partition 1 validates, and full pilot scoring checks cyclic origin.

### Result

Both synthetic cases pass within ~69k total hypotheses. Both local originals pass the representative 11.2 degree / 1.07/0.93 + crop case. `PJ_piccolo.png` reaches held-out validation 0.781974, full-pilot margin 0.410485 and origin `(0,0)`; `PJ_lingua.PNG` reaches 0.812461 / 0.432900 and `(0,0)`. Matching unmarked images remain well separated.

### Decision

Accept Build28 as the first joint blind **affine+negative-crop** checkpoint. Keep prototype-2 non-normative. Do not infer that projective crop or padded-canvas geometry is solved: perspective breaks the simple translation/phase relation used by the affine structural proposal, and positive placement expands the search differently. Those become the next explicit research gates before pilot freeze and v4 codec implementation.

Also document the product architecture decision: Go is retained partly for one cross-platform core (Linux/Windows/Android/iOS), and a future GUI—especially mobile—is planned over that same core.

## build29 — joint projective crop and padded placement

### Question
Can Build28/Build27 be composed when perspective and crop are both unknown, and when affine geometry must be recovered while the carrier is positively displaced inside a larger canvas?

### Rejected paths
A brute geometry×placement pilot search was too expensive. Sparse repeat-only pruning produced multiple-comparison false maxima. Single-tile pilot ranking overfit local content, and splitting the pilot into independent maxima allowed each half to choose a different false phase. A direct homography correction fitted from per-tile origin drift worsened corner error and was discarded. Threshold relaxation was not used.

### Retained method
Use dimensional feasibility + structural DCT to form a broad bank. Evaluate two 32-pilot halves on the same phase/origin hypothesis and rank by the weaker half. Preserve distinct geometric basins, refine only the shortlist, use spatial evidence to reduce the number of expensive placement searches, then run Build27 placement and complete-pilot absolute-origin/margin checks. Insets include the previously missing 0.015 half-step. Reuse a precomputed pixel plane so ranking thousands of candidates does not rebuild the image for every hypothesis.

### Result and decision
Synthetic projective+crop and affine+padded cases are accepted. `PJ_lingua` projective+crop is accepted with ~0.384 validation, ~0.180 margin and ~0.0088 corner error. `PJ_piccolo` projective reaches validation ~0.421 but margin only ~0.119 and is therefore SAFE REJECTED. Both photographic padded joint searches are SAFE REJECTED; a known-geometry `PJ_lingua` padded control reaches ~0.991 validation / ~0.554 margin, proving the channel remains available. Accept Build29 as a bounded composition/safe-rejection checkpoint, not as general projective/padded recovery, and keep prototype-2 non-normative.



## build30 — pilot lock readiness, not normative promotion

### Question
After Build29, are the remaining SAFE REJECT cases evidence that `prototype-2-search-p64` itself should still be replaceable, or are they decoder geometry/placement limitations? Is there enough evidence to protect the candidate from accidental mutation without prematurely declaring a normative v4 format?

### Rejected padded-ranking variant
A natural hypothesis was that Build28/29 DCT geometry ranking failed under positive padding because the phase bank sampled only even sub-block offsets while the padded fixture has an odd residual phase. A diagnostic variant scanned all 64 8×8 sub-block phases. It does expose the odd phase, but at 16/32-block support the true photographic geometry is still overwhelmed by natural texture; at higher support the cost rises without producing a robust global ordering. Do not add this bank and do not relax Build29 gates.

### Isolation experiment
Bypass geometry and placement search entirely. For each private development original, render the existing Build29 projective-crop and affine-padded marked/control fixtures and pass the exact known mapping directly to the public pilot detector. If a Build29 SAFE REJECT is caused by a weak pilot, the marked score/margin should also collapse under the true mapping.

### Result
All four marked observations recover origin `(0,0)`. Scores are 0.986195 / 0.988352 on the larger original and 0.928412 / 0.932220 on the smaller original. Margins are 0.559445 / 0.553797 and 0.490432 / 0.484086 respectively. Matching unmarked margins are 0.007237 / 0.000598 / 0.010553 / 0.010305. Thus even the small padded case retains a strong pilot under the true mapping although the Build29 placement/joint search cannot recover it.

### Decision
Lock the exact candidate identity (`858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`) against accidental mutation and add an explicit source check plus local known-mapping corpus gate. **Do not call this a normative pilot freeze.** The physical print-camera/scanner channel is central to PixSeal and has not yet been tested with a real v4 encoder. Normative promotion waits for that evidence; a physical pilot failure must result in a new candidate/version rather than mutation of the Build30 lock.

## build31 — first authenticated v4 frame and encoder

**Question.** Can the locked public pilot be combined with a real authenticated data plane without changing v3, losing the 16/32/64-byte comparison ceilings, or creating an accidental pilot/data overlap?

**Design.** Keep prototype-2 exactly as locked in Build30. Remove its 64 positions from the 37x32 tile and enumerate the remaining 1120 positions row-major. Retain 32/48/80-byte profile frames and Hamming(7,4) as a controlled baseline. Use authenticated header bytes `0x41/0x42/0x43`, whitening domain `pixseal-whiten-v4`, and HMAC domain `pixseal-frame-v4 || 0x00`. Keep stable v3 APIs and commands unchanged; expose only explicit `v4-embed` and aligned `v4-extract` research commands.

**Implementation finding.** The historical helper `experimentalV4PrototypeDataPositions()` still belongs to Build23 prototype-1. Reusing it with the locked prototype-2 immediately triggers a data-position mapping failure because the two pilot masks differ. Do not edit that historical helper: Build31 adds a prototype-2-locked partition and an exhaustive 1184-position regression so every tile position is exactly pilot or data.

**Result.** All three profiles authenticate on a synthetic carrier and on both private development originals. Robust also authenticates after JPEG q82 and a block-aligned crop; the crop pilot recovers the expected cyclic origin. The minimum 296x256 single tile authenticates a capacity payload. Wrong key fails, v3 and v4 frames reject each other, and a deterministic frame/protected-bit vector is fixed for Build31 regression.

**Decision.** Accept Build31 as the first functional **experimental** v4 codec. Do not promote its Hamming choice or frame vector to normative status yet. Do not connect HMAC outcome to geometry selection. Use this exact encoder to create a new private v4 print-camera/scanner corpus; normative pilot promotion requires HMAC-authenticated physical recovery, not pilot-only correlation.


### Build32 — corpus governance before physical v4 qualification

The prior directory-discovery behavior allowed newly added originals to change legacy research suites implicitly. Build32 replaces that with an explicit LQ/MQ/HQ manifest and SHA verification. No pilot/frame/encoder algorithm change was made. LQ exposed a real capacity boundary at 50% resize, so the common baseline moved to 55% while the exploratory limit suite retains 50% and lower points. The next evidence source is physical print/scan/phone acquisition using the new v4 test key.

## 2026-09-15 — Build33 MQ joint-projective ranking observability

Build32 exposed one architectural Format-v4 blocker that must be kept separate from the four legacy/LQ research reds: the MQ joint projective+crop search SAFE REJECTS despite very strong known-mapping pilot evidence. Build33 therefore starts with observability rather than another wider brute-force bank.

A new diagnostic recreates the exact Build29 MQ transform and tracks the known true/near-true basin through the structural bank, half-pilot prefilter and full-proposal ranking. The experiment is repeated with both the historical random synthetic data plane and a real authenticated robust v4 carrier (`v4-b33-auth`, strength 24, key `PixSeal-v4-TestKey-2026`). Payload/header/CRC/HMAC are never consulted for geometry ranking.

Exploratory local variants tested while preparing this checkpoint included translation-aware sparse pilot ranking, three-region spatial proposal/check/held-out scoring, a coupled evolutionary beam over geometry+translation, and a half-pilot-only beam. Some variants recovered the correct basin on one MQ data plane (including the authenticated carrier) but promoted a convincing false basin on the other. They are therefore rejected rather than merged. This is evidence that a green result on one carrier is insufficient and that Build33 must explicitly control data-plane sensitivity and basin diversity.

The next decoder change should be evaluated against both MQ carriers and matched negatives. The success criterion remains correct blind geometry and placement first, followed by data sampling and a valid Format-v4 HMAC. LQ SAFE REJECT remains acceptable; lowering Build29 floors is not an option.

## 2026-09-15 — Build34 authenticated MQ projective recovery

Build33 showed that the correct MQ basin was already structurally observable but was destabilized by the centered pilot proposal. Build34 identifies the concrete phase defect: the former 0/+-4 pixel phase cross skipped a strong +/-2 pixel alignment. The accepted design keeps two distinct structural anchors, searches a bounded coupled local geometry neighborhood, retains candidates by public structural evidence, then ranks only with the public pilot over interior tile rows using 0/+-2/+-4 pixel centered phases. Top/bottom tile rows are not consulted until one geometry is committed.

Non-zero cyclic origin is canonicalized by a domain-side block translation of the homography and must then re-detect as origin `(0,0)` under the unchanged acceptance gate. Only after acceptance are the 1120 data positions sampled and decoded through Hamming, whitening, frame parsing, CRC and HMAC. On the active MQ carrier, two robust payloads with the public test key recover exactly with corner error about 0.001, validation about 0.86-0.88 and pilot margin about 0.43. LQ remains on the Build29 path and SAFE REJECTS.

## 2026-09-15 — Build35 physical qualification handoff

Build34 closes the active MQ digital blocker strongly enough that further decoder tuning without new evidence would risk overfitting the synthetic/corpus transform. Build35 therefore freezes that behavior for the next experiment and makes the physical test itself reproducible.

The preceding checkpoint is Build34 throughout source, test names and documentation. Build35 preserves that result unchanged while preparing the physical gate.

For controlled physical work, Build35 exports the Build34 joint-projective/HMAC path as `ExperimentalV4ExtractProjective` and `v4-extract-projective`. Canonical pre-print dimensions are explicit inputs. Geometry still uses only public structural/pilot evidence and the unchanged Build29 floors; payload/header/CRC/key/HMAC are unavailable until after geometry acceptance.

The first physical campaign is intentionally scanner-first and uses the active MQ source only. The generated pack contains one unmarked control and two marked robust/strength-24 carriers with distinct payloads. Printing and scanning both target 300 ppi/dpi at actual size, and scans are cropped only to the artwork edges. This isolates the paper/printer/scanner channel before adding free-camera scale, framing, lens and perspective variables. A physical PASS requires exact HMAC recovery of both marked payloads and rejection of the control; pilot-only evidence remains diagnostic.

## Build36 — first v4 physical scanner channel

The Build35 paper pack was acquired on an office scanner whose color path outputs JPEG. The useful corpus is the 600-dpi JPEG set; PDF/PPT export modes were inspected but contained lower-quality embedded raster data. Blind Build35 projective recovery SAFE-REJECTs the control as expected but also rejects marked A/B due geometry ranking.

A private reference-assisted registration diagnostic isolates geometry from channel capacity. Under independently supplied geometry, both marked scans recover their exact authenticated payloads with deterministic soft Hamming decoding. This is not a blind PASS and cannot promote the pilot, but it demonstrates that strength 24 and the existing v4 Hamming/frame/HMAC path survive the real print/scan/JPEG channel. Build36 therefore promotes soft Hamming only after geometry acceptance and leaves encoder/format thresholds untouched.
