#!/usr/bin/env bash
set -euo pipefail

PIXSEAL="${PIXSEAL:?set PIXSEAL to the built pixseal executable}"
ACQUISITION_DIR="${V4_PHONE_ACQUISITION_DIR:-v4-phone private/build38-acquired}"
GENERATED_DIR="${V4_PHONE_OUTPUT_DIR:-v4-phone private/build38-generated}"
OUTPUT_DIR="${V4_PHONE_BUILD49_DIAGNOSTIC_DIR:-v4-phone private/build49-diagnostics}"
CANONICAL_WIDTH="${V4_PHONE_CANONICAL_WIDTH:-1632}"
CANONICAL_HEIGHT="${V4_PHONE_CANONICAL_HEIGHT:-1632}"
PHONE_TIMEOUT="${V4_PHONE_TIMEOUT:-1200}"
REFERENCE="$GENERATED_DIR/pixseal-build38-mq-marked-b.png"

command -v timeout >/dev/null 2>&1 || { echo "error: GNU timeout is required" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "error: python3 is required" >&2; exit 1; }
[[ -x "$PIXSEAL" ]] || { echo "error: PIXSEAL is not executable: $PIXSEAL" >&2; exit 1; }
mkdir -p "$OUTPUT_DIR/blind" "$OUTPUT_DIR/oracle"

# Phase 1: blind proposal observability only. No key/reference/oracle is available.
for file in phone-b-mild.jpg phone-b-angle.jpg; do
  input="$ACQUISITION_DIR/$file"; [[ -f "$input" ]] || { echo "error: missing acquisition: $input" >&2; exit 1; }
  stem="${file%.*}"
  echo "Build49 blind proposal-ranking diagnostic: $file"
  timeout --foreground "${PHONE_TIMEOUT}s" "$PIXSEAL" v4-diagnose-phone-ranking \
    -in "$input" -width "$CANONICAL_WIDTH" -height "$CANONICAL_HEIGHT" -json > "$OUTPUT_DIR/blind/$stem-ranking.json"
done

# Phase 2: oracle only after both blind JSON files are fixed.
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
  echo "Build49 note: reference/OpenCV-SIFT unavailable; oracle ranking summary cannot be produced." >&2
fi

summary="$OUTPUT_DIR/build49-proposal-ranking-summary.tsv"
depth="$OUTPUT_DIR/build49-seed-depth.tsv"
corr="$OUTPUT_DIR/build49-observable-correlation.tsv"
detail="$OUTPUT_DIR/build49-proposal-ranking.tsv"
md="$OUTPUT_DIR/build49-proposal-ranking.md"
printf 'image\tcandidate\tpair\tpair_rank\ttier\tcell_rank\tproposal\tproposal_rank_global\tproposal_rank_pair\tfold1\tfold2\tfold_min\tfold_gap\tfold_min_rank_pair\tbalanced\tbalanced_rank_pair\ttile_mean\ttile_std\ttile_consistency\ttile_rank_pair\tpair_score\tcell_mean\tcell_robust\tqualified\toracle_mean_px\toracle_max_px\n' > "$detail"
printf 'image\tcandidates\tqualified\toracle_nearest\toracle_mean_px\tpair\tpair_rank\ttier\tcell_rank\tproposal_rank_global\tproposal_rank_pair\tfold_min_rank_pair\tbalanced_rank_pair\ttile_rank_pair\n' > "$summary"
printf 'image\ttop_per_pair\tselected\tnearest_candidate\tnearest_mean_px\toracle_nearest_included\n' > "$depth"
printf 'image\tobservable\tspearman_vs_negative_oracle_error\n' > "$corr"

if [[ "$oracle_available" == true ]]; then
for file in phone-b-mild.jpg phone-b-angle.jpg; do
  stem="${file%.*}"
  python3 - "$file" "$OUTPUT_DIR/blind/$stem-ranking.json" "$OUTPUT_DIR/oracle/$stem-quad.json" "$detail" "$summary" "$depth" "$corr" <<'PY'
import json, math, sys
name,bpath,opath,detail_path,summary_path,depth_path,corr_path=sys.argv[1:]
b=json.load(open(bpath,encoding='utf-8'))['ranking']
o=json.load(open(opath,encoding='utf-8'))['quad']
oracle=[(q['x'],q['y']) for q in o]
cands=b.get('candidates',[])
for c in cands:
    q=[(p['X'],p['Y']) if 'X' in p else (p['x'],p['y']) for p in c['source_quad']]
    ds=[math.hypot(a[0]-z[0],a[1]-z[1]) for a,z in zip(q,oracle)]
    c['_omean']=sum(ds)/4; c['_omax']=max(ds)
