## v0.3.0-build64 — 2026-09-24

- Closes Build63 as `fourth-pair-continuation-recovery`: B/mild freezes 937 fourth-pair continuation states, 935 qualify and exactly one authenticates; the authenticated state has proposal 0.302815, validation 0.240486 and post-hoc oracle error 28.868 px. B/angle freezes 6198 continuation states with 0 qualified and 0 authenticated.
- Adds the first production-candidate deep recovery fallback. Every Build44-qualified path remains first and unchanged; Build64 runs only after the historical phone path fails or cannot form its required geometry ensemble.
- Reproduces the complete Build63-derived geometry bank proposal-only, freezes it before held-out/full-pilot qualification, and attempts protected-data/list decoding only after the complete qualified subset is fixed. HMAC remains final authentication only.
- Adds Build64 recovery telemetry to `v4-extract-phone`, `v4-build64-phone-recovery-test`, a complete nine-photo `v4-build64-phone-physical-test`, and `docs/V4_BUILD64_PHONE_RECOVERY_CANDIDATE.md`.
- The Build64 physical gate promotes B/mild from informational reject to required authenticated recovery while requiring B/angle and all three controls to remain rejected. Build44 remains the latest qualified production milestone until that gate passes.

## v0.3.0-build63 — 2026-09-24

- Closes Build62 as `fourth-pair-qualified-gain`: B/mild retains 492 fourth-pair states, 491 qualified and 0 authenticated; the best target state reaches 24.056 px from its exact 25.641 px parent. B/angle retains 3746 fourth-pair states with 0 qualified and 0 authenticated.
- Adds proposal-only 1px continuation from every frozen Build62 fourth-pair state, retaining every accepted intermediate before qualification/HMAC.
- Adds `v4-diagnose-phone-fourth-pair-continue`, `v4-build63-phone-fourth-pair-continuation-test`, private `v4-build63-phone-fourth-pair-continuation-diagnostic`, and `docs/V4_BUILD63_FOURTH_PAIR_CONTINUATION.md`.
- Production Build44 behavior remains unchanged.

## v0.3.0-build62 — 2026-09-24

- Closes Build61 on the qualified Go 1.26.0 host as `third-pair-sibling-qualified-gain`: B/mild freezes 593 sibling states, 591 qualify, none authenticate, and the best qualified sibling reaches 24.614 px from its exact 25.359 px parent while proposal improves; B/angle freezes 1902 siblings with 0 qualified and 0 authenticated.
- Adds a research-only fourth-pair escape stage: every frozen Build61 third-pair sibling is probed with the complete independent +/-1px one-coordinate stencil, and only proposal-local siblings receive the bounded 112-combination coupled pair stencil with at most eight retained proposal-ranked states.
- Freezes the complete B/mild+B/angle Build62 fourth-pair bank before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry remains post-hoc only.
- Adds `v4-diagnose-phone-fourth-pair-escape`, `v4-build62-phone-fourth-pair-escape-test`, private `v4-build62-phone-fourth-pair-escape-diagnostic`, and `docs/V4_BUILD62_FOURTH_PAIR_ESCAPE.md`.
- Keeps Build44 as the latest qualified production milestone; no encoder, wire-format, pilot, proposal-score, threshold, ECC/Hamming, whitening or HMAC-domain changes.

## v0.3.0-build61 — 2026-09-24

- Closes Build60 on the qualified Go 1.26.0 host as `third-pair-continuation-qualified-gain`: B/mild retains 406 post-third-pair continuation states, 397 qualified and 0 authenticated, with the best qualified state at 25.150 px; target B/angle remains not triggered, while two non-target B/angle continuation states pass qualification at ~5665 px and still fail HMAC.
- Adds a research-only third-pair continuation sibling stencil: every frozen Build60 continuation state is probed independently in all eight coordinates at +/-1 px from the identical parent using only the unchanged proposal objective.
- Freezes the complete B/mild+B/angle Build61 sibling bank before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry remains external and post-hoc.
- Adds `v4-diagnose-phone-third-pair-sibling`, `v4-build61-phone-third-pair-sibling-stencil-test`, private `v4-build61-phone-third-pair-sibling-stencil-diagnostic`, and `docs/V4_BUILD61_THIRD_PAIR_SIBLING_STENCIL.md`.
- Tracks full-bank sibling qualification counts explicitly so the first non-target B/angle qualification leakage observed in Build60 cannot be hidden by target-only summaries.
- Changes no production encoder/decoder path, proposal score, threshold, quorum, pilot, ECC/Hamming, whitening or HMAC domain. Build44 remains the qualified production baseline.

## v0.3.0-build60 — 2026-09-24

- Closes the qualified-host Build59 study as `third-pair-qualified-gain`: B/mild freezes 180 third-pair states, 175 qualified and 0 authenticated; target candidate 10 reaches a best qualified post-hoc geometry of 25.407 px, 1.664 px better than its exact 27.071 px Build58 sibling parent.
- Records B/angle as the control: target candidate 18 remains `third-pair-not-triggered`; the complete B/angle Build59 third-pair bank contains 914 states, 0 qualified and 0 authenticated.
- Adds a research-only continuation from **every** retained Build59 third-pair state using the unchanged proposal-only 1px coordinate descent for at most eight passes, retaining every accepted intermediate.
- Freezes the complete B/mild+B/angle Build60 continuation bank before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry remains external and post-hoc.
- Adds `v4-diagnose-phone-third-pair-continue`, `v4-build60-phone-third-pair-continuation-test`, private `v4-build60-phone-third-pair-continuation-diagnostic`, and `docs/V4_BUILD60_THIRD_PAIR_CONTINUATION.md`.
- No production encoder/decoder, score, threshold, quorum, pilot, ECC/Hamming or HMAC domain changes. Build44 remains the qualified production baseline.

## v0.3.0-build59 — 2026-09-24

- Closes the qualified-host Build58 study as `second-pair-sibling-qualified-gain`: B/mild contains 275 frozen sibling states, 271 qualified and 0 authenticated; target candidate 10 improves from 26.968 px to 26.149 px while proposal rises from 0.276536 to 0.285435.
- Records B/angle as the control: 542 Build58 sibling states, 0 qualified and 0 authenticated; target candidate 18 remains not triggered.
- Adds a research-only local-max probe around every frozen Build58 sibling and applies the unchanged bounded 112-combination coupled +/-1px pair stencil only at proposal-local siblings, retaining at most eight states by proposal.
- Freezes the complete B/mild+B/angle Build59 bank before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry remains external and post-hoc.
- Adds `v4-diagnose-phone-third-pair-escape`, `v4-build59-phone-third-pair-escape-test`, private `v4-build59-phone-third-pair-escape-diagnostic`, and `docs/V4_BUILD59_THIRD_PAIR_ESCAPE.md`.
- No production encoder/decoder, score, threshold, quorum, pilot, ECC/Hamming or HMAC domain changes. Build44 remains the qualified production baseline.

## v0.3.0-build58 — 2026-09-24

### Research / Format-v4 phone geometry

- Closes the qualified-host Build57 study: B/mild is `second-pair-continuation-qualified-gain`; target candidate 10 retains 154 post-second-pair continuation states, 147 qualified and 0 authenticated. The best qualified intermediate reaches 26.968 px post-hoc oracle error, `-2.581 px` relative to its exact Build56 second-pair parent and `-7.080 px` relative to the original 34.049 px seed.
- Records B/angle as the control: target candidate 18 remains `second-pair-continuation-not-triggered`; the complete B/angle Build57 bank contains 347 continuation states, 0 qualified and 0 authenticated.
- Adds a research-only independent +/-1px sibling stencil from every retained Build57 post-second-pair continuation state. All 16 siblings use the identical frozen parent; every proposal-improving sibling is retained.
- Freezes the complete B/mild+B/angle Build58 bank before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry remains external and post-hoc.
- Adds `v4-diagnose-phone-second-pair-sibling`, `v4-build58-phone-second-pair-sibling-stencil-test`, private `v4-build58-phone-second-pair-sibling-stencil-diagnostic`, and `docs/V4_BUILD58_SECOND_PAIR_SIBLING_STENCIL.md`.
- No production path changes. Build44 remains the latest qualified milestone.

## v0.3.0-build57 — 2026-09-24

### Research / Format-v4 phone geometry

- Closes the qualified-host Build56 study: B/mild is `sibling-pair-qualified-gain`; target candidate 10 retains 65 second-pair states and a qualified second-pair state improves oracle error `29.778 -> 28.514 px` while proposal rises `0.266832 -> 0.273793`; HMAC remains false. B/angle target remains `sibling-pair-not-triggered`.
- Records complete-bank Build56 counts: B/mild freezes 68 second-pair states (51 qualified, 0 authenticated); B/angle freezes 176 (0 qualified, 0 authenticated).
- Adds a research-only continuation bank from **every** retained Build56 second-pair state using the unchanged 1px proposal-only coordinate descent for at most eight passes, retaining every accepted intermediate rather than privileging the endpoint.
- Freezes the complete B/mild+B/angle second-pair + continuation bank before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry remains external and post-hoc.
- Adds `v4-diagnose-phone-sibling-pair-continue`, `v4-build57-phone-sibling-pair-continuation-test`, private `v4-build57-phone-sibling-pair-continuation-diagnostic`, and `docs/V4_BUILD57_SECOND_PAIR_CONTINUATION.md`.
- No production decoder, encoder, Format-v4 field, pilot, strength, ECC/Hamming, whitening/HMAC domain, score, threshold or quorum changes. Build44 remains the latest qualified milestone.

## v0.3.0-build56 — 2026-09-24

### Research / Format-v4 phone geometry

- Closes the qualified Build55 host study: B/mild is `sibling-qualified-gain` with 91 target sibling states and a best held-out-qualified sibling at 29.778 px; HMAC remains false. The B/angle target is `sibling-not-triggered`.
- Records complete-bank Build55 counts: B/mild freezes 93 sibling states (76 qualified, 0 authenticated); B/angle freezes 106 (0 qualified, 0 authenticated).
- Adds a research-only second pair-escape diagnostic: every frozen Build55 sibling is first tested against all 16 independent +/-1px single-coordinate neighbors using proposal only; only proposal-local siblings receive the deterministic 112-combination coupled +/-1px two-coordinate stencil.
- Retains at most eight proposal-improving pair states per proposal-local sibling, ordered only by the unchanged proposal score. The complete bank is frozen before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference oracle remains external and post-hoc.
- Adds `v4-diagnose-phone-sibling-pair-escape`, `v4-build56-phone-sibling-pair-escape-test`, private `v4-build56-phone-sibling-pair-escape-diagnostic`, and `docs/V4_BUILD56_SIBLING_PAIR_ESCAPE.md`.
- Design probe (not qualification evidence): the Build55 29.778 px sibling is coordinate-local under all 16 single moves, yet five coupled proposal-improving escapes exist; one retained state reaches about 28.514 px while remaining qualified, with HMAC still false.
- No production decoder, encoder, Format-v4 field, pilot, strength, ECC/Hamming, whitening/HMAC domain, score, threshold or quorum changes. Build44 remains the latest qualified milestone.

## v0.3.0-build55 — 2026-09-24

### Research / Format-v4 phone geometry

- Closes the Build54 host diagnostic: B/mild is `continuation-qualified-gain` with 34 retained continuation states and a best qualified intermediate at 30.539 px; HMAC remains false. The B/angle target is `continuation-not-triggered`.
- Adds a research-only independent sibling stencil around every retained Build54 continuation state: all 8 corner coordinates are evaluated at both `-1 px` and `+1 px` from the exact same parent geometry.
- Retains every proposal-improving sibling before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference oracle remains external and post-hoc.
- Adds `v4-diagnose-phone-sibling-stencil`, `v4-build55-phone-sibling-stencil-test`, private `v4-build55-phone-sibling-stencil-diagnostic`, and `docs/V4_BUILD55_SIBLING_STENCIL.md`.
- Qualified Go 1.26.0 host result: B/mild target candidate 10 freezes 91 siblings; the best qualified sibling improves oracle error `30.539 -> 29.778 px` and proposal `0.257772 -> 0.266832`, with HMAC false. B/angle target remains `sibling-not-triggered`. Complete-bank counts are 93/76/0 frozen-qualified-authenticated on B/mild and 106/0/0 on B/angle.
- No production decoder, encoder, Format-v4 field, pilot, strength, ECC/Hamming, whitening/HMAC domain, score, threshold or quorum changes. Build44 remains the latest qualified milestone.

