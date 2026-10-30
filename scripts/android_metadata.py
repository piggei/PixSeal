#!/usr/bin/env python3
"""Validate Android source/package identity against the release VERSION."""
import argparse
from pathlib import Path
import re
import subprocess
import xml.etree.ElementTree as ET

APP_ID = "io.github.piggei.pixseal"
ANDROID = "{http://schemas.android.com/apk/res/android}"


def source_metadata(root):
    match = re.fullmatch(r"PixSeal v(\d+\.\d+\.\d+)-build([1-9]\d*)", (root / "VERSION").read_text().strip())
    if not match:
        raise ValueError("VERSION must contain PixSeal vX.Y.Z-buildN")
    version, build = match.groups()
    config = (root / "gui/FyneApp.toml").read_text()
    for name, value in (("Version", f'"{version}"'), ("Build", build), ("ID", f'"{APP_ID}"')):
        if not re.search(rf"^{name}\s*=\s*{re.escape(value)}\s*$", config, re.M):
            raise ValueError(f"FyneApp.toml {name} differs from VERSION/app identity")
    manifest = ET.parse(root / "gui/AndroidManifest.xml").getroot()
    if (manifest.get("package"), manifest.get(ANDROID + "versionName"), manifest.get(ANDROID + "versionCode")) != (APP_ID, version, build):
        raise ValueError("AndroidManifest.xml identity differs from VERSION")
    if any(e.get(ANDROID + "name") == "android.permission.INTERNET" for e in manifest.findall("uses-permission")):
        raise ValueError("Unexpected network permission in offline source manifest")
    return version, build


def check_package(badging, permissions, version, build):
    match = re.search(r"^package: name='([^']+)' versionCode='([^']+)' versionName='([^']+)'", badging, re.M)
    if not match or match.groups() != (APP_ID, build, version):
        raise ValueError("APK identity/version differs from VERSION")
    if "native-code: 'arm64-v8a'" not in badging:
        raise ValueError("APK must contain the Android ARM64 target")
    if "android.permission.INTERNET" in permissions:
        raise ValueError("Unexpected network permission in offline package")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--apk", type=Path)
    parser.add_argument("--aapt")
    args = parser.parse_args()
    version, build = source_metadata(Path(__file__).resolve().parent.parent)
    if args.apk:
        if not args.aapt:
            parser.error("--apk requires --aapt")
        badging = subprocess.check_output([args.aapt, "dump", "badging", str(args.apk)], text=True)
        permissions = subprocess.check_output([args.aapt, "dump", "permissions", str(args.apk)], text=True)
        check_package(badging, permissions, version, build)
    print(version, build)


if __name__ == "__main__":
    main()
