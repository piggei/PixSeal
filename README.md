# PixSeal

![PixSeal — Hide messages. Keep the picture.](docs/assets/pixseal-banner.png)

**Hide messages. Keep the picture.**

PixSeal is an experimental, pure-Go **robust image steganography** tool for
hiding short authenticated messages inside images while preserving their visual
appearance as much as possible. It embeds protected payload bits in luminance
DCT coefficients and is designed to investigate how well hidden data can survive
real-world transformations such as resizing, JPEG recompression, printing,
scanning and smartphone photography.

PixSeal is a hidden-data channel, not an ownership-marking product. Digital
watermarking is the robustness mechanism; the project goal is robust
steganography for short messages. PixSeal does **not** claim statistical
undetectability and has not undergone a professional cryptographic or
steganalytic audit.

The project is **source-available for noncommercial use** under the PolyForm
Noncommercial License 1.0.0 (`PolyForm-Noncommercial-1.0.0`). Commercial use
requires a separate commercial license. See
[`docs/LICENSING.md`](docs/LICENSING.md), [`LICENSE`](LICENSE) and
[`NOTICE`](NOTICE). Historical copies already released under MIT retain the MIT
rights granted with those copies; the licensing change is prospective.

## Project status and development

Current development snapshot: **v0.3.0-build49**.

Latest qualified milestone: **v0.3.0-build44**. Build49 is research-only and does not change the qualified production decoder.

Stable release baseline: **v0.2.0**.

> **Build44 qualified milestone:** JPEG input now uses a project-controlled pure-Go pre-Go-1.26 decoder, eliminating the accidental dependency on the compiler standard library without retuning geometry. **Go 1.26.0 is now the qualified Build44 toolchain**: the deterministic-raster regression and the complete private Build43 smartphone matrix both pass unchanged. See [`docs/V4_BUILD44_DETERMINISTIC_JPEG.md`](docs/V4_BUILD44_DETERMINISTIC_JPEG.md) and [`docs/GO_TOOLCHAIN_COMPATIBILITY.md`](docs/GO_TOOLCHAIN_COMPATIBILITY.md).

> **Build45 completed diagnostic checkpoint:** Build44 remains the latest qualified milestone. Build45 changes no production decoder decision; it decomposes the remaining `B/mild` failure into blind geometry, held-out qualification or post-geometry data-channel stages, and adds an explicitly lab-only supplied-geometry oracle. Retained-corpus oracle evidence already authenticates both `B/mild` and `B/angle` exactly under independently supplied geometry, so the next justified work is blind geometry rather than strength/ECC/HMAC changes. See [`docs/V4_BUILD45_PHONE_FAILURE_DECOMPOSITION.md`](docs/V4_BUILD45_PHONE_FAILURE_DECOMPOSITION.md).

> **Build46 research snapshot:** the qualification-host Build45 run shows 32/28 frozen Build43 candidates but only one held-out-qualified candidate for `B/mild`/`B/angle`, while the Build42 bank remains empty because production requires a two-geometry direct ensemble before promoting the Build43 bank. Build46 therefore inspects that already-qualified singleton without changing production thresholds: it tests post-qualification single-candidate channel viability and, only afterwards, compares its corners to the isolated reference oracle. See [`docs/V4_BUILD46_QUALIFIED_HANDOFF.md`](docs/V4_BUILD46_QUALIFIED_HANDOFF.md).

Format v3 remains the implemented interoperability baseline. Its on-image layout,
deterministic encoder fingerprints and production decoder are **frozen**. After builds
16–22, the v0.3 absolute-cycle research on v3 is considered **closed**: the available
public structure is measurable but not sufficiently discriminative on the physical
negative control to justify another threshold or voting round. Future v3 work is
maintenance/regression only unless genuinely new independent evidence appears.

Current Format-v4 development is focused on blind physical recovery from smartphone
photographs. Build39 introduced the global projective phone path, Build40 added a
strictly bounded public-pilot-only residual field, Build41 closed the first blind A/B
basin milestone, Build42 improved only post-geometry data recovery, and Build43 added
a bounded side-pair geometry fallback for A/front. The qualified strength-48 matrix is
controls 3/3 reject; `A/front`, `A/mild`, `A/angle` and `B/front` authenticate;
`B/mild` and `B/angle` remain informational rejects. Build44 changes only JPEG ingest
so that this matrix can be reproduced across Go toolchains. Build45 is diagnostic-only:
it records where the unchanged phone path stops and uses an isolated reference-assisted
oracle to prove both difficult B captures remain inside the protected data-channel envelope
when geometry is supplied independently. The qualification-host blind run further shows one
held-out-qualified Build43 candidate in each difficult B capture, but not the two candidates
required by the production ensemble. Build46 measures whether that qualified singleton is
already data-viable and how far it lies from the oracle, without changing any production
quorum. Protected payload contents, secret key and HMAC remain excluded from production
geometry search/ranking; HMAC is only the final frame authenticator.

> **Build42 post-geometry milestone:** Build41 geometry is unchanged. The normal two-hypothesis soft-Hamming decode runs first. Only when already-qualified geometry reaches data decode but fails authentication does Build42 retain the complete frozen Build41-qualified bank, enumerate deterministic three-geometry data ensembles, and apply a bounded Hamming list decoder ordered only by ML score gaps. `A/mild` now authenticates `v4-b38-phone-a`; `A/angle` and `B/front` remain direct Build41 passes; all three controls still reject before data decode. A/front, B/mild and B/angle remain geometry-reject research cases. See [`docs/V4_BUILD42_PHONE_DATA_LIST.md`](docs/V4_BUILD42_PHONE_DATA_LIST.md).
>
> **Build32 active private corpus:** qualification uses only the explicit LQ/MQ/HQ entries in `private-corpus-active.tsv`; extra files in `original pics/` are ignored. New v4 development tests/fixtures use the public test key `PixSeal-v4-TestKey-2026`. Historical `Piccotti` vectors remain only for reproducibility. The corpus images themselves are private and are not shipped.

Project history and future work are kept in [`HISTORY.md`](HISTORY.md) and
[`TODO.md`](TODO.md). Release-facing changes are in [`CHANGELOG.md`](CHANGELOG.md).
The append-only research notebook in [`docs/RESEARCH_LOG.md`](docs/RESEARCH_LOG.md)
records hypotheses, rejected variants, threshold decisions and negative results so
future builds do not silently repeat abandoned experiments.


## What v0.3.0-build49 adds

Build49 is a diagnostic-only **proposal-ranking observability** study. The Build48 qualification-host run showed that B/mild local refinement selected 12 seeds, increased held-out-qualified candidates from 4 to 7, but authenticated none; its best selected geometry improved only 46.202 -> 44.436 px. Build47 had already exposed a better 34.049 px all-pair candidate, proving that the top-2-per-pair proposal selector discarded the best known basin before refinement.

Build49 therefore changes no geometry. It freezes the corrected Build47 extended bank and measures proposal fold agreement, per-tile stability, side-pair/cell evidence and within-pair ranks for every candidate. The CLI command accepts no secret key. Only after both blind outputs are frozen does the private lab script compute SIFT/reference error, top-2/4/6/8-per-pair coverage and post-hoc rank/error correlations. See `docs/V4_BUILD49_PROPOSAL_RANKING_OBSERVABILITY.md`.

## What v0.3.0-build47 adds

Build47 adds diagnostic-only proposal-breadth observability while keeping Build43 production unchanged. Stage 32 is the exact production proposal tier; stage 64 adds deeper cells on the same two selected side pairs; stage 128 also adds lower-ranked side pairs. This distinguishes cell-depth pruning from side-pair pruning and from a genuinely absent basin. Oracle geometry remains post-hoc only. See `docs/V4_BUILD47_FROZEN_BANK_OBSERVABILITY.md`.

## Build46 handoff diagnostic

Build46 is a diagnostic handoff checkpoint. It does **not** change Build43 candidate generation, held-out qualification, the two-geometry production ensemble, Build42's minimum three-geometry bank, the residual fallback, ECC or HMAC. The new `v4-diagnose-phone-handoff` command reruns only the public-evidence Build41/43 search, enumerates the geometries that have already passed held-out qualification, and records the exact production quorum that prevents or permits each downstream path.

For each qualified Build43 candidate Build46 performs one explicitly diagnostic post-qualification decode using the unchanged Build42 soft-Hamming list. HMAC is observed only after that candidate is already frozen and qualified; it cannot create, move, rank or qualify geometry. An optional oracle quadrilateral is likewise comparison-only and is used after blind search to report corner error.

Use `make v4-build46-phone-handoff-test` for the source regression and `make v4-build46-phone-handoff-diagnostic` for the private B/mild/B/angle study. See [`docs/V4_BUILD46_QUALIFIED_HANDOFF.md`](docs/V4_BUILD46_QUALIFIED_HANDOFF.md).

## What v0.3.0-build45 adds

Build45 is a diagnostic research checkpoint, not a new production decoder. It keeps the qualified Build44 deterministic JPEG ingest and the complete Build41/42/43 smartphone path unchanged, while adding `v4-diagnose-phone` and private study targets that expose Build41 qualified-bank state, the complete six-pair Build43 proposal ranking, frozen/held-out-qualified counts, Build42 list activity and final HMAC result. A failed acquisition is classified as `geometry`, `qualification` or `data-channel`; a normal HMAC success is `recovered`.

An optional lab-only SIFT/RANSAC helper may register the known digital marked-B carrier to an acquisition and supply exactly one quadrilateral to `v4-diagnose-phone -oracle-only`. That geometry cannot enter `v4-extract-phone`, cannot be used to rank production candidates and introduces no OpenCV runtime dependency into PixSeal. On retained private B/mild and B/angle captures, the supplied geometry yields origin `(0,0)` and exact `v4-b38-phone-b` authentication on the first soft-Hamming list frame. This is oracle evidence only, but it establishes that the protected strength-48 channel is sufficient in both difficult B images.

The completed qualification-host blind run freezes 32 candidates for B/mild and 28 for B/angle; exactly one candidate in each image survives Build43 held-out qualification. The Build42 bank remains zero because the production Build43 fallback requires a two-geometry ensemble before promotion, while Build42 independently requires three qualified geometries. Build46 refines this singleton state rather than treating it as a generic qualification failure.

Use `make v4-build45-phone-diagnostic` for the blind decomposition, `make v4-build45-phone-oracle-diagnostic` for the isolated reference-assisted control, or `make v4-build45-phone-study` for both. See [`docs/V4_BUILD45_PHONE_FAILURE_DECOMPOSITION.md`](docs/V4_BUILD45_PHONE_FAILURE_DECOMPOSITION.md).

## What v0.3.0-build44 adds

Build44 makes the JPEG raster part of the PixSeal implementation contract instead of inheriting it from the Go compiler version. The CLI now identifies JPEG/PNG by magic bytes; JPEG `DecodeConfig` and full decode use the pure-Go `internal/jpeglegacy` decoder, while PNG remains on `image/png`. The decoder is derived from the pre-Go-1.26 Go standard-library implementation and retains its BSD-style upstream license.

A small public JPEG fixture locks Y/Cb/Cr hashes and a CLI regression proves that normal image opening uses the deterministic decoder. Phone diagnostics expose `input-decoder: pixseal-jpeg-pre-go1.26-v1`. No Format-v3/v4 encoding, pilot, strength, ECC, whitening/HMAC, geometry, data/list decoding or HMAC semantics change.

