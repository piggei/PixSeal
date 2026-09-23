# Format-v4 Build46 — qualified-geometry handoff diagnostic

Build46 is a research-only diagnostic checkpoint. **Build44 remains the latest qualified production milestone.** No production geometry threshold, data decoder, Format-v4 field, encoder parameter, pilot, ECC/Hamming rule, whitening/HMAC domain or strength value changes in this build.

## Starting evidence

The qualification-host Build45 study produced the same blind pattern on both difficult marked-B captures:

| image | Build43 frozen | Build43 held-out-qualified | production Build42 bank | HMAC |
|---|---:|---:|---:|---:|
| `phone-b-mild.jpg` | 32 | 1 | 0 | false |
| `phone-b-angle.jpg` | 28 | 1 | 0 | false |

The isolated reference-assisted oracle authenticates `v4-b38-phone-b` on the first Build42 list frame for both captures. Therefore the protected strength-48 data channel is demonstrably viable when geometry is supplied independently.

The important refinement is that `Build43QualifiedCandidates=1` does **not** imply that a candidate is lost by an implementation bug. The production phone decoder requires two qualified geometries for its direct ensemble. Only when that ensemble exists is the Build43 bank promoted downstream; Build42's list-bank decoder separately requires at least three qualified geometries. A singleton therefore stops before both production data paths by design.

## Build46 question

For each held-out-qualified Build43 singleton:

1. does the candidate authenticate if decoded alone **after** geometry is already frozen and qualified; and
2. how far is its quadrilateral from the independently supplied oracle geometry?

These measurements distinguish two materially different cases:

- `qualified-ensemble-shortfall`: the singleton itself carries a valid frame, but the production quorum correctly prevents one geometry from becoming a production acceptance path;
- `qualified-candidate-mismatch`: the singleton passes held-out pilot qualification but still cannot authenticate, indicating that the surviving geometry is not precise enough for the protected data plane.

Other diagnostic labels are `no-build43-candidate`, `heldout-qualification`, `build42-bank-shortfall`, `qualified-handoff-available`, and `build41-accepted`.

## Security / methodology separation

`v4-diagnose-phone-handoff` reruns Build41/43 blind search using exactly the existing structure/public-pilot evidence. Only after Build43 has frozen and held-out-qualified its candidates does Build46 inspect data margins and final HMAC for the single-candidate laboratory decode.

HMAC therefore cannot:

- create a geometry candidate;
- move a corner or side;
- rank side pairs;
- select a proposal basin;
- change held-out qualification; or
- change any production acceptance threshold.

The optional oracle quadrilateral is also comparison-only. It is never passed into Build41/43 search. The private helper script completes both blind handoff diagnostics first, then—if OpenCV/SIFT and the marked-B reference are available—generates the reference quadrilaterals and computes post-hoc corner distances from the already-written blind JSON.

## Commands

Source regression:

```bash
make v4-build46-phone-handoff-test
```

Private B/mild + B/angle study:

```bash
make v4-build46-phone-handoff-diagnostic
```

Outputs remain private under:

```text
v4-phone private/build46-diagnostics/
```

The main files are:

```text
build46-qualified-handoff.tsv
build46-qualified-handoff.md
blind/phone-b-mild-handoff.json
blind/phone-b-angle-handoff.json
```

If the reference and Python OpenCV/SIFT are available, `oracle/*-quad.json` is generated only after the blind JSON files have been completed.

## Production invariants

Build46 leaves unchanged:

- Format-v3 frozen core;
- Format-v4 encoder and wire format;
- public pilot identity;
- strength 48;
- protected data mapping;
- Hamming(7,4) and Build42 ML list rules;
- whitening/HMAC domains;
- Build37 scanner path;
- Build41 global phone basin search;
- Build43 side-pair candidate generation and held-out qualification;
- production two-geometry direct ensemble requirement;
- production Build42 minimum three-geometry bank;
- Build40 residual fallback; and
- deterministic Build44 JPEG ingest / Go 1.26.0 qualification.
