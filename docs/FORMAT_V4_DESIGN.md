# Format v4 design study — absolute pilot branch

Status: **experimental, non-normative design branch**. Sizing was first quantified in v0.3.0-build22; build23 adds the first isolated prototype pilot foundation.
Format v3 remains frozen and is still the only implemented/interoperable format.

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

The current winner is **non-normative** `prototype-2-search-p64`:

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

