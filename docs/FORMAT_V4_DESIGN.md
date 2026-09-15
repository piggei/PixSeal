# Format v4 design study — absolute pilot branch

Status: **experimental, non-normative design branch**. Sizing was first quantified in v0.3.0-build22; build23 added the isolated pilot foundation, and Build31 adds the first explicit experimental v4 frame/encoder.
Format v3 remains the frozen production/interoperability baseline. Build31 v4 carriers are research artifacts and may still change before normative promotion.

## Why consider v4 now

Builds 16–22 progressively isolated the physical print-camera failure from solver details.
Exact integer-cycle search, held-out repetition cross-fit, multi-partition stability and an
independent pairwise image-domain anchor all reject unsafe promotion. Build21 then showed that
Format v3 is not uniformly cycle-observable: its strongest public structures leave weak aliases,
and capacity has a real vertical Hamming blind spot. Build22 measures the residual repetition
asymmetry in the physical images with held-out topology-exclusive pairs. That signal is present,
but its winner phase is not stable enough to distinguish the inclined smartphone case from the
scanner negative control.

A future format should therefore provide an **intentional absolute-origin signal**, rather than
asking payload redundancy/ECC to serve two unrelated jobs at once.

## Design goals

A candidate v4 should:

1. preserve authenticated payload semantics and a clear HMAC-only success criterion;
2. provide a **public, key-independent absolute pilot** that is structurally distinct from data;
3. let geometry/phase be estimated before payload/key-dependent sync work;
4. avoid the v3 vertical whole-Hamming-word alias;
5. keep current 16/32/64-byte payload ceilings if practical;
6. retain roughly the current robust/balanced/capacity data redundancy;
7. add only a small carrier-geometry and computational cost;
8. remain separately versioned so existing v3 carriers and corpus never become ambiguous.

## Build22 sizing model

For an apples-to-apples first comparison only, the study assumes v3's 16-byte framing overhead
(8-byte header + 8-byte truncated HMAC) and Hamming(7,4), which costs 14 carrier positions per
source byte. Those assumptions are **not** a final v4 ECC/framing specification.

| candidate | tile | pilot | data positions | ideal pilot z* | max payload under v3 framing/ECC | geometry vs v3 |
|---|---:|---:|---:|---:|---:|---:|
| compact | 35×32 | 64 | 1056 | 8.00 | 59 B | 1.000× |
| preserve-capacity | **37×32** | **64** | **1120** | **8.00** | **64 B** | **1.057×** |
| strong-pilot | 38×32 | 96 | 1120 | 9.80 | 64 B | 1.086× |

\* `sqrt(Npilot)` is only the idealized random-wrong-phase separation scale. It is not a physical
recovery guarantee.

### Recommended first prototype: 37×32 + 64 pilot positions

This layout is the current engineering favorite because it separates the new feature from the
data without sacrificing the v3 data budget:

- tile grows from 1120 to 1184 blocks (+5.7%);
- minimum carrier geometry grows from 280×256 to **296×256** pixels;
- 64 blocks (5.4% of the v4 tile) are reserved for the public absolute pilot;
- exactly 1120 data positions remain, so current profile capacities and nominal redundancy can
  be preserved if v3 framing/ECC were retained;
- an exhaustive pilot phase correlation requires about 75,776 pilot samples per geometry
  hypothesis (`1184 × 64`), small compared with current projective/full-grid work;
- because the pilot is key-independent, it can reject poor phase/geometry candidates before
  expensive key/profile decode attempts.

## Pilot requirements

The final 64 pilot positions and signs must be designed, not chosen ad hoc. A v4 pilot should have:

- low cyclic autocorrelation for every non-zero tile shift, especially ±1 X/Y;
- low overlap between shifted pilot masks;
- approximately balanced positive/negative DCT signs;
- spatial dispersion across the whole tile, avoiding a visually concentrated patch;
- enough local coverage that modest crop/occlusion does not erase all pilot evidence;
- a deterministic public definition independent of key and payload;
- explicit regression tests for every cyclic shift and for simulated print-camera noise.

A sparse pseudo-random mask plus a public pseudo-random ±1 sequence is a plausible starting
point, but the sequence/mask should be selected by an exhaustive deterministic search against
these metrics and then frozen in the format specification.

