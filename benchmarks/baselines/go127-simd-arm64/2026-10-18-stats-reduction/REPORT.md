# Optional ARM49 Stats reduction candidate

Private candidate versus frozen `31123a9ecb7fa3b90afdb293beb06dadb2a6f857`, Apple M4 Max, Go1.27.2 archsimd, recordedCGO1/PGOoff/GOMAXPROCS1. Status is **optional/pending**: this package neither rejects the candidate nor declares AMD/whole-codec acceptance.

One production source changes final lane reductions to modular vector sums. Original signed32lane sums are widened toint64 and finally cast to uint32/int32; reducing within the vector gives the same result modulo2^32, including the odd-row tail additions. Loads, bounds, dispatch and arithmetic terms stay unchanged. The file is shared ARM/AMD code, so ARM results do not establish native AMD performance.

Each separate window preserves216 timed samples plus36 excluded1x oracles, all zero-allocation. The initial six balanced100ms rounds start parent-first; fresh six start candidate-first. All28processes exit0/PASS. The windows are not pooled. The identical18-case public Stats fixture exercises widths8/32/64 and heights1/7/8/64/4096/4097, uses independent StatsScalar checks, a live return sink and complete input-immutability checks.

Actual ARM code changes Rows16 from124→109 instructions and Pairs8 from185→169, with two ADDV instructions replacing scalar lane extraction in each function. Call-symbol sets remain unchanged. Instruction differences are separate from native measurements. Savedpixel/VP8DSP race/checkptr logs pass; the VP9DSP line explicitly says no tests to run. Existing extreme and tall scalar-oracle tests are preserved.

| Public fixture | Initial parent/candidate | Fresh parent/candidate | Initial faster/slower pairs | Fresh faster/slower pairs |
|---|---:|---:|---:|---:|
| BenchmarkStatsReductionNative/w32/h1 | 1.020427× | 1.033925× | 6/0 | 6/0 |
| BenchmarkStatsReductionNative/w32/h4096 | 1.011123× | 1.009640× | 4/2 | 4/2 |
| BenchmarkStatsReductionNative/w32/h4097 | 1.015169× | 1.004869× | 4/2 | 3/3 |
| BenchmarkStatsReductionNative/w32/h64 | 1.006761× | 1.016494× | 5/1 | 4/1 |
| BenchmarkStatsReductionNative/w32/h7 | 1.045272× | 1.024250× | 5/1 | 6/0 |
| BenchmarkStatsReductionNative/w32/h8 | 1.018830× | 1.029993× | 4/2 | 6/0 |
| BenchmarkStatsReductionNative/w64/h1 | 1.036479× | 1.041345× | 6/0 | 6/0 |
| BenchmarkStatsReductionNative/w64/h4096 | 0.998236× | 1.001796× | 3/3 | 4/2 |
| BenchmarkStatsReductionNative/w64/h4097 | 0.994573× | 0.998519× | 2/4 | 4/2 |
| BenchmarkStatsReductionNative/w64/h64 | 0.992211× | 1.004505× | 1/5 | 4/2 |
| BenchmarkStatsReductionNative/w64/h7 | 1.011482× | 1.009195× | 5/0 | 5/1 |
| BenchmarkStatsReductionNative/w64/h8 | 1.009274× | 1.004886× | 6/0 | 5/1 |
| BenchmarkStatsReductionNative/w8/h1 | 1.039455× | 1.021848× | 5/1 | 4/2 |
| BenchmarkStatsReductionNative/w8/h4096 | 1.002320× | 1.001638× | 4/2 | 3/2 |
| BenchmarkStatsReductionNative/w8/h4097 | 1.006129× | 0.995242× | 5/1 | 0/5 |
| BenchmarkStatsReductionNative/w8/h64 | 1.005680× | 1.000576× | 4/2 | 3/3 |
| BenchmarkStatsReductionNative/w8/h7 | 1.047447× | 1.053892× | 6/0 | 6/0 |
| BenchmarkStatsReductionNative/w8/h8 | 1.036347× | 1.047932× | 6/0 | 6/0 |

The short8-row7/8 results improve modestly in both windows; these are below1.5×, not threshold-qualified claims. Long/tall near-flat or slower controls remain explicit in the table and complete paired arrays. No candidate-specific AMD or whole-codec gain has been measured.

All1935present baseline Go files match Git plus the one identical new benchmark fixture; two unrelated benchmark-package tests are absent. Candidate differs only in the reduction file. Exact relevant source/test/fixture bytes, patch, rawstdout/stderr, bothdrivers, commands/timestamps/exits, codegen and binarySHA/size/buildinfo are preserved. Full snapshots and executables are omitted. All raw-member and outer SHA256 checksums verify. No new Go/build/test/profile/timing or ROOT/README/Git changes occurred.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
