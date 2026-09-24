# Private physical corpus

The physical print/camera and scanner acquisitions are research material and are
**not distributed with PixSeal source archives**. The PolyForm Noncommercial
license applies to the PixSeal source distribution; it does **not** make these
private photographs, scans, generated fixtures or acquisition archives part of
the licensed source package. Unless a separate notice explicitly grants rights
to a private corpus, keep it private and out of releases. See
[`LICENSING.md`](LICENSING.md). Their test key is intentionally
public and reproducible: `Piccotti`. Build17 keeps it as the default
`PRINT_CAMERA_KEY` so the historical corpus can be decoded without extra setup.
It is a project test constant, not a personal or production credential. If PixSeal is
ever promoted to a production deployment, replace this Makefile/script default with a
different production key. The private research corpus itself is not part of production
artifacts, so production compatibility with `Piccotti` is not required.

Canonical smartphone filenames:

```text
foto stampa.jpg
foto stampa storta.jpg
foto bici dritta.jpg
foto bici storta.jpg
```

Canonical scanner filenames:

```text
0270_001.jpg
0270_002.jpg
```

The two historical `foto stampa*.PNG` files contained JPEG data despite their
extension. The canonical corpus names use `.jpg`; renaming changes only the filename,
not the image bytes.

Generate a local SHA-256 manifest after placing the private files in the default
folders:

```bash
make private-corpus-manifest
```

The output `private-corpus-manifest.sha256` is ignored by Git and records canonical
relative labels rather than machine-specific absolute paths. Keep the manifest with
the private corpus if you need to verify that two research runs used byte-identical
acquisitions.

Physical regression examples:

```bash
make print-camera-test        # defaults to PRINT_CAMERA_KEY=Piccotti
make print-scan-test           # defaults to PRINT_CAMERA_KEY=Piccotti
PRINT_CAMERA_KEY='other-test-key' make print-camera-test  # optional override
```

Only a valid Format-v3 HMAC counts as a PASS. Geometry, lattice, error-count and
unwrap diagnostics remain research evidence.

## Build20

Build20 continues to use `Piccotti` as the intentional public/reproducible research key. The
private acquisition files remain excluded from source archives. The independent pairwise cycle
anchor was evaluated on the same four available physical acquisitions; the two bicycle photos
are still absent from the current corpus. A future production release must replace the default
research key and will not distribute this private corpus.


## Build21

Build21 still uses the intentional public research key `Piccotti`; no acquisition file is added
to the repository or source archive. The new Format-v3 observability audit is static and can run
without the private corpus. Physical print-camera/scan tests remain useful only to confirm that
adding this telemetry does not change the four existing diagnostic verdicts or HMAC behavior.
The two bicycle photographs remain required for a complete six-image physical qualification.


### Build22 note

The held-out physical-topology probe and any future Format-v4 pilot prototype may be evaluated
against the private acquisition corpus, but photographs/scans remain excluded from source/evidence
archives. Format-v4 design telemetry is static and may be distributed; only derived numeric JSON/CSV
evidence from private acquisitions may be packaged.

## Format-v4 corpus separation

The existing smartphone/scanner acquisitions are Format-v3 evidence only. They must not be reused as
if they tested Format v4, because the printed carrier does not contain the v4 pilot. When an experimental
v4 encoder exists, create a separate private v4 acquisition set by re-embedding, reprinting and reacquiring
the carrier. Keep v3 and v4 manifests/directories distinct so a future result cannot mix formats.

### Build24 local-original note

The two `original pics` images used by `v4-pilot-corpus-test` are digital source images, not physical Format-v4
acquisitions. Build24 may generate temporary synthetic pilot/data-plane carriers from them and report derived
metrics, but source ZIPs/evidence ZIPs should not include the image files themselves. They do not satisfy the v4
physical-corpus requirement. After a real experimental v4 payload encoder is available, create new embedded
originals, print them and acquire separate smartphone/scanner samples under a v4-specific manifest.



### Build25 local-original geometry note

Build25 reuses the local original images only as backgrounds for synthetic v4 pilot/data-plane carriers and deterministic geometric transforms. The originals and transformed derivatives remain local and must not be included in source/evidence archives. Evidence may contain only textual metrics/logs. These experiments are not physical v4 print-camera evidence.


## Build31 private Format-v4 physical corpus

