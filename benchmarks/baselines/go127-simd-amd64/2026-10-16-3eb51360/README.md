# Native amd64 assembly/parent/current comparison, 2026-10-16

- Assembly-era baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
- Previous SIMD: `5a0cc1222ef4c00d4ed5e2580d61770f117f64cc`.
- Current SIMD: `3eb51360906b9762e26c38e96092a5147d26e4c4`.
- [Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38962395792)
  and [complete hosted CI](https://github.com/thesyncim/govpx/actions/runs/38962395882)
  both succeeded.

All three variants share one EPYC 9V74 VM (AVX2 and AVX512 available),
official Go 1.27.2, `GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`,
and PGO disabled. Six rounds cover every variant permutation once,
placing each variant first, second, and third twice. `orders.txt`
retains the actual schedule. Both SIMD variants assert native dispatch.

Only shared batch-DCT, SSE-strip, and DCT32 benchmark/test files are
overlaid on both reference worktrees. Production sources remain unchanged
at their pinned commits. All builds use the root project's Go/cache
wrapper. This run measures parent/current differences on the same host;
its absolute times cannot isolate changes from earlier runs on other
CPU models or other hosts with the same model.

| Public workload | Assembly | Parent SIMD | Current SIMD | Current vs parent | Current vs assembly |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP9 320x180 | 5.661 ms | 3.409 ms | 3.409 ms | +0.01% | -39.77% |
| VP8 720p | 14.789 ms | 20.539 ms | 20.645 ms | +0.52% | +39.59% |
| VP9 720p | 48.317 ms | 27.491 ms | 27.483 ms | -0.03% | -43.12% |

Values are medians of six samples in ms per input frame. None of the
public parent/current changes establishes an improvement: paired ranges
cross zero for all three workloads (-1.69% to +1.29%, -0.78% to +1.97%,
and -1.11% to +0.16%, respectively). VP8 remains slower than assembly.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only
mode, and no external encoder comparison. VP9 320x180 uses 120 inputs
at 500 kbps; the 720p workloads use 60 inputs at 2,500 kbps.

| Kernel | Assembly | Parent SIMD | Current SIMD | Current vs parent |
| --- | ---: | ---: | ---: | ---: |
| VP9 DCT32 precise | 7,616 ns | 7,616.5 ns | 2,717 ns | -64.33% |
| VP9 DCT32 RD | 7,998 ns | 7,995 ns | 2,726.5 ns | -65.90% |
| SSE16x720 strip | 858.25 ns | 1,430.5 ns | 1,154 ns | -19.33% |
| SSE16x16 | 12.605 ns | 33.61 ns | 27.585 ns | -17.93% |
| SSE8x8 | 9.832 ns | 15.435 ns | 15.44 ns | +0.03% |
| Horizontal luma filter edges | 33.55 ns | 40.465 ns | 40.245 ns | -0.54% |
| Vertical luma filter edges | 103.9 ns | 107.3 ns | 107.95 ns | +0.61% |
| Residual gather 16x16 | 18.61 ns | 90.29 ns | 90.425 ns | +0.15% |
| Residual gather 8x8 | 7.317 ns | 50.33 ns | 50.315 ns | -0.03% |

Each kernel has six 150-ms samples per variant and reports zero bytes
and allocations. DCT32 and SSE16 improve substantially over their parent,
while these changes do not produce an established public realtime gain.
Complete results and controls are retained in raw logs and `summary.json`.

Every output byte count, encoded/dropped count, quantizer statistic and
histogram, keyframe/interframe size, and allocation count matches across
all **18 samples** per public workload. There are no allocation outliers.
Exact bitstreams are asserted separately by codec oracle gates; benchmark
JSON reports retain byte counts.

| Workload | Output bytes | Encoded / dropped | Allocations per input |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Separate 120-frame CPU profiles run for all three variants and both codecs
after all timing rounds. Profile timings are excluded from benchmark
medians. Raw `.pgo` files and flat/cumulative top tables are preserved.

Check out the current checkpoint and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR HEAD^`
on native Linux/amd64 to reproduce. Omitting the last argument retains
the two-variant comparison. Every raw workflow file is copied unchanged;
`summary.json` includes every sample, both comparisons, and parity fields.
