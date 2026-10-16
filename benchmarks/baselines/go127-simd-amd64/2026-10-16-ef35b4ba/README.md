# Native amd64 assembly/parent/current comparison, 2026-10-16

- Assembly-era baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
- Previous SIMD: `3eb51360906b9762e26c38e96092a5147d26e4c4`.
- Current SIMD: `ef35b4bafd2822817b96c82815255f2097bfe661`.
- [Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38965788305)
  and [complete hosted CI](https://github.com/thesyncim/govpx/actions/runs/38965788353)
  both succeeded.

All variants share one EPYC 7763 VM (AVX2, no AVX512), official Go 1.27.2,
`GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`, and PGO disabled.
Six rounds cover every permutation once, placing each variant first,
second, and third twice. `orders.txt` retains the schedule. Both SIMD
variants assert native dispatch.

Only shared batch-DCT, SSE-strip, DCT32, and in-place paired-DC test/benchmark
fixtures are overlaid on reference worktrees; production sources stay
unchanged at their pinned commits. Builds use the root project's Go/cache
wrapper. Parent/current comparisons share the same host. Comparing absolute
times across separate runs does not isolate code changes, even for hosts
with the same CPU model.

| Public workload | Assembly | Parent SIMD | Current SIMD | Current vs parent | Current vs assembly |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP9 320x180 | 7.014 ms | 4.488 ms | 4.494 ms | +0.13% | -35.93% |
| VP8 720p | 23.064 ms | 25.100 ms | 23.595 ms | -5.99% | +2.30% |
| VP9 720p | 59.821 ms | 34.565 ms | 35.036 ms | +1.36% | -41.43% |

These are six-sample medians in ms per input frame. VP8 improves against
its parent in every pair (-6.65% to -4.99%), while remaining 2.30% slower
than assembly. VP9 parent/current paired ranges cross zero (-0.15% to
+1.61% at 320x180; -0.45% to +2.48% at 720p), so these samples do not
establish a VP9 gain from this checkpoint.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only
mode, and no external encoder comparison. VP9 320x180 uses 120 inputs
at 500 kbps; the 720p workloads use 60 inputs at 2,500 kbps.

| Kernel | Assembly | Parent SIMD | Current SIMD | Current vs parent |
| --- | ---: | ---: | ---: | ---: |
| Paired DC reconstruction | 11.85 ns | 39.305 ns | 15.29 ns | -61.10% |
| Paired DC, exact in-place | 13.435 ns | 40.605 ns | 17.505 ns | -56.89% |
| Ordinary UV filter, horizontal | 28.985 ns | 62.38 ns | 28.095 ns | -54.96% |
| Ordinary UV filter, vertical | 36.73 ns | 86.13 ns | 44.945 ns | -47.82% |
| Macroblock UV filter, horizontal | 31.345 ns | 70.505 ns | 31.52 ns | -55.29% |
| Macroblock UV filter, vertical | 39.165 ns | 102.85 ns | 58.935 ns | -42.70% |

Each kernel has six 150-ms samples per variant, with zero bytes and
allocations. Full results, unchanged controls, and samples are retained in
raw logs and `summary.json`; remaining kernel gaps versus assembly are
visible alongside the parent improvements.

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

Separate 120-frame CPU profiles run for every variant and both codecs
after all timing rounds. Profiles are excluded from timing medians.
Raw `.pgo` files and flat/cumulative top tables are retained.

Check out this checkpoint and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR HEAD^`
on native Linux/amd64 to reproduce. Omitting the last argument preserves
the two-variant comparison. Every raw artifact file is copied unchanged;
`summary.json` retains every sample, both comparisons, and parity fields.
