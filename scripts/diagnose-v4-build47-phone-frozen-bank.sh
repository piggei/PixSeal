#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD47_DIAGNOSTIC_DIR:-v4-phone private/build47-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-900}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind extended bank. No reference/oracle exists in this phase.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-frozen.json"
    echo "Build47 blind frozen-bank diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-frozen \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
done

# Phase 2: optional post-hoc reference geometry, only after blind JSON is fixed.
oracle_available=false
if [[ -f "$REFERENCE" ]] && python3 - <<'PY' >/dev/null 2>&1
import cv2
assert hasattr(cv2, 'SIFT_create')
PY
then
    oracle_available=true
    for file in phone-b-mild.jpg phone-b-angle.jpg; do
        input="$ACQUISITION_DIR/$file"; stem="${file%.*}"
        python3 ./scripts/build45-reference-register.py --reference "$REFERENCE" --acquired "$input" --out "$OUTPUT_DIR/oracle/$stem-quad.json"
    done
else
    echo "Build47 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

tsv="$OUTPUT_DIR/build47-frozen-bank.tsv"
md="$OUTPUT_DIR/build47-frozen-bank.md"
printf 'image\tlimit\tstage\tavailable\tqualified\tproduction_best\tproduction_best_validation\toracle_nearest\toracle_mean_px\toracle_max_px\toracle_ratio\toracle_nearest_qualified\toracle_nearest_pair\toracle_nearest_pair_rank\toracle_nearest_tier\toracle_nearest_cell_rank\n' > "$tsv"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    stem="${file%.*}"; json="$OUTPUT_DIR/blind/$stem-frozen.json"; quad="$OUTPUT_DIR/oracle/$stem-quad.json"
    python3 - "$json" "$quad" "$file" "$oracle_available" >> "$tsv" <<'PY'
import json, math, os, sys
j=json.load(open(sys.argv[1],encoding='utf-8')); r=j.get('frozen') or {}
cands=r.get('candidates') or []
name=sys.argv[3]; use_oracle=sys.argv[4].lower()=='true' and os.path.isfile(sys.argv[2])
oracle=None
if use_oracle:
    oracle=(json.load(open(sys.argv[2],encoding='utf-8')).get('quad') or [])
def err(c):
    if not oracle or len(oracle)!=4: return None
    q=c.get('source_quad') or []
    if len(q)!=4: return None
    ds=[math.hypot(float(q[i]['x'])-float(oracle[i]['x']),float(q[i]['y'])-float(oracle[i]['y'])) for i in range(4)]
    mean=sum(ds)/4.; mx=max(ds)
    d1=math.hypot(float(oracle[3]['x'])-float(oracle[0]['x']),float(oracle[3]['y'])-float(oracle[0]['y']))
    d2=math.hypot(float(oracle[2]['x'])-float(oracle[1]['x']),float(oracle[2]['y'])-float(oracle[1]['y']))
    diag=(d1+d2)/2.; return mean,mx,(mean/diag if diag else 0.)
stages=((32,1,'production'),(64,2,'selected-pair-depth'),(128,3,'all-pair-extension'))
for limit,max_tier,stage in stages:
    sub=[c for c in cands if int(c.get('diagnostic_tier',1)) <= max_tier]
    q=sum(1 for c in sub if c.get('qualified'))
    best=max(sub,key=lambda c:c.get('validation_score',float('-inf'))) if sub else None
    nearest=None; ne=None
    if oracle:
        for c in sub:
            e=err(c)
            if e is not None and (ne is None or e[0]<ne[0]): nearest,ne=c,e
    vals=[name,limit,stage,len(sub),q,
          '' if best is None else best.get('index',''), '' if best is None else f"{best.get('validation_score',0):.6f}",
          '' if nearest is None else nearest.get('index',''), '' if ne is None else f'{ne[0]:.3f}', '' if ne is None else f'{ne[1]:.3f}', '' if ne is None else f'{ne[2]:.6f}',
          '' if nearest is None else str(nearest.get('qualified',False)).lower(), '' if nearest is None else nearest.get('source_pair',''), '' if nearest is None else nearest.get('source_pair_rank',''),
          '' if nearest is None else nearest.get('diagnostic_tier_name',''), '' if nearest is None else nearest.get('cell_rank','')]
    print('\t'.join(map(str,vals)))
PY
done

{
    echo '# PixSeal Build47 frozen-candidate bank diagnostic'
    echo
    echo 'Diagnostic only. Production Build43 remains unchanged. The 32 stage is the exact production proposal tier, even when fewer than 32 candidates are available.'
    echo
    echo '| image | cap | stage | available | qualified | production-best | validation | oracle-nearest | mean px | max px | ratio | nearest qualified | pair | pair rank | nearest tier | cell rank |'
    echo '|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|---:|---|---:|'
    tail -n +2 "$tsv" | while IFS=$'\t' read -r image cap stage avail qual pb pbv on om ox orat oq op opr ot oc; do
      echo "| $image | $cap | $stage | $avail | $qual | $pb | $pbv | ${on:-n/a} | ${om:-n/a} | ${ox:-n/a} | ${orat:-n/a} | ${oq:-n/a} | ${op:-n/a} | ${opr:-n/a} | ${ot:-n/a} | ${oc:-n/a} |"
    done
    echo
    echo 'Tier 1 is the exact Build43 production proposal structure. Tier 2 adds cell ranks 1 and 2 on the same two selected side pairs. Tier 3 adds side-pair ranks 3 through 6 using production cell ranks 0 and 3.'
    echo
    echo 'The oracle is computed only after both blind JSON files are complete. It never participates in candidate generation, ranking or qualification.'
} > "$md"

echo "Build47 frozen-bank diagnostic written to:"
echo "  $tsv"
echo "  $md"
cat "$tsv"