## Decoder architecture opportunity

A v4 decoder can invert the current dependency order:

```text
image / projective hypothesis
        ↓
public absolute pilot correlation
        ↓
absolute tile origin + pilot confidence
        ↓
data-only sampling / profile detection
        ↓
ECC / whitening / frame parse
        ↓
HMAC authentication
```

This can improve both reliability and performance: absolute phase is estimated from a signal
specifically designed for that task, and low-confidence geometries can be rejected before full
data decoding. It also keeps pilot confidence separate from authentication: only HMAC can still
make a payload valid.

## What v4 can and cannot guarantee

A dedicated pilot can **structurally guarantee that the format no longer has the specific v3
cycle symmetry**, once the pilot mask/sequence is chosen so every cyclic shift is distinct. It
cannot guarantee recovery from every photograph, printer, blur, crop or perspective condition.
Physical reliability must still be demonstrated on a corpus with false-positive controls.

Promotion criteria for a v4 prototype should therefore include, at minimum:

- unique correct pilot phase on synthetic distortion matrices;
- a large, predeclared margin over every wrong phase;
- stable phase across independent physical acquisitions;
- negative controls that do not cross the acceptance threshold;
- HMAC-authenticated recovery on the real print-camera/scanner corpus;
- measured runtime no worse than v3 diagnosis for equivalent images, preferably lower because
  of early pilot rejection.

## Build22 decision

The user-host qualification confirmed the development result: the weak v3 physical topology signal
does not distinguish scanner 002 negative control from useful acquisitions. The v3 absolute-cycle
research branch is therefore closed in build23; see `V3_FINAL_STATUS.md`.

## Build23 prototype foundation

Build23 still does **not** enable a v4 encoder or decoder. It introduces a separate code-level
prototype so the design can accumulate deterministic invariants before any interoperability promise.
The provisional candidate is named `prototype-1-stratified-p64` and uses the recommended 37×32 tile.

Structural properties frozen only for this prototype checkpoint:

- 64 unique pilot positions and exactly 1120 remaining data positions;
- one pilot in each cell of an 8×8 spatial stratification;
- exactly 32 `+1` and 32 `-1` pilot signs;
- exactly 16 pilots in each image-space quadrant;
- exhaustive toroidal audit over all 1183 non-zero shifts;
- no perfect non-zero cyclic alias;
- maximum shifted-mask overlap: **8/64 = 12.5%**;
- maximum absolute wrong-shift signed correlation: **5/64 = 7.8125%**.

These numbers are not yet Format-v4 constants. Prototype-1 may be replaced after stronger mask/sign
search, partial-crop scoring and simulated print-camera qualification. `watermark/experimental_v4.go`
is intentionally disconnected from `EmbedWithInfo`, `ExtractWithInfo` and CLI format selection.

## Build24 pilot search and initial qualification

Build24 makes the pilot search reproducible instead of treating prototype-1 as an offline result. The current
search is deliberately bounded:

```text
stage 1: joint coordinate + sign search
seed:    0x50585345414c5634
budget:  200000 candidates

stage 2: balanced-sign refinement on winning mask
seed:    0x76347369676e7331
budget:  200000 candidates
```

The hard family constraints remain 64 unique positions, one position in every 8x8 spatial stratum, 32 positive
and 32 negative signs, 1120 remaining data positions and zero exact non-zero cyclic aliases. Stage 1 admits only
candidates no worse than prototype-1's global overlap/correlation baseline before evaluating crop margins.

The current winner is `prototype-2-search-p64`. As of Build30 its exact identity is **development-locked but still non-normative**:

```text
candidate SHA-256                 858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174
maximum cyclic mask overlap       7 / 64
maximum wrong signed correlation  4 / 64
runner-up wrong correlation       4 / 64
perfect non-zero aliases          0
```

Prototype-1 remains preserved as the Build23 baseline (8/64 overlap, 5/64 correlation). The mean mask overlap is
3.408284 for both candidates because a 64-position mask has a fixed total number of ordered pilot pairs; the
useful improvement is in the worst shifts, not that mean.

### Partial visibility

