#!/usr/bin/env bash

# Shared helpers for PixSeal's local ImageMagick-based test suites.
# The scripts normally ask ImageMagick for dimensions. Very large PNG files can
# exceed ImageMagick resource-policy limits even with -ping; in that case read
# the PNG IHDR width/height directly without decoding pixel data.
read_image_dimensions() {
    local image="$1" output width height
    if output="$("${identify_tool[@]}" -ping -format '%w %h\n' "$image" 2>/dev/null)"; then
        read -r width height <<< "$output"
        if [[ "$width" =~ ^[0-9]+$ && "$height" =~ ^[0-9]+$ ]]; then
            printf '%s %s\n' "$width" "$height"
            return 0
        fi
    fi

    # PNG signature + first IHDR chunk. Never trust the extension alone: camera
    # files can be misnamed. Use decimal bytes for the dimensions so this path
    # is independent of shell hexadecimal parsing quirks on older hosts.
    local sig dims=()
    sig="$(od -An -v -tx1 -N8 -- "$image" 2>/dev/null | tr -d ' \n')"
    if [[ "$sig" == "89504e470d0a1a0a" ]]; then
        read -r -a dims <<< "$(od -An -v -tu1 -j16 -N8 -- "$image" 2>/dev/null)"
        if (( ${#dims[@]} == 8 )); then
            width=$(( dims[0] * 16777216 + dims[1] * 65536 + dims[2] * 256 + dims[3] ))
            height=$(( dims[4] * 16777216 + dims[5] * 65536 + dims[6] * 256 + dims[7] ))
            if (( width > 0 && height > 0 )); then
                printf '%s %s\n' "$width" "$height"
                return 0
            fi
        fi
    fi
    return 1
}


load_active_corpus_images() {
    local pics_dir="$1"
    local manifest="${ACTIVE_CORPUS_MANIFEST:-private-corpus-active.tsv}"
    images=()
    if [[ ! -f "$manifest" ]]; then
        echo "error: active corpus manifest not found: $manifest" >&2
        return 1
    fi
    while IFS=$'\t' read -r role filename width height sha extra; do
        [[ -z "$role" || "$role" == \#* ]] && continue
        if [[ -n "${extra:-}" || -z "$filename" || -z "$sha" ]]; then
            echo "error: invalid active corpus manifest row for role $role" >&2
            return 1
        fi
        local path="$pics_dir/$filename"
        if [[ ! -f "$path" ]]; then
            echo "error: active corpus file missing: $path" >&2
            return 1
        fi
        images+=("$path")
    done < "$manifest"
    if (( ${#images[@]} == 0 )); then
        echo "error: active corpus manifest contains no images" >&2
        return 1
    fi
}

verify_active_corpus_manifest() {
    local pics_dir="$1"
    local manifest="${ACTIVE_CORPUS_MANIFEST:-private-corpus-active.tsv}"
    local failures=0
    while IFS=$'\t' read -r role filename expected_w expected_h expected_sha extra; do
        [[ -z "$role" || "$role" == \#* ]] && continue
        local path="$pics_dir/$filename"
        if [[ ! -f "$path" ]]; then
            echo "FAIL $role missing $path"
            failures=$((failures+1)); continue
        fi
        local got_sha got_w got_h
        got_sha="$(sha256sum "$path" | awk '{print $1}')"
        read -r got_w got_h < <(read_image_dimensions "$path") || true
        if [[ "$got_sha" != "$expected_sha" || "$got_w" != "$expected_w" || "$got_h" != "$expected_h" ]]; then
            echo "FAIL $role $filename expected=${expected_w}x${expected_h}/$expected_sha got=${got_w}x${got_h}/$got_sha"
            failures=$((failures+1))
        else
            echo "OK   $role $filename ${got_w}x${got_h} sha256=$got_sha"
        fi
    done < "$manifest"
    (( failures == 0 ))
}
