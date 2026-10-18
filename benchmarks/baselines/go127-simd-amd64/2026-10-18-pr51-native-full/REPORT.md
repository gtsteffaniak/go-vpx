# PR51 exactbc4 full-native independent audit

Run39126944619/job117673948005 SUCCESS at2026-10-18T13:06:48Z. Hardware IntelXeonPlatinum8573C, Go1.27.2/SIMD/CGO1/PGOoff/GOMAX1. Raw verifies currentbc4a12a2, previousd08c4f72(tree10c) and original691cd051. All132source hashes bind explicitGit blobs, not later worktree state.

Strict independent PASS:131raw=99top+32nestedStats,132SHAs,35commonfixtures,561×18=10098main timings,527excluded one-shots,396guard+66excluded,792optional+132excluded. Newmacro12fixtures and12288full-state/error/panic/backing oracle pass ordinary/race on all3references; newStats state tests pass selected/off ordinary/race. ActualSDK512true: selected32ordinary+race and optional native oracles PASS. All performance claims remain scoped to hardware/workload.

Statsw8h1 improves17.155→9.2095ns,1.8628×; all6pairs1.817–1.894×. w8h7 is1.2164×, while other shape controls retain regressions down to0.9777×. Prepared-zero token macro cases improve1.7225×/1.6957×; dense-record controls0.971–0.986× regress. All561rows including300median losses and37all-six-slower cases remain. All48optionalFP32 medians0.9510–1.0025×;n2560.9748–1.0017× andn64minimum0.9510×. No feature-first dispatch speedup or causal claim is made.

Whole medians vs50: VP9320+2.264%(4/6slower),VP8+.504%(3/6slower),VP9720−1.400%(4/6faster). Versus original691:−44.882%,−32.851%,−56.861% respectively. Every paired delta is retained; no crossCPU/window pooling or source-identity excuse.

All22semantic fields match. VP9720allocation exceptions: oldround6,previousround4,currentround6 each76/60 versususual75/60; all are retained. Remaining allocation samples/profiles are preserved; all kernel/component timings zeroalloc.

Actual native Stats code capture is independently verified for both variants:9live symbols each, explicit executableSHA/buildinfo, GNUversion/commands, Go nm addresses/sizes/alignment, exact source/fixture maps, rawASM spans and completed optional timing protocolSHA/order. Eachvariant has16files, total32nested. Missing/inlined owners are explicit. Binaryblobs are absent by CI protocol; receipts/raw machinecode are preserved and not replaced by reconstructed codegen.

No new local timing, source/Git/PR edits or merge action. Lean package preserves all rawlogs/profiles/codegen, complete independent auditors/proofs, all losses and exact132tracked source files. No executables or full snapshots.


Publication: `raw-evidence.tar.gz` preserves all 284 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
