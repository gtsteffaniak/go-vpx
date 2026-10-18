# Collapse identity bilinear axes and SIMD the small predictors

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkDSPPRBilinear/w4/0_0` | 2.527x | 2.533 / 2.521 |
| `BenchmarkDSPPRBilinear/w4/4_0` | 2.490x | 2.525 / 2.456 |
| `BenchmarkDSPPRBilinear/w4/0_4` | 2.522x | 2.531 / 2.512 |
| `BenchmarkDSPPRBilinear/w4/3_5` | 1.255x | 1.255 / 1.255 |
| `BenchmarkDSPPRBilinear/w8/0_0` | 3.889x | 3.976 / 3.801 |
| `BenchmarkDSPPRBilinear/w8/4_0` | 3.976x | 4.017 / 3.934 |
| `BenchmarkDSPPRBilinear/w8/0_4` | 3.911x | 3.980 / 3.841 |
| `BenchmarkDSPPRBilinear/w8/3_5` | 1.986x | 1.936 / 2.036 |

8 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The optimized route retains scalar fallback semantics for unsupported and overlapping windows; focused fixtures include the applicable numeric and alias boundaries.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
