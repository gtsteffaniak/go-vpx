# Fifty qualifying performance PRs: merged

All **50 PRs #57–#106** are merged after their current fast CI passed. Final main is `715645fd212cfc65295f4f6fb1002dcc60856758`. Every independently accepted PR has at least one representative native component case above 1.15x in both balanced short pairs (the initial two PRs also exceed this threshold). Native hosts and all controls/losses remain scoped in the individual reports. No compounded or whole-codec gain is inferred from kernel ratios.

PR69's sole merge conflict was an import declaration in dequant.go. The composition retains both vector dequant multiplication and sparse DC-only reconstruction; affected DSP and root packet checks passed before its green merge. All branches remain retained. No README or handwritten-assembly changes were introduced by this series.

## Combined focused checks and short codec screen

The 56 exact test names from the published fixtures and root streams pass on the merged tree, covering five affected packages. This is a targeted composition check, not a full suite. All 96 VP8 packet/flag hashes and 72 VP9 ML packets remain unchanged.

Apple M4 Max, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Two short samples per reference/workload; balanced three-reference order. Throughput changes below are preliminary and contain mixed-sign controls:

| Workload | Versus batch parent98 | Versus fixed c4 baseline |
|---|---:|---:|
| vp8-720 | +0.398% | +0.559% |
| vp9-320 | +1.636% | +4.063% |
| vp9-720 | -0.758% | +6.845% |

All 22 nonperformance fields and allocation rates agree in every comparison. The fixed-baseline **1.5x whole-codec gate is not met**. Full CI, matrices, new C comparisons, profiles and broad audits remain deferred. No additional timing is needed to establish the already retained kernel wins.

## Libvpx gap: historical, not remeasured after this series

The last completed x64 comparison measured exact c4 on AMD EPYC 9V74 against optimized libvpx 1.16.0 /1024874c5919305883187e2953de8fcb4c3d7fa6. Go/C time ratios were VP8 adaptive1.60385x, VP8 fixed-speed2.16491x, VP9 320p2.59273x and VP9 720p2.35105x. Only fixed-speed VP8 matched measured output and decoded quality. Higher ratios mean libvpx is faster. Those results predate the 50 PRs and cannot be updated using M4 kernel/codec ratios or another EPYC model.

The earlier M4 comparison at ARM48 matched output/quality for both VP8 scopes: adaptive1.53588x and fixed-speed2.04583x Go/C time. VP9 was unequal output/quality work. These older results are reference context, not a current merged-main gap claim.

The raw package preserves remote merge receipts, combined checks, CLI source/binary bindings, every short sample and the unchanged historical C proof. Retained research archives keep additional local kernel wins, rejected candidates and limitations. The extra VP9 DSP source is pending integration with explicit scratch-alias/malformed-input caveats; it is not shipped code.
