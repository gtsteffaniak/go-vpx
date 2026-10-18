# Rejected ARM49 experiment: lf-inline

Private candidate versus exact frozen `31123a9ecb7fa3b90afdb293beb06dadb2a6f857` Go sources on Apple M4 Max. Six balanced parent/candidate rounds at100ms, GOMAXPROCS=1, recordedCGO=1/PGOoff, Go1.27.2 archsimd identified in both saved executable buildinfo records. All14 processes exit zero with PASS. All 120 measured samples and 20 excluded one-shot oracle samples allocate zero. One-shot ns/op values are preserved but never used as performance evidence.

Both parent snapshots omit only two unrelated benchmarks-package test files; all1,935 present Go files independently match baseline Git blobs. Candidate differs in exactly one production source file, recorded with full before/after bytes and patch. Full snapshots and executables are omitted from the lean archive; their source digests and binary hashes/sizes/build information remain reviewable.

| Group | Parent/candidate median range | Faster pairs | Slower pairs | All-six slower cases |
|---|---:|---:|---:|---:|
| all | 0.991265–1.031851× | 41 | 19 | 0 |
| flat | 1.006623–1.012418× | 8 | 4 | 0 |
| ramp | 1.004751–1.011144× | 11 | 1 | 0 |
| smooth | 0.997744–1.003447× | 6 | 6 | 0 |
| mixed | 0.991265–1.006428× | 4 | 8 | 0 |
| noop | 1.021638–1.031851× | 12 | 0 | 0 |

Trivial helper calls are expanded directly into typed SIMD APIs:48abs,12LE/swappedGE,6GT,4shift1,20shift3,6half and10not sites. Horizontal instructions fall364→344, but NOPs fall47→27; vertical instructions fall698→664 while NOPs fall73→39. Non-NOP counts remain317/625 and call targets remain unchanged. This is not evidence of fewer arithmetic operations. Public results are small/mixed, including textured regressions, with no qualifying1.5× gain. Reject the source expansion without a robust meaningful improvement.

| Fixture | Parent/candidate | Min paired ratio | Faster/slower pairs |
|---|---:|---:|---:|
| BenchmarkVP8PickerFilterNative/pitch1312/flat | 1.006623× | 0.993690× | 4/2 |
| BenchmarkVP8PickerFilterNative/pitch1312/mixed | 1.006428× | 0.977768× | 3/3 |
| BenchmarkVP8PickerFilterNative/pitch1312/noop | 1.021638× | 1.000806× | 6/0 |
| BenchmarkVP8PickerFilterNative/pitch1312/ramp | 1.011144× | 1.006686× | 6/0 |
| BenchmarkVP8PickerFilterNative/pitch1312/smooth | 0.997744× | 0.976953× | 1/5 |
| BenchmarkVP8PickerFilterNative/pitch32/flat | 1.012418× | 0.858567× | 4/2 |
| BenchmarkVP8PickerFilterNative/pitch32/mixed | 0.991265× | 0.986360× | 1/5 |
| BenchmarkVP8PickerFilterNative/pitch32/noop | 1.031851× | 1.012048× | 6/0 |
| BenchmarkVP8PickerFilterNative/pitch32/ramp | 1.004751× | 0.992728× | 5/1 |
| BenchmarkVP8PickerFilterNative/pitch32/smooth | 1.003447× | 0.994452× | 5/1 |

Every paired value, allocation, one-shot, command, stdout/stderr, available preparation/codegen failure note, correctness/race log and source/codegen artifact is preserved. Empty build logs are retained honestly; compiler command/exit details are not invented. No binaries, full snapshots, runtime reruns, root source/Git/README changes or old-assembly/codec performance claims are introduced. All decompressed member and outer SHA256 hashes verify.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
