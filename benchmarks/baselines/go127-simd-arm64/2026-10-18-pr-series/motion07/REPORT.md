# Reuse extended source block throughout split subpel search

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkMotionPR07/partial_8x8` | 1.820x | 1.728 / 1.911 |
| `BenchmarkMotionPR07/partial_4x4` | 1.626x | 1.658 / 1.594 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The optimized route retains scalar fallback semantics for unsupported and overlapping windows; focused fixtures include the applicable numeric and alias boundaries.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