nearest=min(cands,key=lambda x:x['_omean']) if cands else None
with open(detail_path,'a',encoding='utf-8') as f:
    for c in cands:
        vals=[name,c['index'],c['source_pair'],c['source_pair_rank'],c['source_tier'],c['cell_rank'],
              f"{c['proposal_score']:.6f}",c['proposal_rank_global'],c['proposal_rank_within_pair'],
              f"{c['fold1_score']:.6f}",f"{c['fold2_score']:.6f}",f"{c['fold_min_score']:.6f}",f"{c['fold_gap']:.6f}",c['fold_min_rank_within_pair'],
              f"{c['balanced_fold_score']:.6f}",c['balanced_rank_within_pair'],f"{c['proposal_tile_mean']:.6f}",f"{c['proposal_tile_stddev']:.6f}",
              f"{c['tile_consistency_score']:.6f}",c['tile_rank_within_pair'],f"{c['pair_robust_score']:.6f}",f"{c['cell_mean_score']:.6f}",f"{c['cell_robust_score']:.6f}",
              str(c['qualified']).lower(),f"{c['_omean']:.3f}",f"{c['_omax']:.3f}"]
        f.write('\t'.join(map(str,vals))+'\n')
if nearest:
    vals=[name,len(cands),sum(bool(c['qualified']) for c in cands),nearest['index'],f"{nearest['_omean']:.3f}",nearest['source_pair'],nearest['source_pair_rank'],nearest['source_tier'],nearest['cell_rank'],nearest['proposal_rank_global'],nearest['proposal_rank_within_pair'],nearest['fold_min_rank_within_pair'],nearest['balanced_rank_within_pair'],nearest['tile_rank_within_pair']]
    with open(summary_path,'a',encoding='utf-8') as f:f.write('\t'.join(map(str,vals))+'\n')
# Proposal top-N within each pair coverage.
for n in (2,4,6,8):
    selected=[c for c in cands if c['proposal_rank_within_pair']<=n]
    best=min(selected,key=lambda x:x['_omean']) if selected else None
    vals=[name,n,len(selected),best['index'] if best else '',f"{best['_omean']:.3f}" if best else '',str(bool(nearest and nearest in selected)).lower()]
    with open(depth_path,'a',encoding='utf-8') as f:f.write('\t'.join(map(str,vals))+'\n')
# Spearman correlation (positive means larger observable tracks lower error).
def ranks(vals):
    order=sorted(range(len(vals)),key=lambda i:vals[i])
    out=[0.0]*len(vals); i=0
    while i<len(order):
        j=i+1
        while j<len(order) and vals[order[j]]==vals[order[i]]: j+=1
        r=(i+j-1)/2+1
        for k in range(i,j): out[order[k]]=r
        i=j
    return out
def corr(a,b):
    if len(a)<2:return float('nan')
    ra,rb=ranks(a),ranks(b); ma=sum(ra)/len(ra); mb=sum(rb)/len(rb)
    num=sum((x-ma)*(y-mb) for x,y in zip(ra,rb)); da=sum((x-ma)**2 for x in ra); db=sum((y-mb)**2 for y in rb)
    return num/math.sqrt(da*db) if da>0 and db>0 else float('nan')
y=[-c['_omean'] for c in cands]
metrics={
 'proposal':lambda c:c['proposal_score'],'fold_min':lambda c:c['fold_min_score'],'balanced_fold':lambda c:c['balanced_fold_score'],
 'tile_mean':lambda c:c['proposal_tile_mean'],'tile_consistency':lambda c:c['tile_consistency_score'],'pair_score':lambda c:c['pair_robust_score'],
 'cell_mean':lambda c:c['cell_mean_score'],'cell_robust':lambda c:c['cell_robust_score'],'negative_fold_gap':lambda c:-c['fold_gap'],
 'negative_tile_std':lambda c:-c['proposal_tile_stddev']}
with open(corr_path,'a',encoding='utf-8') as f:
    for k,fn in metrics.items(): f.write(f"{name}\t{k}\t{corr([fn(c) for c in cands],y):.6f}\n")
PY
done
fi

{
 echo '# PixSeal Build49 proposal-ranking observability'
 echo
 echo 'Blind candidate generation and all proposal ranks are fixed before any oracle exists. The oracle is post-hoc only. No secret key or HMAC is used by the Build49 command.'
 echo
 echo '## Oracle-nearest candidate'
 echo
 echo '| image | candidates | qualified | nearest | mean px | pair | pair rank | tier | cell | proposal global rank | proposal pair rank | fold-min pair rank | balanced pair rank | tile pair rank |'
 echo '|---|---:|---:|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|'
 if [[ -s "$summary" ]]; then tail -n +2 "$summary" | while IFS=$'\t' read -r a b c d e f g h i j k l m n; do echo "| $a | $b | $c | $d | $e | $f | $g | $h | $i | $j | $k | $l | $m | $n |"; done; fi
 echo
 echo '## Proposal top-N per pair coverage'
 echo
 echo '| image | top N/pair | selected | nearest candidate | nearest mean px | oracle-nearest included |'
 echo '|---|---:|---:|---:|---:|---|'
 if [[ -s "$depth" ]]; then tail -n +2 "$depth" | while IFS=$'\t' read -r a b c d e f; do echo "| $a | $b | $c | $d | $e | $f |"; done; fi
 echo
 echo 'See the TSV files for candidate-level metrics and Spearman correlations.'
} > "$md"

echo "Build49 proposal-ranking diagnostic written to:"
echo "  $detail"
echo "  $summary"
echo "  $depth"
echo "  $corr"
echo "  $md"
[[ -s "$summary" ]] && cat "$summary"
