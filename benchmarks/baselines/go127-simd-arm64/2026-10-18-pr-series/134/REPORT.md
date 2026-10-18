# Reduce paired4x4 variance averages in one row pass

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3Variance4/keyfalse/x0` | 1.896x | 1.887 / 1.905 |
| `BenchmarkRound3Variance4/keyfalse/x60` | 0.978x | 0.979 / 0.976 |
| `BenchmarkRound3Variance4/keyfalse/x7` | 1.889x | 1.891 / 1.887 |
| `BenchmarkRound3Variance4/keytrue/x0` | 1.714x | 1.716 / 1.712 |
| `BenchmarkRound3Variance4/keytrue/x60` | 0.972x | 0.974 / 0.969 |
| `BenchmarkRound3Variance4/keytrue/x7` | 1.716x | 1.711 / 1.721 |

4 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The optimized route retains scalar fallback semantics for unsupported and overlapping windows; focused fixtures include the applicable numeric and alias boundaries.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
