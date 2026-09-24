#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD56_DIAGNOSTIC_DIR:-v4-phone private/build56-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD56_TIMEOUT:-7200}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: reproduce Build55 blind bank, then scan bounded pair escapes only
# from Build55 siblings that are proposal-local under the full +/-1px single
# coordinate stencil. Both images are completed before any oracle generation.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-sibling-pair-escape.json"
    echo "Build56 blind sibling pair-escape diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-sibling-pair-escape \
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
    echo "Build56 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv, json, math, sys
from pathlib import Path

out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']
targets={'phone-b-mild.jpg':10,'phone-b-angle.jpg':18}
parent_path=out/'build56-sibling-local.tsv'
pair_path=out/'build56-pair-states.tsv'
summary_path=out/'build56-sibling-pair-summary.tsv'
md_path=out/'build56-sibling-pair-escape.md'

def xy(p): return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
def err(q, oracle):
    if not oracle or len(q or [])!=4: return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,oracle))/4.0
def f3(v): return '' if v is None else f'{v:.3f}'
def f6(v): return '' if v is None else f'{float(v):.6f}'
def b(v): return str(bool(v)).lower()

parent_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','parent_pair_rank','continuation_state_index','sibling_rank','sibling_dimension','sibling_delta_px','parent_proposal','parent_validation','parent_qualified','parent_hmac','parent_oracle_px','single_evaluations','single_improving','coordinate_local','pair_evaluations','pair_improving','pair_states']
pair_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','parent_pair_rank','parent_state_index','parent_sibling_rank','pair_escape_rank','dimension_a','corner_a','axis_a','delta_a_px','dimension_b','corner_b','axis_b','delta_b_px','proposal','validation','qualified','hmac','movement_from_seed_px','seed_oracle_px','parent_oracle_px','pair_oracle_px','change_from_parent_px','change_from_seed_px']
summary_fields=['image','target_candidate','pair','pair_rank','seed_rank_pair','seed_oracle_px','baseline_oracle_px','baseline_proposal','baseline_validation','baseline_qualified','baseline_hmac','sibling_states','coordinate_local_siblings','sibling_pair_states','best_parent_root_index','best_parent_pair_rank','best_parent_state_index','best_parent_sibling_rank','best_parent_oracle_px','best_parent_proposal','best_parent_qualified','best_pair_root_index','best_pair_parent_pair_rank','best_pair_parent_state_index','best_pair_parent_sibling_rank','best_pair_rank','best_pair_dimension_a','best_pair_delta_a_px','best_pair_dimension_b','best_pair_delta_b_px','best_pair_oracle_px','best_pair_change_from_parent_px','best_pair_proposal','best_pair_validation','best_pair_qualified','best_pair_hmac','best_qualified_pair_root_index','best_qualified_pair_parent_pair_rank','best_qualified_pair_parent_state_index','best_qualified_pair_parent_sibling_rank','best_qualified_pair_rank','best_qualified_pair_oracle_px','any_pair_better_both','any_pair_qualified_gain','any_pair_hmac','classification']

