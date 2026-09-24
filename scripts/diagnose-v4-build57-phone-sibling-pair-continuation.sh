#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD57_DIAGNOSTIC_DIR:-v4-phone private/build57-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD57_TIMEOUT:-10800}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: reproduce the complete Build56 blind bank and continue every retained
# second-pair state with bounded 1px proposal-only coordinate descent. Both
# images finish before any private oracle is generated.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-sibling-pair-continuation.json"
    echo "Build57 blind second-pair continuation diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-sibling-pair-continue \
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
    echo "Build57 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv, json, math, sys
from pathlib import Path

out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']
targets={'phone-b-mild.jpg':10,'phone-b-angle.jpg':18}
parent_path=out/'build57-second-pair-parents.tsv'
state_path=out/'build57-continuation-states.tsv'
summary_path=out/'build57-second-pair-continuation-summary.tsv'
md_path=out/'build57-second-pair-continuation.md'

def xy(p): return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
def err(q, oracle):
    if not oracle or len(q or [])!=4: return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,oracle))/4.0
def f3(v): return '' if v is None else f'{v:.3f}'
def f6(v): return '' if v is None else f'{float(v):.6f}'
def b(v): return str(bool(v)).lower()

parent_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','dimension_a','delta_a_px','dimension_b','delta_b_px','parent_proposal','parent_validation','parent_qualified','parent_hmac','parent_oracle_px','continuation_evaluations','continuation_states']
state_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','continuation_index','pass','dimension','corner','axis','delta_px','proposal','validation','qualified','hmac','movement_from_seed_px','seed_oracle_px','parent_oracle_px','continuation_oracle_px','change_from_parent_px','change_from_seed_px']
summary_fields=['image','target_candidate','pair','pair_rank','seed_rank_pair','seed_oracle_px','baseline_oracle_px','baseline_proposal','baseline_validation','baseline_qualified','baseline_hmac','second_pair_states','continuation_states','best_parent_root_index','best_parent_first_pair_rank','best_parent_state_index','best_parent_sibling_rank','best_parent_second_pair_rank','best_parent_oracle_px','best_parent_proposal','best_parent_qualified','best_cont_root_index','best_cont_first_pair_rank','best_cont_parent_state_index','best_cont_parent_sibling_rank','best_cont_second_pair_rank','best_cont_index','best_cont_dimension','best_cont_delta_px','best_cont_oracle_px','best_cont_change_from_parent_px','best_cont_proposal','best_cont_validation','best_cont_qualified','best_cont_hmac','best_qualified_cont_oracle_px','any_cont_better_both','any_cont_qualified_gain','any_cont_hmac','classification']

