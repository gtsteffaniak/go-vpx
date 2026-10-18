# TM unroll scalar-contract ordering supplement

This supplements the immutable measured package at /private/tmp/govpx-arm-tm-unroll-clean-ready. The original source was valid for every currently supported selection: TMNative=true occurs only on SIMD ARM where Enabled is const true. The change is clearer ordering and future-proofing, not a fix for a reachable current bug.

The generic TMPredict hook now follows its original !Enabled→TMPredictScalar branch. TMPredictNative also explicitly delegates to that scalar implementation when disabled. On present ARM both branches constant-fold away. The actual native body remains 364 instructions, public TM4 remains 44 instructions, and the timed fixture remains 568 instructions with its single native call for TM8/16. All normalized native/public instructions match measured code. The only fixture differences are two post-loop checksum-sink relocation immediates, each verified against its named symbol target. AMD public bodies are normalized-identical. No retiming or new performance claim is needed.

Dedicated disabled direct-pixel tests cover widths 1/3/8/16, independent short above/left lengths, full backing aliases and exact panic/partial-write state. They compare both TMPredict and TMPredictNative to frozen original scalar loops that reread left per pixel. Each selection checks 15,120 configurations across two entries. They pass for purego, nosimd, AMD forced-disabled, actual AMD AVX/AVX2-off, and a clearly test-only ARM Enabled=false Go source overlay under race/checkptr=2. Production Enabled remains const true on ARM. Native full packages, all original SIMD/scalar state tests and race/checkptr suites also pass.

Correction to the frozen REPORT's retention attribution: the user explicitly mandated keeping >=1.5× component wins and separately sought assembly-beating performance. Retaining these smaller repeatable wins against an already optimized baseline is our engineering judgment consistent with that goal; the user did not explicitly state a below 1.5× retention rule. The measured package and its report remain unmodified so this factual correction stays reviewable.

The overlay contains exactly five production paths changed relative to the frozen left-packet baseline plus five added tests. correction.patch changes only two production files and adds four disabled-scalar test files relative to the measured clean winner. Exact sources, all logs, source-overlay input, source/codegen proofs and binary omission hashes are included. No root files, README, dependencies, assembly or toolchain were changed.

Production changes versus the left-packet winner:

- internal/pixel/kernels_simd.go
- internal/pixel/tmpredict_bytes_arm64.go
- internal/pixel/tmpredict_native_arm64.go
- internal/pixel/tmpredict_native_other.go
- internal/vp8/dsp/intra_dispatch_other.go


Publication: `raw-evidence.tar.gz` preserves all 18 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
