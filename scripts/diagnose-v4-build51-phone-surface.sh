#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD51_DIAGNOSTIC_DIR:-v4-phone private/build51-diagnostics}"
KEY="${V4_PHONE_KEY:-PixSeal-v4-TestKey-2026}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_BUILD51_TIMEOUT:-2400}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind top-4 seed selection, exact coordinate-descent trace, and
# deterministic local stencil. No reference/oracle image is used here.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
    input="$ACQUISITION_DIR/$file"
    [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
    stem="${file%.*}"
    json="$OUTPUT_DIR/blind/$stem-surface.json"
    echo "Build51 blind local proposal-surface diagnostic: $file"
    timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-surface \
        -in "$input" -key "$KEY" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$json"
done

# Phase 2: reference/SIFT oracle is generated only after both blind outputs
# already exist. It can annotate but cannot alter any geometry or score.
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
    echo "Build51 note: reference/OpenCV-SIFT unavailable; oracle columns remain blank." >&2
fi

seeds_tsv="$OUTPUT_DIR/build51-surface-seeds.tsv"
trace_tsv="$OUTPUT_DIR/build51-refinement-trajectory.tsv"
stencil_tsv="$OUTPUT_DIR/build51-local-stencil.tsv"
summary_tsv="$OUTPUT_DIR/build51-surface-summary.tsv"
md="$OUTPUT_DIR/build51-local-surface.md"

printf 'image\tcandidate\tseed\tseed_rank_pair\tpair\tpair_rank\ttier\tcell_rank\tpre_proposal\tpre_validation\tpre_qualified\tpost_proposal\tpost_validation\tpost_qualified\tpost_hmac\tpre_oracle_mean_px\tpost_oracle_mean_px\toracle_change_px\n' > "$seeds_tsv"
printf 'image\tcandidate\tseed_rank_pair\tpair\teval\tstep_px\tpass\tdim\tcorner\taxis\tdelta_px\tproposal_before\tproposal_after\tproposal_delta\tproposal_available\taccepted\tvalidation\tvalidation_available\toracle_before_px\toracle_after_px\toracle_change_px\n' > "$trace_tsv"
printf 'image\tcandidate\tseed_rank_pair\tpair\tbasis\tdelta_px\tproposal\tproposal_delta_center\tvalidation\toracle_mean_px\toracle_delta_center_px\tbetter_proposal_and_oracle\n' > "$stencil_tsv"
printf 'image\tseeds\ttarget_candidate\tpair\tpair_rank\tseed_rank_pair\tpre_mean_px\tpost_mean_px\tpost_change_px\taccepted_steps\taccepted_oracle_improving\taccepted_oracle_worsening\ttrace_best_oracle_px\tstencil_best_oracle_px\tstencil_best_basis\tstencil_best_delta\tstencil_best_proposal_delta\tstencil_best_proposal_oracle_px\tany_stencil_better_both\tspearman_proposal_vs_oracle_error\tpost_qualified\tpost_hmac\tclassification\n' > "$summary_tsv"

for file in phone-b-mild.jpg phone-b-angle.jpg; do
    stem="${file%.*}"; json="$OUTPUT_DIR/blind/$stem-surface.json"; quad="$OUTPUT_DIR/oracle/$stem-quad.json"
    python3 - "$json" "$quad" "$file" "$oracle_available" "$seeds_tsv" "$trace_tsv" "$stencil_tsv" "$summary_tsv" <<'PY'
import json, math, os, sys

j=json.load(open(sys.argv[1],encoding='utf-8')); r=j.get('surface') or {}
cands=r.get('candidates') or []
name=sys.argv[3]; use_oracle=sys.argv[4].lower()=='true' and os.path.isfile(sys.argv[2])
seed_path, trace_path, stencil_path, summary_path=sys.argv[5:9]
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

def ranks(vals):
    order=sorted(range(len(vals)), key=lambda i: vals[i])
    out=[0.0]*len(vals); k=0
    while k < len(order):
        z=k+1
        while z < len(order) and vals[order[z]]==vals[order[k]]: z+=1
        avg=(k+1+z)/2.0
        for p in range(k,z): out[order[p]]=avg
        k=z
    return out

def spearman(xs, ys):
    if len(xs)<3: return None
    rx,ry=ranks(xs),ranks(ys)
    mx=sum(rx)/len(rx); my=sum(ry)/len(ry)
    num=sum((a-mx)*(b-my) for a,b in zip(rx,ry))
    dx=sum((a-mx)**2 for a in rx); dy=sum((b-my)**2 for b in ry)
    if dx<=0 or dy<=0: return None
    return num/math.sqrt(dx*dy)

for c in cands:
    c['_pre_err']=err(c.get('source_quad_before'))
    c['_post_err']=err(c.get('source_quad_after'))
    for t in c.get('trace') or []:
        t['_before_err']=err(t.get('source_quad_before'))
        t['_after_err']=err(t.get('source_quad_after'))
    for s in c.get('stencil') or []:
        s['_err']=err(s.get('source_quad'))

with open(seed_path,'a',encoding='utf-8') as f:
    for c in cands:
        pe,qe=c['_pre_err'],c['_post_err']
        vals=[name,c.get('index',''),c.get('seed_index',''),c.get('seed_rank_within_pair',''),c.get('source_pair',''),c.get('source_pair_rank',''),c.get('source_tier',''),c.get('cell_rank',''),
              f"{c.get('pre_proposal',0):.6f}",f"{c.get('pre_validation',0):.6f}",str(c.get('pre_qualified',False)).lower(),
              f"{c.get('post_proposal',0):.6f}",f"{c.get('post_validation',0):.6f}",str(c.get('post_qualified',False)).lower(),str(c.get('single_hmac_authenticated',False)).lower(),
              '' if pe is None else f'{pe:.3f}','' if qe is None else f'{qe:.3f}','' if pe is None or qe is None else f'{qe-pe:.3f}']
        f.write('\t'.join(map(str,vals))+'\n')

with open(trace_path,'a',encoding='utf-8') as f:
    for c in cands:
        for t in c.get('trace') or []:
            be,ae=t['_before_err'],t['_after_err']
            pav=bool(t.get('proposal_available',True)); vav=bool(t.get('validation_available',True))
            pa=float(t.get('proposal_after',0)) if pav else None
            vals=[name,c.get('index',''),c.get('seed_rank_within_pair',''),c.get('source_pair',''),t.get('eval_index',''),t.get('step_size_px',''),t.get('pass',''),t.get('dimension',''),t.get('corner',''),t.get('axis',''),t.get('delta_px',''),
                  f"{t.get('proposal_before',0):.6f}",'' if pa is None else f"{pa:.6f}",'' if pa is None else f"{pa-t.get('proposal_before',0):.6f}",str(pav).lower(),str(t.get('accepted',False)).lower(),'' if not vav else f"{t.get('validation_score',0):.6f}",str(vav).lower(),
                  '' if be is None else f'{be:.3f}','' if ae is None else f'{ae:.3f}','' if be is None or ae is None else f'{ae-be:.3f}']
            f.write('\t'.join(map(str,vals))+'\n')

with open(stencil_path,'a',encoding='utf-8') as f:
    for c in cands:
        sts=c.get('stencil') or []
        center=next((s for s in sts if s.get('basis')=='center'),None)
        cp=float(center.get('proposal_score',0)) if center else float(c.get('pre_proposal',0))
        ce=center.get('_err') if center else c['_pre_err']
        for s in sts:
            e=s['_err']; pd=float(s.get('proposal_score',0))-cp
            ed=None if e is None or ce is None else e-ce
            both=(e is not None and ce is not None and pd>1e-7 and e<ce-1e-9)
            vals=[name,c.get('index',''),c.get('seed_rank_within_pair',''),c.get('source_pair',''),s.get('basis',''),s.get('delta_px',''),f"{s.get('proposal_score',0):.6f}",f'{pd:.6f}',f"{s.get('validation_score',0):.6f}",
                  '' if e is None else f'{e:.3f}','' if ed is None else f'{ed:.3f}',str(both).lower()]
            f.write('\t'.join(map(str,vals))+'\n')

# Post-hoc focus: candidate whose unrefined seed is nearest the oracle.
target=min((c for c in cands if c['_pre_err'] is not None), key=lambda c:c['_pre_err'], default=None)
if target is None:
    vals=[name,len(cands),'','','','','','','','','','','','','','','','','','','','','oracle-unavailable']
else:
    pre=target['_pre_err']; post=target['_post_err']
    accepted=[t for t in target.get('trace') or [] if t.get('accepted')]
    improving=sum(1 for t in accepted if t['_before_err'] is not None and t['_after_err'] is not None and t['_after_err'] < t['_before_err']-1e-9)
    worsening=sum(1 for t in accepted if t['_before_err'] is not None and t['_after_err'] is not None and t['_after_err'] > t['_before_err']+1e-9)
    all_trace=[t for t in target.get('trace') or [] if t['_after_err'] is not None and bool(t.get('proposal_available',True))]
    trace_best=min((t['_after_err'] for t in all_trace), default=None)
    sts=target.get('stencil') or []
    center=next((s for s in sts if s.get('basis')=='center'),None)
    cp=float(center.get('proposal_score',0)) if center else float(target.get('pre_proposal',0))
    ce=center.get('_err') if center else pre
    valid=[s for s in sts if s['_err'] is not None]
    best_err=min(valid,key=lambda s:s['_err'],default=None)
    best_prop=max(valid,key=lambda s:float(s.get('proposal_score',0)),default=None)
    both=[s for s in valid if ce is not None and float(s.get('proposal_score',0))>cp+1e-7 and s['_err']<ce-1e-9]
    rho=spearman([float(s.get('proposal_score',0)) for s in valid],[s['_err'] for s in valid]) if valid else None
    if both:
        classification='optimizer-opportunity'
    elif post is not None and pre is not None and post > pre+1e-6 and worsening > improving:
        classification='proposal-surface-misaligned'
    else:
        classification='mixed-or-inconclusive'
    vals=[name,len(cands),target.get('index',''),target.get('source_pair',''),target.get('source_pair_rank',''),target.get('seed_rank_within_pair',''),
          f'{pre:.3f}', '' if post is None else f'{post:.3f}', '' if post is None else f'{post-pre:.3f}',len(accepted),improving,worsening,
          '' if trace_best is None else f'{trace_best:.3f}', '' if best_err is None else f"{best_err['_err']:.3f}", '' if best_err is None else best_err.get('basis',''), '' if best_err is None else best_err.get('delta_px',''),
          '' if best_err is None else f"{float(best_err.get('proposal_score',0))-cp:.6f}", '' if best_prop is None else f"{best_prop['_err']:.3f}",str(bool(both)).lower(),'' if rho is None else f'{rho:.6f}',
          str(target.get('post_qualified',False)).lower(),str(target.get('single_hmac_authenticated',False)).lower(),classification]
with open(summary_path,'a',encoding='utf-8') as f:
    f.write('\t'.join(map(str,vals))+'\n')
PY
done

{
    echo '# PixSeal Build51 local proposal-surface / refinement-trajectory diagnostic'
    echo
    echo 'Research-only. The blind command records the exact Build41 proposal-only coordinate-descent attempts and a fixed local stencil around every top-4-per-pair seed. All geometry is frozen before held-out annotation. Full-pilot qualification and diagnostic HMAC are evaluated only for the original/final seed states. SIFT/reference oracle geometry is generated only after both blind JSON files already exist.'
    echo
    echo '## Post-hoc summary'
    echo
    echo '| image | target | pair | pair rank | seed rank | pre px | post px | change | accepted | accepted improve | accepted worsen | best trace px | best stencil px | stencil mode | d | best stencil proposal delta | best-proposal stencil px | better-both | rho(proposal,error) | qualified | HMAC | interpretation |'
    echo '|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|---:|---|---:|---|---|---|'
    tail -n +2 "$summary_tsv" | while IFS=$'\t' read -r image seeds target pair prank srank pre post change acc imp wors tbest sbest sbasis sdelta spdelta bp_err both rho qual hmac class; do
      echo "| $image | ${target:-n/a} | ${pair:-n/a} | ${prank:-n/a} | ${srank:-n/a} | ${pre:-n/a} | ${post:-n/a} | ${change:-n/a} | ${acc:-n/a} | ${imp:-n/a} | ${wors:-n/a} | ${tbest:-n/a} | ${sbest:-n/a} | ${sbasis:-n/a} | ${sdelta:-n/a} | ${spdelta:-n/a} | ${bp_err:-n/a} | ${both:-n/a} | ${rho:-n/a} | ${qual:-n/a} | ${hmac:-n/a} | ${class:-n/a} |"
    done
    echo
    echo 'Interpretation labels are post-hoc laboratory summaries only. `optimizer-opportunity` means the fixed stencil contains at least one sample with both higher proposal and lower oracle error than the seed. `proposal-surface-misaligned` means the accepted proposal ascent predominantly moves away from the oracle and the final refined geometry is worse. Neither label changes production behavior.'
    echo
    echo 'See `build51-refinement-trajectory.tsv` for every evaluated +/- coordinate move and `build51-local-stencil.tsv` for every deterministic local sample.'
} > "$md"

echo "Build51 local proposal-surface diagnostic written to:"
echo "  $seeds_tsv"
echo "  $trace_tsv"
echo "  $stencil_tsv"
echo "  $summary_tsv"
echo "  $md"
cat "$summary_tsv"
