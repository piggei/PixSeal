#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD54_DIAGNOSTIC_DIR:-v4-phone private/build54-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD54_TIMEOUT:-3600}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind Build53 pair generation plus proposal-only 1px continuation.
# Every accepted continuation intermediate is retained. Qualification/HMAC are
# attached only after the full geometry bank for both images is frozen.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-pair-continuation.json"
    echo "Build54 blind post-pair continuation diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-pair-continue \
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
    echo "Build54 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv, json, math, sys
from pathlib import Path

out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']
targets={'phone-b-mild.jpg':10,'phone-b-angle.jpg':18}
branch_path=out/'build54-branches.tsv'
state_path=out/'build54-continuation-states.tsv'
summary_path=out/'build54-pair-continuation-summary.tsv'
md_path=out/'build54-pair-continuation.md'

def xy(p): return float(p.get('x',p.get('X'))), float(p.get('y',p.get('Y')))
def err(q, oracle):
    if not oracle or len(q or [])!=4: return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,oracle))/4.0
def f3(v): return '' if v is None else f'{v:.3f}'
def f6(v): return '' if v is None else f'{float(v):.6f}'
def b(v): return str(bool(v)).lower()

branch_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','parent_pair_rank','parent_pair_proposal','parent_pair_validation','parent_pair_qualified','parent_pair_hmac','parent_pair_oracle_px','continuation_evaluations','continuation_states']
state_fields=['image','candidate','seed_rank_pair','pair','pair_rank','root_index','parent_pair_rank','state_index','pass','dimension','corner','axis','delta_px','proposal','validation','qualified','hmac','movement_from_seed_px','seed_oracle_px','parent_pair_oracle_px','state_oracle_px','change_from_parent_px','change_from_seed_px']
summary_fields=['image','target_candidate','pair','pair_rank','seed_rank_pair','seed_oracle_px','baseline_oracle_px','baseline_change_px','baseline_proposal','baseline_validation','baseline_qualified','baseline_hmac','coordinate_local_roots','pair_states','continuation_states','best_parent_root_index','best_parent_pair_rank','best_parent_oracle_px','best_parent_qualified','best_continuation_root_index','best_continuation_parent_pair_rank','best_continuation_index','best_continuation_oracle_px','best_continuation_change_from_parent_px','best_continuation_proposal','best_continuation_validation','best_continuation_qualified','best_continuation_hmac','best_qualified_continuation_root_index','best_qualified_continuation_parent_pair_rank','best_qualified_continuation_index','best_qualified_continuation_oracle_px','any_continuation_better_both','any_continuation_qualified_gain','any_continuation_hmac','classification']

