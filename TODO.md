# Build44 status

- [x] Preserve the qualified Build43 watermark/geometry/data path unchanged.
- [x] Add a project-controlled pure-Go pre-Go-1.26 JPEG decoder and deterministic raster fixture.
- [x] Route CLI JPEG ingest through the deterministic decoder instead of the toolchain standard library.
- [x] Qualify the unchanged private strength-48 phone matrix with Go 1.26.0; controls 3/3 reject and A/front, A/mild, A/angle, B/front authenticate.
- [x] Promote Go 1.26.0 as the qualified Build44 toolchain after the deterministic JPEG and physical gates pass.
- [ ] Resume geometry research on B/mild and B/angle from the now-qualified deterministic JPEG / Go 1.26 baseline, without relaxing Build41 gates or using secret/data evidence.

# PixSeal TODO

This file contains open work only. Completed milestones belong in `HISTORY.md`
and `CHANGELOG.md`.

## v0.3.0 — geometry research

- [x] Build a bounded diagnostic direct/local lattice estimator separated from
      `ExtractWithInfo`.
- [x] Estimate local transformed lattice vectors `u` and `v` and combine multiple
      regions without allowing one region to set the global geometry.
- [x] Add a coarse print-boundary initializer that is never treated as watermark
      evidence.
- [x] Produce a bounded projective scale shortlist and explicit homographies.
- [x] Add a virtual projective DCT sampler and synthetic authenticated end-to-end
      regression, capped at four complete v3 decode attempts.
- [x] Add bounded spatial phase-consensus ranking and ±1%/±2% scale refinement
      without photo-specific target dimensions.
- [x] Add explicit fundamental-period/scale selection using multilevel and
      multiregion support before phase refinement, without photo-specific target
      dimensions.
- [x] Fit a bounded deterministic DLT homography from four local sync-phase
      correspondences when available.
- [x] Add bounded modulo-8 block-origin/sub-pixel refinement with monotonic
      full-probe verification.
- [x] Add a small deterministic post-geometry photometric bank (`raw`,
      `local-normalize`, `mild-highpass`) with bounded probes and unchanged HMAC
      semantics; build5 shows a modest frontal sync improvement but no HMAC.
- [x] Add a bounded adaptive finer lattice level for large no-boundary borderline
      cases and allow a lattice-positive weak-boundary seed to reach the projective
      decoder without making the physical boundary a detection gate.
- [x] Quantify the protected-bit channel on the already-bounded full-decode
      candidates using the six exactly-known v3 header Hamming words, post-ECC
      known-header errors, syndrome occupancy and DCT reliability margins.
- [x] Test a bounded reliability-aware Hamming decoder on the existing Format v3
      code without increasing the four full-decode/HMAC-slot ceiling. The synthetic
      weak-double-error case authenticates, but no real acquisition does; the simple
      per-word soft model is therefore not promoted to production extraction.
- [x] Use spatial/tile agreement to determine whether remaining errors are local
      or systematic. Build9 finds key-independent tile-sign agreement only about
      0.64--0.68 and almost all known prefix bits mixed across spatial cells.
- [x] Fit/test a **smooth low-DOF spatial phase or residual-warp field** using the
      build9 cell evidence. Build10 uses a six-parameter affine field with separate
      diagnostic/decode gates; one real candidate is promoted but no HMAC succeeds.
- [x] Improve spatial control estimation without increasing independent per-cell
      freedom. Build11 adds fractional phase-surface interpolation, explicit
      confidence, search-boundary penalties and robust weighted fitting.
- [x] Evaluate a regularized low-order extension under the same bounded four-HMAC-slot
      rule. Build11 compares a ridge-regularized quadratic field by leave-one-out; it
      wins one real model comparison but fails the bounded correction gate and is not
      promoted.
- [x] Derive continuous local-phase controls from key-independent structural evidence.
      Build12 uses robust/balanced intra-tile repetition coherence, with pairwise
      cross-cell correlation as a bounded fallback; the real smooth-field path no
      longer uses key-assisted oracle controls.
- [x] Test a second key-independent image-domain observer for **cycle-slip resolution**.
      Build13 uses repetition-guided cross-cell correlation in a bounded +/-1 window;
      real mean pair peaks remain only about 0.053--0.065 versus about 0.61 synthetic,
      so the observer is kept as diagnostic-only and applies no real cycle-slip fixes.