Build31 is the first version capable of creating a genuine v4 carrier. The existing private v3 print-camera/scanner photographs **must not** be relabeled as v4 evidence; they were printed from a different on-image format and contain no prototype-2 pilot or Build31 v4 frame.

For the first v4 physical experiment, generate fresh PNG carriers with the explicit experimental command, preferably using the reproducible test key `Piccotti`, robust profile and default strength 24:

```sh
pixseal v4-embed \
  -in INPUT.png \
  -out OUTPUT-v4-build31.png \
  -key Piccotti \
  -message "PixSeal Build31 physical v4" \
  -profile robust \
  -strength 24
```

Print those exact generated PNGs, then photograph and/or scan the paper. Keep the generated digital carrier together with each acquisition and record printer, paper, print scaling, acquisition device and whether any editor rescaled/cropped the image. Physical source/acquisition files remain private and must not enter source/evidence release archives.

A future physical result counts as v4 success only when the Build31 frame reaches **HMAC-authenticated payload recovery**. Public pilot detection, geometry score, CRC or ECC success alone are diagnostics. Build31 `v4-extract` currently handles aligned native lattices only, so the first physical decoder integration remains a later research checkpoint.

Build31 also provides a reproducible helper:

```sh
make v4-physical-fixtures
```

By default it reads `original pics/` and writes only generated private carriers plus a TSV SHA-256 manifest under `v4-physical private/build31-generated/`. That output tree is git-ignored. Override `V4_PHYSICAL_SOURCE_DIR`, `V4_PHYSICAL_OUTPUT_DIR`, key/message/profile/strength Make variables when a different experiment is intended.


## Build32 active corpus policy

The active v4 development/qualification corpus is now exactly the three anonymous originals listed in root `private-corpus-active.tsv`: LQ, MQ and HQ. Corpus-driven tests must consume that manifest and must not discover arbitrary directory contents. `corpus-manifest-check` verifies dimensions and SHA-256. HQ is approximately 200.5 MP and is intentionally skipped by bounded tests whose configured image budget is lower; this is a SKIP, not a PASS.

New v4 experiments and physical fixtures use the intentionally public development key `PixSeal-v4-TestKey-2026`. The historical key `Piccotti` remains only where older regression vectors require it. Neither the three originals nor generated print fixtures belong in source/evidence release archives.


## Build41 smartphone qualification corpus

Build41 reuses the nine original strength-48 files under `v4-phone private/build38-acquired/`. They remain private qualification material and are never copied into source/release archives. `make v4-build41-phone-physical-test` may write only derived TSV/Markdown telemetry and stderr logs under `v4-phone private/build41-diagnostics/`; those outputs are private/generated evidence and are covered by the existing `/v4-phone private/` ignore rule.

The software's PolyForm Noncommercial License does not grant separate rights to private photographs, scans, source images or qualification fixtures. Keep those corpora outside release archives and public commits unless their owner explicitly chooses to publish them under a separate license.


## Build42 smartphone qualification corpus

Build42 uses the same nine private strength-48 originals under `v4-phone private/build38-acquired/`; no new acquisition is required. `make v4-build42-phone-physical-test` writes only derived telemetry and stderr logs under `v4-phone private/build42-diagnostics/`. The required gate is: all three controls reject before data decode; A/angle and B/front retain direct Build41 authentication; A/mild authenticates through the Build42 data-list fallback. A/front, B/mild and B/angle remain informational geometry-reject rows. No image bytes are copied into source or release archives.


## Build43 private smartphone gate

The official Build43 physical gate must be built/run through the repository Makefile, which selects the qualified Go 1.25.1 toolchain. Go 1.26+ is not a qualified Build43 JPEG decoder for this corpus; see `docs/GO_TOOLCHAIN_COMPATIBILITY.md`.


Build43 reuses the same nine original Build38 strength-48 photographs. `v4-build43-phone-physical-test` writes derived logs/matrices only. Required gate: controls reject; A/front, A/mild, A/angle and B/front authenticate; B/mild/B-angle are informational. No private image is distributed.

## Build44 Go 1.26 ingest qualification

Build44 reuses the exact same nine Build38 strength-48 smartphone photographs;
no reprint or reacquisition is required. The candidate target is:

```sh
make v4-build44-go126-phone-physical-test
```

