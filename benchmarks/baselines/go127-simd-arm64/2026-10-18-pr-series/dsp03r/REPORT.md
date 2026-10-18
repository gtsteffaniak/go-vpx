# Skip zero-DC in-place reconstruction with a small dispatch predicate

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkDSPPRZeroDC/dc0/inplacetrue` | 3.367x | 3.355 / 3.379 |
| `BenchmarkDSPPRZeroDC/dc0/inplacefalse` | 0.812x | 0.806 / 0.819 |
| `BenchmarkDSPPRZeroDC/dc8/inplacetrue` | 0.914x | 0.907 / 0.920 |
| `BenchmarkDSPPRZeroDC/dc8/inplacefalse` | 0.886x | 0.878 / 0.893 |

1 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The optimized route retains scalar fallback semantics for unsupported and overlapping windows; focused fixtures include the applicable numeric and alias boundaries.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