- [x] Add a genuinely complementary **local lattice fractional-phase** observer.
      Build14 projects the existing key-independent 3x3 lattice estimates into the
      candidate canonical plane, measures modulo-8 phase and tests a bounded +/-1
      smooth unwrap. It produces several strong same-candidate diagnostics but also
      degradations; no real field passes all HMAC-promotion gates.
- [x] Add a first **global discrete integer phase unwrap** without using the known header.
      Build15 introduced the axis-separable bounded solver and transactional rollback.
- [x] Replace build15 beam ranking with an **exact bounded top-2 certificate**.
      Build16 enumerates at most 3^9 assignments per axis, closes a deterministic beam
      false-accept regression, reports X/Y status separately and keeps no-eligible JSON finite.
- [x] Add a first **split-repetition top-1/top-2 consistency diagnostic**.
      Build16 partitions key-independent repetition pairs deterministically into two folds.
      It remains diagnostic-only because the primary repetition controls used the full pair set.
- [x] Build a genuinely **held-out integer-cycle discriminator experiment** with
      proposal and validation evidence disjoint by coded-bit group. Build17 performs
      symmetric A->B/B->A exact-top2 cross-fit. It is informative but does not yet
      resolve a stable real cycle field, so it remains diagnostic-only.
- [x] Localize build17 cross-fold instability cell-by-cell and remove duplicate cross-fit
      work from lower-ranked bit candidates. Build18 shows the two agreeing cells on
      `foto stampa storta.jpg` are not the strongest controls, so naive partial consensus
      is not justified; held-out work now runs only on the final best bit candidate.
- [x] Test stability/reconciliation across held-out folds and an independent image-domain
      cycle anchor. Build19 shows repetition-only fields are non-persistent across 16 trials;
      build20 shows unguided pairwise registration prefers top-1 continuously but gives 0/9
      integer-cycle agreement on every ambiguous physical case, so neither route is promoted.
- [x] Audit the existing v3 public structure for an **absolute key-independent cycle reference**.
      Build21 finds weak but real repetition-topology asymmetry in robust/balanced, no repetition
      anchor in capacity, and Hamming vertical aliases (exact for capacity). The audited mechanisms
      therefore do not provide uniform absolute cycle observability.
- [x] Decide the v3 absolute-cycle branch: build22 physical topology observability does not separate
      scanner 002 from useful acquisitions, so no further threshold/voting round is justified; preserve
      v3 as the frozen compatibility baseline and move new origin work to experimental Format v4.
- [ ] Test whether a jointly regularized spatial model can improve several acquisitions
      from blind controls without per-cell oracle freedom; do not increase HMAC slots or
      model order unless cross-validation improves out-of-sample controls.
- [ ] Reduce negative-case cost of the mature arbitrary-geometry research paths;
      `geometry-test` remains the dominant `all-test` runtime.
- [ ] Decide whether/when the diagnostic projective path has enough evidence to
      be promoted into `ExtractWithInfo`; do not integrate before the complete
      synthetic/real/negative/regression promotion sequence.

## v0.3.0 — print-camera research

- [x] Clarify that the physical acquisition files are private but the test key `Piccotti` is intentionally public/reproducible; keep it as the Makefile/script default and ignore both acquisition directories.
- [x] Add a Git-ignored local SHA-256 corpus manifest target and canonical `.jpg` smartphone names.
- [ ] Optionally complete the historical six-image v3 corpus if the missing bicycle captures become available;
      treat this as archival comparison data, not a prerequisite for v4 development.

- [x] Provide a private full-resolution smartphone regression target that never
      ships the corpus, SKIPs when absent and PASSes only on valid v3 HMAC.
- [x] Run the bounded lattice/boundary/projective diagnostics on both original
      real smartphone captures; build2 obtains lattice evidence on both, build3 adds
      phase-aware refinement and build4 separates the fundamental scale from
      supported aliases/harmonics.
