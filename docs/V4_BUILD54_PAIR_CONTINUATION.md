# Format-v4 Build54 — post-pair accepted-state continuation

Build54 is a **research-only optimizer diagnostic**. It does not change the qualified production decoder. Build44 remains the latest qualified milestone.

## Why Build54 exists

The qualified Build53 host run confirmed the Build52 optimizer diagnosis on the primary `B/mild` target:

- seed error: `34.049 px`;
- standard coarse-to-fine baseline: `43.168 px`;
- one 2px root is 1px coordinate-local at `32.442 px` and is held-out-qualified;
- Build53 finds four proposal-improving coupled 1px escapes from that root;
- all four pair states are held-out-qualified;
- the best post-hoc geometric pair state is pair rank 2 at `31.842 px` (`-0.600 px` from its root);
- no pair state authenticates;
- `B/angle` does not trigger pair scanning on its target candidate and remains the distant false-basin control.

Build53 therefore proves that a coordinate-local maximum can be escaped without changing the proposal objective. It does **not** answer what happens after the escape, because pair states are frozen endpoints in Build53.

A local deterministic preflight, used only to design this experiment, continued all four retained `B/mild` pair states with the unchanged proposal objective. Pair rank 1 admits a sequence of ordinary 1px coordinate improvements, all held-out-qualified in the preflight. Retaining intermediates exposes a state at about `30.539 px`, while later proposal ascent begins to move away from the oracle again. No preflight continuation state authenticates. This is motivation, not qualification evidence.

## Hypothesis

After a coupled escape crosses the first coordinate-local barrier, ordinary 1px proposal ascent may again expose useful states. As in Build52, keeping only the greedy endpoint can discard a geometrically better intermediate.

Build54 therefore asks a narrow question:

> If every retained Build53 pair state is continued with bounded proposal-only 1px coordinate descent, does retaining every accepted intermediate expose a qualified state closer to the physical geometry, or an authenticated state, without changing the score?

## Blind geometry procedure

For every unchanged top4-per-side-pair seed:

1. reproduce the Build53 proposal-only 2px accepted-state root path;
2. probe each root with the complete 16-point `+/-1 px` single-coordinate stencil;
3. only at a 1px coordinate-local root, reproduce the Build53 112-state coupled two-coordinate stencil;
4. retain up to eight proposal-improving pair states ranked by proposal only;
5. **continue every retained pair state**, not a post-hoc selected subset;
6. run fixed 1px coordinate descent for at most eight passes;
7. retain **every accepted continuation intermediate**, not only the endpoint;
8. freeze the complete root/pair/continuation geometry bank for both target images;
9. only after the freeze, annotate held-out/full-pilot qualification and diagnostic HMAC;
10. generate SIFT/reference geometry only after both blind JSON outputs exist, for post-hoc measurement only.

The continuation is bounded at eight dimensions x two signs x eight passes per pair state. No qualification, secret key, payload, ECC, HMAC or oracle value participates in geometry generation or ranking.

## Commands

Regression only:

```bash
make v4-build54-phone-pair-continuation-test
```

Private qualification-host study:

```bash
make v4-build54-phone-pair-continuation-diagnostic
```

The CLI diagnostic is:

```bash
pixseal v4-diagnose-phone-pair-continue \
  -in phone-b-mild.jpg \
  -key "PixSeal-v4-TestKey-2026" \
  -width 1632 -height 1632 -json
```

## Private outputs

Private outputs belong under `v4-phone private/build54-diagnostics/` and are not packaged:

- `blind/phone-b-mild-pair-continuation.json`
- `blind/phone-b-angle-pair-continuation.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build54-branches.tsv`
- `build54-continuation-states.tsv`
- `build54-pair-continuation-summary.tsv`
- `build54-pair-continuation.md`

## Interpretation

For the primary B/mild target:

- `continuation-recovery`: at least one frozen continuation state authenticates;
- `continuation-qualified-gain`: no HMAC, but a continuation state improves proposal and post-hoc oracle error relative to its frozen parent pair state and passes unchanged qualification;
- `continuation-geometric-gain`: proposal and oracle improve but the state does not qualify;
- `continuation-no-oracle-gain`: accepted continuation states exist but none improve independent geometry;
- `continuation-no-accepted-state`: a pair branch exists but 1px continuation accepts nothing;
- `continuation-not-triggered`: no pair branch exists for the target.

Only `continuation-recovery` would justify designing a **separate**, fully gated production candidate. Any non-HMAC result remains research evidence only.

## Frozen invariants

Build54 does not change:

- encoder or Format-v4 wire format;
- locked public pilot;
- strength 48;
- ECC/Hamming/list-decode behavior;
- whitening/HMAC domains;
- Build41 proposal objective;
- Build43/44 production candidate depth, side-pair limits, qualification thresholds or quorum;
- deterministic JPEG compatibility path.

Build44 remains the production qualification baseline.