Go 1.26.0 is the qualified Build44 build/test toolchain. `make v4-build44-jpeg-compat-test` protects the deterministic raster contract and `make v4-build44-phone-physical-test` verifies the unchanged private Build43 physical matrix. The historical `v4-build44-go126-*` targets remain available as explicit qualification regressions.

See [`docs/V4_BUILD44_DETERMINISTIC_JPEG.md`](docs/V4_BUILD44_DETERMINISTIC_JPEG.md).

## What v0.3.0-build43 adds

Build43 adds a bounded **blind side-pair geometry fallback** for smartphone Format-v4 recovery. It runs only after the qualified Build41 geometry path rejects, uses visible image structure plus the public pilot for proposal/refinement, freezes a maximum 32-candidate bank before held-out evaluation, and reuses the unchanged Build42 data/list decoder and HMAC authentication.

The private Build38 strength-48 milestone is now: all three controls reject; `A/front` authenticates through Build43; `A/mild`, `A/angle` and `B/front` retain their prior authenticated paths; `B/mild` and `B/angle` remain research/INFO rejects.

Use `make v4-build43-phone-side-pair-test` for the public regression and opt-in `make v4-build43-phone-physical-test` for the private nine-photo matrix. See [`docs/V4_BUILD43_PHONE_SIDE_PAIR.md`](docs/V4_BUILD43_PHONE_SIDE_PAIR.md).

## What v0.3.0-build42 adds

Build42 deliberately does **not** reopen smartphone geometry. The Build41 boundary, fixed 2/3 proposal versus 1/3 held-out split, projective qualification floors and two-geometry direct decoder remain the first path. Encoder strength 48, pilot identity, 1120-position data layout, Hamming(7,4), whitening/HMAC domains and Format-v3 are unchanged.

When Build41 has already accepted geometry but its normal two-hypothesis soft decoder fails HMAC, Build42 keeps the complete geometry bank that was frozen and qualified by exactly the same public-pilot rules. At most six already-qualified geometries are retained. All deterministic three-member data ensembles are formed without key or payload feedback, their protected DCT margins are averaged, and the normal soft-Hamming ML result is attempted first. If that frame does not authenticate, Build42 considers only the second-best nibble for the ten Hamming words with the smallest ML score gaps, producing a tightly bounded list. HMAC validates complete frames only; it does not create, refine or rank geometry.

The decoder order is now `Build41 global decode -> Build42 qualified-bank/list decode -> Build40 residual fallback`. This avoids spending minutes fitting a residual field when the accepted global geometry already contains enough data evidence. On the private strength-48 corpus, `A/angle` and `B/front` still authenticate directly without entering Build42; `A/mild` now authenticates `v4-b38-phone-a` through the new data fallback; all three controls reject before data decode. `A/front`, `B/mild` and `B/angle` remain geometry rejects and are explicitly deferred rather than hidden by weaker thresholds.

Use `make v4-build42-phone-data-test` for the public synthetic/list-decoder gate and opt-in `make v4-build42-phone-physical-test` for the private nine-photo matrix. See [`docs/V4_BUILD42_PHONE_DATA_LIST.md`](docs/V4_BUILD42_PHONE_DATA_LIST.md).

## What v0.3.0-build41 adds

Build41 is the first blind smartphone checkpoint to close the staged physical A/B milestone on the existing strength-48 corpus. It changes only the global phone registration front end; the Format-v4 encoder, locked pilot, data layout, strength 48, Hamming(7,4), whitening/HMAC domains, Build40 residual layer and frozen Format-v3 core are unchanged.

A structure-only artwork boundary anchors absolute phase. The canonical-to-phone quadrilateral is then refined with eight bounded corner coordinates using a fixed 3-way spatial split of the public pilot: folds 1+2 are the only proposal evidence, while fold 0 is never sampled until the complete proposal shortlist has been frozen. Shape fitting may ignore cyclic origin while finding the local basin; a bounded proposal-only translation restores physical origin `(0,0)`. The strongest phase candidates receive a small proposal-only corner polish before freeze. Held-out pilot evidence can only accept/reject or rank frozen candidates; it cannot generate new geometry.

The final phone ensemble contains two pilot-qualified geometries separated by at least `max(0.75 px, 0.20 * observed DCT-block scale)`. This diversity requirement prevents two near-identical local optima from masquerading as an ensemble and brackets the remaining sub-pixel registration uncertainty. Only after the pair is frozen does the unchanged Build40 residual/Build36 soft-Hamming/HMAC path run.

`make v4-build41-phone-basin-test` provides a public synthetic end-to-end gate, including an unmarked control. The opt-in private `make v4-build41-phone-physical-test` runs the nine original strength-48 photos and applies the staged milestone: all three controls must reject and at least one A plus one B capture must authenticate the exact expected payload. Build41 checkpoint result: `A/angle` and `B/front` HMAC PASS; `A/mild` reaches qualified geometry/data decode but does not authenticate; A/front, B/mild and B/angle remain outside the accepted basin. See [`docs/V4_BUILD41_PHONE_BASIN.md`](docs/V4_BUILD41_PHONE_BASIN.md).

## What v0.3.0-build40 adds

Build40 tests the next registration hypothesis without changing the Format-v4 carrier. The encoder, locked pilot, 1120 data positions, robust strength 48, Hamming(7,4), whitening/HMAC domains and frozen Format-v3 implementation are unchanged. The existing nine-photo strength-48 corpus is reused; no new printing or acquisition is required.

After Build39 has produced a canonical-to-observed homography, Build40 may fit a six-coefficient quadratic displacement field independently for X and Y. The field is intentionally a **residual** model: corrections are capped at six canonical pixels and cannot reopen an arbitrary projective or whole-tile search. Local controls are measured only on checkerboard-A repetitions of the public pilot. Ambiguous local peaks are rejected, the remaining controls are robustly fitted, and checkerboard-B repetitions stay completely held out until acceptance. A residual field must improve held-out validation and then recover complete-pilot origin `(0,0)` with explicit score/margin floors before protected data may be sampled. HMAC is still only the final verifier.

The new `v4-build40-phone-residual-test` creates a synthetic strength-48 carrier with a known smooth non-projective deformation. The uncorrected carrier remains geometrically imperfect; the pilot-only field is recovered, held-out validation improves, soft-Hamming data sampling through the warped mapper recovers the exact HMAC-authenticated payload, and an unmarked control does not qualify. The residual sampler is integrated into `ExperimentalV4ExtractPhone` / `v4-extract-phone`, with automatic fallback to the untouched Build39 ensemble if a residual-qualified attempt does not authenticate.

The Build40 private real-phone result was intentionally negative but informative: the then-current Build39 basins did not produce a coherent <=6 px residual field. Candidate fits could improve proposal tiles while held-out validation fell and the complete pilot retained a non-zero origin, so Build40 **safe-rejected** those fields rather than overfitting the image. That result directly motivated Build41 global basin recovery. See [`docs/V4_BUILD40_PHONE_RESIDUAL.md`](docs/V4_BUILD40_PHONE_RESIDUAL.md).

Before changing the Build39/40 algorithm, use the opt-in private-corpus target `make v4-build40-phone-corpus-diagnostic`. It runs the unchanged `v4-extract-phone` path over the canonical nine strength-48 captures and writes a TSV plus Markdown matrix under `v4-phone private/build40-diagnostics/`. The matrix records boundary detection, presence of a Build39 projective basin, accepted ensemble size, proposal/held-out pilot scores, pilot origin, residual fit/application and RMS, held-out delta, soft-Hamming profile attempts, maximum data confidence, HMAC/fallback result and exact-payload qualification. Per-image stderr logs are retained; image bytes are never copied into the report directory or release artifacts. This diagnostic target is intentionally excluded from `all-test` because the corpus is private.

The Build40 diagnostic harness remains available as a historical pre-Build41 matrix tool. Build41 adds a separate staged private qualification target and does not rewrite the Build40 measurements.

## What v0.3.0-build39 adds

Build39 consumes the completed strength-48 smartphone acquisition rather than changing the carrier again. The private nine-photo set remains control / marked A / marked B x front / mild / angle at native ~200 MP. Reference-assisted registration now proves that strength 48 is sufficient: representative marked captures reach the authenticated v4 frame, while A/mild is within one post-soft-Hamming bit. The blocker is therefore registration, not another strength/ECC change.

A new experimental `v4-extract-phone` API/CLI path performs bounded internal downsampling (working long side <= 4600 px), rejects scanner-style page/background fits that touch the camera frame, estimates the printed artwork quadrilateral with local-paper and edge-gradient evidence, and refines a canonical-to-observed homography with the public pilot. Proposal and validation use spatially disjoint checkerboard tile partitions so payload/header/ECC/HMAC remain unavailable to geometry selection. A small pilot-qualified geometry ensemble is then eligible for the unchanged Build36 soft-Hamming/HMAC decoder.

The new `v4-build39-phone-registration-test` is deliberately a **checkpoint gate**: it qualifies the downsample bound, perspective boundary extraction, disjoint-pilot registration evidence and negative-control separation on synthetic phone geometry. It does not claim blind HMAC closure on the private real-photo corpus. Current real captures show that a single projective mapping can still leave local residual misregistration (camera optics / print flatness / resampling), so the next research step is a bounded pilot-only residual warp after projective registration. Format v3 remains frozen and Build37 scanner recovery remains unchanged.

## What v0.3.0-build38 adds

Build38 begins the smartphone-specific physical channel without altering the v4 framing, pilot, data mapping, Hamming code, HMAC domains or frozen Format v3. The first private phone corpus contains nine native ~200 MP JPEGs from the existing Build35 strength-24 paper set: control / marked A / marked B, each photographed frontally, at mild perspective and at stronger perspective.

That corpus is an important negative measurement. Build37's scanner boundary assumptions do not transfer directly to free-camera framing, but even reference-assisted diagnostics that supply the digital original for geometry cannot recover the strength-24 payloads. With dense non-rigid residual registration, representative captures still show about 14.5–17.9% protected coded-bit error and 20–42 wrong bits after soft Hamming. A pilot-only local warp can overfit natural image texture without improving protected data, so Build38 explicitly rejects that as a production shortcut.

The next physical experiment therefore changes only embedding strength. `make v4-phone-fixtures` produces a private MQ pack at robust strength **48**, retaining the same 1632x1632 canonical carrier, 300-ppi print geometry, locked pilot and v4 frame. It emits control plus `v4-b38-phone-a` / `v4-b38-phone-b`, a SHA-256 manifest, acquisition plan and instructions under `v4-phone private/build38-generated/`. The new `v4-build38-phone-channel-test` protects aligned/JPEG decoding at the chosen strength. On the canonical MQ image, strength 48 measures about 32.5 dB PSNR against the control and remains visually subtle.

The successful Build35/37 scanner corpus is preserved unchanged. The original strength-24 phone photographs remain a historical private corpus and must not be overwritten. A blind phone decoder will be developed only after the strength-48 photographs prove that the physical data channel itself is inside the Hamming/HMAC recovery envelope. See [`docs/V4_BUILD38_PHONE_QUALIFICATION.md`](docs/V4_BUILD38_PHONE_QUALIFICATION.md).

## What v0.3.0-build37 adds

Build37 closes the first fully blind Format-v4 **print -> paper -> scanner -> JPEG -> authenticated payload** checkpoint on the existing Build35 physical corpus. No carrier is reprinted, embedding strength remains 24, the locked `prototype-2-search-p64` pilot and 37x32/64+1120 tile are unchanged, Hamming(7,4) remains the current experimental ECC, and frozen Format v3 is untouched.

