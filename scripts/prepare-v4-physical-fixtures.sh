#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
SOURCE_DIR="${V4_PHYSICAL_SOURCE_DIR:-original pics}"
OUTPUT_DIR="${V4_PHYSICAL_OUTPUT_DIR:-v4-physical private/build35-generated}"
KEY="${V4_PHYSICAL_KEY:-PixSeal-v4-TestKey-2026}"
ROLE="${V4_PHYSICAL_ROLE:-MQ}"
MESSAGE_A="${V4_PHYSICAL_MESSAGE_A:-v4-b35-phys-a}"
MESSAGE_B="${V4_PHYSICAL_MESSAGE_B:-v4-b35-phys-b}"
PROFILE="${V4_PHYSICAL_PROFILE:-robust}"
STRENGTH="${V4_PHYSICAL_STRENGTH:-24}"
PRINT_PPI="${V4_PHYSICAL_PRINT_PPI:-300}"
ACTIVE_CORPUS_MANIFEST="${ACTIVE_CORPUS_MANIFEST:-private-corpus-active.tsv}"

if [[ ! -x "$PIXSEAL" ]]; then
    echo "error: PIXSEAL is not executable: $PIXSEAL" >&2
    exit 1
fi
if [[ ! -d "$SOURCE_DIR" ]]; then
    echo "error: v4 physical source directory not found: $SOURCE_DIR" >&2
    exit 1
fi
if [[ ! -f "$ACTIVE_CORPUS_MANIFEST" ]]; then
    echo "error: active corpus manifest not found: $ACTIVE_CORPUS_MANIFEST" >&2
    exit 1
