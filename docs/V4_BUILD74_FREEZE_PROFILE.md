# Format-v4 Build74 — Build47 freeze-stage profiling

## Status

Build74 is **observability-only** over the qualified Build73 smartphone baseline. It is not a performance candidate and cannot supersede Build73.

## Why Build74 exists

Two independent Build73 qualification runs leave the Build47 freeze as the most stable residual deep-geometry component: approximately 49–64 seconds across the retained deep cases. `B/angle` still spends more wall time in generation four, but that phase already compresses worker time by almost the full eight-core factor. Another parallel barrier is therefore not justified without understanding the freeze itself.

## Exact freeze decomposition

The historical Build47 freeze performs, in order:

1. structural side-line seeding;
2. a private `pixelPlane` construction;
3. robust scoring/ranking of all side pairs;
4. cached cell ranking for visited pair ranks;
5. basin generation in three ordered tiers:
   - production (`top two pairs`, cell ranks `0,3`);
   - selected-pair depth (`top two pairs`, cell ranks `1,2`);
   - all-pair extension (`remaining pairs`, cell ranks `0,3`).

Build74 reproduces that code path line-for-line in a separate helper and records wall time, helper-evaluation counts and task counts for each region. Basin calls additionally record min/median/max elapsed, evaluation count and output count.

The historical third return value of `experimentalV4PhoneBuild47Freeze` is a **proposal count**, not the internal helper-evaluation total. Build74 preserves that return contract exactly so the qualified Build64/73 `geometry-evals` telemetry remains unchanged. New internal evaluation counts are diagnostic only.

## Equivalence requirement

A direct synthetic regression runs both the historical Build47 freeze and the profiled Build74 freeze and requires exact equality of:

- side-pair ranking and scores;
- proposal count;
- frozen bank length and order;
- tier/cell provenance;
- quad, homography and proposal score for every frozen candidate.

Build74 then reuses the qualified Build73 prefix2/gen3/gen4 scheduler, Build68 qualification plane reuse and Build66 ordered protected-data decode unchanged.

## Physical gate

```bash
make clean
make v4-build74-phone-freeze-profile-test
make v4-build74-phone-physical-test
```

Diagnostics:

```text
v4-phone private/build74-diagnostics/build74-phone-freeze-profile.tsv
v4-phone private/build74-diagnostics/build74-phone-freeze-profile.md
```

Correctness is authoritative. Fine-grained freeze timing is used only to select a later Build75 target. Build73 remains the qualified baseline throughout.

## Retained physical result

The qualified Go 1.26.0 host completed the full nine-photo gate with exact Build73 semantic equivalence. Deep-case freeze totals are 54,945 / 55,203 / 62,086 / 61,776 / 55,501 ms for control/front, control/mild, control/angle, B/mild and B/angle. Basin generation consumes 42,956 / 43,126 / 47,749 / 47,712 / 43,836 ms respectively, consistently ~77–79% of freeze. Pair scoring is the next serial region at 10,656–13,403 ms; cells are ~0.7–0.8 s and plane preparation ~0.2 s. Every deep case executes exactly 16 basin tasks (4 production, 4 depth, 8 all-pairs). This retained evidence selects Build75 ordered-parallel basin generation.
