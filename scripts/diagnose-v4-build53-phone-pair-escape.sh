#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD53_DIAGNOSTIC_DIR:-v4-phone private/build53-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD53_TIMEOUT:-3600}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind proposal-only root/locality/pair geometry. Pair scanning is
# triggered only by a 1px coordinate-local proposal maximum. All retained pair
# states are proposal-ranked and frozen before qualification/HMAC.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-pair-escape.json"
    echo "Build53 blind coupled pair-escape diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-pair-escape \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
done

# Phase 2: private oracle only after both blind JSON files are complete.
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
    echo "Build53 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv, json, math, os, sys
from pathlib import Path

out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']
root_path=out/'build53-roots.tsv'
pair_path=out/'build53-pair-states.tsv'
summary_path=out/'build53-pair-escape-summary.tsv'
md_path=out/'build53-pair-escape.md'

def xy(p): return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
def err(q, oracle):
    if not oracle or len(q or [])!=4: return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,oracle))/4.0

def f3(v): return '' if v is None else f'{v:.3f}'
def f6(v): return '' if v is None else f'{float(v):.6f}'
def b(v): return str(bool(v)).lower()

root_fields=['image','candidate','seed_rank_pair','pair','pair_rank','tier','cell_rank','root_index','root_step_px','root_pass','root_dim','root_corner','root_axis','root_delta_px','root_proposal','root_validation','root_qualified','root_hmac','single_evaluations','single_improving','coordinate_local','pair_evaluations','pair_improving','pair_retained','movement_from_seed_px','seed_oracle_px','root_oracle_px','root_change_from_seed_px']
pair_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','pair_rank_within_root','dim_a','corner_a','axis_a','delta_a_px','dim_b','corner_b','axis_b','delta_b_px','proposal','validation','qualified','hmac','movement_from_seed_px','seed_oracle_px','root_oracle_px','pair_oracle_px','pair_change_from_root_px','pair_change_from_seed_px']
summary_fields=['image','target_candidate','pair','pair_rank','seed_rank_pair','seed_oracle_px','baseline_oracle_px','baseline_change_px','baseline_proposal','baseline_validation','baseline_qualified','baseline_hmac','coordinate_local_roots','pair_states','best_local_root_index','best_local_root_oracle_px','best_local_root_qualified','best_pair_root_index','best_pair_rank','best_pair_oracle_px','best_pair_change_from_root_px','best_pair_proposal','best_pair_validation','best_pair_qualified','best_pair_hmac','best_qualified_pair_root_index','best_qualified_pair_rank','best_qualified_pair_oracle_px','any_pair_better_both','any_pair_qualified_gain','any_pair_hmac','classification']

