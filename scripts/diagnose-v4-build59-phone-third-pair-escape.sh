#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD59_DIAGNOSTIC_DIR:-v4-phone private/build59-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD59_TIMEOUT:-18000}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-third-pair-escape.json"
    echo "Build59 blind third-pair escape diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-third-pair-escape \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
done

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
    echo "Build59 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv, json, math, sys
from pathlib import Path
out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']; targets={'phone-b-mild.jpg':10,'phone-b-angle.jpg':18}
parent_path=out/'build59-third-pair-parents.tsv'; state_path=out/'build59-third-pair-states.tsv'; summary_path=out/'build59-third-pair-summary.tsv'; md_path=out/'build59-third-pair-escape.md'

def xy(p): return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
def err(q,o):
    if not o or len(q or [])!=4: return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,o))/4.0
def f3(v): return '' if v is None else f'{v:.3f}'
def f6(v): return '' if v is None else f'{float(v):.6f}'
def b(v): return str(bool(v)).lower()

parent_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','parent_continuation_index','second_sibling_rank','parent_proposal','parent_validation','parent_qualified','parent_hmac','parent_oracle_px','single_evaluations','single_improving','coordinate_local','pair_evaluations','pair_improving','pair_states']
state_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','parent_continuation_index','second_sibling_rank','third_pair_rank','dimension_a','delta_a_px','dimension_b','delta_b_px','proposal','validation','qualified','hmac','seed_oracle_px','parent_oracle_px','pair_oracle_px','change_from_parent_px','change_from_seed_px']
summary_fields=['image','target_candidate','pair','pair_rank','seed_rank_pair','seed_oracle_px','baseline_oracle_px','baseline_proposal','baseline_validation','baseline_qualified','baseline_hmac','second_sibling_states','coordinate_local_second_siblings','third_pair_states','best_parent_oracle_px','best_parent_qualified','best_pair_root_index','best_pair_first_pair_rank','best_pair_parent_state_index','best_pair_parent_sibling_rank','best_pair_second_pair_rank','best_pair_parent_continuation_index','best_pair_parent_second_sibling_rank','best_pair_rank','best_pair_oracle_px','best_pair_change_from_parent_px','best_pair_proposal','best_pair_validation','best_pair_qualified','best_pair_hmac','best_qualified_pair_oracle_px','any_pair_better_both','any_pair_qualified_gain','any_pair_hmac','classification']
parents=[]; states=[]; summaries=[]
for name in images:
    stem=Path(name).stem; rep=json.load(open(out/'blind'/f'{stem}-third-pair-escape.json',encoding='utf-8'))['third_pair_escape']; oracle=None
    if use_oracle: oracle=json.load(open(out/'oracle'/f'{stem}-quad.json',encoding='utf-8')).get('quad')
    candidates=rep.get('candidates') or []
    flat=[]
    for c in candidates:
        seed_e=err(c.get('seed_source_quad'),oracle)
        for root in c.get('roots') or []:
          for br in root.get('branches') or []:
           for cs in br.get('continuation_states') or []:
            cp=cs.get('parent') or {}
            for sib in cs.get('siblings') or []:
             for ps in sib.get('pair_states') or []:
              for cont in ps.get('continuation_states') or []:
               for ss in cont.get('sibling_states') or []:
                par=ss.get('parent') or {}; pe=err(par.get('source_quad'),oracle); pss=ss.get('pair_states') or []
                meta=(c,root,br,cp,sib,ps,cont,ss,par,pe,seed_e)
                flat.append(meta)
                parents.append({'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':br.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'parent_continuation_index':cont.get('index',''),'second_sibling_rank':ss.get('rank',''),'parent_proposal':f6(par.get('proposal')),'parent_validation':f6(par.get('validation')),'parent_qualified':b(par.get('qualified')),'parent_hmac':b(par.get('single_hmac_authenticated')),'parent_oracle_px':f3(pe),'single_evaluations':ss.get('single_evaluations',''),'single_improving':ss.get('single_improving',''),'coordinate_local':b(ss.get('coordinate_local')),'pair_evaluations':ss.get('pair_evaluations',''),'pair_improving':ss.get('pair_improving',''),'pair_states':len(pss)})
                for tp in pss:
                    te=err(tp.get('source_quad'),oracle)
                    states.append({'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':br.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'parent_continuation_index':cont.get('index',''),'second_sibling_rank':ss.get('rank',''),'third_pair_rank':tp.get('rank',''),'dimension_a':tp.get('dimension_a',''),'delta_a_px':tp.get('delta_a_px',''),'dimension_b':tp.get('dimension_b',''),'delta_b_px':tp.get('delta_b_px',''),'proposal':f6(tp.get('proposal')),'validation':f6(tp.get('validation')),'qualified':b(tp.get('qualified')),'hmac':b(tp.get('single_hmac_authenticated')),'seed_oracle_px':f3(seed_e),'parent_oracle_px':f3(pe),'pair_oracle_px':f3(te),'change_from_parent_px':f3(None if te is None or pe is None else te-pe),'change_from_seed_px':f3(None if te is None or seed_e is None else te-seed_e)})
    target_idx=targets[name]; target=next((c for c in candidates if c.get('index')==target_idx),None)
    if target is None:
        summaries.append({'image':name,'target_candidate':target_idx,'classification':'target-missing','any_pair_better_both':'false','any_pair_qualified_gain':'false','any_pair_hmac':'false'}); continue
    seed_e=err(target.get('seed_source_quad'),oracle); base=target.get('baseline') or {}; base_e=err(base.get('source_quad'),oracle)
    tf=[x for x in flat if x[0].get('index')==target_idx]; tpairs=[]
    for x in tf:
        c,root,br,cp,sib,ps,cont,ss,par,pe,se=x
        for tp in ss.get('pair_states') or []: tpairs.append((x,tp,err(tp.get('source_quad'),oracle)))
    row={'image':name,'target_candidate':target_idx,'pair':target.get('source_pair',''),'pair_rank':target.get('source_pair_rank',''),'seed_rank_pair':target.get('seed_rank_within_pair',''),'seed_oracle_px':f3(seed_e),'baseline_oracle_px':f3(base_e),'baseline_proposal':f6(base.get('proposal')),'baseline_validation':f6(base.get('validation')),'baseline_qualified':b(base.get('qualified')),'baseline_hmac':b(base.get('single_hmac_authenticated')),'second_sibling_states':len(tf),'coordinate_local_second_siblings':sum(bool(x[7].get('coordinate_local')) for x in tf),'third_pair_states':len(tpairs)}
    best_parent=min((x for x in tf if x[9] is not None),default=None,key=lambda x:x[9]); best_pair=min((x for x in tpairs if x[2] is not None),default=None,key=lambda x:x[2]); best_q=min((x for x in tpairs if x[2] is not None and bool(x[1].get('qualified'))),default=None,key=lambda x:x[2])
    better=[]; qgain=[]
    for x,tp,te in tpairs:
        pe=x[9]; par=x[8]
        if te is not None and pe is not None and tp.get('proposal') is not None and par.get('proposal') is not None and float(tp['proposal'])>float(par['proposal'])+1e-7 and te<pe-1e-9:
            better.append((x,tp,te))
            if bool(tp.get('qualified')): qgain.append((x,tp,te))
    any_hmac=any(bool(tp.get('single_hmac_authenticated')) for _,tp,_ in tpairs)
    row.update({'any_pair_better_both':b(bool(better)),'any_pair_qualified_gain':b(bool(qgain)),'any_pair_hmac':b(any_hmac)})
    if best_parent: row.update({'best_parent_oracle_px':f3(best_parent[9]),'best_parent_qualified':b(best_parent[8].get('qualified'))})
    if best_pair:
        x,tp,te=best_pair; c,root,br,cp,sib,ps,cont,ss,par,pe,se=x
        row.update({'best_pair_root_index':root.get('root_index',''),'best_pair_first_pair_rank':br.get('pair_rank',''),'best_pair_parent_state_index':cp.get('index',''),'best_pair_parent_sibling_rank':sib.get('rank',''),'best_pair_second_pair_rank':ps.get('rank',''),'best_pair_parent_continuation_index':cont.get('index',''),'best_pair_parent_second_sibling_rank':ss.get('rank',''),'best_pair_rank':tp.get('rank',''),'best_pair_oracle_px':f3(te),'best_pair_change_from_parent_px':f3(None if te is None or pe is None else te-pe),'best_pair_proposal':f6(tp.get('proposal')),'best_pair_validation':f6(tp.get('validation')),'best_pair_qualified':b(tp.get('qualified')),'best_pair_hmac':b(tp.get('single_hmac_authenticated'))})
    if best_q: row['best_qualified_pair_oracle_px']=f3(best_q[2])
    row['classification']='third-pair-recovery' if any_hmac else 'third-pair-qualified-gain' if qgain else 'third-pair-geometric-gain' if better else 'third-pair-proposal-only' if tpairs else 'third-pair-not-triggered'
    summaries.append(row)
