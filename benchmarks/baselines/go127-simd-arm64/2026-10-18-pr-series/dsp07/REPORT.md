# Fuse split-block one-axis filtering and variance statistics

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkDSPPRSplit/4x4/4_0` | 3.080x | 3.065 / 3.096 |
| `BenchmarkDSPPRSplit/4x4/0_4` | 3.056x | 3.020 / 3.093 |
| `BenchmarkDSPPRSplit/4x4/3_5` | 1.011x | 1.006 / 1.017 |
| `BenchmarkDSPPRSplit/4x8/4_0` | 3.501x | 3.451 / 3.550 |
| `BenchmarkDSPPRSplit/4x8/0_4` | 3.385x | 3.328 / 3.442 |
| `BenchmarkDSPPRSplit/4x8/3_5` | 0.992x | 0.986 / 0.998 |
| `BenchmarkDSPPRSplit/8x4/4_0` | 4.454x | 4.408 / 4.500 |
| `BenchmarkDSPPRSplit/8x4/0_4` | 4.373x | 4.329 / 4.418 |
| `BenchmarkDSPPRSplit/8x4/3_5` | 1.004x | 1.003 / 1.005 |
| `BenchmarkDSPPRSplit/8x8/4_0` | 5.381x | 5.311 / 5.450 |
| `BenchmarkDSPPRSplit/8x8/0_4` | 5.304x | 5.216 / 5.391 |
| `BenchmarkDSPPRSplit/8x8/3_5` | 1.007x | 1.003 / 1.012 |
| `BenchmarkDSPPRSplit/8x16/4_0` | 5.168x | 5.176 / 5.159 |
| `BenchmarkDSPPRSplit/8x16/0_4` | 5.146x | 5.109 / 5.182 |
| `BenchmarkDSPPRSplit/8x16/3_5` | 1.007x | 1.009 / 1.005 |

10 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The optimized route retains scalar fallback semantics for unsupported and overlapping windows; focused fixtures include the applicable numeric and alias boundaries.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
