# Rejected ARM49 experiment: lf-mask-reduce

Private candidate versus exact frozen `31123a9ecb7fa3b90afdb293beb06dadb2a6f857` Go sources on Apple M4 Max. Six balanced parent/candidate rounds at100ms, GOMAXPROCS=1, recordedCGO=1/PGOoff, Go1.27.2 archsimd identified in both saved executable buildinfo records. All14 processes exit zero with PASS. All 120 measured samples and 20 excluded one-shot oracle samples allocate zero. One-shot ns/op values are preserved but never used as performance evidence.

Both parent snapshots omit only two unrelated benchmarks-package test files; all1,935 present Go files independently match baseline Git blobs. Candidate differs in exactly one production source file, recorded with full before/after bytes and patch. Full snapshots and executables are omitted from the lean archive; their source digests and binary hashes/sizes/build information remain reviewable.

| Group | Parent/candidate median range | Faster pairs | Slower pairs | All-six slower cases |
|---|---:|---:|---:|---:|
| all | 0.983671–1.012212× | 32 | 28 | 1 |
| flat | 1.001167–1.002465× | 7 | 5 | 0 |
| ramp | 1.010462–1.012212× | 9 | 3 | 0 |
| smooth | 1.000288–1.001614× | 6 | 6 | 0 |
| mixed | 0.999312–1.003499× | 8 | 4 | 0 |
| noop | 0.983671–0.989374× | 2 | 10 | 1 |

Six full-mask tests change to ReduceMax: **two horizontal and four vertical sites**. The first preparation assertion expected five and failed before a candidate write; the corrected six-site inventory is preserved. A first codegen lookup omitted the Darwin underscore and failed; its failure note and all later successful stdout/stderr remain preserved. These notes are not invented compiler exit receipts. Horizontal bodies remain364instructions and vertical698; instruction counts alone do not establish speed. The public picker remains near-flat with noop regressions and no ≥1.5× gain, so this experiment is rejected.

| Fixture | Parent/candidate | Min paired ratio | Faster/slower pairs |
|---|---:|---:|---:|
| BenchmarkVP8PickerFilterNative/pitch1312/flat | 1.002465× | 0.986688× | 4/2 |
| BenchmarkVP8PickerFilterNative/pitch1312/mixed | 0.999312× | 0.976236× | 5/1 |
| BenchmarkVP8PickerFilterNative/pitch1312/noop | 0.989374× | 0.970297× | 2/4 |
| BenchmarkVP8PickerFilterNative/pitch1312/ramp | 1.012212× | 0.984226× | 5/1 |
| BenchmarkVP8PickerFilterNative/pitch1312/smooth | 1.001614× | 0.988416× | 3/3 |
| BenchmarkVP8PickerFilterNative/pitch32/flat | 1.001167× | 0.860868× | 3/3 |
| BenchmarkVP8PickerFilterNative/pitch32/mixed | 1.003499× | 0.979471× | 3/3 |
| BenchmarkVP8PickerFilterNative/pitch32/noop | 0.983671× | 0.978801× | 0/6 |
| BenchmarkVP8PickerFilterNative/pitch32/ramp | 1.010462× | 0.985654× | 4/2 |
| BenchmarkVP8PickerFilterNative/pitch32/smooth | 1.000288× | 0.990358× | 3/3 |

Every paired value, allocation, one-shot, command, stdout/stderr, available preparation/codegen failure note, correctness/race log and source/codegen artifact is preserved. Empty build logs are retained honestly; compiler command/exit details are not invented. No binaries, full snapshots, runtime reruns, root source/Git/README changes or old-assembly/codec performance claims are introduced. All decompressed member and outer SHA256 hashes verify.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