The `ExperimentalV4ExtractScanner` API and `v4-extract-scanner` CLI expect a full-page scan with visible white paper around the printed artwork plus the canonical pre-print dimensions. The visible paper/artwork boundary is used only as an independent geometric initializer. A deliberately narrow scanner-affine bank is proposed with pilot partition A; disjoint pilot partition B ranks the surviving basin. The held-out winner must also win a complete-pilot cyclic-origin competition at `(0,0)` with explicit score/margin floors. The five best pilot-qualified geometries are then frozen as an ensemble. Only after that freeze are protected data margins sampled, validation-weighted across the ensemble, soft-Hamming decoded and authenticated by the unchanged v4 frame/HMAC. Payload/header/CRC/key/HMAC evidence is unavailable to boundary detection, geometry search and pilot acceptance.

On the first private scanner corpus (actual-size 300-ppi print, smooth low-porosity white paper, office scanner, 600-dpi color JPEG), Build37 produces: control REJECT; `v4-b35-phys-a` PASS; `v4-b35-phys-b` PASS.

## What v0.3.0-build36 adds

Build36 keeps the Build34/35 blind geometry gate unchanged and improves only the **post-geometry data decoder**. The accepted projective path now retains signed DCT evidence for every protected bit and performs deterministic maximum-likelihood Hamming(7,4) decoding before the legacy hard-decision fallback. This is specifically targeted at print/scan/JPEG cases where two weak sign errors can occur in one Hamming word. HMAC remains the only payload success criterion and no frame/key/HMAC evidence participates in geometry proposal or acceptance.

The first Build35 scanner corpus has now been acquired. The available office scanner produced 600-dpi color JPEG rather than the planned 300-dpi lossless PNG. Blind Build35 geometry still SAFE-REJECTs the two marked scans, so this is **not yet a normative physical PASS**. A private reference-assisted geometry diagnostic, used only to isolate channel capacity, shows that the same physical scans authenticate both `v4-b35-phys-a` and `v4-b35-phys-b` exactly when Build36 soft Hamming is applied. This proves that the current strength-24 encoder, locked pilot, data mapping and existing Hamming code survive the print/scan channel; the remaining blocker is blind scanner registration. See [`docs/V4_BUILD36_PHYSICAL_CHANNEL.md`](docs/V4_BUILD36_PHYSICAL_CHANNEL.md).

## What v0.3.0-build35 adds

Build35 is the controlled handoff to physical Format-v4 qualification. The Build34 geometry/HMAC algorithm is intentionally unchanged. The preceding checkpoint is consistently Build34 in code, test names, Make targets and documentation.

The new exported `ExperimentalV4ExtractProjective` API and `v4-extract-projective` CLI expose the Build34 blind projective+crop path for physical captures. The caller supplies the **block-aligned canonical dimensions of the digital carrier before printing**; geometry remains selected exclusively from public structural/pilot evidence, and Hamming/frame/CRC/HMAC processing begins only after the unchanged Build29 acceptance gate.

`make v4-physical-fixtures` now creates a private scanner-first MQ qualification pack: one unmarked control plus two robust/strength-24 carriers with different payloads (`v4-b35-phys-a`, `v4-b35-phys-b`), an SHA-256 manifest and an acquisition plan. `make v4-physical-qualification` requires exact HMAC recovery of both marked payloads and rejection of the control. The initial protocol is actual-size 300-ppi printing followed by 300-dpi lossless scanning cropped to the artwork edges. Smartphone photographs are deliberately deferred until this controlled print/scan gate is measured. See [`docs/V4_BUILD35_PHYSICAL_QUALIFICATION.md`](docs/V4_BUILD35_PHYSICAL_QUALIFICATION.md).

## What v0.3.0-build34 adds

Build34 resolves the MQ blind projective+crop ranking blocker identified by Build33. For carriers with at least 4x4 complete v4 tiles, the decoder now preserves two distinct structural anchors, explores a bounded coupled 243-point local geometry neighborhood around each anchor, retains the 64 strongest local structural candidates, and ranks them with the public pilot using centered translation phases at 0, +/-2 and +/-4 pixels. All complete interior tile rows participate in proposal; first/last tile rows remain spatially held out until after one geometry is committed.

A non-zero cyclic pilot origin is not accepted by fiat. Build34 composes the detected block-cycle offset into the canonical side of the homography and re-runs the full pilot detector, restoring canonical origin `(0,0)` before acceptance. The original Build29 projective validation >=0.35 and pilot-margin >=0.15 floors remain unchanged. Small carriers continue through the previous Build29 path, so the active LQ corpus remains a SAFE REJECT.

After geometry is accepted, Build34 samples the locked 1120-position v4 data partition through that homography, then performs the existing Hamming decode, whitening, frame parsing, CRC and HMAC verification. Payload/header/key/CRC/HMAC evidence never participates in geometry ranking. The Build34 corpus gate verifies exact authenticated recovery for two different robust payloads on MQ.

## What v0.3.0-build31 adds

Build31 is the first checkpoint that emits a **real authenticated Format-v4 carrier**. The stable v3 APIs and CLI remain unchanged; new functionality is exposed only through `ExperimentalV4EmbedWithInfo`, `ExperimentalV4ExtractAligned`, `v4-embed` and `v4-extract`. The latter intentionally assumes the native 8-pixel lattice is already aligned, so Build31 qualifies framing/authentication without prematurely coupling it to the still-experimental Build29 geometry search.

The Build30 locked pilot is unchanged: 37x32 blocks, 64 public pilot positions, 1120 data positions and SHA-256 `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`. Data ordinals are the ascending row-major non-pilot positions; ordinal `d` maps to protected bit `(d*251) mod codedBits`. The three current v4 profiles intentionally retain the v3-sized 32/48/80-byte frames and Hamming(7,4) protected lengths 448/672/1120 for an apples-to-apples baseline. This is **not** a final ECC decision.

The authenticated frame uses `PS`, version/profile bytes `0x41/0x42/0x43`, payload length, CRC32, payload and an 8-byte truncated HMAC-SHA256. Whitening is domain-separated as `pixseal-whiten-v4`; HMAC input is domain-separated with `pixseal-frame-v4`. HMAC is checked before a payload is accepted; pilot confidence, ECC, CRC and header plausibility remain non-authenticating evidence. A deterministic Build31 vector and explicit v3/v4 cross-version rejection regressions protect this candidate framing from accidental mutation. See [`docs/V4_BUILD31_FRAME.md`](docs/V4_BUILD31_FRAME.md).

On both private photographic development originals, robust/balanced/capacity native round-trips authenticate. Robust also authenticates after JPEG quality 82 and after a block-aligned crop with the public pilot recovering cyclic origin. The synthetic minimum-size 296x256 carrier authenticates a capacity-profile payload. At the Build31 checkpoint this supplied the encoder needed to create a **new physical v4 corpus**, while physical print-camera/scanner evidence and a general geometry+payload decoder were still absent. Build34 later connects blind projective geometry to authenticated frame recovery; Build35 turns that path into the controlled physical-qualification interface.

The Build31 qualification matrix grew from 47 to **49 targets** with `v4-frame-test` and `v4-frame-corpus-test`. Its physical requirement remains unchanged: a real paper acquisition must yield a valid v4 HMAC before the locked pilot can be considered normative.

## What v0.3.0-build30 adds

Build30 performs a **pilot freeze-readiness audit** rather than widening the geometry search or loosening Build29 SAFE-REJECT gates. The exact Build24 winner `prototype-2-search-p64` is now protected by a semantic development lock: tile geometry, ordered 64-position mask, ordered 32+/32- sign sequence and candidate identity SHA-256 `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174` must remain unchanged. `make v4-pilot-lock-check` fails if that identity or its core structural invariants change. See [`docs/V4_PILOT_CANDIDATE_LOCK.md`](docs/V4_PILOT_CANDIDATE_LOCK.md).

The new local-corpus audit deliberately bypasses joint geometry/placement ranking and supplies the exact mapping used to create the Build29 projective-crop and affine-padded fixtures. This asks a narrower question: **when the mapping is correct, is the public pilot itself still strong in cases that Build29 may SAFE REJECT?** Across both private photographic development originals and both transform families, marked score remains at least 0.928, marked margin at least 0.484, origin is always `(0,0)`, and the largest matching unmarked-control margin is 0.0106. In particular, the small padded case that the Build29 placement/joint search cannot recover still gives roughly 0.932 score / 0.484 margin under the true mapping. The remaining failures are therefore search/placement limitations, not evidence that the pilot pattern should be replaced.

Build30 **does not promote the candidate to a normative Format-v4 pilot yet**. PixSeal still lacks physical v4 print-camera/scanner evidence produced by a real v4 encoder, and no v4 framing/data/ECC/HMAC codec exists. The candidate is locked so continued decoder research cannot accidentally mutate it; normative promotion requires physical-channel qualification. If that future evidence exposes a pilot-level weakness, a replacement must receive a new identity instead of silently editing the locked candidate.

Build30 also records a rejected improvement attempt for the photographic padded SAFE-REJECT cases: scanning all 64 sub-block DCT phases removes the even-phase bias but remains dominated by natural photographic texture, so it is not adopted. The Build29 acceptance floors and SAFE-REJECT behavior remain unchanged. The default qualification matrix grows from 45 to **47 targets** with the candidate-lock identity check and known-mapping corpus audit.

## What v0.3.0-build29 adds

Build29 extends the joint v4 research in two directions that Build28 intentionally left open: **mild projective geometry under unknown crop** and **affine geometry with positive padded-canvas placement**. The search remains bounded to the development domain (`+/-12 degrees`, bounded anisotropic scale and projective inset banks) and never consults payload, header, CRC, ECC, key material or HMAC. Public pilot evidence is allowed for proposal/ranking because the pilot is intentionally public; final acceptance additionally requires a complete-pilot absolute-origin/margin gate.

The projective branch uses a public dimensional feasibility filter, structural DCT proposal, a two-half pilot consistency prefilter evaluated on the same phase/origin hypothesis, basin diversity, local refinement and then the already-qualified Build27 placement search on only a few finalists. The padded branch uses an analogous affine bank. Build29 introduces explicit development acceptance gates: projective requires placement validation >= 0.35, full-pilot margin >= 0.15 and origin `(0,0)`; padded requires validation >= 0.55, margin >= 0.20 and origin `(0,0)`. These are regression gates, not normative decoder thresholds.

The synthetic projective+crop and affine+padded cases are accepted with low corner error. On the private photographic development corpus, `PJ_lingua.PNG` projective+crop is accepted, while `PJ_piccolo.png` is deliberately **SAFE REJECTED** because its full-pilot margin remains below the projective floor. The joint padded search is also SAFE REJECTED on both photographs; for `PJ_lingua.PNG`, a known-geometry control still reaches very strong placement/pilot scores, proving that the remaining limitation is the joint geometry ranking rather than loss of the public-pilot channel. Safe rejection is therefore part of the qualified behavior, not a failure to be hidden by looser thresholds.

Build29 still does not freeze `prototype-2-search-p64`, does not implement a v4 payload/HMAC path, and does not claim general projective or padded-canvas recovery on all carrier sizes. The next gate is improving the photographic joint padded/projective ranking while preserving the current rejection behavior.

## What v0.3.0-build28 adds

Build28 combines the two problems that Build26 and Build27 qualified separately for the **affine crop** family. The observed image is an arbitrary non-block-aligned crop of a carrier with unknown rotation and unknown anisotropic X/Y scale. Neither geometry nor crop translation is supplied to the search.

