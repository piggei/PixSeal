# Format-v4 Build48 — proposal-only local projective refinement

Build48 is a research-only diagnostic checkpoint. **Build44 remains the latest qualified production milestone.** No Format-v4 encoding, strength, pilot, protected-data mapping, ECC/Hamming, whitening/HMAC, Build42 list logic, Build43 production candidate cap, side-pair count, qualification threshold or quorum changes.

## Evidence entering Build48

The corrected Build47 host run shows that B/mild is not limited by deeper cells on the same two production-selected side pairs: oracle-nearest mean error stays at 59.734 px from the production tier through selected-pair depth. The all-pair diagnostic tier finds a closer 34.049 px geometry from `top+left`, side-pair rank 3, cell rank 0, but it does not pass the unchanged held-out qualification gate. B/angle remains thousands of pixels from the oracle and is informational.

## Hypothesis

A lower-ranked B/mild basin may already be close enough for a bounded local projective refinement to converge toward the correct geometry. The refinement must remain blind and proposal-only. Held-out evidence can only accept/reject the frozen result.

## Blind method

1. Generate the corrected Build47 extended proposal bank.
2. Group candidates by side-pair rank.
3. Keep at most two seeds per side pair using **proposal score only**.
4. Refine all eight quadrilateral coordinates with the established bounded Build41 coordinate descent, still scoring only proposal folds.
5. Freeze the complete refined bank.
6. Evaluate pre/post held-out validation and complete public-pilot score/margin/origin.
7. Only for post-qualified candidates, run the unchanged diagnostic Build42 list/HMAC check.
8. Only after both blind B/mild and B/angle JSON files exist, generate the SIFT/reference oracle and calculate before/after corner error.

The secret key, payload, ECC result, HMAC and oracle quadrilateral are unavailable to seed selection and geometry refinement.

## Commands

```bash
make v4-build48-phone-local-refine-test
make v4-build48-phone-local-refine-diagnostic
```

The private diagnostic writes under `v4-phone private/build48-diagnostics/` and is never packaged.

## Decision gate

B/mild is primary. A material reduction from the Build47 ~34.049 px basin toward a few pixels, especially if the refined geometry independently passes held-out/full-pilot qualification and authenticates, supports a later bounded production experiment. If proposal score rises while oracle error stays large, the proposal objective itself is displaced and further bank expansion is not justified. B/angle remains informational and cannot alone drive production changes.