branch_rows=[]; state_rows=[]; summaries=[]
for name in images:
    stem=Path(name).stem
    j=json.load(open(out/'blind'/f'{stem}-pair-continuation.json',encoding='utf-8'))
    rep=j['continuation']
    oracle=None
    if use_oracle:
        oracle=json.load(open(out/'oracle'/f'{stem}-quad.json',encoding='utf-8')).get('quad')
    candidates=rep.get('candidates') or []
    for c in candidates:
        seed_err=err(c.get('seed_source_quad'),oracle)
        for root in c.get('roots') or []:
            for branch in root.get('branches') or []:
                parent=branch.get('parent_pair') or {}
                pe=err(parent.get('source_quad'),oracle)
                branch_rows.append({
                    'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'parent_pair_rank':branch.get('pair_rank',''),
                    'parent_pair_proposal':f6(parent.get('proposal')),'parent_pair_validation':f6(parent.get('validation')),'parent_pair_qualified':b(parent.get('qualified')),'parent_pair_hmac':b(parent.get('single_hmac_authenticated')),'parent_pair_oracle_px':f3(pe),'continuation_evaluations':branch.get('continuation_evaluations',''),'continuation_states':len(branch.get('continuation_states') or [])})
                for st in branch.get('continuation_states') or []:
                    se=err(st.get('source_quad'),oracle)
                    state_rows.append({
                        'image':name,'candidate':c.get('index',''),'seed_rank_pair':c.get('seed_rank_within_pair',''),'pair':c.get('source_pair',''),'pair_rank':c.get('source_pair_rank',''),'root_index':root.get('root_index',''),'parent_pair_rank':branch.get('pair_rank',''),'state_index':st.get('index',''),'pass':st.get('pass',''),'dimension':st.get('dimension',''),'corner':st.get('corner',''),'axis':st.get('axis',''),'delta_px':st.get('delta_px',''),
                        'proposal':f6(st.get('proposal')),'validation':f6(st.get('validation')),'qualified':b(st.get('qualified')),'hmac':b(st.get('single_hmac_authenticated')),'movement_from_seed_px':f3(st.get('mean_movement_from_seed_px')),'seed_oracle_px':f3(seed_err),'parent_pair_oracle_px':f3(pe),'state_oracle_px':f3(se),'change_from_parent_px':f3(None if se is None or pe is None else se-pe),'change_from_seed_px':f3(None if se is None or seed_err is None else se-seed_err),
                        '_oracle':se,'_parent_oracle':pe,'_proposal_raw':st.get('proposal'),'_parent_proposal_raw':parent.get('proposal'),'_qualified':bool(st.get('qualified')),'_hmac':bool(st.get('single_hmac_authenticated'))})

    target_idx=targets[name]
    target=next((c for c in candidates if c.get('index')==target_idx),None)
    if target is None:
        summaries.append({'image':name,'target_candidate':target_idx,'classification':'target-missing','any_continuation_better_both':'false','any_continuation_qualified_gain':'false','any_continuation_hmac':'false'})
        continue
    seed_e=err(target.get('seed_source_quad'),oracle); base=target.get('baseline') or {}; base_e=err(base.get('source_quad'),oracle)
    tbranches=[]; tstates=[]
    for root in target.get('roots') or []:
        for branch in root.get('branches') or []:
            parent=branch.get('parent_pair') or {}; pe=err(parent.get('source_quad'),oracle)
            tbranches.append((root,branch,parent,pe))
            for st in branch.get('continuation_states') or []:
                se=err(st.get('source_quad'),oracle)
                tstates.append((root,branch,parent,pe,st,se))
    locals_count=sum(1 for r in target.get('roots') or [] if r.get('coordinate_local'))
    pair_count=len(tbranches)
    row={'image':name,'target_candidate':target_idx,'pair':target.get('source_pair',''),'pair_rank':target.get('source_pair_rank',''),'seed_rank_pair':target.get('seed_rank_within_pair',''),'seed_oracle_px':f3(seed_e),'baseline_oracle_px':f3(base_e),'baseline_change_px':f3(None if base_e is None or seed_e is None else base_e-seed_e),'baseline_proposal':f6(base.get('proposal')),'baseline_validation':f6(base.get('validation')),'baseline_qualified':b(base.get('qualified')),'baseline_hmac':b(base.get('single_hmac_authenticated')),'coordinate_local_roots':locals_count,'pair_states':pair_count,'continuation_states':len(tstates)}
    if tbranches and use_oracle:
        bp=min(tbranches,key=lambda x:x[3] if x[3] is not None else float('inf'))
        row.update({'best_parent_root_index':bp[0].get('root_index',''),'best_parent_pair_rank':bp[1].get('pair_rank',''),'best_parent_oracle_px':f3(bp[3]),'best_parent_qualified':b(bp[2].get('qualified'))})
    if tstates and use_oracle:
        bs=min(tstates,key=lambda x:x[5] if x[5] is not None else float('inf'))
        row.update({'best_continuation_root_index':bs[0].get('root_index',''),'best_continuation_parent_pair_rank':bs[1].get('pair_rank',''),'best_continuation_index':bs[4].get('index',''),'best_continuation_oracle_px':f3(bs[5]),'best_continuation_change_from_parent_px':f3(None if bs[5] is None or bs[3] is None else bs[5]-bs[3]),'best_continuation_proposal':f6(bs[4].get('proposal')),'best_continuation_validation':f6(bs[4].get('validation')),'best_continuation_qualified':b(bs[4].get('qualified')),'best_continuation_hmac':b(bs[4].get('single_hmac_authenticated'))})
        qs=[x for x in tstates if x[4].get('qualified') and x[5] is not None]
        if qs:
            bq=min(qs,key=lambda x:x[5])
            row.update({'best_qualified_continuation_root_index':bq[0].get('root_index',''),'best_qualified_continuation_parent_pair_rank':bq[1].get('pair_rank',''),'best_qualified_continuation_index':bq[4].get('index',''),'best_qualified_continuation_oracle_px':f3(bq[5])})
    better=any(se is not None and pe is not None and se < pe-1e-6 and float(st.get('proposal',-1e9)) > float(parent.get('proposal',-1e9))+1e-7 for _,_,parent,pe,st,se in tstates)
    qgain=any(se is not None and pe is not None and se < pe-1e-6 and st.get('qualified') and float(st.get('proposal',-1e9)) > float(parent.get('proposal',-1e9))+1e-7 for _,_,parent,pe,st,se in tstates)
    anyh=any(bool(st.get('single_hmac_authenticated')) for *_,st,se in tstates)
    if anyh: cls='continuation-recovery'
    elif qgain: cls='continuation-qualified-gain'
    elif better: cls='continuation-geometric-gain'
    elif tstates: cls='continuation-no-oracle-gain'
    elif tbranches: cls='continuation-no-accepted-state'
    else: cls='continuation-not-triggered'
    row.update({'any_continuation_better_both':b(better),'any_continuation_qualified_gain':b(qgain),'any_continuation_hmac':b(anyh),'classification':cls})
    summaries.append(row)