- [x] Test integer-cycle stability across 8 deterministic coded-bit-group partitions / 16 held-out directions; all three ambiguous physical cases produce a different complete field in every trial, so repetition-only repartition/voting is not a viable promotion path on the current corpus.
- [x] Test an independent image-domain cycle anchor (build20) and audit existing Format-v3 structural observability (build21); neither yields a uniform safe absolute-cycle gate.
- [x] Close the v3 absolute-cycle promotion branch after build22: the held-out physical topology probe
      is measurable but non-discriminative against scanner 002. Keep Format v3 frozen and require new
      absolute-origin work to use the isolated experimental v4 branch.
- [x] Close the current v3 physical research checkpoint **without** authenticated recovery: no supplied
      smartphone/scanner acquisition produced a valid v3 HMAC, and build22 negative-control evidence
      does not justify further cycle-threshold tuning. Preserve this as an explicit unresolved v3 limit.
- [ ] Revisit v3 physical recovery only if a genuinely independent observable, new acquisition regime or
      clear implementation bug appears; do not resume by retuning existing topology/cross-fit evidence.
- [x] Add known-header protected-bit/ECC diagnostics without using the hidden
      payload as an oracle.
- [x] Add a small deterministic photometric-normalization bank only after
      geometry/phase refinement shows it is needed.
- [ ] Move controlled synthetic camera-channel expansion to the v4 branch first; use v3 only as a frozen
      comparison baseline unless the experiment is explicitly format-agnostic.
- [x] Exercise a frozen build on a new held-out print-camera pair; build6 recovers
      the projective path on the previously boundary-blocked frontal case without
      photo-specific tuning.
- [x] Make the private print-camera harness automatically test every PNG/JPEG in
      the corpus rather than hardcoding the first two filenames.
- [x] Generalize on an additional held-out smartphone pair and add two scanner
      acquisitions of the same printed-image set as a distinct physical channel.
- [x] Bound diagnostic authentication on scanner-sized inputs; build8 reduces both
      supplied 34.8 MP scans from >120 s to about 10--11 s without changing the
      production extractor.
- [x] Do not change Format v3 on the current physical evidence; new print/acquisition work now belongs to
      a separately embedded/printed v4 corpus, with v3 retained as a comparison baseline.

## Future platform work

- [x] Document Go as an intentional portability choice: one reusable core is cross-compiled by `core-target-check` for Linux/amd64, Windows/amd64, Android/arm64 and iOS/arm64.
- [ ] Design and implement a graphical frontend over the reusable Go core, prioritizing smartphone/tablet workflows (camera/gallery selection, embed/extract and clear result/diagnostic presentation).
- [ ] Decide the Android/iOS packaging/binding strategy and platform-specific camera/gallery/file permissions without duplicating the watermark algorithm outside the Go core.
- [ ] Run the no-clobber CLI regressions on a real Windows host when convenient;
      v0.2.0 qualification already includes Windows/amd64 cross-compilation.
- [ ] Decide explicit EXIF Orientation normalization policy.
- [ ] Decide ICC/color-management preservation policy.
- [ ] Consider tiled/lazy pixel access to reduce peak memory on very large images.
- [x] Build44 deterministic JPEG ingest is physically qualified on Go 1.26.0; the Go 1.25.1 pin is no longer required for Build44 and remains only as Build43 historical qualification evidence.

## Format v3 closure / maintenance

- [x] Freeze the qualified v3 interoperability baseline after build22; keep golden encoder fingerprints,
      profiles, Hamming/whitening/HMAC rules and production extraction behavior unchanged. Protect the five
      core source files with `docs/V3_FROZEN_CORE_SHA256.txt` / `make v3-freeze-check`.
- [x] Record the final absolute-cycle research conclusion and negative-control evidence in
      `docs/V3_FINAL_STATUS.md`.
- [ ] Keep the historical v3 geometry/affine research failures visible; optimize runtime only if behavior
      and qualification counts remain unchanged.
- [ ] Re-run the frozen v3 qualification matrix when shared core utilities change for v4 work.
- [ ] Do not add new v3 HMAC slots, topology thresholds or cycle voting without genuinely independent evidence.

## Format v4 experimental branch

- [x] Quantify the 37×32 / 64-public-pilot / 1120-data layout without reducing v3-like payload ceilings.
- [x] Add build23 prototype foundations completely isolated from production v3 encode/decode.
- [x] Add a provisional spatially stratified, sign-balanced 64-pilot candidate and exhaustive cyclic-alias tests.
- [x] Expand pilot optimization beyond prototype-1 with a reproducible fixed-seed joint mask/sign search,
      sign refinement, exhaustive cyclic metrics and partial-visibility qualification; build24 selects the
      non-normative `prototype-2-search-p64` candidate.
