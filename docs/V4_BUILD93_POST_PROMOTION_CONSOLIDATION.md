# Build93 post-promotion consolidation

**Status:** **CLOSED / PASS** on 2026-10-16. The qualified Go 1.26.0 source gate and the normal-production-entry-point 9/9 confirmation both completed successfully. Build94 is unblocked for observability-only profiling.

- **Qualified runtime:** Build93 — Exact Protected Residue Reuse
- **Previous qualified runtime:** Build90 — Ordered Parallel Qualification
- **Qualified toolchain:** Go 1.26.0 / linux amd64
- **Production entry point tested:** `v4-extract-phone`
- **Gate purpose:** `production-entrypoint-integration`

## Purpose

The two independent retained Build93 physical runs already qualified the algorithm and its exact semantics. The post-promotion task was narrower: remove integration ambiguity after the normal `v4-extract-phone` command was rewired from Build90 to Build93, without changing the qualified decoder.

This consolidation is not a new optimization and does not reinterpret Build93 qualification. It closes source and command wiring and verifies the full retained matrix through the production command.

## Source consolidation

The physically qualified Build93 implementation remains byte-for-byte in its retained entry point:

```text
ExperimentalV4ExtractPhoneBuild93Candidate
```

The normal production API is a narrow wrapper:

```text
ExperimentalV4ExtractPhone
  -> ExperimentalV4ExtractPhoneBuild93Candidate
```

The explicit Build93 entry point remains pinned for regression/requalification, while a later runtime build may redirect only the normal entry point after separate qualification. The two pre-consolidation bodies were source-equivalent after comments and function names were removed. Consolidation deleted only the duplicate production body; it did not move or rewrite the function that passed the two historical physical gates.

## Unchanged algorithm and security boundary

No qualified behavior changed:

- direct Build41/43 recovery ordering;
- Build42 post-geometry data recovery;
- Build40 residual fallback;
- Build90 frozen geometry bank and ordered-parallel public-pilot qualification;
- Build93 candidate-local protected-residue sampling and Robust -> Balanced -> Capacity derivation order;
- candidate order, ordered decode commit and first-logical-HMAC behavior;
- encoder, Format-v4 wire format, public pilot, strength 48, ECC/Hamming, whitening/HMAC domains, thresholds and quorums.

Secret key, payload, ECC, HMAC and correctness remain unable to guide geometry proposal, ranking, pruning, scheduling or early exit. HMAC remains final authentication only.

## Completed source gate

The consolidated package was tested with the qualified Go 1.26.0 toolchain:

```bash
make clean
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build93-residue-reuse-test
```

The target completed the Build42/47/68/71/73/75/76/82-93 regression set, the targeted Build93 race fixture, CLI-help regression, VERSION/buildinfo consistency check and `go vet` without failure.

## Completed production integration gate

The retained physical harness accepts only these explicit entry points:

```text
v4-extract-phone
v4-extract-phone-build93
```

The Makefile gate deliberately selected normal `v4-extract-phone`:

```bash
PIXSEAL_GO_TOOLCHAIN=go1.26.0 make v4-build93-phone-production-integration-test
```

The backward-compatible `make v4-build93-phone-physical-test` target invokes the same production gate. Diagnostics were published separately under:

```text
v4-phone private/build93-production-integration/
```

The console result was:

```text
Build93 staged result: controls=3/3 A/front=1/1 A/mild=1/1 A/angle=1/1 B/front=1/1 B/mild=1/1 B/angle-reject=1/1
Build93 semantic-equivalence physical gate: PASS entrypoint=v4-extract-phone purpose=production-entrypoint-integration
```

Run metadata records:

```text
created_utc=2026-10-16T08:42:16Z
gate_purpose=production-entrypoint-integration
entrypoint=v4-extract-phone
qualified_baseline=Build93
previous_qualified_baseline=Build90
toolchain=go1.26.0
pixseal_sha256=c1b69423b692861c80e3cc34eeeaadfd5f4f02c73742599e2e1dbff09819396a
canonical=1632x1632
case_count=9
```

The artifact contains the 9-row TSV, Markdown summary, metadata and one stderr log per image. It contains no default key or payload values.

## Exact retained physical result

The production-entry-point confirmation preserved the qualified matrix:

- three controls: REJECT;
- A/front: Build43 authenticated;
- A/mild: authenticated through its historical path;
- A/angle: direct authenticated;
- B/front: direct authenticated;
- B/mild: Build93 deep authenticated;
- B/angle: REJECT.

All nine rows reported `telemetry_equivalent=true`, `qualification=PASS` and `gate_reason=ok`.

Deep invariants remained exact:

```text
B/mild geometry evaluations       79,259
B/mild bank                          937
B/mild qualification evaluations  1,872
B/mild qualified                     935
B/mild logical decode candidates     691
B/mild logical list frames     2,120,047
B/mild physical decode               696
B/mild speculative                     5
B/mild protected sampling passes     696
B/mild profile derivations         2,086
B/mild HMAC / payload            PASS / PASS

B/angle geometry evaluations      334,857
B/angle bank                         6,198
B/angle qualification evaluations   6,198
B/angle qualified                        0
B/angle protected sampling passes        0
B/angle result                      REJECT
```

The integration-run diagnostic timing was:

| case | deep total | geometry | qualification | decode |
|---|---:|---:|---:|---:|
| B/mild | 102,586 ms | 74,394 ms | 15,658 ms | 12,249 ms |
| B/angle | 223,318 ms | 221,201 ms | 1,702 ms | 0 ms |

Timing is retained for completeness only. This run is not a third promotion-performance sample, and host drift in unchanged phases must not be interpreted as a Build93 algorithmic change.

## Closure decision

**Build93 post-promotion integration is CLOSED / PASS.** The normal production command delegates to the exact implementation used by the two already-qualified runs, and the additional full matrix excluded a rewiring/integration mistake. No further Build93 physical rerun is required for this closure.

The source-controlled non-sensitive summary is `docs/qualified-baselines/build93-production-integration.tsv`. Private images and raw logs remain outside public source archives.

Build94 must not begin with an optimization. It is now unblocked only as observability-only work: re-profile consolidated Build93 on B/mild and B/angle, then select any later target from fresh evidence.
