# Native amd64 squared-error/Walsh checkpoint, 2026-10-16

Assembly baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
SIMD checkpoint: `bd25795a1f9774255de64b3c742e113f8ccbf6e8`.
[Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38958062568)
and [complete hosted CI](https://github.com/thesyncim/govpx/actions/runs/38958062570)
both succeeded.

Both implementations share one EPYC 9V74 hosted VM (AVX2 and AVX512
available), official Go 1.27.2, `GOEXPERIMENT=simd`, `GOAMD64=v1`,
`GOMAXPROCS=1`, and PGO disabled. Six serial pairs alternate A/B and B/A;
kernel samples run for 150 ms. Native SIMD selection is asserted in
`dispatch.txt`. Only the batch-DCT and SSE-strip test files are overlaid
on baseline; its production source stays at the pinned commit.

The preceding `3c94d2cc` runner used EPYC 7763. Its absolute timings cannot
isolate these SSE/Walsh changes on this different CPU. Every comparison
below uses its baseline on the same host.

| Workload | Assembly median | SIMD median | Time change | Paired time changes |
| --- | ---: | ---: | ---: | ---: |
| VP9 320x180 | 5.505 ms/frame | 3.380 ms/frame | -38.6% | -38.9% to -38.0% |
| VP8 720p | 13.981 ms/frame | 20.683 ms/frame | +47.9% | +47.7% to +48.8% |
| VP9 720p | 47.022 ms/frame | 25.998 ms/frame | -44.7% | -44.9% to -44.6% |

VP8 remains slower in every pair. VP9 320x180 uses 120 inputs at 500 kbps;
the 720p workloads use 60 inputs at 2,500 kbps. All use realtime CPU-used 8,
30 fps, one thread, encode-only mode, and no external encoder comparison.

| Kernel | Assembly median | SIMD median | Time change |
| --- | ---: | ---: | ---: |
| SSE16x720 strip | 862.35 ns | 1,435.00 ns | +66.4% |
| SSE16x16 | 12.60 ns | 33.585 ns | +166.5% |
| SSE8x8 | 9.834 ns | 15.420 ns | +56.8% |
| Residual DCT, 25 blocks | 163.70 ns | 162.95 ns | -0.5% |

Every kernel sample reports zero bytes/allocations. Every public sample
matches output byte counts, encoded/dropped counts, quantizer
statistics/histograms, keyframe/interframe sizes, and allocations. There
are no allocation outliers. Exact bitstreams are tested separately by
codec oracles; these JSON reports retain byte counts.

| Workload | Output bytes | Encoded / dropped | Allocations per input |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Separate 120-frame profiles run after every timing pair is complete.
The SIMD VP8 profile attributes 30.9% cumulative CPU time to the byte-wide
loop filter, 4.1% to SSE16, and 4.9% to gathering. Profile timings are
excluded from benchmark medians.

Check out the recorded SIMD checkpoint and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR`
on native Linux/amd64. Raw artifact files are preserved unchanged;
`summary.json` retains every sample, paired change, and parity check.