## v0.3.0-build54 — 2026-09-24

### Research / Format-v4 phone geometry

- Closes the qualified Build53 host study: B/mild is `pair-qualified-gain` (one 1px coordinate-local root, four retained pair states, all four qualified; best post-hoc pair geometry 31.842 px; no HMAC), while the B/angle target is `pair-not-triggered`.
- Adds a research-only post-pair continuation diagnostic that reproduces the unchanged Build53 pair bank and continues every retained pair state with bounded 1px proposal-only coordinate descent.
- Retains every accepted continuation intermediate rather than privileging the greedy endpoint. The complete root/pair/continuation bank is frozen before held-out/full-pilot qualification and diagnostic HMAC; SIFT/reference oracle remains external and post-hoc.
- Adds `v4-diagnose-phone-pair-continue`, `v4-build54-phone-pair-continuation-test`, private `v4-build54-phone-pair-continuation-diagnostic`, and `docs/V4_BUILD54_PAIR_CONTINUATION.md`.
- Qualification-host Build54 result: B/mild retains 34 continuation states and reaches a held-out-qualified 30.539 px intermediate, but no continuation state authenticates; B/angle remains `continuation-not-triggered`. Later proposal ascent can move away from the best retained geometry.
- No production decoder, encoder, Format-v4 field, pilot, strength, ECC/Hamming, whitening/HMAC domain, geometry threshold or quorum changes. Build44 remains the latest qualified milestone.

## v0.3.0-build53 — 2026-09-24

### Research / Format-v4 phone geometry

- Closes the qualified Build52 host study: B/mild is `optimizer-partial-gain` (best retained oracle 31.747 px; best retained qualified state 32.442 px; no HMAC), while B/angle is `optimizer-no-gain`.
- Records the qualified Build53 host result: the B/mild target has one coordinate-local root and four retained pair states; all four are qualified, pair rank 2 reaches 31.842 px, and no pair state authenticates. The B/angle target does not trigger pair scanning.
- Adds a research-only coupled pair-escape diagnostic that preserves the unchanged top4 seed bank and proposal score, retains proposal-only 2px roots, identifies 1px coordinate-local roots, and scans bounded coupled +/-1px two-coordinate moves only at those roots.
- Pair generation/ranking is proposal-only. At most eight pair states per local root are retained; the full geometry bank is frozen before held-out/full-pilot qualification and diagnostic HMAC. SIFT/reference oracle remains external and post-hoc.
- Adds `v4-diagnose-phone-pair-escape`, `v4-build53-phone-pair-escape-test`, private `v4-build53-phone-pair-escape-diagnostic`, and `docs/V4_BUILD53_PAIR_ESCAPE.md`.
- No production decoder, encoder, Format-v4 field, pilot, strength, ECC/Hamming, whitening/HMAC domain, geometry threshold or quorum changes. Build44 remains the latest qualified milestone.

## v0.3.0-build52 — 2026-09-24

- Closes the Build51 qualification-host local-surface study: B/mild candidate 10 (`top+left`, pair rank 3, seed rank 3) starts at 34.049 px mean oracle error and the ordinary `16,8,4,2,1` coordinate descent ends at 43.168 px after 9 accepted moves; only 2 accepted moves improve oracle error and 7 worsen it.
- Records the decisive Build51 stencil sample `corner-2-x +2`: proposal `0.201986 -> 0.208953` (`+0.006966`) while post-hoc oracle error improves `34.049 -> 32.442 px` (`-1.607 px`). B/mild is therefore classified `optimizer-opportunity`; B/angle remains `proposal-surface-misaligned`.
- Adds a research-only fine-restart optimizer diagnostic that keeps the unchanged Build50/51 top4 seeds and proposal score, independently restarts each untouched seed at `2 px -> 1 px`, and retains the seed plus every accepted proposal-improving intermediate state.
- Adds `v4-diagnose-phone-restart`, `v4-build52-phone-optimizer-test`, private `v4-build52-phone-optimizer-diagnostic`, and `docs/V4_BUILD52_FINE_RESTART.md`.
- Freezes the complete baseline + restart state bank before held-out/full-pilot qualification; key/HMAC is downstream of qualification and SIFT/reference oracle remains external and post-hoc.
- Production Build43/42 behavior, Build44 qualified baseline, Format-v4, strength 48, pilot, ECC/Hamming, whitening/HMAC domains and all production thresholds/quorums remain unchanged.

## v0.3.0-build51 — 2026-09-23

- Adds research-only local proposal-surface/refinement-trajectory observability for the blind top4-per-side-pair seeds established by Build50.
- Adds `v4-diagnose-phone-surface`, which records every evaluated Build41 coordinate-descent +/- move with proposal score and accepted/rejected state, without changing the refiner.
- Adds a deterministic up-to-53-point-per-seed local stencil spanning individual corner axes plus bounded translation/scale/shear/perspective-like modes. Stencil generation is proposal-only and oracle-free.
- Freezes complete trace/stencil geometry before held-out validation annotation. Full-pilot qualification and diagnostic HMAC remain limited to original/final seed states; SIFT/reference oracle remains external and post-hoc.
- Records the qualification-host Build50 result: B/mild top4 includes the 34.049 px seed, but refinement moves the nearest candidate to 43.168 px; qualified count rises 7 -> 11 and HMAC stays 0. B/angle remains a distant false basin.
- Records the completed Build51 host result: for B/mild candidate 10 the refiner accepts 9 moves (2 oracle-improving, 7 oracle-worsening); a rejected `corner-2-x +16` state reaches 25.172 px but lowers proposal, while stencil point `corner-2-x +2` improves both proposal (`+0.006966`) and oracle (`-1.607 px`), establishing `optimizer-opportunity`. B/angle has no better-both stencil point and remains `proposal-surface-misaligned`.
- Adds `v4-build51-phone-surface-test`, private `v4-build51-phone-surface-diagnostic`, and `docs/V4_BUILD51_LOCAL_SURFACE.md`.
- Production Build43/42 behavior, Format-v4, strength 48, pilot, ECC/Hamming, whitening/HMAC domains and Build44 qualification remain frozen.

## v0.3.0-build50 — 2026-09-23

- Adds research-only top-4-per-side-pair local projective refinement for the difficult smartphone B cases.
- Preserves the unchanged Build49/raw proposal score; Build49 showed the B/mild 34.049 px oracle-nearest basin is proposal rank 3 within its pair, while tested alternative observables rank it worse.
- Keeps historical Build48 fixed at top 2 seeds per pair and keeps Build43 production fixed at two side pairs / at most 32 frozen candidates.
- Adds `seed_rank_within_pair` to diagnostic refinement records so one blind top-4 run can be analyzed as exact nested top2/top4 subsets.
- Keeps all geometry creation/refinement proposal-only, freezes the bank before held-out qualification, uses the key only after qualification for diagnostic HMAC, and generates SIFT/reference oracle geometry only after both blind outputs exist.
- Adds `v4-diagnose-phone-refine4`, `v4-build50-phone-top4-refine-test`, private `v4-build50-phone-top4-refine-diagnostic`, and `docs/V4_BUILD50_TOP4_REFINEMENT.md`.
- Build44 remains the latest qualified milestone; no production decoder or Format-v4 wire-format change.

## v0.3.0-build49 — 2026-09-23

- Adds research-only proposal-ranking observability across the corrected Build47 extended bank.
- Records proposal fold agreement, per-tile stability, pair/cell scores and multiple proposal-only ranks without changing geometry.
- Adds `v4-diagnose-phone-ranking`; the command deliberately accepts no secret key.
- Adds post-hoc oracle rank, top-2/4/6/8-per-pair coverage and observable/error Spearman summaries in the private lab script.
- Records the Build48 host result: B/mild selected refinement improved 46.202 -> 44.436 px, 7 candidates qualified, 0 authenticated; the better 34.049 px Build47 basin was not selected.
- Build44 remains the latest qualified production milestone.

## v0.3.0-build48 — 2026-09-23

- Adds a research-only local projective refinement stage over the corrected Build47 extended side-pair bank; production Build43/42 behavior is unchanged.
- Selects at most two seeds per side-pair rank using proposal score only. Held-out validation, full-pilot qualification, key/HMAC and oracle geometry cannot influence seed selection or refinement.
- Refines all selected seeds on proposal folds, freezes the complete refined bank, then reports pre/post held-out validation, pilot score/margin/origin and diagnostic single-candidate HMAC only for post-qualified geometry.
- Adds private post-hoc oracle reporting for pre/post corner error after both blind B/mild and B/angle refinement JSON files are complete.
- Records the Build47 host result motivating the experiment: B/mild improves from 59.734 px oracle-nearest mean error in the production/depth tiers to 34.049 px only through side-pair rank 3 (`top+left`), while B/angle remains a distant false-basin case.
- Adds `v4-diagnose-phone-refine`, `v4-build48-phone-local-refine-test`, private `v4-build48-phone-local-refine-diagnostic`, and `docs/V4_BUILD48_LOCAL_REFINEMENT.md`.
- Keeps Build44 as the latest qualified milestone and leaves Format-v4, strength 48, pilot, ECC/Hamming, whitening/HMAC, Build42 and all production geometry thresholds frozen.

## v0.3.0-build47 — frozen candidate bank observability

- Added a diagnostic-only Build47 extended proposal bank with three explicit tiers: the exact Build43 production tier, deeper cells on the same selected side pairs, and lower-ranked side-pair extension. Production remains unchanged.
- Added per-candidate proposal, held-out, pilot, side-pair provenance, qualification and post-hoc oracle-distance telemetry.
- Kept single-candidate data/HMAC inspection behind unchanged Build43 qualification gates.
- Added a private Build47 study that completes blind diagnostics before generating SIFT/reference geometry.
- Build44 remains the latest qualified milestone; Build47 is research-only.

## v0.3.0-build46 — qualified-geometry handoff diagnostic

- Keeps Build44 as the latest qualified production milestone and changes no production phone decision.
- Records the Build45 qualification-host finding that B/mild and B/angle each retain exactly one Build43 held-out-qualified candidate, below the production two-geometry ensemble quorum.
- Adds `v4-diagnose-phone-handoff` to inspect already-qualified Build43 candidates, production quorum availability and a diagnostic-only single-candidate Build42 list decode.
- Adds post-hoc oracle corner-distance comparison that is computed only after blind search/qualification and cannot influence geometry.
- Adds `make v4-build46-phone-handoff-test` and private `make v4-build46-phone-handoff-diagnostic`.
- Leaves Format-v4, strength 48, pilot, Build42, ECC/Hamming, whitening/HMAC and all production thresholds unchanged.

## v0.3.0-build45 — phone failure decomposition

- Starts from the fully qualified Build44 deterministic-JPEG / Go 1.26 baseline without changing production watermark or phone-decoder decisions.
- Adds `v4-diagnose-phone`, a diagnostic-only view of Build41 qualification, the complete Build43 six-side-pair ranking, frozen/qualified geometry counts, Build42 data-list telemetry and final HMAC result.
- Adds deterministic failure classes: `geometry`, `qualification`, `data-channel`, and `recovered`.
- Adds an explicitly isolated reference-assisted lab oracle: OpenCV SIFT/RANSAC may supply one acquisition-space quadrilateral, but that geometry can be consumed only by the diagnostic command and can never become a production fallback.
- Preliminary retained-corpus oracle evidence recovers exact `v4-b38-phone-b` from both B/mild and B/angle at the first Build42 list frame, proving that the difficult B cases remain inside the protected data-channel envelope when geometry is correct.
- Keeps B/mild as the primary Build45 research target and B/angle informational until the blind failure matrix identifies the bounded geometry change justified by evidence.

## v0.3.0-build44 — 2026-09-23

- Starts deterministic JPEG ingest without changing Format-v4, the encoder, public pilot, carrier/data mapping, strength 48, Hamming/ECC, whitening/HMAC domains, Build41/43 geometry or the Build42 data/list decoder.
- Vendors a pure-Go pre-Go-1.26 JPEG decoder under `internal/jpeglegacy`, with the upstream Go BSD-style license retained, so JPEG rasterization no longer depends on the compiler standard library.
- Routes CLI JPEG `DecodeConfig` and full decode through the PixSeal-controlled decoder; PNG continues to use `image/png`. Magic-byte detection remains extension-independent.
- Adds a deterministic JPEG fixture with locked Y/Cb/Cr SHA-256 vectors plus a CLI-ingest regression. The decoder identity is exposed in phone diagnostics as `pixseal-jpeg-pre-go1.26-v1`.
- Qualifies Go 1.26.0 as the Build44 default toolchain after both deterministic JPEG regressions and the unchanged private Build43 smartphone matrix pass on the real Surface/WSL2 qualification host.
- Physical qualification result: controls 3/3 reject; A/front, A/mild, A/angle and B/front authenticate; B/mild and B/angle remain informational rejects.
- Promotes `make v4-build44-phone-physical-test` as the normal Build44 physical gate; retains the explicit `v4-build44-go126-*` targets as historical qualification regressions.
- Adds `v4-build44-jpeg-compat-test`, deterministic decoder fixtures and `docs/V4_BUILD44_DETERMINISTIC_JPEG.md`.

