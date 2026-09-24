#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD58_DIAGNOSTIC_DIR:-v4-phone private/build58-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD58_TIMEOUT:-14400}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: reproduce the complete Build57 blind bank, then evaluate every
# independent +/-1px sibling from every retained post-second-pair continuation
# state. Both image banks are frozen before any oracle is generated.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-second-pair-sibling-stencil.json"
    echo "Build58 blind second-pair sibling-stencil diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-second-pair-sibling \
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
    echo "Build58 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv, json, math, sys
from pathlib import Path

out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']
targets={'phone-b-mild.jpg':10,'phone-b-angle.jpg':18}
parent_path=out/'build58-second-pair-sibling-parents.tsv'
state_path=out/'build58-second-pair-sibling-states.tsv'
summary_path=out/'build58-second-pair-sibling-summary.tsv'
md_path=out/'build58-second-pair-sibling-stencil.md'

def xy(p): return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
def err(q, oracle):
    if not oracle or len(q or [])!=4: return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,oracle))/4.0
def f3(v): return '' if v is None else f'{v:.3f}'
def f6(v): return '' if v is None else f'{float(v):.6f}'
def b(v): return str(bool(v)).lower()

parent_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','continuation_index','continuation_pass','continuation_dimension','parent_proposal','parent_validation','parent_qualified','parent_hmac','parent_oracle_px','sibling_evaluations','sibling_states']
state_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','parent_continuation_index','sibling_rank','dimension','corner','axis','delta_px','proposal','validation','qualified','hmac','movement_from_seed_px','seed_oracle_px','parent_oracle_px','sibling_oracle_px','change_from_parent_px','change_from_seed_px']
summary_fields=['image','target_candidate','pair','pair_rank','seed_rank_pair','seed_oracle_px','baseline_oracle_px','baseline_proposal','baseline_validation','baseline_qualified','baseline_hmac','continuation_states','sibling_states','best_parent_root_index','best_parent_first_pair_rank','best_parent_state_index','best_parent_sibling_rank','best_parent_second_pair_rank','best_parent_continuation_index','best_parent_oracle_px','best_parent_qualified','best_sibling_root_index','best_sibling_first_pair_rank','best_sibling_parent_state_index','best_sibling_parent_sibling_rank','best_sibling_second_pair_rank','best_sibling_parent_continuation_index','best_sibling_rank','best_sibling_dimension','best_sibling_delta_px','best_sibling_oracle_px','best_sibling_change_from_parent_px','best_sibling_proposal','best_sibling_validation','best_sibling_qualified','best_sibling_hmac','best_qualified_sibling_oracle_px','any_sibling_better_both','any_sibling_qualified_gain','any_sibling_hmac','classification']