parents=[]; pairs=[]; summaries=[]
for name in images:
    stem=Path(name).stem
    rep=json.load(open(out/'blind'/f'{stem}-sibling-pair-escape.json',encoding='utf-8'))['pair_escape']
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
                        parent=sib.get('parent') or {}
                        pe=err(parent.get('source_quad'),oracle)
                        ps=sib.get('pair_states') or []
                        parents.append({
                            'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'parent_pair_rank':branch.get('pair_rank',''),'continuation_state_index':cp.get('index',''),'sibling_rank':sib.get('rank',''),'sibling_dimension':parent.get('dimension',''),'sibling_delta_px':parent.get('delta_px',''),'parent_proposal':f6(parent.get('proposal')),'parent_validation':f6(parent.get('validation')),'parent_qualified':b(parent.get('qualified')),'parent_hmac':b(parent.get('single_hmac_authenticated')),'parent_oracle_px':f3(pe),'single_evaluations':sib.get('single_evaluations',''),'single_improving':sib.get('single_improving',''),'coordinate_local':b(sib.get('coordinate_local')),'pair_evaluations':sib.get('pair_evaluations',''),'pair_improving':sib.get('pair_improving',''),'pair_states':len(ps)})
                        for st in ps:
                            se=err(st.get('source_quad'),oracle)
                            row={
                                'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'parent_pair_rank':branch.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'pair_escape_rank':st.get('rank',''),'dimension_a':st.get('dimension_a',''),'corner_a':st.get('corner_a',''),'axis_a':st.get('axis_a',''),'delta_a_px':st.get('delta_a_px',''),'dimension_b':st.get('dimension_b',''),'corner_b':st.get('corner_b',''),'axis_b':st.get('axis_b',''),'delta_b_px':st.get('delta_b_px',''),'proposal':f6(st.get('proposal')),'validation':f6(st.get('validation')),'qualified':b(st.get('qualified')),'hmac':b(st.get('single_hmac_authenticated')),'movement_from_seed_px':f3(st.get('mean_movement_from_seed_px')),'seed_oracle_px':f3(seed_e),'parent_oracle_px':f3(pe),'pair_oracle_px':f3(se),'change_from_parent_px':f3(None if se is None or pe is None else se-pe),'change_from_seed_px':f3(None if se is None or seed_e is None else se-seed_e),
                                '_oracle':se,'_parent_oracle':pe,'_proposal':st.get('proposal'),'_parent_proposal':parent.get('proposal'),'_qualified':bool(st.get('qualified')),'_hmac':bool(st.get('single_hmac_authenticated')),'_root_index':root.get('root_index'),'_branch_rank':branch.get('pair_rank'),'_state_index':cp.get('index'),'_sib_rank':sib.get('rank'),'_pair':st,'_parent':parent}
                            pairs.append(row)

    target_idx=targets[name]
    target=next((c for c in candidates if c.get('index')==target_idx),None)
    if target is None:
        summaries.append({'image':name,'target_candidate':target_idx,'classification':'target-missing','any_pair_better_both':'false','any_pair_qualified_gain':'false','any_pair_hmac':'false'})
        continue
    seed_e=err(target.get('seed_source_quad'),oracle); base=target.get('baseline') or {}; base_e=err(base.get('source_quad'),oracle)
    tparents=[]; tpairs=[]
    for root in target.get('roots') or []:
        for branch in root.get('branches') or []:
            for cs in branch.get('continuation_states') or []:
                cp=cs.get('parent') or {}
                for sib in cs.get('siblings') or []:
                    parent=sib.get('parent') or {}; pe=err(parent.get('source_quad'),oracle)
                    tparents.append((pe,root,branch,cp,sib,parent))
                    for st in sib.get('pair_states') or []:
                        se=err(st.get('source_quad'),oracle)
                        tpairs.append((se,pe,root,branch,cp,sib,parent,st))
    local_count=sum(1 for x in tparents if bool(x[4].get('coordinate_local')))
    best_parent=min((x for x in tparents if x[0] is not None), default=None, key=lambda x:x[0])
    best_pair=min((x for x in tpairs if x[0] is not None), default=None, key=lambda x:x[0])
    best_q=min((x for x in tpairs if x[0] is not None and bool(x[7].get('qualified'))), default=None, key=lambda x:x[0])
    any_hmac=any(bool(x[7].get('single_hmac_authenticated')) for x in tpairs)
    better=[]; qgain=[]
    for x in tpairs:
        se,pe,_,_,_,_,parent,st=x
        if se is not None and pe is not None and st.get('proposal') is not None and parent.get('proposal') is not None and st['proposal']>parent['proposal']+1e-7 and se<pe-1e-9:
            better.append(x)
            if bool(st.get('qualified')): qgain.append(x)
    if any_hmac: cls='sibling-pair-recovery'
    elif qgain: cls='sibling-pair-qualified-gain'
    elif better: cls='sibling-pair-geometric-gain'
    elif tpairs: cls='sibling-pair-proposal-only'
    else: cls='sibling-pair-not-triggered'
    row={'image':name,'target_candidate':target_idx,'pair':target.get('source_pair',''),'pair_rank':target.get('source_pair_rank',''),'seed_rank_pair':target.get('seed_rank_within_pair',''),'seed_oracle_px':f3(seed_e),'baseline_oracle_px':f3(base_e),'baseline_proposal':f6(base.get('proposal')),'baseline_validation':f6(base.get('validation')),'baseline_qualified':b(base.get('qualified')),'baseline_hmac':b(base.get('single_hmac_authenticated')),'sibling_states':len(tparents),'coordinate_local_siblings':local_count,'sibling_pair_states':len(tpairs),'any_pair_better_both':b(bool(better)),'any_pair_qualified_gain':b(bool(qgain)),'any_pair_hmac':b(any_hmac),'classification':cls}
    if best_parent:
        pe,root,branch,cp,sib,parent=best_parent
        row.update({'best_parent_root_index':root.get('root_index',''),'best_parent_pair_rank':branch.get('pair_rank',''),'best_parent_state_index':cp.get('index',''),'best_parent_sibling_rank':sib.get('rank',''),'best_parent_oracle_px':f3(pe),'best_parent_proposal':f6(parent.get('proposal')),'best_parent_qualified':b(parent.get('qualified'))})
    if best_pair:
        se,pe,root,branch,cp,sib,parent,st=best_pair
        row.update({'best_pair_root_index':root.get('root_index',''),'best_pair_parent_pair_rank':branch.get('pair_rank',''),'best_pair_parent_state_index':cp.get('index',''),'best_pair_parent_sibling_rank':sib.get('rank',''),'best_pair_rank':st.get('rank',''),'best_pair_dimension_a':st.get('dimension_a',''),'best_pair_delta_a_px':st.get('delta_a_px',''),'best_pair_dimension_b':st.get('dimension_b',''),'best_pair_delta_b_px':st.get('delta_b_px',''),'best_pair_oracle_px':f3(se),'best_pair_change_from_parent_px':f3(None if se is None or pe is None else se-pe),'best_pair_proposal':f6(st.get('proposal')),'best_pair_validation':f6(st.get('validation')),'best_pair_qualified':b(st.get('qualified')),'best_pair_hmac':b(st.get('single_hmac_authenticated'))})
    if best_q:
        se,pe,root,branch,cp,sib,parent,st=best_q
        row.update({'best_qualified_pair_root_index':root.get('root_index',''),'best_qualified_pair_parent_pair_rank':branch.get('pair_rank',''),'best_qualified_pair_parent_state_index':cp.get('index',''),'best_qualified_pair_parent_sibling_rank':sib.get('rank',''),'best_qualified_pair_rank':st.get('rank',''),'best_qualified_pair_oracle_px':f3(se)})
    summaries.append(row)

with open(parent_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=parent_fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(parents)
with open(pair_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=pair_fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(pairs)
with open(summary_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=summary_fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(summaries)
with open(md_path,'w',encoding='utf-8') as f:
    f.write('# Build56 sibling pair-escape diagnostic\n\n')
    f.write('Build56 reproduces the complete Build55 blind sibling bank. Every frozen sibling is first tested against the complete independent +/-1px single-coordinate stencil using the unchanged proposal score. Only siblings with zero proposal-improving single-coordinate neighbors receive the bounded 112-state coupled +/-1px pair stencil; at most eight pair states are retained by proposal. The complete bank is frozen before held-out/full-pilot qualification or HMAC. SIFT/reference geometry is post-hoc only.\n\n')
    f.write('## Post-hoc summary\n\n| image | target | pair | seed px | baseline px | siblings | local siblings | pair states | best parent px | best pair px | pair qualified | pair HMAC | qualified gain | interpretation |\n|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---|---|---|---|\n')
    for r in summaries:
        f.write(f"| {r.get('image','')} | {r.get('target_candidate','')} | {r.get('pair','')} | {r.get('seed_oracle_px','n/a') or 'n/a'} | {r.get('baseline_oracle_px','n/a') or 'n/a'} | {r.get('sibling_states',0)} | {r.get('coordinate_local_siblings',0)} | {r.get('sibling_pair_states',0)} | {r.get('best_parent_oracle_px','n/a') or 'n/a'} | {r.get('best_pair_oracle_px','n/a') or 'n/a'} | {r.get('best_pair_qualified','n/a') or 'n/a'} | {r.get('best_pair_hmac','n/a') or 'n/a'} | {r.get('any_pair_qualified_gain','false')} | {r.get('classification','')} |\n")
    f.write('\nInterpretation is post-hoc only. `sibling-pair-recovery` requires HMAC authentication. `sibling-pair-qualified-gain` requires a frozen pair state to improve both proposal and independent oracle error relative to its exact Build55 sibling parent while passing unchanged qualification. No label changes production behavior.\n')

print('Build56 sibling pair-escape diagnostic written to:')
print(' ',parent_path); print(' ',pair_path); print(' ',summary_path); print(' ',md_path)
print('\t'.join(summary_fields))
for r in summaries:
    print('\t'.join(str(r.get(k,'')) for k in summary_fields))
PY
