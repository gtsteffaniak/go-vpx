# Native forward32 and fixed-width filter loads

Apple M4 Max, macOS ARM64, Go1.27.2, GOEXPERIMENT=simd, GOMAXPROCS=1, PGO off. Builds/tests completed before serialized timings.

Six alternating pairs compare the unchanged scalar32 network at5a0cc122 with four-lane native SIMD32. Precise median4699.5→2084.5ns (-55.64%), RD4709.5→1953.5ns (-58.52%), all zero allocations. Numeric parity includes independent scalar, existing pinned libvpx fixtures, complete signed-input fallback, exact/misaligned buffers and in-place operation. Additional deterministic cosine-sign extrema cover all1024 frequency pairs.

Fixed four/eight-byte column tiles remove eight per-row width tests and omit the unused upper transpose for simple filters. Six pairs give vertical internal edges74.06→72.28ns (-2.40%); a fresh six pairs give75.46→74.275ns (-1.57%). Horizontal changes are small (-1.60% and -0.98%). All allocate zero.

Six balanced public triples compare perf22=5a0cc122, perf23=SIMD32 only, perf24=fixed-width filter loads only. Workloads: VP8/VP9 1280x720,60frames,2500kbps; VP9 320x180,120frames,500kbps; realtime,cpu-used8,threads1,30fps,encode-only,auto-libvpx=false. Changes are small/noisy: SIMD32 VP8 +0.47%, VP9 small -0.54%, VP9 720 +0.08%; fixed loads -0.29%,+0.35%,-0.09%. These do NOT establish a public encoder speedup. All output/count/quantizer/allocation fields match in every sample.

A separate fixed-register-shift32 prototype was rejected: precise+0.93%, RD-0.74% did not establish a gain. It is not retained. The production32 implementation matches perf23.

`dct32-bounds.py` is a diagnostic affine coefficient/rounding-error bound calculation; invoke with `/usr/bin/python3 benchmarks/baselines/go127-simd-arm64/2026-10-16-dct32-filter-loads/dct32-bounds.py` from the repository root. Rounding is bounded by0.5 at each14-bit network normalization. A conservative intermediate bound5771 yields maximum pre-round arithmetic2139425120; plus8192 fits signed int32. Larger source values retain scalar Go.