The geometry stage is deliberately independent of pilot symbols and payload content. It measures the absolute DCT carrier observable `|DCT(2,3)-DCT(3,2)|` over many blocks and scores how sharply that energy concentrates at one 8-pixel sampling phase. A hierarchical bank searches a bounded `+/-12 degree` rotation range and anisotropic scales, then a compact **coupled** angle/scale refinement is applied only to the final structural basins. This coupled final stage was added because coordinate-wise refinement could remain trapped at a nearby scale on `PJ_piccolo.png` even though the correct basin was already present.

Only after one geometry has been selected structurally does the already-qualified Build27 placement search run. Pilot half A proposes crop translation; disjoint pilot half B supplies held-out placement validation; the full public pilot then reports cyclic origin and margin. Header, payload, CRC, ECC, secret key and HMAC are absent from geometry and placement selection. On the synthetic matrix and both local development originals, the joint search restores origin `(0,0)` with clear matched-negative separation.

Build28 intentionally qualifies **rotation + anisotropic scale + negative crop/translation** only. Joint shear, joint projective/perspective recovery and joint padded-canvas placement remain open even though those components are already qualified separately in Build26/27. `prototype-2-search-p64` therefore remains non-normative and no v4 payload encoder/decoder is enabled yet.

## What v0.3.0-build27 adds

Build27 isolates **unknown carrier placement** before attempting the harder joint geometry+placement problem. The experiment receives the canonical-to-full-frame rotation/scale/shear/projective mapping, but it is not told where that transformed carrier lies in the observed image. The observed image may be an arbitrary non-block-aligned crop or a larger padded canvas with the carrier shifted inside it.

The search remains public and key-independent. Pilot symbols with even indices propose a bounded translation on a 4-pixel grid followed by integer-pixel refinement; the disjoint odd-index pilot symbols validate the surviving placements. Payload, header, CRC, ECC, key and HMAC are absent. Search bounds are derived only from the known full transformed extent versus the observed extent and are capped so the experiment cannot silently become an unbounded image scan. Whole-tile-equivalent translations are treated as expected equivalences rather than false ambiguity because every v4 tile repeats the same pilot/data layout.

The synthetic matrix covers arbitrary crop, deeper crop, padded-canvas placement, affine geometry and mild perspective. The two local development originals also pass representative crop and padded-perspective cases. A deliberately wrong geometry loses strongly even after placement optimization, confirming that the placement search does not mask geometric error.

Build27 is **not** the final blind physical decoder. Geometry is still supplied independently for this placement qualification. The next gate is a crop/placement-tolerant coarse lattice/extent estimator that can feed geometry and placement jointly before the pilot is frozen. `prototype-2-search-p64` therefore remains non-normative.

## What v0.3.0-build26 adds

Build26 removes the known-homography assumption from the previous pilot experiment for an **auto-framed, bounded rotation/affine/mild-perspective search**. Geometry proposal is deliberately separated from pilot evidence: 64 deterministic data-plane coordinates are sampled at homologous positions of adjacent repeated 37x32 tiles, and only their self-consistency is scored. The data symbols themselves, pilot signs, payload, key, header, ECC and HMAC are never consulted. A small repeat-ranked bank is then scored by the public pilot on central tile rows and finally validated on held-out corner repetitions.

The synthetic Build26 matrix covers positive/negative rotation, anisotropic rotation+scale, X/Y shear, mild perspective, rotation+perspective+anisotropic scale, and rotation+shear. Correct geometry remains bounded to at most 6,400 evaluated repeat hypotheses in the current tests. The local corpus adds representative rotation+anisotropic-scale and combined-perspective cases on both development originals; all recover origin `(0,0)` and the intended geometry to the declared corner-error bound.

Build26 is **not** a general physical decoder yet. It assumes the canonical carrier extent is known and the transformed image is auto-framed; arbitrary crop/translation, print-camera acquisition, format framing/version negotiation, ECC choice and a real v4 payload encoder remain open. `prototype-2-search-p64` remains non-normative.

## What v0.3.0-build25 adds

Build25 asks a narrower question before attempting blind recovery: **if the correct projective mapping is known, does the public v4 pilot still identify the absolute cyclic origin after realistic geometric resampling?** A new experimental projective pilot sampler evaluates DCT evidence through a supplied homography while remaining completely disconnected from payload, key and HMAC evidence.

The deterministic synthetic matrix covers 16 cases: arbitrary rotation, anisotropic scale, X/Y shear, mild perspective, rotation+scale, rotation+shear, rotated 75% resize, perspective combined with JPEG/blur/noise, and perspective+crop with/without JPEG. All cases recover origin `(0,0)`. The same matrix is applied to the two local development originals, for 32 additional image/case observations; all recover the correct origin. The weakest local positive margin is **0.322069**, versus **0.063984** on the corresponding unmarked control.

Build25 uses a predeclared **0.10 positive-minus-negative margin separation floor** only as a development regression gate. It is not a decoder acceptance threshold and does not authenticate anything. Prototype-2 therefore remains non-normative: Build25 demonstrates pilot survivability under a correct geometric mapping, but blind pilot-assisted estimation of that mapping remains open.

## What v0.3.0-build24 adds

Build24 is the first dedicated Format-v4 pilot qualification checkpoint. It keeps the five frozen v3 core
files unchanged and searches the v4 pilot independently of payload, header, CRC, key and HMAC evidence.
The bounded search uses fixed seeds and explicit 200,000-candidate budgets for both the joint mask/sign stage
and the sign-refinement stage, so the result can be regenerated rather than accepted as an opaque constant.

The current non-normative winner is `prototype-2-search-p64`, candidate SHA-256
`858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`. Relative to the Build23 prototype,
maximum cyclic mask overlap falls from **8 to 7** and maximum absolute wrong-shift signed correlation from
**5 to 4**. Exhaustive contiguous partial-visibility tests retain worst correct-vs-wrong symbol margins of
**60/43/28/20/12** with 64/48/32/24/16 visible pilots; deterministic random-subset qualification reports no
false-origin ties/wins in the current 256-case-per-level set.

Build24 also adds a pilot-only synthetic carrier and detector. The carrier embeds the public pilot plus a
pseudo-random data plane so wrong-origin scoring sees realistic data interference, but it contains no message,
frame, ECC or authentication tag. The initial image-channel suite recovers the correct origin after JPEG,
small blur/noise/gamma changes, exact 75%/50% resize and aligned crop. The local `original pics` corpus can be
run through the same experiment with `make v4-pilot-corpus-test`; those image files are not shipped in source
archives.

The pilot is **not frozen yet**. Rotation, affine/shear, perspective, combined print-camera-like degradation,
version/framing design, ECC selection and real v4 print/scan acquisitions remain open before a functional v4
payload encoder is promoted.

## What v0.2.0 provides

The release baseline covers:

- adaptive Format v3 profiles: `robust`, `balanced`, `capacity`, `auto`;
- authenticated payloads up to 64 bytes;
- CRC32, truncated HMAC-SHA256, key-derived whitening and Hamming(7,4);
- JPEG/PNG input and PNG output;
- robust digital JPEG, resize and crop recovery on the qualification corpus;
- v3-only extraction with automatic profile detection;
- a reusable pure-Go core compile-checked for Linux, Windows, Android and iOS.

The decoder also contains **experimental**, bounded recovery for arbitrary
rotation, axis-aligned affine transforms, a fixed direct lattice-basis bank and
two mild projective/keystone hypotheses. These are measured research
capabilities, not universal guarantees. General homography estimation and the
physical print-camera channel are outside v0.2.0 scope.

## What v0.3.0-build23 adds

Build23 is the transition checkpoint between the completed v3 absolute-cycle research and the
experimental v4 branch. It does **not** change v3 embedding, extraction, search banks, Hamming,
whitening, profiles, HMAC semantics or physical PASS criteria. Instead it makes the project state
explicit and testable.

For v3, build23 adds [`docs/V3_FINAL_STATUS.md`](docs/V3_FINAL_STATUS.md), which records the
qualified baseline, the physical research outcome, the evidence that closes the absolute-cycle
branch, and the rules for future maintenance. The final build22 host qualification remained
`Release baseline: PASS` and `Qualification corpus: PASS`; only the historical experimental
`geometry-test` and `affine-test` boundaries remained red. The build22 physical topology probe
was non-discriminative: scanner 002 matched scanner 001 at 0.4375 mean modal phase fraction and
was more stable than the inclined smartphone case (0.3125), so the residual v3 topology signal
is preserved as evidence but not promoted.

For v4, build23 adds a completely separate provisional pilot foundation in
`watermark/experimental_v4.go`. The first candidate keeps the 37×32 / 64-pilot / 1120-data
layout from build22, reserves exactly one pilot in each of 64 spatial strata, uses 32 `+1` and
32 `-1` signs, and has no perfect non-zero cyclic alias. Exhaustive structural regression fixes
the current prototype metrics at maximum cyclic mask overlap 8/64 and maximum wrong-shift signed
correlation 5/64. These are **prototype metrics, not a Format-v4 specification and not a physical
recovery guarantee**. `make v4-foundation-test` protects the isolation and invariants. `make v3-freeze-check` verifies
`docs/V3_FROZEN_CORE_SHA256.txt`, making accidental v3 core drift visible. `make all-test` now
contains 32 targets.

## What v0.3.0-build22 adds

Build22 keeps Format v3, its encoder, production extraction, exact unwrap, cross-fit,
multi-partition stability and build20/21 diagnostics frozen. It performs the missing physical
follow-up to build21: instead of asking only whether repetition topology is mathematically
non-invariant, it measures whether that weak asymmetry is actually observable in the acquired
margin grids.

For every robust/balanced profile and every unit shift, the probe splits the repetition graph
into a large **shared registration set** and a small **held-out exclusive set**. Shared edges
alone register two phase hypotheses separated by the tested unit shift; the exclusive edges,
which were not used for registration, choose between them. The report exports winner-phase
stability, cell-direction consistency and absolute held-out effect size. The probe is
key-independent, diagnostic-only and cannot change sampling, candidate ranking or HMAC attempts.

On the current physical corpus the exclusive signal is measurable but not a safe absolute anchor.
The inclined smartphone image has mean winner modal-phase fraction **0.3125** and mean directional
cell consistency **0.6944**. Scanner 001 reaches **0.4375 / 0.7292** and the scanner-002 negative
control reaches **0.4375 / 0.6597**. Because the negative control is at least as phase-stable as
the useful smartphone case, the weak v3 topology residual cannot be promoted into a decode gate.
The frontal photograph remains not-applicable and the probe is skipped.

Build22 also starts a serious, still non-normative Format-v4 design branch. The current preferred
prototype is a **37×32 block tile with 64 public absolute-pilot positions**: tile area grows only
5.7% (minimum geometry 280×256 → 296×256 px), while 1120 data positions remain, preserving the
current 16/32/64-byte payload ceilings and nominal profile redundancy under an apples-to-apples
v3 framing/ECC assumption. See [`docs/FORMAT_V4_DESIGN.md`](docs/FORMAT_V4_DESIGN.md). No v4 encoder
or decoder is enabled in build22. `make all-test` now contains 30 targets.

## What v0.3.0-build21 adds

Build21 keeps Format v3, the encoder, production extractor and every build20 decision rule
frozen. It adds a **static key-independent observability audit** of the public Format-v3
tile mapping. The audit asks a narrower question than the physical decoder: for the exact
`{-1,0,+1}` block-cycle slips explored by the discrete unwrap, which public structural
constraints actually change when the tile origin moves by one block?

