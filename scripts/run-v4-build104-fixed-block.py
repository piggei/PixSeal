#!/usr/bin/env python3
"""Run the predeclared public fixed8 benchmark; retain every slot durably."""
import csv
import os
from pathlib import Path
import re
import signal
import subprocess
import time

stage=Path(os.environ['PIXSEAL_BUILD104_OUTPUT_DIR'])
base=dict(os.environ,GODEBUG='asyncpreemptoff=0',GOMAXPROCS='8',GOGC='100',GOMEMLIMIT='off',GOTRACEBACK='single')
order=['Q','C','C','Q','Q','Q','C','Q','Q','C']
columns=['run','session','metric','slot','kind','pair','variant','pid','iterations','ns_per_op','bytes_per_op','allocs_per_op']
pattern=re.compile(r'^BenchmarkExperimentalV4Build104Slot-8\s+(\d+)\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op\s*$',re.M)
active=None
def interrupted(signum,frame): raise SystemExit(f'benchmark interrupted by signal {signum}')
signal.signal(signal.SIGTERM,interrupted);signal.signal(signal.SIGINT,interrupted)
with (stage/'build104-samples.tsv').open('x') as table:
 writer=csv.writer(table,delimiter='\t',lineterminator='\n');writer.writerow(columns);table.flush();os.fsync(table.fileno())
 run=0
 for session in (1,2):
  for metric in ('ReaderFront','ReaderMild','ReaderAngle','Generation4','Pool8'):
   for slot,variant in enumerate(order,1):
    run+=1
    print(f'build104-start\t{run}/100\tsession={session}\t{metric}\tslot={slot}\t{variant}',flush=True)
    time.sleep(1)
    env=dict(base,PIXSEAL_BUILD104_METRIC=metric,PIXSEAL_BUILD104_VARIANT=variant)
    raw=stage/f'run{run:03d}-s{session}-{metric}-{slot}-{variant}.txt'
    with raw.open('x') as log:
     active=subprocess.Popen([str(stage/'build104-watermark.test'),'-test.run=^$','-test.bench=^BenchmarkExperimentalV4Build104Slot$','-test.benchmem','-test.benchtime=1s','-test.count=1','-test.timeout=0'],env=env,stdin=subprocess.DEVNULL,stdout=log,stderr=subprocess.STDOUT,start_new_session=True)
     pid=active.pid
     try:
      rc=active.wait(timeout=180)
      if rc: raise RuntimeError(f'{raw.name}: child failed {rc}')
     finally:
      if active.poll() is None:
       try: os.killpg(active.pid,signal.SIGKILL)
       except ProcessLookupError: pass
      active.wait();active=None
    values=pattern.findall(raw.read_text())
    if len(values)!=1: raise RuntimeError('missing/ambiguous benchmark result: '+raw.name)
    n,ns,byte_count,alloc_count=values[0]
    if int(n)<1 or float(ns)<=0: raise RuntimeError('nonpositive benchmark result')
    kind='control' if slot in (5,6) else 'comparison'
    pair={1:1,2:1,3:2,4:2,5:0,6:0,7:3,8:3,9:4,10:4}[slot]
    writer.writerow([run,session,metric,slot,kind,pair,variant,pid,n,ns,byte_count,alloc_count]);table.flush();os.fsync(table.fileno())
    print(f'build104-done\t{run}\t{metric}\t{variant}\t{ns}\tns/op\t{byte_count}\tB/op\t{alloc_count}\tallocs/op',flush=True)
 if run!=100: raise RuntimeError('incomplete protocol')
print('Build104 public benchmark protocol: COMPLETE 100/100',flush=True)
