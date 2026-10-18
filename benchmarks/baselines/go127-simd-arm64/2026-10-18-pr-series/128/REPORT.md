# Reuse coefficient probability decisions during header emission

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3CoefUpdates/p0/m0` | 0.951x | 0.970 / 0.931 |
| `BenchmarkRound3CoefUpdates/p0/m1` | 1.003x | 1.003 / 1.002 |
| `BenchmarkRound3CoefUpdates/p1/m0` | 1.826x | 1.827 / 1.825 |
| `BenchmarkRound3CoefUpdates/p1/m1` | 1.001x | 1.000 / 1.001 |
| `BenchmarkRound3CoefUpdates/p2/m0` | 1.865x | 1.867 / 1.864 |
| `BenchmarkRound3CoefUpdates/p2/m1` | 1.000x | 1.001 / 0.998 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The optimized route retains scalar fallback semantics for unsupported and overlapping windows; focused fixtures include the applicable numeric and alias boundaries.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
