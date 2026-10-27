#!/usr/bin/env python3
"""Derive public Build103 tables; profiling rows never enter timing summaries."""
import argparse
import csv
from pathlib import Path
import statistics


def read(path):
    with path.open() as f:
        return list(csv.DictReader(f, delimiter='\t'))


def write(path, rows):
    with path.open('x') as f:
        w = csv.DictWriter(f, fieldnames=list(rows[0]), delimiter='\t', lineterminator='\n')
        w.writeheader()
        w.writerows(rows)


def change(default, off):
    return 100 * (off / default - 1)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('input', type=Path, help='completed private evidence directory')
    parser.add_argument('output', type=Path, help='new output directory, must not exist')
    args = parser.parse_args()
    captures = read(args.input / 'build103-captures.tsv')
    rows = read(args.input / 'build103-runs.tsv')
    cases = ('phone-b-mild.jpg', 'phone-b-angle.jpg')
    if len(captures) != 2 or {r['image'] for r in captures} != set(cases):
        raise ValueError('expected two distinct captures')
    if len(rows) != 20 or {int(r['run']) for r in rows} != set(range(1,21)) or len({r['pid'] for r in rows}) != 20:
        raise ValueError('expected twenty separate replay processes')
    references = {r['image']: r for r in captures}
    for r in rows:
        c = references[r['image']]
        if r['variant'] not in ('default','async_off') or r['godebug'] != 'asyncpreemptoff='+('1' if r['variant']=='async_off' else '0'):
            raise ValueError('invalid runtime variant')
        if r['semantic_gate'] != 'PASS' or c['semantic_gate'] != 'PASS' or any(r[k] != c[k] for k in ('snapshot_sha256','task_sha256','plane_sha256','bank_sha256','bank','geometry_evals')) or r['task_count'] != c['tasks']:
            raise ValueError('frozen semantics or fingerprints differ')
        if r['workers'] != '8' or int(r['clock_ticks_per_s']) <= 0:
            raise ValueError('runtime worker/clock protocol differs')
        r['process_cpu_ms'] = f"{float(r['process_user_cpu_ms']) + float(r['process_system_cpu_ms']):.3f}"
    pairs, means, decisions = [], [], []
    metrics = ('gen4_pool_ms','process_cpu_ms','replay_ms','gen4_worker_ms')
    for image in cases:
        timing = [r for r in rows if r['image']==image and r['pass']=='timing']
        profiled = [r for r in rows if r['image']==image and r['pass']=='cpu']
        timing.sort(key=lambda r:int(r['run']))
        expected = ['default','async_off','async_off','default','async_off','default','default','async_off']
        profile_order = ['default','async_off'] if image==cases[0] else ['async_off','default']
        if [r['variant'] for r in timing] != expected or [r['variant'] for r in profiled] != profile_order:
            raise ValueError('incomplete or reordered timing/profile protocol')
        favorable = True
        for i in range(4):
            group = timing[i*2:i*2+2]
            d = next(r for r in group if r['variant']=='default')
            o = next(r for r in group if r['variant']=='async_off')
            if any(int(r['pair']) != i+1 for r in group):
                raise ValueError('pair identifier differs')
            p = dict(image=image, pair=i+1, order='D-O' if group[0]['variant']=='default' else 'O-D', default_run=d['run'], async_off_run=o['run'])
            for metric in metrics:
                dv, ov = float(d[metric]),float(o[metric])
                if dv <= 0 or ov <= 0: raise ValueError('nonpositive timing')
                p['default_'+metric] = f'{dv:.3f}'
                p['async_off_'+metric] = f'{ov:.3f}'
                p[metric+'_change_pct'] = f'{change(dv,ov):.6f}'
            pair_ok = float(o['gen4_pool_ms']) < float(d['gen4_pool_ms']) and float(o['process_cpu_ms']) < float(d['process_cpu_ms'])
            p['pool_and_cpu_favorable'] = str(pair_ok).lower()
            favorable &= pair_ok
            pairs.append(p)
        improvements = {}
        for metric in metrics:
            d = [float(r[metric]) for r in timing if r['variant']=='default']
            o = [float(r[metric]) for r in timing if r['variant']=='async_off']
            dm,om = statistics.mean(d),statistics.mean(o)
            pct = change(dm,om)
            improvements[metric] = -pct
            means.append(dict(image=image,metric=metric,default_n=len(d),async_off_n=len(o),default_mean=f'{dm:.3f}',async_off_mean=f'{om:.3f}',change_pct=f'{pct:.6f}',default_min=f'{min(d):.3f}',default_max=f'{max(d):.3f}',async_off_min=f'{min(o):.3f}',async_off_max=f'{max(o):.3f}',default_population_cv_pct=f'{100*statistics.pstdev(d)/dm:.6f}',async_off_population_cv_pct=f'{100*statistics.pstdev(o)/om:.6f}'))
        wall_ok, cpu_ok = improvements['gen4_pool_ms'] >= 5, improvements['process_cpu_ms'] >= 5
        decisions.append(dict(image=image,pool_improvement_pct=f"{improvements['gen4_pool_ms']:.6f}",cpu_improvement_pct=f"{improvements['process_cpu_ms']:.6f}",pool_at_least_5pct=str(wall_ok).lower(),cpu_at_least_5pct=str(cpu_ok).lower(),all_four_pairs_favorable=str(favorable).lower(),case_supports_material_effect=str(wall_ok and cpu_ok and favorable).lower()))
    supported = all(r['case_supports_material_effect']=='true' for r in decisions)
    args.output.mkdir()
    write(args.output/'build103-captures.tsv',captures)
    write(args.output/'build103-runs.tsv',rows)
    write(args.output/'build103-pairs.tsv',pairs)
    write(args.output/'build103-means.tsv',means)
    write(args.output/'build103-decision.tsv',decisions)
    print('Build103 semantic gate: PASS 20/20')
    print('Predeclared material diagnostic effect: '+('SUPPORTED' if supported else 'MIXED/UNSUPPORTED'))
    print('Qualified Build98 / previous Build93 / Build100 non-promotion unchanged.')


if __name__ == '__main__':
    main()
