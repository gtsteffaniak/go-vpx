# Extend partial split source once before native SAD

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkMotionPRs/motion04/partial_16x8` | 1.623x | 1.701 / 1.545 |
| `BenchmarkMotionPRs/motion04/partial_8x16` | 1.739x | 1.976 / 1.503 |
| `BenchmarkMotionPRs/motion04/partial_8x8` | 1.768x | 1.934 / 1.601 |
| `BenchmarkMotionPRs/motion04/partial_4x4` | 1.038x | 1.101 / 0.975 |

3 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
