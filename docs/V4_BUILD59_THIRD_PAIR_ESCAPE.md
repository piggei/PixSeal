# Build59 — Third Pair Escape from Proposal-Local Build58 Siblings

## Status

Completed research-only host study. Build44 remains the latest qualified production milestone. Build59 changes no production encoder/decoder path, proposal score, threshold, quorum, pilot, ECC/Hamming, whitening or HMAC domain.

## Evidence entering Build59

The qualified Go 1.26.0 Build58 host artifact closes B/mild as `second-pair-sibling-qualified-gain`:

- complete B/mild sibling bank: 275 states, 271 qualified, 0 HMAC;
- target candidate 10: exact Build57 parent 26.968 px;
- best qualified Build58 sibling anywhere in the target bank: 26.149 px;
- proposal: 0.276536 -> 0.285435;
- HMAC: false;
- B/angle: 542 sibling states, 0 qualified, 0 HMAC.

The 26.149 px state is post-hoc evidence only. It is not selected to guide Build59.

## Blind experiment

For every frozen Build58 second-pair sibling, Build59:

1. evaluates all eight coordinates at both -1 px and +1 px from the exact same sibling parent (16 independent probes);
2. uses only the unchanged public-pilot proposal objective;
3. marks the sibling proposal-local only when none of the 16 probes improves proposal by more than 1e-7;
4. only for proposal-local siblings evaluates the 112 unordered two-coordinate +/-1 px coupled moves inherited from Build53;
5. retains at most eight proposal-improving pair states, ordered only by proposal;
6. freezes the complete B/mild and B/angle geometry banks before qualification or protected-data work.

After the freeze barrier, unchanged held-out/full-pilot qualification and diagnostic single-candidate HMAC are attached. SIFT/reference geometry is generated only after both blind JSON outputs are complete and is used exclusively for post-hoc interpretation.

## Qualified-host result

On Go 1.26.0, B/mild freezes 180 third-pair states; 175 pass the unchanged qualification gate and none authenticates. Target candidate 10 closes as `third-pair-qualified-gain`: the best qualified third-pair state reaches 25.407 px post-hoc oracle error, 1.664 px better than its exact 27.071 px Build58 sibling parent, while proposal is also improved. B/angle target candidate 18 remains `third-pair-not-triggered`; the full B/angle bank contains 914 third-pair states, 0 qualified and 0 authenticated.

The blind B/mild JSON from the original run and the later oracle/report artifact were reconciled after a subsequent rerun accidentally truncated that one JSON file to 0 bytes. Counts and target telemetry match exactly; this transfer artifact does not change the research result.

## Outputs

The private diagnostic writes:

- `blind/phone-b-mild-third-pair-escape.json`
- `blind/phone-b-angle-third-pair-escape.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build59-third-pair-parents.tsv`
- `build59-third-pair-states.tsv`
- `build59-third-pair-summary.tsv`
- `build59-third-pair-escape.md`

## Classification

Post-hoc labels are diagnostic only:

- `third-pair-recovery`: at least one frozen third-pair state authenticates by HMAC;
- `third-pair-qualified-gain`: a frozen third-pair state improves proposal and oracle error relative to its exact Build58 sibling parent and passes unchanged qualification;
- `third-pair-geometric-gain`: proposal and oracle improve but qualification does not pass;
- `third-pair-proposal-only`: pair states exist but no oracle improvement is observed;
- `third-pair-not-triggered`: no pair state is retained on the target.

No classification changes production behavior.