- [x] Extend the synthetic-channel qualification to arbitrary rotation, anisotropic scaling, shear, mild perspective, crop and combined transformations with independently supplied geometry; Build25 adds 16 deterministic cases plus the same matrix on the local originals.
- [x] Define and document the first experimental v4 framing/version marker. Build31 uses authenticated header/profile bytes `0x41/0x42/0x43`, v4-specific whitening/HMAC domains and a deterministic compatibility vector; it remains non-normative.
- [ ] Decide whether normative v4 retains Hamming(7,4) or adopts a stronger ECC. Build31 intentionally keeps Hamming(7,4) only as an apples-to-apples framing/channel baseline.
- [x] Implement an **experimental-only** v4 encoder/decoder path behind explicit `v4-embed` / `v4-extract` commands; stable v3 `embed` / `extract` remain unchanged. Build31 extraction is native-lattice/aligned only.
- [x] Implement an experimental pilot-only detector for an already-resolved 8/6/4-pixel lattice; it reports
      absolute cyclic origin, score, runner-up and margin without payload/HMAC evidence.
- [x] Generalize the pilot-only detector to a first bounded blind rotation/affine/shear/mild-perspective recovery path: Build26 uses repeated data-plane self-consistency for coarse proposal, then public-pilot ranking and held-out pilot validation.
- [x] Complete the first synthetic v4 image-channel matrix for JPEG/resize/crop/rotation/affine/perspective/blur/gamma/noise/combined transformations. Build24 covers the aligned photometric/sampling half; Build25 covers the geometric half with known homographies and unmarked controls.
- [x] Generalize Build25 from known geometry to **bounded blind pilot-assisted geometry estimation** for auto-framed transforms, with deterministic repeat-based proposal and explicit runner-up/negative-control telemetry (Build26).
- [x] Qualify arbitrary crop/translation and padded-canvas placement with geometry supplied independently; Build27 uses disjoint pilot halves and no payload/header/HMAC evidence.
- [x] Combine blind affine geometry with Build27 placement under arbitrary negative crop: Build28 jointly recovers rotation, anisotropic scale and crop/translation with structural-only geometry selection and split-pilot placement validation.
- [x] Add a bounded Build29 joint **projective/perspective + crop** search with explicit ACCEPT/SAFE-REJECT gates; synthetic and `PJ_lingua` qualify, while the smaller `PJ_piccolo` case is intentionally rejected rather than falsely accepted.
- [x] Add a bounded Build29 joint **affine + positive padded-canvas** search and qualify the synthetic case.
- [x] Add a Build30 **semantic candidate lock** for `prototype-2-search-p64`: exact 37×32 geometry, ordered 64-position mask/sign sequence and SHA-256 identity are now protected against accidental mutation. This is a development lock, not normative Format-v4 promotion.
- [x] Add a Build30 known-mapping corpus audit covering projective-crop and affine-padded fixtures on both local originals; all four marked cases recover origin `(0,0)` with margin >= 0.484 while matching unmarked controls remain <= 0.0106.
- [ ] Improve photographic padded-canvas geometry/placement ranking: both current local originals remain SAFE REJECTED by the joint search even though Build30 confirms strong pilot evidence under the exact mapping.
- [ ] Improve projective ranking on small/low-redundancy carriers so `PJ_piccolo` can be accepted without weakening the Build29 margin gate.
- [ ] Re-qualify shear jointly with placement and broaden padded/projective joint coverage before claiming a general blind v4 geometry decoder.
- [ ] Replace the current known canonical-extent assumption with a coarse tile/lattice extent estimator suitable for physical captures before pilot ranking.
- [x] Add Build33 stage-by-stage MQ joint-projective ranking observability: compare random-data and authenticated-v4 carriers under the same transform, record truth-basin rank after structural/half-pilot/full-proposal stages, and dump distinct basins without using HMAC as an oracle.
- [x] Build34 replaces the MQ joint-projective ranking failure with bounded multi-anchor, structural-preserving, translation-phase-aware refinement and validates two different authenticated v4 payloads end-to-end.
- [x] Connect the qualified public geometry/placement proposal chain to Build31 data-only sampling and HMAC frame authentication without allowing payload/HMAC evidence to select geometry. Build34 closes this on MQ with two authenticated payloads.
- [x] Define the Build35 private v4 physical-corpus protocol: MQ block-normalized control + two robust/strength-24 authenticated carriers, exact hashes/acquisition plan, scanner-first 300 ppi/dpi workflow and HMAC-only PASS criterion.
- [x] Acquire the Build35 physical scanner corpus; the office scanner produced 600-dpi color JPEG captures of control / marked-a / marked-b. This is real v4 physical evidence, but blind qualification is still open because A/B currently SAFE-REJECT geometrically.
- [ ] Promote v4 only after physical pilot detection beats negative controls and a valid HMAC-authenticated
      payload is recovered. Pilot confidence alone must never authenticate.


