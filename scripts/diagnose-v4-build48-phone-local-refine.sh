#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD48_DIAGNOSTIC_DIR:-v4-phone private/build48-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-1200}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind proposal-only refinement. No reference/oracle exists here.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-refine.json"
    echo "Build48 blind local-refinement diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-refine \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
done

# Phase 2: post-hoc oracle only after both blind JSON files are fixed.
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
    echo "Build48 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

detail="$OUTPUT_DIR/build48-local-refinement.tsv"
summary="$OUTPUT_DIR/build48-local-refinement-summary.tsv"
md="$OUTPUT_DIR/build48-local-refinement.md"
printf 'image\tcandidate\tseed\tpair\tpair_rank\ttier\tcell_rank\tpre_proposal\tpre_validation\tpre_qualified\tpost_proposal\tpost_validation\tpost_pilot\tpost_margin\tpost_qualified\tmovement_px\tsingle_hmac\tpre_oracle_mean_px\tpost_oracle_mean_px\tpost_oracle_max_px\tpost_oracle_ratio\toracle_improvement_px\n' > "$detail"
printf 'image\tseeds\tpre_qualified\tpost_qualified\tpost_authenticated\tpre_nearest\tpre_mean_px\tpost_nearest\tpost_mean_px\tpost_max_px\tpost_ratio\tpost_nearest_qualified\tpost_nearest_hmac\tpost_pair\tpost_pair_rank\tpost_tier\tpost_cell_rank\n' > "$summary"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    stem="${file%.*}"; json="$OUTPUT_DIR/blind/$stem-refine.json"; quad="$OUTPUT_DIR/oracle/$stem-quad.json"
    python3 - "$json" "$quad" "$file" "$oracle_available" "$detail" "$summary" <<'PY'
import json, math, os, sys
j=json.load(open(sys.argv[1],encoding='utf-8')); r=j.get('refine') or {}
cands=r.get('candidates') or []
name=sys.argv[3]; use_oracle=sys.argv[4].lower()=='true' and os.path.isfile(sys.argv[2])
detail_path, summary_path=sys.argv[5],sys.argv[6]
oracle=None
if use_oracle:
    oracle=(json.load(open(sys.argv[2],encoding='utf-8')).get('quad') or [])
def err(q):
    if not oracle or len(oracle)!=4 or len(q or [])!=4: return None
    ds=[math.hypot(float(q[i]['x'])-float(oracle[i]['x']),float(q[i]['y'])-float(oracle[i]['y'])) for i in range(4)]
    mean=sum(ds)/4.; mx=max(ds)
    d1=math.hypot(float(oracle[3]['x'])-float(oracle[0]['x']),float(oracle[3]['y'])-float(oracle[0]['y']))
    d2=math.hypot(float(oracle[2]['x'])-float(oracle[1]['x']),float(oracle[2]['y'])-float(oracle[1]['y']))
    diag=(d1+d2)/2.
    return mean,mx,(mean/diag if diag else 0.)
rows=[]; pre_best=None; post_best=None
for c in cands:
    pe=err(c.get('source_quad_before')); qe=err(c.get('source_quad_after'))
    if pe is not None and (pre_best is None or pe[0] < pre_best[1][0]): pre_best=(c,pe)
    if qe is not None and (post_best is None or qe[0] < post_best[1][0]): post_best=(c,qe)
    improvement='' if pe is None or qe is None else f'{pe[0]-qe[0]:.3f}'
    rows.append([name,c.get('index',''),c.get('seed_index',''),c.get('source_pair',''),c.get('source_pair_rank',''),c.get('source_tier',''),c.get('cell_rank',''),
        f"{c.get('pre_proposal',0):.6f}",f"{c.get('pre_validation',0):.6f}",str(c.get('pre_qualified',False)).lower(),
        f"{c.get('post_proposal',0):.6f}",f"{c.get('post_validation',0):.6f}",f"{c.get('post_pilot_score',0):.6f}",f"{c.get('post_pilot_margin',0):.6f}",str(c.get('post_qualified',False)).lower(),
        f"{c.get('refine_mean_movement_px',0):.3f}",str(c.get('single_hmac_authenticated',False)).lower(),
        '' if pe is None else f'{pe[0]:.3f}','' if qe is None else f'{qe[0]:.3f}','' if qe is None else f'{qe[1]:.3f}','' if qe is None else f'{qe[2]:.6f}',improvement])
with open(detail_path,'a',encoding='utf-8') as f:
    for row in rows: f.write('\t'.join(map(str,row))+'\n')
if post_best:
    c,e=post_best
    vals=[name,r.get('seeds_selected',len(cands)),r.get('pre_qualified',0),r.get('post_qualified',0),r.get('post_authenticated',0),
          '' if pre_best is None else pre_best[0].get('index',''),'' if pre_best is None else f'{pre_best[1][0]:.3f}',
          c.get('index',''),f'{e[0]:.3f}',f'{e[1]:.3f}',f'{e[2]:.6f}',str(c.get('post_qualified',False)).lower(),str(c.get('single_hmac_authenticated',False)).lower(),
          c.get('source_pair',''),c.get('source_pair_rank',''),c.get('source_tier',''),c.get('cell_rank','')]
else:
    vals=[name,r.get('seeds_selected',len(cands)),r.get('pre_qualified',0),r.get('post_qualified',0),r.get('post_authenticated',0),'','','','','','','','','','','','']
with open(summary_path,'a',encoding='utf-8') as f: f.write('\t'.join(map(str,vals))+'\n')
PY
done

{
    echo '# PixSeal Build48 local projective refinement diagnostic'
    echo
    echo 'Diagnostic only. Seed selection and local geometry refinement use proposal evidence only. The refined bank is frozen before held-out qualification. Oracle geometry is computed only after both blind JSON files exist.'
    echo
    echo '## Summary'
    echo
    echo '| image | seeds | pre qualified | post qualified | post HMAC | pre nearest | pre mean px | post nearest | post mean px | post max px | ratio | post nearest qualified | post nearest HMAC | pair | pair rank | tier | cell |'
    echo '|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|---|---:|---|---:|'
    tail -n +2 "$summary" | while IFS=$'\t' read -r image seeds preq postq auth pren prem postn postm postx ratio pq ph pair prank tier cell; do
      echo "| $image | $seeds | $preq | $postq | $auth | ${pren:-n/a} | ${prem:-n/a} | ${postn:-n/a} | ${postm:-n/a} | ${postx:-n/a} | ${ratio:-n/a} | ${pq:-n/a} | ${ph:-n/a} | ${pair:-n/a} | ${prank:-n/a} | ${tier:-n/a} | ${cell:-n/a} |"
    done
    echo
    echo '## Candidate detail'
    echo
    echo 'See `build48-local-refinement.tsv` for every selected seed and its pre/post-refinement metrics.'
} > "$md"

echo "Build48 local-refinement diagnostic written to:"
echo "  $detail"
echo "  $summary"
echo "  $md"
cat "$summary"
