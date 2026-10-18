# Rejected ARM49 LF tail-store experiment

Private candidate versus frozen `31123a9ecb7fa3b90afdb293beb06dadb2a6f857` on Apple M4 Max, Go1.27.2 archsimd, recordedCGO1/PGOoff/GOMAXPROCS1. Six balanced100ms rounds preserve120 measured samples plus20 excluded1x oracle samples, all zero-allocation and all14processes exit0/PASS.

Parent/candidate median ratios range 0.992576–1.018766×; 26 pairs improve and 34 regress. Small mixed public changes do not justify added unsafe operations; no ≥1.5× gain is demonstrated. This experiment is rejected.

Exactly16 checked tail byte stores at column17 become exact pointer stores after the worker’s complete20×16/pitch≥20 guard. Source and generator are the only two changed Go paths. Independent static normalization recovers both parent files exactly, verifies rowindices0..15 and preserves capture/store order. The saved regeneration log is empty and root reports successful regeneration; packaging does not rerun Go.

The actual store helper falls260→192 instructions, removes16bounds-panic calls plus morestack and has no calls. Its own16-byte frame disappears; one argument-home store to[sp+8] remains. Reduced codegen cost is not a public performance claim. All1935present parent Go files match baseline Git; only two unrelated benchmark-package tests are absent. Candidate source hashes, patch, generator, caller guards, fixtures, tests/race logs, disassembly, binarySHA/size/buildinfo and allstdout/stderr remain preserved. Full snapshots and executables are omitted.

| Fixture | Parent/candidate | Faster/slower pairs |
|---|---:|---:|
| BenchmarkVP8PickerFilterNative/pitch1312/flat | 1.000000× | 2/4 |
| BenchmarkVP8PickerFilterNative/pitch1312/mixed | 1.005937× | 3/3 |
| BenchmarkVP8PickerFilterNative/pitch1312/noop | 0.992576× | 2/4 |
| BenchmarkVP8PickerFilterNative/pitch1312/ramp | 0.993376× | 2/4 |
| BenchmarkVP8PickerFilterNative/pitch1312/smooth | 1.018766× | 4/2 |
| BenchmarkVP8PickerFilterNative/pitch32/flat | 0.997439× | 2/4 |
| BenchmarkVP8PickerFilterNative/pitch32/mixed | 0.999770× | 3/3 |
| BenchmarkVP8PickerFilterNative/pitch32/noop | 0.995213× | 2/4 |
| BenchmarkVP8PickerFilterNative/pitch32/ramp | 0.995209× | 3/3 |
| BenchmarkVP8PickerFilterNative/pitch32/smooth | 1.001390× | 3/3 |

All raw-member and outer checksums verify. No root/README/Git writes or new build/test/timing/profile jobs occurred.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
