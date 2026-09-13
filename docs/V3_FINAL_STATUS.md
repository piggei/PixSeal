# Format v3 final research status

Status: **implemented, interoperable and frozen compatibility baseline** as of v0.3.0-build23.

This document closes the v0.3 Format-v3 absolute-cycle research line. It does not deprecate v3.
Existing v3 carriers remain the production/interoperability baseline, and all v3 release regressions
continue to be mandatory while Format v4 is developed separately.

## What is frozen

The following remain normative for v3 and must not change accidentally:

- tile geometry: 35×32 blocks = 1120 carrier positions;
- block size and DCT embedding rule;
- `robust`, `balanced` and `capacity` profile capacities/redundancy;
- deterministic tile/code mapping and encoder fingerprints;
- 8-byte header, CRC32, truncated HMAC-SHA256 and HMAC-only authenticity criterion;
- key-derived whitening;
- Hamming(7,4) protected-bit mapping;
- production `EmbedWithInfo` / `ExtractWithInfo` behavior and bounded search banks.

The source baseline is recorded in `V3_FROZEN_CORE_SHA256.txt` and verified by `make v3-freeze-check`.
Updating that manifest is an explicit v3 maintenance event and requires requalification.

## Qualified digital baseline

Build22 host qualification on 2026-09-12 reported:

```text
version-check          PASS
vet                    PASS
release-unit           PASS
test-images            PASS
deep-test              PASS (72/72)
extreme-test           PASS / progressive map
research-unit          PASS
... all dedicated v0.3 diagnostics PASS ...
geometry-test          FAIL(2)
affine-test            FAIL(2)
composition-test       PASS
lattice-test           PASS
perspective-test       PASS
core-target-check      PASS

Release baseline: PASS
Qualification corpus: PASS
Research suites: ATTENTION (2 historical experimental targets)
```

The two red targets are long-standing experimental geometry boundaries and are not regressions of
the authenticated digital release baseline.

## Physical-channel conclusion

No supplied v3 physical acquisition has produced a valid HMAC. Across builds 16–22 the project
therefore tested, in increasingly independent ways, whether the remaining failure was an integer-cycle
selection problem that could be solved safely without changing the format.

The key results are:

1. exact bounded top-2 enumeration removed beam-search uncertainty but remained ambiguous;
2. disjoint repetition cross-fit did not produce a reproducible integer field;
3. 8 deterministic partitions / 16 held-out trials produced a different complete field every time;
4. an independent cross-cell pairwise anchor measured continuous shape but agreed with neither exact
   integer field (0/9 cells on all ambiguous cases);
5. the static v3 audit found only weak/non-uniform public absolute-origin structure and a capacity
   vertical Hamming alias;
6. build22 measured the residual topology asymmetry directly in physical margin grids, but scanner 002
   negative control was at least as stable as useful acquisitions.

Representative build22 physical topology metrics:

| acquisition | mean winner modal phase | mean directional cell consistency |
|---|---:|---:|
| inclined smartphone | 0.3125 | 0.6944 |
| scanner 001 | 0.4375 | 0.7292 |
| scanner 002 negative control | 0.4375 | 0.6597 |

This evidence closes the v3 absolute-cycle branch: another threshold, vote or repartition would reuse
the same insufficient signal and risk converting structured noise into confidence.

## What “closed” means

Closed means:

- no more v3 cycle-threshold tuning on the current evidence family;
- no additional v3 HMAC slots or candidate expansion merely to chase the private corpus;
- no known-header/HMAC oracle may steer geometry;
- all historical diagnostic telemetry remains available for regression and scientific traceability;
- future v3 changes are maintenance, portability, performance-with-identical-behavior or genuine bug fixes.

Closed does **not** mean:

- v3 is removed or deprecated;
- v3 digital robustness claims are withdrawn;
- future independent physical evidence is forbidden;
- Format v4 is already stable.

## Success criterion remains unchanged

For v3 physical tests, only a valid HMAC-authenticated v3 payload is success. Header distance, ECC
statistics, pilot-like diagnostics, topology consistency or any other research metric remain diagnostics.
