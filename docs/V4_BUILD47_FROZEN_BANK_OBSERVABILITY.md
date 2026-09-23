# Format-v4 Build47 — frozen candidate bank observability

Build47 is a research-only diagnostic checkpoint. **Build44 remains the latest qualified production milestone.** Build43 production behavior is unchanged: it still keeps the same two side pairs and freezes at most 32 candidates before held-out qualification.

## Evidence entering Build47

Build46 showed that the single held-out-qualified candidates for both B/mild and B/angle fail a diagnostic singleton decode. Post-hoc oracle distance measured about 62.7 px mean error for B/mild and about 5585 px for B/angle. The existing ensemble and Build42 bank quorums therefore remain frozen.

## Question

Build47 asks whether a substantially better geometry is excluded by Build43 proposal breadth: first by cell-depth pruning within the same two selected side pairs, then by side-pair ranking. It does not assume that the numeric 32-candidate cap alone is the limiting factor.

## Diagnostic bank

The Build47 laboratory path preserves the exact Build43 production proposal tier first (top two side pairs, cell ranks 0 and 3). It then appends a diagnostic selected-pair-depth tier using cell ranks 1 and 2 on those same pairs, followed by an all-pair-extension tier using side-pair ranks 3 through 6 with cell ranks 0 and 3. The reporting stages remain labelled 32/64/128 as their theoretical ceilings; if a tier generates fewer valid candidates, `available` remains below that ceiling. This does **not** change `experimentalV4PhoneBuild43MaxFrozen`, which remains 32 in production.

For every frozen candidate Build47 records proposal, held-out validation, full pilot score/margin/origin, qualification result, source side pair/rank, and source/working quadrilateral. Single-candidate HMAC is attempted only for candidates that already pass the unchanged Build43 qualification gates.

## Oracle discipline

Both blind JSON files are completed before SIFT/reference registration is run. Oracle geometry is then used only post-hoc to measure candidate distance. It cannot create, rank, refine or qualify a blind candidate.

## Private study

```bash
make v4-build47-phone-frozen-bank-test
make v4-build47-phone-frozen-bank-diagnostic
```

The private output is written under `v4-phone private/build47-diagnostics/` and must never be packaged.

Interpretation:

- if the oracle-nearest error drops sharply between 32 and 64/128, candidate pruning is implicated;
- if no candidate approaches the oracle even at 128, proposal generation/basin coverage is the next research target;
- B/mild remains primary; B/angle remains informational.

## Diagnostic revision after first physical run

The first Build47 checkpoint accepted a `maxFrozen=128` argument but still reused the Build43 traversal shape: two selected side pairs, two cell ranks per pair, and at most eight candidates per basin. That traversal can emit at most 32 candidates, so the first physical 32/64/128 report necessarily repeated the same bank (`32/32/32` for B/mild and `28/28/28` for B/angle). Those results remain valid measurements of the production-tier bank, but they do **not** test proposal breadth beyond Build43.

The corrected Build47 diagnostic keeps that production tier intact and adds the two explicit diagnostic tiers described above. This correction changes laboratory observability only; Build43 production code and all qualified Build44 behavior remain unchanged.

## Qualification-host result

The corrected three-tier study was then run on the Surface/WSL2 qualification host.

- **B/mild:** production tier 32 candidates / 1 qualified / oracle-nearest mean error 59.734 px; selected-pair-depth 64 / 2 qualified with the same 59.734 px nearest error; all-pair-extension 128 / 12 qualified with a closer 34.049 px candidate from `top+left`, side-pair rank 3, cell rank 0. That oracle-nearest candidate is **not** qualified.
- **B/angle:** 28 / 52 / 108 candidates across the three tiers; oracle-nearest error changes only from 5499.824 px to 5448.391 px and remains a distant false-basin case.

This rejects simple depth pruning on the two production-selected pairs as the B/mild explanation, while also showing that blindly promoting all six pairs would multiply qualified false basins. Build48 therefore studies bounded proposal-only local refinement of independently selected per-pair seeds rather than weaker qualification or larger production banks.
