# Rejected ARM49 experiment: tm-interval

Private candidate versus exact frozen `31123a9ecb7fa3b90afdb293beb06dadb2a6f857` Go sources on Apple M4 Max. Six balanced parent/candidate rounds at100ms, GOMAXPROCS=1, recordedCGO=1/PGOoff, Go1.27.2 archsimd identified in both saved executable buildinfo records. All14 processes exit zero with PASS. All 324 measured samples and 54 excluded one-shot oracle samples allocate zero. One-shot ns/op values are preserved but never used as performance evidence.

Both parent snapshots omit only two unrelated benchmarks-package test files; all1,935 present Go files independently match baseline Git blobs. Candidate differs in exactly one production source file, recorded with full before/after bytes and patch. Full snapshots and executables are omitted from the lean archive; their source digests and binary hashes/sizes/build information remain reviewable.

| Group | Parent/candidate median range | Faster pairs | Slower pairs | All-six slower cases |
|---|---:|---:|---:|---:|
| all | 0.964931–1.038407× | 63 | 98 | 6 |
| 4x4 | 0.978026–1.021025× | 24 | 29 | 0 |
| 8x8 | 0.990070–1.038407× | 35 | 19 | 0 |
| 16x16 | 0.964931–0.977931× | 4 | 50 | 6 |

Two8/16-byte overlap guards switch to unsigned interval predicates without arithmetic/store changes. LLVM body count falls364→352, while the16x16 public cases all regress on medians and8x8 is mixed. The4x4 cases are unchanged controls. Fewer instructions do not establish a speedup. No ≥1.5× component gain or robust public improvement is demonstrated, so this experiment is rejected.

| Fixture | Parent/candidate | Min paired ratio | Faster/slower pairs |
|---|---:|---:|---:|
| BenchmarkTMPredictRowsNative/16x16/pitch1312/clipped | 0.977931× | 0.959771× | 2/4 |
| BenchmarkTMPredictRowsNative/16x16/pitch1312/flat | 0.974244× | 0.959184× | 1/5 |
| BenchmarkTMPredictRowsNative/16x16/pitch1312/mixed | 0.964931× | 0.951154× | 0/6 |
| BenchmarkTMPredictRowsNative/16x16/pitch16/clipped | 0.976684× | 0.953760× | 0/6 |
| BenchmarkTMPredictRowsNative/16x16/pitch16/flat | 0.976722× | 0.945833× | 0/6 |
| BenchmarkTMPredictRowsNative/16x16/pitch16/mixed | 0.976148× | 0.962369× | 1/5 |
| BenchmarkTMPredictRowsNative/16x16/pitch37/clipped | 0.969451× | 0.959774× | 0/6 |
| BenchmarkTMPredictRowsNative/16x16/pitch37/flat | 0.976909× | 0.967080× | 0/6 |
| BenchmarkTMPredictRowsNative/16x16/pitch37/mixed | 0.975658× | 0.963624× | 0/6 |
| BenchmarkTMPredictRowsNative/4x4/pitch1312/clipped | 0.995591× | 0.966582× | 1/4 |
| BenchmarkTMPredictRowsNative/4x4/pitch1312/flat | 1.001663× | 0.992279× | 4/2 |
| BenchmarkTMPredictRowsNative/4x4/pitch1312/mixed | 1.002764× | 0.986726× | 4/2 |
| BenchmarkTMPredictRowsNative/4x4/pitch37/clipped | 0.997069× | 0.984593× | 1/5 |
| BenchmarkTMPredictRowsNative/4x4/pitch37/flat | 0.989124× | 0.984949× | 2/4 |
| BenchmarkTMPredictRowsNative/4x4/pitch37/mixed | 0.999449× | 0.978467× | 4/2 |
| BenchmarkTMPredictRowsNative/4x4/pitch4/clipped | 0.978026× | 0.922528× | 1/5 |
| BenchmarkTMPredictRowsNative/4x4/pitch4/flat | 1.021025× | 0.998894× | 5/1 |
| BenchmarkTMPredictRowsNative/4x4/pitch4/mixed | 0.989868× | 0.953394× | 2/4 |
| BenchmarkTMPredictRowsNative/8x8/pitch1312/clipped | 1.020506× | 0.988509× | 5/1 |
| BenchmarkTMPredictRowsNative/8x8/pitch1312/flat | 1.034144× | 1.016508× | 6/0 |
| BenchmarkTMPredictRowsNative/8x8/pitch1312/mixed | 1.026917× | 0.984397× | 4/2 |
| BenchmarkTMPredictRowsNative/8x8/pitch37/clipped | 1.000279× | 0.974116× | 3/3 |
| BenchmarkTMPredictRowsNative/8x8/pitch37/flat | 0.990070× | 0.979843× | 2/4 |
| BenchmarkTMPredictRowsNative/8x8/pitch37/mixed | 0.995238× | 0.982222× | 3/3 |
| BenchmarkTMPredictRowsNative/8x8/pitch8/clipped | 0.995646× | 0.982947× | 2/4 |
| BenchmarkTMPredictRowsNative/8x8/pitch8/flat | 1.038407× | 1.018172× | 6/0 |
| BenchmarkTMPredictRowsNative/8x8/pitch8/mixed | 1.007575× | 0.972838× | 4/2 |

Every paired value, allocation, one-shot, command, stdout/stderr, available preparation/codegen failure note, correctness/race log and source/codegen artifact is preserved. Empty build logs are retained honestly; compiler command/exit details are not invented. No binaries, full snapshots, runtime reruns, root source/Git/README changes or old-assembly/codec performance claims are introduced. All decompressed member and outer SHA256 hashes verify.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
