# Final combined ARM49 whole-codec evidence

Apple M4 Max, official Go1.27.2 SIMD, CGO1, PGOoff, runtimeGOMAXPROCS1, VP8 realtimeCPU8, one thread,1280×720,2500kbps,300inputframes. Six balanced three-way permutations produced18 JSON samples. No new profiles were run in this final window. The exact current candidate keeps owner158, Stats reduction and AMD-only FP32 changes; rejected guard156 was not integrated.

| Source | Median ns/input |
|---|---:|
| Exact PR48 |7,273,517.5|
| Combined candidate |7,206,938.5|
| Old691 |6,856,174.0|

Current median is0.91536179% faster than PR48; five of six pairs are faster, paired changes−0.95994%,−0.96019%,+0.12148%,−0.36481%,−1.13427%,−3.19873%. This is a small measured combined gain with one slower pair, not an attribution to an individual kernel.

Current median remains5.11603848% slower than old691. Four paired gaps are+4.718–5.375%; two pairs favor current by0.76281% and0.01910%. All raw samples, including old slower observations, are retained. This window does not establish that the old assembly whole-codec gap has closed. No cross-window gap subtraction or multiplied stage speedups are used.

Every JSON matches all22 non-timing semantic fields and has0allocations/input. Quality was skipped: equal recorded fields do not assert bitstream identity or measured quality equality. Exact encoded/dropped counts, output bytes and quantizer statistics are retained in the compact audit and raw JSON.

All1,946 current Go source hashes and reference1,937/1,730 hashes were checked against receipts before and after timing. Current has exactly8 production and7 test deltas relative to PR48. All five DC production files match the frozen158-owner candidate. Three CLI executables are omitted, with exact sizes, private paths and SHA256 receipts; all original binaries remain private. Current was rebuilt with buildGOMAXPROCS2, references with1; runtime is uniformly1. SDK, CGO and simd are inspected in embedded metadata, while PGOoff is proven from saved build commands rather than inferred from metadata.

Root full native suite, focused three-package race/checkptr2, no-assembly/build-selection/format checks, PGO refresh and pre-commit all exited0. Complete logs preserve nonfatal sandbox stat-cache warnings. The private DC stage package separately preserves all2,592 component samples, full-state tests, LLVM/caller proofs, proof/parser failures and rejected156 guard; its tar and outer-manifest hashes are linked here. No source, README, dependency, branch or Git mutation was made by this evidence freeze.


Publication: `raw-evidence.tar.gz` preserves all 8 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