for path,fields,rows in [(parent_path,parent_fields,parents),(state_path,state_fields,states),(summary_path,summary_fields,summaries)]:
    with open(path,'w',newline='',encoding='utf-8') as f: w=csv.DictWriter(f,fieldnames=fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(rows)
with open(md_path,'w',encoding='utf-8') as f:
    f.write('# Build59 third-pair escape diagnostic\n\nBuild59 reproduces the complete Build58 blind bank, probes every frozen second-pair sibling for one-coordinate +/-1px locality using proposal only, and applies the bounded coupled pair stencil only at proposal-local siblings. The entire B/mild+B/angle bank is frozen before qualification/HMAC; SIFT/reference geometry is post-hoc only.\n\n')
    f.write('## Post-hoc summary\n\n| image | target | pair | seed px | baseline px | sibling parents | local parents | third-pair states | best parent px | best pair px | pair qualified | pair HMAC | qualified gain | interpretation |\n|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---|---|---|---|\n')
    for r in summaries: f.write(f"| {r.get('image','')} | {r.get('target_candidate','')} | {r.get('pair','')} | {r.get('seed_oracle_px','n/a') or 'n/a'} | {r.get('baseline_oracle_px','n/a') or 'n/a'} | {r.get('second_sibling_states',0)} | {r.get('coordinate_local_second_siblings',0)} | {r.get('third_pair_states',0)} | {r.get('best_parent_oracle_px','n/a') or 'n/a'} | {r.get('best_pair_oracle_px','n/a') or 'n/a'} | {r.get('best_pair_qualified','n/a') or 'n/a'} | {r.get('best_pair_hmac','n/a') or 'n/a'} | {r.get('any_pair_qualified_gain','false')} | {r.get('classification','')} |\n")
    f.write('\nInterpretation is post-hoc only. `third-pair-recovery` requires HMAC authentication. `third-pair-qualified-gain` requires a frozen pair escape to improve both proposal and independent oracle error relative to its exact Build58 sibling parent while passing unchanged qualification. No label changes production behavior.\n')
print('Build59 third-pair escape diagnostic written to:')
for p in [parent_path,state_path,summary_path,md_path]: print(' ',p)
print('\t'.join(summary_fields))
for r in summaries: print('\t'.join(str(r.get(k,'')) for k in summary_fields))
PY