## v0.3.0-build43 — 2026-09-23

- Adds a proposal-only blind side-pair geometry fallback for smartphone Format-v4 recovery after Build41 geometry rejection.
- Preserves Format-v4 wire format, encoder, strength-48 qualification signal, pilot, ECC/Hamming, whitening, HMAC and the Build42 data/list decoder.
- Adds structural paper→artwork side proposals with conservative gating, multi-pair/angular/fine/complement diversity, a frozen bank capped at 32 candidates, and held-out qualification only after freeze.
- Recovers the private Build38 `A/front` smartphone capture while preserving the existing `A/mild`, `A/angle`, `B/front` passes and 3/3 control rejections.
- Adds Build43 CLI telemetry, `v4-build43-phone-side-pair-test`, opt-in `v4-build43-phone-physical-test`, and `docs/V4_BUILD43_PHONE_SIDE_PAIR.md`.
- Pins Build43 Makefile/Windows build commands to the qualified Go 1.25.1 toolchain and forces native rebuilds so a stale Go 1.26 binary cannot be reused.
- Documents the Go 1.26 `image/jpeg` decoder replacement discovered by the physical A/mild corpus: identical JPEG bytes rasterize differently and change the geometry result. Go 1.26+ is therefore not a qualified Build43 JPEG toolchain.
- Adds `docs/GO_TOOLCHAIN_COMPATIBILITY.md`; no watermark, Format-v4, ECC/HMAC or Build43 geometry retuning is made to compensate for the upstream decoder change.

## v0.3.0-build42 — 2026-09-18

- Keeps the Build41 smartphone geometry search, proposal/held-out split and acceptance thresholds unchanged.
- Retains the complete already-qualified Build41 geometry bank for a post-geometry data fallback when the normal two-hypothesis soft decoder fails authentication.
- Adds deterministic three-geometry margin ensembles and a bounded Hamming list decoder over the ten weakest ML word gaps; HMAC remains final frame authentication and never feeds geometry search/ranking.
- Reorders the phone pipeline to `Build41 global decode -> Build42 data fallback -> Build40 residual fallback`, avoiding expensive residual fitting when accepted global geometry is already sufficient.
- Physical result on the unchanged strength-48 corpus: all three controls REJECT before data decode; `A/angle` and `B/front` remain direct PASS; `A/mild` newly authenticates `v4-b38-phone-a` through Build42; A/front, B/mild and B/angle remain geometry rejects.
- Adds `v4-build42-phone-data-test`, `v4-build42-phone-physical-test`, Build42 CLI telemetry and `docs/V4_BUILD42_PHONE_DATA_LIST.md`.
- No encoder, pilot, strength, data layout, ECC, whitening/HMAC domain, scanner path or frozen Format-v3 change.

## v0.3.0-build41 — 2026-09-18

- Closes the first staged blind-smartphone physical milestone on the existing private strength-48 corpus without changing the Format-v4 encoder, locked pilot, 1120-position data partition, Hamming(7,4), whitening/HMAC domains, robust strength 48, Build40 residual model or frozen Format-v3 core.
- Adds a structure-only inner-artwork boundary seed for difficult free-camera framing and a bounded eight-coordinate projective-basin search anchored to that physical boundary.
- Uses a deterministic three-way spatial public-pilot split: folds 1+2 generate/refine/freeze the geometry shortlist; fold 0 remains completely held out until qualification. Payload bytes, ECC result, secret key and HMAC are unavailable to geometry generation/ranking.
- Adds proposal-only phase restoration and fine variants, then requires held-out/full-pilot qualification at canonical origin `(0,0)`. The final two-geometry ensemble must also satisfy a scale-normalized diversity floor so it brackets sub-pixel registration uncertainty instead of duplicating one local optimum.
- Adds `v4-build41-phone-basin-test`, an end-to-end synthetic marked/control gate over the public `ExperimentalV4ExtractPhone` API.
- Adds opt-in `v4-build41-phone-physical-test` for the private nine-photo strength-48 corpus. Build41 checkpoint result: all three controls REJECT; `A/angle` authenticates `v4-b38-phone-a`; `B/front` authenticates `v4-b38-phone-b`; `A/mild` reaches data decode but does not authenticate; the other three marked captures remain outside the accepted basin.
- The Build40 residual is attempted where applicable but is not applied in either physical PASS case, so the two authenticated results come from Build41 global registration plus the unchanged soft-Hamming/HMAC channel.
- Keeps the project source-available for noncommercial use under PolyForm Noncommercial License 1.0.0.

## Licensing update — 2026-09-18

- Changes the license for current and future source distributions from MIT to
  **PolyForm Noncommercial License 1.0.0**
  (`PolyForm-Noncommercial-1.0.0`).
- Describes the project as **source-available for noncommercial use** rather
  than OSI open source. Commercial use requires a separate commercial license
  from the applicable copyright holder(s).
- Adds `NOTICE` and `docs/LICENSING.md`, including redistribution requirements,
  historical-MIT continuity and the explicit separation of the software license
  from private qualification corpora.
- At the time of the relicensing checkpoint, this was a licensing/documentation-only change and `VERSION` remained `v0.3.0-build40`; no encoder, decoder, format, pilot, ECC, geometry or HMAC behavior changed in that relicensing step. Build41 was developed subsequently under the new license.
- Copies already distributed under MIT retain the MIT rights granted with those
  copies; the new license applies prospectively to distributions carrying it.

## v0.3.0-build40 — 2026-09-17

- Adds a bounded quadratic smartphone residual field on top of the Build39 homography; the maximum correction is six canonical pixels.
- Fits residual controls exclusively from checkerboard-A repetitions of the public Format-v4 pilot and keeps checkerboard-B repetitions held out for acceptance. Payload bytes, frame header, ECC outcome and HMAC are not geometry inputs.
- Rejects ambiguous local pilot peaks, robustly fits the remaining controls, and requires held-out validation gain plus complete-pilot `(0,0)` origin/score/margin before the field can be used.
- Integrates residual-aware protected-margin sampling into `ExperimentalV4ExtractPhone` / `v4-extract-phone`; if the residual path does not authenticate, the untouched Build39 ensemble is retried once. HMAC never scores or reorders geometry.
- Adds `v4-build40-phone-residual-test`: a smooth non-projective strength-48 synthetic channel is corrected from public-pilot evidence and reaches exact HMAC recovery; an unmarked control remains rejected.
- Records the private real-phone result as a safe negative: current Build39 basins do not yield a coherent <=6 px field, so the remaining blocker is global projective-basin selection rather than justification for a larger local warp.
- Keeps Format-v3 frozen, strength 48, the v4 pilot/data layout, Hamming(7,4) and HMAC domains unchanged.

## v0.3.0-build39 — 2026-09-17

Build39 starts blind smartphone registration on the completed strength-48 physical corpus without changing Format-v4 encoding, strength, pilot, data mapping, Hamming(7,4), HMAC domains or the frozen Format-v3 core.

- Adds experimental API `ExperimentalV4ExtractPhone` and CLI `v4-extract-phone`.
- Adds bounded internal downsampling for very large camera captures (long side capped at 4600 px for the registration/data working image).
- Adds a phone-specific artwork-boundary path: scanner boundary is reused only when it is safely interior; otherwise local-paper evidence plus edge-gradient refinement estimates the four artwork edges under perspective.
- Adds projective refinement driven only by the public v4 pilot. Proposal and held-out validation use spatially disjoint checkerboard tile partitions; payload/header/ECC/HMAC are excluded from geometry selection.
- Adds `v4-build39-phone-registration-test`, a synthetic registration checkpoint covering perspective boundary evidence, marked/control separation, the phone downsample bound and invalid canonical dimensions.
- The real strength-48 corpus confirms that the data channel itself is sufficient under reference-assisted registration (multiple marked captures authenticate by HMAC). Blind phone HMAC closure is **not** claimed in Build39: residual local/lens/print warp remains after a single homography and is the next research target.
- Keeps the corrected private-data `.gitignore` baseline, including explicit `/private-fixtures/` and `/v4-phone private/` exclusions and no `/internal/` exclusion.

## v0.3.0-build38 — 2026-09-17

Build38 starts smartphone-channel qualification after Build37 closed blind scanner recovery. It deliberately keeps the Format-v4 wire format, locked pilot, 1120-position data mapping, Hamming(7,4), HMAC domains and frozen Format-v3 core unchanged.

- Adds the first private 9-photo smartphone corpus protocol (`front`, `mild`, `angle` for control/A/B) and records the initial strength-24 result as a **negative but informative physical measurement**.
- Reference-assisted diagnostics isolate the channel limit: even after dense non-rigid registration, representative strength-24 captures retain about 14.5–17.9% protected coded-bit error and 20–42/256 post-soft-Hamming bit errors. This is far outside exact HMAC recovery and shows that geometry alone is not the remaining blocker.
- Rejects increasingly flexible pilot-only local warp fitting as a production direction for the strength-24 photographs because it can overfit natural image texture while the protected data plane remains weak.
- Defines a dedicated phone qualification carrier at robust **strength 48**, preserving every other Format-v4 parameter. Digital MQ comparison gives about 32.5 dB PSNR versus the control while remaining visually subtle in the qualification image.
- Adds `make v4-phone-fixtures`, producing private Build38 control/A/B print masters, SHA-256 manifest, acquisition plan and a phone-specific protocol. The output defaults to `v4-phone private/build38-generated/`.
- Adds `v4-build38-phone-channel-test` to protect the unchanged v4 framing/ECC path at strength 48 through aligned and JPEG-q82 round trips.
- Hardens `.gitignore`: explicitly excludes `private-fixtures/`, `v4-physical private/`, the new `v4-phone private/`, the historical private corpora and OS `Zone.Identifier` artifacts; `internal/` remains tracked.
- The Build35 strength-24 scanner prints remain the qualified scanner corpus and are not overwritten. Build38 requires a new phone-specific marked print pair before any blind camera decoder claim.

## v0.3.0-build37 — 2026-09-15

- Adds `ExperimentalV4ExtractScanner` and the explicit `v4-extract-scanner` CLI for full-page scanner captures with visible white paper around the printed artwork.
- Adds a paper/artwork boundary estimator that is used only as a geometry initializer, never as watermark/authentication evidence.
- Adds a deterministic two-stage scanner-affine refinement: a 18,750-hypothesis bounded coarse bank and a 15,625-hypothesis fine neighborhood. Pilot partition A proposes; disjoint partition B ranks.
- Requires the held-out winner to recover complete-pilot cyclic origin `(0,0)` with scanner pilot score >= 0.15 and margin >= 0.05. Proposal/validation floors are 0.30/0.15.
- Freezes the top five pilot-qualified geometries before any frame/key/HMAC work. Protected DCT margins are validation-weighted across that ensemble, then passed to the Build36 soft-Hamming decoder and unchanged v4 frame/HMAC.
- Adds `v4-build37-scanner-registration-test` with a synthetic full-page JPEG scanner channel and an unmarked control.
- Adds opt-in `v4-build37-physical-scanner-test` for private full-page captures; the private files remain excluded from source archives.
- First real blind physical scanner result on the existing Build35 corpus: control REJECT; marked A authenticates `v4-b35-phys-a`; marked B authenticates `v4-b35-phys-b`. No reprint, strength change, pilot change, ECC change or Format-v3 change is required.
- Documents scanner success as a channel-specific checkpoint only; smartphone frontal/perspective capture remains a separate future gate.

## v0.3.0-build36 — 2026-09-15

Build36 narrows the first real Format-v4 print/scan result to a data-reliability problem after accurate geometry. It does not change the locked pilot, encoder, strength, frame format, ECC code, Build29 acceptance floors or frozen Format-v3 core.