parents=[]; states=[]; summaries=[]
for name in images:
    stem=Path(name).stem
    rep=json.load(open(out/'blind'/f'{stem}-second-pair-sibling-stencil.json',encoding='utf-8'))['second_pair_sibling_stencil']
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
                            for parent in ps.get('continuation_states') or []:
                                pe=err(parent.get('source_quad'),oracle); ss=parent.get('sibling_states') or []
                                parents.append({'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':branch.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'continuation_index':parent.get('index',''),'continuation_pass':parent.get('pass',''),'continuation_dimension':parent.get('dimension',''),'parent_proposal':f6(parent.get('proposal')),'parent_validation':f6(parent.get('validation')),'parent_qualified':b(parent.get('qualified')),'parent_hmac':b(parent.get('single_hmac_authenticated')),'parent_oracle_px':f3(pe),'sibling_evaluations':parent.get('sibling_evaluations',''),'sibling_states':len(ss)})
                                for st in ss:
                                    se=err(st.get('source_quad'),oracle)
                                    states.append({'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':branch.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'parent_continuation_index':parent.get('index',''),'sibling_rank':st.get('rank',''),'dimension':st.get('dimension',''),'corner':st.get('corner',''),'axis':st.get('axis',''),'delta_px':st.get('delta_px',''),'proposal':f6(st.get('proposal')),'validation':f6(st.get('validation')),'qualified':b(st.get('qualified')),'hmac':b(st.get('single_hmac_authenticated')),'movement_from_seed_px':f3(st.get('mean_movement_from_seed_px')),'seed_oracle_px':f3(seed_e),'parent_oracle_px':f3(pe),'sibling_oracle_px':f3(se),'change_from_parent_px':f3(None if se is None or pe is None else se-pe),'change_from_seed_px':f3(None if se is None or seed_e is None else se-seed_e)})

    target_idx=targets[name]
    target=next((c for c in candidates if c.get('index')==target_idx),None)
    if target is None:
        summaries.append({'image':name,'target_candidate':target_idx,'classification':'target-missing','any_sibling_better_both':'false','any_sibling_qualified_gain':'false','any_sibling_hmac':'false'}); continue
    seed_e=err(target.get('seed_source_quad'),oracle); base=target.get('baseline') or {}; base_e=err(base.get('source_quad'),oracle)
    tparents=[]; tstates=[]
    for root in target.get('roots') or []:
        for branch in root.get('branches') or []:
            for cs in branch.get('continuation_states') or []:
                cp=cs.get('parent') or {}
                for sib in cs.get('siblings') or []:
                    for ps in sib.get('pair_states') or []:
                        for parent in ps.get('continuation_states') or []:
                            pe=err(parent.get('source_quad'),oracle); tparents.append((root,branch,cp,sib,ps,parent,pe))
                            for st in parent.get('sibling_states') or []:
                                se=err(st.get('source_quad'),oracle); tstates.append((root,branch,cp,sib,ps,parent,pe,st,se))
    row={'image':name,'target_candidate':target_idx,'pair':target.get('source_pair',''),'pair_rank':target.get('source_pair_rank',''),'seed_rank_pair':target.get('seed_rank_within_pair',''),'seed_oracle_px':f3(seed_e),'baseline_oracle_px':f3(base_e),'baseline_proposal':f6(base.get('proposal')),'baseline_validation':f6(base.get('validation')),'baseline_qualified':b(base.get('qualified')),'baseline_hmac':b(base.get('single_hmac_authenticated')),'continuation_states':len(tparents),'sibling_states':len(tstates)}
    best_parent=min((x for x in tparents if x[6] is not None),default=None,key=lambda x:x[6])
    best_sib=min((x for x in tstates if x[8] is not None),default=None,key=lambda x:x[8])
    best_q=min((x for x in tstates if x[8] is not None and bool(x[7].get('qualified'))),default=None,key=lambda x:x[8])
    better=[]; qgain=[]
    for x in tstates:
        pe=x[6]; st=x[7]; se=x[8]; parent=x[5]
        if se is not None and pe is not None and st.get('proposal') is not None and parent.get('proposal') is not None and float(st['proposal'])>float(parent['proposal'])+1e-7 and se<pe-1e-9:
            better.append(x)
            if bool(st.get('qualified')): qgain.append(x)
    any_hmac=any(bool(x[7].get('single_hmac_authenticated')) for x in tstates)
    row.update({'any_sibling_better_both':b(bool(better)),'any_sibling_qualified_gain':b(bool(qgain)),'any_sibling_hmac':b(any_hmac)})
    if best_parent:
        root,br,cp,sib,ps,parent,pe=best_parent
        row.update({'best_parent_root_index':root.get('root_index',''),'best_parent_first_pair_rank':br.get('pair_rank',''),'best_parent_state_index':cp.get('index',''),'best_parent_sibling_rank':sib.get('rank',''),'best_parent_second_pair_rank':ps.get('rank',''),'best_parent_continuation_index':parent.get('index',''),'best_parent_oracle_px':f3(pe),'best_parent_qualified':b(parent.get('qualified'))})
    if best_sib:
        root,br,cp,sib,ps,parent,pe,st,se=best_sib
        row.update({'best_sibling_root_index':root.get('root_index',''),'best_sibling_first_pair_rank':br.get('pair_rank',''),'best_sibling_parent_state_index':cp.get('index',''),'best_sibling_parent_sibling_rank':sib.get('rank',''),'best_sibling_second_pair_rank':ps.get('rank',''),'best_sibling_parent_continuation_index':parent.get('index',''),'best_sibling_rank':st.get('rank',''),'best_sibling_dimension':st.get('dimension',''),'best_sibling_delta_px':st.get('delta_px',''),'best_sibling_oracle_px':f3(se),'best_sibling_change_from_parent_px':f3(None if se is None or pe is None else se-pe),'best_sibling_proposal':f6(st.get('proposal')),'best_sibling_validation':f6(st.get('validation')),'best_sibling_qualified':b(st.get('qualified')),'best_sibling_hmac':b(st.get('single_hmac_authenticated'))})
    if best_q: row['best_qualified_sibling_oracle_px']=f3(best_q[8])
    if any_hmac: cls='second-pair-sibling-recovery'
    elif qgain: cls='second-pair-sibling-qualified-gain'
    elif better: cls='second-pair-sibling-geometric-gain'
    elif tstates: cls='second-pair-sibling-proposal-only'
    else: cls='second-pair-sibling-not-triggered'
    row['classification']=cls; summaries.append(row)