## v0.3.0 — Build37/38 physical v4 campaign

- [x] Add a full-page scanner-specific boundary initializer that ignores the filename label and treats the paper/artwork edge only as geometry.
- [x] Refine that seed with disjoint public-pilot proposal/validation and require a complete-pilot `(0,0)` origin/margin gate.
- [x] Freeze a five-geometry pilot-qualified ensemble before reading protected data; aggregate signed DCT margins and use Build36 soft Hamming without allowing HMAC to choose geometry.
- [x] Blindly HMAC-authenticate both existing marked scanner JPEGs while rejecting the unmarked control.
- [x] Build38: acquire the first nine-photo smartphone corpus from the strength-24 scanner-qualified paper set (control/A/B x front/mild/angle) and preserve the native ~200 MP originals privately.
- [x] Use reference-assisted camera diagnostics to separate geometry from data-channel strength; strength 24 remains outside the Hamming/HMAC envelope even after dense residual registration.
- [x] Print the Build38 phone-specific robust/strength-48 control/A/B pack and acquire front/mild/angle native phone files for both marked carriers plus control.
- [x] Prove on the strength-48 corpus that independently supplied/reference-assisted geometry brings multiple marked captures inside exact HMAC recovery range; this established geometry, not embedding energy, as the remaining blocker.
- [x] Determine whether phone-camera decoding needs more than Build34 projective geometry: the strength-24 corpus shows residual lens/local warp, but pilot-only local fitting can overfit texture and does not rescue the data channel.
- [x] Build39/40: add blind camera boundary/projective registration and a bounded pilot-only residual model with disjoint spatial proposal/validation; keep HMAC out of geometry selection.

## v0.3.0 — Build36 physical v4 campaign

- [x] Prove channel sufficiency independently of blind geometry: reference-assisted registration + Build36 soft Hamming authenticates both existing marked scans exactly; keep this explicitly non-normative.
- [x] Recover the same Build35 scans with a **blind scanner registration** front-end using only public boundary/pilot evidence; Build37 rejects the control and HMAC-authenticates both marked scans.
- [x] Qualify the physical scanner corpus only after blind HMAC recovery succeeds without original-carrier registration; Build37 closes this scanner-specific gate. This does not yet make Format v4 normative.

## v0.3.0 — Build35 physical v4 campaign

- [x] Reset the active private corpus to explicit anonymous LQ/MQ/HQ entries with dimensions/SHA-256 (Build32).
- [x] Close the active MQ blind projective+crop + authenticated recovery blocker without lowering Build29 gates (Build34).
- [x] Expose the Build34 projective decoder through `ExperimentalV4ExtractProjective` / `v4-extract-projective` with explicit canonical dimensions (Build35).
- [x] Define the scanner-first MQ fixture pack: control + `v4-b35-phys-a` + `v4-b35-phys-b`, robust/strength 24, public test key, exact hashes and acquisition plan.
- [x] Print the Build35 pack at actual size / 300 ppi. The available office scanner cannot produce color lossless PNG/TIFF; retain the 600-dpi JPEG captures as the first physical corpus and document the deviation from the planned acquisition mode.
- [x] Require both marked scans to HMAC-authenticate their exact payloads and the control to reject; Build37 achieves this blindly on the original 600-dpi JPEG captures.
- [x] Acquire frontal/mild/angle phone captures of the same Build35 strength-24 sheets; preserve them as historical phone-channel evidence.
- [x] Acquire the new Build38 strength-48 phone qualification pack; preserve it separately from the strength-24 corpus.

