# Format-v4 Build76 — ordered generation-two parallel recovery

Build76 is the **current qualified smartphone baseline**. It is an equivalence-preserving performance successor to Build75 and changes only the physical scheduling of the public proposal-only geometry search.

## Why generation two

Build75 made Build47 freeze substantially cheaper. In its retained second run, B/mild geometry was 76,556 ms with prefix2 21,876 ms; B/angle geometry was 209,184 ms with prefix2 38,015 ms. Prefix2 accumulated only 39,198 ms and 69,071 ms worker time respectively, so effective concurrency is only about 1.79x/1.82x. Gen3 and gen4 already run near the eight-worker limit.

Retained Build72 stage profiling identifies the internal source. For the deterministic dominant B/angle seed, second-generation `single2 + pair2 + continuation2 + sibling2` takes about 33,968 ms of its ~37,846 ms prefix2 time. B/mild shows the same shape at about 17,909 ms of ~20,989 ms.

## Scheduling rule

```text
per-seed exact Build64 prefix through sibling1
        ↓
freeze all sibling1 outputs in seed/traversal order
        ↓
one bounded global pool
 single2 → pair2 → continuation2 → sibling2
        ↓
commit strictly by frozen task index
        ↓
qualified Build73 generation-three pool
        ↓
qualified Build71 generation-four pool
        ↓
unchanged Build68 qualification + Build66 ordered decode
```

Build75 ordered-parallel Build47 freeze is reused unchanged. No pruning, score, threshold, seed depth, bank order, qualification rule, ECC/Hamming, whitening or HMAC domain changes. Protected data never guides geometry scheduling.

## Regression contract

A direct synthetic regression requires Build76 and qualified Build75 to produce exactly the same geometry evaluation count, seed count, worker-compatible view, blind bank size/order, quad, homography, proposal and validation values. The compatibility telemetry also reconstructs the historical Build73 `prefix2` view as `prefix1 + gen2`.

## Physical gate

```bash
make clean
make v4-build76-phone-gen2-parallel-test
make v4-build76-phone-physical-test
```

Diagnostics:

```text
v4-phone private/build76-diagnostics/build76-phone-gen2-parallel.tsv
v4-phone private/build76-diagnostics/build76-phone-gen2-parallel.md
```

The physical report compares geometry and end-to-end wall-clock against retained qualified Build75. Correctness is authoritative.

## Qualification result

Two independent retained Go 1.26.0 runs passed the complete nine-photo matrix with exact logical telemetry.

| metric | Build76 run 1 | Build76 run 2 | Build76 mean | Build75 qualified mean |
|---|---:|---:|---:|---:|
| full matrix elapsed | 764,104 ms | 752,100 ms | **758,102 ms** | 836,722.5 ms |
| B/mild elapsed | 189,947 ms | 187,554 ms | **188,750.5 ms** | 202,688.5 ms |
| B/mild geometry | 60,675 ms | 60,028 ms | **60,351.5 ms** | 76,149 ms |
| B/angle elapsed | 228,044 ms | 225,582 ms | **226,813 ms** | 251,578 ms |
| B/angle geometry | 184,577 ms | 183,036 ms | **183,806.5 ms** | 209,554.5 ms |

The full-matrix mean improves by about **9.4%** versus Build75. B/mild end-to-end improves about **6.9%**, B/angle about **9.8%**; geometry improves about **20.7%** and **12.3%** respectively.

Semantic counters remain frozen: B/mild is `79259 / 937 / 935 / 691 / 2120047`, physical decode `696`, speculative `5`, HMAC/payload PASS; B/angle is `334857 / 6198 / 0`, REJECT. Controls remain rejects and all historical positives remain on their qualified paths.

## Diagnostic atomicity hardening

The second retained run completed successfully and produced the full Markdown matrix. A later interrupted invocation then started at the same output path, truncated the TSV header and cleared the first log before stopping. This exposed an artifact-management weakness rather than a recovery defect. The finalized Build76 script writes TSV, Markdown and logs into a temporary sibling directory and publishes the set only after the complete 9/9 gate passes. Interrupted or failed reruns therefore leave the last complete diagnostics untouched.

**Decision:** Build76 is promoted as the current qualified smartphone baseline. Build75 becomes the historical qualified ordered-basin milestone.
