#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD50_DIAGNOSTIC_DIR:-v4-phone private/build50-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-1200}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind top-4 proposal selection and proposal-only refinement. The
# reference image/oracle is not used by this command. The secret key is read
# only after geometry freeze and held-out qualification for diagnostic HMAC.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-refine4.json"
    echo "Build50 blind top-4 local-refinement diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-refine4 \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
done

# Phase 2: post-hoc oracle only after both blind outputs are fixed.
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
    echo "Build50 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

detail="$OUTPUT_DIR/build50-top4-refinement.tsv"
summary="$OUTPUT_DIR/build50-top4-refinement-summary.tsv"
md="$OUTPUT_DIR/build50-top4-refinement.md"
printf 'image\tcandidate\tseed\tseed_rank_pair\tpair\tpair_rank\ttier\tcell_rank\tpre_proposal\tpre_validation\tpre_qualified\tpost_proposal\tpost_validation\tpost_pilot\tpost_margin\tpost_qualified\tmovement_px\tsingle_hmac\tpre_oracle_mean_px\tpost_oracle_mean_px\tpost_oracle_max_px\tpost_oracle_ratio\toracle_improvement_px\n' > "$detail"
printf 'image\tselection\tseeds\tpre_qualified\tpost_qualified\tpost_authenticated\tpre_nearest\tpre_mean_px\tpost_nearest\tpost_mean_px\tpost_max_px\tpost_ratio\tpost_nearest_qualified\tpost_nearest_hmac\tpost_pair\tpost_pair_rank\tpost_tier\tpost_cell_rank\tpost_seed_rank_pair\n' > "$summary"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    stem="${file%.*}"; json="$OUTPUT_DIR/blind/$stem-refine4.json"; quad="$OUTPUT_DIR/oracle/$stem-quad.json"
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
    def xy(p):
        return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
    ds=[]
    for i in range(4):
        qx,qy=xy(q[i]); ox,oy=xy(oracle[i]); ds.append(math.hypot(qx-ox,qy-oy))
    mean=sum(ds)/4.; mx=max(ds)
    o0=xy(oracle[0]); o1=xy(oracle[1]); o2=xy(oracle[2]); o3=xy(oracle[3])
    d1=math.hypot(o3[0]-o0[0],o3[1]-o0[1]); d2=math.hypot(o2[0]-o1[0],o2[1]-o1[1]); diag=(d1+d2)/2.
    return mean,mx,(mean/diag if diag else 0.)
for c in cands:
    c['_pre']=err(c.get('source_quad_before')); c['_post']=err(c.get('source_quad_after'))
with open(detail_path,'a',encoding='utf-8') as f:
    for c in cands:
        pe,qe=c['_pre'],c['_post']
        improvement='' if pe is None or qe is None else f'{pe[0]-qe[0]:.3f}'
        row=[name,c.get('index',''),c.get('seed_index',''),c.get('seed_rank_within_pair',''),c.get('source_pair',''),c.get('source_pair_rank',''),c.get('source_tier',''),c.get('cell_rank',''),
             f"{c.get('pre_proposal',0):.6f}",f"{c.get('pre_validation',0):.6f}",str(c.get('pre_qualified',False)).lower(),
             f"{c.get('post_proposal',0):.6f}",f"{c.get('post_validation',0):.6f}",f"{c.get('post_pilot_score',0):.6f}",f"{c.get('post_pilot_margin',0):.6f}",str(c.get('post_qualified',False)).lower(),
             f"{c.get('refine_mean_movement_px',0):.3f}",str(c.get('single_hmac_authenticated',False)).lower(),
             '' if pe is None else f'{pe[0]:.3f}','' if qe is None else f'{qe[0]:.3f}','' if qe is None else f'{qe[1]:.3f}','' if qe is None else f'{qe[2]:.6f}',improvement]
        f.write('\t'.join(map(str,row))+'\n')

def emit(selection, subset):
    pre_best=min((c for c in subset if c['_pre'] is not None),key=lambda c:c['_pre'][0],default=None)
    post_best=min((c for c in subset if c['_post'] is not None),key=lambda c:c['_post'][0],default=None)
    if post_best:
        e=post_best['_post']
        vals=[name,selection,len(subset),sum(bool(c.get('pre_qualified')) for c in subset),sum(bool(c.get('post_qualified')) for c in subset),sum(bool(c.get('single_hmac_authenticated')) for c in subset),
              '' if pre_best is None else pre_best.get('index',''),'' if pre_best is None else f"{pre_best['_pre'][0]:.3f}",
              post_best.get('index',''),f'{e[0]:.3f}',f'{e[1]:.3f}',f'{e[2]:.6f}',str(post_best.get('post_qualified',False)).lower(),str(post_best.get('single_hmac_authenticated',False)).lower(),
              post_best.get('source_pair',''),post_best.get('source_pair_rank',''),post_best.get('source_tier',''),post_best.get('cell_rank',''),post_best.get('seed_rank_within_pair','')]
    else:
        vals=[name,selection,len(subset),sum(bool(c.get('pre_qualified')) for c in subset),sum(bool(c.get('post_qualified')) for c in subset),sum(bool(c.get('single_hmac_authenticated')) for c in subset),'','','','','','','','','','','','','']
    with open(summary_path,'a',encoding='utf-8') as f:f.write('\t'.join(map(str,vals))+'\n')

# Same top-4 blind run, compared as nested prefixes. This avoids rerunning the
# geometry/refinement and makes the Build48-like top2 vs Build50 top4 contrast exact.
emit('top2', [c for c in cands if int(c.get('seed_rank_within_pair') or 999) <= 2])
emit('top4', [c for c in cands if int(c.get('seed_rank_within_pair') or 999) <= 4])
PY
done

{
    echo '# PixSeal Build50 top-4-per-pair local refinement diagnostic'
    echo
    echo 'Diagnostic only. The unchanged Build49/raw proposal score selects up to four seeds per side-pair. All local projective refinement is proposal-only; the complete top-4 refined bank is frozen before held-out qualification. The secret key is used only for post-qualification diagnostic HMAC. Oracle geometry is computed only after both blind JSON files exist.'
    echo
    echo 'The top2 row is the nested rank<=2 subset of the same blind top4 run, so top2/top4 are directly comparable without a second geometry search.'
    echo
    echo '## Summary'
    echo
    echo '| image | selection | seeds | pre qualified | post qualified | post HMAC | pre nearest | pre mean px | post nearest | post mean px | post max px | ratio | nearest qualified | nearest HMAC | pair | pair rank | tier | cell | seed rank/pair |'
    echo '|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|---|---:|---|---:|---:|'
    tail -n +2 "$summary" | while IFS=$'\t' read -r image sel seeds preq postq auth pren prem postn postm postx ratio pq ph pair prank tier cell srank; do
      echo "| $image | $sel | $seeds | $preq | $postq | $auth | ${pren:-n/a} | ${prem:-n/a} | ${postn:-n/a} | ${postm:-n/a} | ${postx:-n/a} | ${ratio:-n/a} | ${pq:-n/a} | ${ph:-n/a} | ${pair:-n/a} | ${prank:-n/a} | ${tier:-n/a} | ${cell:-n/a} | ${srank:-n/a} |"
    done
    echo
    echo '## Candidate detail'
    echo
    echo 'See `build50-top4-refinement.tsv` for every selected seed, including its proposal rank within the side-pair and pre/post-refinement oracle error.'
} > "$md"

echo "Build50 top-4 refinement diagnostic written to:"
echo "  $detail"
echo "  $summary"
echo "  $md"
cat "$summary"