Two observables are audited without using a key, known header, payload, CRC or HMAC:

- the profile-specific repetition-pair topology induced by `v3CodeIndex`; and
- Hamming(7,4) word validity after applying the same wrong-origin tile mapping and hard
  aggregation used by the decoder, measured with deterministic synthetic valid codewords.

The result is mixed. `robust` and `balanced` have a mathematically non-invariant repetition
pair graph, so Format v3 does contain some absolute-origin asymmetry in those profiles.
However the nearest one-block alias preserves about **94.2--94.4%** of repetition pairs in X
and **92.2%** in Y, leaving only a roughly **5.6--7.8%** structural gap before camera noise.
`capacity` has no repetition pairs at all. Hamming structure strongly distinguishes most
horizontal/diagonal unit shifts, but vertical `±1` is a synthetic zero-syndrome alias for
`robust` and `capacity`; for `capacity` it is an exact whole-Hamming-word permutation.
`balanced` breaks that alias only weakly (about 0.088--0.093 synthetic nonzero-syndrome
fraction versus ~0.77--0.79 for typical horizontal/diagonal shifts).

Therefore the current v3 format is **not uniformly cycle-observable across profiles and
axes** under the audited key-independent mechanisms. Build21 does not add a new decoder gate:
the audit is static telemetry and regression evidence only. `make observability-audit-test`
protects these structural facts, and `make all-test` now contains 28 targets.

## What v0.3.0-build20 adds

Build20 keeps Format v3, deterministic encoding, production extraction, exact all-pairs
unwrap, held-out cross-fit and build19 multi-partition diagnostics frozen. It tests a new
**independent cycle-anchor hypothesis** using unguided cross-cell image-domain registration.
The observer never sees repetition groups, key material, known header bits, payload, CRC or
HMAC. It compares the exact top-1 and exact runner-up fields with a gauge-invariant weighted
relative-offset objective. The result is telemetry only: it cannot alter sampling or create
an authentication candidate.

On the current four-image physical corpus the anchor is unavailable for frontal
`foto stampa` because the exact unwrap is not applicable. On all three ambiguous cases it is
available and prefers the exact top-1 continuously, but the rounded integer-cycle agreement
is **0/9 for both top-1 and runner-up**. Mean anchor confidence is low: about 0.0575 on the
inclined smartphone image, 0.0168 on scanner 001 and 0.0798 on scanner 002. Scanner 001 is
a particularly important contradiction because held-out repetition evidence had preferred
the runner-up, while the independent pairwise anchor prefers top-1. Scanner 002, the negative
control, also prefers top-1. Build20 therefore does **not** promote this observer as an
absolute cycle anchor; it records that pairwise registration measures useful continuous
relative shape but does not yet identify the integer field. `make all-test` now contains 27
targets including `cycle-anchor-test`.

## What v0.3.0-build19 adds

Build19 keeps build18 scheduling and every decoder decision rule frozen, then asks whether
the integer-cycle field becomes reproducible when the same key-independent repetition
evidence is repartitioned many different ways. It is a diagnostic-only
**multi-partition stability experiment**: 8 fixed coded-bit-group partitions are tested
in both proposal/validation directions, for at most 16 held-out exact-unwrap trials on
the single best ambiguous bit candidate.

Each partition keeps all repeated positions for one coded bit together, so proposal and
validation remain disjoint exactly as in build17. Partition 0 deliberately reproduces
the original A/B split for continuity; seven additional deterministic, non-complementary
partitions probe sensitivity to the evidence split. For every available trial build19
records held-out support and the recentered local integer-cycle field. It then reports
full-field uniqueness/modal frequency, per-cell modal cycles and fractions, supported-only
stability, and mean pairwise field agreement. No modal field, vote or support fraction is
ever applied to the sampler or HMAC path.

On all three currently ambiguous physical cases all 16 trials are available, but all 16
complete fields are distinct. `foto stampa storta.jpg` has 9/16 held-out-supported trials,
yet those nine fields are also all distinct; mean per-cell modal fraction is only 0.2153
(0.2716 among supported trials). Scanner 001 has 6/16 supported trials and scanner 002
has 12/16, but again every supported complete field is unique. Mean pairwise cycle
agreement is only about 0.07--0.08. The frontal photograph remains `2/24` and
`not-applicable`, so the stability experiment is skipped.

This falsifies the hypothesis that the build17 A/B disagreement was merely an unlucky
split. On the present corpus, repartitioning repetition evidence does not reveal a hidden
persistent integer-cycle field. Future work should therefore seek a genuinely independent
cycle anchor rather than tune voting thresholds or add more repetition partitions.

`make stability-unwrap-test` protects partition disjointness/diversity, bounded synthetic
stability and public per-cell modal telemetry. `make all-test` now contains 26 targets.

## What v0.3.0-build18 adds

Build18 keeps every build17 decision rule frozen and addresses the two questions left by
that checkpoint: where independent cycle fields disagree, and how to avoid paying the
held-out cross-fit cost on bit candidates that cannot become the reported best result.

The cross-fit is now **best-candidate-only**. Every bounded full-decode candidate still
receives the normal blind observers, fractional lattice evidence and build16 all-pairs
exact unwrap. The held-out A->B/B->A experiment is deferred until bit-channel ranking is
complete, then run at most once on the single best bit candidate and only if its all-pairs
unwrap is genuinely ambiguous with an exact runner-up. Diagnostic JSON reports cross-fit
wall time, attempts and skipped candidates. No HMAC/full-decode/smooth-resample ceiling
changes.

Build18 also exports cell-level instability telemetry: the all-pairs proposed X/Y shift,
A->B and B->A local cycle coordinates, each fold-specific confidence, and per-cell cycle
agreement. Aggregate summaries compare the mean joint cross-fit confidence and lattice
confidence of agreeing versus disagreeing cells.

On the four available physical acquisitions the scientific verdict is identical to
build17. `foto stampa.jpg` remains `2/24` and not-applicable, so cross-fit is skipped.
`foto stampa storta.jpg` remains ambiguous (margin `0.02113`), with both held-out deltas
positive but only `2/9` cycle agreement. Those two agreeing cells have mean joint
cross-fit confidence `0.1337`, well below `0.3576` for the seven disagreeing cells; the
small agreement set is therefore not evidence for a high-confidence partial consensus.
Scanner `0270_001.jpg` remains `4/24`, both held-out directions reject top-1 and agreement
is `0/9`. Scanner `0270_002.jpg` remains `5/24`, the two directions split and agreement is
`0/9`.

In the development environment the selected cross-fit itself costs only tens of
milliseconds (`26--73 ms` on the three ambiguous best candidates) and is run once while
three lower-ranked bit candidates are skipped. End-to-end physical diagnostics returned
to roughly the pre-cross-fit range in this environment; absolute runtime must still be
rechecked on the qualification host.

## What v0.3.0-build17 adds

Build17 turns build16's post-hoc split-repetition check into a genuinely held-out
integer-cycle experiment. Repetition evidence is partitioned deterministically by
**coded-bit group**, not by individual pair: every repeated-position constraint for a
given coded bit belongs to exactly one fold, and no logical repetition position is
shared between proposal and validation folds. Two symmetric directions are evaluated:
fold A estimates the repetition controls and exact top-2 while fold B scores that
choice, then fold B proposes and fold A validates.

The proposal path remains key-independent. Each fold-specific repetition observer is
combined with the same independent fractional lattice measurement, then ranked by the
build16 exact `{-1,0,+1}` global solver. The held-out fold compares only the exact
geometric top-1/top-2. Build17 reports both directional validation deltas, exact-search
budgets, fold pair counts and agreement between the two independently proposed
**recentered local integer-cycle fields**. Different equivalent global tile anchors are
therefore not mistaken for local-cycle disagreement.

The cross-fit is deliberately diagnostic-only and runs only after the normal all-pairs
exact solver has already reported a genuine ambiguous top-2. It never creates another
HMAC candidate, never changes Format v3 and never overrides transactional rollback.
`make crossfit-unwrap-test` adds dedicated regressions and `make all-test` now contains
25 targets.

On the four physical acquisitions available in this session the experiment is
informative but not yet promotable. `foto stampa storta.jpg` is supported by the
opposite held-out fold in both directions (`+0.00465`, `+0.20522`) but the two proposed
local cycle fields agree in only `2/9` cells. Scanner `0270_001.jpg` is rejected in both
directions (`-0.18283`, `-0.02287`) and its proposals agree in `0/9`; scanner
`0270_002.jpg` is split (`+0.05639`, `-0.24050`) with `0/9` agreement. The frontal
`foto stampa.jpg` remains `2/24` and `not-applicable`, so cross-fit is not run there.
No physical payload authenticates and no smooth HMAC slot is consumed.

## What v0.3.0-build16 adds

Build16 hardens the global integer-cycle barrier introduced in build15. The former
width-64 beam search could discard a partial state that later became the true optimum
or runner-up, so its reported first/second margin was not a global uniqueness
certificate. Build16 removes that approximation: for each axis it enumerates the full
bounded `{-1,0,+1}` state space over eligible cells, at most `3^9 = 19683`
assignments per axis (`39366` across X/Y), and ranks the exact top-1/top-2 using a
deterministic tie break. The objective, thresholds and transactional rollback remain
key-independent and unchanged in meaning.

The exact solver closes a real false-accept hole found during the build15 audit. A
fixed regression reproduces a field for which the old beam margin would have accepted
an unwrap while the exact runner-up makes it ambiguous. Build16 also reports separate
X/Y eligibility and status, distinguishes `not-applicable`, `no-change`,
`rejected-improvement`, `ambiguous` and `accepted`, and never exports `Inf`/`NaN` for
the no-eligible path. A JSON regression protects that invariant. Mixed-axis cases are
also atomic: ambiguity on either axis prevents a partial integer-cycle commit and the
global result remains `ambiguous`.

As a second research diagnostic, build16 compares the exact geometric top-1/top-2
assignments using two deterministic disjoint subsets of the existing key-independent
Format-v3 repetition pairs. This `split-repetition-top2` signal is intentionally
**diagnostic-only**: the primary repetition controls were estimated from the complete
pair set, so the split is consistency evidence, not a sufficiently independent
held-out certificate and cannot override an exact ambiguous result.

The physical-corpus contract is also clarified. The acquisition files remain private,
while the intentionally public/reproducible test key defaults to `Piccotti` via
`PRINT_CAMERA_KEY`. Both smartphone and scanner directories are ignored by Git, and
`make private-corpus-manifest` can generate a local SHA-256 manifest. The canonical
smartphone filenames now use `.jpg`, matching their actual JPEG encoding; see
[`docs/PRIVATE_CORPUS.md`](docs/PRIVATE_CORPUS.md).

Format v3, deterministic encoding, the production extractor, profiles, whitening,
ECC, tile mapping and HMAC remain frozen. No HMAC or smooth-resample budget is
increased, and no diagnostic error-count improvement is treated as authentication.

## What v0.3.0-build14 adds

Build14 changes the second blind signal rather than expanding the HMAC or geometry
candidate bank. It reuses the already-selected 3x3 `LocalLatticeEstimate` records and
projects one measured lattice intersection per region back into each candidate's
canonical plane. The residual position modulo the frozen 8-pixel Format-v3 block grid
is reported as a **local fractional lattice phase**. It is key-independent: no magic,
header bits, payload data, whitening sequence or HMAC result enters the measurement.

