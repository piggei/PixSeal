#!/usr/bin/env python3
"""Apply the fixed Build104 benchmark gate, excluding identical-code controls."""
import argparse
import csv
from pathlib import Path
import statistics

METRICS=('ReaderFront','ReaderMild','ReaderAngle','Generation4','Pool8')
ORDER=['Q','C','C','Q','Q','Q','C','Q','Q','C']
def read(path):
 with path.open() as f: return list(csv.DictReader(f,delimiter='\t'))
def write(path,rows):
 with path.open('x') as f:
  w=csv.DictWriter(f,fieldnames=list(rows[0]),delimiter='\t',lineterminator='\n');w.writeheader();w.writerows(rows)
def pct(q,c): return 100*(c/q-1)
def main():
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('input',type=Path);p.add_argument('output',type=Path);args=p.parse_args()
 samples=read(args.input/'build104-samples.tsv')
 if len(samples)!=100 or [int(r['run']) for r in samples]!=list(range(1,101)): raise ValueError('incomplete/duplicate/reordered slots')
 means,pairs,controls,decisions=[],[],[],[]
 for session in ('1','2'):
  current_means={};current_pairs={}
  for metric in METRICS:
   rows=[r for r in samples if r['session']==session and r['metric']==metric]
   if len(rows)!=10 or [int(r['slot']) for r in rows]!=list(range(1,11)) or [r['variant'] for r in rows]!=ORDER: raise ValueError('protocol order differs')
   comparison=[r for r in rows if r['kind']=='comparison'];control=[r for r in rows if r['kind']=='control']
   if len(comparison)!=8 or len(control)!=2 or [int(r['slot']) for r in control]!=[5,6]: raise ValueError('comparison/control labels differ')
   if any(int(r['iterations'])<1 or float(r['ns_per_op'])<=0 for r in rows): raise ValueError('invalid benchmark sample')
   q=[r for r in comparison if r['variant']=='Q'];c=[r for r in comparison if r['variant']=='C']
   qtime=[float(r['ns_per_op']) for r in q];ctime=[float(r['ns_per_op']) for r in c]
   qmean,cmean=statistics.mean(qtime),statistics.mean(ctime)
   qb=statistics.mean(float(r['bytes_per_op']) for r in q);cb=statistics.mean(float(r['bytes_per_op']) for r in c)
   qa=statistics.mean(float(r['allocs_per_op']) for r in q);ca=statistics.mean(float(r['allocs_per_op']) for r in c)
   memory_ok=qa==ca and (abs(cb-qb)<=qb*.01 if metric=='Pool8' else qb==cb)
   change=pct(qmean,cmean)
   m=dict(session=session,metric=metric,q_n=4,c_n=4,q_mean_ns=f'{qmean:.3f}',c_mean_ns=f'{cmean:.3f}',change_pct=f'{change:.6f}',q_min_ns=f'{min(qtime):.3f}',q_max_ns=f'{max(qtime):.3f}',c_min_ns=f'{min(ctime):.3f}',c_max_ns=f'{max(ctime):.3f}',q_mean_bytes=f'{qb:.3f}',c_mean_bytes=f'{cb:.3f}',q_mean_allocs=f'{qa:.3f}',c_mean_allocs=f'{ca:.3f}',allocation_gate=str(memory_ok).lower())
   means.append(m);current_means[metric]=(change,memory_ok)
   favorable=True
   for pair in ('1','2','3','4'):
    group=[r for r in comparison if r['pair']==pair]
    if len(group)!=2 or {r['variant'] for r in group}!={'Q','C'}: raise ValueError('invalid pair')
    qrow=next(r for r in group if r['variant']=='Q');crow=next(r for r in group if r['variant']=='C')
    d=pct(float(qrow['ns_per_op']),float(crow['ns_per_op']));ok=d<0;favorable &= ok
    pairs.append(dict(session=session,metric=metric,pair=pair,order='-'.join(r['variant'] for r in group),q_run=qrow['run'],c_run=crow['run'],q_ns=qrow['ns_per_op'],c_ns=crow['ns_per_op'],change_pct=f'{d:.6f}',favorable=str(ok).lower()))
   current_pairs[metric]=favorable
   controls.append(dict(session=session,metric=metric,first_run=control[0]['run'],second_run=control[1]['run'],first_ns=control[0]['ns_per_op'],second_ns=control[1]['ns_per_op'],apparent_change_pct=f"{pct(float(control[0]['ns_per_op']),float(control[1]['ns_per_op'])):.6f}"))
  stage_mean=all(current_means[m][0]<=-3 for m in ('Generation4','Pool8'))
  stage_pairs=all(current_pairs[m] for m in ('Generation4','Pool8'))
  readers=all(current_means[m][0]<=0 for m in METRICS[:3])
  memory=all(value[1] for value in current_means.values())
  passed=stage_mean and stage_pairs and readers and memory
  decisions.append(dict(session=session,gen4_at_least_3pct=str(current_means['Generation4'][0]<=-3).lower(),pool_at_least_3pct=str(current_means['Pool8'][0]<=-3).lower(),all_stage_pairs_favorable=str(stage_pairs).lower(),reader_means_nonregressive=str(readers).lower(),allocation_gates=str(memory).lower(),authorizes_retained_comparison=str(passed).lower()))
 args.output.mkdir()
 for name,rows in (('samples',samples),('means',means),('pairs',pairs),('controls',controls),('decision',decisions)): write(args.output/f'build104-{name}.tsv',rows)
 passed=all(r['authorizes_retained_comparison']=='true' for r in decisions)
 print('Build104 complete benchmark data: PASS 100/100')
 print('Predeclared performance gate: '+('PASS; later retained comparison authorized, no promotion' if passed else 'FAIL; close not promoted; no private campaign'))
if __name__=='__main__': main()
