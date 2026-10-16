# Native amd64 assembly/parent/current comparison, 2026-10-16

- Assembly-era baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
- Previous SIMD: `649237be2ee45d6580f31b6e7d52b9e3c0cef66a`.
- Current SIMD: `dff92311d0034a565166f9f84a5510e9dfb6f9fd`.
- [Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38974750401) and [complete hosted CI](https://github.com/thesyncim/govpx/actions/runs/38974750415) both succeeded.

All variants share one EPYC 7763 VM (AVX2, no AVX512), official Go 1.27.2,
`GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`, and PGO disabled.
Six rounds cover every permutation once, placing each variant first,
second, and third twice. `orders.txt` retains the schedule. Both SIMD
variants assert native dispatch.

Shared test-only fixtures cover batch DCT, SSE strips, DCT32, paired DC,
six-tap phases, split subpixel, column projections and quadrant averages.
These fixtures are overlaid on reference worktrees; production sources
stay unchanged at their pinned commits. Builds use the root project's
Go/cache wrapper. Parent/current comparisons share the same host.
Absolute times across separate runs do not isolate code changes, even for
hosts with the same CPU model.

| Public workload | Assembly | Parent SIMD | Current SIMD | Current vs parent | Current vs assembly |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP9 320x180 | 6.996 ms | 4.389 ms | 4.324 ms | -1.48% | -38.19% |
| VP8 720p | 23.122 ms | 22.213 ms | 22.295 ms | +0.37% | -3.58% |
| VP9 720p | 60.133 ms | 33.117 ms | 32.097 ms | -3.08% | -46.62% |

These are six-sample medians in ms per input frame. VP9 improves against
its parent in every pair: -2.12% to -0.98% at 320x180 and -3.43% to
-2.46% at 720p. VP8's parent/current paired range crosses zero (-0.06%
to +0.77%), so these samples do not establish a VP8 change. The public
gain belongs to the complete checkpoint; kernel speedups alone do not
isolate the contribution of individual changes.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only mode,
and no external encoder comparison. VP9 320x180 uses 120 inputs at 500 kbps;
the 720p workloads use 60 inputs at 2,500 kbps.

| Kernel | Assembly | Parent SIMD | Current SIMD | Current vs parent |
| --- | ---: | ---: | ---: | ---: |
| Four 8x8 quadrant averages | 69.260 ns | 108.850 ns | 20.055 ns | -81.58% |
| Column projection, 16x1 | 10.665 ns | 14.045 ns | 13.070 ns | -6.94% |
| Column projection, 16x7 | 73.675 ns | 58.995 ns | 25.345 ns | -57.04% |
| Column projection, 16x128 | 1,236.500 ns | 973.350 ns | 258.350 ns | -73.46% |
| Column projection, 32x128 | 2,462.000 ns | 1,284.500 ns | 371.600 ns | -71.07% |
| Column projection, 64x128 | 5,879.500 ns | 1,879.500 ns | 588.400 ns | -68.69% |
| Restored vertical internal luma edges | 117.000 ns | 128.250 ns | 124.450 ns | -2.96% |
| Unchanged horizontal internal-edge control | 35.220 ns | 44.785 ns | 43.295 ns | -3.33% |

The x64 vertical-filter path restores the previous generic SIMD body after
6492's ARM-tuned stores regressed it. The horizontal implementation is
unchanged; its measured variation is a control, not an additional algorithm
improvement attributed here. Short column cases remain in the measurements:
16x1 still trails assembly. All 63 kernels have six 150-ms samples per
variant, with zero bytes and allocations. Full results, controls and paired
samples are retained in raw logs and `summary.json`.

Every output byte count, encoded/dropped count, quantizer statistic and
histogram, and keyframe/interframe size matches across all 18 samples per
public workload. Exact bitstreams are asserted separately by codec oracle
gates; benchmark JSON reports retain byte counts.

Allocation medians match. Current VP9 720p round 2 reports one extra
allocation per 60-frame run (1.266667 allocations/input); all other samples,
including every assembly and parent sample, report 1.25. VP9 320x180
remains 0.616667 and VP8 remains zero in every sample. The single outlier
is retained rather than described as exact allocation parity.

| Workload | Output bytes | Encoded / dropped | Median allocations per input |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Full local native/purego coverage, five-package ARM/Rosetta parity,
race/checkptr=1, strict encoder byte parity, the 18/18 parity report, and
source/build-selection/fmt/PGO checks pass. The companion ARM reductions
folder preserves exact validation coverage and the instruction-region proof
for the x64-only scope fix; complete hosted production gates also pass.

Separate 120-frame CPU profiles run for every variant and both codecs
after all timing rounds. Profiles are excluded from timing medians.
Raw `.pgo` files and flat/cumulative top tables are retained.

Check out this checkpoint and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR HEAD^`
on native Linux/amd64 to reproduce. Omitting the last argument preserves
the two-variant comparison. Every raw artifact file is copied unchanged;
`summary.json` retains every sample, both comparisons, and parity fields.