- Adds reliability-aware projective data sampling that preserves signed DCT margins instead of collapsing immediately to hard bits.
- Promotes deterministic maximum-likelihood Hamming(7,4) decoding to the accepted v4 projective data path, with the Build31/34 hard-decision decoder retained as fallback.
- Adds `v4-build36-soft-channel-test`, including a crafted two-weak-error Hamming word that hard decoding cannot recover but the soft decoder recovers deterministically.
- Records the first Build35 physical scanner corpus: office-scanner 600-dpi color JPEG acquisitions of control / marked-a / marked-b. The blind Build35 geometry path still rejects A/B, so normative physical qualification remains open.
- Records a strictly diagnostic reference-assisted result: with geometry supplied independently from the original carrier, both marked physical scans authenticate their exact HMAC-protected payloads under the new soft decoder. This is channel-capacity evidence only, not a blind decoder PASS.
- Leaves print strength at 24 and avoids an unnecessary ECC/format change until blind scanner registration is addressed.

# Changelog

## v0.3.0-build35 — 2026-09-15

Build35 is the physical-qualification handoff after Build34 closed the active MQ blind projective+crop blocker. It intentionally does **not** retune the Build34 decoder. It also normalizes the preceding checkpoint consistently as Build34 everywhere in code, tests and documentation.

- Renumbers the preceding projective/HMAC checkpoint consistently as Build34, including Go symbols, test filenames, Make targets, documentation and payload labels.
- Exposes `ExperimentalV4ExtractProjective` plus the explicit `v4-extract-projective` CLI. The caller supplies the block-aligned canonical pre-print dimensions; Build34 public structure/pilot evidence selects geometry before any frame/HMAC attempt.
- Adds `v4-build35-projective-api-test` to protect the exported projective API/CLI and the canonical-dimension contract.
- Replaces the old Build32 fixture helper with a scanner-first Build35 physical pack generated from the active MQ source: one unmarked control plus two robust/strength-24 carriers with different payloads (`v4-b35-phys-a`, `v4-b35-phys-b`).
- Adds a generated acquisition plan and `make v4-physical-qualification`; marked scans must recover the exact HMAC-authenticated payload while the control must reject.
- Defines the first physical gate as 300-ppi actual-size print followed by 300-dpi lossless scan cropped to the artwork edges. Smartphone acquisition is deliberately deferred until the controlled scanner gate separates print-channel failure from free-camera geometry failure.
- Stable Format v3, the locked pilot identity, Build31 frame/HMAC domains, Hamming baseline and Build29 acceptance floors remain unchanged.

## v0.3.0-build34 — 2026-09-15

Build34 closes the active MQ Format-v4 blind projective+crop blocker without changing the frozen v3 core, locked pilot identity, 37x32 tile partition, Build31 framing/HMAC domains, Hamming baseline or Build29 acceptance thresholds.

- Adds bounded multi-anchor local projective refinement for >=4x4-tile carriers: 2 distinct structural anchors, coupled 243-point local neighborhoods, structural top-64 retention and public-pilot ranking with 0/+-2/+-4 pixel centered phases.
- Uses all interior tile rows for proposal while reserving first/last tile rows as spatially held-out validation after the geometry decision.
- Canonicalizes non-zero cyclic pilot origin by composing the corresponding block translation into the canonical side of the homography, then requires the unchanged `(0,0)` acceptance gate.
- Adds projective sampling of the locked 1120 data positions after geometry acceptance and completes Hamming -> v4 frame -> CRC -> HMAC recovery.
- Adds an opt-in MQ authenticated corpus regression with two payloads (`v4-b34-auth`, `v4-b34-alt`) to guard against data-plane-specific geometry ranking.
- Keeps small/LQ carriers on the Build29 path; LQ projective remains a documented SAFE REJECT.

## v0.3.0-build33 — 2026-09-15

Build33 is a **ranking-observability checkpoint**, not yet the final joint-projective decoder fix. It preserves the frozen Format-v3 core, the locked `prototype-2-search-p64` pilot, the Build31 v4 frame and all Build29 acceptance floors.

- Added `v4-pilot-joint-projective-rank-diagnostic`, focused on the active MQ corpus entry and the exact Build29 projective+crop transform.
- The diagnostic measures the rank of the known true/near-true projective basin after the structural bank, half-pilot prefilter and full-proposal ranking instead of treating the final SAFE REJECT as an undifferentiated failure.
- Runs the same geometry twice: once with the historical synthetic random data-plane carrier and once with a real authenticated Format-v4 robust carrier (`v4-b33-auth`, strength 24, public development key). This exposes ranking sensitivity to the 1120 data positions without ever using payload/header/CRC/HMAC as a geometry oracle.
- Dumps the top 20 **distinct parameter-space basins** after full-proposal ranking so texture-induced basin crowding can be distinguished from raw-candidate multiplicity.
- Explicitly records that exploratory half-pilot/spatial beam variants were rejected because they could still promote convincing false geometry on one data plane. Those prototypes are not shipped.
- The Build33 completion criterion remains unchanged: MQ blind projective+crop must recover the correct mapping and only then authenticate the v4 payload; LQ may continue to SAFE REJECT.
- `all-test` grows from 50 to 51 targets. A PASS of the new diagnostic means the true basin remained observable through the measured stages; it does **not** mean the joint projective decoder is qualified.

## v0.3.0-build32 — 2026-09-14

Build32 resets the active private development corpus to three anonymous quality tiers (`Immagine LQ.png`, `Immagine MQ.png`, `Immagine HQ.png`) and makes corpus membership explicit and hash-checked. It does not change the frozen Format-v3 core, the locked v4 pilot, or the Build31 v4 frame/encoder.

- Added `private-corpus-active.tsv` with role, dimensions and SHA-256 for LQ/MQ/HQ.
- Added `corpus-manifest-check`; corpus-driven shell/Go tests no longer discover arbitrary files from the directory.
- Go v4 corpus tests skip entries above the 50 MP research budget before bitmap decode; HQ remains part of corpus identity.
- Changed the active v4 public development key to `PixSeal-v4-TestKey-2026`; historical deterministic vectors using `Piccotti` remain unchanged for reproducibility.
- Changed the common deep-test resize floor from 50% to 55% after LQ/capacity demonstrated an authenticated boundary at 50%; `extreme-test` still maps 50% and lower limits.
- Updated physical-fixture defaults to Build32, message `v4-b32-physical`, robust/strength 24 and the new development key.
- Physical fixtures remain private/non-release artifacts.

## v0.3.0-build31 — 2026-09-14

- Add the first real **experimental Format-v4 authenticated frame and encoder** while leaving stable Format-v3 `EmbedWithInfo` / `ExtractWithInfo` and CLI behavior unchanged.
- Add explicit experimental `v4-embed` and native-lattice `v4-extract` CLI commands; `v4-extract` intentionally does not claim Build29 geometry recovery yet.
- Keep the Build30 locked `prototype-2-search-p64` pilot unchanged at SHA-256 `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`.
- Define Build31 v4 header/profile bytes `0x41/0x42/0x43`, v4-specific whitening label `pixseal-whiten-v4`, and domain-separated frame HMAC label `pixseal-frame-v4`.
- Map the 1120 ascending non-pilot data ordinals with `(ordinal * 251) mod codedBits`; add an explicit partition regression so no locked prototype-2 pilot block can be treated as data.
- Retain the 32/48/80-byte profile frames and Hamming(7,4) protected lengths only as an apples-to-apples experimental ECC baseline; normative v4 ECC selection remains open.
- Add a deterministic Build31 frame/protected-bit vector, v3/v4 cross-version rejection, wrong-key rejection, one-error-per-Hamming-word correction, minimum 296x256 capacity round-trip, JPEG-q82 and aligned-crop tests.
- Add private-corpus frame qualification on both development originals: robust/balanced/capacity native round-trips authenticate, and robust additionally authenticates after JPEG-q82 and block-aligned crop.
- Add `docs/V4_BUILD31_FRAME.md`, `v4-frame-test` and `v4-frame-corpus-test`; the default `all-test` matrix grows from 47 to 49 targets.
- Build31 now enables creation of a separate private physical v4 print-camera/scanner corpus. Normative pilot promotion still requires HMAC-authenticated physical v4 recovery.

## v0.3.0-build30 — 2026-09-14

- Add a semantic development **candidate lock** for the exact `prototype-2-search-p64` pilot identity; the locked SHA-256 remains `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`.
- Add `v4-pilot-lock-check` so any change to the locked mask/sign identity or its core structural invariants fails immediately.
- Add a known-mapping Build30 corpus audit over the Build29 projective-crop and affine-padded fixtures. Across both local originals, marked score is >= 0.928412, marked margin >= 0.484086, origin is `(0,0)`, and matching negative margin is <= 0.010553.
- Confirm that the Build29 SAFE-REJECT cases are decoder search/placement limitations rather than evidence that the public pilot pattern itself has failed.
- Reject an attempted all-64-sub-block-phase DCT padded-ranking variant: it removes the even-phase bias but remains dominated by natural photographic texture. No Build29 acceptance floor is relaxed.
- Add `docs/V4_PILOT_CANDIDATE_LOCK.md` documenting exactly what is locked and what remains experimental. Normative pilot promotion is explicitly deferred until physical v4 print-camera/scanner qualification exists.
- Add `v4-pilot-lock-corpus-test`; the default all-test matrix grows from 45 to 47 targets.
- Keep the frozen Format-v3 core unchanged; no v4 payload/framing/ECC/HMAC encoder path is introduced.

## v0.3.0-build29 — 2026-09-14

- Add bounded experimental joint **projective+crop** and **affine+padded-canvas** searches for Format v4.
- Add explicit ACCEPT/SAFE-REJECT development gates based on held-out placement validation, complete public-pilot margin and absolute origin; no payload/key/ECC/HMAC evidence participates.
- Qualify synthetic projective+crop and padded cases; qualify `PJ_lingua.PNG` projective+crop while intentionally SAFE REJECTING the smaller `PJ_piccolo.png` projective case and both photographic padded joint-search cases.
- Add a known-geometry padded control showing that `PJ_lingua.PNG` retains a strong public-pilot channel even when the joint padded geometry ranking rejects.
- Add `v4-pilot-joint-projective-test` and `v4-pilot-joint-projective-corpus-test`; the default all-test matrix grows from 43 to 45 targets.
- Keep the frozen Format-v3 core and `prototype-2-search-p64` unchanged.

## v0.3.0-build28 — 2026-09-13

- Add the first bounded **joint blind affine+crop** Format-v4 search: rotation, anisotropic scale and negative crop/translation are unknown in the same experiment.
- Keep geometry proposal independent of pilot signs, payload, key and HMAC by using absolute DCT carrier energy plus 8-pixel phase contrast.
- Add hierarchical coarse/mid/full geometry search and a compact coupled angle/scale refinement only on the final two structural basins; this fixes the `PJ_piccolo.png` scale-coupling failure without loosening acceptance floors.
- Select geometry before exposing pilot symbols; reuse the Build27 split-pilot placement proposal/held-out validation path only after geometry is fixed.
- Qualify two synthetic joint-affine crop cases and one representative joint case on each local development original, all with cyclic origin `(0,0)` and matched-negative separation.
- Add `v4-pilot-joint-affine-test` and `v4-pilot-joint-affine-corpus-test`; `all-test` now contains 43 targets.
- Document why Go was selected for a single portable core and record the future GUI/mobile-app direction; `core-target-check` continues to cover Linux, Windows, Android and iOS targets.
- Keep joint shear, joint projective/perspective and joint padded-canvas placement explicitly open; `prototype-2-search-p64` remains non-normative and no v4 payload codec is enabled.

## v0.3.0-build27 — 2026-09-13

- Add an isolated Format-v4 unknown-placement search while preserving the frozen v3 core and non-normative `prototype-2-search-p64`.
- Qualify arbitrary non-block-aligned crop and padded-canvas carrier placement with the geometric transform supplied independently.
- Split the public 64-symbol pilot into disjoint proposal/validation halves so placement ranking does not reuse the same pilot symbols for acceptance.
- Derive bounded translation ranges from the known full transformed extent versus observed extent; use a 4-pixel proposal grid followed by integer-pixel refinement.
- Treat whole-tile-equivalent translations as expected periodic equivalence rather than requiring an artificial top-2 placement gap.
- Add negative controls and a wrong-geometry control showing placement optimization cannot compensate for a materially incorrect transform.
- Add local-corpus placement qualification on both development originals; local images remain excluded from source/evidence archives.
- Add `v4-pilot-placement-test` and `v4-pilot-placement-corpus-test`; `all-test` now contains 41 targets.

