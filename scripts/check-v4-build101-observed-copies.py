#!/usr/bin/env python3
"""Check test-only observation copies against pinned Build100 bodies."""
from pathlib import Path
import subprocess

root = Path(__file__).resolve().parent.parent
w = root / "watermark"
bank = (w / "experimental_v4_phone_build100_bank.go").read_text()
bank = bank[bank.index("func experimentalV4PhoneBuild100BlindBank"):]
bank = bank.replace("experimentalV4PhoneBuild100BlindBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int)", "experimentalV4Build101ObservedBank(work image.Image, boundary PrintBoundaryEstimate, cw, ch int, candidate bool, observer *experimentalV4Build101Observer)")
bank = bank.replace("\tfrozen, _, evals, freezeProfile, pairProfile := experimentalV4PhoneBuild100FreezeParallel(work, boundary, cw, ch, 128)\n\tprofile.PairScoring = pairProfile", '''\tobserver.begin("freeze")
\tvar frozen []experimentalV4PhoneBuild47Frozen
\tvar evals int
\tif candidate {
\t\tfrozen, _, evals, freezeProfile, profile.PairScoring = experimentalV4PhoneBuild100FreezeParallel(work, boundary, cw, ch, 128)
\t} else {
\t\tfrozen, _, evals, freezeProfile = experimentalV4PhoneBuild75FreezeParallel(work, boundary, cw, ch, 128)
\t\tprofile.PairScoring = experimentalV4PhoneBuild100PairTelemetry{Workers: 1, Tasks: freezeProfile.PairScoreTasks, Evaluations: freezeProfile.PairScoreEvaluations, WallElapsed: freezeProfile.PairScoreElapsed, WorkerElapsed: freezeProfile.PairScoreElapsed}
\t}
\tobserver.end()''')
bank = bank.replace("\tworkers := runtime.GOMAXPROCS(0)", '\tobserver.begin("prefix1")\n\tworkers := runtime.GOMAXPROCS(0)', 1)
bank = bank.replace("\tprofile.Prefix1WallElapsed = time.Since(started)", "\tprofile.Prefix1WallElapsed = time.Since(started)\n\tobserver.end()")
for g in (2, 3, 4):
    bank = bank.replace(f"\tgen{g}Workers := runtime.GOMAXPROCS(0)", f'\tobserver.begin("gen{g}")\n\tgen{g}Workers := runtime.GOMAXPROCS(0)')
    bank = bank.replace(f"\tprofile.Gen{g}WallElapsed = time.Since(started)", f"\tprofile.Gen{g}WallElapsed = time.Since(started)\n\tobserver.end()")
recover = (w / "experimental_v4_phone_build100_recover.go").read_text()
recover = recover[recover.index("func experimentalV4PhoneBuild100Recover"):recover.index("func experimentalV4PhoneBuild100ApplyTelemetry")]
recover = recover.replace("experimentalV4PhoneBuild100Recover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int)", "experimentalV4Build101ObservedRecover(work image.Image, boundary PrintBoundaryEstimate, key []byte, cw, ch int, candidate bool, observer *experimentalV4Build101Observer)")
recover = recover.replace("experimentalV4PhoneBuild100BlindBank(work, boundary, cw, ch)", "experimentalV4Build101ObservedBank(work, boundary, cw, ch, candidate, observer)")
recover = recover.replace("\tqualificationStarted := time.Now()", '\tobserver.begin("qualification")\n\tqualificationStarted := time.Now()')
recover = recover.replace("\ttelemetry.QualificationElapsed = time.Since(qualificationStarted)", "\ttelemetry.QualificationElapsed = time.Since(qualificationStarted)\n\tobserver.end()")
recover = recover.replace("\tdecodeStarted := time.Now()", '\tobserver.begin("decode")\n\tdefer observer.end()\n\tdecodeStarted := time.Now()')
for name, expected in (("bank", bank), ("recover", recover)):
    actual = (w / f"experimental_v4_build101_observed_{name}_test.go").read_text()
    actual = actual[actual.index("func experimentalV4Build101Observed"):]
    def formatted(body):
        return subprocess.run(["gofmt"], input="package watermark\n"+body, text=True, capture_output=True, check=True).stdout
    if formatted(actual) != formatted(expected):
        raise SystemExit(f"FAIL: observed {name} drifted from pinned source plus declared probes")
print("Build101 observed-copy structural check: PASS")