root_rows=[]; pair_rows=[]; summary_rows=[]
for name in images:
    stem=name.rsplit('.',1)[0]
    j=json.load(open(out/'blind'/f'{stem}-pair-escape.json',encoding='utf-8'))
    r=j.get('pair_escape') or {}; cands=r.get('candidates') or []
    oracle=None
    op=out/'oracle'/f'{stem}-quad.json'
    if use_oracle and op.is_file(): oracle=(json.load(open(op,encoding='utf-8')).get('quad') or [])
    for c in cands:
        c['_seed_err']=err(c.get('seed_source_quad'),oracle)
        base=c.get('baseline') or {}; base['_err']=err(base.get('source_quad'),oracle)
        for root in c.get('roots') or []:
            rs=root.get('root') or {}; root['_err']=err(rs.get('source_quad'),oracle)
            for ps in root.get('pair_states') or []: ps['_err']=err(ps.get('source_quad'),oracle)
            root_rows.append({
                'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'tier':c.get('source_tier',''),'cell_rank':c.get('cell_rank',''),
                'root_index':root.get('root_index',''),'root_step_px':rs.get('step_size_px',''),'root_pass':rs.get('pass',''),'root_dim':rs.get('dimension',''),'root_corner':rs.get('corner',''),'root_axis':rs.get('axis',''),'root_delta_px':rs.get('delta_px',''),
                'root_proposal':f6(rs.get('proposal')),'root_validation':f6(rs.get('validation')),'root_qualified':b(rs.get('qualified')),'root_hmac':b(rs.get('single_hmac_authenticated')),
                'single_evaluations':root.get('single_evaluations',''),'single_improving':root.get('single_improving',''),'coordinate_local':b(root.get('coordinate_local')),'pair_evaluations':root.get('pair_evaluations',''),'pair_improving':root.get('pair_improving',''),'pair_retained':root.get('pair_retained',''),
                'movement_from_seed_px':f3(rs.get('mean_movement_from_seed_px')),'seed_oracle_px':f3(c['_seed_err']),'root_oracle_px':f3(root['_err']),'root_change_from_seed_px':f3(None if root['_err'] is None or c['_seed_err'] is None else root['_err']-c['_seed_err'])})
            for ps in root.get('pair_states') or []:
                pair_rows.append({
                    'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'pair_rank_within_root':ps.get('rank',''),
                    'dim_a':ps.get('dimension_a',''),'corner_a':ps.get('corner_a',''),'axis_a':ps.get('axis_a',''),'delta_a_px':ps.get('delta_a_px',''),'dim_b':ps.get('dimension_b',''),'corner_b':ps.get('corner_b',''),'axis_b':ps.get('axis_b',''),'delta_b_px':ps.get('delta_b_px',''),
                    'proposal':f6(ps.get('proposal')),'validation':f6(ps.get('validation')),'qualified':b(ps.get('qualified')),'hmac':b(ps.get('single_hmac_authenticated')),'movement_from_seed_px':f3(ps.get('mean_movement_from_seed_px')),
                    'seed_oracle_px':f3(c['_seed_err']),'root_oracle_px':f3(root['_err']),'pair_oracle_px':f3(ps['_err']),'pair_change_from_root_px':f3(None if ps['_err'] is None or root['_err'] is None else ps['_err']-root['_err']),'pair_change_from_seed_px':f3(None if ps['_err'] is None or c['_seed_err'] is None else ps['_err']-c['_seed_err'])})
    if not oracle or not cands:
        summary_rows.append({'image':name,'classification':'oracle-unavailable'})
        continue
    target=min((c for c in cands if c.get('_seed_err') is not None),key=lambda c:c['_seed_err'])
    locals_=[rt for rt in (target.get('roots') or []) if rt.get('coordinate_local')]
    pairs=[(rt,ps) for rt in locals_ for ps in (rt.get('pair_states') or [])]
    bestroot=min(locals_,key=lambda rt:rt['_err']) if locals_ else None
    bestpair=min(pairs,key=lambda rp:rp[1]['_err']) if pairs else None
    qpairs=[rp for rp in pairs if rp[1].get('qualified')]
    bestq=min(qpairs,key=lambda rp:rp[1]['_err']) if qpairs else None
    better=[rp for rp in pairs if rp[1]['_err'] < rp[0]['_err']-1e-9 and float(rp[1].get('proposal',0)) > float((rp[0].get('root') or {}).get('proposal',0))+1e-7]
    qgain=[rp for rp in better if rp[1].get('qualified')]
    anyh=any(ps.get('single_hmac_authenticated') for _,ps in pairs)
    if anyh: cls='pair-recovery'
    elif qgain: cls='pair-qualified-gain'
    elif better: cls='pair-geometric-gain'
    elif locals_ and not pairs: cls='pair-no-proposal-escape'
    elif locals_: cls='pair-no-oracle-gain'
    else: cls='pair-not-triggered'
    base=target.get('baseline') or {}
    row={'image':name,'target_candidate':target.get('index',''),'pair':target.get('source_pair',''),'pair_rank':target.get('source_pair_rank',''),'seed_rank_pair':target.get('seed_rank_within_pair',''),'seed_oracle_px':f3(target['_seed_err']),'baseline_oracle_px':f3(base.get('_err')),'baseline_change_px':f3(None if base.get('_err') is None else base['_err']-target['_seed_err']),'baseline_proposal':f6(base.get('proposal')),'baseline_validation':f6(base.get('validation')),'baseline_qualified':b(base.get('qualified')),'baseline_hmac':b(base.get('single_hmac_authenticated')),'coordinate_local_roots':len(locals_),'pair_states':len(pairs),'any_pair_better_both':b(better),'any_pair_qualified_gain':b(qgain),'any_pair_hmac':b(anyh),'classification':cls}
    if bestroot:
        row.update({'best_local_root_index':bestroot.get('root_index',''),'best_local_root_oracle_px':f3(bestroot['_err']),'best_local_root_qualified':b((bestroot.get('root') or {}).get('qualified'))})
    if bestpair:
        rt,ps=bestpair; row.update({'best_pair_root_index':rt.get('root_index',''),'best_pair_rank':ps.get('rank',''),'best_pair_oracle_px':f3(ps['_err']),'best_pair_change_from_root_px':f3(ps['_err']-rt['_err']),'best_pair_proposal':f6(ps.get('proposal')),'best_pair_validation':f6(ps.get('validation')),'best_pair_qualified':b(ps.get('qualified')),'best_pair_hmac':b(ps.get('single_hmac_authenticated'))})
    if bestq:
        rt,ps=bestq; row.update({'best_qualified_pair_root_index':rt.get('root_index',''),'best_qualified_pair_rank':ps.get('rank',''),'best_qualified_pair_oracle_px':f3(ps['_err'])})
    summary_rows.append(row)

