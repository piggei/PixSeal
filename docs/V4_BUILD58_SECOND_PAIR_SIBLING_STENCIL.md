# Format-v4 Build58 — sibling stencil after Build57 continuation

Build58 is a **research-only optimizer diagnostic**. At this research checkpoint, Build44 was the latest qualified production milestone. No production path changes.

## Evidence entering Build58

The qualified-host Build57 artifact closes B/mild as `second-pair-continuation-qualified-gain`:

- target candidate 10 retains 65 second-pair parents and 154 continuation states;
- 147 target continuation states are qualified and none authenticate;
- the best qualified continuation reaches **26.968 px** post-hoc oracle error, `-7.080 px` from the original 34.049 px seed;
- proposal continues rising after that geometric minimum, while oracle error first worsens to 27.820 px and then changes again on later accepted coordinates;
- B/angle target candidate 18 remains not triggered; its complete continuation bank has no qualified state.

This reproduces the same structural warning previously observed after Build54: immediate Gauss-Seidel acceptance can change the parent geometry before later coordinates are evaluated. Build58 therefore tests coordinate-order dependence again **without changing the proposal score**.

## Question

For every frozen Build57 post-second-pair continuation state, do independent +/-1px siblings evaluated from the identical parent expose additional qualified geometry gain or HMAC recovery that the ordered continuation misses?

Build58 does not select the 26.968 px state. Every Build57 continuation state is treated identically.

## Blind procedure

Build58 reproduces the complete Build57 blind bank unchanged. For **every retained Build57 post-second-pair continuation state**:

1. keep the exact frozen parent geometry and proposal score;
2. evaluate all 8 corner coordinates in both `-1 px` and `+1 px` directions;
3. every one of the 16 evaluations starts from the **same parent**, so no sibling inherits an earlier coordinate update;
4. retain every sibling whose unchanged proposal score improves by more than `1e-7`;
5. do not use held-out validation, full-pilot qualification, payload, ECC, HMAC, reference geometry or oracle error to generate or retain siblings;
6. freeze the complete B/mild+B/angle Build57 + sibling bank;
7. only after the freeze, annotate held-out/full-pilot qualification and run diagnostic single-candidate HMAC;
8. generate SIFT/reference oracle geometry only after both blind JSON outputs are complete.

The added Build58 work is bounded to at most **16 sibling evaluations per retained Build57 continuation state**.

## Commands

```bash
make v4-build58-phone-second-pair-sibling-stencil-test
make v4-build58-phone-second-pair-sibling-stencil-diagnostic
```

CLI:

```bash
pixseal v4-diagnose-phone-second-pair-sibling \
  -in phone-b-mild.jpg \
  -key "PixSeal-v4-TestKey-2026" \
  -width 1632 -height 1632 -json
```

Private outputs are written under `v4-phone private/build58-diagnostics/`:

- `blind/phone-b-mild-second-pair-sibling-stencil.json`
- `blind/phone-b-angle-second-pair-sibling-stencil.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build58-second-pair-sibling-parents.tsv`
- `build58-second-pair-sibling-states.tsv`
- `build58-second-pair-sibling-summary.tsv`
- `build58-second-pair-sibling-stencil.md`

## Interpretation

For the primary B/mild target:

- `second-pair-sibling-recovery`: at least one frozen sibling authenticates;
- `second-pair-sibling-qualified-gain`: no HMAC, but at least one frozen sibling improves both proposal and post-hoc oracle error relative to its exact Build57 continuation parent and passes unchanged qualification;
- `second-pair-sibling-geometric-gain`: proposal and oracle improve but the sibling does not qualify;
- `second-pair-sibling-proposal-only`: proposal-improving siblings exist but none improve independent geometry;
- `second-pair-sibling-not-triggered`: no retained Build57 continuation state emits a proposal-improving sibling.

Only HMAC-authenticated blind recovery can justify a later production candidate. Non-HMAC outcomes remain research evidence.

## Frozen invariants

Build58 does not change the encoder, Format-v4 wire format, public pilot, strength 48, ECC/Hamming/list decoder, whitening/HMAC domains, proposal score, qualification thresholds, production candidate depth/quorum, Build37 scanner path or deterministic Build44 JPEG/toolchain baseline.

## Qualified host result — 2026-09-24

The Go 1.26.0 host closes Build58 as `second-pair-sibling-qualified-gain`.

- B/mild complete sibling bank: **275 states, 271 qualified, 0 authenticated**.
- Target candidate 10 parent: **26.968 px**, proposal `0.276536`.
- Best qualified target sibling: **26.149 px**, proposal `0.285435`, validation `0.200719`, HMAC false.
- Exact target sibling movement: dimension 6 / corner 3 / x `-1 px`.
- B/angle complete sibling bank: **542 states, 0 qualified, 0 authenticated**; target candidate 18 remains not triggered.

This confirms another optimizer-only qualified geometry gain and motivates Build59's proposal-local third pair escape. It does not promote Build58 to production; At that checkpoint, Build44 was the qualified milestone.
