# VP8 two-block DCT and fixed ARM shifts

Apple M4 Max, darwin/arm64, official Go 1.27.2, GOEXPERIMENT=simd,
GOMAXPROCS=1 / test.cpu=1, PGO=off. No competing build/test processes ran
while these six alternating comparisons were timed.

Kernel variants: ff = unchanged ff9395d0 production code, shifts = only
fixed register-shift helpers, batch = shifts plus the eight-lane two-block
DCT. The same added benchmark runs on each variant with 25 consecutive
blocks spanning valid VP8 residuals [-255,255]. Every sample allocates zero.
Median batch times are approximately 221 ns, 186 ns, and 107 ns respectively.
The older Batch25 benchmark uses out-of-range inputs, so it cannot establish
the native residual path's throughput; it remains as fallback coverage.

The separate public VP8 720p comparison improves 11.684 to 11.165 ms/frame
(-4.44%) over six pairs. perf17 is the ff9395d0 production source with PGO
off; perf18 adds the fixed ARM shifts and batched DCT. Output bytes,
encoded/dropped frames, quantizer distribution and zero allocations match
in every sample. The original assembly-era ARM encoder is still faster.

Independent scalar tests cover zero/odd/even block counts, full signed
inputs, realistic residuals, and exact in-place operation. Values outside
[-255,255] retain scalar arithmetic. The amd64 implementation instead uses
128-bit pairwise dot products and AVX2-safe narrowing; Rosetta correctness
passes, with native x64 performance pending the next hosted run.
