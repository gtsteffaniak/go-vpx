# Native amd64 assembly/parent/current comparison, 2026-10-16

- Assembly-era baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
- Previous SIMD: `08c2b96ed68fd41e46d46112781e88d7408f7551`.
- Current SIMD: `3792c61706a100f22e993636a6ba697fa4c091c6`.
- Runtime-equivalent PGO/CI metadata repair: `44090994360fe6d441cb09c0d5b0262004c42f2e`.
- [Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38970557378) succeeded.
- [Repaired hosted CI](https://github.com/thesyncim/govpx/actions/runs/38971105373) is tracked separately.

Original 3792 CI rejected a stale PGO fingerprint after two comment-spacing
changes made after profiling. Reverting only those spaces in an independent
fingerprint calculation exactly reproduces the stored hash. 4409 fully
regenerates PGO and its fingerprint on final source and passes pre-commit;
it changes no Go runtime source. Benchmark builds disable PGO, so this
metadata issue does not affect the measurements below.

All variants share one EPYC 7763 VM (AVX2, no AVX512), official Go 1.27.2,
`GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`, and PGO disabled.
Six rounds cover every permutation once, placing each variant first,
second, and third twice. `orders.txt` retains the schedule. Both SIMD
variants assert native dispatch.

Only shared batch-DCT, SSE-strip, DCT32, in-place paired-DC, six-tap phases,
and split-subpixel test/benchmark fixtures are overlaid on reference
worktrees; production sources stay unchanged at their pinned commits.
Builds use the root project's Go/cache wrapper. Parent/current comparisons
share the same host. Comparing absolute times across separate runs does
not isolate code changes, even for hosts with the same CPU model.

| Public workload | Assembly | Parent SIMD | Current SIMD | Current vs parent | Current vs assembly |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP9 320x180 | 7.022 ms | 4.413 ms | 4.402 ms | -0.25% | -37.31% |
| VP8 720p | 23.452 ms | 23.716 ms | 23.047 ms | -2.82% | -1.73% |
| VP9 720p | 60.736 ms | 34.432 ms | 34.261 ms | -0.50% | -43.59% |

These are six-sample medians in ms per input frame. VP8 improves against
its parent in every pair (-4.45% to -2.08%), and takes 1.73% less median
time than assembly on this runner. Both VP9 paired ranges cross zero
(-1.38% to +0.34% at 320x180; -1.93% to +0.64% at 720p), so these samples
do not establish a VP9 gain from this checkpoint.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only mode,
and no external encoder comparison. VP9 320x180 uses 120 inputs at 500 kbps;
the 720p workloads use 60 inputs at 2,500 kbps.

| Kernel | Assembly | Parent SIMD | Current SIMD | Current vs parent |
| --- | ---: | ---: | ---: | ---: |
| Six-tap 4x4, horizontal phase 4 | 36.650 ns | 68.940 ns | 43.420 ns | -37.02% |
| Six-tap 4x4, two-axis phase 2/4 | 36.660 ns | 189.700 ns | 124.900 ns | -34.16% |
| Six-tap 8x8, two-axis phase 2/4 | 71.410 ns | 359.600 ns | 178.400 ns | -50.39% |
| Six-tap 16x16, horizontal phase 4 | 704.700 ns | 392.150 ns | 171.350 ns | -56.30% |
| Six-tap 16x16, two-axis phase 2/4 | 704.700 ns | 1,091.500 ns | 514.100 ns | -52.90% |
| Split subpixel 16x8 | 1,592.000 ns | 313.100 ns | 46.055 ns | -85.29% |
| Split subpixel 8x16 | 1,053.500 ns | 316.700 ns | 49.640 ns | -84.33% |
| Split subpixel 8x8 | 52.385 ns | 182.450 ns | 32.035 ns | -82.44% |
| Split subpixel 8x4 | 46.865 ns | 115.200 ns | 23.050 ns | -79.99% |
| Split subpixel 4x8 | 51.380 ns | 119.550 ns | 31.700 ns | -73.48% |
| Split subpixel 4x4 | 46.825 ns | 80.735 ns | 21.865 ns | -72.92% |

Every kernel has six 150-ms samples per variant, with zero bytes and
allocations. The 12 six-tap cases improve 11.65–58.85% against the parent;
some 4x4/8x8 phases still trail the original assembly substantially. The
split cases improve 72.92–85.29%; these microbenchmarks alone do not
establish a public encoder gain from the split change. Full results,
unchanged controls and paired samples are retained in raw logs and
`summary.json`, including remaining gaps against assembly.

Every output byte count, encoded/dropped count, quantizer statistic and
histogram, keyframe/interframe size, and allocation count matches across
all 18 samples per public workload, with no allocation outliers. Exact
bitstreams are asserted separately by codec oracle gates; benchmark JSON
reports retain byte counts. Full local native and purego suites pass at 3792.

| Workload | Output bytes | Encoded / dropped | Allocations per input |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Separate 120-frame CPU profiles run for every variant and both codecs
after all timing rounds. Profiles are excluded from timing medians.
Raw `.pgo` files and flat/cumulative top tables are retained.

Check out 3792 and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR HEAD^`
on native Linux/amd64 to reproduce. Omitting the last argument preserves
the two-variant comparison. Every raw artifact file is copied unchanged;
`summary.json` retains every sample, both comparisons, and parity fields.