for path,fields,rows in [(parent_path,parent_fields,parents),(state_path,state_fields,states),(summary_path,summary_fields,summaries)]:
    with open(path,'w',newline='',encoding='utf-8') as f:
        w=csv.DictWriter(f,fieldnames=fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(rows)
with open(md_path,'w',encoding='utf-8') as f:
    f.write('# Build58 second-pair continuation sibling-stencil diagnostic\n\n')
    f.write('Build58 reproduces the complete Build57 blind bank and evaluates every independent +/-1px coordinate sibling from each retained post-second-pair continuation state. Each sibling starts from the identical frozen parent. The complete B/mild+B/angle bank is frozen before held-out/full-pilot qualification or HMAC; SIFT/reference geometry is post-hoc only.\n\n')
    f.write('## Post-hoc summary\n\n| image | target | pair | seed px | baseline px | continuation parents | sibling states | best parent px | best sibling px | sibling qualified | sibling HMAC | qualified gain | interpretation |\n|---|---:|---|---:|---:|---:|---:|---:|---:|---|---|---|---|\n')
    for r in summaries:
        f.write(f"| {r.get('image','')} | {r.get('target_candidate','')} | {r.get('pair','')} | {r.get('seed_oracle_px','n/a') or 'n/a'} | {r.get('baseline_oracle_px','n/a') or 'n/a'} | {r.get('continuation_states',0)} | {r.get('sibling_states',0)} | {r.get('best_parent_oracle_px','n/a') or 'n/a'} | {r.get('best_sibling_oracle_px','n/a') or 'n/a'} | {r.get('best_sibling_qualified','n/a') or 'n/a'} | {r.get('best_sibling_hmac','n/a') or 'n/a'} | {r.get('any_sibling_qualified_gain','false')} | {r.get('classification','')} |\n")
    f.write('\nInterpretation is post-hoc only. `second-pair-sibling-recovery` requires HMAC authentication. `second-pair-sibling-qualified-gain` requires a frozen sibling to improve both proposal and independent oracle error relative to its exact Build57 continuation parent while passing unchanged qualification. No label changes production behavior.\n')

print('Build58 second-pair sibling-stencil diagnostic written to:')
for p in [parent_path,state_path,summary_path,md_path]: print(' ',p)
print('\t'.join(summary_fields))
for r in summaries: print('\t'.join(str(r.get(k,'')) for k in summary_fields))
PY
