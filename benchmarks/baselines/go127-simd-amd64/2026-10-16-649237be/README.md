# Native amd64 assembly/parent/current comparison, 2026-10-16

- Assembly-era baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
- Previous SIMD: `44090994360fe6d441cb09c0d5b0262004c42f2e`.
- Current SIMD: `649237be2ee45d6580f31b6e7d52b9e3c0cef66a`.
- [Native benchmark](https://github.com/thesyncim/govpx/actions/runs/38972549766) and [complete hosted CI](https://github.com/thesyncim/govpx/actions/runs/38972548234) both succeeded.

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
| VP9 320x180 | 7.036 ms | 4.397 ms | 4.424 ms | +0.62% | -37.13% |
| VP8 720p | 23.402 ms | 23.002 ms | 22.862 ms | -0.61% | -2.31% |
| VP9 720p | 60.583 ms | 34.027 ms | 34.147 ms | +0.35% | -43.64% |

These are six-sample medians in ms per input frame. All parent/current
paired ranges cross zero: -0.07% to +1.79% at VP9 320x180, -1.20% to
+0.14% at VP8 720p, and -0.47% to +1.24% at VP9 720p. These samples do
not establish a public encoder gain from this checkpoint. Kernel results
below include both the horizontal gain and vertical regression on x64.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only mode,
and no external encoder comparison. VP9 320x180 uses 120 inputs at 500 kbps;
the 720p workloads use 60 inputs at 2,500 kbps.

| Kernel | Assembly | Parent SIMD | Current SIMD | Current vs parent |
| --- | ---: | ---: | ---: | ---: |
| Internal luma edges, horizontal | 35.420 ns | 51.540 ns | 44.795 ns | -13.09% |
| Internal luma edges, vertical | 117.100 ns | 124.550 ns | 128.100 ns | +2.85% |

The horizontal kernel improves in every pair (-13.16% to -12.80%), while
the vertical kernel regresses in every pair (+2.64% to +3.22%). This
architecture difference remains unresolved. Every kernel has six 150-ms
samples per variant, with zero bytes and allocations. Full results,
unchanged controls and paired samples are retained in raw logs and
`summary.json`, including remaining gaps against assembly.

Every output byte count, encoded/dropped count, quantizer statistic and
histogram, and keyframe/interframe size matches across all 18 samples per
public workload. Exact bitstreams are asserted separately by codec oracle
gates; benchmark JSON reports retain byte counts.

Allocation medians match, but two VP9 720p samples contain one extra
allocation per 60-frame run: assembly round 4 and current round 1 report
1.266667 allocations/input instead of the usual 1.25. All other samples,
including every parent sample, report 1.25. VP9 320x180 remains 0.616667
and VP8 remains zero in every sample. These outliers are retained rather
than described as exact allocation parity.

| Workload | Output bytes | Encoded / dropped | Median allocations per input |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Full local native/purego, five-package race/checkptr=1, strict encoder byte
parity, the 18/18 parity report, and source/build-selection/fmt/PGO checks
all pass alongside the complete hosted production gates. Exact local gate
logs are under the companion ARM `2026-10-16-internal-filter/` directory.

Separate 120-frame CPU profiles run for every variant and both codecs
after all timing rounds. Profiles are excluded from timing medians.
Raw `.pgo` files and flat/cumulative top tables are retained.

Check out this checkpoint and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR HEAD^`
on native Linux/amd64 to reproduce. Omitting the last argument preserves
the two-variant comparison. Every raw artifact file is copied unchanged;
`summary.json` retains every sample, both comparisons, and parity fields.
