# Format-v4 Build55 — continuation sibling-stencil observability

Build55 is a **research-only optimizer diagnostic**. Build44 remains the latest qualified production milestone. No production path changes.

## Evidence entering Build55

The qualification-host Build54 run confirms `continuation-qualified-gain` on the primary `B/mild` target. The unchanged Build54 continuation bank retains a qualified intermediate at **30.539 px** mean corner error, down from the 32.447 px pair-rank-1 parent and the 34.049 px seed. HMAC remains false. `B/angle` remains `continuation-not-triggered` on its target candidate.

Build54 also exposes a more specific optimizer defect. Its continuation is Gauss-Seidel: when one coordinate improves proposal it is accepted immediately, so later coordinates are evaluated from the modified state. A post-hoc design probe at the frozen 30.539 px state shows that an independent `corner 2 / x +1 px` move improves both proposal (`0.257772 -> 0.266832`) and oracle error (`30.539 -> ~29.778 px`) and remains qualified, but the ordinary continuation reaches another coordinate first. The sibling itself does not authenticate in the design probe.

## Question

Does evaluating every one-step coordinate neighbor **from the same frozen continuation parent**, rather than inheriting earlier coordinate updates, expose additional proposal-improving qualified geometry or HMAC recovery?

## Blind procedure

Build55 reproduces Build54 exactly through the retained continuation bank. Then, for **every** retained continuation state:

1. freeze that continuation state as the parent;
2. evaluate all 8 corner coordinates in both `-1 px` and `+1 px` directions;
3. each of the 16 evaluations starts from the exact same parent geometry;
4. retain every sibling whose unchanged proposal score improves by more than `1e-7`;
5. do not rank or prune siblings using held-out validation, full-pilot detection, payload, ECC, HMAC or oracle geometry;
6. freeze the complete continuation+sibling bank;
7. only then annotate held-out/full-pilot qualification and run diagnostic single-candidate HMAC;
8. generate SIFT/reference oracle geometry only after both blind JSON outputs are complete.

The sibling stencil is bounded at 16 proposal evaluations per continuation parent. Build54 itself remains unchanged and reproducible.

## Commands

```bash
make v4-build55-phone-sibling-stencil-test
make v4-build55-phone-sibling-stencil-diagnostic
```

CLI:

```bash
pixseal v4-diagnose-phone-sibling-stencil \
  -in phone-b-mild.jpg \
  -key "PixSeal-v4-TestKey-2026" \
  -width 1632 -height 1632 -json
```

Private outputs are written under `v4-phone private/build55-diagnostics/`:

- `blind/phone-b-mild-sibling-stencil.json`
- `blind/phone-b-angle-sibling-stencil.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build55-sibling-parents.tsv`
- `build55-sibling-states.tsv`
- `build55-sibling-summary.tsv`
- `build55-sibling-stencil.md`

## Interpretation

For the primary B/mild target:

- `sibling-recovery`: at least one frozen sibling authenticates;
- `sibling-qualified-gain`: no HMAC, but at least one sibling improves proposal and post-hoc oracle error relative to its exact parent and passes unchanged qualification;
- `sibling-geometric-gain`: proposal and oracle improve, but the sibling does not qualify;
- `sibling-proposal-only`: proposal-improving siblings exist but none improve independent geometry;
- `sibling-not-triggered`: the target has no continuation parent with proposal-improving siblings.

Only HMAC-authenticated blind recovery can justify a later production candidate. Non-HMAC outcomes remain research evidence.

## Frozen invariants

Build55 does not change the encoder, Format-v4 wire format, public pilot, strength 48, ECC/Hamming/list decoder, whitening/HMAC domains, proposal score, qualification thresholds, production candidate depth/quorum, Build37 scanner path or deterministic Build44 JPEG/toolchain baseline.

## Qualified host result — 2026-09-24

The Go 1.26.0 qualification host confirms the sibling-stencil hypothesis on the primary `B/mild` target (candidate 10, `top+left`, side-pair rank 3, seed rank 3):

- retained target continuation states: 34;
- proposal-improving target siblings: **91**;
- best parent geometry: continuation state 7 of pair-rank 1 at **30.539 px** post-hoc mean corner error;
- best sibling: sibling rank 3, `dimension 4`, `+1 px` (`corner 2 / x +1 px` in the diagnostic coordinate convention);
- proposal: `0.257772 -> 0.266832`;
- post-hoc oracle error: `30.539 -> 29.778 px` (`-0.761 px`);
- held-out validation: `0.163706`;
- qualified: **true**;
- HMAC: **false**.

The target classification is therefore **`sibling-qualified-gain`**. `B/angle` target candidate 18 has no continuation bank and remains **`sibling-not-triggered`**.

Across the complete blind outputs, not only the target candidate:

- B/mild: 93 frozen siblings, 76 qualified, 0 authenticated;
- B/angle: 106 frozen siblings, 0 qualified, 0 authenticated.

This closes Build55 as evidence that coordinate-order bias was real and that the unchanged proposal objective still contains locally useful geometry. It does **not** justify production promotion because no sibling authenticates.

A follow-up post-hoc design probe shows that the 29.778 px sibling is itself a one-coordinate local maximum: none of the complete 16 independent +/-1px single-coordinate moves improves proposal. A bounded coupled two-coordinate stencil around exactly that frozen sibling nevertheless contains proposal-improving qualified states, motivating Build56.