- [x] Build39: add bounded phone downsampling, perspective-tolerant artwork boundary estimation, spatially split public-pilot projective refinement, `ExperimentalV4ExtractPhone` and `v4-extract-phone`.
- [x] Build40: implement and synthetically qualify the bounded pilot-only residual field after the Build39 homography.
- [x] Run the Build40 decoder over all nine strength-48 smartphone originals and preserve the stage-by-stage matrix (boundary, projective basin, proposal/held-out pilot, origin, residual fit/application, data confidence, soft-Hamming attempts and HMAC). The matrix identified global phone basin recovery as the next blocker.
- [x] Build41: improve only the global smartphone projective-basin stage identified by the complete Build40 matrix. The staged physical milestone is closed: all controls reject, `A/angle` authenticates A and `B/front` authenticates B, with no encoder/ECC/HMAC/residual change.
- [x] Build42: preserve Build41 geometry unchanged and recover `A/mild` only in the post-geometry data plane using the complete qualified bank, deterministic three-geometry ensembles and bounded soft-Hamming list decoding. `A/angle` and `B/front` remain direct passes; all controls still reject before data decode.
- [ ] Build44+: after deterministic ingest is physically qualified, broaden the blind smartphone **geometry** envelope for B/mild and B/angle while preserving A/front/A/mild/A/angle/B/front and all control rejections. Do not increase strength, weaken pilot floors or enlarge the residual model without new evidence.

## Build45 phone research

- [x] Freeze Build44 deterministic JPEG ingest and Go 1.26 qualification as the new baseline.
- [x] Add structured blind failure decomposition for B/mild/B/angle without changing production decisions.
- [x] Expose the complete Build43 six-side-pair proposal ranking for diagnostic inspection.
- [x] Add an isolated supplied-geometry oracle that cannot be reached from `v4-extract-phone`.
- [x] Confirm in retained-corpus lab analysis that B/mild and B/angle authenticate exactly under independently supplied geometry.
- [x] Run `make v4-build45-phone-diagnostic` on the qualification host and record the blind B/mild/B/angle matrix: 32/28 frozen candidates, exactly one held-out-qualified candidate in each, production Build42 bank 0.
- [x] Use the B/mild blind matrix to isolate the next question: determine whether the single held-out-qualified Build43 candidate is data-viable before changing geometry. Keep B/angle informational.
- [ ] Do not modify strength 48, Format-v4 framing, pilot, ECC/Hamming, whitening/HMAC or Build42 data-list recovery unless new channel evidence contradicts the Build45 oracle result.


## Build46 qualified-geometry handoff research

- [x] Preserve Build44/Build43 production thresholds unchanged; do not promote a singleton into the production decoder.
- [x] Add a diagnostic-only single-candidate Build42 list decode after Build43 freeze + held-out qualification.
- [x] Record the production quorum explicitly: two qualified geometries for direct phone acceptance, three for the Build42 list bank.
- [x] Add post-hoc oracle corner-distance comparison that cannot influence blind search or ranking.
- [x] Run `make v4-build46-phone-handoff-diagnostic` on the qualification host: both B singletons fail HMAC; B/mild is ~62.7 px mean from oracle and B/angle is a distant false basin.
- [x] Reject lowering quorum: B/mild singleton HMAC fails. Use oracle distance to study proposal breadth/refinement instead.
- [ ] Keep strength 48, Format-v4 framing, pilot, ECC/Hamming, whitening/HMAC and Build42 data-list rules frozen unless the Build46 evidence contradicts the Build45 oracle.

### Build47 research

- [x] Run the corrected three-tier frozen-bank study on B/mild (primary) and B/angle (informational).
- [x] B/mild selected-pair depth does not improve oracle-nearest error (59.734 px); all-pair extension finds 34.049 px from `top+left`, side-pair rank 3, cell rank 0.
- [x] B/angle remains a distant false-basin case (5499.824 px production -> 5448.391 px all-pair extension) and stays informational.
- [x] Do not promote six side pairs or the 128-bank directly: qualified candidates grow sharply while the oracle-nearest B/mild candidate still does not qualify.

