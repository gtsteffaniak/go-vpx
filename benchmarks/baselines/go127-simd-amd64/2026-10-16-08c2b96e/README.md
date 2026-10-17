# Native amd64 assembly/parent/current comparison, 2026-10-16

- Assembly-era baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
- Previous SIMD: `ef35b4bafd2822817b96c82815255f2097bfe661`.
- Current SIMD: `08c2b96ed68fd41e46d46112781e88d7408f7551`.
- [Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38968105288) succeeded.
- [Complete hosted CI](https://github.com/thesyncim/govpx/actions/runs/38968105289) also succeeded, including strict byte parity and the regression report.

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
| VP9 320x180 | 7.009 ms | 4.493 ms | 4.428 ms | -1.44% | -36.83% |
| VP8 720p | 23.568 ms | 23.808 ms | 23.908 ms | +0.42% | +1.44% |
| VP9 720p | 60.443 ms | 35.497 ms | 34.590 ms | -2.55% | -42.77% |

These are six-sample medians in ms per input frame. VP9 improves against
its parent in every pair: -1.92% to -1.31% at 320x180, and -4.74% to
-0.59% at 720p. The VP8 paired range crosses zero (-1.27% to +1.44%), so
these samples do not establish a VP8 gain or regression from this checkpoint.
VP8 remains 1.44% slower than assembly on this runner.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only
mode, and no external encoder comparison. VP9 320x180 uses 120 inputs
at 500 kbps; the 720p workloads use 60 inputs at 2,500 kbps.

| Kernel | Assembly | Parent SIMD | Current SIMD | Current vs parent |
| --- | ---: | ---: | ---: | ---: |
| VP9 subpixel variance, 8x8 control | 168.100 ns | 32.895 ns | 33.150 ns | +0.78% |
| VP9 subpixel variance, 16x16 | 207.950 ns | 85.560 ns | 73.310 ns | -14.32% |
| VP9 subpixel variance, 32x32 | 421.250 ns | 290.000 ns | 229.400 ns | -20.90% |
| VP9 subpixel variance, 64x64 | 1,190.500 ns | 1,115.500 ns | 880.400 ns | -21.08% |
| VP9 non-flat filter, horizontal 8 | 28.960 ns | 42.500 ns | 26.515 ns | -37.61% |
| VP9 non-flat filter, vertical 8 | 36.040 ns | 44.725 ns | 31.825 ns | -28.84% |
| VP9 non-flat filter, horizontal dual 8 | 63.255 ns | 83.210 ns | 50.935 ns | -38.79% |
| VP9 non-flat filter, vertical dual 8 | 75.345 ns | 87.690 ns | 61.610 ns | -29.74% |

Each kernel has six 150-ms samples per variant, with zero bytes and
allocations. Full results, unchanged controls, and samples are retained in
raw logs and `summary.json`; remaining kernel gaps versus assembly are
visible alongside the parent improvements. Non-flat filter fixtures target
the new branch; flat-filter behavior is covered separately by correctness
and ARM timing evidence, and public workloads include production dispatch.

Every output byte count, encoded/dropped count, quantizer statistic and
histogram, keyframe/interframe size, and allocation count matches across
all **18 samples** per public workload. There are no allocation outliers.
Exact bitstreams are asserted separately by codec oracle gates; benchmark
JSON reports retain byte counts. At this checkpoint, full local native/purego,
five-package race/checkptr=1, strict encoder byte parity, and the 18/18
parity report pass alongside source/build-selection/PGO checks and the complete
hosted production gates.

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
