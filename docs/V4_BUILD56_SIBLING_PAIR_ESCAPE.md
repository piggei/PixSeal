# Format-v4 Build56 — pair escape from proposal-local Build55 siblings

Build56 is a **research-only optimizer diagnostic**. At this research checkpoint, Build44 was the latest qualified production milestone. No production path changes.

## Evidence entering Build56

The qualified Build55 host artifact closes the primary `B/mild` target as `sibling-qualified-gain`:

- candidate 10 / `top+left` / pair rank 3 / seed rank 3;
- 91 proposal-improving target siblings;
- best retained qualified sibling: 29.778 px post-hoc mean corner error;
- parent continuation state: 30.539 px;
- proposal: `0.257772 -> 0.266832`;
- HMAC: false.

The `B/angle` target remains `sibling-not-triggered`. Across the full Build55 blind banks, B/mild freezes 93 sibling states (76 qualified, 0 authenticated) and B/angle freezes 106 (0 qualified, 0 authenticated).

A deterministic post-hoc design probe around the exact 29.778 px Build55 sibling finds **zero** proposal-improving one-coordinate +/-1px moves. The sibling is therefore a local maximum under the complete Build55 single-coordinate stencil. The same frozen parent nevertheless has five proposal-improving coupled +/-1px two-coordinate moves. All five pass unchanged qualification in the local probe and none authenticates. The best geometric retained probe state is pair rank 2:

- dimension 2: `-1 px`;
- dimension 4: `+1 px`;
- proposal: `0.266832 -> 0.273793`;
- post-hoc oracle error: `29.778 -> ~28.514 px`;
- qualified: true;
- HMAC: false.

This probe is motivation only, not qualification evidence.

## Qualified host result

The Build56 host diagnostic confirms the probe without oracle-guided selection. For the primary B/mild target, 16 of 91 frozen siblings are proposal-local and emit 65 retained second-pair states. The best qualified state is second-pair rank 2 (`dimension 2 -1 px` + `dimension 4 +1 px`): proposal improves `0.266832 -> 0.273793`, post-hoc oracle error improves `29.778 -> 28.514 px`, qualification remains true and HMAC remains false. The target classification is `sibling-pair-qualified-gain`.

Across the complete Build56 bank, B/mild freezes 68 second-pair states (51 qualified, 0 authenticated). B/angle freezes 176 second-pair states (0 qualified, 0 authenticated), while target candidate 18 remains `sibling-pair-not-triggered`. This result motivates Build57 continuation from every retained second-pair state.

## Question

Do proposal-local Build55 siblings hide additional coupled two-coordinate directions under the **same unchanged proposal score**, analogous to the earlier Build53 escape but deeper in the optimizer trajectory?

## Blind procedure

Build56 reproduces the complete Build55 blind bank. For **every frozen Build55 sibling**:

1. evaluate all 8 corner coordinates in both `-1 px` and `+1 px` directions from that exact sibling parent;
2. count proposal-improving single-coordinate neighbors using only the unchanged proposal fold;
3. if at least one single-coordinate neighbor improves proposal, do not run a coupled scan from that sibling;
4. if zero single-coordinate neighbors improve proposal, mark the sibling proposal-local and evaluate all 112 deterministic coupled +/-1px two-coordinate combinations;
5. retain at most eight proposal-improving pair states, ordered only by proposal score;
6. do not use held-out validation, full-pilot detection, payload, ECC, HMAC, reference geometry or oracle error in locality classification, pair generation or pair ranking;
7. freeze the complete root/pair/continuation/sibling/second-pair bank for both B/mild and B/angle;
8. only then annotate held-out/full-pilot qualification and run diagnostic single-candidate HMAC;
9. generate SIFT/reference oracle geometry only after both blind JSON outputs are complete.

The additional search is bounded at 16 single-coordinate proposal evaluations per Build55 sibling and, only for proposal-local siblings, 112 coupled pair evaluations with at most eight retained pair states.

## Commands

```bash
make v4-build56-phone-sibling-pair-escape-test
make v4-build56-phone-sibling-pair-escape-diagnostic
```

CLI:

```bash
pixseal v4-diagnose-phone-sibling-pair-escape \
  -in phone-b-mild.jpg \
  -key "PixSeal-v4-TestKey-2026" \
  -width 1632 -height 1632 -json
```

Private outputs are written under `v4-phone private/build56-diagnostics/`:

- `blind/phone-b-mild-sibling-pair-escape.json`
- `blind/phone-b-angle-sibling-pair-escape.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build56-sibling-local.tsv`
- `build56-pair-states.tsv`
- `build56-sibling-pair-summary.tsv`
- `build56-sibling-pair-escape.md`

## Interpretation

For the primary B/mild target:

- `sibling-pair-recovery`: at least one frozen second pair state authenticates;
- `sibling-pair-qualified-gain`: no HMAC, but at least one second pair state improves both proposal and post-hoc oracle error relative to its exact Build55 sibling parent and passes unchanged qualification;
- `sibling-pair-geometric-gain`: proposal and oracle improve but the pair state does not qualify;
- `sibling-pair-proposal-only`: proposal-improving second pair states exist but none improve independent geometry;
- `sibling-pair-not-triggered`: no target sibling is proposal-local with a retained pair escape.

Only HMAC-authenticated blind recovery can justify a later production candidate. Non-HMAC outcomes remain research evidence.

## Frozen invariants

Build56 does not change the encoder, Format-v4 wire format, public pilot, strength 48, ECC/Hamming/list decoder, whitening/HMAC domains, proposal score, qualification thresholds, production candidate depth/quorum, Build37 scanner path or deterministic Build44 JPEG/toolchain baseline.
