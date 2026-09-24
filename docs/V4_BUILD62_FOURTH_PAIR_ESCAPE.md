# Build62 — Fourth-Pair Escape from Proposal-Local Build61 Siblings

## Status

Completed research diagnostic on the qualified Go 1.26.0 host. At this research checkpoint, Build44 was the latest qualified production milestone. Build62 changes no production encoder/decoder path, proposal score, threshold, quorum, pilot, ECC/Hamming, whitening or HMAC domain.

## Evidence entering Build62

The qualified Go 1.26.0 Build61 artifact closes B/mild as `third-pair-sibling-qualified-gain`:

- 593 proposal-improving sibling states are frozen on B/mild;
- 591 pass unchanged qualification;
- 0 authenticate by HMAC;
- the best qualified sibling reaches **24.614 px** post-hoc oracle error;
- that sibling's exact Build60 continuation parent is **25.359 px**, so the local improvement is **-0.745 px** while proposal rises to `0.288064`.

The separate best Build60 continuation state is 25.150 px; it is not the exact genealogy parent of the 24.614 px sibling and is not used to define the Build62 search.

B/angle freezes 1902 Build61 sibling states with 0 qualified and 0 authenticated. This is encouraging, but Build60 previously observed two unauthenticated qualified states in a larger continuation bank, so Build62 continues to report whole-bank control qualification rates explicitly.

## Hypothesis

A frozen Build61 sibling can be a local maximum under all independent one-coordinate +/-1px moves while a coupled two-coordinate +/-1px move still improves the unchanged proposal objective. This is the same optimizer failure mode previously demonstrated by Build53, Build56 and Build59 at earlier search depths.

## Blind algorithm

For every frozen Build61 third-pair sibling, Build62:

1. evaluates all 16 independent one-coordinate +/-1px neighbors from the identical parent;
2. counts proposal-improving neighbors using the unchanged proposal objective;
3. marks the parent `coordinate_local` only if that count is zero;
4. only for proposal-local parents, evaluates the 112 coupled two-coordinate +/-1px combinations;
5. retains at most eight proposal-improving pair states, ordered only by proposal;
6. freezes the complete B/mild+B/angle bank before qualification, protected-data decode or HMAC;
7. attaches qualification/HMAC only after the freeze barrier;
8. generates SIFT/reference oracle geometry only after both blind JSON files exist.

No secret key, payload, ECC result, HMAC result, qualification result or oracle geometry influences geometry generation, locality classification, acceptance, stopping or ranking.

## Control reporting

Build62 records:

- total Build61 sibling parents;
- total proposal-local sibling parents;
- fourth-pair evaluations and proposal-improving states;
- total retained fourth-pair states;
- whole-bank qualified fourth-pair states;
- whole-bank authenticated fourth-pair states;
- whole-bank qualification rate for B/mild and B/angle.

## Commands

```sh
make v4-build62-phone-fourth-pair-escape-test
make v4-build62-phone-fourth-pair-escape-diagnostic
```

The diagnostic writes under `v4-phone private/build62-diagnostics/`:

- `blind/phone-b-mild-fourth-pair-escape.json`
- `blind/phone-b-angle-fourth-pair-escape.json`
- `oracle/phone-b-mild-quad.json`
- `oracle/phone-b-angle-quad.json`
- `build62-fourth-pair-parents.tsv`
- `build62-fourth-pair-states.tsv`
- `build62-fourth-pair-summary.tsv`
- `build62-fourth-pair-escape.md`

## Post-hoc classifications

- `fourth-pair-recovery`: at least one frozen fourth-pair state authenticates by HMAC;
- `fourth-pair-qualified-gain`: a frozen fourth-pair state improves proposal and oracle error relative to its exact Build61 sibling parent and passes unchanged qualification;
- `fourth-pair-geometric-gain`: proposal and oracle improve but qualification does not pass;
- `fourth-pair-proposal-only`: fourth-pair states exist but no oracle improvement is observed;
- `fourth-pair-not-triggered`: no fourth-pair state is retained on the target.

These labels are diagnostic only and do not change production behavior.


## Qualified-host result

The Build62 host run closes B/mild as `fourth-pair-qualified-gain`: 593 target sibling parents include 141 proposal-local parents and retain 492 fourth-pair states; 491/492 qualify and none authenticates. The best qualified target fourth-pair state reaches **24.056 px** from its exact **25.641 px** parent (`-1.585 px`), with proposal `0.274071`, validation `0.221098`, qualification true and HMAC false. B/angle target remains `fourth-pair-not-triggered`; the full B/angle bank contains **3746** fourth-pair states, **0 qualified**, **0 authenticated**.