with open(branch_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=branch_fields,delimiter='\t',extrasaction='ignore');w.writeheader();w.writerows(branch_rows)
with open(state_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=state_fields,delimiter='\t',extrasaction='ignore');w.writeheader();w.writerows(state_rows)
with open(summary_path,'w',newline='',encoding='utf-8') as f:
    w=csv.DictWriter(f,fieldnames=summary_fields,delimiter='\t',extrasaction='ignore');w.writeheader();w.writerows(summaries)
with open(md_path,'w',encoding='utf-8') as f:
    f.write('# PixSeal Build54 post-pair continuation diagnostic\n\n')
    f.write('Research-only. Build54 reproduces the unchanged Build53 pair bank and continues every retained pair state with bounded 1px proposal-only coordinate descent, retaining every accepted intermediate. The complete geometry bank is frozen before held-out/full-pilot qualification or diagnostic HMAC; SIFT/reference geometry is generated only after both blind JSON files exist.\n\n')
    f.write('## Post-hoc summary\n\n')
    f.write('| image | target | pair | seed px | baseline px | pair states | continuation states | best parent px | best continuation px | continuation qualified | continuation HMAC | qualified gain | interpretation |\n')
    f.write('|---|---:|---|---:|---:|---:|---:|---:|---:|---|---|---|---|\n')
    for r in summaries:
        f.write(f"| {r.get('image','')} | {r.get('target_candidate','')} | {r.get('pair','')} | {r.get('seed_oracle_px','n/a')} | {r.get('baseline_oracle_px','n/a')} | {r.get('pair_states','0')} | {r.get('continuation_states','0')} | {r.get('best_parent_oracle_px','n/a')} | {r.get('best_continuation_oracle_px','n/a')} | {r.get('best_continuation_qualified','n/a')} | {r.get('best_continuation_hmac','n/a')} | {r.get('any_continuation_qualified_gain','false')} | {r.get('classification','')} |\n")
    f.write('\nInterpretation is post-hoc only. `continuation-recovery` requires HMAC authentication. `continuation-qualified-gain` requires a retained continuation state to improve proposal and independent oracle error relative to its frozen parent pair state while passing unchanged qualification. No label changes production behavior.\n')
PY

echo "Build54 post-pair continuation diagnostic written to:"
echo "  $OUTPUT_DIR/build54-branches.tsv"
echo "  $OUTPUT_DIR/build54-continuation-states.tsv"
echo "  $OUTPUT_DIR/build54-pair-continuation-summary.tsv"
echo "  $OUTPUT_DIR/build54-pair-continuation.md"
cat "$OUTPUT_DIR/build54-pair-continuation-summary.tsv"