## v0.3.0-build26 — 2026-09-13

- Add an isolated bounded blind Format-v4 geometry search without modifying the frozen v3 core or production APIs.
- Add a key/payload/pilot-symbol-independent coarse proposal observable based on DCT-sign self-consistency at homologous data-plane positions of repeated 37x32 tiles.
- Use a three-stage evidence split: data-repeat geometry proposal, central-row public-pilot ranking, then held-out corner-pilot validation.
- Cover rotation, anisotropic scale, X/Y shear, mild perspective and combined projective/affine cases with explicit wrong-origin/negative-control telemetry.
- Bound the current synthetic search to <=6400 evaluated repeat hypotheses per positive/negative image and preserve deterministic behavior.
- Add local-corpus blind tests for rotation+anisotropic-scale and combined perspective on both development originals.
- Keep `prototype-2-search-p64` non-normative; arbitrary crop/translation, physical v4 acquisition, framing/version marker, ECC and payload encoder remain open.

## v0.3.0-build25 — 2026-09-13

- Keep the five-file Format-v3 core byte-for-byte frozen; no production encoder/decoder or CLI behavior changes.
- Add a known-geometry projective pilot sampler that measures prototype-2 after a separately supplied canonical-to-observed homography. This deliberately tests pilot survivability without claiming blind geometry recovery.
- Add a 16-case deterministic geometry matrix covering +12.3/-17 degree rotation, 110x90/90x110 anisotropic scale, X/Y shear, two mild perspective shapes, rotation+anisotropic scale, rotation+shear, 75% rotated resize, perspective combined with JPEG/blur/noise, and perspective+crop with/without JPEG.
- Require correct cyclic origin and a predeclared 0.10 positive-vs-negative runner-up-margin separation floor for the Build25 qualification tests; this is a development gate, not a normative decoder acceptance threshold.
- Add the same known-geometry matrix to the opt-in local original-image corpus. All 32 image/case combinations pass. The weakest positive margin is 0.322069 on PJ_piccolo under the combined perspective+blur case; its negative-control margin is 0.063984, leaving 0.258085 separation.
- Add `v4-pilot-geometry-test` and `v4-pilot-geometry-corpus-test`; `all-test` now contains 37 targets.
- Prototype-2 remains non-normative. Build25 establishes geometric-channel survivability given a correct mapping; the next open problem is blind pilot-assisted geometry estimation before any pilot freeze or real v4 payload encoder.

## v0.3.0-build24 — 2026-09-13

- Keep the five-file Format-v3 core byte-for-byte frozen and continue all v4 work in isolated experimental files.
- Add a reproducible two-stage Build24 pilot search: 200,000 joint coordinate/sign candidates followed by
  200,000 balanced-sign refinements on the winning mask, with fixed seeds, explicit budgets and candidate hashes.
- Select non-normative `prototype-2-search-p64` (SHA-256
  `858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174`). Compared with prototype-1,
  maximum cyclic mask overlap improves from 8 to 7 and maximum absolute wrong-shift correlation from 5 to 4.
- Add partial-visibility qualification at 64/48/32/24/16 visible pilots, including exhaustive contiguous
  8x8-stratum crops and 256 deterministic random subsets per level. Prototype-2 worst contiguous margins are
  60/43/28/20/12 symbols; no random-subset false origin occurs in the qualification set.
- Add an experimental synthetic v4 carrier renderer and pilot-only detector. These helpers embed a public pilot
  plus deterministic pseudo-random data-plane signs solely for qualification; they do not encode a payload and
  are not reachable from production `EmbedWithInfo`, `ExtractWithInfo` or the CLI.
- Add deterministic image-channel regressions for JPEG, blur, noise, gamma, 75%/50% resize and aligned crop,
  plus an opt-in local original-image corpus target. Both supplied local originals select the correct pilot
  origin in every initial transform case; unmarked originals remain negative controls.
- Add `v4-pilot-search-test`, `v4-pilot-channel-test` and `v4-pilot-corpus-test`; `all-test` now contains 35 targets.
- Do not freeze the v4 pilot yet: arbitrary rotation, anisotropic/affine distortion, perspective, combined
  transformations and real printed/scanned v4 acquisitions remain required before encoder promotion.

## v0.3.0-build23 — 2026-09-12

- Formally close the Format-v3 absolute-cycle research branch after build22 physical qualification;
  v3 remains the frozen production/interoperability baseline and future work is regression/maintenance.
- Add `docs/V3_FINAL_STATUS.md` with the final v3 baseline, physical-channel evidence, closure rationale,
  non-goals and maintenance rules.
- Start an isolated experimental Format-v4 foundation without wiring it into `EmbedWithInfo`,
  `ExtractWithInfo`, CLI version selection or authentication.
- Add provisional 37×32 / 64-pilot / 1120-data prototype constants and a deterministic 64-position,
  sign-balanced pilot candidate selected from a stratified search family.
- Add exhaustive structural pilot invariants: 64 unique pilot positions, 1120 non-pilot data positions,
  32/32 sign balance, 16 pilots per quadrant, no perfect non-zero cyclic alias, maximum cyclic overlap
  8/64 and maximum wrong-shift signed correlation 5/64 for prototype-1.
- Add `v3-freeze-check` backed by `docs/V3_FROZEN_CORE_SHA256.txt`, plus `v4-foundation-test`;
  extend `all-test` to 32 targets.
- Update README, HISTORY, TODO, ALGORITHM, RESULTS, RESEARCH_LOG, PRIVATE_CORPUS and the v4 design study
  so v3 closure and v4 experimental status are unambiguous.

## v0.3.0-build22 — 2026-09-12

Twenty-second print-acquisition research checkpoint. Format v3 and every production decoding rule
remain frozen. This build measures build21's weak repetition-topology asymmetry in real image-domain
margin grids and begins a quantified, non-normative Format-v4 absolute-pilot design study.

### Added

- `v3-image-domain-heldout-repetition-topology`: shared repetition edges register two unit-shifted
  phase hypotheses; topology-exclusive edges are held out and used only for discrimination.
- Per-shift aggregate/cell deltas, directional cell consistency and winner-phase stability.
- `physical_topology_ms` / attempt counters and flattened summary telemetry.
- `format_v4_design_study` with three pilot/tile sizing candidates.
- `docs/FORMAT_V4_DESIGN.md`, including explicit performance/reliability goals and prototype gates.
- `make physical-topology-test` and `make v4-design-study-test`; `make all-test` now contains 30 targets.

### Research result

- Inclined smartphone: winner modal-phase mean 0.3125, directional cell consistency 0.6944.
- Scanner 001: 0.4375 / 0.7292.
- Scanner 002 negative control: 0.4375 / 0.6597.
- Frontal smartphone: exact unwrap remains not-applicable; physical topology probe is skipped.
- The held-out topology signal is measurable but does not separate useful acquisition from the
  negative control, so it is not promoted as an absolute-cycle gate.
- The preferred v4 prototype is 37×32 with 64 public pilot blocks: +5.7% tile area, 1120 data
  positions retained, current v3 capacity ceilings preserved under the comparison assumptions.
- No v4 production path, HMAC budget or Format-v3 rule changes.

## v0.3.0-build21 — 2026-09-12

Twenty-first print-acquisition research checkpoint. Build20 decoding and cycle-anchor behavior
remain frozen; this build audits whether public Format-v3 structure can theoretically break a
one-block integer-cycle ambiguity without secret or oracle information.

### Added

- Static `v3-key-independent-structural-observability` report in `diagnose -json`.
- Exhaustive unit-shift (`dx,dy in {-1,0,+1}` except zero) repetition-pair topology overlap for
  all three v3 profiles.
- Deterministic valid-Hamming-frame probe through the real tile mapping/hard aggregation path,
  plus exact detection of pure code-index permutations that preserve whole Hamming words.
- `make observability-audit-test`; `make all-test` now contains 28 targets.

### Research result

- Robust: 896 repetition pairs; nearest unit-shift pair-topology gap 0.0580. Vertical ±1 has
  zero syndrome in the deterministic Hamming probe despite a non-invariant repetition graph.
- Balanced: 448 repetition pairs; nearest gap 0.0558. Vertical ±1 remains only weakly separated
  by Hamming syndrome (~0.088--0.093), while horizontal/diagonal shifts are near random-code
  syndrome occupancy (~0.77--0.79).
- Capacity: no repetition pairs. Vertical ±1 is an exact code-index permutation that preserves
  Hamming word boundaries and bit roles, so Hamming syndrome is structurally invariant there.
- The audited mechanisms therefore do not separate every unit-cycle shift for every profile.
  No decoder, HMAC budget, encoder or Format-v3 rule changes.

## v0.3.0-build20 — 2026-09-12

Twentieth print-acquisition research checkpoint. Build19 repetition-only conclusions remain
frozen; this build tests an independent image-domain cycle anchor without changing Format v3
or production decode.

### Added

- Diagnostic-only `cross-cell-pairwise-cycle-anchor` based on unguided cross-cell margin-grid
  registration, independent of repetition partitions, key/header/HMAC or payload semantics.
- Gauge-invariant top-1 versus exact-runner-up objective, per-cell anchor telemetry and
  explicit runtime/attempt counters.
- `make cycle-anchor-test`; `make all-test` now contains 27 targets.

### Research result

- Inclined smartphone: anchor available, 12 pair constraints, mean confidence 0.0575; top-1
  objective 5.9922 vs runner-up 6.5005 (delta +0.5083), but 0/9 rounded-cycle agreement.
- Scanner 001: mean confidence 0.0168; top-1 8.7807 vs runner-up 9.0084 (delta +0.2277), yet
  0/9 agreement. This conflicts with the held-out repetition result that preferred runner-up.
- Scanner 002 negative control: mean confidence 0.0798; top-1 14.6072 vs runner-up 14.7028
  (delta +0.0956), also 0/9 agreement.
- Frontal smartphone remains `2/24`, unwrap not-applicable and anchor not run.
- No physical HMAC authenticates. The pairwise observer remains diagnostic-only and is not a
  sufficient absolute integer-cycle anchor.

## v0.3.0-build19 — 2026-09-12

Nineteenth print-acquisition research checkpoint. Build18 decisions and lazy scheduling
remain frozen; this checkpoint tests integer-cycle stability across multiple disjoint
repetition partitions.

### Added

- Eight deterministic coded-bit-group partitions, each evaluated in both held-out
  directions for at most 16 exact integer-cycle trials on the final best ambiguous bit
  candidate.
- Full-field stability telemetry: available/supported trials, complete-field count,
  unique fields, modal-field frequency, supported-only field statistics and mean pairwise
  cycle agreement.
- Per-cell modal cycle, unique-cycle count and modal fraction across all trials and across
  held-out-supported trials only.
- `partition_stability_ms` / attempt counters, explicit bounded state budgets and
  `make stability-unwrap-test`; `make all-test` now contains 26 targets.

### Research result

- `foto stampa.jpg` remains `2/24` and `not-applicable`; stability is not run.
- `foto stampa storta.jpg`: 16/16 trials available, 9 supported, **16 distinct complete
  fields**; supported trials also produce 9/9 distinct fields. Mean cell modal fraction
  is 0.2153 (0.2716 supported-only), minimum 0.125, mean pairwise agreement 0.0731.
- Scanner `0270_001.jpg`: 16/16 available, 6 supported, 16 distinct fields; all six
  supported fields distinct. Mean cell modal fraction 0.2153, pairwise agreement 0.0759.
- Scanner `0270_002.jpg`: 16/16 available, 12 supported, 16 distinct fields; all twelve
  supported fields distinct. Mean cell modal fraction 0.2361, pairwise agreement 0.0824.
- No physical acquisition authenticates. Multi-partition evidence therefore rejects the
  hypothesis that build17 instability was caused by one unlucky A/B split. Repetition
  evidence alone is not a stable cycle anchor on the current corpus.

## v0.3.0-build18 — 2026-09-12

Eighteenth print-acquisition research checkpoint. Build17 decisions are frozen; this
checkpoint localizes cross-fold instability and makes held-out work lazy.

### Added

- Per-cell integer-cycle telemetry for the all-pairs exact proposal and both held-out
  cross-fit directions, including fold confidence and A/B agreement.
- Aggregate agreeing-vs-disagreeing mean joint cross-fit confidence and lattice
  confidence.
- Explicit `crossfit_unwrap_ms`, attempt and skipped-candidate counters.
- Diagnostic budget flag documenting that cross-fit is best-candidate-only.

### Changed

