# Qualified smartphone performance references

These files freeze source-controlled numeric references without shipping the private physical images, secret key, payload files, or raw diagnostic logs.

- `build76-phone-performance.tsv` — historical Build76 two-run reference.
- `build84-phone-performance.tsv` — earlier qualified Build84 two-run command-matrix reference.
- `build90-phone-performance.tsv` — previous qualified Build90 evidence: all three full semantic physical runs are retained, including the host-degraded run, plus the contemporaneous Build91 paired means used for its promotion.
- `build91-build84-build90-paired.tsv` — exact authoritative rows from the completed Build91 AB/BA paired physical profile.
- `build93-phone-performance.tsv` — current qualified Build93 evidence: two independent retained 9/9 physical PASS runs, protected sampling/derivation counts, and frozen B/mild promotion comparisons.
- `build93-production-integration.tsv` — completed Go 1.26.0 post-promotion wiring confirmation through normal `v4-extract-phone`; nine semantic rows plus non-sensitive entry-point/toolchain/executable-hash metadata.
- `build94-profile-summary.tsv` — non-sensitive low-overhead timing and key pprof attribution from the completed Build94 profile of unchanged qualified Build93; this selects the bounded Build95 benchmark candidate but does not change the qualified baseline.

Build93 promotion is based on exact source/race equivalence, the public residue-reuse benchmark, and two independent retained 9/9 semantic physical gates with repeatable B/mild decode improvement. The later production-entry-point matrix is integration evidence only and is kept separate from promotion-performance evidence. Build90 remains preserved as the previous qualified reference. Cross-run absolute geometry timing is retained rather than cherry-picked.
