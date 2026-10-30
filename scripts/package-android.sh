#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
METADATA="$(python3 "$ROOT/scripts/android_metadata.py")"
read -r APP_VERSION APP_BUILD <<< "$METADATA"
ARTIFACT="PixSeal-build${APP_BUILD}-android-arm64.apk"
FYNE="${PIXSEAL_FYNE_TOOL:-$ROOT/.tools/gui/bin/fyne}"
if [[ ! -x "$FYNE" ]];then FYNE="$(command -v fyne || true)";fi
[[ -n "$FYNE" && -x "$FYNE" ]] || { echo 'Install the pinned Fyne tool: go install fyne.io/tools/cmd/fyne@v1.7.4' >&2;exit 1; }
if [[ -z "${JAVA_HOME:-}" && -x "$ROOT/.tools/gui/java/usr/lib/jvm/java-17-openjdk-amd64/bin/javac" ]];then export JAVA_HOME="$ROOT/.tools/gui/java/usr/lib/jvm/java-17-openjdk-amd64";fi
[[ -z "${JAVA_HOME:-}" ]] || export PATH="$JAVA_HOME/bin:$PATH"
command -v javac >/dev/null || { echo 'Java JDK 17 or newer is required (JAVA_HOME/bin/javac).' >&2;exit 1; }
export ANDROID_HOME="${ANDROID_HOME:-$HOME/Android/Sdk}"
export ANDROID_NDK_HOME="${ANDROID_NDK_HOME:-$ANDROID_HOME/ndk/26.1.10909125}"
[[ -d "$ANDROID_HOME/platforms" && -d "$ANDROID_NDK_HOME/toolchains" ]] || { echo 'Set ANDROID_HOME and ANDROID_NDK_HOME to the Android SDK and NDK.' >&2;exit 1; }
export CGO_LDFLAGS="${CGO_LDFLAGS:-} -Wl,-z,max-page-size=16384"
# Fyne increments FyneApp.toml while packaging. Build in an isolated copy so
# both successful and failed builds leave the checked-in metadata untouched.
mkdir -p "$ROOT/.tools/gui"
STAGING="$(mktemp -d "$ROOT/.tools/gui/package.XXXXXX")"
trap 'rm -rf "$STAGING"' EXIT
cp "$ROOT"/gui/*.go "$ROOT/gui/go.mod" "$ROOT/gui/go.sum" "$ROOT/gui/FyneApp.toml" "$ROOT/gui/AndroidManifest.xml" "$STAGING/"
cp -R "$ROOT/gui/assets" "$STAGING/assets"
cd "$STAGING"
go mod edit -replace "github.com/pj/pixseal=$ROOT"
"$FYNE" package --target android/arm64 --name PixSeal --app-id io.github.piggei.pixseal --app-version "$APP_VERSION" --app-build "$APP_BUILD" --icon assets/Icon.png
python3 - <<'PY'
from pathlib import Path
import zipfile
p=Path('PixSeal.apk')
with zipfile.ZipFile(p) as z:
 assert z.testzip() is None
 assert {'AndroidManifest.xml','classes.dex','lib/arm64-v8a/libpixseal.so'} <= set(z.namelist())
PY
SIGNER="${PIXSEAL_APKSIGNER:-$ANDROID_HOME/build-tools/34.0.0/apksigner}"
[[ -x "$SIGNER" ]] || { echo 'Android apksigner is required; set PIXSEAL_APKSIGNER.' >&2;exit 1; }
# Local test identity; never commit the keystore. Reuse it for later test updates.
KEYSTORE="${PIXSEAL_ANDROID_KEYSTORE:-$ROOT/.tools/gui/android-debug.keystore}"
if [[ ! -f "$KEYSTORE" ]];then
 mkdir -p "$(dirname "$KEYSTORE")"
 (umask 077; keytool -genkeypair -keystore "$KEYSTORE" -alias pixseal-test -storepass android -keypass android -keyalg RSA -keysize 2048 -validity 10000 -dname 'CN=PixSeal Local Test')
fi
ZIPALIGN="${PIXSEAL_ZIPALIGN:-$ANDROID_HOME/build-tools/34.0.0/zipalign}"
AAPT="${PIXSEAL_AAPT:-$ANDROID_HOME/build-tools/34.0.0/aapt}"
"$ZIPALIGN" -f 4 PixSeal.apk PixSeal-aligned.apk
"$SIGNER" sign --ks "$KEYSTORE" --ks-key-alias pixseal-test --ks-pass pass:android --key-pass pass:android --v1-signing-enabled true --v2-signing-enabled true --v3-signing-enabled true --v4-signing-enabled false --out PixSeal-signed.apk PixSeal-aligned.apk
"$SIGNER" verify --verbose PixSeal-signed.apk
"$ZIPALIGN" -c 4 PixSeal-signed.apk
python3 "$ROOT/scripts/android_metadata.py" --apk "$STAGING/PixSeal-signed.apk" --aapt "$AAPT"
mkdir -p "$ROOT/dist"
cp PixSeal-signed.apk "$ROOT/dist/.$ARTIFACT.tmp"
mv -f "$ROOT/dist/.$ARTIFACT.tmp" "$ROOT/dist/$ARTIFACT"
(cd "$ROOT/dist" && sha256sum "$ARTIFACT" > "$ARTIFACT.sha256")
echo "Verified test APK: $ROOT/dist/$ARTIFACT"