- A/B cross-fit is no longer repeated for every bit diagnostic. All candidates retain
  all-pairs exact unwrap, but held-out cross-fit runs at most once on the final best bit
  candidate, and still only for an exact ambiguous top-2.
- `crossfit-unwrap-test` now also covers public per-cell telemetry.
- Private acquisition summaries expose the new runtime and confidence aggregates.

### Research result

- Physical decisions are unchanged from build17: frontal `2/24` not-applicable;
  inclined `7/24` ambiguous with 2/9 A/B agreement; scanner 001 `4/24` with 0/9;
  scanner 002 `5/24` with 0/9.
- The two agreeing inclined-photo cells are not the strongest cross-fit controls: mean
  joint confidence is `0.1337` for agreeing cells versus `0.3576` for disagreements.
  A naive partial-consensus promotion is therefore rejected.
- On the local development host the selected cross-fit costs only `26--73 ms` and is
  executed once, with three lower-ranked candidates skipped in each ambiguous case.

## v0.3.0-build17 — 2026-09-11

Seventeenth print-acquisition research checkpoint. Format v3, deterministic encoding,
production extraction and all HMAC/smooth-resample ceilings remain frozen.

### Added

- True held-out repetition cross-fit for ambiguous integer-cycle fields. Repetition
  constraints are partitioned by coded-bit group so proposal and validation folds do
  not reuse a repeated-position group or logical repetition position.
- Symmetric `A -> B` and `B -> A` experiments: one fold estimates blind repetition
  controls and exact top-1/top-2; only the opposite fold scores that distinction.
- Cross-fit telemetry for fold profiles/pair budgets, exact states, directional
  margins/deltas, directional support, local-cycle proposal agreement and aggregate
  support.
- `make crossfit-unwrap-test`; `make all-test` now contains 25 targets.
- Explicit cross-fit budgets in diagnostic JSON: two folds, coded-bit-group partition,
  ambiguous-only execution and a declared worst-case exact-state ceiling.

### Changed

- Cross-fit runs only when the normal all-pairs exact solver is already ambiguous with
  an available second solution. This keeps the new research work away from accepted,
  no-change and not-applicable candidates and preserves bounded runtime.
- Proposal agreement compares recentered local integer cycles rather than fold-specific
  absolute tile phase, because equivalent global tile anchors are not the cycle field
  under study.
- Private print-camera/scan summaries now expose the build17 cross-fit fields.

### Research result

- `foto stampa.jpg`: unchanged best hard `2/24`, global unwrap `not-applicable`; no
  cross-fit is run and no HMAC authenticates.
- `foto stampa storta.jpg`: both held-out directions favor their own exact top-1
  (`+0.00465`, `+0.20522`), but the independently proposed local cycle fields agree in
  only `2/9` cells. The all-pairs unwrap therefore remains ambiguous and rolled back.
- Scanner `0270_001.jpg`: both held-out directions prefer the exact runner-up
  (`-0.18283`, `-0.02287`) and the proposed local cycle fields agree in `0/9` cells.
- Scanner `0270_002.jpg`: the two directions disagree (`+0.05639`, `-0.24050`) and
  proposal agreement is `0/9`. The historical bad-unwrap counterexample remains safely
  rejected.
- No available real acquisition authenticates. Build17 therefore demonstrates that
  genuinely disjoint repetition evidence is informative but still insufficient to
  identify one stable integer-cycle field across the current physical corpus.

## v0.3.0-build16 — 2026-09-11

Sixteenth print-acquisition research checkpoint. Format v3, deterministic encoding,
production extraction and the existing HMAC/smooth-resample ceilings remain frozen.

### Fixed

- Replaced build15's width-64 beam search with exact bounded per-axis enumeration,
  yielding a true top-1/top-2 over at most `3^9 = 19683` assignments per axis.
- Added a deterministic regression for a build15 beam false accept: the beam missed a
  closer runner-up and overstated uniqueness; build16 correctly reports ambiguity.
- Removed non-finite `+Inf` second-objective reporting from no-eligible/degenerate
  unwrap paths and added a JSON-serialization regression.
- Made mixed-axis ambiguity atomic: if either X or Y is ambiguous, no otherwise
  acceptable axis is partially committed and the global status remains `ambiguous`.
- Clarified the physical-corpus credential model: acquisition files remain private, while the intentional test key `Piccotti` remains a public/reproducible Makefile/script default.
- Added `/print-scan private/` and local private-manifest/config artifacts to `.gitignore`.
- Updated the two historical smartphone filenames to canonical `.jpg` names matching
  their actual JPEG content.

### Added

- Explicit X/Y unwrap statuses and eligibility counts plus second-solution availability.
- Deterministic `split-repetition-top2` diagnostics comparing exact top-1/top-2 cycle
  assignments on two disjoint key-independent repetition-pair folds. This signal is
  diagnostic-only and cannot override exact ambiguity in build16.
- `docs/PRIVATE_CORPUS.md` and `make private-corpus-manifest` for local corpus integrity.

### Research result

- `foto stampa storta.jpg` remains conservatively ambiguous with exact proposed/second
  objectives about 0.80446 / 0.82560; no cycle is committed.
- Scanner `0270_002.jpg` tightens from the build15 beam margin ~0.02634 to an exact
  margin ~0.01780 while remaining ambiguous. This directly validates the audit finding
  without changing the safe decision.
- No real physical acquisition authenticated in the four-case corpus available during
  the build16 session. The two bicycle cases were unavailable and remain required for
  full 6/6 validation.


## v0.3.0-build15 — 2026-09-11

Fifteenth print-acquisition research checkpoint. Format v3, deterministic encoding,
production extraction and the maximum four complete/HMAC decode slots remain frozen.

### Added

- Key-independent global discrete integer phase-unwrapping solver for the 3x3 blind
  control field. X/Y are solved independently with a deterministic bounded beam search.
- Candidate cycle shifts are limited to `-1, 0, +1` block. The objective combines
  robust affine residual, confidence-weighted cycle-change cost and second-difference
  curvature; no key, known header bit, CRC or HMAC result participates.
- Explicit ambiguity reporting: baseline/proposed/applied/second-best objectives,
  absolute/relative improvement, margin, evaluated states, eligible/changed cells and
  accepted axes.
- Transactional lattice fusion: fractional lattice phase and integer unwrap are applied
  together only when the global unwrap is accepted. Ambiguous/rejected proposals roll
  back completely to the pre-lattice blind controls.
- `make global-unwrap-test`, including known-slip recovery, high-confidence locking,
  ambiguous-layout rejection, rollback, budget checks and end-to-end synthetic HMAC.

### Changed

- `make all-test` now contains 24 targets.
- `make test-list` documents the new global unwrap target separately from fractional
  lattice-phase measurement.

### Research result

- The build14 per-cell unwrap can yield large same-candidate gains, but build15 shows
  that those cycle choices are often not uniquely identifiable by a key-independent
  global objective.
- On `foto stampa storta.jpg`, a representative final global proposal changes all nine
  cells and greatly lowers the geometric objective, but the first/second solution margin
  is only about 0.021; the proposal is therefore marked ambiguous and rolled back. The
  safe blind path returns to the build12/13-like 5 -> 4 diagnostic result, with no HMAC.
- Representative `foto bici dritta.jpg` and scanner `0270_002.jpg` proposals are also
  rejected as ambiguous. This is evidence against promoting build14's local unwrap
  merely because it improves the known prefix.
- No production decoder path or HMAC budget is changed.

## v0.3.0-build14 — 2026-09-11

Fourteenth print-acquisition research checkpoint. Format v3, deterministic encoding,
production extraction and the maximum four complete/HMAC decode slots remain frozen.

### Added

- Key-independent local fractional lattice-phase observer derived from the already
  selected 3x3 `LocalLatticeEstimate` records and current candidate geometry.
- Approximate inverse mapping through projective/canonical-offset/residual-warp state,
  followed by modulo-8 canonical phase measurement.
- Confidence-gated fractional fusion with the repetition observer.
- Bounded smooth integer unwrap: at most +/-1 block per axis, two passes, only for
  low-confidence primary controls and only when an affine robust prediction provides
  a large residual improvement.
- JSON evidence for lattice-phase availability, confidence, modulo-phase distance,
  consensus cells and cycle-slip corrections.
- `make lattice-phase-test`.
- `docs/RESEARCH_LOG.md`, an append-only record of hypotheses, rejected variants and
  negative results.

### Changed

- `make test-list` now reports each target's input dependency (`source`, `pics`,
  `private`, or aggregate) and warns that corpus-sensitive research/qualification
  failures are not automatically production regressions.
- `make all-test` now contains 23 targets.

### Research result

- The build13 guided cross-cell observer remains diagnostic; its weak real score is
  not used to drive the lattice-phase fusion.
- The local lattice observer is materially stronger on the current real corpus, with
  mean lattice confidence about 0.38--0.52 on the reported best-hard spatial cases.
- An initial inverse-correction sign interpretation and a fractional-only/no-unwrap
  variant were explicitly tested and rejected. The final convention reports observed
  lattice drift; the smooth fitter applies the inverse correction.
- With bounded unwrap, same-candidate diagnostics include 5 -> 2 on
  `foto stampa storta.jpg`, 7 -> 3 on `foto bici dritta.jpg`, and 9 -> 4 on scanner
  `0270_001.jpg`. Other candidates worsen, including scanner `0270_002` 7 -> 10.
- No real field passes all decode-promotion gates; smooth HMAC attempts remain zero and
  no physical acquisition authenticates.

## v0.3.0-build13 — 2026-09-11

Thirteenth print-acquisition research checkpoint. Build13 tests whether an independent
guided cross-cell image-domain observation can disambiguate integer cycle slips in the
build12 repetition-coherence controls without using key/header evidence or increasing
HMAC slots.

### Added

- Guided cross-cell correlation around the repetition-predicted relative shift: +/-1
  block, at most 324 pair probes for nine cells.
- Explicit secondary-observer, observer-distance, consensus-cell and cycle-slip
  reporting.
- Independent secondary-strength gate before any fusion or +/-1 cycle-slip correction.
- Synthetic regressions proving that strong guided evidence participates in consensus
  and that cycle-slip correction remains bounded and cannot override a high-confidence
  primary control.
- `make test-list`, with categorized target descriptions and RELEASE / QUALIFICATION /
  RESEARCH / PRIVATE / PORTABILITY / AGGREGATE labels.

### Changed

- Full +/-4 cross-cell correlation is now fallback-only when repetition evidence is
  unavailable; normal dual-observer cost is capped at 324 pair probes.
- Weak secondary evidence is diagnostic-only and cannot veto or modify the primary
  repetition controls.
- The complete `all-test` qualification matrix remains 22 targets; `test-list` is
  informational.

### Real-corpus result

- No smartphone or scanner acquisition authenticates.
- Guided cross-cell mean pair peaks are only about 0.053--0.065 on the six reported
  best-hard cases, versus about 0.61 on the controlled synthetic consensus regression.
- No real case clears the 0.10 consensus-strength gate; zero cycle-slip fixes are
  applied and no secondary observation changes an HMAC candidate.
- `foto stampa storta.jpg` therefore preserves the build12 blind 5 -> 4 measured
  same-candidate improvement and still fails the strict LOO decode gate.
- `foto stampa.jpg` remains the best hard real case at 2/24.

## v0.3.0-build12 — 2026-09-11

Twelfth print-acquisition research checkpoint. Build12 replaces key-assisted local
phase controls in the smooth-field path with a key-independent observer derived from
Format-v3 structural repetition. Format v3, the encoder, production extraction and
the four complete/HMAC decode slots remain frozen.

### Added

- Key-independent intra-tile repetition-coherence phase scoring for robust and
  balanced profiles; capacity intentionally supplies no repetition evidence.
- Blind aggregate phase search followed by bounded +/-2-block per-cell refinement
  with fractional peaks and confidence inherited from the build11 surface logic.
- Pairwise cross-cell DCT-pattern correlation as an explicitly bounded fallback.
- Ex-post blind-vs-key-assisted-oracle distance reporting; the oracle never feeds the
  blind fit.
- Explicit blind search budgets: 2240 global structural probes, 225 local probes and
  2916 pairwise-fallback probes, all on already-collected grids.
- `make blind-phase-test`.

### Changed

- The real smooth-field path now fits only `blind-self-registration` controls.
- The build11 key-assisted local phase remains diagnostic/oracle evidence only.
- `make all-test` now contains 22 targets.

### Real-corpus result

