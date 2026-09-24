# Format-v4 Build57 — continuation after Build56 second-pair escape

Build57 is a **research-only optimizer diagnostic**. Build44 remains the latest qualified production milestone. No production path changes.

## Evidence entering Build57

The qualified-host Build56 diagnostic closes the primary `B/mild` target as `sibling-pair-qualified-gain`:

- candidate 10 / `top+left` / pair rank 3 / seed rank 3;
- 91 target Build55 siblings, 16 proposal-local siblings and 65 retained second-pair states;
- best exact second-pair parent: root 1 / first-pair rank 1 / continuation state 7 / sibling rank 3 / second-pair rank 2;
- coupled move: dimension 2 `-1 px` + dimension 4 `+1 px`;
- proposal: `0.266832 -> 0.273793`;
- post-hoc oracle error: `29.778 -> 28.514 px`;
- held-out/full-pilot qualification: true;
- HMAC: false.

Across the complete Build56 blind banks, B/mild freezes 68 second-pair states (51 qualified, 0 authenticated). B/angle freezes 176 second-pair states (0 qualified, 0 authenticated). The B/angle target candidate 18 remains `sibling-pair-not-triggered`.

Build56 therefore proves that the unchanged proposal objective still exposes a qualified direction beyond the Build55 one-coordinate local maximum. The next optimizer-only question is what happens **after** each retained second-pair escape.

## Question

If every retained Build56 second-pair state is continued under the same bounded 1px proposal-only Gauss-Seidel rule, does the bank contain additional qualified geometry gain or an authenticated state?

Build57 deliberately does not select the 28.514 px state using oracle or held-out evidence. Every retained Build56 second-pair state is treated identically.

## Blind procedure

Build57 reproduces the complete Build56 blind bank unchanged. For **every retained Build56 second-pair state**:

1. start from the exact frozen second-pair geometry and its unchanged proposal score;
2. run bounded 1px proposal-only coordinate descent over all 8 corner coordinates;
3. use at most eight complete passes;
4. retain **every accepted proposal-improving intermediate**, not only the greedy endpoint;
5. do not use held-out validation, full-pilot qualification, payload, ECC, HMAC, reference geometry or oracle error to generate, accept, stop or rank continuation states;
6. freeze the complete B/mild+B/angle second-pair + continuation bank;
7. only after the freeze, annotate held-out/full-pilot qualification and run diagnostic single-candidate HMAC;
8. generate SIFT/reference oracle geometry only after both blind JSON outputs are complete.

Each second-pair branch is bounded to at most 64 accepted continuation states (8 coordinates x 8 passes), with at most 16 proposal evaluations per complete pass.

## Commands

```bash
make v4-build57-phone-sibling-pair-continuation-test
make v4-build57-phone-sibling-pair-continuation-diagnostic
```

CLI:

```bash
pixseal v4-diagnose-phone-sibling-pair-continue \
  -in phone-b-mild.jpg \
  -key "PixSeal-v4-TestKey-2026" \
  -width 1632 -height 1632 -json
```

Private outputs are written under `v4-phone private/build57-diagnostics/`:

- `blind/phone-b-mild-sibling-pair-continuation.json`
- `blind/phone-b-angle-sibling-pair-continuation.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build57-second-pair-parents.tsv`
- `build57-continuation-states.tsv`
- `build57-second-pair-continuation-summary.tsv`
- `build57-second-pair-continuation.md`

## Interpretation

For the primary B/mild target:

- `second-pair-continuation-recovery`: at least one frozen continuation state authenticates;
- `second-pair-continuation-qualified-gain`: no HMAC, but at least one continuation state improves both proposal and post-hoc oracle error relative to its exact Build56 second-pair parent and passes unchanged qualification;
- `second-pair-continuation-geometric-gain`: proposal and oracle improve but the continuation state does not qualify;
- `second-pair-continuation-proposal-only`: proposal-improving continuation states exist but none improve independent geometry;
- `second-pair-continuation-not-triggered`: no retained second-pair state emits an accepted continuation state.

Only HMAC-authenticated blind recovery can justify a later production candidate. Non-HMAC outcomes remain research evidence.

## Frozen invariants

Build57 does not change the encoder, Format-v4 wire format, public pilot, strength 48, ECC/Hamming/list decoder, whitening/HMAC domains, proposal score, qualification thresholds, production candidate depth/quorum, Build37 scanner path or deterministic Build44 JPEG/toolchain baseline.

## Qualified-host result

The Go 1.26.0 host run closes Build57 as `second-pair-continuation-qualified-gain` on B/mild:

- target candidate 10 retains 65 Build56 second-pair parents and 154 Build57 continuation states;
- 147 of the 154 target continuation states pass unchanged qualification; none authenticate;
- the best qualified continuation state is root 1 / first-pair rank 2 / parent state 5 / sibling rank 3 / second-pair rank 3 / continuation index 9;
- proposal rises from the exact second-pair parent while post-hoc oracle error reaches **26.968 px**;
- this is `-2.581 px` relative to that exact Build56 second-pair parent and `-7.080 px` relative to the original 34.049 px seed;
- a later accepted state on the same branch raises proposal further but moves oracle error back to 27.820 px, demonstrating another order-dependent divergence between greedy proposal ascent and independent geometry;
- B/angle target candidate 18 remains `second-pair-continuation-not-triggered`; across the complete B/angle bank 347 continuation states are produced, none qualified or authenticated.

No Build57 state authenticates, so Build44 remains the latest qualified production milestone.
