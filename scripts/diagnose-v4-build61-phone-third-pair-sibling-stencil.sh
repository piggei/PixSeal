#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD61_DIAGNOSTIC_DIR:-v4-phone private/build61-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD61_TIMEOUT:-43200}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: reproduce the complete Build60 blind bank and evaluate every
# independent +/-1px sibling from every retained post-third-pair continuation
# state. Both image banks finish before any private oracle is generated.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-third-pair-sibling-stencil.json"
    echo "Build61 blind third-pair sibling-stencil diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-third-pair-sibling \
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
    echo "Build61 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv, json, math, sys
from pathlib import Path
out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']; targets={'phone-b-mild.jpg':10,'phone-b-angle.jpg':18}
parent_path=out/'build61-third-pair-sibling-parents.tsv'
state_path=out/'build61-third-pair-sibling-states.tsv'
summary_path=out/'build61-third-pair-sibling-summary.tsv'
md_path=out/'build61-third-pair-sibling-stencil.md'

def xy(p): return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
def err(q,o):
    if not o or len(q or [])!=4: return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,o))/4.0
def f3(v): return '' if v is None else f'{v:.3f}'
def f6(v): return '' if v is None else f'{float(v):.6f}'
def b(v): return str(bool(v)).lower()
def rate(n,d): return '' if not d else f'{float(n)/float(d):.6f}'

