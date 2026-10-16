# Native amd64 assembly/parent/current comparison, 2026-10-16

- Assembly-era baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
- Previous SIMD: `bd25795a1f9774255de64b3c742e113f8ccbf6e8`.
- Current SIMD: `5a0cc1222ef4c00d4ed5e2580d61770f117f64cc`.
- [Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38960429221)
  and [complete hosted CI](https://github.com/thesyncim/govpx/actions/runs/38960429206)
  both succeeded.

All three variants share one EPYC 7763 VM (AVX2, no AVX512), official
Go 1.27.2, `GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`, and PGO
disabled. Six rounds cover every variant permutation once, placing each
variant first, second, and third twice. `orders.txt` retains the actual
schedule. Both SIMD variants assert native backend selection.

Only the shared batch-DCT and SSE-strip test files are overlaid on both
reference worktrees. Production sources stay unchanged at their pinned
commits. All builds use the root project's Go/cache wrapper. This run
measures the parent/current difference on the same host, avoiding the
changing CPU models in earlier separate runs.

| Public workload | Assembly | Parent SIMD | Current SIMD | Current vs parent | Current vs assembly |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP9 320x180 | 6.977 ms | 4.474 ms | 4.459 ms | -0.33% | -36.09% |
| VP8 720p | 23.191 ms | 27.003 ms | 25.051 ms | -7.23% | +8.02% |
| VP9 720p | 60.159 ms | 34.793 ms | 35.053 ms | +0.75% | -41.73% |

Values are medians of six samples in ms per input frame. VP8 improves
against its parent in every round, with paired changes from -9.12% to
-6.50%, while remaining slower than assembly. VP9 parent/current changes
range from -1.08% to +0.29% at 320x180 and -0.62% to +2.38% at 720p;
these samples do not establish a VP9 improvement from this checkpoint.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only
mode, and no external encoder comparison. VP9 320x180 uses 120 inputs
at 500 kbps; the 720p workloads use 60 inputs at 2,500 kbps.

| Kernel | Parent SIMD | Current SIMD | Current vs parent |
| --- | ---: | ---: | ---: |
| Horizontal luma loop-filter edges | 91.04 ns | 50.89 ns | -44.10% |
| Vertical luma loop-filter edges | 170.35 ns | 130.90 ns | -23.16% |
| Residual gather 16x16 | 184.95 ns | 105.20 ns | -43.12% |
| Residual gather 8x8 | 62.66 ns | 65.18 ns | +4.02% |
| SSE16x720 strip | 1,914 ns | 1,896 ns | -0.94% |
| SSE16x16 | 40.255 ns | 39.370 ns | -2.20% |
| SSE8x8 | 17.78 ns | 18.40 ns | +3.49% |
| Residual DCT, 25 blocks | 216.0 ns | 216.5 ns | +0.23% |

Each kernel has six 150-ms samples per variant and reports zero bytes
and allocations. Complete assembly/parent/current kernel measurements,
including unchanged controls and other small regressions, are retained
in the raw logs and `summary.json`.

Every output byte count, encoded/dropped count, quantizer
statistic/histogram, keyframe/interframe size, and allocation count
matches across all **18 samples** per public workload. There are no
allocation outliers. Exact bitstreams are tested separately by codec
oracles; these benchmark JSON reports retain byte counts.

| Workload | Output bytes | Encoded / dropped | Allocations per input |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Separate 120-frame CPU profiles run for all three variants and both codecs
after every timing round is complete. Profile timings are excluded from
benchmark medians. Raw `.pgo` files and flat/cumulative top tables are
retained for each variant.

Check out the recorded current checkpoint and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR HEAD^`
on native Linux/amd64 to reproduce the harness. Omitting the last argument
retains the original two-variant alternating comparison. All raw workflow
files are preserved unchanged. `summary.json` includes every sample,
assembly/current and parent/current deltas, and both parity comparisons.