### Build48 research

- [x] Add proposal-only local projective refinement over at most two independently selected seeds per side-pair rank.
- [x] Freeze every refined geometry before held-out/public-pilot qualification; keep key/HMAC and oracle unavailable during generation/refinement.
- [x] Add pre/post refinement source-quadrilateral telemetry and private post-hoc oracle comparison.
- [x] Run `make v4-build48-phone-local-refine-diagnostic` on the qualification host: B/mild top2 selected 12 seeds, 4 -> 7 qualified, HMAC 0, nearest selected geometry 46.202 -> 44.436 px.
- [x] Primary gate resolved negatively: the Build47 34.049 px basin was not selected by top2, so Build48 could not test its refinement reach.
- [x] Keep B/angle informational; it remained a distant false basin.
- [x] Keep strength 48, Format-v4, pilot, ECC/Hamming, whitening/HMAC, Build42 and Build43 production thresholds frozen.


### Build49 research

- [x] Run `make v4-build49-phone-proposal-ranking-diagnostic` on the qualification host.
- [x] B/mild oracle-nearest 34.049 px candidate is raw-proposal rank 3 within `top+left` (global 17), so top4 retains it.
- [x] Alternative proposal-only observables rank that candidate worse: fold-min 7, balanced 12, tile-consistency 13. Keep raw proposal unchanged.
- [x] Keep Build43/42 production, all thresholds/quorums, strength 48, Format-v4, ECC/Hamming and HMAC frozen.


### Build50 research

- [x] Run `make v4-build50-phone-top4-refine-test` on the qualification host.
- [x] Run `make v4-build50-phone-top4-refine-diagnostic` with retained B/mild/B-angle.
- [x] B/mild top4 includes the 34.049 px rank-3 basin, but refinement moves the nearest result to 43.168 px; qualified count rises 7 -> 11 and HMAC remains 0.
- [x] Do not change production seed depth, proposal scoring or qualification thresholds from Build50: the immediate question is local score-surface/refinement behavior.


### Build51 research

- [x] Run `make v4-build51-phone-surface-test` on the qualification host.
- [x] Run `make v4-build51-phone-surface-diagnostic` with retained B/mild and B/angle captures.
- [x] Inspect the post-hoc B/mild target seed trajectory and fixed local stencil: 9 accepted moves, 2 oracle-improving / 7 worsening; standard refinement 34.049 -> 43.168 px.
- [x] Classify B/mild as `optimizer-opportunity`: `corner-2-x +2` improves proposal by +0.006966 and oracle error by -1.607 px to 32.442 px.
- [x] Keep B/angle informational and classify it `proposal-surface-misaligned`; no better-both stencil point exists.
- [x] Keep Build43/42 production, strength 48, Format-v4, pilot, ECC/Hamming, whitening/HMAC and qualification/quorum thresholds frozen.

### Build52 research

- [x] Close the Build51 trajectory/stencil interpretation without changing the proposal score.
- [x] Implement a bounded proposal-only `2 px -> 1 px` restart from every untouched top4 seed and retain every accepted intermediate state.
- [x] Freeze baseline + restart geometry before held-out/full-pilot annotation; keep key/HMAC and SIFT/reference oracle strictly downstream.
- [x] Add `v4-diagnose-phone-restart`, source regressions, private diagnostic script and Build52 documentation.
- [x] Run `make v4-build52-phone-optimizer-test` on the qualified Go 1.26.0 host: PASS.
- [x] Run `make v4-build52-phone-optimizer-diagnostic` on retained B/mild and B/angle.
- [x] Analyze Build52: B/mild = `optimizer-partial-gain` (31.747 px best retained; 32.442 px best qualified; no HMAC), B/angle = `optimizer-no-gain`.
- [x] Do not promote Build52 into production; HMAC remains zero, so the Build44 production gate is not rerun as a promotion test.

### Build53 research

