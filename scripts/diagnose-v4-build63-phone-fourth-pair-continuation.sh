#!/usr/bin/env bash
set -euo pipefail
PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD63_DIAGNOSTIC_DIR:-v4-phone private/build63-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD63_TIMEOUT:-86400}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"
command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"
for file in phone-b-mild.jpg phone-b-angle.jpg; do
  input="$ACQUISITION_DIR/$file"; [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
  stem="${file%.*}"; json="$OUTPUT_DIR/blind/$stem-fourth-pair-continuation.json"
  echo "Build63 blind fourth-pair continuation diagnostic: $file"
  timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-fourth-pair-continue -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
done
oracle_available=false
if [[ -f "$REFERENCE" ]] && python3 - <<'PY' >/dev/null 2>&1
import cv2
assert hasattr(cv2,'SIFT_create')
PY
then
  oracle_available=true
  for file in phone-b-mild.jpg phone-b-angle.jpg; do input="$ACQUISITION_DIR/$file"; stem="${file%.*}"; python3 ./scripts/build45-reference-register.py --reference "$REFERENCE" --acquired "$input" --out "$OUTPUT_DIR/oracle/$stem-quad.json"; done
else
  echo "Build63 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi
python3 - "$OUTPUT_DIR" "$oracle_available" <<'PY'
import csv,json,math,sys
from pathlib import Path
out=Path(sys.argv[1]); use_oracle=sys.argv[2].lower()=='true'
images=['phone-b-mild.jpg','phone-b-angle.jpg']; targets={'phone-b-mild.jpg':10,'phone-b-angle.jpg':18}
def xy(p): return float(p.get('x',p.get('X'))),float(p.get('y',p.get('Y')))
def err(q,o):
    if not o or len(q or [])!=4:return None
    return sum(math.hypot(xy(a)[0]-xy(b)[0],xy(a)[1]-xy(b)[1]) for a,b in zip(q,o))/4
f3=lambda v:'' if v is None else f'{v:.3f}'; f6=lambda v:'' if v is None else f'{float(v):.6f}'; b=lambda v:str(bool(v)).lower()
parents=[]; states=[]; summaries=[]
for name in images:
    stem=Path(name).stem; rep=json.load(open(out/'blind'/f'{stem}-fourth-pair-continuation.json'))['fourth_pair_continuation']; oracle=None
    if use_oracle: oracle=json.load(open(out/'oracle'/f'{stem}-quad.json')).get('quad')
    bank=[]; target_parents=[]
    for c in rep.get('candidates') or []:
      seed_e=err(c.get('seed_source_quad'),oracle)
      for root in c.get('roots') or []:
       for br in root.get('branches') or []:
        for cs in br.get('continuation_states') or []:
         for s1 in cs.get('siblings') or []:
          for p2 in s1.get('pair_states') or []:
           for c2 in p2.get('continuation_states') or []:
            for s2 in c2.get('sibling_states') or []:
             for p3 in s2.get('pair_states') or []:
              for c3 in p3.get('continuation_states') or []:
               for s3 in c3.get('sibling_states') or []:
                for p4 in s3.get('pair_states') or []:
                 pe=err(p4.get('source_quad'),oracle); cont=p4.get('continuation_states') or []
                 meta={'image':name,'candidate':c.get('index'),'pair_rank':c.get('source_pair_rank'),'seed_rank_pair':c.get('seed_rank_within_pair'),'fourth_pair_rank':p4.get('rank'),'parent_proposal':p4.get('proposal'),'parent_validation':p4.get('validation'),'parent_qualified':p4.get('qualified'),'parent_hmac':p4.get('single_hmac_authenticated'),'parent_oracle_px':pe,'continuation_states':len(cont)}
                 parents.append({**meta,'parent_proposal':f6(meta['parent_proposal']),'parent_validation':f6(meta['parent_validation']),'parent_qualified':b(meta['parent_qualified']),'parent_hmac':b(meta['parent_hmac']),'parent_oracle_px':f3(pe)})
                 if c.get('index')==targets[name]: target_parents.append((meta,p4,pe))
                 for st in cont:
                    se=err(st.get('source_quad'),oracle); row={**meta,'continuation_index':st.get('index'),'pass':st.get('pass'),'dimension':st.get('dimension'),'delta_px':st.get('delta_px'),'proposal':st.get('proposal'),'validation':st.get('validation'),'qualified':st.get('qualified'),'hmac':st.get('single_hmac_authenticated'),'state_oracle_px':se,'change_from_parent_px':None if se is None or pe is None else se-pe,'seed_oracle_px':seed_e}; bank.append((c,meta,st,se,pe))
                    states.append({k:(f3(v) if k in ('parent_oracle_px','state_oracle_px','change_from_parent_px','seed_oracle_px') else f6(v) if k in ('parent_proposal','parent_validation','proposal','validation') else b(v) if k in ('parent_qualified','parent_hmac','qualified','hmac') else v) for k,v in row.items()})
    tgt=targets[name]; tb=[x for x in bank if x[0].get('index')==tgt]; q=[x for x in bank if x[2].get('qualified')]; h=[x for x in bank if x[2].get('single_hmac_authenticated')]; tq=[x for x in tb if x[2].get('qualified')]
    best=min((x for x in tb if x[3] is not None),default=None,key=lambda x:x[3]); bestq=min((x for x in tq if x[3] is not None),default=None,key=lambda x:x[3]); better=[x for x in tb if x[3] is not None and x[4] is not None and float(x[2].get('proposal',0))>float(x[1].get('parent_proposal',0))+1e-7 and x[3]<x[4]-1e-9]; qgain=[x for x in better if x[2].get('qualified')]
    target_c=next((c for c in rep.get('candidates') or [] if c.get('index')==tgt),{})
    r={'image':name,'target_candidate':tgt,'seed_oracle_px':f3(err(target_c.get('seed_source_quad'),oracle)),'fourth_pair_parents':len(target_parents),'continuation_states':len(tb),'bank_continuation_states':len(bank),'bank_continuation_qualified':len(q),'bank_continuation_hmac':len(h),'bank_continuation_qualified_rate':f6(len(q)/len(bank) if bank else 0),'any_continuation_better_both':b(bool(better)),'any_continuation_qualified_gain':b(bool(qgain)),'any_continuation_hmac':b(bool(h))}
    if best:
      _,m,st,se,pe=best; r.update({'best_parent_oracle_px':f3(pe),'best_state_index':st.get('index'),'best_state_oracle_px':f3(se),'best_state_change_from_parent_px':f3(se-pe if se is not None and pe is not None else None),'best_state_proposal':f6(st.get('proposal')),'best_state_validation':f6(st.get('validation')),'best_state_qualified':b(st.get('qualified')),'best_state_hmac':b(st.get('single_hmac_authenticated'))})
    if bestq:r['best_qualified_state_oracle_px']=f3(bestq[3])
    r['classification']='fourth-pair-continuation-recovery' if h else 'fourth-pair-continuation-qualified-gain' if qgain else 'fourth-pair-continuation-geometric-gain' if better else 'fourth-pair-continuation-proposal-only' if tb else 'fourth-pair-continuation-not-triggered'; summaries.append(r)
parent_fields=['image','candidate','pair_rank','seed_rank_pair','fourth_pair_rank','parent_proposal','parent_validation','parent_qualified','parent_hmac','parent_oracle_px','continuation_states']
state_fields=['image','candidate','pair_rank','seed_rank_pair','fourth_pair_rank','parent_proposal','parent_validation','parent_qualified','parent_hmac','parent_oracle_px','continuation_states','continuation_index','pass','dimension','delta_px','proposal','validation','qualified','hmac','state_oracle_px','change_from_parent_px','seed_oracle_px']
summary_fields=['image','target_candidate','seed_oracle_px','fourth_pair_parents','continuation_states','bank_continuation_states','bank_continuation_qualified','bank_continuation_hmac','bank_continuation_qualified_rate','best_parent_oracle_px','best_state_index','best_state_oracle_px','best_state_change_from_parent_px','best_state_proposal','best_state_validation','best_state_qualified','best_state_hmac','best_qualified_state_oracle_px','any_continuation_better_both','any_continuation_qualified_gain','any_continuation_hmac','classification']
for fn,fields,rows in [('build63-fourth-pair-parents.tsv',parent_fields,parents),('build63-continuation-states.tsv',state_fields,states),('build63-fourth-pair-continuation-summary.tsv',summary_fields,summaries)]:
  with open(out/fn,'w',newline='') as f:w=csv.DictWriter(f,fieldnames=fields,delimiter='\t',extrasaction='ignore');w.writeheader();w.writerows(rows)
with open(out/'build63-fourth-pair-continuation.md','w') as f:
  f.write('# Build63 fourth-pair continuation diagnostic\n\nBuild63 continues every frozen Build62 fourth-pair state with bounded proposal-only 1px coordinate descent and retains every accepted intermediate. Qualification/HMAC occur only after the complete B/mild+B/angle bank is frozen; SIFT/reference geometry is post-hoc only.\n\n')
  f.write('| image | target | continuation states | bank qualified/total | bank HMAC | best parent px | best state px | qualified | HMAC | interpretation |\n|---|---:|---:|---:|---:|---:|---:|---|---|---|\n')
  for r in summaries:f.write(f"| {r['image']} | {r['target_candidate']} | {r['continuation_states']} | {r['bank_continuation_qualified']}/{r['bank_continuation_states']} | {r['bank_continuation_hmac']} | {r.get('best_parent_oracle_px','n/a') or 'n/a'} | {r.get('best_state_oracle_px','n/a') or 'n/a'} | {r.get('best_state_qualified','n/a') or 'n/a'} | {r.get('best_state_hmac','n/a') or 'n/a'} | {r['classification']} |\n")
print('Build63 fourth-pair continuation diagnostic written to:'); [print(' ',out/x) for x in ['build63-fourth-pair-parents.tsv','build63-continuation-states.tsv','build63-fourth-pair-continuation-summary.tsv','build63-fourth-pair-continuation.md']]
print('\t'.join(summary_fields)); [print('\t'.join(str(r.get(k,'')) for k in summary_fields)) for r in summaries]
PY
