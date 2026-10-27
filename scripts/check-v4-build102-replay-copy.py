#!/usr/bin/env python3
"""Verify the split test-only bank against the source-checked Build101 copy."""
from pathlib import Path
import subprocess
import sys

root = Path(__file__).resolve().parent.parent
src = (root / 'watermark/experimental_v4_build101_observed_bank_test.go').read_text()
body = src[src.index('\tvar profile'):]
body = '\n'.join(line for line in body.split('\n') if not line.startswith('\tobserver.'))
prefix, tail = body.split('\tgen4Workers := runtime.GOMAXPROCS(0)', 1)
pack = 'experimentalV4Build102Pack(plane, pilot, cw, ch, anchor, nil, evals, {seeds}, {workers}, profile, freezeProfile)'
prefix = prefix.replace('return nil, 0, 0, 0, profile, freezeProfile', 'return experimentalV4Build102Frozen{}')
prefix = prefix.replace('return nil, evals, 0, 0, profile, freezeProfile', 'return '+pack.format(seeds='0', workers='0'))
prefix = prefix.replace('return nil, evals, len(seeds), workers, profile, freezeProfile', 'return '+pack.format(seeds='len(seeds)', workers='workers'))
expected = '''func experimentalV4Build102Prepare(work image.Image, boundary PrintBoundaryEstimate, cw, ch int, candidate bool) experimentalV4Build102Frozen {
'''+prefix+'\treturn '+pack.replace('anchor, nil,', 'anchor, flat4,').format(seeds='len(seeds)', workers='workers')+'\n}\n\n'
expected += '''func experimentalV4Build102Replay(f experimentalV4Build102Frozen) ([]experimentalV4PhoneHypothesis, int, int, int, experimentalV4PhoneBuild100GeometryTelemetry, experimentalV4PhoneBuild75FreezeTelemetry) {
    profile, freezeProfile := f.profile, f.freeze
    plane, pilot, cw, ch, anchor, flat4 := f.plane, f.pilot, f.cw, f.ch, f.anchor, f.tasks
    evals, workers := f.evals, f.workers
    if len(flat4) == 0 { return nil, evals, f.seeds, workers, profile, freezeProfile }
    var wg sync.WaitGroup
    var started time.Time
\tgen4Workers := runtime.GOMAXPROCS(0)'''+tail.replace('len(seeds)', 'f.seeds')
def formatted(body):
    return subprocess.run(['gofmt'],input='package watermark\n'+body,text=True,capture_output=True,check=True).stdout
path = root / 'watermark/experimental_v4_build102_split_bank_test.go'
if '--write' in sys.argv:
    path.write_text('''package watermark

// Test-only split at the Gen4 barrier. Both pieces retain pinned kernels,
// queues, original task order and ordered result collection.
import ("image"; "runtime"; "sort"; "sync"; "time")

'''+expected)
actual = path.read_text();actual = actual[actual.index('func experimentalV4Build102Prepare'):]
if formatted(actual) != formatted(expected):
    raise SystemExit('FAIL: Build102 split bank drifted from verified Build101 bodies')
print('Build102 split-bank structural check: PASS')