fi
if [[ "$PROFILE" == "robust" && ( ${#MESSAGE_A} -gt 16 || ${#MESSAGE_B} -gt 16 ) ]]; then
    echo "error: robust physical payloads must be at most 16 bytes" >&2
    exit 1
fi
if ! [[ "$PRINT_PPI" =~ ^[0-9]+$ ]] || (( PRINT_PPI < 72 || PRINT_PPI > 1200 )); then
    echo "error: V4_PHYSICAL_PRINT_PPI must be an integer from 72 to 1200" >&2
    exit 1
fi

row="$(awk -F '\t' -v role="$ROLE" '$1==role {print; exit}' "$ACTIVE_CORPUS_MANIFEST")"
if [[ -z "$row" ]]; then
    echo "error: role $ROLE not found in $ACTIVE_CORPUS_MANIFEST" >&2
    exit 1
fi
IFS=$'\t' read -r role source_basename source_w source_h source_sha <<< "$row"
src="$SOURCE_DIR/$source_basename"
if [[ ! -f "$src" ]]; then
    echo "error: active corpus source missing: $src" >&2
    exit 1
fi
actual_sha="$(sha256sum "$src" | awk '{print $1}')"
if [[ "$actual_sha" != "$source_sha" ]]; then
    echo "error: source SHA-256 mismatch for $source_basename" >&2
    exit 1
fi

canonical_w=$(( source_w / 8 * 8 ))
canonical_h=$(( source_h / 8 * 8 ))
if (( canonical_w < 1184 || canonical_h < 1024 )); then
    echo "error: role $ROLE normalizes to ${canonical_w}x${canonical_h}; Build34 physical projective path requires at least 4x4 complete v4 tiles (1184x1024)" >&2
    exit 1
fi

if command -v magick >/dev/null 2>&1; then
    crop_cmd=(magick)
elif command -v convert >/dev/null 2>&1; then
    crop_cmd=(convert)
else
    echo "error: ImageMagick (magick or convert) is required to make a block-aligned print master" >&2
    exit 1
fi

mkdir -p "$OUTPUT_DIR"
normalized="$OUTPUT_DIR/pixseal-build35-${ROLE,,}-canonical-${canonical_w}x${canonical_h}.png"
"${crop_cmd[@]}" "$src" -crop "${canonical_w}x${canonical_h}+0+0" +repage "$normalized"

control="$OUTPUT_DIR/pixseal-build35-${ROLE,,}-control.png"
marked_a="$OUTPUT_DIR/pixseal-build35-${ROLE,,}-marked-a.png"
marked_b="$OUTPUT_DIR/pixseal-build35-${ROLE,,}-marked-b.png"
cp "$normalized" "$control"

"$PIXSEAL" v4-embed -in "$normalized" -out "$marked_a" -key "$KEY" -message "$MESSAGE_A" -profile "$PROFILE" -strength "$STRENGTH" -force >/dev/null
"$PIXSEAL" v4-embed -in "$normalized" -out "$marked_b" -key "$KEY" -message "$MESSAGE_B" -profile "$PROFILE" -strength "$STRENGTH" -force >/dev/null

# Store the intended physical density in the PNG pHYs metadata without changing
# pixel samples. The print protocol still requires ACTUAL SIZE / 100% because
# some applications ignore PNG density metadata.
for file in "$control" "$marked_a" "$marked_b"; do
    density_tmp="$file.density.png"
    "${crop_cmd[@]}" "$file" -units PixelsPerInch -density "$PRINT_PPI" "$density_tmp"
    mv "$density_tmp" "$file"
done

# Digital preflight: these exact print files must be valid before paper enters
# the experiment. The aligned decoder is intentional here because no physical
# geometry has been introduced yet.
for spec in "$marked_a|$MESSAGE_A" "$marked_b|$MESSAGE_B"; do
    IFS='|' read -r file expected_payload <<< "$spec"
    got="$("$PIXSEAL" v4-extract -in "$file" -key "$KEY" -raw 2>/dev/null)" || {
        echo "error: generated print carrier failed aligned v4 preflight: $file" >&2
        exit 1
    }
    if [[ "$got" != "$expected_payload" ]]; then
        echo "error: generated carrier payload mismatch for $file" >&2
        exit 1
    fi
done
if "$PIXSEAL" v4-extract -in "$control" -key "$KEY" -raw >/dev/null 2>&1; then
    echo "error: unmarked physical control unexpectedly authenticated before printing" >&2
    exit 1
fi

# The control is already the normalized canonical image; do not leave a fourth
# visually duplicate file in the print pack.
rm -f "$normalized"

pilot_hash='858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174'
manifest="$OUTPUT_DIR/build35-v4-physical-fixtures.tsv"
printf 'fixture\tsource_role\tsource_basename\twidth\theight\tsha256\texpected\tpayload\tprofile\tstrength\tprint_ppi\tpilot_hash\n' > "$manifest"
for spec in \
    "$control|REJECT|-|-|-" \
    "$marked_a|ACCEPT|$MESSAGE_A|$PROFILE|$STRENGTH" \
    "$marked_b|ACCEPT|$MESSAGE_B|$PROFILE|$STRENGTH"; do
    IFS='|' read -r file expected payload profile strength <<< "$spec"
    sha="$(sha256sum "$file" | awk '{print $1}')"
    printf '%s\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$(basename "$file")" "$ROLE" "$source_basename" "$canonical_w" "$canonical_h" "$sha" \
        "$expected" "$payload" "$profile" "$strength" "$PRINT_PPI" "$pilot_hash" >> "$manifest"
done

plan="$OUTPUT_DIR/acquisition-plan.tsv"
printf 'fixture\tcapture\twidth\theight\texpected\tpayload\n' > "$plan"
printf '%s\t%s\t%d\t%d\tREJECT\t-\n' "$(basename "$control")" 'scan-control.png' "$canonical_w" "$canonical_h" >> "$plan"
printf '%s\t%s\t%d\t%d\tACCEPT\t%s\n' "$(basename "$marked_a")" 'scan-marked-a.png' "$canonical_w" "$canonical_h" "$MESSAGE_A" >> "$plan"
printf '%s\t%s\t%d\t%d\tACCEPT\t%s\n' "$(basename "$marked_b")" 'scan-marked-b.png' "$canonical_w" "$canonical_h" "$MESSAGE_B" >> "$plan"

width_cm="$(awk -v px="$canonical_w" -v ppi="$PRINT_PPI" 'BEGIN {printf "%.2f", px/ppi*2.54}')"
height_cm="$(awk -v px="$canonical_h" -v ppi="$PRINT_PPI" 'BEGIN {printf "%.2f", px/ppi*2.54}')"
cat > "$OUTPUT_DIR/README.txt" <<EOF2
PixSeal v0.3.0-build35 private Format-v4 physical qualification pack

Purpose
-------
This pack is the first physical-channel gate after Build34 closed blind MQ projective+crop + HMAC recovery digitally.
The files are PRIVATE TEST ARTIFACTS and are intentionally excluded from source/evidence archives.

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
Key:      $KEY  (intentional public/reproducible test key)
Pilot:    prototype-2-search-p64
Pilot SHA-256: $pilot_hash

Initial acquisition protocol (scanner gate)
-------------------------------------------
1. Print control, marked-a and marked-b at ACTUAL SIZE / 100%, targeting $PRINT_PPI ppi. Disable fit-to-page and photo auto-enhancement if possible.
2. Use the same printer, paper and settings for all three files. Record printer model, paper type and print quality.
3. Scan each print at $PRINT_PPI dpi with geometric/photo enhancement disabled if possible. Save losslessly as PNG.
4. Crop each scan tightly to the printed artwork edges. Do NOT manually perspective-correct, sharpen or resize it.
5. Save the three captures exactly as listed in acquisition-plan.tsv inside the acquisition directory.
6. Run make v4-physical-qualification. A valid marked result requires exact HMAC-authenticated payload recovery; the control must reject.

The first gate is deliberately scanner-first. Smartphone photographs should be acquired only after this controlled print/scan gate, so printer/channel failure can be separated from free-camera geometry failure.
EOF2

echo "created Build35 physical pack in $OUTPUT_DIR"
echo "canonical carrier: ${canonical_w}x${canonical_h} px; print at $PRINT_PPI ppi -> ${width_cm} x ${height_cm} cm"
echo "fixture manifest: $manifest"
echo "acquisition plan: $plan"