for path,fields,rows in [(root_path,root_fields,root_rows),(pair_path,pair_fields,pair_rows),(summary_path,summary_fields,summary_rows)]:
    with open(path,'w',encoding='utf-8',newline='') as f:
        w=csv.DictWriter(f,fieldnames=fields,delimiter='\t',extrasaction='ignore',lineterminator='\n'); w.writeheader(); w.writerows(rows)

with open(md_path,'w',encoding='utf-8') as f:
    f.write('# PixSeal Build53 coupled pair-escape diagnostic\n\n')
    f.write('Research-only. Build53 retains the unchanged top4-per-side-pair seed bank and proposal score. It creates a 2px accepted-state root path, probes every root with the complete +/-1px single-coordinate stencil, and scans coupled +/-1px two-coordinate moves only at roots with zero single-coordinate proposal improvements. Up to eight pair states per root are retained by proposal rank. The complete geometry bank is frozen before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry is generated only after both blind JSON files exist.\n\n')
    f.write('## Post-hoc summary\n\n')
    f.write('| image | target | pair | seed px | baseline px | local roots | pair states | best local root px | best pair px | pair change | pair qualified | pair HMAC | qualified gain | interpretation |\n')
    f.write('|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---|---|---|---|\n')
    for r in summary_rows:
        f.write(f"| {r.get('image','')} | {r.get('target_candidate','n/a')} | {r.get('pair','n/a')} | {r.get('seed_oracle_px','n/a')} | {r.get('baseline_oracle_px','n/a')} | {r.get('coordinate_local_roots','n/a')} | {r.get('pair_states','n/a')} | {r.get('best_local_root_oracle_px','n/a')} | {r.get('best_pair_oracle_px','n/a')} | {r.get('best_pair_change_from_root_px','n/a')} | {r.get('best_pair_qualified','n/a')} | {r.get('best_pair_hmac','n/a')} | {r.get('any_pair_qualified_gain','n/a')} | {r.get('classification','')} |\n")
    f.write('\nInterpretation is post-hoc only. `pair-recovery` requires HMAC authentication. `pair-qualified-gain` requires a retained pair state to improve proposal and independent oracle error relative to its coordinate-local root while also passing the unchanged held-out/full-pilot qualification. `pair-geometric-gain` improves proposal and oracle but does not qualify. `pair-no-proposal-escape` means no coupled move beats the root proposal; `pair-not-triggered` means the target candidate has no 1px coordinate-local 2px root. No label changes production behavior.\n')
PY

echo "Build53 pair-escape diagnostic written to:"
echo "  $OUTPUT_DIR/build53-roots.tsv"
echo "  $OUTPUT_DIR/build53-pair-states.tsv"
echo "  $OUTPUT_DIR/build53-pair-escape-summary.tsv"
echo "  $OUTPUT_DIR/build53-pair-escape.md"
cat "$OUTPUT_DIR/build53-pair-escape-summary.tsv"
