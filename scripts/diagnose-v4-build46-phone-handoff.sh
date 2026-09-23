#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD46_DIAGNOSTIC_DIR:-v4-phone private/build46-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-600}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Blind diagnostics are always completed before any reference-assisted geometry
# is generated. This keeps the production/public-evidence observation separate
# from the optional post-hoc oracle comparison.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-handoff.json"
    echo "Build46 blind handoff diagnostic: $file"
    set +e
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-handoff \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
    rc=$?
    set -e
    if [[ $rc -eq 124 ]]; then
        echo "error: Build46 handoff diagnostic timed out for $file" >&2
        exit 1
    elif [[ $rc -ne 0 ]]; then
        echo "error: Build46 handoff diagnostic failed for $file (rc=$rc)" >&2
        exit "$rc"
    fi
done

oracle_available=false
if [[ -f "$REFERENCE" ]]; then
    if python3 - <<'PY' >/dev/null 2>&1
import cv2
assert hasattr(cv2, 'SIFT_create')
PY
    then
        oracle_available=true
        for file in phone-b-mild.jpg phone-b-angle.jpg; do
            input="$ACQUISITION_DIR/$file"
            stem="${file%.*}"
            python3 ./scripts/build45-reference-register.py \
                --reference "$REFERENCE" --acquired "$input" \
                --out "$OUTPUT_DIR/oracle/$stem-quad.json"
        done
    else
        echo "Build46 note: Python OpenCV/SIFT unavailable; oracle-distance columns will remain blank." >&2
    fi
else
    echo "Build46 note: marked-B reference not found; oracle-distance columns will remain blank: $REFERENCE" >&2
fi

tsv="$OUTPUT_DIR/build46-qualified-handoff.tsv"
md="$OUTPUT_DIR/build46-qualified-handoff.md"
printf 'image\tclassification\tfrozen\tqualified\tdirect_required\tdirect_available\tbuild42_required\tbuild42_available\tcandidate\tpair\tpair_rank\tproposal\tvalidation\tpilot_score\tpilot_margin\torigin_x\torigin_y\tsingle_profiles\tsingle_list_frames\tsingle_confidence\tsingle_hmac\toracle_mean_px\toracle_max_px\toracle_ratio\n' > "$tsv"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-handoff.json"
    quad="$OUTPUT_DIR/oracle/$stem-quad.json"
    python3 - "$json" "$quad" "$file" "$oracle_available" >> "$tsv" <<'PY'
import json, math, os, sys
j=json.load(open(sys.argv[1], encoding='utf-8'))
h=j.get('handoff') or {}
qfile=sys.argv[2]
name=sys.argv[3]
oracle_enabled=sys.argv[4].lower() == 'true' and os.path.isfile(qfile)
oracle=None
if oracle_enabled:
    qj=json.load(open(qfile, encoding='utf-8'))
    oracle=qj.get('quad') or []
for c in h.get('candidates') or []:
    mean=maxe=ratio=''
    if oracle and len(oracle)==4:
        cq=c.get('source_quad') or []
        if len(cq)==4:
            ds=[math.hypot(float(cq[i]['x'])-float(oracle[i]['x']), float(cq[i]['y'])-float(oracle[i]['y'])) for i in range(4)]
            mean=sum(ds)/4.0
            maxe=max(ds)
            xs=[float(p['x']) for p in oracle]; ys=[float(p['y']) for p in oracle]
            # Normalize to the oracle artwork diagonal, not the phone frame.
            d1=math.hypot(xs[3]-xs[0], ys[3]-ys[0])
            d2=math.hypot(xs[2]-xs[1], ys[2]-ys[1])
            diag=(d1+d2)/2.0
            ratio=mean/diag if diag else 0.0
            mean=f'{mean:.3f}'; maxe=f'{maxe:.3f}'; ratio=f'{ratio:.6f}'
    vals=[
        name,h.get('classification',''),h.get('build43_frozen_candidates',0),h.get('build43_qualified_candidates',0),
        h.get('direct_ensemble_required',0),str(h.get('direct_ensemble_available',False)).lower(),
        h.get('build42_bank_required',0),str(h.get('build42_bank_available',False)).lower(),
        c.get('index',0),c.get('source_pair',''),c.get('source_pair_rank',0),
        f"{c.get('proposal_score',0):.6f}",f"{c.get('validation_score',0):.6f}",
        f"{c.get('pilot_score',0):.6f}",f"{c.get('pilot_margin',0):.6f}",
        c.get('origin_x_blocks',0),c.get('origin_y_blocks',0),c.get('single_profiles_tried',0),
        c.get('single_list_frames_tried',0),f"{c.get('single_max_data_confidence',0):.3f}",
        str(c.get('single_hmac_authenticated',False)).lower(),mean,maxe,ratio,
    ]
    print('\t'.join(map(str,vals)))
if not (h.get('candidates') or []):
    vals=[name,h.get('classification',''),h.get('build43_frozen_candidates',0),h.get('build43_qualified_candidates',0),
          h.get('direct_ensemble_required',0),str(h.get('direct_ensemble_available',False)).lower(),
          h.get('build42_bank_required',0),str(h.get('build42_bank_available',False)).lower()] + ['']*16
    print('\t'.join(map(str,vals)))
PY
done

{
    echo '# PixSeal Build46 qualified-geometry handoff diagnostic'
    echo
    echo 'Diagnostic only. Build44 production behavior and Build43/Build42 thresholds are unchanged.'
    echo
    echo '| image | class | frozen | qualified | direct req/avail | B42 req/avail | cand | pair | rank | proposal | validation | pilot | margin | origin | profiles | list frames | conf | single HMAC | oracle mean/max/ratio |'
    echo '|---|---|---:|---:|---|---|---:|---|---:|---:|---:|---:|---:|---|---:|---:|---:|---:|---|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image class frozen qualified dreq davail breq bavail cand pair prank prop val ps pm ox oy prof lf conf shmac ome omax orat; do
        echo "| $image | $class | $frozen | $qualified | $dreq/$davail | $breq/$bavail | $cand | $pair | $prank | $prop | $val | $ps | $pm | ($ox,$oy) | $prof | $lf | $conf | $shmac | ${ome:-n/a}/${omax:-n/a}/${orat:-n/a} |"
    done
    echo
    echo 'Interpretation rules:'
    echo '- qualified-ensemble-shortfall: at least one held-out-qualified candidate authenticates alone diagnostically, but production correctly lacks the required two-geometry ensemble.'
    echo '- qualified-candidate-mismatch: a held-out-qualified singleton exists but does not authenticate even when decoded alone.'
    echo '- build42-bank-shortfall: at least two qualified candidates exist but fewer than the three required by the Build42 list bank.'
    echo
    echo 'Oracle distance is post-hoc only. The reference is generated after all blind diagnostics complete and cannot influence search, ranking or qualification.'
} > "$md"

echo "Build46 qualified-handoff diagnostic written to:"
echo "  $tsv"
echo "  $md"
cat "$tsv"
