# Build61 — Post-Third-Pair Continuation Sibling Stencil

## Status

Research-only. At this research checkpoint, Build44 was the latest qualified production milestone. Build61 changes no production encoder/decoder path, proposal score, threshold, quorum, pilot, ECC/Hamming, whitening or HMAC domain.

## Evidence entering Build61

The qualified Go 1.26.0 Build60 host artifact closes B/mild as `third-pair-continuation-qualified-gain`:

- 180 Build59 third-pair parents;
- 406 retained B/mild continuation states;
- 397 qualified, 0 authenticated;
- best qualified continuation: 25.150 px post-hoc oracle error;
- exact parent for that state: 29.938 px; improvement 4.788 px while proposal also improves.

B/angle target candidate 18 remains not triggered. The complete non-target B/angle continuation bank contains 1467 states, of which 2 pass qualification at about 5665 px oracle error and 0 authenticate. This is treated as a qualification-leakage control, not as recovery evidence.

The 25.150 px state is post-hoc evidence only and is not selected to guide Build61.

## Blind experiment

For every retained Build60 post-third-pair continuation state, Build61:

1. freezes that exact parent geometry;
2. evaluates all eight coordinates in both -1 px and +1 px directions from that same parent;
3. retains only siblings whose unchanged proposal score improves by more than 1e-7;
4. does not inherit an accepted update from an earlier coordinate, eliminating Gauss-Seidel order dependence for this stencil;
5. uses no held-out qualification, payload, ECC, key, HMAC or oracle geometry during sibling generation or retention;
6. freezes the complete B/mild and B/angle sibling banks before downstream annotation.

After the freeze barrier, unchanged held-out/full-pilot qualification and diagnostic single-candidate HMAC are attached. SIFT/reference geometry is generated only after both blind JSON outputs are complete and is used exclusively for post-hoc interpretation.

## Additional control

Because Build60 first observed 2 qualified non-target B/angle continuation states, Build61 records total sibling states, total qualified siblings, total authenticated siblings and the whole-bank qualification rate for both images. Any growth in unauthenticated B/angle qualification must be reviewed before further search expansion.

## Outputs

The private diagnostic writes:

- `blind/phone-b-mild-third-pair-sibling-stencil.json`
- `blind/phone-b-angle-third-pair-sibling-stencil.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build61-third-pair-sibling-parents.tsv`
- `build61-third-pair-sibling-states.tsv`
- `build61-third-pair-sibling-summary.tsv`
- `build61-third-pair-sibling-stencil.md`

## Classification

Post-hoc labels are diagnostic only:

- `third-pair-sibling-recovery`: at least one frozen sibling authenticates by HMAC;
- `third-pair-sibling-qualified-gain`: a frozen sibling improves proposal and oracle error relative to its exact Build60 continuation parent and passes unchanged qualification;
- `third-pair-sibling-geometric-gain`: proposal and oracle improve but qualification does not pass;
- `third-pair-sibling-proposal-only`: sibling states exist but no oracle improvement is observed;
- `third-pair-sibling-not-triggered`: no proposal-improving sibling is retained on the target.

No classification changes production behavior.

## Qualified host result

The Go 1.26.0 host run closes Build61 as `third-pair-sibling-qualified-gain`.
B/mild freezes 593 proposal-improving sibling states, 591 qualify and 0 authenticate.
The best qualified sibling reaches **24.614 px** from its exact **25.359 px**
Build60 continuation parent, improving oracle error by **0.745 px** while proposal
rises to `0.288064`. The separate best Build60 continuation parent is 25.150 px
and is not the genealogy parent of the winning sibling.

B/angle freezes 1902 siblings with 0 qualified and 0 authenticated. Thus the two
rare non-target qualified continuation states seen in Build60 do not propagate as
qualified Build61 siblings.