It builds a separate Go 1.26.0 binary using the deterministic PixSeal JPEG
decoder and writes derived logs/matrices under
`v4-phone private/build44-diagnostics/`. The required matrix is identical to
Build43: controls 3/3 reject; A/front, A/mild, A/angle and B/front authenticate;
B/mild/B-angle remain informational. The private image bytes are never packaged.

### Build46 diagnostics

Build46 writes only local research output under `v4-phone private/build46-diagnostics/`. The directory may contain blind handoff JSON/TSV/Markdown plus optional reference-assisted quadrilateral files generated after the blind runs. These files are private evidence and must not be included in source archives or releases.
## Build48 private local-refinement diagnostics

Build48 reuses the same retained Build38 B/mild and B/angle JPEGs and marked-B digital reference. Private output belongs under `v4-phone private/build48-diagnostics/` and must never be packaged. Blind `*-refine.json` files are written first. Only after both exist may the lab helper create oracle quadrilaterals and the pre/post error tables. B/mild is the primary Build48 research case; B/angle remains informational.



## Build49 private proposal-ranking diagnostics

Build49 reuses private `phone-b-mild.jpg`, `phone-b-angle.jpg` and the marked-B digital reference. Blind ranking JSON files are written first under `v4-phone private/build49-diagnostics/blind/` without any secret key. Only after both exist may the SIFT/reference helper create oracle quadrilaterals and post-hoc rank/error summaries. No Build49 private output is shipped.


## Build50 private top-4 refinement diagnostics

Build50 reuses only the retained private `phone-b-mild.jpg`, `phone-b-angle.jpg` acquisitions and the marked-B digital reference. The blind top4 refinement JSON for both captures is written under `v4-phone private/build50-diagnostics/blind/` before any SIFT/reference oracle is created. The secret key is available to the blind command only for diagnostic single-candidate HMAC after geometry is frozen and held-out-qualified; it cannot guide proposal ranking or local refinement. Oracle files and reports remain private and are never packaged.

## Build51 private local-surface diagnostics

Build51 reuses only the retained `phone-b-mild.jpg`, `phone-b-angle.jpg` acquisitions and the marked-B digital reference already required by Builds45–50. Blind top4 refinement traces and deterministic local stencils for both captures are written under `v4-phone private/build51-diagnostics/blind/` before the reference/SIFT helper is invoked. The key is used only after geometry freeze for diagnostic authentication of final held-out-qualified states. The oracle is post-hoc only. Trace, stencil, oracle and summary outputs are private research artifacts and must never be packaged.


## Build52 private fine-restart diagnostics

Build52 reuses only retained `phone-b-mild.jpg`, `phone-b-angle.jpg` and the marked-B digital reference already required by Builds45–51. The blind `v4-diagnose-phone-restart` outputs for both captures are written under `v4-phone private/build52-diagnostics/blind/` before the reference/SIFT helper is invoked. Each blind record contains the unchanged coarse-to-fine baseline plus the proposal-only `2 -> 1 px` restart state bank generated from the same untouched top4 seeds. Held-out/full-pilot qualification and diagnostic HMAC are annotated only after geometry freeze; SIFT/reference error is post-hoc only. All Build52 JSON/TSV/Markdown diagnostic outputs remain private and must never be packaged in a source release.

## Build61/62 private deep-optimizer diagnostics

Build61 writes private sibling-stencil outputs under `v4-phone private/build61-diagnostics/`.
Build62 writes private fourth-pair escape outputs under `v4-phone private/build62-diagnostics/`.
Both reuse only the retained B/mild and B/angle acquisitions plus the marked-B reference.
Blind JSON for both images must exist before SIFT/reference oracle generation. Qualification
and diagnostic HMAC are attached only after the complete geometry bank is frozen. None of
these JSON/TSV/Markdown artifacts may be included in public source archives.

Build63 writes private post-fourth-pair continuation outputs under `v4-phone private/build63-diagnostics/`. These remain private research artifacts and are excluded from source/release archives.

## Build64 private production-candidate qualification

Build64 reuses the complete retained nine-photo Build38 smartphone corpus. Its physical-gate output belongs under `v4-phone private/build64-diagnostics/` and contains per-image stderr logs plus `build64-phone-matrix.tsv` / `.md`. No acquisition image, payload dump or Build64 private matrix is included in source/release archives. The Build64 candidate gate is explicitly private because it validates the new deep fallback against controls and the difficult B captures before any production-baseline promotion.
