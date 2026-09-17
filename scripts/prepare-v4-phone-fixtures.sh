#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
SOURCE_DIR="${V4_PHONE_SOURCE_DIR:-original pics}"
OUTPUT_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
ROLE="${V4_PHONE_ROLE:-MQ}"
MESSAGE_A="${V4_PHONE_MESSAGE_A:-v4-b38-phone-a}"
MESSAGE_B="${V4_PHONE_MESSAGE_B:-v4-b38-phone-b}"
PROFILE="${V4_PHONE_PROFILE:-robust}"
STRENGTH="${V4_PHONE_STRENGTH:-48}"
PRINT_PPI="${V4_PHONE_PRINT_PPI:-300}"
ACTIVE_CORPUS_MANIFEST="${ACTIVE_CORPUS_MANIFEST:-private-corpus-active.tsv}"

[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
[[ -d "$SOURCE_DIR" ]] || { echo "error: v4 phone source directory not found: $SOURCE_DIR" >&2; exit 1; }
[[ -f "$ACTIVE_CORPUS_MANIFEST" ]] || { echo "error: active corpus manifest not found: $ACTIVE_CORPUS_MANIFEST" >&2; exit 1; }
if [[ "$PROFILE" == "robust" && ( ${#MESSAGE_A} -gt 16 || ${#MESSAGE_B} -gt 16 ) ]]; then
    echo "error: robust phone payloads must be at most 16 bytes" >&2
    exit 1
fi
if ! [[ "$PRINT_PPI" =~ ^[0-9]+$ ]] || (( PRINT_PPI < 72 || PRINT_PPI > 1200 )); then
    echo "error: V4_PHONE_PRINT_PPI must be an integer from 72 to 1200" >&2
    exit 1
fi

row="$(awk -F '\t' -v role="$ROLE" '$1==role {print; exit}' "$ACTIVE_CORPUS_MANIFEST")"
[[ -n "$row" ]] || { echo "error: role $ROLE not found in $ACTIVE_CORPUS_MANIFEST" >&2; exit 1; }
IFS=$'\t' read -r role source_basename source_w source_h source_sha <<< "$row"
src="$SOURCE_DIR/$source_basename"
[[ -f "$src" ]] || { echo "error: active corpus source missing: $src" >&2; exit 1; }
actual_sha="$(sha256sum "$src" | awk '{print $1}')"
[[ "$actual_sha" == "$source_sha" ]] || { echo "error: source SHA-256 mismatch for $source_basename" >&2; exit 1; }

canonical_w=$(( source_w / 8 * 8 ))
canonical_h=$(( source_h / 8 * 8 ))
if (( canonical_w < 1184 || canonical_h < 1024 )); then
    echo "error: role $ROLE normalizes to ${canonical_w}x${canonical_h}; phone qualification requires at least 4x4 complete v4 tiles" >&2
    exit 1
fi

if command -v magick >/dev/null 2>&1; then crop_cmd=(magick)
elif command -v convert >/dev/null 2>&1; then crop_cmd=(convert)
else echo "error: ImageMagick (magick or convert) is required" >&2; exit 1
fi

mkdir -p "$OUTPUT_DIR"
normalized="$OUTPUT_DIR/pixseal-build38-${ROLE,,}-canonical-${canonical_w}x${canonical_h}.png"
"${crop_cmd[@]}" "$src" -crop "${canonical_w}x${canonical_h}+0+0" +repage "$normalized"
control="$OUTPUT_DIR/pixseal-build38-${ROLE,,}-control.png"
marked_a="$OUTPUT_DIR/pixseal-build38-${ROLE,,}-marked-a.png"
marked_b="$OUTPUT_DIR/pixseal-build38-${ROLE,,}-marked-b.png"
cp "$normalized" "$control"

"$PIXSEAL" v4-embed -in "$normalized" -out "$marked_a" -key "$KEY" -message "$MESSAGE_A" -profile "$PROFILE" -strength "$STRENGTH" -force >/dev/null
"$PIXSEAL" v4-embed -in "$normalized" -out "$marked_b" -key "$KEY" -message "$MESSAGE_B" -profile "$PROFILE" -strength "$STRENGTH" -force >/dev/null

for file in "$control" "$marked_a" "$marked_b"; do
    tmp="$file.density.png"
    "${crop_cmd[@]}" "$file" -units PixelsPerInch -density "$PRINT_PPI" "$tmp"
    mv "$tmp" "$file"
done

for spec in "$marked_a|$MESSAGE_A" "$marked_b|$MESSAGE_B"; do
    IFS='|' read -r file expected <<< "$spec"
    got="$("$PIXSEAL" v4-extract -in "$file" -key "$KEY" -raw 2>/dev/null)" || { echo "error: phone print carrier failed digital preflight: $file" >&2; exit 1; }
    [[ "$got" == "$expected" ]] || { echo "error: phone carrier payload mismatch for $file" >&2; exit 1; }
done
if "$PIXSEAL" v4-extract -in "$control" -key "$KEY" -raw >/dev/null 2>&1; then
    echo "error: unmarked phone control unexpectedly authenticated before printing" >&2
    exit 1
fi
rm -f "$normalized"

pilot_hash='858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174'
manifest="$OUTPUT_DIR/build38-v4-phone-fixtures.tsv"
printf 'fixture\tsource_role\tsource_basename\twidth\theight\tsha256\texpected\tpayload\tprofile\tstrength\tprint_ppi\tpilot_hash\n' > "$manifest"
for spec in "$control|REJECT|-|-|-" "$marked_a|ACCEPT|$MESSAGE_A|$PROFILE|$STRENGTH" "$marked_b|ACCEPT|$MESSAGE_B|$PROFILE|$STRENGTH"; do
    IFS='|' read -r file expected payload profile strength <<< "$spec"
    sha="$(sha256sum "$file" | awk '{print $1}')"
    printf '%s\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$(basename "$file")" "$ROLE" "$source_basename" "$canonical_w" "$canonical_h" "$sha" "$expected" "$payload" "$profile" "$strength" "$PRINT_PPI" "$pilot_hash" >> "$manifest"
done

plan="$OUTPUT_DIR/acquisition-plan.tsv"
printf 'fixture\tfront\tmild\tangle\texpected\tpayload\n' > "$plan"
printf '%s\tphone-control-front.jpg\tphone-control-mild.jpg\tphone-control-angle.jpg\tREJECT\t-\n' "$(basename "$control")" >> "$plan"
printf '%s\tphone-a-front.jpg\tphone-a-mild.jpg\tphone-a-angle.jpg\tACCEPT\t%s\n' "$(basename "$marked_a")" "$MESSAGE_A" >> "$plan"
printf '%s\tphone-b-front.jpg\tphone-b-mild.jpg\tphone-b-angle.jpg\tACCEPT\t%s\n' "$(basename "$marked_b")" "$MESSAGE_B" >> "$plan"

width_cm="$(awk -v px="$canonical_w" -v ppi="$PRINT_PPI" 'BEGIN {printf "%.2f", px/ppi*2.54}')"
height_cm="$(awk -v px="$canonical_h" -v ppi="$PRINT_PPI" 'BEGIN {printf "%.2f", px/ppi*2.54}')"
cat > "$OUTPUT_DIR/README.txt" <<EOF2
PixSeal v0.3.0-build38 private Format-v4 smartphone qualification pack

Purpose
-------
Build37 closed blind print -> scanner -> HMAC at strength 24. The first 9-photo smartphone corpus showed that strength 24 is not sufficient for a single-shot phone channel: even reference-assisted dense registration left substantial protected-bit errors. Build38 therefore keeps the v4 wire format, pilot, Hamming and HMAC unchanged and raises only the physical phone qualification strength to $STRENGTH.

Canonical carrier
-----------------
Source role: $ROLE
Source: $source_basename
Block-normalized size: ${canonical_w}x${canonical_h} pixels
Print density target: $PRINT_PPI ppi
Target printed artwork size: ${width_cm} x ${height_cm} cm

Marked carriers
---------------
Profile:  $PROFILE
Strength: $STRENGTH
Payload A: $MESSAGE_A
Payload B: $MESSAGE_B
Key:      $KEY (intentional reproducible test key)
Pilot:    prototype-2-search-p64
Pilot SHA-256: $pilot_hash

Acquisition protocol
--------------------
1. Print control, marked-a and marked-b at ACTUAL SIZE / 100% using the same paper/printer settings as the successful Build35/37 scanner run. Disable fit-to-page and automatic photo enhancement where possible.
2. Use the main 1x phone camera, native maximum resolution, no digital zoom, no flash and no document mode/filter.
3. For each sheet take three original camera files: front, mild perspective (~10-15 degrees), and angle (~25-30 degrees). Keep the whole artwork plus visible white paper around all four sides.
4. Do not crop, rotate, resave, message-app compress or screenshot the files. Preserve the original JPEG/HEIC and metadata.
5. Store acquisitions under v4-phone private/build38-acquired/ with exactly the names in acquisition-plan.tsv.
6. Do not interpret pilot correlation alone as PASS. A marked photo passes only when the exact payload authenticates through the unchanged Format-v4 HMAC; the control must reject.

The previous strength-24 phone photographs remain a private historical corpus. Do not overwrite them.
EOF2

echo "created Build38 phone pack in $OUTPUT_DIR"
echo "canonical carrier: ${canonical_w}x${canonical_h} px; strength $STRENGTH; print at $PRINT_PPI ppi -> ${width_cm} x ${height_cm} cm"
echo "fixture manifest: $manifest"
echo "acquisition plan: $plan"
