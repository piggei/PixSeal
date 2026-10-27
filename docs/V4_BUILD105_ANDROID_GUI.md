# Build105 — native Android GUI and local trial workflow

**Native Fyne GUI / Android ARM64 test package. Qualified desktop phone algorithm remains Build98, previous Build93. Android physical qualification remains pending.**

The user changed the immediate research priority from speed micro-optimizations to functional recovery of B/angle, and requested a practical GUI operating on Android without a PC. Build105 delivers the native app workflow first, so new photographic trials can be performed directly on the phone. It does not implement a new B/angle geometry search or change the qualified recovery kernels.

## Architecture and behavior

`gui/` is a separate Go module using pinned Fyne2.8.2 and tools1.7.4. `internal/appservice` adapts the existing v4 APIs without any graphical or external dependencies in the root module. No qualified `watermark/` source is changed. The CLI remains independently buildable.

The Italian interface imports a JPEG/PNG via Android's document picker, displays a small preview, creates an authenticated v4 carrier, reads either the aligned digital carrier or a photograph of a print, and saves a PNG or JSON trial report through the platform document picker. It shows elapsed time and indeterminate activity, permits one job at a time and clears previous output before a new trial. Background engine work updates widgets through `fyne.Do`. Report output excludes the key and payload; the key is not written to preferences, reports or logs. Only public canonical dimensions are remembered. Payload copy requires an explicit user action after successful authentication.

Creation prepares a new image with preserved aspect ratio, at most1632 pixels on its long edge, dimensions aligned to eight and alpha composited onto white. It embeds with strength48/automatic profile and verifies an authenticated aligned round trip before offering the PNG. These preparation choices apply to newly created carriers and preview thumbnails. Imported acquired photographs reach the existing Build98 API without recompression, resizing or orientation changes by the app adapter; existing decoder behavior remains authoritative.

The import boundary checks actual PNG/JPEG signatures, uses the pinned legacy JPEG decoder, checks dimensions before full raster decoding, and caps imports at64 MiB /24 million pixels. Photo-mode canonical dimensions must be block aligned, at least296×256, no more than4096 on either axis and at most4 million pixels. These are new application resource limits, not changes to CLI limits. The create workflow naturally supplies compatible dimensions.

## Offline Android packaging

App ID: `io.github.piggei.pixseal`; app version0.3.0 / versionCode105; ARM64 native library. SDK34 build-tools, NDK26.1.10909125, local JDK17 and Fyne's mobile packager build the test APK. The custom manifest excludes Internet and broad storage permissions; user-selected content URIs provide file access. The package is aligned before signing, uses a locally generated test keystore excluded from Git, and is verified with Android apksigner. Native linking requests16 KiB segment alignment; this does not replace testing on actual4/16 KiB devices.

The final output is `dist/PixSeal-build105-android-arm64.apk` with a SHA256 sidecar. Intermediate binaries and tool caches stay out of source control. This is a sideloadable development test package, not a Play Store release or a production signing setup. Fyne/third-party licence notices are included with the application assets.

## Validation and limits

- Inherited Build98 exactness/race/CLI/buildinfo/vet, frozen v3 and pilot-lock gates.
- App service authenticated create/save/reload/read, wrong-key rejection, report privacy, JPEG raster equivalence, invalid canonical extents, UTF-8 byte limits, oversized-header rejection, downscale/alpha/bounds handling and race checks.
- Fyne headless UI tests for stale-output clearing, authentication enforcement even when a backend erroneously returns data, busy-state exclusion, result-copy action, and mobile layouts. Layout captures are software-rendered previews, not device screenshots.
- APK ZIP/native-library presence, signature checks and absence of Internet permission.

A real Android device has not been exercised by these automated checks. Import/save providers, keyboard/lifecycle behavior, device memory/thermal limits, and physical photographic acceptance must still be tested. Source equivalence to Build98 does not transfer its Go1.26.0/x86 host qualification automatically to Android ARM64. The known B/angle failure and all historical non-promotion decisions remain unchanged.

The first concrete device trial is the local digital create/read round trip. Next run the existing five authenticated photographs and three negative controls, while retaining B/angle as the known functional gap. B/angle recovery becomes the next geometry milestone after establishing that the native app reproduces the expected baseline on the target phone.

See [the Italian app guide](../gui/README.md) and `make gui-test` / `make gui-android`.