The lattice observer cannot know the absolute integer block cycle by itself. Build14
therefore keeps repetition coherence as the cycle-index prior and uses lattice phase
only as modulo-one-block evidence. Fractional agreement can refine a control. A
bounded +/-1 integer unwrap is considered only for a low-confidence repetition
control when an affine robust prediction of the other cells shows at least 0.45 block
residual improvement and leaves at most 0.55 block residual. At most two unwrap passes
are allowed. No additional source-image DCT sampling is introduced.

The real result is promising but not promotable. Lattice-phase confidence is much
stronger than the build13 cross-cell score (roughly 0.38--0.52 mean confidence on the
reported best-hard spatial cases). With the final observed-drift sign convention and
bounded unwrap, several **same-candidate** diagnostic corrections improve sharply:
`foto stampa storta.jpg` 5 -> 2, `foto bici dritta.jpg` 7 -> 3 and scanner
`0270_001.jpg` 9 -> 4. The global hard baselines remain 7, 6 and 4 respectively; these
numbers must not be compared as if they were the same candidate.

None of those fields consumes an HMAC slot. Their mean fitted-control confidence stays
below the strict 0.20 decode gate and/or leave-one-out remains above 2.0 blocks. Other
candidates also worsen (for example scanner `0270_002` 7 -> 10), so the bounded unwrap
is retained as research evidence rather than promoted. All six physical acquisitions
remain unauthenticated.

Build14 adds `make lattice-phase-test` and expands `make test-list` with an `INPUT`
column and explicit notes that `pics/` qualification/research suites are corpus-sensitive
while `release-check` remains the production gate. `all-test` now contains 23 targets.

## What v0.3.0-build13 adds

Build13 keeps the build12 repetition-coherence observer as the primary key-independent
source of spatial controls and adds a second **guided cross-cell image-domain observer**.
For each pair of populated 3x3 cells, cross-correlation is searched only in a +/-1
integer-block neighbourhood around the relative offset predicted by repetition
coherence. This costs at most 36*9 = 324 pair probes when the primary observer exists;
the older +/-4 / 2916-probe pairwise search is retained only as fallback if repetition
evidence is unavailable. No additional source-image DCT sampling is introduced.

The secondary observer is deliberately unable to override the primary merely because a
local cross-correlation maximum exists. Its mean pair peak must exceed an independent
strength gate before the two observers can be fused or a bounded +/-1 cycle-slip
correction can be applied. Weak secondary evidence is still reported, including
observer distance and per-cell confidence, but it cannot lower the primary confidence,
change the smooth field or consume an HMAC slot.

On the six current physical acquisitions the guided observer follows the primary much
more closely than the old unrestricted pairwise search (mean observer distances about
0.51--0.85 blocks on the best hard candidates), but its mean pair peak is only about
0.053--0.065. A controlled synthetic field produces about 0.61, so no real case clears
the 0.10 consensus-strength gate. Consequently build13 applies zero cycle-slip fixes
on the real corpus and preserves the build12 blind result: `foto stampa storta.jpg`
still measures a same-candidate 5 -> 4 post-ECC improvement, remains above the strict
LOO decode gate and consumes no HMAC slot. No real acquisition authenticates.

Build13 also adds `make test-list`, a categorized index of release, qualification,
research, private physical-channel, portability and aggregate test targets. It is an
informational target and is not itself part of the 22-target `all-test` matrix.

## What v0.3.0-build12 adds

Build12 changes the **source of the spatial phase controls**, not the Format-v3
layout or the smooth-field model. The projective diagnostic now has a
key-independent blind registration observer. Its primary score uses the structural
redundancy already present in the robust and balanced Format-v3 profiles: positions
that map to the same coded bit should have coherent observed DCT margins at the
correct tile phase even when the bit value, key and payload are unknown. Capacity
has no repetition and is intentionally not treated as structural evidence.

The observer first searches the complete 35x32 tile phase on the already-built
aggregate grid for robust and balanced repetition coherence (at most 2240 score
probes). Each populated 3x3 spatial cell then searches only +/-2 blocks around that
blind global phase (at most 225 local probes total), with the same sub-block peak
interpolation and boundary-aware confidence logic introduced in build11. The common
translation is removed before fitting so the smooth field models only residual
spatial drift. If structural repetition is unavailable, a separately reported
pairwise cross-cell correlation fallback is bounded to at most 2916 score probes.
None of these blind probes resamples the source image or uses the key, expected
header bits, CRC or HMAC.

The build11 key-assisted local phase remains in the report only as an ex-post oracle
comparison. The smooth fitter used by the real projective path now accepts controls
only from `blind-self-registration`; oracle controls cannot influence the blind fit.
The affine/regularized-quadratic gates, two-resample ceiling and four complete/HMAC
decode slots remain unchanged.

On the six current real acquisitions the direct cross-cell correlation control is
weak, while intra-tile repetition is consistently measurable. Only
`foto stampa storta.jpg` produces blind fields that pass the diagnostic resample
gate in the final code. Its best measured blind correction changes the **same
candidate** from 5 to 4 known post-ECC header errors, with fit RMS about 1.27 blocks
and leave-one-out RMS about 2.04 blocks. That LOO value remains just above the strict
2.0 decode limit, so no smooth HMAC slot is consumed. The best frontal hard candidate
remains untouched at 2/24. No real acquisition authenticates.

`make blind-phase-test` covers key-independent relative registration, structural
repetition phase recovery, the deliberate absence of capacity repetition evidence,
blind-control isolation from the key-assisted oracle, explicit blind budgets and a
controlled blind-control affine correction that restores a valid Format-v3 HMAC.
The latter verifies the downstream blind-control fitting/correction path; it does not
claim that the real blind observer itself has already recovered a complete physical
print-acquisition warp.

## What v0.3.0-build11 adds

Build11 improves the quality of the build10 spatial phase controls before changing
model order. Format v3, the deterministic encoder, production `ExtractWithInfo`,
photometric/reliability logic and the maximum four complete/HMAC decode slots remain
unchanged.

Each bounded 3x3 spatial cell still searches only the existing +/-2-block phase
neighbourhood, but build11 now evaluates a clipped signed-sync correlation surface
around the integer maximum and uses parabolic interpolation to estimate a fractional
block offset. Every control gets an explicit confidence derived from local correlation
strength, peak prominence and curvature. A maximum that lands on the +/-2 search
boundary is flagged and automatically down-weighted because the true peak may lie
outside the bounded window.

The smooth fitter is now confidence-weighted and Huber-robust. It always fits the
six-parameter affine field first. With at least eight controls it also fits a
ridge-regularized quadratic field (six basis terms per axis), but the quadratic can
win only when leave-one-out RMS improves by both an absolute and relative margin.
The decode gate remains as strict as build10 (fit RMS <= 1.25 blocks and leave-one-out
RMS <= 2.0 blocks), now with an additional mean-control-confidence requirement and
the existing bounded-correction check. At most two corrected-grid resamples are still
allowed and no fifth HMAC slot can be created.

The six real acquisitions provide a useful negative result. The build10 promoted
`foto bici dritta` field no longer improves once boundary-truncated controls are
down-weighted: its eligible field changes 6->8 known post-ECC header errors and is
therefore rejected before HMAC. No real acquisition consumes a smooth-field HMAC slot
in build11. The strongest remaining diagnostic improvement is `foto bici storta`,
4->3 on one same-candidate field, but leave-one-out RMS remains about 2.18 blocks so
it is not decode-eligible. One scanner candidate selects the regularized quadratic
model by cross-validation, but its predicted correction exceeds the bounded correction
cap and is rejected. No real acquisition authenticates.

`make phase-surface-test` covers fractional peak interpolation, search-boundary
confidence reduction, confidence weighting, robust outlier rejection and the
quadratic cross-validation rule. `make smooth-phase-test` continues to verify exact
synthetic affine recovery through ordinary Format-v3 HMAC and the no-extra-slot budget.

## What v0.3.0-build10 adds

Build10 tests the build9 hypothesis that residual print-acquisition errors contain a
smooth spatial phase component. Format v3, the deterministic encoder, production
`ExtractWithInfo`, the photometric bank and the maximum four complete/HMAC decode
slots remain unchanged.

For each already-retained full-decode candidate, the bounded 3x3 spatial diagnostic
can fit a six-parameter affine field in canonical block coordinates:

```text
dx(nx,ny) = a0 + a1*nx + a2*ny
dy(nx,ny) = b0 + b1*nx + b2*ny
```

The observations are the existing +/-2-block local-phase measurements; no hidden
payload bit is used. A field needs at least six populated cells. Diagnostic resampling
requires fit RMS <= 2 blocks and leave-one-out RMS <= 3 blocks. Promotion into an
existing HMAC slot is stricter (RMS <= 1.25 and leave-one-out RMS <= 2), and the
corrected known prefix must improve without increasing known coded-bit errors. At
most two smooth-field resamples are allowed, in the existing photometric candidate
order, and a corrected grid replaces rather than adds a full decode slot.

A deterministic synthetic regression injects a known affine phase drift and verifies
that the fitted inverse field restores an authenticated v3 payload. A deliberately
non-smooth checkerboard field is measurable but rejected by the strict decode gate.

On the six current real acquisitions, the best measured same-candidate smooth changes
are mixed: `foto bici dritta` 7->3 post-ECC known-prefix errors, `foto stampa` 7->5,
`foto stampa storta` 8->6 and scanner `0270_002` 8->6, while scanner `0270_001`
changes 4->9 and `foto bici storta` 3->8. Only one real candidate passes the strict
decode gate: the mild-highpass fundamental candidate on `foto bici dritta`, which
changes 6->5 and consumes one of the existing four HMAC slots. It still does not
authenticate. Thus a smooth component is measurable, but a single global affine
field is not an adequate general correction model.

Build10 also makes spatial-cell collection ignore incomplete edge tiles for the 3x3
diagnostic only. Those edge blocks still contribute to the ordinary aggregate decode
grid, so this change does not alter Format v3 or the production decoder.

`make smooth-phase-test` covers affine-field fitting, exact synthetic recovery,
non-smooth rejection and budget invariants.

## What v0.3.0-build9 adds

Build9 keeps every build8 decoder budget and the production Format-v3 path
unchanged. It instruments spatial repetition already present in the virtual grid
so the same protected carrier positions can be compared across coarse 3x3 image
regions without re-reading a single DCT block. The spatial buckets are collected
during the existing full-grid sampling pass and are research evidence only.

For the best bit-channel candidate, JSON now reports a key-independent sign
agreement over all 1120 tile positions, plus known-prefix stability using the six
fixed v3 header Hamming words. A known coded bit is classified as stable-correct,
stable-wrong, or mixed across spatial cells. A simple per-cell majority is measured
only as an oracle diagnostic; it does not receive an HMAC slot and is not promoted
to extraction.

Build9 also measures a tightly bounded local phase check. Each populated spatial
cell may test only a +/-2-block neighbourhood around the already-selected global
phase (at most 25 phase positions per cell, 225 positions across 9 cells). This
uses the known v3 prefix and is therefore explicitly key-assisted/multiple-tested
research evidence. It is never used to authenticate or to change a decode grid.

Across the six current real acquisitions, mean key-independent tile-position sign
agreement is only about 0.635--0.683 even though scanner lattice consistency is
about 0.87. Most known header bits are mixed across regions and stable-wrong bits
are nearly absent. This strongly favors spatially varying misregistration/channel
damage over a small fixed set of systematically inverted bits. Small local phase
changes help some cases but hurt others, so independent per-cell re-phasing is not
promoted as a decoder.

