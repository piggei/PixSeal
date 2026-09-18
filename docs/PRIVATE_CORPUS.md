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
