#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD52_DIAGNOSTIC_DIR:-v4-phone private/build52-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD52_TIMEOUT:-3600}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind optimizer-only geometry. The unchanged top-4 proposal seeds,
# baseline coarse-to-fine endpoint and fine-first accepted-state bank are all
# created before any reference/oracle geometry exists. The key is used only
# after the complete geometry bank is frozen and a state passes qualification.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-restart.json"
    echo "Build52 blind fine-restart optimizer diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-restart \
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
    echo "Build52 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

detail="$OUTPUT_DIR/build52-optimizer-states.tsv"
summary="$OUTPUT_DIR/build52-optimizer-summary.tsv"
md="$OUTPUT_DIR/build52-optimizer.md"
printf 'image\tcandidate\tseed_rank_pair\tpair\tpair_rank\ttier\tcell_rank\tstate_kind\tstate_index\tstep_px\tpass\tdim\tcorner\taxis\tdelta_px\tproposal\tvalidation\tpilot\tmargin\tqualified\thmac\tmovement_from_seed_px\toracle_mean_px\toracle_change_from_seed_px\n' > "$detail"
printf 'image\ttarget_candidate\tpair\tpair_rank\tseed_rank_pair\tseed_oracle_px\tbaseline_oracle_px\tbaseline_change_px\tbaseline_proposal\tbaseline_validation\tbaseline_qualified\tbaseline_hmac\tfine_states\tfine_best_index\tfine_best_oracle_px\tfine_best_change_px\tfine_best_proposal\tfine_best_validation\tfine_best_qualified\tfine_best_hmac\tfine_best_qualified_index\tfine_best_qualified_oracle_px\tany_fine_better_both\tany_fine_hmac\tclassification\n' > "$summary"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    stem="${file%.*}"; json="$OUTPUT_DIR/blind/$stem-restart.json"; quad="$OUTPUT_DIR/oracle/$stem-quad.json"
    python3 - "$json" "$quad" "$file" "$oracle_available" "$detail" "$summary" <<'PY'
import json, math, os, sys

j=json.load(open(sys.argv[1],encoding='utf-8')); r=j.get('restart') or {}
cands=r.get('candidates') or []
name=sys.argv[3]; use_oracle=sys.argv[4].lower()=='true' and os.path.isfile(sys.argv[2])
detail_path, summary_path=sys.argv[5],sys.argv[6]
oracle=None
if use_oracle:
    oracle=(json.load(open(sys.argv[2],encoding='utf-8')).get('quad') or [])

def xy(p):
    return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))

def err(q):
    if not oracle or len(oracle)!=4 or len(q or [])!=4: return None
    ds=[]
    for i in range(4):
        qx,qy=xy(q[i]); ox,oy=xy(oracle[i]); ds.append(math.hypot(qx-ox,qy-oy))
    return sum(ds)/4.0

for c in cands:
    c['_seed_err']=err(c.get('seed_source_quad'))
    b=c.get('baseline') or {}
    b['_err']=err(b.get('source_quad'))
    for s in c.get('fine_restart') or []:
        s['_err']=err(s.get('source_quad'))

with open(detail_path,'a',encoding='utf-8') as f:
    for c in cands:
        seed_err=c['_seed_err']
        b=c.get('baseline') or {}
        rows=[('baseline',b)] + [('fine',s) for s in (c.get('fine_restart') or [])]
        for kind,s in rows:
            e=s.get('_err')
            vals=[name,c.get('index',''),c.get('seed_rank_within_pair',''),c.get('source_pair',''),c.get('source_pair_rank',''),c.get('source_tier',''),c.get('cell_rank',''),
                  kind,s.get('index',''),s.get('step_size_px',''),s.get('pass',''),s.get('dimension',''),s.get('corner',''),s.get('axis',''),s.get('delta_px',''),
                  f"{float(s.get('proposal',0)):.6f}",f"{float(s.get('validation',0)):.6f}",f"{float(s.get('pilot_score',0)):.6f}",f"{float(s.get('pilot_margin',0)):.6f}",
                  str(bool(s.get('qualified'))).lower(),str(bool(s.get('single_hmac_authenticated'))).lower(),f"{float(s.get('mean_movement_from_seed_px',0)):.3f}",
                  '' if e is None else f'{e:.3f}','' if e is None or seed_err is None else f'{e-seed_err:.3f}']
            f.write('\t'.join(map(str,vals))+'\n')

# Post-hoc focus remains the seed nearest the independently supplied oracle.
target=min((c for c in cands if c['_seed_err'] is not None), key=lambda c:c['_seed_err'], default=None)
if target is None:
    vals=[name,'','','','','','','','','','','','','','','','','','','','','','','','oracle-unavailable']
