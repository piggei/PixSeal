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


## v0.3.0 — Build36 physical v4 campaign

- [x] Prove channel sufficiency independently of blind geometry: reference-assisted registration + Build36 soft Hamming authenticates both existing marked scans exactly; keep this explicitly non-normative.
- [ ] Recover the same Build35 scans with a **blind scanner registration** front-end using only public boundary/lattice/pilot evidence; control must remain rejected.
- [ ] Promote the physical scanner corpus to a normative v4 PASS only after blind HMAC recovery succeeds without original-carrier registration.

## v0.3.0 — Build35 physical v4 campaign

- [x] Reset the active private corpus to explicit anonymous LQ/MQ/HQ entries with dimensions/SHA-256 (Build32).
- [x] Close the active MQ blind projective+crop + authenticated recovery blocker without lowering Build29 gates (Build34).
- [x] Expose the Build34 projective decoder through `ExperimentalV4ExtractProjective` / `v4-extract-projective` with explicit canonical dimensions (Build35).
- [x] Define the scanner-first MQ fixture pack: control + `v4-b35-phys-a` + `v4-b35-phys-b`, robust/strength 24, public test key, exact hashes and acquisition plan.
- [x] Print the Build35 pack at actual size / 300 ppi. The available office scanner cannot produce color lossless PNG/TIFF; retain the 600-dpi JPEG captures as the first physical corpus and document the deviation from the planned acquisition mode.
- [ ] Require both marked scans to HMAC-authenticate their exact payloads and the control to reject.
- [ ] Only after the scanner gate, acquire frontal-phone and perspective-phone captures of the same three printed sheets.
