# Build106 — project audit and build consistency

Date: 2026-10-30. Snapshot: **v0.3.0-build106**. This maintenance build preserves the qualified Build98 engine and previous Build93 baseline. It adds no new physical recovery claim.

## Findings and corrections

- Android metadata drifted: root version/buildinfo and the manifest declared Build105, the checked-in Fyne configuration declared Build107, and the packager hard-coded Build105. The existing APK actually contained versionCode105. Build106 synchronizes all source identities and derives package arguments/output names from `VERSION`.
- Fyne modifies its configuration during packaging. The packager now builds an isolated temporary copy with an absolute local module replacement and removes that copy on success/failure. Checked-in configuration is preserved.
- ZIP/signature/permission checks alone did not verify the APK's release identity. Source checks now reject mismatched manifest/Fyne versions; final APK validation also requires the expected application ID, versionName, versionCode and ARM64 target before installation in `dist`.
- README, architecture/design documentation and Wiki FAQ/onboarding/portability pages still described the GUI as planned or absent. These now describe the existing native Android test app and distinguish package verification from real-device qualification.

## Verification

| Check | Result |
|---|---|
| `make gui-test` | PASS: source/package metadata regressions, inherited Build98 exactness/race, frozen v3, v4 pilot lock, service/UI race tests, headless layouts, GUI vet, shell syntax |
| `make vet release-unit` | PASS |
| `make release-check` | PASS: qualified Go1.26.0, CLI/version/JPEG and scoped release regressions, v3/pilot locks, corpus integrity, image round trips, JPEG compatibility, bounded transformations, core portability |
| Digital image round trips | 2 eligible images × 3 profiles PASS; the 201MP source is skipped by the configured 100MP test bound |
| Baseline transformations | 66 PASS, 0 failures, 0 timeouts, 0 errors, 33 expected size-bound SKIPs; 99 total |
| Core portability | linux/amd64, windows/amd64, android/arm64, ios/arm64 PASS |
| APK | ZIP integrity, ARM64 library/notices, versionCode106/versionName0.3.0/app ID, alignment, v1/v2/v3 signatures and absence of Internet permission PASS |
| Update signing identity | Certificate matches the retained Build105 APK |
| Interrupted packaging | Injected Fyne configuration mutation followed by exit42: original metadata/APK preserved and staging removed |
| `make build build-all` | Native Linux CLI and Linux/Windows/macOS desktop artifacts created |

The complete historical research suite and private nine-photo physical campaign were not rerun. No watermark engine source was modified. These checks preserve source compatibility and establish build health, not Android physical qualification.

## Artifacts

- `dist/PixSeal-build106-android-arm64.apk` and `.sha256` sidecar, locally signed for testing.
- `dist/pixseal` and the `pixseal-linux-amd64`, `pixseal-windows-amd64.exe`, `pixseal-macos-amd64`, `pixseal-macos-arm64` desktop files.
- Artifact hashes: [build106-artifacts.sha256](qualified-baselines/build106-artifacts.sha256).
- Local detailed logs: ignored `reports/build106/` directory.

## Outstanding work

Installation, Android document providers, keyboard/lifecycle, memory use and a create/read round trip on a real phone remain unverified. The Android photographic matrix must then be established. B/angle remains the known blind-registration gap; five known positive cases and three negative controls retain their historical desktop expectations. Android/iOS core compilation is separate from app/device qualification. No new recovery algorithm or promotion is asserted.

## Local storage cleanup — 2026-10-30

After build verification, remove 11 large extracted historical files only after matching their SHA256 against the corresponding intact retained ZIP archive. Historical test binaries and Build103 replay snapshots remain recoverable from those archives; restore them before using checks that require the original extracted paths. Keep all original photos, fixtures, small raw tables and evidence archives.

Remove the superseded Build105 APK/sidecar, the stale root test executable and two accidentally named Build89 output files (`clean`, `report_build_89_source.txtmake`). Build106 deliverables and all six artifact hashes remain valid. The project cleanup reclaims approximately 232 MB; deleting the temporary Go build cache used during verification reclaims a further 1.07 GB. Preserve the signing keystore and the Android/Fyne/JDK build tools. Detailed original paths, archive members and hashes are recorded locally in ignored `reports/build106/cleanup.json`.
