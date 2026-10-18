# VP8 zero-residual FDCT packets

Private component prototype; not yet integrated. Both ARM64 and x64 two-block kernels skip the full transform when all32 residuals are zero, preserving libvpx1.16 integer-rounding coefficient1=1 in each block. ARM combines zero/out-of-range classification in one unsigned range gate. SIMD packet load/store and overlap order stay unchanged; tails/fallback retain existing paths.

Apple M4 Max, Go1.27.2 SIMD, CGO1, PGOoff, GOMAXPROCS1. Initial two balanced60ms rounds: zero16-block batch2.436x, sparse.977x,dense.994x. Refined ARM range gate in a second two-round screen: zero2.994x median, paired2.586x/3.420x. Sparse parent second sample157.7ns vs67.17ns is a retained outlier; its inflated1.701x median is not a qualifying sparse gain. Dense1.134x is preliminary. No end-to-end or x64 performance claim. Native x64 execution of the new gate remains to CI.

Affected scalar/differential/sentinel tests pass, including literalzero transform and210 overlapping-window cases at counts1,2,3,8,16,25. Alias-order fixture is SIMD-specific and skips when that backend is disabled. Full original raw screens, codegen, hashes and targeted logs remain alongside handoff/. Tests were strengthened/formatted after timing; production mathematical expressions are unchanged. Keep the verified zero-path component win under the user policy even if whole throughput is flat.


Publication: `raw-evidence.tar.gz` preserves all 27 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
