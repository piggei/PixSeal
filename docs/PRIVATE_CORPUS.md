# Private physical corpus

The physical print/camera and scanner acquisitions are research material and are
**not distributed with PixSeal source archives**. Their test key is intentionally
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