- No current smartphone or scanner acquisition authenticates.
- Direct cross-cell blind correlation is weak on the real corpus; intra-tile
  structural repetition is materially stronger and is used as the primary observer.
- Only `foto stampa storta.jpg` produces blind fields that pass the diagnostic
  resample gate. The best measured same-candidate correction changes 5 -> 4 known
  post-ECC header errors (fit RMS about 1.27 blocks, LOO about 2.04 blocks).
- The strict decode gate rejects that field because LOO remains above 2.0 blocks, so
  no blind/smooth field consumes an HMAC slot.
- `foto stampa.jpg` remains the best hard real case at 2/24 and is left untouched.
- Wrong-key real control remains unauthenticated; blind structure scores are not
  watermark detection.

## v0.3.0-build11 — 2026-09-10

Eleventh print-acquisition research checkpoint. Build11 improves spatial phase-control
quality before increasing model freedom. Format v3, the encoder, production extraction
and the four complete/HMAC decode slots remain frozen.

### Added

- Fractional-block local phase estimates from a bounded clipped correlation surface
  with parabolic peak interpolation.
- Per-cell phase confidence from correlation strength, peak prominence and curvature.
- Explicit detection and confidence penalty for local phase maxima at the +/-2-block
  search boundary.
- Confidence-weighted Huber-robust smooth-field fitting.
- Ridge-regularized quadratic comparison (six basis terms per axis) with strict
  leave-one-out model selection; affine remains the default.
- JSON reporting for phase confidence, robust outliers, affine/quadratic fit metrics
  and the selected smooth model.
- `make phase-surface-test`.

### Changed

- The build10 smooth decode gate keeps its RMS/LOO limits but now also requires a
  minimum mean control confidence.
- Boundary-truncated controls are down-weighted instead of being treated as precise
  integer observations.
- Quadratic fields are research-only unless cross-validation improves by both an
  absolute and relative margin and all existing correction/budget gates pass.

### Real-corpus result

- No current smartphone or scanner acquisition authenticates.
- No real case consumes a smooth-field HMAC slot in the final build11 code.
- The build10 `foto bici dritta` promoted field becomes 6->8 after confidence-aware
  control estimation and is therefore rejected, indicating that its earlier 6->5
  improvement was not robust to better control modelling.
- `foto bici storta` retains the strongest same-candidate diagnostic improvement,
  4->3, but leave-one-out RMS is about 2.18 blocks and fails the decode gate.
- One `0270_002` candidate selects the regularized quadratic model by cross-validation,
  but its maximum predicted correction exceeds the bounded cap, so it is rejected.

## v0.3.0-build10 — 2026-09-10

Tenth print-acquisition research checkpoint. Build10 tests whether the spatial
phase drift measured in build9 can be represented by one bounded, smooth affine
field without changing Format v3 or increasing the four HMAC-decode slots.

### Added

- Six-parameter affine phase-field fit over the existing 3x3 spatial cells.
- Separate diagnostic and strict decode gates using fit RMS and leave-one-out RMS.
- At most two bounded smooth-field resamples on already-ranked decode candidates.
- Same-candidate before/after known-prefix metrics and explicit smooth HMAC-slot count.
- `make smooth-phase-test` with exact affine recovery, HMAC authentication,
  non-smooth rejection and budget checks.

### Changed

- Spatial diagnostic cells now ignore incomplete edge tiles while the ordinary
  aggregate decode grid still uses those edge blocks.
- A smooth corrected grid can replace an existing full-decode grid only when the
  strict fit gate passes and known-prefix hard-ECC errors improve without adding
  known coded-bit errors. No new HMAC slot is created.

### Real-corpus result

- Same-candidate smooth-field measurements improve four of six current private
  acquisitions and worsen two; no acquisition authenticates.
- The strongest diagnostic change is `foto bici dritta` 7->3 known header errors
  after ECC, but that field does not pass the strict leave-one-out gate.
- One different `foto bici dritta` candidate passes the strict gate, changes 6->5
  and consumes one existing HMAC slot; HMAC still fails.
- The result supports a spatially varying residual component but rejects a single
  global affine field as a general solution.

## v0.3.0-build9 — 2026-09-10

Ninth print-acquisition research checkpoint. Build9 keeps Format v3, the encoder,
production extraction, geometry budgets and the four complete HMAC-decode slots
unchanged. It measures whether protected carrier errors are spatially local or
systematic across repeated tiles.

### Added

- Zero-extra-read 3x3 spatial buckets collected during the existing full-grid
  sampling pass.
- Key-independent sign-agreement metrics over all 1120 v3 tile positions.
- Known-prefix classification of the 42 protected header bits as stable-correct,
  stable-wrong or mixed across spatial cells, plus diagnostic majority results.
- A bounded +/-2-block local phase oracle (max 25 phase positions per cell, max
  225 across 9 cells) that never changes the decode grid or consumes HMAC slots.
- Candidate-paired hard/soft reporting so hard-to-soft changes are never compared
  across different geometry/photometric candidates.
- `make spatial-channel-test`.

### Six-acquisition findings

- Mean key-independent tile-position sign agreement is only about 0.635--0.683
  across the four smartphone photographs and two 600-dpi scanner JPEGs.
- For the best hard bit-channel candidate, 38--42 of the 42 known protected header
  bits are mixed across spatial cells in every case. Stable-wrong bits are zero in
  five cases and one in the remaining case.
- Scanner lattice consistency remains about 0.87, yet spatial sign agreement is
  only about 0.655--0.659, proving that a strong macro-lattice estimate does not
  imply spatially stable protected-bit observations.
- The bounded local phase check usually prefers offsets of roughly 1--1.8 blocks,
  but its known-header majority improves some cases and degrades others. It is
  retained only as evidence for a future smooth low-DOF phase/warp model.
- No real acquisition authenticates; HMAC remains the sole success criterion.

### Preserved

- Maximum four complete diagnostic Format-v3/HMAC decode slots.
- Build8 scanner/full-grid/reliability budgets.
- Frozen Format v3, deterministic encoder fingerprints and production extractor.

## v0.3.0-build8 — 2026-09-10

Eighth print-acquisition research checkpoint. Build8 keeps Format v3, the encoder
and production extraction unchanged, makes diagnostic authentication truly bounded
on scanner-sized inputs, and evaluates a deterministic reliability-aware Hamming
decoder without increasing the four full-decode/HMAC-slot ceiling.

### Added

- A 16 Mi-pixel diagnostic-only budget for the mature baseline extractor. Larger
  scanner/camera inputs skip that pre-pass and enter the bounded projective path.
- Full-grid work bounding: candidates above 50,000 canonical blocks use at most 25
  spatially distributed complete v3 tiles.
- Authentication sub-stage timings and explicit JSON counters for baseline skips,
  bounded grid use, sampled blocks/tiles and reliability-decode attempts.
- Deterministic soft Hamming(7,4) maximum-likelihood decoding from signed DCT
  margins, gated only by the fixed known v3 prefix and consuming an existing decode
  slot rather than adding a candidate.
- `make reliability-test` and optional `make print-scan-test`.

### Scanner/real-corpus findings

- The two 600-dpi 4960x7015 scanner JPEGs complete full diagnostic authentication
  in about 10--11 s after previously exceeding 120 s. Their lattice consistency is
  about 0.869/0.875, but their best known protected headers still contain roughly
  8--9 coded errors and 2--3 multi-error Hamming words.
- The original frontal smartphone capture gets one reliability decode, but soft
  Hamming does not improve its two remaining known-header errors and does not pass
  HMAC. Other real captures do not satisfy the conservative reliability gate.
- No smartphone or scanner acquisition authenticates. High lattice consistency and
  header z-score remain research evidence only.

### Preserved

- Format v3 and deterministic encoder fingerprints.
- Stable `ExtractWithInfo` implementation/search behavior.
- Maximum four complete diagnostic decode/HMAC slots and HMAC-only success.

## v0.3.0-build7 — 2026-09-10

Seventh print-camera research checkpoint. Build6 geometry/fallback behavior,
Build5 photometric modes and all Format-v3/production extraction rules remain
unchanged. Build7 adds bounded protected-bit/ECC diagnostics and makes the
private print-camera harness corpus-driven.

### Added

- `DiagnosticBitChannelEvidence` on the same candidates already selected for at
  most four complete virtual v3 decodes. It reports the strongest profile/phase,
  protected coded-bit margins, 42 known coded bits derived from the fixed first
  three v3 header bytes, the six exactly-known Hamming(7,4) words, post-ECC
  known-header bit errors, Hamming syndrome occupancy and correct-vs-wrong known
  bit margin ratios.
- `make bit-channel-test` with exact synthetic and deliberate two-errors-in-one-
  codeword regressions.
- Corpus-driven `make print-camera-test`: every `.png`, `.jpg` or `.jpeg` file
  in the private directory is tested in deterministic order and summarized.

### Real-corpus research status

- Original frontal capture: best known-header channel has 7/42 protected-bit
  errors, 1/6 known Hamming words beyond the one-bit correction radius and 2/24
  known-header errors after ECC; wrong-bit margin ratio ≈0.67.
- Original inclined capture: 11/42, 3/6 and 6/24 respectively; ratio ≈1.09.
- Held-out frontal bicycle capture: 12/42, 3/6 and 4/24; ratio ≈0.88.
- Held-out inclined bicycle capture: 8/42, 2/6 and 4/24; ratio ≈1.57.
- No real capture produces a valid Format-v3 HMAC. These measurements do not
  authenticate or reveal the hidden payload.

### Preserved

- Four-full-decode ceiling; bit diagnostics reuse the already materialized grids
  and do not multiply geometry/photometric search.
- Format v3, deterministic encoder fingerprints, stable production extractor and
  HMAC-only success semantics.

## v0.3.0-build6 — 2026-09-10

Sixth print-camera research checkpoint. Build5 photometric views and all frozen
Format-v3/production-decoder behavior are preserved. The diagnostic geometry no
longer requires a strong detected print boundary before a lattice-supported
projective attempt can be made.

### Added

- One bounded adaptive pyramid escalation for large, no-boundary, borderline
  lattice cases. The default 1/8 + 1/16 analysis may add only one finer level,
  normally 1/4, capped at a 4096-pixel diagnostic dimension.
- Lattice-first projective fallback: after independent lattice evidence exists,
  a sane weak print quadrilateral may seed the existing bounded projective
  estimator even when `print_boundary.detected=false`.
- Explicit `adaptive_escalated`, `adaptive_divisor`, `adaptive_reason` and
  `lattice_first_fallback_used` diagnostics.
- Regression tests proving the adaptive trigger remains bounded and that weak
  quadrilaterals require minimum geometric/confidence sanity before homography
  use.

### Held-out research status

- `foto bici dritta.jpg`: build5 stopped at consistency≈0.533 with no projective
  decode. Build6 adds 1/4, reaches consistency≈0.615, promotes lattice evidence,
  uses the weak boundary only after that evidence, and performs four projective
  decodes. Best bounded sync is balanced z≈4.454; no HMAC authenticates.
- `foto bici storta.jpg`: strong-boundary path is unchanged (consistency≈0.781,
  best photometric z≈4.079); no HMAC authenticates.
- The two original private photographs retain their build5 best results
  (frontal mild-highpass z≈4.695; inclined raw z≈3.973).
- A 200 MP unmarked qualification image can produce diagnostic lattice evidence
  after adaptive escalation, but its boundary confidence is below the weak-seed
  gate and no projective fallback is opened. This reinforces HMAC-only success
  semantics.

### Preserved

- Maximum four complete projective/HMAC decode attempts.
- Build5 three-view photometric bank and build4 geometry/refinement budgets.
- Format v3, encoder fingerprints, production `ExtractWithInfo` and HMAC-only
  authentication semantics.

## v0.3.0-build5 — 2026-09-10

Fifth print-camera research checkpoint. Build4 geometry is preserved while a
bounded photometric experiment is added after geometry selection. Format v3,
the deterministic encoder and production `ExtractWithInfo` remain unchanged.

### Added

- Fixed three-view photometric bank: `raw`, `local-normalize` and
  `mild-highpass`.
- At most three previously selected geometries are evaluated, for a maximum of
  nine photometric header probes; the complete HMAC decode budget remains four.
- Per-mode known-header fraction/z-score, mean absolute DCT margin and normalized
  signed sync margin in JSON diagnostics and `print-camera-test`.
- `make photometric-test` with pristine-carrier authentication and explicit bank
  budget regressions.
- Release/qualification/research classification in `make all-test` and uniform
  large-image SKIP behavior in bounded geometry scripts.

