# ARM64 byte filters and tall statistics

Apple M4 Max, darwin/arm64, official Go1.27.2, GOEXPERIMENT=simd,
GOMAXPROCS=1, PGO=off. Six alternating encoder pairs compare commit6a7b0345
(`perf6a`) with this checkpoint (`perf17`): taller statistics dispatch with
bounded signed accumulation, direct exact-width transposed stores, and
byte-domain saturated filter correction. No competing local build/test jobs
ran during these encoder timings. See encoder-pairs.json for every workload,
count, quantizer statistic, latency, and allocation result.

VP8 720p median12.443→11.202ms/frame (-10.0%); VP9 changes+0.8% at320x180
and+0.6% at720p are small and not evidence of improvement. All output-byte,
encoded/dropped-frame, quantizer, and allocation fields agree across each
workload. The assembly-era ARM baseline remains faster; this checkpoint
does not establish ARM performance parity.

Separate six-pair kernel comparisons: tall16x720 statistics 1,980→1,090ns,
zero allocations; removing the transposed row array improves macroblock
vertical filtering about29.4→25.7ns. A second comparison specializes the
2/4-byte stores, improving vertical luma edges100.3→90.8ns. Raw kernels
are retained separately; do not combine their absolute values across runs.

The amd64 byte path was checked for correctness through Rosetta, including
all130,816 signed outer/delta combinations, but native x64 performance is
pending the new hosted run. No Rosetta timings support performance claims.
