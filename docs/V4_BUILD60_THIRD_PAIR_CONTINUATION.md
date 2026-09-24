# Build60 — Post-Third-Pair Continuation

## Status

Research-only. At this research checkpoint, Build44 was the latest qualified production milestone. Build60 changes no production encoder/decoder path, proposal score, threshold, quorum, pilot, ECC/Hamming, whitening or HMAC domain.

## Evidence entering Build60

The qualified Go 1.26.0 Build59 host artifact closes B/mild as `third-pair-qualified-gain`:

- 180 retained third-pair states on B/mild;
- 175 qualified, 0 authenticated;
- target candidate 10 best qualified third-pair state: 25.407 px post-hoc oracle error;
- exact Build58 sibling parent for that winning third-pair state: 27.071 px;
- improvement: 1.664 px while proposal also improves;
- B/angle target candidate 18: `third-pair-not-triggered`; full B/angle third-pair bank 914 states, 0 qualified, 0 authenticated.

The 25.407 px state is evidence only. It is not selected to guide Build60.

## Blind experiment

For every retained Build59 third-pair state, Build60:

1. starts from that exact frozen third-pair geometry;
2. runs the established 1 px proposal-only coordinate descent over all eight corner coordinates;
3. allows at most eight passes and therefore at most 64 accepted intermediate states per branch;
4. retains every accepted intermediate rather than privileging the greedy endpoint;
5. uses no held-out qualification, payload, ECC, key, HMAC or oracle geometry during continuation generation, acceptance, stopping or ranking;
6. freezes the complete B/mild and B/angle geometry banks before downstream annotation.

After the freeze barrier, unchanged held-out/full-pilot qualification and diagnostic single-candidate HMAC are attached. SIFT/reference geometry is generated only after both blind JSON outputs are complete and is used exclusively for post-hoc interpretation.

## Outputs

The private diagnostic writes:

- `blind/phone-b-mild-third-pair-continuation.json`
- `blind/phone-b-angle-third-pair-continuation.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build60-third-pair-parents.tsv`
- `build60-continuation-states.tsv`
- `build60-third-pair-continuation-summary.tsv`
- `build60-third-pair-continuation.md`

## Classification

Post-hoc labels are diagnostic only:

- `third-pair-continuation-recovery`: at least one frozen continuation state authenticates by HMAC;
- `third-pair-continuation-qualified-gain`: a frozen continuation state improves proposal and oracle error relative to its exact Build59 third-pair parent and passes unchanged qualification;
- `third-pair-continuation-geometric-gain`: proposal and oracle improve but qualification does not pass;
- `third-pair-continuation-proposal-only`: continuation states exist but no oracle improvement is observed;
- `third-pair-continuation-not-triggered`: no continuation state is retained on the target.

No classification changes production behavior.

## Qualified host result

The Go 1.26.0 host run closes Build60 as `third-pair-continuation-qualified-gain` on B/mild:

- 180 retained Build59 third-pair parents;
- 406 retained B/mild continuation states;
- 397 qualified, 0 authenticated;
- best qualified continuation: **25.150 px** post-hoc oracle error;
- that best state improves its exact **29.938 px** Build59 parent by **4.788 px** while proposal also improves;
- the original B/mild seed remains 34.048 px.

B/angle target candidate 18 remains `third-pair-continuation-not-triggered`. Across the complete non-target B/angle continuation bank, however, 1467 states are retained and **2 pass qualification**, both around **5664.707 px** from oracle and neither authenticates. This is the first observed qualification leakage in the expanded B/angle optimizer bank. It is not a recovery, but future optimizer builds must report whole-bank qualification counts rather than relying only on the target classification.
