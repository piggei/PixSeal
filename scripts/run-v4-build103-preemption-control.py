#!/usr/bin/env python3
"""Fresh-process, balanced default/diagnostic frozen Gen4 study."""
import csv
import hashlib
import os
from pathlib import Path
import signal
import subprocess
import sys

stage = Path(os.environ['PIXSEAL_BUILD103_OUTPUT_DIR'])
testbin = stage / 'build103-watermark.test'
base = dict(os.environ, GOMAXPROCS='8', GOGC='100', GOMEMLIMIT='off', GOTRACEBACK='single')
active = None

def terminate(signum, frame):
    raise SystemExit(f'Build103 interrupted by signal {signum}')

signal.signal(signal.SIGTERM, terminate)
signal.signal(signal.SIGINT, terminate)

def invoke(name, test, env, seconds):
    global active
    # Direct file output survives a timeout/signal and exposes preparation progress.
    with (stage / name).open('x') as log:
        active = subprocess.Popen([str(testbin), '-test.run', '^'+test+'$', '-test.v', '-test.count=1', '-test.timeout=0'], env=env, stdin=subprocess.DEVNULL, stdout=log, stderr=subprocess.STDOUT, text=True, start_new_session=True)
        try:
            active.wait(timeout=seconds)
            if active.returncode != 0:
                raise RuntimeError(f'{name} failed with exit {active.returncode}; evidence preserved')
        except subprocess.TimeoutExpired as error:
            log.write('\nTIMEOUT\n'); log.flush()
            raise RuntimeError(f'{name} exceeded external timeout; evidence preserved') from error
        finally:
            if active.poll() is None:
                try: os.killpg(active.pid, signal.SIGKILL)
                except ProcessLookupError: pass
            active.wait()
            active = None
    return (stage / name).read_text().splitlines()

def extract(lines, header, row):
    headers = [line.split('\t')[1:] for line in lines if line.startswith(header+'\t')]
    values = [line.split('\t')[1:] for line in lines if line.startswith(row+'\t')]
    if len(headers) != 1 or not values or any(len(v) != len(headers[0]) for v in values):
        raise RuntimeError('missing or invalid study rows')
    return headers[0], values

def write_table(path, columns, rows):
    with path.open('x') as f:
        writer = csv.writer(f, delimiter='\t', lineterminator='\n')
        writer.writerow(columns); writer.writerows(rows)

print('build103-start\tprepare-snapshots', flush=True)
prep_env = dict(base, GODEBUG='asyncpreemptoff=0')
lines = invoke('build103-prepare-raw.txt', 'TestExperimentalV4Build103PrepareSnapshots', prep_env, 1800)
capture_columns, capture_values = extract(lines, 'build103-capture-header', 'build103-capture')
if len(capture_values) != 2: raise RuntimeError('expected two snapshots')
write_table(stage/'build103-captures.tsv', capture_columns, capture_values)
captures = {row['image']: row for row in (dict(zip(capture_columns, r)) for r in capture_values)}
for row in capture_values: print('build103-capture\t'+'\t'.join(row), flush=True)

run = 0
columns = None
rows = []
for image in ('phone-b-mild.jpg', 'phone-b-angle.jpg'):
    capture = captures[image]
    snapshot = stage / (image+'.gob')
    if hashlib.sha256(snapshot.read_bytes()).hexdigest() != capture['snapshot_sha256']:
        raise RuntimeError('snapshot hash changed before launch')
    for phase in ('timing', 'cpu'):
        order = ['default','async_off','async_off','default','async_off','default','default','async_off']
        if phase == 'cpu': order = ['default','async_off'] if image == 'phone-b-mild.jpg' else ['async_off','default']
        for i, variant in enumerate(order):
            run += 1
            debug = 'asyncpreemptoff='+('1' if variant=='async_off' else '0')
            print(f'build103-start\t{phase}\t{run}\t{image}\t{variant}', flush=True)
            env = dict(base, GODEBUG=debug, PIXSEAL_BUILD103_SNAPSHOT=str(snapshot), PIXSEAL_BUILD103_SNAPSHOT_SHA256=capture['snapshot_sha256'], PIXSEAL_BUILD103_VARIANT=variant, PIXSEAL_BUILD103_PASS=phase, PIXSEAL_BUILD103_RUN=str(run), PIXSEAL_BUILD103_BLOCK=str(i//4+1), PIXSEAL_BUILD103_PAIR=str(i//2+1), PIXSEAL_BUILD103_POSITION=str(i%2+1), PIXSEAL_BUILD103_CPU_PROFILE=str(stage/f'run{run:02d}-{image}-{variant}.cpu.prof'))
            lines = invoke(f'run{run:02d}-{image}-{variant}-raw.txt', 'TestExperimentalV4Build103ReplayProcess', env, 600)
            names, vals = extract(lines,'build103-header','build103-row')
            if columns is None: columns = names; print('build103-header\t'+'\t'.join(names),flush=True)
            if names != columns or len(vals) != 1: raise RuntimeError('inconsistent replay schema')
            row = dict(zip(names,vals[0]))
            if row['variant'] != variant or row['godebug'] != debug or row['semantic_gate'] != 'PASS' or row['image'] != image or row['snapshot_sha256'] != capture['snapshot_sha256']:
                raise RuntimeError('runtime label or snapshot identity mismatch')
            for field in ('task_sha256','plane_sha256','bank_sha256'):
                if row[field] != capture[field]: raise RuntimeError('cross-process frozen fingerprint mismatch')
            rows.append(vals[0]);print('build103-row\t'+'\t'.join(vals[0]),flush=True)
if run != 20: raise RuntimeError('missing study runs')
write_table(stage/'build103-runs.tsv', columns, rows)
print('Build103 fresh-process semantic control: PASS 20/20',flush=True)