parent_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','parent_continuation_index','second_sibling_rank','third_pair_rank','third_pair_continuation_index','continuation_pass','continuation_dimension','parent_proposal','parent_validation','parent_qualified','parent_hmac','parent_oracle_px','sibling_evaluations','sibling_states']
state_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','first_pair_rank','parent_state_index','parent_sibling_rank','second_pair_rank','parent_continuation_index','second_sibling_rank','third_pair_rank','third_pair_continuation_index','sibling_rank','dimension','corner','axis','delta_px','proposal','validation','qualified','hmac','movement_from_seed_px','seed_oracle_px','parent_oracle_px','sibling_oracle_px','change_from_parent_px','change_from_seed_px']
summary_fields=['image','target_candidate','pair','pair_rank','seed_rank_pair','seed_oracle_px','baseline_oracle_px','baseline_proposal','baseline_validation','baseline_qualified','baseline_hmac','continuation_states','sibling_states','bank_sibling_states','bank_sibling_qualified','bank_sibling_hmac','bank_sibling_qualified_rate','best_parent_root_index','best_parent_first_pair_rank','best_parent_state_index','best_parent_sibling_rank','best_parent_second_pair_rank','best_parent_continuation_index','best_parent_second_sibling_rank','best_parent_third_pair_rank','best_parent_third_pair_continuation_index','best_parent_oracle_px','best_parent_proposal','best_parent_qualified','best_sibling_root_index','best_sibling_first_pair_rank','best_sibling_parent_state_index','best_sibling_parent_sibling_rank','best_sibling_second_pair_rank','best_sibling_parent_continuation_index','best_sibling_parent_second_sibling_rank','best_sibling_third_pair_rank','best_sibling_parent_third_pair_continuation_index','best_sibling_rank','best_sibling_dimension','best_sibling_delta_px','best_sibling_oracle_px','best_sibling_change_from_parent_px','best_sibling_proposal','best_sibling_validation','best_sibling_qualified','best_sibling_hmac','best_qualified_sibling_oracle_px','any_sibling_better_both','any_sibling_qualified_gain','any_sibling_hmac','classification']
parents=[]; states=[]; summaries=[]
for name in images:
    stem=Path(name).stem
    rep=json.load(open(out/'blind'/f'{stem}-third-pair-sibling-stencil.json',encoding='utf-8'))['third_pair_sibling_stencil']
    oracle=None
    if use_oracle: oracle=json.load(open(out/'oracle'/f'{stem}-quad.json',encoding='utf-8')).get('quad')
    candidates=rep.get('candidates') or []
    all_sibling_states=0; all_sibling_qualified=0; all_sibling_hmac=0
    flat_parents=[]; flat_states=[]
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
                for tp in ss.get('pair_states') or []:
                 for parent in tp.get('continuation_states') or []:
                    pe=err(parent.get('source_quad'),oracle); sibs3=parent.get('sibling_states') or []
                    meta=(c,root,br,cp,sib,ps,cont,ss,tp,parent,pe,seed_e)
                    flat_parents.append(meta)
                    all_sibling_states += len(sibs3)
                    all_sibling_qualified += sum(1 for x in sibs3 if bool(x.get('qualified')))
                    all_sibling_hmac += sum(1 for x in sibs3 if bool(x.get('single_hmac_authenticated')))
                    parents.append({'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':br.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'parent_continuation_index':cont.get('index',''),'second_sibling_rank':ss.get('rank',''),'third_pair_rank':tp.get('rank',''),'third_pair_continuation_index':parent.get('index',''),'continuation_pass':parent.get('pass',''),'continuation_dimension':parent.get('dimension',''),'parent_proposal':f6(parent.get('proposal')),'parent_validation':f6(parent.get('validation')),'parent_qualified':b(parent.get('qualified')),'parent_hmac':b(parent.get('single_hmac_authenticated')),'parent_oracle_px':f3(pe),'sibling_evaluations':parent.get('sibling_evaluations',''),'sibling_states':len(sibs3)})
                    for st in sibs3:
                        se=err(st.get('source_quad'),oracle)
                        x=(meta,st,se); flat_states.append(x)
                        states.append({'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'first_pair_rank':br.get('pair_rank',''),'parent_state_index':cp.get('index',''),'parent_sibling_rank':sib.get('rank',''),'second_pair_rank':ps.get('rank',''),'parent_continuation_index':cont.get('index',''),'second_sibling_rank':ss.get('rank',''),'third_pair_rank':tp.get('rank',''),'third_pair_continuation_index':parent.get('index',''),'sibling_rank':st.get('rank',''),'dimension':st.get('dimension',''),'corner':st.get('corner',''),'axis':st.get('axis',''),'delta_px':st.get('delta_px',''),'proposal':f6(st.get('proposal')),'validation':f6(st.get('validation')),'qualified':b(st.get('qualified')),'hmac':b(st.get('single_hmac_authenticated')),'movement_from_seed_px':f3(st.get('mean_movement_from_seed_px')),'seed_oracle_px':f3(seed_e),'parent_oracle_px':f3(pe),'sibling_oracle_px':f3(se),'change_from_parent_px':f3(None if se is None or pe is None else se-pe),'change_from_seed_px':f3(None if se is None or seed_e is None else se-seed_e)})
    target_idx=targets[name]; target=next((c for c in candidates if c.get('index')==target_idx),None)
    if target is None:
        summaries.append({'image':name,'target_candidate':target_idx,'bank_sibling_states':all_sibling_states,'bank_sibling_qualified':all_sibling_qualified,'bank_sibling_hmac':all_sibling_hmac,'bank_sibling_qualified_rate':rate(all_sibling_qualified,all_sibling_states),'classification':'target-missing','any_sibling_better_both':'false','any_sibling_qualified_gain':'false','any_sibling_hmac':'false'}); continue
    seed_e=err(target.get('seed_source_quad'),oracle); base=target.get('baseline') or {}; base_e=err(base.get('source_quad'),oracle)
    tp=[x for x in flat_parents if x[0].get('index')==target_idx]
    ts=[x for x in flat_states if x[0][0].get('index')==target_idx]
    row={'image':name,'target_candidate':target_idx,'pair':target.get('source_pair',''),'pair_rank':target.get('source_pair_rank',''),'seed_rank_pair':target.get('seed_rank_within_pair',''),'seed_oracle_px':f3(seed_e),'baseline_oracle_px':f3(base_e),'baseline_proposal':f6(base.get('proposal')),'baseline_validation':f6(base.get('validation')),'baseline_qualified':b(base.get('qualified')),'baseline_hmac':b(base.get('single_hmac_authenticated')),'continuation_states':len(tp),'sibling_states':len(ts),'bank_sibling_states':all_sibling_states,'bank_sibling_qualified':all_sibling_qualified,'bank_sibling_hmac':all_sibling_hmac,'bank_sibling_qualified_rate':rate(all_sibling_qualified,all_sibling_states)}
    best_parent=min((x for x in tp if x[10] is not None),default=None,key=lambda x:x[10])
    best_sib=min((x for x in ts if x[2] is not None),default=None,key=lambda x:x[2])
    best_q=min((x for x in ts if x[2] is not None and bool(x[1].get('qualified'))),default=None,key=lambda x:x[2])
    better=[]; qgain=[]
    for meta,st,se in ts:
        pe=meta[10]; parent=meta[9]
        if se is not None and pe is not None and st.get('proposal') is not None and parent.get('proposal') is not None and float(st['proposal'])>float(parent['proposal'])+1e-7 and se<pe-1e-9:
            better.append((meta,st,se))
            if bool(st.get('qualified')): qgain.append((meta,st,se))
    any_hmac=any(bool(st.get('single_hmac_authenticated')) for _,st,_ in ts)
    row.update({'any_sibling_better_both':b(bool(better)),'any_sibling_qualified_gain':b(bool(qgain)),'any_sibling_hmac':b(any_hmac)})
    if best_parent:
        c,root,br,cp,sib,ps,cont,ss,tp3,parent,pe,seed=best_parent
        row.update({'best_parent_root_index':root.get('root_index',''),'best_parent_first_pair_rank':br.get('pair_rank',''),'best_parent_state_index':cp.get('index',''),'best_parent_sibling_rank':sib.get('rank',''),'best_parent_second_pair_rank':ps.get('rank',''),'best_parent_continuation_index':cont.get('index',''),'best_parent_second_sibling_rank':ss.get('rank',''),'best_parent_third_pair_rank':tp3.get('rank',''),'best_parent_third_pair_continuation_index':parent.get('index',''),'best_parent_oracle_px':f3(pe),'best_parent_proposal':f6(parent.get('proposal')),'best_parent_qualified':b(parent.get('qualified'))})
    if best_sib:
        meta,st,se=best_sib; c,root,br,cp,sib,ps,cont,ss,tp3,parent,pe,seed=meta
        row.update({'best_sibling_root_index':root.get('root_index',''),'best_sibling_first_pair_rank':br.get('pair_rank',''),'best_sibling_parent_state_index':cp.get('index',''),'best_sibling_parent_sibling_rank':sib.get('rank',''),'best_sibling_second_pair_rank':ps.get('rank',''),'best_sibling_parent_continuation_index':cont.get('index',''),'best_sibling_parent_second_sibling_rank':ss.get('rank',''),'best_sibling_third_pair_rank':tp3.get('rank',''),'best_sibling_parent_third_pair_continuation_index':parent.get('index',''),'best_sibling_rank':st.get('rank',''),'best_sibling_dimension':st.get('dimension',''),'best_sibling_delta_px':st.get('delta_px',''),'best_sibling_oracle_px':f3(se),'best_sibling_change_from_parent_px':f3(None if se is None or pe is None else se-pe),'best_sibling_proposal':f6(st.get('proposal')),'best_sibling_validation':f6(st.get('validation')),'best_sibling_qualified':b(st.get('qualified')),'best_sibling_hmac':b(st.get('single_hmac_authenticated'))})
    if best_q: row['best_qualified_sibling_oracle_px']=f3(best_q[2])
    row['classification']='third-pair-sibling-recovery' if any_hmac else 'third-pair-sibling-qualified-gain' if qgain else 'third-pair-sibling-geometric-gain' if better else 'third-pair-sibling-proposal-only' if ts else 'third-pair-sibling-not-triggered'
    summaries.append(row)
for path,fields,rows in [(parent_path,parent_fields,parents),(state_path,state_fields,states),(summary_path,summary_fields,summaries)]:
    with open(path,'w',newline='',encoding='utf-8') as f: w=csv.DictWriter(f,fieldnames=fields,delimiter='\t',extrasaction='ignore'); w.writeheader(); w.writerows(rows)
with open(md_path,'w',encoding='utf-8') as f:
    f.write('# Build61 third-pair continuation sibling-stencil diagnostic\n\nBuild61 reproduces the complete Build60 blind bank and evaluates every independent +/-1px coordinate sibling from each retained post-third-pair continuation state. Each sibling starts from the identical frozen parent. The complete B/mild+B/angle bank is frozen before qualification/HMAC; SIFT/reference geometry is post-hoc only.\n\n')
    f.write('## Post-hoc summary\n\n| image | target | seed px | continuation parents | sibling states | bank qualified/total | bank HMAC | best parent px | best sibling px | sibling qualified | sibling HMAC | qualified gain | interpretation |\n|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|---|---|\n')
    for r in summaries:
        f.write(f"| {r.get('image','')} | {r.get('target_candidate','')} | {r.get('seed_oracle_px','n/a') or 'n/a'} | {r.get('continuation_states',0)} | {r.get('sibling_states',0)} | {r.get('bank_sibling_qualified',0)}/{r.get('bank_sibling_states',0)} | {r.get('bank_sibling_hmac',0)} | {r.get('best_parent_oracle_px','n/a') or 'n/a'} | {r.get('best_sibling_oracle_px','n/a') or 'n/a'} | {r.get('best_sibling_qualified','n/a') or 'n/a'} | {r.get('best_sibling_hmac','n/a') or 'n/a'} | {r.get('any_sibling_qualified_gain','false')} | {r.get('classification','')} |\n")
    f.write('\nInterpretation is post-hoc only. `third-pair-sibling-recovery` requires HMAC authentication. `third-pair-sibling-qualified-gain` requires a frozen sibling to improve both proposal and independent oracle error relative to its exact Build60 continuation parent while passing unchanged qualification. `bank qualified/total` is tracked explicitly because Build60 first observed two non-target B/angle continuation states passing qualification despite remaining thousands of pixels from oracle and unauthenticated. No label changes production behavior.\n')
print('Build61 third-pair sibling-stencil diagnostic written to:')
for p in [parent_path,state_path,summary_path,md_path]: print(' ',p)
print('\t'.join(summary_fields))
for r in summaries: print('\t'.join(str(r.get(k,'')) for k in summary_fields))
PY