parents=[]; states=[]; summaries=[]
for name in images:
    stem=Path(name).stem
    rep=json.load(open(out/'blind'/f'{stem}-sibling-pair-continuation.json',encoding='utf-8'))['second_pair_continuation']
    oracle=None
    if use_oracle:
        oracle=json.load(open(out/'oracle'/f'{stem}-quad.json',encoding='utf-8')).get('quad')
    candidates=rep.get('candidates') or []
    for c in candidates:
        seed_e=err(c.get('seed_source_quad'),oracle)
        for root in c.get('roots') or []:
            for branch in root.get('branches') or []:
                for cs in branch.get('continuation_states') or []:
                    cp=cs.get('parent') or {}
                    for sib in cs.get('siblings') or []:
                        for ps in sib.get('pair_states') or []:
                            pe=err(ps.get('source_quad'),oracle)
                            conts=ps.get('continuation_states') or []
                            parents.append({
                                'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':branch.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'dimension_a':ps.get('dimension_a',''),'delta_a_px':ps.get('delta_a_px',''),'dimension_b':ps.get('dimension_b',''),'delta_b_px':ps.get('delta_b_px',''),'parent_proposal':f6(ps.get('proposal')),'parent_validation':f6(ps.get('validation')),'parent_qualified':b(ps.get('qualified')),'parent_hmac':b(ps.get('single_hmac_authenticated')),'parent_oracle_px':f3(pe),'continuation_evaluations':ps.get('continuation_evaluations',''),'continuation_states':len(conts)})
                            for st in conts:
                                se=err(st.get('source_quad'),oracle)
                                states.append({
                                    'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':branch.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'continuation_index':st.get('index',''),'pass':st.get('pass',''),'dimension':st.get('dimension',''),'corner':st.get('corner',''),'axis':st.get('axis',''),'delta_px':st.get('delta_px',''),'proposal':f6(st.get('proposal')),'validation':f6(st.get('validation')),'qualified':b(st.get('qualified')),'hmac':b(st.get('single_hmac_authenticated')),'movement_from_seed_px':f3(st.get('mean_movement_from_seed_px')),'seed_oracle_px':f3(seed_e),'parent_oracle_px':f3(pe),'continuation_oracle_px':f3(se),'change_from_parent_px':f3(None if se is None or pe is None else se-pe),'change_from_seed_px':f3(None if se is None or seed_e is None else se-seed_e),
                                    '_oracle':se,'_parent_oracle':pe,'_proposal':st.get('proposal'),'_parent_proposal':ps.get('proposal'),'_qualified':bool(st.get('qualified')),'_hmac':bool(st.get('single_hmac_authenticated')),'_root':root,'_branch':branch,'_cp':cp,'_sib':sib,'_ps':ps,'_st':st})

    target_idx=targets[name]
    target=next((c for c in candidates if c.get('index')==target_idx),None)
    if target is None:
        summaries.append({'image':name,'target_candidate':target_idx,'classification':'target-missing','any_cont_better_both':'false','any_cont_qualified_gain':'false','any_cont_hmac':'false'})
        continue
    seed_e=err(target.get('seed_source_quad'),oracle); base=target.get('baseline') or {}; base_e=err(base.get('source_quad'),oracle)
    tparents=[]; tstates=[]
    for root in target.get('roots') or []:
        for branch in root.get('branches') or []:
            for cs in branch.get('continuation_states') or []:
                cp=cs.get('parent') or {}
                for sib in cs.get('siblings') or []:
                    for ps in sib.get('pair_states') or []:
                        pe=err(ps.get('source_quad'),oracle)
                        tparents.append((pe,root,branch,cp,sib,ps))
                        for st in ps.get('continuation_states') or []:
                            se=err(st.get('source_quad'),oracle)
                            tstates.append((se,pe,root,branch,cp,sib,ps,st))
    best_parent=min((x for x in tparents if x[0] is not None),default=None,key=lambda x:x[0])
    best_cont=min((x for x in tstates if x[0] is not None),default=None,key=lambda x:x[0])
    best_q=min((x for x in tstates if x[0] is not None and bool(x[7].get('qualified'))),default=None,key=lambda x:x[0])
    better=[]; qgain=[]
    for x in tstates:
        se,pe,_,_,_,_,ps,st=x
        if se is not None and pe is not None and st.get('proposal') is not None and ps.get('proposal') is not None and st['proposal']>ps['proposal']+1e-7 and se<pe-1e-9:
            better.append(x)
            if bool(st.get('qualified')): qgain.append(x)
    any_hmac=any(bool(x[7].get('single_hmac_authenticated')) for x in tstates)
    if any_hmac: cls='second-pair-continuation-recovery'
    elif qgain: cls='second-pair-continuation-qualified-gain'
    elif better: cls='second-pair-continuation-geometric-gain'
    elif tstates: cls='second-pair-continuation-proposal-only'
    else: cls='second-pair-continuation-not-triggered'
    row={'image':name,'target_candidate':target_idx,'pair':target.get('source_pair',''),'pair_rank':target.get('source_pair_rank',''),'seed_rank_pair':target.get('seed_rank_within_pair',''),'seed_oracle_px':f3(seed_e),'baseline_oracle_px':f3(base_e),'baseline_proposal':f6(base.get('proposal')),'baseline_validation':f6(base.get('validation')),'baseline_qualified':b(base.get('qualified')),'baseline_hmac':b(base.get('single_hmac_authenticated')),'second_pair_states':len(tparents),'continuation_states':len(tstates),'any_cont_better_both':b(bool(better)),'any_cont_qualified_gain':b(bool(qgain)),'any_cont_hmac':b(any_hmac),'classification':cls}
    if best_parent:
        pe,root,branch,cp,sib,ps=best_parent
        row.update({'best_parent_root_index':root.get('root_index',''),'best_parent_first_pair_rank':branch.get('pair_rank',''),'best_parent_state_index':cp.get('index',''),'best_parent_sibling_rank':sib.get('rank',''),'best_parent_second_pair_rank':ps.get('rank',''),'best_parent_oracle_px':f3(pe),'best_parent_proposal':f6(ps.get('proposal')),'best_parent_qualified':b(ps.get('qualified'))})
    if best_cont:
        se,pe,root,branch,cp,sib,ps,st=best_cont
        row.update({'best_cont_root_index':root.get('root_index',''),'best_cont_first_pair_rank':branch.get('pair_rank',''),'best_cont_parent_state_index':cp.get('index',''),'best_cont_parent_sibling_rank':sib.get('rank',''),'best_cont_second_pair_rank':ps.get('rank',''),'best_cont_index':st.get('index',''),'best_cont_dimension':st.get('dimension',''),'best_cont_delta_px':st.get('delta_px',''),'best_cont_oracle_px':f3(se),'best_cont_change_from_parent_px':f3(None if se is None or pe is None else se-pe),'best_cont_proposal':f6(st.get('proposal')),'best_cont_validation':f6(st.get('validation')),'best_cont_qualified':b(st.get('qualified')),'best_cont_hmac':b(st.get('single_hmac_authenticated'))})
    if best_q: row['best_qualified_cont_oracle_px']=f3(best_q[0])
    summaries.append(row)