The compact soft-Hamming report is also disambiguated: the soft result on the
best-hard candidate and the hard result on the best-soft candidate are now exposed
separately, avoiding comparisons between different geometry/view candidates.

`make spatial-channel-test` covers the spatial classification, key-independent tile
agreement, bounded local-phase neighbourhood and no-extra-read accounting.

## What v0.3.0-build8 adds

Build8 keeps Format v3, the deterministic encoder and production `ExtractWithInfo`
unchanged. It addresses two research findings from build7: scanner inputs below the
production 50 MP search limit could spend more than two minutes in the mature
baseline extractor before projective diagnostics started, and the original frontal
smartphone capture suggested that low-margin errors might benefit from a soft
Hamming experiment.

`diagnose` now has its own 16 Mi-pixel baseline-authentication budget. Larger
acquisition images skip the production baseline only inside diagnostics and proceed
directly to the already-bounded projective path. Complete virtual grids are also
area-bounded: candidates requiring more than 50,000 canonical blocks are aggregated
from at most 25 spatially distributed complete v3 tiles. The JSON report exposes
which baseline/full-grid budgets fired plus coarse/refinement/photometric/full-decode
stage timings. The production extractor itself is untouched.

A diagnostic reliability decoder performs one deterministic maximum-likelihood
choice among the 16 valid Hamming(7,4) codewords for each seven-bit word using the
signed DCT margins. It never creates a list of payload candidates. A full-decode
slot uses this path only when the fixed known v3 prefix shows at most one multi-error
word and wrong known bits are materially weaker than correct ones; the total remains
maximum four complete/HMAC decode slots. Synthetic regression proves that two weak
errors in one Hamming word can be recovered and authenticated by this path while the
ordinary hard decoder fails.

On the real corpus, the simple soft-Hamming model does not authenticate any capture.
The original frontal photograph activates one reliability decode, but its best known
header remains two bits wrong after both hard and soft decoding. Scanner acquisition
adds an important control: the two supplied 600-dpi 4960x7015 JPEG scans have very
strong lattice consistency (about 0.87) yet still show substantial protected-bit
errors. After the diagnostic baseline budget fix, both complete in about 10--11 s
instead of exceeding 120 s. This shows that strong lattice geometry alone does not
guarantee a clean v3 protected-bit channel.

`make reliability-test` covers soft-Hamming and bounded full-grid behavior. An
optional `make print-scan-test` uses the same HMAC-only corpus harness with
`PRINT_SCAN_DIR` (default `print-scan private`) so scanner and smartphone corpora can
remain separate.

## What v0.3.0-build7 adds

Build7 keeps the build6 lattice-first geometry and build5 photometric bank
unchanged, but instruments the exact protected-bit channel seen by the same
four candidates that already reach a complete virtual Format-v3 decode. No new
geometry candidates and no additional HMAC attempts are introduced.

For each full-decode candidate, the diagnostic selects the strongest bounded v3
profile/phase and reconstructs the protected coded bits together with their
signed DCT margins. The first three Format-v3 frame bytes (magic plus
version/profile) are known before payload recovery, so their 24 raw bits and six
Hamming(7,4) codewords can be compared exactly without knowing the hidden
message. JSON reports known protected-bit errors, how many of those six known
codewords contain 0, 1 or more than 1 error, post-ECC known-header bit errors,
global Hamming syndrome occupancy, coded-margin quantiles, and the mean margin of
correct versus wrong known bits. These are research diagnostics only; HMAC is
still the sole authentication criterion.

On the four current private smartphone captures, the best bounded candidate has
1/6, 3/6, 3/6 and 2/6 known header codewords respectively with more than one bit
error. The original frontal capture is closest to Hamming capacity: its best
candidate has 7/42 known coded-bit errors, only one known multi-error word and
2/24 known-header errors after ordinary Hamming correction. Its wrong known bits
also have substantially lower DCT margins than its correct bits, making it the
strongest future candidate for a bounded soft-decoding experiment. No real
capture authenticates in build7.

`make print-camera-test` is now corpus-driven: every PNG/JPEG directly inside the
private print-camera directory is tested in deterministic filename order. The
test no longer hardcodes the two original filenames, reports a corpus summary,
cleanly SKIPs an absent/empty corpus, and still marks each image PASS only after
a valid Format-v3 HMAC.

## What v0.3.0-build6 adds

Build6 keeps the build5 photometric bank and build4 geometry budgets but removes
the strong print-boundary detector as a mandatory gateway to projective
analysis. The boundary remains the preferred initializer when reliable; it is
never watermark evidence.

On large images with no strong boundary, a borderline lattice result can trigger
one bounded adaptive escalation to a finer diagnostic pyramid level (normally
1/4 after the default 1/8 + 1/16 pass). If that second look produces lattice
evidence, a geometrically sane weak boundary quadrilateral may be used only as a
projective seed. The report explicitly records `adaptive_escalated`, the added
divisor/reason and `lattice_first_fallback_used`.

This changes the held-out `foto bici dritta.jpg` case from a build5 early stop
(consistency≈0.533, no projective attempts) to an adaptive 1/4+1/8+1/16 result
with consistency≈0.615, lattice evidence and four bounded projective/HMAC decode
attempts. No HMAC authenticates. The companion held-out inclined capture and the
two original private captures keep their previous build5 geometry/sync results.

A large unmarked qualification image can still show diagnostic lattice evidence
after adaptive escalation; its weak-boundary confidence is too low to open the
projective fallback. This is intentional evidence that lattice/sync scores are
not watermark detection. Only the existing Format-v3 HMAC can authenticate.

## What v0.3.0-build5 adds

Build5 keeps the build4 geometry fixed and adds the first small deterministic
print-camera photometric bank. Format v3, the deterministic encoder and the
production `ExtractWithInfo` search order remain unchanged.

The diagnostic path now evaluates at most three already-selected projective
geometries through exactly three photometric views: `raw`, `local-normalize` and
`mild-highpass`. This adds at most nine known-header probes and still permits at
most four complete Format-v3/HMAC decodes. Photometric scoring never creates new
geometry candidates and never counts as watermark detection.

On the private frontal capture, `mild-highpass` improves the strongest retained
fundamental/phase-DLT known-header probe from z≈4.454 / fraction≈0.768 to
z≈4.695 / fraction≈0.783. On the inclined capture no photometric mode beats the
build4 raw global maximum z≈3.973. Neither photograph authenticates, so the
hidden real payload remains unknown.

Build5 also separates release invariants from qualification-corpus observations
in `make all-test`. Additional images in `original pics/` remain useful and are
reported, but a content-specific robustness miss no longer rewrites the meaning
of the release baseline. Very large/undimensionable inputs are consistently
SKIPped by bounded geometry suites instead of being reported as decoder errors.

## Why Go, portability and the planned GUI

Go is an intentional architectural choice for PixSeal, not only an implementation convenience. The project keeps the steganography/watermarking engine in a reusable Go package so the same core can be compiled for multiple operating-system/CPU targets instead of maintaining separate algorithm implementations. `make core-target-check` currently verifies compilation of that reusable core for **Linux/amd64, Windows/amd64, Android/arm64 and iOS/arm64** on every qualification run. `make build-all` also emits desktop CLI binaries for the supported desktop targets.

On mobile platforms the long-term goal is not to expose a command-line workflow to normal users. A separate graphical frontend is planned, with particular emphasis on **smartphone/tablet use**: choose or capture an image, embed/extract a message, and present the result and diagnostics without requiring shell commands. The GUI is intentionally kept outside the algorithmic core so desktop CLI tools, future Android/iOS applications and other frontends can share the same Go implementation and qualification corpus. Mobile packaging, platform permissions, camera/gallery integration and final UI technology are future product work; current Android/iOS checks validate core portability rather than claiming a finished mobile application.

## Build

Build44 is qualified with **Go 1.26.0** using deterministic project-controlled JPEG ingest. Build43 remains historically qualified with Go 1.25.1. Use the Makefile rather than a bare `go build`:

```sh
make
```

The normal Makefile path selects `GOTOOLCHAIN=go1.26.0`, verifies the qualified Build44 toolchain and **always rebuilds** `dist/pixseal`. To inspect the selected toolchain:

```sh
make toolchain-check
```

For a direct Go command outside the Makefile, select the qualified toolchain explicitly:

```sh
GOTOOLCHAIN=go1.26.0 go build -o pixseal ./cmd/pixseal
```

Go 1.26+ was not qualified for Build43 because its standard `image/jpeg` raster differs. Build44 removes that standard-library dependency through `internal/jpeglegacy`; the deterministic-raster and private phone physical gates now pass under Go 1.26.0, which is therefore the qualified Build44 toolchain. See [`docs/GO_TOOLCHAIN_COMPATIBILITY.md`](docs/GO_TOOLCHAIN_COMPATIBILITY.md).

The project has no external runtime dependencies. ImageMagick and GNU
`timeout` are required only by the shell robustness suites.

## Basic usage

```sh
pixseal embed \
  -in photo.png \
  -out sealed.png \
  -key "a long secret" \
  -message "hidden message" \
  -profile auto
```

A successful embed reports the selected profile:

```text
embedded 14 bytes using profile robust in sealed.png
```

Extraction does not require a profile:

```sh
pixseal extract -in sealed.png -key "a long secret"
```

The authenticated payload is written to stdout. Human-readable diagnostics such
as confidence, recovered profile and geometry correction are written to stderr.
For scripts or multiline payloads, use:

```sh
pixseal extract -in sealed.png -key "a long secret" -raw
```

`-raw` writes **only** the authenticated payload bytes to stdout, without an
added newline; decoder diagnostics remain available separately on stderr.

### Experimental Format v4 commands

Build31 adds an explicit experimental encoder without changing the v3 defaults:

```sh
pixseal v4-embed \
  -in photo.png \
  -out sealed-v4.png \
  -key "PixSeal-v4-TestKey-2026" \
  -message "Build31 physical test" \
  -profile robust
```

The original Build31 decoder remains available for an already aligned native 8-pixel lattice:

```sh
pixseal v4-extract -in sealed-v4.png -key "PixSeal-v4-TestKey-2026"
```

Build35 additionally exposes the qualified Build34 blind projective+crop path for controlled physical captures:

```sh
pixseal v4-extract-projective \
  -in scan-marked-a.png \
  -key "PixSeal-v4-TestKey-2026" \
  -width 1632 \
  -height 1632
```

`-width` and `-height` are the block-aligned dimensions of the digital carrier **before printing**, not the acquired image dimensions. This command still does not use HMAC as a geometry oracle: public structure/pilot evidence must pass the unchanged geometry gate first. The explicit command split prevents the experimental v4 branch from silently changing stable v3 `embed` / `extract` behavior.

### Strength

The public CLI accepts finite strength values from **4 through 120**. The
default is 24 when `-strength` is omitted. Explicit `0`, `NaN` and infinities
are rejected. The Go API retains `Options.Strength == 0` as the internal
"use default" sentinel for compatibility.

### Profiles

| Profile | Maximum payload | Protected frame | Hamming-coded bits | Tile redundancy |
|---|---:|---:|---:|---:|
| `robust` | 16 B | 32 B | 448 | 2.50x |
| `balanced` | 32 B | 48 B | 672 | 1.67x |
| `capacity` | 64 B | 80 B | 1120 | 1.00x |

`auto` chooses the most robust profile capable of containing the actual payload.
Payload limits are measured in bytes, not Unicode characters.

## Format v3 frame

Every profile has a fixed frame size, but the authentication tag follows the
**actual payload**, not the end of the profile capacity region:

```text
magic(2)
+ version/profile(1)
+ payload length(1)
+ CRC32(4)
+ payload(actual length)
+ HMAC-SHA256 tag(8, truncated)
+ zero padding to the profile frame size
```

The tag offset is therefore:

```text
8 + actual_payload_length
```

The HMAC covers the 8-byte header and actual payload. The completed frame is
whitened, Hamming(7,4)-protected and mapped onto the fixed 35x32 = 1120-position
DCT tile with a deterministic modular stride.

Whitening is **not encryption**. Encrypt sensitive content before embedding if
confidentiality is required.

See [`docs/ALGORITHM.md`](docs/ALGORITHM.md) for the precise specification.

## Geometry diagnostics (v0.3 research)

```sh
pixseal diagnose -in capture.jpg
pixseal diagnose -in capture.jpg -json
pixseal diagnose -in capture.jpg -key "a long secret" -json
```

The diagnostic command reports local/global lattice evidence, the optional print
boundary initializer, bounded projective candidates, budgets and timings. When a
key is supplied, small ordinary carriers still use the mature baseline extractor;
the research path may additionally rank bounded projective hypotheses with the
key-known Format-v3 header and run at most four complete virtual projective
decodes. Those probe scores are diagnostic only. A valid Format v3 HMAC remains
the only success criterion.

## Analyze and capacity

```sh
pixseal analyze -in photo.png -message "hidden message"
pixseal analyze -in photo.png -bytes 18
pixseal capacity -in photo.png
pixseal capacity -in photo.png -details
```

`analyze` clearly labels deterministic facts and heuristic recommendations. Its
image-detail estimator uses the same white-background alpha flattening as the
encoder. Recommendations are advisory and are not recovery guarantees.

## Image pipeline and safety

Supported decoded input formats are PNG and JPEG; detection is based on file
content rather than filename extension. Output is always PNG.

If the requested output has no extension, `.png` is appended. If it has another
extension, it is replaced with `.png`. Unix dotfiles are handled as basenames:
`.sealed` becomes `.sealed.png`, not `.png`.

Output safety rules:

- without `-force`, any existing directory entry is protected;
- `-force` can replace **regular files only**;
- directories, symlinks and other special files are rejected;
- no-clobber publication uses an atomic hard-link commit when supported, with an
  exclusive-create fallback where hard links are unavailable;
- newly created files respect the process umask;
- replacing an existing regular file preserves its permission bits;
- the temporary file is `Sync()`ed before commit.

On Windows the fallback replacement path uses backup-and-restore because
`os.Rename` cannot always replace an existing file. This reduces replacement
risk but is not claimed to provide filesystem-level crash durability; directory
`fsync` is not attempted.

### Image-size policy

Before decoding pixel data, the CLI uses `image.DecodeConfig` and rejects inputs
above **300,000,000 decoded pixels**. The core applies the same working-image
limit before allocating its compact extraction pixel plane. This keeps the
verified ~200 MP class usable while avoiding obviously unbounded allocations.

Some **generated geometry candidates** use the stricter 50,000,000-pixel search
bound. Such candidates are skipped rather than materialized.

The v0.3 diagnostic path now builds sampled luminance analysis planes rather than
an additional full RGB extraction plane. The Go image decoder still materializes
the decoded source image; build4 avoids additional full-size rectified/projective
working copies but is not yet a streaming/tiled source decoder.

## Decoder order

The current v3 extractor is bounded and ordered so speculative geometry does not
suppress established recovery paths:

1. direct grids at apparent block sizes 8, 6 and 4 pixels, including pixel phase;
2. gated lossless quarter turns (90/180/270 degrees);
3. pure isotropic fractional-scale recovery through virtual affine sampling;
4. one targeted physical bilinear normalization fallback when the scale evidence
   is decisive;
5. early axis-aligned affine recovery when zero-degree evidence supports it;
6. two mild projective/vertical-keystone hypotheses;
7. fixed direct lattice-basis composition bank;
8. arbitrary-angle rotation estimation and rectification;
9. final axis-aligned affine fallback when rotation evidence is weak.

HMAC authentication is always the final success criterion.

### Pure fractional resize budget

The non-integer direct-scale path has:

- 13 fixed scale hypotheses: 95, 90, 85, 80, 70, 65, 60, 55, 45, 40, 35, 30, 25%;
- up to 3 candidate phases per scale for single-tile probing;
- at most 6 shortlisted scales x 3 phases = **18 full-carrier virtual grids**;
- at most one decisive physical normalization candidate, with the base target
  size plus its eight ±1-pixel neighbours.

100%, 75% and 50% are already handled directly by apparent block sizes 8, 6 and
4 and are not duplicated in that list.

### Experimental geometry budgets

Arbitrary rotation:

- 2 zero-degree probes (8/6 px);
- at most 720 non-zero coarse probes;
- at most 33 fine probes;
- at most 2 refined angle candidates reach rectification.

Fixed direct lattice bank:

- 4 basis shapes: 110x90, 90x110, 105x95, 95x105%;
- 361 angles per shape (-45..+45 at 0.25 degree);
- **1444 sparse probes** maximum;
- at most **48 candidates per shape / 192 total** receive stronger coherence
  evaluation;
- at most 4 matrices x 3 phases per shape group = **12 full authenticated grids per group**,
  with two sequential groups (±10% then ±5%) for a **24-grid overall worst case**.

Mild perspective:

- exactly 2 fixed vertical-keystone hypotheses (`top-narrow-4`,
  `bottom-narrow-4`);
- one aligned phase each;
- at most **2 projective full-grid aggregations**.

These budgets describe the implemented search space; they are not statements of
universal recovery probability.

## Metadata and color/orientation limitations

PixSeal works on decoded pixels and outputs 8-bit NRGBA flattened against white.
Consequently:

- 16-bit PNG input becomes 8-bit output;
- EXIF, XMP, comments and ICC/profile metadata are not preserved;
- Go's JPEG decoder does not automatically apply EXIF Orientation, so a JPEG
  shown rotated by a metadata-aware viewer may be processed according to its
  stored pixel orientation;
- losing an ICC profile can change appearance in color-managed workflows, not
  merely remove descriptive metadata.

EXIF-orientation normalization and color-management preservation are deferred to
v0.3 because they require an explicit image-pipeline policy.

## Tests

```sh
make                 # build with qualified Go 1.25.1
make toolchain-check # verify qualified Go toolchain selection
make test            # complete Go tests + local image round trips
make release-unit    # release-gate Go regressions only
make research-unit   # experimental geometry Go regressions
make lattice-estimator-test # v0.3 bounded local-lattice diagnostics
make homography-test   # v0.3 boundary/homography/virtual decoder regressions
make photometric-test  # v0.3 bounded raw/normalize/high-pass diagnostics
make bit-channel-test  # v0.3 protected-bit/Hamming channel diagnostics
make reliability-test  # v0.3 bounded soft-Hamming/full-grid diagnostics
make spatial-channel-test # v0.3 spatial/tile protected-bit diagnostics
make print-camera-test # private real print-camera corpus; SKIP when absent
make print-scan-test   # optional separate private scanner corpus
make private-corpus-manifest # local SHA-256 manifest for both private physical corpora
make deep-test       # stable JPEG/resize/crop baseline
make extreme-test    # non-strict progressive resize/crop limit map
make geometry-test   # experimental rotation/combined geometry
make affine-test     # experimental axis-aligned affine matrix
make composition-test # historical composed regression
make lattice-test     # direct lattice-basis regression
make perspective-test # two-hypothesis mild projective regression
make core-target-check # reusable core: Linux/Windows/Android/iOS
make version-check    # VERSION vs buildinfo consistency
make release-check    # release baseline gate
make all-test         # every suite sequentially + summary
```

`extreme-test` is intentionally **non-strict**: individual FAIL rows identify
measured limits and do not make the target fail. `deep-test`, `geometry-test`,
`affine-test`, `composition-test`, `lattice-test` and `perspective-test` support
strict mode, and `make all-test` invokes those suites with strict mode by default.

`make release-check` is self-contained and forces strict semantics for the stable
`deep-test` baseline. It includes `release-unit`, private-corpus round trips and
core portability while excluding `research-unit`. `make all-test` keeps research
failures visible separately from the release baseline and reports partial target
sets as `PARTIAL`/`NOT RUN` rather than as a full qualification PASS.

The local shell suites use `original pics/`, which is private and excluded from
source archives. An empty qualification corpus is an error rather than a false
PASS. `print-camera-test` separately scans every PNG/JPEG directly inside the
private print-camera directory and cleanly reports `SKIP` when that corpus is
absent or empty; when present, the test key defaults to `Piccotti` and can still be
overridden with `PRINT_CAMERA_KEY=...`. Each image can report `PASS` only if the
authenticated Format v3 payload is recovered. The same key policy applies to
`print-scan-test`. `Piccotti` is an intentional corpus test key, not a private credential.

The full shell qualification harness is intended for **Linux/WSL with Bash >= 4,
GNU-compatible userland (including `timeout` and `sort -z`) and ImageMagick**.
This harness requirement is separate from reusable-core portability.

## Current limitations

- Maximum v3 payload: 64 bytes.
- Minimum aligned embed geometry: 280x256 pixels.
- CLI embedding currently accepts text via `-message`; the core stores bytes.
- General affine/projective print-camera recovery is not a v0.2.0 guarantee.
- v0.3.0-build17 retains the bounded lattice/projective/photometric/reliability
  research path and build16 exact top-2, then adds a coded-bit-group-disjoint
  A->B/B->A repetition cross-fit as diagnostic telemetry. No real smartphone or
  scanner acquisition has authenticated; lattice/sync/ECC/spatial/blind/unwrap/
  cross-fit/smooth-field evidence remains diagnostic only.
- Experimental rotation/affine performance is image-content dependent.
- Failed extraction can be substantially more expensive than successful
  extraction because bounded candidates must be exhausted.
- No neural model is used.
- No professional cryptographic or steganalytic audit has been performed.

## Release status

**v0.3.0-build17** is an experimental development snapshot. It preserves the
qualified v0.2.0 encoder/Format-v3/production-extractor baseline and changes only the
separate diagnostic research path. Build17 retains build16 exact bounded enumeration
and adds symmetric coded-bit-group-disjoint repetition proposal/validation cross-fit.
The cross-fit is diagnostic-only and cannot promote an ambiguous unwrap.

**v0.2.0** remains the qualified stable release of the Format v3 line. It was
promoted from RC4 after the stable release baseline passed on the private
qualification corpus with 72/72 baseline transformations recovered. The
experimental geometry measurements remain explicitly non-normative and are
recorded separately from the release gate.

Qualification results are recorded in [`docs/RESULTS.md`](docs/RESULTS.md).

## License

PixSeal is **source-available for noncommercial use** under the
**PolyForm Noncommercial License 1.0.0**. The SPDX identifier is
`PolyForm-Noncommercial-1.0.0`. Noncommercial use, modification and
redistribution are permitted only under the license terms; commercial use
requires a separate commercial license from the applicable copyright holder(s).

This licensing model is intentionally described as *source-available*, not as
OSI open source. Historical copies that were already distributed under MIT keep
the MIT rights granted with those copies; the current license applies
prospectively to distributions carrying it.

See [`LICENSE`](LICENSE) for the complete license text, [`NOTICE`](NOTICE) for
the required project notice, and [`docs/LICENSING.md`](docs/LICENSING.md) for
licensing scope, redistribution guidance and the separation between source code
and private research corpora.