else:
    seed_err=target['_seed_err']
    base=target.get('baseline') or {}; base_err=base.get('_err')
    fine=target.get('fine_restart') or []
    valid=[s for s in fine if s.get('_err') is not None]
    best=min(valid,key=lambda s:s['_err'],default=None)
    qualified=[s for s in valid if s.get('qualified')]
    bestq=min(qualified,key=lambda s:s['_err'],default=None)
    better_both=[s for s in valid if s['_err'] < seed_err-1e-9 and float(s.get('proposal',0)) > float(target.get('seed_proposal',0))+1e-7]
    any_hmac=any(bool(s.get('single_hmac_authenticated')) for s in fine)
    if any_hmac:
        classification='optimizer-recovery'
    elif better_both:
        classification='optimizer-partial-gain'
    else:
        classification='optimizer-no-gain'
    vals=[name,target.get('index',''),target.get('source_pair',''),target.get('source_pair_rank',''),target.get('seed_rank_within_pair',''),
          f'{seed_err:.3f}', '' if base_err is None else f'{base_err:.3f}', '' if base_err is None else f'{base_err-seed_err:.3f}',
          f"{float(base.get('proposal',0)):.6f}",f"{float(base.get('validation',0)):.6f}",str(bool(base.get('qualified'))).lower(),str(bool(base.get('single_hmac_authenticated'))).lower(),
          len(fine), '' if best is None else best.get('index',''), '' if best is None else f"{best['_err']:.3f}", '' if best is None else f"{best['_err']-seed_err:.3f}",
          '' if best is None else f"{float(best.get('proposal',0)):.6f}", '' if best is None else f"{float(best.get('validation',0)):.6f}", '' if best is None else str(bool(best.get('qualified'))).lower(), '' if best is None else str(bool(best.get('single_hmac_authenticated'))).lower(),
          '' if bestq is None else bestq.get('index',''), '' if bestq is None else f"{bestq['_err']:.3f}",str(bool(better_both)).lower(),str(any_hmac).lower(),classification]
with open(summary_path,'a',encoding='utf-8') as f:
    f.write('\t'.join(map(str,vals))+'\n')
PY
done

{
    echo '# PixSeal Build52 fine-restart optimizer diagnostic'
    echo
    echo 'Research-only. The unchanged Build50/51 top-4 seeds are reused. For each seed, Build52 records the unchanged Build41 coarse-to-fine endpoint and a separate proposal-only fine restart with the fixed 2px -> 1px schedule. The fine bank retains the untouched seed plus every accepted intermediate state. All geometry is frozen before held-out/full-pilot qualification and diagnostic HMAC. SIFT/reference oracle geometry is generated only after both blind JSON files exist.'
    echo
    echo '## Post-hoc summary'
    echo
    echo '| image | target | pair | pair rank | seed rank | seed px | baseline px | baseline change | baseline proposal | baseline validation | baseline qualified | baseline HMAC | fine states | fine best | fine best px | fine change | fine proposal | fine validation | fine qualified | fine HMAC | best qualified | best qualified px | better-both | any fine HMAC | interpretation |'
    echo '|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---|---|---:|---:|---:|---:|---:|---:|---|---|---:|---:|---|---|---|'
    tail -n +2 "$summary" | while IFS=$'\t' read -r image target pair prank srank seed base bchg bp bv bq bh n best beste bestchg bestp bestv bestq besth qb qbe both anyh class; do
      echo "| $image | ${target:-n/a} | ${pair:-n/a} | ${prank:-n/a} | ${srank:-n/a} | ${seed:-n/a} | ${base:-n/a} | ${bchg:-n/a} | ${bp:-n/a} | ${bv:-n/a} | ${bq:-n/a} | ${bh:-n/a} | ${n:-n/a} | ${best:-n/a} | ${beste:-n/a} | ${bestchg:-n/a} | ${bestp:-n/a} | ${bestv:-n/a} | ${bestq:-n/a} | ${besth:-n/a} | ${qb:-n/a} | ${qbe:-n/a} | ${both:-n/a} | ${anyh:-n/a} | ${class:-n/a} |"
    done
    echo
    echo 'Interpretation labels are post-hoc laboratory summaries only. `optimizer-recovery` requires a fine-restart state to authenticate. `optimizer-partial-gain` means at least one frozen fine state improves both proposal and independent oracle error versus the untouched seed but no state authenticates. `optimizer-no-gain` means the fine restart does not preserve a better-both state. None of these labels changes production behavior.'
    echo
    echo 'See `build52-optimizer-states.tsv` for every baseline/fine state.'
} > "$md"

echo "Build52 optimizer diagnostic written to:"
echo "  $detail"
echo "  $summary"
echo "  $md"
cat "$summary"