with open(parent_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=parent_fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(parents)
with open(state_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=state_fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(states)
with open(summary_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=summary_fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(summaries)
with open(md_path,'w',encoding='utf-8') as f:
    f.write('# Build57 second-pair continuation diagnostic\n\n')
    f.write('Build57 reproduces the complete Build56 blind bank and continues every retained second-pair state with bounded proposal-only 1px coordinate descent. Every accepted intermediate is retained. The complete B/mild+B/angle bank is frozen before held-out/full-pilot qualification or HMAC; SIFT/reference geometry is post-hoc only.\n\n')
    f.write('## Post-hoc summary\n\n| image | target | pair | seed px | baseline px | second pairs | continuation states | best parent px | best continuation px | continuation qualified | continuation HMAC | qualified gain | interpretation |\n|---|---:|---|---:|---:|---:|---:|---:|---:|---|---|---|---|\n')
    for r in summaries:
        f.write(f"| {r.get('image','')} | {r.get('target_candidate','')} | {r.get('pair','')} | {r.get('seed_oracle_px','n/a') or 'n/a'} | {r.get('baseline_oracle_px','n/a') or 'n/a'} | {r.get('second_pair_states',0)} | {r.get('continuation_states',0)} | {r.get('best_parent_oracle_px','n/a') or 'n/a'} | {r.get('best_cont_oracle_px','n/a') or 'n/a'} | {r.get('best_cont_qualified','n/a') or 'n/a'} | {r.get('best_cont_hmac','n/a') or 'n/a'} | {r.get('any_cont_qualified_gain','false')} | {r.get('classification','')} |\n")
    f.write('\nInterpretation is post-hoc only. `second-pair-continuation-recovery` requires HMAC authentication. `second-pair-continuation-qualified-gain` requires a frozen continuation state to improve both proposal and independent oracle error relative to its exact Build56 second-pair parent while passing unchanged qualification. No label changes production behavior.\n')

print('Build57 second-pair continuation diagnostic written to:')
print(' ',parent_path); print(' ',state_path); print(' ',summary_path); print(' ',md_path)
print('\t'.join(summary_fields))
for r in summaries:
    print('\t'.join(str(r.get(k,'')) for k in summary_fields))
PY