- [x] Keep the Build50/51/52 top4 seed bank, proposal score and Build44 production path frozen.
- [x] Implement proposal-only 2px root retention plus a complete +/-1px single-coordinate locality probe.
- [x] At 1px coordinate-local roots only, scan all 112 coupled +/-1px two-coordinate moves and retain at most the top 8 proposal-improving states by proposal alone.
- [x] Freeze all root/pair geometry before held-out/full-pilot qualification; keep secret key/HMAC and SIFT/reference oracle strictly downstream.
- [x] Add `v4-diagnose-phone-pair-escape`, source regressions, private diagnostic script and Build53 documentation.
- [x] Run `make v4-build53-phone-pair-escape-test` on the qualified Go 1.26.0 host: PASS.
- [x] Run `make v4-build53-phone-pair-escape-diagnostic` on retained B/mild and B/angle.
- [x] Classify Build53: B/mild = `pair-qualified-gain` (four qualified pair states; best post-hoc pair 31.842 px; no HMAC), B/angle target = `pair-not-triggered`.
- [x] Do not promote Build53 into production; HMAC remains zero.

### Build54 research

- [x] Keep Build53 pair generation, the top4 seed bank, proposal score and Build44 production path frozen.
- [x] Continue every retained Build53 pair state with bounded 1px proposal-only coordinate descent and retain every accepted intermediate.
- [x] Freeze all root/pair/continuation geometry before held-out/full-pilot qualification; keep secret key/HMAC and SIFT/reference oracle strictly downstream.
- [x] Add `v4-diagnose-phone-pair-continue`, source regressions, private diagnostic script and Build54 documentation.
- [ ] Run `make v4-build54-phone-pair-continuation-test` on the qualified Go 1.26.0 host.
- [x] Run `make v4-build54-phone-pair-continuation-diagnostic` on retained B/mild and B/angle; diagnostic artifacts received and analyzed.
- [x] Classify Build54: B/mild = `continuation-qualified-gain` (best qualified retained state 30.539 px; no HMAC), B/angle target = `continuation-not-triggered`.
- [x] Do not promote Build54: no blind retained continuation state authenticates, so no production candidate is opened from Build54.


### Build55 research

- [x] Keep Build54 continuation, top4 seed bank, proposal score and Build44 production path frozen.
- [x] Add a complete independent +/-1px sibling stencil around every retained Build54 continuation state; all 16 moves must start from the same parent geometry.
- [x] Freeze all sibling geometry before held-out/full-pilot qualification; keep key/HMAC and SIFT/reference oracle strictly downstream.
- [x] Add `v4-diagnose-phone-sibling-stencil`, source regressions, private diagnostic script and Build55 documentation.
- [ ] Run `make v4-build55-phone-sibling-stencil-test` on the qualified Go 1.26.0 host.
- [x] Run `make v4-build55-phone-sibling-stencil-diagnostic` on retained B/mild and B/angle; diagnostic artifacts received and analyzed.
- [x] Classify Build55: B/mild = `sibling-qualified-gain` (91 target siblings; best qualified 29.778 px; no HMAC), B/angle target = `sibling-not-triggered`.
- [x] Do not promote Build55: no blind frozen sibling authenticates, so no production candidate is opened from Build55.

### Build56 research

- [x] Keep the complete Build55 sibling bank, top4 seed bank, proposal score and Build44 production path frozen.
- [x] Probe every frozen Build55 sibling with the complete 16-neighbor independent +/-1px single-coordinate stencil using proposal only.
- [x] Only for siblings with zero proposal-improving single-coordinate neighbors, scan all 112 coupled +/-1px two-coordinate moves and retain at most eight proposal-improving states by proposal alone.
- [x] Freeze all second-pair geometry before held-out/full-pilot qualification; keep key/HMAC and SIFT/reference oracle strictly downstream.
- [x] Add `v4-diagnose-phone-sibling-pair-escape`, source regressions, private diagnostic script and Build56 documentation.
- [ ] Run `make v4-build56-phone-sibling-pair-escape-test` on the qualified Go 1.26.0 host.
- [ ] Run `make v4-build56-phone-sibling-pair-escape-diagnostic` on retained B/mild and B/angle.
- [ ] Classify B/mild as `sibling-pair-recovery`, `sibling-pair-qualified-gain`, `sibling-pair-geometric-gain`, `sibling-pair-proposal-only` or `sibling-pair-not-triggered`; keep B/angle informational.
- [ ] Only if a blind frozen second-pair state authenticates, design a separate production candidate and rerun the complete unchanged Build44 physical gate before promotion.