### Real-corpus research status

- Frontal: `mild-highpass` on the retained fundamental phase-DLT geometry raises
  the known-header score from z≈4.454 / fraction≈0.768 to z≈4.695 /
  fraction≈0.783.
- Inclined: the photometric bank does not beat the build4 raw global maximum
  z≈3.973; this negative result is retained rather than tuned away.
- Neither photograph produces a valid Format-v3 HMAC. The hidden payload remains
  unknown.

### Preserved

- Build4 geometry candidate generation/refinement and four-full-decode ceiling.
- Format v3, encoder fingerprints, production extraction order and HMAC-only
  success semantics.

## v0.3.0-build4 — 2026-09-10

Fourth print-camera research checkpoint. Format v3, the deterministic encoder and
the production `ExtractWithInfo` path remain unchanged from v0.2.0.

### Added

- Explicit fundamental-scale selection after projective scale clustering. A
  candidate must have multi-level and multi-region support; among sufficiently
  supported candidates the lowest spatial frequency is preferred, preventing a
  high-frequency phase alias from automatically becoming the carrier scale.
- A bounded key-assisted refinement around only the selected fundamental family:
  width and height are probed independently at -2%, -1%, -0.5%, 0, +0.5%, +1%
  and +2% (maximum 16 sync probes including full-resolution verification).
- Modulo-8 DCT block-origin refinement. At most three candidates search one
  periodic 8-pixel cell with coarse 2-pixel and fine 0.5-pixel offsets; sparse
  winners are rechecked with the ordinary probe and cannot silently degrade the
  starting candidate.
- One bounded residual-warp diagnostic derived from the 3x3 key-independent
  local lattice field. Control residuals are reduced modulo the 8-pixel block
  grid before fitting a smooth quadratic correction.
- Explicit fundamental/subpixel/residual counters, budgets and evidence in
  `diagnose -json` and `make print-camera-test`.
- Regression coverage for harmonic rejection, modulo-8 subpixel refinement and
  smooth residual-warp fitting.

### Real-corpus research status

- The selector independently promotes the matching low-frequency family near
  1668x1250 on `foto stampa.jpg` and 1653x1254 on `foto stampa storta.jpg`; no
  private-corpus dimensions are hard-coded.
- On the frontal photograph the fundamental branch remains near z=3.50 and does
  not beat the build3 global best z=4.454. The best subpixel result is the
  unshifted origin, so build4 does not claim a synthetic improvement there.
- On the inclined photograph the fundamental-only bounded scale refinement moves
  to about 1686x1279 and reaches z=3.732, while the global best remains the
  build3 value near z=3.973.
- The lattice-derived residual fits use nine controls (RMS about 2.28 canonical
  pixels frontal and 1.24 inclined) but do not yet produce a valid HMAC.
- Neither photograph authenticates. The real hidden payload therefore remains
  unknown and `print-camera-test` correctly remains a research failure.

### Preserved

- Format v3, encoder fingerprints and production extraction search order.
- Eight coarse projective candidates, three bounded refinement seeds and at most
  four complete virtual Format-v3/HMAC decodes.
- HMAC-only success semantics and the private-corpus exclusion policy.

## v0.3.0-build3 — 2026-09-10

Third print-camera research checkpoint. Format v3, the deterministic encoder and
the production `ExtractWithInfo` path remain unchanged from v0.2.0.

### Added

- Spatial Format-v3 phase-consensus measurement across independent complete
  tiles for every retained projective scale candidate. Aggregate sync peaks no
  longer decide geometry by themselves.
- Bounded phase-aware scale refinement: at most three seeds, each evaluated at
  the base scale plus ±1%/±2% width and height perturbations (maximum 27 scale
  evaluations total).
- Deterministic phase-correspondence DLT refinement. Four local phase
  observations can refine the coarse homography, with canonical corrections
  bounded to 64 pixels per axis.
- Phase coherence, refinement budgets, refined candidates and phase-DLT evidence
  in `diagnose -json`.
- Full virtual Format-v3 decodes remain capped at four and HMAC remains the only
  success criterion.

### Real-corpus research status

- The build2 lattice/boundary results are preserved: about 0.760 global
  consistency on `foto stampa.jpg` and 0.812 on `foto stampa storta.jpg`, both
  with lattice evidence.
- Two independent captures contain a closely matching coarse scale family near
  1668x1250 / 1653x1254. Build3 discovers and measures that family from local
  phase consistency; it is not hard-coded and is not assumed to be the true
  carrier size.
- On the frontal capture, phase-DLT refinement of that family raises its sparse
  known-header probe from z≈3.50 to z≈4.45. On the inclined capture the bounded
  phase search also improves the best sparse probe (to about z≈3.97).
- Neither photograph authenticates yet. The hidden real payload therefore
  remains unknown and `print-camera-test` correctly remains a research failure.

### Preserved

- Format v3 and deterministic encoder fingerprints.
- Production `ExtractWithInfo` behavior/search order.
- Eight coarse projective candidates and at most four complete virtual v3
  decodes.
- Private-corpus policy: real smartphone photographs are never distributed.

## v0.3.0-build2 — 2026-09-09

Second print-camera research checkpoint. Format v3, the deterministic encoder and
the production `ExtractWithInfo` search banks remain unchanged from v0.2.0.

### Added

- Pure-Go coarse print-boundary estimator with explicit confidence and corners.
  The quadrilateral is an initializer only and never watermark evidence.
- Boundary-aware normalization of local lattice bases and cross-level native-scale
  support to reduce scene-texture and single-pyramid-level harmonics.
- Bounded projective scale clustering and explicit canonical-to-photo homography.
- Virtual projective DCT sampler that avoids materializing a full rectified copy.
- Diagnostic key-known Format-v3 header probes for candidate ordering, capped at
  eight geometry probes and four complete virtual Format-v3 decodes.
- Synthetic end-to-end projective authentication regression, including wrong-key
  rejection.
- `make homography-test`; it is included in `make all-test` as a research target.
- Additional projective-fit/authentication timings and budgets in diagnose JSON.

### Real-corpus research status

- `foto stampa.jpg`: automatic print boundary, global lattice consistency about
  0.760, `lattice_evidence=true`.
- `foto stampa storta.jpg`: automatic print boundary, global lattice consistency
  about 0.812, `lattice_evidence=true`; build1 was about 0.505/false.
- Both 200 MP captures now reach the bounded virtual projective authentication
  path without a full rectified image. No Format-v3 HMAC is valid yet, therefore
  **the hidden real payload remains unrecovered and print-camera PASS is not
  claimed**.
- Key-known header scores are explicitly treated as multiple-testing-affected
  research evidence, not as proof that a watermark is present.

### Preserved

- Format v3 and deterministic encoder fingerprints.
- Production `ExtractWithInfo` behavior/search order.
- Private-corpus policy: real smartphone photographs are never distributed.

## v0.3.0-build1 — 2026-09-09

First research checkpoint toward print -> paper -> smartphone recovery. Format
v3, the deterministic encoder and the production `ExtractWithInfo` decoder are
unchanged from the qualified v0.2.0 baseline.

### Added

- Separate `watermark.DiagnoseGeometry` API and `pixseal diagnose` CLI command.
- Bounded sampled-luminance diagnostic pyramid for large inputs.
- Local `u/v` basis, phase, DCT differential margin, periodic coherence and
  35x32 v3 tile-repetition measurements.
- Per-region candidate pools with cross-region geometric consensus and explicit
  handling of 90-degree square-lattice basis equivalence.
- Explicit coarse/refined/repetition/phase/sample budgets and per-stage timings
  in JSON output.
- Optional independent baseline HMAC attempt when a key is supplied. Diagnostic
  lattice evidence cannot by itself authenticate a payload.
- `make lattice-estimator-test`.
- `make print-camera-test`, which uses the two private original smartphone
  photographs when present, SKIPs cleanly when absent, and reports PASS only for
  an authenticated Format v3 payload.

### Research status

- Canonical marked synthetic/digital carriers can produce coherent local lattice
  evidence while corresponding unmarked inputs remain negative in the build1
  regression path.
- Arbitrary-angle automatic local-basis estimation is not yet reliable enough to
  feed homography fitting; this remains the next geometry milestone.
- No claim of real print-camera payload recovery is made by build1.

## v0.2.0 — 2026-09-09

Final release of the Format v3 line, promoted from RC4 after qualification on
the private two-image corpus. No Format v3, encoder fingerprint or decoder
search-bank change was made after RC4.

### Qualification

- `release-unit`: PASS.
- `test-images`: PASS on both private qualification images and all three explicit
  profiles.
- `deep-test`: 72/72 PASS with zero failures, timeouts, skips or errors.
- `core-target-check`: PASS for linux/amd64, windows/amd64, android/arm64 and
  ios/arm64.
- Experimental reference results: affine 22/24, composition 2/2, lattice 8/8,
  perspective 4/4.
- The first RC4 all-test run produced two 60-second geometry timeouts; a targeted
  rerun with a 120-second geometry extraction timeout recovered both and restored
  the RC2 reference result of 65/120 with zero timeouts.

### Finalization

- Promoted version metadata from `v0.2.0-rc4` to `v0.2.0`.
- Recorded final qualification and retained the release/research separation.
- Set the geometry research harness default extraction timeout to 120 seconds to
  avoid misclassifying known recoverable cases as wall-clock regressions on the
  qualification host.

## v0.2.0-rc4 — 2026-09-09

Release-engineering reconciliation candidate. **Format v3, deterministic encoder
fingerprints and decoder search banks are unchanged.**

### Fixed

- Restored the RC2 `release-unit` / `research-unit` split so experimental Go
  geometry regressions cannot make the stable release baseline red.
- Made `make release-check` self-contained with target-specific `STRICT=1`.
- Made strict baseline qualification fail when every configured transformation
  is skipped and zero baseline transformations are executed.
- Made partial `ALL_TEST_TARGETS` runs report `PARTIAL` / `NOT RUN` instead of a
  misleading full qualification PASS.
- Handle subcommand `flag.ErrHelp` as successful CLI help (exit 0).
- Make `capacity` use `image.DecodeConfig` only; it no longer decodes the full
  bitmap merely to calculate geometry capacity.
- Restored `extract -raw` contract: exact payload bytes on stdout, diagnostics on
  stderr; geometric shell harnesses again capture the two streams separately.
- Restored numeric validation of `PERSPECTIVE_MAX_MPIX`.
- Restored RC2 CLI oversized-config and forced-permission regression coverage and
  added a true 64-bit pixel-count overflow regression.
- Treat post-publication temporary-file removal as best-effort cleanup rather
  than converting a successful no-clobber commit into a false write failure.
- Restore the v0.2 public/core source-size policy to 300,000,000 pixels while
  retaining RC3 overflow-safe arithmetic.
- Correct direct-lattice documentation to a 24 full-grid overall worst case
  (12 per sequential shape group), not 12 globally.
- Restore accurate RC1 → RC2 → RC3 chronology and attribute the qualification
  report to RC2.

### Preserved from RC3

- Cross-platform no-clobber publication with hard-link first and exclusive-create
  copy fallback when hard links are unavailable.
- Regular-file-only forced replacement, umask-safe new files and permission
  preservation on replacement.
- Shared alpha-on-white flattening between encoder and analyzer.
- Overflow-safe source-size checks and strengthened projective E2E regression.

## v0.2.0-rc3 — 2026-09-09

Second hardening candidate. RC3 improved cross-platform output publication,
overflow checking, alpha flattening and the synthetic perspective regression.
A subsequent audit found that assembly of RC3 had also reverted several valid
RC2 release-engineering and shell-harness changes; RC4 reconciles those changes.

## v0.2.0-rc2 — 2026-09-09

Implemented the first independent audit findings without changing Format v3:
finite strength validation, explicit CLI strength contract, dotfile handling,
`extract -raw`, large-image preflight, safer output permissions/no-clobber logic,
perspective harness fixes, an end-to-end projective regression, release/research
Go-test separation and major documentation cleanup.

PJ's RC2 qualification run reported release baseline PASS, `deep-test` 72/72,
composition 2/2, lattice 8/8 and perspective 4/4. Experimental geometry and
affine limits remained visible at 65/120 and 22/24 respectively.

## v0.2.0-rc1 — 2026-09-08

Consolidated the qualified build-11 algorithmic line and froze Format v3 encoder
fingerprints for release hardening.

## Development builds 1–12

See [`HISTORY.md`](HISTORY.md) for the complete technical progression.