Build24 evaluates contiguous regions of the 8x8 stratum grid and 256 deterministic random subsets at every
requested visibility level:

| visible pilots | prototype-1 worst contiguous margin | prototype-2 worst contiguous margin | prototype-2 worst random margin | random false origins |
|---:|---:|---:|---:|---:|
| 64 | 59 | **60** | 60 | 0/256 |
| 48 | 42 | **43** | 43 | 0/256 |
| 32 | 27 | **28** | 27 | 0/256 |
| 24 | 19 | **20** | 19 | 0/256 |
| 16 | 12 | 12 | 11 | 0/256 |

A margin here is `correct structural score - best absolute wrong-origin structural score`; it is not an HMAC or
physical-recovery confidence value.

### Initial image-domain channel

Build24 adds a synthetic qualification carrier that embeds the public pilot plus deterministic pseudo-random
signs in all 1120 data positions using the same DCT coefficient pair as PixSeal. This intentionally makes wrong
origins sample strong data-plane interference instead of an untouched image. It is not a v4 payload encoder:
there is no frame, ECC, whitening, key or HMAC.

A pilot-only detector then assumes an already-resolved 8/6/4-pixel lattice, aggregates DCT margins by logical
residue and exhaustively scores all 1184 cyclic origins. On the deterministic source regression it selects the
correct origin after native embedding, JPEG quality 82, 3x3 blur, deterministic +/-4 RGB noise, gamma 1.15,
exact 75%/50% resize and an aligned crop. The unmarked negative control has a much smaller best-vs-runner-up
margin; no acceptance threshold is promoted from this experiment.

The two local original development images supplied for Build24 also select the correct origin for every case in
that initial transform matrix. These originals are local qualification inputs and are not packaged in source ZIPs.

### Build24 limit

This is still insufficient to freeze the pilot. Arbitrary rotation, anisotropic scale/shear, mild perspective,
combined transformations, blind geometry estimation and real print-camera/scanner v4 carriers remain required.
Prototype-2 is therefore a stronger candidate, not an interoperability constant.

## Before a normative v4 specification

The next steps must be completed in order:

1. optimize and compare pilot candidates under cyclic, partial-crop and physical-channel simulation;
2. define an explicit v4 version/framing marker that can never be confused with v3;
3. decide ECC separately from pilot design instead of assuming Hamming(7,4) by inertia;
4. implement an experimental encoder and pilot-only detector behind an explicit switch;
5. qualify synthetic geometry and negative controls;
6. create a new v4 physical corpus by embedding/printing the v4 carrier;
7. require HMAC-authenticated physical recovery before any production promotion.



## Build25 known-geometry qualification

Build25 keeps `prototype-2-search-p64` non-normative and tests a prerequisite for blind physical decoding: once a correct canonical-to-observed homography is supplied independently, the pilot must still identify the correct absolute cyclic origin after geometric resampling. The projective pilot sampler reads canonical 8x8 DCT blocks directly through that mapping and reports the same score/runner-up/margin telemetry as the aligned detector.

The fixed matrix has 16 cases covering arbitrary rotation, anisotropic scale, X/Y shear, mild perspective, rotation+scale, rotation+shear, rotated 75% resize, perspective plus JPEG/blur/noise, and perspective+crop with/without JPEG. All synthetic cases recover `(0,0)`. The two local originals contribute 32 additional transformed-carrier cases and all recover `(0,0)` as well. The weakest local positive margin is 0.322069 (`PJ_piccolo`, combined perspective+blur), while that case's unmarked control has margin 0.063984.

A development-only separation floor of 0.10 between positive and matching negative-control margins is enforced to prevent a trivially ambiguous "correct" winner. This is deliberately **not** a normative pilot-detection threshold. Build25 therefore establishes geometric survivability, not blind geometry discovery. The next design gate is a bounded pilot-assisted rotation/affine/perspective search that preserves runner-up evidence and never uses payload/HMAC as an oracle.

## Build26 bounded blind geometry qualification

Build26 removes Build25's supplied-homography assumption for a deliberately bounded, auto-framed experiment. The search follows the intended v4 dependency order rather than asking the pilot to brute-force every projective parameter:

```text
observed image
  -> public repeated-tile geometry proposal (no symbol knowledge)
  -> small bounded candidate bank
  -> public pilot ranking on central repetitions
  -> held-out pilot validation on corner repetitions
  -> absolute cyclic origin
```

### Coarse proposal observable

The current synthetic v4 carrier repeats the same 37x32 data plane spatially. Build26 samples 64 deterministic **non-pilot** residues and compares DCT-sign evidence at homologous positions of adjacent repeated tiles. Only self-consistency is used: the expected data sign is never consulted. Flipping every pilot sign leaves this proposal score bit-for-bit unchanged. The canvas dimensions are allowed only as a weak public prior for scale/inset enumeration.

### Bounded family

The current development bank covers approximately +/-20 degrees rotation, anisotropic scale, X/Y shear up to +/-10 degrees and mild top/bottom perspective insets up to 0.06. Perspective families use a coarse 0.25-degree repeat scan followed by local refinement; pilot evidence is only spent on the repeat-ranked survivors. Current regression cases require <=6400 repeat hypotheses per image.

### Evidence separation

Data-repeat evidence proposes geometry. Central-row pilot samples rank the small proposal bank and optimize sub-block translation. The final winner is rescored on held-out corner pilot repetitions; the repeat score and central pilot score are not reused as the validation objective. Payload/header/CRC/ECC/key/HMAC are absent throughout.

### Build26 result and limit

The synthetic matrix recovers rotation, anisotropic rotation+scale, X/Y shear, mild perspective, combined perspective+anisotropic scale and combined rotation+shear within the declared corner-error bound. Both local originals recover the representative rotate-scale and combined-perspective cases with cyclic origin `(0,0)`.

This is still not a production/physical decoder. Build26 assumes the canonical carrier extent is known and the transform is auto-framed. Arbitrary crop/translation, unknown carrier placement, print-camera acquisition, v4 framing/version identification, ECC and authenticated payload decoding remain open. Prototype-2 therefore remains non-normative.



## Build27 unknown-placement qualification

Build27 keeps `prototype-2-search-p64` unchanged and isolates crop/translation from geometry. With the geometric warp supplied independently, the public pilot recovers non-block-aligned crop offsets and padded-canvas placement on the synthetic matrix and both local development originals. Proposal and validation use disjoint pilot-symbol halves; payload and HMAC remain absent.

Because the v4 tile repeats, absolute translation by complete tile periods is not necessarily observable or necessary. The useful state is a mapping that restores block phase and cyclic tile origin for subsequent data sampling. Build27 therefore does not promote a placement top-2 threshold.

This checkpoint does **not** close blind physical geometry. The next gate is a crop/placement-tolerant coarse lattice/extent estimator followed by the Build27 placement/pilot validation path. Only after that joint synthetic qualification should the project decide whether prototype-2 is ready to freeze before implementing the v4 encoder.


## Build28 joint affine+crop qualification

Build28 combines unknown affine geometry and unknown negative crop placement in one bounded experiment. Geometry selection is now **pilot-symbol independent**: the search scores absolute DCT carrier differential energy and the contrast between sub-block phases, refines angle/X/Y scale hierarchically, and applies a compact coupled refinement only to the final two basins. The public pilot is not exposed until one geometry is fixed.

The selected geometry then enters the Build27 split-pilot placement path. One pilot half proposes translation, the disjoint half validates it, and the complete pilot reports cyclic origin. Two deterministic synthetic cases and representative cases on both local originals recover origin `(0,0)` with matched-negative separation. `PJ_piccolo.png`, which defeated earlier crop-tolerant geometry proposals, now selects the positive ~11.25 degree / 1.0675 / 0.93 basin structurally before placement.

This checkpoint is intentionally affine/crop-only. Projective crop is not equivalent to a constant canonical phase shift, so the affine proposal must not simply be generalized by adding perspective parameters to the same global search. Joint projective recovery, positive padded-canvas placement, shear+placement requalification and unknown physical extent remain required before pilot freeze.

### Portability and planned UI

The implementation remains in Go partly to preserve one portable algorithmic core. `core-target-check` cross-compiles that core for Linux, Windows, Android and iOS targets. A future GUI is planned, especially for mobile use, but it will be a frontend over the same qualified core rather than a separate watermark implementation.

## Build29 joint projective/padded qualification

Build29 does not change the Format-v4 layout or `prototype-2-search-p64`; it studies decoder observability only. The public pilot may participate in geometry proposal because it is format metadata, but proposal and acceptance remain separated: bounded geometry/phase proposal first, Build27 placement validation next, then complete-pilot origin/margin. No authenticated payload evidence is allowed to select geometry.

The important new policy is **safe rejection**. A weakly supported geometry is not considered a decode merely because one local pilot view correlates. Projective acceptance currently requires validation >=0.35, margin >=0.15 and origin `(0,0)`; padded acceptance requires >=0.55, >=0.20 and `(0,0)`. These are development values derived for regression stability, not normative thresholds.

Results show that the present 64-pilot design still has useful signal under the new compositions: one photographic projective case qualifies and a known-geometry padded photographic control is very strong. The remaining failures are therefore search/ranking limitations, not evidence that the pilot pattern must be replaced. The pilot remains non-normative until broader physical and small-carrier evidence exists.



## Build30 pilot candidate lock and normative-freeze boundary

Build30 keeps the 37×32 / 64-pilot / 1120-data layout unchanged and does not alter `prototype-2-search-p64`. The candidate identity is now **development-locked** at SHA-256 `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`. `v4-pilot-lock-check` recomputes this semantic identity and the core cyclic invariants, so later decoder work cannot silently change the public mask/sign sequence.

The lock is supported by a Build30 known-mapping corpus audit. For the same projective-crop and affine-padded fixtures used in Build29, the exact mapping is supplied directly and the public pilot alone is scored. Both development originals retain marked score >= 0.928412, marked margin >= 0.484086 and origin `(0,0)`; matching unmarked controls remain <= 0.010553 margin. This includes Build29 SAFE-REJECT cases and therefore separates **pilot survivability** from **geometry/placement search coverage**.

This is not yet a normative Format-v4 pilot declaration. The missing gate is physical v4 print-camera/scanner evidence from an actual v4 encoder. Until that exists, the lock prevents accidental mutation but does not create an interoperability promise. Decoder search, confidence thresholds, framing/version marker, data mapping, ECC, whitening/interleaving, payload layout and authentication framing remain experimental. See `docs/V4_PILOT_CANDIDATE_LOCK.md` for the exact contract and mutation policy.


## Build31 first authenticated v4 framing candidate

Build31 is the first checkpoint that writes the locked public pilot **and** a real authenticated payload into the same carrier. It deliberately does not redesign ECC at the same time. The three profile frame sizes remain 32/48/80 bytes and Hamming(7,4) remains the protected-bit baseline so pilot/framing behavior can be compared directly with the frozen v3 arithmetic. This is not a final ECC decision.

The public pre-key marker is the development-locked `prototype-2-search-p64` pilot. Once data is sampled with a key, the authenticated frame distinguishes v4 and profile with bytes `0x41`, `0x42`, `0x43`. Whitening uses `pixseal-whiten-v4`; HMAC uses the domain `pixseal-frame-v4 || 0x00` before the authenticated frame prefix. This prevents a v4 protected bitstream from being interpreted as a v3 frame under the same key.

The exact 64 prototype-2 pilot positions are removed from the 1184-position tile. The remaining positions, in ascending row-major order, are data ordinals 0..1119. For coded length `C`, ordinal `d` maps to `(251*d) mod C`. This gives 2/3 observations per robust coded bit, 1/2 per balanced bit and exactly one per capacity bit. A Build31 regression covers the full 64+1120 partition so a historical prototype-1 helper cannot accidentally leak a pilot block into the data plane.

The first aligned decoder follows the intended dependency order:

```text
public pilot -> cyclic origin -> data-only aggregation -> Hamming -> v4 dewhitening -> version/profile parse -> HMAC -> CRC consistency
```

Only HMAC produces success. Build31 does not yet connect Build29 blind projective/placement search to this authenticated decoder. The purpose of `v4-embed` is now also practical: it can create the separately embedded v4 images needed for the first real print-camera/scanner corpus. See `docs/V4_BUILD31_FRAME.md` for the exact experimental layout and deterministic vector.
